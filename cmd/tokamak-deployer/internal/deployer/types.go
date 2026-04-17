package deployer

import "math/big"

// DeployOutput is the result of deploy-contracts — serialized to deploy-output.json
type DeployOutput struct {
	L1ChainID uint64 `json:"l1ChainId"`
	L2ChainID uint64 `json:"l2ChainId"`

	AddressManager                   string `json:"AddressManager"`
	L1CrossDomainMessengerProxy      string `json:"L1CrossDomainMessengerProxy"`
	L1ERC721BridgeProxy              string `json:"L1ERC721BridgeProxy"`
	L1StandardBridgeProxy            string `json:"L1StandardBridgeProxy"`
	L2OutputOracleProxy              string `json:"L2OutputOracleProxy"`
	OptimismMintableERC20FactoryProxy string `json:"OptimismMintableERC20FactoryProxy"`
	OptimismPortalProxy              string `json:"OptimismPortalProxy"`
	ProxyAdmin                       string `json:"ProxyAdmin"`
	SystemConfigProxy                string `json:"SystemConfigProxy"`
	SuperchainConfigProxy            string `json:"SuperchainConfigProxy"`
	// Fault proof only
	DisputeGameFactoryProxy  string `json:"DisputeGameFactoryProxy,omitempty"`
	AnchorStateRegistryProxy string `json:"AnchorStateRegistryProxy,omitempty"`
}

// DeployConfig is the input configuration for deploy-contracts
type DeployConfig struct {
	L1RPCURL         string
	PrivateKey       string
	L2ChainID        uint64
	EnableFaultProof bool
	FinalSystemOwner string
	L2OutputOracleSubmissionInterval uint64

	// Gas price control for L1 deployment transactions.
	//
	// Previous versions called SuggestGasPrice per transaction, which added
	// an RPC round-trip per TX (26-32 extra calls) and let each TX race the
	// mempool at slightly different prices. The deploy now resolves a single
	// gas price at startup and reuses it for every TX so the in-built
	// bump-on-timeout retry path rarely fires.
	//
	//   - FixedGasPrice: if set (>0), use it directly. Mirrors forge's
	//     --with-gas-price flag. Still clamped to [Floor, Ceil].
	//   - GasPriceMultiplier: percent applied to SuggestGasPrice when
	//     FixedGasPrice is nil. 0 or unset → 200 (i.e., 2× suggested).
	//   - GasPriceFloor / GasPriceCeil: safety clamps. Default 1 Gwei / 100 Gwei.
	FixedGasPrice      *big.Int
	GasPriceMultiplier int
	GasPriceFloor      *big.Int
	GasPriceCeil       *big.Int
}
