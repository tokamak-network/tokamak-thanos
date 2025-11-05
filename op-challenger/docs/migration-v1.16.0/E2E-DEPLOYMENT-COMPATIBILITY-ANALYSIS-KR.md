# E2E 테스트 배포 호환성 분석

**작성일**: 2025-11-04
**컨텍스트**: Tokamak-Thanos v1.16.0 마이그레이션
**이슈**: E2E 테스트 실행 시 "revision id 1 cannot be reverted" 에러 발생
**상태**: 🔴 심각 - Optimism v1.16.0 배포 스크립트로 E2E 테스트 실행 불가

## 요약

Optimism v1.16.0의 E2E 테스트는 테스트 환경을 생성하기 위해 배포 스크립트(`DeploySuperchain.s.sol`, `DeployImplementations.s.sol` 등)가 필요합니다. 그러나 이 스크립트들이 Tokamak-Thanos 컨트랙트에 대해 실행되면 컨트랙트 초기화 중에 `"revision id 1 cannot be reverted"` 에러가 발생하며 실패합니다.

**근본 원인**: Tokamak의 커스텀 컨트랙트 아키텍처는 다층 초기화 패턴(OpenZeppelin `Initializable` + `AccessControlUpgradeable` + 커스텀 로직)을 사용하는데, 이것이 Optimism의 프록시 업그레이드 플로우를 통해 실행될 때 EVM 상태 스냅샷 충돌을 일으킵니다.

**영향**: 이 비호환성이 해결되기 전까지 E2E 테스트(`TestOutputAlphabetGame_ReclaimBond` 및 기타)를 실행할 수 없습니다.

## 에러 분석

### 에러 내용

```
panic: revision id 1 cannot be reverted [recovered]
	panic: revision id 1 cannot be reverted

goroutine 53 [running]:
github.com/tokamak-network/tokamak-thanos/op-chain-ops/script.(*Host).Call.func1()
	/Users/zena/tokamak-projects/tokamak-thanos/op-chain-ops/script/script.go:362 +0x214
```

**위치**: `op-chain-ops/script/script.go:352-388`

### 에러의 의미

이 에러는 Go-Ethereum의 EVM 상태 관리에서 발생합니다:

1. **스냅샷(리비전)**: EVM은 트랜잭션 실행 전에 상태의 스냅샷을 생성하여 실패 시 롤백을 가능하게 합니다
2. **스냅샷 스택**: 스냅샷은 리비전 ID(0, 1, 2 등)를 가진 스택으로 관리됩니다
3. **문제점**: 코드가 리비전 ID 1로 되돌리려고 하지만, 해당 스냅샷이 이미 소비/무효화되었습니다

### 에러 핸들러 코드

**파일**: `op-chain-ops/script/script.go:352-388`

```go
func (h *Host) Call(from common.Address, to common.Address, input []byte, gas uint64, value *uint256.Int) (returnData []byte, leftOverGas uint64, err error) {
	h.prelude(from, &to)

	defer func() {
		if r := recover(); r != nil {
			rStr, ok := r.(string)
			if !ok || !strings.Contains(strings.ToLower(rStr), "revision id") {
				fmt.Printf("Unexpected panic in script execution: %v\n", r)
				panic(r)
			}

			// "revision id 1 cannot be reverted" 에러를 여기서 캐치
			fmt.Printf("Caught revision id error: %s\n", rStr)

			if h.evmRevertErr != nil {
				err = h.evmRevertErr
			} else {
				err = errors.New("execution reverted, check logs")
			}
		}
		h.evmRevertErr = nil
	}()

	returnData, leftOverGas, err = h.env.Call(from, to, input, gas, value)
	// ...
}
```

코드에 이 특정 에러를 위한 panic 복구 핸들러가 이미 있지만, 근본적인 문제로 인해 배포가 성공할 수 없습니다.

## Tokamak의 커스텀 컨트랙트 아키텍처

### 1. 커스텀 검증 시스템 (신규 - Optimism에 없음)

**파일**: `packages/tokamak/contracts-bedrock/src/tokamak-contracts/verification/L1ContractVerification.sol`

**주요 특징**:

```solidity
contract L1ContractVerification is
  IL1ContractVerification,
  Initializable,                    // OpenZeppelin 업그레이더블 패턴
  AccessControlUpgradeable          // 역할 기반 접근 제어
{
    function initialize(
        address _tokenAddress,
        address _initialAdmin
    ) public initializer {            // 다중 초기화 레이어
        __AccessControl_init();       // AccessControl 먼저 초기화
        _setupRole(DEFAULT_ADMIN_ROLE, _initialAdmin);
        _setupRole(ADMIN_ROLE, _initialAdmin);
        expectedNativeToken = _tokenAddress;
        isVerificationPossible = false;
    }

    /**
     * 중요 호환성 참고사항:
     * _proxyAdmin 파라미터는 반드시 OpenZeppelin v4.9.x 또는 이전 버전의
     * ProxyAdmin 컨트랙트여야 합니다. OpenZeppelin v5.x ProxyAdmin 컨트랙트는
     * getProxyImplementation을 지원하지 않아 이 함수가 revert됩니다.
     */
    function verifyL1Contracts(
        address _proxyAdmin,
        address _l1StandardBridgeProxy,
        // ...
    ) external onlyRole(ADMIN_ROLE) {
        // 커스텀 검증 로직
    }
}
```

**배포 스크립트에 대한 이슈**:
1. `Initializable` 패턴 필요 - `initialize()`를 정확히 한 번만 호출해야 함
2. `AccessControlUpgradeable` 필요 - 초기화 중 추가 상태 스냅샷 생성
3. ProxyAdmin 버전 제약 - OpenZeppelin v4.9.x 이하와만 호환
4. 검증 함수 호출 전에 역할 설정 필요

### 2. 커스텀 네이티브 토큰 (TON)

**파일**: `packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol:1091-1102`

```solidity
contract L2NativeToken is Ownable, ERC20Detailed, SeigToken {
    constructor() ERC20Detailed("Tokamak Network Token", "TON", 18) { }

    function setSeigManager(SeigManagerI) external pure override {
        revert("TON: TON doesn't allow setSeigManager");
    }

    function faucet(uint256 _amount) external {
        _mint(msg.sender, _amount);
    }
}
```

**배포 스크립트에 대한 이슈**:
1. 다중 상속 체인: `Ownable` → `ERC20Detailed` → `SeigToken` → `ERC20OnApprove`
2. 커스텀 초기화 순서 요구사항
3. 일부 함수가 의도적으로 비활성화됨(`setSeigManager`에서 revert)
4. Optimism의 표준 ETH 처리 방식과 다름

### 3. ProxyAdmin 호환성 요구사항

**표준 Optimism 패턴**:
```solidity
// packages/contracts-bedrock/src/universal/ProxyAdmin.sol:184-202
function upgradeAndCall(
    address payable _proxy,
    address _implementation,
    bytes memory _data
) external payable onlyOwner {
    ProxyType ptype = proxyType[_proxy];
    if (ptype == ProxyType.ERC1967) {
        Proxy(_proxy).upgradeToAndCall{ value: msg.value }(_implementation, _data);
    } else {
        upgrade(_proxy, _implementation);
        (bool success,) = _proxy.call{ value: msg.value }(_data);
        require(success, "ProxyAdmin: call to proxy after upgrade failed");
    }
}
```

**Tokamak의 추가 요구사항**:

**파일**: `packages/contracts-bedrock/scripts/Deploy.s.sol:1066-1073`

```solidity
// upgradeAndCall 전에 ProxyType을 반드시 설정해야 함
uint256 proxyType = uint256(proxyAdmin.proxyType(l1StandardBridgeProxy));
Safe safe = Safe(mustGetAddress("SystemOwnerSafe"));
if (proxyType != uint256(ProxyAdmin.ProxyType.CHUGSPLASH)) {
    _callViaSafe({
        _safe: safe,
        _target: address(proxyAdmin),
        _data: abi.encodeCall(ProxyAdmin.setProxyType,
              (l1StandardBridgeProxy, ProxyAdmin.ProxyType.CHUGSPLASH))
    });
}
```

**이슈**: Optimism의 배포 스크립트는 초기화 전에 프록시 타입을 설정하지 않습니다.

## 근본 원인: EVM 스냅샷 스택 불일치

### 정상적인 Optimism 플로우 (작동함)

```
1. ProxyAdmin.upgradeAndCall() 호출
2.   스냅샷 #0 생성
3.   Proxy.upgradeToAndCall() 호출
4.     스냅샷 #1 생성
5.     _implementation.delegatecall(_data) - 단순 초기화
6.       단일 단계 초기화 함수 실행
7.       성공적으로 반환
8.     스냅샷 #1로 되돌림 (성공) ✅
9.   스냅샷 #0으로 되돌림 (성공) ✅
```

### Tokamak 플로우 (실패함)

```
1. ProxyAdmin.upgradeAndCall() 호출
2.   스냅샷 #0 생성
3.   Proxy.upgradeToAndCall() 호출
4.     스냅샷 #1 생성
5.     _implementation.delegatecall(_data) - 다층 초기화
6.       initialize() 호출 (initializer modifier)
7.         스냅샷 #2 생성
8.         __AccessControl_init() 호출
9.           스냅샷 #3 생성
10.          _setupRole()이 더 많은 스냅샷 생성
11.          커스텀 로직이 더 많은 스냅샷 생성
12.        스냅샷 #1로 되돌리려고 시도
13.        ❌ 에러: 스냅샷 #1이 이미 소비/무효화됨
14.        panic: revision id 1 cannot be reverted
```

**문제점**: 다층 초기화(Initializable + AccessControlUpgradeable + 커스텀)가 프록시 업그레이드의 스냅샷 관리와 충돌하는 깊은 스냅샷 스택을 생성합니다.

## 아키텍처 차이점 요약

| 항목 | Optimism v1.16.0 | Tokamak-Thanos | 호환성 |
|------|------------------|----------------|--------|
| **검증 레이어** | 없음 | AccessControl을 사용한 커스텀 `L1ContractVerification.sol` | ❌ 비호환 |
| **네이티브 토큰** | 표준 ETH | SeigManager를 사용한 커스텀 `L2NativeToken` (TON) | ❌ 비호환 |
| **초기화 패턴** | 단순 단일 단계 | 다층 (Initializable + AccessControl + 커스텀) | ❌ 비호환 |
| **ProxyAdmin 버전** | 모든 버전 | OpenZeppelin v4.9.x 이하 필요 | ⚠️ 제약됨 |
| **프록시 타입 설정** | 자동 | 초기화 전에 수동 `setProxyType` 호출 필요 | ❌ 비호환 |
| **Safe 월렛 통합** | 선택적 | 3-of-3 멀티시그 검증과 함께 필수 | ❌ 다름 |
| **배포 플로우** | OPCM 기반 Solidity 스크립트 | Python 기반 (`bedrock-devnet/main.py`) | ❌ 근본적으로 다름 |

## 호환성을 위한 필요한 수정사항

Tokamak을 Optimism v1.16.0의 E2E 배포 스크립트와 호환되도록 하려면 다음 수정이 필요합니다:

### 옵션 A: Tokamak 컨트랙트 수정 (권장하지 않음)

**⚠️ 경고**: 이는 상당한 컨트랙트 변경이 필요하며 각 수정에 대해 사용자 동의가 필요합니다.

1. **L1ContractVerification 초기화 단순화**
   - 다층 초기화 패턴 제거
   - AccessControlUpgradeable 없이 단일 단계 초기화 사용
   - 예상 작업량: 3-5일 + 보안 감사

2. **L2NativeToken을 표준 패턴으로 적응**
   - 상속 체인 단순화
   - 단일 delegatecall 패턴과 호환되는 초기화
   - 예상 작업량: 2-3일 + 보안 감사

3. **ProxyType 요구사항 제거**
   - Optimism처럼 프록시 타입 감지를 자동화
   - 예상 작업량: 1-2일

**총 예상 작업량**: 2-3주 개발 + 보안 감사 + 테스트

**리스크**: 🔴 높음
- 프로덕션 컨트랙트에 대한 중대한 변경
- 보안 감사 필요
- Tokamak 특화 기능 손실 가능
- 현재 Python 시스템으로 이미 배포가 작동함

### 옵션 B: 배포 스크립트 수정 (부분적으로 가능)

Tokamak의 패턴을 수용하도록 Optimism의 배포 스크립트 커스터마이징:

1. **DeploySuperchain.s.sol 업데이트**
   - 초기화 호출 전에 프록시 타입 설정 추가
   - 초기화를 여러 단계로 분할
   - Safe 월렛 배포 및 소유권 이전 추가

   **파일**: `packages/tokamak/contracts-bedrock/scripts/deploy/DeploySuperchain.s.sol:83-98`

   ```solidity
   function deploySuperchainProxyAdmin(InternalInput memory, Output memory _output) private {
       vm.broadcast(msg.sender);
       IProxyAdmin superchainProxyAdmin = IProxyAdmin(
           DeployUtils.create1({
               _name: "ProxyAdmin",
               _args: DeployUtils.encodeConstructor(abi.encodeCall(IProxyAdmin.__constructor__, (msg.sender)))
           })
       );

       // TOKAMAK 수정: 초기화 전에 프록시 타입 설정
       vm.broadcast(msg.sender);
       superchainProxyAdmin.setProxyType(/* proxy address */, ProxyAdmin.ProxyType.ERC1967);

       vm.label(address(superchainProxyAdmin), "SuperchainProxyAdmin");
       _output.superchainProxyAdmin = superchainProxyAdmin;
   }
   ```

2. **DeployImplementations.s.sol 업데이트**
   - Tokamak 컨트랙트를 위한 커스텀 초기화 시퀀스 추가
   - 다층 초기화를 별도로 처리

3. **TokamakDeploymentHelpers.sol 추가**
   - Tokamak 특화 배포 로직을 위한 새 헬퍼 컨트랙트
   - AccessControl 역할 설정 처리
   - Safe 월렛 통합 관리

**예상 작업량**: 1-2주 개발 + 테스트

**리스크**: 🟡 중간
- 배포 스크립트의 포크를 유지해야 함
- 업스트림 Optimism과 분기됨
- Optimism이 배포 스크립트를 업데이트할 때마다 업데이트 필요

### 옵션 C: Tokamak E2E 테스트 어댑터 생성 (권장)

컨트랙트나 배포 스크립트를 호환되게 만드는 대신, Optimism E2E 테스트 요구사항을 Tokamak의 배포 시스템으로 변환하는 어댑터 레이어를 생성합니다.

**아키텍처**:

```
Optimism E2E Test
       ↓
TokamakE2EAdapter
       ↓
Python 배포 (bedrock-devnet/main.py)
       ↓
Tokamak 컨트랙트 (올바르게 배포 및 초기화됨)
       ↓
어댑터가 E2E 테스트에 주소 반환
```

**구현**:

1. **`op-e2e/tokamak/adapter.go` 생성**
   ```go
   package tokamak

   import (
       "github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils"
   )

   // TokamakDeploymentAdapter는 Tokamak의 Python 배포를 래핑
   type TokamakDeploymentAdapter struct {
       pythonDeployer *PythonDeploymentWrapper
   }

   // DeployL1Contracts는 Tokamak의 Python 시스템을 사용하여 배포
   func (a *TokamakDeploymentAdapter) DeployL1Contracts(cfg *e2eutils.DeployConfig) (*e2eutils.DeployResult, error) {
       // Python 배포 스크립트 호출
       result, err := a.pythonDeployer.Deploy(cfg)
       if err != nil {
           return nil, err
       }

       // Python 배포 출력을 E2E 테스트 형식으로 변환
       return &e2eutils.DeployResult{
           L1Deployments: translateAddresses(result),
           // ... 기타 필드
       }, nil
   }
   ```

2. **E2E 테스트 초기화 업데이트**

   **파일**: `op-e2e/config/init.go`

   ```go
   // Tokamak 특화 배포 경로 추가
   func InitL1(t *testing.T) (*L1Deployment, func()) {
       if useTokamakDeployment() {
           adapter := tokamak.NewDeploymentAdapter()
           return adapter.DeployL1Contracts(cfg)
       }

       // 표준 Optimism 배포
       return standardDeployL1(t)
   }
   ```

3. **Python 배포용 래퍼**

   ```go
   // PythonDeploymentWrapper는 bedrock-devnet/main.py를 실행
   func (w *PythonDeploymentWrapper) Deploy(cfg *Config) (*Result, error) {
       cmd := exec.Command("python3", "bedrock-devnet/main.py", "--config", cfg.ToJSON())
       output, err := cmd.CombinedOutput()
       if err != nil {
           return nil, err
       }
       return parseDeploymentOutput(output)
   }
   ```

**장점**:
- ✅ 컨트랙트 수정 불필요
- ✅ 배포 스크립트 수정 불필요
- ✅ Tokamak의 검증된 Python 배포 시스템 사용
- ✅ 프로덕션 배포를 손상시키지 않고 E2E 테스트 작동
- ✅ 명확한 관심사 분리

**예상 작업량**: 3-5일 개발 + 테스트

**리스크**: 🟢 낮음
- 기존 시스템에 대한 비침습적
- 테스트 및 검증 용이
- 보안 감사 불필요
- 독립적으로 업데이트 가능

### 옵션 D: 문제되는 E2E 테스트 스킵 (임시)

즉각적인 진행을 위해 전체 배포가 필요한 E2E 테스트를 스킵:

**파일**: `op-e2e/faultproofs/output_alphabet_bond_test.go`

```go
func TestOutputAlphabetGame_ReclaimBond(t *testing.T) {
    // 임시: 배포 어댑터가 구현될 때까지 스킵
    t.Skip("Skipping E2E test - Tokamak deployment scripts not yet compatible")

    // ... 나머지 테스트
}
```

환경 변수 제어 추가:

```go
func TestOutputAlphabetGame_ReclaimBond(t *testing.T) {
    if os.Getenv("TOKAMAK_SKIP_E2E_DEPLOYMENT") == "true" {
        t.Skip("Tokamak deployment E2E tests disabled")
    }
    // ... 나머지 테스트
}
```

**장점**:
- ✅ 다른 테스트 실행 가능
- ✅ 마이그레이션 진행 차단 해제
- ✅ 개발 노력 제로

**단점**:
- ❌ E2E 테스트가 실행되지 않음
- ❌ 임시 솔루션일 뿐

## 권장사항

**권장 접근법**: **옵션 C - Tokamak E2E 테스트 어댑터 생성**

**근거**:
1. **Tokamak의 아키텍처 보존**: 컨트랙트 변경 불필요
2. **프로덕션 배포 유지**: Python 배포 시스템 계속 작동
3. **E2E 테스트 활성화**: 적절한 Tokamak 배포를 사용하여 테스트 실행 가능
4. **낮은 리스크**: 어댑터가 프로덕션 시스템과 격리됨
5. **합리적인 작업량**: 다른 옵션의 2-3주 대비 3-5일

**단기**: **옵션 D**를 사용하여 옵션 C 개발 중 마이그레이션 차단 해제

**구현 계획**:

### 1단계: 즉시 (옵션 D)
1. 전체 배포가 필요한 E2E 테스트에 스킵 조건 추가
2. 어떤 테스트가 스킵되고 이유를 문서화
3. v1.16.0 마이그레이션의 나머지 계속

### 2단계: 어댑터 개발 (옵션 C) - 1주
1. `op-e2e/tokamak/adapter.go` 패키지 생성
2. Python 배포 래퍼 구현
3. 주소 변환 레이어 추가
4. 어댑터를 사용하도록 `op-e2e/config/init.go` 업데이트

### 3단계: 테스트 - 3-5일
1. 간단한 E2E 테스트로 어댑터 테스트
2. 주소 매핑 검증
3. 전체 E2E 테스트 스위트 실행
4. 남은 비호환성 문서화

### 4단계: 문서화 - 1일
1. 어댑터 아키텍처 문서화
2. 새 E2E 테스트 추가를 위한 개발자 가이드 추가
3. 마이그레이션 가이드 업데이트

**총 타임라인**: ~2주 (버퍼 포함)

## 대안: Tokamak 특화 E2E 테스트

Optimism의 E2E 테스트를 적응시키는 대신 Tokamak 특화 E2E 테스트 생성:

**위치**: `op-e2e/tokamak/`

**예시**: `op-e2e/tokamak/output_game_test.go`

```go
package tokamak_test

import (
    "testing"
    "github.com/tokamak-network/tokamak-thanos/op-e2e/tokamak"
)

func TestTokamakOutputAlphabetGame_ReclaimBond(t *testing.T) {
    // Tokamak의 배포 시스템을 직접 사용
    deployment := tokamak.SetupTokamakDevnet(t)

    // 게임 로직 테스트 실행
    // ... 테스트 구현
}
```

**장점**:
- Tokamak의 아키텍처를 위해 특별히 설계된 테스트
- Optimism의 변경되는 E2E 테스트와의 호환성 유지 불필요
- Tokamak 특화 기능 테스트 가능 (검증, TON 토큰 등)

**단점**:
- 유지보수할 테스트 코드 증가
- 업스트림 Optimism 변경으로 인한 회귀를 잡지 못함

## 주의가 필요한 파일

### 핵심 파일 (근본 원인)

1. **`packages/tokamak/contracts-bedrock/src/tokamak-contracts/verification/L1ContractVerification.sol`**
   - 라인 25-29: 다중 상속 초기화 패턴
   - 라인 113-128: 다층 `initialize()` 함수
   - 라인 166-169: ProxyAdmin 호환성 제약

2. **`packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol`**
   - 라인 1091-1102: 복잡한 상속을 가진 커스텀 토큰

3. **`op-chain-ops/script/script.go`**
   - 라인 352-388: revision ID panic을 캐치하는 에러 핸들러

### 배포 스크립트 (비호환성 소스)

4. **`packages/tokamak/contracts-bedrock/scripts/deploy/DeploySuperchain.s.sol`**
   - 라인 83-98: ProxyAdmin 배포 (프록시 타입 설정 필요)
   - 라인 100-138: 구현 배포 (커스텀 초기화 필요)

5. **`packages/tokamak/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`**
   - 전체 파일: 표준 Optimism 컨트랙트 초기화 예상

### E2E 테스트 파일 (어댑터 필요)

6. **`op-e2e/config/init.go`**
   - 라인 198: 아티팩트 경로 (이미 수정됨)
   - 라인 50-100: L1 배포 초기화 (어댑터 훅 필요)

7. **`op-e2e/faultproofs/output_alphabet_bond_test.go`**
   - 현재 실패하는 테스트
   - 어댑터 또는 스킵 조건 필요

### 프록시 관리

8. **`packages/contracts-bedrock/src/universal/ProxyAdmin.sol`**
   - 라인 184-202: `upgradeAndCall` 함수
   - 라인 113-124: `setProxyType` 함수

9. **`packages/contracts-bedrock/src/universal/Proxy.sol`**
   - 라인 67-81: `upgradeToAndCall` (스냅샷 생성)

## 테스트 전략

어댑터 구현 후:

1. **어댑터 단위 테스트**
   ```bash
   go test ./op-e2e/tokamak/... -v
   ```

2. **Python 배포와의 통합 테스트**
   ```bash
   go test -run TestTokamakDeploymentAdapter ./op-e2e/tokamak/... -v
   ```

3. **E2E 테스트 스위트**
   ```bash
   go test -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs/... -v -timeout 30m
   ```

4. **전체 테스트 스위트**
   ```bash
   make test-e2e
   ```

## 보안 고려사항

### 컨트랙트 수정 시 (옵션 A)

- ⚠️ **보안 감사 필요**: 모든 컨트랙트 초기화 변경은 전문 감사 필요
- ⚠️ **마이그레이션 리스크**: 기존 배포가 호환되지 않을 수 있음
- ⚠️ **기능 손실**: Tokamak 특화 기능이 제거될 수 있음

### 어댑터 사용 시 (옵션 C - 권장)

- ✅ **컨트랙트 변경 없음**: 기존 보안 보장 유지
- ✅ **격리된 리스크**: 어댑터는 테스트 환경에만 영향
- ✅ **쉬운 롤백**: 문제 발생 시 어댑터 비활성화 가능

## 모니터링 및 검증

솔루션 구현 후:

1. **E2E 테스트 성공률 검증**
   ```bash
   go test ./op-e2e/... -v | grep -E "(PASS|FAIL)"
   ```

2. **배포 아티팩트 검증**
   ```bash
   ls -la packages/tokamak/contracts-bedrock/forge-artifacts/Deploy*.s.sol/
   ```

3. **Python 배포가 여전히 작동하는지 확인**
   ```bash
   make devnet-allocs
   ```

4. **테스트 실행 시간 모니터링**
   - 어댑터가 테스트를 크게 느리게 하면 안 됨
   - 목표: 직접 배포 대비 < 5% 오버헤드

## 결론

"revision id 1 cannot be reverted" 에러는 Tokamak의 다층 컨트랙트 초기화 시스템과 Optimism의 더 단순한 배포 패턴 간의 근본적인 아키텍처 차이로 인해 발생합니다.

**권장 솔루션**: Optimism의 E2E 테스트를 Tokamak의 Python 배포 시스템과 연결하는 E2E 테스트 어댑터(옵션 C)를 생성합니다. 이 접근법은:

- ✅ Tokamak의 프로덕션 배포 시스템 보존
- ✅ 컨트랙트 수정 없이 E2E 테스트 활성화
- ✅ 낮은 리스크와 합리적인 개발 노력
- ✅ 기존 컨트랙트의 보안 유지

**단기**: v1.16.0 마이그레이션 차단을 해제하기 위해 문제되는 E2E 테스트 스킵(옵션 D)

**장기**: Optimism의 테스트를 적응시키는 대신 Tokamak의 아키텍처를 위해 설계된 Tokamak 특화 E2E 테스트 개발 고려

## 다음 단계

1. **즉시**: E2E 테스트에 스킵 조건 추가 (옵션 D)
2. **1주차**: E2E 테스트 어댑터 설계 및 구현 (옵션 C)
3. **2주차**: 전체 E2E 스위트로 어댑터 테스트
4. **3주차**: 솔루션 문서화 및 배포

## 참고자료

- [배포 스크립트 비교](./DEPLOYMENT-SCRIPTS-COMPARISON.md)
- [마이그레이션 가이드](./MIGRATION-GUIDE.md)
- [OpenZeppelin 프록시 패턴](https://docs.openzeppelin.com/contracts/4.x/api/proxy)
- [Go-Ethereum 상태 관리](https://github.com/ethereum/go-ethereum/blob/master/core/state/journal.go)

---

**문서 버전**: 1.0
**최종 업데이트**: 2025-11-04
**다음 리뷰**: 어댑터 구현 후
