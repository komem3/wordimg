package wordimg_test

import (
	"bytes"
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
			width   int
			height  int
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
				width:   512,
				height:  512,
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
				width:   1024,
				height:  1024,
			},
			want{
				imagePath: "goodorbad.png",
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
			err := gen.GenerateImage(b, tt.given.message,
				wordimg.WithWidth(tt.given.width),
				wordimg.WithHeight(tt.given.height),
			)
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
