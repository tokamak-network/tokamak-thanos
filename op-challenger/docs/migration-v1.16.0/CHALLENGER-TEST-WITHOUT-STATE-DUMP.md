# Challenger 테스트를 위한 State-dump 없는 접근법

**작성일**: 2025-11-05
**결론**: State-dump 파일 없이도 Challenger 테스트 가능

## 핵심 발견

### TRH-SDK의 접근 방식
- **State-dump를 사용하지 않음**
- 배포된 컨트랙트 주소만 사용 (`.deploy` 파일)
- Terraform과 deployment JSON 사용

### 실제 필요한 것

| 항목 | 필요성 | 현재 상태 | 위치 |
|------|--------|----------|------|
| L1 컨트랙트 주소 | ✅ 필수 | ✅ 있음 | `deployments/devnetL1/.deploy` |
| Deploy Config | ✅ 필수 | ✅ 있음 | `deploy-config/devnetL1.json` |
| State-dump | ❌ 불필요 | ❌ 없음 | N/A |
| L2 Genesis | ✅ 필수 | 🔄 자동생성 | op-node가 생성 |

## 권장 구현 방안

### Option A: 배포 주소 직접 사용 (권장)

```go
// op-e2e/config/init.go 수정

func initAllocType(root string, allocType AllocType) {
    // 1. 기존 deployment 파일 로드
    deploymentPath := path.Join(root, "packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy")
    addresses := loadDeploymentAddresses(deploymentPath)

    // 2. L1Deployments 구조체 생성
    l1Deployments := createL1DeploymentsFromAddresses(addresses)

    // 3. Deploy config 로드
    deployConfig := loadDeployConfig(path.Join(root, "packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json"))

    // 4. L2 Genesis는 op-node가 생성
    // state-dump 없이도 작동

    // 5. 전역 변수 설정
    l1DeploymentsByType[allocType] = l1Deployments
    deployConfigsByType[allocType] = deployConfig
}
```

### Option B: Mock State 생성

```go
// 최소한의 state만 생성
func createMinimalL1State(deployments *L1Deployments) *foundry.ForgeAllocs {
    allocs := &foundry.ForgeAllocs{
        Accounts: make(map[common.Address]foundry.Account),
    }

    // 배포된 컨트랙트 주소에 빈 계정 추가
    allocs.Accounts[deployments.DisputeGameFactoryProxy] = foundry.Account{
        Balance: "0x0",
        Code: "0x", // 실제 코드는 체인에서 읽음
    }
    // ... 다른 컨트랙트들도 유사하게

    return allocs
}
```

## 구현 단계

### 1. Deployment 로더 구현

**파일**: `op-e2e/config/deployment_loader.go`

```go
package config

import (
    "encoding/json"
    "os"
    "path"
)

func LoadTokamakDeployment(root string) (*genesis.L1Deployments, error) {
    deployPath := path.Join(root, "packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy")

    data, err := os.ReadFile(deployPath)
    if err != nil {
        return nil, err
    }

    var addresses map[string]string
    if err := json.Unmarshal(data, &addresses); err != nil {
        return nil, err
    }

    return mapToL1Deployments(addresses), nil
}
```

### 2. E2E 테스트 수정

```go
// op-e2e/faultproofs/output_alphabet_test.go

func TestOutputAlphabetGame_ChallengerWins(t *testing.T) {
    // State-dump 없이 테스트
    cfg := DefaultSystemConfig(t)
    cfg.UseExistingDeployment = true // 기존 deployment 사용

    sys, err := cfg.Start(t)
    require.NoError(t, err)

    // Challenger 테스트 진행
    // ...
}
```

### 3. 환경 변수 설정

```bash
# State-dump 없이 실행
export USE_EXISTING_DEPLOYMENT=true
export SKIP_STATE_DUMP=true

# E2E 테스트 실행
go test -v ./op-e2e/faultproofs -timeout 30m
```

## 검증 방법

### 1. 기본 연결 테스트

```bash
# L1 컨트랙트 주소 확인
cat packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy | jq

# 컨트랙트 접근 테스트
cast call --rpc-url http://localhost:8545 \
  $(cat .deploy | jq -r .DisputeGameFactoryProxy) \
  "owner()"
```

### 2. Challenger 테스트

```bash
# DisputeGame 테스트
go test -v -run TestDisputeGame ./op-e2e/faultproofs

# Output proposal 테스트
go test -v -run TestOutputProposal ./op-e2e/faultproofs
```

## 장점

✅ **State-dump 생성 문제 우회**
✅ **TRH-SDK와 일치하는 접근법**
✅ **더 간단한 구현**
✅ **실제 배포 환경과 유사**

## 주의사항

⚠️ **초기 state 불일치 가능**
- 하지만 실행 중인 체인을 사용하므로 문제없음

⚠️ **일부 테스트 수정 필요**
- State root 검증 테스트는 조정 필요

## 결론

**State-dump 파일 없이도 Challenger 테스트 가능합니다.**

TRH-SDK도 state-dump를 사용하지 않으며, 실제로 필요한 것은:
1. 컨트랙트 주소 (`.deploy` 파일)
2. Deploy config
3. 실행 중인 L1/L2 체인

이 접근법으로 즉시 Challenger 테스트를 진행할 수 있습니다.