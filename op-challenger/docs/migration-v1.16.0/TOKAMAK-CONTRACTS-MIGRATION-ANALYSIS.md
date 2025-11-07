# Tokamak Contracts 마이그레이션 분석

## 📌 개요

Optimism의 `packages/contracts-bedrock`와 Tokamak의 `packages/tokamak/contracts-bedrock`를 비교 분석하여 E2E 테스트에 필요한 마이그레이션 작업을 정리합니다.

**분석 일자**: 2025-11-06
**Optimism 버전**: v1.16.0
**Tokamak-Thanos 버전**: v1.16.0 migration

---

## 🔍 디렉토리 구조 비교

### Deploy Scripts (`scripts/deploy/`)

| 항목 | Optimism | Tokamak | 차이 |
|------|----------|---------|------|
| 파일 수 | 32개 | 10개 | **-22개** |
| 빌드 상태 | ✅ 정상 | ❌ 실패 | |

### Interfaces (`interfaces/`)

| 서브디렉토리 | Optimism | Tokamak | 상태 |
|-------------|----------|---------|------|
| `L1/` | ✅ 있음 | ❌ **없음** | **필수** |
| `L2/` | ✅ 있음 | ❌ **없음** | **필수** |
| `governance/` | ✅ 있음 | ❌ **없음** | 필수 |
| `safe/` | ✅ 있음 | ❌ **없음** | 필수 |
| `cannon/` | ✅ 있음 | ✅ 있음 | 동일 |
| `dispute/` | ✅ 있음 | ✅ 있음 | 일부 누락 |
| `legacy/` | ✅ 있음 | ✅ 있음 | 일부 누락 |
| `universal/` | ✅ 있음 | ✅ 있음 | 일부 누락 |

---

## 📋 Optimism에만 있는 Deploy 스크립트 (22개)

### 필수 스크립트 (op-deployer가 사용)

1. **DeployProxy.s.sol** ⭐ 최우선
   - 용도: Proxy 컨트랙트 배포
   - 사용처: 모든 proxy 배포 시

2. **DeployProxy2.s.sol** ⭐ 최우선
   - 용도: Proxy 컨트랙트 배포 (v2)

3. **DeployMIPS.s.sol** ⭐ 중요
   - 용도: MIPS VM 배포
   - 사용처: Cannon fault proof

4. **DeployMIPS2.s.sol** ⭐ 중요
   - 용도: MIPS VM 배포 (v2)

5. **DeployPreimageOracle.s.sol** ⭐ 중요
   - 용도: PreimageOracle 배포
   - 사용처: Cannon fault proof

6. **DeployPreimageOracle2.s.sol** ⭐ 중요
   - 용도: PreimageOracle 배포 (v2)

7. **DeployDisputeGame.s.sol** ⭐ 중요
   - 용도: DisputeGame 배포
   - 사용처: Fault proof system

8. **DeployDisputeGame2.s.sol** ⭐ 중요
   - 용도: DisputeGame 배포 (v2)

9. **DeployAlphabetVM.s.sol**
   - 용도: AlphabetVM 배포
   - 사용처: Alphabet game (테스트용)

10. **DeployAltDA.s.sol**
    - 용도: Alt-DA 배포
    - 사용처: Alt-DA 사용 시

### 관리/운영 스크립트

11. **AddGameType.s.sol**
    - 용도: DisputeGameFactory에 게임 타입 추가

12. **DeployAuthSystem.s.sol**
    - 용도: Auth system 배포

13. **DeployAuthSystem2.s.sol**
    - 용도: Auth system 배포 (v2)

14. **ReadImplementationAddresses.s.sol**
    - 용도: Implementation 주소 읽기

15. **ReadSuperchainDeployment.s.sol**
    - 용도: Superchain 배포 정보 읽기

16. **SetDisputeGameImpl.s.sol**
    - 용도: DisputeGame implementation 설정

17. **UpgradeOPChain.s.sol**
    - 용도: OPChain 업그레이드

18. **VerifyOPCM.s.sol**
    - 용도: OPCM 검증

19. **InteropMigration.s.sol**
    - 용도: Interop 마이그레이션

20. **GenerateOPCMMigrateCalldata.sol**
    - 용도: OPCM 마이그레이션 calldata 생성

21. **StandardConstants.sol**
    - 용도: 표준 상수 정의

22. **deploy.sh**
    - 용도: 배포 스크립트

---

## 📋 Optimism에만 있는 Interface 파일

### `interfaces/L1/` 디렉토리 전체 ⭐ 최우선

필수 파일:
- `ISuperchainConfig.sol` ← **DeploySuperchain.s.sol에서 필요**
- `IProtocolVersions.sol` ← **DeploySuperchain.s.sol에서 필요**
- `ISystemConfig.sol`
- `IOptimismPortal.sol`
- `IL1CrossDomainMessenger.sol`
- `IL1StandardBridge.sol`
- `IL1ERC721Bridge.sol`
- 등등...

### `interfaces/L2/` 디렉토리 전체

- `IL2ToL1MessagePasser.sol`
- `IL2CrossDomainMessenger.sol`
- `IL2StandardBridge.sol`
- 등등...

### `interfaces/governance/` 디렉토리

- Governance 관련 인터페이스

### `interfaces/safe/` 디렉토리

- Safe 관련 인터페이스

### `interfaces/universal/` 추가 파일

- `IProxyAdmin.sol` ← **DeploySuperchain.s.sol에서 필요**
- `ICrossDomainMessenger.sol`
- `IEIP712.sol`
- `IERC721Bridge.sol`
- `IOptimismMintableERC20.sol`
- `IOptimismMintableERC20Factory.sol`
- `IStandardBridge.sol`
- `IStaticERC1967Proxy.sol`
- `IWETH98.sol`

---

## 🚨 현재 빌드 에러

### `packages/tokamak/contracts-bedrock`에서 `forge build` 실행 시:

```
Error (6275): Source "interfaces/L1/ISuperchainConfig.sol" not found
Error (6275): Source "interfaces/L1/IProtocolVersions.sol" not found
Error (6275): Source "interfaces/universal/IProxyAdmin.sol" not found
Error (2904): Declaration "Proposal" not found in "src/dispute/lib/Types.sol"
```

**원인**:
1. ❌ `interfaces/L1/` 디렉토리 전체 누락
2. ❌ `interfaces/universal/IProxyAdmin.sol` 누락
3. ❌ `src/dispute/lib/Types.sol`에 `Proposal` 타입 없음

---

## 🎯 마이그레이션 필요 작업

### Phase 1: 필수 Interface 파일 복사 ⭐ 최우선

**E2E 테스트에 즉시 필요한 파일**:

```bash
# L1 interfaces
cp -r optimism/packages/contracts-bedrock/interfaces/L1 \
      tokamak-thanos/packages/tokamak/contracts-bedrock/interfaces/

# L2 interfaces
cp -r optimism/packages/contracts-bedrock/interfaces/L2 \
      tokamak-thanos/packages/tokamak/contracts-bedrock/interfaces/

# Universal interfaces 추가
cp optimism/packages/contracts-bedrock/interfaces/universal/IProxyAdmin.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/interfaces/universal/

cp optimism/packages/contracts-bedrock/interfaces/universal/ICrossDomainMessenger.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/interfaces/universal/

# (기타 필요한 파일들...)
```

### Phase 2: 필수 Deploy 스크립트 복사 ⭐ 중요

**op-deployer가 사용하는 스크립트**:

```bash
cd optimism/packages/contracts-bedrock/scripts/deploy

# 1. Proxy 배포
cp DeployProxy.s.sol DeployProxy2.s.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# 2. MIPS VM
cp DeployMIPS.s.sol DeployMIPS2.s.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# 3. PreimageOracle
cp DeployPreimageOracle.s.sol DeployPreimageOracle2.s.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# 4. DisputeGame
cp DeployDisputeGame.s.sol DeployDisputeGame2.s.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# 5. AlphabetVM (테스트용)
cp DeployAlphabetVM.s.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# 6. Alt-DA (optional)
cp DeployAltDA.s.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# 7. 관리 스크립트
cp ReadImplementationAddresses.s.sol ReadSuperchainDeployment.s.sol \
   tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/
```

### Phase 3: Types.sol 수정

**파일**: `packages/tokamak/contracts-bedrock/src/dispute/lib/Types.sol`

Optimism의 `Proposal` 타입 정의 추가:

```solidity
/// @notice A proposal in the Anchor State Registry
struct Proposal {
    Hash root;
    uint256 l2SequenceNumber;
}
```

### Phase 4: 빌드 및 검증

```bash
cd packages/tokamak/contracts-bedrock
forge build
# forge-artifacts 디렉토리 생성 확인
ls -la forge-artifacts/ | wc -l
```

---

## 📊 Tokamak 전용 파일 분석

### Tokamak에만 있는 파일

**src 디렉토리**:
- `src/L1/L2NativeToken.sol` - Native token 지원
- `src/L1/OnApprove.sol` - Token approval 핸들러
- `src/L2/ETH.sol` - Native ETH 처리
- `src/Safe/SafeExtender.sol` - Safe 확장
- `src/cannon/libraries/MIPSInstructions.sol` - MIPS 명령어
- `src/cannon/libraries/MIPSMemory.sol` - MIPS 메모리
- `src/cannon/libraries/MIPSState.sol` - MIPS 상태

**수정된 파일** (Optimism과 다름):
- `src/L1/L1CrossDomainMessenger.sol`
- `src/L1/L1StandardBridge.sol`
- `src/L1/OptimismPortal.sol`
- `src/L1/OptimismPortal2.sol`
- `src/L1/SystemConfig.sol`
- `src/L2/L2CrossDomainMessenger.sol`
- `src/L2/L2StandardBridge.sol`
- `src/L2/L2ToL1MessagePasser.sol`
- `src/Safe/DeputyGuardianModule.sol`
- `src/cannon/MIPS.sol`
- `src/cannon/PreimageOracle.sol`

**Tokamak 수정 사항**:
- ✅ Native Token 지원 추가
- ✅ ETH 처리 로직 변경
- ✅ MIPS 최적화

---

## 🎯 E2E 테스트 실행을 위한 최소 요구사항

### 즉시 필요한 파일 (빌드 에러 해결)

1. **interfaces/L1/** 디렉토리 전체
   - `ISuperchainConfig.sol`
   - `IProtocolVersions.sol`
   - 등등...

2. **interfaces/universal/IProxyAdmin.sol**

3. **scripts/deploy/DeployProxy.s.sol**
4. **scripts/deploy/DeployProxy2.s.sol**
5. **scripts/deploy/DeployMIPS.s.sol**
6. **scripts/deploy/DeployMIPS2.s.sol**
7. **scripts/deploy/DeployPreimageOracle.s.sol**
8. **scripts/deploy/DeployPreimageOracle2.s.sol**
9. **scripts/deploy/DeployDisputeGame.s.sol**
10. **scripts/deploy/DeployDisputeGame2.s.sol**

### 빌드 후 필요한 작업

```bash
# 1. Forge build
cd packages/tokamak/contracts-bedrock
forge build

# 2. Artifacts 검증
ls -la forge-artifacts/ | grep -i deploysuperchain
ls -la forge-artifacts/ | grep -i deployimplementations

# 3. E2E 테스트 실행
cd ../../op-e2e
go test -v ./faultproofs -run TestOutputCannonGame
```

---

## 🔄 마이그레이션 전략

### 전략 A: Optimism 파일 복사 + Tokamak 수정 병합 ⭐ 권장

**장점**:
- ✅ 빠른 해결
- ✅ Optimism과 호환성 유지
- ✅ E2E 테스트 즉시 실행 가능

**단점**:
- ⚠️ Tokamak 전용 기능 충돌 가능성
- ⚠️ 수동 병합 필요

**작업 순서**:
1. Optimism interfaces 복사
2. Optimism deploy 스크립트 복사
3. Tokamak 전용 수정사항 확인 및 병합
4. Forge build
5. E2E 테스트

### 전략 B: Tokamak 파일 수정 (점진적)

**장점**:
- ✅ Tokamak 수정사항 보존
- ✅ 깔끔한 마이그레이션

**단점**:
- ❌ 시간이 오래 걸림
- ❌ 복잡한 의존성 추적 필요

**작업 순서**:
1. 빌드 에러 하나씩 해결
2. 누락된 파일 추가
3. Import 경로 수정
4. 반복...

---

## 📝 현재 E2E 테스트 실행을 위한 블로커

### 블로커 1: forge-artifacts 부재

**문제**:
```
panic: invalid artifacts path:
stat packages/tokamak/contracts-bedrock/forge-artifacts: no such file or directory
```

**원인**: `forge build` 미실행

**해결**: Interfaces 파일 복사 후 `forge build`

### 블로커 2: Deploy 스크립트 artifacts 부재

**문제**:
```
failed to open artifact "DeployImplementations.s.sol/DeployImplementations.json"
```

**원인**: Deploy 스크립트가 빌드되지 않음

**해결**: 필수 deploy 스크립트 복사 후 빌드

### 블로커 3: Interface 파일 누락

**문제**:
```
Source "interfaces/L1/ISuperchainConfig.sol" not found
Source "interfaces/universal/IProxyAdmin.sol" not found
```

**원인**: Optimism v1.16.0의 interfaces 구조 변경

**해결**: Optimism interfaces 복사

---

## 🔧 즉시 실행 가능한 해결 방법

### 옵션 A: Optimism interfaces + deploy 스크립트 복사

```bash
#!/bin/bash
cd /Users/zena/tokamak-projects/tokamak-thanos

# 1. Interfaces 복사
cp -r /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/L1 \
      packages/tokamak/contracts-bedrock/interfaces/

cp -r /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/L2 \
      packages/tokamak/contracts-bedrock/interfaces/

cp -r /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/governance \
      packages/tokamak/contracts-bedrock/interfaces/

cp -r /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/safe \
      packages/tokamak/contracts-bedrock/interfaces/

# 2. Universal interfaces 업데이트 (덮어쓰기)
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/universal/*.sol \
   packages/tokamak/contracts-bedrock/interfaces/universal/

# 3. 필수 deploy 스크립트 복사
cd /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy
cp DeployProxy*.s.sol \
   DeployMIPS*.s.sol \
   DeployPreimageOracle*.s.sol \
   DeployDisputeGame*.s.sol \
   DeployAlphabetVM.s.sol \
   ReadImplementationAddresses.s.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# 4. Forge build
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
forge build

# 5. 검증
ls -la forge-artifacts/ | wc -l
```

### 옵션 B: packages/contracts-bedrock 사용 (임시)

```go
// op-e2e/config/init.go
artifactsPath := path.Join(root, "packages", "contracts-bedrock", "forge-artifacts")
```

**문제**: `packages/contracts-bedrock`에는 `scripts/deploy` 디렉토리가 없음!

---

## 📌 권장 사항

**즉시 실행**: **옵션 A** (파일 복사 + 빌드)

**이유**:
1. ✅ E2E 테스트를 빠르게 실행 가능
2. ✅ Optimism과 호환성 유지
3. ✅ Tokamak 전용 기능도 유지 가능

**단계**:
1. Interfaces 파일 복사 (L1, L2, governance, safe, universal)
2. Deploy 스크립트 복사 (Proxy, MIPS, PreimageOracle, DisputeGame)
3. Forge build
4. E2E 테스트 실행

---

## 🔍 추가 조사 필요 사항

### 1. Tokamak 수정 파일의 호환성

**수정된 src 파일들** (15개):
- Native Token 관련 수정이 E2E 테스트와 충돌하는가?
- Optimism interfaces와 호환되는가?

### 2. Deploy 스크립트의 Tokamak 커스터마이징

**현재 있는 10개 스크립트**:
- DeployConfig.s.sol
- Deploy.s.sol
- DeployImplementations.s.sol
- DeploySuperchain.s.sol
- DeployOPChain.s.sol
- 등등...

**확인 필요**:
- 이 파일들이 Optimism 버전과 어떻게 다른가?
- Tokamak 전용 수정이 있는가?
- 덮어쓰면 문제가 되는가?

### 3. foundry.toml 설정

**remappings 확인 필요**:
- interfaces 경로가 올바른가?
- lib 경로가 올바른가?

---

## 📌 결론

**E2E 테스트 실행을 위한 최소 요구사항**:

1. ✅ `interfaces/L1/` 전체 복사
2. ✅ `interfaces/L2/` 전체 복사
3. ✅ `interfaces/universal/IProxyAdmin.sol` 복사
4. ✅ 필수 deploy 스크립트 6개 복사
5. ✅ `forge build` 실행
6. ✅ E2E 테스트 재실행

**예상 소요 시간**: 10-15분 (복사 + 빌드)

**다음 단계**: 사용자 승인 후 파일 복사 진행

