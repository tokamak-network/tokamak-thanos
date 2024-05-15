package op_e2e

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/predeploys"
	"github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/transactions"
	"github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/wait"
	"github.com/tokamak-network/tokamak-thanos/op-node/rollup/derive"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"
)

func TestUsdcDepositsAndWithdrawal(t *testing.T) {
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
	var withdrawalAmount int64 = 8

	bobAddress := sys.cfg.Secrets.Addresses().Bob

	opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L1ChainIDBig())
	require.Nil(t, err)

	// Deploy FiatToken and FiatTokenproxy in L1
	FiatTokenAddr, tx, _, err := bindings.DeployFiatTokenV22(opts, l1Client)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err, "Waiting for deposit tx on L1")

	FiatTokenProxyAddr, tx, FiatTokenProxyContract, err := bindings.DeployFiatTokenProxy(opts, l1Client, FiatTokenAddr)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err, "Waiting for deposit tx on L1")

	tx, err = FiatTokenProxyContract.ChangeAdmin(opts, bobAddress)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Deploy MasterMinter in L1
	MasterMinterAddr, tx, MasterMinterContract, err := bindings.DeployMasterMinter(opts, l1Client, FiatTokenProxyAddr)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err, "Waiting for deposit tx on L1")

	// FiatToken L1
	FiatTokenContract, err := bindings.NewFiatTokenV22(FiatTokenProxyAddr, l1Client)
	require.NoError(t, err)

	// initialize FiatToken L1
	tx, err = FiatTokenContract.Initialize(opts, "Bridged USDC (Tokamak Network)", "USDC.e", "USD", 6, MasterMinterAddr, opts.From, opts.From, opts.From)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// FiatToken L2
	FiatTokenContractL2, err := bindings.NewFiatTokenV22(predeploys.FiatTokenV2_2Addr, l2Client)
	require.NoError(t, err)

	// set - MasterMinter ConfigureController
	tx, err = MasterMinterContract.ConfigureController(opts, opts.From, opts.From)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// set - MasterMinter ConfigureMinter
	tx, err = MasterMinterContract.ConfigureMinter(opts, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Mint
	tx, err = FiatTokenContract.Mint(opts, opts.From, big.NewInt(depositedAmount))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// Approve
	tx, err = FiatTokenContract.Approve(opts, cfg.L1Deployments.L1UsdcBridgeProxy, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)

	// balance check (l1BalanceBeforeDeposit)
	_, err = FiatTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	// Init L1UsdcBridge contract
	l1UsdcBridge, err := bindings.NewL1UsdcBridge(cfg.L1Deployments.L1UsdcBridgeProxy, l1Client)
	require.NoError(t, err)

	// balance check
	l2BalanceBeforeDeposit, err := FiatTokenContractL2.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	// Deposit FiatToken
	tx, err = transactions.PadGasEstimate(opts, 1.1, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return l1UsdcBridge.DepositERC20(opts, FiatTokenProxyAddr, predeploys.FiatTokenV2_2Addr, big.NewInt(depositedAmount), 200000, []byte{})
	})
	require.NoError(t, err)

	depositReceipt, err := wait.ForReceiptOK(context.Background(), l1Client, tx.Hash())
	require.NoError(t, err)
	t.Log("Deposit through L1UsdcBridge", "gas used", depositReceipt.GasUsed)

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

	// balance check
	l2BalanceAfterDeposit, err := FiatTokenContractL2.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	l1BalanceAfterDeposit, err := FiatTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	l2BalanceExpectedAmount := l2BalanceBeforeDeposit.Add(l2BalanceBeforeDeposit, big.NewInt(depositedAmount))
	require.Equal(t, l2BalanceExpectedAmount, l2BalanceAfterDeposit)

	// Init L2UsdcBridge contract
	l2UsdcBridge, err := bindings.NewL2UsdcBridge(predeploys.L2UsdcBridgeAddr, l2Client)
	require.NoError(t, err)

	l2Opts, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err)

	// Approve
	tx, err = FiatTokenContractL2.Approve(l2Opts, predeploys.L2UsdcBridgeAddr, new(big.Int).SetUint64(math.MaxUint64))
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), l2Client, tx.Hash())
	require.NoError(t, err)

	// withdraw
	withdrawalTx, err := l2UsdcBridge.Withdraw(l2Opts, predeploys.FiatTokenV2_2Addr, big.NewInt(withdrawalAmount), 200000, []byte{})
	require.NoError(t, err)
	withdrawalReceipt, err := wait.ForReceiptOK(context.Background(), l2Client, withdrawalTx.Hash())
	require.NoError(t, err)

	// balance check
	l2BalanceAfterWithdraw, err := FiatTokenContractL2.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	provedReceipt, finalizedReceipt := ProveAndFinalizeWithdrawal(t, cfg, l1Client, sys.EthInstances["sequencer"], cfg.Secrets.Alice, withdrawalReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, provedReceipt.Status)
	require.Equal(t, types.ReceiptStatusSuccessful, finalizedReceipt.Status)

	// balance check
	l1BalanceAfterFinalizingWithdraw, err := FiatTokenContract.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)

	l1BalanceExpectedAmount := l1BalanceAfterDeposit.Add(l1BalanceAfterDeposit, big.NewInt(withdrawalAmount))
	require.Equal(t, l1BalanceExpectedAmount, l1BalanceAfterFinalizingWithdraw)

	l2BalanceAfterFinalizingWithdraw, err := FiatTokenContractL2.BalanceOf(&bind.CallOpts{}, opts.From)
	require.NoError(t, err)
	require.Equal(t, l2BalanceAfterFinalizingWithdraw, l2BalanceAfterWithdraw)
	require.Equal(t, l2BalanceAfterFinalizingWithdraw, big.NewInt(1))
}
