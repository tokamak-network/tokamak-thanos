#!/usr/bin/env bash

##############################################################################
# Cleanup Script
#
# 배포된 L2 시스템을 정리하고 초기화합니다.
##############################################################################

set -euo pipefail

# Source common library
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${SCRIPT_DIR}/lib/common.sh"
source "${SCRIPT_DIR}/modules/cleanup.sh"

##############################################################################
# 추가 정리 함수
##############################################################################

cleanup_vm_builds() {
    log_info "=========================================="
    log_info "Cleanup VM Binaries"
    log_info "=========================================="
    echo ""

    local files_removed=0

    # Cannon
    if [ -f "${PROJECT_ROOT}/cannon/bin/cannon" ]; then
        rm -f "${PROJECT_ROOT}/cannon/bin/cannon" && ((files_removed++))
        log_info "  ✓ Removed cannon binary"
    fi

    # Asterisc
    if [ -d "${PROJECT_ROOT}/asterisc/bin" ]; then
        rm -rf "${PROJECT_ROOT}/asterisc/bin" && ((files_removed++))
        log_info "  ✓ Removed asterisc binaries"
    fi

    # op-program
    if [ -d "${PROJECT_ROOT}/op-program/bin" ]; then
        rm -rf "${PROJECT_ROOT}/op-program/bin" && ((files_removed++))
        log_info "  ✓ Removed op-program binaries"
    fi

    # kona-client
    if [ -f "${PROJECT_ROOT}/bin/kona-client" ]; then
        rm -f "${PROJECT_ROOT}/bin/kona-client" && ((files_removed++))
        log_info "  ✓ Removed kona-client binary"
    fi

    if [ -f "${PROJECT_ROOT}/bin/prestate.bin.gz" ]; then
        rm -f "${PROJECT_ROOT}/bin/prestate.bin.gz" && ((files_removed++))
        log_info "  ✓ Removed kona runtime prestate"
    fi

    if [ $files_removed -gt 0 ]; then
        log_success "Removed $files_removed VM binary/prestate file(s)"
    else
        log_info "No VM binaries to remove"
    fi

    echo ""
    return 0
}

cleanup_pulled_images() {
    log_info "=========================================="
    log_info "Cleanup Pulled Docker Images"
    log_info "=========================================="
    echo ""

    local images_removed=0

    # ghcr.io 레지스트리에서 pull한 이미지들
    local registry_images=(
        "ghcr.io/zena-park/vm-cannon"
        "ghcr.io/zena-park/vm-asterisc"
        "ghcr.io/zena-park/vm-op-program"
        "ghcr.io/zena-park/vm-kona-client"
        "ghcr.io/zena-park/op-challenger"
        "ghcr.io/zena-park/op-node"
        "ghcr.io/zena-park/op-batcher"
        "ghcr.io/zena-park/op-proposer"
    )

    for img_pattern in "${registry_images[@]}"; do
        # 모든 태그 삭제
        local found_images=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep "^${img_pattern}:" || true)
        if [ -n "$found_images" ]; then
            echo "$found_images" | xargs -r docker rmi -f 2>/dev/null && ((images_removed++))
            log_info "  ✓ Removed: $img_pattern (all tags)"
        fi
    done

    if [ $images_removed -gt 0 ]; then
        log_success "Removed $images_removed pulled image(s)"
    else
        log_info "No pulled images to remove"
    fi

    echo ""
    return 0
}

cleanup_local_images() {
    log_info "=========================================="
    log_info "Cleanup Local Build Images"
    log_info "=========================================="
    echo ""

    local images_removed=0

    # 로컬 빌드 이미지들
    local local_images=(
        "scripts-l1"
        "scripts-challenger-l2"
        "scripts-sequencer-l2"
        "tokamaknetwork/thanos-op-geth"
        "ops-bedrock-l1"
        "ops-bedrock-l2"
    )

    for img_pattern in "${local_images[@]}"; do
        local found_images=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep "^${img_pattern}" || true)
        if [ -n "$found_images" ]; then
            echo "$found_images" | xargs -r docker rmi -f 2>/dev/null && ((images_removed++))
            log_info "  ✓ Removed: $img_pattern"
        fi
    done

    # Dangling images (<none>)
    local dangling=$(docker images -f "dangling=true" -q)
    if [ -n "$dangling" ]; then
        echo "$dangling" | xargs -r docker rmi -f 2>/dev/null && ((images_removed++))
        log_info "  ✓ Removed dangling images"
    fi

    if [ $images_removed -gt 0 ]; then
        log_success "Removed $images_removed local build image(s)"
    else
        log_info "No local images to remove"
    fi

    echo ""
    return 0
}

cleanup_build_cache() {
    log_info "Pruning Docker build cache..."

    docker builder prune -af 2>/dev/null || docker buildx prune -af 2>/dev/null || true

    log_success "Docker build cache pruned"
    echo ""
    return 0
}

cleanup_all_containers() {
    log_info "=========================================="
    log_info "Cleanup ALL Related Containers"
    log_info "=========================================="
    echo ""

    # scripts-, ops-bedrock-, devnet-, kurtosis- 등 모든 관련 컨테이너
    local all_containers=$(docker ps -a --format '{{.Names}}' | grep -E "^(scripts-|ops-bedrock-|devnet-|kurtosis-)" || true)

    if [ -n "$all_containers" ]; then
        echo "$all_containers" | xargs -r docker stop 2>/dev/null || true
        echo "$all_containers" | xargs -r docker rm 2>/dev/null || true
        log_success "All related containers removed"
    else
        log_info "No related containers found"
    fi

    echo ""
    return 0
}

cleanup_env_file() {
    log_info "Removing .env file..."

    if [ -f "${PROJECT_ROOT}/.env" ]; then
        rm -f "${PROJECT_ROOT}/.env"
        log_success ".env file removed"
    else
        log_info "No .env file found"
    fi

    echo ""
    return 0
}

cleanup_genesis_files() {
    log_info "Removing genesis files..."

    local devnet_dir="${PROJECT_ROOT}/.devnet"
    local files_removed=0

    if [ -d "$devnet_dir" ]; then
        for item in genesis-l1.json genesis-l2.json rollup.json addresses.json allocs*.json deploy-config.json; do
            if [ -f "$devnet_dir/$item" ]; then
                rm -f "$devnet_dir/$item" && ((files_removed++))
            fi
        done
    fi

    if [ $files_removed -gt 0 ]; then
        log_success "Removed $files_removed genesis file(s)"
    else
        log_info "No genesis files to remove"
    fi

    echo ""
    return 0
}

##############################################################################
# CLI Interface
##############################################################################

# 기본값 (README와 동일하게 보호 모드)
CLEANUP_ALL_CONTAINERS=false
REBUILD_MODE=false
KEEP_ENV=false
KEEP_GENESIS=false
CLEAN_VM_BUILDS=false
CLEAN_PULLED_IMAGES=false

# 인자 파싱
while [[ $# -gt 0 ]]; do
    case $1 in
        --all-containers)
            CLEANUP_ALL_CONTAINERS=true
            shift
            ;;
        --rebuild)
            REBUILD_MODE=true
            shift
            ;;
        --keep-env)
            KEEP_ENV=true
            shift
            ;;
        --keep-genesis)
            KEEP_GENESIS=true
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
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --all-containers       Clean all related containers (devnet, kurtosis, etc.)"
            echo "  --rebuild              Remove Docker images and build cache (force full rebuild)"
            echo "  --keep-env             Keep environment variable file (.env)"
            echo "  --keep-genesis         Keep genesis files"
            echo "  --clean-vm-builds      Remove VM binaries (default: protected)"
            echo "  --clean-pulled-images  Remove pulled images from registry (default: protected)"
            echo "  --help                 Show this help"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Basic cleanup (protects VM & images)"
            echo "  $0 --rebuild                          # Remove local build images"
            echo "  $0 --clean-vm-builds                  # Also remove VM binaries"
            echo "  $0 --clean-pulled-images              # Also remove pulled images"
            echo "  $0 --clean-vm-builds --clean-pulled-images --all-containers --rebuild"
            echo "                                        # Full cleanup"
            echo "  $0 --keep-env                         # Keep .env file"
            echo "  $0 --keep-genesis                     # Keep genesis files"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

##############################################################################
# 메인 실행
##############################################################################

init_logging

log_warn "=========================================="
log_warn "Deployment Cleanup"
log_warn "=========================================="
echo ""

log_warn "This will clean:"
echo "  - Docker containers (scripts-*)"
[ "$CLEANUP_ALL_CONTAINERS" = true ] && echo "  - ALL related containers (devnet, kurtosis, etc.)"
echo "  - Docker volumes"
echo "  - Docker networks"
[ "$REBUILD_MODE" = true ] && echo "  - Locally built Docker images (tokamaknetwork/thanos-*, ops-bedrock-*)"
[ "$REBUILD_MODE" = true ] && echo "  - Dangling images (<none>)"
[ "$REBUILD_MODE" = true ] && echo "  - Docker build cache (all unused cache)"
[ "$CLEAN_PULLED_IMAGES" = true ] && echo "  - Pulled Docker images from ghcr.io/zena-park"
[ "$CLEAN_VM_BUILDS" = true ] && echo "  - VM binaries (op-program, Cannon, Asterisc, Kona)"
[ "$KEEP_ENV" = false ] && echo "  - Environment variable file (.env)"
[ "$KEEP_GENESIS" = false ] && echo "  - Genesis files (.devnet/*.json)"
echo ""

read -p "Continue? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_info "Cleanup cancelled"
    exit 0
fi

echo ""

# 1. 컨테이너 정리
if [ "$CLEANUP_ALL_CONTAINERS" = true ]; then
    cleanup_all_containers
else
    cleanup_containers
fi

# 2. 볼륨 정리
cleanup_volumes

# 3. .devnet 파일 정리
cleanup_devnet_files

# 4. Genesis 파일 정리 (선택적)
if [ "$KEEP_GENESIS" = false ]; then
    cleanup_genesis_files
fi

# 5. .env 파일 정리 (선택적)
if [ "$KEEP_ENV" = false ]; then
    cleanup_env_file
fi

# 6. VM 바이너리 정리 (선택적)
if [ "$CLEAN_VM_BUILDS" = true ]; then
    cleanup_vm_builds
fi

# 7. Pulled 이미지 정리 (선택적)
if [ "$CLEAN_PULLED_IMAGES" = true ]; then
    cleanup_pulled_images
fi

# 8. 로컬 빌드 이미지 및 캐시 정리 (선택적)
if [ "$REBUILD_MODE" = true ]; then
    cleanup_local_images
    cleanup_build_cache
fi

echo ""
log_success "=========================================="
log_success "Cleanup Complete!"
log_success "=========================================="
echo ""

# 다음 단계 안내
if [ "$CLEAN_VM_BUILDS" = true ]; then
    log_info "Next steps:"
    log_info "  1. Download VM binaries:"
    log_info "     ./op-challenger/scripts/pull-vm-images.sh --tag latest"
    log_info "  2. Deploy:"
    log_info "     ./op-challenger/scripts/deploy-modular.sh --dg-type 3"
elif [ "$CLEAN_PULLED_IMAGES" = true ]; then
    log_info "Next steps:"
    log_info "  1. Re-download images:"
    log_info "     ./op-challenger/scripts/pull-vm-images.sh --tag latest"
    log_info "  2. Deploy:"
    log_info "     ./op-challenger/scripts/deploy-modular.sh --dg-type 3"
else
    log_info "Next steps:"
    log_info "  Deploy: ./op-challenger/scripts/deploy-modular.sh --dg-type 3"
fi
