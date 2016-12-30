// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package artist is a bot pixel artist.
package artist

import (
	"bufio"
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"os"

	"github.com/tswast/pixelsketches/village/gui"
)

func tryWriteFrame(frame int, app *gui.AppState) {
	scr := app.DrawScreen()
	// Write timeline image if we can.
	f, err := os.Create(fmt.Sprintf("out/out-%04d.png", frame))
	if err != nil {
		log.Printf("Could not create out/out-%04d.png %s\n", frame, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := png.Encode(w, scr); err != nil {
		log.Printf("Could not encode out/out-%04d.png %s\n", frame, err)
	}
	w.Flush()
}

func randomWalk(a *gui.Action) {
	a.Horizontal = rand.Intn(3) - 1
	a.Vertical = rand.Intn(3) - 1
	a.Pressed = rand.Intn(2) == 1
}

// Main draws a picture, writes it, and exits.
func Main(seed int64, doTimeLapse bool) error {
	rand.Seed(seed)

	app := gui.NewAppState()

	for frame := 0; ; frame++ {
		if app.Mode != gui.MODE_DRAWING {
			break
		}

		if frame%10 == 0 {
			if doTimeLapse {
				tryWriteFrame(frame/10, app)
			}
		}

		action := &gui.Action{}
		randomWalk(action)
		app.ApplyAction(action)
	}

	f, err := os.Create("out.png")
	if err != nil {
		return fmt.Errorf("Error creating out.png: %s", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := png.Encode(w, app.Image); err != nil {
		return fmt.Errorf("Error encoding out.png: %s", err)
	}
	w.Flush()
	return nil
}
