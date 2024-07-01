// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { CrossDomainMessenger } from "src/universal/CrossDomainMessenger.sol";
import { ISemver } from "src/universal/ISemver.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { Constants } from "src/libraries/Constants.sol";
import { SafeCall } from "src/libraries/SafeCall.sol";
import { Hashing } from "src/libraries/Hashing.sol";
import { Encoding } from "src/libraries/Encoding.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { OnApprove } from "./OnApprove.sol";

/// @custom:proxied
/// @title L1CrossDomainMessenger
/// @notice The L1CrossDomainMessenger is a message passing interface between L1 and L2 responsible
///         for sending and receiving data on the L1 side. Users are encouraged to use this
///         interface instead of interacting with lower-level contracts directly.
contract L1CrossDomainMessenger is CrossDomainMessenger, OnApprove, ISemver {
    using SafeERC20 for IERC20;

    /// @notice Contract of the SuperchainConfig.
    SuperchainConfig public superchainConfig;

    /// @notice Contract of the OptimismPortal.
    /// @custom:network-specific
    OptimismPortal public portal;

    /// @notice Address of the SystemConfig contract.
    /// @custom:network-specific
    SystemConfig public systemConfig;

    // /// @notice Address of native token (ERC-20 token)
    // address public nativeTokenAddress;

    /// @notice Semantic version.
    /// @custom:semver 2.3.0
    string public constant version = "2.3.0";

    /// @notice Constructs the L1CrossDomainMessenger contract.
    constructor() CrossDomainMessenger() {
        initialize({
            _superchainConfig: SuperchainConfig(address(0)),
            _portal: OptimismPortal(payable(0)),
            _systemConfig: SystemConfig(address(0))
        });
    }

    /// @notice Initializes the contract.
    /// @param _superchainConfig Contract of the SuperchainConfig contract on this network.
    /// @param _portal Contract of the OptimismPortal contract on this network.
    function initialize(
        SuperchainConfig _superchainConfig,
        OptimismPortal _portal,
        SystemConfig _systemConfig
    )
        public
        reinitializer(Constants.INITIALIZER)
    {
        superchainConfig = _superchainConfig;
        portal = _portal;
        systemConfig = _systemConfig;
        __CrossDomainMessenger_init({ _otherMessenger: CrossDomainMessenger(Predeploys.L2_CROSS_DOMAIN_MESSENGER) });
    }

    /// @notice Getter for the OptimismPortal contract.
    /// @return Contract of the OptimismPortal on this chain.
    function PORTAL() external view returns (OptimismPortal) {
        return portal;
    }

    /// @inheritdoc CrossDomainMessenger
    function _sendMessage(address _to, uint64 _gasLimit, uint256 _value, bytes memory _data) internal override {
        require(msg.value == 0, "Deny depositing ETH");
        portal.depositTransaction(_to, _value, _value, _gasLimit, false, _data);
    }

    /// @notice unpack onApprove data
    /// @param _data     Data used in OnApprove contract
    function unpackOnApproveData(bytes calldata _data)
        public
        pure
        returns (address _from, address _to, uint256 _amount, uint32 _minGasLimit, bytes calldata _message)
    {
        require(_data.length >= 76, "Invalid onApprove data for L1CrossDomainMessenger");
        assembly {
            // The layout of a "bytes calldata" is:
            // The first 20 bytes: _from
            // The next 20 bytes: _to
            // The next 32 bytes: _amount
            // The next 4 bytes: _minGasLimit
            // The rest: _message
            _from := shr(96, calldataload(_data.offset))
            _to := shr(96, calldataload(add(_data.offset, 20)))
            _amount := calldataload(add(_data.offset, 40))
            _minGasLimit := shr(224, calldataload(add(_data.offset, 72)))
            _message.offset := add(_data.offset, 76)
            _message.length := sub(_data.length, 76)
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
        require(msg.sender == address(nativeTokenAddress()), "only accept native token approve callback");
        (address from, address to, uint256 amount, uint32 minGasLimit, bytes calldata message) =
            unpackOnApproveData(_data);
        require(_owner == from && _amount == amount && amount > 0, "invalid onApprove data");
        _sendNativeTokenMessage(from, to, amount, minGasLimit, message);
        return true;
    }

    /// @inheritdoc CrossDomainMessenger
    function _isOtherMessenger() internal view override returns (bool) {
        return msg.sender == address(portal) && portal.l2Sender() == address(otherMessenger);
    }

    /// @inheritdoc CrossDomainMessenger
    function _isUnsafeTarget(address _target) internal view override returns (bool) {
        return _target == address(this) || _target == address(portal);
    }

    function nativeTokenAddress() public view returns (address) {
        return systemConfig.nativeTokenAddress();
    }

    /// @inheritdoc CrossDomainMessenger
    function paused() public view override returns (bool) {
        return superchainConfig.paused();
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
            address _nativeTokenAddress = nativeTokenAddress();
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

        address _nativeTokenAddress = nativeTokenAddress();
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
        // 1. Call the target contract (_minGasLimit + RELAY_CALL_OVERHEAD + RELAY_GAS_CHECK_BUFFER)
        //   1.a. The RELAY_CALL_OVERHEAD is included in `hasMinGas`.
        // 2. Finish the execution after the external call (RELAY_RESERVED_GAS).
        //
        // If `xDomainMsgSender` is not the default L2 sender, this function
        // is being re-entered. This marks the message as failed to allow it to be replayed.
        if (
            !SafeCall.hasMinGas(_minGasLimit, RELAY_RESERVED_GAS + RELAY_GAS_CHECK_BUFFER)
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
        bool approvalStatus = true;
        if (_value != 0) {
            approvalStatus = IERC20(_nativeTokenAddress).approve(_target, _value);
        }
        bool success = SafeCall.call(_target, gasleft() - RELAY_RESERVED_GAS, 0, _message);
        xDomainMsgSender = Constants.DEFAULT_L2_SENDER;
        if (_value != 0 && !success) {
            approvalStatus = IERC20(_nativeTokenAddress).approve(_target, 0);
        }

        if (success && approvalStatus) {
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
