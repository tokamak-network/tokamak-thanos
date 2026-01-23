// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {UpgradeL1BridgeV1} from '../../src/shutdown/ForceWithdrawBridge.sol';
import {IERC20} from '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import {ProxyAdmin} from '../../src/universal/ProxyAdmin.sol';
import {IGnosisSafe, Enum} from '../interfaces/IGnosisSafe.sol';
import {IL1UsdcBridge} from './interfaces/IL1UsdcBridge.sol';
import {IOptimismPortal} from './interfaces/IOptimismPortal.sol';

/**
 * @title ExecuteL1Withdrawal
 * @notice Phase 2 of L2 shutdown: activation and execution (steps 8-9).
 */
contract ExecuteL1Withdrawal is Script {
  using stdJson for string;

  string public assetsJsonContent;
  address public deployedStorage;

  function run() public {
    address bridgeProxy = vm.envAddress('BRIDGE_PROXY');
    address storageAddress = vm.envAddress('STORAGE_ADDRESS');
    string memory dataPath = vm.envString('DATA_PATH');
    address systemOwnerSafe = vm.envOr('SYSTEM_OWNER_SAFE', address(0));

    deployedStorage = storageAddress;
    assetsJsonContent = vm.readFile(dataPath);

    uint256 closerPrivateKey = vm.envUint('PRIVATE_KEY');
    address closerAddress = vm.addr(closerPrivateKey);

    vm.startBroadcast(closerPrivateKey);

    // Step 8: Activate force withdrawal mode
    _step8_activate(bridgeProxy, closerAddress);

    // Step 9: Sweep liquidity to bridge
    _step9_sweepLiquidity(bridgeProxy, systemOwnerSafe, closerAddress);

    // Step 9: Execute claims (optional)
    bool executeClaims = vm.envOr('EXECUTE_CLAIMS', false);
    if (executeClaims) {
      _step9_executeClaims(bridgeProxy);
    } else {
      console.log('[INFO] EXECUTE_CLAIMS=false, skipping claims');
    }

    vm.stopBroadcast();
  }

  function _step8_activate(address _bridge, address _closer) internal {
    console.log('------------------------------------------');
    console.log('[STEP 8] Activating Force Withdrawal Mode');
    console.log('------------------------------------------');

    UpgradeL1BridgeV1 bridge = UpgradeL1BridgeV1(payable(_bridge));
    bool currentState = bridge.active();
    console.log(
      '[INFO] Current active state:',
      currentState ? 'ACTIVE' : 'INACTIVE'
    );

    if (currentState) {
      console.log('[INFO] Force withdrawal already active, skipping');
    } else {
      console.log('[ACTION] Activating force withdrawal...');
      bridge.forceActive(true);
      console.log('[SUCCESS] Force withdrawal mode activated');
    }

    console.log('[INFO] Closer:', _closer);
    console.log('------------------------------------------\n');
  }

  function _step9_sweepLiquidity(
    address _bridgeProxy,
    address _systemOwnerSafe,
    address _deployerAddress
  ) internal {
    console.log('------------------------------------------');
    console.log('[STEP 9] Sweeping Liquidity to Bridge');
    console.log('------------------------------------------');

    address optimismPortal = vm.envOr('OPTIMISM_PORTAL_PROXY', address(0));
    address l1NativeToken = vm.envOr('L1_NATIVE_TOKEN', address(0));
    address l1UsdcBridge = vm.envOr('L1_USDC_BRIDGE_PROXY', address(0));
    address l1UsdcBridgeAdmin = vm.envOr('L1_USDC_BRIDGE_ADMIN', address(0));

    address l1UsdcToken = address(0);
    if (l1UsdcBridge != address(0)) {
      l1UsdcToken = IL1UsdcBridge(l1UsdcBridge).l1Usdc();
    }

    address ownerToUse = _systemOwnerSafe != address(0)
      ? _systemOwnerSafe
      : ProxyAdmin(vm.envAddress('PROXY_ADMIN')).owner();

    uint256 requiredNative = _sumClaimsForToken(l1NativeToken);
    uint256 requiredUsdc = _sumClaimsForToken(l1UsdcToken);

    console.log('[INFO] Bridge proxy:', _bridgeProxy);
    console.log('[INFO] Owner to Use:', ownerToUse);
    console.log('[INFO] Required native token:', requiredNative);
    console.log('[INFO] Required USDC:', requiredUsdc);

    if (l1NativeToken != address(0)) {
      uint256 bridgeBalance = _getBalance(_bridgeProxy, l1NativeToken);
      uint256 remaining = requiredNative > bridgeBalance
        ? requiredNative - bridgeBalance
        : 0;
      if (remaining > 0) {
        _sweepFromPortal(
          optimismPortal,
          ownerToUse,
          _deployerAddress,
          _bridgeProxy,
          remaining,
          'OptimismPortal'
        );
      } else {
        console.log('[INFO] Native token already sufficient on bridge');
      }
    } else {
      console.log('[WARN] L1_NATIVE_TOKEN not set, skipping sweep');
    }

    if (l1UsdcToken != address(0) && l1UsdcBridge != address(0)) {
      uint256 bridgeBalance = _getBalance(_bridgeProxy, l1UsdcToken);
      uint256 remaining = requiredUsdc > bridgeBalance
        ? requiredUsdc - bridgeBalance
        : 0;
      if (remaining > 0) {
        uint256 escrowBalance = _getBalance(l1UsdcBridge, l1UsdcToken);
        uint256 amount = remaining < escrowBalance ? remaining : escrowBalance;
        if (amount > 0) {
          address adminFromSlot = _getEip1967Admin(l1UsdcBridge);
          address sweepOwner = ownerToUse;
          if (l1UsdcBridgeAdmin != address(0)) {
            sweepOwner = l1UsdcBridgeAdmin;
          } else if (adminFromSlot != address(0)) {
            sweepOwner = adminFromSlot;
          }
          console.log('[INFO] Sweeping USDC amount:', amount);
          _execSweepUsdc(
            sweepOwner,
            _deployerAddress,
            l1UsdcBridge,
            _bridgeProxy,
            amount
          );
        } else {
          console.log('[WARN] L1 USDC bridge balance insufficient');
        }
      } else {
        console.log('[INFO] USDC already sufficient on bridge');
      }
    } else {
      console.log('[WARN] L1_USDC_BRIDGE_PROXY not set, skipping sweep');
    }

    console.log('------------------------------------------\n');
  }

  function _step9_executeClaims(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 9] Executing Force Withdrawal Claims');
    console.log('------------------------------------------');

    UpgradeL1BridgeV1 bridge = UpgradeL1BridgeV1(payable(_bridgeProxy));
    uint256 tokenCount = _countTokens(assetsJsonContent);

    console.log('[INFO] Bridge proxy:', _bridgeProxy);
    console.log('[INFO] Position contract:', deployedStorage);
    console.log('[INFO] Token sets to process:', tokenCount);

    uint256 totalClaims = 0;
    for (uint256 i = 0; i < tokenCount; i++) {
      address l1Token = vm.parseJsonAddress(
        assetsJsonContent,
        string.concat('[', vm.toString(i), '].l1Token')
      );
      uint256 claimCount = _countClaimsInToken(assetsJsonContent, i);

      console.log('[INFO] Processing token:', l1Token);
      console.log('[INFO] Claim count:', claimCount);

      for (uint256 j = 0; j < claimCount; j++) {
        string memory base = string.concat(
          '[',
          vm.toString(i),
          '].data[',
          vm.toString(j),
          ']'
        );
        string memory hashString = _strip0x(
          vm.parseJsonString(assetsJsonContent, string.concat(base, '.hash'))
        );
        uint256 amount = vm.parseUint(
          vm.parseJsonString(assetsJsonContent, string.concat(base, '.amount'))
        );
        address claimer = vm.parseJsonAddress(
          assetsJsonContent,
          string.concat(base, '.claimer')
        );

        bridge.forceWithdrawClaim(
          deployedStorage,
          hashString,
          l1Token,
          amount,
          claimer
        );
        totalClaims++;
      }
    }

    console.log('[SUCCESS] Total claims executed:', totalClaims);
    console.log('------------------------------------------\n');
  }

  function _getBalance(
    address holder,
    address token
  ) internal view returns (uint256) {
    if (token == address(0)) {
      return holder.balance;
    }
    return IERC20(token).balanceOf(holder);
  }

  function _sumClaimsForToken(address token) internal view returns (uint256) {
    if (token == address(0)) {
      return 0;
    }
    uint256 tokenCount = _countTokens(assetsJsonContent);
    uint256 total = 0;
    for (uint256 i = 0; i < tokenCount; i++) {
      address l1Token = vm.parseJsonAddress(
        assetsJsonContent,
        string.concat('[', vm.toString(i), '].l1Token')
      );
      if (l1Token != token) continue;
      uint256 claimCount = _countClaimsInToken(assetsJsonContent, i);
      for (uint256 j = 0; j < claimCount; j++) {
        string memory base = string.concat(
          '[',
          vm.toString(i),
          '].data[',
          vm.toString(j),
          ']'
        );
        uint256 amount = vm.parseUint(
          vm.parseJsonString(assetsJsonContent, string.concat(base, '.amount'))
        );
        total += amount;
      }
    }
    return total;
  }

  function _sweepFromPortal(
    address portal,
    address ownerToUse,
    address deployerAddress,
    address bridgeProxy,
    uint256 remaining,
    string memory label
  ) internal returns (uint256 remaining_) {
    if (portal == address(0)) {
      console.log('[WARN] Portal not set for:', label);
      return remaining;
    }

    address l1NativeToken = vm.envOr('L1_NATIVE_TOKEN', address(0));
    uint256 escrowBalance = _getBalance(portal, l1NativeToken);
    if (escrowBalance == 0) {
      console.log('[WARN] Portal balance is zero for:', label);
      return remaining;
    }

    uint256 amount = remaining < escrowBalance ? remaining : escrowBalance;
    console.log('[INFO] Sweeping native token from:', label);
    console.log('[INFO] Sweep amount:', amount);

    bytes memory sweepData = abi.encodeWithSignature(
      'sweepNativeToken(address,uint256)',
      bridgeProxy,
      amount
    );
    _execWithOwner(ownerToUse, deployerAddress, portal, sweepData, label);
    remaining_ = remaining - amount;
  }

  function _execSweepUsdc(
    address ownerToUse,
    address deployerAddress,
    address l1UsdcBridge,
    address bridgeProxy,
    uint256 amount
  ) internal {
    bytes memory sweepData = abi.encodeWithSignature(
      'sweepUSDC(address,uint256)',
      bridgeProxy,
      amount
    );
    _execWithOwner(
      ownerToUse,
      deployerAddress,
      l1UsdcBridge,
      sweepData,
      'L1 USDC Sweep'
    );
  }

  function _execWithOwner(
    address ownerToUse,
    address deployerAddress,
    address target,
    bytes memory data,
    string memory label
  ) internal {
    if (_isContract(ownerToUse)) {
      bool executed = _execViaSafe(IGnosisSafe(ownerToUse), target, data, label);
      require(executed, 'D1: safe execution failed');
      return;
    }

    require(ownerToUse == deployerAddress, 'D1: caller is not owner');
    (bool success, ) = target.call(data);
    require(success, 'D1: sweep call failed');
  }

  function _execViaSafe(
    IGnosisSafe _safe,
    address _target,
    bytes memory _data,
    string memory _label
  ) internal returns (bool executed_) {
    uint256 threshold = _safe.getThreshold();
    address caller = vm.addr(vm.envUint('PRIVATE_KEY'));

    require(threshold == 1, 'D1: safe threshold > 1');
    require(_safe.isOwner(caller), 'D1: caller not safe owner');

    bytes memory signature = abi.encodePacked(
      uint256(uint160(caller)),
      bytes32(0),
      uint8(1)
    );
    executed_ = _safe.execTransaction({
      to: _target,
      value: 0,
      data: _data,
      operation: Enum.Operation.Call,
      safeTxGas: 0,
      baseGas: 0,
      gasPrice: 0,
      gasToken: address(0),
      refundReceiver: payable(address(0)),
      signatures: signature
    });
    require(executed_, 'D1: safe execTransaction failed');
  }

  function _isContract(address _addr) internal view returns (bool) {
    return _addr.code.length > 0;
  }

  function _getEip1967Admin(
    address _proxy
  ) internal view returns (address admin_) {
    bytes32 slot = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;
    bytes32 data = vm.load(_proxy, slot);
    admin_ = address(uint160(uint256(data)));
  }

  function _strip0x(string memory value) internal pure returns (string memory) {
    bytes memory data = bytes(value);
    if (
      data.length >= 2 &&
      data[0] == bytes1('0') &&
      (data[1] == bytes1('x') || data[1] == bytes1('X'))
    ) {
      bytes memory trimmed = new bytes(data.length - 2);
      for (uint256 i = 2; i < data.length; i++) {
        trimmed[i - 2] = data[i];
      }
      return string(trimmed);
    }
    return value;
  }

  function _countTokens(string memory json) internal pure returns (uint256) {
    return _countOccurrences(json, '"l1Token"');
  }

  function _countClaimsInToken(
    string memory json,
    uint256 tokenIdx
  ) internal pure returns (uint256) {
    return _countOccurrencesInTokenData(json, tokenIdx, '"hash"');
  }

  function _countOccurrencesInTokenData(
    string memory json,
    uint256 tokenIdx,
    string memory pattern
  ) internal pure returns (uint256) {
    bytes memory bJson = bytes(json);
    bytes memory bPattern = bytes(pattern);
    bytes memory bL1Token = bytes('"l1Token"');
    uint256 count = 0;
    uint256 tokenFoundCount = 0;

    for (uint256 i = 0; i <= bJson.length - bL1Token.length; i++) {
      if (_isMatch(bJson, i, bL1Token)) {
        if (tokenFoundCount == tokenIdx) {
          for (
            uint256 j = i + bL1Token.length;
            j <= bJson.length - bPattern.length;
            j++
          ) {
            if (_isMatch(bJson, j, bL1Token)) break;

            if (_isMatch(bJson, j, bPattern)) {
              count++;
            }
          }
          return count;
        }
        tokenFoundCount++;
      }
    }
    return 0;
  }

  function _countOccurrences(
    string memory json,
    string memory pattern
  ) internal pure returns (uint256) {
    bytes memory bJson = bytes(json);
    bytes memory bPattern = bytes(pattern);
    uint256 count = 0;
    if (bJson.length < bPattern.length) return 0;
    for (uint256 i = 0; i <= bJson.length - bPattern.length; i++) {
      if (_isMatch(bJson, i, bPattern)) {
        count++;
      }
    }
    return count;
  }

  function _isMatch(
    bytes memory data,
    uint256 index,
    bytes memory pattern
  ) internal pure returns (bool) {
    if (index + pattern.length > data.length) return false;
    for (uint256 i = 0; i < pattern.length; i++) {
      if (data[index + i] != pattern[i]) return false;
    }
    return true;
  }
}
