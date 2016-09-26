package common

import (
	"fmt"
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

func newColor(r, g, b int) *Color {
	c := new(Color)
	c.r, c.g, c.b = r, g, b
	return c
}

func init() {
	Red = newColor(255, 0, 0)
	Green = newColor(0, 255, 0)
	Blue = newColor(0, 0, 255)
	White = newColor(255, 255, 255)
	Black = newColor(0, 0, 0)
}

// DrawColor simply calls fmt.Fprintf and writes the RGB components of c to file.
func DrawColor(file *os.File, c *Color) {
	fmt.Fprintf(file, "%d %d %d", c.r, c.g, c.b)
}
