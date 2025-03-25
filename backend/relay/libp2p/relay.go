package libp2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
	"github.com/dirty-bro-tech/peers-touch-station/relay"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	relayLib "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	ma "github.com/multiformats/go-multiaddr"
)

type Relay struct {
	opts *relay.Options

	initiated  bool
	initDoOnce sync.Once
}

func (r *Relay) Handlers() []server.Handler {
	return []server.Handler{}
}

func (r *Relay) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r *Relay) Name() string {
	return "libp2p-relay"
}

func (r *Relay) Port() int {
	//TODO implement me
	panic("implement me")
}

func (r *Relay) Status() server.ServerStatus {
	//TODO implement me
	panic("implement me")
}

func (r *Relay) Init(ctx context.Context, opts ...option.Option) error {
	for _, o := range opts {
		r.opts.Apply(o)
	}

	if r.opts.KeyFile == "" {
		return fmt.Errorf("no key file provided")
	}

	r.initiated = true
	return nil
}

func (r *Relay) Start(ctx context.Context, opts ...option.Option) error {
	var h host.Host
	var cancel context.CancelFunc

	r.initDoOnce.Do(func() {
		ctx, cancel = context.WithCancel(ctx)

		go func(h host.Host) {
			defer cancel()

			if !r.initiated {
				log.Warn(ctx, "libp2p registry server should be initiated first.")
				return
			}

			for _, opt := range opts {
				r.opts.Apply(opt)
			}

			// Load or generate private key
			privKey, err := loadOrGenerateKey(r.opts.KeyFile)
			if err != nil {
				log.Fatalf(ctx, "Failed to handle private key[%s]: %v", r.opts.KeyFile, err)
			}

			// Create host with custom identity
			h, err = libp2p.New(
				// libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port)),
				libp2p.ListenAddrStrings(r.opts.Addresses.String()...),
				libp2p.Identity(privKey),
				libp2p.EnableRelay(),
				libp2p.EnableNATService(),
			)
			if err != nil {
				log.Fatalf(ctx, "Failed to create host: %v", err)
			}

			// Create and start relay service
			_, err = relayLib.New(h)
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
			go r.discoverPeers(ctx, h, discovery)

			// Print server information
			log.Infof(ctx, "Relay server running with:")
			log.Infof(ctx, " - Peer ID: %s", h.ID())
			for _, addr := range h.Addrs() {
				log.Infof(ctx, " - Address: %s/p2p/%s", addr, h.ID())
			}

			go func(h host.Host) {
				ticker := time.NewTicker(5 * time.Minute)
				defer ticker.Stop()
				defer h.Close()

				for {
					select {
					case <-ticker.C:
						for _, pid := range h.Peerstore().Peers() {
							if r.isRegisteredWithRelay(h, pid) {
								log.Infof(ctx, "Active relay registration: %s", pid)
							}
						}
					case <-ctx.Done():
						log.Info(ctx, "Stopping relay monitoring")
						return
					}
				}
			}(h)

			// Modified server loop
			select {
			case <-ctx.Done():
				log.Info(ctx, "Relay server shutting down")
				err = h.Close()
				if err != nil {
					log.Fatalf(ctx, "Failed to close libp2p host: %v", err)
				}
			}
		}(h)
	})

	return nil
}

func (r *Relay) Options() *relay.Options {
	return r.opts
}

func (r *Relay) List(ctx context.Context, opts ...relay.GetOption) ([]relay.Peer, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Relay) isRegisteredWithRelay(h host.Host, relayID peer.ID) bool {
	for _, conn := range h.Network().Conns() {
		// Check connection direction and protocols
		if conn.RemotePeer() == relayID {
			for _, proto := range conn.RemoteMultiaddr().Protocols() {
				if proto.Code == ma.P_CIRCUIT {
					return true
				}
			}
		}
	}
	return false
}

func NewRelay(opts ...option.Option) server.SubServer {
	rs := &Relay{
		opts: relay.BootstrapOptions(),
	}

	rs.opts.Apply(opts...)
	return rs
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

func (r *Relay) discoverPeers(ctx context.Context, h host.Host, discovery *routing.RoutingDiscovery) {
	for {
		peerChan, err := discovery.FindPeers(ctx, "peers-network")
		if err != nil {
			log.Infof(ctx, "Error finding peers: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		for peer := range peerChan {
			if r.isRegisteredWithRelay(h, peer.ID) {
				log.Infof(ctx, "Already registered with relay: %s", peer.ID)
				continue
			}
			// ... rest of existing peer handling code ...
		}

		time.Sleep(1 * time.Minute)
	}
}

func connectToBootstrap(h host.Host, peers string) {
	for _, addr := range strings.Split(peers, ",") {
		maddr, _ := ma.NewMultiaddr(addr)
		pi, _ := peer.AddrInfoFromP2pAddr(maddr)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := h.Connect(ctx, *pi); err != nil {
			fmt.Printf("Failed to connect to bootstrap %s: %v\n", addr, err)
		} else {
			fmt.Printf("Connected to bootstrap: %s\n", addr)
		}
	}
}
