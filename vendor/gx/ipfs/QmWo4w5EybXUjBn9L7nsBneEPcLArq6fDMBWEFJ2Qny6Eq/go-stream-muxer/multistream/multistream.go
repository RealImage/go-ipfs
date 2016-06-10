// package multistream implements a peerstream transport using
// go-multistream to select the underlying stream muxer
package multistream

import (
	"net"

	mss "gx/ipfs/Qmf91yhgRLo2dhhbc5zZ7TxjMaR1oxaWaoc9zRZdi1kU4a/go-multistream"

	smux "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer"
	multiplex "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/multiplex"
	spdy "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/spdystream"
	yamux "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/yamux"
)

type Transport struct {
	mux *mss.MultistreamMuxer

	tpts map[string]smux.Transport

	OrderPreference []string
}

func NewTransport() *Transport {
	mux := mss.NewMultistreamMuxer()
	mux.AddHandler("/multiplex", nil)
	mux.AddHandler("/spdystream", nil)
	mux.AddHandler("/yamux", nil)

	tpts := map[string]smux.Transport{
		"/multiplex":  multiplex.DefaultTransport,
		"/spdystream": spdy.Transport,
		"/yamux":      yamux.DefaultTransport,
	}

	return &Transport{
		mux:             mux,
		tpts:            tpts,
		OrderPreference: []string{"/yamux", "/spdystream", "/multiplex"},
	}
}

func NewBlankTransport() *Transport {
	return &Transport{
		mux:  mss.NewMultistreamMuxer(),
		tpts: make(map[string]smux.Transport),
	}
}

func (t *Transport) AddTransport(path string, tpt smux.Transport) {
	t.mux.AddHandler(path, nil)
	t.tpts[path] = tpt
	t.OrderPreference = append(t.OrderPreference, path)
}

func (t *Transport) NewConn(nc net.Conn, isServer bool) (smux.Conn, error) {
	var proto string
	if isServer {
		selected, _, err := t.mux.Negotiate(nc)
		if err != nil {
			return nil, err
		}
		proto = selected
	} else {
		// prefer yamux
		selected, err := mss.SelectOneOf(t.OrderPreference, nc)
		if err != nil {
			return nil, err
		}
		proto = selected
	}

	tpt := t.tpts[proto]

	return tpt.NewConn(nc, isServer)
}
