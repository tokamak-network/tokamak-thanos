package upgrades

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-e2e/actions/helpers"
	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/geth"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	isthmusL1BlockCodeHash          = common.HexToHash("0xe59074b8d4c08924ce463087b05485b835650652528383a32ef009fb2b6d4050")
	isthmusGasPriceOracleCodeHash   = common.HexToHash("0x9279e9e0535a7b63939d670e7faec536256b63d4ff353eb521a3342c51ce26e5")
	isthmusOperatorFeeVaultCodeHash = common.HexToHash("0x9ee0fa5ab86f13f58fcb8f798f3a74401a8493d99d1c5b3bad19a8dff4b3194f")
)

func TestIsthmusActivationAtGenesis(gt *testing.T) {
	t := helpers.NewDefaultTesting(gt)
	env := helpers.SetupEnv(t, helpers.WithActiveGenesisFork(rollup.Isthmus))

	// Start op-nodes
	env.Seq.ActL2PipelineFull(t)
	env.Verifier.ActL2PipelineFull(t)

	// Verify Isthmus is active at genesis
	l2Head := env.Seq.L2Unsafe()
	require.NotZero(t, l2Head.Hash)
	require.True(t, env.SetupData.RollupCfg.IsIsthmus(l2Head.Time), "Isthmus should be active at genesis")

	// build empty L1 block
	env.Miner.ActEmptyBlock(t)

	// Build L2 chain and advance safe head
	env.Seq.ActL1HeadSignal(t)
	env.Seq.ActBuildToL1Head(t)

	// Make verifier sync, then check the block
	env.Verifier.ActL2PipelineFull(t)
	block := env.VerifEngine.L2Chain().CurrentBlock()
	verifyIsthmusHeaderWithdrawalsRoot(gt, env.SeqEngine.RPCClient(), block, true)
	require.Equal(t, types.EmptyRequestsHash, *block.RequestsHash, "isthmus block must have requests hash")

	// Check genesis config type can convert to a valid block
	genesisBlock, err := env.VerifEngine.EthClient().BlockByNumber(t.Ctx(), big.NewInt(0))
	require.NoError(t, err)
	reproduced := env.SetupData.L2Cfg.ToBlock()
	require.Equal(t, genesisBlock.WithdrawalsRoot(), reproduced.WithdrawalsRoot(), "genesis.ToBlock withdrawals-hash must match as expected")
	require.Equal(t, genesisBlock.Hash(), reproduced.Hash(), "genesis.ToBlock block hash must match")

	require.Equal(t, types.EmptyRequestsHash, *genesisBlock.RequestsHash(), "isthmus retrieved genesis block must have a requests-hash")
	require.Equal(t, types.EmptyRequestsHash, *reproduced.RequestsHash(), "isthmus generated genesis block have a requests-hash")

	// Check that the RPC client can handle block-hash verification of the genesis block
	cfg := sources.EngineClientDefaultConfig(env.SetupData.RollupCfg)
	cfg.TrustRPC = false // Make the RPC client check the block contents fully.
	l2Cl, err := sources.NewEngineClient(env.VerifEngine.RPCClient(), testlog.Logger(t, log.LevelInfo), nil, cfg)
	require.NoError(t, err)
	genesisPayload, err := l2Cl.PayloadByNumber(t.Ctx(), 0)
	require.NoError(t, err)
	require.Equal(t, genesisPayload.ExecutionPayload.WithdrawalsRoot, genesisBlock.WithdrawalsRoot())
	require.Equal(t, types.EmptyRequestsHash, *genesisPayload.RequestsHash)
	got, ok := genesisPayload.CheckBlockHash()
	require.Equal(t, got, reproduced.Hash())
	require.True(t, ok, "CheckBlockHash must pass")
}

// There are 2 stages pre-Isthmus that we need to test:
// 1. Pre-Canyon: withdrawals root should be nil
// 2. Post-Canyon: withdrawals root should be EmptyWithdrawalsHash
func TestWithdrawlsRootPreCanyonAndIsthmus(gt *testing.T) {
	t := helpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, helpers.DefaultRollupTestParams())
	genesisBlock := hexutil.Uint64(0)
	canyonOffset := hexutil.Uint64(2)

	log := testlog.Logger(t, log.LvlDebug)

	dp.DeployConfig.L1CancunTimeOffset = &canyonOffset

	// Activate pre-canyon forks at genesis, and schedule Canyon the block after
	dp.DeployConfig.L2GenesisRegolithTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisCanyonTimeOffset = &canyonOffset
	dp.DeployConfig.L2GenesisDeltaTimeOffset = nil
	dp.DeployConfig.L2GenesisEcotoneTimeOffset = nil
	dp.DeployConfig.L2GenesisFjordTimeOffset = nil
	dp.DeployConfig.L2GenesisGraniteTimeOffset = nil
	dp.DeployConfig.L2GenesisHoloceneTimeOffset = nil
	dp.DeployConfig.L2GenesisIsthmusTimeOffset = nil
	require.NoError(t, dp.DeployConfig.Check(log), "must have valid config")

	sd := e2eutils.Setup(t, dp, helpers.DefaultAlloc)
	_, _, _, sequencer, engine, verifier, _, _ := helpers.SetupReorgTestActors(t, dp, sd, log)

	// start op-nodes
	sequencer.ActL2PipelineFull(t)
	verifier.ActL2PipelineFull(t)

	verifyPreCanyonHeaderWithdrawalsRoot(gt, engine.L2Chain().CurrentBlock())

	// build blocks until canyon activates
	sequencer.ActBuildL2ToCanyon(t)

	// Send withdrawal transaction
	// Bind L2 Withdrawer Contract
	ethCl := engine.EthClient()
	l2withdrawer, err := bindings.NewL2ToL1MessagePasser(predeploys.L2ToL1MessagePasserAddr, ethCl)
	require.Nil(t, err, "binding withdrawer on L2")

	// Initiate Withdrawal
	l2opts, err := bind.NewKeyedTransactorWithChainID(dp.Secrets.Alice, new(big.Int).SetUint64(dp.DeployConfig.L2ChainID))
	require.Nil(t, err)
	l2opts.Value = big.NewInt(500)

	_, err = l2withdrawer.Receive(l2opts)
	require.Nil(t, err)

	// mine blocks
	sequencer.ActL2EmptyBlock(t)
	sequencer.ActL2EmptyBlock(t)

	verifyPreIsthmusHeaderWithdrawalsRoot(gt, engine.L2Chain().CurrentBlock())
}

// In this section, we will test the following combinations
// 1. Withdrawals root before isthmus w/ and w/o L2toL1 withdrawal
// 2. Withdrawals root at isthmus w/ and w/o L2toL1 withdrawal
// 3. Withdrawals root after isthmus w/ and w/o L2toL1 withdrawal
func TestWithdrawalsRootBeforeAtAndAfterIsthmus(t *testing.T) {
	tests := []struct {
		name              string
		f                 func(gt *testing.T, withdrawalTx bool, withdrawalTxBlock, totalBlocks int)
		withdrawalTx      bool
		withdrawalTxBlock int
		totalBlocks       int
	}{
		{"BeforeIsthmusWithoutWithdrawalTx", testWithdrawlsRootAtIsthmus, false, 0, 1},
		{"BeforeIsthmusWithWithdrawalTx", testWithdrawlsRootAtIsthmus, true, 1, 1},
		{"AtIsthmusWithoutWithdrawalTx", testWithdrawlsRootAtIsthmus, false, 0, 2},
		{"AtIsthmusWithWithdrawalTx", testWithdrawlsRootAtIsthmus, true, 2, 2},
		{"AfterIsthmusWithoutWithdrawalTx", testWithdrawlsRootAtIsthmus, false, 0, 3},
		{"AfterIsthmusWithWithdrawalTx", testWithdrawlsRootAtIsthmus, true, 3, 3},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			test.f(t, test.withdrawalTx, test.withdrawalTxBlock, test.totalBlocks)
		})
	}
}

func testWithdrawlsRootAtIsthmus(gt *testing.T, withdrawalTx bool, withdrawalTxBlock, totalBlocks int) {
	t := helpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, helpers.DefaultRollupTestParams())
	genesisBlock := hexutil.Uint64(0)
	isthmusOffset := hexutil.Uint64(2)

	log := testlog.Logger(t, log.LvlDebug)

	dp.DeployConfig.L2GenesisRegolithTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisCanyonTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisIsthmusTimeOffset = &isthmusOffset
	dp.DeployConfig.L2GenesisDeltaTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisEcotoneTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisFjordTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisGraniteTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisHoloceneTimeOffset = &genesisBlock
	require.NoError(t, dp.DeployConfig.Check(log), "must have valid config")

	sd := e2eutils.Setup(t, dp, helpers.DefaultAlloc)
	_, _, _, sequencer, engine, verifier, _, _ := helpers.SetupReorgTestActors(t, dp, sd, log)

	// start op-nodes
	sequencer.ActL2PipelineFull(t)
	verifier.ActL2PipelineFull(t)

	verifyPreIsthmusHeaderWithdrawalsRoot(gt, engine.L2Chain().CurrentBlock())

	ethCl := engine.EthClient()
	for i := 1; i <= totalBlocks; i++ {
		var tx *types.Transaction

		sequencer.ActL2StartBlock(t)

		if withdrawalTx && withdrawalTxBlock == i {
			l2withdrawer, err := bindings.NewL2ToL1MessagePasser(predeploys.L2ToL1MessagePasserAddr, ethCl)
			require.Nil(t, err, "binding withdrawer on L2")

			// Initiate Withdrawal
			// Bind L2 Withdrawer Contract and invoke the Receive function
			l2opts, err := bind.NewKeyedTransactorWithChainID(dp.Secrets.Alice, new(big.Int).SetUint64(dp.DeployConfig.L2ChainID))
			require.Nil(t, err)
			l2opts.Value = big.NewInt(500)
			tx, err = l2withdrawer.Receive(l2opts)
			require.Nil(t, err)

			// include the transaction
			engine.ActL2IncludeTx(dp.Addresses.Alice)(t)
		}
		sequencer.ActL2EndBlock(t)

		if withdrawalTx && withdrawalTxBlock == i {
			// wait for withdrawal to be included in a block
			receipt, err := geth.WaitForTransaction(tx.Hash(), ethCl, 10*time.Duration(dp.DeployConfig.L2BlockTime)*time.Second)
			require.Nil(t, err, "withdrawal initiated on L2 sequencer")
			require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "transaction had incorrect status")
		}
	}
	rpcCl := engine.RPCClient()

	// we set withdrawals root only at or after isthmus
	if totalBlocks >= 2 {
		verifyIsthmusHeaderWithdrawalsRoot(gt, rpcCl, engine.L2Chain().CurrentBlock(), true)
	}
}

func TestWithdrawlsRootPostIsthmus(gt *testing.T) {
	t := helpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, helpers.DefaultRollupTestParams())
	genesisBlock := hexutil.Uint64(0)
	isthmusOffset := hexutil.Uint64(2)

	log := testlog.Logger(t, log.LvlDebug)

	dp.DeployConfig.L2GenesisRegolithTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisCanyonTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisIsthmusTimeOffset = &isthmusOffset
	dp.DeployConfig.L2GenesisDeltaTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisEcotoneTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisFjordTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisGraniteTimeOffset = &genesisBlock
	dp.DeployConfig.L2GenesisHoloceneTimeOffset = &genesisBlock
	require.NoError(t, dp.DeployConfig.Check(log), "must have valid config")

	sd := e2eutils.Setup(t, dp, helpers.DefaultAlloc)
	_, _, _, sequencer, engine, verifier, _, _ := helpers.SetupReorgTestActors(t, dp, sd, log)

	// start op-nodes
	sequencer.ActL2PipelineFull(t)
	verifier.ActL2PipelineFull(t)

	verifyPreIsthmusHeaderWithdrawalsRoot(gt, engine.L2Chain().CurrentBlock())

	rpcCl := engine.RPCClient()
	verifyIsthmusHeaderWithdrawalsRoot(gt, rpcCl, engine.L2Chain().CurrentBlock(), false)

	// Send withdrawal transaction
	// Bind L2 Withdrawer Contract
	ethCl := engine.EthClient()
	l2withdrawer, err := bindings.NewL2ToL1MessagePasser(predeploys.L2ToL1MessagePasserAddr, ethCl)
	require.Nil(t, err, "binding withdrawer on L2")

	// Initiate Withdrawal
	l2opts, err := bind.NewKeyedTransactorWithChainID(dp.Secrets.Alice, new(big.Int).SetUint64(dp.DeployConfig.L2ChainID))
	require.Nil(t, err)
	l2opts.Value = big.NewInt(500)

	tx, err := l2withdrawer.Receive(l2opts)
	require.Nil(t, err)

	// build blocks until Isthmus activates
	sequencer.ActL2StartBlock(t)
	sequencer.ActL2EndBlock(t)
	sequencer.ActL2StartBlock(t)
	sequencer.ActL2EndBlock(t)
	sequencer.ActL2StartBlock(t)
	engine.ActL2IncludeTx(dp.Addresses.Alice)(t)
	sequencer.ActL2EndBlock(t)

	// wait for withdrawal to be included in a block
	receipt, err := geth.WaitForTransaction(tx.Hash(), ethCl, 10*time.Duration(dp.DeployConfig.L2BlockTime)*time.Second)
	require.Nil(t, err, "withdrawal initiated on L2 sequencer")
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "transaction had incorrect status")

	verifyIsthmusHeaderWithdrawalsRoot(gt, rpcCl, engine.L2Chain().CurrentBlock(), true)
}

// Pre-Canyon, the withdrawals root field in the header should be nil
func verifyPreCanyonHeaderWithdrawalsRoot(gt *testing.T, header *types.Header) {
	require.Nil(gt, header.WithdrawalsHash)
}

// Post-Canyon, the withdrawals root field in the header should be EmptyWithdrawalsHash
func verifyPreIsthmusHeaderWithdrawalsRoot(gt *testing.T, header *types.Header) {
	require.Equal(gt, types.EmptyWithdrawalsHash, *header.WithdrawalsHash)
}

func verifyIsthmusHeaderWithdrawalsRoot(gt *testing.T, rpcCl client.RPC, header *types.Header, l2toL1MPPresent bool) {
	getStorageRoot := func(rpcCl client.RPC, ctx context.Context, address common.Address, blockTag string) common.Hash {
		var getProofResponse *eth.AccountResult
		err := rpcCl.CallContext(ctx, &getProofResponse, "eth_getProof", address, []common.Hash{}, blockTag)
		assert.Nil(gt, err)
		assert.NotNil(gt, getProofResponse)
		return getProofResponse.StorageHash
	}

	if !l2toL1MPPresent {
		require.Equal(gt, types.EmptyWithdrawalsHash, *header.WithdrawalsHash)
	} else {
		storageHash := getStorageRoot(rpcCl, context.Background(), predeploys.L2ToL1MessagePasserAddr, "latest")
		require.Equal(gt, *header.WithdrawalsHash, storageHash)
	}
}

func TestIsthmusNetworkUpgradeTransactions(gt *testing.T) {
	t := helpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, helpers.DefaultRollupTestParams())
	isthmusOffset := hexutil.Uint64(4)

	log := testlog.Logger(t, log.LevelDebug)

	zero := hexutil.Uint64(0)

	// Activate all forks at genesis, and schedule Isthmus the block after
	dp.DeployConfig.L2GenesisHoloceneTimeOffset = &zero
	dp.DeployConfig.L2GenesisIsthmusTimeOffset = &isthmusOffset
	dp.DeployConfig.L1PragueTimeOffset = &zero
	// New forks have to be added here...
	require.NoError(t, dp.DeployConfig.Check(log), "must have valid config")

	sd := e2eutils.Setup(t, dp, helpers.DefaultAlloc)
	_, _, _, sequencer, engine, verifier, _, _ := helpers.SetupReorgTestActors(t, dp, sd, log)
	ethCl := engine.EthClient()

	// build a single block to move away from the genesis with 0-values in L1Block contract
	sequencer.ActL2StartBlock(t)
	sequencer.ActL2EndBlock(t)

	// start op-nodes
	sequencer.ActL2PipelineFull(t)
	verifier.ActL2PipelineFull(t)

	// Get gas price from oracle
	gasPriceOracle, err := bindings.NewGasPriceOracleCaller(predeploys.GasPriceOracleAddr, ethCl)
	require.NoError(t, err)

	// Get current implementations addresses (by slot) for L1Block + GasPriceOracle
	initialL1BlockAddress, err := ethCl.StorageAt(context.Background(), predeploys.L1BlockAddr, genesis.ImplementationSlot, nil)
	require.NoError(t, err)
	initialGasPriceOracleAddress, err := ethCl.StorageAt(context.Background(), predeploys.GasPriceOracleAddr, genesis.ImplementationSlot, nil)
	require.NoError(t, err)

	// Build to the isthmus block
	sequencer.ActBuildL2ToIsthmus(t)

	// get latest block
	latestBlock, err := ethCl.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, sequencer.L2Unsafe().Number, latestBlock.Number().Uint64())

	transactions := latestBlock.Transactions()

	// L1Block: 1 set-L1-info + 1 deploy
	// See [derive.IsthmusNetworkUpgradeTransactions]
	require.Equal(t, 9, len(transactions))

	// All transactions are successful
	for i := 1; i < 9; i++ {
		txn := transactions[i]
		receipt, err := ethCl.TransactionReceipt(context.Background(), txn.Hash())
		require.NoError(t, err)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
		require.NotEmpty(t, txn.Data(), "upgrade tx must provide input data")
	}

	expectedL1BlockAddress := crypto.CreateAddress(derive.L1BlockIsthmusDeployerAddress, 0)

	// L1 Block Proxy is updated
	updatedL1BlockAddress, err := ethCl.StorageAt(context.Background(), predeploys.L1BlockAddr, genesis.ImplementationSlot, latestBlock.Number())
	require.NoError(t, err)
	require.Equal(t, expectedL1BlockAddress, common.BytesToAddress(updatedL1BlockAddress))
	require.NotEqualf(t, initialL1BlockAddress, updatedL1BlockAddress, "Gas L1 Block address should have changed")
	verifyCodeHashMatches(t, ethCl, expectedL1BlockAddress, isthmusL1BlockCodeHash)

	expectedGasPriceOracleAddress := crypto.CreateAddress(derive.GasPriceOracleIsthmusDeployerAddress, 0)

	// Gas Price Oracle Proxy is updated
	updatedGasPriceOracleAddress, err := ethCl.StorageAt(context.Background(), predeploys.GasPriceOracleAddr, genesis.ImplementationSlot, latestBlock.Number())
	require.NoError(t, err)
	require.Equal(t, expectedGasPriceOracleAddress, common.BytesToAddress(updatedGasPriceOracleAddress))
	require.NotEqualf(t, initialGasPriceOracleAddress, updatedGasPriceOracleAddress, "Gas Price Oracle Proxy address should have changed")
	verifyCodeHashMatches(t, ethCl, expectedGasPriceOracleAddress, isthmusGasPriceOracleCodeHash)

	// Check that Isthmus was activated
	isIsthmus, err := gasPriceOracle.IsIsthmus(nil)
	require.NoError(t, err)
	require.True(t, isIsthmus)

	expectedOperatorFeeVaultAddress := crypto.CreateAddress(derive.OperatorFeeVaultDeployerAddress, 0)

	// Operator Fee vault is updated
	updatedOperatorFeeVaultAddress, err := ethCl.StorageAt(context.Background(), predeploys.OperatorFeeVaultAddr, genesis.ImplementationSlot, latestBlock.Number())
	require.NoError(t, err)
	require.Equal(t, expectedOperatorFeeVaultAddress, common.BytesToAddress(updatedOperatorFeeVaultAddress))
	verifyCodeHashMatches(t, ethCl, expectedOperatorFeeVaultAddress, isthmusOperatorFeeVaultCodeHash)

	// EIP-2935 contract is deployed
	expectedBlockHashAddress := crypto.CreateAddress(predeploys.EIP2935ContractDeployer, 0)
	require.Equal(t, predeploys.EIP2935ContractAddr, expectedBlockHashAddress)
	code := verifyCodeHashMatches(t, ethCl, predeploys.EIP2935ContractAddr, predeploys.EIP2935ContractCodeHash)
	require.Equal(t, predeploys.EIP2935ContractCode, code)

	// Test that the beacon-block-root has been set
	checkRecentBlockHash := func(blockNumber uint64, expectedHash common.Hash, msg string) {
		historyBufferLength := uint64(8191)
		bufferIdx := common.BigToHash(new(big.Int).SetUint64(blockNumber % historyBufferLength))

		rootValue, err := ethCl.StorageAt(context.Background(), predeploys.EIP2935ContractAddr, bufferIdx, nil)
		require.NoError(t, err)
		require.Equal(t, expectedHash, common.BytesToHash(rootValue), msg)
	}

	// Legacy check:
	// > The first block is an exception in upgrade-networks,
	// > since the recent-block-hash contract isn't there at Isthmus activation,
	// > and the recent-block-hash insertion is processed at the start of the block before deposit txs.
	// > If the contract was permissionlessly deployed before, the contract storage will be updated (but not in this test).
	checkRecentBlockHash(latestBlock.NumberU64(), common.Hash{}, "isthmus activation block has no data yet (since contract wasn't there)")

	// Build empty L2 block, to pass Isthmus activation
	sequencer.ActL2StartBlock(t)
	sequencer.ActL2EndBlock(t)

	// Test the L2 block after activation: it should have the most recent block hash
	latestBlock, err = ethCl.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)
	checkRecentBlockHash(latestBlock.NumberU64()-1, latestBlock.Header().ParentHash, "post-activation")
}
