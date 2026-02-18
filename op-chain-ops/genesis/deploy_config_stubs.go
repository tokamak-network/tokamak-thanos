package genesis

// Stub types for FP upstream compatibility.
// These deploy config sub-types are used by op-deployer but don't exist in old tokamak-thanos.

type AltDADeployConfig struct{}

type DevDeployConfig struct{}

type EIP1559DeployConfig struct {
	EIP1559Denominator       uint64 `json:"eip1559Denominator"`
	EIP1559DenominatorCanyon uint64 `json:"eip1559DenominatorCanyon"`
	EIP1559Elasticity        uint64 `json:"eip1559Elasticity"`
}

type GasPriceOracleDeployConfig struct {
	GasPriceOracleOverhead      uint64 `json:"gasPriceOracleOverhead"`
	GasPriceOracleScalar        uint64 `json:"gasPriceOracleScalar"`
	GasPriceOracleBaseFeeScalar uint32 `json:"gasPriceOracleBaseFeeScalar"`
	GasPriceOracleBlobBaseFeeScalar uint32 `json:"gasPriceOracleBlobBaseFeeScalar"`
}

type GovernanceDeployConfig struct{}

type L1DependenciesConfig struct{}

type L2GenesisBlockDeployConfig struct {
	L2GenesisBlockNonce         uint64 `json:"l2GenesisBlockNonce"`
	L2GenesisBlockGasLimit      uint64 `json:"l2GenesisBlockGasLimit"`
	L2GenesisBlockDifficulty    uint64 `json:"l2GenesisBlockDifficulty"`
	L2GenesisBlockMixHash       string `json:"l2GenesisBlockMixHash"`
	L2GenesisBlockNumber        uint64 `json:"l2GenesisBlockNumber"`
	L2GenesisBlockGasUsed       uint64 `json:"l2GenesisBlockGasUsed"`
	L2GenesisBlockParentHash    string `json:"l2GenesisBlockParentHash"`
	L2GenesisBlockBaseFeePerGas uint64 `json:"l2GenesisBlockBaseFeePerGas"`
}

type L2InitializationConfig struct {
	L2GenesisBlockDeployConfig
	L2VaultsDeployConfig
	GasPriceOracleDeployConfig
	EIP1559DeployConfig
}

type L2VaultsDeployConfig struct{}
