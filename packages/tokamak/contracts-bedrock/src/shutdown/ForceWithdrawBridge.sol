// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {L1StandardBridge} from 'src/L1/L1StandardBridge.sol';
import {IERC20} from '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import {
  SafeERC20
} from '@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol';
import {
  ReentrancyGuard
} from '@openzeppelin/contracts/security/ReentrancyGuard.sol';

interface IProxyAdminOwner {
  function owner() external view returns (address);
}

/// @title ForceWithdrawBridge
/// @notice Extension of L1StandardBridge that adds force withdrawal functionality for L2 shutdown scenarios.
///         This contract enables users to claim their L2 assets directly from the L1 bridge when the L2 network
///         is no longer operational. It verifies claims against on-chain hash storage (GenFWStorage contracts)
///         to ensure only legitimate asset owners can withdraw their funds.
/// @dev Based on Titan Legacy UpgradeL1BridgeV1 but adapted for Optimism Bedrock architecture.
///      Key differences from Titan Legacy:
///      - Uses Bedrock's L1StandardBridge as base (instead of legacy version)
///      - Compatible with L1ChugSplashProxy (setCode upgrade pattern)
///      - Simplified modifiers using assembly for proxy storage slots
///      - Added getProxyOwner() for compatibility with upgrade tasks
///      - Added ReentrancyGuard for protection against reentrancy attacks
contract ForceWithdrawBridge is L1StandardBridge, ReentrancyGuard {
  using SafeERC20 for IERC20;

  /// @notice Error thrown when caller is not the proxy owner
  error FW_ONLY_OWNER();
  /// @notice Error thrown when proxy admin owner lookup fails
  error FW_OWNER_LOOKUP_FAILED();

  /// @notice Error thrown when caller is not the designated closer
  error FW_ONLY_CLOSER();

  /// @notice Error thrown when the position contract is not registered or inactive
  error FW_NOT_AVAILABLE_POSITION();

  /// @notice Error thrown when hash lookup fails in position contract
  error FW_NOT_SEARCH_POSITION();

  /// @notice Error thrown when the provided hash doesn't match the computed hash
  error FW_INVALID_HASH();

  /// @notice Error thrown when ETH transfer fails during force withdrawal
  error FW_FAIL_TRANSFER_ETH();

  /// @notice Error thrown when trying to set the same storage value
  error ER_SAME_STORAGE();

  /// @notice Emitted when a force withdrawal claim is successfully processed
  /// @param _index Hash of (token, claimer, amount) - unique identifier for the claim
  /// @param _token L1 token address (address(0) for ETH)
  /// @param amount Amount of tokens withdrawn
  /// @param _claimer Address receiving the tokens
  /// @param _requester Address that initiated the transaction (can be different from claimer)
  event ForceWithdraw(
    bytes32 indexed _index,
    address indexed _token,
    uint256 amount,
    address _claimer,
    address _requester
  );

  /// @notice EIP-1967 implementation slot
  /// @dev keccak256("eip1967.proxy.implementation") - 1
  bytes32 internal constant IMPLEMENTATION_KEY =
    0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

  /// @notice EIP-1967 admin slot
  /// @dev keccak256("eip1967.proxy.admin") - 1
  bytes32 internal constant OWNER_KEY =
    0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

  /// @notice Parameter structure for batch force withdrawal claims
  /// @param position Contract address where the hash value is stored (GenFWStorage contract)
  /// @param hashed String representation of the hash (without 0x prefix)
  /// @param token L1 token address to receive (address(0) for ETH)
  /// @param amount Amount of tokens to receive
  /// @param getAddress Address to receive the tokens
  struct ForceClaimParam {
    address position;
    string hashed;
    address token;
    uint256 amount;
    address getAddress;
  }

  /// @notice Parameter structure for modifying position registry state
  /// @param position Address of the GenFWStorage contract
  /// @param state New active state (true = active, false = inactive)
  struct ForceRegistryParam {
    address position;
    bool state;
  }

  /// @notice Maps hash => claimer address (legacy, kept for compatibility but not actively used)
  /// @dev This was used in Titan Legacy but claim validation now uses claimState mapping
  mapping(bytes32 => address) public gb;

  /// @notice Maps GenFWStorage contract addresses to their active status
  /// @dev Only active positions (true) can be used for force withdrawal claims
  mapping(address => bool) public position;

  /// @notice Array of all registered GenFWStorage contract addresses
  /// @dev Used for iterating positions in getForcePosition()
  address[] public positions;

  /// @notice Address authorized to manage force withdrawal settings
  /// @dev Can activate/deactivate force withdrawals and register positions
  ///      Separate from owner to allow delegation of operational control
  address public closer;

  /// @notice Signature for calling getOwner() on the proxy
  /// @dev Used in onlyOwner modifier to retrieve proxy admin via delegatecall
  bytes constant SIG_GETOWNER = abi.encodeWithSignature('getOwner()');

  /// @notice Tracks which hashes have already been claimed
  /// @dev Prevents double-claiming the same asset
  ///      Key: keccak256(abi.encodePacked(token, claimer, amount))
  mapping(bytes32 => bool) public claimState;

  /// @notice Controls whether force withdrawal functionality is active
  /// @dev When false, all force withdrawal claims will fail
  ///      Toggled by closer via forceActive()
  bool public active;

  /// @notice Modifier restricting access to the designated closer address
  modifier onlyCloser() virtual {
    if (msg.sender != closer) revert FW_ONLY_CLOSER();
    _;
  }

  /// @notice Modifier restricting access to the proxy owner
  /// @dev Tries getOwner() for legacy proxies; falls back to ERC1967 admin slot.
  modifier onlyOwner() virtual {
    address owner = _getProxyOwner();
    if (msg.sender == owner) {
      _;
      return;
    }

    if (_isContract(owner)) {
      try IProxyAdminOwner(owner).owner() returns (address adminOwner) {
        if (msg.sender == adminOwner) {
          _;
          return;
        }
      } catch {
        revert FW_OWNER_LOOKUP_FAILED();
      }
    }

    revert FW_ONLY_OWNER();
  }

  function _getProxyOwner() internal returns (address owner_) {
    bytes32 slot = OWNER_KEY;
    assembly {
      owner_ := sload(slot)
    }
    if (owner_ != address(0)) {
      return owner_;
    }

    (bool success, bytes memory data) = address(this).delegatecall(SIG_GETOWNER);
    if (success && data.length >= 32) {
      owner_ = abi.decode(data, (address));
    }
  }

  function _isContract(address _addr) internal view returns (bool) {
    return _addr.code.length > 0;
  }

  /// @notice Activates or deactivates the force withdrawal functionality
  /// @dev Only callable by the closer address
  /// @param _state New active state (true to enable, false to disable)
  function forceActive(bool _state) external onlyCloser {
    require(active != _state, 'Same State');
    active = _state;
  }

  /// @notice Sets the closer address
  /// @dev Only callable by the proxy owner
  /// @param _closer New closer address
  function setCloser(address _closer) external onlyOwner {
    require(closer != _closer, 'Same Address');
    closer = _closer;
  }

  /// @notice Sets both closer and active state in a single transaction
  /// @dev Only callable by the proxy owner. Reverts if both values are already set to the provided values.
  /// @param _closer New closer address
  /// @param _state New active state
  function setCloserAndActive(address _closer, bool _state) external onlyOwner {
    if (closer == _closer && active == _state) {
      revert ER_SAME_STORAGE();
    }
    closer = _closer;
    active = _state;
  }

  /// @notice Registers GenFWStorage contract addresses as valid position contracts
  /// @dev Only callable by closer. All provided addresses are set to active (true).
  ///      GenFWStorage contracts contain hash constants used to verify force withdrawal claims.
  /// @param _position Array of GenFWStorage contract addresses to register
  function forceRegistry(address[] calldata _position) external onlyCloser {
    for (uint256 i = 0; i < _position.length; i++) {
      position[_position[i]] = true;
      positions.push(_position[i]);
    }
  }

  /// @notice Modifies the active state of registered position contracts
  /// @dev Only callable by closer. Allows activating/deactivating positions without removing them.
  /// @param _data Array of position addresses and their new states
  function forceModify(
    ForceRegistryParam[] calldata _data
  ) external onlyCloser {
    for (uint256 i = 0; i < _data.length; i++) {
      position[_data[i].position] = _data[i].state;
    }
  }

  /// @notice Finds which registered position contract contains a given hash
  /// @dev Iterates through all active positions and calls the hash function via staticcall
  /// @param _hash Hash value to search for (without 0x prefix)
  /// @return Address of the position contract containing the hash, or address(0) if not found
  function getForcePosition(
    string memory _hash
  ) external view returns (address) {
    string memory f = string(abi.encodePacked('_', _hash, '()'));
    for (uint256 i = 0; i < positions.length; i++) {
      address p = positions[i];

      if (position[p] == false) continue;

      (bool success, bytes memory data) = p.staticcall(
        abi.encodeWithSignature(f)
      );

      if (success) {
        bytes32 r = abi.decode(data, (bytes32));
        if (r == 0) {
          continue;
        }
        return p;
      }
    }
    return address(0);
  }

  /// @notice Processes multiple force withdrawal claims in a single transaction
  /// @dev Iterates through all params and calls internal claim() function for each
  /// @param params Array of claim parameters (position, hash, token, amount, recipient)
  function forceWithdrawClaimAll(ForceClaimParam[] calldata params) external nonReentrant {
    for (uint256 i = 0; i < params.length; i++) {
      claim(
        params[i].position,
        params[i].hashed,
        params[i].token,
        params[i].amount,
        params[i].getAddress
      );
    }
  }

  /// @notice Processes a single force withdrawal claim
  /// @dev Public wrapper around internal claim() function
  /// @param _position GenFWStorage contract address containing the hash
  /// @param _hash Hash value (without 0x prefix) to verify against
  /// @param _token L1 token address (address(0) for ETH)
  /// @param _amount Amount to withdraw
  /// @param _address Recipient address
  function forceWithdrawClaim(
    address _position,
    string calldata _hash,
    address _token,
    uint256 _amount,
    address _address
  ) external nonReentrant {
    claim(_position, _hash, _token, _amount, _address);
  }

  /// @notice Internal function that validates and processes a force withdrawal claim
  /// @dev Verification steps:
  ///      1. Check position is registered and active
  ///      2. Look up hash constant in GenFWStorage via staticcall
  ///      3. Compute expected hash: keccak256(abi.encodePacked(token, claimer, amount))
  ///      4. Verify computed hash matches stored hash
  ///      5. Check hash hasn't been claimed before
  ///      6. Mark hash as claimed
  ///      7. Transfer tokens (ETH via call, ERC20 via safeTransfer)
  /// @param _position GenFWStorage contract address
  /// @param _hash Hash value (without 0x prefix)
  /// @param _token Token address (address(0) for ETH)
  /// @param _amount Amount to transfer
  /// @param _address Recipient address
  function claim(
    address _position,
    string calldata _hash,
    address _token,
    uint256 _amount,
    address _address
  ) internal {
    if (!position[_position]) revert FW_NOT_AVAILABLE_POSITION();

    // Construct function signature: _<hash>()
    string memory f = string(abi.encodePacked('_', _hash, '()'));
    (bool s, bytes memory d) = _position.staticcall(abi.encodeWithSignature(f));

    if (!s || d.length == 0) {
      revert FW_NOT_SEARCH_POSITION();
    }

    // Compute expected hash
    bytes32 v = keccak256(abi.encodePacked(_token, _address, _amount));
    bytes32 r = abi.decode(d, (bytes32));

    require(claimState[r] == false, 'already claim Hash');

    if (v != r) {
      revert FW_INVALID_HASH();
    }

    claimState[r] = true;

    // Transfer tokens
    if (_token == address(0)) {
      (s, ) = (_address).call{value: _amount}(new bytes(0));
      if (!s) revert FW_FAIL_TRANSFER_ETH();
    } else {
      IERC20(_token).safeTransfer(_address, _amount);
    }

    emit ForceWithdraw(r, _token, _amount, _address, msg.sender);
  }

  /// @notice Returns the proxy owner address from EIP-1967 storage slot
  /// @dev Used by tasks to verify ownership before performing upgrades
  /// @return owner Proxy admin address
  function getProxyOwner() external view returns (address owner) {
    assembly {
      owner := sload(OWNER_KEY)
    }
  }

  /// @notice Returns the proxy implementation address from EIP-1967 storage slot
  /// @dev Useful for verifying successful upgrades
  /// @return implementation Current implementation contract address
  function getProxyImplementation()
    external
    view
    returns (address implementation)
  {
    assembly {
      implementation := sload(IMPLEMENTATION_KEY)
    }
  }
}
