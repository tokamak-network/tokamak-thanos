import fs from 'fs'
import path from 'path'

import { ethers } from 'ethers'
import dotenv from 'dotenv'

dotenv.config()

/**
 * ENV
 * Required:
 *  - L1_RPC_URL
 *  - SAFE_ADDRESS
 *
 * Tx mode A (raw):
 *  - TX_TO
 *  - TX_DATA
 *  - TX_VALUE_WEI (optional, integer string)
 *    OR TX_VALUE_ETH (optional, decimal string)
 *
 * Tx mode B (ChangeProxyAdmin helper):
 *  - PROXY_ADMIN
 *  - TARGET_PROXY
 *  - NEW_ADMIN
 *  - (optional) TX_VALUE_WEI / TX_VALUE_ETH (normally 0)
 *
 * Signing:
 *  - PRIVATE_KEY_1, PRIVATE_KEY_2, ... (any keys starting with PRIVATE_KEY_)
 *
 * Optional:
 *  - TX_NONCE (override Safe nonce)
 *  - OUTPUT_FILE (default: safe_tx_data.json)
 *
 * Safe Transaction Service integration (optional):
 *  - SAFE_TX_SERVICE_URL (e.g., https://safe-transaction-sepolia.safe.global)
 *  - SAFE_TX_SERVICE_TOKEN (optional Authorization: Bearer value)
 *  - SAFE_TX_SERVICE_ORIGIN (optional string shown in Safe UI)
 */

type SignatureEntry = { signer: string; signature: string }

type SafeTxOutput = {
  safeAddress: string
  chainId: string
  threshold: string
  nonce: string

  to: string
  value: string // wei (decimal string)
  data: string
  operation: number

  safeTxGas: string
  baseGas: string
  gasPrice: string
  gasToken: string
  refundReceiver: string

  txHash: string // safeTxHash
  owners: string[]

  signatures: SignatureEntry[]
  signaturesBytes: string // concatenated signatures for execTransaction

  mode: 'RAW' | 'CHANGE_PROXY_ADMIN'
  timestamp: string
}

function requireEnv(name: string): string {
  const v = process.env[name]
  if (!v || v.trim().length === 0) {
    throw new Error(`Missing env: ${name}`)
  }
  return v.trim()
}

function optionalEnv(name: string): string | undefined {
  const v = process.env[name]
  if (!v || v.trim().length === 0) {
    return undefined
  }
  return v.trim()
}

function mustBeAddress(label: string, value: string): string {
  if (!ethers.isAddress(value)) {
    throw new Error(`${label} is not a valid address: ${value}`)
  }
  return ethers.getAddress(value) // checksum
}

function mustBeHexData(label: string, value: string): string {
  if (!ethers.isHexString(value)) {
    throw new Error(`${label} must be a hex string (0x...): ${value}`)
  }
  return value
}

function getSignerPrivateKeys(): string[] {
  return Object.keys(process.env)
    .filter((k) => k.startsWith('PRIVATE_KEY_'))
    .map((k) => (process.env[k] ?? '').trim())
    .filter((v) => v.length > 0)
}

function validatePrivateKey(pk: string): string {
  // ethers.Wallet will also validate, but we fail early.
  if (!/^0x[0-9a-fA-F]{64}$/.test(pk)) {
    throw new Error(`Invalid private key format: ${pk.slice(0, 10)}...`)
  }
  return pk
}

function parseValueWei(): bigint {
  const wei = optionalEnv('TX_VALUE_WEI')
  const eth = optionalEnv('TX_VALUE_ETH')

  if (wei && eth) {
    throw new Error(
      'Provide only one of TX_VALUE_WEI or TX_VALUE_ETH (not both).'
    )
  }
  if (wei) {
    if (!/^\d+$/.test(wei)) {
      throw new Error(
        'TX_VALUE_WEI must be an integer string in wei (e.g., 0, 1000000000000000).'
      )
    }
    return BigInt(wei)
  }
  if (eth) {
    // decimal eth → wei
    return ethers.parseEther(eth)
  }
  return 0n
}

function concatSignatures(signatures: SignatureEntry[]): string {
  // Each signature is 65 bytes hex string (0x + 130 hex)
  // Safe expects concatenated signatures sorted by signer (ascending)
  const chunks = signatures.map((s) => s.signature.replace(/^0x/, ''))
  return '0x' + chunks.join('')
}

function sortBySigner(signatures: SignatureEntry[]): SignatureEntry[] {
  return [...signatures].sort((a, b) =>
    a.signer.toLowerCase().localeCompare(b.signer.toLowerCase())
  )
}

type FetchFn = (url: string, init?: any) => Promise<any>

async function optionalFetchJson(url: string, init?: any): Promise<any> {
  const fetchFn: FetchFn | undefined = (globalThis as any).fetch
  if (!fetchFn) {
    throw new Error(
      'Safe API upload requested but fetch is unavailable. Use Node.js 18+.'
    )
  }
  return fetchFn(url, init)
}

async function registerWithSafeService(params: {
  baseUrl: string
  token?: string
  origin?: string
  safeTx: SafeTxOutput
  signatures: SignatureEntry[]
}) {
  const { baseUrl, token, origin, safeTx, signatures } = params
  if (signatures.length === 0) {
    console.warn(
      '⚠️  Skipping Safe API upload: no signatures available to submit.'
    )
    return
  }

  const apiRoot = baseUrl.replace(/\/+$/, '')
  const txEndpoint = `${apiRoot}/api/v1/safes/${safeTx.safeAddress}/multisig-transactions/`
  console.log(`🔎 Safe API tx endpoint: ${txEndpoint}`)

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const primarySig = signatures[0]
  const payload: Record<string, any> = {
    safe: safeTx.safeAddress,
    to: safeTx.to,
    value: safeTx.value,
    data: safeTx.data,
    operation: safeTx.operation,
    safeTxGas: safeTx.safeTxGas,
    baseGas: safeTx.baseGas,
    gasPrice: safeTx.gasPrice,
    gasToken: safeTx.gasToken,
    refundReceiver: safeTx.refundReceiver,
    nonce: Number(safeTx.nonce),
    contractTransactionHash: safeTx.txHash,
    sender: primarySig.signer,
    signature: primarySig.signature,
  }
  if (origin) {
    payload.origin = origin
  }

  console.log('🌐 Uploading transaction to Safe Transaction Service...')
  const res = await optionalFetchJson(txEndpoint, {
    method: 'POST',
    headers,
    body: JSON.stringify(payload),
  })

  if (!res.ok && res.status !== 409) {
    const text = await res.text()
    throw new Error(
      `Safe API transaction upload failed (${res.status}): ${text}`
    )
  }

  if (res.status === 201) {
    console.log('✅ Safe API: transaction registered')
  } else if (res.status === 409) {
    console.log('ℹ️  Safe API: transaction already exists, adding confirmations')
  }

  for (const sig of signatures.slice(1)) {
    const confirmEndpoint = `${apiRoot}/api/v1/multisig-transactions/${safeTx.txHash}/confirmations/`
    console.log(`🔎 Safe API confirm endpoint: ${confirmEndpoint}`)
    const confirmRes = await optionalFetchJson(confirmEndpoint, {
      method: 'POST',
      headers,
      body: JSON.stringify({
        owner: sig.signer,
        signature: sig.signature,
      }),
    })
    if (confirmRes.ok) {
      console.log(`   ➕ Confirmation recorded for ${sig.signer}`)
      continue
    }
    if (confirmRes.status === 409 || confirmRes.status === 400) {
      const text = await confirmRes.text()
      console.warn(
        `   ⚠️  Failed to add confirmation for ${sig.signer} (${confirmRes.status}): ${text}`
      )
      continue
    }
    const text = await confirmRes.text()
    throw new Error(
      `Safe API confirmation failed for ${sig.signer} (${confirmRes.status}): ${text}`
    )
  }
}

/**
 * Produce an EOA signature over a digest (bytes32) WITHOUT the "\x19Ethereum Signed Message" prefix.
 * This is the safest default for Safe: it will be validated against safeTxHash directly.
 */
function signDigest(wallet: ethers.Wallet, digestHex32: string): string {
  if (!/^0x[0-9a-fA-F]{64}$/.test(digestHex32)) {
    throw new Error(`digest must be bytes32 hex: ${digestHex32}`)
  }
  const sig = wallet.signingKey.sign(digestHex32)
  return sig.serialized // 0x{r}{s}{v}
}

async function main() {
  console.log('🚀 Safe Transaction Signer (Refactored)\n')

  const L1_RPC_URL = requireEnv('L1_RPC_URL')
  const SAFE_ADDRESS = mustBeAddress('SAFE_ADDRESS', requireEnv('SAFE_ADDRESS'))

  const provider = new ethers.JsonRpcProvider(L1_RPC_URL)

  const safeAbi = [
    'function getThreshold() view returns (uint256)',
    'function nonce() view returns (uint256)',
    'function getOwners() view returns (address[])',
    'function getTransactionHash(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 _nonce) view returns (bytes32)',
  ] as const

  const safe = new ethers.Contract(SAFE_ADDRESS, safeAbi, provider)

  const network = await provider.getNetwork()
  const chainId = network.chainId.toString()

  const threshold: bigint = await safe.getThreshold()
  const onchainNonce: bigint = await safe.nonce()
  const owners: string[] = (await safe.getOwners()).map((o: string) =>
    ethers.getAddress(o)
  )

  const ownerSet = new Set(owners.map((o) => o.toLowerCase()))

  const overrideNonce = optionalEnv('TX_NONCE')
  let nonceToUse = onchainNonce
  if (overrideNonce !== undefined) {
    if (!/^\d+$/.test(overrideNonce)) {
      throw new Error('TX_NONCE must be an integer string.')
    }
    nonceToUse = BigInt(overrideNonce)
  }

  // ---- Determine tx details ----
  const rawTo = optionalEnv('TX_TO')
  const rawData = optionalEnv('TX_DATA')

  const PROXY_ADMIN = optionalEnv('PROXY_ADMIN')
  const TARGET_PROXY = optionalEnv('TARGET_PROXY')
  const NEW_ADMIN = optionalEnv('NEW_ADMIN')

  let mode: 'RAW' | 'CHANGE_PROXY_ADMIN' = 'RAW'
  let to: string
  let data: string

  if (!rawTo && PROXY_ADMIN && TARGET_PROXY && NEW_ADMIN) {
    mode = 'CHANGE_PROXY_ADMIN'
    console.log('ℹ️  Mode: ChangeProxyAdmin')

    to = mustBeAddress('PROXY_ADMIN', PROXY_ADMIN)

    const proxyAdminInterface = new ethers.Interface([
      'function changeProxyAdmin(address proxy, address newAdmin)',
    ])
    data = proxyAdminInterface.encodeFunctionData('changeProxyAdmin', [
      mustBeAddress('TARGET_PROXY', TARGET_PROXY),
      mustBeAddress('NEW_ADMIN', NEW_ADMIN),
    ])
  } else {
    // RAW mode
    if (!rawTo) {
      throw new Error('TX_TO is required in RAW mode.')
    }
    if (!rawData) {
      throw new Error('TX_DATA is required in RAW mode.')
    }
    to = mustBeAddress('TX_TO', rawTo)
    data = mustBeHexData('TX_DATA', rawData)
    if (data === '0x') {
      throw new Error(
        'TX_DATA must not be empty (0x). Refusing to sign an empty calldata transaction.'
      )
    }
  }

  const valueWei = parseValueWei() // bigint
  const value = valueWei.toString()

  // ---- Safe meta params (set to 0 by default) ----
  const operation = 0 // Call
  const safeTxGas = 0n
  const baseGas = 0n
  const gasPrice = 0n
  const gasToken = ethers.ZeroAddress
  const refundReceiver = ethers.ZeroAddress

  // ---- Compute safeTxHash ----
  const txHash: string = await safe.getTransactionHash(
    to,
    valueWei,
    data,
    operation,
    safeTxGas,
    baseGas,
    gasPrice,
    gasToken,
    refundReceiver,
    nonceToUse
  )

  console.log('📊 Safe Status:')
  console.log('  Safe:', SAFE_ADDRESS)
  console.log('  ChainId:', chainId)
  console.log('  Threshold:', threshold.toString())
  console.log('  Nonce (on-chain):', onchainNonce.toString())
  console.log('  Nonce (used):', nonceToUse.toString())
  console.log('  Owners:', owners.length)
  console.log('')

  console.log('📋 Transaction Details:')
  console.log('  Mode:', mode)
  console.log('  To:', to)
  console.log('  Value (wei):', value)
  console.log('  Data length:', data.length)
  console.log('')

  console.log('🔐 SafeTxHash (getTransactionHash):')
  console.log(' ', txHash)
  console.log('')

  // ---- Load keys & sign ----
  const privateKeys = getSignerPrivateKeys().map(validatePrivateKey)
  if (privateKeys.length === 0) {
    console.warn(
      '⚠️  No PRIVATE_KEY_* found in .env. Only generating txHash + tx payload.'
    )
  } else {
    console.log(
      `✍️  Signing with ${privateKeys.length} keys found in .env...\n`
    )
  }

  const signatures: SignatureEntry[] = []

  for (const pk of privateKeys) {
    try {
      const wallet = new ethers.Wallet(pk, provider)
      const signer = ethers.getAddress(await wallet.getAddress())

      if (!ownerSet.has(signer.toLowerCase())) {
        console.warn(`⚠️  Skipping ${signer}: Not a Safe owner`)
        continue
      }

      // IMPORTANT: sign digest directly (no prefix)
      const signature = signDigest(wallet, txHash)

      signatures.push({ signer, signature })
      console.log(`✅ Signed by: ${signer}`)
    } catch (e: any) {
      console.error(`❌ Failed to sign with a key: ${e?.message ?? String(e)}`)
    }
  }

  const sorted = sortBySigner(signatures)
  const signaturesBytes = concatSignatures(sorted)

  // ---- Save output ----
  const output: SafeTxOutput = {
    safeAddress: SAFE_ADDRESS,
    chainId,
    threshold: threshold.toString(),
    nonce: nonceToUse.toString(),

    to,
    value,
    data,
    operation,

    safeTxGas: safeTxGas.toString(),
    baseGas: baseGas.toString(),
    gasPrice: gasPrice.toString(),
    gasToken,
    refundReceiver,

    txHash,
    owners,

    signatures: sorted,
    signaturesBytes,

    mode,
    timestamp: new Date().toISOString(),
  }

  const outputDir = path.join(process.cwd(), 'data')
  fs.mkdirSync(outputDir, { recursive: true })

  const baseName = optionalEnv('OUTPUT_FILE') ?? 'safe_tx_data.json'
  const safeBaseName = baseName.replace(/\.json$/i, '')
  const timestamp = new Date()
    .toISOString()
    .replace(/[:.]/g, '-')
    .replace(/Z$/, 'Z')
  const outputFile = path.join(
    outputDir,
    `${safeBaseName}_${timestamp}.json`
  )
  fs.writeFileSync(outputFile, JSON.stringify(output, null, 2))

  console.log('')
  console.log('💾 Saved:', outputFile)
  console.log(
    `📊 Progress: ${sorted.length}/${threshold.toString()} signatures`
  )

  if (BigInt(sorted.length) >= threshold) {
    console.log(
      '🚀 Ready to execute: you can use output.signaturesBytes in execTransaction.'
    )
  } else {
    console.log(
      `⚠️  Need ${threshold - BigInt(sorted.length)} more signature(s).`
    )
  }

  const safeServiceUrl = optionalEnv('SAFE_TX_SERVICE_URL')
  const safeServiceToken = optionalEnv('SAFE_TX_SERVICE_TOKEN')
  const safeServiceOrigin = optionalEnv('SAFE_TX_SERVICE_ORIGIN')
  if (safeServiceUrl) {
    try {
      await registerWithSafeService({
        baseUrl: safeServiceUrl,
        token: safeServiceToken,
        origin: safeServiceOrigin,
        safeTx: output,
        signatures: sorted,
      })
    } catch (err: any) {
      console.error(`❌ Safe API upload failed: ${err?.message ?? String(err)}`)
    }
  } else {
    console.log('ℹ️  SAFE_TX_SERVICE_URL not set. Skipping Safe API upload.')
  }
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
