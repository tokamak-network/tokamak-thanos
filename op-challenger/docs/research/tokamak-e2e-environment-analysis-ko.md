# Tokamak Thanos E2E 테스트 환경 분석

> **작성일**: 2025-10-31
> **목적**: E2E 테스트 환경의 구조와 설정을 이해하고, 발생한 sourceHash 에러의 원인을 분석

## 1. 개요

### 1.1 E2E 테스트의 목적

E2E (End-to-End) 테스트는 **전체 Optimism 스택을 메모리에서 실행**하여 실제 환경과 동일한 조건에서 시스템을 검증합니다.

**주요 목표:**
- L1 ↔ L2 상호작용 검증
- Sequencer, Batcher, Proposer, Challenger 통합 테스트
- Fault Proof 시스템 검증
- 실제 배포 전 버그 발견

**장점:**
- 빠른 실행 (실제 블록체인 불필요)
- 격리된 환경 (외부 의존성 없음)
- 재현 가능 (항상 동일한 결과)

**단점:**
- 실제 네트워크 환경과 완전히 동일하지 않음
- 메모리 기반이라 대규모 테스트 어려움

---

## 2. E2E 환경 구성 요소

### 2.0 E2E 환경 설정 위치

E2E 테스트 환경은 다음 파일들에서 설정됩니다:

```
op-e2e/
├── setup.go                      ⭐ 핵심 설정 파일
│   ├── DefaultSystemConfig()     - 기본 시스템 구성 생성
│   ├── SystemConfig.Start()      - 시스템 초기화 및 시작
│   └── System 구조체             - 실행 중인 시스템 상태 관리
│
├── config/
│   ├── init.go                   - L1/L2 allocs 로딩
│   ├── config.go                 - Deploy config
│   └── deployments.go            - L1 배포 주소
│
├── faultproofs/
│   ├── output_cannon_test.go     - Cannon fault proof 테스트
│   ├── output_alphabet_test.go   - Alphabet fault proof 테스트
│   └── helpers.go                - 테스트 헬퍼 (Fault Proof 특화 설정)
│
└── e2eutils/
    ├── geth/                     - Geth 초기화 유틸
    ├── batcher/                  - Batcher 설정
    └── ...                       - 기타 유틸리티
```

**주요 설정 함수:**

| 함수 | 위치 | 용도 |
|------|------|------|
| `DefaultSystemConfig()` | `op-e2e/setup.go` | 기본 E2E 환경 구성 |
| `SystemConfig.Start()` | `op-e2e/setup.go` | 시스템 초기화 및 시작 |
| `StartFaultDisputeSystem()` | `op-e2e/faultproofs/helpers.go` | Fault Proof 테스트 환경 (Tokamak 커스텀) |

### 2.1 전체 아키텍처

```
┌─────────────────────────────────────────────────────────┐
│            Tokamak Thanos E2E Test System                │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────────────────────────────────────┐     │
│  │              L1 Layer (Ethereum)                │     │
│  ├────────────────────────────────────────────────┤     │
│  │  • Geth (in-memory)                             │     │
│  │  • DisputeGameFactory 컨트랙트                  │     │
│  │  • OptimismPortal 컨트랙트                      │     │
│  │  • L2OutputOracle 컨트랙트                      │     │
│  └────────────┬───────────────────────────────────┘     │
│               │                                          │
│               │ (L1 ↔ L2 통신)                          │
│               │                                          │
│  ┌────────────▼───────────────────────────────────┐     │
│  │              L2 Layer (Sequencer)               │     │
│  ├────────────────────────────────────────────────┤     │
│  │  • op-geth (execution)                          │     │
│  │  • op-node (consensus/derivation)               │     │
│  │  • op-batcher (batch submission)                │     │
│  │  • op-proposer (output proposal)                │     │
│  └─────────────────────────────────────────────────┘     │
│                                                          │
│  ┌─────────────────────────────────────────────────┐     │
│  │              L2 Layer (Verifier)                │     │
│  ├─────────────────────────────────────────────────┤     │
│  │  • op-geth (execution)                          │     │
│  │  • op-node (consensus/derivation)               │     │
│  └─────────────────────────────────────────────────┘     │
│                                                          │
│  ⚠️ 주의: op-challenger는 기본 구성에 포함되지 않음!   │
│  Fault Proof 테스트에서만 별도로 추가됨                 │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 2.2 DefaultSystemConfig (기본 구성)

```go
// op-e2e/setup.go
func DefaultSystemConfig(t testing.TB) SystemConfig {
    // Forge allocs 로드 (현재는 사용 안 함, Optimism과의 차이점)
    deployConfig := config.DeployConfig.Copy()
    l1Deployments := config.L1Deployments.Copy()

    return SystemConfig{
        Secrets:       secrets,
        Premine:       premine,
        DeployConfig:  deployConfig,
        L1Deployments: config.L1Deployments,
        BlobsPath:     t.TempDir(),

        // 노드 구성: Sequencer + Verifier
        Nodes: map[string]*rollupNode.Config{
            "sequencer": {
                Driver: driver.Config{
                    VerifierConfDepth:  0,
                    SequencerConfDepth: 0,
                    SequencerEnabled:   true,  // 시퀀서 활성화
                },
            },
            "verifier": {
                Driver: driver.Config{
                    VerifierConfDepth:  0,
                    SequencerConfDepth: 0,
                    SequencerEnabled:   false,  // 검증만
                },
            },
        },

        // 로깅 설정
        Loggers: map[string]log.Logger{
            "sequencer": testlog.Logger(t, log.LevelInfo),
            "verifier":  testlog.Logger(t, log.LevelInfo),
            "batcher":   testlog.Logger(t, log.LevelInfo),
            "proposer":  testlog.Logger(t, log.LevelCrit),
        },

        // ...
    }
}
```

**노드 구성:**
- **L1 노드**: 1개 (Geth)
- **L2 Sequencer**: op-geth + op-node + op-batcher + op-proposer
- **L2 Verifier**: op-geth + op-node (읽기 전용)
- **총 3개 Geth 인스턴스**

### 2.3 Fault Proof 테스트 구성 (`StartFaultDisputeSystem`)

```go
func StartFaultDisputeSystem(t *testing.T) (*System, *ethclient.Client) {
    cfg := op_e2e.DefaultSystemConfig(t)

    // Verifier 제거 (테스트 단순화)
    delete(cfg.Nodes, "verifier")

    // Sequencer 설정
    cfg.Nodes["sequencer"].SafeDBPath = t.TempDir()
    cfg.DeployConfig.SequencerWindowSize = 4
    cfg.DeployConfig.FinalizationPeriodSeconds = 2
    cfg.SupportL1TimeTravel = true  // L1 시간 조작 가능
    cfg.DeployConfig.L2OutputOracleSubmissionInterval = 1
    cfg.NonFinalizedProposals = true

    sys, err := cfg.Start(t)
    require.Nil(t, err)
    return sys, sys.Clients["l1"]
}
```

**간소화된 구성:**
- **L1 노드**: 1개
- **L2 Sequencer**: 1개
- **총 2개 Geth 인스턴스** ← 예상

---

## 3. Genesis 생성 프로세스

### 3.1 L2 Genesis 생성 흐름

```
┌─────────────────────────────────────────────────────┐
│ 1. DefaultSystemConfig 생성                          │
│    - DeployConfig 로드                               │
│    - L1Deployments 로드                              │
│    - Secrets (키) 생성                               │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ 2. cfg.Start(t) 호출                                 │
│    ├─ BuildL1DeveloperGenesis(...)                  │
│    │   └─ L1 genesis block 생성                     │
│    │                                                  │
│    └─ BuildL2Genesis(config, l1Block)  ⚠️ 문제!    │
│        │                                              │
│        ├─ ❌ Tokamak: Forge allocs 안 받음          │
│        │   └─ NewL2StorageConfig() 호출              │
│        │       └─ L1Block storage 설정 (문제!)      │
│        │                                              │
│        └─ ✅ Optimism: Forge allocs 받음            │
│            └─ allocs 그대로 사용                     │
│                └─ L1Block storage 비어있음          │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│ 3. Geth 초기화                                       │
│    - InitL1(..., l1Genesis)                          │
│    - InitL2(..., l2Genesis)                          │
│    - 각 노드 데이터베이스에 genesis 저장             │
└─────────────────────────────────────────────────────┘
```

**⚠️ 핵심 차이점:**

| 항목 | Tokamak (현재) | Optimism | 결과 |
|------|----------------|----------|------|
| **함수 시그니처** | `BuildL2Genesis(config, l1Block)` | `BuildL2Genesis(config, l2Allocs, l1Block)` | 다름 ❌ |
| **Forge allocs 사용** | ❌ 안 받음 | ✅ 받음 | **문제의 핵심** |
| **L1Block storage** | Go 코드로 설정 | Forge allocs (비어있음) | sourceHash 에러 발생 |
| **Devnet** | Forge allocs 사용 (정상) | Forge allocs 사용 (정상) | 동일 ✅ |
| **E2E 테스트** | Go 코드 사용 (에러) | Forge allocs 사용 (정상) | **불일치가 문제!** |

### 3.2 L1Block Storage 설정 (문제 발생 지점)

**기존 Tokamak 코드** (`op-chain-ops/genesis/config.go`):

```go
func NewL2StorageConfig(config *DeployConfig, block *types.Block) (state.StorageConfig, error) {
    storage := make(state.StorageConfig)

    // ❌ 문제: L1Block storage를 genesis에서 설정
    storage["L1Block"] = state.StorageValues{
        "number":         block.Number(),
        "timestamp":      block.Time(),
        "basefee":        block.BaseFee(),
        "hash":           block.Hash(),
        "sequenceNumber": 0,
        "batcherHash":    eth.AddressAsLeftPaddedHash(config.BatchSenderAddress),
        "l1FeeOverhead":  config.GasPriceOracleOverhead,
        "l1FeeScalar":    config.GasPriceOracleScalar,
    }

    // ... 다른 predeploy storage 설정
    return storage, nil
}
```

**Optimism 방식** (Forge 스크립트):

```solidity
// L1Block storage는 비워둠
// 블록 1의 deposit transaction이 초기화
```

---

## 4. 블록 생성 및 Deposit Transaction

### 4.1 블록 1 생성 프로세스

```
Genesis (블록 0)
  │
  │ L1Block storage = 설정됨 (Tokamak) ❌
  │ L1Block storage = 비어있음 (Optimism) ✅
  │
  ▼
Sequencer 시작
  │
  ▼
블록 1 생성 시도
  ├─ L1 attributes deposit transaction 생성
  │  ├─ SourceHash: 계산됨 ✅
  │  ├─ From: L1Block Address
  │  ├─ To: L1Block Address
  │  ├─ Data: L1 정보 (number, timestamp, basefee 등)
  │  └─ IsSystemTransaction: true
  │
  ├─ Payload 생성 (블록 1)
  │  └─ Deposit tx 포함
  │
  ├─ 블록 실행
  │  └─ L1Block 컨트랙트 업데이트
  │     ├─ Tokamak: 이미 설정됨 → 업데이트 ❌
  │     └─ Optimism: 비어있음 → 초기화 ✅
  │
  └─ 블록 저장 (DB)
```

### 4.2 Deposit Transaction 구조

```go
type DepositTx struct {
    SourceHash          common.Hash    // ✅ 항상 설정됨
    From                common.Address
    To                  *common.Address
    Mint                *big.Int
    Value               *big.Int
    Gas                 uint64
    IsSystemTransaction bool
    Data                []byte
}
```

**SourceHash 생성:**
```go
// L1 attributes deposit
source := L1InfoDepositSource{
    L1BlockHash: block.Hash(),
    SeqNumber:   seqNumber,
}
dep.SourceHash = source.SourceHash()  // keccak256(l1BlockHash, seqNumber)
```

---

## 5. 발생한 문제: sourceHash 에러

### 5.1 에러 메시지

```
ERROR: sequencer temporarily failed to start building new block
err="temp: failed to retrieve L2 parent block:
     failed to determine block-hash of hash 0xc09d2b...,
     could not get payload: missing required field 'sourceHash' in transaction"
```

### 5.2 에러 발생 시나리오

```
Timeline:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

15:16:37.062 ✅ Sequencer 블록 1 생성 시작
             payload_id=0x0237c2d922881d89
             parent=2da9d2..ab6d94:0 (genesis)

15:16:37.064 ✅ Payload 업데이트 완료
             hash=c09d2b..852d5b
             txs=1 (deposit tx with sourceHash)

15:16:37.069 ✅ 블록 1 임포트 성공
             number=1 hash=c09d2b..852d5b

15:16:37.073 ✅ 블록 1 삽입 완료
             (DB에 저장됨)

15:16:37.074 ✅ "successfully built a new block"

15:16:37.075 ℹ️  블록 2 생성 시작
             parent=c09d2b..852d5b:1

15:16:37.078 ❌ ERROR!
             "failed to retrieve L2 parent block"
             "missing required field 'sourceHash'"
```

**문제 분석:**
1. 블록 1을 **만들 때**: sourceHash 있음 ✅
2. 블록 1을 **저장할 때**: 정상 ✅
3. 블록 1을 **읽을 때**: sourceHash 없음 ❌

### 5.3 가능한 원인들

#### 원인 1: JSON Serialization 문제 (가능성 낮음)

```go
// op-geth의 DepositTx 구조체에서
type DepositTx struct {
    SourceHash common.Hash `json:"sourceHash"`  // ✅ 정상
    // 또는
    SourceHash common.Hash `json:"-"`  // ❌ JSON에서 제외
}
```

**왜 가능성 낮은가:**
- Devnet은 정상 작동 (같은 op-geth 사용)
- 다른 많은 테스트도 통과

#### 원인 2: 데이터베이스 불일치 (가능성 높음) ⚠️

```
가설: 블록을 저장한 DB ≠ 읽으려는 DB

시나리오:
1. Sequencer A가 블록 1을 DB_A에 저장
2. Sequencer B가 블록 2를 만들려고 DB_B에서 블록 1 읽기 시도
3. DB_B에는 블록 1이 없거나 오래된 버전
4. 에러 발생
```

**증거:**
- 로그에서 4개의 geth 인스턴스 초기화 확인
- 블록 1이 여러 번 import됨 (다른 hash):
  - `572600..7d45b9`
  - `f263e8..669fa8`
  - `c09d2b..852d5b` ← 에러 발생

#### 원인 3: 이전 테스트 데이터 (가능성 중간)

```
가설: t.TempDir() 정리 실패로 이전 데이터 남음

증거:
- 테스트가 15:16에 실행
- Genesis 코드 수정은 14:52
- 그 사이에 다른 테스트 실행 가능성
```

#### 원인 4: Genesis L1Block Storage 설정 (원래 가설)

```
가설: L1Block이 genesis에 설정되어 deposit tx 실패

문제점:
- 하지만 블록 1은 성공적으로 생성됨
- 읽을 때 문제 발생
- Genesis 문제라면 생성 단계에서 실패해야 함
```

---

## 6. 로그 분석

### 6.1 Geth 인스턴스 초기화

```
15:16:31.797 Initialising Ethereum protocol network=0   dbversion=<nil>
15:16:31.852 Initialising Ethereum protocol network=0   dbversion=<nil>
15:16:31.906 Initialising Ethereum protocol network=900 dbversion=<nil>
15:16:32.360 Initialising Ethereum protocol network=901 dbversion=<nil>
```

**분석:**
- `network=0`: L1 노드 (2개?) 🤔
- `network=900`: L2 노드 (ChainID?)
- `network=901`: L2 노드 (또 다른 ChainID?)

**예상 구성과 불일치:**
- 예상: L1(1개) + L2 Sequencer(1개) = 2개
- 실제: 4개

### 6.2 블록 1 Import 이벤트

```
15:16:31.811 Imported number=1 hash=572600..7d45b9  ← L2 체인 A
15:16:31.859 Imported number=1 hash=f263e8..669fa8  ← L2 체인 B
15:16:33.047 Imported number=1 hash=d2cea0..a4d0ce  ← L1
15:16:37.069 Imported number=1 hash=c09d2b..852d5b  ← L2 체인 C (에러!)
```

**결론:**
- **최소 3개의 다른 블록체인이 실행 중**
- 마지막 것에서 sourceHash 에러 발생

---

## 7. 해결 방향

### 7.1 단기 해결책 (검증 필요)

#### 1. 완전히 새로운 환경에서 테스트

```bash
# 모든 프로세스 종료
ps aux | grep -E "geth|op-node|go test" | grep -v grep
kill -9 <PID>

# 임시 디렉토리 정리
rm -rf /tmp/TestMultiple*

# Go 캐시 정리
go clean -testcache -cache

# 새로운 터미널에서 실행
cd /Users/zena/tokamak-projects/tokamak-thanos/op-e2e
go test -v -run TestMultipleGameTypes -timeout 10m ./faultproofs
```

#### 2. 단일 테스트 실행

```bash
# 가장 간단한 테스트부터
go test -v -run TestSystemE2E_Simple -timeout 3m
```

#### 3. Genesis 수정 적용 확인

```go
// op-chain-ops/genesis/config.go에서
// L1Block storage 설정 제거 확인
storage["L1Block"] = state.StorageValues{...}  // 이 부분 있으면 ❌
```

### 7.2 중기 해결책 (근본 해결)

#### 0. 코어 코드 변경 영향 범위 분석

**⚠️ 중요: `BuildL2Genesis` 함수는 코어 함수입니다!**

**현재 사용 위치:**
```bash
# 사용하는 곳 찾기
$ grep -r "BuildL2Genesis(" --include="*.go"

op-node/cmd/genesis/cmd.go          ← 실제 배포 CLI 명령어
op-e2e/setup.go                     ← E2E 테스트
op-e2e/op_geth.go                   ← E2E 테스트
op-e2e/e2eutils/setup.go            ← E2E 테스트
```

**실제 배포에서의 사용:**
```bash
# packages/tokamak/contracts-bedrock/scripts/start-deploy.sh
op-node genesis l2 \
  --deploy-config $DEPLOY_CONFIG_PATH \
  --l1-deployments $deployResultFile \
  --outfile.l2 $outdir/genesis.json

# 이것도 BuildL2Genesis(config, l1Block) 호출!
# ⚠️ 실제 배포도 같은 문제 가능성!
```

**영향도 평가:**

| 위치 | 사용 여부 | 실제 에러 | 영향도 | 비고 |
|------|-----------|-----------|--------|------|
| **Devnet** | ❌ 사용 안 함 | ❌ 없음 | 없음 | `make devnet-allocs` (Forge 스크립트 사용) |
| **E2E 테스트** | ✅ 사용 중 | ✅ **에러 발생** | 높음 | sourceHash 에러 발생 중 |
| **실제 배포** | ✅ 사용 중 | ❌ **없음** | 낮음 | ✅ 정상 작동 확인됨 |
| **테스트넷 배포** | ✅ 사용 중 | ❌ 없음 (추정) | 낮음 | 실제 배포와 동일 |

**⚠️ 중요한 발견: 실제 배포는 정상 작동!**

**왜 실제 배포는 괜찮고 E2E 테스트만 에러인가?**

**핵심 차이: Genesis 생성 방식!**

```
실제 배포 (deploy-modular.sh) ✅:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
make devnet-up
  ├─ make devnet-allocs  ← Forge 스크립트!
  │  └─ .devnet/allocs-l1.json 생성
  │     └─ L1Block storage 비어있음 ✅
  │
  └─ bedrock-devnet/main.py
     └─ Forge allocs 사용
        └─ L1Block storage 비어있음 유지 ✅

E2E 테스트 (op-e2e/setup.go) ❌:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
BuildL2Genesis(config, l1Block)
  └─ NewL2StorageConfig()  ← Go 코드!
     └─ L1Block storage 설정됨 ❌
        └─ sourceHash 에러 발생!
```

**결론:**
- **실제 배포: Forge allocs 사용** → 정상 ✅
- **E2E 테스트: Go 코드 사용** → 에러 ❌
- **Devnet: Forge allocs 사용** → 정상 ✅
- **문제는 E2E 테스트만!** ⚠️

**왜 이런 차이가 생겼나?**
- 실제 배포와 Devnet은 **Forge 스크립트 기반** (`make devnet-allocs`)
- E2E 테스트는 **Go 코드 기반** (`BuildL2Genesis`)
- Optimism도 E2E 테스트에서 **Forge allocs 사용**
- Tokamak만 E2E에서 **Go 코드 사용** ← 불일치!

**해결 방향:**
- E2E 테스트도 Forge allocs를 사용하도록 변경 (Optimism과 동일하게)
- 모든 환경(Devnet/E2E/Production)에서 **동일한 Genesis 생성 방식** 사용 ✅

---

## 🚧 수정 작업 진행 중 (TODO)

### ✅ Phase 1: BuildL2Genesis 함수 수정 (완료)

**목표**: Forge allocs를 선택적으로 받을 수 있도록 함수 시그니처 변경

**파일**: `op-chain-ops/genesis/layer_two.go`

- [x] **BuildL2Genesis 함수 시그니처 변경**
  ```go
  // Before
  func BuildL2Genesis(config *DeployConfig, l1StartBlock *types.Block) (*core.Genesis, error)

  // After (variadic parameter for backward compatibility)
  func BuildL2Genesis(config *DeployConfig, l1StartBlock *types.Block, dump ...*ForgeAllocs) (*core.Genesis, error)
  ```

- [x] **Forge allocs 사용 로직 추가**
  - Forge allocs가 제공되면 직접 사용
  - 제공되지 않으면 기존 Go 코드 사용 (하위 호환성)

- [x] **로깅 추가**
  - Optimism 방식 사용 시: "Using Forge allocs for L2 genesis (Optimism style)"
  - Legacy 방식 사용 시: "Using Go code for L2 genesis (legacy Tokamak style)"

**장점**:
- ✅ 기존 코드 영향 없음 (variadic parameter로 하위 호환성 유지)
- ✅ 점진적 마이그레이션 가능

---

### ✅ Phase 2: E2E 테스트 setup.go 수정 (완료)

**목표**: E2E 테스트에서 Forge allocs 사용

**파일**: `op-e2e/setup.go`

- [x] **1. L2 allocs 로딩 추가**
  ```go
  // op-e2e/setup.go - Start() 함수 수정

  // After
  l2Allocs := config.L2Allocs(genesis.L2AllocsFjord)
  l2Genesis, err := genesis.BuildL2Genesis(cfg.DeployConfig, l1Block, l2Allocs)
  ```

- [x] **2. 에러 핸들링 개선**
  ```go
  if err != nil {
      return nil, fmt.Errorf("failed to build L2 genesis: %w", err)
  }
  ```

- [x] **3. 주석 추가**
  - "Use Forge allocs for L2 genesis (Optimism style)"
  - "This ensures L1Block storage is empty and initialized by deposit tx"

**예상 결과**:
- ✅ E2E 테스트에서 Forge allocs 사용
- ✅ L1Block storage 비어있음 유지
- ✅ sourceHash 에러 해결

---

### ✅ Phase 3: 다른 호출 위치 수정 (완료)

**파일**: `op-e2e/op_geth.go`, `op-e2e/e2eutils/setup.go`

- [x] **1. op-e2e/op_geth.go 수정**
  - Forge allocs 사용 추가
  - `config.L2Allocs(genesis.L2AllocsFjord)` 사용

- [x] **2. op-e2e/e2eutils/setup.go 수정**
  - Forge allocs 사용 추가
  - `alloc.L2Allocs(genesis.L2AllocsFjord)` 사용

- [x] **3. 모든 E2E 관련 파일 수정 완료**
  - `op-e2e/setup.go` ✅
  - `op-e2e/op_geth.go` ✅
  - `op-e2e/e2eutils/setup.go` ✅

---

### ✅ Phase 4: 테스트 및 검증 (완료)

- [x] **1. 단위 테스트 실행**
  ```bash
  cd op-challenger
  go test -short ./...
  ```
  - ✅ 모든 테스트 통과!

- [x] **2. E2E 테스트 환경 설정**
  ```bash
  cd op-e2e
  ls -lh .devnet/allocs-l2.json
  ```
  - ✅ Forge allocs 파일 존재 확인 (8.9MB)

- [x] **3. E2E 테스트 실행**
  ```bash
  cd op-e2e
  go test -v -run TestMultipleGameTypes -timeout 10m ./faultproofs
  ```
  - ✅ **sourceHash 에러 완전히 해결!**
  - ✅ **"Using Forge allocs for L2 genesis (Optimism style)" 로그 확인**

- [x] **4. 로그 분석**
  - ✅ Genesis: "L2 genesis built from Forge allocs accounts=2367"
  - ✅ Geth 정상 초기화: "Initialising Ethereum protocol network=900"
  - ✅ **sourceHash 에러 0건** (완전 해결!)
  - ✅ L1Block storage 비어있음 (Forge allocs 그대로 사용)

**🎉 핵심 성과:**
- E2E 테스트가 이제 Optimism과 동일한 방식으로 작동
- L1Block storage가 비어있어 deposit tx로 초기화
- sourceHash 에러 완전히 해결
- Devnet, 실제 배포, E2E 테스트 모두 **동일한 Genesis 생성 방식** 사용 ✅

---

### ⏳ Phase 5: 문서 업데이트 및 정리

**TODO**:

- [ ] **1. testing-guide-ko.md 업데이트**
  - sourceHash 에러 섹션 수정
  - 해결 완료 표시

- [ ] **2. 이 TODO 섹션 제거**
  - 모든 작업 완료 후 제거

- [ ] **3. 결론 섹션 업데이트**
  - "E2E 테스트도 Forge allocs 사용으로 변경 완료" 추가

---

## 💡 작업 가이드

### 빌드 및 테스트 명령어

```bash
# 1. 수정 후 빌드 확인
cd /Users/zena/tokamak-projects/tokamak-thanos
make build-go

# 2. 단위 테스트
cd op-challenger
go test -short ./...

# 3. E2E 테스트 준비
cd op-e2e
make devnet-allocs

# 4. E2E 테스트 실행
go test -v -run TestMultipleGameTypes -timeout 10m ./faultproofs 2>&1 | tee e2e-test.log

# 5. 로그 확인
grep -i "forge allocs" e2e-test.log
grep -i "sourcehash" e2e-test.log
```

### 롤백 방법

만약 문제 발생 시:

```bash
# Git으로 롤백
cd /Users/zena/tokamak-projects/tokamak-thanos
git diff op-chain-ops/genesis/layer_two.go  # 변경사항 확인
git checkout op-chain-ops/genesis/layer_two.go  # 롤백
```

---

#### 1. 안전한 방법: 기존 함수 유지 + 새 함수 추가 ⭐ (대체 방안)

**Step 1: 새로운 함수 추가 (Breaking Change 없음)**

```go
// op-chain-ops/genesis/layer_two.go

// 기존 함수 유지 (Deprecated 표시)
// Deprecated: Use BuildL2GenesisFromAllocs instead
func BuildL2Genesis(config *DeployConfig, l1StartBlock *types.Block) (*core.Genesis, error) {
    // 기존 로직 유지 (하위 호환성)
    genspec, err := NewL2Genesis(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    db := state.NewMemoryStateDB(genspec)
    storage, err := NewL2StorageConfig(config, l1StartBlock)
    if err != nil {
        return nil, err
    }
    // ... 기존 로직
}

// ✅ 새로운 함수 추가 (Optimism 방식)
func BuildL2GenesisFromAllocs(
    config *DeployConfig,
    dump *ForgeAllocs,
    l1StartBlock *types.Block,
) (*core.Genesis, error) {
    genspec, err := NewL2Genesis(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    // Forge allocs 사용
    for addr, val := range dump.Copy().Accounts {
        genspec.Alloc[addr] = val  // ← L1Block storage 비어있음 유지
    }

    return genspec, nil
}
```

**Step 2: E2E 테스트만 새 함수로 변경**

```go
// op-e2e/setup.go
func (cfg SystemConfig) Start(t *testing.T) (*System, error) {
    // ...

    // ✅ 새 함수 사용
    l2Allocs := config.L2Allocs(genesis.L2AllocsFjord)
    l2Genesis, err := genesis.BuildL2GenesisFromAllocs(
        cfg.DeployConfig,
        l2Allocs,
        l1Block,
    )
    if err != nil {
        return nil, err
    }

    // ...
}
```

**장점:**
- ✅ 기존 코드 영향 없음 (하위 호환성 유지)
- ✅ E2E 테스트 문제 해결
- ✅ 실제 배포 스크립트 영향 없음
- ✅ 점진적 마이그레이션 가능

**단점:**
- ⚠️ 함수 중복 (향후 정리 필요)

#### 2. 적극적 방법: Optimism처럼 함수 시그니처 변경 (권장하지만 신중히)

**현재 Tokamak 코드:**
```go
// op-e2e/setup.go
func (cfg SystemConfig) Start(t *testing.T) (*System, error) {
    // ...

    // ❌ 문제: Forge allocs를 안 받음
    l2Genesis, err := genesis.BuildL2Genesis(cfg.DeployConfig, l1Block)
    if err != nil {
        return nil, err
    }

    // ...
}
```

**Optimism 방식으로 수정:**
```go
// op-e2e/setup.go
func (cfg SystemConfig) Start(t *testing.T) (*System, error) {
    // ...

    // ✅ 해결: Forge allocs 사용
    l2Allocs := config.L2Allocs(genesis.L2AllocsFjord)
    l2Genesis, err := genesis.BuildL2Genesis(
        cfg.DeployConfig,
        l2Allocs,      // ← Forge allocs 추가!
        l1Block,
    )
    if err != nil {
        return nil, err
    }

    // ...
}
```

**op-chain-ops/genesis/layer_two.go 수정:**
```go
// 현재 시그니처
func BuildL2Genesis(config *DeployConfig, l1StartBlock *types.Block) (*core.Genesis, error)

// 변경 후 시그니처 (Optimism과 동일)
func BuildL2Genesis(
    config *DeployConfig,
    dump *ForgeAllocs,        // ← 추가!
    l1StartBlock *types.Block,
) (*core.Genesis, error) {
    genspec, err := NewL2Genesis(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    // Forge allocs 사용
    if dump != nil {
        for addr, val := range dump.Copy().Accounts {
            genspec.Alloc[addr] = val  // ← L1Block storage 비어있음 유지
        }
        return genspec, nil
    }

    // Fallback: 기존 방식 (deprecated, 향후 제거)
    db := state.NewMemoryStateDB(genspec)
    storage, err := NewL2StorageConfig(config, l1StartBlock)
    if err != nil {
        return nil, err
    }
    // ...
}
```

#### 2. E2E 테스트 환경 재구성

```go
// 명확한 노드 구성 확인
func StartFaultDisputeSystem(t *testing.T) (*System, *ethclient.Client) {
    cfg := op_e2e.DefaultSystemConfig(t)

    // 명확하게 verifier 제거
    delete(cfg.Nodes, "verifier")

    // 각 노드의 datadir 명확히 분리
    cfg.Nodes["sequencer"].SafeDBPath = filepath.Join(t.TempDir(), "sequencer-db")

    // 로깅 강화
    cfg.Loggers["sequencer"] = testlog.Logger(t, log.LevelDebug)

    sys, err := cfg.Start(t)
    require.Nil(t, err)

    // 시스템 상태 검증
    require.Len(t, sys.EthInstances, 2, "Should have exactly 2 eth instances (L1 + Sequencer)")

    return sys, sys.Clients["l1"]
}
```

#### 2. Genesis 코드를 Optimism 방식으로 완전히 변경

```go
// op-chain-ops/genesis/config.go
func NewL2StorageConfig(config *DeployConfig, block *types.Block) (state.StorageConfig, error) {
    storage := make(state.StorageConfig)

    // L1Block storage는 비워둠 (Optimism 방식)
    // 첫 번째 deposit transaction이 초기화

    storage["L2CrossDomainMessenger"] = state.StorageValues{...}
    storage["L2StandardBridge"] = state.StorageValues{...}
    // ... 다른 predeploys

    return storage, nil
}
```

### 7.3 장기 해결책

#### 1. Forge 스크립트로 Genesis 생성 통일

```bash
# Devnet과 E2E 테스트 모두 동일한 방식 사용
make devnet-allocs  # Forge 스크립트 사용

# E2E 테스트도 생성된 allocs 사용
op-e2e/config/init.go:
- L2Allocs 사용 (Forge 생성)
- NewL2StorageConfig 사용 안 함
```

#### 2. 통합 테스트 강화

```go
// genesis 생성 후 즉시 검증
func TestL2Genesis(t *testing.T) {
    genesis, err := BuildL2Genesis(config, l1Block)
    require.NoError(t, err)

    // L1Block storage 확인
    l1BlockStorage := genesis.Alloc[predeploys.L1BlockAddr].Storage

    // Optimism 방식: 비어있어야 함
    require.Empty(t, l1BlockStorage, "L1Block storage should be empty in genesis")

    // 블록 1 시뮬레이션
    block1 := simulateBlock1(genesis)

    // 블록 1의 deposit tx에 sourceHash 확인
    require.NotEmpty(t, block1.Transactions()[0].SourceHash())
}
```

---

## 8. Optimism과의 차이점

### 8.1 Genesis 생성 방식

| 항목 | Tokamak Thanos | Optimism |
|------|----------------|----------|
| **방식** | Go 코드 (`NewL2StorageConfig`) | Forge 스크립트 (`L2Genesis.s.sol`) |
| **L1Block Storage** | Genesis에서 설정 ❌ | 비어있음 ✅ |
| **초기화 시점** | Genesis 블록 0 | 블록 1 (deposit tx) |
| **Devnet** | Forge 스크립트 (비어있음) ✅ | Forge 스크립트 (비어있음) ✅ |
| **E2E 테스트** | Go 코드 (설정됨) ❌ | Forge allocs 사용 (비어있음) ✅ |

### 8.2 E2E 테스트 구성

| 항목 | Tokamak Thanos | Optimism |
|------|----------------|----------|
| **Genesis 소스** | Go 코드로 생성 | Forge allocs 파일 로드 |
| **일관성** | Devnet ≠ E2E | Devnet = E2E |
| **문제** | sourceHash 에러 | 없음 |

---

## 9. 다음 단계

### 9.1 즉시 확인할 사항

1. ✅ **Genesis 코드 수정 확인**
   - `op-chain-ops/genesis/config.go`에서 L1Block storage 설정 제거

2. ⏳ **빌드 확인**
   ```bash
   cd /Users/zena/tokamak-projects/tokamak-thanos
   make build-go
   ```

3. ⏳ **테스트 환경 완전 정리**
   ```bash
   # 모든 관련 프로세스 종료
   # 캐시 정리
   # 새 터미널에서 테스트
   ```

4. ⏳ **단순한 테스트부터 시작**
   ```bash
   cd op-e2e
   go test -v -run TestSystemE2E_Simple -timeout 3m
   ```

### 9.2 근본 원인 파악

1. **4개의 Geth 인스턴스 원인 규명**
   - 왜 예상(2개)과 다른가?
   - 이전 테스트의 잔여물인가?
   - 테스트 코드 문제인가?

2. **데이터베이스 경로 추적**
   - 각 geth 인스턴스의 datadir 로깅
   - 블록 저장/읽기 경로 일치 확인

3. **SourceHash 직렬화 확인**
   - op-geth의 DepositTx JSON 태그 확인
   - 저장/읽기 과정 로깅

---

## 10. 결론

### 10.1 핵심 발견

1. **⭐ 실제 배포는 정상, E2E 테스트만 에러 발생**
   - Devnet: 정상 작동 ✅ (Forge 스크립트 사용)
   - 실제 배포: 정상 작동 ✅ (BuildL2Genesis 사용하지만 문제 없음)
   - E2E 테스트: sourceHash 에러 ❌ (BuildL2Genesis 사용)

2. **Tokamak vs Optimism Genesis 생성 방식 차이**
   - Tokamak E2E: Go 코드 (`NewL2StorageConfig`) → L1Block storage 설정됨
   - Optimism E2E: Forge allocs 사용 → L1Block storage 비어있음
   - Tokamak Devnet: Forge allocs 사용 → L1Block storage 비어있음 (Optimism과 동일)

3. **E2E 테스트 환경 특유의 문제**
   - 블록 1 생성: 성공 ✅
   - 블록 1 저장: 성공 ✅
   - 블록 1 읽기: sourceHash 없음 ❌
   - 원인: 데이터베이스 불일치, 타이밍 이슈, 메모리 기반 환경 등 추정

4. **예상과 다른 geth 인스턴스 개수**
   - 예상: 2개 (L1 + Sequencer)
   - 실제 로그: 4개 초기화 확인
   - 원인: 이전 테스트 잔여물 또는 테스트 설정 문제 가능성

### 10.2 권장 사항

**단기 (테스트 통과):**
- Genesis 코드 수정 (L1Block storage 제거)
- 환경 완전 정리 후 재테스트

**중기 (안정성):**
- E2E 테스트를 Forge allocs 기반으로 변경
- Devnet과 E2E의 genesis 생성 방식 통일

**장기 (유지보수성):**
- Optimism upstream과 완전 동기화
- Genesis 생성 로직 단순화
- 통합 테스트 강화

### 10.3 여전히 답이 필요한 질문

1. ❓ **왜 4개의 geth 인스턴스가 실행되는가?**
2. ❓ **블록 1은 성공적으로 생성되었는데, 왜 읽을 때 sourceHash가 없는가?**
3. ❓ **데이터베이스 불일치가 실제 원인인가?**
4. ❓ **Genesis L1Block storage가 실제 sourceHash 에러의 원인인가?**

이 질문들에 대한 답을 찾기 위해서는 **추가 디버깅과 로깅**이 필요합니다.

---

## 참고 자료

- [Optimism Specs - Deposits](https://specs.optimism.io/protocol/deposits.html)
- [op-geth DepositTx Implementation](https://github.com/ethereum-optimism/op-geth)
- [Tokamak Genesis Generation Code](../../op-chain-ops/genesis/)
- [E2E Test Setup Code](../../op-e2e/setup.go)

