// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { CrossDomainMessenger } from "src/universal/CrossDomainMessenger.sol";
import { ISemver } from "src/universal/ISemver.sol";
import { Constants } from "src/libraries/Constants.sol";

/// @custom:proxied
/// @title L1CrossDomainMessenger
/// @notice The L1CrossDomainMessenger is a message passing interface between L1 and L2 responsible
///         for sending and receiving data on the L1 side. Users are encouraged to use this
///         interface instead of interacting with lower-level contracts directly.
contract L1CrossDomainMessenger is CrossDomainMessenger, ISemver {
    using SafeERC20 for IERC20;

    /// @notice Address of the OptimismPortal. The public getter for this
    ///         is legacy and will be removed in the future. Use `portal()` instead.
    /// @custom:network-specific
    /// @custom:legacy
    OptimismPortal public PORTAL;

    /// @notice Address of TON (ERC-20 token)
    address public tonAddress;

    /// @notice Semantic version.
    /// @custom:semver 1.7.1
    string public constant version = "1.7.1";

    /// @notice Constructs the L1CrossDomainMessenger contract.
    constructor() CrossDomainMessenger(Predeploys.L2_CROSS_DOMAIN_MESSENGER) {
        initialize({ _portal: OptimismPortal(payable(0)),  _tonAddress: address(0) });
    }

    /// @notice Initializes the contract.
    /// @param _portal Address of the OptimismPortal contract on this network.
    function initialize(OptimismPortal _portal, address _tonAddress) public reinitializer(Constants.INITIALIZER) {
        PORTAL = _portal;
        tonAddress = _tonAddress;
        __CrossDomainMessenger_init();
    }

    /// @notice Getter for the OptimismPortal address.
    function portal() external view returns (address) {
        return address(PORTAL);
    }

    /// @inheritdoc CrossDomainMessenger
    function _sendMessage(address _to, uint64 _gasLimit, uint256 _value, bytes memory _data) internal override {
        PORTAL.depositTransaction(_to, _value, _gasLimit, false, _data);
    }

    function _sendDepositTONMessage(address _to, uint64 _gasLimit, uint256 _value, bytes memory _data) internal override {
        IERC20(tonAddress).safeTransferFrom(msg.sender, address(this), _value);
        IERC20(tonAddress).approve(address(PORTAL), _value);
        PORTAL.depositTONTransaction(_to, _value, _gasLimit, false, _data);
    }


    /// @inheritdoc CrossDomainMessenger
    function _isOtherMessenger() internal view override returns (bool) {
        return msg.sender == address(PORTAL) && PORTAL.l2Sender() == OTHER_MESSENGER;
    }

    /// @inheritdoc CrossDomainMessenger
    function _isUnsafeTarget(address _target) internal view override returns (bool) {
        return _target == address(this) || _target == address(PORTAL);
    }
}
