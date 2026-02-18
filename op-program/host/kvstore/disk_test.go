package kvstore

import (
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-program/host/types"
)

func TestDiskKV(t *testing.T) {
	tmp := t.TempDir()
	kv, err := NewDiskKV(log.New(), tmp, types.DataFormatPebble)
	require.NoError(t, err)
	kvTest(t, kv)
}

func TestCreateMissingDirectory(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "data")
	kv, err := NewDiskKV(log.New(), dir, types.DataFormatPebble)
	require.NoError(t, err)
	val := []byte{1, 2, 3, 4}
	key := crypto.Keccak256Hash(val)
	require.NoError(t, kv.Put(key, val))
}
