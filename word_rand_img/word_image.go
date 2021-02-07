// Package wordrandimg implements generator to create an image from text.
package wordrandimg

import (
	"image"
	"image/png"
	"io"
	"math"
	"math/rand"
	"unicode/utf8"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type generator struct {
	randFunc func() float64
	colorGen *colorGenerator
}

type wordDrawer struct {
	*font.Drawer
	widthFix fixed.Int26_6
	startX   int
	wordSize int
}

// NewGenerator create Generator.
func NewGenerator(rand *rand.Rand) *generator {
	return &generator{
		randFunc: rand.Float64,
		colorGen: &colorGenerator{
			randFunc: rand.Intn,
		},
	}
}

// GenerateImage generate image from given message.
func (g *generator) GenerateImage(w io.Writer, width, height int, message string) error {
	drawer, err := g.prepareDrawer(width, height, message)
	if err != nil {
		return err
	}
	drawer.drawMessage(message)
	if err := drawer.writeImage(w); err != nil {
		return err
	}
	return nil
}

// rand retun 0.8 ~ 1.2.
func (g *generator) rand() float64 {
	return 0.8 + 0.4*g.randFunc()
}

func (g *generator) prepareDrawer(width, height int, message string) (drawer *wordDrawer, err error) {
	wordSize := math.Sqrt(float64(width*height/utf8.RuneCountInString(message))) * g.rand()
	col := g.colorGen.randColor()

	fontSet, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	drawer = &wordDrawer{
		Drawer: &font.Drawer{
			Dst: image.NewRGBA(image.Rect(0, 0, width, height)),
			Src: image.NewUniform(col),
			Face: truetype.NewFace(fontSet, &truetype.Options{
				Size: wordSize,
			}),
		},
		wordSize: int(wordSize),
	}

	advance := drawer.MeasureString(message)
	widthFix := fixed.Int26_6(int(float64(width)*0.95) << 6)
	rowNum := int(advance/widthFix) + 1
	startY := height/(rowNum+1) + int(wordSize/2)
	startX := int(0.05 * float64(width))

	drawer.Dot = freetype.Pt(startX, startY)
	drawer.startX = startX
	drawer.widthFix = widthFix
	return
}

func (d *wordDrawer) drawMessage(message string) {
	var runes []byte
	for _, char := range message {
		advance := d.MeasureBytes(append(runes, byte(char)))
		if advance > d.widthFix {
			d.DrawBytes(runes)

			runes = append(runes[:0], byte(char))
			d.Dot = d.Dot.Add(freetype.Pt(d.startX-d.Dot.X.Ceil(), int(d.wordSize)))
			continue
		}
		runes = append(runes, byte(char))
	}
	d.DrawBytes(runes)
}

func (d *wordDrawer) writeImage(w io.Writer) error {
	return png.Encode(w, d.Dst)
}
