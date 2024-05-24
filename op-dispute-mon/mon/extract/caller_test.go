package extract

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"
	faultTypes "github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/game/types"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources/batching"
	batchingTest "github.com/tokamak-network/tokamak-thanos/op-service/sources/batching/test"
)

var (
	fdgAddr = common.HexToAddress("0x24112842371dFC380576ebb09Ae16Cb6B6caD7CB")
)

func TestMetadataCreator_CreateContract(t *testing.T) {
	tests := []struct {
		name        string
		game        types.GameMetadata
		expectedErr error
	}{
		{
			name: "validCannonGameType",
			game: types.GameMetadata{GameType: faultTypes.CannonGameType, Proxy: fdgAddr},
		},
		{
			name: "validAlphabetGameType",
			game: types.GameMetadata{GameType: faultTypes.AlphabetGameType, Proxy: fdgAddr},
		},
		{
			name:        "InvalidGameType",
			game:        types.GameMetadata{GameType: 2, Proxy: fdgAddr},
			expectedErr: fmt.Errorf("unsupported game type: 2"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			caller, metrics := setupMetadataLoaderTest(t)
			creator := NewGameCallerCreator(metrics, caller)
			_, err := creator.CreateContract(test.game)
			require.Equal(t, test.expectedErr, err)
			if test.expectedErr == nil {
				require.Equal(t, 1, metrics.cacheAddCalls)
				require.Equal(t, 1, metrics.cacheGetCalls)
			}
			_, err = creator.CreateContract(test.game)
			require.Equal(t, test.expectedErr, err)
			if test.expectedErr == nil {
				require.Equal(t, 1, metrics.cacheAddCalls)
				require.Equal(t, 2, metrics.cacheGetCalls)
			}
		})
	}
}

func setupMetadataLoaderTest(t *testing.T) (*batching.MultiCaller, *mockCacheMetrics) {
	fdgAbi, err := bindings.FaultDisputeGameMetaData.GetAbi()
	require.NoError(t, err)
	stubRpc := batchingTest.NewAbiBasedRpc(t, fdgAddr, fdgAbi)
	caller := batching.NewMultiCaller(stubRpc, batching.DefaultBatchSize)
	return caller, &mockCacheMetrics{}
}

type mockCacheMetrics struct {
	cacheAddCalls int
	cacheGetCalls int
}

func (m *mockCacheMetrics) CacheAdd(_ string, _ int, _ bool) {
	m.cacheAddCalls++
}
func (m *mockCacheMetrics) CacheGet(_ string, _ bool) {
	m.cacheGetCalls++
}
