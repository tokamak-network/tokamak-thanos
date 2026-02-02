// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {OptimismPortal} from "src/L1/OptimismPortal.sol";

interface IShutdownProxyAdminOwner {
    function owner() external view returns (address);
}

contract ShutdownOptimismPortal is OptimismPortal {
    using SafeERC20 for IERC20;

    /// @dev EIP-1967 admin slot for proxy.
    bytes32 internal constant ADMIN_SLOT =
        0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    modifier onlyProxyAdminOwner() {
        address admin = _getProxyAdmin();
        require(admin != address(0), "ShutdownOptimismPortal: admin not set");
        if (admin.code.length > 0) {
            try IShutdownProxyAdminOwner(admin).owner() returns (address adminOwner) {
                require(msg.sender == adminOwner, "ShutdownOptimismPortal: unauthorized");
            } catch {
                revert("ShutdownOptimismPortal: admin owner lookup failed");
            }
        } else {
            require(msg.sender == admin, "ShutdownOptimismPortal: unauthorized");
        }
        _;
    }

    /// @notice Sweep native token (ERC20) to a target address.
    /// @dev Restricted to proxy admin owner.
    function sweepNativeToken(address _to, uint256 _amount) external onlyProxyAdminOwner {
        address nativeToken = _nativeToken();
        require(nativeToken != address(0), "ShutdownOptimismPortal: native token is ETH");
        IERC20(nativeToken).safeTransfer(_to, _amount);
    }

    function _getProxyAdmin() internal view returns (address admin_) {
        bytes32 slot = ADMIN_SLOT;
        assembly {
            admin_ := sload(slot)
        }
    }
}
