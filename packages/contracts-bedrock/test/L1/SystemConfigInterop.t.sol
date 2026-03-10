// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import { CommonTest } from "test/setup/CommonTest.sol";

// Target contract dependencies
import { SystemConfigInterop } from "src/L1/SystemConfigInterop.sol";

contract SystemConfigInterop_Test is CommonTest {
    /// @notice Marked virtual to be overridden in
    ///         test/kontrol/deployment/DeploymentSummary.t.sol
    function setUp() public virtual override {
        super.enableInterop();
        super.setUp();
    }

    /// @dev Tests that a dependency can be added.
    function testFuzz_addDependency_succeeds(uint256 _chainId) public {
        vm.prank(systemConfig.owner());
        _systemConfigInterop().addDependency(_chainId);
    }

    /// @dev Tests that adding a dependency as not the owner reverts.
    function testFuzz_addDependency_notOwner_reverts(uint256 _chainId) public {
        vm.expectRevert("Ownable: caller is not the owner");
        _systemConfigInterop().addDependency(_chainId);
    }

    /// @dev Tests that a dependency can be removed.
    function testFuzz_removeDependency_succeeds(uint256 _chainId) public {
        vm.prank(systemConfig.owner());
        _systemConfigInterop().removeDependency(_chainId);
    }

    /// @dev Tests that removing a dependency as not the owner reverts.
    function testFuzz_removeDependency_notOwner_reverts(uint256 _chainId) public {
        vm.expectRevert("Ownable: caller is not the owner");
        _systemConfigInterop().removeDependency(_chainId);
    }

    /// @dev Returns the SystemConfigInterop instance.
    function _systemConfigInterop() internal view returns (SystemConfigInterop) {
        return SystemConfigInterop(address(systemConfig));
    }
}
