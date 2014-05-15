package dotbot

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
				partitions <- mask.buildPartition(row, col)
			}
		}
	}
	close(partitions)
}

func (mask *Mask) buildPartition(row, col int) (partition Mask) {
	stack := new(stack)

	visit := func(row, col int) {
		if mask.Contains(row, col) {
			mask.Remove(row, col)
			stack.push(row, col)
		}
	}

	visit(row, col)

	for stack.size > 0 {
		row, col = stack.pop()
		partition.Add(row, col)

		visit(row-1, col)
		visit(row+1, col)
		visit(row, col-1)
		visit(row, col+1)
	}
	return
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

// Internal data structure for buildPartition().

type stack struct {
	data [NumDots]uint
	size int
}

func (s *stack) push(row, col int) {
	s.data[s.size] = index(row, col)
	s.size++
}

func (s *stack) pop() (row, col int) {
	s.size--
	r, c := unIndex(s.data[s.size])
	return int(r), int(c)
}
