# ForceWithdraw Shutdown E2E Test Guide

This document describes the E2E test procedure based on the consolidated shutdown workflow (`Prepare` -> `Execute`).

## 5-Step Shutdown Process

The current codebase has been optimized from individual script execution to a **consolidated step-by-step workflow**.

### Step 1: Block Deposits and Withdrawals (Block)
- **Script**: `BlockDepositsWithdrawals.s.sol`
- **Purpose**: Halt asset movement between L1 and L2 to ensure data consistency.
- **Outcome**: `SuperchainConfig` is set to the Paused state.

### Step 2: Data Collection (Fetch)
- **Scripts**: `fetch_explorer_assets.py`, `compute_l2_burns.py`, `compute_finalized_native_withdrawals.py`
- **Purpose**: Collect L2 asset holder lists and burn history via block explorer APIs and RPC.

### Step 3: Snapshot Generation (Gen)
- **Script**: `GenerateAssetSnapshot.s.sol`
- **Purpose**: Create the final asset ledger (`generate-assets-<chainId>.json`) including cryptographic hashes based on the collected data.

### Step 4: L1 Preparation and Activation (Prepare)
- **Script**: `PrepareL1Withdrawal.s.sol` (Consolidates multiple previous deployment scripts)
- **Key Functions**:
    - Upgrade `ForceWithdrawBridge` and `OptimismPortal`.
    - Deploy `GenFWStorage` and register snapshot hashes.
    - Register the storage address with the bridge and activate Force Withdrawal mode.
- **Critical Output**: The deployed `GenFWStorage` address.

### Step 5: Settlement and Claims (Execute)
- **Script**: `ExecuteL1Withdrawal.s.sol`
- **Key Functions**:
    - Sweep bridge liquidity from the Portal.
    - Execute user-specific claims based on snapshot data.
- **Required Variable**: `STORAGE_ADDRESS` (The address generated in Step 4).

---

## How to Run Tests

### Environment Variable Setup (.env)
The following environment variables must be configured:
```env
BRIDGE_PROXY=0x...
OPTIMISM_PORTAL_PROXY=0x...
SYSTEM_OWNER_SAFE=0x...
PRIVATE_KEY=0x...
L1_RPC_URL=...
L2_RPC_URL=...
```

### Script Execution
```bash
# 1. Simulation (Dry-run)
./scripts/e2e-shutdown-test.sh --dry-run

# 2. Actual Execution (Local/Devnet/Testnet)
./scripts/e2e-shutdown-test.sh
```

## Important Notes
- **Storage Address Passing**: The storage address output from the `PrepareL1Withdrawal` step must be correctly passed as `STORAGE_ADDRESS` for the `ExecuteL1Withdrawal` step (Automated in the `.sh` script).
- **Python Dependencies**: A Python 3 environment with standard libraries (including `urllib3`) is required for Step 2.
- **Safe Compatibility**: `SafeUtils` library supports simulation and transaction generation for multisig accounts (Safe).
