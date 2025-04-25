package libp2p

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dirty-bro-tech/peers-touch-go/core/logger"
	gormGen "github.com/dirty-bro-tech/peers-touch-station/gen/gorm"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"gorm.io/gorm"
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
			ps := rt.ListPeers()
			logger.Infof(ctx, "DHT routing table status, peerCount=[%d]-[%d], latency=[%d]", rt.Size(), len(ps), bs.calculatePeerLatency(rt))
		case connection := <-connectChan:
			// todo new context
			logger.Infof(ctx, "DHT request received, type=[%s], peerId=[%s]", connection.msg.Type, connection.peerID[:8]+"...")

			// New code: Insert into bootstrap_nodes
			addrs := bs.host.Peerstore().Addrs(connection.peerID)
			addrsStr, _ := json.Marshal(addrs)

			node := gormGen.BootstrapNode{
				PeerID:                   connection.peerID.String(),
				MultiAddresses:           string(addrsStr),
				ProtocolVersion:          "ipfs/0.1.0", // Default or extract from metadata
				LastSuccessfulConnection: time.Now(),
			}

			err := bs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				// delete old records if exists
				if err := tx.Where("peer_id = ?", node.PeerID).Delete(&gormGen.BootstrapNode{}).Error; err != nil {
					logger.Errorf(ctx, "Failed to delete old bootstrap node: %v", err)
					return err
				}

				if err := bs.db.WithContext(ctx).Create(&node).Error; err != nil {
					logger.Errorf(ctx, "Failed to insert bootstrap node: %v", err)
					return err
				}

				return nil
			})
			if err != nil {
				logger.Errorf(ctx, "Failed to update peers bootstrap node: %v", err)
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
