# Lessons learned

## 2026-02-23 CI drift handling
- Root-cause first: when CI fails across many packages, separate deterministic infra drift (dependency/API mismatches) from product-runtime modules.
- Avoid brittle global gates: root `go mod tidy` and root `golangci-lint ./...` are not safe in mixed-stability monorepos.
- Toolchain alignment matters: if `go.mod` targets `1.24`, lint tooling must also be built/executed with Go 1.24-compatible settings.
- Prefer explicit module matrices in CI so failure domains are visible and independently maintainable.

## Self-rules
- Before changing CI scope, re-run failing commands locally and keep exact failure output.
- If excluding CI targets, keep coverage for active runtime modules and document excluded scope with reason.
- Any CI scope reduction must include a follow-up backlog item to restore coverage after dependency alignment.

## 2026-02-24 Fault-proof devnet checklist verification
- Treat reproducibility as a first-class gate: one successful `devnet-up` is not sufficient; always run a second cycle with explicit `devnet-down -> devnet-up`.
- Distinguish startup-time compose warnings from runtime faults: missing interpolation vars may still be resolved by devnet orchestration, so validate final container env before declaring misconfiguration.
- Non-fatal warning storms can hide in otherwise healthy flows: proposer may continue publishing while emitting repeated estimate-gas warnings, so judge by both throughput and warning profile.

## Self-rules
- For devnet validation, always capture command logs under `/tmp/fault-proof-*.log` with timestamps for reproducible evidence.
- When a service appears healthy, verify both restart count and function signal (block progression, confirmed tx logs), not just container `Up` status.

## 2026-02-24 Fault-proof rerun (12:30 KST)
- A successful `devnet-up` does not imply a clean proposer lane; classify warning storms separately from boot health.
- For repeated gas-estimate failures, track the full cycle: warning burst count, permanent-attempt failure count, and post-failure recovery signal.
- Preflight "clean worktree" must be evaluated at the exact start timestamp of the cycle, not inferred from earlier runs.

## Self-rules
- In checklist reports, always split verdicts into `boot health` and `runtime warning profile` so partial risk is explicit.
- When proposer warnings are present, include at least one recovery proof log (`proposer tx successfully published` after failure) before closing a cycle.

## 2026-02-24 Dispute lifecycle and bond refund validation
- `devnet-down` is stop-only: it does not reset chain state/volumes, so config changes tied to genesis/allocs require `devnet-clean`.
- `faultGameMaxClockDuration` must satisfy `max(clockExtension*2, clockExtension+challengePeriod)`; with `clockExtension=0` and challenge period `120s`, minimum valid value is `120`.
- In manual dispute resolution flows, `resolve` can keep reverting until all required subgame claims are explicitly resolved (including root claim index `0` in this scenario).

## Self-rules
- When validating config changes that affect dispute timers, always prove activation with fresh-state evidence (`devnet-clean` + new game resolution timestamps).
- If `resolve` emits repeated estimate-gas reverts, inspect and complete `resolve-claim` ordering before classifying as infra instability.

## 2026-02-24 Sepolia one-click automation
- Non-interactive deploy scripts should default to deterministic outputs (explicit workdir, generated runtime env path, and chain-id normalization).
- For long automation chains, include a dry-run mode that validates control flow without triggering on-chain writes.
- When runtime keys and deployment roles diverge, startup succeeds but protocol actions fail later; role-address injection from private keys should be explicit in automation.

## Self-rules
- Any deploy automation touching real networks must provide `--dry-run` and explicit required-env validation before execution.
- Generate runtime env from deployment state (`state.json`) rather than hardcoded addresses to avoid stale wiring.

## 2026-02-24 One-click runbook authoring
- Execution runbooks should mirror real script control-flow order (`dry-run -> apply -> artifact generation -> runtime up`) to reduce operator ambiguity.
- Document partial rerun flags (`DEPLOYER_SKIP_APPLY`, `DEPLOYER_FORCE_REINIT`, `RUNTIME_DOWN_FIRST`) explicitly; these are high-impact operational toggles.

## Self-rules
- For every operational script addition, ship a separate "execution guide" doc with required vars, happy path, and rollback/rerun path.
