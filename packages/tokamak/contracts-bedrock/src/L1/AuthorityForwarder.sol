// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

/// @title AuthorityForwarder
/// @notice Intermediary contract that manages ownership transfer from L2 Operator to DAO.
///         Enables phased authority transfer: Phase 1 (full operator control) -> Phase 2 (DAO-controlled dangerous functions).
/// @dev This contract acts as the owner of critical L1 contracts (SystemConfig, ProxyAdmin, etc.)
///      and forwards calls from authorized parties while enforcing access control based on the current phase.
contract AuthorityForwarder {
    /// @notice The L2 operator address (immutable, set at deployment)
    address public immutable OPERATOR;

    /// @notice The DAO address (initially zero, set once via setDAO)
    address public DAO;

    /// @notice Mapping of delegated executors: [Target Contract][Function Selector] -> [Executor Address]
    /// @dev When set, only the delegated executor can call the specific function, excluding even the operator
    mapping(address => mapping(bytes4 => address)) public delegatedExecutors;

    /// @notice Enum representing the current phase of authority control
    enum Phase {
        Initial,        // DAO not set, operator has full control
        DAOControlled   // DAO set, dangerous functions blocked for operator
    }

    /// @notice Emitted when a call is successfully forwarded to a target contract
    /// @param target The address of the target contract
    /// @param success Whether the call succeeded
    event CallForwarded(address indexed target, bool success);

    /// @notice Emitted when the DAO address is set
    /// @param dao The address of the DAO
    event DAOSet(address indexed dao);

    /// @notice Emitted when a function is delegated to a specific executor
    /// @param target The target contract address
    /// @param selector The function selector
    /// @param executor The executor address
    event DelegationSet(address indexed target, bytes4 indexed selector, address indexed executor);

    /// @notice Constructor to initialize the AuthorityForwarder
    /// @param _operator The address of the L2 operator (typically a Safe multisig)
    constructor(address _operator) {
        require(_operator != address(0), "AuthorityForwarder: invalid operator address");
        OPERATOR = _operator;
    }

    /// @notice Forwards a call to a target contract with access control
    /// @param _target The address of the target contract
    /// @param _data The calldata to forward
    /// @dev Access control priority: 1. Delegated Executor, 2. Operator (with phase checks), 3. DAO (unrestricted)
    function forwardCall(address _target, bytes calldata _data) external payable {
        bytes4 selector = bytes4(_data[:4]);
        address delegatee = delegatedExecutors[_target][selector];

        // Priority 1: Delegated Executor (if set, has exclusive access)
        if (delegatee == msg.sender) {
            _executeCall(_target, _data);
            return;
        }

        // Priority 2: Operator (with restrictions in Phase 2)
        if (msg.sender == OPERATOR) {
            // If function is delegated, operator cannot call it
            require(delegatee == address(0), "AuthorityForwarder: function delegated exclusively");

            // In Phase 2 (DAO set), block dangerous functions for operator
            if (DAO != address(0) && _isDangerousFunction(selector)) {
                revert("AuthorityForwarder: dangerous operation blocked, use DAO governance");
            }

            _executeCall(_target, _data);
            return;
        }

        // Priority 3: DAO (unrestricted access to all functions)
        if (msg.sender == DAO) {
            _executeCall(_target, _data);
            return;
        }

        revert("AuthorityForwarder: caller not authorized");
    }

    /// @notice Sets the DAO address (one-time operation, irreversible)
    /// @param _dao The address of the DAO contract
    /// @dev Can only be called by the operator, and only once. DAO must be a contract (not EOA).
    function setDAO(address _dao) external {
        require(msg.sender == OPERATOR, "AuthorityForwarder: caller not operator");
        require(DAO == address(0), "AuthorityForwarder: DAO already set");
        require(_dao != address(0), "AuthorityForwarder: invalid DAO address");

        // Security: Ensure DAO is a contract (prevent accidental EOA setting)
        require(_dao.code.length > 0, "AuthorityForwarder: DAO must be a contract");

        DAO = _dao;
        emit DAOSet(_dao);
    }

    /// @notice Delegates execution of a specific function to a third party
    /// @param _target The target contract address
    /// @param _selector The function selector to delegate
    /// @param _executor The address that will have exclusive execution rights
    /// @dev Can only be called by the DAO. Once delegated, even the operator cannot call the function.
    function setDelegatedExecutor(address _target, bytes4 _selector, address _executor) external {
        require(msg.sender == DAO, "AuthorityForwarder: caller not DAO");
        delegatedExecutors[_target][_selector] = _executor;
        emit DelegationSet(_target, _selector, _executor);
    }

    /// @notice Returns the current phase of authority control
    /// @return The current phase (Initial or DAOControlled)
    function currentPhase() public view returns (Phase) {
        return DAO == address(0) ? Phase.Initial : Phase.DAOControlled;
    }

    /// @notice Internal function to execute a call to a target contract
    /// @param _target The target contract address
    /// @param _data The calldata to execute
    /// @dev Reverts with the original error if the call fails
    function _executeCall(address _target, bytes memory _data) internal {
        (bool success, bytes memory result) = _target.call{ value: msg.value }(_data);
        if (!success) {
            // Bubble up the revert reason
            assembly {
                revert(add(result, 32), mload(result))
            }
        }
        emit CallForwarded(_target, success);
    }

    /// @notice Checks if a function selector corresponds to a dangerous function
    /// @param selector The function selector to check
    /// @return True if the function is dangerous, false otherwise
    /// @dev Dangerous functions include: upgrade, changeProxyAdmin, setAddressManager, setBatcherHash, recover, hold, etc.
    function _isDangerousFunction(bytes4 selector) internal pure returns (bool) {
        return
            selector == bytes4(keccak256("upgrade(address,address)")) ||
            selector == bytes4(keccak256("upgradeAndCall(address,address,bytes)")) ||
            selector == bytes4(keccak256("changeProxyAdmin(address,address)")) ||
            selector == bytes4(keccak256("setAddress(string,address)")) ||
            selector == bytes4(keccak256("setAddressManager(address)")) ||
            selector == bytes4(keccak256("setBatcherHash(bytes32)")) ||
            selector == bytes4(keccak256("setImplementation(uint32,address)")) ||
            selector == bytes4(keccak256("transferOwnership(address)")) ||
            selector == bytes4(keccak256("recover(uint256)")) ||
            selector == bytes4(keccak256("hold(address,uint256)"));
    }
}
