package main

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/nejstastnejsistene/dotbot-go/screencap"
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

	f, err := os.Create("/data/local/DotBot/screencap.png")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))
}
