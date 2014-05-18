package main

import (
	"fmt"
	"time"

	"github.com/nejstastnejsistene/dotbot-go/screencap"
)

func main() {
	start := time.Now()
	img, err := screencap.Mmap("screencap.raw")
	if err != nil {
		panic(err)
	} else {
		defer func() {
			if err := img.Close(); err != nil {
				panic(err)
			}
		}()
	}
	fmt.Println(time.Since(start))
}
