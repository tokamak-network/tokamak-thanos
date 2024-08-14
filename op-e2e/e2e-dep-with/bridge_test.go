package e2e_dep_with

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/predeploys"
	e2e "github.com/tokamak-network/tokamak-thanos/op-e2e"
	"github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/receipts"
	"github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/transactions"
	"github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/wait"
	"github.com/tokamak-network/tokamak-thanos/op-node/rollup/derive"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"
)

// TestERC20BridgeDeposits tests the the L1StandardBridge bridge ERC20
// functionality.
func TestERC20BridgeDeposits(t *testing.T) {
	e2e.InitParallel(t)

	cfg := e2e.DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LevelInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	// Deploy wNativeToken
	wNativeTokenAddress, tx, wNativeToken, err := bindings.DeployWNativeToken(opts, l1Client)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err, "Waiting for deposit tx on L1")

	// Get some token
	opts.Value = big.NewInt(params.Ether)
	tx, err = wNativeToken.Fallback(opts, []byte{})
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	opts.Value = nil
	wNativeBalance, err := wNativeToken.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(params.Ether), wNativeBalance)

	l2Opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err)
	optimismMintableTokenFactory, err := bindings.NewOptimismMintableERC20Factory(predeploys.OptimismMintableERC20FactoryAddr, l2Client)
	require.NoError(t, err)
	tx, err = optimismMintableTokenFactory.CreateOptimismMintableERC20(l2Opts, wNativeTokenAddress, "L2-WETH", "L2-WETH")
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, tx.Hash())
	require.NoError(t, err)

	// Get the deployment event to have access to the L2 WNativeToken address
	it, err := optimismMintableTokenFactory.FilterOptimismMintableERC20Created(&bind.FilterOpts{Start: 0}, nil, nil)
	require.NoError(t, err)
	var event *bindings.OptimismMintableERC20FactoryOptimismMintableERC20Created
	for it.Next() {
		event = it.Event
	}
	require.NotNil(t, event)

	// Approve token with the bridge
	tx, err = wNativeToken.Approve(opts, cfg.L1Deployments.L1StandardBridgeProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Bridge the erc20 token
	l1StandardBridge, err := bindings.NewL1StandardBridge(cfg.L1Deployments.L1StandardBridgeProxy, l1Client)
	require.NoError(t, err)
	tx, err = transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1StandardBridge.BridgeERC20(opts, wNativeTokenAddress, event.LocalToken, big.NewInt(100), 100000, []byte{})
	})
	require.NoError(t, err)
	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	t.Log("Deposit through L1StandardBridge", "gas used", depositReceipt.GasUsed)

	// compute the deposit transaction hash + poll for it
	portal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	depositEvent, err := receipts.FindLog(depositReceipt.Logs, portal.ParseTransactionDeposited)
	require.NoError(t, err, "Should emit deposit event")

	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	// Ensure that the deposit went through
	optimismMintableToken, err := bindings.NewOptimismMintableERC20(event.LocalToken, l2Client)
	require.NoError(t, err)

	// Should have balance on L2
	l2Balance, err := optimismMintableToken.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l2Balance, big.NewInt(100))

	// Init L2StandardBridge contract
	l2StandardBridge, err := bindings.NewL2StandardBridge(predeploys.L2StandardBridgeAddr, l2Client)
	require.NoError(t, err)

	withdrawalTx, err := l2StandardBridge.BridgeERC20(l2Opts, event.LocalToken, event.RemoteToken, l2Balance, 200000, []byte{})
	require.NoError(t, err)
	withdrawalReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, withdrawalTx.Hash())
	require.NoError(t, err)

	l1BalanceBeforeFinalizingWithdraw, err := wNativeToken.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	provedReceipt, finalizedReceipt, _, _ := e2e.ProveAndFinalizeWithdrawal(t, cfg, sys, "sequencer", cfg.Secrets.Alice, withdrawalReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, provedReceipt.Status)
	require.Equal(t, types.ReceiptStatusSuccessful, finalizedReceipt.Status)

	l1BalanceAfterFinalizingWithdraw, err := wNativeToken.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	require.Equal(t, l1BalanceAfterFinalizingWithdraw, l1BalanceBeforeFinalizingWithdraw.Add(l1BalanceBeforeFinalizingWithdraw, l2Balance))
}

// TestETHBridgeDeposits tests the the L1StandardBridge bridge ETH
// functionality.
func TestETHBridgeDeposits(t *testing.T) {
	e2e.InitParallel(t)

	cfg := e2e.DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	l1StandardBridge, err := bindings.NewL1StandardBridge(cfg.L1Deployments.L1StandardBridgeProxy, l1Client)
	require.NoError(t, err)

	l1ETHBalanceBeforeDeposit, err := l1Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	opts.Value = big.NewInt(params.Ether)
	tx, err := transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1StandardBridge.BridgeETH(opts, opts.Value, 200000, []byte{})
	})
	require.NoError(t, err)
	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	t.Log("Deposit through L1StandardBridge", "gas used", depositReceipt.GasUsed)

	l1ETHBalanceAfterDeposit, err := l1Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	// Check balances
	usedETH := big.NewInt(params.Ether)
	// usedETH = usedETH.Add(usedETH, big.NewInt(int64(depositReceipt.GasUsed)))
	cost := big.NewInt(int64(depositReceipt.GasUsed))
	cost = cost.Mul(cost, depositReceipt.EffectiveGasPrice)
	usedETH = usedETH.Add(usedETH, cost)
	require.Equal(t, l1ETHBalanceBeforeDeposit, l1ETHBalanceAfterDeposit.Add(l1ETHBalanceAfterDeposit, usedETH))

	// compute the deposit transaction hash + poll for it
	portal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	depIt, err := portal.FilterTransactionDeposited(&bind.FilterOpts{Start: 0}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	l2ETH, err := bindings.NewETH(predeploys.ETHAddr, l2Client)
	require.NoError(t, err)

	l2ETHBalance, err := l2ETH.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l2ETHBalance, big.NewInt(params.Ether))

	// Init L2StandardBridge contract
	l2StandardBridge, err := bindings.NewL2StandardBridge(predeploys.L2StandardBridgeAddr, l2Client)
	require.NoError(t, err)

	l2Opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err)

	withdrawalTx, err := l2StandardBridge.Withdraw(l2Opts, predeploys.ETHAddr, big.NewInt(params.Ether), 200000, []byte{})
	require.NoError(t, err)
	withdrawalReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, withdrawalTx.Hash())
	require.NoError(t, err)

	l1ETHBalanceAfterDeposit, err = l1Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	provedReceipt, finalizedReceipt, _, _ := e2e.ProveAndFinalizeWithdrawal(t, cfg, sys, "sequencer", cfg.Secrets.Alice, withdrawalReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, provedReceipt.Status)
	require.Equal(t, types.ReceiptStatusSuccessful, finalizedReceipt.Status)

	l1ETHBalanceAfterFinalizingWithdraw, err := l1Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	usedETH = big.NewInt(params.Ether)

	provedCost := big.NewInt(int64(provedReceipt.GasUsed))
	provedCost.Mul(provedCost, provedReceipt.EffectiveGasPrice)
	usedETH.Sub(usedETH, provedCost)

	finalizedCost := big.NewInt(int64(finalizedReceipt.GasUsed))
	finalizedCost.Mul(finalizedCost, finalizedReceipt.EffectiveGasPrice)
	usedETH.Sub(usedETH, finalizedCost)

	require.Equal(t, l1ETHBalanceAfterFinalizingWithdraw, usedETH.Add(l1ETHBalanceAfterDeposit, usedETH))
}

// TestNativeTokenBridgeDeposits tests the the L1StandardBridge bridgeNativeToken
// functionality.
func TestNativeTokenBridgeDeposits(t *testing.T) {
	e2e.InitParallel(t)

	cfg := e2e.DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount int64 = 9

	opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve NativeToken with the bridge
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.L1StandardBridgeProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	l1BalanceBeforeDeposit, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	// Init L1StandardBridge contract
	l1StandardBridge, err := bindings.NewL1StandardBridge(cfg.L1Deployments.L1StandardBridgeProxy, l1Client)
	require.NoError(t, err)

	l2BalanceBeforeDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	// Deposit Native token
	tx, err = transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1StandardBridge.BridgeNativeToken(opts, big.NewInt(depositedAmount), 200000, []byte{})
	})
	require.NoError(t, err)
	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through L1StandardBridge", "gas used", depositReceipt.GasUsed)

	l1BalanceAfterDeposit, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l1BalanceBeforeDeposit, l1BalanceAfterDeposit.Add(l1BalanceAfterDeposit, big.NewInt(depositedAmount)))

	// compute the deposit transaction hash + poll for it
	portal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	depIt, err := portal.FilterTransactionDeposited(&bind.FilterOpts{Start: 0}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	l2BalanceAfterDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	l2BalanceExpectedAmount := l2BalanceBeforeDeposit.Add(l2BalanceBeforeDeposit, big.NewInt(depositedAmount))

	require.Equal(t, l2BalanceExpectedAmount, l2BalanceAfterDeposit)

	// Init L2StandardBridge contract
	l2StandardBridge, err := bindings.NewL2StandardBridge(predeploys.L2StandardBridgeAddr, l2Client)
	require.NoError(t, err)

	l2Opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err)

	l2Opts.Value = big.NewInt(depositedAmount)
	withdrawalTx, err := l2StandardBridge.WithdrawNativeToken(l2Opts, 200000, []byte{})
	require.NoError(t, err)
	withdrawalReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, withdrawalTx.Hash())
	require.NoError(t, err)

	provedReceipt, finalizedReceipt, _, _ := e2e.ProveAndFinalizeWithdrawal(t, cfg, sys, "sequencer", cfg.Secrets.Alice, withdrawalReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, provedReceipt.Status)
	require.Equal(t, types.ReceiptStatusSuccessful, finalizedReceipt.Status)

	l2BalanceAfterFinalizingWithdraw, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l2BalanceAfterFinalizingWithdraw, l1BalanceBeforeDeposit)
}

func TestWithdrawNativeTokenTo(t *testing.T) {
	e2e.InitParallel(t)

	cfg := e2e.DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	var depositedAmount int64 = 9

	opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	nativeTokenContract, err := bindings.NewL2NativeToken(cfg.L1Deployments.L2NativeToken, l1Client)
	require.NoError(t, err)

	// faucet NativeToken
	tx, err := nativeTokenContract.Faucet(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve NativeToken with the bridge
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.L1StandardBridgeProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	l1BalanceBeforeDeposit, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	// Init L1StandardBridge contract
	l1StandardBridge, err := bindings.NewL1StandardBridge(cfg.L1Deployments.L1StandardBridgeProxy, l1Client)
	require.NoError(t, err)

	l2BalanceBeforeDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	// Deposit Native token
	tx, err = transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1StandardBridge.BridgeNativeToken(opts, big.NewInt(depositedAmount), 200000, []byte{})
	})
	require.NoError(t, err)
	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through L1StandardBridge", "gas used", depositReceipt.GasUsed)

	l1BalanceAfterDeposit, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l1BalanceBeforeDeposit, l1BalanceAfterDeposit.Add(l1BalanceAfterDeposit, big.NewInt(depositedAmount)))

	// compute the deposit transaction hash + poll for it
	portal, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
	require.NoError(t, err)

	depIt, err := portal.FilterTransactionDeposited(&bind.FilterOpts{Start: 0}, nil, nil, nil)
	require.NoError(t, err)
	var depositEvent *bindings.OptimismPortalTransactionDeposited
	for depIt.Next() {
		depositEvent = depIt.Event
	}
	require.NotNil(t, depositEvent)

	depositTx, err := derive.UnmarshalDepositLogEvent(&depositEvent.Raw)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, types.NewTx(depositTx).Hash())
	require.NoError(t, err)

	l2BalanceAfterDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	l2BalanceExpectedAmount := l2BalanceBeforeDeposit.Add(l2BalanceBeforeDeposit, big.NewInt(depositedAmount))

	require.Equal(t, l2BalanceExpectedAmount, l2BalanceAfterDeposit)

	// Init L2StandardBridge contract
	l2StandardBridge, err := bindings.NewL2StandardBridge(predeploys.L2StandardBridgeAddr, l2Client)
	require.NoError(t, err)

	l2Opts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err)

	bobAddress := sys.Cfg.Secrets.Addresses().Bob
	l2BalanceBeforeFinalizingWithdraw, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, bobAddress)
	require.NoError(t, err)

	l2Opts.Value = big.NewInt(depositedAmount)
	withdrawalTx, err := l2StandardBridge.WithdrawNativeTokenTo(l2Opts, bobAddress, 200000, []byte{})
	require.NoError(t, err)
	withdrawalReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, withdrawalTx.Hash())
	require.NoError(t, err)

	provedReceipt, finalizedReceipt, _, _ := e2e.ProveAndFinalizeWithdrawal(t, cfg, sys, "sequencer", cfg.Secrets.Alice, withdrawalReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, provedReceipt.Status)
	require.Equal(t, types.ReceiptStatusSuccessful, finalizedReceipt.Status)

	l2BalanceAfterFinalizingWithdraw, err := nativeTokenContract.BalanceOf(&bind.CallOpts{}, bobAddress)
	require.NoError(t, err)
	require.Equal(t, l2BalanceAfterFinalizingWithdraw, l2BalanceBeforeFinalizingWithdraw.Add(l2BalanceBeforeFinalizingWithdraw, big.NewInt(depositedAmount)))
}
