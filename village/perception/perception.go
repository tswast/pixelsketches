// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package perception

import (
	"image"
	"image/color"
	"math"

	"github.com/tswast/pixelsketches/palettes"
)

// Interests were randomly generated using:
// for i := 0; i < 16; i++ {
// 	fmt.Printf("%f\n", rand.Float32())
// }
const (
	IdealBlack      = 0.604660
	IdealDarkBlue   = 0.940509
	IdealDarkPurple = 0.664560
	IdealDarkGreen  = 0.437714
	IdealBrown      = 0.424637
	IdealDarkGray   = 0.686823
	IdealLightGray  = 0.065637
	IdealWhite      = 0.156519
	IdealRed        = 0.096970
	IdealOrange     = 0.300912
	IdealYellow     = 0.515213
	IdealGreen      = 0.813640
	IdealBlue       = 0.214264
	IdealIndigo     = 0.380657
	IdealPink       = 0.318058
	IdealPeach      = 0.468890
)

// Rating rates an image.
//
// The return value *should* be in [0, 1].
type Rating func(im image.Image) float64

func CountColors(im image.Image) map[color.Color]int {
	m := make(map[color.Color]int)
	b := im.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y
	pxls := w * h
	if pxls == 0 {
		return m
	}
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			m[im.At(x, y)] += 1
		}
	}
	return m
}

// RateImage rates an image from 0 to 1 based on perception of color.
//
// Value is 0 at the endpoints and 1 at the ideal value.
func RateImage(pxls, cnt int, ideal float64) float64 {
	x := float64(cnt) / float64(pxls)
	m := 1.0 / ideal
	b := 0.0
	// If ideal is exactly 0, make sure the line slopes down
	if x > ideal || ideal == 0.0 {
		m = -1.0 / (1.0 - ideal)
		b = -(-1.0 / (1.0 - ideal))
	}
	return m*x + b
}

// NewRating creates a rating function which desires an ideal amount of a color.
func NewRating(ideal float64, c color.Color) Rating {
	r := func(im image.Image) float64 {
		b := im.Bounds()
		w := b.Max.X - b.Min.X
		h := b.Max.Y - b.Min.Y
		pxls := w * h
		cnt := CountColors(im)[c]
		return RateImage(pxls, cnt, ideal)
	}
	return r
}

// RateBlack rates an image according to black's interest.
func RateBlack(im image.Image) float64 {
	b := im.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y
	pxls := w * h
	cnt := CountColors(im)[palettes.PICO8_BLACK]
	return RateImage(pxls, cnt, IdealBlack)
}

// RateWholeImage rates an image according to predefined "interests".
func RateWholeImage(im image.Image) float64 {
	b := im.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y
	pxls := w * h
	cnts := CountColors(im)
	rt := 0.0
	r := RateImage(pxls, cnts[palettes.PICO8_BLACK], IdealBlack)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_DARK_BLUE], IdealDarkBlue)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_DARK_PURPLE], IdealDarkPurple)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_DARK_GREEN], IdealDarkGreen)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_BROWN], IdealBrown)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_DARK_GRAY], IdealDarkGray)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_LIGHT_GRAY], IdealLightGray)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_WHITE], IdealWhite)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_RED], IdealRed)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_ORANGE], IdealOrange)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_YELLOW], IdealYellow)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_GREEN], IdealGreen)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_BLUE], IdealBlue)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_INDIGO], IdealIndigo)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_PINK], IdealPink)
	rt += r
	r = RateImage(pxls, cnts[palettes.PICO8_PEACH], IdealPeach)
	rt += r
	rt /= 16.0

	return rt
}

// colorDist calculates the distance between two colors.
func colorDist(a, b color.Color) float64 {
	r1d, g1d, b1d, _ := a.RGBA()
	r2d, g2d, b2d, _ := b.RGBA()
	// Scale colors to be between 0 and 1
	r1 := float64(r1d) / 65535
	r2 := float64(r2d) / 65535
	g1 := float64(g1d) / 65535
	g2 := float64(g2d) / 65535
	b1 := float64(b1d) / 65535
	b2 := float64(b2d) / 65535
	d := 0.0
	d += math.Pow(r1-r2, 2)
	d += math.Pow(g1-g2, 2)
	d += math.Pow(b1-b2, 2)
	// Scale by the maximum possible distance.
	d /= 3
	d = math.Sqrt(d)
	return d
}

// perceiveTLCorner checks if the point at x, y is a top-left corner.
func perceiveTLCorner(x, y int, im image.Image) float64 {
	b := im.Bounds()
	// Out of boutnds?
	if x < b.Min.X || x >= b.Max.X || y < b.Min.Y || y >= b.Max.Y {
		return 0.0
	}
	v := 1.0
	c := im.At(x, y)
	if x > b.Min.X {
		// Colors should be as different as possible to the left.
		v *= colorDist(im.At(x-1, y), c)
	}
	if y > b.Min.Y {
		// Colors should be as different as possible to the top.
		v *= colorDist(im.At(x, y-1), c)
	}
	if x < b.Max.X-1 {
		// Colors should be as similar as possible to the right.
		v *= (1.0 - colorDist(im.At(x+1, y), c))
	}
	if y < b.Max.Y-1 {
		// Colors should be as similar as possible to the bottom.
		v *= (1.0 - colorDist(im.At(x, y+1), c))
	}
	return v
}

func countTLCorners(im image.Image) float64 {
	b := im.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y
	pxls := w * h
	// If x, y is a corner, then just to the right and just below cannot be corners.
	maxCorners := float64(pxls) / 3.0
	corners := 0.0
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			corners += perceiveTLCorner(x, y, im)
		}
	}
	return corners / maxCorners
}
