# Tokamak Thanos 배포 스크립트 마이그레이션 가이드

## 개요

이 문서는 Tokamak Thanos의 배포 스크립트를 Optimism 스타일로 마이그레이션하기 위한 단계별 실행 가이드입니다.

**전제 조건**:
- [배포 스크립트 마이그레이션 분석](./DEPLOYMENT-SCRIPT-MIGRATION-ANALYSIS.md) 문서를 먼저 읽으세요.

---

## 목차

1. [사전 준비](#1-사전-준비)
2. [Phase 1: 코드 정리 (1-2일)](#phase-1-코드-정리-1-2일)
3. [Phase 2: 핵심 스크립트 추가 (1주)](#phase-2-핵심-스크립트-추가-1주)
4. [Phase 3: 보조 스크립트 추가 (1주)](#phase-3-보조-스크립트-추가-1주)
5. [Phase 4: 테스트 및 검증 (1주)](#phase-4-테스트-및-검증-1주)
6. [트러블슈팅](#트러블슈팅)

---

## 1. 사전 준비

### 1.1 환경 설정

```bash
# 작업 디렉토리 이동
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# Git 브랜치 생성
git checkout -b feature/migrate-deployment-scripts

# Optimism 소스 접근 확인
ls /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/
```

### 1.2 백업

```bash
# 현재 스크립트 백업
cp -r scripts/deploy scripts/deploy.backup-$(date +%Y%m%d)

# Git 커밋
git add .
git commit -m "backup: Save current deployment scripts before migration"
```

### 1.3 의존성 확인

```bash
# Foundry 버전 확인
forge --version

# 컴파일 테스트
forge build

# 기존 테스트 실행
forge test
```

---

## Phase 1: 코드 정리 (1-2일)

### Step 1.1: DeploySuperchain.s.sol 정리

**목표**: 불필요한 디버그 로깅 제거

#### 현재 상태 (Tokamak Thanos)
```solidity
function run(Input memory _input) public returns (Output memory output_) {
    console.log("=== DeploySuperchain.run START ===");

    InternalInput memory internalInput = toInternalInput(_input);
    console.log("Step 1: Converted input");

    assertValidInput(internalInput);
    console.log("Step 2: Validated input");

    deploySuperchainProxyAdmin(internalInput, output_);
    console.log("Step 3: Deployed SuperchainProxyAdmin");

    // ... more console.log statements
}
```

#### 변경 후 (Optimism 스타일)
```solidity
function run(Input memory _input) public returns (Output memory output_) {
    // Convert the external Input to InternalInput
    InternalInput memory internalInput = toInternalInput(_input);

    // Make sure the inputs are all set
    assertValidInput(internalInput);

    // Deploy the proxy admin, with the owner set to the deployer.
    deploySuperchainProxyAdmin(internalInput, output_);

    // Deploy and initialize the superchain contracts.
    deploySuperchainImplementationContracts(internalInput, output_);
    deployAndInitializeSuperchainConfig(internalInput, output_);
    deployAndInitializeProtocolVersions(internalInput, output_);

    // Transfer ownership of the ProxyAdmin from the deployer to the specified owner.
    transferProxyAdminOwnership(internalInput, output_);

    // Output assertions, to make sure outputs were assigned correctly.
    assertValidOutput(internalInput, output_);
}
```

#### 실행 명령

```bash
# 파일 수정
vim scripts/deploy/DeploySuperchain.s.sol

# 또는 자동화 스크립트 사용
cat > /tmp/cleanup-logs.sh << 'EOF'
#!/bin/bash
FILE="scripts/deploy/DeploySuperchain.s.sol"

# console.log 라인 제거 (60-92줄)
sed -i.bak '/console.log("=== DeploySuperchain.run START ===/d' "$FILE"
sed -i.bak '/console.log("Step [0-9]/d' "$FILE"
sed -i.bak '/console.log("=== DeploySuperchain.run SUCCESS ===/d' "$FILE"

echo "Cleaned up console.log statements"
EOF

chmod +x /tmp/cleanup-logs.sh
/tmp/cleanup-logs.sh
```

#### 검증

```bash
# 컴파일 확인
forge build

# 차이점 확인
git diff scripts/deploy/DeploySuperchain.s.sol
```

### Step 1.2: 커밋

```bash
git add scripts/deploy/DeploySuperchain.s.sol
git commit -m "refactor: Remove debug logs from DeploySuperchain to match Optimism style"
```

---

## Phase 2: 핵심 스크립트 추가 (1주)

### Step 2.1: StandardConstants.sol 추가

**우선순위**: ⭐⭐⭐ High

#### 파일 복사

```bash
# Optimism에서 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/StandardConstants.sol \
   scripts/deploy/StandardConstants.sol
```

#### 파일 내용

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

library StandardConstants {
    /// @notice The semantic version of the MIPS VM.
    uint8 public constant MIPS_VERSION = 1;
}
```

#### Deploy.s.sol 수정

기존:
```solidity
// Line 277
mipsVersion: 1,  // 하드코딩됨
```

변경:
```solidity
import { StandardConstants } from "scripts/deploy/StandardConstants.sol";

// Line 277
mipsVersion: StandardConstants.MIPS_VERSION,
```

#### 검증

```bash
forge build
forge test --match-contract Deploy
```

#### 커밋

```bash
git add scripts/deploy/StandardConstants.sol
git add scripts/deploy/Deploy.s.sol
git commit -m "feat: Add StandardConstants.sol for consistent version management"
```

### Step 2.2: DeployMIPS2.s.sol 추가

**우선순위**: ⭐⭐⭐ High (Fault Proof 핵심)

#### 파일 복사

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployMIPS2.s.sol \
   scripts/deploy/DeployMIPS2.s.sol
```

#### 의존성 확인

DeployMIPS2.s.sol이 의존하는 파일들:
- `interfaces/cannon/IMIPS64.sol` ✅ (이미 존재)
- `interfaces/cannon/IPreimageOracle.sol` ✅ (이미 존재)
- `src/cannon/MIPS64.sol` ✅ (이미 존재)

#### 컴파일 확인

```bash
forge build --force

# 특정 스크립트만 컴파일
forge build --contracts scripts/deploy/DeployMIPS2.s.sol
```

#### 테스트 작성

```solidity
// test/deploy/DeployMIPS2.t.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { Test } from "forge-std/Test.sol";
import { DeployMIPS2 } from "scripts/deploy/DeployMIPS2.s.sol";

contract DeployMIPS2Test is Test {
    DeployMIPS2 deployer;

    function setUp() public {
        deployer = new DeployMIPS2();
    }

    function test_deploy_succeeds() public {
        // TODO: 테스트 구현
    }
}
```

#### 커밋

```bash
git add scripts/deploy/DeployMIPS2.s.sol
git add test/deploy/DeployMIPS2.t.sol
git commit -m "feat: Add DeployMIPS2.s.sol for MIPS64 VM deployment"
```

### Step 2.3: DeployPreimageOracle2.s.sol 추가

**우선순위**: ⭐⭐⭐ High

#### 파일 복사

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployPreimageOracle2.s.sol \
   scripts/deploy/DeployPreimageOracle2.s.sol
```

#### 커밋

```bash
git add scripts/deploy/DeployPreimageOracle2.s.sol
git commit -m "feat: Add DeployPreimageOracle2.s.sol for preimage oracle deployment"
```

### Step 2.4: DeployDisputeGame.s.sol 추가

**우선순위**: ⭐⭐⭐ High

#### 파일 복사

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployDisputeGame.s.sol \
   scripts/deploy/DeployDisputeGame.s.sol
```

#### 의존성 확인

```bash
# 필요한 인터페이스 확인
grep -r "IFaultDisputeGame" scripts/deploy/DeployDisputeGame.s.sol
grep -r "IPermissionedDisputeGame" scripts/deploy/DeployDisputeGame.s.sol
```

#### 컴파일 및 검증

```bash
forge build
forge test --match-contract DisputeGame
```

#### 커밋

```bash
git add scripts/deploy/DeployDisputeGame.s.sol
git commit -m "feat: Add DeployDisputeGame.s.sol for dispute game deployment"
```

### Step 2.5: SetDisputeGameImpl.s.sol 추가

**우선순위**: ⭐⭐⭐ High

#### 파일 복사

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/SetDisputeGameImpl.s.sol \
   scripts/deploy/SetDisputeGameImpl.s.sol
```

#### 커밋

```bash
git add scripts/deploy/SetDisputeGameImpl.s.sol
git commit -m "feat: Add SetDisputeGameImpl.s.sol for setting dispute game implementations"
```

### Step 2.6: Phase 2 마무리

```bash
# 전체 빌드 확인
forge build

# 전체 테스트 실행
forge test

# Phase 2 완료 커밋
git add .
git commit -m "feat: Complete Phase 2 - Add core deployment scripts"
```

---

## Phase 3: 보조 스크립트 추가 (1주)

### Step 3.1: DeployProxy2.s.sol 추가

**우선순위**: ⭐⭐ Medium

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployProxy2.s.sol \
   scripts/deploy/DeployProxy2.s.sol

git add scripts/deploy/DeployProxy2.s.sol
git commit -m "feat: Add DeployProxy2.s.sol for proxy deployment"
```

### Step 3.2: 유틸리티 스크립트 추가

#### ReadImplementationAddresses.s.sol

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/ReadImplementationAddresses.s.sol \
   scripts/deploy/ReadImplementationAddresses.s.sol

git add scripts/deploy/ReadImplementationAddresses.s.sol
git commit -m "feat: Add ReadImplementationAddresses.s.sol utility"
```

#### ReadSuperchainDeployment.s.sol

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/ReadSuperchainDeployment.s.sol \
   scripts/deploy/ReadSuperchainDeployment.s.sol

git add scripts/deploy/ReadSuperchainDeployment.s.sol
git commit -m "feat: Add ReadSuperchainDeployment.s.sol utility"
```

### Step 3.3: DeployAuthSystem2.s.sol 추가

**우선순위**: ⭐⭐ Medium

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAuthSystem2.s.sol \
   scripts/deploy/DeployAuthSystem2.s.sol

git add scripts/deploy/DeployAuthSystem2.s.sol
git commit -m "feat: Add DeployAuthSystem2.s.sol for auth system"
```

### Step 3.4: 선택적 스크립트 추가

#### DeployAltDA.s.sol (이미 존재하는지 확인)

```bash
# Optimism 버전과 비교
diff /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAltDA.s.sol \
     scripts/deploy/DeployAltDA.s.sol

# 차이가 있다면 업데이트
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAltDA.s.sol \
   scripts/deploy/DeployAltDA.s.sol
```

#### DeployAlphabetVM.s.sol (테스트 전용)

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAlphabetVM.s.sol \
   scripts/deploy/DeployAlphabetVM.s.sol

git add scripts/deploy/DeployAlphabetVM.s.sol
git commit -m "feat: Add DeployAlphabetVM.s.sol for testing"
```

### Step 3.5: 업그레이드 스크립트 추가

#### UpgradeOPChain.s.sol

```bash
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/UpgradeOPChain.s.sol \
   scripts/deploy/UpgradeOPChain.s.sol

git add scripts/deploy/UpgradeOPChain.s.sol
git commit -m "feat: Add UpgradeOPChain.s.sol for chain upgrades"
```

### Step 3.6: Phase 3 마무리

```bash
# 전체 빌드
forge build

# 스크립트 개수 확인
find scripts/deploy -name "*.s.sol" | wc -l

# Phase 3 완료
git commit -am "feat: Complete Phase 3 - Add auxiliary scripts"
```

---

## Phase 4: 테스트 및 검증 (1주)

### Step 4.1: 단위 테스트 작성

#### 테스트 구조 생성

```bash
mkdir -p test/deploy

# 테스트 템플릿 생성
cat > test/deploy/DeployScripts.t.sol << 'EOF'
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { Test } from "forge-std/Test.sol";
import { Deploy } from "scripts/deploy/Deploy.s.sol";
import { DeploySuperchain } from "scripts/deploy/DeploySuperchain.s.sol";
import { DeployImplementations } from "scripts/deploy/DeployImplementations.s.sol";

contract DeployScriptsTest is Test {
    Deploy deploy;

    function setUp() public {
        deploy = new Deploy();
    }

    function test_deploySuperchain_succeeds() public {
        // TODO: 구현
    }

    function test_deployImplementations_succeeds() public {
        // TODO: 구현
    }
}
EOF
```

#### 테스트 실행

```bash
forge test --match-contract DeployScripts -vv
```

### Step 4.2: 통합 테스트

#### E2E 배포 테스트

```bash
# 테스트넷 설정
cat > .env.test << 'EOF'
L1_RPC_URL=http://localhost:8545
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
EOF

# E2E 스크립트 생성
cat > scripts/test-deployment.sh << 'EOF'
#!/bin/bash
set -e

echo "Starting E2E deployment test..."

# 1. Deploy Superchain
forge script scripts/deploy/DeploySuperchain.s.sol:DeploySuperchain \
    --rpc-url $L1_RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast

# 2. Deploy Implementations
forge script scripts/deploy/DeployImplementations.s.sol:DeployImplementations \
    --rpc-url $L1_RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast

# 3. Deploy OP Chain
forge script scripts/deploy/DeployOPChain.s.sol:DeployOPChain \
    --rpc-url $L1_RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast

echo "E2E deployment test completed successfully!"
EOF

chmod +x scripts/test-deployment.sh
```

#### 로컬 테스트넷에서 실행

```bash
# Anvil 시작
anvil &

# 배포 테스트 실행
./scripts/test-deployment.sh

# Anvil 종료
killall anvil
```

### Step 4.3: 회귀 테스트

#### 기존 기능 확인

```bash
# 기존 테스트가 여전히 통과하는지 확인
forge test

# 특정 컨트랙트 테스트
forge test --match-contract L1CrossDomainMessenger
forge test --match-contract OptimismPortal
forge test --match-contract SystemConfig
```

### Step 4.4: 스크립트 개수 검증

```bash
# 스크립트 개수 확인
echo "Optimism scripts:"
find /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy -name "*.s.sol" | wc -l

echo "Tokamak Thanos scripts:"
find scripts/deploy -name "*.s.sol" | wc -l

# 누락된 스크립트 확인
comm -23 \
  <(find /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy -name "*.s.sol" -exec basename {} \; | sort) \
  <(find scripts/deploy -name "*.s.sol" -exec basename {} \; | sort)
```

### Step 4.5: 문서 업데이트

```bash
# README 업데이트
cat >> README.md << 'EOF'

## Deployment Scripts

This repository contains deployment scripts aligned with Optimism's architecture:

### Core Scripts
- `Deploy.s.sol` - Main deployment orchestrator
- `DeploySuperchain.s.sol` - Superchain contracts
- `DeployImplementations.s.sol` - Implementation contracts
- `DeployOPChain.s.sol` - OP Chain specific contracts

### Component Scripts
- `DeployMIPS2.s.sol` - MIPS64 VM deployment
- `DeployPreimageOracle2.s.sol` - Preimage oracle
- `DeployDisputeGame.s.sol` - Dispute game system
- `SetDisputeGameImpl.s.sol` - Set dispute game implementations

### Utility Scripts
- `ReadImplementationAddresses.s.sol` - Read deployment addresses
- `ReadSuperchainDeployment.s.sol` - Read superchain config
- `UpgradeOPChain.s.sol` - Chain upgrades

For detailed migration information, see [docs/DEPLOYMENT-SCRIPT-MIGRATION-ANALYSIS.md](./docs/DEPLOYMENT-SCRIPT-MIGRATION-ANALYSIS.md)
EOF

git add README.md
git commit -m "docs: Update README with deployment scripts info"
```

### Step 4.6: Phase 4 완료

```bash
git add .
git commit -m "test: Complete Phase 4 - Testing and validation"
```

---

## 최종 체크리스트

### 코드 품질

- [ ] 모든 스크립트가 컴파일됨 (`forge build`)
- [ ] 불필요한 로그 제거됨
- [ ] Optimism 코딩 스타일 준수
- [ ] 주석이 적절히 작성됨

### 기능성

- [ ] StandardConstants.sol 추가됨
- [ ] 핵심 배포 스크립트 추가됨 (MIPS, DisputeGame 등)
- [ ] 유틸리티 스크립트 추가됨
- [ ] 업그레이드 스크립트 추가됨

### 테스트

- [ ] 단위 테스트 통과
- [ ] 통합 테스트 통과
- [ ] E2E 테스트 통과
- [ ] 회귀 테스트 통과

### 문서

- [ ] 마이그레이션 분석 문서 작성
- [ ] 마이그레이션 가이드 작성
- [ ] README 업데이트
- [ ] CHANGELOG 작성

### Git

- [ ] 의미 있는 커밋 메시지
- [ ] 브랜치 정리
- [ ] PR 준비 완료

---

## 트러블슈팅

### 문제 1: 컴파일 실패 - 인터페이스 누락

**증상**:
```
Error: Interface IMIPS64 not found
```

**해결**:
```bash
# 인터페이스 존재 확인
find . -name "IMIPS64.sol"

# Optimism에서 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/cannon/IMIPS64.sol \
   interfaces/cannon/IMIPS64.sol
```

### 문제 2: 테스트 실패 - Deployment 주소 불일치

**증상**:
```
Error: Address mismatch in deployment
```

**해결**:
```bash
# Artifacts 디렉토리 정리
rm -rf deployments/
rm -rf broadcast/

# 재배포
forge script scripts/deploy/Deploy.s.sol:Deploy --broadcast
```

### 문제 3: Gas 부족

**증상**:
```
Error: out of gas
```

**해결**:
```bash
# foundry.toml 수정
[profile.default]
gas_limit = "100000000"
gas_price = "1000000000"
```

### 문제 4: Nonce 불일치

**증상**:
```
Error: nonce too low
```

**해결**:
```bash
# Nonce 리셋 (테스트넷만)
anvil --reset

# 또는 강제로 nonce 지정
forge script ... --with-gas-price 0 --nonce 0
```

### 문제 5: 프록시 초기화 실패

**증상**:
```
Error: Proxy initialization failed
```

**해결**:
```solidity
// 초기화 체크
require(
    implementation != address(0),
    "Implementation not set"
);
```

---

## 참고 자료

### Optimism 문서
- [Optimism Deployment Docs](https://docs.optimism.io/)
- [OP Stack Specs](https://specs.optimism.io/)

### Foundry 문서
- [Foundry Book](https://book.getfoundry.sh/)
- [Forge Scripting](https://book.getfoundry.sh/tutorials/solidity-scripting)

### 관련 파일
- [마이그레이션 분석 문서](./DEPLOYMENT-SCRIPT-MIGRATION-ANALYSIS.md)
- [Optimism 소스 코드](https://github.com/ethereum-optimism/optimism)

---

## 다음 단계

마이그레이션 완료 후:

1. **PR 생성**
   ```bash
   git push origin feature/migrate-deployment-scripts
   # GitHub에서 PR 생성
   ```

2. **코드 리뷰 요청**
   - 팀원들에게 리뷰 요청
   - CI/CD 파이프라인 통과 확인

3. **배포 테스트**
   - Testnet에 배포
   - 모니터링 설정
   - 롤백 계획 수립

4. **메인넷 준비**
   - 보안 감사
   - 최종 테스트
   - 배포 체크리스트 완성

---

**작성자**: Claude Code
**마지막 업데이트**: 2025-11-06
**예상 완료 시간**: 3-4주
