//go:build opdeployer

package interopgen

import (
	"github.com/ethereum/go-ethereum/core"

	"github.com/tokamak-network/tokamak-thanos/op-node/rollup"
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
	"github.com/tokamak-network/tokamak-thanos/op-supervisor/supervisor/backend/depset"
)

type L1Output struct {
	Genesis *core.Genesis
}

type L2Output struct {
	Genesis   *core.Genesis
	RollupCfg *rollup.Config
}

type WorldOutput struct {
	L1  *L1Output
	L2s map[string]*L2Output
}

func (wo *WorldOutput) RollupConfigSet() depset.StaticRollupConfigSet {
	rcfgs := make(map[eth.ChainID]*rollup.Config)
	for _, rcfg := range wo.L2s {
		rcfgs[eth.ChainIDFromBig(rcfg.RollupCfg.L2ChainID)] = rcfg.RollupCfg
	}
	return depset.StaticRollupConfigSetFromRollupConfigMap(rcfgs, depset.StaticTimestamp(wo.L1.Genesis.Timestamp))
}
