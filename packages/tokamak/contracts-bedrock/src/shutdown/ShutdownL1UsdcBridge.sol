// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC20} from "src/USDC/L1/libraries/IERC20.sol";
import {SafeERC20} from "src/USDC/L1/libraries/SafeERC20.sol";
import {L1UsdcBridge} from "src/USDC/L1/tokamak-UsdcBridge/L1UsdcBridge.sol";

interface IShutdownProxyAdminOwnerUsdc {
    function owner() external view returns (address);
}

contract ShutdownL1UsdcBridge is L1UsdcBridge {
    using SafeERC20 for IERC20;

    /// @dev EIP-1967 admin slot for proxy.
    bytes32 internal constant ADMIN_SLOT =
        0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    modifier onlyProxyAdminOwner() {
        address admin = _getProxyAdmin();
        require(admin != address(0), "ShutdownL1UsdcBridge: admin not set");
        if (admin.code.length > 0) {
            try IShutdownProxyAdminOwnerUsdc(admin).owner() returns (address adminOwner) {
                require(msg.sender == adminOwner, "ShutdownL1UsdcBridge: unauthorized");
            } catch {
                revert("ShutdownL1UsdcBridge: admin owner lookup failed");
            }
        } else {
            require(msg.sender == admin, "ShutdownL1UsdcBridge: unauthorized");
        }
        _;
    }

    /// @notice Sweep USDC to a target address.
    /// @dev Restricted to proxy admin owner.
    function sweepUSDC(address _to, uint256 _amount) external onlyProxyAdminOwner {
        require(l1Usdc != address(0), "ShutdownL1UsdcBridge: l1Usdc not set");
        IERC20(l1Usdc).safeTransfer(_to, _amount);
    }

    function _getProxyAdmin() internal view returns (address admin_) {
        bytes32 slot = ADMIN_SLOT;
        assembly {
            admin_ := sload(slot)
        }
    }
}
