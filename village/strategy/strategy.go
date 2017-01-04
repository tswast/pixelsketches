// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package strategy defines different strategies for drawing images.
package strategy

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sync"

	"github.com/tswast/pixelsketches/village/gui"
	"github.com/tswast/pixelsketches/village/perception"
)

// A Strategy returns a new action for a given app state.
type Strategy func(*gui.AppState) (gui.Action, Rating)

// RandomWalk chooses the next action completely randomly.
func RandomWalk(_ *gui.AppState) (gui.Action, Rating) {
	return gui.Action{
		Horizontal: rand.Intn(3) - 1,
		Vertical:   rand.Intn(3) - 1,
		Pressed:    rand.Intn(2) == 1,
	}, Rating{}
}

func actionDistance(src, tgt image.Point) int {
	h := tgt.X - src.X
	if h < 0 {
		h *= -1
	}
	v := tgt.Y - src.Y
	if v < 0 {
		v *= -1
	}
	// Return the maximum dimension, since the target can be reached by
	// diagonal actions and actions in just that direction.
	if h > v {
		return h
	}
	return v
}

// simPaint returns maximum Rating if can paint in direction, otherwise -1.
//
// Also, return the minimum number of actions needed to get to that position and paint.
func simPaint(app *gui.AppState, act gui.Action) Rating {
	// Select bounds to search for non-selected color.
	imX := app.Cursor.Pos.X - gui.ImageX + act.Horizontal
	if imX < 0 {
		imX = 0
	}
	if imX >= gui.ImageWidth {
		imX = gui.ImageWidth - 1
	}
	startX := 0
	maxX := gui.ImageX
	if act.Horizontal > 0 {
		startX = imX
	}
	if act.Horizontal < 0 {
		maxX = imX
	}
	imY := app.Cursor.Pos.Y + act.Vertical
	startY := 0
	maxY := gui.ImageHeight
	if act.Vertical > 0 {
		startY = imY
	}
	if act.Vertical < 0 {
		maxY = imY
	}

	// Try to replace one of each color.
	var colors map[color.Color]image.Point
	colors = make(map[color.Color]image.Point)
	for x := startX; x < maxX; x++ {
		for y := startY; y < maxY; y++ {
			clr := app.Image.At(x, y)
			if clr == app.Color {
				continue
			}
			// Found an existing color? Only override the point if the distance
			// to the new point is less than the distance to the old one.
			tgt := image.Point{X: x, Y: y}
			pt, ok := colors[clr]
			if ok && actionDistance(app.Cursor.Pos, pt) <= actionDistance(app.Cursor.Pos, tgt) {
				continue
			}
			colors[clr] = tgt
		}
	}
	max := Rating{rate: -1.0, reason: "no-different-colors-found"}
	for clr, pt := range colors {
		// Set the color, rate, then undo. (Should be faster than copying and applying actions.)
		app.Image.Set(pt.X, pt.Y, app.Color)
		rate := perception.RateWholeImage(app.Image)
		app.Image.Set(pt.X, pt.Y, clr)
		dist := actionDistance(app.Cursor.Pos, pt)
		if (rate == max.rate && dist < max.dist) || rate > max.rate {
			max.rate = rate
			max.dist = dist
			max.reason = fmt.Sprintf("selected: %v overwrite: %v pos: %v", app.Color, clr, pt)
		}
	}
	return max
}

// simChooseColor returns maximum Rating if can choose a color in that direction, otherwise -1.
//
// Also returns the number of actions needed to select the color then paint.
func simChooseColor(app *gui.AppState, act gui.Action) Rating {
	// When can't choose some color?
	// When going right and to the right of the buttons.
	if act.Horizontal > 0 && app.Cursor.Pos.X >= gui.ImageX-gui.ButtonBuffer {
		return Rating{rate: -1, reason: "no-color-to-right"}
	}
	drawAct := gui.Action{Horizontal: 1}

	// Which colors can we select in this direction?
	cMin := 0
	cMax := len(app.Image.Palette) - 1
	if act.Vertical < 0 {
		cMax = app.Cursor.Pos.Y / gui.ButtonHeight
	} else if act.Vertical > 0 {
		cMin = app.Cursor.Pos.Y / gui.ButtonHeight
	}

	// Which colors could we pick?
	max := Rating{rate: -1}
	simApp := gui.CopyAppState(app)
	simApp.ApplyAction(&act)
	for c := cMin; c <= cMax; c++ {
		// Try drawing from each color choice.
		// Skip the color that was previously picked. The actions to draw with
		// that color will already be considered in simPaint.
		if app.Image.Palette[c] == app.Color {
			continue
		}

		justSelected := simApp.Color == app.Image.Palette[c]
		if !justSelected {
			simApp.Color = app.Image.Palette[c]
			// Give an additional -1 button buffer, since releasing without
			// moving is not an action the bot considers.
			simApp.Cursor.Pos.X = gui.ImageX - gui.ButtonBuffer - 1
			simApp.Cursor.Pos.Y = c*gui.ButtonHeight + gui.ButtonHeight/2
		}

		v := simPaint(simApp, drawAct)
		rate := v.rate
		dist := v.dist + actionDistance(app.Cursor.Pos, simApp.Cursor.Pos) + 1
		// Add an action to click the button if we aren't pressing. Release
		// will happen on the move out, on the button boundary.
		if ((app.Cursor.Pos.Y/gui.ButtonHeight) == c && app.Cursor.Pos.X < simApp.Cursor.Pos.X && act.Pressed && !app.Cursor.Pressed) ||
			// Just released the button.
			justSelected {
			// Remove the extra action if already clicked the button.
			dist -= 1
		}
		if (rate == max.rate && dist < max.dist) || rate > max.rate {
			max.rate = rate
			max.dist = dist
			max.reason = fmt.Sprintf("color-%d", c)
		}
	}
	return max
}

// simExit returns Rating if can exit in direction, otherwise -1.
func simExit(app *gui.AppState, act gui.Action) (float64, int) {
	rate := -1.0
	dist := 0
	// Going right.
	if (act.Horizontal > 0 && act.Vertical == 0) ||
		// Going down-right.
		(act.Horizontal > 0 && act.Vertical > 0) ||
		// Going down, but not down-left.
		(act.Horizontal == 0 && act.Vertical > 0) ||
		// Going in any direction when completely in the button boundary.
		(app.Cursor.Pos.X > gui.ExitX && app.Cursor.Pos.Y > gui.ExitY) {
		// The Rating for choosing the exit action is whatever Rating the image
		// would get now.
		app.ApplyAction(&act)
		rate = perception.RateWholeImage(app.Image)

		// The distance is only 1 if this action causes an exit.
		if app.Mode != gui.MODE_DRAWING {
			dist = 1
		} else if app.Cursor.Pos.X >= gui.ExitX && app.Cursor.Pos.Y >= gui.ExitY {
			if app.Cursor.Pressed && app.Cursor.PressPos.X >= gui.ExitX && app.Cursor.PressPos.Y >= gui.ExitY {
				dist = 2
			} else {
				dist = 3
			}
		} else {
			// The distance is the distance to the button (and press when reaching) + 1 to release.
			dist = actionDistance(app.Cursor.Pos, image.Point{gui.ExitX, gui.ExitY}) + 1
		}
	}
	return rate, dist
}

// simAction returns the maximum expected Rating for a given action & direction.
//
// Modifies app, so send a copy.
func simAction(app *gui.AppState, act gui.Action) Rating {
	// Can't move left from the left edge of the screen.
	if (app.Cursor.Pos.X <= 0 && act.Horizontal < 0) ||
		// Can't move right from the right edge of the screen.
		(app.Cursor.Pos.X >= gui.ScreenWidth-1 && act.Horizontal > 0) ||
		// Can't move up from top edge of the screen.
		(app.Cursor.Pos.Y <= 0 && act.Vertical < 0) ||
		// Can't move down from bottom edge of the screen.
		(app.Cursor.Pos.Y >= gui.ScreenHeight-1 && act.Vertical > 0) ||
		// There is nothing to click in the upper-right quadrant once outside of the image.
		(app.Cursor.Pos.X >= gui.ImageX+gui.ImageWidth && app.Cursor.Pos.Y <= gui.ExitY && act.Horizontal > 0 && act.Vertical < 0) {
		// Return -1 to discourage from picking this action.
		return Rating{rate: -1, dist: 0, reason: "no-op"}
	}

	// Already painting this action? Return the new Rating. Don't simulate
	// anything else since already painted once for this action.
	if act.Pressed {
		simApp := gui.CopyAppState(app)
		simApp.ApplyAction(&act)
		imX := simApp.Cursor.Pos.X - gui.ImageX
		imY := simApp.Cursor.Pos.Y
		if imX >= 0 && imX < gui.ImageWidth &&
			simApp.Image.At(imX, imY) != app.Image.At(imX, imY) {
			return Rating{
				rate:   perception.RateWholeImage(simApp.Image),
				dist:   1,
				reason: "already-painting",
			}
		}
	}

	max := Rating{rate: -1}

	// What are the actions that are possible in this direction? Always end
	// on a paint or exit button so we can see how the Rating will change.
	// The exit button should be the most likely choice for a Rating that
	// stays exactly the same.

	// Can we reach the exit button in the lower-right corner?
	rate, dist := simExit(gui.CopyAppState(app), act)
	if (rate == max.rate && dist < max.dist) || rate > max.rate {
		max.rate = rate
		max.dist = dist
		max.reason = "exit"
	}

	// Can we paint the selected color somewhere different?
	v := simPaint(gui.CopyAppState(app), act)
	if (v.rate == max.rate && v.dist < max.dist) || v.rate > max.rate {
		max.rate = v.rate
		max.dist = v.dist
		max.reason = "paint-" + v.reason
	}

	// Can we pick a new color and paint somewhere with that?
	v = simChooseColor(gui.CopyAppState(app), act)
	if (v.rate == max.rate && v.dist < max.dist) || v.rate > max.rate {
		max.rate = v.rate
		max.dist = v.dist
		max.reason = "choose-color-" + v.reason
	}
	return max
}

var directions = []struct {
	h int
	v int
}{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	// Skip the middle so it doesn't get stuck (no movement)
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

type Rating struct {
	rate   float64
	dist   int
	reason string
}

// Ideal chooses the next action which has the highest expected overall Rating.
func Ideal(app *gui.AppState) (gui.Action, Rating) {
	var results map[gui.Action]Rating
	results = make(map[gui.Action]Rating)

	// Check each possible action and do the one with the highest expected value.
	var wg sync.WaitGroup
	lock := sync.Mutex{}
	for _, dir := range directions {
		wg.Add(2)
		a := gui.Action{
			Horizontal: dir.h,
			Vertical:   dir.v,
		}
		calculateResult := func(a gui.Action) {
			defer wg.Done()
			v := simAction(gui.CopyAppState(app), a)
			// Write results.
			lock.Lock()
			results[a] = v
			lock.Unlock()
		}
		go calculateResult(a)
		a.Pressed = true
		go calculateResult(a)
	}
	wg.Wait()

	var maxAct gui.Action
	max := Rating{rate: -1.0}
	for k, v := range results {
		if (v.rate == max.rate && v.dist < max.dist) || v.rate > max.rate {
			max = v
			maxAct = k
		}
	}
	return maxAct, max
}
