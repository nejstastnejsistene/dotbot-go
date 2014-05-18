package screencap

import (
	"encoding/binary"
	"errors"
	"image"
	"os"
	"syscall"
)

const (
	rgba8888 = 1
	rgb888   = 3
	rgb565   = 4
)

var NotImplemented = errors.New("not implemented: RGB_888 and RGB_565")

func Mmap(filename string) (ScreenCap, error) {
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
		data, err := syscall.Mmap(int(f.Fd()), 0, 4*width*height,
			syscall.PROT_READ, syscall.MAP_PRIVATE)
		if err != nil {
			return nil, err
		}
		return RGBA_8888{
			data,
			image.RGBA{
				data[12:],
				4 * width,
				image.Rect(0, 0, width, height),
			},
		}, err
	case rgb888:
	case rgb565:
		return nil, NotImplemented
	}
	return nil, image.ErrFormat
}
