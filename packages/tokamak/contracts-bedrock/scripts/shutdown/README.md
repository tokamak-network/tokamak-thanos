# L2 Shutdown Script Guide

This directory contains scripts for generating snapshots, executing L1 force withdrawals, and blocking deposits/withdrawals for L2 shutdown.

## Configuration Files

- `fetch_explorer_assets.py`
  - Collects L2 holders/contracts/tokens and unclaimed withdrawal data based on the Explorer V2 API.
- `compute_finalized_native_withdrawals.py`
  - Calculates the sum of `NativeTokenWithdrawalFinalized` on L1.
- `compute_l2_burns.py`
  - Calculates the sum of L2 burns/withdrawals and generates `data/l2-burns-<chainId>.json`.
- `GenerateAssetSnapshot.s.sol`
  - Validates L1/L2 data and generates the final snapshot JSON.
- `L1Withdrawal.s.sol`
  - Performs L1 Bridge/Portal upgrades and activates force withdrawals.
- `BlockDepositsWithdrawals.s.sol`
  - Blocks deposits and withdrawals via `OptimismPortalClosing` upgrade and Superchain pause.
- Documentation
  - `USER_WITHDRAWAL_GUIDE.md`: Guide for user force withdrawals.

## Output Files

- `data/l2-holders-<chainId>.json`: List of L2 holders (accounts + token holders)
- `data/l2-contracts-<chainId>.json`: List of L2 contracts
- `data/l2-tokens-<chainId>.json`: List of L2 tokens
- `data/unclaimed-withdrawals-<chainId>.json`: List of unclaimed withdrawals
- `data/l2-burns-<chainId>.json`: Sum of L2 burns/withdrawals
- `data/generate-assets-<chainId>.json`: Final asset snapshot

## Pre-execution (Recommended)

1) Collect Explorer Data
```bash
L1_RPC_URL=... \
L2_RPC_URL=... \
L2_START_BLOCK=0 \
python3 scripts/shutdown/fetch_explorer_assets.py <chainId>
```

2) Generate L2 Burn Sum
```bash
L2_RPC_URL=... \
python3 scripts/shutdown/compute_l2_burns.py $L2_RPC_URL <chainId> $L2_START_BLOCK
```

3) Calculate L1 Finalized Native Withdrawal Sum
```bash
L1_RPC_URL=... \
python3 scripts/shutdown/compute_finalized_native_withdrawals.py $L1_RPC_URL <bridgeAddress>
```

**Checkpoints:**
- Result files from the 3 scripts above are generated in `data/`.
- The snapshot script will fail if `unclaimed-withdrawals-<chainId>.json` is missing.
- If `unclaimed-withdrawals-<chainId>.json` is `[]`, it's treated as 0 unclaimed withdrawals.
- If `L2_START_BLOCK` is not specified, it defaults to `0`.

## Quick Start

1) Collect L2 Data
```bash
L1_RPC_URL=... \
L2_RPC_URL=... \
BRIDGE_PROXY=... \
OPTIMISM_PORTAL_PROXY=... \
L1_USDC_BRIDGE_PROXY=... \
L2_USDC_BRIDGE_PROXY=... \
L2_START_BLOCK=0 \
python3 scripts/shutdown/fetch_explorer_assets.py <chainId>
```
**Checkpoints:**
- Look for logs confirming completion of holder/contract/token collection.
- Verify generation of `data/l2-holders-<chainId>.json`, `l2-contracts-<chainId>.json`, `l2-tokens-<chainId>.json`, and `unclaimed-withdrawals-<chainId>.json`.

2) Generate Snapshot
```bash
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
**Checkpoints:**
- Ensure `[OK] Match` for all tokens in Phase 2 and Phase 3.
- Verify generation of `data/generate-assets-<chainId>.json`.
- If `SKIP_FETCH=true`, all required files must already exist in `data/`.

3) Block Deposits and Withdrawals (L1)
```bash
L1_RPC_URL=... \
PRIVATE_KEY=... \
OPTIMISM_PORTAL_PROXY=... \
SUPERCHAIN_CONFIG_PROXY=... \
PROXY_ADMIN=... \
SYSTEM_OWNER_SAFE=... \
GUARDIAN_SAFE=... \
forge script scripts/shutdown/BlockDepositsWithdrawals.s.sol:BlockDepositsWithdrawals --fork-url $L1_RPC_URL
```
**Checkpoints:**
- Portal version should be `2.8.1-closing`.
- Superchain pause status should be `Yes`.
- Ensure blocking validation for deposit/receive/onApprove passes.

4) Activate L1 Force Withdrawals
```bash
L1_RPC_URL=... \
PRIVATE_KEY=... \
BRIDGE_PROXY=... \
PROXY_ADMIN=... \
SYSTEM_OWNER_SAFE=... \
OPTIMISM_PORTAL_PROXY=... \
L1_NATIVE_TOKEN=... \
L1_USDC_BRIDGE_PROXY=... \
L1_USDC_BRIDGE_ADMIN=... \
DATA_PATH=data/generate-assets-<chainId>.json \
EXECUTE_CLAIMS=false \
forge script scripts/shutdown/L1Withdrawal.s.sol:L1Withdrawal --fork-url $L1_RPC_URL --via-ir
```
**Checkpoints:**
- Verify `SafeTxHash` log output.
- Force withdrawal mode should be `ACTIVE`.
- If `L1_USDC_BRIDGE_ADMIN` is an EOA, the bridge `upgradeTo` is executed directly by the EOA.

## Modular Execution Options

- `PrepareL1Withdrawal.s.sol`: Executes steps 1-7 (Upgrade/Registration/Activation).
- `ExecuteL1Withdrawal.s.sol`: Executes steps 8-9 (Sweep/Claims).

**Checkpoints:**
- For paths via Safe, verify the `SafeTxHash` logs.
- If `EXECUTE_CLAIMS=false`, claim execution is skipped.

## Precautions

- Snapshot generation will fail if `unclaimed-withdrawals-<chainId>.json` is missing.
- If `L1_USDC_BRIDGE_ADMIN` is an EOA, `upgradeTo` will be executed directly via the EOA instead of the Safe.
