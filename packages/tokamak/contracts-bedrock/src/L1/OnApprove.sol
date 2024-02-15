// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { IERC165 } from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

abstract contract OnApprove is IERC165 {
    function supportsInterface(bytes4 interfaceId) external view returns (bool) {
        return interfaceId == (OnApprove(this).onApprove.selector);
    }

    function pack(uint32 _minGasLimit, bytes memory _message) public pure returns(bytes memory) {
        bytes memory packedData = abi.encode(_minGasLimit, _message);
        return packedData;
    }

    function unpackOnApproveData(bytes memory _data) public pure returns(uint32, bytes memory) {
        require(_data.length >= 4, "");
        uint32 _minGasLimit;
        bytes memory _message =_data;
        uint256 _messageLength;

        assembly {
            let _pos := _data
            _messageLength := sub(mload(_pos), 32)
            _pos := add(_pos, 32)
            _minGasLimit := shr(224, mload(_pos))
            _message := add(_pos, 4)
        }
        return (_minGasLimit, _message);
    }

    function onApprove(
        address _owner,
        address _spender,
        uint256 _amount,
        bytes memory _data
    )
        external
        virtual
        returns (bool);
}
