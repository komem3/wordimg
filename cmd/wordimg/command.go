package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/komem3/wordimg"
	"golang.org/x/image/font/gofont/goregular"
)

type Align string

const (
	AlignLight  = "light"
	AlignCenter = "center"
)

type options struct {
	Message   string `short:"m" long:"message" description:"Message to write to image. required." required:"true"`
	ImagePath string `short:"i" long:"imagePath" description:"Path of the image to write. Default is 'unix_time.png"`
	FontPath  string `short:"f" long:"font" description:"Path to font file. Only support ttf."`
	FontSize  string `long:"size" description:"Font size. Supports j${line} as a special format. If you specify 'j1', the font size will fit on one line."`
	Align     Align  `short:"a" long:"align" description:"Word postion. 'left' or 'center'"`
	Width     int    `short:"w" long:"width" description:"Width of the generated image." default:"512"`
	Height    int    `short:"h" long:"height" description:"Height of the generated image." default:"512"`
	Color     string `short:"c" long:"color" description:"Text color. Format is 'R:G:B'. Example: 255:255:0(yellow)"`
}

type commandLine struct {
	options

	size     int
	justLine int
}

func (c *commandLine) parse(args []string) error {
	_, err := flags.ParseArgs(&c.options, args)
	if err != nil {
		return err
	}

	if c.ImagePath == "" {
		c.ImagePath = fmt.Sprintf("%d.png", time.Now().Unix())
	}

	if c.FontSize == "" {
		return nil
	}
	if c.FontSize[0] == 'j' {
		if len(c.FontSize) < 2 {
			return fmt.Errorf("%s of the format is not 'j${line}'", c.FontSize)
		}
		justLine, e := strconv.Atoi(c.FontSize[1:])
		if e != nil {
			return fmt.Errorf("%s of the format is not 'j${line}': %w", c.FontSize, e)
		}
		c.justLine = justLine
		return nil
	}

	fontSize, err := strconv.Atoi(c.FontSize)
	if err != nil {
		return fmt.Errorf("%s is not int: %w", c.FontSize, err)
	}
	c.size = fontSize

	return nil
}

func (c *commandLine) exec() error {
	f, err := os.Create(c.ImagePath)
	if err != nil {
		return fmt.Errorf("create %s: %w", c.ImagePath, err)
	}
	defer f.Close()

	var fontData []byte
	if c.FontPath == "" {
		fontData = goregular.TTF
	} else {
		tfFile, e := os.Open(c.FontPath)
		if e != nil {
			return fmt.Errorf("open font file %s: %w", c.FontPath, e)
		}
		defer f.Close()
		fontData, err = ioutil.ReadAll(tfFile)
		if err != nil {
			return fmt.Errorf("read font data: %w", err)
		}
	}

	color, err := wordimg.ConvertColor(c.Color)
	if err != nil {
		return fmt.Errorf("convert color.RGBA from %s: %w", c.Color, err)
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	generator := wordimg.NewGenerator(rand, fontData)

	opts := []wordimg.Option{
		wordimg.WithWidth(c.Width),
		wordimg.WithHeight(c.Height),
		wordimg.WithFontSize(c.size),
		wordimg.WithJustLine(c.justLine),
		wordimg.WithColor(color),
	}
	if c.Align == AlignCenter {
		opts = append(opts, wordimg.WithAlignCenter())
	}

	return generator.GenerateImage(f, c.Message, opts...)
}
