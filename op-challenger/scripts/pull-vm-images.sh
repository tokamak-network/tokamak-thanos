#!/usr/bin/env bash

##############################################################################
# VM 이미지 Pull 스크립트
#
# 사용법:
#   ./pull-vm-images.sh                    # 현재 커밋의 이미지
#   ./pull-vm-images.sh --commit HASH      # 특정 커밋의 이미지
#   ./pull-vm-images.sh --tag v1.0.0       # 특정 태그의 이미지
##############################################################################

set -euo pipefail

# 환경 변수 기본값 설정 (zsh 호환성)
: "${GVM_DEBUG:=}"
: "${GOPATH:=}"
: "${GOROOT:=}"

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $*"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

##############################################################################
# ⚙️  설정 (여기를 수정하세요)
##############################################################################

# 프로젝트 루트
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# 🔧 레지스트리 설정 (필요시 수정)
# - 개인 계정: ghcr.io/zena-park
# - Organization: ghcr.io/tokamak-network
DEFAULT_REGISTRY="ghcr.io/zena-park"
REGISTRY="${VM_REGISTRY:-$DEFAULT_REGISTRY}"

# 이미지 태그 (기본값: 현재 Git 커밋 해시)
IMAGE_TAG="${VM_IMAGE_TAG:-$(git rev-parse --short HEAD 2>/dev/null || echo 'latest')}"

##############################################################################
# 인자 파싱
##############################################################################

while [[ $# -gt 0 ]]; do
    case $1 in
        --commit)
            IMAGE_TAG="$2"
            shift 2
            ;;
        --tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        --registry)
            REGISTRY="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --commit HASH       Use specific Git commit (default: current)"
            echo "  --tag TAG           Use specific tag (default: current commit)"
            echo "  --registry REGISTRY Registry URL (default: ghcr.io/tokamak-network)"
            echo "  -h, --help          Show this help"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Pull current commit images"
            echo "  $0 --commit 47efca11                  # Pull specific commit"
            echo "  $0 --tag v1.0.0                       # Pull tagged version"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

##############################################################################
# 이미지 Pull 및 바이너리 추출
##############################################################################

pull_and_extract() {
    local image_name="$1"
    local target_dir="$2"
    local binary_path="$3"

    local full_image="${REGISTRY}/${image_name}:${IMAGE_TAG}"

    log_info "Pulling: $full_image"

    if ! docker pull "$full_image"; then
        log_error "Failed to pull: $full_image"
        log_warn "이미지가 존재하지 않습니다. 다음을 확인하세요:"
        log_warn "  1. Git 커밋: $IMAGE_TAG"
        log_warn "  2. Registry: $REGISTRY"
        log_warn "  3. 이미지가 빌드되었는지 확인"
        return 1
    fi

    log_info "Extracting binary from image..."

    # 임시 컨테이너 생성
    # scratch 기반 이미지는 명령어가 없으므로 /bin/sh를 시도하고, 실패하면 true 명령어 사용
    local container_id=$(docker create "$full_image" /bin/sh 2>/dev/null || docker create "$full_image" true 2>/dev/null || docker create "$full_image")

    # 바이너리 복사
    mkdir -p "$target_dir"
    if docker cp "${container_id}:${binary_path}" "$target_dir/"; then
        log_success "Extracted: $target_dir/$(basename $binary_path)"
    else
        log_error "Failed to extract: $binary_path"
        docker rm "$container_id" >/dev/null 2>&1
        return 1
    fi

    # 컨테이너 제거
    docker rm "$container_id" >/dev/null 2>&1
}

##############################################################################
# 메인 실행
##############################################################################

main() {
    log_info "=========================================="
    log_info "VM 이미지 Pull"
    log_info "=========================================="
    log_info "Registry: $REGISTRY"
    log_info "Tag/Commit: $IMAGE_TAG"
    echo ""

    local failed=0

    # Cannon VM (builder 이미지는 /app/cannon/bin/cannon에 바이너리 저장)
    if ! pull_and_extract "vm-cannon" "${PROJECT_ROOT}/cannon/bin" "/app/cannon/bin/cannon"; then
        ((failed++))
    fi
    echo ""

    # Asterisc VM (builder 이미지는 /app/asterisc/bin/asterisc에 바이너리 저장)
    if ! pull_and_extract "vm-asterisc" "${PROJECT_ROOT}/asterisc/bin" "/app/asterisc/bin/asterisc"; then
        ((failed++))
    fi

    # Extract Asterisc prestate files (GameType 2)
    log_info "Extracting Asterisc prestate files..."
    local asterisc_image="${REGISTRY}/vm-asterisc:${IMAGE_TAG}"
    local container_id=$(docker create "$asterisc_image" true 2>/dev/null)

    # Extract prestate-proof.json (배포용 - .pre 포맷)
    if docker cp "${container_id}:/app/asterisc/bin/prestate-proof.json" "${PROJECT_ROOT}/asterisc/bin/" 2>/dev/null; then
        log_success "  ✓ prestate-proof.json extracted (deployment format)"

        # Extract hash
        local asterisc_hash=$(cat "${PROJECT_ROOT}/asterisc/bin/prestate-proof.json" | jq -r '.pre' 2>/dev/null || echo "")
        if [ -n "$asterisc_hash" ] && [ "$asterisc_hash" != "null" ]; then
            log_success "  ✓ Asterisc deployment prestate hash: ${asterisc_hash:0:10}...${asterisc_hash: -8}"
        fi
    else
        log_warn "  ⚠️  prestate-proof.json not found"
    fi

    # Extract prestate.json (런타임용 - full VMState)
    if docker cp "${container_id}:/app/asterisc/bin/prestate.json" "${PROJECT_ROOT}/asterisc/bin/" 2>/dev/null; then
        log_success "  ✓ prestate.json extracted (runtime full VMState) ⭐"
    else
        log_warn "  ⚠️  prestate.json not found in image"
        log_warn "      GameType 2 may not work properly at runtime"
    fi

    docker rm "$container_id" >/dev/null 2>&1
    echo ""

    # op-program (builder 이미지는 /app/op-program/bin/op-program에 바이너리 저장)
    if ! pull_and_extract "vm-op-program" "${PROJECT_ROOT}/op-program/bin" "/app/op-program/bin/op-program"; then
        ((failed++))
    fi

    # Extract prestate files (required for contract deployment)
    log_info "Extracting prestate files..."
    local op_prog_image="${REGISTRY}/vm-op-program:${IMAGE_TAG}"
    local container_id=$(docker create "$op_prog_image" true 2>/dev/null)

    # Extract prestate-proof.json (contains absolute pre-state hash)
    if docker cp "${container_id}:/app/op-program/bin/prestate-proof.json" "${PROJECT_ROOT}/op-program/bin/" 2>/dev/null; then
        log_success "  ✓ prestate-proof.json extracted"
    else
        log_warn "  ⚠️  prestate-proof.json not found in image"
    fi

    # Extract prestate.json
    if docker cp "${container_id}:/app/op-program/bin/prestate.json" "${PROJECT_ROOT}/op-program/bin/" 2>/dev/null; then
        log_success "  ✓ prestate.json extracted"
    else
        log_warn "  ⚠️  prestate.json not found in image"
    fi

    docker rm "$container_id" >/dev/null 2>&1
    echo ""

    # kona-client (GameType 3용 - RISC-V 바이너리 + prestate)
    log_info "Pulling kona-client (GameType 3 support)..."

    if pull_and_extract "vm-kona-client" "${PROJECT_ROOT}/bin" "/kona-client-elf"; then
        log_success "  ✓ kona-client downloaded: ${PROJECT_ROOT}/bin/kona-client-elf"

        # Rename to kona-client for consistency
        mv "${PROJECT_ROOT}/bin/kona-client-elf" "${PROJECT_ROOT}/bin/kona-client"
        log_success "  ✓ Renamed to kona-client"

        # Extract kona prestate files from image (built by kona's official system)
        log_info "Extracting kona prestate files from image..."
        local kona_image="${REGISTRY}/vm-kona-client:${IMAGE_TAG}"
        local container_id=$(docker create "$kona_image" true 2>/dev/null)

        # Extract prestate-proof.json (배포용 - .pre 포맷)
        if docker cp "${container_id}:/prestate-proof.json" "${PROJECT_ROOT}/op-program/bin/prestate-kona.json" 2>/dev/null; then
            log_success "  ✓ prestate-kona.json extracted (deployment format)"

            # Extract hash
            local kona_hash=$(cat "${PROJECT_ROOT}/op-program/bin/prestate-kona.json" | jq -r '.pre' 2>/dev/null || echo "")
            if [ -n "$kona_hash" ] && [ "$kona_hash" != "null" ]; then
                log_success "  ✓ Kona deployment prestate hash: ${kona_hash:0:10}...${kona_hash: -8}"
            fi
        else
            log_warn "  ⚠️  prestate-proof.json not found in image"
        fi

        # Extract prestate.bin.gz (런타임용 - full VMState, 압축)
        if docker cp "${container_id}:/prestate.bin.gz" "${PROJECT_ROOT}/bin/prestate.bin.gz" 2>/dev/null; then
            log_success "  ✓ prestate.bin.gz extracted (runtime full VMState) ⭐"
        else
            log_warn "  ⚠️  prestate.bin.gz not found in image"
            log_warn "      GameType 3 may not work properly"
        fi

        docker rm "$container_id" >/dev/null 2>&1

        # prestate 생성을 위해 kona 디렉토리에도 복사
        local kona_dir="${PROJECT_ROOT}/../kona"
        local kona_target_dir="$kona_dir/target/riscv64imac-unknown-none-elf/release-client-lto"

        mkdir -p "$kona_target_dir"
        cp "${PROJECT_ROOT}/bin/kona-client" "$kona_target_dir/kona-client"
        log_success "  ✓ RISC-V binary copied (for prestate generation)"
        log_info "    GameType 3 (AsteriscKona) available ✅"
    else
        log_warn "  ⚠️  kona-client not available (GameType 3 지원 안됨)"
        log_warn "      GameType 0, 1, 2는 정상 작동합니다"
    fi
    echo ""

    # Pull OP Stack Docker images (배포에 필요한 서비스 이미지들)
    log_info "=========================================="
    log_info "Pulling OP Stack Docker Images"
    log_info "=========================================="
    echo ""

    # Default L2 image (from docker-compose)
    local l2_image="${L2_IMAGE:-tokamaknetwork/thanos-op-geth:nightly}"

    local op_stack_images=(
        "${REGISTRY}/op-challenger:${IMAGE_TAG}"
        "${REGISTRY}/op-node:${IMAGE_TAG}"
        "${REGISTRY}/op-proposer:${IMAGE_TAG}"
        "${REGISTRY}/op-batcher:${IMAGE_TAG}"
        "$l2_image"
    )

    local pull_failed=0
    for img in "${op_stack_images[@]}"; do
        log_info "Pulling: $img"
        if docker pull "$img" 2>&1 | grep -q "manifest unknown\|not found"; then
            log_warn "  ⚠️  Image not found: $img (may need to build locally)"
            ((pull_failed++))
        elif docker pull "$img"; then
            log_success "  ✓ Pulled: $img"
        else
            log_warn "  ⚠️  Failed to pull: $img (network issue or not available)"
            ((pull_failed++))
        fi
    done
    echo ""

    if [ $pull_failed -gt 0 ]; then
        log_warn "⚠️  $pull_failed OP Stack images could not be pulled"
        log_warn "   These will be built/pulled during deployment"
    else
        log_success "✅ All OP Stack images pulled successfully!"
    fi
    echo ""

    # 결과
    if [ $failed -eq 0 ]; then
        log_success "=========================================="
        log_success "모든 필수 VM 바이너리 다운로드 완료!"
        log_success "=========================================="
        echo ""
        log_info "다음 단계:"
        log_info "  1. 배포: ./op-challenger/scripts/deploy-modular.sh --dg-type 3 --clean"
        log_info "  2. 모니터링: ./op-challenger/scripts/monitor-challenger.sh"
        return 0
    else
        log_error "=========================================="
        log_error "$failed 개 필수 바이너리 다운로드 실패"
        log_error "=========================================="
        return 1
    fi
}

main "$@"

