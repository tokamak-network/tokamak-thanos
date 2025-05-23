// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { VmSafe } from "forge-std/Vm.sol";
import { Script } from "forge-std/Script.sol";

import { console2 as console } from "forge-std/console2.sol";
import { stdJson } from "forge-std/StdJson.sol";

import { GnosisSafe as Safe } from "safe-contracts/GnosisSafe.sol";
import { OwnerManager } from "safe-contracts/base/OwnerManager.sol";
import { GnosisSafeProxyFactory as SafeProxyFactory } from "safe-contracts/proxies/GnosisSafeProxyFactory.sol";
import { Enum as SafeOps } from "safe-contracts/common/Enum.sol";

import { Deployer } from "scripts/Deployer.sol";

import { ProxyAdmin } from "src/universal/ProxyAdmin.sol";
import { AddressManager } from "src/legacy/AddressManager.sol";
import { Proxy } from "src/universal/Proxy.sol";
import { L1StandardBridge } from "src/L1/L1StandardBridge.sol";
import { StandardBridge } from "src/universal/StandardBridge.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { OptimismPortal2 } from "src/L1/OptimismPortal2.sol";
import { OptimismPortalInterop } from "src/L1/OptimismPortalInterop.sol";
import { L1ChugSplashProxy } from "src/legacy/L1ChugSplashProxy.sol";
import { ResolvedDelegateProxy } from "src/legacy/ResolvedDelegateProxy.sol";
import { L1CrossDomainMessenger } from "src/L1/L1CrossDomainMessenger.sol";
import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
import { OptimismMintableERC20Factory } from "src/universal/OptimismMintableERC20Factory.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { SystemConfigInterop } from "src/L1/SystemConfigInterop.sol";
import { ResourceMetering } from "src/L1/ResourceMetering.sol";
import { DataAvailabilityChallenge } from "src/L1/DataAvailabilityChallenge.sol";
import { Constants } from "src/libraries/Constants.sol";
import { DisputeGameFactory } from "src/dispute/DisputeGameFactory.sol";
import { FaultDisputeGame } from "src/dispute/FaultDisputeGame.sol";
import { PermissionedDisputeGame } from "src/dispute/PermissionedDisputeGame.sol";
import { DelayedWETH } from "src/dispute/weth/DelayedWETH.sol";
import { AnchorStateRegistry } from "src/dispute/AnchorStateRegistry.sol";
import { PreimageOracle } from "src/cannon/PreimageOracle.sol";
import { MIPS } from "src/cannon/MIPS.sol";
import { L1ERC721Bridge } from "src/L1/L1ERC721Bridge.sol";
import { ProtocolVersions, ProtocolVersion } from "src/L1/ProtocolVersions.sol";
import { StorageSetter } from "src/universal/StorageSetter.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Chains } from "scripts/Chains.sol";
import { Config } from "scripts/Config.sol";

import { IBigStepper } from "src/dispute/interfaces/IBigStepper.sol";
import { IPreimageOracle } from "src/cannon/interfaces/IPreimageOracle.sol";
import { AlphabetVM } from "test/mocks/AlphabetVM.sol";
import "src/dispute/lib/Types.sol";
import { ChainAssertions } from "scripts/ChainAssertions.sol";
import { Types } from "scripts/Types.sol";
import { LibStateDiff } from "scripts/libraries/LibStateDiff.sol";
import { EIP1967Helper } from "test/mocks/EIP1967Helper.sol";
import { ForgeArtifacts } from "scripts/ForgeArtifacts.sol";
import { Process } from "scripts/libraries/Process.sol";

import { L2NativeToken } from "src/L1/L2NativeToken.sol";
import { L1UsdcBridge } from "src/tokamak-contracts/USDC/L1//tokamak-UsdcBridge/L1UsdcBridge.sol";
import { L1UsdcBridgeProxy } from "src/tokamak-contracts/USDC/L1/tokamak-UsdcBridge/L1UsdcBridgeProxy.sol";

/// @title Deploy
/// @notice Script used to deploy a bedrock system. The entire system is deployed within the `run` function.
///         To add a new contract to the system, add a public function that deploys that individual contract.
///         Then add a call to that function inside of `run`. Be sure to call the `save` function after each
///         deployment so that hardhat-deploy style artifacts can be generated using a call to `sync()`.
///         The `CONTRACT_ADDRESSES_PATH` environment variable can be set to a path that contains a JSON file full of
///         contract name to address pairs. That enables this script to be much more flexible in the way it is used.
///         This contract must not have constructor logic because it is set into state using `etch`.
contract Deploy is Deployer {
    using stdJson for string;

    // A state variable to store the contents of the JSON file
    string jsonDeployment;

    /// @notice FaultDisputeGameParams is a struct that contains the parameters necessary to call
    ///         the function _setFaultGameImplementation. This struct exists because the EVM needs
    ///         to finally adopt PUSHN and get rid of stack too deep once and for all.
    ///         Someday we will look back and laugh about stack too deep, today is not that day.
    struct FaultDisputeGameParams {
        AnchorStateRegistry anchorStateRegistry;
        DelayedWETH weth;
        GameType gameType;
        Claim absolutePrestate;
        IBigStepper faultVm;
        uint256 maxGameDepth;
        Duration maxClockDuration;
    }

    ////////////////////////////////////////////////////////////////
    //                        Modifiers                           //
    ////////////////////////////////////////////////////////////////

    /// @notice Modifier that wraps a function in broadcasting.
    modifier broadcast() {
        vm.startBroadcast(msg.sender);
        _;
        vm.stopBroadcast();
    }

    /// @notice Modifier that will only allow a function to be called on devnet.
    modifier onlyDevnet() {
        uint256 chainid = block.chainid;
        if (chainid == Chains.LocalDevnet || chainid == Chains.GethDevnet) {
            _;
        }
    }

    /// @notice Modifier that will only allow a function to be called on a public
    ///         testnet or devnet.
    modifier onlyTestnetOrDevnet() {
        uint256 chainid = block.chainid;
        if (
            chainid == Chains.Goerli || chainid == Chains.Sepolia || chainid == Chains.LocalDevnet
                || chainid == Chains.GethDevnet
        ) {
            _;
        }
    }

    /// @notice Modifier that wraps a function with statediff recording.
    ///         The returned AccountAccess[] array is then written to
    ///         the `snapshots/state-diff/<name>.json` output file.
    modifier stateDiff() {
        vm.startStateDiffRecording();
        _;
        VmSafe.AccountAccess[] memory accesses = vm.stopAndReturnStateDiff();
        console.log(
            "Writing %d state diff account accesses to snapshots/state-diff/%s.json",
            accesses.length,
            vm.toString(block.chainid)
        );
        string memory json = LibStateDiff.encodeAccountAccesses(accesses);
        string memory statediffPath =
            string.concat(vm.projectRoot(), "/snapshots/state-diff/", vm.toString(block.chainid), ".json");
        vm.writeJson({ json: json, path: statediffPath });
    }

    ////////////////////////////////////////////////////////////////
    //                        Accessors                           //
    ////////////////////////////////////////////////////////////////

    /// @notice The create2 salt used for deployment of the contract implementations.
    ///         Using this helps to reduce config across networks as the implementation
    ///         addresses will be the same across networks when deployed with create2.
    function _implSalt() internal view returns (bytes32) {
        return keccak256(bytes(Config.implSalt()));
    }

    /// @notice Returns the proxy addresses. If a proxy is not found, it will have address(0).
    function _proxies() internal view returns (Types.ContractSet memory proxies_) {
        proxies_ = Types.ContractSet({
            L1CrossDomainMessenger: mustGetAddress("L1CrossDomainMessengerProxy"),
            L1StandardBridge: mustGetAddress("L1StandardBridgeProxy"),
            L2OutputOracle: mustGetAddress("L2OutputOracleProxy"),
            DisputeGameFactory: mustGetAddress("DisputeGameFactoryProxy"),
            DelayedWETH: mustGetAddress("DelayedWETHProxy"),
            PermissionedDelayedWETH: mustGetAddress("PermissionedDelayedWETHProxy"),
            AnchorStateRegistry: mustGetAddress("AnchorStateRegistryProxy"),
            OptimismMintableERC20Factory: mustGetAddress("OptimismMintableERC20FactoryProxy"),
            OptimismPortal: mustGetAddress("OptimismPortalProxy"),
            OptimismPortal2: mustGetAddress("OptimismPortalProxy"),
            SystemConfig: mustGetAddress("SystemConfigProxy"),
            L1ERC721Bridge: mustGetAddress("L1ERC721BridgeProxy"),
            ProtocolVersions: mustGetAddress("ProtocolVersionsProxy"),
            SuperchainConfig: mustGetAddress("SuperchainConfigProxy")
        });
    }

    /// @notice Returns the proxy addresses, not reverting if any are unset.
    function _proxiesUnstrict() internal view returns (Types.ContractSet memory proxies_) {
        proxies_ = Types.ContractSet({
            L1CrossDomainMessenger: getAddress("L1CrossDomainMessengerProxy"),
            L1StandardBridge: getAddress("L1StandardBridgeProxy"),
            L2OutputOracle: getAddress("L2OutputOracleProxy"),
            DisputeGameFactory: getAddress("DisputeGameFactoryProxy"),
            DelayedWETH: getAddress("DelayedWETHProxy"),
            PermissionedDelayedWETH: getAddress("PermissionedDelayedWETHProxy"),
            AnchorStateRegistry: getAddress("AnchorStateRegistryProxy"),
            OptimismMintableERC20Factory: getAddress("OptimismMintableERC20FactoryProxy"),
            OptimismPortal: getAddress("OptimismPortalProxy"),
            OptimismPortal2: getAddress("OptimismPortalProxy"),
            SystemConfig: getAddress("SystemConfigProxy"),
            L1ERC721Bridge: getAddress("L1ERC721BridgeProxy"),
            ProtocolVersions: getAddress("ProtocolVersionsProxy"),
            SuperchainConfig: getAddress("SuperchainConfigProxy")
        });
    }

    ////////////////////////////////////////////////////////////////
    //            State Changing Helper Functions                 //
    ////////////////////////////////////////////////////////////////

    /// @notice Gets the address of the SafeProxyFactory and Safe singleton for use in deploying a new GnosisSafe.
    function _getSafeFactory() internal returns (SafeProxyFactory safeProxyFactory_, Safe safeSingleton_) {
        if (getAddress("SafeProxyFactory") != address(0)) {
            // The SafeProxyFactory is already saved, we can just use it.
            safeProxyFactory_ = SafeProxyFactory(getAddress("SafeProxyFactory"));
            safeSingleton_ = Safe(getAddress("SafeSingleton"));
            return (safeProxyFactory_, safeSingleton_);
        }

        // These are the standard create2 deployed contracts. First we'll check if they are deployed,
        // if not we'll deploy new ones, though not at these addresses.
        address safeProxyFactory = 0xa6B71E26C5e0845f74c812102Ca7114b6a896AB2;
        address safeSingleton = 0xd9Db270c1B5E3Bd161E8c8503c55cEABeE709552;

        safeProxyFactory.code.length == 0
            ? safeProxyFactory_ = new SafeProxyFactory()
            : safeProxyFactory_ = SafeProxyFactory(safeProxyFactory);

        safeSingleton.code.length == 0 ? safeSingleton_ = new Safe() : safeSingleton_ = Safe(payable(safeSingleton));

        save("SafeProxyFactory", address(safeProxyFactory_));
        save("SafeSingleton", address(safeSingleton_));
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

    /// @notice Call from the Safe contract to the Proxy Admin's upgrade and call method
    function _upgradeAndCallViaSafe(address _proxy, address _implementation, bytes memory _innerCallData) internal {
        address proxyAdmin = mustGetAddress("ProxyAdmin");

        bytes memory data =
            abi.encodeCall(ProxyAdmin.upgradeAndCall, (payable(_proxy), _implementation, _innerCallData));

        Safe safe = Safe(mustGetAddress("SystemOwnerSafe"));
        _callViaSafe({ _safe: safe, _target: proxyAdmin, _data: data });
    }

    function verifyImplementationExists(address impl) internal view returns (bool) {
        uint256 codeSize;
        assembly {
            codeSize := extcodesize(impl)
        }
        return codeSize > 0;
    }

    /// @notice Transfer ownership of the ProxyAdmin contract to the final system owner
    function transferProxyAdminOwnership() public broadcast {
        ProxyAdmin proxyAdmin = ProxyAdmin(mustGetAddress("ProxyAdmin"));
        address owner = proxyAdmin.owner();
        address safe = mustGetAddress("SystemOwnerSafe");
        if (owner != safe) {
            proxyAdmin.transferOwnership(safe);
            console.log("ProxyAdmin ownership transferred to Safe at: %s", safe);
        }
    }

    /// @notice Transfer ownership of a Proxy to the ProxyAdmin contract
    ///         This is expected to be used in conjusting with deployERC1967ProxyWithOwner after setup actions
    ///         have been performed on the proxy.
    /// @param _name The name of the proxy to transfer ownership of.
    function transferProxyToProxyAdmin(string memory _name) public broadcast {
        Proxy proxy = Proxy(mustGetAddress(_name));
        address proxyAdmin = mustGetAddress("ProxyAdmin");
        proxy.changeAdmin(proxyAdmin);
        console.log("Proxy %s ownership transferred to ProxyAdmin at: %s", _name, proxyAdmin);
    }

    ////////////////////////////////////////////////////////////////
    //                    SetUp and Run                           //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploy all of the L1 contracts necessary for a full Superchain with a single Op Chain.
    function run() public {
        console.log("Deploying a fresh OP Stack including SuperchainConfig");
        _run();
    }

    function runWithStateDump() public {
        vm.chainId(cfg.l1ChainID());
        _run();
        vm.dumpState(Config.stateDumpPath(""));
    }

    /// @notice Deploy all L1 contracts and write the state diff to a file.DeployConfig
    function runWithStateDiff() public stateDiff {
        _run();
    }

    /// @notice Read the JSON file once and store it in jsonCache to get the deployed implementation contract address
    function getDeploymentReUse() public {
        uint256 chainid = block.chainid;
        /// After deployment to mainnet we will add path of deployment file
        // if (chainid == Chains.Mainnet) {
        //     jsonCache = vm.readFile("./")
        // }
        if (cfg.reuseDeployment()) {
            if (chainid == Chains.Sepolia) {
                jsonDeployment = vm.readFile("./deployments/thanos-stack-sepolia/address.json");
            }
            // for test
            else if (chainid == Chains.LocalDevnet) {
                jsonDeployment = vm.readFile("./deployments/devnetL1/addresses.json");
            }
        }
    }

    /// @notice Internal function containing the deploy logic.
    function _run() internal virtual {
        console.log("start of L1 Deploy!");
        deploySafe("SystemOwnerSafe");
        console.log("deployed Safe!");
        setupSuperchain();
        console.log("set up superchain!");
        if (cfg.usePlasma()) {
            bytes32 typeHash = keccak256(bytes(cfg.daCommitmentType()));
            bytes32 keccakHash = keccak256(bytes("KeccakCommitment"));
            if (typeHash == keccakHash) {
                setupOpPlasma();
            }
        }
        // Get the deployment following the deploying network
        if (cfg.reuseDeployment()) {
            getDeploymentReUse();
        }
        setupOpChain();
        console.log("set up op chain!");
    }

    ////////////////////////////////////////////////////////////////
    //           High Level Deployment Functions                  //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploy a full system with a new SuperchainConfig
    ///         The Superchain system has 2 singleton contracts which lie outside of an OP Chain:
    ///         1. The SuperchainConfig contract
    ///         2. The ProtocolVersions contract
    function setupSuperchain() public {
        console.log("Setting up Superchain");

        // Deploy a new ProxyAdmin and AddressManager
        // This proxy will be used on the SuperchainConfig and ProtocolVersions contracts, as well as the contracts
        // in the OP Chain system.
        deployAddressManager();
        deployProxyAdmin();
        transferProxyAdminOwnership();

        // Deploy the SuperchainConfigProxy
        deployERC1967Proxy("SuperchainConfigProxy");
        deploySuperchainConfig();
        initializeSuperchainConfig();

        // Deploy the ProtocolVersionsProxy
        deployERC1967Proxy("ProtocolVersionsProxy");
        deployProtocolVersions();
        initializeProtocolVersions();
    }

    /// @notice Deploy a new OP Chain, with an existing SuperchainConfig provided
    function setupOpChain() public {
        console.log("Deploying OP Chain");

        // Ensure that the requisite contracts are deployed
        mustGetAddress("SuperchainConfigProxy");
        mustGetAddress("SystemOwnerSafe");
        mustGetAddress("AddressManager");
        mustGetAddress("ProxyAdmin");

        if (cfg.devnet()) {
            deployL2NativeToken();
        }
        deployProxies();
        deployImplementations();
        deployL1UsdcBridgeProxy();
        setL1UsdcBridge();
        initializeImplementations();

        setAlphabetFaultGameImplementation({ _allowUpgrade: false });
        setFastFaultGameImplementation({ _allowUpgrade: false });
        setCannonFaultGameImplementation({ _allowUpgrade: false });
        setPermissionedCannonFaultGameImplementation({ _allowUpgrade: false });

        transferDisputeGameFactoryOwnership();
        transferDelayedWETHOwnership();
    }

    /// @notice Deploy all of the proxies
    function deployProxies() public {
        console.log("Deploying proxies");

        deployERC1967Proxy("OptimismPortalProxy");
        deployERC1967Proxy("SystemConfigProxy");
        deployL1StandardBridgeProxy();
        deployL1CrossDomainMessengerProxy();
        deployERC1967Proxy("OptimismMintableERC20FactoryProxy");
        deployERC1967Proxy("L1ERC721BridgeProxy");

        // Both the DisputeGameFactory and L2OutputOracle proxies are deployed regardless of whether fault proofs is
        // enabled to prevent a nastier refactor to the deploy scripts. In the future, the L2OutputOracle will be
        // removed. If fault proofs are not enabled, the DisputeGameFactory proxy will be unused.
        deployERC1967Proxy("DisputeGameFactoryProxy");
        deployERC1967Proxy("L2OutputOracleProxy");
        deployERC1967Proxy("DelayedWETHProxy");
        deployERC1967Proxy("PermissionedDelayedWETHProxy");
        deployERC1967Proxy("AnchorStateRegistryProxy");

        transferAddressManagerOwnership(); // to the ProxyAdmin
    }

    /// @notice Deploy all of the implementations
    function deployImplementations() public {
        console.log("Deploying implementations");

        if (cfg.reuseDeployment() == false) {
            deploySystemConfig();
            deployL1StandardBridge();
            deployL1ERC721Bridge();
            deployOptimismMintableERC20Factory();
            deployL1CrossDomainMessenger();
            deployL2OutputOracle();
            deployOptimismPortal();
            // Fault proofs
            deployOptimismPortal2();
            deployDisputeGameFactory();
            deployDelayedWETH();
            deployAnchorStateRegistry();
        }

        deployPreimageOracle();
        deployMips();
        // USDC bridge
        deployL1UsdcBridge();
    }

    /// @notice Initialize all of the implementations
    function initializeImplementations() public {
        console.log("Initializing implementations");
        // Selectively initialize either the original OptimismPortal or the new OptimismPortal2. Since this will upgrade
        // the proxy, we cannot initialize both.
        if (cfg.useFaultProofs()) {
            console.log("Fault proofs enabled. Initializing the OptimismPortal proxy with the OptimismPortal2.");
            initializeOptimismPortal2();
        } else {
            initializeOptimismPortal();
        }

        initializeSystemConfig();
        initializeL1StandardBridge();
        initializeL1ERC721Bridge();
        initializeOptimismMintableERC20Factory();
        initializeL1CrossDomainMessenger();
        initializeL2OutputOracle();
        initializeDisputeGameFactory();
        initializeDelayedWETH();
        initializePermissionedDelayedWETH();
        initializeAnchorStateRegistry();
    }

    /// @notice Add Plasma setup to the OP chain
    function setupOpPlasma() public {
        console.log("Deploying OP Plasma");
        deployDataAvailabilityChallengeProxy();
        deployDataAvailabilityChallenge();
        initializeDataAvailabilityChallenge();
    }

    ////////////////////////////////////////////////////////////////
    //              Non-Proxied Deployment Functions              //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploy the Safe
    function deploySafe(string memory _name) public broadcast returns (address addr_) {
        address[] memory owners = new address[](0);
        addr_ = deploySafe(_name, owners, 1, true);
    }

    /// @notice Deploy a new Safe contract. If the keepDeployer option is used to enable further setup actions, then
    ///         the removeDeployerFromSafe() function should be called on that safe after setup is complete.
    ///         Note this function does not have the broadcast modifier.
    /// @param _name The name of the Safe to deploy.
    /// @param _owners The owners of the Safe.
    /// @param _threshold The threshold of the Safe.
    /// @param _keepDeployer Wether or not the deployer address will be added as an owner of the Safe.
    function deploySafe(
        string memory _name,
        address[] memory _owners,
        uint256 _threshold,
        bool _keepDeployer
    )
        public
        returns (address addr_)
    {
        bytes32 salt = keccak256(abi.encode(_name, _implSalt()));
        console.log("Deploying safe: %s with salt %s", _name, vm.toString(salt));
        (SafeProxyFactory safeProxyFactory, Safe safeSingleton) = _getSafeFactory();

        address[] memory expandedOwners = new address[](_owners.length + 1);
        if (_keepDeployer) {
            // By always adding msg.sender first we know that the previousOwner will be SENTINEL_OWNERS, which makes it
            // easier to call removeOwner later.
            expandedOwners[0] = msg.sender;
            for (uint256 i = 0; i < _owners.length; i++) {
                expandedOwners[i + 1] = _owners[i];
            }
            _owners = expandedOwners;
        }

        bytes memory initData = abi.encodeCall(
            Safe.setup, (_owners, _threshold, address(0), hex"", address(0), address(0), 0, payable(address(0)))
        );
        addr_ = address(safeProxyFactory.createProxyWithNonce(address(safeSingleton), initData, uint256(salt)));

        save(_name, addr_);
        console.log("New safe: %s deployed at %s\n    Note that this safe is owned by the deployer key", _name, addr_);
    }

    /// @notice If the keepDeployer option was used with deploySafe(), this function can be used to remove the deployer.
    ///         Note this function does not have the broadcast modifier.
    function removeDeployerFromSafe(string memory _name, uint256 _newThreshold) public {
        Safe safe = Safe(mustGetAddress(_name));

        // The sentinel address is used to mark the start and end of the linked list of owners in the Safe.
        address sentinelOwners = address(0x1);

        // Because deploySafe() always adds msg.sender first (if keepDeployer is true), we know that the previousOwner
        // will be sentinelOwners.
        _callViaSafe({
            _safe: safe,
            _target: address(safe),
            _data: abi.encodeCall(OwnerManager.removeOwner, (sentinelOwners, msg.sender, _newThreshold))
        });
        console.log("Removed deployer owner from ", _name);
    }

    /// @notice Deploy the AddressManager
    function deployAddressManager() public broadcast returns (address addr_) {
        console.log("Deploying AddressManager");
        AddressManager manager = new AddressManager();
        require(manager.owner() == msg.sender);

        save("AddressManager", address(manager));
        console.log("AddressManager deployed at %s", address(manager));
        addr_ = address(manager);
    }

    /// @notice Deploy the ProxyAdmin
    function deployProxyAdmin() public broadcast returns (address addr_) {
        console.log("Deploying ProxyAdmin");
        ProxyAdmin admin = new ProxyAdmin({ _owner: msg.sender });
        require(admin.owner() == msg.sender);

        AddressManager addressManager = AddressManager(mustGetAddress("AddressManager"));
        if (admin.addressManager() != addressManager) {
            admin.setAddressManager(addressManager);
        }

        require(admin.addressManager() == addressManager);

        save("ProxyAdmin", address(admin));
        console.log("ProxyAdmin deployed at %s", address(admin));
        addr_ = address(admin);
    }

    /// @notice Deploy the StorageSetter contract, used for upgrades.
    function deployStorageSetter() public broadcast returns (address addr_) {
        console.log("Deploying StorageSetter");
        StorageSetter setter = new StorageSetter{ salt: _implSalt() }();
        console.log("StorageSetter deployed at: %s", address(setter));
        string memory version = setter.version();
        console.log("StorageSetter version: %s", version);
        addr_ = address(setter);
    }

    ////////////////////////////////////////////////////////////////
    //                Proxy Deployment Functions                  //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploy the L1StandardBridgeProxy using a ChugSplashProxy
    function deployL1StandardBridgeProxy() public broadcast returns (address addr_) {
        console.log("Deploying proxy for L1StandardBridge");
        address proxyAdmin = mustGetAddress("ProxyAdmin");
        L1ChugSplashProxy proxy = new L1ChugSplashProxy(proxyAdmin);

        require(EIP1967Helper.getAdmin(address(proxy)) == proxyAdmin);

        save("L1StandardBridgeProxy", address(proxy));
        console.log("L1StandardBridgeProxy deployed at %s", address(proxy));
        addr_ = address(proxy);
    }

    /// @notice Deploy the L1CrossDomainMessengerProxy using a ResolvedDelegateProxy
    function deployL1CrossDomainMessengerProxy() public broadcast returns (address addr_) {
        console.log("Deploying proxy for L1CrossDomainMessenger");
        AddressManager addressManager = AddressManager(mustGetAddress("AddressManager"));
        ResolvedDelegateProxy proxy = new ResolvedDelegateProxy(addressManager, "OVM_L1CrossDomainMessenger");

        save("L1CrossDomainMessengerProxy", address(proxy));
        console.log("L1CrossDomainMessengerProxy deployed at %s", address(proxy));

        addr_ = address(proxy);
    }

    /// @notice Deploys an ERC1967Proxy contract with the ProxyAdmin as the owner.
    /// @param _name The name of the proxy contract to be deployed.
    /// @return addr_ The address of the deployed proxy contract.
    function deployERC1967Proxy(string memory _name) public returns (address addr_) {
        addr_ = deployERC1967ProxyWithOwner(_name, mustGetAddress("ProxyAdmin"));
    }

    /// @notice Deploys an ERC1967Proxy contract with a specified owner.
    /// @param _name The name of the proxy contract to be deployed.
    /// @param _proxyOwner The address of the owner of the proxy contract.
    /// @return addr_ The address of the deployed proxy contract.
    function deployERC1967ProxyWithOwner(
        string memory _name,
        address _proxyOwner
    )
        public
        broadcast
        returns (address addr_)
    {
        console.log(string.concat("Deploying ERC1967 proxy for ", _name));
        Proxy proxy = new Proxy({ _admin: _proxyOwner });

        require(EIP1967Helper.getAdmin(address(proxy)) == _proxyOwner);

        save(_name, address(proxy));
        console.log("   at %s", address(proxy));
        addr_ = address(proxy);
    }

    /// @notice Deploy the DataAvailabilityChallengeProxy
    function deployDataAvailabilityChallengeProxy() public broadcast returns (address addr_) {
        console.log("Deploying proxy for DataAvailabilityChallenge");
        address proxyAdmin = mustGetAddress("ProxyAdmin");
        Proxy proxy = new Proxy({ _admin: proxyAdmin });

        require(EIP1967Helper.getAdmin(address(proxy)) == proxyAdmin);

        save("DataAvailabilityChallengeProxy", address(proxy));
        console.log("DataAvailabilityChallengeProxy deployed at %s", address(proxy));

        addr_ = address(proxy);
    }

    ////////////////////////////////////////////////////////////////
    //             Implementation Deployment Functions            //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploy the SuperchainConfig contract
    function deploySuperchainConfig() public broadcast {
        SuperchainConfig superchainConfig = new SuperchainConfig{ salt: _implSalt() }();

        require(superchainConfig.guardian() == address(0));
        bytes32 initialized = vm.load(address(superchainConfig), bytes32(0));
        require(initialized != 0);

        save("SuperchainConfig", address(superchainConfig));
        console.log("SuperchainConfig deployed at %s", address(superchainConfig));
    }

    /// @notice Deploy the L1CrossDomainMessenger
    function deployL1CrossDomainMessenger() public broadcast returns (address addr_) {
        console.log("Deploying L1CrossDomainMessenger implementation");
        L1CrossDomainMessenger messenger = new L1CrossDomainMessenger{ salt: _implSalt() }();

        save("L1CrossDomainMessenger", address(messenger));
        console.log("L1CrossDomainMessenger deployed at %s", address(messenger));

        // Override the `L1CrossDomainMessenger` contract to the deployed implementation. This is necessary
        // to check the `L1CrossDomainMessenger` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.L1CrossDomainMessenger = address(messenger);
        ChainAssertions.checkL1CrossDomainMessenger({ _contracts: contracts, _vm: vm, _isProxy: false });

        addr_ = address(messenger);
    }

    /// @notice Deploy the OptimismPortal
    function deployOptimismPortal() public broadcast returns (address addr_) {
        console.log("Deploying OptimismPortal implementation");
        if (cfg.useInterop()) {
            addr_ = address(new OptimismPortalInterop{ salt: _implSalt() }());
        } else {
            addr_ = address(new OptimismPortal{ salt: _implSalt() }());
        }
        save("OptimismPortal", addr_);
        console.log("OptimismPortal deployed at %s", addr_);

        // Override the `OptimismPortal` contract to the deployed implementation. This is necessary
        // to check the `OptimismPortal` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.OptimismPortal = addr_;
        ChainAssertions.checkOptimismPortal({ _contracts: contracts, _cfg: cfg, _isProxy: false });
    }

    /// @notice Deploy the OptimismPortal2
    function deployOptimismPortal2() public broadcast returns (address addr_) {
        console.log("Deploying OptimismPortal2 implementation");

        // Could also verify this inside DeployConfig but doing it here is a bit more reliable.
        require(
            uint32(cfg.respectedGameType()) == cfg.respectedGameType(), "Deploy: respectedGameType must fit into uint32"
        );

        OptimismPortal2 portal = new OptimismPortal2{ salt: _implSalt() }({
            _proofMaturityDelaySeconds: cfg.proofMaturityDelaySeconds(),
            _disputeGameFinalityDelaySeconds: cfg.disputeGameFinalityDelaySeconds()
        });

        save("OptimismPortal2", address(portal));
        console.log("OptimismPortal2 deployed at %s", address(portal));

        // Override the `OptimismPortal2` contract to the deployed implementation. This is necessary
        // to check the `OptimismPortal2` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.OptimismPortal2 = address(portal);
        ChainAssertions.checkOptimismPortal2({ _contracts: contracts, _cfg: cfg, _isProxy: false });

        addr_ = address(portal);
    }

    /// @notice Deploy the L2OutputOracle
    function deployL2OutputOracle() public broadcast returns (address addr_) {
        console.log("Deploying L2OutputOracle implementation");
        L2OutputOracle oracle = new L2OutputOracle{ salt: _implSalt() }();

        save("L2OutputOracle", address(oracle));
        console.log("L2OutputOracle deployed at %s", address(oracle));

        // Override the `L2OutputOracle` contract to the deployed implementation. This is necessary
        // to check the `L2OutputOracle` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.L2OutputOracle = address(oracle);
        ChainAssertions.checkL2OutputOracle({
            _contracts: contracts,
            _cfg: cfg,
            _l2OutputOracleStartingTimestamp: 0,
            _isProxy: false
        });

        addr_ = address(oracle);
    }

    /// @notice Deploy the OptimismMintableERC20Factory
    function deployOptimismMintableERC20Factory() public broadcast returns (address addr_) {
        console.log("Deploying OptimismMintableERC20Factory implementation");
        OptimismMintableERC20Factory factory = new OptimismMintableERC20Factory{ salt: _implSalt() }();

        save("OptimismMintableERC20Factory", address(factory));
        console.log("OptimismMintableERC20Factory deployed at %s", address(factory));

        // Override the `OptimismMintableERC20Factory` contract to the deployed implementation. This is necessary
        // to check the `OptimismMintableERC20Factory` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.OptimismMintableERC20Factory = address(factory);
        ChainAssertions.checkOptimismMintableERC20Factory({ _contracts: contracts, _isProxy: false });

        addr_ = address(factory);
    }

    /// @notice Deploy the DisputeGameFactory
    function deployDisputeGameFactory() public broadcast returns (address addr_) {
        console.log("Deploying DisputeGameFactory implementation");
        DisputeGameFactory factory = new DisputeGameFactory{ salt: _implSalt() }();
        save("DisputeGameFactory", address(factory));
        console.log("DisputeGameFactory deployed at %s", address(factory));

        // Override the `DisputeGameFactory` contract to the deployed implementation. This is necessary to check the
        // `DisputeGameFactory` implementation alongside dependent contracts, which are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.DisputeGameFactory = address(factory);
        ChainAssertions.checkDisputeGameFactory({ _contracts: contracts, _expectedOwner: address(0) });

        addr_ = address(factory);
    }

    function deployDelayedWETH() public broadcast returns (address addr_) {
        console.log("Deploying DelayedWETH implementation");
        DelayedWETH weth = new DelayedWETH{ salt: _implSalt() }(cfg.faultGameWithdrawalDelay());
        save("DelayedWETH", address(weth));
        console.log("DelayedWETH deployed at %s", address(weth));

        // Override the `DelayedWETH` contract to the deployed implementation. This is necessary
        // to check the `DelayedWETH` implementation alongside dependent contracts, which are
        // always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.DelayedWETH = address(weth);
        ChainAssertions.checkDelayedWETH({
            _contracts: contracts,
            _cfg: cfg,
            _isProxy: false,
            _expectedOwner: address(0)
        });

        addr_ = address(weth);
    }

    /// @notice Deploy the ProtocolVersions
    function deployProtocolVersions() public broadcast returns (address addr_) {
        console.log("Deploying ProtocolVersions implementation");
        ProtocolVersions versions = new ProtocolVersions{ salt: _implSalt() }();
        save("ProtocolVersions", address(versions));
        console.log("ProtocolVersions deployed at %s", address(versions));

        // Override the `ProtocolVersions` contract to the deployed implementation. This is necessary
        // to check the `ProtocolVersions` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.ProtocolVersions = address(versions);
        ChainAssertions.checkProtocolVersions({ _contracts: contracts, _cfg: cfg, _isProxy: false });

        addr_ = address(versions);
    }

    function getL2NativeToken() public view returns (address) {
        bool isForkPublicNetwork = vm.envOr("FORK_PUBLIC_NETWORK", false);
        if (isForkPublicNetwork) {
            address addr_ = vm.envAddress("L2_NATIVE_TOKEN");
            return addr_;
        }
        return cfg.nativeTokenAddress();
    }

    /// @notice Deploy the Safe
    function deployL2NativeToken() public broadcast {
        string memory path = Config.deployConfigPath();
        L2NativeToken token = new L2NativeToken{ salt: _implSalt() }();
        address addr_ = address(token);
        cfg.setNativeTokenAddress(addr_, path);
        console.log("Native token deployed at", addr_);
        save("L2NativeToken", addr_);
    }

    /// @notice Deploy the L1UsdcBridge
    function deployL1UsdcBridge() public broadcast returns (address addr_) {
        L1UsdcBridge bridge = new L1UsdcBridge{ salt: _implSalt() }();

        require(address(bridge.messenger()) == address(0));
        require(address(bridge.otherBridge()) == address(0));
        require(address(bridge.l1Usdc()) == address(0));
        require(address(bridge.l2Usdc()) == address(0));

        save("L1UsdcBridge", address(bridge));
        console.log("L1UsdcBridge deployed at %s", address(bridge));

        addr_ = address(bridge);
    }

    /// @notice Deploy the L1UsdcBridgeProxy
    function deployL1UsdcBridgeProxy() public broadcast returns (address addr_) {
        address l1UsdcBridge = mustGetAddress("L1UsdcBridge");
        address l1CrossDomainMessengerProxy = mustGetAddress("L1CrossDomainMessengerProxy");
        L1UsdcBridgeProxy proxy =
            new L1UsdcBridgeProxy({ _logic: l1UsdcBridge, initialOwner: msg.sender, _data: abi.encode() });

        require(EIP1967Helper.getAdmin(address(proxy)) == address(msg.sender));

        proxy.setAddress(
            l1CrossDomainMessengerProxy, Predeploys.L2_USDC_BRIDGE, cfg.l1UsdcAddr(), Predeploys.FIATTOKENV2_2
        );
        proxy.upgradeTo(l1UsdcBridge);

        save("L1UsdcBridgeProxy", address(proxy));
        console.log("L1UsdcBridgeProxy deployed at %s", address(proxy));
        addr_ = address(proxy);
    }

    function setL1UsdcBridge() public broadcast {
        address l1UsdcBridgeProxy = mustGetAddress("L1UsdcBridgeProxy");
        address l1CrossDomainMessengerProxy = mustGetAddress("L1CrossDomainMessengerProxy");

        L1UsdcBridge bridge = L1UsdcBridge(l1UsdcBridgeProxy);

        require(address(bridge.messenger()) == l1CrossDomainMessengerProxy);
        require(address(bridge.otherBridge()) == Predeploys.L2_USDC_BRIDGE);
        require(address(bridge.l1Usdc()) == cfg.l1UsdcAddr());
        require(address(bridge.l2Usdc()) == Predeploys.FIATTOKENV2_2);
    }

    /// @notice Deploy the PreimageOracle
    function deployPreimageOracle() public broadcast returns (address addr_) {
        console.log("Deploying PreimageOracle implementation");
        PreimageOracle preimageOracle = new PreimageOracle{ salt: _implSalt() }({
            _minProposalSize: cfg.preimageOracleMinProposalSize(),
            _challengePeriod: cfg.preimageOracleChallengePeriod()
        });
        save("PreimageOracle", address(preimageOracle));
        console.log("PreimageOracle deployed at %s", address(preimageOracle));

        addr_ = address(preimageOracle);
    }

    /// @notice Deploy Mips
    function deployMips() public broadcast returns (address addr_) {
        console.log("Deploying Mips implementation");
        MIPS mips = new MIPS{ salt: _implSalt() }(IPreimageOracle(mustGetAddress("PreimageOracle")));
        save("Mips", address(mips));
        console.log("MIPS deployed at %s", address(mips));

        addr_ = address(mips);
    }

    /// @notice Deploy the AnchorStateRegistry
    function deployAnchorStateRegistry() public broadcast returns (address addr_) {
        console.log("Deploying AnchorStateRegistry implementation");
        AnchorStateRegistry anchorStateRegistry =
            new AnchorStateRegistry{ salt: _implSalt() }(DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy")));
        save("AnchorStateRegistry", address(anchorStateRegistry));
        console.log("AnchorStateRegistry deployed at %s", address(anchorStateRegistry));

        addr_ = address(anchorStateRegistry);
    }

    /// @notice Deploy the SystemConfig
    function deploySystemConfig() public broadcast returns (address addr_) {
        console.log("Deploying SystemConfig implementation");
        if (cfg.useInterop()) {
            addr_ = address(new SystemConfigInterop{ salt: _implSalt() }());
        } else {
            addr_ = address(new SystemConfig{ salt: _implSalt() }());
        }
        save("SystemConfig", addr_);
        console.log("SystemConfig deployed at %s", addr_);

        // Override the `SystemConfig` contract to the deployed implementation. This is necessary
        // to check the `SystemConfig` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.SystemConfig = addr_;
        ChainAssertions.checkSystemConfig({ _contracts: contracts, _cfg: cfg, _isProxy: false });
    }

    /// @notice Deploy the L1StandardBridge
    function deployL1StandardBridge() public broadcast returns (address addr_) {
        console.log("Deploying L1StandardBridge implementation");

        L1StandardBridge bridge = new L1StandardBridge{ salt: _implSalt() }();

        save("L1StandardBridge", address(bridge));
        console.log("L1StandardBridge deployed at %s", address(bridge));

        // Override the `L1StandardBridge` contract to the deployed implementation. This is necessary
        // to check the `L1StandardBridge` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.L1StandardBridge = address(bridge);
        ChainAssertions.checkL1StandardBridge({ _contracts: contracts, _isProxy: false });

        addr_ = address(bridge);
    }

    /// @notice Deploy the L1ERC721Bridge
    function deployL1ERC721Bridge() public broadcast returns (address addr_) {
        console.log("Deploying L1ERC721Bridge implementation");
        L1ERC721Bridge bridge = new L1ERC721Bridge{ salt: _implSalt() }();

        save("L1ERC721Bridge", address(bridge));
        console.log("L1ERC721Bridge deployed at %s", address(bridge));

        // Override the `L1ERC721Bridge` contract to the deployed implementation. This is necessary
        // to check the `L1ERC721Bridge` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.L1ERC721Bridge = address(bridge);

        ChainAssertions.checkL1ERC721Bridge({ _contracts: contracts, _isProxy: false });

        addr_ = address(bridge);
    }

    /// @notice Transfer ownership of the address manager to the ProxyAdmin
    function transferAddressManagerOwnership() public broadcast {
        console.log("Transferring AddressManager ownership to ProxyAdmin");
        AddressManager addressManager = AddressManager(mustGetAddress("AddressManager"));
        address owner = addressManager.owner();
        address proxyAdmin = mustGetAddress("ProxyAdmin");
        if (owner != proxyAdmin) {
            addressManager.transferOwnership(proxyAdmin);
            console.log("AddressManager ownership transferred to %s", proxyAdmin);
        }

        require(addressManager.owner() == proxyAdmin);
    }

    /// @notice Deploy the DataAvailabilityChallenge
    function deployDataAvailabilityChallenge() public broadcast returns (address addr_) {
        console.log("Deploying DataAvailabilityChallenge implementation");
        DataAvailabilityChallenge dac = new DataAvailabilityChallenge();
        save("DataAvailabilityChallenge", address(dac));
        console.log("DataAvailabilityChallenge deployed at %s", address(dac));

        addr_ = address(dac);
    }

    ////////////////////////////////////////////////////////////////
    //                    Initialize Functions                    //
    ////////////////////////////////////////////////////////////////

    /// @notice Initialize the SuperchainConfig
    function initializeSuperchainConfig() public broadcast {
        address payable superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");
        address payable superchainConfig = mustGetAddress("SuperchainConfig");
        _upgradeAndCallViaSafe({
            _proxy: superchainConfigProxy,
            _implementation: superchainConfig,
            _innerCallData: abi.encodeCall(SuperchainConfig.initialize, (cfg.superchainConfigGuardian(), false))
        });

        ChainAssertions.checkSuperchainConfig({ _contracts: _proxiesUnstrict(), _cfg: cfg, _isPaused: false });
    }

    /// @notice Initialize the DisputeGameFactory
    function initializeDisputeGameFactory() public broadcast {
        console.log("Upgrading and initializing DisputeGameFactory proxy");
        address disputeGameFactoryProxy = mustGetAddress("DisputeGameFactoryProxy");

        address disputeGameFactory;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".DisputeGameFactory");
            console.log("DisputeGameFactory address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                disputeGameFactory = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                disputeGameFactory = address(new DisputeGameFactory{ salt: _implSalt() }());
                require(disputeGameFactory != address(0), "DisputeGameFactory deployment failed");
                save("DisputeGameFactory", disputeGameFactory);
                console.log("DisputeGameFactory deployed at %s", disputeGameFactory);
            }
        } else {
            disputeGameFactory = mustGetAddress("DisputeGameFactory");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(disputeGameFactoryProxy),
            _implementation: disputeGameFactory,
            _innerCallData: abi.encodeCall(DisputeGameFactory.initialize, (msg.sender))
        });

        string memory version = DisputeGameFactory(disputeGameFactoryProxy).version();
        console.log("DisputeGameFactory version: %s", version);

        ChainAssertions.checkDisputeGameFactory({ _contracts: _proxiesUnstrict(), _expectedOwner: msg.sender });
    }

    function initializeDelayedWETH() public broadcast {
        console.log("Upgrading and initializing DelayedWETH proxy");
        address delayedWETHProxy = mustGetAddress("DelayedWETHProxy");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");

        address delayedWETH;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".DelayedWETH");
            console.log("DelayedWETH address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                delayedWETH = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                delayedWETH = address(new DelayedWETH{ salt: _implSalt() }(cfg.faultGameWithdrawalDelay()));
                require(delayedWETH != address(0), "DelayedWETH deployment failed");
                save("DelayedWETH", delayedWETH);
                console.log("DelayedWETH deployed at %s", delayedWETH);
            }
        } else {
            delayedWETH = mustGetAddress("DelayedWETH");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(delayedWETHProxy),
            _implementation: delayedWETH,
            _innerCallData: abi.encodeCall(DelayedWETH.initialize, (msg.sender, SuperchainConfig(superchainConfigProxy)))
        });

        string memory version = DelayedWETH(payable(delayedWETHProxy)).version();
        console.log("DelayedWETH version: %s", version);

        ChainAssertions.checkDelayedWETH({
            _contracts: _proxiesUnstrict(),
            _cfg: cfg,
            _isProxy: true,
            _expectedOwner: msg.sender
        });
    }

    function initializePermissionedDelayedWETH() public broadcast {
        console.log("Upgrading and initializing permissioned DelayedWETH proxy");
        address delayedWETHProxy = mustGetAddress("PermissionedDelayedWETHProxy");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");

        address delayedWETH;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".DelayedWETH");
            console.log("DelayedWETH address from JSON: %s", savedAddress);

            // Re-use on-chain deployment
            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                delayedWETH = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                // If we didn't deploy the DelayedWETH contract in initializeDelayedWETH(), deploy it here.
                // if (load("DelayedWETH") == address(0)) {
                    console.log("Implementation from JSON not found on-chain, deploying new one");
                    delayedWETH = address(new DelayedWETH{ salt: _implSalt() }(cfg.faultGameWithdrawalDelay()));
                    require(delayedWETH != address(0), "DelayedWETH deployment failed");
                    save("DelayedWETH", delayedWETH);
                    console.log("DelayedWETH deployed at %s", delayedWETH);
                    // Import from deployment
                // }
                // else {
                //     delayedWETH = mustGetAddress("DelayedWETH");
                // }
            }
        } else {
            delayedWETH = mustGetAddress("DelayedWETH");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(delayedWETHProxy),
            _implementation: delayedWETH,
            _innerCallData: abi.encodeCall(DelayedWETH.initialize, (msg.sender, SuperchainConfig(superchainConfigProxy)))
        });

        string memory version = DelayedWETH(payable(delayedWETHProxy)).version();
        console.log("DelayedWETH version: %s", version);

        ChainAssertions.checkPermissionedDelayedWETH({
            _contracts: _proxiesUnstrict(),
            _cfg: cfg,
            _isProxy: true,
            _expectedOwner: msg.sender
        });
    }

    function initializeAnchorStateRegistry() public broadcast {
        console.log("Upgrading and initializing AnchorStateRegistry proxy");
        address anchorStateRegistryProxy = mustGetAddress("AnchorStateRegistryProxy");
        SuperchainConfig superchainConfig = SuperchainConfig(mustGetAddress("SuperchainConfigProxy"));

        AnchorStateRegistry.StartingAnchorRoot[] memory roots = new AnchorStateRegistry.StartingAnchorRoot[](5);
        roots[0] = AnchorStateRegistry.StartingAnchorRoot({
            gameType: GameTypes.CANNON,
            outputRoot: OutputRoot({
                root: Hash.wrap(cfg.faultGameGenesisOutputRoot()),
                l2BlockNumber: cfg.faultGameGenesisBlock()
            })
        });
        roots[1] = AnchorStateRegistry.StartingAnchorRoot({
            gameType: GameTypes.PERMISSIONED_CANNON,
            outputRoot: OutputRoot({
                root: Hash.wrap(cfg.faultGameGenesisOutputRoot()),
                l2BlockNumber: cfg.faultGameGenesisBlock()
            })
        });
        roots[2] = AnchorStateRegistry.StartingAnchorRoot({
            gameType: GameTypes.ALPHABET,
            outputRoot: OutputRoot({
                root: Hash.wrap(cfg.faultGameGenesisOutputRoot()),
                l2BlockNumber: cfg.faultGameGenesisBlock()
            })
        });
        roots[3] = AnchorStateRegistry.StartingAnchorRoot({
            gameType: GameTypes.ASTERISC,
            outputRoot: OutputRoot({
                root: Hash.wrap(cfg.faultGameGenesisOutputRoot()),
                l2BlockNumber: cfg.faultGameGenesisBlock()
            })
        });
        roots[4] = AnchorStateRegistry.StartingAnchorRoot({
            gameType: GameTypes.FAST,
            outputRoot: OutputRoot({
                root: Hash.wrap(cfg.faultGameGenesisOutputRoot()),
                l2BlockNumber: cfg.faultGameGenesisBlock()
            })
        });

        address anchorStateRegistry;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".AnchorStateRegistry");
            console.log("AnchorStateRegistry address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                anchorStateRegistry = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                anchorStateRegistry = address(
                    new AnchorStateRegistry{ salt: _implSalt() }(
                        DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy"))
                    )
                );
                require(anchorStateRegistry != address(0), "AnchorStateRegistry deployment failed");
                save("AnchorStateRegistry", anchorStateRegistry);
                console.log("AnchorStateRegistry deployed at %s", anchorStateRegistry);
            }
        } else {
             anchorStateRegistry = mustGetAddress("AnchorStateRegistry");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(anchorStateRegistryProxy),
            _implementation: anchorStateRegistry,
            _innerCallData: abi.encodeCall(AnchorStateRegistry.initialize, (roots, superchainConfig))
        });

        string memory version = AnchorStateRegistry(payable(anchorStateRegistryProxy)).version();
        console.log("AnchorStateRegistry version: %s", version);
    }

    /// @notice Initialize the SystemConfig
    function initializeSystemConfig() public broadcast {
        console.log("Upgrading and initializing SystemConfig proxy");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");

        bytes32 batcherHash = bytes32(uint256(uint160(cfg.batchSenderAddress())));

        address customGasTokenAddress = Constants.ETHER;
        if (cfg.useCustomGasToken()) {
            customGasTokenAddress = cfg.customGasTokenAddress();
        }

        address l2NativeTokenAddress = cfg.nativeTokenAddress();
        if (l2NativeTokenAddress == address(0)) {
            L2NativeToken token = new L2NativeToken{ salt: _implSalt() }();
            l2NativeTokenAddress = address(token);
        }
        console.log("Address of l2NativeToken: ", l2NativeTokenAddress);

        address systemConfig;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".SystemConfig");
            console.log("SystemConfig address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                systemConfig = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                systemConfig = address(new SystemConfig{ salt: _implSalt() }());
                require(systemConfig != address(0), "SystemConfig deployment failed");
                save("SystemConfig", systemConfig);
                console.log("SystemConfig deployed at %s", systemConfig);
            }
        } else {
            systemConfig = mustGetAddress("SystemConfig");
        }

        bytes memory _innerCallData = abi.encodeCall(
            SystemConfig.initialize,
            (
                cfg.finalSystemOwner(),
                cfg.basefeeScalar(),
                cfg.blobbasefeeScalar(),
                batcherHash,
                uint64(cfg.l2GenesisBlockGasLimit()),
                cfg.p2pSequencerAddress(),
                Constants.DEFAULT_RESOURCE_CONFIG(),
                cfg.batchInboxAddress(),
                SystemConfig.Addresses({
                    l1CrossDomainMessenger: mustGetAddress("L1CrossDomainMessengerProxy"),
                    l1ERC721Bridge: mustGetAddress("L1ERC721BridgeProxy"),
                    l1StandardBridge: mustGetAddress("L1StandardBridgeProxy"),
                    disputeGameFactory: mustGetAddress("DisputeGameFactoryProxy"),
                    optimismPortal: mustGetAddress("OptimismPortalProxy"),
                    optimismMintableERC20Factory: mustGetAddress("OptimismMintableERC20FactoryProxy"),
                    gasPayingToken: customGasTokenAddress,
                    nativeTokenAddress: l2NativeTokenAddress
                })
            )
        );

        _upgradeAndCallViaSafe({
            _proxy: payable(systemConfigProxy),
            _implementation: systemConfig,
            _innerCallData: _innerCallData
        });
        SystemConfig config = SystemConfig(systemConfigProxy);
        string memory version = config.version();
        console.log("SystemConfig version: %s", version);

        ChainAssertions.checkSystemConfig({ _contracts: _proxies(), _cfg: cfg, _isProxy: true });
    }

    /// @notice Initialize L1StandardBridge
    function initializeL1StandardBridge() public broadcast {
        console.log("Upgrading and initializing L1StandardBridge proxy");

        // Retrieve the ProxyAdmin contract instance and addresses
        address l1StandardBridgeProxy = mustGetAddress("L1StandardBridgeProxy");
        address l1CrossDomainMessengerProxy = mustGetAddress("L1CrossDomainMessengerProxy");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");

        ProxyAdmin proxyAdmin = ProxyAdmin(mustGetAddress("ProxyAdmin"));
        Safe safe = Safe(mustGetAddress("SystemOwnerSafe"));

        address l1StandardBridge;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".L1StandardBridge");
            console.log("L1StandardBridge address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                l1StandardBridge = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                l1StandardBridge = address(new L1StandardBridge{ salt: _implSalt() }());
                require(l1StandardBridge != address(0), "L1StandardBridge deployment failed");
                save("L1StandardBridge", l1StandardBridge);
                console.log("L1StandardBridge deployed at %s", l1StandardBridge);
            }
        } else {
            l1StandardBridge = mustGetAddress("L1StandardBridge");
        }

        // Set the proxy type.
        uint256 proxyType = uint256(proxyAdmin.proxyType(l1StandardBridgeProxy));
        if (proxyType != uint256(ProxyAdmin.ProxyType.CHUGSPLASH)) {
            _callViaSafe({
                _safe: safe,
                _target: address(proxyAdmin),
                _data: abi.encodeCall(ProxyAdmin.setProxyType, (l1StandardBridgeProxy, ProxyAdmin.ProxyType.CHUGSPLASH))
            });
        }
        require(uint256(proxyAdmin.proxyType(l1StandardBridgeProxy)) == uint256(ProxyAdmin.ProxyType.CHUGSPLASH));

        _upgradeAndCallViaSafe({
            _proxy: payable(l1StandardBridgeProxy),
            _implementation: l1StandardBridge,
            _innerCallData: abi.encodeCall(
                L1StandardBridge.initialize,
                (
                    L1CrossDomainMessenger(l1CrossDomainMessengerProxy),
                    SuperchainConfig(superchainConfigProxy),
                    SystemConfig(systemConfigProxy)
                )
            )
        });
        // Retrieve and log the version of the upgraded L1StandardBridge.
        string memory version = L1StandardBridge(payable(l1StandardBridgeProxy)).version();
        console.log("L1StandardBridge version: %s", version);

        // Run chain assertions to verify the deployment.
        ChainAssertions.checkL1StandardBridge({ _contracts: _proxies(), _isProxy: true });
    }

    /// @notice Initialize the L1ERC721Bridge
    function initializeL1ERC721Bridge() public broadcast {
        console.log("Upgrading and initializing L1ERC721Bridge proxy");
        address l1ERC721BridgeProxy = mustGetAddress("L1ERC721BridgeProxy");
        address l1CrossDomainMessengerProxy = mustGetAddress("L1CrossDomainMessengerProxy");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");

        address l1ERC721Bridge;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".L1ERC721Bridge");
            console.log("L1ERC721Bridge address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                l1ERC721Bridge = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                l1ERC721Bridge = address(new L1ERC721Bridge{ salt: _implSalt() }());
                require(l1ERC721Bridge != address(0), "L1ERC721Bridge deployment failed");
                save("L1ERC721Bridge", l1ERC721Bridge);
                console.log("L1ERC721Bridge deployed at %s", l1ERC721Bridge);
            }
        } else {
            l1ERC721Bridge = mustGetAddress("L1ERC721Bridge");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(l1ERC721BridgeProxy),
            _implementation: l1ERC721Bridge,
            _innerCallData: abi.encodeCall(
                L1ERC721Bridge.initialize,
                (L1CrossDomainMessenger(payable(l1CrossDomainMessengerProxy)), SuperchainConfig(superchainConfigProxy))
            )
        });

        L1ERC721Bridge bridge = L1ERC721Bridge(l1ERC721BridgeProxy);
        string memory version = bridge.version();
        console.log("L1ERC721Bridge version: %s", version);

        ChainAssertions.checkL1ERC721Bridge({ _contracts: _proxies(), _isProxy: true });
    }

    /// @notice Initialize the OptimismMintableERC20Factory
    function initializeOptimismMintableERC20Factory() public broadcast {
        console.log("Upgrading and initializing OptimismMintableERC20Factory proxy");
        address optimismMintableERC20FactoryProxy = mustGetAddress("OptimismMintableERC20FactoryProxy");
        address l1StandardBridgeProxy = mustGetAddress("L1StandardBridgeProxy");

        address optimismMintableERC20Factory;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".OptimismMintableERC20Factory");
            console.log("OptimismMintableERC20Factory address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                optimismMintableERC20Factory = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                optimismMintableERC20Factory = address(new OptimismMintableERC20Factory{ salt: _implSalt() }());
                require(optimismMintableERC20Factory != address(0), "OptimismMintableERC20Factory deployment failed");
                save("OptimismMintableERC20Factory", optimismMintableERC20Factory);
                console.log("OptimismMintableERC20Factory deployed at %s", optimismMintableERC20Factory);
            }
        } else {
            optimismMintableERC20Factory = mustGetAddress("OptimismMintableERC20Factory");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(optimismMintableERC20FactoryProxy),
            _implementation: optimismMintableERC20Factory,
            _innerCallData: abi.encodeCall(OptimismMintableERC20Factory.initialize, (l1StandardBridgeProxy))
        });

        OptimismMintableERC20Factory factory = OptimismMintableERC20Factory(optimismMintableERC20FactoryProxy);
        string memory version = factory.version();
        console.log("OptimismMintableERC20Factory version: %s", version);

        ChainAssertions.checkOptimismMintableERC20Factory({ _contracts: _proxies(), _isProxy: true });
    }

    /// @notice initialize L1CrossDomainMessenger
    function initializeL1CrossDomainMessenger() public broadcast {
        console.log("Upgrading and initializing L1CrossDomainMessenger proxy");

        // Retrieve the contract instance and addresses
        address l1CrossDomainMessengerProxy = mustGetAddress("L1CrossDomainMessengerProxy");
        // Retrieve dependent contract addresses.
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");
        address optimismPortalProxy = mustGetAddress("OptimismPortalProxy");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");

        ProxyAdmin proxyAdmin = ProxyAdmin(mustGetAddress("ProxyAdmin"));
        Safe safe = Safe(mustGetAddress("SystemOwnerSafe"));

        address l1CrossDomainMessenger;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".L1CrossDomainMessenger");
            console.log("L1CrossDomainMessenger address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                l1CrossDomainMessenger = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                l1CrossDomainMessenger = address(new L1CrossDomainMessenger{ salt: _implSalt() }());
                require(l1CrossDomainMessenger != address(0), "L1CrossDomainMessenger deployment failed");
                save("L1CrossDomainMessenger", l1CrossDomainMessenger);
                console.log("L1CrossDomainMessenger deployed at %s", l1CrossDomainMessenger);
            }
        } else {
            l1CrossDomainMessenger = mustGetAddress("L1CrossDomainMessenger");
        }

        // Set the proxy type to RESOLVED if not already set.
        uint256 proxyType = uint256(proxyAdmin.proxyType(l1CrossDomainMessengerProxy));

        if (proxyType != uint256(ProxyAdmin.ProxyType.RESOLVED)) {
            _callViaSafe({
                _safe: safe,
                _target: address(proxyAdmin),
                _data: abi.encodeCall(ProxyAdmin.setProxyType, (l1CrossDomainMessengerProxy, ProxyAdmin.ProxyType.RESOLVED))
            });
        }
        require(uint256(proxyAdmin.proxyType(l1CrossDomainMessengerProxy)) == uint256(ProxyAdmin.ProxyType.RESOLVED));

        // Set the implementation name if needed.
        string memory contractName = "OVM_L1CrossDomainMessenger";
        string memory implName = proxyAdmin.implementationName(l1CrossDomainMessenger);
        if (keccak256(bytes(contractName)) != keccak256(bytes(implName))) {
            _callViaSafe({
                _safe: safe,
                _target: address(proxyAdmin),
                _data: abi.encodeCall(ProxyAdmin.setImplementationName, (l1CrossDomainMessengerProxy, contractName))
            });
        }
        require(
            keccak256(bytes(proxyAdmin.implementationName(l1CrossDomainMessengerProxy)))
                == keccak256(bytes(contractName))
        );

        _upgradeAndCallViaSafe({
            _proxy: payable(l1CrossDomainMessengerProxy),
            _implementation: l1CrossDomainMessenger,
            _innerCallData: abi.encodeCall(
                L1CrossDomainMessenger.initialize,
                (
                    SuperchainConfig(superchainConfigProxy),
                    OptimismPortal(payable(optimismPortalProxy)),
                    SystemConfig(systemConfigProxy)
                )
            )
        });
        // Retrieve and log the version of the upgraded L1CrossDomainMessenger.
        L1CrossDomainMessenger messenger = L1CrossDomainMessenger(l1CrossDomainMessengerProxy);
        string memory version = messenger.version();
        console.log("L1CrossDomainMessenger version: %s", version);

        // Verify the deployment using chain assertions.
        ChainAssertions.checkL1CrossDomainMessenger({ _contracts: _proxies(), _vm: vm, _isProxy: true });
    }

    /// @notice Initialize the L2OutputOracle
    function initializeL2OutputOracle() public broadcast {
        console.log("Upgrading and initializing L2OutputOracle proxy");
        address l2OutputOracleProxy = mustGetAddress("L2OutputOracleProxy");

        address l2OutputOracle;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".L2OutputOracle");
            console.log("L2OutputOracle address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                l2OutputOracle = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                l2OutputOracle = address(new L2OutputOracle{ salt: _implSalt() }());
                require(l2OutputOracle != address(0), "L2OutputOracle deployment failed");
                save("L2OutputOracle", l2OutputOracle);
                console.log("L2OutputOracle deployed at %s", l2OutputOracle);
            }
        } else {
            l2OutputOracle = mustGetAddress("L2OutputOracle");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(l2OutputOracleProxy),
            _implementation: l2OutputOracle,
            _innerCallData: abi.encodeCall(
                L2OutputOracle.initialize,
                (
                    cfg.l2OutputOracleSubmissionInterval(),
                    cfg.l2BlockTime(),
                    cfg.l2OutputOracleStartingBlockNumber(),
                    cfg.l2OutputOracleStartingTimestamp(),
                    cfg.l2OutputOracleProposer(),
                    cfg.l2OutputOracleChallenger(),
                    cfg.finalizationPeriodSeconds()
                )
            )
        });

        L2OutputOracle oracle = L2OutputOracle(l2OutputOracleProxy);
        string memory version = oracle.version();
        console.log("L2OutputOracle version: %s", version);

        ChainAssertions.checkL2OutputOracle({
            _contracts: _proxies(),
            _cfg: cfg,
            _l2OutputOracleStartingTimestamp: cfg.l2OutputOracleStartingTimestamp(),
            _isProxy: true
        });
    }

    /// @notice Initialize the OptimismPortal
    function initializeOptimismPortal() public broadcast {
        console.log("Upgrading and initializing OptimismPortal proxy");

        // Retrieve common contract addresses.
        address optimismPortalProxy = mustGetAddress("OptimismPortalProxy");
        address l2OutputOracleProxy = mustGetAddress("L2OutputOracleProxy");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");

        address optimismPortal;

        // Get the OptimismPortal implementation address based on deployment type.
        if (cfg.reuseDeployment()) {
            // For duplicate deployments, try to read from JSON
            address savedAddress = jsonDeployment.readAddress(".OptimismPortal");
            console.log("OptimismPortal address from JSON: %s", savedAddress);
            // Check if the implementation exists on-chain
            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                optimismPortal = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                optimismPortal = address(new OptimismPortal{ salt: _implSalt() }());
                require(optimismPortal != address(0), "OptimismPortal deployment failed");
                save("OptimismPortal", optimismPortal);
                console.log("OptimismPortal deployed at %s", optimismPortal);
            }
        } else {
            optimismPortal = mustGetAddress("OptimismPortal");
        }
        // Upgrade the proxy with the new implementation and initialization data.
        _upgradeAndCallViaSafe({
            _proxy: payable(optimismPortalProxy),
            _implementation: optimismPortal,
            _innerCallData: abi.encodeCall(
                OptimismPortal.initialize,
                (
                    L2OutputOracle(l2OutputOracleProxy),
                    SystemConfig(systemConfigProxy),
                    SuperchainConfig(superchainConfigProxy)
                )
            )
        });

        // Retrieve and log the version of the upgraded OptimismPortal.
        OptimismPortal portal = OptimismPortal(payable(optimismPortalProxy));
        string memory version = portal.version();
        console.log("OptimismPortal version: %s", version);

        // Verify the deployment using chain assertions.
        ChainAssertions.checkOptimismPortal({ _contracts: _proxies(), _cfg: cfg, _isProxy: true });
    }

    /// @notice Initialize the OptimismPortal2
    function initializeOptimismPortal2() public broadcast {
        console.log("Upgrading and initializing OptimismPortal2 proxy");
        address optimismPortalProxy = mustGetAddress("OptimismPortalProxy");
        address disputeGameFactoryProxy = mustGetAddress("DisputeGameFactoryProxy");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");

        address optimismPortal2;

        if (cfg.reuseDeployment()) {
            address savedAddress = jsonDeployment.readAddress(".OptimismPortal2");
            console.log("OptimismPortal2 address from JSON: %s", savedAddress);

            if (savedAddress != address(0) && verifyImplementationExists(savedAddress)) {
                optimismPortal2 = savedAddress;
                console.log("Using existing implementation from JSON");
            } else {
                console.log("Implementation from JSON not found on-chain, deploying new one");
                optimismPortal2 = address(
                    new OptimismPortal2{ salt: _implSalt() }({
                        _proofMaturityDelaySeconds: cfg.proofMaturityDelaySeconds(),
                        _disputeGameFinalityDelaySeconds: cfg.disputeGameFinalityDelaySeconds()
                    })
                );
                require(optimismPortal2 != address(0), "OptimismPortal2 deployment failed");
                save("OptimismPortal2", optimismPortal2);
                console.log("OptimismPortal2 deployed at %s", optimismPortal2);
            }
        } else {
            optimismPortal2 = mustGetAddress("OptimismPortal2");
        }

        _upgradeAndCallViaSafe({
            _proxy: payable(optimismPortalProxy),
            _implementation: optimismPortal2,
            _innerCallData: abi.encodeCall(
                OptimismPortal2.initialize,
                (
                    DisputeGameFactory(disputeGameFactoryProxy),
                    SystemConfig(systemConfigProxy),
                    SuperchainConfig(superchainConfigProxy),
                    GameType.wrap(uint32(cfg.respectedGameType()))
                )
            )
        });

        OptimismPortal2 portal = OptimismPortal2(payable(optimismPortalProxy));
        string memory version = portal.version();
        console.log("OptimismPortal2 version: %s", version);

        ChainAssertions.checkOptimismPortal2({ _contracts: _proxies(), _cfg: cfg, _isProxy: true });
    }

    function initializeProtocolVersions() public broadcast {
        console.log("Upgrading and initializing ProtocolVersions proxy");
        address protocolVersionsProxy = mustGetAddress("ProtocolVersionsProxy");
        address protocolVersions = mustGetAddress("ProtocolVersions");

        address finalSystemOwner = cfg.finalSystemOwner();
        uint256 requiredProtocolVersion = cfg.requiredProtocolVersion();
        uint256 recommendedProtocolVersion = cfg.recommendedProtocolVersion();

        _upgradeAndCallViaSafe({
            _proxy: payable(protocolVersionsProxy),
            _implementation: protocolVersions,
            _innerCallData: abi.encodeCall(
                ProtocolVersions.initialize,
                (
                    finalSystemOwner,
                    ProtocolVersion.wrap(requiredProtocolVersion),
                    ProtocolVersion.wrap(recommendedProtocolVersion)
                )
            )
        });

        ProtocolVersions versions = ProtocolVersions(protocolVersionsProxy);
        string memory version = versions.version();
        console.log("ProtocolVersions version: %s", version);

        ChainAssertions.checkProtocolVersions({ _contracts: _proxiesUnstrict(), _cfg: cfg, _isProxy: true });
    }

    /// @notice Transfer ownership of the DisputeGameFactory contract to the final system owner
    function transferDisputeGameFactoryOwnership() public broadcast {
        console.log("Transferring DisputeGameFactory ownership to Safe");
        DisputeGameFactory disputeGameFactory = DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy"));
        address owner = disputeGameFactory.owner();

        address safe = mustGetAddress("SystemOwnerSafe");
        if (owner != safe) {
            disputeGameFactory.transferOwnership(safe);
            console.log("DisputeGameFactory ownership transferred to Safe at: %s", safe);
        }
        ChainAssertions.checkDisputeGameFactory({ _contracts: _proxies(), _expectedOwner: safe });
    }

    /// @notice Transfer ownership of the DelayedWETH contract to the final system owner
    function transferDelayedWETHOwnership() public broadcast {
        console.log("Transferring DelayedWETH ownership to Safe");
        DelayedWETH weth = DelayedWETH(mustGetAddress("DelayedWETHProxy"));
        address owner = weth.owner();

        address safe = mustGetAddress("SystemOwnerSafe");
        if (owner != safe) {
            weth.transferOwnership(safe);
            console.log("DelayedWETH ownership transferred to Safe at: %s", safe);
        }
        ChainAssertions.checkDelayedWETH({ _contracts: _proxies(), _cfg: cfg, _isProxy: true, _expectedOwner: safe });
    }

    /// @notice Transfer ownership of the permissioned DelayedWETH contract to the final system owner
    function transferPermissionedDelayedWETHOwnership() public broadcast {
        console.log("Transferring permissioned DelayedWETH ownership to Safe");
        DelayedWETH weth = DelayedWETH(mustGetAddress("PermissionedDelayedWETHProxy"));
        address owner = weth.owner();

        address safe = mustGetAddress("SystemOwnerSafe");
        if (owner != safe) {
            weth.transferOwnership(safe);
            console.log("DelayedWETH ownership transferred to Safe at: %s", safe);
        }
        ChainAssertions.checkPermissionedDelayedWETH({
            _contracts: _proxies(),
            _cfg: cfg,
            _isProxy: true,
            _expectedOwner: safe
        });
    }

    /// @notice Loads the mips absolute prestate from the prestate-proof for devnets otherwise
    ///         from the config.
    function loadMipsAbsolutePrestate() internal returns (Claim mipsAbsolutePrestate_) {
        if (block.chainid == Chains.LocalDevnet || block.chainid == Chains.GethDevnet) {
            // Fetch the absolute prestate dump
            string memory filePath = string.concat(vm.projectRoot(), "/../../../op-program/bin/prestate-proof.json");
            string[] memory commands = new string[](3);
            commands[0] = "bash";
            commands[1] = "-c";
            commands[2] = string.concat("[[ -f ", filePath, " ]] && echo \"present\"");
            if (Process.run(commands).length == 0) {
                revert("Cannon prestate dump not found, generate it with `make cannon-prestate` in the monorepo root.");
            }
            commands[2] = string.concat("cat ", filePath, " | jq -r .pre");
            mipsAbsolutePrestate_ = Claim.wrap(abi.decode(Process.run(commands), (bytes32)));
            console.log(
                "[Cannon Dispute Game] Using devnet MIPS Absolute prestate: %s",
                vm.toString(Claim.unwrap(mipsAbsolutePrestate_))
            );
        } else {
            console.log(
                "[Cannon Dispute Game] Using absolute prestate from config: %x", cfg.faultGameAbsolutePrestate()
            );
            mipsAbsolutePrestate_ = Claim.wrap(bytes32(cfg.faultGameAbsolutePrestate()));
        }
    }

    /// @notice Sets the implementation for the `CANNON` game type in the `DisputeGameFactory`
    function setCannonFaultGameImplementation(bool _allowUpgrade) public broadcast {
        console.log("Setting Cannon FaultDisputeGame implementation");
        DisputeGameFactory factory = DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy"));
        DelayedWETH weth = DelayedWETH(mustGetAddress("DelayedWETHProxy"));

        // Set the Cannon FaultDisputeGame implementation in the factory.
        _setFaultGameImplementation({
            _factory: factory,
            _allowUpgrade: _allowUpgrade,
            _params: FaultDisputeGameParams({
                anchorStateRegistry: AnchorStateRegistry(mustGetAddress("AnchorStateRegistryProxy")),
                weth: weth,
                gameType: GameTypes.CANNON,
                absolutePrestate: loadMipsAbsolutePrestate(),
                faultVm: IBigStepper(mustGetAddress("Mips")),
                maxGameDepth: cfg.faultGameMaxDepth(),
                maxClockDuration: Duration.wrap(uint64(cfg.faultGameMaxClockDuration()))
            })
        });
    }

    /// @notice Sets the implementation for the `PERMISSIONED_CANNON` game type in the `DisputeGameFactory`
    function setPermissionedCannonFaultGameImplementation(bool _allowUpgrade) public broadcast {
        console.log("Setting Cannon PermissionedDisputeGame implementation");
        DisputeGameFactory factory = DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy"));
        DelayedWETH weth = DelayedWETH(mustGetAddress("PermissionedDelayedWETHProxy"));

        // Set the Cannon FaultDisputeGame implementation in the factory.
        _setFaultGameImplementation({
            _factory: factory,
            _allowUpgrade: _allowUpgrade,
            _params: FaultDisputeGameParams({
                anchorStateRegistry: AnchorStateRegistry(mustGetAddress("AnchorStateRegistryProxy")),
                weth: weth,
                gameType: GameTypes.PERMISSIONED_CANNON,
                absolutePrestate: loadMipsAbsolutePrestate(),
                faultVm: IBigStepper(mustGetAddress("Mips")),
                maxGameDepth: cfg.faultGameMaxDepth(),
                maxClockDuration: Duration.wrap(uint64(cfg.faultGameMaxClockDuration()))
            })
        });
    }

    /// @notice Sets the implementation for the `ALPHABET` game type in the `DisputeGameFactory`
    function setAlphabetFaultGameImplementation(bool _allowUpgrade) public onlyDevnet broadcast {
        console.log("Setting Alphabet FaultDisputeGame implementation");
        DisputeGameFactory factory = DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy"));
        DelayedWETH weth = DelayedWETH(mustGetAddress("DelayedWETHProxy"));

        Claim outputAbsolutePrestate = Claim.wrap(bytes32(cfg.faultGameAbsolutePrestate()));
        _setFaultGameImplementation({
            _factory: factory,
            _allowUpgrade: _allowUpgrade,
            _params: FaultDisputeGameParams({
                anchorStateRegistry: AnchorStateRegistry(mustGetAddress("AnchorStateRegistryProxy")),
                weth: weth,
                gameType: GameTypes.ALPHABET,
                absolutePrestate: outputAbsolutePrestate,
                faultVm: IBigStepper(new AlphabetVM(outputAbsolutePrestate, PreimageOracle(mustGetAddress("PreimageOracle")))),
                // The max depth for the alphabet trace is always 3. Add 1 because split depth is fully inclusive.
                maxGameDepth: cfg.faultGameSplitDepth() + 3 + 1,
                maxClockDuration: Duration.wrap(uint64(cfg.faultGameMaxClockDuration()))
            })
        });
    }

    /// @notice Sets the implementation for the `ALPHABET` game type in the `DisputeGameFactory`
    function setFastFaultGameImplementation(bool _allowUpgrade) public onlyDevnet broadcast {
        console.log("Setting Fast FaultDisputeGame implementation");
        DisputeGameFactory factory = DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy"));
        DelayedWETH weth = DelayedWETH(mustGetAddress("DelayedWETHProxy"));

        Claim outputAbsolutePrestate = Claim.wrap(bytes32(cfg.faultGameAbsolutePrestate()));
        PreimageOracle fastOracle = new PreimageOracle(cfg.preimageOracleMinProposalSize(), 0);
        _setFaultGameImplementation({
            _factory: factory,
            _allowUpgrade: _allowUpgrade,
            _params: FaultDisputeGameParams({
                anchorStateRegistry: AnchorStateRegistry(mustGetAddress("AnchorStateRegistryProxy")),
                weth: weth,
                gameType: GameTypes.FAST,
                absolutePrestate: outputAbsolutePrestate,
                faultVm: IBigStepper(new AlphabetVM(outputAbsolutePrestate, fastOracle)),
                // The max depth for the alphabet trace is always 3. Add 1 because split depth is fully inclusive.
                maxGameDepth: cfg.faultGameSplitDepth() + 3 + 1,
                maxClockDuration: Duration.wrap(0) // Resolvable immediately
             })
        });
    }

    /// @notice Sets the implementation for the given fault game type in the `DisputeGameFactory`.
    function _setFaultGameImplementation(
        DisputeGameFactory _factory,
        bool _allowUpgrade,
        FaultDisputeGameParams memory _params
    )
        internal
    {
        if (address(_factory.gameImpls(_params.gameType)) != address(0) && !_allowUpgrade) {
            console.log(
                "[WARN] DisputeGameFactoryProxy: `FaultDisputeGame` implementation already set for game type: %s",
                vm.toString(GameType.unwrap(_params.gameType))
            );
            return;
        }

        uint32 rawGameType = GameType.unwrap(_params.gameType);
        if (rawGameType != GameTypes.PERMISSIONED_CANNON.raw()) {
            _factory.setImplementation(
                _params.gameType,
                new FaultDisputeGame({
                    _gameType: _params.gameType,
                    _absolutePrestate: _params.absolutePrestate,
                    _maxGameDepth: _params.maxGameDepth,
                    _splitDepth: cfg.faultGameSplitDepth(),
                    _clockExtension: Duration.wrap(uint64(cfg.faultGameClockExtension())),
                    _maxClockDuration: _params.maxClockDuration,
                    _vm: _params.faultVm,
                    _weth: _params.weth,
                    _anchorStateRegistry: _params.anchorStateRegistry,
                    _l2ChainId: cfg.l2ChainID()
                })
            );
        } else {
            _factory.setImplementation(
                _params.gameType,
                new PermissionedDisputeGame({
                    _gameType: _params.gameType,
                    _absolutePrestate: _params.absolutePrestate,
                    _maxGameDepth: _params.maxGameDepth,
                    _splitDepth: cfg.faultGameSplitDepth(),
                    _clockExtension: Duration.wrap(uint64(cfg.faultGameClockExtension())),
                    _maxClockDuration: Duration.wrap(uint64(cfg.faultGameMaxClockDuration())),
                    _vm: _params.faultVm,
                    _weth: _params.weth,
                    _anchorStateRegistry: _params.anchorStateRegistry,
                    _l2ChainId: cfg.l2ChainID(),
                    _proposer: cfg.l2OutputOracleProposer(),
                    _challenger: cfg.l2OutputOracleChallenger()
                })
            );
        }

        string memory gameTypeString;
        if (rawGameType == GameTypes.CANNON.raw()) {
            gameTypeString = "Cannon";
        } else if (rawGameType == GameTypes.PERMISSIONED_CANNON.raw()) {
            gameTypeString = "PermissionedCannon";
        } else if (rawGameType == GameTypes.ALPHABET.raw()) {
            gameTypeString = "Alphabet";
        } else {
            gameTypeString = "Unknown";
        }

        console.log(
            "DisputeGameFactoryProxy: set `FaultDisputeGame` implementation (Backend: %s | GameType: %s)",
            gameTypeString,
            vm.toString(rawGameType)
        );
    }

    /// @notice Initialize the DataAvailabilityChallenge
    function initializeDataAvailabilityChallenge() public broadcast {
        console.log("Upgrading and initializing DataAvailabilityChallenge proxy");
        address dataAvailabilityChallengeProxy = mustGetAddress("DataAvailabilityChallengeProxy");
        address dataAvailabilityChallenge = mustGetAddress("DataAvailabilityChallenge");

        address finalSystemOwner = cfg.finalSystemOwner();
        uint256 daChallengeWindow = cfg.daChallengeWindow();
        uint256 daResolveWindow = cfg.daResolveWindow();
        uint256 daBondSize = cfg.daBondSize();
        uint256 daResolverRefundPercentage = cfg.daResolverRefundPercentage();

        _upgradeAndCallViaSafe({
            _proxy: payable(dataAvailabilityChallengeProxy),
            _implementation: dataAvailabilityChallenge,
            _innerCallData: abi.encodeCall(
                DataAvailabilityChallenge.initialize,
                (finalSystemOwner, daChallengeWindow, daResolveWindow, daBondSize, daResolverRefundPercentage)
            )
        });

        DataAvailabilityChallenge dac = DataAvailabilityChallenge(payable(dataAvailabilityChallengeProxy));
        string memory version = dac.version();
        console.log("DataAvailabilityChallenge version: %s", version);

        require(dac.owner() == finalSystemOwner);
        require(dac.challengeWindow() == daChallengeWindow);
        require(dac.resolveWindow() == daResolveWindow);
        require(dac.bondSize() == daBondSize);
        require(dac.resolverRefundPercentage() == daResolverRefundPercentage);
    }
}
