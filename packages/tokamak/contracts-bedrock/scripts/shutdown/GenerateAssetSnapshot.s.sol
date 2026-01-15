// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {VmSafe} from 'forge-std/Vm.sol';

interface IERC20 {
  function balanceOf(address account) external view returns (uint256);
  function name() external view returns (string memory);
}

interface IERC721 {
  function balanceOf(address owner) external view returns (uint256);
  function name() external view returns (string memory);
}

interface IL2StandardToken {
  function l1Token() external view returns (address);
}

interface INonfungiblePositionManager {
  function balanceOf(address owner) external view returns (uint256);
  function tokenOfOwnerByIndex(
    address owner,
    uint256 index
  ) external view returns (uint256);
  function positions(
    uint256 tokenId
  )
    external
    view
    returns (
      uint96 nonce,
      address operator,
      address token0,
      address token1,
      uint24 fee,
      int24 tickLower,
      int24 tickUpper,
      uint128 liquidity,
      uint256 feeGrowthGlobal0X128,
      uint256 feeGrowthGlobal1X128,
      uint128 tokensOwed0,
      uint128 tokensOwed1
    );
}

contract GenerateAssetSnapshot is Script {
  using stdJson for string;

  // L2 System Addresses
  address l2StandardBridge = 0x4200000000000000000000000000000000000010;
  address l2ERC721Bridge = 0x4200000000000000000000000000000000000014;
  address messagePasser = 0x4200000000000000000000000000000000000016;

  // Token Sentinel & Specifics
  address nativeToken = 0xdeaDDeADDEaDdeaDdEAddEADDEAdDeadDEADDEaD;
  address ethToken = 0x4200000000000000000000000000000000000486;
  address wethToken = 0x4200000000000000000000000000000000000006;
  address usdcToken = 0x4200000000000000000000000000000000000778;

  uint256 l2StartBlock;
  uint256 l2EndBlock;

  mapping(address => address) l2ToL1;
  mapping(address => string) l2TokenNames;
  mapping(address => string) tokenType;
  mapping(address => address[]) tokenUsers;
  address[] uniqueL2Tokens;

  address[] discoveredUsers;
  mapping(address => bool) isUserKnown;

  function setUp() public {
    l2StartBlock = vm.envUint('L2_START_BLOCK');
    l2EndBlock = vm.envOr('L2_END_BLOCK', block.number);
  }

  function _addToken(
    address l1,
    address l2,
    string memory name,
    string memory tType
  ) internal {
    for (uint i = 0; i < uniqueL2Tokens.length; i++) {
      if (uniqueL2Tokens[i] == l2) return;
    }
    l2ToL1[l2] = l1;
    l2TokenNames[l2] = name;
    tokenType[l2] = tType;
    uniqueL2Tokens.push(l2);
  }

  function _addUser(address user) internal {
    if (user == address(0) || user.code.length > 0) return;
    // Exclude system addresses
    if (
      user == l2StandardBridge ||
      user == l2ERC721Bridge ||
      user == messagePasser
    ) return;

    if (!isUserKnown[user]) {
      isUserKnown[user] = true;
      discoveredUsers.push(user);
    }
  }

  function discoverAllUsersAndTokens() internal {
    console.log(
      'Global Scan: Discovering all active addresses from bridge history...'
    );

    bytes32 withdrawalEvent = keccak256(
      'WithdrawalInitiated(address,address,address,address,uint256,bytes)'
    );
    bytes32 depositEvent = keccak256(
      'DepositFinalized(address,address,address,address,uint256,bytes)'
    );

    // Fetch all events from the bridge without topic filters
    bytes32[] memory topics = new bytes32[](0);
    VmSafe.EthGetLogs[] memory logs = vm.eth_getLogs(
      l2StartBlock,
      l2EndBlock,
      l2StandardBridge,
      topics
    );

    for (uint i = 0; i < logs.length; i++) {
      bytes32 sig = logs[i].topics[0];
      if (sig == withdrawalEvent || sig == depositEvent) {
        // topics: [sig, l1Token, l2Token, from]
        address l2T = address(uint160(uint256(logs[i].topics[2])));
        address from = address(uint160(uint256(logs[i].topics[3])));

        // data: [to, amount, extraData]
        address to = abi.decode(logs[i].data, (address));

        _addUser(from);
        _addUser(to);

        if (l2T == nativeToken) {
          _addToken(address(0), l2T, 'Tokamak Network', 'ERC20');
        } else {
          try IERC20(l2T).name() returns (string memory name) {
            address l1T = address(0);
            try IL2StandardToken(l2T).l1Token() returns (address addr) {
              l1T = addr;
            } catch {}
            _addToken(l1T, l2T, name, 'ERC20');
          } catch {}
        }
      }
    }

    // Scan MessagePasser logs for user discovery
    bytes32 msgPassed = keccak256(
      'MessagePassed(uint256,address,address,uint256,uint256,bytes)'
    );
    topics = new bytes32[](0);
    logs = vm.eth_getLogs(l2StartBlock, l2EndBlock, messagePasser, topics);
    for (uint i = 0; i < logs.length; i++) {
      if (logs[i].topics[0] == msgPassed) {
        address sender = address(uint160(uint256(logs[i].topics[2])));
        _addUser(sender);
      }
    }

    console.log('   Discovered unique users:', discoveredUsers.length);
    console.log('   Discovered unique tokens:', uniqueL2Tokens.length);
  }

  function collectBalances() internal {
    console.log('Collecting balances across discovered users...');
    for (uint i = 0; i < uniqueL2Tokens.length; i++) {
      address l2T = uniqueL2Tokens[i];
      string memory tType = tokenType[l2T];
      for (uint j = 0; j < discoveredUsers.length; j++) {
        address user = discoveredUsers[j];
        uint256 bal = 0;

        if (l2T == nativeToken) {
          bal = user.balance;
        } else if (keccak256(bytes(tType)) == keccak256('ERC20')) {
          try IERC20(l2T).balanceOf(user) returns (uint256 b) {
            bal = b;
          } catch {}
        } else if (keccak256(bytes(tType)) == keccak256('ERC721')) {
          try IERC721(l2T).balanceOf(user) returns (uint256 b) {
            bal = b;
          } catch {}
        }

        if (bal > 0) {
          tokenUsers[l2T].push(user);
        }
      }
    }
  }

  function run() public {
    console.log('--- Comprehensive Asset Snapshot Generation ---');

    discoverAllUsersAndTokens();
    collectBalances();
    generateFiles();
  }

  function generateFiles() internal {
    string memory json3 = '[\n';
    uint256 tokenCount = 0;
    for (uint i = 0; i < uniqueL2Tokens.length; i++) {
      address l2T = uniqueL2Tokens[i];
      address[] memory users = tokenUsers[l2T];
      if (users.length == 0) continue;

      if (tokenCount > 0) json3 = string.concat(json3, ',\n');
      json3 = string.concat(
        json3,
        '  {\n',
        '    "l1Token": "',
        vm.toString(l2ToL1[l2T]),
        '",\n',
        '    "l2Token": "',
        vm.toString(l2T),
        '",\n',
        '    "tokenName": "',
        l2TokenNames[l2T],
        '",\n',
        '    "data": [\n'
      );

      for (uint j = 0; j < users.length; j++) {
        address u = users[j];
        uint256 bal;
        if (l2T == nativeToken) bal = u.balance;
        else {
          if (keccak256(bytes(tokenType[l2T])) == keccak256('ERC20'))
            bal = IERC20(l2T).balanceOf(u);
          else bal = IERC721(l2T).balanceOf(u);
        }

        if (j > 0) json3 = string.concat(json3, ',\n');
        bytes32 h = keccak256(abi.encodePacked(l2ToL1[l2T], u, bal));
        json3 = string.concat(
          json3,
          '      {\n',
          '        "claimer": "',
          vm.toString(u),
          '",\n',
          '        "amount": "',
          vm.toString(bal),
          '",\n',
          '        "hash": "',
          vm.toString(h),
          '"\n',
          '      }'
        );
      }
      json3 = string.concat(json3, '\n    ]\n  }');
      tokenCount++;
    }
    json3 = string.concat(json3, '\n]');

    string memory fileName = string.concat(
      'data/generate-assets-',
      vm.toString(block.chainid),
      '.json'
    );
    vm.writeFile(fileName, json3);
    console.log('Success: Written', fileName);
    console.log('Total users with assets:', discoveredUsers.length);
  }
}
