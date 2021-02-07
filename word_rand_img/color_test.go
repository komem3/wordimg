package wordrandimg

import (
	"image/color"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestColorGenerator_randColor(t *testing.T) {
	type (
		given struct {
			randValue int
		}
		want struct {
			color color.RGBA
		}
	)
	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			"rand return 0 == red",
			given{
				randValue: 0,
			},
			want{
				color: color.RGBA{255, 0, 0, 0},
			},
		},
		{
			"rand return 5",
			given{
				randValue: 5,
			},
			want{
				color: color.RGBA{5, 250, 250, 0},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gen := &colorGenerator{
				randFunc: func(_ int) int {
					return tt.given.randValue
				},
			}
			c := gen.randColor()
			if diff := cmp.Diff(tt.want.color, c); diff != "" {
				t.Errorf("randColor return: want(-), got(+)\n%s\n", diff)
			}
		})
	}
}
