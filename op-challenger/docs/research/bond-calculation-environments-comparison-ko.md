# Fault Proof 보증금 계산 가이드 (Optimism Mainnet 기준)

> **작성일**: 2025-11-04
> **기준**: Optimism Mainnet (`mainnet.json`) 공식 설정
> **목적**: 환경설정이 보증금에 미치는 영향을 정확히 계산하고 이해

## 📋 목차

1. [Optimism Mainnet 표준 설정](#1-optimism-mainnet-표준-설정)
2. [보증금 계산 공식](#2-보증금-계산-공식)
3. [제출 간격과 게임 범위](#3-제출-간격과-게임-범위)
4. [Depth별 보증금 금액](#4-depth별-보증금-금액)
5. [필요 자본 계산](#5-필요-자본-계산)
6. [환경별 비교 (참고)](#6-환경별-비교-참고)

---

## 1. Optimism Mainnet 표준 설정

### 1.1 핵심 설정값 (`mainnet.json`)

```json
{
  // Fault Game 설정
  "faultGameMaxDepth": 73,                    // 최대 게임 깊이
  "faultGameSplitDepth": 30,                  // Output Root/Trace 분리점
  "faultGameMaxClockDuration": 302400,        // 3.5일 (각 팀)
  "faultGameClockExtension": 10800,           // 3시간
  "faultGameWithdrawalDelay": 604800,         // 7일

  // Proposer 설정
  "l2OutputOracleSubmissionInterval": 1800,   // 1800 블록마다 제출

  // 블록 설정
  "l2BlockTime": 2,                           // 2초/블록
  "l1BlockTime": 12,                          // 12초/블록

  // GameType
  "respectedGameType": 0                      // Cannon (MIPS VM)
}
```

### 1.2 설정의 의미

| 설정 | 값 | 실제 의미 | 영향 |
|------|-----|----------|------|
| **l2OutputOracleSubmissionInterval** | 1800 blocks | 60분마다 제출 | **게임 1개가 1800 블록 커버** |
| **faultGameMaxDepth** | 73 | 최대 깊이 | 2^73 instructions까지 검증 가능 |
| **faultGameSplitDepth** | 30 | 분리점 | Depth 0-30: 블록 단위, 31-73: instruction 단위 |
| **faultGameMaxClockDuration** | 302400초 | 3.5일 | 각 팀의 최대 응답 시간 |

### 1.3 게임 범위 계산

```
제출 간격: 1800 blocks
L2 Block Time: 2초
    ↓
실제 시간: 1800 × 2초 = 3600초 = 60분
    ↓
게임 1개 = 60분 동안의 L2 블록을 검증
         = 1800개 블록을 커버
```

---

## 2. 보증금 계산 공식

### 2.1 온체인 계산 로직

보증금은 **FaultDisputeGame.sol**에서 다음과 같이 계산됩니다:

```solidity
// packages/tokamak/contracts-bedrock/src/dispute/FaultDisputeGame.sol:742-783

function getRequiredBond(Position _position) public view returns (uint256 requiredBond_) {
    uint256 depth = uint256(_position.depth());

    // 고정 상수 (Big Bonds v1.5 스펙)
    uint256 assumedBaseFee = 200 gwei;      // ← 고정값!
    uint256 baseGasCharged = 400_000;       // 초기 가스량
    uint256 highGasCharged = 300_000_000;   // 최대 가스량

    // 지수 승수 계산
    // multiplier = (highGasCharged/baseGasCharged)^(1/MAX_GAME_DEPTH)
    uint256 a = highGasCharged / baseGasCharged;  // 750
    uint256 multiplier = a ** (1/MAX_GAME_DEPTH);

    // 필요 가스 계산
    uint256 requiredGas = baseGasCharged × (multiplier^depth);

    // 최종 보증금
    requiredBond_ = assumedBaseFee × requiredGas;
}
```

### 2.2 핵심 상수 (Optimism Mainnet 기준)

| 상수 | 값 | 의미 | 변경 가능? |
|------|-----|------|-----------|
| **assumedBaseFee** | 200 gwei | 가정된 L1 base fee | ❌ 하드코딩 |
| **baseGasCharged** | 400,000 gas | 초기 가스량 (Depth 0) | ❌ 하드코딩 |
| **highGasCharged** | 300,000,000 gas | 최대 가스량 (Depth MAX) | ❌ 하드코딩 |
| **MAX_GAME_DEPTH** | **73** (Optimism Mainnet 기준) | 최대 게임 깊이 | ✅ 환경별 설정 |

**환경별 MAX_GAME_DEPTH**:
- Optimism Mainnet: **73** (표준)
- Tokamak Sepolia: **73** (표준 준수)
- Tokamak Devnet: **50** (테스트용, 빠른 실행)

**⚠️ 중요**:
- `assumedBaseFee`는 실제 L1 가스비와 무관하게 **200 gwei로 고정**
- `MAX_GAME_DEPTH`는 **Optimism Mainnet 기준 73**이 표준

### 2.3 실제 L1 가스비와의 관계

```
현재 L1 실제 가스비 (2025년 11월 기준):
- 평소: 1-5 gwei (Dencun 업그레이드 후 극적으로 감소)
- 혼잡: 10-20 gwei
- 극한: 50-100 gwei

보증금 계산에 사용되는 가스비:
- 항상: 200 gwei (고정)

→ 보증금은 L1 가스비 변동에 영향받지 않음!
→ 예측 가능한 고정 비용
→ 실제 가스비가 낮아도 보증금은 동일 (200 gwei 기준)

```

**💡 핵심 인사이트**:
- 실제 L1 가스비는 **1-20 gwei** (평소)
- 보증금 계산은 **200 gwei 고정**
- **보증금이 실제 가스비보다 10-200배 높게 책정**
- 이는 의도된 설계:
  - ✅ 보안을 위한 보수적 설정
  - ✅ 가스비 급등 시에도 충분한 보증금 유지
  - ✅ 예측 가능성 (시장 변동 무관)

---

## 3. 제출 간격과 게임 범위

### 3.1 Optimism Mainnet의 게임 범위

**설정**:
```
l2OutputOracleSubmissionInterval: 1800 blocks
l2BlockTime: 2초
```

**계산**:
```
게임 1개가 커버하는 범위:
= 1800 blocks
= 1800 × 2초 = 3600초 = 60분

결론: 게임 하나당 60분 분량의 L2 블록을 검증
```

### 3.2 게임 범위와 Depth의 관계

**핵심**: 게임 범위는 고정(1800 블록), Depth는 **그 범위를 얼마나 세분화하는가**

```
게임 범위: 1800 블록 (Block #1000 ~ #2800)

Depth 0 (Root):
└─ 전체 1800 블록 커버

Depth 1:
├─ Left:  900 블록 (Block #1000 ~ #1900)
└─ Right: 900 블록 (Block #1900 ~ #2800)

Depth 2:
├─ 450 블록
├─ 450 블록
├─ 450 블록
└─ 450 블록

Depth 10:
└─ 각 claim이 1800/2^10 = 1.76 블록씩 커버

Depth 20:
└─ 각 claim이 1800/2^20 = 0.0017 블록 (블록보다 작음!)

Depth 30 (SPLIT_DEPTH): ⭐ 여기서 전환!
└─ 1800/2^30 = 0.0000017 블록
    → 블록 단위에서 Instruction 단위로 전환

Depth 31-73:
└─ VM Instruction 단위로 Binary Search
    └─ Depth 73: 단일 instruction
```

### 3.3 왜 Instruction까지 가야 하는가? ⭐⭐⭐

#### 핵심 질문

> "1800 블록 특정까지는 Depth 11이면 충분한데,
> 왜 Depth 73까지 가야 하나?
> 블록만 특정하면 되지 않나?"

#### 답변: L1에서 직접 검증하기 위해서

**블록 레벨에서 멈추면 검증 불가능**:

```
만약 Depth 11에서 멈춘다면:

Challenger: "Block #1523이 잘못됐어!"
L1 Contract: "증명해봐"
Challenger: "Block #1523 전체를 L1에서 재실행하면 돼"

문제:
├─ Block #1523 실행 = 30M gas 필요
├─ L1 transaction limit = ~30M gas
├─ 검증 로직 자체도 gas 필요
└─ 불가능! ❌

결과:
└─ L1에서 검증 못 함
    └─ "믿어줘"에 의존
        └─ Fault Proof 아님! ❌
```

**Instruction까지 가면 검증 가능**:

```
Depth 42-43 (단일 Instruction):

Challenger: "Block #1523의 Instruction #12,345,678이 잘못됐어!"
           "Prestate: 0xABCD,
            예상 Poststate: 0x1111,
            Proposer 주장: 0x2222"

L1 Contract: "검증하자"
    ↓
step() 함수 호출:
    ├─ 온체인 VM (IBigStepper) 로드
    ├─ 단일 Instruction 실행 (MIPS ADD 등)
    ├─ Gas 비용: ~500k (감당 가능 ✅)
    ├─ 실제 Poststate 계산: 0x1111
    └─ Proposer 주장(0x2222)과 비교
        └─ 불일치! Proposer가 틀렸음 증명 ✅

결과:
└─ L1에서 직접 검증 완료
    └─ 수학적 증명
        └─ 진짜 Fault Proof! ✅
```

#### 계산: 1800 블록 → 실제 필요 Depth

```
1단계: 블록 특정
└─ log₂(1800) ≈ 11 depth
   └─ "Block #1523이 문제" 특정

2단계: Instruction 특정 (핵심!)
└─ 1800 블록의 총 Instructions (최악):
    ├─ 1800 blocks × 30M gas/block = 54B gas
    ├─ 54B gas × 100 instructions/gas = 5.4T instructions
    └─ log₂(5.4T) ≈ 42 depth
        └─ "Instruction #12,345,678이 문제" 특정

3단계: 안전 마진
└─ 실제 필요: 42 depth
    ├─ 설정값: 73 depth
    └─ 여유: 31 depth (미래 확장, 복잡한 블록 대비)
```

#### 실제 도달 Depth 분포

| 블록 복잡도 | 1800 블록의 총 Instructions | 필요 Depth | 실제 발생 확률 |
|-----------|---------------------------|-----------|-------------|
| **간단** (전송 위주) | ~1.8B | ~31 | 80% ⭐ |
| **일반** (컨트랙트 호출) | ~180B | ~38 | 15% |
| **복잡** (DEX, NFT 등) | ~1.8T | ~41 | 4% |
| **최악** (모든 블록 최대) | ~5.4T | **~42** | 1% |
| **이론적 최대** | 2^73 | 73 | 0% (불가능) |

**결론**:
- 대부분의 게임: Depth 31-38에서 종료 (99%)
- 복잡한 경우: Depth 38-42에서 종료 (0.9%)
- Depth 73: 이론적 상한선 (실제로는 도달 안 함)

#### 온체인 검증: step() 함수

**FaultDisputeGame.sol의 핵심**:

```solidity
function step(
    uint256 _claimIndex,
    bool _isAttack,
    bytes calldata _stateData,   // Prestate (instruction 실행 전 VM 상태)
    bytes calldata _proof         // Merkle proof
) public {
    // 1. MAX_GAME_DEPTH인지 확인 (단일 instruction만)
    require(position.depth() == MAX_GAME_DEPTH, "Not at max depth");

    // 2. 온체인 VM으로 단일 instruction 실행
    bytes32 postState = VM.step(_stateData, _proof);
    //                   ↑ IBigStepper (Cannon/Asterisc VM)
    //                   ↑ Gas 비용: ~500k

    // 3. 결과 비교
    if (postState == claim.value) {
        // Claim이 맞음 → 반박 실패
        counter claim loses
    } else {
        // Claim이 틀림 → 반박 성공
        claim is countered ✅
    }
}
```

**왜 이것이 가능한가?**:
- ✅ 단일 Instruction 실행 = ~500k gas (L1에서 가능)
- ✅ 결과가 확정적 (deterministic)
- ✅ 누구나 검증 가능 (trustless)

#### 핵심 정리

```
블록 특정 (Depth 11):
❌ L1에서 검증 불가능 (30M gas = 블록 전체 실행 필요)
❌ "믿어줘" 시스템
❌ Fault Proof 아님

Instruction 특정 (Depth 42):
✅ L1에서 검증 가능 (500k gas = 단일 instruction 실행)
✅ 온체인 수학적 증명 (step() 함수)
✅ 진짜 Fault Proof!

MAX_GAME_DEPTH = 73:
└─ 실제 필요(42) + 안전 마진(31)
    └─ "갈 수 있다"는 상한선
    └─ 미래의 더 복잡한 블록 대비
```

**💡 Fault Proof의 본질**:
> "L1에서 직접 실행해서 검증할 수 있을 만큼 작은 단위(단일 instruction)까지 좁히는 것"
> 이것이 없으면 그냥 "믿어줘" 시스템이 됩니다.

### 3.4 환경별 게임 범위 비교

| 환경 | 제출 간격 (blocks) | 게임 범위 (시간) | 게임당 블록 수 |
|------|--------------------|-----------------|---------------|
| **Mainnet** | 1800 | 60분 | **1800 블록** |
| **Sepolia** | 120 | 24분 | **120 블록** |
| **Devnet** | 10 | 20초 | **10 블록** |

**결론**:
- 제출 간격이 길수록 → 게임당 더 많은 블록 커버
- 더 많은 블록 → 더 많은 Instructions → 더 깊은 Depth 필요 → 더 많은 보증금

---

## 4. Depth별 보증금 금액

### 6.1 전체 설정 비교표

| 설정 항목 | Devnet (Tokamak) | Sepolia (Tokamak) | Mainnet (Optimism) |
|----------|-----------------|-------------------|-------------------|
| **L1 Chain** | Local (Anvil) | Sepolia | Ethereum Mainnet |
| **L2 Chain** | Local | Thanos Sepolia | OP Mainnet |
| **L1 Block Time** | 3초 | 12초 | 12초 |
| **L2 Block Time** | 2초 | 12초 | 2초 |
| **l2OutputOracleSubmissionInterval** | 10 blocks | 120 blocks | 1800 blocks |
| **실제 제출 주기 (시간)** | **20초** | **24분** | **60분** |
| **faultGameMaxDepth** | **50** | **73** | **73** |
| **faultGameSplitDepth** | **14** | **30** | **30** |
| **faultGameMaxClockDuration** | 150초 (2.5분) | 302400초 (3.5일) | 302400초 (3.5일) |
| **faultGameClockExtension** | 0초 | 10800초 (3시간) | 10800초 (3시간) |
| **faultGameWithdrawalDelay** | 3600초 (1시간) | 604800초 (7일) | 604800초 (7일) |
| **respectedGameType** | 3 (AsteriscKona) | 0 (Cannon) | 0 (Cannon) |
| **faultGameAbsolutePrestate** | 0x03ce7bc8... | 0x03ab262c... | 0x037ef3c1... |
| **proofMaturityDelaySeconds** | 12초 | 604800초 (7일) | 604800초 (7일) |
| **disputeGameFinalityDelaySeconds** | 6초 | 302400초 (3.5일) | 302400초 (3.5일) |

**출처**:
- Devnet: `packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json`
- Sepolia: `packages/tokamak/contracts-bedrock/deploy-config/thanos-sepolia.json`
- Mainnet: `packages/contracts-bedrock/deploy-config/mainnet.json` (Optimism 공식)

### 6.2 파라미터 간 관계 분석

#### 제출 간격 → MAX_GAME_DEPTH 관계

```
제출 간격이 결정하는 것:
└─> 게임 범위 (몇 개의 L2 블록을 검증)
    └─> 필요한 Depth = log₂(총 instruction 수)
        └─> MAX_GAME_DEPTH >= 필요 Depth 이어야 함

예시 검증:
┌──────────┬───────────┬──────────────┬─────────────┬──────────────┬─────────┐
│   환경    │ 제출 간격  │  게임 범위   │ 예상 Instr  │ 필요 Depth  │ 설정값  │
├──────────┼───────────┼──────────────┼─────────────┼──────────────┼─────────┤
│ Devnet   │ 20초      │ 10 blocks    │ ~100M       │ ~27         │ 50 ✅   │
│ Sepolia  │ 24분      │ 120 blocks   │ ~1.2B       │ ~31         │ 73 ✅   │
│ Mainnet  │ 60분      │ 1800 blocks  │ ~18B        │ ~35         │ 73 ✅   │
└──────────┴───────────┴──────────────┴─────────────┴──────────────┴─────────┘

계산:
- 가정: 블록당 평균 10M instructions
- 필요 Depth: log₂(instruction 수)
- MAX_DEPTH는 필요 Depth보다 충분히 커야 함
```

**핵심**: 제출 간격이 길수록 → 게임 범위 증가 → 더 높은 MAX_DEPTH 필요

#### "검사 보장 시간"의 진실 ⭐

**❌ 오해**: "MAX_GAME_DEPTH + 제출간격 = 검사 시간"
**✅ 실제**: 검사 보장 시간은 **faultGameMaxClockDuration + Challenge Period**

```
실제 검사 보장 시간:

┌─ Challenge Period (7일, Sepolia/Mainnet)
│  └─ Output Root 제출 후 이의제기 가능 기간
│
├─ Dispute Game Duration (최대 3.5일 × 2 = 7일)
│  ├─ Proposer Clock: 3.5일 (faultGameMaxClockDuration)
│  └─ Challenger Clock: 3.5일 (faultGameMaxClockDuration)
│
└─ Total: 7일 + 7일 = 최대 14일

Devnet (테스트):
└─ Total: ~5분 (빠른 검증용)
```

| 환경 | Challenge Period | Game Duration | **총 보장 시간** |
|------|-----------------|---------------|----------------|
| **Devnet** | ~0 (테스트) | 2.5분 | **~5분** |
| **Sepolia** | 7일 | 최대 7일 | **최대 14일** |
| **Mainnet** | 7일 | 최대 7일 | **최대 14일** |

**결론**:
- MAX_GAME_DEPTH는 "얼마나 깊이" 검증할 수 있는지 (instruction 수)
- faultGameMaxClockDuration은 "얼마나 오래" 검증할 수 있는지 (시간)
- 둘은 **서로 다른 축**의 보안 파라미터!

### 6.3 주요 차이점 분석

#### 🔵 Devnet (로컬 개발)

```yaml
목적: 빠른 테스트 및 개발
특징:
  - 매우 짧은 제출 주기 (20초)
  - 빠른 게임 진행 (2.5분 제한)
  - 낮은 Max Depth (50)
  - 즉시 출금 가능 (1시간)

용도:
  - 기능 테스트
  - 통합 테스트
  - 프로토타입 개발
```

#### 🟢 Sepolia Testnet (공개 테스트넷)

```yaml
목적: 실제 환경에 가까운 테스트
특징:
  - 중간 제출 주기 (24분)
  - 현실적 게임 진행 (3.5일)
  - 높은 Max Depth (73)
  - 안전한 출금 대기 (7일)

용도:
  - 보안 테스트
  - 부하 테스트
  - 실전 시뮬레이션
```

#### 🟡 Mainnet (Optimism 공식 설정)

```yaml
목적: 프로덕션 환경 (Optimism Mainnet 기준)
특징:
  - 긴 제출 주기 (60분, 1800 blocks)
  - 안전한 게임 진행 (3.5일)
  - 높은 Max Depth (73)
  - 보수적 출금 대기 (7일)
  - Cannon VM 사용 (GameType 0)

용도:
  - 실제 자금 보호
  - 탈중앙화 검증

출처: packages/contracts-bedrock/deploy-config/mainnet.json
```

---

## 4. Depth별 보증금 금액

### 4.1 Optimism Mainnet 표준 (MAX_GAME_DEPTH = 73) ⭐

**Optimism Mainnet 공식 설정** (`mainnet.json`):

```
보증금 계산 (Optimism 표준):
- assumedBaseFee: 200 gwei (하드코딩, 고정)
- baseGasCharged: 400,000 gas (하드코딩, 고정)
- highGasCharged: 300,000,000 gas (하드코딩, 고정)
- MAX_GAME_DEPTH: 73
- multiplier: (750)^(1/73) ≈ 1.0935
```

| Depth | 보증금 (ETH) | 보증금 (Gwei) | 커버 범위 | 구간 | 비고 |
|-------|-------------|---------------|----------|------|------|
| **0** | 0.08 | 80,000,000 | 1800 블록 전체 | Root | 초기 진입 |
| **10** | 0.65 | 646,928,486 | 1.76 블록 | Output Root | |
| **20** | 3.16 | 3,157,978,182 | 0.0017 블록 | Output Root | |
| **30** | 15.42 | 15,415,844,048 | 0.0000017 블록 | **Split Depth** ⭐ | 블록→Instruction 전환 |
| **31** | 16.86 | 16,856,754,530 | Instructions/2 | Trace 시작 | |
| **35** | 24.90 | 24,897,459,123 | Instructions/32 | Trace | 일반 분쟁 예상 |
| **40** | 50.91 | 50,910,264,694 | Instructions/1024 | Trace | |
| **50** | 240.60 | 240,603,842,312 | Instructions/1M | Trace 깊음 | |
| **60** | 1,138 | 1,138,183,836,054 | Instructions/1B | 매우 깊음 | 거의 도달 안 함 |
| **73** | 10,707 | 10,707,396,688,564 | 1 instruction | **최대** | 사실상 불가능 |

### 4.2 주요 Depth 분석 (Mainnet 기준)

#### Depth 30 (SPLIT_DEPTH) - 가장 중요! ⭐

```
의미: Output Root 검증 종료, Trace 시작
보증금: 15.42 ETH
커버: 0.0000017 블록 (블록보다 작아서 전환)

이 지점까지:
- 1800 블록을 2^30 = 10억 조각으로 나눔
- 블록 단위로는 더 이상 의미 없음
- Instruction 단위로 전환
```

#### Depth 35 (일반 분쟁 예상)

```
보증금: 24.90 ETH
커버: Instructions/32

실제 의미:
- 대부분의 분쟁이 여기서 해결
- 1800 블록 중 정확한 instruction 범위 특정
```

#### Depth 42 (최악의 경우)

```
보증금: 126.41 ETH
커버: Instructions/4096

실제 의미:
- 매우 복잡한 블록 (1800개)
- 최악의 경우 이 정도까지 도달
- 누적 보증금: ~300 ETH (한쪽)
```

---

## 5. 필요 자본 계산

### 4.1 제출 주기와 게임 범위

#### VM Instruction이란? ⭐

**핵심 개념**: Fault Proof에서 검증하는 것은 **VM instructions**입니다.

```
L2 블록 실행 과정:

1. L2 블록 = 여러 트랜잭션 묶음
   ├─ Transaction 1: ETH 전송
   ├─ Transaction 2: 스마트 컨트랙트 호출
   └─ Transaction 3: DEX swap

2. op-program이 블록을 재실행
   ├─ MIPS VM (Cannon) 또는
   └─ RISC-V VM (Asterisc)

3. VM이 실행하는 각 명령어 = VM Instruction
   ├─ MIPS: ADD, SUB, LW, SW, BRANCH 등
   └─ RISC-V: ADDI, LOAD, STORE, JAL 등

4. 한 블록의 총 VM Instruction 수
   ├─ 간단한 블록: ~1M instructions
   ├─ 일반 블록: ~10M instructions
   └─ 복잡한 블록: ~100M instructions
```

**중요**:
- EVM opcode ≠ VM instruction
- EVM의 하나의 opcode가 수백~수천 개의 VM instruction으로 변환됨

**예시**:
```
EVM SSTORE (스토리지 쓰기):
  ↓ op-program이 MIPS/RISC-V로 변환
  ↓
~1000-10000 MIPS/RISC-V instructions
  (메모리 할당, 해시 계산, 상태 업데이트 등)
```

#### 블록 Gas Limit vs VM Instructions

```
L2 Block Gas Limit: 30,000,000 gas (0x1c9c380)

하지만 이것은 EVM gas이지 VM instruction이 아님!

실제 변환:
┌─────────────────────────────────────────────────────┐
│ EVM Gas → VM Instructions 변환 비율                  │
├─────────────────────────────────────────────────────┤
│ 간단한 연산 (ADD, MUL):                              │
│   - EVM: 3-5 gas                                     │
│   - VM: ~100-1000 instructions                       │
│   - 비율: 1 gas → ~200-300 instructions              │
│                                                       │
│ 중간 연산 (SLOAD):                                   │
│   - EVM: 200-2100 gas                                │
│   - VM: ~5000-50000 instructions                     │
│   - 비율: 1 gas → ~25-100 instructions               │
│                                                       │
│ 무거운 연산 (SSTORE):                                │
│   - EVM: 20000 gas                                   │
│   - VM: ~100000-1000000 instructions                 │
│   - 비율: 1 gas → ~5-50 instructions                 │
└─────────────────────────────────────────────────────┘

평균 비율 (추정): 1 gas ≈ 10-100 VM instructions
```

#### 실제 계산

```
Devnet 예시:
- 10 블록 × 30M gas/블록 = 300M gas
- 변환 비율: 1 gas ≈ 10-100 instructions (평균)
- 예상 VM instructions: 3B - 30B
- 필요 Depth: log₂(30B) ≈ 35

하지만 실제로는:
- 블록이 비어있을 수도 있음 (거의 0 instructions)
- DEX 거래가 많으면 복잡 (수억 instructions)
- 평균을 내기 어려움

→ 문서에서는 보수적으로 추정
```

#### 제출 주기와의 관계

```
Proposer 제출 주기 = L2 Output 제출 간격
→ 게임 범위 = 제출 간격 동안의 L2 블록 수
→ Trace 길이 = 게임 범위의 총 VM instruction 수 (변동 큼!)
→ 필요 Depth = log₂(Trace 길이)
→ MAX_GAME_DEPTH ≥ 필요 Depth (충분한 여유 필요)
```

### 4.2 환경별 게임 범위

| 환경 | 제출 주기 | L2 Block Time | 게임 범위 (블록) | VM Instruction 수* | 필요 Depth** |
|------|----------|---------------|------------------|--------------------|-------------|
| **Devnet** | 20초 | 2초 | 10 blocks | ~100M - 1B | ~27 - 30 |
| **Sepolia** | 24분 | 12초 | 120 blocks | ~1B - 10B | ~30 - 34 |
| **Mainnet** | 60분 | 2초 | 1800 blocks | ~10B - 100B | ~34 - 37 |

*VM Instruction 수: **실제 trace 길이는 블록 내용에 따라 크게 변동**
- 간단한 블록 (전송만): ~1M instructions/block
- 복잡한 블록 (컨트랙트 실행): ~100M instructions/block
- 평균 추정치는 실제 네트워크 활동에 따라 다름

**필요 Depth: log₂(총 VM instruction 수)

**⚠️ 문서 전반의 "10M instructions/block"은 개략적 추정치**:
- 실제로는 1M - 100M instructions/block 사이에서 크게 변동
- 네트워크 활동, 트랜잭션 복잡도에 따라 다름
- MAX_GAME_DEPTH는 최악의 경우를 대비한 충분한 여유를 제공

### 4.3 실제 영향 분석

#### Case 1: Devnet (20초 주기)

```
게임 범위: 10 L2 blocks
├─ Block 1000 ~ Block 1010
├─ Gas Limit per Block: 30M gas (0x1c9c380)
├─ Total Gas: ~300M gas
└─ 예상 VM Instructions: ~100M - 1B (블록 내용에 따라 변동)

실제 필요 Depth (보수적 추정):
├─ Output Root: 0 → 14 (2.92 ETH)
├─ Trace: 15 → 30 (96 ETH)
└─ 총 예상 비용: ~18 - 100 ETH (한쪽, 복잡도에 따라)

⚠️ 중요:
- 실제 trace 길이는 블록 내 트랜잭션 복잡도에 따라 크게 변동
- Depth 50은 최대 2^50 instructions까지 커버 (충분한 여유)
```

#### Case 2: Sepolia (24분 주기)

```
게임 범위: 120 L2 blocks
├─ Block 1000 ~ Block 1120
├─ Gas Limit per Block: 30M gas
├─ Total Gas: ~3.6B gas
└─ 예상 VM Instructions: ~1B - 10B (블록 내용에 따라 변동)

실제 필요 Depth (보수적 추정):
├─ Output Root: 0 → 30 (15.42 ETH)
├─ Trace: 31 → 34 (32.24 ETH)
└─ 총 예상 비용: ~32 - 47 ETH (한쪽, 복잡도에 따라)

⚠️ 중요:
- MAX_DEPTH 73은 2^73 instructions까지 커버
- 실제로는 Depth 34 정도면 대부분의 블록 커버
- 충분한 안전 마진 확보
```

#### Case 3: Mainnet (60분 주기, Optimism 기준)

```
게임 범위: 1800 L2 blocks (Optimism Mainnet 설정)
├─ Block 1000 ~ Block 2800
├─ Gas Limit per Block: 30M gas (Optimism)
├─ Total Gas: ~54B gas
└─ 예상 VM Instructions: ~10B - 100B (블록 내용에 따라 변동)

실제 필요 Depth (보수적 추정):
├─ Output Root: 0 → 30 (15.42 ETH)
├─ Trace: 31 → 37 (56.17 ETH)
└─ 총 예상 비용: ~56 - 71 ETH (한쪽, 복잡도에 따라)

⚠️ 중요:
- 1800 블록은 매우 넓은 범위
- 실제 trace 길이는 네트워크 활동에 따라 크게 변동
- MAX_DEPTH 73은 충분한 여유 제공

출처: Optimism mainnet.json (l2GenesisBlockGasLimit: 30M)
```

### 4.4 제출 주기 최적화

#### 🔵 짧은 주기 (< 1분)

```
장점:
✅ 빠른 finality
✅ 작은 게임 범위 → 낮은 depth 필요
✅ 낮은 보증금 요구
✅ 빠른 분쟁 해결

단점:
❌ Proposer 부담 증가 (빈번한 트랜잭션)
❌ L1 가스비 누적
❌ 네트워크 부하 증가
```

#### 🟢 중간 주기 (10-30분)

```
장점:
✅ 균형잡힌 finality
✅ 적절한 게임 범위
✅ 합리적 보증금 요구
✅ Proposer 부담 적정

단점:
⚠️ 일부 trade-off 필요
```

#### 🟡 긴 주기 (> 1시간)

```
장점:
✅ Proposer 부담 최소화
✅ L1 가스비 절감
✅ 네트워크 효율성

단점:
❌ 느린 finality (출금 시간 증가)
❌ 큰 게임 범위 → 높은 depth 필요
❌ 높은 보증금 요구
❌ 분쟁 해결 시간 증가
```

---

## 5. 실제 비용 시나리오

### 5.1 정직한 분쟁 시나리오

**가정**: 악의적 Proposer가 잘못된 Output 제출, 정직한 Challenger가 반박

#### Devnet에서의 비용

```
악의적 Proposer:
├─ Depth 0-27 참여
├─ 짝수 depth 담당
├─ 예상 보증금: ~18 ETH
├─ 게임 패배
└─ 손실: -18 ETH 💸

정직한 Challenger:
├─ Depth 1-27 참여
├─ 홀수 depth 담당
├─ 예상 보증금: ~18 ETH 예치
├─ 게임 승리
├─ 자신의 보증금 회수: +18 ETH
├─ 상대방 보증금 획득: +18 ETH
└─ 순이익: +18 ETH 💰

실질 비용:
- Challenger: 0 ETH (오히려 +18 ETH)
- 시스템: 안전하게 보호됨 ✅
```

#### Sepolia에서의 비용

```
악의적 Proposer:
├─ Depth 0-31 참여
├─ 예상 보증금: ~32 ETH
└─ 손실: -32 ETH 💸

정직한 Challenger:
├─ Depth 1-31 참여
├─ 예상 보증금: ~32 ETH 예치
└─ 순이익: +32 ETH 💰

실질 비용:
- Challenger: 0 ETH (오히려 +32 ETH)
```

#### Mainnet에서의 비용 (Optimism 기준)

```
악의적 Proposer:
├─ Depth 0-35 참여 (Optimism mainnet 설정)
├─ 예상 보증금: ~56 ETH
└─ 손실: -56 ETH 💸

정직한 Challenger:
├─ Depth 1-35 참여
├─ 예상 보증금: ~56 ETH 예치
└─ 순이익: +56 ETH 💰

실질 비용:
- Challenger: 0 ETH (오히려 +56 ETH)

출처: Optimism mainnet.json (l2OutputOracleSubmissionInterval: 1800)
```

### 5.2 필요 자본 계산 (설정 기반)

#### 설정값이 자본을 결정하는 과정 ⭐

Challenger의 필요 자본은 다음 설정들로부터 **논리적으로 계산**됩니다:

```
[설정 1] MAX_GAME_DEPTH
    ↓
최대 검증 가능 instruction 수 결정 (2^MAX_GAME_DEPTH)
    ↓
[설정 2] SPLIT_DEPTH
    ↓
Output Root 검증 범위 결정 (0 ~ SPLIT_DEPTH)
    ↓
[설정 3] l2OutputOracleSubmissionInterval
    ↓
게임당 검증해야 할 블록 수 결정
    ↓
[계산] 실제 도달 가능한 Depth 추정
    ↓
[계산] 해당 Depth까지 누적 보증금 계산
    ↓
필요 자본 도출
```

#### Devnet 자본 계산 과정

**설정값**:
```json
{
  "faultGameMaxDepth": 50,
  "faultGameSplitDepth": 14,
  "l2OutputOracleSubmissionInterval": 10,
  "l2BlockTime": 2
}
```

**계산**:
```
1. 게임 범위
   = 10 blocks × 2초 = 20초 분량

2. 예상 VM Instructions (최악의 경우)
   = 10 blocks × 30M gas/block × 100 instructions/gas
   = 30B instructions

3. 필요 Depth
   = log₂(30B) ≈ 35

4. 하지만 MAX_GAME_DEPTH = 50이므로
   → Depth 35까지만 실제 도달 가능

5. Depth 35까지 누적 보증금 (양측 합계, 최악)
   = Depth 0-35의 보증금 합
   ≈ 100 ETH (한쪽 50 ETH)

6. 최악의 경우 대비 자본
   = 50 ETH (한쪽 부담)
```

**결론**: Devnet에서 **최악의 경우 대비 자본 = 50 ETH**

#### Sepolia 자본 계산 과정

**설정값**:
```json
{
  "faultGameMaxDepth": 73,
  "faultGameSplitDepth": 30,
  "l2OutputOracleSubmissionInterval": 120,
  "l2BlockTime": 12
}
```

**계산**:
```
1. 게임 범위
   = 120 blocks × 12초 = 24분 분량

2. 예상 VM Instructions (최악의 경우)
   = 120 blocks × 30M gas/block × 100 instructions/gas
   = 360B instructions

3. 필요 Depth
   = log₂(360B) ≈ 39

4. MAX_GAME_DEPTH = 73이므로
   → Depth 39까지 실제 도달 가능

5. Depth 39까지 누적 보증금 (양측 합계, 최악)
   ≈ 300 ETH (한쪽 150 ETH)

6. 최악의 경우 대비 자본
   = 150 ETH (한쪽 부담)
```

**결론**: Sepolia에서 **최악의 경우 대비 자본 = 150 ETH**

#### Mainnet 자본 계산 과정

**설정값**:
```json
{
  "faultGameMaxDepth": 73,
  "faultGameSplitDepth": 30,
  "l2OutputOracleSubmissionInterval": 1800,
  "l2BlockTime": 2
}
```

**계산**:
```
1. 게임 범위
   = 1800 blocks × 2초 = 60분 분량

2. 예상 VM Instructions (최악의 경우)
   = 1800 blocks × 30M gas/block × 100 instructions/gas
   = 5.4T instructions

3. 필요 Depth
   = log₂(5.4T) ≈ 42

4. MAX_GAME_DEPTH = 73이므로
   → Depth 42까지 실제 도달 가능

5. Depth 42까지 누적 보증금 (양측 합계, 최악)
   ≈ 600 ETH (한쪽 300 ETH)

6. 최악의 경우 대비 자본
   = 300 ETH (한쪽 부담)
```

**결론**: Mainnet에서 **최악의 경우 대비 자본 = 300 ETH**

---

#### 환경별 필요 자본 요약 (설정 기반 계산)

| 환경 | 설정 출처 | MAX_DEPTH | 제출 간격 | 최악 도달 Depth | 필요 자본 |
|------|----------|-----------|----------|----------------|----------|
| **Devnet** | devnetL1.json | 50 | 10 blocks (20초) | ~35 | **50 ETH** |
| **Sepolia** | thanos-sepolia.json | 73 | 120 blocks (24분) | ~39 | **150 ETH** |
| **Mainnet** | mainnet.json | 73 | 1800 blocks (60분) | ~42 | **300 ETH** |

**계산 근거**:
- 보증금 공식: `assumedBaseFee × requiredGas(depth)`
- 누적 보증금: Depth 0부터 해당 depth까지의 합
- 한쪽 부담: 홀수 또는 짝수 depth만 담당 (Binary Search)

**⚠️ 이것은 추정이 아닌 설정 기반 계산입니다**:
- 각 환경의 실제 설정값 사용
- 보증금 공식 적용
- 최악의 경우(모든 depth까지 진행) 가정
- Challenger는 이 자본을 확보해야 안전하게 운영 가능

### 5.3 ROI 분석

#### Challenger의 예상 수익

```
시나리오: 월 10개 악의적 제안 적발

Devnet:
- 월 수익: 18 ETH × 10 = 180 ETH
- 초기 자본: 50 ETH
- ROI: 360% / month

Sepolia:
- 월 수익: 32 ETH × 10 = 320 ETH
- 초기 자본: 100 ETH
- ROI: 320% / month

Mainnet:
- 월 수익: 56 ETH × 10 = 560 ETH
- 초기 자본: 200 ETH
- ROI: 280% / month

⚠️ 실제로는 악의적 제안이 거의 없음
→ 보증금 시스템의 억제 효과
→ Challenger 수익 모델은 비현실적
→ Foundation/DAO 지원 필요
```

---

## 6. 권장사항

### 6.1 환경별 최적 설정

#### Devnet 권장사항

```yaml
제출 주기: 20-60초
  - 빠른 테스트 사이클
  - 낮은 보증금 부담

MAX_GAME_DEPTH: 50
  - 테스트에 충분
  - 빠른 게임 진행

SPLIT_DEPTH: 14
  - Output Root 빠른 검증
  - Trace 테스트 가능

권장 자본: 10-20 ETH
  - 대부분의 시나리오 커버
```

#### Sepolia 권장사항

```yaml
제출 주기: 10-30분
  - 현실적 테스트
  - 합리적 보증금

MAX_GAME_DEPTH: 73
  - 프로덕션 준비
  - 완전한 보안 검증

SPLIT_DEPTH: 30
  - 충분한 Output 검증
  - 실전 Trace 테스트

권장 자본: 50-100 ETH
  - 안전한 여유분
```

#### Mainnet 권장사항 (Optimism 공식 설정 기준)

```yaml
제출 주기: 60분 (1800 blocks)
  - Optimism mainnet.json 공식 설정
  - 가스비 효율
  - 보수적 운영

MAX_GAME_DEPTH: 73
  - 최대 보안
  - Optimism 표준

SPLIT_DEPTH: 30
  - 균형잡힌 검증
  - Optimism 표준

권장 자본: 200-500 ETH
  - 전문 운영자 기준
  - 위기 대응 가능

출처: packages/contracts-bedrock/deploy-config/mainnet.json
```

### 6.2 제출 주기 선택 가이드

```
선택 기준:

1. Finality 요구사항
   빠른 출금 필요 → 짧은 주기 (5-15분)
   안정성 우선 → 긴 주기 (30-60분)

2. 가스비 예산
   낮은 예산 → 긴 주기
   충분한 예산 → 짧은 주기

3. 네트워크 특성
   L1 혼잡도 높음 → 긴 주기
   L1 안정적 → 짧은 주기

4. 보안 요구사항
   높은 보안 → 짧은 주기 (빠른 검증)
   균형 → 중간 주기
```

### 6.3 보증금 준비 전략

#### 개인 Challenger

```
1단계: 소규모 시작 (10-20 ETH)
  └─ Devnet/Sepolia에서 경험 축적

2단계: 점진적 확대 (50-100 ETH)
  └─ 실전 Sepolia 참여

3단계: 전문 운영 (200+ ETH)
  └─ Mainnet 준비
```

#### 기관/재단

```
초기 자본: 500-1000 ETH
운영 모델:
  - 24/7 모니터링
  - 자동화된 Challenger
  - 리스크 관리
  - DAO Treasury 연동
```

### 6.4 실무 체크리스트

#### Proposer 운영 시

- [ ] 제출 주기 결정 (가스비 vs finality)
- [ ] 보증금 자본 준비 (환경별 권장치)
- [ ] 모니터링 시스템 구축
- [ ] 백업 자금 확보 (2-3배)
- [ ] 네트워크 상태 체크 (동기화, 가스비)

#### Challenger 운영 시

- [ ] 노드 동기화 확인 (L1, L2, Rollup)
- [ ] 충분한 자본 준비 (환경별 권장치)
- [ ] 자동화 스크립트 준비
- [ ] Prestate 검증 (올바른 버전)
- [ ] 테스트 실행 (Devnet/Sepolia)

---

## 7. 요약 비교표

### 7.1 환경별 핵심 비교

| 항목 | Devnet (Tokamak) | Sepolia (Tokamak) | Mainnet (Optimism) |
|------|-----------------|-------------------|-------------------|
| **목적** | 개발/테스트 | 공개 테스트 | 프로덕션 |
| **출처** | devnetL1.json | thanos-sepolia.json | mainnet.json |
| **제출 주기** | 20초 | 24분 | 60분 |
| **MAX_DEPTH** | 50 | 73 | 73 |
| **SPLIT_DEPTH** | 14 | 30 | 30 |
| **일반 분쟁 비용** | ~18 ETH | ~32 ETH | ~56 ETH |
| **권장 자본** | 20-50 ETH | 50-150 ETH | 100-300 ETH |
| **게임 시간** | 2.5분 | 3.5일 | 3.5일 |
| **출금 대기** | 1시간 | 7일 | 7일 |

### 7.2 제출 주기별 영향

| 제출 주기 | 게임 범위 | 필요 Depth | 보증금 수준 | 권장 환경 |
|----------|----------|-----------|-----------|----------|
| **< 1분** | 매우 작음 | 낮음 (~25) | 낮음 (~10 ETH) | Devnet |
| **10-30분** | 중간 | 중간 (~31) | 중간 (~32 ETH) | Sepolia |
| **> 60분** | 큼 | 높음 (~35) | 높음 (~56 ETH) | Mainnet |

### 7.3 최종 권장사항

```
🎯 Optimism 표준: MAX_GAME_DEPTH = 73 ⭐
   - Mainnet: mainnet.json
   - Sepolia: thanos-sepolia.json
   - 이것이 프로덕션 표준

⚠️ Tokamak Devnet: MAX_GAME_DEPTH = 50 (테스트 전용)
   - 빠른 테스트를 위한 예외
   - 프로덕션에서는 사용 금지

환경별 권장:

✅ Devnet (Tokamak):
   - 설정: 빠른 주기 (20초) + 낮은 depth (50)
   - 자본: 50 ETH
   - 용도: 개발 효율 최대화

✅ Sepolia (Tokamak):
   - 설정: 중간 주기 (24분) + 표준 depth (73)
   - 자본: 150 ETH
   - 용도: 현실적 테스트

✅ Mainnet (Optimism):
   - 설정: 긴 주기 (60분) + 표준 depth (73)
   - 자본: 300 ETH
   - 용도: 보안 + 효율 균형
   - 출처: mainnet.json (공식)
```

---

## 📚 참고 자료

### 설정 파일 출처

1. **Tokamak Devnet**
   - 파일: `packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json`
   - 용도: 로컬 개발 및 테스트

2. **Tokamak Sepolia Testnet**
   - 파일: `packages/tokamak/contracts-bedrock/deploy-config/thanos-sepolia.json`
   - 용도: 공개 테스트넷

3. **Optimism Mainnet** ⭐
   - 파일: `packages/contracts-bedrock/deploy-config/mainnet.json`
   - 용도: 프로덕션 환경 참조
   - 버전: `@eth-optimism/contracts-bedrock@0.17.3`

### 보증금 계산 로직

- **소스 코드**: `packages/contracts-bedrock/src/dispute/FaultDisputeGame.sol:702-743`
- **스펙**: Big Bonds v1.5 (TM)
- **상수**:
  - `assumedBaseFee = 200 gwei` (하드코딩)
  - `baseGasCharged = 400,000 gas`
  - `highGasCharged = 300,000,000 gas`

---
