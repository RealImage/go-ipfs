package conreq

import (
	"sync"

	peer "gx/ipfs/QmZpD74pUj6vuxTp1o6LhA3JavC2Bvh9fsWPPVvHnD9sE7/go-libp2p-peer"
	host "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/host"
	net "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/net"
	protocol "gx/ipfs/QmcQTVCQWCN2MYgBHpFXE5S56rcg2mRsxaRgMYmA1UWgA8/go-libp2p/p2p/protocol"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	logging "gx/ipfs/Qmazh5oNUVsDZTs2g59rq8aYQqwpss8tcUWQzor5sCCEuH/go-log"
)

var log = logging.Logger("conreq")

const ID protocol.ID = "/ipfs/conreq1.0.0"

type ConReqService struct {
	host host.Host

	reqs    map[peer.ID]*conReq
	reqlock sync.Mutex
}

type conReq struct {
	respChan chan net.Conn

	active int
	lk     sync.Mutex
}

func newConReq() *conReq {
	return &conReq{
		respChan: make(chan net.Conn, 1),
	}
}

func (crs *ConReqService) ConnectTo(ctx context.Context, pid peer.ID) (net.Conn, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	crs.reqlock.Lock()
	active, ok := crs.reqs[pid]
	if ok {
		_ = active
		panic("TODO")
	}
	req := newConReq()
	crs.reqs[pid] = req
	crs.reqlock.Unlock()

	peers := crs.getPeersToAsk(ctx, pid)

	for {
		select {
		case p, ok := <-peers:
			if !ok {
				panic("DONE?")
			}
			go crs.sendReqToPeer(ctx, p, pid)
		case <-ctx.Done():
			panic("DONE?")
		case con := <-req.respChan:
			// woohoo! got one!
			return con, nil
		}
	}
}

func (crs *ConReqService) getPeersToAsk(ctx context.Context, trgt peer.ID) <-chan peer.ID {
	panic("NYI")
}

func (crs *ConReqService) sendReqToPeer(ctx context.Context, p, target peer.ID) {
	s, err := crs.host.NewStream(ctx, ID, p)
	if err != nil {
		log.Error(err)
	}

	_ = s
}
