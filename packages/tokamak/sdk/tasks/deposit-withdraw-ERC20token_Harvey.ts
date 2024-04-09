import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys } from '@tokamak-network/core-utils'
import OptimismMintableERC20TokenFactoryABI from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json'
import OptimismMintableERC20TokenABI from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'
import { BytesLike, ethers, Event } from 'ethers'

import { CrossChainMessenger, MessageStatus, NativeTokenBridgeAdapter, NumberLike } from '../src'
import L1StandardBridgeABI from '../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
import L2StandardBridgeABI from '../../contracts-bedrock/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'
import MockERC20ABI from '../../contracts-bedrock/forge-artifacts/MockERC20Token.sol/MockERC20Token.json'
// import * as OptimismMintableERC20TokenFactoryABI from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json'
// import * as OptimismMintableERC20TokenABI from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'
// import OptimismPortalABI from '../../contracts-bedrock/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'

// const OptimismPortalABI = require("../../contracts-bedrock/forge-artifacts/OptimismPortal.sol/OptimismPortal.json")
// const L1StandardBridgeABI = require("../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json")

console.log('Setup task...')

const privateKey = process.env.PRIVATE_KEY as BytesLike

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
console.log('l1Wallet :', l1Wallet.address)
const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

// const erc20ABI = [
//   {
//     inputs: [
//       { internalType: 'address', name: '_spender', type: 'address' },
//       { internalType: 'uint256', name: '_value', type: 'uint256' },
//     ],
//     name: 'approve',
//     outputs: [{ internalType: 'bool', name: '', type: 'bool' }],
//     stateMutability: 'nonpayable',
//     type: 'function',
//   },
//   {
//     constant: true,
//     inputs: [{ name: '_owner', type: 'address' }],
//     name: 'balanceOf',
//     outputs: [{ name: 'balance', type: 'uint256' }],
//     type: 'function',
//   },
//   {
//     inputs: [{ internalType: 'uint256', name: 'amount', type: 'uint256' }],
//     name: 'faucet',
//     outputs: [],
//     stateMutability: 'nonpayable',
//     type: 'function',
//   },
// ]

// const ETH = '0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000'

const oneETH = ethers.utils.parseUnits('1', 18)
const twoETH = ethers.utils.parseUnits('2', 18)

const zeroAddr = '0x'.padEnd(42, '0')

let l2NativeToken = process.env.NATIVE_TOKEN || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''
let MockERC20
let l2MockERC20
let l1MockAddress
let l2MockAddress

const name = 'Mock'
const symbol = 'MTK'

const l2name = 'L2Mock'
const l2symbol = 'LTK'

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

const depositERC20Token = async (amount: NumberLike) => {
  console.log('Deposit Native token:', amount)
  console.log('Native token address:', l2NativeToken)

  const factory_MockERC20 = new ethers.ContractFactory(
    MockERC20ABI.abi,
    MockERC20ABI.bytecode,
    l1Wallet
  )

  MockERC20 = await factory_MockERC20.deploy(name,symbol)
  await MockERC20.deployed()

  await MockERC20.mint(l1Wallet.address, amount)

  const factory_OptimismMintable = new ethers.Contract(
    predeploys.OptimismMintableERC20Factory,
    OptimismMintableERC20TokenFactoryABI.abi,
    l2Wallet
  )

  let tx = await factory_OptimismMintable.createOptimismMintableERC20(
    MockERC20.address,
    l2name,
    l2symbol
  )
  await tx.wait()

  const receipt = await tx.wait()
  const event = receipt.events.find(
    (e: Event) => e.event === 'OptimismMintableERC20Created'
  )

  if (!event) {
    throw new Error('Unable to find OptimismMintableERC20Created event')
  }

  l2MockAddress = event.args.localToken
  console.log('l1MockAddress:', MockERC20.address)
  console.log('l2MockAddress:', l2MockAddress)

  l2MockERC20 = new ethers.Contract(
    l2MockAddress,
    OptimismMintableERC20TokenABI.abi,
    l2Wallet
  )

  // const l2NativeTokenContract = new ethers.Contract(
  //   l2NativeToken,
  //   erc20ABI,
  //   l1Wallet
  // )

  // const OptimismPortalContract = new ethers.Contract(
  //   optimismPortal,
  //   OptimismPortalABI.abi,
  //   l1Wallet
  // )

  const L1StandardBridgeContract = new ethers.Contract(
    l1StandardBridge,
    L1StandardBridgeABI.abi,
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
  // console.log('l1 contracts:', l1Contracts)

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

  let l1MockBalance = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance: ', l1MockBalance.toString())

  let l2MockBalance = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance: ', l2MockBalance.toString())

  // let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
  //   l1Wallet.address
  // )
  // console.log('native token(TON) balance in L1:', Number(l2NativeTokenBalance.toString()))

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  // let optimismPortalStorage = await OptimismPortalContract.depositedAmount()
  // console.log('optimismPortalStorage: ', optimismPortalStorage.toString())

  let standardStorage = await L1StandardBridgeContract.deposits(MockERC20.address,l2MockAddress)
  console.log('standardStorage: ', standardStorage.toString())

  // if (Number(l2NativeTokenBalance.toString()) === 0) {
  //   console.log('start faucet')
  //   const tx = await l2NativeTokenContract.connect(l1Wallet).faucet(amount)
  //   await tx.wait()
  //   const l2NativeTokenBalance2 = await l2NativeTokenContract.balanceOf(
  //     l1Wallet.address
  //   )
  //   console.log('after faucet l2 native token balance in L1:', l2NativeTokenBalance2.toString())
  // }

  tx = await MockERC20.connect(l1Wallet).approve(L1StandardBridgeContract.address, amount)
  await tx.wait()

  tx = await L1StandardBridgeContract.connect(l1Wallet).depositERC20(
    MockERC20.address,
    l2MockAddress,
    amount,
    20000,
    '0x'
  )
  const depositTx = await tx.wait()
  console.log(
    'depositTx Tx:',
    depositTx.transactionHash,
    ' Block',
    depositTx.blockNumber,
    ' hash',
    tx.hash
  )

  await messenger.waitForMessageStatus(depositTx.transactionHash, MessageStatus.RELAYED)

  l1MockBalance = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance: ', l1MockBalance.toString())
  l2MockBalance = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance: ', l2MockBalance.toString())

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())
  // l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  // console.log(
  //   'native token(TON) balance in L1: ',
  //   l2NativeTokenBalance.toString()
  // )

  // optimismPortalStorage = await OptimismPortalContract.depositedAmount()
  // console.log('optimismPortalStorage: ', optimismPortalStorage.toString())

  standardStorage = await L1StandardBridgeContract.deposits(MockERC20.address,l2MockAddress)
  console.log('standardStorage: ', standardStorage.toString())
}

const withdrawERC20Token = async (amount: NumberLike) => {
  console.log('Withdraw Native token:', amount)

  // const l2NativeTokenContract = new ethers.Contract(
  //   l2NativeToken,
  //   erc20ABI,
  //   l1Wallet
  // )

  // const OptimismPortalContract = new ethers.Contract(
  //   optimismPortal,
  //   OptimismPortalABI.abi,
  //   l1Wallet
  // )

  const L1StandardBridgeContract = new ethers.Contract(
    l1StandardBridge,
    L1StandardBridgeABI.abi,
    l1Wallet
  )

  const L2StandardBridgeContract = new ethers.Contract(
    predeploys.L2StandardBridge,
    L2StandardBridgeABI.abi,
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

  l1MockAddress = '0xBD20Cd677066ced4099A479feBD2A416686917BB'
  l2MockAddress = '0x2957d6E10BC45CcE14c4FfC256dCb11E267d2D06'

  MockERC20 = new ethers.Contract(
    l1MockAddress,
    MockERC20ABI.abi,
    l1Wallet
  )

  l2MockERC20 = new ethers.Contract(
    l2MockAddress,
    OptimismMintableERC20TokenABI.abi,
    l2Wallet
  )

  let l1MockBalance = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance: ', l1MockBalance.toString())
  let l2MockBalance = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance: ', l2MockBalance.toString())

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  let standardStorage = await L1StandardBridgeContract.deposits(MockERC20.address,l2MockERC20.address)
  console.log('standardStorage: ', standardStorage.toString())

  const withdrawal = await L2StandardBridgeContract.connect(l2Wallet).withdraw(
    l2MockERC20.address,
    amount,
    20000,
    '0x'
  )
  const withdrawalTx = await withdrawal.wait()
  console.log(
    'withdrawal Tx:',
    withdrawalTx.transactionHash,
    ' Block',
    withdrawalTx.blockNumber,
    ' hash',
    withdrawal.hash
  )

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

  const tx = await messenger.finalizeMessage(withdrawalTx.transactionHash)
  const receipt = await tx.wait()
  console.log('Finalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  l1MockBalance = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance: ', l1MockBalance.toString())
  l2MockBalance = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance: ', l2MockBalance.toString())

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance:', l2Balance.toString())

  standardStorage = await L1StandardBridgeContract.deposits(MockERC20.address,l2MockERC20.address)
  console.log('standardStorage: ', standardStorage.toString())
}

task('deposit-erc20-token', 'Deposits ERC20Token to L2.')
  .addParam('amount', 'Deposit amount', twoETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositERC20Token(args.amount)
  })

task('withdraw-erc20-token', 'Withdraw ERC20token from L2.')
  .addParam('amount', 'Withdrawal amount', oneETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await withdrawERC20Token(args.amount)
  })
