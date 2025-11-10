# 다음 세션 작업 가이드

**작성일**: 2025-11-10
**상태**: L1CrossDomainMessenger, L1ERC721Bridge 완료, L1StandardBridge 작업 필요

---

## 📊 현재 진행 상황

### ✅ 완료된 작업

#### 1. L1CrossDomainMessenger (완료)
- ✅ ProxyAdminOwnedBase, ReinitializableBase 상속 추가
- ✅ Constructor: `_disableInitializers()` 패턴 적용
- ✅ Storage layout spacer 추가
- ✅ Initialize 시그니처: 3개 파라미터 유지 (Tokamak 호환성)
- ✅ Tokamak Native Token 기능 100% 보존:
  - `sendNativeTokenMessage()`
  - `_sendNativeTokenMessage()`
  - `onApprove()`
  - `relayMessage()`
  - `nativeTokenAddress()`
- ✅ depositTransaction 6-parameter 패턴 유지 (**중요!**)
- ✅ 컴파일 성공
- ✅ ChainAssertions 통과

#### 2. DisputeGameFactory (완료)
- ✅ Constructor: `initialize(address(0))` → `_disableInitializers()`
- ✅ 컴파일 성공
- ✅ ChainAssertions 통과

#### 3. L1ERC721Bridge (완료)
- ✅ ProxyAdminOwnedBase, ReinitializableBase 상속 추가
- ✅ Constructor: `_disableInitializers()` 패턴 적용
- ✅ Storage layout spacer 추가
- ✅ Initialize 시그니처: SystemConfig 사용
- ✅ `paused()`, `superchainConfig()` 함수 업데이트
- ✅ Deploy.s.sol, Initializable.t.sol 업데이트
- ✅ 컴파일 성공
- ✅ ChainAssertions 통과

#### 4. 인터페이스 업데이트 (완료)
- ✅ ISystemConfig: `nativeTokenAddress()` 추가
- ✅ IOptimismPortal2: 6-parameter `depositTransaction()` 추가

#### 5. 문서화 (완료)
- ✅ L1CROSSDOMAINMESSENGER-INTEGRATION-GUIDE.md 작성

---

## ⏳ 다음 세션 작업: L1StandardBridge

### 현재 에러
```
ERROR: DeployUtils: storage value is not 0xff at the given slot and offset
Location: checkL1StandardBridgeImpl
```

### 필요한 작업

#### 1단계: Optimism 원본 확인
```bash
cat /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/L1/L1StandardBridge.sol
```

#### 2단계: L1StandardBridge 업데이트

**예상 변경사항** (L1CrossDomainMessenger, L1ERC721Bridge와 동일한 패턴):

```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Contracts
import { ProxyAdminOwnedBase } from "src/L1/ProxyAdminOwnedBase.sol";
import { ReinitializableBase } from "src/universal/ReinitializableBase.sol";
import { StandardBridge } from "src/universal/StandardBridge.sol";

// Interfaces
import { ISuperchainConfig } from "interfaces/L1/ISuperchainConfig.sol";
import { ISystemConfig } from "interfaces/L1/ISystemConfig.sol";
import { ISemver } from "src/universal/ISemver.sol";

contract L1StandardBridge is
    StandardBridge,
    ProxyAdminOwnedBase,
    ReinitializableBase,
    ISemver
{
    /// @custom:legacy
    /// @custom:spacer superchainConfig
    address private spacer_60_0_20;  // 또는 다른 slot 번호

    /// @notice Semantic version.
    string public constant version = "2.x.0-tokamak";

    /// @notice Contract of the SystemConfig.
    ISystemConfig public systemConfig;

    /// @notice Constructs the L1StandardBridge contract.
    constructor() StandardBridge() ReinitializableBase(2) {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    function initialize(
        ICrossDomainMessenger _messenger,
        ISystemConfig _systemConfig
    )
        external
        reinitializer(initVersion())
    {
        _assertOnlyProxyAdminOrProxyAdminOwner();
        systemConfig = _systemConfig;
        __StandardBridge_init({
            _messenger: _messenger,
            _otherBridge: StandardBridge(payable(Predeploys.L2_STANDARD_BRIDGE))
        });
    }

    function paused() public view override returns (bool) {
        return superchainConfig().paused();
    }

    function superchainConfig() public view returns (ISuperchainConfig) {
        return systemConfig.superchainConfig();
    }

    // Tokamak Native Token 관련 함수들 유지
    // ...
}
```

#### 3단계: Deploy 스크립트 업데이트

**파일**: `scripts/Deploy.s.sol`

찾을 위치:
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
grep -n "L1StandardBridge.*initialize" scripts/Deploy.s.sol
```

업데이트:
```solidity
// OLD
L1StandardBridge.initialize(
    L1CrossDomainMessenger(payable(l1CrossDomainMessengerProxy)),
    SuperchainConfig(superchainConfigProxy)
)

// NEW
L1StandardBridge.initialize(
    L1CrossDomainMessenger(payable(l1CrossDomainMessengerProxy)),
    ISystemConfig(systemConfigProxy)
)
```

#### 4단계: 테스트 파일 업데이트

**파일**: `test/vendor/Initializable.t.sol`

찾을 위치:
```bash
grep -n "l1StandardBridge.initialize" test/vendor/Initializable.t.sol
```

업데이트:
```solidity
// OLD
abi.encodeCall(l1StandardBridge.initialize, (l1CrossDomainMessenger, superchainConfig))

// NEW
abi.encodeCall(l1StandardBridge.initialize, (l1CrossDomainMessenger, ISystemConfig(address(systemConfig))))
```

#### 5단계: 컴파일 및 테스트
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
forge build

# 성공하면
cd /Users/zena/tokamak-projects/tokamak-thanos
go test -v ./op-e2e/faultproofs -run TestOutputCannonGame -timeout 30m
```

---

## 🔑 핵심 패턴 (재사용 가능)

### 모든 L1 컨트랙트에 적용할 패턴:

1. **상속 추가**:
   ```solidity
   contract ContractName is
       BaseContract,
       ProxyAdminOwnedBase,
       ReinitializableBase,
       ISemver
   ```

2. **Storage spacer 추가**:
   ```solidity
   /// @custom:legacy
   /// @custom:spacer superchainConfig
   address private spacer_XX_0_20;  // XX는 slot 번호

   // 기존 storage 변수들...

   ISystemConfig public systemConfig;  // 새로운 위치
   ```

3. **Constructor 패턴**:
   ```solidity
   constructor() BaseContract() ReinitializableBase(2) {
       _disableInitializers();
   }
   ```

4. **Initialize 패턴**:
   ```solidity
   function initialize(..., ISystemConfig _systemConfig)
       external
       reinitializer(initVersion())
   {
       _assertOnlyProxyAdminOrProxyAdminOwner();
       systemConfig = _systemConfig;
       // ... initialization logic
   }
   ```

5. **paused() 함수**:
   ```solidity
   function paused() public view override returns (bool) {
       return superchainConfig().paused();
   }

   function superchainConfig() public view returns (ISuperchainConfig) {
       return systemConfig.superchainConfig();
   }
   ```

---

## 🚨 주의사항

### Tokamak Native Token 기능 보존

L1StandardBridge에서 **반드시 확인해야 할 함수들**:

1. **depositERC20 관련**:
   - Native Token 전송 로직
   - ERC20 approve 패턴

2. **_initiateBridgeERC20**:
   - Native Token 처리
   - Portal 호출 시 파라미터

3. **finalizeBridgeERC20**:
   - Native Token 수신
   - ERC20 transfer 로직

### Critical: depositTransaction 파라미터

L1CrossDomainMessenger에서 확인했듯이, **6-parameter 패턴 유지**:
```solidity
portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
//                              ^^^^^^  ^^^^^^
//                              mint    value (BOTH REQUIRED!)
```

---

## 📝 체크리스트

작업 시작 전 확인:

- [ ] Optimism 원본 L1StandardBridge 확인
- [ ] Tokamak 현재 L1StandardBridge 백업
- [ ] Native Token 관련 함수 목록 작성

작업 중 확인:

- [ ] ProxyAdminOwnedBase, ReinitializableBase 상속 추가
- [ ] Constructor `_disableInitializers()` 패턴 적용
- [ ] Storage spacer 추가
- [ ] Initialize 시그니처 업데이트
- [ ] `paused()`, `superchainConfig()` 함수 추가
- [ ] Tokamak Native Token 함수들 모두 유지
- [ ] Deploy.s.sol 업데이트
- [ ] Initializable.t.sol 업데이트

작업 완료 후 확인:

- [ ] `forge build` 성공
- [ ] ChainAssertions 통과 (E2E 테스트)
- [ ] Native Token 기능 검증

---

## 🔍 트러블슈팅

### 문제 1: "storage value is not 0xff"
**원인**: Constructor에서 `initialize()` 직접 호출
**해결**: Constructor에서 `_disableInitializers()` 호출

### 문제 2: "Cannot implicitly convert"
**원인**: Initialize 시그니처 불일치
**해결**: Deploy/Test 파일에서 타입 캐스팅 수정

### 문제 3: "ProxyAdminOwnedBase_NotResolvedDelegateProxy"
**원인**: ChainAssertions가 올바르게 체크 중
**해결**: 이 에러는 예상된 동작 (Implementation contract에서 정상)

### 문제 4: Native Token 기능 손상
**원인**: 함수 삭제 또는 로직 변경
**해결**: 기존 Tokamak 버전과 diff 비교, 모든 Native Token 함수 유지

---

## 📚 참고 자료

### 완료된 컨트랙트 예시
- `src/L1/L1CrossDomainMessenger.sol` - Native Token 기능 포함
- `src/L1/L1ERC721Bridge.sol` - 기본 패턴
- `src/dispute/DisputeGameFactory.sol` - 간단한 예시

### Optimism 원본
- `/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/L1/L1StandardBridge.sol`

### 문서
- `op-challenger/docs/migration-v1.16.0/L1CROSSDOMAINMESSENGER-INTEGRATION-GUIDE.md`
- `op-challenger/docs/migration-v1.16.0/DEPLOYMENT-COMPARISON-DEPLOY-VS-DEPLOYIMPLEMENTATIONS.md`

---

## 🎯 다음 세션 시작 명령어

```bash
# 1. 프로젝트 디렉토리로 이동
cd /Users/zena/tokamak-projects/tokamak-thanos

# 2. 현재 상태 확인
cd packages/tokamak/contracts-bedrock
forge build 2>&1 | tail -20

# 3. Optimism 원본 확인
cat /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/L1/L1StandardBridge.sol | head -100

# 4. Tokamak 현재 버전 확인
cat src/L1/L1StandardBridge.sol | head -100

# 5. 이 문서 열기
cat /Users/zena/tokamak-projects/tokamak-thanos/op-challenger/docs/migration-v1.16.0/NEXT-SESSION-TODO.md
```

---

## 💡 예상 소요 시간

**L1StandardBridge 작업**: 1-2시간
- 코드 작성: 30분
- 컴파일/수정: 30분
- 테스트: 30-60분

**추가 컨트랙트** (필요시):
- 각 컨트랙트당 30분-1시간

---

## ✅ 최종 목표

모든 L1 컨트랙트를 Optimism v1.16.0 패턴으로 업그레이드하면서 **Tokamak Native Token 기능 100% 보존**

**성공 기준**:
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos
go test -v ./op-e2e/faultproofs -run TestOutputCannonGame -timeout 30m
# 출력: PASS (에러 없음)
```

---

**작성자**: AI Assistant
**다음 세션 시**: 이 문서를 먼저 읽고 시작하세요!

