#!/usr/bin/env bash

##############################################################################
# Modular Deployment Script
#
# This is a refactored version of deploy-full-stack.sh that uses modularized
# components. It demonstrates how to orchestrate the deployment using the
# new module system.
#
# Usage:
#   ./deploy-modular.sh [options]
#
# Options:
#   --mode local|existing    Deployment mode (default: local)
#   --dg-type 0|1|2|254|255  Dispute game type (default: 0)
#   --help                   Show this help
#
##############################################################################

set -euo pipefail

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
export PROJECT_ROOT

# Source common library
source "${SCRIPT_DIR}/lib/common.sh"

# Source modules
source "${SCRIPT_DIR}/modules/cleanup.sh"
source "${SCRIPT_DIR}/modules/vm-build.sh"
source "${SCRIPT_DIR}/modules/genesis.sh"
source "${SCRIPT_DIR}/modules/env-setup.sh"
source "${SCRIPT_DIR}/modules/docker-deploy.sh"

##############################################################################
# Configuration
##############################################################################

DEPLOY_MODE="${DEPLOY_MODE:-local}"
DG_TYPE="${DG_TYPE:-0}"
CHALLENGER_TRACE_TYPE="${CHALLENGER_TRACE_TYPE:-cannon}"

# Docker compose command - check which version is available
if docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    echo "ERROR: Neither 'docker compose' nor 'docker-compose' is available"
    exit 1
fi

##############################################################################
# Argument Parsing
##############################################################################

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --mode)
                DEPLOY_MODE="$2"
                shift 2
                ;;
            --dg-type)
                DG_TYPE="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --mode local|existing    Deployment mode (default: local)"
                echo "  --dg-type 0|1|2|254|255  Dispute game type (default: 0)"
                echo "  --help                   Show this help"
                echo ""
                echo "Environment Variables:"
                echo "  FAULT_GAME_MAX_CLOCK_DURATION    Max clock duration in seconds"
                echo "  FAULT_GAME_WITHDRAWAL_DELAY      Withdrawal delay in seconds"
                echo "  PROPOSAL_INTERVAL                Proposal interval"
                echo ""
                echo "Examples:"
                echo "  $0                              # Deploy with default settings (GameType 0)"
                echo "  $0 --dg-type 2                  # Deploy with GameType 2 (Asterisc)"
                echo "  FAULT_GAME_MAX_CLOCK_DURATION=120 $0  # Custom clock duration"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
}

##############################################################################
# Main Deployment Flow
##############################################################################

get_prestate_hash_for_gametype() {
    local dg_type="$1"
    local prestate_hash=""

    case "$dg_type" in
        0|1|254|255)
            # Cannon (MIPS) prestate
            local cannon_prestate="${PROJECT_ROOT}/op-program/bin/prestate-proof.json"
            if [ -f "$cannon_prestate" ]; then
                prestate_hash=$(jq -r '.pre' "$cannon_prestate" 2>/dev/null || echo "")
                if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                    # Log to stderr to avoid polluting stdout (function return value)
                    log_info "Using Cannon (MIPS) prestate: $prestate_hash" >&2
                else
                    log_error "Failed to extract prestate from $cannon_prestate" >&2
                fi
            else
                log_error "Cannon prestate file not found: $cannon_prestate" >&2
            fi
            ;;
        2)
            # Asterisc (RISC-V) prestate - read directly from Docker build output
            local asterisc_prestate="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
            if [ -f "$asterisc_prestate" ]; then
                # Read stateHash directly from prestate-proof.json (Docker reproducible build output)
                # Note: Asterisc uses .stateHash field (not .pre like Cannon)
                prestate_hash=$(cat "$asterisc_prestate" | jq -r '.stateHash' 2>/dev/null || echo "")
                if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                    # Log to stderr to avoid polluting stdout (function return value)
                    log_info "Using Asterisc (RISC-V) prestate from Docker build: $prestate_hash" >&2
                else
                    log_error "Failed to extract .stateHash from $asterisc_prestate" >&2
                fi
            else
                log_error "Asterisc prestate not found: $asterisc_prestate" >&2
                log_error "❌ GameType 2 requires RISC-V prestate!" >&2
                log_error "   Cannot use MIPS prestate for RISC-V architecture." >&2
                echo "" >&2
                log_error "Please ensure Asterisc VM build succeeded:" >&2
                log_error "  - Check VM build logs above" >&2
                log_error "  - Verify asterisc/bin/prestate-proof.json exists" >&2
                log_error "  - Or use GameType 0/1 (Cannon) instead" >&2
                # Do NOT fallback to Cannon prestate - return empty
            fi
            ;;
    esac

    echo "$prestate_hash"
}

main() {
    parse_arguments "$@"

    # Initialize logging
    init_logging

    log_info "=========================================="
    log_info "Tokamak L2 Full Stack Deployment (Modular)"
    log_info "=========================================="
    echo ""
    log_info "Deployment Mode: $DEPLOY_MODE"
    log_info "Dispute Game Type: $DG_TYPE ($(get_gametype_name $DG_TYPE))"
    log_info "Trace Type: $CHALLENGER_TRACE_TYPE"
    echo ""

    # Validate GameType and TraceType compatibility
    validate_gametype_tracetype "$DG_TYPE" "$CHALLENGER_TRACE_TYPE"

    # Step 0: Cleanup existing containers and devnet files
    cleanup_containers
    cleanup_devnet_files

    # Step 1: Build VMs
    log_info "=========================================="
    log_info "Step 1: Build VMs"
    log_info "=========================================="
    echo ""

    local build_cannon="false"
    local build_asterisc="false"

    case "$DG_TYPE" in
        0|1|254|255)
            build_cannon="true"
            ;;
        2)
            build_asterisc="true"
            ;;
    esac

    if ! build_vms "$build_cannon" "$build_asterisc"; then
        log_error "VM build failed"
        exit 1
    fi

    # Get the correct prestate hash for this GameType
    local prestate_hash=$(get_prestate_hash_for_gametype "$DG_TYPE")

    if [ -z "$prestate_hash" ] || [ "$prestate_hash" = "null" ]; then
        log_error "Failed to get prestate hash for GameType $DG_TYPE"
        exit 1
    fi

    log_success "Prestate hash for GameType $DG_TYPE: $prestate_hash"
    echo ""

    # Step 2: Prepare environment FILES FIRST (before Docker starts!)
    # This prevents Docker from creating directories when files don't exist
    prepare_required_files

    # Step 3: Generate Genesis (if needed)
    # WARNING: This will start Docker (make devnet-up), so files must exist first!
    if [ "$DEPLOY_MODE" = "local" ]; then
        log_info "=========================================="
        log_info "Step 3: Generate Genesis"
        log_info "=========================================="
        echo ""

        if ! generate_genesis "$prestate_hash"; then
            log_error "Genesis generation failed"
            exit 1
        fi
    else
        log_info "Skipping Genesis generation (existing mode)"
    fi

    # Step 4: Load environment variables
    load_environment_variables

    # Step 5: Deploy services and run health check
    if ! deploy_and_check "$DG_TYPE" true; then
        log_error "Deployment failed"
        exit 1
    fi
}

##############################################################################
# Execute Main
##############################################################################

main "$@"
