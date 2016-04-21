package s3datastore

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/golang/groupcache/lru"
	datastore "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	query "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/query"
	"github.com/rlmcpherson/s3gof3r"
)

var _ datastore.ThreadSafeDatastore = &S3Datastore{}

var ErrInvalidType = errors.New("s3 datastore: invalid type error")

const (
	maxConcurrentCalls = 128
)

type S3Datastore struct {
	bucket *s3gof3r.Bucket
	cache  *lru.Cache
}

func New(domain, bucketName string) (*S3Datastore, error) {
	k, err := s3gof3r.EnvKeys()
	if err != nil {
		return nil, err
	}

	ds := &S3Datastore{}

	s3 := s3gof3r.New(domain, k)
	if s3 == nil {
		return nil, errors.New("nil s3 object")
	}

	ds.bucket = s3.Bucket(bucketName)
	if ds.bucket == nil {
		return nil, errors.New("nil bucket object")
	}

	ds.bucket.Config = &s3gof3r.Config{}
	*ds.bucket.Config = *s3gof3r.DefaultConfig
	ds.bucket.Config.Md5Check = false

	ds.cache = lru.New(32)

	return ds, nil
}

func (ds *S3Datastore) encode(key datastore.Key) string {
	s := hex.EncodeToString(key.Bytes())

	if len(s) >= 8 {
		s = s[6:8] + s[1:]
	}

	return s
}

func (ds *S3Datastore) decode(raw string) (datastore.Key, bool) {
	k, err := hex.DecodeString(raw)
	if err != nil {
		return datastore.Key{}, false
	}
	return datastore.NewKey(string(k)), true
}

func (ds *S3Datastore) Put(key datastore.Key, value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return ErrInvalidType
	}

	k := ds.encode(key)

	w, err := ds.bucket.PutWriter(k, nil, nil)
	if err != nil {
		return err
	} else if w == nil {
		return errors.New("nil writer")
	}

	defer w.Close()

	buf := bytes.NewBuffer(v)

	n, err := io.Copy(w, buf)
	if err != nil {
		return err
	} else if n != int64(len(v)) {
		return errors.New("value not written fully")
	}

	ds.cache.Add(k, v)

	return nil
}

func (ds *S3Datastore) Get(key datastore.Key) (interface{}, error) {
	k := ds.encode(key)

	b, ok := ds.cache.Get(k)
	if ok {
		return b, nil
	}

	r, _, err := ds.bucket.GetReader(k, nil)
	if err != nil {
		return nil, err
	} else if r == nil {
		return nil, errors.New("nil reader")
	}

	defer r.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	b = buf.Bytes()
	ds.cache.Add(k, b)

	return b, nil
}

func (ds *S3Datastore) Has(key datastore.Key) (bool, error) {
	_, err := ds.Get(key)
	if err != nil {
		respErr, ok := err.(*s3gof3r.RespError)
		if ok {
			if respErr.StatusCode == http.StatusNotFound {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func (ds *S3Datastore) Delete(key datastore.Key) error {
	k := ds.encode(key)
	return ds.bucket.Delete(k)
}

type s3Batch struct {
	puts    map[datastore.Key]interface{}
	deletes map[datastore.Key]struct{}

	ds *S3Datastore
}

func (ds *S3Datastore) Batch() (datastore.Batch, error) {
	return &s3Batch{
		puts:    make(map[datastore.Key]interface{}),
		deletes: make(map[datastore.Key]struct{}),
		ds:      ds,
	}, nil
}

func (bt *s3Batch) Put(key datastore.Key, val interface{}) error {
	bt.puts[key] = val
	return nil
}

func (bt *s3Batch) Delete(key datastore.Key) error {
	bt.deletes[key] = struct{}{}
	return nil
}

func (bt *s3Batch) Commit() error {
	var wg sync.WaitGroup

	errCh := make(chan error, maxConcurrentCalls)
	defer close(errCh)

	i := 0
	for k, v := range bt.puts {
		wg.Add(1)
		go func() {
			err := bt.ds.Put(k, v)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}()
		i++

		if (i%maxConcurrentCalls) == 0 || i == len(bt.puts) {
			wg.Wait()

			select {
			case err := <-errCh:
				return err
			default:
			}
		}
	}

	i = 0
	for k, _ := range bt.deletes {
		wg.Add(1)
		go func() {
			err := bt.ds.Delete(k)
			if err != nil {
				errCh <- err
			}
		}()
		i++

		if (i%maxConcurrentCalls) == 0 || i == len(bt.deletes) {
			wg.Wait()

			select {
			case err := <-errCh:
				return err
			default:
			}
		}
	}

	return nil
}

func (ds *S3Datastore) Query(q query.Query) (query.Results, error) {
	return nil, errors.New("TODO implement query for s3 datastore?")
}

func (ds *S3Datastore) IsThreadSafe() {}
