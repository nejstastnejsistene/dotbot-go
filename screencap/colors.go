package screencap

import (
	"image"
	"image/color"
)

type RGB struct {
	BytesPerPixel int
	Pix           []uint8
	Stride        int
	Rect          image.Rectangle
}

func (p *RGB) ColorModel() color.Model {
	switch p.BytesPerPixel {
	case 24:
		return RGB888Model
	case 16:
		return RGB565Model
	default:
		panic("Unknown BytesPerPixel")
	}
	return nil
}

func (p *RGB) Bounds() image.Rectangle {
	return p.Rect
}

func (p *RGB) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	switch p.BytesPerPixel {
	case 24:
		return &RGB888{p.Pix[i+0], p.Pix[i+1], p.Pix[i+2]}
	case 16:
		c := new(RGB565)
		c[0] = p.Pix[i+0]
		c[1] = p.Pix[i+1]
		return c
	default:
		panic("Unknown BytesPerPixel")
	}
	return color.RGBA{}
}

func (p *RGB) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*p.BytesPerPixel
}

type RGB888 struct {
	R, G, B uint8
}

func (c *RGB888) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = 0xffff
	return
}

type RGB565 [2]uint8

func (c *RGB565) RGBA() (r, g, b, a uint32) {
	r = uint32(c[0]>>3) & 31
	r <<= 3
	r |= r << 8
	g = uint32(c[0]&7) << 3
	g |= uint32(c[1]>>5) & 7
	g <<= 2
	g |= g << 8
	b = uint32(c[1] & 31)
	b <<= 3
	b |= b << 8
	a = 0xffff
	return
}

var RGB888Model = color.ModelFunc(func(c color.Color) color.Color {
	return nil
})

var RGB565Model = color.ModelFunc(func(c color.Color) color.Color {
	return nil
})
