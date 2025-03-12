package relay

import (
	"context"
)

type Registry interface {
	Init(ctx context.Context, opts ...Option) error
	Start(ctx context.Context, opts ...Option) error
	Options() Options
	List(ctx context.Context, opts ...GetOption) ([]Peer, error)
}
