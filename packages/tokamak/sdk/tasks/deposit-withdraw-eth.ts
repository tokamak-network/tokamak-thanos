import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BytesLike, ethers } from 'ethers'
import { predeploys } from '@tokamak-network/core-utils'

import {
  CrossChainMessenger,
  ETHBridgeAdapter,
  MessageStatus,
  NumberLike,
  Portals,
} from '../src'

console.log('Setup task...')

const privateKey = process.env.PRIVATE_KEY as BytesLike

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

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

let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L1_OUTPUT_ORACLE || ''

const updateAddresses = async (hre: HardhatRuntimeEnvironment) => {
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

const depositETH = async (amount: NumberLike) => {
  console.log('Deposit ETH:', amount)
  console.log('l1 address:', l1Wallet.address)
  console.log('l2 address:', l2Wallet.address)

  const ethContract = new ethers.Contract(predeploys.ETH, erc20ABI, l2Wallet)

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
    ETH: {
      l1Bridge: l1Contracts.L1StandardBridge,
      l2Bridge: predeploys.L2StandardBridge,
      Adapter: ETHBridgeAdapter,
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

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 eth balance: ', l1Balance.toString())

  let l2Balance = await ethContract.balanceOf(l2Wallet.address)
  console.log('l2 eth balance:', l2Balance.toString())

  const depositTx = await messenger.depositETH(amount)
  await depositTx.wait()
  console.log('depositTx:', depositTx.hash)

  const portals = new Portals({
    contracts: {
      l1: l1Contracts,
    },
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  const relayedDepositTx =
    await portals.waitingDepositTransactionRelayedUsingL1Tx(depositTx.hash)
  console.log('relayed tx:', relayedDepositTx)

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 eth balance: ', l1Balance.toString())

  l2Balance = await ethContract.balanceOf(l2Wallet.address)
  console.log('l2 eth balance:', l2Balance.toString())
}

const withdrawETH = async (amount: NumberLike) => {
  console.log('Withdraw ETH:', amount)
  console.log('l1 address:', l1Wallet.address)
  console.log('l2 address:', l2Wallet.address)

  const ethContract = new ethers.Contract(predeploys.ETH, erc20ABI, l2Wallet)

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
    ETH: {
      l1Bridge: l1Contracts.L1StandardBridge,
      l2Bridge: predeploys.L2StandardBridge,
      Adapter: ETHBridgeAdapter,
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

  console.log('l1 eth balance: ', (await l1Wallet.getBalance()).toString())

  let l2Balance = await ethContract.balanceOf(l2Wallet.address)
  console.log('l2 eth balance: ', l2Balance.toString())

  const withdraw = await messenger.withdrawETH(amount)
  const withdrawalTx = await withdraw.wait()
  console.log(
    ' Withdrawal Tx:',
    withdrawalTx.transactionHash,
    ' Block',
    withdrawalTx.blockNumber,
    ' hash',
    withdraw.hash
  )

  l2Balance = await ethContract.balanceOf(l2Wallet.address)
  console.log('l2 eth balance: ', l2Balance.toString())

  await messenger.waitForMessageStatus(
    withdrawalTx,
    MessageStatus.READY_TO_PROVE
  )
  console.log('Prove the message')
  const proveTx = await messenger.proveMessage(withdrawalTx)
  const proveReceipt = await proveTx.wait(3)
  console.log('Proved the message:', proveReceipt.transactionHash)

  const finalizeInterval = setInterval(async () => {
    const currentStatus = await messenger.getMessageStatus(withdrawalTx)
    console.log('Message status:', currentStatus)
  }, 3000)

  try {
    await messenger.waitForMessageStatus(
      withdrawalTx,
      MessageStatus.READY_FOR_RELAY
    )
  } finally {
    clearInterval(finalizeInterval)
  }

  console.log('l1 eth balance:', (await l1Wallet.getBalance()).toString())
  const tx = await messenger.finalizeMessage(withdrawalTx)
  const receipt = await tx.wait()
  console.log('Finalized message tx', receipt.transactionHash)
  console.log('l1 eth balance:', (await l1Wallet.getBalance()).toString())
}

task('deposit-eth', 'Deposits native ETH to L2.')
  .addParam('amount', 'Deposit amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositETH(args.amount)
  })

task('withdraw-eth', 'Withdraw ERC20 ETH from L2.')
  .addParam('amount', 'Withdrawal amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await withdrawETH(args.amount)
  })
