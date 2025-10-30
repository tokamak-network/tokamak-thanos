# L2 시스템 배포 가이드

## 목차
1. [개요](#개요)
2. [완전한 시스템 아키텍처](#완전한-시스템-아키텍처)
3. [실제 배포 구성](#실제-배포-구성)
4. [리소스 요구사항](#리소스-요구사항)
5. [환경별 L1 설정](#환경별-l1-설정)
6. [핵심 포인트 요약](#핵심-포인트-요약)

---

## 개요

이 문서는 **Optimism L2 Rollup 시스템**을 실제로 배포하는 방법을 다룹니다.

### 전체 시스템 구성

```
완전한 L2 시스템 =
  L1 Ethereum (공유 인프라) +
  Sequencer 스택 (블록 생성) +
  Challenger 스택 (검증 및 챌린지)
```

### 사전 지식

이 문서를 읽기 전에 먼저 L2 아키텍처를 이해하시는 것을 권장합니다:
- **[L2 시스템 아키텍처 가이드](./l2-system-architecture-ko.md)** ⭐

---

## 완전한 시스템 아키텍처

### L2 + Challenger 시스템 구성

```
┌────────────────────────────────────────────────────────────┐
│                     완전한 L2 시스템                        │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌─────────────────────────────────────┐                  │
│  │  L1 Ethereum (필수 인프라!) ⭐      │                  │
│  ├─────────────────────────────────────┤                  │
│  │  • Mainnet: 실제 Ethereum           │                  │
│  │  • Sepolia: 테스트넷                 │                  │
│  │  • Anvil/Geth: 로컬 개발용          │                  │
│  └───────────┬─────────────────────────┘                  │
│              │                                            │
│              │ RPC 연결 (모든 컴포넌트가 공유!)            │
│              │                                            │
│  ┌───────────┴─────────────────────────┐                  │
│  │  Sequencer 스택                     │                  │
│  ├─────────────────────────────────────┤                  │
│  │  op-node (sequencer mode)           │                  │
│  │  op-geth (L2 실행)                  │                  │
│  │  op-batcher (L1에 배치 제출)        │                  │
│  │  op-proposer (L1에 Output 제출)    │                  │
│  └───────────┬─────────────────────────┘                  │
│              │                                            │
│              │ (별도 서버/컨테이너)                        │
│              │                                            │
│  ┌───────────┴─────────────────────────┐                  │
│  │  Challenger 스택 (독립 검증!)       │                  │
│  ├─────────────────────────────────────┤                  │
│  │  op-node (follower mode) ⭐         │                  │
│  │  op-geth (L2 검증용) ⭐              │                  │
│  │  op-challenger (챌린지 수행)        │                  │
│  └─────────────────────────────────────┘                  │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### 상세 데이터 흐름

```
┌────────────────────────────────────────────────────────────┐
│                        L1 Ethereum                         │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌─────────────────┐  ┌──────────────────────────┐        │
│  │ Batch Inbox     │  │ DisputeGameFactory       │        │
│  │ 0xff00...       │  │ (Fault Proof 게임 관리)   │        │
│  └────▲────────────┘  └────▲─────────────────────┘        │
│       │                    │                              │
└───────┼────────────────────┼──────────────────────────────┘
        │                    │
        │ write              │ write (게임 생성)
        │                    │
┌───────┴────────────────────┴──────────────────────────────┐
│              Sequencer 스택                               │
├───────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐              ┌──────────────┐          │
│  │ op-batcher   │              │ op-proposer  │          │
│  │ (L1에 쓰기)  │              │ (게임 생성)  │          │
│  └──────▲───────┘              └──────▲───────┘          │
│         │                             │                  │
│         │ read L2 blocks              │ read state       │
│         │                             │                  │
│  ┌──────┴─────────────────────────────┴───────┐          │
│  │      op-node (sequencer mode)             │          │
│  │              ↕ Engine API                  │          │
│  │      op-geth (L2 실행 엔진)                │          │
│  │              ↑ user txs                    │          │
│  └────────────────────────────────────────────┘          │
│                                                           │
└───────────────────────────────────────────────────────────┘
        │
        │ P2P Gossip (optional)
        │
┌───────▼───────────────────────────────────────────────────┐
│              Challenger 스택 ⭐                           │
├───────────────────────────────────────────────────────────┤
│                                                           │
│  ┌─────────────────────────────────────────┐             │
│  │      op-challenger                      │             │
│  │      (검증 & 챌린지)                     │             │
│  └──────▲──────────────────────────────────┘             │
│         │ read L2 state + 게임 참여                      │
│         │                                                │
│  ┌──────┴──────────────────────────────────┐             │
│  │      op-node (follower mode)            │             │
│  │              ↕ Engine API                │             │
│  │      op-geth (L2 검증용, 별도 DB!)       │             │
│  └────────────────▲──────────────────────────┘            │
│                   │                                      │
│                   │ read batches from L1                 │
└───────────────────┼──────────────────────────────────────┘
                    │
                    ▼
        L1 Ethereum (읽기: Batch Inbox, DisputeGameFactory)
```

### 왜 별도의 Challenger 스택이 필요한가?

**독립적 검증을 위해!**

```
Sequencer 노드를 공유하는 경우 (❌ 비권장):
├─ Sequencer가 악의적이면 Challenger도 속을 수 있음
├─ 독립적 검증 불가능
└─ Challenger의 신뢰성 보장 불가

별도 Challenger 스택 (✅ 권장):
├─ L1 데이터로부터 완전히 독립적으로 L2 상태 재구성
├─ Sequencer를 신뢰할 필요 없음
├─ Sequencer의 잘못된 Output을 감지 가능
└─ 진정한 Fault Proof 시스템
```

---

## 실제 배포 구성

### docker-compose.yml (완전판)

```yaml
version: '3.8'

services:
  # ============================================
  # L1 Ethereum (로컬 개발용)
  # ============================================
  l1-geth:
    image: ethereum/client-go:latest  # ⭐ 일반 Ethereum geth!
    command:
      - --dev
      - --http
      - --http.addr=0.0.0.0
      - --http.port=8545
      - --http.api=eth,net,web3,debug
      - --ws
      - --ws.addr=0.0.0.0
      - --ws.port=8546
    ports:
      - "8545:8545"  # HTTP RPC
      - "8546:8546"  # WebSocket
    volumes:
      - l1_data:/data

  # ============================================
  # Sequencer 스택
  # ============================================
  sequencer-op-node:
    image: tokamaknetwork/thanos-op-node:latest
    command:
      - op-node
      - --l1=http://l1-geth:8545              # ⭐ L1 연결
      - --l1.beacon=http://l1-beacon:5052     # ⭐ L1 Beacon (블롭용)
      - --l2=http://sequencer-op-geth:8551
      - --l2.jwt-secret=/jwt.hex
      - --sequencer.enabled=true              # Sequencer 모드
      - --sequencer.l1-confs=4
      - --rollup.config=/rollup.json
      - --rpc.addr=0.0.0.0
      - --rpc.port=8545                       # 내부 포트 8545
      - --p2p.listen.tcp=9003
      - --p2p.listen.udp=9003
    ports:
      - "7545:8545"  # 외부:내부
      - "9003:9003"
    volumes:
      - ./rollup.json:/rollup.json
      - ./jwt.hex:/jwt.hex
    depends_on:
      - l1-geth                               # ⭐ L1 필요!
      - sequencer-op-geth
    environment:
      - OP_NODE_L1_ETH_RPC=http://l1-geth:8545

  sequencer-op-geth:
    image: tokamaknetwork/thanos-op-geth:latest  # ⭐ L2 op-geth
    command:
      - geth
      - --datadir=/data
      - --http
      - --http.addr=0.0.0.0
      - --http.port=8545
      - --http.api=web3,eth,debug,txpool,net,engine
      - --authrpc.addr=0.0.0.0
      - --authrpc.port=8551
      - --authrpc.jwtsecret=/jwt.hex
      - --syncmode=full
      - --nodiscover
      - --maxpeers=0
    ports:
      - "8555:8545"  # 사용자 RPC
    volumes:
      - sequencer_l2_data:/data
      - ./jwt.hex:/jwt.hex
      - ./genesis-l2.json:/genesis.json

  op-batcher:
    image: tokamaknetwork/thanos-op-batcher:latest
    command:
      - op-batcher
      - --l1-eth-rpc=http://l1-geth:8545           # ⭐ L1에 쓰기
      - --l2-eth-rpc=http://sequencer-op-geth:8545
      - --rollup-rpc=http://sequencer-op-node:8545  # 내부 포트
      - --poll-interval=1s
      - --sub-safety-margin=6
      - --num-confirmations=1
      - --safe-abort-nonce-too-low-count=3
      - --resubmission-timeout=30s
      - --rpc.addr=0.0.0.0
      - --rpc.port=8548
      - --private-key=${BATCHER_PRIVATE_KEY}
    depends_on:
      - l1-geth                                    # ⭐ L1 필요!
      - sequencer-op-node
      - sequencer-op-geth

  op-proposer:
    image: tokamaknetwork/thanos-op-proposer:latest
    command:
      - op-proposer
      - --l1-eth-rpc=http://l1-geth:8545           # ⭐ L1에 쓰기
      - --l2-eth-rpc=http://sequencer-op-geth:8545
      - --rollup-rpc=http://sequencer-op-node:8545  # 내부 포트
      - --poll-interval=12s
      - --rpc.addr=0.0.0.0
      - --rpc.port=8545                            # 내부 포트 8545
      - --game-factory-address=${GAME_FACTORY_ADDRESS}  # DisputeGameFactory
      - --proposal-interval=${PROPOSAL_INTERVAL:-30s}
      - --private-key=${PROPOSER_PRIVATE_KEY}
    depends_on:
      - l1-geth                                    # ⭐ L1 필요!
      - sequencer-op-node

  # ============================================
  # Challenger 스택 (독립적!)
  # ============================================
  challenger-op-node:
    image: tokamaknetwork/thanos-op-node:latest
    command:
      - op-node
      - --l1=http://l1-geth:8545                   # ⭐ 같은 L1!
      - --l1.beacon=http://l1-beacon:5052
      - --l2=http://challenger-op-geth:8551        # ⭐ 별도 L2!
      - --l2.jwt-secret=/jwt.hex
      - --sequencer.enabled=false                  # Follower 모드
      - --rollup.config=/rollup.json
      - --rpc.addr=0.0.0.0
      - --rpc.port=8545                            # 내부 포트 8545
      - --p2p.listen.tcp=9004
      - --p2p.listen.udp=9004
      - --p2p.bootnodes=${SEQUENCER_P2P_ENR}       # P2P 연결
    ports:
      - "7546:8545"  # 외부:내부 (Sequencer와 다른 포트!)
      - "9004:9004"
    volumes:
      - ./rollup.json:/rollup.json
      - ./jwt.hex:/jwt.hex
    depends_on:
      - l1-geth                                    # ⭐ L1 공유!
      - challenger-op-geth

  challenger-op-geth:
    image: tokamaknetwork/thanos-op-geth:latest    # ⭐ 별도 L2 geth!
    command:
      - geth
      - --datadir=/data
      - --http
      - --http.addr=0.0.0.0
      - --http.port=8545
      - --http.api=web3,eth,debug,txpool,net,engine
      - --authrpc.addr=0.0.0.0
      - --authrpc.port=8551
      - --authrpc.jwtsecret=/jwt.hex
      - --syncmode=full
      - --nodiscover
      - --maxpeers=0
    ports:
      - "8556:8545"  # 다른 포트
    volumes:
      - challenger_l2_data:/data                   # ⭐ 별도 DB!
      - ./jwt.hex:/jwt.hex
      - ./genesis-l2.json:/genesis.json

  op-challenger:
    image: tokamaknetwork/thanos-op-challenger:latest
    command:
      - op-challenger
      - --l1-eth-rpc=http://l1-geth:8545                  # ⭐ L1 읽기
      - --l1-beacon=http://l1-beacon:5052
      - --l2-eth-rpc=http://challenger-op-geth:8545       # ⭐ 자기 L2
      - --rollup-rpc=http://challenger-op-node:8545       # ⭐ 자기 op-node (내부 포트)
      - --game-factory-address=${GAME_FACTORY_ADDRESS}
      - --trace-type=${CHALLENGER_TRACE_TYPE}             # ⭐ 모든 GameType 지원
      - --datadir=/data
      - --private-key=${CHALLENGER_PRIVATE_KEY}
    volumes:
      - challenger_data:/data
      - ${PROJECT_ROOT}/cannon/bin:/cannon                # VM 바이너리 마운트
      - ${PROJECT_ROOT}/asterisc/bin:/asterisc
      - ${PROJECT_ROOT}/op-program/bin:/op-program
      - ${PROJECT_ROOT}/bin:/kona                         # kona-client (GameType 3)
    depends_on:
      - l1-geth                                           # ⭐ L1 공유!
      - challenger-op-node
      - challenger-op-geth
    environment:
      # Challenger는 모든 GameType 지원
      # CHALLENGER_TRACE_TYPE 예시: "asterisc-kona,asterisc,cannon"

volumes:
  l1_data:
  sequencer_l2_data:
  challenger_l2_data:
  challenger_data:
```

### 환경 변수 설정 (.env)

```bash
# Batcher (Sequencer가 L1에 배치 제출)
BATCHER_PRIVATE_KEY=0x...

# Proposer (Sequencer가 L1에 Output 제출)
PROPOSER_PRIVATE_KEY=0x...
GAME_FACTORY_ADDRESS=0x...  # DisputeGameFactory 컨트랙트 주소

# Challenger
CHALLENGER_PRIVATE_KEY=0x...
CHALLENGER_TRACE_TYPE=cannon,asterisc,asterisc-kona  # 모든 GameType 지원

# GameType 설정
DG_TYPE=0  # Proposer가 사용할 GameType (0, 1, 2, 3)

# P2P (Follower가 Sequencer와 연결, 선택사항)
SEQUENCER_P2P_ENR=enr://...  # Sequencer의 P2P ENR

# VM 바이너리 경로
PROJECT_ROOT=/path/to/tokamak-thanos
```

> 📚 **실제 배포**: 이 문서의 docker-compose는 참고용입니다. 실제 배포는 [배포 스크립트](../../scripts/README-KR.md)를 사용하세요.

---

## 리소스 요구사항

### 프로덕션 환경

| 컴포넌트 | CPU | 메모리 | 디스크 | 필수도 |
|---------|-----|--------|--------|--------|
| **L1 geth** | 8 Core | 16 GB | 2 TB | ✅ 필수 (또는 외부 RPC) |
| **Sequencer op-node** | 2 Core | 4 GB | 100 GB | ✅ 필수 |
| **Sequencer op-geth** | 4 Core | 8 GB | 500 GB | ✅ 필수 |
| **op-batcher** | 1 Core | 2 GB | 10 GB | ✅ 필수 |
| **op-proposer** | 1 Core | 2 GB | 10 GB | ✅ 필수 |
| **Challenger op-node** | 2 Core | 4 GB | 100 GB | ⭐ 검증용 |
| **Challenger op-geth** | 4 Core | 8 GB | 500 GB | ⭐ 검증용 |
| **op-challenger** | 1 Core | 4 GB | 50 GB | ⭐ 챌린지용 |

**총 리소스 (완전 스택):**
- CPU: 19 Core
- 메모리: 44 GB
- 디스크: 3.27 TB

### 개발/테스트 환경

| 컴포넌트 | CPU | 메모리 | 디스크 |
|---------|-----|--------|--------|
| **L1 geth (dev)** | 2 Core | 4 GB | 50 GB |
| **Sequencer op-node** | 1 Core | 2 GB | 20 GB |
| **Sequencer op-geth** | 2 Core | 4 GB | 50 GB |
| **op-batcher** | 0.5 Core | 1 GB | 5 GB |
| **op-proposer** | 0.5 Core | 1 GB | 5 GB |
| **Challenger op-node** | 1 Core | 2 GB | 20 GB |
| **Challenger op-geth** | 2 Core | 4 GB | 50 GB |
| **op-challenger** | 0.5 Core | 2 GB | 10 GB |

**총 리소스 (개발):**
- CPU: 9.5 Core
- 메모리: 20 GB
- 디스크: 210 GB

### 간소화 옵션

#### 옵션 1: 외부 L1 RPC 사용 (권장)

```yaml
# L1 geth 없이 외부 RPC 사용
services:
  sequencer-op-node:
    environment:
      - OP_NODE_L1_ETH_RPC=https://sepolia.infura.io/v3/${API_KEY}
      - OP_NODE_L1_BEACON=https://sepolia.beaconcha.in

# 절약: L1 geth 리소스 (8 Core, 16 GB, 2 TB)
```

**절약 효과:**
- 프로덕션: 11 Core, 28 GB, 1.27 TB
- 개발: 7.5 Core, 16 GB, 160 GB

#### 옵션 2: Challenger가 Sequencer geth 공유 (비권장)

```yaml
op-challenger:
  environment:
    - OP_CHALLENGER_L2_ETH_RPC=http://sequencer-op-geth:8545  # ⚠️

# 절약: Challenger op-geth + op-node 리소스 (6 Core, 12 GB, 600 GB)
```

**단점:**
- ❌ 독립적 검증 불가
- ❌ Sequencer가 잘못된 데이터를 제공하면 Challenger도 속음
- ❌ 진정한 Fault Proof 불가능

**⚠️ 프로덕션 환경에서는 절대 사용하지 마세요!**

---

## 환경별 L1 설정

### 로컬 개발 환경

```yaml
# 완전히 독립적인 로컬 devnet
l1-geth:
  image: ethereum/client-go:latest
  command:
    - --dev                    # 개발 모드
    - --dev.period=1           # 1초마다 블록
    - --http
    - --http.addr=0.0.0.0
    - --http.port=8545

# 장점: 완전한 제어, 빠른 테스트
# 단점: 실제 환경과 차이
```

**사용 시나리오:**
- 로컬 개발 및 테스트
- CI/CD 파이프라인
- 단위 테스트

### 테스트넷 환경 (Sepolia)

```yaml
# L1 geth 없이 공개 RPC 사용
services:
  sequencer-op-node:
    environment:
      - OP_NODE_L1_ETH_RPC=https://sepolia.infura.io/v3/${INFURA_API_KEY}
      - OP_NODE_L1_BEACON=https://sepolia.beaconcha.in

  op-batcher:
    environment:
      - OP_BATCHER_L1_ETH_RPC=https://sepolia.infura.io/v3/${INFURA_API_KEY}

  op-proposer:
    environment:
      - OP_PROPOSER_L1_ETH_RPC=https://sepolia.infura.io/v3/${INFURA_API_KEY}

  challenger-op-node:
    environment:
      - OP_NODE_L1_ETH_RPC=https://sepolia.infura.io/v3/${INFURA_API_KEY}
      - OP_NODE_L1_BEACON=https://sepolia.beaconcha.in

  op-challenger:
    environment:
      - OP_CHALLENGER_L1_ETH_RPC=https://sepolia.infura.io/v3/${INFURA_API_KEY}

# 장점: 실제 네트워크, 리소스 절약
# 단점: API 키 필요, 속도 제한
```

**공개 RPC 제공자:**
- Infura: https://infura.io
- Alchemy: https://alchemy.com
- Ankr: https://ankr.com
- LlamaNodes: https://llamanodes.com

**사용 시나리오:**
- 통합 테스트
- 스테이징 환경
- 베타 테스트

### 프로덕션 환경 (Mainnet)

```yaml
# 자체 L1 Full Node 운영 (권장)
l1-geth:
  image: ethereum/client-go:latest
  command:
    - --mainnet
    - --syncmode=snap         # 또는 full
    - --http
    - --http.addr=0.0.0.0
    - --http.port=8545
    - --http.api=eth,net,web3
    - --ws
    - --ws.addr=0.0.0.0
    - --ws.port=8546
    - --authrpc.jwtsecret=/jwt.hex
    - --datadir=/data
  volumes:
    - l1_mainnet_data:/data   # 2TB+ 필요

# 또는 외부 RPC (백업용)
services:
  sequencer-op-node:
    environment:
      - OP_NODE_L1_ETH_RPC=http://l1-geth:8545
      - OP_NODE_L1_ETH_RPC_BACKUP=https://eth.llamarpc.com

# 장점: 안정성, 제어권
# 단점: 높은 리소스 요구
```

**권장 설정:**
- ✅ 자체 L1 Full Node (주 RPC)
- ✅ 외부 RPC (백업)
- ✅ 모니터링 및 알림 설정
- ✅ 자동 장애 조치 (Failover)

---

## 핵심 포인트 요약

### 1. L2 = op-node + op-geth

```
L2는 단일 프로세스가 아닙니다!

op-node:  Rollup 로직, L1↔L2 연결
op-geth:  EVM 실행, State 관리

둘 다 필요하며, Engine API로 통신합니다.
```

### 2. op-node의 두 모드

```
Sequencer 모드:
├─ 블록 생성 (Producer)
├─ --sequencer.enabled=true
└─ 1개만 운영 (현재)

Follower 모드:
├─ 블록 검증 (Verifier)
├─ --sequencer.enabled=false (기본값)
└─ 무제한 운영 가능
```

### 3. L1은 모두가 공유

```
L1 Ethereum:
├─ Sequencer 스택 (읽기 + 쓰기)
├─ Challenger 스택 (읽기)
└─ 하나의 L1을 모두 사용!

L2 op-geth:
├─ Sequencer의 L2 (별도 DB)
├─ Challenger의 L2 (별도 DB)
└─ 각자 독립적으로 운영!
```

### 4. Challenger는 독립 스택 필요

```
진정한 Fault Proof를 위해:

op-challenger는 독립된 op-node + op-geth 사용
├─ L1 데이터로부터 완전히 독립적으로 L2 재구성
├─ Sequencer 신뢰 불필요
└─ 악의적 Output 감지 가능 ✅
```

---

## 빠른 시작 가이드

### 1. 사전 준비

```bash
# 필수 파일 준비
touch .env
touch jwt.hex
touch rollup.json
touch genesis-l2.json

# JWT 시크릿 생성
openssl rand -hex 32 > jwt.hex
```

### 2. 환경 변수 설정

```bash
# .env 파일 편집
cat > .env << EOF
BATCHER_PRIVATE_KEY=0x...
PROPOSER_PRIVATE_KEY=0x...
CHALLENGER_PRIVATE_KEY=0x...
GAME_FACTORY_ADDRESS=0x...      # DisputeGameFactory 주소
CHALLENGER_TRACE_TYPE=cannon,asterisc,asterisc-kona  # 모든 GameType 지원
DG_TYPE=0                       # Proposer가 사용할 GameType (0, 1, 2, 3)
PROPOSAL_INTERVAL=30s            # 제안 간격
SEQUENCER_P2P_ENR=enr://...     # Sequencer P2P (선택사항)
INFURA_API_KEY=...              # 외부 RPC 사용 시
EOF
```

### 3. 시스템 시작

```bash
# 모든 서비스 시작
docker-compose up -d

# 로그 확인
docker-compose logs -f

# 특정 서비스만 시작
docker-compose up -d l1-geth sequencer-op-node sequencer-op-geth
```

### 4. 상태 확인

```bash
# Sequencer 상태
curl http://localhost:9545/health

# L2 블록 높이
cast block-number --rpc-url http://localhost:8555

# Challenger 동기화 상태
curl http://localhost:9546/health
```

---

## 관련 문서

### 📘 학습 자료

- **[L2 시스템 아키텍처 가이드](./l2-system-architecture-ko.md)**
  - L2의 기본 구조 이해
  - op-node와 op-geth의 관계
  - Engine API 통신 방식

- **[Batch Inbox Address 상세 분석](./research/batch-inbox-address-ko.md)**
  - L1↔L2 데이터 전달 방식
  - Batch Inbox의 구조와 보안

### 🛠️ 운영 가이드

- **[Blob Pruning Risk Analysis](./research/blob-pruning-risk-analysis-ko.md)**
  - Blob 데이터 관리
  - 위험 분석 및 대응

### 🔗 외부 리소스

- [Optimism Specs](https://specs.optimism.io)
- [OP Stack GitHub](https://github.com/ethereum-optimism/optimism)
- [Optimism Docs](https://docs.optimism.io)

---