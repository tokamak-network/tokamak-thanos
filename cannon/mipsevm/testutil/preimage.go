package testutil

import (
	"github.com/ethereum/go-ethereum/crypto"

	preimage "github.com/tokamak-network/tokamak-thanos/op-preimage"
)

func Keccak256Preimage(data []byte) [32]byte {
	return preimage.Keccak256Key(crypto.Keccak256Hash(data)).PreimageKey()
}
