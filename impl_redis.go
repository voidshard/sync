package sync

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Redis is a PubSub implementation using Redis.
type Redis struct {
	rdb redis.UniversalClient
}

// RedisSub is a subscription to a list of channels.
type RedisSub struct {
	ps       *redis.PubSub
	read     chan []byte
	channels []string
}

// New returns a new PubSubRedis instance.
func newRedis(opts *Options) (*Redis, error) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    opts.Addresses,
		Username: opts.Username,
		Password: opts.Password,
	})
	return &Redis{rdb: rdb}, nil
}

// Publish publishes a message to a channel.
func (ps *Redis) Publish(ctx context.Context, channel string, msg []byte) error {
	return ps.rdb.Publish(ctx, channel, msg).Err()
}

// Subscribe subscribes to a list of channels.
func (ps *Redis) Subscribe(ctx context.Context, channels ...string) (Subscription, error) {
	if len(channels) == 0 {
		// technically redis does allow this
		return nil, fmt.Errorf("no channels provided")
	}
	pubsub := ps.rdb.Subscribe(ctx, channels...)

	rs := &RedisSub{
		ps:       pubsub,
		channels: channels,
		read:     make(chan []byte),
	}
	go func() {
		for msg := range rs.ps.Channel() {
			rs.read <- []byte(msg.Payload)
		}
	}()

	return rs, nil
}

// Unsubscribe unsubscribes from a list of channels.
func (s *RedisSub) Unsubscribe(ctx context.Context, channels ...string) error {
	if len(channels) == 0 {
		channels = s.channels
	}
	return s.ps.Unsubscribe(ctx, channels...)
}

// Read returns a channel of messages published to this channel.
func (s *RedisSub) Read() <-chan []byte {
	return s.read
}
