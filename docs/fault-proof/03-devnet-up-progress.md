# Fault-proof Devnet Validation Progress

## Run: 2026-02-23 14:58 KST
- Branch: `feat/fault-proof`
- Commit (start): `45ddee89cddd74bd856e6d9ff07c0a5cd8f45aea`
- Result: **FAIL (hard blocker)**

## PASS/FAIL Matrix

| Section | Status | Evidence |
|---|---|---|
| 1) Preflight | PARTIAL | Branch fixed at `feat/fault-proof`; clean tree before run; tool versions captured |
| 2) Build & Compatibility Gate | FAIL | `make devnet-up` initially failed at submodule resolution (`solady-v0.0.245`, then `superchain-registry`) |
| 3) Boot Validation (`make devnet-up`) | FAIL | After submodule fix, pre-devnet fails because `go` is not installed (`/bin/sh: go: not found`) |
| 4) Fault-Proof Runtime Validation | BLOCKED | Challenger runtime cannot be observed before successful boot |
| 5) Regression Gate | BLOCKED | L1/L2 services not started |
| 6) Reproducibility Re-run | BLOCKED | First full successful run unavailable |

## Commands + Key Output

1. `git pull --rebase`
   - `Already up to date.`

2. First `make devnet-up`
   - `fatal: No url found for submodule path 'packages/contracts-bedrock/lib/solady-v0.0.245' in .gitmodules`

3. Added missing `.gitmodules` entries:
   - `packages/contracts-bedrock/lib/solady-v0.0.245`
   - `packages/contracts-bedrock/lib/superchain-registry`

4. Second `make devnet-up`
   - Submodule init proceeds successfully.
   - Fails in `pre-devnet`:
     - `./ops/scripts/geth-version-checker.sh: line 7: geth: command not found`
     - `/bin/sh: 4: go: not found`
     - `env: ‘go’: No such file or directory`

## Environment Snapshot
- `go version`: not found
- `docker --version`: not found
- `docker compose version`: not found
- `python3 --version`: 3.12.3
- `node --version`: v22.22.0
- `pnpm --version`: 10.28.2

## Root Cause
1. **Repository config issue on branch**: `.gitmodules` omitted two gitlink paths present in tree (`solady-v0.0.245`, `superchain-registry`), causing deterministic submodule bootstrap failure.
2. **Host dependency blocker**: required devnet toolchain is absent (`go`, likely Docker daemon/CLI), so `pre-devnet` cannot build/install prerequisites or start services.

## Immediate Remediation Applied
- Minimal scoped fix to `.gitmodules` adding the two missing submodule entries.

## Required Decision / Next Action
- Install and provision required local dependencies for this runner:
  - Go (>= version expected by repo; geth install path currently expects Go toolchain)
  - Docker + Docker Compose (required by devnet services)
- After environment provisioning, re-run full checklist including 2 successful `make devnet-up` cycles for reproducibility.

---

## Re-check: 2026-02-23 15:13 KST
- Branch: `feat/fault-proof`
- Commit: `58e950300cdf3c23e2d7307122dfa45bda742524`
- Result: **FAIL (same hard blocker)**

### Evidence
- `git pull --rebase` → `Already up to date.`
- `make devnet-up` now passes submodule init step, then fails in `pre-devnet`:
  - `./ops/scripts/geth-version-checker.sh: line 7: geth: command not found`
  - `/bin/sh: 4: go: not found`
  - `env: ‘go’: No such file or directory`

### Delta vs previous run
- No new code-level failure discovered after `.gitmodules` repair.
- Blocker remains environment provisioning (Go/Docker missing), preventing challenger runtime and reproducibility validation.
