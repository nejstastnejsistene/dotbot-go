package solver

// #cgo CFLAGS: -O3
// #include "moves.h"
import "C"

type Move struct {
	cMove  C.Move
	Path   Mask
	Color  Color
	Cyclic bool
}

func NewMove(path Mask, color Color, cyclic bool) Move {
	mColor := Mask(color)
	mCyclic := Mask(0)
	if cyclic {
		mCyclic = 1
	}
	cMove := C.Move(path)
	cMove |= C.Move(mColor << C.COLOR_SHIFT)
	cMove |= C.Move(mCyclic << C.CYCLIC_SHIFT)
	return Move{cMove, path, color, cyclic}
}

func (board Board) ChooseMove(turnsRemaining int) Move {
	move := C.ChooseMove(&board[0], C.int(turnsRemaining))
	return Move{
		move,
		Mask(move) & AllDots,
		Color(move >> C.COLOR_SHIFT),
		((move >> C.CYCLIC_SHIFT) & 1) == 1,
	}
}

type Point struct{ Row, Col int }

func (move Move) ConstructPath() []Point {
	q := (*queue)(C.ConstructPath(move.cMove))
	defer q.free()
	cPoints := q.slice()

	points := make([]Point, len(cPoints))
	for i, point := range cPoints {
		points[i] = Point{
			int(point) % BoardSize,
			int(point) / BoardSize,
		}
	}
	return points
}
