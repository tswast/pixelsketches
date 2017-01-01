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
	"github.com/tswast/pixelsketches/village/strategy"
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

// Main draws a picture, writes it, and exits.
func Main(outPath string, seed int64, doTimeLapse bool, maxIter int, s strategy.Strategy) error {
	rand.Seed(seed)

	app := gui.NewAppState()

	frame := 0
	for ; ; frame++ {
		if frame > maxIter {
			log.Printf("reached max iterations %d\n", maxIter)
			break
		}
		if app.Mode != gui.MODE_DRAWING {
			break
		}

		if frame%10 == 0 {
			if doTimeLapse {
				tryWriteFrame(frame/10, app)
			}
		}

		a := s(app)
		app.ApplyAction(&a)
	}
	fmt.Printf("frames: %d\n", frame)

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("Error creating %s: %s", outPath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := png.Encode(w, app.Image); err != nil {
		return fmt.Errorf("Error encoding %s: %s", outPath, err)
	}
	w.Flush()
	return nil
}
