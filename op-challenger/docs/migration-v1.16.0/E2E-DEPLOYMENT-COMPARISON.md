# E2E 테스트 배포 방식 비교: Optimism vs Tokamak-Thanos

## 개요
이 문서는 Optimism과 Tokamak-Thanos의 E2E 테스트에서 컨트랙트 배포 방식의 차이점을 분석합니다.

## 1. 공통점

### 1.1 배포 스크립트 사용
- **둘 다 `DeployImplementations.s.sol` 사용**
  - Optimism: `/packages/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`
  - Tokamak: `/packages/tokamak/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`
  - **Note**: Tokamak의 `DeployImplementations.s.sol`은 Optimism에서 복사해온 파일

### 1.2 배포 파이프라인
- **둘 다 `op-deployer` 사용**
  - `op-deployer/pkg/deployer/pipeline/implementations.go`
  - `DeployImplementations.Run()` 호출

### 1.3 E2E 초기화 흐름
```
op-e2e/config/init.go
  ↓
initAllocType() - 각 L2 모드별 초기화
  ↓
ApplyPipelineE2E()
  ↓
deployer.ApplyPipeline()
  ↓
DeployImplementations.Run()
```

## 2. 주요 차이점

**핵심**: 배포 스크립트(`DeployImplementations.s.sol`)는 동일하지만, **컨트랙트 구현(`SystemConfig.sol`)이 다름**

### 2.1 SystemConfig 생성자

#### Optimism (올바른 구현)
```solidity
// /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/L1/SystemConfig.sol
constructor() ReinitializableBase(2) {
    Storage.setUint(START_BLOCK_SLOT, type(uint256).max);
    _disableInitializers();  // ✅ Implementation 컨트랙트는 초기화 비활성화
}
```

#### Tokamak-Thanos (문제 있는 구현)
```solidity
// /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/src/L1/SystemConfig.sol
constructor() {
    Storage.setUint(START_BLOCK_SLOT, type(uint256).max);
    initialize({  // ❌ Implementation 배포 시 initialize 호출 → revert 발생
        _owner: address(0xdEaD),
        _basefeeScalar: 0,
        _blobbasefeeScalar: 0,
        _batcherHash: bytes32(0),
        _gasLimit: 1,
        _unsafeBlockSigner: address(0),
        _config: ResourceMetering.ResourceConfig({...}),
        _batchInbox: address(0),
        _addresses: SystemConfig.Addresses({...})
    });
}
```

**문제점:**
- Implementation 컨트랙트는 **초기화되지 않은 상태**로 배포되어야 함
- Proxy를 통해 배포할 때만 `initialize()` 호출
- Tokamak의 생성자는 배포 시점에 `initialize()`를 호출하여 **execution reverted** 발생

### 2.2 에러 발생 지점

#### 에러 로그 분석
```
INFO "=== Deploy deploySuperchainConfigImpl ==="
INFO "=== Deploy deployProtocolVersionsImpl ==="
INFO "=== Deploy deploySystemConfigImpl ==="
WARN Fault addr=0x60fE235DebffaAb783f3aa237C3333C7C1dA973d  <-- SystemConfig 배포 실패
     err="execution reverted"
```

**실패 원인:**
1. `DeployImplementations.s.sol`의 `deploySystemConfigImpl()` 함수 호출
2. `SystemConfig` 생성자 실행
3. 생성자 내부에서 `initialize()` 호출 시도
4. `initialize()`는 proxy를 통해서만 호출되어야 하므로 revert

### 2.3 배포 스크립트 구조

#### Optimism & Tokamak (동일)
```solidity
// DeployImplementations.s.sol (Tokamak이 Optimism에서 복사)
function deploySystemConfigImpl(Output memory _output) private {
    ISystemConfig impl = ISystemConfig(
        DeployUtils.createDeterministic({
            _name: "SystemConfig",
            _args: DeployUtils.encodeConstructor(abi.encodeCall(ISystemConfig.__constructor__, ())),
            _salt: _salt
        })
    );
    vm.label(address(impl), "SystemConfigImpl");
    _output.systemConfigImpl = impl;
}
```

**배포 스크립트 동작:**
1. `ISystemConfig.__constructor__()` - 빈 파라미터로 생성자 호출
2. `SystemConfig` 컨트랙트의 생성자 실행
3. 여기서 Optimism과 Tokamak의 차이 발생!

**Optimism의 경우:**
- `SystemConfig` 생성자가 `_disableInitializers()` 호출
- ✅ 성공적으로 배포됨

**Tokamak의 경우 (수정 전):**
- `SystemConfig` 생성자가 `initialize()` 호출
- `initialize()`는 `onlyInitializing` modifier로 보호됨
- Implementation 배포 시에는 초기화 상태가 아니므로 ❌ revert

## 3. 해결 방법

### 3.1 SystemConfig 생성자 수정
```solidity
// 수정 후
constructor() {
    // NOTE: For implementation contracts, we disable initializers to prevent
    // anyone from initializing the implementation contract itself.
    // The proxy will call initialize() when it's deployed.
    Storage.setUint(START_BLOCK_SLOT, type(uint256).max);
    _disableInitializers();
}
```

### 3.2 수정 이유
- **Implementation 패턴**: Implementation 컨트랙트는 로직만 제공
- **Proxy 패턴**: Proxy가 Implementation을 delegatecall로 호출
- **초기화 분리**:
  - Implementation 배포 시: 초기화 비활성화
  - Proxy 배포 시: `initialize()` 호출하여 상태 설정

## 4. E2E 테스트 배포 흐름 상세

### 4.1 L2 모드별 배포
```go
// op-e2e/config/init.go
allocModes := []genesis.L2AllocsMode{
    genesis.L2AllocsGranite,  // 현재 디버깅용으로 1개만 활성화
    // genesis.L2AllocsInterop,
    // genesis.L2AllocsIsthmus,
    // genesis.L2AllocsHolocene,
    // genesis.L2AllocsFjord,
    // genesis.L2AllocsEcotone,
    // genesis.L2AllocsDelta,
}
```

**특징:**
- 각 L2 모드별로 독립적인 goroutine 실행
- 병렬로 배포 진행
- CREATE2 사용으로 동일한 주소에 배포 (중복 방지)

### 4.2 배포 단계
```
1. DeploySuperchain
   ↓
2. DeployImplementations  <-- 여기서 실패
   ├─ deploySuperchainConfigImpl ✅
   ├─ deployProtocolVersionsImpl ✅
   ├─ deploySystemConfigImpl ❌ (SystemConfig 생성자 문제)
   ├─ deployL1CrossDomainMessengerImpl
   ├─ deployL1ERC721BridgeImpl
   ├─ deployL1StandardBridgeImpl
   ├─ deployOptimismMintableERC20FactoryImpl
   ├─ deployOptimismPortalImpl
   ├─ deployETHLockboxImpl
   ├─ deployDelayedWETHImpl
   ├─ deployPreimageOracleSingleton
   ├─ deployMipsSingleton
   ├─ deployDisputeGameFactoryImpl
   ├─ deployAnchorStateRegistryImpl
   └─ deployOPContractsManager
   ↓
3. DeployOPChain
   ↓
4. L2 Genesis 생성
```

## 5. 다른 컨트랙트 확인 필요 사항

### 5.1 확인 완료
- ✅ `SuperchainConfig`: Optimism도 `_disableInitializers()` 사용
- ✅ `ProtocolVersions`: 문제 없음

### 5.2 추가 확인 필요
다른 Implementation 컨트랙트들도 동일한 패턴을 따르는지 확인:
- `L1CrossDomainMessenger`
- `L1ERC721Bridge`
- `L1StandardBridge`
- `OptimismMintableERC20Factory`
- `OptimismPortal`
- `DisputeGameFactory`
- `AnchorStateRegistry`

## 6. 결론

### 6.1 핵심 문제
**배포 스크립트는 동일하지만, Tokamak의 `SystemConfig.sol` 컨트랙트 구현이 Implementation 배포 패턴을 따르지 않음**

- ✅ `DeployImplementations.s.sol`: Optimism에서 복사 (동일)
- ❌ `SystemConfig.sol`: Tokamak 자체 구현 (생성자에서 `initialize()` 호출)

### 6.2 해결책
**`SystemConfig.sol` 생성자를 Optimism처럼 `_disableInitializers()` 사용하도록 수정**

### 6.3 영향 범위
- E2E 테스트: ✅ 수정 후 정상 동작 예상
- 실제 배포: ✅ Proxy 패턴 사용 시 문제 없음
- 기존 배포: ⚠️ 이미 배포된 컨트랙트는 영향 없음 (새로운 Implementation 배포 시에만 적용)

## 7. 다음 단계

1. ✅ `SystemConfig.sol` 생성자 수정 완료
2. ✅ 컨트랙트 컴파일 완료 (성공)
3. 🔄 E2E 테스트 실행 중
4. ✅ Implementation 컨트랙트 생성자 확인 완료
   - ✅ `SystemConfig`: `_disableInitializers()` 적용
   - ✅ `OptimismPortal2`: `_disableInitializers()` 적용
   - ✅ `ETHLockbox`: Tokamak은 ERC20 ETH 사용으로 불필요 (address(0) 설정)
   - ✅ `DelayedWETH`: 정상 배포 (Fault Proof bond 관리에 필수)
5. ✅ 토큰 아키텍처 분석 완료 (TOKAMAK-TOKEN-ARCHITECTURE.md 참조)
6. ⬜ 전체 E2E 테스트 통과 확인

## 8. Tokamak 토큰 아키텍처 요약

### 8.1 핵심 차이점
- **Optimism**: Native ETH 사용
- **Tokamak**: TON (ERC20 형태) + ERC20 ETH 사용

### 8.2 영향받는 컨트랙트
| 컨트랙트 | Optimism | Tokamak | 비고 |
|---------|----------|---------|------|
| ETHLockbox | ✅ 배포 | ❌ 스킵 | Native ETH lock 용도, Tokamak 불필요 |
| DelayedWETH | ✅ 배포 | ✅ 배포 | Fault Proof bond 관리, 필수 |
| L1StandardBridge | ETH 처리 | TON + ERC20 ETH 처리 | 2단계 전송 로직 |
| OptimismPortal2 | ETH 입금 | 메시지 전달만 | ETH 직접 처리 안 함 |

### 8.3 상세 분석
자세한 내용은 `TOKAMAK-TOKEN-ARCHITECTURE.md` 문서 참조:
- TON, ETH, WETH의 역할 구분
- 각 토큰의 브릿지 메커니즘
- Fault Proof bond token 분석

