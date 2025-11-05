# Optimism E2E 테스트 환경 분석

> **작성일**: 2025-10-31
> **목적**: Optimism upstream의 E2E 테스트 환경을 분석하여 Tokamak과 비교

## 1. 개요

### 1.1 Optimism E2E 테스트의 특징

Optimism의 E2E 테스트는 **Forge allocs 파일 기반**으로 genesis를 생성하여, **실제 배포 환경과 동일한 초기 상태**를 보장합니다.

**핵심 차이점:**
- Genesis 생성: Solidity Forge 스크립트
- 일관성: Devnet = E2E = Production
- L1Block 초기화: 블록 1의 deposit transaction

---

## 2. Genesis 생성 방식

### 2.1 Forge 스크립트 기반 Genesis

```
┌──────────────────────────────────────────────────┐
│ Forge Script (Solidity)                           │
├──────────────────────────────────────────────────┤
│                                                   │
│  scripts/L2Genesis.s.sol:L2Genesis                │
│  ├─ Deploy all predeploy contracts                │
│  ├─ Initialize storage (except L1Block)           │
│  ├─ Generate state dump                           │
│  └─ Output: state-dump-{chainId}.json            │
│                                                    │
└────────────────┬─────────────────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────────────────┐
│ Allocs File                                       │
├──────────────────────────────────────────────────┤
│  .devnet/allocs-l2.json                           │
│  ├─ All predeploy contracts                       │
│  ├─ Contract bytecode                             │
│  ├─ Storage values                                │
│  └─ L1Block: storage = {} (비어있음!) ✅         │
└────────────────┬─────────────────────────────────┘
                 │
                 ▼
┌──────────────────────────────────────────────────┐
│ Usage                                             │
├──────────────────────────────────────────────────┤
│  - Devnet: Load allocs directly                   │
│  - E2E Tests: Load allocs directly                │
│  - Production: Load allocs directly               │
│                                                    │
│  → 모두 동일한 genesis 사용 ✅                    │
└───────────────────────────────────────────────────┘
```

### 2.2 Forge 스크립트 예시

```solidity
// scripts/L2Genesis.s.sol
contract L2Genesis is Script {
    function runWithAllUpgrades() public {
        // Deploy L1Block predeploy (proxy)
        address l1Block = 0x4200000000000000000000000000000000000015;

        // L1Block의 storage는 설정하지 않음!
        // 블록 1의 deposit transaction이 초기화할 것임

        // 다른 predeploys 설정
        setupL2CrossDomainMessenger();
        setupL2StandardBridge();
        // ...

        // State dump 생성
        vm.dumpState("state-dump.json");
    }
}
```

### 2.3 Allocs 파일 구조

```json
{
  "0x4200000000000000000000000000000000000015": {
    "balance": "0x0",
    "code": "0x60806040...",  // Proxy bytecode
    "nonce": "0x0",
    "storage": {
      // Implementation slot만 설정
      "0x360894a13ba1...": "0x...c0d30015",
      // L1 정보 관련 storage는 비어있음!
    }
  }
}
```

---

## 3. E2E 테스트 Genesis 로딩

### 3.1 Allocs 로딩 프로세스

```go
// op-e2e/config/init.go
func init() {
    // Forge allocs 파일 로드
    l2Allocs = make(map[genesis.L2AllocsMode]*genesis.ForgeAllocs)

    mustL2Allocs := func(mode genesis.L2AllocsMode) {
        name := "allocs-l2"
        if mode != "" {
            name += "-" + string(mode)
        }
        allocs, err := genesis.LoadForgeAllocs(
            filepath.Join(l2AllocsDir, name+".json")
        )
        if err != nil {
            panic(err)
        }
        l2Allocs[mode] = allocs
    }

    // 여러 fork 버전의 allocs 로드
    mustL2Allocs(genesis.L2AllocsFjord)
    mustL2Allocs(genesis.L2AllocsEcotone)
    mustL2Allocs(genesis.L2AllocsDelta)
}
```

### 3.2 BuildL2Genesis (Optimism 방식)

```go
// op-chain-ops/genesis/layer_two.go
func BuildL2Genesis(
    config *DeployConfig,
    dump *foundry.ForgeAllocs,  // ← Forge allocs 사용!
    l1StartBlock *eth.BlockRef,
) (*core.Genesis, error) {
    genspec, err := NewL2Genesis(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    // Forge allocs를 직접 복사
    for addr, val := range dump.Copy().Accounts {
        genspec.Alloc[addr] = val  // ← L1Block storage 비어있음 유지
    }

    // L1Block storage를 Go 코드로 설정하지 않음!
    // NewL2StorageConfig() 호출 없음!

    return genspec, nil
}
```

**핵심 차이:**
- ❌ `NewL2StorageConfig(config, l1Block)` 호출 없음
- ✅ Forge allocs를 그대로 사용
- ✅ L1Block storage 비어있음 유지

---

## 4. 블록 1 초기화

### 4.1 L1 Attributes Deposit Transaction

```
블록 0 (Genesis)
  │
  │ L1Block storage = {} (비어있음) ✅
  │
  ▼
Sequencer 시작
  │
  ▼
블록 1 생성
  ├─ L1 attributes deposit transaction 생성
  │  ├─ SourceHash:
  │  │   keccak256(
  │  │     L1_INFO_TX_HASH,
  │  │     l1BlockHash,
  │  │     seqNumber
  │  │   )
  │  ├─ From: 0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001
  │  ├─ To: 0x4200000000000000000000000000000000000015
  │  ├─ Data: abi.encode(
  │  │   number, timestamp, basefee, hash,
  │  │   sequenceNumber, batcherHash, ...
  │  │ )
  │  └─ IsSystemTransaction: false (Regolith 이후)
  │
  ├─ Deposit tx 실행
  │  └─ L1Block.setL1BlockValues(...) 호출
  │     └─ storage 초기화 ✅
  │
  └─ 블록 저장
     └─ Deposit tx with sourceHash 포함 ✅
```

### 4.2 SourceHash 계산

```go
// op-node/rollup/derive/deposit_source.go
type L1InfoDepositSource struct {
    L1BlockHash common.Hash
    SeqNumber   uint64
}

func (dep *L1InfoDepositSource) SourceHash() common.Hash {
    input := make([]byte, 32+32+32+32)

    // Domain identifier
    copy(input[:32], L1InfoDepositedTxType[:])

    // L1 block hash
    copy(input[32:64], dep.L1BlockHash.Bytes())

    // Sequence number
    binary.BigEndian.PutUint64(input[56:64], dep.SeqNumber)

    return crypto.Keccak256Hash(input)
}
```

---

## 5. E2E 시스템 구성

### 5.1 전체 아키텍처

```
┌─────────────────────────────────────────────────────────┐
│               Optimism E2E Test System                   │
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
│  ┌─────────────────────────────────────────────────┐     │
│  │         Fault Proof System (Optional)           │     │
│  ├─────────────────────────────────────────────────┤     │
│  │  • op-challenger (dispute resolution)           │     │
│  │  • op-program (fault proof VM)                  │     │
│  │  • cannon/asterisc (trace generation)           │     │
│  └─────────────────────────────────────────────────┘     │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### 5.2 기본 구성 다이어그램

```
사용자 트랜잭션 제출
        │
        ▼
┌──────────────────┐
│ L2 Sequencer     │
│ (op-geth)        │  1. 트랜잭션 실행
│                  │  2. 블록 생성
└────────┬─────────┘
         │
         ├────────────────────────────────┐
         │                                │
         ▼                                ▼
┌──────────────────┐            ┌──────────────────┐
│ op-batcher       │            │ op-proposer      │
│                  │            │                  │
│ Batch 생성       │            │ Output 계산      │
└────────┬─────────┘            └────────┬─────────┘
         │                                │
         ▼                                ▼
┌─────────────────────────────────────────────────┐
│              L1 Ethereum                         │
│                                                  │
│  • Batch 데이터 저장 (calldata/blob)             │
│  • Output Root 저장 (DisputeGameFactory)         │
└────────┬────────────────────────────────────────┘
         │
         ▼
┌──────────────────┐
│ L2 Verifier      │
│ (op-node)        │  1. L1에서 batch 읽기
│                  │  2. Block derivation
│ (op-geth)        │  3. 트랜잭션 재실행
│                  │  4. 상태 검증
└──────────────────┘
```

### 5.3 DefaultSystemConfig (기본 구성)

```go
// op-e2e/setup.go
func DefaultSystemConfig(t testing.TB) SystemConfig {
    // Forge allocs 로드
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
            "sequencer": sequencerConfig(),
            "verifier":  verifierConfig(),
        },

        // 로깅 설정
        Loggers: map[string]log.Logger{
            "sequencer": testlog.Logger(t, log.LevelInfo),
            "verifier":  testlog.Logger(t, log.LevelInfo),
        },

        // ...
    }
}
```

**노드 구성:**
- **L1 노드**: 1개 (Geth)
- **L2 Sequencer**: op-geth + op-node
- **L2 Verifier**: op-geth + op-node
- **총 3개 Geth 인스턴스**

### 5.4 Fault Proof 테스트 구성

```go
// Optimism op-e2e/faultproofs/helpers.go
func setupDisputeGameTest(t *testing.T) (*op_e2e.System, *ethclient.Client) {
    cfg := op_e2e.DefaultSystemConfig(t)

    // Verifier 제거 (테스트 단순화)
    delete(cfg.Nodes, "verifier")

    // Fault Proof 특화 설정
    cfg.DeployConfig.SequencerWindowSize = 4
    cfg.DeployConfig.FinalizationPeriodSeconds = 2
    cfg.SupportL1TimeTravel = true
    cfg.DeployConfig.L2OutputOracleSubmissionInterval = 1
    cfg.NonFinalizedProposals = true

    // 시스템 시작
    sys, err := cfg.Start(t)
    require.NoError(t, err)

    return sys, sys.Clients["l1"]
}
```

**간소화된 구성:**
- **L1 노드**: 1개
- **L2 Sequencer**: 1개 (op-geth + op-node)
- **총 2개 Geth 인스턴스**

### 5.5 각 컴포넌트 역할

#### L1 Layer
```
┌──────────────────────────────────────┐
│ L1 Geth                               │
├──────────────────────────────────────┤
│ • 스마트 컨트랙트 호스팅             │
│ • Batch 데이터 저장 (calldata/blob)  │
│ • Output Root 저장                    │
│ • Dispute Game 실행                   │
└──────────────────────────────────────┘
```

#### L2 Sequencer
```
┌──────────────────────────────────────┐
│ op-geth (Execution Engine)            │
├──────────────────────────────────────┤
│ • 트랜잭션 실행                       │
│ • 상태 변경                           │
│ • 블록 생성                           │
└─────────────┬────────────────────────┘
              │
┌─────────────▼────────────────────────┐
│ op-node (Consensus)                   │
├──────────────────────────────────────┤
│ • L1 데이터 derivation                │
│ • 블록 생성 트리거                    │
│ • P2P 네트워크 관리                   │
└─────────────┬────────────────────────┘
              │
┌─────────────▼────────────────────────┐
│ op-batcher (Batch Submitter)          │
├──────────────────────────────────────┤
│ • L2 트랜잭션 배치화                  │
│ • L1에 batch 제출                     │
└──────────────────────────────────────┘
              │
┌─────────────▼────────────────────────┐
│ op-proposer (Output Proposer)         │
├──────────────────────────────────────┤
│ • L2 상태 root 계산                   │
│ • L1에 output root 제출               │
└──────────────────────────────────────┘
```

#### L2 Verifier
```
┌──────────────────────────────────────┐
│ op-geth (Execution Engine)            │
├──────────────────────────────────────┤
│ • 트랜잭션 검증 실행                  │
│ • 상태 일치 확인                      │
└─────────────┬────────────────────────┘
              │
┌─────────────▼────────────────────────┐
│ op-node (Consensus)                   │
├──────────────────────────────────────┤
│ • L1에서 batch 가져오기               │
│ • 블록 derivation                     │
│ • Sequencer와 상태 비교               │
└──────────────────────────────────────┘
```

### 5.5 System Start (시스템 초기화)

```go
// op-e2e/setup.go
func (cfg SystemConfig) Start(t *testing.T) (*System, error) {
    // ========================================
    // 1. L1 Genesis 생성
    // ========================================
    l1Genesis, err := genesis.BuildL1DeveloperGenesis(
        cfg.DeployConfig,
        config.L1Allocs,    // ← Forge allocs
        config.L1Deployments,
    )
    if err != nil {
        return nil, err
    }

    l1Block := l1Genesis.ToBlock()

    // ========================================
    // 2. L2 Genesis 생성
    // ========================================
    l2Allocs := config.L2Allocs(genesis.L2AllocsFjord)  // ← Forge allocs
    l2Genesis, err := genesis.BuildL2Genesis(
        cfg.DeployConfig,
        l2Allocs,    // ← Forge allocs 사용!
        &eth.BlockRef{...},
    )
    if err != nil {
        return nil, err
    }

    // ========================================
    // 3. L1 Geth 초기화
    // ========================================
    l1Node, l1Backend, err := geth.InitL1(
        cfg.L1ChainIDBig().Uint64(),
        cfg.L1BlockTime,
        l1Genesis,
        cfg.JWTFilePath,
        cfg.JWTSecret,
        cfg.Loggers["l1"],
    )
    if err != nil {
        return nil, err
    }

    // ========================================
    // 4. L2 Geth 인스턴스들 초기화
    // ========================================
    for name := range cfg.Nodes {
        chainID := cfg.L2ChainIDBig(name).Uint64()

        l2Node, l2Backend, err := geth.InitL2(
            name,
            chainID,
            l2Genesis,
            cfg.JWTFilePath,
            cfg.JWTSecret,
            cfg.Loggers[name],
        )
        if err != nil {
            return nil, fmt.Errorf("failed to init L2 %s: %w", name, err)
        }

        ethInstances[name] = &opGeth.GethInstance{
            Backend: l2Backend,
            Node:    l2Node,
        }
    }

    // ========================================
    // 5. op-node 인스턴스들 시작
    // ========================================
    for name := range cfg.Nodes {
        opNode, err := startOpNode(name, cfg, l1Backend)
        if err != nil {
            return nil, fmt.Errorf("failed to start op-node %s: %w", name, err)
        }
        opNodes[name] = opNode
    }

    // ========================================
    // 6. op-batcher 시작 (Sequencer만)
    // ========================================
    if cfg.EnableBatcher {
        batcher, err := startBatcher(cfg, l1Backend)
        if err != nil {
            return nil, fmt.Errorf("failed to start batcher: %w", err)
        }
        sys.Batcher = batcher
    }

    // ========================================
    // 7. op-proposer 시작 (Sequencer만)
    // ========================================
    if cfg.EnableProposer {
        proposer, err := startProposer(cfg, l1Backend)
        if err != nil {
            return nil, fmt.Errorf("failed to start proposer: %w", err)
        }
        sys.Proposer = proposer
    }

    // ========================================
    // 8. op-challenger 시작 (Optional)
    // ========================================
    if cfg.EnableChallenger {
        challenger, err := startChallenger(cfg, l1Backend)
        if err != nil {
            return nil, fmt.Errorf("failed to start challenger: %w", err)
        }
        sys.Challenger = challenger
    }

    return sys, nil
}
```

### 5.6 노드 간 통신 흐름

```
사용자 트랜잭션 제출
        │
        ▼
┌──────────────────┐
│ L2 Sequencer     │
│ (op-geth)        │  1. 트랜잭션 실행
│                  │  2. 블록 생성
└────────┬─────────┘
         │
         ├────────────────────────────────┐
         │                                │
         ▼                                ▼
┌──────────────────┐            ┌──────────────────┐
│ op-batcher       │            │ op-proposer      │
│                  │            │                  │
│ Batch 생성       │            │ Output 계산      │
└────────┬─────────┘            └────────┬─────────┘
         │                                │
         ▼                                ▼
┌─────────────────────────────────────────────────┐
│              L1 Ethereum                         │
│                                                  │
│  • Batch 데이터 저장 (calldata/blob)             │
│  • Output Root 저장 (DisputeGameFactory)         │
└────────┬────────────────────────────────────────┘
         │
         ▼
┌──────────────────┐
│ L2 Verifier      │
│ (op-node)        │  1. L1에서 batch 읽기
│                  │  2. Block derivation
│ (op-geth)        │  3. 트랜잭션 재실행
│                  │  4. 상태 검증
└──────────────────┘
```

## 6. E2E 테스트 Genesis 로딩

### 6.1 Allocs 로딩 프로세스

```go
// op-e2e/setup.go
func (cfg SystemConfig) Start(t *testing.T) (*System, error) {
    // L1 Genesis 생성
    l1Genesis, err := genesis.BuildL1DeveloperGenesis(
        cfg.DeployConfig,
        config.L1Allocs,    // ← Forge allocs
        config.L1Deployments,
    )
    if err != nil {
        return nil, err
    }

    l1Block := l1Genesis.ToBlock()

    // L2 Genesis 생성
    l2Allocs := config.L2Allocs(genesis.L2AllocsFjord)  // ← Forge allocs
    l2Genesis, err := genesis.BuildL2Genesis(
        cfg.DeployConfig,
        l2Allocs,    // ← Forge allocs 사용!
        &eth.BlockRef{...},
    )
    if err != nil {
        return nil, err
    }

    // Geth 초기화
    l1Node, l1Backend, err := geth.InitL1(...)
    l2Node, l2Backend, err := geth.InitL2(name, chainID, l2Genesis, ...)

    // ...
}
```

---

## 7. Tokamak과의 비교

### 7.1 시스템 구성 비교

| 항목 | Optimism | Tokamak | 비고 |
|------|----------|---------|------|
| **기본 노드 구성** | L1 + Sequencer + Verifier | L1 + Sequencer + Verifier | 동일 ✅ |
| **Fault Proof 구성** | L1 + Sequencer (2개) | L1 + Sequencer (2개) | 동일 ✅ |
| **Geth 인스턴스 수** | 2-3개 (명확) | 4개 (불명확) ⚠️ | Tokamak에서 추가 인스턴스 발견 |
| **노드 관리** | `cfg.Nodes` map으로 관리 | `cfg.Nodes` map으로 관리 | 동일 ✅ |
| **로깅** | 노드별 Logger 설정 | 노드별 Logger 설정 | 동일 ✅ |

**Tokamak의 문제점:**
- E2E 테스트 실행 시 예상(2개)보다 많은 geth 인스턴스(4개) 실행
- 원인 불명 (이전 테스트 잔여물 또는 설정 오류 가능성)

### 7.2 Genesis 생성 비교

| 단계 | Optimism | Tokamak |
|------|----------|---------|
| **1. Forge 스크립트 실행** | ✅ 항상 실행 | ✅ Devnet만 실행 |
| **2. Allocs 파일 생성** | ✅ 생성 | ✅ 생성 |
| **3. Allocs 로딩** | ✅ Devnet + E2E | ✅ Devnet만 |
| **4. Go 코드 생성** | ❌ 안 함 | ⚠️ E2E에서 사용 |
| **5. L1Block Storage** | ✅ 비어있음 | ❌ 설정됨 (E2E) |

### 7.3 코드 비교

#### Optimism
```go
// op-chain-ops/genesis/layer_two.go
func BuildL2Genesis(
    config *DeployConfig,
    dump *foundry.ForgeAllocs,  // ← Forge allocs
    l1StartBlock *eth.BlockRef,
) (*core.Genesis, error) {
    genspec, err := NewL2Genesis(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    // Forge allocs 그대로 사용
    for addr, val := range dump.Copy().Accounts {
        genspec.Alloc[addr] = val
    }

    // ✅ L1Block storage 건드리지 않음

    return genspec, nil
}
```

#### Tokamak
```go
// op-chain-ops/genesis/layer_two.go
func BuildL2Genesis(
    config *DeployConfig,
    l1StartBlock *types.Block,  // ← Forge allocs 안 받음
) (*core.Genesis, error) {
    genspec, err := NewL2Genesis(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    db := state.NewMemoryStateDB(genspec)

    // ❌ Go 코드로 storage 설정
    storage, err := NewL2StorageConfig(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    // ❌ L1Block storage 설정됨!

    return db.Genesis(), nil
}
```

---

## 8. Optimism의 장점

### 8.1 일관성

```
Forge Script
     │
     ├─────────┬─────────┬─────────┐
     │         │         │         │
  Devnet    E2E Test  Testnet  Mainnet
     │         │         │         │
     └─────────┴─────────┴─────────┘
              모두 동일한 Genesis
```

**이점:**
- 버그 감소
- 예측 가능성
- 디버깅 용이

### 8.2 Solidity 기반

```
장점:
✅ 스마트 컨트랙트와 동일한 언어
✅ Storage layout 자동 계산
✅ 컨트랙트 업그레이드 시뮬레이션 가능
✅ Gas 계산 정확
```

### 8.3 Fork 버전 관리

```
allocs-l2-delta.json     // Delta fork
allocs-l2-ecotone.json   // Ecotone fork
allocs-l2-fjord.json     // Fjord fork
allocs-l2.json           // Latest
```

**이점:**
- Fork별 genesis 관리
- 하위 호환성 테스트
- 업그레이드 시뮬레이션

---

## 9. 권장 사항: Tokamak이 Optimism 방식을 채택해야 하는 이유

### 9.1 단기 이점

1. **sourceHash 에러 해결**
   - L1Block storage 비어있음 보장
   - Deposit transaction 정상 작동

2. **Devnet과 E2E 일관성**
   - 버그 감소
   - 디버깅 시간 단축

### 9.2 중기 이점

1. **코드 유지보수 간소화**
   - `NewL2StorageConfig()` 제거 가능
   - Storage 로직 중복 제거

2. **Upstream 동기화 용이**
   - Optimism 업데이트 반영 쉬움
   - 차이점 최소화

### 9.3 장기 이점

1. **Production 배포 안정성**
   - 테스트 환경 = 실제 환경
   - 예상치 못한 문제 감소

2. **커뮤니티 호환성**
   - Optimism 생태계와 일치
   - 도구 및 스크립트 공유 가능

---

## 10. 마이그레이션 가이드

### 10.1 Phase 1: E2E 테스트를 Allocs 기반으로 변경

```go
// op-e2e/config/init.go - 이미 존재하는 allocs 로딩 로직 활용
func L2Allocs(mode genesis.L2AllocsMode) *genesis.ForgeAllocs {
    allocs, ok := l2Allocs[mode]
    if !ok {
        panic(fmt.Errorf("unknown L2 allocs mode: %q", mode))
    }
    return allocs.Copy()
}
```

```go
// op-e2e/setup.go - BuildL2Genesis 호출 변경
func (cfg SystemConfig) Start(t *testing.T) (*System, error) {
    // ...

    l1Block := l1Genesis.ToBlock()

    // ✅ Forge allocs 사용
    l2Allocs := config.L2Allocs(genesis.L2AllocsFjord)
    l2Genesis, err := genesis.BuildL2Genesis(
        cfg.DeployConfig,
        l2Allocs,  // ← Forge allocs 전달
        l1Block.ToRef(),
    )

    // ...
}
```

### 9.2 Phase 2: BuildL2Genesis 함수 시그니처 변경

```go
// op-chain-ops/genesis/layer_two.go
func BuildL2Genesis(
    config *DeployConfig,
    dump *ForgeAllocs,      // ← allocs 추가
    l1StartBlock *types.Block,
) (*core.Genesis, error) {
    genspec, err := NewL2Genesis(config, l1StartBlock)
    if err != nil {
        return nil, err
    }

    // Forge allocs 사용
    if dump != nil {
        for addr, val := range dump.Copy().Accounts {
            genspec.Alloc[addr] = val
        }
        return genspec, nil
    }

    // Fallback: 기존 방식 (deprecated)
    db := state.NewMemoryStateDB(genspec)
    storage, err := NewL2StorageConfig(config, l1StartBlock)
    // ...
}
```

### 9.3 Phase 3: NewL2StorageConfig 제거

```go
// op-chain-ops/genesis/config.go
// func NewL2StorageConfig() 함수 삭제
// 모든 호출부를 Forge allocs 기반으로 변경
```

---

## 10. 결론

### 10.1 핵심 차이점

| 항목 | Optimism | Tokamak |
|------|----------|---------|
| **Genesis 소스** | Forge allocs | Go 코드 + Forge allocs |
| **일관성** | Devnet = E2E = Production | Devnet ≠ E2E |
| **L1Block 초기화** | 블록 1 (deposit tx) | Genesis (E2E) / 블록 1 (Devnet) |
| **유지보수성** | 높음 | 중간 |
| **Upstream 동기화** | 쉬움 | 어려움 |

### 10.2 권장 사항

**즉시 조치:**
1. E2E 테스트를 Forge allocs 기반으로 변경
2. Genesis 생성 로직 통일

**장기 목표:**
1. Optimism upstream 완전 동기화
2. Tokamak 특화 기능은 별도 레이어로 분리
3. 코드 중복 최소화

### 10.3 기대 효과

1. **sourceHash 에러 해결** ✅
2. **테스트 안정성 향상** ✅
3. **개발 속도 증가** ✅
4. **유지보수 비용 감소** ✅
5. **Upstream 업데이트 반영 용이** ✅

---

## 참고 자료

- [Optimism Genesis Generation](https://github.com/ethereum-optimism/optimism/tree/develop/op-chain-ops/genesis)
- [Optimism E2E Tests](https://github.com/ethereum-optimism/optimism/tree/develop/op-e2e)
- [Forge Scripts](https://github.com/ethereum-optimism/optimism/tree/develop/packages/contracts-bedrock/scripts)
- [L2 Genesis Script](https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/scripts/L2Genesis.s.sol)

