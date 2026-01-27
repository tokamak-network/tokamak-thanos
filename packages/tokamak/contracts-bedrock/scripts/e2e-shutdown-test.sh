#!/bin/bash

# ForceWithdraw Shutdown E2E Script (Consolidated Workflow)
# Usage: ./e2e-shutdown-test.sh [--dry-run]
#
# Required env:
# - BRIDGE_PROXY (L1 bridge proxy)
# - SYSTEM_OWNER_SAFE
# - PRIVATE_KEY
# - L1_RPC_URL
# - L2_RPC_URL
# - OPTIMISM_PORTAL_PROXY

set -e

# --- Color Constants ---
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}===============================================${NC}"
echo -e "${BLUE}     Consolidated Shutdown E2E Test Script     ${NC}"
echo -e "${BLUE}===============================================${NC}"

# Dry-run mode
DRY_RUN=false
if [[ "$1" == "--dry-run" ]]; then
  DRY_RUN=true
  export DRY_RUN=true
  echo -e "${YELLOW}🧪 Dry-run mode enabled.${NC}"
fi

# 1. Load .env
if [ -f .env ]; then
  export $(echo $(grep -v '^#' .env | xargs) | envsubst)
fi

# 2. Prerequisites check
if [ -z "$BRIDGE_PROXY" ] || [ -z "$SYSTEM_OWNER_SAFE" ] || [ -z "$PRIVATE_KEY" ]; then
  echo -e "${RED}[Error] BRIDGE_PROXY, SYSTEM_OWNER_SAFE, and PRIVATE_KEY must be set.${NC}"
  exit 1
fi

# 3. Step 1: Block Deposits (L1)
echo -e "\n${GREEN}[Step 1] Blocking Deposits and Withdrawals (L1)...${NC}"
forge script scripts/shutdown/BlockDepositsWithdrawals.s.sol:BlockDepositsWithdrawals \
  --fork-url $L1_RPC_URL --broadcast --private-key $PRIVATE_KEY

# 4. Step 2: Data Fetching (L2)
echo -e "\n${GREEN}[Step 2] Fetching Asset Data from L2...${NC}"
L2_CHAIN_ID=$(cast chain-id --rpc-url $L2_RPC_URL)
python3 scripts/shutdown/fetch_explorer_assets.py $L2_CHAIN_ID
python3 scripts/shutdown/compute_l2_burns.py $L2_RPC_URL $L2_CHAIN_ID
python3 scripts/shutdown/compute_finalized_native_withdrawals.py $L1_RPC_URL $BRIDGE_PROXY

# 5. Step 3: Generate Snapshot (L2)
echo -e "\n${GREEN}[Step 3] Generating Asset Snapshot...${NC}"
export SKIP_FETCH=true
forge script scripts/shutdown/GenerateAssetSnapshot.s.sol:GenerateAssetSnapshot \
  --fork-url $L2_RPC_URL --ffi

SNAPSHOT_PATH="data/generate-assets-${L2_CHAIN_ID}.json"

# 6. Step 4: Prepare & Activate (L1)
echo -e "\n${GREEN}[Step 4] Preparing L1 Withdrawal (Upgrades & Storage)...${NC}"
export DATA_PATH=$SNAPSHOT_PATH
# Capture GenFWStorage address from output
OUTPUT=$(forge script scripts/shutdown/PrepareL1Withdrawal.s.sol:PrepareL1Withdrawal \
  --fork-url $L1_RPC_URL --broadcast --private-key $PRIVATE_KEY)

echo "$OUTPUT"

STORAGE_ADDR=$(echo "$OUTPUT" | grep "Deployed GenFWStorage at:" | awk '{print $NF}')

if [ -z "$STORAGE_ADDR" ] && [ "$DRY_RUN" = false ]; then
  echo -e "${RED}[Error] Failed to extract GenFWStorage address.${NC}"
  exit 1
fi

echo -e "📍 Deployed Storage: ${BLUE}$STORAGE_ADDR${NC}"

# 7. Step 5: Execute Claims (L1)
echo -e "\n${GREEN}[Step 5] Executing Claims and Liquidity Sweep (L1)...${NC}"
export STORAGE_ADDRESS=$STORAGE_ADDR
export EXECUTE_CLAIMS=true

forge script scripts/shutdown/ExecuteL1Withdrawal.s.sol:ExecuteL1Withdrawal \
  --fork-url $L1_RPC_URL --broadcast --private-key $PRIVATE_KEY

echo -e "\n${BLUE}===============================================${NC}"
echo -e "${GREEN}          E2E Shutdown Test Completed          ${NC}"
echo -e "${BLUE}===============================================${NC}"
