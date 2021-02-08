// Package wordimg implements generator to create an image from text.
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

type generator struct {
	randFunc func() float64
	colorGen *colorGenerator
	font     []byte
}

type wordDrawer struct {
	*font.Drawer
	widthFix fixed.Int26_6
	startX   int
	fontSize int
}

// NewGenerator create Generator.
func NewGenerator(rand *rand.Rand, font []byte) *generator {
	return &generator{
		randFunc: rand.Float64,
		colorGen: &colorGenerator{
			randFunc: rand.Intn,
		},
		font: font,
	}
}

// GenerateImage generate image from given message.
func (g *generator) GenerateImage(w io.Writer, message string, op ...option) error {
	conf := newConfig(op...)
	drawer, err := g.prepareDrawer(message, conf)
	if err != nil {
		return err
	}
	drawer.drawMessage(message)
	if err := drawer.writeImage(w); err != nil {
		return err
	}
	return nil
}

// rand retun 0.7 ~ 1.0.
func (g *generator) rand() float64 {
	return 0.7 + 0.3*g.randFunc()
}

func (g *generator) calcFontSize(message string, width, height int) float64 {
	wordSize := math.Sqrt(float64(width * height / utf8.RuneCountInString(message)))
	ww := width / int(wordSize) * int(wordSize)
	if width == ww {
		return wordSize * g.rand()
	}
	return math.Sqrt(float64(ww*height/utf8.RuneCountInString(message))) * g.rand()
}

func (g *generator) prepareDrawer(message string, config config) (drawer *wordDrawer, err error) {
	// wordSize := math.Sqrt(float64(config.width*config.height/utf8.RuneCountInString(message))) * g.rand()
	var fontSize float64
	if config.fontSize > 0 {
		fontSize = float64(config.fontSize)
	} else {
		fontSize = g.calcFontSize(message, config.width, config.height)
	}

	var col color.RGBA
	if config.color != nil {
		col = *config.color
	} else {
		col = g.colorGen.randColor()
	}

	fontSet, err := truetype.Parse(g.font)
	if err != nil {
		return nil, err
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
	widthFix := fixed.Int26_6(int(float64(config.width)*0.95) << 6)
	rowNum := int(advance/widthFix) + 1

	startY := (config.height-(rowNum*int(fontSize)))/2 + int(fontSize*0.75)
	startX := int(0.05 * float64(config.width))

	drawer.Dot = freetype.Pt(startX, startY)
	drawer.startX = startX
	drawer.widthFix = widthFix
	return
}

func (d *wordDrawer) drawMessage(message string) {
	var sb strings.Builder
	for _, char := range message {
		advance := d.MeasureString(sb.String() + string(char))
		if advance > d.widthFix {
			d.DrawString(sb.String())

			sb.Reset()
			d.Dot = d.Dot.Add(freetype.Pt(d.startX-d.Dot.X.Ceil(), int(d.fontSize)))
		}
		sb.WriteRune(char)
	}
	d.DrawString(sb.String())
}

func (d *wordDrawer) writeImage(w io.Writer) error {
	return png.Encode(w, d.Dst)
}
