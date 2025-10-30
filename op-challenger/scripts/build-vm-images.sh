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
NO_CACHE=""
BUILD_ASTERISC_ONLY=false
BUILD_KONA_ONLY=false

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
        --no-cache)
            NO_CACHE="--no-cache"
            shift
            ;;
        --asterisc-only)
            BUILD_ASTERISC_ONLY=true
            shift
            ;;
        --kona-only)
            BUILD_KONA_ONLY=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --push              Push images to registry after build"
            echo "  --registry REGISTRY Registry URL (default: ghcr.io/zena-park)"
            echo "  --tag TAG           Override image tag (default: auto-detect)"
            echo "  --no-cache          Force rebuild without using cache"
            echo "  --asterisc-only     Build only vm-asterisc image"
            echo "  --kona-only         Build only vm-kona-client image"
            echo "  -h, --help          Show this help"
            echo ""
            echo "Tagging Strategy:"
            echo "  - Auto-detects Git tag or commit hash"
            echo "  - If Git tag exists: uses tag (e.g., v1.0.0)"
            echo "  - Otherwise: uses commit hash (e.g., 47efca11)"
            echo "  - Also tags as 'latest' when appropriate"
            echo ""
            echo "Examples:"
            echo "  $0                           # Build all images"
            echo "  $0 --push                   # Build and push to registry"
            echo "  $0 --no-cache --push        # Rebuild without cache and push"
            echo "  $0 --kona-only --push       # Build only kona and push"
            echo "  $0 --asterisc-only --push   # Build only asterisc and push"
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
        $NO_CACHE \
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

    # 2. latest 태그 (push 시 항상 업데이트)
    if [ "$PUSH_IMAGES" = true ] && [ "$IMAGE_TAG" != "latest" ]; then
        additional_tags+=("${REGISTRY}/${image_name}:latest")
    elif [ -n "$GIT_TAG" ] || [ "$IMAGE_TAG" = "latest" ]; then
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
# Asterisc 이미지 빌드 (외부 repo 사용)
##############################################################################

build_asterisc_image() {
    local image_name="vm-asterisc"
    local primary_image="${REGISTRY}/${image_name}:${IMAGE_TAG}"

    log_info "=========================================="
    log_info "Building: $image_name (RISC-V + Prestate)"
    log_info "Primary tag: $IMAGE_TAG"
    log_info "=========================================="

    # Check for asterisc repository
    local asterisc_dir="${ASTERISC_DIR:-${PROJECT_ROOT}/../asterisc}"

    if [ ! -d "$asterisc_dir" ]; then
        log_warn "asterisc repository not found at: $asterisc_dir"
        log_info "Cloning asterisc repository..."

        local parent_dir="$(dirname ${PROJECT_ROOT})"
        cd "$parent_dir"

        if git clone https://github.com/ethereum-optimism/asterisc.git; then
            log_success "asterisc repository cloned successfully"
            asterisc_dir="${parent_dir}/asterisc"
        else
            log_error "Failed to clone asterisc repository"
            return 1
        fi
    fi

    cd "$asterisc_dir"

    # Checkout to master branch
    local ASTERISC_TAG="${ASTERISC_TAG:-master}"
    log_info "Using asterisc tag: $ASTERISC_TAG"
    git checkout "$ASTERISC_TAG" 2>/dev/null || log_warn "Could not checkout $ASTERISC_TAG, using current branch"
    git submodule update --init --recursive

    # Build using asterisc's reproducible-prestate target
    log_info "Building asterisc with reproducible prestate..."
    log_warn "⏳ This may take 5-10 minutes..."
    echo ""

    # Use make reproducible-prestate to build in Docker (Linux x86-64)
    if make reproducible-prestate; then
        log_success "Built asterisc binary and prestate ⭐"

        # Files are now in bin/
        if [ ! -f "bin/asterisc" ]; then
            log_error "asterisc binary not found after build"
            cd "${PROJECT_ROOT}"
            return 1
        fi

        if [ ! -f "bin/prestate.json" ]; then
            log_error "prestate.json not found after build"
            cd "${PROJECT_ROOT}"
            return 1
        fi

        # Create a minimal Docker image with the binaries
        local asterisc_output_dir="${PROJECT_ROOT}/asterisc/bin"
        mkdir -p "$asterisc_output_dir"

        # Copy binaries
        cp "bin/asterisc" "$asterisc_output_dir/"
        cp "bin/prestate.json" "$asterisc_output_dir/"

        # Convert prestate.json (.stateHash) to prestate-proof.json (.pre) format
        if [ -f "bin/prestate.json" ]; then
            local state_hash=$(cat "bin/prestate.json" | jq -r '.stateHash' 2>/dev/null || echo "")
            if [ -n "$state_hash" ] && [ "$state_hash" != "null" ]; then
                echo "{\"pre\": \"$state_hash\"}" | jq . > "$asterisc_output_dir/prestate-proof.json"
                log_success "  ✓ prestate-proof.json created (deployment format)"
            fi
        fi

        # Create a scratch image
        cat > "${PROJECT_ROOT}/.Dockerfile.asterisc-export" << 'EOF'
FROM scratch
COPY asterisc/bin/asterisc /app/asterisc/bin/asterisc
COPY asterisc/bin/prestate.json /app/asterisc/bin/prestate.json
COPY asterisc/bin/prestate-proof.json /app/asterisc/bin/prestate-proof.json
EOF

        cd "${PROJECT_ROOT}"

        if docker build --platform linux/amd64 $NO_CACHE -f "${PROJECT_ROOT}/.Dockerfile.asterisc-export" -t "$primary_image" .; then
            log_success "Built: $primary_image ⭐"
            rm -f "${PROJECT_ROOT}/.Dockerfile.asterisc-export"

            # 추가 태그 생성
            local additional_tags=()

            if [ "$IMAGE_TAG" != "$GIT_COMMIT_SHORT" ] && [ "$GIT_COMMIT_SHORT" != "unknown" ]; then
                additional_tags+=("${REGISTRY}/${image_name}:${GIT_COMMIT_SHORT}")
            fi

            # latest 태그 (push 시 항상 업데이트)
            if [ "$PUSH_IMAGES" = true ] && [ "$IMAGE_TAG" != "latest" ]; then
                additional_tags+=("${REGISTRY}/${image_name}:latest")
            elif [ -n "$GIT_TAG" ] || [ "$IMAGE_TAG" = "latest" ]; then
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
        else
            log_error "Failed to create Docker image"
            rm -f "${PROJECT_ROOT}/.Dockerfile.asterisc-export"
            return 1
        fi
    else
        log_error "Failed to build asterisc"
        cd "${PROJECT_ROOT}"
        return 1
    fi
}

##############################################################################
# Kona Client 이미지 빌드 (Rust + RISC-V)
##############################################################################

build_kona_image() {
    local image_name="vm-kona-client"
    local primary_image="${REGISTRY}/${image_name}:${IMAGE_TAG}"

    log_info "=========================================="
    log_info "Building: $image_name (Rust + RISC-V + Prestate)"
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

    # Use kona's official asterisc-repro.dockerfile (includes prestate generation!)
    log_info "Using kona's official prestate generation system..."
    log_info "  Dockerfile: kona/docker/fpvm-prestates/asterisc-repro.dockerfile"
    log_info "  This builds: asterisc + kona-client + prestate files"
    echo ""

    local dockerfile_path="docker/fpvm-prestates/asterisc-repro.dockerfile"

    if [ ! -f "$dockerfile_path" ]; then
        log_error "Kona's official dockerfile not found: $dockerfile_path"
        return 1
    fi

    # Build using kona's official Dockerfile (includes asterisc + kona-client + prestate!)
    log_info "Running Docker build for kona with prestate (may take 15-20 minutes)..."
    log_warn "⏳ This builds: asterisc + kona-client + prestate-proof.json + prestate.bin.gz"
    echo ""

    # Build arguments for kona's Dockerfile
    local ASTERISC_TAG="${ASTERISC_TAG:-master}"
    local CLIENT_BIN="${CLIENT_BIN:-kona-client}"
    local CLIENT_TAG="${CLIENT_TAG:-main}"

    log_info "Build arguments:"
    log_info "  ASTERISC_TAG=$ASTERISC_TAG"
    log_info "  CLIENT_BIN=$CLIENT_BIN"
    log_info "  CLIENT_TAG=$CLIENT_TAG"
    echo ""

    if docker build --platform linux/amd64 $NO_CACHE \
        --build-arg ASTERISC_TAG="$ASTERISC_TAG" \
        --build-arg CLIENT_BIN="$CLIENT_BIN" \
        --build-arg CLIENT_TAG="$CLIENT_TAG" \
        -f "$dockerfile_path" --progress=plain -t "$primary_image" .; then
        log_success "Built: $primary_image (with prestate files!) ⭐"

        # 이미지에서 prestate 파일들을 로컬로 추출
        log_info "Extracting prestate files from image..."
        local container_id=$(docker create "$primary_image" true 2>/dev/null)

        if [ -n "$container_id" ]; then
            # Extract prestate-proof.json (배포용 - .pre 포맷)
            mkdir -p "${PROJECT_ROOT}/op-program/bin"
            if docker cp "${container_id}:/prestate-proof.json" "${PROJECT_ROOT}/op-program/bin/prestate-kona.json" 2>/dev/null; then
                log_success "  ✓ prestate-kona.json extracted (deployment format)"

                # Extract and display hash
                local kona_hash=$(cat "${PROJECT_ROOT}/op-program/bin/prestate-kona.json" | jq -r '.pre' 2>/dev/null || echo "")
                if [ -n "$kona_hash" ] && [ "$kona_hash" != "null" ]; then
                    log_success "  ✓ Kona deployment prestate hash: ${kona_hash:0:10}...${kona_hash: -8}"
                fi
            else
                log_warn "  ⚠️  prestate-proof.json not found in image"
            fi

            # Extract prestate.bin.gz (런타임용 - full VMState, 압축)
            mkdir -p "${PROJECT_ROOT}/bin"
            if docker cp "${container_id}:/prestate.bin.gz" "${PROJECT_ROOT}/bin/prestate.bin.gz" 2>/dev/null; then
                log_success "  ✓ prestate.bin.gz extracted (runtime full VMState) ⭐"
            else
                log_warn "  ⚠️  prestate.bin.gz not found in image"
                log_warn "      GameType 3 may not work properly without this file"
            fi

            # Extract kona-client binary as well
            if docker cp "${container_id}:/kona-client-elf" "${PROJECT_ROOT}/bin/kona-client" 2>/dev/null; then
                log_success "  ✓ kona-client binary extracted"
            else
                log_warn "  ⚠️  kona-client binary not found in image"
            fi

            docker rm "$container_id" >/dev/null 2>&1
        else
            log_error "Failed to create temporary container from image"
        fi

        # 추가 태그 생성
        local additional_tags=()

        if [ "$IMAGE_TAG" != "$GIT_COMMIT_SHORT" ] && [ "$GIT_COMMIT_SHORT" != "unknown" ]; then
            additional_tags+=("${REGISTRY}/${image_name}:${GIT_COMMIT_SHORT}")
        fi

        # latest 태그 (push 시 항상 업데이트)
        if [ "$PUSH_IMAGES" = true ] && [ "$IMAGE_TAG" != "latest" ]; then
            additional_tags+=("${REGISTRY}/${image_name}:latest")
        elif [ -n "$GIT_TAG" ] || [ "$IMAGE_TAG" = "latest" ]; then
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

    # 빌드할 이미지 목록 (target:image_name 형식) - asterisc 제외!
    local images=(
        "cannon-builder:vm-cannon"
        "op-program-builder:vm-op-program"
        "op-challenger-target:op-challenger"
        "op-node-target:op-node"
        "op-batcher-target:op-batcher"
        "op-proposer-target:op-proposer"
    )

    local failed=0
    local built_images=()

    # Selective build mode
    if [ "$BUILD_KONA_ONLY" = true ]; then
        log_info "Building kona-client only (--kona-only mode)..."
        if build_kona_image; then
            built_images+=("vm-kona-client")
        else
            log_error "Kona build failed"
            ((failed++))
        fi
    elif [ "$BUILD_ASTERISC_ONLY" = true ]; then
        log_info "Building asterisc only (--asterisc-only mode)..."
        if build_asterisc_image; then
            built_images+=("vm-asterisc")
        else
            log_error "Asterisc build failed"
            ((failed++))
        fi
    else
        # Full build mode (Go 기반 이미지)
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

        # Asterisc 빌드 (외부 repo, 별도 처리)
        log_info "Building asterisc (External RISC-V repo)..."
        if build_asterisc_image; then
            built_images+=("vm-asterisc")
        else
            log_warn "⚠️  asterisc build failed (GameType 2 지원 불가)"
            log_warn "    GameType 0, 1, 3은 정상 작동합니다"
        fi
        echo ""

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
    fi

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

