package store

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/ipfs/go-cid"
	ipldgit "github.com/ipfs/go-ipld-git"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage"
	"github.com/multiformats/go-multihash"
	ipldgitprime "github.com/peerforge/go-ipld-gitstore"
)

var _ ipldgitprime.Store = &Blob{}

type Blob struct {
	linkSys ipld.LinkSystem
	store   storage.StreamingWritableStorage
}

func (b *Blob) PutBlob(ctx context.Context, blob ipldgit.Blob) error {
	lb := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    cid.GitRaw,
		MhType:   0x11,
		MhLength: 20,
	}}

	lnk := b.linkSys.MustComputeLink(lb, blob)
	cid := lnk.(cidlink.Link).Cid
	sha := cid.Hash()
	mh, err := multihash.Decode(sha)
	if err != nil {
		return err
	}

	spew.Dump(sha.String(), string(mh.Digest))

	wr, wrCommitter, err := b.store.PutStream(ctx)
	if err != nil {
		return err
	}

	data, err := blob.AsBytes()
	if err != nil {
		return err
	}

	_, err = wr.Write(data)
	if err != nil {
		err = wrCommitter("")
		if err != nil {
			return err
		}

		return err
	}

	err = wrCommitter(cid.String())
	if err != nil {
		return err
	}

	spew.Dump(string(data))
	//res, err := b.store.PutStream(raw, "git-raw", "sha1", -1)
	//if err != nil {
	//	p.errCh <- fmt.Errorf("push/put: %v", err)
	//	return
	//}
	return nil
}
