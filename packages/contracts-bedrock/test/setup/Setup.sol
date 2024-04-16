// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Predeploys } from "src/libraries/Predeploys.sol";
import { L2CrossDomainMessenger } from "src/L2/L2CrossDomainMessenger.sol";
import { L2StandardBridge } from "src/L2/L2StandardBridge.sol";
import { L2ToL1MessagePasser } from "src/L2/L2ToL1MessagePasser.sol";
import { L2ERC721Bridge } from "src/L2/L2ERC721Bridge.sol";
import { BaseFeeVault } from "src/L2/BaseFeeVault.sol";
import { SequencerFeeVault } from "src/L2/SequencerFeeVault.sol";
import { L1FeeVault } from "src/L2/L1FeeVault.sol";
import { GasPriceOracle } from "src/L2/GasPriceOracle.sol";
import { L1Block } from "src/L2/L1Block.sol";
import { LegacyMessagePasser } from "src/legacy/LegacyMessagePasser.sol";
import { GovernanceToken } from "src/governance/GovernanceToken.sol";
import { OptimismMintableERC20Factory } from "src/universal/OptimismMintableERC20Factory.sol";
import { StandardBridge } from "src/universal/StandardBridge.sol";
import { FeeVault } from "src/universal/FeeVault.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { OptimismPortal2 } from "src/L1/OptimismPortal2.sol";
import { DisputeGameFactory } from "src/dispute/DisputeGameFactory.sol";
import { DelayedWETH } from "src/dispute/weth/DelayedWETH.sol";
import { AnchorStateRegistry } from "src/dispute/AnchorStateRegistry.sol";
import { L1CrossDomainMessenger } from "src/L1/L1CrossDomainMessenger.sol";
import { DeployConfig } from "scripts/DeployConfig.s.sol";
import { Deploy } from "scripts/Deploy.s.sol";
import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
import { ProtocolVersions } from "src/L1/ProtocolVersions.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { L1StandardBridge } from "src/L1/L1StandardBridge.sol";
import { AddressManager } from "src/legacy/AddressManager.sol";
import { L1ERC721Bridge } from "src/L1/L1ERC721Bridge.sol";
import { AddressAliasHelper } from "src/vendor/AddressAliasHelper.sol";
import { Executables } from "scripts/Executables.sol";
import { Vm } from "forge-std/Vm.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { DataAvailabilityChallenge } from "src/L1/DataAvailabilityChallenge.sol";

/// @title Setup
/// @dev This contact is responsible for setting up the contracts in state. It currently
///      sets the L2 contracts directly at the predeploy addresses instead of setting them
///      up behind proxies. In the future we will migrate to importing the genesis JSON
///      file that is created to set up the L2 contracts instead of setting them up manually.
contract Setup {
    error FfiFailed(string);

    /// @notice The address of the foundry Vm contract.
    Vm private constant vm = Vm(0x7109709ECfa91a80626fF3989D68f67F5b1DD12D);

    /// @notice The address of the Deploy contract. Set into state with `etch` to avoid
    ///         mutating any nonces. MUST not have constructor logic.
    Deploy internal constant deploy = Deploy(address(uint160(uint256(keccak256(abi.encode("optimism.deploy"))))));

    OptimismPortal optimismPortal;
    OptimismPortal2 optimismPortal2;
    DisputeGameFactory disputeGameFactory;
    DelayedWETH delayedWeth;
    L2OutputOracle l2OutputOracle;
    SystemConfig systemConfig;
    L1StandardBridge l1StandardBridge;
    L1CrossDomainMessenger l1CrossDomainMessenger;
    AddressManager addressManager;
    L1ERC721Bridge l1ERC721Bridge;
    OptimismMintableERC20Factory l1OptimismMintableERC20Factory;
    ProtocolVersions protocolVersions;
    SuperchainConfig superchainConfig;
    DataAvailabilityChallenge dataAvailabilityChallenge;
    AnchorStateRegistry anchorStateRegistry;

    L2CrossDomainMessenger l2CrossDomainMessenger =
        L2CrossDomainMessenger(payable(Predeploys.L2_CROSS_DOMAIN_MESSENGER));
    L2StandardBridge l2StandardBridge = L2StandardBridge(payable(Predeploys.L2_STANDARD_BRIDGE));
    L2ToL1MessagePasser l2ToL1MessagePasser = L2ToL1MessagePasser(payable(Predeploys.L2_TO_L1_MESSAGE_PASSER));
    OptimismMintableERC20Factory l2OptimismMintableERC20Factory =
        OptimismMintableERC20Factory(Predeploys.OPTIMISM_MINTABLE_ERC20_FACTORY);
    L2ERC721Bridge l2ERC721Bridge = L2ERC721Bridge(Predeploys.L2_ERC721_BRIDGE);
    BaseFeeVault baseFeeVault = BaseFeeVault(payable(Predeploys.BASE_FEE_VAULT));
    SequencerFeeVault sequencerFeeVault = SequencerFeeVault(payable(Predeploys.SEQUENCER_FEE_WALLET));
    L1FeeVault l1FeeVault = L1FeeVault(payable(Predeploys.L1_FEE_VAULT));
    GasPriceOracle gasPriceOracle = GasPriceOracle(Predeploys.GAS_PRICE_ORACLE);
    L1Block l1Block = L1Block(Predeploys.L1_BLOCK_ATTRIBUTES);
    LegacyMessagePasser legacyMessagePasser = LegacyMessagePasser(Predeploys.LEGACY_MESSAGE_PASSER);
    GovernanceToken governanceToken = GovernanceToken(Predeploys.GOVERNANCE_TOKEN);

    /// @dev Deploys the Deploy contract without including its bytecode in the bytecode
    ///      of this contract by fetching the bytecode dynamically using `vm.getCode()`.
    ///      If the Deploy bytecode is included in this contract, then it will double
    ///      the compile time and bloat all of the test contract artifacts since they
    ///      will also need to include the bytecode for the Deploy contract.
    ///      This is a hack as we are pushing solidity to the edge.
    function setUp() public virtual {
        vm.etch(address(deploy), vm.getDeployedCode("Deploy.s.sol:Deploy"));
        vm.allowCheatcodes(address(deploy));
        deploy.setUp();
    }

    /// @dev Sets up the L1 contracts.
    function L1() public {
        // Set the deterministic deployer in state to ensure that it is there
        vm.etch(
            0x4e59b44847b379578588920cA78FbF26c0B4956C,
            hex"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf3"
        );

        deploy.run();

        optimismPortal = OptimismPortal(deploy.mustGetAddress("OptimismPortalProxy"));
        optimismPortal2 = OptimismPortal2(deploy.mustGetAddress("OptimismPortalProxy"));
        disputeGameFactory = DisputeGameFactory(deploy.mustGetAddress("DisputeGameFactoryProxy"));
        delayedWeth = DelayedWETH(deploy.mustGetAddress("DelayedWETHProxy"));
        l2OutputOracle = L2OutputOracle(deploy.mustGetAddress("L2OutputOracleProxy"));
        systemConfig = SystemConfig(deploy.mustGetAddress("SystemConfigProxy"));
        l1StandardBridge = L1StandardBridge(deploy.mustGetAddress("L1StandardBridgeProxy"));
        l1CrossDomainMessenger = L1CrossDomainMessenger(deploy.mustGetAddress("L1CrossDomainMessengerProxy"));
        addressManager = AddressManager(deploy.mustGetAddress("AddressManager"));
        l1ERC721Bridge = L1ERC721Bridge(deploy.mustGetAddress("L1ERC721BridgeProxy"));
        l1OptimismMintableERC20Factory =
            OptimismMintableERC20Factory(deploy.mustGetAddress("OptimismMintableERC20FactoryProxy"));
        protocolVersions = ProtocolVersions(deploy.mustGetAddress("ProtocolVersionsProxy"));
        superchainConfig = SuperchainConfig(deploy.mustGetAddress("SuperchainConfigProxy"));
        anchorStateRegistry = AnchorStateRegistry(deploy.mustGetAddress("AnchorStateRegistryProxy"));

        vm.label(address(l2OutputOracle), "L2OutputOracle");
        vm.label(deploy.mustGetAddress("L2OutputOracleProxy"), "L2OutputOracleProxy");
        vm.label(address(optimismPortal), "OptimismPortal");
        vm.label(deploy.mustGetAddress("OptimismPortalProxy"), "OptimismPortalProxy");
        vm.label(address(disputeGameFactory), "DisputeGameFactory");
        vm.label(deploy.mustGetAddress("DisputeGameFactoryProxy"), "DisputeGameFactoryProxy");
        vm.label(address(delayedWeth), "DelayedWETH");
        vm.label(deploy.mustGetAddress("DelayedWETHProxy"), "DelayedWETHProxy");
        vm.label(address(systemConfig), "SystemConfig");
        vm.label(deploy.mustGetAddress("SystemConfigProxy"), "SystemConfigProxy");
        vm.label(address(l1StandardBridge), "L1StandardBridge");
        vm.label(deploy.mustGetAddress("L1StandardBridgeProxy"), "L1StandardBridgeProxy");
        vm.label(address(l1CrossDomainMessenger), "L1CrossDomainMessenger");
        vm.label(deploy.mustGetAddress("L1CrossDomainMessengerProxy"), "L1CrossDomainMessengerProxy");
        vm.label(address(addressManager), "AddressManager");
        vm.label(address(l1ERC721Bridge), "L1ERC721Bridge");
        vm.label(deploy.mustGetAddress("L1ERC721BridgeProxy"), "L1ERC721BridgeProxy");
        vm.label(address(l1OptimismMintableERC20Factory), "OptimismMintableERC20Factory");
        vm.label(deploy.mustGetAddress("OptimismMintableERC20FactoryProxy"), "OptimismMintableERC20FactoryProxy");
        vm.label(address(protocolVersions), "ProtocolVersions");
        vm.label(deploy.mustGetAddress("ProtocolVersionsProxy"), "ProtocolVersionsProxy");
        vm.label(address(superchainConfig), "SuperchainConfig");
        vm.label(deploy.mustGetAddress("SuperchainConfigProxy"), "SuperchainConfigProxy");
        vm.label(AddressAliasHelper.applyL1ToL2Alias(address(l1CrossDomainMessenger)), "L1CrossDomainMessenger_aliased");

        if (deploy.cfg().usePlasma()) {
            dataAvailabilityChallenge =
                DataAvailabilityChallenge(deploy.mustGetAddress("DataAvailabilityChallengeProxy"));
            vm.label(address(dataAvailabilityChallenge), "DataAvailabilityChallengeProxy");
            vm.label(deploy.mustGetAddress("DataAvailabilityChallenge"), "DataAvailabilityChallenge");
        }
    }

    /// @dev Sets up the L2 contracts. Depends on `L1()` being called first.
    function L2() public {
        string memory allocsPath = string.concat(vm.projectRoot(), "/.testdata/genesis.json");
        if (vm.isFile(allocsPath) == false) {
            string[] memory args = new string[](3);
            args[0] = Executables.bash;
            args[1] = "-c";
            args[2] = string.concat(vm.projectRoot(), "/scripts/generate-l2-genesis.sh");
            Vm.FfiResult memory result = vm.tryFfi(args);
            if (result.exitCode != 0) {
                revert FfiFailed(
                    string.concat(
                        "FFI call to generate genesis.json failed with exit code: ",
                        string(abi.encodePacked(result.exitCode)),
                        ".\nCommand: ",
                        Executables.bash,
                        " -c ",
                        vm.projectRoot(),
                        "/scripts/generate-l2-genesis.sh",
                        ".\nOutput: ",
                        string(result.stdout),
                        "\nError: ",
                        string(result.stderr)
                    )
                );
            }
        }

        // Prevent race condition where the genesis.json file is not yet created
        while (vm.isFile(allocsPath) == false) {
            vm.sleep(1);
        }

        vm.loadAllocs(allocsPath);

        // Set the governance token's owner to be the final system owner
        address finalSystemOwner = deploy.cfg().finalSystemOwner();
        vm.prank(governanceToken.owner());
        governanceToken.transferOwnership(finalSystemOwner);

        vm.label(Predeploys.OPTIMISM_MINTABLE_ERC20_FACTORY, "OptimismMintableERC20Factory");
        vm.label(Predeploys.L2_STANDARD_BRIDGE, "L2StandardBridge");
        vm.label(Predeploys.L2_CROSS_DOMAIN_MESSENGER, "L2CrossDomainMessenger");
        vm.label(Predeploys.L2_TO_L1_MESSAGE_PASSER, "L2ToL1MessagePasser");
        vm.label(Predeploys.SEQUENCER_FEE_WALLET, "SequencerFeeVault");
        vm.label(Predeploys.L2_ERC721_BRIDGE, "L2ERC721Bridge");
        vm.label(Predeploys.BASE_FEE_VAULT, "BaseFeeVault");
        vm.label(Predeploys.L1_FEE_VAULT, "L1FeeVault");
        vm.label(Predeploys.L1_BLOCK_ATTRIBUTES, "L1Block");
        vm.label(Predeploys.GAS_PRICE_ORACLE, "GasPriceOracle");
        vm.label(Predeploys.LEGACY_MESSAGE_PASSER, "LegacyMessagePasser");
        vm.label(Predeploys.GOVERNANCE_TOKEN, "GovernanceToken");
        vm.label(Predeploys.EAS, "EAS");
        vm.label(Predeploys.SCHEMA_REGISTRY, "SchemaRegistry");
    }
}
