import '@nomiclabs/hardhat-ethers'
import { predeploys, sleep } from '@tokamak-network/core-utils'
import { BytesLike, ethers } from 'ethers'
import 'hardhat-deploy'
import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'

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

let l2NativeToken = process.env.NATIVE_TOKEN || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''

const updateAddresses = async (hre: HardhatRuntimeEnvironment) => {
  if (l2NativeToken === '') {
    const Deployment__l2NativeToken = await hre.deployments.get('L2NativeToken')
    l2NativeToken = Deployment__l2NativeToken.address
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
  console.log('TON address:', l2NativeToken)

  const tonContract = new ethers.Contract(l2NativeToken, erc20ABI, l1Wallet)

  const l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
    L1UsdcBridge: zeroAddr,
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
    nativeTokenAddress: l2NativeToken,
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
  console.log('l2 native balance before depositing: ', l2BalancePrev.toString())

  const data = ethers.utils.solidityPack(
    ['address', 'uint32', 'bytes'],
    [l2Wallet.address, 200000, '0x']
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
  console.log('TON address:', l2NativeToken)

  const tonContract = new ethers.Contract(l2NativeToken, erc20ABI, l1Wallet)

  const l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
    L1UsdcBridge: zeroAddr,
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
    nativeTokenAddress: l2NativeToken,
    bridges,
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  let l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  const l2TONBalanceBefore = await l2Wallet.getBalance()

  console.log('l2 ton balance: ', l2TONBalanceBefore.toString())

  const data = ethers.utils.solidityPack(
    ['address', 'uint32', 'bytes'],
    [l2Wallet.address, 200000, '0xd0e30db0']
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

  l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  const l2TONBalanceAfter = await l2Wallet.getBalance()
  console.log('l1 ton balance after: ', l1TONBalance.toString())
  console.log('l2 ton balance after: ', l2TONBalanceAfter.toString())

  console.log(
    'added ton balance: ',
    l2TONBalanceAfter.sub(l2TONBalanceBefore).toString()
  )
}

const approveAndDepositTONViaOP = async (amount: NumberLike) => {
  console.log('Deposit TON via Portal:', amount)
  console.log('TON address:', l2NativeToken)

  const tonContract = new ethers.Contract(l2NativeToken, erc20ABI, l1Wallet)
  const LegacyERC20NativeToken = new ethers.Contract(
    predeploys.LegacyERC20NativeToken,
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
    L1UsdcBridge: zeroAddr,
  }
  console.log('l1 contracts:', l1Contracts)

  let l1TONBalance = await tonContract.balanceOf(l1Wallet.address)
  console.log('l1 ton balance:', l1TONBalance.toString())

  const l2BalancePrev = await LegacyERC20NativeToken.balanceOf(l1Wallet.address)
  console.log('l2 wton balance: ', l2BalancePrev.toString())

  const data = ethers.utils.solidityPack(
    ['address', 'uint256', 'uint32', 'bytes'],
    [l1Wallet.address, amount, 200000, '0xd0e30db0']
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
    const l2BalanceAfter = await LegacyERC20NativeToken.balanceOf(
      l1Wallet.address
    )
    if (l2BalanceAfter.eq(l2BalancePrev)) {
      await sleep(1000)
      continue
    }
    console.log(
      'l2 LegacyERC20NativeToken balance: ',
      l2BalanceAfter.toString()
    )
    console.log(
      'added LegacyERC20NativeToken balance: ',
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
