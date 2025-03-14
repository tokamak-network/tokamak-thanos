package empty

import (
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/interfaces"
	"github.com/tokamak-network/tokamak-thanos/devnet-sdk/types"
)

// EmptyRegistry represents a registry that returns not found errors for all contract accesses
type EmptyRegistry struct{}

var _ interfaces.ContractsRegistry = (*EmptyRegistry)(nil)

func (r *EmptyRegistry) SuperchainWETH(address types.Address) (interfaces.SuperchainWETH, error) {
	return nil, &interfaces.ErrContractNotFound{
		ContractType: "SuperchainWETH",
		Address:      address,
	}
}
