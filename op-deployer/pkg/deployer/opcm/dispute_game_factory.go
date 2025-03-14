package opcm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/script"
)

type SetDisputeGameImplInput struct {
	Factory             common.Address
	Impl                common.Address
	Portal              common.Address
	AnchorStateRegistry common.Address
	GameType            uint32
}

func SetDisputeGameImpl(
	h *script.Host,
	input SetDisputeGameImplInput,
) error {
	return RunScriptVoid[SetDisputeGameImplInput](
		h,
		input,
		"SetDisputeGameImpl.s.sol",
		"SetDisputeGameImpl",
	)
}
