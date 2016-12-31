// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package artist

import (
	"image"
	"image/color"
	"math"
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
func RateImage(im image.Image, c color.Color, ideal float64) float64 {
	obs := perceiveColor(im, c)
	diff := math.Abs(ideal - obs)
	return 1.0 - diff
}
