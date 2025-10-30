# L1에 제출한 데이터가 정말 맞는 정보인지 증명하는 방법

## 목차
1. [핵심 질문](#핵심-질문)
2. [Optimistic Rollup의 검증 메커니즘](#optimistic-rollup의-검증-메커니즘)
3. [재파생(Derivation)의 원리](#재파생derivation의-원리)
4. [Output Root의 역할](#output-root의-역할)
5. [Fault Proof가 증명하는 것](#fault-proof가-증명하는-것)
6. [실제 코드로 보는 검증 과정](#실제-코드로-보는-검증-과정)
7. [공격 시나리오와 방어](#공격-시나리오와-방어)
8. [결론](#결론)

---

## 핵심 질문

> **"L1에 제출한 데이터가 정말 맞는 정보라는 것을 어떻게 증명하지?"**

이 질문은 Optimistic Rollup의 **핵심 보안 메커니즘**에 관한 것입니다.

---

## Optimistic Rollup의 검증 메커니즘

### 기본 원칙

```
┌─────────────────────────────────────────────────────────────┐
│      Optimistic Rollup의 3단계 검증 시스템                    │
└─────────────────────────────────────────────────────────────┘

1. 결정론적 재파생 (Deterministic Derivation)
   → L1 데이터 → L2 상태 (항상 동일한 결과)

2. Output Root 제출
   → Sequencer가 주장하는 L2 상태

3. 챌린지 메커니즘
   → 누구나 잘못된 Output Root에 이의 제기 가능
```

### 핵심 개념

**"L1 데이터 자체를 검증하는 것이 아니라, L1 데이터로부터 파생된 L2 상태를 검증합니다!"**

```
┌─────────────────────────────────────────────────────────────┐
│               검증의 대상                                     │
└─────────────────────────────────────────────────────────────┘

❌ 검증하지 않는 것:
   - L1 배치 데이터 자체의 유효성
   (L1에 있으면 무조건 유효한 것으로 간주)

✅ 검증하는 것:
   - L1 배치 데이터로부터 재파생한 L2 상태
   - Sequencer가 제출한 Output Root
   - Output Root = Hash(StateRoot + StorageRoot + BlockHash)
```

---

## 재파생(Derivation)의 원리

### 1. 재파생이란?

**L1 데이터만으로 L2 상태를 완전히 재구성하는 과정**

```
┌─────────────────────────────────────────────────────────────┐
│                  재파생 프로세스                              │
└─────────────────────────────────────────────────────────────┘

입력: L1 배치 데이터
  ↓
┌──────────────────────────────────────┐
│ 1. L1 스캔                            │
│    - BatchInbox TX 찾기               │
│    - Calldata/Blob에서 배치 추출      │
└──────────────────────────────────────┘
  ↓
┌──────────────────────────────────────┐
│ 2. 배치 디코딩                        │
│    - 압축 해제                        │
│    - 트랜잭션 목록 추출                │
└──────────────────────────────────────┘
  ↓
┌──────────────────────────────────────┐
│ 3. L2 블록 재구성                     │
│    - 트랜잭션 순서대로 실행            │
│    - 상태 전이 계산                   │
└──────────────────────────────────────┘
  ↓
출력: L2 State Root
```

### 2. 결정론적 실행

**핵심: 같은 입력(L1 데이터) → 항상 같은 출력(L2 상태)**

```go
// op-node/rollup/derive/data_source.go
// L1 데이터로부터 L2 배치를 추출하는 과정

func ProcessBatchData(l1Data []byte) (L2Blocks, error) {
    // 1. 배치 디코딩
    batches := DecodeBatches(l1Data)

    // 2. 각 배치를 L2 블록으로 변환
    for _, batch := range batches {
        // ⭐️ 결정론적 실행: 항상 같은 결과
        l2Block := DeriveL2Block(batch)

        // 3. 트랜잭션 실행
        for _, tx := range l2Block.Transactions {
            // ⭐️ EVM 실행: 결정론적
            ExecuteTransaction(tx, currentState)
        }
    }

    return l2Blocks, nil
}
```

### 3. 재파생의 보장

| 요소 | 설명 |
|------|------|
| **L1 데이터 불변성** | L1 블록체인에 저장되면 변경 불가 |
| **실행 결정론** | EVM은 결정론적 상태 머신 |
| **규칙 명확성** | 배치 형식, 실행 규칙 모두 공개 |
| **누구나 검증 가능** | L1 데이터만 있으면 누구나 재파생 가능 |

---

## Output Root의 역할

### 1. Output Root란?

**L2 상태의 암호학적 요약**

```go
// Output Root 구조
OutputRoot = Hash(
    Version,
    StateRoot,        // L2 상태 트리 루트
    StorageRoot,      // 메시지 패서 스토리지 루트
    BlockHash,        // L2 블록 해시
)
```

### 2. Output Root 제출 과정

```
┌─────────────────────────────────────────────────────────────┐
│            Sequencer의 Output Root 제출                       │
└─────────────────────────────────────────────────────────────┘

Sequencer:
  ↓
1. L1 배치 데이터 게시
   ├─ TX to BatchInboxAddress
   └─ Calldata: [압축된 L2 블록들]
  ↓
2. 자체적으로 L2 상태 계산
   ├─ 트랜잭션 실행
   └─ 상태 전이
  ↓
3. Output Root 계산
   ├─ StateRoot 계산
   ├─ StorageRoot 계산
   └─ OutputRoot = Hash(...)
  ↓
4. L2OutputOracle에 제출
   └─ proposeL2Output(OutputRoot, L2BlockNumber, ...)
```

### 3. Output Root의 검증 가능성

**핵심: Output Root는 주장일 뿐, 증명이 아닙니다!**

```
┌─────────────────────────────────────────────────────────────┐
│         Output Root의 검증 (Challenger가 수행)               │
└─────────────────────────────────────────────────────────────┘

Challenger:
  ↓
1. L1에서 배치 데이터 가져오기
   └─ 동일한 L1 TX 읽기
  ↓
2. 로컬에서 L2 재파생
   └─ 동일한 배치 데이터 실행
  ↓
3. 로컬 Output Root 계산
   └─ OutputRoot_local = Hash(StateRoot_local, ...)
  ↓
4. 비교
   ┌────────────────────────────────────────┐
   │ OutputRoot_sequencer vs OutputRoot_local │
   └────────────────────────────────────────┘
       ↓                    ↓
    일치                  불일치
       ↓                    ↓
    OK!               Fraud Proof 제출!
```

---

## Fault Proof가 증명하는 것

### 1. 증명의 대상

```
┌─────────────────────────────────────────────────────────────┐
│            Fault Proof가 증명하는 것                          │
└─────────────────────────────────────────────────────────────┘

❌ 증명하지 않는 것:
   - L1 배치 데이터의 정확성
   - L1 배치 데이터의 유효성

✅ 증명하는 것:
   - "이 L1 배치 데이터로부터 재파생한 L2 상태는 X이다"
   - "Sequencer가 주장한 Output Root Y는 틀렸다"
   - "올바른 Output Root는 X이다"
```

### 2. 증명 과정

#### A. 이분 탐색 (Bisection)

```
┌─────────────────────────────────────────────────────────────┐
│          Fault Dispute Game 이분 탐색 과정                   │
└─────────────────────────────────────────────────────────────┘

Root Claim: OutputRoot_sequencer
  ↓
Challenger: "이것은 틀렸다!"
  ↓
┌──────────────────────────────────────┐
│ Level 0: 전체 trace (2^30 steps)     │
│   Sequencer: "최종 상태는 S_final"   │
│   Challenger: "아니다, S_final'"     │
└──────────────────────────────────────┘
  ↓ 이분 탐색
┌──────────────────────────────────────┐
│ Level 1: 중간 지점 (2^29)            │
│   Sequencer: "중간 상태는 S_mid"     │
│   Challenger: "동의" 또는 "불일치"   │
└──────────────────────────────────────┘
  ↓ 계속 이분 탐색
┌──────────────────────────────────────┐
│ Level 30: 단일 instruction           │
│   Sequencer: "Step N: A → B"        │
│   Challenger: "아니다, A → B'"      │
└──────────────────────────────────────┘
  ↓
┌──────────────────────────────────────┐
│ 온체인 검증                           │
│ - Step N의 pre-state: A              │
│ - Step N의 instruction: I            │
│ - MIPS VM 실행: Execute(A, I)        │
│ - 결과 확인: B vs B'                 │
└──────────────────────────────────────┘
```

#### B. 단일 Step 검증

```go
// op-challenger/game/fault/trace/cannon/provider.go
func (p *CannonTraceProvider) GetStepData(ctx context.Context, pos types.Position)
    ([]byte, []byte, *types.PreimageOracleData, error) {

    // 1. Proof 로드 (로컬에서 생성)
    proof, err := p.loadProof(ctx, traceIndex.Uint64())

    // 2. State Data: pre-state
    value := ([]byte)(proof.StateData)

    // 3. Proof Data: instruction + witness
    data := ([]byte)(proof.ProofData)

    // 4. Preimage Data: 필요한 외부 데이터
    oracleData, err := p.preimageLoader.LoadPreimage(proof)

    return value, data, oracleData, nil
}
```

**온체인 검증:**

```solidity
// FaultDisputeGame.sol (개념적)
function step(uint256 claimIndex, bytes memory stateData, bytes memory proofData) {
    // 1. Pre-state 검증
    bytes32 preStateHash = keccak256(stateData);

    // 2. MIPS VM 실행 (온체인)
    bytes32 postStateHash = MIPS.step(stateData, proofData);

    // 3. 결과 비교
    if (postStateHash == claim.value) {
        // Sequencer 승리
    } else {
        // Challenger 승리
    }
}
```

### 3. 검증의 보장

| 보장 사항 | 설명 |
|----------|------|
| **결정론적 실행** | 같은 L1 데이터 → 항상 같은 L2 상태 |
| **온체인 검증** | 단일 step은 L1에서 직접 실행하여 검증 |
| **암호학적 보장** | Merkle proof로 전체 상태 연결 |
| **경제적 인센티브** | 잘못된 주장은 bond 손실 |

---

## 실제 코드로 보는 검증 과정

### 1. Challenger의 검증 흐름

#### A. Output Root 가져오기 (provider.go:88-94)

```go
func (o *OutputTraceProvider) Get(ctx context.Context, pos types.Position)
    (common.Hash, error) {

    // ⭐️ Honest Block Number: L1 데이터로 도달 가능한 최대 블록
    outputBlock, err := o.HonestBlockNumber(ctx, pos)
    if err != nil {
        return common.Hash{}, err
    }

    // ⭐️ Output Root 계산: 로컬에서 재파생한 결과
    return o.outputAtBlock(ctx, outputBlock)
}
```

#### B. Honest Block Number 계산 (provider.go:72-86)

```go
func (o *OutputTraceProvider) HonestBlockNumber(ctx context.Context,
    pos types.Position) (uint64, error) {

    // 1. Claimed Block Number: Sequencer 주장
    outputBlock, err := o.ClaimedBlockNumber(pos)

    // 2. ⭐️ Safe Head: L1 데이터로 도달 가능한 블록
    resp, err := o.rollupProvider.SafeHeadAtL1Block(ctx, o.l1Head.Number)
    maxSafeHead := resp.SafeHead.Number

    // 3. 최소값 선택 (보수적 접근)
    if outputBlock > maxSafeHead {
        outputBlock = maxSafeHead  // ⭐️ L1 데이터 기반
    }

    return outputBlock, nil
}
```

**핵심:**
- `SafeHeadAtL1Block`: L1 배치 데이터를 재파생한 결과
- Challenger는 L1 데이터로 도달 가능한 블록까지만 주장
- L1 데이터 기반 = 검증 가능 보장

#### C. Output Root 계산 (provider.go:124-130)

```go
func (o *OutputTraceProvider) outputAtBlock(ctx context.Context, block uint64)
    (common.Hash, error) {

    // ⭐️ Rollup Provider: 로컬 op-node가 재파생한 결과
    output, err := o.rollupProvider.OutputAtBlock(ctx, block)
    if err != nil {
        return common.Hash{}, fmt.Errorf("failed to fetch output at block %v: %w",
            block, err)
    }

    return common.Hash(output.OutputRoot), nil
}
```

**데이터 흐름:**

```
Challenger의 op-node
  ↓
L1 RPC로 배치 데이터 가져오기
  ↓
배치 데이터 재파생
  ↓
L2 상태 계산
  ↓
Output Root 생성
  ↓
Sequencer의 Output Root와 비교
  ↓
불일치 시 Fault Proof 제출
```

### 2. Cannon Trace 생성

#### A. 로컬 입력 가져오기 (local.go:47-61)

```go
func FetchLocalInputsFromProposals(
    ctx context.Context,
    l1Head common.Hash,  // ⭐️ L1 Head
    l2Client L2HeaderSource,
    agreedOutput contracts.Proposal,
    claimedOutput contracts.Proposal
) (LocalGameInputs, error) {

    // ⭐️ L2 Header: 합의된 블록의 헤더
    agreedHeader, err := l2Client.HeaderByNumber(ctx, agreedOutput.L2BlockNumber)
    l2Head := agreedHeader.Hash()

    return LocalGameInputs{
        L1Head:        l1Head,                      // ⭐️ L1 기준점
        L2Head:        l2Head,                      // L2 시작점
        L2OutputRoot:  agreedOutput.OutputRoot,     // 합의된 상태
        L2Claim:       claimedOutput.OutputRoot,    // 검증할 주장
        L2BlockNumber: claimedOutput.L2BlockNumber, // 목표 블록
    }, nil
}
```

#### B. Cannon 실행 (executor.go:66-102)

```go
func (e *Executor) generateProof(...) error {
    // ... (생략) ...

    args = append(args,
        "--",
        e.server, "--server",
        "--l1", e.l1,              // ⭐️ L1 RPC
        "--l1.beacon", e.l1Beacon, // ⭐️ L1 Beacon RPC
        "--l2", e.l2,              // L2 RPC (상태 확인용)
        "--datadir", dataDir,

        // ⭐️ 재파생 입력값
        "--l1.head", e.inputs.L1Head.Hex(),
        "--l2.head", e.inputs.L2Head.Hex(),
        "--l2.outputroot", e.inputs.L2OutputRoot.Hex(),
        "--l2.claim", e.inputs.L2Claim.Hex(),
        "--l2.blocknumber", e.inputs.L2BlockNumber.Text(10),
    )

    // ⭐️ op-program 실행: L1 데이터로 L2 재파생
    err = e.cmdExecutor(ctx, e.logger, e.cannon, args...)
    return err
}
```

**op-program의 동작:**

```
1. L1Head부터 L1 스캔
   ↓
2. BatchInbox TX 찾기
   ↓
3. 배치 데이터 추출 (Calldata/Blob)
   ↓
4. L2Head부터 재파생 시작
   ↓
5. 배치 데이터 실행
   ↓
6. L2BlockNumber까지 도달
   ↓
7. Output Root 계산
   ↓
8. L2Claim과 비교
   ↓
9. 불일치 시: Fault Proof 데이터 생성
```

---

## 공격 시나리오와 방어

### 시나리오 1: Sequencer가 잘못된 Output Root 제출

```
┌─────────────────────────────────────────────────────────────┐
│             악의적인 Sequencer                                │
└─────────────────────────────────────────────────────────────┘

Day 0:
  악의적인 Sequencer
  ↓
  1. L1에 정상 배치 데이터 게시
     TX calldata: [Alice: 100 ETH, Bob: 50 ETH]
  ↓
  2. 하지만 잘못된 Output Root 제출
     주장: "Alice: 0 ETH, Bob: 150 ETH"  ← 거짓!
  ↓
  3. L2OutputOracle에 제출
     proposeL2Output(FakeOutputRoot, ...)

Day 1-7 (Challenge Period):
  정직한 Challenger
  ↓
  1. L1에서 배치 데이터 가져오기
     TX calldata: [Alice: 100 ETH, Bob: 50 ETH]
  ↓
  2. 로컬에서 재파생
     Result: "Alice: 100 ETH, Bob: 50 ETH"
  ↓
  3. Output Root 계산
     CorrectOutputRoot = Hash(...)
  ↓
  4. 비교
     FakeOutputRoot ≠ CorrectOutputRoot  ← 불일치!
  ↓
  5. Dispute Game 생성
     create(gameType, rootClaim=FakeOutputRoot, ...)
  ↓
  6. 이분 탐색으로 불일치 지점 찾기
  ↓
  7. 단일 Step 검증
     ├─ Pre-state: Alice 100 ETH
     ├─ Instruction: Transfer 0 ETH
     └─ Post-state: Alice 100 ETH (not 0!)
  ↓
  8. 온체인 검증
     MIPS.step() 실행 → Challenger 승리!

결과:
  ✅ FakeOutputRoot 무효화
  ✅ Sequencer의 bond 몰수
  ✅ CorrectOutputRoot로 교체
  ✅ 시스템 보안 유지
```

### 시나리오 2: Sequencer가 L1에 잘못된 배치 데이터 게시

```
┌─────────────────────────────────────────────────────────────┐
│          L1 데이터 자체가 잘못된 경우                         │
└─────────────────────────────────────────────────────────────┘

Day 0:
  악의적인 Sequencer
  ↓
  1. L1에 잘못된 배치 데이터 게시
     TX calldata: [Bob: +1000000 ETH]  ← 위조 트랜잭션!
  ↓
  2. Output Root 제출
     OutputRoot = Hash(Bob: 1000050 ETH, ...)

문제:
  ❓ L1 데이터 자체가 잘못되었는데 어떻게 검증?

해답:
  ✅ 배치 유효성 규칙 (Batch Validity Rules)

배치 유효성 검사:
  ┌──────────────────────────────────────┐
  │ 1. 트랜잭션 서명 검증                 │
  │    → Bob이 직접 서명한 TX인가?        │
  │    → 서명 검증 실패 → TX 무효         │
  ├──────────────────────────────────────┤
  │ 2. Nonce 검증                        │
  │    → Bob의 올바른 nonce인가?          │
  │    → 아니면 TX 무효                   │
  ├──────────────────────────────────────┤
  │ 3. Gas 검증                          │
  │    → 충분한 gas가 있는가?             │
  │    → 없으면 TX 무효                   │
  ├──────────────────────────────────────┤
  │ 4. 잔액 검증                          │
  │    → Bob이 충분한 잔액이 있는가?      │
  │    → 없으면 TX 무효                   │
  └──────────────────────────────────────┘

재파생 결과:
  Challenger
  ↓
  1. L1에서 배치 데이터 가져오기
     TX calldata: [Bob: +1000000 ETH]
  ↓
  2. 재파생 시도
     ├─ 서명 검증 실패!
     ├─ 또는 잔액 부족!
     └─ TX 무효 처리
  ↓
  3. Output Root 계산
     CorrectOutputRoot = Hash(Bob: 50 ETH, ...)  ← 위조 TX 무시
  ↓
  4. 비교
     SequencerOutputRoot ≠ CorrectOutputRoot
  ↓
  5. Fault Proof 제출 및 승리!

결과:
  ✅ 잘못된 배치도 검증 가능
  ✅ 유효성 규칙으로 자동 필터링
  ✅ 올바른 상태만 반영
```

### 시나리오 3: 모든 Challenger가 공모하여 거짓 증명

```
┌─────────────────────────────────────────────────────────────┐
│         51% 공격: Challenger들이 모두 공모                    │
└─────────────────────────────────────────────────────────────┘

가정:
  - 모든 Challenger가 악의적
  - 모두 거짓 Output Root에 동의

문제:
  ❓ 정직한 Challenger가 없으면?

해답:
  ✅ 단 한 명의 정직한 Challenger만 있으면 충분!

이유:
  1. 온체인 검증은 다수결이 아님
     ├─ MIPS VM이 직접 실행
     ├─ 수학적으로 올바른 쪽이 승리
     └─ 투표나 합의가 아님

  2. 경제적 인센티브
     ├─ 정직한 Challenger: Bond 회수 + 보상
     ├─ 악의적 Challenger: Bond 몰수
     └─ 정직한 행동이 경제적으로 유리

  3. 무허가 참여 (Permissionless)
     ├─ 누구나 Challenger가 될 수 있음
     ├─ 새로운 정직한 참여자 등장 가능
     └─ 영구적인 51% 공격 불가능

결과:
  ✅ 단 1명의 정직한 Challenger로 충분
  ✅ 수학적 정확성이 승리 결정
  ✅ 다수결이 아닌 진실성 기반
```

---

## 결론

### L1 데이터 검증의 핵심

```
┌─────────────────────────────────────────────────────────────┐
│            핵심 원리 정리                                     │
└─────────────────────────────────────────────────────────────┘

1. L1 데이터 자체는 검증 대상이 아님
   ├─ L1에 있으면 "존재"가 보장됨
   └─ 변경 불가능 (Immutable)

2. 검증 대상은 L1 데이터의 "해석"
   ├─ L1 데이터 → L2 상태 (재파생)
   ├─ Sequencer의 Output Root
   └─ 이 둘의 일치 여부

3. 재파생은 결정론적
   ├─ 같은 L1 데이터 → 항상 같은 L2 상태
   ├─ 누구나 독립적으로 검증 가능
   └─ 수학적 정확성 보장

4. 온체인 최종 검증
   ├─ 단일 VM step을 L1에서 실행
   ├─ 다수결이 아닌 계산 결과로 판단
   └─ 정직한 쪽이 반드시 승리
```

### 보안 보장

| 보장 | 메커니즘 |
|------|----------|
| **데이터 가용성** | L1 블록체인 (영구 저장) |
| **실행 정확성** | 결정론적 EVM + MIPS VM |
| **검증 가능성** | 누구나 L1 데이터로 재파생 가능 |
| **경제적 안전성** | Bond 시스템 + 인센티브 |
| **탈중앙화** | 무허가 Challenger 참여 |

### 최종 답변

> **"L1에 제출한 데이터가 정말 맞는 정보라는 것을 어떻게 증명하지?"**

**답변:**

```
1. L1 데이터 자체를 증명하는 것이 아니라,
2. L1 데이터로부터 파생된 L2 상태를 증명합니다.
3. 누구나 L1 데이터를 재파생하여 검증할 수 있습니다.
4. 불일치 발견 시 Fault Proof로 증명합니다.
5. 온체인에서 수학적으로 올바른 쪽이 승리합니다.

핵심: L1은 "데이터 저장소"이고,
      Fault Proof는 "데이터 해석의 정확성"을 검증합니다.
```

---

**작성일**: 2025-10-17
**기반 코드**: tokamak-thanos (commit: ef63c0e65)
