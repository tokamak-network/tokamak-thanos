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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

// TestERC20BridgeDeposits tests the the L1StandardBridge bridge ERC20
// functionality.
func TestERC20BridgeDeposits(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	// Deploy WETH9
	weth9Address, tx, WETH9, err := bindings.DeployWTON(opts, l1Client)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err, "Waiting for deposit tx on L1")

	// Get some WETH
	opts.Value = big.NewInt(params.Ether)
	tx, err = WETH9.Fallback(opts, []byte{})
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	opts.Value = nil
	wethBalance, err := WETH9.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(params.Ether), wethBalance)

	// Deploy L2 WETH9
	l2Opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err)
	optimismMintableTokenFactory, err := bindings.NewOptimismMintableERC20Factory(predeploys.OptimismMintableERC20FactoryAddr, l2Client)
	require.NoError(t, err)
	tx, err = optimismMintableTokenFactory.CreateOptimismMintableERC20(l2Opts, weth9Address, "L2-WETH", "L2-WETH")
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, tx.Hash())
	require.NoError(t, err)

	// Get the deployment event to have access to the L2 WETH9 address
	it, err := optimismMintableTokenFactory.FilterOptimismMintableERC20Created(&bind.FilterOpts{Start: 0}, nil, nil)
	require.NoError(t, err)
	var event *bindings.OptimismMintableERC20FactoryOptimismMintableERC20Created
	for it.Next() {
		event = it.Event
	}
	require.NotNil(t, event)

	// Approve WETH9 with the bridge
	tx, err = WETH9.Approve(opts, cfg.L1Deployments.L1StandardBridgeProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Bridge the WETH9
	l1StandardBridge, err := bindings.NewL1StandardBridge(cfg.L1Deployments.L1StandardBridgeProxy, l1Client)
	require.NoError(t, err)
	tx, err = transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1StandardBridge.BridgeERC20(opts, weth9Address, event.LocalToken, big.NewInt(100), 100000, []byte{})
	})
	require.NoError(t, err)
	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	t.Log("Deposit through L1StandardBridge", "gas used", depositReceipt.GasUsed)

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

	// Ensure that the deposit went through
	optimismMintableToken, err := bindings.NewOptimismMintableERC20(event.LocalToken, l2Client)
	require.NoError(t, err)

	// Should have balance on L2
	l2Balance, err := optimismMintableToken.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l2Balance, big.NewInt(100))
}

// TestETHBridgeDeposits tests the the L1StandardBridge bridge ETH
// functionality.
func TestETHBridgeDeposits(t *testing.T) {
	InitParallel(t)

	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	log := testlog.Logger(t, log.LvlInfo)
	log.Info("genesis", "l2", sys.RollupConfig.Genesis.L2, "l1", sys.RollupConfig.Genesis.L1, "l2_time", sys.RollupConfig.Genesis.L2Time)

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	opts.Value = big.NewInt(params.Ether)
	l1StandardBridge, err := bindings.NewL1StandardBridge(cfg.L1Deployments.L1StandardBridgeProxy, l1Client)
	require.NoError(t, err)
	tx, err := transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1StandardBridge.BridgeETH(opts, 200000, []byte{})
	})
	require.NoError(t, err)
	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	t.Log("Deposit through L1StandardBridge", "gas used", depositReceipt.GasUsed)

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
}

// TestNativeTokenBridgeDeposits tests the the L1StandardBridge bridgeNativeToken
// functionality.
func TestNativeTokenBridgeDeposits(t *testing.T) {
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

	// Approve NativeToken with the bridge
	tx, err = nativeTokenContract.Approve(opts, cfg.L1Deployments.L1StandardBridgeProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Init L1StandardBridge contract
	l1StandardBridge, err := bindings.NewL1StandardBridge(cfg.L1Deployments.L1StandardBridgeProxy, l1Client)
	require.NoError(t, err)

	l2BalanceBeforeDeposit, err := l2Client.BalanceAt(context.Background(), opts.From, nil)
	require.NoError(t, err)

	// Deposit TON
	tx, err = transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1StandardBridge.BridgeNativeToken(opts, big.NewInt(depositedAmount), 200000, []byte{})
	})
	require.NoError(t, err)
	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through L1StandardBridge", "gas used", depositReceipt.GasUsed)

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
}
