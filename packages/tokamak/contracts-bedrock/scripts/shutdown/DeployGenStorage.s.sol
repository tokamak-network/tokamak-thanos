// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {GenFWStorage1} from '../../test/shutdown/GenFWStorage1.sol';

/**
 * @title DeployGenStorage
 * @notice Forge script to deploy GenFWStorage contracts with constant hashes.
 *
 * Usage:
 * forge script scripts/shutdown/DeployGenStorage.s.sol --rpc-url <rpc> --broadcast --sig "run(string)" <path_to_assets_json>
 */
contract DeployGenStorage is Script {
  using stdJson for string;

  function run(string memory _assetsPath) public returns (address[] memory) {
    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');

    // Read asset JSON (data logic needed - actually requires a hash list)
    // Here, as an example, deploy only one or follow hardcoded logic

    vm.startBroadcast(deployerPrivateKey);

    // TODO: Logic to deploy as many as needed based on actual asset data
    // Currently, this is an example of deploying only the first GenFWStorage1
    GenFWStorage1 storage1 = new GenFWStorage1();

    address[] memory deployed = new address[](1);
    deployed[0] = address(storage1);

    console.log('Deployed GenFWStorage1 at:', address(storage1));

    vm.stopBroadcast();
    return deployed;
  }
}
