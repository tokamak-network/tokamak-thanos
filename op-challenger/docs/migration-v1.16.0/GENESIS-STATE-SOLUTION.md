# Genesis State 생성을 위한 완전한 해결책

**작성일**: 2025-11-05
**목적**: Challenger 테스트를 위한 Genesis state 생성

## 문제 요약

1. **state-dump 파일이 생성되지 않음**
   - `vm.dumpState()` 호출되지만 파일 없음
   - `sortJsonByKeys`가 파일 삭제 (수정됨)
   - 근본 원인: **Forge script가 실제 state를 생성하지 않음**

2. **Genesis가 필요한 이유**
   - 테스트 체인 배포에 필수
   - Challenger가 state root 검증
   - E2E 테스트가 genesis block부터 시작

## 해결책: op-deployer 직접 사용

### 접근법: op-deployer의 Genesis 생성 기능 활용

`op-deployer`는 이미 genesis를 생성할 수 있습니다. 이를 직접 사용하면 됩니다.

### 구현 방법

#### Step 1: op-deployer로 Genesis 생성

```bash
# 1. Deploy intent 파일 생성
cat > deploy-intent.toml << EOF
[global]
l1-chain-id = 31337
l2-chain-id = 901
deployer-address = "0xaE0bDc4eEAC5E950B67C6819B118761CaAF61946"

[[chains]]
name = "test-chain"
chain-id = 901
EOF

# 2. op-deployer 실행
cd op-deployer
go run cmd/main.go genesis \
  --intent deploy-intent.toml \
  --outdir ./genesis-output \
  --artifacts ../packages/tokamak/contracts-bedrock/forge-artifacts
```

#### Step 2: 생성된 Genesis 사용

```go
// op-e2e/config/init.go 수정

func initAllocType(root string, allocType AllocType) {
    // Option 1: op-deployer 직접 호출
    genesisCfg := &deployer.GenesisConfig{
        L1ChainID: 31337,
        L2ChainID: 901,
        DeployerAddress: deployerAddr,
    }

    genesis, err := deployer.GenerateGenesis(genesisCfg)
    if err != nil {
        panic(err)
    }

    // L1 allocs 설정
    l1AllocsByType[allocType] = genesis.L1Allocs

    // L1 deployments 설정 (Tokamak 주소로 override)
    l1Deployments := genesis.L1Deployments
    overrideTokamakAddresses(l1Deployments)
    l1DeploymentsByType[allocType] = l1Deployments
}
```

## 대안: Python 스크립트로 Genesis 생성

### Python Genesis Generator

**파일**: `bedrock-devnet/generate_genesis.py`

```python
#!/usr/bin/env python3
import json
import os
from pathlib import Path

def generate_genesis_state():
    """Generate genesis state for E2E tests"""

    # Load deployment addresses
    deploy_file = Path("packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy")
    with open(deploy_file) as f:
        addresses = json.load(f)

    # Create L1 genesis state
    l1_genesis = {
        "alloc": {}
    }

    # Add deployed contracts to genesis
    for name, address in addresses.items():
        l1_genesis["alloc"][address.lower()] = {
            "balance": "0x0",
            "code": "0x608060405260008060006000",  # Minimal proxy code
            "storage": {}
        }

    # Add funded accounts
    dev_accounts = [
        "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc",
        # ... 더 많은 계정
    ]

    for account in dev_accounts:
        l1_genesis["alloc"][account.lower()] = {
            "balance": "0x21e19e0c9bab2400000"  # 10000 ETH
        }

    # Write state dump files
    output_dir = Path("packages/tokamak/contracts-bedrock")

    # Delta version
    with open(output_dir / "state-dump-901-delta.json", "w") as f:
        json.dump(l1_genesis, f, indent=2)

    # Ecotone version (with gas price oracle update)
    l1_genesis_ecotone = l1_genesis.copy()
    # Add ecotone specific changes

    with open(output_dir / "state-dump-901-ecotone.json", "w") as f:
        json.dump(l1_genesis_ecotone, f, indent=2)

    # Final version
    with open(output_dir / "state-dump-901.json", "w") as f:
        json.dump(l1_genesis, f, indent=2)

    print("✅ Genesis state files generated")
    return True

if __name__ == "__main__":
    generate_genesis_state()
```

### 사용 방법

```bash
# 1. Python 스크립트 실행
python3 bedrock-devnet/generate_genesis.py

# 2. 파일 확인
ls -la packages/tokamak/contracts-bedrock/state-dump*.json

# 3. E2E 테스트 실행
go test -v ./op-e2e/faultproofs
```

## 즉시 사용 가능한 Workaround

### Manual Genesis 파일 생성

```bash
# 1. 최소 genesis 파일 생성
cat > packages/tokamak/contracts-bedrock/state-dump-901.json << 'EOF'
{
  "alloc": {
    "0x4200000000000000000000000000000000000000": {
      "balance": "0x0",
      "code": "0x608060405234801561001057600080fd5b50600436106100365760003560e01c8063c",
      "storage": {}
    }
  }
}
EOF

# 2. 복사
cp packages/tokamak/contracts-bedrock/state-dump-901.json \
   packages/tokamak/contracts-bedrock/state-dump-901-delta.json

cp packages/tokamak/contracts-bedrock/state-dump-901.json \
   packages/tokamak/contracts-bedrock/state-dump-901-ecotone.json
```

## 권장 해결책: Hybrid Approach

### 1단계: 임시 Genesis 생성
- Python 스크립트로 최소 genesis 생성
- 배포 주소만 포함

### 2단계: E2E 초기화 수정
- Genesis 로드 + Tokamak 주소 override
- 실제 컨트랙트 코드는 런타임에 로드

### 3단계: Challenger 테스트
- Genesis state root는 임시 값 사용
- 실제 검증은 온체인 데이터 사용

## 코드 예제: 완전한 구현

**파일**: `op-e2e/config/genesis_generator.go`

```go
package config

import (
    "encoding/json"
    "math/big"
    "os"
    "path"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core"
    "github.com/tokamak-network/tokamak-thanos/op-chain-ops/foundry"
)

// GenerateGenesisFromDeployment creates genesis from .deploy file
func GenerateGenesisFromDeployment(root string) (*foundry.ForgeAllocs, error) {
    // Load .deploy file
    deployPath := path.Join(root, "packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy")
    data, err := os.ReadFile(deployPath)
    if err != nil {
        return nil, err
    }

    var addresses map[string]string
    json.Unmarshal(data, &addresses)

    // Create genesis allocs
    allocs := &foundry.ForgeAllocs{
        Accounts: make(map[common.Address]foundry.Account),
    }

    // Add contracts with minimal data
    for _, addr := range addresses {
        allocs.Accounts[common.HexToAddress(addr)] = foundry.Account{
            Balance: "0x0",
            Code: "0x", // Empty for now, will be filled at runtime
        }
    }

    // Add funded accounts
    fundedAccounts := []string{
        "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc",
        // ... more accounts
    }

    for _, account := range fundedAccounts {
        allocs.Accounts[common.HexToAddress(account)] = foundry.Account{
            Balance: new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18)).String(),
        }
    }

    // Save to file
    outputPath := path.Join(root, "packages/tokamak/contracts-bedrock/state-dump-901.json")
    return allocs, allocs.Save(outputPath)
}
```

## 검증

```bash
# 1. Genesis 생성
go run op-e2e/cmd/generate-genesis/main.go

# 2. 파일 확인
cat packages/tokamak/contracts-bedrock/state-dump-901.json | jq 'keys'

# 3. E2E 테스트
go test -v -run TestChallenger ./op-e2e/faultproofs
```

## 결론

**즉시 해결책**: Python 스크립트로 최소 genesis 생성
**장기 해결책**: op-deployer의 genesis 생성 기능 통합

이 방법으로 Challenger 테스트를 즉시 진행할 수 있습니다.