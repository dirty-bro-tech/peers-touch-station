package libp2p

import (
	"context"
	"time"

	"github.com/dirty-bro-tech/peers-touch-go/core/logger"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	kb "github.com/libp2p/go-libp2p-kbucket"
)

// New method to monitor DHT health
func (bs *BootstrapServer) monitorRoutingTable(ctx context.Context, d *dht.IpfsDHT) {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rt := d.RoutingTable()
			logger.Infof(ctx, "DHT routing table status, peerCount=[%d], latency=[%d]", rt.Size(), bs.calculatePeerLatency(rt))
		case <-ctx.Done():
			logger.Warnf(ctx, "DHT routing table monitoring context done by %s", ctx.Err())
			return
		}
	}
}

func (bs *BootstrapServer) calculatePeerLatency(rt *kb.RoutingTable) time.Duration {
	var total time.Duration
	count := 0

	// Iterate through all peers in routing table
	for _, pid := range rt.ListPeers() {
		// Get latency from peerstore's EWMA (Exponentially Weighted Moving Average)
		total += bs.host.Peerstore().LatencyEWMA(pid)
		count++
	}

	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}
