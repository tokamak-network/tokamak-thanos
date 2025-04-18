import { add0x, hashWithdrawal } from '@tokamak-network/core-utils'
import { BigNumber, utils, ethers } from 'ethers'
import { Log, TransactionReceipt } from '@ethersproject/abstract-provider'
import { hexDataSlice } from 'ethers/lib/utils'

import {
  L1ChainID,
  LowLevelMessage,
  WithdrawalMessageInfo,
} from '../interfaces'

const { hexDataLength } = utils

// Constants used by `CrossDomainMessenger.baseGas`
const RELAY_CONSTANT_OVERHEAD = BigNumber.from(200_000)
const RELAY_PER_BYTE_DATA_COST = BigNumber.from(16)
const MIN_GAS_DYNAMIC_OVERHEAD_NUMERATOR = BigNumber.from(64)
const MIN_GAS_DYNAMIC_OVERHEAD_DENOMINATOR = BigNumber.from(63)
const RELAY_CALL_OVERHEAD = BigNumber.from(40_000)
const RELAY_RESERVED_GAS = BigNumber.from(40_000)
const RELAY_GAS_CHECK_BUFFER = BigNumber.from(5_000)
const RELAY_GAS_CHECK_BUFFER_INCLUDING_APPROVAL = BigNumber.from(40_000)

/**
 * Utility for hashing a LowLevelMessage object.
 *
 * @param message LowLevelMessage object to hash.
 * @returns Hash of the given LowLevelMessage.
 */
export const hashLowLevelMessage = (message: LowLevelMessage): string => {
  return hashWithdrawal(
    message.messageNonce,
    message.sender,
    message.target,
    message.value,
    message.minGasLimit,
    message.message
  )
}

/**
 * Utility for hashing a message hash. This computes the storage slot
 * where the message hash will be stored in state. HashZero is used
 * because the first mapping in the contract is used.
 *
 * @param messageHash Message hash to hash.
 * @returns Hash of the given message hash.
 */
export const hashMessageHash = (messageHash: string): string => {
  const data = ethers.utils.defaultAbiCoder.encode(
    ['bytes32', 'uint256'],
    [messageHash, ethers.constants.HashZero]
  )
  return ethers.utils.keccak256(data)
}

/**
 * Compute the min gas limit for a migrated withdrawal.
 */
export const migratedWithdrawalGasLimit = (
  data: string,
  chainID: number
): BigNumber => {
  // Compute the gas limit and cap at 25 million
  const dataCost = BigNumber.from(hexDataLength(data)).mul(
    RELAY_PER_BYTE_DATA_COST
  )
  let overhead: BigNumber
  if (chainID === 420) {
    overhead = BigNumber.from(200_000)
  } else {
    const relayGasBuffer = Object.values(L1ChainID).includes(chainID)
      ? RELAY_GAS_CHECK_BUFFER
      : RELAY_GAS_CHECK_BUFFER_INCLUDING_APPROVAL

    // Dynamic overhead (EIP-150)
    // We use a constant 1 million gas limit due to the overhead of simulating all migrated withdrawal
    // transactions during the migration. This is a conservative estimate, and if a withdrawal
    // uses more than the minimum gas limit, it will fail and need to be replayed with a higher
    // gas limit.
    const dynamicOverhead = MIN_GAS_DYNAMIC_OVERHEAD_NUMERATOR.mul(
      1_000_000
    ).div(MIN_GAS_DYNAMIC_OVERHEAD_DENOMINATOR)

    // Constant overhead
    overhead = RELAY_CONSTANT_OVERHEAD.add(dynamicOverhead)
      .add(RELAY_CALL_OVERHEAD)
      // Gas reserved for the worst-case cost of 3/5 of the `CALL` opcode's dynamic gas
      // factors. (Conservative)
      // Relay reserved gas (to ensure execution of `relayMessage` completes after the
      // subcontext finishes executing) (Conservative)
      .add(RELAY_RESERVED_GAS)
      // Gas reserved for the execution between the `hasMinGas` check and the `CALL`
      // opcode. (Conservative)
      .add(relayGasBuffer)
  }

  let minGasLimit = dataCost.add(overhead)
  if (minGasLimit.gt(25_000_000)) {
    minGasLimit = BigNumber.from(25_000_000)
  }
  return minGasLimit
}

export const calculateWithdrawalMessage = (log: Log): WithdrawalMessageInfo => {
  if (typeof log.topics[1] === 'undefined') {
    throw new Error('"nonce" undefined')
  }
  const messageNonce = BigNumber.from(log.topics[1])

  if (typeof log.topics[2] === 'undefined') {
    throw new Error('"sender" undefined')
  }
  const sender = log.topics[2].substring(26)

  if (typeof log.topics[3] === 'undefined') {
    throw new Error(`"target" undefined`)
  }
  const target = log.topics[3].substring(26)

  const data = log.data
  if (data.length < 320) {
    throw new Error('Bad data')
  }

  const value = BigNumber.from(hexDataSlice(data, 0, 32))
  const minGasLimit = BigNumber.from(hexDataSlice(data, 32, 64))
  const dataOffset = BigNumber.from(hexDataSlice(data, 64, 96)).toNumber()
  const withdrawalHash = hexDataSlice(data, 96, 128)
  const dataLength = BigNumber.from(
    hexDataSlice(data, dataOffset, dataOffset + 32)
  ).toNumber()
  const message = hexDataSlice(
    data,
    dataOffset + 32,
    dataOffset + 32 + dataLength
  )

  return {
    messageNonce,
    sender: add0x(sender),
    target: add0x(target),
    value,
    minGasLimit,
    message,
    withdrawalHash,
    l2BlockNumber: log.blockNumber,
  }
}

export const calculateWithdrawalMessageUsingRecept = (
  txReceipt: TransactionReceipt
): WithdrawalMessageInfo => {
  if (txReceipt.status !== 1) {
    return undefined
  }
  let withdrawalMessage: WithdrawalMessageInfo
  txReceipt.logs.forEach((log) => {
    if (
      log.topics[0] ===
      ethers.utils.id(
        'MessagePassed(uint256,address,address,uint256,uint256,bytes,bytes32)'
      )
    ) {
      withdrawalMessage = calculateWithdrawalMessage(log)
    }
  })
  return withdrawalMessage
}
