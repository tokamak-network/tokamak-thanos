import { promises as fs } from 'fs'

import { task, types } from 'hardhat/config'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import { SignerWithAddress } from '@nomiclabs/hardhat-ethers/signers'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { Event, Contract, Wallet, providers, utils, ethers } from 'ethers'
import { predeploys, sleep } from '@tokamak-network/core-utils'
import Artifact__OptimismMintableERC20TokenFactory from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json'
import Artifact__OptimismMintableERC20Token from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'
import Artifact__L2ToL1MessagePasser from '@tokamak-network/thanos-contracts/forge-artifacts/L2ToL1MessagePasser.sol/L2ToL1MessagePasser.json'
import Artifact__L2CrossDomainMessenger from '@tokamak-network/thanos-contracts/forge-artifacts/L2CrossDomainMessenger.sol/L2CrossDomainMessenger.json'
import Artifact__L2StandardBridge from '@tokamak-network/thanos-contracts/forge-artifacts/L2StandardBridge.sol/L2StandardBridge.json'
import Artifact__OptimismPortal from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismPortal.sol/OptimismPortal.json'
import Artifact__L1CrossDomainMessenger from '@tokamak-network/thanos-contracts/forge-artifacts/L1CrossDomainMessenger.sol/L1CrossDomainMessenger.json'
import Artifact__L1StandardBridge from '@tokamak-network/thanos-contracts/forge-artifacts/L1StandardBridge.sol/L1StandardBridge.json'
import Artifact__L2OutputOracle from '@tokamak-network/thanos-contracts/forge-artifacts/L2OutputOracle.sol/L2OutputOracle.json'
import Artifact__WNativeToken from '@tokamak-network/thanos-contracts/forge-artifacts/WNativeToken.sol/WNativeToken.json'

import {
  CrossChainMessenger,
  MessageStatus,
  CONTRACT_ADDRESSES,
  OEContractsLike,
  DEFAULT_L2_CONTRACT_ADDRESSES,
} from '../src'

const deployWTON = async (
  hre: HardhatRuntimeEnvironment,
  signer: SignerWithAddress,
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

const depositWTON = async (
  hre: HardhatRuntimeEnvironment,
  l2ProviderUrl: string,
  l1ContractsJsonPath: string,
  signerIndex: number
) => {
  const signers = await hre.ethers.getSigners()
  if (signers.length === 0) {
    throw new Error('No configured signers')
  }
  if (signerIndex < 0 || signers.length <= signerIndex) {
    throw new Error('Invalid signer index')
  }
  const signer = signers[signerIndex]
  const address = await signer.getAddress()
  console.log(`Using signer ${address}`)

  // Ensure that the signer has a balance before trying to
  // do anything
  const balance = await signer.getBalance()
  if (balance.eq(0)) {
    throw new Error('Signer has no balance')
  }

  const l2NativeToken = process.env.NATIVE_TOKEN || ''

  const l2Provider = new providers.StaticJsonRpcProvider(l2ProviderUrl)

  const l2Signer = new hre.ethers.Wallet(
    hre.network.config.accounts[signerIndex],
    l2Provider
  )

  const l2ChainId = await l2Signer.getChainId()
  let contractAddrs = CONTRACT_ADDRESSES[l2ChainId]
  if (l1ContractsJsonPath) {
    const data = await fs.readFile(l1ContractsJsonPath)
    const json = JSON.parse(data.toString())
    contractAddrs = {
      l1: {
        AddressManager: json.AddressManager,
        L1CrossDomainMessenger: json.L1CrossDomainMessengerProxy,
        L1StandardBridge: json.L1StandardBridgeProxy,
        StateCommitmentChain: ethers.constants.AddressZero,
        CanonicalTransactionChain: ethers.constants.AddressZero,
        BondManager: ethers.constants.AddressZero,
        OptimismPortal: json.OptimismPortalProxy,
        L2OutputOracle: json.L2OutputOracleProxy,
      },
      l2: DEFAULT_L2_CONTRACT_ADDRESSES,
    } as OEContractsLike
  }

  const messenger = new CrossChainMessenger({
    l1SignerOrProvider: signer,
    l2SignerOrProvider: l2Signer,
    l1ChainId: await signer.getChainId(),
    l2ChainId,
    nativeTokenAddress: l2NativeToken,
    bedrock: true,
    contracts: contractAddrs,
  })
  // Ensure deployments module is initialized
  await hre.deployments.all()

  console.log('Deploying WTON to L1')
  const WTON = await deployWTON(hre, signer, true)
  console.log(`Deployed to ${WTON.address}`)

  // Save WTON address to deployments
  await hre.deployments.save('WTON', {
    address: WTON.address,
    abi: [WTON.interface.format('json')],
  })

  console.log('Creating L2 WTON')
  const OptimismMintableERC20 = await createOptimismMintableERC20(
    WTON,
    l2Signer
  )

  // Save OptimismMintableERC20 address to deployments
  await hre.deployments.save('OptimismMintableERC20', {
    address: OptimismMintableERC20.address,
    abi: [OptimismMintableERC20.interface.format('json')],
  })

  const l1WTONAddress = (await hre.deployments.get('WTON')).address
  const l2WTONAddress = (await hre.deployments.get('OptimismMintableERC20'))
    .address
  console.log(l1WTONAddress, l2WTONAddress)

  console.log(`Approving WTON for deposit`)
  const approvalTx = await messenger.approveERC20(
    WTON.address,
    OptimismMintableERC20.address,
    hre.ethers.constants.MaxUint256
  )
  await approvalTx.wait()
  console.log('WTON approved')

  console.log('Depositing WTON to L2')
  const depositTx = await messenger.depositERC20(
    WTON.address,
    OptimismMintableERC20.address,
    utils.parseEther('1')
  )
  await depositTx.wait()
  console.log(`ERC20 deposited - ${depositTx.hash}`)

  console.log('Checking to make sure deposit was successful')
  // Deposit might get reorged, wait and also log for reorgs.
  let prevBlockHash: string = ''
  for (let i = 0; i < 12; i++) {
    const messageReceipt = await signer.provider!.getTransactionReceipt(
      depositTx.hash
    )
    if (messageReceipt.status !== 1) {
      console.log(`Deposit failed, retrying...`)
    }

    // Wait for stability, we want some amount of time after any reorg
    if (prevBlockHash !== '' && messageReceipt.blockHash !== prevBlockHash) {
      console.log(
        `Block hash changed from ${prevBlockHash} to ${messageReceipt.blockHash}`
      )
      i = 0
    } else if (prevBlockHash !== '') {
      console.log(`No reorg detected: ${i}`)
    }

    prevBlockHash = messageReceipt.blockHash
    await sleep(1000)
  }
  console.log(`Deposit confirmed`)

  const l2Balance = await OptimismMintableERC20.balanceOf(address)
  if (l2Balance.lt(utils.parseEther('1'))) {
    throw new Error(
      `bad deposit. recipient balance on L2: ${utils.formatEther(l2Balance)}`
    )
  }
  console.log(`Deposit success`)
}

const withdrawWTON = async (
  hre: HardhatRuntimeEnvironment,
  l2ProviderUrl: string,
  l1ContractsJsonPath: string,
  signerIndex: number
) => {
  const signers = await hre.ethers.getSigners()
  if (signers.length === 0) {
    throw new Error('No configured signers')
  }
  if (signerIndex < 0 || signers.length <= signerIndex) {
    throw new Error('Invalid signer index')
  }
  const signer = signers[signerIndex]
  const address = await signer.getAddress()
  console.log(`Using signer ${address}`)

  // Ensure that the signer has a balance before trying to
  // do anything
  const balance = await signer.getBalance()
  if (balance.eq(0)) {
    throw new Error('Signer has no balance')
  }

  const l2NativeToken = process.env.NATIVE_TOKEN || ''

  const l2Provider = new providers.StaticJsonRpcProvider(l2ProviderUrl)

  const l2Signer = new hre.ethers.Wallet(
    hre.network.config.accounts[signerIndex],
    l2Provider
  )

  const l2ChainId = await l2Signer.getChainId()
  let contractAddrs = CONTRACT_ADDRESSES[l2ChainId]
  if (l1ContractsJsonPath) {
    const data = await fs.readFile(l1ContractsJsonPath)
    const json = JSON.parse(data.toString())
    contractAddrs = {
      l1: {
        AddressManager: json.AddressManager,
        L1CrossDomainMessenger: json.L1CrossDomainMessengerProxy,
        L1StandardBridge: json.L1StandardBridgeProxy,
        StateCommitmentChain: ethers.constants.AddressZero,
        CanonicalTransactionChain: ethers.constants.AddressZero,
        BondManager: ethers.constants.AddressZero,
        OptimismPortal: json.OptimismPortalProxy,
        L2OutputOracle: json.L2OutputOracleProxy,
      },
      l2: DEFAULT_L2_CONTRACT_ADDRESSES,
    } as OEContractsLike
  }

  const OptimismPortal = new hre.ethers.Contract(
    contractAddrs.l1.OptimismPortal,
    Artifact__OptimismPortal.abi,
    signer
  )

  const L1CrossDomainMessenger = new hre.ethers.Contract(
    contractAddrs.l1.L1CrossDomainMessenger,
    Artifact__L1CrossDomainMessenger.abi,
    signer
  )

  const L1StandardBridge = new hre.ethers.Contract(
    contractAddrs.l1.L1StandardBridge,
    Artifact__L1StandardBridge.abi,
    signer
  )

  const L2OutputOracle = new hre.ethers.Contract(
    contractAddrs.l1.L2OutputOracle,
    Artifact__L2OutputOracle.abi,
    signer
  )

  const L2ToL1MessagePasser = new hre.ethers.Contract(
    predeploys.L2ToL1MessagePasser,
    Artifact__L2ToL1MessagePasser.abi
  )

  const L2CrossDomainMessenger = new hre.ethers.Contract(
    predeploys.L2CrossDomainMessenger,
    Artifact__L2CrossDomainMessenger.abi
  )

  const L2StandardBridge = new hre.ethers.Contract(
    predeploys.L2StandardBridge,
    Artifact__L2StandardBridge.abi
  )

  const messenger = new CrossChainMessenger({
    l1SignerOrProvider: signer,
    l2SignerOrProvider: l2Signer,
    l1ChainId: await signer.getChainId(),
    l2ChainId,
    nativeTokenAddress: l2NativeToken,
    bedrock: true,
    contracts: contractAddrs,
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
    signer
  )
  const OptimismMintableERC20 = new hre.ethers.Contract(
    l2WTONAddress,
    Artifact__OptimismMintableERC20Token.abi,
    l2Signer
  )

  console.log('Starting withdrawal')

  const preBalance = await WTON.balanceOf(signer.address)
  const withdraw = await messenger.withdrawERC20(
    WTON.address,
    OptimismMintableERC20.address,
    utils.parseEther('1')
  )
  const withdrawalReceipt = await withdraw.wait()
  for (const log of withdrawalReceipt.logs) {
    switch (log.address) {
      case L2ToL1MessagePasser.address: {
        const parsed = L2ToL1MessagePasser.interface.parseLog(log)
        console.log(`Log ${parsed.name} from ${log.address}`)
        console.log(parsed.args)
        console.log()
        break
      }
      case L2StandardBridge.address: {
        const parsed = L2StandardBridge.interface.parseLog(log)
        console.log(`Log ${parsed.name} from ${log.address}`)
        console.log(parsed.args)
        console.log()
        break
      }
      case L2CrossDomainMessenger.address: {
        const parsed = L2CrossDomainMessenger.interface.parseLog(log)
        console.log(`Log ${parsed.name} from ${log.address}`)
        console.log(parsed.args)
        console.log()
        break
      }
      default: {
        console.log(`Unknown log from ${log.address} - ${log.topics[0]}`)
      }
    }
  }

  setInterval(async () => {
    const currentStatus = await messenger.getMessageStatus(withdraw)
    console.log(`Message status: ${MessageStatus[currentStatus]}`)
    const latest = await L2OutputOracle.latestBlockNumber()
    console.log(`Latest L2OutputOracle commitment number: ${latest.toString()}`)
    const tip = await signer.provider!.getBlockNumber()
    console.log(`L1 chain tip: ${tip.toString()}`)
  }, 3000)

  const now = Math.floor(Date.now() / 1000)

  console.log('Waiting for message to be able to be proved')
  await messenger.waitForMessageStatus(withdraw, MessageStatus.READY_TO_PROVE)

  console.log('Proving withdrawal...')
  const prove = await messenger.proveMessage(withdraw)
  const proveReceipt = await prove.wait()
  console.log(proveReceipt)
  if (proveReceipt.status !== 1) {
    throw new Error('Prove withdrawal transaction reverted')
  }

  console.log('Waiting for message to be able to be relayed')
  await messenger.waitForMessageStatus(withdraw, MessageStatus.READY_FOR_RELAY)

  console.log('Finalizing withdrawal...')
  // TODO: Update SDK to properly estimate gas
  const finalize = await messenger.finalizeMessage(withdraw, {
    overrides: { gasLimit: 500_000 },
  })
  const finalizeReceipt = await finalize.wait()
  console.log('finalizeReceipt:', finalizeReceipt)
  console.log(`Took ${Math.floor(Date.now() / 1000) - now} seconds`)

  for (const log of finalizeReceipt.logs) {
    switch (log.address) {
      case OptimismPortal.address: {
        const parsed = OptimismPortal.interface.parseLog(log)
        console.log(`Log ${parsed.name} from OptimismPortal (${log.address})`)
        console.log(parsed.args)
        console.log()
        break
      }
      case L1CrossDomainMessenger.address: {
        const parsed = L1CrossDomainMessenger.interface.parseLog(log)
        console.log(
          `Log ${parsed.name} from L1CrossDomainMessenger (${log.address})`
        )
        console.log(parsed.args)
        console.log()
        break
      }
      case L1StandardBridge.address: {
        const parsed = L1StandardBridge.interface.parseLog(log)
        console.log(`Log ${parsed.name} from L1StandardBridge (${log.address})`)
        console.log(parsed.args)
        console.log()
        break
      }
      case WTON.address: {
        const parsed = WTON.interface.parseLog(log)
        console.log(`Log ${parsed.name} from WTON (${log.address})`)
        console.log(parsed.args)
        console.log()
        break
      }
      default:
        console.log(
          `Unknown log emitted from ${log.address} - ${log.topics[0]}`
        )
    }
  }

  const postBalance = await WTON.balanceOf(signer.address)

  const expectedBalance = preBalance.add(utils.parseEther('1'))
  if (!expectedBalance.eq(postBalance)) {
    throw new Error(
      `Balance mismatch, expected: ${expectedBalance}, actual: ${postBalance}`
    )
  }
  console.log('Withdrawal success')
}

// TODO(tynes): this task could be modularized in the future
// so that it can deposit an arbitrary token. Right now it
// deploys a WTON contract, mints some WTON and then
// deposits that into L2 through the StandardBridge.
task('deposit-erc20', 'Deposit WTON onto L2.')
  .addParam(
    'l2ProviderUrl',
    'L2 provider URL.',
    'http://localhost:9545',
    types.string
  )
  .addOptionalParam(
    'l1ContractsJsonPath',
    'Path to a JSON with L1 contract addresses in it',
    '',
    types.string
  )
  .addOptionalParam('signerIndex', 'Index of signer to use', 0, types.int)
  .setAction(async (args, hre) => {
    await depositWTON(
      hre,
      args.l2ProviderUrl,
      args.l1ContractsJsonPath,
      args.signerIndex
    )
  })

task('withdraw-erc20', 'Withdraw WTON from L2 to L1')
  .addParam(
    'l2ProviderUrl',
    'L2 provider URL.',
    'http://localhost:9545',
    types.string
  )
  .addOptionalParam(
    'l1ContractsJsonPath',
    'Path to a JSON with L1 contract addresses in it',
    '',
    types.string
  )
  .addOptionalParam('signerIndex', 'Index of signer to use', 0, types.int)
  .setAction(async (args, hre) => {
    await withdrawWTON(
      hre,
      args.l2ProviderUrl,
      args.l1ContractsJsonPath,
      args.signerIndex
    )
  })
