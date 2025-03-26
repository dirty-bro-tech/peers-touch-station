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
			logger.Info(ctx, "DHT routing table status",
				"peerCount", rt.Size(),
				"networkSize", bs.estimateNetworkSize(rt),
				"latency", bs.calculatePeerLatency(rt))
		case <-ctx.Done():
			return
		}
	}
}

func (bs *BootstrapServer) estimateNetworkSize(rt *kb.RoutingTable) int {
	// Use Kademlia's k-bucket structure to estimate network size
	// Each successive bucket represents a doubling of the address space
	buckets := rt.Buckets()
	if len(buckets) == 0 {
		return 0
	}

	// Last bucket depth indicates network size order of magnitude
	lastBucketDepth := len(buckets) - 1
	return (1 << uint(lastBucketDepth)) * rt.BucketSize()
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
