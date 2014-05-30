package solver

// #cgo CFLAGS: -O3
// #include "mask.h"
import "C"
import (
	"fmt"
	"regexp"
	"strings"
)

type Mask C.Mask

const (
	NumDots = C.NumDots
	AllDots = Mask(C.AllDots)
)

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

func DotMask(row, col int) Mask {
	return Mask(1) << index(row, col)
}

func (mask Mask) Matches(pattern Mask) bool {
	return (mask & pattern) == pattern
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

func (mask Mask) Count() int {
	return int(C.Count(C.Mask(mask)))
}

func (mask Mask) CountNeighbors(row, col int) int {
	return int(C.CountNeighbors(C.Mask(mask), C.int(row), C.int(col)))
}

func (mask Mask) Partition(partitions chan Mask) {
	q := newQueue()
	defer q.free()
	C.Partition(C.Mask(mask), (*C.Queue)(q))
	for _, partition := range q.slice() {
		partitions <- partition
	}
	close(partitions)
}

func (mask Mask) DFS(paths chan Mask) {
	q := newQueue()
	defer q.free()
	C.DFS(C.Mask(mask), (*C.Queue)(q))
	for _, path := range q.slice() {
		paths <- path
	}
	close(paths)
}

// Function for creating a mask from a string, for testing.
// Highly panicky. Can handle leading whitespace, as long as
// its identical for all lines starting with the first line
// with content. Every second character can be an X to indicate
// that it is set. Everything else should be a space. Example:
//
//	s := `
//	X X X
//	X   X X
//	X X X X X
//	    X   X
//	    X X X`
//
func maskFromString(s string) (mask Mask) {
	// Skip leading empty lines.
	lines := strings.Split(s, "\n")
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	if len(lines) == 0 {
		panic("no lines")
	}
	// Find leading whitespace for the first line.
	// We'll assume all lines have it.
	p := regexp.MustCompile(`^[\s]*`)
	n := len(p.Find([]byte(lines[0])))
	for row, line := range lines {
		for col, char := range line[n:] {
			switch {
			case !InBounds(row, col/2):
				panic("out of bounds")
			case col%2 == 0 && char == 'X':
				mask.Add(row, col/2)
			case char != ' ':
				panic(fmt.Sprintf("unexpected character: %#+v", string(char)))
			}
		}
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
