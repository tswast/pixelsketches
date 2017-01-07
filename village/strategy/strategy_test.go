// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package strategy

import (
	"fmt"
	"image"
	"testing"

	"github.com/tswast/pixelsketches/palettes"
	"github.com/tswast/pixelsketches/village/gui"
	"github.com/tswast/pixelsketches/village/perception"
)

func checkRating(t *testing.T, fn string, app *gui.AppState, im string, action gui.Action, got Rating, expected *Rating) {
	if got.reason.explain() != expected.reason.explain() || got.dist != expected.dist {
		press := ""
		if app.Cursor.Pressed {
			press = fmt.Sprintf(" [pressed @ %v]", app.Cursor.PressPos)
		}
		t.Errorf(
			"%s(%s @ %v%s, %#v, rating)\n\t=> dist: %d reason: %q,\n\twant dist: %d, reason: %v",
			fn,
			im,
			app.Cursor.Pos,
			press,
			action,
			got.dist,
			got.reason.explain(),
			expected.dist,
			expected.reason.explain())
	}
}

func checkSimPaint(t *testing.T, app *gui.AppState, im string, action gui.Action, got Rating, expected *Rating) {
	checkRating(t, "simPaint", app, im, action, got, expected)
}

func checkNoColors(t *testing.T, app *gui.AppState, im string, action gui.Action, got Rating) {
	checkSimPaint(t, app, im, action, got, &Rating{rate: -1, dist: 0, reason: &simpleReason{"no-different-colors-found"}})
}

var (
	toRight     gui.Action = gui.Action{Horizontal: 1}
	toDownRight gui.Action = gui.Action{Horizontal: 1, Vertical: 1}
	toDown      gui.Action = gui.Action{Vertical: 1}
	toDownLeft  gui.Action = gui.Action{Horizontal: -1, Vertical: 1}
	toLeft      gui.Action = gui.Action{Horizontal: -1}
	toUpLeft    gui.Action = gui.Action{Horizontal: -1, Vertical: -1}
	toUp        gui.Action = gui.Action{Vertical: -1}
	toUpRight   gui.Action = gui.Action{Horizontal: 1, Vertical: -1}
	toCenter    gui.Action = gui.Action{}
)

func TestSimPaintNoColors(t *testing.T) {
	// No colors to possibly replace.
	app := gui.NewAppState()
	dirs := []gui.Action{toRight, toDownRight, toDown, toDownLeft, toLeft, toUpLeft, toUp, toUpRight, toCenter}
	for _, dir := range dirs {
		// Center of screen
		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth/2
		app.Cursor.Pos.Y = gui.ImageHeight / 2
		got := simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "all black canvas, black selected", dir, got)

		// Left of screen
		app.Cursor.Pos.X = gui.ImageX - 1
		got = simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "all black canvas, black selected", dir, got)

		// Right of screen
		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth
		got = simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "all black canvas, black selected", dir, got)
	}

	// Going away from image when off the image.
	app.Color = palettes.PICO8_PINK
	app.Cursor.Pos = image.Point{10, 63}
	got := simPaint(app, toLeft, perception.RateWholeImage)
	checkNoColors(t, app, "all black canvas, pink selected", toLeft, got)

	app.Cursor.Pos = image.Point{gui.ImageX + gui.ImageWidth + 2, 63}
	got = simPaint(app, toRight, perception.RateWholeImage)
	checkNoColors(t, app, "all black canvas, pink selected", toRight, got)

	// Going horizontally, but no non-black pixels outside current column.
	app = gui.NewAppState()
	for y := app.Image.Bounds().Min.Y; y < app.Image.Bounds().Max.Y; y++ {
		app.Image.Set(gui.ImageWidth/2, y, palettes.PICO8_PINK)
	}
	for _, dir := range dirs {
		if dir.Horizontal == 0 {
			continue
		}
		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth/2
		app.Cursor.Pos.Y = 0
		got := simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "pink column, black selected", dir, got)

		app.Cursor.Pos.Y = gui.ImageHeight / 2
		got = simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "pink column, black selected", dir, got)

		app.Cursor.Pos.Y = gui.ImageHeight - 1
		got = simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "pink column, black selected", dir, got)
	}

	// Going vertically, but no non-black pixels outside the current row.
	app = gui.NewAppState()
	for x := app.Image.Bounds().Min.X; x < app.Image.Bounds().Max.X; x++ {
		app.Image.Set(x, gui.ImageHeight/2, palettes.PICO8_PINK)
	}
	for _, dir := range dirs {
		if dir.Vertical == 0 {
			continue
		}
		app.Cursor.Pos.X = gui.ImageX
		app.Cursor.Pos.Y = gui.ImageHeight / 2
		got := simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "pink row, black selected", dir, got)

		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth/2
		got = simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "pink row, black selected", dir, got)

		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth - 1
		got = simPaint(app, dir, perception.RateWholeImage)
		checkNoColors(t, app, "pink row, black selected", dir, got)
	}
}

// newAppStatePinkBlock makes a new canvase with a 5x5 block centered @ (cx, cy).
func newAppStatePinkBlock(cx, cy int) *gui.AppState {
	app := gui.NewAppState()
	// Leave a buffer to check that not selecting locations further out.
	for x := cx - 2; x < cx+3; x++ {
		for y := cy - 2; y < cy+3; y++ {
			app.Image.Set(x, y, palettes.PICO8_PINK)
		}
	}
	app.Cursor.Pos = imCoordToGuiCoord(image.Point{3, 3})
	app.Color = palettes.PICO8_PINK
	return app
}

func TestSimPaintAtPoint(t *testing.T) {
	// Painting 2 actions away.
	app := newAppStatePinkBlock(3, 3)
	app.Image.Set(3, 1, palettes.PICO8_BLACK)
	got := simPaint(app, toUp, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, top-middle black",
		toUp,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 1}}})

	app = newAppStatePinkBlock(3, 3)
	app.Image.Set(4, 1, palettes.PICO8_BLACK)
	got = simPaint(app, toUpRight, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, top-right black",
		toUpRight,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 4, Y: 1}}})

	app = newAppStatePinkBlock(3, 3)
	app.Image.Set(5, 3, palettes.PICO8_BLACK)
	got = simPaint(app, toRight, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, middle-right black",
		toRight,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 5, Y: 3}}})

	app = newAppStatePinkBlock(3, 3)
	app.Image.Set(4, 5, palettes.PICO8_BLACK)
	got = simPaint(app, toDownRight, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, bottom-right black",
		toDownRight,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 4, Y: 5}}})

	app = newAppStatePinkBlock(3, 3)
	app.Image.Set(3, 5, palettes.PICO8_BLACK)
	got = simPaint(app, toDown, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, bottom-middle black",
		toDown,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 5}}})

	app = newAppStatePinkBlock(3, 3)
	app.Image.Set(2, 5, palettes.PICO8_BLACK)
	got = simPaint(app, toDownLeft, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, bottom-left black",
		toDownLeft,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 2, Y: 5}}})

	app = newAppStatePinkBlock(3, 3)
	app.Image.Set(1, 3, palettes.PICO8_BLACK)
	got = simPaint(app, toLeft, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, middle-left black",
		toLeft,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 1, Y: 3}}})

	app = newAppStatePinkBlock(3, 3)
	app.Image.Set(2, 1, palettes.PICO8_BLACK)
	got = simPaint(app, toUpLeft, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, top-left black",
		toUpLeft,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 2, Y: 1}}})

	// Block on right side of screen.
	app = newAppStatePinkBlock(60, 3)
	app.Image.Set(62, 3, palettes.PICO8_BLACK)
	app.Cursor.Pos.X = gui.ImageX + 60
	app.Cursor.Pos.Y = 3
	got = simPaint(app, toRight, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block @ (60, 3), middle-right black",
		toRight,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 62, Y: 3}}})

	// Painting 1 action away, but not this action, so 2 or more actions away.
	app = newAppStatePinkBlock(3, 3)
	app.Cursor.Pos.X = gui.ImageX + 3
	app.Cursor.Pos.Y = 3
	app.Image.Set(2, 2, palettes.PICO8_BLACK)
	got = simPaint(app, toUp, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, top-left dist=1 black",
		toUp,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 2, Y: 2}}})

	// Distance 1 away, but not pressing.
	app = newAppStatePinkBlock(3, 3)
	app.Cursor.Pos.X = gui.ImageX + 3
	app.Cursor.Pos.Y = 3
	app.Image.Set(3, 2, palettes.PICO8_BLACK)
	got = simPaint(app, toUp, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, top-middle dist=1, black",
		toUp,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 2}}})

	// Distance 1 away, including this action, but have to release then paint again. So, actually distance 3.
	app = newAppStatePinkBlock(3, 3)
	app.Cursor.Pos.X = gui.ImageX + 3
	app.Cursor.Pos.Y = 3
	// Pressed off-canvas
	app.Cursor.PressPos.X = 0
	app.Cursor.PressPos.Y = 0
	app.Cursor.Pressed = true
	app.Image.Set(3, 2, palettes.PICO8_BLACK)
	act := gui.Action{Vertical: -1, Pressed: true}
	got = simPaint(app, act, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, top-middle dist=1, black",
		act,
		got,
		&Rating{dist: 3, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 2}}})
	app.Cursor.Pressed = false

	// Distance 1 away, including this action. (Same as "still painting")
	app = newAppStatePinkBlock(3, 3)
	app.Cursor.Pos.X = gui.ImageX + 3
	app.Cursor.Pos.Y = 3
	app.Image.Set(3, 2, palettes.PICO8_BLACK)
	act = gui.Action{Vertical: -1, Pressed: true}
	got = simPaint(app, act, perception.RateWholeImage)
	checkSimPaint(
		t,
		app,
		"pink-block, top-middle dist=1, black",
		act,
		got,
		&Rating{dist: 1, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 2}}})
}

func checkSimChooseColor(t *testing.T, app *gui.AppState, im string, action gui.Action, got Rating, expected *Rating) {
	checkRating(t, "simChooseColor", app, im, action, got, expected)
}

func TestSimChooseColor(t *testing.T) {
	app := gui.NewAppState()
	app.Cursor.Pos.X = gui.ImageX
	app.Cursor.Pos.Y = gui.ButtonHeight / 2
	app.Image.Set(3, gui.ButtonHeight/2, palettes.PICO8_ORANGE)
	app.Color = palettes.PICO8_PINK
	act := toLeft
	got := simChooseColor(app, act, func(im image.Image) float64 {
		return perception.RateImage(im, palettes.PICO8_BLACK, 1.0)
	})
	checkSimChooseColor(
		t,
		app,
		"all-black, except pink @ (3, 2), want all-black",
		act,
		got,
		&Rating{
			// 3 from clicking new color. 6 to paint pixel from button.
			dist: 3 + 6,
			reason: &paintReason{
				newColor: palettes.PICO8_BLACK,
				oldColor: palettes.PICO8_ORANGE,
				pos:      image.Point{X: 3, Y: gui.ButtonHeight / 2}}})

	// Moving up or down, which doesn't get you closer to the button.
	app = gui.NewAppState()
	app.Cursor.Pos.X = gui.ImageX
	app.Cursor.Pos.Y = gui.ButtonHeight / 2
	app.Image.Set(3, gui.ButtonHeight/2, palettes.PICO8_ORANGE)
	app.Color = palettes.PICO8_PINK
	act = toUp
	got = simChooseColor(app, act, func(im image.Image) float64 {
		return perception.RateImage(im, palettes.PICO8_BLACK, 1.0)
	})
	checkSimChooseColor(
		t,
		app,
		"all-black, except pink @ (3, 2), want all-black",
		act,
		got,
		&Rating{
			// 1 for going wrong direction. 3 from clicking new color. 6 to paint pixel from button.
			dist: 1 + 3 + 6,
			reason: &paintReason{
				newColor: palettes.PICO8_BLACK,
				oldColor: palettes.PICO8_ORANGE,
				pos:      image.Point{X: 3, Y: gui.ButtonHeight / 2}}})
}
