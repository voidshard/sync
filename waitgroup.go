package sync

import (
	"context"
	"fmt"
	s "sync" // the actual sync package
	"time"

	"github.com/golang/groupcache/lru"
)

var (
	// ErrTimeout is returned when a WaitGroup times out on a Wait() call.
	ErrTimeout = fmt.Errorf("timeout")
)

// WaitGroup is a wrapper around the sync.WaitGroup that uses a pub/sub.
//
// It is used in a similar way, where one caller uses a WaitGroup to wait for
// some number of unique task(s) to complete. The addition of a Pub/Sub pattern
// allows these tasks to be completed by any number of workers over the network.
//
// The idea here is to have at most one WaitGroup of a given name active, we
// don't handle spinning up multiple WaitGroups with the same name at the same
// time.
type WaitGroup struct {
	sub Subscription

	wg    *s.WaitGroup
	cache *lru.Cache
	kill  chan bool
}

// newWaitGroup returns a new WaitGroup with the given subscription and cache size.
func newWaitGroup(sub Subscription, size int) *WaitGroup {
	if size < 1 {
		size = 1
	}

	wg := &WaitGroup{
		sub:   sub,
		wg:    &s.WaitGroup{},
		cache: lru.New(size),
		kill:  make(chan bool),
	}

	go func() {
		defer wg.cache.Clear()
		defer wg.sub.Unsubscribe(context.Background())

		for {
			select {
			case <-wg.kill:
				return
			case msg := <-wg.sub.Read():
				// TODO: Possible if the cache is too small we could count a key twice if it
				// has been evicted. Technically it can be avoided if the size is chosen well.
				data := string(msg[:])

				_, seen := wg.cache.Get(data)
				if seen {
					continue
				}

				wg.wg.Done()
				wg.cache.Add(data, true)
			}
		}
	}()

	return wg
}

// Add adds the given number of tasks to the WaitGroup.
func (wg *WaitGroup) Add(i int) {
	wg.wg.Add(i)
}

// Wait blocks until the WaitGroup has received the expected number of Done() calls or
// the given timeout has elapsed.
//
// If the timeout is reached, ErrTimeout is returned.
func (wg *WaitGroup) Wait(t time.Duration) error {
	timeout := make(chan bool)
	done := make(chan bool)

	go func() {
		// TODO: Possible blocked go routine leakage. Reasoned to be better than
		// the reverse of having possible eternally blocked user routines.
		wg.wg.Wait()
		done <- true
	}()

	go func() {
		time.Sleep(t)
		timeout <- true
	}()

	select {
	case <-done:
		wg.kill <- true
		return nil
	case <-timeout:
		return ErrTimeout
	}
}
