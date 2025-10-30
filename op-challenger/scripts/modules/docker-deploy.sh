#!/usr/bin/env bash

##############################################################################
# Docker Deploy Module
#
# This module handles Docker Compose service deployment and health checking.
##############################################################################

set -euo pipefail

# Source common library
_MODULE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${_MODULE_DIR}/../lib/common.sh"

##############################################################################
# Configuration
##############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="${SCRIPT_DIR}/docker-compose-full.yml"

# Docker compose command - check which version is available
if docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    log_error "Neither 'docker compose' nor 'docker-compose' is available"
    return 1
fi

##############################################################################
# Docker Deployment Functions
##############################################################################

# get_sequencer_enr() - DISABLED
# ENR-based P2P discovery doesn't work reliably in Docker due to missing IP addresses
# Challenger will sync from L1 batch data instead, which is more reliable
#
# get_sequencer_enr() {
#     local max_attempts=30
#     local attempt=0
#     log_info "Waiting for sequencer op-node to be ready..." >&2
#     ...
# }

export_docker_env_vars() {
    local dg_type="${1:-0}"

    log_info "Exporting environment variables for Docker Compose..."

    # Load contract addresses first (needed before docker-compose)
    source "${SCRIPT_DIR}/modules/env-setup.sh"
    load_contract_addresses || log_warn "Failed to load contract addresses"
    load_private_keys || log_warn "Failed to load private keys"

    # Export common variables
    export PROJECT_ROOT
    export DG_TYPE="$dg_type"

    # Set VM_REGISTRY and IMAGE_TAG for docker-compose to use pre-built images
    export VM_REGISTRY="${VM_REGISTRY:-ghcr.io/zena-park}"
    export IMAGE_TAG="${IMAGE_TAG:-latest}"

    # Set L2_IMAGE for L2 geth (Dockerfile.l2 requires this)
    export L2_IMAGE="${L2_IMAGE:-tokamaknetwork/thanos-op-geth:nightly}"

    # Set CHALLENGER_TRACE_TYPE to support ALL trace types (Independent Verifier)
    # Set SEQUENCER_CHALLENGER_TRACE_TYPE to support ONLY proposer's type (For Sequencer)
    # The order matters - put the proposer's game type first for priority
    case "$dg_type" in
        0|1|254|255)
            # Cannon-based games
            export SEQUENCER_CHALLENGER_TRACE_TYPE="cannon"  # ⭐ 시퀀서: Proposer 타입만
            export CHALLENGER_TRACE_TYPE="cannon,asterisc,asterisc-kona"  # 독립: 모든 타입
            log_info "Proposer GameType $dg_type (Cannon)"
            ;;
        2)
            # Asterisc game
            export SEQUENCER_CHALLENGER_TRACE_TYPE="asterisc"  # ⭐ 시퀀서: Proposer 타입만
            export CHALLENGER_TRACE_TYPE="asterisc,cannon,asterisc-kona"  # 독립: 모든 타입
            log_info "Proposer GameType $dg_type (Asterisc)"
            ;;
        3)
            # AsteriscKona game
            export SEQUENCER_CHALLENGER_TRACE_TYPE="asterisc-kona"  # ⭐ 시퀀서: Proposer 타입만
            export CHALLENGER_TRACE_TYPE="asterisc-kona,asterisc,cannon"  # 독립: 모든 타입
            log_info "Proposer GameType $dg_type (AsteriscKona)"
            ;;
        *)
            # Default: support all with cannon priority
            export SEQUENCER_CHALLENGER_TRACE_TYPE="cannon"
            export CHALLENGER_TRACE_TYPE="cannon,asterisc,asterisc-kona"
            log_info "Default GameType"
            ;;
    esac

    log_info "Challenger configuration:"
    log_info "  🔹 sequencer-challenger (For Sequencer):"
    log_info "     - Trace type: $SEQUENCER_CHALLENGER_TRACE_TYPE (Proposer's GameType only)"
    log_info "     - Role: Close games for user withdrawals"
    log_info "  🔹 op-challenger (Independent Verifier):"
    log_info "     - Trace types: $CHALLENGER_TRACE_TYPE (All GameTypes)"
    log_info "     - Role: Decentralized verification and challenge"

    # Debug: Verify export worked
    log_info "DEBUG: Exported SEQUENCER_CHALLENGER_TRACE_TYPE=$SEQUENCER_CHALLENGER_TRACE_TYPE"
    log_info "DEBUG: Exported CHALLENGER_TRACE_TYPE=$CHALLENGER_TRACE_TYPE"
    log_info "DEBUG: Exported DG_TYPE=$DG_TYPE"

    # Set Asterisc-specific paths if GameType 2
    if [ "$dg_type" = "2" ]; then
        # Host path for volume mount
        export ASTERISC_BIN="${PROJECT_ROOT}/asterisc/bin"

        # Container path for environment variable (NOT host path!)
        # Docker Compose will mount ASTERISC_BIN to /asterisc, so use container path
        export ASTERISC_PRESTATE="/asterisc/prestate-proof.json"

        log_info "GameType 2: Asterisc paths configured"
        log_info "  Host: ${ASTERISC_BIN}"
        log_info "  Container: ${ASTERISC_PRESTATE}"
    fi

    # Note: SEQUENCER_P2P_ENR is not set here
    # Challenger will sync from L1 batch data, not P2P
    # This is more reliable in Docker environments

    log_success "Environment variables exported"
    return 0
}

deploy_docker_services() {
    local dg_type="${1:-0}"

    log_info "=========================================="
    log_info "Deploy Docker Services"
    log_info "=========================================="
    echo ""

    # Export environment variables (without ENR first)
    export_docker_env_vars "$dg_type"

    # Change to script directory for docker-compose context
    cd "${SCRIPT_DIR}"

    log_info "Starting Docker services..."
    log_info "Compose file: docker-compose-full.yml"
    echo ""

    # Debug: Check environment variables before docker-compose
    log_info "DEBUG: Environment variables before docker-compose:"
    log_info "  VM_REGISTRY=$VM_REGISTRY"
    log_info "  IMAGE_TAG=$IMAGE_TAG"
    log_info "  L2_IMAGE=$L2_IMAGE"
    log_info "  CHALLENGER_TRACE_TYPE=$CHALLENGER_TRACE_TYPE"
    log_info "  DG_TYPE=$DG_TYPE"
    log_info "  GAME_FACTORY_ADDRESS=$GAME_FACTORY_ADDRESS"
    log_info "  PROJECT_ROOT=$PROJECT_ROOT"
    echo ""

    # Pre-build l1 and l2 images (not available in registry)
    log_info "Building l1 and l2 images (local only)..."
    # IMPORTANT: Disable BuildKit to avoid "file already closed" error
    # Docker Compose v2.39+ tries to read bake definitions from stdin with BuildKit
    # This causes issues when stdin is not available in script environment
    export DOCKER_BUILDKIT=0
    export COMPOSE_DOCKER_CLI_BUILD=0

    # Build images using legacy builder (more stable for script execution)
    if $DOCKER_COMPOSE -f docker-compose-full.yml build l1 challenger-l2 sequencer-l2 2>&1; then
        log_success "  ✓ l1 and l2 images built"
    else
        log_error "Failed to build l1/l2 images"
        return 1
    fi
    echo ""

    # Use --no-build for OP Stack services (use pre-built images from registry)
    log_info "Starting all services (using registry images for OP Stack)..."
    if $DOCKER_COMPOSE -f docker-compose-full.yml up -d --no-build 2>&1; then
        log_success "Services started successfully"
        echo ""

        # Debug: Verify challengers received correct environment variables
        log_info "DEBUG: Verifying challenger environment variables..."

        # Check sequencer-challenger
        if docker inspect scripts-sequencer-challenger-1 --format='{{range .Config.Env}}{{println .}}{{end}}' 2>/dev/null | grep "OP_CHALLENGER_TRACE_TYPE" > /dev/null; then
            local seq_trace=$(docker inspect scripts-sequencer-challenger-1 --format='{{range .Config.Env}}{{println .}}{{end}}' 2>/dev/null | grep "OP_CHALLENGER_TRACE_TYPE" | cut -d= -f2)
            log_info "  sequencer-challenger TRACE_TYPE: $seq_trace"
        else
            log_warn "  Could not verify sequencer-challenger TRACE_TYPE"
        fi

        # Check op-challenger (독립 검증자)
        if docker inspect scripts-op-challenger-1 --format='{{range .Config.Env}}{{println .}}{{end}}' 2>/dev/null | grep "OP_CHALLENGER_TRACE_TYPE" > /dev/null; then
            local ind_trace=$(docker inspect scripts-op-challenger-1 --format='{{range .Config.Env}}{{println .}}{{end}}' 2>/dev/null | grep "OP_CHALLENGER_TRACE_TYPE" | cut -d= -f2)
            log_info "  op-challenger TRACE_TYPE: $ind_trace"
        else
            log_warn "  Could not verify op-challenger TRACE_TYPE"
        fi
    else
        log_error "Failed to start services"
        return 1
    fi

    # Note: Challenger will sync from L1 batch data (not P2P)
    # In Docker environment, ENR-based P2P discovery is unreliable due to dynamic IPs
    # The challenger op-node will derive L2 state from L1 batches instead

    return 0
}

run_health_check() {
    local wait_time="${1:-5}"

    log_info "=========================================="
    log_info "Health Check"
    log_info "=========================================="
    echo ""

    log_info "Waiting for services to initialize..."
    sleep "$wait_time"

    local health_check_script="${SCRIPT_DIR}/health-check.sh"

    if [ -f "$health_check_script" ]; then
        log_info "Running health check..."
        echo ""
        bash "$health_check_script"
        return $?
    else
        log_warn "Health check script not found: $health_check_script"
        log_info "Skipping health check"
        return 0
    fi
}

show_deployment_info() {
    echo ""
    log_success "=========================================="
    log_success "Deployment Complete!"
    log_success "=========================================="
    echo ""
    log_info "Services are now running. Check status with:"
    log_info "  docker ps"
    echo ""
    log_info "View logs with:"
    log_info "  $DOCKER_COMPOSE -f ${COMPOSE_FILE} logs -f [service-name]"
    echo ""
    log_info "To monitor the challenger:"
    log_info "  ${SCRIPT_DIR}/monitor-challenger.sh"
    echo ""
    log_info "To re-run health check:"
    log_info "  ${SCRIPT_DIR}/health-check.sh"
    echo ""
}

deploy_and_check() {
    local dg_type="${1:-0}"
    local run_health="${2:-true}"

    # Deploy services
    if ! deploy_docker_services "$dg_type"; then
        log_error "Docker service deployment failed"
        return 1
    fi

    # Run health check if requested
    if [ "$run_health" = "true" ]; then
        run_health_check
    fi

    # Show deployment info
    show_deployment_info

    return 0
}

##############################################################################
# CLI Interface (if run directly)
##############################################################################

if [ "${BASH_SOURCE[0]}" -ef "$0" ]; then
    # Script is being run directly, not sourced

    DG_TYPE="${DG_TYPE:-0}"
    RUN_HEALTH_CHECK=true

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --dg-type)
                DG_TYPE="$2"
                shift 2
                ;;
            --no-health-check)
                RUN_HEALTH_CHECK=false
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --dg-type 0|1|2|254|255    Dispute game type (default: 0)"
                echo "  --no-health-check           Skip health check"
                echo "  --help                      Show this help"
                echo ""
                echo "Examples:"
                echo "  $0                      # Deploy with default settings"
                echo "  $0 --dg-type 2          # Deploy for GameType 2"
                echo "  $0 --no-health-check    # Deploy without health check"
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

    # Run deployment
    if deploy_and_check "$DG_TYPE" "$RUN_HEALTH_CHECK"; then
        exit 0
    else
        exit 1
    fi
fi
