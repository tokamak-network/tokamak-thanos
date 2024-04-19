package op_e2e

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/transactions"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

// TestDepositAndWithdrawPortalsSuccessfully tests deposit via OptimismPortal and withdraw via L2ToL1MessagePasser successfully
func TestDepositAndWithdrawSuccessfully(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount int64 = 9

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve NativeToken with the OP
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.OptimismPortalProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Init OptiismPortal contract
	optimismPortal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	l2BalanceBeforeDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	// Deposit NativeToken
	tx, err = optimismPortal.DepositTransaction(opts, opts.From, big.NewInt(depositedAmount), 200000, []byte{})
	require.NoError(t, err)

	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through OptiismPortal", "gas used", depositReceipt.GasUsed)

	depIt, err := optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: 0}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	l2BalanceAfterDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	l2BalanceExpectedAmount := l2BalanceBeforeDeposit.Add(l2BalanceBeforeDeposit, big.NewInt(depositedAmount))

	require.Equal(t, l2BalanceExpectedAmount, l2BalanceAfterDeposit)

	l2ToL1MessagePasser, err := bindings.NewL2ToL1MessagePasser(predeploys.L2ToL1MessagePasserAddr, l2Client)
	require.NoError(t, err)

	l2Opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err)
	l2Opts.Value = big.NewInt(depositedAmount)

	// init a withdraw
	tx, err = l2ToL1MessagePasser.InitiateWithdrawal(l2Opts, opts.From, big.NewInt(200000), []byte{})
	require.NoError(t, err)

	withdrawalReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, withdrawalReceipt.Status)

	l1BalanceBeforeFinalizeWithdraw, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	ProveAndFinalizeWithdrawal(t, cfg, l1Client, sys.EthInstances["sequencer"], cfg.Secrets.Alice, withdrawalReceipt)

	tx, err = nativeTokenContract.TransferFrom(opts, cfg.L1Deployments.OptimismPortalProxy, opts.From, big.NewInt(depositedAmount))
	require.NoError(t, err)

	withdrawTokenReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, withdrawTokenReceipt.Status)

	l1BalanceAfterFinalizeWithdraw, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l1BalanceAfterFinalizeWithdraw, l1BalanceBeforeFinalizeWithdraw.Add(l1BalanceBeforeFinalizeWithdraw, big.NewInt(depositedAmount)))
}

func TestOnApproveSuccessfully(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount = big.NewInt(9)

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, depositedAmount)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	calldata := EncodeCallData(opts.From, opts.From, depositedAmount, uint32(200000), []byte{})

	// Approve NativeToken with the OP
	tx, err = nativeTokenContract.ApproveAndCall(opts, cfg.L1Deployments.OptimismPortalProxy, depositedAmount, calldata)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Init OptimismPortal contract
	optimismPortal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	l2BalanceBeforeDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	depIt, err := optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: 0}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	l2BalanceAfterDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	l2BalanceExpectedAmount := l2BalanceBeforeDeposit.Add(l2BalanceBeforeDeposit, depositedAmount)

	require.Equal(t, l2BalanceExpectedAmount, l2BalanceAfterDeposit)
}

// TestDepositAndCallContractL2ViaPortal tests successfully
func TestDepositAndCallL2Successfully(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount int64 = 9

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve NativeToken with the OP
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.OptimismPortalProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Init OptiismPortal contract
	optimismPortal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	wnativetokenABI, err := bindings.WNativeTokenMetaData.GetAbi()
	require.NoError(t, err)
	calldata, err := wnativetokenABI.Pack("deposit")
	require.NoError(t, err)

	// Deposit NativeToken
	tx, err = optimismPortal.DepositTransaction(opts, predeploys.WNativeTokenAddr, big.NewInt(depositedAmount), 200000, calldata)
	require.NoError(t, err)

	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through OptiismPortal", "gas used", depositReceipt.GasUsed)

	depIt, err := optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: 0}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	wnativetokenContract, err := bindings.NewL2NativeToken(predeploys.WNativeTokenAddr, l2Client)
	require.NoError(t, err)
	wnativetokenBalance, err := wnativetokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, wnativetokenBalance, big.NewInt(depositedAmount))
}

func TestDeployContractFailedNonPayableConstructor(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount int64 = 9

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve NativeToken with the OP
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.OptimismPortalProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Init OptiismPortal contract
	optimismPortal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	// Deposit NativeToken
	tx, err = optimismPortal.DepositTransaction(opts, common.Address{}, big.NewInt(depositedAmount), 800000, []byte(bindings.WNativeTokenMetaData.Bin))
	require.NoError(t, err)

	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	depIt, err := optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: 0}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	relayedTxReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NotNil(t, err)
	require.Equal(t, types.ReceiptStatusFailed, relayedTxReceipt.Status)
}

func TestDeployContractFailedOutOfGas(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount int64 = 2000000000000000000

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve NativeToken with the OP
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.OptimismPortalProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Init OptiismPortal contract
	optimismPortal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	// Deposit NativeToken
	tx, err = optimismPortal.DepositTransaction(opts, opts.From, big.NewInt(depositedAmount), 200000, []byte{})
	require.NoError(t, err)

	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through OptiismPortal", "gas used", depositReceipt.GasUsed)

	depIt, err := optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: depositReceipt.BlockNumber.Uint64()}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)

	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	// Estimate data
	_, estimatedData, _, _ := bindings.DeployWNativeToken(opts, l1Client)

	// Deploy contract
	deployedTx, err := transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return optimismPortal.DepositTransaction(opts, common.Address{}, big.NewInt(0), 200000, estimatedData.Data())
	})
	require.NoError(t, err)

	depositDeployedReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, deployedTx.Hash())
	require.NoError(t, err)
	depIt, err = optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: depositDeployedReceipt.BlockNumber.Uint64()}, nil, nil, nil)
	require.NoError(t, err)
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err = derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	relayedTxReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NotNil(t, err)
	require.Equal(t, types.ReceiptStatusFailed, relayedTxReceipt.Status)
}

// TestDeployContractSuccessfully tests successfully
func TestDeployContractSuccessfully(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount int64 = 2000000000000000000

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve NativeToken with the OP
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.OptimismPortalProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Init OptiismPortal contract
	optimismPortal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	// Deposit NativeToken
	tx, err = optimismPortal.DepositTransaction(opts, opts.From, big.NewInt(depositedAmount), 200000, []byte{})
	require.NoError(t, err)

	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through OptiismPortal", "gas used", depositReceipt.GasUsed)

	depIt, err := optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: depositReceipt.BlockNumber.Uint64()}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)

	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	// Estimate data
	_, estimatedData, _, _ := bindings.DeployWNativeToken(opts, l1Client)
	t.Log("Estimated gas:", estimatedData.Gas())

	// Deploy contract
	deployedTx, err := transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return optimismPortal.DepositTransaction(opts, common.Address{}, big.NewInt(0), estimatedData.Gas(), estimatedData.Data())
	})
	require.NoError(t, err)

	depositDeployedReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, deployedTx.Hash())
	require.NoError(t, err)

	depIt, err = optimismPortal.FilterTransactionDeposited(&bind.FilterOpts{Start: depositDeployedReceipt.BlockNumber.Uint64()}, nil, nil, nil)
	require.NoError(t, err)
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	// Calculate relayed depositTx
	depositTx, err = derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	relayedTxReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, relayedTxReceipt.Status)
}
