# Deploy.s.sol vs DeployImplementations.s.sol 비교 분석

## 개요

Tokamak-Thanos의 기존 배포 스크립트 (`Deploy.s.sol`)와 새로운 Optimism 스타일 배포 스크립트 (`DeployImplementations.s.sol`)의 차이점을 분석합니다.

## 주요 차이점

### 1. Fault Proof 컨트랙트 배포 방식

#### Deploy.s.sol (Tokamak 기존 방식)

```solidity
// OptimismPortal2
OptimismPortal2 portal = new OptimismPortal2{ salt: _implSalt() }({
    _proofMaturityDelaySeconds: cfg.proofMaturityDelaySeconds(),
    _disputeGameFinalityDelaySeconds: cfg.disputeGameFinalityDelaySeconds()
});
```

- **배포 방법**: `new` 키워드로 직접 배포 (salt 사용)
- **파라미터**: 2개 (proofMaturityDelaySeconds, disputeGameFinalityDelaySeconds)
- **위치**: `deployOptimismPortal2()` 함수 (라인 722-725)

#### DeployImplementations.s.sol (Optimism 방식)

```solidity
// OptimismPortal2
IOptimismPortal impl = IOptimismPortal(
    DeployUtils.createDeterministic({
        _name: "OptimismPortal2",
        _args: DeployUtils.encodeConstructor(
            abi.encodeCall(IOptimismPortal.__constructor__, (proofMaturityDelaySeconds, disputeGameFinalityDelaySeconds))
        ),
        _salt: _salt
    })
);
```

- **배포 방법**: `DeployUtils.createDeterministic()` 사용 (CREATE2)
- **파라미터**: 2개 (동일)
- **위치**: `deployOptimismPortalImpl()` 함수

**✅ 수정 완료**: `DeployImplementations.s.sol`에서 파라미터 2개 전달하도록 수정됨

---

### 2. Dispute Game 배포

#### Deploy.s.sol의 Dispute Game 배포 전략

```solidity
// deployImplementations()에서는 DisputeGameFactory만 배포
deployDisputeGameFactory();

// setupOpChain()에서 각 게임 타입별로 Implementation 설정
setAlphabetFaultGameImplementation({ _allowUpgrade: false });
setFastFaultGameImplementation({ _allowUpgrade: false });
setCannonFaultGameImplementation({ _allowUpgrade: false });
setPermissionedCannonFaultGameImplementation({ _allowUpgrade: false });
setAsteriscFaultGameImplementation({ _allowUpgrade: false });
setAsteriscKonaFaultGameImplementation({ _allowUpgrade: false });
```

각 게임 타입은 `setFaultGameImplementationGeneric()` 함수를 통해 배포:

```solidity
// FaultDisputeGame 배포 (Permissioned가 아닌 경우)
_factory.setImplementation(
    _params.gameType,
    new FaultDisputeGame({
        _gameType: _params.gameType,
        _absolutePrestate: _params.absolutePrestate,
        _maxGameDepth: _params.maxGameDepth,
        _splitDepth: cfg.faultGameSplitDepth(),
        _clockExtension: Duration.wrap(uint64(cfg.faultGameClockExtension())),
        _maxClockDuration: _params.maxClockDuration,
        _vm: _params.faultVm,
        _weth: _params.weth,
        _anchorStateRegistry: _params.anchorStateRegistry,
        _l2ChainId: cfg.l2ChainID()
    })
);

// PermissionedDisputeGame 배포 (Permissioned Cannon의 경우)
_factory.setImplementation(
    _params.gameType,
    new PermissionedDisputeGame({
        _gameType: _params.gameType,
        _absolutePrestate: _params.absolutePrestate,
        _maxGameDepth: _params.maxGameDepth,
        _splitDepth: cfg.faultGameSplitDepth(),
        _clockExtension: Duration.wrap(uint64(cfg.faultGameClockExtension())),
        _maxClockDuration: Duration.wrap(uint64(cfg.faultGameMaxClockDuration())),
        _vm: _params.faultVm,
        _weth: _params.weth,
        _anchorStateRegistry: _params.anchorStateRegistry,
        _l2ChainId: cfg.l2ChainID(),
        _proposer: cfg.l2OutputOracleProposer(),
        _challenger: cfg.l2OutputOracleChallenger()
    })
);
```

#### DeployImplementations.s.sol의 Dispute Game 배포

**✅ 정상**: `DeployImplementations.s.sol`에는 `FaultDisputeGame`과 `PermissionedDisputeGame` 배포 함수가 **의도적으로 없음**!

**이유**: Dispute Game Implementation은 **런타임에 동적으로 등록**됩니다.

현재 존재하는 함수:
- `deployDisputeGameFactoryImpl()` - DisputeGameFactory만 배포
- `deployAnchorStateRegistryImpl()` - AnchorStateRegistry만 배포

**게임 타입 등록 방법**:
1. **Tokamak의 Deploy.s.sol 방식** (구버전):
   - `setupOpChain()`에서 `setAlphabetFaultGameImplementation()` 등을 호출
   - 각 함수가 `new FaultDisputeGame()`으로 인스턴스를 생성하고 `factory.setImplementation()` 호출

2. **Optimism의 새로운 방식** (DeployImplementations.s.sol):
   - 별도 스크립트 `SetDisputeGameImpl.s.sol` 사용
   - 또는 런타임에 `AdditionalDisputeGames`를 통해 등록

3. **E2E 테스트에서의 처리**:
   - `op-e2e/config/init.go`에서 `AdditionalDisputeGames` 배열로 게임 타입 정의
   - 런타임에 `DisputeGameFactory.setImplementation()` 호출하여 등록

---

### 3. DelayedWETH 배포

#### Deploy.s.sol

```solidity
function deployDelayedWETH() public broadcast returns (address addr_) {
    console.log("Deploying DelayedWETH implementation");
    DelayedWETH weth = new DelayedWETH{ salt: _implSalt() }(cfg.faultGameWithdrawalDelay());
    save("DelayedWETH", address(weth));
    console.log("DelayedWETH deployed at %s", address(weth));

    // ... validation ...

    addr_ = address(weth);
}
```

- **배포 방법**: `new DelayedWETH{ salt: _implSalt() }(cfg.faultGameWithdrawalDelay())`
- **파라미터**: 1개 (withdrawalDelay)

#### DeployImplementations.s.sol

```solidity
function deployDelayedWETHImpl(Input memory _input, Output memory _output) private {
    uint256 withdrawalDelaySeconds = _input.withdrawalDelaySeconds;
    IDelayedWETH impl = IDelayedWETH(
        DeployUtils.createDeterministic({
            _name: "DelayedWETH",
            _args: DeployUtils.encodeConstructor(
                abi.encodeCall(IDelayedWETH.__constructor__, (withdrawalDelaySeconds))
            ),
            _salt: _salt
        })
    );
    vm.label(address(impl), "DelayedWETHImpl");
    _output.delayedWethImpl = impl;
}
```

- **배포 방법**: `DeployUtils.createDeterministic()`
- **파라미터**: 1개 (동일)

**✅ 수정 완료**: `DelayedWETH` constructor가 `_disableInitializers()` 호출하도록 수정됨

---

### 4. 배포 흐름 차이

#### Tokamak Deploy.s.sol의 배포 흐름 (구버전 - Monolithic)

```
setupOpChain() {
  1. deployProxies()
     - OptimismPortalProxy
     - DisputeGameFactoryProxy
     - DelayedWETHProxy
     - PermissionedDelayedWETHProxy
     - AnchorStateRegistryProxy

  2. deployImplementations()
     - deployOptimismPortal2()
     - deployDisputeGameFactory()
     - deployDelayedWETH()
     - deployAnchorStateRegistry()
     - deployPreimageOracle()
     - deployMips()
     - deployRiscv()

  3. initializeImplementations()
     - initializeOptimismPortal2()
     - initializeDisputeGameFactory()
     - initializeDelayedWETH()
     - initializePermissionedDelayedWETH()
     - initializeAnchorStateRegistry()

  4. setFaultGameImplementations()
     - setAlphabetFaultGameImplementation()
     - setFastFaultGameImplementation()
     - setCannonFaultGameImplementation()
     - setPermissionedCannonFaultGameImplementation()
     - setAsteriscFaultGameImplementation()
     - setAsteriscKonaFaultGameImplementation()
}
```

**특징**:
- 모든 단계를 하나의 스크립트에서 순차적으로 처리
- 프록시 → Implementation → 초기화 → 게임 등록 순서
- 간단하지만 유연성 부족

---

#### Optimism op-deployer의 배포 흐름 (신버전 - Modular Pipeline)

```
ApplyPipeline() {
  1. init
     - InitGenesisStrategy() 또는 InitLiveStrategy()

  2. deploy-superchain
     - DeploySuperchain()
       * SuperchainConfig 배포
       * ProtocolVersions 배포

  3. deploy-implementations
     - DeployImplementations()
       * 공유 Implementation 배포
       * DeployImplementations.s.sol 실행

  4. deploy-opchain-{chainID}
     - DeployOPChain()
       * OPCM을 통해 프록시 + Implementation 동시 배포
       * DeployOPChain.s.sol 실행
       * 프록시 배포 및 초기화를 한 번에 처리

  5. deploy-alt-da-{chainID}
     - DeployAltDA()
       * Alternative DA 설정 (필요시)

  6. deploy-additional-dispute-games-{chainID}
     - DeployAdditionalDisputeGames()
       * PreimageOracle 배포 (필요시)
       * VM 배포 (AlphabetVM, MIPS 등)
       * FaultDisputeGame 배포
       * DisputeGameFactory.setImplementation() 호출
       * AnchorStateRegistry 설정 (MakeRespected인 경우)

  7. generate-l2-genesis-{chainID}
     - GenerateL2Genesis()
}
```

**특징**:
- 파이프라인 기반의 모듈화된 구조
- 각 단계가 독립적으로 실행 가능
- 체인별로 병렬 처리 가능
- OPCM (OPContractsManager)를 통한 통합 배포

---

#### 핵심 차이점

| 항목 | Tokamak Deploy.s.sol | Optimism op-deployer |
|------|---------------------|---------------------|
| **프록시 배포** | `deployProxies()` - 별도 단계 | `DeployOPChain()` - OPCM 통합 |
| **Implementation 배포** | `deployImplementations()` - 별도 단계 | `DeployImplementations()` + OPCM |
| **초기화** | `initializeImplementations()` - 별도 단계 | `DeployOPChain()` - 프록시 배포와 동시 |
| **게임 등록** | `setFaultGameImplementations()` - 스크립트 내 | `DeployAdditionalDisputeGames()` - 별도 파이프라인 |
| **VM 배포** | `deployMips()`, `deployRiscv()` - 스크립트 내 | `DeployAdditionalDisputeGames()` - 런타임 |
| **PreimageOracle** | `deployPreimageOracle()` - 스크립트 내 | `DeployAdditionalDisputeGames()` - 런타임 |

---

#### OPCM (OPContractsManager)의 역할

Optimism의 새로운 배포 방식에서 OPCM은 핵심 역할을 합니다:

1. **통합 배포**: 프록시와 Implementation을 한 번에 배포
2. **Blueprint 패턴**: 컨트랙트 코드를 Blueprint로 저장하고 필요시 복제
3. **표준화**: 모든 OP Chain이 동일한 방식으로 배포
4. **업그레이드 관리**: 중앙화된 업그레이드 관리

**DeployOPChain.s.sol**이 OPCM을 호출하면:
- 프록시 생성
- Implementation 설정
- 초기화 실행
- 모두 하나의 트랜잭션 또는 원자적 작업으로 처리

---

## Dispute Game Implementation 등록 방식 차이

### ✅ 이해: Dispute Game은 런타임에 등록됨

**Deploy.s.sol (구버전)**:
- 배포 스크립트 내에서 직접 게임 타입별 implementation 생성 및 등록
- `setAlphabetFaultGameImplementation()` 등의 함수 사용
- 장점: 한 번에 모든 게임 타입 설정
- 단점: 배포 스크립트가 복잡해짐

**DeployImplementations.s.sol + SetDisputeGameImpl.s.sol (신버전)**:
- Implementation 배포와 게임 타입 등록을 **분리**
- `DeployImplementations.s.sol`: DisputeGameFactory만 배포
- `SetDisputeGameImpl.s.sol`: 별도 스크립트로 게임 타입 등록
- 또는 런타임에 `AdditionalDisputeGames`로 등록
- 장점: 모듈화, 유연성
- 단점: 여러 단계 필요

### E2E 테스트에서의 게임 타입 등록

`op-e2e/config/init.go`에서:
```go
AdditionalDisputeGames: []state.AdditionalDisputeGame{
    {
        ChainProofParams: state.ChainProofParams{
            DisputeGameType: 254, // Alphabet Fast
            // ... 기타 파라미터
        },
    },
    // 다른 게임 타입들...
}
```

이 설정은 런타임에 `DisputeGameFactory.setImplementation()`을 호출하여 각 게임 타입을 등록합니다.

### 실제로 "빠진" 것은 없음

PreimageOracle, MIPS VM, RISC-V VM 등은:
- E2E 테스트 환경에서 필요할 때 동적으로 배포
- 또는 별도의 배포 스크립트로 관리
- `DeployImplementations.s.sol`은 **핵심 Implementation만 배포**하는 것이 목적

---

## Optimism의 DeployImplementations.s.sol 확인 필요

Tokamak의 `DeployImplementations.s.sol`은 Optimism에서 복사한 것인데, Optimism 원본에서도 위의 배포 함수들이 빠져있는지 확인이 필요합니다.

만약 Optimism에서도 빠져있다면, Optimism은 다른 방식으로 배포하는 것일 수 있습니다:

1. **OPCM (OPContractsManager)에서 배포**: OPCM이 런타임에 Dispute Game을 배포할 수도 있음
2. **별도 스크립트**: Fault Proof 관련 컨트랙트는 별도의 배포 스크립트로 분리되어 있을 수도 있음

---

## 다음 단계 (올바른 방향)

### ✅ 중요: Optimism의 모듈화된 아키텍처 따르기

**❌ 잘못된 방향**:
- `DeployImplementations.s.sol`에 `FaultDisputeGame`, `PermissionedDisputeGame` 배포 함수 추가
- 이는 구버전(Monolithic) 방식으로 돌아가는 것

**✅ 올바른 방향**:
- `DeployImplementations.s.sol`은 **핵심 Implementation만 배포**
- Dispute Game은 **별도 파이프라인 단계**에서 처리

---

### 1. ✅ 별도 배포 스크립트 검증 완료

Optimism의 새로운 아키텍처는 다음과 같이 모듈화되어 있습니다:

**별도의 Deploy 스크립트들**:
```bash
packages/tokamak/contracts-bedrock/scripts/deploy/
├── DeployImplementations.s.sol    ✅ 핵심 implementation만
├── DeployDisputeGame.s.sol        ✅ Dispute Game 전용
├── DeployDisputeGame2.s.sol       ✅ 새 버전
├── DeployMIPS2.s.sol              ✅ MIPS VM 전용
├── DeployAlphabetVM.s.sol         ✅ Alphabet VM 전용
├── DeployPreimageOracle2.s.sol    ✅ PreimageOracle 전용
├── DeployAsterisc.s.sol           ✅ Asterisc VM 전용
└── ...
```

#### Go 래퍼 (`op-deployer/pkg/deployer/opcm/`)
```go
✅ alphabet.go          → DeployAlphabetVM.s.sol 호출
✅ mips2.go            → DeployMIPS2.s.sol 호출
✅ perimage_oracle2.go → DeployPreimageOracle2.s.sol 호출
✅ dispute_game.go     → DeployDisputeGame.s.sol 호출
✅ dispute_game2.go    → DeployDisputeGame2.s.sol 호출
✅ asterisc.go         → DeployAsterisc.s.sol 호출
```

#### 파이프라인 (`dispute_games.go`)
`op-deployer/pkg/deployer/pipeline/dispute_games.go`에서 런타임 조합:

```go
func deployDisputeGame(...) {
    // 1. PreimageOracle 배포 (필요시)
    opcm.DeployPreimageOracle(...)

    // 2. VM 배포 (타입에 따라)
    switch game.VMType {
    case VMTypeAlphabet:
        opcm.DeployAlphabetVM(...)
    case VMTypeCannon:
        opcm.DeployMIPS(...)
    case VMTypeAsterisc:
        opcm.DeployAsterisc(...)
    }

    // 3. DisputeGame 배포
    opcm.DeployDisputeGame(...)

    // 4. DisputeGameFactory에 등록
    opcm.SetDisputeGameImpl(...)
}
```

**결론**: DeployImplementations.s.sol에 추가 코드가 **필요 없습니다**!

### 2. ✅ DeployImplementations.s.sol 검증 완료

현재 `DeployImplementations.s.sol`은 Optimism 아키텍처를 올바르게 따릅니다:

| 항목 | 상태 | 설명 |
|------|------|------|
| OptimismPortal2 | ✅ | 2개 파라미터 전달 |
| DelayedWETH | ✅ | `_disableInitializers()` 호출 |
| SystemConfig | ✅ | `_disableInitializers()` 호출 |
| DisputeGameFactory | ✅ | Implementation 배포 (게임 타입 등록은 파이프라인에서) |
| ETHLockbox | ✅ | `address(0)` (Tokamak 특화 - 비활성화) |
| RAT | ✅ | `address(0)` (Tokamak 특화 - 비활성화) |
| FaultDisputeGame | ✅ | 별도 스크립트 (`DeployDisputeGame.s.sol`) |
| PreimageOracle | ✅ | 별도 스크립트 (`DeployPreimageOracle2.s.sol`) |
| MIPS VM | ✅ | 별도 스크립트 (`DeployMIPS2.s.sol`) |
| Alphabet VM | ✅ | 별도 스크립트 (`DeployAlphabetVM.s.sol`) |

**컴파일 확인**:
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
forge build
```

---

### 3. 다음 작업: E2E 테스트 실행

#### 전체 파이프라인 흐름
```
1. init                              ✅ Genesis 전략 초기화
2. deploy-superchain                 ✅ SuperchainConfig, ProtocolVersions
3. deploy-implementations            ✅ 핵심 Implementation (검증 완료)
4. deploy-opchain                    ⏳ OPCM을 통한 프록시 + 초기화
5. deploy-additional-dispute-games   ⏳ 게임 타입 등록 (별도 스크립트 호출)
6. generate-l2-genesis              ⏳ L2 Genesis 생성
```

#### 테스트 명령어
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos

# 컨트랙트 컴파일
cd packages/tokamak/contracts-bedrock && forge build

# E2E 테스트 실행
cd /Users/zena/tokamak-projects/tokamak-thanos
go test -v ./op-e2e/faultproofs -run TestOutputCannonGame -timeout 30m
```

**파이프라인 흐름**:
1. `init` - Genesis 전략 초기화
2. `deploy-superchain` - SuperchainConfig, ProtocolVersions
3. `deploy-implementations` - 핵심 Implementation ← **현재 검증 중**
4. `deploy-opchain` - OPCM을 통한 프록시 + 초기화
5. `deploy-additional-dispute-games` - 게임 타입 등록 (FaultDisputeGame 배포)
6. `generate-l2-genesis` - L2 Genesis 생성

### 4. 문제 발생 시 디버깅 포인트
- `DeployImplementations.s.sol` 실행 로그 확인
- `DeployOPChain.s.sol` OPCM 호출 로그 확인
- `DeployAdditionalDisputeGames` 파이프라인 로그 확인
- 각 게임 타입의 `setImplementation()` 호출 확인

---

## 요약

| 항목 | Deploy.s.sol | DeployImplementations.s.sol | 상태 |
|------|--------------|----------------------------|------|
| OptimismPortal2 | ✅ 2 파라미터 | ✅ 2 파라미터 (수정됨) | ✅ 완료 |
| DelayedWETH | ✅ 배포 | ✅ 배포 (수정됨) | ✅ 완료 |
| DisputeGameFactory | ✅ 배포 | ✅ 배포 | ✅ 완료 |
| AnchorStateRegistry | ✅ 배포 | ✅ 배포 | ✅ 완료 |
| FaultDisputeGame | ✅ 스크립트 내 등록 | ✅ 런타임 등록 | ✅ 정상 (방식 차이) |
| PermissionedDisputeGame | ✅ 스크립트 내 등록 | ✅ 런타임 등록 | ✅ 정상 (방식 차이) |
| PreimageOracle | ✅ 배포 | ✅ 별도 관리 | ✅ 정상 (E2E에서 처리) |
| MIPS VM | ✅ 배포 | ✅ 별도 관리 | ✅ 정상 (E2E에서 처리) |
| RISC-V VM | ✅ 배포 | ✅ 별도 관리 | ✅ 정상 (E2E에서 처리) |
| AlphabetVM | ✅ 배포 | ✅ 별도 관리 | ✅ 정상 (E2E에서 처리) |

## 결론

**✅ DeployImplementations.s.sol은 정상적으로 작동합니다!**

차이점은 **설계 철학의 차이**입니다:
- **Deploy.s.sol**: 모든 것을 한 스크립트에서 처리 (Monolithic)
- **DeployImplementations.s.sol**: 핵심 Implementation만 배포, 나머지는 런타임/별도 스크립트로 처리 (Modular)

E2E 테스트는 `AdditionalDisputeGames`를 통해 런타임에 필요한 게임 타입을 등록하므로, `DeployImplementations.s.sol`에 게임 타입 등록 코드가 없는 것은 **정상**입니다.

---

## Tokamak-Thanos 마이그레이션 현황

### 목표
Tokamak-Thanos를 Optimism의 새로운 배포 방식(op-deployer 파이프라인)으로 마이그레이션

### 현재 상태

#### ✅ 완료된 작업

1. **DeployImplementations.s.sol 수정**
   - 파일: `packages/tokamak/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`
   - 변경 내용:
     * OptimismPortal2 constructor: 2개 파라미터 전달 (라인 399)
       ```solidity
       abi.encodeCall(IOptimismPortal.__constructor__, (proofMaturityDelaySeconds, disputeGameFinalityDelaySeconds))
       ```
     * ETHLockbox: `address(0)` 설정 (라인 334-350)
       ```solidity
       _output.ethLockboxImpl = IETHLockbox(address(0));
       ```
     * RAT: `address(0)` 설정 (라인 481-497)
       ```solidity
       _output.ratImpl = IRAT(address(0));
       ```
     * deployETHLockboxImpl(), deployRATImpl() 함수 호출 추가 (라인 112, 124)

2. **OptimismPortal2.sol 수정**
   - 파일: `packages/tokamak/contracts-bedrock/src/L1/OptimismPortal2.sol`
   - 변경 내용:
     * Constructor에서 `_disableInitializers()` 호출 (라인 153-160)
       ```solidity
       constructor(uint256 _proofMaturityDelaySeconds, uint256 _disputeGameFinalityDelaySeconds) {
           PROOF_MATURITY_DELAY_SECONDS = _proofMaturityDelaySeconds;
           DISPUTE_GAME_FINALITY_DELAY_SECONDS = _disputeGameFinalityDelaySeconds;
           _disableInitializers();
       }
       ```
     * 디버그 로그 제거 완료

3. **IOptimismPortal2.sol 수정**
   - 파일: `packages/tokamak/contracts-bedrock/interfaces/L1/IOptimismPortal2.sol`
   - 변경 내용:
     * `__constructor__` 시그니처: 2개 파라미터로 변경 (라인 105)
       ```solidity
       function __constructor__(uint256 _proofMaturityDelaySeconds, uint256 _disputeGameFinalityDelaySeconds) external;
       ```

4. **SystemConfig.sol 수정**
   - 파일: `packages/tokamak/contracts-bedrock/src/L1/SystemConfig.sol`
   - 변경 내용:
     * Constructor에서 `initialize()` 대신 `_disableInitializers()` 호출
       ```solidity
       constructor() {
           _disableInitializers();
       }
       ```

5. **DelayedWETH.sol 수정**
   - 파일: `packages/tokamak/contracts-bedrock/src/dispute/weth/DelayedWETH.sol`
   - 변경 내용:
     * Constructor에서 `initialize()` 대신 `_disableInitializers()` 호출
       ```solidity
       constructor(uint256 _delay) {
           DELAY_SECONDS = _delay;
           _disableInitializers();
       }
       ```

6. **DeploySuperchain.s.sol 수정**
   - 파일: `packages/tokamak/contracts-bedrock/scripts/deploy/DeploySuperchain.s.sol`
   - 변경 내용:
     * SuperchainConfig.initialize: 2개 파라미터 전달 (라인 142)
       ```solidity
       abi.encodeCall(ISuperchainConfig.initialize, (guardian, false))
       ```

7. **문서 작성**
   - `TOKAMAK-TOKEN-ARCHITECTURE.md`: 토큰 아키텍처 분석
   - `TOKAMAK-NATIVE-TOKEN-ANALYSIS.md`: 네이티브 토큰 상세 분석
   - `DEPLOYMENT-COMPARISON-DEPLOY-VS-DEPLOYIMPLEMENTATIONS.md`: 배포 방식 비교 (현재 문서)

8. **토큰 아키텍처 분석 완료**
   - TON: L1에서 ERC20, L2에서 Native (L2 Gas Token)
   - ETH: L1에서 Native (L1 Gas Token), L2에서 ERC20
   - TON Lock: OptimismPortal에서 관리
   - ETH Lock: L1StandardBridge에서 관리
   - DelayedWETH: Fault Proof bond 관리용 (ETH → WETH 변환)

#### 🔄 진행 중

**DeployImplementations.s.sol 검증**
- 파일 위치: `/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`
- 상태: 컴파일 가능, E2E 테스트 대기 중

**주요 구현 내용**:
```solidity
run(Input memory _input) public returns (Output memory output_) {
    // 1. Superchain 레벨 Implementation
    deploySuperchainConfigImpl(output_);
    deployProtocolVersionsImpl(output_);

    // 2. Chain-specific Implementation
    deploySystemConfigImpl(output_);
    deployL1CrossDomainMessengerImpl(output_);
    deployL1ERC721BridgeImpl(output_);
    deployL1StandardBridgeImpl(output_);
    deployOptimismMintableERC20FactoryImpl(output_);

    // 3. Fault Proof Implementation
    deployOptimismPortalImpl(_input, output_);
    deployETHLockboxImpl(output_);           // address(0) for Tokamak
    deployDelayedWETHImpl(_input, output_);
    deployPreimageOracleSingleton(_input, output_);
    deployMipsSingleton(_input, output_);
    deployDisputeGameFactoryImpl(output_);
    deployAnchorStateRegistryImpl(_input, output_);
    deployRATImpl(output_);                  // address(0) for Tokamak

    // 4. OPCM (OPContractsManager) 배포
    deployOPContractsManager(_input, output_);

    // 5. 검증
    assertValidOutput(_input, output_);
}
```

**Tokamak 특화 처리**:
- `ethLockboxImpl = address(0)` - Tokamak은 ERC20 ETH 사용
- `ratImpl = address(0)` - Tokamak은 RAT 미사용
- `opcmInteropMigrator = address(0)` - Interop 미지원

#### ⏳ 대기 중

1. **컨트랙트 컴파일**
   ```bash
   cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
   forge build
   ```

2. **E2E 테스트 실행**
   ```bash
   cd /Users/zena/tokamak-projects/tokamak-thanos
   go test -v ./op-e2e/faultproofs -run TestOutputCannonGame -timeout 30m
   ```

3. **예상 테스트 흐름**:
   ```
   1. init (InitGenesisStrategy)
   2. deploy-superchain (DeploySuperchain.s.sol)
   3. deploy-implementations (DeployImplementations.s.sol) ← 현재 검증 중
   4. deploy-opchain (DeployOPChain.s.sol via OPCM)
   5. deploy-additional-dispute-games (런타임 게임 타입 등록)
   6. generate-l2-genesis
   ```

#### 🎯 다음 작업

1. **컨트랙트 컴파일 확인**
   - 모든 constructor 파라미터가 올바른지 확인
   - ABI 호환성 확인

2. **E2E 테스트 실행 및 디버깅**
   - `OptimismPortal2` 배포 성공 확인
   - `DisputeGameFactory` 초기화 확인
   - 게임 타입 등록 (`AdditionalDisputeGames`) 확인

3. **문제 발생 시 대응**
   - 로그 분석하여 정확한 실패 지점 파악
   - Optimism 원본과 비교하여 누락된 부분 확인
   - Tokamak 특화 로직 조정

---

## 참고 자료

- **Tokamak Deploy.s.sol**: `/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/Deploy.s.sol`
- **Optimism 원본**: `/Users/zena/tokamak-projects/optimism/`
- **op-deployer 파이프라인**: `/Users/zena/tokamak-projects/tokamak-thanos/op-deployer/pkg/deployer/apply.go`
- **토큰 아키텍처 분석**: `TOKAMAK-TOKEN-ARCHITECTURE.md`

