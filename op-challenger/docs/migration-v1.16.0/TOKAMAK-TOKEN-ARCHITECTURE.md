# Tokamak Token Architecture Analysis

## 개요
이 문서는 Tokamak의 토큰 구조를 분석하고, Optimism과의 차이점을 명확히 합니다.

## 1. Tokamak vs Optimism 토큰 구조

### 1.1 Optimism의 토큰 구조
```
L1: Native ETH (gas token)
    ↓ (bridge via OptimismPortal)
L2: Native ETH (gas token)
```

**특징:**
- L1과 L2 모두 **Native ETH**를 gas token으로 사용
- ETH는 `OptimismPortal`에 lock됨
- `ETHLockbox` 사용 (Super 업그레이드에서 도입)

### 1.2 Tokamak의 토큰 구조
```
L1: TON (ERC20)
    ↓ (L1StandardBridge → OptimismPortal에 lock)
L2: TON (Native gas token)

L1: ETH (Native)
    ↓ (L1StandardBridge에 lock)
L2: ETH (ERC20)
```

**특징:**
- **L1 Gas Token**: Native ETH (Ethereum 트랜잭션 수수료)
- **L2 Gas Token**: TON (Native)
- **TON**: L1에서 ERC20, L2에서 Native
- **ETH**: L1에서 Native, L2에서 ERC20
- **TON Lock**: `OptimismPortal`에 lock됨 (Portal의 `_depositTransaction`에서 처리)
- **ETH Lock**: `L1StandardBridge`에 lock됨 (`deposits` mapping에 저장)
- `ETHLockbox`: `OPContractsManager`에서 배포는 하지만 **실제로 사용하지 않음**
  - `OptimismPortal`, `L1StandardBridge`에서 ETHLockbox 참조 없음

## 2. 핵심 컨트랙트 분석

### 2.1 L2NativeToken.sol (TON)

```solidity
// packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol
contract L2NativeToken is Ownable, ERC20Detailed, SeigToken {
    constructor() ERC20Detailed("Tokamak Network Token", "TON", 18) { }

    function setSeigManager(SeigManagerI) external pure override {
        revert("TON: TON doesn't allow setSeigManager");
    }

    function faucet(uint256 _amount) external {
        _mint(msg.sender, _amount);
    }
}
```

**특징:**
- TON은 ERC20 토큰
- `SeigToken`을 상속 (Tokamak의 staking 메커니즘)
- L1에서 ERC20, L2에서 Native Token으로 사용

### 2.2 L1StandardBridge.sol

**ETH 브릿지 (L1 → L2):**
```solidity
// packages/tokamak/contracts-bedrock/src/L1/L1StandardBridge.sol (Line 238-258)
function _initiateBridgeETH(
    address _from,
    address _to,
    uint256 _amount,
    uint32 _minGasLimit,
    bytes memory _extraData
)
    internal
    override
{
    require(msg.value != 0, "StandardBridge: msg.value is zero amount");
    require(msg.value == _amount, "StandardBridge: bridging ETH must include sufficient ETH value");

    // ETH를 L1StandardBridge에 lock
    deposits[address(0)][Predeploys.ETH] = deposits[address(0)][Predeploys.ETH] + _amount;

    _emitETHBridgeInitiated(_from, _to, _amount, _extraData);

    messenger.sendMessage(
        address(otherBridge),
        abi.encodeWithSelector(this.finalizeBridgeETH.selector, _from, _to, _amount, _extraData),
        _minGasLimit
    );
}
```

**TON 브릿지 (L2 → L1 출금 완료):**
```solidity
// packages/tokamak/contracts-bedrock/src/L1/L1StandardBridge.sol (Line 449-463)
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

**핵심 차이:**
- **ETH**: `msg.value`로 받아서 `L1StandardBridge`의 `deposits` mapping에 저장
- **TON**: ERC20이므로 `safeTransferFrom` 사용, 2단계 전송 (messenger → bridge → recipient)

### 2.3 OptimismPortal의 TON Lock 메커니즘

Tokamak의 `OptimismPortal`은 **TON만** lock합니다:

```solidity
// packages/tokamak/contracts-bedrock/src/L1/OptimismPortal.sol (Line 520-538)
function _depositTransaction(
    address _sender,
    address _to,
    uint256 _mint,
    uint256 _value,
    uint64 _gasLimit,
    bool _isCreation,
    bytes calldata _data
)
    internal
    metered(_gasLimit)
{
    address _nativeTokenAddress = _nativeToken();

    // Lock TON in this contract
    if (_mint > 0) {
        IERC20(_nativeTokenAddress).safeTransferFrom(_sender, address(this), _mint);
    }
    // ... rest of function
}
```

**핵심:**
- `_nativeToken()`: TON 주소 반환
- TON은 ERC20이므로 `safeTransferFrom`으로 Portal에 lock
- **ETH는 Portal에 lock되지 않음** (L1StandardBridge에 lock됨)
- **ETHLockbox 없이 각 컨트랙트가 직접 관리**

### 2.4 SystemConfig.sol

Tokamak의 `SystemConfig`는 Native Token 주소를 관리:

```solidity
struct Addresses {
    address l1CrossDomainMessenger;
    address l1ERC721Bridge;
    address l1StandardBridge;
    address disputeGameFactory;
    address optimismPortal;
    address optimismMintableERC20Factory;
    address gasPayingToken;       // TON 주소
    address nativeTokenAddress;   // TON 주소
}
```

## 3. Optimism vs Tokamak 비교표

| 항목 | Optimism | Tokamak |
|------|----------|---------|
| L1 Gas Token | Native ETH | Native ETH |
| L2 Gas Token | Native ETH | Native TON |
| ETHLockbox | ✅ 사용 | ❌ 사용 안 함 |
| DelayedWETH | ✅ 사용 | ✅ 사용 |
| Bridge 방식 | OptimismPortal (ETH) | L1StandardBridge (TON + ETH) |
| Native Token 컨트랙트 | 없음 (Native) | L2NativeToken.sol (TON) |

## 4. DeployImplementations.s.sol 수정 필요 사항

### 4.1 ETHLockbox 관련

**Optimism:**
```solidity
function deployETHLockboxImpl(Output memory _output) private {
    IETHLockbox impl = IETHLockbox(
        DeployUtils.createDeterministic({
            _name: "ETHLockbox",
            _args: DeployUtils.encodeConstructor(abi.encodeCall(IETHLockbox.__constructor__, ())),
            _salt: _salt
        })
    );
    vm.label(address(impl), "ETHLockboxImpl");
    _output.ethLockboxImpl = impl;
}
```

**Tokamak (수정 필요):**
```solidity
function deployETHLockboxImpl(Output memory _output) private {
    // Tokamak uses ERC20 form of ETH, not ETHLockbox
    // Set to address(0) to indicate it's not used
    _output.ethLockboxImpl = IETHLockbox(address(0));
}
```

**이유:**
- Tokamak은 L1에서 Native ETH를 사용하지만, L2에서는 ERC20로 변환
- ETHLockbox는 Optimism의 Super 업그레이드에서 도입된 컨트랙트
- Tokamak의 브릿지 메커니즘은 `L1StandardBridge`에서 직접 처리하므로 불필요

### 4.2 DelayedWETH 관련

**확인 결과: ✅ Tokamak은 DelayedWETH를 사용합니다**

**DelayedWETH의 역할:**
```solidity
// packages/tokamak/contracts-bedrock/src/dispute/FaultDisputeGame.sol
contract FaultDisputeGame is IFaultDisputeGame, Clone, ISemver {
    /// @notice WETH contract for holding ETH.
    IDelayedWETH internal immutable WETH;

    // ...
}
```

**사용 목적:**
1. **Fault Proof 시스템의 Bond 관리**
   - Dispute Game에서 참가자들이 bond를 예치
   - 분쟁 해결 후 지연된 출금 (delayed withdrawal)
   - 악의적인 즉시 출금 방지

2. **OPContractsManager에서의 배포:**
```solidity
// packages/tokamak/contracts-bedrock/src/L1/OPContractsManager.sol
output.delayedWETHPermissionedGameProxy = IDelayedWETH(
    payable(
        deployProxy(_input.l2ChainId, output.opChainProxyAdmin, _input.saltMixer, "DelayedWETHPermissionedGame")
    )
);

// FaultDisputeGame 생성 시 DelayedWETH 전달
IFaultDisputeGame.GameConstructorParams({
    gameType: _input.disputeGameType,
    absolutePrestate: _input.disputeAbsolutePrestate,
    maxGameDepth: _input.disputeMaxGameDepth,
    splitDepth: _input.disputeSplitDepth,
    clockExtension: _input.disputeClockExtension,
    maxClockDuration: _input.disputeMaxClockDuration,
    vm: IBigStepper(implementation.mipsImpl),
    weth: IDelayedWETH(payable(address(output.delayedWETHPermissionedGameProxy))),  // ← 여기
    anchorStateRegistry: IAnchorStateRegistry(address(output.anchorStateRegistryProxy)),
    l2ChainId: _input.l2ChainId
})
```

**중요 사항:**
- DelayedWETH는 **Fault Proof 시스템에서 필수**
- 게임 참여 시 Native ETH를 WETH로 변환하여 bond로 사용
- Tokamak에서도 **반드시 배포해야 함**
- `DeployImplementations.s.sol`에서 `deployDelayedWETHImpl` 호출 필요

**Bond 제공 방식:**
```solidity
// FaultDisputeGame.sol - initialize 함수
WETH.deposit{ value: msg.value }();  // ETH → WETH 변환
```

**Optimism과 Tokamak 동일:**
- 둘 다 Native ETH → WETH.deposit() → DelayedWETH
- Bond는 WETH 형태로 저장됨

**파일 위치:**
```
packages/tokamak/contracts-bedrock/src/dispute/weth/DelayedWETH.sol
packages/tokamak/contracts-bedrock/src/dispute/weth/WETH98.sol
```

## 5. Fault Proof Bond Token 분석

### 5.1 Bond Token 확인 결과

**Optimism:**
- Bond Token: Native ETH
- DelayedWETH를 통해 관리

**Tokamak:**
- Bond Token: **Native ETH → WETH로 변환하여 bond 제공**
- DelayedWETH를 통해 관리 (Optimism과 동일)
- 게임 참여 시 `msg.value`로 ETH 전송 → `WETH.deposit()` 호출하여 WETH로 변환

**핵심 차이:**
```
Optimism: L1 Native ETH → WETH.deposit() → DelayedWETH (Fault Proof bond)
Tokamak:  L1 Native ETH → WETH.deposit() → DelayedWETH (Fault Proof bond)
          (동일한 메커니즘)
```

### 5.2 OptimismPortal2의 토큰 처리

**확인 결과:**
- `OptimismPortal2`는 **ETH를 직접 처리하지 않음**
- ETH 처리는 `L1StandardBridge`가 담당
- Portal은 주로 메시지 전달과 출금 증명 검증을 담당

**토큰 흐름:**
```
L1 → L2 입금:
  TON: L1StandardBridge.bridgeNativeToken() → OptimismPortal (메시지 전달)
  ETH: L1StandardBridge.bridgeERC20() → OptimismPortal (메시지 전달)

L2 → L1 출금:
  OptimismPortal.proveWithdrawalTransaction() → 증명 검증
  OptimismPortal.finalizeWithdrawalTransaction() → L1StandardBridge로 전달
```

### 5.3 Tokamak의 3가지 토큰 역할

| 토큰 | L1 형태 | L2 형태 | 용도 | 관련 컨트랙트 |
|------|---------|---------|------|--------------|
| TON | ERC20 | Native | Gas Token, Bridge | L1StandardBridge, L2NativeToken |
| ETH | Native | ERC20 | Bridge, General Use | L1StandardBridge |
| WETH | ERC20 (wrapped) | - | Fault Proof Bond (L1 only) | DelayedWETH, FaultDisputeGame |

## 6. DeployImplementations.s.sol 수정 사항 정리

### 6.1 ETHLockbox 관련 수정

**현재 상태:**
```solidity
// DeployImplementations.s.sol (Line 348)
function deployETHLockboxImpl(Output memory _output) private {
    // Tokamak uses ERC20 form of ETH, not ETHLockbox
    // Set to address(0) to indicate it's not used
    _output.ethLockboxImpl = IETHLockbox(address(0));

    /* Original Optimism implementation:
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

**수정 이유:**
- ✅ Tokamak의 `OPContractsManager`는 Optimism 코드를 복사했기 때문에 ETHLockbox 배포 코드가 있음
- ✅ 하지만 실제 컨트랙트(`OptimismPortal`, `L1StandardBridge`)에서는 **ETHLockbox를 사용하지 않음**
- ✅ Tokamak은 각 컨트랙트에서 직접 lock 처리:
  - **TON**: `OptimismPortal`에 lock
  - **ETH**: `L1StandardBridge`에 lock
- ✅ `DeployImplementations.s.sol`에서는 `address(0)` 설정으로 배포 스킵
- ⚠️ 주의: `OPContractsManager.deployOPChain()`을 사용하면 ETHLockbox가 배포되지만 사용되지 않음

### 6.2 DelayedWETH 관련 확인

**현재 상태:**
```solidity
// DeployImplementations.s.sol (Line 352)
console.log("=== Deploy deployDelayedWETHImpl ===");
deployDelayedWETHImpl(_input, output_);
```

**확인 결과:**
- ✅ DelayedWETH는 **반드시 배포 필요**
- ✅ Fault Proof 시스템에서 bond 관리에 사용
- ✅ 현재 코드 정상 (수정 불필요)

**DelayedWETH 배포 함수:**
```solidity
function deployDelayedWETHImpl(Input memory _input, Output memory _output) private {
    IDelayedWETH impl = IDelayedWETH(
        DeployUtils.createDeterministic({
            _name: "DelayedWETH",
            _args: DeployUtils.encodeConstructor(
                abi.encodeCall(IDelayedWETH.__constructor__, (_input.withdrawalDelaySeconds))
            ),
            _salt: _salt
        })
    );
    vm.label(address(impl), "DelayedWETHImpl");
    _output.delayedWETHImpl = impl;
}
```

### 6.3 검증 로직 수정

**assertValidOutput 수정:**
```solidity
function assertValidOutput(Output memory _output) internal view {
    // ... other validations ...

    // address(_output.ethLockboxImpl)  // Tokamak doesn't use ETHLockbox
    // address(_output.ratImpl)  // RAT not used in Tokamak
}
```

**ChainAssertions 수정:**
```solidity
// ChainAssertions.checkETHLockboxImpl(_output.ethLockboxImpl, _output.optimismPortalImpl);
// Tokamak doesn't use ETHLockbox
```

## 7. 결론 및 다음 단계

### 7.1 핵심 차이점 정리

| 항목 | Optimism | Tokamak | 설명 |
|------|----------|---------|------|
| **L1에서 트랜잭션 수수료** | ETH | ETH | 둘 다 이더리움 네트워크 사용 |
| **L2에서 트랜잭션 수수료** | ETH | TON | Tokamak은 TON으로 수수료 지불 |
| **TON 토큰** | 없음 | L1: ERC20<br>L2: Native | Tokamak만 사용 |
| **ETH 토큰** | L1: Native<br>L2: Native | L1: Native<br>L2: ERC20 | Tokamak L2에서는 ERC20 |
| **TON Lock 위치** | - | OptimismPortal | - |
| **ETH Lock 위치** | OptimismPortal | L1StandardBridge | 위치가 다름 |
| **ETHLockbox** | ✅ 사용 | ❌ 사용 안 함 | Tokamak은 불필요 |
| **DelayedWETH** | ✅ 사용 | ✅ 사용 | 게임 참여 시 ETH를 WETH로 bond 제공 |

### 7.2 수정 완료 사항
1. ✅ `ETHLockbox` 배포 스킵 (address(0) 설정)
2. ✅ `ETHLockbox` 검증 로직 주석 처리
3. ✅ `DelayedWETH` 배포 확인 (정상 작동)
4. ✅ 토큰 아키텍처 분석 완료

### 7.3 다음 단계
1. ⬜ 컨트랙트 컴파일 확인
2. ⬜ E2E 테스트 재실행
3. ⬜ 에러 발생 시 추가 디버깅
4. ⬜ 필요시 다른 Implementation 컨트랙트 생성자 확인

## 7. 참고 파일

- `packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol` - TON 컨트랙트
- `packages/tokamak/contracts-bedrock/src/L1/L1StandardBridge.sol` - TON 브릿지 로직
- `packages/tokamak/contracts-bedrock/src/L1/SystemConfig.sol` - Native Token 설정
- `packages/tokamak/contracts-bedrock/src/L1/OptimismPortal2.sol` - Portal 로직

