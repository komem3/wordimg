package wordimg

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/rand"
	"strings"
	"unicode/utf8"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	paddingLeft  = 0.05
	paddingRight = 0.9
	padding      = paddingRight - paddingLeft
)

// Generator contain font data.
type Generator struct {
	randFunc func() float64
	colorGen *colorGenerator
	font     []byte
}

type wordDrawer struct {
	*font.Drawer
	widthMax int
	startX   int
	fontSize int
}

// NewGenerator create Generator.
func NewGenerator(rand *rand.Rand, font []byte) *Generator {
	return &Generator{
		randFunc: rand.Float64,
		colorGen: &colorGenerator{
			randFunc: rand.Intn,
		},
		font: font,
	}
}

// GenerateImage generate image from given message and write the image file to io.Writer.
func (g *Generator) GenerateImage(w io.Writer, message string, op ...Option) error {
	conf := newConfig(op...)
	drawer, err := g.prepareDrawer(message, conf)
	if err != nil {
		return err
	}
	drawer.drawMessage(message, conf)
	if err := drawer.writeImage(w); err != nil {
		return err
	}
	return nil
}

// rand retun 0.5 ~ 1.0.
func (g *Generator) rand() float64 {
	return 0.5 + 0.5*g.randFunc()
}

func (g *Generator) calcFontSize(message string, c config) float64 {
	wordSize := math.Sqrt(float64(c.width * c.height / utf8.RuneCountInString(message)))
	ww := c.width / int(wordSize) * int(wordSize)
	if c.width == ww {
		return wordSize * g.rand()
	}
	return math.Sqrt(float64(ww*c.height/utf8.RuneCountInString(message))) * g.rand()
}

func (*Generator) justFontSize(message string, fontSet *truetype.Font, c config) float64 {
	wordSize := float64(c.width * c.justLine / utf8.RuneCountInString(message))
	widthFix := fixed.Int26_6(int(float64(c.width*c.justLine)*padding) << 6)

	for face := truetype.NewFace(fontSet, &truetype.Options{Size: wordSize}); font.MeasureString(face, message) < widthFix; wordSize++ {
		face = truetype.NewFace(fontSet, &truetype.Options{Size: wordSize})
	}
	return wordSize - 1
}

func (g *Generator) prepareDrawer(message string, config config) (drawer *wordDrawer, err error) {
	fontSet, err := truetype.Parse(g.font)
	if err != nil {
		return nil, err
	}

	var fontSize float64
	switch {
	case config.justLine > 0:
		fontSize = g.justFontSize(message, fontSet, config)
	case config.fontSize > 0:
		fontSize = float64(config.fontSize)
	default:
		fontSize = g.calcFontSize(message, config)
	}

	var col color.RGBA
	if config.color != nil {
		col = *config.color
	} else {
		col = g.colorGen.randColor()
	}

	drawer = &wordDrawer{
		Drawer: &font.Drawer{
			Dst: image.NewRGBA(image.Rect(0, 0, config.width, config.height)),
			Src: image.NewUniform(col),
			Face: truetype.NewFace(fontSet, &truetype.Options{
				Size: fontSize,
			}),
		},
		fontSize: int(fontSize),
	}

	advance := drawer.MeasureString(message)
	widthMax := int(float64(config.width) * paddingRight)
	rowNum := int(advance.Ceil()/widthMax) + 1

	startY := (config.height-(rowNum*int(fontSize)))/2 + int(fontSize*0.75)
	startX := int(paddingLeft * float64(config.width))

	drawer.Dot = freetype.Pt(startX, startY)
	drawer.startX = startX
	drawer.widthMax = widthMax
	return
}

func (d *wordDrawer) drawMessage(message string, c config) {
	var sb strings.Builder
	for _, char := range message {
		advance := d.MeasureString(sb.String() + string(char))
		if advance.Ceil() > d.widthMax {
			d.DrawString(sb.String())
			sb.Reset()
			d.Dot = d.Dot.Add(freetype.Pt(d.startX-d.Dot.X.Ceil(), int(d.fontSize)))
		}
		sb.WriteRune(char)
	}
	if c.align == alignLeft {
		d.DrawString(sb.String())
		return
	}

	advance := d.MeasureString(sb.String())
	paddingLeft := (d.widthMax - advance.Ceil()) / 2
	d.Dot = d.Dot.Add(freetype.Pt(paddingLeft, 0))
	d.DrawString(sb.String())
}

func (d *wordDrawer) writeImage(w io.Writer) error {
	return png.Encode(w, d.Dst)
}
