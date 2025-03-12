package main

import (
	"context"

	"github.com/dirty-bro-tech/peers-touch-station/registry"
	"github.com/dirty-bro-tech/peers-touch-station/registry/libp2p"
)

func main() {
	reg, err := libp2p.NewRegistry()
	if err != nil {
		return
	}

	err = reg.Init(context.Background(),
		relay.KeyFile("demo.key"),
		relay.Addresses(relay.Addr{
			HeadProtocol:      relay.HeadProtocolIP4,
			Address:           "0.0.0.0",
			TransportProtocol: relay.TransportProtocolTCP,
			Port:              4002,
		}))
	if err != nil {
		return
	}

	err = reg.Start(context.Background())
	if err != nil {
		return
	}
}
