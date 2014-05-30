package solver

// #cgo CFLAGS: -O3
// #include "cycles.h"
// #include "mask.h"
import "C"

var db [BoardSize + 1][BoardSize + 1][]Mask

func (mask Mask) Cycles(cycles chan Mask, colorMask Mask) {
	q := newQueue()
	defer q.free()
	C.Cycles(C.Mask(mask), C.Mask(colorMask), (*C.Queue)(q))
	for _, cycle := range q.slice() {
		cycles <- cycle
	}
	close(cycles)
}

const Square = C.Square

func (mask Mask) findSquare(r0, c0, r1, c1 int) Mask {
	return Mask(C.findSquare(C.Mask(mask),
		C.int(r0), C.int(c0), C.int(r1), C.int(c1)))
}

func perimeter(rows, cols int) int {
	return 2*(rows+cols) - 4
}

func (mask Mask) ConvexHull() (int, int, int, int) {
	var r0, c0, r1, c1 C.int
	C.ConvexHull(C.Mask(mask), &r0, &c0, &r1, &c1)
	return int(r0), int(c0), int(r1), int(c1)
}

func (mask Mask) Encircled() Mask {
	return Mask(C.Encircled(C.Mask(mask)))
}

func init() {
	C.init()
}
