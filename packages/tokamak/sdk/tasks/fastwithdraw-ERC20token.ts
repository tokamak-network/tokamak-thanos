import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys } from '@eth-optimism/core-utils'
import OptimismMintableERC20TokenFactoryABI from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json'
import OptimismMintableERC20TokenABI from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'
import { BytesLike, ethers, Event } from 'ethers'

import { CrossChainMessenger, MessageStatus, NativeTokenBridgeAdapter, NumberLike } from '../src'
import MockERC20ABI from '../../contracts-bedrock/forge-artifacts/MockERC20Token.sol/MockERC20Token.json'
import L1FastWithdrawABI from '../../contracts-bedrock/forge-artifacts/L1FastWithdraw.sol/L1FastWithdraw.json'
import L2FastWithdrawABI from '../../contracts-bedrock/forge-artifacts/L2FastWithdraw.sol/L2FastWithdraw.json'
import L1StandardBridgeABI from '../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
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
  await MockERC20.mint(l1user1.address, amount)

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

  let l1MockBalanceUser1 = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance: ', l1MockBalanceUser1.toString())

  let l2MockBalanceUser1 = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance: ', l2MockBalanceUser1.toString())

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  // let optimismPortalStorage = await OptimismPortalContract.depositedAmount()
  // console.log('optimismPortalStorage: ', optimismPortalStorage.toString())

  // let standardStorage = await L1StandardBridgeContract.deposits(l2NativeToken,predeploys.LegacyERC20ETH)
  // console.log('standardStorage: ', standardStorage.toString())

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

  l1MockBalanceUser1 = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance: ', l1MockBalanceUser1.toString())

  l2MockBalanceUser1 = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance: ', l2MockBalanceUser1.toString())

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

}

const fastwithdrawERC20Token = async (amount: NumberLike) => {
  console.log('Withdraw Native token:', amount)

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

  l1MockAddress = '0x9f4282cea29432724BbefF6ab4394B338e0fabB6'
  l2MockAddress = '0x85248c7403da665a13436C4B5120d0EB2dc80F96'

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
  let l1MockBalance = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance(L1Wallet): ', l1MockBalance.toString())
  let l2MockBalance = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance(L2Wallet): ', l2MockBalance.toString())

  let l1MockBalanceUser1 = await MockERC20.balanceOf(l1user1.address)
  console.log('l1MockBalance(User1): ', l1MockBalanceUser1.toString())
  let l2MockBalanceUser1 = await l2MockERC20.balanceOf(l2user1.address)
  console.log('l2MockBalance(User1): ', l2MockBalanceUser1.toString())

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance (ETH) (l1Wallet): ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance (TON) (l1Wallet): ', l2Balance.toString())

  let l1BalanceUser1 = await l1user1.getBalance()
  console.log('l1 native balance (ETH) (User1): ', l1BalanceUser1.toString())
  let l2BalanceUser1 = await l2user1.getBalance()
  console.log('l2 native balance (TON) (User1): ', l2BalanceUser1.toString())

  //request L2
  let L2FastWithdrawBalance = await l2MockERC20.balanceOf(L2FastWithDrawContract.address)
  console.log('before L2 ERC20 (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  // l1MockBalanceUser1 = await MockERC20.balanceOf(l1user1.address)
  // console.log('l1MockBalance(User1): ', l1MockBalanceUser1.toString())
  // l2MockBalanceUser1 = await l2MockERC20.balanceOf(l2user1.address)
  // console.log('l2MockBalance(User1): ', l2MockBalanceUser1.toString())
  // l1MockBalance = await MockERC20.balanceOf(l1Wallet.address)
  // console.log('l1MockBalance(L1Wallet): ', l1MockBalance.toString())
  // l2MockBalance = await l2MockERC20.balanceOf(l2Wallet.address)
  // console.log('l2MockBalance(L2Wallet): ', l2MockBalance.toString())

  // let tx = await l2MockERC20.connect(l2Wallet).approve(l2user1.address, threeETH)
  // await tx.wait()
  // console.log('pass the approve(transferFromTest)')

  // tx = await l2MockERC20.connect(l2user1).transferFrom(l2Wallet.address, l2user1.address,threeETH)
  // await tx.wait()
  // console.log('pass the transferFrom(transferFromTest)')

  const tx = await l2MockERC20.connect(l2Wallet).approve(L2FastWithDrawContract.address, threeETH)
  await tx.wait()
  console.log('pass the approve')

  await (await L2FastWithDrawContract.connect(l2Wallet).requestFW(
    l2MockERC20.address,
    threeETH,
    twoETH
  )).wait()
  console.log('pass the request')

  const saleCount = await L2FastWithDrawContract.salecount()
  console.log('saleCount : ', saleCount)
  let saleInformation = await L2FastWithDrawContract.dealData(saleCount)
  console.log('saleInformation : ', saleInformation)

  L2FastWithdrawBalance = await l2MockERC20.balanceOf(L2FastWithDrawContract.address)
  console.log('after L2 ERC20 (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  // let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
  //   l1user1.address
  // )

  // console.log(
  //   'native token(TON) balance in L1 (user1): ',
  //   l2NativeTokenBalance.toString()
  // )

  // let l2NativeTokenBalanceWallet = await l2NativeTokenContract.balanceOf(
  //   l1Wallet.address
  // )

  // console.log(
  //   'native token(TON) balance in L1 (Wallet): ',
  //   l2NativeTokenBalanceWallet.toString()
  // )

  const providerApproveTx = await MockERC20.connect(l1user1).approve(L1FastWithDrawContract.address, twoETH)
  await providerApproveTx.wait()
  console.log('pass the L1 TON approve')


  const providerTx = await L1FastWithDrawContract.connect(l1user1).provideFW(
    MockERC20.address,
    l2Wallet.address,
    twoETH,
    saleCount,
    200000
  )
  await providerTx.wait()
  console.log('providerTx : ', providerTx.hash)

  await messenger.waitForMessageStatus(providerTx.hash, MessageStatus.RELAYED)

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance:', l2Balance.toString())

  l1BalanceUser1 = await l1user1.getBalance()
  console.log('l1 native balance (ETH) (User1): ', l1BalanceUser1.toString())
  l2BalanceUser1 = await l2user1.getBalance()
  console.log('l2 native balance (TON) (User1): ', l2BalanceUser1.toString())

  l1MockBalance = await MockERC20.balanceOf(l1Wallet.address)
  console.log('l1MockBalance(L1Wallet): ', l1MockBalance.toString())
  l2MockBalance = await l2MockERC20.balanceOf(l2Wallet.address)
  console.log('l2MockBalance(L2Wallet): ', l2MockBalance.toString())

  l1MockBalanceUser1 = await MockERC20.balanceOf(l1user1.address)
  console.log('l1MockBalance(User1): ', l1MockBalanceUser1.toString())
  l2MockBalanceUser1 = await l2MockERC20.balanceOf(l2user1.address)
  console.log('l2MockBalance(User1): ', l2MockBalanceUser1.toString())

  // l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
  //   l1user1.address
  // )
  // console.log(
  //   'native token(TON) balance in L1 (user1): ',
  //   l2NativeTokenBalance.toString()
  // )

  // l2NativeTokenBalanceWallet = await l2NativeTokenContract.balanceOf(
  //   l1Wallet.address
  // )
  // console.log(
  //   'native token(TON) balance in L1 (Wallet): ',
  //   l2NativeTokenBalanceWallet.toString()
  // )

  L2FastWithdrawBalance = await l2Provider.getBalance(L2FastWithDrawContract.address)
  console.log('provider after l2 native balance (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  saleInformation = await L2FastWithDrawContract.dealData(1)
  console.log('saleInformation : ', saleInformation)
}

task('deposit-ERC20token', 'deposit L2ERC20 to L1 -> L2.')
  .addParam('amount', 'Deposit amount', fiveETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositERC20Token(args.amount)
  })

task('fastwithdraw-ERC20token', 'fastWithdraw L2ERC20 test L1 <-> L2.')
  .addParam('amount', 'Withdrawal amount', oneETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await fastwithdrawERC20Token(args.amount)
  })
