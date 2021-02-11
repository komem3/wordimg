package wordimg

import "image/color"

type align int

const (
	alignLeft align = iota + 1
	alignCenter
)

type (
	// Option is a setting value to image generation.
	Option func(*config)

	config struct {
		width    int
		height   int
		color    *color.RGBA
		fontSize int
		justLine int
		align    align
	}
)

func newConfig(os ...Option) config {
	c := config{
		width:  512,
		height: 512,
		align:  alignLeft,
	}
	for _, o := range os {
		o(&c)
	}
	return c
}

// WithWidth set width.
func WithWidth(w int) Option {
	return func(c *config) {
		c.width = w
	}
}

// WithHeight set height.
func WithHeight(w int) Option {
	return func(c *config) {
		c.height = w
	}
}

// WithColor set color
func WithColor(col *color.RGBA) Option {
	return func(c *config) {
		c.color = col
	}
}

// WithFontSize set font size.
func WithFontSize(f int) Option {
	return func(c *config) {
		c.fontSize = f
	}
}

// WithJustLine specify how many lines.
func WithJustLine(l int) Option {
	return func(c *config) {
		c.justLine = l
	}
}

// WithAlignCenter center the text.
func WithAlignCenter() Option {
	return func(c *config) {
		c.align = alignCenter
	}
}
