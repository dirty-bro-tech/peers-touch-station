package bootstrap

import (
	"context"

	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	opts *Options
)

type optionsKey struct{}

var OptionWrapper = option.NewWrapper[Options](optionsKey{}, func(options *option.Options) *Options {
	return BootstrapOptions()
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
}

// Option defines a function type for setting options
type Option func(*Options)

// WithListenAddr sets the listen address for the bootstrap server
func WithListenAddr(addr string) option.Option {
	return OptionWrapper.Wrap(func(opts *Options) {
		opts.ListenAddr = addr
	})
}

func BootstrapOptions() *Options {
	if opts == nil {
		opts = &Options{
			SubServerOptions: server.NewSubServerOptionsFromRoot(),
		}
	}

	return opts
}
