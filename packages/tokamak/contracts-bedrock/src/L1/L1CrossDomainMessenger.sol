// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { CrossDomainMessenger } from "src/universal/CrossDomainMessenger.sol";
import { ISemver } from "src/universal/ISemver.sol";
import { Constants } from "src/libraries/Constants.sol";
import { SafeCall } from "src/libraries/SafeCall.sol";
import { Hashing } from "src/libraries/Hashing.sol";
import { Encoding } from "src/libraries/Encoding.sol";
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
    address public l1TONAddress;

    /// @notice Semantic version.
    /// @custom:semver 1.7.1
    string public constant version = "1.7.1";

    /// @notice Constructs the L1CrossDomainMessenger contract.
    constructor() CrossDomainMessenger(Predeploys.L2_CROSS_DOMAIN_MESSENGER) {
        initialize({ _portal: OptimismPortal(payable(0)),  _l1TONAddress: address(0) });
    }

    /// @notice Initializes the contract.
    /// @param _portal Address of the OptimismPortal contract on this network.
    function initialize(OptimismPortal _portal, address _l1TONAddress) public reinitializer(Constants.INITIALIZER) {
        PORTAL = _portal;
        l1TONAddress = _l1TONAddress;
        __CrossDomainMessenger_init();
        if (address(PORTAL) != address(0) && address(l1TONAddress) != address(0)) {
            IERC20(l1TONAddress).approve(address(PORTAL), 2**256 - 1);
        }
    }

    /// @notice Getter for the OptimismPortal address.
    function portal() external view returns (address) {
        return address(PORTAL);
    }

    // /// @inheritdoc CrossDomainMessenger
    // function _sendMessage(address _to, uint64 _gasLimit, uint256 _value, bytes memory _data) internal override {
    //     require(msg.value == 0, "Deny depositing ETH");
    //     PORTAL.depositTransaction(_to, 0, _gasLimit, false, _data);
    // }

    /// @inheritdoc CrossDomainMessenger
    function _sendMessage(address _to, uint64 _gasLimit, uint256 _value, bytes memory _data) internal override{
        require(msg.value == 0, "Deny depositing ETH");
        if (_value > 0) {
            IERC20(l1TONAddress).safeTransferFrom(msg.sender, address(this), _value);
        }
        PORTAL.depositTransaction(_to, _value, _gasLimit, false, _data);
    }

    /// @inheritdoc CrossDomainMessenger
    function _isOtherMessenger() internal view override returns (bool) {
        return msg.sender == address(PORTAL) && PORTAL.l2Sender() == OTHER_MESSENGER;
    }

    /// @inheritdoc CrossDomainMessenger
    function _isUnsafeTarget(address _target) internal view override returns (bool) {
        return _target == address(this) || _target == address(PORTAL);
    }

    /// @notice Sends a deposit ton message to some target address on the other chain. Note that if the call
    ///         always reverts, then the message will be unrelayable, and any ETH sent will be
    ///         permanently locked. The same will occur if the target on the other chain is
    ///         considered unsafe (see the _isUnsafeTarget() function).
    /// @param _target      Target contract or wallet address.
    /// @param _amount      Amount of deposit TON.
    /// @param _message     Message to trigger the target address with.
    /// @param _minGasLimit Minimum gas limit that the message can be executed with.
    function sendTONMessage(address _target, uint256 _amount, bytes calldata _message, uint32 _minGasLimit) external {
        // Triggers a message to the other messenger. Note that the amount of gas provided to the
        // message is the amount of gas requested by the user PLUS the base gas value. We want to
        // guarantee the property that the call to the target contract will always have at least
        // the minimum gas limit specified by the user.
        _sendMessage(
            OTHER_MESSENGER,
            baseGas(_message, _minGasLimit),
            _amount,
            abi.encodeWithSelector(
                this.relayMessage.selector, messageNonce(), msg.sender, _target, _amount, _minGasLimit, _message
            )
        );

        emit SentMessage(_target, msg.sender, _message, messageNonce(), _minGasLimit);
        emit SentMessageExtension1(msg.sender, _amount);

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
    /// @param _value       TON value to send with the message.
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

        if (_isOtherMessenger()) {
            // These properties should always hold when the message is first submitted (as
            // opposed to being replayed).
            assert(msg.value == 0);
            assert(!failedMessages[versionedHash]);
        } else {
            require(msg.value == 0, "CrossDomainMessenger: value must be zero unless message is from a system address");

            require(failedMessages[versionedHash], "CrossDomainMessenger: message cannot be replayed");
        }

        require(
            _isUnsafeTarget(_target) == false, "CrossDomainMessenger: cannot send message to blocked system address"
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
        bool transferStatus = true;
        if (_value != 0){
            transferStatus = IERC20(l1TONAddress).transfer(_target, _value);
        }
        bool success = SafeCall.call(_target, gasleft() - RELAY_RESERVED_GAS, 0, _message);
        xDomainMsgSender = Constants.DEFAULT_L2_SENDER;

        if (success && transferStatus) {
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
