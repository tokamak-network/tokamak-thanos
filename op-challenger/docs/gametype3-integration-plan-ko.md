# GameType 3 (ASTERISC_KONA) 통합 계획

## 개요

이 문서는 Optimism의 **GameType 3 (ASTERISC_KONA)**를 Tokamak-Thanos에 통합하기 위한 실제 코드 분석 기반의 상세 계획입니다.

**분석 기준**: Optimism v2.0.0-alpha.8 (실제 코드)

### 🎯 Tokamak-Thanos 현재 구현 상태

**진행률**: ✅ **약 50% 완료** (Optimism 머지를 통해 자동 구현됨)

**이미 구현된 것**:
- ✅ TraceType 상수 (`TraceTypeAsteriscKona`, `TraceTypeSuperAsteriscKona`)
- ✅ KonaExecutor 완전 구현 (`vm/kona_server_executor.go`, 61 lines)
- ✅ SuperKonaExecutor 구현 (`vm/kona_super_server_executor.go`, 54 lines)
- ✅ 테스트 코드 (`vm/kona_server_executor_test.go`, 46 lines)

**남은 작업**:
- ❌ GameType 상수 (`AsteriscKonaGameType = 3`)
- ❌ Config 설정 및 검증 함수
- ❌ 등록 코드 (`register.go`, `register_task.go`)
- ❌ CLI 플래그
- ❌ kona-client 바이너리 빌드
- ❌ Prestate 준비

**예상 완료 기간**: 9-11일 (Phase 0-9 전체, KonaExecutor가 이미 구현되어 약 1일 단축됨)

---

## 📊 핵심 발견사항 (코드 분석)

### GameType 3의 정체

**Optimism 실제 코드 분석 결과**:

```go
// optimism/op-challenger/game/fault/types/types.go:33
AsteriscKonaGameType GameType = 3
```

**GameType 3 = ASTERISC_KONA**는:
- ✅ **온체인 VM**: GameType 2와 **동일한 RISCV.sol 사용**
- ✅ **StateConverter**: GameType 2의 `asterisc.NewStateConverter` 재사용
- ✅ **TraceAccessor**: GameType 2의 `NewOutputAsteriscTraceAccessor` 재사용
- ❌ **차이점**: 오직 **executor**와 **server binary**만 다름

```go
// optimism/op-challenger/game/fault/register_task.go:161-192
func NewAsteriscKonaRegisterTask(...) *RegisterTask {
    stateConverter := asterisc.NewStateConverter(cfg.Asterisc)  // ← GameType 2와 동일!
    ...
    return outputs.NewOutputAsteriscTraceAccessor(              // ← GameType 2와 동일!
        logger, m, cfg.AsteriscKona, serverExecutor, ...        // ← cfg만 다름
    )
}
```

**즉, GameType 3 = GameType 2 + Rust 기반 kona-client**

---

## 🔍 GameType 2 vs GameType 3 비교

### 실제 코드 비교

| 구성 요소 | GameType 2 (ASTERISC) | GameType 3 (ASTERISC_KONA) |
|----------|---------------------|--------------------------|
| **GameType 상수** | `AsteriscGameType = 2` | `AsteriscKonaGameType = 3` |
| **TraceType** | `TraceTypeAsterisc` | `TraceTypeAsteriscKona` |
| **온체인 VM** | `RISCV.sol` | **동일: RISCV.sol** ✅ |
| **StateConverter** | `asterisc.NewStateConverter(cfg.Asterisc)` | **동일** ✅ |
| **TraceAccessor** | `NewOutputAsteriscTraceAccessor` | **동일** ✅ |
| **Executor** | `vm.NewOpProgramServerExecutor()` | `vm.NewKonaExecutor()` ❌ |
| **Server Binary** | `op-program` (Go) | `kona-client` (Rust) ❌ |
| **Config 필드** | `cfg.Asterisc` | `cfg.AsteriscKona` ❌ |
| **CLI 플래그** | `--asterisc-*` | `--asterisc-kona-*` ❌ |

### 등록 코드 비교

**GameType 2**:
```go
// optimism/op-challenger/game/fault/register.go:91-96
if cfg.TraceTypeEnabled(faultTypes.TraceTypeAsterisc) {
    l2HeaderSource, rollupClient, syncValidator, err := clients.SingleChainClients()
    registerTasks = append(registerTasks, NewAsteriscRegisterTask(
        faultTypes.AsteriscGameType,             // GameType 2
        cfg, m,
        vm.NewOpProgramServerExecutor(logger),   // Go 기반 op-program
        l2HeaderSource, rollupClient, syncValidator
    ))
}
```

**GameType 3**:
```go
// optimism/op-challenger/game/fault/register.go:98-103
if cfg.TraceTypeEnabled(faultTypes.TraceTypeAsteriscKona) {
    l2HeaderSource, rollupClient, syncValidator, err := clients.SingleChainClients()
    registerTasks = append(registerTasks, NewAsteriscKonaRegisterTask(
        faultTypes.AsteriscKonaGameType,         // GameType 3
        cfg, m,
        vm.NewKonaExecutor(),                    // Rust 기반 kona-client
        l2HeaderSource, rollupClient, syncValidator
    ))
}
```

**핵심 차이**: `NewAsteriscKonaRegisterTask` + `vm.NewKonaExecutor()`

---

## 📦 kona-client란?

### 외부 Rust 프로젝트

**공식 저장소**: https://github.com/op-rs/kona (확인됨: optimism 코드에서 참조)

```yaml
# optimism/kurtosis-devnet/interop.yaml:154
- github.com/op-rs/kona/docker/recipes/kona-node/grafana
```

**kona-client 특징**:
- **언어**: Rust (`no_std` 지원)
- **크기**: Go 구현 대비 ~80% 작음 (~8,000 라인)
- **장점**:
  - 경량화 및 최적화
  - ZK proof 시스템 통합 용이 (SP1, RISC Zero)
  - 포터블 (다양한 환경 실행 가능)
- **단점**:
  - 상대적으로 새로운 구현 (2024~2025년 추가)
  - 별도 빌드 필요 (Rust 툴체인)

### kona-client 실행 방식

```go
// optimism/op-challenger/game/fault/trace/vm/kona_server_executor.go:29-40
func (s *KonaExecutor) OracleCommand(cfg Config, dataDir string, inputs utils.LocalGameInputs) ([]string, error) {
    args := []string{
        cfg.Server,                                    // kona-client 바이너리 경로
        "single",                                       // 단일 L2 체인 모드
        "--l1-node-address", cfg.L1,
        "--l1-beacon-address", cfg.L1Beacon,
        "--l2-node-address", cfg.L2s[0],
        "--l1-head", inputs.L1Head.Hex(),
        "--l2-head", inputs.L2Head.Hex(),               // ← op-program과 다름
        "--l2-output-root", inputs.L2OutputRoot.Hex(),  // ← op-program과 다름
        "--l2-claim", inputs.L2Claim.Hex(),
        "--l2-block-number", inputs.L2SequenceNumber.Text(10),
    }

    if s.nativeMode {
        args = append(args, "--native")
    } else {
        args = append(args, "--server")
        args = append(args, "--data-dir", dataDir)
    }

    return args, nil
}
```

**비교: op-program (GameType 2)**:
```bash
# GameType 2
op-program --server --l1 <URL> --l2 <URL> --l2.claim <HASH> --l2.blocknumber <NUM>

# GameType 3
kona-client single --l1-node-address <URL> --l2-node-address <URL> --l2-head <HASH> --l2-output-root <HASH> --l2-claim <HASH> --l2-block-number <NUM>
```

---

## 🏗️ Tokamak-Thanos 현재 상태 분석

### 있는 것 ✅

**GameType 2 (ASTERISC)**: 100% 완료
- ✅ `types.go`: `AsteriscGameType uint32 = 2` 정의
- ✅ `config.go`: `Asterisc vm.Config` 필드
- ✅ `register.go`: `registerAsterisc()` 함수
- ✅ `vm/`: Executor, StateConverter 등
- ✅ `RISCV.sol`: 온체인 검증 컨트랙트 (DevNet 배포 완료)

### 없는 것 ❌

**GameType 3 (ASTERISC_KONA)**: 부분 구현 상태 (약 50%)

#### ✅ 이미 구현된 것 (Optimism v2.0.0-alpha.8에서 머지됨)

**Git 히스토리 확인**:
```bash
7457c5689 - feat(op-challenger): Add TraceTypeAsteriscKona to default --trace-type option
c36de049f - chore(ops): Support kona + asterisc in the op-challenger
bd7c16b87 - feat(op-challenger): Kona interop executor
3cc36be2e - chore(op-challenger): Update kona executor to use subcommand
```

**구현된 컴포넌트**:

1. **TraceType 상수** (`types/types.go:225-229`):
   ```go
   TraceTypeAsteriscKona      TraceType = "asterisc-kona"      // ✅ 구현됨!
   TraceTypeSuperAsteriscKona TraceType = "super-asterisc-kona" // ✅ 구현됨!
   ```

2. **vm/kona_server_executor.go** (`op-challenger/game/fault/trace/vm/`):
   - ✅ **파일이 이미 존재하고 완전히 구현되어 있음!**
   - 위치: `/op-challenger/game/fault/trace/vm/kona_server_executor.go`
   - 61 lines, 완전한 구현
   - Import 경로: `github.com/tokamak-network/tokamak-thanos` (올바름)
   - 구현 내용:
     - `NewKonaExecutor()` 함수 ✅
     - `NewNativeKonaExecutor()` 함수 ✅
     - `OracleCommand()` 메서드 ✅
     - kona-client 실행 로직 완성 ✅

3. **vm/kona_super_server_executor.go**:
   - ✅ GameType 7 (SuperAsteriscKona) 지원 코드
   - 54 lines, 멀티체인 interop 지원

4. **vm/kona_server_executor_test.go**:
   - ✅ 테스트 코드 46 lines

#### ❌ 실제로 없는 것 (구현 필요)

1. **types.go - GameType 상수**:
   ```go
   // ❌ 없음 - 추가 필요!
   AsteriscKonaGameType uint32 = 3
   ```

2. **config.go - Config 필드**:
   ```go
   // ❌ 없음 - 추가 필요!
   AsteriscKona                        vm.Config
   AsteriscKonaAbsolutePreState        string
   AsteriscKonaAbsolutePreStateBaseURL *url.URL
   ```

3. **config.go - TraceTypes 배열** (Line 65):
   ```go
   // 현재
   var TraceTypes = []TraceType{TraceTypeAlphabet, TraceTypeCannon,
                                 TraceTypePermissioned, TraceTypeAsterisc, TraceTypeFast}

   // ❌ TraceTypeAsteriscKona 누락! 추가 필요
   ```

4. **register.go**:
   ```go
   // ❌ 등록 코드 없음 - 추가 필요!
   if cfg.TraceTypeEnabled(config.TraceTypeAsteriscKona) {
       ...
   }
   ```

5. **register_task.go**:
   ```go
   // ❌ 함수 없음 - 추가 필요!
   func NewAsteriscKonaRegisterTask(...) *RegisterTask {
       ...
   }
   ```

6. **flags/flags.go**:
   ```go
   // ❌ CLI 플래그 없음 - 추가 필요!
   AsteriscKonaServerFlag
   AsteriscKonaPreStateFlag
   AsteriscKonaL2CustomFlag
   ```

7. **kona-client 바이너리**:
   - ❌ Rust 프로젝트 빌드 필요
   - 외부 의존성

---

## 📋 통합 작업 목록

### ✅ 이미 완료된 것 (Optimism v2.0.0-alpha.8 머지)

**확인된 Git 커밋**:
```bash
7457c5689 - feat(op-challenger): Add TraceTypeAsteriscKona to default --trace-type option
c36de049f - chore(ops): Support kona + asterisc in the op-challenger
bd7c16b87 - feat(op-challenger): Kona interop executor
3cc36be2e - chore(op-challenger): Update kona executor to use subcommand
```

**완료된 컴포넌트**:
- [x] `types/types.go`: TraceType 상수 정의 (`TraceTypeAsteriscKona`, `TraceTypeSuperAsteriscKona`)
- [x] `vm/kona_server_executor.go`: KonaExecutor 완전 구현 (61 lines)
- [x] `vm/kona_super_server_executor.go`: SuperKonaExecutor 구현 (54 lines)
- [x] `vm/kona_server_executor_test.go`: 테스트 코드 (46 lines)

---

### 🔄 남은 작업

### Phase 0: 온체인 컨트랙트 설정 (0.5일) 🔴 최우선

#### 0.1 현재 상태 확인

**✅ RISCV.sol은 이미 배포되어 있음** (GameType 2 작업 시 완료):

```bash
# DevNet 주소 확인
RISCV_ADDRESS=0xEad59ca6b38c83EcD7735006Db68a29c5e8A96A2
DISPUTE_GAME_FACTORY=0x11c81c1A7979cdd309096D1ea53F887EA9f8D14d

# RISCV.sol 배포 확인
cast call $RISCV_ADDRESS "version()" --rpc-url http://localhost:8545
# 출력: "1.2.0"
```

**⚠️ 중요**: GameType 3는 **새로운 컨트랙트 배포가 필요 없습니다!**
- GameType 2와 **동일한 RISCV.sol 재사용**
- DisputeGameFactory에 **등록만** 하면 됨

#### 0.2 DisputeGameFactory에 GameType 3 등록

**왜 필요한가?**

```solidity
// DisputeGameFactory.sol:28
mapping(GameType => IDisputeGame) public gameImpls;

// DisputeGameFactory.sol:84-97
function create(GameType _gameType, ...) external payable returns (IDisputeGame proxy_) {
    // ⚠️ gameImpls 매핑에서 구현 컨트랙트 찾기
    IDisputeGame impl = gameImpls[_gameType];

    // ⚠️ 등록되지 않았으면 에러!
    if (address(impl) == address(0)) revert NoImplementation(_gameType);

    // 구현 컨트랙트를 clone하여 새 게임 생성
    // ...
}
```

**등록 방법**:

```bash
# Admin 계정 확인
ADMIN_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# GameType 3를 RISCV.sol에 등록
cast send $DISPUTE_GAME_FACTORY \
  "setImplementation(uint32,address)" \
  3 \
  $RISCV_ADDRESS \
  --private-key $ADMIN_KEY \
  --rpc-url http://localhost:8545

# 트랜잭션 확인
# 출력:
# blockHash               0x...
# transactionHash         0x...
# status                  1 (success)
```

#### 0.3 등록 확인

```bash
# GameType 3의 구현 컨트랙트 주소 확인
cast call $DISPUTE_GAME_FACTORY "gameImpls(uint32)" 3 --rpc-url http://localhost:8545
# 출력: 0x000000000000000000000000ead59ca6b38c83ecd7735006db68a29c5e8a96a2
#       (RISCV.sol 주소 - GameType 2와 동일!)

# GameType 2와 비교
cast call $DISPUTE_GAME_FACTORY "gameImpls(uint32)" 2 --rpc-url http://localhost:8545
# 출력: 0x000000000000000000000000ead59ca6b38c83ecd7735006db68a29c5e8a96a2
#       (동일한 주소!)
```

#### 0.4 테스트 게임 생성

```bash
# GameType 3로 테스트 게임 생성
OUTPUT_ROOT=0x0000000000000000000000000000000000000000000000000000000000000123

cast send $DISPUTE_GAME_FACTORY \
  "create(uint32,bytes32,bytes)" \
  3 \
  $OUTPUT_ROOT \
  0x \
  --private-key $ADMIN_KEY \
  --rpc-url http://localhost:8545

# 성공 시 출력:
# transactionHash  0x...
# status           1 (success)

# 생성된 게임 확인
cast call $DISPUTE_GAME_FACTORY "gameCount()" --rpc-url http://localhost:8545
# 출력: 1 (또는 증가된 숫자)

# 최신 게임 정보 확인
GAME_INDEX=0  # 또는 gameCount() - 1
cast call $DISPUTE_GAME_FACTORY "gameAtIndex(uint256)" $GAME_INDEX --rpc-url http://localhost:8545
# 출력: (gameType=3, timestamp=..., proxy=0x...)
```

#### 0.5 배포 스크립트 수정 (선택)

**파일**: `packages/tokamak/contracts-bedrock/scripts/deploy/Deploy.s.sol` (또는 관련 스크립트)

```solidity
// GameType 등록 부분에 추가

// GameType 2 등록 (기존)
disputeGameFactory.setImplementation(GameType.wrap(2), riscv);

// ✅ GameType 3 등록 추가 (RISCV.sol 재사용!)
disputeGameFactory.setImplementation(GameType.wrap(3), riscv);  // 동일한 riscv 주소
```

**또는 Bash 스크립트로 자동화**:

**파일**: `op-challenger/scripts/setup-gametype3.sh` (신규 생성)

```bash
#!/usr/bin/env bash
set -euo pipefail

# GameType 3 온체인 설정 스크립트

echo "=========================================="
echo "GameType 3 온체인 설정"
echo "=========================================="

# 환경 변수 로드
source .env

# RISCV.sol 주소 확인
RISCV_ADDRESS=$(jq -r .RISCV .devnet/addresses.json)
if [ "$RISCV_ADDRESS" == "null" ] || [ -z "$RISCV_ADDRESS" ]; then
    echo "ERROR: RISCV.sol not deployed"
    exit 1
fi

echo "RISCV.sol address: $RISCV_ADDRESS"

# DisputeGameFactory 주소
DISPUTE_GAME_FACTORY=$(jq -r .DisputeGameFactoryProxy .devnet/addresses.json)
echo "DisputeGameFactory: $DISPUTE_GAME_FACTORY"

# Admin private key
ADMIN_KEY=${ADMIN_PRIVATE_KEY:-0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}

# GameType 3 등록
echo ""
echo "Registering GameType 3 to DisputeGameFactory..."

cast send $DISPUTE_GAME_FACTORY \
  "setImplementation(uint32,address)" \
  3 \
  $RISCV_ADDRESS \
  --private-key $ADMIN_KEY \
  --rpc-url http://localhost:8545

# 등록 확인
echo ""
echo "Verifying GameType 3 registration..."
IMPL_ADDRESS=$(cast call $DISPUTE_GAME_FACTORY "gameImpls(uint32)" 3 --rpc-url http://localhost:8545)

if [ "$IMPL_ADDRESS" == "0x000000000000000000000000${RISCV_ADDRESS:2}" ]; then
    echo "✅ GameType 3 registered successfully!"
    echo "   Implementation: $RISCV_ADDRESS"
else
    echo "❌ GameType 3 registration failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "GameType 3 설정 완료!"
echo "=========================================="
```

**실행**:
```bash
chmod +x op-challenger/scripts/setup-gametype3.sh
./op-challenger/scripts/setup-gametype3.sh
```

#### 0.6 검증 체크리스트

- [ ] RISCV.sol 배포 확인 (`cast call $RISCV_ADDRESS "version()"`)
- [ ] GameType 2 등록 확인 (`gameImpls(2)` = RISCV 주소)
- [ ] GameType 3 등록 완료 (`gameImpls(3)` = RISCV 주소)
- [ ] GameType 2와 3이 **동일한 주소** 확인
- [ ] 테스트 게임 생성 성공

---

### Phase 1: 준비 작업 (1일)

#### 1.1 kona-client 빌드 환경 확인

**kona 프로젝트 클론**:
```bash
cd /Users/zena/tokamak-projects/
git clone https://github.com/op-rs/kona.git
cd kona

# Rust 툴체인 확인
rustc --version
cargo --version

# RISC-V 타겟 추가
rustup target add riscv64gc-unknown-linux-gnu
```

**kona-client 빌드**:
```bash
# 빌드 방법 확인 (README 참조)
# 예상: cargo build --release --bin kona-client
```

**참고**: kona는 외부 프로젝트이므로 버전 호환성 주의

---

### Phase 2: 코드 상수 및 타입 추가 (0.5일)

#### 1.1 GameType 상수 추가

**파일**: `op-challenger/game/fault/types/types.go`

```go
const (
    CannonGameType       uint32 = 0
    PermissionedGameType uint32 = 1
    AsteriscGameType     uint32 = 2
+   AsteriscKonaGameType uint32 = 3    // 추가
    FastGameType         uint32 = 254
    AlphabetGameType     uint32 = 255
)
```

#### 1.2 TraceType 추가

**파일**: `op-challenger/game/fault/types/types.go`

```go
const (
    TraceTypeAlphabet          TraceType = "alphabet"
    TraceTypeFast              TraceType = "fast"
    TraceTypeCannon            TraceType = "cannon"
    TraceTypeAsterisc          TraceType = "asterisc"
+   TraceTypeAsteriscKona      TraceType = "asterisc-kona"  // 추가
    TraceTypePermissioned      TraceType = "permissioned"
    TraceTypeSuperCannon       TraceType = "super-cannon"
    TraceTypeSuperPermissioned TraceType = "super-permissioned"
+   TraceTypeSuperAsteriscKona TraceType = "super-asterisc-kona"  // 추가 (선택)
)
```

---

### Phase 3: Config 설정 추가 (1일)

#### 2.1 Config 구조체 필드 추가

**파일**: `op-challenger/config/config.go`

```go
type Config struct {
    // ... 기존 필드

    // Asterisc (GameType 2)
    Asterisc                        vm.Config
    AsteriscAbsolutePreState        string
    AsteriscAbsolutePreStateBaseURL *url.URL

+   // AsteriscKona (GameType 3) - 추가
+   AsteriscKona                        vm.Config
+   AsteriscKonaAbsolutePreState        string
+   AsteriscKonaAbsolutePreStateBaseURL *url.URL

    // ...
}
```

#### 2.2 Error 상수 추가

```go
var (
    // ...
+   ErrMissingAsteriscKonaAbsolutePreState = errors.New("missing asterisc kona absolute pre-state")
+   ErrMissingAsteriscKonaSnapshotFreq     = errors.New("missing asterisc kona snapshot freq")
+   ErrMissingAsteriscKonaInfoFreq         = errors.New("missing asterisc kona info freq")
)
```

#### 2.3 Check() 함수 업데이트

```go
func (c Config) Check() error {
    // ...

+   if c.TraceTypeEnabled(types.TraceTypeAsteriscKona) {
+       if c.RollupRpc == "" {
+           return ErrMissingRollupRpc
+       }
+       if err := c.validateBaseAsteriscKonaOptions(); err != nil {
+           return err
+       }
+   }

    return nil
}

+ func (c Config) validateBaseAsteriscKonaOptions() error {
+     if err := c.AsteriscKona.Check(); err != nil {
+         return fmt.Errorf("asterisc kona: %w", err)
+     }
+     if c.AsteriscKonaAbsolutePreState == "" && c.AsteriscKonaAbsolutePreStateBaseURL == nil {
+         return ErrMissingAsteriscKonaAbsolutePreState
+     }
+     if c.AsteriscKona.SnapshotFreq == 0 {
+         return ErrMissingAsteriscKonaSnapshotFreq
+     }
+     if c.AsteriscKona.InfoFreq == 0 {
+         return ErrMissingAsteriscKonaInfoFreq
+     }
+     return nil
+ }
```

---

### Phase 4: KonaExecutor 검증 및 업데이트 (0.5일) ✅ 이미 존재!

#### 3.1 기존 구현 확인

**✅ kona_server_executor.go 이미 존재함!**

**파일 위치**:
```bash
/Users/zena/tokamak-projects/tokamak-thanos/op-challenger/game/fault/trace/vm/
├── kona_server_executor.go         (61 lines)
├── kona_super_server_executor.go   (54 lines)
└── kona_server_executor_test.go    (46 lines)
```

**확인 명령어**:
```bash
# 파일 존재 확인
ls -lh /Users/zena/tokamak-projects/tokamak-thanos/op-challenger/game/fault/trace/vm/kona*.go

# 파일 내용 확인
cat /Users/zena/tokamak-projects/tokamak-thanos/op-challenger/game/fault/trace/vm/kona_server_executor.go
```

**주요 내용** (이미 완전히 구현됨):
```go
package vm

import (
    "errors"
    "strconv"
    // ...
)

type KonaExecutor struct {
    nativeMode bool
}

func NewKonaExecutor() *KonaExecutor {
    return &KonaExecutor{nativeMode: false}
}

func (s *KonaExecutor) OracleCommand(cfg Config, dataDir string, inputs utils.LocalGameInputs) ([]string, error) {
    // 단일 L2만 지원
    if len(cfg.L2s) != 1 || len(cfg.RollupConfigPaths) > 1 || len(cfg.Networks) > 1 {
        return nil, errors.New("multiple L2s specified but only one supported")
    }

    args := []string{
        cfg.Server,
        "single",                                     // kona 서브커맨드
        "--l1-node-address", cfg.L1,
        "--l1-beacon-address", cfg.L1Beacon,
        "--l2-node-address", cfg.L2s[0],
        "--l1-head", inputs.L1Head.Hex(),
        "--l2-head", inputs.L2Head.Hex(),             // ← op-program과 다름
        "--l2-output-root", inputs.L2OutputRoot.Hex(), // ← op-program과 다름
        "--l2-claim", inputs.L2Claim.Hex(),
        "--l2-block-number", inputs.L2SequenceNumber.Text(10),
    }

    if s.nativeMode {
        args = append(args, "--native")
    } else {
        args = append(args, "--server")
        args = append(args, "--data-dir", dataDir)
    }

    // rollup config 설정
    if len(cfg.RollupConfigPaths) > 0 {
        args = append(args, "--rollup-config-path", cfg.RollupConfigPaths[0])
    } else {
        if len(cfg.Networks) == 0 {
            return nil, errors.New("network is not defined")
        }
        chainCfg := chaincfg.ChainByName(cfg.Networks[0])
        args = append(args, "--l2-chain-id", strconv.FormatUint(chainCfg.ChainID, 10))
    }

    return args, nil
}
```

#### 3.2 Optimism 최신 버전과 비교 (선택 사항)

**Optimism 버전과 diff 비교**:
```bash
# Optimism과 비교하여 차이점 확인
diff /Users/zena/tokamak-projects/optimism/op-challenger/game/fault/trace/vm/kona_server_executor.go \
     /Users/zena/tokamak-projects/tokamak-thanos/op-challenger/game/fault/trace/vm/kona_server_executor.go

# 차이가 있다면 최신 버전으로 업데이트 고려
```

#### 3.3 검증 체크리스트

**확인 사항**:
- [x] `NewKonaExecutor()` 함수 존재
- [x] `NewNativeKonaExecutor()` 함수 존재 (Native 모드)
- [x] `OracleCommand()` 메서드 구현
- [x] kona-client 서브커맨드 `single` 사용
- [x] Import 경로가 `github.com/tokamak-network/tokamak-thanos`로 설정됨
- [x] 테스트 코드 존재 (`kona_server_executor_test.go`)

**결과**: ✅ **모든 항목이 이미 완전히 구현되어 있습니다!**

---

### Phase 5: RegisterTask 추가 (1일)

#### 4.1 NewAsteriscKonaRegisterTask 함수 추가

**파일**: `op-challenger/game/fault/register_task.go`

```go
+ func NewAsteriscKonaRegisterTask(
+     gameType faultTypes.GameType,
+     cfg *config.Config,
+     m caching.Metrics,
+     serverExecutor vm.OracleServerExecutor,
+     l2Client utils.L2HeaderSource,
+     rollupClient outputs.OutputRollupClient,
+     syncValidator SyncValidator,
+ ) *RegisterTask {
+     // ⚠️ 중요: Asterisc의 StateConverter 재사용!
+     stateConverter := asterisc.NewStateConverter(cfg.Asterisc)
+
+     return &RegisterTask{
+         gameType:      gameType,
+         syncValidator: syncValidator,
+         getTopPrestateProvider: func(ctx context.Context, prestateBlock uint64) (faultTypes.PrestateProvider, error) {
+             return outputs.NewPrestateProvider(rollupClient, prestateBlock), nil
+         },
+         getBottomPrestateProvider: cachePrestates(
+             gameType,
+             stateConverter,
+             m,
+             cfg.AsteriscKonaAbsolutePreStateBaseURL,  // ← Kona 전용 설정
+             cfg.AsteriscKonaAbsolutePreState,         // ← Kona 전용 설정
+             filepath.Join(cfg.Datadir, "asterisc-kona-prestates"),
+             func(ctx context.Context, path string) faultTypes.PrestateProvider {
+                 return vm.NewPrestateProvider(path, stateConverter)
+             }),
+         newTraceAccessor: func(
+             logger log.Logger,
+             m metrics.Metricer,
+             prestateProvider faultTypes.PrestateProvider,
+             vmPrestateProvider faultTypes.PrestateProvider,
+             dir string,
+             l1Head eth.BlockID,
+             splitDepth faultTypes.Depth,
+             prestateBlock uint64,
+             poststateBlock uint64) (*trace.Accessor, error) {
+             provider := vmPrestateProvider.(*vm.PrestateProvider)
+             // ⚠️ 중요: Asterisc의 TraceAccessor 재사용!
+             return outputs.NewOutputAsteriscTraceAccessor(
+                 logger, m,
+                 cfg.AsteriscKona,  // ← Kona 전용 config
+                 serverExecutor,    // ← KonaExecutor
+                 l2Client,
+                 prestateProvider,
+                 provider.PrestatePath(),
+                 rollupClient,
+                 dir,
+                 l1Head,
+                 splitDepth,
+                 prestateBlock,
+                 poststateBlock,
+             )
+         },
+     }
+ }
```

**핵심**:
- `asterisc.NewStateConverter(cfg.Asterisc)` 재사용 ✅
- `NewOutputAsteriscTraceAccessor` 재사용 ✅
- `cfg.AsteriscKona` 사용하여 차별화 ✅

---

### Phase 6: register.go 업데이트 (1일)

#### 5.1 등록 코드 추가

**파일**: `op-challenger/game/fault/register.go`

```go
func RegisterGameTypes(...) (CloseFunc, error) {
    // ...

    if cfg.TraceTypeEnabled(config.TraceTypeAsterisc) {
        if err := registerAsterisc(...); err != nil {
            return nil, fmt.Errorf("failed to register asterisc game type: %w", err)
        }
    }

+   // GameType 3: AsteriscKona 등록
+   if cfg.TraceTypeEnabled(config.TraceTypeAsteriscKona) {
+       if err := registerAsteriscKona(faultTypes.AsteriscKonaGameType, registry, oracles, ...); err != nil {
+           return nil, fmt.Errorf("failed to register asterisc kona game type: %w", err)
+       }
+   }

    // ...
}

+ func registerAsteriscKona(
+     gameType uint32,
+     registry Registry,
+     oracles OracleRegistry,
+     // ... 동일한 파라미터들
+ ) error {
+     playerCreator := func(game types.GameMetadata, dir string) (scheduler.GamePlayer, error) {
+         // NewAsteriscKonaRegisterTask 사용
+         task := NewAsteriscKonaRegisterTask(
+             gameType,
+             cfg,
+             m,
+             vm.NewKonaExecutor(),  // ← KonaExecutor 사용!
+             l2Client,
+             rollupClient,
+             syncValidator,
+         )
+         return task.Register(ctx, registry, oracles, ...)
+     }
+     registry.RegisterGameType(gameType, playerCreator)
+     return nil
+ }
```

---

### Phase 7: CLI 플래그 추가 (0.5일)

#### 7.1 flags.go 업데이트

**파일**: `op-challenger/flags/flags.go`

```go
var (
    // ... 기존 플래그

+   AsteriscKonaServerFlag = &cli.StringFlag{
+       Name:    "asterisc-kona-server",
+       Usage:   "Path to kona-client executable to use as pre-image oracle server (asterisc-kona trace type only)",
+       EnvVars: prefixEnvVars("ASTERISC_KONA_SERVER"),
+   }

+   AsteriscKonaPreStateFlag = &cli.StringFlag{
+       Name:    "asterisc-kona-prestate",
+       Usage:   "Path to absolute prestate for asterisc-kona trace type",
+       EnvVars: prefixEnvVars("ASTERISC_KONA_PRESTATE"),
+   }

+   AsteriscKonaL2CustomFlag = &cli.BoolFlag{
+       Name:    "asterisc-kona-l2-custom",
+       Usage:   "Notify kona-host that L2 uses custom config",
+       EnvVars: prefixEnvVars("ASTERISC_KONA_L2_CUSTOM"),
+       Value:   false,
+       Hidden:  true,
+   }
)

// optionalFlags에 추가
var optionalFlags = []cli.Flag{
    // ...
+   AsteriscKonaServerFlag,
+   AsteriscKonaPreStateFlag,
+   AsteriscKonaL2CustomFlag,
}
```

#### 7.2 플래그 검증 함수 추가

```go
+ func CheckAsteriscKonaFlags(ctx *cli.Context) error {
+     if err := checkOutputProviderFlags(ctx); err != nil {
+         return err
+     }
+     if err := CheckAsteriscBaseFlags(ctx, types.TraceTypeAsteriscKona); err != nil {
+         return err
+     }
+     if !ctx.IsSet(AsteriscKonaServerFlag.Name) {
+         return fmt.Errorf("flag %s is required", AsteriscKonaServerFlag.Name)
+     }
+     if !PreStatesURLFlag.IsSet(ctx, types.TraceTypeAsteriscKona) && !ctx.IsSet(AsteriscKonaPreStateFlag.Name) {
+         return fmt.Errorf("flag %s or %s is required",
+             PreStatesURLFlag.EitherFlagName(types.TraceTypeAsteriscKona),
+             AsteriscKonaPreStateFlag.Name)
+     }
+     return nil
+ }

func CheckRequired(ctx *cli.Context, traceTypes []types.TraceType) error {
    // ...
    for _, traceType := range traceTypes {
        switch traceType {
        // ...
+       case types.TraceTypeAsteriscKona:
+           if err := CheckAsteriscKonaFlags(ctx); err != nil {
+               return err
+           }
        // ...
        }
    }
    return nil
}
```

---

### Phase 8: kona prestate 준비 (1-2일)

#### 7.1 kona-client로 prestate 생성

**⚠️ 중요**: GameType 3는 **Asterisc와 동일한 RISCV.sol**을 사용하지만, **prestate는 별도 생성**

**방법: kona 직접 빌드 및 prestate 생성**:

**Step 1: kona 프로젝트 클론**
```bash
cd /Users/zena/tokamak-projects/
git clone https://github.com/ethereum-optimism/kona.git
cd kona

# 최신 stable 태그 확인
git tag | grep -v "rc" | tail -5
git checkout <LATEST_TAG>  # 예: v0.1.0
```

**Step 2: Rust 환경 확인**
```bash
# Rust 버전 확인 (1.70+ 필요)
rustc --version

# RISC-V 타겟 추가
rustup target add riscv64gc-unknown-linux-gnu
rustup target add riscv64imac-unknown-none-elf
```

**Step 3: kona-client 빌드**
```bash
cd /Users/zena/tokamak-projects/kona

# 전체 워크스페이스 빌드
cargo build --release

# 또는 kona-client만 빌드
cargo build --release -p kona-client

# 빌드 결과 확인
ls -lh target/release/kona-client
# 출력: kona-client 바이너리 (~10-50MB)
```

**Step 4: op-program-client (RISC-V) 빌드**
```bash
# kona의 RISC-V 프로그램 빌드
cd /Users/zena/tokamak-projects/kona
make build-rv64

# 또는
cargo build --release --target riscv64gc-unknown-linux-gnu -p op-program-client

# 결과 확인
ls -lh target/riscv64gc-unknown-linux-gnu/release/op-program-client
```

**Step 5: Prestate 생성**
```bash
cd /Users/zena/tokamak-projects/kona

# op-program-client를 Asterisc로 실행하여 prestate 생성
/Users/zena/tokamak-projects/tokamak-thanos/asterisc/bin/asterisc load-elf \
  --path target/riscv64gc-unknown-linux-gnu/release/op-program-client \
  --out /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/prestate-kona.json

# prestate 크기 확인
du -h /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/prestate-kona.json
# 출력: ~200MB
```

**Step 6: Prestate를 Tokamak-Thanos에 복사**
```bash
# kona-client 바이너리 복사
cp /Users/zena/tokamak-projects/kona/target/release/kona-client \
   /Users/zena/tokamak-projects/tokamak-thanos/bin/

# prestate는 이미 생성됨
ls -lh /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/prestate-kona.json
```

**Alternative: Optimism에서 다운로드** (빠른 방법):
```bash
# Optimism 공식 prestate (사용 가능한 경우)
# 예: https://prestates.optimism.io/asterisc-kona/

mkdir -p /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/
wget <PRESTATE_URL> -O op-program/bin/prestate-kona.json
```

#### 7.2 prestate hash 확인

```bash
# prestate hash 계산
cat prestate-kona.json | jq -r '.stateHash'
```

---

### Phase 9: 테스트 및 검증 (2-3일)

#### 8.1 Unit 테스트

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos

# KonaExecutor 테스트
go test ./op-challenger/game/fault/trace/vm -run TestKonaExecutor -v

# Config 테스트
go test ./op-challenger/config -run TestAsteriscKona -v

# RegisterTask 테스트
go test ./op-challenger/game/fault -run TestAsteriscKona -v
```

#### 8.2 통합 테스트

```bash
# 빌드
go build ./op-challenger/cmd/op-challenger

# 실행 테스트 (dry-run)
./bin/op-challenger \
  --trace-type asterisc-kona \
  --asterisc-kona-server /path/to/kona-client \
  --asterisc-kona-prestate ./prestate-kona.json \
  --asterisc-bin ./bin/asterisc \
  --l1-eth-rpc http://localhost:8545 \
  --l2-eth-rpc http://localhost:9545 \
  --rollup-rpc http://localhost:9546 \
  --game-factory-address 0x... \
  --datadir ./temp/challenger-kona \
  # ... 기타 설정
```

#### 8.3 DevNet 테스트

**⚠️ 온체인 컨트랙트는 GameType 2와 동일**:
- RISCV.sol 이미 배포됨 ✅
- DisputeGameFactory에 GameType 3 등록만 필요

```bash
# GameType 3 게임 생성 (L1에서)
cast send <DisputeGameFactory> \
  "create(uint32,bytes32,bytes)" \
  3 \  # GameType 3
  <OUTPUT_ROOT> \
  <EXTRA_DATA>

# Challenger 실행
./bin/op-challenger --trace-type asterisc-kona ...
```

---

## 🔧 핵심 수정사항 요약

### 새로 추가할 파일

**없음!** ✅ 모든 필요한 파일이 이미 Optimism 머지를 통해 존재합니다.

**이미 존재하는 파일**:
- ✅ `vm/kona_server_executor.go` (61 lines, 완전히 구현됨)
- ✅ `vm/kona_super_server_executor.go` (54 lines, GameType 7용)
- ✅ `vm/kona_server_executor_test.go` (46 lines, 테스트 코드)

### 수정할 파일

| 파일 | 추가 내용 | 라인 수 | 상태 |
|-----|----------|--------|------|
| `types/types.go` | GameType 3 상수 추가 | +1 line | ❌ 필요 |
| `config/config.go` | TraceTypes 배열에 AsteriscKona 추가 | +1 line | ❌ 필요 |
| `config/config.go` | AsteriscKona 필드, 검증 함수 | +30 lines | ❌ 필요 |
| `register.go` | registerAsteriscKona() 등록 | +10 lines | ❌ 필요 |
| `register_task.go` | NewAsteriscKonaRegisterTask() | +60 lines | ❌ 필요 |
| `flags/flags.go` | CLI 플래그 3개, 검증 함수 | +40 lines | ❌ 필요 |
| `vm/kona_server_executor.go` | KonaExecutor 구현 | 61 lines | ✅ **이미 존재** |

**총 추가 필요**: ~142 lines (kona_server_executor.go는 이미 존재하므로 제외)

---

## ⚠️ 주의사항 및 제약사항

### 1. 온체인 컨트랙트 등록 (필수!)

**⚠️ 가장 중요**: Off-chain 코드를 다 구현해도, **온체인 등록이 없으면 게임이 생성되지 않습니다!**

#### 동작 원리

```solidity
// DisputeGameFactory.sol:84-97
function create(GameType _gameType, ...) external payable {
    // 1. gameImpls 매핑에서 구현 컨트랙트 조회
    IDisputeGame impl = gameImpls[_gameType];  // ← GameType 3로 조회

    // 2. 등록되지 않았으면 즉시 에러!
    if (address(impl) == address(0)) revert NoImplementation(_gameType);
    //                                       ↑
    //                                    게임 생성 실패!

    // 3. 등록되어 있으면 clone하여 새 게임 생성
    proxy_ = impl.cloneDeterministic(...);
}
```

#### 등록 상태별 결과

**GameType 2 (등록됨)** ✅:
```bash
gameImpls[2] = 0xEad59ca6b38c83EcD7735006Db68a29c5e8A96A2  # RISCV.sol

# 게임 생성 시도
cast send $FACTORY "create(uint32,...)" 2 ...
# 결과: ✅ 성공! 새 게임 생성됨
```

**GameType 3 (등록 안됨)** ❌:
```bash
gameImpls[3] = 0x0000000000000000000000000000000000000000  # 기본값 (미등록)

# 게임 생성 시도
cast send $FACTORY "create(uint32,...)" 3 ...
# 결과: ❌ ERROR: NoImplementation(3)
```

**GameType 3 (등록 완료)** ✅:
```bash
# DisputeGameFactory에 등록
cast send $FACTORY "setImplementation(uint32,address)" 3 $RISCV_ADDRESS

gameImpls[3] = 0xEad59ca6b38c83EcD7735006Db68a29c5e8A96A2  # RISCV.sol

# 게임 생성 시도
cast send $FACTORY "create(uint32,...)" 3 ...
# 결과: ✅ 성공! 새 게임 생성됨
```

#### 등록 방법

```bash
# RISCV.sol 주소 (GameType 2와 동일)
RISCV_ADDRESS=0xEad59ca6b38c83EcD7735006Db68a29c5e8A96A2

# DisputeGameFactory에 GameType 3 등록
cast send $DISPUTE_GAME_FACTORY \
  "setImplementation(uint32,address)" \
  3 \
  $RISCV_ADDRESS \
  --private-key $ADMIN_KEY \
  --rpc-url http://localhost:8545

# 등록 확인
cast call $DISPUTE_GAME_FACTORY "gameImpls(uint32)" 3
# 출력: 0x000...ead59ca6b38c83ecd7735006db68a29c5e8a96a2 ✅
```

**핵심**:
- ✅ **RISCV.sol은 이미 있음** (GameType 2와 공유)
- ✅ **추가 배포 불필요**
- 🔴 **DisputeGameFactory에 GameType 3 등록 필수!**

### 2. kona-client 의존성

- ❌ **외부 Rust 프로젝트** (Optimism 저장소 밖)
- ⚠️ **별도 빌드 필요**
- ⚠️ **버전 호환성 주의**

**권장 방법**:
```bash
# kona를 git submodule로 추가
cd /Users/zena/tokamak-projects/tokamak-thanos
git submodule add https://github.com/op-rs/kona.git kona

# 또는 별도 디렉토리에서 관리
cd /Users/zena/tokamak-projects/
git clone https://github.com/op-rs/kona.git
```

### 3. Prestate 관리

- ⚠️ **GameType 3 전용 prestate 필요**
- ⚠️ **Asterisc prestate와 다를 수 있음**
- ⚠️ **Optimism 공식 prestate 확인 필수**

### 4. StateConverter 재사용

```go
// ✅ 올바른 구현
stateConverter := asterisc.NewStateConverter(cfg.Asterisc)  // Asterisc 것 재사용

// ❌ 잘못된 구현
stateConverter := kona.NewStateConverter(cfg.AsteriscKona)  // 없는 함수!
```

---

## 📅 예상 일정

| Phase | 작업 | 예상 기간 | 상태 | 우선순위 |
|-------|------|----------|------|---------|
| **Phase 0** | **온체인 컨트랙트 설정** | **0.5일** | ❌ 필요 | 🔴 **최우선** |
| **Phase 1** | kona-client 빌드 환경 | 1일 | ❌ 필요 | 🔴 최우선 |
| **Phase 2** | 상수 및 타입 추가 | 0.5일 | ❌ 필요 | 🔴 높음 |
| **Phase 3** | Config 설정 | 1일 | ❌ 필요 | 🔴 높음 |
| **Phase 4** | KonaExecutor 검증 | 0.5일 | ✅ **이미 존재** | 🟢 검증만 |
| **Phase 5** | RegisterTask 추가 | 1일 | ❌ 필요 | 🔴 높음 |
| **Phase 6** | register.go 업데이트 | 1일 | ❌ 필요 | 🔴 높음 |
| **Phase 7** | CLI 플래그 추가 | 0.5일 | ❌ 필요 | 🟡 중간 |
| **Phase 8** | Prestate 준비 | 1-2일 | ❌ 필요 | 🔴 높음 |
| **Phase 9** | 테스트 및 검증 | 2-3일 | ❌ 필요 | 🔴 높음 |

**총 예상 기간**: 9-11일 (Phase 0-9 전체, KonaExecutor 검증 0.5일 포함)

**⚠️ 단축 가능**: KonaExecutor가 이미 구현되어 있어 약 1일 단축됨 (기존 10-12일 예상)

**⚠️ Phase 0 (온체인 설정)을 먼저 완료하면**:
- GameType 3 게임 생성 가능 (온체인)
- Off-chain 코드는 나중에 구현해도 됨
- Phase 1-9만 작업: 8.5-10.5일
- 단계적 개발 가능

---

## ✅ 완료 기준

### 이미 완료된 것 (Optimism 머지로 자동 완료) ✅

- [x] TraceType 상수 정의 (`TraceTypeAsteriscKona`, `TraceTypeSuperAsteriscKona`)
- [x] KonaExecutor 구현 (`vm/kona_server_executor.go`, 61 lines)
- [x] SuperKonaExecutor 구현 (`vm/kona_super_server_executor.go`, 54 lines)
- [x] KonaExecutor 테스트 코드 (`vm/kona_server_executor_test.go`, 46 lines)
- [x] `NewKonaExecutor()` 함수 구현
- [x] `NewNativeKonaExecutor()` 함수 구현 (Native 모드 지원)
- [x] `OracleCommand()` 메서드 구현 (kona-client 실행 로직)

**Git 히스토리**:
```bash
7457c5689 - feat(op-challenger): Add TraceTypeAsteriscKona to default --trace-type option
c36de049f - chore(ops): Support kona + asterisc in the op-challenger
bd7c16b87 - feat(op-challenger): Kona interop executor
3cc36be2e - chore(op-challenger): Update kona executor to use subcommand
```

### 필수 (MVP) - 실제 작업 필요

#### 온체인 (Phase 0) - 가장 먼저!
- [ ] **DisputeGameFactory에 GameType 3 등록** 🔴 최우선
- [ ] 등록 확인 (`gameImpls(3)` = RISCV 주소)
- [ ] 테스트 게임 생성 확인

#### Off-chain (Phase 1-9)
- [ ] kona-client 바이너리 빌드 및 환경 설정
- [ ] GameType 3 상수 정의 (`AsteriscKonaGameType = 3`)
- [ ] config.TraceTypes 배열에 `TraceTypeAsteriscKona` 추가
- [ ] Config.AsteriscKona 필드 추가
- [ ] Config 검증 함수 추가 (`validateBaseAsteriscKonaOptions()`)
- [ ] NewAsteriscKonaRegisterTask 구현
- [ ] register.go에 등록 코드 추가
- [ ] CLI 플래그 추가 (`--asterisc-kona-*`)
- [ ] kona prestate 준비
- [ ] Unit 테스트 통과
- [ ] DevNet에서 Challenger가 GameType 3 게임 해결 성공

### 선택 (Nice to Have)

- [ ] GameType 7 (SuperAsteriscKona) 지원
- [ ] kona-client 자동 빌드 스크립트
- [ ] Multi-prestate URL 지원
- [ ] E2E 테스트 추가

---

## 🏗️ 배포 환경 추가 구성

### 필요한 추가 파일

```bash
tokamak-thanos/
├── bin/
│   ├── kona-client            # ✅ 새로 추가 (Rust 바이너리)
│   ├── asterisc               # 기존 (공통 사용)
│   └── op-challenger          # 기존
│
├── op-program/bin/
│   ├── prestate.json          # 기존 (GameType 2)
│   └── prestate-kona.json     # ✅ 새로 추가 (GameType 3)
```

### deploy-modular.sh 수정

**파일**: `op-challenger/scripts/deploy-modular.sh`

```bash
# 기존
--dg-type 0|1|2|254|255  # Dispute game type

# 수정 후
--dg-type 0|1|2|3|254|255  # ✅ GameType 3 추가
```

**파일**: `op-challenger/scripts/modules/vm-build.sh`

```bash
build_vm_binaries() {
    local dg_type=$1

    case ${dg_type} in
        0|1)  # Cannon
            build_cannon
            ;;
        2)    # Asterisc
            build_asterisc
            ;;
        3)    # ✅ AsteriscKona 추가
            build_asterisc           # VM은 동일
            build_kona_client        # kona-client 추가 빌드
            ;;
        254|255)  # Test
            ;;
    esac
}

# ✅ 새로운 함수 추가
build_kona_client() {
    log_info "Building kona-client..."

    # kona 프로젝트가 있는지 확인
    if [ ! -d "/Users/zena/tokamak-projects/kona" ]; then
        log_error "kona project not found. Please clone first:"
        log_error "  cd /Users/zena/tokamak-projects"
        log_error "  git clone https://github.com/ethereum-optimism/kona.git"
        exit 1
    fi

    # kona-client 빌드
    cd /Users/zena/tokamak-projects/kona
    cargo build --release -p kona-client

    # 바이너리 복사
    cp target/release/kona-client "${PROJECT_ROOT}/bin/"

    log_success "kona-client built successfully"
}

# ✅ prestate 생성 추가
generate_prestates() {
    local dg_type=$1

    case ${dg_type} in
        0|1)
            generate_cannon_prestate
            ;;
        2)
            generate_asterisc_prestate
            ;;
        3)
            generate_kona_prestate  # ✅ 새로운 함수
            ;;
    esac
}

generate_kona_prestate() {
    log_info "Generating kona prestate..."

    # op-program-client (RISC-V) 빌드
    cd /Users/zena/tokamak-projects/kona
    make build-rv64 || cargo build --release --target riscv64gc-unknown-linux-gnu -p op-program-client

    # asterisc로 prestate 생성
    "${PROJECT_ROOT}/asterisc/bin/asterisc" load-elf \
        --path target/riscv64gc-unknown-linux-gnu/release/op-program-client \
        --out "${PROJECT_ROOT}/op-program/bin/prestate-kona.json"

    log_success "kona prestate generated"
}
```

### env-setup.sh 수정

**파일**: `op-challenger/scripts/modules/env-setup.sh`

```bash
setup_challenger_env() {
    local dg_type=$1

    case ${dg_type} in
        0|1)
            CHALLENGER_TRACE_TYPE="cannon"
            ASTERISC_BIN=""
            KONA_BIN=""
            ;;
        2)
            CHALLENGER_TRACE_TYPE="asterisc"
            ASTERISC_BIN="${PROJECT_ROOT}/asterisc/bin/asterisc"
            ASTERISC_SERVER="${PROJECT_ROOT}/op-program/bin/op-program"
            ASTERISC_PRESTATE="${PROJECT_ROOT}/op-program/bin/prestate-rv64.json"
            KONA_BIN=""
            ;;
        3)  # ✅ GameType 3 추가
            CHALLENGER_TRACE_TYPE="asterisc-kona"
            ASTERISC_BIN="${PROJECT_ROOT}/asterisc/bin/asterisc"
            KONA_SERVER="${PROJECT_ROOT}/bin/kona-client"
            KONA_PRESTATE="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
            ;;
        254|255)
            CHALLENGER_TRACE_TYPE="alphabet"
            ;;
    esac

    # .env 파일 생성
    cat > "${PROJECT_ROOT}/.env.challenger" <<EOF
DG_TYPE=${dg_type}
CHALLENGER_TRACE_TYPE=${CHALLENGER_TRACE_TYPE}

# GameType 2 (Asterisc)
ASTERISC_BIN=${ASTERISC_BIN}
ASTERISC_SERVER=${ASTERISC_SERVER}
ASTERISC_PRESTATE=${ASTERISC_PRESTATE}

# GameType 3 (AsteriscKona) - ✅ 추가
KONA_SERVER=${KONA_SERVER}
KONA_PRESTATE=${KONA_PRESTATE}
EOF
}
```

### docker-compose-full.yml 수정 ⭐

**파일**: `op-challenger/scripts/docker-compose-full.yml`

#### 변경사항 요약

GameType 3 지원을 위해 `op-challenger` 서비스에 다음을 추가했습니다:

**1. Volume 마운트 추가** (Line 272):
```yaml
volumes:
  - '${KONA_BIN:-${PROJECT_ROOT}/bin}:/kona'  # ✅ kona-client 바이너리 추가
```

**2. 환경변수 추가** (Lines 295-300):
```yaml
environment:
  # ✅ AsteriscKona (GameType 3) 설정 - 새로 추가
  OP_CHALLENGER_ASTERISC_KONA_ROLLUP_CONFIG: /devnet/rollup.json
  OP_CHALLENGER_ASTERISC_KONA_L2_GENESIS: /devnet/genesis-l2.json
  OP_CHALLENGER_ASTERISC_KONA_BIN: /asterisc/asterisc        # GameType 2와 동일
  OP_CHALLENGER_ASTERISC_KONA_SERVER: /kona/kona-client      # kona-client
  OP_CHALLENGER_ASTERISC_KONA_PRESTATE: ${KONA_PRESTATE:-/op-program/prestate-kona.json}
```

**3. Platform 주석 업데이트** (Line 259):
```yaml
platform: linux/amd64  # kona-client 포함
```

#### 전체 op-challenger 서비스 설정

```yaml
op-challenger:
  platform: linux/amd64  # Required for x86_64 binaries (asterisc, cannon, op-program, kona-client)
  build:
    context: ${PROJECT_ROOT}
    dockerfile: ops/docker/op-stack-go/Dockerfile
    target: op-challenger-target
  image: tokamaknetwork/thanos-op-challenger:${IMAGE_TAG:-latest}
  volumes:
    - 'challenger_data:/db'
    - '${PROJECT_ROOT}/op-program/bin:/op-program'
    - '${PROJECT_ROOT}/cannon/bin:/cannon'
    - '${ASTERISC_BIN:-${PROJECT_ROOT}/asterisc/bin}:/asterisc'
    - '${KONA_BIN:-${PROJECT_ROOT}/bin}:/kona'  # ✅ 추가
    - '${PROJECT_ROOT}/.devnet:/devnet'
  environment:
    # 공통 설정
    OP_CHALLENGER_L1_ETH_RPC: http://l1:8545
    OP_CHALLENGER_ROLLUP_RPC: http://challenger-op-node:8545
    OP_CHALLENGER_L2_ETH_RPC: http://challenger-l2:8545
    OP_CHALLENGER_TRACE_TYPE: ${CHALLENGER_TRACE_TYPE:-cannon}

    # Cannon (GameType 0, 1)
    OP_CHALLENGER_CANNON_BIN: /cannon/cannon
    OP_CHALLENGER_CANNON_SERVER: /op-program/op-program
    OP_CHALLENGER_CANNON_PRESTATE: /op-program/prestate.json

    # Asterisc (GameType 2)
    OP_CHALLENGER_ASTERISC_BIN: /asterisc/asterisc
    OP_CHALLENGER_ASTERISC_SERVER: /op-program/op-program
    OP_CHALLENGER_ASTERISC_PRESTATE: ${ASTERISC_PRESTATE:-/asterisc/prestate.json}

    # ✅ AsteriscKona (GameType 3) - 새로 추가
    OP_CHALLENGER_ASTERISC_KONA_ROLLUP_CONFIG: /devnet/rollup.json
    OP_CHALLENGER_ASTERISC_KONA_L2_GENESIS: /devnet/genesis-l2.json
    OP_CHALLENGER_ASTERISC_KONA_BIN: /asterisc/asterisc
    OP_CHALLENGER_ASTERISC_KONA_SERVER: /kona/kona-client
    OP_CHALLENGER_ASTERISC_KONA_PRESTATE: ${KONA_PRESTATE:-/op-program/prestate-kona.json}
```

#### 실행 방법

**Step 1: 환경변수 파일 생성** (`.env` 또는 `.env.challenger`):

```bash
# 프로젝트 루트 경로
PROJECT_ROOT=/Users/zena/tokamak-projects/tokamak-thanos

# L1/L2 설정
L1_BEACON_URL=http://localhost:5052

# 계정 private keys
CHALLENGER_PRIVATE_KEY=0x...
PROPOSER_PRIVATE_KEY=0x...
BATCHER_PRIVATE_KEY=0x...

# GameType 및 컨트랙트
DG_TYPE=3  # ✅ GameType 3
GAME_FACTORY_ADDRESS=0x11c81c1A7979cdd309096D1ea53F887EA9f8D14d

# Challenger 설정
CHALLENGER_TRACE_TYPE=asterisc-kona  # ✅ GameType 3

# ✅ GameType 3 전용 경로
KONA_BIN=${PROJECT_ROOT}/bin
KONA_PRESTATE=/op-program/prestate-kona.json

# GameType 2 경로 (공통으로 사용)
ASTERISC_BIN=${PROJECT_ROOT}/asterisc/bin
ASTERISC_PRESTATE=/asterisc/prestate.json
```

**Step 2: kona-client 빌드 확인**:

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos

# kona-client가 있는지 확인
ls -lh bin/kona-client
# 출력: -rwxr-xr-x ... bin/kona-client

# prestate-kona.json 확인
ls -lh op-program/bin/prestate-kona.json
# 출력: -rw-r--r-- ... op-program/bin/prestate-kona.json
```

**Step 3: Docker Compose 실행**:

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos

# GameType 3로 전체 스택 실행
docker compose -f op-challenger/scripts/docker-compose-full.yml up -d

# 로그 확인
docker compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-challenger
```

#### 로그 확인

GameType 3가 정상 동작하는지 확인:

```bash
# Challenger 로그에서 trace-type 확인
docker compose -f op-challenger/scripts/docker-compose-full.yml logs op-challenger | grep -i "trace"

# 출력 예시:
# INFO [op-challenger] Trace type: asterisc-kona
# INFO [op-challenger] Loaded kona-client: /kona/kona-client
# INFO [op-challenger] Prestate: /op-program/prestate-kona.json

# kona-client 실행 확인
docker compose -f op-challenger/scripts/docker-compose-full.yml exec op-challenger ls -lh /kona/kona-client
```

#### GameType별 실행 예시

**GameType 0 (Cannon)**:
```bash
CHALLENGER_TRACE_TYPE=cannon \
DG_TYPE=0 \
docker compose -f op-challenger/scripts/docker-compose-full.yml up -d
```

**GameType 2 (Asterisc)**:
```bash
CHALLENGER_TRACE_TYPE=asterisc \
DG_TYPE=2 \
ASTERISC_PRESTATE=/asterisc/prestate.json \
docker compose -f op-challenger/scripts/docker-compose-full.yml up -d
```

**GameType 3 (AsteriscKona)** ✅:
```bash
CHALLENGER_TRACE_TYPE=asterisc-kona \
DG_TYPE=3 \
KONA_BIN=/Users/zena/tokamak-projects/tokamak-thanos/bin \
KONA_PRESTATE=/op-program/prestate-kona.json \
docker compose -f op-challenger/scripts/docker-compose-full.yml up -d
```

#### 트러블슈팅

**문제 1: kona-client not found**
```bash
# 증상
ERROR: /kona/kona-client: no such file or directory

# 해결
# 1. kona-client 빌드 확인
ls -lh /Users/zena/tokamak-projects/tokamak-thanos/bin/kona-client

# 2. KONA_BIN 환경변수 확인
echo $KONA_BIN

# 3. Docker volume 재마운트
docker compose -f op-challenger/scripts/docker-compose-full.yml down -v
docker compose -f op-challenger/scripts/docker-compose-full.yml up -d
```

**문제 2: prestate-kona.json not found**
```bash
# 증상
ERROR: Failed to load prestate: /op-program/prestate-kona.json

# 해결
# 1. prestate 생성 확인
ls -lh /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/prestate-kona.json

# 2. prestate 생성 (없는 경우)
cd /Users/zena/tokamak-projects/kona
./asterisc/bin/asterisc load-elf \
  --path target/riscv64gc-unknown-linux-gnu/release/op-program-client \
  --out /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/prestate-kona.json
```

**문제 3: 잘못된 trace-type**
```bash
# 증상
ERROR: unsupported trace type: asterisc-kona

# 해결
# GameType 3 코드가 구현되지 않은 경우입니다.
# Phase 1-6의 코드 수정을 먼저 완료해야 합니다.
```

### Makefile 추가

**파일**: `Makefile` (프로젝트 루트)

```makefile
# 기존
build-asterisc:
	cd asterisc && make asterisc
	cd op-program && make op-program-rv64

# ✅ 새로 추가
build-kona:
	@echo "Building kona-client..."
	@if [ ! -d "../kona" ]; then \
		echo "ERROR: kona project not found"; \
		echo "Please clone: git clone https://github.com/ethereum-optimism/kona.git"; \
		exit 1; \
	fi
	cd ../kona && cargo build --release -p kona-client
	cp ../kona/target/release/kona-client ./bin/

build-kona-prestate: build-kona
	@echo "Generating kona prestate..."
	cd ../kona && make build-rv64
	./asterisc/bin/asterisc load-elf \
		--path ../kona/target/riscv64gc-unknown-linux-gnu/release/op-program-client \
		--out ./op-program/bin/prestate-kona.json

# 전체 빌드 타겟 수정
build-all: build-go build-contracts build-asterisc build-kona  # ✅ build-kona 추가
```

### 배포 스크립트 실행 예시

```bash
# GameType 2 배포 (기존)
./op-challenger/scripts/deploy-modular.sh --dg-type 2

# GameType 3 배포 (새로 추가)
./op-challenger/scripts/deploy-modular.sh --dg-type 3
  ↓
1. kona-client 빌드 확인
2. kona prestate 생성
3. .env.challenger 생성 (KONA_* 변수 포함)
4. docker-compose up (kona-client 마운트)
```

### CI/CD 파이프라인 수정

**파일**: `.github/workflows/build.yml` (있다면)

```yaml
jobs:
  build:
    steps:
      # 기존 빌드
      - name: Build Cannon
        run: make build-cannon

      - name: Build Asterisc
        run: make build-asterisc

      # ✅ 새로 추가
      - name: Setup Rust
        uses: actions-rs/toolchain@v1
        with:
          toolchain: stable
          target: riscv64gc-unknown-linux-gnu

      - name: Clone Kona
        run: |
          cd ..
          git clone https://github.com/ethereum-optimism/kona.git

      - name: Build Kona
        run: make build-kona

      - name: Build Kona Prestate
        run: make build-kona-prestate
```

---

## 🔍 검증 체크리스트

### 빌드 검증

```bash
# 빌드 성공
go build ./op-challenger/...

# 타입 체크
go vet ./op-challenger/...

# 테스트 실행
go test ./op-challenger/game/fault/... -v
```

### 실행 검증

```bash
# --help 확인
./bin/op-challenger --help | grep -A 5 "asterisc-kona"

# Config 파싱 테스트
./bin/op-challenger \
  --trace-type asterisc-kona \
  --asterisc-kona-server /path/to/kona-client \
  --asterisc-kona-prestate ./prestate-kona.json \
  --dry-run

# DevNet 실행
make devnet-up
./bin/op-challenger --trace-type asterisc-kona ...
```

### 온체인 검증

```bash
# DisputeGameFactory에서 GameType 3 구현 확인
cast call <DisputeGameFactory> "gameImpls(uint32)" 3
# 출력: <RISCV.sol 주소>

# GameType 3 게임 생성
cast send <DisputeGameFactory> "create(uint32,bytes32,bytes)" 3 <ROOT> 0x
# 출력: tx hash

# 게임 상태 확인
cast call <GameProxy> "gameType()"
# 출력: 3
```

---

## 📚 참고 자료

### 코드 분석 기준

- **Optimism 저장소**: `/Users/zena/tokamak-projects/optimism`
- **분석 파일**:
  - `op-challenger/game/fault/types/types.go`
  - `op-challenger/config/config.go`
  - `op-challenger/game/fault/register.go`
  - `op-challenger/game/fault/register_task.go`
  - `op-challenger/game/fault/trace/vm/kona_server_executor.go`
  - `op-challenger/flags/flags.go`

### 외부 리소스

- **kona GitHub**: https://github.com/op-rs/kona
- **kona 문서**: https://devdocs.optimism.io/book/ (예상)
- **Optimism 블로그**: Kona 관련 포스트
- **RISC-V Spec**: https://riscv.org/

### 관련 문서

- [GameType 2 통합 계획](./gametype2-integration-plan-ko.md)
- [RISC-V GameTypes 비교](./risc-v-gametypes-comparison-ko.md)
- [Asterisc vs Tokamak 비교](./asterisc-comparison-optimism-vs-tokamak-ko.md)
- [Game Types와 VMs](./game-types-and-vms-ko.md)

---

## 🚀 Quick Start

### 최소 작업으로 GameType 3 지원하기

**✅ 좋은 소식: KonaExecutor가 이미 구현되어 있어서 작업량이 줄어들었습니다!**

```bash
# 1. 기존 KonaExecutor 확인
cd /Users/zena/tokamak-projects/tokamak-thanos
ls -lh op-challenger/game/fault/trace/vm/kona*.go
# 출력:
# ✅ kona_server_executor.go (61 lines)
# ✅ kona_super_server_executor.go (54 lines)
# ✅ kona_server_executor_test.go (46 lines)

# 2. kona-client 빌드
cd /Users/zena/tokamak-projects/
git clone https://github.com/op-rs/kona.git
cd kona
cargo build --release --bin kona-client

# 3. 코드 수정 (핵심 5개 파일만 - KonaExecutor는 이미 있음!)
cd /Users/zena/tokamak-projects/tokamak-thanos

# types.go: GameType 3 상수 추가 (1 line)
# config.go: TraceTypes 배열에 추가 (1 line)
# config.go: AsteriscKona 필드 추가 (~30 lines)
# register.go: 등록 코드 추가 (~10 lines)
# register_task.go: NewAsteriscKonaRegisterTask 추가 (~60 lines)
# flags.go: CLI 플래그 추가 (~40 lines)

# 4. 빌드 및 테스트
go build ./op-challenger/cmd/op-challenger
go test ./op-challenger/game/fault/... -v

# 5. DevNet 테스트
make devnet-up
./bin/op-challenger --trace-type asterisc-kona \
  --asterisc-kona-server /path/to/kona-client \
  --asterisc-kona-prestate ./prestate-kona.json \
  # ...
```

**예상 소요 시간**: 9-11일 (Phase 0-9 전체, KonaExecutor는 이미 구현되어 약 1일 단축)

---

**⚠️ 핵심 요약**:

### Tokamak-Thanos 현재 상태 (2025-10-26 기준)

**✅ 이미 구현된 것 (약 50%)**:
1. **TraceType 상수**: `TraceTypeAsteriscKona`, `TraceTypeSuperAsteriscKona` ✅
2. **KonaExecutor**: `vm/kona_server_executor.go` (61 lines) ✅
3. **SuperKonaExecutor**: `vm/kona_super_server_executor.go` (54 lines) ✅
4. **테스트 코드**: `vm/kona_server_executor_test.go` (46 lines) ✅

**❌ 남은 작업 (약 50%)**:
1. **GameType 상수**: `AsteriscKonaGameType = 3` 추가 필요
2. **Config 설정**: `AsteriscKona` 필드 및 검증 함수 추가
3. **등록 코드**: `register.go`, `register_task.go` 추가
4. **CLI 플래그**: `flags.go` 추가
5. **kona-client 빌드**: Rust 프로젝트 빌드
6. **Prestate 준비**: kona 전용 prestate 생성

### 핵심 특징

1. **GameType 3 = GameType 2 + Rust 기반 kona-client**
2. **온체인 컨트랙트는 동일** (RISCV.sol 재사용)
3. **StateConverter와 TraceAccessor 재사용** (asterisc 것 사용)
4. **차이점**: Executor (이미 구현됨 ✅), Server Binary, Config 설정만
5. **외부 의존성**: kona-client (별도 Rust 프로젝트 빌드 필요)

**예상 작업 기간**: 9-11일 (Phase 0-9 전체, KonaExecutor가 이미 구현되어 약 1일 단축)

---

## 🎯 GameType 3 통합 요약

### 온체인 vs Off-chain 구분

| 영역 | GameType 2 | GameType 3 | 작업 필요 여부 |
|-----|-----------|-----------|-------------|
| **온체인** | RISCV.sol | **동일한 RISCV.sol** | ⚠️ **등록만 필요** |
| **Off-chain** | op-program (Go) | kona-client (Rust) | ❌ **새로운 빌드 필요** |
| **StateConverter** | asterisc | **동일** | ✅ 재사용 |
| **TraceAccessor** | OutputAsteriscTraceAccessor | **동일** | ✅ 재사용 |

### 작업 우선순위

**1순위 (Phase 0)**: 온체인 설정 (0.5일)
```bash
# GameType 3를 DisputeGameFactory에 등록
cast send $FACTORY "setImplementation(uint32,address)" 3 $RISCV_ADDRESS
```
→ 이것만 하면 **GameType 3 게임 생성 가능**

**2순위 (Phase 1)**: kona-client 빌드 (1일)
```bash
cd /Users/zena/tokamak-projects/kona
cargo build --release -p kona-client
```
→ Challenger가 사용할 바이너리

**3순위 (Phase 2-7)**: Off-chain 코드 (3-4일)
- 상수, Config, RegisterTask, 플래그 추가
→ Challenger가 GameType 3 게임 참여 가능

**4순위 (Phase 8-9)**: Prestate 및 테스트 (3-5일)
- kona prestate 생성 및 검증
→ 프로덕션 준비 완료

### 단계적 배포 전략

**Step 1: 온체인만 (Phase 0)** - 0.5일
```
✅ GameType 3 등록
→ Proposer가 GameType 3 게임 생성 가능
→ 다른 사람의 Challenger가 게임 해결 가능
```

**Step 2: Challenger 코드 (Phase 1-7)** - 4-5일
```
✅ Off-chain 구현
→ 자체 Challenger가 GameType 3 게임 참여 가능
```

**Step 3: 프로덕션 (Phase 8-9)** - 3-5일
```
✅ Prestate 생성 및 전체 테스트
→ 프로덕션 배포 가능
```

### GameType별 필요 컴포넌트

```
GameType 0 (Cannon):
├── 온체인: MIPS.sol
├── VM: cannon
├── Server: op-program
└── Prestate: prestate.json

GameType 2 (Asterisc):
├── 온체인: RISCV.sol ✅
├── VM: asterisc ✅
├── Server: op-program (Go)
└── Prestate: prestate-rv64.json ✅

GameType 3 (AsteriscKona):
├── 온체인: RISCV.sol ← ✅ GameType 2와 동일 (재사용!)
├── VM: asterisc ← ✅ GameType 2와 동일 (재사용!)
├── Server: kona-client (Rust) ← ❌ 새로 빌드 필요!
└── Prestate: prestate-kona.json ← ❌ 새로 생성 필요!
```

### 핵심 메시지

**GameType 3는 "GameType 2의 변형"입니다**:
- 온체인 검증: 동일 (RISCV.sol)
- 증명 생성 방식: 다름 (kona-client)
- 목적: Rust 기반 경량화, ZK 통합 준비

**통합 난이도**: ⭐⭐⭐☆☆ (중간)
- 온체인: 매우 쉬움 (등록만)
- Off-chain: 중간 (코드 추가 ~142 lines)
- 외부 의존성: kona-client 빌드

---

## 📝 다음 단계

이 문서가 완성되면:

1. **Phase 0부터 시작** (온체인 설정)
   ```bash
   ./op-challenger/scripts/setup-gametype3.sh
   ```

2. **Phase 1-9 순차 진행** (Off-chain 구현)
   - 각 Phase마다 테스트
   - 문제 발생 시 롤백

3. **DevNet 전체 테스트**
   ```bash
   docker compose -f op-challenger/scripts/docker-compose-full.yml up -d
   ```

4. **프로덕션 배포** (선택)

---

**마지막 업데이트**: 2025-10-26
**문서 버전**: v1.0
**검토 상태**: ✅ 코드 분석 완료, 실행 가능한 계획

이 계획을 따라 단계별로 진행하면 GameType 3를 안전하게 통합할 수 있습니다.

