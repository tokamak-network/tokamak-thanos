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

---

## 검증 결과 (2026-02-24 12:30 KST, Cycle 3/4)
- **결과:** PARTIAL (**NO-GO**; proposer 경고 폭증 리스크 미해소)
- **환경 + 커밋:** `feat/fault-proof` @ `f44c3b92e37dff23aa52f4a9958db5791c969775`
- **환경 스냅샷:** Go `1.24.0`, Docker `27.4.0`, Compose `v2.31.0-desktop.2`, Python `3.14.2`, Node `v20.16.0`, pnpm `10.8.0`
- **사전 점검 특이사항:** 워킹 트리는 검증 시작 시점부터 문서 파일 변경 상태였음 (`docs/fault-proof/03-devnet-up-progress.md`, `docs/lessons.md`, `docs/todo.md`)

### 체크리스트 판정
| 섹션 | 상태 | 증적 |
|---|---|---|
| 1) 사전 점검 | PARTIAL | 브랜치/커밋/도구 버전 고정은 완료, 워킹 트리는 clean 조건 미충족 |
| 2) 빌드 및 호환성 게이트 | PASS | `make pre-devnet` 성공 (`/tmp/fault-proof-pre-devnet-20260224-123055.log`), challenger 이미지 `tokamaknetwork/thanos-op-challenger:latest`, unknown flag 시그니처 없음 |
| 3) 부팅 검증 (`make devnet-up`) | PASS | 1차/2차 모두 `make devnet-up` 성공, 핵심 서비스 `Up` |
| 4) Fault-Proof 런타임 검증 | PARTIAL | challenger 시작/스케줄러/모니터링/서비스 시작 완료 로그 확인, restart `0`; 활성 dispute 부재로 polling 동작의 풍부한 런타임 로그는 제한적 |
| 5) 회귀 게이트 | PARTIAL | L1/L2 블록 지속 증가, batcher publish/confirm 정상 지속. proposer는 `failed to estimate gas: execution reverted` 경고 폭증(사이클당 31회) + `Failed to send proposal transaction` 1회 후 회복 |
| 6) 재현성 재실행 | PASS | `make devnet-down` 후 동일 config로 `make devnet-up` 재실행 성공, 서비스/블록 진행 신호 재현 |

### 실행 증적 (주요 로그)
- `/tmp/fault-proof-devnet-down-20260224-123038.log`
- `/tmp/fault-proof-pre-devnet-20260224-123055.log`
- `/tmp/fault-proof-devnet-up-20260224-123101.log`
- `/tmp/fault-proof-runtime-observe-20260224-123150.log`
- `/tmp/fault-proof-proposer-logs-20260224-123607.log`
- `/tmp/fault-proof-challenger-logs-20260224-123607.log`
- `/tmp/fault-proof-batcher-logs-20260224-123616.log`
- `/tmp/fault-proof-devnet-down-20260224-123636.log`
- `/tmp/fault-proof-devnet-up-20260224-123651.log`
- `/tmp/fault-proof-repro-samples-20260224-123732.log`
- `/tmp/fault-proof-proposer-recovery-20260224-123932.log`

### 권고
- **권고:** `fix-before-merge`  
  - 근거: proposer의 estimate-gas revert 경고 폭증 및 주기적 영구 실패 로그(`operation failed permanently after 30 attempts`)가 반복 재현됨

---

## 검증 결과 (2026-02-24 13:59 KST, Cycle 5 - dispute lifecycle)
- **결과:** PARTIAL (**NO-GO**; dispute lifecycle 완주 성공, proposer 경고 리스크 지속)
- **환경 + 커밋:** `feat/fault-proof` @ `f44c3b92e37dff23aa52f4a9958db5791c969775`
- **핵심 설정:** `faultGameMaxClockDuration=120`, `faultGameWithdrawalDelay=0`, `faultGameClockExtension=0`

### 체크리스트 판정
| 섹션 | 상태 | 증적 |
|---|---|---|
| 1) 사전 점검 | PARTIAL | 워킹 트리 변경 상태 유지, 커밋/브랜치/도구는 고정 |
| 2) 빌드 및 호환성 게이트 | PASS | `make devnet-clean` 후 `make devnet-up` 성공 |
| 3) 부팅 검증 (`make devnet-up`) | PASS | 핵심 서비스(`l1/l2/op-node/op-batcher/op-proposer/op-challenger`) `Up` |
| 4) Fault-Proof 런타임 검증 | PASS | dispute game 생성→move(bond deposit)→resolve-claim→resolve→claimCredit(refund) 완주 |
| 5) 회귀 게이트 | PARTIAL | proposer의 `failed to estimate gas: execution reverted` 경고 패턴은 동일하게 관측 |
| 6) 재현성 재실행 | PASS | `devnet-clean`으로 상태 초기화 후 동일 시나리오 재수행 가능성 확인 |

### dispute lifecycle 증적
- 수집 로그: `/tmp/fault-proof-dispute-bond-20260224-1358.log`
- 게임 주소: `0x4D67490e0D3FE0f3Ca16C7d0E6D64785E553c612` (Type `0`)
- 생성 tx: `0x53db92ef3dc5b296debf4e9d5a6e9df69d4b3f5ea23469192ed5d18f05768bd3`
- 공격(move) tx: `0xab96797d9df0aabecb465e27665e41eff1fdcf23ff76f6ad1dda5c4355df1ffc`
  - claim #1 bond: `0.09132520 ETH`
- `resolve` 1차 시도: `operation failed permanently after 30 attempts: failed to estimate gas: execution reverted`
- 해소 순서:
  - `resolve-claim(1)` tx: `0x7a2a0f713d8497fd5d175f7b08ec9ab25348678e8ef1a4abc629190406e6fb76`
  - `resolve-claim(0)` tx: `0x94dc310ef3c5072dccaa3b910c11eed01c5038a0000e550df3a79f0fdb316ec0`
  - `resolve` tx: `0xfdc4f83745e3d7f1595cabcb528bda0d8e96478f9f489ef0b4cb204d0ab01d27`
- 환급(claimCredit) tx: `0xd63c72a0f630b8e2a9842eb7e12b15bc3facddb2b9c2c2b5178d0ab0454c386d`
  - post-check: `credit(attacker)=0`, `DelayedWETH.balanceOf(game)=0`, game status=`1`(`Challenger Won`)

### 권고
- **권고:** `fix-before-merge`
  - 근거: dispute lifecycle 자체는 완주되었지만 proposer estimate-gas revert 경고 폭증 및 주기적 영구 실패 패턴은 여전히 남아 있음
