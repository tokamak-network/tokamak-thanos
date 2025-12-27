import { ethers } from 'ethers'
import dotenv from 'dotenv'

dotenv.config()

/**
 * Check the owner of a contract (supports various proxy patterns)
 *
 * Usage:
 *   npx tsx checkOwner.ts
 *   npx tsx checkOwner.ts 0x277a690e99C4197d07c29eD57090441Dcb384b30
 */

// Owner storage slot for L1ChugSplashProxy
const OWNER_SLOT = '0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103'

// ProxyAdmin slot for standard proxies (EIP-1967)
const ADMIN_SLOT = '0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103'

async function getOwnerFromSlot(
  provider: ethers.Provider,
  address: string,
  slot: string
): Promise<string | null> {
  try {
    const value = await provider.getStorage(address, slot)
    if (value === ethers.ZeroHash) {
      return null
    }
    // Storage value is a bytes32, extract address from last 20 bytes
    return ethers.getAddress('0x' + value.slice(-40))
  } catch (err) {
    return null
  }
}

async function getOwnerFromFunction(
  provider: ethers.Provider,
  address: string
): Promise<string | null> {
  // Try common owner/admin functions
  const functions = [
    'function owner() view returns (address)',
    'function admin() view returns (address)',
    'function getOwner() view returns (address)',
    'function getProxyAdmin() view returns (address)',
  ]

  for (const funcSig of functions) {
    try {
      const contract = new ethers.Contract(address, [funcSig], provider)
      const functionName = funcSig.split('function ')[1].split('(')[0]
      const owner = await contract[functionName]()
      if (owner && owner !== ethers.ZeroAddress) {
        return ethers.getAddress(owner)
      }
    } catch (err) {
      // Function doesn't exist or call failed, try next
      continue
    }
  }
  return null
}

async function main() {
  const L1_RPC_URL = process.env.L1_RPC_URL
  if (!L1_RPC_URL) {
    throw new Error('L1_RPC_URL not set in .env')
  }

  // Get target address from command line or use default
  const targetAddress = process.argv[2] || '0x277a690e99C4197d07c29eD57090441Dcb384b30'

  if (!ethers.isAddress(targetAddress)) {
    throw new Error(`Invalid address: ${targetAddress}`)
  }

  console.log('🔍 Checking Owner/Admin Information\n')
  console.log('Target Address:', ethers.getAddress(targetAddress))
  console.log('Network:', L1_RPC_URL.includes('sepolia') ? 'Sepolia' : 'Unknown')
  console.log('')

  const provider = new ethers.JsonRpcProvider(L1_RPC_URL)

  // Check contract bytecode
  const code = await provider.getCode(targetAddress)
  if (code === '0x') {
    console.log('❌ Error: Address is not a contract (no bytecode)')
    return
  }

  console.log('✅ Contract confirmed')
  console.log('')

  // Method 1: Check owner from storage slot (L1ChugSplashProxy)
  console.log('📊 Method 1: Reading from storage slot')
  console.log(`   Slot: ${OWNER_SLOT}`)
  const ownerFromSlot = await getOwnerFromSlot(provider, targetAddress, OWNER_SLOT)
  if (ownerFromSlot) {
    console.log(`   ✅ Owner found: ${ownerFromSlot}`)
  } else {
    console.log('   ⚠️  No owner in this slot')
  }
  console.log('')

  // Method 2: Try calling owner functions
  console.log('📊 Method 2: Calling owner/admin functions')
  const ownerFromFunction = await getOwnerFromFunction(provider, targetAddress)
  if (ownerFromFunction) {
    console.log(`   ✅ Owner/Admin found: ${ownerFromFunction}`)
  } else {
    console.log('   ⚠️  No owner function available')
  }
  console.log('')

  // Summary
  console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
  console.log('📋 Summary:')
  if (ownerFromSlot || ownerFromFunction) {
    const owner = ownerFromSlot || ownerFromFunction
    console.log(`   Owner/Admin: ${owner}`)

    // Check if it's a Safe
    const ownerCode = await provider.getCode(owner!)
    if (ownerCode !== '0x') {
      console.log(`   Type: Smart Contract (possibly a Safe multisig)`)

      // Try to get Safe info
      try {
        const safeAbi = [
          'function getThreshold() view returns (uint256)',
          'function getOwners() view returns (address[])',
        ]
        const safe = new ethers.Contract(owner!, safeAbi, provider)
        const threshold = await safe.getThreshold()
        const owners = await safe.getOwners()
        console.log(`   Safe Threshold: ${threshold}/${owners.length}`)
        console.log(`   Safe Owners:`)
        owners.forEach((addr: string, i: number) => {
          console.log(`     ${i + 1}. ${addr}`)
        })
      } catch (err) {
        console.log(`   (Not a Safe or unable to read Safe data)`)
      }
    } else {
      console.log(`   Type: EOA (Externally Owned Account)`)
    }
  } else {
    console.log('   ❌ Owner/Admin not found')
  }
  console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error('❌ Error:', error.message)
    process.exit(1)
  })
