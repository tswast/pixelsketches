// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package palettes provides color palettes for images.
//
// Individual colors are prefixed by the color palette's name. For example, PICO8_.
package palettes

import (
	"image/color"
)

var PICO8_BLACK = color.RGBA{0, 0, 0, 255}
var PICO8_DARK_BLUE = color.RGBA{29, 43, 83, 255}
var PICO8_DARK_PURPLE = color.RGBA{126, 37, 83, 255}
var PICO8_DARK_GREEN = color.RGBA{0, 135, 81, 255}
var PICO8_BROWN = color.RGBA{171, 82, 54, 255}
var PICO8_DARK_GRAY = color.RGBA{95, 87, 79, 255}
var PICO8_LIGHT_GRAY = color.RGBA{194, 195, 199, 255}
var PICO8_WHITE = color.RGBA{255, 241, 232, 255}
var PICO8_RED = color.RGBA{255, 0, 77, 255}
var PICO8_ORANGE = color.RGBA{255, 163, 0, 255}
var PICO8_YELLOW = color.RGBA{255, 236, 39, 255}
var PICO8_GREEN = color.RGBA{0, 228, 54, 255}
var PICO8_BLUE = color.RGBA{41, 173, 255, 255}
var PICO8_INDIGO = color.RGBA{131, 118, 156, 255}
var PICO8_PINK = color.RGBA{255, 119, 168, 255}
var PICO8_PEACH = color.RGBA{255, 204, 170, 255}

// The PICO8 color palette.
//
// PICO-8 is a "fantasy console" developed by Joseph White (a.k.a. zep). It is
// a  game console / computer platform with carefully constructed limitations,
// including this 16-color palette. http://www.lexaloffle.com/pico-8.php
//
// See: http://www.romanzolotarev.com/pico-8-color-palette/ for a good
// reference chart of the pico-8 color palette.
var PICO8 = []color.Color{
	PICO8_BLACK,
	PICO8_DARK_BLUE,
	PICO8_DARK_PURPLE,
	PICO8_DARK_GREEN,
	PICO8_BROWN,
	PICO8_DARK_GRAY,
	PICO8_LIGHT_GRAY,
	PICO8_WHITE,
	PICO8_RED,
	PICO8_ORANGE,
	PICO8_YELLOW,
	PICO8_GREEN,
	PICO8_BLUE,
	PICO8_INDIGO,
	PICO8_PINK,
	PICO8_PEACH,
}
