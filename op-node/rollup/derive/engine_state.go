package derive

import "github.com/tokamak-network/tokamak-thanos/op-service/eth"

// EngineState provides a read-only view of the engine's block head references.
type EngineState interface {
	SafeL2Head() eth.L2BlockRef
	UnsafeL2Head() eth.L2BlockRef
}
