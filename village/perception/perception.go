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

// ColorDist calculates the distance between two colors.
func ColorDist(a, b color.Color) float64 {
	r1, g1, b1, _ := a.RGBA()
	r2, g2, b2, _ := b.RGBA()
	d := 0.0
	d += math.Pow(float64(r1-r2), 2)
	d += math.Pow(float64(g1-g2), 2)
	d += math.Pow(float64(b1-b2), 2)
	d = math.Sqrt(d)
	d /= math.Pow(2.0, 16)
	return d
}

// PerceiveTLCorner checks if the point at x, y is a top-left corner.
func PerceiveTLCorner(x, y int, im image.Image) float64 {
	b := im.Bounds()
	// Out of boutnds?
	if x < b.Min.X || x >= b.Max.X || y < b.Min.Y || y >= b.Max.Y {
		return 0.0
	}
	v := 1.0
	c := im.At(x, y)
	if x >= b.Min.X {

	}
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y
	im.At(x, y)
}
