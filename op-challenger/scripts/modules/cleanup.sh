#!/usr/bin/env bash

##############################################################################
# Cleanup Module
#
# This module handles cleanup of existing Docker containers and volumes
# to prevent conflicts during deployment.
##############################################################################

set -euo pipefail

# Source common library
_MODULE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${_MODULE_DIR}/../lib/common.sh"

##############################################################################
# Cleanup Functions
##############################################################################

cleanup_containers() {
    log_info "=========================================="
    log_info "Cleanup Existing Containers"
    log_info "=========================================="
    echo ""

    log_info "Checking for existing containers that might conflict..."

    # Check if any containers are using our ports
    local cleanup_needed=false
    if docker ps -a --format '{{.Names}}' | grep -E "^(scripts-|ops-bedrock-)" > /dev/null 2>&1; then
        cleanup_needed=true
        log_warn "Found existing containers that might conflict"
        log_info "Stopping and removing conflicting containers..."

        # Stop containers
        docker ps -a --format '{{.Names}}' | grep -E "^(scripts-|ops-bedrock-)" | xargs -r docker stop 2>/dev/null || true

        # Remove containers
        docker ps -a --format '{{.Names}}' | grep -E "^(scripts-|ops-bedrock-)" | xargs -r docker rm 2>/dev/null || true

        log_success "Cleanup completed"
    else
        log_info "No conflicting containers found"
    fi

    echo ""
    return 0
}

cleanup_volumes() {
    log_info "Cleaning up Docker volumes..."

    local volumes_removed=0
    for vol in scripts_l1_data scripts_sequencer_l2_data scripts_challenger_l2_data \
               scripts_sequencer_safedb_data scripts_challenger_safedb_data \
               scripts_challenger_data scripts_op_log \
               ops-bedrock_l1_data ops-bedrock_l2_data ops-bedrock_safedb_data \
               ops-bedrock_challenger_data ops-bedrock_op_log; do
        if docker volume rm "$vol" 2>/dev/null; then
            volumes_removed=$((volumes_removed + 1))
        fi
    done

    if [ $volumes_removed -gt 0 ]; then
        log_success "Removed $volumes_removed Docker volume(s)"
    else
        log_info "No volumes to remove"
    fi

    return 0
}

cleanup_devnet_files() {
    log_info "Cleaning up .devnet directory..."

    local devnet_dir="${PROJECT_ROOT}/.devnet"

    if [ ! -d "$devnet_dir" ]; then
        log_info "No .devnet directory found"
        return 0
    fi

    local files_removed=0

    # Remove files/directories that might be incorrectly created by Docker
    for item in jwt-secret.txt p2p-node-key.txt p2p-challenger-key.txt; do
        local path="$devnet_dir/$item"
        if [ -e "$path" ]; then
            if [ -d "$path" ]; then
                log_warn "Removing directory (should be file): $item"
                rm -rf "$path" && ((files_removed++))
            elif [ -f "$path" ]; then
                log_info "Removing file: $item"
                rm -f "$path" && ((files_removed++))
            fi
        fi
    done

    # Optionally remove genesis files (if you want to regenerate)
    # Uncomment these lines if you want to also clean genesis files
    # for item in genesis-l1.json genesis-l2.json rollup.json addresses.json allocs*.json; do
    #     rm -f "$devnet_dir/$item" 2>/dev/null && ((files_removed++))
    # done

    if [ $files_removed -gt 0 ]; then
        log_success "Removed $files_removed .devnet file(s)"
    else
        log_info "No .devnet files to remove"
    fi

    return 0
}

cleanup_all() {
    cleanup_containers
    cleanup_volumes
    cleanup_devnet_files
    log_success "Full cleanup completed"
}

##############################################################################
# CLI Interface (if run directly)
##############################################################################

if [ "${BASH_SOURCE[0]}" -ef "$0" ]; then
    # Script is being run directly, not sourced

    CLEANUP_VOLUMES=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --volumes)
                CLEANUP_VOLUMES=true
                shift
                ;;
            --all)
                CLEANUP_VOLUMES=true
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --volumes    Also remove Docker volumes"
                echo "  --all        Remove containers and volumes"
                echo "  --help       Show this help"
                echo ""
                echo "Examples:"
                echo "  $0              # Cleanup containers only"
                echo "  $0 --volumes    # Cleanup containers and volumes"
                echo "  $0 --all        # Full cleanup"
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

    # Run cleanup
    if [ "$CLEANUP_VOLUMES" = true ]; then
        cleanup_all
    else
        cleanup_containers
    fi

    exit 0
fi
