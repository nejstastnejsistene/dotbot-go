package solver

// #cgo CFLAGS: -O3
// #include "cycles.h"
// #include "mask.h"
import "C"

// A database of all possible non-square cycles. The slice of
// cycles at db[rows][cols] is the list of all cycles of size
// rows x cols that doesn't contain a square. The logic behind
// this is that any cycle that has a square in it doesn't add
// anything to that same cycle with the square removed. E.g.
//
// X X X X                                           X X X
// X   X X will have the same affect on the board as X   X
// X X X X                                           X X X
//
// This is because any cycle will remove all dots of that color,
// so the only real differentiating factor between cycles is how
// many and which dots they encircle. This database only has data
// for sizes 3x3 to 6x6, because anything less than that either
// contains a square or no cycles at all. It is populated in init().
var db [BoardSize + 1][BoardSize + 1][]Mask

// Find all of the cycles in this mask. It will only yield
// unique masks, with respect to the effect that this mask
// would have to the board. It uses colorMask to help determine
// this. This means that if there are multiple squares, it will
// only return one of them, because all squares have the same
// effect on the board.
func (mask Mask) Cycles(cycles chan Mask, colorMask Mask) {
	q := newQueue()
	defer q.free()
	C.Cycles(C.Mask(mask), C.Mask(colorMask), (*C.Queue)(q))
	for _, cycle := range q.slice() {
		cycles <- cycle
	}
	close(cycles)
}

// A suare in the upper left corner of the board.
const Square = C.Square

// Returns the first square in this mask. Returns 0 if there is none.
func (mask Mask) findSquare(r0, c0, r1, c1 int) Mask {
	return Mask(C.findSquare(C.Mask(mask),
		C.int(r0), C.int(c0), C.int(r1), C.int(c1)))
}

// Returns the perimeter of a rows x cols cycle.
func perimeter(rows, cols int) int {
	return 2*(rows+cols) - 4
}

// Return the minimum and maximum rows and columns. In other words,
// return the convex hull, where (r0, c0), (r1, c1) are the
// top left and bottom right coordinates. Convex hull is a
// fancy math term to indicate the smallest convex set that
// contains all of a set of points.
// https://en.wikipedia.org/wiki/Convex_hull
func (mask Mask) ConvexHull() (int, int, int, int) {
	var r0, c0, r1, c1 C.int
	C.ConvexHull(C.Mask(mask), &r0, &c0, &r1, &c1)
	return int(r0), int(c0), int(r1), int(c1)
}

// Returns a Mask of any dots encircled by this cyclic Mask.
// Using this on masks that are cyclic or contain squares will have
// undefined/meaningless return values.
func (mask Mask) Encircled() Mask {
	return Mask(C.Encircled(C.Mask(mask)))
}

func init() {
	C.init()
	// Mirror db.
	for rows := 3; rows <= BoardSize; rows++ {
		for cols := 3; cols <= BoardSize; cols++ {
			q := (*queue)(C.db[rows][cols])
			db[rows][cols] = q.slice()
		}
	}
}
