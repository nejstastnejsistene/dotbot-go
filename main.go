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

func main() {
	start := time.Now()

	img, err := screencap.NewScreenCap()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := img.Close(); err != nil {
			panic(err)
		}
	}()

	board, err := screenreader.ReadScreen(img)
	fmt.Println(board)

	t, err := touchscreen.New("/dev/input/event0", img.Bounds())
	if err != nil {
		fmt.Println(err)
	} else {
		for r := 0; r < solver.BoardSize; r++ {
			for c := 0; c < solver.BoardSize; c++ {
				err = t.Tap(image.Point{
					215 + 154*c, 525 + 154*r,
				})
				if err != nil {
					panic(err)
				}
			}
		}
	}

	fmt.Println(time.Since(start))
}
