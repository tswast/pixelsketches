// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/tswast/pixelsketches/palettes"
	"github.com/tswast/pixelsketches/village/perception"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Got unexpected number of arguments %d\n", len(os.Args)-1)
	}
	p := os.Args[1]
	f, err := os.Open(p)
	if err != nil {
		log.Fatalf("Error opening %s: %s", p, err)
	}
	defer f.Close()
	im, err := png.Decode(f)
	if err != nil {
		log.Fatalf("Error decoding %s: %s", p, err)
	}

	rt := 0.0
	// Interests were randomly generated using:
	// for i := 0; i < 16; i++ {
	// 	fmt.Printf("%f\n", rand.Float32())
	// }
	r := perception.RateImage(im, palettes.PICO8_BLACK, 0.604660)
	fmt.Printf("black: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_DARK_BLUE, 0.940509)
	fmt.Printf("dark-blue: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_DARK_PURPLE, 0.664560)
	fmt.Printf("dark-purple: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_DARK_GREEN, 0.437714)
	fmt.Printf("dark-green: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_BROWN, 0.424637)
	fmt.Printf("brown: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_DARK_GRAY, 0.686823)
	fmt.Printf("dark-gray: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_LIGHT_GRAY, 0.065637)
	fmt.Printf("light-gray: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_WHITE, 0.156519)
	fmt.Printf("white: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_RED, 0.096970)
	fmt.Printf("red: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_ORANGE, 0.300912)
	fmt.Printf("orange: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_YELLOW, 0.515213)
	fmt.Printf("yellow: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_GREEN, 0.813640)
	fmt.Printf("green: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_BLUE, 0.214264)
	fmt.Printf("blue: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_INDIGO, 0.380657)
	fmt.Printf("indigo: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_PINK, 0.318058)
	fmt.Printf("pink: %f\n", r)
	rt += r
	r = perception.RateImage(im, palettes.PICO8_PEACH, 0.468890)
	fmt.Printf("peach: %f\n", r)
	rt += r
	rt /= 16.0
	fmt.Printf("\nfinal: %f\n", rt)
}
