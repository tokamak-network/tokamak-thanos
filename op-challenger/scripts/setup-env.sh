#!/usr/bin/env bash
set -euo pipefail

##############################################################################
# Environment Variable Setup Script
#
# Generates environment variables and wallets required for L2 system deployment.
#
# Reference: /op-challenger/docs/l2-system-deployment-ko.md
##############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help message
usage() {
    cat << EOF
Environment Variable Setup Script

Usage:
    $0 [options]

Options:
    --mode MODE             Deployment mode: local|sepolia|mainnet (default: local)
    --output FILE           Environment variable file path (default: .env)
    --help                  Show this help message

Examples:
    # Local development environment setup
    $0 --mode local

    # Sepolia testnet setup
    $0 --mode sepolia

EOF
    exit 0
}

# 기본값 설정
DEPLOY_MODE="local"
ENV_FILE="${PROJECT_ROOT}/.env"

# 인자 파싱
while [[ $# -gt 0 ]]; do
    case $1 in
        --mode)
            DEPLOY_MODE="$2"
            shift 2
            ;;
        --output)
            ENV_FILE="$2"
            shift 2
            ;;
        --help)
            usage
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            ;;
    esac
done

##############################################################################
# 1. Check required tools
##############################################################################

log_info "Checking required tools..."

check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "$1 is not installed."
        exit 1
    fi
}

check_command cast
check_command openssl

log_success "All required tools are installed."

##############################################################################
# 2. Generate wallets
##############################################################################

log_info "Setting up wallets for deployment..."

# ⭐ Use Hardhat's default test accounts (same as devnet-up)
# This ensures Genesis batcherAddr matches .env BATCHER_ADDRESS

# Admin 지갑 (Hardhat Account #0)
admin_address="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
admin_key="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

# Proposer 지갑 (Hardhat Account #1)
proposer_address="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
proposer_key="0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

# Batcher 지갑 (Hardhat Account #2) ⭐ 중요!
batcher_address="0x3C44CdDdB6a900fa2b585dd299e03d12fa4293BC"
batcher_key="0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"

# Sequencer 지갑 (Hardhat Account #3)
sequencer_address="0x90F79bf6EB2c4f870365E785982E1f101E93b906"
sequencer_key="0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6"

# Challenger 지갑 (Hardhat Account #4)
challenger_address="0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"
challenger_key="0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a"

log_success "Wallet configuration completed (using Hardhat test accounts)"
log_info "These accounts match devnet-up Genesis generation"

##############################################################################
# 3. Configuration by deployment mode
##############################################################################

log_info "Deployment mode: $DEPLOY_MODE"

case $DEPLOY_MODE in
    local)
        L1_CHAIN_ID=900
        L2_CHAIN_ID=901
        L1_BLOCK_TIME=12
        L2_BLOCK_TIME=2
        L1_RPC_URL="http://localhost:8545"
        L1_BEACON_URL=""

        # Test mnemonic for local development
        MNEMONIC="test test test test test test test test test test test junk"
        ;;
    sepolia)
        L1_CHAIN_ID=11155111
        L2_CHAIN_ID=11155420
        L1_BLOCK_TIME=12
        L2_BLOCK_TIME=2

        log_warn "Please enter the following information for Sepolia deployment:"
        read -p "L1 RPC URL (e.g., https://sepolia.infura.io/v3/YOUR_KEY): " L1_RPC_URL
        read -p "L1 Beacon URL (e.g., https://sepolia.beaconcha.in): " L1_BEACON_URL

        MNEMONIC="test test test test test test test test test test test junk"
        log_warn "⚠️  Use a secure mnemonic in production environments!"
        ;;
    mainnet)
        L1_CHAIN_ID=1
        L2_CHAIN_ID=10
        L1_BLOCK_TIME=12
        L2_BLOCK_TIME=2

        log_warn "⚠️  This is a Mainnet deployment. Proceed with caution!"
        read -p "L1 RPC URL: " L1_RPC_URL
        read -p "L1 Beacon URL: " L1_BEACON_URL

        MNEMONIC="test test test test test test test test test test test junk"
        log_error "⚠️  You MUST use a secure mnemonic in production environments!"
        ;;
    *)
        log_error "Unsupported deployment mode: $DEPLOY_MODE"
        exit 1
        ;;
esac

##############################################################################
# 4. Set contract addresses
##############################################################################

# Note: For local deployment, contract addresses will be loaded from
# .devnet/addresses.json during deploy-full-stack.sh execution.
# We do NOT include them in .env for local mode!

log_info "Configuring contract addresses..."

if [ "$DEPLOY_MODE" = "local" ]; then
    log_info "Local mode: Contract addresses will be loaded from addresses.json during deployment"
    log_info "Make sure to run 'make devnet-up' before deployment to generate contracts"

    # For local, only set batch inbox address
    BATCH_INBOX_ADDRESS="0xff00000000000000000000000000000000042069"

    # DO NOT set these - deploy-full-stack.sh will load them from addresses.json
    # This prevents stale or zero addresses in .env
else
    # For testnet/mainnet, user must provide contract addresses
    log_warn "Please enter contract addresses:"
    read -p "L2OutputOracle address: " L2OO_ADDRESS
    read -p "DisputeGameFactory address: " GAME_FACTORY_ADDRESS
    read -p "Batch Inbox address (default: 0xff00...042069): " BATCH_INBOX_ADDRESS
    BATCH_INBOX_ADDRESS=${BATCH_INBOX_ADDRESS:-0xff00000000000000000000000000000000042069}
fi

##############################################################################
# 5. Generate environment variable file
##############################################################################

log_info "Generating environment variable file: $ENV_FILE"

cat > "$ENV_FILE" << EOF
##############################################################################
# L2 System Deployment Environment Variables
# Generated: $(date)
# Deployment Mode: $DEPLOY_MODE
##############################################################################

# ============================================================================
# Basic Configuration
# ============================================================================
DEPLOY_MODE=$DEPLOY_MODE

# ============================================================================
# Network Configuration
# ============================================================================
L1_CHAIN_ID=$L1_CHAIN_ID
L2_CHAIN_ID=$L2_CHAIN_ID
L1_BLOCK_TIME=$L1_BLOCK_TIME
L2_BLOCK_TIME=$L2_BLOCK_TIME

# L1 Connection Information
L1_RPC_URL=$L1_RPC_URL
L1_BEACON_URL=$L1_BEACON_URL

# ============================================================================
# Wallet Information (Admin)
# ============================================================================
ADMIN_ADDRESS=$admin_address
ADMIN_PRIVATE_KEY=$admin_key

# ============================================================================
# Wallet Information (Batcher)
# ============================================================================
BATCHER_ADDRESS=$batcher_address
BATCHER_PRIVATE_KEY=$batcher_key

# ============================================================================
# Wallet Information (Proposer)
# ============================================================================
PROPOSER_ADDRESS=$proposer_address
PROPOSER_PRIVATE_KEY=$proposer_key

# ============================================================================
# Wallet Information (Sequencer)
# ============================================================================
SEQUENCER_ADDRESS=$sequencer_address
SEQUENCER_PRIVATE_KEY=$sequencer_key

# ============================================================================
# Wallet Information (Challenger)
# ============================================================================
CHALLENGER_ADDRESS=$challenger_address
CHALLENGER_PRIVATE_KEY=$challenger_key

# ============================================================================
# Mnemonic (For Development)
# ============================================================================
MNEMONIC="$MNEMONIC"

# ============================================================================
# Contract Addresses
# ============================================================================
# Note: For local mode, these will be automatically loaded from
# .devnet/addresses.json during deploy-full-stack.sh execution
# You do NOT need to manually update these for local deployment!
BATCH_INBOX_ADDRESS=$BATCH_INBOX_ADDRESS
EOF

# Contract addresses - only for non-local mode
if [ "$DEPLOY_MODE" != "local" ]; then
    cat >> "$ENV_FILE" << EOF
L2OO_ADDRESS=$L2OO_ADDRESS
GAME_FACTORY_ADDRESS=$GAME_FACTORY_ADDRESS
EOF
fi

cat >> "$ENV_FILE" << 'EOF'

# ============================================================================
# Docker Image Tag
# ============================================================================
IMAGE_TAG=latest

# ============================================================================
# Batcher Configuration
# ============================================================================
BATCH_TYPE=0

# ============================================================================
# Plasma Configuration (Optional)
# ============================================================================
PLASMA_ENABLED=false
PLASMA_DA_SERVICE=false
PLASMA_GENERIC_DA=false

# ============================================================================
# Fault Proof Configuration
# ============================================================================
# GameType 선택:
#   0 = CANNON (MIPS VM) - Default
#   1 = PERMISSIONED_CANNON (Permissioned MIPS VM)
#   2 = ASTERISC (RISC-V VM) - New!
#   254 = FAST (Test only)
#   255 = ALPHABET (Test only)
DG_TYPE=0

# Trace Type (challenger가 사용할 VM):
#   cannon = MIPS VM (GameType 0, 1)
#   asterisc = RISC-V VM (GameType 2)
#   alphabet = Alphabet VM (GameType 255)
CHALLENGER_TRACE_TYPE=cannon

PROPOSAL_INTERVAL=12s

# Asterisc 관련 경로 (GameType 2 사용 시)
# ASTERISC_BIN=./asterisc/bin
# ASTERISC_PRESTATE=./asterisc/prestate.json

EOF

log_success "Environment variable file generated: $ENV_FILE"

##############################################################################
# 6. Display wallet information
##############################################################################

echo ""
log_info "=========================================="
log_info "Generated Wallet Information"
log_info "=========================================="
echo ""
echo -e "${YELLOW}⚠️  Backup the following information in a secure location!${NC}"
echo ""

echo "Admin Account:"
echo "  Address:     $admin_address"
echo "  Private Key: $admin_key"
echo ""

echo "Batcher Account:"
echo "  Address:     $batcher_address"
echo "  Private Key: $batcher_key"
echo ""

echo "Proposer Account:"
echo "  Address:     $proposer_address"
echo "  Private Key: $proposer_key"
echo ""

echo "Sequencer Account:"
echo "  Address:     $sequencer_address"
echo "  Private Key: $sequencer_key"
echo ""

echo "Challenger Account:"
echo "  Address:     $challenger_address"
echo "  Private Key: $challenger_key"
echo ""

##############################################################################
# 7. Next steps guide
##############################################################################

echo ""
log_info "=========================================="
log_info "Next Steps"
log_info "=========================================="
echo ""

if [ "$DEPLOY_MODE" = "local" ]; then
    echo "1. Check environment variable file:"
    echo "   cat $ENV_FILE"
    echo ""
    echo "2. Deploy full system:"
    echo "   ./op-challenger/scripts/deploy-full-stack.sh --mode local"
    echo ""
else
    echo "1. Send ETH to each wallet (for L1 gas fees)"
    echo "   - Admin:     $admin_address"
    echo "   - Batcher:   $batcher_address"
    echo "   - Proposer:  $proposer_address"
    echo "   - Sequencer: $sequencer_address"
    echo "   - Challenger: $challenger_address"
    echo ""
    echo "2. Deploy L1 contracts (if needed)"
    echo ""
    echo "3. Update environment variable file:"
    echo "   - L2OO_ADDRESS"
    echo "   - GAME_FACTORY_ADDRESS"
    echo ""
    echo "4. Deploy full system:"
    echo "   ./op-challenger/scripts/deploy-full-stack.sh --mode $DEPLOY_MODE --skip-l1"
    echo ""
fi

log_success "Environment setup completed!"

exit 0
