package muxtest

import (
	multiplex "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/multiplex"
	multistream "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/multistream"
	muxado "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/muxado"
	spdy "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/spdystream"
	yamux "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/yamux"
)

var _ = multiplex.DefaultTransport
var _ = multistream.NewTransport
var _ = muxado.Transport
var _ = spdy.Transport
var _ = yamux.DefaultTransport
