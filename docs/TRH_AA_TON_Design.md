# ERC-4337 Account Abstraction — Thanos TON Native 환경 상세 설계서

> **대상 레포**: `tokamak-thanos` (Gaming/Full Preset)
> **기반 코드**: eth-infinitism/account-abstraction v0.7.0
> **작성일**: 2026.03.16 | Tokamak Network | v1.0

---

## 1. 왜 AA가 필요한가

Thanos에서 TON이 이미 네이티브 토큰이므로 EOA로 TX를 보낼 수 있다. AA의 가치는 "TON을 쓰는 것"이 아니라 **"EOA의 한계를 제거하는 것"**이다.

### 1.1 EOA vs Smart Account

| 기능 | EOA + TON Native | AA Smart Account |
|------|-----------------|-----------------|
| 가스비 지불 | 유저가 직접 TON으로 지불 | Paymaster가 대납 가능 (gasless) |
| TX 서명 | 매 TX마다 Private Key 서명 필수 | 세션키로 일정 기간/범위 내 무서명 TX |
| 다중 호출 | 각각 별도 TX + 각각 가스비 | `executeBatch()`로 1 TX로 묶음 |
| 지갑 복구 | Private Key 분실 = 자산 상실 | Social recovery, 멀티시그 가능 |
| 가스비 토큰 | TON만 가능 | L2 ETH, USDC 등 ERC-20으로 가능 |
| 온보딩 | 지갑 설치 → TON 구매 → 가스비 납부 | 소셜 로그인 → 즉시 사용 |

### 1.2 Gaming Preset에서의 핵심 시나리오

**Gasless TX**: 게임 유저가 메타마스크를 열고 TON을 사서 가스비를 내는 것은 비현실적이다. Paymaster가 가스비를 대납하면 유저는 지갑만 연결하고 바로 게임을 플레이한다. 게임 운영사가 Paymaster에 TON을 넣어두고 유저의 가스비를 대신 내주는 구조로, Web2 게임처럼 유저는 가스비를 인식하지 못한다.

**Session Key**: EOA는 모든 TX마다 Private Key 서명이 필요하다. 게임에서 칼을 휘두를 때마다 메타마스크 팝업이 뜨면 아무도 안 한다. AA의 Smart Account는 "세션키"를 발급할 수 있다 — 1시간 동안, 이 게임 컨트랙트에 대해서만, 가스비 상한 0.1 TON까지 서명 없이 TX를 보낼 수 있는 임시 키. EOA로는 불가능하다.

**Batch TX**: 게임에서 "아이템 구매 → 장착 → 스탯 업데이트"는 3개 컨트랙트 호출이다. EOA는 3개 TX를 각각 보내고 각각 가스비를 내지만, AA Smart Account는 `executeBatch()`로 3개를 하나의 TX로 묶는다.

**요약**: AA가 없으면 Thanos L2는 "TON으로 가스비 내는 일반 L2"이고, AA가 있으면 "게임 유저가 블록체인을 인식하지 못하는 L2"가 된다. 이 차이가 Gaming Preset의 존재 이유다.

---

## 2. 핵심 문제: ETH → TON 인터페이스

ERC-4337의 모든 컨트랙트는 `msg.value = ETH`를 전제로 설계되어 있다. Thanos에서는 `msg.value = TON`이므로, ETH와 관련된 모든 deposit/withdraw/gas 결제 로직을 점검해야 한다.

### 2.1 Thanos 네이티브 토큰 매핑

| 항목 | Ethereum (원본) | Thanos |
|------|----------------|--------|
| msg.value | ETH | TON |
| Native Token | ETH (18 decimals) | TON (18 decimals) |
| deposit() | ETH payable | TON payable |
| withdrawTo() | ETH transfer | TON transfer |
| Gas fee 정산 | ETH 기준 gasCost | TON 기준 gasCost |
| Paymaster deposit | ETH stake | TON stake |
| Bundler 보상 | ETH (actualGasCost) | TON (actualGasCost) |
| WNativeToken | WETH | WTON (0x4200…0006) |
| ETH on L2 | 네이티브 | ERC-20 (0x4200…0031) |

### 2.2 핵심 인사이트: 변경이 의외로 적다

Thanos에서 `msg.value`가 이미 TON이다. 따라서 **StakeManager의 deposit/withdraw 로직은 코드 변경 없이 그대로 작동**한다. `msg.value`로 받고 `payable(addr).call{value: amount}("")`로 보내는 패턴은 네이티브 토큰이 무엇이든 동일하게 작동하기 때문이다.

진짜 변경이 필요한 곳은 3가지뿐이다:

1. **TonTokenPaymaster** (신규): L2 ETH(ERC-20)로 가스비를 받고, TON deposit에서 대납하는 정산 로직
2. **Bundler** (설정): 가스 가격 오라클을 TON 기준으로 변경
3. **프론트엔드/SDK** (설정): 가스 추정 호출의 단위 인식

---

## 3. 컨트랙트별 변경 상세

### 3.1 StakeManager.sol — 변경 없음

StakeManager는 `msg.value`를 직접 다루는 유일한 ERC-4337 코어 컨트랙트이다. Thanos에서 `msg.value = TON`이므로, deposit/withdraw/addStake 모든 함수가 코드 변경 없이 TON을 처리한다.

```solidity
// depositTo: msg.value로 deposit — Thanos에서 msg.value = TON
function depositTo(address account) public virtual payable {
    deposits[account].deposit += msg.value;  // ← TON이 들어옴
    emit Deposited(account, deposits[account].deposit);
}

// receive: 직접 전송 시 자동 deposit
receive() external payable {
    depositTo(msg.sender);  // ← TON 수신
}

// withdrawTo: 네이티브 토큰 인출
function withdrawTo(address payable withdrawAddress, uint256 amount) external {
    DepositInfo storage info = deposits[msg.sender];
    require(amount <= info.deposit, "Withdraw amount too large");
    info.deposit -= amount;
    (bool success,) = withdrawAddress.call{value: amount}("");  // ← TON 전송
    require(success, "failed to withdraw");
    emit Withdrawn(msg.sender, withdrawAddress, amount);
}

// addStake: Paymaster/Factory가 TON으로 스테이킹
function addStake(uint32 unstakeDelaySec) public payable {
    DepositInfo storage info = deposits[msg.sender];
    info.stake += msg.value;  // ← TON이 스테이크됨
    ...
}
```

**왜 변경이 필요 없는가**: op-geth의 state transition 레이어에서 네이티브 토큰이 TON으로 매핑되었으므로, Solidity의 `payable`, `msg.value`, `address.transfer()`, `address.call{value: x}("")` 등 모든 네이티브 토큰 관련 opcode가 TON 기준으로 작동한다.

### 3.2 EntryPoint.sol — 최소 변경 (주석/에러메시지만)

EntryPoint의 핵심 로직(handleOps, innerHandleOp)은 `gasCost = gasUsed * gasPrice`로 계산한다. Thanos에서 `tx.gasprice`는 이미 TON 단위이므로, 가스비 계산 자체는 변경이 필요 없다.

**변경 대상: 주석과 에러 메시지의 ETH 참조 정리**

```diff
- /// @dev compensate the caller's beneficiary with the collected ETH fees
+ /// @dev compensate the caller's beneficiary with the collected native token (TON) fees
  function _compensate(address payable beneficiary, uint256 amount) internal {
-     require(amount > 0, "AA90 ETH amount is zero");
+     require(amount > 0, "AA90 native token amount is zero");
      ...
  }
```

가스비 정산 플로우 (변경 없이 작동):

```
1. Bundler가 handleOps(ops, beneficiary) 호출
   → Bundler는 TON으로 L2 가스비를 지불

2. 각 UserOp에 대해:
   a. validateUserOp: Account가 missingAccountFunds를 EntryPoint에 TON으로 전송
   b. executeUserOp: Account의 callData 실행
   c. 가스 정산: actualGasCost = actualGas * tx.gasprice (TON 단위)

3. Paymaster가 있는 경우:
   a. validatePaymasterUserOp: Paymaster의 TON deposit 확인
   b. postOp: Paymaster의 TON deposit에서 가스비 차감

4. _compensate: 남은 TON을 beneficiary(Bundler)에게 전송

⇒ 모든 단계에서 네이티브 토큰(TON)이 자동으로 사용됨
```

### 3.3 SimpleAccount.sol — 변경 없음

SimpleAccount는 `validateUserOp`에서 EntryPoint에 `missingAccountFunds`를 전송한다. 이 전송은 `assembly { pop(call(..., missingAmount, ...)) }`로 네이티브 토큰을 보내므로, Thanos에서 자동으로 TON이 전송된다.

```solidity
// _payPrefund: EntryPoint에 가스비 선납 — 변경 불필요
function _payPrefund(uint256 missingAccountFunds) internal virtual {
    if (missingAccountFunds != 0) {
        (bool success,) = payable(msg.sender).call{
            value: missingAccountFunds, gas: type(uint256).max
        }("");  // ← Thanos에서 TON 전송
        (success);
    }
}

// execute: 외부 호출 — value도 TON
function execute(address dest, uint256 value, bytes calldata func) external {
    _requireFromEntryPointOrOwner();
    _call(dest, value, func);  // ← value = TON
}

// addDeposit: EntryPoint에 TON deposit
function addDeposit() public payable {
    entryPoint().depositTo{value: msg.value}(address(this));  // ← TON
}
```

### 3.4 VerifyingPaymaster.sol — 최소 변경

서명 기반 가스비 대납. 오프체인 서버가 UserOp에 서명하면 Paymaster가 TON deposit에서 가스비를 지불한다. deposit이 TON이므로 로직 변경 없이 작동한다. 주석만 정리.

### 3.5 TonTokenPaymaster.sol — 신규 컨트랙트 (핵심 변경)

**이 컨트랙트가 전체 설계의 핵심이다.** 사용자가 L2 ETH(ERC-20)로 가스비를 지불하고, Paymaster가 TON deposit에서 대납하는 패턴이다.

#### 3.5.1 왜 필요한가

Gaming Preset의 핵심 시나리오: 게임 유저가 ETH만 가지고 있고 TON은 없는 상태에서 게임을 플레이할 수 있어야 한다. TonTokenPaymaster가 유저의 L2 ETH를 받고, TON으로 가스비를 대신 내준다.

#### 3.5.2 컨트랙트 구조

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "../core/BasePaymaster.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title TonTokenPaymaster
/// @notice Thanos L2 전용. 사용자가 L2 ETH(ERC-20)로 가스비를 지불.
///         Paymaster는 TON(네이티브)으로 EntryPoint에 deposit하고,
///         사용자의 L2 ETH를 수거하여 정산.
contract TonTokenPaymaster is BasePaymaster {

    // L2 ETH ERC-20 (Predeploy: 0x4200000000000000000000000000000000000031)
    IERC20 public immutable l2Eth;

    // TON/ETH 환율 오라클
    ITokenPriceOracle public oracle;

    // 마크업: 10% = 가스비의 10%를 수수료로 수거
    uint256 public markupPercent = 10;

    constructor(
        IEntryPoint _entryPoint,
        IERC20 _l2Eth,
        ITokenPriceOracle _oracle
    ) BasePaymaster(_entryPoint) {
        l2Eth = _l2Eth;
        oracle = _oracle;
    }

    /// @notice 검증 단계: 사용자의 L2 ETH allowance 확인 + 사전 수거
    function _validatePaymasterUserOp(
        PackedUserOperation calldata userOp,
        bytes32 userOpHash,
        uint256 maxCost          // ← TON 단위 (네이티브)
    ) internal override returns (bytes memory context, uint256 validationData) {

        address sender = userOp.getSender();

        // 1. maxCost(TON) → L2 ETH로 환산
        uint256 ethCost = _tonToEth(maxCost);
        uint256 ethCostWithMarkup = ethCost * (100 + markupPercent) / 100;

        // 2. 사용자의 L2 ETH allowance 확인
        uint256 allowance = l2Eth.allowance(sender, address(this));
        require(allowance >= ethCostWithMarkup, "PM: insufficient ETH allowance");

        // 3. L2 ETH 사전 수거 (pre-charge)
        l2Eth.transferFrom(sender, address(this), ethCostWithMarkup);

        // 4. context에 사전 수거 금액 기록 (postOp에서 사용)
        context = abi.encode(sender, ethCostWithMarkup, maxCost);
        validationData = 0;
    }

    /// @notice 실행 후 정산: 실제 가스비만큼 수거하고 나머지 환불
    function _postOp(
        PostOpMode mode,
        bytes calldata context,
        uint256 actualGasCost,   // ← TON 단위
        uint256 actualUserOpFeePerGas
    ) internal override {
        (address sender, uint256 preCharged, ) =
            abi.decode(context, (address, uint256, uint256));

        // 1. 실제 가스비(TON) → L2 ETH로 환산
        uint256 actualEthCost = _tonToEth(actualGasCost);
        uint256 actualWithMarkup = actualEthCost * (100 + markupPercent) / 100;

        // 2. 과다 수거분 환불
        if (preCharged > actualWithMarkup) {
            l2Eth.transfer(sender, preCharged - actualWithMarkup);
        }
    }

    /// @notice TON → L2 ETH 환율 변환
    function _tonToEth(uint256 tonAmount) internal view returns (uint256) {
        uint256 price = oracle.getPrice(); // 1 TON당 ETH 가격 (18 decimals)
        return tonAmount * price / 1e18;
    }
}
```

#### 3.5.3 환율 오라클

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

/// @title ITokenPriceOracle
/// @notice TON/ETH 환율을 제공하는 오라클 인터페이스
interface ITokenPriceOracle {
    /// @notice 1 TON당 ETH 가격 (18 decimals)
    /// @return price 예: 0.0005e18 = 1 TON = 0.0005 ETH
    function getPrice() external view returns (uint256 price);

    /// @notice 가격 업데이트 시간
    function lastUpdated() external view returns (uint256);
}
```

**Phase 1 구현 — SimplePriceOracle**: Owner가 수동 업데이트. 24시간 미갱신 시 stale price revert.

```solidity
contract SimplePriceOracle is ITokenPriceOracle {
    uint256 public override lastUpdated;
    uint256 private _price;
    address public owner;

    constructor(uint256 initialPrice) {
        owner = msg.sender;
        _price = initialPrice;
        lastUpdated = block.timestamp;
    }

    function getPrice() external view override returns (uint256) {
        require(block.timestamp - lastUpdated < 86400, "stale price");
        return _price;
    }

    function updatePrice(uint256 newPrice) external {
        require(msg.sender == owner, "only owner");
        _price = newPrice;
        lastUpdated = block.timestamp;
    }
}
```

**Phase 2 확장 — UniswapV3TwapOracle**: DeFi Preset의 Uniswap V3 predeploy에서 TON/ETH 풀의 TWAP을 자동으로 읽는 오라클. Full Preset에서 DeFi + Gaming predeploy가 모두 있으면 이 패턴을 사용.

```
TON/WTON Pool → TWAP → oracle.getPrice()
별도 외부 트리거 불필요. 온체인 자동화.
```

#### 3.5.4 정산 플로우 다이어그램

```
사용자 (L2 ETH 보유, TON 없음)
    │
    ├─ 1. L2 ETH approve → TonTokenPaymaster
    │
    ├─ 2. UserOp 생성 (paymasterAndData = Paymaster 주소)
    │
    └─ 3. Bundler에게 UserOp 제출
            │
            ▼
      Bundler: handleOps(ops, beneficiary)
            │
            ├─ 4. validatePaymasterUserOp()
            │     ├─ maxCost(TON) → L2 ETH 환산
            │     ├─ l2Eth.transferFrom(sender, paymaster, ethCostWithMarkup)
            │     └─ context = (sender, preCharged, maxCost)
            │
            ├─ 5. executeUserOp() — 사용자의 callData 실행
            │
            ├─ 6. _postOp()
            │     ├─ actualGasCost(TON) → L2 ETH 환산
            │     ├─ 과다 수거분 → l2Eth.transfer(sender, refund)
            │     └─ Paymaster의 TON deposit에서 가스비 차감 (EntryPoint가 처리)
            │
            └─ 7. _compensate(beneficiary, collected)
                  └─ Bundler에게 TON 보상 전송
```

#### 3.5.5 감사 필요 포인트

| 포인트 | 상세 | 심각도 |
|--------|------|--------|
| Pre-charge & Refund 패턴 | transferFrom 후 postOp에서 환불. PostOpMode.opReverted 시 사용자 자금 동결 가능성. revert 케이스 전수 검증 필요. | High |
| 환율 오라클 조작 | 가격 갱신 지연 시 arbitrage 가능. stale price 체크 필수. 향후 Uniswap V3 TWAP으로 교체 권장. | High |
| ERC-20 reentrancy | l2Eth.transferFrom → _postOp 사이에 reentrancy 가능성. L2 ETH는 표준 ERC-20이므로 위험 낮지만 check-effects-interaction 패턴 준수. | Medium |

---

## 4. Bundler 변경

Bundler는 UserOp을 수집하여 `EntryPoint.handleOps()`를 호출하는 오프체인 서비스이다. Thanos에서 가스 가격은 TON 단위이므로, 설정 변경이 필요하다.

### 4.1 설정 변경

```typescript
// bundler.config.ts
export const thanosConfig = {
  // 가스 가격: TON 단위 — Thanos RPC가 이미 TON 기준 반환
  minGasPrice: parseUnits("0.001", "gwei"),

  // 수익성 계산: TON 기준
  minBundleProfit: parseUnits("0.01", "ether"),  // 0.01 TON

  // EntryPoint: Genesis predeploy 주소
  entryPoint: "0x4200000000000000000000000000000000000063",

  // Beneficiary: Bundler의 TON 수령 주소
  beneficiary: process.env.BUNDLER_ADDRESS,
};
```

### 4.2 변경 체크리스트

| 항목 | 상태 |
|------|------|
| 가스 가격 단위를 TON으로 인식 | 자동 (RPC가 이미 TON 반환) |
| 수익성 계산 기준을 TON으로 변경 | config 변경 필요 |
| EntryPoint 주소를 predeploy 주소로 변경 | config 변경 필요 |
| `eth_estimateGas` → TON 가스비 반환 인식 | 자동 (op-geth 레벨) |
| ERC-4337 RPC API (`eth_sendUserOperation` 등) | 변경 없음 (스펙 동일) |

---

## 5. SDK / 프론트엔드 변경

### 5.1 UserOp 생성

```typescript
class TonUserOpBuilder {
  async buildUserOp(params: {
    sender: string;
    callData: string;
    paymaster?: string;
  }): Promise<UserOperation> {

    // 현재 TON 가스 가격 조회 — RPC가 이미 TON 단위 반환
    const feeData = await provider.getFeeData();

    return {
      sender: params.sender,
      nonce: await entryPoint.getNonce(params.sender, 0),
      callData: params.callData,
      maxFeePerGas: feeData.maxFeePerGas,             // TON gwei
      maxPriorityFeePerGas: feeData.maxPriorityFeePerGas,  // TON gwei
      paymasterAndData: params.paymaster
        ? await getPaymasterSignature(params.paymaster)
        : "0x",
    };
  }
}
```

### 5.2 TonTokenPaymaster 사용 시 사전 approve

```typescript
// 사용자가 TonTokenPaymaster를 쓰려면 L2 ETH approve가 필요
const l2Eth = new Contract(L2_ETH_ADDRESS, ERC20_ABI, signer);
await l2Eth.approve(PAYMASTER_ADDRESS, ethers.MaxUint256);

// 이후 UserOp에 paymasterAndData만 설정하면 gasless TX 가능
```

---

## 6. Genesis Predeploy 배치

Gaming/Full Preset의 genesis에 다음 3개 컨트랙트를 predeploy로 탑재한다.

| 컨트랙트 | Predeploy 주소 | Proxy | Admin |
|----------|---------------|-------|-------|
| EntryPoint.sol | `0x4200...0063` | Transparent Proxy | SystemConfig owner |
| VerifyingPaymaster.sol | `0x4200...0064` | Transparent Proxy | SystemConfig owner |
| SimpleAccountFactory.sol | `0x4200...0065` | Transparent Proxy | SystemConfig owner |

TonTokenPaymaster와 SimplePriceOracle은 genesis predeploy가 아닌 **배포 후 모듈**로 제공한다. 이유: 환율 오라클은 운영사마다 다른 설정이 필요하고, 마크업 비율도 커스터마이징이 필요하기 때문이다.

```
Genesis Predeploy (Gaming Preset):
  ├── EntryPoint          — ERC-4337 코어. 모든 체인에 동일.
  ├── VerifyingPaymaster  — 서명 기반 gasless. 범용.
  └── SimpleAccountFactory — Smart Account 생성. 범용.

배포 후 모듈 (trh-sdk integrate):
  ├── TonTokenPaymaster   — L2 ETH로 가스비 지불. 운영사 커스터마이징.
  └── SimplePriceOracle   — TON/ETH 환율. 운영사 관리.
```

---

## 7. 변경 파일 총정리

| 파일 | 변경 수준 | 내용 | 감사 필요 |
|------|----------|------|----------|
| StakeManager.sol | 없음 | 코드 변경 불필요 | 불필요 |
| NonceManager.sol | 없음 | 코드 변경 불필요 | 불필요 |
| Helpers.sol | 없음 | 코드 변경 불필요 | 불필요 |
| UserOperationLib.sol | 없음 | 코드 변경 불필요 | 불필요 |
| EntryPoint.sol | 최소 | 주석/에러메시지 ETH → native token | 불필요 |
| EntryPointSimulations.sol | 최소 | 주석 ETH 참조 정리 | 불필요 |
| BaseAccount.sol | 없음 | 코드 변경 불필요 | 불필요 |
| BasePaymaster.sol | 없음 | 코드 변경 불필요 | 불필요 |
| SimpleAccount.sol | 최소 | 주석 정리 | 불필요 |
| SimpleAccountFactory.sol | 최소 | 주석 정리 | 불필요 |
| VerifyingPaymaster.sol | 최소 | 주석 정리. 로직 변경 없음 | 불필요 |
| **TonTokenPaymaster.sol** | **신규** | **L2 ETH → TON 가스비 정산. 핵심.** | **필수** |
| **ITokenPriceOracle.sol** | **신규** | **TON/ETH 환율 오라클 인터페이스** | **필수** |
| **SimplePriceOracle.sol** | **신규** | **오라클 기본 구현** | **필수** |
| Bundler config | 설정 | TON 가스 가격 기준 설정 | 불필요 |
| SDK / 프론트엔드 | 설정 | 가스 추정 단위 인식 | 불필요 |

**감사 대상: 신규 3개 파일만. 기존 eth-infinitism 코드는 변경 없음.**

---

## 8. 테스트 시나리오

### T1: 기본 UserOp (Paymaster 없음)
1. SimpleAccount를 Gaming Preset genesis에서 확인
2. Account에 TON deposit (`entryPoint.depositTo`)
3. UserOp 생성: callData = ERC-20 transfer
4. Bundler가 handleOps 호출
5. **검증**: Account의 TON deposit이 가스비만큼 차감
6. **검증**: Bundler beneficiary에 TON이 보상으로 전송

### T2: VerifyingPaymaster (서명 기반 gasless)
1. VerifyingPaymaster에 TON deposit
2. 사용자가 TON 없이 UserOp 생성
3. 오프체인 서버가 paymasterAndData에 서명
4. handleOps 호출 → Paymaster의 TON deposit에서 가스비 차감
5. **검증**: 사용자의 TON 잔액 변화 없음 (gasless)

### T3: TonTokenPaymaster (L2 ETH로 가스비 지불)
1. TonTokenPaymaster에 TON deposit + oracle 설정
2. 사용자가 L2 ETH를 보유하고 TON은 0인 상태
3. 사용자가 L2 ETH를 Paymaster에 approve
4. UserOp 생성 (paymasterAndData = Paymaster 주소)
5. handleOps → Paymaster가 TON으로 가스비 대납
6. **검증**: 사용자의 L2 ETH가 가스비 + 마크업만큼 차감
7. **검증**: Paymaster에 L2 ETH가 수거됨
8. **검증**: Paymaster의 TON deposit이 가스비만큼 차감

### T4: 오라클 장애 시 failsafe
1. `oracle.updatePrice`를 24시간 이상 안 함
2. TonTokenPaymaster를 통한 UserOp 제출
3. **검증**: `validatePaymasterUserOp`에서 "stale price" revert
4. **검증**: 사용자의 L2 ETH가 차감되지 않음 (안전)

### T5: Bundler 수익성
1. 10개 UserOp 배치 → handleOps 호출
2. TX receipt에서 actualGasCost 합산
3. Bundler beneficiary의 TON 잔액 변화 확인
4. **검증**: 수익 = 수거 TON - TX 가스비 > 0

### T6: 호환성 — 기존 OP Stack E2E 테스트 통과
1. Gaming Preset genesis로 체인 구동
2. 기존 OP Stack E2E 테스트(deposit/withdraw/fault proof) 실행
3. **검증**: AA predeploy 추가가 기존 기능에 영향 없음

---

## 9. 구현 일정

| 단계 | 기간 | 작업 | 전제조건 |
|------|------|------|----------|
| 9-1 | 1주 | EntryPoint + VerifyingPaymaster + SimpleAccountFactory를 genesis predeploy로 탑재. 주석 정리. | Phase 1 (VRF predeploy) 완료 |
| 9-2 | 2주 | TonTokenPaymaster + SimplePriceOracle 개발. 단위 테스트. | 9-1 완료 |
| 9-3 | 1주 | Bundler 설정 변경 + SDK 가스 추정 연동 | 9-2 완료 |
| 9-4 | 2주 | E2E 테스트 T1~T6 전수 실행. 감사 준비. | 9-3 완료 |
| 9-5 | 병렬 | TonTokenPaymaster + Oracle 외부 감사 | 9-2 완료 시 시작 |

**총 예상: 6주 (감사 병렬 진행)**

---

## 10. 위험 및 완화

| 위험 | 심각도 | 완화 방안 |
|------|--------|----------|
| TonTokenPaymaster의 pre-charge 실패 시 자금 동결 | High | PostOpMode별 환불 로직 전수 검증. revert 시에도 사전 수거분이 환불되도록 try-catch 패턴 적용. |
| 환율 오라클 조작 / 지연 | High | stale price 체크(24h). Phase 2에서 Uniswap V3 TWAP으로 교체. 마크업 10%가 환율 변동 버퍼 역할. |
| ERC-20 reentrancy | Medium | L2 ETH는 표준 OptimismMintableERC20이므로 hook 없음. 그래도 check-effects-interaction 패턴 준수. |
| eth-infinitism 업스트림 업데이트 | Low | 기존 코드 변경이 거의 없으므로(주석만) upstream merge 용이. TonTokenPaymaster는 독립 컨트랙트라 충돌 없음. |
| Bundler 호환성 | Low | ERC-4337 RPC 스펙 변경 없음. 가스 단위만 TON으로 바뀌고 API는 동일. 기존 Bundler(Stackup, Pimlico 등) 설정만 변경하면 사용 가능. |

---

## 부록 A: 파일 디렉토리 구조

```
packages/tokamak/contracts-bedrock/src/tokamak-contracts/AA/
├── EntryPoint.sol                  ← eth-infinitism v0.7 포팅 (주석만 변경)
├── VerifyingPaymaster.sol          ← eth-infinitism v0.7 포팅 (주석만 변경)
├── SimpleAccount.sol               ← eth-infinitism v0.7 포팅 (주석만 변경)
├── SimpleAccountFactory.sol        ← eth-infinitism v0.7 포팅 (주석만 변경)
├── TonTokenPaymaster.sol           ← 신규 (핵심)
├── interfaces/
│   └── ITokenPriceOracle.sol       ← 신규
├── oracles/
│   ├── SimplePriceOracle.sol       ← 신규 (Phase 1)
│   └── UniswapV3TwapOracle.sol     ← 신규 (Phase 2, Full Preset)
└── lib/
    ├── BaseAccount.sol             ← eth-infinitism (변경 없음)
    ├── BasePaymaster.sol           ← eth-infinitism (변경 없음)
    ├── StakeManager.sol            ← eth-infinitism (변경 없음)
    ├── NonceManager.sol            ← eth-infinitism (변경 없음)
    └── Helpers.sol                 ← eth-infinitism (변경 없음)
```

## 부록 B: Predeploy 주소 할당

```
0x4200000000000000000000000000000000000060  VRFPredeploy      (Gaming/Full)
0x4200000000000000000000000000000000000061  VRFCoordinator     (Gaming/Full)
0x4200000000000000000000000000000000000062  VRFConsumerBase    (Gaming/Full)
0x4200000000000000000000000000000000000063  EntryPoint         (Gaming/Full)
0x4200000000000000000000000000000000000064  VerifyingPaymaster (Gaming/Full)
0x4200000000000000000000000000000000000065  SimpleAccountFactory (Gaming/Full)
```

---

*End of Document*
