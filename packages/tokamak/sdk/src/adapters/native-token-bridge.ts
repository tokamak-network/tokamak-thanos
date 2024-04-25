/* eslint-disable @typescript-eslint/no-unused-vars */
import {
  ethers,
  Contract,
  Overrides,
  Signer,
  BigNumber,
  CallOverrides,
} from 'ethers'
import {
  TransactionRequest,
  TransactionResponse,
  BlockTag,
} from '@ethersproject/abstract-provider'
import optimismMintableERC20 from '@tokamak-network/thanos-contracts/forge-artifacts/OptimismMintableERC20.sol/OptimismMintableERC20.json'

import {
  TokenBridgeMessage,
  NumberLike,
  AddressLike,
  MessageDirection,
} from '../interfaces'
import { toAddress, omit } from '../utils'
import { StandardBridgeAdapter } from './standard-bridge'

/**
 * Bridge adapter for any token bridge that uses the standard token bridge interface.
 */
export class NativeTokenBridgeAdapter extends StandardBridgeAdapter {
  public async getDepositsByAddress(
    address: AddressLike,
    opts?: {
      fromBlock?: BlockTag
      toBlock?: BlockTag
    }
  ): Promise<TokenBridgeMessage[]> {
    const events = await this.l1Bridge.queryFilter(
      this.l1Bridge.filters.ERC20DepositInitiated(
        undefined,
        undefined,
        address
      ),
      opts?.fromBlock,
      opts?.toBlock
    )

    return events
      .filter((event) =>
        this.supportsTokenPair(event.args.l1Token, event.args.l2Token)
      )
      .map((event) => {
        return {
          direction: MessageDirection.L1_TO_L2,
          from: event.args.from,
          to: event.args.to,
          l1Token: event.args.l1Token,
          l2Token: event.args.l2Token,
          amount: event.args.amount,
          data: event.args.extraData,
          logIndex: event.logIndex,
          blockNumber: event.blockNumber,
          transactionHash: event.transactionHash,
        }
      })
      .sort((a, b) => {
        // Sort descending by block number
        return b.blockNumber - a.blockNumber
      })
  }

  public async getWithdrawalsByAddress(
    address: AddressLike,
    opts?: {
      fromBlock?: BlockTag
      toBlock?: BlockTag
    }
  ): Promise<TokenBridgeMessage[]> {
    const events = await this.l2Bridge.queryFilter(
      this.l2Bridge.filters.WithdrawalInitiated(undefined, undefined, address),
      opts?.fromBlock,
      opts?.toBlock
    )

    return events
      .filter((event) =>
        this.supportsTokenPair(event.args.l1Token, event.args.l2Token)
      )
      .map((event) => {
        return {
          direction: MessageDirection.L2_TO_L1,
          from: event.args.from,
          to: event.args.to,
          l1Token: event.args.l1Token,
          l2Token: event.args.l2Token,
          amount: event.args.amount,
          data: event.args.extraData,
          logIndex: event.logIndex,
          blockNumber: event.blockNumber,
          transactionHash: event.transactionHash,
        }
      })
      .sort((a, b) => {
        // Sort descending by block number
        return b.blockNumber - a.blockNumber
      })
  }

  public async supportsTokenPair(
    l1Token: AddressLike,
    l2Token: AddressLike
  ): Promise<boolean> {
    // Only support L2 native token deposits and withdrawals.
    return this.filterL2NativeTokenDepositsAndWithdrawls(l1Token, l2Token)
  }

  public async approval(
    l1Token: AddressLike,
    l2Token: AddressLike,
    signer: ethers.Signer
  ): Promise<BigNumber> {
    if (!(await this.supportsTokenPair(l1Token, l2Token))) {
      throw new Error(`token pair not supported by bridge`)
    }

    const token = new Contract(
      toAddress(l1Token),
      optimismMintableERC20.abi,
      this.messenger.l1Provider
    )

    return token.allowance(await signer.getAddress(), this.l1Bridge.address)
  }

  public async approve(
    l1Token: AddressLike,
    l2Token: AddressLike,
    amount: NumberLike,
    signer: Signer,
    opts?: {
      overrides?: Overrides
    }
  ): Promise<TransactionResponse> {
    return signer.sendTransaction(
      await this.populateTransaction.approve(l1Token, l2Token, amount, opts)
    )
  }

  public async deposit(
    l1Token: AddressLike,
    l2Token: AddressLike,
    amount: NumberLike,
    signer: Signer,
    opts?: {
      recipient?: AddressLike
      l2GasLimit?: NumberLike
      overrides?: Overrides
    }
  ): Promise<TransactionResponse> {
    return signer.sendTransaction(
      await this.populateTransaction.deposit(l1Token, l2Token, amount, opts)
    )
  }

  public async withdraw(
    l1Token: AddressLike,
    l2Token: AddressLike,
    amount: NumberLike,
    signer: Signer,
    opts?: {
      recipient?: AddressLike
      overrides?: Overrides
    }
  ): Promise<TransactionResponse> {
    return signer.sendTransaction(
      await this.populateTransaction.withdraw(l1Token, l2Token, amount, opts)
    )
  }

  populateTransaction = {
    approve: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        overrides?: Overrides
      }
    ): Promise<TransactionRequest> => {
      if (!(await this.supportsTokenPair(l1Token, l2Token))) {
        throw new Error(`token pair not supported by bridge`)
      }

      const token = new Contract(
        toAddress(l1Token),
        optimismMintableERC20.abi,
        this.messenger.l1Provider
      )

      return token.populateTransaction.approve(
        this.l1Bridge.address,
        amount,
        opts?.overrides || {}
      )
    },

    deposit: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike
        l2GasLimit?: NumberLike
        overrides?: Overrides
      }
    ): Promise<TransactionRequest> => {
      if (!(await this.supportsTokenPair(l1Token, l2Token))) {
        throw new Error(`token pair not supported by bridge`)
      }

      // NOTE: don't set the "value" because on ETH is native token in L1
      if (opts?.recipient === undefined) {
        return this.l1Bridge.populateTransaction.depositERC20(
          toAddress(l1Token),
          toAddress(l2Token),
          amount,
          opts?.l2GasLimit || 200_000, // Default to 200k gas limit.
          '0x', // No data.
          opts?.overrides || {}
        )
      } else {
        return this.l1Bridge.populateTransaction.depositERC20To(
          toAddress(l1Token),
          toAddress(l2Token),
          toAddress(opts.recipient),
          amount,
          opts?.l2GasLimit || 200_000, // Default to 200k gas limit.
          '0x', // No data.
          opts?.overrides || {}
        )
      }
    },

    withdraw: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike
        overrides?: Overrides
      }
    ): Promise<TransactionRequest> => {
      if (!(await this.supportsTokenPair(l1Token, l2Token))) {
        throw new Error(`token pair not supported by bridge`)
      }

      // we set this value to "amount" because we are withdrawing the native token
      if (opts?.recipient === undefined) {
        return this.l2Bridge.populateTransaction.withdraw(
          toAddress(l2Token),
          amount,
          0, // L1 gas not required.
          '0x', // No data.
          {
            ...omit(opts?.overrides || {}, 'value'),
            value: amount,
          }
        )
      } else {
        return this.l2Bridge.populateTransaction.withdrawTo(
          toAddress(l2Token),
          toAddress(opts.recipient),
          amount,
          0, // L1 gas not required.
          '0x', // No data.
          {
            ...omit(opts?.overrides || {}, 'value'),
            value: amount,
          }
        )
      }
    },
  }

  estimateGas = {
    approve: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        overrides?: CallOverrides
      }
    ): Promise<BigNumber> => {
      return this.messenger.l1Provider.estimateGas(
        await this.populateTransaction.approve(l1Token, l2Token, amount, opts)
      )
    },

    deposit: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike
        l2GasLimit?: NumberLike
        overrides?: CallOverrides
      }
    ): Promise<BigNumber> => {
      return this.messenger.l1Provider.estimateGas(
        await this.populateTransaction.deposit(l1Token, l2Token, amount, opts)
      )
    },

    withdraw: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike
        overrides?: CallOverrides
      }
    ): Promise<BigNumber> => {
      return this.messenger.l2Provider.estimateGas(
        await this.populateTransaction.withdraw(l1Token, l2Token, amount, opts)
      )
    },
  }
}
