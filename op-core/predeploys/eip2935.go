package predeploys

import (
	"github.com/ethereum/go-ethereum/common"
)

// EIP-2935 defines a deterministic deployment transaction that deploys the recent block hashes contract.
// See https://eips.ethereum.org/EIPS/eip-2935
var (
	EIP2935ContractAddr     = common.HexToAddress("0x0F792be4B0c0cb4DAE440Ef133E90C0eCD48CCCC")
	EIP2935ContractCode     = common.FromHex("3373fffffffffffffffffffffffffffffffffffffffe14604457602036146024575f5ffd5b620180005f350680545f35146037575f5ffd5b6201800001545f5260205ff35b42620180004206555f3562018000420662018000015500")
	EIP2935ContractCodeHash = common.HexToHash("0x6e49e66782037c0555897870e29fa5e552daf4719552131a0abce779daec0a5d")
	EIP2935ContractDeployer = common.HexToAddress("0x3462413Af4609098e1E27A490f554f260213D685")
)
