# L2 Shutdown Scripts Analysis

## Overview

The `scripts/shutdown` directory contains a collection of scripts designed to securely facilitate "Force Withdrawals" of user assets from L2 to L1 during a network shutdown.

### Core Objectives

1. Comprehensive L2 asset data collection and rigorous validation.
2. Generation of an asset snapshot ledger based on cryptographic proofs.
3. L1 Bridge upgrades and activation of force withdrawal permissions.
4. Enabling users to directly claim their assets on L1 via the bridge.

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          L2 Shutdown Flow (V2)                              │
│                                                                             │
│  [Phase 1: Data Collection] (Python)                                        │
│  ┌──────────────────────┐    ┌──────────────────────────────────┐           │
│  │ fetch_explorer_      │───>│ l2-holders / l2-tokens /         │           │
│  │ assets.py            │    │ unclaimed-withdrawals (JSON)     │           │
│  └──────────────────────┘    └──────────────────────────────────┘           │
│  ┌──────────────────────┐    ┌──────────────────────┐                      │
│  │ compute_l2_burns.py  │───>│ l2-burns-{id}.json   │                      │
│  └──────────────────────┘    └──────────────────────┘                      │
│                                                                             │
│  [Phase 2: Snapshot Generation] (Solidity - L2 Fork)                        │
│  ┌───────────────────────────┐      ┌───────────────────────────┐           │
│  │ GenerateAssetSnapshot.s.sol│ ───> │ generate-assets-{id}.json │           │
│  │ (Validation & Hashing)    │      │ (Merkle-like Claim Proofs)│           │
│  └───────────────────────────┘      └───────────────────────────┘           │
│                                                                             │
│  [Phase 3: Preparation & Infrastructure] (Solidity - L1)                    │
│  ┌───────────────────────────┐      ┌───────────────────────────┐           │
│  │ PrepareL1Withdrawal.s.sol │ ───> │ 1. Upgrade Proxies        │           │
│  │ (Steps 1-8)               │      │ 2. Deploy GenFWStorage    │           │
│  │                           │      │ 3. Activate Force Mode    │           │
│  └───────────────────────────┘      └───────────────────────────┘           │
│                                                                             │
│  [Phase 4: Execution & Settlement] (Solidity - L1)                          │
│  ┌───────────────────────────┐      ┌───────────────────────────┐           │
│  │ ExecuteL1Withdrawal.s.sol │ ───> │ 1. Sweep Portal Liquidity │           │
│  │ (Steps 9-10)              │      │ 2. Admin/Test Claims      │           │
│  └───────────────────────────┘      └───────────────────────────┘           │
│                                                                             │
│  [Phase 5: User Action] (L1 Etherscan/CLI)                                  │
│  ┌───────────────────────────┐      ┌───────────────────────────┐           │
│  │ forceWithdrawClaim()      │ ───> │ User receives L1 assets   │           │
│  └───────────────────────────┘      └───────────────────────────┘           │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## File Structure & Roles

### Python Scripts (Off-chain Processing)

| File | Role | Input | Output |
| :--- | :--- | :--- | :--- |
| `fetch_explorer_assets.py` | Collects all balance data from Explorer V2 API | `chainId` | `l2-holders`, `l2-tokens`, `l2-contracts` |
| `compute_l2_burns.py` | Calculates tokens burned without withdrawal (fees, etc.) | `L2_RPC_URL` | `l2-burns-{id}.json` |
| `compute_finalized_native_withdrawals.py` | Aggregates finalized native token withdrawals on L1 | `L1_RPC_URL` | Stdout (Used for validation) |

### Solidity Scripts (On-chain Logic)

| File | Role | Target Network |
| :--- | :--- | :--- |
| `GenerateAssetSnapshot.s.sol` | Validates L1-L2 consistency & generates claim ledger | L2 (Fork) |
| `BlockDepositsWithdrawals.s.sol` | Blocks new bridge deposits & pauses Superchain | L1 |
| `PrepareL1Withdrawal.s.sol` | Upgrades, Storage deployment, and Mode activation | L1 |
| `ExecuteL1Withdrawal.s.sol` | Liquidity sweep from Portal & Settlement (Claims) | L1 |

---

## Core Mechanism Analysis

### 1. Asset Consistency Validation (Phase 2)
`GenerateAssetSnapshot` does more than just collect balances; it mathematically verifies that assets held by the L1 Bridge match the L2 supply.

*   **Native Token (TON) Validation**:
    `expectedL1 = L2TotalSupply + UnclaimedWithdrawals + BurnAdjustment`
*   **ERC20/ETH Validation**:
    `L1Bridge.deposits(l1Token, l2Token) == L2Token.totalSupply() + burns`

### 2. Multisig Nonce Management (`SafeNonceOffset`)
`PrepareL1Withdrawal` sends multiple transactions to a Gnosis Safe. To queue multiple transactions at once, use `safeNonceOffset`.
*   **Usage**: Start from 0 for the first batch, then set the offset to the number of transactions created in previous runs for subsequent batches.

### 3. Claim Hash Generation & Registration
To protect user assets, each entry is hashed as follows:
*   **Hash**: `keccak256(abi.encodePacked(l1Token, claimer, amount))`
*   **Selector**: `bytes4(keccak256(abi.encodePacked("_", stripped_hash, "()")))`
These hash-selector pairs are registered in the `GenFWStorage` contract as proof for bridge withdrawal authorization.

### 4. Closer and Permission Management
*   **Closer**: An account authorized to finalize the shutdown procedure. Configured during the `Prepare` phase.
*   **forceActive**: Disables standard bridge mode and enables Force Withdrawal mode.
*   **forceRegistry**: Officially registers the generated snapshot (Storage) as the bridge's authorized ledger.

---

## Environment Variables (Standard V2)

### Common Required Variables
| Variable | Description |
| :--- | :--- |
| `L1_RPC_URL` / `L2_RPC_URL` | RPC Endpoints |
| `PRIVATE_KEY` | Deployer/Executor Private Key |
| `BRIDGE_PROXY` | L1StandardBridge Proxy Address |
| `DATA_PATH` | Path to the asset snapshot JSON file |

### Operational Variables (Safe Environment)
| Variable | Description |
| :--- | :--- |
| `SYSTEM_OWNER_SAFE` | Address of the Safe that owns ProxyAdmin and Bridge |
| `SAFE_NONCE_OFFSET` | Nonce offset for continuous Safe transaction generation (Default: 0) |
| `STORAGE_ADDRESS` | GenFWStorage address deployed in `Prepare` (Required for `Execute`) |

---

## Execution Guidelines

1.  **Enable FFI**: `GenerateAssetSnapshot` calls Python scripts, so `ffi = true` must be set in `foundry.toml`.
2.  **Native Token Mapping**: The Native Token (TON) on Thanos is mapped to a specifically designated L1 TON address.
3.  **Manual Claims**: When users claim manually, they must input the hash value **without** the `0x` prefix.
4.  **Transaction Simulation**: All scripts should be reviewed via Forge Simulation before actual broadcasting to avoid unintended side effects.

---
**Version Info**: As of 2026.01.27 (Reflecting consolidated workflow)
