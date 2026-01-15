# Asset Snapshot Verification Guide

## Overview

This guide explains how to verify that generated asset snapshots match on-chain balances using the Python verification script.

## Prerequisites

- Python 3.7 or higher
- Access to an L2 RPC endpoint
- Generated asset snapshot JSON file

## Installation

### 1. Navigate to the contracts-bedrock directory

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos/packages/tokamak/contracts-bedrock
```

### 2. Set up Python virtual environment (first time only)

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install web3 eth-abi eth-utils
```

## Usage

### Basic Command

```bash
# Activate virtual environment
source .venv/bin/activate

# Run verification
python3 scripts/shutdown/verify_assets.py <snapshot_file> <rpc_url>
```

### Example

```bash
python3 scripts/shutdown/verify_assets.py \
  data/generate-assets-111551119090.json \
  https://rpc.thanos-sepolia.tokamak.network
```

## Output

The script provides detailed verification results:

### Success Output

```
============================================================
       Asset Snapshot Verification Report
============================================================
Target File:  data/generate-assets-111551119090.json
Chain ID:     111551119090
Block Number: 5936060
------------------------------------------------------------

[Token 1] Tokamak Network Token
  L1: 0x0000000000000000000000000000000000000000
  L2: 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000
  Claims: 29
    [OK] 0xF25c18d9e7Da113CFF654546DD7f2b6f7Bb3fC88 | 808563999762483883112308
    [OK] 0x961b6fb7D210298B88d7E4491E907cf09c9cD61d | 73403998
    ...

============================================================
              Verification Summary
============================================================
Total Tokens:          5
Total Claims:          58
Total Amount:          893453000815169815564064615281482 wei
Balance Mismatches:    0
Hash Mismatches:       0
Balance Query Failures: 0

[SUCCESS] All verifications passed!
============================================================
```

### Status Indicators

- `[OK]` - Balance matches snapshot ✅
- `[MISMATCH]` - Balance differs from snapshot ⚠️
- `[WARN]` - Hash mismatch or balance query failed ⚠️

## Verification Process

The script performs three types of verification:

1. **Hash Verification**
   - Recalculates `keccak256(abi.encodePacked(l1Token, claimer, amount))`
   - Compares with stored hash in snapshot

2. **Balance Query**
   - Queries actual on-chain balance for each claimer
   - Supports both native tokens and ERC20 tokens

3. **Amount Comparison**
   - Compares snapshot amount with actual balance
   - Reports any discrepancies

## Exit Codes

- `0` - All verifications passed
- `1` - Verification failures detected

## Troubleshooting

### Connection Error

```
Error: Could not connect to RPC: <url>
```

**Solution**: Check RPC URL and network connectivity

### File Not Found

```
Error: Snapshot file not found: <path>
```

**Solution**: Verify the snapshot file path is correct

### Invalid JSON

```
Error: Invalid JSON in snapshot file
```

**Solution**: Check JSON file format using `jq` or a JSON validator:

```bash
jq . data/generate-assets-111551119090.json
```

## Integration with Shutdown Flow

This verification script is part of the L2 shutdown process:

1. **Generate Snapshot** - `GenerateAssetSnapshot.s.sol` creates the snapshot
2. **Verify Snapshot** - `verify_assets.py` validates the snapshot ← You are here
3. **Deploy Storage** - Deploy GenFWStorage contracts
4. **Register Positions** - Register storage contracts with bridge
5. **Activate Force Withdraw** - Enable force withdrawal mechanism

## Advanced Usage

### Verify Multiple Snapshots

```bash
for snapshot in data/generate-assets-*.json; do
  echo "Verifying $snapshot..."
  python3 scripts/shutdown/verify_assets.py "$snapshot" https://rpc.thanos-sepolia.tokamak.network
done
```

### Save Verification Report

```bash
python3 scripts/shutdown/verify_assets.py \
  data/generate-assets-111551119090.json \
  https://rpc.thanos-sepolia.tokamak.network \
  | tee verification-report.txt
```

## Script Location

- **Script**: `scripts/shutdown/verify_assets.py`
- **Snapshots**: `data/generate-assets-*.json`
- **Virtual Environment**: `.venv/`

## Related Documentation

- [Shutdown Flow Guide](../../docs/llm/shutdown-flow.md)
- [GenerateAssetSnapshot Script](./GenerateAssetSnapshot.s.sol)
- [E2E Shutdown Test](../e2e-shutdown-test_devnet.sh)
