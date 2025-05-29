package libp2p

import (
	"context"

	log "github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/dirty-bro-tech/peers-touch-go/core/plugin/registry/native"
	dht_pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"github.com/libp2p/go-libp2p/core/network"
)

func Init() {
	native.AppendDhtRequestHook(func(ctx context.Context, s network.Stream, req *dht_pb.Message) {
		/*connectChan <- connectEvent{
			peerID: s.Conn().RemotePeer(),
			msg:    req,
		}*/

		log.Infof(ctx, "got a dht request from: %s; type: %s; msg: %+v", s.Conn().RemotePeer().String(), req.Type, req)
	})
}
