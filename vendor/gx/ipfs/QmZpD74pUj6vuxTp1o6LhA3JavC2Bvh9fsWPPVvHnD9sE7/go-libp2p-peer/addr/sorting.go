package addr

import (
	"bytes"

	manet "gx/ipfs/QmUBa4w6CbHJUMeGJPDiMEDWsM93xToK1fTnFXnrC8Hksw/go-multiaddr-net"
	ma "gx/ipfs/QmYzDkkgAEmrcNzFCiYo6L1dTX4EAG1gZkbtdbd9trL4vd/go-multiaddr"
	mafmt "gx/ipfs/QmeLQ13LftT9XhNn22piZc3GP56fGqhijuL5Y8KdUaRn1g/mafmt"
)

func isFDCostlyTransport(a ma.Multiaddr) bool {
	return mafmt.TCP.Matches(a)
}

type AddrList []ma.Multiaddr

func (al AddrList) Len() int {
	return len(al)
}

func (al AddrList) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}

func (al AddrList) Less(i, j int) bool {
	a := al[i]
	b := al[j]

	// dial localhost addresses next, they should fail immediately
	lba := manet.IsIPLoopback(a)
	lbb := manet.IsIPLoopback(b)
	if lba {
		if !lbb {
			return true
		}
	}

	// dial utp and similar 'non-fd-consuming' addresses first
	fda := isFDCostlyTransport(a)
	fdb := isFDCostlyTransport(b)
	if !fda {
		if fdb {
			return true
		}

		// if neither consume fd's, assume equal ordering
		return false
	}

	// if 'b' doesnt take a file descriptor
	if !fdb {
		return false
	}

	// if 'b' is loopback and both take file descriptors
	if lbb {
		return false
	}

	// for the rest, just sort by bytes
	return bytes.Compare(a.Bytes(), b.Bytes()) > 0
}
