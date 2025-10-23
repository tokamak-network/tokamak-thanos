# Dispute Game 타입과 VM 매핑

## 목차
1. [개요](#개요)
2. [게임 타입 정의](#게임-타입-정의)
3. [지원하는 VM](#지원하는-vm)
4. [게임 타입별 상세 설명](#게임-타입별-상세-설명)
5. [VM별 TraceAccessor](#vm별-traceaccessor)
6. [등록 및 설정](#등록-및-설정)
7. [실전 사용 예시](#실전-사용-예시)
8. [FAQ](#faq)

---

## 개요

op-challenger의 fault proof 시스템에서는 **게임 타입(Game Type)**에 따라 서로 다른 **가상 머신(VM)**을 사용합니다. 이는 Ethereum의 클라이언트 다양성(Client Diversity) 원칙과 유사하게, 여러 VM 구현체를 지원함으로써 시스템의 안정성과 보안을 강화합니다.

### 주요 특징

- **다중 VM 지원**: MIPS, RISC-V 등 여러 아키텍처 지원
- **게임 타입별 격리**: 각 게임 타입은 독립적인 PlayerCreator와 TraceAccessor를 가짐
- **동적 등록**: 설정에 따라 원하는 게임 타입만 활성화 가능
- **확장 가능성**: 새로운 VM 추가가 용이한 구조

---

## 게임 타입 정의

게임 타입은 `uint32` 값으로 정의되며, 온체인 `DisputeGameFactory`에서 게임 생성 시 사용됩니다.

```go
// op-challenger/game/fault/types/types.go

const (
    CannonGameType       uint32 = 0    // MIPS VM (프로덕션)
    PermissionedGameType uint32 = 1    // MIPS VM (권한 기반)
    AsteriscGameType     uint32 = 2    // RISC-V VM
    FastGameType         uint32 = 254  // Alphabet (빠른 테스트)
    AlphabetGameType     uint32 = 255  // Alphabet (테스트)
)
```

### 게임 타입 구조

```go
type GameMetadata struct {
    Index     uint64          // 게임 인덱스
    GameType  uint32          // 게임 타입 (위의 상수 중 하나)
    Timestamp uint64          // 생성 타임스탬프
    Proxy     common.Address  // 게임 컨트랙트 주소
}
```

---

## 지원하는 VM

### 1. MIPS VM (Cannon)

**아키텍처**: MIPS (32-bit)

**특징**:
- 온체인 검증: `MIPS.sol` 스마트 컨트랙트
- 단일 instruction 실행을 온체인에서 재현
- 결정론적 실행 보장
- 프로덕션 환경에서 검증된 안정성

**실행 파일**:
- `cannon`: MIPS VM 에뮬레이터
- `op-program`: fault proof 프로그램 (MIPS 바이너리)

**prestate**:
- `prestate.json`: MIPS VM의 초기 상태
- 약 190MB 크기의 VM 메모리 스냅샷

### 2. RISC-V VM (Asterisc)

**아키텍처**: RISC-V (64-bit)

**특징**:
- 온체인 검증: `RiscV.sol` 스마트 컨트랙트
- RISC-V는 오픈소스 ISA (Instruction Set Architecture)
- Cannon의 대안으로 클라이언트 다양성 제공
- 동일한 op-program을 RISC-V로 컴파일하여 실행

**실행 파일**:
- `asterisc`: RISC-V VM 에뮬레이터
- `op-program-rv64`: fault proof 프로그램 (RISC-V 바이너리)

**prestate**:
- `prestate-rv64.json`: RISC-V VM의 초기 상태

### 3. Alphabet Provider (테스트용)

**타입**: 간소화된 trace provider

**특징**:
- 실제 VM 없이 L2 output root만 검증
- 빠른 증명 생성 (VM 실행 오버헤드 없음)
- 개발 및 테스트 환경 전용

**사용 케이스**:
- 로컬 devnet 테스트
- Fault proof 로직 검증
- CI/CD 파이프라인

---

## 게임 타입별 상세 설명

### GameType 0: Cannon (MIPS)

**용도**: 프로덕션 fault proof

**등록 코드**: `op-challenger/game/fault/register.go:76-79`

```go
if cfg.TraceTypeEnabled(config.TraceTypeCannon) {
    registerCannon(faultTypes.CannonGameType, registry, oracles, ...)
}
```

**TraceAccessor 생성**: `register.go:347-357`

```go
creator := func(ctx context.Context, logger log.Logger, gameDepth faultTypes.Depth, dir string) (faultTypes.TraceAccessor, error) {
    cannonPrestate, err := prestateSource.PrestatePath(requiredPrestatehash)
    if err != nil {
        return nil, fmt.Errorf("failed to get cannon prestate: %w", err)
    }
    accessor, err := outputs.NewOutputCannonTraceAccessor(
        logger, m, cfg, l2Client, prestateProvider,
        cannonPrestate, rollupClient, dir, l1HeadID,
        splitDepth, prestateBlock, poststateBlock
    )
    return accessor, nil
}
```

**주요 설정**:
```bash
--trace-type cannon
--cannon-bin ./cannon/bin/cannon
--cannon-server ./op-program/bin/op-program
--cannon-prestate ./op-program/bin/prestate.json
--cannon-rollup-config ./.devnet/rollup.json
--cannon-l2-genesis ./.devnet/genesis-l2.json
--cannon-snapshot-freq 1000000000  # 1B instructions
```

**온체인 검증 컨트랙트**:
- `MIPS.sol`: 단일 MIPS instruction 실행
- `PreimageOracle.sol`: Preimage 데이터 제공

---

### GameType 1: Permissioned Cannon

**용도**: 권한 기반 프라이빗 네트워크

**차이점**:
- Cannon과 동일한 MIPS VM 사용
- 특정 주소만 게임 참여 가능
- 초기 테스트 단계 또는 프라이빗 롤업용

**등록 코드**: `register.go:81-84`

```go
if cfg.TraceTypeEnabled(config.TraceTypePermissioned) {
    registerCannon(faultTypes.PermissionedGameType, registry, oracles, ...)
}
```

**설정**:
```bash
--trace-type permissioned
# Cannon과 동일한 설정 사용
```

---

### GameType 2: Asterisc (RISC-V)

**용도**: Cannon의 대안, 클라이언트 다양성

**등록 코드**: `register.go:86-89`

```go
if cfg.TraceTypeEnabled(config.TraceTypeAsterisc) {
    registerAsterisc(faultTypes.AsteriscGameType, registry, oracles, ...)
}
```

**TraceAccessor 생성**: `register.go:252-262`

```go
creator := func(ctx context.Context, logger log.Logger, gameDepth faultTypes.Depth, dir string) (faultTypes.TraceAccessor, error) {
    asteriscPrestate, err := prestateSource.PrestatePath(requiredPrestatehash)
    if err != nil {
        return nil, fmt.Errorf("failed to get asterisc prestate: %w", err)
    }
    accessor, err := outputs.NewOutputAsteriscTraceAccessor(
        logger, m, cfg, l2Client, prestateProvider,
        asteriscPrestate, rollupClient, dir, l1HeadID,
        splitDepth, prestateBlock, poststateBlock
    )
    return accessor, nil
}
```

**주요 설정**:
```bash
--trace-type asterisc
--asterisc-bin ./asterisc/bin/asterisc
--asterisc-server ./op-program/bin/op-program-rv64
--asterisc-prestate ./op-program/bin/prestate-rv64.json
--asterisc-rollup-config ./.devnet/rollup.json
--asterisc-l2-genesis ./.devnet/genesis-l2.json
--asterisc-snapshot-freq 1000000000
```

**온체인 검증 컨트랙트**:
- `RiscV.sol`: 단일 RISC-V instruction 실행
- `PreimageOracle.sol`: Cannon과 동일

**RISC-V vs MIPS**:

| 항목 | MIPS | RISC-V |
|-----|------|--------|
| 아키텍처 | 32-bit | 64-bit |
| ISA | 독점적 | 오픈소스 |
| 온체인 가스 비용 | ~5M gas | ~5M gas (유사) |
| 성숙도 | 높음 (검증됨) | 중간 (활발히 개발 중) |
| 생태계 | 제한적 | 활발 (RISC-V Foundation) |

---

### GameType 254: Fast (Alphabet)

**용도**: 빠른 테스트 및 개발

**등록 코드**: `register.go:91-94`

```go
if cfg.TraceTypeEnabled(config.TraceTypeFast) {
    registerAlphabet(faultTypes.FastGameType, registry, oracles, ...)
}
```

**TraceAccessor 생성**: `register.go:146-152`

```go
creator := func(ctx context.Context, logger log.Logger, gameDepth faultTypes.Depth, dir string) (faultTypes.TraceAccessor, error) {
    accessor, err := outputs.NewOutputAlphabetTraceAccessor(
        logger, m, prestateProvider, rollupClient, l2Client,
        l1Head, splitDepth, prestateBlock, poststateBlock
    )
    return accessor, nil
}
```

**특징**:
- VM 실행 없음 (L2 output만 검증)
- 증명 생성 시간: 수 초
- 프로덕션 사용 불가

**설정**:
```bash
--trace-type fast
# VM 관련 설정 불필요
```

---

### GameType 255: Alphabet (테스트)

**용도**: Fault proof 로직 테스트

**등록 코드**: `register.go:96-99`

```go
if cfg.TraceTypeEnabled(config.TraceTypeAlphabet) {
    registerAlphabet(faultTypes.AlphabetGameType, registry, oracles, ...)
}
```

**동작 방식**:
```go
// alphabet/provider.go
type AlphabetTraceProvider struct {
    state string  // 예: "abcdefgh"
}

func (p *AlphabetTraceProvider) Get(ctx, pos Position) (common.Hash, error) {
    index := pos.TraceIndex()
    letter := p.state[index]
    return common.BytesToHash([]byte{letter}), nil
}
```

**사용 예시**:
- 게임 트리 순회 로직 테스트
- Bisection 알고리즘 검증
- Claim/Counter-claim 시뮬레이션

---

## VM별 TraceAccessor

### TraceAccessor 인터페이스

```go
// op-challenger/game/fault/types/types.go

type TraceAccessor interface {
    // 특정 position의 상태 해시 반환
    Get(ctx context.Context, pos Position) (common.Hash, error)

    // Step 증명 데이터 반환 (리프 노드)
    GetStepData(ctx context.Context, pos Position) (
        prestate []byte,
        proofData []byte,
        oracleData *PreimageOracleData,
        error
    )

    // L2 block number 챌린지 데이터
    GetL2BlockNumberChallenge(ctx context.Context) (*InvalidL2BlockNumberChallenge, error)
}
```

### 게임 타입별 TraceAccessor 구현

| 게임 타입 | TraceAccessor 구현 | VM 실행 | 위치 |
|----------|-------------------|---------|------|
| Cannon (0) | `OutputCannonTraceAccessor` | MIPS | `trace/outputs/output_cannon.go` |
| Permissioned (1) | `OutputCannonTraceAccessor` | MIPS | `trace/outputs/output_cannon.go` |
| Asterisc (2) | `OutputAsteriscTraceAccessor` | RISC-V | `trace/outputs/output_asterisc.go` |
| Fast (254) | `OutputAlphabetTraceAccessor` | - | `trace/outputs/output_alphabet.go` |
| Alphabet (255) | `OutputAlphabetTraceAccessor` | - | `trace/outputs/output_alphabet.go` |

### Split Depth 개념

모든 TraceAccessor는 **split depth**를 사용합니다:

```
Depth < splitDepth: L2 output root 검증
Depth >= splitDepth: VM 실행 검증
```

**예시** (splitDepth = 30, maxDepth = 73):

```
Position Depth 0-29:  L2 output 검증
    ├─ rollupClient.OutputAtBlock(blockNum)
    └─ 빠른 응답 (RPC 호출만)

Position Depth 30-73: VM 실행 검증
    ├─ Cannon 또는 Asterisc 실행
    └─ 느린 응답 (VM 시뮬레이션)
```

---

## 등록 및 설정

### GameTypeRegistry 구조

```go
// op-challenger/game/registry/registry.go

type GameTypeRegistry struct {
    types        map[uint32]scheduler.PlayerCreator      // 게임 타입 → 플레이어 생성자
    bondCreators map[uint32]claims.BondContractCreator  // 게임 타입 → 본드 컨트랙트 생성자
}

func (r *GameTypeRegistry) RegisterGameType(gameType uint32, creator scheduler.PlayerCreator) {
    if _, ok := r.types[gameType]; ok {
        panic(fmt.Errorf("duplicate creator registered for game type: %v", gameType))
    }
    r.types[gameType] = creator
}

func (r *GameTypeRegistry) CreatePlayer(game types.GameMetadata, dir string) (scheduler.GamePlayer, error) {
    creator, ok := r.types[game.GameType]
    if !ok {
        return nil, fmt.Errorf("%w: %v", ErrUnsupportedGameType, game.GameType)
    }
    return creator(game, dir)
}
```

### 등록 흐름

```
1. Service 초기화 (game/service.go:88-129)
   └─> RegisterGameTypes() 호출

2. RegisterGameTypes() (game/fault/register.go:53-102)
   ├─ cfg.TraceTypeEnabled() 확인
   ├─ registerCannon() 또는 registerAsterisc() 호출
   └─ registry.RegisterGameType(gameType, playerCreator)

3. PlayerCreator 함수 생성
   ├─ Contract 바인딩
   ├─ Prestate 로드
   ├─ TraceAccessor 생성자 정의
   └─ GamePlayer 반환

4. 게임 발견 시 (scheduler/coordinator.go)
   ├─ registry.CreatePlayer(game, dir) 호출
   └─ 적절한 PlayerCreator로 GamePlayer 생성
```

### TraceType vs GameType 매핑

**설정 레벨** (CLI 플래그):
```bash
--trace-type cannon,asterisc  # TraceType 설정
```

**코드 매핑**:
```go
// config/config.go:53-63
type TraceType string

const (
    TraceTypeAlphabet     TraceType = "alphabet"      // → GameType 255
    TraceTypeFast         TraceType = "fast"          // → GameType 254
    TraceTypeCannon       TraceType = "cannon"        // → GameType 0
    TraceTypeAsterisc     TraceType = "asterisc"      // → GameType 2
    TraceTypePermissioned TraceType = "permissioned"  // → GameType 1
)
```

**등록 시 변환**:
```go
// register.go:76-99
if cfg.TraceTypeEnabled(config.TraceTypeCannon) {
    registerCannon(faultTypes.CannonGameType, ...)  // "cannon" → 0
}
if cfg.TraceTypeEnabled(config.TraceTypeAsterisc) {
    registerAsterisc(faultTypes.AsteriscGameType, ...)  // "asterisc" → 2
}
```

---

## 실전 사용 예시

### 1. 로컬 Devnet (Cannon만 사용)

```bash
#!/bin/bash

DISPUTE_GAME_FACTORY=$(jq -r .DisputeGameFactoryProxy .devnet/addresses.json)

./op-challenger/bin/op-challenger \
  --trace-type cannon \
  --l1-eth-rpc http://localhost:8545 \
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

**결과**:
- GameType 0 (Cannon) 게임만 처리
- GameType 2 (Asterisc) 게임은 `ErrUnsupportedGameType` 에러

---

### 2. 프로덕션 (Cannon + Asterisc 다중 지원)

```bash
#!/bin/bash

./op-challenger/bin/op-challenger \
  --trace-type cannon,asterisc \
  --l1-eth-rpc https://mainnet.infura.io/v3/YOUR_KEY \
  --l1-beacon https://beacon-api.example.com \
  --rollup-rpc https://rollup-node.example.com \
  --l2-eth-rpc https://l2-node.example.com \
  --game-factory-address 0x1234... \
  --datadir /data/challenger \
  \
  --cannon-bin /usr/local/bin/cannon \
  --cannon-server /usr/local/bin/op-program \
  --cannon-prestate /data/prestates/prestate.json \
  --cannon-rollup-config /config/rollup.json \
  --cannon-l2-genesis /config/genesis-l2.json \
  --cannon-snapshot-freq 1000000000 \
  \
  --asterisc-bin /usr/local/bin/asterisc \
  --asterisc-server /usr/local/bin/op-program-rv64 \
  --asterisc-prestate /data/prestates/prestate-rv64.json \
  --asterisc-rollup-config /config/rollup.json \
  --asterisc-l2-genesis /config/genesis-l2.json \
  --asterisc-snapshot-freq 1000000000 \
  \
  --private-key $PRIVATE_KEY \
  --max-concurrency 4 \
  --max-pending-tx 10 \
  --game-window 672h \
  --metrics.enabled \
  --metrics.port 7300
```

**결과**:
- GameType 0 (Cannon) 게임 처리 가능
- GameType 2 (Asterisc) 게임 처리 가능
- GameType 1, 254, 255 게임은 무시됨

---

### 3. 빠른 테스트 (Fast 모드)

```bash
#!/bin/bash

./op-challenger/bin/op-challenger \
  --trace-type fast \
  --l1-eth-rpc http://localhost:8545 \
  --rollup-rpc http://localhost:9546 \
  --l2-eth-rpc http://localhost:9545 \
  --game-factory-address $DISPUTE_GAME_FACTORY \
  --datadir temp/fast-challenger \
  --mnemonic "test test test test test test test test test test test junk" \
  --hd-path "m/44'/60'/0'/0/8"
```

**특징**:
- VM 바이너리 불필요
- 증명 생성이 매우 빠름 (수 초)
- CI/CD 테스트에 적합

---

### 4. 특정 게임만 처리 (Allowlist)

```bash
#!/bin/bash

./op-challenger/bin/op-challenger \
  --trace-type cannon,asterisc \
  --game-allowlist 0xGameAddress1,0xGameAddress2 \
  --l1-eth-rpc http://localhost:8545 \
  --rollup-rpc http://localhost:9546 \
  --game-factory-address $DISPUTE_GAME_FACTORY \
  # ... 나머지 설정
```

**용도**:
- 특정 게임만 선택적으로 모니터링
- 테스트 게임 격리
- 리소스 제한 환경

---

### 5. 다중 Prestate 지원 (URL 기반)

```bash
#!/bin/bash

./op-challenger/bin/op-challenger \
  --trace-type cannon \
  --cannon-prestate-base-url https://prestates.example.com/cannon/ \
  --l1-eth-rpc http://localhost:8545 \
  --rollup-rpc http://localhost:9546 \
  --game-factory-address $DISPUTE_GAME_FACTORY \
  # ... 나머지 설정
```

**동작**:
```go
// register.go:300-305
if cfg.CannonAbsolutePreStateBaseURL != nil {
    prestateSource = prestates.NewMultiPrestateProvider(
        cfg.CannonAbsolutePreStateBaseURL,
        filepath.Join(cfg.Datadir, "cannon-prestates")
    )
}
```

**Prestate URL 구조**:
```
https://prestates.example.com/cannon/
  ├─ 0xabcd1234.json  # Prestate hash → 파일명
  ├─ 0x5678efab.json
  └─ ...
```

**장점**:
- 여러 prestate 버전 지원
- 자동 다운로드 및 캐싱
- 네트워크 업그레이드 시 유연성

---

## 게임 타입별 성능 비교

### 증명 생성 시간

| 게임 타입 | VM | Trace 생성 시간 | Step 증명 시간 | 총 시간 (추정) |
|----------|-----|---------------|--------------|--------------|
| Cannon (0) | MIPS | 5-30분 | 수 초 | 5-30분 |
| Asterisc (2) | RISC-V | 5-30분 | 수 초 | 5-30분 |
| Fast (254) | - | 수 초 | - | 수 초 |
| Alphabet (255) | - | 즉시 | - | 즉시 |

*실제 시간은 하드웨어, L2 블록 범위, snapshot 가용성에 따라 다름*

### 디스크 사용량

| 게임 타입 | Prestate 크기 | Proof 캐시 (게임당) | 총 예상 |
|----------|--------------|-------------------|---------|
| Cannon (0) | ~190MB | 1-10GB | ~200GB (100개 게임) |
| Asterisc (2) | ~200MB | 1-10GB | ~200GB (100개 게임) |
| Fast (254) | - | <1MB | ~1MB |
| Alphabet (255) | - | <1KB | <1KB |

### 온체인 검증 비용

| 게임 타입 | Step 가스 비용 | Resolve 가스 비용 |
|----------|---------------|-----------------|
| Cannon (0) | ~5M gas | ~500K gas |
| Asterisc (2) | ~5M gas | ~500K gas |
| Fast (254) | N/A | ~300K gas |
| Alphabet (255) | N/A | ~300K gas |

---

## FAQ

### Q1: 왜 여러 VM을 지원하나요?

**A**: 세 가지 주요 이유가 있습니다:

1. **클라이언트 다양성 (Client Diversity)**
   - Ethereum처럼 여러 구현체 지원
   - 한 VM에 버그가 있어도 다른 VM으로 검증 가능
   - 단일 장애점(Single Point of Failure) 제거

2. **보안 강화**
   - 서로 다른 아키텍처로 동일한 결과 검증
   - 구현 버그 발견 확률 증가
   - 공격 표면 분산

3. **유연성**
   - 네트워크마다 선호하는 VM 선택 가능
   - 새로운 VM 추가 용이
   - 점진적 마이그레이션 가능

---

### Q2: Cannon과 Asterisc 중 어떤 것을 사용해야 하나요?

**A**: 상황에 따라 다릅니다:

**Cannon (MIPS) 추천**:
- 프로덕션 환경 (검증된 안정성)
- 보수적 접근
- 현재 Optimism 메인넷과 동일한 환경

**Asterisc (RISC-V) 추천**:
- 클라이언트 다양성 확보
- RISC-V 생태계 활용
- 미래 지향적 선택

**둘 다 사용** (권장):
- `--trace-type cannon,asterisc`
- 최대 보안 및 안정성
- 프로덕션 환경에서 권장

---

### Q3: 새로운 게임 타입이 추가되면 어떻게 하나요?

**A**: 다음 단계를 따릅니다:

1. **Prestate 업데이트**
   ```bash
   # 새 prestate 다운로드
   wget https://example.com/new-prestate.json -O /data/prestates/new-prestate.json
   ```

2. **설정 변경 없음**
   - 기존 `--trace-type cannon` 설정 유지
   - Multi-prestate 모드라면 자동 감지

3. **재시작**
   ```bash
   systemctl restart op-challenger
   ```

4. **확인**
   ```bash
   # 로그에서 새 prestate 로드 확인
   journalctl -u op-challenger -f | grep "prestate"
   ```

---

### Q4: GameType과 TraceType의 차이는 무엇인가요?

**A**:

**GameType** (온체인):
- `uint32` 값 (0, 1, 2, 254, 255)
- 스마트 컨트랙트에서 사용
- `DisputeGameFactory.create(gameType, ...)`

**TraceType** (오프체인):
- 문자열 값 ("cannon", "asterisc", ...)
- CLI 설정에서 사용
- `--trace-type cannon,asterisc`

**매핑**:
```
TraceType "cannon"       → GameType 0
TraceType "permissioned" → GameType 1
TraceType "asterisc"     → GameType 2
TraceType "fast"         → GameType 254
TraceType "alphabet"     → GameType 255
```

---

### Q5: 지원하지 않는 게임 타입을 만나면 어떻게 되나요?

**A**:

```go
// scheduler/coordinator.go
func (r *GameTypeRegistry) CreatePlayer(game types.GameMetadata, dir string) (scheduler.GamePlayer, error) {
    creator, ok := r.types[game.GameType]
    if !ok {
        return nil, fmt.Errorf("%w: %v", ErrUnsupportedGameType, game.GameType)
    }
    return creator(game, dir)
}
```

**결과**:
- 해당 게임은 스킵됨
- 로그에 경고 출력
- 다른 게임은 정상 처리
- 프로세스는 계속 실행

**로그 예시**:
```
WARN [01-15|12:00:00] Unsupported game type game=0x1234... type=99
```

---

### Q6: VM 바이너리는 어디서 빌드하나요?

**A**:

**Cannon**:
```bash
cd cannon
make cannon
# 출력: ./bin/cannon
```

**op-program (MIPS)**:
```bash
cd op-program
make op-program
# 출력: ./bin/op-program
```

**Asterisc**:
```bash
cd asterisc
make asterisc
# 출력: ./bin/asterisc
```

**op-program (RISC-V)**:
```bash
cd op-program
make op-program-rv64
# 출력: ./bin/op-program-rv64
```

**Prestate 생성**:
```bash
make cannon-prestate
# 출력: op-program/bin/prestate.json

make asterisc-prestate
# 출력: op-program/bin/prestate-rv64.json
```

---

### Q7: Split Depth는 어떻게 결정되나요?

**A**: 온체인 게임 컨트랙트에서 설정됩니다:

```solidity
// FaultDisputeGame.sol
uint256 public immutable SPLIT_DEPTH;

constructor(uint256 _splitDepth, ...) {
    SPLIT_DEPTH = _splitDepth;
}
```

**일반적인 값**:
- Mainnet: `splitDepth = 30`
- Testnet: `splitDepth = 30`
- Devnet: `splitDepth = 14` (빠른 테스트)

**의미**:
```
maxDepth = 73 (예시)
splitDepth = 30

Depth 0-29:  2^30 = 1,073,741,824 blocks 범위
             → L2 output 검증 (빠름)

Depth 30-73: 2^43 = 8,796,093,022,208 instructions 범위
             → VM 실행 검증 (느림)
```

**트레이드오프**:
- Split depth ↑: L2 output 검증 구간 확대 → 빠름
- Split depth ↓: VM 검증 구간 확대 → 느림, 정확함

---

### Q8: 여러 challenger가 다른 VM을 사용할 수 있나요?

**A**: 예, 가능하며 권장됩니다!

**시나리오**:

```
Challenger A: --trace-type cannon
Challenger B: --trace-type asterisc
Challenger C: --trace-type cannon,asterisc
```

**결과**:
- 각 challenger는 독립적으로 동작
- 같은 게임에 대해 서로 다른 VM으로 검증
- 최종 온체인 결과는 동일해야 함
- 결과가 다르면 VM 구현 버그 의심

**장점**:
- 클라이언트 다양성 극대화
- 상호 검증 (Cross-validation)
- 버그 발견 가능성 증가

---

## 관련 파일

### 핵심 코드

- **게임 타입 정의**: `op-challenger/game/fault/types/types.go`
- **등록 로직**: `op-challenger/game/fault/register.go`
- **레지스트리**: `op-challenger/game/registry/registry.go`
- **TraceType 설정**: `op-challenger/config/config.go`

### TraceAccessor 구현

- **Cannon**: `op-challenger/game/fault/trace/cannon/provider.go`
- **Asterisc**: `op-challenger/game/fault/trace/asterisc/provider.go`
- **Alphabet**: `op-challenger/game/fault/trace/alphabet/provider.go`
- **Output 래퍼**: `op-challenger/game/fault/trace/outputs/`

### VM 바이너리

- **Cannon**: `cannon/`
- **Asterisc**: `asterisc/`
- **op-program**: `op-program/`

### 스마트 컨트랙트

- **MIPS.sol**: `packages/contracts-bedrock/src/cannon/MIPS.sol`
- **RiscV.sol**: `packages/contracts-bedrock/src/cannon/MIPS2.sol` (예정)
- **FaultDisputeGame.sol**: `packages/contracts-bedrock/src/dispute/FaultDisputeGame.sol`
- **DisputeGameFactory.sol**: `packages/contracts-bedrock/src/dispute/DisputeGameFactory.sol`

---

## 참고 자료

### 공식 문서

- [Optimism Fault Proof Specs](https://specs.optimism.io/experimental/fault-proof/)
- [Cannon Documentation](../../cannon/README.md)
- [Asterisc Documentation](../../asterisc/README.md)
- [op-program Documentation](../../op-program/README.md)

### 관련 문서

- [op-challenger 아키텍처](./op-challenger-architecture-ko.md)
- [배포 가이드](./deployment-guide-ko.md)

### 외부 참고

- [RISC-V Foundation](https://riscv.org/)
- [MIPS Architecture](https://www.mips.com/)
- [Ethereum Client Diversity](https://ethereum.org/en/developers/docs/nodes-and-clients/client-diversity/)

---

**마지막 업데이트**: 2025-01-21
**버전**: v1.7.7 기준
