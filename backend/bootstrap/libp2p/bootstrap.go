package libp2p

import (
	"context"
	"fmt"
	"time"

	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/plugin"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-go/core/store"
	"github.com/dirty-bro-tech/peers-touch-station/bootstrap"
	dbModels "github.com/dirty-bro-tech/peers-touch-station/gen/gorm"
	"github.com/dirty-bro-tech/peers-touch-station/utils"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dht_pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"gorm.io/gorm"
)

var (
	_ server.SubServer = (*BootstrapServer)(nil)
)

type connectEvent struct {
	peerID peer.ID
	msg    *dht_pb.Message
}

type BootstrapServer struct {
	opts *bootstrap.Options

	connectChan chan connectEvent
	host        host.Host
	dht         *dht.IpfsDHT

	db *gorm.DB
}

func (bs *BootstrapServer) Options() *server.SubServerOptions {
	return bs.opts.SubServerOptions
}

func (bs *BootstrapServer) Handlers() []server.Handler {
	return []server.Handler{
		bs.ListPeersHandler(),
	}
}

func (bs *BootstrapServer) Name() string {
	return "libp2p-bootstrap"
}

func (bs *BootstrapServer) Port() int {
	// todo implement me
	return 0
}

func (bs *BootstrapServer) Status() server.ServerStatus {
	//TODO implement me
	panic("implement me")
}

func NewBootstrapServer(opts ...option.Option) server.SubServer {
	bs := &BootstrapServer{
		opts: bootstrap.BootstrapOptions(),
	}

	bs.opts.Apply(opts...)
	return bs
}

// Init initializes the bootstrap server
func (bs *BootstrapServer) Init(ctx context.Context, opts ...option.Option) error {
	for _, o := range opts {
		bs.opts.Apply(o)
	}

	// Load or generate private key
	privKey, err := utils.LoadOrGenerateKey(bs.opts.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to handle private key: %w", err)
	}

	// Create host
	bs.host, err = libp2p.New(
		libp2p.ListenAddrStrings(bs.opts.ListenAddr),
		libp2p.Identity(privKey),
	)
	if err != nil {
		return fmt.Errorf("failed to create host: %w", err)
	}

	bs.connectChan = make(chan connectEvent)

	// After creating host
	if bs.host != nil {
		bs.host.Network().Notify(bs) // Register connection listener
	}

	// init rds
	// todo get rds by config
	bs.db, err = store.GetRDS(ctx,
		store.WithQueryStore(plugin.NativePluginName),
		store.WithRDSName("sqlite"),
		store.WithRDSDBName("main"),
	)
	if err != nil {
		return fmt.Errorf("failed to get rds: %w", err)
	}

	// create tables
	if err = bs.db.AutoMigrate(&dbModels.BootstrapNode{}, &dbModels.BootstrapNodesHistory{}); err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}
	return nil
}

func (bs *BootstrapServer) Start(ctx context.Context, opts ...option.Option) error {
	go func() {
		// Initialize DHT in server mode
		bs.dht = bs.initDHT(ctx, bs.host, dht.ModeServer)

		// Print server information
		fmt.Println("Bootstrap server running with:")
		fmt.Printf(" - Peer ID: %s\n", bs.host.ID())
		for _, addr := range bs.host.Addrs() {
			fmt.Printf(" - Address: %s/p2p/%s\n", addr, bs.host.ID())
		}

		// Keep the server running
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Infof(ctx, "bootstrap server stop, ctx done, reason[%s]", ctx.Err())
				return
			}
		}
	}()

	return nil
}

// ListPeers returns a list of all peer IDs currently connected to the bootstrap server
func (bs *BootstrapServer) ListPeers(ctx context.Context) []peer.ID {
	var connectedPeers []peer.ID

	// Get all peers from peerstore
	allPeers := bs.host.Peerstore().Peers()

	// Filter out our own ID and disconnected peers
	for _, pid := range allPeers {
		if pid == bs.host.ID() {
			continue
		}
		if bs.host.Network().Connectedness(pid) == 1 { // 1 means Connected
			connectedPeers = append(connectedPeers, pid)
		}
	}

	return connectedPeers
}

func (bs *BootstrapServer) GetAddrInfo(ctx context.Context) peer.AddrInfo {
	return peer.AddrInfo{
		ID:    bs.host.ID(),
		Addrs: bs.host.Addrs(),
	}
}

func (bs *BootstrapServer) Stop(ctx context.Context) error {
	if err := bs.dht.Close(); err != nil {
		return err
	}
	return bs.host.Close()
}

func (bs *BootstrapServer) initDHT(ctx context.Context, h host.Host, mode dht.ModeOpt) *dht.IpfsDHT {
	kdht, err := dht.New(ctx, h, dht.Mode(mode), dht.OnRequestHook(func(ctx context.Context, s network.Stream, req *dht_pb.Message) {
		bs.connectChan <- connectEvent{
			peerID: s.Conn().RemotePeer(),
			msg:    req,
		}
	}))
	if err != nil {
		log.Fatal(ctx, err)
	}

	// Enhanced bootstrap with retry logic
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		if err = kdht.Bootstrap(ctx); err == nil {
			break
		}
		log.Warnf(ctx, "DHT bootstrap attempt failed, attempt=[%d], error: %s", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		log.Fatal(ctx, "Final DHT bootstrap failed after retries", err)
	}

	// Add periodic routing table logging
	go bs.monitorRoutingTable(ctx, kdht)

	return kdht
}
