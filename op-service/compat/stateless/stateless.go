// Package stateless provides shims for go-ethereum/core/stateless types
// that don't exist in the tokamak-thanos-geth fork.
package stateless

// Witness holds stateless execution witness data.
type Witness struct {
	Headers [][]byte
	State   map[string][]byte
}
