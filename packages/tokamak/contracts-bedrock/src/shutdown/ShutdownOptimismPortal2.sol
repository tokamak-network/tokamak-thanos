// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {OptimismPortal2} from "src/L1/OptimismPortal2.sol";

interface IShutdownProxyAdminOwner2 {
    function owner() external view returns (address);
}

contract ShutdownOptimismPortal2 is OptimismPortal2 {
    using SafeERC20 for IERC20;

    /// @dev EIP-1967 admin slot for proxy.
    bytes32 internal constant ADMIN_SLOT =
        0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    constructor(uint256 _proofMaturityDelaySeconds, uint256 _disputeGameFinalityDelaySeconds)
        OptimismPortal2(_proofMaturityDelaySeconds, _disputeGameFinalityDelaySeconds)
    {}

    modifier onlyProxyAdminOwner() {
        address admin = _getProxyAdmin();
        require(admin != address(0), "ShutdownOptimismPortal2: admin not set");
        if (admin.code.length > 0) {
            try IShutdownProxyAdminOwner2(admin).owner() returns (address adminOwner) {
                require(msg.sender == adminOwner, "ShutdownOptimismPortal2: unauthorized");
            } catch {
                revert("ShutdownOptimismPortal2: admin owner lookup failed");
            }
        } else {
            require(msg.sender == admin, "ShutdownOptimismPortal2: unauthorized");
        }
        _;
    }

    /// @notice Sweep native token (ERC20) to a target address.
    /// @dev Restricted to proxy admin owner.
    function sweepNativeToken(address _to, uint256 _amount) external onlyProxyAdminOwner {
        address nativeToken = _nativeToken();
        require(nativeToken != address(0), "ShutdownOptimismPortal2: native token is ETH");
        IERC20(nativeToken).safeTransfer(_to, _amount);
    }

    function _getProxyAdmin() internal view returns (address admin_) {
        bytes32 slot = ADMIN_SLOT;
        assembly {
            admin_ := sload(slot)
        }
    }
}
