# Shutdown Scripts Guide

This directory contains scripts for L2 shutdown, asset snapshot generation, L1 force withdrawals, and blocking deposits/withdrawals.

## Configuration Files

- `fetch_explorer_assets.py`
  - Collects L2 holder/contract/token and unclaimed withdrawal data based on the Explorer V2 API.
- `GenerateAssetSnapshot.s.sol`
  - Validates L1/L2 data and generates the final snapshot JSON.
- `L1Withdrawal.s.sol`
  - Performs L1 bridge/portal upgrades and activates Force Withdrawal.
- `BlockDepositsWithdrawals.s.sol`
  - Blocks deposits and withdrawals through the `OptimismPortalClosing` upgrade and Superchain pause.
- Documentation
  - `SHUTDOWN_SCRIPT_GUIDE.md`: Summary of features, core usage, and operational procedures.
  - `DRY_RUN_RESULTS.md`: Summary of dry-run execution results.
  - `PHASE2_BURN_RECONCILIATION.md`: Summary of causes and resolutions for Phase 2 burn reconciliation.
  - `USER_WITHDRAWAL_GUIDE.md`: User guide for performing force withdrawals.

## Quick Start

1) L2 Data Collection
```
L1_RPC_URL=... \
L2_RPC_URL=... \
BRIDGE_PROXY=... \
OPTIMISM_PORTAL_PROXY=... \
L1_USDC_BRIDGE_PROXY=... \
L2_USDC_BRIDGE_PROXY=... \
L2_START_BLOCK=0 \
python3 scripts/shutdown/fetch_explorer_assets.py <chainId>
```
Checkpoints:
- Logs indicating completion of holder/contract/token collection.
- Generation of `data/l2-holders-<chainId>.json`.
- Generation of `data/l2-contracts-<chainId>.json`.
- Generation of `data/l2-tokens-<chainId>.json`.
- Generation of `data/unclaimed-withdrawals-<chainId>.json`.

2) Snapshot Generation
```
L1_RPC_URL=... \
L2_RPC_URL=... \
BRIDGE_PROXY=... \
OPTIMISM_PORTAL_PROXY=... \
L1_USDC_BRIDGE_PROXY=... \
L2_USDC_BRIDGE_PROXY=... \
SKIP_FETCH=false \
L2_START_BLOCK=0 \
forge script scripts/shutdown/GenerateAssetSnapshot.s.sol:GenerateAssetSnapshot --fork-url $L2_RPC_URL --ffi
```
Checkpoints:
- Phase 2: All tokens show `[OK] Match`.
- Phase 3: All tokens show `[OK] Match`.
- Generation of `data/generate-assets-<chainId>.json`.

3) Block Deposits and Withdrawals (L1)
```
L1_RPC_URL=... \
PRIVATE_KEY=... \
OPTIMISM_PORTAL_PROXY=... \
SUPERCHAIN_CONFIG_PROXY=... \
PROXY_ADMIN=... \
SYSTEM_OWNER_SAFE=... \
GUARDIAN_SAFE=... \
forge script scripts/shutdown/BlockDepositsWithdrawals.s.sol:BlockDepositsWithdrawals --fork-url $L1_RPC_URL
```
Checkpoints:
- Portal version = `2.8.1-closing`.
- Superchain paused = `Yes`.
- Verification of deposit/receive/onApprove blocking passed.

4) Enable L1 Force Withdrawal
```
L1_RPC_URL=... \
PRIVATE_KEY=... \
BRIDGE_PROXY=... \
PROXY_ADMIN=... \
SYSTEM_OWNER_SAFE=... \
OPTIMISM_PORTAL_PROXY=... \
L1_NATIVE_TOKEN=... \
L1_USDC_BRIDGE_PROXY=... \
DATA_PATH=data/generate-assets-<chainId>.json \
EXECUTE_CLAIMS=false \
forge script scripts/shutdown/L1Withdrawal.s.sol:L1Withdrawal --fork-url $L1_RPC_URL --via-ir
```
Checkpoints:
- `SafeTxHash` logs displayed.
- Force withdrawal mode = `ACTIVE`.