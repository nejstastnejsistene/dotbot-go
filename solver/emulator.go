package solver

func PlayGame(numTurns int) (scores []int, totalScore int) {
	scores = make([]int, numTurns)
	board := RandomBoard()
	for i := 0; i < numTurns; i++ {
		move := board.ChooseMove(numTurns - i)
		scores[i] = board.MakeMove(move)
		totalScore += scores[i]
		if move.Cyclic {
			board.FillEmptyExcluding(move.Color)
		} else {
			board.FillEmpty()
		}
	}
	return
}
