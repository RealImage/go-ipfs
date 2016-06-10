package peerstream_spdystream

import (
	"testing"

	test "gx/ipfs/QmWo4w5EybXUjBn9L7nsBneEPcLArq6fDMBWEFJ2Qny6Eq/go-stream-muxer/test"
)

func TestSpdyStreamTransport(t *testing.T) {
	test.SubtestAll(t, Transport)
}
