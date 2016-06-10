package conreq

import (
	ma "gx/ipfs/QmYzDkkgAEmrcNzFCiYo6L1dTX4EAG1gZkbtdbd9trL4vd/go-multiaddr"
	net "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net"
)

type notifee ConReqService

func (n *notifee) Connected(nw net.Network, c net.Conn) {

}

func (n *notifee) ClosedStream(nw net.Network, s net.Stream) {

}

func (n *notifee) Disconnected(nw net.Network, c net.Conn) {

}

func (n *notifee) Listen(nw net.Network, a ma.Multiaddr) {

}

func (n *notifee) ListenClose(nw net.Network, a ma.Multiaddr) {

}

func (n *notifee) OpenedStream(nw net.Network, s net.Stream) {

}

var _ net.Notifiee = (*notifee)(nil)
