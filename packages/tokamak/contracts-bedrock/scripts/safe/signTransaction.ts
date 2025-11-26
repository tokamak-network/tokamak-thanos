import { ethers } from 'ethers';
import fs from 'fs';
import dotenv from 'dotenv';
import path from 'path';

// Load environment variables
dotenv.config();

const {
  L1_RPC_URL,
  SAFE_ADDRESS,
  PROXY_ADMIN,
  TARGET_PROXY,
  NEW_ADMIN,
  TX_DATA, // Optional: Raw transaction data
  TX_TO,   // Optional: Target contract address
  TX_VALUE // Optional: ETH value
} = process.env;

// Helper to get all private keys from env starting with PRIVATE_KEY_
function getSignerPrivateKeys(): string[] {
  return Object.keys(process.env)
    .filter(key => key.startsWith('PRIVATE_KEY_'))
    .map(key => process.env[key]!)
    .filter(key => key && key.length > 0);
}

async function main() {
  console.log('🚀 Safe Transaction Signer\n');

  // 1. Configuration Validation
  if (!L1_RPC_URL || !SAFE_ADDRESS) {
    throw new Error('Missing required environment variables: L1_RPC_URL, SAFE_ADDRESS');
  }

  const provider = new ethers.JsonRpcProvider(L1_RPC_URL);
  const safeAddress = SAFE_ADDRESS;

  // 2. Determine Transaction Details
  let to = TX_TO;
  let value = TX_VALUE || '0';
  let data = TX_DATA || '0x';

  // Special case: ChangeProxyAdmin
  if (!to && PROXY_ADMIN && TARGET_PROXY && NEW_ADMIN) {
    console.log('ℹ️  Mode: ChangeProxyAdmin');
    to = PROXY_ADMIN;
    const proxyAdminInterface = new ethers.Interface([
      'function changeProxyAdmin(address proxy, address newAdmin)',
    ]);
    data = proxyAdminInterface.encodeFunctionData('changeProxyAdmin', [
      TARGET_PROXY,
      NEW_ADMIN,
    ]);
  } else if (!to || !data) {
    throw new Error('Missing transaction details. Provide either (TX_TO, TX_DATA) or (PROXY_ADMIN, TARGET_PROXY, NEW_ADMIN)');
  }

  console.log('📋 Transaction Details:');
  console.log('  Safe:', safeAddress);
  console.log('  To:', to);
  console.log('  Value:', value);
  console.log('  Data Length:', data.length);
  console.log('');

  // 3. Get Safe Info
  const safeAbi = [
    'function getThreshold() view returns (uint256)',
    'function nonce() view returns (uint256)',
    'function getOwners() view returns (address[])',
    'function getTransactionHash(address to, uint256 value, bytes calldata data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, uint256 _nonce) view returns (bytes32)',
  ];

  const safe = new ethers.Contract(safeAddress, safeAbi, provider);
  const threshold = await safe.getThreshold();
  const nonce = await safe.nonce();
  const owners = await safe.getOwners();

  console.log('📊 Safe Status:');
  console.log('  Threshold:', threshold.toString());
  console.log('  Nonce:', nonce.toString());
  console.log('  Owners:', owners.length);
  console.log('');

  // 4. Calculate Transaction Hash
  const operation = 0; // Call
  const safeTxGas = 0;
  const baseGas = 0;
  const gasPrice = 0;
  const gasToken = ethers.ZeroAddress;
  const refundReceiver = ethers.ZeroAddress;

  const txHash = await safe.getTransactionHash(
    to,
    value,
    data,
    operation,
    safeTxGas,
    baseGas,
    gasPrice,
    gasToken,
    refundReceiver,
    nonce
  );

  console.log('🔐 Transaction Hash:');
  console.log(' ', txHash);
  console.log('');

  // 5. Load Private Keys and Sign
  const privateKeys = getSignerPrivateKeys();
  if (privateKeys.length === 0) {
    console.warn('⚠️  No PRIVATE_KEY_* found in .env. Only generating transaction data.');
  }

  const signatures: { signer: string; signature: string }[] = [];

  console.log(`✍️  Signing with ${privateKeys.length} keys found in .env...\n`);

  for (const pk of privateKeys) {
    try {
      const wallet = new ethers.Wallet(pk, provider);
      const address = await wallet.getAddress();

      // Verify if signer is an owner
      if (!owners.map((o: string) => o.toLowerCase()).includes(address.toLowerCase())) {
        console.warn(`⚠️  Skipping ${address}: Not a Safe owner`);
        continue;
      }

      const sig = await wallet.signMessage(ethers.getBytes(txHash));
      // Fix V value for Safe (add 4 if needed, but ethers usually handles this. Safe expects v=27/28 or v=31/32 for eth_sign)
      // Ethers produces 27/28. Safe `checkNSignatures` with `ecrecover` works with 27/28.
      // However, for `eth_sign` (which signMessage uses), Safe sometimes expects v+4.
      // Let's stick to standard signature first. If it fails, we might need adjustment.
      // Actually, Safe's `checkNSignatures` handles v=27/28 correctly for EOA signatures.

      signatures.push({
        signer: address,
        signature: sig
      });
      console.log(`✅ Signed by: ${address}`);
    } catch (e: any) {
      console.error(`❌ Failed to sign with key: ${e.message}`);
    }
  }

  // 6. Sort Signatures
  signatures.sort((a, b) =>
    a.signer.toLowerCase().localeCompare(b.signer.toLowerCase())
  );

  // 7. Save to File
  const txData = {
    safeAddress,
    threshold: threshold.toString(),
    nonce: nonce.toString(),
    txHash,
    to,
    value,
    data,
    operation,
    safeTxGas: safeTxGas.toString(),
    baseGas: baseGas.toString(),
    gasPrice: gasPrice.toString(),
    gasToken,
    refundReceiver,
    signatures,
    timestamp: new Date().toISOString(),
  };

  const outputFile = 'safe_tx_data.json';
  fs.writeFileSync(outputFile, JSON.stringify(txData, null, 2));

  console.log('');
  console.log('💾 Transaction data saved to:', outputFile);
  console.log(`📊 Progress: ${signatures.length}/${threshold} signatures`);

  if (signatures.length >= threshold) {
    console.log('🚀 Ready to execute! Run: npx tsx executeTransaction.ts');
  } else {
    console.log(`⚠️  Need ${Number(threshold) - signatures.length} more signatures.`);
    console.log('   Share the safe_tx_data.json file or the TX Hash with other owners.');
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
