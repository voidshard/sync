package sync

import (
	"os"
	"strings"
)

// Engine is a type of PubSub engine.
type Engine string

const (
	// EngineRedis is the redis engine.
	EngineRedis Engine = "redis"
)

const (
	// EnvPubSubAddresses is the environment variable name for the addresses.
	// Addresses are separated by colons (ie the `:` symbol).
	EnvPubSubAddresses = "PUBSUB_ADDRESSES"

	// EnvPubSubUsername is the environment variable name for the username
	EnvPubSubUsername = "PUBSUB_USERNAME"

	// EnvPubSubPassword is the environment variable name for the password
	EnvPubSubPassword = "PUBSUB_PASSWORD"

	// EnvPubSubEngine is the environment variable name for the engine
	EnvPubSubEngine = "PUBSUB_ENGINE"
)

var (
	// defaultAddress is a map of default addresses for each engine.
	defaultAddress = map[Engine]string{
		EngineRedis: "localhost:6379",
	}
)

// Options holds generic options for connecting to something
type Options struct {
	Engine    Engine
	Addresses []string
	Username  string
	Password  string
}

// NewOptions returns a new Options struct.
// We read default variables from the environment.
// - If not given we default to redis.
// - If no addresses are given we default to localhost.
func NewOptions() *Options {
	eng := Engine(os.Getenv(EnvPubSubEngine))
	if eng == "" {
		eng = EngineRedis
	}

	addresses := []string{}
	for _, address := range strings.Split(os.Getenv(EnvPubSubAddresses), ":") {
		if address != "" {
			addresses = append(addresses, address)
		}
	}
	if len(addresses) == 0 {
		defaultAddr, ok := defaultAddress[eng]
		if ok {
			addresses = append(addresses, defaultAddr)
		}
	}

	return &Options{
		Engine:    eng,
		Addresses: addresses,
		Username:  os.Getenv(EnvPubSubUsername),
		Password:  os.Getenv(EnvPubSubPassword),
	}
}
