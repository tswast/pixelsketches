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

func checkPaintAtPoint(t *testing.T, app string, action gui.Action, got Rating, expected *paintReason) {
	pr, ok := got.reason.(*paintReason)
	if !ok || *pr != *expected {
		t.Errorf("simPaint(%s, %#v) => %#v reason: %q, want reason: %v", app, action, got.reason.explain())
	}
}

func TestSimPaint(t *testing.T) {
	// No colors to possibly replace.
	app := gui.NewAppState()
	toRight := gui.Action{Horizontal: 1}
	toDownRight := gui.Action{Horizontal: 1, Vertical: 1}
	toDown := gui.Action{Vertical: 1}
	toDownLeft := gui.Action{Horizontal: -1, Vertical: 1}
	toLeft := gui.Action{Horizontal: -1}
	toUpLeft := gui.Action{Horizontal: -1, Vertical: -1}
	toUp := gui.Action{Vertical: -1}
	toUpRight := gui.Action{Horizontal: 1, Vertical: -1}
	toCenter := gui.Action{}
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

	// Painting 2 actions away.
	app = gui.NewAppState()
	// Leave a buffer to check that not selecting locations further out.
	for x := 1; x < 6; x++ {
		for y := 1; y < 6; y++ {
			app.Image.Set(x, y, palettes.PICO8_PINK)
		}
	}
	app.Image.Set(3, 1, palettes.PICO8_BLACK)
	app.Cursor.Pos.X = 3
	app.Cursor.Pos.Y = 3
	app.Color = palettes.PICO8_PINK
	got := simPaint(app, toUp)
	checkPaintAtPoint(
		t,
		"AppState{pink-block, top square black}",
		toUp,
		got,
		&paintReason{newColor: palettes.PICO8_PINK, oldColor: palettes.PICO8_BLACK, pos: image.Point{X: 3, Y: 1}})
}
