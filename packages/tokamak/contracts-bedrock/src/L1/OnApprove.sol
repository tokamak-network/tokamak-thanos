// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

abstract contract OnApprove {
    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        return interfaceId == OnApprove.onApprove.selector || interfaceId == OnApprove.supportsInterface.selector;
    }

    function onApprove(
        address _owner,
        address _spender,
        uint256 _amount,
        bytes calldata _data
    )
        external
        virtual
        returns (bool);
}
