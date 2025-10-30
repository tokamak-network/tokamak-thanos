# Challenger 시스템 아키텍처

> **전체 시스템 구성 및 Fault Proof Challenger 시스템 상세 가이드**

---

## 📚 목차

1. [개요](#개요)
2. [전체 시스템 구성](#전체-시스템-구성)
3. [컴포넌트 상세 설명](#컴포넌트-상세-설명)
4. [GameType별 아키텍처](#gametype별-아키텍처)
5. [데이터 흐름](#데이터-흐름)
6. [배포 아키텍처](#배포-아키텍처)
7. [보안 및 독립성](#보안-및-독립성)

---

## 개요

Optimism 기반의 L2 Rollup 시스템은, **Fault Proof 시스템**을 통해 L2 상태의 정확성을 검증합니다. 이 문서는 전체 시스템의 아키텍처와 특히 **Challenger 시스템**의 구성을 자세히 설명합니다.

### 핵심 특징

- ✅ **다중 GameType 지원**: Cannon (MIPS), Asterisc (RISC-V), AsteriscKona (RISC-V + Rust)
- ✅ **독립적인 Challenger 스택**: Sequencer와 완전히 분리된 검증 시스템
- ✅ **모듈식 아키텍처**: GameType별 VM과 Executor를 선택적으로 배포
- ✅ **재현 가능한 빌드**: Docker를 통한 deterministic 빌드

---

## 전체 시스템 구성

### 1. High-Level 아키텍처

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                         L1 (Ethereum)                                         │
│  ┌────────────────┐  ┌──────────────────┐   ┌────────────────────────────┐   │
│  │ Batcher Inbox  │  │ OptimismPortal2  │   │ DisputeGameFactory         │   │
│  │ (Batch Data)   │  │ (사용자 입출금)      │  │  (Fault Proof Games)        │   │
│  └────────────────┘  └─────────┬────────┘   │ - Output 검증              │   │
│         ▲                  ▲   │   │        │ - GameType 0,1,2,3...      │   │
│         │                  │   │   │        └────────────────────────────┘   │
│         │                  │   │   │                   ▲           ▲         │
│         │                  │   │   └──────참조──────────┘           │         │
│         │                  │   │     (respectedGameType 게임 검증)   │         │
└─────────┼──────────────────┼───┼───────────────────────────────────┼──────────┘
          │                  │   │                                   │
          │ Batch 제출   L1→L2 │ L2→L1                       Output/Challenge
          │               입금 │ 출금                         제출/이의 제기
          │                  │   │                                   │
┌─────────┴─────────┐     [사용자]      ┌─────────────────────────────┴─────┐
│ Sequencer Stack   │                   │  op-proposer / op-challenger      │
│                   │                   │                                   │
├───────────────────┤                   │  ┌──────────────┐  ┌────────────┐ │
│ sequencer-l2      │                   │  │ op-proposer  │  │ Challenger │ │
│ (L2 geth)         │                   │  │ (게임 생성)  │  │   Stack    │ │
├───────────────────┤                   │  └──────────────┘  │ (독립적!)⭐│ │
│ sequencer-op-node │                   │                    └────────────┘ │
│ (Sequencer 모드)  │                   └───────────────────────────────────┘
├───────────────────┤                          │                      │
│ op-batcher        │                          │                      │
│ (Batch 제출)      │                   DisputeGameFactory    challenger-l2
└───────────────────┘                       .create()        challenger-op-node
                                          (새 게임 생성)      op-challenger
```

**주요 컴포넌트 역할**:

| 컴포넌트 | 역할 | 사용자 | 상호작용 |
|---------|------|--------|---------|
| **Batcher Inbox** | L2 트랜잭션 배치 데이터 저장 | op-batcher | Batch 데이터 제출 |
| **OptimismPortal2** | L1↔L2 메시지/자산 전달 | 일반 사용자 | 입금(`depositTransaction`), 출금 증명(`proveWithdrawalTransaction`), 출금 완료(`finalizeWithdrawalTransaction`) |
| **DisputeGameFactory** | Fault Proof 게임 관리 | op-proposer, op-challenger | 게임 생성(`create`), 게임 참여(`move`, `step`) |

**핵심 포인트**:
- ✅ **OptimismPortal2 ← DisputeGameFactory**: OptimismPortal2는 `respectedGameType` 게임의 출력(rootClaim)을 신뢰하여 출금 검증
- ✅ **op-proposer → DisputeGameFactory**: `.create()`로 새 게임 생성 (L2 Output 제안)
- ✅ **op-challenger → DisputeGameFactory**: 게임에 참여하여 잘못된 Output 검증
- ✅ **사용자 → OptimismPortal2**: 입출금 사용 (DisputeGameFactory와는 간접적으로만 연결)

### 2. 시스템 계층 구조

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ op-batcher   │  │ op-proposer  │  │ op-challenger    │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │
┌─────────────────────────────────────────────────────────────┐
│                    Rollup Protocol Layer                     │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              op-node (Rollup 합의)                     │ │
│  │  - 파생 (Derivation): L1 배치 → L2 블록               │ │
│  │  - 동기화 (Sync): 다른 노드와 상태 동기화              │ │
│  │  - P2P: 블록 전파 및 수신                              │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │
┌─────────────────────────────────────────────────────────────┐
│                    Execution Layer                           │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              L2 geth (EVM 실행)                        │ │
│  │  - 트랜잭션 실행                                        │ │
│  │  - 상태 저장 (독립 DB) ⭐                              │ │
│  │  - RPC 제공                                             │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │
┌─────────────────────────────────────────────────────────────┐
│                    Data Availability Layer                   │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                L1 Ethereum                             │ │
│  │  - Batch Data (calldata/blobs)                        │ │
│  │  - Output Roots                                        │ │
│  │  - Dispute Game 상태                                   │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 컴포넌트 상세 설명

### L1 Layer

#### 1. **Batcher Inbox (SystemConfig에서 관리)**
- **역할**: L2 트랜잭션 배치 데이터 저장 주소
- **데이터**: 압축된 L2 블록 데이터 (calldata 또는 blobs)
- **접근**: op-node가 L1의 Batcher Inbox 주소로 전송된 트랜잭션 데이터를 읽어서 L2 블록으로 재구성
- **주의**: OptimismPortal은 L1↔L2 메시지 전달용이며, Batcher Inbox와는 다름

#### 2. **OptimismPortal2**
- **역할**: L1과 L2 간 메시지 및 자산 전달 (브릿지)
- **사용자**: 일반 사용자 (입출금하는 사람들)
- **주요 함수**:
  - `depositTransaction()`: L1 → L2 자산 입금 (ETH/ERC20)
  - `proveWithdrawalTransaction(_disputeGameIndex)`: L2 → L1 출금 증명
  - `finalizeWithdrawalTransaction()`: 출금 완료 (증명 후 지연 시간 경과)

- **DisputeGameFactory와의 관계** ⭐:
  ```solidity
  // OptimismPortal2 내부 코드
  GameType public respectedGameType;  // 신뢰하는 GameType (예: 3)

  function proveWithdrawalTransaction(..., uint256 _disputeGameIndex, ...) {
      // DisputeGameFactory에서 게임 가져오기
      (GameType gameType,, IDisputeGame gameProxy) =
          disputeGameFactory.gameAtIndex(_disputeGameIndex);

      // respectedGameType만 신뢰
      require(gameType.raw() == respectedGameType.raw());

      // CHALLENGER_WINS 게임은 사용 불가
      require(gameProxy.status() != GameStatus.CHALLENGER_WINS);

      // 게임의 rootClaim을 출금 증명에 사용
      Claim outputRoot = gameProxy.rootClaim();
  }
  ```

- **핵심**: OptimismPortal2는 DisputeGameFactory에 **의존하여 출금을 검증**합니다
  - ✅ 사용자가 OptimismPortal2로 입출금함
  - ✅ OptimismPortal2는 `respectedGameType` 게임만 신뢰함
  - ✅ 출금 시 DisputeGameFactory의 게임 결과(rootClaim, status)를 검증에 사용함
  - ❌ op-proposer는 OptimismPortal2와 직접 상호작용하지 않음

#### 3. **Output Oracle (DisputeGameFactory 기반)**
- **역할**: L2 상태 루트 제출 및 검증
- **주기**: Proposer가 정기적으로 제출 (예: 30초마다)
- **검증**: DisputeGameFactory를 통한 Fault Proof 게임
- **주의**: 최신 시스템은 L2OutputOracle 대신 DisputeGameFactory 사용

#### 4. **Dispute Games (DisputeGameFactory)**
- **역할**: Fault Proof 게임 관리
- **GameType**: 0, 1, 2, 3, 254, 255 지원
- **구현**: 각 GameType별로 다른 VM (MIPS.sol, RISCV.sol)
- **게임 생성**: `create(GameType, Claim, extraData)`
- **게임 검증**: 각 GameType의 구현 컨트랙트 (`gameImpls` 매핑)

### Sequencer Stack

#### 1. **sequencer-l2 (L2 geth)**
```yaml
역할: L2 트랜잭션 실행 및 블록 생성
모드: Sequencer (Leader)
DB: sequencer_l2_data (독립 볼륨)
포트: 9545 (RPC)
```

#### 2. **sequencer-op-node**
```yaml
역할: Rollup 합의 및 L1 동기화
모드: Sequencer (Leader)
연결:
  - L1 (http://l1:8545)
  - sequencer-l2 (Engine API)
포트: 7545 (RPC)
```

#### 3. **op-batcher**
```yaml
역할: L2 블록을 L1에 배치로 제출
연결:
  - L1 (배치 제출)
  - sequencer-op-node (L2 블록 가져오기)
키: BATCHER_PRIVATE_KEY
```

#### 4. **op-proposer**
```yaml
역할: L2 Output을 DisputeGameFactory에 제출
연결:
  - L1 DisputeGameFactory (게임 생성)
  - sequencer-op-node (상태 루트 가져오기)
키: PROPOSER_PRIVATE_KEY
주기: PROPOSAL_INTERVAL (예: 30s)
방식: DisputeGameFactory.create()로 새 게임 생성
```

### Challenger Stack (독립적!) ⭐

#### 1. **challenger-l2 (L2 geth)**
```yaml
역할: 독립적인 L2 실행 검증
모드: Follower (읽기 전용)
DB: challenger_l2_data (독립 볼륨) ⭐
포트: 9546 (RPC)
중요: Sequencer와 완전히 분리된 DB!
```

**왜 독립적이어야 하나?**
- Sequencer가 악의적일 경우를 대비
- L1에서 독립적으로 L2 재구성
- 진정한 검증 (trustless)

#### 2. **challenger-op-node**
```yaml
역할: 독립적인 Rollup 합의
모드: Follower (검증 전용) ⭐
연결:
  - L1 (http://l1:8545)
  - challenger-l2 (Engine API)
포트: 7546 (RPC)
중요: L1에서 직접 배치 읽어서 재구성!
```

#### 3. **op-challenger**
```yaml
역할: 잘못된 Output에 대한 Challenge 실행
연결:
  - L1 (게임 참여)
  - challenger-op-node (검증된 상태)
  - challenger-l2 (독립 실행 결과)
키: CHALLENGER_PRIVATE_KEY
TraceType: cannon, asterisc, asterisc-kona
```

**op-challenger의 동작**:
1. L1의 Dispute Games 모니터링
2. 잘못된 Output 감지
3. Binary search로 첫 번째 불일치 찾기
4. VM (Cannon/Asterisc) 실행하여 증명 생성
5. L1에 증명 제출

---

## GameType별 아키텍처

### GameType 0/1: Cannon (MIPS VM)

```
┌─────────────────────────────────────────────────────────┐
│                  GameType 0/1 (Cannon)                  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  On-Chain (L1):                                         │
│  ┌────────────────────────────────────────────────┐    │
│  │           MIPS.sol (MIPS VM)                   │    │
│  │  - Single-step instruction 검증                │    │
│  │  - Memory merkle proof 검증                    │    │
│  │  - PreimageOracle 지원                         │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Off-Chain (Challenger):                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  Cannon VM (MIPS Emulator)                     │    │
│  │  - op-program 실행 (Go)                        │    │
│  │  - Trace 생성 (각 instruction)                 │    │
│  │  - Prestate: prestate.json                     │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  빌드:                                                   │
│  ┌────────────────────────────────────────────────┐    │
│  │  make reproducible-prestate                    │    │
│  │  → Docker로 Linux 바이너리 생성                │    │
│  │  → cannon binary                               │    │
│  │  → op-program binary                           │    │
│  │  → prestate-proof.json                         │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  파일 구조:                                             │
│  /cannon/bin/cannon             ← MIPS VM              │
│  /op-program/bin/op-program     ← Server (Go)          │
│  /op-program/bin/prestate.json  ← Prestate             │
└─────────────────────────────────────────────────────────┘
```

### GameType 2: Asterisc (RISC-V VM)

```
┌─────────────────────────────────────────────────────────┐
│                  GameType 2 (Asterisc)                  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  On-Chain (L1):                                         │
│  ┌────────────────────────────────────────────────┐    │
│  │          RISCV.sol (RISC-V VM) ⭐               │    │
│  │  - RV64GC instruction 검증                     │    │
│  │  - Memory merkle proof 검증                    │    │
│  │  - PreimageOracle 지원                         │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Off-Chain (Challenger):                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  Asterisc VM (RISC-V Emulator)                 │    │
│  │  - op-program 실행 (Go, RISC-V 컴파일)        │    │
│  │  - Trace 생성 (각 instruction)                 │    │
│  │  - Prestate: prestate-proof-rv64.json          │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  빌드:                                                   │
│  ┌────────────────────────────────────────────────┐    │
│  │  Asterisc 소스 클론:                           │    │
│  │  git clone github.com/ethereum-optimism/       │    │
│  │             asterisc                           │    │
│  │                                                 │    │
│  │  make reproducible-prestate                    │    │
│  │  → Docker로 Linux 바이너리 생성                │    │
│  │  → asterisc binary (RISC-V VM)                 │    │
│  │  → op-program binary (RISC-V 타겟)            │    │
│  │  → prestate-proof.json (.pre)                  │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  파일 구조:                                             │
│  /asterisc/bin/asterisc                ← RISC-V VM     │
│  /op-program/bin/op-program            ← Server (Go)   │
│  /asterisc/bin/prestate-proof.json     ← Prestate      │
│  /asterisc/bin/prestate.json → prestate-proof.json     │
│                                   (심볼릭 링크)         │
│                                                          │
│  특징:                                                   │
│  - RISC-V 아키텍처 (Cannon의 MIPS 대체)              │
│  - 더 현대적인 ISA                                     │
│  - 동일한 op-program 사용 (다른 타겟)                 │
└─────────────────────────────────────────────────────────┘
```

### GameType 3: AsteriscKona (RISC-V + Rust) 🆕

```
┌─────────────────────────────────────────────────────────┐
│              GameType 3 (AsteriscKona) 🆕               │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  On-Chain (L1):                                         │
│  ┌────────────────────────────────────────────────┐    │
│  │     RISCV.sol (GameType 2와 동일!) ⭐⭐        │    │
│  │  - GameType 2와 동일한 구현 재사용             │    │
│  │  - 온체인 검증 로직 100% 동일                  │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Off-Chain (Challenger):                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  Asterisc VM (RISC-V Emulator)                 │    │
│  │  - kona-client 실행 (Rust!) ⭐⭐              │    │
│  │  - Trace 생성 (각 instruction)                 │    │
│  │  - Prestate: prestate-kona.json                │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  빌드:                                                   │
│  ┌────────────────────────────────────────────────┐    │
│  │  Kona 소스 클론:                               │    │
│  │  git clone github.com/op-rs/kona               │    │
│  │                                                 │    │
│  │  Docker 빌드 (자동화됨):                      │    │
│  │  docker run ghcr.io/op-rs/kona/              │    │
│  │    asterisc-builder:0.3.0                     │    │
│  │    cargo build -Zbuild-std=core,alloc         │    │
│  │      -p kona-client                            │    │
│  │      --profile release-client-lto             │    │
│  │  → target/riscv64imac-unknown-none-elf/       │    │
│  │            release-client-lto/kona-client     │    │
│  │                                                 │    │
│  │  Prestate 생성 (Docker 내부):                  │    │
│  │  docker run golang:1.22-bookworm              │    │
│  │    ./asterisc load-elf                         │    │
│  │      --path /kona/.../kona-client             │    │
│  │      --out prestate-kona.json                 │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  파일 구조 (간략):                                      │
│  /asterisc/bin/asterisc             ← RISC-V VM        │
│  /bin/kona-client                   ← Server (Rust) ⭐ │
│  /op-program/bin/prestate-kona.json ← Prestate         │
│                                                          │
│  차이점 (vs GameType 2):                               │
│  - 온체인 VM: 동일 (RISCV.sol)                        │
│  - Server: kona-client (Rust) vs op-program (Go)      │
│  - 장점: ~80% 경량화, ZK proof 통합 준비             │
│                                                          │
│  Fallback:                                              │
│  - kona prestate 없으면 Asterisc prestate 사용        │
│    (동일한 RISCV.sol이므로 호환 가능)                 │
└─────────────────────────────────────────────────────────┘
```

---

## 📂 GameType 2 vs GameType 3 상세 폴더 구조

### GameType 2 (Asterisc - Go 기반)

```
프로젝트 디렉토리 구조:
tokamak-projects/                        # 작업 디렉토리
│
└── tokamak-thanos/                      # 프로젝트 루트 (이 저장소)
    │
    ├── op-program/                      # ← GameType 2 소스 코드 (Go, 프로젝트 내부)
    │   ├── client/                      # op-program 클라이언트
    │   │   ├── main.go                  # Go 언어로 작성
    │   │   ├── driver.go
    │   │   └── ...
    │   │
    │   ├── host/                        # op-program 호스트
    │   │
    │   └── bin/
    │       ├── op-program               # ← GameType 2 바이너리 (Go 빌드)
    │       └── prestate.json            # Cannon용 prestate
    │
    └── asterisc/                        # ← VM (GameType 2, 3 공통)
        └── bin/
            ├── asterisc                 # RISC-V VM 실행 파일 (12MB)
            └── prestate-proof.json      # ← GameType 2 prestate (63MB)
```

**빌드 방법**:
```bash
# Go 기반 op-program 빌드
cd tokamak-thanos/op-program
go build -o bin/op-program ./client

# Asterisc VM 빌드 (Docker)
cd tokamak-thanos/asterisc
make reproducible-prestate
```

---

### GameType 3 (AsteriscKona - Rust 기반)

```
프로젝트 디렉토리 구조:
tokamak-projects/                        # 작업 디렉토리
│
├── kona/                                # ← GameType 3 소스 코드 (별도 저장소!)
│   │                                    # GitHub: https://github.com/op-rs/kona
│   │                                    # deploy-modular.sh가 자동으로 클론
│   │
│   ├── bin/
│   │   ├── client/                      # ← kona-client 소스 (Rust)
│   │   │   ├── src/
│   │   │   │   ├── main.rs            # Rust 언어로 작성 ⭐
│   │   │   │   └── lib.rs
│   │   │   └── Cargo.toml             # Rust 패키지 설정
│   │   │
│   │   ├── host/                        # kona-host (실행 환경)
│   │   ├── node/                        # kona-node
│   │   └── ...
│   │
│   ├── crates/                          # Rust 크레이트들
│   │   ├── executor/
│   │   ├── preimage/
│   │   └── ...
│   │
│   └── target/                          # Rust 빌드 결과물
│       └── riscv64imac-unknown-none-elf/  # ← RISC-V bare metal 타겟
│           └── release-client-lto/       # ← LTO 최적화 프로파일
│               └── kona-client           # ← RISC-V ELF (prestate 생성 + 실행용)
│
└── tokamak-thanos/                      # 프로젝트 루트
    ├── bin/                             # ← GameType 3 배포 시 생성
    │   └── kona-client                  # ← GameType 3 바이너리 (kona에서 복사, 배포 후 유지)
    │
    ├── op-program/
    │   └── bin/
    │       ├── prestate.json            # Cannon용
    │       └── prestate-kona.json       # ← GameType 3 prestate (배포 시 생성, 유지)
    │
    └── asterisc/                        # ← VM (GameType 2, 3 공통)
        └── bin/
            ├── asterisc                 # RISC-V VM (GameType 2와 동일!)
            └── prestate-proof.json      # GameType 2 prestate
```


### 핵심 차이점 정리

| 항목 | GameType 2 (Asterisc) | GameType 3 (AsteriscKona) |
|------|----------------------|--------------------------|
| **프로그램 언어** | Go | Rust ⭐ |
| **소스 코드 위치** | `tokamak-thanos/op-program/` | `../kona/` (별도 저장소) ⭐ |
| **빌드된 바이너리** | `tokamak-thanos/op-program/bin/op-program` | `tokamak-thanos/bin/kona-client` |
| **빌드 도구** | Go compiler | Cargo (Rust) |
| **VM 실행 파일** | `tokamak-thanos/asterisc/bin/asterisc` | `tokamak-thanos/asterisc/bin/asterisc` (동일!) ✅ |
| **Prestate 위치** | `tokamak-thanos/asterisc/bin/prestate-proof.json` | `tokamak-thanos/op-program/bin/prestate-kona.json` |
| **바이너리 크기** | ~100MB | ~20MB (~80% 감소) ⭐ |

**위치 관계**:
```
tokamak-projects/          # 작업 디렉토리
├── kona/                  # GameType 3 소스 (별도 저장소)
└── tokamak-thanos/        # 프로젝트 루트
    ├── op-program/        # GameType 2 소스
    ├── bin/kona-client    # GameType 3 바이너리 (kona에서 복사)
    └── asterisc/bin/      # VM (공통)
```

**자동 배포**:
```bash
# deploy-modular.sh가 자동으로 처리
cd tokamak-thanos
./op-challenger/scripts/deploy-modular.sh --dg-type 3

# 배포 후 생성되는 파일:
# ✅ bin/kona-client (kona에서 복사, 유지됨)
# ✅ op-program/bin/prestate-kona.json (asterisc로 생성, 유지됨)
# ✅ ../kona/ (저장소 클론, 유지됨)
```

**배포 전/후 비교**:
```
배포 전:                        배포 후:
tokamak-thanos/                 tokamak-thanos/
├── op-program/                 ├── bin/                    ← 생성됨
│   └── bin/                    │   └── kona-client         ← 생성됨
│       └── prestate.json       ├── op-program/
└── asterisc/                   │   └── bin/
    └── bin/                    │       ├── prestate.json
        └── asterisc            │       └── prestate-kona.json  ← 생성됨
                                └── asterisc/
                                    └── bin/
                                        └── asterisc

kona/ 없음                       ../kona/                    ← 클론됨
                                ├── bin/client/
                                └── target/riscv64imac-unknown-none-elf/
                                    └── release-client-lto/
                                        └── kona-client
```

---

### GameType 비교표

| 구성 요소 | GameType 0/1 (Cannon) | GameType 2 (Asterisc) | GameType 3 (AsteriscKona) |
|----------|---------------------|---------------------|-------------------------|
| **온체인 VM** | MIPS.sol | RISCV.sol | RISCV.sol (동일) ⭐ |
| **ISA** | MIPS | RISC-V (RV64GC) | RISC-V (RV64GC) |
| **Off-chain VM** | cannon | asterisc | asterisc (동일) |
| **Server** | op-program (Go) | op-program (Go) | kona-client (Rust) ⭐ |
| **Server 언어** | Go | Go | Rust |
| **Prestate 필드 (배포용)** | `.pre` | `.pre` | `.pre` ⭐ |
| **빌드 도구** | Go + Docker | Go + Docker | Rust + Cargo + Docker |
| **바이너리 크기** | ~100MB | ~100MB | ~20MB (~80% 감소) ⭐ |
| **장점** | 안정성, Go 에코시스템 | 현대적 ISA, RISC-V 표준 | 경량화, ZK 통합 준비 |
| **상태** | ✅ 프로덕션 | ✅ 프로덕션 | 🆕 신규 통합 |

---

## 데이터 흐름

### 1. 정상 흐름 (Happy Path)

```
┌─────────────────────────────────────────────────────────────┐
│                   1. L2 트랜잭션 실행                        │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        User Transaction → Sequencer L2 → 블록 생성
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              2. L1에 배치 데이터 제출 (op-batcher)           │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
              Sequencer op-node → 블록 데이터 수집
                          │
                          ▼
              op-batcher → L1 BatcherInbox에 제출
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│           3. L2 Output 제출 (op-proposer)                    │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
              op-proposer → DisputeGameFactory에 제출
                          (새로운 게임 생성: 제안된 Output)
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│           4. Challenger 독립 검증 (병렬)                     │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐              ┌──────────────────┐
│ L1에서 배치 읽기  │              │ 제출된 Output    │
│ (challenger-     │              │ 읽기             │
│  op-node)        │              │                  │
└──────────────────┘              └──────────────────┘
        │                                   │
        ▼                                   ▼
┌──────────────────┐              ┌──────────────────┐
│ L2 재구성        │              │ 로컬 계산한      │
│ (challenger-l2)  │              │ Output와 비교    │
└──────────────────┘              └──────────────────┘
        │                                   │
        └─────────────────┬─────────────────┘
                          │
                          ▼
                    ┌─────────────┐
                    │  일치?       │
                    └─────────────┘
                      │         │
              일치 ✅ │         │ 불일치 ❌
                      │         │
                      ▼         ▼
                  (정상)    Challenge 시작!
```

### 2. Challenge 흐름 (Dispute)

```
┌─────────────────────────────────────────────────────────────┐
│              1. 잘못된 Output 감지 (op-challenger)           │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        Challenger의 로컬 상태 ≠ L1 제출된 Output
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   2. Dispute Game 시작                       │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        L1 DisputeGameFactory.create()
          → 새로운 게임 인스턴스 생성
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              3. Binary Search (Bisection)                    │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────────┐
        │  Output Root (전체 블록 범위)       │
        │  ↓                                   │
        │  Claim at depth 1 (절반)            │
        │  ↓                                   │
        │  Claim at depth 2 (1/4)             │
        │  ↓                                   │
        │  ...                                 │
        │  ↓                                   │
        │  Single Instruction (최종 depth)    │
        └─────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│         4. VM 실행 및 증명 생성 (Single Step)               │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        GameType에 따라 다른 VM 사용:
        ┌───────────────────────────────────┐
        │ GameType 0/1: cannon (MIPS)       │
        │ GameType 2:   asterisc (RISC-V)   │
        │ GameType 3:   asterisc + kona     │
        └───────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────────┐
        │  1. VM이 해당 instruction 실행      │
        │  2. Pre-state 준비                  │
        │  3. Post-state 계산                 │
        │  4. Memory merkle proof 생성        │
        │  5. Preimage 데이터 준비            │
        └─────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              5. L1에 증명 제출 (step())                      │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        L1 FaultDisputeGame.step()
          → VM (MIPS.sol / RISCV.sol)에서 검증
          → 온체인 실행 결과와 비교
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   6. 게임 해결 (resolve)                     │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
  Challenger 승리 ✅             Proposer 승리 ❌
  (잘못된 Output 제거)          (Output 유지)
```

### 3. Challenger의 독립 검증 상세

```
┌─────────────────────────────────────────────────────────────┐
│              Challenger의 독립 검증 프로세스                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Step 1: L1에서 배치 데이터 읽기                            │
│  ┌────────────────────────────────────────────────────┐    │
│  │  challenger-op-node                                │    │
│  │  → L1 BatcherInbox 모니터링                       │    │
│  │  → 배치 데이터 다운로드                            │    │
│  │  → 압축 해제                                       │    │
│  │  → L2 트랜잭션 복원                                │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  Step 2: L2 독립 재구성                                     │
│  ┌────────────────────────────────────────────────────┐    │
│  │  challenger-l2 (독립 geth DB!)                     │    │
│  │  → 복원된 트랜잭션 순서대로 실행                   │    │
│  │  → 독립 상태 DB에 저장                             │    │
│  │  → Output Root 계산                                │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  Step 3: Output 비교                                        │
│  ┌────────────────────────────────────────────────────┐    │
│  │  로컬 계산 Output Root                             │    │
│  │         vs                                          │    │
│  │  L1 제출된 Output Root (DisputeGameFactory)       │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  Step 4: 불일치 감지 시                                     │
│  ┌────────────────────────────────────────────────────┐    │
│  │  1. Dispute Game 생성                              │    │
│  │  2. Binary search로 첫 불일치 instruction 찾기     │    │
│  │  3. VM 실행하여 증명 생성                          │    │
│  │  4. L1에 증명 제출                                 │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  중요: Sequencer 컴포넌트와 0% 공유! ⭐⭐                   │
│  - 독립 L2 geth                                             │
│  - 독립 op-node                                             │
│  - 독립 DB                                                  │
│  → 진정한 검증 (trustless)                                 │
└─────────────────────────────────────────────────────────────┘
```

---

## 배포 아키텍처

### Docker Compose 구조

**개념적 구조** (실제 설정은 더 복잡합니다):

```yaml
services:
  # L1 Layer
  l1:
    ports: ["8545:8545"]
    volumes: [l1_data:/db]

  # ============================================================
  # Sequencer Stack
  # ============================================================
  sequencer-l2:
    ports: ["9545:8545"]              # Sequencer L2 RPC
    volumes: [sequencer_l2_data:/db]  # Sequencer DB
    environment:
      ROLLUP_CLIENT_HTTP: http://sequencer-op-node:8545

  sequencer-op-node:
    ports: ["7545:8545"]              # Sequencer op-node RPC
    environment:
      OP_NODE_L1_ETH_RPC: http://l1:8545
      OP_NODE_L2_ENGINE_RPC: http://sequencer-l2:8551
      OP_NODE_SEQUENCER_ENABLED: true  # ← Sequencer 모드

  op-batcher:
    environment:
      OP_BATCHER_ROLLUP_RPC: http://sequencer-op-node:8545

  op-proposer:
    environment:
      OP_PROPOSER_ROLLUP_RPC: http://sequencer-op-node:8545

  # ============================================================
  # Challenger Stack (완전히 독립!) ⭐
  # ============================================================
  challenger-l2:
    ports: ["9546:8545"]                  # Challenger L2 RPC (다른 포트!)
    volumes: [challenger_l2_data:/db]    # ← 독립 DB ⭐⭐
    environment:
      ROLLUP_CLIENT_HTTP: http://challenger-op-node:8545

  challenger-op-node:
    ports: ["7546:8545"]                  # Challenger op-node RPC (다른 포트!)
    environment:
      OP_NODE_L1_ETH_RPC: http://l1:8545  # ← L1에서 직접 읽기
      OP_NODE_L2_ENGINE_RPC: http://challenger-l2:8551
      OP_NODE_SEQUENCER_ENABLED: false    # ← Follower 모드 ⭐

  op-challenger:
    volumes:
      # 모든 GameType의 VM 바이너리 마운트
      - ${PROJECT_ROOT}/cannon/bin:/cannon        # GameType 0/1
      - ${PROJECT_ROOT}/asterisc/bin:/asterisc    # GameType 2, 3
      - ${PROJECT_ROOT}/bin:/kona                 # GameType 3 (kona-client)
      - ${PROJECT_ROOT}/op-program/bin:/op-program
    environment:
      OP_CHALLENGER_L1_ETH_RPC: http://l1:8545
      OP_CHALLENGER_ROLLUP_RPC: http://challenger-op-node:8545  # ⭐ 자신의 op-node
      OP_CHALLENGER_L2_ETH_RPC: http://challenger-l2:8545        # ⭐ 자신의 L2 geth
      OP_CHALLENGER_TRACE_TYPE: ${CHALLENGER_TRACE_TYPE}  # 모든 GameType 지원
      # ... GameType별 상세 환경 변수 ...

volumes:
  l1_data:
  sequencer_l2_data:      # ← Sequencer L2 DB
  challenger_l2_data:     # ← Challenger L2 DB (독립!) ⭐⭐
  challenger_data:        # ← Challenger 작업 디렉토리
```

**핵심 포인트**:
- ✅ Challenger는 **별도 L2 geth** (`challenger-l2`) 사용
- ✅ Challenger는 **별도 op-node** (`challenger-op-node`, Follower 모드) 사용
- ✅ Challenger는 **별도 DB** (`challenger_l2_data`) 사용
- ✅ 모든 GameType의 VM 바이너리를 마운트하여 **모든 게임 타입 지원**
- ✅ Sequencer와 **0% 공유** → 진정한 독립 검증

> 📚 **실제 Docker Compose 설정**: [docker-compose-full.yml](../scripts/docker-compose-full.yml)

### 포트 매핑

| 서비스 | 포트 | 설명 |
|--------|------|------|
| **L1** | 8545 | L1 Ethereum RPC |
| **Sequencer L2** | 9545 | 사용자 트랜잭션 제출 |
| **Sequencer op-node** | 7545 | Rollup 정보 조회 |
| **Challenger L2** | 9546 | 독립 검증 RPC ⭐ |
| **Challenger op-node** | 7546 | 독립 Rollup 정보 ⭐ |
| **op-batcher** | 6545 | Batcher 관리 |
| **op-proposer** | 6546 | Proposer 관리 |
| **op-challenger** | 6547 | Challenger 관리 |

---

## 보안 및 독립성

### 1. Challenger 독립성의 중요성

#### 왜 독립적이어야 하나?

```
시나리오: Sequencer가 악의적으로 행동

┌─────────────────────────────────────────────────────────┐
│  만약 Challenger가 Sequencer와 노드/DB를 공유한다면...   │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  악의적 Sequencer:                                      │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. 잘못된 L2 블록 생성                        │    │
│  │  2. 공유 DB에 잘못된 상태 저장                 │    │
│  │  3. op-node도 공유하므로 잘못된 정보 제공      │    │
│  └────────────────────────────────────────────────┘    │
│            ▼                                             │
│  Challenger (공유 시):                                  │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. 공유 DB에서 상태 읽기                      │    │
│  │     → 이미 오염된 데이터!                      │    │
│  │  2. 공유 op-node에서 정보 가져오기             │    │
│  │     → 이미 오염된 정보!                        │    │
│  │  3. 잘못된 Output도 "정상"으로 판단 ❌        │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  결과: Challenger가 속는다! ❌❌❌                      │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  독립적인 Challenger Stack을 사용하면...                │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  악의적 Sequencer:                                      │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. 잘못된 L2 블록 생성                        │    │
│  │  2. 잘못된 Output을 L1에 제출                  │    │
│  └────────────────────────────────────────────────┘    │
│            ▼                                             │
│  독립 Challenger:                                        │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. L1에서 직접 배치 데이터 읽기 ⭐           │    │
│  │     → Sequencer 우회, L1만 신뢰               │    │
│  │  2. 독립 DB에서 L2 재구성 ⭐                   │    │
│  │     → 깨끗한 상태에서 처음부터 재계산         │    │
│  │  3. 로컬 Output vs L1 제출 Output 비교        │    │
│  │     → 불일치 감지! ✅                          │    │
│  │  4. Dispute Game 시작                          │    │
│  │     → 잘못된 Output 제거 ✅                    │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  결과: Challenger가 정확히 검증! ✅✅✅                │
└─────────────────────────────────────────────────────────┘
```

### 2. 독립성 검증 체크리스트

#### 배포 시 반드시 확인:

```bash
# ✅ 독립 DB 확인
docker volume ls | grep l2_data
# 출력:
# scripts_sequencer_l2_data    ← Sequencer DB
# scripts_challenger_l2_data   ← Challenger DB (다른 볼륨!)

# ✅ 독립 op-node 확인
docker exec challenger-op-node ps aux | grep op-node
# 환경변수 확인:
# OP_NODE_SEQUENCER_ENABLED=false  ← Follower 모드

# ✅ 연결 확인
docker logs challenger-op-node | grep "L1 source"
# 출력:
# L1 source: http://l1:8545  ← L1에서 직접 읽기

# ✅ DB 격리 확인
docker exec sequencer-l2 ls /db      # Sequencer DB
docker exec challenger-l2 ls /db     # Challenger DB (다른 디렉토리!)
```

### 3. Trust Model

```
┌─────────────────────────────────────────────────────────┐
│                    Trust Model                           │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  신뢰 필요:                                             │
│  ┌────────────────────────────────────────────────┐    │
│  │  ✅ L1 Ethereum (공개 블록체인)                │    │
│  │  ✅ Consensus 레이어 (검증됨)                  │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  신뢰 불필요 (검증 가능):                              │
│  ┌────────────────────────────────────────────────┐    │
│  │  ⭐ Sequencer (Challenger가 검증)             │    │
│  │  ⭐ Proposer (Challenger가 검증)              │    │
│  │  ⭐ Batcher (L1에 기록, 누구나 읽기 가능)     │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Fault Proof의 핵심:                                    │
│  ┌────────────────────────────────────────────────┐    │
│  │  "신뢰하지 말고, 검증하라"                      │    │
│  │  (Don't trust, verify)                         │    │
│  │                                                 │    │
│  │  → Challenger가 모든 것을 독립 검증            │    │
│  │  → 잘못된 행위는 증명으로 제거                 │    │
│  └────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## 요약

### 핵심 아키텍처 원칙

1. **독립성**: Challenger는 Sequencer와 완전히 분리
2. **검증 가능성**: 모든 데이터는 L1에서 직접 검증
3. **모듈성**: GameType별로 다른 VM 선택 가능
4. **재현성**: Docker로 deterministic 빌드

### GameType별 특징

- **GameType 0/1 (Cannon)**: MIPS VM, 안정적, Go 기반
- **GameType 2 (Asterisc)**: RISC-V VM, 현대적, Go 기반
- **GameType 3 (AsteriscKona)**: RISC-V VM, 경량화, Rust 기반 🆕

### 배포 시 주의사항

1. ✅ Challenger 독립 DB 확인
2. ✅ Follower 모드 확인
3. ✅ L1 직접 연결 확인
4. ✅ GameType별 VM 바이너리 준비
5. ✅ Prestate 올바르게 생성

---

## 참고 자료

- [Optimism Specs](https://specs.optimism.io)
- [OP Stack GitHub](https://github.com/ethereum-optimism/optimism)
- [Fault Proof 문서](https://specs.optimism.io/fault-proof/index.html)
- [배포 가이드](../scripts/README-KR.md)
---
