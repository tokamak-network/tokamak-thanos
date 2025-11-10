# OptimismPortal2 마이그레이션 체크리스트

## ✅ 완료된 항목

### 1. 상속 및 기본 구조
- [x] `ProxyAdminOwnedBase` 상속 추가
- [x] `ReinitializableBase` 상속 추가
- [x] `OnApprove` 인터페이스 유지
- [x] `SafeERC20` 사용 추가
- [x] Constructor: `ReinitializableBase(2)` 호출
- [x] Constructor: `_disableInitializers()` 호출
- [x] Constructor: DISPUTE_GAME_FINALITY_DELAY_SECONDS 제거 (AnchorStateRegistry로 이동)

### 2. 스토리지 변수
- [x] `IAnchorStateRegistry public anchorStateRegistry` 추가
- [x] `IETHLockbox public ethLockbox` 추가 (배포만, 사용 안 함)
- [x] `bool public superRootsActive` 추가
- [x] `ISystemConfig public systemConfig` (인터페이스 타입으로)
- [x] Spacers 올바르게 설정:
  - `spacer_53_1_20` (superchainConfig - 삭제되고 getter로 대체)
  - `spacer_56_0_20` (disputeGameFactory - AnchorStateRegistry로 이동)
  - `spacer_58_0_32` (disputeGameBlacklist - AnchorStateRegistry로 이동)
  - `spacer_59_0_4` (respectedGameType - AnchorStateRegistry로 이동)
  - `spacer_59_4_8` (respectedGameTypeUpdatedAt - AnchorStateRegistry로 이동)
  - `spacer_61_0_32` (_balance)

### 3. Initialize 함수
- [x] 시그니처: `initialize(ISystemConfig, IAnchorStateRegistry, IETHLockbox)`
- [x] `reinitializer(initVersion())` 사용
- [x] `_assertOnlyProxyAdminOrProxyAdminOwner()` 체크
- [x] `systemConfig`, `anchorStateRegistry`, `ethLockbox` 설정
- [x] `l2Sender` 초기화 로직 유지

### 4. Upgrade 함수
- [x] `upgrade(IAnchorStateRegistry, IETHLockbox)` 추가
- [x] `reinitializer(initVersion())` 사용
- [x] `_assertOnlyProxyAdminOrProxyAdminOwner()` 체크

### 5. Getter 함수들
- [x] `version()` 함수 (4.6.0)
- [x] `paused()` - systemConfig.superchainConfig().paused()
- [x] `superchainConfig()` - systemConfig.superchainConfig()
- [x] `guardian()` - systemConfig.guardian()
- [x] `disputeGameFactory()` - anchorStateRegistry.disputeGameFactory()
- [x] `disputeGameFinalityDelaySeconds()` - anchorStateRegistry를 통해
- [x] `respectedGameType()` - anchorStateRegistry를 통해
- [x] `respectedGameTypeUpdatedAt()` - anchorStateRegistry를 통해
- [x] `disputeGameBlacklist()` - anchorStateRegistry를 통해
- [x] `nativeTokenAddress()` - Tokamak 함수 유지
- [x] `_nativeToken()` - Tokamak 내부 함수 유지

### 6. Tokamak 네이티브 토큰 함수들
- [x] `depositTransaction(address, uint256 _mint, uint256 _value, uint64, bool, bytes)` - 6 파라미터!
  - [x] `external` (not payable - ERC20 사용)
  - [x] `_mint` 파라미터로 L1에서 lock할 토큰량 지정
  - [x] `_value` 파라미터로 L2 msg.value 지정
  - [x] Portal에 직접 토큰 lock: `IERC20(_nativeToken()).safeTransferFrom(msg.sender, address(this), _mint)`
  - [x] ETHLockbox 사용 안 함!
  - [x] `opaqueData = abi.encodePacked(_mint, _value, ...)`

- [x] `finalizeWithdrawalTransaction` 네이티브 토큰 처리
  - [x] Portal에서 직접 토큰 전송
  - [x] `if (_tx.data.length != 0)`: approve 방식
  - [x] `else`: safeTransfer 방식
  - [x] `SafeCall.callWithMinGas(..., 0, ...)` - value=0 (ERC20 사용)
  - [x] ETHLockbox 사용 안 함!
  - [x] 승인 cleanup 로직 추가

- [x] `onApprove(address _owner, address, uint256 _amount, bytes calldata _data)` 콜백
  - [x] `override` 키워드
  - [x] `msg.sender == _nativeToken()` 체크
  - [x] `unpackOnApproveData` 사용
  - [x] Portal에 토큰 전송
  - [x] TransactionDeposited 이벤트 발생

- [x] `unpackOnApproveData(bytes calldata _data)` 헬퍼
  - [x] assembly로 데이터 언팩
  - [x] uint64 gasLimit으로 반환 (백업은 uint32였음)

### 7. 기타 함수들
- [x] `donateETH()` - 유지
- [x] `migrateLiquidity()` - Optimism v1.16.0 추가
- [x] `migrateToSuperRoots()` - Optimism v1.16.0 추가
- [x] `proveWithdrawalTransaction` 오버로드들
- [x] `_proveWithdrawalTransaction` 내부 함수
- [x] `checkWithdrawal` - AnchorStateRegistry 사용
- [x] `numProofSubmitters` getter
- [x] `_assertNotPaused()` 내부 함수 (whenNotPaused modifier 대체)
- [x] `_isUnsafeTarget()` 내부 함수
- [x] `_resourceConfig()` 오버라이드

### 8. 이벤트
- [x] `TransactionDeposited` - 유지
- [x] `WithdrawalProven` - 유지
- [x] `WithdrawalProvenExtension1` - 유지
- [x] `WithdrawalFinalized` - 유지
- [x] `ETHMigrated` - Optimism v1.16.0 추가
- [x] `PortalMigrated` - Optimism v1.16.0 추가
- ~~`DisputeGameBlacklisted`~~ - AnchorStateRegistry로 이동
- ~~`RespectedGameTypeSet`~~ - AnchorStateRegistry로 이동

### 9. 에러
- [x] 모든 Optimism v1.16.0 에러들 포함
- [x] `OptimismPortal_*` 패턴 사용

### 10. 인터페이스 업데이트
- [x] `IOptimismPortal2.depositTransaction` - 6 파라미터로 업데이트됨 확인

## ❌ 의도적으로 제거된 항목 (Optimism v1.16.0에서 AnchorStateRegistry로 이동)

- ~~`blacklistDisputeGame(IDisputeGame)`~~ - AnchorStateRegistry에서 관리
- ~~`setRespectedGameType(GameType)`~~ - AnchorStateRegistry에서 관리
- ~~`disputeGameBlacklist` mapping~~ - AnchorStateRegistry에서 관리
- ~~`respectedGameType` 변수~~ - AnchorStateRegistry에서 관리
- ~~`respectedGameTypeUpdatedAt` 변수~~ - AnchorStateRegistry에서 관리
- ~~`superchainConfig` 변수~~ - systemConfig를 통해 접근
- ~~`disputeGameFactory` 변수~~ - anchorStateRegistry를 통해 접근
- ~~`whenNotPaused` modifier~~ - `_assertNotPaused()` 함수로 대체
- ~~`DISPUTE_GAME_FINALITY_DELAY_SECONDS` immutable~~ - AnchorStateRegistry로 이동

## 🎯 핵심 차이점 요약

### Optimism v1.16.0 (ETH)
```solidity
function depositTransaction(
    address _to,
    uint256 _value,     // ETH amount
    uint64 _gasLimit,
    bool _isCreation,
    bytes memory _data
) public payable {
    if (msg.value > 0) ethLockbox.lockETH{ value: msg.value }();
    // ...
    opaqueData = abi.encodePacked(msg.value, _value, ...);
}

function finalizeWithdrawalTransactionExternalProof(...) {
    if (_tx.value > 0) ethLockbox.unlockETH(_tx.value);
    SafeCall.callWithMinGas(_tx.target, _tx.gasLimit, _tx.value, _tx.data);
    // ...
}
```

### Tokamak (Native Token ERC20)
```solidity
function depositTransaction(
    address _to,
    uint256 _mint,      // L1에서 lock → L2 계정에 mint
    uint256 _value,     // L2 트랜잭션 msg.value
    uint64 _gasLimit,
    bool _isCreation,
    bytes memory _data
) external {           // NOT payable
    if (_mint > 0) {
        IERC20(_nativeToken()).safeTransferFrom(msg.sender, address(this), _mint);
    }
    // ...
    opaqueData = abi.encodePacked(_mint, _value, ...);
}

function finalizeWithdrawalTransactionExternalProof(...) {
    address _nativeTokenAddress = _nativeToken();
    if (_tx.value > 0) {
        if (_tx.data.length != 0) {
            IERC20(_nativeTokenAddress).approve(_tx.target, _tx.value);
        } else {
            IERC20(_nativeTokenAddress).safeTransfer(_tx.target, _tx.value);
        }
    }
    SafeCall.callWithMinGas(_tx.target, _tx.gasLimit, 0, _tx.data); // value=0!
    // Cleanup approval
    if (_tx.data.length != 0 && _tx.value != 0) {
        IERC20(_nativeTokenAddress).approve(_tx.target, 0);
    }
}
```

## 📋 다음 단계

1. ✅ OptimismPortal2.sol 마이그레이션 완료
2. ⏭️ 컴파일 확인
3. ⏭️ ChainAssertions.sol 업데이트
4. ⏭️ Deploy.s.sol 업데이트
5. ⏭️ 테스트 파일들 업데이트
6. ⏭️ E2E 테스트 실행

---

**작성일**: 2025-11-10
**상태**: ✅ OptimismPortal2 마이그레이션 완료

