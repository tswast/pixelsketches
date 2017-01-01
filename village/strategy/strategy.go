// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package strategy defines different strategies for drawing images.
package strategy

import (
	"math/rand"
	"sync"

	"github.com/tswast/pixelsketches/village/gui"
	"github.com/tswast/pixelsketches/village/perception"
)

const (
	// Number of simulations to do per possible action.
	simDepth int = 64
	// Choose simLength so it is possible to go from one
	// corner of the screen to the other.
	simLength int = gui.ScreenWidth + gui.ScreenHeight
)

// A Strategy returns a new action for a given app state.
type Strategy func(*gui.AppState) gui.Action

// RandomWalk chooses the next action completely randomly.
func RandomWalk(_ *gui.AppState) gui.Action {
	return gui.Action{
		Horizontal: rand.Intn(3) - 1,
		Vertical:   rand.Intn(3) - 1,
		Pressed:    rand.Intn(2) == 1,
	}
}

// simAction returns the expected rating for a given action.
func simAction(app *gui.AppState, act gui.Action) float64 {
	rate := 0.0
	for sim := 1; sim <= simDepth; sim++ {
		simApp := gui.CopyAppState(app)
		simApp.ApplyAction(&act)
		// Repeat the move in the same direction for subsequent simulations.
		a := gui.Action{
			Horizontal: act.Horizontal * sim,
			Vertical:   act.Vertical * sim,
		}
		simApp.ApplyAction(&a)
		for i := 2; i < simLength-sim; i++ {
			if simApp.Mode != gui.MODE_DRAWING {
				break
			}
			a = RandomWalk(simApp)
			simApp.ApplyAction(&a)
		}
		rate += perception.RateWholeImage(simApp.Image)
	}
	return rate / float64(simDepth)
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

// Ideal chooses the next action which has the highest expected overall rating.
func Ideal(app *gui.AppState) gui.Action {
	var results map[gui.Action]float64
	results = make(map[gui.Action]float64)

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
	max := -1.0 // Ratings are always in [0, 1]
	for k, v := range results {
		if v > max {
			max = v
			maxAct = k
		}
	}
	return maxAct
}
