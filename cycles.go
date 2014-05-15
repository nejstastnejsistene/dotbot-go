package dotbot

// A database of all possible non-square cycles. The slice of
// cycles at db[rows][cols] is the list of all cycles of size
// rows x cols that doesn't contain a square. The logic behind
// this is that any cycle that has a square in it doesn't add
// anything to that same cycle with the square removed. E.g.
//
// X X X X											 X X X
// X   X X will have the same affect on the board as X   X
// X X X X 											 X X X
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
	// Mark the end of cycles whenever this returns.
	defer func() { close(cycles) }()

	if !colorMask.Matches(mask) {
		panic("mask is not contained within colorMask")
	}

	// Cycles have at least 4 dots.
	numDots := mask.Count()
	if numDots < 4 {
		return
	}

	r0, c0, r1, c1 := mask.ConvexHull()

	// Cycles are at least 2x2.
	numRows := r1 - r0 + 1
	numCols := c1 - c0 + 1
	if numRows < 2 || numCols < 2 {
		return
	}

	seen := make(map[Mask]bool)

	findPattern := func(pattern Mask) {
		// Translate this pattern throughout the convex hull.
		for r := r0; r < r1; r++ {
			for c := c0; c < c1; c++ {
				cycle := pattern << index(r, c)
				if mask.Matches(cycle) {
					// Calculate the resulting board.
					result := colorMask | cycle.Encircled()
					// Yield cycles which haven't been seen yet.
					if !seen[result] {
						seen[result] = true
						cycles <- cycle
					}
				}
			}
		}
	}

	// Compare all cycles from size 3x3 to numRows x numCols.
	for rows := 3; rows <= numRows; rows++ {
		for cols := 3; cols <= numCols; cols++ {
			// This prevents us from checking cycles that we don't
			// have enough dots to form. The perimeter is the number
			// of dots in any cycle, unless it crosses over itself.
			// In that case, there will be one less because a dot
			// is server as two corners. Example:
			//
			// X X X		This cycle is 5x5, so the perimeter
			// X   X		is 16. That center dot with four
			// X X X X X	neighbors rather than the typical two
			//     X   X	brings the actual number of dots to 15.
			//     X X X
			if numDots >= perimeter(rows, cols)-1 {
				for _, cycle := range db[rows][cols] {
					findPattern(cycle)
				}
			}
		}
	}

	// If there is a square, yield the first one found.
	square := mask.findSquare(r0, c0, r1, c1)
	if square != 0 {
		cycles <- square
	}
}

const Square = Mask(3 | 3<<BoardSize)

// Returns the first square in this mask. Returns 0 if there is none.
func (mask Mask) findSquare(r0, c0, r1, c1 int) Mask {
	for r := r0; r < r1; r++ {
		for c := c0; c < c1; c++ {
			square := Square << index(r, c)
			if mask.Matches(Square << index(r, c)) {
				return square
			}
		}
	}
	return 0
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
func (mask Mask) ConvexHull() (r0, c0, r1, c1 int) {
	if mask == 0 {
		return 0, 0, 0, 0
	}
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

func init() {
	for rows := 3; rows <= BoardSize; rows++ {
		for cols := 3; cols <= BoardSize; cols++ {
			// Start goroutines to compute all of the cycles
			// of dimension rows x cols.
			cycles := make(chan Mask)
			go func(rows, cols int) {
				computeCyclesOfSize(cycles, rows, cols)
				close(cycles)
			}(rows, cols)
			// As the results become available, fill the database
			// with these computed cycles.
			defer func(rows, cols int) {
				db[rows][cols] = make([]Mask, 0)
				for cycle := range cycles {
					db[rows][cols] = append(db[rows][cols], cycle)
				}
			}(rows, cols)
		}
	}
}

// Compute all of the cycles of dimension rows x cols.
// This involves generating potential cycles and then filtering
// them as described above.
func computeCyclesOfSize(cycles chan Mask, rows, cols int) {

	// Build cycles from left to right, column by column. It does that
	// by recognizing that all convex cycles without squares can be
	// represented by the top and bottom dot for each column, or in
	// other words, the start and end row for each column. For example,
	//
	// X X X X X X		X X X X X X
	//            		X		  X
	// X X        	=>	X X X     X
	//     X      			X X   X
	//		 X X X			  X X X
	//
	// This algorithm recursively selects these start and end points
	// and fills in the appropriate dots. This relatively small set
	// of potential cycles is small enough that it can be filtered
	// in a reasonable amount of time (about 1/3 seconds on my machine).
	var buildCandidateCycles func(chan Mask, Mask, int, int, int)
	buildCandidateCycles = func(cycles chan Mask,
		cycle Mask, col, prevStart, prevEnd int) {

		// Go through all the possible pairs of starts and ends. The starts
		// begin at 0 and can go until 2 less than the maximum size, go give
		// room for the end. The ends go from 2 more than the start until
		// the maximum size. The significance of this buffer of 2, is that it
		// is the smallest area that you can fit a corner into a cycle without
		// folding upon itself and creating a square.
		for start := 0; start < rows-2; start++ {
			for end := start + 2; end < rows; end++ {

				// Make a copy of the cycle and add the start and end points.
				newCycle := cycle
				newCycle.Add(start, col)
				newCycle.Add(end, col)

				// For the first and last columsn, connect the dots
				// between the start and end rows.
				if col == 0 || col == cols-1 {
					for row := start + 1; row < end; row++ {
						newCycle.Add(row, col)
					}
				}

				// This forms corners between the current and previous
				// start and end. For whichever of the starts is highest,
				// dots are placed below it in its column up until the other
				// start. For the ends, it is the same but from bottom to
				// top instead.
				//
				// Examples:
				//
				// X		X		  X		  X
				//   X		X X		X 		X X
				//		=>				=>
				// X		X X		  X		X X
				//   X		  X		X		X
				//
				if col > 0 {
					for row := start + 1; row <= prevStart; row++ {
						newCycle.Add(row, col)
					}
					for row := prevStart + 1; row <= start; row++ {
						newCycle.Add(row, col-1)
					}
					for row := prevEnd; row < end; row++ {
						newCycle.Add(row, col)
					}
					for row := end; row < prevEnd; row++ {
						newCycle.Add(row, col-1)
					}
				}

				// Yield the generated cycle if on the last column, otherwise
				// recur to the next column.
				if col+1 == cols {
					cycles <- newCycle
				} else {
					buildCandidateCycles(cycles, newCycle, col+1, start, end)
				}
			}
		}
	}

	// Check if a cycle fits our criterea for a unique cycle with no squares.
	isValidCycle := func(cycle Mask) bool {
		if cycle.findSquare(0, 0, rows, cols) != 0 {
			return false
		}
		// Keep track of the number of dots it should have.
		numDots := perimeter(rows, cols)
		// Make sure each dot has two or more neighbors.
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				if cycle.Contains(row, col) {
					n := cycle.CountNeighbors(row, col)
					if n == 4 {
						// Four neighbors means the cycle
						// crossed over itself, and needs one
						// less corner dot.
						numDots--
					} else if n < 2 {
						// Zero neighbors means it's in isolation.
						// One neighbor means it's an endpoint and
						// is not necessary.
						return false
					}
				}
			}
		}
		// Make sure it has the expected number of dots.
		if cycle.Count() != numDots {
			return false
		}
		return true
	}

	// Generate the potential cycles.
	candidates := make(chan Mask)
	go func() {
		buildCandidateCycles(candidates, Mask(0), 0, 0, 0)
		close(candidates)
	}()

	// Filter out the invalid candidates.
	for cycle := range candidates {
		if isValidCycle(cycle) {
			cycles <- cycle
		}
	}
}
