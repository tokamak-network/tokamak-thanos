import hardhat from 'hardhat'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BigNumber, BytesLike, Wallet, ethers } from 'ethers'
import * as l1StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
import * as l2StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'
import * as OptimismPortalAbi from '@tokamak-network/titan2-contracts/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'
// import * as l2CrossDomainMessengerAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2CrossDomainMessenger.sol/L2CrossDomainMessenger.json'
import {
  sleep,
  toRpcHexString,
  remove0x,
  toHexString,
} from '@eth-optimism/core-utils'

import * as l2ToL1MessagePasserAbi from '../../../contracts-bedrock/forge-artifacts/L2ToL1MessagePasser.sol/L2ToL1MessagePasser.json'
import { CrossChainMessenger, MessageStatus } from '../../src'
import Artifact__MockHello from '../../../contracts-bedrock/forge-artifacts/MockHello.sol/MockHello.json'
import l1CrossDomainMessengerAbi from '../../../contracts-bedrock/forge-artifacts/L1CrossDomainMessenger.sol/L1CrossDomainMessenger.json'
import { makeStateTrieProof, hashMessageHash } from '../../src/utils'
import {
  erc20ABI,
  deployHello,
  getBalances,
  differenceLog,
  getMessageOfHello,
  // deployERC20,
  // createOptimismMintableERC20,
} from '../shared'
const privateKey = process.env.PRIVATE_KEY as BytesLike

const l1Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L1_URL || 'http://localhost:8545'
)
const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL || 'http://localhost:9545'
)
const l1Wallet = new ethers.Wallet(privateKey, l1Provider)
const l2Wallet = new ethers.Wallet(privateKey, l2Provider)

const zeroAddr = '0x'.padEnd(42, '0')

// const l2CrossDomainMessenger =
//   process.env.L2_CROSS_DOMAIN_MESSENGER ||
//   '0x4200000000000000000000000000000000000007'
let l2NativeToken = process.env.L2NativeToken || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''
const l2ToL1MessagePasser =
  process.env.L2ToL1MessagePasser ||
  '0x4200000000000000000000000000000000000016'

const l2StandardBridge = process.env.L2_STANDARD_BRIDGE || ''
const legacy_ERC20_ETH = process.env.LEGACY_ERC20_ETH || ''
const l2_ERC20_ETH =
  process.env.ETH || '0x4200000000000000000000000000000000000486'
const l2EthContract = new ethers.Contract(l2_ERC20_ETH, erc20ABI, l2Wallet)

let l1BridgeContract
let l1CrossDomainMessengerContract
let OptomismPortalContract
// let l2CrossDomainMessengerContract
// let l2ToL1MessagePasserContract
// let l2OutputOracleContract

let l1Contracts
let messenger: CrossChainMessenger
let tonContract
let helloContractL1
let helloContractL2
// let l1ERC20Token
// let l2ERC20Token

const calculateSourceHash = (l1BlockHash, logIndex) => {
  const input = ethers.utils.concat([
    ethers.utils.arrayify(l1BlockHash),
    ethers.utils.hexZeroPad(ethers.BigNumber.from(logIndex).toHexString(), 32),
  ])
  const depositIDHash = ethers.utils.keccak256(input)

  const domainInput = ethers.utils.concat([
    ethers.utils.hexZeroPad(ethers.BigNumber.from(0).toHexString(), 32),
    ethers.utils.arrayify(depositIDHash),
  ])

  return ethers.utils.keccak256(domainInput)
}
const toMinimumEvenHex = (hexString) => {
  hexString = ethers.utils.hexStripZeros(hexString)
  return hexString.length % 2 === 0 ? hexString : `0x0${hexString.slice(2)}`
}

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
      ['sendNativeTokenMessage(address,uint256,bytes,uint32)'](
        l1Wallet.address,
        amount,
        '0x',
        20000
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

const portal_1_depositTON_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== portal_1_depositTON_L1_TO_L2  ====== ')

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
    OptomismPortalContract.address
  )
  if (allowanceAmount < amount) {
    await (
      await tonContract
        .connect(l1Wallet)
        .approve(OptomismPortalContract.address, amount)
    ).wait()
  }

  try {
    const sendTx = await (
      await OptomismPortalContract.connect(l1Wallet).depositTransaction(
        l1Wallet.address,
        amount,
        '21000',
        false,
        '0x'
      )
    ).wait()
    console.log('\nsendTx:', sendTx.transactionHash)

    await messenger.waitForMessageStatus(
      sendTx.transactionHash,
      MessageStatus.RELAYED
    )
  } catch (e) {
    console.log(e)

    console.log('\n sleep ... ')
    await sleep(60000)
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
}

const portal_2_sendMessage_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== portal_2_sendMessage_L1_TO_L2  ====== ')

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
    OptomismPortalContract.address
  )
  if (allowanceAmount < amount) {
    await (
      await tonContract
        .connect(l1Wallet)
        .approve(OptomismPortalContract.address, amount)
    ).wait()
  }
  const hello_prev = await getMessageOfHello(helloContractL1)
  const message = 'hi. from L1:' + hello_prev.blockNumber

  const tonBalanceHelloL2_prev = await l2Wallet.provider.getBalance(
    helloContractL2.address
  )

  const callData = helloContractL2.interface.encodeFunctionData('sayPayable', [
    message,
  ])

  const _gasLimit = (callData.length * 16 + 21000) * 3

  try {
    const sendTx = await (
      await OptomismPortalContract.connect(l1Wallet).depositTransaction(
        helloContractL2.address,
        amount,
        _gasLimit,
        false,
        callData
      )
    ).wait()
    console.log('\nsendTx:', sendTx.transactionHash)

    await messenger.waitForMessageStatus(
      sendTx.transactionHash,
      MessageStatus.RELAYED
    )
  } catch (e) {
    console.log(e)

    console.log('\n sleep ... ')
    await sleep(60000)
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

  const tonBalanceHelloL2_after = await l2Wallet.provider.getBalance(
    helloContractL2.address
  )

  console.log(
    'L2 ContracT TON Changed : ',
    ethers.utils.formatEther(
      tonBalanceHelloL2_after.sub(tonBalanceHelloL2_prev)
    )
  )
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

const passer_withdrawTON_L2_TO_L1 = async (amount: BigNumber) => {
  console.log('\n==== passer_withdrawTON_L2_TO_L1  ====== ')

  const beforeBalances = await getBalances(
    l1Wallet,
    l2Wallet,
    tonContract,
    l2EthContract,
    l1BridgeContract,
    l1CrossDomainMessengerContract,
    OptomismPortalContract
  )

  const l2ToL1MessagePasserContract = new ethers.Contract(
    l2ToL1MessagePasser,
    l2ToL1MessagePasserAbi.abi,
    l2Wallet
  )

  const withdrawal = await l2ToL1MessagePasserContract
    .connect(l2Wallet)
    .initiateWithdrawal(l1Wallet.address, 21000, '0x', {
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

  const messagePassedEvent = withdrawalTx.events[0]
  await sleep(60000)

  const l2OutputIndex =
    await messenger.contracts.l1.L2OutputOracle.getL2OutputIndexAfter(
      withdrawalTx.blockNumber
    )

  // Now pull the proposal out given the output index. Should always work as long as the above
  // codepath completed successfully.
  const proposal = await messenger.contracts.l1.L2OutputOracle.getL2Output(
    l2OutputIndex
  )

  const l2_BlockNumber = proposal.l2BlockNumber.toNumber()

  const stateTrieProof = await makeStateTrieProof(
    messenger.l2Provider as ethers.providers.JsonRpcProvider,
    l2_BlockNumber,
    messenger.contracts.l2.BedrockMessagePasser.address,
    hashMessageHash(messagePassedEvent.args.withdrawalHash)
  )

  const block = await (
    messenger.l2Provider as ethers.providers.JsonRpcProvider
  ).send('eth_getBlockByNumber', [toRpcHexString(l2_BlockNumber), false])

  const proof = {
    outputRootProof: {
      version: ethers.constants.HashZero,
      stateRoot: block.stateRoot,
      messagePasserStorageRoot: stateTrieProof.storageRoot,
      latestBlockhash: block.hash,
    },
    withdrawalProof: stateTrieProof.storageProof,
    l2OutputIndex,
  }

  const args = [
    [
      messagePassedEvent.args.nonce,
      messagePassedEvent.args.sender,
      messagePassedEvent.args.target,
      messagePassedEvent.args.value,
      messagePassedEvent.args.gasLimit,
      messagePassedEvent.args.data,
    ],
    proof.l2OutputIndex,
    [
      proof.outputRootProof.version,
      proof.outputRootProof.stateRoot,
      proof.outputRootProof.messagePasserStorageRoot,
      proof.outputRootProof.latestBlockhash,
    ],
    proof.withdrawalProof,
    {},
  ] as const

  const sendTx = await (
    await OptomismPortalContract.connect(l1Wallet).proveWithdrawalTransaction(
      ...args
    )
  ).wait()
  console.log('Proved the message: ', sendTx.transactionHash)

  await sleep(10000)

  const tx =
    await messenger.contracts.l1.OptimismPortal.finalizeWithdrawalTransaction([
      messagePassedEvent.args.nonce,
      messagePassedEvent.args.sender,
      messagePassedEvent.args.target,
      messagePassedEvent.args.value,
      messagePassedEvent.args.gasLimit,
      messagePassedEvent.args.data,
    ])
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

const portal_3_createContract_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== portal_3_createContract_L1_TO_L2  ====== ')

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
    OptomismPortalContract.address
  )
  if (allowanceAmount < amount) {
    await (
      await tonContract
        .connect(l1Wallet)
        .approve(OptomismPortalContract.address, amount)
    ).wait()
  }

  const callData = Artifact__MockHello.bytecode.object
  const _gasLimit = (callData.length * 16 + 21000) * 3

  const sendTx = await (
    await OptomismPortalContract.connect(l1Wallet).depositTransaction(
      ethers.constants.AddressZero,
      amount,
      _gasLimit,
      true,
      callData
    )
  ).wait()
  console.log('\nsendTx:', sendTx.transactionHash)

  const sourceHash = calculateSourceHash(
    sendTx.blockHash,
    sendTx.events[0].logIndex
  )

  const txData = [
    sourceHash, // sourceHash
    l1Wallet.address, // from
    '0x', // to, in case of contract creation, it is "0x"
    toMinimumEvenHex(amount.toHexString()), // mint
    toMinimumEvenHex(amount.toHexString()), // value
    toMinimumEvenHex(toHexString(_gasLimit)), // gasLimit
    '0x', // isSystemTx, always "0x"
    callData, // callData
  ]

  const rawTxL2 = '0x7e' + remove0x(ethers.utils.RLP.encode(txData))
  const txHashL2 = ethers.utils.keccak256(rawTxL2)

  while (true) {
    const l2Tx = await l2Wallet.provider.getTransactionReceipt(txHashL2)
    if (l2Tx) {
      console.log('\nl2Tx')
      console.log(l2Tx)
      break
    }

    sleep(1000)
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
}

const portal_4_onApproveDepositTon_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== portal_4_onApproveDepositTon_L1_TO_L2  ====== ')

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
    ['address', 'uint32', 'bool', 'bytes'],
    [l1Wallet.address, 30000, false, '0x']
  )

  // const unpackOnApproveData2 = await OptomismPortalContract["unpackOnApproveData2(bytes)"](data)
  // console.log('unpackOnApproveData2', unpackOnApproveData2)

  try {
    const sendTx = await (
      await tonContract
        .connect(l1Wallet)
        .approveAndCall(OptomismPortalContract.address, amount, data)
    ).wait()
    console.log('\nsendTx:', sendTx.transactionHash)

    await messenger.waitForMessageStatus(
      sendTx.transactionHash,
      MessageStatus.RELAYED
    )
  } catch (e) {
    console.log(e)

    console.log('\n sleep ... ')
    await sleep(30000)
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
}

const portal_5_approveAndCallWithMessage_L1_TO_L2 = async (
  amount: BigNumber
) => {
  console.log('\n==== portal_5_approveAndCallWithMessage_L1_TO_L2  ====== ')

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

  const tonBalanceHelloL2_prev = await l2Wallet.provider.getBalance(
    helloContractL2.address
  )

  const callData = await helloContractL2.interface.encodeFunctionData(
    'sayPayable',
    [message]
  )

  const _gasLimit = callData.length * 16 + 21000
  console.log('_gasLimit', _gasLimit)

  const data = ethers.utils.solidityPack(
    ['address', 'uint32', 'bool', 'bytes'],
    [helloContractL2.address, _gasLimit * 3, false, callData]
  )

  // const unpackOnApproveData2 = await OptomismPortalContract["unpackOnApproveData2(bytes)"](data)
  // console.log('unpackOnApproveData2', unpackOnApproveData2)

  try {
    const sendTx = await (
      await tonContract
        .connect(l1Wallet)
        .approveAndCall(OptomismPortalContract.address, amount, data)
    ).wait()
    console.log('\napproveAndCallTx:', sendTx.transactionHash)

    await messenger.waitForMessageStatus(
      sendTx.transactionHash,
      MessageStatus.RELAYED
    )
  } catch (e) {
    console.log(e)

    console.log('\n sleep ... ')
    await sleep(30000)
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

  const tonBalanceHelloL2_after = await l2Wallet.provider.getBalance(
    helloContractL2.address
  )

  console.log(
    'L2 ContracT TON Changed : ',
    ethers.utils.formatEther(
      tonBalanceHelloL2_after.sub(tonBalanceHelloL2_prev)
    )
  )
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

  // l2ToL1MessagePasserContract = new ethers.Contract(
  //   l2ToL1MessagePasser,
  //   l2ToL1MessagePasserAbi.abi,
  //   l2Wallet
  // )

  // l2CrossDomainMessengerContract = new ethers.Contract(
  //   l2CrossDomainMessenger,
  //   l2CrossDomainMessengerAbi.abi,
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

  // 1. deposit from L1 to L2 with portal.depositTransaction
  // 2. send message from L1 to L2 with portal.depositTransaction
  // 3. send create contract to L2 with portal.depositTransaction
  // 4. deposit from L1 to L2 with portal.OnApprove
  // 5. send message from L1 to L2 with portal.OnApprove
  // 6. send create contract to L2 with portal.OnApprove

  await messenger_1_depositTON_L1_TO_L2(depositAmount)
  await portal_1_depositTON_L1_TO_L2(
    depositAmount.mul(ethers.BigNumber.from('5'))
  )

  helloContractL1 = await deployHello(hardhat, l1Wallet)
  helloContractL2 = await deployHello(hardhat, l2Wallet)

  await portal_2_sendMessage_L1_TO_L2(depositAmount)
  await bridge_2_withdrawTON_L2_TO_L1(withdrawAmount)
  await passer_withdrawTON_L2_TO_L1(withdrawAmount)

  await portal_3_createContract_L1_TO_L2(ethers.constants.Zero) // on holding : I can't check the tx on L2 (create contract)

  await portal_4_onApproveDepositTon_L1_TO_L2(depositAmount)

  helloContractL1 = await deployHello(hardhat, l1Wallet)
  helloContractL2 = await deployHello(hardhat, l2Wallet)

  await portal_5_approveAndCallWithMessage_L1_TO_L2(depositAmount)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
