package reprovide_test

import (
	"testing"

	ds "github.com/RealImage/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	dssync "github.com/RealImage/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/sync"
	blocks "github.com/RealImage/go-ipfs/blocks"
	blockstore "github.com/RealImage/go-ipfs/blocks/blockstore"
	mock "github.com/RealImage/go-ipfs/routing/mock"
	testutil "github.com/RealImage/go-ipfs/thirdparty/testutil"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"

	. "github.com/RealImage/go-ipfs/exchange/reprovide"
)

func TestReprovide(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mrserv := mock.NewServer()

	idA := testutil.RandIdentityOrFatal(t)
	idB := testutil.RandIdentityOrFatal(t)

	clA := mrserv.Client(idA)
	clB := mrserv.Client(idB)

	bstore := blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))

	blk := blocks.NewBlock([]byte("this is a test"))
	bstore.Put(blk)

	reprov := NewReprovider(clA, bstore)
	err := reprov.Reprovide(ctx)
	if err != nil {
		t.Fatal(err)
	}

	provs, err := clB.FindProviders(ctx, blk.Key())
	if err != nil {
		t.Fatal(err)
	}

	if len(provs) == 0 {
		t.Fatal("Should have gotten a provider")
	}

	if provs[0].ID != idA.ID() {
		t.Fatal("Somehow got the wrong peer back as a provider.")
	}
}
