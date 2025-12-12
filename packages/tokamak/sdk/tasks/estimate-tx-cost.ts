import { task } from 'hardhat/config'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BytesLike, ethers } from 'ethers'

import { asL2Provider, L2Provider } from '../src'

console.log('Setup task...')

// Lazy initialization to avoid module-level errors
let _l2Provider: L2Provider<ethers.providers.StaticJsonRpcProvider> | null =
  null
let _l2Wallet: ethers.Wallet | null = null

const getL2Provider = () => {
  if (!_l2Provider) {
    const baseProvider = new ethers.providers.StaticJsonRpcProvider(
      process.env.L2_URL
    )
    _l2Provider = asL2Provider(baseProvider)
  }
  return _l2Provider
}

const getL2Wallet = () => {
  const privateKey = process.env.PRIVATE_KEY as BytesLike
  if (!privateKey) {
    throw new Error('PRIVATE_KEY environment variable is not set')
  }
  const provider = getL2Provider()
  if (!_l2Wallet) {
    _l2Wallet = new ethers.Wallet(privateKey, provider)
  }
  return _l2Wallet
}

const estimateL1Gas = async () => {
  const l2Wallet = getL2Wallet()
  const provider = getL2Provider()

  const tx = await l2Wallet.populateTransaction({
    to: '0x1000000000000000000000000000000000000000',
    value: ethers.utils.parseEther('0.01'),
    gasPrice: await provider.getGasPrice(),
  })
  const l1CostEstimate = await provider.estimateL1GasCost(tx)
  console.log(ethers.utils.formatEther(l1CostEstimate))
  const gasLimit = tx.gasLimit
  if (!gasLimit) {
    console.error(`gasLimit undefined`)
    return
  }
  const gasPrice = tx.maxFeePerGas
  const l2CostEstimate = gasLimit.mul(gasPrice)
  console.log(`L2 cost estimated: ${ethers.utils.formatEther(l2CostEstimate)}`)

  const totalSum = l2CostEstimate.add(l1CostEstimate)
  console.log(`Total sum: ${ethers.utils.formatEther(totalSum)}`)

  const res = await l2Wallet.sendTransaction(tx)
  const receipt = await res.wait()
  console.log(`Receipt: ${JSON.stringify(receipt)}`)

  const l2CostActual = receipt.gasUsed.mul(receipt.effectiveGasPrice)
  console.log(`L2 cost actual: ${ethers.utils.formatEther(l2CostActual)}`)

  const l1CostActual = receipt.l1Fee
  console.log(`L1 cost actual: ${ethers.utils.formatEther(l1CostActual)}`)

  const totalActual = l2CostActual.add(l1CostActual)
  console.log(`Total cost actual: ${ethers.utils.formatEther(totalActual)}`)

  const difference = totalActual.sub(totalSum).abs()
  console.log(
    `Difference between actual and estimate: ${ethers.utils.formatEther(
      difference
    )}`
  )
}

task('estimate-tx-cost', 'Estimate transaction cost.').setAction(async () => {
  await estimateL1Gas()
})
