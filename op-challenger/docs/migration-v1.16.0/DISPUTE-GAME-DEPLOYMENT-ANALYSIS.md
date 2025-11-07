# Tokamak Thanos DisputeGame 배포 방식 분석

**작성일**: 2025-11-07
**목적**: 타노스의 현재 DisputeGame 배포 방식 분석 및 DeployDisputeGame.s.sol 필요성 검토

---

## 📊 요약

- **현황**: 타노스는 이미 DisputeGame을 배포하고 있음 (OPContractsManager 통해)
- **문제**: 독립적인 배포 불가, 유연성 부족, 재배포 어려움
- **해결책**: 옵티미즘의 `DeployDisputeGame.s.sol` 스크립트 추가 권장

---

## 1. 타노스의 현재 DisputeGame 배포 방식

### 1.1 배포 흐름

```
Deploy.s.sol
  └─ deployOpChain()
      └─ DeployOPChain.run()
          └─ opcm.deploy(deployInput)  ← OPContractsManager
              └─ 자동으로 다음을 배포:
                  ├─ DisputeGameFactory (프록시)
                  ├─ FaultDisputeGame
                  ├─ PermissionedDisputeGame
                  ├─ AnchorStateRegistry (프록시)
                  ├─ DelayedWETH (Permissioned용)
                  └─ DelayedWETH (Permissionless용)
```

### 1.2 배포 코드 위치

**파일**: `packages/tokamak/contracts-bedrock/scripts/deploy/DeployOPChain.s.sol:459-462`

```solidity
console.log("DeployOPChain: About to call opcm.deploy()");
vm.broadcast(msg.sender);
IOPContractsManager.DeployOutput memory deployOutput = opcm.deploy(deployInput);
console.log("DeployOPChain: opcm.deploy() completed successfully");
```

### 1.3 배포되는 컨트랙트들

**출력 구조**: `IOPContractsManager.DeployOutput`

```solidity
struct DeployOutput {
    IProxyAdmin opChainProxyAdmin;
    IAddressManager addressManager;
    IL1ERC721Bridge l1ERC721BridgeProxy;
    ISystemConfig systemConfigProxy;
    IOptimismMintableERC20Factory optimismMintableERC20FactoryProxy;
    IL1StandardBridge l1StandardBridgeProxy;
    IL1CrossDomainMessenger l1CrossDomainMessengerProxy;
    IETHLockbox ethLockboxProxy;

    // Fault Proof 관련 컨트랙트들
    IOptimismPortal2 optimismPortalProxy;
    IDisputeGameFactory disputeGameFactoryProxy;           // ✅ 배포됨
    IAnchorStateRegistry anchorStateRegistryProxy;         // ✅ 배포됨
    IFaultDisputeGame faultDisputeGame;                    // ✅ 배포됨
    IPermissionedDisputeGame permissionedDisputeGame;      // ✅ 배포됨
    IDelayedWETH delayedWETHPermissionedGameProxy;        // ✅ 배포됨
    IDelayedWETH delayedWETHPermissionlessGameProxy;      // ✅ 배포됨
    IRAT ratProxy;
}
```

**참조**: `/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/L1/IOPContractsManager.sol:159-178`

### 1.4 배포 입력 파라미터

```solidity
struct DeployInput {
    Roles roles;                          // opChainProxyAdminOwner, systemConfigOwner, etc.
    uint32 basefeeScalar;
    uint32 blobBasefeeScalar;
    uint256 l2ChainId;
    bytes startingAnchorRoot;
    string saltMixer;
    uint64 gasLimit;

    // Dispute Game 파라미터
    GameType disputeGameType;             // CANNON, PERMISSIONED_CANNON, etc.
    Claim disputeAbsolutePrestate;        // Prestate hash
    uint256 disputeMaxGameDepth;          // 게임 최대 깊이 (예: 73)
    uint256 disputeSplitDepth;            // Split 깊이 (예: 30)
    Duration disputeClockExtension;       // 시계 연장 (예: 3시간)
    Duration disputeMaxClockDuration;     // 최대 시간 (예: 3.5일)

    // RAT 파라미터
    bool deployRAT;
    uint256 perTestBondAmount;
    uint256 evidenceSubmissionPeriod;
    uint256 minimumStakingBalance;
    uint256 ratTriggerProbability;
    address ratManager;
}
```

### 1.5 현재 타노스 설정

**파일**: `packages/tokamak/contracts-bedrock/scripts/deploy/DeployOPChain.s.sol:186-198`

```solidity
function startingAnchorRoot() public pure returns (bytes memory) {
    // WARNING: 현재는 permissioned game만 지원
    // 항상 하드코딩된 시작 루트(0xdead)를 사용
    // permissionless game으로 업데이트하려면:
    //   1. 새로운 starting anchor root 설정
    //   2. 새로운 permissioned dispute game 배포 필요

    return abi.encode(ScriptConstants.DEFAULT_OUTPUT_ROOT());
}
```

**의미**:
- 현재는 **Permissioned Game**만 사용
- Permissionless Game 전환 시 재배포 필요

---

## 2. 타노스 방식의 문제점

### 2.1 OPContractsManager 종속성

#### ❌ 문제:
- DisputeGame을 배포하려면 **전체 OP Chain을 배포**해야 함
- `opcm.deploy()`는 모든 컨트랙트를 한 번에 배포
- 게임만 독립적으로 배포할 수 없음

#### 예시:
```bash
# 새로운 게임 타입 추가하고 싶을 때
❌ 현재: 전체 OP Chain 재배포 필요 (매우 비효율적)
✅ 원하는: 게임만 배포하고 Factory에 등록
```

### 2.2 커스텀 게임 배포 불가

#### ❌ 문제:
- 테스트용 게임 (AlphabetVM, 다른 파라미터)을 배포할 수 없음
- 업그레이드된 게임을 독립적으로 배포 불가
- 다양한 GameType을 실험할 수 없음

#### 예시:
```bash
# E2E 테스트에서 AlphabetVM 게임 테스트하고 싶을 때
❌ 현재: 불가능 (OPContractsManager는 Cannon만 지원)
✅ 원하는: AlphabetVM으로 간단한 테스트 게임 배포
```

### 2.3 파라미터 유연성 부족

#### ❌ 문제:
- 게임 파라미터는 `DeployInput`에 고정됨
- 동일한 체인에 다른 파라미터로 여러 게임 배포 불가
- 파라미터 실험이 어려움

#### 예시:
```bash
# 다양한 maxGameDepth로 성능 테스트하고 싶을 때
❌ 현재: DeployInput 수정 → 전체 재배포
✅ 원하는: 각 깊이별로 게임 독립 배포
```

### 2.4 재배포 어려움

#### ❌ 문제:
- 게임 업그레이드 시 전체 시스템 재배포 필요
- 개발 속도 저하
- 테스트 사이클 길어짐

#### 예시:
```bash
# FaultDisputeGame 로직 수정 후 테스트
❌ 현재: opcm.deploy() 재실행 → 10-20분 소요
✅ 원하는: 게임만 재배포 → 1-2분 소요
```

### 2.5 VM 선택 불가

#### ❌ 문제:
- OPContractsManager는 고정된 VM(MIPS) 사용
- Alphabet, Asterisc(RISC-V), CannonNext 등 선택 불가
- 개발/테스트에서 간단한 VM 사용 불가

#### 예시:
```bash
# 빠른 테스트를 위해 AlphabetVM 사용하고 싶을 때
❌ 현재: 불가능
✅ 원하는: VM을 파라미터로 선택 가능
```

---

## 3. 옵티미즘의 DeployDisputeGame.s.sol

### 3.1 파일 위치

**옵티미즘**: `/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployDisputeGame.s.sol`

**타노스**: 현재 없음 ❌

### 3.2 기능

#### Input 구조
```solidity
struct Input {
    string gameKind;                    // "FaultDisputeGame" 또는 "PermissionedDisputeGame"
    GameType gameType;                  // CANNON, PERMISSIONED_CANNON, ALPHABET, etc.
    Claim absolutePrestate;             // Prestate hash
    uint256 maxGameDepth;               // 게임 최대 깊이
    uint256 splitDepth;                 // Split 깊이
    Duration clockExtension;            // 시계 연장
    Duration maxClockDuration;          // 최대 시간
    IDelayedWETH delayedWeth;          // DelayedWETH 컨트랙트
    IAnchorStateRegistry anchorStateRegistry;  // AnchorStateRegistry
    uint256 l2ChainId;                  // L2 Chain ID
    IBigStepper vm;                     // VM 선택 (MIPS, Alphabet, etc.)
}
```

#### Output 구조
```solidity
struct Output {
    IFaultDisputeGame faultDisputeGame;           // 배포된 게임 (FaultDisputeGame)
    IPermissionedDisputeGame permissionedGame;    // 배포된 게임 (PermissionedDisputeGame)
}
```

### 3.3 사용 예시

#### 예시 1: FaultDisputeGame (Permissionless) 배포
```solidity
DeployDisputeGame deployer = new DeployDisputeGame();
DeployDisputeGame.Output memory output = deployer.run(
    DeployDisputeGame.Input({
        gameKind: "FaultDisputeGame",
        gameType: GameTypes.CANNON,
        absolutePrestate: Claim.wrap(0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98),
        maxGameDepth: 73,
        splitDepth: 30,
        clockExtension: Duration.wrap(3 hours),
        maxClockDuration: Duration.wrap(3.5 days),
        delayedWeth: delayedWETH,
        anchorStateRegistry: registry,
        l2ChainId: 901,
        vm: mipsVM  // MIPS64 VM
    })
);

// 배포된 게임을 Factory에 등록
disputeGameFactory.setImplementation(GameTypes.CANNON, output.faultDisputeGame);
```

#### 예시 2: PermissionedDisputeGame 배포
```solidity
DeployDisputeGame deployer = new DeployDisputeGame();
DeployDisputeGame.Output memory output = deployer.run(
    DeployDisputeGame.Input({
        gameKind: "PermissionedDisputeGame",
        gameType: GameTypes.PERMISSIONED_CANNON,
        absolutePrestate: Claim.wrap(0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98),
        maxGameDepth: 73,
        splitDepth: 30,
        clockExtension: Duration.wrap(3 hours),
        maxClockDuration: Duration.wrap(3.5 days),
        delayedWeth: delayedWETH,
        anchorStateRegistry: registry,
        l2ChainId: 901,
        vm: mipsVM
    })
);

// 배포된 게임을 Factory에 등록
disputeGameFactory.setImplementation(GameTypes.PERMISSIONED_CANNON, output.permissionedGame);
```

#### 예시 3: AlphabetVM 테스트 게임 배포
```solidity
// AlphabetVM 배포 (간단한 테스트용 VM)
DeployAlphabetVM alphabetDeployer = new DeployAlphabetVM();
DeployAlphabetVM.Output memory alphabetOutput = alphabetDeployer.run(
    DeployAlphabetVM.Input({
        absolutePrestate: Claim.wrap(0x0000000000000000000000000000000000000000000000000000000000000060)
    })
);

// AlphabetVM으로 테스트 게임 배포
DeployDisputeGame deployer = new DeployDisputeGame();
DeployDisputeGame.Output memory output = deployer.run(
    DeployDisputeGame.Input({
        gameKind: "FaultDisputeGame",
        gameType: GameTypes.ALPHABET,
        absolutePrestate: Claim.wrap(0x0000000000000000000000000000000000000000000000000000000000000060),
        maxGameDepth: 4,  // Alphabet은 깊이가 얕음
        splitDepth: 2,
        clockExtension: Duration.wrap(30 minutes),
        maxClockDuration: Duration.wrap(1 days),
        delayedWeth: delayedWETH,
        anchorStateRegistry: registry,
        l2ChainId: 901,
        vm: alphabetOutput.alphabetVM  // AlphabetVM 사용
    })
);
```

### 3.4 장점 정리

| 기능 | 설명 |
|------|------|
| **독립 배포** | OP Chain 전체 배포 없이 게임만 배포 가능 |
| **파라미터 자유** | 다양한 파라미터로 여러 게임 배포 가능 |
| **VM 선택** | MIPS, Alphabet, Asterisc 등 자유 선택 |
| **빠른 반복** | 게임만 재배포 → 1-2분 (vs 전체 재배포 10-20분) |
| **테스트 용이** | 테스트용 게임을 쉽게 배포 |
| **재사용성** | 다른 프로젝트에서도 재사용 가능 |

---

## 4. 비교표

| 항목 | 타노스 (OPCM 자동) | 옵티미즘 (DeployDisputeGame) |
|------|-------------------|------------------------------|
| **DisputeGame 배포** | ✅ 가능 | ✅ 가능 |
| **독립 배포** | ❌ OP Chain과 함께만 | ✅ 언제든 독립 배포 |
| **커스텀 파라미터** | ⚠️ DeployInput에 고정 | ✅ 자유롭게 설정 |
| **VM 선택** | ❌ 고정 (Cannon만) | ✅ MIPS/Alphabet/Asterisc 선택 |
| **재배포 속도** | ❌ 느림 (10-20분) | ✅ 빠름 (1-2분) |
| **테스트 용이성** | ⚠️ 어려움 | ✅ 쉬움 |
| **유지보수** | ⚠️ 전체 시스템 수정 필요 | ✅ 게임만 수정 |
| **GameType 추가** | ❌ OPCM 재배포 | ✅ 스크립트 재실행 |
| **다중 게임 배포** | ❌ 불가능 | ✅ 여러 게임 독립 배포 |
| **하위 호환성** | ✅ 기존 배포 유지 | ✅ 기존 배포에 영향 없음 |

---

## 5. 실제 사용 시나리오

### 시나리오 1: 새로운 GameType 추가

#### 현재 타노스 방식:
```bash
1. DeployInput에 새 GameType 추가
2. opcm.deploy() 재실행
3. 전체 OP Chain 재배포 (10-20분)
4. 모든 컨트랙트 주소 변경 → 설정 업데이트 필요
```

#### DeployDisputeGame 사용:
```bash
1. DeployDisputeGame.run() 실행 (1-2분)
2. 새 게임만 배포
3. disputeGameFactory.setImplementation() 호출
4. 완료 → 다른 컨트랙트는 영향 없음
```

### 시나리오 2: 게임 파라미터 업그레이드

#### 현재 타노스 방식:
```bash
1. DeployInput의 disputeMaxGameDepth 수정
2. opcm.deploy() 재실행
3. 전체 시스템 재배포
4. 테스트 → 수정 → 재배포 반복 (매우 느림)
```

#### DeployDisputeGame 사용:
```bash
1. 새 maxGameDepth로 DeployDisputeGame.run() 실행
2. 새 게임 배포 (1-2분)
3. 테스트 → 만족하면 Factory 업데이트
4. 빠른 반복 가능
```

### 시나리오 3: E2E 테스트

#### 현재 타노스 방식:
```bash
# E2E 테스트에서 AlphabetVM 게임 테스트하고 싶을 때
❌ 불가능 → OPCM은 Cannon만 지원
⚠️ 전체 시스템을 다시 배포해야 테스트 가능
```

#### DeployDisputeGame 사용:
```bash
# E2E 테스트에서 다양한 VM 테스트
1. AlphabetVM 게임 배포 (빠름)
2. 테스트 실행
3. MIPS VM 게임 배포
4. 테스트 실행
5. 빠른 검증 가능
```

### 시나리오 4: 프로덕션 업그레이드

#### 현재 타노스 방식:
```bash
# Permissioned → Permissionless 전환
1. DeployInput 수정
2. opcm.deploy() 재실행
3. 모든 프록시 주소 업데이트
4. 다운타임 발생 가능
```

#### DeployDisputeGame 사용:
```bash
# Permissioned → Permissionless 전환
1. 새 Permissionless 게임 배포 (기존 시스템 영향 없음)
2. 테스트넷에서 검증
3. disputeGameFactory.setImplementation() 호출 (원자적 업데이트)
4. 다운타임 최소화
```

---

## 6. 타노스에서 DisputeGame을 독립적으로 배포하는 방법 (현재)

### 6.1 Go 코드에서 추가 게임 배포

타노스는 Go 코드에서 추가 DisputeGame 배포를 지원합니다.

**파일**: `op-deployer/pkg/deployer/pipeline/dispute_games.go`

```go
func DeployAdditionalDisputeGames(
    ctx context.Context,
    env *pipeline.Env,
    intent *state.Intent,
    st *state.State,
    chainID common.Hash,
) error {
    lgr := env.Logger.New("stage", "deploy-dispute-games")

    // 1. Oracle 배포 또는 기존 사용
    var oracle common.Address
    if intent.UseStandardPreimageOracle {
        oracle = st.ImplementationsDeployment.PreimageOracleSingletonAddr
    } else {
        // 커스텀 Oracle 배포
        oracle = deployCustomOracle(...)
    }

    // 2. VM 배포
    var vm common.Address
    switch vmType {
    case "alphabet":
        vm = deployAlphabetVM(...)
    case "cannon":
        vm = deployCannonVM(...)
    case "cannon-next":
        vm = deployCannonNextVM(...)
    }

    // 3. DisputeGame 배포
    game := opcm.DeployDisputeGame(ctx, DisputeGameInput{
        GameType: gameType,
        AbsolutePrestate: prestate,
        VM: vm,
        ...
    })

    // 4. Factory에 등록
    opcm.SetDisputeGameImpl(ctx, gameType, game)
}
```

**한계**:
- Go 코드 작성 필요 (Solidity 스크립트보다 복잡)
- op-deployer CLI를 통해서만 실행 가능
- Forge 스크립트 생태계와 분리됨

### 6.2 Solidity 스크립트로 배포 (권장)

**현재 없음** → `DeployDisputeGame.s.sol` 추가 필요

---

## 7. 권장 사항

### 7.1 DeployDisputeGame.s.sol 추가 권장

#### 이유:

1. **기존 배포 유지**
   - OPCM 자동 배포는 그대로 유지 (하위 호환성)
   - 추가로 독립 배포 기능 제공 (확장성)

2. **개발 속도 향상**
   - 게임만 빠르게 재배포 (1-2분)
   - 파라미터 실험 용이
   - 빠른 테스트 사이클

3. **유연성 증가**
   - 다양한 VM 선택 가능
   - 커스텀 파라미터 자유 설정
   - 여러 GameType 독립 관리

4. **표준 준수**
   - 옵티미즘 v1.16.0 표준 방식
   - 마이그레이션 가이드 호환
   - 커뮤니티 지원

5. **E2E 테스트 개선**
   - 테스트용 게임 쉽게 배포
   - AlphabetVM으로 빠른 검증
   - 다양한 시나리오 테스트

### 7.2 구현 계획

#### Phase 1: 기본 스크립트 추가
```bash
1. DeployDisputeGame.s.sol 복사
2. SetDisputeGameImpl.s.sol 복사 (Factory 등록용)
3. 빌드 테스트
```

#### Phase 2: VM 스크립트 추가 (선택)
```bash
4. DeployAlphabetVM.s.sol 복사 (테스트용)
5. DeployAsterisc.s.sol 이미 있음 ✅
```

#### Phase 3: 통합 테스트
```bash
6. E2E 테스트에서 독립 배포 검증
7. 문서 업데이트
```

### 7.3 예상 소요 시간

- **스크립트 복사**: 10분
- **빌드 테스트**: 5분
- **통합 테스트**: 30분
- **문서 작성**: 20분
- **총**: 약 1시간

---

## 8. 결론

### 타노스는 DisputeGame을 배포하지만...

#### ✅ 현재 동작하는 것:
- OPContractsManager가 DisputeGame 자동 배포
- FaultDisputeGame, PermissionedDisputeGame 모두 배포됨
- 기본 기능은 정상 작동

#### ❌ 개선이 필요한 것:
- 독립적인 배포 불가
- 유연성 부족 (파라미터, VM 선택)
- 재배포 속도 느림
- 테스트 어려움

### DeployDisputeGame.s.sol이 필요한 이유

1. ✅ **모듈화**: 게임만 독립적으로 관리
2. ✅ **유연성**: 파라미터와 VM 자유 선택
3. ✅ **효율성**: 빠른 개발/테스트 사이클
4. ✅ **확장성**: 여러 GameType 동시 지원
5. ✅ **표준 준수**: 옵티미즘 v1.16.0 호환

### 최종 권장사항

**DeployDisputeGame.s.sol을 추가하는 것을 강력히 권장합니다.**

이를 통해:
- 기존 OPCM 배포 방식은 유지 (하위 호환성)
- 추가로 독립 배포 기능 제공 (확장성)
- 개발 생산성 대폭 향상
- 옵티미즘 표준 완전 준수

---

## 참고 문서

- [ESSENTIAL-DEPLOYMENT-SCRIPTS-FOR-THANOS.md](./ESSENTIAL-DEPLOYMENT-SCRIPTS-FOR-THANOS.md)
- [DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md](./DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md)
- Optimism Contracts Bedrock: `/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/`

---

**작성자**: Claude Code
**다음 단계**: DeployDisputeGame.s.sol 및 SetDisputeGameImpl.s.sol 추가
