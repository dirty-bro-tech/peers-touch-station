package libp2p

// Add these methods to implement network.Notifiee interface

import (
	"context"

	"github.com/dirty-bro-tech/peers-touch-go/core/logger"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

func (bs *BootstrapServer) Listen(network.Network, ma.Multiaddr)      {}
func (bs *BootstrapServer) ListenClose(network.Network, ma.Multiaddr) {}
func (bs *BootstrapServer) Connected(net network.Network, conn network.Conn) {
	logger.Info(context.Background(), "New client connection",
		"peerID", conn.RemotePeer(),
		"remoteAddress", conn.RemoteMultiaddr())
}
func (bs *BootstrapServer) Disconnected(net network.Network, conn network.Conn) {}
func (bs *BootstrapServer) OpenedStream(net network.Network, s network.Stream)  {}
func (bs *BootstrapServer) ClosedStream(net network.Network, s network.Stream)  {}
