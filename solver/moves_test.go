package solver

import (
	"fmt"
	"testing"
)

func TestMakeMove(t *testing.T) {
}

func TestChooseMove(t *testing.T) {
}

func TestMoves(t *testing.T) {
	for i := 0; i < 1000; i++ {
		board := RandomBoard()
		moves := make(chan Move)
		go board.Moves(moves)
		for move := range moves {
			if move.Path == 0 {
				t.Fatal("path should not be empty")
			}
			if move.Color == Empty || move.Color == NotEmpty {
				t.Fatal("color should be a real color")
			}
			cycles := make(chan Mask)
			go move.Path.Cycles(cycles, move.Path)
			if _, ok := <-cycles; move.Cyclic != ok {
				t.Fatal("move.Cyclic is incorrect")
			}
		}
	}
}

// This hugely purple board triggered a bug where
// I wasn't properly calculating the number of empty dots,
// and then the cutoff wasn't working right so the huge
// purple area wasn't being weighted correctly. This test
// makes sure ChooseMove() selects one of the many obvious
// purple cycles available.
func TestPurple(t *testing.T) {
	board := Board{
		{6, 5, 3, 6, 6, 6},
		{6, 6, 5, 3, 5, 6},
		{6, 5, 4, 6, 6, 6},
		{6, 6, 5, 6, 6, 6},
		{6, 6, 6, 6, 6, 6},
		{5, 4, 5, 6, 6, 6},
	}
	move := board.ChooseMove(-1)
	if !(move.Cyclic && move.Color == Purple) {
		t.Fatalf("missed obvious cycles here:\n%v", board)
	}
}

// Similar problem to TestPurple. Although the cutoff was
// a real bug, it was simply causing the algorithm to avoid
// this bug. The depth of the default chosen move was zero
// rather than the current depth.
func TestPurple2(t *testing.T) {
	board := Board{
		{5, 4, 5, 2, 5, 2},
		{4, 6, 6, 6, 6, 5},
		{4, 5, 4, 5, 2, 6},
		{4, 6, 6, 6, 6, 6},
		{5, 4, 5, 6, 6, 6},
		{4, 3, 6, 2, 6, 6},
	}
	move := board.ChooseMove(-1)
	if !(move.Cyclic && move.Color == Purple) {
		t.Fatalf("missed obvious cycles here:\n%v", board)
	}
}

// This test case caught a but in chooseMove which forgot to
// initialize a struct with zero values.
func TestYellow(t *testing.T) {
	fmt.Println("TestFoo")
	board := Board{
		{3, 2, 4, 4, 4, 3},
		{2, 3, 6, 4, 3, 3},
		{5, 4, 3, 3, 6, 6},
		{4, 2, 2, 3, 3, 3},
		{3, 6, 4, 3, 6, 3},
		{3, 6, 3, 4, 3, 3},
	}
	move := board.ChooseMove(-1)
	board.MakeMove(move)
	move = board.ChooseMove(-1)
	if !(move.Cyclic && move.Color == Yellow) {
		t.Fatal("missed available yellow cycles on 2nd turn")
	}
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
				move := NewMove(cycle, Red, true)
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
			t.Fatalf("length of path should be the number of dots:\n%v", move.Path)
		}
	}
}
