import { task } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { Event, Contract, Wallet, utils, ethers, BytesLike } from 'ethers'
import { predeploys } from '@tokamak-network/core-utils'
import Artifact__OptimismMintableERC20TokenFactory from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json'
import Artifact__OptimismMintableERC20Token from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'
import Artifact__WNativeToken from '@tokamak-network/thanos-contracts/forge-artifacts/WNativeToken.sol/WNativeToken.json'

import { CrossChainMessenger, MessageStatus, Portals } from '../src'

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

let nativeTokenAddress = process.env.NATIVE_TOKEN || ''
let addressManager = process.env.ADDRESS_MANAGER || ''
let l1CrossDomainMessenger = process.env.L1_CROSS_DOMAIN_MESSENGER || ''
let l1StandardBridge = process.env.L1_STANDARD_BRIDGE || ''
let optimismPortal = process.env.OPTIMISM_PORTAL || ''
let l2OutputOracle = process.env.L2_OUTPUT_ORACLE || ''

const updateAddresses = async (hre: HardhatRuntimeEnvironment) => {
  if (nativeTokenAddress === '') {
    const Deployment__NativeToken = await hre.deployments.get('L2NativeToken')
    nativeTokenAddress = Deployment__NativeToken.address
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

const deployWTON = async (
  hre: HardhatRuntimeEnvironment,
  signer: Wallet,
  wrap: boolean
): Promise<Contract> => {
  const Factory__WTON = new hre.ethers.ContractFactory(
    Artifact__WNativeToken.abi,
    Artifact__WNativeToken.bytecode.object,
    signer
  )

  console.log('Sending deployment transaction')
  const WTON = await Factory__WTON.deploy()
  const receipt = await WTON.deployTransaction.wait()
  console.log(`WTON deployed: ${receipt.transactionHash}`)

  if (wrap) {
    const deposit = await signer.sendTransaction({
      value: utils.parseEther('1'),
      to: WTON.address,
    })
    await deposit.wait()
  }

  return WTON
}

const createOptimismMintableERC20 = async (
  L1ERC20: Contract,
  l2Signer: Wallet
): Promise<Contract> => {
  const OptimismMintableERC20TokenFactory = new Contract(
    predeploys.OptimismMintableERC20Factory,
    Artifact__OptimismMintableERC20TokenFactory.abi,
    l2Signer
  )

  const name = await L1ERC20.name()
  const symbol = await L1ERC20.symbol()

  const tx =
    await OptimismMintableERC20TokenFactory.createOptimismMintableERC20(
      L1ERC20.address,
      `L2 ${name}`,
      `L2-${symbol}`
    )

  const receipt = await tx.wait()
  const event = receipt.events.find(
    (e: Event) => e.event === 'OptimismMintableERC20Created'
  )

  if (!event) {
    throw new Error('Unable to find OptimismMintableERC20Created event')
  }

  const l2WTONAddress = event.args.localToken
  console.log(`Deployed to ${l2WTONAddress}`)

  return new Contract(
    l2WTONAddress,
    Artifact__OptimismMintableERC20Token.abi,
    l2Signer
  )
}

const depositWTON = async (hre: HardhatRuntimeEnvironment) => {
  const l1ChainId = await l1Wallet.getChainId()
  const l2ChainId = await l2Wallet.getChainId()

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
  const messenger = new CrossChainMessenger({
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
    l1ChainId,
    l2ChainId,
    nativeTokenAddress,
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
  })
  // Ensure deployments module is initialized
  await hre.deployments.all()

  console.log('Deploying WTON to L1')
  const WTON = await deployWTON(hre, l1Wallet, true)
  console.log(`Deployed to ${WTON.address}`)

  // Save WTON address to deployments
  await hre.deployments.save('WTON', {
    address: WTON.address,
    abi: [WTON.interface.format('json')],
  })

  console.log('Creating L2 WTON')
  const OptimismMintableERC20 = await createOptimismMintableERC20(
    WTON,
    l2Wallet
  )

  // Save OptimismMintableERC20 address to deployments
  await hre.deployments.save('OptimismMintableERC20', {
    address: OptimismMintableERC20.address,
    abi: [OptimismMintableERC20.interface.format('json')],
  })

  console.log(`Approving WTON for deposit`)
  const approvalTx = await messenger.approveERC20(
    WTON.address,
    OptimismMintableERC20.address,
    hre.ethers.constants.MaxUint256
  )
  await approvalTx.wait()
  console.log('WTON approved')

  // report balances
  console.log(`Balance WTON before depositing...`)
  let l1Balance = await WTON.balanceOf(l1Wallet.address)
  console.log('l1 WTON balance: ', l1Balance.toString())

  let l2Balance = await OptimismMintableERC20.balanceOf(l2Wallet.address)
  console.log('l2 WTON balance:', l2Balance.toString())

  // deposit WTON
  console.log('Depositing WTON to L2')
  const depositTx = await messenger.bridgeERC20(
    WTON.address,
    OptimismMintableERC20.address,
    utils.parseEther('1')
  )
  await depositTx.wait()
  console.log(`ERC20 deposited - ${depositTx.hash}`)

  const portals = new Portals({
    contracts: {
      l1: l1Contracts,
    },
    l1ChainId,
    l2ChainId,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  const relayedDepositTx =
    await portals.waitingDepositTransactionRelayedUsingL1Tx(depositTx.hash)
  console.log('relayed tx:', relayedDepositTx)

  console.log(`Balance WTON after depositing...`)

  l1Balance = await WTON.balanceOf(l1Wallet.address)
  console.log('l1 WTON balance: ', l1Balance.toString())

  l2Balance = await OptimismMintableERC20.balanceOf(l2Wallet.address)
  console.log('l2 WTON balance:', l2Balance.toString())
  console.log(`Deposit success`)
}

const withdrawWTON = async (hre: HardhatRuntimeEnvironment) => {
  const l1ChainId = (await l1Provider.getNetwork()).chainId
  const l2ChainId = (await l2Provider.getNetwork()).chainId

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

  const messenger = new CrossChainMessenger({
    bedrock: true,
    contracts: {
      l1: l1Contracts,
    },
    l1ChainId,
    l2ChainId,
    nativeTokenAddress,
    l1SignerOrProvider: l1Wallet,
    l2SignerOrProvider: l2Wallet,
  })

  // Access the stored contract addresses
  const l1WTONAddress = (await hre.deployments.get('WTON')).address
  const l2WTONAddress = (await hre.deployments.get('OptimismMintableERC20'))
    .address

  if (!l1WTONAddress || !l2WTONAddress) {
    throw new Error(
      'WTON contract addresses are not set. Make sure to run deposit-erc20 first.'
    )
  }

  const WTON = new hre.ethers.Contract(
    l1WTONAddress,
    Artifact__WNativeToken.abi,
    l1Wallet
  )
  const OptimismMintableERC20 = new hre.ethers.Contract(
    l2WTONAddress,
    Artifact__OptimismMintableERC20Token.abi,
    l2Wallet
  )

  console.log('Starting withdrawal')
  console.log(`Balance WTON before withdrawing...`)
  let l1Balance = await WTON.balanceOf(l1Wallet.address)
  console.log('l1 WTON balance: ', l1Balance.toString())

  let l2Balance = await OptimismMintableERC20.balanceOf(l2Wallet.address)
  console.log('l2 WTON balance:', l2Balance.toString())

  // report balances

  const withdraw = await messenger.withdrawERC20(
    WTON.address,
    OptimismMintableERC20.address,
    utils.parseEther('1')
  )
  const withdrawalReceipt = await withdraw.wait()

  console.log(
    ' Withdrawal Tx:',
    withdrawalReceipt.transactionHash,
    ' Block',
    withdrawalReceipt.blockNumber,
    ' hash',
    withdraw.hash
  )

  await messenger.waitForMessageStatus(
    withdrawalReceipt,
    MessageStatus.READY_TO_PROVE
  )
  console.log('Prove the message')
  const proveTx = await messenger.proveMessage(withdrawalReceipt)
  const proveReceipt = await proveTx.wait(3)
  console.log('Proved the message:', proveReceipt.transactionHash)

  const finalizeInterval = setInterval(async () => {
    const currentStatus = await messenger.getMessageStatus(withdrawalReceipt)
    console.log('Message status:', currentStatus)
  }, 3000)

  try {
    await messenger.waitForMessageStatus(
      withdrawalReceipt,
      MessageStatus.READY_FOR_RELAY
    )
  } finally {
    clearInterval(finalizeInterval)
  }

  console.log(`Balance WTON before finalizing...`)
  l1Balance = await WTON.balanceOf(l1Wallet.address)
  console.log('l1 WTON balance: ', l1Balance.toString())

  l2Balance = await OptimismMintableERC20.balanceOf(l2Wallet.address)
  console.log('l2 WTON balance:', l2Balance.toString())

  const tx = await messenger.finalizeMessage(withdrawalReceipt)
  const receipt = await tx.wait()
  console.log('Finalized message tx', receipt.transactionHash)

  console.log(`Balance WTON after withdrawing...`)
  l1Balance = await WTON.balanceOf(l1Wallet.address)
  console.log('l1 WTON balance: ', l1Balance.toString())

  l2Balance = await OptimismMintableERC20.balanceOf(l2Wallet.address)
  console.log('l2 WTON balance:', l2Balance.toString())
  console.log('Withdrawal success')
}

// deploys a WTON contract, mints some WTON and then
// deposits that into L2 through the StandardBridge.
task('deposit-erc20', 'Deposit WTON onto L2.').setAction(async (args, hre) => {
  await updateAddresses(hre)
  await depositWTON(hre)
})

task('withdraw-erc20', 'Withdraw WTON from L2 to L1').setAction(
  async (args, hre) => {
    await updateAddresses(hre)
    await withdrawWTON(hre)
  }
)
