// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package gui defines the interface used by artists in the village.
package gui

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/tswast/pixelsketches/palettes"
)

const (
	ScreenWidth  int = 114
	ScreenHeight int = 64
	imageHeight  int = 64
	imageWidth   int = 64
	imageX       int = 25
	buttonHeight int = 4
	buttonBuffer int = 1
	exitX        int = imageX + imageHeight + buttonBuffer
	exitY        int = ScreenHeight - buttonHeight
)

// Struct Cursor is a cursor position and button state.
//
// It remembers where a button press first started, since that information is
// used to determine if button actions will happen or not.
type Cursor struct {
	Pos      image.Point
	Pressed  bool
	PressPos image.Point
}

const (
	MODE_DRAWING = 0
	MODE_DONE    = 1
)

// Struct AppState represents the drawing application state.
type AppState struct {
	Cursor Cursor
	Color  color.Color
	Image  *image.Paletted
	Mode   int
}

// Stuct Action specifies how the app state should change.
type Action struct {
	Pressed    bool
	Horizontal int
	Vertical   int
}

func newImage() *image.Paletted {
	r := image.Rect(0, 0, imageWidth, imageHeight)
	im := image.NewPaletted(r, palettes.PICO8)
	return im
}

// NewAppState creates a new AppState.
func NewAppState() *AppState {
	app := &AppState{}
	app.Image = newImage()
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			app.Image.Set(x, y, palettes.PICO8_BLACK)
		}
	}
	app.Color = app.Image.Palette[0]
	return app
}

// CopyAppState makes a deep copy of an AppState.
func CopyAppState(app *AppState) *AppState {
	out := *app
	// Copy the image.
	out.Image = newImage()
	draw.Draw(out.Image, out.Image.Bounds(), app.Image, image.ZP, draw.Src)
	return &out
}

// ApplyAction modifies an AppState by an action.
func (app *AppState) ApplyAction(act *Action) {
	if app.Mode != MODE_DRAWING {
		return
	}

	// Moving?
	app.Cursor.Pos.Y = act.Vertical + app.Cursor.Pos.Y
	if app.Cursor.Pos.Y < 0 {
		app.Cursor.Pos.Y = 0
	}
	if app.Cursor.Pos.Y >= ScreenHeight {
		app.Cursor.Pos.Y = ScreenHeight - 1
	}
	app.Cursor.Pos.X = act.Horizontal + app.Cursor.Pos.X
	if app.Cursor.Pos.X < 0 {
		app.Cursor.Pos.X = 0
	}
	if app.Cursor.Pos.X >= ScreenWidth {
		app.Cursor.Pos.X = ScreenWidth - 1
	}

	// Just pressed?
	prevPressed := app.Cursor.Pressed
	app.Cursor.Pressed = act.Pressed
	pressed := act.Pressed
	if pressed && !prevPressed {
		app.Cursor.PressPos.X = app.Cursor.Pos.X
		app.Cursor.PressPos.Y = app.Cursor.Pos.Y
	}

	// Drawing?
	if pressed && app.Cursor.PressPos.X >= imageX && app.Cursor.PressPos.X < imageX+imageWidth &&
		app.Cursor.Pos.X >= imageX && app.Cursor.Pos.X < imageX+imageWidth {
		app.Image.Set(app.Cursor.Pos.X-imageX, app.Cursor.Pos.Y, app.Color)
	}

	// Clicked a button?
	if prevPressed && !pressed {
		// Done drawing?
		if app.Cursor.PressPos.X >= exitX && app.Cursor.PressPos.Y >= exitY && app.Cursor.Pos.X >= exitX && app.Cursor.Pos.Y >= exitY {
			app.Mode = MODE_DONE
			return
		}
		// New color?
		if app.Cursor.PressPos.X < imageX-buttonBuffer &&
			app.Cursor.Pos.X < imageX-buttonBuffer && (app.Cursor.PressPos.Y/4) == (app.Cursor.Pos.Y/4) {
			app.Color = app.Image.Palette[app.Cursor.Pos.Y/4]
		}
	}
}
