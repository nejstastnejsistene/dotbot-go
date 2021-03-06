package solver

import (
	"math/rand"
	"testing"
)

func TestRandomColor(t *testing.T) {
	seen := make(map[Color]bool)
	for i := 0; i < 1000; i++ {
		seen[RandomColor()] = true
	}
	for _, color := range Colors {
		if !seen[color] {
			t.Error("Color not represented:", color)
		}
	}
}

func TestFillEmptyExcluding(t *testing.T) {
	for i := 0; i < 1000; i++ {
		for _, color := range Colors {
			board := new(Board)
			board.FillEmptyExcluding(color)
			for row := 0; row < BoardSize; row++ {
				for col := 0; col < BoardSize; col++ {
					switch board.Color(row, col) {
					case Empty, NotEmpty:
						t.Fatal("Empty and NotEmpty are not valid colors")
					case color:
						t.Fatalf("%s should have been excluded", color)
					}
				}
			}
		}
	}
}

func TestSetColor(t *testing.T) {
	board := new(Board)
	for i := 0; i < 1000; i++ {
		row := rand.Intn(BoardSize)
		col := rand.Intn(BoardSize)
		color := RandomColor()
		board.SetColor(row, col, color)
		if board.Color(row, col) != color {
			t.Fatal("SetColor() and Color() do not match up")
		}
	}
}

func TestShrink(t *testing.T) {
	for i := 0; i < 1000; i++ {
		board := RandomBoard()
		copy := board

		row := rand.Intn(BoardSize)
		col := rand.Intn(BoardSize)

		board.Shrink(row, col)
		if board.Color(0, col) != Empty {
			t.Fatal("Top dot should be empty")
		}
		for r := row; r > 0; r-- {
			if board.Color(r, col) != copy.Color(r-1, col) {
				t.Errorf("Shrink(%v, %v)", row, col)
				t.Fatalf("Dot should have fallen:\n%v\n%v", copy, board)
			}
		}
	}
}

func TestColorMask(t *testing.T) {
	for i := 0; i < 1000; i++ {
		board := RandomBoard()
		count := 0
		allColors := Mask(0)
		for _, color := range Colors {
			mask := board.ColorMask(color)
			count += mask.Count()
			allColors |= mask
			for row := 0; row < BoardSize; row++ {
				for col := 0; col < BoardSize; col++ {
					boardTrue := board.Color(row, col) == color
					maskTrue := mask.Contains(row, col)
					if boardTrue != maskTrue {
						t.Fatal("Color mask doesn't match board")
					}
				}
			}
		}
		if allColors != AllDots || count != NumDots {
			t.Fatal("The color masks don't add up to AllDots!")
		}
	}
}
