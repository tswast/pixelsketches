// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package strategy

import (
	"testing"

	"github.com/tswast/pixelsketches/palettes"
	"github.com/tswast/pixelsketches/village/gui"
)

func checkNoColors(t *testing.T, app string, action gui.Action, got Rating) {
	if got.rate > 0 || got.reason != "no-different-colors-found" {
		t.Errorf("simPaint(%s, %#v) => %#v, want rate: 0, reason: no-different-colors-found", app, action, got)
	}
}

func TestSimPaintNoDifferentColors(t *testing.T) {
	// No colors to possibly replace.
	app := gui.NewAppState()
	toRight := gui.Action{Horizontal: 1}
	got := simPaint(app, toRight)
	checkNoColors(t, "all black canvas, black selected", toRight, got)

	// Going right, but no non-black pixels to the right.
	app.Cursor.Pos.X = gui.ImageX
	for y := app.Image.Bounds().Min.Y; y < app.Image.Bounds().Max.Y; y++ {
		app.Image.Set(0, y, palettes.PICO8_PINK)
	}
	got = simPaint(app, toRight)
	checkNoColors(t, "pink column, black selected", toRight, got)
}
