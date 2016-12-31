// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package artist

import (
	"image"
	"image/color"
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
