package bootstrap

import (
	"context"

	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/libp2p/go-libp2p/core/peer"
)

type optionsKey struct{}

var wrapper = option.NewWrapper[Options](optionsKey{}, func(options *option.Options) *Options {
	return &Options{
		SubServerOptions: server.NewSubServerOptionsFromRoot(),
	}
})

// Bootstrap defines the interface for bootstrap server functionality
type Bootstrap interface {
	server.SubServer
	ListPeers(ctx context.Context) []peer.ID
	GetAddrInfo(ctx context.Context) peer.AddrInfo
}

// Options holds configuration options for the bootstrap server
type Options struct {
	*server.SubServerOptions

	ListenAddr string
	KeyFile    string
}

// Option defines a function type for setting options
type Option func(*Options)

// WithListenAddr sets the listen address for the bootstrap server
func WithListenAddr(addr string) option.Option {
	return wrapper.Wrap(func(opts *Options) {
		opts.ListenAddr = addr
	})
}

// WithKeyFile sets the path to the private key file
func WithKeyFile(keyFile string) option.Option {
	return wrapper.Wrap(func(opts *Options) {
		opts.KeyFile = keyFile
	})
}

// NewOptions creates a new Options instance with default values
func NewOptions(opts ...Option) *Options {
	options := &Options{
		ListenAddr: "/ip4/0.0.0.0/tcp/4001", // Default listen address
		KeyFile:    "bootstrap.key",         // Default key file
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
