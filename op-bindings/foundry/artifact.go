package foundry

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tokamak-network/tokamak-thanos/op-bindings/solc"
)

// Artifact represents a foundry compilation artifact.
// The Abi is specifically left as a json.RawMessage because
// round trip marshaling/unmarshaling of the abi.ABI type
// causes issues.
type Artifact struct {
	Abi              json.RawMessage    `json:"abi"`
	StorageLayout    solc.StorageLayout `json:"storageLayout"`
	DeployedBytecode DeployedBytecode   `json:"deployedBytecode"`
	Bytecode         Bytecode           `json:"bytecode"`
}

type DeployedBytecode struct {
	SourceMap           string          `json:"sourceMap"`
	Object              hexutil.Bytes   `json:"object"`
	LinkReferences      json.RawMessage `json:"linkReferences"`
	ImmutableReferences json.RawMessage `json:"immutableReferences"`
}

type Bytecode struct {
	SourceMap      string          `json:"sourceMap"`
	Object         hexutil.Bytes   `json:"object"`
	LinkReferences json.RawMessage `json:"linkReferences"`
}
