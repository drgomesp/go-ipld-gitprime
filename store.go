package ipldgitprime

import (
	"context"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/storage"
)

type Store interface {
	storage.ReadableStorage
	storage.WritableStorage

	GetObject(ctx context.Context, key string) ([]byte, error)
	PutObject(ctx context.Context, obj ipld.Node) (datamodel.Link, error)
}
