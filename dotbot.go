package main

import (
	"fmt"
	"image"
	"time"

	"github.com/nejstastnejsistene/dotbot-go/screencap"
	"github.com/nejstastnejsistene/dotbot-go/screenreader"
	"github.com/nejstastnejsistene/dotbot-go/solver"
	"github.com/nejstastnejsistene/dotbot-go/touchscreen"
)

var TS *touchscreen.TouchScreen

func main() {
	defer TS.Device.Close()
	for {
		PlayOneTurn()
		// Give time for new dots to fall in.
		time.Sleep(750 * time.Millisecond)
	}
}

func PlayOneTurn() (err error) {
	// Release the touchscreen in case interrupted mid event.
	defer TS.FingerUp()

	img, err := screencap.NewScreenCap()
	if err != nil {
		return
	}
	defer func() { err = img.Close() }()

	// Read the screen.
	board, err := screenreader.ReadScreen(img)
	if err != nil {
		return
	}

	fmt.Println(board)

	// Play a move.
	err = MakeMove(board.ChooseMove(-1))
	if err != nil {
		fmt.Println(err)
	}
	return
}

func MakeMove(move solver.Move) (err error) {
	path := move.ConstructPath()
	if len(path) == 1 {
		p := screenreader.Grid.Coordinate(path[0])
		err = TS.Tap(p)
		err = TS.Tap(p)
	} else {
		points := make([]image.Point, len(path))
		for i, p := range path {
			points[i] = screenreader.Grid.Coordinate(p)
		}
		err = TS.Gesture(points)
	}
	return
}

func init() {
	// Take a screencap to get the screen's size.
	img, err := screencap.NewScreenCap()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := img.Close(); err != nil {
			panic(err)
		}
	}()
	// Open the touchscreen device.
	TS, err = touchscreen.New("/dev/input/event0", img.Bounds())
	if err != nil {
		TS.Device.Close()
		panic(err)
	}
}
