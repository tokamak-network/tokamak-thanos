// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.15;

import "../libraries/SafeERC20.sol";

import { IL2FastWithdraw } from "../interfaces/IL2FastWithdraw.sol";
import { IL1CrossDomainMessenger } from "../interfaces/IL1CrossDomainMessenger.sol";
import { AccessibleCommon } from "../common/AccessibleCommon.sol";
import { L1FastWithdrawStorage } from "./L1FastWithdrawStorage.sol";

contract L1FastWithdraw is AccessibleCommon, L1FastWithdrawStorage {

    using SafeERC20 for IERC20;

    function initialize(
        address _crossDomainMessenger,
        address _l2fastWithdraw,
        address _legacyERC20,
        address _l1legacyERC20
    )
        public
    {
        crossDomainMessenger = _crossDomainMessenger;
        l2fastWithdrawContract = _l2fastWithdraw;
        LEGACY_ERC20_ETH = _legacyERC20;
        LEGACY_l1token = _l1legacyERC20;
    }

    function provideFW(
        address _l1token,
        address _to,
        uint256 _amount,
        uint256 _saleCount,
        uint32 _minGasLimit
    )
        external
        payable
    {
        //need the check cancel or editing (L1FW)
        bytes memory message;

        message = abi.encodeWithSignature("claimFW(address,address,uint256,uint256)",
            msg.sender,
            _to,
            _amount,
            _saleCount
        );

        // message = abi.encodeWithSelector(
        //     IL2FastWithdraw.claimFW.selector,
        //     msg.sender,
        //     _to,
        //     _amount,
        //     _saleCount
        // );

        //need to approve
        if (LEGACY_l1token == _l1token) {
            IERC20(_l1token).transferFrom(msg.sender, address(this), _amount);
            IERC20(_l1token).transfer(_to,_amount);
        } else {
            IERC20(_l1token).transferFrom(msg.sender, _to, _amount);
        }


        IL1CrossDomainMessenger(crossDomainMessenger).sendMessage(
            l2fastWithdrawContract,
            message,
            _minGasLimit
        );
    }

    function cancelFW(

    )
        external
    {

    }

    function editFW(

    )
        external
    {

    }

}