package utils

import (
	"context"

	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

func InitDHT(ctx context.Context, h host.Host, mode dht.ModeOpt) *dht.IpfsDHT {
	kdht, err := dht.New(ctx, h, dht.Mode(mode))
	if err != nil {
		log.Fatal(ctx, err)
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		log.Fatal(ctx, err)
	}
	return kdht
}
