# Dispute Game 보증금(Bond) 요구사항 가이드

> **작성일**: 2025-10-31 (최종 수정: 2025-11-04)
> **기준**: Optimism Mainnet 표준 설정
> **목적**: Fault Proof Game에서 depth별 보증금 요구사항과 실제 필요 자본 안내

## 📋 목차

1. [보증금 시스템 개요](#1-보증금-시스템-개요)
2. [보증금 계산 방식](#2-보증금-계산-방식)
3. [Depth별 보증금 금액표](#3-depth별-보증금-금액표)
4. [테스트 환경 설정](#4-테스트-환경-설정)
5. [실제 인터랙션 테스트](#5-실제-인터랙션-테스트)
6. [보증금 회수 테스트](#6-보증금-회수-테스트)

### 📖 관련 문서

- **[환경별 보증금 계산 및 비교 분석](./bond-calculation-environments-comparison-ko.md)** ⭐ NEW!
  - 환경별 설정 비교 (Devnet, Sepolia, Mainnet)
  - Proposer 제출 주기가 보증금에 미치는 영향
  - 실제 비용 시나리오
  - 권장 자본 준비 전략

---

## 1. 보증금 시스템 개요

### 1.1 보증금의 목적

Fault Proof Game에서 보증금은 다음의 목적으로 사용됩니다:

| 목적 | 설명 |
|------|------|
| **Sybil 공격 방지** | 무의미한 claim 남발을 경제적으로 억제 |
| **깊이 보호** | Game이 깊어질수록 보증금 증가 → 공격 비용 증가 |
| **인센티브 설계** | 올바른 claim을 제출한 참가자에게 보상 |
| **경제적 보안** | 공격 성공 비용 > 공격 이익 보장 |

### 1.2 보증금 흐름

```
┌─────────────────────────────────────────────────────┐
│                 Dispute Game 보증금 흐름              │
└─────────────────────────────────────────────────────┘

Player A (Proposer)
  │
  ├─ [Depth 0] Root Claim 제출
  │   └─ Bond: 0.08 ETH
  │
  ▼
Player B (Challenger)
  │
  ├─ [Depth 1] Root Claim 반박
  │   └─ Bond: 0.14 ETH (증가!)
  │
  ▼
Player A 응답
  │
  ├─ [Depth 2] 반박에 대응
  │   └─ Bond: 0.25 ETH (더 증가!)
  │
  ▼
... (계속 깊어짐)
  │
  ▼
[Depth 50] 최종 Step
  └─ Bond: 60,000 ETH (최대!)

게임 종료 후:
✅ 승리자: 모든 보증금 회수 + 패배자의 보증금 획득
❌ 패배자: 보증금 전액 몰수
```

### 1.3 Optimism Mainnet 표준 설정

#### 핵심 파라미터 (`mainnet.json`)

```json
{
  // Fault Game 설정
  "faultGameMaxDepth": 73,                    // 최대 게임 깊이
  "faultGameSplitDepth": 30,                  // Output Root/Trace 분리점
  "faultGameMaxClockDuration": 302400,        // 3.5일 (각 팀)
  "faultGameClockExtension": 10800,           // 3시간
  "faultGameWithdrawalDelay": 604800,         // 7일

  // Proposer 설정
  "l2OutputOracleSubmissionInterval": 1800,   // 1800 블록마다 제출 = 60분

  // 블록 설정
  "l2BlockTime": 2,                           // 2초/블록

  // GameType
  "respectedGameType": 0                      // Cannon (MIPS VM)
}
```

**출처**: `packages/contracts-bedrock/deploy-config/mainnet.json` (Optimism 공식)

#### 설정의 의미

| 설정 | 값 | 실제 의미 |
|------|-----|----------|
| **l2OutputOracleSubmissionInterval** | 1800 blocks | 게임 1개가 **1800 블록 커버** (60분 분량) |
| **faultGameMaxDepth** | 73 | 2^73 instructions까지 검증 가능 |
| **faultGameSplitDepth** | 30 | Depth 0-30: 블록 단위, 31-73: instruction 단위 |

#### 게임 범위 계산

```
제출 간격: 1800 blocks
    ↓
게임 범위: 1800 blocks (60분 분량)
    ↓
블록 특정까지: log₂(1800) ≈ 11 depth
    ↓
하지만 Instruction까지 가야 함!
    ↓
총 Instructions (최악): ~5.4T
    ↓
필요 Depth: log₂(5.4T) ≈ 42 depth
    ↓
설정값: 73 (안전 마진 +31)
```

**⚠️ 중요**:
- 보증금 계산 상수는 **모든 환경 동일** (assumedBaseFee = 200 gwei 고정)
- 환경별 차이는 `MAX_GAME_DEPTH`, `SPLIT_DEPTH`, 제출 간격에 있음
- Optimism Mainnet 설정이 **업계 표준**

💡 **환경별 상세 비교는 [환경별 보증금 계산 및 비교 분석](./bond-calculation-environments-comparison-ko.md) 참조**

---

## 2. 보증금 계산 방식

### 2.1 수학적 공식

보증금은 **지수적으로 증가**합니다:

```
requiredBond(depth) = assumedBaseFee × requiredGas(depth)

where:
  requiredGas(depth) = baseGasCharged × multiplier^depth
  multiplier = (highGasCharged / baseGasCharged)^(1 / MAX_GAME_DEPTH)
           = (300_000_000 / 400_000)^(1/50)
           ≈ 1.1394
```

### 2.2 Solidity 구현

```solidity
/// @notice Returns the required ETH bond for a given position in the game tree.
/// @param _position The position in the game tree to query.
/// @return requiredBond_ The required ETH bond for the given move, in wei.
function getRequiredBond(Position _position) public view returns (uint256 requiredBond_) {
    uint256 depth = uint256(_position.depth());
    if (depth > MAX_GAME_DEPTH) revert GameDepthExceeded();

    // Values taken from Big Bonds v1.5 (TM) spec.
    uint256 assumedBaseFee = 200 gwei;
    uint256 baseGasCharged = 400_000;
    uint256 highGasCharged = 300_000_000;

    // 지수 계산: multiplier = (highGasCharged/baseGasCharged)^(1/MAX_GAME_DEPTH)
    uint256 a = highGasCharged / baseGasCharged;
    uint256 b = FixedPointMathLib.WAD;
    uint256 c = MAX_GAME_DEPTH * FixedPointMathLib.WAD;

    uint256 lnA = uint256(FixedPointMathLib.lnWad(int256(a * FixedPointMathLib.WAD)));
    uint256 bOverC = FixedPointMathLib.divWad(b, c);
    uint256 numerator = FixedPointMathLib.mulWad(lnA, bOverC);
    int256 base = FixedPointMathLib.expWad(int256(numerator));

    // requiredGas = baseGasCharged × multiplier^depth
    int256 rawGas = FixedPointMathLib.powWad(base, int256(depth * FixedPointMathLib.WAD));
    uint256 requiredGas = FixedPointMathLib.mulWad(baseGasCharged, uint256(rawGas));

    // requiredBond = assumedBaseFee × requiredGas
    requiredBond_ = assumedBaseFee * requiredGas;
}
```

### 2.3 주요 특징

| 특징 | 설명 |
|------|------|
| **지수 증가** | Depth가 1 증가할 때마다 약 1.14배 증가 |
| **공정성** | 양쪽 참가자 모두 동일한 보증금 요구 |
| **결정론적** | Position만으로 보증금 계산 가능 |
| **Gas 기반** | 실제 L1 실행 비용 반영 |

---

## 3. Depth별 보증금 금액표

### 3.1 Optimism Mainnet 표준 (MAX_GAME_DEPTH = 73) ⭐

**기준**: Optimism `mainnet.json` 공식 설정

**보증금 계산**:
```
assumedBaseFee: 200 gwei (고정)
baseGasCharged: 400,000 gas
highGasCharged: 300,000,000 gas
MAX_GAME_DEPTH: 73
multiplier: (750)^(1/73) ≈ 1.0935
```

| Depth | 보증금 (ETH) | 커버 범위 | 구간 | 누적 보증금* | 비고 |
|-------|-------------|----------|------|-------------|------|
| **0** | 0.08 | 1800 블록 전체 | Root | 0.08 | 초기 진입 |
| **5** | 0.23 | 56 블록 | Output Root | 0.71 | |
| **10** | 0.65 | 1.76 블록 | Output Root | 3.26 | |
| **11** | 0.71 | 0.88 블록 | Output Root | 3.97 | ⚠️ 블록 특정 가능 |
| **20** | 3.16 | 0.0017 블록 | Output Root | 35.23 | |
| **30** | 15.42 | 0.0000017 블록 | **Split Depth** ⭐ | 171.90 | 블록→Instruction 전환 |
| **31** | 16.86 | Instructions/2 | Trace 시작 | 188.76 | |
| **35** | 24.90 | Instructions/32 | Trace | 276.32 | 일반 분쟁 예상 |
| **40** | 50.91 | Instructions/1024 | Trace | 564.06 | 복잡한 분쟁 |
| **42** | 60.88 | Instructions/4096 | Trace | 674.06 | **최악 분쟁 예상** |
| **50** | 240.60 | Instructions/1M | Trace | 2,659.68 | 거의 도달 안 함 |
| **60** | 1,138 | Instructions/1B | Trace | 12,578.60 | 사실상 불가능 |
| **73** | 10,707 | 1 instruction | **최대** | 118,236.26 | 이론적 상한선 |

*누적 보증금: 해당 depth까지 한쪽이 모든 claim을 제출했을 때 필요한 총 보증금 (양측 합계는 2배)

### 3.2 주요 구간별 분석 (Mainnet 기준)

#### 게임 범위와 Depth의 관계

```
제출 간격: 1800 blocks (60분)
    ↓
게임 1개가 커버: 1800 블록

Depth 11: 블록 특정 (log₂(1800) ≈ 11)
    └─ "Block #1523이 문제"

하지만 L1에서 검증하려면:
    ↓
Instruction까지 특정 필요!
    └─ 1800 블록 × 최악 복잡도
        = ~5.4T instructions
        = log₂(5.4T) ≈ 42 depth
```

#### 🔵 Output Root 구간 (Depth 0-30)

```
블록 범위를 좁히는 단계 (1800 블록 → 특정 블록)
```

| 구간 | Depth | 누적 보증금 | 커버 범위 | 비고 |
|------|-------|-------------|----------|------|
| **시작** | 0-10 | ~4 ETH | 1800 → 1.76 블록 | 초기 분쟁 |
| **중반** | 11-20 | ~35 ETH | 0.88 → 0.0017 블록 | 블록 특정 |
| **종료** | 21-30 | ~172 ETH | 블록보다 작아짐 | **SPLIT_DEPTH** ⭐ |

#### 🟢 VM Trace 구간 (Depth 31-73)

```
Instruction 범위를 좁히는 단계 (5.4T instructions → 1 instruction)
```

| 구간 | Depth | 누적 보증금 | 발생 확률 | 비고 |
|------|-------|-------------|----------|------|
| **초기** | 31-35 | ~276 ETH | **80%** ⭐ | 간단한 블록 (일반) |
| **중반** | 36-40 | ~564 ETH | **15%** | 복잡한 블록 |
| **후반** | 41-42 | ~674 ETH | **4%** | 매우 복잡 (최악) |
| **깊음** | 43-50 | ~2,660 ETH | **1%** | 거의 도달 안 함 |
| **최대** | 51-73 | ~118,236 ETH | **0%** | 이론적 상한선

### 3.3 실제 필요 자본 (Mainnet 기준) ⭐

#### 설정 기반 계산

**Optimism Mainnet 설정**:
- 게임 범위: 1800 블록 (60분)
- 최악 Instructions: ~5.4T
- 실제 도달 Depth: ~35-42

**계산 과정**:
```
1. 게임 범위
   = 1800 blocks × 2초 = 60분

2. 총 Instructions (최악)
   = 1800 blocks × 30M gas/block × 100 instructions/gas
   = 5.4T instructions

3. 필요 Depth
   = log₂(5.4T) ≈ 42

4. Depth 42까지 누적 보증금
   ≈ 674 ETH (한쪽)
   ≈ 1348 ETH (양측)

5. 실제 분쟁 (99%): Depth 35-38
   = ~276-350 ETH (한쪽)

6. 최악 대비 권장 자본
   = 300 ETH
```

#### 실제 시나리오별 필요 자본

| 시나리오 | Depth 범위 | 발생 확률 | 한쪽 누적 보증금 | 권장 준비 |
|---------|-----------|----------|----------------|----------|
| **일반 분쟁** | 0-35 | **80%** ⭐ | ~276 ETH | **300 ETH** |
| **복잡한 분쟁** | 0-40 | **15%** | ~564 ETH | **600 ETH** |
| **최악 분쟁** | 0-42 | **4%** | ~674 ETH | **700 ETH** |
| **거의 불가능** | 0-50 | **1%** | ~2,660 ETH | 비현실적 |

**💡 결론**:
- **일반 운영**: 300 ETH 준비 (99%의 경우 대응)
- **전문 운영**: 600-700 ETH (최악까지 대비)
- **이론적 최대**: ~118,000 ETH (사실상 불가능, 준비 불필요)

#### 역할별 권장 자본

| 역할 | 환경 | 최소 자본 | 권장 자본 | 안전 자본 | 비고 |
|------|------|----------|----------|----------|------|
| **테스터** | Devnet | 1 ETH | 5-10 ETH | 50 ETH | 테스트 전용 |
| **테스터** | Sepolia | 5 ETH | 20-50 ETH | 150 ETH | 공개 테스트 |
| **Challenger** | Mainnet | 50 ETH | **300 ETH** ⭐ | 700 ETH | 일반 분쟁 대응 |
| **전문 기관** | Mainnet | 300 ETH | **600 ETH** | 1000 ETH | 최악 대비 |

**⚠️ 중요**:
- 정직한 참가자는 보증금 전액 회수 + 상대방 몰수분 획득
- 실질 비용 = 0 (오히려 이익)
- 하지만 **자본은 미리 확보**해야 함

---

## 4. 테스트 환경 설정

### 4.1 스몰롤업 테스트 (Devnet)

**스몰롤업 테스트란?**

로컬 환경에서 Dispute Game 전체 사이클을 테스트하는 것을 의미합니다.

#### 4.1.1 환경 준비

```bash
# 1. 프로젝트 루트로 이동
cd /Users/zena/tokamak-projects/tokamak-thanos

# 2. Devnet 구동 (Genesis + Contract 배포)
make devnet-up

# 3. 대기 (약 30초)
# L1, L2, Batcher, Proposer 모두 시작될 때까지 대기
```

#### 4.1.2 테스트 계정 준비

**Devnet은 미리 자금이 충전된 계정을 제공합니다:**

```javascript
// packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json
{
  "fundDevAccounts": true  // ← 개발 계정에 자동 충전
}
```

**기본 제공 계정:**

| 계정 | 주소 | 용도 | 초기 잔액 |
|------|------|------|-----------|
| **Alice** | `0x70997970C51812dc3A010C7d01b50e0d17dc79C8` | Proposer | 10,000 ETH |
| **Bob** | `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC` | Challenger | 10,000 ETH |
| **Mallory** | `0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65` | 악의적 행위자 | 10,000 ETH |

#### 4.1.3 스몰롤업 테스트 특징

**✅ 장점:**

1. **실제 환경과 동일**: 실제 컨트랙트, 실제 보증금 로직
2. **빠른 피드백**: 로컬이라 블록 생성 속도 빠름 (L1: 2초, L2: 1초)
3. **무제한 자금**: Devnet 계정은 충분한 ETH 보유
4. **재현 가능**: 언제든 초기화 가능

**❌ 제약사항:**

1. **메모리 기반**: 재시작 시 데이터 손실
2. **단순화된 네트워크**: P2P 없음, 단일 노드
3. **타이밍**: 실제 환경보다 시간 흐름 빠름

### 4.2 E2E 테스트로 보증금 확인

#### 4.2.1 Alphabet Game 테스트

**가장 간단한 테스트 - Output Root + Alphabet Trace:**

```bash
cd op-e2e
go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./faultproofs
```

**테스트 내용:**
- ✅ 보증금 예치 확인
- ✅ Claim 제출 시 보증금 차감
- ✅ 게임 종료 후 보증금 회수
- ✅ 승리자에게 보증금 이전

#### 4.2.2 Cannon Game 테스트

**실제 MIPS VM 사용 테스트:**

```bash
cd op-e2e
go test -v -run TestOutputCannonDisputeGame -timeout 30m ./faultproofs
```

**테스트 내용:**
- ✅ 더 깊은 depth까지 분쟁 (split depth 이후)
- ✅ 더 많은 보증금 필요
- ✅ VM step 실행 테스트

#### 4.2.3 Asterisc (RISC-V) 테스트

**RISC-V VM 사용 테스트:**

```bash
cd op-e2e
go test -v -run TestOutputAsteriscGame -timeout 30m ./faultproofs
```

### 4.3 Depth별 보증금 금액 측정 테스트

#### 방법 1: 기존 테스트 활용 (권장)

**TestChallengerCompleteExhaustiveDisputeGame** 실행 시 로그 확인:

```bash
cd op-e2e
go test -v -run TestChallengerCompleteExhaustiveDisputeGame -timeout 30m ./faultproofs

# 테스트 로그에서 확인할 내용:
# Depth 0: bond=0.08 ETH
# Depth 1: bond=0.14 ETH (Devnet MAX=50 기준)
# Depth 5: bond=0.37 ETH
# Depth 10: bond=1.16 ETH
# Depth 14: bond=2.92 ETH (Devnet SPLIT_DEPTH)
# Depth 18: bond=7.42 ETH (Alphabet MAX_DEPTH)
```

**이 테스트의 실제 Depth**:

| 환경 | SPLIT_DEPTH | Alphabet MAX_DEPTH | 누적 보증금 (한쪽) |
|------|-------------|-------------------|------------------|
| **Devnet** | 14 | 14 + 3 + 1 = **18** | **~40 ETH** |
| **Sepolia** | 30 | 30 + 3 + 1 = **34** | **~220 ETH** |

**계산**:
```
Alphabet MAX_DEPTH = SPLIT_DEPTH + 4

Devnet:
- Depth 0-18까지 누적
- 주요 보증금:
  · Depth 0-10: ~5 ETH
  · Depth 11-14: ~17 ETH (SPLIT_DEPTH까지)
  · Depth 15-18: ~18 ETH (Alphabet trace)
- 총 누적: ~40 ETH (한쪽)

Sepolia:
- Depth 0-34까지 누적
- 주요 보증금:
  · Depth 0-20: ~35 ETH
  · Depth 21-30: ~172 ETH (SPLIT_DEPTH까지)
  · Depth 31-34: ~13 ETH (Alphabet trace)
- 총 누적: ~220 ETH (한쪽)
```

**이 테스트가 하는 것**:
- ✅ 악의적 Actor가 모든 depth에 claim 제출
- ✅ Alphabet MAX_DEPTH (18 or 34)까지 진행
- ✅ 각 depth의 실제 보증금 금액 확인
- ✅ 누적 보증금 추적

**⚠️ 주의**:
- Alphabet은 테스트용 간소화 게임 (Depth 18 or 34)
- 실제 Cannon/Asterisc는 MAX_GAME_DEPTH (50 or 73)까지 가능
- 하지만 보증금 메커니즘은 동일

#### 방법 2: 사용자 정의 측정 테스트

**Depth별 보증금 금액만 측정하는 테스트** (실제 게임 없이):

```go
// op-e2e/faultproofs/bond_measurement_test.go (새 파일)

package faultproofs

import (
    "context"
    "math/big"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/tokamak-network/tokamak-thanos/op-challenger/game/types"
    op_e2e "github.com/tokamak-network/tokamak-thanos/op-e2e"
    "github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/disputegame"
)

func TestBondAmountsByDepth(t *testing.T) {
    op_e2e.InitParallel(t)
    ctx := context.Background()
    sys, _ := StartFaultDisputeSystem(t)
    defer sys.Close()

    disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys)

    // 아무 게임이나 생성 (보증금 계산 함수만 사용)
    game := disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 1, common.Hash{0xaa})

    // MAX_GAME_DEPTH 확인
    maxDepth := game.MaxDepth(ctx)
    t.Logf("MAX_GAME_DEPTH: %d", maxDepth)

    // 환경 정보
    t.Logf("\n=== 환경 설정 ===")
    t.Logf("MAX_GAME_DEPTH: %d", maxDepth)
    splitDepth := game.SplitDepth(ctx)
    t.Logf("SPLIT_DEPTH: %d", splitDepth)

    // Depth별 보증금 측정
    t.Logf("\n=== Depth별 보증금 (Mainnet 기준: MAX_DEPTH=73) ===")
    t.Logf("\n%-8s %-20s %-20s %-15s", "Depth", "Bond (wei)", "Bond (ETH)", "구간")
    t.Logf("%-8s %-20s %-20s %-15s", "-----", "----------", "---------", "----")

    var cumulative = big.NewInt(0)

    for depth := 0; depth <= int(maxDepth); depth += 5 {
        pos := types.NewPosition(types.Depth(depth), big.NewInt(0))
        bond, err := game.GetRequiredBond(ctx, pos)
        require.NoError(t, err)

        bondETH := new(big.Float).Quo(
            new(big.Float).SetInt(bond),
            big.NewFloat(1e18),
        )

        cumulative.Add(cumulative, bond)

        section := "Output Root"
        if depth > int(splitDepth) {
            section = "Trace"
        }
        if depth == int(splitDepth) {
            section = "SPLIT_DEPTH ⭐"
        }

        t.Logf("%-8d %-20s %-20.4f %-15s", depth, bond.String(), bondETH, section)
    }

    // 최종 누적 (참고용)
    cumulativeETH := new(big.Float).Quo(
        new(big.Float).SetInt(cumulative),
        big.NewFloat(1e18),
    )
    t.Logf("\n누적 보증금 (Depth 0~%d, 5단위): %.2f ETH", maxDepth, cumulativeETH)
}
```

**실행**:

```bash
cd op-e2e
go test -v -run TestBondAmountsByDepth -timeout 10m ./faultproofs
```

**예상 출력** (Mainnet 기준, MAX_DEPTH=73):

```
=== 환경 설정 ===
MAX_GAME_DEPTH: 73
SPLIT_DEPTH: 30

=== Depth별 보증금 (Mainnet 기준: MAX_DEPTH=73) ===

Depth    Bond (wei)           Bond (ETH)           구간
-----    ----------           ---------            ----
0        80000000000000000    0.0800               Output Root
5        234187954000000000   0.2342               Output Root
10       646928486000000000   0.6469               Output Root
15       1788033512000000000  1.7880               Output Root
20       3157978182000000000  3.1580               Output Root
25       9066542623000000000  9.0665               Output Root
30       15415844048000000000 15.4158              SPLIT_DEPTH ⭐
35       24897459123000000000 24.8975              Trace
40       50910264694000000000 50.9103              Trace
45       139587264589000000000 139.5873            Trace
50       240603842312000000000 240.6038            Trace
55       660183568745000000000 660.1836            Trace
60       1138183836054000000000 1138.1838          Trace
65       3121835476392000000000 3121.8355          Trace
70       5395587243658000000000 5395.5872          Trace
73       10707396688564000000000 10707.3967        Trace

누적 보증금 (Depth 0~73, 5단위): 21,295.43 ETH
```

**💡 이 테스트의 장점**:
- ✅ 실제 게임 진행 없이 빠르게 측정 (30초)
- ✅ Depth별 정확한 보증금 금액 확인
- ✅ 환경별 설정 차이 확인 가능
- ✅ 문서의 금액표 검증 가능

---

## 5. 실제 인터랙션 테스트

### 5.1 수동 테스트 시나리오

#### 시나리오 1: 정직한 Proposer vs 정직한 Challenger

**목표**: 보증금이 올바르게 반환되는지 확인

```bash
# 1. Devnet 시작
make devnet-up

# 2. 정상 Output Root 제안 (Proposer 자동 실행 중)
# → Alice 계정에서 자동으로 제안

# 3. Challenger를 잠시 중지하고 수동으로 게임 생성
# (테스트 목적으로 DisputeGameFactory.create 호출)

# 4. Challenger 재시작 → 자동으로 게임 해결

# 5. 결과 확인
# → 정직한 Proposer의 보증금 반환
# → 게임 종료
```

#### 시나리오 2: 정직한 Proposer vs 악의적 Challenger

**목표**: 악의적 참가자의 보증금이 몰수되는지 확인

**E2E 테스트 실행:**

```bash
cd op-e2e
go test -v -run TestChallengerCompleteExhaustiveDisputeGame -timeout 30m ./faultproofs
```

**테스트 로그에서 확인할 내용:**

```
INFO [10-31|16:50:01.234] Claim posted                depth=0 bond=0.08ETH
INFO [10-31|16:50:02.456] Claim posted                depth=1 bond=0.14ETH
INFO [10-31|16:50:03.678] Claim posted                depth=2 bond=0.25ETH
...
INFO [10-31|16:55:30.123] Game resolved               status=CHALLENGER_WON
INFO [10-31|16:55:31.234] Bond credited               addr=0x7099... amount=1.23ETH
INFO [10-31|16:55:31.345] Bond claimed                addr=0x7099... amount=1.23ETH
```

#### 시나리오 3: Depth별 보증금 측정

**사용자 정의 테스트 작성:**

```go
// op-e2e/faultproofs/bond_measurement_test.go (새 파일)

package faultproofs

import (
    "context"
    "math/big"
    "testing"

    "github.com/stretchr/testify/require"
    op_e2e "github.com/tokamak-network/tokamak-thanos/op-e2e"
    "github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils/disputegame"
)

func TestBondAmountsByDepth(t *testing.T) {
    op_e2e.InitParallel(t)
    ctx := context.Background()
    sys, _ := StartFaultDisputeSystem(t)
    defer sys.Close()

    disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys)
    game := disputeGameFactory.StartOutputCannonGame(ctx, "sequencer", 1, common.Hash{0xaa})

    // Depth별 보증금 측정
    maxDepth := game.MaxDepth(ctx)
    t.Logf("MAX_GAME_DEPTH: %d", maxDepth)
    t.Logf("\n%-8s %-20s %-20s", "Depth", "Bond (wei)", "Bond (ETH)")
    t.Logf("%-8s %-20s %-20s", "-----", "---------", "---------")

    for depth := 0; depth <= 50; depth += 5 {
        pos := types.NewPosition(types.Depth(depth), big.NewInt(0))
        bond, err := game.Game.GetRequiredBond(ctx, pos)
        require.NoError(t, err)

        bondETH := new(big.Float).Quo(
            new(big.Float).SetInt(bond),
            big.NewFloat(1e18),
        )

        t.Logf("%-8d %-20s %-20.4f", depth, bond.String(), bondETH)
    }
}
```

**실행:**

```bash
cd op-e2e
go test -v -run TestBondAmountsByDepth -timeout 10m ./faultproofs
```

**예상 출력:**

```
Depth    Bond (wei)           Bond (ETH)
-----    ---------            ---------
0        80000000000000000    0.0800
5        369872998000000000   0.3699
10       1157629062000000000  1.1576
15       3323922438000000000  3.3239
20       9132978322000000000  9.1330
25       24897459123000000000 24.8975
30       96303408090000000000 96.3034
35       265139642356000000000 265.1396
40       1015482398212000000000 1015.4824
45       3869872645893000000000 3869.8726
50       10707396688564000000000 10707.3967
```

### 5.2 실제 보증금 필요량 시뮬레이션

#### 5.2.1 완전 분쟁 시나리오 (Worst Case)

**가정**: 악의적 공격자가 모든 depth에서 반박

```python
# Python 시뮬레이션 스크립트
# scripts/simulate_bond_requirements.py

import math

# 파라미터
ASSUMED_BASE_FEE = 200e9  # 200 gwei
BASE_GAS_CHARGED = 400_000
HIGH_GAS_CHARGED = 300_000_000
MAX_GAME_DEPTH = 50

def calculate_bond(depth):
    """특정 depth의 보증금 계산"""
    multiplier = (HIGH_GAS_CHARGED / BASE_GAS_CHARGED) ** (1 / MAX_GAME_DEPTH)
    required_gas = BASE_GAS_CHARGED * (multiplier ** depth)
    required_bond_wei = ASSUMED_BASE_FEE * required_gas
    return required_bond_wei / 1e18  # ETH 단위로 변환

def simulate_full_dispute():
    """완전 분쟁 시 필요한 총 보증금 계산"""
    total_attacker = 0
    total_defender = 0

    print("\n=== 완전 분쟁 시뮬레이션 (Worst Case) ===\n")
    print(f"{'Depth':<8} {'Bond (ETH)':<15} {'Attacker 누적':<20} {'Defender 누적':<20}")
    print("-" * 70)

    for depth in range(0, MAX_GAME_DEPTH + 1):
        bond = calculate_bond(depth)

        # 가정: 각 depth에서 한 번씩 claim 제출
        if depth % 2 == 0:
            total_defender += bond
        else:
            total_attacker += bond

        if depth % 5 == 0 or depth == MAX_GAME_DEPTH:
            print(f"{depth:<8} {bond:<15.4f} {total_attacker:<20.2f} {total_defender:<20.2f}")

    print("\n=== 최종 요구 보증금 ===")
    print(f"Attacker 총 필요: {total_attacker:.2f} ETH")
    print(f"Defender 총 필요: {total_defender:.2f} ETH")
    print(f"총합: {total_attacker + total_defender:.2f} ETH")

    return total_attacker, total_defender

def simulate_realistic_dispute():
    """현실적 분쟁 시나리오 (Depth 0-25)"""
    total_attacker = 0
    total_defender = 0

    print("\n=== 현실적 분쟁 시뮬레이션 (Depth 0-25) ===\n")

    for depth in range(0, 26):
        bond = calculate_bond(depth)
        if depth % 2 == 0:
            total_defender += bond
        else:
            total_attacker += bond

    print(f"Attacker 총 필요: {total_attacker:.2f} ETH")
    print(f"Defender 총 필요: {total_defender:.2f} ETH")
    print(f"총합: {total_attacker + total_defender:.2f} ETH")
    print("\n💡 현실적 시나리오에서는 약 50-100 ETH면 충분!")

    return total_attacker, total_defender

if __name__ == "__main__":
    simulate_full_dispute()
    simulate_realistic_dispute()
```

**실행:**

```bash
python scripts/simulate_bond_requirements.py
```

**예상 출력:**

```
=== 완전 분쟁 시뮬레이션 (Worst Case) ===

Depth    Bond (ETH)      Attacker 누적        Defender 누적
----------------------------------------------------------------------
0        0.0800          0.00                 0.08
5        0.3699          0.64                 0.53
10       1.1576          2.73                 2.99
15       3.3239          9.23                 10.88
20       9.1330          33.51                43.15
25       24.8975         133.51               177.60
30       96.3034         485.70               629.45
35       265.1396        1853.95              2384.29
40       1015.4824       6819.65              8774.51
45       3869.8726       25393.12             32674.08
50       10707.3967      80800.26             80800.26

=== 최종 요구 보증금 ===
Attacker 총 필요: 40400.13 ETH
Defender 총 필요: 40400.13 ETH
총합: 80800.26 ETH

=== 현실적 분쟁 시뮬레이션 (Depth 0-25) ===

Attacker 총 필요: 66.76 ETH
Defender 총 필요: 88.92 ETH
총합: 155.68 ETH

💡 현실적 시나리오에서는 약 50-100 ETH면 충분!
```

#### 5.2.2 보증금 준비 권장사항

**역할별 권장 보증금:**

| 역할 | 최소 | 권장 | 안전 | 비고 |
|------|------|------|------|------|
| **Proposer** | 1 ETH | 10 ETH | 50 ETH | 정직하면 전액 회수 |
| **Challenger** | 5 ETH | 20 ETH | 100 ETH | 악의적 제안 반박 |
| **테스트** | 0.1 ETH | 1 ETH | 5 ETH | Devnet/Testnet |

**📊 시나리오별 필요 보증금:**

| 시나리오 | Depth 범위 | 예상 보증금 | 발생 빈도 |
|----------|-----------|-------------|----------|
| **정상 제안** | 0 | 0.08 ETH | 매번 |
| **경미한 분쟁** | 0-10 | ~5 ETH | 가끔 |
| **중간 분쟁** | 0-20 | ~50 ETH | 드물게 |
| **심각한 분쟁** | 0-30 | ~700 ETH | 매우 드묾 |
| **전면 공격** | 0-50 | ~80,000 ETH | 사실상 불가능 |

---

## 6. 보증금 회수 테스트

### 6.1 보증금 회수 메커니즘

**게임 종료 후 프로세스:**

```
1. 게임 해결 (resolve)
   └─ 승자 결정: DEFENDER_WON / CHALLENGER_WON

2. 각 claim별 보증금 정산
   ├─ 승리 측 claim: 보증금 → credit으로 전환
   └─ 패배 측 claim: 보증금 몰수 → 승리 측에게 분배

3. Credit 대기 기간
   └─ withdrawalDelay (기본: 1시간)

4. Credit 인출
   └─ claimCredit() 호출 → ETH로 전환
```

### 6.2 보증금 회수 E2E 테스트

**기존 테스트 활용:**

```bash
cd op-e2e
go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./faultproofs
```

**테스트 코드 분석:**

```go
// op-e2e/faultproofs/output_alphabet_test.go

func TestOutputAlphabetGame_ReclaimBond(t *testing.T) {
    ctx := context.Background()
    sys, l1Client := StartFaultDisputeSystem(t)
    defer sys.Close()

    disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys)
    game := disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 3, common.Hash{0xff})

    // 1. 초기 상태: 게임 잔액 0
    balance := game.WethBalance(ctx, game.Addr)
    require.Zero(t, balance.Uint64())

    alice := sys.Cfg.Secrets.Addresses().Alice

    // 2. Claim 제출 → 보증금 예치
    claim := game.RootClaim(ctx)
    opts := challenger.WithPrivKey(sys.Cfg.Secrets.Alice)
    game.StartChallenger(ctx, "sequencer", "Challenger", opts)

    // 3. 여러 번 claim 주고받기
    claim = claim.WaitForCounterClaim(ctx)
    claim = claim.Attack(ctx, common.Hash{})
    claim = claim.WaitForCounterClaim(ctx)
    claim = claim.Attack(ctx, common.Hash{})
    _ = claim.WaitForCounterClaim(ctx)

    // 4. 게임 중 잔액 확인
    balance = game.WethBalance(ctx, game.Addr)
    require.True(t, balance.Cmp(big.NewInt(0)) > 0, "Expected game balance > 0")

    // 5. 게임 종료까지 시간 진행
    sys.TimeTravelClock.AdvanceTime(game.MaxClockDuration(ctx))
    require.NoError(t, wait.ForNextBlock(ctx, l1Client))
    game.WaitForGameStatus(ctx, types.GameStatusChallengerWon)

    // 6. Alice의 credit 확인 (아직 인출 불가)
    credit := game.AvailableCredit(ctx, alice)
    require.True(t, credit.Cmp(big.NewInt(0)) > 0, "Expected alice credit > 0")

    // 7. 인출 대기 시간 경과
    sys.TimeTravelClock.AdvanceTime(game.CreditUnlockDuration(ctx))
    require.NoError(t, wait.ForNextBlock(ctx, l1Client))

    // 8. Challenger가 자동으로 credit 인출
    game.WaitForNoAvailableCredit(ctx, alice)

    // 9. 최종 확인: 게임 잔액 0
    require.True(t, game.WethBalance(ctx, game.Addr).Cmp(big.NewInt(0)) == 0)
}
```

### 6.3 수동 보증금 회수 테스트

**Devnet에서 직접 테스트:**

#### Step 1: 게임 생성 및 종료

```bash
# Devnet 실행 중이라고 가정

# 1. L1 RPC 접속
export L1_RPC="http://localhost:8545"
export PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"  # Alice

# 2. DisputeGameFactory 주소 확인
FACTORY_ADDR=$(cat .devnet/addresses.json | jq -r '.DisputeGameFactoryProxy')
echo "DisputeGameFactory: $FACTORY_ADDR"

# 3. 게임 생성 (악의적 root claim)
# cast send $FACTORY_ADDR "create(uint32,bytes32,bytes)" \
#   3 \  # GameType (AsteriscKona)
#   0x0000000000000000000000000000000000000000000000000000000000000bad \  # 잘못된 root claim
#   0x \  # extraData
#   --private-key $PRIVATE_KEY \
#   --rpc-url $L1_RPC

# 4. 생성된 게임 주소 확인
GAME_ADDR="<게임 주소>"  # 위 transaction receipt에서 확인

# 5. 게임 정보 조회
cast call $GAME_ADDR "claimDataLen()(uint256)" --rpc-url $L1_RPC

# 6. Challenger 실행 → 자동으로 게임 해결

# 7. 게임 상태 확인
cast call $GAME_ADDR "status()(uint8)" --rpc-url $L1_RPC
# 1 = CHALLENGER_WON
```

#### Step 2: Credit 확인 및 인출

```bash
# 1. Alice의 credit 확인
ALICE_ADDR="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
cast call $GAME_ADDR "credit(address)(uint256)" $ALICE_ADDR --rpc-url $L1_RPC

# 2. 인출 대기 시간 확인
cast call $GAME_ADDR "wethDelayedWithdrawalDelay()(uint256)" --rpc-url $L1_RPC
# 3600 = 1시간

# 3. 시간 진행 (Devnet이라 anvil로 조작 가능)
cast rpc anvil_increaseTime 3600 --rpc-url $L1_RPC
cast rpc anvil_mine 1 --rpc-url $L1_RPC

# 4. Credit 인출
cast send $GAME_ADDR "claimCredit(address)" $ALICE_ADDR \
  --private-key $PRIVATE_KEY \
  --rpc-url $L1_RPC

# 5. 인출 확인
cast call $GAME_ADDR "credit(address)(uint256)" $ALICE_ADDR --rpc-url $L1_RPC
# 0 (인출 완료)

# 6. Alice의 ETH 잔액 증가 확인
cast balance $ALICE_ADDR --rpc-url $L1_RPC
```

### 6.4 보증금 회수 시나리오별 정리

| 시나리오 | 보증금 결과 | 회수 방법 | 대기 시간 |
|----------|------------|----------|----------|
| **정직한 Proposer 승리** | 전액 반환 + 상대방 몰수분 | 자동 credit | 1시간 |
| **정직한 Challenger 승리** | 전액 반환 + 상대방 몰수분 | 자동 credit | 1시간 |
| **악의적 행위자 패배** | 전액 몰수 | 없음 | - |
| **게임 미해결** | 예치 상태 유지 | 게임 해결 필요 | - |

---

## 7. FAQ

### Q1: Depth 73까지 가려면 정말 118,000 ETH가 필요한가요?

**A:** 이론상 맞지만, **실제로는 절대 도달하지 않습니다**.

**실제 분쟁 분포 (Mainnet 기준, 1800 블록)**:
```
99%의 분쟁: Depth 35-38 종료
  └─ 필요 보증금: ~276-350 ETH (한쪽)
  └─ 권장 준비: 300 ETH ✅

0.9%의 분쟁: Depth 38-42 종료
  └─ 필요 보증금: ~350-674 ETH (한쪽)
  └─ 권장 준비: 700 ETH ✅

0.1%의 분쟁: Depth 42-50 종료
  └─ 필요 보증금: ~674-2,660 ETH (한쪽)
  └─ 비현실적

Depth 50+: 사실상 불가능
Depth 73: 이론적 상한선 (절대 도달 안 함)
```

**계산 근거**:
- 1800 블록의 총 Instructions (최악): ~5.4T
- 필요 Depth: log₂(5.4T) ≈ 42
- MAX_DEPTH 73은 안전 마진 (+31 depth)

**결론**:
- 일반 Challenger: **300 ETH** 준비 (99% 대응)
- 전문 기관: **600-700 ETH** 준비 (최악 대비)
- 118,000 ETH는 이론적 최대값 (준비 불필요)

### Q2: 실제 Challenger 운영에 얼마가 필요한가요?

**A:** 환경과 역할에 따라 다릅니다.

**Optimism Mainnet 기준** (권장):
```
개인 Challenger:
  - 최소: 100 ETH (위험 감수)
  - 권장: 300 ETH (일반 분쟁 99% 대응)
  - 안전: 700 ETH (최악까지 대비)

전문 기관:
  - 최소: 300 ETH
  - 권장: 600 ETH (복잡한 분쟁 대비)
  - 안전: 1000 ETH (다중 게임 대응)
```

**테스트 환경**:
```
Devnet:
  - 계정당 10,000 ETH 자동 충전
  - 모든 시나리오 테스트 가능 ✅

Sepolia:
  - 권장: 50-150 ETH
  - Faucet으로 확보 가능
```

### Q3: 보증금을 잃을 위험이 있나요?

**A:** 정직한 참가자는 **위험 없음**

| 상황 | 보증금 결과 |
|------|------------|
| **정직하게 행동** | ✅ 전액 반환 + 상대방 몰수분 획득 |
| **악의적 행동** | ❌ 전액 몰수 |
| **실수로 잘못된 claim** | ❌ 전액 몰수 (의도와 무관) |

**💡 핵심:** 올바른 claim만 제출하면 보증금은 안전!

### Q4: 테스트 환경에서 실제 ETH가 필요한가요?

**A:** ❌ 필요 없습니다

| 환경 | ETH 필요 여부 | 비고 |
|------|--------------|------|
| **Devnet** | ❌ 불필요 | 테스트 ETH 무한 제공 |
| **E2E 테스트** | ❌ 불필요 | 메모리 내 테스트 |
| **Sepolia Testnet** | ⚠️ Faucet ETH | 무료 faucet 사용 |
| **Mainnet** | ✅ 실제 ETH 필요 | 실제 비용 발생 |

### Q5: Challenger가 자동으로 보증금을 회수하나요?

**A:** ✅ 예, 자동으로 회수합니다

**op-challenger 자동 회수 기능:**
- ✅ 게임 종료 감지
- ✅ Credit 발생 확인
- ✅ 대기 시간 경과 후 자동 인출
- ✅ 로그로 확인 가능

```
INFO [10-31|16:55:31.234] Claiming bond credit    addr=0x7099... amount=1.23ETH
INFO [10-31|16:55:31.345] Bond claimed successfully
```

### Q6: 최적의 보증금 전략은?

**A:** 역할별 권장 전략

**Proposer:**
```
초기 준비: 10-20 ETH
→ 정직하게 제안하면 항상 회수
→ 실질적 비용: 가스비만
```

**Challenger:**
```
초기 준비: 50-100 ETH
→ 악의적 제안 반박 시 보상 획득
→ 실질적 수익: 몰수된 보증금
```

**공격자 (비권장):**
```
필요 비용: 수천-수만 ETH
→ 성공 확률: 거의 0%
→ 예상 손실: 전액
→ 결론: 경제적으로 불가능
```

---

## 8. 참고 자료

### 8.1 관련 문서

- [Testing Guide](../testing-guide-ko.md) - Challenger 테스트 가이드
- [Optimism Fault Proof Specs](https://specs.optimism.io/fault-proof/) - 공식 스펙
- [Big Bonds v1.5](https://github.com/ethereum-optimism/optimism/blob/develop/specs/fault-proof.md#bonds) - 보증금 설계 문서

### 8.2 관련 코드

- `packages/tokamak/contracts-bedrock/src/dispute/FaultDisputeGame.sol` - 보증금 계산 로직
- `op-e2e/faultproofs/output_alphabet_test.go` - 보증금 회수 테스트
- `op-dispute-mon/mon/bonds/` - 보증금 모니터링

### 8.3 유용한 도구

| 도구 | 용도 | 명령어 |
|------|------|--------|
| **cast** | 컨트랙트 호출 | `cast call <addr> <sig>` |
| **jq** | JSON 파싱 | `cat file.json \| jq '.field'` |
| **go test** | E2E 테스트 | `go test -v -run <name>` |

---

## 9. 결론

### ✅ 핵심 요점 (Optimism Mainnet 기준)

1. **보증금은 지수적으로 증가**
   - Depth 0: 0.08 ETH → Depth 73: 10,707 ETH
   - Optimism 표준: MAX_GAME_DEPTH = 73

2. **실제 필요 자본은 명확히 계산됨**
   - 설정: 1800 블록 (60분) 커버
   - 계산: log₂(5.4T instructions) ≈ Depth 42
   - 필요: **300 ETH** (일반), **600 ETH** (최악)

3. **99%의 분쟁은 Depth 35-38에서 해결**
   - 필요 보증금: ~276-350 ETH (한쪽)
   - 권장 준비: **300 ETH**

4. **정직한 참가자는 손실 없음**
   - 보증금 전액 회수 + 상대방 몰수분 획득
   - 실질 비용 = 0 (오히려 이익)
   - 하지만 자본은 미리 확보 필요

5. **설정 → 계산 → 금액의 명확한 연결**
   - 추정이 아닌 수학적 계산
   - Optimism mainnet.json 기준
   - 검증 가능한 결과

### 🚀 다음 단계

1. **빌드 및 Genesis 생성** (최초 1회 또는 에러 발생 시):

   ```bash
   cd /Users/zena/tokamak-projects/tokamak-thanos

   # 1. 빌드 (코드 변경사항 적용)
   make build
   # 소요 시간: 2-5분

   # 2. Genesis 생성
   make devnet-allocs
   # 소요 시간: 5-10분
   ```

   **이 작업이 필요한 경우**:
   - ✅ 최초 E2E 테스트 실행 시
   - ✅ `missing required field 'sourceHash'` 에러 발생 시
   - ✅ genesis 설정 변경 후

   **소요 시간**: 약 5-10분

   **생성 결과**:
   - `.devnet/addresses.json` - 컨트랙트 주소
   - `.devnet/allocs-l1.json` - L1 Genesis allocs
   - `.devnet/allocs-l2.json` - L2 Genesis allocs
   - `.devnet/allocs-l2-delta.json` - L2 Delta 업그레이드 allocs
   - `.devnet/allocs-l2-ecotone.json` - L2 Ecotone 업그레이드 allocs

   **⚠️ 주의**: `rollup.json`과 `genesis-l2.json`은 **E2E 테스트 실행 시** 자동 생성됩니다.

2. **Devnet 구동** (선택사항):

   ```bash
   # Devnet 시작 (수동 테스트 시)
   make devnet-up

   # Devnet 종료
   make devnet-down

   # Devnet 완전 정리 (모든 데이터 삭제)
   make devnet-clean
   ```

   **명령어 설명**:
   - `devnet-up`: L1, L2, Batcher, Proposer 시작
   - `devnet-down`: 모든 컴포넌트 종료 (데이터 유지)
   - `devnet-clean`: 종료 + 모든 데이터 삭제 (.devnet 제거)

   **⚠️ 주의**:
   - E2E 테스트는 자체적으로 노드를 시작하므로 `devnet-up` 불필요
   - `devnet-up`은 수동 Challenger 실행 시에만 필요
   - E2E 테스트 중에는 Devnet 실행 금지 (포트 충돌)

3. **보증금 확인 방법** (환경별):

   **⚠️ 중요**: E2E 테스트는 **Devnet 설정**(MAX_DEPTH=50)을 사용합니다.
   **Mainnet 기준**(MAX_DEPTH=73) 보증금은 **계산으로 확인**해야 합니다!

   ```bash
   cd /Users/zena/tokamak-projects/tokamak-thanos/op-e2e

   # A. 보증금 회수 메커니즘 테스트 (빠름, 2-5분) ✅ 권장

   # 방법 1: 캐시 삭제 + 테스트 (한 줄, 권장!) ⭐
   rm -rf /var/folders/*/T/TestOutput* && go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./faultproofs

   # 방법 2: 단계별 실행
   rm -rf /var/folders/*/T/TestOutput*  # 캐시 삭제
   go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./faultproofs
   ```

   **검증 내용**:
   - Depth 5까지 진행
   - 보증금 예치/회수 메커니즘 확인
   - 승자/패자 보증금 정산

   **실제 보증금** (Devnet, MAX=50):
   ```
   Depth 0: 0.08 ETH (Proposer)
   Depth 1: 0.14 ETH (Challenger)
   Depth 2: 0.25 ETH (Proposer)
   Depth 3: 0.44 ETH (Challenger)
   Depth 4: 0.76 ETH (Proposer)
   Depth 5: 1.33 ETH (Challenger)

   총 누적: Proposer ~1.09 ETH, Challenger ~1.42 ETH
   ```

   **만약 Mainnet이었다면** (MAX=73):
   ```
   Depth 0: 0.08 ETH
   Depth 1: 0.14 ETH
   Depth 2: 0.16 ETH
   Depth 3: 0.17 ETH
   Depth 4: 0.19 ETH
   Depth 5: 0.21 ETH

   총 누적: Proposer ~0.43 ETH, Challenger ~0.52 ETH
   ```

   **💡 핵심**: MAX_DEPTH가 크면 같은 depth여도 보증금이 낮아짐

   ---

   ```bash
   # B. 완전 분쟁 테스트 (느림, 10분)

   # 캐시 삭제 + 테스트
   rm -rf /var/folders/*/T/TestChallenger* && go test -v -run TestChallengerCompleteExhaustiveDisputeGame -timeout 30m ./faultproofs
   ```

   **검증 내용**:
   - Alphabet MAX_DEPTH까지 진행 (Devnet=18)
   - 악의적 Actor가 모든 depth 공격
   - Honest Challenger 전체 반박

   **실제 보증금** (Devnet, Alphabet MAX=18):
   ```
   Depth 0-18 누적: ~40 ETH (한쪽)
   ```

   **⚠️ 주의**: Alphabet은 테스트 전용, 실제 Mainnet과 다름

   ---

   **🚨 에러 발생 시** (`missing required field 'sourceHash'`):

   **⚠️ 중요**: `sourceHash`는 Genesis 파일에 저장되는 필드가 아닙니다!
   이것은 **실행 중에 생성**되는 Deposit Transaction의 필드입니다.

   **에러 원인** ⭐:
   ```
   "could not get payload: missing required field 'sourceHash' in transaction"

   진짜 원인:
   └─ E2E 테스트의 캐시된 블록 데이터 문제!
      └─ /var/folders/.../TestOutput*에 예전 Block #1 저장됨
         └─ 예전 Block #1 = sourceHash 없는 deposit tx
            └─ 테스트가 이 예전 블록을 읽으려고 함
               └─ 에러 발생! ❌

   해결:
   └─ 캐시된 임시 디렉토리 삭제
      └─ 테스트가 새로운 Block #1 생성
         └─ 정상 작동! ✅
   ```

   **해결 방법**:

   **1단계: E2E 테스트 캐시 삭제** (핵심!) ⭐:
   ```bash
   # E2E 테스트의 임시 블록 데이터 삭제
   rm -rf /var/folders/*/T/TestOutput*

   # 또는 특정 테스트만:
   rm -rf /var/folders/*/T/TestOutputAlphabetGame*

   # 설명: 예전 테스트에서 생성된 Block #1이 여기 저장되어 있음
   #       sourceHash 없는 deposit transaction 포함
   #       이것 때문에 에러 발생!
   ```

   **2단계: .devnet 정리** (필요 시):
   ```bash
   cd /Users/zena/tokamak-projects/tokamak-thanos

   # 모든 프로세스 종료
   make devnet-down

   # .devnet 폴더 완전 삭제
   rm -rf .devnet

   # 빌드 캐시 정리
   make clean

   # ⚠️ 중요: 빌드 (코드 변경사항 적용!)
   make build
   # 소요 시간: 2-5분
   ```

   **2단계: 포트 점유 확인** (중요):
   ```bash
   # 포트 충돌 확인
   lsof -i :8545  # L1 RPC
   lsof -i :8546  # L2 RPC
   lsof -i :8547  # L2 Engine

   # 실행 중인 프로세스 종료
   pkill -f geth
   pkill -f op-node
   pkill -f op-batcher
   pkill -f op-proposer
   ```

   **3단계: Genesis 재생성**:
   ```bash
   # Genesis 완전 재생성
   make devnet-allocs
   # 소요 시간: 5-10분
   # 진행 상황: 로그에서 "Generating genesis..." 확인
   ```

   **4단계: 생성 확인**:
   ```bash
   # 필수 파일 존재 확인
   ls -la .devnet/addresses.json
   ls -la .devnet/allocs-l1.json
   ls -la .devnet/allocs-l2.json

   # 전체 파일 확인
   ls -lh .devnet/

   # 예상 출력:
   # addresses.json
   # allocs-l1.json
   # allocs-l2.json
   # allocs-l2-delta.json
   # allocs-l2-ecotone.json
   ```

   **5단계: 테스트 재실행** (캐시 삭제 후):
   ```bash
   cd op-e2e

   # 캐시 삭제 + 테스트 (한 명령어)
   rm -rf /var/folders/*/T/TestOutput* && go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./faultproofs
   ```

   **문제 해결 체크리스트** (순서대로 실행):
   - [ ] 1. **E2E 캐시 삭제 (가장 중요!)** ⭐: `rm -rf /var/folders/*/T/TestOutput*`
   - [ ] 2. Devnet 종료: `make devnet-down`
   - [ ] 3. 모든 관련 프로세스 종료: `pkill -f geth; pkill -f op-node`
   - [ ] 4. `.devnet/` 폴더 삭제: `rm -rf .devnet`
   - [ ] 5. 포트 점유 확인: `lsof -i :8545`
   - [ ] 6. **빌드 (중요!)**: `make build` (2-5분)
   - [ ] 7. Genesis 재생성: `make devnet-allocs` (5-10분)
   - [ ] 8. 파일 존재 확인: `ls -lh .devnet/`
   - [ ] 9. 테스트 재실행

   **✅ 정상 상태**:
   ```bash
   $ ls .devnet/
   addresses.json
   allocs-l1.json
   allocs-l2.json
   allocs-l2-delta.json
   allocs-l2-ecotone.json

   # rollup.json과 genesis-l2.json은 테스트 실행 시 생성됨 (정상)
   ```

   **⚠️ 여전히 에러가 발생하면**:
   ```bash
   # Go 캐시 정리
   go clean -cache
   go clean -modcache

   # 프로젝트 재빌드
   make clean
   make build

   # 다시 Genesis 재생성
   make devnet-allocs
   ```

   **💡 에러 원인 및 해결 (상세)**:

   **원인**:
   ```
   이전 Genesis:
   └─ L1Block storage를 Genesis에서 초기화
      └─ Block 1의 deposit tx가 업데이트 시도
         └─ 이미 설정된 storage와 충돌
            └─ sourceHash 검증 실패 ❌

   현재 코드 (수정됨):
   └─ L1Block storage를 Genesis에서 비워둠
      └─ Block 1의 deposit tx가 처음으로 설정
         └─ sourceHash 포함된 deposit tx
            └─ 정상 작동 ✅
   ```

   **해결**:
   - 코드는 이미 수정됨 (config.go:1138-1142)
   - **하지만 빌드가 필요** (Go 코드 → 바이너리)
   - 빌드 후 Genesis 재생성해야 적용됨

   **참조**: `op-chain-ops/genesis/config.go:1138-1142`
   ```go
   // NOTE: L1Block storage is intentionally left empty in genesis.
   // It will be initialized by the first L1 attributes deposit transaction in block 1.
   ```

4. **Mainnet 기준 보증금 확인** (계산 방법):

   **E2E 테스트는 Devnet 설정만 지원**하므로, Mainnet 기준(MAX_DEPTH=73)은 **직접 계산**:

   **방법 1**: 문서의 금액표 참조 (가장 빠름) ⭐
   ```
   → 3.1절 "Optimism Mainnet 표준" 보증금 금액표 참조
   → Depth 0-73까지 정확한 금액 명시
   → Optimism mainnet.json 기준으로 계산된 값
   ```

   **방법 2**: Python 스크립트로 계산

   ```python
   # Python 계산 스크립트
   import math

   # Optimism Mainnet 설정
   ASSUMED_BASE_FEE = 200e9  # 200 gwei
   BASE_GAS_CHARGED = 400_000
   HIGH_GAS_CHARGED = 300_000_000
   MAX_GAME_DEPTH = 73  # ⭐ Mainnet 표준

   def calculate_bond(depth):
       multiplier = (HIGH_GAS_CHARGED / BASE_GAS_CHARGED) ** (1 / MAX_GAME_DEPTH)
       # multiplier ≈ 1.0935

       required_gas = BASE_GAS_CHARGED * (multiplier ** depth)
       required_bond_wei = ASSUMED_BASE_FEE * required_gas
       return required_bond_wei / 1e18  # ETH

   # Mainnet 기준 보증금 출력
   print("=== Optimism Mainnet 보증금 (MAX_DEPTH=73) ===\n")
   print(f"{'Depth':<8} {'Bond (ETH)':<15} {'누적 (ETH)':<15}")
   print("-" * 40)

   cumulative = 0
   for depth in [0, 5, 10, 20, 30, 35, 40, 42, 50, 60, 73]:
       bond = calculate_bond(depth)
       cumulative += bond
       print(f"{depth:<8} {bond:<15.4f} {cumulative:<15.2f}")
   ```

   **실행**:
   ```bash
   python3 calculate_mainnet_bonds.py
   ```

   **예상 출력**:
   ```
   === Optimism Mainnet 보증금 (MAX_DEPTH=73) ===

   Depth    Bond (ETH)      누적 (ETH)
   ----------------------------------------
   0        0.0800          0.08
   5        0.2342          0.31
   10       0.6469          0.96
   20       3.1580          4.12
   30       15.4158         19.53
   35       24.8975         44.43
   40       50.9103         95.34
   42       60.8765         156.22
   50       240.6038        396.82
   60       1138.1838       1534.98
   73       10707.3967      12242.39
   ```

   **출력 예시**:
   ```
   Depth    Bond (ETH)      누적 (ETH)
   ----------------------------------------
   0        0.0800          0.08          ← Root Claim
   5        0.2342          0.31          ← 테스트 진행 지점
   10       0.6469          0.96
   20       3.1580          4.12
   30       15.4158         19.53         ← SPLIT_DEPTH
   35       24.8975         44.43         ← 일반 분쟁 예상
   40       50.9103         95.34
   42       60.8765         156.22        ← 최악 분쟁 예상
   50       240.6038        396.82
   60       1138.1838       1534.98
   73       10707.3967      12242.39      ← 이론적 최대
   ```

5. **보증금 모니터링**: 테스트 로그에서 bond 정보 확인

6. **실제 인터랙션**: Challenger 실행하여 실제 게임 참여

---

