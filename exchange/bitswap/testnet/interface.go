package bitswap

import (
	bsnet "github.com/RealImage/go-ipfs/exchange/bitswap/network"
	"github.com/RealImage/go-ipfs/thirdparty/testutil"
	peer "gx/ipfs/QmbyvM8zRFDkbFdYyt1MnevUMJ62SiSGbfDFZ3Z8nkrzr4/go-libp2p-peer"
)

type Network interface {
	Adapter(testutil.Identity) bsnet.BitSwapNetwork

	HasPeer(peer.ID) bool
}
