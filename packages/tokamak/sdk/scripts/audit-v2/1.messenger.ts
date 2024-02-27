import hardhat from 'hardhat'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BigNumber, BytesLike, Wallet, ethers } from 'ethers'
import * as l1StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
import * as l2StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'
import * as l1CrossDomainMessengerAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L1CrossDomainMessenger.sol/L1CrossDomainMessenger.json'
import * as OptimismPortalAbi from '@tokamak-network/titan2-contracts/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'

import { CrossChainMessenger, MessageStatus } from '../../src'
// import Artifact__ERC20 from '../../../contracts-bedrock/forge-artifacts/MockERC20.sol/MockERC20.json'
import {
  erc20ABI,
  getPortalDepositedAmount,
  getL1Balance,
  getL1ContractBalance,
  getL2Balance,
  differenceTonBalance,
  differenceEthBalance,
  differenceErc20Balance,
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
const l2EthContract = new ethers.Contract(l2_ERC20_ETH, erc20ABI, l2Wallet)

let l1BridgeContract
let l1CrossDomainMessengerContract
let OptomismPortalContract

let l1Contracts
let messenger
let tonContract

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
  const l1balance = await getL1Balance(l1Wallet, tonContract)
  const l2balance = await getL2Balance(l2Wallet, l2EthContract)

  const l1BridgeBalance = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal = await getPortalDepositedAmount(OptomismPortalContract)

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
  const l1balance_after = await getL1Balance(l1Wallet, tonContract)
  const l2balance_after = await getL2Balance(l2Wallet, l2EthContract)
  const l1BridgeBalance_after = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance_after = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance_after = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal_after = await getPortalDepositedAmount(OptomismPortalContract)

  await differenceTonBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet TON Changed : '
  )
  await differenceTonBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet TON Changed : '
  )

  await differenceEthBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet ETH Changed : '
  )
  await differenceEthBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet ETH Changed : '
  )

  await differenceTonBalance(
    l1BridgeBalance,
    l1BridgeBalance_after,
    'l1BridgeBalance TON Changed : '
  )
  await differenceTonBalance(
    l1MessengerBalance,
    l1MessengerBalance_after,
    'l1CrossDomainMessenger TON Changed : '
  )
  await differenceTonBalance(
    OptomismPortalBalance,
    OptomismPortalBalance_after,
    'OptomismPortalContract TON Changed : '
  )
  await differenceErc20Balance(
    portal,
    portal_after,
    'OptomismPortal depositAmount Changed : '
  )
}

const messenger_2_depositTON_L1_TO_L2 = async (amount: BigNumber) => {
  console.log('\n==== messenger_2_depositTON_L1_TO_L2  ====== ')
  const l1balance = await getL1Balance(l1Wallet, tonContract)
  const l2balance = await getL2Balance(l2Wallet, l2EthContract)

  const l1BridgeBalance = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal = await getPortalDepositedAmount(OptomismPortalContract)

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
  const l1balance_after = await getL1Balance(l1Wallet, tonContract)
  const l2balance_after = await getL2Balance(l2Wallet, l2EthContract)
  const l1BridgeBalance_after = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance_after = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance_after = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal_after = await getPortalDepositedAmount(OptomismPortalContract)

  await differenceTonBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet TON Changed : '
  )
  await differenceTonBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet TON Changed : '
  )

  await differenceEthBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet ETH Changed : '
  )
  await differenceEthBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet ETH Changed : '
  )

  await differenceTonBalance(
    l1BridgeBalance,
    l1BridgeBalance_after,
    'l1BridgeBalance TON Changed : '
  )
  await differenceTonBalance(
    l1MessengerBalance,
    l1MessengerBalance_after,
    'l1CrossDomainMessenger TON Changed : '
  )
  await differenceTonBalance(
    OptomismPortalBalance,
    OptomismPortalBalance_after,
    'OptomismPortalContract TON Changed : '
  )
  await differenceErc20Balance(
    portal,
    portal_after,
    'OptomismPortal depositAmount Changed : '
  )
}

const bridge_2_withdrawTON_L2_TO_L1 = async (amount: BigNumber) => {
  console.log('\n==== bridge_2_withdrawTON_L2_TO_L1  ====== ')
  const l1balance = await getL1Balance(l1Wallet, tonContract)
  const l2balance = await getL2Balance(l2Wallet, l2EthContract)

  const l1BridgeBalance = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal = await getPortalDepositedAmount(OptomismPortalContract)

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
  const l1balance_after = await getL1Balance(l1Wallet, tonContract)
  const l2balance_after = await getL2Balance(l2Wallet, l2EthContract)

  await differenceTonBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet TON Changed : '
  )
  await differenceTonBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet TON Changed : '
  )

  await differenceEthBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet ETH Changed : '
  )
  await differenceEthBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet ETH Changed : '
  )

  const l1BridgeBalance_after = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance_after = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance_after = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal_after = await getPortalDepositedAmount(OptomismPortalContract)

  await differenceTonBalance(
    l1BridgeBalance,
    l1BridgeBalance_after,
    'l1BridgeBalance TON Changed : '
  )
  await differenceTonBalance(
    l1MessengerBalance,
    l1MessengerBalance_after,
    'l1CrossDomainMessenger TON Changed : '
  )
  await differenceTonBalance(
    OptomismPortalBalance,
    OptomismPortalBalance_after,
    'OptomismPortalContract TON Changed : '
  )
  await differenceErc20Balance(
    portal,
    portal_after,
    'OptomismPortal depositAmount Changed : '
  )
}

const messenger_3_createContract_L1_TO_L2 = async () => {
  console.log('\n==== messenger_3_createContract_L1_TO_L2  ====== ')
  const l1balance = await getL1Balance(l1Wallet, tonContract)
  const l2balance = await getL2Balance(l2Wallet, l2EthContract)

  const l1BridgeBalance = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal = await getPortalDepositedAmount(OptomismPortalContract)

  // OptomismPortalContract.connect(l1Wallet).depositTransaction(
  //   ethers.constants.AddressZero,
  //   ethers.constants.Zero,
  //   20000,
  //   true,
  //   Artifact__ERC20.bytecode
  // )

  // const sendTx = await (
  //   await l1CrossDomainMessengerContract
  //     .connect(l1Wallet)
  //     .sendMessage(ethers.constants.AddressZero, _data, 20000)
  // ).wait()
  // console.log('\nsendTx:', sendTx.transactionHash)

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
  const l1balance_after = await getL1Balance(l1Wallet, tonContract)
  const l2balance_after = await getL2Balance(l2Wallet, l2EthContract)
  const l1BridgeBalance_after = await getL1ContractBalance(
    l1BridgeContract,
    tonContract
  )
  const l1MessengerBalance_after = await getL1ContractBalance(
    l1CrossDomainMessengerContract,
    tonContract
  )
  const OptomismPortalBalance_after = await getL1ContractBalance(
    OptomismPortalContract,
    tonContract
  )
  const portal_after = await getPortalDepositedAmount(OptomismPortalContract)

  await differenceTonBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet TON Changed : '
  )
  await differenceTonBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet TON Changed : '
  )

  await differenceEthBalance(
    l1balance,
    l1balance_after,
    'L1 Wallet ETH Changed : '
  )
  await differenceEthBalance(
    l2balance,
    l2balance_after,
    'L2 Wallet ETH Changed : '
  )

  await differenceTonBalance(
    l1BridgeBalance,
    l1BridgeBalance_after,
    'l1BridgeBalance TON Changed : '
  )
  await differenceTonBalance(
    l1MessengerBalance,
    l1MessengerBalance_after,
    'l1CrossDomainMessenger TON Changed : '
  )
  await differenceTonBalance(
    OptomismPortalBalance,
    OptomismPortalBalance_after,
    'OptomismPortalContract TON Changed : '
  )
  await differenceErc20Balance(
    portal,
    portal_after,
    'OptomismPortal depositAmount Changed : '
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

  // 1. deposit TON L1 to L2
  // 2. withdraw TON L2 to L1
  // 3. send create contract message L1 to L2
  // 4. send message L1 to L2
  // 5. send message L2 to L1
  // 6. deposit ETH L1 to L2
  // 7. withdraw ETH L2 to L1
  // 8. deposit ERC20 L2 to L1
  // 9. withdraw ERC20 L2 to L1

  await messenger_1_depositTON_L1_TO_L2(depositAmount)

  await messenger_2_depositTON_L1_TO_L2(depositAmount)

  await bridge_2_withdrawTON_L2_TO_L1(withdrawAmount)

  await messenger_3_createContract_L1_TO_L2()
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
