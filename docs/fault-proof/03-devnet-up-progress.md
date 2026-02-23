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
