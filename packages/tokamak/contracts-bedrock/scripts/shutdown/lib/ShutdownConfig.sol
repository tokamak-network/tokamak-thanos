// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;



/// @title ShutdownConfig
/// @notice Centralized configuration for L2 shutdown scripts
/// @dev Re-exports constants from existing libraries and defines shutdown-specific configurations
library ShutdownConfig {
    // ========== Re-export from Predeploys ==========

    /// @notice ETH predeploy on L2 (see Predeploys.ETH)
    address internal constant PREDEPLOY_ETH = 0x4200000000000000000000000000000000000486;

    /// @notice USDC predeploy on L2 (FiatTokenV2_2) (see Predeploys.FIATTOKENV2_2)
    address internal constant PREDEPLOY_USDC = 0x4200000000000000000000000000000000000778;

    /// @notice L2 Standard Bridge predeploy (see Predeploys.L2_STANDARD_BRIDGE)
    address internal constant L2_STANDARD_BRIDGE = 0x4200000000000000000000000000000000000010;

    /// @notice L2 USDC Bridge predeploy (see Predeploys.L2_USDC_BRIDGE)
    address internal constant L2_USDC_BRIDGE = 0x4200000000000000000000000000000000000775;

    /// @notice L2 to L1 Message Passer predeploy (see Predeploys.L2_TO_L1_MESSAGE_PASSER)
    address internal constant L2_TO_L1_MESSAGE_PASSER = 0x4200000000000000000000000000000000000016;

    // ========== Re-export from Constants ==========

    /// @notice EIP-1967 admin slot (see Constants.PROXY_OWNER_ADDRESS)
    bytes32 internal constant PROXY_ADMIN_SLOT =
        0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    /// @notice EIP-1967 implementation slot (see Constants.PROXY_IMPLEMENTATION_ADDRESS)
    bytes32 internal constant PROXY_IMPL_SLOT =
        0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

    // ========== Shutdown-specific Constants ==========

    /// @notice Multicall3 address (same on all EVM chains)
    address internal constant MULTICALL3 = 0xcA11bde05977b3631167028862bE2a173976CA11;

    /// @notice Native token address placeholder (address(0) represents native)
    address internal constant NATIVE_TOKEN = address(0);

    /// @notice Forge default sender address (used in tests and scripts)
    /// @dev This is the default msg.sender when running forge scripts without --sender flag
    ///      See: https://book.getfoundry.sh/reference/config/testing#sender
    address internal constant FORGE_DEFAULT_SENDER = 0x1804c8AB1F12E6bbf3894d4083f33e07309d1f38;

    // ========== Default Values ==========

    /// @notice Default batch size for multicall operations
    uint256 internal constant DEFAULT_MULTICALL_BATCH_SIZE = 200;

    /// @notice Default tolerance for native balance comparison (0.01 ETH)
    uint256 internal constant DEFAULT_NATIVE_BALANCE_TOLERANCE = 1e16;

    // ========== Environment Variable Keys ==========

    string internal constant ENV_L1_RPC_URL = "L1_RPC_URL";
    string internal constant ENV_L2_RPC_URL = "L2_RPC_URL";
    string internal constant ENV_BRIDGE_PROXY = "BRIDGE_PROXY";
    string internal constant ENV_OPTIMISM_PORTAL_PROXY = "OPTIMISM_PORTAL_PROXY";
    string internal constant ENV_L1_NATIVE_TOKEN = "L1_NATIVE_TOKEN";
    string internal constant ENV_L1_USDC_BRIDGE_PROXY = "L1_USDC_BRIDGE_PROXY";
    string internal constant ENV_L2_USDC_BRIDGE_PROXY = "L2_USDC_BRIDGE_PROXY";
    string internal constant ENV_PROXY_ADMIN = "PROXY_ADMIN";
    string internal constant ENV_SYSTEM_OWNER_SAFE = "SYSTEM_OWNER_SAFE";
    string internal constant ENV_PRIVATE_KEY = "PRIVATE_KEY";
    string internal constant ENV_CLOSER_ADDRESS_PRIVATE_KEY = "CLOSER_ADDRESS_PRIVATE_KEY";
    string internal constant ENV_DATA_PATH = "DATA_PATH";
    string internal constant ENV_DRY_RUN = "DRY_RUN";

    // ========== Configuration Structs ==========

    /// @notice RPC endpoint configuration
    struct RpcConfig {
        string l1RpcUrl;
        string l2RpcUrl;
    }

    /// @notice Bridge addresses configuration
    struct BridgeConfig {
        address l1Bridge;
        address optimismPortal;
        address l1UsdcBridge;
        address l2UsdcBridge;
        address l1NativeToken;
        address proxyAdmin;
        address systemOwnerSafe;
    }

    /// @notice Chain-specific configuration
    struct ChainConfig {
        uint256 l2ChainId;
        uint256 finalizedNativeWithdrawals;
        uint256 multicallBatchSize;
        uint256 nativeBalanceTolerance;
    }

    /// @notice Complete configuration for shutdown scripts
    struct FullConfig {
        RpcConfig rpc;
        BridgeConfig bridge;
        ChainConfig chain;
    }

    // ========== Helper Functions ==========

    /// @notice Check if address is a contract
    /// @param addr Address to check
    /// @return True if address has code
    function isContract(address addr) internal view returns (bool) {
        return addr.code.length > 0;
    }

    /// @notice Get EIP-1967 admin from proxy using vm.load
    /// @dev This function signature is for documentation; actual implementation
    ///      requires vm.load() which is only available in Forge scripts
    /// @param proxy Proxy contract address
    /// @return admin Admin address from storage slot
    function getEip1967AdminSlot(
        address proxy
    ) internal pure returns (bytes32) {
        // Returns the slot to be used with vm.load(proxy, slot)
        // Usage: address admin = address(uint160(uint256(vm.load(proxy, ShutdownConfig.PROXY_ADMIN_SLOT))));
        return PROXY_ADMIN_SLOT;
    }

    /// @notice Calculate absolute difference between two values
    /// @param a First value
    /// @param b Second value
    /// @return Absolute difference
    function abs(uint256 a, uint256 b) internal pure returns (uint256) {
        return a > b ? a - b : b - a;
    }
}
