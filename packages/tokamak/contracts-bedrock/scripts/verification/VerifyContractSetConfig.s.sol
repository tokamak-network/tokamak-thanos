// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import 'src/tokamak-contracts/verification/L1ContractVerification.sol';

contract VerifyContractSetConfig is Script {
  // Contract IDs
  bytes32 constant SYSTEM_CONFIG_ID = keccak256('SYSTEM_CONFIG');
  bytes32 constant L1_STANDARD_BRIDGE_ID = keccak256('L1_STANDARD_BRIDGE');
  bytes32 constant L1_CROSS_DOMAIN_MESSENGER_ID =
    keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 constant OPTIMISM_PORTAL_ID = keccak256('OPTIMISM_PORTAL');

  // Environment variables
  address systemConfigProxy;
  address l1StandardBridge;
  address l1CrossDomainMessenger;
  address optimismPortal;
  address systemConfigImpl;
  address l1StandardBridgeImpl;
  address l1CrossDomainMessengerImpl;
  address optimismPortalImpl;
  address nativeTokenAddress;
  address proxyAdmin;
  address verifierAddress;
  address tokamakDAO;
  address foundation;
  address thirdOwner;
  address bridgeAddress;
  address l2TonAddress;
  uint256 threshold;

  function setUp() public {
    // Load environment variables
    systemConfigProxy = vm.envAddress('SYSTEM_CONFIG_PROXY');
    l1StandardBridge = vm.envAddress('L1_STANDARD_BRIDGE');
    l1CrossDomainMessenger = vm.envAddress('L1_CROSS_DOMAIN_MESSENGER');
    optimismPortal = vm.envAddress('OPTIMISM_PORTAL');
    systemConfigImpl = vm.envAddress('SYSTEM_CONFIG_IMPL');
    l1StandardBridgeImpl = vm.envAddress('L1_STANDARD_BRIDGE_IMPL');
    l1CrossDomainMessengerImpl = vm.envAddress(
      'L1_CROSS_DOMAIN_MESSENGER_IMPL'
    );
    optimismPortalImpl = vm.envAddress('OPTIMISM_PORTAL_IMPL');
    nativeTokenAddress = vm.envAddress('NATIVE_TOKEN_ADDRESS');
    proxyAdmin = vm.envAddress('PROXY_ADMIN_CONTRACT_ADDRESS');
    verifierAddress = vm.envAddress('VERIFIER_ADDRESS');
    tokamakDAO = vm.envAddress('TOKAMAK_DAO_ADDRESS');
    foundation = vm.envAddress('FOUNDATION_ADDRESS');
    thirdOwner = vm.envAddress('THIRD_OWNER_ADDRESS');
    bridgeAddress = vm.envAddress('L1_BRIDGE_REGISTRY_ADDRESS');
    l2TonAddress = vm.envAddress('L2_TON_ADDRESS');
    threshold = vm.envUint('THRESHOLD');
  }

  function setSafeConfig(L1ContractVerification verifier) internal {
    verifier.setSafeConfig(tokamakDAO, foundation, thirdOwner, threshold);
  }

  function setContractConfig(
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

  function verifyL1Contracts(
    L1ContractVerification verifier
  ) internal returns (bool) {
    return verifier.verifyL1Contracts(systemConfigProxy);
  }

  function run() external {
    setUp();
    vm.startBroadcast();

    L1ContractVerification verifier = L1ContractVerification(verifierAddress);

    // Configure the verifier
    setSafeConfig(verifier);
    verifier.setSafeVerificationRequired(true);
    verifier.setExpectedNativeToken(nativeTokenAddress);

    // Set contract configurations
    setContractConfig(
      systemConfigProxy,
      systemConfigImpl,
      SYSTEM_CONFIG_ID,
      verifier
    );

    setContractConfig(
      l1StandardBridge,
      l1StandardBridgeImpl,
      L1_STANDARD_BRIDGE_ID,
      verifier
    );

    setContractConfig(
      l1CrossDomainMessenger,
      l1CrossDomainMessengerImpl,
      L1_CROSS_DOMAIN_MESSENGER_ID,
      verifier
    );

    setContractConfig(
      optimismPortal,
      optimismPortalImpl,
      OPTIMISM_PORTAL_ID,
      verifier
    );

    // Verify L1 contracts
    bool success = verifyL1Contracts(verifier);
    console.log('Verification result:', success);

    verifier.setBridgeRegistryAddress(bridgeAddress);

    // Verify and setup config
    bool verifyAndRegisterRollupConfigResult = verifier
      .verifyAndRegisterRollupConfig(
        systemConfigProxy,
        2,
        l2TonAddress,
        'TestRollup'
      );
    console.log(
      'Verify and register rollup config result:',
      verifyAndRegisterRollupConfigResult
    );

    vm.stopBroadcast();
  }
}
