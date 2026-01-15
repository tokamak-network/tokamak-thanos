// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {UpgradeL1BridgeV1} from '../../src/shutdown/UpgradeL1BridgeV1.sol';

/**
 * @title ActivateForceWithdraw
 * @notice Forge script to activate or deactivate force withdrawal functionality.
 *
 * Usage:
 * forge script scripts/shutdown/ActivateForceWithdraw.s.sol --rpc-url <rpc> --broadcast --sig "run(address,bool)" <bridge_proxy> <true/false>
 */
contract ActivateForceWithdraw is Script {
  function run(address _bridgeProxy, bool _shouldActivate) public {
    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');
    vm.startBroadcast(deployerPrivateKey);

    UpgradeL1BridgeV1 bridge = UpgradeL1BridgeV1(payable(_bridgeProxy));

    console.log('------------------------------------------');
    console.log('[INFO] Force Withdrawal Activation/Deactivation');
    console.log('Bridge Proxy:', _bridgeProxy);
    console.log('Target State:', _shouldActivate ? 'ACTIVATE' : 'DEACTIVATE');

    bool currentState = bridge.active();
    console.log('Current Active State:', currentState ? 'ACTIVE' : 'INACTIVE');

    if (currentState == _shouldActivate) {
      console.log('[INFO] Bridge is already in the target state. Skipping.');
    } else {
      console.log('[ACTION] Calling forceActive(%s)...', _shouldActivate);
      bridge.forceActive(_shouldActivate);
      console.log('[SUCCESS] Done!');
    }

    console.log('------------------------------------------');

    vm.stopBroadcast();
  }
}
