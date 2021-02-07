package main

import (
	"math/rand"
	"os"
	"time"

	wordrandimg "github.com/komom3/playtools/word_rand_img"
)

func main() {
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	generator := wordrandimg.NewGenerator(rand)
	err = generator.GenerateImage(f, 512, 512, "Hello World")
	if err != nil {
		panic(err)
	}
}
