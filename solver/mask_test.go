package solver

import (
	"math/rand"
	"testing"
)

func RandMask() Mask {
	return Mask(rand.Int63n(int64(AllDots + 1)))
}

func TestInBounds(t *testing.T) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if !InBounds(row, col) {
				t.Error("Unexpected out of bounds", row, col)
			}
		}
	}
	for i := 0; i < 1000; i++ {
		row := rand.Int() + BoardSize
		col := rand.Int() + BoardSize
		if InBounds(row, col) || InBounds(-row, -col) {
			t.Error("Expected to be out of bounds:", row, col)
		}
	}
}

func TestMatches(t *testing.T) {
	for i := 0; i < 1000; i++ {
		original := RandMask()
		mask := original
		n := rand.Intn(NumDots)
		for j := 0; j < n; j++ {
			row := rand.Intn(BoardSize)
			col := rand.Intn(BoardSize)
			mask.Add(row, col)
		}
		if !mask.Matches(original) {
			t.Fatal("Masks do not match")
		}
	}
}

func TestAddRemove(t *testing.T) {
	for i := 0; i < 1000; i++ {
		mask := Mask(0)
		row := rand.Intn(BoardSize)
		col := rand.Intn(BoardSize)
		mask.Add(row, col)
		if !mask.Contains(row, col) {
			t.Fatal("Mask should contain", row, col)
		}
		mask.Remove(row, col)
		if mask.Contains(row, col) {
			t.Fatal("Mask shouldn't contain", row, col)
		}
	}
}

func TestCount(t *testing.T) {
	for i := 0; i < 1000; i++ {
		points := rand.Perm(NumDots)
		count := rand.Intn(NumDots)
		mask := Mask(0)
		for j := 0; j < count; j++ {
			row, col := unIndex(uint(points[j]))
			mask.Add(row, col)
		}
		if mask.Count() != count {
			t.Fatal("Count() is incorrect")
		}
	}
}

func TestPartition(t *testing.T) {
	for i := 0; i < 1000; i++ {
		partitions := make(chan Mask)
		mask := RandMask()
		go mask.Partition(partitions)
		count := 0
		for p := range partitions {
			count += p.Count()
			for row := 0; row < BoardSize; row++ {
				for col := 0; col < BoardSize; col++ {
					if p.Contains(row, col) && !mask.Contains(row, col) {
						t.Fatal("Partition contains dot not in original")
					}
				}
			}
		}
		if count != mask.Count() {
			t.Fatal("Total number of dots is incorrect")
		}
	}
}
