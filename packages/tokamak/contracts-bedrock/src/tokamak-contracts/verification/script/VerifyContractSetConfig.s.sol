// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import '../L1ContractVerification.sol';

contract VerifyContractSetConfig is Script {
  // Contract IDs
  bytes32 constant SYSTEM_CONFIG_ID = keccak256('SYSTEM_CONFIG');
  bytes32 constant L1_STANDARD_BRIDGE_ID = keccak256('L1_STANDARD_BRIDGE');
  bytes32 constant L1_CROSS_DOMAIN_MESSENGER_ID = keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 constant OPTIMISM_PORTAL_ID = keccak256('OPTIMISM_PORTAL');

  // Environment variables
  uint256 chainId;
  address systemConfigProxy;
  address l1StandardBridge;
  address l1CrossDomainMessenger;
  address optimismPortal;
  address safeWallet;
  address systemConfigImpl;
  address l1StandardBridgeImpl;
  address l1CrossDomainMessengerImpl;
  address optimismPortalImpl;
  address nativeTokenAddress;
  address proxyAdmin;
  address verifierAddress;

  function setUp() public {
    // Load environment variables
    chainId = vm.envUint("CHAIN_ID");
    systemConfigProxy = vm.envAddress("SYSTEM_CONFIG_PROXY");
    l1StandardBridge = vm.envAddress("L1_STANDARD_BRIDGE");
    l1CrossDomainMessenger = vm.envAddress("L1_CROSS_DOMAIN_MESSENGER");
    optimismPortal = vm.envAddress("OPTIMISM_PORTAL");
    safeWallet = vm.envAddress("SAFE_WALLET");
    systemConfigImpl = vm.envAddress("SYSTEM_CONFIG_IMPL");
    l1StandardBridgeImpl = vm.envAddress("L1_STANDARD_BRIDGE_IMPL");
    l1CrossDomainMessengerImpl = vm.envAddress("L1_CROSS_DOMAIN_MESSENGER_IMPL");
    optimismPortalImpl = vm.envAddress("OPTIMISM_PORTAL_IMPL");
    nativeTokenAddress = vm.envAddress("NATIVE_TOKEN_ADDRESS");
    proxyAdmin = vm.envAddress("PROXY_ADMIN_CONTRACT_ADDRESS");
    verifierAddress = vm.envAddress("VERIFIER_ADDRESS");
  }

  function setSafeConfig(L1ContractVerification verifier) internal {
    verifier.setSafeConfig(
      safeWallet,  // TokamakDAO address (using Safe address for now)
      safeWallet,  // Foundation address (using Safe address for now)
      1            // Threshold (matches the actual threshold)
    );
  }

  function verifyAndSetContractConfig(
    address proxyAddress,
    address implAddress,
    bytes32 contractId,
    L1ContractVerification verifier
  ) internal {
    bytes32 implementationHash = implAddress.codehash;
    bytes32 proxyHash = proxyAddress.codehash;

    verifier.setContractConfig(
      contractId,
      implementationHash,
      proxyHash,
      proxyAdmin
    );
  }

  function verifyL1Contracts(L1ContractVerification verifier) internal returns (bool) {
    return verifier.verifyL1Contracts(
      systemConfigProxy,
      safeWallet
    );
  }

  function run() external {
    setUp();
    vm.startBroadcast();

    L1ContractVerification verifier = L1ContractVerification(verifierAddress);

    // Configure the verifier
    setSafeConfig(verifier);
    verifier.setSafeVerificationRequired(chainId, false);
    verifier.setExpectedNativeToken(chainId, nativeTokenAddress);

    // Set contract configurations
    verifyAndSetContractConfig(
      systemConfigProxy,
      systemConfigImpl,
      SYSTEM_CONFIG_ID,
      verifier
    );

    verifyAndSetContractConfig(
      l1StandardBridge,
      l1StandardBridgeImpl,
      L1_STANDARD_BRIDGE_ID,
      verifier
    );

    verifyAndSetContractConfig(
      l1CrossDomainMessenger,
      l1CrossDomainMessengerImpl,
      L1_CROSS_DOMAIN_MESSENGER_ID,
      verifier
    );

    verifyAndSetContractConfig(
      optimismPortal,
      optimismPortalImpl,
      OPTIMISM_PORTAL_ID,
      verifier
    );

    // Verify L1 contracts
    bool success = verifyL1Contracts(verifier);
    console.log('Verification result:', success);

    vm.stopBroadcast();
  }
}
