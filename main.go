package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/nejstastnejsistene/dotbot-go/screencap"
)

func main() {
	start := time.Now()
	f, err := os.Open("screencap.raw")
	if err != nil {
		panic(err)
	}
	img, err := screencap.Decode(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(img.(*image.RGBA).Pix))
	fmt.Println(time.Since(start))
}
