package opcm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/foundry"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/script"
	opcrypto "github.com/tokamak-network/tokamak-thanos/op-service/crypto"
)

type DeploySuperchainInput struct {
	SuperchainProxyAdminOwner  common.Address         `toml:"superchainProxyAdminOwner"`
	ProtocolVersionsOwner      common.Address         `toml:"protocolVersionsOwner"`
	Guardian                   common.Address         `toml:"guardian"`
	Paused                     bool                   `toml:"paused"`
	RequiredProtocolVersion    params.ProtocolVersion `toml:"requiredProtocolVersion"`
	RecommendedProtocolVersion params.ProtocolVersion `toml:"recommendedProtocolVersion"`
}

func (dsi *DeploySuperchainInput) InputSet() bool {
	return true
}

type DeploySuperchainOutput struct {
	SuperchainProxyAdmin  common.Address `json:"proxyAdminAddress"`
	SuperchainConfigImpl  common.Address `json:"superchainConfigImplAddress"`
	SuperchainConfigProxy common.Address `json:"superchainConfigProxyAddress"`
	ProtocolVersionsImpl  common.Address `json:"protocolVersionsImplAddress"`
	ProtocolVersionsProxy common.Address `json:"protocolVersionsProxyAddress"`
}

func (output *DeploySuperchainOutput) CheckOutput(input common.Address) error {
	return nil
}

type DeploySuperchainOpts struct {
	ChainID     *big.Int
	ArtifactsFS foundry.StatDirFs
	Deployer    common.Address
	Signer      opcrypto.SignerFn
	Input       DeploySuperchainInput
	Client      *ethclient.Client
	Logger      log.Logger
}

func DeploySuperchain(h *script.Host, input DeploySuperchainInput) (DeploySuperchainOutput, error) {
	return RunScriptSingle[DeploySuperchainInput, DeploySuperchainOutput](h, input, "DeploySuperchain.s.sol", "DeploySuperchain")
}
