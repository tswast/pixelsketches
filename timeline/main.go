// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

const (
	screenWidth  int = 114
	screenHeight int = 64
	imageHeight  int = 64
	imageWidth   int = 64
	imageX       int = 25
	buttonHeight int = 4
	buttonBuffer int = 2
	exitX        int = imageX + imageHeight + buttonBuffer
	exitY        int = screenHeight - buttonHeight
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	rand.Seed(19700101)

	// json.Unmarshal()
	black := color.RGBA{0, 0, 0, 255}
	darkBlue := color.RGBA{29, 43, 83, 255}
	darkPurple := color.RGBA{126, 37, 83, 255}
	darkGreen := color.RGBA{0, 135, 81, 255}
	brown := color.RGBA{171, 82, 54, 255}
	darkGray := color.RGBA{95, 87, 79, 255}
	lightGray := color.RGBA{194, 195, 199, 255}
	white := color.RGBA{255, 241, 232, 255}
	red := color.RGBA{255, 0, 77, 255}
	orange := color.RGBA{255, 163, 0, 255}
	yellow := color.RGBA{255, 236, 39, 255}
	green := color.RGBA{0, 228, 54, 255}
	blue := color.RGBA{41, 173, 255, 255}
	indigo := color.RGBA{131, 118, 156, 255}
	pink := color.RGBA{255, 119, 168, 255}
	peach := color.RGBA{255, 204, 170, 255}

	pal := []color.Color{
		black,
		darkBlue,
		darkPurple,
		darkGreen,
		brown,
		darkGray,
		lightGray,
		white,
		red,
		orange,
		yellow,
		green,
		blue,
		indigo,
		pink,
		peach,
	}
	r := image.Rect(0, 0, imageWidth, imageHeight)
	im := image.NewPaletted(r, pal)
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			im.Set(x, y, black)
		}
	}

	cx := imageX + rand.Intn(imageWidth)
	cy := rand.Intn(imageHeight)
	color := pal[rand.Intn(len(pal))]
	pressX := 0
	pressY := 0
	prevPressed := false
	pressed := false
	vert := 0
	horiz := 0
	for {
		// Moving?
		cy = vert + cy
		if cy < 0 {
			cy = 0
		}
		if cy >= screenHeight {
			cy = screenHeight - 1
		}
		cx = horiz + cx
		if cx < 0 {
			cx = 0
		}
		if cx >= screenWidth {
			cx = screenWidth - 1
		}

		// Just pressed?
		if pressed && !prevPressed {
			pressX = cx
			pressY = cy
		}

		// Drawing?
		if pressed && pressX >= imageX && pressX < imageX+imageWidth &&
			cx >= imageX && cx < imageX+imageWidth {
			im.Set(cx-imageX, cy, color)
		}

		// Clicked a button?
		if prevPressed && !pressed {
			// Done drawing?
			if pressX >= exitX && pressY >= exitY && cx >= exitX && cy >= exitY {
				break
			}
			// New color?
			if pressX < imageX-buttonBuffer &&
				cx < imageX-buttonBuffer && (pressY/4) == (cy/4) {
				color = pal[cy/4]
			}
		}

		prevPressed = pressed
		vert = rand.Intn(3) - 1
		horiz = rand.Intn(3) - 1
		pressed = rand.Intn(2) == 1
	}

	f, err := os.Create("out.png")
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := png.Encode(w, im); err != nil {
		panic(err)
	}
	w.Flush()
}
