# Dispute Game 마이그레이션 분석

**작성일**: 2025-11-07
**목적**: 타노스 Dispute Game 컨트랙트와 Optimism 최신 버전 비교 분석

---

## 📊 1. 주요 차이점

### A. Types.sol 구조체 이름 변경

**타노스 (현재)**:
```solidity
struct OutputRoot {
    Hash root;
    uint256 l2BlockNumber;
}
```

**Optimism (최신)**:
```solidity
struct Proposal {
    Hash root;
    uint256 l2SequenceNumber;
}
```

**변경 이유**:
- L2 블록 번호뿐만 아니라 타임스탬프 기반 시퀀스도 지원하기 위한 일반화
- 더 유연한 아키텍처를 위한 네이밍 변경

---

### B. 파일 구조 비교

#### 타노스에 있는 파일 (4개)
```
packages/tokamak/contracts-bedrock/src/dispute/
├── AnchorStateRegistry.sol
├── DisputeGameFactory.sol
├── FaultDisputeGame.sol
└── PermissionedDisputeGame.sol
```

#### Optimism에 있는 파일 (7개)
```
/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/dispute/
├── AnchorStateRegistry.sol
├── DelayedWETH.sol                    ⭐ 새로 추가
├── DisputeGameFactory.sol
├── FaultDisputeGame.sol
├── PermissionedDisputeGame.sol
├── SuperFaultDisputeGame.sol          ⭐ 새로 추가
└── SuperPermissionedDisputeGame.sol   ⭐ 새로 추가
```

#### 새로운 파일들
1. **DelayedWETH.sol**
   - 이전: `src/dispute/weth/DelayedWETH.sol`
   - 현재: `src/dispute/DelayedWETH.sol` (한 단계 위로 이동)
   - 용도: Bond 지연 지급 메커니즘

2. **SuperFaultDisputeGame.sol**
   - Super Root를 사용하는 새로운 게임 타입
   - 크로스체인 증명 지원

3. **SuperPermissionedDisputeGame.sol**
   - 권한 관리가 있는 Super 게임 타입

---

### C. lib 파일 차이

모든 lib 파일이 다름:
```bash
Files differ:
- lib/Errors.sol
- lib/LibUDT.sol
- lib/Types.sol
```

**주요 변경**:
- 더 많은 에러 타입 추가
- 새로운 GameType 추가 (SUPER_CANNON, SUPER_PERMISSIONED_CANNON, OP_SUCCINCT, KAILUA)
- BondDistributionMode enum 추가

---

## 🔍 2. 영향 받는 파일 분석

### A. Constants.sol

**현재 상태** (Optimism 버전으로 업데이트됨):
```solidity
import { Proposal, Hash } from "src/dispute/lib/Types.sol";

function DEFAULT_OUTPUT_ROOT() internal pure returns (Proposal memory) {
    return Proposal({ root: Hash.wrap(bytes32(hex"dead")), l2SequenceNumber: 0 });
}
```

**문제**: 타노스의 Types.sol은 `OutputRoot` 구조체를 사용하므로 충돌 발생

**해결**: Types.sol을 원본으로 revert 필요

---

### B. FaultDisputeGame.sol

**타노스**:
```solidity
OutputRoot public startingOutputRoot;
```

**Optimism**:
```solidity
Proposal public startingOutputRoot;
```

**영향**: Types.sol을 Optimism 버전으로 업데이트하면 FaultDisputeGame.sol도 함께 업데이트 필요

---

### C. 우리가 추가한 배포 스크립트들

다행히 영향 받지 않음:

| 스크립트 | 의존성 | 영향도 |
|---------|--------|--------|
| DeployMIPS2.s.sol | MIPS, PreimageOracle | ✅ 없음 |
| DeployPreimageOracle2.s.sol | PreimageOracle | ✅ 없음 |
| SetDisputeGameImpl.s.sol | 인터페이스만 사용 | ✅ 없음 |
| DeployAlphabetVM.s.sol | AlphabetVM | ✅ 없음 |
| DeployDisputeGame2.s.sol | 인터페이스만 사용 | ✅ 없음 |

**이유**: 이 스크립트들은 인터페이스를 통해서만 컨트랙트와 상호작용하므로 내부 구조 변경에 영향 받지 않음

---

## 🎯 3. 마이그레이션 전략

### 옵션 A: 전체 업데이트 ❌ (권장하지 않음)

**내용**:
- 모든 dispute 컨트랙트를 Optimism 최신 버전으로 교체
- lib 파일들도 모두 교체

**장점**:
- ✅ 최신 기능 사용 가능 (SuperFaultDisputeGame 등)
- ✅ Optimism과 완전히 동기화

**단점**:
- ❌ 타노스의 커스터마이징 손실 위험
- ❌ 대규모 변경으로 인한 불안정성
- ❌ 테스트 범위 확대 필요

---

### 옵션 B: 선택적 업데이트 ⚠️

**내용**:
1. ✅ interfaces만 업데이트 (이미 완료)
2. lib 파일들을 Optimism 버전으로 업데이트
3. 핵심 컨트랙트 (FaultDisputeGame 등)는 선택적 업데이트

**장점**:
- ✅ 필요한 부분만 업데이트
- ✅ 타노스 커스터마이징 부분 보존 가능

**단점**:
- ⚠️ 수동 병합 필요
- ⚠️ 호환성 테스트 필요

---

### 옵션 C: 최소 변경 ✅ (권장)

**내용**:
1. ✅ interfaces 디렉토리만 Optimism 최신으로 업데이트 (완료)
2. ✅ IRAT.sol, DeployAltDA.s.sol 추가 (완료)
3. ✅ Types.sol은 타노스 원본 유지 (OutputRoot 구조)
4. 컴파일 에러만 최소한으로 수정
5. E2E 테스트 실행

**장점**:
- ✅ 가장 안전한 방법
- ✅ 타노스 커스터마이징 완전 보존
- ✅ 변경 범위 최소화
- ✅ 빠른 E2E 테스트 진행 가능

**단점**:
- ❌ 최신 기능 (SuperFaultDisputeGame) 사용 불가
- (하지만 현재 목표는 E2E 테스트이므로 문제없음)

---

## 📋 4. 현재 상태

### 완료된 작업 ✅

1. **interfaces 디렉토리 업데이트**
   ```bash
   ✅ Optimism upstream에서 전체 interfaces/ 복사 완료
   ✅ packages/tokamak/contracts-bedrock/interfaces/L1/
   ✅ packages/tokamak/contracts-bedrock/interfaces/dispute/
   ✅ packages/tokamak/contracts-bedrock/interfaces/universal/
   ✅ packages/tokamak/contracts-bedrock/interfaces/cannon/
   ```

2. **누락된 파일 추가**
   ```bash
   ✅ IRAT.sol - Optimism에서 복사
   ✅ DeployAltDA.s.sol - Optimism에서 복사
   ```

3. **Types.sol 원복**
   ```bash
   ✅ git checkout packages/tokamak/contracts-bedrock/src/dispute/lib/Types.sol
   ✅ OutputRoot 구조체 유지
   ```

### 현재 파일 상태

| 파일 | 버전 | 상태 |
|------|------|------|
| interfaces/** | Optimism 최신 | ✅ 업데이트 완료 |
| src/dispute/lib/Types.sol | 타노스 원본 | ✅ 유지 |
| src/dispute/FaultDisputeGame.sol | 타노스 원본 | ✅ 유지 |
| src/dispute/DisputeGameFactory.sol | 타노스 원본 | ✅ 유지 |
| scripts/libraries/Constants.sol | 타노스 원본 | ✅ 유지 필요 |
| interfaces/L1/IRAT.sol | Optimism 최신 | ✅ 추가 완료 |
| scripts/deploy/DeployAltDA.s.sol | Optimism 최신 | ✅ 추가 완료 |

---

## 🎬 5. 권장 액션 플랜

### 단계 1: 현재 상태 확인 ✅
```bash
cd packages/tokamak/contracts-bedrock
forge build
```

**예상 결과**:
- ✅ 컴파일 성공 또는
- ⚠️ 최소한의 에러만 발생

### 단계 2: 컴파일 에러 수정
필요시 최소한의 수정만 진행

### 단계 3: forge-artifacts 생성 확인
```bash
ls packages/tokamak/contracts-bedrock/forge-artifacts
```

### 단계 4: E2E 테스트 실행
```bash
cd op-e2e
go test -v ./faultproofs -run TestOutputCannonGame -timeout 30m
```

### 단계 5: 테스트 결과 분석
- ✅ 성공: 현재 상태 유지
- ❌ 실패: 필요한 부분만 선택적 업데이트 (옵션 B)

---

## 📝 6. 결론 및 권고사항

### 현재 전략: 옵션 C (최소 변경) ✅

**이유**:
1. **목표 명확성**: E2E 테스트 실행이 목표 (전체 마이그레이션 아님)
2. **안정성 우선**: 타노스의 검증된 코드 유지
3. **위험 최소화**: 변경 범위를 최소화하여 예상치 못한 버그 방지
4. **점진적 접근**: E2E 테스트 후 필요시 점진적으로 업데이트 가능

### 향후 고려사항

**E2E 테스트 성공 후**:
- SuperFaultDisputeGame 등 신기능이 필요한 경우에만 선택적 업데이트 검토
- 타노스 팀과 협의하여 커스터마이징 요구사항 확인
- 단계적 마이그레이션 계획 수립

**E2E 테스트 실패 시**:
- 에러 로그 분석
- 필요한 최소 파일만 Optimism 버전으로 업데이트
- 호환성 테스트 재실행

---

## 🔗 7. 관련 문서

| 문서 | 내용 |
|------|------|
| `ROOT-CAUSE-ANALYSIS-KR.md` | E2E 테스트 실패 근본 원인 분석 |
| `DEPLOYMENT-SCRIPT-MIGRATION-ANALYSIS.md` | 배포 스크립트 마이그레이션 분석 |
| `PIPELINE-INTEGRATION-COMPLETE.md` | 파이프라인 통합 완료 보고서 |
| `ADDRESS-MAPPING-COMPLETE.md` | L1Deployments 주소 매핑 완료 |

---

**작성자**: Claude Code
**최종 업데이트**: 2025-11-07
**상태**: ✅ 분석 완료, 권장 전략 수립 완료
