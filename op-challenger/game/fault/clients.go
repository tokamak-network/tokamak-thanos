package fault

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/config"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/trace/utils"
	"github.com/tokamak-network/tokamak-thanos/op-service/dial"
)

type clientProvider struct {
	ctx            context.Context
	logger         log.Logger
	cfg            *config.Config
	l2HeaderSource utils.L2HeaderSource
	rollupClient   RollupClient
	syncValidator  *syncStatusValidator
	toClose        []CloseFunc
}

func (c *clientProvider) Close() {
	for _, closeFunc := range c.toClose {
		closeFunc()
	}
}

func (c *clientProvider) SingleChainClients() (utils.L2HeaderSource, RollupClient, *syncStatusValidator, error) {
	headers, err := c.L2HeaderSource()
	if err != nil {
		return nil, nil, nil, err
	}
	rollup, err := c.RollupClient()
	if err != nil {
		return nil, nil, nil, err
	}
	return headers, rollup, c.syncValidator, nil
}

func (c *clientProvider) L2HeaderSource() (utils.L2HeaderSource, error) {
	if c.l2HeaderSource != nil {
		return c.l2HeaderSource, nil
	}

	// Tokamak-Thanos uses L2Rpc (single string) instead of L2Rpcs (array)
	l2Client, err := ethclient.DialContext(c.ctx, c.cfg.L2Rpc)
	if err != nil {
		return nil, fmt.Errorf("dial l2 client %v: %w", c.cfg.L2Rpc, err)
	}
	c.l2HeaderSource = l2Client
	c.toClose = append(c.toClose, l2Client.Close)
	return l2Client, nil
}

func (c *clientProvider) RollupClient() (RollupClient, error) {
	if c.rollupClient != nil {
		return c.rollupClient, nil
	}
	// Tokamak-Thanos requires timeout parameter
	rollupClient, err := dial.DialRollupClientWithTimeout(c.ctx, 30*time.Second, c.logger, c.cfg.RollupRpc)
	if err != nil {
		return nil, fmt.Errorf("dial rollup client %v: %w", c.cfg.RollupRpc, err)
	}
	c.rollupClient = rollupClient
	c.syncValidator = newSyncStatusValidator(rollupClient)
	c.toClose = append(c.toClose, rollupClient.Close)
	return rollupClient, nil
}
