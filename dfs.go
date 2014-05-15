package dotbot

func (mask Mask) DFS(paths chan Mask) {
	seen := make(map[Mask]bool)
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if mask.Contains(row, col) && mask.CountNeighbors(row, col) == 1 {
				mask.dfsHelper(paths, seen, row, col, Mask(0))
			}
		}
	}
	close(paths)
}

func (mask Mask) dfsHelper(paths chan Mask,
	seen map[Mask]bool, row, col int, path Mask) {

	visit := func(row, col int) {
		if mask.Contains(row, col) {
			mask.dfsHelper(paths, seen, row, col, path)
		}
	}

	mask.Remove(row, col)
	path.Add(row, col)
	if !seen[path] {
		seen[path] = true
		paths <- path
	}

	visit(row-1, col)
	visit(row+1, col)
	visit(row, col-1)
	visit(row, col+1)
}
