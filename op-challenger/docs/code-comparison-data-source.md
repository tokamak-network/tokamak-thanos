# tokamak-thanos vs Optimism op-challenger 코드 비교: 게임 시작 시 데이터 소스

## 비교 목적

게임을 시작할 때 원본 데이터를 어디서 가져오는지 두 코드베이스를 비교하여 확인합니다.

---

## 핵심 결론

**양쪽 코드 모두 동일한 소스에서 데이터를 가져옵니다:**

### 데이터 소스

| 데이터 종류 | 소스 | 이유 |
|------------|------|------|
| **L1 Header** | L1 RPC (`cfg.L1EthRpc`) | L1 상태 확인 |
| **L2 Header** | L2 RPC (`cfg.L2Rpc`) | L2 블록 정보 |
| **배치 데이터 (Batch Data)** | **L1 Ethereum** | 단일 진실의 원천 (Single Source of Truth) |

---

## 1. tokamak-thanos/op-challenger 코드 분석

### A. 게임 플레이어 생성 (player.go:69-159)

```go
func NewGamePlayer(
    ctx context.Context,
    systemClock clock.Clock,
    l1Clock types.ClockReader,
    logger log.Logger,
    m metrics.Metricer,
    dir string,
    addr common.Address,
    txSender TxSender,
    loader GameContract,
    syncValidator SyncValidator,
    validators []Validator,
    creator resourceCreator,  // ⭐️ TraceAccessor 생성 함수
    l1HeaderSource L1HeaderSource,
    selective bool,
    claimants []common.Address,
) (*GamePlayer, error) {
    // ... (생략) ...

    // ⭐️ L1 Head 가져오기
    l1HeadHash, err := loader.GetL1Head(ctx)
    l1Header, err := l1HeaderSource.HeaderByHash(ctx, l1HeadHash)
    l1Head := eth.HeaderBlockID(l1Header)

    // ⭐️ TraceAccessor 생성 (여기서 데이터 소스 결정)
    accessor, err := creator(ctx, logger, gameDepth, dir)

    // ... (생략) ...
}
```

**핵심**: `creator` 함수가 TraceAccessor를 생성하며, 여기서 데이터 소스가 결정됩니다.

### B. Cannon Trace Provider 등록 (register.go:280-373)

```go
func registerCannon(
    gameType uint32,
    registry Registry,
    oracles OracleRegistry,
    ctx context.Context,
    systemClock clock.Clock,
    l1Clock faultTypes.ClockReader,
    logger log.Logger,
    m metrics.Metricer,
    cfg *config.Config,  // ⭐️ L1/L2 RPC URL 포함
    syncValidator SyncValidator,
    rollupClient outputs.OutputRollupClient,
    txSender TxSender,
    gameFactory *contracts.DisputeGameFactoryContract,
    caller *batching.MultiCaller,
    l2Client utils.L2HeaderSource,  // ⭐️ L2 RPC Client
    l1HeaderSource L1HeaderSource,  // ⭐️ L1 RPC Client
    selective bool,
    claimants []common.Address,
) error {
    // ... (생략) ...

    playerCreator := func(game types.GameMetadata, dir string) (scheduler.GamePlayer, error) {
        // ... (생략) ...

        // ⭐️ L1 Head 로드 (L1 RPC 사용)
        l1HeadID, err := loadL1Head(contract, ctx, l1HeaderSource)

        // ⭐️ TraceAccessor creator 정의
        creator := func(ctx context.Context, logger log.Logger, gameDepth faultTypes.Depth, dir string) (faultTypes.TraceAccessor, error) {
            cannonPrestate, err := prestateSource.PrestatePath(requiredPrestatehash)

            // ⭐️ Cannon TraceAccessor 생성
            accessor, err := outputs.NewOutputCannonTraceAccessor(
                logger, m, cfg,
                l2Client,           // ⭐️ L2 RPC Client
                prestateProvider,
                cannonPrestate,
                rollupClient,       // ⭐️ Rollup Client (L2 상태)
                dir,
                l1HeadID,           // ⭐️ L1 Head
                splitDepth,
                prestateBlock,
                poststateBlock)
            return accessor, nil
        }

        return NewGamePlayer(ctx, systemClock, l1Clock, logger, m, dir, game.Proxy,
            txSender, contract, syncValidator,
            []Validator{prestateValidator, startingValidator},
            creator,  // ⭐️ creator 함수 전달
            l1HeaderSource, selective, claimants)
    }
    // ... (생략) ...
}
```

### C. Output Cannon TraceAccessor 생성 (output_cannon.go:21-50)

```go
func NewOutputCannonTraceAccessor(
    logger log.Logger,
    m metrics.Metricer,
    cfg *config.Config,  // ⭐️ L1/L2 RPC URL 포함
    l2Client utils.L2HeaderSource,  // ⭐️ L2 RPC
    prestateProvider types.PrestateProvider,
    cannonPrestate string,
    rollupClient OutputRollupClient,  // ⭐️ Rollup RPC
    dir string,
    l1Head eth.BlockID,  // ⭐️ L1 Head Hash
    splitDepth types.Depth,
    prestateBlock uint64,
    poststateBlock uint64,
) (*trace.Accessor, error) {
    // ⭐️ Output Provider 생성 (L2 상태 추적)
    outputProvider := NewTraceProvider(logger, prestateProvider, rollupClient,
        l2Client, l1Head, splitDepth, prestateBlock, poststateBlock)

    // ⭐️ Cannon Creator: 로컬 입력을 가져와 Cannon 실행
    cannonCreator := func(ctx context.Context, localContext common.Hash, depth types.Depth,
        agreed contracts.Proposal, claimed contracts.Proposal) (types.TraceProvider, error) {

        subdir := filepath.Join(dir, localContext.Hex())

        // ⭐️ 로컬 입력 가져오기 (L2 Client 사용)
        localInputs, err := utils.FetchLocalInputsFromProposals(ctx, l1Head.Hash,
            l2Client, agreed, claimed)

        // ⭐️ Cannon TraceProvider 생성
        provider := cannon.NewTraceProvider(logger, m, cfg, prestateProvider,
            cannonPrestate, localInputs, subdir, depth)
        return provider, nil
    }

    cache := NewProviderCache(m, "output_cannon_provider", cannonCreator)
    selector := split.NewSplitProviderSelector(outputProvider, splitDepth,
        OutputRootSplitAdapter(outputProvider, cache.GetOrCreate))
    return trace.NewAccessor(selector), nil
}
```

### D. 로컬 입력 가져오기 (utils/local.go:47-61)

```go
func FetchLocalInputsFromProposals(
    ctx context.Context,
    l1Head common.Hash,  // ⭐️ L1 Head Hash
    l2Client L2HeaderSource,  // ⭐️ L2 RPC Client
    agreedOutput contracts.Proposal,
    claimedOutput contracts.Proposal
) (LocalGameInputs, error) {
    // ⭐️ L2 Header 가져오기 (L2 RPC 사용)
    agreedHeader, err := l2Client.HeaderByNumber(ctx, agreedOutput.L2BlockNumber)
    if err != nil {
        return LocalGameInputs{}, fmt.Errorf("fetch L2 block header %v: %w",
            agreedOutput.L2BlockNumber, err)
    }
    l2Head := agreedHeader.Hash()

    return LocalGameInputs{
        L1Head:        l1Head,  // L1 Hash
        L2Head:        l2Head,  // L2 Hash
        L2OutputRoot:  agreedOutput.OutputRoot,
        L2Claim:       claimedOutput.OutputRoot,
        L2BlockNumber: claimedOutput.L2BlockNumber,
    }, nil
}
```

**핵심**:
- L1 Head: L1 RPC에서 가져옴
- L2 Head: L2 RPC에서 가져옴
- Output Root: 게임 컨트랙트에서 가져옴

### E. Cannon Executor (executor.go:37-102)

```go
func NewExecutor(logger log.Logger, m CannonMetricer, cfg *config.Config,
    prestate string, inputs utils.LocalGameInputs) *Executor {
    return &Executor{
        logger:           logger,
        metrics:          m,
        l1:               cfg.L1EthRpc,      // ⭐️ L1 RPC URL
        l1Beacon:         cfg.L1Beacon,      // ⭐️ L1 Beacon RPC URL
        l2:               cfg.L2Rpc,         // ⭐️ L2 RPC URL
        inputs:           inputs,
        cannon:           cfg.CannonBin,
        server:           cfg.CannonServer,
        network:          cfg.CannonNetwork,
        rollupConfig:     cfg.CannonRollupConfigPath,
        l2Genesis:        cfg.CannonL2GenesisPath,
        absolutePreState: prestate,
        snapshotFreq:     cfg.CannonSnapshotFreq,
        infoFreq:         cfg.CannonInfoFreq,
        selectSnapshot:   utils.FindStartingSnapshot,
        cmdExecutor:      utils.RunCmd,
    }
}

func (e *Executor) generateProof(...) error {
    // ... (생략) ...

    args = append(args,
        "--",
        e.server, "--server",
        "--l1", e.l1,              // ⭐️ L1 RPC URL
        "--l1.beacon", e.l1Beacon, // ⭐️ L1 Beacon RPC URL
        "--l2", e.l2,              // ⭐️ L2 RPC URL
        "--datadir", dataDir,
        "--l1.head", e.inputs.L1Head.Hex(),
        "--l2.head", e.inputs.L2Head.Hex(),
        "--l2.outputroot", e.inputs.L2OutputRoot.Hex(),
        "--l2.claim", e.inputs.L2Claim.Hex(),
        "--l2.blocknumber", e.inputs.L2BlockNumber.Text(10),
    )
    // ... (생략) ...
}
```

**핵심**:
- `--l1`: L1 Execution Layer RPC (Calldata 가져오기용)
- `--l1.beacon`: L1 Consensus Layer RPC (Blob 가져오기용)
- `--l2`: L2 RPC (현재 L2 상태 확인용)

---

## 2. Optimism/op-challenger 코드 분석

### A. 게임 플레이어 생성 (player.go:75-162)

```go
func NewGamePlayer(
    ctx context.Context,
    systemClock clock.Clock,
    l1Clock types.ClockReader,
    logger log.Logger,
    m metrics.Metricer,
    dir string,
    addr common.Address,
    txSender TxSender,
    loader GameContract,
    syncValidator SyncValidator,
    validators []Validator,
    creator resourceCreator,  // ⭐️ TraceAccessor 생성 함수
    l1HeaderSource L1HeaderSource,
    selective bool,
    claimants []common.Address,
) (*GamePlayer, error) {
    // ... (생략) ...

    // ⭐️ L1 Head 가져오기
    l1HeadHash, err := loader.GetL1Head(ctx)
    l1Header, err := l1HeaderSource.HeaderByHash(ctx, l1HeadHash)
    l1Head := eth.HeaderBlockID(l1Header)

    // ⭐️ TraceAccessor 생성
    accessor, err := creator(ctx, logger, gameDepth, dir)

    // ... (생략) ...
}
```

**tokamak-thanos와 동일**

### B. 코드 구조 차이점

Optimism은 리팩토링하여 `RegisterTask` 패턴을 사용하지만, 핵심 로직은 동일:

```go
// register.go:63-68
if cfg.TraceTypeEnabled(faultTypes.TraceTypeCannon) {
    l2HeaderSource, rollupClient, syncValidator, err := clients.SingleChainClients()
    registerTasks = append(registerTasks, NewCannonRegisterTask(
        faultTypes.CannonGameType, cfg, m,
        vm.NewOpProgramServerExecutor(logger),
        l2HeaderSource, rollupClient, syncValidator))
}
```

---

## 3. 데이터 소스 비교 요약

### tokamak-thanos와 Optimism 공통점

| 데이터 | 소스 | 파일 위치 | 용도 |
|--------|------|-----------|------|
| **L1 Head Hash** | L1 RPC | executor.go:41 (`cfg.L1EthRpc`) | L1 상태 기준점 |
| **L1 Beacon** | L1 Beacon RPC | executor.go:42 (`cfg.L1Beacon`) | Blob 데이터 가져오기 |
| **L2 Head Hash** | L2 RPC | local.go:48 (`l2Client.HeaderByNumber`) | L2 블록 정보 |
| **L2 Output Root** | 게임 컨트랙트 | local.go:57-58 (`agreedOutput.OutputRoot`) | 합의된 상태 |
| **배치 데이터** | **L1 Ethereum** | op-program (prefetcher) | 재파생용 원본 데이터 |

### 배치 데이터의 실제 소스

```
┌─────────────────────────────────────────────────────────────┐
│             배치 데이터 가져오기 흐름                         │
└─────────────────────────────────────────────────────────────┘

1. op-challenger가 Cannon 실행
   ↓
2. Cannon이 op-program 서버 실행
   ↓
3. op-program이 prefetcher 실행
   ↓
4. prefetcher가 L1 접근:
   ┌──────────────────────────────────────┐
   │ L1 Execution Layer (Calldata)        │
   │ - BatchInboxAddress TX 조회           │
   │ - calldata 필드에서 배치 데이터 추출  │
   │                                       │
   │ L1 Consensus Layer (Blob)            │
   │ - Beacon API로 Blob 조회              │
   │ - Blob 데이터 다운로드                │
   └──────────────────────────────────────┘
   ↓
5. 배치 데이터로 L2 재파생
   ↓
6. 재파생된 L2 상태와 게임 claim 비교
```

---

## 4. 핵심 차이점

### 구조적 차이

| 항목 | tokamak-thanos | Optimism |
|------|----------------|----------|
| **등록 방식** | 함수 기반 (registerCannon) | RegisterTask 클래스 |
| **VM Executor** | 직접 Executor 생성 | vm.NewOpProgramServerExecutor |
| **Preimage Loader** | kvstore.NewDiskKV().Get | 함수 래퍼로 lazy 로딩 |
| **코드 구조** | 단순 함수 호출 | 객체 지향 Task 패턴 |

### 데이터 소스는 동일!

**양쪽 모두:**
- L1 RPC를 통해 L1 데이터 가져옴
- L1 Beacon RPC를 통해 Blob 가져옴 (사용 시)
- L2 RPC는 현재 L2 상태 확인용 (배치 데이터 소스 아님!)
- **배치 데이터는 무조건 L1에서 가져옴**

---

## 5. 코드 경로 추적 예시

### tokamak-thanos 경로

```
1. main() → RegisterGameTypes()
   ↓ (register.go:53)

2. RegisterGameTypes() → registerCannon()
   ↓ (register.go:280)

3. registerCannon() → playerCreator()
   ↓ (register.go:313)

4. playerCreator() → outputs.NewOutputCannonTraceAccessor()
   ↓ (register.go:352)

5. NewOutputCannonTraceAccessor() → cannon.NewTraceProvider()
   ↓ (output_cannon.go:43)

6. NewTraceProvider() → NewExecutor()
   ↓ (provider.go:44)

7. NewExecutor() → cfg.L1EthRpc, cfg.L1Beacon, cfg.L2Rpc 저장
   ↓ (executor.go:37-54)

8. generateProof() → op-program 실행
   ↓ (executor.go:66-102)

9. op-program → L1/L2/Beacon RPC 사용하여 데이터 가져오기
   ↓ (host/host.go:186-211)

10. L1에서 배치 데이터 가져오기
    ↓ (prefetcher/prefetcher.go)

✅ 최종: L1 Ethereum이 데이터 소스!
```

### Optimism 경로 (유사)

```
1. main() → RegisterGameTypes()
   ↓ (register.go:45)

2. RegisterGameTypes() → NewCannonRegisterTask()
   ↓ (register.go:63-68)

3. RegisterTask.Register() → playerCreator()
   ↓ (task_cannon.go)

4. 이후 경로는 tokamak-thanos와 동일
   ↓

✅ 최종: L1 Ethereum이 데이터 소스!
```

---

## 6. 결론

### 양쪽 코드베이스 모두:

1. **L1 RPC**: L1 상태 확인 및 Calldata 배치 데이터 가져오기
2. **L1 Beacon RPC**: Blob 배치 데이터 가져오기 (EIP-4844 사용 시)
3. **L2 RPC**: 현재 L2 상태 확인 (배치 데이터 소스 아님!)
4. **Rollup RPC**: L2 output root 가져오기

### 배치 데이터 소스

```
❌ L2 아카이브 노드에서 가져오지 않음
✅ L1 Ethereum에서 직접 가져옴
   - Calldata: L1 Execution Layer
   - Blob: L1 Consensus Layer (비콘체인)

이유: Optimistic Rollup의 보안 모델
→ L1이 단일 진실의 원천 (Single Source of Truth)
→ 탈중앙화된 검증 보장
```

### 코드 차이점

- **구조**: Optimism이 더 모듈화 (RegisterTask 패턴)
- **데이터 소스**: **완전히 동일**
- **보안 모델**: **완전히 동일**

---

**작성일**: 2025-10-17
**기반 코드**:
- tokamak-thanos: commit ef63c0e65
- optimism: latest main branch
