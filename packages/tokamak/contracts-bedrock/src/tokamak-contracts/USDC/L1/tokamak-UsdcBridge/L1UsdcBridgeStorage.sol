// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { StandardBridgeStorage } from "../../universal/StandardBridgeStorage.sol";

contract L1UsdcBridgeStorage is StandardBridgeStorage {

    mapping(address => mapping(address => uint256)) public deposits;

}
