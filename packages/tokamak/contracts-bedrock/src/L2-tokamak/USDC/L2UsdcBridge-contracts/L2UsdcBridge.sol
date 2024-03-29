// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import { IERC20 } from "@openzeppelin/contracts_v5.0.1/token/ERC20/IERC20.sol";
import { Address } from "../../../libraries/Address.sol";
import { SafeERC20 } from "@openzeppelin/contracts_v5.0.1/token/ERC20/utils/SafeERC20.sol";

import { ERC1967Utils } from "@openzeppelin/contracts_v5.0.1/proxy/ERC1967/ERC1967Utils.sol";
import { L2UsdcBridgeStorage } from "./L2UsdcBridgeStorage.sol";

interface ICrossDomainMessenger {
    function xDomainMessageSender() external view returns (address) ;
    function sendMessage(
        address _target,
        bytes calldata _message,
        uint32 _minGasLimit
    ) external payable;
}

interface IUSDC {
    function mint(address to, uint256 amount) external;
    function burn(uint256 amount) external;
    function minterAllowance(address minter) external view returns (uint256) ;
    function balanceOf(address account) external view returns (uint256) ;

}

interface IL1USDCBridge {
     function finalizeERC20Withdrawal(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    ) external;
}

interface IMasterMinter {
    function configureMinter(uint256 _newAllowance) external returns (bool);

}

contract L2UsdcBridge is L2UsdcBridgeStorage {
    using SafeERC20 for IERC20;

    /**
     * @custom:legacy
     * @notice Emitted whenever a withdrawal from L2 to L1 is initiated.
     *
     * @param l1Token   Address of the token on L1.
     * @param l2Token   Address of the corresponding token on L2.
     * @param from      Address of the withdrawer.
     * @param to        Address of the recipient on L1.
     * @param amount    Amount of the ERC20 withdrawn.
     * @param extraData Extra data attached to the withdrawal.
     */
    event WithdrawalInitiated(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /**
     * @custom:legacy
     * @notice Emitted whenever an ERC20 deposit is finalized.
     *
     * @param l1Token   Address of the token on L1.
     * @param l2Token   Address of the corresponding token on L2.
     * @param from      Address of the depositor.
     * @param to        Address of the recipient on L2.
     * @param amount    Amount of the ERC20 deposited.
     * @param extraData Extra data attached to the deposit.
     */
    event DepositFinalized(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /**
     * @notice Only allow EOAs to call the functions. Note that this is not safe against contracts
     *         calling code within their constructors, but also doesn't really matter since we're
     *         just trying to prevent users accidentally depositing with smart contract wallets.
     */
    modifier onlyEOA() {
        require(
            !Address.isContract(msg.sender),
            "StandardBridge: function can only be called from an EOA"
        );
        _;
    }

    /**
     * @notice Ensures that the caller is a cross-chain message from the other bridge.
     */
    modifier onlyOtherBridge() {
        require(
            msg.sender == messenger &&
                ICrossDomainMessenger(messenger).xDomainMessageSender() == otherBridge,
            "StandardBridge: function can only be called from the other bridge"
        );
        _;
    }

    /**
     * @custom:legacy
     * @notice Initiates a withdrawal from L2 to L1.
     *
     * @param _l2Token     Address of the L2 token to withdraw.
     * @param _amount      Amount of the L2 token to withdraw.
     * @param _minGasLimit Minimum gas limit to use for the transaction.
     * @param _extraData   Extra data attached to the withdrawal.
     */
    function withdraw(
        address _l2Token,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) external payable virtual onlyEOA {
        _initiateWithdrawal(_l2Token, msg.sender, msg.sender, _amount, _minGasLimit, _extraData);
    }

    /**
     * @custom:legacy
     * @notice Finalizes a deposit from L1 to L2.
     *
     * @param _l1Token   Address of the L1 token to deposit.
     * @param _l2Token   Address of the corresponding L2 token.
     * @param _from      Address of the depositor.
     * @param _to        Address of the recipient.
     * @param _amount    Amount of the tokens being deposited.
     * @param _extraData Extra data attached to the deposit.
     */
    function finalizeDeposit(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    ) external onlyOtherBridge onlyL1Usdc(_l1Token) onlyL2Usdc(_l2Token) {

        uint256 allowance = IUSDC(_l2Token).minterAllowance(address(this));

        if (allowance < _amount) IMasterMinter(l2UsdcMasterMinter).configureMinter(type(uint256).max);

        IUSDC(_l2Token).mint(_to, _amount);

        emit DepositFinalized(_l1Token, _l2Token, _from, _to, _amount, _extraData);
    }

    /**
     * @custom:legacy
     * @notice Internal function to a withdrawal from L2 to L1 to a target account on L1.
     *
     * @param _l2Token     Address of the L2 token to withdraw.
     * @param _from        Address of the withdrawer.
     * @param _to          Recipient account on L1.
     * @param _amount      Amount of the L2 token to withdraw.
     * @param _minGasLimit Minimum gas limit to use for the transaction.
     * @param _extraData   Extra data attached to the withdrawal.
     */
    function _initiateWithdrawal(
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) internal onlyL2Usdc(_l2Token) {

        uint256 allowance = IUSDC(_l2Token).minterAllowance(address(this));
        if (allowance < _amount) IMasterMinter(l2UsdcMasterMinter).configureMinter(type(uint256).max);

        IERC20(_l2Token).safeTransferFrom(_from, address(this), _amount);
        IUSDC(_l2Token).burn(_amount);

        ICrossDomainMessenger(messenger).sendMessage(
            otherBridge,
            abi.encodeWithSelector(
                IL1USDCBridge.finalizeERC20Withdrawal.selector,
                // Because this call will be executed on the remote chain, we reverse the order of
                // the remote and local token addresses relative to their order in the
                // finalizeBridgeERC20 function.
                l1Usdc,
                _l2Token,
                _from,
                _to,
                _amount,
                _extraData
            ),
            _minGasLimit
        );

        emit WithdrawalInitiated(l1Usdc, _l2Token, _from, _to, _amount, _extraData);
    }


}
