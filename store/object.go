package store

import (
	"context"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage"
	"github.com/ipld/go-ipld-prime/storage/fsstore"
	"github.com/multiformats/go-base32"
	"github.com/rs/zerolog/log"

	ipldgitprime "github.com/drgomesp/go-ipld-gitprime"
)

var _ ipldgitprime.Store = &ObjectStore{}

type ObjectStore struct {
	linkSys *ipld.LinkSystem
	store   storage.StreamingWritableStorage
}

func (b *ObjectStore) Has(ctx context.Context, key string) (bool, error) {
	return b.store.(storage.ReadableStorage).Has(ctx, key)
}

func (b *ObjectStore) Get(ctx context.Context, key string) ([]byte, error) {
	return b.store.(storage.ReadableStorage).Get(ctx, key)
}

func (b *ObjectStore) Put(ctx context.Context, key string, data []byte) error {
	log.Trace().Msgf("Put(key=%v, data=%v)", key, string(data))
	return b.store.(storage.WritableStorage).Put(ctx, key, data)
}

var b32encoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func b32enc(in string) string {
	return b32encoder.EncodeToString([]byte(in))
}

func NewObjectStore(ls *ipld.LinkSystem) (*ObjectStore, error) {
	st := &fsstore.Store{}
	err := os.MkdirAll(".ipld", os.ModePerm)
	if err != nil {
		return nil, err
	}

	err = st.Init(".ipld/", b32enc, func(key string, shards *[]string) {
		*shards = append(*shards, key)
	})

	if err != nil {
		return nil, err
	}

	return &ObjectStore{
		linkSys: ls,
		store:   st,
	}, nil
}

func (b *ObjectStore) PutObject(ctx context.Context, node ipld.Node) (datamodel.Link, error) {
	lb := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    cid.GitRaw,
		MhType:   0x11,
		MhLength: 20,
	}}

	lnk, err := b.linkSys.Store(
		linking.LinkContext{},
		lb,
		node,
	)

	if err != nil {
		return nil, err
	}

	return lnk, nil
}

func (b *ObjectStore) GetObject(ctx context.Context, key string) ([]byte, error) {
	return b.store.(storage.ReadableStorage).Get(ctx, key)
}
