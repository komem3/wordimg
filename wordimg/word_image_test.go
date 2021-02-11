package wordimg_test

import (
	"bytes"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/komem3/word_rand_img/wordimg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font/gofont/goregular"
)

func TestGenerator_GenrateImage(t *testing.T) {
	type (
		given struct {
			message string
			randI   int
			randF   float64
			options []wordimg.Option
		}
		want struct {
			imagePath string
			err       error
		}
	)
	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			"512*512",
			given{
				message: "Hello World",
				randI:   0,
				randF:   0,
			},
			want{
				imagePath: "helloworld.png",
				err:       nil,
			},
		},
		{
			"1024*1024",
			given{
				message: "There is nothing either good or bad, but thinking makes it so.",
				randI:   2,
				randF:   0.9,
				options: []wordimg.Option{
					wordimg.WithWidth(1024),
					wordimg.WithHeight(1024),
				},
			},
			want{
				imagePath: "goodorbad.png",
				err:       nil,
			},
		},
		{
			"green and font 24px",
			given{
				message: "Hello World",
				randI:   0,
				randF:   0,
				options: []wordimg.Option{
					wordimg.WithColor(&color.RGBA{0, 255, 0, 255}),
					wordimg.WithFontSize(24),
				},
			},
			want{
				imagePath: "green24.png",
				err:       nil,
			},
		},
		{
			"just 1 line",
			given{
				message: "goodbye thank you",
				randI:   0,
				randF:   0,
				options: []wordimg.Option{
					wordimg.WithColor(&color.RGBA{255, 0, 255, 255}),
					wordimg.WithJustLine(1),
				},
			},
			want{
				imagePath: "just1line.png",
				err:       nil,
			},
		},
		{
			"just 2 line",
			given{
				message: "goodbye thank you",
				randI:   0,
				randF:   0,
				options: []wordimg.Option{
					wordimg.WithColor(&color.RGBA{255, 0, 255, 255}),
					wordimg.WithJustLine(2),
				},
			},
			want{
				imagePath: "just2line.png",
				err:       nil,
			},
		},
		{
			"text align center",
			given{
				message: "goodbye thank you",
				randI:   0,
				randF:   0,
				options: []wordimg.Option{
					wordimg.WithAlignCenter(),
				},
			},
			want{
				imagePath: "align_center.png",
				err:       nil,
			},
		},
	}
	rand := rand.New(rand.NewSource(0))
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gen := wordimg.NewGenerator(rand, goregular.TTF)
			gen.SetRand(tt.given.randI, tt.given.randF)

			b := new(bytes.Buffer)
			err := gen.GenerateImage(b, tt.given.message, tt.given.options...)
			require.NoError(t, err)

			wantFile, err := os.Open(filepath.Join("testdata", tt.want.imagePath))
			require.NoError(t, err)
			wantBytes, err := ioutil.ReadAll(wantFile)
			require.NoError(t, err)

			if ok := reflect.DeepEqual(wantBytes, b.Bytes()); !ok {
				tmpFile, err := ioutil.TempFile(".", "*.png")
				require.NoError(t, err)
				defer func() {
					err = tmpFile.Close()
					assert.NoError(t, err)
				}()
				_, err = tmpFile.Write(b.Bytes())
				assert.NoError(t, err)
				t.Errorf("GenerateImage: %s and %s are different", tt.want.imagePath, tmpFile.Name())
			}
		})
	}
}
