# ERC-4337 + EIP-7702 Account Abstraction — Thanos TON Native 환경 상세 설계서 (v2)

> **대상 레포**: `tokamak-thanos`, `tokamak-thanos-geth`
> **기반 코드**: eth-infinitism/account-abstraction v0.8+ / EIP-7702 (Pectra)
> **작성일**: 2026.03.16 | Tokamak Network | v2.0
> **변경 사유**: tokamak-thanos-geth에 Pectra(Isthmus) 업그레이드 도입 확인. EIP-7702가 L2 EVM 레벨에서 네이티브 지원됨에 따라 AA 전략 전면 수정.

---

## 변경 이력

| 버전 | 날짜 | 변경 내용 |
|------|------|----------|
| v1.0 | 2026.03.16 | 최초 작성. ERC-4337 v0.7 기반. EIP-7702를 "OP Stack이 지원하면 나중에 도입" 전략. |
| **v2.0** | **2026.03.16** | **전면 수정. tokamak-thanos-geth Pectra 도입 확인. EIP-7702를 기본 전략으로 채택. EntryPoint v0.8+로 업그레이드. Simple7702Account 중심 설계.** |

---

## 1. 전략 변경 배경

### 1.1 무엇이 바뀌었는가

v1.0 작성 시점에서는 EIP-7702를 "Level 3 (EVM 실행 레이어 변경) = 위험/비추천" 영역으로 분류했다. Thanos가 OP Stack v1.7.7 기반이고, EIP-7702의 `SET_CODE_TX_TYPE (0x04)` 트랜잭션을 지원하려면 op-geth를 직접 포크해야 한다고 판단했기 때문이다.

그러나 다음이 확인되었다:

1. **OP Stack Isthmus 하드포크**가 Pectra의 L2 관련 기능(EIP-7702 포함)을 공식 도입
2. **Optimism op-geth**가 EIP-7702 트랜잭션 타입을 네이티브 지원
3. **tokamak-thanos-geth**가 이 Pectra/Isthmus 업그레이드를 도입

따라서 EIP-7702는 더 이상 "Thanos 자체 포크"가 아니라 **OP Stack 업스트림에서 제공하는 표준 기능**이다. Level 3 변경이 아니라 Level 0 (업스트림 머지)이다.

### 1.2 전략 수정 요약

| 항목 | v1.0 (이전) | v2.0 (수정) |
|------|------------|------------|
| EVM 레벨 | EIP-7702 미지원 전제 | EIP-7702 네이티브 지원 |
| EntryPoint 버전 | v0.7 | **v0.8+** (7702 통합) |
| 유저 계정 | SimpleAccount (새 주소) | **Simple7702Account (기존 EOA 유지)** |
| 계정 생성 | SimpleAccountFactory로 CREATE2 배포 | **EOA에 7702 delegation 설정** (배포 불필요) |
| 유저 주소 | 새 Smart Account 주소 | **기존 EOA 주소 그대로** |
| 기존 자산 | Smart Account로 이전 필요 | **이전 불필요** |
| Paymaster | 동일 | 동일 (여전히 필요) |
| Bundler | 동일 | 동일 (여전히 필요) |
| TonTokenPaymaster | 동일 | 동일 (여전히 필요) |
| genesis predeploy | EntryPoint + VerifyingPaymaster + SimpleAccountFactory | **EntryPoint v0.8 + VerifyingPaymaster + Simple7702Account** |

---

## 2. 왜 AA가 필요한가 (변경 없음)

Thanos에서 TON이 이미 네이티브 토큰이므로 EOA로 TX를 보낼 수 있다. AA의 가치는 "TON을 쓰는 것"이 아니라 **"EOA의 한계를 제거하는 것"**이다.

| 기능 | EOA만 사용 | EIP-7702 + ERC-4337 |
|------|-----------|---------------------|
| 가스비 지불 | 유저가 직접 TON으로 | Paymaster가 대납 (gasless) |
| TX 서명 | 매 TX마다 Private Key 서명 | 세션키로 무서명 TX |
| 다중 호출 | 각각 별도 TX | `executeBatch()`로 1 TX |
| 지갑 복구 | Private Key 분실 = 자산 상실 | Social recovery 가능 |
| 가스비 토큰 | TON만 가능 | L2 ETH, USDC 등 ERC-20 가능 |
| 유저 주소 | EOA 주소 | **동일한 EOA 주소 (7702 덕분)** |

Gaming Preset의 핵심 시나리오는 v1.0과 동일하다: gasless TX, session key, batch TX. 달라지는 것은 **유저가 새 Smart Account 주소를 만들 필요 없이 기존 EOA에서 바로 이 모든 기능을 쓸 수 있다**는 점이다.

---

## 3. EIP-7702가 바꾸는 것과 바꾸지 않는 것

### 3.1 EIP-7702가 바꾸는 것

**EOA가 Smart Account처럼 동작한다.** `SET_CODE_TX_TYPE (0x04)` 트랜잭션으로 EOA에 delegation target 컨트랙트 코드를 설정하면, 해당 EOA는 그 컨트랙트의 로직을 실행할 수 있다. 이것은 영구적이 아니라 유저가 언제든 해제할 수 있다.

```
// EIP-7702 delegation 설정
// 유저의 EOA (0xAbcd...) → Simple7702Account 컨트랙트 코드를 위임
//
// 이후 EOA가 직접 executeBatch(), validateUserOp() 등을 실행 가능
// 주소는 그대로 0xAbcd...
```

구체적 변화:

- **계정 배포 불필요**: SimpleAccountFactory로 CREATE2 배포하는 단계가 사라짐. 7702 delegation TX 한 번이면 끝.
- **주소 유지**: 유저의 기존 EOA 주소에 이미 있는 TON, L2 ETH, NFT 등이 그대로 유지됨. Smart Account로 자산 이전 불필요.
- **가스비 절감**: Smart Account 배포 가스(~300K gas)가 사라짐. 7702 delegation TX는 ~50K gas.
- **UX 단순화**: "새 지갑 주소가 생긴다"는 혼란이 사라짐.

### 3.2 EIP-7702가 바꾸지 않는 것

**Paymaster와 Bundler는 여전히 필요하다.** 7702는 "EOA가 코드를 실행할 수 있게 해주는 것"이지, "누가 가스비를 내는가"를 해결하지 않는다.

| 기능 | EIP-7702만으로 가능? | 추가로 필요한 것 |
|------|-------------------|----------------|
| Batch TX | **가능** (delegation target의 executeBatch) | 없음 |
| Session Key | **가능** (delegation target의 validateUserOp) | 없음 |
| Gasless TX | **불가능** — 누군가 가스비를 내야 함 | **Paymaster + EntryPoint** |
| L2 ETH로 가스비 지불 | **불가능** — 네이티브 토큰(TON)만 가스비로 사용 가능 | **TonTokenPaymaster** |
| Bundling | **불가능** — UserOp을 모아서 제출하는 오프체인 서비스 필요 | **Bundler** |

따라서 EIP-7702 + ERC-4337의 조합이 최종 아키텍처다. 7702가 "계정"을 담당하고, 4337이 "가스비 추상화"를 담당한다.

---

## 4. 아키텍처 개요

### 4.1 컴포넌트 구성

```
┌─────────────────────────────────────────────────────────┐
│                    Gaming Preset L2                       │
│                                                           │
│  ┌──────────────────────────────────────────────────┐    │
│  │  Genesis Predeploy (tokamak-thanos)              │    │
│  │                                                    │    │
│  │  0x4200...0063  EntryPoint v0.8+                  │    │
│  │  0x4200...0064  VerifyingPaymaster                │    │
│  │  0x4200...0065  Simple7702Account (delegation)    │    │
│  └──────────────────────────────────────────────────┘    │
│                                                           │
│  ┌──────────────────────────────────────────────────┐    │
│  │  배포 후 모듈 (trh-sdk integrate)                 │    │
│  │                                                    │    │
│  │  TonTokenPaymaster  (L2 ETH 가스비 정산)          │    │
│  │  SimplePriceOracle  (TON/ETH 환율)                │    │
│  └──────────────────────────────────────────────────┘    │
│                                                           │
│  ┌──────────────────────────────────────────────────┐    │
│  │  EVM 레벨 (tokamak-thanos-geth, Pectra)          │    │
│  │                                                    │    │
│  │  SET_CODE_TX_TYPE (0x04) 네이티브 지원            │    │
│  │  EIP-7702 authorization list 처리                 │    │
│  │  tx.origin / EXTCODESIZE 시맨틱 변경              │    │
│  └──────────────────────────────────────────────────┘    │
│                                                           │
│  ┌──────────────────────────────────────────────────┐    │
│  │  오프체인 (Helm/K8s 모듈)                         │    │
│  │                                                    │    │
│  │  Bundler (ERC-4337 RPC + 7702 UserOp 처리)       │    │
│  │  Paymaster Signer (서명 서비스)                    │    │
│  └──────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

### 4.2 유저 플로우 (v1.0 vs v2.0)

**v1.0 (이전 — ERC-4337 v0.7만)**
```
1. 유저가 SimpleAccountFactory에 createAccount() 호출 → 새 주소 생성
2. 기존 EOA에서 새 Smart Account로 자산 이전
3. UserOp 생성 (sender = Smart Account 주소)
4. Bundler → EntryPoint → Smart Account 실행
```

**v2.0 (수정 — EIP-7702 + ERC-4337 v0.8)**
```
1. 유저가 7702 delegation TX 전송 → 기존 EOA에 Simple7702Account 코드 위임
2. (자산 이전 불필요 — 기존 EOA 주소와 자산 그대로)
3. UserOp 생성 (sender = 기존 EOA 주소)
4. Bundler → EntryPoint v0.8 → EOA가 Smart Account 로직 실행
```

단계가 줄었고, 유저 경험이 훨씬 단순해졌다.

---

## 5. Thanos TON 환경에서의 네이티브 토큰 매핑 (변경 없음)

v1.0의 핵심 인사이트가 그대로 유효하다: Thanos에서 `msg.value = TON`이므로, StakeManager/EntryPoint의 deposit/withdraw/gas 정산 로직은 코드 변경 없이 TON을 처리한다.

| 항목 | Ethereum (원본) | Thanos |
|------|----------------|--------|
| msg.value | ETH | TON |
| Native Token | ETH (18 decimals) | TON (18 decimals) |
| deposit/withdraw | ETH | TON (자동 작동) |
| Gas fee 정산 | ETH 기준 | TON 기준 (자동 작동) |
| WNativeToken | WETH | WTON (0x4200…0006) |
| ETH on L2 | 네이티브 | ERC-20 (0x4200…0031) |

---

## 6. 컨트랙트별 상세

### 6.1 EntryPoint v0.8+ — 업그레이드

v0.7 → v0.8+로 업그레이드한다. v0.8의 핵심 변경:

- **EIP-7702 호환**: UserOp의 sender가 7702 delegation이 설정된 EOA일 수 있음. EntryPoint가 delegation address를 userOpHash에 포함하여 replay attack 방지.
- **ABI 호환**: v0.8은 v0.7과 ABI 호환. 기존 Account/Paymaster가 수정 없이 작동.
- **paymasterSignature 분리**: `paymasterAndData`에서 Paymaster 서명을 분리. 유저가 UserOp에 서명한 후에도 Paymaster 서명을 추가할 수 있음. UX 개선.

```diff
- // v0.7: userOpHash에 delegation address 미포함
+ // v0.8: userOpHash에 EIP-7702 delegation address 포함
+ //        → 7702 EOA의 delegation이 변경되면 기존 UserOp 무효화 (replay 방지)
```

Thanos TON 환경에서의 추가 변경: v1.0과 동일하게 주석/에러메시지의 ETH 참조를 native token으로 정리. 로직 변경 없음.

### 6.2 Simple7702Account — 새 delegation target (v1.0의 SimpleAccount/Factory 대체)

eth-infinitism v0.8+에 포함된 `Simple7702Account`를 기반으로 한다. 이 컨트랙트가 genesis predeploy로 탑재되어, 유저의 EOA가 delegation하는 대상이 된다.

```solidity
// Simple7702Account — EIP-7702 delegation target
// 유저의 EOA가 이 컨트랙트에 delegation하면,
// EOA 주소에서 직접 다음 기능을 사용할 수 있음:
//
// - validateUserOp(): EntryPoint와 연동하여 가스비 추상화
// - execute(dest, value, data): 단일 호출 실행
// - executeBatch(dests[], values[], datas[]): 배치 호출
// - ERC-1271 서명 검증: 스마트 컨트랙트 서명
//
// 기존 SimpleAccount와 차이점:
// - Proxy/Factory 배포 불필요 (7702 delegation으로 대체)
// - 유저 EOA 주소 유지
// - owner = EOA의 Private Key (자동)
```

**핵심**: Simple7702Account는 **stateless delegation target**이다. 상태(잔액, nonce 등)는 유저의 EOA에 저장되고, 로직만 Simple7702Account에서 가져온다. 따라서 predeploy로 하나만 배포하면 모든 유저가 공유할 수 있다.

### 6.3 VerifyingPaymaster — 최소 변경 (v1.0과 동일)

서명 기반 gasless TX. TON deposit에서 가스비 대납. 7702 도입과 무관하게 동일하게 작동한다.

유일한 차이: UserOp의 sender가 이제 EOA 주소이므로, Paymaster가 서명 검증 시 sender 주소를 확인하는 로직은 변경 없이 작동한다 (주소가 달라지는 것이 아니라 같은 주소에서 코드가 실행되는 것이므로).

### 6.4 TonTokenPaymaster — 변경 없음 (v1.0과 동일)

L2 ETH(ERC-20)로 가스비를 받고, TON deposit에서 대납하는 정산 로직. EIP-7702와 무관하게 동일하게 작동한다. 전체 구현은 v1.0 문서의 Section 3.5 참조.

핵심 플로우:
```
유저 (L2 ETH 보유, TON 없음)
  → L2 ETH approve → TonTokenPaymaster
  → UserOp 전송 (sender = 기존 EOA, 7702 delegation 활성)
  → Paymaster가 TON으로 가스비 대납
  → 유저의 L2 ETH에서 가스비 + 마크업 수거
```

### 6.5 ITokenPriceOracle / SimplePriceOracle — 변경 없음

v1.0과 동일. TON/ETH 환율 오라클. Phase 1은 owner 수동 업데이트, Phase 2는 Uniswap V3 TWAP.

### 6.6 SimpleAccountFactory — 제거

7702 도입으로 Smart Account 배포가 불필요해졌으므로, **genesis predeploy에서 SimpleAccountFactory를 제거**한다.

기존 ERC-4337 v0.7 방식(Factory로 Smart Account 배포)을 원하는 개발자를 위해 SimpleAccountFactory는 npm 패키지로 제공하되, genesis predeploy에는 포함하지 않는다.

---

## 7. 유저 온보딩 플로우

### 7.1 Gaming 유저 (gasless, 7702)

```
1. 유저가 게임 웹사이트에서 "지갑 연결" (MetaMask/Rabby/etc)
   → EOA 주소 확인

2. 게임이 7702 delegation TX를 생성
   → authorization_list: [{address: Simple7702Account_predeploy}]
   → 유저가 서명 (1회)
   → EOA에 Smart Account 기능이 활성화됨

3. 이후 모든 게임 액션은 UserOp으로 처리
   → Paymaster가 가스비 대납 (유저는 TON 불필요)
   → Session Key로 서명 없이 게임 플레이
   → Batch TX로 여러 액션을 1 TX로 묶음

4. 유저의 EOA 주소에 게임 아이템(NFT), 보상(토큰)이 쌓임
   → 별도 주소로 이전할 필요 없음
```

### 7.2 DeFi 유저 (L2 ETH 가스비)

```
1. 유저가 L2 ETH만 보유하고 TON은 없는 상태
2. 7702 delegation 설정 (1회)
3. L2 ETH를 TonTokenPaymaster에 approve (1회)
4. 이후 모든 DeFi TX는 L2 ETH로 가스비 지불
   → Paymaster가 TON으로 대납하고, 유저의 L2 ETH 수거
```

### 7.3 일반 유저 (EOA 그대로)

7702 delegation을 설정하지 않으면 기존 EOA와 완전히 동일하게 작동한다. AA는 opt-in이다. 기존 Thanos 유저에게 아무런 영향 없음.

---

## 8. Pectra 도입으로 인한 주의사항

### 8.1 EXTCODESIZE 변경

EIP-7702 이후, delegation이 설정된 EOA의 `EXTCODESIZE`가 0이 아닌 값을 반환한다. 기존에 `EXTCODESIZE == 0`으로 "이 주소가 EOA인가?"를 확인하던 컨트랙트는 오작동할 수 있다.

**Thanos에서의 영향**: 기존 Thanos predeploy 컨트랙트(USDC Bridge, Uniswap V3 등)가 이 패턴을 사용하는지 점검 필요. OP Standard predeploy는 Optimism이 이미 검증했으므로 문제 없음.

### 8.2 tx.origin 시맨틱

7702 delegation TX에서 `tx.origin`은 여전히 EOA 주소이다. 7702는 `msg.sender`와 `tx.origin`의 관계를 변경하지 않는다. 다만, `tx.origin == msg.sender`로 "직접 호출인가?"를 확인하던 패턴은 7702에서도 true가 될 수 있으므로, 이 패턴에 의존하는 컨트랙트는 검토 필요.

### 8.3 Delegation 해제

유저는 언제든 7702 delegation을 해제할 수 있다. 해제 후 EOA는 원래의 일반 EOA로 복귀한다. 이것은 7702의 핵심 특성: **비영구적**이다. Smart Account로 "전환"하는 것이 아니라, 임시로 코드를 "위임"하는 것이다.

### 8.4 Nonce 관리

7702 delegation TX는 EOA의 nonce를 소비한다. EntryPoint의 2차원 nonce(key + seq)와 EOA의 L1 nonce가 별개로 관리된다. v0.8+에서는 미배포 Account의 2차원 nonce 사용도 허용(`initCode`가 이미 배포된 경우 무시).

---

## 9. Genesis Predeploy 배치 (수정)

### 9.1 v1.0 → v2.0 변경

| 주소 | v1.0 | v2.0 |
|------|------|------|
| 0x4200…0063 | EntryPoint v0.7 | **EntryPoint v0.8+** |
| 0x4200…0064 | VerifyingPaymaster | VerifyingPaymaster (동일) |
| 0x4200…0065 | SimpleAccountFactory | **Simple7702Account** (delegation target) |

### 9.2 배포 구조

```
Genesis Predeploy (Gaming/Full Preset):
  ├── EntryPoint v0.8+          ← ERC-4337 코어 + 7702 통합
  ├── VerifyingPaymaster        ← 서명 기반 gasless TX
  └── Simple7702Account         ← 7702 delegation target (stateless)

배포 후 모듈 (trh-sdk integrate):
  ├── TonTokenPaymaster         ← L2 ETH로 가스비 지불
  └── SimplePriceOracle         ← TON/ETH 환율
```

---

## 10. 변경 파일 총정리 (v2.0)

| 파일 | 변경 수준 | v1.0 대비 변화 | 내용 | 감사 |
|------|----------|---------------|------|------|
| StakeManager.sol | 없음 | 동일 | 코드 변경 불필요 | 불필요 |
| NonceManager.sol | 없음 | 동일 | 코드 변경 불필요 | 불필요 |
| Helpers.sol | 없음 | 동일 | 코드 변경 불필요 | 불필요 |
| UserOperationLib.sol | 없음 | 동일 | 코드 변경 불필요 | 불필요 |
| **EntryPoint.sol** | **업그레이드** | v0.7→v0.8+ | 7702 호환. userOpHash에 delegation address 포함. | eth-infinitism 감사 완료 |
| BaseAccount.sol | 없음 | 동일 | 코드 변경 불필요 | 불필요 |
| BasePaymaster.sol | 없음 | 동일 | 코드 변경 불필요 | 불필요 |
| ~~SimpleAccount.sol~~ | **제거** | genesis에서 제거 | 7702로 대체. npm 패키지로만 제공. | — |
| ~~SimpleAccountFactory.sol~~ | **제거** | genesis에서 제거 | 7702로 대체. npm 패키지로만 제공. | — |
| **Simple7702Account.sol** | **신규 (대체)** | Factory 대신 도입 | 7702 delegation target. stateless. | eth-infinitism 감사 완료 |
| VerifyingPaymaster.sol | 최소 | 동일 | 주석 정리 | 불필요 |
| **TonTokenPaymaster.sol** | **신규** | 동일 | L2 ETH→TON 가스비 정산 | **필수** |
| **ITokenPriceOracle.sol** | **신규** | 동일 | TON/ETH 환율 오라클 인터페이스 | **필수** |
| **SimplePriceOracle.sol** | **신규** | 동일 | 오라클 기본 구현 | **필수** |
| Bundler config | 설정 | 동일 | TON 가스 가격 기준 | 불필요 |
| tokamak-thanos-geth | 없음 | 신규 확인 | Pectra/Isthmus 업스트림 머지 완료 | 불필요 (업스트림) |

**감사 대상**: TonTokenPaymaster, ITokenPriceOracle, SimplePriceOracle (3개만). EntryPoint v0.8과 Simple7702Account는 eth-infinitism에서 이미 감사 완료.

---

## 11. Bundler 변경 (v1.0에서 추가 변경)

### 11.1 7702 UserOp 지원

Bundler가 7702 delegation이 설정된 EOA의 UserOp을 처리할 수 있어야 한다. 핵심 변경:

- **sender 코드 확인**: UserOp의 sender EOA에 7702 delegation이 설정되어 있는지 `EXTCODESIZE`로 확인
- **delegation address 추출**: sender의 delegation target이 알려진 Simple7702Account인지 검증
- **userOpHash 계산**: v0.8+에서 delegation address를 hash에 포함

```typescript
// bundler.config.ts (v2.0 추가)
export const thanosConfig = {
  entryPoint: "0x4200000000000000000000000000000000000063",
  
  // 허용된 7702 delegation target 목록
  allowed7702Targets: [
    "0x4200000000000000000000000000000000000065", // Simple7702Account predeploy
  ],
  
  // 가스 가격: TON 단위 (v1.0과 동일)
  minGasPrice: parseUnits("0.001", "gwei"),
  minBundleProfit: parseUnits("0.01", "ether"),
  beneficiary: process.env.BUNDLER_ADDRESS,
};
```

### 11.2 ERC-4337 RPC API

7702 관련 추가 RPC 변경은 없다. `eth_sendUserOperation`, `eth_estimateUserOperationGas` 등 기존 API가 동일하게 작동한다. Bundler 내부에서 sender의 7702 상태를 확인하는 것은 validation 로직의 일부이다.

---

## 12. 테스트 시나리오 (수정)

### T1: 7702 Delegation 설정 + 기본 UserOp
1. 유저 EOA에 7702 delegation TX 전송 (target = Simple7702Account predeploy)
2. EOA에 TON deposit (`entryPoint.depositTo`)
3. UserOp 생성 (sender = EOA 주소, callData = ERC-20 transfer)
4. Bundler가 handleOps 호출
5. **검증**: EOA의 TON deposit이 가스비만큼 차감
6. **검증**: ERC-20 transfer가 EOA 주소에서 실행됨

### T2: 7702 + VerifyingPaymaster (gasless)
1. EOA에 7702 delegation 설정
2. VerifyingPaymaster에 TON deposit
3. 유저가 TON 없이 UserOp 생성
4. 오프체인 서버가 paymasterAndData에 서명
5. handleOps → Paymaster의 TON deposit에서 가스비 차감
6. **검증**: 유저 EOA의 TON 잔액 변화 없음

### T3: 7702 + TonTokenPaymaster (L2 ETH 가스비)
1. EOA에 7702 delegation 설정
2. TonTokenPaymaster에 TON deposit + oracle 설정
3. 유저 EOA가 L2 ETH 보유, TON은 0
4. L2 ETH를 Paymaster에 approve
5. UserOp → Paymaster가 TON으로 대납, 유저의 L2 ETH 수거
6. **검증**: 유저의 L2 ETH가 가스비+마크업만큼 차감
7. **검증**: Paymaster의 TON deposit이 차감

### T4: 7702 + Batch TX
1. EOA에 7702 delegation 설정
2. UserOp의 callData = `executeBatch([NFT.mint(), Game.equip(), Stats.update()])`
3. handleOps 실행
4. **검증**: 3개 호출이 1 TX에서 모두 실행
5. **검증**: 가스비가 1 TX 분만 차감 (개별 3 TX 대비 절감)

### T5: Delegation 해제 후 일반 EOA로 복귀
1. 7702 delegation이 설정된 EOA에서 delegation 해제 TX 전송
2. 이후 UserOp 제출
3. **검증**: EntryPoint에서 validation 실패 (더 이상 Smart Account 아님)
4. **검증**: 일반 EOA TX는 정상 작동

### T6: 기존 EOA 자산 유지 확인
1. 유저 EOA에 TON, L2 ETH, NFT 보유
2. 7702 delegation 설정
3. Batch TX로 NFT transfer 실행
4. **검증**: 모든 자산이 동일한 EOA 주소에서 접근 가능
5. **검증**: NFT가 EOA 주소에서 전송됨 (새 주소가 아님)

### T7: 호환성 — 기존 OP Stack E2E 테스트 통과
1. Gaming Preset genesis로 체인 구동
2. 기존 OP Stack E2E 테스트(deposit/withdraw/fault proof) 실행
3. **검증**: AA predeploy 추가가 기존 기능에 영향 없음

### T8: EXTCODESIZE 호환성
1. 7702 delegation이 설정된 EOA의 EXTCODESIZE 확인
2. **검증**: EXTCODESIZE > 0 반환
3. 기존 predeploy(USDC Bridge, Uniswap V3)에서 해당 EOA와 상호작용
4. **검증**: EXTCODESIZE 체크에 의한 오작동 없음

---

## 13. 구현 일정 (수정)

| 단계 | 기간 | 작업 | 전제조건 |
|------|------|------|----------|
| 13-1 | 1주 | tokamak-thanos-geth의 Pectra/Isthmus 업그레이드 상태 확인. EIP-7702 TX type이 L2에서 정상 작동하는지 검증. | — |
| 13-2 | 1주 | EntryPoint v0.8+ 포팅. Simple7702Account predeploy 추가. VerifyingPaymaster 주석 정리. genesis 생성 로직에 AA predeploy 추가. | 13-1 완료 |
| 13-3 | 2주 | TonTokenPaymaster + SimplePriceOracle 개발. 단위 테스트. | 13-2 완료 |
| 13-4 | 1주 | Bundler 설정 변경 + 7702 UserOp validation 로직 추가. SDK 가스 추정 연동. | 13-3 완료 |
| 13-5 | 2주 | E2E 테스트 T1~T8 전수 실행. EXTCODESIZE 호환성 점검. 감사 준비. | 13-4 완료 |
| 13-6 | 병렬 | TonTokenPaymaster + Oracle 외부 감사 | 13-3 완료 시 시작 |

**총 예상: 7주 (감사 병렬 진행)**

v1.0 대비 1주 증가: 7702 TX 작동 검증(13-1)과 EXTCODESIZE 호환성 점검(T8)이 추가됨. 그러나 SimpleAccountFactory 개발이 사라지고 Simple7702Account는 eth-infinitism에서 가져오므로 실질적 개발량은 비슷하다.

---

## 14. 위험 및 완화 (수정)

| 위험 | 심각도 | 완화 방안 |
|------|--------|----------|
| tokamak-thanos-geth Pectra 머지 불완전 | **High** | 13-1에서 7702 TX 전수 검증. 실패 시 v1.0 (ERC-4337 v0.7) 으로 fallback. |
| EXTCODESIZE 호환성 깨짐 | **Medium** | 기존 predeploy 컨트랙트에서 EXTCODESIZE==0 체크 사용 여부 전수 점검. OP Standard predeploy는 Optimism이 이미 검증 완료. Tokamak 자체 predeploy(USDC Bridge, Uniswap V3)만 점검 필요. |
| Simple7702Account delegation 악용 | **Medium** | delegation target을 predeploy 주소(0x4200...0065)로 고정. Bundler의 `allowed7702Targets` whitelist로 알 수 없는 delegation 거부. |
| TonTokenPaymaster pre-charge 실패 | **High** | v1.0과 동일. PostOpMode별 환불 로직 전수 검증. |
| 환율 오라클 조작/지연 | **High** | v1.0과 동일. stale price 체크(24h). Phase 2에서 Uniswap V3 TWAP. |
| eth-infinitism v0.8 미감사 기능 | **Low** | v0.8은 감사 완료. Simple7702Account도 감사 대상에 포함됨. 추가 감사 불필요. |
| Bundler 7702 호환성 | **Low** | 주요 Bundler(Stackup, Pimlico, Alto)가 이미 v0.8+/7702 지원. 설정만 변경하면 사용 가능. |

---

## 15. v1.0 → v2.0 마이그레이션 경로

v1.0으로 이미 배포된 체인이 있는 경우:

1. **EntryPoint 업그레이드**: Transparent Proxy이므로 admin이 implementation을 v0.8+로 변경
2. **Simple7702Account predeploy 추가**: genesis 업데이트 없이 일반 배포로 가능 (predeploy 주소가 아닌 일반 주소)
3. **기존 SimpleAccount 유저**: 기존 Smart Account가 그대로 작동. v0.8은 v0.7과 ABI 호환.
4. **신규 유저**: 7702 delegation으로 온보딩. 기존 유저와 공존 가능.

즉, v1.0과 v2.0은 **공존 가능**하다. 기존 유저를 강제 마이그레이션할 필요 없음.

---

## 부록 A: 디렉토리 구조 (수정)

```
packages/tokamak/contracts-bedrock/src/tokamak-contracts/AA/
├── EntryPoint.sol                  ← eth-infinitism v0.8+ 포팅 (7702 통합)
├── Simple7702Account.sol           ← eth-infinitism v0.8+ (delegation target)
├── VerifyingPaymaster.sol          ← eth-infinitism v0.8+ 포팅 (주석만 변경)
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

// 제거됨 (v1.0에서 존재, v2.0에서 제거):
// ├── SimpleAccount.sol            ← npm 패키지로만 제공
// └── SimpleAccountFactory.sol     ← npm 패키지로만 제공
```

## 부록 B: Predeploy 주소 할당 (수정)

```
0x4200000000000000000000000000000000000060  VRFPredeploy        (Gaming/Full)
0x4200000000000000000000000000000000000061  VRFCoordinator       (Gaming/Full)
0x4200000000000000000000000000000000000062  VRFConsumerBase      (Gaming/Full)
0x4200000000000000000000000000000000000063  EntryPoint v0.8+     (Gaming/Full)  ← 업그레이드
0x4200000000000000000000000000000000000064  VerifyingPaymaster   (Gaming/Full)  ← 동일
0x4200000000000000000000000000000000000065  Simple7702Account    (Gaming/Full)  ← 변경 (Factory→7702)
```

## 부록 C: v1.0 대비 v2.0 장점 요약

| 측면 | v1.0 (ERC-4337 v0.7) | v2.0 (EIP-7702 + v0.8) | 개선 |
|------|---------------------|------------------------|------|
| 유저 주소 | 새 Smart Account 주소 | 기존 EOA 주소 유지 | 혼란 제거 |
| 온보딩 | Factory 배포 (~300K gas) | 7702 delegation (~50K gas) | 가스비 83% 절감 |
| 자산 이전 | 기존 EOA → Smart Account 이전 필요 | 불필요 | UX 대폭 개선 |
| 감사 범위 | EntryPoint v0.7 + SimpleAccount + Factory | TonTokenPaymaster + Oracle (3개만) | 감사 비용 절감 |
| 생태계 호환 | ERC-4337 지갑만 지원 | 모든 EOA 지갑 지원 (MetaMask, Rabby 등) | 생태계 접근성 |
| 업스트림 추적 | v0.7은 레거시화 진행 중 | v0.8+은 현재 메인 브랜치 | 장기 유지보수 유리 |
| OP Stack 호환 | Level 1 변경 (predeploy만) | Level 0 (업스트림 머지) + Level 1 | 더 안전 |

---

*End of Document*
