// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.28;

// Ported from eth-infinitism/account-abstraction v0.8.0
// https://github.com/eth-infinitism/account-abstraction
// License: GPL-3.0
// Thanos modifications: constructor → initialize() for predeploy compatibility (0x4200...0064).
//   Uses BasePaymaster._setEntryPoint() (Thanos lib modification).

import "@openzeppelin/contracts_v5.0.1/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts_v5.0.1/utils/cryptography/MessageHashUtils.sol";
import "@openzeppelin/contracts_v5.0.1/proxy/utils/Initializable.sol";
import "./lib/BasePaymaster.sol";
import "./lib/Helpers.sol";
import "./lib/UserOperationLib.sol";
import "./interfaces/IEntryPoint.sol";

/* solhint-disable not-rely-on-time */

/**
 * A sample paymaster that uses external service to decide whether to pay for the UserOp.
 * The paymaster trusts an external signer to sign the transaction.
 * The calling user must pass the UserOp to that external signer first, which performs
 * whatever off-chain verification before signing the UserOp.
 * Note that this signature is NOT a replacement for the account-specific signature:
 * - the paymaster checks a signature to agree to PAY for GAS.
 * - the account checks a signature to prove identity and authorize a specific transaction.
 */
contract VerifyingPaymaster is BasePaymaster, Initializable {
    using ECDSA for bytes32;
    using UserOperationLib for PackedUserOperation;

    address public verifyingSigner;

    uint256 private constant VALID_TIMESTAMP_OFFSET = PAYMASTER_DATA_OFFSET;
    uint256 private constant SIGNATURE_OFFSET = VALID_TIMESTAMP_OFFSET + 64;

    event VerifyingSignerUpdated(address indexed oldSigner, address indexed newSigner);

    // constructor() inherited from BasePaymaster: empty, owner = msg.sender

    /**
     * Initialize the paymaster. Called once after deployment (predeploy pattern).
     * @param _entryPoint      - The EntryPoint contract address.
     * @param _verifyingSigner - The off-chain signer that signs paymaster approvals.
     */
    function initialize(
        IEntryPoint _entryPoint,
        address _verifyingSigner
    ) external initializer {
        require(_verifyingSigner != address(0), "VerifyingPaymaster: zero signer");
        _setEntryPoint(_entryPoint); // Thanos-adapted BasePaymaster method
        verifyingSigner = _verifyingSigner;
    }

    /**
     * Update the verifying signer address (only owner).
     * @param _newSigner - The new signer address.
     */
    function setVerifyingSigner(address _newSigner) external onlyOwner {
        require(_newSigner != address(0), "VerifyingPaymaster: zero signer");
        emit VerifyingSignerUpdated(verifyingSigner, _newSigner);
        verifyingSigner = _newSigner;
    }

    /**
     * Return the hash we're going to sign off-chain (and validate on-chain).
     * This method is called by the off-chain service, to sign the request.
     * It is called on-chain from the validatePaymasterUserOp, to validate the signature.
     * Note that this signature covers all fields of the UserOperation, except the "paymasterAndData",
     * which will carry the signature itself.
     *
     * @param userOp     - The user operation.
     * @param validUntil - The expiry timestamp.
     * @param validAfter - The start timestamp.
     * @return - The hash to be signed.
     */
    function getHash(
        PackedUserOperation calldata userOp,
        uint48 validUntil,
        uint48 validAfter
    ) public view returns (bytes32) {
        // Can't use userOp.hash(), since it contains also the paymasterAndData itself.
        address sender = userOp.sender;
        return
            keccak256(
                abi.encode(
                    sender,
                    userOp.nonce,
                    keccak256(userOp.initCode),
                    keccak256(userOp.callData),
                    userOp.accountGasLimits,
                    userOp.preVerificationGas,
                    userOp.gasFees,
                    // Paymaster fields: only address + gas limits (no data/sig)
                    block.chainid,
                    address(this),
                    validUntil,
                    validAfter
                )
            );
    }

    /**
     * Verify our external signer signed this request.
     * The "paymasterAndData" is formatted as:
     *   [paymaster address (20 bytes)][validationGasLimit (16 bytes)][postOpGasLimit (16 bytes)]
     *   [validUntil (6 bytes)][validAfter (6 bytes)][signature (65+ bytes)]
     */
    function _validatePaymasterUserOp(
        PackedUserOperation calldata userOp,
        bytes32 /*userOpHash*/,
        uint256 /*maxCost*/
    ) internal view override returns (bytes memory context, uint256 validationData) {
        (uint48 validUntil, uint48 validAfter) = parsePaymasterAndData(userOp.paymasterAndData);
        // Only sign the hash of valid fields (we will pass validUntil,validAfter in the paymasterAndData,
        // so they must be signed with as well.)
        bytes32 hash = MessageHashUtils.toEthSignedMessageHash(getHash(userOp, validUntil, validAfter));
        bytes calldata paymasterSignature = userOp.paymasterAndData[SIGNATURE_OFFSET:];
        // Don't revert on signature failure: return SIG_VALIDATION_FAILED with a time range.
        if (verifyingSigner != ECDSA.recover(hash, paymasterSignature)) {
            return ("", _packValidationData(true, validUntil, validAfter));
        }
        // No need for other validation: paymaster is not using userOp.callData.
        return ("", _packValidationData(false, validUntil, validAfter));
    }

    /**
     * Parse the paymaster data: extract validUntil and validAfter timestamps.
     * @param paymasterAndData - Raw paymasterAndData bytes.
     * @return validUntil - The expiry timestamp.
     * @return validAfter - The start timestamp.
     */
    function parsePaymasterAndData(
        bytes calldata paymasterAndData
    ) public pure returns (uint48 validUntil, uint48 validAfter) {
        (validUntil, validAfter) = abi.decode(
            paymasterAndData[VALID_TIMESTAMP_OFFSET:SIGNATURE_OFFSET],
            (uint48, uint48)
        );
    }
}
