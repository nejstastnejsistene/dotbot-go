package dotbot

type Mask uint64

const AllDots = (Mask(1) << (BoardSize * BoardSize)) - 1

func InBounds(row, col int) bool {
	return 0 <= row && row < BoardSize && 0 <= col && col < BoardSize
}

func index(row, col int) uint {
	if !InBounds(row, col) {
		panic("Index out of bounds")
	}
	return uint(BoardSize)*uint(col) + uint(row)
}

func (mask Mask) Matches(pattern Mask) bool {
	return (mask & pattern) == pattern
}

func DotMask(row, col int) Mask {
	return Mask(1) << index(row, col)
}

func (mask Mask) Contains(row, col int) bool {
	return mask.Matches(DotMask(row, col))
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
		if InBounds(row, col) && mask.Contains(row, col) {
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

// Internal data structures for buildPartition().

type point struct{ row, col int }

type stack struct {
	data [BoardSize * BoardSize]point
	size int
}

func (s *stack) push(row, col int) {
	s.data[s.size] = point{row, col}
	s.size++
}

func (s *stack) pop() (row, col int) {
	s.size--
	row = s.data[s.size].row
	col = s.data[s.size].col
	return
}
