// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { L2Genesis } from "scripts/L2Genesis.s.sol";
import { OutputMode } from "scripts/libraries/Config.sol";
import { Fork } from "scripts/libraries/Config.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Constants } from "src/libraries/Constants.sol";
import { Process } from "scripts/libraries/Process.sol";

/// @title L2GenesisTest
/// @notice Test suite for L2Genesis script.
contract L2GenesisTest is Test {
    L2Genesis genesis;

    function setUp() public {
        genesis = new L2Genesis();
    }

    /// @notice Creates a default Input for testing.
    function _defaultInput() internal pure returns (L2Genesis.Input memory) {
        return L2Genesis.Input({
            l1ChainID: 1,
            l2ChainID: 10,
            l1CrossDomainMessengerProxy: payable(address(0x100000)),
            l1StandardBridgeProxy: payable(address(0x100001)),
            l1ERC721BridgeProxy: payable(address(0x100002)),
            opChainProxyAdminOwner: address(0x200000),
            sequencerFeeVaultRecipient: address(0x300000),
            sequencerFeeVaultMinimumWithdrawalAmount: 0,
            sequencerFeeVaultWithdrawalNetwork: 0,
            baseFeeVaultRecipient: address(0x400000),
            baseFeeVaultMinimumWithdrawalAmount: 0,
            baseFeeVaultWithdrawalNetwork: 0,
            l1FeeVaultRecipient: address(0x500000),
            l1FeeVaultMinimumWithdrawalAmount: 0,
            l1FeeVaultWithdrawalNetwork: 0,
            operatorFeeVaultRecipient: address(0x600000),
            operatorFeeVaultMinimumWithdrawalAmount: 0,
            operatorFeeVaultWithdrawalNetwork: 0,
            governanceTokenOwner: address(0x700000),
            fork: uint256(Fork.DELTA),
            deployCrossL2Inbox: false,
            enableGovernance: true,
            fundDevAccounts: false,
            useRevenueShare: false,
            chainFeesRecipient: address(0),
            l1FeesDepositor: address(0),
            useCustomGasToken: false,
            gasPayingTokenName: "",
            gasPayingTokenSymbol: "",
            nativeAssetLiquidityAmount: 0,
            liquidityControllerOwner: address(0)
        });
    }

    /// @notice Creates a temp file and returns the path to it.
    function tmpfile() internal returns (string memory) {
        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] = "mktemp";
        bytes memory result = Process.run(commands);
        return string(result);
    }

    /// @notice Deletes a file at a given filesystem path. Does not force delete
    ///         and does not recursively delete.
    function deleteFile(string memory path) internal {
        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] = string.concat("rm ", path);
        Process.run(commands);
    }

    /// @notice Returns the number of top level keys in a JSON object at a given
    ///         file path.
    function getJSONKeyCount(string memory path) internal returns (uint256) {
        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] = string.concat("jq 'keys | length' < ", path, " | xargs cast abi-encode 'f(uint256)'");
        return abi.decode(Process.run(commands), (uint256));
    }

    /// @notice Helper function to run a function with a temporary dump file.
    function withTempDump(function (string memory) internal f) internal {
        string memory path = tmpfile();
        f(path);
        deleteFile(path);
    }

    /// @notice Helper function for reading the number of storage keys for a given account.
    function getStorageKeysCount(string memory _path, address _addr) internal returns (uint256) {
        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] =
            string.concat("jq -r '.[\"", vm.toLowercase(vm.toString(_addr)), "\"].storage | length' < ", _path);
        return vm.parseUint(string(Process.run(commands)));
    }

    /// @notice Returns the number of accounts that contain particular code at a given path to a genesis file.
    function getCodeCount(string memory path, string memory name) internal returns (uint256) {
        bytes memory code = vm.getDeployedCode(name);
        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] = string.concat(
            "jq -r 'map_values(select(.code == \"",
            vm.toString(code),
            "\")) | length' < ",
            path,
            " | xargs cast abi-encode 'f(uint256)'"
        );
        return abi.decode(Process.run(commands), (uint256));
    }

    /// @notice Returns the number of accounts that have a particular slot set.
    function getPredeployCountWithSlotSet(string memory path, bytes32 slot) internal returns (uint256) {
        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] = string.concat(
            "jq 'map_values(.storage | select(has(\"",
            vm.toString(slot),
            "\"))) | keys | length' < ",
            path,
            " | xargs cast abi-encode 'f(uint256)'"
        );
        return abi.decode(Process.run(commands), (uint256));
    }

    /// @notice Returns the number of accounts that have a particular slot set to a particular value.
    function getPredeployCountWithSlotSetToValue(
        string memory path,
        bytes32 slot,
        bytes32 value
    )
        internal
        returns (uint256)
    {
        string[] memory commands = new string[](3);
        commands[0] = "bash";
        commands[1] = "-c";
        commands[2] = string.concat(
            "jq 'map_values(.storage | select(.\"",
            vm.toString(slot),
            "\" == \"",
            vm.toString(value),
            "\")) | length' < ",
            path,
            " | xargs cast abi-encode 'f(uint256)'"
        );
        return abi.decode(Process.run(commands), (uint256));
    }

    /// @notice Tests the genesis predeploys setup using a temp file for the case where deployCrossL2Inbox is false.
    function test_genesis_predeploys_notUsingInterop() external {
        string memory path = tmpfile();
        _test_genesis_predeploys(path, false);
        deleteFile(path);
    }

    /// @notice Tests the genesis predeploys setup using a temp file for the case where deployCrossL2Inbox is true.
    function test_genesis_predeploys_usingInterop() external {
        string memory path = tmpfile();
        _test_genesis_predeploys(path, true);
        deleteFile(path);
    }

    /// @notice Tests the genesis predeploys setup.
    function _test_genesis_predeploys(string memory _path, bool _useInterop) internal {
        L2Genesis.Input memory input = _defaultInput();
        input.deployCrossL2Inbox = _useInterop;
        if (_useInterop) {
            input.fork = uint256(Fork.INTEROP);
        }

        genesis.run(input);
        vm.dumpState(_path);

        // 2 predeploys do not have proxies (WETH and GovernanceToken)
        assertEq(getCodeCount(_path, "Proxy.sol:Proxy"), Predeploys.PREDEPLOY_COUNT - 2);

        // Count predeploys with implementation slot set.
        // Tokamak adds OPERATOR_FEE_VAULT and FEE_SPLITTER over upstream (17 -> 19 without interop).
        // With interop, CrossL2Inbox and L2ToL2CrossDomainMessenger are also added (19 -> 21).
        assertEq(getPredeployCountWithSlotSet(_path, Constants.PROXY_IMPLEMENTATION_ADDRESS), _useInterop ? 21 : 19);

        // All proxies except 2 have the proxy 1967 admin slot set to the proxy admin
        assertEq(
            getPredeployCountWithSlotSetToValue(
                _path, Constants.PROXY_OWNER_ADDRESS, bytes32(uint256(uint160(Predeploys.PROXY_ADMIN)))
            ),
            Predeploys.PREDEPLOY_COUNT - 2
        );

        // Also see Predeploys.t.test_predeploysSet_succeeds which uses L1Genesis for the CommonTest prestate.
    }

    /// @notice Tests the number of accounts in the genesis setup
    function test_allocs_size() external {
        withTempDump(_test_allocs_size);
    }

    /// @notice Tests the number of accounts in the genesis setup
    function _test_allocs_size(string memory _path) internal {
        L2Genesis.Input memory input = _defaultInput();
        input.fundDevAccounts = false;

        genesis.run(input);
        vm.dumpState(_path);

        uint256 expected = 0;
        expected += 2048 - 2; // predeploy proxies (excludes WETH and GovernanceToken which are direct)
        // predeploy implementations: 19 namespace addresses + WETH (direct) + GovernanceToken (direct, enableGovernance=true)
        expected += 19; // namespace implementation addresses
        expected += 2; // WETH and GovernanceToken at direct predeploy addresses
        expected += 256; // precompiles
        expected += 16; // preinstalls (MultiCall3, Create2Deployer, Safe_v130, SafeL2_v130, MultiSendCallOnly_v130,
                        // SafeSingletonFactory, DeterministicDeploymentProxy, MultiSend_v130, Permit2,
                        // SenderCreator_v060, EntryPoint_v060, SenderCreator_v070, EntryPoint_v070,
                        // BeaconBlockRoots, HistoryStorage, CreateX)
        expected += 2; // BeaconBlockRootsSender and HistoryStorageSender (nonce bumped)
        // 16 prefunded dev accounts are excluded (fundDevAccounts=false)
        assertEq(expected, getJSONKeyCount(_path), "key count check");

        // 3 slots: implementation, owner, admin
        assertEq(3, getStorageKeysCount(_path, Predeploys.PROXY_ADMIN), "proxy admin storage check");
    }
}
