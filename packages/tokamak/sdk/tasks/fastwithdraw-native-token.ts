import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { predeploys } from '@eth-optimism/core-utils'
import { BytesLike, ethers } from 'ethers'

import { CrossChainMessenger, MessageStatus, NativeTokenBridgeAdapter, NumberLike } from '../src'
import L1FastWithdrawABI from '../../contracts-bedrock/forge-artifacts/L1FastWithdraw.sol/L1FastWithdraw.json'
import L2FastWithdrawABI from '../../contracts-bedrock/forge-artifacts/L2FastWithdraw.sol/L2FastWithdraw.json'
// import * as OptimismPortalABI from '../../contracts-bedrock/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'
// import * as L1StandardBridgeABI from '../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
// import L1FastWithdrawProxyABI from '../../contracts-bedrock/forge-artifacts/L1FastWithdrawProxy.sol/L1FastWithdrawProxy.json'
// import L2FastWithdrawProxyABI from '../../contracts-bedrock/forge-artifacts/L2FastWithdrawProxy.sol/L2FastWithdrawProxy.json'

// const OptimismPortalABI = require("../../contracts-bedrock/forge-artifacts/OptimismPortal.sol/OptimismPortal.json")
// const L1StandardBridgeABI = require("../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json")

console.log('Setup task...')

const privateKey = process.env.PRIVATE_KEY as BytesLike
const privateKey2 = process.env.PRIVATE_KEY2 as BytesLike

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
console.log('l1Wallet :', l1Wallet.address)
const l1user1 = new ethers.Wallet(privateKey2, l1Provider)
console.log('l1user1 :', l1user1.address)
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
        type: 'address'
      },
      {
        internalType: 'address',
        name: 'to',
        type: 'address'
      },
      {
        internalType: 'uint256',
        name: 'amount',
        type: 'uint256'
      }
    ],
    name: 'transferFrom',
    outputs: [
      {
        internalType: 'bool',
        name: '',
        type: 'bool'
      }
    ],
    stateMutability: 'nonpayable',
    type: 'function'
  },
]

const ETH = '0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000'

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

const depositNativeToken = async (amount: NumberLike) => {
  console.log('Deposit Native token:', amount)
  console.log('Native token address:', l2NativeToken)

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

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log('native token(TON) balance in L1:', Number(l2NativeTokenBalance.toString()))

  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())

  // let optimismPortalStorage = await OptimismPortalContract.depositedAmount()
  // console.log('optimismPortalStorage: ', optimismPortalStorage.toString())

  // let standardStorage = await L1StandardBridgeContract.deposits(l2NativeToken,predeploys.LegacyERC20ETH)
  // console.log('standardStorage: ', standardStorage.toString())

  if (Number(l2NativeTokenBalance.toString()) < Number(amount)) {
    console.log('start faucet')
    const tx = await l2NativeTokenContract.connect(l1Wallet).faucet(amount)
    await tx.wait()
    const l2NativeTokenBalance2 = await l2NativeTokenContract.balanceOf(
      l1Wallet.address
    )
    console.log('after faucet l2 native token(TON) balance in L1:', l2NativeTokenBalance2.toString())
  }

  const approveTx = await messenger.approveERC20(l2NativeToken, ETH, amount)
  await approveTx.wait()
  console.log('approveTx:', approveTx.hash)

  const depositTx = await messenger.depositERC20(l2NativeToken, ETH, amount)
  await depositTx.wait()
  console.log('depositTx:', depositTx.hash)

  await messenger.waitForMessageStatus(depositTx.hash, MessageStatus.RELAYED)

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance: ', l2Balance.toString())
  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(l1Wallet.address)
  console.log(
    'native token(TON) balance in L1: ',
    l2NativeTokenBalance.toString()
  )

  // optimismPortalStorage = await OptimismPortalContract.depositedAmount()
  // console.log('optimismPortalStorage: ', optimismPortalStorage.toString())

  // standardStorage = await L1StandardBridgeContract.deposits(l2NativeToken,predeploys.LegacyERC20ETH)
  // console.log('standardStorage: ', standardStorage.toString())
}

const fastwithdrawNativeToken = async (amount: NumberLike) => {
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

  // //deploy the L1, L2FastWithdraw
  // const L1FastWithDrawDep = new ethers.ContractFactory(
  //   L1FastWithdrawABI.abi,
  //   L1FastWithdrawABI.bytecode,
  //   l1Wallet
  // )

  // L1FastWithDraw = await L1FastWithDrawDep.deploy()
  // await L1FastWithDraw.deployed()

  // const L1FastWithdrawProxyDep = new ethers.ContractFactory(
  //   L1FastWithdrawProxyABI.abi,
  //   L1FastWithdrawProxyABI.bytecode,
  //   l1Wallet
  // )

  // L1FastWithDrawProxy = await L1FastWithdrawProxyDep.deploy()
  // await L1FastWithDrawProxy.deployed()
  // console.log("L1FastWithDrawProxy :", L1FastWithDrawProxy.address);

  // await (
  //   await L1FastWithDrawProxy.upgradeTo(L1FastWithDraw.address)
  // ).wait();
  // let imp2 = await L1FastWithDrawProxy.implementation()
  // console.log('check upgradeAddress : ', imp2)
  // console.log("L1FastWithDraw Address : ", L1FastWithDraw.address)
  // console.log('L1FastWithDrawProxy upgradeTo done')

  // const L2FastWithDrawDep = new ethers.ContractFactory(
  //   L2FastWithdrawABI.abi,
  //   L2FastWithdrawABI.bytecode,
  //   l2Wallet
  // )

  // L2FastWithDraw = await L2FastWithDrawDep.deploy()
  // await L2FastWithDraw.deployed()

  // const L2FastWithdrawProxyDep = new ethers.ContractFactory(
  //   L2FastWithdrawProxyABI.abi,
  //   L2FastWithdrawProxyABI.bytecode,
  //   l2Wallet
  // )

  // L2FastWithDrawProxy = await L2FastWithdrawProxyDep.deploy()
  // await L2FastWithDrawProxy.deployed()
  // console.log("L2FastWithDrawProxy :", L2FastWithDrawProxy.address);

  // await (
  //   await L2FastWithDrawProxy.upgradeTo(L2FastWithDraw.address)
  // ).wait();
  // imp2 = await L2FastWithDrawProxy.implementation()
  // console.log('check upgradeAddress : ', imp2)
  // console.log("L2FastWithDraw Address : ", L2FastWithDraw.address)
  // console.log('L2FastWithDrawProxy upgradeTo done')

  // //L1, L2 initialize
  // await (await L1FastWithDrawProxy.connect(l1Wallet).initialize(
  //   l1Contracts.L1CrossDomainMessenger,
  //   L2FastWithDrawProxy.address,
  //   zeroAddr,
  //   l2NativeTokenContract.address
  // )).wait();
  // console.log("L1FastWithdraw initialize done")

  // const checkL1Inform = await L1FastWithDrawProxy.LEGACY_l1token();
  // console.log("checkL1Inform :", checkL1Inform)
  // console.log("l2NativeTokenContract.address :", l2NativeTokenContract.address)

  // await (await L2FastWithDrawProxy.connect(l2Wallet).initialize(
  //   l2CrossDomainMessengerAddr,
  //   L1FastWithDrawProxy.address,
  //   predeploys.LegacyERC20ETH,
  //   l2NativeTokenContract.address
  // )).wait();
  // console.log("L2FastWithdraw initialize done")

  // const checkL2Inform = await L2FastWithDrawProxy.LEGACY_l1token();
  // console.log("checkL2Inform :", checkL2Inform)
  // console.log("l2NativeTokenContract.address :", l2NativeTokenContract.address)

  // //L1, L2 Contract set
  // L1FastWithDrawContract = new ethers.Contract(
  //   L1FastWithDrawProxy.address,
  //   L1FastWithdrawABI.abi,
  //   l1Wallet
  // )

  // L2FastWithDrawContract = new ethers.Contract(
  //   L2FastWithDrawProxy.address,
  //   L2FastWithdrawABI.abi,
  //   l2Wallet
  // )



  // deploy the L1, L2FastWithdraw
  const L1FastWithDrawDep = new ethers.ContractFactory(
    L1FastWithdrawABI.abi,
    L1FastWithdrawABI.bytecode,
    l1Wallet
  )

  L1FastWithDrawContract = await L1FastWithDrawDep.deploy()
  await L1FastWithDrawContract.deployed()
  console.log('L1FastWithDrawContract.address :', L1FastWithDrawContract.address)

  const L2FastWithDrawDep = new ethers.ContractFactory(
    L2FastWithdrawABI.abi,
    L2FastWithdrawABI.bytecode,
    l2Wallet
  )

  L2FastWithDrawContract = await L2FastWithDrawDep.deploy()
  await L2FastWithDrawContract.deployed()
  // console.log(await l2Provider.getCode(L2FastWithDrawContract.address))
  console.log('L2FastWithDrawContract.address :', L2FastWithDrawContract.address)

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
  console.log("L2FastWithdraw initialize done")

  const checkL2Inform = await L2FastWithDrawContract.crossDomainMessenger()
  console.log("checkL2Inform :", checkL2Inform)
  console.log("l2CrossDomainMessengerAddr :", l2CrossDomainMessengerAddr)

  // start the test
  let l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance (ETH) (l1Wallet): ', l1Balance.toString())
  let l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance (TON) (l1Wallet): ', l2Balance.toString())

  let l1BalanceUser1 = await l1user1.getBalance()
  console.log('l1 native balance (ETH) (User1): ', l1BalanceUser1.toString())
  let l2BalanceUser1 = await l2user1.getBalance()
  console.log('l2 native balance (TON) (User1): ', l2BalanceUser1.toString())

  //request L2
  let L2FastWithdrawBalance = await l2Provider.getBalance(L2FastWithDrawContract.address)
  console.log('before l2 native balance (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  await (await L2FastWithDrawContract.connect(l2Wallet).requestFW(
    predeploys.LegacyERC20ETH,
    threeETH,
    twoETH,
    {
      value: threeETH
    }
  )).wait()
  const saleCount = await L2FastWithDrawContract.salecount()
  console.log('saleCount : ', saleCount);
  let saleInformation = await L2FastWithDrawContract.dealData(saleCount)
  console.log('saleInformation : ', saleInformation);

  L2FastWithdrawBalance = await l2Provider.getBalance(L2FastWithDrawContract.address);
  console.log('after l2 native balance (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  let l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1user1.address
  )
  console.log(
    'native token(TON) balance in L1 (user1): ',
    l2NativeTokenBalance.toString()
  )

  let l2NativeTokenBalanceWallet = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log(
    'native token(TON) balance in L1 (Wallet): ',
    l2NativeTokenBalanceWallet.toString()
  )

  // //transferFrom test
  // let tx = await l2NativeTokenContract.connect(l1user1).approve(l1Wallet.address, twoETH)
  // await tx.wait();
  // console.log("private approve is done")
  // tx = await l2NativeTokenContract.connect(l1Wallet).transferFrom(l1user1.address,l1Wallet.address,twoETH)
  // await tx.wait();
  // console.log("private transferFrom is done")

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

  // if (Number(l2NativeTokenBalance.toString()) === 0) {
  //   console.log('start faucet')
  //   const tx = await l2NativeTokenContract.connect(l1user1).faucet(twoETH)
  //   await tx.wait()
  //   const l2NativeTokenBalance2 = await l2NativeTokenContract.balanceOf(
  //     l1user1.address
  //   )
  //   console.log('after faucet l2 native token(TON) balance in L1:', l2NativeTokenBalance2.toString())
  // }


  //provider L1
  if (Number(l2NativeTokenBalance.toString()) === 0) {
    console.log('start faucet')
    const tx = await l2NativeTokenContract.connect(l1user1).faucet(twoETH)
    await tx.wait()
    const l2NativeTokenBalance2 = await l2NativeTokenContract.balanceOf(
      l1user1.address
    )
    console.log('after faucet l2 native token(TON) balance in L1 (user1):', l2NativeTokenBalance2.toString())
  }

  const providerApproveTx = await l2NativeTokenContract.connect(l1user1).approve(L1FastWithDrawContract.address, twoETH)
  await providerApproveTx.wait()
  console.log("pass the L1 TON approve")


  const providerTx = await L1FastWithDrawContract.connect(l1user1).provideFW(
    l2NativeToken,
    l2Wallet.address,
    twoETH,
    saleCount,
    200000
  )
  await providerTx.wait()
  console.log("providerTx : ", providerTx.hash)

  await messenger.waitForMessageStatus(providerTx.hash, MessageStatus.RELAYED)

  l1Balance = await l1Wallet.getBalance()
  console.log('l1 native balance: ', l1Balance.toString())
  l2Balance = await l2Wallet.getBalance()
  console.log('l2 native balance:', l2Balance.toString())

  l1BalanceUser1 = await l1user1.getBalance()
  console.log('l1 native balance (ETH) (User1): ', l1BalanceUser1.toString())
  l2BalanceUser1 = await l2user1.getBalance()
  console.log('l2 native balance (TON) (User1): ', l2BalanceUser1.toString())

  l2NativeTokenBalance = await l2NativeTokenContract.balanceOf(
    l1user1.address
  )
  console.log(
    'native token(TON) balance in L1 (user1): ',
    l2NativeTokenBalance.toString()
  )

  l2NativeTokenBalanceWallet = await l2NativeTokenContract.balanceOf(
    l1Wallet.address
  )
  console.log(
    'native token(TON) balance in L1 (Wallet): ',
    l2NativeTokenBalanceWallet.toString()
  )

  L2FastWithdrawBalance = await l2Provider.getBalance(L2FastWithDrawContract.address)
  console.log('provider after l2 native balance (L2FastWithdrawBalance): ', L2FastWithdrawBalance.toString())

  saleInformation = await L2FastWithDrawContract.dealData(1)
  console.log("saleInformation : ", saleInformation)
}

task('deposit-nativetoken', 'request L2NativeToken to L2.')
  .addParam('amount', 'Deposit amount', fiveETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await depositNativeToken(args.amount)
  })

task('fastwithdraw-native-token', 'Withdraw native token from L2.')
  .addParam('amount', 'Withdrawal amount', oneETH.toString(), types.string)
  .setAction(async (args, hre) => {
    await updateAddresses(hre)
    await fastwithdrawNativeToken(args.amount)
  })
