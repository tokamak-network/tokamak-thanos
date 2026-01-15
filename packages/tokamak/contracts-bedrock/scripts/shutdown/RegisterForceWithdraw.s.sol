// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {UpgradeL1BridgeV1} from '../../src/shutdown/UpgradeL1BridgeV1.sol';

/**
 * @title RegisterForceWithdraw
 * @notice Forge script to register GenFWStorage position contracts in the Bridge.
 *
 * Usage:
 * forge script scripts/shutdown/RegisterForceWithdraw.s.sol --rpc-url <rpc> --broadcast --sig "run(address,string)" <bridge_proxy> <path_to_positions_json>
 */
contract RegisterForceWithdraw is Script {
  using stdJson for string;

  function run(address _bridgeProxy, string memory _positionsPath) public {
    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');

    // Read JSON file
    string memory json = vm.readFile(_positionsPath);
    address[] memory positions = json.readAddressArray('$');

    console.log('------------------------------------------');
    console.log('[INFO] Bridge Position Registration');
    console.log('Bridge Proxy:', _bridgeProxy);
    console.log('Loading positions from:', _positionsPath);
    console.log('Loaded %s position(s)', positions.length);

    vm.startBroadcast(deployerPrivateKey);

    UpgradeL1BridgeV1 bridge = UpgradeL1BridgeV1(payable(_bridgeProxy));

    console.log('\n[ACTION] Calling forceRegistry with positions...');
    bridge.forceRegistry(positions);

    console.log('\n[INFO] Verifying registration status...');
    bool allGood = true;
    for (uint i = 0; i < positions.length; i++) {
      bool isRegistered = bridge.position(positions[i]);
      if (isRegistered) {
        console.log('   Registration success for:', positions[i]);
      } else {
        console.log('   Registration failed for:', positions[i]);
        allGood = false;
      }
    }

    if (!allGood) {
      revert('Registration verification failed!');
    }

    console.log('\n[SUCCESS] All positions successfully registered');
    console.log('------------------------------------------');

    vm.stopBroadcast();
  }
}
