package fsrepo

import (
	"encoding/json"
	"fmt"
	"path"

	ds "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	levelds "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/leveldb"
	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/measure"
	mount "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/syncmount"
	ldbopts "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/syndtr/goleveldb/leveldb/opt"
	repo "github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/thirdparty/s3datastore"
)

type S3Params struct {
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
}

func openS3Datastore(r *FSRepo) (repo.Datastore, error) {
	leveldbPath := path.Join(r.path, leveldbDirectory)

	// save leveldb reference so it can be neatly closed afterward
	leveldbDS, err := levelds.NewDatastore(leveldbPath, &levelds.Options{
		Compression: ldbopts.NoCompression,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to open leveldb datastore: %v", err)
	}

	config, err := r.Config()
	if err != nil {
		return nil, fmt.Errorf("unable to get config: %v", err)
	}

	if config == nil {
		return nil, fmt.Errorf("config is empty")
	}

	if config.Datastore.Params == nil {
		return nil, fmt.Errorf("no params specified")
	}

	p := S3Params{}
	err = json.Unmarshal(*config.Datastore.Params, &p)
	if err != nil {
		return nil, fmt.Errorf("unable to parse S3 params: %v", err)
	}

	blocksDS, err := s3datastore.New("s3-"+p.Region+".amazonaws.com", p.Bucket, p.AccessKey, p.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("unable to open s3 datastore: %v", err)
	}

	// Add our PeerID to metrics paths to keep them unique
	//
	// As some tests just pass a zero-value Config to fsrepo.Init,
	// cope with missing PeerID.
	id := r.config.Identity.PeerID
	if id == "" {
		// the tests pass in a zero Config; cope with it
		id = fmt.Sprintf("uninitialized_%p", r)
	}
	prefix := "fsrepo." + id + ".datastore."
	metricsBlocks := measure.New(prefix+"blocks", blocksDS)
	metricsLevelDB := measure.New(prefix+"leveldb", leveldbDS)
	mountDS := mount.New([]mount.Mount{
		{
			Prefix:    ds.NewKey("/blocks"),
			Datastore: metricsBlocks,
		},
		{
			Prefix:    ds.NewKey("/"),
			Datastore: metricsLevelDB,
		},
	})

	return mountDS, nil
}
