# Option C 구현 계획: E2E 테스트 어댑터

**작성일**: 2025-11-04
**예상 작업 기간**: 3-5일
**난이도**: 중간
**리스크**: 🟢 낮음

## 요약

이 문서는 Tokamak-Thanos의 E2E 테스트가 Optimism v1.16.0의 배포 스크립트 대신 Tokamak의 기존 Python 기반 배포 시스템을 사용할 수 있도록 하는 어댑터의 상세 구현 계획을 제공합니다.

## 현재 상황 분석

### 문제점

**현재 초기화 흐름**:
```
op-e2e/config/init.go (init function)
    ↓
defaultIntent() - 배포 설정 생성
    ↓
deployer.ApplyPipeline() - op-deployer 사용
    ↓
DeploySuperchain.s.sol 실행 시도
    ↓
❌ "revision id 1 cannot be reverted" 에러 발생
```

**핵심 파일**: `op-e2e/config/init.go:272-285`
```go
if err := deployer.ApplyPipeline(
    context.Background(),
    deployer.ApplyPipelineOpts{
        DeploymentTarget:   deployer.DeploymentTargetGenesis,
        L1RPCUrl:           "",
        DeployerPrivateKey: pk,
        Intent:             intent,
        State:              st,
        Logger:             lgr,
        StateWriter:        pipeline.NoopStateWriter(),
    },
); err != nil {
    panic(fmt.Errorf("failed to apply pipeline: %w", err))
}
```

이 `ApplyPipeline` 호출이 Optimism의 Solidity 배포 스크립트를 실행하려고 시도하면서 에러가 발생합니다.

### Tokamak의 배포 시스템

**Tokamak Python 배포 흐름**:
```
bedrock-devnet/main.py
    ↓
bedrock-devnet/devnet/__init__.py
    ↓
Python 배포 로직 실행
    ↓
✅ contracts-bedrock/state-dump-900.json 생성
✅ contracts-bedrock/deployments/devnetL1/ 생성
```

**Python 배포의 핵심**:
- **위치**: `bedrock-devnet/devnet/__init__.py`
- **출력**:
  - L1 state dump: `packages/tokamak/contracts-bedrock/state-dump-900.json`
  - 배포 주소들: `packages/tokamak/contracts-bedrock/deployments/devnetL1/`
- **특징**: Tokamak의 커스텀 컨트랙트 (L1ContractVerification, L2NativeToken 등)을 올바르게 배포

## 해결 방안: 간소화된 접근법

원래 계획했던 복잡한 어댑터 대신, **더 간단하고 실용적인 접근법**을 사용합니다:

### 핵심 아이디어

`op-deployer`의 `ApplyPipeline` 호출을 우회하고, **이미 생성된 allocations과 배포 정보를 직접 로드**합니다.

```
현재 (실패):
op-e2e/config/init.go
    ↓
deployer.ApplyPipeline()
    ↓
❌ 배포 스크립트 실행 실패

새로운 (성공):
op-e2e/config/init.go
    ↓
Tokamak Adapter 체크
    ↓
기존 state-dump와 deployment 파일 로드
    ↓
✅ E2E 테스트 실행
```

## 아키텍처 설계

### 1. 파일 구조

```
op-e2e/
├── config/
│   ├── init.go                 (수정 필요)
│   └── tokamak_adapter.go      (신규)
└── tokamak/
    ├── README.md               (신규)
    └── deployment_loader.go    (신규)
```

### 2. 핵심 컴포넌트

#### A. 환경 변수 제어

**파일**: `op-e2e/config/init.go` (수정)

```go
// 라인 144 근처 init() 함수 시작 부분에 추가
func init() {
    // Tokamak 어댑터 활성화 확인
    if os.Getenv("TOKAMAK_USE_EXISTING_DEPLOYMENT") == "true" {
        initWithTokamakDeployment()
        return
    }

    // 기존 코드 계속...
    if os.Getenv("DISABLE_OP_E2E_LEGACY") == "true" {
        return
    }
    // ...
}
```

#### B. Tokamak 배포 로더

**파일**: `op-e2e/config/tokamak_adapter.go` (신규)

```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "path"

    "github.com/tokamak-network/tokamak-thanos/op-chain-ops/foundry"
    "github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis"
    op_service "github.com/tokamak-network/tokamak-thanos/op-service"
)

// initWithTokamakDeployment는 기존 Tokamak 배포 파일을 로드하여 초기화
func initWithTokamakDeployment() {
    cwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    root, err := op_service.FindMonorepoRoot(cwd)
    if err != nil {
        panic(err)
    }

    // Tokamak 배포 파일 경로
    contractsDir := path.Join(root, "packages", "tokamak", "contracts-bedrock")
    stateDumpPath := path.Join(contractsDir, "state-dump-900.json")
    deploymentDir := path.Join(contractsDir, "deployments", "devnetL1")
    deployConfigPath := path.Join(contractsDir, "deploy-config", "devnetL1.json")

    // 파일 존재 확인
    if err := ensureDir(contractsDir); err != nil {
        panic(fmt.Errorf("contracts directory not found: %w", err))
    }
    if _, err := os.Stat(stateDumpPath); os.IsNotExist(err) {
        panic(fmt.Errorf("state dump not found at %s. Run 'make devnet-allocs' first", stateDumpPath))
    }
    if err := ensureDir(deploymentDir); err != nil {
        panic(fmt.Errorf("deployment directory not found at %s. Run 'make devnet-allocs' first", deploymentDir))
    }

    // L1 State Dump 로드
    l1Allocs, err := loadStateDump(stateDumpPath)
    if err != nil {
        panic(fmt.Errorf("failed to load state dump: %w", err))
    }

    // L1 Deployments 로드
    l1Deployments, err := loadDeployments(deploymentDir)
    if err != nil {
        panic(fmt.Errorf("failed to load deployments: %w", err))
    }

    // Deploy Config 로드
    deployConfig, err := loadDeployConfig(deployConfigPath)
    if err != nil {
        panic(fmt.Errorf("failed to load deploy config: %w", err))
    }

    // 전역 변수에 할당 (현재 init()과 동일한 방식)
    allocType := DefaultAllocType
    mtx.Lock()
    l1AllocsByType[allocType] = l1Allocs
    l1DeploymentsByType[allocType] = l1Deployments
    deployConfigsByType[allocType] = deployConfig

    // L2 Allocs도 로드 (Tokamak의 기존 L2 genesis 파일 사용)
    l2AllocsModes := []genesis.L2AllocsMode{
        genesis.L2AllocsGranite, // 현재 Tokamak이 사용하는 모드
        // 필요한 다른 모드들 추가
    }

    l2Allocs := make(genesis.L2AllocsModeMap)
    for _, mode := range l2AllocsModes {
        l2AllocsPath := path.Join(contractsDir, fmt.Sprintf("l2-allocs-%s.json", mode))
        if allocs, err := loadStateDump(l2AllocsPath); err == nil {
            l2Allocs[mode] = allocs
        }
    }
    l2AllocsByType[allocType] = l2Allocs
    mtx.Unlock()

    log.Info("Tokamak deployment loaded successfully")
}

// loadStateDump는 state dump JSON 파일을 로드
func loadStateDump(path string) (*foundry.ForgeAllocs, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var allocs foundry.ForgeAllocs
    if err := json.Unmarshal(data, &allocs); err != nil {
        return nil, err
    }

    return &allocs, nil
}

// loadDeployments는 deployment 디렉토리에서 컨트랙트 주소를 로드
func loadDeployments(deploymentDir string) (*genesis.L1Deployments, error) {
    // .deploy 파일 읽기 (Tokamak 배포 시스템이 생성)
    deployPath := path.Join(deploymentDir, ".deploy")
    data, err := os.ReadFile(deployPath)
    if err != nil {
        return nil, err
    }

    // JSON 파싱
    var addresses map[string]string
    if err := json.Unmarshal(data, &addresses); err != nil {
        return nil, err
    }

    // L1Deployments 구조체로 변환
    deployments := &genesis.L1Deployments{
        L1CrossDomainMessengerProxy: common.HexToAddress(addresses["L1CrossDomainMessengerProxy"]),
        L1StandardBridgeProxy:       common.HexToAddress(addresses["L1StandardBridgeProxy"]),
        L2OutputOracleProxy:         common.HexToAddress(addresses["L2OutputOracleProxy"]),
        OptimismPortalProxy:         common.HexToAddress(addresses["OptimismPortalProxy"]),
        // ... 나머지 컨트랙트 주소들
    }

    return deployments, nil
}

// loadDeployConfig는 deploy config JSON 파일을 로드
func loadDeployConfig(path string) (*genesis.DeployConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config genesis.DeployConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}

// ensureDir는 디렉토리가 존재하는지 확인
func ensureDir(dir string) error {
    info, err := os.Stat(dir)
    if err != nil {
        return err
    }
    if !info.IsDir() {
        return fmt.Errorf("%s is not a directory", dir)
    }
    return nil
}
```

#### C. 사용 문서

**파일**: `op-e2e/tokamak/README.md` (신규)

```markdown
# Tokamak E2E Testing with Existing Deployment

## Overview

This adapter allows E2E tests to use Tokamak's pre-existing deployment files instead of running Optimism's deployment scripts during test initialization.

## Why?

Tokamak-Thanos uses custom contracts (L1ContractVerification, L2NativeToken) with multi-layer initialization patterns that are incompatible with Optimism v1.16.0's deployment scripts. Instead of modifying the contracts or scripts, we load pre-deployed contract state.

## Prerequisites

Before running E2E tests, generate the deployment files:

```bash
make devnet-allocs
```

This creates:
- `packages/tokamak/contracts-bedrock/state-dump-900.json` - L1 state
- `packages/tokamak/contracts-bedrock/deployments/devnetL1/` - Contract addresses

## Usage

### Enable Tokamak Adapter

Set the environment variable:

```bash
export TOKAMAK_USE_EXISTING_DEPLOYMENT=true
```

### Run Tests

```bash
go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs -timeout 10m
```

### Disable Adapter (Use Standard Optimism Flow)

```bash
unset TOKAMAK_USE_EXISTING_DEPLOYMENT
# OR
export TOKAMAK_USE_EXISTING_DEPLOYMENT=false
```

## Architecture

```
E2E Test Start
    ↓
Check TOKAMAK_USE_EXISTING_DEPLOYMENT env var
    ↓
If "true":
    ├─ Load state-dump-900.json
    ├─ Load deployments/devnetL1/
    ├─ Load deploy-config/devnetL1.json
    └─ Populate global test variables
    ↓
If "false":
    └─ Use standard op-deployer pipeline (may fail)
    ↓
E2E Test Execution
```

## Files Loaded

1. **L1 State Dump**: `packages/tokamak/contracts-bedrock/state-dump-900.json`
   - Contains L1 genesis state with all deployed contracts

2. **Contract Addresses**: `packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy`
   - JSON mapping of contract names to addresses

3. **Deploy Config**: `packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json`
   - Configuration parameters for the deployment

## Troubleshooting

### Error: "state dump not found"

**Solution**: Run `make devnet-allocs` to generate deployment files.

### Error: "deployment directory not found"

**Solution**: Verify that Python deployment completed successfully:
```bash
cd bedrock-devnet
python3 main.py --allocs
```

### Tests still using Optimism deployment

**Solution**: Verify environment variable is set:
```bash
echo $TOKAMAK_USE_EXISTING_DEPLOYMENT
# Should print: true
```

## Advantages

✅ **No Contract Modifications**: Tokamak contracts remain unchanged
✅ **No Script Modifications**: Optimism deployment scripts untouched
✅ **Fast Test Startup**: No deployment during test init
✅ **Low Risk**: Uses production deployment system
✅ **Easy Rollback**: Just unset the environment variable

## Limitations

⚠️ **Requires Pre-deployment**: Must run `make devnet-allocs` first
⚠️ **Static State**: Uses fixed deployment, not dynamic
⚠️ **Cache Invalidation**: Re-run `make devnet-allocs` after contract changes

## Development

### Adding New Tests

New E2E tests automatically use the adapter when the environment variable is set. No additional code needed.

### Updating Deployment Files

After modifying contracts:

```bash
# Clean old deployment
rm -rf packages/tokamak/contracts-bedrock/deployments/devnetL1
rm packages/tokamak/contracts-bedrock/state-dump-900.json

# Generate new deployment
make devnet-allocs

# Run tests
export TOKAMAK_USE_EXISTING_DEPLOYMENT=true
go test -v ./op-e2e/faultproofs
```

## See Also

- [E2E Deployment Compatibility Analysis](../migration-v1.16.0/E2E-DEPLOYMENT-COMPATIBILITY-ANALYSIS-KR.md)
- [Migration Guide](../migration-v1.16.0/MIGRATION-GUIDE.md)
```

## 구현 단계

### Phase 1: 기본 구조 생성 (1시간)

**Task 1.1**: 디렉토리 생성
```bash
mkdir -p op-e2e/tokamak
```

**Task 1.2**: `tokamak_adapter.go` 파일 생성
- 기본 구조 작성
- `initWithTokamakDeployment()` 함수 스켈레톤

**Task 1.3**: `README.md` 작성
- 사용 방법 문서화

### Phase 2: 로더 함수 구현 (2-3시간)

**Task 2.1**: `loadStateDump()` 구현
- JSON 파일 읽기
- `ForgeAllocs` 구조체로 파싱
- 에러 처리

**Task 2.2**: `loadDeployments()` 구현
- `.deploy` 파일 읽기
- 주소 매핑
- `L1Deployments` 구조체 생성

**Task 2.3**: `loadDeployConfig()` 구현
- Deploy config JSON 로드
- `DeployConfig` 구조체로 파싱

### Phase 3: init.go 수정 (1시간)

**Task 3.1**: 환경 변수 체크 추가
```go
// op-e2e/config/init.go:144 근처
func init() {
    if os.Getenv("TOKAMAK_USE_EXISTING_DEPLOYMENT") == "true" {
        initWithTokamakDeployment()
        return
    }
    // 기존 코드...
}
```

**Task 3.2**: Import 추가
```go
import (
    // 기존 imports...
    "os" // 이미 있을 수 있음
)
```

### Phase 4: Python 배포 실행 (30분)

**Task 4.1**: 배포 파일 생성
```bash
make devnet-allocs
```

**Task 4.2**: 파일 확인
```bash
ls -la packages/tokamak/contracts-bedrock/state-dump-900.json
ls -la packages/tokamak/contracts-bedrock/deployments/devnetL1/
```

### Phase 5: 테스트 (2-3시간)

**Task 5.1**: 컴파일 테스트
```bash
go build ./op-e2e/config
```

**Task 5.2**: 간단한 유닛 테스트 작성
```bash
# op-e2e/config/tokamak_adapter_test.go
```

**Task 5.3**: E2E 테스트 실행
```bash
export TOKAMAK_USE_EXISTING_DEPLOYMENT=true
go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs -timeout 10m
```

**Task 5.4**: 결과 검증
- 에러 로그 확인
- 배포 주소 검증
- 테스트 통과 확인

### Phase 6: 문서화 및 정리 (1-2시간)

**Task 6.1**: 코드 주석 추가

**Task 6.2**: README 업데이트

**Task 6.3**: 마이그레이션 가이드 업데이트

## 주요 파일 및 위치

### 생성할 파일

1. **`op-e2e/config/tokamak_adapter.go`** (~200-250 lines)
   - `initWithTokamakDeployment()` - 메인 초기화 함수
   - `loadStateDump()` - state dump 로더
   - `loadDeployments()` - 배포 주소 로더
   - `loadDeployConfig()` - config 로더
   - `ensureDir()` - 유틸리티

2. **`op-e2e/tokamak/README.md`** (~150 lines)
   - 사용 방법
   - 아키텍처 설명
   - 트러블슈팅

### 수정할 파일

1. **`op-e2e/config/init.go`** (최소 수정)
   - 라인 144-148: 환경 변수 체크 추가
   - Import 추가 (필요시)

### 사용할 기존 파일

1. **`packages/tokamak/contracts-bedrock/state-dump-900.json`**
   - Python 배포가 생성
   - L1 genesis state

2. **`packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy`**
   - Python 배포가 생성
   - 컨트랙트 주소 매핑

3. **`packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json`**
   - 기존 파일
   - 배포 설정

## 예상 문제 및 해결 방안

### 문제 1: State Dump 형식 불일치

**증상**: `loadStateDump()`에서 JSON 파싱 에러

**원인**: Tokamak의 state dump 형식이 `ForgeAllocs`와 다를 수 있음

**해결**:
```go
// 중간 구조체 사용
type TokamakStateDump struct {
    // Tokamak 특화 필드
}

func loadStateDump(path string) (*foundry.ForgeAllocs, error) {
    var tokamakDump TokamakStateDump
    // 먼저 Tokamak 형식으로 로드
    json.Unmarshal(data, &tokamakDump)

    // ForgeAllocs로 변환
    allocs := convertToForgeAllocs(tokamakDump)
    return allocs, nil
}
```

### 문제 2: 주소 매핑 불완전

**증상**: 일부 컨트랙트 주소가 누락

**원인**: `.deploy` 파일에 모든 주소가 없을 수 있음

**해결**:
```go
// 개별 컨트랙트 JSON 파일도 읽기
func loadDeployments(deploymentDir string) (*genesis.L1Deployments, error) {
    // .deploy 파일 먼저 시도
    deployments, err := loadFromDeployFile(deploymentDir)
    if err == nil {
        return deployments, nil
    }

    // 실패하면 개별 JSON 파일들 읽기
    return loadFromIndividualFiles(deploymentDir)
}
```

### 문제 3: L2 Allocs 누락

**증상**: L2 genesis 생성 시 에러

**원인**: L2 allocs 파일이 없음

**해결**:
```go
// L2 allocs가 없으면 동적 생성
func generateL2Allocs(mode genesis.L2AllocsMode, deployConfig *genesis.DeployConfig) (*foundry.ForgeAllocs, error) {
    // genesis.BuildL2Genesis() 사용하여 생성
}
```

## 성공 기준

### Phase 1-3 완료 (첫날 끝)
- ✅ 코드 컴파일 성공
- ✅ 환경 변수 체크 작동
- ✅ 파일 로드 에러 없음

### Phase 4-5 완료 (둘째날 끝)
- ✅ `make devnet-allocs` 성공
- ✅ State dump 로드 성공
- ✅ Deployments 로드 성공
- ✅ E2E 테스트가 초기화 단계를 통과

### Phase 6 완료 (셋째날 끝)
- ✅ 최소 1개 E2E 테스트 통과
- ✅ 문서 완성
- ✅ 코드 리뷰 준비 완료

## 테스트 계획

### 단위 테스트

**파일**: `op-e2e/config/tokamak_adapter_test.go`

```go
package config

import (
    "os"
    "path"
    "testing"
)

func TestLoadStateDump(t *testing.T) {
    // 테스트 데이터 준비
    testData := `{"0x...": {"balance": "0x1"}}`
    tmpfile := writeTempFile(t, testData)
    defer os.Remove(tmpfile)

    // 로드 테스트
    allocs, err := loadStateDump(tmpfile)
    if err != nil {
        t.Fatalf("loadStateDump failed: %v", err)
    }

    // 검증
    if len(*allocs) == 0 {
        t.Error("Expected non-empty allocs")
    }
}

func TestLoadDeployments(t *testing.T) {
    // 유사하게 구현
}

func TestInitWithTokamakDeployment(t *testing.T) {
    // 통합 테스트
}
```

### 통합 테스트

```bash
# 1. 배포 파일 생성
make devnet-allocs

# 2. 환경 변수 설정
export TOKAMAK_USE_EXISTING_DEPLOYMENT=true

# 3. E2E 테스트 실행
go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs -timeout 10m

# 4. 결과 확인
# - "Tokamak deployment loaded successfully" 로그 확인
# - 테스트가 초기화 단계를 통과하는지 확인
```

## 롤백 계획

문제 발생 시:

1. **환경 변수 제거**:
   ```bash
   unset TOKAMAK_USE_EXISTING_DEPLOYMENT
   ```

2. **코드 되돌리기**:
   ```bash
   git checkout HEAD -- op-e2e/config/init.go
   rm op-e2e/config/tokamak_adapter.go
   rm -rf op-e2e/tokamak/
   ```

3. **원래 상태로 테스트**:
   ```bash
   go test -v ./op-e2e/config
   ```

## 다음 단계 (완료 후)

1. **더 많은 E2E 테스트 실행**
   - `TestOutputCannonGame_*`
   - `TestOutputAlphabetGame_*`
   - 기타 faultproofs 테스트들

2. **성능 최적화**
   - 파일 로드 캐싱
   - 병렬 로딩

3. **CI/CD 통합**
   - GitHub Actions 워크플로우 업데이트
   - 자동화된 테스트 실행

4. **프로덕션 준비**
   - 에러 처리 강화
   - 로깅 개선
   - 모니터링 추가

## 참고 자료

- [E2E Deployment Compatibility Analysis](./E2E-DEPLOYMENT-COMPATIBILITY-ANALYSIS-KR.md)
- [Migration Guide](./MIGRATION-GUIDE.md)
- [Deployment Scripts Comparison](./DEPLOYMENT-SCRIPTS-COMPARISON.md)

## 체크리스트

### 구현 전
- [ ] Python 배포 실행: `make devnet-allocs`
- [ ] 배포 파일 확인
- [ ] 현재 E2E 테스트 에러 재현

### 구현 중
- [ ] Phase 1: 기본 구조 생성
- [ ] Phase 2: 로더 함수 구현
- [ ] Phase 3: init.go 수정
- [ ] Phase 4: Python 배포 실행
- [ ] Phase 5: 테스트
- [ ] Phase 6: 문서화

### 구현 후
- [ ] 최소 1개 E2E 테스트 통과
- [ ] 코드 리뷰
- [ ] 문서 검토
- [ ] PR 생성

---

**문서 버전**: 1.0
**최종 업데이트**: 2025-11-04
**다음 리뷰**: Phase 5 완료 후
