package libp2p

import (
	"context"
	"fmt"
	"time"

	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-station/utils"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type BootstrapServer struct {
	host host.Host
	dht  *dht.IpfsDHT
}

func (bs *BootstrapServer) Name() string {
	return "libp2p-bootstrap"
}

func (bs *BootstrapServer) Port() int {
	return bs.host
}

func (bs *BootstrapServer) Status() server.ServerStatus {
	//TODO implement me
	panic("implement me")
}

func NewBootstrapServer(ctx context.Context, listenAddr string, keyFile string) (*BootstrapServer, error) {
	// Load or generate private key
	privKey, err := utils.LoadOrGenerateKey(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to handle private key: %w", err)
	}

	// Create host
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.Identity(privKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	return &BootstrapServer{
		host: h,
	}, nil
}

// Init initializes the bootstrap server
func (bs *BootstrapServer) Init(ctx context.Context, opts ...server.SubServerOption) error {
	// todo
	return nil
}

func (bs *BootstrapServer) Start(ctx context.Context, opts ...server.SubServerOption) {
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
			return
		case <-ticker.C:
			// Print connected peers
			peers := bs.host.Peerstore().Peers()
			fmt.Printf("Connected peers (%d):\n", len(peers))
			for _, pid := range peers {
				if pid == bs.host.ID() {
					continue
				}
				fmt.Printf(" - %s\n", pid)
			}
		}
	}
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
	kdht, err := dht.New(ctx, h, dht.Mode(mode))
	if err != nil {
		log.Fatal(ctx, err)
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		log.Fatal(ctx, err)
	}
	return kdht
}
