package contracts

import (
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/contracts/registry/client"
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/contracts/registry/empty"
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/interfaces"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NewClientRegistry creates a new Registry that uses the provided client
func NewClientRegistry(c *ethclient.Client) interfaces.ContractsRegistry {
	return &client.ClientRegistry{Client: c}
}

func NewEmptyRegistry() interfaces.ContractsRegistry {
	return &empty.EmptyRegistry{}
}
