# Shutdown E2E Test Plan (Pre-Implementation)

This document describes the intended end-to-end test flow using the Forge shutdown scripts.
It focuses on sequencing, required inputs, and verification checkpoints before implementation.

## Required Order

The E2E test must run in this order:
1) Block deposits and withdrawals
2) Collect asset data
3) Generate asset snapshot (off-chain)
4) Generate asset snapshot (on-chain registration)
5) Bridge upgrades and activation
6) Final withdrawals and claims

## Contract Location Check

- `ForceWithdrawBridge` is implemented in `src/shutdown/ForceWithdrawBridge.sol`.
- `GenFWStorage` lives in `src/shutdown/GenFWStorage.sol` and is used by the scripts for snapshot registration.
- No missing contracts were found for the current shutdown script imports.

## Preconditions

- L1/L2 RPC URLs are available and synced.
- Required addresses are known (bridge proxy, portal proxy, proxy admin, safes).
- Execution key or Safe flow is configured.
- `data/` directory is writable and accessible.

## Step-by-Step Plan

### Step 1: Block Deposits and Withdrawals (L1)
Script: `scripts/shutdown/BlockDepositsWithdrawals.s.sol:BlockDepositsWithdrawals`

Inputs:
- `L1_RPC_URL`, `PRIVATE_KEY`
- `OPTIMISM_PORTAL_PROXY`, `SUPERCHAIN_CONFIG_PROXY`
- `PROXY_ADMIN`, `SYSTEM_OWNER_SAFE`, `GUARDIAN_SAFE`

Expected Results:
- OptimismPortal implementation upgraded to `OptimismPortalClosing`.
- Portal `version()` returns `2.8.1-closing`.
- Deposit/receive/onApprove calls revert with shutdown message.
- Superchain paused status is `true`.

### Step 2: Collect Asset Data (L2)
Scripts:
- `scripts/shutdown/fetch_explorer_assets.py`
- `scripts/shutdown/compute_l2_burns.py`
- `scripts/shutdown/compute_finalized_native_withdrawals.py`

Inputs:
- `L1_RPC_URL`, `L2_RPC_URL`, `L2_START_BLOCK`
- `BRIDGE_PROXY`, `OPTIMISM_PORTAL_PROXY`, `L1_USDC_BRIDGE_PROXY`, `L2_USDC_BRIDGE_PROXY`

Outputs:
- `data/l2-holders-<chainId>.json`
- `data/l2-contracts-<chainId>.json`
- `data/l2-tokens-<chainId>.json`
- `data/unclaimed-withdrawals-<chainId>.json`
- `data/l2-burns-<chainId>.json`

### Step 3: Generate Asset Snapshot (Off-Chain)
Script: `scripts/shutdown/GenerateAssetSnapshot.s.sol:GenerateAssetSnapshot`

Inputs:
- Same RPC and proxy inputs as Step 2
- `SKIP_FETCH=false`

Output:
- `data/generate-assets-<chainId>.json`

Checkpoints:
- Phase 2/3 validations show `[OK] Match` for tokens.
- Snapshot file exists and contains holders/tokens/unclaimed data.

### Step 4: Generate Asset Snapshot (On-Chain Registration)
Script: `scripts/shutdown/PrepareL1Withdrawal.s.sol:PrepareL1Withdrawal`

Purpose:
- Deploy `GenFWStorage` and register snapshot hashes on L1.

Inputs:
- `L1_RPC_URL`, `PRIVATE_KEY`
- `BRIDGE_PROXY`, `PROXY_ADMIN`, `SYSTEM_OWNER_SAFE`
- `OPTIMISM_PORTAL_PROXY`, `L1_USDC_BRIDGE_PROXY`, `L1_USDC_BRIDGE_ADMIN`
- `DATA_PATH=data/generate-assets-<chainId>.json`

Checkpoints:
- `GenFWStorage` deployed and storage address logged.
- Registered position is active on `ForceWithdrawBridge`.

### Step 5: Bridge Upgrades and Activation
Scripts:
- `PrepareL1Withdrawal.s.sol` (if not already run)
- `L1Withdrawal.s.sol` (combined flow alternative)

Upgrades:
- `ForceWithdrawBridge` implementation for L1 bridge
- `ShutdownOptimismPortal` or `ShutdownOptimismPortal2` for portal
- `ShutdownL1UsdcBridge` for L1 USDC bridge

Checkpoints:
- Force withdrawal mode is `ACTIVE`.
- Proxy admin ownership checks pass (EOA or Safe).

### Step 6: Final Withdrawals and Claims
Script: `scripts/shutdown/ExecuteL1Withdrawal.s.sol:ExecuteL1Withdrawal`

Inputs:
- Same L1 proxy/admin inputs as Step 4
- `DATA_PATH=data/generate-assets-<chainId>.json`
- `EXECUTE_CLAIMS=true`

Checkpoints:
- Sweep/claim transactions succeed.
- `claimState` for processed hashes is `true`.
- Recipient balances updated for ETH/ERC20.

## E2E Validation Matrix

- Deposits blocked: `depositTransaction`, `receive`, `onApprove` revert.
- Snapshot integrity: computed totals match for all tokens.
- Registration: `position(storage)` is true; `active()` is true.
- Claims: balances and `claimState` updated as expected.

## Notes

- For Safe-based execution, validate `SafeTxHash` output at each step.
- If re-running, ensure idempotent checks (skip if already upgraded/active).
- The on-chain snapshot registration is separate from the off-chain snapshot file generation.
