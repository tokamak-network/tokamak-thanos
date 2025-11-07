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

# 5. Prestate Proof 생성 (mt64)
log_info "5/7: Prestate Proof 생성 중..."
if ./cannon/bin/cannon run \
    --proof-at '=0' \
    --stop-at '=1' \
    --input op-program/bin/prestate-mt64.bin.gz \
    --meta op-program/bin/meta-mt64.json \
    --proof-fmt 'op-program/bin/%d.json' \
    --output "" > /tmp/prestate-proof.log 2>&1; then
    mv op-program/bin/0.json op-program/bin/prestate-proof-mt64.json
    # Deploy.s.sol이 prestate-proof.json을 읽으므로 복사
    cp op-program/bin/prestate-proof-mt64.json op-program/bin/prestate-proof.json
    log_success "Prestate Proof (mt64) 생성 완료"
    log_info "  ✓ prestate-proof-mt64.json 생성"
    log_info "  ✓ prestate-proof.json 복사 (Deploy.s.sol용)"
else
    log_error "Prestate Proof 생성 실패. 로그: /tmp/prestate-proof.log"
    exit 1
fi

# mt64Next도 생성 (E2E 테스트가 사용)
if [ -f op-program/bin/prestate-mt64Next.bin.gz ]; then
    log_info "  prestate-proof-mt64Next.json 생성 중..."
    if ./cannon/bin/cannon run \
        --proof-at '=0' \
        --stop-at '=1' \
        --input op-program/bin/prestate-mt64Next.bin.gz \
        --meta op-program/bin/meta-mt64Next.json \
        --proof-fmt 'op-program/bin/%d.json' \
        --output "" > /tmp/prestate-proof-next.log 2>&1; then
        mv op-program/bin/0.json op-program/bin/prestate-proof-mt64Next.json
        log_info "  ✓ prestate-proof-mt64Next.json 생성"
    else
        log_warning "  prestate-proof-mt64Next.json 생성 실패 (선택사항)"
    fi
fi
echo ""

# 6. Prestate Hash 추출
log_info "6/9: Prestate Hash 추출..."
PRESTATE_HASH=$(cat op-program/bin/prestate-proof-mt64.json | jq -r '.pre')
log_success "Prestate Hash: $PRESTATE_HASH"
echo ""

# 7. 1차 devnet-allocs 실행 (더미 genesis output root로)
log_info "7/9: 1차 .devnet 파일 생성 중 (genesis output root 계산용)..."
log_warning "  이 단계는 약 5분 소요됩니다..."

DEPLOY_CONFIG_TEMPLATE="$PROJECT_ROOT/packages/tokamak/contracts-bedrock/deploy-config/devnetL1-template.json"
if [ ! -f "$DEPLOY_CONFIG_TEMPLATE" ]; then
    log_error "템플릿 파일을 찾을 수 없습니다: $DEPLOY_CONFIG_TEMPLATE"
    exit 1
fi

# 템플릿 백업
cp "$DEPLOY_CONFIG_TEMPLATE" "${DEPLOY_CONFIG_TEMPLATE}.backup-script"
log_info "  템플릿 백업 완료"

# 템플릿에 prestate hash만 업데이트
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/\"faultGameAbsolutePrestate\": \"0x[^\"]*\"/\"faultGameAbsolutePrestate\": \"$PRESTATE_HASH\"/" "$DEPLOY_CONFIG_TEMPLATE"
else
    sed -i "s/\"faultGameAbsolutePrestate\": \"0x[^\"]*\"/\"faultGameAbsolutePrestate\": \"$PRESTATE_HASH\"/" "$DEPLOY_CONFIG_TEMPLATE"
fi

cd "$PROJECT_ROOT"
if make devnet-allocs > /tmp/devnet-allocs-1st.log 2>&1; then
    log_success "1차 .devnet 파일 생성 완료"
else
    log_error "1차 devnet-allocs 실패. 로그: /tmp/devnet-allocs-1st.log"
    mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
    exit 1
fi
echo ""

# 8. Genesis Output Root 계산
log_info "8/9: Genesis Output Root 계산 중..."
GENESIS_OUTPUT_ROOT=$(go run scripts/calc-genesis-output-root.go 2>/dev/null)
CALC_EXIT_CODE=$?
if [ $CALC_EXIT_CODE -ne 0 ]; then
    log_error "Genesis output root 계산 실패 (exit code: $CALC_EXIT_CODE)"
    go run scripts/calc-genesis-output-root.go 2>&1 | head -20  # 에러 출력
    mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
    exit 1
fi
log_success "Genesis Output Root: $GENESIS_OUTPUT_ROOT"

# 템플릿에 genesis output root도 업데이트
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/\"faultGameGenesisOutputRoot\": \"0x[^\"]*\"/\"faultGameGenesisOutputRoot\": \"$GENESIS_OUTPUT_ROOT\"/" "$DEPLOY_CONFIG_TEMPLATE"
else
    sed -i "s/\"faultGameGenesisOutputRoot\": \"0x[^\"]*\"/\"faultGameGenesisOutputRoot\": \"$GENESIS_OUTPUT_ROOT\"/" "$DEPLOY_CONFIG_TEMPLATE"
fi
log_success "  템플릿에 prestate hash + genesis output root 적용 완료"
echo ""

# 9. 2차 devnet-allocs 실행 (올바른 값으로)
log_info "9/9: 최종 .devnet 파일 재생성 중..."
log_warning "  이 단계는 약 5분 소요됩니다..."
rm -rf .devnet

if make devnet-allocs > /tmp/devnet-allocs-2nd.log 2>&1; then
    log_success "최종 .devnet 파일 생성 완료"

    # 생성된 devnetL1.json 검증
    DEPLOY_CONFIG="$PROJECT_ROOT/packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json"
    if [ -f "$DEPLOY_CONFIG" ]; then
        DEPLOYED_PRESTATE=$(cat "$DEPLOY_CONFIG" | jq -r '.faultGameAbsolutePrestate')
        DEPLOYED_GENESIS_ROOT=$(cat "$DEPLOY_CONFIG" | jq -r '.faultGameGenesisOutputRoot')

        log_info "  검증 중..."
        log_info "    Prestate Hash:         $DEPLOYED_PRESTATE"
        log_info "    Genesis Output Root:   $DEPLOYED_GENESIS_ROOT"

        if [ "$DEPLOYED_PRESTATE" == "$PRESTATE_HASH" ] && [ "$DEPLOYED_GENESIS_ROOT" == "$GENESIS_OUTPUT_ROOT" ]; then
            log_success "  ✓ 모든 값 일치 확인!"
        else
            log_error "  ✗ 값 불일치"
            [ "$DEPLOYED_PRESTATE" != "$PRESTATE_HASH" ] && log_error "    Prestate Hash 불일치"
            [ "$DEPLOYED_GENESIS_ROOT" != "$GENESIS_OUTPUT_ROOT" ] && log_error "    Genesis Output Root 불일치"
            mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
            exit 1
        fi
    else
        log_error "devnetL1.json이 생성되지 않았습니다"
        mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
        exit 1
    fi
else
    log_error "2차 devnet-allocs 실패. 로그: /tmp/devnet-allocs-2nd.log"
    mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
    exit 1
fi

# 템플릿 복원
mv "${DEPLOY_CONFIG_TEMPLATE}.backup-script" "$DEPLOY_CONFIG_TEMPLATE"
log_info "  템플릿 원래 상태로 복원 완료"
echo ""

# 10. E2E 테스트 전 최종 검증
log_info "=========================================="
log_info "E2E 테스트 전 최종 검증"
log_info "=========================================="

# prestate-proof-mt64.json과 devnetL1.json의 일치 확인
PROOF_PRESTATE=$(cat op-program/bin/prestate-proof-mt64.json | jq -r '.pre')
CONFIG_PRESTATE=$(cat "$DEPLOY_CONFIG" | jq -r '.faultGameAbsolutePrestate')
CONFIG_GENESIS_ROOT=$(cat "$DEPLOY_CONFIG" | jq -r '.faultGameGenesisOutputRoot')

log_info "Prestate Hash (proof):        $PROOF_PRESTATE"
log_info "Prestate Hash (config):       $CONFIG_PRESTATE"
log_info "Genesis Output Root (config): $CONFIG_GENESIS_ROOT"

ALL_MATCH=true
if [ "$PROOF_PRESTATE" != "$CONFIG_PRESTATE" ]; then
    log_error "✗ Prestate hash 불일치!"
    ALL_MATCH=false
fi

if [ "$ALL_MATCH" == "true" ]; then
    log_success "✓ 모든 설정 검증 완료!"
else
    log_error "E2E 테스트 전에 문제를 해결해야 합니다."
    exit 1
fi
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
echo "  - .devnet/allocs-l1.json (올바른 genesis output root 포함)"
echo "  - .devnet/allocs-l2-*.json"
echo "  - .devnet/addresses.json"
echo ""
log_info "검증된 값:"
echo "  - Prestate Hash:         $PRESTATE_HASH"
echo "  - Genesis Output Root:   $GENESIS_OUTPUT_ROOT"
echo ""

