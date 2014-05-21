package main

import (
	"fmt"
	"time"

	"github.com/nejstastnejsistene/dotbot-go/screencap"
	"github.com/nejstastnejsistene/dotbot-go/screenreader"
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

	fmt.Println(screenreader.ReadScreen(img))

	fmt.Println(time.Since(start))
}
