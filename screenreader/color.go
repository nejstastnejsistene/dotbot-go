package screenreader

import (
	"image/color"
	"math"

	"github.com/nejstastnejsistene/dotbot-go/solver"
)

// Determine if a color is a background color by checking its chroma.
func isBackground(col color.Color) bool {
	return Chroma(col) < 0.1
}

// Convert Go colors to RGB values in range [0,1).
func normalize(col color.Color) (r, g, b float64) {
	ri, gi, bi, _ := col.RGBA()
	r = float64(ri) / float64(0x10000)
	g = float64(gi) / float64(0x10000)
	b = float64(bi) / float64(0x10000)
	return
}

// Calculate the chroma of a color, from [0,1).
// https://en.wikipedia.org/wiki/HSL_and_HSV
func Chroma(col color.Color) (v float64) {
	r, g, b := normalize(col)
	M := math.Max(r, math.Max(g, b))
	m := math.Min(r, math.Min(g, b))
	return M - m
}

// Calculate the hue of a color, from [0,2pi).
// https://en.wikipedia.org/wiki/Hue
func Hue(col color.Color) (hue float64) {
	r, g, b := normalize(col)
	hue = math.Pi / 3
	switch {
	case r >= g && g >= b:
		hue *= (g - b) / (r - b)
	case g > r && r >= b:
		hue *= 2 - (r-b)/(g-b)
	case g >= b && b > r:
		hue *= 2 + (b-r)/(g-r)
	case b > g && g > r:
		hue *= 4 - (g-r)/(b-r)
	case b > r && r >= g:
		hue *= 4 + (r-g)/(b-g)
	case r >= b && b > g:
		hue *= 6 - (b-g)/(r-g)
	}
	return
}

// Convert a hue to the nearest dot color.
func HueToColor(hue float64) solver.Color {
	n := len(solver.Colors)
	index := int(float64(n)*hue/(2*math.Pi) + 0.5)
	return solver.Colors[index%n]
}
