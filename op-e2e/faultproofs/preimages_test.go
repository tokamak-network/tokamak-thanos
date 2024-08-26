package faultproofs

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	op_e2e "github.com/tokamak-network/tokamak-thanos/op-e2e"
	"github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/disputegame"
	preimage "github.com/tokamak-network/tokamak-thanos/op-preimage"
	"github.com/tokamak-network/tokamak-thanos/op-program/client"
)

func TestLocalPreimages(t *testing.T) {
	op_e2e.InitParallel(t, op_e2e.UsesCannon)
	tests := []struct {
		key preimage.Key
	}{
		{key: client.L1HeadLocalIndex},
		{key: client.L2OutputRootLocalIndex},
		{key: client.L2ClaimLocalIndex},
		{key: client.L2ClaimBlockNumberLocalIndex},
		// We don't check client.L2ChainIDLocalIndex because e2e tests use a custom chain configuration
		// which requires using a custom chain ID indicator so op-program will load the full rollup config and
		// genesis from the preimage oracle
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("preimage-%v", test.key), func(t *testing.T) {
			op_e2e.InitParallel(t, op_e2e.UsesCannon)

			ctx := context.Background()
			sys, _ := StartFaultDisputeSystem(t)
			t.Cleanup(sys.Close)

			disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys)
			game := disputeGameFactory.StartOutputCannonGame(ctx, "sequencer", 3, common.Hash{0x01, 0xaa})
			require.NotNil(t, game)
			claim := game.DisputeLastBlock(ctx)

			// Create the root of the cannon trace.
			claim = claim.Attack(ctx, common.Hash{0x01})

			game.LogGameData(ctx)

			game.VerifyPreimage(ctx, claim, test.key)

			game.LogGameData(ctx)
		})
	}
}
