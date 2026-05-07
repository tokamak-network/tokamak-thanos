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
	DelayedWETHProxy         string `json:"DelayedWETHProxy,omitempty"`

	// Implementations records the runtime impl address for each reuse-target contract.
	// When reuse hits, this is the on-chain address from the registry; otherwise it is the
	// freshly-deployed address. Useful for curating a future registry update.
	Implementations map[string]string `json:"implementations,omitempty"`
}

// DeployConfig is the input configuration for deploy-contracts
type DeployConfig struct {
	L1RPCURL         string
	PrivateKey       string
	L2ChainID        uint64
	EnableFaultProof bool
	// DelayedWETHDelay is the withdrawal delay in seconds for the DelayedWETH bond escrow.
	// 0 is valid for local testnets (no enforced delay).
	DelayedWETHDelay uint64
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

	// Reuse-deployment knobs.
	//
	// ReuseDeployment: master toggle. When false (default), behavior is identical to
	// pre-reuse releases — every impl is freshly deployed.
	//
	// RegistryPath: optional override. When non-empty, the file at this path is loaded
	// instead of the embedded registry/{l1ChainID}.json.
	//
	// ReuseStrict: if true, any registry verification failure (missing on-chain code,
	// bytecode hash mismatch, invalid hex) aborts the deploy instead of falling back to
	// a fresh deploy for that entry.
	ReuseDeployment bool
	RegistryPath    string
	ReuseStrict     bool
}
