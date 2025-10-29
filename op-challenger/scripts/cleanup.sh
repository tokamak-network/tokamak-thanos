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
    --all-containers      Clean ALL related containers (including devnet, kurtosis, etc.)
    --rebuild             Remove locally built Docker images and build cache
    --clean-vm-builds     Remove VM binaries (use if VM code changed)
    --clean-pulled-images Remove pulled Docker images from registry (ghcr.io)
    --keep-env            Keep environment variable file (.env)
    --keep-genesis        Keep Genesis files in .devnet
    --help                Show this help message

Default Behavior:
    By default, this script cleans:
    - Docker containers and volumes
    - Genesis files (.devnet/*.json)
    - Environment variable file (.env)
    - Genesis hash caches

    NOT cleaned by default:
    - VM binaries (use --clean-vm-builds to remove)
    - Pulled Docker images from ghcr.io (use --clean-pulled-images to remove)

Examples:
    # Basic cleanup (keeps VM binaries and Docker images)
    $0

    # Clean and force rebuild of local images
    $0 --rebuild

    # Clean VM binaries (if you want to rebuild from source)
    $0 --clean-vm-builds

    # Clean pulled Docker images (if you want to re-pull)
    $0 --clean-pulled-images

    # Clean but keep .env file
    $0 --keep-env

    # Clean but keep Genesis files
    $0 --keep-genesis

    # Clean ALL containers (this script + devnet + kurtosis)
    $0 --all-containers

    # Full clean (everything including VM binaries and pulled images)
    $0 --rebuild --clean-vm-builds --clean-pulled-images

    # Minimal clean (keeps .env, Genesis, VM binaries, Docker images)
    $0 --keep-env --keep-genesis

EOF
    exit 0
}

# Default options
CLEAN_ENV=true
CLEAN_GENESIS=true
CLEAN_ALL_CONTAINERS=false
REBUILD=false
CLEAN_VM_BUILDS=false          # Keep VM binaries by default
CLEAN_PULLED_IMAGES=false      # Keep pulled Docker images by default

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
        --rebuild)
            REBUILD=true
            shift
            ;;
        --clean-vm-builds)
            CLEAN_VM_BUILDS=true
            shift
            ;;
        --clean-pulled-images)
            CLEAN_PULLED_IMAGES=true
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
if [ "$REBUILD" = true ]; then
    echo "  - Locally built Docker images (tokamaknetwork/thanos-*, ops-bedrock-*)"
    echo "  - Dangling images (<none>)"
    echo "  - Docker build cache (all unused cache)"
fi
if [ "$CLEAN_PULLED_IMAGES" = true ]; then
    REGISTRY="${VM_REGISTRY:-ghcr.io/zena-park}"
    echo "  - Pulled Docker images from $REGISTRY"
fi
if [ "$CLEAN_VM_BUILDS" = true ]; then
    echo "  - VM binaries (op-program, Cannon, Asterisc, Kona)"
fi
if [ "$CLEAN_ENV" = true ]; then
    echo "  - Environment variable file (.env)"
fi
if [ "$CLEAN_GENESIS" = true ]; then
    echo "  - Genesis files (.devnet/*.json)"
fi
echo ""
if [ "$CLEAN_PULLED_IMAGES" = false ]; then
    log_info "Keeping pulled Docker images (use --clean-pulled-images to remove)"
fi
if [ "$CLEAN_VM_BUILDS" = false ]; then
    log_info "Keeping VM binaries (use --clean-vm-builds to remove)"
fi
echo ""

read -p "Continue? (Y/n): " -r confirm
confirm=${confirm:-Y}  # Default to Y if empty (just press Enter)
if [[ ! $confirm =~ ^[Yy]$ ]]; then
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
# 3. Clean Docker images and build cache (if --rebuild)
##############################################################################

if [ "$REBUILD" = true ]; then
    log_info "=========================================="
    log_info "Removing Docker images and build cache..."
    log_info "=========================================="
    echo ""

    # Remove locally built images only (keep ghcr.io images from pull-vm-images.sh)
    log_info "Removing locally built images..."
    log_warn "  Keeping ghcr.io images (downloaded via pull-vm-images.sh)"

    # Remove only tokamaknetwork/thanos images (locally built)
    docker images --format '{{.Repository}}:{{.Tag}}' | grep '^tokamaknetwork/thanos-' | while read image; do
        log_info "  Removing image: $image"
        docker rmi "$image" 2>/dev/null || true
    done

    # Remove ops-bedrock images (devnet L1/L2)
    log_info "Removing ops-bedrock-* images..."
    docker images --format '{{.Repository}}:{{.Tag}}' | grep 'ops-bedrock' | while read image; do
        log_info "  Removing image: $image"
        docker rmi "$image" 2>/dev/null || true
    done

    # Remove dangling images (<none>)
    log_info "Removing dangling images (<none>)..."
    docker images -f "dangling=true" -q | while read image_id; do
        if [ -n "$image_id" ]; then
            log_info "  Removing dangling image: $image_id"
            docker rmi "$image_id" 2>/dev/null || true
        fi
    done

    # Remove builder cache for op-challenger and related services
    log_info "Pruning Docker build cache..."
    docker builder prune -f --filter "label=stage=op-challenger-builder" 2>/dev/null || true
    docker builder prune -f --filter "label=stage=op-node-builder" 2>/dev/null || true
    docker builder prune -f --filter "label=stage=op-batcher-builder" 2>/dev/null || true
    docker builder prune -f --filter "label=stage=op-proposer-builder" 2>/dev/null || true

    # Additional aggressive cache pruning (removes all unused build cache)
    log_warn "Removing ALL unused build cache (aggressive)..."
    docker builder prune -af 2>/dev/null || true

    log_success "Docker images and build cache removed"
    log_success "  ✓ tokamaknetwork/thanos-* images (locally built)"
    log_success "  ✓ ops-bedrock-* images (devnet)"
    log_success "  ✓ Dangling images (<none>)"
    log_success "  ✓ All build cache"
    log_warn "  ⚠️  Kept ghcr.io images (use --clean-pulled-images to remove)"
    log_success "Next deployment will rebuild everything from source"
    echo ""
fi

##############################################################################
# 3.5. Clean pulled Docker images from registry (if requested)
##############################################################################

if [ "$CLEAN_PULLED_IMAGES" = true ]; then
    log_info "=========================================="
    log_info "Removing pulled Docker images..."
    log_info "=========================================="
    echo ""

    # Get the registry from environment or default
    REGISTRY="${VM_REGISTRY:-ghcr.io/zena-park}"

    log_warn "Removing images from: $REGISTRY"

    # Remove all images from the registry
    docker images --format '{{.Repository}}:{{.Tag}}' | grep "^${REGISTRY}/" | while read image; do
        log_info "  Removing: $image"
        docker rmi "$image" 2>/dev/null || true
    done

    log_success "Pulled Docker images removed"
    log_success "  ✓ All $REGISTRY/* images deleted"
    log_warn "  ⚠️  Next deployment will re-pull images from registry"
    echo ""
fi

##############################################################################
# 4. Clean VM build artifacts (if enabled)
##############################################################################

if [ "$CLEAN_VM_BUILDS" = true ]; then
    log_info "=========================================="
    log_info "Removing VM build artifacts..."
    log_info "=========================================="
    echo ""

    files_removed=0

    # Remove op-program binaries
    if [ -f "${PROJECT_ROOT}/op-program/bin/op-program" ]; then
        log_info "Removing op-program binary..."
        rm -f "${PROJECT_ROOT}/op-program/bin/op-program" && ((files_removed++))
    fi

    # Remove Cannon VM and prestate
    if [ -f "${PROJECT_ROOT}/cannon/bin/cannon" ]; then
        log_info "Removing Cannon binary..."
        rm -f "${PROJECT_ROOT}/cannon/bin/cannon" && ((files_removed++))
    fi
    if [ -f "${PROJECT_ROOT}/op-program/bin/prestate-proof.json" ]; then
        log_info "Removing Cannon prestate..."
        rm -f "${PROJECT_ROOT}/op-program/bin/prestate-proof.json" && ((files_removed++))
        rm -f "${PROJECT_ROOT}/op-program/bin/prestate.json" 2>/dev/null
    fi

    # Remove Asterisc VM and prestate
    if [ -f "${PROJECT_ROOT}/asterisc/bin/asterisc" ]; then
        log_info "Removing Asterisc binary..."
        rm -f "${PROJECT_ROOT}/asterisc/bin/asterisc" && ((files_removed++))
    fi
    if [ -f "${PROJECT_ROOT}/asterisc/bin/prestate-proof.json" ]; then
        log_info "Removing Asterisc prestate..."
        rm -f "${PROJECT_ROOT}/asterisc/bin/prestate-proof.json" && ((files_removed++))
        rm -f "${PROJECT_ROOT}/asterisc/bin/prestate.json" 2>/dev/null
        rm -f "${PROJECT_ROOT}/asterisc/bin/meta.json" 2>/dev/null
    fi

    # Remove Asterisc source directory
    if [ -d "${PROJECT_ROOT}/.asterisc-src" ]; then
        log_info "Removing Asterisc source directory..."
        rm -rf "${PROJECT_ROOT}/.asterisc-src" && ((files_removed++))
    fi

    # Remove Kona client binary and prestate
    if [ -f "${PROJECT_ROOT}/bin/kona-client" ]; then
        log_info "Removing Kona client binary..."
        rm -f "${PROJECT_ROOT}/bin/kona-client" && ((files_removed++))
    fi
    if [ -f "${PROJECT_ROOT}/op-program/bin/prestate-kona.json" ]; then
        log_info "Removing Kona prestate..."
        rm -f "${PROJECT_ROOT}/op-program/bin/prestate-kona.json" && ((files_removed++))
    fi

    # Remove Kona build cache (external project)
    kona_dir="${PROJECT_ROOT}/../kona"
    if [ -d "$kona_dir/target" ]; then
        log_info "Removing Kona build cache (external project)..."
        rm -rf "$kona_dir/target" && ((files_removed++))
    fi
    if [ -d "$kona_dir/bin-docker" ]; then
        rm -rf "$kona_dir/bin-docker" 2>/dev/null
    fi

    if [ $files_removed -gt 0 ]; then
        log_success "VM build artifacts removed"
        log_success "  ✓ Removed $files_removed item(s)"
        log_success "  ✓ Next deployment will rebuild all VMs from source"
    else
        log_info "No VM build artifacts to remove"
    fi
    echo ""
fi

##############################################################################
# 5. Clean environment variable file
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
# 6. Summary
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

if [ "$REBUILD" = true ]; then
    log_info "=========================================="
    log_info "Rebuild Mode Enabled"
    log_info "=========================================="
    echo ""
    log_success "✅ Docker images and build cache removed"
    log_success "✅ Next deployment will rebuild all services from source"
    echo ""
    log_warn "Note: First build after --rebuild will take longer (5-10 minutes)"
    log_warn "      but ensures all code changes are properly compiled"
    echo ""
fi

if [ "$CLEAN_VM_BUILDS" = true ]; then
    log_info "=========================================="
    log_info "VM Builds Cleaned"
    log_info "=========================================="
    echo ""
    log_success "✅ VM build artifacts removed (op-program, Cannon, Asterisc, Kona)"
    log_success "✅ Next deployment will rebuild all VMs from source"
    echo ""
    log_warn "Note: First VM build will take longer (10-30 minutes total)"
    log_warn "      - op-program: ~2-5 min"
    log_warn "      - Cannon: ~3-5 min"
    log_warn "      - Asterisc: ~5-10 min"
    log_warn "      - Kona: ~5-10 min"
    log_warn "      This ensures latest code from external projects (Kona, Asterisc)"
    echo ""
fi

exit 0
