// cmd/tokamak-deployer/cmd/assets.go
package cmd

import "embed"

// DeployArtifactsFS contains the ABI+bytecode for all L1 contracts needed for deployment.
// Populated by scripts/extract-artifacts.sh before goreleaser build.
//
//go:embed deploy-artifacts
var DeployArtifactsFS embed.FS
