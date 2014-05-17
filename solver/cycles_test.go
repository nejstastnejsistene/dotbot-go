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

func TestFindSquare(t *testing.T) {
}

func TestEncircled(t *testing.T) {
	var mask, encircled Mask
	assert := func() {
		if mask.Encircled() != encircled {
			t.Errorf("Mask:\n%v", mask)
			t.Errorf("Expected encircled:\n%v", encircled)
			t.Errorf("Actual encircled:\n%v", mask.Encircled())
			t.Errorf("%#+v\n", mask)
		}
	}

	mask = maskFromString(`
	X X X
	X   X
	X   X`)
	encircled = 0
	assert()

	mask = maskFromString(`
	X   X
	X   X
	X X X
	`)
	encircled = 0
	assert()

	mask = maskFromString(`
	X X X
	X    
	X X X
	`)
	encircled = 0
	assert()

	mask = maskFromString(`
	X X X
	    X
	X X X
	`)
	encircled = 0
	assert()

	mask = maskFromString(`
	X X X
	X   X
	X X X`)
	encircled = 0
	encircled.Add(1, 1)
	assert()

	mask = maskFromString(`
	X X X X X X
	X     X   X
	X X   X X X
	  X X X X
	    X   X
	    X X X`)
	encircled = 0
	encircled.Add(1, 1)
	encircled.Add(1, 2)
	encircled.Add(2, 2)
	encircled.Add(1, 4)
	encircled.Add(4, 3)
	assert()

	mask = maskFromString(`
	X X X X X
	X   X   X X
	X X X X   X
	  X   X X X
	  X X X`)
	encircled = 0
	encircled.Add(1, 1)
	encircled.Add(1, 3)
	encircled.Add(2, 4)
	encircled.Add(3, 2)
	assert()
}

// TODO: Think of a way to test the completeness of the generated cycles.
func TestDB(t *testing.T) {
	for rows := 0; rows < len(db); rows++ {
		for cols := 0; cols < len(db[rows]); cols++ {
			for _, cycle := range db[rows][cols] {

				// Make sure there are the expected number of dots.
				diff := perimeter(rows, cols) - cycle.Count()
				if !(diff == 0 || diff == 1) {
					t.Errorf("Incorrect number of dots.")
				}

				// Make sure an appropriate number of dots are encircled.
				n := cycle.Encircled().Count()
				switch {
				case n == 0:
					t.Errorf("encircled == %d\n%v", n, cycle)
				case rows == 3 && cols == 3:
					if n != 1 {
						t.Errorf("encircled != %d\n%v", n, cycle)
					}
				default:
					lowerLimit := 2
					if rows == 6 || cols == 6 {
						lowerLimit = 3
					}
					upperLimit := (rows - 2) * (cols - 2)
					if !(lowerLimit <= n && n <= upperLimit) {
						msgFmt := "!(%d <= encircled <= %d)\n%v"
						t.Errorf(msgFmt, lowerLimit, upperLimit, cycle)
					}
				}

				// Make sure there are no squares!
				if cycle.findSquare(0, 0, rows-1, cols-1) != 0 {
					t.Error("There shouldn't be any squares here.")
				}
			}
		}
	}
}
