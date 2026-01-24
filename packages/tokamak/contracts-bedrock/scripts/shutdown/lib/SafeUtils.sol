// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {console} from "forge-std/console.sol";
import {Vm} from "forge-std/Vm.sol";
import {IGnosisSafe, Enum} from "../../interfaces/IGnosisSafe.sol";

/// @title SafeUtils
/// @notice Gnosis Safe execution utilities for shutdown scripts
/// @dev Provides helpers for executing transactions via Safe with threshold checks
library SafeUtils {
    /// @notice Foundry VM address
    address internal constant VM_ADDRESS =
        address(uint160(uint256(keccak256("hevm cheat code"))));

    // ========== Safe Transaction Execution ==========

    /// @notice Execute transaction via Gnosis Safe (threshold=1 only)
    /// @dev Only executes if threshold is 1 and caller is owner
    /// @param safe Safe contract instance
    /// @param target Target contract address
    /// @param data Calldata to execute
    /// @param label Label for logging
    /// @param callerPrivateKey Private key of the caller
    /// @return executed True if transaction was executed
    function execViaSafe(
        IGnosisSafe safe,
        address target,
        bytes memory data,
        string memory label,
        uint256 callerPrivateKey
    ) internal returns (bool executed) {
        uint256 threshold = safe.getThreshold();
        address caller = Vm(VM_ADDRESS).addr(callerPrivateKey);

        if (threshold != 1) {
            console.log(
                "Safe execution skipped (threshold > 1) for:",
                label
            );
            return false;
        }

        if (!safe.isOwner(caller)) {
            console.log(
                "Safe execution skipped (caller not owner) for:",
                label
            );
            return false;
        }

        // Build approval signature for threshold=1
        bytes memory signature = abi.encodePacked(
            uint256(uint160(caller)),
            bytes32(0),
            uint8(1)
        );

        executed = safe.execTransaction({
            to: target,
            value: 0,
            data: data,
            operation: Enum.Operation.Call,
            safeTxGas: 0,
            baseGas: 0,
            gasPrice: 0,
            gasToken: address(0),
            refundReceiver: payable(address(0)),
            signatures: signature
        });
    }

    /// @notice Execute transaction via Safe using caller address (from env)
    /// @param safe Safe contract instance
    /// @param target Target contract address
    /// @param data Calldata to execute
    /// @param label Label for logging
    /// @return executed True if transaction was executed
    function execViaSafeFromEnv(
        IGnosisSafe safe,
        address target,
        bytes memory data,
        string memory label
    ) internal returns (bool executed) {
        uint256 privateKey = Vm(VM_ADDRESS).envUint("PRIVATE_KEY");
        return execViaSafe(safe, target, data, label, privateKey);
    }

    // ========== Safe Transaction Logging ==========

    /// @notice Log Safe transaction details for manual signing
    /// @param safe Safe contract instance
    /// @param target Target contract address
    /// @param data Calldata to execute
    /// @param label Label for the transaction
    function logSafeTx(
        IGnosisSafe safe,
        address target,
        bytes memory data,
        string memory label
    ) internal view {
        uint256 nonce = safe.nonce();
        bytes32 txHash = safe.getTransactionHash(
            target,
            0,
            data,
            Enum.Operation.Call,
            0,
            0,
            0,
            address(0),
            address(0),
            nonce
        );

        console.log("--- Safe TX Preview ---");
        console.log("Label:", label);
        console.log("Safe:", address(safe));
        console.log("Target:", target);
        console.log("Nonce:", Vm(VM_ADDRESS).toString(nonce));
        console.log("SafeTxHash:");
        console.logBytes32(txHash);
    }

    /// @notice Log Safe ownership information
    /// @param safe Safe contract instance
    /// @param label Label for logging
    /// @param caller Caller address to check
    function logSafeOwnership(
        IGnosisSafe safe,
        string memory label,
        address caller
    ) internal view {
        console.log("--- Safe Ownership Info ---");
        console.log("Label:", label);
        console.log("Safe:", address(safe));
        console.log("Threshold:", safe.getThreshold());
        console.log("Caller:", caller);
        console.log("Is Owner:", safe.isOwner(caller));

        address[] memory owners = safe.getOwners();
        console.log("Owner Count:", owners.length);
        for (uint256 i = 0; i < owners.length; i++) {
            console.log("  Owner", i, ":", owners[i]);
        }
    }

    // ========== Owner-based Execution ==========

    /// @notice Execute function via owner (Safe or EOA)
    /// @dev Routes to Safe execution if owner is contract, direct call if EOA
    /// @param ownerToUse Owner address (Safe or EOA)
    /// @param deployerAddress Deployer's EOA address
    /// @param target Target contract address
    /// @param data Calldata to execute
    /// @param label Label for logging
    /// @param callerPrivateKey Private key of the caller
    /// @return success True if execution succeeded
    function execWithOwner(
        address ownerToUse,
        address deployerAddress,
        address target,
        bytes memory data,
        string memory label,
        uint256 callerPrivateKey
    ) internal returns (bool success) {
        if (_isContract(ownerToUse)) {
            console.log("[INFO] Owner is a contract (Safe). Executing via Safe...");
            IGnosisSafe safe = IGnosisSafe(ownerToUse);

            logSafeTx(safe, target, data, label);

            success = execViaSafe(safe, target, data, label, callerPrivateKey);
            if (success) {
                console.log("[SUCCESS]", label, "executed via Safe");
            } else {
                console.log(
                    "[WARN] Safe execution skipped. Manual signing required."
                );
            }
        } else {
            bool isDryRun = Vm(VM_ADDRESS).envOr("DRY_RUN", false);
            if (ownerToUse != deployerAddress && !isDryRun) {
                revert("SafeUtils: caller is not owner");
            }
            console.log("[ACTION] Executing direct call via EOA...");
            (success, ) = target.call(data);
            require(success, "SafeUtils: direct call failed");
            console.log("[SUCCESS]", label, "executed via EOA");
        }
    }

    // ========== Proxy Upgrade Helpers ==========

    /// @notice Upgrade proxy via Safe or EOA
    /// @param proxy Proxy contract address
    /// @param proxyAdmin ProxyAdmin contract address
    /// @param systemOwnerSafe System owner Safe address (or address(0) for EOA)
    /// @param newImpl New implementation address
    /// @param label Label for logging
    /// @param callerPrivateKey Private key of the caller
    function upgradeProxyWithSafe(
        address proxy,
        address proxyAdmin,
        address systemOwnerSafe,
        address newImpl,
        string memory label,
        uint256 callerPrivateKey
    ) internal {
        console.log("[INFO] Proxy:", proxy);
        console.log("[INFO] Proxy Admin:", proxyAdmin);
        console.log("[INFO] New Implementation:", newImpl);

        bytes memory upgradeData = abi.encodeWithSignature(
            "upgrade(address,address)",
            proxy,
            newImpl
        );

        // Determine the owner to use
        address adminOwner = _getProxyAdminOwner(proxyAdmin);
        address ownerToUse = systemOwnerSafe != address(0)
            ? systemOwnerSafe
            : adminOwner;

        console.log("[INFO] ProxyAdmin Owner:", adminOwner);
        console.log("[INFO] Using Owner:", ownerToUse);

        address deployerAddress = Vm(VM_ADDRESS).addr(callerPrivateKey);

        if (_isContract(ownerToUse)) {
            console.log("[INFO] Owner is a contract (Safe). Preparing Safe TX...");
            IGnosisSafe safe = IGnosisSafe(ownerToUse);

            logSafeTx(safe, proxyAdmin, upgradeData, label);

            bool success = execViaSafe(
                safe,
                proxyAdmin,
                upgradeData,
                label,
                callerPrivateKey
            );
            if (success) {
                console.log("[SUCCESS] Proxy upgraded via Safe");
            } else {
                console.log(
                    "[WARN] Safe execution skipped (Threshold > 1 or simulation). Manual signing required."
                );
            }
        } else {
            bool isDryRun = Vm(VM_ADDRESS).envOr("DRY_RUN", false);
            if (ownerToUse != deployerAddress && !isDryRun) {
                revert("SafeUtils: caller is not owner");
            }
            console.log("[ACTION] Executing direct upgrade via EOA...");
            (bool success, ) = proxyAdmin.call(upgradeData);
            require(success, "SafeUtils: proxy upgrade failed");
            console.log("[SUCCESS] Proxy upgraded successfully");
        }
    }

    /// @notice Upgrade proxy directly (using upgradeTo on proxy)
    /// @param proxy Proxy contract address
    /// @param admin Admin address (Safe or EOA)
    /// @param newImpl New implementation address
    /// @param deployerAddress Deployer's EOA address
    /// @param label Label for logging
    /// @param callerPrivateKey Private key of the caller
    function upgradeProxyDirect(
        address proxy,
        address admin,
        address newImpl,
        address deployerAddress,
        string memory label,
        uint256 callerPrivateKey
    ) internal {
        console.log("[INFO] Proxy:", proxy);
        console.log("[INFO] Admin:", admin);
        console.log("[INFO] New Implementation:", newImpl);

        bytes memory upgradeData = abi.encodeWithSignature(
            "upgradeTo(address)",
            newImpl
        );

        if (_isContract(admin)) {
            IGnosisSafe safe = IGnosisSafe(admin);
            logSafeTx(safe, proxy, upgradeData, label);

            bool success = execViaSafe(
                safe,
                proxy,
                upgradeData,
                label,
                callerPrivateKey
            );
            require(success, "SafeUtils: safe execTransaction failed");
            console.log("[SUCCESS] Proxy upgraded via Safe (direct admin)");
            return;
        }

        require(admin == deployerAddress, "SafeUtils: caller is not admin");
        console.log("[ACTION] Executing direct upgrade via EOA...");
        (bool success, ) = proxy.call(upgradeData);
        require(success, "SafeUtils: proxy upgradeTo failed");
        console.log("[SUCCESS] Proxy upgraded via admin EOA");
    }

    // ========== Internal Helpers ==========

    /// @notice Check if address is a contract
    /// @param addr Address to check
    /// @return True if address has code
    function _isContract(address addr) internal view returns (bool) {
        return addr.code.length > 0;
    }

    /// @notice Get ProxyAdmin owner
    /// @param proxyAdmin ProxyAdmin contract address
    /// @return Owner address
    function _getProxyAdminOwner(
        address proxyAdmin
    ) internal view returns (address) {
        // ProxyAdmin.owner() selector: 0x8da5cb5b
        (bool success, bytes memory data) = proxyAdmin.staticcall(
            abi.encodeWithSignature("owner()")
        );
        require(success, "SafeUtils: failed to get ProxyAdmin owner");
        return abi.decode(data, (address));
    }

    /// @notice Public wrapper for contract check
    /// @param addr Address to check
    /// @return True if address has code
    function isContract(address addr) internal view returns (bool) {
        return _isContract(addr);
    }
}
