package screencap

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"syscall"
)

type ScreenCap interface {
	image.Image
	io.Closer
}

type RGBA_8888 struct {
	data []byte
	image.RGBA
}

func (img RGBA_8888) ColorModel() color.Model {
	return img.RGBA.ColorModel()
}

func (img RGBA_8888) Bounds() image.Rectangle {
	return img.RGBA.Bounds()
}

func (img RGBA_8888) At(x, y int) color.Color {
	return img.RGBA.At(x, y)
}

func (img RGBA_8888) Close() error {
	fmt.Println("Closing!")
	return syscall.Munmap(img.data)
}
