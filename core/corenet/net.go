package corenet

import (
	"time"

	core "github.com/RealImage/go-ipfs/core"
	peer "gx/ipfs/QmZpD74pUj6vuxTp1o6LhA3JavC2Bvh9fsWPPVvHnD9sE7/go-libp2p-peer"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	net "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net"
	pro "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/protocol"
)

type ipfsListener struct {
	conCh  chan net.Stream
	proto  pro.ID
	ctx    context.Context
	cancel func()
}

func (il *ipfsListener) Accept() (net.Stream, error) {
	select {
	case c := <-il.conCh:
		return c, nil
	case <-il.ctx.Done():
		return nil, il.ctx.Err()
	}
}

func (il *ipfsListener) Close() error {
	il.cancel()
	// TODO: unregister handler from peerhost
	return nil
}

func Listen(nd *core.IpfsNode, protocol string) (*ipfsListener, error) {
	ctx, cancel := context.WithCancel(nd.Context())

	list := &ipfsListener{
		proto:  pro.ID(protocol),
		conCh:  make(chan net.Stream),
		ctx:    ctx,
		cancel: cancel,
	}

	nd.PeerHost.SetStreamHandler(list.proto, func(s net.Stream) {
		select {
		case list.conCh <- s:
		case <-ctx.Done():
			s.Close()
		}
	})

	return list, nil
}

func Dial(nd *core.IpfsNode, p peer.ID, protocol string) (net.Stream, error) {
	ctx, cancel := context.WithTimeout(nd.Context(), time.Second*30)
	defer cancel()
	err := nd.PeerHost.Connect(ctx, peer.PeerInfo{ID: p})
	if err != nil {
		return nil, err
	}
	return nd.PeerHost.NewStream(nd.Context(), pro.ID(protocol), p)
}
