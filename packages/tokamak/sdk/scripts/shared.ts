// import hardhat from 'hardhat'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import { BigNumber, Wallet, Contract, Event, utils } from 'ethers'
import { predeploys } from '@eth-optimism/core-utils'
import Artifact__OptimismMintableERC20TokenFactory from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json'
import Artifact__OptimismMintableERC20Token from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'

import Artifact__ERC20 from '../../contracts-bedrock/forge-artifacts/MockERC20.sol/MockERC20.json'
// import Artifact__L2NativeToken from '../../contracts-bedrock/forge-artifacts/L2NativeToken.sol/L2NativeToken.json'
import Artifact__MockHello from '../../contracts-bedrock/forge-artifacts/MockHello.sol/MockHello.json'

export const deployERC20 = async (
  hre: HardhatRuntimeEnvironment,
  signer: Wallet,
  name: string,
  symbol: string,
  total: BigNumber
): Promise<Contract> => {
  const Factory__ERC20 = new hre.ethers.ContractFactory(
    Artifact__ERC20.abi,
    Artifact__ERC20.bytecode.object,
    signer
  )

  console.log('Sending deployment transaction')
  const token = await Factory__ERC20.deploy(name, symbol, signer.address, total)
  const receipt = await token.deployTransaction.wait()
  console.log(`token deployed: ${receipt.transactionHash}`)

  return token
}

export const createOptimismMintableERC20 = async (
  hre: HardhatRuntimeEnvironment,
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

  const l2WethAddress = event.args.localToken
  console.log(`Deployed to ${l2WethAddress}`)

  return new Contract(
    l2WethAddress,
    Artifact__OptimismMintableERC20Token.abi,
    l2Signer
  )
}

export const deployHello = async (
  hre: HardhatRuntimeEnvironment,
  signer: Wallet
): Promise<Contract> => {
  const Factory__Hello = new hre.ethers.ContractFactory(
    Artifact__MockHello.abi,
    Artifact__MockHello.bytecode.object,
    signer
  )

  console.log('Sending deployment transaction')
  const hello = await Factory__Hello.deploy()
  const receipt = await hello.deployTransaction.wait()
  console.log(`Hello deployed: ${receipt.transactionHash}`)

  return hello
}

export const getMessageOfHello = async (helloContract) => {
  const blockNumber = await helloContract.blockNumber()
  const message = await helloContract.message()

  return {
    blockNumber,
    message,
  }
}

export const differenceLog = async (prevBalance, afterBalance) => {
  await differenceTonBalance(
    prevBalance.l1Balance,
    afterBalance.l1Balance,
    'L1 Wallet TON Changed : '
  )
  await differenceTonBalance(
    prevBalance.l2Balance,
    afterBalance.l2Balance,
    'L2 Wallet TON Changed : '
  )

  await differenceEthBalance(
    prevBalance.l1Balance,
    afterBalance.l1Balance,
    'L1 Wallet ETH Changed : '
  )
  await differenceEthBalance(
    prevBalance.l2Balance,
    afterBalance.l2Balance,
    'L2 Wallet ETH Changed : '
  )

  await differenceTonBalance(
    prevBalance.l1BridgeBalance,
    afterBalance.l1BridgeBalance,
    'l1BridgeBalance TON Changed : '
  )
  await differenceTonBalance(
    prevBalance.l1MessengerBalance,
    afterBalance.l1MessengerBalance,
    'l1CrossDomainMessenger TON Changed : '
  )
  await differenceTonBalance(
    prevBalance.OptomismPortalBalance,
    afterBalance.OptomismPortalBalance,
    'OptomismPortalContract TON Changed : '
  )
  await differenceErc20Balance(
    prevBalance.portal,
    afterBalance.portal,
    'OptomismPortal depositAmount Changed : '
  )
}

export const getBalances = async (
  l1wallet,
  l2wallet,
  tonContract,
  l2EthContract,
  l1BridgeContract,
  l1CrossDomainMessengerContract,
  OptomismPortalContract
) => {
  const l1Balance = await getL1Balance(l1wallet, tonContract)
  const l2Balance = await getL2Balance(l2wallet, l2EthContract)

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

  return {
    l1Balance,
    l2Balance,
    l1BridgeBalance,
    l1MessengerBalance,
    OptomismPortalBalance,
    portal,
  }
}

export const getL1Balance = async (account: Wallet, tonContract) => {
  const tonBalance = await tonContract.balanceOf(account.address)
  const ethBalance = await account.getBalance()
  return {
    tonBalance,
    ethBalance,
  }
}

export const getL1ContractBalance = async (contract, tonContract) => {
  const tonBalance = await tonContract.balanceOf(contract.address)
  return {
    tonBalance,
  }
}

export const getL2Balance = async (account: Wallet, l2EthContract) => {
  const tonBalance = await account.getBalance()
  const ethBalance = await l2EthContract.balanceOf(account.address)

  return {
    tonBalance,
    ethBalance,
  }
}

export const getErc20Balance = async (account: Wallet, token) => {
  const balance = await token.balanceOf(account.address)
  return {
    balance,
  }
}

export const getPortalDepositedAmount = async (portal) => {
  const balance = await portal.depositedAmount()
  return {
    balance,
  }
}

export const differenceTonBalance = async (
  l1Log: any,
  l2Log: any,
  title: string
) => {
  console.log(title, utils.formatEther(l2Log.tonBalance.sub(l1Log.tonBalance)))
}

export const differenceEthBalance = async (
  l1Log: any,
  l2Log: any,
  title: string
) => {
  console.log(title, utils.formatEther(l2Log.ethBalance.sub(l1Log.ethBalance)))
}

export const differenceErc20Balance = async (
  l1Log: any,
  l2Log: any,
  title: string
) => {
  console.log(title, utils.formatEther(l2Log.balance.sub(l1Log.balance)))
}

export const logEvent = async (receipt, topic, con, title) => {
  // const topic = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessage');
  // const topic1 = l1CrossDomainMessengerContract.interface.getEventTopic('SentMessageExtension1');
  const log = receipt.logs.find((x) => x.topics.indexOf(topic) >= 0)
  const deployedEvent = con.interface.parseLog(log)
  console.log(title, deployedEvent)
}

export const erc20ABI = [
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
    inputs: [
      { internalType: 'address', name: '_owner', type: 'address' },
      { internalType: 'address', name: '_spender', type: 'address' },
    ],
    name: 'allowance',
    outputs: [
      {
        internalType: 'uint256',
        name: '',
        type: 'uint256',
      },
    ],
    stateMutability: 'view',
    type: 'function',
  },
  {
    constant: true,
    inputs: [{ name: 'account', type: 'address' }],
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
    inputs: [],
    name: 'totalSupply',
    outputs: [
      {
        internalType: 'uint256',
        name: '',
        type: 'uint256',
      },
    ],
    stateMutability: 'view',
    type: 'function',
  },
  {
    inputs: [
      {
        internalType: 'address',
        name: 'spender',
        type: 'address',
      },
      {
        internalType: 'uint256',
        name: 'amount',
        type: 'uint256',
      },
      {
        internalType: 'bytes',
        name: 'data',
        type: 'bytes',
      },
    ],
    name: 'approveAndCall',
    outputs: [
      {
        internalType: 'bool',
        name: '',
        type: 'bool',
      },
    ],
    stateMutability: 'nonpayable',
    type: 'function',
  },
]

export const wtonABI = [
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
    inputs: [{ name: 'account', type: 'address' }],
    name: 'balanceOf',
    outputs: [{ name: 'balance', type: 'uint256' }],
    type: 'function',
  },
  {
    inputs: [],
    name: 'totalSupply',
    outputs: [
      {
        internalType: 'uint256',
        name: '',
        type: 'uint256',
      },
    ],
    stateMutability: 'view',
    type: 'function',
  },
  {
    inputs: [],
    name: 'deposit',
    outputs: [],
    stateMutability: 'payable',
    type: 'function',
  },
  {
    inputs: [{ name: 'wad', type: 'uint' }],
    name: 'withdraw',
    outputs: [],
    stateMutability: 'nonpayable',
    type: 'function',
  },
]
