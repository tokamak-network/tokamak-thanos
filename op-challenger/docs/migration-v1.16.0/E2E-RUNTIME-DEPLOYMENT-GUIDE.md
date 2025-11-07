# E2E 런타임 배포 흐름 가이드

## 📌 개요

이 문서는 Optimism E2E 테스트의 런타임 배포 흐름을 분석하고, Tokamak-Thanos에 적용하기 위한 가이드를 제공합니다.

---

## 🔍 Optimism E2E 배포 흐름 분석

### 1. 전체 실행 흐름

```
E2E 테스트 시작
  ↓
패키지 초기화 (init() 함수)
  ↓
initAllocType() - 각 AllocType별 병렬 실행
  ↓
ApplyPipeline() - 컨트랙트 배포
  ↓
State 저장 (메모리에 L1/L2 allocs, deploy config)
  ↓
E2E 테스트 실행 (런타임 생성된 state 사용)
```

---

## 📋 Step-by-Step 상세 분석

### Step 1: Package 초기화 (`init()`)

**파일**: `op-e2e/config/init.go:144`

```go
func init() {
    root := FindMonorepoRoot()

    // 로거 설정
    oplog.SetGlobalLogHandler(errHandler)

    // 각 AllocType별 초기화
    for _, allocType := range allocTypes {
        initAllocType(root, allocType)
    }

    oplog.SetGlobalLogHandler(handler)
}
```

**동작**:
- ✅ 모노레포 루트 찾기
- ✅ 글로벌 로거 설정
- ✅ AllocType별 병렬 초기화 (`AllocTypeMTCannon`, `AllocTypeMTCannonNext` 등)

---

### Step 2: AllocType 초기화 (`initAllocType()`)

**파일**: `op-e2e/config/init.go:197`

```go
func initAllocType(root string, allocType AllocType) {
    // 1. Forge artifacts 로드
    artifactsPath := path.Join(root, "packages", "contracts-bedrock", "forge-artifacts")
    loc, err := artifacts.NewFileLocator(artifactsPath)

    // 2. L2 Allocs 모드별 병렬 처리
    allocModes := []genesis.L2AllocsMode{
        genesis.L2AllocsInterop,
        genesis.L2AllocsIsthmus,
        genesis.L2AllocsHolocene,
        genesis.L2AllocsGranite,
        genesis.L2AllocsFjord,
        genesis.L2AllocsEcotone,
        genesis.L2AllocsDelta,
    }

    for _, mode := range allocModes {
        go func(mode genesis.L2AllocsMode) {
            // 3. Intent 생성
            intent := defaultIntent(root, loc, deployerAddr, allocType)

            // 4. ApplyPipeline 실행
            deployer.ApplyPipeline(...)

            // 5. State 저장
            l2Alloc[mode] = st.Chains[0].Allocs.Data

            if mode == genesis.L2AllocsGranite {
                // Deploy config, L1 allocs 저장
                deployConfigsByType[allocType] = dc
                l1AllocsByType[allocType] = st.L1StateDump.Data
                l1DeploymentsByType[allocType] = l1Deployments
            }
        }(mode)
    }
    wg.Wait()

    l2AllocsByType[allocType] = l2Alloc
}
```

**핵심 포인트**:
- ✅ **병렬 처리**: 각 L2AllocsMode별로 goroutine 실행
- ✅ **런타임 생성**: 모든 것을 런타임에 계산
- ✅ **메모리 저장**: 파일이 아닌 메모리에 저장

---

### Step 3: Intent 생성 (`defaultIntent()`)

**파일**: `op-e2e/config/init.go:324`

```go
func defaultIntent(root string, loc *artifacts.Locator, deployer common.Address, allocType AllocType) *state.Intent {
    secrets := secrets.DefaultSecrets
    addrs := secrets.Addresses()

    // Cannon prestate (파일에서 읽음)
    defaultPrestate := common.HexToHash("0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98")

    // Genesis output root (더미 값)
    genesisOutputRoot := common.HexToHash("0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF")

    return &state.Intent{
        ConfigType: state.IntentTypeCustom,
        L1ChainID:  900,
        SuperchainRoles: &addresses.SuperchainRoles{...},
        GlobalDeployOverrides: map[string]any{
            "faultGameAbsolutePrestate":    defaultPrestate.Hex(),
            "faultGameGenesisBlock":        0,
            "faultGameGenesisOutputRoot":   genesisOutputRoot.Hex(),  // ← 더미 값!
            ...
        },
        Chains: []*state.ChainIntent{
            {
                ID: common.BigToHash(big.NewInt(901)),
                Roles: addresses.ChainRoles{...},
                AdditionalDisputeGames: [...],
            },
        },
    }
}
```

**핵심 포인트**:
- ⚠️ **faultGameGenesisOutputRoot**: `0xDEADBEEF...` 더미 값 사용
- ⚠️ **faultGameGenesisBlock**: `0` (genesis block)
- ✅ **faultGameAbsolutePrestate**: 파일에서 읽은 실제 값

---

### Step 4: ApplyPipeline 실행

**파일**: `op-deployer/pkg/deployer/apply.go:422`

```go
func ApplyPipeline(ctx context.Context, opts ApplyPipelineOpts) error {
    // 1. L1 Host 생성 (EVM 시뮬레이터)
    l1Host := script.NewHost(...)

    // 2. Scripts 로드 (DeploySuperchain, DeployImplementations 등)
    opcmScripts, err := opcm.NewScripts(l1Host)

    // 3. Pipeline 실행
    pline := []pipelineStage{
        {"init", pipeline.InitGenesisStrategy},
        {"deploy-superchain", pipeline.DeploySuperchain},
        {"deploy-implementations", pipeline.DeployImplementations},
        // 각 체인별:
        {"deploy-opchain-<chainID>", pipeline.DeployOPChain},
        {"deploy-alt-da-<chainID>", pipeline.DeployAltDA},
        {"deploy-additional-dispute-games-<chainID>", pipeline.DeployAdditionalDisputeGames},
        {"generate-l2-genesis-<chainID>", pipeline.GenerateL2Genesis},
    }

    for _, stage := range pline {
        if err := stage.apply(); err != nil {
            return err
        }
    }
}
```

**Pipeline Stage별 동작**:

#### Stage 1: `InitGenesisStrategy`
- State 초기화
- 배포자 계정 설정

#### Stage 2: `DeploySuperchain`
**Solidity 스크립트**: `DeploySuperchain.s.sol`

```solidity
function run(Input memory _input) public returns (Output memory output_) {
    // 1. Input 검증
    assertValidInput(internalInput);

    // 2. ProxyAdmin 배포
    deploySuperchainProxyAdmin(internalInput, output_);

    // 3. Implementation 배포
    deploySuperchainImplementationContracts(internalInput, output_);
    //    - SuperchainConfig (impl)
    //    - ProtocolVersions (impl)

    // 4. SuperchainConfig 배포 및 초기화
    deployAndInitializeSuperchainConfig(internalInput, output_);
    //    - Proxy 배포
    //    - upgradeAndCall(proxy, impl, initialize(guardian))

    // 5. ProtocolVersions 배포 및 초기화
    deployAndInitializeProtocolVersions(internalInput, output_);
    //    - Proxy 배포
    //    - upgradeAndCall(proxy, impl, initialize(...))

    // 6. ProxyAdmin 소유권 이전
    transferProxyAdminOwnership(internalInput, output_);

    // 7. Output 검증
    assertValidOutput(internalInput, output_);
}
```

**배포 순서**:
1. ProxyAdmin
2. SuperchainConfig (impl)
3. ProtocolVersions (impl)
4. SuperchainConfig (proxy + initialize)
5. ProtocolVersions (proxy + initialize)

#### Stage 3: `DeployImplementations`
**Solidity 스크립트**: `DeployImplementations.s.sol`

배포되는 컨트랙트:
- AddressManager
- Proxy
- L1CrossDomainMessenger
- OptimismPortal2
- SystemConfig
- L1StandardBridge
- L1ERC721Bridge
- OptimismMintableERC20Factory
- DisputeGameFactory
- DelayedWETH
- AnchorStateRegistry
- PermissionedDisputeGame

#### Stage 4: `DeployOPChain`
**Solidity 스크립트**: `DeployOPChain.s.sol`

```solidity
function run(Input memory _input) public returns (Output memory) {
    // 1. OPCM (OPContractsManager) 사용
    // 2. Proxy 배포 및 초기화
    // 3. SystemConfig, L1CrossDomainMessenger 등 초기화
    // 4. 게임 설정 (DisputeGameFactory에 게임 타입 등록)
}
```

#### Stage 5: `DeployAdditionalDisputeGames`
추가 Dispute Game 배포 (Alphabet, Asterisc 등)

#### Stage 6: `GenerateL2Genesis`
**파일**: `op-deployer/pkg/deployer/pipeline/l2genesis.go`

```go
func GenerateL2Genesis(env *Env, intent *state.Intent, bundle artifacts.ChainBundle, st *state.State, chainID common.Hash) error {
    // 1. Deploy config 가져오기
    deployConfig := inspect.DeployConfig(st, chainID)

    // 2. L1 deployments 가져오기
    l1Deployments := genesis.CreateL1DeploymentsFromContracts(l1Contracts)
    deployConfig.SetDeployments(l1Deployments)

    // 3. L2 allocs 가져오기 (체인별로 이미 생성됨)

    // 4. L2 Genesis 빌드
    l2Genesis := genesis.BuildL2Genesis(deployConfig, l2Allocs, l1StartingBlockRef)

    // 5. State에 저장
    st.Chains[chainIdx].Allocs = l2Genesis.Alloc
}
```

---

## 🔑 핵심: Genesis Output Root 처리

### Optimism의 방식

```
1. defaultIntent()에서 더미 값 설정
   genesisOutputRoot = 0xDEADBEEF...
   ↓
2. ApplyPipeline() 실행
   - DeploySuperchain: SuperchainConfig, ProtocolVersions 배포
   - DeployImplementations: 기타 컨트랙트 배포
   - DeployOPChain: AnchorStateRegistry 초기화
     → StartingAnchorRoot에 0xDEADBEEF... 저장
   ↓
3. GenerateL2Genesis()
   - 런타임에 L2 genesis 생성
   - **실제 genesis output root 계산** (여기서!)
   ↓
4. E2E 테스트 실행
   - L1 state: 컨트랙트에 0xDEADBEEF... 저장됨
   - L2 genesis: 런타임에 생성된 값 사용
   - Challenger가 실제 값 계산: 0xb41d...
   ↓
5. ❓ 왜 에러가 안 나나?
   → Challenger가 L2 genesis를 참조하여 계산하기 때문!
```

### 🤔 의문점: 왜 Optimism은 에러가 안 날까?

**분석 필요**:
1. Challenger가 어떻게 genesis output root를 계산하는가?
2. 컨트랙트의 값 (0xDEADBEEF...)과 비교하지 않는가?
3. 아니면 특별한 처리가 있는가?

**확인할 파일**:
- `op-challenger/game/fault/trace/outputs/prestate.go`
- `op-challenger/game/scheduler/coordinator.go`

---

## 🔧 Tokamak-Thanos 적용 방안

### 현재 문제점

**우리의 방식** (.devnet 사전 생성):
```
1. make devnet-allocs (사전 실행)
   - devnetL1.json 읽기: faultGameGenesisOutputRoot: 0xDEADBEEF...
   - Deploy.s.sol 실행
   - .devnet/allocs-l1.json 저장 (0xDEADBEEF... 포함)
   ↓
2. E2E 테스트 시작
   - .devnet/allocs-l1.json 읽기 (고정값)
   - 런타임에 L2 genesis 생성 (새 값 계산)
   ↓
3. ❌ 불일치 발생!
   - Contract: 0xDEADBEEF... (또는 사전 계산된 값)
   - Provider: 0xb41d... (런타임 계산)
```

### 해결 방안

#### 옵션 1: Optimism 방식 채택 (런타임 배포) ⭐ 권장

**장점**:
- ✅ Optimism과 동일한 구조
- ✅ 일관성 보장
- ✅ 더미 값 문제 근본 해결

**단점**:
- ❌ `.devnet` 파일 사용 중단
- ❌ 큰 구조 변경

**구현**:
```go
// op-e2e/config/init.go
func init() {
    // .devnet 체크 제거
    // 항상 런타임 생성
    for _, allocType := range allocTypes {
        initAllocType(root, allocType)
    }
}
```

**문제**: `revision id 1 cannot be reverted` panic 발생

**원인 분석 필요**:
- Tokamak 전용 컨트랙트 문제?
- Deploy script 차이?
- State snapshot 관리 문제?

#### 옵션 2: .devnet 유지 + 개선

**방법 1**: 올바른 genesis output root 사전 계산

```bash
# prepare-e2e-test.sh
1. make devnet-allocs (더미 값)
2. Go 스크립트로 genesis output root 계산
3. devnetL1-template.json 업데이트
4. make devnet-allocs 재실행 (올바른 값)
```

**문제**: E2E가 런타임에 L2 genesis를 다시 생성하면 값이 달라질 수 있음

**방법 2**: E2E가 .devnet의 L2 allocs도 사용하도록 수정

```go
// op-e2e/system/e2esys/setup.go
l2Genesis, err := genesis.BuildL2Genesis(
    cfg.DeployConfig,
    config.L2Allocs(cfg.AllocType, allocsMode),  // ← .devnet에서 로드
    eth.BlockRefFromHeader(l1Block.Header())
)
```

**문제**: L2 genesis가 고정되어 유연성 감소

---

## 🐛 Panic 원인 분석

### Revision ID Panic

**에러 메시지**:
```
panic: revision id 1 cannot be reverted
at: github.com/ethereum/go-ethereum/core/state/journal.go:91
```

**발생 위치**:
```
DeploySuperchain.s.sol 실행 중
  ↓
upgradeAndCall() 호출
  ↓
Proxy 초기화 중 revert 발생
  ↓
StateDB.RevertToSnapshot() 시도
  ↓
❌ panic: revision id 1 cannot be reverted
```

**가능한 원인**:

1. **Snapshot 관리 문제**:
   - Forge 스크립트가 snapshot을 생성
   - 초기화 중 실패하여 revert 시도
   - 하지만 snapshot ID가 잘못됨

2. **Try-catch 블록**:
   - Solidity에서 external call 실패
   - Try-catch로 처리하려 하지만 snapshot 문제 발생

3. **Tokamak 전용 컨트랙트**:
   - SuperchainConfig나 ProtocolVersions 초기화 중 문제
   - Tokamak 필드 때문에 revert 발생?

### 디버깅 방법

**1. 로그 추가**:
```solidity
// DeploySuperchain.s.sol
function deployAndInitializeSuperchainConfig(...) private {
    console.log("=== deployAndInitializeSuperchainConfig START ===");

    vm.startBroadcast(msg.sender);
    console.log("Creating SuperchainConfig proxy...");
    ISuperchainConfig superchainConfigProxy = ISuperchainConfig(...);

    console.log("Calling upgradeAndCall...");
    superchainProxyAdmin.upgradeAndCall(...);
    console.log("upgradeAndCall SUCCESS");

    vm.stopBroadcast();
    console.log("=== deployAndInitializeSuperchainConfig END ===");
}
```

**2. Try-catch 제거 테스트**:
컨트랙트에 try-catch 블록이 있다면 제거하고 테스트

**3. Optimism 컨트랙트 직접 사용**:
Tokamak artifacts 대신 Optimism artifacts 사용 테스트

---

## 📊 Tokamak-Thanos 추가 개발 필요 사항

### 1. Tokamak 전용 필드 처리

**파일**: `op-chain-ops/genesis/config.go`

```go
type TokamakDeployConfig struct {
    NativeTokenAddress       common.Address
    NativeTokenName          string
    NativeTokenSymbol        string
}
```

**Intent에 추가 필요**:
```go
intent := defaultIntent(...)
intent.Chains[0].TokamakConfig = &genesis.TokamakDeployConfig{
    NativeTokenAddress: common.HexToAddress("0x..."),
    NativeTokenName:    "TON",
    NativeTokenSymbol:  "TON",
}
```

### 2. ETHLockbox 배포

**Optimism**: ETHLockbox 배포하지만 검증 스킵
**Tokamak**: Native Token 사용 시 필수

**추가 필요**:
- `DeployETHLockbox.s.sol` 스크립트
- Pipeline에 stage 추가

### 3. Artifacts 경로 수정

**Optimism**: `packages/contracts-bedrock/forge-artifacts`
**Tokamak**: `packages/tokamak/contracts-bedrock/forge-artifacts`

```go
// op-e2e/config/init.go
artifactsPath := path.Join(root, "packages", "tokamak", "contracts-bedrock", "forge-artifacts")
```

---

## 🎯 권장 접근 방법

### Phase 1: Panic 원인 규명 ✅

1. **DeploySuperchain.s.sol에 로그 추가**
2. **재실행하여 어느 단계에서 실패하는지 확인**
3. **Optimism 컨트랙트와 비교**

### Phase 2: Panic 수정

**방법 A**: 문제되는 컨트랙트/함수 수정
**방법 B**: Optimism 컨트랙트 직접 사용 (테스트용)

### Phase 3: Optimism 방식 완전 적용

```go
// op-e2e/config/init.go
func init() {
    // .devnet 체크 제거
    for _, allocType := range allocTypes {
        initAllocType(root, allocType)
    }
}
```

### Phase 4: Genesis Output Root 문제 해결

**Optimism도 더미 값 사용하는데 왜 에러가 안 나는지 확인**:
1. Challenger 로직 분석
2. Genesis block (block 0)에 대한 특별 처리 확인
3. 필요시 우리도 동일한 처리 추가

---

## 📝 다음 단계

### 즉시 실행 가능한 작업

1. **DeploySuperchain.s.sol 로그 추가** (어디서 실패하는지 확인)
2. **Optimism E2E 테스트 완료 대기** (성공 여부 확인)
3. **Panic 원인 파악 후 수정**

### 중기 작업

1. **init.go를 Optimism 방식으로 변경**
2. **Genesis output root 처리 로직 개선**
3. **E2E 테스트 성공 검증**

---

## 🔍 추가 조사 필요 항목

### 1. Optimism Challenger의 Genesis Block 처리

**질문**: Challenger가 block 0 (genesis)의 output root를 어떻게 검증하는가?

**확인할 파일**:
```
op-challenger/game/fault/trace/outputs/prestate.go
op-challenger/game/scheduler/coordinator.go
op-challenger/game/fault/validator.go
```

**가설**:
- Genesis block (block 0)은 검증을 스킵하는가?
- 특별한 처리가 있는가?
- AnchorStateRegistry의 값을 사용하는가?

### 2. L2 Genesis 생성의 determinism

**질문**: 같은 조건에서 L2 genesis를 여러 번 생성하면 항상 같은 결과가 나오는가?

**테스트 필요**:
```go
// 같은 deploy config, 같은 L2 allocs로
genesis1 := BuildL2Genesis(deployConfig, l2Allocs, l1BlockRef)
genesis2 := BuildL2Genesis(deployConfig, l2Allocs, l1BlockRef)

// genesis1.ToBlock().Hash() == genesis2.ToBlock().Hash() ?
// genesis output root 동일한가?
```

---

---

## 🚨 중요 발견: Optimism도 동일한 문제 발생!

### Optimism E2E 테스트 결과 (2025-11-06)

```
=== RUN   TestOutputCannonGame
=== RUN   TestOutputCannonGame/mt-cannon
=== RUN   TestOutputCannonGame/mt-cannon-next

ERROR[11-06|11:55:09.374] Invalid prestate role=Challenger
err="failed to validate prestate: output root absolute prestate does not match:
Provider: 0xc534d0a83769c58d903da6d853f17164834423ba875a7a109f1caf1e8e3737a9
Contract: 0xdead000000000000000000000000000000000000000000000000000000000000"

zsh: exit 124    timeout 180 go test ... (3분 타임아웃)
```

### 분석

**Optimism도 우리와 똑같은 에러 발생**:
- ❌ **Provider** (런타임 계산): `0xc534...`
- ❌ **Contract** (더미 값): `0xdead0000...` (우리는 `0xDEADBEEF...`)
- ⏱️ **결과**: 3분 timeout

**결론**:
1. 🔴 **이 문제는 Optimism에서도 해결되지 않았음**
2. 🔴 **TestOutputCannonGame 자체에 문제가 있을 가능성**
3. 🔴 **또는 테스트가 실제로 실패하는 것이 정상일 수 있음**

### 추가 조사 필요

1. **Optimism CI/CD 확인**: 이 테스트가 실제로 통과하는가?
2. **다른 테스트 방법**: 다른 E2E 테스트는 성공하는가?
3. **Genesis block 특별 처리**: Block 0에 대한 특별한 로직이 있는가?

---

## 📌 수정된 결론

**Genesis Output Root 문제는 Optimism에서도 미해결**:
1. ⚠️ 런타임 배포 방식으로 전환해도 동일한 문제 발생
2. ⚠️ Panic 문제 해결이 우선이 아님
3. ⚠️ Genesis output root validation 로직 자체를 재검토 필요

**다음 단계**:
1. Genesis block (block 0)에 대한 validation을 스킵하는 방법 검토
2. 또는 올바른 genesis output root를 사전 계산하여 설정
3. Optimism의 다른 E2E 테스트 확인

**현재 상태**:
- `.devnet` 방식이 틀린 것이 아님
- 문제는 genesis output root validation 로직 자체에 있음

