// Package superchain provides a shim for the op-geth superchain registry.
// This replaces github.com/ethereum/go-ethereum/superchain which only exists in op-geth.
// For tokamak-thanos, chain configs come from the superchain-registry or are hardcoded.
package superchain

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var ErrUnknownChain = errors.New("unknown chain")

func addrPtr(hex string) *common.Address {
	a := common.HexToAddress(hex)
	return &a
}

// Superchain describes a superchain network (e.g. mainnet, sepolia).
type Superchain struct {
	Name                   string         `toml:"name"`
	ProtocolVersionsAddr   common.Address `toml:"protocol_versions_addr"`
	SuperchainConfigAddr   common.Address `toml:"superchain_config_addr"`
	OpContractsManagerAddr common.Address `toml:"op_contracts_manager_addr"`
	Hardforks              HardforkConfig
	L1                     L1Config
}

type L1Config struct {
	ChainID   uint64 `toml:"chain_id"`
	PublicRPC string `toml:"public_rpc"`
	Explorer  string `toml:"explorer"`
}

type ChainConfig struct {
	Name                 string       `toml:"name"`
	PublicRPC            string       `toml:"public_rpc"`
	SequencerRPC         string       `toml:"sequencer_rpc"`
	Explorer             string       `toml:"explorer"`
	SuperchainLevel      int          `toml:"superchain_level"`
	GovernedByOptimism   bool         `toml:"governed_by_optimism"`
	SuperchainTime       *uint64      `toml:"superchain_time"`
	DataAvailabilityType string       `toml:"data_availability_type"`
	DeploymentTxHash     *common.Hash `toml:"deployment_tx_hash"`

	ChainID           uint64          `toml:"chain_id"`
	BatchInboxAddr    common.Address  `toml:"batch_inbox_addr"`
	BlockTime         uint64          `toml:"block_time"`
	SeqWindowSize     uint64          `toml:"seq_window_size"`
	MaxSequencerDrift uint64          `toml:"max_sequencer_drift"`
	GasPayingToken    *common.Address `toml:"gas_paying_token"`
	Hardforks         HardforkConfig  `toml:"hardforks"`
	Interop           *Interop        `toml:"interop,omitempty"`
	Optimism          *OptimismConfig `toml:"optimism,omitempty"`
	AltDA             *AltDAConfig    `toml:"alt_da,omitempty"`
	Genesis           GenesisConfig   `toml:"genesis"`
	Roles             RolesConfig     `toml:"roles"`
	Addresses         AddressesConfig `toml:"addresses"`
	Network           string          // populated at load time
}

type Dependency struct{}

type Interop struct {
	Dependencies map[string]Dependency `json:"dependencies" toml:"dependencies"`
}

type HardforkConfig struct {
	CanyonTime             *uint64 `toml:"canyon_time"`
	DeltaTime              *uint64 `toml:"delta_time"`
	EcotoneTime            *uint64 `toml:"ecotone_time"`
	FjordTime              *uint64 `toml:"fjord_time"`
	GraniteTime            *uint64 `toml:"granite_time"`
	HoloceneTime           *uint64 `toml:"holocene_time"`
	IsthmusTime            *uint64 `toml:"isthmus_time"`
	JovianTime             *uint64 `toml:"jovian_time"`
	InteropTime            *uint64 `toml:"interop_time"`
	PectraBlobScheduleTime *uint64 `toml:"pectra_blob_schedule_time,omitempty"`
}

type OptimismConfig struct {
	EIP1559Elasticity        uint64  `toml:"eip1559_elasticity"`
	EIP1559Denominator       uint64  `toml:"eip1559_denominator"`
	EIP1559DenominatorCanyon *uint64 `toml:"eip1559_denominator_canyon"`
}

type AltDAConfig struct {
	DaChallengeContractAddress common.Address `toml:"da_challenge_contract_address"`
	DaChallengeWindow          uint64         `toml:"da_challenge_window"`
	DaResolveWindow            uint64         `toml:"da_resolve_window"`
	DaCommitmentType           string         `toml:"da_commitment_type"`
}

type GenesisConfig struct {
	L2Time       uint64       `toml:"l2_time"`
	L1           GenesisRef   `toml:"l1"`
	L2           GenesisRef   `toml:"l2"`
	SystemConfig SystemConfig `toml:"system_config"`
}

type GenesisRef struct {
	Hash   common.Hash `toml:"hash"`
	Number uint64      `toml:"number"`
}

type SystemConfig struct {
	BatcherAddr       common.Address `json:"batcherAddr" toml:"batcherAddress"`
	Overhead          common.Hash    `json:"overhead" toml:"overhead"`
	Scalar            common.Hash    `json:"scalar" toml:"scalar"`
	GasLimit          uint64         `json:"gasLimit" toml:"gasLimit"`
	BaseFeeScalar     *uint64        `json:"baseFeeScalar,omitempty" toml:"baseFeeScalar,omitempty"`
	BlobBaseFeeScalar *uint64        `json:"blobBaseFeeScalar,omitempty" toml:"blobBaseFeeScalar,omitempty"`
}

type RolesConfig struct {
	SystemConfigOwner *common.Address `json:"SystemConfigOwner" toml:"SystemConfigOwner"`
	ProxyAdminOwner   *common.Address `json:"ProxyAdminOwner" toml:"ProxyAdminOwner"`
	Guardian          *common.Address `json:"Guardian" toml:"Guardian"`
	Challenger        *common.Address `json:"Challenger" toml:"Challenger"`
	Proposer          *common.Address `json:"Proposer,omitempty" toml:"Proposer,omitempty"`
	UnsafeBlockSigner *common.Address `json:"UnsafeBlockSigner,omitempty" toml:"UnsafeBlockSigner,omitempty"`
	BatchSubmitter    *common.Address `json:"BatchSubmitter" toml:"BatchSubmitter"`
}

type AddressesConfig struct {
	AddressManager                    *common.Address `toml:"AddressManager,omitempty" json:"AddressManager,omitempty"`
	L1CrossDomainMessengerProxy       *common.Address `toml:"L1CrossDomainMessengerProxy,omitempty" json:"L1CrossDomainMessengerProxy,omitempty"`
	L1ERC721BridgeProxy               *common.Address `toml:"L1ERC721BridgeProxy,omitempty" json:"L1ERC721BridgeProxy,omitempty"`
	L1StandardBridgeProxy             *common.Address `toml:"L1StandardBridgeProxy,omitempty" json:"L1StandardBridgeProxy,omitempty"`
	L2OutputOracleProxy               *common.Address `toml:"L2OutputOracleProxy,omitempty" json:"L2OutputOracleProxy,omitempty"`
	OptimismMintableERC20FactoryProxy *common.Address `toml:"OptimismMintableERC20FactoryProxy,omitempty" json:"OptimismMintableERC20FactoryProxy,omitempty"`
	OptimismPortalProxy               *common.Address `toml:"OptimismPortalProxy,omitempty" json:"OptimismPortalProxy,omitempty"`
	SystemConfigProxy                 *common.Address `toml:"SystemConfigProxy,omitempty" json:"SystemConfigProxy,omitempty"`
	ProxyAdmin                        *common.Address `toml:"ProxyAdmin,omitempty" json:"ProxyAdmin,omitempty"`
	SuperchainConfig                  *common.Address `toml:"SuperchainConfig,omitempty" json:"SuperchainConfig,omitempty"`
	AnchorStateRegistryProxy          *common.Address `toml:"AnchorStateRegistryProxy,omitempty" json:"AnchorStateRegistryProxy,omitempty"`
	DelayedWETHProxy                  *common.Address `toml:"DelayedWETHProxy,omitempty" json:"DelayedWETHProxy,omitempty"`
	DisputeGameFactoryProxy           *common.Address `toml:"DisputeGameFactoryProxy,omitempty" json:"DisputeGameFactoryProxy,omitempty"`
	FaultDisputeGame                  *common.Address `toml:"FaultDisputeGame,omitempty" json:"FaultDisputeGame,omitempty"`
	MIPS                              *common.Address `toml:"MIPS,omitempty" json:"MIPS,omitempty"`
	PermissionedDisputeGame           *common.Address `toml:"PermissionedDisputeGame,omitempty" json:"PermissionedDisputeGame,omitempty"`
	PreimageOracle                    *common.Address `toml:"PreimageOracle,omitempty" json:"PreimageOracle,omitempty"`
	DAChallengeAddress                *common.Address `toml:"DAChallengeAddress,omitempty" json:"DAChallengeAddress,omitempty"`
}

// Chain represents a registered chain in the superchain.
type Chain struct {
	Name    string `json:"name"`
	Network string `json:"network"`
	config  *ChainConfig
}

func (c *Chain) Config() (*ChainConfig, error) {
	if c.config == nil {
		return nil, fmt.Errorf("config not loaded for chain %s", c.Name)
	}
	return c.config, nil
}

// Chains is the global registry of known chains.
var Chains map[uint64]*Chain

// Superchains holds known superchain configs by network name.
var Superchains map[string]Superchain

func init() {
	Chains = make(map[uint64]*Chain)
	Superchains = make(map[string]Superchain)

	// Register OP Mainnet superchain
	Superchains["mainnet"] = Superchain{
		Name: "mainnet",
		L1:   L1Config{ChainID: 1},
	}
	// Register OP Sepolia superchain
	Superchains["sepolia"] = Superchain{
		Name: "sepolia",
		L1:   L1Config{ChainID: 11155111},
	}

	zero := uint64(0)
	denomCanyon := uint64(250)

	portalMainnet := common.HexToAddress("0xbEb5Fc579115071764c7423A4f12eDde41f106Ed")
	sysConfigMainnet := common.HexToAddress("0x229047fed2591dbec1eF1118d64F7aF3dB9EB290")
	portalSepolia := common.HexToAddress("0x16Fc5058F25648194471939df75CF27A2fdC48BC")
	sysConfigSepolia := common.HexToAddress("0x034edD2A225f7f429A63E0f1D2084B9E0A93b538")

	// OP Mainnet (chain ID 10)
	Chains[10] = &Chain{
		Name:    "op",
		Network: "mainnet",
		config: &ChainConfig{
			Name:              "OP Mainnet",
			ChainID:           10,
			BlockTime:         2,
			SeqWindowSize:     3600,
			MaxSequencerDrift: 600,
			BatchInboxAddr:    common.HexToAddress("0xFF00000000000000000000000000000000000010"),
			Optimism:          &OptimismConfig{EIP1559Elasticity: 6, EIP1559Denominator: 50, EIP1559DenominatorCanyon: &denomCanyon},
			Hardforks: HardforkConfig{
				CanyonTime:  &zero,
				DeltaTime:   &zero,
				EcotoneTime: &zero,
				FjordTime:   &zero,
				GraniteTime: &zero,
			},
			Addresses: AddressesConfig{
				OptimismPortalProxy:     &portalMainnet,
				SystemConfigProxy:       &sysConfigMainnet,
				DisputeGameFactoryProxy: addrPtr("0xe5965Ab5962eDc7477C8520243A95517CD252fA9"),
			},
			Genesis: GenesisConfig{
				L2Time: 1686068903,
				L1:     GenesisRef{Hash: common.HexToHash("0x438335a20d98863a4c0c97999eb2481921ccd28553eac6f913af7c12aec04108"), Number: 17422590},
				L2:     GenesisRef{Hash: common.HexToHash("0xdbf6a80fef073de06add9b0d14026d6e5a86c85f6d102c36d3d8e9cf89c2afd3"), Number: 105235063},
				SystemConfig: SystemConfig{
					BatcherAddr: common.HexToAddress("0x6887246668a3b87F54DeB3b94Ba47a6f63F32985"),
					Overhead:    common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc"),
					Scalar:      common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0"),
					GasLimit:    30000000,
				},
			},
			Network: "mainnet",
		},
	}

	// OP Sepolia (chain ID 11155420)
	Chains[11155420] = &Chain{
		Name:    "op",
		Network: "sepolia",
		config: &ChainConfig{
			Name:              "OP Sepolia",
			ChainID:           11155420,
			BlockTime:         2,
			SeqWindowSize:     3600,
			MaxSequencerDrift: 600,
			BatchInboxAddr:    common.HexToAddress("0xFF00000000000000000000000000000000042069"),
			Optimism:          &OptimismConfig{EIP1559Elasticity: 6, EIP1559Denominator: 50, EIP1559DenominatorCanyon: &denomCanyon},
			Hardforks: HardforkConfig{
				CanyonTime:  &zero,
				DeltaTime:   &zero,
				EcotoneTime: &zero,
				FjordTime:   &zero,
				GraniteTime: &zero,
			},
			Addresses: AddressesConfig{
				OptimismPortalProxy:     &portalSepolia,
				SystemConfigProxy:       &sysConfigSepolia,
				DisputeGameFactoryProxy: addrPtr("0x05F9613aDB30026FFd634f38e5C4dFd30a197Fa1"),
			},
			Genesis: GenesisConfig{
				L2Time: 1691802540,
				L1:     GenesisRef{Hash: common.HexToHash("0x48f520cf4ddaf34c8336e6e490571f5e0fd16de3f16b0d18e93c7e95299675d2"), Number: 4071248},
				L2:     GenesisRef{Hash: common.HexToHash("0x102de6ffb001480cc9b8b548fd05c34cd4f46ae4aa91759393db90ea0409887d"), Number: 0},
				SystemConfig: SystemConfig{
					BatcherAddr: common.HexToAddress("0x8F23BB38F531600e5d8FDDaAEC41F13FaB46E98c"),
					Overhead:    common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc"),
					Scalar:      common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0"),
					GasLimit:    30000000,
				},
			},
			Network: "sepolia",
		},
	}
}

// GetSuperchain returns the superchain config for a given network.
func GetSuperchain(network string) (Superchain, error) {
	sc, ok := Superchains[network]
	if !ok {
		return Superchain{}, fmt.Errorf("%w: superchain registry not available for network %q", ErrUnknownChain, network)
	}
	return sc, nil
}

// GetChain returns chain config by chain ID.
func GetChain(chainID uint64) (*Chain, error) {
	chain, ok := Chains[chainID]
	if !ok {
		return nil, fmt.Errorf("%w ID: %d", ErrUnknownChain, chainID)
	}
	return chain, nil
}

// GetDepset returns the dependency set for a chain.
func GetDepset(chainID uint64) (map[string]Dependency, error) {
	chain, err := GetChain(chainID)
	if err != nil {
		return nil, err
	}
	cfg, err := chain.Config()
	if err != nil {
		return nil, err
	}
	if cfg.Interop == nil {
		cfg.Interop = &Interop{Dependencies: make(map[string]Dependency)}
		cfg.Interop.Dependencies[fmt.Sprintf("%d", cfg.ChainID)] = Dependency{}
	}
	return cfg.Interop.Dependencies, nil
}

// ChainIDByName returns a chain ID by name.
func ChainIDByName(name string) (uint64, error) {
	for id, chain := range Chains {
		if chain.Name+"-"+chain.Network == name {
			return id, nil
		}
	}
	return 0, fmt.Errorf("%w %q", ErrUnknownChain, name)
}

// ChainNames returns all registered chain names.
func ChainNames() []string {
	var out []string
	for _, ch := range Chains {
		out = append(out, ch.Name+"-"+ch.Network)
	}
	return out
}
