package solver

import "testing"

func TestMakeMove(t *testing.T) {
}

func TestChooseMove(t *testing.T) {
}

func TestMoves(t *testing.T) {
}

func TestPathOnCycles(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()
	for rows := 3; rows < BoardSize; rows++ {
		for cols := 3; cols < BoardSize; cols++ {
			for _, cycle := range db[rows][cols] {
				move := Move{cycle, Red, true}
				path := move.ConstructPath()
				if len(path) != move.Path.Count()+1 {
					t.Fatalf("length of path should be numDots + 1:\n%v", cycle)
				}
				if path[0] != path[len(path)-1] {
					t.Fatalf("path does not connect to itself")
				}
			}
		}
	}
}

func TestPathOnRandomBoards(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()
	for i := 0; i < 1000; i++ {
		board := RandomBoard()
		move := board.ChooseMove(-1)
		path := move.ConstructPath()
		if move.Cyclic {
			if len(path) != move.Path.Count()+1 {
				t.Fatalf("length of path should be numDots + 1:\n%v", move.Path)
			}
			if path[0] != path[len(path)-1] {
				t.Fatalf("path does not connect to itself")
			}
		} else if len(path) != move.Path.Count() {
			t.Fatalf("length of path should be the number of dots\n%v", move.Path)
		}
	}
}
