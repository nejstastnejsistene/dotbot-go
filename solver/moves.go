package solver

const (
	MaxDepth    = 3
	Cutoff      = NumDots / 2
	Decay       = 0.5
	CycleWeight = 1 / Decay
)

type Move struct {
	path   Mask
	Color  Color
	Cyclic bool
}

type weightedMove struct {
	weight float64
	depth  int
	move   Move
}

func (board *Board) MakeMove(move Move) (score int) {
	dots := move.path
	if move.Cyclic {
		dots |= dots.Encircled()
	}
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
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
	if maxDepth > MaxDepth {
		maxDepth = MaxDepth
	}
	return board.chooseMove(moves, 0, 1, maxDepth).move
}

func (board Board) chooseMove(moves chan Move, numEmpty, depth, maxDepth int) (chosen weightedMove) {
	for move := range moves {
		if depth > 1 && move.path.Count() == 1 {
			continue
		}
		//
		newBoard := board
		score := newBoard.MakeMove(move)
		//
		weight := float64(score)
		deepest := depth
		//
		if move.Cyclic {
			weight *= CycleWeight
		}
		//
		if numEmpty < Cutoff && depth < maxDepth {
			newMoves := make(chan Move)
			go newBoard.Moves(newMoves)
			result := newBoard.chooseMove(newMoves, numEmpty+score, depth+1, maxDepth)
			weight += Decay * result.weight
			deepest = result.depth
		}
		if depth == 1 {
			weight /= float64(deepest)
		}
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
