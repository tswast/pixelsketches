// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gui

import (
	"image"
	"image/color"
	"testing"

	"github.com/tswast/pixelsketches/palettes"
)

func TestNewAppState(t *testing.T) {
	got := NewAppState()
	if got == nil {
		t.Fatal("Expected *AppState, got nil")
	}
	if got.Image == nil {
		t.Fatal("Expected *image.Paletted, got nil")
	}
	colorInPal := false
	for _, color := range got.Image.Palette {
		if got.Color == color {
			colorInPal = true
			break
		}
	}
	if !colorInPal {
		t.Errorf("Expected %q in got.Image.Palette, but not present", got.Color)
	}
}

func TestCopyAppState(t *testing.T) {
	app := NewAppState()
	app.Image.Set(16, 37, palettes.PICO8_DARK_BLUE)

	got := CopyAppState(app)

	if got == nil {
		t.Fatal("Expected *AppState, got nil")
	}
	if got.Image == nil {
		t.Error("Expected *image.Paletted, got nil")
	}
	app.Image.Set(16, 37, palettes.PICO8_PINK)
	if got.Image.At(16, 37) != palettes.PICO8_DARK_BLUE {
		t.Errorf("got.Image.At(16, 37) => %v, expected %v", got.Image.At(16, 37), palettes.PICO8_DARK_BLUE)
	}
	app.Cursor.Pressed = true
	if got.Cursor.Pressed {
		t.Error("Expected got.Cursor.Pressed == false, got true.")
	}
}

var cursortests = []struct {
	action Action
	prev   Cursor
	next   Cursor
}{
	// Empty action doesn't change cursor position.
	{Action{}, Cursor{}, Cursor{}},
	{
		action: Action{},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 31, Y: 42}},
	},
	{
		action: Action{},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: 0}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: 0}},
	},
	{
		action: Action{},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: ScreenHeight - 1}},
	},
	{
		action: Action{},
		prev:   Cursor{Pos: image.Point{X: 0, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: 0, Y: ScreenHeight - 1}},
	},
	// Moving past the edge of the screen doesn't change cursor position.
	{
		action: Action{Horizontal: -1, Vertical: -1},
		prev:   Cursor{},
		next:   Cursor{},
	},
	{
		action: Action{Horizontal: 1, Vertical: -1},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: 0}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: 0}},
	},
	{
		action: Action{Horizontal: 1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: ScreenHeight - 1}},
	},
	{
		action: Action{Horizontal: -1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: 0, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: 0, Y: ScreenHeight - 1}},
	},
	// Moving along the edge of the screen changes cursor position in the allowed direction.
	{
		action: Action{Horizontal: 1, Vertical: -1},
		prev:   Cursor{},
		next:   Cursor{Pos: image.Point{X: 1}},
	},
	{
		action: Action{Horizontal: -1, Vertical: 1},
		prev:   Cursor{},
		next:   Cursor{Pos: image.Point{Y: 1}},
	},
	{
		action: Action{Horizontal: -1, Vertical: -1},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: 0}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 2, Y: 0}},
	},
	{
		action: Action{Horizontal: 1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: 0}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: 1}},
	},
	{
		action: Action{Horizontal: -1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 2, Y: ScreenHeight - 1}},
	},
	{
		action: Action{Horizontal: 1, Vertical: -1},
		prev:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: ScreenWidth - 1, Y: ScreenHeight - 2}},
	},
	{
		action: Action{Horizontal: 1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: 0, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: 1, Y: ScreenHeight - 1}},
	},
	{
		action: Action{Horizontal: -1, Vertical: -1},
		prev:   Cursor{Pos: image.Point{X: 0, Y: ScreenHeight - 1}},
		next:   Cursor{Pos: image.Point{X: 0, Y: ScreenHeight - 2}},
	},
	// Normal movement in all directions.
	{
		action: Action{Horizontal: -1, Vertical: -1},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 30, Y: 41}},
	},
	{
		action: Action{Horizontal: -1, Vertical: 0},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 30, Y: 42}},
	},
	{
		action: Action{Horizontal: -1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 30, Y: 43}},
	},
	{
		action: Action{Horizontal: 0, Vertical: -1},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 31, Y: 41}},
	},
	{
		action: Action{Horizontal: 0, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 31, Y: 43}},
	},
	{
		action: Action{Horizontal: 1, Vertical: -1},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 32, Y: 41}},
	},
	{
		action: Action{Horizontal: 1, Vertical: 0},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 32, Y: 42}},
	},
	{
		action: Action{Horizontal: 1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next:   Cursor{Pos: image.Point{X: 32, Y: 43}},
	},
	// Button presses work as expected.
	{
		action: Action{Pressed: true},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next: Cursor{
			Pos:      image.Point{X: 31, Y: 42},
			Pressed:  true,
			PressPos: image.Point{X: 31, Y: 42}},
	},
	{
		action: Action{Pressed: true, Horizontal: 1, Vertical: 1},
		prev:   Cursor{Pos: image.Point{X: 31, Y: 42}},
		next: Cursor{
			Pos:      image.Point{X: 32, Y: 43},
			Pressed:  true,
			PressPos: image.Point{X: 32, Y: 43}},
	},
	{
		action: Action{Pressed: false},
		prev: Cursor{
			Pos:      image.Point{X: 31, Y: 42},
			Pressed:  true,
			PressPos: image.Point{X: 31, Y: 42}},
		next: Cursor{
			Pos:      image.Point{X: 31, Y: 42},
			Pressed:  false,
			PressPos: image.Point{X: 31, Y: 42}},
	},
}

func TestApplyActionMovesCursor(t *testing.T) {
	app := NewAppState()
	for _, tt := range cursortests {
		app.Cursor = tt.prev
		app.ApplyAction(&tt.action)
		if app.Cursor != tt.next {
			t.Errorf(
				"(&AppState{Cursor: %v}).ApplyAction(&%v) =>\n\tAppState{Cursor: %v},\twant AppState{Cursor: %v}",
				tt.prev, tt.action, app.Cursor, tt.next)
		}
	}
}

var colortests = []struct {
	cursor   Cursor
	prev     color.Color
	expected color.Color
}{
	// Not pressed means release can't change selected color.
	{Cursor{}, palettes.PICO8_BLACK, palettes.PICO8_BLACK},
	// Releasing after pressing color selects it.
	{
		Cursor{
			Pos:      image.Point{X: 11, Y: 5},
			Pressed:  true,
			PressPos: image.Point{X: 13, Y: 6}},
		palettes.PICO8_BLACK,
		palettes.PICO8_DARK_BLUE,
	},
	// Releasing outside the button doesn't change selected color.
	{
		Cursor{
			Pos:      image.Point{X: 11, Y: 2},
			Pressed:  true,
			PressPos: image.Point{X: 13, Y: 6}},
		palettes.PICO8_BLACK,
		palettes.PICO8_BLACK,
	},
}

func TestApplyActionSelectsColor(t *testing.T) {
	app := NewAppState()
	for _, tt := range colortests {
		app.Cursor = tt.cursor
		app.Color = tt.prev
		app.ApplyAction(&Action{})
		if app.Color != tt.expected {
			t.Errorf(
				"(&AppState{Cursor: %v, Color: %v}).ApplyAction(&Action{}) =>\n\tAppState{Color: %v},\twant %v",
				tt.cursor, tt.prev, app.Color, tt.expected)
		}
	}
}

var painttests = []struct {
	action   Action
	cursor   Cursor
	color    color.Color
	expected color.Color
}{
	// Releasing doesn't paint. Paint only on press down.
	{
		Action{},
		Cursor{Pos: image.Point{X: ImageX + 12, Y: 36}},
		palettes.PICO8_PINK,
		palettes.PICO8_BLACK,
	},
	{
		Action{},
		Cursor{
			Pos:      image.Point{X: ImageX + 12, Y: 36},
			Pressed:  true,
			PressPos: image.Point{X: ImageX + 12, Y: 36},
		},
		palettes.PICO8_PINK,
		palettes.PICO8_BLACK,
	},
	// Pressing with original press off-image doesn't paint.
	{
		Action{Pressed: true},
		Cursor{
			Pos:      image.Point{X: ImageX + 12, Y: 36},
			Pressed:  true,
			PressPos: image.Point{X: ImageX - 12, Y: 36},
		},
		palettes.PICO8_PINK,
		palettes.PICO8_BLACK,
	},
	// Pressing with original press position on image paints.
	{
		Action{Pressed: true},
		Cursor{Pos: image.Point{X: ImageX + 12, Y: 36}},
		palettes.PICO8_PINK,
		palettes.PICO8_PINK,
	},
	{
		Action{Pressed: true},
		Cursor{
			Pos:      image.Point{X: ImageX + 12, Y: 36},
			Pressed:  true,
			PressPos: image.Point{X: ImageX + 12, Y: 36},
		},
		palettes.PICO8_PINK,
		palettes.PICO8_PINK,
	},
	{
		Action{Pressed: true},
		Cursor{
			Pos:      image.Point{X: ImageX + 12, Y: 36},
			Pressed:  true,
			PressPos: image.Point{X: ImageX, Y: 0},
		},
		palettes.PICO8_PINK,
		palettes.PICO8_PINK,
	},
}

func TestApplyActionPaintsColor(t *testing.T) {
	app := NewAppState()
	for _, tt := range painttests {
		app.Cursor = tt.cursor
		app.Color = tt.color
		app.ApplyAction(&tt.action)
		x := app.Cursor.Pos.X - ImageX
		y := app.Cursor.Pos.Y
		color := app.Image.At(x, y)
		if color != tt.expected {
			t.Errorf(
				"(&AppState{Cursor: %v}).ApplyAction(&%v) =>\n\tpixel @ %d %d = %v,\twant %v",
				tt.cursor, tt.action, x, y, color, tt.expected)
		}
	}
}

func TestApplyActionExitsDrawingMode(t *testing.T) {
	app := NewAppState()
	app.Cursor = Cursor{
		Pos:      image.Point{X: 110, Y: 62},
		Pressed:  true,
		PressPos: image.Point{X: 107, Y: 61},
	}
	app.ApplyAction(&Action{})
	if app.Mode != MODE_DONE {
		t.Errorf(
			"(&AppState{Cursor: %v}).ApplyAction(&Action{}) =>\n\tAppState{Mode: %v},\twant AppState{Mode: MODE_DONE}",
			app.Cursor)
	}
}
