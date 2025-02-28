import { ethers } from 'ethers'
import { task, types } from 'hardhat/config'
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
  // Execute addOwnerWithThreshold
  const tx = await executeContractCallWithSigners(
    safeContract,
    safeContract,
    'addOwnerWithThreshold',
    [owner, threshold],
    signers
  )
  console.log(`Tx Hash for adding owner ${owner}:`, tx.hash)
  await tx.wait()

  // Get Safe owner
  const safeOwners = await safeContract.getOwners()
  console.log(`Safe owners after adding owner ${owner}:`, safeOwners)
  if (safeOwners.includes(owner)) {
    console.log(`Successfully added owner: ${owner}`)
  } else {
    console.log(`Failed to add owner: ${owner}`)
  }
}

// Task
task('set-safe-wallet', 'Set Safe Wallet for the TokamakDAO')
  .addParam('rpc', 'L1 RPC endpoint', '', types.string)
  .addParam('chainid', 'L1 chain id', '', types.int)
  .addParam('privatekey', 'Admin Private key', '', types.string)
  .addParam('address', 'Gnosis Safe contract address', '', types.string)
  .setAction(async (args) => {
    const l1Provider = new ethers.providers.StaticJsonRpcProvider(args.rpc)

    // Create the signer
    // TODO: Get the Admin's private key from the seed phrase (user's input)
    const ownerAPrivateKey = args.privatekey
    const signer = new ethers.Wallet(ownerAPrivateKey, l1Provider)

    // Get predefined owners
    const owners = getDAOMembers(args.chainid)

    // ABIs of Gnosis Safe
    const gnosisSafeAbi = [
      'function getThreshold() view returns (uint256)',
      'function addOwnerWithThreshold(address owner, uint256 _threshold) external',
      'function changeThreshold(uint256 _threshold) external',
      'function execTransaction(address to, uint256 value, bytes calldata data, uint8 operation, uint256 safeTxGas, uint256 baseGas, uint256 gasPrice, address gasToken, address refundReceiver, bytes calldata signatures) external returns (bool success)',
      'function nonce() view returns (uint256)',
      'function getOwners() view returns (address[])',
    ]

    // Create contract instance
    const gnosisSafeContract = new ethers.Contract(
      args.address,
      gnosisSafeAbi,
      signer
    )
    console.log('Gnosis Safe Contract:', gnosisSafeContract.address)

    // Execute
    try {
      // Add 2 owners to Safe
      await addOwnerAndVerify(gnosisSafeContract, owners[0], 1, [signer])
      await addOwnerAndVerify(gnosisSafeContract, owners[1], 1, [signer])

      // Change threshold to 3
      const tx3 = await executeContractCallWithSigners(
        gnosisSafeContract,
        gnosisSafeContract,
        'changeThreshold',
        [3],
        [signer]
      )
      console.log('Tx 3 Hash:', tx3.hash)
      await tx3.wait()
      const newSafeThreshold = await gnosisSafeContract.getThreshold()
      console.log('New threshold:', String(newSafeThreshold))
    } catch (error) {
      console.error('Got the error running Safe contract:', error)
    }
  })
