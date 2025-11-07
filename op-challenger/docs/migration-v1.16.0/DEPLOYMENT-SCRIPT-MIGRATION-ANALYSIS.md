# Tokamak Thanos 배포 스크립트 마이그레이션 분석

## 개요

이 문서는 Optimism contracts-bedrock의 배포 스크립트 구조와 Tokamak Thanos의 배포 스크립트 구조를 비교 분석하고, Tokamak Thanos를 Optimism 스타일로 마이그레이션하기 위한 요구사항을 정리한 문서입니다.

**분석 일자**: 2025-11-06
**Optimism 경로**: `/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock`
**Tokamak Thanos 경로**: `/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock`

---

## 1. 전체 구조 비교

### 1.1 배포 스크립트 개수

| 프로젝트 | 배포 스크립트 수 | 위치 |
|---------|----------------|------|
| **Optimism** | 26개 | `scripts/deploy/*.s.sol` |
| **Tokamak Thanos** | 7개 | `scripts/deploy/*.s.sol` |

### 1.2 배포 스크립트 목록 비교

#### Optimism 보유 스크립트 (26개)
```
✓ Deploy.s.sol
✓ DeploySuperchain.s.sol
✓ DeployImplementations.s.sol
✓ DeployOPChain.s.sol
✓ DeployOwnership.s.sol
✓ DeployConfig.s.sol
✓ DeployAsterisc.s.sol
✓ DeployAuthSystem.s.sol
✓ DeployAuthSystem2.s.sol
✓ DeployAlphabetVM.s.sol
✓ DeployAltDA.s.sol
✓ DeployDisputeGame.s.sol
✓ DeployDisputeGame2.s.sol
✓ DeployMIPS.s.sol
✓ DeployMIPS2.s.sol
✓ DeployPreimageOracle.s.sol
✓ DeployPreimageOracle2.s.sol
✓ DeployProxy.s.sol
✓ DeployProxy2.s.sol
✓ SetDisputeGameImpl.s.sol
✓ UpgradeOPChain.s.sol
✓ InteropMigration.s.sol
✓ ReadImplementationAddresses.s.sol
✓ ReadSuperchainDeployment.s.sol
✓ GenerateOPCMMigrateCalldata.sol
✓ VerifyOPCM.s.sol
```

#### Tokamak Thanos 보유 스크립트 (7개)
```
✓ Deploy.s.sol
✓ DeploySuperchain.s.sol
✓ DeployImplementations.s.sol
✓ DeployOPChain.s.sol
✓ DeployOwnership.s.sol
✓ DeployConfig.s.sol
✓ DeployAsterisc.s.sol
```

#### 누락된 스크립트 (19개)
```
✗ DeployAuthSystem.s.sol
✗ DeployAuthSystem2.s.sol
✗ DeployAlphabetVM.s.sol
✗ DeployAltDA.s.sol
✗ DeployDisputeGame.s.sol
✗ DeployDisputeGame2.s.sol
✗ DeployMIPS.s.sol
✗ DeployMIPS2.s.sol
✗ DeployPreimageOracle.s.sol
✗ DeployPreimageOracle2.s.sol
✗ DeployProxy.s.sol
✗ DeployProxy2.s.sol
✗ SetDisputeGameImpl.s.sol
✗ UpgradeOPChain.s.sol
✗ InteropMigration.s.sol
✗ ReadImplementationAddresses.s.sol
✗ ReadSuperchainDeployment.s.sol
✗ GenerateOPCMMigrateCalldata.sol
✗ VerifyOPCM.s.sol
```

---

## 2. 주요 스크립트 상세 비교

### 2.1 DeploySuperchain.s.sol

#### 코드 차이점
| 구분 | Optimism | Tokamak Thanos |
|------|----------|----------------|
| **코드 라인 수** | 299줄 | 313줄 |
| **주요 차이점** | 깔끔한 구조 | console.log 디버깅 코드 추가 (60-92줄) |

#### Tokamak Thanos의 추가 코드
```solidity
// Line 60-92: 각 배포 단계마다 console.log 추가
console.log("=== DeploySuperchain.run START ===");
console.log("Step 1: Converted input");
console.log("Step 2: Validated input");
console.log("Step 3: Deployed SuperchainProxyAdmin");
console.log("Step 4: Deployed implementation contracts");
console.log("Step 5: Deployed and initialized SuperchainConfig");
console.log("Step 6: Deployed and initialized ProtocolVersions");
console.log("Step 7: Transferred ProxyAdmin ownership");
console.log("Step 8: Validated output");
console.log("=== DeploySuperchain.run SUCCESS ===");
```

**분석**:
- Tokamak Thanos는 디버깅을 위한 로깅이 추가됨
- 핵심 로직은 동일
- **권장사항**: Optimism 스타일로 통일하여 불필요한 로깅 제거

### 2.2 Deploy.s.sol

#### 구조 비교
| 구분 | Optimism | Tokamak Thanos |
|------|----------|----------------|
| **코드 라인 수** | 451줄 | 451줄 |
| **핵심 로직** | 완전 동일 | 완전 동일 |
| **상태** | ✅ 동일 | ✅ 동일 |

**분석**:
- 두 프로젝트의 Deploy.s.sol은 완전히 동일
- 마이그레이션 불필요

### 2.3 DeployOPChain.s.sol

#### 구조 비교
| 구분 | Optimism | Tokamak Thanos |
|------|----------|----------------|
| **코드 라인 수** | 1000+ 줄 | 1000+ 줄 |
| **첫 150줄** | 완전 동일 | 완전 동일 |
| **상태** | ✅ 동일 | ✅ 동일 |

**분석**:
- DeployOPChainInput 클래스 구조 동일
- set() 함수들 동일
- 마이그레이션 불필요

---

## 3. 라이브러리 구조 비교

### 3.1 scripts/libraries/ 디렉토리

#### Optimism 라이브러리
```
Chains.sol
Config.sol
Constants.sol
DeployUtils.sol
ForgeArtifacts.sol
Process.sol
Solarray.sol
StateDiff.sol
Types.sol
```

#### Tokamak Thanos 라이브러리
```
Chains.sol
Config.sol
Constants.sol
DeployUtils.sol
ForgeArtifacts.sol
LibStateDiff.sol      ← 추가된 파일
Process.sol
Solarray.sol
StateDiff.sol
Types.sol
```

**차이점**:
- Tokamak Thanos에는 `LibStateDiff.sol` 추가
- 나머지 파일은 동일

---

## 4. 보조 스크립트 및 유틸리티

### 4.1 Optimism의 추가 파일들

```
BaseDeployIO.sol           - 배포 I/O 베이스 클래스
ChainAssertions.sol        - 체인 상태 검증
Deployer.sol              - 배포자 베이스 클래스
StandardConstants.sol     - 표준 상수 정의
AddGameType.s.sol         - 게임 타입 추가
```

### 4.2 Tokamak Thanos 보유 여부

| 파일 | Tokamak Thanos |
|------|----------------|
| BaseDeployIO.sol | ✅ 있음 |
| ChainAssertions.sol | ✅ 있음 |
| Deployer.sol | ✅ 있음 |
| StandardConstants.sol | ❌ 없음 |
| AddGameType.s.sol | ❌ 없음 |

---

## 5. 주요 차이점 요약

### 5.1 누락된 기능

#### 1. Dispute Game 관련 배포 스크립트
- `DeployDisputeGame.s.sol`
- `DeployDisputeGame2.s.sol`
- `SetDisputeGameImpl.s.sol`
- `AddGameType.s.sol`

#### 2. Fault Proof VM 관련 스크립트
- `DeployAlphabetVM.s.sol`
- `DeployMIPS.s.sol`
- `DeployMIPS2.s.sol`
- `DeployPreimageOracle.s.sol`
- `DeployPreimageOracle2.s.sol`

#### 3. 권한 및 인증 시스템
- `DeployAuthSystem.s.sol`
- `DeployAuthSystem2.s.sol`

#### 4. Alternative DA (Data Availability)
- `DeployAltDA.s.sol`

#### 5. 프록시 배포 스크립트
- `DeployProxy.s.sol`
- `DeployProxy2.s.sol`

#### 6. 업그레이드 및 마이그레이션
- `UpgradeOPChain.s.sol`
- `InteropMigration.s.sol`

#### 7. 유틸리티 스크립트
- `ReadImplementationAddresses.s.sol`
- `ReadSuperchainDeployment.s.sol`
- `GenerateOPCMMigrateCalldata.sol`
- `VerifyOPCM.s.sol`

### 5.2 코드 스타일 차이

| 항목 | Optimism | Tokamak Thanos |
|------|----------|----------------|
| **로깅** | 최소화 | 상세한 디버그 로그 |
| **주석** | 표준화됨 | 표준화됨 |
| **구조** | 모듈화 | 모듈화 |

---

## 6. 아키텍처 차이 분석

### 6.1 배포 흐름 비교

#### Optimism 배포 흐름
```
Deploy.run()
  ├─ deploySuperchain()
  │   └─ DeploySuperchain.run()
  │       ├─ deploySuperchainProxyAdmin()
  │       ├─ deploySuperchainImplementationContracts()
  │       ├─ deployAndInitializeSuperchainConfig()
  │       ├─ deployAndInitializeProtocolVersions()
  │       └─ transferProxyAdminOwnership()
  │
  ├─ deployImplementations()
  │   └─ DeployImplementations.run()
  │       ├─ deployOPContractsManager()
  │       ├─ deploySharedContracts()
  │       └─ deployGameImplementations()
  │
  └─ deployOpChain()
      └─ OPContractsManager.deploy()
          ├─ deployProxies()
          ├─ initializeSystemConfig()
          ├─ initializeBridges()
          └─ initializeDisputeGames()
```

#### Tokamak Thanos 배포 흐름
```
Deploy.run()
  ├─ deploySuperchain()      ✅ 동일
  ├─ deployImplementations() ✅ 동일
  └─ deployOpChain()         ✅ 동일
```

**분석**: 핵심 배포 흐름은 동일하나, Optimism에는 더 많은 보조 스크립트 존재

### 6.2 모듈화 전략

#### Optimism의 모듈화
- **1단계 스크립트**: 핵심 배포 (Deploy, DeploySuperchain, DeployImplementations, DeployOPChain)
- **2단계 스크립트**: 특정 컴포넌트 배포 (DeployMIPS, DeployDisputeGame 등)
- **3단계 스크립트**: 유틸리티 (Read*, Verify*, Generate*)

#### Tokamak Thanos의 모듈화
- **1단계 스크립트만 보유**: 핵심 배포 스크립트만 존재
- **2/3단계 미보유**: 세부 컴포넌트 및 유틸리티 스크립트 부재

---

## 7. 마이그레이션 우선순위

### 7.1 즉시 필요 (High Priority)

#### 1. StandardConstants.sol
- **이유**: Deploy.s.sol에서 참조됨
- **위치**: `scripts/deploy/StandardConstants.sol`
- **내용**: MIPS_VERSION 등 표준 상수 정의

#### 2. Dispute Game 배포 스크립트
- **DeployDisputeGame.s.sol**: 분쟁 게임 구현 배포
- **SetDisputeGameImpl.s.sol**: 게임 구현 설정
- **이유**: Fault Proof 시스템의 핵심

#### 3. Fault Proof VM 스크립트
- **DeployMIPS2.s.sol**: MIPS VM 배포 (최신 버전)
- **DeployPreimageOracle2.s.sol**: Preimage Oracle 배포
- **이유**: Cannon Fault Proof 시스템 필수

### 7.2 중요도 중간 (Medium Priority)

#### 4. 프록시 배포 스크립트
- **DeployProxy2.s.sol**: ERC1967 Proxy 배포
- **이유**: 업그레이더블 컨트랙트 패턴 지원

#### 5. 유틸리티 스크립트
- **ReadImplementationAddresses.s.sol**: 구현 주소 읽기
- **ReadSuperchainDeployment.s.sol**: Superchain 배포 정보 읽기
- **이유**: 배포 후 검증 및 디버깅

#### 6. AuthSystem
- **DeployAuthSystem2.s.sol**: 인증 시스템 배포
- **이유**: 권한 관리 강화

### 7.3 선택적 (Low Priority)

#### 7. Alternative DA
- **DeployAltDA.s.sol**: 대체 DA 솔루션
- **이유**: 특정 사용 케이스에만 필요

#### 8. 마이그레이션 및 업그레이드
- **UpgradeOPChain.s.sol**: OPChain 업그레이드
- **InteropMigration.s.sol**: Interop 마이그레이션
- **이유**: 기존 체인 업그레이드 시에만 필요

#### 9. 기타 VM
- **DeployAlphabetVM.s.sol**: 알파벳 VM (테스트용)
- **이유**: 개발/테스트 환경에서만 사용

---

## 8. 호환성 분석

### 8.1 인터페이스 호환성

| 컴포넌트 | Optimism | Tokamak Thanos | 호환성 |
|---------|----------|----------------|--------|
| ISuperchainConfig | ✅ | ✅ | ✅ 호환 |
| IProtocolVersions | ✅ | ✅ | ✅ 호환 |
| IOPContractsManager | ✅ | ✅ | ✅ 호환 |
| IProxyAdmin | ✅ | ✅ | ✅ 호환 |
| DeployUtils | ✅ | ✅ | ✅ 호환 |

### 8.2 Config 호환성

Deploy Config 파일 구조는 양측이 동일:
```json
{
  "l1ChainID": ...,
  "l2ChainID": ...,
  "finalSystemOwner": ...,
  "superchainConfigGuardian": ...,
  ...
}
```

---

## 9. 리스크 분석

### 9.1 마이그레이션 리스크

#### High Risk
1. **Dispute Game 시스템 변경**
   - 분쟁 해결 메커니즘의 핵심
   - 잘못된 마이그레이션 시 보안 취약점 발생 가능

2. **Fault Proof VM 호환성**
   - MIPS 구현의 정확성 필수
   - 검증 실패 시 체인 무결성 문제

#### Medium Risk
3. **프록시 패턴 변경**
   - 업그레이드 로직 영향
   - 기존 배포된 컨트랙트와의 호환성

4. **권한 관리 시스템**
   - 잘못된 권한 설정 시 거버넌스 문제

#### Low Risk
5. **로깅 및 디버깅 코드**
   - 성능 영향 미미
   - 가독성 문제

### 9.2 버전 호환성 리스크

| 컴포넌트 | 리스크 레벨 | 이유 |
|---------|-----------|------|
| Core Contracts | Low | 동일한 베이스 사용 |
| Dispute Games | High | 새로운 스크립트 필요 |
| Fault Proof | High | VM 구현 정확성 필수 |
| Proxies | Medium | 업그레이드 패턴 변경 |

---

## 10. 테스트 전략

### 10.1 마이그레이션 후 테스트 항목

#### Unit Tests
- [ ] DeploySuperchain 단위 테스트
- [ ] DeployImplementations 단위 테스트
- [ ] DeployOPChain 단위 테스트
- [ ] 새로 추가된 스크립트 단위 테스트

#### Integration Tests
- [ ] 전체 배포 플로우 테스트
- [ ] Superchain → Implementations → OPChain 순차 배포
- [ ] 프록시 업그레이드 테스트
- [ ] Dispute Game 초기화 테스트

#### E2E Tests
- [ ] Devnet 배포 테스트
- [ ] Fault Proof 동작 확인
- [ ] 크로스 체인 메시지 전송
- [ ] Withdrawal 프로세스

### 10.2 회귀 테스트

기존 Tokamak Thanos 기능이 마이그레이션 후에도 동작하는지 확인:
- [ ] 기존 배포 스크립트 동작 확인
- [ ] Deploy Config 호환성
- [ ] 기존 테스트 통과 여부

---

## 11. 결론

### 11.1 주요 발견사항

1. **핵심 배포 스크립트는 동일**: Deploy.s.sol, DeployOPChain.s.sol 등은 거의 동일
2. **보조 스크립트 부재**: 19개의 보조 배포 스크립트가 Tokamak Thanos에 없음
3. **로깅 차이**: Tokamak Thanos는 디버그 로깅이 더 많음
4. **구조적 호환성**: 전반적인 아키텍처는 호환 가능

### 11.2 권장사항

#### 단기 (1-2주)
1. StandardConstants.sol 추가
2. 불필요한 console.log 제거 (Optimism 스타일로 통일)
3. DeploySuperchain.s.sol 정리

#### 중기 (1개월)
1. Dispute Game 배포 스크립트 추가
2. Fault Proof VM 스크립트 추가
3. 유틸리티 스크립트 추가

#### 장기 (2-3개월)
1. 전체 스크립트 세트 동기화
2. 자동화된 테스트 구축
3. CI/CD 파이프라인 통합

### 11.3 예상 작업량

| 작업 | 예상 시간 | 난이도 |
|------|----------|--------|
| StandardConstants 추가 | 1일 | 하 |
| 로깅 정리 | 2일 | 하 |
| Dispute Game 스크립트 | 1주 | 중 |
| Fault Proof VM 스크립트 | 1주 | 중 |
| 유틸리티 스크립트 | 3일 | 하 |
| 테스트 구축 | 1주 | 중 |
| **총합** | **3-4주** | **중** |

---

## 12. 다음 단계

1. ✅ 이 분석 문서 검토
2. ⏳ 마이그레이션 가이드 참조 (별도 문서)
3. ⏳ StandardConstants.sol부터 시작
4. ⏳ 단계별 마이그레이션 실행
5. ⏳ 각 단계마다 테스트 수행

---

**작성자**: Claude Code
**마지막 업데이트**: 2025-11-06
