// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

contract MockL1StandardBridge {
  // Empty implementation
}

contract MockL1CrossDomainMessenger {
  // Empty implementation
}

contract MockOptimismPortal {
  // Empty implementation
}

contract MockBridgeRegistry {
  event RollupConfigRegistered(
    address rollupConfig,
    uint8 tokenType,
    address l2Token,
    string name
  );

  function registerRollupConfig(
    address rollupConfig,
    uint8 tokenType,
    address l2Token,
    string calldata name
  ) external {
    emit RollupConfigRegistered(rollupConfig, tokenType, l2Token, name);
  }

  function availableForRegistration(
    address /* rollupConfig */,
    uint8 /* tokenType */
  ) external pure returns (bool) {
    return true;
  }
}
