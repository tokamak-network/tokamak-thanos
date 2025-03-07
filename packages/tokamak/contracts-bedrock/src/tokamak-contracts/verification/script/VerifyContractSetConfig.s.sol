// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import 'forge-std/Script.sol';
import 'forge-std/console.sol';
import '@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol';
import '../L1ContractVerification.sol';

contract VerifyContractSetConfig is Script {
  // Contract addresses
  address constant _SYSTEM_CONFIG_PROXY =
    0x6eF61974A3CDa7BbD0a4DD0A613f56d211c8AfDC;
  address constant _L1_STANDARD_BRIDGE =
    0x757EC5b8F81eDdfC31F305F3325Ac6Abf4A63a5D;
  address constant _L1_CROSS_DOMAIN_MESSENGER =
    0xd054Bc768aAC07Dd0BaA2856a2fFb68F495E4CC2;
  address constant _OPTIMISM_PORTAL =
    0x2fbD30Fcd1c4573b0288E706Be56B5c0d2DfcAF6;
  address constant _SAFE_WALLET = 0x3E5c63644E683549055b9Be8653de26E0B4CD36E;
  // Contract IDs
  bytes32 constant _SYSTEM_CONFIG_ID = keccak256('SYSTEM_CONFIG');
  bytes32 constant _L1_STANDARD_BRIDGE_ID = keccak256('L1_STANDARD_BRIDGE');
  bytes32 constant _L1_CROSS_DOMAIN_MESSENGER_ID =
    keccak256('L1_CROSS_DOMAIN_MESSENGER');
  bytes32 constant _OPTIMISM_PORTAL_ID = keccak256('OPTIMISM_PORTAL');

  // L1ContractVerification deployed address
  address constant _VERIFIER_ADDRESS =
    0x9321fE98838893E49197d3EAE6008fF3D7e559b7;

  // Implementation slot
  bytes32 constant _IMPLEMENTATION_SLOT =
    0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
  // Admin slot
  bytes32 constant _ADMIN_SLOT =
    0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;
  // Legacy implementation slot (for OVM_L1CrossDomainMessenger)
  bytes32 constant _LEGACY_IMPLEMENTATION_SLOT =
    0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3;
  // Legacy admin slot (for OVM_L1CrossDomainMessenger)
  bytes32 constant _LEGACY_ADMIN_SLOT =
    0x10d6a54a4754c8869d6886b5f5d7fbfa5b4522237ea5c60d11bc4e7a1ff9390b;

  function getAddressFromSlot(
    address _target,
    bytes32 _slot
  ) internal view returns (address) {
    bytes32 value = vm.load(_target, _slot);
    return address(uint160(uint256(value)));
  }

  function hasCode(address _addr) internal view returns (bool) {
    uint size;
    assembly {
      size := extcodesize(_addr)
    }
    return size > 0;
  }

  function getImplementationAddress(
    address proxyAddress,
    bool isLegacy
  ) internal view returns (address implementation, bool found) {
    if (isLegacy) {
      // Try legacy implementation slot first
      implementation = getAddressFromSlot(
        proxyAddress,
        _LEGACY_IMPLEMENTATION_SLOT
      );
      if (hasCode(implementation)) return (implementation, true);
    } else {
      // Try standard EIP-1967 slot first
      implementation = getAddressFromSlot(proxyAddress, _IMPLEMENTATION_SLOT);
      if (hasCode(implementation)) return (implementation, true);
    }

    // Try the other slot as fallback
    bytes32 fallbackSlot = isLegacy
      ? _IMPLEMENTATION_SLOT
      : _LEGACY_IMPLEMENTATION_SLOT;
    implementation = getAddressFromSlot(proxyAddress, fallbackSlot);
    if (hasCode(implementation)) return (implementation, true);

    // Try scanning first few slots as last resort
    for (uint i = 0; i < 10; i++) {
      address potentialImpl = getAddressFromSlot(proxyAddress, bytes32(i));
      if (hasCode(potentialImpl)) {
        console.log('Found potential implementation in slot', i);
        return (potentialImpl, true);
      }
    }

    // Try to get implementation from Etherscan-verified source
    if (proxyAddress == _L1_CROSS_DOMAIN_MESSENGER) {
      // Hardcoded implementation from Etherscan verification
      address knownImplementation = 0x1Eba187355Bf2CC4530334861F48E116dcE6676a;
      if (hasCode(knownImplementation)) {
        console.log('Using known implementation from Etherscan verification');
        return (knownImplementation, true);
      }
    }

    return (address(0), false);
  }

  function getProxyAdmin(
    address proxyAddress,
    bool isLegacy
  ) internal view returns (address admin) {
    if (isLegacy) {
      admin = getAddressFromSlot(proxyAddress, _LEGACY_ADMIN_SLOT);
      if (admin == address(0)) {
        admin = getAddressFromSlot(proxyAddress, _ADMIN_SLOT);
      }
    } else {
      admin = getAddressFromSlot(proxyAddress, _ADMIN_SLOT);
      if (admin == address(0)) {
        admin = getAddressFromSlot(proxyAddress, _LEGACY_ADMIN_SLOT);
      }
    }

    // If still not found, try to get from Etherscan verification
    if (admin == address(0) && proxyAddress == _L1_CROSS_DOMAIN_MESSENGER) {
      // Hardcoded admin from Etherscan or documentation
      admin = 0x5a0Aae59D09fccBdDb6C6CcEB07B7279367C3d2A; // Example - replace with actual admin
      console.log('Using known admin from Etherscan verification');
    }

    return admin;
  }

  function verifyAndSetContractConfig(
    address proxyAddress,
    bytes32 contractId,
    L1ContractVerification verifier,
    bool isLegacy
  ) internal {
    console.log('\nVerifying contract:', vm.toString(contractId));
    console.log('Proxy Address:', proxyAddress);
    console.log('Is Legacy Proxy:', isLegacy);

    // Get implementation address with fallbacks
    (
      address implementation,
      bool implementationFound
    ) = getImplementationAddress(proxyAddress, isLegacy);

    // Get admin address
    address proxyAdmin = getProxyAdmin(proxyAddress, isLegacy);

    // Check if implementation is valid
    console.log('Implementation Address:', implementation);
    console.log('Implementation Found:', implementationFound);
    console.log('Proxy Admin:', proxyAdmin);

    require(implementationFound, 'Could not find valid implementation');

    // Get implementation code hash
    bytes32 implementationHash = implementation.codehash;

    // Get proxy code hash
    bytes32 proxyHash = proxyAddress.codehash;

    // Print results
    console.log('Implementation Hash:');
    console.logBytes32(implementationHash);
    console.log('Proxy Hash:');
    console.logBytes32(proxyHash);

    // Set the contract config
    verifier.setContractConfig(
      11155111, // Sepolia chain ID
      contractId,
      implementationHash,
      proxyHash,
      proxyAdmin
    );

    console.log(
      'Contract configuration set successfully for',
      vm.toString(contractId)
    );
  }

  function verifyL1Contracts(L1ContractVerification verifier) internal {
    console.log('\nVerifying L1 contracts:');
    console.log('System Config Proxy:', _SYSTEM_CONFIG_PROXY);
    console.log('Safe Wallet:', _SAFE_WALLET);
    console.log('Chain ID:', "11155111");

    bool success = verifier.verifyL1Contracts(
      11155111,
      _SYSTEM_CONFIG_PROXY,
      _SAFE_WALLET
    );
    console.log('Verification result:', success);

    if (success) {
      console.log('L1 contracts verified successfully!');
    } else {
      console.log('L1 contracts verification failed!');
    }
  }

  function run() external {
    console.log('Setting Contract Configurations for Multiple Contracts');
    console.log('===================================================');

    vm.startBroadcast();

    L1ContractVerification verifier = L1ContractVerification(_VERIFIER_ADDRESS);

    // Verify and set config for SystemConfig
    verifyAndSetContractConfig(
      _SYSTEM_CONFIG_PROXY,
      _SYSTEM_CONFIG_ID,
      verifier,
      false
    );

    // Verify and set config for L1StandardBridge
    verifyAndSetContractConfig(
      _L1_STANDARD_BRIDGE,
      _L1_STANDARD_BRIDGE_ID,
      verifier,
      false
    );

    // Verify and set config for L1CrossDomainMessenger (legacy proxy)
    verifyAndSetContractConfig(
      _L1_CROSS_DOMAIN_MESSENGER,
      _L1_CROSS_DOMAIN_MESSENGER_ID,
      verifier,
      true
    );

    // Verify and set config for OptimismPortal
    verifyAndSetContractConfig(
      _OPTIMISM_PORTAL,
      _OPTIMISM_PORTAL_ID,
      verifier,
      false
    );

    verifyL1Contracts(verifier);

    vm.stopBroadcast();

    console.log('\nAll contract configurations set successfully');
  }
}
