package screencap

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"os"
	"os/exec"
	"syscall"
)

const (
	rgba8888 = 1
	rgb888   = 3
	rgb565   = 4
)

var NotImplemented = errors.New("not implemented: RGB_888 and RGB_565")

type ScreenCap struct {
	data []byte
	img  image.Image
}

func (img *ScreenCap) ColorModel() color.Model {
	return img.img.ColorModel()
}

func (img *ScreenCap) Bounds() image.Rectangle {
	return img.img.Bounds()
}

func (img *ScreenCap) At(x, y int) color.Color {
	return img.img.At(x, y)
}

func (img *ScreenCap) Close() error {
	return syscall.Munmap(img.data)
}

func NewScreenCap() (*ScreenCap, error) {
	filename := "/data/local/DotBot/tmp.screencap"
	err := TakeScreenCap(filename)
	if err != nil {
		return nil, err
	}
	return Mmap(filename)
}

func TakeScreenCap(filename string) error {
	return exec.Command("screencap", filename).Run()
}

func Mmap(filename string) (*ScreenCap, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	var buf [3]int32
	err = binary.Read(f, binary.LittleEndian, &buf)
	if err != nil {
		return nil, err
	}
	width := int(buf[0])
	height := int(buf[1])
	format := int(buf[2])
	switch format {
	case rgba8888:
		bpp := 4
		size := width * height * bpp
		offset := 12
		data, err := syscall.Mmap(int(f.Fd()), 0, size+offset,
			syscall.PROT_READ, syscall.MAP_PRIVATE)
		if err != nil {
			return nil, err
		}
		return &ScreenCap{
			data,
			&image.RGBA{
				data[offset:],
				width * bpp,
				image.Rect(0, 0, width, height),
			},
		}, err
	case rgb888:
	case rgb565:
		return nil, NotImplemented
	}
	return nil, image.ErrFormat
}
