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
)

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/field.png", fieldHandler)
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<!DOCTYPE html>\n")
	fmt.Fprint(w, "<html>Hello!<br><img src='field.png'>\n")
}

func fieldHandler(w http.ResponseWriter, r *http.Request) {
	f := cell.RandomField(128, 128)
	img := cell.ToImage(f)
	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
