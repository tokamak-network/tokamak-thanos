# Quickstart: BlockDepositsWithdrawals Script

This document summarizes the procedure for blocking OptimismPortal deposits and pausing SuperchainConfig using the `scripts/shutdown/BlockDepositsWithdrawals.s.sol` script.

## Prerequisites

- `ProxyAdmin.owner()` must be set to `SystemOwnerSafe`.
- `SuperchainConfig.guardian()` must be an **EOA** (an address signing with a private key).
- The `PRIVATE_KEY` must satisfy the following two conditions:
  - It must be the account of the sole owner (1-of-1) of `SystemOwnerSafe`.
  - It must be the same address as `SuperchainConfig.guardian()`.

## Environment Variables

The script reads the following values from `packages/tokamak/contracts-bedrock/.env`:

- `PRIVATE_KEY`: Private key of the guardian EOA and the SystemOwnerSafe owner.
- `L1_RPC_URL`: L1 RPC URL (used for `--fork-url` during execution).
- `OPTIMISM_PORTAL_PROXY`: Proxy address of OptimismPortal.
- `SUPERCHAIN_CONFIG_PROXY`: Proxy address of SuperchainConfig.
- `PROXY_ADMIN`: Address of ProxyAdmin.
- `SYSTEM_OWNER_SAFE`: Address of SystemOwnerSafe.
- `GUARDIAN_SAFE`: (Optional) If the guardian is an EOA, set this to the **guardian's EOA address**.

If `GUARDIAN_SAFE` is not specified, it defaults to `SYSTEM_OWNER_SAFE`. If the guardian is an EOA, you must specify the EOA address.

## Dry-run

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos/packages/tokamak/contracts-bedrock
set -a; source .env; set +a
forge script scripts/shutdown/BlockDepositsWithdrawals.s.sol --fork-url "$L1_RPC_URL" --sig "run()"
```

Items to check in the dry-run log:

- `Derived caller from PRIVATE_KEY` matches the guardian address.
- `Is caller owner: true`
- `Deposit block reason for depositTransaction/receive/onApprove` all show the same blocking reason.

## Actual Execution

If the dry-run is successful, add `--broadcast`.

```bash
forge script scripts/shutdown/BlockDepositsWithdrawals.s.sol --fork-url "$L1_RPC_URL" --sig "run()" --broadcast
```

## Verification

After execution, the script verifies the following:

- OptimismPortal implementation version is `2.8.1-closing`.
- SuperchainConfig is in a `paused` state.
- `depositTransaction`, `receive`, and `onApprove` are all blocked with the same revert reason.

## Troubleshooting

- If `Is caller owner: false` appears, the `PRIVATE_KEY` does not match the owner of `SYSTEM_OWNER_SAFE`.
- If `GUARDIAN_SAFE is not a contract` appears, the guardian is an EOA, and you must specify the EOA address in `GUARDIAN_SAFE`.
