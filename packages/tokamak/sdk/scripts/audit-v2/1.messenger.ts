import hardhat from 'hardhat'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BigNumber, BytesLike, Wallet, ethers } from 'ethers'
import * as l1StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
import * as l2StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'
import * as l1CrossDomainMessengerAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L1CrossDomainMessenger.sol/L1CrossDomainMessenger.json'
import * as OptimismPortalAbi from '@tokamak-network/titan2-contracts/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'
import * as l2CrossDomainMessengerAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2CrossDomainMessenger.sol/L2CrossDomainMessenger.json'
// import * as l2OutputOracleAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2OutputOracle.sol/L2OutputOracle.json'

// import * as l2ToL1MessagePasserAbi from '../../../contracts-bedrock/forge-artifacts/L2ToL1MessagePasser.sol/L2ToL1MessagePasser.json'
import { CrossChainMessenger, MessageStatus } from '../../src'
import Artifact__MockHello from '../../../contracts-bedrock/forge-artifacts/MockHello.sol/MockHello.json'
import {
  erc20ABI,
  deployHello,
  getBalances,
  differenceLog,
  // deployERC20,
  // createOptimismMintableERC20,
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

const l2CrossDomainMessenger =
  process.env.L2_CROSS_DOMAIN_MESSENGER ||
  '0x4200000000000000000000000000000000000007'
let l2NativeToken = process.env.L2NativeToken || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''
// const l2ToL1MessagePasser =
//   process.env.L2ToL1MessagePasser ||
//   '0x4200000000000000000000000000000000000016'

const l2StandardBridge = process.env.L2_STANDARD_BRIDGE || ''
const legacy_ERC20_ETH = process.env.LEGACY_ERC20_ETH || ''
const l2_ERC20_ETH = process.env.ETH || ''
const l2EthContract = new ethers.Contract(l2_ERC20_ETH, erc20ABI, l2Wallet)

let l1BridgeContract
let l1CrossDomainMessengerContract
let OptomismPortalContract
let l2CrossDomainMessengerContract
// let l2ToL1MessagePasserContract
// let l2OutputOracleContract

let l1Contracts
let messenger
let tonContract
let helloContractL1
let helloContractL2
// let l1ERC20Token
// let l2ERC20Token

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

const messenger_1_depositTON_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== messenger_1_depositTON_L1_TO_L2  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const allowanceAmount = await tonContract.allowance(
    l1Wallet.address,
    l1CrossDomainMessenger
  )
  if (allowanceAmount < amount) {
    await (
      await tonContract
        .connect(l1Wallet)
        .approve(l1CrossDomainMessenger, amount)
    ).wait()
  }

  const sendTx = await (
    await l1CrossDomainMessengerContract
      .connect(l1Wallet)
      .sendNativeTokenMessage(l1Wallet.address, amount, '0x', 20000)
  ).wait()
  console.log('\nsendTx:', sendTx.transactionHash)

  // const topic = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessage');
  // const topic1 = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessageExtension1');
  // const topic2 = OptomismPortalContract.interface.getEventTopic('TransactionDeposited');

  // await logEvent(sendTx, topic, l1CrossDomainMessengerContract , 'SentMessage ' );
  // await logEvent(sendTx, topic1, l1CrossDomainMessengerContract , 'SentMessageExtension1 ' );
  // await logEvent(sendTx, topic2, OptomismPortalContract , 'TransactionDeposited ' );

  await messenger.waitForMessageStatus(
    sendTx.transactionHash,
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

const messenger_2_depositTON_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== messenger_2_depositTON_L1_TO_L2  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const data = ethers.utils.solidityPack(
    ['address', 'uint32', 'bytes'],
    [l1Wallet.address, 20000, '0x']
  )

  const sendTx = await (
    await tonContract
      .connect(l1Wallet)
      .approveAndCall(l1CrossDomainMessenger, amount, data)
  ).wait()
  console.log('\napproveAndCallTx:', sendTx.transactionHash)

  // const topic = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessage');
  // const topic1 = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessageExtension1');
  // const topic2 = OptomismPortalContract.interface.getEventTopic('TransactionDeposited');

  // await logEvent(sendTx, topic, l1CrossDomainMessengerContract , 'SentMessage ' );
  // await logEvent(sendTx, topic1, l1CrossDomainMessengerContract , 'SentMessageExtension1 ' );
  // await logEvent(sendTx, topic2, OptomismPortalContract , 'TransactionDeposited ' );

  await messenger.waitForMessageStatus(
    sendTx.transactionHash,
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

const messenger_3_createContract_L1_TO_L2 = async () => {
  console.log('\n==== messenger_3_createContract_L1_TO_L2  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  // console.log('Artifact__MockHello.bytecode', Artifact__MockHello.bytecode.object)
  const _byteCode = Artifact__MockHello.bytecode.object
  let _gasLimit = _byteCode.length * 16 + 21000
  _gasLimit = 120000
  console.log('_gasLimit', _gasLimit)

  const sendTx = await (
    await OptomismPortalContract.connect(l1Wallet).depositTransaction(
      ethers.constants.AddressZero,
      ethers.constants.Zero,
      _gasLimit * 3,
      true,
      _byteCode
    )
  ).wait()

  console.log('\nsendTx:', sendTx.transactionHash)

  // const topic = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessage');
  // const topic1 = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessageExtension1');
  // const topic2 = OptomismPortalContract.interface.getEventTopic('TransactionDeposited');

  // await logEvent(sendTx, topic, l1CrossDomainMessengerContract , 'SentMessage ' );
  // await logEvent(sendTx, topic1, l1CrossDomainMessengerContract , 'SentMessageExtension1 ' );
  // await logEvent(sendTx, topic2, OptomismPortalContract , 'TransactionDeposited ' );

  await messenger.waitForMessageStatus(
    sendTx.transactionHash,
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

const getMessageOfHello = async (helloContract) => {
  const blockNumber = await helloContract.blockNumber()
  const message = await helloContract.message()

  return {
    blockNumber,
    message,
  }
}

const messenger_4_sendMessage_L1_TO_L2 = async () => {
  console.log('\n==== messenger_4_sendMessage_L1_TO_L2  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const hello_prev = await getMessageOfHello(helloContractL1)
  const message = 'hi. from L1:' + hello_prev.blockNumber

  const callData = await helloContractL2.interface.encodeFunctionData('say', [
    message,
  ])
  const _gasLimit = callData.length * 16 + 21000
  // _gasLimit = 120000;
  console.log('_gasLimit', _gasLimit)

  const sendTx = await (
    await l1CrossDomainMessengerContract
      .connect(l1Wallet)
      .sendMessage(helloContractL2.address, callData, _gasLimit * 2)
  ).wait()

  console.log('\nsendTx:', sendTx.transactionHash)

  // const topic = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessage');
  // const topic1 = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessageExtension1');
  // const topic2 = OptomismPortalContract.interface.getEventTopic('TransactionDeposited');

  // await logEvent(sendTx, topic, l1CrossDomainMessengerContract , 'SentMessage ' );
  // await logEvent(sendTx, topic1, l1CrossDomainMessengerContract , 'SentMessageExtension1 ' );
  // await logEvent(sendTx, topic2, OptomismPortalContract , 'TransactionDeposited ' );

  await messenger.waitForMessageStatus(
    sendTx.transactionHash,
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

  const hello_after = await getMessageOfHello(helloContractL2)

  if (hello_after.message.localeCompare(message) === 0) {
    console.log('.. success sendMessage !! ')
  } else {
    console.log('.. fail sendMessage !! ')
  }
}

const messenger_5_sendMessage_L2_TO_L1 = async () => {
  console.log('\n==== messenger_5_sendMessage_L2_TO_L1  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const hello_prev = await getMessageOfHello(helloContractL2)

  const message = 'nice to meet you. from L2:' + hello_prev.blockNumber

  const callData = await helloContractL1.interface.encodeFunctionData('say', [
    message,
  ])
  const _gasLimit = callData.length * 16 + 21000
  console.log('_gasLimit', _gasLimit)

  const sendTx = await (
    await l2CrossDomainMessengerContract
      .connect(l2Wallet)
      .sendMessage(helloContractL1.address, callData, _gasLimit * 10)
  ).wait()

  console.log('\nsendTx:', sendTx.transactionHash)

  await messenger.waitForMessageStatus(
    sendTx.transactionHash,
    MessageStatus.READY_TO_PROVE
  )

  console.log('\nProve the message')
  const proveTx = await messenger.proveMessage(sendTx.transactionHash)
  const proveReceipt = await proveTx.wait(3)
  console.log('Proved the message: ', proveReceipt.transactionHash)

  const finalizeInterval = setInterval(async () => {
    const currentStatus = await messenger.getMessageStatus(
      sendTx.transactionHash
    )
    console.log('Message status: ', currentStatus)
  }, 3000)

  try {
    await messenger.waitForMessageStatus(
      sendTx.transactionHash,
      MessageStatus.READY_FOR_RELAY
    )
  } finally {
    clearInterval(finalizeInterval)
  }

  const tx = await messenger.finalizeMessage(sendTx.transactionHash)
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
  const hello_after = await getMessageOfHello(helloContractL1)

  if (hello_after.message.localeCompare(message) === 0) {
    console.log('.. success sendMessage !! ')
  } else {
    console.log('.. fail sendMessage !! ')
  }
}

const messenger_6_sendNativeTokenMessage_L1_TO_L2 = async (
  amount: BigNumber
) => {
  console.log('\n==== messenger_6_sendNativeTokenMessage_L1_TO_L2  ====== ')
  console.log('\n amount: ', ethers.utils.formatEther(amount))

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const hello_prev = await getMessageOfHello(helloContractL1)
  const message = 'hi. from L1:' + hello_prev.blockNumber

  const allowanceAmount = await tonContract.allowance(
    l1Wallet.address,
    l1CrossDomainMessenger
  )
  if (allowanceAmount < amount) {
    await (
      await tonContract
        .connect(l1Wallet)
        .approve(l1CrossDomainMessenger, amount)
    ).wait()
  }

  const callData = await helloContractL2.interface.encodeFunctionData('say', [
    message,
  ])
  const _gasLimit = callData.length * 16 + 21000
  // _gasLimit = 120000;
  console.log('_gasLimit', _gasLimit)

  const sendTx = await (
    await l1CrossDomainMessengerContract
      .connect(l1Wallet)
      .sendNativeTokenMessage(
        helloContractL2.address,
        amount,
        callData,
        _gasLimit * 10
      )
  ).wait()

  console.log('\nsendTx:', sendTx.transactionHash)

  try {
    await messenger.waitForMessageStatus(
      sendTx.transactionHash,
      MessageStatus.RELAYED
    )
  } catch (e) {
    console.log('\nerror', e)
    console.log('\n')
  }

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

  const hello_after = await getMessageOfHello(helloContractL2)

  console.log('hello_after.message', hello_after.message)
  console.log('message', message)

  if (hello_after.message.localeCompare(message) === 0) {
    console.log('.. success sendMessage !! ')
  } else {
    console.log('.. fail sendMessage !! ')
  }
}

const messenger_6_depositETH_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== messenger_6_depositETH_L1_TO_L2  ====== ')

  let err = true
  try {
    await (
      await l1CrossDomainMessengerContract
        .connect(l1Wallet)
        .sendMessage(l1Wallet.address, '0x', 20000, { value: amount })
    ).wait()
  } catch (e) {
    err = true
  }

  if (err) {
    console.log(
      ' Successfully occur error : execution reverted: Deny depositing ETH'
    )
  } else {
    console.log(
      '*** Error : The use of Ether in sendMessage of l1CrossDomainMessengerContract is prohibited. '
    )
  }
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

const messenger_7_withdrawETH_L2_TO_L1 = async (amount: BigNumber) => {
  console.log('\n==== messenger_7_withdrawETH_L2_TO_L1  ====== ')

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

  const localToken = l2_ERC20_ETH
  const remoteToken = ethers.constants.AddressZero
  const from = l2Wallet.address
  const to = l2Wallet.address

  const callData = await l2BridgeContract.interface.encodeFunctionData(
    'finalizeBridgeERC20',
    [localToken, remoteToken, from, to, amount, '0x']
  )

  const withdrawal = await l2CrossDomainMessengerContract
    .connect(l2Wallet)
    .sendMessage(l1StandardBridge, callData, 20000)

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

  l2CrossDomainMessengerContract = new ethers.Contract(
    l2CrossDomainMessenger,
    l2CrossDomainMessengerAbi.abi,
    l2Wallet
  )

  helloContractL1 = await deployHello(hardhat, l1Wallet)
  helloContractL2 = await deployHello(hardhat, l2Wallet)

  // const name = 'Test'
  // const symbol = 'TST'
  // const initialSupply = ethers.utils.parseEther('100000')

  // l1ERC20Token = await deployERC20(
  //   hardhat,
  //   l1Wallet,
  //   name,
  //   symbol,
  //   initialSupply
  // )
  // l2ERC20Token = await createOptimismMintableERC20(
  //   hardhat,
  //   l1ERC20Token,
  //   l2Wallet
  // )

  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

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

  // 1. deposit TON L1 to L2
  // 2. withdraw TON L2 to L1
  // 3. send create contract message L1 to L2
  // 4. send message L1 to L2
  // 5. send message L2 to L1

  // 6. sendNativeTokenMessage

  // 7. deposit ETH L1 to L2 -> L1에서 이더 입력 못함. 리버트
  // 8. withdraw ETH L2 to L1 -> 되어야 함.

  await messenger_1_depositTON_L1_TO_L2(depositAmount)

  await messenger_2_depositTON_L1_TO_L2(depositAmount)

  await bridge_2_withdrawTON_L2_TO_L1(withdrawAmount)
  await messenger_3_createContract_L1_TO_L2()

  await messenger_4_sendMessage_L1_TO_L2()
  await messenger_5_sendMessage_L2_TO_L1()

  await messenger_6_sendNativeTokenMessage_L1_TO_L2(depositAmount)
  await messenger_6_sendNativeTokenMessage_L1_TO_L2(ethers.constants.Zero)

  await messenger_6_depositETH_L1_TO_L2(depositAmount)
  await bridge_3_depositETH_L1_TO_L2(depositAmount)
  await messenger_7_withdrawETH_L2_TO_L1(withdrawAmount)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
