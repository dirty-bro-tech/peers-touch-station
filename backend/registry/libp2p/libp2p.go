package libp2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-station/registry"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

type Registry struct {
	opts *registry.Options

	initiated  bool
	initDoOnce sync.Once
}

func (r *Registry) Init(ctx context.Context, opts ...registry.Option) error {
	if r.opts == nil {
		r.opts = &registry.Options{}
	}

	for _, o := range opts {
		o(r.opts)
	}

	r.initiated = true
	return nil
}

func (r *Registry) Start(ctx context.Context, opts ...registry.Option) error {
	r.initDoOnce.Do(func() {
		if !r.initiated {
			log.Warn(ctx, "libp2p registry server should be initiated first.")
			return
		}

		for _, opt := range opts {
			opt(r.opts)
		}

		// Load or generate private key
		privKey, err := loadOrGenerateKey(r.opts.KeyFile)
		if err != nil {
			log.Fatalf(ctx, "Failed to handle private key: %v", err)
		}

		// Create host with custom identity
		h, err := libp2p.New(
			// libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port)),
			libp2p.ListenAddrStrings(r.opts.Addresses.String()...),
			libp2p.Identity(privKey),
		)
		if err != nil {
			log.Fatalf(ctx, "Failed to create host: %v", err)
		}

		// Create and start relay service
		_, err = relay.New(h)
		if err != nil {
			log.Fatalf(ctx, "Failed to start relay service: %v", err)
		}

		// Initialize DHT in server mode
		kdht := initDHT(context.Background(), h, dht.ModeServer)
		// Create routing discovery
		discovery := routing.NewRoutingDiscovery(kdht)
		// Advertise our presence
		util.Advertise(context.Background(), discovery, "peers-network")
		// Start peer discovery
		go discoverPeers(ctx, h, discovery)

		// Print server information
		fmt.Println("Relay and bootstrap server running with:")
		fmt.Printf(" - Peer ID: %s\n", h.ID())
		for _, addr := range h.Addrs() {
			fmt.Printf(" - Address: %s/p2p/%s\n", addr, h.ID())
		}

		// Keep the server running
		select {}
	})

	return nil
}

func (r *Registry) Options() *registry.Options {
	return r.opts
}

func (r *Registry) List(ctx context.Context, opts ...registry.GetOption) ([]registry.Peer, error) {
	//TODO implement me
	panic("implement me")
}

func NewRegistry(opts ...registry.Option) (*Registry, error) {
	r := &Registry{}
	for _, opt := range opts {
		opt(r.opts)
	}
	return r, nil
}

func loadOrGenerateKey(keyFile string) (crypto.PrivKey, error) {
	// Try to load existing key
	if data, err := os.ReadFile(keyFile); err == nil {
		return crypto.UnmarshalPrivateKey(data)
	}

	// Generate new key
	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Save the key
	data, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(keyFile, data, 0600); err != nil {
		return nil, err
	}

	return privKey, nil
}

func initDHT(ctx context.Context, h host.Host, mode dht.ModeOpt) *dht.IpfsDHT {
	kdht, err := dht.New(ctx, h, dht.Mode(mode))
	if err != nil {
		log.Fatal(ctx, err)
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		log.Fatal(ctx, err)
	}
	return kdht
}

func discoverPeers(ctx context.Context, h host.Host, discovery *routing.RoutingDiscovery) {
	for {
		peerChan, err := discovery.FindPeers(context.Background(), "peers-network")
		if err != nil {
			log.Infof(ctx, "Error finding peers: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		for peer := range peerChan {
			if peer.ID == h.ID() || len(peer.Addrs) == 0 {
				continue
			}
			fmt.Printf("Discovered peer: %s\n", peer.ID)
		}

		time.Sleep(1 * time.Minute)
	}
}
