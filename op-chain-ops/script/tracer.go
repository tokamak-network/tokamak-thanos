package script

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// scriptTracer implements vm.EVMLogger for the script Host.
// It delegates to the Host's hook methods.
type scriptTracer struct {
	host *Host
}

var _ vm.EVMLogger = (*scriptTracer)(nil)

func (t *scriptTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	var typ byte
	if create {
		typ = byte(vm.CREATE)
	} else {
		typ = byte(vm.CALL)
	}
	t.host.onEnter(0, typ, from, to, input, gas, value)
}

func (t *scriptTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
	t.host.onExit(0, output, gasUsed, err, false)
}

func (t *scriptTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	t.host.onEnter(1, byte(typ), from, to, input, gas, value)
}

func (t *scriptTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
	t.host.onExit(1, output, gasUsed, err, false)
}

func (t *scriptTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	// onOpcode and onStorageChange/onLog are handled here
	// We can't easily call onOpcode with tracing.OpContext since old geth uses ScopeContext
}

func (t *scriptTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

func (t *scriptTracer) CaptureTxStart(gasLimit uint64) {}
func (t *scriptTracer) CaptureTxEnd(restGas uint64)    {}
