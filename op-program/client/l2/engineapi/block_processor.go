package engineapi

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/tokamak-network/tokamak-thanos/op-service/eip1559"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

var (
	ErrExceedsGasLimit  = errors.New("tx gas exceeds block gas limit")
	ErrUsesTooMuchGas   = errors.New("action takes too much gas")
	errInvalidGasLimit  = errors.New("invalid gas limit")
	errInvalidTimestamp = errors.New("invalid timestamp")
)

type BlockDataProvider interface {
	StateAt(root common.Hash) (*state.StateDB, error)
	GetHeader(common.Hash, uint64) *types.Header
	Engine() consensus.Engine
	GetVMConfig() *vm.Config
	Config() *params.ChainConfig
	consensus.ChainHeaderReader
}

type BlockProcessor struct {
	header       *types.Header
	state        *state.StateDB
	receipts     types.Receipts
	transactions types.Transactions
	gasPool      *core.GasPool
	dataProvider BlockDataProvider
	evm          *vm.EVM
}

func NewBlockProcessorFromPayloadAttributes(provider BlockDataProvider, parent common.Hash, attrs *eth.PayloadAttributes) (*BlockProcessor, error) {
	header := &types.Header{
		ParentHash:       parent,
		Coinbase:         attrs.SuggestedFeeRecipient,
		Difficulty:       common.Big0,
		GasLimit:         uint64(*attrs.GasLimit),
		Time:             uint64(attrs.Timestamp),
		Extra:            nil,
		MixDigest:        common.Hash(attrs.PrevRandao),
		Nonce:            types.EncodeNonce(0),
		ParentBeaconRoot: attrs.ParentBeaconBlockRoot,
	}
	if attrs.EIP1559Params != nil {
		d, e := eip1559.DecodeHolocene1559Params(attrs.EIP1559Params[:])
		if d == 0 {
			d = provider.Config().BaseFeeChangeDenominator(header.Time)
			e = provider.Config().ElasticityMultiplier()
		}
		// Encode extra data for Holocene-aware chains
		header.Extra = eip1559.EncodeHolocene1559Params(d, e)
	}

	return NewBlockProcessorFromHeader(provider, header)
}

func NewBlockProcessorFromHeader(provider BlockDataProvider, h *types.Header) (*BlockProcessor, error) {
	header := types.CopyHeader(h)

	if header.GasLimit > params.MaxGasLimit {
		return nil, fmt.Errorf("%w: have %v, max %v", errInvalidGasLimit, header.GasLimit, params.MaxGasLimit)
	}
	parentHeader := provider.GetHeaderByHash(header.ParentHash)
	if header.Time <= parentHeader.Time {
		return nil, errInvalidTimestamp
	}
	statedb, err := provider.StateAt(parentHeader.Root)
	if err != nil {
		return nil, fmt.Errorf("get parent state: %w", err)
	}
	header.Number = new(big.Int).Add(parentHeader.Number, common.Big1)
	header.BaseFee = eip1559.CalcBaseFee(provider.Config(), parentHeader, header.Time)
	header.GasUsed = 0
	gasPool := new(core.GasPool).AddGas(header.GasLimit)

	// Create block context and EVM
	context := core.NewEVMBlockContext(header, provider, nil, provider.Config(), statedb)
	txContext := vm.TxContext{}
	vmCfg := vm.Config{}
	if c := provider.GetVMConfig(); c != nil {
		vmCfg = *c
	}
	vmenv := vm.NewEVM(context, txContext, statedb, provider.Config(), vmCfg)

	if h.ParentBeaconRoot != nil {
		if provider.Config().IsCancun(header.Number, header.Time) {
			zero := uint64(0)
			header.BlobGasUsed = &zero
			header.ExcessBlobGas = &zero
		}
		// Re-create context after blob gas fields are set
		context = core.NewEVMBlockContext(header, provider, nil, provider.Config(), statedb)
		vmenv = vm.NewEVM(context, txContext, statedb, provider.Config(), vmCfg)
		core.ProcessBeaconBlockRoot(*header.ParentBeaconRoot, vmenv, statedb)
	}

	return &BlockProcessor{
		header:       header,
		state:        statedb,
		gasPool:      gasPool,
		dataProvider: provider,
		evm:          vmenv,
	}, nil
}

func (b *BlockProcessor) CheckTxWithinGasLimit(tx *types.Transaction) error {
	if tx.Gas() > b.header.GasLimit {
		return fmt.Errorf("%w tx gas: %d, block gas limit: %d", ErrExceedsGasLimit, tx.Gas(), b.header.GasLimit)
	}
	if tx.Gas() > b.gasPool.Gas() {
		return fmt.Errorf("%w: %d, only have %d", ErrUsesTooMuchGas, tx.Gas(), b.gasPool.Gas())
	}
	return nil
}

func (b *BlockProcessor) AddTx(tx *types.Transaction) (*types.Receipt, error) {
	txIndex := len(b.transactions)
	b.state.SetTxContext(tx.Hash(), txIndex)
	receipt, err := core.ApplyTransaction(
		b.dataProvider.Config(),
		b.dataProvider,
		nil, // author (coinbase already in header)
		b.gasPool,
		b.state,
		b.header,
		tx,
		&b.header.GasUsed,
		b.evm.Config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to apply transaction to L2 block (tx %d): %w", txIndex, err)
	}
	b.receipts = append(b.receipts, receipt)
	b.transactions = append(b.transactions, tx)
	return receipt, nil
}

func (b *BlockProcessor) Assemble() (*types.Block, types.Receipts, error) {
	var withdrawals []*types.Withdrawal
	if b.header.WithdrawalsHash != nil {
		withdrawals = []*types.Withdrawal{}
	}
	block, err := b.dataProvider.Engine().FinalizeAndAssemble(b.dataProvider, b.header, b.state, b.transactions, nil, b.receipts, withdrawals)
	if err != nil {
		return nil, nil, err
	}
	return block, b.receipts, nil
}

func (b *BlockProcessor) Commit() error {
	root, err := b.state.Commit(b.header.Number.Uint64(), b.dataProvider.Config().IsEIP158(b.header.Number))
	if err != nil {
		return fmt.Errorf("state write error: %w", err)
	}
	if err := b.state.Database().TrieDB().Commit(root, false); err != nil {
		return fmt.Errorf("trie write error: %w", err)
	}
	return nil
}
