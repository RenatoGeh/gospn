package common

import (
	"fmt"
	"math/rand"
	"os"
)

// Color (RGB).
type Color struct {
	r int
	g int
	b int
}

// Red color.
var Red *Color

// Green color.
var Green *Color

// Blue color.
var Blue *Color

// White color.
var White *Color

// Black color.
var Black *Color

// NewColor creates a new color given components.
func NewColor(r, g, b int) *Color {
	c := new(Color)
	c.r, c.g, c.b = r, g, b
	return c
}

func init() {
	Red = NewColor(255, 0, 0)
	Green = NewColor(0, 255, 0)
	Blue = NewColor(0, 0, 255)
	White = NewColor(255, 255, 255)
	Black = NewColor(0, 0, 0)
}

type rgb struct {
	r float64
	g float64
	b float64
}

type hsv struct {
	h float64
	s float64
	v float64
}

//func rgb2hsv(in rgb) hsv {
//var out hsv
//var min, max, delta float64

//if in.r < in.g {
//min = in.r
//} else {
//min = in.g
//}
//if min >= in.b {
//min = in.b
//}

//if in.r > in.g {
//max = in.r
//} else {
//max = in.g
//}
//if max <= in.b {
//max = in.b
//}

//out.v = max
//delta = max - min

//if delta < 0.00001 {
//out.s = 0
//out.h = 0
//return out
//}
//if max > 0 {
//out.s = (delta / max)
//} else {
//out.s = 0
//out.h = math.NaN()
//return out
//}
//if in.r >= max {
//out.h = (in.g - in.b) / delta
//} else if in.g >= max {
//out.h = 2.0 + (in.b-in.r)/delta
//} else {
//out.h = 4.0 + (in.r-in.g)/delta
//}
//out.h *= 60.0

//if out.h < 0.0 {
//out.h += 360.0
//}

//return out
//}

func hsv2rgb(in hsv) rgb {
	var hh, p, q, t, ff float64
	var i int
	var out rgb

	if in.s <= 0.0 {
		out.r = in.v
		out.g = in.v
		out.b = in.v
		return out
	}
	hh = in.h
	if hh >= 360.0 {
		hh = 0.0
	}
	hh /= 60.0
	i = int(hh)
	ff = hh - float64(i)
	p = in.v * (1.0 - in.s)
	q = in.v * (1.0 - (in.s * ff))
	t = in.v * (1.0 - (in.s * (1.0 - ff)))

	switch i {
	case 0:
		out.r = in.v
		out.g = t
		out.b = p
	case 1:
		out.r = q
		out.g = in.v
		out.b = p
	case 2:
		out.r = p
		out.g = in.v
		out.b = t
	case 3:
		out.r = p
		out.g = q
		out.b = in.v
	case 4:
		out.r = t
		out.g = p
		out.b = in.v
	case 5:
	default:
		out.r = in.v
		out.g = p
		out.b = q
	}
	return out
}

// DrawColor simply calls fmt.Fprintf and writes the RGB components of c to file.
func DrawColor(file *os.File, c *Color) {
	fmt.Fprintf(file, "%d %d %d", c.r, c.g, c.b)
}

// DrawColorRGB simply calls fmt.Fprintf and writes the RGB components (r, g, b) to file.
func DrawColorRGB(file *os.File, r, g, b int) {
	fmt.Fprintf(file, "%d %d %d", r, g, b)
}

// String returns a color string representation.
func (c *Color) String() string { return fmt.Sprintf("%d %d %d", c.r, c.g, c.b) }

// RandColor returns a random color.
func RandColor() *Color { return NewColor(rand.Intn(256), rand.Intn(256), rand.Intn(256)) }

// RandTone returns a random tone for a given hue tone (HSV).
func RandTone(h float64) *Color {
	var c hsv
	c.h = h
	c.s, c.v = rand.Float64(), rand.Float64()
	nc := hsv2rgb(c)
	return &Color{r: int(255 * nc.r), g: int(255 * nc.g), b: int(255 * nc.b)}
}

// RandColorScale returns a scaled color according to the interval [0,max], point and tone given.
func RandColorScale(p, max, h, s, minB float64) *Color {
	var c hsv
	c.h = h
	c.s = s
	c.v = (1-minB)*p/max + minB
	nc := hsv2rgb(c)
	return &Color{r: int(255 * nc.r), g: int(255 * nc.g), b: int(255 * nc.b)}
}
