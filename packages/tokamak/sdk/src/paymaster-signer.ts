import { ethers } from 'ethers'

// ---------------------------------------------------------------------------
// Predeploy addresses (all Thanos presets)
// ---------------------------------------------------------------------------
const ENTRYPOINT_V08 = '0x4200000000000000000000000000000000000063'
const MULTI_TOKEN_PAYMASTER = '0x4200000000000000000000000000000000000067'
const USDC_ADDRESS = '0x4200000000000000000000000000000000000778'
const WETH_ADDRESS = '0x4200000000000000000000000000000000000006'
const DEFAULT_BUNDLER_URL = 'http://localhost:4337'

// ---------------------------------------------------------------------------
// ABIs
// ---------------------------------------------------------------------------
const SIMPLE_7702_ACCOUNT_ABI = [
  'function execute(address target, uint256 value, bytes calldata data)',
  'function executeBatch((address target, uint256 value, bytes data)[] calls)',
]

const ERC20_ABI = [
  'function approve(address spender, uint256 amount) returns (bool)',
  'function allowance(address owner, address spender) view returns (uint256)',
]

const ENTRYPOINT_ABI = [
  'function getNonce(address sender, uint192 key) view returns (uint256)',
]

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------
export interface PaymasterOptions {
  feeToken: 'USDC' | 'WETH' | 'native'
  bundlerUrl?: string
}

interface PackedUserOperation {
  sender: string
  nonce: ethers.BigNumber
  initCode: string
  callData: string
  accountGasLimits: string
  preVerificationGas: ethers.BigNumber
  gasFees: string
  paymasterAndData: string
  signature: string
}

// ---------------------------------------------------------------------------
// Internal helpers (arrow functions to satisfy prefer-arrow lint rule)
// ---------------------------------------------------------------------------

/**
 * Packs two uint128 values into a single bytes32 value.
 * high occupies the upper 128 bits, low occupies the lower 128 bits.
 */
const packUint128x2 = (high: bigint, low: bigint): string => {
  const packed = (BigInt(high) << BigInt(128)) | BigInt(low)
  return ethers.utils.hexZeroPad(
    ethers.BigNumber.from(packed).toHexString(),
    32
  )
}

const buildPaymasterAndData = (tokenAddr: string): string =>
  ethers.utils.hexConcat([MULTI_TOKEN_PAYMASTER, tokenAddr])

const buildUserOpHash = (
  userOp: PackedUserOperation,
  chainId: number
): string => {
  const encodedFields = ethers.utils.defaultAbiCoder.encode(
    [
      'address', // sender
      'uint256', // nonce
      'bytes32', // initCode hash
      'bytes32', // callData hash
      'bytes32', // accountGasLimits
      'uint256', // preVerificationGas
      'bytes32', // gasFees
      'bytes32', // paymasterAndData hash
    ],
    [
      userOp.sender,
      userOp.nonce,
      ethers.utils.keccak256(userOp.initCode),
      ethers.utils.keccak256(userOp.callData),
      userOp.accountGasLimits,
      userOp.preVerificationGas,
      userOp.gasFees,
      ethers.utils.keccak256(userOp.paymasterAndData),
    ]
  )

  const innerHash = ethers.utils.keccak256(encodedFields)

  const outerEncoded = ethers.utils.defaultAbiCoder.encode(
    ['bytes32', 'address', 'uint256'],
    [innerHash, ENTRYPOINT_V08, chainId]
  )

  return ethers.utils.keccak256(outerEncoded)
}

const sendRawUserOp = async (
  bundlerUrl: string,
  userOp: PackedUserOperation
): Promise<string> => {
  const serializedUserOp = {
    sender: userOp.sender,
    nonce: ethers.utils.hexlify(userOp.nonce),
    initCode: userOp.initCode,
    callData: userOp.callData,
    accountGasLimits: userOp.accountGasLimits,
    preVerificationGas: ethers.utils.hexlify(userOp.preVerificationGas),
    gasFees: userOp.gasFees,
    paymasterAndData: userOp.paymasterAndData,
    signature: userOp.signature,
  }

  const response = await fetch(bundlerUrl, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: 1,
      method: 'eth_sendUserOperation',
      params: [serializedUserOp, ENTRYPOINT_V08],
    }),
  })

  const data = (await response.json()) as {
    result?: string
    error?: { message: string; code?: number }
  }

  if (data.error) {
    throw new Error(
      `Bundler error: ${data.error.message} (code: ${
        data.error.code ?? 'unknown'
      })`
    )
  }

  if (!data.result) {
    throw new Error('Bundler returned empty result for eth_sendUserOperation')
  }

  return data.result
}

const waitForReceipt = async (
  bundlerUrl: string,
  userOpHash: string
): Promise<string> => {
  const MAX_ATTEMPTS = 20
  const POLL_INTERVAL_MS = 3000

  for (let i = 0; i < MAX_ATTEMPTS; i++) {
    const response = await fetch(bundlerUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        jsonrpc: '2.0',
        id: 1,
        method: 'eth_getUserOperationReceipt',
        params: [userOpHash],
      }),
    })

    const data = (await response.json()) as {
      result?: { receipt: { transactionHash: string } } | null
      error?: { message: string }
    }

    if (data.error) {
      throw new Error(`Bundler receipt error: ${data.error.message}`)
    }

    if (data.result?.receipt?.transactionHash) {
      return data.result.receipt.transactionHash
    }

    // Wait before next poll (except last iteration)
    if (i < MAX_ATTEMPTS - 1) {
      await new Promise<void>((resolve) =>
        setTimeout(resolve, POLL_INTERVAL_MS)
      )
    }
  }

  throw new Error(
    `UserOperation receipt timeout: no receipt after ${MAX_ATTEMPTS} attempts for ${userOpHash}`
  )
}

// ---------------------------------------------------------------------------
// Core sendAsUserOp implementation
// ---------------------------------------------------------------------------

const sendAsUserOp = async (
  signer: ethers.Signer,
  tx: ethers.providers.TransactionRequest,
  options: PaymasterOptions
): Promise<ethers.providers.TransactionResponse> => {
  const provider = signer.provider
  if (!provider) {
    throw new Error('Signer must have a provider attached for UserOp sending')
  }

  const bundlerUrl = options.bundlerUrl ?? DEFAULT_BUNDLER_URL

  // 1. Get sender address
  const sender = await signer.getAddress()

  // 2. Check EIP-7702 delegation
  const code = await provider.getCode(sender)
  if (!code.startsWith('0xef0100')) {
    throw new Error(
      `EIP-7702 delegation not set for ${sender}. Run setupAAPaymaster first.`
    )
  }

  // 3. Resolve token address
  const tokenAddresses: Record<string, string> = {
    USDC: USDC_ADDRESS,
    WETH: WETH_ADDRESS,
  }
  const tokenAddr = tokenAddresses[options.feeToken]
  if (!tokenAddr) {
    throw new Error(`Unsupported feeToken: ${options.feeToken}`)
  }

  const simple7702Iface = new ethers.utils.Interface(SIMPLE_7702_ACCOUNT_ABI)
  const erc20Iface = new ethers.utils.Interface(ERC20_ABI)
  const entrypointContract = new ethers.Contract(
    ENTRYPOINT_V08,
    ENTRYPOINT_ABI,
    provider
  )
  const erc20Contract = new ethers.Contract(tokenAddr, ERC20_ABI, provider)

  // 4. Check token allowance and approve if needed
  const currentAllowance = (await erc20Contract.allowance(
    sender,
    MULTI_TOKEN_PAYMASTER
  )) as ethers.BigNumber

  const network = await provider.getNetwork()
  const chainId = network.chainId

  if (currentAllowance.isZero()) {
    // Build and send approve UserOp first
    const approveCalldata = erc20Iface.encodeFunctionData('approve', [
      MULTI_TOKEN_PAYMASTER,
      ethers.constants.MaxUint256,
    ])
    const approveExecuteCalldata = simple7702Iface.encodeFunctionData(
      'execute',
      [tokenAddr, 0, approveCalldata]
    )

    const approveNonce = (await entrypointContract.getNonce(
      sender,
      0
    )) as ethers.BigNumber
    const approveFeeData = await provider.getFeeData()

    const approveMaxFeePerGas =
      approveFeeData.maxFeePerGas ?? ethers.utils.parseUnits('1', 'gwei')
    const approveMaxPriorityFeePerGas =
      approveFeeData.maxPriorityFeePerGas ??
      ethers.utils.parseUnits('1', 'gwei')

    const approveGasLimits = packUint128x2(
      BigInt(150000), // verificationGasLimit
      BigInt(200000) // callGasLimit
    )
    const approveGasFees = packUint128x2(
      BigInt(approveMaxPriorityFeePerGas.toString()),
      BigInt(approveMaxFeePerGas.toString())
    )

    const approvePaymasterAndData = buildPaymasterAndData(tokenAddr)

    const approveUserOpForHash: PackedUserOperation = {
      sender,
      nonce: approveNonce,
      initCode: '0x',
      callData: approveExecuteCalldata,
      accountGasLimits: approveGasLimits,
      preVerificationGas: ethers.BigNumber.from(21000),
      gasFees: approveGasFees,
      paymasterAndData: approvePaymasterAndData,
      signature: '0x',
    }

    const approveUserOpHash = buildUserOpHash(approveUserOpForHash, chainId)
    const approveSignature = await signer.signMessage(
      ethers.utils.arrayify(approveUserOpHash)
    )

    const approveUserOp: PackedUserOperation = {
      ...approveUserOpForHash,
      signature: approveSignature,
    }

    const approveOpHash = await sendRawUserOp(bundlerUrl, approveUserOp)
    // Wait for approve receipt before proceeding
    await waitForReceipt(bundlerUrl, approveOpHash)
  }

  // 5. Build bridge tx callData
  const bridgeCallData = simple7702Iface.encodeFunctionData('execute', [
    tx.to ?? ethers.constants.AddressZero,
    tx.value ?? 0,
    tx.data ?? '0x',
  ])

  // 6. Get nonce for bridge UserOp
  const bridgeNonce = (await entrypointContract.getNonce(
    sender,
    0
  )) as ethers.BigNumber

  // 7. Gas parameters
  const feeData = await provider.getFeeData()
  const maxFeePerGas =
    feeData.maxFeePerGas ?? ethers.utils.parseUnits('1', 'gwei')
  const maxPriorityFeePerGas =
    feeData.maxPriorityFeePerGas ?? ethers.utils.parseUnits('1', 'gwei')

  const accountGasLimits = packUint128x2(
    BigInt(150000), // verificationGasLimit
    BigInt(200000) // callGasLimit
  )
  const gasFees = packUint128x2(
    BigInt(maxPriorityFeePerGas.toString()),
    BigInt(maxFeePerGas.toString())
  )

  // 8. paymasterAndData (Phase 1: 40 bytes)
  const paymasterAndData = buildPaymasterAndData(tokenAddr)

  // 9. Build UserOp without signature, compute hash
  const userOpForHash: PackedUserOperation = {
    sender,
    nonce: bridgeNonce,
    initCode: '0x',
    callData: bridgeCallData,
    accountGasLimits,
    preVerificationGas: ethers.BigNumber.from(21000),
    gasFees,
    paymasterAndData,
    signature: '0x',
  }

  const userOpHash = buildUserOpHash(userOpForHash, chainId)

  // 10. Sign UserOp hash
  const signature = await signer.signMessage(ethers.utils.arrayify(userOpHash))

  const userOp: PackedUserOperation = {
    ...userOpForHash,
    signature,
  }

  // 11. Send UserOp to bundler
  const bundlerUserOpHash = await sendRawUserOp(bundlerUrl, userOp)

  // 12. Poll for receipt
  const txHash = await waitForReceipt(bundlerUrl, bundlerUserOpHash)

  // 13. Return TransactionResponse-like object
  return {
    hash: txHash,
    wait: (confirmations?: number) =>
      provider.waitForTransaction(txHash, confirmations),
    confirmations: 0,
    from: sender,
    nonce: 0,
    gasLimit: ethers.BigNumber.from(0),
    data: '0x',
    value: ethers.BigNumber.from(0),
    chainId,
  } as ethers.providers.TransactionResponse
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/**
 * Wraps a signer with a Proxy that intercepts sendTransaction calls and
 * routes them through the Alto bundler as ERC-4337 PackedUserOperations.
 *
 * When feeToken is 'native', the original signer is returned unchanged.
 */
export const wrapWithPaymaster = (
  signer: ethers.Signer,
  options: PaymasterOptions
): ethers.Signer => {
  if (options.feeToken === 'native') {
    return signer
  }

  return new Proxy(signer, {
    get: (target: ethers.Signer, prop: string | symbol) => {
      if (prop === 'sendTransaction') {
        return (
          tx: ethers.providers.TransactionRequest
        ): Promise<ethers.providers.TransactionResponse> =>
          sendAsUserOp(target, tx, options)
      }
      return Reflect.get(target, prop)
    },
  })
}
