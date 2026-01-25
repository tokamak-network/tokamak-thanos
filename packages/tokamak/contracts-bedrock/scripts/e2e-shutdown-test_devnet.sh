#!/bin/bash

# ForceWithdraw Shutdown E2E Devnet Script
# Usage: ./e2e-shutdown-test_devnet.sh [--dry-run]

set -e

# --- Color Constants ---
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}===============================================${NC}"
echo -e "${BLUE}   ForceWithdraw Shutdown E2E Devnet Script    ${NC}"
echo -e "${BLUE}===============================================${NC}"

# Dry-run mode (optional flag only)
DRY_RUN=false
if [[ $# -gt 0 ]]; then
  if [[ "$1" == "--dry-run" ]] && [[ $# -eq 1 ]]; then
    DRY_RUN=true
  else
    echo -e "${RED}[Error] Usage: ./e2e-shutdown-test_devnet.sh [--dry-run]${NC}"
    exit 1
  fi
fi
if $DRY_RUN; then
  echo -e "${YELLOW}🧪 Dry-run mode enabled. No network calls or transactions will be executed.${NC}"
fi

run_cmd() {
  if $DRY_RUN; then
    echo -e "${YELLOW}[DRY-RUN]${NC} $*"
  else
    "$@"
  fi
}

run_forge() {
  if $DRY_RUN; then
    echo -e "${YELLOW}[DRY-RUN]${NC} $*"
    return 0
  fi
  "$@"
}

# --- Path Configuration ---
REPO_ROOT="/Users/theo/workspace_tokamak/tokamak-thanos"
CONTRACTS_DIR="$REPO_ROOT/packages/tokamak/contracts-bedrock"
ADDRESSES_JSON="$REPO_ROOT/.devnet/addresses.json"
DEPLOY_CONFIG="$CONTRACTS_DIR/deploy-config/devnetL1.json"

# --- Helper: Read JSON value ---
read_json() {
  local file="$1"
  local key="$2"
  if [ -f "$file" ]; then
    grep -o "\"$key\": *\"[^\"]*\"" "$file" | head -1 | cut -d'"' -f4
  fi
}

read_json_raw() {
  local file="$1"
  local key="$2"
  if [ -f "$file" ]; then
    grep -o "\"$key\": *[^,}]*" "$file" | head -1 | sed 's/.*: *//' | tr -d ' "'
  fi
}

# 1. Hardcoded PRIVATE_KEY for devnet (Foundry test account #4)
export PRIVATE_KEY="0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6"

# 3. Auto-detect from devnet artifacts
echo -e "\n🔍 Loading devnet configuration..."

# --- From addresses.json (L1 contract addresses) ---
if [ -f "$ADDRESSES_JSON" ]; then
  echo -e "  📄 Found: $ADDRESSES_JSON"

  [ -z "$BRIDGE_PROXY" ] && BRIDGE_PROXY=$(read_json "$ADDRESSES_JSON" "L1StandardBridgeProxy")
  [ -z "$PROXY_ADMIN" ] && PROXY_ADMIN=$(read_json "$ADDRESSES_JSON" "ProxyAdmin")
  [ -z "$SYSTEM_OWNER_SAFE" ] && SYSTEM_OWNER_SAFE=$(read_json "$ADDRESSES_JSON" "SystemOwnerSafe")
  [ -z "$OPTIMISM_PORTAL_PROXY" ] && OPTIMISM_PORTAL_PROXY=$(read_json "$ADDRESSES_JSON" "OptimismPortalProxy")
  [ -z "$SUPERCHAIN_CONFIG_PROXY" ] && SUPERCHAIN_CONFIG_PROXY=$(read_json "$ADDRESSES_JSON" "SuperchainConfigProxy")
  [ -z "$L1_USDC_BRIDGE_PROXY" ] && L1_USDC_BRIDGE_PROXY=$(read_json "$ADDRESSES_JSON" "L1UsdcBridgeProxy")
else
  echo -e "${YELLOW}  [WARN] addresses.json not found at $ADDRESSES_JSON${NC}"
fi

# --- From deploy-config/devnetL1.json ---
if [ -f "$DEPLOY_CONFIG" ]; then
  echo -e "  📄 Found: $DEPLOY_CONFIG"

  [ -z "$L1_NATIVE_TOKEN" ] && L1_NATIVE_TOKEN=$(read_json "$DEPLOY_CONFIG" "nativeTokenAddress")
  [ -z "$L1_USDC" ] && L1_USDC=$(read_json "$DEPLOY_CONFIG" "l1UsdcAddr")
  [ -z "$GUARDIAN_SAFE" ] && GUARDIAN_SAFE=$(read_json "$DEPLOY_CONFIG" "superchainConfigGuardian")
else
  echo -e "${YELLOW}  [WARN] deploy-config not found at $DEPLOY_CONFIG${NC}"
fi

# --- L2 Predeploy addresses (constants) ---
L2_USDC_BRIDGE_PROXY=${L2_USDC_BRIDGE_PROXY:-"0x4200000000000000000000000000000000000775"}

# --- Default RPC URLs for devnet ---
export L1_RPC_URL=${L1_RPC_URL:-http://localhost:8545}
export L2_RPC_URL=${L2_RPC_URL:-http://localhost:9545}

# --- Validate required addresses ---
MISSING_VARS=""
[ -z "$BRIDGE_PROXY" ] && MISSING_VARS="$MISSING_VARS BRIDGE_PROXY"
[ -z "$PROXY_ADMIN" ] && MISSING_VARS="$MISSING_VARS PROXY_ADMIN"
[ -z "$SYSTEM_OWNER_SAFE" ] && MISSING_VARS="$MISSING_VARS SYSTEM_OWNER_SAFE"
[ -z "$L1_NATIVE_TOKEN" ] && MISSING_VARS="$MISSING_VARS L1_NATIVE_TOKEN"

if [ -n "$MISSING_VARS" ]; then
  echo -e "${RED}[Error] Missing required environment variables:$MISSING_VARS${NC}"
  echo -e "${YELLOW}Please ensure devnet is running or set these in .env${NC}"
  exit 1
fi

# --- Export all environment variables ---
export BRIDGE_PROXY
export PROXY_ADMIN
export SYSTEM_OWNER_SAFE
export OPTIMISM_PORTAL_PROXY
export SUPERCHAIN_CONFIG_PROXY
export L1_USDC_BRIDGE_PROXY
export L2_USDC_BRIDGE_PROXY
export L1_NATIVE_TOKEN
export L1_USDC
export GUARDIAN_SAFE=${GUARDIAN_SAFE:-$SYSTEM_OWNER_SAFE}

# --- Print loaded configuration ---
echo -e "\n${GREEN}✅ Configuration loaded:${NC}"
echo -e "  BRIDGE_PROXY:           $BRIDGE_PROXY"
echo -e "  PROXY_ADMIN:            $PROXY_ADMIN"
echo -e "  SYSTEM_OWNER_SAFE:      $SYSTEM_OWNER_SAFE"
echo -e "  OPTIMISM_PORTAL_PROXY:  $OPTIMISM_PORTAL_PROXY"
echo -e "  SUPERCHAIN_CONFIG_PROXY: $SUPERCHAIN_CONFIG_PROXY"
echo -e "  L1_NATIVE_TOKEN:        $L1_NATIVE_TOKEN"
echo -e "  L1_USDC_BRIDGE_PROXY:   $L1_USDC_BRIDGE_PROXY"
echo -e "  L2_USDC_BRIDGE_PROXY:   $L2_USDC_BRIDGE_PROXY"
echo -e "  L1_RPC_URL:             $L1_RPC_URL"
echo -e "  L2_RPC_URL:             $L2_RPC_URL"

# --- Derive execution account ---
if $DRY_RUN; then
  PRIVATE_KEY_ADDRESS="0x[DRY-RUN]"
else
  PRIVATE_KEY_ADDRESS=$(cast wallet address $PRIVATE_KEY)
fi
echo -e "\n🔑 Execution Account: $PRIVATE_KEY_ADDRESS"

# --- Change to contracts-bedrock directory ---
cd "$CONTRACTS_DIR"

# Step 1: Block Deposits and Withdrawals (L1)
echo -e "\n${GREEN}[Step 1] Blocking Deposits and Withdrawals (L1)...${NC}"
run_forge forge script scripts/shutdown/BlockDepositsWithdrawals.s.sol:BlockDepositsWithdrawals \
  --fork-url $L1_RPC_URL \
  --broadcast \
  --slow \
  --private-key $PRIVATE_KEY

# Step 2: Collect Asset Data (L2)
echo -e "\n${GREEN}[Step 2] Collecting Asset Data (L2)...${NC}"
if $DRY_RUN; then
  L2_CHAIN_ID=${L2_CHAIN_ID:-901}
else
  L2_CHAIN_ID=$(cast chain-id --rpc-url $L2_RPC_URL)
  if [ -z "$L2_CHAIN_ID" ]; then
    echo -e "${RED}[Error] Failed to read L2 chain id from RPC.${NC}"
    exit 1
  fi
fi
export L2_CHAIN_ID

export L2_START_BLOCK=${L2_START_BLOCK:-0}
run_cmd python3 scripts/shutdown/fetch_explorer_assets.py $L2_CHAIN_ID
run_cmd python3 scripts/shutdown/compute_l2_burns.py $L2_RPC_URL $L2_CHAIN_ID $L2_START_BLOCK
run_cmd python3 scripts/shutdown/compute_finalized_native_withdrawals.py $L1_RPC_URL $BRIDGE_PROXY

# Step 3: Generate Asset Snapshot (Off-chain)
echo -e "\n${GREEN}[Step 3] Generating Asset Snapshot (Off-chain)...${NC}"
export SKIP_FETCH=true

run_forge forge script scripts/shutdown/GenerateAssetSnapshot.s.sol:GenerateAssetSnapshot \
  --fork-url $L2_RPC_URL \
  --ffi

SNAPSHOT_ASSETS_PATH="data/generate-assets-${L2_CHAIN_ID}.json"
if ! $DRY_RUN; then
  if [ ! -f "$SNAPSHOT_ASSETS_PATH" ]; then
    echo -e "${RED}[Error] Snapshot file not found: $SNAPSHOT_ASSETS_PATH${NC}"
    exit 1
  fi
fi

echo "Snapshot output path:"
echo "  - assets:     $SNAPSHOT_ASSETS_PATH"

# Step 4: Snapshot Registration + Bridge Upgrades + Activation (L1)
echo -e "\n${GREEN}[Step 4] Registering Snapshot + Upgrading Bridge...${NC}"
export DATA_PATH=$SNAPSHOT_ASSETS_PATH

# Capture output to extract GenFWStorage address
OUTPUT=$(run_forge forge script scripts/shutdown/PrepareL1Withdrawal.s.sol:PrepareL1Withdrawal \
  --fork-url $L1_RPC_URL \
  --via-ir \
  --broadcast \
  --slow \
  --private-key $PRIVATE_KEY)

echo "$OUTPUT"

# Extract GenFWStorage address from logs
# Pattern: [INFO] Deployed GenFWStorage at: 0x...
STORAGE_ADDRESS=$(echo "$OUTPUT" | grep "Deployed GenFWStorage at:" | awk '{print $NF}')

if [ -n "$STORAGE_ADDRESS" ]; then
  echo -e "\n${GREEN}✅ Auto-detected GenFWStorage: $STORAGE_ADDRESS${NC}"
  export STORAGE_ADDRESS
else
  echo -e "\n${RED}[Error] Failed to detect GenFWStorage address from output.${NC}"
  # Don't exit here to allow debugging, but warn heavily
fi

# Step 5: Final Withdrawals and Claims (L1)
echo -e "\n${GREEN}[Step 5] Executing Final Withdrawals and Claims (L1)...${NC}"
export EXECUTE_CLAIMS=${EXECUTE_CLAIMS:-true}

if [ -z "$STORAGE_ADDRESS" ]; then
  echo -e "${YELLOW}[WARN] STORAGE_ADDRESS not set. ExecuteL1Withdrawal may fail.${NC}"
fi

run_forge forge script scripts/shutdown/ExecuteL1Withdrawal.s.sol:ExecuteL1Withdrawal \
  --fork-url $L1_RPC_URL \
  --via-ir \
  --broadcast \
  --slow \
  --private-key $PRIVATE_KEY

echo -e "\n${BLUE}===============================================${NC}"
echo -e "${GREEN}           Devnet E2E Completed               ${NC}"
echo -e "${BLUE}===============================================${NC}"
