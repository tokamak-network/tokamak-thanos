// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { IERC165 } from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { IOptimismMintableERC20 } from "src/universal/IOptimismMintableERC20.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { StandardBridge } from "src/universal/StandardBridge.sol";
import { ISemver } from "src/universal/ISemver.sol";
import { CrossDomainMessenger } from "src/universal/CrossDomainMessenger.sol";
import { L1CrossDomainMessenger } from "src/L1/L1CrossDomainMessenger.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { Constants } from "src/libraries/Constants.sol";
import { SafeCall } from "src/libraries/SafeCall.sol";
import { OptimismMintableERC20 } from "src/universal/OptimismMintableERC20.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { OnApprove } from "./OnApprove.sol";

/// @custom:proxied
/// @title L1StandardBridge
/// @notice The L1StandardBridge is responsible for transfering ETH and ERC20 tokens between L1 and
///         L2. In the case that an ERC20 token is native to L1, it will be escrowed within this
///         contract. If the ERC20 token is native to L2, it will be burnt. Before Bedrock, ETH was
///         stored within this contract. After Bedrock, ETH is instead stored inside the
///         OptimismPortal contract.
///         NOTE: this contract is not intended to support all variations of ERC20 tokens. Examples
///         of some token types that may not be properly supported by this contract include, but are
///         not limited to: tokens with transfer fees, rebasing tokens, and tokens with blocklists.
contract L1StandardBridge is StandardBridge, OnApprove, ISemver {
    using SafeERC20 for IERC20;

    // address public nativeTokenAddress;

    SystemConfig public systemConfig;

    /// @custom:legacy
    /// @notice Emitted whenever a deposit of ETH from L1 into L2 is initiated.
    /// @param from      Address of the depositor.
    /// @param to        Address of the recipient on L2.
    /// @param amount    Amount of ETH deposited.
    /// @param extraData Extra data attached to the deposit.
    event ETHDepositInitiated(address indexed from, address indexed to, uint256 amount, bytes extraData);

    /// @custom:legacy
    /// @notice Emitted whenever a withdrawal of ETH from L2 to L1 is finalized.
    /// @param from      Address of the withdrawer.
    /// @param to        Address of the recipient on L1.
    /// @param amount    Amount of ETH withdrawn.
    /// @param extraData Extra data attached to the withdrawal.
    event ETHWithdrawalFinalized(address indexed from, address indexed to, uint256 amount, bytes extraData);

    /// @custom:legacy
    /// @notice Emitted whenever an ERC20 deposit is initiated.
    /// @param l1Token   Address of the token on L1.
    /// @param l2Token   Address of the corresponding token on L2.
    /// @param from      Address of the depositor.
    /// @param to        Address of the recipient on L2.
    /// @param amount    Amount of the ERC20 deposited.
    /// @param extraData Extra data attached to the deposit.
    event ERC20DepositInitiated(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /// @custom:legacy
    /// @notice Emitted whenever an ERC20 withdrawal is finalized.
    /// @param l1Token   Address of the token on L1.
    /// @param l2Token   Address of the corresponding token on L2.
    /// @param from      Address of the withdrawer.
    /// @param to        Address of the recipient on L1.
    /// @param amount    Amount of the ERC20 withdrawn.
    /// @param extraData Extra data attached to the withdrawal.
    event ERC20WithdrawalFinalized(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /// @notice Semantic version.
    /// @custom:semver 2.1.0
    string public constant version = "2.1.0";

    /// @notice Address of the SuperchainConfig contract.
    SuperchainConfig public superchainConfig;

    /// @notice Constructs the L1StandardBridge contract.
    constructor() StandardBridge() {
        initialize({
            _messenger: CrossDomainMessenger(address(0)),
            _superchainConfig: SuperchainConfig(address(0)),
            _systemConfig: SystemConfig(address(0))
        });
    }

    /// @notice Initializer
    /// @param _messenger        Contract for the CrossDomainMessenger on this network.
    /// @param _systemConfig     Contract for SystemConfig on L1
    /// @param _superchainConfig Contract for the SuperchainConfig on this network.

    function initialize(
        CrossDomainMessenger _messenger,
        SystemConfig _systemConfig,
        SuperchainConfig _superchainConfig
    )
        public
        reinitializer(Constants.INITIALIZER)
    {
        superchainConfig = _superchainConfig;
        systemConfig = _systemConfig;
        __StandardBridge_init({
            _messenger: _messenger,
            _otherBridge: StandardBridge(payable(Predeploys.L2_STANDARD_BRIDGE))
        });
    }

    /// @inheritdoc StandardBridge
    function paused() public view override returns (bool) {
        return superchainConfig.paused();
    }

    /// @notice Deposit ETH on L1 and receive ETH on L2.
    receive() external payable override onlyEOA {
        _initiateETHDeposit(msg.sender, msg.sender, RECEIVE_DEFAULT_GAS_LIMIT, bytes(""));
    }

    /// @notice unpack onApprove data
    /// @param _data     Data used in OnApprove contract
    function unpackOnApproveData(bytes calldata _data)
        public
        pure
        returns (address _from, address _to, uint256 _amount, uint32 _minGasLimit, bytes calldata _message)
    {
        require(_data.length >= 76, "Invalid onApprove data for L1StandardBridge");
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
        address _nativeTokenAddress = nativeTokenAddress();
        require(msg.sender == _nativeTokenAddress, "only accept native token approve callback");
        (address from, address to, uint256 amount, uint32 minGasLimit, bytes memory message) =
            unpackOnApproveData(_data);
        require(_owner == from && _amount == amount && amount > 0, "invalid onApprove data");
        _initiateBridgeNativeToken(from, to, amount, minGasLimit, message);
        return true;
    }

    function nativeTokenAddress() public view returns (address) {
        return systemConfig.nativeTokenAddress();
    }

    /// @notice Deposits some amount of token into the sender's native's account on L2.
    /// @param _amount      Amount of native being bridged.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function depositNativeToken(uint256 _amount, uint32 _minGasLimit, bytes calldata _extraData) external onlyEOA {
        _initiateBridgeNativeToken(msg.sender, msg.sender, _amount, _minGasLimit, _extraData);
    }

    /// @notice Deposits some amount of token into the sender's native's account on L2.
    /// @param _to          Address of the recipient on L2.
    /// @param _amount      Amount of native being bridged.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function depositNativeTokenTo(
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    )
        external
    {
        _initiateBridgeNativeToken(msg.sender, _to, _amount, _minGasLimit, _extraData);
    }

    /// @notice Deposits some amount of token into the sender's native's account on L2.
    /// @param _amount      Amount of native being bridged.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function bridgeNativeToken(uint256 _amount, uint32 _minGasLimit, bytes calldata _extraData) external onlyEOA {
        _initiateBridgeNativeToken(msg.sender, msg.sender, _amount, _minGasLimit, _extraData);
    }

    /// @notice Deposits some amount of ERC20 token into a target Native's account on L2.
    /// @param _to          Address of the recipient on L2.
    /// @param _amount      Amount of native being bridged.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function bridgeNativeTokenTo(
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    )
        external
    {
        _initiateBridgeNativeToken(msg.sender, _to, _amount, _minGasLimit, _extraData);
    }

    /// @notice Deposits some amount of ETH into the sender's ETH's account on L2.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function depositETH(uint32 _minGasLimit, bytes calldata _extraData) external payable onlyEOA {
        _initiateETHDeposit(msg.sender, msg.sender, _minGasLimit, _extraData);
    }

    /// @notice Deposits some amount of ETH into a target ETH's account on L2.
    ///         Note that if ETH is sent to a contract on L2 and the call fails, then that ETH will
    ///         be locked in the L2StandardBridge. ETH may be recoverable if the call can be
    ///         successfully replayed by increasing the amount of gas supplied to the call. If the
    ///         call will fail for any amount of gas, then the ETH will be locked permanently.
    /// @param _to          Address of the recipient on L2.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function depositETHTo(address _to, uint32 _minGasLimit, bytes calldata _extraData) external payable {
        _initiateETHDeposit(msg.sender, _to, _minGasLimit, _extraData);
    }

    /// @notice Initiates a bridge of ETH through the CrossDomainMessenger. Receive ETH on L2
    /// @param _from        Address of the sender.
    /// @param _to          Address of the receiver.
    /// @param _amount      Amount of ETH being bridged.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function _initiateBridgeETH(
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    )
        internal
        override
    {
        require(msg.value == _amount, "StandardBridge: bridging ETH must include sufficient ETH value");
        deposits[address(0)][Predeploys.ETH] = deposits[address(0)][Predeploys.ETH] + _amount;
        _emitETHBridgeInitiated(_from, _to, _amount, _extraData);

        messenger.sendMessage(
            address(otherBridge),
            abi.encodeWithSelector(
                this.finalizeBridgeERC20.selector,
                // Because this call will be executed on the remote chain, we reverse the order of
                // the remote and local token addresses relative to their order in the
                // finalizeBridgeERC20 function.
                Predeploys.ETH,
                address(0),
                _from,
                _to,
                _amount,
                _extraData
            ),
            _minGasLimit
        );
    }

    /// @custom:legacy
    /// @notice Deposits some amount of ERC20 tokens into the sender's account on L2.
    /// @param _l1Token     Address of the L1 token being deposited.
    /// @param _l2Token     Address of the corresponding token on L2.
    /// @param _amount      Amount of the ERC20 to deposit.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function depositERC20(
        address _l1Token,
        address _l2Token,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    )
        external
        virtual
        onlyEOA
    {
        _initiateERC20Deposit(_l1Token, _l2Token, msg.sender, msg.sender, _amount, _minGasLimit, _extraData);
    }

    /// @custom:legacy
    /// @notice Deposits some amount of ERC20 tokens into a target account on L2.
    /// @param _l1Token     Address of the L1 token being deposited.
    /// @param _l2Token     Address of the corresponding token on L2.
    /// @param _to          Address of the recipient on L2.
    /// @param _amount      Amount of the ERC20 to deposit.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    ///                     Data supplied here will not be used to execute any code on L2 and is
    ///                     only emitted as extra data for the convenience of off-chain tooling.
    function depositERC20To(
        address _l1Token,
        address _l2Token,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    )
        external
        virtual
    {
        _initiateERC20Deposit(_l1Token, _l2Token, msg.sender, _to, _amount, _minGasLimit, _extraData);
    }

    /// @custom:legacy
    /// @notice Finalizes a withdrawal of ETH from L2 => receive L2's native token in L1
    /// @param _from      Address of the withdrawer on L2.
    /// @param _to        Address of the recipient on L1.
    /// @param _amount    Amount of ETH to withdraw.
    /// @param _extraData Optional data forwarded from L2.
    function finalizeETHWithdrawal(
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    )
        external
        payable
    {
        finalizeBridgeETH(_from, _to, _amount, _extraData);
    }

    /// @custom:legacy
    /// @notice Finalizes a withdrawal of ERC20 tokens from L2.
    /// @param _l1Token   Address of the token on L1. L1 must be different than l1TokenAddress
    /// @param _l2Token   Address of the corresponding token on L2.
    /// @param _from      Address of the withdrawer on L2.
    /// @param _to        Address of the recipient on L1.
    /// @param _amount    Amount of the ERC20 to withdraw.
    /// @param _extraData Optional data forwarded from L2.
    function finalizeERC20Withdrawal(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    )
        external
    {
        // TODO
        finalizeBridgeERC20(_l1Token, _l2Token, _from, _to, _amount, _extraData);
    }

    /// @custom:legacy
    /// @notice Retrieves the access of the corresponding L2 bridge contract.
    /// @return Address of the corresponding L2 bridge contract.
    function l2TokenBridge() external view returns (address) {
        return address(otherBridge);
    }

    /// @notice Internal function for initiating an ETH deposit.
    /// @param _from        Address of the sender on L1.
    /// @param _to          Address of the recipient on L2.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    function _initiateETHDeposit(address _from, address _to, uint32 _minGasLimit, bytes memory _extraData) internal {
        _initiateBridgeETH(_from, _to, msg.value, _minGasLimit, _extraData);
    }

    /// @notice Internal function for initiating an ERC20 deposit.
    /// @param _l1Token     Address of the L1 token being deposited.
    /// @param _l2Token     Address of the corresponding token on L2.
    /// @param _from        Address of the sender on L1.
    /// @param _to          Address of the recipient on L2.
    /// @param _amount      Amount of the ERC20 to deposit.
    /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
    /// @param _extraData   Optional data to forward to L2.
    function _initiateERC20Deposit(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    )
        internal
    {
        _initiateBridgeERC20(_l1Token, _l2Token, _from, _to, _amount, _minGasLimit, _extraData);
    }

    /// @notice Sends ERC20 tokens to a receiver's address on the other chain.
    /// @param _localToken  Address of the ERC20 on this chain.
    /// @param _remoteToken Address of the corresponding token on the remote chain.
    /// @param _to          Address of the receiver.
    /// @param _amount      Amount of local tokens to deposit.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function _initiateBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    )
        internal
        override
    {
        if (nativeTokenAddress() == _localToken) {
            _initiateBridgeNativeToken(_from, _to, _amount, _minGasLimit, _extraData);
        } else {
            super._initiateBridgeERC20(_localToken, _remoteToken, _from, _to, _amount, _minGasLimit, _extraData);
        }
    }

    /// @notice Sends native tokens to a receiver's address on the other chain.
    /// @param _to          Address of the receiver.
    /// @param _amount      Amount of local tokens to deposit.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function _initiateBridgeNativeToken(
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    )
        internal
    {
        address _nativeTokenAddress = nativeTokenAddress();
        IERC20(_nativeTokenAddress).safeTransferFrom(_from, address(this), _amount);
        IERC20(_nativeTokenAddress).approve(address(messenger), _amount);

        _emitERC20BridgeInitiated(
            _nativeTokenAddress, Predeploys.LEGACY_ERC20_NATIVE_TOKEN, _from, _to, _amount, _extraData
        );

        L1CrossDomainMessenger(address(messenger)).sendNativeTokenMessage(
            address(otherBridge),
            _amount,
            abi.encodeWithSelector(this.finalizeBridgeETH.selector, _from, _to, _amount, _extraData),
            _minGasLimit
        );
    }

    /// @inheritdoc StandardBridge
    /// @notice Emits the legacy ETHDepositInitiated event followed by the ETHBridgeInitiated event.
    ///         This is necessary for backwards compatibility with the legacy bridge.
    function _emitETHBridgeInitiated(
        address _from,
        address _to,
        uint256 _amount,
        bytes memory _extraData
    )
        internal
        override
    {
        emit ETHDepositInitiated(_from, _to, _amount, _extraData);
        super._emitETHBridgeInitiated(_from, _to, _amount, _extraData);
    }

    /// @inheritdoc StandardBridge
    /// @notice Emits the legacy ERC20DepositInitiated event followed by the ERC20BridgeInitiated
    ///         event. This is necessary for backwards compatibility with the legacy bridge.
    function _emitETHBridgeFinalized(
        address _from,
        address _to,
        uint256 _amount,
        bytes memory _extraData
    )
        internal
        override
    {
        emit ETHWithdrawalFinalized(_from, _to, _amount, _extraData);
        super._emitETHBridgeFinalized(_from, _to, _amount, _extraData);
    }

    /// @inheritdoc StandardBridge
    /// @notice Emits the legacy ERC20WithdrawalFinalized event followed by the ERC20BridgeFinalized
    ///         event. This is necessary for backwards compatibility with the legacy bridge.
    function _emitERC20BridgeInitiated(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        bytes memory _extraData
    )
        internal
        override
    {
        emit ERC20DepositInitiated(_localToken, _remoteToken, _from, _to, _amount, _extraData);
        super._emitERC20BridgeInitiated(_localToken, _remoteToken, _from, _to, _amount, _extraData);
    }

    /// @inheritdoc StandardBridge
    /// @notice Emits the legacy ERC20WithdrawalFinalized event followed by the ERC20BridgeFinalized
    ///         event. This is necessary for backwards compatibility with the legacy bridge.
    function _emitERC20BridgeFinalized(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        bytes memory _extraData
    )
        internal
        override
    {
        emit ERC20WithdrawalFinalized(_localToken, _remoteToken, _from, _to, _amount, _extraData);
        super._emitERC20BridgeFinalized(_localToken, _remoteToken, _from, _to, _amount, _extraData);
    }

    /// @notice Finalizes an ETH bridge on this chain. Can only be triggered by the other
    ///         StandardBridge contract on the remote chain.
    /// @param _from      Address of the sender.
    /// @param _to        Address of the receiver.
    /// @param _amount    Amount of ETH being bridged.
    /// @param _extraData Extra data to be sent with the transaction. Note that the recipient will
    ///                   not be triggered with this data, but it will be emitted and can be used
    ///                   to identify the transaction.
    function finalizeBridgeETH(
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    )
        public
        payable
        override
        onlyOtherBridge
    {
        require(paused() == false, "L1 StandardBridge: paused");
        address _nativeTokenAddress = nativeTokenAddress();
        require(msg.value == 0, "StandardBridge: must not receive Ether even if this is payable");
        require(_to != address(this), "StandardBridge: cannot send to self");
        require(_to != address(messenger), "StandardBridge: cannot send to messenger");

        IERC20(_nativeTokenAddress).safeTransferFrom(address(messenger), address(this), _amount);

        // Emit the correct events. By default this will be _amount, but child
        // contracts may override this function in order to emit legacy events as well.
        finalizeBridgeERC20(_nativeTokenAddress, Predeploys.LEGACY_ERC20_NATIVE_TOKEN, _from, _to, _amount, _extraData);
    }

    function finalizeBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    )
        public
        override
        onlyOtherBridge
    {
        require(paused() == false, "L1 StandardBridge: paused");

        address _nativeTokenAddress = nativeTokenAddress();
        if (_isOptimismMintableERC20(_localToken) && _localToken != _nativeTokenAddress) {
            require(
                _isCorrectTokenPair(_localToken, _remoteToken),
                "StandardBridge: wrong remote token for Optimism Mintable ERC20 local token"
            );
            OptimismMintableERC20(_localToken).mint(_to, _amount);
        } else {
            // Case 1: Withdraw ETH
            // Case 2: Withdraw normal Token
            if (_localToken == address(0) && _remoteToken == Predeploys.ETH) {
                bool success = SafeCall.call(_to, gasleft(), _amount, hex"");
                require(success, "StandardBridge: ETH transfer failed");
                _emitETHBridgeFinalized(_from, _to, _amount, _extraData);
            } else {
                IERC20(_localToken).safeTransfer(_to, _amount);

                // Emit the correct events. By default this will be ERC20BridgeFinalized, but child
                // contracts may override this function in order to emit legacy events as well.
                _emitERC20BridgeFinalized(_localToken, _remoteToken, _from, _to, _amount, _extraData);
            }
            if (_nativeTokenAddress != _localToken) {
                deposits[_localToken][_remoteToken] = deposits[_localToken][_remoteToken] - _amount;
            }
        }
    }
}
