// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {UpgradeL1BridgeV1} from '../../src/shutdown/UpgradeL1BridgeV1.sol';

/**
 * @title ExecuteForceWithdraw
 * @notice Forge script to execute force withdrawal claims.
 */
contract ExecuteForceWithdraw is Script {
  using stdJson for string;

  struct Withdrawal {
    address target;
    uint256 amount;
  }
  // ... necessary fields

  function run(
    address _bridgeProxy,
    string memory _assetsPath,
    address _positionAddr
  ) public {
    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');

    vm.startBroadcast(deployerPrivateKey);

    UpgradeL1BridgeV1 bridge = UpgradeL1BridgeV1(payable(_bridgeProxy));

    // TODO: Logic to read individual claim data from JSON and execute iteratively
    // forge script is suitable for sending a large number of transactions in batches

    console.log('Executing force withdrawals via bridge:', _bridgeProxy);
    console.log('Using position contract:', _positionAddr);

    vm.stopBroadcast();
  }
}
