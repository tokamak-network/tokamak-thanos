# Fault Proof 업데이트 1차 분석 (업스트림 커밋 기반)

- 기준 저장소: `ethereum-optimism/optimism`
- 기준 커밋: `5401e3a46d4e193bbe5400b9c3948d8098c8692d`
- 분석 브랜치: `feat/fault-proof`
- 분석 일시: 2026-02-18 (KST)

---

## 1) 업스트림 커밋 영향도 매핑

## 1-1. 업스트림 커밋 변경 파일
해당 커밋은 아래 **1개 파일만** 변경한다.

- `kona-proofs/version.json`

커밋 메시지:
- `proofs: Update pinned kona client version (#18611)`

## 1-2. tokamak-thanos에 동일 경로 존재 여부
`tokamak-thanos` 기준으로 `kona-proofs/version.json` 경로가 **존재하지 않음**.

### 결론
- **직접 머지 영향도: 낮음 (Low)**
- 이유: 변경 대상 파일 경로 자체가 현재 레포에 없음.

## 1-3. 간접 영향(주의 포인트)
해당 커밋의 본질은 “FP 경로에서 사용하는 proofs client(kona) 버전 핀 업데이트”다.
즉, 파일이 없더라도 아래 관점은 확인 필요:

1. 현재 tokamak-thanos가 kona 기반 proofs 경로를 사용 중인지
2. 미사용이라면 cannon/op-program 기반 FP 운영 정책으로 고정할지
3. 향후 kona 도입 계획이 있다면, 별도 디렉토리/버전 핀 체계를 신규 도입해야 하는지

## 1-4. Tokamak 커스텀 영역별 영향도

### A. `packages/tokamak/contracts-bedrock`
- **직접 영향 없음**
- 단, FP 정책 변경 시 `DisputeGameFactory`, `OptimismPortal2` 연동 설정/권한 정책은 간접 영향 가능

### B. `packages/tokamak/sdk`
- **직접 영향 없음**
- SDK는 브릿지/메시징 계층이 중심이라 proofs client 버전 핀 변경과 직접 결합도 낮음

### C. 배포/운영 설정(`deploy-config`, `ops*`, `op-*`)
- **간접 영향 가능 (중간)**
- 이유: FP 실행 경로를 kona로 확장/전환할 경우 런타임 바이너리/버전 관리 포인트가 생김

---

## 2) FP 필수 컴포넌트 체크리스트

아래는 “Optimism latest Fault Proof 지원” 목적의 필수 체크 항목이다.

상태 기준:
- ✅ 존재/확인됨
- ⚠️ 존재는 하나 정책/설정 검증 필요
- ❌ 부재/추가 필요

## 2-1. 실행 바이너리/서비스 계층

- ✅ `op-challenger` 디렉토리 존재
- ✅ `op-dispute-mon` 디렉토리 존재
- ✅ `op-program` 디렉토리 존재
- ✅ `cannon` 디렉토리 존재
- ❌ `kona-proofs/version.json` (업스트림 커밋 대상) 경로 부재

## 2-2. L1 컨트랙트/FP 코어

- ⚠️ `DisputeGameFactory` 주소/연결이 체인별로 일관되게 설정되어 있는지 점검 필요
- ⚠️ `OptimismPortal2` 사용 경로(Portal/Portal2 동시 표기) 정합성 점검 필요
- ⚠️ `SystemConfig` ↔ `ProxyAdmin` ↔ Safe 검증 로직이 FP 배포/업그레이드 시나리오와 충돌 없는지 점검 필요

## 2-3. 체인 설정/파라미터

- ⚠️ 체인별 `finalizationPeriodSeconds`가 FP 운영 정책과 일치하는지 확인
- ⚠️ proposer/challenger 권한 분리(EOA/멀티시그) 검증
- ⚠️ `DisputeGameFactory` 주소가 zero-address로 남아있는 네트워크 정리

## 2-4. 운영/관측

- ⚠️ `op-dispute-mon` 알림 기준(경보 임계치/중복 억제) 명문화 필요
- ⚠️ 장애 시 롤백 절차(Portal2/DisputeGameFactory/Challenger) 런북 필요

## 2-5. 테스트/검증

- ⚠️ devnet에서 dispute 생성 → 반박/해결 → 최종 상태 전이 시나리오 점검 필요
- ⚠️ e2e에 Tokamak 커스텀(USDC bridge, verification) 회귀 테스트 포함 필요
- ⚠️ FP 경로 변경이 브릿지/메시징 성능·안정성에 미치는 영향 회귀 검증 필요

---

## 3) 바로 실행할 액션 (다음 단계)

1. `kona-proofs` 채택 여부 결정
   - (A) 미채택: cannon/op-program 경로 유지 선언
   - (B) 채택: `kona-proofs` 디렉토리/버전 핀 체계 신규 도입

2. 네트워크별 FP 핵심 주소 테이블 확정
   - `OptimismPortal2`, `DisputeGameFactory`, challenger/proposer 권한

3. 검증 체크리스트를 CI 게이트로 일부 승격
   - 최소: 설정 정합성 + zero-address 방지 + 핵심 컴포넌트 존재성

---

## 부록) 확인 명령(재현용)

```bash
# 업스트림 커밋 변경 파일 확인
git fetch https://github.com/ethereum-optimism/optimism.git 5401e3a46d4e193bbe5400b9c3948d8098c8692d
git show --name-only --pretty=format:'' FETCH_HEAD

# tokamak-thanos 내 kona 관련 경로 확인
find . -maxdepth 3 -type d -name '*kona*' -o -type f -name 'version.json'
```
