// SPDX-License-Identifier: Unlicense
pragma solidity 0.8.15;

import "./L2FastWithdrawStorage";

interface IL2StandardBridge {
    function withdrawTo(
        address _l2Token,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    )
        external
        payable;
}


contract L2FastWithdraw is L2FastWithdrawStorage {
    //use Owner
    function setL2StandardBridge(
        address _l2StandardBridge
    ) external {
        L2StandradBridge = _l2StandardBridge;
    }

    function setL1FwContract(
        address _l1fwAddr
    ) external {
        L1fwContract = _l1fwAddr;
    }


    //use User
    function applicationFW(
        address _l2Token,
        uint256 _amount,
        uint256 _minAmount,
        uint32 _minGasLimit
    ) external {
        //make the extraData (msg.sender, amount)
        IL2StandardBridge(L2StandradBridge).withdrawTo(_l2Token, L1fwContract, _amount, _minGasLimit, _extraData);

        //sendMessage to L1FwContract;
    }
}
