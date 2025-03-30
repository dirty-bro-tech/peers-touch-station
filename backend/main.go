package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dirty-bro-tech/peers-touch-go"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-station/bootstrap"
	bootstrapP2p "github.com/dirty-bro-tech/peers-touch-station/bootstrap/libp2p"
	"github.com/dirty-bro-tech/peers-touch-station/relay"
	"github.com/dirty-bro-tech/peers-touch-station/relay/libp2p"
	"net/http"

	// default plugins
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/native"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	/*// Start bootstrap server
	bootstrapServer := bootstrapP2p.NewBootstrapServer(bootstrap.WithListenAddr("/ip4/0.0.0.0/tcp/4001"), bootstrap.WithKeyFile("demo.key"))
	// Add peer printer ticker
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// wait the bootstrap server completed init
		time.Sleep(200 * time.Second)

		for {
			select {
			case <-ticker.C:
				listPeers := bootstrapServer.ListPeers(ctx)
				fmt.Printf("Connected listPeers (%d):\n", len(listPeers))
				for _, peer := range listPeers {
					fmt.Println(" -", peer)
				}
			case <-ctx.Done():
				return
			}
		}
	}()*/

	p := peers.NewPeer()
	err := p.Init(
		ctx,
		peers.WithName("hello-world"),
		server.WithHandlers(
			server.NewHandler("hello-world", "/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello world, from native handler"))
			})),
			server.NewHandler("hello-world-hertz", "/hello-hz",
				func(c context.Context, ctx *app.RequestContext) {
					ctx.String(http.StatusOK, "hello world, from hertz handler")
				},
			),
		),
		server.WithSubServer("bootstrapServer",
			bootstrapP2p.NewBootstrapServer,
			bootstrap.WithListenAddr("/ip4/0.0.0.0/tcp/4001"),
			bootstrap.WithKeyFile("demo.key")),
		server.WithSubServer("relyServer", libp2p.NewRelay,
			relay.KeyFile("demo.key"),
			relay.Addresses(relay.Addr{
				HeadProtocol:      relay.HeadProtocolIP4,
				Address:           "0.0.0.0",
				TransportProtocol: relay.TransportProtocolTCP,
				Port:              4002,
			},
			),
		),
	)
	if err != nil {
		return
		panic(err)
	}

	err = p.Start()
	if err != nil {
		panic(err)
	}
}
