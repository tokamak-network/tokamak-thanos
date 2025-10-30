# Optimism 구조 마이그레이션 분석 (GameType 3 통합용)

## 개요

GameType 3 (ASTERISC_KONA)를 Optimism 방식으로 통합하려다 발견한 **Tokamak-Thanos와 Optimism 간의 구조적 차이점**을 분석합니다.

**작성일**: 2025-10-26
**목적**: Optimism 구조로 완전 마이그레이션 계획 수립
**분석 범위**: op-challenger 패키지

---

## 🔍 핵심 차이점 발견

### 1. GameType 타입 정의

**Optimism (신버전)**:
```go
// op-challenger/game/fault/types/types.go:27-28
type GameType uint32

const (
    CannonGameType            GameType = 0
    PermissionedGameType      GameType = 1
    AsteriscGameType          GameType = 2
    AsteriscKonaGameType      GameType = 3
    // ...
)
```

**Tokamak-Thanos (구버전)**:
```go
// op-challenger/game/fault/types/types.go:21-28
const (
    CannonGameType       uint32 = 0
    PermissionedGameType uint32 = 1
    AsteriscGameType     uint32 = 2
    AsteriscKonaGameType uint32 = 3  // ← 추가함
    FastGameType         uint32 = 254
    AlphabetGameType     uint32 = 255
)
```

**차이**: Optimism은 `type GameType uint32` 사용, Tokamak은 `uint32` 직접 사용

**영향**:
- Registry 인터페이스 시그니처
- registerOracle 함수 시그니처
- 모든 GameType 관련 함수 파라미터

---

### 2. Config 구조

#### 2.1 Cannon Config

**Optimism**:
```go
// op-challenger/config/config.go:85-87
Cannon                        vm.Config  // ← vm.Config 사용
CannonAbsolutePreState        string
CannonAbsolutePreStateBaseURL *url.URL
```

**Tokamak-Thanos**:
```go
// op-challenger/config/config.go:139-147
CannonBin                     string   // ← 개별 필드
CannonServer                  string   // ← 개별 필드
CannonAbsolutePreState        string
CannonAbsolutePreStateBaseURL *url.URL
CannonNetwork                 string
CannonRollupConfigPath        string
CannonL2GenesisPath           string
CannonSnapshotFreq            uint
CannonInfoFreq                uint
```

**차이**:
- Optimism: `Cannon vm.Config`로 통합
- Tokamak: 개별 필드로 관리

**현재 Tokamak 상태**:
- ✅ `Asterisc vm.Config` **이미 적용됨** (GameType 2 작업 시)
- ✅ `AsteriscKona vm.Config` **이미 추가함** (방금 작업)
- ❌ `Cannon vm.Config` **아직 구버전** 사용 중

#### 2.2 Asterisc Config

**Optimism**:
```go
Asterisc                            vm.Config
AsteriscAbsolutePreState            string
AsteriscAbsolutePreStateBaseURL     *url.URL
AsteriscKona                        vm.Config
AsteriscKonaAbsolutePreState        string
AsteriscKonaAbsolutePreStateBaseURL *url.URL
```

**Tokamak-Thanos**:
```go
Asterisc                        vm.Config  // ✅ 같음
AsteriscBin                     string     // ← 추가 필드
AsteriscServer                  string     // ← 추가 필드
AsteriscAbsolutePreState        string     // ✅ 같음
AsteriscAbsolutePreStateBaseURL *url.URL   // ✅ 같음
AsteriscNetwork                 string     // ← 추가 필드
AsteriscRollupConfigPath        string     // ← 추가 필드
AsteriscL2GenesisPath           string     // ← 추가 필드
AsteriscSnapshotFreq            uint       // ← 추가 필드
AsteriscInfoFreq                uint       // ← 추가 필드

AsteriscKona                        vm.Config  // ✅ 같음
AsteriscKonaAbsolutePreState        string     // ✅ 같음
AsteriscKonaAbsolutePreStateBaseURL *url.URL   // ✅ 같음
```

**차이**:
- Optimism: `vm.Config`만 사용, 나머지는 vm.Config 내부 필드로
- Tokamak: `vm.Config` + 추가 개별 필드들

**하이브리드 상태**:
- GameType 2 (Asterisc) 작업 시 `vm.Config` 도입
- 하지만 하위 호환성을 위해 개별 필드도 유지
- Check() 함수에서 개별 필드를 vm.Config로 복사

---

### 3. Prestates 패키지 API

**Optimism**:
```go
// register_task.go에서 사용
prestateSource.PrestatePath(ctx, prestateHash)  // ← ctx 파라미터 있음

// 인터페이스
type PrestateSource interface {
    PrestatePath(ctx context.Context, prestateHash common.Hash) (string, error)
}
```

**Tokamak-Thanos**:
```go
// register.go에서 사용
prestateSource.PrestatePath(prestateHash)  // ← ctx 파라미터 없음

// 인터페이스
type PrestateSource interface {
    PrestatePath(prestateHash common.Hash) (string, error)
}
```

**영향 범위**:
- `prestates/single.go`
- `prestates/multi.go`
- `register.go` (PrestateSource 인터페이스)

---

### 4. Prestates 캐싱 구조

**Optimism**:
```go
// register_task.go:264-274
prestateProviderCache := prestates.NewPrestateProviderCache(m, fmt.Sprintf("prestates-%v", gameType),
    func(ctx context.Context, prestateHash common.Hash) (faultTypes.PrestateProvider, error) {
        // ← ctx 파라미터 있음
        prestatePath, err := prestateSource.PrestatePath(ctx, prestateHash)
        // ...
    })
return prestateProviderCache.GetOrCreate  // ← ctx 있는 버전
```

**Tokamak-Thanos**:
```go
// register.go:214
prestateProviderCache := prestates.NewPrestateProviderCache(m, fmt.Sprintf("prestates-%v", gameType),
    func(prestateHash common.Hash) (faultTypes.PrestateProvider, error) {
        // ← ctx 파라미터 없음
        prestatePath, err := prestateSource.PrestatePath(prestateHash)
        // ...
    })
```

**차이**:
- Optimism: Context 전달 (취소 가능, 타임아웃 지원)
- Tokamak: Context 없음

---

### 5. Contract 메서드

**Optimism**:
```go
// register_task.go:312
prestateBlock, poststateBlock, err := contract.GetGameRange(ctx)
```

**Tokamak-Thanos**:
```go
// register.go (기존 코드)
prestateBlock, poststateBlock, err := contract.GetBlockRange(ctx)
```

**차이**: 메서드 이름이 다름

**확인 필요**:
```bash
grep -r "GetGameRange\|GetBlockRange" contracts/*.go
```

---

### 6. Registry 인터페이스

**Optimism**:
```go
type Registry interface {
    RegisterGameType(gameType faultTypes.GameType, creator scheduler.PlayerCreator)  // ← GameType 타입
    RegisterBondContract(gameType faultTypes.GameType, creator claims.BondContractCreator)
}
```

**Tokamak-Thanos**:
```go
type Registry interface {
    RegisterGameType(gameType uint32, creator scheduler.PlayerCreator)  // ← uint32
    RegisterBondContract(gameType uint32, creator claims.BondContractCreator)
}
```

---

### 7. DialRollupClientWithTimeout 시그니처

**에러 메시지**:
```
not enough arguments in call to dial.DialRollupClientWithTimeout
have (context.Context, log.Logger, string)
want (context.Context, time.Duration, log.Logger, string)
```

**Optimism**:
```go
dial.DialRollupClientWithTimeout(ctx, logger, rpc)
```

**Tokamak-Thanos 실제**:
```go
dial.DialRollupClientWithTimeout(ctx, timeout, logger, rpc)  // ← timeout 파라미터 추가
```

---

## 📋 마이그레이션 필요 항목

### Level 1: 타입 시스템 변경 (핵심)

| 항목 | 파일 | 변경 내용 | 영향도 |
|-----|------|----------|-------|
| **GameType 타입** | `types/types.go` | `type GameType uint32` 추가 | 🔴 높음 |
| **Registry** | `register.go` | 인터페이스 시그니처 변경 | 🔴 높음 |
| **모든 register 함수** | `register.go` | GameType 파라미터 타입 변경 | 🔴 높음 |

### Level 2: Prestates 패키지 업데이트

| 항목 | 파일 | 변경 내용 | 영향도 |
|-----|------|----------|-------|
| **PrestateSource 인터페이스** | `register.go` | ctx 파라미터 추가 | 🟡 중간 |
| **SinglePrestateSource** | `prestates/single.go` | PrestatePath에 ctx 추가 | 🟡 중간 |
| **MultiPrestateProvider** | `prestates/multi.go` | PrestatePath에 ctx 추가 | 🟡 중간 |
| **NewPrestateSource** | `prestates/source.go` | 함수 추가 (Optimism에 있음) | 🟡 중간 |

### Level 3: Contract 메서드 호환성

| 항목 | 파일 | 변경 내용 | 영향도 |
|-----|------|----------|-------|
| **GetGameRange** | `contracts/faultdisputegame.go` | GetBlockRange → GetGameRange | 🟢 낮음 |
| 또는 | `register_task.go` | GetGameRange → GetBlockRange | 🟢 낮음 |

### Level 4: Dial 함수 시그니처

| 항목 | 파일 | 변경 내용 | 영향도 |
|-----|------|----------|-------|
| **DialRollupClient** | `clients.go` | timeout 파라미터 확인/추가 | 🟢 낮음 |

---

## 🔄 마이그레이션 전략

### 전략 A: 전면 마이그레이션 (Optimism 완전 적용)

**장점**:
- ✅ Optimism과 100% 동일한 구조
- ✅ 향후 Optimism 업데이트 쉬움
- ✅ GameType 4, 6, 7 추가 용이

**단점**:
- ❌ 작업량 많음 (3-5일)
- ❌ 기존 코드 대규모 수정
- ❌ 회귀 테스트 필요

**작업 순서**:
1. **GameType 타입 변경** (1일)
   - `type GameType uint32` 추가
   - Registry 인터페이스 변경
   - 모든 호출부 수정

2. **Prestates API 업데이트** (1일)
   - ctx 파라미터 추가
   - NewPrestateSource 함수 추가

3. **Register 구조 변경** (1-2일)
   - register_task.go 사용
   - clients.go 사용
   - 기존 register* 함수 제거

4. **테스트 및 검증** (1일)
   - 모든 GameType 회귀 테스트
   - 기존 기능 정상 동작 확인

---

### 전략 B: 점진적 마이그레이션 (하이브리드)

**현재 상태 유지**:
- ✅ GameType은 uint32로 유지
- ✅ Prestates API는 구버전 유지
- ✅ 기존 코드 안정성 유지

**GameType 3만 추가**:
- ✅ register_task.go 사용 (새로운 스타일)
- ✅ 기존 register* 함수와 공존
- ✅ 점진적 전환

**장점**:
- ✅ 작업량 적음 (1일)
- ✅ 리스크 낮음
- ✅ 빠른 GameType 3 지원

**단점**:
- ❌ 하이브리드 구조 (일관성 떨어짐)
- ❌ 향후 Optimism 동기화 복잡

---

## 📊 파일별 변경 분석

### 필수 변경 파일 (전략 A)

#### 1. types/types.go

**현재**:
```go
const (
    CannonGameType       uint32 = 0
    // ...
)
```

**변경 후** (Optimism 스타일):
```go
type GameType uint32

const (
    CannonGameType            GameType = 0
    PermissionedGameType      GameType = 1
    AsteriscGameType          GameType = 2
    AsteriscKonaGameType      GameType = 3
    FastGameType              GameType = 254
    AlphabetGameType          GameType = 255
    UnknownGameType           GameType = math.MaxUint32
)

func (t GameType) MarshalText() ([]byte, error) {
    return []byte(t.String()), nil
}

func (t GameType) String() string {
    switch t {
    case CannonGameType:
        return "Cannon"
    case PermissionedGameType:
        return "PermissionedCannon"
    case AsteriscGameType:
        return "Asterisc"
    case AsteriscKonaGameType:
        return "AsteriscKona"
    case FastGameType:
        return "Fast"
    case AlphabetGameType:
        return "Alphabet"
    default:
        return fmt.Sprintf("Unknown(%d)", t)
    }
}
```

**영향**:
- 모든 `uint32` GameType → `faultTypes.GameType`로 변경 필요
- 약 50+ 곳 수정 예상

---

#### 2. register.go

**현재** (387 lines, 복잡한 register* 함수들):
```go
func RegisterGameTypes(...) (CloseFunc, error) {
    l2Client, err := ethclient.DialContext(ctx, cfg.L2Rpc)
    syncValidator := newSyncStatusValidator(rollupClient)

    if cfg.TraceTypeEnabled(config.TraceTypeCannon) {
        if err := registerCannon(...); err != nil {  // ← 긴 파라미터 리스트
            return nil, err
        }
    }
    // ... registerAsterisc, registerAlphabet 등
}

func registerCannon(
    gameType uint32,  // ← 17개 파라미터!
    registry Registry,
    oracles OracleRegistry,
    // ... 14개 더
) error {
    // 200+ lines 복잡한 로직
}
```

**변경 후** (Optimism 스타일, 133 lines, 간단):
```go
func RegisterGameTypes(...) (CloseFunc, error) {
    clients := &clientProvider{ctx: ctx, logger: logger, cfg: cfg}
    var registerTasks []*RegisterTask

    if cfg.TraceTypeEnabled(faultTypes.TraceTypeCannon) {
        l2HeaderSource, rollupClient, syncValidator, err := clients.SingleChainClients()
        registerTasks = append(registerTasks, NewCannonRegisterTask(...))  // ← 간단!
    }

    for _, task := range registerTasks {
        if err := task.Register(...); err != nil {
            return clients.Close, err
        }
    }
    return clients.Close, nil
}

// register* 함수들 모두 제거, register_task.go로 이동
```

**장점**:
- ✅ 코드 50% 감소 (387 lines → 133 lines)
- ✅ 중복 제거
- ✅ RegisterTask 패턴으로 통일

---

#### 3. register_task.go (신규 생성)

**Optimism 구조** (387 lines):
```go
type RegisterTask struct {
    gameType               faultTypes.GameType
    skipPrestateValidation bool
    syncValidator          SyncValidator
    getTopPrestateProvider    func(...) (...)
    getBottomPrestateProvider func(...) (...)
    newTraceAccessor          func(...) (...)
}

func NewCannonRegisterTask(...) *RegisterTask { ... }
func NewAsteriscRegisterTask(...) *RegisterTask { ... }
func NewAsteriscKonaRegisterTask(...) *RegisterTask { ... }  // ← GameType 3
func NewAlphabetRegisterTask(...) *RegisterTask { ... }

func (e *RegisterTask) Register(...) error { ... }
func cachePrestates(...) func(...) { ... }
func registerOracle(...) error { ... }
func loadL1Head(...) (eth.BlockID, error) { ... }
```

**장점**:
- ✅ 각 GameType별 설정 격리
- ✅ 공통 로직 재사용 (Register 메서드)
- ✅ 새로운 GameType 추가 쉬움

---

#### 4. clients.go (신규 생성)

**Optimism 구조** (84 lines):
```go
type clientProvider struct {
    ctx            context.Context
    logger         log.Logger
    cfg            *config.Config
    l2HeaderSource utils.L2HeaderSource
    rollupClient   RollupClient
    syncValidator  *syncStatusValidator
    rootProvider   super.RootProvider  // ← Tokamak에 없음
    toClose        []CloseFunc
}

func (c *clientProvider) SingleChainClients() (...)
func (c *clientProvider) SuperchainClients() (...)  // ← Tokamak에 없음
func (c *clientProvider) Close()
```

**Tokamak 적용 시**:
- SuperchainClients는 제외 (super 패키지 없음)
- SingleChainClients만 구현

---

#### 5. prestates/source.go (확인 필요)

**Optimism**:
```go
func NewPrestateSource(
    baseURL *url.URL,
    preStatePath string,
    prestateDir string,
    stateConverter vm.StateConverter,
) PrestateSource { ... }
```

**Tokamak-Thanos**:
- 이 함수가 있는지 확인 필요
- 없으면 추가 또는 기존 방식 유지

---

### 6. outputs/output_cannon.go 시그니처

**에러**:
```
too many arguments in call to outputs.NewOutputCannonTraceAccessor
```

**Optimism 시그니처 확인 필요**:
```go
func NewOutputCannonTraceAccessor(
    logger log.Logger,
    m Metricer,
    cfg vm.Config,              // ← 이 파라미터가 있는지?
    serverExecutor vm.OracleServerExecutor,
    l2Client utils.L2HeaderSource,
    prestateProvider types.PrestateProvider,
    prestatePath string,
    rollupClient OutputRollupClient,
    dir string,
    l1Head eth.BlockID,
    splitDepth types.Depth,
    prestateBlock uint64,
    poststateBlock uint64,
) (*trace.Accessor, error)
```

---

## 🎯 권장 마이그레이션 계획

### 단계별 접근 (최소 리스크)

#### 단계 1: 타입 시스템 마이그레이션 (1-2일)

**작업**:
1. `type GameType uint32` 추가
2. 상수들을 GameType 타입으로 변경
3. Registry 인터페이스 변경
4. 모든 호출부 수정

**검증**:
- 모든 GameType (0, 1, 2) 빌드 성공
- 기존 테스트 통과

---

#### 단계 2: Prestates API 업데이트 (1일)

**작업**:
1. PrestateSource 인터페이스에 ctx 추가
2. single.go, multi.go PrestatePath 시그니처 변경
3. source.go 추가 (NewPrestateSource)
4. 캐시 함수 ctx 지원

**검증**:
- Prestate 로딩 테스트
- 멀티 prestate 테스트

---

#### 단계 3: Register 구조 변경 (1-2일)

**작업**:
1. register_task.go 생성
2. clients.go 생성
3. register.go 간소화
4. 기존 register* 함수 제거

**검증**:
- GameType 0, 1, 2 등록 테스트
- 회귀 테스트

---

#### 단계 4: GameType 3 추가 (0.5일)

**작업**:
1. NewAsteriscKonaRegisterTask 추가
2. RegisterGameTypes에 등록

**검증**:
- GameType 3 빌드 성공
- Unit 테스트

---

### 예상 일정

| 단계 | 작업 | 예상 기간 | 리스크 |
|-----|------|----------|-------|
| 단계 1 | 타입 시스템 | 1-2일 | 🔴 높음 |
| 단계 2 | Prestates API | 1일 | 🟡 중간 |
| 단계 3 | Register 구조 | 1-2일 | 🟡 중간 |
| 단계 4 | GameType 3 추가 | 0.5일 | 🟢 낮음 |
| **총계** | | **3.5-5.5일** | |

---

## ⚠️ 리스크 분석

### 높은 리스크

1. **GameType 타입 변경**
   - 영향 범위: 전체 codebase
   - 회귀 가능성: 높음
   - 완전한 테스트 필요

2. **기존 기능 깨짐**
   - GameType 0, 1, 2가 정상 작동해야 함
   - DevNet 전체 테스트 필수

### 중간 리스크

1. **Prestates API 변경**
   - 캐싱 로직 변경
   - Context 전달 구조

2. **Register 구조 변경**
   - 기존 패턴과 다름
   - 학습 비용

---

## 🔍 조사 결과 (실제 코드 분석)

### 1. Tokamak-Thanos의 Cannon Config ✅ 조사 완료

**발견**: ✅ **하이브리드 상태** - 일부는 vm.Config, 일부는 개별 필드

**사용 현황**:
```go
// register_task.go (방금 생성) - Optimism 스타일
cfg.Cannon  // ← vm.Config 사용

// cannon/executor_test.go - 구버전 스타일
cfg.CannonBin
cfg.CannonServer
cfg.CannonSnapshotFreq
```

**결론**:
- ✅ `cfg.Cannon vm.Config` 이미 존재
- ⚠️ 하지만 개별 필드도 여전히 사용 중
- 📝 Check() 함수에서 개별 필드 → vm.Config 복사

---

### 2. GetGameRange vs GetBlockRange ✅ 조사 완료

**Tokamak-Thanos**:
```go
// contracts/faultdisputegame.go:155-157
func (f *FaultDisputeGameContractLatest) GetBlockRange(ctx context.Context) (prestateBlock uint64, poststateBlock uint64, retErr error)

// contracts/faultdisputegame.go:596
type FaultDisputeGameContract interface {
    GetBlockRange(ctx context.Context) (prestateBlock uint64, poststateBlock uint64, retErr error)
}
```

**Optimism**:
```go
// contracts/superfaultdisputegame.go:78-79
func (f *SuperFaultDisputeGameContractLatest) GetGameRange(ctx context.Context) (prestateBlock uint64, poststateBlock uint64, retErr error)
```

**결론**:
- ❌ **메서드 이름 다름**: Tokamak은 `GetBlockRange()`, Optimism은 `GetGameRange()`
- ⚠️ Optimism의 `GetGameRange()`는 **Super 전용**일 수 있음
- 📝 **해결**: register_task.go에서 `GetGameRange()` → `GetBlockRange()` 사용

---

### 3. prestates/source.go 존재 여부 ✅ 조사 완료

**Tokamak-Thanos**:
```bash
prestates/
├── cache.go       ✅
├── multi.go       ✅
├── single.go      ✅
└── source.go      ❌ 없음!
```

**Optimism**:
```bash
prestates/
├── cache.go       ✅
├── multi.go       ✅
├── single.go      ✅
└── source.go      ✅ 있음!
```

**결론**:
- ❌ **Tokamak에 source.go 없음**
- 📝 **NewPrestateSource 함수 없음** → 추가 필요

---

### 4. dial 패키지 시그니처 ✅ 조사 완료

**Tokamak-Thanos**:
```go
// op-service/dial/dial.go:38
func DialRollupClientWithTimeout(
    ctx context.Context,
    timeout time.Duration,  // ← timeout 필요!
    log log.Logger,
    url string,
) (*sources.RollupClient, error)
```

**Optimism**:
```go
// op-service/dial/dial.go:50
func DialRollupClientWithTimeout(
    ctx context.Context,
    log log.Logger,  // ← timeout 없음!
    url string,
    callerOpts ...client.RPCOption,
) (*sources.RollupClient, error)
```

**결론**:
- ❌ **시그니처 완전히 다름**
- Tokamak: timeout 필수
- Optimism: callerOpts 선택적
- 📝 **해결**: clients.go에서 timeout 추가

---

### 5. NewPrestateProviderCache 시그니처 차이 ✅ 조사 완료

**Tokamak-Thanos**:
```go
// prestates/cache.go:21
func NewPrestateProviderCache(
    m caching.Metrics,
    label string,
    createProvider func(prestateHash common.Hash) (types.PrestateProvider, error),  // ← ctx 없음
) *PrestateProviderCache

func (p *PrestateProviderCache) GetOrCreate(prestateHash common.Hash) (types.PrestateProvider, error)  // ← ctx 없음
```

**Optimism**:
```go
// prestates/cache.go:23
func NewPrestateProviderCache(
    m caching.Metrics,
    label string,
    createProvider func(ctx context.Context, prestateHash common.Hash) (types.PrestateProvider, error),  // ← ctx 있음
) *PrestateProviderCache

func (p *PrestateProviderCache) GetOrCreate(ctx context.Context, prestateHash common.Hash) (types.PrestateProvider, error)  // ← ctx 있음
```

**결론**:
- ❌ **ctx 파라미터 차이**
- 📝 **전체 prestates 패키지 업데이트 필요**

---

## 📋 상세 마이그레이션 계획 (전략 A)

### 전체 작업 목록

| 단계 | 파일 | 작업 내용 | 예상 시간 | 우선순위 |
|-----|------|----------|---------|---------|
| **1단계** | `prestates/source.go` | Optimism에서 복사 | 0.5시간 | 🔴 최우선 |
| **2단계** | `prestates/cache.go` | ctx 파라미터 추가 | 1시간 | 🔴 최우선 |
| **3단계** | `prestates/single.go` | PrestatePath ctx 추가 | 0.5시간 | 🔴 최우선 |
| **4단계** | `prestates/multi.go` | PrestatePath ctx 추가 | 0.5시간 | 🔴 최우선 |
| **5단계** | `clients.go` | timeout 파라미터 수정 | 0.5시간 | 🔴 최우선 |
| **6단계** | `register_task.go` | GetGameRange → GetBlockRange | 0.5시간 | 🔴 최우선 |
| **7단계** | `register.go` | PrestateSource ctx 추가 | 0.5시간 | 🔴 최우선 |
| **8단계** | **테스트** | **모든 변경사항 테스트** | **2-3시간** | 🔴 **최우선** |
| **총계** | | | **6-8시간 (1일)** | |

### 단계별 상세 작업

#### 🔴 단계 1: prestates/source.go 추가 (0.5시간)

**작업**:
```bash
# Optimism에서 복사
cp /Users/zena/tokamak-projects/optimism/op-challenger/game/fault/trace/prestates/source.go \
   /Users/zena/tokamak-projects/tokamak-thanos/op-challenger/game/fault/trace/prestates/

# Import 경로 수정
sed -i '' 's/ethereum-optimism\/optimism/tokamak-network\/tokamak-thanos/g' \
    op-challenger/game/fault/trace/prestates/source.go
```

**검증**:
```bash
go build ./op-challenger/game/fault/trace/prestates
```

---

#### 🔴 단계 2: prestates/cache.go 업데이트 (1시간)

**변경**:
```go
// 변경 전
func NewPrestateProviderCache(
    m caching.Metrics,
    label string,
    createProvider func(prestateHash common.Hash) (types.PrestateProvider, error),
) *PrestateProviderCache

// 변경 후
func NewPrestateProviderCache(
    m caching.Metrics,
    label string,
    createProvider func(ctx context.Context, prestateHash common.Hash) (types.PrestateProvider, error),
) *PrestateProviderCache
```

**영향 받는 함수**:
- `GetOrCreate(prestateHash)` → `GetOrCreate(ctx, prestateHash)`

---

#### 🔴 단계 3-4: single.go, multi.go 업데이트 (1시간)

**single.go**:
```go
// 변경 전
func (s *SinglePrestateSource) PrestatePath(_ common.Hash) (string, error) {
    return s.path, nil
}

// 변경 후
func (s *SinglePrestateSource) PrestatePath(_ context.Context, _ common.Hash) (string, error) {
    return s.path, nil
}
```

**multi.go**:
```go
// 변경 전
func (m *MultiPrestateProvider) PrestatePath(hash common.Hash) (string, error) {
    // ...
}

// 변경 후
func (m *MultiPrestateProvider) PrestatePath(ctx context.Context, hash common.Hash) (string, error) {
    // ...
}
```

---

#### 🔴 단계 5: clients.go 수정 (0.5시간)

**변경**:
```go
// 변경 전
rollupClient, err := dial.DialRollupClientWithTimeout(c.ctx, c.logger, c.cfg.RollupRpc)

// 변경 후 (Tokamak 버전에 맞게)
rollupClient, err := dial.DialRollupClientWithTimeout(c.ctx, 30*time.Second, c.logger, c.cfg.RollupRpc)
```

---

#### 🔴 단계 6: register_task.go 수정 (0.5시간)

**변경**:
```go
// Line 237: GetGameRange → GetBlockRange
prestateBlock, poststateBlock, err := contract.GetBlockRange(ctx)
```

---

#### 🔴 단계 7: register.go의 PrestateSource 인터페이스 (0.5시간)

**변경**:
```go
// 변경 전
type PrestateSource interface {
    PrestatePath(prestateHash common.Hash) (string, error)
}

// 변경 후
type PrestateSource interface {
    PrestatePath(ctx context.Context, prestateHash common.Hash) (string, error)
}
```

---

### 🧪 단계 8: 테스트 계획 (2-3시간)

#### 8.1 Unit 테스트
```bash
# Prestates 패키지
go test ./op-challenger/game/fault/trace/prestates/... -v

# Register 패키지
go test ./op-challenger/game/fault/... -run TestRegister -v

# 전체
go test ./op-challenger/... -short
```

#### 8.2 빌드 테스트
```bash
go build ./op-challenger/cmd/op-challenger
```

#### 8.3 회귀 테스트
- GameType 0 (Cannon) 정상 작동 확인
- GameType 2 (Asterisc) 정상 작동 확인
- GameType 255 (Alphabet) 정상 작동 확인

---

## ✅ 성공 기준

### 필수 (단계 1-7)
- [ ] prestates/source.go 추가 완료
- [ ] prestates 패키지 ctx 파라미터 지원
- [ ] clients.go 빌드 성공
- [ ] register_task.go 빌드 성공
- [ ] register.go 빌드 성공
- [ ] Linter 에러 0개

### 필수 (단계 8)
- [ ] 모든 unit 테스트 통과
- [ ] op-challenger 바이너리 빌드 성공
- [ ] GameType 0, 2 회귀 테스트 통과

### 최종
- [ ] GameType 3 코드 추가 가능 상태
- [ ] Optimism 구조 완전 적용

---

## 🚀 실행 계획

### 오늘 (1일차)

**오전** (4시간):
- 단계 1-4: prestates 패키지 완전 업데이트
- 중간 테스트

**오후** (4시간):
- 단계 5-7: clients.go, register_task.go, register.go 수정
- 단계 8: 전체 테스트

**목표**: 🎯 Optimism 구조 완전 적용 완료

---

### 내일 (2일차)

**GameType 3 추가**:
- NewAsteriscKonaRegisterTask 최종 검증
- CLI 플래그 추가
- 전체 통합 테스트

**목표**: 🎯 GameType 3 완전 작동

---

## 📊 요약

### 발견한 주요 차이점 (5개)

1. ❌ **GameType 타입**: `uint32` vs `faultTypes.GameType` → 보류 (나중에)
2. ❌ **Prestates ctx**: 없음 vs 있음 → **수정 필요**
3. ❌ **source.go**: 없음 vs 있음 → **추가 필요**
4. ❌ **Dial timeout**: 필수 vs 선택 → **수정 필요**
5. ⚠️ **GetGameRange**: GetBlockRange vs GetGameRange → **메서드명만 변경**

### 최소 작업 (GameType 3만 추가)

**필수 수정 (1일)**:
- prestates/source.go 추가
- prestates 패키지 ctx 지원
- clients.go timeout 추가
- register_task.go GetBlockRange 사용

**나중에 (선택)**:
- GameType 타입 변경 (큰 작업)

---

**마지막 업데이트**: 2025-10-26
**상태**: ✅ 조사 완료
**다음**: 단계 1-7 실행 (예상 6-8시간)

