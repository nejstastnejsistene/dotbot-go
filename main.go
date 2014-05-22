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
	for {
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

			start = time.Now()
			t, err := touchscreen.New("/dev/input/event0", img.Bounds())
			if err != nil {
				panic(err)
			}
			fmt.Println("touchscreen.New():", time.Since(start))

			start = time.Now()
			board, err := screenreader.ReadScreen(img)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("screenreader.ReadScreen():", time.Since(start))

			fmt.Println(board)

			start = time.Now()
			path := board.ChooseMove(-1).ConstructPath()
			fmt.Println("screenreader.ReadScreen():", time.Since(start))

			if len(path) == 1 {
				r := path[0].Row
				c := path[0].Col
				start = time.Now()
				err = t.Tap(image.Point{
					215 + 154*c, 525 + 154*r,
				})
				fmt.Println("touchscreen.Tap():", time.Since(start))
				start = time.Now()
				err = t.Tap(image.Point{
					215 + 154*c, 525 + 154*r,
				})
				fmt.Println("touchscreen.Tap():", time.Since(start))
			} else {
				start = time.Now()
				points := make([]image.Point, len(path))
				for i, p := range path {
					r := p.Row
					c := p.Col
					points[i] = image.Point{215 + 154*c, 525 + 154*r}
				}
				t.Gesture(points)
				fmt.Println("touchscreen.Gesture():", time.Since(start))
			}
		}()

		time.Sleep(750 * time.Millisecond)
	}
}
