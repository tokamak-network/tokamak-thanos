// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.15;

import { Proxy } from "../proxy/Proxy.sol";
import { L2FastWithdrawStorage } from "./L2FastWithdrawStorage.sol";

contract L2FastWithdrawProxy is Proxy, L2FastWithdrawStorage {

}
