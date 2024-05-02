import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BigNumber, BytesLike, ethers } from 'ethers'

import { Portals, NumberLike, asL2Provider, MessageStatus } from '../src'

console.log('Setup task...')

const privateKey = process.env.PRIVATE_KEY as BytesLike

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)

const erc20ABI = [
  {
    inputs: [
      { internalType: 'address', name: '_spender', type: 'address' },
      { internalType: 'uint256', name: '_value', type: 'uint256' },
    ],
    name: 'approve',
    outputs: [{ internalType: 'bool', name: '', type: 'bool' }],
    stateMutability: 'nonpayable',
    type: 'function',
  },
  {
    constant: true,
    inputs: [{ name: '_owner', type: 'address' }],
    name: 'balanceOf',
    outputs: [{ name: 'balance', type: 'uint256' }],
    type: 'function',
  },
  {
    inputs: [{ internalType: 'uint256', name: 'amount', type: 'uint256' }],
    name: 'faucet',
    outputs: [],
    stateMutability: 'nonpayable',
    type: 'function',
  },
  {
    constant: false,
    inputs: [
      {
        internalType: 'address',
        name: 'src',
        type: 'address',
      },
      {
        internalType: 'address',
        name: 'dst',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'amount',
        type: 'uint256',
      },
    ],
    name: 'transferFrom',
    outputs: [{ internalType: 'bool', name: '', type: 'bool' }],
    payable: false,
    stateMutability: 'nonpayable',
    type: 'function',
  },
]

let l2NativeToken = process.env.NATIVE_TOKEN || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''

const updateAddresses = async (hre: HardhatRuntimeEnvironment) => {
  if (l2NativeToken === '') {
    const Deployment__L2NativeToken = await hre.deployments.get('L2NativeToken')
    l2NativeToken = Deployment__L2NativeToken.address
  }

  if (addressManager === '') {
    const Deployment__AddressManager = await hre.deployments.get(
      'AddressManager'
    )
    addressManager = Deployment__AddressManager.address
  }

  if (optimismPortal === '') {
    const Deployment__OptimismPortal = await hre.deployments.get(
      'OptimismPortalProxy'
    )
    optimismPortal = Deployment__OptimismPortal.address
  }

  if (l2OutputOracle === '') {
    const Deployment__L2OutputOracle = await hre.deployments.get(
      'L2OutputOracleProxy'
    )
    l2OutputOracle = Deployment__L2OutputOracle.address
  }
}

const depositViaOP = async (amount: NumberLike) => {
  console.log('Deposit Native token:', amount)
  console.log('Native token address:', l2NativeToken)

  const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
  const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

  const l2NativeTokenContract = new ethers.Contract(
    l2NativeToken,
    erc20ABI,
    l1Wallet
  )

  const l1Contracts = {
    AddressManager: addressManager,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }
  console.log('l1 contracts:', l1Contracts)

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const portals = new Portals({
    contracts: {
      l1: l1Contracts,
    },
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log('l2 native token balance in L1:', l2NativeTokenBalance.toString())

  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  const approveTx = await l2NativeTokenContract.approve(optimismPortal, amount)
  await approveTx.wait()
  console.log('approveTx:', approveTx.hash)

  const depositTx = await portals.depositTransaction({
    to: l2Wallet.address,
    value: BigNumber.from(amount),
    gasLimit: BigNumber.from('200000'),
    data: '0x',
  })
  const depositReceipt = await depositTx.wait()
  console.log('depositTx:', depositReceipt.transactionHash)

  const relayedTxHash = await portals.waitingDepositTransactionRelayed(
    depositReceipt,
    {}
  )
  const status = await portals.getMessageStatus(depositReceipt)
  console.log('deposit status relayed:', status === MessageStatus.RELAYED)
  console.log('relayedTxHash:', relayedTxHash)
  const depositedTxReceipt = await l2Provider.getTransactionReceipt(
    relayedTxHash
  )
  console.log('depositedTxReceipt:', depositedTxReceipt)

  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'l2 native token balance in L1: ',
    l2NativeTokenBalance.toString()
  )
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())
}

const withdrawViaBedrockMessagePasser = async (amount: NumberLike) => {
  console.log('Withdraw Native token:', amount)
  console.log('Native token address:', l2NativeToken)

  const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
  const l2Wallet = new ethers.Wallet(privateKey, asL2Provider(l2Provider))

  const l2NativeTokenContract = new ethers.Contract(
    l2NativeToken,
    erc20ABI,
    l1Wallet
  )

  const l1Contracts = {
    AddressManager: addressManager,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }
  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const portals = new Portals({
    contracts: {
      l1: l1Contracts,
    },
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log(
    'l2 native token balance before withdraw in L1: ',
    l2NativeTokenBalance.toString()
  )

  const withdrawalTx = await portals.initiateWithdrawal({
    target: l1Wallet.address,
    value: BigNumber.from(amount),
    gasLimit: BigNumber.from('200000'),
    data: '0x12345678',
  })
  const withdrawalReceipt = await withdrawalTx.wait()
  const withdrawalMessageInfo = await portals.calculateWithdrawalMessage(
    withdrawalReceipt
  )
  console.log('withdrawalMessageInfo:', withdrawalMessageInfo)

  let status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check publish L2\'s root:', status)

  await portals.waitForWithdrawalTxReadyForRelay(withdrawalReceipt)

  status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check ready for proving:', status)

  const proveTransaction = await portals.proveWithdrawalTransaction(
    withdrawalMessageInfo
  )
  await proveTransaction.wait()

  status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check in challenging:', status)

  await portals.waitForFinalization(withdrawalMessageInfo)
  const finalizedTransaction = await portals.finalizeWithdrawalTransaction(
    withdrawalMessageInfo
  )
  const finalizedTransactionReceipt = await finalizedTransaction.wait()
  console.log('finalized transaction receipt:', finalizedTransactionReceipt)

  status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check relayed:', status)

  const transferTx = await l2NativeTokenContract.transferFrom(
    l1Contracts.OptimismPortal,
    l1Wallet.address,
    amount
  )
  await transferTx.wait()
  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'l2 native token balance after withdraw in L1: ',
    l2NativeTokenBalance.toString()
  )
}

const withdrawViaBedrockMessagePasserV2 = async (amount: NumberLike) => {
  console.log('Withdraw Native token:', amount)
  console.log('Native token address:', l2NativeToken)

  const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
  const l2Wallet = new ethers.Wallet(privateKey, asL2Provider(l2Provider))

  const l2NativeTokenContract = new ethers.Contract(
    l2NativeToken,
    erc20ABI,
    l1Wallet
  )

  const l1Contracts = {
    AddressManager: addressManager,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }
  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const portals = new Portals({
    contracts: {
      l1: l1Contracts,
    },
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log(
    'l2 native token balance before withdraw in L1: ',
    l2NativeTokenBalance.toString()
  )

  const withdrawalTx = await portals.initiateWithdrawal({
    target: l1Wallet.address,
    value: BigNumber.from(amount),
    gasLimit: BigNumber.from('200000'),
    data: '0x12345678',
  })
  const withdrawalReceipt = await withdrawalTx.wait()

  let status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check publish L2\'s root:', status)

  await portals.waitForWithdrawalTxReadyForRelayUsingL2Tx(
    withdrawalReceipt.transactionHash
  )

  status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check ready for proving:', status)

  const proveTransaction = await portals.proveWithdrawalTransactionUsingL2Tx(
    withdrawalReceipt.transactionHash
  )
  await proveTransaction.wait()

  status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check in challenging:', status)

  await portals.waitForFinalizationUsingL2Tx(withdrawalReceipt.transactionHash)
  status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check ready for relay:', status)

  const finalizedTransaction =
    await portals.finalizeWithdrawalTransactionUsingL2Tx(
      withdrawalReceipt.transactionHash
    )
  const finalizedTransactionReceipt = await finalizedTransaction.wait()
  console.log('finalized transaction receipt:', finalizedTransactionReceipt)

  status = await portals.getL2ToL1MessageStatusByReceipt(withdrawalReceipt)
  console.log('[Withdrawal Status] check relayed:', status)

  const transferTx = await l2NativeTokenContract.transferFrom(
    l1Contracts.OptimismPortal,
    l1Wallet.address,
    amount
  )
  await transferTx.wait()
  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'l2 native token balance after withdraw in L1: ',
    l2NativeTokenBalance.toString()
  )
}

task('deposit-portal', 'Deposit L2NativeToken to L2 via OP.')
  .addParam('amount', 'Deposit amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositViaOP(args.amount)
  })

task(
  'withdraw-portal',
  'Withdraw L2NativeToken to L1 via BedrockMessagePasser.'
)
  .addParam('amount', 'Withdrawal amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await withdrawViaBedrockMessagePasser(args.amount)
  })

task(
  'withdraw-portal-v2',
  'Withdraw L2NativeToken to L1 via BedrockMessagePasser.'
)
  .addParam('amount', 'Withdrawal amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await withdrawViaBedrockMessagePasserV2(args.amount)
  })
