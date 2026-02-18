package p2p

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SigningDomainBlocksV1 is the domain used for signing blocks in v1.
var SigningDomainBlocksV1 = [32]byte{1}

// SigningHash computes the hash used for p2p block signing.
func SigningHash(domain [32]byte, chainID *big.Int, payloadHash []byte) (common.Hash, error) {
	if chainID.BitLen() > 256 {
		return common.Hash{}, fmt.Errorf("chain_id is too large: %d bits", chainID.BitLen())
	}
	var buf []byte
	buf = append(buf, domain[:]...)
	buf = append(buf, chainID.Bytes()...)
	buf = append(buf, payloadHash...)
	return crypto.Keccak256Hash(buf), nil
}
