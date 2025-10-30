# op-challenger 아키텍처 분석 문서

## 목차
1. [개요](#개요)
2. [핵심 개념](#핵심-개념)
3. [시스템 아키텍처](#시스템-아키텍처)
4. [주요 컴포넌트](#주요-컴포넌트)
5. [Dispute Game 메커니즘](#dispute-game-메커니즘)
6. [Trace Provider 시스템](#trace-provider-시스템)
7. [실행 흐름](#실행-흐름)
8. [설정 및 배포](#설정-및-배포)

---

## 개요

`op-challenger`는 Optimism 스택의 Fault Proof 시스템에서 **무효한 상태 주장(invalid claims)을 감지하고 이의를 제기**하는 모듈식 챌린저 에이전트입니다. Go 언어로 작성되었으며, 다양한 유형의 dispute game(증명 게임)을 지원합니다.

### 주요 역할
- L2 output root의 유효성 검증
- 무효한 클레임에 대한 자동 대응
- 분쟁 게임의 진행 상태 모니터링
- 보상금(bond) 청구 및 회수

### 지원하는 Trace Types
- **Cannon**: MIPS 기반 VM을 사용한 fault proof
- **Asterisc**: RISC-V 기반 VM을 사용한 fault proof
- **Alphabet**: 테스트용 간소화된 trace provider
- **Permissioned**: 권한 기반 게임
- **Fast**: 고속 trace 생성 모드

---

## 핵심 개념

### Fault Dispute Game
분쟁 해결을 위한 이진 검색 기반 게임으로, 챌린저와 제안자가 특정 상태 전환의 정확성을 놓고 겨루는 시스템입니다.

**게임 트리 구조:**
```
Root Claim (L2 Output Root)
├─ Attack (disagree)
│  ├─ Defend (agree with parent)
│  └─ Attack (disagree with parent)
└─ Defend (agree)
   └─ ...
```

### Position과 Depth
- **Position**: 게임 트리에서의 위치를 나타내는 gindex (generalized index)
- **Depth**: 트리의 깊이 (0 = root, maxDepth = leaf)
- **MaxDepth**: 게임의 최대 깊이 (보통 73)

### Clock과 타이밍
- 각 클레임은 Chess Clock 방식의 타이머를 가짐
- 시간 내에 대응하지 않으면 상대방 승리
- `maxClockDuration`: 게임당 최대 허용 시간

---

## 시스템 아키텍처

### 전체 구조

```
┌─────────────────────────────────────────────────────────────┐
│                        op-challenger                        │
├─────────────────────────────────────────────────────────────┤
│                          Service                            │
│  ┌────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │  Monitor   │  │  Scheduler   │  │  BondClaimer     │   │
│  └─────┬──────┘  └──────┬───────┘  └────────┬─────────┘   │
│        │                │                    │              │
│  ┌─────▼────────────────▼────────────────────▼─────────┐   │
│  │            GameTypeRegistry                         │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐         │   │
│  │  │ Cannon   │  │ Asterisc │  │ Alphabet │  ...    │   │
│  │  │ Player   │  │ Player   │  │ Player   │         │   │
│  │  └──────────┘  └──────────┘  └──────────┘         │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Fault Game Components                   │   │
│  │  ┌────────┐  ┌────────┐  ┌──────────┐  ┌─────────┐ │   │
│  │  │ Agent  │  │ Solver │  │Responder │  │Contract │ │   │
│  │  └────────┘  └────────┘  └──────────┘  └─────────┘ │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
         │                    │                    │
    ┌────▼────┐         ┌─────▼──────┐      ┌─────▼──────┐
    │ L1 Node │         │ L2 Node    │      │ Contracts  │
    └─────────┘         └────────────┘      └────────────┘
```

### 디렉토리 구조

```
op-challenger/
├── cmd/                    # CLI 명령어 (create-game, move, resolve 등)
├── config/                 # 설정 관리
├── game/                   # 게임 관련 핵심 로직
│   ├── scheduler/          # 게임 스케줄링 및 작업 관리
│   ├── registry/           # 게임 타입 레지스트리
│   ├── fault/              # Fault dispute game 구현
│   │   ├── agent.go        # 게임 진행 에이전트
│   │   ├── player.go       # 게임 플레이어
│   │   ├── solver/         # 게임 솔버 (전략 계산)
│   │   ├── trace/          # Trace provider들
│   │   │   ├── cannon/     # Cannon VM trace
│   │   │   ├── asterisc/   # Asterisc VM trace
│   │   │   ├── alphabet/   # 테스트용 trace
│   │   │   └── outputs/    # L2 output 기반 trace
│   │   ├── contracts/      # 스마트 컨트랙트 바인딩
│   │   ├── responder/      # 트랜잭션 실행
│   │   └── claims/         # Bond 청구 로직
│   ├── keccak/             # Large preimage (Keccak) 처리
│   └── monitor.go          # L1 블록 모니터링
├── sender/                 # 트랜잭션 전송 관리
└── metrics/                # 메트릭 수집
```

---

## 주요 컴포넌트

### 1. Service (`game/service.go`)

전체 챌린저 서비스의 진입점이자 라이프사이클 관리자입니다.

**주요 책임:**
- L1/L2 클라이언트 연결 관리
- 게임 타입 레지스트리 초기화
- 모니터, 스케줄러, 본드 클레이머 조정
- 메트릭 및 프로파일링 서비스 관리

**초기화 과정:**
```go
// service.go:88-129
func (s *Service) initFromConfig(ctx context.Context, cfg *config.Config) error {
    1. TxManager 초기화 (트랜잭션 관리)
    2. L1 클라이언트 연결
    3. Rollup 클라이언트 연결
    4. Factory Contract 바인딩
    5. 게임 타입 등록 (RegisterGameTypes)
    6. Bond 청구 시스템 초기화
    7. Scheduler 초기화
    8. Large Preimage 스케줄러 초기화
    9. Monitor 초기화
}
```

**실행 흐름:**
```go
// service.go:257-266
func (s *Service) Start(ctx context.Context) error {
    s.sched.Start(ctx)          // 스케줄러 시작
    s.claimer.Start(ctx)        // 본드 청구자 시작
    s.preimages.Start(ctx)      // Preimage 검증 시작
    s.monitor.StartMonitoring() // L1 블록 모니터링 시작
}
```

### 2. GameMonitor (`game/monitor.go`)

L1 체인을 모니터링하고 새로운 게임 및 블록을 감지합니다.

**핵심 기능:**

```go
// monitor.go:136-144
func (m *gameMonitor) onNewL1Head(ctx context.Context, sig eth.L1BlockRef) {
    1. L1 시계 업데이트
    2. 진행 중인 게임들 스케줄링 (progressGames)
    3. Large preimage 검증 예약
}
```

**게임 필터링:**
```go
// monitor.go:111-134
func (m *gameMonitor) progressGames(...) error {
    // GameWindow 기준으로 최근 게임만 필터링
    minGameTimestamp := clock.MinCheckedTimestamp(m.clock, m.gameWindow)
    games := m.source.GetGamesAtOrAfter(ctx, blockHash, minGameTimestamp)

    // Allowlist 확인
    for _, game := range games {
        if !m.allowedGame(game.Proxy) {
            continue
        }
        gamesToPlay = append(gamesToPlay, game)
    }

    // 스케줄러에 전달
    m.scheduler.Schedule(gamesToPlay, blockNumber)
}
```

### 3. Scheduler (`game/scheduler/`)

여러 게임을 동시에 처리하기 위한 작업 스케줄링 시스템입니다.

**구조:**
```
Scheduler
├── Coordinator    # 게임 상태 관리 및 작업 스케줄링
├── Worker Pool    # 병렬 작업 실행
└── DiskManager    # 게임 데이터 디렉토리 관리
```

**Coordinator 역할:**

```go
// coordinator.go:60-115
func (c *coordinator) schedule(ctx context.Context, games []GameMetadata, blockNumber uint64) {
    1. 더 이상 필요 없는 게임 상태 제거
    2. 각 게임에 대한 job 생성
       - 이미 진행 중(inflight)인 게임은 스킵
       - 해결된 게임은 스킵
       - 처음 보는 게임은 player 생성
    3. Job을 작업 큐에 enqueue
    4. 완료된 결과를 result 큐에서 처리
}
```

**게임 상태 관리:**
```go
type gameState struct {
    player                GamePlayer
    inflight              bool              // 현재 작업 중인지
    lastProcessedBlockNum uint64            // 마지막 처리 블록
    status                types.GameStatus  // 게임 상태
}
```

### 4. GamePlayer (`game/fault/player.go`)

개별 게임을 진행하는 플레이어입니다.

**생성 과정:**

```go
// player.go:69-159
func NewGamePlayer(...) (*GamePlayer, error) {
    1. 게임 상태 확인 (이미 해결됨?)
    2. 게임 파라미터 로드
       - maxClockDuration
       - gameDepth
    3. TraceAccessor 생성 (Cannon/Asterisc 등)
    4. Oracle contract 로드
    5. L1 head 확인
    6. Preimage uploader 초기화
    7. Responder 생성
    8. Agent 생성 및 연결
}
```

**게임 진행:**

```go
// player.go:174-199
func (g *GamePlayer) ProgressGame(ctx context.Context) GameStatus {
    1. 게임이 이미 완료되었는지 확인
    2. 로컬 노드 동기화 검증
    3. Agent의 Act() 호출 (다음 액션 수행)
    4. 게임 상태 업데이트
    5. 로그 기록
}
```

### 5. Agent (`game/fault/agent.go`)

게임의 실제 로직을 실행하는 핵심 에이전트입니다.

**Act 메소드 흐름:**

```go
// agent.go:79-113
func (a *Agent) Act(ctx context.Context) error {
    1. 게임 해결 시도 (tryResolve)
       - 해결 가능하면 즉시 종료

    2. L2 블록 번호 챌린지 확인
       - 이미 챌린지된 경우 스킵

    3. 컨트랙트에서 게임 상태 로드
       - 모든 클레임 가져오기

    4. Solver에게 다음 액션 계산 요청
       - CalculateNextActions()

    5. 모든 액션을 병렬로 실행
       - Move (attack/defend)
       - Step (증명 제출)
       - ChallengeL2BlockNumber
}
```

**클레임 해결:**

```go
// agent.go:171-213
func (a *Agent) tryResolveClaims(ctx context.Context) error {
    1. 모든 클레임 로드
    2. 각 클레임에 대해:
       - Clock 만료 확인
       - Selective 모드: 인센티브 확인 (individual claims만 해당)
       - CallResolveClaim으로 해결 가능 여부 확인
    3. 해결 가능한 클레임들을 일괄 해결
}
```

**⚠️ Selective Claim Resolution 주의사항:**
- Individual claims resolve: selective 설정 적용됨 ✅
- Game resolve (tryResolve): selective 체크 누락 (코드 이슈) ❌
- 자세한 내용: `op-challenger/scripts/README-KR.md` 참조

### 6. GameSolver (`game/fault/solver/`)

최적의 게임 전략을 계산합니다.

**핵심 알고리즘:**

```go
// game_solver.go:26-74
func (s *GameSolver) CalculateNextActions(ctx, game) ([]Action, error) {
    1. Root claim과의 동의 여부 판단
       - agreeWithClaim()으로 로컬 trace와 비교

    2. L2 블록 번호 검증
       - 동의하는 경우, L2 블록 유효성 확인
       - 무효하면 ChallengeL2BlockNumber 반환

    3. 모든 클레임 순회
       - maxDepth인 경우: calculateStep()
       - 중간 depth: calculateMove()

    4. 정직한 클레임 추적 (honestClaimTracker)
       - 이미 올바른 대응이 있는지 확인
       - 중복 방지
}
```

**Move vs Step:**

- **Move**: 중간 노드에 대한 공격/방어 (새로운 해시 제시)
- **Step**: 리프 노드에 대한 최종 증명 (VM 실행 증명)

```go
// claim_solver.go 주요 로직
func (s *claimSolver) NextMove(claim, game) (*Claim, error) {
    1. 부모 클레임과 동의 여부 판단
    2. Position 계산 (attack or defend)
    3. TraceAccessor.Get(position)으로 올바른 값 조회
    4. 기존에 올바른 클레임이 있으면 nil 반환
    5. 없으면 새 클레임 생성
}
```

### 7. TraceAccessor (`game/fault/trace/`)

게임의 각 포지션에서 정확한 상태 값을 제공합니다.

**인터페이스:**

```go
type TraceAccessor interface {
    Get(ctx context.Context, pos Position) (common.Hash, error)
    GetStepData(ctx context.Context, pos Position) (
        prestate []byte,
        proofData []byte,
        oracleData *PreimageOracleData,
        error
    )
    GetL2BlockNumberChallenge(ctx) (*InvalidL2BlockNumberChallenge, error)
}
```

**Cannon Trace Provider 구조:**

```
CannonTraceProvider
├── Executor        # Cannon VM 실행
├── PreimageLoader  # Preimage 데이터 로드
└── PrestateProvider # 초기 상태 제공
```

**증명 생성 흐름:**

```go
// cannon/provider.go:103-170
func (p *CannonTraceProvider) loadProof(ctx, i uint64) (*ProofData, error) {
    1. 디스크 캐시에서 증명 파일 찾기
       - 경로: {dir}/proofs/{i}.json.gz

    2. 없으면 Generator로 생성
       - p.generator.GenerateProof(ctx, dir, i)
       - Cannon VM 실행하여 step i까지 진행

    3. 프로그램 종료 확인
       - 마지막 step 이후는 no-op으로 확장

    4. ProofData 반환
       - ClaimValue: 해당 step의 상태 해시
       - StateData: VM 상태 인코딩
       - ProofData: 증명 데이터
       - OracleData: Preimage oracle 정보
}
```

### 8. Responder (`game/fault/responder/`)

계산된 액션을 실제 트랜잭션으로 변환하고 실행합니다.

**핵심 메소드:**

```go
// responder.go
type Responder interface {
    PerformAction(ctx context.Context, action types.Action) error
    Resolve() error
    ResolveClaims(claimIdx ...uint64) error
    CallResolve(ctx context.Context) (GameStatus, error)
    CallResolveClaim(ctx context.Context, claimIdx uint64) error
}
```

**PerformAction 처리:**

```go
func (r *FaultResponder) PerformAction(ctx, action) error {
    switch action.Type {
    case ActionTypeMove:
        // Attack or Defend 트랜잭션 생성
        if action.IsAttack {
            return r.Attack(parent, action.Value)
        }
        return r.Defend(parent, action.Value)

    case ActionTypeStep:
        // Preimage oracle 업데이트 (필요 시)
        if action.OracleData != nil {
            r.contract.UpdateOracleTx(...)
        }
        // Step 트랜잭션 생성
        return r.Step(claimIdx, isAttack, prestate, proof)

    case ActionTypeChallengeL2BlockNumber:
        return r.ChallengeL2BlockNumber(challenge)
    }
}
```

### 9. Contract Bindings (`game/fault/contracts/`)

온체인 dispute game contract와 상호작용합니다.

**주요 Contract:**

1. **FaultDisputeGameContract**: 메인 게임 로직
2. **PreimageOracleContract**: Preimage 데이터 관리
3. **DelayedWETHContract**: Bond WETH 관리
4. **DisputeGameFactoryContract**: 게임 생성 및 조회

**주요 메소드:**

```go
// FaultDisputeGameContract 인터페이스
GetAllClaims(ctx, block) ([]Claim, error)
GetGameMetadata(ctx, block) (GameMetadata, error)
AttackTx(ctx, parent, pivot) (TxCandidate, error)
DefendTx(ctx, parent, pivot) (TxCandidate, error)
StepTx(claimIdx, isAttack, stateData, proof) (TxCandidate, error)
CallResolve(ctx) (GameStatus, error)
ResolveTx() (TxCandidate, error)
```

**버전 호환성:**

```go
// contracts/faultdisputegame.go:83-129
func NewFaultDisputeGameContract(...) (FaultDisputeGameContract, error) {
    // 버전 조회
    version := caller.SingleCall(methodVersion)

    // 버전별 호환 레이어 선택
    if strings.HasPrefix(version, "0.8.") {
        return &FaultDisputeGameContract080{...}
    } else if strings.HasPrefix(version, "0.18.") || ... {
        return &FaultDisputeGameContract0180{...}
    } else if strings.HasPrefix(version, "1.1.") {
        return &FaultDisputeGameContract111{...}
    } else {
        return &FaultDisputeGameContractLatest{...}
    }
}
```

### 10. BondClaimer (`game/fault/claims/`)

게임 종료 후 bond를 회수합니다.

**동작 방식:**

```go
// BondClaimScheduler
1. 주기적으로 게임 스캔
2. 해결된 게임 확인
3. Credit 잔액 조회
4. ClaimCredit 트랜잭션 전송
```

---

## Dispute Game 메커니즘

### 게임 생명주기

```
1. [CREATED]
   ├─ Factory에서 게임 생성
   └─ Root claim 설정

2. [IN_PROGRESS]
   ├─ Claim/Counter-claim 교환
   ├─ Bisection으로 분쟁 범위 축소
   └─ 최종적으로 단일 step까지 축소

3. [RESOLVED]
   ├─ 모든 클레임 해결
   ├─ 게임 상태 확정
   └─ Bond 분배
```

### Bisection Protocol

게임은 이진 검색 방식으로 정확한 불일치 지점을 찾습니다:

```
예시: 1000 step의 실행 trace

Depth 0 (Root):  [0 -------- 1000] (Output Root)
                        ↓ 불일치
Depth 1:        [0 --- 500 --- 1000]
                      ↓ 불일치
Depth 2:        [0 - 250 - 500]
                      ...
Depth 73 (Leaf): [250] <- 단일 instruction
```

**각 depth에서:**
- 공격자(Attacker): 부모 클레임이 틀렸다고 주장
- 방어자(Defender): 부모 클레임이 맞다고 주장

### Position System

Position은 gindex로 인코딩됩니다:

```go
// types/position.go
type Position struct {
    depth      Depth
    indexAtDepth *big.Int
}

// Attack: 왼쪽 자식으로 이동
func (p Position) Attack() Position {
    return NewPosition(p.depth+1, p.indexAtDepth*2)
}

// Defend: 오른쪽 자식으로 이동
func (p Position) Defend() Position {
    return NewPosition(p.depth+1, p.indexAtDepth*2+1)
}

// TraceIndex: 실제 VM step 인덱스 계산
func (p Position) TraceIndex(maxDepth Depth) *big.Int {
    // (indexAtDepth + 1) * 2^(maxDepth - depth) - 1
}
```

### Chess Clock 메커니즘

각 플레이어는 제한된 시간을 가지며, 차례가 오면 시간이 소모됩니다:

```go
// types/claim.go
type Clock struct {
    Duration  time.Duration  // 남은 시간
    Timestamp time.Time      // 마지막 행동 시각
}

// ChessClock 계산
func ChessClock(now time.Time, claim Claim, parent Claim) time.Duration {
    elapsed := now.Sub(claim.Clock.Timestamp)
    if claim.IsRoot() {
        return claim.Clock.Duration - elapsed
    }
    // 부모의 남은 시간에서 경과 시간 차감
    return parent.Clock.Duration - claim.Clock.Duration - elapsed
}
```

시간이 만료되면:
- 해당 클레임은 자동 패배
- 카운터된 클레임이 승리
- ResolveClaim으로 해결 가능

### Step Execution (리프 노드)

최대 depth에 도달하면 실제 VM 실행을 증명합니다:

```go
// Step 증명 구성요소
type StepData struct {
    PreState   []byte  // step 실행 전 VM 상태
    ProofData  []byte  // Merkle proof 등
    OracleData *PreimageOracleData  // Preimage 정보 (필요 시)
}
```

**검증 과정:**

1. PreState를 VM에 로드
2. 단일 instruction 실행
3. PostState 계산
4. PostState가 클레임 값과 일치하는지 검증

온체인에서 VM을 실행하므로 **결정론적 증명**이 가능합니다.

### Preimage Oracle

VM 실행 중 필요한 외부 데이터를 제공합니다:

**종류:**
1. **Local Preimage**: 게임별로 고유 (L2 block header 등)
2. **Global Preimage**: 모든 게임에서 공유 (L1 block data 등)
3. **Large Preimage**: 큰 데이터 (Keccak256 등)

**Large Preimage 처리:**

```go
// keccak/ 디렉토리
1. LargePreimageScheduler: 주기적으로 검증
2. PreimageChallenger: 무효한 preimage 챌린지
3. PreimageVerifier: L1에서 데이터 검증
```

---

## Trace Provider 시스템

### Cannon Provider

MIPS 아키텍처 기반 VM으로 `op-program`을 실행합니다.

**실행 프로세스:**

```go
// cannon/executor.go
type Executor struct {
    cannon     string  // cannon 바이너리 경로
    server     string  // op-program 경로
    prestate   string  // 초기 상태 파일
    inputs     LocalGameInputs
}

func (e *Executor) GenerateProof(ctx, dir, step uint64) error {
    1. Snapshot 경로 계산
       - 가장 가까운 이전 snapshot 찾기

    2. Cannon 실행
       cmd := exec.Command(e.cannon, "run",
           "--input", snapshot,
           "--output", finalState,
           "--meta", metaPath,
           "--proof-at", step,
           "--proof-fmt", proofPath,
           "--",
           e.server, serverArgs...)

    3. Proof 파일 생성
       - {step}.json.gz에 저장
       - StateData, ProofData 포함
}
```

**Snapshot 메커니즘:**

```
Snapshot Frequency = 1,000,000,000 instructions

0 (prestate)
└─> 1B snapshot
    └─> 2B snapshot
        └─> ...
```

빠른 증명 생성을 위해 중간 상태를 저장합니다.

### Output Provider

L2 output root를 검증하는 trace provider입니다.

**Split Depth:**

```
Depth < splitDepth: L2 output 검증
Depth >= splitDepth: VM 실행 검증 (Cannon/Asterisc)
```

**동작:**

```go
// outputs/provider.go
func (p *OutputTraceProvider) Get(ctx, pos Position) (common.Hash, error) {
    if pos.Depth() <= p.prestateDepth {
        // 초기 상태 반환
        return p.prestateProvider.AbsolutePreStateCommitment(ctx)
    }

    // 해당 position의 L2 block 번호 계산
    outputBlock := p.outputAtBlock(pos.TraceIndex())

    // Rollup client에서 output 조회
    output := p.rollupClient.OutputAtBlock(outputBlock)

    return output.OutputRoot, nil
}
```

### Alphabet Provider (테스트용)

간단한 문자열 기반 trace로 테스트에 사용됩니다:

```go
// alphabet/provider.go
// "abcdefgh" 같은 문자열을 trace로 사용
func (p *AlphabetTraceProvider) Get(ctx, pos Position) (common.Hash, error) {
    index := pos.TraceIndex()
    letter := p.state[index]
    return common.BytesToHash([]byte{letter}), nil
}
```

---

## 실행 흐름

### 전체 프로세스

```
1. 초기화
   ├─ Config 로드 및 검증
   ├─ L1/L2 클라이언트 연결
   ├─ Contract 바인딩
   └─ 게임 타입 등록

2. 모니터링 시작
   ├─ L1 newHeads 구독
   └─ 새 블록마다 onNewL1Head 호출

3. 게임 발견
   ├─ DisputeGameFactory.GetGamesAtOrAfter()
   ├─ GameWindow 필터링
   └─ Allowlist 확인

4. 게임 스케줄링
   ├─ Scheduler.Schedule(games, blockNumber)
   ├─ 각 게임에 대해 GamePlayer 생성
   └─ Worker pool에 작업 분배

5. 게임 진행 (Worker)
   ├─ Player.ProgressGame()
   ├─ Agent.Act()
   │   ├─ Solver.CalculateNextActions()
   │   └─ Responder.PerformAction()
   └─ 상태 업데이트

6. 트랜잭션 전송
   ├─ TxSender에 전달
   ├─ TxManager로 전송
   └─ Receipt 확인

7. 게임 해결
   ├─ Agent.tryResolve()
   ├─ ResolveClaims()
   └─ Resolve()

8. Bond 회수
   ├─ BondClaimScheduler
   ├─ Credit 확인
   └─ ClaimCredit 트랜잭션
```

### 단일 게임의 Act 사이클

```
┌─────────────────────────────────────────────┐
│         Agent.Act() - 매 L1 블록마다         │
└─────────────────────────────────────────────┘
                    ↓
        ┌──────────────────────┐
        │  1. tryResolve()?    │
        │  - Yes → 게임 종료    │
        │  - No  → 계속        │
        └──────────┬───────────┘
                   ↓
        ┌──────────────────────┐
        │ 2. L2 Block 체크     │
        │  - 이미 챌린지됨?     │
        └──────────┬───────────┘
                   ↓
        ┌──────────────────────┐
        │ 3. Claims 로드       │
        │  - GetAllClaims()    │
        └──────────┬───────────┘
                   ↓
        ┌─────────────────────────────────┐
        │ 4. Solver.CalculateNextActions() │
        │  ┌──────────────────────────┐   │
        │  │ Root claim 동의 여부?     │   │
        │  └───────┬──────────────────┘   │
        │          ↓                       │
        │  ┌──────────────────────────┐   │
        │  │ L2 block valid?          │   │
        │  │  Invalid → Challenge      │   │
        │  └───────┬──────────────────┘   │
        │          ↓                       │
        │  ┌──────────────────────────┐   │
        │  │ 각 Claim 순회            │   │
        │  │  - maxDepth: Step        │   │
        │  │  - else: Move            │   │
        │  └──────────────────────────┘   │
        └─────────────┬───────────────────┘
                      ↓
        ┌─────────────────────────────┐
        │ 5. Actions 병렬 실행         │
        │  ┌────────────────────┐     │
        │  │ Move: Attack/Defend │     │
        │  └────────────────────┘     │
        │  ┌────────────────────┐     │
        │  │ Step: Proof 제출    │     │
        │  └────────────────────┘     │
        │  ┌────────────────────┐     │
        │  │ Challenge L2 Block  │     │
        │  └────────────────────┘     │
        └─────────────────────────────┘
```

### Solver 알고리즘 상세

```go
// ClaimSolver 핵심 로직
for each claim in game.Claims() {
    // 1. 이 클레임과 동의하는가?
    agree := agreeWithClaim(claim)

    // 2. 부모와의 관계
    parent := game.Claims()[claim.ParentContractIndex]
    agreeWithParent := agreeWithClaim(parent)

    // 3. 논리적 대응 결정
    if agree == agreeWithParent {
        // 모순: 부모는 맞는데 자식이 맞다?
        // 또는 부모는 틀린데 자식이 틀리다?
        // → Defend (같은 편)
        position = claim.Position.Defend()
    } else {
        // 정상: 부모 틀림 & 자식 맞음
        // 또는 부모 맞음 & 자식 틀림
        // → Attack (반대편)
        position = claim.Position.Attack()
    }

    // 4. 올바른 값 계산
    correctValue := traceAccessor.Get(position)

    // 5. 이미 올바른 클레임이 있는지 확인
    if honestClaimTracker.Contains(claim, position, correctValue) {
        continue  // 이미 대응됨
    }

    // 6. 새 클레임 생성
    actions = append(actions, Action{
        Type: ActionTypeMove,
        IsAttack: shouldAttack,
        ParentClaim: claim,
        Value: correctValue,
    })

    honestClaimTracker.Add(claim, newClaim)
}
```

### TraceAccessor 호출 흐름

```
Position 계산
    ↓
TraceAccessor.Get(position)
    ↓
┌─────────────────────────────────┐
│ OutputTraceProvider             │
│  - Depth < splitDepth?          │
│    Yes: L2 output 반환          │
│    No: Delegate to Cannon/...   │
└─────────────┬───────────────────┘
              ↓
┌─────────────────────────────────┐
│ CannonTraceProvider             │
│  - TraceIndex 계산              │
│  - loadProof(traceIndex)        │
└─────────────┬───────────────────┘
              ↓
┌─────────────────────────────────┐
│ Proof 파일 존재?                │
│  Yes: 파일 읽기                 │
│  No: Executor.GenerateProof()   │
└─────────────┬───────────────────┘
              ↓
┌─────────────────────────────────┐
│ Cannon VM 실행                  │
│  1. Prestate 로드               │
│  2. Snapshot에서 시작           │
│  3. Step까지 실행               │
│  4. Proof 저장                  │
└─────────────┬───────────────────┘
              ↓
          ClaimValue 반환
```

---

## 설정 및 배포

### 필수 설정 항목

```go
// config/config.go
type Config struct {
    // L1 연결
    L1EthRpc   string           // L1 RPC endpoint
    L1Beacon   string           // L1 Beacon API

    // L2 연결
    RollupRpc  string           // Rollup node RPC
    L2Rpc      string           // L2 execution node RPC

    // Contract 주소
    GameFactoryAddress common.Address

    // 게임 필터
    GameAllowlist []common.Address  // 특정 게임만 참여
    GameWindow    time.Duration     // 최근 N일 게임만 모니터링

    // 성능
    MaxConcurrency uint             // 동시 처리 게임 수
    PollInterval   time.Duration    // RPC 폴링 간격

    // Trace 타입
    TraceTypes []TraceType          // 지원할 trace 종류

    // Cannon 설정
    CannonBin                string
    CannonServer             string
    CannonAbsolutePreState   string
    CannonRollupConfigPath   string
    CannonL2GenesisPath      string
    CannonSnapshotFreq       uint

    // Bond 청구
    AdditionalBondClaimants  []common.Address
    SelectiveClaimResolution bool

    // 트랜잭션
    TxMgrConfig   txmgr.CLIConfig
    MaxPendingTx  uint64
}
```

### 실행 예시

**Local Devnet:**

```bash
DISPUTE_GAME_FACTORY=$(jq -r .DisputeGameFactoryProxy .devnet/addresses.json)

./op-challenger/bin/op-challenger \
  --trace-type cannon \
  --l1-eth-rpc http://localhost:8545 \
  --l1-beacon http://localhost:5052 \
  --rollup-rpc http://localhost:9546 \
  --l2-eth-rpc http://localhost:9545 \
  --game-factory-address $DISPUTE_GAME_FACTORY \
  --datadir temp/challenger-data \
  --cannon-rollup-config .devnet/rollup.json \
  --cannon-l2-genesis .devnet/genesis-l2.json \
  --cannon-bin ./cannon/bin/cannon \
  --cannon-server ./op-program/bin/op-program \
  --cannon-prestate ./op-program/bin/prestate.json \
  --mnemonic "test test test test test test test test test test test junk" \
  --hd-path "m/44'/60'/0'/0/8" \
  --num-confirmations 1
```

**Production 고려사항:**

```bash
./op-challenger/bin/op-challenger \
  # 보안
  --private-key $PRIVATE_KEY \

  # 성능
  --max-concurrency 4 \
  --max-pending-tx 10 \

  # 모니터링
  --metrics.enabled \
  --metrics.port 7300 \
  --pprof.enabled \

  # 게임 필터
  --game-window 672h \  # 28 days

  # 고급 설정
  --selective-claim-resolution \
  --additional-bond-claimants $TREASURY_ADDRESS \

  # 백업 RPC
  --l1-eth-rpc https://mainnet.infura.io/... \
  --l1-eth-rpc https://eth.llamarpc.com/...
```

### 디렉토리 구조

```
{datadir}/
├── {game-address-1}/
│   ├── proofs/
│   │   ├── 0.json.gz
│   │   ├── 1000000000.json.gz  # snapshot
│   │   └── ...
│   ├── preimages/
│   │   └── kv/
│   ├── state.json              # final state
│   └── meta.json               # last step info
├── {game-address-2}/
│   └── ...
└── ...
```

게임이 해결되면 해당 디렉토리는 자동으로 삭제됩니다.

### 모니터링

**메트릭:**

```go
// metrics/metrics.go
- games_in_progress
- games_defender_won
- games_challenger_won
- game_act_time
- game_move_count
- game_step_count
- claim_resolution_time
- cannon_execution_time
```

**로그 예시:**

```
INFO [01-15|12:34:56.789] Starting op-challenger       version=v1.0.0
INFO [01-15|12:34:57.123] Starting monitoring
INFO [01-15|12:35:00.456] Game info                    game=0x1234... claims=5 status=IN_PROGRESS
INFO [01-15|12:35:01.789] Performing action            action=move is_attack=true parent=2 value=0xabcd...
INFO [01-15|12:36:00.123] Resolving game               game=0x1234...
INFO [01-15|12:36:01.456] Game resolved                game=0x1234... status=CHALLENGER_WINS
```

### CLI 서브커맨드

**게임 생성:**

```bash
./bin/op-challenger create-game \
  --l1-eth-rpc http://localhost:8545 \
  --game-factory-address $FACTORY \
  --output-root $ROOT_HASH \
  --l2-block-num 1000 \
  --private-key $KEY
```

**게임 목록:**

```bash
./bin/op-challenger list-games \
  --l1-eth-rpc http://localhost:8545 \
  --game-factory-address $FACTORY
```

**클레임 목록:**

```bash
./bin/op-challenger list-claims \
  --l1-eth-rpc http://localhost:8545 \
  --game-address $GAME
```

**수동 Move:**

```bash
./bin/op-challenger move \
  --l1-eth-rpc http://localhost:8545 \
  --game-address $GAME \
  --attack \
  --parent-index 5 \
  --claim $CLAIM_HASH \
  --private-key $KEY
```

**게임 해결:**

```bash
./bin/op-challenger resolve \
  --l1-eth-rpc http://localhost:8545 \
  --game-address $GAME \
  --private-key $KEY
```

---

## 부록

### 주요 타입 정의

```go
// game/types/types.go
type GameMetadata struct {
    GameType  uint32
    Timestamp uint64
    Proxy     common.Address
}

// game/fault/types/types.go
type Claim struct {
    ClaimData
    CounteredBy         common.Address
    Claimant            common.Address
    Clock               Clock
    ContractIndex       int
    ParentContractIndex int
}

type ClaimData struct {
    Value    common.Hash
    Position Position
    Bond     *big.Int
}

type Position struct {
    depth        Depth
    indexAtDepth *big.Int
}

type Action struct {
    Type                          ActionType
    ParentClaim                   Claim
    IsAttack                      bool
    Value                         common.Hash
    PreState                      []byte
    ProofData                     []byte
    OracleData                    *PreimageOracleData
    InvalidL2BlockNumberChallenge *InvalidL2BlockNumberChallenge
}

type ActionType string
const (
    ActionTypeMove                       ActionType = "move"
    ActionTypeStep                       ActionType = "step"
    ActionTypeChallengeL2BlockNumber     ActionType = "challenge-l2-block"
)
```

### 주요 인터페이스

```go
// scheduler/types.go
type GamePlayer interface {
    ValidatePrestate(ctx context.Context) error
    ProgressGame(ctx context.Context) GameStatus
    Status() GameStatus
}

// game/fault/types/types.go
type TraceAccessor interface {
    Get(ctx context.Context, pos Position) (common.Hash, error)
    GetStepData(ctx context.Context, pos Position) (prestate []byte, proofData []byte, oracleData *PreimageOracleData, err error)
    GetL2BlockNumberChallenge(ctx context.Context) (*InvalidL2BlockNumberChallenge, error)
}

type Game interface {
    Claims() []Claim
    MaxDepth() Depth
    IsDuplicate(claim Claim) bool
    DefendsParent(claim Claim) bool
}
```

### 에러 처리

```go
// 주요 에러 타입
var (
    ErrInvalidPrestate        = errors.New("invalid prestate")
    ErrL2BlockNumberValid     = errors.New("l2 block number is valid")
    ErrGameDepthReached       = errors.New("game depth reached")
    ErrClaimNotCountered      = errors.New("claim has not been countered")
)
```

### 성능 최적화

1. **Snapshot 활용**: Cannon 실행 시 중간 상태 저장
2. **병렬 처리**: Worker pool로 여러 게임 동시 처리
3. **캐싱**:
   - Proof 파일 디스크 캐싱
   - Last step 메타데이터 캐싱
4. **Batching**: MultiCaller로 RPC 호출 일괄 처리

### 보안 고려사항

1. **Private Key 관리**: 환경변수 또는 KMS 사용
2. **Gas Price 제한**: TxManager 설정
3. **재진입 공격 방지**: Contract 레벨에서 처리
4. **타임아웃 설정**: Context timeout 적절히 설정

---

## 참고 자료

- [Optimism Fault Proof Specs](https://specs.optimism.io/experimental/fault-proof/index.html)
- [op-challenger README](../README.md)
- [Cannon VM Documentation](../../cannon/README.md)
- [op-program Documentation](../../op-program/README.md)
