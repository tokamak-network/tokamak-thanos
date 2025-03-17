package fees

import (
	"context"
	"math/big"
	"testing"

	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"

	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/geth"
	"github.com/ethereum-optimism/optimism/op-e2e/system/e2esys"
	"github.com/ethereum-optimism/optimism/op-e2e/system/helpers"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	op_e2e.RunMain(m)
}

type stateGetterAdapter struct {
	ctx      context.Context
	t        *testing.T
	client   *ethclient.Client
	blockNum *big.Int
}

func (sga *stateGetterAdapter) GetState(addr common.Address, key common.Hash) common.Hash {
	sga.t.Helper()
	val, err := sga.client.StorageAt(sga.ctx, addr, key, sga.blockNum)
	require.NoError(sga.t, err)
	var res common.Hash
	copy(res[:], val)
	return res
}

// TestFees checks that L1/L2 fees are handled.
func TestFees(t *testing.T) {
	t.Run("pre-regolith", func(t *testing.T) {
		op_e2e.InitParallel(t)
		cfg := e2esys.RegolithSystemConfig(t, nil)
		cfg.DeployConfig.L1GenesisBlockBaseFeePerGas = (*hexutil.Big)(big.NewInt(7))

		testFees(t, cfg)
	})
	t.Run("regolith", func(t *testing.T) {
		op_e2e.InitParallel(t)
		cfg := e2esys.RegolithSystemConfig(t, new(hexutil.Uint64))
		cfg.DeployConfig.L1GenesisBlockBaseFeePerGas = (*hexutil.Big)(big.NewInt(7))

		testFees(t, cfg)
	})
	t.Run("ecotone", func(t *testing.T) {
		op_e2e.InitParallel(t)
		cfg := e2esys.EcotoneSystemConfig(t, new(hexutil.Uint64))
		cfg.DeployConfig.L1GenesisBlockBaseFeePerGas = (*hexutil.Big)(big.NewInt(7))

		testFees(t, cfg)
	})
	t.Run("fjord", func(t *testing.T) {
		op_e2e.InitParallel(t)
		cfg := e2esys.FjordSystemConfig(t, new(hexutil.Uint64))
		cfg.DeployConfig.L1GenesisBlockBaseFeePerGas = (*hexutil.Big)(big.NewInt(7))

		testFees(t, cfg)
	})

	t.Run("isthmus without operator fee", func(t *testing.T) {
		op_e2e.InitParallel(t)
		cfg := e2esys.IsthmusSystemConfig(t, new(hexutil.Uint64))
		cfg.DeployConfig.L1GenesisBlockBaseFeePerGas = (*hexutil.Big)(big.NewInt(7))

		testFees(t, cfg)
	})

	t.Run("isthmus with operator fee", func(t *testing.T) {
		op_e2e.InitParallel(t)
		cfg := e2esys.IsthmusSystemConfig(t, new(hexutil.Uint64))
		cfg.DeployConfig.L1GenesisBlockBaseFeePerGas = (*hexutil.Big)(big.NewInt(7))
		cfg.DeployConfig.GasPriceOracleOperatorFeeScalar = 1439103868
		cfg.DeployConfig.GasPriceOracleOperatorFeeConstant = 12564178266093314607

		testFees(t, cfg)
	})
}

func testFees(t *testing.T, cfg e2esys.SystemConfig) {
	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")

	l2Seq := sys.NodeClient("sequencer")
	l2Verif := sys.NodeClient("verifier")
	l1 := sys.NodeClient("l1")

	balanceAt := func(a common.Address, blockNumber *big.Int) *big.Int {
		t.Helper()
		bal, err := l2Seq.BalanceAt(context.Background(), a, blockNumber)
		require.NoError(t, err)
		return bal
	}

	// Wait for first block after genesis. The genesis block has zero L1Block values and will throw off the GPO checks
	_, err = geth.WaitForBlock(big.NewInt(1), l2Verif)
	require.NoError(t, err)

	config := sys.L2Genesis().Config

	sga := &stateGetterAdapter{
		ctx:    context.Background(),
		t:      t,
		client: l2Seq,
	}

	l1CostFn := types.NewL1CostFunc(config, sga)
	operatorFeeFn := types.NewOperatorCostFunc(config, sga)

	// Transactor Account
	ethPrivKey := cfg.Secrets.Alice
	fromAddr := crypto.PubkeyToAddress(ethPrivKey.PublicKey)

	require.NotEqual(t, cfg.DeployConfig.L2OutputOracleProposer, fromAddr)
	require.NotEqual(t, cfg.DeployConfig.BatchSenderAddress, fromAddr)

	// Find gaspriceoracle contract
	gpoContract, err := bindings.NewGasPriceOracle(predeploys.GasPriceOracleAddr, l2Seq)
	require.Nil(t, err)

	if !sys.RollupConfig.IsEcotone(sys.L2GenesisCfg.Timestamp) {
		overhead, err := gpoContract.Overhead(&bind.CallOpts{})
		require.Nil(t, err, "reading gpo overhead")
		require.Equal(t, overhead.Uint64(), cfg.DeployConfig.GasPriceOracleOverhead, "wrong gpo overhead")

		scalar, err := gpoContract.Scalar(&bind.CallOpts{})
		require.Nil(t, err, "reading gpo scalar")
		feeScalar := cfg.DeployConfig.FeeScalar()
		require.Equal(t, scalar, new(big.Int).SetBytes(feeScalar[:]), "wrong gpo scalar")
	} else {
		_, err := gpoContract.Overhead(&bind.CallOpts{})
		require.ErrorContains(t, err, "deprecated")
		_, err = gpoContract.Scalar(&bind.CallOpts{})
		require.ErrorContains(t, err, "deprecated")
	}

	decimals, err := gpoContract.Decimals(&bind.CallOpts{})
	require.Nil(t, err, "reading gpo decimals")

	require.Equal(t, decimals.Uint64(), uint64(6), "wrong gpo decimals")

	baseFeeRecipientStartBalance := balanceAt(predeploys.BaseFeeVaultAddr, big.NewInt(rpc.EarliestBlockNumber.Int64()))
	l1FeeRecipientStartBalance := balanceAt(predeploys.L1FeeVaultAddr, big.NewInt(rpc.EarliestBlockNumber.Int64()))
	sequencerFeeVaultStartBalance := balanceAt(predeploys.SequencerFeeVaultAddr, big.NewInt(rpc.EarliestBlockNumber.Int64()))
	operatorFeeVaultStartBalance := balanceAt(predeploys.OperatorFeeVaultAddr, big.NewInt(rpc.EarliestBlockNumber.Int64()))

	require.Equal(t, baseFeeRecipientStartBalance.Sign(), 0, "base fee vault should be empty")
	require.Equal(t, l1FeeRecipientStartBalance.Sign(), 0, "l1 fee vault should be empty")
	require.Equal(t, sequencerFeeVaultStartBalance.Sign(), 0, "sequencer fee vault should be empty")
	require.Equal(t, operatorFeeVaultStartBalance.Sign(), 0, "operator fee vault should be empty")

	genesisBlock, err := l2Seq.BlockByNumber(context.Background(), big.NewInt(rpc.EarliestBlockNumber.Int64()))
	require.NoError(t, err)

	coinbaseStartBalance := balanceAt(genesisBlock.Coinbase(), big.NewInt(rpc.EarliestBlockNumber.Int64()))

	// Simple transfer from signer to random account
	startBalance := balanceAt(fromAddr, big.NewInt(rpc.EarliestBlockNumber.Int64()))

	require.Greater(t, startBalance.Uint64(), big.NewInt(params.Ether).Uint64())

	transferAmount := big.NewInt(params.Ether)
	gasTip := big.NewInt(10)
	receipt := helpers.SendL2Tx(t, cfg, l2Seq, ethPrivKey, func(opts *helpers.TxOpts) {
		opts.ToAddr = &common.Address{0xff, 0xff}
		opts.Value = transferAmount
		opts.GasTipCap = gasTip
		opts.Gas = 21000
		opts.GasFeeCap = big.NewInt(200)
		opts.VerifyOnClients(l2Verif)
	})

	require.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)

	header, err := l2Seq.HeaderByNumber(context.Background(), receipt.BlockNumber)
	require.Nil(t, err)

	coinbaseEndBalance := balanceAt(header.Coinbase, header.Number)
	endBalance := balanceAt(fromAddr, header.Number)
	baseFeeRecipientEndBalance := balanceAt(predeploys.BaseFeeVaultAddr, header.Number)

	l1Header, err := l1.HeaderByNumber(context.Background(), nil)
	require.Nil(t, err)

	l1FeeRecipientEndBalance := balanceAt(predeploys.L1FeeVaultAddr, header.Number)
	sequencerFeeVaultEndBalance := balanceAt(predeploys.SequencerFeeVaultAddr, header.Number)
	operatorFeeVaultEndBalance := balanceAt(predeploys.OperatorFeeVaultAddr, header.Number)

	// Diff fee recipient + coinbase balances
	baseFeeRecipientDiff := new(big.Int).Sub(baseFeeRecipientEndBalance, baseFeeRecipientStartBalance)
	l1FeeRecipientDiff := new(big.Int).Sub(l1FeeRecipientEndBalance, l1FeeRecipientStartBalance)
	sequencerFeeVaultDiff := new(big.Int).Sub(sequencerFeeVaultEndBalance, sequencerFeeVaultStartBalance)
	operatorFeeVaultDiff := new(big.Int).Sub(operatorFeeVaultEndBalance, operatorFeeVaultStartBalance)
	coinbaseDiff := new(big.Int).Sub(coinbaseEndBalance, coinbaseStartBalance)

	// Tally L2 Fee
	gasUsed := new(big.Int).SetUint64(receipt.GasUsed)
	l2Fee := gasTip.Mul(gasTip, gasUsed)
	require.Equal(t, sequencerFeeVaultDiff, coinbaseDiff, "coinbase is always sequencer fee vault")
	require.Equal(t, l2Fee, coinbaseDiff, "l2 fee mismatch")
	require.Equal(t, l2Fee, sequencerFeeVaultDiff)

	// Tally BaseFee
	baseFee := new(big.Int).Mul(header.BaseFee, gasUsed)
	require.Equal(t, baseFee, baseFeeRecipientDiff, "base fee mismatch")

	// Tally L1 Fee
	tx, _, err := l2Seq.TransactionByHash(context.Background(), receipt.TxHash)
	require.NoError(t, err, "Should be able to get transaction")
	bytes, err := tx.MarshalBinary()
	require.Nil(t, err)

	l1Fee := l1CostFn(tx.RollupCostData(), header.Time)
	require.Equalf(t, l1Fee, l1FeeRecipientDiff, "L1 fee mismatch: start balance %v, end balance %v", l1FeeRecipientStartBalance, l1FeeRecipientEndBalance)

	operatorFee := operatorFeeFn(receipt.GasUsed, header.Time)
	require.True(t, operatorFeeVaultDiff.Cmp(operatorFee.ToBig()) == 0, "operator fee mismatch: start balance %v, end balance %v", operatorFeeVaultStartBalance, operatorFeeVaultEndBalance)

	gpoEcotone, err := gpoContract.IsEcotone(nil)
	require.NoError(t, err)
	require.Equal(t, sys.RollupConfig.IsEcotone(header.Time), gpoEcotone, "GPO and chain must have same ecotone view")

	gpoFjord, err := gpoContract.IsFjord(nil)
	require.NoError(t, err)
	require.Equal(t, sys.RollupConfig.IsFjord(header.Time), gpoFjord, "GPO and chain must have same fjord view")

	gpoIsthmus, err := gpoContract.IsIsthmus(nil)
	require.NoError(t, err)
	require.Equal(t, sys.RollupConfig.IsIsthmus(header.Time), gpoIsthmus, "GPO and chain must have same isthmus view")

	gpoL1Fee, err := gpoContract.GetL1Fee(&bind.CallOpts{}, bytes)
	require.Nil(t, err)

	adjustedGPOFee := gpoL1Fee
	if sys.RollupConfig.IsFjord(header.Time) {
		// The fastlz size of the transaction is 102 bytes
		require.Equal(t, uint64(102), tx.RollupCostData().FastLzSize)
		// Which results in both the fjord cost function and GPO using the minimum value for the fastlz regression:
		// Geth Linear Regression: -42.5856 + 102 * 0.8365 = 42.7374
		// GPO Linear Regression: -42.5856 + 170 * 0.8365 = 99.6194
		// The additional 68 (170 vs. 102) is due to the GPO adding 68 bytes to account for the signature.
		require.Greater(t, types.MinTransactionSize.Uint64(), uint64(99))
		// Because of this, we don't need to do any adjustment as the GPO and cost func are both bounded to the minimum value.
		// However, if the fastlz regression output is ever larger than the minimum, this will require an adjustment.
	} else if sys.RollupConfig.IsRegolith(header.Time) {
		// if post-regolith, adjust the GPO fee by removing the overhead it adds because of signature data
		artificialGPOOverhead := big.NewInt(68 * 16) // it adds 68 bytes to cover signature and RLP data
		l1BaseFee := big.NewInt(7)                   // we assume the L1 basefee is the minimum, 7
		// in our case we already include that, so we subtract it, to do a 1:1 comparison
		adjustedGPOFee = new(big.Int).Sub(gpoL1Fee, new(big.Int).Mul(artificialGPOOverhead, l1BaseFee))
	}
	require.Equal(t, l1Fee, adjustedGPOFee, "GPO reports L1 fee mismatch")

	require.Equal(t, receipt.L1Fee, l1Fee, "l1 fee in receipt is correct")
	if !sys.RollupConfig.IsEcotone(header.Time) { // FeeScalar receipt attribute is removed as of Ecotone
		require.Equal(t,
			new(big.Float).Mul(
				new(big.Float).SetInt(l1Header.BaseFee),
				new(big.Float).Mul(new(big.Float).SetInt(receipt.L1GasUsed), receipt.FeeScalar),
			),
			new(big.Float).SetInt(receipt.L1Fee), "fee field in receipt matches gas used times scalar times base fee")
	}

	expectedOperatorFee := new(big.Int).Add(
		new(big.Int).Div(
			new(big.Int).Mul(
				gasUsed,
				new(big.Int).SetUint64(uint64(cfg.DeployConfig.GasPriceOracleOperatorFeeScalar)),
			),
			new(big.Int).SetUint64(uint64(1e6)),
		),
		new(big.Int).SetUint64(cfg.DeployConfig.GasPriceOracleOperatorFeeConstant),
	)

	if sys.RollupConfig.IsIsthmus(header.Time) {
		require.True(t, expectedOperatorFee.Cmp(operatorFee.ToBig()) == 0,
			"operator fee is correct",
		)

		if cfg.DeployConfig.GasPriceOracleOperatorFeeScalar == 0 && cfg.DeployConfig.GasPriceOracleOperatorFeeConstant == 0 {
			// if both operator fee params are 0, they aren't added to receipts.
			require.Nil(t, receipt.OperatorFeeScalar, "operator fee constant in receipt is nil")
			require.Nil(t, receipt.OperatorFeeConstant, "operator fee constant in receipt is nil")
		} else {
			require.Equal(t, cfg.DeployConfig.GasPriceOracleOperatorFeeScalar, uint32(*receipt.OperatorFeeScalar), "operator fee constant in receipt is correct")
			require.Equal(t, cfg.DeployConfig.GasPriceOracleOperatorFeeConstant, *receipt.OperatorFeeConstant, "operator fee constant in receipt is correct")
		}
	}

	// Calculate total fee
	baseFeeRecipientDiff.Add(baseFeeRecipientDiff, coinbaseDiff)
	totalFee := new(big.Int).Add(baseFeeRecipientDiff, l1FeeRecipientDiff)
	totalFee = new(big.Int).Add(totalFee, operatorFeeVaultDiff)
	balanceDiff := new(big.Int).Sub(startBalance, endBalance)
	balanceDiff.Sub(balanceDiff, transferAmount)
	require.Equal(t, balanceDiff, totalFee, "balances should add up")
}
