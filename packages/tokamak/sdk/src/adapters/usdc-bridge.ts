import { hexStringEquals } from '@tokamak-network/core-utils'
import { Contract } from 'ethers'
import l1UsdcBridgeArtifact from '@tokamak-network/thanos-contracts/forge-artifacts/L1UsdcBridge.sol/L1UsdcBridge.json'
import l2UsdcBridgeArtifact from '@tokamak-network/thanos-contracts/forge-artifacts/L2UsdcBridge.sol/L2UsdcBridge.json'

import { AddressLike } from '../interfaces'
import { toAddress } from '../utils'
import { StandardBridgeAdapter } from './standard-bridge'
import { CrossChainMessenger } from '../cross-chain-messenger'

export class USDCBridgeAdapter extends StandardBridgeAdapter {
  constructor(opts: {
    messenger: CrossChainMessenger
    l1Bridge: AddressLike
    l2Bridge: AddressLike
  }) {
    super(opts)
    this.l1Bridge = new Contract(
      toAddress(opts.l1Bridge),
      l1UsdcBridgeArtifact.abi,
      this.messenger.l1Provider
    )
    this.l2Bridge = new Contract(
      toAddress(opts.l2Bridge),
      l2UsdcBridgeArtifact.abi,
      this.messenger.l2Provider
    )
  }

  public async supportsTokenPair(
    l1Token: AddressLike,
    l2Token: AddressLike
  ): Promise<boolean> {
    // Just need access to this ABI for this one function.
    const l1Bridge = new Contract(
      this.l1Bridge.address,
      [
        {
          inputs: [],
          name: 'l1Usdc',
          outputs: [
            {
              internalType: 'address',
              name: '',
              type: 'address',
            },
          ],
          stateMutability: 'view',
          type: 'function',
        },
        {
          inputs: [],
          name: 'l2Usdc',
          outputs: [
            {
              internalType: 'address',
              name: '',
              type: 'address',
            },
          ],
          stateMutability: 'view',
          type: 'function',
        },
      ],
      this.messenger.l1Provider
    )

    const allowedL1Token = await l1Bridge.l1Usdc()
    if (!hexStringEquals(allowedL1Token, toAddress(l1Token))) {
      return false
    }

    const allowedL2Token = await l1Bridge.l2Usdc()
    if (!hexStringEquals(allowedL2Token, toAddress(l2Token))) {
      return false
    }

    return true
  }
}
