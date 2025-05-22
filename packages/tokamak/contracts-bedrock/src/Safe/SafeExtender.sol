// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { GnosisSafe as Safe } from "safe-contracts/GnosisSafe.sol";
import { FallbackManager } from "safe-contracts/base/FallbackManager.sol";

/**
 * @title SafeExtender
 * @dev Extends the Safe contract to add a getter for the fallback handler
 * @notice This contract is only meant to be used as a minimal extension for Safe to expose the fallback handler
 */
contract SafeExtender is Safe {

    /**
     * @notice Gets the fallback handler for the Safe
     * @return The address of the current fallback handler
     */
    function getFallbackHandler() public view returns (address) {
        address handler;
        assembly {
            handler := sload(FALLBACK_HANDLER_STORAGE_SLOT)
        }
        return handler;
    }
}