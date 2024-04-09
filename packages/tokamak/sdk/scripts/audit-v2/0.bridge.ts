import hardhat from 'hardhat'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BigNumber, BytesLike, Wallet, ethers } from 'ethers'
import * as l1CrossDomainMessengerAbi from '@tokamak-network/thanos-contracts/forge-artifacts/L1CrossDomainMessenger.sol/L1CrossDomainMessenger.json'
import * as OptimismPortalAbi from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'

import { CrossChainMessenger, MessageStatus } from '../../src'
import * as l1StandardBridgeAbi from '../../../contracts-bedrock/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
import * as l2StandardBridgeAbi from '../../../contracts-bedrock/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'
import {
  erc20ABI,
  wtonABI,
  deployERC20,
  createOptimismMintableERC20,
  getErc20Balance,
  differenceErc20Balance,
  getBalances,
  differenceLog,
} from '../shared'

const privateKey = process.env.PRIVATE_KEY as BytesLike

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

const zeroAddr = '0x'.padEnd(42, '0')

let l2NativeToken = process.env.L2NativeToken || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''

const l2StandardBridge = process.env.L2_STANDARD_BRIDGE || ''
const legacy_ERC20_ETH = process.env.LEGACY_ERC20_ETH || ''
const l2_ERC20_ETH = process.env.ETH || ''
const l2_WTON = process.env.WTON || ''
const l2EthContract = new ethers.Contract(l2_ERC20_ETH, erc20ABI, l2Wallet)

let l1BridgeContract
let l1CrossDomainMessengerContract
let OptomismPortalContract

let l1Contracts
let messenger
let tonContract
let l1ERC20Token
let l2ERC20Token

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

const bridge_1_depositTON_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== bridge_1_depositTON_L1_TO_L2  ====== ')
  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const data = ethers.utils.solidityPack(['uint32', 'bytes'], [2000000, '0x'])

  const approveAndCallTx = await (
    await tonContract
      .connect(l1Wallet)
      .approveAndCall(l1Contracts.L1StandardBridge, amount, data)
  ).wait()
  console.log('\napproveAndCallTx:', approveAndCallTx.transactionHash)

  await messenger.waitForMessageStatus(
    approveAndCallTx.transactionHash,
    MessageStatus.RELAYED
  )

  const afterBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  await differenceLog(beforeBalances, afterBalances)
}

const bridge_2_withdrawTON_L2_TO_L1 = async (amount: BigNumber) => {
  console.log('\n==== bridge_2_withdrawTON_L2_TO_L1  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )
  const l2BridgeContract = new ethers.Contract(
    l2StandardBridge,
    l2StandardBridgeAbi.abi,
    l2Wallet
  )
  const withdrawal = await l2BridgeContract
    .connect(l2Wallet)
    .withdraw(legacy_ERC20_ETH, amount, 20000, '0x', {
      value: amount,
    })
  const withdrawalTx = await withdrawal.wait()
  console.log(
    '\nwithdrawal Tx:',
    withdrawalTx.transactionHash,
    ' Block',
    withdrawalTx.blockNumber,
    ' hash',
    withdrawal.hash
  )

  await messenger.waitForMessageStatus(
    withdrawalTx.transactionHash,
    MessageStatus.READY_TO_PROVE
  )

  console.log('\nProve the message')
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
  console.log('\nFinalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  const afterBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  await differenceLog(beforeBalances, afterBalances)
}

const bridge_3_depositETH_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== bridge_3_depositETH_L1_TO_L2  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )
  const deposition = await l1BridgeContract
    .connect(l1Wallet)
    .depositETH(20000, '0x', {
      value: amount,
    })
  const depositionTx = await deposition.wait()
  console.log(
    '\ndeposit Tx:',
    depositionTx.transactionHash,
    ' Block',
    depositionTx.blockNumber,
    ' hash',
    deposition.hash
  )

  await messenger.waitForMessageStatus(
    depositionTx.transactionHash,
    MessageStatus.RELAYED
  )

  const afterBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  await differenceLog(beforeBalances, afterBalances)
}

const bridge_4_withdrawETH_L2_TO_L1 = async (amount: BigNumber) => {
  console.log('\n==== bridge_4_withdrawETH_L2_TO_L1  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )
  const l2BridgeContract = new ethers.Contract(
    l2StandardBridge,
    l2StandardBridgeAbi.abi,
    l2Wallet
  )
  const withdrawal = await l2BridgeContract
    .connect(l2Wallet)
    .withdraw(l2_ERC20_ETH, amount, 20000, '0x', {
      value: amount,
    })
  const withdrawalTx = await withdrawal.wait()
  console.log(
    '\nwithdrawal Tx:',
    withdrawalTx.transactionHash,
    ' Block',
    withdrawalTx.blockNumber,
    ' hash',
    withdrawal.hash
  )

  await messenger.waitForMessageStatus(
    withdrawalTx.transactionHash,
    MessageStatus.READY_TO_PROVE
  )

  console.log('\nProve the message')
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
  console.log('\nFinalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  const afterBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  await differenceLog(beforeBalances, afterBalances)
}

const bridge_5_withdrawWTON_L2_TO_L1 = async (amount: BigNumber) => {
  console.log('\n==== bridge_5_withdrawWTON_L2_TO_L1 ====== ')

  const l2WtonContract = new ethers.Contract(l2_WTON, wtonABI, l2Wallet)
  const l2WtonBalanceCur = await l2WtonContract.balanceOf(l2Wallet.address)
  console.log(
    'l2 wton balance current: ',
    ethers.utils.formatEther(l2WtonBalanceCur)
  )

  const deposition = await l2WtonContract
    .connect(l2Wallet)
    .deposit({ value: amount })
  const depositionTx = await deposition.wait()
  console.log(
    '\ndeposit 1 ETH Tx:',
    depositionTx.transactionHash,
    ' Block',
    depositionTx.blockNumber,
    ' hash',
    deposition.hash
  )

  const l2WtonBalancePrev = await l2WtonContract.balanceOf(l2Wallet.address)
  console.log('l2 wton balance: ', ethers.utils.formatEther(l2WtonBalancePrev))

  let ret = false
  try {
    const l2BridgeContract = new ethers.Contract(
      l2StandardBridge,
      l2StandardBridgeAbi.abi,
      l2Wallet
    )
    await l2BridgeContract
      .connect(l2Wallet)
      .withdraw(l2_WTON, amount, 20000, '0x', {
        value: amount,
      })
  } catch (err) {
    ret = true
  }

  if (ret) {
    const withdrawal = await l2WtonContract.connect(l2Wallet).withdraw(amount)
    const withdrawTx = await withdrawal.wait()
    console.log(
      '\nwithdrawal 1 ETH Tx:',
      withdrawTx.transactionHash,
      ' Block',
      withdrawTx.blockNumber,
      ' hash',
      withdrawal.hash
    )
    const l2WtonBalanceAfter = await l2WtonContract.balanceOf(l2Wallet.address)
    console.log(
      'l2 wton balance: ',
      ethers.utils.formatEther(l2WtonBalanceAfter)
    )
  } else {
    console.log('wrong execution')
  }
}

const bridge_6_depositERC20_L1_TO_L2 = async (
  hre: HardhatRuntimeEnvironment,
  amount: BigNumber
) => {
  console.log('\n==== bridge_6_depositERC20_L1_TO_L2 ====== ')
  const name = 'Test'
  const symbol = 'TST'
  const initialSupply = ethers.utils.parseEther('100000')

  l1ERC20Token = await deployERC20(hre, l1Wallet, name, symbol, initialSupply)
  l2ERC20Token = await createOptimismMintableERC20(hre, l1ERC20Token, l2Wallet)

  await (
    await l1ERC20Token
      .connect(l1Wallet)
      .approve(l1BridgeContract.address, amount)
  ).wait()

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const l1token = await getErc20Balance(l1Wallet, l1ERC20Token)
  const l2token = await getErc20Balance(l2Wallet, l2ERC20Token)

  const deposition = await l1BridgeContract
    .connect(l1Wallet)
    .depositERC20(
      l1ERC20Token.address,
      l2ERC20Token.address,
      amount,
      20000,
      '0x'
    )
  const depositionTx = await deposition.wait()
  console.log(
    '\ndeposit Tx:',
    depositionTx.transactionHash,
    ' Block',
    depositionTx.blockNumber,
    ' hash',
    deposition.hash
  )

  await messenger.waitForMessageStatus(
    depositionTx.transactionHash,
    MessageStatus.RELAYED
  )

  const l1token_after = await getErc20Balance(l1Wallet, l1ERC20Token)
  const l2token_after = await getErc20Balance(l2Wallet, l2ERC20Token)

  const afterBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  await differenceLog(beforeBalances, afterBalances)

  await differenceErc20Balance(l1token, l1token_after, 'L1 ERC20 Changed : ')
  await differenceErc20Balance(l2token, l2token_after, 'L2 ERC20 Changed : ')
}

const bridge_7_withdrawERC20_L2_TO_L1 = async (amount: BigNumber) => {
  console.log('\n==== bridge_7_withdrawERC20_L2_TO_L1  ====== ')

  const l1token = await getErc20Balance(l1Wallet, l1ERC20Token)
  const l2token = await getErc20Balance(l2Wallet, l2ERC20Token)

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )
  const l2BridgeContract = new ethers.Contract(
    l2StandardBridge,
    l2StandardBridgeAbi.abi,
    l2Wallet
  )
  const withdrawal = await l2BridgeContract
    .connect(l2Wallet)
    .withdraw(l2ERC20Token.address, amount, 20000, '0x', {
      value: amount,
    })
  const withdrawalTx = await withdrawal.wait()
  console.log(
    '\nwithdrawal Tx:',
    withdrawalTx.transactionHash,
    ' Block',
    withdrawalTx.blockNumber,
    ' hash',
    withdrawal.hash
  )

  await messenger.waitForMessageStatus(
    withdrawalTx.transactionHash,
    MessageStatus.READY_TO_PROVE
  )

  console.log('\nProve the message')
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
  console.log('\nFinalized message tx', receipt.transactionHash)
  console.log('Finalized withdrawal')

  const afterBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )
  const l1token_after = await getErc20Balance(l1Wallet, l1ERC20Token)
  const l2token_after = await getErc20Balance(l2Wallet, l2ERC20Token)

  await differenceLog(beforeBalances, afterBalances)

  await differenceErc20Balance(l1token, l1token_after, 'L1 ERC20 Changed : ')
  await differenceErc20Balance(l2token, l2token_after, 'L2 ERC20 Changed : ')
}

const faucet = async (account: Wallet, amount: BigNumber) => {
  await (await tonContract.connect(account).faucet(amount)).wait()

  const l1TONTotalSupply = await tonContract.totalSupply()
  console.log(
    'l1 ton total supply:',
    ethers.utils.formatEther(l1TONTotalSupply)
  )
}

const setup = async () => {
  await updateAddresses(hardhat)

  l1Contracts = {
    StateCommitmentChain: zeroAddr,
    CanonicalTransactionChain: zeroAddr,
    BondManager: zeroAddr,
    AddressManager: addressManager,
    L1CrossDomainMessenger: l1CrossDomainMessenger,
    L1StandardBridge: l1StandardBridge,
    OptimismPortal: optimismPortal,
    L2OutputOracle: l2OutputOracle,
  }

  tonContract = new ethers.Contract(l2NativeToken, erc20ABI, l1Wallet)

  l1BridgeContract = new ethers.Contract(
    l1StandardBridge,
    l1StandardBridgeAbi.abi,
    l1Wallet
  )

  l1CrossDomainMessengerContract = new ethers.Contract(
    l1CrossDomainMessenger,
    l1CrossDomainMessengerAbi.abi,
    l1Wallet
  )

  OptomismPortalContract = new ethers.Contract(
    optimismPortal,
    OptimismPortalAbi.abi,
    l1Wallet
  )

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId
  //   console.log('l1ChainId',l1ChainId)

  messenger = new CrossChainMessenger({
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })
}

const main = async () => {
  await setup()
  await faucet(l1Wallet, ethers.utils.parseEther('100'))

  const depositAmount = ethers.utils.parseEther('2')
  const withdrawAmount = ethers.utils.parseEther('1')

  await bridge_1_depositTON_L1_TO_L2(depositAmount)
  await bridge_2_withdrawTON_L2_TO_L1(withdrawAmount)

  await bridge_3_depositETH_L1_TO_L2(depositAmount)
  await bridge_4_withdrawETH_L2_TO_L1(withdrawAmount)

  await bridge_5_withdrawWTON_L2_TO_L1(withdrawAmount)

  await bridge_6_depositERC20_L1_TO_L2(hardhat, depositAmount)
  await bridge_7_withdrawERC20_L2_TO_L1(withdrawAmount)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
