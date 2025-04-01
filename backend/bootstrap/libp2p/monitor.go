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
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Warnf(ctx, "DHT routing table monitoring context done by %s", ctx.Err())
			return
		case <-ticker.C:
			rt := d.RoutingTable()
			rt.ListPeers()
			logger.Infof(ctx, "DHT routing table status, peerCount=[%d], latency=[%d]", rt.Size(), bs.calculatePeerLatency(rt))
		case connection := <-bs.connectChan:
			logger.Infof(ctx, "DHT request received, type=[%s], peerId=[%s]", connection.msg.Type, connection.peerID[:8]+"...")
			// Print connected peers
			peers := bs.host.Peerstore().Peers()
			logger.Infof(ctx, "Connected peers (%d):\n", len(peers))
			for _, pid := range peers {
				if pid == bs.host.ID() {
					continue
				}

				if bs.host.Network().Connectedness(pid) == 1 { // 1 means Connected
					logger.Infof(ctx, "peer is still connecting, peerId=[%s]", pid)
				} else {
					logger.Infof(ctx, "peer not connected, peerId=[%s]", pid)
				}
			}
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
