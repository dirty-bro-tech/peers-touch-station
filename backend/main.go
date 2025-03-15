package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dirty-bro-tech/peers-touch-go"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	bootstrapP2p "github.com/dirty-bro-tech/peers-touch-station/bootstrap/libp2p"
	"github.com/dirty-bro-tech/peers-touch-station/relay"
	"github.com/dirty-bro-tech/peers-touch-station/relay/libp2p"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	// Start bootstrap server
	bootstrapServer, err := bootstrapP2p.NewBootstrapServer(ctx, "/ip4/0.0.0.0/tcp/4001", "demo.key")
	if err != nil {
		panic(err)
	}

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

	p := peers.NewPeer()
	err = p.Init(
		ctx,
		peers.WithName("hello-world"),
		peers.WithAppendHandlers(
			server.NewHandler("hello-world", "/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello world, from native handler"))
			})),
			server.NewHandler("hello-world-hertz", "/hello-hz",
				func(c context.Context, ctx *app.RequestContext) {
					ctx.String(http.StatusOK, "hello world, from hertz handler")
				},
			),
		),
		peers.WithSubServer(bootstrapServer),
		peers.WithSubServer(reg),
	)
	if err != nil {
		panic(err)
	}

	// Add peer printer ticker
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				peers := bootstrapServer.ListPeers(ctx)
				fmt.Printf("Connected peers (%d):\n", len(peers))
				for _, peer := range peers {
					fmt.Println(" -", peer)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	err = p.Start(ctx)
	if err != nil {
		panic(err)
	}
}
