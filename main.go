package main

import (
	"fmt"
	"image"
	"time"

	"github.com/nejstastnejsistene/dotbot-go/screencap"
	"github.com/nejstastnejsistene/dotbot-go/screenreader"
	"github.com/nejstastnejsistene/dotbot-go/touchscreen"
)

func main() {
	// Open the touchscreen.
	start := time.Now()
	t, err := touchscreen.New("/dev/input/event0", nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := t.Device.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println("touchscreen.New():", time.Since(start))

	for {
		// Take a screenshot.
		func() {
			start := time.Now()
			img, err := screencap.NewScreenCap()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer func() {
				if err := img.Close(); err != nil {
					fmt.Println(err)
				}
			}()
			fmt.Println("screencap.NewScreenCap():", time.Since(start))

			// Tell the touchscreen how large the screen is.
			if t.PixelSize == nil {
				t.SetPixelSize(img.Bounds())
			}

			// Read the screen.
			start = time.Now()
			board, err := screenreader.ReadScreen(img)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("screenreader.ReadScreen():", time.Since(start))

			fmt.Println(board)

			// Choose a move.
			start = time.Now()
			path := board.ChooseMove(-1).ConstructPath()
			fmt.Println("board.ChooseMove():", time.Since(start))

			// Play the move.
			if len(path) == 1 {
				p := screenreader.Grid.Coordinate(path[0])
				start = time.Now()
				err = t.Tap(p)
				fmt.Println("touchscreen.Tap():", time.Since(start))
				start = time.Now()
				err = t.Tap(p)
				fmt.Println("touchscreen.Tap():", time.Since(start))
			} else {
				start = time.Now()
				points := make([]image.Point, len(path))
				for i, p := range path {
					points[i] = screenreader.Grid.Coordinate(p)
				}
				t.Gesture(points)
				fmt.Println("touchscreen.Gesture():", time.Since(start))
			}
		}()

		// Give time for new dots to fall in.
		time.Sleep(750 * time.Millisecond)
	}
}
