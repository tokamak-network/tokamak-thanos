package predeploys

import "github.com/ethereum/go-ethereum/common"

const (
	LegacyERC20NativeToken = "0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000"
)

var (
	LegacyERC20NativeTokenAddr = common.HexToAddress(LegacyERC20NativeToken)
)
