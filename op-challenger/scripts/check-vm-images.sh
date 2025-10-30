#!/usr/bin/env bash

##############################################################################
# VM 이미지/바이너리 확인 스크립트
#
# 사용법:
#   ./check-vm-images.sh                    # 레지스트리 확인 (현재 커밋)
#   ./check-vm-images.sh --tag 47efca11c    # 레지스트리 확인 (특정 태그)
#   ./check-vm-images.sh --local            # 로컬 바이너리 검증
##############################################################################

set -euo pipefail

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

##############################################################################
# ⚙️  설정 (여기를 수정하세요)
##############################################################################

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# 🔧 레지스트리 설정 (필요시 수정)
# - 개인 계정: ghcr.io/zena-park
# - Organization: ghcr.io/tokamak-network
DEFAULT_REGISTRY="ghcr.io/zena-park"
REGISTRY="${VM_REGISTRY:-$DEFAULT_REGISTRY}"

# 이미지 태그 (기본값: 현재 Git 커밋 해시)
IMAGE_TAG="${VM_IMAGE_TAG:-$(git rev-parse --short HEAD 2>/dev/null || echo "latest")}"

# 모드 선택
CHECK_LOCAL=false

##############################################################################

# 인자 파싱
while [[ $# -gt 0 ]]; do
    case $1 in
        --local)
            CHECK_LOCAL=true
            shift
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
            echo "  --local             Check local binaries and prestate files"
            echo "  --tag TAG           Image tag to check (default: current git commit)"
            echo "  --registry URL      Registry URL (default: ghcr.io/zena-park)"
            echo "  -h, --help          Show this help"
            echo ""
            echo "Examples:"
            echo "  $0                        # Check registry (current commit)"
            echo "  $0 --tag 47efca11c        # Check registry (specific tag)"
            echo "  $0 --tag latest           # Check registry (latest tag)"
            echo "  $0 --local                # Verify local binaries"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

##############################################################################
# 로컬 바이너리 검증
##############################################################################

check_local_binaries() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}로컬 VM 바이너리 및 Prestate 검증${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""

    local success=0
    local failed=0

    # 1. Cannon 바이너리
    echo -e "${MAGENTA}1️⃣  Cannon (GameType 0, 1)${NC}"
    local cannon_bin="${PROJECT_ROOT}/cannon/bin/cannon"
    if [ -f "$cannon_bin" ]; then
        local file_type=$(file "$cannon_bin" 2>/dev/null || echo "")
        if echo "$file_type" | grep -q "ELF 64-bit.*x86-64"; then
            echo -e "  ${GREEN}✓${NC} Binary: $cannon_bin (Linux x86-64 ELF)"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} Binary: $cannon_bin (NOT Linux x86-64!)"
            echo "    Actual: $file_type"
            ((failed++))
        fi
    else
        echo -e "  ${RED}✗${NC} Binary not found: $cannon_bin"
        ((failed++))
    fi
    echo ""

    # 2. Asterisc 바이너리
    echo -e "${MAGENTA}2️⃣  Asterisc (GameType 2)${NC}"
    local asterisc_bin="${PROJECT_ROOT}/asterisc/bin/asterisc"
    if [ -f "$asterisc_bin" ]; then
        local file_type=$(file "$asterisc_bin" 2>/dev/null || echo "")
        if echo "$file_type" | grep -q "ELF 64-bit.*x86-64"; then
            echo -e "  ${GREEN}✓${NC} Binary: $asterisc_bin (Linux x86-64 ELF)"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} Binary: $asterisc_bin (NOT Linux x86-64!)"
            echo "    Actual: $file_type"
            ((failed++))
        fi
    else
        echo -e "  ${RED}✗${NC} Binary not found: $asterisc_bin"
        ((failed++))
    fi

    # Asterisc prestate files
    local asterisc_prestate_proof="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
    local asterisc_prestate_runtime="${PROJECT_ROOT}/asterisc/bin/prestate.json"

    if [ -f "$asterisc_prestate_proof" ]; then
        if jq -e '.pre' "$asterisc_prestate_proof" >/dev/null 2>&1; then
            local hash=$(jq -r '.pre' "$asterisc_prestate_proof")
            echo -e "  ${GREEN}✓${NC} Prestate (deployment): prestate-proof.json"
            echo "    Hash: ${hash:0:10}...${hash: -8}"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} prestate-proof.json: Invalid format (missing .pre field)"
            ((failed++))
        fi
    else
        echo -e "  ${YELLOW}⚠${NC}  prestate-proof.json not found"
    fi

    if [ -f "$asterisc_prestate_runtime" ]; then
        if jq -e '.stateHash' "$asterisc_prestate_runtime" >/dev/null 2>&1; then
            echo -e "  ${GREEN}✓${NC} Prestate (runtime): prestate.json (full VMState)"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} prestate.json: Invalid format (missing .stateHash field)"
            ((failed++))
        fi
    else
        echo -e "  ${RED}✗${NC} prestate.json not found (REQUIRED for runtime)"
        ((failed++))
    fi
    echo ""

    # 3. Kona 바이너리
    echo -e "${MAGENTA}3️⃣  Kona (GameType 3)${NC}"
    local kona_bin="${PROJECT_ROOT}/bin/kona-client"
    if [ -f "$kona_bin" ]; then
        local file_type=$(file "$kona_bin" 2>/dev/null || echo "")
        if echo "$file_type" | grep -q "ELF 64-bit.*RISC-V"; then
            echo -e "  ${GREEN}✓${NC} Binary: $kona_bin (RISC-V 64-bit ELF)"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} Binary: $kona_bin (NOT RISC-V!)"
            echo "    Actual: $file_type"
            ((failed++))
        fi
    else
        echo -e "  ${RED}✗${NC} Binary not found: $kona_bin"
        echo -e "      ${YELLOW}→${NC} Run: ./build-vm-images.sh --vm kona"
        ((failed++))
    fi

    # Kona prestate files
    local kona_prestate_proof="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
    local kona_prestate_runtime="${PROJECT_ROOT}/bin/prestate.bin.gz"

    if [ -f "$kona_prestate_proof" ]; then
        if jq -e '.pre' "$kona_prestate_proof" >/dev/null 2>&1; then
            local hash=$(jq -r '.pre' "$kona_prestate_proof")
            echo -e "  ${GREEN}✓${NC} Prestate (deployment): prestate-kona.json"
            echo "    Hash: ${hash:0:10}...${hash: -8}"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} prestate-kona.json: Invalid format (missing .pre field)"
            ((failed++))
        fi
    else
        echo -e "  ${YELLOW}⚠${NC}  prestate-kona.json not found"
    fi

    if [ -f "$kona_prestate_runtime" ]; then
        if file "$kona_prestate_runtime" | grep -q "gzip compressed"; then
            echo -e "  ${GREEN}✓${NC} Prestate (runtime): prestate.bin.gz (gzipped VMState)"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} prestate.bin.gz: Not a gzip file"
            ((failed++))
        fi
    else
        echo -e "  ${RED}✗${NC} prestate.bin.gz not found (REQUIRED for runtime)"
        echo -e "      ${YELLOW}→${NC} Run: ./build-vm-images.sh --vm kona"
        echo -e "      ${YELLOW}→${NC} Or: ./pull-vm-images.sh --vm kona"
        ((failed++))
    fi
    echo ""

    # 4. op-program 바이너리
    echo -e "${MAGENTA}4️⃣  op-program (Server)${NC}"
    local op_program_bin="${PROJECT_ROOT}/op-program/bin/op-program"
    if [ -f "$op_program_bin" ]; then
        local file_type=$(file "$op_program_bin" 2>/dev/null || echo "")
        if echo "$file_type" | grep -q "ELF 64-bit.*x86-64"; then
            echo -e "  ${GREEN}✓${NC} Binary: $op_program_bin (Linux x86-64 ELF)"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} Binary: $op_program_bin (NOT Linux x86-64!)"
            echo "    Actual: $file_type"
            ((failed++))
        fi
    else
        echo -e "  ${RED}✗${NC} Binary not found: $op_program_bin"
        ((failed++))
    fi

    local op_program_prestate="${PROJECT_ROOT}/op-program/bin/prestate.json"
    local op_program_prestate_proof="${PROJECT_ROOT}/op-program/bin/prestate-proof.json"

    if [ -f "$op_program_prestate" ]; then
        # Cannon prestate uses different format: pc, step, memory, etc.
        if jq -e '.pc' "$op_program_prestate" >/dev/null 2>&1; then
            echo -e "  ${GREEN}✓${NC} Prestate (runtime): prestate.json (Cannon MIPS State)"
            ((success++))
        else
            echo -e "  ${RED}✗${NC} prestate.json: Invalid format (missing .pc field)"
            ((failed++))
        fi
    else
        echo -e "  ${YELLOW}⚠${NC}  prestate.json not found"
    fi

    if [ -f "$op_program_prestate_proof" ]; then
        if jq -e '.pre' "$op_program_prestate_proof" >/dev/null 2>&1; then
            local hash=$(jq -r '.pre' "$op_program_prestate_proof")
            echo -e "  ${GREEN}✓${NC} Prestate (deployment): prestate-proof.json"
            echo "    Hash: ${hash:0:10}...${hash: -8}"
            ((success++))
        else
            echo -e "  ${YELLOW}⚠${NC}  prestate-proof.json: Invalid format"
        fi
    else
        echo -e "  ${YELLOW}⚠${NC}  prestate-proof.json not found"
    fi
    echo ""

    # 결과 요약
    echo -e "${BLUE}========================================${NC}"
    echo -e "검증 성공: ${GREEN}$success${NC} / 실패: ${RED}$failed${NC}"
    echo -e "${BLUE}========================================${NC}"

    if [ $failed -eq 0 ]; then
        echo -e "${GREEN}✅ 모든 바이너리와 prestate 파일이 올바릅니다!${NC}"
        echo ""
        echo -e "${BLUE}다음 단계:${NC}"
        echo "  ./op-challenger/scripts/deploy-modular.sh --dg-type 3"
        return 0
    else
        echo -e "${RED}❌ 일부 파일에 문제가 있습니다.${NC}"
        echo ""
        echo -e "${YELLOW}해결 방법:${NC}"
        echo "  1. VM 이미지 다시 받기:"
        echo "     ./op-challenger/scripts/pull-vm-images.sh --tag latest"
        echo ""
        echo "  2. 또는 로컬 빌드:"
        echo "     ./op-challenger/scripts/build-vm-images.sh --no-cache"
        return 1
    fi
}

##############################################################################
# 레지스트리 이미지 확인
##############################################################################

check_registry_images() {
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

    local success=0
    local failed=0

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
        return 0
    else
        echo -e "${RED}❌ 일부 이미지가 레지스트리에 없습니다.${NC}"
        echo ""
        echo -e "${YELLOW}해결 방법:${NC}"
        echo "  1. 이미지 빌드 및 업로드:"
        echo "     ./op-challenger/scripts/build-vm-images.sh --push"
        echo ""
        echo "  2. 다른 태그 확인:"
        echo "     $0 --tag <다른_태그>"
        return 1
    fi
}

##############################################################################
# 메인 실행
##############################################################################

if [ "$CHECK_LOCAL" = true ]; then
    check_local_binaries
else
    check_registry_images
fi

