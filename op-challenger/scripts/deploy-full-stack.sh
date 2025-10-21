#!/usr/bin/env bash
set -euo pipefail

##############################################################################
# L2 시스템 완전 배포 스크립트 (독립적인 Challenger 스택 포함)
#
# 이 스크립트는 다음을 배포합니다:
# - L1 Ethereum (로컬 또는 원격)
# - Sequencer 스택: sequencer-op-node, sequencer-l2, op-batcher, op-proposer
# - Challenger 스택: challenger-op-node, challenger-l2, op-challenger (⭐ 독립적!)
#
# 참고: /op-challenger/docs/l2-system-deployment-ko.md
##############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# 로그 파일 설정
LOG_DIR="${PROJECT_ROOT}/.devnet/logs"
mkdir -p "$LOG_DIR"
DEPLOY_LOG="${LOG_DIR}/deploy-$(date +%Y%m%d-%H%M%S).log"
CONFIG_TRACE_LOG="${LOG_DIR}/config-trace-$(date +%Y%m%d-%H%M%S).log"

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수 (콘솔 + 파일)
log_info() {
    local msg="$1"
    echo -e "${BLUE}[INFO]${NC} $msg"
    echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
}

log_success() {
    local msg="$1"
    echo -e "${GREEN}[SUCCESS]${NC} $msg"
    echo "[SUCCESS] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
}

log_warn() {
    local msg="$1"
    echo -e "${YELLOW}[WARN]${NC} $msg"
    echo "[WARN] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
}

log_error() {
    local msg="$1"
    echo -e "${RED}[ERROR]${NC} $msg"
    echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
}

# 설정 추적 로그 (config-trace.log 전용)
log_config() {
    local msg="$1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $msg" >> "$CONFIG_TRACE_LOG"
}

# 헬프 메시지
usage() {
    cat << EOF
L2 시스템 완전 배포 스크립트 (독립적인 Challenger 스택 포함)

사용법:
    $0 [옵션]

옵션:
    --env-file FILE         환경 변수 파일 경로 (기본: .env)
    --mode MODE             배포 모드: local|sepolia|mainnet (기본: local)
    --skip-build            Docker 이미지 빌드 건너뛰기
    --skip-l1               L1 배포 건너뛰기 (외부 RPC 사용 시)
    --with-indexer          Indexer 스택도 함께 배포
    --help                  이 도움말 표시

환경 변수 (로컬 모드 전용):
    FAULT_GAME_MAX_CLOCK_DURATION    게임 진행 최대 시간 (초) (기본: 1200 = 20분)
    FAULT_GAME_WITHDRAWAL_DELAY      게임 종료 후 인출 대기 시간 (초) (기본: 604800 = 7일)
    PROPOSAL_INTERVAL                Proposer 게임 생성 간격 (기본: 12s)
    DG_TYPE                          게임 타입 (0=CANNON, 1=PERMISSIONED) (기본: 0)

예시:
    # 로컬 개발 환경 전체 배포
    $0 --mode local

    # 게임 진행 시간을 5분으로 설정하여 배포
    FAULT_GAME_MAX_CLOCK_DURATION=300 $0 --mode local

    # 게임 진행 1분, 인출 대기 1일, 제안 간격 30초로 빠른 테스트 환경 배포
    FAULT_GAME_MAX_CLOCK_DURATION=60 FAULT_GAME_WITHDRAWAL_DELAY=86400 PROPOSAL_INTERVAL=30s $0 --mode local

    # Sepolia 테스트넷에 배포 (L1 건너뛰기)
    $0 --mode sepolia --skip-l1

    # Indexer 포함 배포
    $0 --mode local --with-indexer

EOF
    exit 0
}

# 기본값 설정
ENV_FILE="${PROJECT_ROOT}/.env"
DEPLOY_MODE="local"
SKIP_BUILD=false
SKIP_L1=false
WITH_INDEXER=false
DOCKER_COMPOSE_FILE="${SCRIPT_DIR}/docker-compose-full.yml"

# 인자 파싱
while [[ $# -gt 0 ]]; do
    case $1 in
        --env-file)
            ENV_FILE="$2"
            shift 2
            ;;
        --mode)
            DEPLOY_MODE="$2"
            shift 2
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-l1)
            SKIP_L1=true
            shift
            ;;
        --with-indexer)
            WITH_INDEXER=true
            shift
            ;;
        --help)
            usage
            ;;
        *)
            log_error "알 수 없는 옵션: $1"
            usage
            ;;
    esac
done

##############################################################################
# 1. 사전 검증 및 초기화 상태 확인
##############################################################################

log_info "배포 사전 검증을 시작합니다..."

# 필수 도구 확인
check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "$1이(가) 설치되어 있지 않습니다."
        exit 1
    fi
}

check_command docker
check_command jq
check_command openssl

# Docker Compose 버전 확인 (v1: docker-compose, v2: docker compose)
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
    log_info "docker-compose (v1) 사용"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
    log_info "docker compose (v2) 사용"
else
    log_error "docker-compose 또는 docker compose가 설치되어 있지 않습니다."
    log_info "Docker Desktop을 설치하거나 docker-compose를 설치하세요."
    exit 1
fi

log_success "필수 도구가 모두 설치되어 있습니다."

##############################################################################
# 초기화 상태 확인 (컨테이너 정리)
##############################################################################

log_info "Stopping any running containers from previous deployment..."

# Stop and remove all containers from previous deployment
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" down 2>/dev/null || true

log_success "Previous containers stopped"
log_info "Note: Volumes will be checked and removed before each service starts to ensure Genesis consistency"
echo ""

# 환경 변수 파일 확인
if [ ! -f "$ENV_FILE" ]; then
    log_warn "환경 변수 파일이 없습니다: $ENV_FILE"
    log_info "환경 변수를 자동으로 생성합니다..."
    echo ""

    # setup-env.sh 자동 실행
    if [ -f "${SCRIPT_DIR}/setup-env.sh" ]; then
        log_info "Running setup-env.sh..."
        if bash "${SCRIPT_DIR}/setup-env.sh" --mode "$DEPLOY_MODE" --output "$ENV_FILE"; then
            log_success "Environment variables generated successfully"
            echo ""
        else
            log_error "Failed to generate environment variables"
            log_info "Please run manually: ${SCRIPT_DIR}/setup-env.sh --mode $DEPLOY_MODE"
            exit 1
        fi
    else
        log_error "setup-env.sh not found: ${SCRIPT_DIR}/setup-env.sh"
        exit 1
    fi
fi

# Docker Compose 파일 확인
if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
    log_error "Docker Compose 파일을 찾을 수 없습니다: $DOCKER_COMPOSE_FILE"
    exit 1
fi

# 환경 변수 로드
log_info "환경 변수를 로드합니다: $ENV_FILE"
set -a
source "$ENV_FILE"
set +a

# PROJECT_ROOT 환경 변수 설정 (docker-compose에서 사용)
export PROJECT_ROOT

# Docker Compose에서 사용하는 추가 환경 변수 설정
export L2_IMAGE="${L2_IMAGE:-tokamaknetwork/thanos-op-geth:nightly}"
export IMAGE_TAG="${IMAGE_TAG:-nightly}"
export BATCH_TYPE="${BATCH_TYPE:-0}"
export PLASMA_ENABLED="${PLASMA_ENABLED:-false}"
export PLASMA_DA_SERVICE="${PLASMA_DA_SERVICE:-}"
export PLASMA_DA_SERVER="${PLASMA_DA_SERVER:-}"
export CHALLENGER_TRACE_TYPE="${CHALLENGER_TRACE_TYPE:-cannon}"
export DG_TYPE="${DG_TYPE:-0}"
export PROPOSAL_INTERVAL="${PROPOSAL_INTERVAL:-12s}"

# Define paths for contract addresses (will be loaded after Genesis generation)
ADDRESSES_FILE="${PROJECT_ROOT}/.devnet/addresses.json"
ZERO_ADDRESS="0x0000000000000000000000000000000000000000"

log_info "Contract addresses will be loaded from addresses.json after Genesis verification"

# 필수 환경 변수 명시적 export (디버깅 및 명확성 향상)
export MNEMONIC="${MNEMONIC}"
export ADMIN_PRIVATE_KEY="${ADMIN_PRIVATE_KEY}"
export BATCHER_PRIVATE_KEY="${BATCHER_PRIVATE_KEY}"
export PROPOSER_PRIVATE_KEY="${PROPOSER_PRIVATE_KEY}"
export SEQUENCER_PRIVATE_KEY="${SEQUENCER_PRIVATE_KEY}"
export CHALLENGER_PRIVATE_KEY="${CHALLENGER_PRIVATE_KEY}"

# 필수 환경 변수 검증
if [ -z "${MNEMONIC:-}" ]; then
    log_error "MNEMONIC environment variable is not set!"
    log_info "Please run setup-env.sh first to generate environment variables."
    exit 1
fi

if [ -z "${BATCHER_PRIVATE_KEY:-}" ] || [ -z "${PROPOSER_PRIVATE_KEY:-}" ] || [ -z "${CHALLENGER_PRIVATE_KEY:-}" ]; then
    log_warn "Some private keys are not set. This may cause deployment issues."
    log_info "Batcher Key: ${BATCHER_PRIVATE_KEY:+SET}${BATCHER_PRIVATE_KEY:-NOT SET}"
    log_info "Proposer Key: ${PROPOSER_PRIVATE_KEY:+SET}${PROPOSER_PRIVATE_KEY:-NOT SET}"
    log_info "Challenger Key: ${CHALLENGER_PRIVATE_KEY:+SET}${CHALLENGER_PRIVATE_KEY:-NOT SET}"
fi

log_success "Account environment variables loaded"
log_info "Contract addresses will be validated after Genesis generation"

##############################################################################
# 2. 배포 모드별 설정
##############################################################################

log_info "배포 모드: $DEPLOY_MODE"
log_info "=========================================="
echo ""
log_info "📝 로그 파일:"
log_info "  배포 로그: $DEPLOY_LOG"
log_info "  설정 추적 로그: $CONFIG_TRACE_LOG"
echo ""

# 환경 변수 로그 기록
log_config "=== 스크립트 시작 ==="
log_config "Deployment mode: $DEPLOY_MODE"
log_config "Environment variables:"
log_config "  FAULT_GAME_MAX_CLOCK_DURATION: ${FAULT_GAME_MAX_CLOCK_DURATION:-not set}"
log_config "  FAULT_GAME_WITHDRAWAL_DELAY: ${FAULT_GAME_WITHDRAWAL_DELAY:-not set}"
log_config "  PROPOSAL_INTERVAL: ${PROPOSAL_INTERVAL:-not set}"
log_config "  DG_TYPE: ${DG_TYPE:-not set}"

case $DEPLOY_MODE in
    local)
        log_info "로컬 개발 환경으로 배포합니다."
        export L1_RPC_URL="${L1_RPC_URL:-http://localhost:8545}"
        export L1_BEACON_URL="${L1_BEACON_URL:-}"
        SKIP_L1=false

        # cast 명령어 확인 (로컬 모드에서 ETH 전송에 필요)
        check_command cast
        ;;
    sepolia)
        log_info "Sepolia 테스트넷에 배포합니다."
        if [ -z "${L1_RPC_URL:-}" ]; then
            log_error "Sepolia 배포 시 L1_RPC_URL이 필요합니다."
            exit 1
        fi
        if [ -z "${L1_BEACON_URL:-}" ]; then
            log_error "Sepolia 배포 시 L1_BEACON_URL이 필요합니다."
            exit 1
        fi
        SKIP_L1=true
        ;;
    mainnet)
        log_info "Mainnet에 배포합니다."
        log_warn "⚠️  프로덕션 환경입니다. 신중하게 진행하세요!"
        read -p "계속하시겠습니까? (yes/no): " confirm
        if [ "$confirm" != "yes" ]; then
            log_info "배포를 취소합니다."
            exit 0
        fi
        if [ -z "${L1_RPC_URL:-}" ]; then
            log_error "Mainnet 배포 시 L1_RPC_URL이 필요합니다."
            exit 1
        fi
        if [ -z "${L1_BEACON_URL:-}" ]; then
            log_error "Mainnet 배포 시 L1_BEACON_URL이 필요합니다."
            exit 1
        fi
        SKIP_L1=true
        ;;
    *)
        log_error "지원하지 않는 배포 모드: $DEPLOY_MODE"
        log_info "사용 가능한 모드: local, sepolia, mainnet"
        exit 1
        ;;
esac

##############################################################################
# 3. 필수 파일 생성
##############################################################################

log_info "필수 파일을 생성합니다..."

DEVNET_DIR="${PROJECT_ROOT}/.devnet"
mkdir -p "$DEVNET_DIR"

# JWT 시크릿 생성 (없는 경우)
JWT_SECRET="${DEVNET_DIR}/jwt-secret.txt"
if [ ! -f "$JWT_SECRET" ]; then
    log_info "JWT 시크릿을 생성합니다: $JWT_SECRET"
    openssl rand -hex 32 > "$JWT_SECRET"
    log_success "JWT 시크릿 생성 완료"
else
    log_info "기존 JWT 시크릿을 사용합니다: $JWT_SECRET"
fi

# P2P 키 생성 (없는 경우)
P2P_NODE_KEY="${DEVNET_DIR}/p2p-node-key.txt"
if [ ! -f "$P2P_NODE_KEY" ]; then
    log_info "P2P Node 키를 생성합니다: $P2P_NODE_KEY"
    openssl rand -hex 32 > "$P2P_NODE_KEY"
    log_success "P2P Node 키 생성 완료"
fi

P2P_CHALLENGER_KEY="${DEVNET_DIR}/p2p-challenger-key.txt"
if [ ! -f "$P2P_CHALLENGER_KEY" ]; then
    log_info "P2P Challenger 키를 생성합니다: $P2P_CHALLENGER_KEY"
    openssl rand -hex 32 > "$P2P_CHALLENGER_KEY"
    log_success "P2P Challenger 키 생성 완료"
fi

# Genesis 파일 확인 (로컬 모드인 경우 필수)
if [ "$DEPLOY_MODE" = "local" ]; then
    GENESIS_L1="${DEVNET_DIR}/genesis-l1.json"
    GENESIS_L2="${DEVNET_DIR}/genesis-l2.json"
    ROLLUP_CONFIG="${DEVNET_DIR}/rollup.json"

    # Check for directories incorrectly created by Docker volume mounts
    CLEANUP_NEEDED=false
    for genesis_file in "$GENESIS_L1" "$GENESIS_L2" "$ROLLUP_CONFIG"; do
        if [ -d "$genesis_file" ]; then
            log_warn "Found directory instead of file: $genesis_file"
            log_warn "Removing incorrectly created directory..."
            rm -rf "$genesis_file"
            CLEANUP_NEEDED=true
        fi
    done

    if [ "$CLEANUP_NEEDED" = true ]; then
        log_success "Cleaned up incorrectly created directories"
        log_info "This usually happens when Docker volume mount creates directories for missing files"
        echo ""
    fi

    # Check if game configuration changed (via environment variables)
    # If so, remove existing Genesis and deployment cache to force regeneration
    if [ -n "${FAULT_GAME_MAX_CLOCK_DURATION:-}" ] || [ -n "${FAULT_GAME_WITHDRAWAL_DELAY:-}" ]; then
        log_warn "⚙️  Game configuration environment variables detected!"
        log_info "  FAULT_GAME_MAX_CLOCK_DURATION: ${FAULT_GAME_MAX_CLOCK_DURATION:-default}"
        log_info "  FAULT_GAME_WITHDRAWAL_DELAY: ${FAULT_GAME_WITHDRAWAL_DELAY:-default}"
        echo ""
        log_warn "Will apply custom game settings during Genesis generation."
        echo ""

        # Always remove deployment cache and Genesis when custom settings are specified
        # This ensures new settings are applied even on first run
        log_info "🧹 Removing old files to apply new game settings..."
        echo ""

        # Remove deployment cache (critical!)
        log_info "  1️⃣  Removing deployment cache..."
        if [ -d "${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deployments/devnetL1" ]; then
            rm -rf "${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deployments/devnetL1"
            echo "     ✅ Removed: packages/tokamak/contracts-bedrock/deployments/devnetL1"
        else
            echo "     ℹ️  Not found: packages/tokamak/contracts-bedrock/deployments/devnetL1"
        fi

        # Remove devnetL1.json (critical!) - forces template to be used
        log_info "  2️⃣  Removing devnetL1.json (forces template to be used)..."
        DEVNET_CONFIG="${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json"
        if [ -f "$DEVNET_CONFIG" ]; then
            OLD_CLOCK=$(jq -r '.faultGameMaxClockDuration' "$DEVNET_CONFIG" 2>/dev/null || echo "unknown")
            OLD_DELAY=$(jq -r '.faultGameWithdrawalDelay' "$DEVNET_CONFIG" 2>/dev/null || echo "unknown")
            echo "     📋 Old config: clock=${OLD_CLOCK}s, delay=${OLD_DELAY}s"
            rm -f "$DEVNET_CONFIG"
            echo "     ✅ Removed: devnetL1.json"
        else
            echo "     ℹ️  Not found: devnetL1.json"
        fi

        # Remove existing Genesis files (if any)
        log_info "  3️⃣  Removing Genesis and contract files..."
        rm -f "$GENESIS_L1" "$GENESIS_L2" "$ROLLUP_CONFIG" "$ADDRESSES_FILE"
        rm -f "${DEVNET_DIR}/allocs"*.json
        echo "     ✅ Removed: genesis*.json, rollup.json, addresses.json, allocs*.json"

        # Remove Docker volumes to prevent Genesis hash mismatch
        log_info "  4️⃣  Removing Docker volumes..."
        VOLUMES_REMOVED=0
        for vol in scripts_l1_data scripts_sequencer_l2_data scripts_challenger_l2_data \
                   scripts_sequencer_safedb_data scripts_challenger_safedb_data \
                   scripts_challenger_data scripts_op_log; do
            if docker volume rm "$vol" 2>/dev/null; then
                VOLUMES_REMOVED=$((VOLUMES_REMOVED + 1))
            fi
        done
        echo "     ✅ Removed ${VOLUMES_REMOVED} Docker volume(s)"
        echo ""

        log_success "✅ Cleanup complete. Will regenerate with custom game settings."
        echo ""
    fi

    # Check for Genesis files AND addresses.json
    if [ ! -f "$GENESIS_L1" ] || [ ! -f "$GENESIS_L2" ] || [ ! -f "$ROLLUP_CONFIG" ] || [ ! -f "$ADDRESSES_FILE" ]; then
        log_warn "Required Genesis and contract files are missing!"
        echo ""
        log_info "Missing files:"
        [ ! -f "$GENESIS_L1" ] && echo "  ❌ ${GENESIS_L1}"
        [ ! -f "$GENESIS_L2" ] && echo "  ❌ ${GENESIS_L2}"
        [ ! -f "$ROLLUP_CONFIG" ] && echo "  ❌ ${ROLLUP_CONFIG}"
        [ ! -f "$ADDRESSES_FILE" ] && echo "  ❌ ${ADDRESSES_FILE} (contract addresses)"
        echo ""

        log_info "Automatically generating Genesis and contract files..."
        log_warn "This will:"
        echo "  1. Configure game settings (if specified)"
        echo "  2. Start devnet (make devnet-up)"
        echo "  3. Deploy contracts and generate Genesis (~2-3 minutes)"
        echo "  4. Stop devnet (make devnet-down)"
        echo "  5. Continue with deployment"
        echo ""

        cd "${PROJECT_ROOT}"

        # Configure game settings before deployment
        DEPLOY_CONFIG="${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deploy-config/devnetL1-template.json"
        log_config "=== Template 파일 경로 ==="
        log_config "DEPLOY_CONFIG: $DEPLOY_CONFIG"
        log_config "File exists: $([ -f "$DEPLOY_CONFIG" ] && echo 'YES' || echo 'NO')"

        if [ -f "$DEPLOY_CONFIG" ]; then
            # Backup original config
            if [ ! -f "${DEPLOY_CONFIG}.backup" ]; then
                cp "$DEPLOY_CONFIG" "${DEPLOY_CONFIG}.backup"
                log_info "Backup created: ${DEPLOY_CONFIG}.backup"
            fi

            # Check for environment variables to customize game configuration
            CONFIG_MODIFIED=false

            log_config "=== Template 수정 전 상태 ==="
            log_config "Original values:"
            ORIG_CLOCK=$(jq -r '.faultGameMaxClockDuration' "$DEPLOY_CONFIG" 2>/dev/null || echo "unknown")
            ORIG_DELAY=$(jq -r '.faultGameWithdrawalDelay' "$DEPLOY_CONFIG" 2>/dev/null || echo "unknown")
            log_config "  faultGameMaxClockDuration: ${ORIG_CLOCK}s"
            log_config "  faultGameWithdrawalDelay: ${ORIG_DELAY}s"

            if [ -n "${FAULT_GAME_MAX_CLOCK_DURATION:-}" ]; then
                log_info "Configuring faultGameMaxClockDuration: ${FAULT_GAME_MAX_CLOCK_DURATION}s ($((FAULT_GAME_MAX_CLOCK_DURATION / 60)) min)"
                log_config "Modifying faultGameMaxClockDuration: ${ORIG_CLOCK} -> ${FAULT_GAME_MAX_CLOCK_DURATION}"

                # Use jq to update the value
                if command -v jq >/dev/null 2>&1; then
                    TEMP_CONFIG=$(mktemp)
                    jq ".faultGameMaxClockDuration = ${FAULT_GAME_MAX_CLOCK_DURATION}" "$DEPLOY_CONFIG" > "$TEMP_CONFIG"
                    mv "$TEMP_CONFIG" "$DEPLOY_CONFIG"
                    CONFIG_MODIFIED=true

                    # Verify modification
                    NEW_CLOCK=$(jq -r '.faultGameMaxClockDuration' "$DEPLOY_CONFIG")
                    log_config "Verification: faultGameMaxClockDuration = ${NEW_CLOCK}s"
                    if [ "$NEW_CLOCK" = "${FAULT_GAME_MAX_CLOCK_DURATION}" ]; then
                        log_config "✅ faultGameMaxClockDuration successfully updated"
                    else
                        log_config "❌ faultGameMaxClockDuration update FAILED!"
                    fi
                else
                    log_warn "jq not found. Cannot modify faultGameMaxClockDuration"
                    log_config "❌ jq not available - modification skipped"
                fi
            fi

            if [ -n "${FAULT_GAME_WITHDRAWAL_DELAY:-}" ]; then
                log_info "Configuring faultGameWithdrawalDelay: ${FAULT_GAME_WITHDRAWAL_DELAY}s ($((FAULT_GAME_WITHDRAWAL_DELAY / 86400)) days)"
                log_config "Modifying faultGameWithdrawalDelay: ${ORIG_DELAY} -> ${FAULT_GAME_WITHDRAWAL_DELAY}"

                if command -v jq >/dev/null 2>&1; then
                    TEMP_CONFIG=$(mktemp)
                    jq ".faultGameWithdrawalDelay = ${FAULT_GAME_WITHDRAWAL_DELAY}" "$DEPLOY_CONFIG" > "$TEMP_CONFIG"
                    mv "$TEMP_CONFIG" "$DEPLOY_CONFIG"
                    CONFIG_MODIFIED=true

                    # Verify modification
                    NEW_DELAY=$(jq -r '.faultGameWithdrawalDelay' "$DEPLOY_CONFIG")
                    log_config "Verification: faultGameWithdrawalDelay = ${NEW_DELAY}s"
                    if [ "$NEW_DELAY" = "${FAULT_GAME_WITHDRAWAL_DELAY}" ]; then
                        log_config "✅ faultGameWithdrawalDelay successfully updated"
                    else
                        log_config "❌ faultGameWithdrawalDelay update FAILED!"
                    fi
                fi
            fi

            if [ -n "${PROPOSAL_INTERVAL:-}" ]; then
                log_info "Proposal interval will be set to: ${PROPOSAL_INTERVAL}"
                log_config "PROPOSAL_INTERVAL: ${PROPOSAL_INTERVAL}"
            fi

            if [ "$CONFIG_MODIFIED" = true ]; then
                log_success "Game configuration updated in template"
                log_config "=== Template 수정 완료 ==="
                log_config "Modified template file: $DEPLOY_CONFIG"
                echo ""
            fi
        fi

        # Start devnet
        log_info "Step 1/3: Starting devnet (make devnet-up)..."
        log_info "This will take 2-3 minutes. Please wait..."
        echo ""

        if make devnet-up; then
            log_success "Devnet started successfully"
            echo ""

            # Wait for Genesis files and addresses.json to be created
            log_info "Step 2/3: Waiting for Genesis and contract files to be created..."
            max_wait=180  # 3 minutes
            elapsed=0

            while [ $elapsed -lt $max_wait ]; do
                if [ -f "$GENESIS_L1" ] && [ -f "$GENESIS_L2" ] && [ -f "$ROLLUP_CONFIG" ] && [ -f "$ADDRESSES_FILE" ]; then
                    log_success "Genesis and contract files detected!"
                    break
                fi
                sleep 5
                elapsed=$((elapsed + 5))
                echo -n "."
            done
            echo ""

            if [ -f "$GENESIS_L1" ] && [ -f "$GENESIS_L2" ] && [ -f "$ROLLUP_CONFIG" ] && [ -f "$ADDRESSES_FILE" ]; then
                log_success "All Genesis and contract files created successfully!"
                echo ""
                log_info "Generated files:"
                ls -lh .devnet/genesis*.json .devnet/rollup.json .devnet/addresses.json
                echo ""

                # Verify devnetL1.json was created with correct values
                log_config "=== devnetL1.json 생성 확인 (make devnet-up 후) ==="
                DEVNET_CONFIG_ACTUAL="${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json"
                if [ -f "$DEVNET_CONFIG_ACTUAL" ]; then
                    ACTUAL_CLOCK=$(jq -r '.faultGameMaxClockDuration' "$DEVNET_CONFIG_ACTUAL" 2>/dev/null || echo "unknown")
                    ACTUAL_DELAY=$(jq -r '.faultGameWithdrawalDelay' "$DEVNET_CONFIG_ACTUAL" 2>/dev/null || echo "unknown")
                    log_config "devnetL1.json values:"
                    log_config "  faultGameMaxClockDuration: ${ACTUAL_CLOCK}s"
                    log_config "  faultGameWithdrawalDelay: ${ACTUAL_DELAY}s"

                    # Compare with expected values
                    if [ -n "${FAULT_GAME_MAX_CLOCK_DURATION:-}" ]; then
                        if [ "$ACTUAL_CLOCK" = "${FAULT_GAME_MAX_CLOCK_DURATION}" ]; then
                            log_config "✅ Clock duration matches expected: ${FAULT_GAME_MAX_CLOCK_DURATION}s"
                        else
                            log_config "❌ Clock duration mismatch! Expected: ${FAULT_GAME_MAX_CLOCK_DURATION}s, Actual: ${ACTUAL_CLOCK}s"
                        fi
                    fi

                    if [ -n "${FAULT_GAME_WITHDRAWAL_DELAY:-}" ]; then
                        if [ "$ACTUAL_DELAY" = "${FAULT_GAME_WITHDRAWAL_DELAY}" ]; then
                            log_config "✅ Withdrawal delay matches expected: ${FAULT_GAME_WITHDRAWAL_DELAY}s"
                        else
                            log_config "❌ Withdrawal delay mismatch! Expected: ${FAULT_GAME_WITHDRAWAL_DELAY}s, Actual: ${ACTUAL_DELAY}s"
                        fi
                    fi
                else
                    log_config "❌ devnetL1.json not found at: $DEVNET_CONFIG_ACTUAL"
                fi

                # Verify deployed contract settings
                log_config "=== Deployed Contract 설정 확인 (addresses.json) ==="
                if [ -f "$ADDRESSES_FILE" ]; then
                    DGF_ADDR=$(jq -r '.DisputeGameFactoryProxy // empty' "$ADDRESSES_FILE" 2>/dev/null || echo "")
                    log_config "DisputeGameFactoryProxy: ${DGF_ADDR}"
                    if [ -n "$DGF_ADDR" ] && [ "$DGF_ADDR" != "null" ]; then
                        log_config "✅ Contract deployed successfully"
                    else
                        log_config "❌ DisputeGameFactoryProxy address not found!"
                    fi
                fi

                # Stop devnet and remove volumes to ensure clean state
                log_info "Step 3/3: Stopping devnet and removing volumes (make devnet-down)..."
                log_warn "Removing devnet volumes to ensure L1 uses the same genesis hash..."

                # Stop devnet containers and remove volumes
                cd "${PROJECT_ROOT}/ops-bedrock"
                if [ -f "docker-compose.yml" ]; then
                    $DOCKER_COMPOSE -f docker-compose.yml down -v 2>/dev/null || true
                fi
                cd "${PROJECT_ROOT}"

                log_success "Genesis generation complete! Continuing with deployment..."
                echo ""
            else
                log_error "Genesis files were not created within ${max_wait} seconds"
                log_info "Stopping devnet..."
                make devnet-down || true
                echo ""
                log_error "Genesis generation failed. Please check devnet logs."
                exit 1
            fi
        else
            log_error "Failed to start devnet"
            log_info "Please try manually:"
            echo "  cd ${PROJECT_ROOT}"
            echo "  make devnet-up"
            echo "  make devnet-down"
            exit 1
        fi
    else
        log_success "Genesis files verified"
    fi

    # Load contract addresses from addresses.json (after Genesis generation)
    log_info "Loading contract addresses from addresses.json..."

    if [ -f "$ADDRESSES_FILE" ]; then
        # DisputeGameFactoryProxy 주소 읽기
        GAME_FACTORY_ADDRESS=$(jq -r '.DisputeGameFactoryProxy // .DisputeGameFactory' "$ADDRESSES_FILE" 2>/dev/null)
        if [ "$GAME_FACTORY_ADDRESS" != "null" ] && [ -n "$GAME_FACTORY_ADDRESS" ] && [ "$GAME_FACTORY_ADDRESS" != "$ZERO_ADDRESS" ]; then
            export GAME_FACTORY_ADDRESS
            log_success "GAME_FACTORY_ADDRESS loaded: $GAME_FACTORY_ADDRESS"
        else
            log_error "Failed to load DisputeGameFactory address from addresses.json"
            cat "$ADDRESSES_FILE" | jq '.' | head -20
            exit 1
        fi

        # L2OutputOracleProxy 주소 읽기
        L2OO_ADDRESS=$(jq -r '.L2OutputOracleProxy // .L2OutputOracle' "$ADDRESSES_FILE" 2>/dev/null)
        if [ "$L2OO_ADDRESS" != "null" ] && [ -n "$L2OO_ADDRESS" ] && [ "$L2OO_ADDRESS" != "$ZERO_ADDRESS" ]; then
            export L2OO_ADDRESS
            log_info "L2OO_ADDRESS loaded: $L2OO_ADDRESS"
        fi
    else
        log_error "addresses.json not found after Genesis generation!"
        log_error "This should not happen. Please check devnet-up logs."
        exit 1
    fi

    log_success "All contract addresses loaded successfully"
    echo ""
fi

##############################################################################
# 4. Cannon VM 및 op-program 빌드 (현재 환경에 맞는 prestate 생성)
##############################################################################

log_info "=========================================="
log_info "Cannon VM 및 op-program 빌드"
log_info "=========================================="

# op-program/bin 디렉토리 확인
OP_PROGRAM_BIN="${PROJECT_ROOT}/op-program/bin"
CANNON_BIN="${PROJECT_ROOT}/cannon/bin"
PRESTATE_FILE="${OP_PROGRAM_BIN}/prestate.json"
PRESTATE_PROOF_FILE="${OP_PROGRAM_BIN}/prestate-proof.json"

# Genesis 해시 캐시 파일
GENESIS_HASH_CACHE="${OP_PROGRAM_BIN}/.genesis-hash"

# 기존 prestate 파일 확인
if [ -f "$PRESTATE_FILE" ] && [ -f "$PRESTATE_PROOF_FILE" ]; then
    log_info "기존 prestate 파일 발견:"
    ls -lh "$PRESTATE_FILE" "$PRESTATE_PROOF_FILE" 2>/dev/null || true

    # Genesis 파일 내용 해시로 변경 감지 (로컬 모드인 경우)
    if [ "$DEPLOY_MODE" = "local" ] && [ -f "${DEVNET_DIR}/genesis-l2.json" ]; then
        # 현재 Genesis 파일의 해시 계산
        CURRENT_GENESIS_HASH=$(shasum -a 256 "${DEVNET_DIR}/genesis-l2.json" 2>/dev/null | awk '{print $1}' || \
                               sha256sum "${DEVNET_DIR}/genesis-l2.json" 2>/dev/null | awk '{print $1}')

        # 이전 빌드 시 사용한 Genesis 해시 확인
        if [ -f "$GENESIS_HASH_CACHE" ]; then
            CACHED_GENESIS_HASH=$(cat "$GENESIS_HASH_CACHE")

            if [ "$CURRENT_GENESIS_HASH" = "$CACHED_GENESIS_HASH" ]; then
                log_success "✅ Genesis 파일이 변경되지 않았습니다."
                log_info "   현재 해시: ${CURRENT_GENESIS_HASH:0:16}..."
                REBUILD_PRESTATE=false
            else
                log_warn "⚠️  Genesis 파일이 변경되었습니다!"
                log_info "   이전 해시: ${CACHED_GENESIS_HASH:0:16}..."
                log_info "   현재 해시: ${CURRENT_GENESIS_HASH:0:16}..."
                log_info "   현재 환경에 맞는 새로운 prestate를 생성해야 합니다."
                REBUILD_PRESTATE=true
            fi
        else
            log_warn "Genesis 해시 캐시가 없습니다."
            log_info "안전을 위해 prestate를 재생성하는 것을 권장합니다."
            REBUILD_PRESTATE=true
        fi

        # Rollup config도 확인
        if [ -f "${DEVNET_DIR}/rollup.json" ]; then
            ROLLUP_HASH_CACHE="${OP_PROGRAM_BIN}/.rollup-hash"
            CURRENT_ROLLUP_HASH=$(shasum -a 256 "${DEVNET_DIR}/rollup.json" 2>/dev/null | awk '{print $1}' || \
                                  sha256sum "${DEVNET_DIR}/rollup.json" 2>/dev/null | awk '{print $1}')

            if [ -f "$ROLLUP_HASH_CACHE" ]; then
                CACHED_ROLLUP_HASH=$(cat "$ROLLUP_HASH_CACHE")
                if [ "$CURRENT_ROLLUP_HASH" != "$CACHED_ROLLUP_HASH" ]; then
                    log_warn "⚠️  Rollup config가 변경되었습니다!"
                    REBUILD_PRESTATE=true
                fi
            fi
        fi
    else
        # Genesis 파일 없으면 항상 재빌드
        log_warn "Genesis 파일을 찾을 수 없거나 로컬 모드가 아닙니다."
        REBUILD_PRESTATE=true
    fi

    # 사용자에게 재빌드 여부 확인
    if [ "$REBUILD_PRESTATE" = true ]; then
        log_warn "현재 환경에 맞는 prestate를 다시 생성하시겠습니까?"
        log_info "이 작업은 약 5-10분 소요됩니다."
        read -p "Rebuild prestate? (yes/no) [yes]: " rebuild_confirm
        rebuild_confirm=${rebuild_confirm:-yes}

        if [ "$rebuild_confirm" != "yes" ]; then
            log_info "기존 prestate를 사용합니다."
            REBUILD_PRESTATE=false
        fi
    else
        log_info "기존 prestate를 사용합니다. (변경 사항 없음)"
        REBUILD_PRESTATE=false
    fi
else
    log_warn "prestate 파일이 없습니다. 새로 생성합니다."
    REBUILD_PRESTATE=true
fi

# prestate 빌드 실행
if [ "$REBUILD_PRESTATE" = true ]; then
    log_info "현재 환경에 맞는 Cannon VM 및 op-program을 빌드합니다..."
    echo ""

    cd "${PROJECT_ROOT}"

    # 1. op-program 빌드
    log_info "[1/3] op-program 빌드 중... (약 2-5분)"
    if make op-program; then
        log_success "op-program 빌드 완료"
        echo ""
    else
        log_error "op-program 빌드 실패"
        log_info "수동으로 빌드를 시도하세요:"
        echo "  cd ${PROJECT_ROOT}"
        echo "  make op-program"
        exit 1
    fi

    # 2. cannon 빌드
    log_info "[2/3] Cannon VM 빌드 중... (약 1-2분)"
    if make cannon; then
        log_success "Cannon VM 빌드 완료"
        echo ""
    else
        log_error "Cannon VM 빌드 실패"
        log_info "수동으로 빌드를 시도하세요:"
        echo "  cd ${PROJECT_ROOT}"
        echo "  make cannon"
        exit 1
    fi

    # 3. prestate 생성
    log_info "[3/3] 현재 환경에 맞는 prestate 생성 중... (약 10-30초)"
    if make cannon-prestate; then
        log_success "Prestate 생성 완료!"
        echo ""
    else
        log_error "Prestate 생성 실패"
        log_info "수동으로 생성을 시도하세요:"
        echo "  cd ${PROJECT_ROOT}"
        echo "  make cannon-prestate"
        exit 1
    fi

    # 생성된 파일 확인
    if [ -f "$PRESTATE_FILE" ] && [ -f "$PRESTATE_PROOF_FILE" ]; then
        log_success "생성된 파일:"
        ls -lh "$PRESTATE_FILE" "$PRESTATE_PROOF_FILE" "$CANNON_BIN/cannon" "$OP_PROGRAM_BIN/op-program"
        echo ""

        # Absolute prestate hash 출력
        ABSOLUTE_PRESTATE=$(cat "$PRESTATE_PROOF_FILE" | jq -r '.pre' 2>/dev/null || echo "error")
        if [ "$ABSOLUTE_PRESTATE" != "error" ] && [ -n "$ABSOLUTE_PRESTATE" ]; then
            log_success "✅ Absolute Prestate Hash:"
            echo "   $ABSOLUTE_PRESTATE"
            echo ""
            log_info "이 해시가 DisputeGameFactory의 absolute prestate와 일치해야 합니다."
            log_info "op-challenger가 'Invalid prestate' 에러 없이 작동할 것입니다."
            echo ""
        else
            log_warn "Absolute prestate hash를 읽을 수 없습니다."
        fi

        # Genesis 및 Rollup config 해시를 캐시에 저장
        if [ "$DEPLOY_MODE" = "local" ]; then
            if [ -f "${DEVNET_DIR}/genesis-l2.json" ]; then
                CURRENT_GENESIS_HASH=$(shasum -a 256 "${DEVNET_DIR}/genesis-l2.json" 2>/dev/null | awk '{print $1}' || \
                                       sha256sum "${DEVNET_DIR}/genesis-l2.json" 2>/dev/null | awk '{print $1}')
                echo "$CURRENT_GENESIS_HASH" > "$GENESIS_HASH_CACHE"
                log_info "Genesis 해시 캐시 저장: .genesis-hash"
            fi

            if [ -f "${DEVNET_DIR}/rollup.json" ]; then
                CURRENT_ROLLUP_HASH=$(shasum -a 256 "${DEVNET_DIR}/rollup.json" 2>/dev/null | awk '{print $1}' || \
                                      sha256sum "${DEVNET_DIR}/rollup.json" 2>/dev/null | awk '{print $1}')
                echo "$CURRENT_ROLLUP_HASH" > "${OP_PROGRAM_BIN}/.rollup-hash"
                log_info "Rollup config 해시 캐시 저장: .rollup-hash"
            fi
            echo ""
        fi
    else
        log_error "Prestate 파일이 생성되지 않았습니다!"
        exit 1
    fi

    cd "${SCRIPT_DIR}"
else
    log_info "기존 prestate 사용:"

    # 현재 absolute prestate hash 출력
    if [ -f "$PRESTATE_PROOF_FILE" ]; then
        ABSOLUTE_PRESTATE=$(cat "$PRESTATE_PROOF_FILE" | jq -r '.pre' 2>/dev/null || echo "error")
        if [ "$ABSOLUTE_PRESTATE" != "error" ] && [ -n "$ABSOLUTE_PRESTATE" ]; then
            log_info "Current Absolute Prestate Hash:"
            echo "   $ABSOLUTE_PRESTATE"
            echo ""
        fi
    fi
fi

log_success "Cannon VM 및 op-program 준비 완료"
echo ""

##############################################################################
# 5. Docker 이미지 빌드
##############################################################################

if [ "$SKIP_BUILD" = false ]; then
    log_info "Docker 이미지를 빌드합니다..."

    # docker-compose build 실행
    $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" build --parallel

    log_success "Docker 이미지 빌드 완료"
else
    log_info "Docker 이미지 빌드를 건너뜁니다."
fi

##############################################################################
# 6. L1 배포
##############################################################################

if [ "$SKIP_L1" = false ]; then
    log_info "L1 Ethereum을 시작합니다..."

    # CRITICAL: Remove L1 volume to ensure it initializes with the correct genesis
    # L1 geth will ignore genesis.json if it finds existing blockchain data
    if docker volume inspect scripts_l1_data >/dev/null 2>&1; then
        log_warn "Removing existing L1 volume to ensure genesis hash consistency..."
        # Stop L1 if running
        $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" stop l1 2>/dev/null || true
        $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" rm -f l1 2>/dev/null || true
        # Remove L1 volume
        docker volume rm scripts_l1_data 2>/dev/null || true
        log_success "L1 volume removed. L1 will initialize from genesis-l1.json"
    fi

    $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d l1

    log_info "L1 RPC가 준비될 때까지 대기합니다..."
    timeout=120
    elapsed=0
    while ! curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "http://localhost:8545" > /dev/null 2>&1; do
        if [ $elapsed -ge $timeout ]; then
            log_error "L1 RPC 시작 시간 초과 (${timeout}초)"
            $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" logs l1
            exit 1
        fi
        sleep 2
        elapsed=$((elapsed + 2))
        echo -n "."
    done
    echo ""

    log_success "L1 Ethereum 시작 완료"

    # 로컬 모드: 필수 계정에 L1 ETH 전송
    if [ "$DEPLOY_MODE" = "local" ]; then
        log_info "필수 계정에 L1 ETH를 전송합니다..."

        # L1 dev 모드 기본 계정 (충분한 ETH 보유)
        DEV_PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

        # Proposer에 ETH 전송
        if [ -n "${PROPOSER_ADDRESS:-}" ]; then
            log_info "Proposer (${PROPOSER_ADDRESS})에 10 ETH 전송 중..."
            if cast send "${PROPOSER_ADDRESS}" --value 10ether --private-key "${DEV_PRIVATE_KEY}" --rpc-url http://localhost:8545 > /dev/null 2>&1; then
                log_success "Proposer에 ETH 전송 완료"
            else
                log_warn "Proposer ETH 전송 실패 (계속 진행)"
            fi
        fi

        # Batcher에 ETH 전송
        if [ -n "${BATCHER_ADDRESS:-}" ]; then
            log_info "Batcher (${BATCHER_ADDRESS})에 10 ETH 전송 중..."
            if cast send "${BATCHER_ADDRESS}" --value 10ether --private-key "${DEV_PRIVATE_KEY}" --rpc-url http://localhost:8545 > /dev/null 2>&1; then
                log_success "Batcher에 ETH 전송 완료"
            else
                log_warn "Batcher ETH 전송 실패 (계속 진행)"
            fi
        fi

        # Challenger에 ETH 전송
        if [ -n "${CHALLENGER_ADDRESS:-}" ]; then
            log_info "Challenger (${CHALLENGER_ADDRESS})에 10 ETH 전송 중..."
            if cast send "${CHALLENGER_ADDRESS}" --value 10ether --private-key "${DEV_PRIVATE_KEY}" --rpc-url http://localhost:8545 > /dev/null 2>&1; then
                log_success "Challenger에 ETH 전송 완료"
            else
                log_warn "Challenger ETH 전송 실패 (계속 진행)"
            fi
        fi

        log_success "L1 계정 준비 완료"
    fi
else
    log_info "L1 배포를 건너뜁니다 (외부 RPC 사용)"
    log_info "L1 RPC URL: $L1_RPC_URL"
fi

##############################################################################
# 7. Sequencer 스택 배포
##############################################################################

log_info "Sequencer 스택을 배포합니다..."

# CRITICAL: Remove Sequencer L2 volume to ensure it initializes with the correct genesis
if docker volume inspect scripts_sequencer_l2_data >/dev/null 2>&1; then
    log_warn "Removing existing Sequencer L2 volume to ensure genesis consistency..."
    $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" stop sequencer-l2 2>/dev/null || true
    $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" rm -f sequencer-l2 2>/dev/null || true
    docker volume rm scripts_sequencer_l2_data 2>/dev/null || true
    log_success "Sequencer L2 volume removed"
fi

# Sequencer L2 geth 시작
log_info "Sequencer L2 op-geth를 시작합니다..."
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d sequencer-l2

log_info "Sequencer L2 RPC가 준비될 때까지 대기합니다..."
timeout=120
elapsed=0
while ! curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "http://localhost:9545" > /dev/null 2>&1; do
    if [ $elapsed -ge $timeout ]; then
        log_error "Sequencer L2 RPC 시작 시간 초과 (${timeout}초)"
        $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" logs sequencer-l2
        exit 1
    fi
    sleep 2
    elapsed=$((elapsed + 2))
    echo -n "."
done
echo ""

log_success "Sequencer L2 op-geth 시작 완료"

# Remove Sequencer op-node safedb to ensure clean state
if docker volume inspect scripts_sequencer_safedb_data >/dev/null 2>&1; then
    log_info "Removing Sequencer op-node safedb for clean state..."
    docker volume rm scripts_sequencer_safedb_data 2>/dev/null || true
fi

# Sequencer op-node 시작
log_info "Sequencer op-node를 시작합니다..."
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d sequencer-op-node

# op-batcher 시작
log_info "op-batcher를 시작합니다..."
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d op-batcher

# op-proposer 시작
log_info "op-proposer를 시작합니다..."
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d op-proposer

log_success "Sequencer 스택 배포 완료"

##############################################################################
# 8. Challenger 스택 배포 (독립적!) ⭐
##############################################################################

log_info "=========================================="
log_info "독립적인 Challenger 스택을 배포합니다 ⭐"
log_info "=========================================="

# CRITICAL: Remove Challenger L2 volume to ensure it initializes with the correct genesis
if docker volume inspect scripts_challenger_l2_data >/dev/null 2>&1; then
    log_warn "Removing existing Challenger L2 volume to ensure genesis consistency..."
    $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" stop challenger-l2 2>/dev/null || true
    $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" rm -f challenger-l2 2>/dev/null || true
    docker volume rm scripts_challenger_l2_data 2>/dev/null || true
    log_success "Challenger L2 volume removed"
fi

# Challenger L2 geth 시작 (별도 데이터베이스!)
log_info "Challenger L2 op-geth를 시작합니다 (별도 DB!)..."
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d challenger-l2

log_info "Challenger L2 RPC가 준비될 때까지 대기합니다..."
timeout=120
elapsed=0
while ! curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "http://localhost:9546" > /dev/null 2>&1; do
    if [ $elapsed -ge $timeout ]; then
        log_error "Challenger L2 RPC 시작 시간 초과 (${timeout}초)"
        $DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" logs challenger-l2
        exit 1
    fi
    sleep 2
    elapsed=$((elapsed + 2))
    echo -n "."
done
echo ""

log_success "Challenger L2 op-geth 시작 완료"

# Remove Challenger op-node safedb to ensure clean state
if docker volume inspect scripts_challenger_safedb_data >/dev/null 2>&1; then
    log_info "Removing Challenger op-node safedb for clean state..."
    docker volume rm scripts_challenger_safedb_data 2>/dev/null || true
fi

# Challenger op-node 시작 (Follower 모드)
log_info "Challenger op-node를 시작합니다 (Follower 모드)..."
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d challenger-op-node

# op-challenger 시작
log_info "op-challenger를 시작합니다..."
$DOCKER_COMPOSE -f "$DOCKER_COMPOSE_FILE" up -d op-challenger

log_success "독립적인 Challenger 스택 배포 완료! ✅"

##############################################################################
# 9. Indexer 배포 (선택사항)
##############################################################################

if [ "$WITH_INDEXER" = true ]; then
    log_info "Indexer 스택을 배포합니다..."

    if [ -d "${PROJECT_ROOT}/indexer" ]; then
        cd "${PROJECT_ROOT}/indexer"
        $DOCKER_COMPOSE up -d
        log_success "Indexer 스택 배포 완료"
    else
        log_warn "Indexer 디렉토리를 찾을 수 없습니다: ${PROJECT_ROOT}/indexer"
    fi
fi

##############################################################################
# 10. 헬스 체크
##############################################################################

log_info "서비스 헬스 체크를 수행합니다..."
sleep 5  # 서비스가 완전히 시작될 때까지 대기

# 헬스 체크 스크립트 실행
if [ -f "${SCRIPT_DIR}/health-check.sh" ]; then
    bash "${SCRIPT_DIR}/health-check.sh"
else
    log_warn "헬스 체크 스크립트를 찾을 수 없습니다: ${SCRIPT_DIR}/health-check.sh"
fi

##############################################################################
# 11. 배포 완료
##############################################################################

log_success "=========================================="
log_success "L2 시스템 배포가 완료되었습니다!"
log_success "=========================================="

echo ""
log_info "📊 서비스 엔드포인트:"
echo ""
echo "  ┌─ L1 Ethereum"
echo "  │  └─ RPC:                    http://localhost:8545"
echo "  │"
echo "  ┌─ Sequencer 스택"
echo "  │  ├─ L2 RPC (사용자):        http://localhost:9545"
echo "  │  ├─ op-node RPC:            http://localhost:7545"
echo "  │  ├─ op-batcher RPC:         http://localhost:6545"
echo "  │  └─ op-proposer RPC:        http://localhost:6546"
echo "  │"
echo "  └─ Challenger 스택 ⭐ (독립적!)"
echo "     ├─ Challenger L2 RPC:      http://localhost:9546"
echo "     ├─ Challenger op-node RPC: http://localhost:7546"
echo "     └─ op-challenger:          (백그라운드)"
echo ""

log_info "📝 다음 단계:"
echo "  1. 로그 확인:"
echo "     ${DOCKER_COMPOSE} -f ${SCRIPT_DIR}/docker-compose-full.yml logs -f"
echo ""
echo "  2. 특정 서비스 로그:"
echo "     ${DOCKER_COMPOSE} -f ${SCRIPT_DIR}/docker-compose-full.yml logs -f sequencer-op-node"
echo "     ${DOCKER_COMPOSE} -f ${SCRIPT_DIR}/docker-compose-full.yml logs -f challenger-op-node"
echo "     ${DOCKER_COMPOSE} -f ${SCRIPT_DIR}/docker-compose-full.yml logs -f op-challenger"
echo ""
echo "  3. 헬스 체크:"
echo "     ${SCRIPT_DIR}/health-check.sh"
echo ""
echo "  4. 상태 모니터링:"
echo "     ${DOCKER_COMPOSE} -f ${SCRIPT_DIR}/docker-compose-full.yml ps"
echo ""

if [ "$WITH_INDEXER" = true ]; then
    echo "  Indexer API:               http://localhost:8100"
    echo "  Grafana:                   http://localhost:3000"
    echo ""
fi

log_info "🛑 시스템 중지:"
echo "     ${DOCKER_COMPOSE} -f ${SCRIPT_DIR}/docker-compose-full.yml down"
echo ""
log_info "🗑️  데이터 삭제:"
echo "     ${DOCKER_COMPOSE} -f ${SCRIPT_DIR}/docker-compose-full.yml down -v"
echo ""

log_warn "⭐ 중요: Challenger 스택은 Sequencer와 완전히 독립적으로 작동합니다!"
log_info "   - challenger-l2: 별도 데이터베이스로 L1에서 독립적으로 L2 재구성"
log_info "   - challenger-op-node: Follower 모드로 L1 데이터 검증"
log_info "   - op-challenger: 독립적으로 검증된 데이터로 챌린지 수행"
echo ""

log_info "=========================================="
log_info "📝 로그 파일 생성 완료"
log_info "=========================================="
log_info "배포 로그: $DEPLOY_LOG"
log_info "설정 추적 로그: $CONFIG_TRACE_LOG"
echo ""
log_info "설정값 적용 과정을 확인하려면:"
echo "    cat $CONFIG_TRACE_LOG"
echo ""

log_config "=== 스크립트 종료 ==="
log_config "Deployment completed successfully"

exit 0
