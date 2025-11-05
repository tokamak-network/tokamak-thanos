package health

import (
	"context"

	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
)

// SupervisorHealthAPI defines the interface for the supervisor's health check.
type SupervisorHealthAPI interface {
	SyncStatus(ctx context.Context) (eth.SupervisorSyncStatus, error)
}
