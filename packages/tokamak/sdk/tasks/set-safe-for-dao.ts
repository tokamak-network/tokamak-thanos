import { ethers } from 'ethers'
import { task } from 'hardhat/config'
import { executeContractCallWithSigners } from '@tokamak-network/thanos-contracts/lib/safe-contracts/src/utils/execution'

import { getDAOMembers } from '../src/utils/owners'

/**
 * Adds the specified owner to the Gnosis Safe and verifies that the owner has been added.
 *
 * @param safeContract - The Gnosis Safe contract instance (with signer connected)
 * @param owner - The owner address to add
 * @param threshold - The threshold value to apply
 * @param signers - An array of signers (typically a single signer)
 */
export const addOwnerAndVerify = async (
  safeContract: ethers.Contract,
  owner: string,
  threshold: number,
  signers: ethers.Signer[]
): Promise<void> => {
  try {
    // Execute addOwnerWithThreshold
    const tx = await executeContractCallWithSigners(
      safeContract,
      safeContract,
      'addOwnerWithThreshold',
      [owner, threshold],
      signers
    )
    console.log(`Tx Hash for adding owner ${owner}:`, tx.hash)

    // Wait for transaction and check status
    const receipt = await tx.wait()
    if (receipt.status !== 1) {
      throw new Error(
        `Transaction failed for adding owner ${owner}. Tx Hash: ${tx.hash}`
      )
    }

    // Get Safe owners and verify
    const safeOwners = await safeContract.getOwners()
    console.log(`Safe owners after adding owner ${owner}:`, safeOwners)
    if (!safeOwners.includes(owner)) {
      throw new Error(
        `Verification failed: Owner ${owner} was not added to the safe.`
      )
    }
    console.log(`Successfully added owner: ${owner}`)
  } catch (error) {
    console.error(`Error adding owner ${owner}:`, error)
    throw error // re-throw the error for handling to upper level
  }
}

// Task
task('set-safe-wallet', 'Set Safe Wallet for the Tokamak DAO').setAction(
  async () => {
    try {
      // Verify ENVs
      const l1Url = process.env.L1_URL
      const privateKey = process.env.PRIVATE_KEY
      const safeWalletAddress = process.env.SAFE_WALLET_ADDRESS

      if (!l1Url || !privateKey || !safeWalletAddress) {
        throw new Error(
          'Missing required environment variables: L1_URL, PRIVATE_KEY, SAFE_WALLET_ADDRESS'
        )
      }

      const l1Provider = new ethers.providers.StaticJsonRpcProvider(l1Url)
      const network = await l1Provider.getNetwork()
      console.log('L1 Chain ID:', network.chainId)

      // Create the signer
      const signer = new ethers.Wallet(privateKey, l1Provider)

      // Get the desigated owners depending network
      const designatedOwners = getDAOMembers(network.chainId)
      if (!designatedOwners || designatedOwners.length < 2) {
        throw new Error(
          'Insufficient DAO members returned for designated owners'
        )
      }

      // ABIs of Gnosis Safe
      const gnosisSafeAbi = [
        'function getThreshold() view returns (uint256)',
        'function addOwnerWithThreshold(address owner, uint256 _threshold) external',
        'function changeThreshold(uint256 _threshold) external',
        'function execTransaction(address to, uint256 value, bytes calldata data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, bytes calldata signatures) external returns (bool success)',
        'function nonce() view returns (uint256)',
        'function getOwners() view returns (address[])',
      ]

      // Create Gnosis safe contract instance
      const gnosisSafeContract = new ethers.Contract(
        safeWalletAddress,
        gnosisSafeAbi,
        signer
      )
      console.log('Gnosis Safe Contract:', gnosisSafeContract.address)

      // Execute: Add owners to Safe
      await addOwnerAndVerify(gnosisSafeContract, designatedOwners[0], 1, [
        signer,
      ])
      await addOwnerAndVerify(gnosisSafeContract, designatedOwners[1], 1, [
        signer,
      ])

      // Change threshold to 3
      const txChangeThreshold = await executeContractCallWithSigners(
        gnosisSafeContract,
        gnosisSafeContract,
        'changeThreshold',
        [3],
        [signer]
      )
      console.log('Tx Hash for changing threshold:', txChangeThreshold.hash)
      const receiptChangeThreshold = await txChangeThreshold.wait()
      if (receiptChangeThreshold.status !== 1) {
        throw new Error(
          `Threshold change transaction failed. Tx Hash: ${txChangeThreshold.hash}`
        )
      }
      const newSafeThreshold = await gnosisSafeContract.getThreshold()
      console.log('New threshold:', String(newSafeThreshold))
    } catch (error) {
      console.error('Got the error running Safe contract:', error)
      process.exit(1) // exit
    }
  }
)
