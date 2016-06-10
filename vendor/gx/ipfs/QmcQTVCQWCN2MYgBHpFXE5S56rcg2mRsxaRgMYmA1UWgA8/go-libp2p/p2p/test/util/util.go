package testutil

import (
	"testing"

	peer "gx/ipfs/QmZpD74pUj6vuxTp1o6LhA3JavC2Bvh9fsWPPVvHnD9sE7/go-libp2p-peer"
	bhost "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/host/basic"
	metrics "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/metrics"
	inet "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net"
	swarm "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net/swarm"
	tu "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/testutil"

	ma "gx/ipfs/QmYzDkkgAEmrcNzFCiYo6L1dTX4EAG1gZkbtdbd9trL4vd/go-multiaddr"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

func GenSwarmNetwork(t *testing.T, ctx context.Context) *swarm.Network {
	p := tu.RandPeerNetParamsOrFatal(t)
	ps := peer.NewPeerstore()
	ps.AddPubKey(p.ID, p.PubKey)
	ps.AddPrivKey(p.ID, p.PrivKey)
	n, err := swarm.NewNetwork(ctx, []ma.Multiaddr{p.Addr}, p.ID, ps, metrics.NewBandwidthCounter())
	if err != nil {
		t.Fatal(err)
	}
	ps.AddAddrs(p.ID, n.ListenAddresses(), peer.PermanentAddrTTL)
	return n
}

func DivulgeAddresses(a, b inet.Network) {
	id := a.LocalPeer()
	addrs := a.Peerstore().Addrs(id)
	b.Peerstore().AddAddrs(id, addrs, peer.PermanentAddrTTL)
}

func GenHostSwarm(t *testing.T, ctx context.Context) *bhost.BasicHost {
	n := GenSwarmNetwork(t, ctx)
	return bhost.New(n)
}

var RandPeerID = tu.RandPeerID
