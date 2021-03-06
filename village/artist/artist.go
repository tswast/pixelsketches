// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package artist is a bot pixel artist.
package artist

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
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
func Main(inPath, outPath string, seed int64, debug, doTimeLapse bool, maxIter int, s strategy.Strategizer) error {
	rand.Seed(seed)

	app := gui.NewAppState()
	if inPath != "" {
		f, err := os.Open(inPath)
		if err != nil {
			log.Fatalf("Error opening %s: %s", inPath, err)
		}
		defer f.Close()
		im, err := png.Decode(f)
		if err != nil {
			log.Fatalf("Error decoding %s: %s", inPath, err)
		}
		draw.Draw(app.Image, app.Image.Bounds(), im, image.ZP, draw.Src)
	}

	pts := make(map[image.Point]int)
	frame := 0
	for ; ; frame++ {
		if frame > maxIter {
			log.Printf("reached max iterations %d\n", maxIter)
			break
		}
		if app.Mode != gui.MODE_DRAWING {
			break
		}

		if doTimeLapse {
			tryWriteFrame(frame, app)
		}
		if frame%100 == 0 {
			log.Printf("current-frame: %d\n", frame)
		}

		a, r := s.Strategize(app)
		if debug {
			log.Printf(
				"frame: %d\n\tpos: %v\n\timPos: %v\n\tcolor: %v\n\taction: %v\n\trating: %s\n",
				frame,
				app.Cursor.Pos,
				image.Point{X: app.Cursor.Pos.X - gui.ImageX, Y: app.Cursor.Pos.Y},
				app.Color,
				a,
				r.String())
			if frame%100 == 0 {
				f, err := os.Create(outPath)
				if err != nil {
					log.Printf("Error creating %s: %s\n", outPath, err)
				}
				w := bufio.NewWriter(f)
				if err := png.Encode(w, app.Image); err != nil {
					log.Printf("Error encoding %s: %s\n", outPath, err)
				}
				w.Flush()
				f.Close()
			}
		}
		// Stop if we've been at this exact same point before.
		dejavu, ok := pts[app.Cursor.Pos]
		if !ok {
			pts[app.Cursor.Pos] = frame
		} else if frame-dejavu > 20 {
			log.Printf("already been at this position")
			break
		}
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
