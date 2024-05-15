package responder

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
)

// setupFaultDisputeGame deploys the FaultDisputeGame contract to a simulated backend
func setupFaultDisputeGame() (common.Address, *bind.TransactOpts, *backends.SimulatedBackend, *bindings.FaultDisputeGame, error) {
	privateKey, err := crypto.GenerateKey()
	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}
	backend := backends.NewSimulatedBackend(
		core.GenesisAlloc{from: {Balance: big.NewInt(params.Ether)}},
		50_000_000,
	)

	blockHashOracle, _, _, err := bindings.DeployBlockOracle(opts, backend)
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}

	_, _, contract, err := bindings.DeployFaultDisputeGame(
		opts,
		backend,
		uint8(0),             // Game Type ID
		[32]byte{0x01},       // Absolute Prestate Claim
		big.NewInt(15),       // Max Game Depth
		uint64(604800),       // 7 days
		common.Address{0xdd}, // VM
		common.Address{0xee}, // L2OutputOracle (Not used in Alphabet Game)
		blockHashOracle,      // Block hash oracle
	)
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}
	return from, opts, backend, contract, nil
}

// TestBuildFaultDefendData ensures that the manual ABI packing is the same as going through the bound contract.
func TestBuildFaultDefendData(t *testing.T) {
	_, opts, _, contract, err := setupFaultDisputeGame()
	require.NoError(t, err)

	resp, _ := newTestFaultResponder(t)

	data, err := resp.buildFaultDefendData(1, [32]byte{0x02, 0x03})
	require.NoError(t, err)

	opts.GasLimit = 100_000
	tx, err := contract.Defend(opts, big.NewInt(1), [32]byte{0x02, 0x03})
	require.NoError(t, err)

	require.Equal(t, data, tx.Data())
}

// TestBuildFaultAttackData ensures that the manual ABI packing is the same as going through the bound contract.
func TestBuildFaultAttackData(t *testing.T) {
	_, opts, _, contract, err := setupFaultDisputeGame()
	require.NoError(t, err)

	resp, _ := newTestFaultResponder(t)

	data, err := resp.buildFaultAttackData(1, [32]byte{0x02, 0x03})
	require.NoError(t, err)

	opts.GasLimit = 100_000
	tx, err := contract.Attack(opts, big.NewInt(1), [32]byte{0x02, 0x03})
	require.NoError(t, err)

	require.Equal(t, data, tx.Data())
}

// TestBuildFaultStepData ensures that the manual ABI packing is the same as going through the bound contract.
func TestBuildFaultStepData(t *testing.T) {
	_, opts, _, contract, err := setupFaultDisputeGame()
	require.NoError(t, err)

	resp, _ := newTestFaultResponder(t)

	data, err := resp.buildStepTxData(2, false, []byte{0x01}, []byte{0x02})
	require.NoError(t, err)

	opts.GasLimit = 100_000
	tx, err := contract.Step(opts, big.NewInt(2), false, []byte{0x01}, []byte{0x02})
	require.NoError(t, err)

	require.Equal(t, data, tx.Data())
}
