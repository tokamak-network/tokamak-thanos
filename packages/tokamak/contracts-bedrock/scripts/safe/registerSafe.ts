import dotenv from 'dotenv'
import { ethers } from 'ethers'
import Safe from '@safe-global/protocol-kit'
import SafeApiKit from '@safe-global/api-kit'

dotenv.config()

/**
 * Register a Safe with the Safe Transaction Service by proposing a dummy transaction.
 * This script is useful when a Safe exists on-chain but is not indexed by the Transaction Service.
 *
 * ENV Required:
 *  - L1_RPC_URL: RPC endpoint URL
 *  - SAFE_ADDRESS: Safe contract address
 *  - PRIVATE_KEY_1: At least one owner's private key
 *  - SAFE_TX_SERVICE_URL: Transaction Service URL (e.g., https://api.safe.global/tx-service/sep)
 */

function requireEnv(name: string): string {
  const v = process.env[name]
  if (!v || v.trim().length === 0) {
    throw new Error(`Missing env: ${name}`)
  }
  return v.trim()
}

function getFirstPrivateKey(): string {
  const keys = Object.keys(process.env)
    .filter((k) => k.startsWith('PRIVATE_KEY_'))
    .map((k) => (process.env[k] ?? '').trim())
    .filter((v) => v.length > 0)

  if (keys.length === 0) {
    throw new Error('No PRIVATE_KEY_* found in environment variables')
  }

  return keys[0]
}

async function main() {
  console.log('🔧 Safe Transaction Service Registration Tool\n')

  const L1_RPC_URL = requireEnv('L1_RPC_URL')
  const SAFE_ADDRESS = requireEnv('SAFE_ADDRESS')
  const SAFE_TX_SERVICE_URL = requireEnv('SAFE_TX_SERVICE_URL')
  const PRIVATE_KEY = getFirstPrivateKey()

  console.log('📊 Configuration:')
  console.log('  Safe Address:', SAFE_ADDRESS)
  console.log('  RPC URL:', L1_RPC_URL)
  console.log('  Transaction Service:', SAFE_TX_SERVICE_URL)
  console.log('')

  // Initialize provider and signer
  const provider = new ethers.JsonRpcProvider(L1_RPC_URL)
  const signer = new ethers.Wallet(PRIVATE_KEY, provider)
  const signerAddress = await signer.getAddress()

  console.log('👤 Signer:', signerAddress)
  console.log('')

  // Initialize Safe Protocol Kit
  console.log('🔄 Initializing Safe Protocol Kit...')
  const protocolKit = await Safe.default.init({
    provider: L1_RPC_URL,
    signer: PRIVATE_KEY,
    safeAddress: SAFE_ADDRESS,
  })

  const owners = await protocolKit.getOwners()
  const threshold = await protocolKit.getThreshold()
  const nonce = await protocolKit.getNonce()

  console.log('✅ Safe initialized:')
  console.log('  Owners:', owners.length)
  console.log('  Threshold:', threshold)
  console.log('  Nonce:', nonce)
  console.log('')

  // Verify signer is an owner
  if (!owners.map((o) => o.toLowerCase()).includes(signerAddress.toLowerCase())) {
    throw new Error(`Signer ${signerAddress} is not an owner of this Safe`)
  }

  // Initialize Safe API Kit
  console.log('🔄 Initializing Safe API Kit...')
  const apiKit = new SafeApiKit.default({
    chainId: (await provider.getNetwork()).chainId,
    txServiceUrl: SAFE_TX_SERVICE_URL,
  })

  console.log('✅ API Kit initialized')
  console.log('')

  // Create a dummy transaction (Safe sends 0 ETH to itself with minimal data)
  console.log('📝 Creating dummy transaction to register Safe...')

  // Use a simple data payload that is not empty (to pass validation)
  const safeTransactionData = {
    to: SAFE_ADDRESS,
    value: '0',
    data: '0x00', // Minimal non-empty data
    operation: 0, // Call
  }

  const safeTransaction = await protocolKit.createTransaction({
    transactions: [safeTransactionData],
  })

  const safeTxHash = await protocolKit.getTransactionHash(safeTransaction)
  console.log('📋 SafeTxHash:', safeTxHash)

  const senderSignature = await protocolKit.signHash(safeTxHash)
  console.log('✍️  Signed by:', signerAddress)
  console.log('')

  // Propose transaction to Safe Transaction Service using direct API call
  console.log('🌐 Proposing transaction to Safe Transaction Service...')

  try {
    // Build the payload manually
    const txServiceUrl = SAFE_TX_SERVICE_URL.replace(/\/+$/, '')
    const endpoint = `${txServiceUrl}/api/v1/multisig-transactions/`

    const payload = {
      safe: SAFE_ADDRESS,
      to: safeTransaction.data.to,
      value: safeTransaction.data.value,
      data: safeTransaction.data.data || '0x',
      operation: safeTransaction.data.operation,
      safeTxGas: safeTransaction.data.safeTxGas,
      baseGas: safeTransaction.data.baseGas,
      gasPrice: safeTransaction.data.gasPrice,
      gasToken: safeTransaction.data.gasToken,
      refundReceiver: safeTransaction.data.refundReceiver,
      nonce: safeTransaction.data.nonce,
      contractTransactionHash: safeTxHash,
      sender: signerAddress,
      signature: senderSignature.data,
      origin: 'Safe Registration Script',
    }

    console.log('📤 POST', endpoint)

    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })

    if (response.ok || response.status === 201) {
      console.log('✅ Transaction proposed successfully!')
      console.log('')
      console.log('🎉 Safe has been registered with Transaction Service!')
    } else if (response.status === 422) {
      const errorData = await response.json()
      console.log('ℹ️  Response:', JSON.stringify(errorData, null, 2))
      if (
        errorData.nonFieldErrors?.some((e: string) =>
          e.includes('already executed')
        ) ||
        JSON.stringify(errorData).includes('already exists')
      ) {
        console.log('ℹ️  Safe/Transaction already registered')
        console.log('✅ Registration successful!')
      } else {
        throw new Error(
          `API returned 422: ${JSON.stringify(errorData, null, 2)}`
        )
      }
    } else {
      const errorText = await response.text()
      throw new Error(
        `API returned ${response.status}: ${errorText.slice(0, 500)}`
      )
    }

    console.log('')
    console.log('📌 Next steps:')
    console.log('  1. Visit https://app.safe.global')
    console.log(`  2. Add Safe: ${SAFE_ADDRESS}`)
    console.log('  3. The Safe should now be visible in the UI')
    console.log('  4. You can reject/execute the dummy transaction from the UI')
    console.log('')
  } catch (error: any) {
    if (error.message?.includes('Not Found')) {
      console.log('')
      console.log('⚠️  Safe Transaction Service returned 404')
      console.log(
        '    This might mean the Safe is not recognized by the service yet.'
      )
      console.log('')
      console.log('💡 Alternative solutions:')
      console.log(
        '    1. Execute an on-chain transaction from the Safe to trigger indexing'
      )
      console.log('    2. Wait for the indexer to catch up (can take time)')
      console.log(
        '    3. Contact Safe support if this Safe should be supported'
      )
      console.log('')
    }
    throw error
  }
}

main().catch((err) => {
  console.error('❌ Error:', err.message ?? String(err))
  process.exit(1)
})
