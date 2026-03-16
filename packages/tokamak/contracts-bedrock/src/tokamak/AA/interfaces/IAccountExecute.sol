// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// Ported from eth-infinitism/account-abstraction v0.8.0
// https://github.com/eth-infinitism/account-abstraction
// License: GPL-3.0
// Thanos modifications: import path updates, comment/error message updates only.

import "./PackedUserOperation.sol";

interface IAccountExecute {
    /**
     * Account may implement this execute method.
     * passing this methodSig at the beginning of callData will cause the entryPoint to pass the full UserOp (and hash)
     * to the account.
     * The account should skip the methodSig, and use the callData (and optionally, other UserOp fields)
     *
     * @param userOp              - The operation that was just validated.
     * @param userOpHash          - Hash of the user's request data.
     */
    function executeUserOp(
        PackedUserOperation calldata userOp,
        bytes32 userOpHash
    ) external;
}
