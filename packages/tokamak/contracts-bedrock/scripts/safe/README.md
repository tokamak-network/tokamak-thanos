# Safe Transaction Scripts

This directory contains scripts to create, sign, and execute Gnosis Safe multisig transactions.

## 🚀 Key Features

- **`signTransaction.ts`**: Creates a transaction and signs it using private keys from `.env`.
- **`executeTransaction.ts`**: Executes the transaction on-chain using collected signatures.

## 📋 Setup

1. Create `.env` file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` file:
   - `SAFE_ADDRESS`: Safe contract address
   - `L1_RPC_URL`: RPC endpoint
   - `PRIVATE_KEY_OWNER_*`: Private keys of owners to sign (multiple keys supported)

### Transaction Configuration (Choose one)

#### Option A: Change ProxyAdmin (Simple Mode)
Use this to call `changeProxyAdmin` on the `ProxyAdmin` contract.
```bash
PROXY_ADMIN=0x...
TARGET_PROXY=0x...
NEW_ADMIN=0x...
```

#### Option B: Custom Transaction (Advanced Mode)
Use this to call any function on any contract. If `TX_DATA` is set, it overrides Option A.
```bash
TX_TO=0x...      # Target contract address
TX_DATA=0x...    # Encoded function call data (see guide below)
TX_VALUE=0       # ETH value to send (usually 0)
```

## 💡 TX_DATA Generation Guide

You can easily generate `TX_DATA` using the `cast` tool.

### Example 1: Change Gas Limit (SystemConfig)
Calling `setGasLimit(uint64)` to change the gas limit to 30,000,000:

```bash
# Function signature and arguments
cast calldata "setGasLimit(uint64)" 30000000
# Output: 0x... (Copy this value to TX_DATA)
```

### Example 2: Pause
Calling `pause()` function:

```bash
cast calldata "pause()"
# Output: 0x8456cb59
```

### Example 3: Complex Arguments (Address, etc.)
Calling `transferOwnership(address)`:

```bash
cast calldata "transferOwnership(address)" 0x1234...
```

## 🛠 Usage

### 1. Collect Signatures (Sign)

Generate signatures using private keys in `.env` and save them to `safe_tx_data.json`.

```bash
npx tsx signTransaction.ts
```

- Automatically checks Safe's current Nonce and Threshold.
- Verifies if keys in `.env` belong to actual owners.
- Displays a "Ready to execute" message if enough signatures are collected.

### 2. Add External Signature (Optional)

If another owner has generated a signature using CLI or another tool, you can add it during execution.

### 3. Execute Transaction (Execute)

Execute the transaction on-chain once enough signatures (Threshold) are collected.

```bash
npx tsx executeTransaction.ts [EXTRA_SIGNATURE]
```

- Loads `safe_tx_data.json`.
- Checks if the number of signatures meets the threshold.
- Sorts signatures and calls `execTransaction`.
- If you are missing 1 signature, you can pass it as an argument to execute immediately:
  ```bash
  npx tsx executeTransaction.ts 0x1234...
  ```

## 🔒 Security

- **Never commit `.env` to git.**
- `safe_tx_data.json` does not contain sensitive information (Private Keys). It only contains signature values.
