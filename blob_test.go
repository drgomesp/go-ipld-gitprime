package ipldgitprime_test

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ipfs/go-cid"
	ipldgit "github.com/ipfs/go-ipld-git"
	"github.com/multiformats/go-multihash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.NewConsoleWriter())
}

func Test_Parse(t *testing.T) {
	err := filepath.Walk(path.Join("..", "..", ".git", "objects"), func(path string, info os.FileInfo, err error) error {
		assert.NoError(t, err)

		if info.IsDir() {
			return nil
		}

		parts := strings.Split(path, string(filepath.Separator))
		if dir := parts[len(parts)-2]; dir == "info" || dir == "pack" {
			return nil
		}

		fi, err := os.Open(path)
		assert.NoError(t, err)

		nd, err := ipldgit.ParseCompressedObject(fi)

		assert.NoError(t, err)
		assert.NotNil(t, nd)

		return nil
	})

	assert.NoError(t, err)
}

func Test_ParseCommit(t *testing.T) {
	hash := "4f51f04f3f558f70f258b6743faed485b6d75dbb"
	expectedCid, err := CidFromHex(hash)
	assert.NoError(t, err)

	repo, err := git.PlainOpen(path.Join("..", ".."))
	assert.NoError(t, err)

	obj, err := repo.Storer.EncodedObject(plumbing.CommitObject, plumbing.NewHash(hash))
	assert.NoError(t, err)

	log.Debug().
		Str("cid", expectedCid.String()).
		Str("type", obj.Type().String()).
		Str("hash", obj.Hash().String()).
		Msgf("aaa")
}

func CidFromHex(sha string) (cid.Cid, error) {
	mhash, err := multihash.FromHexString("1114" + sha)
	if err != nil {
		return cid.Undef, err
	}

	return cid.NewCidV1(0x78, mhash), nil
}
