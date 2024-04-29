import { task } from 'hardhat/config'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'
import { BytesLike, ethers } from 'ethers'

import { asL2Provider } from '../src'

console.log('Setup task...')

const privateKey = process.env.PRIVATE_KEY as BytesLike

const l2Provider = new ethers.providers.StaticJsonRpcProvider(
  process.env.L2_URL
)
const provider = asL2Provider(l2Provider)

const l2Wallet = new ethers.Wallet(privateKey, provider)

const estimateL1Gas = async () => {
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
