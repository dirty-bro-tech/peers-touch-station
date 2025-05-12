package main

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dirty-bro-tech/peers-touch-go"
	"github.com/dirty-bro-tech/peers-touch-go/core/debug/actuator"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-go/core/service"
	local "github.com/dirty-bro-tech/peers-touch-station/bootstrap/libp2p"
	"github.com/dirty-bro-tech/peers-touch-station/relay"
	"github.com/dirty-bro-tech/peers-touch-station/relay/libp2p"

	// default plugins
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/native"
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/registry/native"
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/store/native"
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/store/rds/postgres"
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/store/rds/sqlite"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	local.Init()

	p := peers.NewPeer()
	err := p.Init(
		ctx,
		service.WithPrivateKey("private.pem"),
		service.Name("peers-touch-station"),
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
		server.WithSubServer("debug", actuator.NewDebugSubServer, actuator.WithDebugServerPath("")),

		/*		server.WithSubServer("bootstrapServer",
				bootstrapP2p.NewBootstrapServer,
				bootstrap.WithListenAddr("/ip4/0.0.0.0/tcp/4001")),*/
		server.WithSubServer("relyServer", libp2p.NewRelay,
			relay.KeyFile("libp2pIdentity.key"),
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
