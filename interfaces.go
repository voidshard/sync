package sync

import (
	"context"
)

// PubSub is a generic interface for a pubsub system.
type PubSub interface {
	// Subscribe subscribes to a list of channels.
	Subscribe(context.Context, ...string) (Subscription, error)

	// Publish publishes a message to a channel.
	Publish(context.Context, string, []byte) error
}

// Subscription is a subscription to a list of channels.
type Subscription interface {
	// Unsubscribe unsubscribes from a list of channels.
	Unsubscribe(context.Context, ...string) error

	// Read returns a channel of messages published to this channel.
	Read() <-chan []byte
}
