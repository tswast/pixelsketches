// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
)

var black = color.RGBA{0, 0, 0, 255}
var darkBlue = color.RGBA{29, 43, 83, 255}
var darkPurple = color.RGBA{126, 37, 83, 255}
var darkGreen = color.RGBA{0, 135, 81, 255}
var brown = color.RGBA{171, 82, 54, 255}
var darkGray = color.RGBA{95, 87, 79, 255}
var lightGray = color.RGBA{194, 195, 199, 255}
var white = color.RGBA{255, 241, 232, 255}
var red = color.RGBA{255, 0, 77, 255}
var orange = color.RGBA{255, 163, 0, 255}
var yellow = color.RGBA{255, 236, 39, 255}
var green = color.RGBA{0, 228, 54, 255}
var blue = color.RGBA{41, 173, 255, 255}
var indigo = color.RGBA{131, 118, 156, 255}
var pink = color.RGBA{255, 119, 168, 255}
var peach = color.RGBA{255, 204, 170, 255}

const (
	screenWidth  int = 114
	screenHeight int = 64
	imageHeight  int = 64
	imageWidth   int = 64
	imageX       int = 25
	buttonHeight int = 4
	buttonBuffer int = 1
	exitX        int = imageX + imageHeight + buttonBuffer
	exitY        int = screenHeight - buttonHeight
)

type ActionBallot struct {
	Pressed    []bool
	Horizontal []int
	Vertical   []int
}

type Action struct {
	Pressed    bool
	Horizontal int
	Vertical   int
}

var (
	targetX int
	targetY int
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func drawPalette(scr draw.Image, pal []color.Color) {
	scrB := scr.Bounds()
	scrH := scrB.Max.Y - scrB.Min.Y
	clrH := scrH / len(pal)
	for ci, clr := range pal {
		draw.Draw(
			scr,
			image.Rectangle{
				image.Point{0, ci * clrH},
				image.Point{imageX - buttonBuffer, (ci + 1) * clrH}},
			&image.Uniform{clr},
			image.ZP,
			draw.Src)
	}
}

func drawColorChoice(scr draw.Image, clr color.Color) {
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{imageX - buttonBuffer, 0},
			image.Point{imageX, imageHeight}},
		&image.Uniform{clr},
		image.ZP,
		draw.Src)
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{imageX + imageWidth, 0},
			image.Point{imageX + imageWidth + buttonBuffer, imageHeight}},
		&image.Uniform{clr},
		image.ZP,
		draw.Src)
}

func drawTools(scr draw.Image) {
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{imageX + imageWidth + buttonBuffer, 0},
			image.Point{screenWidth, screenHeight}},
		&image.Uniform{darkGray},
		image.ZP,
		draw.Src)
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{imageX + imageWidth + buttonBuffer, screenHeight - buttonHeight},
			image.Point{screenWidth, screenHeight}},
		&image.Uniform{black},
		image.ZP,
		draw.Src)
	// Draw EXIT
	// E
	scr.Set(91, 60, white)
	scr.Set(91, 61, white)
	scr.Set(91, 62, white)
	scr.Set(91, 63, white)
	scr.Set(92, 60, white)
	scr.Set(92, 61, white)
	scr.Set(92, 63, white)
	scr.Set(93, 60, white)
	scr.Set(93, 63, white)
	// X
	scr.Set(95, 60, white)
	scr.Set(95, 61, white)
	scr.Set(95, 63, white)
	scr.Set(96, 61, white)
	scr.Set(96, 62, white)
	scr.Set(97, 60, white)
	scr.Set(97, 62, white)
	scr.Set(97, 63, white)
	// I
	scr.Set(99, 60, white)
	scr.Set(99, 63, white)
	scr.Set(100, 60, white)
	scr.Set(100, 61, white)
	scr.Set(100, 62, white)
	scr.Set(100, 63, white)
	scr.Set(101, 60, white)
	scr.Set(101, 63, white)
	// T
	scr.Set(103, 60, white)
	scr.Set(104, 60, white)
	scr.Set(104, 61, white)
	scr.Set(104, 62, white)
	scr.Set(104, 63, white)
	scr.Set(105, 60, white)
}

func tryWriteFrame(frame int, im image.Image, clr color.Color, pal []color.Color, x, y int) {
	r := image.Rect(0, 0, screenWidth, screenHeight)
	scr := image.NewNRGBA(r)
	drawPalette(scr, pal)
	drawColorChoice(scr, clr)
	drawTools(scr)
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{imageX, 0},
			image.Point{imageX + imageWidth, imageHeight}},
		im,
		image.ZP,
		draw.Src)
	// Draw cursor
	scr.Set(x, y, white)
	// Write timeline image if we can.
	f, err := os.Create(fmt.Sprintf("out/out-%04d.png", frame))
	if err != nil {
		log.Printf("Could not create out/out-%04d.png %s\n", frame, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := png.Encode(w, scr); err != nil {
		log.Printf("Could not encode out/out-%04d.png %s\n", frame, err)
	}
	w.Flush()
}

func goToPoint(x, y, tx, ty int) *ActionBallot {
	b := &ActionBallot{}
	if x == tx && y == ty {
		return b
	}
	if y < ty {
		b.Vertical = append(b.Vertical, 1)
	} else if y == ty {
		b.Vertical = append(b.Vertical, 0)
	} else {
		b.Vertical = append(b.Vertical, -1)
	}
	if x < tx {
		b.Horizontal = append(b.Horizontal, 1)
	} else if x == tx {
		b.Horizontal = append(b.Horizontal, 0)
	} else {
		b.Horizontal = append(b.Horizontal, -1)
	}
	return b
}

func randomWalk(a *Action, x, y int) {
	a.Horizontal = rand.Intn(3) - 1
	a.Vertical = rand.Intn(3) - 1
	a.Pressed = rand.Intn(2) == 1
}

func strategize(a *Action, x, y int) {
	b := goToPoint(x, y, targetX, targetY)
	if len(b.Horizontal) == 0 || len(b.Vertical) == 0 {
		a.Horizontal = 0
		a.Vertical = 0
		targetX = rand.Intn(screenWidth)
		targetY = rand.Intn(screenHeight)
		fmt.Printf("got to point %d %d, new target %d %d\n", x, y, targetX, targetY)
	} else {
		a.Horizontal = b.Horizontal[rand.Intn(len(b.Horizontal))]
		a.Vertical = b.Vertical[rand.Intn(len(b.Vertical))]
	}
	a.Pressed = rand.Intn(2) == 1
}

func main() {
	var seed int64
	seed = 19700101
	if len(os.Args) > 2 {
		log.Fatalf("Got unexpected number of arguments %d\n", len(os.Args)-1)
	}
	if len(os.Args) == 2 {
		var err error
		seed, err = strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatalf("Error parsing seed %s\n", err)
		}
	}
	rand.Seed(seed)

	// json.Unmarshal()

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
	targetX = imageX + rand.Intn(imageWidth)
	targetY = rand.Intn(imageHeight)
	color := pal[rand.Intn(len(pal))]
	pressX := 0
	pressY := 0
	prevPressed := false
	action := &Action{}
	for frame := 0; ; frame++ {
		// Moving?
		cy = action.Vertical + cy
		if cy < 0 {
			cy = 0
		}
		if cy >= screenHeight {
			cy = screenHeight - 1
		}
		cx = action.Horizontal + cx
		if cx < 0 {
			cx = 0
		}
		if cx >= screenWidth {
			cx = screenWidth - 1
		}

		// Just pressed?
		pressed := action.Pressed
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

		if frame%10 == 0 {
			tryWriteFrame(frame/10, im, color, pal, cx, cy)
		}

		prevPressed = pressed
		//strategize(action, cx, cy)
		randomWalk(action, cx, cy)
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
