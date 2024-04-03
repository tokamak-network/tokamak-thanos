// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import { StandardBridgeStorage } from "../../universal/StandardBridgeStorage.sol";

contract L2UsdcBridgeStorage is StandardBridgeStorage {

    address public l2UsdcMasterMinter;

}
