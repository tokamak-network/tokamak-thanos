package predeploys

import (
	"github.com/ethereum/go-ethereum/common"
)

// Preset constants for genesis predeploy selection.
const (
	PresetGeneral = "general"
	PresetDeFi    = "defi"
	PresetGaming  = "gaming"
	PresetFull    = "full"
)

type DeployConfig interface {
	GovernanceEnabled() bool
	CanyonTime(genesisTime uint64) *uint64
	PresetID() string
}

type Predeploy struct {
	Address       common.Address
	ProxyDisabled bool
	Enabled       func(config DeployConfig) bool
}
