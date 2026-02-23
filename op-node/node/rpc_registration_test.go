package node

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/tokamak-network/tokamak-thanos/op-node/config"
)

func TestRegisterAPIs_NilHandler(t *testing.T) {
	err := registerAPIs(&config.Config{}, &OpNode{}, nil)
	require.ErrorContains(t, err, "rpc handler is nil")
}

func TestRegisterOptionalAPIs_AdminOnly(t *testing.T) {
	cfg := &config.Config{}
	cfg.RPC.EnableAdmin = true

	node := &OpNode{
		cfg: &config.Config{},
		log: log.Root(),
	}

	var namespaces []string
	err := registerOptionalAPIs(cfg, node, func(api rpc.API) error {
		namespaces = append(namespaces, api.Namespace)
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"admin"}, namespaces)
}

func TestRegisterOptionalAPIs_PropagatesRegistrarError(t *testing.T) {
	cfg := &config.Config{}
	cfg.RPC.EnableAdmin = true

	node := &OpNode{
		cfg: &config.Config{},
		log: log.Root(),
	}

	registerErr := errors.New("boom")
	err := registerOptionalAPIs(cfg, node, func(api rpc.API) error {
		return registerErr
	})
	require.ErrorContains(t, err, "failed to add Admin API")
	require.ErrorIs(t, err, registerErr)
}
