package host

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-node/chaincfg"
	preimage "github.com/tokamak-network/tokamak-thanos/op-preimage"
	"github.com/tokamak-network/tokamak-thanos/op-program/chainconfig"
	"github.com/tokamak-network/tokamak-thanos/op-program/client"
	"github.com/tokamak-network/tokamak-thanos/op-program/client/l1"
	"github.com/tokamak-network/tokamak-thanos/op-program/host/config"
	"github.com/tokamak-network/tokamak-thanos/op-program/host/kvstore"
	"github.com/tokamak-network/tokamak-thanos/op-program/io"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"
)

func TestServerMode(t *testing.T) {
	dir := t.TempDir()

	l1Head := common.Hash{0x11}
	l2OutputRoot := common.Hash{0x33}
	cfg := config.NewConfig(chaincfg.Sepolia, chainconfig.OPSepoliaChainConfig, l1Head, common.Hash{0x22}, l2OutputRoot, common.Hash{0x44}, 1000)
	cfg.DataDir = dir
	cfg.ServerMode = true

	preimageServer, preimageClient, err := io.CreateBidirectionalChannel()
	require.NoError(t, err)
	defer preimageClient.Close()
	hintServer, hintClient, err := io.CreateBidirectionalChannel()
	require.NoError(t, err)
	defer hintClient.Close()
	logger := testlog.Logger(t, log.LevelTrace)
	result := make(chan error)
	go func() {
		result <- PreimageServer(context.Background(), logger, cfg, preimageServer, hintServer)
	}()

	pClient := preimage.NewOracleClient(preimageClient)
	hClient := preimage.NewHintWriter(hintClient)
	l1PreimageOracle := l1.NewPreimageOracle(pClient, hClient)

	require.Equal(t, l1Head.Bytes(), pClient.Get(client.L1HeadLocalIndex), "Should get l1 head preimages")
	require.Equal(t, l2OutputRoot.Bytes(), pClient.Get(client.L2OutputRootLocalIndex), "Should get l2 output root preimages")

	// Should exit when a preimage is unavailable
	require.Panics(t, func() {
		l1PreimageOracle.HeaderByBlockHash(common.HexToHash("0x1234"))
	}, "Preimage should not be available")
	require.ErrorIs(t, waitFor(result), kvstore.ErrNotFound)
}

func waitFor(ch chan error) error {
	timeout := time.After(30 * time.Second)
	select {
	case err := <-ch:
		return err
	case <-timeout:
		return errors.New("timed out")
	}
}
