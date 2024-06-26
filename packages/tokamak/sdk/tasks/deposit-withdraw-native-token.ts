import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys } from '@tokamak-network/core-utils'
import { BytesLike, ethers } from 'ethers'
import Artifact__L1CrossDomainMessenger from '@tokamak-network/thanos-contracts/forge-artifacts/L1CrossDomainMessenger.sol/L1CrossDomainMessenger.json'

import {
  CrossChainMessenger,
  MessageStatus,
  NativeTokenBridgeAdapter,
  NumberLike,
  asL2Provider,
} from '../src'

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
]

const zeroAddr = '0x'.padEnd(42, '0')

let nativeTokenAddress = process.env.NATIVE_TOKEN || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''

const updateAddresses = async (hre: HardhatRuntimeEnvironment) => {
  if (nativeTokenAddress === '') {
    const Deployment__L2NativeToken = await hre.deployments.get('L2NativeToken')
    nativeTokenAddress = Deployment__L2NativeToken.address
  }

  if (addressManager === '') {
    const Deployment__AddressManager = await hre.deployments.get(
      'AddressManager'
    )
    addressManager = Deployment__AddressManager.address
  }

  if (l1CrossDomainMessenger === '') {
    const Deployment__L1CrossDomainMessenger = await hre.deployments.get(
      'L1CrossDomainMessengerProxy'
    )
    l1CrossDomainMessenger = Deployment__L1CrossDomainMessenger.address
  }

  if (l1StandardBridge === '') {
    const Deployment__L1StandardBridge = await hre.deployments.get(
      'L1StandardBridgeProxy'
    )
    l1StandardBridge = Deployment__L1StandardBridge.address
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

const depositNativeToken = async (amount: NumberLike) => {
  console.log('Deposit Native token:', amount)
  console.log('Native token address:', nativeTokenAddress)

  const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
  const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

  const l2NativeTokenContract = new ethers.Contract(
    nativeTokenAddress,
    erc20ABI,
    l1Provider
  )

  const l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }
  console.log('L1 contracts:', l1Contracts)

  const bridges = {
    NativeToken: {
      l1Bridge: l1Contracts.L1StandardBridge,
      l2Bridge: predeploys.L2StandardBridge,
      Adapter: NativeTokenBridgeAdapter,
    },
  }

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const messenger = new CrossChainMessenger({
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
    bridges,
    l1ChainId,
    l2ChainId,
    nativeTokenAddress,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log(
    'L2 native token balance in L1 before depositing:',
    l2NativeTokenBalance.toString()
  )

  let l2Balance = await l2Wallet.getBalance()
  console.log('Balance in L2 before depositing: ', l2Balance.toString())

  const approveTx = await messenger.approveNativeToken(amount)
  await approveTx.wait()
  console.log('Approve tx:', approveTx.hash)

  const depositTx = await messenger.depositNativeToken(amount)
  await depositTx.wait()
  console.log(
    'Deposit native token from L1 to L2, tx_hash:',
    depositTx.hash,
    'amount:',
    amount
  )

  await messenger.waitForMessageStatus(depositTx.hash, MessageStatus.RELAYED)

  console.log('Verify balances')
  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'L2 native token balance in L1 after depositing:',
    l2NativeTokenBalance.toString()
  )
  l2Balance = await l2Wallet.getBalance()
  console.log('Balance in L2 after depositing: ', l2Balance.toString())
}

const depositNativeTokenViaMessenger = async (amount: NumberLike) => {
  console.log('Deposit Native token:', amount)
  console.log('Native token address:', nativeTokenAddress)

  const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
  const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

  const l1CrossDomainMessengerContract = new ethers.Contract(
    l1CrossDomainMessenger,
    Artifact__L1CrossDomainMessenger.abi,
    l1Wallet
  )

  const l2NativeTokenContract = new ethers.Contract(
    nativeTokenAddress,
    erc20ABI,
    l1Wallet
  )

  const l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }
  console.log('l1 contracts:', l1Contracts)

  const bridges = {
    NativeToken: {
      l1Bridge: l1Contracts.L1StandardBridge,
      l2Bridge: predeploys.L2StandardBridge,
      Adapter: NativeTokenBridgeAdapter,
    },
  }

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const messenger = new CrossChainMessenger({
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
    bridges,
    l1ChainId,
    l2ChainId,
    nativeTokenAddress,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log('L2 native token balance in L1:', l2NativeTokenBalance.toString())

  let l2Balance = await l2Wallet.getBalance()
  console.log('L2 native token balance in L2: ', l2Balance.toString())

  const approveTx = await l2NativeTokenContract.approve(
    l1CrossDomainMessenger,
    amount
  )
  await approveTx.wait()
  console.log('Approve transaction hash:', approveTx.hash)

  const depositTx = await l1CrossDomainMessengerContract.sendNativeTokenMessage(
    l2Wallet.address,
    amount,
    '0x',
    21000
  )
  await depositTx.wait()
  console.log('Deposit transaction hash:', depositTx.hash)

  await messenger.waitForMessageStatus(depositTx.hash, MessageStatus.RELAYED)

  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'L2 native token balance in L1: ',
    l2NativeTokenBalance.toString()
  )
  l2Balance = await l2Wallet.getBalance()
  console.log('L2 native balance in L2: ', l2Balance.toString())
}

const withdrawNativeToken = async (amount: NumberLike) => {
  console.log('Withdraw Native token:', amount)

  const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
  const l2Wallet = new ethers.Wallet(privateKey, asL2Provider(l2Provider))

  const l2NativeTokenContract = new ethers.Contract(
    nativeTokenAddress,
    erc20ABI,
    l1Wallet
  )

  const l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }

  const bridges = {
    NativeToken: {
      l1Bridge: l1Contracts.L1StandardBridge,
      l2Bridge: predeploys.L2StandardBridge,
      Adapter: NativeTokenBridgeAdapter,
    },
  }

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

  const messenger = new CrossChainMessenger({
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
    bridges,
    l1ChainId,
    l2ChainId,
    nativeTokenAddress,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log(
    'L2 native token: balance in L1: ',
    l2NativeTokenBalance.toString()
  )

  const l2Balance = await l2Wallet.getBalance()
  console.log('L2 native token: balance on L2: ', l2Balance.toString())

  const withdrawal = await messenger.withdrawNativeToken(amount)
  const withdrawalTx = await withdrawal.wait()

  const updatedL2Balance = await l2Wallet.getBalance()
  console.log(
    'L2 native token: balance after depositing',
    updatedL2Balance.toString()
  )

  const l1CostActual = withdrawalTx['l1Fee']
  const l2CostActual = withdrawalTx.gasUsed.mul(withdrawalTx.effectiveGasPrice)

  console.log(`L1 cost actual: ${l1CostActual.toString()}`)
  console.log(`L2 cost actual: ${l2CostActual.toString()}`)
  console.log(`Total cost: ${l1CostActual.add(l2CostActual).toString()}`)
  console.log('Withdrawal amount', amount.toString())
  console.log(
    `Spent amount = L1 cost actual + L2 cost actual + Withdrawal Amount = ${l1CostActual
      .add(l2CostActual)
      .add(amount)
      .toString()}`
  )
  console.log('Balance changed ', l2Balance.sub(updatedL2Balance).toString())

  // // Check ready for prove
  await messenger.waitForMessageStatus(
    withdrawalTx.transactionHash,
    MessageStatus.READY_TO_PROVE
  )

  console.log('Prove the message')
  const proveTx = await messenger.proveMessage(withdrawalTx.transactionHash)
  const proveReceipt = await proveTx.wait(3)
  console.log('Proved the message: ', proveReceipt.transactionHash)

  const finalizeInterval = setInterval(async () => {
    const currentStatus = await messenger.getMessageStatus(
      withdrawalTx.transactionHash
    )
    console.log('Message status: ', currentStatus)
  }, 3000)

  try {
    await messenger.waitForMessageStatus(
      withdrawalTx.transactionHash,
      MessageStatus.READY_FOR_RELAY
    )
  } finally {
    clearInterval(finalizeInterval)
  }

  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'L2 native token: balance in L1 after proving: ',
    l2NativeTokenBalance.toString()
  )

  const tx = await messenger.finalizeMessage(withdrawalTx.transactionHash)
  const receipt = await tx.wait()
  console.log('Finalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'L2 native token: balance in L1 after finalizing: ',
    l2NativeTokenBalance.toString()
  )
}

task('deposit-native-token', 'Deposits L2NativeToken to L2.')
  .addParam('amount', 'Deposit amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositNativeToken(args.amount)
  })

task('deposit-native-token-via-messenger', 'Deposits L2NativeToken to L2.')
  .addParam('amount', 'Deposit amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositNativeTokenViaMessenger(args.amount)
  })

task('withdraw-native-token', 'Withdraw native token from L2.')
  .addParam('amount', 'Withdrawal amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await withdrawNativeToken(args.amount)
  })
