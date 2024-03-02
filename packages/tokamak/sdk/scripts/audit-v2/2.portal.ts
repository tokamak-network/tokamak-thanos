import hardhat from 'hardhat'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BigNumber, BytesLike, Wallet, ethers } from 'ethers'
import * as l1StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
// import * as l2StandardBridgeAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'
import * as OptimismPortalAbi from '@tokamak-network/titan2-contracts/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'
import * as l2CrossDomainMessengerAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2CrossDomainMessenger.sol/L2CrossDomainMessenger.json'
// import * as l2OutputOracleAbi from '@tokamak-network/titan2-contracts/forge-artifacts/L2OutputOracle.sol/L2OutputOracle.json'
// import * as l2ToL1MessagePasserAbi from '../../../contracts-bedrock/forge-artifacts/L2ToL1MessagePasser.sol/L2ToL1MessagePasser.json'
import { sleep } from '@eth-optimism/core-utils'

import { CrossChainMessenger, MessageStatus } from '../../src'
// import Artifact__MockHello from '../../../contracts-bedrock/forge-artifacts/MockHello.sol/MockHello.json'
import l1CrossDomainMessengerAbi from '../../../contracts-bedrock/forge-artifacts/L1CrossDomainMessenger.sol/L1CrossDomainMessenger.json'
import {
  erc20ABI,
  //   deployHello,
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

// const l2StandardBridge = process.env.L2_STANDARD_BRIDGE || ''
// const legacy_ERC20_ETH = process.env.LEGACY_ERC20_ETH || ''
const l2_ERC20_ETH = process.env.ETH || ''
const l2EthContract = new ethers.Contract(l2_ERC20_ETH, erc20ABI, l2Wallet)

let l1BridgeContract
let l1CrossDomainMessengerContract
let OptomismPortalContract
// let l2CrossDomainMessengerContract
// let l2ToL1MessagePasserContract
// let l2OutputOracleContract

let l1Contracts
let messenger
let tonContract
// let helloContractL1
// let helloContractL2
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

  //   const withdrawAmount = ethers.utils.parseEther('1')

  // 1. deposit from L1 to L2 with portal.depositTransaction
  // 2. send message from L1 to L2 with portal.depositTransaction
  // 3. send create contract to L2 with portal.depositTransaction
  // 4. send message from L1 to L2 with portal.OnApprove
  // 5. send create contract to L2 with portal.OnApprove

  await portal_1_depositTON_L1_TO_L2(depositAmount)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
