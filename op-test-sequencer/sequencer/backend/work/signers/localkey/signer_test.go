package localkey

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer/seqtypes"
)

func TestSigner(t *testing.T) {
	logger := testlog.Logger(t, log.LevelInfo)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	chainID := eth.ChainIDFromUInt64(123)
	id := seqtypes.SignerID("foobar")

	signer := NewSigner(id, logger, chainID, key)
	testSigner(t, signer, chainID, addr)
}
