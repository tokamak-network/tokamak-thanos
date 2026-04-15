#!/usr/bin/env bash
# cmd/tokamak-deployer/scripts/extract-artifacts.sh
# Usage: bash scripts/extract-artifacts.sh <forge-artifacts-dir> <output-dir>
set -euo pipefail

FORGE_ARTIFACTS_DIR="${1:-../../packages/tokamak/contracts-bedrock/forge-artifacts}"
OUTPUT_DIR="${2:-./cmd/deploy-artifacts}"

mkdir -p "$OUTPUT_DIR"

CONTRACTS=(
  "AddressManager"
  "L1CrossDomainMessenger"
  "L1ERC721Bridge"
  "L1StandardBridge"
  "L2OutputOracle"
  "OptimismMintableERC20Factory"
  "OptimismPortal"
  "OptimismPortal2"
  "ProxyAdmin"
  "SystemConfig"
  "SuperchainConfig"
  "L1Block"
  "DisputeGameFactory"
  "AnchorStateRegistry"
  "MIPS"
  "PreimageOracle"
  "Proxy"
)

for contract in "${CONTRACTS[@]}"; do
  src="$FORGE_ARTIFACTS_DIR/$contract.sol/$contract.json"
  if [ -f "$src" ]; then
    cp "$src" "$OUTPUT_DIR/$contract.json"
    echo "✅ $contract"
  else
    echo "⚠️  $contract.json not found at $src"
  fi
done

echo "Extracted $(ls "$OUTPUT_DIR"/*.json 2>/dev/null | wc -l) artifacts to $OUTPUT_DIR"
