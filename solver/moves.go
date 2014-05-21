package solver

const (
	MaxDepth    = 3
	Cutoff      = NumDots / 2
	Decay       = 0.5
	CycleWeight = 1 / Decay
)

type Move struct {
	Path   Mask
	Color  Color
	Cyclic bool
}

type weightedMove struct {
	weight float64
	depth  int
	move   Move
}

// Execute a move against the board.
func (board *Board) MakeMove(move Move) (score int) {
	// Mark the dots in the path, and the dots that they
	// encircle to be shrunk.
	dots := move.Path
	if move.Cyclic {
		dots |= dots.Encircled()
	}
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			// Shrink the marked dots, and also any dots of the same
			// color if it's a cycle.
			if dots.Contains(row, col) ||
				(move.Cyclic && board.Color(row, col) == move.Color) {
				board.Shrink(row, col)
				score++
			}
		}
	}
	return
}

func (board Board) ChooseMove(movesRemaining int) Move {
	moves := make(chan Move)
	go board.Moves(moves)
	maxDepth := movesRemaining
	if maxDepth <= 0 || maxDepth > MaxDepth {
		maxDepth = MaxDepth
	}
	return board.chooseMove(moves, 0, 1, maxDepth).move
}

func (board Board) chooseMove(moves chan Move, numEmpty, depth, maxDepth int) (chosen weightedMove) {
	for move := range moves {
		// Don't consider shrinkers after the first round. There
		// are too many possibilities and they are generally not
		// particularly high scoring, so they waste a lot of time.
		if depth > 1 && move.Path.Count() == 1 {
			continue
		}
		// Apply the move to a copy of the board.
		newBoard := board
		score := newBoard.MakeMove(move)
		// Initialize the weight to the score of the move.
		weight := float64(score)
		deepest := depth
		// Give weight to cycles to account for the decreased
		// entropy in the dots that will be filled in.
		if move.Cyclic {
			weight *= CycleWeight
		}
		// If the bounds haven't been reached, recur.
		if numEmpty < Cutoff && depth < maxDepth {
			newMoves := make(chan Move)
			go newBoard.Moves(newMoves)
			result := newBoard.chooseMove(newMoves, numEmpty+score, depth+1, maxDepth)
			weight += Decay * result.weight
			deepest = result.depth
		}
		// At the first level, how deep it had to look to reach
		// the cutoff factors into the weight. As an extreme example,
		// using 36 shrinkers over 36 turns is a lot less valuable
		// than scoring 36 in a single turn.
		if depth == 1 {
			weight /= float64(deepest)
		}
		// Update the maximally weighted move.
		if weight > chosen.weight {
			chosen = weightedMove{weight, deepest, move}
		}
	}
	return
}

func (board Board) Moves(moves chan Move) {
	// Iterate through colors.
	for _, color := range Colors {
		colorMask := board.ColorMask(color)
		// Iterate through partitions.
		partitions := make(chan Mask)
		go colorMask.Partition(partitions)
		for partition := range partitions {
			// Iterate through cycles, if any.
			cycles := make(chan Mask)
			hadCycles := false
			go partition.Cycles(cycles, colorMask)
			for cycle := range cycles {
				moves <- Move{cycle, color, true}
				hadCycles = true
			}
			// Do a DFS if there weren't any cycles.
			if !hadCycles {
				paths := make(chan Mask)
				go partition.DFS(paths)
				for path := range paths {
					moves <- Move{path, color, false}
				}
			}
		}
	}
	close(moves)
}
