package dotbot

var db [BoardSize + 1][BoardSize + 1][]Mask

func init() {
	for rows := 6; rows <= BoardSize; rows++ {
		for cols := 6; cols <= BoardSize; cols++ {

			cycles := make(chan Mask)
			go func(rows, cols int) {
				computeCyclesOfSize(cycles, rows, cols)
				close(cycles)
			}(rows, cols)

			defer func(rows, cols int) {
				db[rows][cols] = make([]Mask, 0, 100)
				for cycle := range cycles {
					db[rows][cols] = append(db[rows][cols], cycle)
				}
			}(rows, cols)
		}
	}

}

const Square = Mask(3 | 3<<BoardSize)

func perimeter(rows, cols int) int {
	return 2*(rows+cols) - 4
}

func (mask Mask) convexHull() (r0, c0, r1, c1 int) {
	r0 = BoardSize
	c0 = BoardSize
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if mask.Contains(row, col) {
				if row < r0 {
					r0 = row
				}
				if row > r1 {
					r1 = row
				}
				if col < c0 {
					c0 = col
				}
				if col > c1 {
					c1 = col
				}
			}
		}
	}
	return
}

func (mask Mask) findPattern(pattern Mask, cycles chan Mask) {
	r0, c0, r1, c1 := mask.convexHull()
	for r := r0; r < r1; r++ {
		for c := c0; c < c1; c++ {
			cycle := pattern << index(r, c)
			if mask.Matches(cycle) {
				cycles <- cycle
			}
		}
	}
}

// Returns a Mask of any dots encircled by this cyclic Mask.
// Using this on masks that are cyclic or contain squares will have
// undefined/meaningless return values.
func (mask Mask) Encircled() Mask {

	// This works by filling in the dots from each of the four directions,
	// outwards-in. Any dots left unfilled are encircled. This approach
	// would not work for board larger that 6x6 because it would become
	// possible to create concave cyclic paths.
	for r := 0; r < BoardSize; r++ {
		for c := 0; c < BoardSize && !mask.Contains(r, c); c++ {
			mask.Add(r, c)
		}
		for c := BoardSize - 1; c >= 0 && !mask.Contains(r, c); c-- {
			mask.Add(r, c)
		}
	}
	for c := 0; c < BoardSize; c++ {
		for r := 0; r < BoardSize && !mask.Contains(r, c); r++ {
			mask.Add(r, c)
		}
		for r := BoardSize - 1; r >= 0 && !mask.Contains(r, c); r-- {
			mask.Add(r, c)
		}
	}
	return ^mask
}

func computeCyclesOfSize(cycles chan Mask, rows, cols int) {

	var buildCandidateCycles func(chan Mask, Mask, int, int, int)
	buildCandidateCycles = func(cycles chan Mask, cycle Mask, col, prevStart, prevEnd int) {
		if col == cols {
			cycles <- cycle
			return
		}

		for start := 0; start < rows-2; start++ {
			for end := start + 2; end < rows; end++ {
				newCycle := cycle
				newCycle.Add(start, col)
				newCycle.Add(end, col)

				if col == 0 || col == cols-1 {
					for row := start + 1; row < end; row++ {
						newCycle.Add(row, col)
					}
				}

				if col > 0 {
					for row := start + 1; row < prevStart+1; row++ {
						newCycle.Add(row, col)
					}
					for row := prevStart + 1; row < start+1; row++ {
						newCycle.Add(row, col-1)
					}
					for row := prevEnd; row < end; row++ {
						newCycle.Add(row, col)
					}
					for row := end; row < prevEnd; row++ {
						newCycle.Add(row, col-1)
					}
				}

				if col+1 == cols {
					cycles <- newCycle
				} else {
					buildCandidateCycles(cycles, newCycle, col+1, start, end)
				}
			}
		}
	}

	hasSquare := func(mask Mask) bool {
		for r := 0; r < rows-1; r++ {
			for c := 0; c < cols-1; c++ {
				if mask.Matches(Square << index(r, c)) {
					return true
				}
			}
		}
		return false
	}

	isNonSquareCycle := func(cycle Mask) bool {
		if cycle.Count() != perimeter(rows, cols) {
			return false
		}
		if hasSquare(cycle) {
			return false
		}
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				if cycle.Contains(row, col) {
					if cycle.CountNeighbors(row, col) < 2 {
						return false
					}
				}
			}
		}
		return true
	}

	candidates := make(chan Mask)
	go func() {
		buildCandidateCycles(candidates, Mask(0), 0, 0, 0)
		close(cycles)
	}()

	for cycle := range candidates {
		if isNonSquareCycle(cycle) {
			cycles <- cycle
		}
	}
}
