# OptimismPortal2 마이그레이션 전략

## 📊 현재 상태 분석

### Tokamak OptimismPortal2 (Ecotone 기반)
**네이티브 토큰 핵심 기능:**
```solidity
// 1. 네이티브 토큰 주소 가져오기
function _nativeToken() internal view returns (address) {
    return systemConfig.nativeTokenAddress();
}

// 2. Deposit (6개 파라미터) - _mint와 _value의 차이
function depositTransaction(
    address _to,
    uint256 _mint,    // L1에서 lock → L2 recipient 계정으로 mint될 토큰량
    uint256 _value,   // L2 트랜잭션 실행 시 msg.value (callvalue)
    uint64 _gasLimit,
    bool _isCreation,
    bytes memory _data
)

/**
 * 🔑 핵심 차이점:
 *
 * _mint:  L1에서 Portal에 lock되는 토큰량
 *         → L2에서 _to 주소의 잔액으로 mint됨
 *         → L2 계정의 "총 받을 토큰량"
 *
 * _value: L2에서 트랜잭션 실행 시 msg.value
 *         → L2 컨트랙트 호출 시 전달되는 값
 *         → "즉시 사용할 토큰량"
 *
 * 💡 예시:
 *
 * Case 1: 단순 토큰 전송
 *   _mint = 100 TON
 *   _value = 0
 *   _to = Alice 계정
 *   _data = ""
 *
 *   → L1: Portal에 100 TON lock
 *   → L2: Alice 계정에 100 TON mint
 *   → L2: 컨트랙트 호출 없음 (value=0)
 *
 * Case 2: 토큰과 함께 컨트랙트 호출
 *   _mint = 100 TON
 *   _value = 50 TON
 *   _to = SomeContract
 *   _data = "0x1234..."
 *
 *   → L1: Portal에 100 TON lock
 *   → L2: SomeContract에 100 TON mint
 *   → L2: SomeContract.call{value: 50 TON}(data)
 *   → L2: 나머지 50 TON은 SomeContract 잔액으로 유지
 *
 * Case 3: 모든 토큰을 바로 사용
 *   _mint = 100 TON
 *   _value = 100 TON
 *   _to = SomeContract
 *   _data = "0x1234..."
 *
 *   → L1: Portal에 100 TON lock
 *   → L2: SomeContract에 100 TON mint
 *   → L2: SomeContract.call{value: 100 TON}(data)
 *   → L2: 모든 토큰이 즉시 사용됨
 */

/**
 * 📊 _mint vs _value 비교 표
 *
 * ┌──────────────┬─────────────────────────┬──────────────────────────┐
 * │    항목      │       _mint             │         _value           │
 * ├──────────────┼─────────────────────────┼──────────────────────────┤
 * │ L1 처리      │ Portal에 lock           │ (처리 없음)              │
 * │              │ IERC20.transferFrom()   │                          │
 * ├──────────────┼─────────────────────────┼──────────────────────────┤
 * │ L2 처리      │ _to 계정 잔액 증가      │ msg.value로 전달         │
 * │              │ (balance += _mint)      │ (트랜잭션 callvalue)     │
 * ├──────────────┼─────────────────────────┼──────────────────────────┤
 * │ 용도         │ 총 받을 토큰량          │ 컨트랙트 호출 시 사용량  │
 * ├──────────────┼─────────────────────────┼──────────────────────────┤
 * │ 제약조건     │ _mint >= _value         │ _value <= _mint          │
 * │              │ (당연히!)               │ (가진 것보다 많이 쓸 수 없음)│
 * ├──────────────┼─────────────────────────┼──────────────────────────┤
 * │ 일반적 사용  │ 항상 지정               │ 컨트랙트 호출 시에만     │
 * └──────────────┴─────────────────────────┴──────────────────────────┘
 *
 * 💡 왜 두 개가 필요한가?
 *
 * Ethereum에서 ETH를 보낼 때:
 *   - msg.value: 트랜잭션과 함께 전송되는 ETH
 *   - 받는 주소의 잔액: msg.value만큼 증가
 *   - 동일한 값!
 *
 * Tokamak에서 Native Token을 보낼 때:
 *   - _mint: L1에서 lock → L2 계정 잔액 증가
 *   - _value: L2 트랜잭션의 msg.value
 *   - 다를 수 있음!
 *
 * 🎯 실제 사용 예시:
 *
 * "Alice가 100 TON을 L2로 보내면서,
 *  그 중 30 TON으로 DEX 컨트랙트를 호출하고 싶다"
 *
 * → _mint = 100 TON  (Alice L2 계정에 100 TON 생성)
 * → _value = 30 TON  (DEX.swap{value: 30 TON}(...) 호출)
 * → 결과: Alice L2 계정에는 70 TON이 남음
 */

// 3. 네이티브 토큰 처리
function _depositTransaction(...) {
    address _nativeTokenAddress = _nativeToken();

    // L1: _mint만큼 Portal에 lock
    if (_mint > 0) {
        IERC20(_nativeTokenAddress).safeTransferFrom(_sender, address(this), _mint);
    }

    // L2: opaqueData에 _mint와 _value 모두 포함
    bytes memory opaqueData = abi.encodePacked(_mint, _value, _gasLimit, _isCreation, _data);
    emit TransactionDeposited(from, _to, DEPOSIT_VERSION, opaqueData);
}

// 4. Withdrawal 시 네이티브 토큰 처리
function finalizeWithdrawalTransaction(...) {
    // Native token transfer/approve logic
    IERC20(_nativeTokenAddress).approve(_tx.target, _tx.value);
    IERC20(_nativeTokenAddress).safeTransfer(_tx.target, _tx.value);
}

// 5. OnApprove 인터페이스 (네이티브 토큰 전송 시 자동 deposit)
function onApprove(address _owner, address _spender, uint256 _amount, bytes calldata _data)
```

**특징:**
- ❌ AnchorStateRegistry 없음
- ❌ ETHLockbox 없음
- ❌ ProxyAdminOwnedBase/ReinitializableBase 없음
- ✅ 네이티브 토큰 (TON 등) 완전 지원
- ✅ OnApprove 패턴

---

### Optimism v1.16.0 OptimismPortal2
**주요 변경사항:**
```solidity
// 1. 새로운 상속
contract OptimismPortal2 is
    Initializable,
    ResourceMetering,
    ReinitializableBase,     // NEW
    ProxyAdminOwnedBase,     // NEW
    ISemver

// 2. AnchorStateRegistry 통합
IAnchorStateRegistry public anchorStateRegistry;  // NEW

// 3. ETHLockbox 통합
IETHLockbox public ethLockbox;  // NEW

// 4. Super Roots 지원
bool public superRootsActive;  // NEW

// 5. 새로운 initialize
function initialize(
    IAnchorStateRegistry _anchorStateRegistry,
    ISystemConfig _systemConfig,
    ISuperchainConfig _superchainConfig
)

// 6. Optimism v1.16.0 depositTransaction (5 파라미터)
function depositTransaction(
    address _to,
    uint256 _value,      // L2에서 사용할 value
    uint64 _gasLimit,
    bool _isCreation,
    bytes memory _data
) public payable {
    // Lock the ETH in the ETHLockbox
    if (msg.value > 0) ethLockbox.lockETH{ value: msg.value }();

    // opaqueData = abi.encodePacked(msg.value, _value, ...)
    // msg.value: 실제 전송된 ETH
}
```

**특징:**
- ✅ AnchorStateRegistry 통합 (dispute game 연동)
- ✅ ETHLockbox 통합 (ETH 관리)
- ✅ ProxyAdminOwnedBase/ReinitializableBase
- ✅ Custom Gas Token 지원 (별도 경로)
- ❌ OnApprove 패턴 없음

---

## ⚠️ **핵심 충돌 발견!**

### **depositTransaction 함수 시그니처 충돌**

#### Optimism v1.16.0 (5 파라미터):
```solidity
function depositTransaction(
    address _to,
    uint256 _value,      // L2에서 사용할 value
    uint64 _gasLimit,
    bool _isCreation,
    bytes memory _data
) public payable {
    // msg.value (ETH)를 ETHLockbox에 lock
    if (msg.value > 0) ethLockbox.lockETH{ value: msg.value }();
}
```

#### Tokamak (6 파라미터):
```solidity
function depositTransaction(
    address _to,
    uint256 _mint,       // 네이티브 토큰 lock 량
    uint256 _value,      // L2로 전송할 량
    uint64 _gasLimit,
    bool _isCreation,
    bytes memory _data
) external payable {
    // 네이티브 토큰을 Portal에 lock
    if (_mint > 0) {
        IERC20(_nativeToken()).safeTransferFrom(msg.sender, address(this), _mint);
    }
}
```

### **차이점:**
| 항목 | Optimism v1.16.0 | Tokamak |
|------|-----------------|---------|
| 파라미터 수 | 5개 | 6개 |
| ETH 처리 | `msg.value` → **ETHLockbox** | ❌ 사용 안 함 |
| 네이티브 토큰 | ❌ 없음 | `_mint` 파라미터 |
| **Lock 위치** | **ETHLockbox (별도)** | **OptimismPortal2 (내부)** |
| Token 보관 | ETH in ETHLockbox | Native Token in Portal |
| opaqueData | `msg.value, _value, ...` | `_mint, _value, ...` |

### **🚨 결론: 근본적인 아키텍처 차이!**

#### **토큰 보관 방식:**
```
┌─────────────────────────────────────────────┐
│ Optimism v1.16.0                            │
├─────────────────────────────────────────────┤
│                                             │
│  User (ETH) ──→ OptimismPortal2             │
│                      │                      │
│                      ↓                      │
│                 ETHLockbox.lockETH()        │
│                      │                      │
│                 [ETH 보관]                  │
│                                             │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│ Tokamak                                     │
├─────────────────────────────────────────────┤
│                                             │
│  User (TON) ──→ OptimismPortal2             │
│                      │                      │
│              IERC20.safeTransferFrom()      │
│                      │                      │
│              [TON을 Portal에 직접 보관]     │
│                                             │
│  ❌ ETHLockbox 사용 안 함!                  │
│                                             │
└─────────────────────────────────────────────┘
```

#### **문제점:**
1. **함수 시그니처 충돌**: 5개 vs 6개 파라미터
2. **토큰 보관 위치 충돌**: ETHLockbox vs Portal 내부
3. **토큰 타입 충돌**: ETH vs Native Token (TON 등)
4. **Withdrawal 로직 충돌**: ETHLockbox에서 unlock vs Portal에서 직접 전송

---

## 🔄 마이그레이션 전략 (최종 수정)

### ✅ 전략: 선택적 통합 (Optimism v1.16.0 Dispute Game + Tokamak Native Token)

**핵심 원칙: 토큰 보관은 토카막 방식 유지, Dispute Game 시스템만 통합**

```solidity
contract OptimismPortal2 is
    Initializable,
    ResourceMetering,
    ReinitializableBase,     // ✅ Optimism v1.16.0
    ProxyAdminOwnedBase,     // ✅ Optimism v1.16.0
    OnApprove,               // ✅ Tokamak 유지
    ISemver
{
    // ✅ Optimism v1.16.0: AnchorStateRegistry (Dispute Game 연동)
    IAnchorStateRegistry public anchorStateRegistry;

    // ⚠️ ETHLockbox: 구조상 필요하지만 사용 안 함
    IETHLockbox public ethLockbox;  // 배포만 하고 실제 사용 X

    // ✅ Optimism v1.16.0: Super Roots 지원
    bool public superRootsActive;

    // ✅ Tokamak: 네이티브 토큰 함수 (그대로 유지)
    function _nativeToken() internal view returns (address) {
        return systemConfig.nativeTokenAddress();
    }

    // ✅ Tokamak: depositTransaction (6 파라미터 유지)
    function depositTransaction(
        address _to,
        uint256 _mint,    // 네이티브 토큰 lock
        uint256 _value,   // L2 전송량
        uint64 _gasLimit,
        bool _isCreation,
        bytes memory _data
    ) external payable {
        // ✅ 토카막 방식: Portal에 직접 보관
        address nativeTokenAddr = _nativeToken();
        if (_mint > 0) {
            IERC20(nativeTokenAddr).safeTransferFrom(
                msg.sender,
                address(this),  // Portal에 직접 lock
                _mint
            );
        }
        // ❌ ETHLockbox 사용 안 함!
    }

    // ✅ Optimism v1.16.0: Withdrawal은 AnchorStateRegistry 사용
    function proveWithdrawalTransaction(...) {
        // anchorStateRegistry를 통한 검증
        if (!anchorStateRegistry.isGameProper(_disputeGameProxy)) {
            revert OptimismPortal_ImproperDisputeGame();
        }
        // ...
    }

    function finalizeWithdrawalTransaction(...) {
        // ✅ 토카막 방식: Portal에서 직접 네이티브 토큰 전송
        if (_tx.value > 0) {
            IERC20(_nativeToken()).approve(_tx.target, _tx.value);
            IERC20(_nativeToken()).safeTransfer(_tx.target, _tx.value);
        }
        // ❌ ETHLockbox 사용 안 함!
    }

    // ✅ Optimism v1.16.0 initialize
    function initialize(
        IAnchorStateRegistry _anchorStateRegistry,
        ISystemConfig _systemConfig,
        ISuperchainConfig _superchainConfig
    ) external reinitializer(initVersion()) {
        _assertOnlyProxyAdminOrProxyAdminOwner();
        anchorStateRegistry = _anchorStateRegistry;
        systemConfig = _systemConfig;
        // ethLockbox는 나중에 upgrade 함수에서 설정 (사용은 안 함)
    }
}
```

### ✅ 장점:
1. **Dispute Game 시스템 완전 통합**: AnchorStateRegistry 기반
2. **토카막 네이티브 토큰 100% 보존**: 변경 없음
3. **명확한 분리**: Dispute Game ≠ Token Management
4. **하위 호환성**: 기존 토카막 deposit/withdrawal 로직 유지

### ⚠️ 주의사항:
1. **ETHLockbox**: 배포는 하되 **실제 사용 안 함**
   - OPContractsManager 호환성을 위해 주소만 유지
   - Portal은 계속 네이티브 토큰 직접 보관

2. **depositTransaction**: 6-parameter **유지**
   - Optimism v1.16.0과 시그니처 다름
   - 기존 토카막 dApp과 호환성 유지

3. **opaqueData 포맷**: 토카막 방식 유지
   - `abi.encodePacked(_mint, _value, ...)`
   - Optimism의 `msg.value, _value, ...`와 다름

---

## 📝 구현 체크리스트

### Phase 1: 기본 구조 업데이트
- [ ] ProxyAdminOwnedBase, ReinitializableBase 상속 추가
- [ ] anchorStateRegistry 필드 추가
- [ ] ethLockbox 필드 추가
- [ ] superRootsActive 필드 추가
- [ ] Storage spacers 추가

### Phase 2: 초기화 패턴 업데이트
- [ ] constructor에 _disableInitializers() 추가
- [ ] initialize() 함수 시그니처 업데이트
- [ ] reinitializer 패턴 적용

### Phase 3: Tokamak 네이티브 토큰 보존
- [ ] _nativeToken() 함수 유지
- [ ] depositTransaction (6 파라미터) 유지
- [ ] _depositTransaction 네이티브 토큰 로직 유지
- [ ] onApprove 인터페이스 유지
- [ ] withdrawal 시 네이티브 토큰 처리 유지

### Phase 4: AnchorStateRegistry 통합
- [ ] disputeGameFactory() → anchorStateRegistry.disputeGameFactory()
- [ ] respectedGameType() → anchorStateRegistry.respectedGameType()
- [ ] isGameBlacklisted() → anchorStateRegistry.disputeGameBlacklist()
- [ ] proveWithdrawalTransaction 로직 업데이트

### Phase 5: 테스트 및 검증
- [ ] 단위 테스트 (native token deposit/withdrawal)
- [ ] E2E 테스트
- [ ] ChainAssertions 업데이트
- [ ] 가스 최적화 검증

---

## ⚠️ 주의사항

### Critical: 절대 손상되면 안 되는 기능
1. **네이티브 토큰 Deposit**
   - `depositTransaction(address, uint256 _mint, uint256 _value, ...)`
   - `_mint` 파라미터를 통한 네이티브 토큰 lock

2. **네이티브 토큰 Withdrawal**
   - `finalizeWithdrawalTransaction`에서 네이티브 토큰 transfer
   - IERC20.approve() 및 safeTransfer() 로직

3. **OnApprove 패턴**
   - 네이티브 토큰 approve 시 자동 deposit 트리거

4. **_nativeToken() 함수**
   - systemConfig.nativeTokenAddress() 호출

### Storage Layout 호환성
- Optimism v1.16.0의 spacer와 Tokamak의 기존 storage가 충돌하지 않도록 주의
- 기존 proxy 컨트랙트와 호환되어야 함

---

## 🎯 다음 단계

1. **Optimism v1.16.0 OptimismPortal2.sol 전체 복사**
2. **Tokamak 네이티브 토큰 함수 통합**
3. **OnApprove 인터페이스 추가**
4. **ChainAssertions 업데이트 (anchorStateRegistry 체크 추가)**
5. **컴파일 및 테스트**

---

## 📌 참고
- Optimism v1.16.0: AnchorStateRegistry 기반 dispute resolution
- Tokamak: Native Token 기반 bridge 시스템
- 두 시스템은 **독립적**이므로 동시 지원 가능

