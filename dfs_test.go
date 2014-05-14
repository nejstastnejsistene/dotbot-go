package dotbot

import (
	"fmt"
	"testing"
)

func TestDFS(t *testing.T) {
	mask := RandomBoard().ColorMask(Red)
	ch := make(chan Mask)
	go mask.Partition(ch)
	<-ch
	<-ch
	<-ch
	<-ch
	p := <-ch
	//fmt.Println(p, p.Contains(4, 4))
	paths := make(chan Mask)
	go p.DFS(paths)
	for x := range paths {
		fmt.Println(x)
	}
}
