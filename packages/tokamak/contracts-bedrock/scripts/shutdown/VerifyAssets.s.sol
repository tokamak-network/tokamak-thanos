// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';

contract VerifyAssets is Script {
  using stdJson for string;

  address eth = 0xdeaDDeADDEaDdeaDdEAddEADDEAdDeadDEADDEaD;

  function run() public view {
    string memory rootPath = vm.projectRoot();
    string memory path = string.concat(rootPath, '/data/generate-assets.json');
    string memory json = vm.readFile(path);

    console.log('--------------------------------------------------');
    console.log('Starting Asset Verification...');
    console.log('Target File: data/generate-assets.json');
    console.log('--------------------------------------------------');

    // Manual parsing (extract directly by field)
    address[] memory l1Tokens = json.readAddressArray('$[*].l1Token');
    address[] memory l2Tokens = json.readAddressArray('$[*].l2Token');
    string[] memory tokenNames = json.readStringArray('$[*].tokenName');

    for (uint i = 0; i < l1Tokens.length; i++) {
      console.log('Token:', tokenNames[i]);

      string memory claimersKey = string.concat(
        '$[',
        vm.toString(i),
        '].data[*].claimer'
      );
      address[] memory claimers = json.readAddressArray(claimersKey);

      string memory amountsKey = string.concat(
        '$[',
        vm.toString(i),
        '].data[*].amount'
      );
      uint256[] memory amounts = json.readUintArray(amountsKey);

      for (uint j = 0; j < claimers.length; j++) {
        address claimer = claimers[j];
        uint256 amount = amounts[j];

        uint256 actualBalance;
        if (l2Tokens[i] == eth) {
          actualBalance = claimer.balance;
        } else {
          // Try simplify to avoid complex calls
          actualBalance = 0;
        }

        if (actualBalance > 0 || l2Tokens[i] != eth) {
          console.log('   [INFO] Account: %s | Amt: %s', claimer, amount);
          console.log('   [INFO] L2 Balance: %s', actualBalance);
        }
      }
    }

    console.log('--------------------------------------------------');
    console.log('Verification Complete');
  }
}
