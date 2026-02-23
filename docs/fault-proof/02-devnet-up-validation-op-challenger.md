# Devnet-up 검증 계획 (op-challenger 업그레이드)

## 목표
`feat/fault-proof` 브랜치가 `make devnet-up`으로 로컬 devnet을 정상 기동할 수 있고, 업그레이드된 `op-challenger`가 기대한 fault-proof 동작과 함께 안정적으로 실행되는지 검증한다.

## 범위
- 브랜치: `feat/fault-proof`
- 중점: `op-challenger` 업그레이드가 로컬 devnet 부팅 및 런타임 동작에 미치는 영향
- 환경: 로컬 개발자 머신

## Go / No-Go 기준
아래 조건을 모두 만족할 때만 **GO**로 판단한다:
1. `make devnet-up`이 devnet을 성공적으로 기동한다 (exit code + 서비스 헬스 OK)
2. `op-challenger`가 기동 후 정상 상태를 유지한다 (crash loop/panic 없음)
3. 핵심 플로우에 치명적 회귀가 없다 (L1/L2 진행, batch/proposer 정상)
4. Fault-proof 관측 지점이 통과한다 (game polling/connection 동작이 기대와 일치)
5. 동일 실행을 최소 2회 재현할 수 있다

## 실행 체크리스트

### 1) 사전 점검 (재현성)
- [ ] 실행 대상 브랜치와 커밋 해시가 고정되었는지 확인
- [ ] 워킹 트리가 깨끗한지 확인
- [ ] devnet-up에 사용한 핵심 env/config 스냅샷 저장
- [ ] 로컬 의존성/도구 버전 기록 (Go, Docker 등)

### 2) 빌드 및 호환성 게이트
- [ ] devnet-up에서 사용하는 빌드 타깃이 compile/link 오류 없이 완료되는지 확인
- [ ] 이 브랜치가 사용하는 `op-challenger` 바이너리/이미지 경로 확인
- [ ] 업그레이드된 flags/env vars(추가/삭제/이름 변경) 검증
- [ ] 업그레이드 이후 미해결 참조가 없는지 확인

### 3) 부팅 검증 (`make devnet-up`)
- [ ] `make devnet-up` 실행
- [ ] 모든 핵심 서비스의 시작 로그 수집
- [ ] 기대 서비스가 모두 올라왔는지 확인 (L1, L2, batcher, proposer, challenger)
- [ ] 최초 5~10분 로그에서 panic/restart storm/retry 폭증 여부 관찰

### 4) Fault-Proof 런타임 검증 (op-challenger)
- [ ] challenger가 설정된 RPC/contracts에 연결 가능한지 확인
- [ ] challenger polling loop가 정상 동작하는지 확인
- [ ] 기본 시나리오 검증: 활성 dispute 없음, 비정상 오류 없음
- [ ] (선택) 로컬 환경에서 dispute 유사 조건을 유발해 동작/로그 확인

### 5) 회귀 게이트
- [ ] L2 블록 진행이 기대대로 지속되는지 확인
- [ ] Batcher/proposer가 정상 동작을 계속하는지 확인
- [ ] 이전 베이스라인 대비 error/warn 프로파일에 치명적 증가가 없는지 확인
- [ ] 실행 중 핵심 RPC endpoint가 안정적으로 응답하는지 확인

### 6) 재현성 재실행
- [ ] 환경 teardown
- [ ] 동일 config로 `make devnet-up` 재실행
- [ ] 동일 성공 기준을 다시 만족하는지 확인

## 수집 증적
- 사용한 커밋 해시
- 실행한 정확한 명령어
- 서비스 상태 스냅샷
- Challenger 로그 발췌 (시작, steady-state, 오류 발생 시)
- 체크리스트 섹션별 최종 pass/fail 매트릭스

## 실패 처리
어떤 체크라도 실패하면 아래를 보고한다:
1. 실패한 단계
2. 증상/로그 시그니처
3. 추정 루트 원인
4. 즉시 조치(remediation)
5. 재검증 결과

## 산출물 형식
각 검증 사이클마다 아래 형식의 짧은 보고서를 작성한다:
- **결과:** PASS / FAIL
- **환경 + 커밋**
- **통과 항목**
- **실패 항목 / 리스크**
- **권고:** merge-forward / fix-before-merge
