package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/komem3/word_rand_img/wordimg"
)

type generator interface {
	GenerateImage(w io.Writer, message string, op ...wordimg.Option) error
}

type wordImgHandler struct {
	generator
}

func wordImgHandleGroup(r chi.Router, fontData []byte) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	h := &wordImgHandler{
		generator: wordimg.NewGenerator(rand, fontData),
	}
	r.Get("/", h.radnImage)
}

func (h *wordImgHandler) radnImage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	message := query.Get("text")
	if message == "" {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("'text' in url query is empty")); err != nil {
			panic(err)
		}
		return
	}
	if len(message) > 512 {
		log.Printf("[WARN] input too long text. (%d chars)\n", len(message))
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("text is too long.")); err != nil {
			panic(err)
		}
		return
	}

	var options []wordimg.Option
	fontSize := query.Get("size")
	if fontSize != "" {
		size, err := strconv.Atoi(fontSize)
		if err != nil {
			log.Printf("[WARN] font convert: %+v\n", err)
		} else {
			options = append(options, wordimg.WithFontSize(size))
		}
	}
	rgb := query.Get("color")
	if rgb != "" {
		color, err := wordimg.ConvertColor(rgb)
		if err != nil {
			log.Printf("[WARN] color convert: %+v\n", err)
		} else {
			options = append(options, wordimg.WithColor(color))
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := h.generator.GenerateImage(w, message, options...); err != nil {
		log.Printf("[ERROR] generate: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("internal server error")); err != nil {
			panic(err)
		}
		return
	}
	w.Header().Add("Content-Type", "image/png")
}
