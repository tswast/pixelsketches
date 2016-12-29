// Copyright 2016 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/tswast/pixelsketches/village/artist"
)

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

	if err := artist.Main(seed, false /* don't write timelapse */); err != nil {
		log.Fatal(err)
	}
}
