import { ethers } from 'ethers';
import fs from 'fs';
import dotenv from 'dotenv';

// Load environment variables
dotenv.config();

const { L1_RPC_URL, PRIVATE_KEY_EXECUTOR } = process.env;

async function main() {
  // Optional: Get extra signature from command line
  const extraSignature = process.argv[2];

  console.log('🚀 Safe Transaction Executor\n');

  if (!L1_RPC_URL) {
    throw new Error('Missing L1_RPC_URL in .env');
  }

  // 1. Load Transaction Data
  const txDataFile = 'safe_tx_data.json';
  if (!fs.existsSync(txDataFile)) {
    throw new Error(`Transaction data file not found: ${txDataFile}. Run signTransaction.ts first.`);
  }

  const txData = JSON.parse(fs.readFileSync(txDataFile, 'utf-8'));

  console.log('📋 Loaded Transaction:');
  console.log('  Safe:', txData.safeAddress);
  console.log('  Nonce:', txData.nonce);
  console.log('  Signatures:', txData.signatures.length);
  console.log('');

  // 2. Handle Extra Signature
  if (extraSignature) {
    console.log('➕ Adding extra signature from CLI...');
    // We don't know the signer address easily without recovering, but Safe needs sorted signatures.
    // So we must recover the address.
    const msgHash = ethers.getBytes(txData.txHash);
    const recoveredAddress = ethers.verifyMessage(msgHash, extraSignature);

    console.log(`   Recovered Signer: ${recoveredAddress}`);

    // Check if already exists
    if (!txData.signatures.find((s: any) => s.signer.toLowerCase() === recoveredAddress.toLowerCase())) {
      txData.signatures.push({
        signer: recoveredAddress,
        signature: extraSignature
      });
      console.log('   Signature added.');
    } else {
      console.log('   Signature already present.');
    }
  }

  // 3. Sort Signatures (Critical for Safe)
  txData.signatures.sort((a: any, b: any) =>
    a.signer.toLowerCase().localeCompare(b.signer.toLowerCase())
  );

  // 4. Check Threshold
  if (txData.signatures.length < Number(txData.threshold)) {
    throw new Error(`Not enough signatures. Need ${txData.threshold}, have ${txData.signatures.length}`);
  }

  // 5. Prepare Execution
  const provider = new ethers.JsonRpcProvider(L1_RPC_URL);

  // Executor can be anyone, but we need a signer to pay for gas.
  // Use PRIVATE_KEY_EXECUTOR if available, otherwise try PRIVATE_KEY_OWNER_1, etc.
  let executorKey = PRIVATE_KEY_EXECUTOR;
  if (!executorKey) {
    // Try to find any private key in env
    const envKeys = Object.keys(process.env).filter(k => k.startsWith('PRIVATE_KEY_'));
    if (envKeys.length > 0) {
      executorKey = process.env[envKeys[0]];
      console.log(`ℹ️  Using ${envKeys[0]} as executor`);
    }
  }

  if (!executorKey) {
    throw new Error('No private key found to execute transaction. Set PRIVATE_KEY_EXECUTOR in .env');
  }

  const executor = new ethers.Wallet(executorKey, provider);
  console.log(`👤 Executor: ${await executor.getAddress()}`);

  // 6. Concatenate Signatures
  const concatenatedSignatures = ethers.concat(txData.signatures.map((s: any) => s.signature));

  // 7. Execute
  const safeAbi = [
    'function execTransaction(address to, uint256 value, bytes calldata data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address payable refundReceiver, bytes memory signatures) payable returns (bool)',
  ];

  const safe = new ethers.Contract(txData.safeAddress, safeAbi, executor);

  console.log('🚀 Sending Transaction...\n');

  try {
    // Estimate gas first
    try {
        await safe.execTransaction.staticCall(
            txData.to,
            txData.value,
            txData.data,
            txData.operation,
            txData.safeTxGas,
            txData.baseGas,
            txData.gasPrice,
            txData.gasToken,
            txData.refundReceiver,
            concatenatedSignatures
        );
        console.log('✅ Simulation successful');
    } catch (e: any) {
        console.warn('⚠️  Simulation failed. Transaction might fail:', e.message);
    }

    const tx = await safe.execTransaction(
      txData.to,
      txData.value,
      txData.data,
      txData.operation,
      txData.safeTxGas,
      txData.baseGas,
      txData.gasPrice,
      txData.gasToken,
      txData.refundReceiver,
      concatenatedSignatures
    );

    console.log('📤 Transaction Submitted:');
    console.log('   TX Hash:', tx.hash);
    console.log('');
    console.log('⏳ Waiting for confirmation...');

    const receipt = await tx.wait();

    console.log('');
    console.log('✅ Transaction Confirmed!');
    console.log('   Block:', receipt?.blockNumber);
    console.log('   Gas Used:', receipt?.gasUsed.toString());
    console.log('');
    console.log('🌐 View on Etherscan:');
    console.log(`   https://sepolia.etherscan.io/tx/${tx.hash}`);

  } catch (error: any) {
    console.error('❌ Transaction Failed:', error.message);
    if (error.data) {
      console.error('   Error Data:', error.data);
    }
    throw error;
  }
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
