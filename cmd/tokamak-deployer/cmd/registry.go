// cmd/tokamak-deployer/cmd/registry.go
package cmd

import "embed"

// DefaultRegistryFS holds per-L1-chainId reuse registries embedded into the binary.
// Loaded at runtime by deployer.loadRegistry when --reuse-deployment is set.
//
//go:embed registry
var DefaultRegistryFS embed.FS
