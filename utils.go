package sync

import (
	"fmt"
)

// NewPubSub returns a new pubsub instance for the given options.
func NewPubSub(opts *Options) (PubSub, error) {
	switch opts.Engine {
	case EngineRedis:
		return newRedis(opts)
	}
	return nil, fmt.Errorf("unknown engine %s", opts.Engine)
}
