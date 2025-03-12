package libp2p

import (
	"context"
	"fmt"
	"time"

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

	// Initialize DHT in server mode
	kdht := utils.InitDHT(ctx, h, dht.ModeServer)

	// Print server information
	fmt.Println("Bootstrap server running with:")
	fmt.Printf(" - Peer ID: %s\n", h.ID())
	for _, addr := range h.Addrs() {
		fmt.Printf(" - Address: %s/p2p/%s\n", addr, h.ID())
	}

	return &BootstrapServer{
		host: h,
		dht:  kdht,
	}, nil
}

func (bs *BootstrapServer) Start(ctx context.Context) {
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
func (bs *BootstrapServer) ListPeers() []peer.ID {
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

func (bs *BootstrapServer) GetAddrInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    bs.host.ID(),
		Addrs: bs.host.Addrs(),
	}
}

func (bs *BootstrapServer) Stop() error {
	if err := bs.dht.Close(); err != nil {
		return err
	}
	return bs.host.Close()
}
