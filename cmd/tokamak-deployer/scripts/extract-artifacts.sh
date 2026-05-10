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
  "MultiTokenPaymaster"
  "FiatTokenV2_2"
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

# FiatTokenProxy is stored under a different name in deploy-artifacts
FIAT_PROXY_SRC="$FORGE_ARTIFACTS_DIR/FiatTokenProxy.sol/FiatTokenProxy.json"
if [ -f "$FIAT_PROXY_SRC" ]; then
  cp "$FIAT_PROXY_SRC" "$OUTPUT_DIR/FiatTokenV2_2Proxy.json"
  echo "✅ FiatTokenProxy → FiatTokenV2_2Proxy"
else
  echo "⚠️  FiatTokenProxy.json not found at $FIAT_PROXY_SRC"
fi

echo "Extracted $(ls "$OUTPUT_DIR"/*.json 2>/dev/null | wc -l) artifacts to $OUTPUT_DIR"
