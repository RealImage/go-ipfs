package sm_yamux

import (
	"testing"

	test "gx/ipfs/Qmb1US8uyZeEpMyc56wVZy2cDFdQjNFojAUYVCoo9ieTqp/go-stream-muxer/test"
)

func TestYamuxTransport(t *testing.T) {
	test.SubtestAll(t, DefaultTransport)
}
