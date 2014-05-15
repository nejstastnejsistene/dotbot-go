package solver

import (
	"fmt"
	"testing"
)

func TestCycles(t *testing.T) {
	mask := Mask(0)
	mask.Add(0, 0)
	mask.Add(1, 0)
	mask.Add(2, 0)
	mask.Add(0, 1)
	mask.Add(2, 1)
	mask.Add(0, 2)
	mask.Add(1, 2)
	mask.Add(2, 2)

	mask.Add(3, 2)
	mask.Add(4, 2)
	mask.Add(2, 3)
	mask.Add(2, 4)
	mask.Add(3, 4)
	mask.Add(4, 4)
	mask.Add(4, 3)

	mask.Add(3, 5)
	mask.Add(4, 5)

	fmt.Println(mask)

	cycles := make(chan Mask)
	go mask.Cycles(cycles, mask)
	for cycle := range cycles {
		fmt.Println("Cycle:")
		fmt.Println(cycle)
	}
}
