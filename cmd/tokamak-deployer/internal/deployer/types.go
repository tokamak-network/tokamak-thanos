package deployer

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
}
