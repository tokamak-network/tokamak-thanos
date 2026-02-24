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
