package solver

// #cgo CFLAGS: -O3
// #include "board.h"
// #include "mask.h"
import "C"
import "fmt"

const BoardSize = C.BoardSize

type Color C.Color
type Board C.Board

const (
	Empty    Color = C.Empty
	NotEmpty       = C.NotEmpty
	Red            = C.Red
	Yellow         = C.Yellow
	Green          = C.Green
	Blue           = C.Blue
	Purple         = C.Purple
)

var Colors = []Color{
	Red,
	Yellow,
	Green,
	Blue,
	Purple,
}

func (c Color) String() string {
	switch c {
	case Empty:
		return "Empty"
	case NotEmpty:
		return "NotEmpty"
	case Red:
		return "Red"
	case Yellow:
		return "Yellow"
	case Green:
		return "Green"
	case Blue:
		return "Blue"
	case Purple:
		return "Purple"
	}
	panic(fmt.Sprintf("Unknown color: %d", c))
	return ""
}

func RandomBoard() (board Board) {
	C.FillEmpty(&board[0])
	return
}

func RandomColor() Color {
	return Color(C.RandomColor())
}

func (board *Board) FillEmpty() {
	board.FillEmptyExcluding(Empty)
}

func (board *Board) FillEmptyExcluding(exclude Color) {
	C.FillEmptyExcluding(&board[0], C.Color(exclude))
}

func (board Board) Color(row, col int) Color {
	return Color(board[col][row])
}

func (board *Board) SetColor(row, col int, color Color) {
	board[col][row] = C.Color(color)
}

func (board *Board) Shrink(row, col int) {
	C.Shrink(&board[0], C.int(row), C.int(col))
}

func (board Board) ColorMask(color Color) Mask {
	return Mask(C.ColorMask(&board[0], C.Color(color)))
}

const dotFmt = " \x1b[%dm\u25cf\x1b[0m"

var colorCodes = map[Color]int{
	NotEmpty: 0,
	Red:      31,
	Yellow:   33,
	Green:    32,
	Blue:     36,
	Purple:   35,
}

func (board Board) String() (s string) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			color := board.Color(row, col)
			if color == Empty {
				s += "  "
			} else {
				s += fmt.Sprintf(dotFmt, colorCodes[color])
			}
		}
		s += "\n"
	}
	return
}
