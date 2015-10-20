// The main package for the Game of Life HTTP server.
//
// This is an implementation of Conway's Game of Life. It is a 2D cellular
// automaton.
package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tswast/gameoflife/cell"
	"html/template"
	"image/png"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Save the last two minutes of history.
const cacheSize = 120

// These variables handle the current Game of Life state.
//
// Use `cond` to make sure the handlers can read and the update loop can write
// without conflict.
var (
	curr   int
	fields [cacheSize]*cell.Field
	cond   = sync.NewCond(&sync.Mutex{})
)
var templates = template.Must(template.ParseFiles("index.html"))
var fieldPath = regexp.MustCompile("^/field/([0-9]+).png$")

func main() {
	cond.L.Lock()
	fields[0] = cell.RandomField(128, 128)
	start := &*fields[0]
	cond.L.Unlock()

	tick := time.Tick(time.Second)
	r := cell.Run(start, tick)
	go func() {
		for {
			f := <-r
			cond.L.Lock()
			curr = (curr + 1) % cacheSize
			fields[curr] = f
			cond.Broadcast()
			cond.L.Unlock()
		}
	}()

	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/field/", fieldHandler)
	http.ListenAndServe(":8080", nil)
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %q", err.Error())
		return
	}

	// A read loop is required to handle ping keep-alive messages.
	go func(c *websocket.Conn) {
		for {
			if _, _, err := c.NextReader(); err != nil {
				c.Close()
				break
			}
		}
	}(conn)

	cond.L.Lock()
	prev := curr - 1
	cond.L.Unlock()
	for {
		cond.L.Lock()
		for prev == curr {
			cond.Wait()
		}
		prev = curr
		cond.L.Unlock()

		if err = conn.WriteJSON(struct {
			SeqNum int
		}{prev}); err != nil {
			log.Printf("Error writing sequence number: %q", err.Error())
			return
		}
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	cond.L.Lock()
	c := curr
	cond.L.Unlock()

	err := templates.ExecuteTemplate(w, "index.html", struct{ Curr int }{c})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fieldHandler(w http.ResponseWriter, r *http.Request) {
	i, err := getFieldIndex(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cond.L.Lock()
	f := &*fields[i]
	cond.L.Unlock()

	img := cell.ToImage(f)
	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getFieldIndex(w http.ResponseWriter, r *http.Request) (int, error) {
	m := fieldPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return 0, errors.New("missing field index")
	}
	i, err := strconv.Atoi(m[1])
	if err != nil {
		return i, err
	}
	if i < 0 || i >= cacheSize {
		return i, fmt.Errorf("invalid field index %d", i)
	}
	return i, nil
}
