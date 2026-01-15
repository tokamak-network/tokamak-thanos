// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {UpgradeL1BridgeV1} from '../../src/shutdown/UpgradeL1BridgeV1.sol';

contract DeployUpgradeL1BridgeV1 is Script {
  function run() external {
    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');
    vm.startBroadcast(deployerPrivateKey);

    UpgradeL1BridgeV1 impl = new UpgradeL1BridgeV1();
    console.log('Deployed to:', address(impl));

    vm.stopBroadcast();
  }
}
