//go:build !rethdb

package sources

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-service/client"
	"github.com/tokamak-network/tokamak-thanos/op-service/sources/caching"
)

const buildRethdb = false

func newRecProviderFromConfig(client client.RPC, log log.Logger, metrics caching.Metrics, config *EthClientConfig) *CachingReceiptsProvider {
	return newRPCRecProviderFromConfig(client, log, metrics, config)
}
