#!/usr/bin/env bash
set -euo pipefail

##############################################################################
# Deployment Cleanup Script
#
# Completely resets the deployment state including:
# - Docker containers, volumes, and networks
# - Incorrectly created directories in .devnet
# - Optionally: environment variables and Genesis files
#
# Reference: /op-challenger/scripts/README.md
##############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DOCKER_COMPOSE_FILE="${SCRIPT_DIR}/docker-compose-full.yml"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Log functions
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
Deployment Cleanup Script

Usage:
    $0 [options]

Options:
    --all-containers    Clean ALL related containers (including devnet, kurtosis, etc.)
    --keep-env          Keep environment variable file (.env)
    --keep-genesis      Keep Genesis files in .devnet
    --help              Show this help message

Default Behavior:
    By default, this script cleans EVERYTHING:
    - Docker containers and volumes
    - Genesis files (.devnet/*.json)
    - Environment variable file (.env)
    - Genesis hash caches

Examples:
    # Clean everything (default)
    $0

    # Clean but keep .env file
    $0 --keep-env

    # Clean but keep Genesis files
    $0 --keep-genesis

    # Clean ALL containers (this script + devnet + kurtosis)
    $0 --all-containers

    # Keep both .env and Genesis
    $0 --keep-env --keep-genesis

EOF
    exit 0
}

# Default options (clean everything by default)
CLEAN_ENV=true
CLEAN_GENESIS=true
CLEAN_ALL_CONTAINERS=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --all)
            CLEAN_ENV=true
            CLEAN_GENESIS=true
            shift
            ;;
        --all-containers)
            CLEAN_ALL_CONTAINERS=true
            shift
            ;;
        --keep-env)
            CLEAN_ENV=false
            shift
            ;;
        --keep-genesis)
            CLEAN_GENESIS=false
            shift
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
# Confirm cleanup
##############################################################################

echo ""
log_warn "=========================================="
log_warn "Deployment Cleanup"
log_warn "=========================================="
echo ""
log_warn "This will clean:"
echo "  - Docker containers (scripts-*)"
if [ "$CLEAN_ALL_CONTAINERS" = true ]; then
    echo "  - ALL related containers (devnet, kurtosis, etc.)"
fi
echo "  - Docker volumes"
echo "  - Docker networks"
if [ "$CLEAN_ENV" = true ]; then
    echo "  - Environment variable file (.env)"
fi
if [ "$CLEAN_GENESIS" = true ]; then
    echo "  - Genesis files (.devnet/*.json)"
fi
echo ""

read -p "Continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    log_info "Cleanup cancelled."
    exit 0
fi

##############################################################################
# 1. Stop and remove Docker resources
##############################################################################

log_info "Stopping and removing Docker containers..."

# Docker Compose version detection
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
else
    log_error "docker-compose or docker compose not found"
    exit 1
fi

cd "$SCRIPT_DIR"

# Stop and remove containers, networks, volumes
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" down -v 2>/dev/null || true

log_success "Docker Compose resources cleaned"

# Clean ALL related containers if requested
if [ "$CLEAN_ALL_CONTAINERS" = true ]; then
    log_info "Cleaning ALL related containers..."

    # Clean ops-bedrock devnet containers
    if [ -f "${PROJECT_ROOT}/ops-bedrock/docker-compose.yml" ]; then
        log_info "Cleaning ops-bedrock devnet..."
        cd "${PROJECT_ROOT}/ops-bedrock"
        $DOCKER_COMPOSE -f docker-compose.yml down -v 2>/dev/null || true
    fi

    # Stop and remove all containers with specific prefixes
    log_info "Removing containers matching patterns..."
    docker ps -a --format '{{.Names}}' | grep -E '^(scripts-|ops-bedrock-|op-|kurtosis-|wait-for-)' | while read container; do
        log_info "  Removing: $container"
        docker rm -f "$container" 2>/dev/null || true
    done

    # Clean dangling volumes
    log_info "Cleaning dangling volumes..."
    docker volume ls -qf dangling=true | while read volume; do
        docker volume rm "$volume" 2>/dev/null || true
    done

    log_success "All related containers cleaned"
fi

##############################################################################
# 2. Clean incorrectly created directories in .devnet
##############################################################################

log_info "Cleaning .devnet directory..."

DEVNET_DIR="${PROJECT_ROOT}/.devnet"

if [ -d "$DEVNET_DIR" ]; then
    # Remove directories that should be files
    for item in genesis-l1.json genesis-l2.json rollup.json; do
        item_path="${DEVNET_DIR}/${item}"
        if [ -d "$item_path" ]; then
            log_warn "Removing incorrectly created directory: $item"
            rm -rf "$item_path"
        fi
    done

    # Optionally clean Genesis files
    if [ "$CLEAN_GENESIS" = true ]; then
        log_warn "Removing all Genesis and allocation files..."
        rm -f "${DEVNET_DIR}"/genesis-*.json
        rm -f "${DEVNET_DIR}"/rollup.json
        rm -f "${DEVNET_DIR}"/allocs-*.json
        rm -f "${DEVNET_DIR}"/addresses.json
        rm -f "${DEVNET_DIR}"/*.json

        # Remove Genesis hash cache files
        log_info "Removing Genesis hash cache files..."
        rm -f "${PROJECT_ROOT}/op-program/bin/.genesis-hash"
        rm -f "${PROJECT_ROOT}/op-program/bin/.rollup-hash"
    fi

    log_success ".devnet directory cleaned"
else
    log_info ".devnet directory does not exist, skipping"
fi

##############################################################################
# 3. Clean environment variable file
##############################################################################

if [ "$CLEAN_ENV" = true ]; then
    ENV_FILE="${PROJECT_ROOT}/.env"
    if [ -f "$ENV_FILE" ]; then
        log_warn "Removing environment variable file: $ENV_FILE"
        rm -f "$ENV_FILE"
        log_success "Environment variable file removed"
    fi
fi

##############################################################################
# 4. Summary
##############################################################################

echo ""
log_success "=========================================="
log_success "Cleanup Completed!"
log_success "=========================================="
echo ""

log_info "Next steps:"
if [ "$CLEAN_ENV" = true ] || [ ! -f "${PROJECT_ROOT}/.env" ]; then
    echo "  1. Generate environment variables:"
    echo "     ${SCRIPT_DIR}/setup-env.sh --mode local"
    echo ""
fi

if [ "$CLEAN_GENESIS" = true ]; then
    echo "  2. Generate Genesis files (required!):"
    echo ""
    log_warn "     ⚠️  Genesis files are created by 'make devnet-up', NOT 'make devnet-allocs'"
    echo ""
    echo "     cd ${PROJECT_ROOT}"
    echo "     make devnet-up        # Start devnet (generates Genesis files)"
    echo "     # Wait 1-2 minutes for completion..."
    echo "     make devnet-down      # Stop devnet"
    echo ""
    echo "     # Verify Genesis files:"
    echo "     ls -la .devnet/genesis*.json .devnet/rollup.json"
    echo ""
fi

echo "  3. Deploy L2 system:"
echo "     ${SCRIPT_DIR}/deploy-full-stack.sh --mode local"
echo ""

exit 0
