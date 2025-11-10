# OPContractsManager Assertion 이슈 디버깅

**작성일**: 2025-11-10
**상태**: 🔍 원인 파악 완료

---

## 🎯 문제 요약

E2E 테스트에서 `DeployImplementations.s.sol`의 `assertValidOutput()` 실행 중 **OPContractsManager assertion에서 실패**

---

## 🔍 디버깅 과정

### 1단계: 로그 추가

**추가한 로그**:
- `ChainAssertions.checkOPContractsManager`: 각 require 문마다 로그 추가
- `DeployImplementations.s.sol`: Deploy 단계 로그

**결과**:
```
INFO [11-10|14:19:16.956] Running chain assertions on the OPContractsManager at 0x9676558753665831C1e29B2f0b05C1dBd7eb9b56
INFO [11-10|14:19:16.956] CHECK-OPCM-10: OK
INFO [11-10|14:19:16.956] "CHECK-OPCM-15: OK, version=3.0.0"
INFO [11-10|14:19:16.956] CHECK-OPCM-17: OK
INFO [11-10|14:19:16.956] CHECK-OPCM-18: OK
INFO [11-10|14:19:16.956] CHECK-OPCM-19: OK
INFO [11-10|14:19:16.956] Getting implementations from OPCM...
WARN [11-10|14:19:16.956] Fault addr=0x8f0818c1e0a4cA9Bc005f61428f7f8a371464D38 label=DeployImplementations err="execution reverted"
```

### 2단계: 실패 지점 특정

**실패 위치**:
```solidity
IOPContractsManager.Implementations memory impls = _opcm.implementations();
```

**증거**:
- ✅ "Getting implementations from OPCM..." 로그 출력됨
- ❌ "Got implementations" 로그 출력되지 않음
- ❌ CHECK-OPCM-50 로그도 출력되지 않음

**결론**: `_opcm.implementations()` 호출이 revert하고 있음

---

## 🚨 핵심 문제

**`IOPContractsManager.implementations()` 함수 호출이 실패**

### 가능한 원인들

1. **OPContractsManager 구현 문제**
   - `implementations()` getter가 제대로 구현되지 않음
   - Implementation 데이터가 제대로 저장되지 않음

2. **Storage layout 이슈**
   - OPContractsManager의 storage layout이 Optimism v1.16.0과 다름
   - `implementations()` 함수가 참조하는 storage slot이 비어있음

3. **Initialization 문제**
   - OPContractsManager가 제대로 초기화되지 않음
   - Implementation 주소들이 설정되지 않음

---

## 📋 검증 완료 항목

### ✅ 성공한 Assertion들

- `address(_opcm) != address(0)` - CHECK-OPCM-10 ✅
- `bytes(_opcm.version()).length > 0` - CHECK-OPCM-15 ✅ (version="3.0.0")
- `address(_opcm.protocolVersions()) == _proxies.ProtocolVersions` - CHECK-OPCM-17 ✅
- `address(_opcm.superchainProxyAdmin()) == address(_superchainProxyAdmin)` - CHECK-OPCM-18 ✅
- `address(_opcm.superchainConfig()) == _proxies.SuperchainConfig` - CHECK-OPCM-19 ✅

**참고**: 위의 5개 getter는 모두 성공했으므로, OPContractsManager 자체는 배포되어 있음.

### ❌ 실패한 Assertion

- `_opcm.implementations()` - **호출 자체가 revert** ❌

---

## 🔬 다음 단계: 상세 조사

### 1. OPContractsManager 소스 코드 확인
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
grep -n "function implementations" src -r
```

**확인 사항**:
- `implementations()` 함수가 존재하는가?
- 함수 시그니처가 올바른가?
- 함수 내부에서 어떤 storage에 접근하는가?

### 2. deployOPContractsManager 함수 확인
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
grep -A 50 "function deployOPContractsManager" scripts/deploy/DeployImplementations.s.sol
```

**확인 사항**:
- OPContractsManager가 어떻게 배포되는가?
- Implementation 주소들이 제대로 설정되는가?
- Initialize 과정이 있는가?

### 3. Optimism 원본과 비교
```bash
cd /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock
grep -n "function implementations" src/L1/OPContractsManager.sol
```

**확인 사항**:
- Tokamak-Thanos와 Optimism 원본의 차이점은?
- `implementations()` 함수 구현이 다른가?
- Storage layout이 다른가?

### 4. Forge trace로 상세 디버깅
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock
forge script scripts/deploy/DeployImplementations.s.sol --sig "run(...)" -vvvv
```

**확인 사항**:
- Revert가 발생하는 정확한 위치
- Revert reason 메시지
- Stack trace

---

## 💡 예상 해결 방법

### 방법 1: OPContractsManager Implementation 수정
OPContractsManager의 `implementations()` 함수 또는 관련 storage가 제대로 설정되지 않은 경우

### 방법 2: deployOPContractsManager 수정
OPContractsManager 배포 시 implementation 주소들을 제대로 설정하지 않은 경우

### 방법 3: Optimism 원본 파일 복사
Tokamak-Thanos의 OPContractsManager가 Optimism v1.16.0과 다른 경우, 원본 파일 복사 필요

---

## 📊 현재 상태

| Component | Status | Notes |
|-----------|--------|-------|
| OPContractsManager 배포 | ✅ 성공 | address != 0, version() 동작 |
| OPContractsManager getters (5개) | ✅ 성공 | protocolVersions, superchainProxyAdmin, superchainConfig 모두 동작 |
| OPContractsManager.implementations() | ❌ 실패 | **호출이 revert** |
| L1StandardBridge 마이그레이션 | ✅ 완료 | ChainAssertions 통과 |

---

## 🎯 긴급 액션 아이템

1. **즉시**: OPContractsManager 소스 코드 확인
   - `implementations()` 함수 존재 여부
   - 함수 구현 내용

2. **다음**: deployOPContractsManager 함수 확인
   - Implementation 설정 로직
   - Initialize 과정

3. **필요시**: Optimism 원본과 비교
   - 차이점 파악
   - 필요한 수정사항 도출

---

## 📚 참고 자료

### 로그 위치
- Full E2E test log: 위 디버깅 과정 참조

### 관련 파일
- `/packages/tokamak/contracts-bedrock/scripts/deploy/ChainAssertions.sol`
- `/packages/tokamak/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`
- `/packages/tokamak/contracts-bedrock/src/L1/OPContractsManager.sol` (확인 필요)

### 에러 메시지
```
Caught revision id error: revision id 33 cannot be reverted
ERROR DeployImplementations.Run FAILED
ERROR execution reverted at address 0x8f0818c1e0a4cA9Bc005f61428f7f8a371464D38
```

---

**작성자**: AI Assistant
**상태**: 원인 파악 완료, 해결 방법 조사 중

