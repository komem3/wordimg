package wordrandimg

import (
	"image/color"
)

const randRange = 100

var (
	red    = color.RGBA{255, 0, 0, 255}
	blue   = color.RGBA{0, 0, 255, 255}
	green  = color.RGBA{0, 255, 0, 255}
	yellow = color.RGBA{255, 255, 0, 255}
	pink   = color.RGBA{255, 0, 255, 255}
	cyan   = color.RGBA{0, 255, 255, 255}

	colors = []color.RGBA{
		red, blue, green, yellow, pink, cyan,
	}
)

type colorGenerator struct {
	randFunc func(max int) int
}

func (c *colorGenerator) randColor() color.RGBA {
	base := colors[c.randFunc(len(colors))]

	ran := c.randFunc(randRange)
	if base.R == 255 {
		ran *= -1
	}
	base.R += uint8(ran)

	ran = c.randFunc(randRange)
	if base.G == 255 {
		ran *= -1
	}
	base.G += uint8(randRange)

	ran = c.randFunc(50)
	if base.B == 255 {
		ran *= -1
	}
	base.B += uint8(randRange)
	return base
}
