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

func RandomBoard() Board {
	var board Board
	board.FillEmptyWithRandom()
	return board
}

func randColor() Color {
	return Color(rand.Intn(int(Purple-Red+1))) + Red
}

func (board *Board) FillEmptyWithRandom() {
	board.FillEmptyWithRandomExcluding(Empty)
}

func (board *Board) FillEmptyWithRandomExcluding(exclude Color) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if board.Color(row, col) == Empty {
				color := randColor()
				for color == exclude {
					color = randColor()
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
