# Tokamak Native Token 구조 분석

## 개요
이 문서는 Tokamak-Thanos의 네이티브 토큰 구조를 분석하고, Optimism과의 차이점을 정리합니다.

## 1. 핵심 차이점

### 1.1 Optimism의 구조
- **Native Token**: ETH (Ethereum의 네이티브 토큰)
- **L1**: ETH를 직접 사용
- **L2**: ETH를 직접 사용
- **Bridge**: ETH를 `OptimismPortal`에 예치

### 1.2 Tokamak의 구조
- **L1 Gas Token**: ETH (Ethereum 네트워크 수수료)
- **L2 Gas Token**: TON (Tokamak L2 네트워크 수수료)
- **TON**: L1에서 ERC20, L2에서 Native 토큰
- **ETH**: L1에서 Native, L2에서 ERC20
- **TON Lock**: `OptimismPortal`에 저장
- **ETH Lock**: `L1StandardBridge`에 저장

## 2. 주요 컨트랙트 분석

### 2.1 L2NativeToken.sol
**위치**: `packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol`

**역할**: TON 토큰의 ERC20 구현

```solidity
contract L2NativeToken is Ownable, ERC20Detailed, SeigToken {
    constructor() ERC20Detailed("Tokamak Network Token", "TON", 18) { }

    function faucet(uint256 _amount) external {
        _mint(msg.sender, _amount);
    }
}
```

**특징**:
- TON은 ERC20 토큰
- `SeigToken`을 상속하여 스테이킹 관련 기능 포함
- 18 decimals 사용

### 2.2 L1StandardBridge.sol
**위치**: `packages/tokamak/contracts-bedrock/src/L1/L1StandardBridge.sol`

**Tokamak 특화 로직**:

```solidity
function finalizeBridgeNativeToken(
    address _from,
    address _to,
    uint256 _amount,
    bytes calldata _extraData
)
    public
    override
    onlyOtherBridge
{
    require(paused() == false, "L1 StandardBridge: paused");
    address _nativeTokenAddress = _nativeToken();
    require(_to != address(this), "StandardBridge: cannot send to self");
    require(_to != address(messenger), "StandardBridge: cannot send to messenger");

    // If TON is chosen as the native token, we have to perform those 2 steps below
    // TON does not allow transferFrom to the recipient directly
    IERC20(_nativeTokenAddress).safeTransferFrom(address(messenger), address(this), _amount);
    IERC20(_nativeTokenAddress).safeTransfer(_to, _amount);

    _emitNativeTokenBridgeFinalized(_from, _to, _amount, _extraData);
}
```

**핵심**:
- TON은 `transferFrom`을 직접 recipient에게 할 수 없음
- 2단계 전송 필요:
  1. `messenger` → `bridge`
  2. `bridge` → `recipient`

### 2.3 ETH 처리 방식

**Optimism**:
- L1, L2 모두 ETH는 네이티브 토큰
- L1에서 `OptimismPortal`에 직접 예치
- `ETHLockbox` 사용 (Super 업그레이드)

**Tokamak**:
- L1에서 ETH는 Native, L2에서는 ERC20
- L1에서 `L1StandardBridge`에 저장 (`deposits` mapping)
- TON이 L2의 네이티브 gas token
- `ETHLockbox` **사용하지 않음** (배포는 되지만 실제 사용 안 함)

## 3. SystemConfig의 Gas Paying Token

### 3.1 Optimism
```solidity
// ETH가 gas paying token
function gasPayingToken() external view returns (address, uint8) {
    return (address(0), 18);  // address(0) = ETH
}
```

### 3.2 Tokamak
```solidity
// TON이 gas paying token
function gasPayingToken() external view returns (address, uint8) {
    return (nativeTokenAddress, 18);  // TON의 ERC20 주소
}
```

## 4. 배포 시 고려사항

### 4.1 ETHLockbox
- **Optimism**: 사용 (Super 업그레이드에서 도입)
- **Tokamak**: **사용하지 않음**
- **이유**:
  - Tokamak의 `OPContractsManager`는 Optimism 코드를 복사했기 때문에 배포 코드는 존재
  - 하지만 `OptimismPortal`, `L1StandardBridge`에서 실제로 사용하지 않음
  - TON은 `OptimismPortal`에, ETH는 `L1StandardBridge`에 각각 lock됨

### 4.2 DelayedWETH
- **Optimism**: 사용 (Fault Proof bond 관리)
- **Tokamak**: ✅ **사용함**
  - Fault Proof 게임 참여 시 ETH를 WETH로 변환하여 bond 제공
  - DelayedWETH 컨트랙트 존재: `src/dispute/weth/DelayedWETH.sol`
  - FaultDisputeGame에서 필수로 사용
  - 생성자 수정 완료: `_disableInitializers()` 패턴 적용

### 4.3 Native Token 설정
- **Deploy Config**: `l1NativeTokenAddress` 필드로 TON 주소 지정
- **SystemConfig**: `nativeTokenAddress` 필드에 TON 주소 저장

## 5. DeployImplementations.s.sol 수정 사항

### 5.1 ETHLockbox 관련

#### 현재 상태 (수정 완료)
```solidity
function deployETHLockboxImpl(Output memory _output) private {
    // Tokamak uses ERC20 form of ETH, not ETHLockbox
    // Set to address(0) to indicate it's not used
    _output.ethLockboxImpl = IETHLockbox(address(0));

    /* Original Optimism implementation (주석 처리):
    IETHLockbox impl = IETHLockbox(
        DeployUtils.createDeterministic({
            _name: "ETHLockbox",
            _args: DeployUtils.encodeConstructor(abi.encodeCall(IETHLockbox.__constructor__, ())),
            _salt: _salt
        })
    );
    vm.label(address(impl), "ETHLockboxImpl");
    _output.ethLockboxImpl = impl;
    */
}
```

#### 검증 제거
```solidity
// ChainAssertions.checkETHLockboxImpl(_output.ethLockboxImpl, _output.optimismPortalImpl);
```

#### Output 검증에서 제외
```solidity
address[] memory addrs2 = Solarray.addresses(
    address(_output.systemConfigImpl),
    address(_output.l1CrossDomainMessengerImpl),
    address(_output.l1ERC721BridgeImpl),
    address(_output.l1StandardBridgeImpl),
    address(_output.optimismMintableERC20FactoryImpl),
    address(_output.disputeGameFactoryImpl),
    address(_output.anchorStateRegistryImpl)
    // address(_output.ethLockboxImpl)  // Tokamak doesn't use ETHLockbox
);
```

### 5.2 DelayedWETH 관련 (수정 완료)

#### 확인 결과
1. ✅ Tokamak의 Fault Proof 시스템은 **WETH를 사용**
2. ✅ DelayedWETH 컨트랙트 존재: `src/dispute/weth/DelayedWETH.sol`
3. ✅ FaultDisputeGame에서 필수로 사용

#### 수정 사항
```solidity
// 원래 Tokamak 코드 (잘못된 패턴):
constructor(uint256 _delay) {
    DELAY_SECONDS = _delay;
    initialize({ _owner: address(0), _config: SuperchainConfig(address(0)) });
}

// 수정 후 (Optimism 패턴):
constructor(uint256 _delay) {
    DELAY_SECONDS = _delay;
    _disableInitializers();  // Implementation 패턴
}
```

## 6. OPContractsManager 수정 사항

### 6.1 Implementations 구조체

```solidity
struct Implementations {
    // ...
    address ethLockboxImpl;  // Tokamak에서는 address(0)
    // ...
}
```

**Tokamak 처리**:
- `ethLockboxImpl`을 `address(0)`으로 설정
- OPCM이 이를 허용하는지 확인 필요
- 필요시 OPCM 로직 수정

### 6.2 검증 로직

OPCM의 `deployOPChain` 함수에서 `ethLockboxImpl`을 사용하는지 확인:
- 사용한다면: 조건부 로직 추가 (address(0)이면 스킵)
- 사용하지 않는다면: 그대로 진행

## 7. 다음 단계

### 7.1 즉시 확인 필요
1. ✅ ETHLockbox: 사용하지 않음 확인 완료
2. ✅ DelayedWETH: Tokamak에서 사용함, 생성자 수정 완료
3. ⬜ OPContractsManager: ethLockboxImpl = address(0) 허용 여부 확인

### 7.2 테스트 필요
1. ⬜ E2E 테스트에서 ETHLockbox 없이 정상 동작하는지 확인
2. ⬜ TON 브리지 기능 정상 동작 확인
3. ⬜ ETH (ERC20) 브리지 기능 정상 동작 확인

### 7.3 문서화 필요
1. ⬜ Tokamak의 네이티브 토큰 아키텍처 다이어그램
2. ⬜ ETH vs TON 브리지 플로우 차이점
3. ⬜ 마이그레이션 가이드 업데이트

## 8. 참고 코드 위치

### 8.1 Tokamak 특화 파일
- `packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol`
- `packages/tokamak/contracts-bedrock/src/L1/L1StandardBridge.sol`
- `packages/tokamak/contracts-bedrock/src/L1/OnApprove.sol`

### 8.2 수정된 파일
- `packages/tokamak/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`
- `packages/tokamak/contracts-bedrock/src/L1/SystemConfig.sol`
- `packages/tokamak/contracts-bedrock/src/L1/OptimismPortal2.sol`

## 9. 결론

**핵심 요약**:
1. ✅ **L1 Gas Token**: ETH (Ethereum 네트워크)
2. ✅ **L2 Gas Token**: TON (Tokamak L2 네트워크)
3. ✅ **TON**: L1에서 ERC20, L2에서 Native → `OptimismPortal`에 lock
4. ✅ **ETH**: L1에서 Native, L2에서 ERC20 → `L1StandardBridge`에 lock
5. ✅ **ETHLockbox**: 사용하지 않음 → `address(0)` 설정
6. ✅ **DelayedWETH**: 사용함 → Fault Proof bond 관리용, 생성자 수정 완료
7. ✅ **SystemConfig**, **OptimismPortal2**: 생성자 수정 완료 (`_disableInitializers()`)

**다음 작업**:
- ✅ DelayedWETH 확인 완료
- ✅ Implementation 컨트랙트 생성자 수정 완료
- 🔄 E2E 테스트 진행
- 🔄 추가 에러 발생 시 분석 및 수정

