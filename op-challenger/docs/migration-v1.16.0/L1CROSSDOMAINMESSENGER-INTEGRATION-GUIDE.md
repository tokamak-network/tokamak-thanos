# L1CrossDomainMessenger: Optimism v1.16.0 + Tokamak Native Token 통합 가이드

**작성일**: 2025-11-10
**버전**: v2.9.0-tokamak
**상태**: ✅ 완료 (컴파일 성공)

---

## 📋 목차

1. [작업 개요](#작업-개요)
2. [문제점 및 해결 전략](#문제점-및-해결-전략)
3. [주요 변경사항](#주요-변경사항)
4. [Tokamak Native Token 기능 보존](#tokamak-native-token-기능-보존)
5. [검증 결과](#검증-결과)
6. [마이그레이션 체크리스트](#마이그레이션-체크리스트)
7. [트러블슈팅](#트러블슈팅)

---

## 작업 개요

### 목표
Optimism v1.16.0의 새로운 아키텍처 패턴을 Tokamak-Thanos에 적용하면서, **Tokamak의 Native Token 기능을 100% 보존**합니다.

### 핵심 요구사항
1. ✅ Optimism v1.16.0의 ProxyAdminOwnedBase, ReinitializableBase 패턴 적용
2. ✅ Tokamak Native Token 전송/수신 기능 유지
3. ✅ 기존 배포와의 호환성 유지 (storage layout)
4. ✅ 업그레이드 가능한 구조 유지

### 배경
- Optimism v1.16.0은 새로운 initialization 패턴 도입
- Tokamak은 ETH 대신 ERC20 Native Token을 사용하는 독자적인 아키텍처
- 두 시스템의 장점을 모두 취하는 통합이 필요

---

## 문제점 및 해결 전략

### 문제 1: 상속 구조 불일치

**기존 (Tokamak)**:
```solidity
contract L1CrossDomainMessenger is CrossDomainMessenger, OnApprove, ISemver
```

**Optimism v1.16.0**:
```solidity
contract L1CrossDomainMessenger is
    CrossDomainMessenger,
    ProxyAdminOwnedBase,
    ReinitializableBase,
    ISemver
```

**해결책**: Tokamak OnApprove 유지하면서 Optimism base contracts 추가
```solidity
contract L1CrossDomainMessenger is
    CrossDomainMessenger,
    ProxyAdminOwnedBase,      // ✅ Optimism v1.16.0
    ReinitializableBase,      // ✅ Optimism v1.16.0
    OnApprove,                // ✅ Tokamak Native Token
    ISemver
```

---

### 문제 2: Constructor/Initialize 패턴 변경

**기존 (Tokamak)**:
```solidity
constructor() CrossDomainMessenger() {
    initialize({
        _superchainConfig: SuperchainConfig(address(0)),
        _portal: OptimismPortal(payable(address(0))),
        _systemConfig: SystemConfig(address(0))
    });
}
```

**문제점**:
- Implementation contract에서 직접 `initialize()` 호출
- Proxy pattern에 맞지 않음
- Optimism v1.16.0의 `_disableInitializers()` 패턴 미적용

**해결책**:
```solidity
constructor() ReinitializableBase(2) {
    _disableInitializers();
}

function initialize(
    ISuperchainConfig _superchainConfig,  // Tokamak 호환성 유지
    IOptimismPortal _portal,
    ISystemConfig _systemConfig
) external reinitializer(initVersion()) {
    _assertOnlyProxyAdminOrProxyAdminOwner();
    portal = _portal;
    systemConfig = _systemConfig;
    __CrossDomainMessenger_init({...});
}
```

---

### 문제 3: Storage Layout 호환성

**문제점**:
- 기존 배포된 Proxy와 storage layout이 달라지면 업그레이드 실패
- superchainConfig 위치 변경 필요

**해결책**: Spacer 패턴 사용
```solidity
/// @custom:legacy
/// @custom:spacer superchainConfig
address private spacer_251_0_20;

IOptimismPortal public portal;

/// @custom:legacy
/// @custom:spacer systemConfig (old location)
address private spacer_253_0_20;

string public constant version = "2.9.0-tokamak";

ISystemConfig public systemConfig;  // 새로운 위치
```

---

### 문제 4: Native Token 기능 손상 위험

**Tokamak의 핵심 로직**:
```solidity
// 6개 파라미터: _to, _mint, _value, _gasLimit, _isCreation, _data
portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
//                              ^^^^^^  ^^^^^^
//                              mint    value (L2 mint + L1 transfer)
```

**Optimism v1.16.0**:
```solidity
// 5개 파라미터 (mint 없음)
portal.depositTransaction(_to, _value, _gasLimit, false, _data);
```

**해결책**:
1. `IOptimismPortal2` 인터페이스에 6-parameter 버전 유지
2. Tokamak의 `_sendMessage` 로직 그대로 유지
3. ERC20 approve 로직 3개소 모두 유지

---

## 주요 변경사항

### 1. L1CrossDomainMessenger.sol

#### 1.1 Import 변경
```solidity
// 추가된 Imports
import { ProxyAdminOwnedBase } from "src/L1/ProxyAdminOwnedBase.sol";
import { ReinitializableBase } from "src/universal/ReinitializableBase.sol";

// 인터페이스로 변경
import { ISuperchainConfig } from "interfaces/L1/ISuperchainConfig.sol";
import { ISystemConfig } from "interfaces/L1/ISystemConfig.sol";
import { IOptimismPortal2 as IOptimismPortal } from "interfaces/L1/IOptimismPortal2.sol";
```

#### 1.2 상속 구조 변경
```solidity
contract L1CrossDomainMessenger is
    CrossDomainMessenger,
    ProxyAdminOwnedBase,      // ✅ NEW
    ReinitializableBase,      // ✅ NEW
    OnApprove,                // ✅ KEPT (Tokamak)
    ISemver
```

#### 1.3 Storage Layout
```solidity
// Legacy slots (spacers)
address private spacer_251_0_20;      // OLD: superchainConfig
IOptimismPortal public portal;
address private spacer_253_0_20;      // OLD: systemConfig location
string public constant version = "2.9.0-tokamak";

// New location
ISystemConfig public systemConfig;
```

#### 1.4 Constructor & Initialize
```solidity
constructor() ReinitializableBase(2) {
    _disableInitializers();
}

function initialize(
    ISuperchainConfig _superchainConfig,  // Kept for compatibility
    IOptimismPortal _portal,
    ISystemConfig _systemConfig
) external reinitializer(initVersion()) {
    _assertOnlyProxyAdminOrProxyAdminOwner();
    portal = _portal;
    systemConfig = _systemConfig;
    __CrossDomainMessenger_init({...});
}
```

#### 1.5 paused() 함수 변경
```solidity
// OLD
function paused() public view override returns (bool) {
    return superchainConfig.paused();
}

// NEW
function paused() public view override returns (bool) {
    return superchainConfig().paused();  // systemConfig를 통해 접근
}

function superchainConfig() public view returns (ISuperchainConfig) {
    return systemConfig.superchainConfig();
}
```

#### 1.6 _nativeToken() 함수
```solidity
function _nativeToken() internal view returns (address) {
    // Tokamak: SystemConfig has nativeTokenAddress() function
    return ISystemConfig(address(systemConfig)).nativeTokenAddress();
}
```

#### 1.7 _sendMessage() - **Critical for Tokamak**
```solidity
function _sendMessage(address _to, uint64 _gasLimit, uint256 _value, bytes memory _data) internal override {
    // Tokamak: Deny direct ETH deposits, only accept Native Token
    require(msg.value == 0, "Deny depositing ETH");
    // Tokamak: OptimismPortal expects 6 params (_to, _mint, _value, _gasLimit, _isCreation, _data)
    // _mint and _value are both set to _value for native token functionality
    portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
    //                              ^^^^^^  ^^^^^^
    //                              BOTH VALUES MUST BE PRESENT!
}
```

#### 1.8 _isOtherMessenger() 구현
```solidity
function _isOtherMessenger() internal view override returns (bool) {
    return msg.sender == address(portal) && portal.l2Sender() == address(otherMessenger);
}
```

---

### 2. 인터페이스 업데이트

#### 2.1 ISystemConfig.sol
```solidity
interface ISystemConfig is IProxyAdminOwnedBase {
    // ... existing functions ...

    function nativeTokenAddress() external view returns (address);  // ✅ ADDED
}
```

#### 2.2 IOptimismPortal2.sol
```solidity
interface IOptimismPortal2 is IProxyAdminOwnedBase {
    // ... existing functions ...

    function depositTransaction(
        address _to,
        uint256 _mint,      // ✅ ADDED for Tokamak
        uint256 _value,
        uint64 _gasLimit,
        bool _isCreation,
        bytes memory _data
    ) external payable;
}
```

---

### 3. 테스트 파일 업데이트

#### 3.1 Deploy.s.sol
```solidity
// Import 추가
import { ISuperchainConfig } from "interfaces/L1/ISuperchainConfig.sol";
import { ISystemConfig } from "interfaces/L1/ISystemConfig.sol";
import { IOptimismPortal2 as IOptimismPortal } from "interfaces/L1/IOptimismPortal2.sol";

// Initialize 호출 수정
_innerCallData: abi.encodeCall(
    L1CrossDomainMessenger.initialize,
    (
        ISuperchainConfig(superchainConfigProxy),
        IOptimismPortal(payable(optimismPortalProxy)),
        ISystemConfig(systemConfigProxy)
    )
)
```

#### 3.2 test/L1/L1CrossDomainMessenger.t.sol
```solidity
// Import 추가
import { ISuperchainConfig } from "interfaces/L1/ISuperchainConfig.sol";
import { ISystemConfig } from "interfaces/L1/ISystemConfig.sol";
import { IOptimismPortal2 as IOptimismPortal } from "interfaces/L1/IOptimismPortal2.sol";

// Initialize 호출 수정
l1CrossDomainMessenger.initialize(
    ISuperchainConfig(address(superchainConfig)),
    IOptimismPortal(payable(address(optimismPortal))),
    ISystemConfig(address(systemConfig))
);
```

#### 3.3 test/vendor/Initializable.t.sol
```solidity
// 동일한 import 및 initialize 호출 수정 적용
```

---

## Tokamak Native Token 기능 보존

### ✅ 보존된 핵심 함수 (7개)

#### 1. sendNativeTokenMessage()
```solidity
function sendNativeTokenMessage(
    address _target,
    uint256 _amount,
    bytes calldata _message,
    uint32 _minGasLimit
) external
```
- **역할**: 사용자가 Native Token을 L2로 전송
- **변경사항**: 없음 (100% 유지)

#### 2. _sendNativeTokenMessage()
```solidity
function _sendNativeTokenMessage(
    address _sender,
    address _target,
    uint256 _amount,
    uint32 _minGasLimit,
    bytes calldata _message
) internal
```
- **역할**: 내부 Native Token 전송 로직
- **핵심**:
  ```solidity
  IERC20(_nativeTokenAddress).safeTransferFrom(_sender, address(this), _amount);
  IERC20(_nativeTokenAddress).approve(address(portal), _amount);
  ```
- **변경사항**: 없음

#### 3. onApprove()
```solidity
function onApprove(
    address _owner,
    address,
    uint256 _amount,
    bytes calldata _data
) external override returns (bool)
```
- **역할**: ERC20 approveAndCall 콜백
- **변경사항**: 없음

#### 4. unpackOnApproveData()
```solidity
function unpackOnApproveData(bytes calldata _data)
    internal pure
    returns (address _to, uint32 _minGasLimit, bytes calldata _message)
```
- **역할**: approve data 파싱
- **변경사항**: 없음

#### 5. relayMessage()
```solidity
function relayMessage(
    uint256 _nonce,
    address _sender,
    address _target,
    uint256 _value,
    uint256 _minGasLimit,
    bytes calldata _message
) external payable override
```
- **역할**: L2→L1 메시지 릴레이 (Native Token 포함)
- **핵심 로직 유지**:
  ```solidity
  // Native Token 수신
  IERC20(_nativeTokenAddress).safeTransferFrom(address(portal), address(this), _value);

  // Target에 approve
  IERC20(_nativeTokenAddress).approve(_target, _value);

  // Call
  SafeCall.call(_target, gasleft() - RELAY_RESERVED_GAS, 0, _message);

  // Approve reset
  IERC20(_nativeTokenAddress).approve(_target, 0);
  ```
- **변경사항**: 없음

#### 6. nativeTokenAddress()
```solidity
function nativeTokenAddress() public view returns (address)
```
- **역할**: Native Token 주소 조회
- **변경사항**: 없음

#### 7. _nativeToken()
```solidity
function _nativeToken() internal view returns (address) {
    return ISystemConfig(address(systemConfig)).nativeTokenAddress();
}
```
- **역할**: 내부 Native Token 주소 조회
- **변경사항**: `systemConfig.nativeTokenAddress()` 호출 방식으로 개선

---

### 🔍 Critical 로직 검증

#### depositTransaction 6-parameter 패턴
```solidity
// ✅ CORRECT (Tokamak)
portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
//                              ^^^^^^  ^^^^^^
//                              mint    value

// 설명:
// - _mint (_value): L2에서 mint할 Native Token 양
// - _value (_value): L1에서 전송하는 Native Token 양
// - 둘 다 같은 값이어야 L1↔L2 Native Token 균형 유지
```

#### ERC20 Approve 흐름
```solidity
// 1. 전송 시 (sendNativeTokenMessage)
IERC20(token).safeTransferFrom(sender, this, amount);  // User → L1XDM
IERC20(token).approve(portal, amount);                 // L1XDM → Portal

// 2. 수신 시 (relayMessage)
IERC20(token).safeTransferFrom(portal, this, value);   // Portal → L1XDM
IERC20(token).approve(target, value);                  // L1XDM → Target (before call)
SafeCall.call(target, ...);                            // Execute
IERC20(token).approve(target, 0);                      // Reset approve (after call)
```

---

## 검증 결과

### 컴파일 성공
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
forge build
# ✅ Success (only warnings, no errors)
```

### 함수 존재 검증
```bash
=== Tokamak Native Token 핵심 함수 검증 ===
📋 sendNativeTokenMessage            ✅
📋 _sendNativeTokenMessage           ✅
📋 onApprove                         ✅
📋 unpackOnApproveData               ✅
📋 relayMessage                      ✅
📋 nativeTokenAddress                ✅
📋 _nativeToken                      ✅
```

### 로직 검증
```bash
=== 중요 로직 검증 ===
📋 depositTransaction 파라미터:
   Old: portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
   New: portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
   ✅ Native Token mint/value 패턴 유지됨

=== ERC20 approve 로직 검증 ===
   Old: 3 개소
   New: 3 개소
   ✅ ERC20 approve 로직 유지됨
```

### 변경된 파일 목록
```
✅ 수정:
  - src/L1/L1CrossDomainMessenger.sol
  - interfaces/L1/ISystemConfig.sol
  - interfaces/L1/IOptimismPortal2.sol
  - scripts/Deploy.s.sol
  - test/L1/L1CrossDomainMessenger.t.sol
  - test/vendor/Initializable.t.sol

✅ 추가 (from Optimism):
  - src/L1/ProxyAdminOwnedBase.sol
  - src/universal/ReinitializableBase.sol
```

---

## 마이그레이션 체크리스트

### Phase 1: 준비 ✅
- [x] Optimism v1.16.0 base contracts 복사
  - [x] ProxyAdminOwnedBase.sol
  - [x] ReinitializableBase.sol
- [x] 기존 Tokamak L1CrossDomainMessenger 백업
- [x] Native Token 기능 목록 작성

### Phase 2: 통합 ✅
- [x] L1CrossDomainMessenger 상속 구조 변경
- [x] Constructor/Initialize 패턴 업데이트
- [x] Storage layout spacer 추가
- [x] paused() 함수 수정
- [x] _nativeToken() 함수 업데이트
- [x] **_sendMessage() 6-parameter 유지**

### Phase 3: 인터페이스 업데이트 ✅
- [x] ISystemConfig에 nativeTokenAddress() 추가
- [x] IOptimismPortal2에 6-parameter depositTransaction 추가

### Phase 4: 테스트 파일 수정 ✅
- [x] Deploy.s.sol initialize 호출 수정
- [x] L1CrossDomainMessenger.t.sol 수정
- [x] Initializable.t.sol 수정

### Phase 5: 검증 ✅
- [x] 컴파일 성공 확인
- [x] Native Token 함수 존재 확인
- [x] depositTransaction 파라미터 검증
- [x] ERC20 approve 로직 검증

### Phase 6: E2E 테스트 (다음 단계)
- [ ] Fault proof 테스트 실행
- [ ] Native Token deposit 테스트
- [ ] Native Token withdrawal 테스트
- [ ] Cross-domain message 테스트

---

## 트러블슈팅

### 문제 1: "Cannot implicitly convert... to ISuperchainConfig"

**증상**:
```solidity
Error (5407): Cannot implicitly convert component at position 0 from
"contract SuperchainConfig" to "contract ISuperchainConfig".
```

**원인**: Initialize 호출 시 concrete type 전달

**해결책**:
```solidity
// ❌ BAD
l1CrossDomainMessenger.initialize(
    SuperchainConfig(addr),
    OptimismPortal(addr),
    SystemConfig(addr)
);

// ✅ GOOD
l1CrossDomainMessenger.initialize(
    ISuperchainConfig(address(addr)),
    IOptimismPortal(payable(address(addr))),
    ISystemConfig(address(addr))
);
```

---

### 문제 2: "storage value is not 0xff at the given slot"

**증상**:
```
DeployUtils: storage value is not 0xff at the given slot and offset
```

**원인**: Implementation contract에서 `initialize()` 직접 호출

**해결책**:
```solidity
// ❌ BAD
constructor() CrossDomainMessenger() {
    initialize(...);  // Don't call initialize in constructor!
}

// ✅ GOOD
constructor() ReinitializableBase(2) {
    _disableInitializers();  // Prevent initialization of implementation
}
```

---

### 문제 3: Native Token 전송 실패

**증상**: depositTransaction에서 Native Token이 L2로 전송되지 않음

**원인**: depositTransaction 파라미터 개수 불일치

**해결책**:
```solidity
// ❌ BAD (Optimism v1.16.0 style - 5 parameters)
portal.depositTransaction(_to, _value, _gasLimit, false, _data);

// ✅ GOOD (Tokamak style - 6 parameters)
portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
//                              ^^^^^^  ^^^^^^
//                              mint    value (BOTH REQUIRED!)
```

`IOptimismPortal2` 인터페이스도 6-parameter 버전으로 업데이트해야 함!

---

### 문제 4: paused() 함수 호출 실패

**증상**:
```solidity
Error: Member "paused" not found in ISystemConfig
```

**원인**: SystemConfig에는 paused() 메서드가 없음

**해결책**:
```solidity
// ❌ BAD
function paused() public view override returns (bool) {
    return systemConfig.paused();  // systemConfig doesn't have paused()
}

// ✅ GOOD
function paused() public view override returns (bool) {
    return superchainConfig().paused();  // Get via superchainConfig()
}

function superchainConfig() public view returns (ISuperchainConfig) {
    return systemConfig.superchainConfig();
}
```

---

## 다음 단계

### 1. E2E 테스트 실행
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos
go test -v ./op-e2e/faultproofs -run TestOutputCannonGame -timeout 30m
```

### 2. Native Token 기능 테스트
```bash
# Deposit 테스트
go test -v ./op-e2e -run TestNativeTokenDeposit

# Withdrawal 테스트
go test -v ./op-e2e -run TestNativeTokenWithdrawal

# Cross-domain message 테스트
go test -v ./op-e2e -run TestNativeTokenMessage
```

### 3. 추가 컨트랙트 검토
다른 컨트랙트들도 v1.16.0 패턴 적용 필요:
- [ ] SystemConfig
- [ ] OptimismPortal2
- [ ] DisputeGameFactory (✅ 이미 완료)
- [ ] AnchorStateRegistry

### 4. 문서화
- [ ] 배포 가이드 업데이트
- [ ] 아키텍처 다이어그램 업데이트
- [ ] API 문서 업데이트

---

## 참고 자료

### Optimism v1.16.0 관련
- [Optimism Monorepo](https://github.com/ethereum-optimism/optimism)
- [ProxyAdminOwnedBase Pattern](https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/src/L1/ProxyAdminOwnedBase.sol)
- [ReinitializableBase Pattern](https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/src/universal/ReinitializableBase.sol)

### Tokamak Native Token 관련
- Tokamak L1CrossDomainMessenger: `src/L1/L1CrossDomainMessenger.sol`
- OnApprove Interface: `src/L1/OnApprove.sol`
- Native Token SystemConfig: `src/L1/SystemConfig.sol`

### 관련 문서
- [DEPLOYMENT-COMPARISON-DEPLOY-VS-DEPLOYIMPLEMENTATIONS.md](./DEPLOYMENT-COMPARISON-DEPLOY-VS-DEPLOYIMPLEMENTATIONS.md)
- [Optimism Dispute Games 배포 파이프라인](../../op-deployer/pkg/deployer/pipeline/dispute_games.go)

---

## 작성자 노트

이 통합 작업의 핵심은 **두 시스템의 장점을 모두 취하는 것**입니다:

1. **Optimism v1.16.0의 장점**:
   - 더 안전한 initialization 패턴
   - ProxyAdmin 권한 관리 개선
   - 버전별 재초기화 지원

2. **Tokamak의 장점**:
   - Native Token (ERC20) 지원
   - 유연한 토큰 경제학
   - ETH 의존성 제거

이 두 가지를 결합함으로써, Tokamak-Thanos는 Optimism 생태계의 최신 보안 패턴을 따르면서도 독자적인 Native Token 기능을 유지할 수 있습니다.

**가장 중요한 포인트**: `depositTransaction(_to, _value, _value, ...)` 패턴에서 `_value`를 두 번 전달하는 것은 실수가 아니라 Tokamak Native Token의 핵심 설계입니다. 이는 L1에서 전송한 토큰 양과 L2에서 mint할 토큰 양이 일치해야 한다는 것을 보장합니다.

---

**문서 버전**: 1.0
**최종 업데이트**: 2025-11-10
**검증 상태**: ✅ 컴파일 성공, E2E 테스트 대기 중

