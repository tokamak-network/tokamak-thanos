import { task, types } from 'hardhat/config'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys } from '@eth-optimism/core-utils'
import { BytesLike, ethers } from 'ethers'

import {
  CrossChainMessenger,
  MessageStatus,
  TONBridgeAdapter,
  AddressLike,
} from '../src'

console.log('Setup task...')

const privateKey = process.env.PRIVATE_KEY

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const l1Wallet = new ethers.Wallet(privateKey as BytesLike, l1Provider)
const l2Wallet = new ethers.Wallet(privateKey as BytesLike, l2Provider)

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

const ETH = '0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000'

const zeroAddr = '0x'.padEnd(42, '0')

let TON = process.env.TON
let addressManager = process.env.ADDRESS_MANAGER
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE
let optimismPortal = process.env.OPTIMISM_PORTAL
let l2OutputOracle = process.env.L1_OUTPUT_ORACLE

const updateAddress = async (hre) => {
  if (TON === undefined || TON === '') {
    const Deployment__TON = await hre.deployments.get('TON')
    TON = Deployment__TON.address
  }

  if (addressManager === undefined || addressManager === '') {
    const Deployment__AddressManager = await hre.deployments.get(
      'AddressManager'
    )
    addressManager = Deployment__AddressManager.address
  }

  if (l1CrossDomainMessenger === undefined || l1CrossDomainMessenger === '') {
    const Deployment__L1CrossDomainMessenger = await hre.deployments.get(
      'L1CrossDomainMessengerProxy'
    )
    l1CrossDomainMessenger = Deployment__L1CrossDomainMessenger.address
  }

  if (l1StandardBridge === undefined || l1StandardBridge === '') {
    const Deployment__L1StandardBridge = await hre.deployments.get(
      'L1StandardBridgeProxy'
    )
    l1StandardBridge = Deployment__L1StandardBridge.address
  }

  if (optimismPortal === undefined || optimismPortal === '') {
    const Deployment__OptimismPortal = await hre.deployments.get(
      'OptimismPortalProxy'
    )
    optimismPortal = Deployment__OptimismPortal.address
  }

  if (l2OutputOracle === undefined || l2OutputOracle === '') {
    const Deployment__L2OutputOracle = await hre.deployments.get(
      'L2OutputOracleProxy'
    )
    l2OutputOracle = Deployment__L2OutputOracle.address
  }
}

const depositTON = async (amount) => {
  console.log('Deposit TON:', amount)
  console.log('TON address:', TON)

  const tonContract = new ethers.Contract(TON as string, erc20ABI, l1Wallet)

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
    TON: {
      l1Bridge: l1Contracts.L1StandardBridge as AddressLike,
      l2Bridge: predeploys.L2StandardBridge as AddressLike,
      Adapter: TONBridgeAdapter,
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
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  const approveTx = await messenger.approveERC20(TON as string, ETH, amount)
  await approveTx.wait()
  console.log('approveTx:', approveTx.hash)

  const depositTx = await messenger.depositERC20(TON as string, ETH, amount)
  await depositTx.wait()
  console.log('depositTx:', depositTx.hash)

  await messenger.waitForMessageStatus(depositTx.hash, MessageStatus.RELAYED)

  l2Balance = await l2Wallet.getBalance()
  l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', l1TONBalance.toString())
  console.log('l2 native balance: ', l2Balance.toString())
}

const withdrawTON = async (amount) => {
  console.log('Withdraw TON:', amount)

  const tonContract = new ethers.Contract(TON as string, erc20ABI, l1Wallet)

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
    TON: {
      l1Bridge: l1Contracts.L1StandardBridge as AddressLike,
      l2Bridge: predeploys.L2StandardBridge as AddressLike,
      Adapter: TONBridgeAdapter,
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
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.toString())

  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  const withdrawal = await messenger.withdrawETH(amount)
  const withdrawalTx = await withdrawal.wait()
  console.log(
    'withdrawal Tx:',
    withdrawalTx.transactionHash,
    ' Block',
    withdrawalTx.blockNumber,
    ' hash',
    withdrawal.hash
  )

  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance:', l2Balance.toString())

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

  tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.toString())

  const tx = await messenger.finalizeMessage(withdrawalTx.transactionHash)
  const receipt = await tx.wait()
  console.log('Finalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  tonBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance: ', tonBalance.toString())
}

task('deposit-ton', 'Deposits ERC20-TON to L2.')
  .addParam('amount', 'Deposit amount', 1, types.int)
  .setAction(async (args, hre) => {
    await updateAddress(hre)
    await depositTON(args.amount)
  })

task('withdraw-ton', 'Withdraw native TON from L2.')
  .addParam('amount', 'Withdrawal amount', 1, types.int)
  .setAction(async (args, hre) => {
    await updateAddress(hre)
    await withdrawTON(args.amount)
  })
