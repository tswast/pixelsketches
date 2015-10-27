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

	update := make(chan *cell.UpdateRequest)
	go startInputServer(update)

	tick := time.Tick(time.Second)
	r := cell.Run(cell.RandomField(128, 128), tick, update)
	seq := int32(0)
	for f := range r {
		pf := cell.ToProto(f)
		pf.Seq = seq

		data, err := proto.Marshal(pf)
		if err != nil {
			log.Fatal("Marshaling error: %q\n\tfor Field %#v\n", err.Error(), *f)
		}
		pub.SendChan <- [][]byte{data}
		seq++
	}
}

func startInputServer(update chan<- *cell.UpdateRequest) {
	rep := goczmq.NewRepChanneler("tcp://*:5001")
	defer rep.Destroy()

	for msg := range rep.RecvChan {
		if len(msg) != 1 {
			log.Printf(
				"Message had unexpected number of frames: %d, want 1.\n",
				len(msg))
			continue
		}
		data := msg[0]

		r := &cell.UpdateRequest{}
		if err := proto.Unmarshal(data, r); err != nil {
			log.Printf(
				"Error unmarshalling message:\n\t%q,\n\t%#v\n",
				err.Error(),
				msg)
			continue
		}
		update <- r

		reply, err := proto.Marshal(&cell.UpdateResponse{})
		if err != nil {
			log.Fatal("Marshaling error: %q\n\tfor update response\n", err.Error())
		}
		rep.SendChan <- [][]byte{reply}
	}
}
