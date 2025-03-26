// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { Script } from "forge-std/Script.sol";
import { console2 as console } from "forge-std/console2.sol";
import { Vm } from "forge-std/Vm.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { stdJson } from "forge-std/StdJson.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { ProxyAdmin } from "src/universal/ProxyAdmin.sol";
import { GnosisSafe as Safe } from "safe-contracts/GnosisSafe.sol";
import { Enum as SafeOps } from "safe-contracts/common/Enum.sol";
import { Config } from "scripts/Config.sol";
import { ForgeArtifacts } from "scripts/ForgeArtifacts.sol";

/// @notice Represents a deployment. Is serialized to JSON as a key/value
///         pair. Can be accessed from within scripts.
struct Deployment {
    string name;
    address payable addr;
}

contract UpgradeOptimismPortal is Script {

    error DeploymentDoesNotExist(string);
    error InvalidDeployment(string);
    mapping(string => Deployment) internal _namedDeployments;
    Deployment[] internal _newDeployments;
    string internal deploymentOutfile;

    function run() external {
        // ENV
        address optimismPortalProxy = vm.envAddress("OPTIMISM_PORTAL_PROXY_ADDRESS");
        address l2OutputOracleProxy = vm.envAddress("L2_OUTPUT_ORACLE_PROXY_ADDRESS");
        address systemConfigProxy = vm.envAddress("SYSTEM_CONFIG_PROXY_ADDRESS");
        address superchainConfigProxy = vm.envAddress("SUPERCHAIN_CONFIG_PROXY_ADDRESS");

        // start broadcast
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // deploy new OptimismPortal
        console.log("Deploying OptimismPortal implementation");
        address addr_ = address(new OptimismPortal{ salt: _implSalt() }());
        save("OptimismPortal", addr_);
        console.log("OptimismPortal deployed at %s", addr_);

        // // UpdateAndCall
        // address optimismPortal = mustGetAddress("OptimismPortal");

        // _upgradeAndCallViaSafe({
        //     _proxy: payable(optimismPortalProxy),
        //     _implementation: optimismPortal,
        //     _innerCallData: abi.encodeCall(
        //         OptimismPortal.initialize,
        //         (
        //             L2OutputOracle(l2OutputOracleProxy),
        //             SystemConfig(systemConfigProxy),
        //             SuperchainConfig(superchainConfigProxy)
        //         )
        //     )
        // });

        // OptimismPortal portal = OptimismPortal(payable(optimismPortalProxy));
        // string memory version = portal.version();
        // console.log("OptimismPortal version: %s", version);

        // stop broadcast
        vm.stopBroadcast();
    }

    /// @notice Call from the Safe contract to the Proxy Admin's upgrade and call method
    function _upgradeAndCallViaSafe(address _proxy, address _implementation, bytes memory _innerCallData) internal {
        address proxyAdmin = vm.envAddress("PROXY_ADMIN_ADDRESS");
        address systemOwnerSafe = vm.envAddress("SYSTEM_OWNER_SAFE_ADDRESS");
        bytes memory data =
            abi.encodeCall(ProxyAdmin.upgradeAndCall, (payable(_proxy), _implementation, _innerCallData));

        Safe safe = Safe(payable(systemOwnerSafe));
        _callViaSafe({ _safe: safe, _target: proxyAdmin, _data: data });
    }

    /// @notice Make a call from the Safe contract to an arbitrary address with arbitrary data
    function _callViaSafe(Safe _safe, address _target, bytes memory _data) internal {
        // This is the signature format used when the caller is also the signer.
        bytes memory signature = abi.encodePacked(uint256(uint160(msg.sender)), bytes32(0), uint8(1));

        _safe.execTransaction({
            to: _target,
            value: 0,
            data: _data,
            operation: SafeOps.Operation.Call,
            safeTxGas: 0,
            baseGas: 0,
            gasPrice: 0,
            gasToken: address(0),
            refundReceiver: payable(address(0)),
            signatures: signature
        });
    }

    function _implSalt() internal view returns (bytes32) {
        return keccak256(bytes(Config.implSalt()));
    }

    /// @notice Returns the address of a deployment and reverts if the deployment
    ///         does not exist.
    /// @return The address of the deployment.
    function mustGetAddress(string memory _name) public view returns (address payable) {
        address addr = getAddress(_name);
        if (addr == address(0)) {
            revert DeploymentDoesNotExist(_name);
        }
        return payable(addr);
    }

    /// @notice Appends a deployment to disk as a JSON deploy artifact.
    /// @param _name The name of the deployment.
    /// @param _deployed The address of the deployment.
    function save(string memory _name, address _deployed) public {
        if (bytes(_name).length == 0) {
            revert InvalidDeployment("EmptyName");
        }
        if (bytes(_namedDeployments[_name].name).length > 0) {
            revert InvalidDeployment("AlreadyExists");
        }

        console.log("Saving %s: %s", _name, _deployed);
        Deployment memory deployment = Deployment({ name: _name, addr: payable(_deployed) });
        _namedDeployments[_name] = deployment;
        _newDeployments.push(deployment);
        _appendDeployment(_name, _deployed);
    }

    /// @notice Returns the address of a deployment. Also handles the predeploys.
    /// @param _name The name of the deployment.
    /// @return The address of the deployment. May be `address(0)` if the deployment does not
    ///         exist.
    function getAddress(string memory _name) public view returns (address payable) {
        Deployment memory existing = _namedDeployments[_name];
        if (existing.addr != address(0)) {
            if (bytes(existing.name).length == 0) {
                return payable(address(0));
            }
            return existing.addr;
        }

        bytes32 digest = keccak256(bytes(_name));
        if (digest == keccak256(bytes("L2CrossDomainMessenger"))) {
            return payable(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        } else if (digest == keccak256(bytes("L2ToL1MessagePasser"))) {
            return payable(Predeploys.L2_TO_L1_MESSAGE_PASSER);
        } else if (digest == keccak256(bytes("L2StandardBridge"))) {
            return payable(Predeploys.L2_STANDARD_BRIDGE);
        } else if (digest == keccak256(bytes("L2ERC721Bridge"))) {
            return payable(Predeploys.L2_ERC721_BRIDGE);
        } else if (digest == keccak256(bytes("SequencerFeeWallet"))) {
            return payable(Predeploys.SEQUENCER_FEE_WALLET);
        } else if (digest == keccak256(bytes("OptimismMintableERC20Factory"))) {
            return payable(Predeploys.OPTIMISM_MINTABLE_ERC20_FACTORY);
        } else if (digest == keccak256(bytes("OptimismMintableERC721Factory"))) {
            return payable(Predeploys.OPTIMISM_MINTABLE_ERC721_FACTORY);
        } else if (digest == keccak256(bytes("L1Block"))) {
            return payable(Predeploys.L1_BLOCK_ATTRIBUTES);
        } else if (digest == keccak256(bytes("GasPriceOracle"))) {
            return payable(Predeploys.GAS_PRICE_ORACLE);
        } else if (digest == keccak256(bytes("L1MessageSender"))) {
            return payable(Predeploys.L1_MESSAGE_SENDER);
        } else if (digest == keccak256(bytes("DeployerWhitelist"))) {
            return payable(Predeploys.DEPLOYER_WHITELIST);
        } else if (digest == keccak256(bytes("WETH"))) {
            return payable(Predeploys.WETH);
        } else if (digest == keccak256(bytes("LegacyERC20NativeToken"))) {
            return payable(Predeploys.LEGACY_ERC20_NATIVE_TOKEN);
        } else if (digest == keccak256(bytes("L1BlockNumber"))) {
            return payable(Predeploys.L1_BLOCK_NUMBER);
        } else if (digest == keccak256(bytes("LegacyMessagePasser"))) {
            return payable(Predeploys.LEGACY_MESSAGE_PASSER);
        } else if (digest == keccak256(bytes("ProxyAdmin"))) {
            return payable(Predeploys.PROXY_ADMIN);
        } else if (digest == keccak256(bytes("BaseFeeVault"))) {
            return payable(Predeploys.BASE_FEE_VAULT);
        } else if (digest == keccak256(bytes("L1FeeVault"))) {
            return payable(Predeploys.L1_FEE_VAULT);
        } else if (digest == keccak256(bytes("GovernanceToken"))) {
            return payable(Predeploys.GOVERNANCE_TOKEN);
        } else if (digest == keccak256(bytes("SchemaRegistry"))) {
            return payable(Predeploys.SCHEMA_REGISTRY);
        } else if (digest == keccak256(bytes("EAS"))) {
            return payable(Predeploys.EAS);
        }
        return payable(address(0));
    }

     /// @notice Adds a deployment to the temp deployments file
    function _appendDeployment(string memory _name, address _deployed) internal {
        deploymentOutfile = Config.deploymentOutfile();
        console.log("Writing artifact to %s", deploymentOutfile);
        ForgeArtifacts.ensurePath(deploymentOutfile);

        vm.writeJson({ json: stdJson.serialize("", _name, _deployed), path: deploymentOutfile });
    }
}
