package opcm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/script"
)

type DeployAlphabetVMInput struct {
	AbsolutePrestate common.Hash
	PreimageOracle   common.Address
}

type DeployAlphabetVMOutput struct {
	AlphabetVM common.Address
}

func DeployAlphabetVM(
	host *script.Host,
	input DeployAlphabetVMInput,
) (DeployAlphabetVMOutput, error) {
	return RunScriptSingle[DeployAlphabetVMInput, DeployAlphabetVMOutput](
		host,
		input,
		"DeployAlphabetVM.s.sol",
		"DeployAlphabetVM",
	)
}
