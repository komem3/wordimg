package wordimg

import "image/color"

type (
	Option func(*config)

	config struct {
		width    int
		height   int
		color    *color.RGBA
		fontSize int
	}
)

func newConfig(os ...Option) config {
	c := config{
		width:  512,
		height: 512,
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
