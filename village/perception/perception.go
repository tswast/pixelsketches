// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package perception

import (
	"image"
	"image/color"

	"github.com/tswast/pixelsketches/palettes"
)

func perceiveColor(im image.Image, c color.Color) float64 {
	b := im.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y
	pxls := w * h
	if pxls == 0 {
		return 0.0
	}
	cnt := 0.0
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if im.At(x, y) == c {
				cnt += 1.0
			}
		}
	}
	return cnt / float64(pxls)
}

// RateImage rates an image from 0 to 1 based on perception of color.
//
// Value is 0 at the endpoints and 1 at the ideal value.
func RateImage(im image.Image, c color.Color, ideal float64) float64 {
	x := perceiveColor(im, c)
	m := 1.0 / ideal
	b := 0.0
	// If ideal is exactly 0, make sure the line slopes down
	if x > ideal || ideal == 0.0 {
		m = -1.0 / (1.0 - ideal)
		b = -(-1.0 / (1.0 - ideal))
	}
	return m*x + b
}

// RateWholeImage rates an image according to predefined "interests".
func RateWholeImage(im image.Image) float64 {
	rt := 0.0
	// Interests were randomly generated using:
	// for i := 0; i < 16; i++ {
	// 	fmt.Printf("%f\n", rand.Float32())
	// }
	r := RateImage(im, palettes.PICO8_BLACK, 0.604660)
	rt += r
	r = RateImage(im, palettes.PICO8_DARK_BLUE, 0.940509)
	rt += r
	r = RateImage(im, palettes.PICO8_DARK_PURPLE, 0.664560)
	rt += r
	r = RateImage(im, palettes.PICO8_DARK_GREEN, 0.437714)
	rt += r
	r = RateImage(im, palettes.PICO8_BROWN, 0.424637)
	rt += r
	r = RateImage(im, palettes.PICO8_DARK_GRAY, 0.686823)
	rt += r
	r = RateImage(im, palettes.PICO8_LIGHT_GRAY, 0.065637)
	rt += r
	r = RateImage(im, palettes.PICO8_WHITE, 0.156519)
	rt += r
	r = RateImage(im, palettes.PICO8_RED, 0.096970)
	rt += r
	r = RateImage(im, palettes.PICO8_ORANGE, 0.300912)
	rt += r
	r = RateImage(im, palettes.PICO8_YELLOW, 0.515213)
	rt += r
	r = RateImage(im, palettes.PICO8_GREEN, 0.813640)
	rt += r
	r = RateImage(im, palettes.PICO8_BLUE, 0.214264)
	rt += r
	r = RateImage(im, palettes.PICO8_INDIGO, 0.380657)
	rt += r
	r = RateImage(im, palettes.PICO8_PINK, 0.318058)
	rt += r
	r = RateImage(im, palettes.PICO8_PEACH, 0.468890)
	rt += r
	rt /= 16.0
	return rt
}
