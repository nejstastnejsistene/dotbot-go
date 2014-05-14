package dotbot

import (
	"fmt"
	"math/rand"
)

const BoardSize = 6

type Color int
type Board [BoardSize][BoardSize]Color

const (
	Empty Color = iota
	NotEmpty
	Red
	Yellow
	Green
	Blue
	Purple
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
	default:
		panic(fmt.Sprintf("Unknown color: %d", c))
	}
}

func RandomBoard() Board {
	var board Board
	board.FillEmptyWithRandom()
	return board
}

func RandomColor() Color {
	return Colors[rand.Intn(len(Colors))]
}

func (board *Board) FillEmptyWithRandom() {
	board.FillEmptyWithRandomExcluding(Empty)
}

func (board *Board) FillEmptyWithRandomExcluding(exclude Color) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if board.Color(row, col) == Empty {
				color := RandomColor()
				for color == exclude {
					color = RandomColor()
				}
				board.SetColor(row, col, color)
			}
		}
	}
}

func (board Board) Color(row, col int) Color {
	return board[col][row]
}

func (board *Board) SetColor(row, col int, color Color) {
	board[col][row] = color
}

func (board *Board) Shrink(row, col int) {
	for row > 0 {
		board.SetColor(row, col, board.Color(row-1, col))
		row--
	}
	board.SetColor(0, col, Empty)
}

func (board Board) ColorMask(color Color) Mask {
	mask := Mask(0)
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if board.Color(row, col) == color {
				mask.Add(row, col)
			}
		}
	}
	return mask
}

func (board Board) Copy() (copy Board) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			copy.SetColor(row, col, board.Color(row, col))
		}
	}
	return
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
