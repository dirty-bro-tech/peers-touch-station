package libp2p

import (
	"context"
	"fmt"

	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/plugin/registry/native"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-go/core/store"
	"github.com/dirty-bro-tech/peers-touch-station/bootstrap"
	dbModels "github.com/dirty-bro-tech/peers-touch-station/gen/gorm"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dht_pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"gorm.io/gorm"
)

var (
	_ server.Subserver = (*BootstrapServer)(nil)
)

var (
	connectChan = make(chan connectEvent)
)

type connectEvent struct {
	peerID peer.ID
	msg    *dht_pb.Message
}

type BootstrapServer struct {
	opts *bootstrap.Options

	host host.Host
	dht  *dht.IpfsDHT

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

func NewBootstrapServer(opts ...option.Option) server.Subserver {
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

	// todo: temporary solution
	bs.host, bs.dht = native.GetLibp2pHost()
	if bs.host == nil {
		return fmt.Errorf("host is nil, import native libp2p registry first")
	}

	// After creating host
	if bs.host != nil {
		bs.host.Network().Notify(bs) // Register connection listener
	}

	// init rds
	var err error
	// todo get rds by config
	bs.db, err = store.GetRDS(ctx,
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

func (bs *BootstrapServer) initMonitor(ctx context.Context) {
	// Add periodic routing table logging
	go bs.monitorRoutingTable(ctx, bs.dht)
}
