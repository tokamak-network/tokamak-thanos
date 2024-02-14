// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { IERC165 } from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

abstract contract OnApprove is IERC165 {
    function supportsInterface(bytes4 interfaceId) external view returns (bool) {
        return interfaceId == (OnApprove(this).onApprove.selector);
    }

    function onApprove(
        address _owner,
        address _spender,
        uint256 _amount,
        bytes calldata data
    )
        external
        virtual
        returns (bool);
}
