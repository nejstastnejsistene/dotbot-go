package dotbot

import (
	"math/rand"
	"testing"
)

func TestDFS(t *testing.T) {
}

func TestUniqFilter(t *testing.T) {
	filter := new(uniqFilter)
	for i := 0; i < 1000; i++ {
		r0 := rand.Intn(BoardSize)
		c0 := rand.Intn(BoardSize)
		r1 := rand.Intn(BoardSize)
		c1 := rand.Intn(BoardSize)
		x := index(r0, c0)
		y := index(r1, c1)
		if filter[x][y] != filter[x][y] {
			t.Fatal("uniqFilter in inconsistent state")
		}
		exists := filter[x][y] || filter[y][x]
		isUnique := filter.check(r0, c0, r1, c1)
		if exists == isUnique {
			t.Fatal("uniqFilter let duplicate through")
		}
	}
}
