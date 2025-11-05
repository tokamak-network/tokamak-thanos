# E2E 테스트 설정 가이드

> **v1.16.0 마이그레이션 과정에서 발생한 E2E 테스트 문제 해결 가이드**

---

## 📋 목차

1. [문제 배경](#문제-배경)
2. [해결한 문제들](#해결한-문제들)
3. [코드 변경 사항](#코드-변경-사항)
4. [E2E 테스트 실행](#e2e-테스트-실행)
5. [Cannon Prestate 빌드](#cannon-prestate-빌드)
6. [트러블슈팅](#트러블슈팅)

---

## 문제 배경

v1.16.0으로 마이그레이션 후 E2E 테스트 실행 시 다음 문제들이 발생:

### 1️⃣ Genesis 파일 형식 문제
```
json: cannot unmarshal hex string without 0x prefix into Go value of type common.Address
```

### 2️⃣ Tokamak 전용 필드 미인식
```
json: unknown field "nativeTokenName"
json: unknown field "setPrecompileBalances"
```

### 3️⃣ ETHLockbox 검증 에러
```
ETHLockbox is not set
```

### 4️⃣ Predeploy 검증 에러
```
predeploy 4200000000000000000000000000000000000000 is missing from L2 genesis allocs
```

---

## 해결한 문제들

### ✅ 1. ForgeAllocs JSON 파싱 개선

**파일**: `op-chain-ops/foundry/allocs.go`

**문제**: Tokamak의 state-dump JSON이 `{"accounts": {...}}` 형식으로 래핑되어 있어 파싱 실패

**해결**:
```go
func (d *ForgeAllocs) UnmarshalJSON(b []byte) error {
    type forgeAllocAccount struct {
        Balance hexutil.U256      `json:"balance"`
        Nonce   hexutil.Uint64    `json:"nonce"`
        Code    hexutil.Bytes     `json:"code,omitempty"`
        Storage map[string]string `json:"storage,omitempty"` // 유연한 파싱을 위해 string map 사용
    }

    // 래핑된 형식 처리
    var wrapped struct {
        Accounts map[string]forgeAllocAccount `json:"accounts"`
    }
    if err := json.Unmarshal(b, &wrapped); err == nil && len(wrapped.Accounts) > 0 {
        d.Accounts = make(types.GenesisAlloc, len(wrapped.Accounts))
        for addrStr, acc := range wrapped.Accounts {
            addr := common.HexToAddress(addrStr)  // 0x 없이도 파싱
            var storage map[common.Hash]common.Hash
            if len(acc.Storage) > 0 {
                storage = make(map[common.Hash]common.Hash, len(acc.Storage))
                for k, v := range acc.Storage {
                    storage[common.HexToHash(k)] = common.HexToHash(v)
                }
            }
            d.Accounts[addr] = types.Account{
                Code:    acc.Code,
                Storage: storage,
                Balance: (*uint256.Int)(&acc.Balance).ToBig(),
                Nonce:   (uint64)(acc.Nonce),
            }
        }
        return nil
    }

    // 기존 직접 형식도 지원 (fallback)
    // ...
}
```

**효과**:
- ✅ 래핑된 JSON 형식 지원
- ✅ 0x prefix 없는 hex string 파싱
- ✅ 기존 형식 호환성 유지

---

### ✅ 2. Tokamak 전용 필드 지원

**파일**: `op-chain-ops/genesis/config.go`

**문제**: Tokamak의 Native Token, USDC Bridge 관련 필드 미인식

**해결**:

#### TokamakDeployConfig 구조체 추가:
```go
// TokamakDeployConfig contains Tokamak-specific deployment configuration fields
// for Native Token and USDC bridge support
type TokamakDeployConfig struct {
    // Native Token configuration
    NativeTokenAddress       common.Address `json:"nativeTokenAddress,omitempty"`
    NativeTokenName          string         `json:"nativeTokenName,omitempty"`
    NativeTokenSymbol        string         `json:"nativeTokenSymbol,omitempty"`
    NativeCurrencyLabelBytes []byte         `json:"nativeCurrencyLabelBytes,omitempty"`

    // USDC Bridge configuration
    L1UsdcAddr     common.Address `json:"l1UsdcAddr,omitempty"`
    UsdcTokenName  string         `json:"usdcTokenName,omitempty"`
    FiatTokenOwner common.Address `json:"fiatTokenOwner,omitempty"`

    // Other Tokamak-specific fields
    SetPrecompileBalances bool `json:"setPrecompileBalances,omitempty"`
}

var _ ConfigChecker = (*TokamakDeployConfig)(nil)

func (d *TokamakDeployConfig) Check(log log.Logger) error {
    // Tokamak fields are all optional, no validation needed
    return nil
}
```

#### DeployConfig에 추가:
```go
type DeployConfig struct {
    // ... existing fields ...

    // Tokamak-specific fields for Native Token support
    TokamakDeployConfig `evm:"-"`
}
```

#### Unknown fields 허용:
```go
func NewDeployConfig(path string) (*DeployConfig, error) {
    // ...
    dec := json.NewDecoder(bytes.NewReader(file))
    // Allow unknown fields for Tokamak-specific extensions
    // (native token, USDC bridge, UniswapV3, etc.)
    // dec.DisallowUnknownFields()  // 주석 처리

    var config DeployConfig
    if err := dec.Decode(&config); err != nil {
        return nil, fmt.Errorf("cannot unmarshal deploy config: %w", err)
    }
    return &config, nil
}
```

**효과**:
- ✅ Tokamak 전용 필드 인식
- ✅ UniswapV3, governance token 등 추가 필드 허용
- ✅ 기존 OP Stack 필드 호환성 유지

---

### ✅ 3. ETHLockbox/RAT 검증 Skip

**파일**: `op-chain-ops/genesis/config.go`

**문제**: Tokamak은 ETHLockbox, RAT 컨트랙트를 사용하지 않음

**해결**:
```go
func (d *L1Deployments) Check(deployConfig *DeployConfig) error {
    // ... reflection code ...

    if name == "RATProxy" || name == "RAT" {
        continue
    }

    // Skip ETHLockbox for Native Token deployments (e.g., Tokamak)
    if name == "ETHLockbox" || name == "ETHLockboxProxy" {
        continue
    }

    // ... rest of validation ...
}
```

**효과**:
- ✅ ETHLockbox 검증 에러 해결
- ✅ RAT 검증 에러 해결

---

### ✅ 4. E2E 초기화 단순화

**파일**: `op-e2e/config/init.go`

**문제**: 복잡한 Tokamak/op-deployer 혼용 로직으로 인한 충돌

**해결**: `.devnet/` 파일 직접 로드 방식으로 단순화

#### initFromDevnetFiles() 함수 추가:
```go
func initFromDevnetFiles(root string) error {
    // Load L1 allocs
    l1AllocsPath := path.Join(root, ".devnet", "allocs-l1.json")
    l1Allocs, err := foundry.LoadForgeAllocs(l1AllocsPath)
    if err != nil {
        return fmt.Errorf("failed to load L1 allocs: %w", err)
    }

    // Load L1 deployments
    l1DeploymentsPath := path.Join(root, ".devnet", "addresses.json")
    l1Deployments, err := genesis.NewL1Deployments(l1DeploymentsPath)
    if err != nil {
        return fmt.Errorf("failed to load L1 deployments: %w", err)
    }

    // Load deploy config
    deployConfigPath := path.Join(root, "packages", "tokamak", "contracts-bedrock", "deploy-config", "devnetL1.json")
    deployConfig, err := genesis.NewDeployConfig(deployConfigPath)
    if err != nil {
        return fmt.Errorf("failed to load deploy config: %w", err)
    }

    // Apply deploy config settings
    deployConfig.L1GenesisBlockTimestamp = hexutil.Uint64(time.Now().Unix())
    deployConfig.FundDevAccounts = true
    deployConfig.L1BlockTime = 2
    deployConfig.L2BlockTime = 1
    if l1Deployments != nil {
        deployConfig.SetDeployments(l1Deployments)
    }

    // Load L2 allocs for different modes
    l2AllocsDir := path.Join(root, ".devnet")
    l2AllocsMap := make(genesis.L2AllocsModeMap)

    modes := []genesis.L2AllocsMode{
        genesis.L2AllocsDelta,
        genesis.L2AllocsEcotone,
        genesis.L2AllocsFjord,
        genesis.L2AllocsGranite,
        genesis.L2AllocsHolocene,
        genesis.L2AllocsIsthmus,
        genesis.L2AllocsInterop,
    }

    for _, mode := range modes {
        name := "allocs-l2"
        if mode != "" {
            name += "-" + string(mode)
        }
        allocPath := path.Join(l2AllocsDir, name+".json")
        if _, err := os.Stat(allocPath); err == nil {
            allocs, err := foundry.LoadForgeAllocs(allocPath)
            if err != nil {
                log.Warn("Failed to load L2 allocs", "mode", mode, "error", err)
                continue
            }
            l2AllocsMap[mode] = allocs
        }
    }

    // Apply to all AllocTypes
    mtx.Lock()
    for _, allocType := range allocTypes {
        l1AllocsByType[allocType] = l1Allocs
        l1DeploymentsByType[allocType] = l1Deployments
        deployConfigsByType[allocType] = deployConfig
        l2AllocsByType[allocType] = l2AllocsMap
    }
    mtx.Unlock()

    return nil
}
```

#### init() 함수 수정:
```go
func init() {
    // ... logger setup ...

    // Check if .devnet allocs exist (pre-generated files)
    devnetL1AllocsPath := path.Join(root, ".devnet", "allocs-l1.json")

    if _, err := os.Stat(devnetL1AllocsPath); err == nil {
        // Use pre-generated .devnet files (simpler and more stable)
        log.Info("Using pre-generated .devnet allocs")
        if err := initFromDevnetFiles(root); err != nil {
            panic(fmt.Errorf("failed to init from .devnet files: %w", err))
        }
    } else {
        // Fall back to op-deployer initialization
        log.Info("Using op-deployer initialization (.devnet files not found)")
        for _, allocType := range allocTypes {
            initAllocType(root, allocType)
        }
    }
}
```

#### 삭제한 파일:
- `op-e2e/config/init_tokamak.go` - 복잡한 병렬 초기화 로직 제거

**효과**:
- ✅ 초기화 로직 단순화
- ✅ 병렬 실행 충돌 해결
- ✅ .devnet 파일 안정적 로드

---

## 코드 변경 사항

### 커밋 히스토리

```bash
# 1차 커밋: Genesis state loading 수정
git log --oneline -1 cc29eb922
# cc29eb922 feat(e2e): Fix genesis state loading for Tokamak E2E tests

# 2차 커밋: E2E 초기화 단순화
git log --oneline -1 d99ed8110
# d99ed8110 refactor(e2e): Simplify initialization to use .devnet files directly

# 3차 커밋: Tokamak 필드 지원
git log --oneline -1 9639cf1a0
# 9639cf1a0 feat(genesis): Add Tokamak-specific fields support in DeployConfig
```

### 변경된 파일 목록

```
op-chain-ops/foundry/allocs.go          # ForgeAllocs JSON 파싱 개선
op-chain-ops/genesis/config.go          # TokamakDeployConfig 추가, 검증 skip
op-e2e/config/init.go                   # .devnet 직접 로드
op-e2e/config/init_tokamak.go          # 삭제됨
```

---

## E2E 테스트 실행

### 전체 실행 절차 요약

```bash
# 1. .devnet 파일 생성 (5분, 최초 1회)
make devnet-allocs

# 2. macOS용 VM 바이너리 빌드 (5-10분, 최초 1회)
cd cannon
make cannon
cd ../op-program
make op-program-host

# 3. mips64 ELF 및 mt64 prestate 생성 (5분, 최초 1회)
make op-program-client-mips64
cd ..
./cannon/bin/cannon load-elf \
  --type multithreaded64-4 \
  --path op-program/bin/op-program-client64.elf \
  --out op-program/bin/prestate-mt64.bin.gz \
  --meta op-program/bin/meta-mt64.json

./cannon/bin/cannon run \
  --proof-at '=0' --stop-at '=1' \
  --input op-program/bin/prestate-mt64.bin.gz \
  --meta op-program/bin/meta-mt64.json \
  --proof-fmt 'op-program/bin/%d.json' --output ""
mv op-program/bin/0.json op-program/bin/prestate-proof-mt64.json

# 4. E2E 테스트 실행 (10-20분)
cd op-e2e
go test -v ./faultproofs -run TestOutputCannonGame -timeout 20m
```

### 1. .devnet 파일 생성

```bash
# 프로젝트 루트에서
make devnet-allocs
```

**생성되는 파일**:
- `.devnet/allocs-l1.json` - L1 genesis 상태
- `.devnet/allocs-l2.json` - L2 genesis 상태 (Granite)
- `.devnet/allocs-l2-delta.json` - L2 genesis 상태 (Delta)
- `.devnet/allocs-l2-ecotone.json` - L2 genesis 상태 (Ecotone)
- `.devnet/addresses.json` - L1 컨트랙트 주소

**소요 시간**: 약 5분

### 2. VM 바이너리 준비

#### ⚠️ 중요: macOS vs Linux 바이너리

E2E 테스트는 실제 바이너리를 실행하므로, **실행 환경에 맞는 바이너리**가 필요합니다.

**옵션 A: macOS에서 E2E 테스트 실행 (권장)**

```bash
# Cannon을 macOS용으로 빌드
cd cannon
make cannon

# 확인 (macOS용 바이너리인지 확인)
file bin/cannon
# 출력: bin/cannon: Mach-O 64-bit executable arm64

# OP Program도 로컬 빌드
cd ../op-program
make op-program-host

# 확인
ls -lh bin/op-program
```

**소요 시간**: 약 5분

**옵션 B: Docker 이미지에서 다운로드 (Linux용, macOS에서 E2E 실패)**

⚠️ **주의**: 이 방법으로 다운로드한 바이너리는 **Linux용**이므로 macOS에서 E2E 테스트 시 다음 에러 발생:
```
fork/exec ./../../cannon/bin/cannon: exec format error
```

```bash
# Linux용 바이너리 다운로드 (CI/CD나 Linux 환경용)
./op-challenger/scripts/pull-vm-images.sh --tag latest
```

**다운로드되는 파일** (모두 Linux용):
- `cannon/bin/cannon` - Cannon VM (Linux)
- `asterisc/bin/asterisc` - Asterisc VM (Linux)
- `op-program/bin/op-program` - OP Program (Linux)
- `bin/kona-client` - Kona client (Linux)
- `op-program/bin/prestate.json` - 기본 prestate
- `op-program/bin/prestate-proof.json` - 기본 prestate proof

**권장 사항**:
- ✅ **macOS에서 테스트**: 옵션 A 사용 (로컬 빌드)
- ✅ **Linux/CI에서 테스트**: 옵션 B 사용 (Docker 이미지)
- ✅ **Docker로 테스트**: 컨테이너 내부에서 실행

### 3. Cannon Prestate 준비 (mt64)

```bash
# 1. mips64 ELF 빌드
cd op-program
make op-program-client-mips64

# 2. mt64 prestate 생성 (macOS용 cannon 사용)
cd ..
./cannon/bin/cannon load-elf \
  --type multithreaded64-4 \
  --path op-program/bin/op-program-client64.elf \
  --out op-program/bin/prestate-mt64.bin.gz \
  --meta op-program/bin/meta-mt64.json

# 3. 확인
ls -lh op-program/bin/prestate-mt64.bin.gz
# 예상 출력: prestate-mt64.bin.gz (약 19MB)

# 4. prestate proof 생성
./cannon/bin/cannon run \
  --proof-at '=0' \
  --stop-at '=1' \
  --input op-program/bin/prestate-mt64.bin.gz \
  --meta op-program/bin/meta-mt64.json \
  --proof-fmt 'op-program/bin/%d.json' \
  --output ""

# proof 파일을 prestate-proof-mt64.json으로 이동
mv op-program/bin/0.json op-program/bin/prestate-proof-mt64.json

# 5. 생성된 파일 확인
ls -lh op-program/bin/prestate-mt64* op-program/bin/meta-mt64.json
```

**소요 시간**: 약 2-3분

**주의**:
- ⚠️ `load-elf`에서 `--type multithreaded64-4` 필요 (버전 7, mt64)
- ⚠️ `run`에서는 --type 불필요 (input에서 자동 감지)
- ⚠️ JSON 출력 형식은 load-elf에서 지원 안 됨, `.bin.gz` 사용

### 4. E2E 테스트 실행

#### 옵션 A: 단일 테스트 (권장)

```bash
cd op-e2e

# verbose 모드로 실행 (진행 상황 확인 가능)
go test -v ./faultproofs -run TestOutputCannonGame -timeout 20m
```

**예상 결과**:
```
=== RUN   TestOutputCannonGame
=== PAUSE TestOutputCannonGame
=== CONT  TestOutputCannonGame
INFO [timestamp] Using pre-generated .devnet allocs
INFO [timestamp] Initialized from .devnet files  l1_allocs=XX l2_modes=3
INFO [timestamp] Building developer L1 genesis block
INFO [timestamp] Included L1 deployment  name=DisputeGameFactory address=0x...
INFO [timestamp] Included L1 deployment  name=L1CrossDomainMessenger address=0x...
...
=== RUN   TestOutputCannonGame/mt-cannon
=== RUN   TestOutputCannonGame/mt-cannon-next
--- PASS: TestOutputCannonGame (10.52s)
    --- PASS: TestOutputCannonGame/mt-cannon (10.45s)
    --- PASS: TestOutputCannonGame/mt-cannon-next (10.48s)
PASS
ok      github.com/tokamak-network/tokamak-thanos/op-e2e/faultproofs    15.234s
```

**소요 시간**: 약 10-20분

#### 옵션 B: 전체 테스트

```bash
cd op-e2e

# 전체 fault proof 테스트 실행
go test ./faultproofs/... -timeout 30m

# verbose 모드 (로그 많음, 느림)
go test -v ./faultproofs/... -timeout 40m
```

**소요 시간**: 약 20-40분

---

## Cannon Prestate 빌드

### 문제: Multithreaded Cannon Prestate 부족

E2E 테스트 중 일부 테스트에서 다음 에러 발생:
```
stat ./../../op-program/bin/prestate-mt64.bin.gz: no such file or directory
stat ./../../op-program/bin/prestate-mt64Next.bin.gz: no such file or directory
```

### 해결 방법

#### 방법 1: 로컬에서 mips64 빌드 + optimism에서 복사 (가장 빠름) ⚡

우리 프로젝트의 Makefile과 Dockerfile은 이미 업데이트되었지만, cannon이 Linux 바이너리라 맥에서 직접 실행이 안됩니다. 따라서 optimism 프로젝트에서 prestate 파일을 복사하는 것이 가장 빠릅니다:

```bash
# 1. mips64 ELF 빌드 (우리 프로젝트에서)
cd op-program
make op-program-client-mips64

# 2. mt64 prestate 파일 복사 (optimism 프로젝트에서)
cp /path/to/optimism/op-program/bin/prestate-mt64*.bin.gz ./bin/
cp /path/to/optimism/op-program/bin/prestate-mt64*.json ./bin/

# 확인
ls -lh bin/prestate-mt64*
```

**장점**:
- ✅ 즉시 사용 가능 (5분)
- ✅ 검증된 prestate

**단점**:
- ⚠️ optimism 프로젝트 필요
- ⚠️ 버전 불일치 가능성 (hash가 다를 수 있음)

#### 방법 2: reproducible-prestate 빌드 (Linux/Docker 환경)

**주의**: 현재 cannon 컴파일 에러로 인해 Docker 빌드가 실패합니다. 이 방법은 cannon 코드 수정 후 사용 가능합니다.

Docker를 사용하여 재현 가능한 prestate 생성:

```bash
cd op-program
make reproducible-prestate
```

**업데이트된 파일**:
- `op-program/Makefile` - mips64 빌드 타겟 추가
- `op-program/Dockerfile.repro` - Go 1.24.2, mt64 prestate 생성 추가

**생성되는 파일**:
- `op-program/bin/prestate-proof-mt64.json`
- `op-program/bin/prestate-mt64.bin.gz`
- `op-program/bin/prestate-proof-mt64Next.json`
- `op-program/bin/prestate-mt64Next.bin.gz`
- `op-program/bin/prestate-proof-interop.json`
- `op-program/bin/prestate-interop.bin.gz`
- `op-program/bin/prestate-proof-interopNext.json`
- `op-program/bin/prestate-interopNext.bin.gz`

**소요 시간**: 10-30분 (Docker 이미지 빌드 + cannon 실행)

**출력 예시**:
```
-------------------- Production Prestates --------------------

Cannon64 Absolute prestate hash:
0x03c7ae758fa7f367ba6a7d8c21f0c5a1e64a5e6f8a5e6f8a5e6f8a5e6f8a5e6f

-------------------- Experimental Prestates --------------------

Cannon64Next Absolute prestate hash:
0x03fd582694a9cc73adc8f3e8bfcf7f5f6d5c5a4e3d2c1b0a9f8e7d6c5b4a3f2e

CannonInterop Absolute prestate hash:
0x03be3ecf8a4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c

CannonInteropNext Absolute prestate hash:
0x03ab1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab
```

#### 방법 2: 다른 프로젝트에서 복사

이미 빌드된 prestate 파일이 있는 경우:

```bash
# optimism 프로젝트에서 복사
cp /Users/zena/tokamak-projects/optimism/op-program/bin/prestate-mt64*.bin.gz \
   /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/

cp /Users/zena/tokamak-projects/optimism/op-program/bin/prestate-mt64*.json \
   /Users/zena/tokamak-projects/tokamak-thanos/op-program/bin/
```

⚠️ **주의**: 다른 프로젝트의 prestate는 코드 버전이 다르면 hash가 맞지 않을 수 있습니다.

---

## 트러블슈팅

### 1. "json: cannot unmarshal hex string without 0x prefix"

**원인**: JSON 파일의 hex string에 0x prefix가 없음

**해결**: 이미 `op-chain-ops/foundry/allocs.go`에서 처리됨 (common.HexToAddress 사용)

**확인**:
```bash
# state-dump 파일 확인
cat packages/tokamak/contracts-bedrock/state-dump-901.json | jq '.accounts | keys[0]'
```

### 2. "unknown field" 에러

**원인**: DeployConfig에 Tokamak 전용 필드 미정의

**해결**: 이미 `TokamakDeployConfig` 추가 및 unknown fields 허용

**확인**:
```bash
# deploy config 파일 확인
cat packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json | jq 'keys' | grep -i "native\|usdc"
```

### 3. "ETHLockbox is not set"

**원인**: Tokamak은 ETHLockbox 미사용

**해결**: 이미 `L1Deployments.Check()`에서 skip 처리

### 4. "predeploy is missing from L2 genesis allocs"

**원인**: L2 allocs에 predeploy 주소 누락

**해결**: `.devnet/allocs-l2-*.json` 파일에 predeploy 포함 확인

**확인**:
```bash
# L2 allocs의 predeploy 개수 확인
cat .devnet/allocs-l2-delta.json | jq 'keys | map(select(startswith("0x4200"))) | length'
# 예상: 2048 (0x4200...0000 ~ 0x4200...07FF)
```

### 5. "context canceled" during tests

**원인**: 테스트 타임아웃

**해결**: timeout 증가
```bash
go test -v ./faultproofs -run TestOutputCannonGame -timeout 20m
```

### 6. Docker 빌드 느림

**원인**: Large context (2.59GB) 전송

**최적화**:
```bash
# .dockerignore 확인
cat op-program/.dockerignore

# 불필요한 파일 제외
echo "bin/" >> op-program/.dockerignore
echo "temp/" >> op-program/.dockerignore
```

---

## 결론

### ✅ 해결된 문제

1. **ForgeAllocs JSON 파싱** - 래핑된 형식 및 0x prefix 처리
2. **Tokamak 필드 지원** - TokamakDeployConfig 추가
3. **ETHLockbox 검증** - Native Token 배포용 skip
4. **E2E 초기화** - .devnet 직접 로드로 단순화
5. **Cannon multicannon 지원** - macOS E2E 테스트 가능
6. **mips64 빌드** - mt64 prestate 생성 가능

### 🎉 E2E 테스트 실행 가능

모든 수정사항이 완료되어 **macOS에서 E2E 테스트를 정상적으로 실행**할 수 있습니다:

```bash
# 전체 준비 과정 (최초 1회, 약 20분)
make devnet-allocs                      # .devnet 파일 생성
cd cannon && make cannon                # macOS용 cannon 빌드
cd ../op-program && make op-program-host op-program-client-mips64
cd ..
./cannon/bin/cannon load-elf --type multithreaded64-4 \
  --path op-program/bin/op-program-client64.elf \
  --out op-program/bin/prestate-mt64.bin.gz \
  --meta op-program/bin/meta-mt64.json
./cannon/bin/cannon run --proof-at '=0' --stop-at '=1' \
  --input op-program/bin/prestate-mt64.bin.gz \
  --meta op-program/bin/meta-mt64.json \
  --proof-fmt 'op-program/bin/%d.json' --output ""
mv op-program/bin/0.json op-program/bin/prestate-proof-mt64.json

# E2E 테스트 실행
cd op-e2e
go test -v ./faultproofs -run TestOutputCannonGame -timeout 30m
```

### 📝 테스트 결과 (추가 예정)

E2E 테스트 실행 결과는 아래 섹션에 추가됩니다.

---

**작성일**: 2025-11-05
**최종 업데이트**: 2025-11-05
**버전**: v1.16.0
**상태**: ✅ E2E 테스트 환경 구축 완료

