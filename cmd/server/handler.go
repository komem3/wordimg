package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/komem3/word_rand_img/wordimg"
)

type wordImgHandler struct {
	wordimg.Generator
}

func wordImgHandleGroup(r chi.Router, fontData []byte) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	h := &wordImgHandler{
		Generator: wordimg.NewGenerator(rand, fontData),
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

	options, err := h.optionParse(query)
	if err != nil {
		log.Printf("[WARN] query parse to option: %v\n", err)
		if merr, ok := err.(*multierror.Error); ok {
			for _, err := range merr.Errors {
				log.Printf("[WARN] %v\n", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := h.GenerateImage(w, message, options...); err != nil {
		log.Printf("[ERROR] generate: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("internal server error")); err != nil {
			panic(err)
		}
		return
	}
	w.Header().Add("Content-Type", "image/png")
}

func (*wordImgHandler) optionParse(query url.Values) (options []wordimg.Option, merr error) {
	fontSize := query.Get("size")
	if fontSize != "" {
		if fontSize[0] == 'j' {
			if len(fontSize) < 2 {
				merr = multierror.Append(merr, fmt.Errorf("%s of the format is not 'j${line}'", fontSize))
			}
			justLine, e := strconv.Atoi(fontSize[1:])
			if e != nil {
				merr = multierror.Append(merr, fmt.Errorf("%s of the format is not 'j${line}': %w", fontSize, e))
			}
			options = append(options, wordimg.WithJustLine(justLine))
		} else {
			size, err := strconv.Atoi(fontSize)
			if err != nil {
				merr = multierror.Append(merr, fmt.Errorf("font convert: %w", err))
			} else {
				options = append(options, wordimg.WithFontSize(size))
			}
		}
	}
	rgb := query.Get("color")
	if rgb != "" {
		color, err := wordimg.ConvertColor(rgb)
		if err != nil {
			merr = multierror.Append(merr, fmt.Errorf("color convert: %w", err))
		} else {
			options = append(options, wordimg.WithColor(color))
		}
	}
	align := query.Get("align")
	if align != "" {
		if align == "center" {
			options = append(options, wordimg.WithAlignCenter())
		} else {
			merr = multierror.Append(merr, fmt.Errorf("align convert %s is not 'center'", align))
		}
	}

	return options, merr
}
