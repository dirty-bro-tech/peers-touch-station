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
		// Use the new router pattern for family endpoints
		server.WithRouters(family.NewFamilyRouter()),
		// Initialize family options
		family.WithPhotoSaveDir("photos-directory"),
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
