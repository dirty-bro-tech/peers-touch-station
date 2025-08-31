package main

import (
	"context"

	"github.com/dirty-bro-tech/peers-touch-go"
	"github.com/dirty-bro-tech/peers-touch-go/core/debug/actuator"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-go/core/service"
	"github.com/dirty-bro-tech/peers-touch-station/subserver/family"

	// default plugins
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/native"
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/native/registry"
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/store/rds/postgres"
	_ "github.com/dirty-bro-tech/peers-touch-go/core/plugin/store/rds/sqlite"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := peers.NewPeer()
	err := p.Init(
		ctx,
		service.WithPrivateKey("private.pem"),
		service.Name("peers-touch-station"),
		server.WithSubServer("debug", actuator.NewDebugSubServer, actuator.WithDebugServerPath("")),
		server.WithSubServer("family", family.NewPhotoSaveSubServer, family.WithPhotoSaveDir("photos-directory")),

		/*		server.WithSubServer("bootstrapServer",
						bootstrapP2p.NewBootstrapServer,
						bootstrap.WithListenAddr("/ip4/0.0.0.0/tcp/4001")),
				server.WithSubServer("relyServer", libp2p.NewRelay,
					relay_.KeyFile("libp2pIdentity.key"),
					relay_.Addresses(relay_.Addr{
						HeadProtocol:      relay_.HeadProtocolIP4,
						Address:           "0.0.0.0",
						TransportProtocol: relay_.TransportProtocolTCP,
						Port:              4002,
					},
					),
				),*/
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
