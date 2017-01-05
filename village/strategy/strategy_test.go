// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package strategy

import (
	"image"
	"testing"

	"github.com/tswast/pixelsketches/palettes"
	"github.com/tswast/pixelsketches/village/gui"
)

func checkNoColors(t *testing.T, app string, action gui.Action, got Rating) {
	cr, ok := got.reason.(*simpleReason)
	if got.rate > 0 || !ok || cr.reason != "no-different-colors-found" {
		t.Errorf("simPaint(%s, %#v) => %#v reason: %q, want rate: 0, reason: no-different-colors-found", app, action, got, got.reason.explain())
	}
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
		got := simPaint(app, dir)
		checkNoColors(t, "AppState{all black canvas, black selected, cursor: centered}", dir, got)

		// Left of screen
		app.Cursor.Pos.X = gui.ImageX - 1
		got = simPaint(app, dir)
		checkNoColors(t, "AppState{all black canvas, black selected, cursor: left}", dir, got)

		// Right of screen
		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth
		got = simPaint(app, dir)
		checkNoColors(t, "AppState{all black canvas, black selected, cursor: right}", dir, got)
	}

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
		got := simPaint(app, dir)
		checkNoColors(t, "AppState{pink column, black selected, cursor: top}", dir, got)

		app.Cursor.Pos.Y = gui.ImageHeight / 2
		got = simPaint(app, dir)
		checkNoColors(t, "AppState{pink column, black selected, cursor: middle}", dir, got)

		app.Cursor.Pos.Y = gui.ImageHeight - 1
		got = simPaint(app, dir)
		checkNoColors(t, "AppState{pink column, black selected, cursor: bottom}", dir, got)
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
		got := simPaint(app, dir)
		checkNoColors(t, "AppState{pink row, black selected, cursor: left}", dir, got)

		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth/2
		got = simPaint(app, dir)
		checkNoColors(t, "AppState{pink row, black selected, cursor: middle}", dir, got)

		app.Cursor.Pos.X = gui.ImageX + gui.ImageWidth - 1
		got = simPaint(app, dir)
		checkNoColors(t, "AppState{pink row, black selected, cursor: right}", dir, got)
	}
}

func checkPaintAtPoint(t *testing.T, app *gui.AppState, im string, action gui.Action, got Rating, expected *Rating) {
	if got.reason.explain() != expected.reason.explain() || got.dist != expected.dist {
		t.Errorf(
			"simPaint(%s @ %v, %#v)\n\t=> dist: %d reason: %q,\n\twant dist: %d, reason: %v",
			im,
			app.Cursor.Pos,
			action,
			got.dist,
			got.reason.explain(),
			expected.dist,
			expected.reason.explain())
	}
}

func newAppStatePinkBlock() *gui.AppState {
	app := gui.NewAppState()
	// Leave a buffer to check that not selecting locations further out.
	for x := 1; x < 6; x++ {
		for y := 1; y < 6; y++ {
			app.Image.Set(x, y, palettes.PICO8_PINK)
		}
	}
	app.Cursor.Pos = imCoordToGuiCoord(image.Point{3, 3})
	app.Color = palettes.PICO8_PINK
	return app
}

func TestSimPaintAtPoint(t *testing.T) {
	// Painting 2 actions away.
	app := newAppStatePinkBlock()
	app.Image.Set(3, 1, palettes.PICO8_BLACK)
	got := simPaint(app, toUp)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, top-middle black",
		toUp,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 1}}})

	app = newAppStatePinkBlock()
	app.Image.Set(4, 1, palettes.PICO8_BLACK)
	got = simPaint(app, toUpRight)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, top-right black",
		toUpRight,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 4, Y: 1}}})

	app = newAppStatePinkBlock()
	app.Image.Set(5, 3, palettes.PICO8_BLACK)
	got = simPaint(app, toRight)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, middle-right black",
		toRight,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 5, Y: 3}}})

	app = newAppStatePinkBlock()
	app.Image.Set(4, 5, palettes.PICO8_BLACK)
	got = simPaint(app, toDownRight)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, bottom-right black",
		toDownRight,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 4, Y: 5}}})

	app = newAppStatePinkBlock()
	app.Image.Set(3, 5, palettes.PICO8_BLACK)
	got = simPaint(app, toDown)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, bottom-middle black",
		toDown,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 5}}})

	app = newAppStatePinkBlock()
	app.Image.Set(2, 5, palettes.PICO8_BLACK)
	got = simPaint(app, toDownLeft)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, bottom-left black",
		toDownLeft,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 2, Y: 5}}})

	app = newAppStatePinkBlock()
	app.Image.Set(1, 3, palettes.PICO8_BLACK)
	got = simPaint(app, toLeft)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, middle-left black",
		toLeft,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 1, Y: 3}}})

	app = newAppStatePinkBlock()
	app.Image.Set(2, 1, palettes.PICO8_BLACK)
	got = simPaint(app, toUpLeft)
	checkPaintAtPoint(
		t,
		app,
		"pink-block, top-left black",
		toUpLeft,
		got,
		&Rating{dist: 2, reason: &paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 2, Y: 1}}})
}
