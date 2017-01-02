// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/tswast/pixelsketches/palettes"
)

func drawPalette(scr draw.Image, pal []color.Color) {
	scrB := scr.Bounds()
	scrH := scrB.Max.Y - scrB.Min.Y
	clrH := scrH / len(pal)
	for ci, clr := range pal {
		draw.Draw(
			scr,
			image.Rectangle{
				image.Point{0, ci * clrH},
				image.Point{ImageX - ButtonBuffer, (ci + 1) * clrH}},
			&image.Uniform{clr},
			image.ZP,
			draw.Src)
	}
}

func drawColorChoice(scr draw.Image, clr color.Color) {
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{ImageX - ButtonBuffer, 0},
			image.Point{ImageX, ImageHeight}},
		&image.Uniform{clr},
		image.ZP,
		draw.Src)
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{ImageX + ImageWidth, 0},
			image.Point{ImageX + ImageWidth + ButtonBuffer, ImageHeight}},
		&image.Uniform{clr},
		image.ZP,
		draw.Src)
}

func drawTools(scr draw.Image) {
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{ImageX + ImageWidth + ButtonBuffer, 0},
			image.Point{ScreenWidth, ScreenHeight}},
		&image.Uniform{palettes.PICO8_DARK_GRAY},
		image.ZP,
		draw.Src)
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{ImageX + ImageWidth + ButtonBuffer, ScreenHeight - ButtonHeight},
			image.Point{ScreenWidth, ScreenHeight}},
		&image.Uniform{palettes.PICO8_BLACK},
		image.ZP,
		draw.Src)
	// Draw EXIT
	// E
	scr.Set(91, 60, palettes.PICO8_WHITE)
	scr.Set(91, 61, palettes.PICO8_WHITE)
	scr.Set(91, 62, palettes.PICO8_WHITE)
	scr.Set(91, 63, palettes.PICO8_WHITE)
	scr.Set(92, 60, palettes.PICO8_WHITE)
	scr.Set(92, 61, palettes.PICO8_WHITE)
	scr.Set(92, 63, palettes.PICO8_WHITE)
	scr.Set(93, 60, palettes.PICO8_WHITE)
	scr.Set(93, 63, palettes.PICO8_WHITE)
	// X
	scr.Set(95, 60, palettes.PICO8_WHITE)
	scr.Set(95, 61, palettes.PICO8_WHITE)
	scr.Set(95, 63, palettes.PICO8_WHITE)
	scr.Set(96, 61, palettes.PICO8_WHITE)
	scr.Set(96, 62, palettes.PICO8_WHITE)
	scr.Set(97, 60, palettes.PICO8_WHITE)
	scr.Set(97, 62, palettes.PICO8_WHITE)
	scr.Set(97, 63, palettes.PICO8_WHITE)
	// I
	scr.Set(99, 60, palettes.PICO8_WHITE)
	scr.Set(99, 63, palettes.PICO8_WHITE)
	scr.Set(100, 60, palettes.PICO8_WHITE)
	scr.Set(100, 61, palettes.PICO8_WHITE)
	scr.Set(100, 62, palettes.PICO8_WHITE)
	scr.Set(100, 63, palettes.PICO8_WHITE)
	scr.Set(101, 60, palettes.PICO8_WHITE)
	scr.Set(101, 63, palettes.PICO8_WHITE)
	// T
	scr.Set(103, 60, palettes.PICO8_WHITE)
	scr.Set(104, 60, palettes.PICO8_WHITE)
	scr.Set(104, 61, palettes.PICO8_WHITE)
	scr.Set(104, 62, palettes.PICO8_WHITE)
	scr.Set(104, 63, palettes.PICO8_WHITE)
	scr.Set(105, 60, palettes.PICO8_WHITE)
}

// DrawScreen draws the user interface of an app.
func (app *AppState) DrawScreen() *image.NRGBA {
	im := app.Image
	pal := im.Palette
	clr := app.Color
	r := image.Rect(0, 0, ScreenWidth, ScreenHeight)
	scr := image.NewNRGBA(r)
	drawPalette(scr, pal)
	drawColorChoice(scr, clr)
	drawTools(scr)
	draw.Draw(
		scr,
		image.Rectangle{
			image.Point{ImageX, 0},
			image.Point{ImageX + ImageWidth, ImageHeight}},
		im,
		image.ZP,
		draw.Src)
	// Draw cursor
	// Choose a different color every time, so it is easier to track where the
	// cursor is.
	csr := (app.Cursor.Pos.X + app.Cursor.Pos.Y) % len(palettes.PICO8)
	scr.Set(app.Cursor.Pos.X, app.Cursor.Pos.Y, palettes.PICO8[csr])
	return scr
}
