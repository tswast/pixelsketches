// The main package for the Game of Life HTTP server.
//
// This is an implementation of Conway's Game of Life. It is a 2D cellular
// automaton.
package main

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/tswast/gameoflife/cell"
	"github.com/zeromq/goczmq"
	"html/template"
	"image/png"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
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

var (
	req  *goczmq.Sock
	reqL = &sync.Mutex{}
)

func main() {
	pub := os.Getenv("PUB_PORT_5000_TCP")
	if pub == "" {
		log.Fatal("Missing PUB_PORT_5000_TCP environment variable")
	}
	rep := os.Getenv("PUB_PORT_5001_TCP")
	if rep == "" {
		log.Fatal("Missing PUB_PORT_5001_TCP environment variable")
	}
	reqL.Lock()
	var err error
	req, err = goczmq.NewReq(rep)
	reqL.Unlock()
	if err != nil {
		log.Fatal("Error connecting to REP %q: %q\n", rep, err.Error())
	}

	fields[0] = cell.RandomField(128, 128)
	go func() {
		sub := goczmq.NewSubChanneler(pub, "" /* all messages, no topic filtering */)
		defer sub.Destroy()

		for msg := range sub.RecvChan {
			if len(msg) != 1 {
				log.Printf(
					"Message had unexpected number of frames: %d, want 1.\n",
					len(msg))
				continue
			}
			data := msg[0]

			pf := &cell.FieldProto{}
			if err := proto.Unmarshal(data, pf); err != nil {
				log.Printf(
					"Error unmarshalling message:\n\t%q,\n\t%#v\n",
					err.Error(),
					msg)
				continue
			}
			seq := int(pf.Seq) % cacheSize
			if seq < 0 {
				seq = seq + cacheSize
			}

			f, err := cell.FromProto(pf)
			if err != nil {
				log.Printf("Got invalid Field:\n\t%q,\n\t%#v\n", err.Error(), f)
				continue
			}

			cond.L.Lock()
			curr = seq
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

	// Listen for update requests over the websocket.
	go func(c *websocket.Conn) {
		for _, r, err := c.NextReader(); err == nil; _, r, err = c.NextReader() {
			reqL.Lock()
			if req == nil {
				log.Printf("Not connected to REQ/REP on pub\n")
				continue
			}
			reqL.Unlock()

			msg := &cell.UpdateRequest{}
			if err := jsonpb.Unmarshal(r, msg); err != nil {
				log.Printf("Error unmarshalling UpdateRequest: %q\n", err.Error())
				continue
			}

			data, err := proto.Marshal(msg)
			if err != nil {
				log.Fatal("Marshaling error: %q\n\tfor %#v\n", err.Error(), *msg)
			}

			reqL.Lock()
			if err = req.SendMessage([][]byte{data}); err != nil {
				log.Printf("Error sending request: %q\n", err.Error())
				continue
			}
			if _, err = req.RecvMessage(); err != nil {
				log.Printf("Error receiving reply: %q\n", err.Error())
				continue
			}
			reqL.Unlock()
		}
		c.Close()
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
	if fields[i] == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		cond.L.Unlock()
		return
	}
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
