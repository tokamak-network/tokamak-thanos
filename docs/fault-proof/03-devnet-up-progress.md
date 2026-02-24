# Fault-proof Devnet 검증 진행 상황

## 실행: 2026-02-23 14:58 KST
- 브랜치: `feat/fault-proof`
- 커밋 (시작): `45ddee89cddd74bd856e6d9ff07c0a5cd8f45aea`
- 결과: **FAIL (하드 블로커)**

## PASS/FAIL 매트릭스

| 섹션 | 상태 | 증적 |
|---|---|---|
| 1) 사전 점검 | PARTIAL | 브랜치를 `feat/fault-proof`로 고정; 실행 전 clean tree 확인; 도구 버전 수집 완료 |
| 2) 빌드 및 호환성 게이트 | FAIL | `make devnet-up`이 초기에는 submodule 해석 단계에서 실패 (`solady-v0.0.245`, 이후 `superchain-registry`) |
| 3) 부팅 검증 (`make devnet-up`) | FAIL | submodule 수정 후 `pre-devnet`에서 `go` 미설치로 실패 (`/bin/sh: go: not found`) |
| 4) Fault-Proof 런타임 검증 | BLOCKED | 부팅 성공 전에는 challenger 런타임 관측 불가 |
| 5) 회귀 게이트 | BLOCKED | L1/L2 서비스가 시작되지 않음 |
| 6) 재현성 재실행 | BLOCKED | 첫 완전 성공 실행이 아직 없음 |

## 명령어 + 주요 출력

1. `git pull --rebase`
   - `Already up to date.`

2. 첫 번째 `make devnet-up`
   - `fatal: No url found for submodule path 'packages/contracts-bedrock/lib/solady-v0.0.245' in .gitmodules`

3. 누락된 `.gitmodules` 항목 추가:
   - `packages/contracts-bedrock/lib/solady-v0.0.245`
   - `packages/contracts-bedrock/lib/superchain-registry`

4. 두 번째 `make devnet-up`
   - Submodule init은 성공적으로 진행됨.
   - `pre-devnet`에서 실패:
     - `./ops/scripts/geth-version-checker.sh: line 7: geth: command not found`
     - `/bin/sh: 4: go: not found`
     - `env: ‘go’: No such file or directory`

## 환경 스냅샷
- `go version`: not found
- `docker --version`: not found
- `docker compose version`: not found
- `python3 --version`: 3.12.3
- `node --version`: v22.22.0
- `pnpm --version`: 10.28.2

## 루트 원인
1. **브랜치의 저장소 설정 이슈**: 트리에 존재하는 두 gitlink 경로(`solady-v0.0.245`, `superchain-registry`)가 `.gitmodules`에서 누락되어, 재현 가능한 submodule bootstrap 실패를 유발했다.
2. **호스트 의존성 블로커**: devnet에 필요한 toolchain(`go`, 그리고 사실상 Docker daemon/CLI)이 없어 `pre-devnet` 단계에서 사전 요구사항 빌드/설치 및 서비스 기동이 불가능했다.

## 즉시 적용한 조치
- 누락된 두 submodule 항목을 `.gitmodules`에 추가하는 최소 범위 수정 적용.

## 필요한 의사결정 / 다음 조치
- 이 러너에 필요한 로컬 의존성 설치 및 프로비저닝:
  - Go (저장소가 기대하는 버전 이상; 현재 geth 설치 경로는 Go toolchain을 요구)
  - Docker + Docker Compose (devnet 서비스 필수)
- 환경 프로비저닝 후, 재현성 확인을 위해 `make devnet-up` 2회 성공을 포함한 전체 체크리스트 재실행.

---

## 재확인: 2026-02-23 15:13 KST
- 브랜치: `feat/fault-proof`
- 커밋: `58e950300cdf3c23e2d7307122dfa45bda742524`
- 결과: **FAIL (동일한 하드 블로커)**

### 증적
- `git pull --rebase` → `Already up to date.`
- `make devnet-up`은 이제 submodule init 단계를 통과하지만, 이후 `pre-devnet`에서 실패:
  - `./ops/scripts/geth-version-checker.sh: line 7: geth: command not found`
  - `/bin/sh: 4: go: not found`
  - `env: ‘go’: No such file or directory`

### 이전 실행 대비 변화
- `.gitmodules` 복구 이후 코드 레벨의 신규 실패는 발견되지 않음.
- 블로커는 여전히 환경 프로비저닝(Go/Docker 미설치)이며, 이로 인해 challenger 런타임 및 재현성 검증이 불가능함.

---

## 재확인: 2026-02-23 15:28 KST
- 브랜치: `feat/fault-proof`
- 커밋: `42d507c69e279e7227dc5292c1398fc40d5b0948`
- 결과: **FAIL (동일한 하드 블로커)**

### 증적
- `git pull --rebase` → `Already up to date.`
- 의존성 점검:
  - `go version` → not found
  - `docker --version` → not found
  - `docker compose version` → not found
- `make devnet-up`이 `pre-devnet`에서 동일 시그니처로 실패:
  - `geth: command not found`
  - `/bin/sh: 4: go: not found`
  - `env: ‘go’: No such file or directory`

### 이전 실행 대비 변화
- 저장소 레벨의 추가 결함은 식별되지 않음.
- 하드 블로커는 변함없음: 호스트 toolchain 미프로비저닝 상태로 인해 devnet 부팅 및 op-challenger 검증이 계속 차단됨.

---

## Re-check: 2026-02-23 15:43 KST
- Branch: `feat/fault-proof`
- Commit: `9e770dcd0dc729c4b73dd27bf736417a8002486c`
- Result: **FAIL (same hard blocker)**

### Evidence
- `git pull --rebase` → `Already up to date.`
- Dependency probe:
  - `go version` → not found
  - `docker --version` → not found
  - `docker compose version` → not found
- `make devnet-up` fails in `pre-devnet` with unchanged signature:
  - `./ops/scripts/geth-version-checker.sh: line 7: geth: command not found`
  - `/bin/sh: 4: go: not found`
  - `env: ‘go’: No such file or directory`

### Delta vs previous run
- No new code-level or config-level defects found.
- Blocker remains strictly environment provisioning (Go + Docker/Compose missing).

---

## 재검증: 2026-02-24 11:32 KST (Cycle 1)
- 브랜치: `feat/fault-proof`
- 커밋: `f44c3b92e37dff23aa52f4a9958db5791c969775`
- 결과: **PARTIAL (핵심 기동/재현성 성공, proposer 경고 패턴 리스크)**

### 체크리스트 상태 (Cycle 1)
| 섹션 | 상태 | 증적 |
|---|---|---|
| 1) 사전 점검 | PASS | `git status --porcelain` clean, 도구 버전 스냅샷 확보 (`go1.24.0`, `Docker 27.4.0`, `compose v2.31.0`) |
| 2) 빌드 및 호환성 게이트 | PASS | `make pre-devnet` 성공 (exit 0), `op-challenger` 이미지/환경변수 확인, unknown flag 시그니처 없음 |
| 3) 부팅 검증 (`make devnet-up`) | PASS | `make devnet-up` 성공 (exit 0), 핵심 서비스(`l1/l2/op-node/op-batcher/op-proposer/op-challenger`) up |
| 4) Fault-Proof 런타임 검증 | PARTIAL | challenger 시작/스케줄러/모니터링 로그 확인, restart 0; 활성 dispute 부재로 polling 동작의 풍부한 로그는 제한적 |
| 5) 회귀 게이트 | PARTIAL | L1/L2 블록 지속 증가, batcher 정상 제출; proposer에서 `failed to estimate gas: execution reverted` 경고 반복 후 정상 제출 지속 |
| 6) 재현성 재실행 | PASS | Cycle 2에서 `devnet-down` + `devnet-up` 동일 성공 재현 |

### 명령어 + 주요 출력
1. `make pre-devnet`  
   - exit: `0`  
   - log: `/tmp/fault-proof-pre-devnet-20260224-113244.log`
2. `make devnet-up`  
   - exit: `0`  
   - log: `/tmp/fault-proof-devnet-up-20260224-113303.log`
3. 안정성 샘플링(4회):  
   - `[11:37:16] l1=0x18b0 l2=0x4a3d restart=0`  
   - `[11:38:16] l1=0x18ba l2=0x4a5b restart=0`  
   - `[11:39:16] l1=0x18c4 l2=0x4a79 restart=0`  
   - `[11:40:16] l1=0x18ce l2=0x4a97 restart=0`
4. challenger 로그(요약):  
   - `Starting op-challenger`  
   - `starting scheduler`  
   - `starting monitoring`  
   - `challenger game service start completed`

### Delta vs previous run
- 기존 하드 블로커(Go/Docker 미설치)는 해소되었고 `make devnet-up`이 정상 완주됨.
- 남은 리스크는 proposer의 반복 경고(`estimate gas reverted`)이며, 현재 관측상 서비스 중단 없이 다음 제출 사이클로 회복됨.

---

## 재검증: 2026-02-24 11:42 KST (Cycle 2, 재현성)
- 브랜치: `feat/fault-proof`
- 커밋: `f44c3b92e37dff23aa52f4a9958db5791c969775`
- 결과: **PASS (재실행 성공)**

### 명령어 + 주요 출력
1. `make devnet-down`  
   - exit: `0`  
   - log: `/tmp/fault-proof-devnet-down-20260224-114210.log`
2. `make devnet-up` (재실행)  
   - exit: `0`  
   - log: `/tmp/fault-proof-devnet-up-20260224-114217.log`
3. 서비스 상태: `docker compose ps`에서 핵심 서비스 전부 `Up`
4. 블록 진행 샘플:
   - `[11:43:03] l1=0x18e7 l2=0x4aeb restart=0`
   - `[11:43:13] l1=0x18e9 l2=0x4af0 restart=0`
   - `[11:43:23] l1=0x18ea l2=0x4af5 restart=0`

### 비고
- Cycle 1과 동일하게 proposer 경고 패턴은 재관측됨.
- challenger는 재실행에서도 즉시 기동되고 restart 없이 유지됨.

---

## 재검증: 2026-02-24 12:30 KST (Cycle 3)
- 브랜치: `feat/fault-proof`
- 커밋: `f44c3b92e37dff23aa52f4a9958db5791c969775`
- 결과: **PARTIAL (기동/재현성 성공, proposer 경고 폭증 리스크 지속)**

### 체크리스트 상태 (Cycle 3)
| 섹션 | 상태 | 증적 |
|---|---|---|
| 1) 사전 점검 | PARTIAL | 브랜치/커밋/도구 버전 고정 완료. 단, 시작 시 워킹 트리는 clean 아님 (`docs/fault-proof/03-devnet-up-progress.md`, `docs/lessons.md`, `docs/todo.md`) |
| 2) 빌드 및 호환성 게이트 | PASS | `make pre-devnet` 성공 (exit 0), challenger 이미지/환경변수 확인, unknown flag 시그니처 없음 |
| 3) 부팅 검증 (`make devnet-up`) | PASS | `make devnet-up` 성공 (exit 0), 핵심 서비스 전부 `Up` |
| 4) Fault-Proof 런타임 검증 | PARTIAL | challenger start/scheduler/monitoring/game service 시작 완료 로그 확인, restart 0; dispute 부재로 polling 상세 로그는 제한적 |
| 5) 회귀 게이트 | PARTIAL | L1/L2 블록 지속 증가, batcher publish/confirm 정상. proposer는 경고 폭증(31회) + 영구 실패 1회 후 회복 |
| 6) 재현성 재실행 | PASS | Cycle 4에서 `devnet-down -> devnet-up` 재실행 성공 |

### 명령어 + 주요 출력
1. `make devnet-down`  
   - exit: `0`  
   - log: `/tmp/fault-proof-devnet-down-20260224-123038.log`
2. `make pre-devnet`  
   - exit: `0`  
   - log: `/tmp/fault-proof-pre-devnet-20260224-123055.log`
3. `make devnet-up`  
   - exit: `0`  
   - log: `/tmp/fault-proof-devnet-up-20260224-123101.log`
4. 안정성 샘플링(5회):  
   - `[12:31:50] sample=1 l1=0x1aca l2=0x50a2 challenger_restart=0 proposer_restart=0`  
   - `[12:32:50] sample=2 l1=0x1ad4 l2=0x50c0 challenger_restart=0 proposer_restart=0`  
   - `[12:33:50] sample=3 l1=0x1ade l2=0x50de challenger_restart=0 proposer_restart=0`  
   - `[12:34:50] sample=4 l1=0x1ae8 l2=0x50fd challenger_restart=0 proposer_restart=0`  
   - `[12:35:51] sample=5 l1=0x1af2 l2=0x511b challenger_restart=0 proposer_restart=0`
5. proposer 로그 집계 (15m 창):  
   - `WARN_COUNT=31` (`failed to estimate gas: execution reverted`)  
   - `PERMANENT_FAIL_COUNT=1` (`operation failed permanently after 30 attempts`)  
   - `SUCCESS_COUNT=16` (`proposer tx successfully published`)  
   - log: `/tmp/fault-proof-proposer-logs-20260224-123607.log`
6. challenger 로그 집계:  
   - `Starting op-challenger`  
   - `starting scheduler`  
   - `starting monitoring`  
   - `challenger game service start completed`  
   - log: `/tmp/fault-proof-challenger-logs-20260224-123607.log`

### Delta vs previous run
- Cycle 1/2와 동일하게 proposer 경고 패턴(경고 폭증 + 영구 실패 1회 + 이후 회복)이 반복 재현됨.
- challenger/l1/l2/batcher 관점의 기동 안정성은 동일하게 유지됨.

---

## 재검증: 2026-02-24 12:36 KST (Cycle 4, 재현성)
- 브랜치: `feat/fault-proof`
- 커밋: `f44c3b92e37dff23aa52f4a9958db5791c969775`
- 결과: **PASS (재실행 성공, 동일 리스크 재현)**

### 명령어 + 주요 출력
1. `make devnet-down`  
   - exit: `0`  
   - log: `/tmp/fault-proof-devnet-down-20260224-123636.log`
2. `make devnet-up` (재실행)  
   - exit: `0`  
   - log: `/tmp/fault-proof-devnet-up-20260224-123651.log`
3. 블록 진행 샘플:  
   - `[12:37:32] sample=1 l1=0x1afe l2=0x514d challenger_restart=0 proposer_restart=0`  
   - `[12:37:42] sample=2 l1=0x1b00 l2=0x5152 challenger_restart=0 proposer_restart=0`  
   - `[12:37:52] sample=3 l1=0x1b02 l2=0x5157 challenger_restart=0 proposer_restart=0`
4. proposer 패턴 재확인(재실행 후 2m 창):  
   - `WARN_COUNT=31`  
   - `PERMANENT_FAIL_COUNT=1`  
   - `SUCCESS_COUNT=4` (영구 실패 이후 회복 확인)  
   - log: `/tmp/fault-proof-proposer-recovery-20260224-123932.log`
5. challenger 재기동 로그:  
   - `Starting op-challenger` / `starting scheduler` / `starting monitoring` / `challenger game service start completed`  
   - log: `/tmp/fault-proof-challenger-repro-logs-20260224-123804.log`

### 비고
- 재현성 관점에서 `devnet-up` 성공 경로는 안정적으로 재현됨.
- proposer 경고 폭증 패턴은 2차 재기동에서도 동일하게 재현되어, 기능 저하 리스크가 지속됨.

---

## 재검증: 2026-02-24 13:59 KST (Cycle 5, dispute lifecycle + bond refund)
- 브랜치: `feat/fault-proof`
- 커밋: `f44c3b92e37dff23aa52f4a9958db5791c969775`
- 결과: **PARTIAL (lifecycle 완주 성공, proposer 경고 리스크 지속)**

### 명령어 + 주요 출력
1. 설정 반영/초기화
   - `make devnet-clean` -> `make devnet-up`
   - `faultGameMaxClockDuration=120`, `faultGameWithdrawalDelay=0`, `faultGameClockExtension=0` 확인
2. dispute game 생성 (type 0)
   - create tx: `0x53db92ef3dc5b296debf4e9d5a6e9df69d4b3f5ea23469192ed5d18f05768bd3`
   - game: `0x4D67490e0D3FE0f3Ca16C7d0E6D64785E553c612`
3. move(attack)로 bond deposit 발생
   - attack tx: `0xab96797d9df0aabecb465e27665e41eff1fdcf23ff76f6ad1dda5c4355df1ffc`
   - claim #1 bond: `0.09132520 ETH`
   - pre-resolve 확인: `DelayedWETH.balanceOf(game)=91325200000000000`
4. 종료 단계
   - `resolve` 1차 직접 시도는 `failed to estimate gas: execution reverted`로 30회 후 영구 실패
   - `resolve-claim(1)` tx: `0x7a2a0f713d8497fd5d175f7b08ec9ab25348678e8ef1a4abc629190406e6fb76`
   - `resolve-claim(0)` tx: `0x94dc310ef3c5072dccaa3b910c11eed01c5038a0000e550df3a79f0fdb316ec0`
   - `resolve` 재시도 tx: `0xfdc4f83745e3d7f1595cabcb528bda0d8e96478f9f489ef0b4cb204d0ab01d27` (성공)
5. bond refund
   - `claimCredit` tx: `0xd63c72a0f630b8e2a9842eb7e12b15bc3facddb2b9c2c2b5178d0ab0454c386d`
   - post-check: `credit(attacker)=0`, `DelayedWETH.balanceOf(game)=0`
   - game status: `Challenger Won`

### 실행 증적
- `/tmp/fault-proof-devnet-clean-20260224-1351.log`
- `/tmp/fault-proof-devnet-up-20260224-1351-clean.log`
- `/tmp/fault-proof-dispute-bond-20260224-1358.log`

### Delta vs previous run
- 기존 Cycle 3/4의 “활성 dispute 부재” 공백을 해소하고, 실제 dispute lifecycle + bond 환급까지 검증 완료.
- proposer 경고 폭증 패턴은 dispute lifecycle 검증 중에도 동일하게 관측되어 잔여 리스크로 유지.
