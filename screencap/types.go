package screencap

import (
	"image"
	"image/color"
	"io"
	"syscall"
)

type ScreenCap interface {
	image.Image
	io.Closer
}

type screencap struct {
	data []byte
	img  image.Image
}

func (img screencap) ColorModel() color.Model {
	return img.img.ColorModel()
}

func (img screencap) Bounds() image.Rectangle {
	return img.img.Bounds()
}

func (img screencap) At(x, y int) color.Color {
	return img.img.At(x, y)
}

func (img screencap) Close() error {
	return syscall.Munmap(img.data)
}
