// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// json.Unmarshal()
	black := color.RGBA{0, 0, 0, 255}
	darkBlue := color.RGBA{29, 43, 83, 255}
	darkPurple := color.RGBA{126, 37, 83, 255}
	darkGreen := color.RGBA{0, 135, 81, 255}
	brown := color.RGBA{171, 82, 54, 255}
	darkGray := color.RGBA{95, 87, 79, 255}

	pal := []color.Color{
		black,
		darkBlue,
		darkPurple,
		darkGreen,
		brown,
		darkGray,
	}
	r := image.Rect(0, 0, 64, 64)
	im := image.NewPaletted(r, pal)
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			im.Set(x, y, black)
		}
	}
	im.Set(42, 42, darkBlue)
	im.Set(41, 41, darkPurple)
	im.Set(40, 40, darkGreen)
	im.Set(39, 39, brown)
	im.Set(38, 38, darkGray)

	f, err := os.Create("out.png")
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := png.Encode(w, im); err != nil {
		panic(err)
	}
	w.Flush()
}
