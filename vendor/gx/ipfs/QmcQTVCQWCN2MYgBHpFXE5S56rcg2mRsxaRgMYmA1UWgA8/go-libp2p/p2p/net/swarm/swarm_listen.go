package swarm

import (
	"fmt"

	ma "gx/ipfs/QmYzDkkgAEmrcNzFCiYo6L1dTX4EAG1gZkbtdbd9trL4vd/go-multiaddr"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	transport "gx/ipfs/QmaLnS2kGBLuGZKJdT5KAyoWEtW3u8CS3h6YKCVge5ohD2/go-libp2p-transport"
	mconn "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/metrics/conn"
	inet "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net"
	conn "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net/conn"
	ps "gx/ipfs/Qmcq3bs1zXXoougbpTtLRp7AravkCGBkVyCqJxZRk5qdwo/go-peerstream"
	lgbl "gx/ipfs/Qmdt3dCBXLwr9TEYbmcx6adfP4mNfNArKyWMmYy4EzpHiu/go-libp2p-loggables"
)

// Open listeners and reuse-dialers for the given addresses
func (s *Swarm) setupInterfaces(addrs []ma.Multiaddr) error {
	errs := make([]error, len(addrs))
	var succeeded int
	for i, a := range addrs {
		tpt := s.transportForAddr(a)
		if tpt == nil {
			errs[i] = fmt.Errorf("no transport for address: %s", a)
			continue
		}

		d, err := tpt.Dialer(a, transport.TimeoutOpt(DialTimeout), transport.ReusePorts)
		if err != nil {
			errs[i] = err
			continue
		}

		s.dialer.AddDialer(d)

		list, err := tpt.Listen(a)
		if err != nil {
			errs[i] = err
			continue
		}

		err = s.addListener(list)
		if err != nil {
			errs[i] = err
			continue
		}
		succeeded++
	}

	for i, e := range errs {
		if e != nil {
			log.Warning("listen on %s failed: %s", addrs[i], errs[i])
		}
	}
	if succeeded == 0 && len(addrs) > 0 {
		return fmt.Errorf("failed to listen on any addresses: %s", errs)
	}

	return nil
}

func (s *Swarm) transportForAddr(a ma.Multiaddr) transport.Transport {
	for _, t := range s.transports {
		if t.Matches(a) {
			return t
		}
	}

	return nil
}

func (s *Swarm) addListener(tptlist transport.Listener) error {

	sk := s.peers.PrivKey(s.local)
	if sk == nil {
		// may be fine for sk to be nil, just log a warning.
		log.Warning("Listener not given PrivateKey, so WILL NOT SECURE conns.")
	}

	list, err := conn.WrapTransportListener(s.Context(), tptlist, s.local, sk)
	if err != nil {
		return err
	}

	list.SetAddrFilters(s.Filters)

	if cw, ok := list.(conn.ListenerConnWrapper); ok {
		cw.SetConnWrapper(func(c transport.Conn) transport.Conn {
			return mconn.WrapConn(s.bwc, c)
		})
	}

	return s.addConnListener(list)
}

func (s *Swarm) addConnListener(list conn.Listener) error {
	// AddListener to the peerstream Listener. this will begin accepting connections
	// and streams!
	sl, err := s.swarm.AddListener(list)
	if err != nil {
		return err
	}
	log.Debugf("Swarm Listeners at %s", s.ListenAddresses())

	maddr := list.Multiaddr()

	// signal to our notifiees on successful conn.
	s.notifyAll(func(n inet.Notifiee) {
		n.Listen((*Network)(s), maddr)
	})

	// go consume peerstream's listen accept errors. note, these ARE errors.
	// they may be killing the listener, and if we get _any_ we should be
	// fixing this in our conn.Listener (to ignore them or handle them
	// differently.)
	go func(ctx context.Context, sl *ps.Listener) {

		// signal to our notifiees closing
		defer s.notifyAll(func(n inet.Notifiee) {
			n.ListenClose((*Network)(s), maddr)
		})

		for {
			select {
			case err, more := <-sl.AcceptErrors():
				if !more {
					return
				}
				log.Warningf("swarm listener accept error: %s", err)
			case <-ctx.Done():
				return
			}
		}
	}(s.Context(), sl)

	return nil
}

// connHandler is called by the StreamSwarm whenever a new connection is added
// here we configure it slightly. Note that this is sequential, so if anything
// will take a while do it in a goroutine.
// See https://godoc.org/github.com/jbenet/go-peerstream for more information
func (s *Swarm) connHandler(c *ps.Conn) *Conn {
	ctx := context.Background()
	// this context is for running the handshake, which -- when receiveing connections
	// -- we have no bound on beyond what the transport protocol bounds it at.
	// note that setup + the handshake are bounded by underlying io.
	// (i.e. if TCP or UDP disconnects (or the swarm closes), we're done.
	// Q: why not have a shorter handshake? think about an HTTP server on really slow conns.
	// as long as the conn is live (TCP says its online), it tries its best. we follow suit.)

	sc, err := s.newConnSetup(ctx, c)
	if err != nil {
		log.Debug(err)
		log.Event(ctx, "newConnHandlerDisconnect", lgbl.NetConn(c.NetConn()), lgbl.Error(err))
		c.Close() // boom. close it.
		return nil
	}

	// if a peer dials us, remove from dial backoff.
	s.backf.Clear(sc.RemotePeer())

	return sc
}
