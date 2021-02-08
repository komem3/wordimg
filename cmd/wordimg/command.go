package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/komom3/word_rand_img/wordimg"
	"golang.org/x/image/font/gofont/goregular"
)

type options struct {
	Message   string `short:"m" long:"message" description:"Message to write to image. required." required:"true"`
	ImagePath string `short:"i" long:"imagePath" description:"Path of the image to write. Default is 'unix_time.png"`
	FontPath  string `short:"f" long:"font" description:"Path to font file. Only support ttf."`
	FontSize  int    `long:"size" description:"Font size."`
	Width     int    `short:"w" long:"width" description:"Width of the generated image." default:"512"`
	Height    int    `short:"h" long:"height" description:"Height of the generated image." default:"512"`
	Color     string `short:"c" long:"color" description:"Text color."`
}

type commandLine struct {
	options
}

func (c *commandLine) parse(args []string) error {
	_, err := flags.ParseArgs(&c.options, args)
	if err != nil {
		return err
	}

	if c.ImagePath == "" {
		c.ImagePath = fmt.Sprintf("%d.png", time.Now().Unix())
	}
	return err
}

func (c *commandLine) exec() error {
	f, err := os.Create(c.ImagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var fontData []byte
	if c.FontPath == "" {
		fontData = goregular.TTF
	} else {
		tfFile, err := os.Open(c.FontPath)
		if err != nil {
			return err
		}
		defer f.Close()
		fontData, err = ioutil.ReadAll(tfFile)
		if err != nil {
			return err
		}
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	generator := wordimg.NewGenerator(rand, fontData)

	return generator.GenerateImage(f, c.Message,
		wordimg.WithWidth(c.Width),
		wordimg.WithHeight(c.Height),
		wordimg.WithFontSize(c.FontSize),
	)
}
