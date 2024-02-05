import {
  hexStringEquals,
  predeploys,
  L1Predeploys,
} from '@eth-optimism/core-utils'
import { ethers } from 'ethers'

import { AddressLike } from '../interfaces'
import { toAddress } from './coercion'

/**
 * Specifically filter out ETH. ETH deposits and withdrawals are handled by the ETH bridge
 * adapter. Bridges that are not the ETH bridge should not be able to handle or even
 * present ETH deposits or withdrawals.
 */
export const filterOutEthDepositsAndWithdrawls = (
  l1Token: AddressLike,
  l2Token: AddressLike
): boolean => {
  if (
    hexStringEquals(toAddress(l1Token), ethers.constants.AddressZero) ||
    hexStringEquals(toAddress(l2Token), predeploys.ETH)
  ) {
    return false
  }
}

/**
 * Specifically filter out TON. TON deposits and withdrawals are handled by the TON bridge
 * adapter. Bridges that are not the TON bridge should not be able to handle or even
 * present TON deposits or withdrawals.
 */
export const filterOutTonDepositsAndWithdrawls = (
  l1Token: AddressLike,
  l2Token: AddressLike
): boolean => {
  if (
    hexStringEquals(toAddress(l1Token), L1Predeploys.L1TonAddress) ||
    hexStringEquals(toAddress(l2Token), ethers.constants.AddressZero)
  ) {
    return false
  }
}

/**
 * Filter ETH deposits and withdrawls
 *
 * @param l1Token
 * @param l2Token
 * @returns
 */
export const filterEthDepositsAndWithdrawls = (
  l1Token: AddressLike,
  l2Token: AddressLike
): boolean => {
  return !filterOutEthDepositsAndWithdrawls(l1Token, l2Token)
}

/**
 * Filter TON deposits and withdrawls
 *
 * @param l1Token
 * @param l2Token
 * @returns
 */
export const filterTonDepositsAndWithdrawls = (
  l1Token: AddressLike,
  l2Token: AddressLike
): boolean => {
  return !filterOutTonDepositsAndWithdrawls(l1Token, l2Token)
}
