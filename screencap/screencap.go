package screencap

import (
	"encoding/binary"
	"errors"
	"image"
	"io"
)

const (
	rgba8888 = 1
	rgb888   = 3
	rgb565   = 4
)

var NotImplemented = errors.New("not implemented: RGB_888 and RGB_565")

func Decode(r io.Reader) (image.Image, error) {
	var buf [3]int32
	err := binary.Read(r, binary.LittleEndian, &buf)
	if err != nil {
		return nil, err
	}
	width := int(buf[0])
	height := int(buf[1])
	format := int(buf[2])
	bounds := image.Rect(0, 0, width, height)
	switch format {
	case rgba8888:
		img := image.NewRGBA(bounds)
		err = binary.Read(r, binary.LittleEndian, &img.Pix)
		if err != nil {
			return nil, err
		}
		return img, err
	case rgb888:
	case rgb565:
		return nil, NotImplemented
	}
	return nil, image.ErrFormat
}
