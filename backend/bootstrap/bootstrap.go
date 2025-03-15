package bootstrap

import (
	"context"

	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/libp2p/go-libp2p/core/peer"
)

type optionsKey struct{}

// Bootstrap defines the interface for bootstrap server functionality
type Bootstrap interface {
	server.SubServer
	ListPeers(ctx context.Context) []peer.ID
	GetAddrInfo(ctx context.Context) peer.AddrInfo
}

// Options holds configuration options for the bootstrap server
type Options struct {
	ListenAddr string
	KeyFile    string
}

// Option defines a function type for setting options
type Option func(*Options)

// WithListenAddr sets the listen address for the bootstrap server
func WithListenAddr(addr string) server.SubServerOption {
	return func(o *server.SubServerOptions) {
		optionWrap(o, func(opts *Options) {
			opts.ListenAddr = addr
		})
	}
}

// WithKeyFile sets the path to the private key file
func WithKeyFile(keyFile string) server.SubServerOption {
	return func(o *server.SubServerOptions) {
		optionWrap(o, func(opts *Options) {
			opts.KeyFile = keyFile
		})
	}
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

func optionWrap(o *server.SubServerOptions, f func(*Options)) {
	if o.Ctx == nil {
		o.Ctx = context.Background()
	}

	var opts *Options
	if o.Ctx.Value(optionsKey{}) == nil {
		opts = &Options{}
		o.Ctx = context.WithValue(o.Ctx, optionsKey{}, opts)
	} else {
		opts = o.Ctx.Value(optionsKey{}).(*Options)
	}

	f(opts)
}
