package screenreader

import (
	"errors"
	"image"

	"github.com/nejstastnejsistene/dotbot-go/solver"
)

type Grid struct {
	topLeft image.Point // The center of the top left dot.
	dist    int         // The distance between centers of adjacent dots.
}

func ReadScreen(img image.Image) (board solver.Board, err error) {
	grid, err := FindGrid(img)
	if err != nil {
		return
	}
	for c := 0; c < solver.BoardSize; c++ {
		for r := 0; r < solver.BoardSize; r++ {
			x := grid.topLeft.X + c*grid.dist
			y := grid.topLeft.Y + r*grid.dist
			col := img.At(x, y)
			if isBackground(col) {
				err = errors.New("screenreader: expecting dot here")
				return
			}
			board.SetColor(r, c, HueToColor(Hue(col)))
		}
	}
	return
}

func FindGrid(img image.Image) (grid Grid, err error) {
	grid.topLeft, err = findTopLeft(img)
	if err != nil {
		return
	}
	// Go from the center of the top left dot to its rightmost edge.
	xMax := img.Bounds().Max.X / 3
	var x0 = grid.topLeft.X
	for ; x0 < xMax; x0++ {
		if isBackground(img.At(x0, grid.topLeft.Y)) {
			break
		}
	}
	// Go to the leftmost edge of the second dot.
	for ; x0 < xMax; x0++ {
		if !isBackground(img.At(x0, grid.topLeft.Y)) {
			break
		}
	}
	// Go to the rightmost edge of the second dot.
	var x = x0
	for ; x < xMax; x++ {
		if isBackground(img.At(x, grid.topLeft.Y)) {
			break
		}
	}
	if x == xMax {
		err = errors.New("screenreader: can't find distance between dots")
	} else {
		// Subtract the center of the second dot from the first.
		grid.dist = (x0+x)/2 - grid.topLeft.X
	}
	return
}

func findTopLeft(img image.Image) (p image.Point, err error) {
	xMax := img.Bounds().Max.X / 4
	yMax := img.Bounds().Max.Y / 3
	for x0 := 0; x0 < xMax; x0++ {
		for y := 0; y < yMax; y++ {
			// Find the leftmost edge.
			if !isBackground(img.At(x0, y)) {
				// Find the rightmost edge.
				var x int
				for x = x0; x < xMax; x++ {
					if isBackground(img.At(x, y)) {
						break
					}
				}
				// Average the two for the x coordinate.
				p.X = (x0 + x) / 2
				// Find the topmost edge.
				var y0 int
				for y0 = y; y0 > 0; y0-- {
					if isBackground(img.At(p.X, y0)) {
						break
					}
				}
				// Find the bottommost edge.
				for ; y < yMax; y++ {
					if isBackground(img.At(p.X, y)) {
						break
					}
				}
				// Average the two for the y coordinate.
				p.Y = (y0 + y) / 2
				return
			}
		}
	}
	err = errors.New("screenreader: unable to find top left dot")
	return
}
