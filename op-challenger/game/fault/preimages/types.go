package preimages

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/keccak/merkle"
	keccakTypes "github.com/tokamak-network/tokamak-thanos/op-challenger/game/keccak/types"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources/batching"
	"github.com/tokamak-network/tokamak-thanos/op-service/txmgr"
)

var ErrNilPreimageData = fmt.Errorf("cannot upload nil preimage data")

// PreimageUploader is responsible for posting preimages.
type PreimageUploader interface {
	// UploadPreimage uploads the provided preimage.
	UploadPreimage(ctx context.Context, claimIdx uint64, data *types.PreimageOracleData) error
}

// PreimageOracleContract is the interface for interacting with the PreimageOracle contract.
type PreimageOracleContract interface {
	InitLargePreimage(uuid *big.Int, partOffset uint32, claimedSize uint32) (txmgr.TxCandidate, error)
	AddLeaves(uuid *big.Int, startingBlockIndex *big.Int, input []byte, commitments []common.Hash, finalize bool) (txmgr.TxCandidate, error)
	Squeeze(claimant common.Address, uuid *big.Int, prestateMatrix keccakTypes.StateSnapshot, preState keccakTypes.Leaf, preStateProof merkle.Proof, postState keccakTypes.Leaf, postStateProof merkle.Proof) (txmgr.TxCandidate, error)
	CallSqueeze(ctx context.Context, claimant common.Address, uuid *big.Int, prestateMatrix keccakTypes.StateSnapshot, preState keccakTypes.Leaf, preStateProof merkle.Proof, postState keccakTypes.Leaf, postStateProof merkle.Proof) error
	GetProposalMetadata(ctx context.Context, block batching.Block, idents ...keccakTypes.LargePreimageIdent) ([]keccakTypes.LargePreimageMetaData, error)
	ChallengePeriod(ctx context.Context) (uint64, error)
	GetMinBondLPP(ctx context.Context) (*big.Int, error)
}
