package main

import (
	"context"
	"fmt"
	"math/rand"
	s "sync"
	"time"

	"github.com/voidshard/sync"
)

const (
	group = "mygroup"
)

func main() {
	// FanOut is some "worker" or async process that sends a Done() message(s)
	// to our listening WaitGroup as items complete.
	ps, err := sync.NewPubSub(sync.NewOptions())
	if err != nil {
		panic(err)
	}

	sink := sync.New(ps)

	// More complex than required. This is to test that we handle multiple
	// messages arriving for the same ID.
	wg := &s.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		work(sink, 0, 7)
	}()

	go func() {
		defer wg.Done()
		work(sink, 4, 10) // deliberate overlap
	}()

	wg.Wait()

}

func work(sink *sync.Sync, start, end int) error {
	for i := start; i < end; i++ {
		time.Sleep(time.Duration(rand.Intn(50)+20) * time.Millisecond)

		id := fmt.Sprintf("id-%d", i) // unique reference

		err := sink.Done(context.TODO(), group, id)
		if err != nil {
			return err
		}

		fmt.Println("- sent", id)
	}
	return nil
}
