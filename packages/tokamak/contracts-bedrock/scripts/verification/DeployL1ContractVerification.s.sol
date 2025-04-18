// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import '@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol';

/**
 * @title SetupL1ContractVerification
 * @notice Script to deploy and configure the L1ContractVerification contract as upgradeable
 * @dev This script is used by the TRH team to deploy and set up the verification contract
 *      before releasing the trh-sdk
 */
contract SetupL1ContractVerification is Script {
  // Environment variables - private variables with leading underscore
  address private _systemConfigProxy;
  address private _tokamakDAO;
  address private _foundation;
  address private _nativeToken; // TON token address
  address private _bridgeRegistry;
  address private _multisigWallet;
  address private _l1ProxyAdmin;

  // Constants - private constants with leading underscore
  uint256 private constant _SAFE_THRESHOLD = 3;

  /**
   * @notice Load configuration from environment variables
   * @dev Called before the main execution to set up required addresses
   */
  function setUp() public {
    // Load environment variables for contract addresses
    _systemConfigProxy = vm.envAddress('SYSTEM_CONFIG_PROXY');
    _tokamakDAO = vm.envAddress('TOKAMAK_DAO_ADDRESS');
    _foundation = vm.envAddress('FOUNDATION_ADDRESS');
    _nativeToken = vm.envAddress('NATIVE_TOKEN_ADDRESS');
    _bridgeRegistry = vm.envAddress('L1_BRIDGE_REGISTRY_ADDRESS');
    _multisigWallet = vm.envAddress('MULTISIG_WALLET');
    _l1ProxyAdmin = vm.envAddress('PROXY_ADMIN_ADDRESS');
  }

  /**
   * @notice Main execution function
   * @dev Deploys and configures the L1ContractVerification contract with a TransparentUpgradeableProxy
   */
  function run() external {
    setUp();

    // Start broadcasting transactions
    vm.startBroadcast();

    // Deploy the implementation contract
    L1ContractVerification implementationContract = new L1ContractVerification();
    console.log('L1ContractVerification implementation deployed at:', address(implementationContract));

    // Deploy a new ProxyAdmin
    ProxyAdmin verificationContractProxyAdmin = new ProxyAdmin();
    console.log('verificationContractProxyAdmin deployed at:', address(verificationContractProxyAdmin));

    // Prepare initialization data for the proxy
    bytes memory initData = abi.encodeWithSelector(
      L1ContractVerification.initialize.selector,
      _nativeToken,
      msg.sender
    );

    // Deploy the TransparentUpgradeableProxy with the new ProxyAdmin
    TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
      address(implementationContract),
      address(verificationContractProxyAdmin),
      initData
    );
    console.log('L1ContractVerification proxy deployed at:', address(proxy));

    // Create a reference to the proxy contract to call its functions
    L1ContractVerification verifier = L1ContractVerification(address(proxy));

    // Set logic contract info
    verifier.setLogicContractInfo(_systemConfigProxy, _l1ProxyAdmin);

    // Get the safe wallet address from the proxy admin
    IProxyAdmin proxyAdminContract = IProxyAdmin(address(_l1ProxyAdmin));
    address safeWalletAddress = proxyAdminContract.owner();

    console.log('Safe wallet address:', safeWalletAddress);
    // Get the implementation address of the safe wallet using masterCopy function
    address implementation = IGnosisSafe(safeWalletAddress).masterCopy();

    // Set Safe wallet info with the implementation address and codehash
    verifier.setSafeConfig(
      _tokamakDAO,
      _foundation,
      _SAFE_THRESHOLD,
      implementation.codehash,
      safeWalletAddress.codehash,
      3
    );

    // Set bridge registry address
    verifier.setBridgeRegistryAddress(_bridgeRegistry);

    // Enable verification
    verifier.setVerificationPossible(true);
    console.log('Verification enabled');

    // Now transfer admin role to multisig
    verifier.grantRole(verifier.ADMIN_ROLE(), _multisigWallet);
    console.log('Admin role granted to multisig wallet:', _multisigWallet);

    // Transfer ownership of the ProxyAdmin to the multisig wallet
    verificationContractProxyAdmin.transferOwnership(_multisigWallet);
    console.log('verificationContractProxyAdmin ownership transferred to multisig wallet:', _multisigWallet);

    address currentOwner = verificationContractProxyAdmin.owner();
    console.log('Current ProxyAdmin owner:', currentOwner);
    console.log('Expected multisig address:', _multisigWallet);

    require(currentOwner == _multisigWallet, "ProxyAdmin not owned by multisig!");

    console.log('L1ContractVerification configuration complete');
    vm.stopBroadcast();
  }
}
