// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package strategy defines different strategies for drawing images.
package strategy

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"sort"
	"sync"

	"github.com/tswast/pixelsketches/village/gui"
	"github.com/tswast/pixelsketches/village/perception"
)

// Strategizer chooses the next action.
type Strategizer interface {
	// Strategize chooses the next action.
	//
	// The Rating is returned to help debug its behavior. It should indicate
	// the desirability (rating and distance to event) of the next action and
	// why it was chosen.
	Strategize(*gui.AppState) (gui.Action, Rating)
}

type reason interface {
	explain() string
}

type simpleReason struct {
	reason string
}

func (r *simpleReason) explain() string {
	return r.reason
}

type paintReason struct {
	newColor color.Color
	oldColor color.Color
	pos      image.Point
}

func (r *paintReason) explain() string {
	return fmt.Sprintf("painting with %v over %v @ %v", r.newColor, r.oldColor, r.pos)
}

type Rating struct {
	rate   float64
	dist   int
	reason reason
}

func (r *Rating) String() string {
	return fmt.Sprintf("{rate: %f dist: %d reason: %q}", r.rate, r.dist, r.reason.explain())
}

func imCoordToGuiCoord(pt image.Point) image.Point {
	return image.Point{X: pt.X + gui.ImageX, Y: pt.Y}
}

type RandomWalk struct{}

// RandomWalk chooses the next action completely randomly.
func (_ *RandomWalk) Strategize(_ *gui.AppState) (gui.Action, Rating) {
	return gui.Action{
		Horizontal: rand.Intn(3) - 1,
		Vertical:   rand.Intn(3) - 1,
		Pressed:    rand.Intn(2) == 1,
	}, Rating{}
}

// actionDistance returns the number of actions to move from src to tgt.
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
func simPaint(app *gui.AppState, act gui.Action, rating perception.Rating) Rating {
	actPt := image.Point{
		X: app.Cursor.Pos.X + act.Horizontal,
		Y: app.Cursor.Pos.Y + act.Vertical,
	}

	// Select bounds to search for non-selected color.
	imX := actPt.X - gui.ImageX
	startX := 0
	if act.Horizontal > 0 {
		startX = imX
	}
	if startX < 0 {
		startX = 0
	}
	maxX := gui.ImageWidth - 1
	if act.Horizontal < 0 {
		maxX = imX
	}
	if maxX >= gui.ImageWidth {
		maxX = gui.ImageWidth - 1
	}

	imY := actPt.Y
	startY := 0
	if act.Vertical > 0 {
		startY = imY
	}
	if startY < 0 {
		startY = 0
	}
	maxY := gui.ImageHeight
	if act.Vertical < 0 {
		maxY = imY
	}
	if maxY >= gui.ImageHeight {
		maxY = gui.ImageHeight - 1
	}

	// Try to replace one of each color.
	var colors map[color.Color]image.Point
	colors = make(map[color.Color]image.Point)
	for x := startX; x <= maxX; x++ {
		for y := startY; y <= maxY; y++ {
			clr := app.Image.At(x, y)
			if clr == app.Color {
				continue
			}
			// Found an existing color? Only override the point if the distance
			// to the new point is less than the distance to the old one.
			npt := image.Point{X: x, Y: y}
			guiNpt := imCoordToGuiCoord(npt)
			pt, ok := colors[clr]
			if ok {
				guiPt := imCoordToGuiCoord(pt)
				if actionDistance(actPt, guiPt) <= actionDistance(actPt, guiNpt) {
					continue
				}
			}
			colors[clr] = npt
		}
	}
	max := Rating{rate: -1.0, reason: &simpleReason{"no-different-colors-found"}}
	for clr, pt := range colors {
		// Set the color, rate, then undo. (Should be faster than copying and applying actions.)
		app.Image.Set(pt.X, pt.Y, app.Color)
		rate := rating(app.Image)
		app.Image.Set(pt.X, pt.Y, clr)

		// Distance to move from cursor to point, including this action.
		dist := actionDistance(actPt, imCoordToGuiCoord(pt)) + 1
		// Special cases are needed for distance == 1.
		if dist == 1 {
			if act.Pressed && app.Cursor.Pressed &&
				!(app.Cursor.PressPos.X >= gui.ImageX && app.Cursor.PressPos.X < gui.ImageX+gui.ImageWidth) {
				// Have to release first then press again because press started
				// off-canvas.
				dist += 2
			} else if !act.Pressed {
				// Not pressing, so must take another action to press.
				dist += 1
			}
		}

		if (rate == max.rate && dist < max.dist) || rate > max.rate {
			max.rate = rate
			max.dist = dist
			max.reason = &paintReason{newColor: app.Color, oldColor: clr, pos: pt}
		}
	}
	return max
}

// simChooseColor returns maximum Rating if can choose a color in that direction, otherwise -1.
//
// Also returns the number of actions needed to select the color then paint.
func simChooseColor(app *gui.AppState, act gui.Action, rating perception.Rating) Rating {
	// When can't choose some color?
	// When going right and to the right of the buttons.
	if act.Horizontal > 0 && app.Cursor.Pos.X >= gui.ImageX-gui.ButtonBuffer {
		return Rating{rate: -1, reason: &simpleReason{"no-color-to-right"}}
	}
	actPt := image.Point{X: app.Cursor.Pos.X + act.Horizontal, Y: app.Cursor.Pos.Y + act.Vertical}
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
	// Apply the action to be certain the latest color is chosen.
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
			// Give an additional 2 past the button buffer, since the bot won't
			// stop and release right on the edge.
			simApp.Cursor.Pos.X = gui.ImageX - gui.ButtonBuffer - 2
			simApp.Cursor.Pos.Y = c*gui.ButtonHeight + gui.ButtonHeight/2
		}

		v := simPaint(simApp, drawAct, rating)
		rate := v.rate
		// One action for current action +
		// Distance from cursor after current action to button and click +
		// Distance from button to paint.
		dist := 1 + actionDistance(actPt, simApp.Cursor.Pos) + v.dist
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
			max.reason = v.reason
		}
	}
	return max
}

// simExit returns Rating if can exit in direction, otherwise -1.
func simExit(app *gui.AppState, act gui.Action, rating perception.Rating) (float64, int) {
	rate := -1.0
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
		rate = rating(app.Image)
	}
	// Use distance -1 so that exit is chosen before any other equivalent action.
	return rate, -1
}

// simAction returns the maximum expected Rating for a given action & direction.
//
// Modifies app, so send a copy.
func simAction(app *gui.AppState, act gui.Action, rating perception.Rating) Rating {
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
		return Rating{rate: -1, dist: 0, reason: &simpleReason{"no-op"}}
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
				rate:   rating(simApp.Image),
				dist:   1,
				reason: &simpleReason{"already-painting"},
			}
		}
	}

	max := Rating{rate: -1}

	// What are the actions that are possible in this direction? Always end
	// on a paint or exit button so we can see how the Rating will change.
	// The exit button should be the most likely choice for a Rating that
	// stays exactly the same.

	// Can we reach the exit button in the lower-right corner?
	rate, dist := simExit(gui.CopyAppState(app), act, rating)
	if (rate == max.rate && dist < max.dist) || rate > max.rate {
		max.rate = rate
		max.dist = dist
		max.reason = &simpleReason{"exit"}
	}

	// Can we paint the selected color somewhere different?
	v := simPaint(gui.CopyAppState(app), act, rating)
	if (v.rate == max.rate && v.dist < max.dist) || v.rate > max.rate {
		max = v
	}

	// Can we pick a new color and paint somewhere with that?
	v = simChooseColor(gui.CopyAppState(app), act, rating)
	if (v.rate == max.rate && v.dist < max.dist) || v.rate > max.rate {
		max.rate = v.rate
		max.dist = v.dist
		max.reason = &simpleReason{"choose-color-" + v.reason.explain()}
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

// Actions lets you sort by gui.Action values.
//
// See: https://gobyexample.com/sorting-by-functions
type Actions []gui.Action

func (s Actions) Len() int {
	return len(s)
}

func (s Actions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Actions) Less(i, j int) bool {
	if s[i].Pressed != s[j].Pressed {
		return s[j].Pressed
	}
	if s[i].Vertical != s[j].Vertical {
		return s[i].Vertical < s[j].Vertical
	}
	return s[i].Horizontal < s[j].Horizontal
}

type Ideal struct {
	Rating perception.Rating
}

// Ideal chooses the next action which has the highest expected overall Rating.
func (s *Ideal) Strategize(app *gui.AppState) (gui.Action, Rating) {
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
			v := simAction(gui.CopyAppState(app), a, s.Rating)
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

	var maxActs []gui.Action
	max := Rating{rate: -1.0}
	for k, v := range results {
		if (v.rate == max.rate && v.dist < max.dist) || v.rate > max.rate {
			max = v
			maxActs = []gui.Action{k}
		} else if v.rate == max.rate && v.dist == max.dist {
			maxActs = append(maxActs, k)
		}
	}
	sort.Sort(Actions(maxActs))
	if len(maxActs) == 0 {
		log.Printf("Oops. I didn't find a maximum action.\n")
		return gui.Action{}, max
	}
	maxAct := maxActs[rand.Intn(len(maxActs))]
	return maxAct, max
}
