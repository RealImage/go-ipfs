package swarm

import (
	"testing"

	peer "gx/ipfs/QmZpD74pUj6vuxTp1o6LhA3JavC2Bvh9fsWPPVvHnD9sE7/go-libp2p-peer"
	metrics "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/metrics"
	addrutil "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net/swarm/addr"
	testutil "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/testutil"

	ma "gx/ipfs/QmYzDkkgAEmrcNzFCiYo6L1dTX4EAG1gZkbtdbd9trL4vd/go-multiaddr"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

func TestFilterAddrs(t *testing.T) {

	m := func(s string) ma.Multiaddr {
		maddr, err := ma.NewMultiaddr(s)
		if err != nil {
			t.Fatal(err)
		}
		return maddr
	}

	bad := []ma.Multiaddr{
		m("/ip4/1.2.3.4/udp/1234"),           // unreliable
		m("/ip4/1.2.3.4/udp/1234/sctp/1234"), // not in manet
		m("/ip4/1.2.3.4/udp/1234/udt"),       // udt is broken on arm
		m("/ip6/fe80::1/tcp/0"),              // link local
		m("/ip6/fe80::100/tcp/1234"),         // link local
	}

	good := []ma.Multiaddr{
		m("/ip4/127.0.0.1/tcp/0"),
		m("/ip6/::1/tcp/0"),
		m("/ip4/1.2.3.4/udp/1234/utp"),
	}

	goodAndBad := append(good, bad...)

	// test filters

	for _, a := range bad {
		if addrutil.AddrUsable(a, false) {
			t.Errorf("addr %s should be unusable", a)
		}
	}

	for _, a := range good {
		if !addrutil.AddrUsable(a, false) {
			t.Errorf("addr %s should be usable", a)
		}
	}

	subtestAddrsEqual(t, addrutil.FilterUsableAddrs(bad), []ma.Multiaddr{})
	subtestAddrsEqual(t, addrutil.FilterUsableAddrs(good), good)
	subtestAddrsEqual(t, addrutil.FilterUsableAddrs(goodAndBad), good)

	// now test it with swarm

	id, err := testutil.RandPeerID()
	if err != nil {
		t.Fatal(err)
	}

	ps := peer.NewPeerstore()
	ctx := context.Background()

	if _, err := NewNetwork(ctx, bad, id, ps, metrics.NewBandwidthCounter()); err == nil {
		t.Fatal("should have failed to create swarm")
	}

	if _, err := NewNetwork(ctx, goodAndBad, id, ps, metrics.NewBandwidthCounter()); err != nil {
		t.Fatal("should have succeeded in creating swarm", err)
	}
}

func subtestAddrsEqual(t *testing.T, a, b []ma.Multiaddr) {
	if len(a) != len(b) {
		t.Error(t)
	}

	in := func(addr ma.Multiaddr, l []ma.Multiaddr) bool {
		for _, addr2 := range l {
			if addr.Equal(addr2) {
				return true
			}
		}
		return false
	}

	for _, aa := range a {
		if !in(aa, b) {
			t.Errorf("%s not in %s", aa, b)
		}
	}
}

func TestDialBadAddrs(t *testing.T) {

	m := func(s string) ma.Multiaddr {
		maddr, err := ma.NewMultiaddr(s)
		if err != nil {
			t.Fatal(err)
		}
		return maddr
	}

	ctx := context.Background()
	s := makeSwarms(ctx, t, 1)[0]

	test := func(a ma.Multiaddr) {
		p := testutil.RandPeerIDFatal(t)
		s.peers.AddAddr(p, a, peer.PermanentAddrTTL)
		if _, err := s.Dial(ctx, p); err == nil {
			t.Error("swarm should not dial: %s", m)
		}
	}

	test(m("/ip6/fe80::1"))                // link local
	test(m("/ip6/fe80::100"))              // link local
	test(m("/ip4/127.0.0.1/udp/1234/utp")) // utp
}
