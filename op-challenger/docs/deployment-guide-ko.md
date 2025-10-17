# op-challenger 전체 시스템 구동 가이드

## 목차
1. [시스템 개요](#시스템-개요)
2. [서비스 구성](#서비스-구성)
3. [사전 요구사항](#사전-요구사항)
4. [개발 환경 설정](#개발-환경-설정)
5. [Docker Compose를 통한 시스템 구동](#docker-compose를-통한-시스템-구동)
6. [각 서비스별 상세 설정](#각-서비스별-상세-설정)
7. [시스템 모니터링 및 디버깅](#시스템-모니터링-및-디버깅)
8. [프로덕션 배포](#프로덕션-배포)

---

## 시스템 개요

Tokamak Thanos (Optimism Fork) 스택은 다음과 같은 핵심 컴포넌트들로 구성됩니다:

```
┌─────────────────────────────────────────────────────────────┐
│                    Thanos Stack Architecture                │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │    L1    │  │    L2    │  │ op-node  │  │op-batcher│   │
│  │  (geth)  │  │  (geth)  │  │(rollup)  │  │(sequencer│   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  │  batcher)│   │
│       │             │             │         └────┬─────┘   │
│       └─────────────┴─────────────┴──────────────┘         │
│                         │                                   │
│       ┌─────────────────┼───────────────────┐              │
│       │                 │                   │              │
│  ┌────▼─────┐    ┌──────▼──────┐    ┌──────▼──────┐       │
│  │op-proposer│    │op-challenger│    │  da-server  │       │
│  │(output   │    │(fault proof)│    │  (plasma)   │       │
│  │proposer) │    └─────────────┘    └─────────────┘       │
│  └──────────┘                                              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 서비스 구성

### 핵심 서비스

#### 1. **L1 (Ethereum Layer 1)**
- **이미지**: 커스텀 geth 이미지 (`Dockerfile.l1`)
- **포트**:
  - `8545`: HTTP JSON-RPC
  - `8546`: WebSocket JSON-RPC
  - `7060`: pprof (프로파일링)
  - `9999`: 추가 RPC
- **역할**: L1 블록체인 시뮬레이션 (로컬 개발용)
- **데이터**: `l1_data` 볼륨에 블록체인 데이터 저장

**환경 변수**:
```bash
GETH_MINER_RECOMMIT=100ms  # 빠른 블록 생성 (개발용)
BLOCK_NUMBER=${BLOCK_NUMBER}
L1_RPC=${L1_RPC}
```

**설정 파일**:
- `/.devnet/genesis-l1.json`: L1 제네시스 설정
- `/config/test-jwt-secret.txt`: JWT 인증 시크릿

---

#### 2. **L2 (Optimism Layer 2)**
- **이미지**: 커스텀 L2 geth 이미지 (`Dockerfile.l2`)
- **포트**:
  - `9545`: HTTP JSON-RPC
  - `8060`: pprof
  - `8551`: Engine API (op-node 연결용)
- **역할**: L2 실행 레이어 (EVM 실행)
- **데이터**: `l2_data` 볼륨

**환경 변수**:
```bash
GETH_MINER_RECOMMIT=100ms
```

**설정 파일**:
- `/.devnet/genesis-l2.json`: L2 제네시스 설정
- `/config/test-jwt-secret.txt`: Engine API JWT 인증

**특이사항**:
- `--authrpc.jwtsecret` 플래그로 Engine API 보안 연결

---

#### 3. **op-node (Rollup Node)**
- **이미지**: `tokamaknetwork/thanos-op-node:latest`
- **빌드 파일**: `ops/docker/op-stack-go/Dockerfile` (target: `op-node-target`)
- **포트**:
  - `7545`: Rollup RPC (외부 접속)
  - `9003`: P2P (TCP/UDP)
  - `7300`: Prometheus 메트릭
  - `6060`: pprof
- **역할**: Rollup 합의 레이어, L1↔L2 동기화

**주요 설정**:
```bash
--l1=http://l1:8545                    # L1 연결
--l2=http://l2:8551                    # L2 Engine API 연결
--l2.jwt-secret=/config/test-jwt-secret.txt
--rollup.config=/rollup.json           # Rollup 설정
--sequencer.enabled                    # 시퀀서 모드
--sequencer.l1-confs=0                 # L1 확인 블록 수 (개발용 0)
--p2p.sequencer.key=...                # P2P 키
--snapshotlog.file=/op_log/snapshot.log
--l1.beacon=${L1_BEACON}               # L1 Beacon API
--plasma.enabled=${PLASMA_ENABLED}     # Plasma 모드
--plasma.da-server=http://da-server:3100
```

**의존성**:
- L1, L2가 먼저 시작되어야 함
- 시작 시 `nc -z l1 8545`로 L1 대기

**볼륨**:
- `safedb_data:/db`: Safe head 데이터베이스
- `op_log:/op_log`: Snapshot 로그

---

#### 4. **op-proposer (Output Proposer)**
- **이미지**: `tokamaknetwork/thanos-op-proposer:latest`
- **빌드 파일**: `ops/docker/op-stack-go/Dockerfile` (target: `op-proposer-target`)
- **포트**:
  - `6062`: pprof
  - `7302`: 메트릭
  - `6546`: RPC
- **역할**: L2 상태 루트를 L1에 제안

**환경 변수**:
```bash
OP_PROPOSER_L1_ETH_RPC=http://l1:8545
OP_PROPOSER_ROLLUP_RPC=http://op-node:8545
OP_PROPOSER_POLL_INTERVAL=1s           # 폴링 간격
OP_PROPOSER_NUM_CONFIRMATIONS=1        # L1 확인 블록 수
OP_PROPOSER_MNEMONIC=test test test... # 테스트 니모닉
OP_PROPOSER_L2_OUTPUT_HD_PATH=m/44'/60'/0'/0/1
OP_PROPOSER_L2OO_ADDRESS=${L2OO_ADDRESS}           # L2OutputOracle 주소
OP_PROPOSER_GAME_FACTORY_ADDRESS=${DGF_ADDRESS}   # DisputeGameFactory 주소
OP_PROPOSER_GAME_TYPE=${DG_TYPE}                  # 게임 타입
OP_PROPOSER_PROPOSAL_INTERVAL=${PROPOSAL_INTERVAL}
OP_PROPOSER_ALLOW_NON_FINALIZED=true
```

**의존성**:
- L1, L2, op-node
- op-node RPC가 준비될 때까지 대기

---

#### 5. **op-batcher (Batch Submitter)**
- **이미지**: `tokamaknetwork/thanos-op-batcher:latest`
- **빌드 파일**: `ops/docker/op-stack-go/Dockerfile` (target: `op-batcher-target`)
- **포트**:
  - `6061`: pprof
  - `7301`: 메트릭
  - `6545`: RPC
- **역할**: L2 트랜잭션을 배치로 묶어 L1에 제출

**환경 변수**:
```bash
OP_BATCHER_L1_ETH_RPC=http://l1:8545
OP_BATCHER_L2_ETH_RPC=http://l2:8545
OP_BATCHER_ROLLUP_RPC=http://op-node:8545
OP_BATCHER_MAX_CHANNEL_DURATION=1      # 채널 최대 기간
OP_BATCHER_SUB_SAFETY_MARGIN=4         # 안전 마진
OP_BATCHER_POLL_INTERVAL=1s
OP_BATCHER_NUM_CONFIRMATIONS=1
OP_BATCHER_MNEMONIC=test test test...
OP_BATCHER_SEQUENCER_HD_PATH=m/44'/60'/0'/0/2
OP_BATCHER_BATCH_TYPE=0                # 배치 타입 (0=legacy, 1=zlib)
OP_BATCHER_PLASMA_ENABLED=${PLASMA_ENABLED}
OP_BATCHER_PLASMA_DA_SERVER=http://da-server:3100
```

---

#### 6. **op-challenger (Fault Proof Challenger)** ⭐
- **이미지**: `tokamaknetwork/thanos-op-challenger:latest`
- **빌드 파일**: `ops/docker/op-stack-go/Dockerfile` (target: `op-challenger-target`)
- **역할**: 무효한 output root에 이의 제기

**환경 변수**:
```bash
OP_CHALLENGER_L1_ETH_RPC=http://l1:8545
OP_CHALLENGER_L1_BEACON=unset          # 개발환경에선 미사용
OP_CHALLENGER_ROLLUP_RPC=http://op-node:8545
OP_CHALLENGER_L2_ETH_RPC=http://l2:8545
OP_CHALLENGER_TRACE_TYPE=cannon,fast   # Trace provider 타입
OP_CHALLENGER_GAME_FACTORY_ADDRESS=${DGF_ADDRESS}
OP_CHALLENGER_UNSAFE_ALLOW_INVALID_PRESTATE=true  # 개발용
OP_CHALLENGER_DATADIR=/db
OP_CHALLENGER_MNEMONIC=test test test...
OP_CHALLENGER_HD_PATH=m/44'/60'/0'/0/4

# Cannon 설정
OP_CHALLENGER_CANNON_ROLLUP_CONFIG=./.devnet/rollup.json
OP_CHALLENGER_CANNON_L2_GENESIS=./.devnet/genesis-l2.json
OP_CHALLENGER_CANNON_BIN=./cannon/bin/cannon
OP_CHALLENGER_CANNON_SERVER=/op-program/op-program
OP_CHALLENGER_CANNON_PRESTATE=/op-program/prestate.json
OP_CHALLENGER_NUM_CONFIRMATIONS=1
```

**볼륨**:
- `challenger_data:/db`: 게임 데이터 및 증명
- `../op-program/bin:/op-program`: op-program 바이너리 마운트

**의존성**:
- L1, L2, op-node
- op-program 바이너리 및 prestate 필요

---

#### 7. **da-server (Data Availability Server)**
- **이미지**: `tokamaknetwork/thanos-da-server:latest`
- **빌드 파일**: `ops/docker/op-stack-go/Dockerfile` (target: `da-server-target`)
- **포트**: `3100`
- **역할**: Plasma 모드에서 데이터 가용성 제공

**설정**:
```bash
da-server \
  --file.path=/data \
  --addr=0.0.0.0 \
  --port=3100 \
  --log.level=debug \
  --generic-commitment=${PLASMA_GENERIC_DA}
```

**볼륨**:
- `da_data:/data`: DA 데이터 저장

---

#### 8. **artifact-server (Static File Server)**
- **이미지**: `nginx:1.25-alpine`
- **포트**: `8080`
- **역할**: 제네시스 및 설정 파일 제공
- **볼륨**: `${PWD}/../.devnet/:/usr/share/nginx/html/:ro`

**제공 파일**:
- `addresses.json`: 컨트랙트 주소 목록
- `rollup.json`: Rollup 설정
- `genesis-l1.json`, `genesis-l2.json`: 제네시스 파일
- `allocs-l1.json`, `allocs-l2.json`: 계정 할당

---

## 사전 요구사항

### 필수 도구

```bash
# 1. Git
git --version  # >= 2.x

# 2. Docker & Docker Compose
docker --version          # >= 20.10
docker compose version    # >= 2.x

# 3. Go
go version  # >= 1.21

# 4. Node.js & pnpm
node --version  # >= 18.x
pnpm --version  # >= 8.x

# 5. Foundry (Solidity 개발)
forge --version
cast --version

# 6. Python
python3 --version  # >= 3.9

# 7. jq (JSON 파싱)
jq --version

# 8. geth (go-ethereum)
geth version
```

### 설치 스크립트 (macOS/Linux)

```bash
# Homebrew (macOS)
brew install git docker docker-compose go node pnpm jq python@3.9

# Foundry 설치
curl -L https://foundry.paradigm.xyz | bash
foundryup

# geth 설치 (프로젝트 Makefile 사용)
cd /path/to/tokamak-thanos
make install-geth
```

---

## 개발 환경 설정

### 1. 저장소 클론

```bash
git clone https://github.com/tokamak-network/tokamak-thanos.git
cd tokamak-thanos
```

### 2. 서브모듈 초기화

```bash
make submodules
# 또는
git submodule update --init --recursive
```

### 3. Node.js 버전 설정

```bash
nvm use
# .nvmrc 파일에 명시된 버전 사용
```

### 4. TypeScript 패키지 빌드

```bash
pnpm clean
pnpm install
pnpm build
```

**주의사항**:
- 브랜치 변경 시 반드시 재빌드 필요
- 패키지 간 호환성 문제 방지

### 5. Go 바이너리 빌드

```bash
# 모든 Go 서비스 빌드
make build-go

# 개별 빌드
make op-node
make op-proposer
make op-batcher
make op-challenger
make op-program
make cannon
```

### 6. Cannon Prestate 생성

op-challenger는 Cannon prestate가 필요합니다:

```bash
make cannon-prestate
```

**생성 파일**:
- `op-program/bin/prestate.json`: 초기 VM 상태
- `op-program/bin/prestate-proof.json`: 증명 데이터
- `op-program/bin/meta.json`: 메타데이터

---

## Docker Compose를 통한 시스템 구동

### 방법 1: Makefile 사용 (권장)

#### 전체 시스템 시작

```bash
# 1단계: Devnet 초기화 및 시작
make devnet-up
```

**내부 동작**:
```bash
# 1. 사전 요구사항 확인 및 준비
make pre-devnet
  ├─ make submodules
  ├─ geth 설치 확인
  └─ cannon-prestate 생성 확인

# 2. Devnet allocs 생성 (필요시)
make devnet-allocs

# 3. Python 스크립트로 devnet 시작
PYTHONPATH=./bedrock-devnet python3 ./bedrock-devnet/main.py --monorepo-dir=.
```

**Python 스크립트 동작**:
1. L1/L2 제네시스 생성
2. 스마트 컨트랙트 배포
3. 설정 파일 생성 (`.devnet/`)
4. Docker Compose 시작

#### 시스템 중지

```bash
make devnet-down
```

#### 완전 초기화 (데이터 삭제)

```bash
make devnet-clean
```

**삭제 대상**:
- `.devnet/` 디렉토리
- Docker 볼륨 (l1_data, l2_data 등)
- Docker 이미지 (ops-bedrock*)
- 배포 데이터

---

### 방법 2: Docker Compose 직접 사용

#### 전제 조건

1. `.devnet/` 디렉토리 생성 및 설정 파일 준비
2. 환경 변수 설정

```bash
# 환경 변수 파일 생성
cat > .env <<EOF
IMAGE_TAG=latest
L2_IMAGE=tokamaknetwork/thanos-geth:latest
PLASMA_ENABLED=false
PLASMA_DA_SERVICE=false
PLASMA_GENERIC_DA=false
L1_FORK_PUBLIC_NETWORK=false
WAITING_L1_PORT=8545

# Contract 주소 (devnet-up 후 자동 생성됨)
DGF_ADDRESS=0x...
L2OO_ADDRESS=0x...
DG_TYPE=0
PROPOSAL_INTERVAL=12s
EOF
```

#### Docker 이미지 빌드

```bash
cd ops-bedrock
export COMPOSE_DOCKER_CLI_BUILD=1
export DOCKER_BUILDKIT=1
docker compose build
```

**빌드 대상**:
- `l1`: L1 geth 이미지
- `l2`: L2 geth 이미지
- `op-node`, `op-proposer`, `op-batcher`, `op-challenger`: Go 서비스들
- `da-server`: DA 서버

#### 서비스 시작

```bash
cd ops-bedrock
docker compose up -d
```

#### 특정 서비스만 시작

```bash
# L1, L2만 시작
docker compose up -d l1 l2

# op-challenger만 재시작
docker compose restart op-challenger
```

---

### 시작 순서 및 의존성

```
1. L1 시작 (독립적)
   └─> 블록 생성 시작

2. L2 시작 (독립적)
   └─> Engine API 대기

3. op-node 시작
   ├─ L1 연결 대기 (nc -z l1 8545)
   ├─ L2 연결 (Engine API)
   └─> Rollup 동기화 시작

4. op-proposer 시작
   ├─ op-node 대기 (nc -z op-node 8545)
   └─> Output 제안 시작

5. op-batcher 시작
   ├─ op-node 대기
   └─> 배치 제출 시작

6. op-challenger 시작
   ├─ L1, L2, op-node 의존
   └─> 게임 모니터링 시작

7. da-server 시작 (plasma 모드)

8. artifact-server 시작 (독립적)
```

**대기 로직 예시**:
```bash
while ! nc -z op-node 8545; do
  sleep 1
done
echo 'op-node is now available'
```

---

## 각 서비스별 상세 설정

### L1 설정 (`Dockerfile.l1`, `entrypoint-l1.sh`)

**Dockerfile.l1**:
```dockerfile
FROM ethereum/client-go:latest
COPY entrypoint-l1.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
```

**entrypoint-l1.sh 주요 로직**:
```bash
#!/bin/sh
set -e

# JWT 시크릿 생성 (없는 경우)
if [ ! -f /config/test-jwt-secret.txt ]; then
  openssl rand -hex 32 > /config/test-jwt-secret.txt
fi

# Geth 초기화
if [ ! -d /db/geth ]; then
  geth init --datadir=/db /genesis.json
fi

# Geth 시작
exec geth \
  --datadir=/db \
  --http \
  --http.addr=0.0.0.0 \
  --http.port=8545 \
  --http.api=eth,net,web3,debug,txpool \
  --ws \
  --ws.addr=0.0.0.0 \
  --ws.port=8546 \
  --ws.api=eth,net,web3,debug,txpool \
  --authrpc.addr=0.0.0.0 \
  --authrpc.port=8551 \
  --authrpc.jwtsecret=/config/test-jwt-secret.txt \
  --nodiscover \
  --maxpeers=0 \
  --mine \
  --miner.etherbase=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 \
  --allow-insecure-unlock \
  --rpc.allow-unprotected-txs
```

---

### op-challenger 상세 설정

#### Dockerfile 타겟 (`ops/docker/op-stack-go/Dockerfile`)

```dockerfile
# op-challenger-target 스테이지
FROM op-stack-go-base AS op-challenger-target

WORKDIR /app

# 바이너리 복사
COPY --from=op-challenger-builder /app/op-challenger/bin/op-challenger /usr/local/bin/

# cannon, op-program 마운트 경로
VOLUME /op-program

ENTRYPOINT ["op-challenger"]
```

#### 컨테이너 내부 구조

```
/
├── usr/local/bin/op-challenger   # 실행 파일
├── db/                            # 데이터 디렉토리
│   ├── {game-addr-1}/
│   │   ├── proofs/
│   │   └── preimages/
│   └── {game-addr-2}/
└── op-program/                    # 마운트
    ├── op-program                 # op-program 바이너리
    └── prestate.json              # 초기 상태
```

#### 환경 변수 매핑

Docker Compose 환경 변수는 op-challenger 플래그로 변환됩니다:

| 환경 변수 | CLI 플래그 | 설명 |
|----------|-----------|------|
| `OP_CHALLENGER_L1_ETH_RPC` | `--l1-eth-rpc` | L1 RPC URL |
| `OP_CHALLENGER_L2_ETH_RPC` | `--l2-eth-rpc` | L2 RPC URL |
| `OP_CHALLENGER_ROLLUP_RPC` | `--rollup-rpc` | Rollup RPC URL |
| `OP_CHALLENGER_TRACE_TYPE` | `--trace-type` | cannon,fast,asterisc |
| `OP_CHALLENGER_GAME_FACTORY_ADDRESS` | `--game-factory-address` | Factory 주소 |
| `OP_CHALLENGER_DATADIR` | `--datadir` | 데이터 디렉토리 |
| `OP_CHALLENGER_CANNON_BIN` | `--cannon-bin` | cannon 경로 |
| `OP_CHALLENGER_CANNON_SERVER` | `--cannon-server` | op-program 경로 |
| `OP_CHALLENGER_CANNON_PRESTATE` | `--cannon-prestate` | prestate.json 경로 |
| `OP_CHALLENGER_MNEMONIC` | `--mnemonic` | 니모닉 |
| `OP_CHALLENGER_HD_PATH` | `--hd-path` | HD 지갑 경로 |

#### 런타임 로그 예시

```
INFO [01-17|12:00:00.123] Starting op-challenger    version=v1.0.0
INFO [01-17|12:00:00.456] Loaded config
  l1-eth-rpc=http://l1:8545
  rollup-rpc=http://op-node:8545
  trace-types=[cannon fast]
  game-factory=0x1234...
INFO [01-17|12:00:01.789] Connected to L1          chain-id=900
INFO [01-17|12:00:01.890] Connected to Rollup Node
INFO [01-17|12:00:02.000] Starting scheduler       max-concurrency=8
INFO [01-17|12:00:02.100] Starting monitoring
INFO [01-17|12:00:03.000] Subscribed to L1 heads
```

---

### 설정 파일 구조

#### `.devnet/` 디렉토리

```
.devnet/
├── addresses.json          # 배포된 컨트랙트 주소
├── rollup.json             # Rollup 설정
├── genesis-l1.json         # L1 제네시스
├── genesis-l2.json         # L2 제네시스
├── allocs-l1.json          # L1 계정 할당
└── allocs-l2*.json         # L2 계정 할당
```

#### `addresses.json` 구조

```json
{
  "AddressManager": "0x...",
  "L1CrossDomainMessengerProxy": "0x...",
  "L1StandardBridgeProxy": "0x...",
  "L2OutputOracleProxy": "0x...",
  "DisputeGameFactoryProxy": "0x...",
  "OptimismPortalProxy": "0x...",
  "SystemConfigProxy": "0x...",
  "ProxyAdmin": "0x...",
  ...
}
```

#### `rollup.json` 주요 필드

```json
{
  "genesis": {
    "l1": {
      "hash": "0x...",
      "number": 0
    },
    "l2": {
      "hash": "0x...",
      "number": 0
    },
    "l2_time": 1234567890,
    "system_config": {
      "batcherAddr": "0x...",
      "overhead": "0x...",
      "scalar": "0x...",
      "gasLimit": 30000000
    }
  },
  "block_time": 2,
  "max_sequencer_drift": 600,
  "seq_window_size": 3600,
  "channel_timeout": 300,
  "l1_chain_id": 900,
  "l2_chain_id": 901,
  "regolith_time": 0,
  "canyon_time": 0,
  "delta_time": 0,
  "ecotone_time": 0,
  "batch_inbox_address": "0xff00000000000000000000000000000000000901",
  "deposit_contract_address": "0x...",
  "l1_system_config_address": "0x..."
}
```

---

## 시스템 모니터링 및 디버깅

### 로그 확인

#### 전체 로그

```bash
cd ops-bedrock
docker compose logs -f
```

#### 서비스별 로그

```bash
# op-challenger 로그
docker compose logs -f op-challenger

# op-node 로그
docker compose logs -f op-node

# L1 로그
docker compose logs -f l1

# 여러 서비스 동시
docker compose logs -f op-challenger op-node op-proposer
```

#### 최근 N줄만 보기

```bash
docker compose logs --tail=100 op-challenger
```

### 컨테이너 상태 확인

```bash
# 실행 중인 컨테이너
docker compose ps

# 상세 정보
docker compose ps -a

# 리소스 사용량
docker stats
```

### 컨테이너 내부 접속

```bash
# op-challenger 쉘
docker compose exec op-challenger sh

# 파일 시스템 확인
docker compose exec op-challenger ls -la /db

# 로그 파일 직접 확인
docker compose exec op-challenger cat /op_log/snapshot.log
```

### 메트릭 확인

각 서비스는 Prometheus 메트릭을 노출합니다:

```bash
# op-node 메트릭
curl http://localhost:7300/metrics

# op-proposer 메트릭
curl http://localhost:7302/metrics

# op-batcher 메트릭
curl http://localhost:7301/metrics
```

**주요 메트릭**:
- `op_node_sync_status`: 동기화 상태
- `op_challenger_games_in_progress`: 진행 중 게임 수
- `op_challenger_games_challenger_won`: 챌린저 승리 게임 수
- `op_proposer_proposals_submitted`: 제출된 제안 수
- `op_batcher_batches_submitted`: 제출된 배치 수

### Pprof 프로파일링

```bash
# CPU 프로파일
curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof

# 메모리 프로파일
curl http://localhost:6060/debug/pprof/heap > heap.prof

# 고루틴
curl http://localhost:6060/debug/pprof/goroutine?debug=2

# 프로파일 분석 (Go 필요)
go tool pprof -http=:8081 cpu.prof
```

### 네트워크 디버깅

```bash
# 컨테이너 네트워크 확인
docker network ls
docker network inspect ops-bedrock_default

# 컨테이너 간 연결 테스트
docker compose exec op-challenger nc -zv l1 8545
docker compose exec op-node nc -zv l2 8551
```

### 일반적인 문제 해결

#### 1. op-challenger가 시작되지 않음

**증상**: 컨테이너가 즉시 종료

**확인**:
```bash
docker compose logs op-challenger
```

**가능한 원인**:
- `op-program/bin/` 디렉토리가 비어있음
  ```bash
  make cannon-prestate
  ```
- `.devnet/addresses.json`이 없음
  ```bash
  make devnet-allocs
  ```
- 환경 변수 누락
  ```bash
  docker compose config | grep OP_CHALLENGER
  ```

#### 2. L1/L2 연결 실패

**증상**: "connection refused" 에러

**확인**:
```bash
# L1 포트 확인
curl http://localhost:8545 -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# L2 포트 확인
curl http://localhost:9545 -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

**해결**:
```bash
# 서비스 재시작
docker compose restart l1 l2

# 포트 충돌 확인
lsof -i :8545
lsof -i :9545
```

#### 3. 디스크 공간 부족

**확인**:
```bash
docker system df
docker volume ls
```

**정리**:
```bash
# 사용하지 않는 볼륨 제거
docker volume prune

# 전체 정리
make devnet-clean
docker system prune -a
```

#### 4. 게임이 생성되지 않음

**op-proposer 로그 확인**:
```bash
docker compose logs op-proposer | grep -i "proposal"
```

**Factory 주소 확인**:
```bash
cat .devnet/addresses.json | jq .DisputeGameFactoryProxy
```

**수동 게임 생성**:
```bash
docker compose exec op-challenger op-challenger create-game \
  --l1-eth-rpc http://l1:8545 \
  --game-factory-address $DGF_ADDRESS \
  --output-root 0x... \
  --l2-block-num 1000 \
  --mnemonic "test test test test test test test test test test test junk" \
  --hd-path "m/44'/60'/0'/0/8"
```

---

## 프로덕션 배포

### 환경 변수 설정

프로덕션 환경에서는 보안을 위해 다음을 변경해야 합니다:

```bash
# .env.production
IMAGE_TAG=v1.0.0  # 특정 버전 태그

# 실제 RPC 엔드포인트
L1_ETH_RPC=https://mainnet.infura.io/v3/YOUR_KEY
L1_BEACON=https://beacon-nd-123-456-789.p2pify.com/YOUR_KEY

# Plasma 설정 (필요시)
PLASMA_ENABLED=true
PLASMA_DA_SERVICE=true

# 보안 설정
OP_CHALLENGER_UNSAFE_ALLOW_INVALID_PRESTATE=false
OP_CHALLENGER_MNEMONIC=  # 사용하지 않음
OP_CHALLENGER_PRIVATE_KEY=${PRIVATE_KEY}  # 환경 변수로 주입

# 성능 튜닝
OP_CHALLENGER_MAX_CONCURRENCY=4
OP_CHALLENGER_MAX_PENDING_TX=10
OP_CHALLENGER_POLL_INTERVAL=12s

# 모니터링
OP_CHALLENGER_METRICS_ENABLED=true
OP_CHALLENGER_METRICS_ADDR=0.0.0.0
OP_CHALLENGER_METRICS_PORT=7300
```

### Docker Compose 오버라이드

`docker-compose.override.yml`:

```yaml
version: '3.4'

services:
  op-challenger:
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "10"
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          cpus: '1'
          memory: 2G
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "-", "http://localhost:7300/metrics"]
      interval: 30s
      timeout: 10s
      retries: 3
    environment:
      # 프로덕션 설정
      OP_CHALLENGER_PRIVATE_KEY: ${OP_CHALLENGER_PRIVATE_KEY}
      OP_CHALLENGER_UNSAFE_ALLOW_INVALID_PRESTATE: "false"
```

### 시작 스크립트

`start-production.sh`:

```bash
#!/bin/bash
set -e

# 환경 변수 로드
source .env.production

# 필수 변수 확인
if [ -z "$OP_CHALLENGER_PRIVATE_KEY" ]; then
  echo "Error: OP_CHALLENGER_PRIVATE_KEY not set"
  exit 1
fi

# 이미지 최신화
docker compose pull

# 서비스 시작
docker compose -f docker-compose.yml -f docker-compose.override.yml up -d

# 헬스 체크
sleep 10
docker compose ps

echo "System started. Checking logs..."
docker compose logs --tail=50 op-challenger
```

### 모니터링 스택 구성

Prometheus + Grafana를 추가하여 모니터링:

`docker-compose.monitoring.yml`:

```yaml
version: '3.4'

services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
    environment:
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_PASSWORD}

volumes:
  prometheus_data:
  grafana_data:
```

`prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'op-challenger'
    static_configs:
      - targets: ['op-challenger:7300']

  - job_name: 'op-node'
    static_configs:
      - targets: ['op-node:7300']

  - job_name: 'op-proposer'
    static_configs:
      - targets: ['op-proposer:7302']

  - job_name: 'op-batcher'
    static_configs:
      - targets: ['op-batcher:7301']
```

### 백업 및 복구

#### 볼륨 백업

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR=./backups/$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

# Challenger 데이터 백업
docker run --rm \
  -v ops-bedrock_challenger_data:/data \
  -v $PWD/$BACKUP_DIR:/backup \
  alpine tar czf /backup/challenger_data.tar.gz -C /data .

# .devnet 백업
tar czf $BACKUP_DIR/devnet_config.tar.gz .devnet/

echo "Backup completed: $BACKUP_DIR"
```

#### 복구

```bash
#!/bin/bash
# restore.sh

BACKUP_DIR=$1

if [ -z "$BACKUP_DIR" ]; then
  echo "Usage: ./restore.sh <backup_directory>"
  exit 1
fi

# 시스템 중지
docker compose down

# Challenger 데이터 복구
docker run --rm \
  -v ops-bedrock_challenger_data:/data \
  -v $PWD/$BACKUP_DIR:/backup \
  alpine tar xzf /backup/challenger_data.tar.gz -C /data

# .devnet 복구
tar xzf $BACKUP_DIR/devnet_config.tar.gz

# 시스템 재시작
docker compose up -d

echo "Restore completed"
```

---

## 유용한 명령어 모음

### 개발 워크플로우

```bash
# 전체 빌드 및 시작
make build && make devnet-up

# 코드 변경 후 재빌드 (op-challenger)
make op-challenger
docker compose build op-challenger
docker compose up -d op-challenger

# 로그 실시간 모니터링
docker compose logs -f op-challenger | grep -i "game"

# 특정 게임 클레임 조회
docker compose exec op-challenger op-challenger list-claims \
  --l1-eth-rpc http://l1:8545 \
  --game-address 0x...
```

### 테스트

```bash
# Go 단위 테스트
cd op-challenger
go test ./... -v

# 통합 테스트
make devnet-test

# TypeScript 테스트
pnpm test

# e2e 테스트
cd op-e2e
go test ./... -v
```

### 정리 명령어

```bash
# Devnet만 정리
make devnet-clean

# 전체 정리 (빌드 아티팩트 포함)
make nuke

# Docker 완전 정리
docker compose down -v
docker system prune -a --volumes
```

---

## 참고 자료

### 공식 문서
- [Optimism Bedrock Specs](https://specs.optimism.io)
- [Fault Proof Specs](https://specs.optimism.io/experimental/fault-proof/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

### 관련 파일
- `Makefile`: 빌드 및 devnet 관리 스크립트
- `ops-bedrock/docker-compose.yml`: 서비스 정의
- `bedrock-devnet/`: Devnet 초기화 스크립트
- `CONTRIBUTING.md`: 개발 가이드

### 추가 도구
- **op-dispute-mon**: 게임 모니터링 대시보드
- **cannon**: MIPS VM 시뮬레이터
- **op-program**: Fault proof 프로그램

---

**대상 프로젝트**: tokamak-thanos (Optimism Fork)
