/* eslint-disable @typescript-eslint/no-unused-vars */
import {
  Provider,
  TransactionResponse,
  TransactionRequest,
  TransactionReceipt,
} from '@ethersproject/abstract-provider'
import { Signer } from '@ethersproject/abstract-signer'
import { BigNumber, PayableOverrides, ethers } from 'ethers'
import {
  sleep,
  DepositTx,
  toRpcHexString,
  predeploys,
  BedrockOutputData,
  BedrockCrossChainMessageProof,
} from '@tokamak-network/core-utils'
import semver from 'semver'

import {
  DepositTransactionRequest,
  OEContracts,
  OEContractsLike,
  NumberLike,
  SignerOrProviderLike,
  ProvenWithdrawal,
  WithdrawalMessageInfo,
  WithdrawalTransactionRequest,
  MessageStatus,
  StateRoot,
  StateRootBatch,
} from './interfaces'
import {
  hashLowLevelMessage,
  hashMessageHash,
  makeStateTrieProof,
  toSignerOrProvider,
  toNumber,
  DeepPartial,
  getPortalsContracts,
  DEPOSIT_CONFIRMATION_BLOCKS,
  CHAIN_BLOCK_TIMES,
  calculateWithdrawalMessageUsingRecept,
  getContractInterfaceBedrock,
  toJsonRpcProvider,
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
   * Whether or not Bedrock compatibility is enabled.
   */
  public bedrock: boolean

  /**
   * Cache for output root validation. Output roots are expensive to verify, so we cache them.
   */
  private _outputCache: Array<{ root: string; valid: boolean }> = []

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
   * @param opts.bedrock Whether or not to enable Bedrock compatibility.
   */
  constructor(opts: {
    l1SignerOrProvider: SignerOrProviderLike
    l2SignerOrProvider: SignerOrProviderLike
    l1ChainId: NumberLike
    l2ChainId: NumberLike
    depositConfirmationBlocks?: NumberLike
    l1BlockTimeSeconds?: NumberLike
    contracts?: DeepPartial<OEContractsLike>
    bedrock?: boolean
  }) {
    this.bedrock = opts.bedrock ?? true
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

  public async getMessageStatus(
    txReceipt: TransactionReceipt
  ): Promise<MessageStatus> {
    return txReceipt.to === predeploys.L2ToL1MessagePasser
      ? this.getL2ToL1MessageStatusByReceipt(txReceipt)
      : this.getL1ToL2MessageStatusByReceipt(txReceipt)
  }

  public async getL1ToL2MessageStatusByReceipt(
    txReceipt: TransactionReceipt
  ): Promise<MessageStatus> {
    const l1BlockNumber = await this.getL1BlockNumber()
    return txReceipt.blockNumber > l1BlockNumber
      ? MessageStatus.UNCONFIRMED_L1_TO_L2_MESSAGE
      : MessageStatus.RELAYED
  }

  /**
   * Returns the StateBatchAppended event that was emitted when the batch with a given index was
   * created. Returns null if no such event exists (the batch has not been submitted).
   *
   * @param batchIndex Index of the batch to find an event for.
   * @returns StateBatchAppended event for the batch, or null if no such batch exists.
   */
  public async getStateBatchAppendedEventByBatchIndex(
    batchIndex: number
  ): Promise<ethers.Event | null> {
    const events = await this.contracts.l1.StateCommitmentChain.queryFilter(
      this.contracts.l1.StateCommitmentChain.filters.StateBatchAppended(
        batchIndex
      )
    )

    if (events.length === 0) {
      return null
    } else if (events.length > 1) {
      // Should never happen!
      throw new Error(`found more than one StateBatchAppended event`)
    } else {
      return events[0]
    }
  }

  /**
   * Returns the StateBatchAppended event for the batch that includes the transaction with the
   * given index. Returns null if no such event exists.
   *
   * @param transactionIndex Index of the L2 transaction to find an event for.
   * @returns StateBatchAppended event for the batch that includes the given transaction by index.
   */
  public async getStateBatchAppendedEventByTransactionIndex(
    transactionIndex: number
  ): Promise<ethers.Event | null> {
    const isEventHi = (event: ethers.Event, index: number) => {
      const prevTotalElements = event.args._prevTotalElements.toNumber()
      return index < prevTotalElements
    }

    const isEventLo = (event: ethers.Event, index: number) => {
      const prevTotalElements = event.args._prevTotalElements.toNumber()
      const batchSize = event.args._batchSize.toNumber()
      return index >= prevTotalElements + batchSize
    }

    const totalBatches: ethers.BigNumber =
      await this.contracts.l1.StateCommitmentChain.getTotalBatches()
    if (totalBatches.eq(0)) {
      return null
    }

    let lowerBound = 0
    let upperBound = totalBatches.toNumber() - 1
    let batchEvent: ethers.Event | null =
      await this.getStateBatchAppendedEventByBatchIndex(upperBound)

    // Only happens when no batches have been submitted yet.
    if (batchEvent === null) {
      return null
    }

    if (isEventLo(batchEvent, transactionIndex)) {
      // Upper bound is too low, means this transaction doesn't have a corresponding state batch yet.
      return null
    } else if (!isEventHi(batchEvent, transactionIndex)) {
      // Upper bound is not too low and also not too high. This means the upper bound event is the
      // one we're looking for! Return it.
      return batchEvent
    }

    // Binary search to find the right event. The above checks will guarantee that the event does
    // exist and that we'll find it during this search.
    while (lowerBound < upperBound) {
      const middleOfBounds = Math.floor((lowerBound + upperBound) / 2)
      batchEvent = await this.getStateBatchAppendedEventByBatchIndex(
        middleOfBounds
      )

      if (isEventHi(batchEvent, transactionIndex)) {
        upperBound = middleOfBounds
      } else if (isEventLo(batchEvent, transactionIndex)) {
        lowerBound = middleOfBounds
      } else {
        break
      }
    }

    return batchEvent
  }

  /**
   * Returns information about the state root batch that included the state root for the given
   * transaction by index. Returns null if no such state root has been published yet.
   *
   * @param transactionIndex Index of the L2 transaction to find a state root batch for.
   * @returns State root batch for the given transaction index, or null if none exists yet.
   */
  public async getStateRootBatchByTransactionIndex(
    transactionIndex: number
  ): Promise<StateRootBatch | null> {
    const stateBatchAppendedEvent =
      await this.getStateBatchAppendedEventByTransactionIndex(transactionIndex)
    if (stateBatchAppendedEvent === null) {
      return null
    }

    const stateBatchTransaction = await stateBatchAppendedEvent.getTransaction()
    const [stateRoots] =
      this.contracts.l1.StateCommitmentChain.interface.decodeFunctionData(
        'appendStateBatch',
        stateBatchTransaction.data
      )

    return {
      blockNumber: stateBatchAppendedEvent.blockNumber,
      stateRoots,
      header: {
        batchIndex: stateBatchAppendedEvent.args._batchIndex,
        batchRoot: stateBatchAppendedEvent.args._batchRoot,
        batchSize: stateBatchAppendedEvent.args._batchSize,
        prevTotalElements: stateBatchAppendedEvent.args._prevTotalElements,
        extraData: stateBatchAppendedEvent.args._extraData,
      },
    }
  }

  /**
   * Returns the state root that corresponds to a given message. This is the state root for the
   * block in which the transaction was included, as published to the StateCommitmentChain. If the
   * state root for the given message has not been published yet, this function returns null.
   *
   * @param message Message to find a state root for.
   * @param messageIndex The index of the message, if multiple exist from multicall
   * @returns State root for the block in which the message was created.
   */
  public async getMessageStateRoot(
    transactionHash: string
  ): Promise<StateRoot | null> {
    // We need the block number of the transaction that triggered the message so we can look up the
    // state root batch that corresponds to that block number.
    const messageTxReceipt = await this.l2Provider.getTransactionReceipt(
      transactionHash
    )

    // Every block has exactly one transaction in it. Since there's a genesis block, the
    // transaction index will always be one less than the block number.
    const messageTxIndex = messageTxReceipt.blockNumber - 1

    // Pull down the state root batch, we'll try to pick out the specific state root that
    // corresponds to our message.
    const stateRootBatch = await this.getStateRootBatchByTransactionIndex(
      messageTxIndex
    )

    // No state root batch, no state root.
    if (stateRootBatch === null) {
      return null
    }

    // We have a state root batch, now we need to find the specific state root for our transaction.
    // First we need to figure out the index of the state root within the batch we found. This is
    // going to be the original transaction index offset by the total number of previous state
    // roots.
    const indexInBatch =
      messageTxIndex - stateRootBatch.header.prevTotalElements.toNumber()

    // Just a sanity check.
    if (stateRootBatch.stateRoots.length <= indexInBatch) {
      // Should never happen!
      throw new Error(`state root does not exist in batch`)
    }

    return {
      stateRoot: stateRootBatch.stateRoots[indexInBatch],
      stateRootIndexInBatch: indexInBatch,
      batch: stateRootBatch,
    }
  }

  /**
   * Returns the Bedrock output root that corresponds to the given message.
   *
   * @param message Message to get the Bedrock output root for.
   * @param messageIndex The index of the message, if multiple exist from multicall
   * @returns Bedrock output root.
   */
  public async getMessageBedrockOutput(
    l2BlockNumber: number
  ): Promise<BedrockOutputData | null> {
    let proposal: any
    let l2OutputIndex: BigNumber
    if (await this.fpac()) {
      // Get the respected game type from the portal.
      const gameType =
        await this.contracts.l1.OptimismPortal2.respectedGameType()

      // Get the total game count from the DisputeGameFactory since that will give us the end of
      // the array that we're searching over. We'll then use that to find the latest games.
      const gameCount = await this.contracts.l1.DisputeGameFactory.gameCount()

      // Find the latest 100 games (or as many as we can up to 100).
      const latestGames =
        await this.contracts.l1.DisputeGameFactory.findLatestGames(
          gameType,
          Math.max(0, gameCount.sub(1).toNumber()),
          Math.min(100, gameCount.toNumber())
        )

      // Find all games that are for proposals about blocks newer than the message block.
      const matches: any[] = []
      for (const game of latestGames) {
        try {
          const [blockNumber] = ethers.utils.defaultAbiCoder.decode(
            ['uint256'],
            game.extraData
          )
          if (blockNumber.gte(l2BlockNumber)) {
            matches.push({
              ...game,
              l2BlockNumber: blockNumber,
            })
          }
        } catch (err) {
          // If we can't decode the extra data then we just skip this game.
          continue
        }
      }

      // Shuffle the list of matches. We shuffle here to avoid potential DoS vectors where the
      // latest games are all invalid and the SDK would be forced to make a bunch of archive calls.
      for (let i = matches.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1))
        ;[matches[i], matches[j]] = [matches[j], matches[i]]
      }

      // Now we verify the proposals in the matches array.
      let match: any
      for (const option of matches) {
        if (
          await this.isValidOutputRoot(option.rootClaim, option.l2BlockNumber)
        ) {
          match = option
          break
        }
      }

      // If there's no match then we can't prove the message to the portal.
      if (!match) {
        return null
      }

      // Put the result into the same format as the old logic for now to reduce added code.
      l2OutputIndex = match.index
      proposal = {
        outputRoot: match.rootClaim,
        timestamp: match.timestamp,
        l2BlockNumber: match.l2BlockNumber,
      }
    } else {
      // Try to find the output index that corresponds to the block number attached to the message.
      // We'll explicitly handle "cannot get output" errors as a null return value, but anything else
      // needs to get thrown. Might need to revisit this in the future to be a little more robust
      // when connected to RPCs that don't return nice error messages.
      try {
        l2OutputIndex =
          await this.contracts.l1.L2OutputOracle.getL2OutputIndexAfter(
            l2BlockNumber
          )
      } catch (err) {
        if (err.message.includes('L2OutputOracle: cannot get output')) {
          return null
        } else {
          throw err
        }
      }

      // Now pull the proposal out given the output index. Should always work as long as the above
      // codepath completed successfully.
      proposal = await this.contracts.l1.L2OutputOracle.getL2Output(
        l2OutputIndex
      )
    }

    // Format everything and return it nicely.
    return {
      outputRoot: proposal.outputRoot,
      l1Timestamp: proposal.timestamp.toNumber(),
      l2BlockNumber: proposal.l2BlockNumber.toNumber(),
      l2OutputIndex: l2OutputIndex.toNumber(),
    }
  }

  public async getL2ToL1MessageStatusByReceipt(
    txReceipt: TransactionReceipt
  ): Promise<MessageStatus> {
    const withdrawalMessageInfo = await this.calculateWithdrawalMessage(
      txReceipt
    )

    if (!withdrawalMessageInfo) {
      throw Error('withdrawal message not found')
    }

    const finalizedMessage = await this.getFinalizedWithdrawalStatus(
      withdrawalMessageInfo.withdrawalHash
    )
    if (finalizedMessage) {
      return MessageStatus.RELAYED
    }

    let timestamp: number
    if (this.bedrock) {
      const output = await this.getMessageBedrockOutput(
        withdrawalMessageInfo.l2BlockNumber
      )
      if (output === null) {
        return MessageStatus.STATE_ROOT_NOT_PUBLISHED
      }

      const provenWithdrawal = await this.getProvenWithdrawal(
        withdrawalMessageInfo.withdrawalHash
      )
      if (!provenWithdrawal || provenWithdrawal.timestamp.toNumber() === 0) {
        return MessageStatus.READY_TO_PROVE
      }

      timestamp = provenWithdrawal.timestamp.toNumber()
    } else {
      const stateRoot = await this.getMessageStateRoot(
        txReceipt.transactionHash
      )
      if (stateRoot === null) {
        return MessageStatus.STATE_ROOT_NOT_PUBLISHED
      }

      const bn = stateRoot.batch.blockNumber
      const block = await this.l1Provider.getBlock(bn)

      timestamp = block.timestamp
    }

    if (await this.fpac()) {
      // Grab the proven withdrawal data.
      const provenWithdrawal = await this.getProvenWithdrawal(
        withdrawalMessageInfo.withdrawalHash
      )

      // Sanity check, should've already happened above but do it just in case.
      if (provenWithdrawal === null) {
        // Ready to prove is the correct status here, we would not expect to hit this code path
        // unless there was an unexpected reorg on L1. Since this is unlikely we log a warning.
        console.warn(
          'Unexpected code path reached in getMessageStatus, returning READY_TO_PROVE'
        )
        return MessageStatus.READY_TO_PROVE
      }

      // Shouldn't happen, but worth checking just in case.
      if (!('proofSubmitter' in provenWithdrawal)) {
        throw new Error(
          `expected to get FPAC withdrawal but got legacy withdrawal`
        )
      }

      try {
        // If this doesn't revert then we should be fine to relay.
        await this.contracts.l1.OptimismPortal2.checkWithdrawal(
          hashLowLevelMessage(withdrawalMessageInfo),
          provenWithdrawal.proofSubmitter
        )

        return MessageStatus.READY_FOR_RELAY
      } catch (err) {
        return MessageStatus.IN_CHALLENGE_PERIOD
      }
    } else {
      const challengePeriod = await this.getChallengePeriodSeconds()
      const latestBlock = await this.l1Provider.getBlock('latest')

      if (timestamp + challengePeriod > latestBlock.timestamp) {
        return MessageStatus.IN_CHALLENGE_PERIOD
      } else {
        return MessageStatus.READY_FOR_RELAY
      }
    }
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

  public async calculateWithdrawalMessage(
    txReceipt: TransactionReceipt
  ): Promise<WithdrawalMessageInfo> {
    if (txReceipt.status !== 1) {
      return null
    }
    let promiseMessage: Promise<WithdrawalMessageInfo> = null
    const withdrawalMessage = calculateWithdrawalMessageUsingRecept(txReceipt)
    promiseMessage = Promise.resolve(withdrawalMessage)
    return promiseMessage
  }

  public async calculateWithdrawalMessageByL2TxHash(
    transactionHash: string
  ): Promise<WithdrawalMessageInfo> {
    const txReceipt = await this.l2Provider.getTransactionReceipt(
      transactionHash
    )
    return this.calculateWithdrawalMessage(txReceipt)
  }

  public async waitForMessageStatus(
    txReceipt: TransactionReceipt,
    status: MessageStatus,
    opts?: {
      pollIntervalMs?: number
      timeoutMs?: number
    }
  ): Promise<void> {
    let totalTimeMs = 0
    while (totalTimeMs < (opts?.timeoutMs || Infinity)) {
      const tick = Date.now()

      const currentStatus = await this.getMessageStatus(txReceipt)
      if (currentStatus >= status) {
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
    if (!this.bedrock) {
      return (
        await this.contracts.l1.StateCommitmentChain.FRAUD_PROOF_WINDOW()
      ).toNumber()
    }

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
   * Generates the bedrock proof required to finalize an L2 to L1 message.
   *
   * @param message Message to generate a proof for.
   * @param messageIndex The index of the message, if multiple exist from multicall
   * @returns Proof that can be used to finalize the message.
   */
  public async getBedrockMessageProof(
    withdrawalMessageInfo: WithdrawalMessageInfo
  ): Promise<BedrockCrossChainMessageProof> {
    const output = await this.getMessageBedrockOutput(
      withdrawalMessageInfo.l2BlockNumber
    )
    if (output === null) {
      throw new Error(`state root for message not yet published`)
    }

    const hash = hashLowLevelMessage(withdrawalMessageInfo)
    const messageSlot = hashMessageHash(hash)

    const provider = toJsonRpcProvider(this.l2Provider)

    const stateTrieProof = await makeStateTrieProof(
      provider,
      output.l2BlockNumber,
      this.contracts.l2.BedrockMessagePasser.address,
      messageSlot
    )

    const block = await provider.send('eth_getBlockByNumber', [
      toRpcHexString(output.l2BlockNumber),
      false,
    ])

    return {
      outputRootProof: {
        version: ethers.constants.HashZero,
        stateRoot: block.stateRoot,
        messagePasserStorageRoot: stateTrieProof.storageRoot,
        latestBlockhash: block.hash,
      },
      withdrawalProof: stateTrieProof.storageProof,
      l2OutputIndex: output.l2OutputIndex,
    }
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
    if (!this.bedrock) {
      throw new Error('message proving only applies after the bedrock upgrade')
    }

    if (await this.fpac()) {
      // Getting the withdrawal is a bit more complicated after FPAC.
      // First we need to get the number of proof submitters for this withdrawal.
      const numProofSubmitters = BigNumber.from(
        await this.contracts.l1.OptimismPortal2.numProofSubmitters(
          withdrawalHash
        )
      ).toNumber()

      // Now we need to find any withdrawal where the output proposal that the withdrawal was proven
      // against is actually valid. We can use the same output validation cache used elsewhere.
      for (let i = 0; i < numProofSubmitters; i++) {
        // Grab the proof submitter.
        const proofSubmitter =
          await this.contracts.l1.OptimismPortal2.proofSubmitters(
            withdrawalHash,
            i
          )

        // Grab the ProvenWithdrawal struct for this proof.
        const provenWithdrawal =
          await this.contracts.l1.OptimismPortal2.provenWithdrawals(
            withdrawalHash,
            proofSubmitter
          )

        // Grab the game that was proven against.
        const game = new ethers.Contract(
          provenWithdrawal.disputeGameProxy,
          getContractInterfaceBedrock('FaultDisputeGame'),
          this.l1SignerOrProvider
        )

        // Check the game status.
        const status = await game.status()
        if (status === 1) {
          // If status is CHALLENGER_WINS then it's no good.
          continue
        } else if (status === 2) {
          // If status is DEFENDER_WINS then it's a valid proof.
          return {
            ...provenWithdrawal,
            proofSubmitter,
          }
        } else if (status > 2) {
          // Shouldn't happen in practice.
          throw new Error('got invalid game status')
        }

        // Otherwise we're IN_PROGRESS.
        // Grab the block number from the extra data. Since this is not a standardized field we need
        // to be defensive and assume that the extra data could be anything. If the extra data does
        // not decode properly then we just skip this game.
        const extraData = await game.extraData()
        let l2BlockNumber: number
        try {
          ;[l2BlockNumber] = ethers.utils.defaultAbiCoder.decode(
            ['uint256'],
            extraData
          )
        } catch (err) {
          // Didn't decode properly, bad game.
          continue
        }

        // Finally we check if the output root is valid. If it is, then we can return the proven
        // withdrawal. If it isn't, then we act as if this proof does not exist because it isn't
        // useful for finalizing the withdrawal.
        if (
          await this.isValidOutputRoot(await game.rootClaim(), l2BlockNumber)
        ) {
          return {
            ...provenWithdrawal,
            proofSubmitter,
          }
        }
      }

      // Return null if we didn't find a valid proof.
      return null
    } else {
      return this.contracts.l1.OptimismPortal.provenWithdrawals(withdrawalHash)
    }
  }

  public async isValidOutputRoot(
    outputRoot: string,
    l2BlockNumber: number
  ): Promise<boolean> {
    // Use the cache if we can.
    const cached = this._outputCache.find((other) => {
      return other.root === outputRoot
    })

    // Skip if we can use the cached.
    if (cached) {
      return cached.valid
    }

    // If the cache ever gets to 10k elements, clear out the first half. Works well enough
    // since the cache will generally tend to be used in a FIFO manner.
    if (this._outputCache.length > 10000) {
      this._outputCache = this._outputCache.slice(5000)
    }

    // We didn't hit the cache so we're going to have to do the work.
    try {
      // Make sure this is a JSON RPC provider.
      const provider = toJsonRpcProvider(this.l2Provider)

      // Grab the block and storage proof at the same time.
      const [block, proof] = await Promise.all([
        provider.send('eth_getBlockByNumber', [
          toRpcHexString(l2BlockNumber),
          false,
        ]),
        makeStateTrieProof(
          provider,
          l2BlockNumber,
          this.contracts.l2.OVM_L2ToL1MessagePasser.address,
          ethers.constants.HashZero
        ),
      ])

      // Compute the output.
      const output = ethers.utils.solidityKeccak256(
        ['bytes32', 'bytes32', 'bytes32', 'bytes32'],
        [
          ethers.constants.HashZero,
          block.stateRoot,
          proof.storageRoot,
          block.hash,
        ]
      )

      // If the output matches the proposal then we're good.
      const valid = output === outputRoot
      this._outputCache.push({ root: outputRoot, valid })
      return valid
    } catch (err) {
      // Assume the game is invalid but don't add it to the cache just in case we had a temp error.
      return false
    }
  }

  public async getFinalizedWithdrawalStatus(
    withdrawalHash: string
  ): Promise<boolean> {
    return this.contracts.l1.OptimismPortal.finalizedWithdrawals(withdrawalHash)
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

  public async proveWithdrawalTransactionUsingL2Tx(
    transactionHash: string,
    opts?: {
      signer?: Signer
    }
  ): Promise<TransactionResponse> {
    const message = await this.calculateWithdrawalMessageByL2TxHash(
      transactionHash
    )
    return this.proveWithdrawalTransaction(message, opts)
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

  public async finalizeWithdrawalTransactionUsingL2Tx(
    transactionHash: string,
    opts?: {
      signer?: Signer
    }
  ): Promise<TransactionResponse> {
    const message = await this.calculateWithdrawalMessageByL2TxHash(
      transactionHash
    )
    return this.finalizeWithdrawalTransaction(message)
  }

  /**
   * Uses portal version to determine if the messenger is using fpac contracts. Better not to cache
   * this value as it will change during the fpac upgrade and we want clients to automatically
   * begin using the new logic without throwing any errors.
   *
   * @returns Whether or not the messenger is using fpac contracts.
   */
  public async fpac(): Promise<boolean> {
    if (
      this.contracts.l1.OptimismPortal.address === ethers.constants.AddressZero
    ) {
      // Only really relevant for certain SDK tests where the portal is not deployed. We should
      // probably just update the tests so the portal gets deployed but feels like it's out of
      // scope for the FPAC changes.
      return false
    } else {
      return semver.gte(
        await this.contracts.l1.OptimismPortal.version(),
        '3.0.0'
      )
    }
  }

  populateTransaction = {
    depositTransaction: async (
      request: DepositTransactionRequest
    ): Promise<TransactionRequest> => {
      return this.contracts.l1.OptimismPortal.populateTransaction.depositTransaction(
        request.to,
        request.mint,
        request.value,
        request.gasLimit,
        request.isCreation,
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
      const proof = await this.getBedrockMessageProof(message)

      const args = [
        [
          message.messageNonce,
          message.sender,
          message.target,
          message.value,
          message.minGasLimit,
          message.message,
        ],
        proof.l2OutputIndex,
        [
          proof.outputRootProof.version,
          proof.outputRootProof.stateRoot,
          proof.outputRootProof.messagePasserStorageRoot,
          proof.outputRootProof.latestBlockhash,
        ],
        proof.withdrawalProof,
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

  /**
   * Object that holds the functions that estimates the gas required for a given transaction.
   * Follows the pattern used by ethers.js.
   */
  estimateGas = {
    depositTransaction: async (
      request: DepositTransactionRequest
    ): Promise<BigNumber> => {
      const tx = await this.populateTransaction.depositTransaction(request)
      return this.l1Provider.estimateGas(tx)
    },

    initiateWithdrawal: async (
      request: WithdrawalTransactionRequest
    ): Promise<BigNumber> => {
      const tx = await this.populateTransaction.initiateWithdrawal(request)
      return this.l2Provider.estimateGas(tx)
    },

    proveWithdrawalTransaction: async (
      message: WithdrawalMessageInfo,
      opts?: {
        overrides?: PayableOverrides
      }
    ): Promise<BigNumber> => {
      const tx = await this.populateTransaction.proveWithdrawalTransaction(
        message,
        opts
      )
      return this.l1Provider.estimateGas(tx)
    },

    finalizeWithdrawalTransaction: async (
      message: WithdrawalMessageInfo
    ): Promise<BigNumber> => {
      const tx = await this.populateTransaction.finalizeWithdrawalTransaction(
        message
      )
      return this.l1Provider.estimateGas(tx)
    },
  }
}
