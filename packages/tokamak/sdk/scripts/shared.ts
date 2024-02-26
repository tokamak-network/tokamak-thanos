// import hardhat from 'hardhat'
import { HardhatRuntimeEnvironment } from 'hardhat/types'
import { BigNumber, Wallet, Contract, Event } from 'ethers'
import { predeploys } from '@eth-optimism/core-utils'
import Artifact__OptimismMintableERC20TokenFactory from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20Factory.sol/OptimismMintableERC20Factory.json'
import Artifact__OptimismMintableERC20Token from '@eth-optimism/contracts-bedrock/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'

import Artifact__ERC20 from '../../contracts-bedrock/forge-artifacts/MockERC20.sol/MockERC20.json'

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
