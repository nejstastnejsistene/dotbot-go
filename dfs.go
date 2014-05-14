package dotbot

func (mask Mask) DFS(paths chan Mask) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if mask.Contains(row, col) {
				mask.dfsHelper(paths, row, col, Mask(0))
			}
		}
	}
	close(paths)
}

func (mask Mask) dfsHelper(paths chan Mask, row, col int, path Mask) {

	visit := func(row, col int) {
		if InBounds(row, col) && mask.Contains(row, col) {
			mask.dfsHelper(paths, row, col, path)
		}
	}

	mask.Remove(row, col)
	path.Add(row, col)

	visit(row-1, col)
	visit(row+1, col)
	visit(row, col-1)
	visit(row, col+1)

	paths <- path
}
