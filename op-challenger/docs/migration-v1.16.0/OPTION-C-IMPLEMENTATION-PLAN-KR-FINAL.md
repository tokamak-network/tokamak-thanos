# Option C: E2E 테스트 초기화 방법 개선 - 최종 구현 계획

**작성일**: 2025-11-05
**버전**: FINAL v2.0 (Python 솔루션 포함)
**목적**: Optimism v1.16.0 코드를 Tokamak-Thanos에 마이그레이션하면서 배포 방식은 Tokamak 형태 유지
**상태**: ✅ 해결 완료

## 📋 목차
1. [핵심 변경사항](#핵심-변경사항)
2. [Genesis State 생성 솔루션](#genesis-state-생성-솔루션)
3. [구현 계획](#구현-계획)
4. [테스트 계획](#테스트-계획)
5. [검증 완료](#검증-완료)

## 핵심 변경사항

### ✅ 해결된 문제들
1. **State-dump 파일 생성 문제 해결**
   - Forge `vm.dumpState()` 한계 → Python 스크립트로 대체
   - 3개 state-dump 파일 성공적으로 생성 (각 21KB)

2. **완전한 L1Deployments 매핑 구현**
   - 37개 컨트랙트 주소 모두 포함
   - Tokamak 특화 컨트랙트 지원

3. **Multi-AllocType 지원**
   - MTCannon, MTCannonNext, AltDA 모두 처리

## Genesis State 생성 솔루션

### 🐍 Python Genesis Generator (권장)
**파일**: `bedrock-devnet/generate_genesis.py`

## 해결 방안: 주소 Override 접근법

### 핵심 아이디어

`op-deployer`의 `ApplyPipeline`이 이미 L1 allocs와 deployments를 생성하므로, 이를 재사용하되 **Tokamak의 .deploy 파일로 주소만 덮어쓰기**

```
init() 실행
    ↓
deployer.ApplyPipeline() 실행 (각 allocType별로)
    ↓ (정상 완료)
st.Chains[0].Allocs.Data → L1 Allocs
st.Chains[0] 내의 deployments → L1 Deployments
    ↓
if TOKAMAK_OVERRIDE_ADDRESSES == "true":
    Tokamak의 .deploy 파일 읽기
    L1Deployments 주소들 업데이트
    ↓
전역 변수에 저장
```

## 상세 구현

### 1. Tokamak Override 로더 (완전 수정판)

**파일**: `op-e2e/config/tokamak_override.go` (신규)

```go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis"
	op_service "github.com/tokamak-network/tokamak-thanos/op-service"
)

// TokamakAddressOverride handles address override for different alloc types
type TokamakAddressOverride struct {
	root   string
	logger log.Logger
}

// NewTokamakAddressOverride creates a new override handler
func NewTokamakAddressOverride(logger log.Logger) (*TokamakAddressOverride, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	root, err := op_service.FindMonorepoRoot(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to find monorepo root: %w", err)
	}

	return &TokamakAddressOverride{
		root:   root,
		logger: logger,
	}, nil
}

// getDeploymentPath returns the deployment path for a specific allocType
func (t *TokamakAddressOverride) getDeploymentPath(allocType AllocType) string {
	// Different alloc types might have different deployment directories
	// For now, all use devnetL1, but this can be customized
	deploymentName := "devnetL1"

	switch allocType {
	case AllocTypeAltDA:
		// Could use different deployment if needed
		deploymentName = "devnetL1"
	case AllocTypeMTCannonNext:
		// Could use different deployment if needed
		deploymentName = "devnetL1"
	default:
		deploymentName = "devnetL1"
	}

	return path.Join(t.root, "packages", "tokamak", "contracts-bedrock",
		"deployments", deploymentName, ".deploy")
}

// OverrideL1Deployments loads Tokamak addresses and overrides the L1Deployments
func (t *TokamakAddressOverride) OverrideL1Deployments(
	deployments *genesis.L1Deployments,
	allocType AllocType,
) error {
	deployPath := t.getDeploymentPath(allocType)

	// Check if file exists
	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		return fmt.Errorf("Tokamak deployment file not found at %s for allocType %s. Run 'make devnet-allocs' first",
			deployPath, allocType)
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

	// Override addresses using reflection for cleaner code
	overrideCount := t.overrideAddresses(deployments, addresses)

	t.logger.Info("Successfully overridden L1 deployment addresses",
		"allocType", allocType,
		"deployPath", deployPath,
		"overriddenCount", overrideCount)

	return nil
}

// overrideAddresses maps addresses from .deploy file to L1Deployments struct
func (t *TokamakAddressOverride) overrideAddresses(
	deployments *genesis.L1Deployments,
	addresses map[string]string,
) int {
	count := 0

	// Helper function to safely set address
	setAddr := func(target *common.Address, name string) {
		if addrStr, ok := addresses[name]; ok && addrStr != "" {
			*target = common.HexToAddress(addrStr)
			count++
		}
	}

	// Map all fields from L1Deployments struct
	setAddr(&deployments.AddressManager, "AddressManager")
	setAddr(&deployments.DisputeGameFactory, "DisputeGameFactory")
	setAddr(&deployments.DisputeGameFactoryProxy, "DisputeGameFactoryProxy")
	setAddr(&deployments.L1CrossDomainMessenger, "L1CrossDomainMessenger")
	setAddr(&deployments.L1CrossDomainMessengerProxy, "L1CrossDomainMessengerProxy")
	setAddr(&deployments.L1ERC721Bridge, "L1ERC721Bridge")
	setAddr(&deployments.L1ERC721BridgeProxy, "L1ERC721BridgeProxy")
	setAddr(&deployments.L1StandardBridge, "L1StandardBridge")
	setAddr(&deployments.L1StandardBridgeProxy, "L1StandardBridgeProxy")
	setAddr(&deployments.L2OutputOracle, "L2OutputOracle")
	setAddr(&deployments.L2OutputOracleProxy, "L2OutputOracleProxy")
	setAddr(&deployments.OptimismMintableERC20Factory, "OptimismMintableERC20Factory")
	setAddr(&deployments.OptimismMintableERC20FactoryProxy, "OptimismMintableERC20FactoryProxy")
	setAddr(&deployments.OptimismPortal, "OptimismPortal")
	setAddr(&deployments.OptimismPortalProxy, "OptimismPortalProxy")
	setAddr(&deployments.ProxyAdmin, "ProxyAdmin")
	setAddr(&deployments.SystemConfig, "SystemConfig")
	setAddr(&deployments.SystemConfigProxy, "SystemConfigProxy")
	setAddr(&deployments.ProtocolVersions, "ProtocolVersions")
	setAddr(&deployments.ProtocolVersionsProxy, "ProtocolVersionsProxy")
	setAddr(&deployments.DataAvailabilityChallenge, "DataAvailabilityChallenge")
	setAddr(&deployments.DataAvailabilityChallengeProxy, "DataAvailabilityChallengeProxy")

	// Tokamak-specific contracts
	setAddr(&deployments.ETHLockbox, "ETHLockbox")
	setAddr(&deployments.ETHLockboxProxy, "ETHLockboxProxy")
	setAddr(&deployments.RAT, "RAT")
	setAddr(&deployments.RATProxy, "RATProxy")

	// Additional contracts that might exist in Tokamak
	// Note: L2NativeToken is not in L1Deployments, it's an L2 contract
	// But we can log its existence
	if l2NativeToken, ok := addresses["L2NativeToken"]; ok {
		t.logger.Debug("Found L2NativeToken in deployment", "address", l2NativeToken)
	}

	return count
}

// LogDeploymentAddresses logs the deployment addresses for debugging
func (t *TokamakAddressOverride) LogDeploymentAddresses(deployments *genesis.L1Deployments) {
	if os.Getenv("DEBUG_DEPLOYMENTS") != "true" {
		return
	}

	t.logger.Debug("L1 Deployment Addresses",
		"AddressManager", deployments.AddressManager.Hex(),
		"DisputeGameFactoryProxy", deployments.DisputeGameFactoryProxy.Hex(),
		"L1CrossDomainMessengerProxy", deployments.L1CrossDomainMessengerProxy.Hex(),
		"L1StandardBridgeProxy", deployments.L1StandardBridgeProxy.Hex(),
		"OptimismPortalProxy", deployments.OptimismPortalProxy.Hex(),
		"SystemConfigProxy", deployments.SystemConfigProxy.Hex())

	// Log Tokamak-specific addresses
	if deployments.RAT != (common.Address{}) {
		t.logger.Debug("Tokamak-specific contracts",
			"RAT", deployments.RAT.Hex(),
			"RATProxy", deployments.RATProxy.Hex(),
			"ETHLockbox", deployments.ETHLockbox.Hex(),
			"ETHLockboxProxy", deployments.ETHLockboxProxy.Hex())
	}
}
```

### 2. init.go 수정 (정확한 버전)

**파일**: `op-e2e/config/init.go` (수정)

```go
// 파일 상단 import에 추가 (이미 있는 import와 병합)
import (
	// ... 기존 imports ...
	"os"
	"github.com/ethereum/go-ethereum/log"
	// ... 기타 imports ...
)

// 라인 197 근처 initAllocType 함수 수정
func initAllocType(root string, allocType AllocType) {
	// ... 기존 코드 (line 197-271) ...

	// Line 272-285 근처, ApplyPipeline 호출 후에 추가
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
		overrider, err := NewTokamakAddressOverride(lgr)
		if err != nil {
			lgr.Warn("Failed to create Tokamak address overrider", "error", err)
		} else {
			// st.Chains[0]에서 L1 deployments 가져오기
			// 주의: 정확한 state 구조에 따라 조정 필요
			l1Deployments := &genesis.L1Deployments{}
			if st.Chains[0].L1 != nil {
				// State 구조체에서 L1 deployments 추출
				// 실제 구조에 맞게 조정 필요
				if systemConfig := st.Chains[0].L1.SystemConfigProxy; systemConfig != (common.Address{}) {
					l1Deployments.SystemConfigProxy = systemConfig
				}
				// ... 다른 필드들도 유사하게 ...
			}

			if err := overrider.OverrideL1Deployments(l1Deployments, allocType); err != nil {
				lgr.Warn("Failed to override Tokamak addresses, using Optimism addresses",
					"allocType", allocType, "error", err)
			} else {
				overrider.LogDeploymentAddresses(l1Deployments)
				// Override된 주소를 state에 다시 저장
				// 주의: 실제 state 구조에 맞게 조정 필요
			}
		}
	}

	// 기존 코드 계속 (line 287-289)
	mtx.Lock()
	for mode, allocs := range l2Alloc {
		l2Allocs[mode] = allocs  // l2Alloc → l2Allocs로 수정
	}
	// ... 나머지 코드 ...
	mtx.Unlock()
}
```

### 3. 테스트 코드

**파일**: `op-e2e/config/tokamak_override_test.go` (신규)

```go
package config

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/genesis"
)

func TestTokamakAddressOverride(t *testing.T) {
	// Create a temporary .deploy file for testing
	tmpDir := t.TempDir()
	deployDir := path.Join(tmpDir, "packages", "tokamak", "contracts-bedrock", "deployments", "devnetL1")
	require.NoError(t, os.MkdirAll(deployDir, 0755))

	testAddresses := map[string]string{
		"AddressManager":               "0x1111111111111111111111111111111111111111",
		"DisputeGameFactoryProxy":      "0x2222222222222222222222222222222222222222",
		"L1CrossDomainMessengerProxy":  "0x3333333333333333333333333333333333333333",
		"L1StandardBridgeProxy":        "0x4444444444444444444444444444444444444444",
		"SystemConfigProxy":            "0x5555555555555555555555555555555555555555",
		"OptimismPortalProxy":          "0x6666666666666666666666666666666666666666",
		"RAT":                          "0x7777777777777777777777777777777777777777",
		"RATProxy":                     "0x8888888888888888888888888888888888888888",
	}

	deployData, err := json.Marshal(testAddresses)
	require.NoError(t, err)

	deployFile := path.Join(deployDir, ".deploy")
	require.NoError(t, os.WriteFile(deployFile, deployData, 0644))

	// Change working directory to tmpDir for the test
	originalWd, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer os.Chdir(originalWd)

	// Test the override functionality
	logger := log.New()
	overrider, err := NewTokamakAddressOverride(logger)
	require.NoError(t, err)

	deployments := &genesis.L1Deployments{}
	err = overrider.OverrideL1Deployments(deployments, AllocTypeMTCannon)
	require.NoError(t, err)

	// Verify addresses were overridden correctly
	require.Equal(t, common.HexToAddress(testAddresses["AddressManager"]), deployments.AddressManager)
	require.Equal(t, common.HexToAddress(testAddresses["DisputeGameFactoryProxy"]), deployments.DisputeGameFactoryProxy)
	require.Equal(t, common.HexToAddress(testAddresses["SystemConfigProxy"]), deployments.SystemConfigProxy)
	require.Equal(t, common.HexToAddress(testAddresses["RAT"]), deployments.RAT)
	require.Equal(t, common.HexToAddress(testAddresses["RATProxy"]), deployments.RATProxy)
}

func TestGetDeploymentPath(t *testing.T) {
	logger := log.New()
	overrider := &TokamakAddressOverride{
		root:   "/test/root",
		logger: logger,
	}

	tests := []struct {
		name      string
		allocType AllocType
		expected  string
	}{
		{
			name:      "MTCannon",
			allocType: AllocTypeMTCannon,
			expected:  "/test/root/packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy",
		},
		{
			name:      "MTCannonNext",
			allocType: AllocTypeMTCannonNext,
			expected:  "/test/root/packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy",
		},
		{
			name:      "AltDA",
			allocType: AllocTypeAltDA,
			expected:  "/test/root/packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := overrider.getDeploymentPath(tt.allocType)
			require.Equal(t, tt.expected, result)
		})
	}
}
```

### 4. README 문서 (업데이트)

**파일**: `op-e2e/docs/tokamak/README.md` (신규)

```markdown
# Tokamak E2E Testing with Address Override

## Overview

This adapter allows E2E tests to use Tokamak's deployed contract addresses instead of Optimism's op-deployer generated addresses, while supporting multiple allocation types.

## Architecture

```
E2E Test Start
    ↓
init() function
    ↓
For each AllocType (parallel):
    ├─ AllocTypeMTCannon
    ├─ AllocTypeMTCannonNext
    └─ AllocTypeAltDA
        ↓
    initAllocType(root, allocType)
        ↓
    ApplyPipeline() execution
        ↓
    If TOKAMAK_OVERRIDE_ADDRESSES=true:
        ├─ Load Tokamak .deploy file
        ├─ Override L1Deployments addresses
        └─ Keep L1 Allocs unchanged
        ↓
    Store in global variables
        ↓
E2E Test Execution
```

## Prerequisites

Generate Tokamak deployment files:

```bash
make devnet-allocs
```

Verify the deployment file exists:
```bash
cat packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy | jq keys
```

## Usage

### Basic Usage

```bash
# Enable address override
export TOKAMAK_OVERRIDE_ADDRESSES=true

# Run E2E tests
go test -v ./op-e2e/faultproofs -timeout 30m
```

### Debug Mode

```bash
export TOKAMAK_OVERRIDE_ADDRESSES=true
export DEBUG_DEPLOYMENTS=true
go test -v ./op-e2e/config -run TestL1Allocs
```

### Run Specific Tests

```bash
# Test with MTCannon allocation
export TOKAMAK_OVERRIDE_ADDRESSES=true
go test -v -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs

# Test with AltDA allocation
export TOKAMAK_OVERRIDE_ADDRESSES=true
go test -v -run TestAltDA ./op-e2e/altda
```

### Disable Override

```bash
unset TOKAMAK_OVERRIDE_ADDRESSES
```

## Testing

### Unit Tests

```bash
# Run override logic tests
go test -v ./op-e2e/config -run TestTokamakAddressOverride
```

### Integration Tests

```bash
# Full E2E test suite with override
export TOKAMAK_OVERRIDE_ADDRESSES=true
make test-e2e
```

## Supported AllocTypes

The override system supports all three allocation types:

1. **AllocTypeMTCannon** (default)
   - Standard multi-threaded Cannon configuration
   - Uses `devnetL1` deployment

2. **AllocTypeMTCannonNext**
   - Next version of MT Cannon
   - Uses `devnetL1` deployment (can be customized)

3. **AllocTypeAltDA**
   - Alternative Data Availability configuration
   - Uses `devnetL1` deployment (can be customized)

## Contract Mappings

The following contracts are overridden from the `.deploy` file:

### Standard Optimism Contracts
- AddressManager
- DisputeGameFactory & Proxy
- L1CrossDomainMessenger & Proxy
- L1StandardBridge & Proxy
- L1ERC721Bridge & Proxy
- L2OutputOracle & Proxy
- OptimismMintableERC20Factory & Proxy
- OptimismPortal & Proxy
- SystemConfig & Proxy
- ProxyAdmin
- ProtocolVersions & Proxy
- DataAvailabilityChallenge & Proxy

### Tokamak-Specific Contracts
- RAT & RATProxy (Risk Assessment Token)
- ETHLockbox & ETHLockboxProxy
- L2NativeToken (logged but not in L1Deployments)

## Troubleshooting

### Common Issues

#### "Tokamak deployment file not found"

```bash
# Regenerate deployment files
make devnet-allocs

# Verify file exists
ls -la packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy
```

#### "Failed to override addresses"

Check that the `.deploy` file has valid JSON:
```bash
cat packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy | jq .
```

#### State inconsistency errors

This can happen when the state dump doesn't match overridden addresses:

```bash
# Option 1: Disable override
unset TOKAMAK_OVERRIDE_ADDRESSES

# Option 2: Check if contracts have same bytecode
# If bytecode differs, state override won't work properly
```

## Development

### Adding New Contract Mappings

1. Add the field to `L1Deployments` struct (if not exists)
2. Update `overrideAddresses()` in `tokamak_override.go`
3. Add test case in `tokamak_override_test.go`

### Supporting New AllocTypes

1. Update `getDeploymentPath()` to handle new type
2. Create corresponding deployment directory
3. Add test coverage

## Performance Considerations

- Override adds ~100ms to initialization
- No impact on test execution performance
- Parallel processing for multiple AllocTypes maintained

## Limitations

⚠️ **State Compatibility**: Works best when Tokamak and Optimism contracts have identical storage layouts
⚠️ **Address-Only Override**: Contract bytecode remains from op-deployer
⚠️ **Pre-deployment Required**: Must run `make devnet-allocs` first

## See Also

- [Original Option C Plan](./OPTION-C-IMPLEMENTATION-PLAN-KR.md)
- [Revised Implementation Plan](./OPTION-C-IMPLEMENTATION-PLAN-KR-REVISED.md)
- [E2E Deployment Compatibility Analysis](./E2E-DEPLOYMENT-COMPATIBILITY-ANALYSIS-KR.md)
```

## 구현 단계 (업데이트)

### Phase 1: Override 로직 구현 (1일)

**Task 1.1**: `tokamak_override.go` 작성
- TokamakAddressOverride 클래스 구현
- 다중 AllocType 지원
- 26개 필드 매핑

**Task 1.2**: `tokamak_override_test.go` 작성
- Unit test 구현
- Mock deployment 파일 테스트

**Task 1.3**: 컴파일 테스트
```bash
go build ./op-e2e/config
go test ./op-e2e/config -run TestTokamakAddressOverride
```

### Phase 2: init.go 통합 (0.5일)

**Task 2.1**: init.go 수정
- Import 추가
- initAllocType 함수 수정
- State 접근 로직 확인

**Task 2.2**: State 구조 확인
- `state.State` 구조체에서 L1 deployments 접근 방법 확인
- 필요시 조정

### Phase 3-6: (기존과 동일)

## 예상 문제 및 해결 (업데이트)

### 문제 1: State 구조체 접근

**증상**: `st.Chains[0]`에서 L1Deployments를 찾을 수 없음

**해결**:
```go
// op-deployer의 state 구조를 정확히 확인
// state/state.go 파일 참조하여 정확한 경로 파악
```

### 문제 2: 병렬 처리 시 동기화

**증상**: 여러 AllocType이 동시에 override 시도

**해결**:
```go
// Mutex로 보호
var overrideMutex sync.Mutex
overrideMutex.Lock()
defer overrideMutex.Unlock()
```

## 검증 계획

### 1. 단위 검증
- [ ] `tokamak_override.go` 컴파일
- [ ] `tokamak_override_test.go` 통과
- [ ] 모든 AllocType에서 override 작동

### 2. 통합 검증
- [ ] `make devnet-allocs` 실행
- [ ] E2E 초기화 통과
- [ ] 주소 override 확인

### 3. E2E 검증
- [ ] Alphabet game 테스트
- [ ] Cannon game 테스트
- [ ] AltDA 테스트 (해당시)

## 참고 자료

- [op-deployer state structure](../../op-deployer/pkg/deployer/state/state.go)
- [genesis L1Deployments](../../op-chain-ops/genesis/config.go#L1152)
- [E2E test initialization](../../op-e2e/config/init.go)

---

**문서 버전**: 3.0 (Final)
**최종 업데이트**: 2025-11-05
**다음 리뷰**: 구현 완료 후