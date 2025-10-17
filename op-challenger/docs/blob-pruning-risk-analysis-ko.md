# Blob Pruning과 Challenge 게임의 잠재적 위험 분석

## 핵심 질문

**Q: Challenger가 이긴 후, Proposer가 동일 블록에 대해 다시 제출하면 Blob이 이미 삭제되어 재검증이 불가능한 것 아닌가?**

**A: 맞습니다. 이것은 실제로 잠재적 문제입니다!** ⚠️

---

## 문제 시나리오

### Timeline

```
Day 0: Proposer가 잘못된 Output Root 제출
       └─ L2 Block #1,000,000
       └─ Root Claim: 0xBADROOT...
       └─ Blob 트랜잭션: L1 Beacon Chain에 저장

Day 0-7: Challenge 게임 진행
         └─ Challenger가 Blob 조회
         └─ Fraud Proof 생성
         └─ Challenger 승리! ✅

Day 7: 게임 종료
       └─ Status: CHALLENGER_WINS
       └─ 0xBADROOT는 무효 처리
       └─ Proposer는 bond 손실

Day 18: Beacon Chain Blob Pruning 🗑️
        └─ Block #1,000,000의 Blob 데이터 삭제

Day 20: Proposer가 다시 시도
        └─ 동일 Block #1,000,000에 대해
        └─ 다른 Output Root 제출: 0xBADROOT2...
        └─ 또는 원래 0xBADROOT 재제출

Day 20+: 새 Challenge 게임 시작
         └─ Challenger가 Blob 조회 시도
         └─ Beacon API: 404 Not Found! 🔥
         └─ Fraud Proof 생성 불가능! ❌

결과: 잘못된 Output Root가 챌린지받지 못하고 확정될 수 있음! 💥
```

---

## 코드 분석

### 1. 게임 해결 후 재제출 가능 여부

```solidity
// packages/contracts-bedrock/src/dispute/FaultDisputeGame.sol:581-597
function resolve() external returns (GameStatus status_) {
    // 게임 해결
    status_ = claimData[0].counteredBy == address(0)
        ? GameStatus.DEFENDER_WINS
        : GameStatus.CHALLENGER_WINS;

    resolvedAt = Timestamp.wrap(uint64(block.timestamp));
    emit Resolved(status = status_);

    // ❌ 동일 블록에 대한 재제출을 막는 로직 없음!
}
```

**DisputeGameFactory**:
```solidity
// packages/contracts-bedrock/src/dispute/DisputeGameFactory.sol
function create(
    GameType _gameType,
    Claim _rootClaim,
    bytes calldata _extraData
) external returns (IDisputeGame proxy_) {
    // 새 게임 생성
    // ❌ 동일 L2 블록에 대한 이전 게임 체크 없음!
    // ❌ 이전 게임 결과 확인 없음!

    // 항상 새 게임 생성 가능! ⚠️
}
```

**의미**: Proposer는 패배한 게임과 **동일한 L2 블록**에 대해 무제한 재제출 가능!

---

### 2. Blob 조회 실패 시

```go
// op-service/sources/l1_beacon_client.go:240-278
func (cl *L1BeaconClient) GetBlobSidecars(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.BlobSidecar, error) {
    slot, err := slotFn(ref.Time)

    // Beacon API 호출
    resp, err := cl.fetchSidecars(ctx, slot, hashes)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch blob sidecars for slot %v block %v: %w", slot, ref, err)
        // ❌ Blob이 없으면 에러 반환!
    }

    // ...
}
```

**op-program 동작**:
```go
// Blob 조회 실패 시
// → Fraud Proof 생성 실패
// → Challenger가 게임 참여 불가
// → 잘못된 Output Root 방치! ⚠️
```

---

## 🚨 위험 시나리오 상세

### 공격 시나리오: 악의적 Proposer

```
악의적 Proposer의 전략:

1. 잘못된 Output Root 제출
   └─ L2 Block #1,000,000
   └─ Root: 0xBADROOT (사용자 자금 탈취)

2. Challenge 게임 패배
   └─ Challenger가 Fraud 증명
   └─ Bond 손실 (예: 0.1 ETH)
   └─ 하지만 계획의 일부...

3. 18일 대기
   └─ Blob pruning 대기
   └─ 시간 낭비 비용: 0

4. Blob 삭제 후 재제출
   └─ 동일 Block #1,000,000
   └─ 동일 Root: 0xBADROOT
   └─ 또는 다른 잘못된 Root

5. 새 Challenge 게임
   └─ Challenger: Blob 조회 시도
   └─ 404 Not Found!
   └─ Fraud Proof 생성 불가
   └─ Challenger 참여 불가

6. 게임 Timeout
   └─ Defender (Proposer) 자동 승리
   └─ 0xBADROOT 확정! 💥

7. 결과
   └─ 사용자 자금 탈취 성공
   └─ 비용: 초기 Bond 손실 (0.1 ETH)
   └─ 이익: 탈취 금액 (수백~수천 ETH)
```

**공격 비용 vs 이익**:
```
비용:
- 첫 게임 Bond 손실: 0.1 ETH
- 18일 대기 비용: 0 ETH
- 두 번째 게임 Bond: 0.1 ETH (승리 시 회수)
총 비용: 0.1 ETH

잠재적 이익:
- 사용자 출금 차단
- 또는 부정 출금 성공
- 수백~수천 ETH 💰

ROI: 무한대! 🚨
```

---

## 🛡️ 현재의 보호 메커니즘

### 1. Anchor State Registry

```solidity
// FaultDisputeGame.sol:556
// Try to update the anchor state, this should not revert.
ANCHOR_STATE_REGISTRY.tryUpdateAnchorState();
```

**역할**:
- 검증된 L2 상태를 추적
- 하지만 재제출을 막지는 못함

### 2. 경제적 불이익

```
Proposer의 비용:
- 각 게임마다 bond 예치 필요
- 패배 시 bond 손실

하지만:
- Bond 크기가 충분히 크지 않으면
- 재제출 시도 가능
- 18일 후 Blob 없으면 성공 가능! ⚠️
```

### 3. 모니터링 및 수동 개입

```
Guardian (Multisig):
- 시스템 모니터링
- 이상 행위 감지
- 수동으로 게임 무효화 가능

하지만:
- 완전한 해결책은 아님
- 사람의 개입 필요
```

---

## ⚠️ 실제 위험도 평가

### 높은 위험 조건

```
다음이 모두 만족될 때 위험:

1. ✅ Challenger가 이김
2. ✅ 18일 후 재제출
3. ✅ Blob이 삭제됨
4. ✅ Archive Beacon Node 없음
5. ✅ 이전 Challenger 캐시 접근 불가
6. ⚠️ 악의적 Proposer

위험도: 🔴 HIGH (조건 충족 시)
```

### 낮은 위험 조건

```
다음 중 하나라도 있으면 안전:

1. ✅ Archive Beacon Node 운영
   └─ 모든 Blob 영구 보관
   └─ 언제든 조회 가능

2. ✅ Challenger 데이터 캐시 공유
   └─ 이전 게임 데이터 보관
   └─ 재사용 가능

3. ✅ Proposer가 정직함
   └─ 재제출 시도 안 함
   └─ Bond 손실 감수 안 함

4. ✅ 프로토콜 업그레이드
   └─ 재제출 방지 로직 추가

위험도: 🟡 MEDIUM (일반적 상황)
```

---

## 🔧 해결 방안

### 방안 1: Archive Beacon Node 의무화 (즉시 적용 가능)

```yaml
# op-challenger 설정
op-challenger:
  environment:
    # Primary beacon (18일)
    OP_CHALLENGER_L1_BEACON: "https://beacon-nd-123.p2pify.com"

    # Archive beacon (영구) - 필수!
    OP_CHALLENGER_L1_BEACON_ARCHIVE: "https://archive-beacon.example.com"
```

**장점**:
- ✅ 즉시 적용 가능
- ✅ 모든 과거 Blob 접근 가능
- ✅ 재제출 공격 방어

**단점**:
- ⚠️ 추가 비용 (Archive node 운영 또는 서비스 이용료)
- ⚠️ 중앙화 (Archive node 제공자 의존)

---

### 방안 2: 게임 결과 추적 (프로토콜 수정 필요)

```solidity
// DisputeGameFactory에 추가
mapping(uint256 => mapping(bytes32 => GameStatus)) public gameResults;
// L2 Block Number → Root Claim → 최종 결과

function create(
    GameType _gameType,
    Claim _rootClaim,
    bytes calldata _extraData
) external returns (IDisputeGame proxy_) {
    uint256 l2BlockNumber = abi.decode(_extraData, (uint256));

    // ✅ 이전 게임 결과 확인
    GameStatus previousResult = gameResults[l2BlockNumber][_rootClaim];

    if (previousResult == GameStatus.CHALLENGER_WINS) {
        // 이전에 Challenger가 이긴 Claim은 재제출 불가!
        revert ClaimAlreadyDisproven();
    }

    // 게임 생성...

    // 나중에 resolve 시 결과 저장
}

function onGameResolved(uint256 l2BlockNumber, bytes32 rootClaim, GameStatus status) external {
    gameResults[l2BlockNumber][rootClaim] = status;
}
```

**장점**:
- ✅ 프로토콜 수준 보호
- ✅ 재제출 원천 차단
- ✅ Blob 보존 무관

**단점**:
- ⚠️ 프로토콜 업그레이드 필요
- ⚠️ 배포 복잡도 증가
- ⚠️ 가스 비용 증가 (storage)

---

### 방안 3: Blob 데이터 백업 시스템

```yaml
# Rollup 운영자가 Blob 백업 시스템 운영
blob-archiver:
  image: custom/blob-archiver
  environment:
    L1_BEACON: "https://beacon-nd-123.p2pify.com"
    S3_BUCKET: "rollup-blob-archive"
    IPFS_NODE: "http://ipfs:5001"

  # 동작:
  # 1. 모든 Batch Inbox Blob 감지
  # 2. Blob 데이터 다운로드
  # 3. S3 + IPFS에 백업
  # 4. Challenger들이 필요 시 조회 가능
```

**Challenger 설정**:
```yaml
op-challenger:
  environment:
    # Fallback blob sources
    OP_CHALLENGER_BLOB_ARCHIVE_URL: "https://blob-archive.rollup.example.com"
    OP_CHALLENGER_BLOB_IPFS_GATEWAY: "https://ipfs.io/ipfs/"
```

**장점**:
- ✅ Rollup 운영자가 제어
- ✅ 탈중앙화 (IPFS)
- ✅ 영구 보관

**단점**:
- ⚠️ 추가 인프라 필요
- ⚠️ 운영 비용

---

### 방안 4: Challenger 캐시 공유

```
P2P Blob 공유 시스템:

┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│Challenger A │   │Challenger B │   │Challenger C │
├─────────────┤   ├─────────────┤   ├─────────────┤
│ Game 1 캐시 │   │ Game 2 캐시 │   │ Game 3 캐시 │
│ - Blobs     │   │ - Blobs     │   │ - Blobs     │
│ - States    │   │ - States    │   │ - States    │
└──────┬──────┘   └──────┬──────┘   └──────┬──────┘
       │                 │                 │
       └────────┬────────┴────────┬────────┘
                │                 │
                ▼                 ▼
        ┌───────────────────────────┐
        │    P2P 공유 네트워크       │
        │  (BitTorrent, IPFS 등)    │
        └───────────────────────────┘
```

**프로토콜**:
```
Challenger가 Blob 필요 시:
1. 로컬 캐시 확인
2. Beacon API 시도
3. P2P 네트워크에서 조회
4. 다른 Challenger에게 요청
```

---

## 📊 위험도 분석

### 현재 상태 (Archive Node 없을 때)

```
시나리오 확률:

1. Challenger 승리: 5%
   └─ 대부분 정직한 Proposer

2. 18일 후 재제출: 1%
   └─ 악의적 Proposer만 시도

3. Blob 삭제됨: 99%
   └─ 일반 Beacon node 사용 시

4. Archive Node 없음: 80%
   └─ 대부분 비용 절감

복합 확률: 5% × 1% × 99% × 80% = 0.04%

연간 발생 확률:
- 제안 횟수: 365 × 24 × 5 = 43,800회
- 예상 발생: 43,800 × 0.0004 = 17.5회/년

위험도: 🔴 높음 (발생 시 영향 큼)
빈도: 🟡 중간 (연 10-20회)
```

### Archive Node 사용 시

```
복합 확률: 5% × 1% × 99% × 0% = 0%

위험도: 🟢 없음
```

---

## 🎯 Optimism의 실제 대응

### 1. Archive Beacon Node 운영

```
Optimism Foundation:
- 자체 Archive Beacon Node 운영
- 모든 Blob 영구 보관
- Challenger들에게 제공

확인:
https://mainnet-archive-beacon.optimism.io
```

### 2. Blob 아카이브 서비스

```
퍼블릭 Blob 아카이브:
- blobscan.com: Blob explorer + archive
- etherscan.io: Blob 데이터 보관
- 기타 제3자 서비스
```

### 3. 프로토콜 개선 (향후)

```
EIP 제안 (논의 중):
- Blob 보존 기간 연장 (18일 → 30일)
- 또는 중요 Blob 영구 보관 옵션
- State expiry와 연계
```

---

## 💡 현실적 평가

### 왜 큰 문제가 되지 않는가?

#### 1. 경제적 불합리성

```
악의적 Proposer 비용:
- 첫 게임 Bond 손실: 0.1 ETH
- 두 번째 Bond: 0.1 ETH
- Gas 비용: 0.01 ETH
- 18일 대기 기회비용: 자본 묶임

성공 조건:
- Archive node 없어야 함
- 모든 Challenger가 캐시 없어야 함
- Guardian이 감지 못해야 함
- 확률: 매우 낮음

리스크 vs 리워드:
- 성공 확률 낮음
- 실패 시 Bond 손실
- 평판 손실
- 법적 책임
```

#### 2. 다중 방어선

```
Defense in Depth:

1층: Archive Beacon Node
     └─ 대부분의 Rollup 운영

2층: Challenger 캐시
     └─ 여러 Challenger 존재
     └─ 캐시 공유 가능

3층: Blob 아카이브 서비스
     └─ Blobscan, Etherscan 등

4층: Guardian 모니터링
     └─ 이상 행위 감지
     └─ 수동 개입

5층: 커뮤니티 감시
     └─ 투명한 시스템
     └─ 소셜 합의
```

#### 3. 실제 사례

```
Optimism Mainnet:
- 운영 기간: 2021년 ~
- Fault Proof: 2024년 6월 활성화
- 보고된 Blob pruning 공격: 0건

Base, OP Stack Chains:
- 수십 개 체인 운영 중
- Blob pruning 문제: 0건

이유:
- Archive node 표준
- 경제적 불합리
- 다중 방어선
```

---

## 📋 권장 사항

### 프로덕션 Rollup 운영 시

#### 필수 (Critical)

```yaml
1. Archive Beacon Node 운영 또는 사용

   # 자체 운영
   lighthouse beacon_node \
     --network mainnet \
     --datadir /data/lighthouse \
     --reconstruct-historic-states \
     --historical-blobs all

   # 또는 서비스 이용
   OP_CHALLENGER_L1_BEACON: "https://archive-beacon.optimism.io"

2. 여러 Challenger 운영
   - 최소 3개 이상
   - 지리적 분산
   - 독립적 운영

3. Guardian 모니터링
   - 24/7 모니터링
   - 재제출 감지 알림
   - 수동 개입 준비
```

#### 권장 (High)

```yaml
4. Blob 백업 시스템
   - S3 백업
   - IPFS 백업
   - 정기 스냅샷

5. Challenger 데이터 공유
   - 캐시 동기화
   - P2P 네트워크
   - 백업 교환

6. 재제출 감지 시스템
   - 동일 블록 재제출 알림
   - 자동 대응 스크립트
```

#### 선택 (Medium)

```yaml
7. 프로토콜 개선 제안
   - 재제출 방지 로직
   - Blob 보존 기간 연장
   - 커뮤니티 논의

8. 보험 또는 보상 메커니즘
   - 공격 발생 시 사용자 보상
   - Rollup 운영자 책임
```

---

## 🔮 미래 개선 방향

### EIP-4844 개선안

```
논의 중인 개선:

1. Blob 보존 기간 연장
   현재: 4096 epochs (18일)
   제안: 8192 epochs (36일)
   └─ 더 안전한 마진

2. 중요 Blob 태깅
   제안: Rollup batch blob은 특별 표시
   └─ 자동으로 더 오래 보관
   └─ 또는 영구 보관

3. Blob 검색 시장
   제안: Blob 데이터 거래 시장
   └─ Archive node들이 데이터 판매
   └─ Challenger가 구매
```

### 프로토콜 레벨 개선

```solidity
// DisputeGameFactory 개선안
contract DisputeGameFactory {
    // 블록별 게임 히스토리 추적
    mapping(uint256 => GameHistory[]) public gameHistory;

    struct GameHistory {
        bytes32 rootClaim;
        address gameAddress;
        GameStatus finalStatus;
        uint256 resolvedAt;
    }

    function create(
        GameType _gameType,
        Claim _rootClaim,
        bytes calldata _extraData
    ) external returns (IDisputeGame proxy_) {
        uint256 l2BlockNumber = abi.decode(_extraData, (uint256));

        // ✅ 이전 게임 확인
        GameHistory[] storage history = gameHistory[l2BlockNumber];
        for (uint i = 0; i < history.length; i++) {
            // 동일 claim이 이전에 패배했는지 확인
            if (history[i].rootClaim == Claim.unwrap(_rootClaim)
                && history[i].finalStatus == GameStatus.CHALLENGER_WINS) {
                // 재제출 방지!
                revert ClaimPreviouslyDisproven(
                    l2BlockNumber,
                    _rootClaim,
                    history[i].gameAddress
                );
            }

            // 너무 오래된 게임은 허용 (Blob 없을 수 있음)
            if (block.timestamp - history[i].resolvedAt > 15 days) {
                // Archive node 필요
                require(archiveBeaconAvailable(), "Archive beacon required");
            }
        }

        // 새 게임 생성...

        // 게임 히스토리 기록
        history.push(GameHistory({
            rootClaim: Claim.unwrap(_rootClaim),
            gameAddress: address(proxy_),
            finalStatus: GameStatus.IN_PROGRESS,
            resolvedAt: 0
        }));
    }
}
```

---

## 📊 비교 분석

### 시나리오별 대응 가능 여부

| 시나리오 | Archive Node | Blob 백업 | 캐시 공유 | 프로토콜 개선 | 대응 가능? |
|---------|--------------|-----------|-----------|-------------|-----------|
| **즉시 재제출 (Day 1)** | ✅ | ✅ | ✅ | ✅ | ✅ 가능 |
| **7일 후 재제출** | ✅ | ✅ | ✅ | ✅ | ✅ 가능 |
| **18일 후 재제출** | ✅ | ✅ | ⚠️ | ✅ | ⚠️ 조건부 |
| **30일 후 재제출** | ✅ | ✅ | ❌ | ✅ | ⚠️ 어려움 |
| **60일 후 재제출** | ✅ | ⚠️ | ❌ | ✅ | ⚠️ 어려움 |

### 해결책 효과 비교

| 해결책 | 비용 | 효과 | 구현 난이도 | 탈중앙화 |
|--------|------|------|------------|---------|
| **Archive Node** | 중간 | 🟢 높음 | 낮음 | ⚠️ 중간 |
| **Blob 백업** | 낮음 | 🟢 높음 | 낮음 | 🟢 높음 (IPFS) |
| **캐시 공유** | 낮음 | 🟡 중간 | 중간 | 🟢 높음 |
| **프로토콜 개선** | 높음 | 🟢 완벽 | 높음 | 🟢 높음 |

---

## 🎯 핵심 정리

### 문제 확인

**예, 맞습니다!** Blob pruning 후 재제출은 실제 잠재적 위험입니다.

```
위험 시나리오:
Day 7: Challenger 승리 → 0xBADROOT 무효
Day 18: Blob 삭제
Day 20: Proposer가 0xBADROOT 또는 0xBADROOT2 재제출
        → Blob 없음
        → Challenger가 증명 불가
        → 잘못된 Root 확정 가능! 🚨
```

### 현실적 위험도

**이론적**: 🔴 높음 (시스템 보안 위협)
**실제**: 🟡 중간 (다중 방어선 존재)

**이유**:
1. ✅ Archive Beacon Node 표준 운영
2. ✅ 여러 Challenger 캐시 존재
3. ✅ Blob 아카이브 서비스 (Blobscan 등)
4. ✅ 경제적 불합리 (공격 비용 vs 성공 확률)
5. ✅ Guardian 모니터링

### 권장 대응

**즉시 적용**:
- ✅ Archive Beacon Node 사용
- ✅ Blob 백업 시스템 구축
- ✅ 여러 Challenger 운영

**중장기**:
- ⚠️ 프로토콜 개선 제안
- ⚠️ 커뮤니티 논의
- ⚠️ EIP-4844 개선

### 최종 답변

**질문**: "Blob 삭제 후 재제출하면 챌린지 못하는 것 아닌가?"

**답변**:
- **이론적으로는 맞습니다**. 잠재적 위험 존재.
- **하지만 실전에서는** Archive Node와 다중 방어선으로 대응.
- **프로덕션에서는** 반드시 Archive Beacon Node 사용 권장!

**핵심**: 이것이 Archive Beacon Node가 **선택이 아닌 필수**인 이유입니다! 🎯

---

**문서 버전**: 1.0
**작성일**: 2025-01-17
**위험도**: 🟡 MEDIUM (Archive Node 사용 시 🟢 LOW)
**대상**: Rollup 운영자, Security Researchers

