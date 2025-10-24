# Asterisc (RISC-V) 완벽 가이드

## 목차
1. [개요](#개요)
2. [RISC-V 아키텍처](#risc-v-아키텍처)
3. [Asterisc vs Cannon 비교](#asterisc-vs-cannon-비교)
4. [Asterisc 동작 원리](#asterisc-동작-원리)
5. [설치 및 빌드](#설치-및-빌드)
6. [설정 및 실행](#설정-및-실행)
7. [증명 생성 프로세스](#증명-생성-프로세스)
8. [온체인 검증](#온체인-검증)
9. [성능 최적화](#성능-최적화)
10. [트러블슈팅](#트러블슈팅)
11. [FAQ](#faq)

---

## 개요

**Asterisc**는 Optimism Fault Proof 시스템에서 **RISC-V 아키텍처**를 사용하는 VM 기반 증명 시스템입니다. Cannon (MIPS)의 대안으로, 클라이언트 다양성(Client Diversity)을 제공하여 시스템의 안정성과 보안을 강화합니다.

### 주요 특징

- **RISC-V 64-bit ISA**: 오픈소스 명령어 세트 아키텍처 사용
- **Fraud Proof 방식**: Optimistic Rollup의 dispute resolution
- **온체인 검증**: 단일 RISC-V instruction을 온체인에서 실행
- **Cannon과 호환**: 동일한 op-program을 다른 ISA로 컴파일

### Asterisc의 역할

```
┌────────────────────────────────────────────────────────────────┐
│                    Dispute Game 생명주기                        │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. Output Proposal (L2 → L1)                                  │
│     └─> Proposer: "Block 1000의 state root는 0xabcd..."       │
│                                                                 │
│  2. Challenge (Challenger 발견)                                │
│     └─> Challenger: "이의 제기! 잘못된 output"                 │
│                                                                 │
│  3. Bisection Game (Asterisc 사용)                             │
│     ┌───────────────────────────────────────────────────┐     │
│     │ Depth 0-29: L2 output 검증 (빠름)                │     │
│     │  └─> RollupClient.OutputAtBlock()                 │     │
│     │                                                     │     │
│     │ Depth 30-73: RISC-V VM 실행 (Asterisc)            │     │
│     │  ├─> op-program 실행 (RISC-V 바이너리)            │     │
│     │  ├─> Instruction-by-instruction 시뮬레이션         │     │
│     │  └─> Proof 생성 (.json.gz)                        │     │
│     └───────────────────────────────────────────────────┘     │
│                                                                 │
│  4. Step Execution (온체인)                                    │
│     ┌───────────────────────────────────────────────────┐     │
│     │ RiscV.sol 컨트랙트                                 │     │
│     │  ├─> Prestate 로드                                 │     │
│     │  ├─> 단일 RISC-V instruction 실행                 │     │
│     │  └─> Poststate 검증                                │     │
│     └───────────────────────────────────────────────────┘     │
│                                                                 │
│  5. Resolution                                                 │
│     └─> 올바른 주장자 승리, Bond 분배                         │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

---

## RISC-V 아키텍처

### RISC-V란?

**RISC-V** (Reduced Instruction Set Computer - Five)는 버클리 대학에서 개발한 **오픈소스 ISA**입니다.

### 주요 특징

1. **오픈소스**: 누구나 무료로 사용 가능 (특허 제약 없음)
2. **단순함**: 최소한의 기본 명령어 세트
3. **확장 가능**: 모듈식 확장 (RV32I, RV64I, M, A, F, D, C 등)
4. **효율성**: 최적화된 instruction encoding

### Asterisc가 사용하는 RISC-V 확장

```
RV64GC = RV64I + M + A + F + D + C

RV64I:  64-bit Integer Base
M:      Integer Multiplication/Division
A:      Atomic Instructions
F:      Single-Precision Floating-Point
D:      Double-Precision Floating-Point
C:      Compressed Instructions (16-bit)
```

### RISC-V Instruction 예시

```assembly
# 간단한 덧셈
addi x1, x0, 10    # x1 = x0 + 10 = 10
addi x2, x0, 20    # x2 = x0 + 20 = 20
add  x3, x1, x2    # x3 = x1 + x2 = 30

# 메모리 접근
ld   x4, 0(x1)     # x4 = Memory[x1 + 0] (Load Doubleword)
sd   x4, 8(x2)     # Memory[x2 + 8] = x4 (Store Doubleword)

# 조건 분기
beq  x1, x2, label # if (x1 == x2) goto label
```

### RISC-V 레지스터

```
x0:  Zero (항상 0)
x1:  ra (Return Address)
x2:  sp (Stack Pointer)
x3:  gp (Global Pointer)
x4:  tp (Thread Pointer)
x5-x7:   Temporaries
x8:  s0/fp (Saved/Frame Pointer)
x9:  s1 (Saved Register)
x10-x11: Function Arguments/Return Values
x12-x17: Function Arguments
x18-x27: Saved Registers
x28-x31: Temporaries

PC:  Program Counter (다음 실행 instruction 주소)
```

---

## Asterisc vs Cannon 비교

### 아키텍처 비교

| 항목 | Cannon (MIPS) | Asterisc (RISC-V) |
|-----|--------------|------------------|
| **ISA** | MIPS32 | RISC-V 64-bit (RV64GC) |
| **레지스터 수** | 32 (32-bit) | 32 (64-bit) + PC |
| **주소 공간** | 32-bit (4GB) | 64-bit (이론상 제한 없음) |
| **Instruction 크기** | 32-bit 고정 | 32-bit + 16-bit (compressed) |
| **Floating Point** | 소프트웨어 에뮬레이션 | 하드웨어 지원 (F, D 확장) |
| **라이선스** | 독점적 (MIPS Technologies) | 오픈소스 (BSD) |

### 구현 비교

| 항목 | Cannon | Asterisc |
|-----|--------|----------|
| **VM 바이너리** | `cannon` | `asterisc` |
| **op-program** | `op-program` (MIPS) | `op-program-rv64` (RISC-V) |
| **Prestate** | `prestate.json` (~190MB) | `prestate-rv64.json` (~200MB) |
| **온체인 컨트랙트** | `MIPS.sol` | `RiscV.sol` (개발 중) |
| **GameType** | 0, 1 | 2 |
| **설정 플래그** | `--cannon-*` | `--asterisc-*` |

### 성능 비교

| 항목 | Cannon | Asterisc | 비고 |
|-----|--------|----------|------|
| **증명 생성 시간** | 5-30분 | 5-30분 | 유사 (하드웨어 의존) |
| **메모리 사용량** | 2-8GB | 2-8GB | 유사 |
| **디스크 사용량** | 1-10GB/게임 | 1-10GB/게임 | 유사 |
| **온체인 가스 비용** | ~5M gas | ~5M gas (예상) | 유사 |
| **Snapshot 빈도** | 1B instructions | 1B instructions | 설정 가능 |

### 코드 호환성

**동일한 op-program 소스 코드**를 사용하지만, 다른 타겟으로 컴파일합니다:

```bash
# Cannon (MIPS)
GOARCH=mips GOOS=linux go build -o op-program-mips

# Asterisc (RISC-V)
GOARCH=riscv64 GOOS=linux go build -o op-program-rv64
```

---

## Asterisc 동작 원리

### 전체 흐름

```
1. Prestate 준비
   ├─ op-program 컴파일 (RISC-V 바이너리)
   └─ Prestate 생성 (초기 VM 상태)

2. 게임 발견
   ├─ DisputeGameFactory 모니터링
   ├─ GameType == 2 (Asterisc) 확인
   └─ GamePlayer 생성

3. Trace 생성 (오프체인)
   ├─ Position 계산
   ├─ TraceIndex 도출
   └─ Asterisc VM 실행
       ├─ Snapshot 로드 (가장 가까운)
       ├─ TraceIndex까지 실행
       └─ Proof 저장 (state.json.gz)

4. Claim 제출 (온체인)
   ├─ Attack 또는 Defend
   └─ ClaimValue = StateHash at Position

5. Step 증명 (온체인, 최종)
   ├─ Prestate 제출
   ├─ RiscV.sol 실행
   ├─ Poststate 계산
   └─ 검증
```

### 핵심 컴포넌트

#### 1. AsteriscTraceProvider

**파일**: `op-challenger/game/fault/trace/asterisc/provider.go`

**역할**: 특정 position의 VM 상태 제공

```go
type AsteriscTraceProvider struct {
    logger         log.Logger
    dir            string              // 게임 데이터 디렉토리
    prestate       string              // prestate.json 경로
    generator      utils.ProofGenerator // Executor
    gameDepth      types.Depth
    preimageLoader *utils.PreimageLoader
    PrestateProvider
    lastStep       uint64              // 캐싱 최적화
}

// 주요 메소드
func (p *AsteriscTraceProvider) Get(ctx, pos Position) (common.Hash, error)
func (p *AsteriscTraceProvider) GetStepData(ctx, pos Position) ([]byte, []byte, *PreimageOracleData, error)
```

**Get() 흐름**:
```go
// provider.go:58-73
func (p *AsteriscTraceProvider) Get(ctx context.Context, pos types.Position) (common.Hash, error) {
    // 1. Position → TraceIndex 변환
    traceIndex := pos.TraceIndex(p.gameDepth)

    // 2. Proof 로드 (없으면 생성)
    proof, err := p.loadProof(ctx, traceIndex.Uint64())
    if err != nil {
        return common.Hash{}, err
    }

    // 3. ClaimValue (state hash) 반환
    value := proof.ClaimValue
    return value, nil
}
```

#### 2. Executor

**파일**: `op-challenger/game/fault/trace/asterisc/executor.go`

**역할**: Asterisc VM 실행 및 proof 생성

```go
type Executor struct {
    logger           log.Logger
    metrics          AsteriscMetricer
    l1               string           // L1 RPC
    l1Beacon         string           // L1 Beacon API
    l2               string           // L2 RPC
    inputs           utils.LocalGameInputs
    asterisc         string           // asterisc 바이너리 경로
    server           string           // op-program 경로
    network          string
    rollupConfig     string
    l2Genesis        string
    absolutePreState string           // prestate.json 경로
    snapshotFreq     uint             // Snapshot 빈도
    infoFreq         uint             // 정보 출력 빈도
}
```

**GenerateProof() 흐름**:

```go
// executor.go:58-127
func (e *Executor) GenerateProof(ctx context.Context, dir string, i uint64) error {
    // 1. Snapshot 선택
    start, err := e.selectSnapshot(logger, snapshotDir, e.absolutePreState, i)
    // 가장 가까운 이전 snapshot (예: 1B, 2B, ...)

    // 2. Asterisc 실행
    args := []string{
        "run",
        "--input", start,                                    // Snapshot 또는 prestate
        "--output", lastGeneratedState,                      // 최종 상태
        "--proof-at", "=" + strconv.FormatUint(i, 10),       // 증명할 step
        "--proof-fmt", filepath.Join(proofDir, "%d.json.gz"), // 출력 형식
        "--snapshot-at", "%" + strconv.FormatUint(snapshotFreq, 10),
        "--snapshot-fmt", filepath.Join(snapshotDir, "%d.json.gz"),
        "--",
        e.server, "--server",                                // op-program 실행
        "--l1", e.l1,
        "--l2", e.l2,
        "--l1.head", e.inputs.L1Head.Hex(),
        "--l2.head", e.inputs.L2Head.Hex(),
        // ... 기타 입력
    }

    // 3. 실행
    err = e.cmdExecutor(ctx, logger, e.asterisc, args...)

    // 4. Proof 파일 생성
    // {dir}/proofs/{i}.json.gz
}
```

#### 3. Proof 데이터 구조

**파일**: `op-challenger/game/fault/trace/utils/proof.go`

```go
type ProofData struct {
    ClaimValue common.Hash `json:"post"`        // VM 상태 해시 (post state)
    StateData  hexutil.Bytes `json:"state"`     // VM 전체 상태 (직렬화)
    ProofData  hexutil.Bytes `json:"proof"`     // Merkle proof 등
}
```

**StateData 내용** (RISC-V 특화):
```json
{
  "pc": "0x...",           // Program Counter
  "registers": [           // x0-x31 레지스터
    "0x0000000000000000",  // x0 (always zero)
    "0x...",               // x1 (ra)
    "0x...",               // x2 (sp)
    // ... x3-x31
  ],
  "memory": "0x...",       // 전체 메모리 (압축)
  "step": 123456789,       // 현재 step
  "heap": "0x...",         // Heap 포인터
  "preimageKey": "0x...",  // Preimage oracle key
  "preimageOffset": 0      // Preimage oracle offset
}
```

---

## 설치 및 빌드

### 사전 요구사항

```bash
# Go 1.21+
go version

# RISC-V 툴체인 (op-program 빌드용)
# macOS
brew install riscv-tools

# Ubuntu
sudo apt-get install gcc-riscv64-linux-gnu

# Git & Make
git --version
make --version
```

### 저장소 클론

```bash
git clone https://github.com/tokamak-network/tokamak-thanos.git
cd tokamak-thanos
```

### op-program 빌드 (RISC-V 타겟)

```bash
cd op-program

# RISC-V 64-bit 바이너리 빌드
make op-program-rv64

# 출력: ./bin/op-program-rv64
```

**Makefile 내부**:
```makefile
op-program-rv64:
	GOOS=linux GOARCH=riscv64 go build \
		-o ./bin/op-program-rv64 \
		-ldflags="-extldflags=-static" \
		./cmd
```

### Asterisc VM 빌드

**참고**: 현재 tokamak-thanos 저장소에 asterisc가 포함되지 않을 수 있습니다. Optimism 원본 저장소를 사용하거나 별도로 빌드해야 합니다.

```bash
# Optimism 저장소에서 asterisc 가져오기
cd /path/to/optimism
cd asterisc

make asterisc

# 출력: ./bin/asterisc
```

### Prestate 생성

```bash
cd op-program

# Prestate 생성 (RISC-V)
make asterisc-prestate

# 출력: ./bin/prestate-rv64.json
```

**내부 동작**:
```bash
# 1. op-program-rv64 빌드
# 2. Asterisc로 초기 실행
asterisc load-elf --path=./bin/op-program-rv64 --out=./bin/prestate-rv64.json

# 3. 검증
asterisc run --input=./bin/prestate-rv64.json --proof-at="=0" --proof-fmt=./bin/prestate-proof-rv64.json
```

### op-challenger 빌드

```bash
cd op-challenger

make op-challenger

# 출력: ./bin/op-challenger
```

---

## 설정 및 실행

### 기본 설정

#### 환경 변수 방식

```bash
#!/bin/bash

# L1/L2 연결
export OP_CHALLENGER_L1_ETH_RPC=http://localhost:8545
export OP_CHALLENGER_L1_BEACON=http://localhost:5052
export OP_CHALLENGER_ROLLUP_RPC=http://localhost:9546
export OP_CHALLENGER_L2_ETH_RPC=http://localhost:9545

# 게임 설정
export OP_CHALLENGER_TRACE_TYPE=asterisc
export OP_CHALLENGER_GAME_FACTORY_ADDRESS=0x...

# Asterisc 특화 설정
export OP_CHALLENGER_ASTERISC_BIN=/path/to/asterisc/bin/asterisc
export OP_CHALLENGER_ASTERISC_SERVER=/path/to/op-program/bin/op-program-rv64
export OP_CHALLENGER_ASTERISC_PRESTATE=/path/to/op-program/bin/prestate-rv64.json
export OP_CHALLENGER_ASTERISC_ROLLUP_CONFIG=/path/to/.devnet/rollup.json
export OP_CHALLENGER_ASTERISC_L2_GENESIS=/path/to/.devnet/genesis-l2.json
export OP_CHALLENGER_ASTERISC_SNAPSHOT_FREQ=1000000000  # 1B instructions

# 데이터 디렉토리
export OP_CHALLENGER_DATADIR=/data/challenger

# 키 관리
export OP_CHALLENGER_PRIVATE_KEY=$PRIVATE_KEY

# 실행
./op-challenger/bin/op-challenger
```

#### CLI 플래그 방식

```bash
#!/bin/bash

./op-challenger/bin/op-challenger \
  --trace-type asterisc \
  --l1-eth-rpc http://localhost:8545 \
  --l1-beacon http://localhost:5052 \
  --rollup-rpc http://localhost:9546 \
  --l2-eth-rpc http://localhost:9545 \
  --game-factory-address 0x... \
  --datadir /data/challenger \
  \
  --asterisc-bin /usr/local/bin/asterisc \
  --asterisc-server /usr/local/bin/op-program-rv64 \
  --asterisc-prestate /data/prestates/prestate-rv64.json \
  --asterisc-rollup-config /config/rollup.json \
  --asterisc-l2-genesis /config/genesis-l2.json \
  --asterisc-snapshot-freq 1000000000 \
  --asterisc-info-freq 10000000 \
  \
  --private-key $PRIVATE_KEY \
  --max-concurrency 4 \
  --max-pending-tx 10
```

### 설정 파일 구조

#### rollup.json

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
      "overhead": "0x0000000000000000000000000000000000000000000000000000000000000834",
      "scalar": "0x00000000000000000000000000000000000000000000000000000000000f4240",
      "gasLimit": 30000000
    }
  },
  "block_time": 2,
  "max_sequencer_drift": 600,
  "seq_window_size": 3600,
  "channel_timeout": 300,
  "l1_chain_id": 900,
  "l2_chain_id": 901,
  "batch_inbox_address": "0xff00000000000000000000000000000000000901",
  "deposit_contract_address": "0x...",
  "l1_system_config_address": "0x..."
}
```

#### genesis-l2.json

```json
{
  "config": {
    "chainId": 901,
    "homesteadBlock": 0,
    "eip150Block": 0,
    // ... EVM 설정
    "optimism": {
      "eip1559Elasticity": 6,
      "eip1559Denominator": 50
    }
  },
  "nonce": "0x0",
  "timestamp": "0x...",
  "extraData": "0x",
  "gasLimit": "0x1c9c380",
  "difficulty": "0x0",
  "mixHash": "0x...",
  "coinbase": "0x4200000000000000000000000000000000000011",
  "alloc": {
    "0x...": {
      "balance": "0x..."
    }
  }
}
```

### 로컬 Devnet 실행 예시

```bash
#!/bin/bash
set -e

# 1. Devnet 시작
cd /path/to/tokamak-thanos
make devnet-up

# 2. Contract 주소 가져오기
DGF_ADDRESS=$(jq -r .DisputeGameFactoryProxy .devnet/addresses.json)
echo "DisputeGameFactory: $DGF_ADDRESS"

# 3. op-challenger 실행 (Asterisc)
./op-challenger/bin/op-challenger \
  --trace-type asterisc \
  --l1-eth-rpc http://localhost:8545 \
  --rollup-rpc http://localhost:9546 \
  --l2-eth-rpc http://localhost:9545 \
  --game-factory-address $DGF_ADDRESS \
  --datadir temp/asterisc-challenger-data \
  --asterisc-bin /usr/local/bin/asterisc \
  --asterisc-server ./op-program/bin/op-program-rv64 \
  --asterisc-prestate ./op-program/bin/prestate-rv64.json \
  --asterisc-rollup-config .devnet/rollup.json \
  --asterisc-l2-genesis .devnet/genesis-l2.json \
  --mnemonic "test test test test test test test test test test test junk" \
  --hd-path "m/44'/60'/0'/0/8" \
  --num-confirmations 1 \
  --log.level debug
```

### 프로덕션 설정 예시

```bash
#!/bin/bash

./op-challenger \
  --trace-type asterisc \
  --l1-eth-rpc https://mainnet.infura.io/v3/YOUR_KEY \
  --l1-beacon https://beacon-api.example.com \
  --rollup-rpc https://rollup-node.example.com \
  --l2-eth-rpc https://l2-node.example.com \
  --game-factory-address 0x... \
  --datadir /data/challenger \
  \
  --asterisc-bin /usr/local/bin/asterisc \
  --asterisc-server /usr/local/bin/op-program-rv64 \
  --asterisc-prestate-base-url https://prestates.example.com/asterisc/ \
  --asterisc-rollup-config /config/rollup.json \
  --asterisc-l2-genesis /config/genesis-l2.json \
  --asterisc-snapshot-freq 1000000000 \
  --asterisc-info-freq 10000000 \
  \
  --private-key $PRIVATE_KEY \
  --max-concurrency 4 \
  --max-pending-tx 10 \
  --game-window 672h \
  --poll-interval 12s \
  \
  --metrics.enabled \
  --metrics.addr 0.0.0.0 \
  --metrics.port 7300 \
  \
  --pprof.enabled \
  --pprof.addr 0.0.0.0 \
  --pprof.port 6060 \
  \
  --log.level info \
  --log.format json
```

---

## 증명 생성 프로세스

### 단계별 흐름

#### 1. Position → TraceIndex 계산

```go
// Position은 게임 트리에서의 위치
type Position struct {
    depth        Depth
    indexAtDepth *big.Int
}

// TraceIndex는 실제 VM step 번호
func (p Position) TraceIndex(maxDepth Depth) *big.Int {
    // 공식: (indexAtDepth + 1) * 2^(maxDepth - depth) - 1
    return new(big.Int).Sub(
        new(big.Int).Mul(
            new(big.Int).Add(p.indexAtDepth, big.NewInt(1)),
            new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(maxDepth - p.depth)), nil),
        ),
        big.NewInt(1),
    )
}
```

**예시**:
```
maxDepth = 73
Position(depth=30, index=0) → TraceIndex = 8,796,093,022,207
Position(depth=73, index=0) → TraceIndex = 0
Position(depth=73, index=1) → TraceIndex = 1
```

#### 2. Snapshot 선택

```go
// utils/snapshots.go
func FindStartingSnapshot(logger log.Logger, snapshotDir, absolutePreState string, i uint64) (string, error) {
    // 1. Snapshot 파일 목록 가져오기
    //    예: [0.json.gz, 1000000000.json.gz, 2000000000.json.gz, ...]

    // 2. i보다 작거나 같은 가장 큰 snapshot 찾기
    //    i = 1,234,567,890 → 선택: 1000000000.json.gz

    // 3. Snapshot 없으면 prestate 사용
    if closestSnapshot == nil {
        return absolutePreState, nil
    }

    return closestSnapshotPath, nil
}
```

**장점**:
- 전체 실행 시간 단축
- 1B instruction마다 checkpoint
- 메모리 효율적

#### 3. Asterisc 실행

```bash
# 실제 명령어 예시
asterisc run \
  --input /data/challenger/0x.../snapshots/1000000000.json.gz \
  --output /data/challenger/0x.../state.json \
  --proof-at =1234567890 \
  --proof-fmt /data/challenger/0x.../proofs/%d.json.gz \
  --snapshot-at %1000000000 \
  --snapshot-fmt /data/challenger/0x.../snapshots/%d.json.gz \
  -- \
  /usr/local/bin/op-program-rv64 --server \
    --l1 http://localhost:8545 \
    --l2 http://localhost:9545 \
    --l1.head 0x... \
    --l2.head 0x... \
    --l2.outputroot 0x... \
    --l2.claim 0x... \
    --l2.blocknumber 1000 \
    --rollup.config /config/rollup.json \
    --l2.genesis /config/genesis-l2.json \
    --datadir /data/challenger/0x.../preimages
```

**내부 동작**:
1. Snapshot 로드 (VM 상태 복원)
2. op-program-rv64를 게스트 프로그램으로 실행
3. RISC-V instruction 단위로 시뮬레이션
4. 지정된 step마다 상태 저장
5. Proof 파일 생성

#### 4. Proof 검증 (로컬)

```go
// 생성된 proof 읽기
proofPath := filepath.Join(dir, proofsDir, fmt.Sprintf("%d.json.gz", traceIndex))
file, _ := os.Open(proofPath)
defer file.Close()

gzReader, _ := gzip.NewReader(file)
defer gzReader.Close()

var proof ProofData
json.NewDecoder(gzReader).Decode(&proof)

// ClaimValue 확인
if proof.ClaimValue == (common.Hash{}) {
    return errors.New("proof missing post hash")
}

// StateData 확인
if len(proof.StateData) == 0 {
    return errors.New("proof missing state data")
}
```

#### 5. 디렉토리 구조

```
/data/challenger/{game-address}/
├── proofs/
│   ├── 0.json.gz
│   ├── 1000000000.json.gz
│   ├── 1234567890.json.gz  ← 요청한 proof
│   └── ...
├── snapshots/
│   ├── 0.json.gz           ← prestate 복사본
│   ├── 1000000000.json.gz  ← 1B step snapshot
│   ├── 2000000000.json.gz
│   └── ...
├── preimages/
│   └── kv/
│       ├── 0x12...34
│       ├── 0xab...cd
│       └── ...
└── state.json              ← 최종 상태
```

---

## 온체인 검증

### RiscV.sol 컨트랙트

**위치**: `packages/contracts-bedrock/src/cannon/RiscV.sol` (예상)

**주요 함수**:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

contract RiscV {
    /// @notice 단일 RISC-V instruction 실행
    /// @param _stateData VM 상태 (레지스터, 메모리 등)
    /// @param _proof Merkle proof
    /// @param _localContext Preimage oracle data
    /// @return postState_ 실행 후 VM 상태 해시
    function step(
        bytes calldata _stateData,
        bytes calldata _proof,
        bytes32 _localContext
    ) external pure returns (bytes32 postState_) {
        // 1. StateData 디코딩
        State memory state = decodeState(_stateData);

        // 2. 현재 instruction 페치
        uint32 insn = fetchInstruction(state);

        // 3. Instruction 디코딩
        (uint8 opcode, uint8 rd, uint8 rs1, uint8 rs2, int32 imm) = decodeInstruction(insn);

        // 4. Instruction 실행
        if (opcode == 0x33) {  // R-type (add, sub, ...)
            if (funct7 == 0x00 && funct3 == 0x0) {
                // ADD rd, rs1, rs2
                state.registers[rd] = state.registers[rs1] + state.registers[rs2];
            } else if (funct7 == 0x20 && funct3 == 0x0) {
                // SUB rd, rs1, rs2
                state.registers[rd] = state.registers[rs1] - state.registers[rs2];
            }
            // ... 기타 R-type instructions
        } else if (opcode == 0x13) {  // I-type (addi, ...)
            if (funct3 == 0x0) {
                // ADDI rd, rs1, imm
                state.registers[rd] = state.registers[rs1] + imm;
            }
            // ... 기타 I-type instructions
        } else if (opcode == 0x03) {  // Load instructions
            // LD, LW, LH, LB, ...
        } else if (opcode == 0x23) {  // Store instructions
            // SD, SW, SH, SB, ...
        } else if (opcode == 0x63) {  // Branch instructions
            // BEQ, BNE, BLT, BGE, ...
        }
        // ... 모든 RISC-V instruction 처리

        // 5. PC 업데이트
        state.pc += 4;  // 대부분의 경우

        // 6. Step 증가
        state.step += 1;

        // 7. 상태 해시 계산
        postState_ = keccak256(encodeState(state));
    }
}
```

### Step Execution 과정

```
1. Claim 제출 (오프체인 증명)
   └─> Attack/Defend: ClaimValue = StateHash

2. 게임 진행 (Bisection)
   └─> Depth가 maxDepth에 도달

3. Step 호출 (온체인)
   ┌─────────────────────────────────────┐
   │ FaultDisputeGame.step()             │
   ├─────────────────────────────────────┤
   │ 1. Claim 검증                        │
   │    - Parent claim 확인               │
   │    - Position 확인                   │
   │                                      │
   │ 2. Preimage Oracle 업데이트 (필요시) │
   │    - PreimageOracle.loadKeccak256()  │
   │    - PreimageOracle.loadLocalData()  │
   │                                      │
   │ 3. VM.step() 호출                    │
   │    └─> RiscV.step(stateData, proof) │
   │         ├─ Instruction 실행          │
   │         └─ PostState 계산            │
   │                                      │
   │ 4. 결과 검증                         │
   │    if (postState == claimValue) {   │
   │        // Claim 정당                 │
   │    } else {                          │
   │        // Claim 부당 → 상대방 승리   │
   │    }                                 │
   └─────────────────────────────────────┘
```

### 가스 비용 분석

| 항목 | 가스 비용 | 비고 |
|-----|---------|------|
| **step() 호출** | ~5,000,000 gas | Instruction에 따라 다름 |
| - State 디코딩 | ~100,000 gas | |
| - Instruction 실행 | ~50,000 - 1,000,000 gas | 메모리 접근 시 증가 |
| - Preimage 로드 | ~500,000 - 2,000,000 gas | 필요 시 |
| - State 인코딩 | ~100,000 gas | |
| **resolve() 호출** | ~500,000 gas | |
| **Total (게임당)** | ~10,000,000 - 50,000,000 gas | Claim 수에 따라 |

---

## 성능 최적화

### Snapshot 전략

#### 기본 설정

```bash
--asterisc-snapshot-freq 1000000000  # 1B instructions
```

**트레이드오프**:
- Frequency ↑ (예: 500M): 더 빠른 증명 생성, 디스크 사용량 ↑
- Frequency ↓ (예: 5B): 느린 증명 생성, 디스크 절약

#### 적응형 Snapshot

복잡한 게임에서는 더 자주 snapshot:

```bash
# 예: 중요 구간만 500M, 나머지 2B
--asterisc-snapshot-freq 500000000 \
--asterisc-custom-snapshot-ranges "0-1000000000:500000000,1000000000-:2000000000"
```

### 병렬 처리

```bash
# 여러 게임 동시 처리
--max-concurrency 8  # CPU 코어 수에 맞춰 조정
```

**리소스 고려**:
```
메모리: 2GB/게임 * max-concurrency
디스크 I/O: Snapshot 읽기/쓰기
CPU: RISC-V 시뮬레이션
```

### 캐싱 최적화

#### Prestate 캐싱

```bash
# 여러 prestate 지원
--asterisc-prestate-base-url https://prestates.example.com/asterisc/

# 로컬 캐시
--datadir /data/challenger
# 자동 생성: /data/challenger/asterisc-prestates/
```

#### Proof 캐싱

```go
// provider.go:103-145
func (p *AsteriscTraceProvider) loadProof(ctx context.Context, i uint64) (*utils.ProofData, error) {
    // 1. 디스크 캐시 확인
    proofPath := filepath.Join(p.dir, proofsDir, fmt.Sprintf("%d.json.gz", i))
    if _, err := os.Stat(proofPath); err == nil {
        // 캐시 히트 → 파일 읽기
        return loadProofFromFile(proofPath)
    }

    // 2. 캐시 미스 → 생성
    if err := p.generator.GenerateProof(ctx, p.dir, i); err != nil {
        return nil, err
    }

    // 3. 재시도
    return loadProofFromFile(proofPath)
}
```

### 메모리 관리

```bash
# JVM 스타일 메모리 제한 (Go는 자동이지만)
export GOMEMLIMIT=16GiB

# 컨테이너 메모리 제한
docker run --memory=16g --memory-swap=16g ...
```

### 디스크 I/O 최적화

#### SSD 사용 권장

```bash
# NVMe SSD에 데이터 디렉토리 배치
--datadir /nvme/challenger-data
```

#### 압축 레벨 조정

```go
// proof 저장 시 gzip 압축
gzWriter, _ := gzip.NewWriterLevel(file, gzip.BestSpeed)  // 빠른 압축
// 또는
gzWriter, _ := gzip.NewWriterLevel(file, gzip.BestCompression)  // 공간 절약
```

---

## 트러블슈팅

### 문제 1: Asterisc 바이너리를 찾을 수 없음

**증상**:
```
ERRO [01-17|12:00:00] Failed to generate proof error="exec: \"asterisc\": executable file not found in $PATH"
```

**해결**:
```bash
# 1. Asterisc 설치 확인
which asterisc

# 2. 절대 경로 사용
--asterisc-bin /usr/local/bin/asterisc

# 3. 빌드 (없는 경우)
cd /path/to/optimism/asterisc
make asterisc
sudo cp ./bin/asterisc /usr/local/bin/
```

### 문제 2: Prestate 불일치

**증상**:
```
WARN [01-17|12:01:00] Prestate mismatch game=0x... required=0xabcd... actual=0x1234...
```

**원인**:
- 온체인 게임이 요구하는 prestate와 로컬 prestate가 다름
- 네트워크 업그레이드 후 발생 가능

**해결**:
```bash
# 1. 요구되는 prestate 다운로드
wget https://prestates.example.com/asterisc/0xabcd...1234.json -O /data/prestates/0xabcd...1234.json

# 2. Multi-prestate 모드 사용
--asterisc-prestate-base-url https://prestates.example.com/asterisc/

# 3. Unsafe 모드 (테스트만!)
--unsafe-allow-invalid-prestate=true
```

### 문제 3: 메모리 부족

**증상**:
```
ERRO [01-17|12:05:00] Proof generation failed error="signal: killed"
```

**해결**:
```bash
# 1. 메모리 증가
docker run --memory=16g ...

# 2. Concurrency 감소
--max-concurrency 2

# 3. Snapshot 빈도 증가 (메모리 footprint 감소)
--asterisc-snapshot-freq 500000000
```

### 문제 4: 증명 생성이 너무 느림

**증상**:
- 한 proof 생성에 30분 이상 소요

**원인 및 해결**:

```bash
# 1. Snapshot 없음 → prestate부터 매번 실행
# 해결: Snapshot 확인
ls -lh /data/challenger/0x.../snapshots/

# 2. 디스크 I/O 병목
# 해결: SSD 사용, I/O 모니터링
iostat -x 1

# 3. CPU 부족
# 해결: 더 강력한 CPU, concurrency 조정
--max-concurrency 4  # CPU 코어 수에 맞춤
```

### 문제 5: RPC 연결 실패

**증상**:
```
ERRO [01-17|12:10:00] Failed to get L2 output error="Post \"http://localhost:9545\": dial tcp: connection refused"
```

**해결**:
```bash
# 1. L2 노드 상태 확인
curl -X POST http://localhost:9545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# 2. 방화벽 확인
sudo ufw status

# 3. 네트워크 연결 확인 (Docker)
docker network inspect ops-bedrock_default
```

### 문제 6: 게임을 찾을 수 없음

**증상**:
```
INFO [01-17|12:15:00] No games to play
```

**원인**:
- DisputeGameFactory에 게임이 없음
- GameType 필터 문제
- Game allowlist 설정

**해결**:
```bash
# 1. 게임 목록 확인
./op-challenger list-games \
  --l1-eth-rpc http://localhost:8545 \
  --game-factory-address 0x...

# 2. Trace type 확인
--trace-type asterisc  # GameType 2만 처리

# 3. Allowlist 제거 (테스트)
# --game-allowlist 플래그 제거
```

---

## FAQ

### Q1: Asterisc와 Cannon을 동시에 실행할 수 있나요?

**A**: 네, 가능하며 권장됩니다!

```bash
./op-challenger \
  --trace-type cannon,asterisc \
  --cannon-bin /usr/local/bin/cannon \
  --cannon-server /usr/local/bin/op-program \
  --cannon-prestate /data/prestates/prestate.json \
  --asterisc-bin /usr/local/bin/asterisc \
  --asterisc-server /usr/local/bin/op-program-rv64 \
  --asterisc-prestate /data/prestates/prestate-rv64.json \
  # ... 기타 설정
```

**결과**:
- GameType 0, 1 (Cannon) 처리
- GameType 2 (Asterisc) 처리
- 최대 보안 및 클라이언트 다양성

---

### Q2: op-program을 두 번 빌드해야 하나요?

**A**: 네, MIPS와 RISC-V 타겟으로 각각 빌드해야 합니다.

```bash
# Cannon용 (MIPS)
cd op-program
make op-program
# 출력: ./bin/op-program

# Asterisc용 (RISC-V)
make op-program-rv64
# 출력: ./bin/op-program-rv64
```

**차이점**:
- 소스 코드 동일
- 컴파일 타겟만 다름 (`GOARCH=mips` vs `GOARCH=riscv64`)

---

### Q3: RISC-V가 MIPS보다 빠른가요?

**A**: 비슷합니다. 실제 성능은 하드웨어와 구현에 따라 다릅니다.

**벤치마크** (예상):
```
증명 생성 시간 (10,000,000 instructions):
- Cannon (MIPS):   ~5분
- Asterisc (RISC-V): ~5분

메모리 사용:
- Cannon:   ~4GB
- Asterisc: ~4GB
```

**실제로는**:
- VM 구현 품질에 더 의존적
- 하드웨어 최적화 차이
- Snapshot 전략의 영향

---

### Q4: 온체인 가스 비용 차이는?

**A**: 유사할 것으로 예상됩니다.

```
Step 실행:
- MIPS.sol:  ~5M gas
- RiscV.sol: ~5M gas (예상)

이유:
- 단일 instruction 실행
- 비슷한 복잡도
- 메모리 접근 패턴 유사
```

**주의**: RiscV.sol은 아직 최적화 중일 수 있습니다.

---

### Q5: RISC-V의 장점은 무엇인가요?

**A**:

1. **오픈소스**:
   - 특허 제약 없음
   - 커뮤니티 주도 개발
   - 자유로운 구현 가능

2. **생태계**:
   - Linux, GCC, LLVM 완전 지원
   - 활발한 RISC-V Foundation
   - 교육 및 연구 자료 풍부

3. **미래 지향적**:
   - 현대적 ISA 설계
   - 확장 가능한 구조
   - 최신 보안 기능

4. **클라이언트 다양성**:
   - Cannon의 대안
   - 단일 장애점 제거
   - 상호 검증 가능

---

### Q6: Asterisc로 ZK proof를 할 수 있나요?

**A**: 현재 tokamak-thanos의 Asterisc는 **fraud proof (optimistic)** 방식입니다.

하지만 **RISC-V + ZK 조합은 가능**합니다:

**기존 프로젝트**:
- **RISC Zero**: RISC-V 기반 zkVM (zkSTARK)
- **SP1 (Succinct)**: RISC-V zkVM (Plonky3)
- **Valida**: RISC-V zkVM (STARK)

**통합 가능성**:
```go
// 미래 가능성
if cfg.TraceTypeEnabled(config.TraceTypeAsteriscZK) {
    registerAsteriscZK(AsteriscZKGameType, ...)
    // RISC Zero 또는 SP1 prover 사용
}
```

자세한 내용은 [RISC-V와 ZK의 미래](./riscv-and-zk-future-ko.md) 문서를 참조하세요.

---

### Q7: Prestate가 너무 큽니다. 줄일 수 있나요?

**A**: Prestate는 op-program의 초기 VM 상태이므로 크기 고정입니다.

**최적화 방법**:

1. **압축** (이미 적용):
   ```bash
   # .json.gz 형식으로 저장
   gzip prestate-rv64.json
   ```

2. **원격 저장소 사용**:
   ```bash
   --asterisc-prestate-base-url https://cdn.example.com/prestates/
   # 필요 시에만 다운로드
   ```

3. **중복 제거**:
   ```bash
   # 같은 prestate를 여러 게임이 공유
   /data/prestates/0xabcd...json  # 모든 게임이 참조
   ```

---

### Q8: 다른 RISC-V 프로그램을 실행할 수 있나요?

**A**: 이론적으로 가능하지만, 현재는 op-program만 지원됩니다.

**제약 사항**:
- Preimage Oracle 인터페이스 필요
- 특정 시스템 콜 지원
- 결정론적 실행 보장

**미래 가능성**:
```rust
// 임의의 RISC-V 바이너리
asterisc run --input prestate.json -- ./my-riscv-program
```

---

### Q9: Asterisc는 언제 메인넷에 배포되나요?

**A**: Optimism의 로드맵에 따라 다릅니다.

**현재 상태**:
- 테스트넷: 일부 활성화
- 메인넷: 아직 미정

**확인 방법**:
```bash
# DisputeGameFactory에서 GameType 2 확인
cast call $DISPUTE_GAME_FACTORY "gameImpls(uint32)" 2 --rpc-url $L1_RPC
# 0x0000...0000 → 미배포
# 0xabcd...1234 → 배포됨
```

---

### Q10: 성능 벤치마크는 어디서 볼 수 있나요?

**A**: 직접 측정하거나 메트릭을 확인하세요.

**메트릭**:
```bash
# Prometheus 메트릭
curl http://localhost:7300/metrics | grep asterisc

# 주요 메트릭
op_challenger_asterisc_execution_time_seconds
op_challenger_games_asterisc_won
op_challenger_games_asterisc_lost
```

**로컬 벤치마크**:
```bash
time asterisc run \
  --input prestate-rv64.json \
  --proof-at =10000000 \
  --proof-fmt proof.json \
  -- \
  /usr/local/bin/op-program-rv64 --server \
    # ... 설정
```

---

## 참고 자료

### 공식 문서

- [Optimism Fault Proof Specs](https://specs.optimism.io/experimental/fault-proof/)
- [RISC-V Specification](https://riscv.org/technical/specifications/)
- [RISC-V Foundation](https://riscv.org/)

### 관련 파일

- **Asterisc Provider**: `op-challenger/game/fault/trace/asterisc/provider.go`
- **Executor**: `op-challenger/game/fault/trace/asterisc/executor.go`
- **등록 로직**: `op-challenger/game/fault/register.go` (line 187-277)
- **op-program**: `op-program/`

### 커뮤니티

- [Optimism Discord](https://discord.optimism.io)
- [RISC-V Discord](https://discord.gg/riscv)
- [GitHub Issues](https://github.com/tokamak-network/tokamak-thanos/issues)

### 추가 문서

- [게임 타입과 VM 매핑](./game-types-and-vms-ko.md)
- [op-challenger 아키텍처](./op-challenger-architecture-ko.md)
- [배포 가이드](./deployment-guide-ko.md)
- [RISC-V와 ZK의 미래](./riscv-and-zk-future-ko.md)

---

**마지막 업데이트**: 2025-01-21
**버전**: v1.7.7 기준
**상태**: Asterisc GameType 2 지원
