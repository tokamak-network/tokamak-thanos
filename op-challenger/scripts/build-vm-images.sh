#!/usr/bin/env bash

##############################################################################
# VM 이미지 빌드 및 푸시 스크립트
#
# 사용법:
#   ./build-vm-images.sh                    # 빌드만
#   ./build-vm-images.sh --push             # 빌드 + 푸시
#   ./build-vm-images.sh --registry ghcr.io # 다른 레지스트리
##############################################################################

set -euo pipefail

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수
log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $*"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

##############################################################################
# ⚙️  설정 (여기를 수정하세요)
##############################################################################

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 🔧 레지스트리 설정 (필요시 수정)
# - 개인 계정: ghcr.io/zena-park
# - Organization: ghcr.io/tokamak-network
DEFAULT_REGISTRY="ghcr.io/zena-park"
REGISTRY="${VM_REGISTRY:-$DEFAULT_REGISTRY}"

# Git 정보 자동 추출
GIT_COMMIT_SHORT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_COMMIT_FULL=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")

# 이미지 태그 결정
if [ -n "${VM_IMAGE_TAG:-}" ]; then
    # 사용자가 명시적으로 지정
    IMAGE_TAG="$VM_IMAGE_TAG"
elif [ -n "$GIT_TAG" ]; then
    # Git 태그가 있으면 사용 (예: v1.0.0)
    IMAGE_TAG="$GIT_TAG"
else
    # 커밋 해시 사용 (같은 커밋 = 같은 이미지 = 덮어쓰기)
    IMAGE_TAG="$GIT_COMMIT_SHORT"
fi

PUSH_IMAGES=false

##############################################################################
# 인자 파싱
##############################################################################

while [[ $# -gt 0 ]]; do
    case $1 in
        --push)
            PUSH_IMAGES=true
            shift
            ;;
        --registry)
            REGISTRY="$2"
            shift 2
            ;;
        --tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --push              Push images to registry after build"
            echo "  --registry REGISTRY Registry URL (default: ghcr.io/zena-park)"
            echo "  --tag TAG           Override image tag (default: auto-detect)"
            echo "  -h, --help          Show this help"
            echo ""
            echo "Tagging Strategy:"
            echo "  - Auto-detects Git tag or commit hash"
            echo "  - If Git tag exists: uses tag (e.g., v1.0.0)"
            echo "  - Otherwise: uses commit hash (e.g., 47efca11)"
            echo "  - Also tags as 'latest' when appropriate"
            echo ""
            echo "Examples:"
            echo "  $0                           # Build with auto-detected tag"
            echo "  $0 --tag v1.0.0             # Build with specific tag"
            echo "  $0 --push                   # Build and push to registry"
            echo "  $0 --tag stable --push      # Build, tag as 'stable', and push"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

##############################################################################
# GitHub Container Registry 로그인
##############################################################################

login_to_ghcr() {
    # 이미 로그인되어 있는지 확인
    if cat ~/.docker/config.json 2>/dev/null | grep -q "ghcr.io"; then
        log_info "Already logged in to ghcr.io"
        return 0
    fi

    log_info "Logging in to GitHub Container Registry..."

    # 토큰 찾기 (우선순위 순서)
    local token=""

    # 1. 환경변수
    if [ -n "${GITHUB_TOKEN:-}" ]; then
        token="$GITHUB_TOKEN"
        log_info "Using GITHUB_TOKEN from environment"
    # 2. .env 파일
    elif [ -f "${PROJECT_ROOT}/.env" ] && grep -q "GITHUB_TOKEN" "${PROJECT_ROOT}/.env"; then
        # export GITHUB_TOKEN="..." 형식 지원
        token=$(grep "GITHUB_TOKEN" "${PROJECT_ROOT}/.env" | sed 's/^export //' | cut -d= -f2 | tr -d '"' | tr -d "'" | xargs)
        log_info "Using GITHUB_TOKEN from .env file"
    # 3. ~/.github-token 파일
    elif [ -f ~/.github-token ]; then
        token=$(cat ~/.github-token | xargs)
        log_info "Using token from ~/.github-token"
    else
        log_error "GitHub token not found!"
        log_error "Please set GITHUB_TOKEN environment variable or create ~/.github-token file"
        log_error ""
        log_error "Get token from: https://github.com/settings/tokens"
        log_error "Required permissions: write:packages, read:packages"
        return 1
    fi

    # GitHub username 추출
    local username="${REGISTRY#ghcr.io/}"
    username="${username%%/*}"

    log_info "Logging in as: $username"

    if echo "$token" | docker login ghcr.io -u "$username" --password-stdin 2>&1 | grep -q "Login Succeeded"; then
        log_success "Successfully logged in to ghcr.io"
        return 0
    else
        log_error "Failed to login to ghcr.io"
        log_error "Please check your token and permissions"
        return 1
    fi
}

##############################################################################
# VM 이미지 빌드
##############################################################################

build_vm_image() {
    local target="$1"
    local image_name="$2"
    local primary_image="${REGISTRY}/${image_name}:${IMAGE_TAG}"

    log_info "=========================================="
    log_info "Building: $image_name"
    log_info "Target: $target"
    log_info "Primary tag: $IMAGE_TAG"
    log_info "=========================================="

    cd "${PROJECT_ROOT}"

    # 빌드 (primary 태그)
    if ! docker build \
        -f ops/docker/op-stack-go/Dockerfile \
        --target "$target" \
        --platform linux/amd64 \
        --build-arg GIT_COMMIT="$GIT_COMMIT_FULL" \
        --build-arg GIT_DATE="$(git show -s --format='%ct')" \
        -t "$primary_image" \
        .; then
        log_error "Failed to build: $primary_image"
        return 1
    fi

    log_success "Built: $primary_image"

    # 추가 태그 생성
    local additional_tags=()

    # 1. 커밋 해시 태그 (primary가 아닌 경우)
    if [ "$IMAGE_TAG" != "$GIT_COMMIT_SHORT" ] && [ "$GIT_COMMIT_SHORT" != "unknown" ]; then
        additional_tags+=("${REGISTRY}/${image_name}:${GIT_COMMIT_SHORT}")
    fi

    # 2. latest 태그 (Git 태그 버전이거나 명시적으로 요청된 경우)
    if [ -n "$GIT_TAG" ] || [ "$IMAGE_TAG" = "latest" ]; then
        if [ "$IMAGE_TAG" != "latest" ]; then
            additional_tags+=("${REGISTRY}/${image_name}:latest")
        fi
    fi

    # 추가 태그 적용
    if [ ${#additional_tags[@]} -gt 0 ]; then
        for tag in "${additional_tags[@]}"; do
            docker tag "$primary_image" "$tag"
            log_info "  ✓ Tagged: $tag"
        done
    fi

    return 0
}

##############################################################################
# Kona Client 이미지 빌드 (Rust + RISC-V)
##############################################################################

build_kona_image() {
    local image_name="vm-kona-client"
    local primary_image="${REGISTRY}/${image_name}:${IMAGE_TAG}"

    log_info "=========================================="
    log_info "Building: $image_name (Rust + RISC-V)"
    log_info "Primary tag: $IMAGE_TAG"
    log_info "=========================================="

    # Check for kona repository
    local kona_dir="${KONA_DIR:-${PROJECT_ROOT}/../kona}"

    if [ ! -d "$kona_dir" ]; then
        log_warn "kona repository not found at: $kona_dir"
        log_info "Cloning kona repository..."

        local parent_dir="$(dirname ${PROJECT_ROOT})"
        cd "$parent_dir"

        if git clone --depth 1 https://github.com/op-rs/kona.git; then
            log_success "kona repository cloned successfully"
            kona_dir="${parent_dir}/kona"
        else
            log_error "Failed to clone kona repository"
            return 1
        fi
    fi

    cd "$kona_dir"

    # Create Dockerfile for kona-client
    local dockerfile_kona="${kona_dir}/Dockerfile.kona-export"
    cat > "$dockerfile_kona" <<'EOF'
# Multi-stage Docker build for kona-client (RISC-V only)
FROM ghcr.io/op-rs/kona/asterisc-builder:0.3.0 AS builder

WORKDIR /workspace

# Copy the entire kona repository
COPY . .

# Build RISC-V kona-client for prestate generation
# Note: This is the native target for kona-client (runs in asterisc VM)
RUN cargo build -Zbuild-std=core,alloc -p kona-client --bin kona-client --profile release-client-lto

# Export stage - copy binary to root for easy extraction
FROM scratch AS export
# RISC-V binary - this is what op-challenger uses
COPY --from=builder /workspace/target/riscv64imac-unknown-none-elf/release-client-lto/kona-client /kona-client
EOF

    log_success "Created: $dockerfile_kona"

    # Build using Docker
    log_info "Running Docker build for kona-client (may take 10-15 minutes)..."
    log_warn "⏳ Rust 컴파일 진행 중... (진행 상황이 보이지 않을 수 있습니다)"
    echo ""

    if docker build --platform linux/amd64 -f "$dockerfile_kona" --progress=plain -t "$primary_image" .; then
        log_success "Built: $primary_image"

        # 추가 태그 생성
        local additional_tags=()

        if [ "$IMAGE_TAG" != "$GIT_COMMIT_SHORT" ] && [ "$GIT_COMMIT_SHORT" != "unknown" ]; then
            additional_tags+=("${REGISTRY}/${image_name}:${GIT_COMMIT_SHORT}")
        fi

        if [ -n "$GIT_TAG" ] || [ "$IMAGE_TAG" = "latest" ]; then
            if [ "$IMAGE_TAG" != "latest" ]; then
                additional_tags+=("${REGISTRY}/${image_name}:latest")
            fi
        fi

        # 추가 태그 적용
        if [ ${#additional_tags[@]} -gt 0 ]; then
            for tag in "${additional_tags[@]}"; do
                docker tag "$primary_image" "$tag"
                log_info "  ✓ Tagged: $tag"
            done
        fi

        cd "${PROJECT_ROOT}"
        return 0
    else
        log_error "Failed to build: $primary_image"
        cd "${PROJECT_ROOT}"
        return 1
    fi
}

push_vm_image() {
    local image_name="$1"

    log_info "Pushing all tags for: $image_name"

    # 해당 이미지의 모든 태그 찾기 (<none> 제외)
    local all_tags=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep "^${REGISTRY}/${image_name}:" | grep -v ":<none>$")

    local push_failed=0

    for full_image in $all_tags; do
        log_info "  Pushing: $full_image"
        if docker push "$full_image"; then
            log_success "    ✓ Pushed: $full_image"
        else
            log_error "    ✗ Failed: $full_image"
            ((push_failed++))
        fi
    done

    if [ $push_failed -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

##############################################################################
# 메인 실행
##############################################################################

main() {
    log_info "VM 이미지 빌드 시작"
    log_info "Registry: $REGISTRY"
    log_info "Primary Tag: $IMAGE_TAG"
    log_info "Git Commit: $GIT_COMMIT_SHORT (full: $GIT_COMMIT_FULL)"
    if [ -n "$GIT_TAG" ]; then
        log_info "Git Tag: $GIT_TAG"
    fi
    log_info "Push: $PUSH_IMAGES"
    echo ""

    log_info "Tagging strategy:"
    log_info "  - Primary: $IMAGE_TAG"
    if [ "$IMAGE_TAG" != "$GIT_COMMIT_SHORT" ]; then
        log_info "  - Commit: $GIT_COMMIT_SHORT"
    fi
    if [ -n "$GIT_TAG" ] && [ "$IMAGE_TAG" != "latest" ]; then
        log_info "  - Latest: latest"
    fi
    echo ""

    # 빌드할 이미지 목록 (target:image_name 형식)
    local images=(
        "cannon-builder:vm-cannon"
        "asterisc-builder:vm-asterisc"
        "op-program-builder:vm-op-program"
        "op-challenger-target:op-challenger"
        "op-node-target:op-node"
        "op-batcher-target:op-batcher"
        "op-proposer-target:op-proposer"
    )

    local failed=0
    local built_images=()

    # 빌드 (Go 기반 이미지)
    for item in "${images[@]}"; do
        local target="${item%%:*}"
        local image_name="${item##*:}"

        if ! build_vm_image "$target" "$image_name"; then
            ((failed++))
        else
            built_images+=("$image_name")
        fi
        echo ""
    done

    # Kona client 빌드 (Rust 기반, 별도 처리)
    log_info "Building kona-client (Rust + RISC-V)..."
    if build_kona_image; then
        built_images+=("vm-kona-client")
    else
        log_warn "⚠️  kona-client build failed (GameType 3 지원 불가)"
        log_warn "    GameType 0, 1, 2는 정상 작동합니다"
        # Don't fail entire build for optional kona
    fi
    echo ""

    # 푸시 (요청된 경우)
    if [ "$PUSH_IMAGES" = true ]; then
        log_info "=========================================="
        log_info "Pushing images to registry..."
        log_info "=========================================="

        # 로그인 확인
        if ! login_to_ghcr; then
            log_error "Cannot push without login"
            return 1
        fi
        echo ""

        if [ ${#built_images[@]} -gt 0 ]; then
            # Push 단계
            for image_name in "${built_images[@]}"; do
                if ! push_vm_image "$image_name"; then
                    ((failed++))
                fi
            done
        else
            log_warn "No images were built successfully, skipping push"
        fi
    fi

    # 결과 요약
    echo ""
    log_info "=========================================="
    log_info "빌드 완료"
    log_info "=========================================="

    if [ $failed -eq 0 ]; then
        log_success "모든 이미지 빌드 성공!"

        if [ "$PUSH_IMAGES" = true ]; then
            log_success "모든 이미지 푸시 완료!"
            echo ""
            log_warn "⚠️  이미지가 Private으로 업로드되었습니다"
            log_warn "    팀원이 사용하려면 수동으로 Public 설정 필요:"
            log_warn "    https://github.com/${REGISTRY#ghcr.io/}?tab=packages"
            log_warn "    각 패키지 → Settings → Change visibility → Make public"
            echo ""
            log_info "사용 방법:"
            log_info "  ./op-challenger/scripts/pull-vm-images.sh --tag $IMAGE_TAG"
            log_info "  ./op-challenger/scripts/deploy-modular.sh --dg-type 3"
        else
            log_warn "이미지를 푸시하려면: $0 --push"
        fi
    else
        log_error "$failed 개 이미지 빌드/푸시 실패"
        return 1
    fi
}

main "$@"

