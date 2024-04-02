/* eslint-disable @typescript-eslint/no-unused-vars */
import {
  Provider,
  TransactionResponse,
  TransactionRequest,
  TransactionReceipt,
} from '@ethersproject/abstract-provider'
import { Signer } from '@ethersproject/abstract-signer'
import { BigNumber, ethers } from 'ethers'
import { sleep, DepositTx } from '@tokamak-network/core-utils'

import {
  DepositTransactionRequest,
  OEContracts,
  OEContractsLike,
  NumberLike,
  SignerOrProviderLike,
  ProvenWithdrawal,
} from './interfaces'
import {
  toSignerOrProvider,
  toNumber,
  DeepPartial,
  getPortalsContracts,
  DEPOSIT_CONFIRMATION_BLOCKS,
  CHAIN_BLOCK_TIMES,
} from './utils'

export class Portals {
  /**
   * Provider connected to the L1 chain.
   */
  public l1SignerOrProvider: Signer | Provider

  /**
   * Provider connected to the L2 chain.
   */
  public l2SignerOrProvider: Signer | Provider

  /**
   * Chain ID for the L1 network.
   */
  public l1ChainId: number

  /**
   * Chain ID for the L2 network.
   */
  public l2ChainId: number

  /**
   * Contract objects attached to their respective providers and addresses.
   */
  public contracts: OEContracts

  /**
   * Number of blocks before a deposit is considered confirmed.
   */
  public depositConfirmationBlocks: number

  /**
   * Estimated average L1 block time in seconds.
   */
  public l1BlockTimeSeconds: number

  /**
   * Creates a new CrossChainProvider instance.
   *
   * @param opts Options for the provider.
   * @param opts.l1SignerOrProvider Signer or Provider for the L1 chain, or a JSON-RPC url.
   * @param opts.l2SignerOrProvider Signer or Provider for the L2 chain, or a JSON-RPC url.
   * @param opts.l1ChainId Chain ID for the L1 chain.
   * @param opts.l2ChainId Chain ID for the L2 chain.
   * @param opts.depositConfirmationBlocks Optional number of blocks before a deposit is confirmed.
   * @param opts.l1BlockTimeSeconds Optional estimated block time in seconds for the L1 chain.
   * @param opts.contracts Optional contract address overrides.
   */
  constructor(opts: {
    l1SignerOrProvider: SignerOrProviderLike
    l2SignerOrProvider: SignerOrProviderLike
    l1ChainId: NumberLike
    l2ChainId: NumberLike
    depositConfirmationBlocks?: NumberLike
    l1BlockTimeSeconds?: NumberLike
    contracts?: DeepPartial<OEContractsLike>
  }) {
    this.l1SignerOrProvider = toSignerOrProvider(opts.l1SignerOrProvider)
    this.l2SignerOrProvider = toSignerOrProvider(opts.l2SignerOrProvider)

    try {
      this.l1ChainId = toNumber(opts.l1ChainId)
    } catch (err) {
      throw new Error(`L1 chain ID is missing or invalid: ${opts.l1ChainId}`)
    }

    try {
      this.l2ChainId = toNumber(opts.l2ChainId)
    } catch (err) {
      throw new Error(`L2 chain ID is missing or invalid: ${opts.l2ChainId}`)
    }

    this.depositConfirmationBlocks =
      opts?.depositConfirmationBlocks !== undefined
        ? toNumber(opts.depositConfirmationBlocks)
        : DEPOSIT_CONFIRMATION_BLOCKS[this.l2ChainId] || 0

    this.l1BlockTimeSeconds =
      opts?.l1BlockTimeSeconds !== undefined
        ? toNumber(opts.l1BlockTimeSeconds)
        : CHAIN_BLOCK_TIMES[this.l1ChainId] || 1

    this.contracts = getPortalsContracts(this.l2ChainId, {
      l1SignerOrProvider: this.l1SignerOrProvider,
      l2SignerOrProvider: this.l2SignerOrProvider,
      overrides: opts.contracts,
    })
  }

  /**
   * Provider connected to the L1 chain.
   */
  get l1Provider(): Provider {
    if (Provider.isProvider(this.l1SignerOrProvider)) {
      return this.l1SignerOrProvider
    } else {
      return this.l1SignerOrProvider.provider
    }
  }

  /**
   * Provider connected to the L2 chain.
   */
  get l2Provider(): Provider {
    if (Provider.isProvider(this.l2SignerOrProvider)) {
      return this.l2SignerOrProvider
    } else {
      return this.l2SignerOrProvider.provider
    }
  }

  /**
   * Signer connected to the L1 chain.
   */
  get l1Signer(): Signer {
    if (Provider.isProvider(this.l1SignerOrProvider)) {
      throw new Error(`messenger has no L1 signer`)
    } else {
      return this.l1SignerOrProvider
    }
  }

  /**
   * Signer connected to the L2 chain.
   */
  get l2Signer(): Signer {
    if (Provider.isProvider(this.l2SignerOrProvider)) {
      throw new Error(`messenger has no L2 signer`)
    } else {
      return this.l2SignerOrProvider
    }
  }

  public async getL1BlockNumber(): Promise<number> {
    return this.contracts.l2.OVM_L1BlockNumber.getL1BlockNumber()
  }

  public async waitingDepositTransactionRelayed(
    txReceipt: TransactionReceipt,
    opts: {
      pollIntervalMs?: number
      timeoutMs?: number
    }
  ): Promise<string> {
    const l1BlockNumber = txReceipt.blockNumber
    let totalTimeMs = 0
    while (totalTimeMs < (opts.timeoutMs || Infinity)) {
      const tick = Date.now()
      if (l1BlockNumber <= (await this.getL1BlockNumber())) {
        return this.calculateRelayedDepositTxID(txReceipt)
      }
      await sleep(opts.pollIntervalMs || 1000)
      totalTimeMs += Date.now() - tick
    }
    throw new Error(`timed out waiting for relayed deposit transaction`)
  }

  public async calculateRelayedDepositTxID(
    txReceipt: TransactionReceipt
  ): Promise<string> {
    if (txReceipt.status !== 1) {
      return null
    }
    let promiseString: Promise<string> = null
    txReceipt.logs.forEach((log) => {
      if (
        log.topics[0] ===
        ethers.utils.id('TransactionDeposited(address,address,uint256,bytes)')
      ) {
        const depositTx = DepositTx.fromL1Log(log)
        promiseString = Promise.resolve(depositTx.hash())
      }
    })
    return promiseString
  }

  /**
   * Queries the current challenge period in seconds from the StateCommitmentChain.
   *
   * @returns Current challenge period in seconds.
   */
  public async getChallengePeriodSeconds(): Promise<number> {
    const oracleVersion = await this.contracts.l1.L2OutputOracle.version()
    const challengePeriod =
      oracleVersion === '1.0.0'
        ? // The ABI in the SDK does not contain FINALIZATION_PERIOD_SECONDS
          // in OptimismPortal, so making an explicit call instead.
          BigNumber.from(
            await this.contracts.l1.OptimismPortal.provider.call({
              to: this.contracts.l1.OptimismPortal.address,
              data: '0xf4daa291', // FINALIZATION_PERIOD_SECONDS
            })
          )
        : await this.contracts.l1.L2OutputOracle.FINALIZATION_PERIOD_SECONDS()
    return challengePeriod.toNumber()
  }

  /**
   * Queries the OptimismPortal contract's `provenWithdrawals` mapping
   * for a ProvenWithdrawal that matches the passed withdrawalHash
   *
   * @bedrock
   * Note: This function is bedrock-specific.
   *
   * @returns A ProvenWithdrawal object
   */
  public async getProvenWithdrawal(
    withdrawalHash: string
  ): Promise<ProvenWithdrawal> {
    return this.contracts.l1.OptimismPortal.provenWithdrawals(withdrawalHash)
  }

  public async depositTransaction(
    request: DepositTransactionRequest,
    opts?: {
      signer?: Signer
    }
  ): Promise<TransactionResponse> {
    const tx = await this.populateTransaction.depositTransaction(request)
    return (opts?.signer || this.l1Signer).sendTransaction(tx)
  }

  populateTransaction = {
    depositTransaction: async (
      request: DepositTransactionRequest
    ): Promise<TransactionRequest> => {
      return this.contracts.l1.OptimismPortal.populateTransaction.depositTransaction(
        request.to,
        request.value,
        request.gasLimit,
        request.data
      )
    },
  }
}
