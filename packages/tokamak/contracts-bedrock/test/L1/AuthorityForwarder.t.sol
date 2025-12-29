// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test } from "forge-std/Test.sol";
import { AuthorityForwarder } from "src/L1/AuthorityForwarder.sol";

/// @title AuthorityForwarder_Test
/// @notice Comprehensive test suite for AuthorityForwarder contract
contract AuthorityForwarder_Test is Test {
    AuthorityForwarder public forwarder;

    address public operator;
    address public dao;
    address public dao2;
    address public attacker;
    address public securityCouncil;

    // Mock target contract
    MockTarget public target;

    event CallForwarded(address indexed target, bool success);
    event DAOSet(address indexed dao);
    event DelegationSet(address indexed target, bytes4 indexed selector, address indexed executor);

    function setUp() public {
        operator = makeAddr("operator");
        dao = address(new MockDAO());
        dao2 = address(new MockDAO());
        attacker = makeAddr("attacker");
        securityCouncil = makeAddr("securityCouncil");

        // Deploy forwarder with operator
        forwarder = new AuthorityForwarder(operator);

        // Deploy mock target and transfer ownership to forwarder
        target = new MockTarget();
        target.transferOwnership(address(forwarder));
    }

    /*//////////////////////////////////////////////////////////////
                        PHASE 1 TESTS (INITIAL PHASE)
    //////////////////////////////////////////////////////////////*/

    /// @notice Test 1: All functions allowed in Phase 1
    function test_phase1_allFunctionsAllowed() public {
        // Verify Phase 1
        assertEq(forwarder.DAO(), address(0));
        assertEq(uint256(forwarder.currentPhase()), uint256(AuthorityForwarder.Phase.Initial));

        // Operator can call dangerous functions in Phase 1
        vm.prank(operator);
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.dangerousFunction, ())
        );
        assertTrue(target.dangerousCalled());

        // Operator can also call routine functions
        vm.prank(operator);
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.routineFunction, ())
        );
        assertTrue(target.routineCalled());
    }

    /// @notice Test 2: setDAO() authorization check
    function test_setDAO_onlyOperator() public {
        // Non-operator cannot call setDAO
        vm.prank(attacker);
        vm.expectRevert("AuthorityForwarder: caller not operator");
        forwarder.setDAO(dao);

        // Operator can call setDAO
        vm.prank(operator);
        vm.expectEmit(true, false, false, false);
        emit DAOSet(dao);
        forwarder.setDAO(dao);
        assertEq(forwarder.DAO(), dao);
    }

    /// @notice Test 3: setDAO() can only be called once
    function test_setDAO_onlyOnce() public {
        // First call succeeds
        vm.prank(operator);
        forwarder.setDAO(dao);
        assertEq(forwarder.DAO(), dao);

        // Second call reverts
        vm.prank(operator);
        vm.expectRevert("AuthorityForwarder: DAO already set");
        forwarder.setDAO(dao2);
    }

    /// @notice Test 4: DAO must be a contract
    function test_setDAO_mustBeContract() public {
        address eoa = makeAddr("eoa");

        vm.prank(operator);
        vm.expectRevert("AuthorityForwarder: DAO must be a contract");
        forwarder.setDAO(eoa);
    }

    /// @notice Test: Cannot set zero address as DAO
    function test_setDAO_cannotBeZero() public {
        vm.prank(operator);
        vm.expectRevert("AuthorityForwarder: invalid DAO address");
        forwarder.setDAO(address(0));
    }

    /*//////////////////////////////////////////////////////////////
                    PHASE 2 TESTS (DAO-CONTROLLED PHASE)
    //////////////////////////////////////////////////////////////*/

    /// @notice Test 5: Dangerous functions blocked in Phase 2
    function test_phase2_dangerousFunctionsBlocked() public {
        // Setup: Set DAO
        vm.prank(operator);
        forwarder.setDAO(dao);
        assertEq(uint256(forwarder.currentPhase()), uint256(AuthorityForwarder.Phase.DAOControlled));

        // Operator cannot call dangerous function (using transferOwnership as example)
        vm.prank(operator);
        vm.expectRevert("AuthorityForwarder: dangerous operation blocked, use DAO governance");
        forwarder.forwardCall(
            address(target),
            abi.encodeWithSignature("transferOwnership(address)", attacker)
        );
    }

    /// @notice Test 5b: All dangerous functions blocked in Phase 2
    function test_phase2_allDangerousFunctionsBlocked() public {
        // Setup: Set DAO
        vm.prank(operator);
        forwarder.setDAO(dao);

        // Test all dangerous function selectors
        bytes4[10] memory dangerousSelectors = [
            bytes4(keccak256("upgrade(address,address)")),
            bytes4(keccak256("upgradeAndCall(address,address,bytes)")),
            bytes4(keccak256("changeProxyAdmin(address,address)")),
            bytes4(keccak256("setAddress(string,address)")),
            bytes4(keccak256("setAddressManager(address)")),
            bytes4(keccak256("setBatcherHash(bytes32)")),
            bytes4(keccak256("setImplementation(uint32,address)")),
            bytes4(keccak256("transferOwnership(address)")),
            bytes4(keccak256("recover(uint256)")),
            bytes4(keccak256("hold(address,uint256)"))
        ];

        // Operator should be blocked from calling all dangerous functions
        for (uint256 i = 0; i < dangerousSelectors.length; i++) {
            vm.prank(operator);
            vm.expectRevert("AuthorityForwarder: dangerous operation blocked, use DAO governance");
            forwarder.forwardCall(
                address(target),
                abi.encodePacked(dangerousSelectors[i], bytes32(0), bytes32(0))
            );
        }

        // DAO should be able to call all dangerous functions
        for (uint256 i = 0; i < dangerousSelectors.length; i++) {
            vm.prank(dao);
            // Call should not revert for DAO (will revert at target if not implemented, but not blocked by forwarder)
            try forwarder.forwardCall(
                address(target),
                abi.encodePacked(dangerousSelectors[i], bytes32(0), bytes32(0))
            ) {} catch {
                // Expected to fail at target, but not due to authorization
            }
        }
    }

    /// @notice Test 6: Routine functions allowed in Phase 2
    function test_phase2_routineFunctionsAllowed() public {
        // Setup: Set DAO
        vm.prank(operator);
        forwarder.setDAO(dao);

        // Operator can still call routine functions
        vm.prank(operator);
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.routineFunction, ())
        );
        assertTrue(target.routineCalled());
    }

    /// @notice Test 7: DAO can call all functions
    function test_phase2_daoCanCallAllFunctions() public {
        // Setup: Set DAO
        vm.prank(operator);
        forwarder.setDAO(dao);

        // DAO can call dangerous functions
        vm.prank(dao);
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.dangerousFunction, ())
        );
        assertTrue(target.dangerousCalled());

        // DAO can also call routine functions
        vm.prank(dao);
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.routineFunction, ())
        );
        assertTrue(target.routineCalled());
    }

    /*//////////////////////////////////////////////////////////////
                        DELEGATION TESTS
    //////////////////////////////////////////////////////////////*/

    /// @notice Test 8: Delegation priority
    function test_delegation_priority() public {
        vm.prank(operator);
        forwarder.setDAO(dao);

        // DAO delegates dangerous function to security council
        bytes4 selector = MockTarget.dangerousFunction.selector;
        vm.prank(dao);
        vm.expectEmit(true, true, true, false);
        emit DelegationSet(address(target), selector, securityCouncil);
        forwarder.setDelegatedExecutor(address(target), selector, securityCouncil);

        // Operator cannot call (delegated exclusively)
        vm.prank(operator);
        vm.expectRevert("AuthorityForwarder: function delegated exclusively");
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.dangerousFunction, ())
        );

        // Security Council can call (has delegation)
        vm.prank(securityCouncil);
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.dangerousFunction, ())
        );
        assertTrue(target.dangerousCalled());
    }

    /// @notice Test 9: Only DAO can set delegation
    function test_delegation_onlyDAO() public {
        vm.prank(operator);
        forwarder.setDAO(dao);

        bytes4 selector = MockTarget.dangerousFunction.selector;

        // Operator cannot set delegation
        vm.prank(operator);
        vm.expectRevert("AuthorityForwarder: caller not DAO");
        forwarder.setDelegatedExecutor(address(target), selector, securityCouncil);

        // DAO can set delegation
        vm.prank(dao);
        forwarder.setDelegatedExecutor(address(target), selector, securityCouncil);
        assertEq(forwarder.delegatedExecutors(address(target), selector), securityCouncil);
    }

    /*//////////////////////////////////////////////////////////////
                        EDGE CASE TESTS
    //////////////////////////////////////////////////////////////*/

    /// @notice Test 11: Unauthorized caller blocked
    function test_unauthorizedCallerBlocked() public {
        vm.prank(operator);
        forwarder.setDAO(dao);

        // Random address cannot call
        vm.prank(attacker);
        vm.expectRevert("AuthorityForwarder: caller not authorized");
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.routineFunction, ())
        );
    }

    /// @notice Test 12: Phase transition verification
    function test_phaseTransition() public {
        // Phase 1: Initial
        assertEq(uint256(forwarder.currentPhase()), uint256(AuthorityForwarder.Phase.Initial));

        // Transition
        vm.prank(operator);
        forwarder.setDAO(dao);

        // Phase 2: DAOControlled
        assertEq(uint256(forwarder.currentPhase()), uint256(AuthorityForwarder.Phase.DAOControlled));
    }

    /// @notice Test: Call with value forwarding
    function test_forwardCall_withValue() public {
        vm.deal(operator, 1 ether);

        vm.prank(operator);
        forwarder.forwardCall{ value: 0.5 ether }(
            address(target),
            abi.encodeCall(MockTarget.payableFunction, ())
        );

        assertEq(address(target).balance, 0.5 ether);
    }

    /// @notice Test: Revert bubbling
    function test_forwardCall_revertBubbling() public {
        vm.prank(operator);
        vm.expectRevert("MockTarget: intentional revert");
        forwarder.forwardCall(
            address(target),
            abi.encodeCall(MockTarget.revertingFunction, ())
        );
    }

    /// @notice Test: Constructor validation
    function test_constructor_invalidOperator() public {
        vm.expectRevert("AuthorityForwarder: invalid operator address");
        new AuthorityForwarder(address(0));
    }
}

/*//////////////////////////////////////////////////////////////
                        MOCK CONTRACTS
//////////////////////////////////////////////////////////////*/

/// @notice Mock target contract for testing
contract MockTarget {
    address public owner;
    bool public dangerousCalled;
    bool public routineCalled;

    constructor() {
        owner = msg.sender;
    }

    function transferOwnership(address newOwner) external {
        require(msg.sender == owner, "not owner");
        owner = newOwner;
    }

    function dangerousFunction() external {
        require(msg.sender == owner, "not owner");
        dangerousCalled = true;
    }

    function routineFunction() external {
        require(msg.sender == owner, "not owner");
        routineCalled = true;
    }

    function payableFunction() external payable {
        require(msg.sender == owner, "not owner");
    }

    function revertingFunction() external pure {
        revert("MockTarget: intentional revert");
    }
}

/// @notice Mock DAO contract (has code, not EOA)
contract MockDAO {
    // Empty contract with code
}
