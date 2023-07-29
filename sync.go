package sync

import (
	"context"
)

// Sync a simple WaitGroup-esque wrapper around a PubSub
type Sync struct {
	ps PubSub
}

// New returns a new Sync instance with the given PubSub
func New(ps PubSub) *Sync {
	return &Sync{ps: ps}
}

// WaitGroup returns a new WaitGroup with the given name and cache size.
func (s *Sync) WaitGroup(ctx context.Context, name string, size int) (*WaitGroup, error) {
	sub, err := s.ps.Subscribe(ctx, name)
	return newWaitGroup(sub, size), err
}

// Done notifies the given WaitGroup that some id has completed.
//
// The ID here is any unique reference; it is used on fan-in to de-duplicate messages
// to help avoid counting Done() messages twice for the same thing (on retries / redelivery etc).
func (s *Sync) Done(ctx context.Context, name, id string) error {
	return s.ps.Publish(ctx, name, []byte(id))
}
