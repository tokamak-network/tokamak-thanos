// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "forge-std/Test.sol";
import "../../src/AA/SimplePriceOracle.sol";

contract SimplePriceOracleTest is Test {
    SimplePriceOracle oracle;
    address owner = address(this);

    function setUp() public {
        oracle = new SimplePriceOracle(0.0005e18); // 1 TON = 0.0005 ETH
    }

    function test_GetPrice_ReturnsInitialPrice() public view {
        assertEq(oracle.getPrice(), 0.0005e18);
    }

    function test_GetPrice_StaleAfter24Hours() public {
        vm.warp(block.timestamp + 86401); // 24h + 1s
        vm.expectRevert("SimplePriceOracle: stale price");
        oracle.getPrice();
    }

    function test_UpdatePrice_ByOwner() public {
        oracle.updatePrice(0.001e18);
        assertEq(oracle.getPrice(), 0.001e18);
    }

    function test_UpdatePrice_NonOwnerReverts() public {
        vm.prank(address(0xBAD));
        vm.expectRevert("only owner");
        oracle.updatePrice(0.001e18);
    }

    function test_UpdatePrice_ZeroPriceReverts() public {
        vm.expectRevert("SimplePriceOracle: zero price");
        oracle.updatePrice(0);
    }

    function test_UpdatePrice_ResetsStaleTimer() public {
        vm.warp(block.timestamp + 50000);
        oracle.updatePrice(0.001e18);
        // now fresh: should not revert
        assertEq(oracle.getPrice(), 0.001e18);
    }

    function test_LastUpdated_MatchesTimestamp() public view {
        assertEq(oracle.lastUpdated(), block.timestamp);
    }

    function test_GetPrice_StaleAtExact24Hours() public {
        vm.warp(block.timestamp + 86400); // exactly 24h — stale (< not <=)
        vm.expectRevert("SimplePriceOracle: stale price");
        oracle.getPrice();
    }

    function test_Constructor_ZeroPriceReverts() public {
        vm.expectRevert("SimplePriceOracle: zero price");
        new SimplePriceOracle(0);
    }

    function test_TransferOwnership_ByOwner() public {
        address newOwner = makeAddr("newOwner");
        oracle.transferOwnership(newOwner);
        assertEq(oracle.owner(), newOwner);
    }

    function test_TransferOwnership_ZeroAddressReverts() public {
        vm.expectRevert("SimplePriceOracle: zero owner");
        oracle.transferOwnership(address(0));
    }

    function test_TransferOwnership_NonOwnerReverts() public {
        vm.prank(makeAddr("attacker"));
        vm.expectRevert("only owner");
        oracle.transferOwnership(makeAddr("newOwner"));
    }
}
