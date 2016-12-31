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

func setProp(im *image.Paletted, clr color.Color, cnt int) {
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			c := x*10 + y + 1
			if c <= cnt {
				im.Set(x, y, clr)
			} else {
				im.Set(x, y, palettes.PICO8_BLACK)
			}
		}
	}
}

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
		setProp(im, tt.clr, tt.cnt)

		got := perceiveColor(im, tt.clr)

		if math.Abs(got-float64(tt.cnt)/100.0) > 0.001 {
			t.Errorf("perceiveColor(im, %v) => %f, but expected %f", tt.clr, got, float64(tt.cnt)/100.0)
		}
	}
}

var ratetests = []struct {
	ideal    int
	actual   int
	expected float64
}{
	// Ideal is an endpoint.
	{100, 0, 0.0},
	{100, 100, 1.0},
	{0, 0, 1.0},
	{0, 50, 0.5},
	{0, 100, 0.0},
	// Ideal is the center.
	{50, 0, 0.0},
	{50, 25, 0.5},
	{50, 50, 1.0},
	{50, 75, 0.5},
	{50, 100, 0.0},
	// Ideal is offset from center.
	{40, 0, 0.0},
	{40, 20, 0.5},
	{40, 40, 1.0},
	{40, 70, 0.5},
	{40, 100, 0.0},
	{60, 0, 0.0},
	{60, 30, 0.5},
	{60, 60, 1.0},
	{60, 80, 0.5},
	{60, 100, 0.0},
}

func TestRateImage(t *testing.T) {
	r := image.Rectangle{image.Point{0, 0}, image.Point{10, 10}}
	im := image.NewPaletted(r, palettes.PICO8)
	for _, tt := range ratetests {
		setProp(im, palettes.PICO8_PINK, tt.actual)

		got := RateImage(im, palettes.PICO8_PINK, float64(tt.ideal)/100)

		if math.Abs(got-tt.expected) > 0.001 {
			t.Errorf(
				"RateColor(im[%d%% pink], PICO8_PINK, %f) => %f, but expected %f",
				tt.actual,
				float64(tt.ideal)/100,
				got,
				tt.expected)
		}
	}
}
