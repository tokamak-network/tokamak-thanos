// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/**
 * @title OptimismPortalClosing
 * @notice Shutdown router that blocks deposits while delegating all other logic to the
 *         audited OptimismPortal implementation.
 */
contract OptimismPortalClosing {
  /// @notice Address of the audited OptimismPortal implementation.
  address public immutable implementation;

  /// @param _implementation Address of the audited OptimismPortal implementation.
  constructor(address _implementation) {
    require(_implementation != address(0), 'OptimismPortalClosing: implementation is zero');
    implementation = _implementation;
  }

  /// @notice Blocks direct deposit calls.
  function depositTransaction(
    address, // _to
    uint256, // _mint
    uint256, // _value
    uint64, // _gasLimit
    bool, // _isCreation
    bytes calldata // _data
  ) external {
    revert('OptimismPortal: deposits are disabled due to chain shutdown');
  }

  /// @notice Blocks approve-and-call deposit path.
  function onApprove(
    address, // _owner
    address, // _spender
    uint256, // _amount
    bytes calldata // _data
  ) external returns (bool) {
    revert('OptimismPortal: deposits are disabled due to chain shutdown');
  }

  /// @notice Blocks direct ETH transfers.
  receive() external payable {
    revert('OptimismPortal: deposits are disabled due to chain shutdown');
  }

  /// @notice Returns a modified version string to indicate this is a closing implementation.
  function version() public pure returns (string memory) {
    return '2.8.1-closing';
  }

  fallback() external payable {
    _delegate(implementation);
  }

  function _delegate(address _implementation) internal {
    assembly {
      calldatacopy(0, 0, calldatasize())
      let result := delegatecall(gas(), _implementation, 0, calldatasize(), 0, 0)
      returndatacopy(0, 0, returndatasize())
      switch result
      case 0 {
        revert(0, returndatasize())
      }
      default {
        return(0, returndatasize())
      }
    }
  }
}
