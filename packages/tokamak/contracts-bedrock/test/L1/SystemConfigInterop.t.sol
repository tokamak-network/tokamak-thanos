// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import { CommonTest } from "test/setup/CommonTest.sol";

// Libraries
import { Constants } from "src/libraries/Constants.sol";
import { StaticConfig } from "src/libraries/StaticConfig.sol";
import { GasPayingToken } from "src/libraries/GasPayingToken.sol";

// Target contract dependencies
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { SystemConfigInterop } from "src/L1/SystemConfigInterop.sol";
import { OptimismPortalInterop } from "src/L1/OptimismPortalInterop.sol";
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import { ConfigType } from "src/L2/L1BlockInterop.sol";

contract SystemConfigInterop_Test is CommonTest {
    /// @notice Marked virtual to be overridden in
    ///         test/kontrol/deployment/DeploymentSummary.t.sol
    function setUp() public virtual override {
        super.enableInterop();
        super.setUp();
    }

    /// @dev Tests that a dependency can be added.
    function testFuzz_addDependency_succeeds(uint256 _chainId) public {
        vm.expectCall(
            address(optimismPortal),
            abi.encodeCall(
                OptimismPortalInterop.setConfig, (ConfigType.ADD_DEPENDENCY, StaticConfig.encodeAddDependency(_chainId))
            )
        );

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
        vm.expectCall(
            address(optimismPortal),
            abi.encodeCall(
                OptimismPortalInterop.setConfig,
                (ConfigType.REMOVE_DEPENDENCY, StaticConfig.encodeRemoveDependency(_chainId))
            )
        );

        vm.prank(systemConfig.owner());
        _systemConfigInterop().removeDependency(_chainId);
    }

    /// @dev Tests that removing a dependency as not the owner reverts.
    function testFuzz_removeDependency_notOwner_reverts(uint256 _chainId) public {
        vm.expectRevert("Ownable: caller is not the owner");
        _systemConfigInterop().removeDependency(_chainId);
    }

    /// @dev Helper to clean storage and then initialize the system config with an arbitrary gas token address.
    function _cleanStorageAndInit(address _token) internal {
        // Wipe out the initialized slot so the proxy can be initialized again
        vm.store(address(systemConfig), bytes32(0), bytes32(0));
        vm.store(address(systemConfig), GasPayingToken.GAS_PAYING_TOKEN_SLOT, bytes32(0));
        vm.store(address(systemConfig), GasPayingToken.GAS_PAYING_TOKEN_NAME_SLOT, bytes32(0));
        vm.store(address(systemConfig), GasPayingToken.GAS_PAYING_TOKEN_SYMBOL_SLOT, bytes32(0));

        systemConfig.initialize({
            _owner: alice,
            _basefeeScalar: 2100,
            _blobbasefeeScalar: 1000000,
            _batcherHash: bytes32(hex"abcd"),
            _gasLimit: 30_000_000,
            _unsafeBlockSigner: address(1),
            _config: Constants.DEFAULT_RESOURCE_CONFIG(),
            _batchInbox: address(0),
            _addresses: SystemConfig.Addresses({
                l1CrossDomainMessenger: address(0),
                l1ERC721Bridge: address(0),
                disputeGameFactory: address(0),
                l1StandardBridge: address(0),
                optimismPortal: address(optimismPortal),
                optimismMintableERC20Factory: address(0),
                nativeTokenAddress: address(0),
                gasPayingToken: _token
            })
        });
    }

    /// @dev Returns the SystemConfigInterop instance.
    function _systemConfigInterop() internal view returns (SystemConfigInterop) {
        return SystemConfigInterop(address(systemConfig));
    }
}
