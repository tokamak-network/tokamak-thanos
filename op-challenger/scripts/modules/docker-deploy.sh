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

    # Export common variables
    export PROJECT_ROOT
    export DG_TYPE="$dg_type"
    export CHALLENGER_TRACE_TYPE

    # Set Asterisc-specific paths if GameType 2
    if [ "$dg_type" = "2" ]; then
        # Host path for volume mount
        export ASTERISC_BIN="${PROJECT_ROOT}/asterisc/bin"

        # Container path for environment variable (NOT host path!)
        # Docker Compose will mount ASTERISC_BIN to /asterisc, so use container path
        export ASTERISC_PRESTATE="/asterisc/prestate.json"

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
    log_info "Compose file: $COMPOSE_FILE"
    echo ""

    if $DOCKER_COMPOSE -f docker-compose-full.yml up -d; then
        log_success "Services started successfully"
        echo ""
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
