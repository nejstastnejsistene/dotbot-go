package dotbot

func (mask Mask) DFS(paths chan Mask) {
	filter := new(uniqFilter)
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if mask.Contains(row, col) && mask.countNeighbors(row, col) == 1 {
				mask.uniqPaths(paths, row, col, filter)
			}
		}
	}
	close(paths)
}

func (mask Mask) countNeighbors(row, col int) int {
	count := 0
	if mask.Contains(row-1, col) {
		count++
	}
	if mask.Contains(row+1, col) {
		count++
	}
	if mask.Contains(row, col-1) {
		count++
	}
	if mask.Contains(row, col+1) {
		count++
	}
	return count
}

func (mask Mask) uniqPaths(paths chan Mask, row, col int, filter *uniqFilter) {
	dfsResults := make(chan dfsResult)
	go mask.dfsHelper(dfsResults, row, col, Mask(0))

	for result := range dfsResults {
		if filter.check(row, col, result.endRow, result.endCol) {
			paths <- result.path
		}
	}
}

func (mask Mask) dfsHelper(results chan dfsResult, row, col int, path Mask) {

	visit := func(row, col int) {
		if mask.Contains(row, col) {
			mask.dfsHelper(results, row, col, path)
		}
	}

	mask.Remove(row, col)
	path.Add(row, col)
	results <- dfsResult{path, row, col}

	visit(row-1, col)
	visit(row+1, col)
	visit(row, col-1)
	visit(row, col+1)

	if path.Count() == 1 {
		close(results)
	}
}

// Internal data structures.

type dfsResult struct {
	path           Mask
	endRow, endCol int
}

type uniqFilter [NumDots][NumDots]bool

func (u *uniqFilter) check(startRow, startCol, endRow, endCol int) bool {
	i := index(startRow, startCol)
	j := index(endRow, endCol)
	exists := u[i][j] || u[j][i]
	if exists {
		return false
	} else {
		u[i][j] = true
		u[j][i] = true
		return true
	}
}
