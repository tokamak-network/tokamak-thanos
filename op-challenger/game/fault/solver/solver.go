package solver

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
)

var (
	ErrStepNonLeafNode       = errors.New("cannot step on non-leaf claims")
	ErrStepAgreedClaim       = errors.New("cannot step on claims we agree with")
	ErrStepIgnoreInvalidPath = errors.New("cannot step on claims that dispute invalid paths")
)

// claimSolver uses a [TraceProvider] to determine the moves to make in a dispute game.
type claimSolver struct {
	trace     types.TraceProvider
	gameDepth int
}

// newClaimSolver creates a new [claimSolver] using the provided [TraceProvider].
func newClaimSolver(gameDepth int, traceProvider types.TraceProvider) *claimSolver {
	return &claimSolver{
		traceProvider,
		gameDepth,
	}
}

// NextMove returns the next move to make given the current state of the game.
func (s *claimSolver) NextMove(ctx context.Context, claim types.Claim, game types.Game) (*types.Claim, error) {
	if claim.Depth() == s.gameDepth {
		return nil, types.ErrGameDepthReached
	}

	// Before challenging this claim, first check that the move wasn't warranted.
	// If the parent claim is on a dishonest path, then we would have moved against it anyways. So we don't move.
	// Avoiding dishonest paths ensures that there's always a valid claim available to support ours during step.
	if !claim.IsRoot() {
		parent, err := game.GetParent(claim)
		if err != nil {
			return nil, err
		}
		agreeWithParent, err := s.agreeWithClaimPath(ctx, game, parent)
		if err != nil {
			return nil, err
		}
		if !agreeWithParent {
			return nil, nil
		}
	}

	agree, err := s.agreeWithClaim(ctx, claim.ClaimData)
	if err != nil {
		return nil, err
	}
	if agree {
		return s.defend(ctx, claim)
	} else {
		return s.attack(ctx, claim)
	}
}

type StepData struct {
	LeafClaim  types.Claim
	IsAttack   bool
	PreState   []byte
	ProofData  []byte
	OracleData *types.PreimageOracleData
}

// AttemptStep determines what step should occur for a given leaf claim.
// An error will be returned if the claim is not at the max depth.
// Returns ErrStepIgnoreInvalidPath if the claim disputes an invalid path
func (s *claimSolver) AttemptStep(ctx context.Context, game types.Game, claim types.Claim) (StepData, error) {
	if claim.Depth() != s.gameDepth {
		return StepData{}, ErrStepNonLeafNode
	}

	// Step only on claims that dispute a valid path
	parent, err := game.GetParent(claim)
	if err != nil {
		return StepData{}, err
	}
	parentValid, err := s.agreeWithClaimPath(ctx, game, parent)
	if err != nil {
		return StepData{}, err
	}
	if !parentValid {
		return StepData{}, ErrStepIgnoreInvalidPath
	}

	claimCorrect, err := s.agreeWithClaim(ctx, claim.ClaimData)
	if err != nil {
		return StepData{}, err
	}
	var preState []byte
	var proofData []byte
	var oracleData *types.PreimageOracleData

	if !claimCorrect {
		// Attack the claim by executing step index, so we need to get the pre-state of that index
		preState, proofData, oracleData, err = s.trace.GetStepData(ctx, claim.Position)
		if err != nil {
			return StepData{}, err
		}
	} else {
		// We agree with the claim so Defend and use this claim as the starting point to
		// execute the step after. Thus we need the pre-state of the next step.
		preState, proofData, oracleData, err = s.trace.GetStepData(ctx, claim.MoveRight())
		if err != nil {
			return StepData{}, err
		}
	}

	return StepData{
		LeafClaim:  claim,
		IsAttack:   !claimCorrect,
		PreState:   preState,
		ProofData:  proofData,
		OracleData: oracleData,
	}, nil
}

// attack returns a response that attacks the claim.
func (s *claimSolver) attack(ctx context.Context, claim types.Claim) (*types.Claim, error) {
	position := claim.Attack()
	value, err := s.traceAtPosition(ctx, position)
	if err != nil {
		return nil, fmt.Errorf("attack claim: %w", err)
	}
	return &types.Claim{
		ClaimData:           types.ClaimData{Value: value, Position: position},
		ParentContractIndex: claim.ContractIndex,
	}, nil
}

// defend returns a response that defends the claim.
func (s *claimSolver) defend(ctx context.Context, claim types.Claim) (*types.Claim, error) {
	if claim.IsRoot() {
		return nil, nil
	}
	position := claim.Defend()
	value, err := s.traceAtPosition(ctx, position)
	if err != nil {
		return nil, fmt.Errorf("defend claim: %w", err)
	}
	return &types.Claim{
		ClaimData:           types.ClaimData{Value: value, Position: position},
		ParentContractIndex: claim.ContractIndex,
	}, nil
}

// agreeWithClaim returns true if the claim is correct according to the internal [TraceProvider].
func (s *claimSolver) agreeWithClaim(ctx context.Context, claim types.ClaimData) (bool, error) {
	ourValue, err := s.traceAtPosition(ctx, claim.Position)
	return bytes.Equal(ourValue[:], claim.Value[:]), err
}

// traceAtPosition returns the [common.Hash] from internal [TraceProvider] at the given [Position].
func (s *claimSolver) traceAtPosition(ctx context.Context, p types.Position) (common.Hash, error) {
	return s.trace.Get(ctx, p)
}

// agreeWithClaimPath returns true if the every other claim in the path to root is correct according to the internal [TraceProvider].
func (s *claimSolver) agreeWithClaimPath(ctx context.Context, game types.Game, claim types.Claim) (bool, error) {
	agree, err := s.agreeWithClaim(ctx, claim.ClaimData)
	if err != nil {
		return false, err
	}
	if !agree {
		return false, nil
	}
	if claim.IsRoot() {
		return true, nil
	}
	parent, err := game.GetParent(claim)
	if err != nil {
		return false, fmt.Errorf("failed to get parent of claim %v: %w", claim.ContractIndex, err)
	}
	if parent.IsRoot() {
		return true, nil
	}
	grandParent, err := game.GetParent(parent)
	if err != nil {
		return false, err
	}
	return s.agreeWithClaimPath(ctx, game, grandParent)
}
