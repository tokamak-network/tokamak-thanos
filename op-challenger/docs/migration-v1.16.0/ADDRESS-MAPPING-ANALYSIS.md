# 주소 매핑 분석 보고서

**작성일**: 2025-11-07
**목적**: L1Deployments 구조체에 Fault Proof VM 주소 추가가 필요한지 분석

---

## 🎯 핵심 질문

**"L1Deployments 구조체에 Mips, PreimageOracle 필드를 추가해야 하는가?"**

---

## 📊 분석 결과

### 1. addresses.json 현재 상태 ✅

```bash
$ cat .devnet/addresses.json | jq 'keys' | grep -E "(Mips|PreimageOracle|AnchorState|DelayedWETH)"
"AnchorStateRegistry"
"AnchorStateRegistryProxy"
"DelayedWETH"
"DelayedWETHProxy"
"Mips"
"PermissionedDelayedWETHProxy"
"PreimageOracle"
"Riscv"
```

**결과**: ✅ Fault Proof 관련 주소들이 이미 addresses.json에 모두 존재

---

### 2. L1Deployments 구조체 현재 상태

#### 기존 구조체 (수정 전)

```go
// op-chain-ops/genesis/config.go:1182
type L1Deployments struct {
    AddressManager                    common.Address `json:"AddressManager"`
    DisputeGameFactory                common.Address `json:"DisputeGameFactory"`
    DisputeGameFactoryProxy           common.Address `json:"DisputeGameFactoryProxy"`
    // ... 기타 필드들 ...
    RAT                               common.Address `json:"RAT"`
    RATProxy                          common.Address `json:"RATProxy"`
    // ❌ Mips, PreimageOracle, Riscv 필드 없음
}
```

**문제점**:
- JSON unmarshal 시 구조체에 없는 필드는 자동으로 무시됨
- addresses.json에 Mips, PreimageOracle이 있어도 로드되지 않음

---

### 3. 기존 타노스는 어떻게 작동했나?

#### E2E 테스트 초기화 흐름

```
op-e2e/config/init.go:init()
  ↓
initAllocType() - Runtime generation 방식
  ↓
deployer.ApplyPipelineE2E()
  ↓
L1 컨트랙트 배포 (MIPS, PreimageOracle 포함)
  ↓
addresses.L1Contracts → genesis.CreateL1DeploymentsFromContracts()
  ↓
L1Deployments 생성 (runtime에서)
```

**핵심 발견**:
- E2E 테스트는 **runtime generation** 방식 사용
- `.devnet/addresses.json`을 직접 로드하지 **않음**
- 대신 `deployer`가 생성한 `addresses.L1Contracts`를 `CreateL1DeploymentsFromContracts()`로 변환
- 이 과정에서 Mips, PreimageOracle은 `addresses.L1Contracts.MipsImpl`, `PreimageOracleImpl`에서 가져옴

#### addresses.L1Contracts 구조체

```go
// op-chain-ops/addresses/contracts.go:10
type L1Contracts struct {
    SuperchainContracts
    ImplementationsContracts  // ← 여기에 MipsImpl, PreimageOracleImpl 포함
    OpChainContracts
}

// op-chain-ops/addresses/contracts.go:29
type ImplementationsContracts struct {
    // ...
    PreimageOracleImpl               common.Address  // ✅ 있음
    MipsImpl                         common.Address  // ✅ 있음
    // ...
}
```

**결과**: ✅ Runtime generation 방식에서는 이미 Mips, PreimageOracle 주소가 전달됨

---

### 4. 문제가 발생하는 경우

#### initFromDevnetFiles() 사용 시

```go
// op-e2e/config/init.go:147
func initFromDevnetFiles(root string) error {
    // Load L1 deployments
    l1DeploymentsPath := path.Join(root, ".devnet", "addresses.json")
    l1Deployments, err := genesis.NewL1Deployments(l1DeploymentsPath)  // ← 문제!
    if err != nil {
        return fmt.Errorf("failed to load L1 deployments: %w", err)
    }
    // ...
}
```

**문제**:
- `genesis.NewL1Deployments()`는 JSON unmarshal 사용
- L1Deployments 구조체에 Mips, PreimageOracle 필드가 없음
- addresses.json의 Mips, PreimageOracle 값이 무시됨

**영향**:
- **현재는 영향 없음**: E2E 테스트가 `initFromDevnetFiles()`를 사용하지 않기 때문
- **잠재적 문제**: 향후 `.devnet` 파일 기반 초기화 사용 시 문제 발생 가능

---

### 5. CreateL1DeploymentsFromContracts() 분석

#### 현재 구현

```go
// op-chain-ops/genesis/config.go:1219
func CreateL1DeploymentsFromContracts(contracts *addresses.L1Contracts) *L1Deployments {
    return &L1Deployments{
        AddressManager:                    contracts.AddressManagerImpl,
        DisputeGameFactory:                contracts.DisputeGameFactoryImpl,
        // ...
        RAT:                               contracts.RATImpl,
        RATProxy:                          contracts.RATProxy,
        // ❌ Mips, PreimageOracle 매핑 없음
    }
}
```

**문제**:
- `contracts.MipsImpl`, `contracts.PreimageOracleImpl`이 존재하지만
- L1Deployments에 매핑되지 않음
- 따라서 L1Deployments를 통해 Mips, PreimageOracle 주소에 접근 불가

---

## 🔍 필요성 판단

### ❓ L1Deployments에 필드 추가가 필요한가?

#### Case 1: Runtime Generation만 사용 (현재 방식)
- **필요성**: ⚠️ **선택사항**
- **이유**: E2E 테스트가 L1Deployments를 통해 Mips/PreimageOracle에 접근하지 않음
- **현재**: 테스트가 `addresses.L1Contracts`를 직접 사용

#### Case 2: .devnet 파일 기반 초기화 사용 (향후)
- **필요성**: ✅ **필수**
- **이유**: `genesis.NewL1Deployments(addresses.json)` 사용 시 필드가 없으면 로드 불가
- **영향**: `.devnet` 파일 기반 테스트/도구가 Mips/PreimageOracle 주소를 찾을 수 없음

#### Case 3: L1Deployments 표준화 (Optimism 호환)
- **필요성**: ✅ **권장**
- **이유**: Optimism v1.16.0에서는 L1Deployments에 Fault Proof VM 주소 포함
- **장점**:
  - Optimism과 구조 일치
  - 향후 업그레이드 용이
  - 다른 도구와 호환성 향상

---

## 💡 결론

### 추천 방안: ✅ **필드 추가 권장**

**이유**:

1. **향후 호환성**: `.devnet` 파일 기반 초기화 지원을 위해 필요
2. **Optimism 표준**: v1.16.0 표준을 따르기 위해
3. **완전성**: addresses.json에 있는 모든 컨트랙트를 L1Deployments에서 접근 가능하게
4. **부작용 없음**: 기존 코드에 영향 없음 (JSON unmarshal은 추가 필드를 자동으로 처리)

### 추가할 필드

```go
type L1Deployments struct {
    // ... 기존 필드들 ...

    // Fault Proof VMs
    Mips                      common.Address `json:"Mips"`
    PreimageOracle            common.Address `json:"PreimageOracle"`
    Riscv                     common.Address `json:"Riscv"`
    AnchorStateRegistry       common.Address `json:"AnchorStateRegistry"`
    AnchorStateRegistryProxy  common.Address `json:"AnchorStateRegistryProxy"`
    DelayedWETH               common.Address `json:"DelayedWETH"`
    DelayedWETHProxy          common.Address `json:"DelayedWETHProxy"`
}
```

### CreateL1DeploymentsFromContracts() 업데이트

```go
func CreateL1DeploymentsFromContracts(contracts *addresses.L1Contracts) *L1Deployments {
    return &L1Deployments{
        // ... 기존 매핑들 ...

        // Fault Proof VMs
        Mips:                      contracts.MipsImpl,
        PreimageOracle:            contracts.PreimageOracleImpl,
        AnchorStateRegistry:       contracts.AnchorStateRegistryImpl,
        AnchorStateRegistryProxy:  contracts.AnchorStateRegistryProxy,
        DelayedWETH:               contracts.DelayedWethImpl,
        DelayedWETHProxy:          contracts.DelayedWethPermissionlessGameProxy,
    }
}
```

---

## 🚀 영향 분석

### ✅ 긍정적 영향

1. **addresses.json 완전 로드**: NewL1Deployments()가 모든 Fault Proof 주소 로드 가능
2. **Optimism 호환성**: v1.16.0 표준 구조와 일치
3. **향후 확장성**: .devnet 파일 기반 도구 개발 가능
4. **일관성**: 배포된 모든 L1 컨트랙트를 L1Deployments를 통해 접근 가능

### ⚠️ 잠재적 이슈

1. **없음**: 기존 코드는 runtime generation 방식 사용하므로 영향 없음
2. **테스트**: Go 컴파일 확인 필요 (구조체 호환성)

---

## 📋 실행 계획

### 1. 구조체 수정 ✅ (완료)
- L1Deployments에 Fault Proof VM 필드 추가

### 2. 매핑 함수 수정 (필요)
- CreateL1DeploymentsFromContracts() 업데이트

### 3. 검증 (필요)
```bash
# Go 빌드 확인
cd op-chain-ops/genesis
go build

# E2E 빌드 확인
cd ../../op-e2e
go build ./...

# 테스트 컴파일 확인
go test -c ./faultproofs
```

---

## 📚 참고

### 관련 커밋
- `27093d507` - "use copy of config.L1Deployments"
- `11843b48e` - "op-deployer: Add deploy mips script"
- `51b263230` - "OPCM upgrade: Add MTCannon to OPCM"

### 관련 파일
- `op-chain-ops/genesis/config.go` - L1Deployments 구조체
- `op-chain-ops/addresses/contracts.go` - L1Contracts 구조체
- `op-e2e/config/init.go` - E2E 초기화
- `.devnet/addresses.json` - 배포된 주소들

---

**작성자**: Claude Code
**마지막 업데이트**: 2025-11-07
**결론**: 필드 추가 권장 (향후 호환성 및 Optimism 표준 준수)
