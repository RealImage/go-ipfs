package peerstream_multiplex

import (
	"testing"

	test "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/test"
)

func TestMultiplexTransport(t *testing.T) {
	test.SubtestAll(t, DefaultTransport)
}
