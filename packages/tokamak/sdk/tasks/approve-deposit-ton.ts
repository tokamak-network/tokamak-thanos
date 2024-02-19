import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys, sleep } from '@tokamak-network/core-utils'
import { BytesLike, ethers } from 'ethers'

import {
  CrossChainMessenger,
  MessageStatus,
  NativeTokenBridgeAdapter,
  NumberLike,
} from '../src'

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
    inputs: [{ name: 'account', type: 'address' }],
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
    inputs: [],
    name: 'totalSupply',
    outputs: [
      {
        internalType: 'uint256',
        name: '',
        type: 'uint256',
      },
    ],
    stateMutability: 'view',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'address',
        name: 'spender',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'amount',
        type: 'uint256',
      },
      {
        internalType: 'bytes',
        name: 'data',
        type: 'bytes',
      },
    ],
    name: 'approveAndCall',
    outputs: [
      {
        internalType: 'bool',
        name: '',
        type: 'bool',
      },
    ],
    stateMutability: 'nonpayable',
    type: 'function',
  },
]

const zeroAddr = '0x'.padEnd(42, '0')

let TON = process.env.TON || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''

const updateAddresses = async (hre: HardhatRuntimeEnvironment) => {
  if (TON === '') {
    const Deployment__TON = await hre.deployments.get('L2NativeToken')
    TON = Deployment__TON.address
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

const approveAndDepositTON = async (amount: NumberLike) => {
  console.log('Deposit TON:', amount)
  console.log('TON address:', TON)

  const tonContract = new ethers.Contract(TON, erc20ABI, l1Wallet)

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
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  // let receipt = await (await tonContract.connect(l1Wallet).faucet(ethers.BigNumber.from(""+amount)))
  // console.log('faucet receipt', receipt)
  // let l1TONTotalSupply = await tonContract.totalSupply()
  // console.log('l1 ton total supply:', l1TONTotalSupply.toString())

  let l1TONBalance = await tonContract
    .connect(l1Wallet)
    .balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  const l2BalancePrev = await l2Wallet.getBalance()
  console.log('l2 native balance prev: ', l2BalancePrev.toString())

  const data = ethers.utils.solidityPack(
    ['address', 'address', 'uint256', 'uint32', 'bytes'],
    [l1Wallet.address, l1Wallet.address, amount, 200000, '0x']
  )
  const approveAndCallTx = await (
    await tonContract
      .connect(l1Wallet)
      .approveAndCall(
        l1Contracts.L1StandardBridge,
        ethers.BigNumber.from('' + amount),
        data
      )
  ).wait()
  console.log('approveAndCallTx:', approveAndCallTx.transactionHash)

  await messenger.waitForMessageStatus(
    approveAndCallTx.transactionHash,
    MessageStatus.RELAYED
  )

  const l2BalanceAfter = await l2Wallet.getBalance()
  l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance after: ', l1TONBalance.toString())
  console.log('l2 native balance: ', l2BalanceAfter.toString())

  console.log(
    'l2 added native balance: ',
    l2BalanceAfter.sub(l2BalancePrev).toString()
  )
}

const approveAndDepositTONViaCDM = async (amount: NumberLike) => {
  console.log('Deposit TON via CDM:', amount)
  console.log('TON address:', TON)

  const tonContract = new ethers.Contract(TON, erc20ABI, l1Wallet)
  const wtonContract = new ethers.Contract(
    predeploys.WNativeToken,
    erc20ABI,
    l2Wallet
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
    TON: {
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
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  const l2CDMBalancePrev = await wtonContract.balanceOf(
    predeploys.L2CrossDomainMessenger
  )
  console.log('l2cdm wton balance: ', l2CDMBalancePrev.toString())

  const data = ethers.utils.solidityPack(
    ['address', 'address', 'uint256', 'uint32', 'bytes'],
    [l1Wallet.address, predeploys.WNativeToken, amount, 200000, '0xd0e30db0']
  )

  console.log('Approve and Call via CDM: ', data)
  const approveAndCallTx = await (
    await tonContract
      .connect(l1Wallet)
      .approveAndCall(
        l1Contracts.L1CrossDomainMessenger,
        ethers.BigNumber.from('' + amount),
        data
      )
  ).wait()
  console.log('approveAndCallTx:', approveAndCallTx.transactionHash)

  await messenger.waitForMessageStatus(
    approveAndCallTx.transactionHash,
    MessageStatus.RELAYED
  )

  const l2CDMBalanceAfter = await wtonContract.balanceOf(
    predeploys.L2CrossDomainMessenger
  )
  l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance after: ', l1TONBalance.toString())
  console.log('l2cdm wton balance: ', l2CDMBalanceAfter.toString())

  console.log(
    'added wton balance: ',
    l2CDMBalanceAfter.sub(l2CDMBalancePrev).toString()
  )
}

const approveAndDepositTONViaOP = async (amount: NumberLike) => {
  console.log('Deposit TON via Portal:', amount)
  console.log('TON address:', TON)

  const tonContract = new ethers.Contract(TON, erc20ABI, l1Wallet)
  const wtonContract = new ethers.Contract(
    predeploys.WNativeToken,
    erc20ABI,
    l2Wallet
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

  let l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  const l2BalancePrev = await wtonContract.balanceOf(l1Wallet.address)
  console.log('l2 wton balance: ', l2BalancePrev.toString())

  const data = ethers.utils.solidityPack(
    ['address', 'address', 'uint256', 'uint32', 'bytes'],
    [l1Wallet.address, predeploys.WNativeToken, amount, 200000, '0xd0e30db0']
  )

  console.log('Approve and Call via Portal: ', data)
  const approveAndCallTx = await (
    await tonContract
      .connect(l1Wallet)
      .approveAndCall(optimismPortal, ethers.BigNumber.from(amount), data)
  ).wait()
  console.log('approveAndCallTx:', approveAndCallTx.transactionHash)
  l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance after: ', l1TONBalance.toString())
  while (true) {
    const l2BalanceAfter = await wtonContract.balanceOf(l1Wallet.address)
    if (l2BalanceAfter.eq(l2BalancePrev)) {
      await sleep(1000)
      continue
    }
    console.log('l2 wton balance: ', l2BalanceAfter.toString())
    console.log(
      'added wton balance: ',
      l2BalanceAfter.sub(l2BalancePrev).toString()
    )
    break
  }
}

task('approve-deposit-ton', 'Deposits ERC20-TON to L2.')
  .addParam('amount', 'Deposit amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await approveAndDepositTON(args.amount)
  })

task('approve-deposit-ton-cdm', 'Deposits ERC20-TON to L2.')
  .addParam('amount', 'Deposit amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await approveAndDepositTONViaCDM(args.amount)
  })

task('approve-deposit-ton-portal', 'Deposits ERC20-TON to L2.')
  .addParam('amount', 'Deposit amount', '1', types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await approveAndDepositTONViaOP(args.amount)
  })
