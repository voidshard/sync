package main

import (
	"context"
	"fmt"
	"time"

	"github.com/voidshard/sync"
)

const (
	group = "mygroup"
)

func main() {
	// FanIn kicks off a WaitGroup behind a Pub/Sub and waits for it to complete.

	ps, err := sync.NewPubSub(sync.NewOptions())
	if err != nil {
		panic(err)
	}
	sink := sync.New(ps)

	wg, err := sink.WaitGroup(context.Background(), group, 512)
	if err != nil {
		panic(err)
	}
	wg.Add(10)

	fmt.Println("waiting ..")
	err = wg.Wait(time.Second * 60)
	if err != nil {
		panic(err)
	}
}
