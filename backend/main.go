package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	bootstrapP2p "github.com/dirty-bro-tech/peers-touch-station/bootstrap/libp2p"
	"github.com/dirty-bro-tech/peers-touch-station/relay"
	"github.com/dirty-bro-tech/peers-touch-station/relay/libp2p"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start bootstrap server
	bootstrapServer, err := bootstrapP2p.NewBootstrapServer(ctx, "/ip4/0.0.0.0/tcp/4001", "demo.key")
	if err != nil {
		panic(err)
	}
	go bootstrapServer.Start(ctx)

	// Add peer printer ticker
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				peers := bootstrapServer.ListPeers()
				fmt.Printf("Connected peers (%d):\n", len(peers))
				for _, peer := range peers {
					fmt.Println(" -", peer)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start relay server
	reg, err := libp2p.NewRegistry()
	if err != nil {
		panic(err)
	}

	err = reg.Init(ctx,
		relay.KeyFile("demo.key"),
		relay.Addresses(relay.Addr{
			HeadProtocol:      relay.HeadProtocolIP4,
			Address:           "0.0.0.0",
			TransportProtocol: relay.TransportProtocolTCP,
			Port:              4002,
		}))
	if err != nil {
		panic(err)
	}

	go func() {
		if err := reg.Start(ctx); err != nil {
			panic(err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	cancel()

	// Graceful shutdown
	/*if err :=reg.Stop(); err != nil {
		panic(err)
	}*/
	if err := bootstrapServer.Stop(); err != nil {
		panic(err)
	}
}
