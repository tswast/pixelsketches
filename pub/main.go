// The main package for the Game of Life ZeroMQ Pub/Sub server.
//
// This is an implementation of Conway's Game of Life. It is a 2D cellular
// automaton.
package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/tswast/gameoflife/cell"
	"github.com/zeromq/goczmq"
	"log"
	"time"
)

func main() {
	pub := goczmq.NewPubChanneler("tcp://*:5000")
	defer pub.Destroy()

	tick := time.Tick(time.Second)
	r := cell.Run(cell.RandomField(128, 128), tick)
	for f := range r {
		data, err := proto.Marshal(cell.ToProto(f))
		if err != nil {
			log.Fatal("marshaling error for Field %#v: ", *f, err)
		}
		pub.SendChan <- [][]byte{data}
	}
}
