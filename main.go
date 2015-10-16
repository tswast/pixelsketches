// The main package for the Game of Life HTTP server.
//
// This is an implementation of Conway's Game of Life. It is a 2D cellular
// automaton.
package main

import (
	"fmt"
	"github.com/tswast/gameoflife/cell"
	"image/png"
	"net/http"
	"sync"
	"time"
)

var (
	curr  *cell.Field
	mutex *sync.Mutex = &sync.Mutex{}
)

func main() {
	mutex.Lock()
	curr = cell.RandomField(128, 128)
	start := &*curr
	mutex.Unlock()

	tick := time.Tick(time.Second)
	r := cell.Run(start, tick)
	go func() {
		for {
			c := <-r
			mutex.Lock()
			curr = c
			mutex.Unlock()
		}
	}()

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/field.png", fieldHandler)
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<!DOCTYPE html>\n<html>")
	fmt.Fprint(w, "Hello!<br>")
	fmt.Fprint(
		w, "<img style='width:80%; image-rendering: pixelated' src='field.png'>\n")
}

func fieldHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	f := &*curr
	mutex.Unlock()

	img := cell.ToImage(f)
	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
