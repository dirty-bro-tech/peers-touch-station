package relay

import (
	"context"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
)

var (
	opts *Options
)

type Relay interface {
	server.Subserver

	Options() Options
	List(ctx context.Context, opts ...GetOption) ([]Peer, error)
}

func BootstrapOptions() *Options {
	if opts == nil {
		opts = &Options{
			SubServerOptions: server.NewSubServerOptionsFromRoot(),
		}
	}

	return opts
}
