package wordrandimg_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	wordrandimg "github.com/komom3/playtools/word_rand_img"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_GenrateImage(t *testing.T) {
	type (
		given struct {
			randValue float64
			message   string
			width     int
			height    int
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
				randValue: 0,
				message:   "Hello World",
				width:     512,
				height:    512,
			},
			want{
				imagePath: "helloworld.png",
				err:       nil,
			},
		},
		{
			"1024*1024",
			given{
				randValue: 0,
				message:   "There is nothing either good or bad, but thinking makes it so.",
				width:     1024,
				height:    1024,
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
			gen := wordrandimg.NewGenerator(rand)
			gen.SetRand(int(tt.given.randValue), tt.given.randValue)

			b := new(bytes.Buffer)
			err := gen.GenerateImage(b, tt.given.width, tt.given.height, tt.given.message)
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
