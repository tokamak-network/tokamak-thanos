import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys } from '@eth-optimism/core-utils'
import { BytesLike, ethers } from 'ethers'

import { CrossChainMessenger, MessageStatus, NativeTokenBridgeAdapter, NumberLike } from '../src'
// import MockERC20ABI from '../../contracts-bedrock/forge-artifacts/MockERC20Token.sol/MockERC20Token.json'
import L1FastWithdrawABI from '../../contracts-bedrock/forge-artifacts/L1FastWithdraw.sol/L1FastWithdraw.json'
import L2FastWithdrawABI from '../../contracts-bedrock/forge-artifacts/L2FastWithdraw.sol/L2FastWithdraw.json'
import L1StandardBridgeABI from '../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
// import OptimismPortalABI from '../../contracts-bedrock/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'
// import L2StandardBridgeABI from '../../contracts-bedrock/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'

// import * as OptimismPortalABI from '../../contracts-bedrock/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'
// import * as L1StandardBridgeABI from '../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
// import L1FastWithdrawProxyABI from '../../contracts-bedrock/forge-artifacts/L1FastWithdrawProxy.sol/L1FastWithdrawProxy.json'
// import L2FastWithdrawProxyABI from '../../contracts-bedrock/forge-artifacts/L2FastWithdrawProxy.sol/L2FastWithdrawProxy.json'

const privateKey = process.env.PRIVATE_KEY as BytesLike
const privateKey2 = process.env.PRIVATE_KEY2 as BytesLike

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
// console.log('l1Wallet :', l1Wallet.address)
const l1user1 = new ethers.Wallet(privateKey2, l1Provider)
// console.log('l1user1 :', l1user1.address)
const l2Wallet = new ethers.Wallet(privateKey, l2Provider)
// console.log('l2Wallet :', l2Wallet.address)
const l2user1 = new ethers.Wallet(privateKey2, l2Provider)
// console.log('l2user1 :', l2user1.address)

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
    inputs: [
      {
        internalType: 'address',
        name: 'from',
        type: 'address',
      },
      {
        internalType: 'address',
        name: 'to',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'amount',
        type: 'uint256',
      }
    ],
    name: 'transferFrom',
    outputs: [
      {
        internalType: 'bool',
        name: '',
        type: 'bool',
      }
    ],
    stateMutability: 'nonpayable',
    type: 'function',
  },
]

// const ETH = '0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000'
// const ETH = '0x0000000000000000000000000000000000000000'

const oneETH = ethers.utils.parseUnits('1', 18)
const twoETH = ethers.utils.parseUnits('2', 18)
const threeETH = ethers.utils.parseUnits('3', 18)
// const fourETH = ethers.utils.parseUnits('4', 18)
const fiveETH = ethers.utils.parseUnits('5', 18)

// const eightETH = ethers.utils.parseUnits('8', 18)
// const tenETH = ethers.utils.parseUnits('10', 18)

const zeroAddr = '0x'.padEnd(42, '0')

let l2NativeToken = process.env.NATIVE_TOKEN || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''

// let L1FastWithDrawProxy
// let L1FastWithDraw
let L1FastWithDrawContract
// let L2FastWithDrawProxy
// let L2FastWithDraw
let L2FastWithDrawContract

// let l1fastWithdrawAddr = ""
// let l2fastWithdrawAddr = ""
const l2CrossDomainMessengerAddr = '0x4200000000000000000000000000000000000007'

let l2ETHERC20

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

const depositETH = async (amount: NumberLike) => {
  console.log('Deposit Native token:', amount)
  console.log('Native token address:', l2NativeToken)

  l2ETHERC20 = new ethers.Contract(
    predeploys.ETH,
    erc20ABI,
    l2Wallet
  )

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

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance(ETH): ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance(TON): ', l2Balance.toString())
  let l2ETHBalance = await l2ETHERC20.balanceOf(l2Wallet.address)
  console.log('l2 ETH(ERC20) balance: ', l2ETHBalance.toString())

  const tx = await L1StandardBridgeContract.connect(l1Wallet).depositETH(
    20000,
    '0x',
    {
      value: amount,
    }
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

  await messenger.waitForMessageStatus(
    depositTx.transactionHash,
    MessageStatus.RELAYED
  )

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance(ETH): ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance(TON): ', l2Balance.toString())
  l2ETHBalance = await l2ETHERC20.balanceOf(l2Wallet.address)
  console.log('l2 ETH(ERC20) balance: ', l2ETHBalance.toString())
}

const fastwithdrawERC20Token = async (amount: NumberLike) => {
  console.log('Withdraw Native token:', amount)

  l2ETHERC20 = new ethers.Contract(
    predeploys.ETH,
    erc20ABI,
    l2Wallet
  )

  const l2NativeTokenContract = new ethers.Contract(
    l2NativeToken,
    erc20ABI,
    l1Wallet
  )

  // const OptimismPortalContract = new ethers.Contract(
  //   optimismPortal,
  //   OptimismPortalABI.abi,
  //   l1Wallet
  // )

  // const L1StandardBridgeContract = new ethers.Contract(
  //   l1StandardBridge,
  //   L1StandardBridgeABI.abi,
  //   l1Wallet
  // )

  // const L2StandardBridgeContract = new ethers.Contract(
  //   predeploys.L2StandardBridge,
  //   L2StandardBridgeABI.abi,
  //   l2Wallet
  // )

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

  // deploy the L1, L2FastWithdraw
  const L1FastWithDrawDep = new ethers.ContractFactory(
    L1FastWithdrawABI.abi,
    L1FastWithdrawABI.bytecode,
    l1Wallet
  )

  L1FastWithDrawContract = await L1FastWithDrawDep.deploy()
  await L1FastWithDrawContract.deployed()
  // console.log('L1FastWithDrawContract.address :', L1FastWithDrawContract.address)

  const L2FastWithDrawDep = new ethers.ContractFactory(
    L2FastWithdrawABI.abi,
    L2FastWithdrawABI.bytecode,
    l2Wallet
  )

  L2FastWithDrawContract = await L2FastWithDrawDep.deploy()
  await L2FastWithDrawContract.deployed()
  // console.log(await l2Provider.getCode(L2FastWithDrawContract.address))
  // console.log('L2FastWithDrawContract.address :', L2FastWithDrawContract.address)

  //L1, L2 initialize
  await (await L1FastWithDrawContract.connect(l1Wallet).initialize(
    l1Contracts.L1CrossDomainMessenger,
    L2FastWithDrawContract.address,
    zeroAddr,
    l2NativeTokenContract.address
  )).wait()
  console.log('L1FastWithdraw initialize done')

  const checkL1Inform = await L1FastWithDrawContract.crossDomainMessenger()
  console.log('checkL1Inform :', checkL1Inform)
  console.log('l1Contracts.L1CrossDomainMessenger :', l1Contracts.L1CrossDomainMessenger)

  await (await L2FastWithDrawContract.connect(l2Wallet).initialize(
    l2CrossDomainMessengerAddr,
    L1FastWithDrawContract.address,
    predeploys.LegacyERC20ETH,
    l2NativeTokenContract.address
  )).wait();
  console.log('L2FastWithdraw initialize done')

  const checkL2Inform = await L2FastWithDrawContract.crossDomainMessenger()
  console.log('checkL2Inform :', checkL2Inform)
  console.log('l2CrossDomainMessengerAddr :', l2CrossDomainMessengerAddr)

  // start the test

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance (ETH) (l1Wallet): ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance (TON) (l1Wallet): ', l2Balance.toString())
  let l2ETHBalance = await l2ETHERC20.balanceOf(l2Wallet.address)
  console.log('l2 ETH(ERC20) balance (l1Wallet): ', l2ETHBalance.toString())

  let l1BalanceUser1 = await l1user1.getBalance()
  console.log('l1 native balance (ETH) (User1): ', l1BalanceUser1.toString())
  let l2BalanceUser1 = await l2user1.getBalance()
  console.log('l2 native balance (TON) (User1): ', l2BalanceUser1.toString())
  let l2ETHBalanceUser1 = await l2ETHERC20.balanceOf(l2user1.address)
  console.log('l2 ETH(ERC20) balance (User1): ', l2ETHBalanceUser1.toString())

  //request L2
  let L2FastWithdrawBalance = await l2ETHERC20.balanceOf(L2FastWithDrawContract.address)
  console.log('before L2 ERC20 (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  const tx = await l2ETHERC20.connect(l2Wallet).approve(L2FastWithDrawContract.address, threeETH)
  await tx.wait()
  console.log('pass the approve')

  await (await L2FastWithDrawContract.connect(l2Wallet).requestFW(
    l2ETHERC20.address,
    threeETH,
    twoETH
  )).wait()
  console.log('pass the request')

  const saleCount = await L2FastWithDrawContract.salecount()
  console.log('saleCount : ', saleCount)
  let saleInformation = await L2FastWithDrawContract.dealData(saleCount)
  console.log('saleInformation : ', saleInformation)

  L2FastWithdrawBalance = await l2ETHERC20.balanceOf(L2FastWithDrawContract.address)
  console.log('after L2 ERC20 (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  //provider L1
  const providerTx = await L1FastWithDrawContract.connect(l1user1).provideFW(
    zeroAddr,
    l2Wallet.address,
    twoETH,
    saleCount,
    200000,
    {
      value: twoETH,
    }
  )
  await providerTx.wait()
  console.log('providerTx : ', providerTx.hash)

  await messenger.waitForMessageStatus(providerTx.hash, MessageStatus.RELAYED)

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance:', l2Balance.toString())
  l2ETHBalance = await l2ETHERC20.balanceOf(l2Wallet.address)
  console.log('l2 ETH(ERC20) balance (l1Wallet): ', l2ETHBalance.toString())

  l1BalanceUser1 = await l1user1.getBalance()
  console.log('l1 native balance (ETH) (User1): ', l1BalanceUser1.toString())
  l2BalanceUser1 = await l2user1.getBalance()
  console.log('l2 native balance (TON) (User1): ', l2BalanceUser1.toString())
  l2ETHBalanceUser1 = await l2ETHERC20.balanceOf(l2user1.address)
  console.log('l2 ETH(ERC20) balance (User1): ', l2ETHBalanceUser1.toString())

  L2FastWithdrawBalance = await l2Provider.getBalance(L2FastWithDrawContract.address)
  console.log('provider after l2 native balance (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  saleInformation = await L2FastWithDrawContract.dealData(1)
  console.log('saleInformation : ', saleInformation)
}

task('deposit-ETH', 'deposit ETH to L1 -> L2.')
  .addParam('amount', 'Deposit amount', fiveETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositETH(args.amount)
  })

task('fastwithdraw-ETH', 'fastWithdraw L2 ETH ERC20 test L1 <-> L2.')
  .addParam('amount', 'Withdrawal amount', oneETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await fastwithdrawERC20Token(args.amount)
  })
