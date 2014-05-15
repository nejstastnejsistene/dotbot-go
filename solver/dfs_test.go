package solver

import (
	"fmt"
	"testing"
)

func TestDFS(t *testing.T) {
	fmt.Println("DFS")
	mask := Mask(0)
	mask.Add(5, 2)
	mask.Add(5, 3)
	mask.Add(5, 4)
	mask.Add(4, 3)
	fmt.Println(mask)
	paths := make(chan Mask)
	go mask.DFS(paths)
	for path := range paths {
		fmt.Println(path)
	}
}
