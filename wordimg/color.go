package wordimg

import (
	"errors"
	"image/color"
	"strconv"
	"strings"
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
	base.G += uint8(ran)

	ran = c.randFunc(50)
	if base.B == 255 {
		ran *= -1
	}
	base.B += uint8(ran)
	return base
}

// ConvertColor convert color.RGBA from "R:G:B" string.
func ConvertColor(str string) (*color.RGBA, error) {
	if str == "" {
		return nil, nil
	}
	rgb := strings.Split(str, ":")
	if len(rgb) != 3 {
		return nil, errors.New("not color format('R:G:B)")
	}

	r, err := strconv.ParseUint(rgb[0], 10, 8)
	if err != nil {
		return nil, errors.New("red property is not integer")
	}
	g, err := strconv.ParseUint(rgb[1], 10, 8)
	if err != nil {
		return nil, errors.New("green property is not integer")
	}
	b, err := strconv.ParseUint(rgb[2], 10, 8)
	if err != nil {
		return nil, errors.New("blue property is not integer")
	}
	return &color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}
