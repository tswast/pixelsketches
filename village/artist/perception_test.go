// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package artist

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/tswast/pixelsketches/palettes"
)

var colortests = []struct {
	cnt int
	clr color.Color
}{
	{100, palettes.PICO8_BLACK},
	{73, palettes.PICO8_PINK},
}

func TestPerceiveColor(t *testing.T) {
	r := image.Rectangle{image.Point{0, 0}, image.Point{10, 10}}
	im := image.NewPaletted(r, palettes.PICO8)
	for _, tt := range colortests {
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				cnt := x*10 + y + 1
				if cnt <= tt.cnt {
					im.Set(x, y, tt.clr)
				} else {
					im.Set(x, y, palettes.PICO8_BLACK)
				}
			}
		}

		got := perceiveColor(im, tt.clr)

		if math.Abs(got-float64(tt.cnt)/100.0) > 0.001 {
			t.Errorf("perceiveColor(im, %v) => %f, but expected %f", tt.clr, got, float64(tt.cnt)/100.0)
		}
	}
}
