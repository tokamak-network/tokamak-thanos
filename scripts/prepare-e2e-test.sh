#!/bin/bash
# E2E 테스트 환경 자동 준비 스크립트
# 사용법: ./scripts/prepare-e2e-test.sh

set -e  # 에러 발생 시 중단

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 프로젝트 루트 확인
if [ ! -f "go.mod" ]; then
    log_error "프로젝트 루트에서 실행해주세요"
    exit 1
fi

PROJECT_ROOT=$(pwd)

log_info "=========================================="
log_info "E2E 테스트 환경 자동 준비 시작"
log_info "=========================================="
log_info "프로젝트 루트: $PROJECT_ROOT"
echo ""

# 1. Cannon 빌드
log_info "1/7: Cannon 빌드 중..."
cd "$PROJECT_ROOT/cannon"
if make cannon > /tmp/cannon-build.log 2>&1; then
    log_success "Cannon 빌드 완료"
    file bin/cannon | grep -q "Mach-O" && log_success "  ✓ macOS 바이너리 확인됨"
else
    log_error "Cannon 빌드 실패. 로그: /tmp/cannon-build.log"
    exit 1
fi
echo ""

# 2. op-program-host 빌드
log_info "2/7: op-program-host 빌드 중..."
cd "$PROJECT_ROOT/op-program"
if make op-program-host > /tmp/op-program-host-build.log 2>&1; then
    log_success "op-program-host 빌드 완료"
else
    log_error "op-program-host 빌드 실패. 로그: /tmp/op-program-host-build.log"
    exit 1
fi
echo ""

# 3. mips64 ELF 빌드
log_info "3/7: mips64 ELF 빌드 중..."
if make op-program-client-mips64 > /tmp/mips64-build.log 2>&1; then
    log_success "mips64 ELF 빌드 완료"
    ls -lh bin/op-program-client64.elf | awk '{print "  ✓ 파일 크기: " $5}'
else
    log_error "mips64 빌드 실패. 로그: /tmp/mips64-build.log"
    exit 1
fi
echo ""

# 4. mt64 Prestate 생성
log_info "4/7: mt64 Prestate 생성 중..."
cd "$PROJECT_ROOT"
if ./cannon/bin/cannon load-elf \
    --type multithreaded64-4 \
    --path op-program/bin/op-program-client64.elf \
    --out op-program/bin/prestate-mt64.bin.gz \
    --meta op-program/bin/meta-mt64.json > /tmp/prestate-load.log 2>&1; then
    log_success "mt64 Prestate 로드 완료"
    ls -lh op-program/bin/prestate-mt64.bin.gz | awk '{print "  ✓ 파일 크기: " $5}'
else
    log_error "Prestate 로드 실패. 로그: /tmp/prestate-load.log"
    exit 1
fi
echo ""

# 5. Prestate Proof 생성
log_info "5/7: Prestate Proof 생성 중..."
if ./cannon/bin/cannon run \
    --proof-at '=0' \
    --stop-at '=1' \
    --input op-program/bin/prestate-mt64.bin.gz \
    --meta op-program/bin/meta-mt64.json \
    --proof-fmt 'op-program/bin/%d.json' \
    --output "" > /tmp/prestate-proof.log 2>&1; then
    mv op-program/bin/0.json op-program/bin/prestate-proof-mt64.json
    log_success "Prestate Proof 생성 완료"
else
    log_error "Prestate Proof 생성 실패. 로그: /tmp/prestate-proof.log"
    exit 1
fi
echo ""

# 6. Prestate Hash 추출 및 템플릿 임시 업데이트
log_info "6/7: Prestate Hash 추출 및 템플릿 임시 업데이트 중..."
PRESTATE_HASH=$(cat op-program/bin/prestate-proof-mt64.json | jq -r '.pre')
log_success "Prestate Hash: $PRESTATE_HASH"

DEPLOY_CONFIG_TEMPLATE="$PROJECT_ROOT/packages/tokamak/contracts-bedrock/deploy-config/devnetL1-template.json"
if [ ! -f "$DEPLOY_CONFIG_TEMPLATE" ]; then
    log_error "템플릿 파일을 찾을 수 없습니다: $DEPLOY_CONFIG_TEMPLATE"
    exit 1
fi

# 템플릿 백업
cp "$DEPLOY_CONFIG_TEMPLATE" "${DEPLOY_CONFIG_TEMPLATE}.backup-script"
log_info "  템플릿 백업 완료: ${DEPLOY_CONFIG_TEMPLATE}.backup-script"

# 템플릿에 새 prestate hash 임시 업데이트
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/\"faultGameAbsolutePrestate\": \"0x[^\"]*\"/\"faultGameAbsolutePrestate\": \"$PRESTATE_HASH\"/" "$DEPLOY_CONFIG_TEMPLATE"
else
    sed -i "s/\"faultGameAbsolutePrestate\": \"0x[^\"]*\"/\"faultGameAbsolutePrestate\": \"$PRESTATE_HASH\"/" "$DEPLOY_CONFIG_TEMPLATE"
fi
log_success "  템플릿에 새 prestate hash 임시 적용"
echo ""

# 7. devnet-allocs 재생성 (새 prestate hash가 반영된 템플릿 사용)
log_info "7/7: .devnet 파일 재생성 중..."
log_warning "  이 단계는 약 5분 소요됩니다..."
log_info "  make devnet-allocs가 템플릿에서 devnetL1.json을 생성합니다"
cd "$PROJECT_ROOT"

if make devnet-allocs > /tmp/devnet-allocs.log 2>&1; then
    log_success ".devnet 파일 생성 완료"

    # 생성된 devnetL1.json 검증
    DEPLOY_CONFIG="$PROJECT_ROOT/packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json"
    if [ -f "$DEPLOY_CONFIG" ]; then
        DEPLOYED_PRESTATE=$(cat "$DEPLOY_CONFIG" | jq -r '.faultGameAbsolutePrestate')
        log_info "  devnetL1.json의 Prestate Hash: $DEPLOYED_PRESTATE"

        if [ "$DEPLOYED_PRESTATE" == "$PRESTATE_HASH" ]; then
            log_success "  ✓ Prestate Hash 일치 확인!"
        else
            log_error "  ✗ Prestate Hash 불일치"
            log_error "    예상: $PRESTATE_HASH"
            log_error "    실제: $DEPLOYED_PRESTATE"
            # 템플릿 복원 후 종료
            mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
            exit 1
        fi
    else
        log_error "devnetL1.json이 생성되지 않았습니다"
        # 템플릿 복원 후 종료
        mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
        exit 1
    fi
else
    log_error "devnet-allocs 실패. 로그: /tmp/devnet-allocs.log"
    # 템플릿 복원 후 종료
    mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
    exit 1
fi

# 템플릿 복원
mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
log_info "  템플릿 원래 상태로 복원 완료"
echo ""

# 완료 메시지
log_success "=========================================="
log_success "E2E 테스트 환경 준비 완료!"
log_success "=========================================="
echo ""
log_info "다음 명령으로 E2E 테스트를 실행하세요:"
echo ""
echo "  cd op-e2e"
echo "  go test -v ./faultproofs -run TestOutputCannonGame -timeout 30m"
echo ""
log_info "생성된 파일:"
echo "  - cannon/bin/cannon (macOS)"
echo "  - op-program/bin/op-program (macOS)"
echo "  - op-program/bin/op-program-client64.elf (mips64)"
echo "  - op-program/bin/prestate-mt64.bin.gz"
echo "  - op-program/bin/prestate-proof-mt64.json"
echo "  - .devnet/allocs-l1.json"
echo "  - .devnet/allocs-l2-*.json"
echo "  - .devnet/addresses.json (prestate hash: $PRESTATE_HASH)"
echo ""

