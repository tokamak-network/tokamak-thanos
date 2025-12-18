// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { console2 as console } from "forge-std/console2.sol";
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// Target contracts
import { UpgradeL1BridgeV1 } from "src/shutdown/UpgradeL1BridgeV1.sol";
import { L1StandardBridge } from "src/L1/L1StandardBridge.sol";
import { L1ChugSplashProxy } from "src/legacy/L1ChugSplashProxy.sol";

// Test helper - GenFWStorage
import { GenFWStorage1 } from "test/shutdown/GenFWStorage1.sol";

/// @title MockERC20
/// @notice Simple ERC20 for testing
contract MockERC20 is ERC20 {
    constructor() ERC20("Mock Token", "MOCK") {
        _mint(msg.sender, 1000000 * 10 ** 18);
    }

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

/// @title UpgradeL1BridgeV1_E2E_Test
/// @notice End-to-end test for the complete force withdrawal workflow
/// @dev Tests the full flow: Deploy → Upgrade → Register → Activate → Claim
contract UpgradeL1BridgeV1_E2E_Test is Test {
    // Contracts
    L1ChugSplashProxy public proxy;
    L1StandardBridge public standardBridge;
    UpgradeL1BridgeV1 public upgradeBridge;
    UpgradeL1BridgeV1 public bridgeProxy; // Proxy pointing to UpgradeL1BridgeV1
    GenFWStorage1 public genStorage1;
    MockERC20 public token;

    // Actors
    address public owner;
    address public closer;
    address public user1;
    address public user2;

    // Test data
    uint256 public constant CLAIM_AMOUNT_1 = 100 ether;
    uint256 public constant CLAIM_AMOUNT_2 = 50 ether;

    // Expected hashes (precomputed for testing)
    // hash1 = keccak256(abi.encodePacked(address(token), user1, CLAIM_AMOUNT_1))
    bytes32 public hash1;
    // hash2 = keccak256(abi.encodePacked(address(0), user2, CLAIM_AMOUNT_2)) // ETH
    bytes32 public hash2;

    function setUp() public {
        // Setup actors
        owner = address(this);
        closer = makeAddr("closer");
        user1 = makeAddr("user1");
        user2 = makeAddr("user2");

        console.log("\n===========================================");
        console.log("E2E Test: Force Withdrawal Workflow");
        console.log("===========================================\n");

        // Step 0: Deploy initial contracts
        _deployInitialContracts();

        // Step 1: Upgrade to UpgradeL1BridgeV1
        _upgradeBridge();

        // Step 2: Deploy GenFWStorage
        _deployGenStorage();

        // Step 3: Calculate test hashes
        _calculateHashes();

        // Step 3.5: Setup hashes in GenFWStorage
        _setupGenStorageHashes();

        // Fund the bridge with tokens and ETH
        _fundBridge();
    }

    /// ========================================
    /// Step 0: Deploy Initial Contracts
    /// ========================================

    function _deployInitialContracts() internal {
        console.log("Step 0: Deploying initial contracts...");

        // Deploy L1StandardBridge implementation
        standardBridge = new L1StandardBridge();
        console.log("  L1StandardBridge deployed at:", address(standardBridge));

        // Deploy L1ChugSplashProxy
        proxy = new L1ChugSplashProxy(owner);
        console.log("  L1ChugSplashProxy deployed at:", address(proxy));

        // Set initial implementation
        bytes memory bytecode = address(standardBridge).code;
        proxy.setCode(bytecode);
        console.log("  Proxy set to L1StandardBridge implementation");

        console.log("  [OK] Step 0 complete\n");
    }

    /// ========================================
    /// Step 1: Upgrade to UpgradeL1BridgeV1
    /// ========================================

    function _upgradeBridge() internal {
        console.log("Step 1: Upgrading bridge to UpgradeL1BridgeV1...");

        // Deploy UpgradeL1BridgeV1 implementation
        upgradeBridge = new UpgradeL1BridgeV1();
        console.log("  UpgradeL1BridgeV1 deployed at:", address(upgradeBridge));

        // Upgrade proxy via setCode
        bytes memory bytecode = address(upgradeBridge).code;
        proxy.setCode(bytecode);
        console.log("  Proxy upgraded to UpgradeL1BridgeV1");

        // Wrap proxy in UpgradeL1BridgeV1 interface
        bridgeProxy = UpgradeL1BridgeV1(payable(address(proxy)));

        // Verify upgrade
        address proxyOwner = bridgeProxy.getProxyOwner();
        console.log("  Proxy owner:", proxyOwner);
        assertEq(proxyOwner, owner, "Proxy owner mismatch");

        console.log("  [OK] Step 1 complete\n");
    }

    /// ========================================
    /// Step 2: Deploy GenFWStorage
    /// ========================================

    function _deployGenStorage() internal {
        console.log("Step 2: Deploying GenFWStorage...");

        genStorage1 = new GenFWStorage1();
        console.log("  GenFWStorage1 deployed at:", address(genStorage1));

        console.log("  [OK] Step 2 complete\n");
    }

    function _setupGenStorageHashes() internal {
        // Setup hash1 in GenFWStorage1
        // Function signature: _<hash>()
        string memory funcName1 = string(abi.encodePacked("_", _bytes32ToHexString(hash1), "()"));
        bytes4 funcSig1 = bytes4(keccak256(bytes(funcName1)));
        genStorage1.setHash(funcSig1, hash1);

        // Setup hash2 in GenFWStorage1
        string memory funcName2 = string(abi.encodePacked("_", _bytes32ToHexString(hash2), "()"));
        bytes4 funcSig2 = bytes4(keccak256(bytes(funcName2)));
        genStorage1.setHash(funcSig2, hash2);

        console.log("  Hashes configured in GenFWStorage1");
    }

    /// ========================================
    /// Step 3: Calculate Test Hashes
    /// ========================================

    function _calculateHashes() internal {
        console.log("Step 3: Calculating test hashes...");

        // Deploy mock token
        token = new MockERC20();
        console.log("  Mock token deployed at:", address(token));

        // Calculate hash1: keccak256(token, user1, CLAIM_AMOUNT_1)
        hash1 = keccak256(abi.encodePacked(address(token), user1, CLAIM_AMOUNT_1));
        console.log("  Hash1 (ERC20):");
        console.logBytes32(hash1);

        // Calculate hash2: keccak256(address(0), user2, CLAIM_AMOUNT_2)
        hash2 = keccak256(abi.encodePacked(address(0), user2, CLAIM_AMOUNT_2));
        console.log("  Hash2 (ETH):");
        console.logBytes32(hash2);

        console.log("  [OK] Step 3 complete\n");
    }

    /// ========================================
    /// Fund Bridge
    /// ========================================

    function _fundBridge() internal {
        console.log("Funding bridge...");

        // Fund with tokens
        token.mint(address(bridgeProxy), 1000 ether);
        console.log("  Bridge token balance:", token.balanceOf(address(bridgeProxy)) / 1e18, "tokens");

        // Fund with ETH
        vm.deal(address(bridgeProxy), 1000 ether);
        console.log("  Bridge ETH balance:", address(bridgeProxy).balance / 1e18, "ETH");

        console.log("  [OK] Funding complete\n");
    }

    /// ========================================
    /// Test: Step 3 - Register Positions
    /// ========================================

    function test_step3_registerPositions() public {
        console.log("=== Test: Step 3 - Register Positions ===\n");

        // Check initial state
        assertFalse(bridgeProxy.position(address(genStorage1)), "Position should not be registered initially");

        // Register position (as closer)
        vm.prank(owner);
        bridgeProxy.setCloser(closer);

        address[] memory positions = new address[](1);
        positions[0] = address(genStorage1);

        vm.prank(closer);
        bridgeProxy.forceRegistry(positions);

        console.log("  Position registered:", address(genStorage1));

        // Verify registration
        assertTrue(bridgeProxy.position(address(genStorage1)), "Position should be registered");

        console.log("  [OK] Position registration verified\n");
    }

    /// ========================================
    /// Test: Step 4 - Activate Force Withdrawal
    /// ========================================

    function test_step4_activateForceWithdrawal() public {
        console.log("=== Test: Step 4 - Activate Force Withdrawal ===\n");

        // Setup: Register position first
        vm.prank(owner);
        bridgeProxy.setCloser(closer);

        // Activate force withdrawal
        vm.prank(closer);
        bridgeProxy.forceActive(true);

        console.log("  Force withdrawal activated");

        // Verify activation
        assertTrue(bridgeProxy.active(), "Force withdrawal should be active");

        console.log("  [OK] Activation verified\n");
    }

    /// ========================================
    /// Test: Step 5 - ERC20 Force Withdrawal Claim
    /// ========================================

    function test_step5_forceWithdrawClaim_ERC20_succeeds() public {
        console.log("=== Test: Step 5 - ERC20 Force Withdrawal ===\n");

        // Setup: Register and activate
        _setupForceWithdrawal();

        // Record initial balance
        uint256 initialBalance = token.balanceOf(user1);
        console.log("  User1 initial balance:", initialBalance / 1e18, "tokens");

        // Execute claim
        string memory hashStr = _bytes32ToHexString(hash1);
        console.log("  Claiming with hash:", hashStr);

        vm.prank(user1);
        bridgeProxy.forceWithdrawClaim(
            address(genStorage1), // position
            hashStr, // hash (without 0x prefix)
            address(token), // token
            CLAIM_AMOUNT_1, // amount
            user1 // recipient
        );

        // Verify claim
        uint256 finalBalance = token.balanceOf(user1);
        console.log("  User1 final balance:", finalBalance / 1e18, "tokens");

        assertEq(finalBalance - initialBalance, CLAIM_AMOUNT_1, "Balance should increase by claim amount");

        // Verify claim state
        assertTrue(bridgeProxy.claimState(hash1), "Claim should be marked as completed");

        console.log("  [OK] ERC20 claim successful\n");
    }

    /// ========================================
    /// Test: Step 6 - ETH Force Withdrawal Claim
    /// ========================================

    function test_step6_forceWithdrawClaim_ETH_succeeds() public {
        console.log("=== Test: Step 6 - ETH Force Withdrawal ===\n");

        // Setup: Register and activate
        _setupForceWithdrawal();

        // Record initial balance
        uint256 initialBalance = user2.balance;
        console.log("  User2 initial balance:", initialBalance / 1e18, "ETH");

        // Execute claim
        string memory hashStr = _bytes32ToHexString(hash2);
        console.log("  Claiming with hash:", hashStr);

        vm.prank(user2);
        bridgeProxy.forceWithdrawClaim(
            address(genStorage1), // position
            hashStr, // hash (without 0x prefix)
            address(0), // ETH
            CLAIM_AMOUNT_2, // amount
            user2 // recipient
        );

        // Verify claim
        uint256 finalBalance = user2.balance;
        console.log("  User2 final balance:", finalBalance / 1e18, "ETH");

        assertEq(finalBalance - initialBalance, CLAIM_AMOUNT_2, "Balance should increase by claim amount");

        // Verify claim state
        assertTrue(bridgeProxy.claimState(hash2), "Claim should be marked as completed");

        console.log("  [OK] ETH claim successful\n");
    }

    /// ========================================
    /// Test: Error Cases
    /// ========================================

    function test_error_invalidHash_reverts() public {
        console.log("=== Test: Invalid Hash Reverts ===\n");

        _setupForceWithdrawal();

        bytes32 wrongHash = keccak256(abi.encodePacked(address(token), user1, uint256(999 ether)));
        string memory hashStr = _bytes32ToHexString(wrongHash);

        vm.expectRevert(UpgradeL1BridgeV1.FW_INVALID_HASH.selector);
        vm.prank(user1);
        bridgeProxy.forceWithdrawClaim(
            address(genStorage1),
            hashStr,
            address(token),
            CLAIM_AMOUNT_1, // Wrong amount
            user1
        );

        console.log("  [OK] Invalid hash correctly reverted\n");
    }

    function test_error_doubleClaimreverts() public {
        console.log("=== Test: Double Claim Reverts ===\n");

        _setupForceWithdrawal();

        string memory hashStr = _bytes32ToHexString(hash1);

        // First claim succeeds
        vm.prank(user1);
        bridgeProxy.forceWithdrawClaim(address(genStorage1), hashStr, address(token), CLAIM_AMOUNT_1, user1);

        console.log("  First claim succeeded");

        // Second claim should fail
        vm.expectRevert("already claim Hash");
        vm.prank(user1);
        bridgeProxy.forceWithdrawClaim(address(genStorage1), hashStr, address(token), CLAIM_AMOUNT_1, user1);

        console.log("  [OK] Double claim correctly reverted\n");
    }

    function test_error_unregisteredPosition_reverts() public {
        console.log("=== Test: Unregistered Position Reverts ===\n");

        _setupForceWithdrawal();

        address fakeStorage = makeAddr("fakeStorage");
        string memory hashStr = _bytes32ToHexString(hash1);

        vm.expectRevert(UpgradeL1BridgeV1.FW_NOT_AVAILABLE_POSITION.selector);
        vm.prank(user1);
        bridgeProxy.forceWithdrawClaim(fakeStorage, hashStr, address(token), CLAIM_AMOUNT_1, user1);

        console.log("  [OK] Unregistered position correctly reverted\n");
    }

    /// ========================================
    /// Test: Batch Claim
    /// ========================================

    function test_step7_batchClaim_succeeds() public {
        console.log("=== Test: Step 7 - Batch Claim ===\n");

        _setupForceWithdrawal();

        // Prepare batch claim params
        UpgradeL1BridgeV1.ForceClaimParam[] memory params = new UpgradeL1BridgeV1.ForceClaimParam[](2);

        params[0] = UpgradeL1BridgeV1.ForceClaimParam({
            position: address(genStorage1),
            hashed: _bytes32ToHexString(hash1),
            token: address(token),
            amount: CLAIM_AMOUNT_1,
            getAddress: user1
        });

        params[1] = UpgradeL1BridgeV1.ForceClaimParam({
            position: address(genStorage1),
            hashed: _bytes32ToHexString(hash2),
            token: address(0),
            amount: CLAIM_AMOUNT_2,
            getAddress: user2
        });

        // Record initial balances
        uint256 user1InitialBalance = token.balanceOf(user1);
        uint256 user2InitialBalance = user2.balance;

        console.log("  Executing batch claim for 2 users...");

        // Execute batch claim
        vm.prank(user1);
        bridgeProxy.forceWithdrawClaimAll(params);

        // Verify both claims
        assertEq(token.balanceOf(user1) - user1InitialBalance, CLAIM_AMOUNT_1, "User1 token claim failed");
        assertEq(user2.balance - user2InitialBalance, CLAIM_AMOUNT_2, "User2 ETH claim failed");

        assertTrue(bridgeProxy.claimState(hash1), "Hash1 should be marked as claimed");
        assertTrue(bridgeProxy.claimState(hash2), "Hash2 should be marked as claimed");

        console.log("  [OK] Batch claim successful for both users\n");
    }

    /// ========================================
    /// Helper Functions
    /// ========================================

    function _setupForceWithdrawal() internal {
        // Set closer
        vm.prank(owner);
        bridgeProxy.setCloser(closer);

        // Register position
        address[] memory positions = new address[](1);
        positions[0] = address(genStorage1);
        vm.prank(closer);
        bridgeProxy.forceRegistry(positions);

        // Activate
        vm.prank(closer);
        bridgeProxy.forceActive(true);
    }

    function _bytes32ToHexString(bytes32 _bytes) internal pure returns (string memory) {
        bytes memory hexString = new bytes(64);
        bytes memory hexAlphabet = "0123456789abcdef";

        for (uint256 i = 0; i < 32; i++) {
            hexString[i * 2] = hexAlphabet[uint8(_bytes[i] >> 4)];
            hexString[i * 2 + 1] = hexAlphabet[uint8(_bytes[i] & 0x0f)];
        }

        return string(hexString);
    }
}
