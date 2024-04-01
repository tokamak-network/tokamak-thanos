// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.15;

import "../libraries/SafeERC20.sol";

import { AccessibleCommon } from "../common/AccessibleCommon.sol";
import { L2FastWithdrawStorage } from "./L2FastWithdrawStorage.sol";
import { IOptimismMintableERC20, ILegacyMintableERC20 } from "../interfaces/IOptimismMintableERC20.sol";

// import "hardhat/console.sol";

contract L2FastWithdraw is AccessibleCommon, L2FastWithdrawStorage {

    using SafeERC20 for IERC20;

    //=======modifier========

    modifier onlyEOA() {
        require(!_isContract(msg.sender), "L2FW: function can only be called from an EOA");
        _;
    }

    //=======external========

    // function registerToken(
    //     address _l1token,
    //     address _l2token
    // )
    //     external
    //     onlyOwner
    // {
    //     enteringToken[_l2token] = _l1token;
    //     checkToken[_l1token][_l2token] = true;
    // }

    // function deleteToken(
    //     address _l1token,
    //     address _l2token
    // )
    //     external
    //     onlyOwner
    // {
    //     // enteringToken[_l2token] = address(0);
    //     checkToken[_l1token][_l2token] = false;
    // }

    function requestFW(
        address _l2token,
        uint256 _totalAmount,
        uint256 _fwAmount
    )
        external
        payable
        onlyEOA
    {
        // address l1token = enteringToken[_l2token];
        // require(checkToken[l1token][_l2token], "not entering token");
        ++salecount;

        if (dealData[salecount].l2token == LEGACY_ERC20_ETH) {
            require(msg.value == _totalAmount, "FW: nativeTON need amount");
            payable(address(this)).call{value: msg.value};
        } else {
            //need to approve
            IERC20(dealData[salecount].l2token).safeTransferFrom(msg.sender,address(this),dealData[salecount].totalAmount);
        }

        dealData[salecount] = RequestData({
            l2token: _l2token,
            requester: msg.sender,
            provider: address(0),
            totalAmount: _totalAmount,
            fwAmount: _fwAmount
        });
    }

    function claimFW(
        address _from,
        address _to,
        uint256 _amount,
        uint256 _saleCount
    )
        external
        payable
    {
        // console.log("msg.sender: ", msg.sender);
        require(dealData[_saleCount].requester == _to, "not match the seller");
        require(dealData[_saleCount].fwAmount <= _amount, "need to over minAmount");
        dealData[_saleCount].provider = _from;
        dealData[_saleCount].fwAmount = _amount;

        if(dealData[_saleCount].l2token == LEGACY_ERC20_ETH) {
            (bool sent, ) = payable(_from).call{value: dealData[_saleCount].totalAmount}("");
            require(sent, "claim fail");
        } else {
            IERC20(dealData[_saleCount].l2token).safeTransferFrom(address(this),_from,dealData[_saleCount].totalAmount);
        }
    }

    function cancelFW(
        uint256 _salecount
    )
        external
        payable
    {
        require(dealData[_salecount].requester == dealData[_salecount].provider && dealData[_salecount].fwAmount == 0, "already been sold");
        require(dealData[_salecount].requester == msg.sender, "your not seller");

        if (dealData[_salecount].l2token == LEGACY_ERC20_ETH) {
            (bool sent, ) = payable(msg.sender).call{value: dealData[_salecount].totalAmount}("");
            require(sent, "cancel refund fail");
        } else {
            IERC20(dealData[_salecount].l2token).safeTransfer(msg.sender,dealData[_salecount].totalAmount);
        }
    }

    function getL1token(
        address _l2token
    )
        external
        view
        returns (address l1Token)
    {
        if (_l2token == LEGACY_ERC20_ETH) {
            return l1Token = LEGACY_l1token;
        } else {
            return l1Token = ILegacyMintableERC20(_l2token).l1Token();
        }
    }

    //=======internal========

    function _isContract(address account) internal view returns (bool) {
        return account.code.length > 0;
    }
}
