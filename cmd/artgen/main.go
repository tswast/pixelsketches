// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/tswast/pixelsketches/village/artist"
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
		s = &strategy.Ideal{}
	} else {
		log.Fatal("Unexpected value for strategy.")
	}

	if err := artist.Main(inp, p, int64(seed), debug, tl, maxIter, s); err != nil {
		log.Fatal(err)
	}
}
