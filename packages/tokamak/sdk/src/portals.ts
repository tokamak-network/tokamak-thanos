/* eslint-disable @typescript-eslint/no-unused-vars */
import {
  Provider,
  TransactionResponse,
  TransactionRequest,
  TransactionReceipt,
} from '@ethersproject/abstract-provider'
import { Signer } from '@ethersproject/abstract-signer'
import { BigNumber, PayableOverrides, ethers } from 'ethers'
import { sleep, DepositTx, toRpcHexString } from '@tokamak-network/core-utils'

import {
  DepositTransactionRequest,
  OEContracts,
  OEContractsLike,
  NumberLike,
  SignerOrProviderLike,
  ProvenWithdrawal,
  WithdrawalMessageInfo,
  WithdrawalTransactionRequest,
} from './interfaces'
import {
  calculateWithdrawalMessage,
  hashLowLevelMessage,
  hashMessageHash,
  makeStateTrieProof,
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

  public async getL2BlockNumber(): Promise<number> {
    return this.contracts.l1.L2OutputOracle.latestBlockNumber()
  }

  public async calculateWithdrawalMessage(
    txReceipt: TransactionReceipt
  ): Promise<WithdrawalMessageInfo> {
    if (txReceipt.status !== 1) {
      return null
    }
    let promiseMessage: Promise<WithdrawalMessageInfo> = null
    txReceipt.logs.forEach((log) => {
      if (
        log.topics[0] ===
        ethers.utils.id(
          'MessagePassed(uint256,address,address,uint256,uint256,bytes,bytes32)'
        )
      ) {
        const withdrawalMessage = calculateWithdrawalMessage(log)
        promiseMessage = Promise.resolve(withdrawalMessage)
      }
    })
    return promiseMessage
  }

  public async waitForWithdrawalTxReadyForRelay(
    txReceipt: TransactionReceipt,
    opts?: {
      pollIntervalMs?: number
      timeoutMs?: number
    }
  ) {
    const l2BlockNumber = txReceipt.blockNumber
    let totalTimeMs = 0
    while (totalTimeMs < (opts?.timeoutMs || Infinity)) {
      const tick = Date.now()
      if (l2BlockNumber <= (await this.getL2BlockNumber())) {
        return
      }
      await sleep(opts?.pollIntervalMs || 1000)
      totalTimeMs += Date.now() - tick
    }
    throw new Error(`timed out waiting for relayed deposit transaction`)
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

  public async initiateWithdrawal(
    request: WithdrawalTransactionRequest,
    opts?: {
      signer?: Signer
    }
  ): Promise<TransactionResponse> {
    const tx = await this.populateTransaction.initiateWithdrawal(request)
    return (opts?.signer || this.l2Signer).sendTransaction(tx)
  }

  public async proveWithdrawalTransaction(
    message: WithdrawalMessageInfo,
    opts?: {
      signer?: Signer
    }
  ): Promise<TransactionResponse> {
    const tx = await this.populateTransaction.proveWithdrawalTransaction(
      message
    )
    return (opts?.signer || this.l1Signer).sendTransaction(tx)
  }


  public async waitForFinalization(
    message: WithdrawalMessageInfo,
    opts?: {
      pollIntervalMs?: number
      timeoutMs?: number
    }
  ) {
    const provenWithdrawal = await this.getProvenWithdrawal(message.withdrawalHash)
    const finalizedPeriod = await this.getChallengePeriodSeconds();
    const BUFFER_TIME = 12
    let totalTimeMs = 0
    while (totalTimeMs < (opts?.timeoutMs || Infinity)) {
      const currentTimestamp = Date.now()
      if (currentTimestamp / 1000 - BUFFER_TIME > provenWithdrawal.timestamp.toNumber() + finalizedPeriod) {
        return
      }
      await sleep(opts?.pollIntervalMs || 1000)
      totalTimeMs += Date.now() - currentTimestamp
    }
    throw new Error(`timed out waiting for relayed deposit transaction`)
  }

  public async finalizeWithdrawalTransaction(
    message: WithdrawalMessageInfo,
    opts?: {
      signer?: Signer
    }
  ): Promise<TransactionResponse> {
    const tx = await this.populateTransaction.finalizeWithdrawalTransaction(
      message
    )
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

    initiateWithdrawal: async (
      request: WithdrawalTransactionRequest
    ): Promise<TransactionRequest> => {
      return this.contracts.l2.BedrockMessagePasser.populateTransaction.initiateWithdrawal(
        request.target,
        request.gasLimit,
        request.data,
        {
          value: request.value,
        }
      )
    },

    proveWithdrawalTransaction: async (
      message: WithdrawalMessageInfo,
      opts?: {
        overrides?: PayableOverrides
      }
    ): Promise<TransactionRequest> => {
      const l2OutputIndex =
        await this.contracts.l1.L2OutputOracle.getL2OutputIndexAfter(
          message.l2BlockNumber
        )
      const proposal = await this.contracts.l1.L2OutputOracle.getL2Output(
        l2OutputIndex
      )

      const withdrawalInfoSlot = hashMessageHash(
        hashLowLevelMessage({
          message: message.message,
          messageNonce: message.messageNonce,
          minGasLimit: message.minGasLimit,
          value: message.value,
          sender: message.sender,
          target: message.target,
        })
      )

      const storageProof = await makeStateTrieProof(
        this.l2Provider as ethers.providers.JsonRpcProvider,
        proposal.l2BlockNumber,
        this.contracts.l2.BedrockMessagePasser.address,
        withdrawalInfoSlot
      )

      const block = await (
        this.l2Provider as ethers.providers.JsonRpcProvider
      ).send('eth_getBlockByNumber', [
        toRpcHexString(proposal.l2BlockNumber),
        false,
      ])

      const args = [
        [
          message.messageNonce,
          message.sender,
          message.target,
          message.value,
          message.minGasLimit,
          message.message,
        ],
        l2OutputIndex.toNumber(),
        [
          ethers.constants.HashZero,
          block.stateRoot,
          storageProof.storageRoot, // for proving storage root is correct
          block.hash,
        ],
        storageProof.storageProof, // for proving withdrawal info
        opts?.overrides || {},
      ] as const
      return this.contracts.l1.OptimismPortal.populateTransaction.proveWithdrawalTransaction(
        ...args
      )
    },

    finalizeWithdrawalTransaction: async (
      message: WithdrawalMessageInfo
    ): Promise<TransactionRequest> => {
      return this.contracts.l1.OptimismPortal.populateTransaction.finalizeWithdrawalTransaction(
        [
          message.messageNonce,
          message.sender,
          message.target,
          message.value,
          message.minGasLimit,
          message.message,
        ]
      )
    },
  }
}
