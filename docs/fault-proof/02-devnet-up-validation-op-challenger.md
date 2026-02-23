# Devnet-up Validation Plan (op-challenger upgrade)

## Objective
Validate that the `feat/fault-proof` branch can bring up local devnet via `make devnet-up` and that the upgraded `op-challenger` runs stably with expected fault-proof behavior.

## Scope
- Branch: `feat/fault-proof`
- Focus: `op-challenger` upgrade impact on local devnet boot + runtime behavior
- Environment: local developer machine

## Go / No-Go Criteria
A run is **GO** only if all conditions pass:
1. `make devnet-up` starts devnet successfully (exit code + service health OK)
2. `op-challenger` starts and remains healthy (no crash loop/panic)
3. No critical regressions in core flow (L1/L2 progression, batch/proposer normal)
4. Fault-proof observation points pass (game polling/connection behavior as expected)
5. Run is reproducible at least 2 times

## Execution Checklist

### 1) Preflight (Reproducibility)
- [ ] Confirm branch and commit hash are fixed for the run
- [ ] Confirm clean working tree
- [ ] Snapshot key env/config used for devnet-up
- [ ] Record local dependency/tool versions (Go, Docker, etc.)

### 2) Build & Compatibility Gate
- [ ] Build targets used by devnet-up complete without compile/link errors
- [ ] Verify `op-challenger` binary/image path used by this branch
- [ ] Verify upgraded flags/env vars (added/removed/renamed)
- [ ] Verify no unresolved references after upgrade

### 3) Boot Validation (`make devnet-up`)
- [ ] Execute `make devnet-up`
- [ ] Capture startup logs for all core services
- [ ] Confirm expected services are up (L1, L2, batcher, proposer, challenger)
- [ ] Observe first 5-10 min logs for panic/restart storms/retry explosions

### 4) Fault-Proof Runtime Validation (op-challenger)
- [ ] Confirm challenger can connect to configured RPC/contracts
- [ ] Confirm challenger polling loop behaves normally
- [ ] Validate baseline scenario: no active dispute, no abnormal errors
- [ ] (Optional) Trigger dispute-like condition in local setup and verify behavior/logging

### 5) Regression Gate
- [ ] L2 block progression continues as expected
- [ ] Batcher/proposer continue normal operation
- [ ] No critical increase in error/warn profile compared to previous baseline
- [ ] Key RPC endpoints respond reliably during the run

### 6) Reproducibility Re-run
- [ ] Tear down environment
- [ ] Re-run `make devnet-up` with same config
- [ ] Confirm same success criteria pass again

## Evidence to Collect
- Commit hash used
- Exact commands executed
- Service status snapshots
- Challenger log excerpts (startup, steady-state, and any errors)
- Final pass/fail matrix for each checklist section

## Failure Handling
If any check fails, report:
1. Failing step
2. Symptom/log signature
3. Suspected root cause
4. Immediate remediation action
5. Re-test result

## Deliverable Format
For each validation cycle, produce a short report with:
- **Result:** PASS / FAIL
- **Environment + commit**
- **What passed**
- **What failed / risks**
- **Recommendation:** merge-forward / fix-before-merge
