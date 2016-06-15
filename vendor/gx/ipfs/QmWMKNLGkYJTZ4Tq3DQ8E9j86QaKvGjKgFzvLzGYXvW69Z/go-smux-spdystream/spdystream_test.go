package peerstream_spdystream

import (
	"testing"

	test "gx/ipfs/Qmb1US8uyZeEpMyc56wVZy2cDFdQjNFojAUYVCoo9ieTqp/go-stream-muxer/test"
)

func TestSpdyStreamTransport(t *testing.T) {
	test.SubtestAll(t, Transport)
}
