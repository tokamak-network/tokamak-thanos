#!/usr/bin/env bash

##############################################################################
# Environment Setup Module
#
# This module handles:
# - JWT secret generation
# - P2P key generation
# - Environment variable loading from addresses.json
# - Private key configuration
##############################################################################

set -euo pipefail

# Source common library
_MODULE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${_MODULE_DIR}/../lib/common.sh"

##############################################################################
# Configuration
##############################################################################

DEVNET_DIR="${PROJECT_ROOT}/.devnet"

##############################################################################
# File Preparation Functions
##############################################################################

fix_directory_mount() {
    local file_path="$1"
    local file_name=$(basename "$file_path")

    if [ -d "$file_path" ]; then
        log_warn "$file_name is a directory (from Docker mount), removing..."
        rm -rf "$file_path"
        return 0
    fi
    return 1
}

ensure_jwt_secret() {
    local jwt_secret="${DEVNET_DIR}/jwt-secret.txt"

    fix_directory_mount "$jwt_secret"

    if [ ! -f "$jwt_secret" ]; then
        log_info "Generating JWT secret: $jwt_secret"
        openssl rand -hex 32 > "$jwt_secret"
        log_success "JWT secret generated"
    else
        log_info "Using existing JWT secret: $jwt_secret"
    fi

    return 0
}

ensure_p2p_key() {
    local key_file="$1"
    local key_name="$2"

    # Fix directory mount issue
    if ! fix_directory_mount "$key_file"; then
        log_warn "$key_file was a directory, removed and will recreate as file"
    fi

    if [ ! -f "$key_file" ]; then
        log_info "Generating $key_name: $key_file"
        if openssl rand -hex 32 > "$key_file" 2>/dev/null; then
            log_success "$key_name generated"
        else
            log_error "Failed to generate $key_name"
            return 1
        fi
    else
        log_info "Using existing $key_name: $key_file"
    fi

    # Verify file was created successfully
    if [ ! -f "$key_file" ] || [ ! -s "$key_file" ]; then
        log_error "$key_name file is missing or empty: $key_file"
        return 1
    fi

    return 0
}

ensure_p2p_keys() {
    ensure_p2p_key "${DEVNET_DIR}/p2p-node-key.txt" "P2P Node key"
    ensure_p2p_key "${DEVNET_DIR}/p2p-challenger-key.txt" "P2P Challenger key"
    return 0
}

prepare_required_files() {
    log_info "=========================================="
    log_info "Prepare Required Files"
    log_info "=========================================="
    echo ""

    # Ensure devnet directory exists
    ensure_dir "$DEVNET_DIR"

    # Generate/check required files
    if ! ensure_jwt_secret; then
        log_error "Failed to prepare JWT secret"
        return 1
    fi

    if ! ensure_p2p_keys; then
        log_error "Failed to prepare P2P keys"
        return 1
    fi

    echo ""
    log_success "Required files prepared"
    echo ""

    return 0
}

##############################################################################
# Environment Variable Functions
##############################################################################

load_contract_addresses() {
    local addresses_file="${DEVNET_DIR}/addresses.json"

    if [ ! -f "$addresses_file" ]; then
        log_warn "Addresses file not found: $addresses_file"
        log_warn "Contract addresses will not be available"
        return 1
    fi

    log_info "Loading contract addresses from: $addresses_file"

    export GAME_FACTORY_ADDRESS=$(jq -r '.DisputeGameFactoryProxy' "$addresses_file" 2>/dev/null || echo "")
    export L2OO_ADDRESS=$(jq -r '.L2OutputOracleProxy' "$addresses_file" 2>/dev/null || echo "")

    if [ -n "$GAME_FACTORY_ADDRESS" ] && [ "$GAME_FACTORY_ADDRESS" != "null" ]; then
        log_success "Contract addresses loaded:"
        log_info "  DisputeGameFactory: $GAME_FACTORY_ADDRESS"
        log_info "  L2OutputOracle: $L2OO_ADDRESS"
        return 0
    else
        log_warn "Failed to load contract addresses"
        return 1
    fi
}

load_private_keys() {
    log_info "Loading private keys (hardhat default accounts)..."

    # These are hardhat's default accounts - ONLY FOR DEVELOPMENT
    # Must match devnetL1-template.json configuration:
    # - l2OutputOracleProposer: Account #1 (0x70997970C51812dc3A010C7d01b50e0d17dc79C8)
    # - batchSenderAddress:     Account #2 (0x3C44CdDdB6a900fa2b585dd299e03d12fa4293BC)
    # - Challenger:             Account #3 (0x90F79bf6EB2c4f870365E785982E1f101E93b906)

    export PROPOSER_PRIVATE_KEY="0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"    # Account #1
    export BATCHER_PRIVATE_KEY="0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"     # Account #2
    export CHALLENGER_PRIVATE_KEY="0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6" # Account #3

    log_success "Private keys loaded (⚠️  development only)"

    return 0
}

load_environment_variables() {
    log_info "=========================================="
    log_info "Load Environment Variables"
    log_info "=========================================="
    echo ""

    load_contract_addresses
    load_private_keys

    echo ""
    log_success "Environment variables loaded"
    echo ""

    return 0
}

setup_environment() {
    prepare_required_files
    load_environment_variables
    return 0
}

##############################################################################
# CLI Interface (if run directly)
##############################################################################

if [ "${BASH_SOURCE[0]}" -ef "$0" ]; then
    # Script is being run directly, not sourced

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "This script prepares required files and loads environment variables."
                echo ""
                echo "Options:"
                echo "  --help    Show this help"
                echo ""
                echo "Examples:"
                echo "  $0    # Setup environment"
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                exit 1
                ;;
        esac
    done

    # Initialize logging
    init_logging

    # Run setup
    setup_environment

    exit 0
fi
