// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/tswast/pixelsketches/village/gui"
)

type ActionBallot struct {
	Pressed    []bool
	Horizontal []int
	Vertical   []int
}

var (
	targetX int
	targetY int
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

func goToPoint(x, y, tx, ty int) *ActionBallot {
	b := &ActionBallot{}
	if x == tx && y == ty {
		return b
	}
	if y < ty {
		b.Vertical = append(b.Vertical, 1)
	} else if y == ty {
		b.Vertical = append(b.Vertical, 0)
	} else {
		b.Vertical = append(b.Vertical, -1)
	}
	if x < tx {
		b.Horizontal = append(b.Horizontal, 1)
	} else if x == tx {
		b.Horizontal = append(b.Horizontal, 0)
	} else {
		b.Horizontal = append(b.Horizontal, -1)
	}
	return b
}

func randomWalk(a *gui.Action) {
	a.Horizontal = rand.Intn(3) - 1
	a.Vertical = rand.Intn(3) - 1
	a.Pressed = rand.Intn(2) == 1
}

func strategize(a *gui.Action, x, y int) {
	b := goToPoint(x, y, targetX, targetY)
	if len(b.Horizontal) == 0 || len(b.Vertical) == 0 {
		a.Horizontal = 0
		a.Vertical = 0
		targetX = rand.Intn(gui.ScreenWidth)
		targetY = rand.Intn(gui.ScreenHeight)
		fmt.Printf("got to point %d %d, new target %d %d\n", x, y, targetX, targetY)
	} else {
		a.Horizontal = b.Horizontal[rand.Intn(len(b.Horizontal))]
		a.Vertical = b.Vertical[rand.Intn(len(b.Vertical))]
	}
	a.Pressed = rand.Intn(2) == 1
}

func main() {
	var seed int64
	seed = 19700101
	if len(os.Args) > 2 {
		log.Fatalf("Got unexpected number of arguments %d\n", len(os.Args)-1)
	}
	if len(os.Args) == 2 {
		var err error
		seed, err = strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatalf("Error parsing seed %s\n", err)
		}
	}
	rand.Seed(seed)

	// json.Unmarshal()
	app := gui.NewAppState()

	for frame := 0; ; frame++ {
		if app.Mode != gui.MODE_DRAWING {
			break
		}

		if frame%10 == 0 {
			//	tryWriteFrame(frame/10, app)
		}

		action := &gui.Action{}
		//strategize(action, cx, cy)
		randomWalk(action)
		app.ApplyAction(action)
	}

	f, err := os.Create("out.png")
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := png.Encode(w, app.Image); err != nil {
		panic(err)
	}
	w.Flush()
}
