# Option C 구현 계획 (수정판): E2E 테스트 어댑터

**작성일**: 2025-11-05
**수정일**: 2025-11-05
**예상 작업 기간**: 4-6일
**난이도**: 중상
**리스크**: 🟡 중간

## ⚠️ 주요 변경사항 (원본 대비)

1. **State Dump 파일 생성 문제 해결 추가**
   - 현재 `make devnet-allocs`가 파일을 생성하지 않음
   - Forge 스크립트 수정 또는 대안 방법 필요

2. **L1Deployments 완전 매핑**
   - 26개 필드 전체 매핑 코드 추가
   - Tokamak 특화 필드 (RAT, ETHLockbox 등) 포함

3. **다중 AllocType 지원**
   - AllocTypeMTCannon, AllocTypeMTCannonNext, AllocTypeAltDA 모두 지원
   - 병렬 처리 유지

4. **실제 파일 경로 검증**
   - 파일 존재 여부 확인 후 진행
   - 없으면 Optimism 방식으로 fallback

## 현재 상황 분석

### 문제점

**현재 초기화 흐름**:
```
op-e2e/config/init.go (init function)
    ↓
for each allocType (병렬 실행):
    - defaultIntent() 생성
    - deployer.ApplyPipeline() 호출
    ↓
❌ "revision id 1 cannot be reverted" 에러 발생
```

**핵심 파일**: `op-e2e/config/init.go:189-285`

### Tokamak의 배포 시스템

**Tokamak 배포 흐름**:
```
bedrock-devnet/main.py
    ↓
Forge 스크립트 실행
    ↓
packages/tokamak/contracts-bedrock/scripts/...
    ↓
❓ state-dump-901.json 생성 (현재는 생성 안 됨)
✅ deployments/devnetL1/.deploy 생성 (정상 작동)
```

**현재 상태**:
- `.deploy` 파일: ✅ 정상 생성됨 (37개 컨트랙트 주소)
- `state-dump-*.json`: ❌ 생성되지 않음
- `deploy-config/devnetL1.json`: ✅ 존재함

## 해결 방안: 단계적 접근법

### 접근법 A: 기존 op-deployer 활용 (권장)

**핵심 아이디어**:
`op-deployer`의 `ApplyPipeline`이 이미 L1 allocs와 deployments를 생성하므로, 이를 재사용하되 **Tokamak의 .deploy 파일로 주소를 덮어쓰기**

```
init() 실행
    ↓
deployer.ApplyPipeline() 실행
    ↓ (정상 완료)
st.Chains[0].Allocs.Data → L1 Allocs
st.Chains[0].L1Deployments → L1 Deployments
    ↓
if TOKAMAK_OVERRIDE_ADDRESSES == "true":
    Tokamak의 .deploy 파일 읽기
    L1Deployments 주소들 업데이트
    ↓
전역 변수에 저장
```

**장점**:
- State dump 파일 생성 문제 우회
- Optimism의 정상 작동하는 로직 활용
- 주소만 교체하므로 리스크 낮음
- 기존 코드 최소 수정

**단점**:
- `ApplyPipeline`이 여전히 실행되어야 함 (시간 소요)
- Tokamak 주소와 Optimism 주소가 다르면 state dump와 불일치 가능

### 접근법 B: 완전 우회 (원본 계획)

State dump 파일을 로드하여 완전히 우회. 단, 파일 생성 문제 해결 필요.

## 권장 구현: 접근법 A

### 아키텍처

```
E2E Test Start
    ↓
Check TOKAMAK_OVERRIDE_ADDRESSES env var
    ↓
If "false" or unset:
    └─ 기존 ApplyPipeline 실행 (Optimism 방식)
    ↓
If "true":
    ├─ ApplyPipeline 실행
    ├─ Load Tokamak .deploy file
    ├─ Override L1Deployments addresses
    └─ Keep L1 Allocs from pipeline
    ↓
E2E Test Execution
```

### 파일 구조

```
op-e2e/
├── config/
│   ├── init.go                    (수정)
│   └── tokamak_override.go        (신규)
└── docs/
    └── tokamak/
        └── README.md              (신규)
```

## 상세 구현

### 1. Tokamak Override 로더

**파일**: `op-e2e/config/tokamak_override.go` (신규)

```go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis"
	op_service "github.com/tokamak-network/tokamak-thanos/op-service"
)

// OverrideTokamakAddresses loads Tokamak deployment addresses and overrides the L1Deployments
func OverrideTokamakAddresses(deployments *genesis.L1Deployments) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	root, err := op_service.FindMonorepoRoot(cwd)
	if err != nil {
		return fmt.Errorf("failed to find monorepo root: %w", err)
	}

	deploymentDir := path.Join(root, "packages", "tokamak", "contracts-bedrock", "deployments", "devnetL1")
	deployPath := path.Join(deploymentDir, ".deploy")

	// Check if file exists
	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		return fmt.Errorf("Tokamak deployment file not found at %s. Run 'make devnet-allocs' first", deployPath)
	}

	// Load .deploy file
	data, err := os.ReadFile(deployPath)
	if err != nil {
		return fmt.Errorf("failed to read deployment file: %w", err)
	}

	var addresses map[string]string
	if err := json.Unmarshal(data, &addresses); err != nil {
		return fmt.Errorf("failed to parse deployment file: %w", err)
	}

	// Override addresses
	// Helper function to safely get address
	getAddr := func(name string) common.Address {
		if addr, ok := addresses[name]; ok {
			return common.HexToAddress(addr)
		}
		return common.Address{}
	}

	// Map all 26+ fields from L1Deployments struct
	if addr := getAddr("AddressManager"); addr != (common.Address{}) {
		deployments.AddressManager = addr
	}
	if addr := getAddr("DisputeGameFactory"); addr != (common.Address{}) {
		deployments.DisputeGameFactory = addr
	}
	if addr := getAddr("DisputeGameFactoryProxy"); addr != (common.Address{}) {
		deployments.DisputeGameFactoryProxy = addr
	}
	if addr := getAddr("L1CrossDomainMessenger"); addr != (common.Address{}) {
		deployments.L1CrossDomainMessenger = addr
	}
	if addr := getAddr("L1CrossDomainMessengerProxy"); addr != (common.Address{}) {
		deployments.L1CrossDomainMessengerProxy = addr
	}
	if addr := getAddr("L1ERC721Bridge"); addr != (common.Address{}) {
		deployments.L1ERC721Bridge = addr
	}
	if addr := getAddr("L1ERC721BridgeProxy"); addr != (common.Address{}) {
		deployments.L1ERC721BridgeProxy = addr
	}
	if addr := getAddr("L1StandardBridge"); addr != (common.Address{}) {
		deployments.L1StandardBridge = addr
	}
	if addr := getAddr("L1StandardBridgeProxy"); addr != (common.Address{}) {
		deployments.L1StandardBridgeProxy = addr
	}
	if addr := getAddr("L2OutputOracle"); addr != (common.Address{}) {
		deployments.L2OutputOracle = addr
	}
	if addr := getAddr("L2OutputOracleProxy"); addr != (common.Address{}) {
		deployments.L2OutputOracleProxy = addr
	}
	if addr := getAddr("OptimismMintableERC20Factory"); addr != (common.Address{}) {
		deployments.OptimismMintableERC20Factory = addr
	}
	if addr := getAddr("OptimismMintableERC20FactoryProxy"); addr != (common.Address{}) {
		deployments.OptimismMintableERC20FactoryProxy = addr
	}
	if addr := getAddr("OptimismPortal"); addr != (common.Address{}) {
		deployments.OptimismPortal = addr
	}
	if addr := getAddr("OptimismPortalProxy"); addr != (common.Address{}) {
		deployments.OptimismPortalProxy = addr
	}
	// Tokamak-specific contracts
	if addr := getAddr("ETHLockbox"); addr != (common.Address{}) {
		deployments.ETHLockbox = addr
	}
	if addr := getAddr("ETHLockboxProxy"); addr != (common.Address{}) {
		deployments.ETHLockboxProxy = addr
	}
	if addr := getAddr("ProxyAdmin"); addr != (common.Address{}) {
		deployments.ProxyAdmin = addr
	}
	if addr := getAddr("SystemConfig"); addr != (common.Address{}) {
		deployments.SystemConfig = addr
	}
	if addr := getAddr("SystemConfigProxy"); addr != (common.Address{}) {
		deployments.SystemConfigProxy = addr
	}
	if addr := getAddr("ProtocolVersions"); addr != (common.Address{}) {
		deployments.ProtocolVersions = addr
	}
	if addr := getAddr("ProtocolVersionsProxy"); addr != (common.Address{}) {
		deployments.ProtocolVersionsProxy = addr
	}
	if addr := getAddr("DataAvailabilityChallenge"); addr != (common.Address{}) {
		deployments.DataAvailabilityChallenge = addr
	}
	if addr := getAddr("DataAvailabilityChallengeProxy"); addr != (common.Address{}) {
		deployments.DataAvailabilityChallengeProxy = addr
	}
	// RAT is Tokamak-specific
	if addr := getAddr("RAT"); addr != (common.Address{}) {
		deployments.RAT = addr
	}
	if addr := getAddr("RATProxy"); addr != (common.Address{}) {
		deployments.RATProxy = addr
	}

	return nil
}

// LogDeploymentAddresses logs the deployment addresses for debugging
func LogDeploymentAddresses(deployments *genesis.L1Deployments) {
	if os.Getenv("DEBUG_DEPLOYMENTS") == "true" {
		fmt.Printf("L1 Deployment Addresses:\n")
		fmt.Printf("  AddressManager: %s\n", deployments.AddressManager.Hex())
		fmt.Printf("  DisputeGameFactoryProxy: %s\n", deployments.DisputeGameFactoryProxy.Hex())
		fmt.Printf("  L1CrossDomainMessengerProxy: %s\n", deployments.L1CrossDomainMessengerProxy.Hex())
		fmt.Printf("  L1StandardBridgeProxy: %s\n", deployments.L1StandardBridgeProxy.Hex())
		fmt.Printf("  OptimismPortalProxy: %s\n", deployments.OptimismPortalProxy.Hex())
		fmt.Printf("  SystemConfigProxy: %s\n", deployments.SystemConfigProxy.Hex())
		// Add more as needed
	}
}
```

### 2. init.go 수정

**파일**: `op-e2e/config/init.go` (수정)

```go
// 라인 189 근처 initAllocType 함수 끝 부분에 추가
func initAllocType(root string, allocType AllocType) {
	// ... 기존 코드 (line 197-285) ...

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

	// 새로 추가: Tokamak 주소 override
	if os.Getenv("TOKAMAK_OVERRIDE_ADDRESSES") == "true" {
		l1Deployments := inspect.GenesisSystemConfigFromState(st)
		if err := OverrideTokamakAddresses(l1Deployments); err != nil {
			lgr.Warn("Failed to override Tokamak addresses, using Optimism addresses", "error", err)
		} else {
			lgr.Info("Successfully overridden L1 deployment addresses with Tokamak deployments")
			LogDeploymentAddresses(l1Deployments)
		}
	}

	mtx.Lock()
	l2Alloc[mode] = st.Chains[0].Allocs.Data
	// ... 기존 코드 계속 ...
```

### 3. README 문서

**파일**: `op-e2e/docs/tokamak/README.md` (신규)

```markdown
# Tokamak E2E Testing with Address Override

## Overview

This adapter allows E2E tests to use Tokamak's deployed contract addresses instead of Optimism's op-deployer generated addresses.

## Why?

Tokamak-Thanos uses custom contracts with multi-layer initialization that may generate different addresses than Optimism's deployment scripts. This adapter overrides the addresses while keeping the genesis state from op-deployer.

## Prerequisites

Generate Tokamak deployment files:

```bash
make devnet-allocs
```

This creates:
- `packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy` - Contract addresses

## Usage

### Enable Address Override

```bash
export TOKAMAK_OVERRIDE_ADDRESSES=true
```

### Run Tests

```bash
go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs -timeout 10m
```

### Debug Mode

To see which addresses are being used:

```bash
export DEBUG_DEPLOYMENTS=true
export TOKAMAK_OVERRIDE_ADDRESSES=true
go test -v ./op-e2e/faultproofs
```

### Disable Override

```bash
unset TOKAMAK_OVERRIDE_ADDRESSES
# OR
export TOKAMAK_OVERRIDE_ADDRESSES=false
```

## How It Works

```
Test Initialization
    ↓
Run op-deployer ApplyPipeline
    ├─ Generate L1 genesis state
    ├─ Generate L1 deployment addresses (Optimism)
    └─ Generate L2 allocs
    ↓
If TOKAMAK_OVERRIDE_ADDRESSES=true:
    ├─ Load Tokamak .deploy file
    ├─ Override L1Deployments addresses
    └─ Keep L1 Allocs unchanged
    ↓
Test Execution with Tokamak addresses
```

## Troubleshooting

### Error: "Tokamak deployment file not found"

**Solution**: Run `make devnet-allocs` to generate deployment files.

```bash
make devnet-allocs
```

### Addresses don't match expected values

**Solution**: Enable debug mode to see which addresses are loaded:

```bash
export DEBUG_DEPLOYMENTS=true
cat packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy | jq
```

### Tests fail with address-related errors

**Cause**: State dump may not match overridden addresses.

**Solution**: This approach works best when:
1. Contract bytecode is the same between Tokamak and Optimism
2. Only deployment addresses differ
3. State layouts are identical

If tests fail, try disabling override:
```bash
unset TOKAMAK_OVERRIDE_ADDRESSES
```

## Advantages

✅ **No Pipeline Changes**: Uses existing op-deployer logic
✅ **Fast Implementation**: Only ~150 lines of code
✅ **Low Risk**: Doesn't modify core deployment logic
✅ **Easy Rollback**: Just unset environment variable
✅ **Works Without State Dumps**: Uses ApplyPipeline's state

## Limitations

⚠️ **Requires Pre-deployment**: Must run `make devnet-allocs` first
⚠️ **Address Override Only**: State dump remains from op-deployer
⚠️ **Assumes Compatible State**: Works if contract state layouts match

## Alternative: Full State Loading

If address override doesn't work (state incompatibility), see:
- `OPTION-C-IMPLEMENTATION-PLAN-KR.md` - Original plan with full state loading
- Requires fixing state-dump file generation first

## See Also

- [E2E Deployment Compatibility Analysis](../migration-v1.16.0/E2E-DEPLOYMENT-COMPATIBILITY-ANALYSIS-KR.md)
- [Migration Guide](../migration-v1.16.0/MIGRATION-GUIDE.md)
```

## 구현 단계

### Phase 1: Override 로직 구현 (1일)

**Task 1.1**: `tokamak_override.go` 작성
- `OverrideTokamakAddresses()` 함수
- `LogDeploymentAddresses()` 함수
- 26개 필드 전부 매핑

**Task 1.2**: 컴파일 테스트
```bash
go build ./op-e2e/config
```

### Phase 2: init.go 통합 (0.5일)

**Task 2.1**: init.go 수정
- 환경 변수 체크 추가
- Override 함수 호출 추가

**Task 2.2**: 컴파일 검증
```bash
go build ./op-e2e
```

### Phase 3: 기본 테스트 (1일)

**Task 3.1**: 배포 파일 생성
```bash
make devnet-allocs
ls -la packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy
```

**Task 3.2**: Override 테스트
```bash
export TOKAMAK_OVERRIDE_ADDRESSES=true
export DEBUG_DEPLOYMENTS=true
go test -v ./op-e2e/config -run TestL1Allocs
```

### Phase 4: E2E 테스트 (2일)

**Task 4.1**: 간단한 E2E 테스트
```bash
export TOKAMAK_OVERRIDE_ADDRESSES=true
go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs -timeout 10m
```

**Task 4.2**: 결과 분석
- 주소가 올바르게 override 되었는지 확인
- State 불일치 문제 있는지 확인

### Phase 5: 문제 해결 및 조정 (1-2일)

**Task 5.1**: 발견된 문제 수정
- State 불일치 시 대응
- 누락된 주소 추가

**Task 5.2**: 추가 테스트
```bash
go test -v ./op-e2e/faultproofs -timeout 30m
```

### Phase 6: 문서화 (0.5일)

**Task 6.1**: README 작성
**Task 6.2**: 코드 주석 추가
**Task 6.3**: 마이그레이션 가이드 업데이트

## 예상 문제 및 해결

### 문제 1: State Dump와 주소 불일치

**증상**: 테스트에서 "contract not found at address" 에러

**원인**: ApplyPipeline의 state dump는 Optimism 주소를 사용하지만, 우리는 Tokamak 주소로 override

**해결방안**:
1. **Option A**: State dump도 Tokamak 주소로 업데이트
2. **Option B**: 완전한 state loading 방식으로 전환 (원본 계획)

### 문제 2: .deploy 파일 파싱 에러

**증상**: JSON unmarshal 에러

**해결**: `.deploy` 파일 형식 확인 후 파서 조정

### 문제 3: 일부 주소 누락

**증상**: 특정 컨트랙트 주소가 zero address

**해결**:
```go
// Optional 주소는 zero address 허용
if addr := getAddr("OptionalContract"); addr != (common.Address{}) {
    deployments.OptionalContract = addr
}
```

## 롤백 계획

문제 발생 시:

1. **환경 변수 제거**:
   ```bash
   unset TOKAMAK_OVERRIDE_ADDRESSES
   ```

2. **코드 되돌리기**:
   ```bash
   git checkout HEAD -- op-e2e/config/init.go
   rm op-e2e/config/tokamak_override.go
   ```

## 대안: State Dump 파일 생성 수정

만약 접근법 A가 작동하지 않으면:

1. **Forge 스크립트 수정**:
   - `packages/tokamak/contracts-bedrock/scripts/`에서 state dump 생성 부분 찾기
   - 왜 파일이 생성되지 않는지 디버깅
   - 파일 생성 로직 수정

2. **원본 계획 (Option C) 실행**:
   - State dump 파일 생성 후
   - 완전한 state loading 구현

## 성공 기준

### Phase 1-2 완료
- ✅ 코드 컴파일 성공
- ✅ Override 함수 작동
- ✅ 주소가 올바르게 로드됨

### Phase 3-4 완료
- ✅ `make devnet-allocs` 성공
- ✅ .deploy 파일 로드 성공
- ✅ E2E 테스트 초기화 통과

### Phase 5-6 완료
- ✅ 최소 1개 E2E 테스트 통과
- ✅ 문서 완성
- ✅ 코드 리뷰 준비 완료

## 참고 자료

- [Original Option C Plan](./OPTION-C-IMPLEMENTATION-PLAN-KR.md)
- [E2E Deployment Compatibility Analysis](./E2E-DEPLOYMENT-COMPATIBILITY-ANALYSIS-KR.md)

---

**문서 버전**: 2.0
**최종 업데이트**: 2025-11-05
**다음 리뷰**: Phase 3 완료 후
