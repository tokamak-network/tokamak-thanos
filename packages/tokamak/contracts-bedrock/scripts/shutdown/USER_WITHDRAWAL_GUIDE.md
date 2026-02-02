# User Self-Withdrawal Guide

This guide explains how **regular users** can directly withdraw their L2 assets from L1 (Ethereum) after the L2 network shutdown.

## Overview

Once the administrator completes the L2 shutdown and registers the asset snapshot, users can withdraw their assets by calling the `forceWithdrawClaim` function of the L1 bridge. **No special permissions or administrator approval are required**, and users receive their assets by executing the transaction directly from their own wallets.

## Withdrawal Permission Verification

**Anyone can withdraw.** The `forceWithdrawClaim` function of the `ForceWithdrawBridge` contract has the following characteristics:
- **External**: It is a function that can be called by anyone from the outside.
- **No Permission Constraints**: There are no administrator-only restrictions like `onlyOwner`.
- **Hash Verification**: Security is guaranteed through cryptographic hash verification (`keccak256`) rather than permissions. In other words, the withdrawal succeeds only when your own asset information recorded in the snapshot is accurate.

```solidity
function forceWithdrawClaim(
  address _position,
  string calldata _hash,
  address _token,
  uint256 _amount,
  address _address
) external {  // ✅ Can be called by anyone
  claim(_position, _hash, _token, _amount, _address);
}
```

---

## Prerequisites

Before starting the withdrawal, please ensure that the following items are completed:
1. **L2 Shutdown Complete**: The L2 network operator must have completed the shutdown procedure.
2. **Snapshot Registered**: The administrator must have registered the asset information (snapshot) in the `GenFWStorage` contract on L1.
3. **Force Withdrawal Activated**: The `active` state of the bridge contract must be `true`.

Please contact the operator to confirm if **"Force Withdrawal is ready"**.

---

## Step-by-Step Withdrawal Method

### Step 1: Verify Your Claim Data

The data required for withdrawal can be found in the **L2 Asset Snapshot JSON file** provided by the service operator.

**Snapshot Data Example:**
```json
{
  "l1Token": "0xa30fe40285B8f5c0457DbC3B7C8A280373c40044",
  "l2Token": "0x0000000000000000000000000000000000000000",
  "tokenName": "Tokamak Network",
  "data": [
    {
      "claimer": "0x049bF8C1291938Ae5A8CBB109062A91af3a153E5",
      "amount": "19999999888886994463890645985991",
      "hash": "0x3e92bf2ea2aa02b9c1d1b023e3d7a1d5930aaacd67cbff6ddd6788bcccbe7cd4"
    }
  ]
}
```

**Information You Must Note Down:**
- `l1Token`: The L1 token address to receive (ETH is `0x0000...0000`)
- `claimer`: Your wallet address (recipient)
- `amount`: The exact amount available for withdrawal (in Wei)
- `hash`: The unique hash value assigned to you

### Step 2: Verify Required Contract Addresses

You must receive the following two addresses from the operator:
- **Bridge Proxy Address**: The target for calling the withdrawal function (e.g., `0x072B...`)
- **Storage Position Address**: The `GenFWStorage` address where your asset information is stored

---

### Step 3: Execute Withdrawal (Etherscan Recommended)

The easiest and safest way for users to withdraw without professional tools is to use Etherscan.

1. **Visit Etherscan**: Search for the Bridge Proxy address on Etherscan.
2. **Check the Contract Tab**: Click the `Contract` menu at the bottom.
3. **Click Write as Proxy**:
   - Since the bridge is a proxy contract, you must select **`Write as Proxy`** instead of `Write Contract`.
   - If a message appears, go through the `Is this a proxy?` -> `Verify` process to activate this tab.
4. **Connect Wallet**: Click the `Connect to Web3` button to connect your wallet (MetaMask, etc.).
5. **Select the forceWithdrawClaim Function**: Locate and click `forceWithdrawClaim` in the list.
6. **Input Parameters (Very Important)**:
   - `_position`: Input the verified **Storage Position address**.
   - `_hash`: Input the `hash` value from the snapshot, but **you must remove the leading `0x`**. (e.g., `0x3e92...` -> `3e92...`)
   - `_token`: Input the **l1Token address** from the snapshot.
   - `_amount`: Input the total **amount value** (all numbers) from the snapshot.
   - `_address`: Input **your wallet address** where you wish to receive the assets.
7. **Execute Write**: Click the button, verify the gas fee in your wallet, and provide final approval.

---

## Important Notes

### ⚠️ Note on Hash Format
When inputting the `_hash` parameter on Etherscan or in scripts, **there should be no `0x` prefix**.
If you include `0x`, the contract will fail to find the information, causing the transaction to revert.

### Withdrawing Multiple Types of Assets
If a user holds both TON and ETH, they must execute `forceWithdrawClaim` **twice in total**, once for each asset type (calling it once per asset type).

### Checking Withdrawal Status in Advance
Before attempting a withdrawal, you can check if someone else has withdrawn it or if it's already complete.
- Go to the **Read as Proxy** tab on Etherscan.
- Query the **claimState** function by inputting your claim hash (including the 0x).
- `false`: Not yet withdrawn (Available for withdrawal)
- `true`: Already withdrawn

---

## Troubleshooting

### Error: "Invalid hash"
- Double-check if you removed the `0x` from the input hash.
- This error occurs if `_token`, `_amount`, or `_address` differs from the snapshot. You must enter them exactly, down to the last digit.

### Error: "Already claimed"
- This asset has already been withdrawn. Please check your wallet balance.

### Error: "Position not registered"
- Ensure the address entered in `_position` matches the official storage address provided by the operator.

### Transaction Fee (Gas Fee)
Withdrawal is an action that occurs on L1 (Mainnet), so a small amount of gas (ETH) is required. Typically, it consumes about 150,000 to 200,000 gas.

---

## Security Notice
- **Verify Official Sources**: Always verify contract addresses and snapshot files through the official Discord or website.
- **No Time Limit**: There is no deadline for withdrawing registered assets. Users can withdraw at their convenience.
- **Private Key Security**: Any site that asks for your private key during this process is a 100% scam. Only use the official Etherscan website.

---

## Developer Tools Usage (Reference)

### Using Cast (Foundry CLI)
```bash
cast send $BRIDGE_PROXY \
  "forceWithdrawClaim(address,string,address,uint256,address)" \
  $STORAGE_POSITION \
  "3e92bf2ea2aa02b9c1d1b023e3d7a1d5930aaacd67cbff6ddd6788bcccbe7cd4" \
  $L1_TOKEN \
  $AMOUNT \
  $MY_ADDRESS \
  --rpc-url $L1_RPC_URL \
  --private-key $YOUR_PRIVATE_KEY
```

### Using Ethers.js
```javascript
const bridge = new ethers.Contract(BRIDGE_PROXY, abi, signer);
await bridge.forceWithdrawClaim(
  STORAGE_POSITION,
  "3e92bf2ea2aa02b9c1d1b023e3d7a1d5930aaacd67cbff6ddd6788bcccbe7cd4", // Excluding 0x
  L1_TOKEN,
  AMOUNT,
  MY_ADDRESS
);
```
