// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/tswast/pixelsketches/palettes"
	"github.com/tswast/pixelsketches/village/artist"
	"github.com/tswast/pixelsketches/village/perception"
	"github.com/tswast/pixelsketches/village/strategy"
)

func main() {
	var seed int
	var maxIter int
	var tl bool
	var debug bool
	var inp string
	var p string
	var st string
	flag.IntVar(&seed, "seed", 19700101, "Seed used for random number generator.")
	flag.IntVar(&maxIter, "max-iter", 1000000, "Maximum number of iterations.")
	flag.BoolVar(&tl, "timelapse", false, "Write timelapse to out/ directory.")
	flag.BoolVar(&debug, "debug", false, "Write current frame and action.")
	flag.StringVar(&inp, "in", "", "Path to input file to start edit with.")
	flag.StringVar(&p, "out", "", "Path to output file.")
	flag.StringVar(&st, "strategy", "random", "Strategy to use: random|ideal")
	flag.Parse()
	if p == "" {
		log.Fatal("Value for -out is missing.")
	}

	var s strategy.Strategizer
	if st == "random" {
		s = &strategy.RandomWalk{}
	} else if st == "ideal" {
		s = &strategy.Ideal{Rating: perception.RateWholeImage}
	} else if st == "dictator" {
		s = &strategy.Ideal{Rating: perception.RateBlack}
	} else if st == "plurality" {
		s = &strategy.Plurality{[]*strategy.Ideal{
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealBlack, palettes.PICO8_BLACK)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealDarkBlue, palettes.PICO8_DARK_BLUE)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealDarkPurple, palettes.PICO8_DARK_PURPLE)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealDarkGreen, palettes.PICO8_DARK_GREEN)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealBrown, palettes.PICO8_BROWN)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealDarkGray, palettes.PICO8_DARK_GRAY)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealLightGray, palettes.PICO8_LIGHT_GRAY)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealWhite, palettes.PICO8_WHITE)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealRed, palettes.PICO8_RED)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealOrange, palettes.PICO8_ORANGE)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealYellow, palettes.PICO8_YELLOW)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealGreen, palettes.PICO8_GREEN)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealBlue, palettes.PICO8_BLUE)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealIndigo, palettes.PICO8_INDIGO)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealPink, palettes.PICO8_PINK)},
			&strategy.Ideal{Rating: perception.NewRating(perception.IdealPeach, palettes.PICO8_PEACH)},
		}}
	} else {
		log.Fatal("Unexpected value for strategy.")
	}

	if err := artist.Main(inp, p, int64(seed), debug, tl, maxIter, s); err != nil {
		log.Fatal(err)
	}
}
