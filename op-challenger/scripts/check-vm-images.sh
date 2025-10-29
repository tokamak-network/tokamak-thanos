#!/usr/bin/env bash

##############################################################################
# VM 이미지 레지스트리 업로드 확인 스크립트
#
# 사용법:
#   ./check-vm-images.sh                    # 현재 커밋 확인
#   ./check-vm-images.sh --tag 47efca11c    # 특정 태그 확인
##############################################################################

set -euo pipefail

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

##############################################################################
# ⚙️  설정 (여기를 수정하세요)
##############################################################################

# 🔧 레지스트리 설정 (필요시 수정)
# - 개인 계정: ghcr.io/zena-park
# - Organization: ghcr.io/tokamak-network
DEFAULT_REGISTRY="ghcr.io/zena-park"
REGISTRY="${VM_REGISTRY:-$DEFAULT_REGISTRY}"

# 이미지 태그 (기본값: 현재 Git 커밋 해시)
IMAGE_TAG="${VM_IMAGE_TAG:-$(git rev-parse --short HEAD 2>/dev/null || echo "latest")}"

##############################################################################

# 인자 파싱
while [[ $# -gt 0 ]]; do
    case $1 in
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
            echo "  --tag TAG           Image tag to check (default: current git commit)"
            echo "  --registry URL      Registry URL (default: ghcr.io/zena-park)"
            echo "  -h, --help          Show this help"
            echo ""
            echo "Examples:"
            echo "  $0                        # Check current commit"
            echo "  $0 --tag 47efca11c        # Check specific tag"
            echo "  $0 --tag latest           # Check latest tag"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}VM 이미지 레지스트리 확인${NC}"
echo -e "${BLUE}========================================${NC}"
echo "Registry: $REGISTRY"
echo "Tag: $IMAGE_TAG"
echo ""

images=(
    "vm-cannon"
    "vm-asterisc"
    "vm-op-program"
    "vm-kona-client"
    "op-challenger"
    "op-node"
    "op-batcher"
    "op-proposer"
)

success=0
failed=0

for img in "${images[@]}"; do
    full_image="$REGISTRY/$img:$IMAGE_TAG"
    echo -n "$img:$IMAGE_TAG ... "

    if docker manifest inspect "$full_image" &>/dev/null; then
        echo -e "${GREEN}✓ 존재${NC}"
        ((success++))
    else
        echo -e "${RED}✗ 없음${NC}"
        ((failed++))
    fi
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "성공: ${GREEN}$success${NC} / 실패: ${RED}$failed${NC}"
echo -e "${BLUE}========================================${NC}"

if [ $failed -eq 0 ]; then
    echo -e "${GREEN}✅ 모든 이미지가 레지스트리에 존재합니다! (총 8개)${NC}"
    echo ""
    echo -e "${BLUE}배포 시 사용 방법:${NC}"
    echo "  # VM 바이너리 다운로드"
    echo "  ./op-challenger/scripts/pull-vm-images.sh --tag $IMAGE_TAG"
    echo ""
    echo "  # 배포"
    echo "  ./op-challenger/scripts/deploy-modular.sh --dg-type 3"
    exit 0
else
    echo -e "${RED}❌ 일부 이미지가 레지스트리에 없습니다.${NC}"
    echo ""
    echo -e "${YELLOW}해결 방법:${NC}"
    echo "  1. 이미지 빌드 및 업로드:"
    echo "     ./op-challenger/scripts/build-vm-images.sh --push"
    echo ""
    echo "  2. 다른 태그 확인:"
    echo "     $0 --tag <다른_태그>"
    exit 1
fi

