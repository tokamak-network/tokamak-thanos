# L2 시스템 아키텍처 가이드

## 목차
1. [개요](#개요)
2. [L2 시스템의 구성 요소](#l2-시스템의-구성-요소)
3. [op-node의 두 가지 모드](#op-node의-두-가지-모드)
4. [Engine API 통신](#engine-api-통신)
5. [다음 단계](#다음-단계)

---

## 개요

Optimism Stack (Tokamak Thanos)에서 **L2는 단일 프로세스가 아닙니다**. L2 시스템은 **op-node**와 **op-geth** 두 개의 프로세스가 협력하여 작동합니다.

```
┌─────────────────────────────────────────────┐
│                L2 System                    │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────────────┐   ┌──────────────┐   │
│  │    op-node       │◄─►│   op-geth    │   │
│  │  (Rollup Logic)  │   │(EVM Execute) │   │
│  └──────────────────┘   └──────────────┘   │
│         ▲                       ▲           │
│         │                       │           │
│         │                       │           │
└─────────┼───────────────────────┼───────────┘
          │                       │
          ▼                       ▼
         L1                    RPC Calls
     (Ethereum)              (사용자 요청)
```

---

## L2 시스템의 구성 요소

### 1. op-geth: EVM 실행 엔진

```
op-geth (Ethereum 실행 클라이언트)
├─ 역할: 트랜잭션 실행 및 상태 관리
├─ 기반: go-ethereum (Geth) 포크
│
├─ 담당 기능:
│  ├─ EVM 실행 (트랜잭션 처리)
│  ├─ State DB 관리 (계정, 잔액, 스토리지)
│  ├─ Mempool 관리
│  ├─ Block 실행
│  └─ JSON-RPC 제공 (eth_*, debug_* 등)
│
└─ 특징: Ethereum과 거의 동일, Deposit TX 지원 추가
```

**op-geth만으로는 불완전:**
- ❌ 어떤 트랜잭션을 언제 실행할지 모름
- ❌ L1에서 데이터를 가져올 수 없음
- ❌ Rollup 로직이 없음

### 2. op-node: Rollup 조정자

```
op-node (Rollup 노드)
├─ 역할: L1↔L2 연결 및 Rollup 로직
├─ 기반: 순수 Optimism 코드
│
├─ 담당 기능:
│  ├─ L1 모니터링 (Batch Inbox 감시)
│  ├─ Batch 데이터 파싱
│  ├─ L2 Block 구성 결정
│  ├─ op-geth에 실행 지시
│  ├─ Derivation Pipeline (L1→L2 변환)
│  └─ Sequencer 역할 (Sequencer 모드 시)
│
└─ 특징: Rollup 프로토콜의 두뇌
```

**op-node만으로는 불완전:**
- ❌ 트랜잭션을 실행할 수 없음
- ❌ State를 관리할 수 없음
- ❌ EVM이 없음

### 3. 둘의 관계

```
op-node (운전자)                op-geth (엔진)
├─ 언제 블록을 만들지 결정       ├─ 실제 트랜잭션 실행
├─ 어떤 트랜잭션을 포함할지      ├─ State 업데이트
├─ L1 데이터 읽기                ├─ Receipt 생성
└─ op-geth에게 명령 전달         └─ op-node의 명령 수행

        │                              ▲
        └──────── Engine API ──────────┘
           (JSON-RPC over HTTP)
```

---

## op-node의 두 가지 모드

op-node는 설정에 따라 **Sequencer 모드** 또는 **Follower 모드**로 동작합니다.

### Sequencer 모드 (생산자)

```
┌─────────────────────────────────────────────┐
│        op-node (Sequencer 모드)             │
├─────────────────────────────────────────────┤
│                                             │
│ 역할: 새로운 L2 블록을 "생성"               │
│                                             │
│ ┌─────────────────────────────────┐         │
│ │ 1. 사용자 트랜잭션 수집          │         │
│ │    (op-geth mempool에서)        │         │
│ └────────────┬────────────────────┘         │
│              ▼                              │
│ ┌─────────────────────────────────┐         │
│ │ 2. 블록 생성 (매 2초)            │         │
│ │    - PayloadAttributes 구성     │         │
│ │    - op-geth에 실행 지시        │         │
│ └────────────┬────────────────────┘         │
│              ▼                              │
│ ┌─────────────────────────────────┐         │
│ │ 3. Unsafe L2 블록 생성           │         │
│ │    - 즉시 사용자에게 응답        │         │
│ └────────────┬────────────────────┘         │
│              ▼                              │
│ ┌─────────────────────────────────┐         │
│ │ 4. op-batcher에게 전달           │         │
│ │    - L1 제출용                   │         │
│ └─────────────────────────────────┘         │
│                                             │
│ 특징:                                       │
│ ✅ 트랜잭션 순서 결정 (Sequencing)          │
│ ✅ 블록 타임스탬프 결정                      │
│ ✅ 즉각적인 L2 확정성                        │
│ ⚠️  단일 Sequencer (중앙화 포인트)          │
│                                             │
└─────────────────────────────────────────────┘
```

**설정 방법:**
```bash
op-node \
  --sequencer.enabled=true              # ⭐ Sequencer 활성화
  --sequencer.l1-confs=4                # L1 확인 블록 수
  --p2p.sequencer.key=<private-key>     # P2P용 키
```

### Follower 모드 (검증자/복제자)

```
┌─────────────────────────────────────────────┐
│        op-node (Follower 모드)              │
├─────────────────────────────────────────────┤
│                                             │
│ 역할: L1에서 L2 블록을 "복원/검증"          │
│                                             │
│ ┌─────────────────────────────────┐         │
│ │ 1. L1 모니터링                   │         │
│ │    - Batch Inbox 감시            │         │
│ │    - 새 배치 데이터 감지         │         │
│ └────────────┬────────────────────┘         │
│              ▼                              │
│ ┌─────────────────────────────────┐         │
│ │ 2. 배치 파싱                     │         │
│ │    - 압축 해제                   │         │
│ │    - 트랜잭션 추출               │         │
│ │    - Batcher 검증                │         │
│ └────────────┬────────────────────┘         │
│              ▼                              │
│ ┌─────────────────────────────────┐         │
│ │ 3. L2 블록 재실행                │         │
│ │    - op-geth에 실행 지시         │         │
│ │    - State Root 검증             │         │
│ └────────────┬────────────────────┘         │
│              ▼                              │
│ ┌─────────────────────────────────┐         │
│ │ 4. Safe/Finalized Head 업데이트  │         │
│ │    - L1 확정성 따름               │         │
│ └─────────────────────────────────┘         │
│                                             │
│ 특징:                                       │
│ ✅ 독립적으로 L2 상태 검증                   │
│ ✅ Sequencer 신뢰 불필요                     │
│ ✅ 누구나 실행 가능 (탈중앙화)               │
│ ⏱️  L1 확정 시간만큼 지연 (~12분)            │
│                                             │
└─────────────────────────────────────────────┘
```

**설정 방법:**
```bash
op-node \
  --rollup.config=./rollup.json
  --l1=https://eth.llamarpc.com
  --l2=http://localhost:8551
  # ⭐ --sequencer.enabled 없음! (기본값: false = Follower)
```

### 모드 비교표

| 항목 | Sequencer 모드 | Follower 모드 |
|------|---------------|--------------|
| **역할** | 블록 생성 (Producer) | 블록 검증 (Verifier) |
| **트랜잭션 수신** | ✅ 사용자로부터 직접 | ❌ (Sequencer에게 전달) |
| **블록 생성** | ✅ 능동적 생성 (2초마다) | ❌ L1에서 복원 |
| **L1 읽기** | ⚠️  Deposit 트랜잭션만 | ✅ 모든 배치 데이터 |
| **L1 쓰기** | ✅ (op-batcher 통해) | ❌ |
| **지연시간** | ~1초 (Unsafe) | ~12분 (Safe) |
| **신뢰 필요** | 사용자는 Sequencer 신뢰 | 신뢰 불필요 (L1 검증) |
| **개수** | 1개 (현재) | 무제한 |
| **운영 주체** | Rollup 운영자 | 누구나 |

---

## Engine API 통신

op-node와 op-geth는 **Engine API**를 통해 통신합니다.

### Engine API 인터페이스

```go
// op-node/rollup/derive/engine_controller.go

type ExecEngine interface {
    // op-geth에게 명령:

    GetPayload(...)       // "빌드한 블록 줘"
    ForkchoiceUpdate(...) // "이 블록이 헤드야"
    NewPayload(...)       // "이 블록 실행해"
    L2BlockRefByLabel(...) // "현재 블록 정보 줘"
}
```

### Sequencer 모드의 블록 생성 흐름

```
사용자 → RPC 요청 → op-geth (Mempool에 저장)
                           ↓
                  op-node가 "블록 만들어" 명령
                           ↓
┌─────────────────────────────────────────────┐
│ 1. op-node (Sequencer 모드)                 │
├─────────────────────────────────────────────┤
│ // 2초마다 블록 생성 결정                    │
│ attrs := PayloadAttributes{                 │
│   timestamp: now,                           │
│   transactions: [tx1, tx2, tx3]             │
│ }                                           │
│                                             │
│ // op-geth에게 명령 ⭐                       │
│ payloadID := engine.ForkchoiceUpdate(       │
│   parent: currentHead,                      │
│   attrs: attrs                              │
│ )                                           │
└──────────────┬──────────────────────────────┘
               │ HTTP JSON-RPC
               │ POST http://op-geth:8551
               ▼
┌─────────────────────────────────────────────┐
│ 2. op-geth (EVM 실행 엔진)                   │
├─────────────────────────────────────────────┤
│ // Mempool에서 트랜잭션 가져오기              │
│ txs := mempool.Get()                        │
│                                             │
│ // EVM으로 실행 ⭐                           │
│ for tx in txs:                              │
│     receipt := EVM.Execute(tx)              │
│     state.Update(receipt)                   │
│                                             │
│ // 블록 생성                                 │
│ block := BuildBlock(txs, receipts)          │
│ payloadID := store(block)                   │
│                                             │
│ return payloadID                            │
└──────────────┬──────────────────────────────┘
               │ 결과 반환
               ▼
┌─────────────────────────────────────────────┐
│ 3. op-node                                  │
│ // 완성된 블록 가져오기                      │
│ payload := engine.GetPayload(payloadID)     │
│                                             │
│ // 최종 확정                                 │
│ status := engine.NewPayload(payload)        │
│ unsafeHead = payload.blockHash              │
│                                             │
│ // op-batcher에게 전달 (L1 제출용)           │
└─────────────────────────────────────────────┘
```

### Follower 모드의 블록 복원 흐름

```
┌─────────────────────────────────────────────┐
│ 1. op-node (Follower 모드)                  │
├─────────────────────────────────────────────┤
│ // L1 모니터링                               │
│ batch := monitorBatchInbox(l1)              │
│                                             │
│ // 배치 파싱                                 │
│ payload := parseBatch(batch)                │
│ // {                                        │
│ //   transactions: [tx1, tx2, tx3],        │
│ //   timestamp: 1234567890,                │
│ //   ...                                   │
│ // }                                        │
│                                             │
│ // op-geth에게 "이거 실행해" 명령 ⭐         │
│ status := engine.NewPayload(payload)        │
└──────────────┬──────────────────────────────┘
               │ HTTP JSON-RPC
               │ POST http://op-geth:8551
               ▼
┌─────────────────────────────────────────────┐
│ 2. op-geth (EVM 실행 엔진)                   │
├─────────────────────────────────────────────┤
│ // 받은 payload 실행 ⭐                      │
│ for tx in payload.transactions:             │
│     receipt := EVM.Execute(tx)              │
│     state.Update(receipt)                   │
│                                             │
│ // State Root 계산                          │
│ stateRoot := state.Root()                   │
│                                             │
│ // 검증 결과 반환                            │
│ return {                                    │
│   status: VALID,                           │
│   latestValidHash: blockHash               │
│ }                                           │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────┐
│ 3. op-node                                  │
│ if status == VALID:                         │
│     safeHead = payload.blockHash ✅         │
│ else:                                       │
│     reorg() ❌                              │
└─────────────────────────────────────────────┘
```

---

## 다음 단계

이제 L2 시스템의 기본 아키텍처를 이해하셨습니다!

실제 시스템 배포 및 운영에 대해서는 다음 문서를 참고하세요:

### 📚 관련 문서

- **[L2 시스템 배포 가이드](./l2-system-deployment-ko.md)** ⭐
  - 완전한 Sequencer + Challenger 스택 구성
  - docker-compose.yml 예시
  - 리소스 요구사항 및 환경별 설정

- **[Batch Inbox Address 상세 분석](./research/batch-inbox-address-ko.md)**
  - L1에서 L2 데이터가 어떻게 전달되는지
  - Batch Inbox의 구조와 보안

### 💡 핵심 요약

```
L2 시스템 = op-node + op-geth

op-node의 두 모드:
├─ Sequencer: 블록 생성 (--sequencer.enabled=true)
└─ Follower:  블록 검증 (기본값)

통신:
└─ Engine API (JSON-RPC over HTTP)
```

### 📖 코드 위치

- `op-node/node/config.go`: op-node 설정
- `op-node/rollup/driver/state.go`: Driver 이벤트 루프
- `op-node/rollup/driver/sequencer.go`: Sequencer 로직
- `op-node/rollup/derive/engine_controller.go`: Engine API 인터페이스

### 🔗 Optimism 스펙

- [Rollup Node Specification](https://specs.optimism.io/protocol/rollup-node.html)
- [Derivation Pipeline](https://specs.optimism.io/protocol/derivation.html)
- [Engine API](https://github.com/ethereum/execution-apis/tree/main/src/engine)
