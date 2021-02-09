package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	fontPath := os.Getenv("FONT_PATH")
	var fontData []byte
	if fontPath == "" {
		fontData = goregular.TTF
	} else {
		f, err := os.Open(fontPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		fontData, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
	}

	r := chi.NewRouter()
	r.Route("/wordimg", func(r chi.Router) {
		wordImgHandleGroup(r, fontData)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("[DEBUG] bind port %s\n", port)
	panic(http.ListenAndServe(":"+port, r))
}
