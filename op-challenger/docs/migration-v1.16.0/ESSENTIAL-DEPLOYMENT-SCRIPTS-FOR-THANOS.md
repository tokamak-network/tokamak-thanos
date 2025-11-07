# Tokamak Thanos 필수 배포 스크립트 분석

## 개요

Optimism의 19개 누락 스크립트 중 **Tokamak Thanos에 반드시 필요한 스크립트**를 우선순위별로 분석한 문서입니다.

**작성 일자**: 2025-11-06
**분석 기준**: Fault Proof 시스템, Dispute Game, E2E 테스트 요구사항

---

## 📊 우선순위 요약

| 우선순위 | 개수 | 스크립트 |
|---------|------|---------|
| 🔴 **Critical (필수)** | 5개 | DeployMIPS2, DeployPreimageOracle2, DeployDisputeGame, SetDisputeGameImpl, StandardConstants |
| 🟡 **High (강력 권장)** | 4개 | DeployProxy2, ReadImplementationAddresses, DeployAuthSystem2, VerifyOPCM |
| 🟢 **Medium (권장)** | 3개 | ReadSuperchainDeployment, UpgradeOPChain, GenerateOPCMMigrateCalldata |
| ⚪ **Low (선택)** | 7개 | 나머지 스크립트 |

---

## 🔴 Critical 우선순위 (5개) - 즉시 필요

### 1. DeployMIPS2.s.sol

**왜 필수인가?**
- **Fault Proof 시스템의 핵심**: MIPS64 VM을 배포하는 스크립트
- **op-challenger가 의존**: Cannon fault proof를 실행하려면 MIPS VM이 필수
- **E2E 테스트 필수**: `op-e2e/faultproofs/` 테스트들이 MIPS를 사용

**기능**:
```solidity
// MIPS64 VM 싱글톤 배포
IMIPS64 mipsSingleton = IMIPS64(
    DeployUtils.createDeterministic({
        _name: "MIPS64",
        _args: abi.encodeCall(IMIPS64.__constructor__, (preimageOracle, mipsVersion)),
        _salt: DeployUtils.DEFAULT_SALT
    })
);
```

**영향 범위**:
- ✅ `DeployImplementations.s.sol`에서 호출됨
- ✅ `op-challenger`가 실행 시 MIPS VM 필요
- ✅ E2E 테스트의 MTCannon, MTCannonNext 타입

**타노스 현황**:
- ❌ 없음
- ⚠️ DeployImplementations에서 MIPS 관련 코드가 있지만 DeployMIPS2는 별도 스크립트

---

### 2. DeployPreimageOracle2.s.sol

**왜 필수인가?**
- **MIPS VM의 의존성**: MIPS는 PreimageOracle 없이 작동 불가
- **Large Preimage 지원**: v1.16.0에서 large preimage 기능 추가됨
- **op-challenger 필수 컴포넌트**

**기능**:
```solidity
// Preimage Oracle 배포
IPreimageOracle preimageOracle = IPreimageOracle(
    DeployUtils.create1({
        _name: "PreimageOracle",
        _args: abi.encodeCall(
            IPreimageOracle.__constructor__,
            (minProposalSize, challengePeriod)
        )
    })
);
```

**영향 범위**:
- ✅ DeployMIPS2의 Input으로 사용
- ✅ op-challenger/game/fault/trace 에서 사용
- ✅ Large preimage 테스트

**타노스 현황**:
- ❌ 없음
- ⚠️ PreimageOracle 컨트랙트는 있지만 배포 스크립트 없음

---

### 3. DeployDisputeGame.s.sol

**왜 필수인가?**
- **Dispute Game 구현 배포**: FaultDisputeGame, PermissionedDisputeGame 배포
- **E2E 테스트 필수**: 모든 dispute game 테스트가 이 스크립트에 의존
- **op-challenger 핵심**: Challenger가 게임과 상호작용

**기능**:
```solidity
// FaultDisputeGame 또는 PermissionedDisputeGame 배포
// gameKind 파라미터로 선택
if (gameKind == "FaultDisputeGame") {
    // FaultDisputeGame 배포
} else if (gameKind == "PermissionedDisputeGame") {
    // PermissionedDisputeGame 배포
}
```

**영향 범위**:
- ✅ DisputeGameFactory에 게임 구현 등록
- ✅ op-challenger가 게임 인스턴스 생성
- ✅ E2E 테스트의 모든 dispute game 시나리오

**타노스 현황**:
- ❌ 없음
- ⚠️ 현재는 OPContractsManager가 자동으로 배포하지만, 커스텀 게임 타입 배포 불가

---

### 4. SetDisputeGameImpl.s.sol

**왜 필수인가?**
- **게임 구현 등록**: DisputeGameFactory에 새로운 게임 타입 설정
- **Respected Game Type 설정**: AnchorStateRegistry에 존중할 게임 타입 지정
- **런타임 필수**: 배포 후 게임 타입을 활성화하는 유일한 방법

**기능**:
```solidity
// DisputeGameFactory에 게임 구현 설정
factory.setImplementation(gameType, impl);

// AnchorStateRegistry에 respected game type 설정
anchorStateRegistry.setRespectedGameType(gameType);
```

**영향 범위**:
- ✅ 새로운 게임 타입 추가 시 필수
- ✅ 게임 타입 업그레이드 시 필수
- ✅ E2E 테스트에서 permissioned/permissionless 게임 전환

**타노스 현황**:
- ❌ 없음
- ⚠️ Deploy.s.sol에서 하드코딩으로 설정하지만, 유연성 부족

---

### 5. StandardConstants.sol (라이브러리)

**왜 필수인가?**
- **버전 관리 일관성**: MIPS_VERSION을 중앙에서 관리
- **Deploy.s.sol에서 참조**: 현재 하드코딩된 값을 상수로 교체
- **업그레이드 용이성**: 버전 변경 시 한 곳만 수정

**내용**:
```solidity
library StandardConstants {
    /// @notice The semantic version of the MIPS VM.
    uint8 public constant MIPS_VERSION = 1;
}
```

**영향 범위**:
- ✅ DeployImplementations.s.sol
- ✅ DeployMIPS2.s.sol
- ✅ 모든 MIPS 관련 코드

**타노스 현황**:
- ❌ 없음
- ⚠️ Deploy.s.sol:277에 `mipsVersion: 1` 하드코딩됨

---

## 🟡 High 우선순위 (4개) - 강력 권장

### 6. DeployProxy2.s.sol

**왜 필요한가?**
- **ERC1967 Proxy 배포**: 표준 프록시 패턴 배포
- **유연한 프록시 관리**: create2를 사용한 결정론적 주소
- **업그레이더블 컨트랙트**: 향후 업그레이드 지원

**기능**:
```solidity
// ERC1967 Proxy 배포 (create2 사용)
IProxy proxy = IProxy(
    DeployUtils.createDeterministic({
        _name: "Proxy",
        _args: abi.encodeCall(IProxy.__constructor__, (admin)),
        _salt: salt
    })
);
```

**타노스 현황**:
- ❌ 없음
- 🟢 Deploy.s.sol에 `deployERC1967ProxyWithOwner()` 함수 있음 (부분적 대체)

---

### 7. ReadImplementationAddresses.s.sol

**왜 필요한가?**
- **배포 검증**: 배포된 구현 컨트랙트 주소 확인
- **디버깅 도구**: 배포 후 문제 진단
- **자동화 스크립트**: CI/CD에서 배포 확인

**기능**:
```solidity
// 모든 구현 컨트랙트 주소 읽기
address l1CrossDomainMessengerImpl = ...
address l1StandardBridgeImpl = ...
address optimismPortalImpl = ...
// ... 등등
```

**타노스 현황**:
- ❌ 없음
- 💡 수동으로 artifacts에서 읽어야 함

---

### 8. DeployAuthSystem2.s.sol

**왜 필요한가?**
- **권한 관리 강화**: 멀티시그, Safe 통합
- **보안 개선**: 단일 EOA 대신 Safe를 finalSystemOwner로 사용
- **프로덕션 필수**: 메인넷 배포 시 필요

**기능**:
```solidity
// Safe, ProxyAdmin 등 권한 시스템 배포
// DeployOwnership.s.sol의 향상된 버전
```

**타노스 현황**:
- ❌ DeployAuthSystem2 없음
- 🟡 DeployOwnership.s.sol은 있음 (기본 기능)

---

### 9. VerifyOPCM.s.sol

**왜 필요한가?**
- **배포 검증**: OPContractsManager 배포 상태 확인
- **무결성 체크**: 모든 구현체가 올바르게 설정되었는지 검증
- **감사 도구**: 보안 감사 시 사용

**타노스 현황**:
- ❌ 없음
- 💡 ChainAssertions.sol이 부분적으로 대체

---

## 🟢 Medium 우선순위 (3개) - 권장

### 10. ReadSuperchainDeployment.s.sol

**이유**: Superchain 배포 정보 읽기 (디버깅 도구)

### 11. UpgradeOPChain.s.sol

**이유**: 체인 업그레이드 스크립트 (운영 단계에서 필요)

### 12. GenerateOPCMMigrateCalldata.sol

**이유**: OPCM 마이그레이션 calldata 생성 (업그레이드 시)

---

## ⚪ Low 우선순위 (7개) - 선택적

### 13. DeployAuthSystem.s.sol
- **상태**: Deprecated (DeployAuthSystem2로 대체)
- **필요성**: ❌ 불필요

### 14. DeployAlphabetVM.s.sol
- **용도**: 테스트 전용 간단한 VM
- **필요성**: ⚪ 개발/테스트 환경에서만

### 15. DeployAltDA.s.sol
- **상태**: ✅ 타노스에 이미 있음
- **필요성**: 🟢 이미 보유

### 16-19. DeployMIPS.s.sol, DeployDisputeGame2.s.sol, DeployPreimageOracle.s.sol, DeployProxy.s.sol
- **상태**: Deprecated (2 버전으로 대체됨)
- **필요성**: ❌ 불필요

### 20. InteropMigration.s.sol
- **용도**: Interop 기능 마이그레이션
- **필요성**: ⚪ Interop 사용 시에만

---

## 📋 우선순위별 액션 플랜

### Phase 1: Critical (1주) - 즉시 착수

#### Week 1: Fault Proof 핵심
```bash
# 1일차: StandardConstants
cp optimism/scripts/deploy/StandardConstants.sol thanos/scripts/deploy/
# Deploy.s.sol 수정: mipsVersion: StandardConstants.MIPS_VERSION

# 2일차: DeployPreimageOracle2
cp optimism/scripts/deploy/DeployPreimageOracle2.s.sol thanos/scripts/deploy/
forge build

# 3일차: DeployMIPS2
cp optimism/scripts/deploy/DeployMIPS2.s.sol thanos/scripts/deploy/
forge build

# 4-5일차: DeployDisputeGame
cp optimism/scripts/deploy/DeployDisputeGame.s.sol thanos/scripts/deploy/
forge build
forge test --match-contract DisputeGame

# 6일차: SetDisputeGameImpl
cp optimism/scripts/deploy/SetDisputeGameImpl.s.sol thanos/scripts/deploy/
forge build
```

**검증**:
```bash
# E2E 테스트 실행
cd op-e2e
go test -v ./faultproofs/... -run TestOutputCannonProposedDisputeGame
```

### Phase 2: High Priority (3-5일) - 1주차 완료 후

```bash
# DeployProxy2
cp optimism/scripts/deploy/DeployProxy2.s.sol thanos/scripts/deploy/

# ReadImplementationAddresses
cp optimism/scripts/deploy/ReadImplementationAddresses.s.sol thanos/scripts/deploy/

# DeployAuthSystem2
cp optimism/scripts/deploy/DeployAuthSystem2.s.sol thanos/scripts/deploy/

# VerifyOPCM
cp optimism/scripts/deploy/VerifyOPCM.s.sol thanos/scripts/deploy/
```

### Phase 3: Medium Priority (1주) - 선택적

```bash
# ReadSuperchainDeployment
cp optimism/scripts/deploy/ReadSuperchainDeployment.s.sol thanos/scripts/deploy/

# UpgradeOPChain
cp optimism/scripts/deploy/UpgradeOPChain.s.sol thanos/scripts/deploy/

# GenerateOPCMMigrateCalldata
cp optimism/scripts/deploy/GenerateOPCMMigrateCalldata.sol thanos/scripts/deploy/
```

---

## 🎯 타노스 특화 요구사항

### op-challenger 통합 관점

Tokamak Thanos는 **op-challenger**를 사용하므로 다음 스크립트가 특히 중요:

| 스크립트 | op-challenger 연관성 | 중요도 |
|---------|---------------------|--------|
| DeployMIPS2 | ✅ Cannon VM 실행 | 🔴 Critical |
| DeployPreimageOracle2 | ✅ Large preimage 지원 | 🔴 Critical |
| DeployDisputeGame | ✅ 게임 생성 및 플레이 | 🔴 Critical |
| SetDisputeGameImpl | ✅ 게임 타입 설정 | 🔴 Critical |

### E2E 테스트 관점

E2E 테스트가 다음을 요구:

```go
// op-e2e/config/init.go
// MTCannon, MTCannonNext 타입 사용
require(mipsSingleton != nil)
require(preimageOracle != nil)
```

**필수 스크립트**:
1. ✅ DeployMIPS2
2. ✅ DeployPreimageOracle2
3. ✅ DeployDisputeGame

---

## 🔍 의존성 그래프

```
StandardConstants.sol (라이브러리)
    ↓
DeployPreimageOracle2.s.sol
    ↓
DeployMIPS2.s.sol
    ↓
DeployDisputeGame.s.sol
    ↓
SetDisputeGameImpl.s.sol
    ↓
[E2E 테스트 실행 가능]
```

**순서**:
1. StandardConstants.sol 추가
2. DeployPreimageOracle2.s.sol 추가
3. DeployMIPS2.s.sol 추가 (StandardConstants 참조)
4. DeployDisputeGame.s.sol 추가 (MIPS, Oracle 사용)
5. SetDisputeGameImpl.s.sol 추가 (Game 등록)

---

## 📊 기대 효과

### Critical 5개 추가 시

**Before**:
- ❌ Cannon fault proof 미작동
- ❌ E2E 테스트 실패 (MTCannon 타입)
- ❌ Large preimage 지원 불가
- ⚠️ 버전 관리 일관성 부족

**After**:
- ✅ Cannon fault proof 완전 지원
- ✅ E2E 테스트 통과 (모든 AllocType)
- ✅ Large preimage 지원
- ✅ 버전 관리 체계화
- ✅ Optimism 호환성 100%

### High 4개 추가 시

**추가 이점**:
- ✅ 배포 검증 자동화
- ✅ 프록시 관리 향상
- ✅ 권한 시스템 강화
- ✅ 디버깅 도구 확보

---

## ✅ 체크리스트

### 필수 (Critical) - 반드시 추가

- [ ] StandardConstants.sol
- [ ] DeployPreimageOracle2.s.sol
- [ ] DeployMIPS2.s.sol
- [ ] DeployDisputeGame.s.sol
- [ ] SetDisputeGameImpl.s.sol

### 강력 권장 (High) - 가능한 빨리 추가

- [ ] DeployProxy2.s.sol
- [ ] ReadImplementationAddresses.s.sol
- [ ] DeployAuthSystem2.s.sol
- [ ] VerifyOPCM.s.sol

### 권장 (Medium) - 운영 단계에서 추가

- [ ] ReadSuperchainDeployment.s.sol
- [ ] UpgradeOPChain.s.sol
- [ ] GenerateOPCMMigrateCalldata.sol

### 선택 (Low) - 필요 시 추가

- [ ] DeployAlphabetVM.s.sol (테스트용)
- [ ] InteropMigration.s.sol (Interop 사용 시)

---

## 🚀 빠른 시작 가이드

### 최소 필수 세트 (1일 작업)

```bash
# Optimism 경로
OPT_PATH=/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock
THANOS_PATH=/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# Critical 5개만 복사
cd $THANOS_PATH

# 1. StandardConstants (1분)
cp $OPT_PATH/scripts/deploy/StandardConstants.sol scripts/deploy/

# 2. DeployPreimageOracle2 (5분)
cp $OPT_PATH/scripts/deploy/DeployPreimageOracle2.s.sol scripts/deploy/

# 3. DeployMIPS2 (5분)
cp $OPT_PATH/scripts/deploy/DeployMIPS2.s.sol scripts/deploy/

# 4. DeployDisputeGame (30분 - 복잡함)
cp $OPT_PATH/scripts/deploy/DeployDisputeGame.s.sol scripts/deploy/

# 5. SetDisputeGameImpl (10분)
cp $OPT_PATH/scripts/deploy/SetDisputeGameImpl.s.sol scripts/deploy/

# 빌드 및 테스트
forge build
forge test

# E2E 테스트
cd ../../../op-e2e
go test -v ./faultproofs/... -run TestOutputCannonProposedDisputeGame
```

---

## 💡 결론

### 반드시 필요한 스크립트: **5개**

1. **StandardConstants.sol** - 버전 관리
2. **DeployPreimageOracle2.s.sol** - Preimage oracle
3. **DeployMIPS2.s.sol** - MIPS VM
4. **DeployDisputeGame.s.sol** - Dispute game 구현
5. **SetDisputeGameImpl.s.sol** - Game 등록

### 이유

- ✅ **op-challenger 작동 필수**: MIPS + Oracle + DisputeGame
- ✅ **E2E 테스트 필수**: MTCannon, MTCannonNext 타입
- ✅ **Fault Proof 시스템 완성**: 5개 스크립트가 전체 시스템 구성
- ✅ **Optimism 호환성**: v1.16.0 표준 준수

### 권장 추가: **4개**

6. DeployProxy2.s.sol
7. ReadImplementationAddresses.s.sol
8. DeployAuthSystem2.s.sol
9. VerifyOPCM.s.sol

### 나머지 10개: 선택적

- 테스트/개발 도구 또는 Deprecated

---

**작성자**: Claude Code
**참고 문서**:
- [DEPLOYMENT-SCRIPT-MIGRATION-ANALYSIS.md](./DEPLOYMENT-SCRIPT-MIGRATION-ANALYSIS.md)
- [DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md](./DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md)
- [E2E-TEST-ARCHITECTURE-KR.md](./E2E-TEST-ARCHITECTURE-KR.md)

**다음 단계**: [DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md](./DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md)의 Phase 1 실행
