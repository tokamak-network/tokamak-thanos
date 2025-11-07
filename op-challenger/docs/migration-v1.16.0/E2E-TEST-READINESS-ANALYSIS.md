# E2E 테스트 준비 상태 분석

**작성일**: 2025-11-07
**목적**: 배포 스크립트 추가 후 E2E 테스트 준비 상태 점검

---

## 🎯 핵심 질문

**"배포 스크립트만 수정하면 E2E 테스트 준비가 완료되었을까?"**

### 📊 결론

**❌ 아니오. 배포 스크립트 외에 추가 작업이 필요합니다.**

---

## ✅ 완료된 작업 (배포 스크립트)

### 1. Critical 배포 스크립트 (100% 완료)

| 스크립트 | 상태 | 역할 |
|---------|------|------|
| StandardConstants.sol | ✅ 완료 | MIPS 버전 관리 |
| DeployPreimageOracle2.s.sol | ✅ 완료 | PreimageOracle 배포 |
| DeployMIPS2.s.sol | ✅ 완료 | MIPS64 VM 배포 |
| DeployDisputeGame.s.sol | ✅ 완료 | DisputeGame 배포 |
| SetDisputeGameImpl.s.sol | ✅ 완료 | Factory 등록 |
| DeployAlphabetVM.s.sol | ✅ 완료 | 테스트용 VM |

**결과**:
- ✅ 컨트랙트 배포 스크립트 완비
- ✅ Solidity 레벨에서 준비 완료
- ✅ `forge build` 성공

---

## ❌ 미완료 작업 (E2E 테스트 통합)

### 1. Go 코드 레벨 작업

E2E 테스트는 **Go 코드**로 실행되므로, Solidity 스크립트만으로는 부족합니다.

#### 필요한 Go 코드 작업:

| 작업 | 상태 | 위치 | 설명 |
|------|------|------|------|
| **1. 배포 파이프라인 통합** | ❌ 필요 | `op-deployer/` | 새 스크립트를 Go 배포 파이프라인에 통합 |
| **2. Genesis 생성** | ⚠️ 부분 | `op-chain-ops/genesis/` | L1/L2 genesis 생성 |
| **3. Allocs 생성** | ⚠️ 부분 | `op-e2e/config/init.go` | State dump 생성 |
| **4. 주소 매핑** | ❌ 필요 | `op-e2e/config/` | 배포된 컨트랙트 주소 로드 |
| **5. Prestate 설정** | ❌ 필요 | `op-program/` | Cannon prestate 파일 |

### 2. 구체적인 미완료 작업

#### 문제 1: Genesis Output Root 계산 ⚠️

**현재 상태**:
```json
// deploy-config/devnetL1.json
{
  "faultGameGenesisOutputRoot": "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF"
}
```

**에러**:
```
ERROR Invalid prestate
err="failed to validate prestate: output root absolute prestate does not match:
Provider: 0x5e276326... (E2E가 계산한 실제 값)
Contract: 0xDEADBEEF... (deploy config의 더미 값)"
```

**필요한 작업**:
```bash
# 1. 실제 genesis output root 계산
go run scripts/calc-genesis-output-root.go

# 2. deploy-config 업데이트
vi packages/tokamak/contracts-bedrock/deploy-config/devnetL1-template.json

# 3. .devnet 재생성
./scripts/prepare-e2e-test.sh
```

**파일 위치**:
- `TODO-NEXT-SESSION.md:14-21` - 해결 방법 상세 설명
- `scripts/calc-genesis-output-root.go` - 계산 스크립트 (생성 필요)

---

#### 문제 2: E2E 배포 파이프라인 통합 ❌

**현재 상태**:
- Solidity 스크립트는 존재 ✅
- Go 배포 파이프라인에서 호출 안 됨 ❌

**E2E 테스트 실행 흐름**:
```
E2E 테스트 시작
  ↓
op-e2e/config/init() - 패키지 초기화
  ↓
initAllocType() - AllocType별 배포
  ↓
deployer.ApplyPipeline() ← 여기서 컨트랙트 배포
  ↓
Forge 스크립트 실행 (Deploy.s.sol)
  ↓
State dump 생성 (.devnet/allocs-l1.json)
  ↓
E2E 테스트 실행 (생성된 state 사용)
```

**필요한 작업**:
1. `op-deployer/pkg/deployer/pipeline/` 수정
   - 새 스크립트를 파이프라인에 통합
   - `DeployDisputeGame`, `SetDisputeGameImpl` 호출 추가

2. `op-e2e/config/init.go` 수정
   - 새 컨트랙트 주소 로드
   - Prestate 검증 로직 업데이트

**참고 문서**:
- `E2E-RUNTIME-DEPLOYMENT-GUIDE.md` - 배포 파이프라인 상세 설명

---

#### 문제 3: Prestate 파일 생성 ⚠️

**현재 상태**:
```bash
# prepare-e2e-test.sh에서 생성
op-program/bin/prestate-proof.json          ✅ 있음
op-program/bin/prestate-proof-mt64.json     ✅ 있음
op-program/bin/prestate-proof-mt64Next.json ✅ 있음
```

**문제**:
- 파일은 존재하지만 **내용이 올바른지 검증 필요**
- MIPS VM 버전과 일치하는지 확인 필요

**검증 방법**:
```bash
# Prestate hash 확인
cat op-program/bin/prestate-proof.json | jq .pre

# Deploy config와 비교
cat packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json | jq .faultGameAbsolutePrestate
```

**일치해야 하는 값**:
- `prestate-proof.json`의 `pre` 필드
- `devnetL1.json`의 `faultGameAbsolutePrestate`

---

#### 문제 4: 컨트랙트 주소 매핑 ❌

**현재 상태**:
- 새 스크립트로 배포된 컨트랙트 주소가 E2E에서 로드되지 않음

**필요한 파일**:
```
.devnet/addresses.json
```

**필요한 주소**:
```json
{
  "MipsSingleton": "0x...",
  "PreimageOracle": "0x...",
  "FaultDisputeGame": "0x...",
  "PermissionedDisputeGame": "0x...",
  "AlphabetVM": "0x..."
}
```

**필요한 작업**:
1. `op-e2e/config/deployment_loader.go` 수정
   - 새 컨트랙트 주소 로드 로직 추가

2. `op-e2e/config/init.go` 수정
   - 주소를 E2E 테스트에서 사용 가능하도록 설정

---

#### 문제 5: AllocType별 배포 설정 ⚠️

**AllocType이란?**
- E2E 테스트는 여러 배포 모드를 테스트
- `AllocTypeMTCannon`, `AllocTypeMTCannonNext`, `AllocTypeAltDA`

**현재 상태**:
```go
// op-e2e/config/init.go
const (
    AllocTypeAltDA        AllocType = "alt-da"
    AllocTypeMTCannon     AllocType = "mt-cannon"       // ✅ 준비 완료
    AllocTypeMTCannonNext AllocType = "mt-cannon-next"  // ⚠️ 스킵됨
    DefaultAllocType = AllocTypeMTCannon
)
```

**필요한 작업**:
- `MTCannon`에 대한 배포 설정 완료 확인
- `MTCannonNext`는 현재 스킵되어 있음 (OK)

---

## 📋 E2E 테스트 준비 체크리스트

### ✅ Solidity 레벨 (완료)

- ✅ StandardConstants.sol
- ✅ DeployPreimageOracle2.s.sol
- ✅ DeployMIPS2.s.sol
- ✅ DeployDisputeGame.s.sol
- ✅ SetDisputeGameImpl.s.sol
- ✅ DeployAlphabetVM.s.sol
- ✅ DeployImplementations.s.sol 통합
- ✅ forge build 성공

### ❌ Go 레벨 (미완료)

#### Critical (필수)
- ❌ Genesis output root 계산 및 설정
- ❌ E2E 배포 파이프라인 통합 (`op-deployer`)
- ❌ 컨트랙트 주소 매핑 (`op-e2e/config`)
- ⚠️ Prestate 파일 검증

#### High (강력 권장)
- ❌ 배포 자동화 스크립트 개선
- ❌ E2E 테스트 실행 가이드 작성
- ❌ CI/CD 통합

#### Medium (권장)
- ❌ 로깅 및 디버깅 도구
- ❌ 성능 모니터링

---

## 🔄 작업 순서

### Phase 1: Genesis Output Root (필수) - 30분

```bash
# 1. 계산 스크립트 작성
vi scripts/calc-genesis-output-root.go

# 2. Output root 계산
go run scripts/calc-genesis-output-root.go

# 3. Deploy config 업데이트
vi packages/tokamak/contracts-bedrock/deploy-config/devnetL1-template.json

# 4. 재배포
rm -rf .devnet
./scripts/prepare-e2e-test.sh
```

**참고**: `TODO-NEXT-SESSION.md`에 상세 가이드 있음

---

### Phase 2: E2E 배포 통합 (필수) - 2-3시간

```bash
# 1. op-deployer 수정
vi op-deployer/pkg/deployer/pipeline/implementations.go
vi op-deployer/pkg/deployer/pipeline/dispute_games.go

# 2. op-e2e 수정
vi op-e2e/config/init.go
vi op-e2e/config/deployment_loader.go

# 3. 테스트
cd op-e2e
go test -v ./faultproofs -run TestOutputCannonGame -timeout 30m
```

---

### Phase 3: Prestate 검증 (권장) - 30분

```bash
# 1. Prestate hash 추출
PRESTATE_HASH=$(cat op-program/bin/prestate-proof.json | jq -r .pre)

# 2. Deploy config와 비교
DEPLOY_PRESTATE=$(cat packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json | jq -r .faultGameAbsolutePrestate)

# 3. 일치 확인
if [ "$PRESTATE_HASH" == "$DEPLOY_PRESTATE" ]; then
    echo "✅ Prestate 일치"
else
    echo "❌ Prestate 불일치"
    echo "Expected: $DEPLOY_PRESTATE"
    echo "Got: $PRESTATE_HASH"
fi
```

---

### Phase 4: E2E 테스트 실행 (검증) - 30-40분

```bash
# 전체 faultproofs 테스트 실행
cd op-e2e
go test -v ./faultproofs -timeout 1h

# 개별 테스트 실행
go test -v ./faultproofs -run TestOutputCannonGame -timeout 30m
go test -v ./faultproofs -run TestPermissionedGame -timeout 30m
```

---

## 📊 예상 소요 시간

| Phase | 작업 | 소요 시간 |
|-------|------|-----------|
| 1 | Genesis output root 설정 | 30분 |
| 2 | E2E 배포 파이프라인 통합 | 2-3시간 |
| 3 | Prestate 검증 | 30분 |
| 4 | E2E 테스트 실행 | 30-40분 |
| **총계** | | **4-5시간** |

---

## 🎯 현재 상태 요약

### ✅ 준비된 것 (30%)

1. **Solidity 배포 스크립트** - 100% 완료
   - 모든 필수 스크립트 추가됨
   - Forge 빌드 성공
   - 기존 코드와 호환

2. **문서화** - 90% 완료
   - 배포 가이드 작성
   - 분석 문서 작성
   - TODO 문서 정리

### ❌ 필요한 것 (70%)

1. **Go 코드 통합** - 0% 완료
   - 배포 파이프라인 수정
   - 주소 로딩 로직
   - E2E 설정

2. **Genesis 설정** - 30% 완료
   - Output root 계산 필요
   - Deploy config 업데이트 필요
   - Prestate 검증 필요

3. **E2E 테스트 실행** - 0% 완료
   - 테스트 실행 안 됨
   - 에러 디버깅 필요

---

## 💡 핵심 인사이트

### 1. Solidity vs Go

**중요한 차이점**:
```
Solidity 스크립트 (완료)
  └─ 컨트랙트 배포 로직
  └─ Forge로 실행

Go 코드 (미완료)
  └─ E2E 테스트 실행 로직
  └─ Genesis 생성
  └─ State dump 생성
  └─ 테스트 오케스트레이션
```

**결론**: Solidity만으로는 E2E 테스트를 실행할 수 없음!

### 2. 배포 스크립트 ≠ E2E 준비

**배포 스크립트 추가**:
- ✅ 컨트랙트를 배포할 **방법** 제공
- ✅ 모듈화, 재사용성 개선
- ✅ 옵티미즘 표준 준수

**E2E 테스트 준비**:
- ❌ 배포 스크립트를 **자동으로 실행**하는 로직 필요
- ❌ 배포된 컨트랙트를 **테스트에서 사용**하는 로직 필요
- ❌ Genesis 생성 및 검증 로직 필요

### 3. 남은 작업의 성격

**배포 스크립트 추가** (완료):
- 난이도: 중간
- 작업량: 적음 (파일 복사 + 통합)
- 영향: 큼 (표준 준수)

**E2E 테스트 준비** (미완료):
- 난이도: 높음
- 작업량: 많음 (Go 코드 수정)
- 영향: 큼 (테스트 가능 여부)

---

## 🚀 다음 단계

### 즉시 시작 (Critical)

1. **Genesis output root 계산**
   - 참고: `TODO-NEXT-SESSION.md`
   - 소요 시간: 30분

2. **E2E 배포 파이프라인 통합**
   - 참고: `E2E-RUNTIME-DEPLOYMENT-GUIDE.md`
   - 소요 시간: 2-3시간

### 검증 (High)

3. **Prestate 검증**
   - 파일이 올바른지 확인
   - 소요 시간: 30분

4. **E2E 테스트 실행**
   - 전체 통합 테스트
   - 소요 시간: 30-40분

---

## 📚 관련 문서

| 문서 | 내용 | 상태 |
|------|------|------|
| `TODO-NEXT-SESSION.md` | Genesis output root 해결 가이드 | ✅ 작성 완료 |
| `E2E-RUNTIME-DEPLOYMENT-GUIDE.md` | 배포 파이프라인 상세 설명 | ✅ 작성 완료 |
| `ESSENTIAL-DEPLOYMENT-SCRIPTS-FOR-THANOS.md` | 배포 스크립트 목록 | ✅ 업데이트 필요 |
| `DISPUTE-GAME-DEPLOYMENT-GUIDE.md` | DisputeGame 사용 가이드 | ✅ 작성 완료 |

---

## ⚠️ 주의사항

### 1. 배포 스크립트만으로는 부족

**배포 스크립트**:
- ✅ 개발/테스트 환경에서 **수동 배포** 가능
- ✅ 프로덕션 배포에 **사용 가능**
- ❌ E2E 테스트 **자동 실행** 불가

### 2. Go 코드 수정 필수

E2E 테스트는 Go로 작성되어 있으므로:
- Go 코드 수정 없이는 E2E 테스트 실행 불가
- 배포 파이프라인 통합 필수

### 3. 단계적 접근 권장

모든 작업을 한 번에 하지 말고:
1. Genesis output root 먼저 해결 (간단)
2. E2E 파이프라인 통합 (복잡)
3. 테스트 실행 및 디버깅 (반복)

---

## 🎉 결론

### 질문: "배포 스크립트만 수정하면 E2E 테스트 준비가 완료되었을까?"

### 답변: **아니오 (30% 완료)**

**완료된 것**:
- ✅ Solidity 배포 스크립트 (100%)
- ✅ 문서화 (90%)

**필요한 것**:
- ❌ Go 배포 파이프라인 통합 (0%)
- ⚠️ Genesis output root 설정 (30%)
- ❌ E2E 테스트 실행 (0%)

**예상 작업 시간**: 4-5시간

**다음 단계**:
1. Genesis output root 계산 (30분)
2. E2E 배포 파이프라인 통합 (2-3시간)
3. Prestate 검증 (30분)
4. E2E 테스트 실행 (30-40분)

---

**작성자**: Claude Code
**마지막 업데이트**: 2025-11-07
**다음 작업**: `TODO-NEXT-SESSION.md` 참고
