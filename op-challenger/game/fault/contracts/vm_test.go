package contracts

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	preimage "github.com/tokamak-network/tokamak-thanos/op-preimage"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources/batching"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources/batching/rpcblock"
	batchingTest "github.com/tokamak-network/tokamak-thanos/op-service/sources/batching/test"
	"github.com/tokamak-network/tokamak-thanos/packages/contracts-bedrock/snapshots"
)

func TestVMContract_Oracle(t *testing.T) {
	vmAbi := snapshots.LoadMIPSABI()

	stubRpc := batchingTest.NewAbiBasedRpc(t, vmAddr, vmAbi)
	vmContract := NewVMContract(vmAddr, batching.NewMultiCaller(stubRpc, batching.DefaultBatchSize))

	stubRpc.SetResponse(vmAddr, methodOracle, rpcblock.Latest, nil, []interface{}{oracleAddr})

	oracleContract, err := vmContract.Oracle(context.Background())
	require.NoError(t, err)
	tx, err := oracleContract.AddGlobalDataTx(types.NewPreimageOracleData(common.Hash{byte(preimage.Keccak256KeyType)}.Bytes(), make([]byte, 20), 0))
	require.NoError(t, err)
	// This test doesn't care about all the tx details, we just want to confirm the contract binding is using the
	// correct address
	require.Equal(t, &oracleAddr, tx.To)
}
