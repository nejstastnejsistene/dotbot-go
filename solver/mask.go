package solver

type Mask uint64

const NumDots = BoardSize * BoardSize
const AllDots = (Mask(1) << NumDots) - 1

func InBounds(row, col int) bool {
	return 0 <= row && row < BoardSize && 0 <= col && col < BoardSize
}

func index(row, col int) uint {
	if !InBounds(row, col) {
		panic("Index out of bounds")
	}
	return uint(BoardSize)*uint(col) + uint(row)
}

func unIndex(index uint) (row, col int) {
	return int(index % uint(BoardSize)), int(index / uint(BoardSize))
}

func (mask Mask) Matches(pattern Mask) bool {
	return (mask & pattern) == pattern
}

func DotMask(row, col int) Mask {
	return Mask(1) << index(row, col)
}

func (mask Mask) Contains(row, col int) bool {
	return InBounds(row, col) && mask.Matches(DotMask(row, col))
}

func (mask *Mask) Add(row, col int) {
	*mask |= DotMask(row, col)
}

func (mask *Mask) Remove(row, col int) {
	*mask &= ^DotMask(row, col)
}

func (mask Mask) Count() (count int) {
	for mask != 0 {
		mask ^= (mask & -mask)
		count++
	}
	return
}

func (mask Mask) CountNeighbors(row, col int) int {
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

func (mask Mask) Partition(partitions chan Mask) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if mask.Contains(row, col) {
				partitions <- mask.buildPartition(Mask(0), row, col)
			}
		}
	}
	close(partitions)
}

func (mask *Mask) buildPartition(p Mask, row, col int) Mask {

	visit := func(p Mask, row, col int) Mask {
		if mask.Contains(row, col) {
			p = mask.buildPartition(p, row, col)
		}
		return p
	}

	mask.Remove(row, col)
	p.Add(row, col)

	p = visit(p, row-1, col)
	p = visit(p, row+1, col)
	p = visit(p, row, col-1)
	p = visit(p, row, col+1)
	return p
}

func (mask Mask) DFS(paths chan Mask) {
	seen := make(map[Mask]bool)
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if mask.Contains(row, col) && mask.CountNeighbors(row, col) == 1 {
				mask.buildPaths(paths, seen, row, col, Mask(0))
			}
		}
	}
	close(paths)
}

func (mask Mask) buildPaths(paths chan Mask,
	seen map[Mask]bool, row, col int, path Mask) {

	visit := func(row, col int) {
		if mask.Contains(row, col) {
			mask.buildPaths(paths, seen, row, col, path)
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

func (mask Mask) String() string {
	board := new(Board)
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if mask.Contains(row, col) {
				board.SetColor(row, col, NotEmpty)
			}
		}
	}
	return board.String()
}
