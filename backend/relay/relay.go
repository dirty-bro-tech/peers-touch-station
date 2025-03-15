package relay

import (
	"context"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
)

type Relay interface {
	server.SubServer

	Options() Options
	List(ctx context.Context, opts ...GetOption) ([]Peer, error)
}
