// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// OpenZeppelin
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

// Contracts
import { ProxyAdminOwnedBase } from "src/L1/ProxyAdminOwnedBase.sol";
import { ReinitializableBase } from "src/universal/ReinitializableBase.sol";
import { CrossDomainMessenger } from "src/universal/CrossDomainMessenger.sol";
import { OnApprove } from "./OnApprove.sol";

// Libraries
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Constants } from "src/libraries/Constants.sol";
import { SafeCall } from "src/libraries/SafeCall.sol";
import { Hashing } from "src/libraries/Hashing.sol";
import { Encoding } from "src/libraries/Encoding.sol";

// Interfaces
import { ISemver } from "src/universal/ISemver.sol";
import { ISystemConfig } from "interfaces/L1/ISystemConfig.sol";
import { IOptimismPortal2 as IOptimismPortal } from "interfaces/L1/IOptimismPortal2.sol";
import { ISuperchainConfig } from "interfaces/L1/ISuperchainConfig.sol";

/// @custom:proxied true
/// @title L1CrossDomainMessenger
/// @notice The L1CrossDomainMessenger is a message passing interface between L1 and L2 responsible
///         for sending and receiving data on the L1 side. Users are encouraged to use this
///         interface instead of interacting with lower-level contracts directly.
/// @dev This version includes Tokamak's Native Token functionality while following Optimism v1.16.0 patterns.
contract L1CrossDomainMessenger is
    CrossDomainMessenger,
    ProxyAdminOwnedBase,
    ReinitializableBase,
    OnApprove,
    ISemver
{
    using SafeERC20 for IERC20;

    /// @custom:legacy
    /// @custom:spacer superchainConfig
    /// @notice Spacer taking up the legacy `superchainConfig` slot.
    address private spacer_251_0_20;

    /// @notice Contract of the OptimismPortal.
    /// @custom:network-specific
    IOptimismPortal public portal;

    /// @custom:legacy
    /// @custom:spacer systemConfig (old location)
    /// @notice Spacer taking up the legacy `systemConfig` slot.
    address private spacer_253_0_20;

    /// @notice Semantic version.
    /// @custom:semver 2.9.0-tokamak
    string public constant version = "2.9.0-tokamak";

    /// @notice Contract of the SystemConfig.
    ISystemConfig public systemConfig;

    /// @notice Constructs the L1CrossDomainMessenger contract.
    constructor() ReinitializableBase(2) {
        _disableInitializers();
    }

    /// @notice Initializes the contract.
    /// @param _superchainConfig Contract of the SuperchainConfig contract on this network (kept for Tokamak compatibility).
    /// @param _portal Contract of the OptimismPortal contract on this network.
    /// @param _systemConfig Contract of the SystemConfig contract on this network.
    function initialize(
        ISuperchainConfig _superchainConfig,
        IOptimismPortal _portal,
        ISystemConfig _systemConfig
    )
        external
        reinitializer(initVersion())
    {
        // Initialization transactions must come from the ProxyAdmin or its owner.
        _assertOnlyProxyAdminOrProxyAdminOwner();

        // Now perform initialization logic.
        // Note: spacer_251_0_20 is where superchainConfig used to be stored, now we use SystemConfig
        portal = _portal;
        systemConfig = _systemConfig;
        __CrossDomainMessenger_init({ _otherMessenger: CrossDomainMessenger(Predeploys.L2_CROSS_DOMAIN_MESSENGER) });
    }

    /// @notice Upgrades the contract to have a reference to the SystemConfig.
    /// @param _systemConfig The new SystemConfig contract.
    function upgrade(ISystemConfig _systemConfig) external reinitializer(initVersion()) {
        // Upgrade transactions must come from the ProxyAdmin or its owner.
        _assertOnlyProxyAdminOrProxyAdminOwner();

        // Now perform upgrade logic.
        systemConfig = _systemConfig;
    }

    /// @inheritdoc CrossDomainMessenger
    function paused() public view override returns (bool) {
        return superchainConfig().paused();
    }

    /// @notice Returns the SuperchainConfig contract.
    /// @return ISuperchainConfig The SuperchainConfig contract.
    function superchainConfig() public view returns (ISuperchainConfig) {
        return systemConfig.superchainConfig();
    }

    /// @notice Getter function for the OptimismPortal contract on this chain.
    ///         Public getter is legacy and will be removed in the future. Use `portal()` instead.
    /// @return Contract of the OptimismPortal on this chain.
    /// @custom:legacy
    function PORTAL() external view returns (IOptimismPortal) {
        return portal;
    }

    /// @inheritdoc CrossDomainMessenger
    function _sendMessage(address _to, uint64 _gasLimit, uint256 _value, bytes memory _data) internal override {
        // Tokamak: Deny direct ETH deposits, only accept Native Token
        require(msg.value == 0, "Deny depositing ETH");
        // Tokamak: OptimismPortal expects 6 params (_to, _mint, _value, _gasLimit, _isCreation, _data)
        // _mint and _value are both set to _value for native token functionality
        portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
    }

    /// @inheritdoc CrossDomainMessenger
    function _isOtherMessenger() internal view override returns (bool) {
        return msg.sender == address(portal) && portal.l2Sender() == address(otherMessenger);
    }

    /// @notice Checks whether the target address is excluded from receiving messages.
    /// @param _target Address to check.
    /// @return Whether or not the target address is excluded.
    function _isUnsafeTarget(address _target) internal view override returns (bool) {
        return _target == address(this) || _target == address(portal);
    }

    /// @notice Getter function for address of native token on this network
    /// @return address The address of native token
    function nativeTokenAddress() public view returns (address) {
        return _nativeToken();
    }

    function _nativeToken() internal view returns (address) {
        // Tokamak: SystemConfig has nativeTokenAddress() function
        return ISystemConfig(address(systemConfig)).nativeTokenAddress();
    }

    // ==================== TOKAMAK NATIVE TOKEN FUNCTIONALITY ====================

    /// @notice unpack onApprove data
    /// @param _data     Data used in OnApprove contract
    function unpackOnApproveData(bytes calldata _data)
        internal
        pure
        returns (address _to, uint32 _minGasLimit, bytes calldata _message)
    {
        require(_data.length >= 24, "Invalid onApprove data for L1CrossDomainMessenger");
        assembly {
            // The layout of a "bytes calldata" is:
            // The next 20 bytes: _to
            // The next 4 bytes: _minGasLimit
            // The rest: _message
            _to := shr(96, calldataload(_data.offset))
            _minGasLimit := shr(224, calldataload(add(_data.offset, 20)))
            _message.offset := add(_data.offset, 24)
            _message.length := sub(_data.length, 24)
        }
    }

    /// @notice ERC20 onApprove callback
    /// @param _owner    Account that called approveAndCall
    /// @param _amount   Approved amount
    /// @param _data     Data used in OnApprove contract
    function onApprove(
        address _owner,
        address,
        uint256 _amount,
        bytes calldata _data
    )
        external
        override
        returns (bool)
    {
        require(msg.sender == address(_nativeToken()), "only accept native token approve callback");
        (address to, uint32 minGasLimit, bytes calldata message) = unpackOnApproveData(_data);
        _sendNativeTokenMessage(_owner, to, _amount, minGasLimit, message);
        return true;
    }

    /// @notice Sends a deposit native token message to some target address on the other chain. Note that if the call
    ///         always reverts, then the message will be unrelayable, and any ETH sent will be
    ///         permanently locked. The same will occur if the target on the other chain is
    ///         considered unsafe (see the _isUnsafeTarget() function).
    /// @param _target      Target contract or wallet address.
    /// @param _amount      Amount of deposit native token.
    /// @param _message     Message to trigger the target address with.
    /// @param _minGasLimit Minimum gas limit that the message can be executed with.
    function sendNativeTokenMessage(
        address _target,
        uint256 _amount,
        bytes calldata _message,
        uint32 _minGasLimit
    )
        external
    {
        // Triggers a message to the other messenger. Note that the amount of gas provided to the
        // message is the amount of gas requested by the user PLUS the base gas value. We want to
        // guarantee the property that the call to the target contract will always have at least
        // the minimum gas limit specified by the user.
        _sendNativeTokenMessage(msg.sender, _target, _amount, _minGasLimit, _message);
    }

    /// @notice Sends a deposit native token message internally to some target address on the other chain. Note that if
    /// the call
    ///         always reverts, then the message will be unrelayable, and any ETH sent will be
    ///         permanently locked. The same will occur if the target on the other chain is
    ///         considered unsafe (see the _isUnsafeTarget() function).
    /// @param _sender      Sender address.
    /// @param _target      Target contract or wallet address.
    /// @param _amount      Amount of deposit native token.
    /// @param _message     Message to trigger the target address with.
    /// @param _minGasLimit Minimum gas limit that the message can be executed with.
    function _sendNativeTokenMessage(
        address _sender,
        address _target,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _message
    )
        internal
    {
        // Collect native token
        if (_amount > 0) {
            address _nativeTokenAddress = _nativeToken();
            IERC20(_nativeTokenAddress).safeTransferFrom(_sender, address(this), _amount);
            IERC20(_nativeTokenAddress).approve(address(portal), _amount);
        }

        // Triggers a message to the other messenger. Note that the amount of gas provided to the
        // message is the amount of gas requested by the user PLUS the base gas value. We want to
        // guarantee the property that the call to the target contract will always have at least
        // the minimum gas limit specified by the user.
        _sendMessage(
            address(otherMessenger),
            baseGas(_message, _minGasLimit),
            _amount,
            abi.encodeWithSelector(
                this.relayMessage.selector, messageNonce(), _sender, _target, _amount, _minGasLimit, _message
            )
        );

        emit SentMessage(_target, _sender, _message, messageNonce(), _minGasLimit);
        emit SentMessageExtension1(_sender, _amount);

        unchecked {
            ++msgNonce;
        }
    }

    /// @notice Relays a message that was sent by the other CrossDomainMessenger contract. Can only
    ///         be executed via cross-chain call from the other messenger OR if the message was
    ///         already received once and is currently being replayed.
    /// @param _nonce       Nonce of the message being relayed.
    /// @param _sender      Address of the user who sent the message.
    /// @param _target      Address that the message is targeted at.
    /// @param _value       Native token value to send with the message.
    /// @param _minGasLimit Minimum amount of gas that the message can be executed with.
    /// @param _message     Message to send to the target.
    function relayMessage(
        uint256 _nonce,
        address _sender,
        address _target,
        uint256 _value,
        uint256 _minGasLimit,
        bytes calldata _message
    )
        external
        payable
        override
    {
        require(paused() == false, "L1 CrossDomainMessenger: paused");
        require(msg.value == 0, "CrossDomainMessenger: value must be zero");

        (, uint16 _nonceVersion) = Encoding.decodeVersionedNonce(_nonce);
        require(_nonceVersion < 2, "CrossDomainMessenger: only version 0 or 1 messages are supported at this time");

        // If the message is version 0, then it's a migrated legacy withdrawal. We therefore need
        // to check that the legacy version of the message has not already been relayed.
        if (_nonceVersion == 0) {
            bytes32 oldHash = Hashing.hashCrossDomainMessageV0(_target, _sender, _message, _nonce);
            require(successfulMessages[oldHash] == false, "CrossDomainMessenger: legacy withdrawal already relayed");
        }

        // We use the v1 message hash as the unique identifier for the message because it commits
        // to the value and minimum gas limit of the message.
        bytes32 versionedHash =
            Hashing.hashCrossDomainMessageV1(_nonce, _sender, _target, _value, _minGasLimit, _message);

        address _nativeTokenAddress = _nativeToken();
        if (_isOtherMessenger()) {
            // These properties should always hold when the message is first submitted (as
            // opposed to being replayed).
            assert(!failedMessages[versionedHash]);
            if (_value > 0) {
                IERC20(_nativeTokenAddress).safeTransferFrom(address(portal), address(this), _value);
            }
        } else {
            require(failedMessages[versionedHash], "CrossDomainMessenger: message cannot be replayed");
        }

        require(
            _isUnsafeTarget(_target) == false && _target != _nativeTokenAddress,
            "CrossDomainMessenger: cannot send message to blocked system address or nativeTokenAddress"
        );

        require(successfulMessages[versionedHash] == false, "CrossDomainMessenger: message has already been relayed");

        // If there is not enough gas left to perform the external call and finish the execution,
        // return early and assign the message to the failedMessages mapping.
        // We are asserting that we have enough gas to:
        // 1. Call the target contract (_minGasLimit + RELAY_CALL_OVERHEAD + RELAY_GAS_CHECK_BUFFER_INCLUDING_APPROVAL)
        //   1.a. The RELAY_CALL_OVERHEAD is included in `hasMinGas`.
        // 2. Finish the execution after the external call (RELAY_RESERVED_GAS).
        //
        // If `xDomainMsgSender` is not the default L2 sender, this function
        // is being re-entered. This marks the message as failed to allow it to be replayed.
        if (
            !SafeCall.hasMinGas(_minGasLimit, RELAY_RESERVED_GAS + RELAY_GAS_CHECK_BUFFER_INCLUDING_APPROVAL)
                || xDomainMsgSender != Constants.DEFAULT_L2_SENDER
        ) {
            failedMessages[versionedHash] = true;
            emit FailedRelayedMessage(versionedHash);

            // Revert in this case if the transaction was triggered by the estimation address. This
            // should only be possible during gas estimation or we have bigger problems. Reverting
            // here will make the behavior of gas estimation change such that the gas limit
            // computed will be the amount required to relay the message, even if that amount is
            // greater than the minimum gas limit specified by the user.
            if (tx.origin == Constants.ESTIMATION_ADDRESS) {
                revert("CrossDomainMessenger: failed to relay message");
            }
            return;
        }

        xDomainMsgSender = _sender;
        // _target must not be address(0). otherwise, this transaction could be reverted
        if (_value != 0 && _target != address(0)) {
            IERC20(_nativeTokenAddress).approve(_target, _value);
        }
        // _target is expected to perform a transferFrom to collect token
        bool success = SafeCall.call(_target, gasleft() - RELAY_RESERVED_GAS, 0, _message);
        if (_value != 0 && _target != address(0)) {
            IERC20(_nativeTokenAddress).approve(_target, 0);
        }
        xDomainMsgSender = Constants.DEFAULT_L2_SENDER;

        if (success) {
            // This check is identical to the one above, but it ensures that the same message cannot be relayed
            // twice, and adds a layer of protection against reentrancy.
            assert(successfulMessages[versionedHash] == false);
            successfulMessages[versionedHash] = true;
            emit RelayedMessage(versionedHash);
        } else {
            failedMessages[versionedHash] = true;
            emit FailedRelayedMessage(versionedHash);

            // Revert in this case if the transaction was triggered by the estimation address. This
            // should only be possible during gas estimation or we have bigger problems. Reverting
            // here will make the behavior of gas estimation change such that the gas limit
            // computed will be the amount required to relay the message, even if that amount is
            // greater than the minimum gas limit specified by the user.
            if (tx.origin == Constants.ESTIMATION_ADDRESS) {
                revert("CrossDomainMessenger: failed to relay message");
            }
        }
    }
}
