package mdutils

import (
	ds "github.com/RealImage/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	dssync "github.com/RealImage/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/sync"
	"github.com/RealImage/go-ipfs/blocks/blockstore"
	bsrv "github.com/RealImage/go-ipfs/blockservice"
	"github.com/RealImage/go-ipfs/exchange/offline"
	dag "github.com/RealImage/go-ipfs/merkledag"
)

func Mock() dag.DAGService {
	bstore := blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))
	bserv := bsrv.New(bstore, offline.Exchange(bstore))
	return dag.NewDAGService(bserv)
}
