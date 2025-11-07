# 근본 원인 분석: E2E 테스트 AllocType 문제

## 문제 요약

E2E 테스트에서 `AllocTypeMTCannon`과 `AllocTypeMTCannonNext` 두 가지 게임 타입이 **동일한 배포 설정을 공유**하는 문제가 발생합니다. 이는 Deploy.s.sol의 devnet 모드 로직이 하드코딩된 prestate 파일 경로를 사용하기 때문입니다.

## 핵심 근본 원인

### 1. Deploy.s.sol의 prestate 로딩 로직

**위치:** `/packages/tokamak/contracts-bedrock/scripts/Deploy.s.sol:1804-1827`

```solidity
function loadMipsAbsolutePrestate() internal returns (Claim mipsAbsolutePrestate_) {
    if (block.chainid == Chains.LocalDevnet || block.chainid == Chains.GethDevnet) {
        // Fetch the absolute prestate dump
        string memory filePath = string.concat(vm.projectRoot(), "/../../../op-program/bin/prestate-proof.json");
        // ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
        // 🔴 문제: 항상 동일한 파일 경로를 사용함!

        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] = string.concat("[[ -f ", filePath, " ]] && echo \"present\"");
        if (Process.run(commands).length == 0) {
            revert("Cannon prestate dump not found, generate it with `make cannon-prestate` in the monorepo root.");
        }
        commands[2] = string.concat("cat ", filePath, " | jq -r .pre");
        mipsAbsolutePrestate_ = Claim.wrap(abi.decode(Process.run(commands), (bytes32)));
        console.log(
            "[Cannon Dispute Game] Using devnet MIPS Absolute prestate: %s",
            vm.toString(Claim.unwrap(mipsAbsolutePrestate_))
        );
    } else {
        console.log(
            "[Cannon Dispute Game] Using absolute prestate from config: %x", cfg.faultGameAbsolutePrestate()
        );
        mipsAbsolutePrestate_ = Claim.wrap(bytes32(cfg.faultGameAbsolutePrestate()));
    }
}
```

**문제점:**
- `loadMipsAbsolutePrestate()`는 **항상** `op-program/bin/prestate-proof.json` 파일을 읽음
- GameType (CANNON, PERMISSIONED_CANNON 등)에 관계없이 동일한 파일 사용
- MTCannon vs MTCannonNext를 구별하는 로직이 없음

### 2. GameType 배포 설정

**위치:** `/packages/tokamak/contracts-bedrock/scripts/Deploy.s.sol:1889-1909`

```solidity
function setCannonFaultGameImplementation(bool _allowUpgrade) public broadcast {
    console.log("Setting Cannon FaultDisputeGame implementation");
    DisputeGameFactory factory = DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy"));
    DelayedWETH weth = DelayedWETH(mustGetAddress("DelayedWETHProxy"));

    // Set the Cannon FaultDisputeGame implementation in the factory.
    _setFaultGameImplementation({
        _factory: factory,
        _allowUpgrade: _allowUpgrade,
        _params: FaultDisputeGameParams({
            anchorStateRegistry: AnchorStateRegistry(mustGetAddress("AnchorStateRegistryProxy")),
            weth: weth,
            gameType: GameTypes.CANNON,  // GameType 0
            absolutePrestate: loadMipsAbsolutePrestate(),  // 🔴 항상 동일한 함수 호출
            faultVm: IBigStepper(mustGetAddress("Mips")),
            maxGameDepth: cfg.faultGameMaxDepth(),
            maxClockDuration: Duration.wrap(uint64(cfg.faultGameMaxClockDuration()))
        })
    });
}
```

### 3. 현재 존재하는 prestate 파일들

```bash
$ ls -la op-program/bin/prestate-proof*.json
-rwxr-xr-x  prestate-proof.json           # MTCannon용
-rwxr-xr-x  prestate-proof-mt64.json      # ?
-rwxr-xr-x  prestate-proof-mt64Next.json  # MTCannonNext용
```

**문제:** Deploy.s.sol은 `prestate-proof-mt64Next.json`의 존재를 모르고 항상 `prestate-proof.json`만 읽습니다.

## AllocType 시스템 설계

### E2E 테스트의 AllocType 정의

**위치:** `/op-e2e/config/init.go:51-57`

```go
type AllocType string

const (
    AllocTypeAltDA        AllocType = "alt-da"
    AllocTypeMTCannon     AllocType = "mt-cannon"       // 기본 MTCannon (MIPS mt64)
    AllocTypeMTCannonNext AllocType = "mt-cannon-next"  // 새 버전 (MIPS mt64Next)

    DefaultAllocType = AllocTypeMTCannon
)
```

### AllocType별 VM 타입 매핑

**위치:** `/op-e2e/config/init.go.backup:462-470`

```go
func cannonVMType(allocType AllocType) state.VMType {
    if allocType == AllocTypeMTCannonNext {
        return state.VMTypeCannonNext  // mt64Next 버전
    }
    return state.VMTypeCannon  // mt64 기본 버전
}

func IsCannonInDevelopment() bool {
    return cannonVMType(AllocTypeMTCannonNext).MipsVersion() != cannonVMType(AllocTypeMTCannon).MipsVersion()
}
```

### 테스트에서의 PrestateVariant 선택

**위치:** `/op-e2e/system/e2esys/setup.go:394-401`

```go
func (sys *System) PrestateVariant() shared.PrestateVariant {
    switch sys.AllocType() {
    case config.AllocTypeMTCannonNext:
        return shared.MTCannonNextVariant  // ✅ 올바르게 구별함
    default:
        return shared.MTCannonVariant      // ✅ 올바르게 구별함
    }
}
```

## Optimism vs Tokamak 배포 차이점

### Optimism의 배포 방식: op-deployer

**위치:** `/optimism/op-deployer/pkg/deployer/pipeline/pre_state.go:17-74`

Optimism은 **op-deployer**를 사용하여 배포 시 동적으로 prestate를 생성합니다:

```go
func GeneratePreState(ctx context.Context, pEnv *Env, globalIntent *state.Intent, st *state.State, preStateBuilder PreStateBuilder) error {
    lgr := pEnv.Logger.New("stage", "generate-pre-state")

    if preStateBuilder == nil {
        lgr.Debug("preStateBuilder not found, skipping prestate generation")
        return nil
    }
    lgr.Info("preStateBuilder found, proceeding with prestate generation")

    // 각 체인별로 genesis와 rollup 설정을 렌더링
    for _, chain := range globalIntent.Chains {
        genesis, rollup, err := RenderGenesisAndRollup(st, chain.ID, globalIntent)
        // ...
        prestateBuilderOpts = append(prestateBuilderOpts, prestate.WithChainConfig(
            chain.ID.Big().String(),
            bytes.NewReader(rollupJSON),
            bytes.NewReader(genesisJSON),
        ))
    }

    // 동적으로 prestate 빌드
    manifest, err := preStateBuilder.BuildPrestate(ctx, prestateBuilderOpts...)
    if err != nil {
        return fmt.Errorf("failed to build prestate: %w", err)
    }

    st.PrestateManifest = &manifest  // 각 AllocType별로 다른 manifest
    // ...
}
```

**특징:**
- ✅ 배포 시마다 동적으로 prestate 생성
- ✅ 각 AllocType마다 독립적인 PrestateManifest 생성
- ✅ 체인별, 설정별로 다른 prestate 사용 가능

### Tokamak의 배포 방식: Forge Script

**위치:** `/packages/tokamak/contracts-bedrock/scripts/Deploy.s.sol`

Tokamak은 **Forge Script**를 사용하여 정적 배포를 수행합니다:

```solidity
// 단일 배포 스크립트가 모든 설정을 처리
function run() public {
    deployProxies();
    deployImplementations();
    initializeSystemConfig();
    // ...
    setCannonFaultGameImplementation(false);  // 🔴 항상 동일한 prestate 사용
}
```

**특징:**
- ⚠️ **사전에 생성된 단일 .devnet 파일** 사용
- ⚠️ 모든 AllocType이 **동일한 배포 설정 공유**
- ⚠️ 런타임에 AllocType별로 다른 prestate를 선택할 수 없음

## 문제가 발생하는 정확한 흐름

### 1단계: Genesis 생성 (`make devnet-allocs`)

```bash
$ make devnet-allocs
Deploying contracts...
# Deploy.s.sol 실행:
#   - loadMipsAbsolutePrestate() 호출
#   - op-program/bin/prestate-proof.json 읽기
#   - 모든 GameType에 동일한 prestate 적용
Generating .devnet files...
# 결과: 단일 .devnet 파일 세트 생성
```

### 2단계: E2E 초기화 (`op-e2e/config/init.go`)

**위치:** `op-e2e/config/init.go:144-221`

```go
func initFromDevnetFiles() error {
    // .devnet 파일들을 읽어옴
    l1Allocs := readJSON(".devnet/genesis-l1-allocs.json")
    l1Deployments := readJSON(".devnet/genesis-l1-deployments.json")
    deployConfig := readJSON(".devnet/deploy-config.json")
    l2AllocsMap := readJSON(".devnet/genesis-l2-allocs.json")

    // 🔴 문제: 모든 AllocType에 동일하게 적용!
    mtx.Lock()
    for _, allocType := range allocTypes {  // [AltDA, MTCannon, MTCannonNext]
        l1AllocsByType[allocType] = l1Allocs          // ❌ 모두 같음
        l1DeploymentsByType[allocType] = l1Deployments
        deployConfigsByType[allocType] = deployConfig
        l2AllocsByType[allocType] = l2AllocsMap
    }
    mtx.Unlock()
}
```

### 3단계: 테스트 실행

```go
// MTCannon 테스트
cfg := e2esys.DefaultSystemConfig(t, e2esys.WithAllocType(config.AllocTypeMTCannon))
sys1, err := cfg.Start(t)
// → prestate-proof.json의 prestate 사용 ✅

// MTCannonNext 테스트
cfg := e2esys.DefaultSystemConfig(t, e2esys.WithAllocType(config.AllocTypeMTCannonNext))
sys2, err := cfg.Start(t)
// → prestate-proof.json의 prestate 사용 ❌ (prestate-proof-mt64Next.json을 사용해야 함)
```

**결과:** MTCannonNext 테스트가 잘못된 prestate를 사용하여 실패할 수 있음

## 해결 방법 분석

### 방법 1: Prestate 파일 복사 (사용자 제안)

가장 간단한 해결책: 현재 prestate-proof.json이 올바른 버전이라면, 이를 복사만 하면 됨

```bash
# MTCannonNext도 같은 prestate를 사용하도록 복사
cp op-program/bin/prestate-proof.json op-program/bin/prestate-proof-mt64Next.json
```

**장점:**
- ✅ 즉시 적용 가능
- ✅ 코드 수정 불필요
- ✅ 두 AllocType이 현재 같은 MIPS 버전을 사용한다면 올바른 해결책

**단점:**
- ⚠️ MTCannon과 MTCannonNext가 **정말로 다른 MIPS 버전**을 사용해야 한다면 잘못된 해결책
- ⚠️ 근본적인 문제(Deploy.s.sol의 정적 파일 로딩) 해결 안 됨

### 방법 2: Deploy.s.sol 수정 (근본 해결)

Deploy.s.sol에 GameType별 prestate 선택 로직 추가:

```solidity
// 새로운 함수 추가
function loadMipsAbsolutePrestateForGameType(GameType gameType) internal returns (Claim) {
    if (block.chainid == Chains.LocalDevnet || block.chainid == Chains.GethDevnet) {
        string memory fileName;

        // GameType에 따라 다른 파일 선택
        uint32 rawGameType = GameType.unwrap(gameType);
        if (rawGameType == 3) {  // MTCannonNext
            fileName = "prestate-proof-mt64Next.json";
        } else {
            fileName = "prestate-proof.json";
        }

        string memory filePath = string.concat(
            vm.projectRoot(),
            "/../../../op-program/bin/",
            fileName
        );

        // 기존 로직과 동일...
    } else {
        return Claim.wrap(bytes32(cfg.faultGameAbsolutePrestate()));
    }
}

// 기존 함수 수정
function setCannonFaultGameImplementation(bool _allowUpgrade) public broadcast {
    _setFaultGameImplementation({
        // ...
        absolutePrestate: loadMipsAbsolutePrestateForGameType(GameTypes.CANNON),
        // ...
    });
}
```

**장점:**
- ✅ 근본적인 문제 해결
- ✅ 각 GameType이 올바른 prestate 파일 사용
- ✅ 확장 가능 (새로운 GameType 추가 시)

**단점:**
- ⚠️ 코드 수정 필요
- ⚠️ GameType 3이 MTCannonNext인지 확인 필요

### 방법 3: 동적 배포로 전환 (Optimism 방식)

Tokamak도 op-deployer 방식으로 전환:

**장점:**
- ✅ Optimism과 동일한 배포 방식
- ✅ AllocType별로 완전히 독립적인 배포
- ✅ 유지보수 용이

**단점:**
- ❌ 대규모 리팩토링 필요
- ❌ 기존 Forge Script 배포 방식 전면 수정
- ❌ 개발 시간 많이 소요

## 권장 해결 순서

### 즉시 조치 (테스트 통과용)

1. **prestate 파일 상태 확인:**
   ```bash
   # 두 파일의 해시 비교
   jq -r '.pre' op-program/bin/prestate-proof.json
   jq -r '.pre' op-program/bin/prestate-proof-mt64Next.json
   ```

2. **현재 상태 확인:**
   - 두 파일이 **같은 해시**를 가지고 있다면 → 문제 없음
   - 두 파일이 **다른 해시**를 가지고 있다면 → Deploy.s.sol 수정 필요

3. **IsCannonInDevelopment() 확인:**
   ```go
   // op-e2e/config/init.go.backup:469-471
   func IsCannonInDevelopment() bool {
       return cannonVMType(AllocTypeMTCannonNext).MipsVersion() !=
              cannonVMType(AllocTypeMTCannon).MipsVersion()
   }
   ```
   - `true`를 반환하면 → 두 타입이 다른 버전 사용, Deploy.s.sol 수정 필요
   - `false`를 반환하면 → 같은 버전 사용, 파일 복사로 해결 가능

### 단기 조치 (방법 1 또는 2)

**방법 1 적용 조건:**
- `IsCannonInDevelopment() == false`
- 두 prestate 파일의 해시가 같음
- 테스트만 빠르게 통과시키면 됨

**방법 2 적용 조건:**
- `IsCannonInDevelopment() == true`
- 두 AllocType이 실제로 다른 MIPS 버전 사용
- 근본적인 문제 해결 필요

### 장기 조치 (방법 3)

- v1.16.0 마이그레이션 완료 후
- Optimism과의 호환성 향상 필요 시
- 동적 배포 방식으로 전환 검토

## 추가 발견 사항

### GameType 정의

**위치:** `/packages/tokamak/contracts-bedrock/src/dispute/lib/Types.sol:25-34`

```solidity
library GameTypes {
    GameType internal constant CANNON = GameType.wrap(0);
    GameType internal constant PERMISSIONED_CANNON = GameType.wrap(1);
    GameType internal constant ASTERISC = GameType.wrap(2);
    // GameType 3은 정의되지 않음!
}
```

**문제:** Tokamak 코드에는 GameType 3 (MTCannonNext)가 명시적으로 정의되어 있지 않습니다.

### 필요한 확인 사항

1. **GameType 3 정의 추가 필요 여부:**
   ```solidity
   library GameTypes {
       GameType internal constant CANNON = GameType.wrap(0);
       GameType internal constant PERMISSIONED_CANNON = GameType.wrap(1);
       GameType internal constant ASTERISC = GameType.wrap(2);
       GameType internal constant MTCANNON_NEXT = GameType.wrap(3);  // ← 추가 필요?
   }
   ```

2. **Prestate 파일 생성 스크립트 확인:**
   - `make cannon-prestate`가 두 파일 모두 생성하는지 확인
   - 또는 별도 타겟 필요: `make cannon-prestate-mt64next`

3. **E2E 테스트의 기대값 확인:**
   - MTCannon과 MTCannonNext 테스트가 실제로 다른 결과를 기대하는지
   - 아니면 같은 prestate를 사용해도 되는지

## 결론

**근본 원인:**
- Deploy.s.sol의 `loadMipsAbsolutePrestate()`가 하드코딩된 단일 파일 경로를 사용
- AllocType에 따라 다른 prestate 파일을 선택하는 로직 부재

**즉시 해결책:**
- 현재 두 AllocType이 같은 MIPS 버전을 사용한다면 → 파일 복사
- 다른 MIPS 버전을 사용한다면 → Deploy.s.sol에 선택 로직 추가

**장기 해결책:**
- Optimism의 op-deployer 방식 도입 검토

## 참고 파일 위치

1. **Deploy.s.sol prestate 로딩:**
   - `/packages/tokamak/contracts-bedrock/scripts/Deploy.s.sol:1804-1827`

2. **E2E AllocType 시스템:**
   - `/op-e2e/config/init.go:51-75`
   - `/op-e2e/config/init.go:144-221` (initFromDevnetFiles)

3. **Optimism 동적 배포:**
   - `/optimism/op-deployer/pkg/deployer/pipeline/pre_state.go:17-74`

4. **GameType 정의:**
   - `/packages/tokamak/contracts-bedrock/src/dispute/lib/Types.sol:25-34`

5. **Prestate 파일:**
   - `/op-program/bin/prestate-proof.json`
   - `/op-program/bin/prestate-proof-mt64Next.json`
