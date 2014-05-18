package screencap

import (
	"encoding/binary"
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

type ScreenCap struct {
	mapBase []byte
	img     image.Image
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
	return syscall.Munmap(img.mapBase)
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
	return exec.Command("/system/bin/screencap", filename).Run()
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

	var bytesPerPixel int
	switch format {
	case rgba8888:
		bytesPerPixel = 4
	case rgb888:
		bytesPerPixel = 3
	case rgb565:
		bytesPerPixel = 2
	default:
		return nil, image.ErrFormat
	}

	size := width * height * bytesPerPixel
	offset := 12
	mapBase, err := syscall.Mmap(int(f.Fd()), 0, size+offset,
		syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}

	var img image.Image
	pix := mapBase[offset:]
	stride := width * bytesPerPixel
	rect := image.Rect(0, 0, width, height)
	switch format {
	case rgba8888:
		img = &image.RGBA{pix, stride, rect}
	case rgb888:
	case rgb565:
		img = &RGB{bytesPerPixel, pix, stride, rect}
	}
	return &ScreenCap{mapBase, img}, err
}
