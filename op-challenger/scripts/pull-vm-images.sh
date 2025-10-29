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
    echo ""

    # op-program (builder 이미지는 /app/op-program/bin/op-program에 바이너리 저장)
    if ! pull_and_extract "vm-op-program" "${PROJECT_ROOT}/op-program/bin" "/app/op-program/bin/op-program"; then
        ((failed++))
    fi
    echo ""

    # kona-client (GameType 3용 - 선택사항, RISC-V 바이너리)
    log_info "Pulling kona-client (GameType 3 support)..."

    # kona-client는 RISC-V 바이너리를 두 곳에 저장
    # 1. bin/kona-client - deploy가 사용
    # 2. ../kona/target/.../kona-client - prestate 생성용 (같은 파일)
    if pull_and_extract "vm-kona-client" "${PROJECT_ROOT}/bin" "/kona-client"; then
        log_success "  ✓ kona-client downloaded: ${PROJECT_ROOT}/bin/kona-client"

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

    # 결과
    if [ $failed -eq 0 ]; then
        log_success "=========================================="
        log_success "모든 필수 VM 바이너리 다운로드 완료!"
        log_success "=========================================="
        echo ""
        log_info "다음 단계:"
        log_info "  1. Prestate 생성 필요 (컨트랙 배포 후)"
        log_info "  2. 배포: ./op-challenger/scripts/deploy-modular.sh --dg-type 3"
        return 0
    else
        log_error "=========================================="
        log_error "$failed 개 필수 바이너리 다운로드 실패"
        log_error "=========================================="
        return 1
    fi
}

main "$@"

