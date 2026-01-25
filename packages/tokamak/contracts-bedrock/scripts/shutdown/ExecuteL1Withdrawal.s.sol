// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {ForceWithdrawBridge} from '../../src/shutdown/ForceWithdrawBridge.sol';
import {IERC20} from '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import {ProxyAdmin} from '../../src/universal/ProxyAdmin.sol';
import {IGnosisSafe, Enum} from '../interfaces/IGnosisSafe.sol';
import {IL1UsdcBridge} from './interfaces/IL1UsdcBridge.sol';
import {IOptimismPortal} from './interfaces/IOptimismPortal.sol';

import {ShutdownConfig} from './lib/ShutdownConfig.sol';
import {JsonUtils} from './lib/JsonUtils.sol';
import {SafeUtils} from './lib/SafeUtils.sol';

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

    ForceWithdrawBridge bridge = ForceWithdrawBridge(payable(_bridge));
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
    console.log('[INFO] Bridge proxy:', _bridgeProxy);
    console.log('[INFO] Owner to Use:', _resolveOwner(_systemOwnerSafe));

    _sweepNativeLiquidity(_bridgeProxy, _systemOwnerSafe, _deployerAddress);
    _sweepUsdcLiquidity(_bridgeProxy, _systemOwnerSafe, _deployerAddress);

    console.log('------------------------------------------\n');
  }

  function _resolveOwner(
    address _systemOwnerSafe
  ) internal view returns (address) {
    return
      _systemOwnerSafe != address(0)
        ? _systemOwnerSafe
        : ProxyAdmin(vm.envAddress('PROXY_ADMIN')).owner();
  }

  function _sweepNativeLiquidity(
    address _bridgeProxy,
    address _systemOwnerSafe,
    address _deployerAddress
  ) internal {
    address l1NativeToken = vm.envOr('L1_NATIVE_TOKEN', address(0));
    if (l1NativeToken == address(0)) {
      console.log('[WARN] L1_NATIVE_TOKEN not set, skipping sweep');
      return;
    }

    address optimismPortal = vm.envOr('OPTIMISM_PORTAL_PROXY', address(0));
    address ownerToUse = _resolveOwner(_systemOwnerSafe);
    uint256 requiredNative = _sumClaimsForToken(l1NativeToken);

    console.log('[INFO] Required native token:', requiredNative);

    uint256 remaining = _remainingAmount(
      _bridgeProxy,
      l1NativeToken,
      requiredNative
    );
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
  }

  function _sweepUsdcLiquidity(
    address _bridgeProxy,
    address _systemOwnerSafe,
    address _deployerAddress
  ) internal {
    address l1UsdcBridge = vm.envOr('L1_USDC_BRIDGE_PROXY', address(0));
    if (l1UsdcBridge == address(0)) {
      console.log('[WARN] L1_USDC_BRIDGE_PROXY not set, skipping sweep');
      return;
    }

    address l1UsdcToken = IL1UsdcBridge(l1UsdcBridge).l1Usdc();

    // Check if L1 USDC token contract exists
    if (l1UsdcToken.code.length == 0) {
      console.log('[WARN] L1 USDC token contract not found at:', l1UsdcToken);
      console.log(
        '[WARN] Skipping USDC sweep (for devnet. USDC token not deployed to L1).'
      );
      return;
    }

    uint256 requiredUsdc = _sumClaimsForToken(l1UsdcToken);

    console.log('[INFO] Required USDC:', requiredUsdc);

    uint256 remaining = _remainingAmount(
      _bridgeProxy,
      l1UsdcToken,
      requiredUsdc
    );
    if (remaining == 0) {
      console.log('[INFO] USDC already sufficient on bridge');
      return;
    }

    uint256 amount = _min(remaining, _getBalance(l1UsdcBridge, l1UsdcToken));
    if (amount == 0) {
      console.log('[WARN] L1 USDC bridge balance insufficient');
      return;
    }

    console.log('[INFO] Sweeping USDC amount:', amount);
    _execSweepUsdc(
      _resolveUsdcSweepOwner(l1UsdcBridge, _systemOwnerSafe),
      _deployerAddress,
      l1UsdcBridge,
      _bridgeProxy,
      amount
    );
  }

  function _resolveUsdcSweepOwner(
    address _l1UsdcBridge,
    address _systemOwnerSafe
  ) internal view returns (address) {
    address l1UsdcBridgeAdmin = vm.envOr('L1_USDC_BRIDGE_ADMIN', address(0));
    if (l1UsdcBridgeAdmin != address(0)) {
      return l1UsdcBridgeAdmin;
    }

    address adminFromSlot = _getEip1967Admin(_l1UsdcBridge);
    if (adminFromSlot != address(0)) {
      return adminFromSlot;
    }

    return _resolveOwner(_systemOwnerSafe);
  }

  function _remainingAmount(
    address holder,
    address token,
    uint256 required
  ) internal view returns (uint256) {
    uint256 balance = _getBalance(holder, token);
    return required > balance ? required - balance : 0;
  }

  function _min(uint256 a, uint256 b) internal pure returns (uint256) {
    return a < b ? a : b;
  }

  function _step9_executeClaims(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 9] Executing Force Withdrawal Claims');
    console.log('------------------------------------------');

    ForceWithdrawBridge bridge = ForceWithdrawBridge(payable(_bridgeProxy));
    uint256 tokenCount = JsonUtils.countTokens(assetsJsonContent);

    console.log('[INFO] Bridge proxy:', _bridgeProxy);
    console.log('[INFO] Position contract:', deployedStorage);
    console.log('[INFO] Token sets to process:', tokenCount);

    uint256 totalClaims = 0;
    for (uint256 i = 0; i < tokenCount; i++) {
      address l1Token = vm.parseJsonAddress(
        assetsJsonContent,
        string.concat('[', vm.toString(i), '].l1Token')
      );
      uint256 claimCount = JsonUtils.countClaimsInToken(assetsJsonContent, i);

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
        string memory hashString = JsonUtils.strip0x(
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
    uint256 tokenCount = JsonUtils.countTokens(assetsJsonContent);
    uint256 total = 0;
    for (uint256 i = 0; i < tokenCount; i++) {
      address l1Token = vm.parseJsonAddress(
        assetsJsonContent,
        string.concat('[', vm.toString(i), '].l1Token')
      );
      if (l1Token != token) continue;
      uint256 claimCount = JsonUtils.countClaimsInToken(assetsJsonContent, i);
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
    if (SafeUtils.isContract(ownerToUse)) {
      bool executed = SafeUtils.execViaSafeFromEnv(
        IGnosisSafe(ownerToUse),
        target,
        data,
        label
      );
      require(executed, 'D1: safe execution failed');
      return;
    }

    require(ownerToUse == deployerAddress, 'D1: caller is not owner');
    (bool success, ) = target.call(data);
    require(success, 'D1: sweep call failed');
  }

  function _getEip1967Admin(
    address _proxy
  ) internal view returns (address admin_) {
    bytes32 data = vm.load(_proxy, ShutdownConfig.PROXY_ADMIN_SLOT);
    admin_ = address(uint160(uint256(data)));
  }
}
