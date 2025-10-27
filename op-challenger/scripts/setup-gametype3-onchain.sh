#!/usr/bin/env bash

##############################################################################
# GameType 3 (AsteriscKona) On-Chain Setup Script
#
# This script registers GameType 3 to DisputeGameFactory.
# GameType 3 shares the same RISCV.sol implementation with GameType 2,
# but uses kona-client (Rust) instead of op-program (Go) for off-chain proving.
#
# Usage:
#   ./setup-gametype3-onchain.sh
#
# Requirements:
#   - RISCV.sol already deployed (from GameType 2 setup)
#   - DisputeGameFactory deployed
#   - Admin private key with permissions
##############################################################################

set -euo pipefail

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Source common library
source "${SCRIPT_DIR}/lib/common.sh"

##############################################################################
# Configuration
##############################################################################

# Load environment variables
ENV_FILE="${PROJECT_ROOT}/.env"
if [ ! -f "$ENV_FILE" ]; then
    log_error "Environment file not found: $ENV_FILE"
    log_error "Please run deploy-modular.sh first to generate .env"
    exit 1
fi

source "$ENV_FILE"

# Required variables
ADMIN_PRIVATE_KEY="${ADMIN_PRIVATE_KEY:-}"
L1_RPC_URL="${L1_RPC_URL:-http://localhost:8545}"
DISPUTE_GAME_FACTORY="${GAME_FACTORY_ADDRESS:-}"

# Addresses file
ADDRESSES_FILE="${PROJECT_ROOT}/.devnet/addresses.json"

##############################################################################
# Functions
##############################################################################

get_riscv_address() {
    if [ -f "$ADDRESSES_FILE" ]; then
        local riscv_addr=$(jq -r '.RISCV // .RISCVProxy' "$ADDRESSES_FILE" 2>/dev/null || echo "")
        if [ -n "$riscv_addr" ] && [ "$riscv_addr" != "null" ]; then
            echo "$riscv_addr"
            return 0
        fi
    fi

    log_error "RISCV.sol address not found in $ADDRESSES_FILE"
    return 1
}

check_gametype_impl() {
    local game_type="$1"
    local factory_addr="$2"

    log_info "Checking GameType $game_type implementation..."

    local impl_addr=$(cast call "$factory_addr" "gameImpls(uint32)" "$game_type" --rpc-url "$L1_RPC_URL" 2>/dev/null || echo "")

    if [ -z "$impl_addr" ] || [ "$impl_addr" = "0x0000000000000000000000000000000000000000000000000000000000000000" ]; then
        echo ""
        return 1
    else
        # Convert to checksummed address
        echo "$impl_addr" | sed 's/^0x000000000000000000000000/0x/'
        return 0
    fi
}

register_gametype() {
    local game_type="$1"
    local impl_addr="$2"
    local factory_addr="$3"

    log_info "Registering GameType $game_type..."
    log_info "  Implementation: $impl_addr"
    log_info "  Factory: $factory_addr"

    if [ -z "$ADMIN_PRIVATE_KEY" ]; then
        log_error "ADMIN_PRIVATE_KEY not set"
        return 1
    fi

    local tx_hash=$(cast send "$factory_addr" \
        "setImplementation(uint32,address)" \
        "$game_type" \
        "$impl_addr" \
        --private-key "$ADMIN_PRIVATE_KEY" \
        --rpc-url "$L1_RPC_URL" \
        --json 2>/dev/null | jq -r '.transactionHash' || echo "")

    if [ -n "$tx_hash" ] && [ "$tx_hash" != "null" ]; then
        log_success "Transaction sent: $tx_hash"
        log_info "Waiting for confirmation..."
        sleep 3
        return 0
    else
        log_error "Transaction failed"
        return 1
    fi
}

##############################################################################
# Main
##############################################################################

main() {
    log_info "=========================================="
    log_info "GameType 3 (AsteriscKona) On-Chain Setup"
    log_info "=========================================="
    echo ""

    # Check required variables
    if [ -z "$DISPUTE_GAME_FACTORY" ]; then
        log_error "GAME_FACTORY_ADDRESS not set in .env"
        log_error "Please run deploy-modular.sh first"
        exit 1
    fi

    log_info "DisputeGameFactory: $DISPUTE_GAME_FACTORY"
    log_info "L1 RPC: $L1_RPC_URL"
    echo ""

    # Get RISCV.sol address
    log_info "Step 1: Get RISCV.sol address"
    log_info "=========================================="
    RISCV_ADDRESS=$(get_riscv_address)
    if [ -z "$RISCV_ADDRESS" ]; then
        log_error "Failed to get RISCV.sol address"
        log_error ""
        log_error "Please ensure GameType 2 (Asterisc) is deployed first:"
        log_error "  ./deploy-modular.sh --dg-type 2"
        exit 1
    fi

    log_success "RISCV.sol found: $RISCV_ADDRESS"
    echo ""

    # Verify RISCV.sol is deployed
    log_info "Step 2: Verify RISCV.sol deployment"
    log_info "=========================================="
    local version=$(cast call "$RISCV_ADDRESS" "version()" --rpc-url "$L1_RPC_URL" 2>/dev/null || echo "")
    if [ -n "$version" ]; then
        log_success "RISCV.sol is deployed (version: $version)"
    else
        log_error "RISCV.sol not found at $RISCV_ADDRESS"
        exit 1
    fi
    echo ""

    # Check GameType 2 registration (should exist)
    log_info "Step 3: Verify GameType 2 registration"
    log_info "=========================================="
    local gt2_impl=$(check_gametype_impl 2 "$DISPUTE_GAME_FACTORY")
    if [ -n "$gt2_impl" ]; then
        log_success "GameType 2 is registered: $gt2_impl"

        if [ "$gt2_impl" != "$RISCV_ADDRESS" ]; then
            log_warn "⚠️  GameType 2 implementation mismatch!"
            log_warn "   Expected: $RISCV_ADDRESS"
            log_warn "   Got: $gt2_impl"
        fi
    else
        log_warn "GameType 2 is not registered (this is unusual)"
        log_warn "GameType 3 will still work if RISCV.sol is deployed"
    fi
    echo ""

    # Check if GameType 3 is already registered
    log_info "Step 4: Check GameType 3 registration"
    log_info "=========================================="
    local gt3_impl=$(check_gametype_impl 3 "$DISPUTE_GAME_FACTORY")
    if [ -n "$gt3_impl" ]; then
        log_success "GameType 3 is already registered: $gt3_impl"

        if [ "$gt3_impl" = "$RISCV_ADDRESS" ]; then
            log_success "✅ GameType 3 correctly points to RISCV.sol"
            log_info ""
            log_info "No action needed. GameType 3 is ready to use!"
            exit 0
        else
            log_warn "⚠️  GameType 3 points to different implementation:"
            log_warn "   Current: $gt3_impl"
            log_warn "   Expected: $RISCV_ADDRESS"
            echo ""
            read -p "Do you want to update it? (y/N) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_info "Skipping update"
                exit 0
            fi
        fi
    else
        log_info "GameType 3 is not registered yet"
    fi
    echo ""

    # Register GameType 3
    log_info "Step 5: Register GameType 3"
    log_info "=========================================="
    log_info "This will register GameType 3 (AsteriscKona) to use RISCV.sol"
    log_info "Same on-chain VM as GameType 2, but different off-chain prover (kona-client)"
    echo ""

    if register_gametype 3 "$RISCV_ADDRESS" "$DISPUTE_GAME_FACTORY"; then
        log_success "GameType 3 registration transaction sent"
    else
        log_error "Failed to register GameType 3"
        exit 1
    fi
    echo ""

    # Verify registration
    log_info "Step 6: Verify registration"
    log_info "=========================================="
    sleep 2

    local gt3_impl_after=$(check_gametype_impl 3 "$DISPUTE_GAME_FACTORY")
    if [ -n "$gt3_impl_after" ] && [ "$gt3_impl_after" = "$RISCV_ADDRESS" ]; then
        log_success "✅ GameType 3 successfully registered!"
        log_success "   Implementation: $gt3_impl_after"
        echo ""
        log_success "=========================================="
        log_success "GameType 3 On-Chain Setup Complete!"
        log_success "=========================================="
        echo ""
        log_info "Next steps:"
        log_info "  1. Deploy Challenger with GameType 3:"
        log_info "     ./deploy-modular.sh --dg-type 3"
        log_info ""
        log_info "  2. Create a test game:"
        log_info "     cast send $DISPUTE_GAME_FACTORY \\"
        log_info "       \"create(uint32,bytes32,bytes)\" \\"
        log_info "       3 \\"
        log_info "       <OUTPUT_ROOT> \\"
        log_info "       0x \\"
        log_info "       --private-key \$ADMIN_PRIVATE_KEY \\"
        log_info "       --rpc-url $L1_RPC_URL"
    else
        log_error "❌ GameType 3 registration verification failed"
        log_error "   Expected: $RISCV_ADDRESS"
        log_error "   Got: $gt3_impl_after"
        exit 1
    fi
}

##############################################################################
# Execute
##############################################################################

main "$@"
