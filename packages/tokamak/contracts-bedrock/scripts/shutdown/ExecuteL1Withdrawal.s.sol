// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
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
 * @notice Phase 2 of L2 shutdown: liquidity sweep and claims execution (step 9).
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

    bool isDryRun = vm.envOr('DRY_RUN', false);
    if (!isDryRun) {
      vm.startBroadcast(closerPrivateKey);
    }

    // Step 9: Sweep liquidity to bridge
    _step9_sweepLiquidity(bridgeProxy, systemOwnerSafe, closerAddress);

    // Step 10: Execute claims (optional)
    bool executeClaims = vm.envOr('EXECUTE_CLAIMS', true);
    if (executeClaims) {
      _step10_executeClaims(bridgeProxy);
    } else {
      console.log('[INFO] EXECUTE_CLAIMS=true, skipping claims');
    }

    if (!isDryRun) {
      vm.stopBroadcast();
    }
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

    console.log('[INFO] Required native token: ', requiredNative);

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

    console.log('[INFO] Required USDC: ', requiredUsdc);

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

  /// @notice Claim result status
  uint8 internal constant CLAIM_SUCCESS = 0;
  uint8 internal constant CLAIM_ALREADY = 1;
  uint8 internal constant CLAIM_SKIPPED = 2;
  uint8 internal constant CLAIM_FAILED = 3;

  function _step10_executeClaims(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 10] Executing Force Withdrawal Claims');
    console.log('------------------------------------------');

    uint256 tokenCount = JsonUtils.countTokens(assetsJsonContent);
    bool isDryRun = vm.envOr('DRY_RUN', false);

    console.log('[INFO] Bridge proxy:', _bridgeProxy);
    console.log('[INFO] Position contract:', deployedStorage);
    console.log('[INFO] Token sets to process: ', tokenCount);

    uint256 totalClaims = 0;
    uint256 skippedClaims = 0;
    uint256 alreadyClaimed = 0;

    for (uint256 i = 0; i < tokenCount; i++) {
      (
        uint256 executed,
        uint256 already,
        uint256 skipped
      ) = _processTokenClaims(_bridgeProxy, i, isDryRun);
      totalClaims += executed;
      alreadyClaimed += already;
      skippedClaims += skipped;
    }

    console.log('[SUCCESS] Total claims executed: ', totalClaims);
    if (alreadyClaimed > 0) {
      console.log('[INFO] Already claimed (skipped): ', alreadyClaimed);
    }
    if (isDryRun && skippedClaims > 0) {
      console.log(
        '[WARN] DRY_RUN: claims skipped (bridge may not be upgraded or funded): ',
        skippedClaims
      );
    }
    console.log('------------------------------------------\n');
  }

  function _processTokenClaims(
    address _bridgeProxy,
    uint256 _tokenIndex,
    bool _isDryRun
  ) internal returns (uint256 executed_, uint256 already_, uint256 skipped_) {
    address l1Token = vm.parseJsonAddress(
      assetsJsonContent,
      string.concat('[', vm.toString(_tokenIndex), '].l1Token')
    );
    uint256 claimCount = JsonUtils.countClaimsInToken(
      assetsJsonContent,
      _tokenIndex
    );

    console.log('[INFO] Processing token:', l1Token);
    console.log('[INFO] Claim count: ', claimCount);

    for (uint256 j = 0; j < claimCount; j++) {
      uint8 result = _processSingleClaim(
        _bridgeProxy,
        l1Token,
        _tokenIndex,
        j,
        _isDryRun
      );
      if (result == CLAIM_SUCCESS) {
        executed_++;
      } else if (result == CLAIM_ALREADY) {
        already_++;
      } else if (result == CLAIM_SKIPPED) {
        skipped_++;
      }
      // CLAIM_FAILED reverts inside _processSingleClaim
    }
  }

  function _processSingleClaim(
    address _bridgeProxy,
    address _l1Token,
    uint256 _tokenIndex,
    uint256 _claimIndex,
    bool _isDryRun
  ) internal returns (uint8) {
    string memory base = string.concat(
      '[',
      vm.toString(_tokenIndex),
      '].data[',
      vm.toString(_claimIndex),
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

    // Check if already claimed
    if (_isAlreadyClaimed(_bridgeProxy, _l1Token, claimer, amount)) {
      console.log('[INFO] Already claimed, skipping claimer:', claimer);
      return CLAIM_ALREADY;
    }

    // Execute claim
    bytes memory claimData = abi.encodeWithSignature(
      'forceWithdrawClaim(address,string,address,uint256,address)',
      deployedStorage,
      hashString,
      _l1Token,
      amount,
      claimer
    );
    (bool success, ) = _bridgeProxy.call(claimData);

    if (success) {
      return CLAIM_SUCCESS;
    } else if (_isDryRun) {
      return CLAIM_SKIPPED;
    } else {
      revert('D1: claim execution failed');
    }
  }

  function _isAlreadyClaimed(
    address _bridgeProxy,
    address _token,
    address _claimer,
    uint256 _amount
  ) internal view returns (bool) {
    bytes32 claimHash = keccak256(abi.encodePacked(_token, _claimer, _amount));
    (bool success, bytes memory data) = _bridgeProxy.staticcall(
      abi.encodeWithSignature('claimState(bytes32)', claimHash)
    );
    if (success && data.length >= 32) {
      return abi.decode(data, (bool));
    }
    return false;
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
    console.log('[INFO] Sweep amount: ', amount);

    bytes memory sweepData = abi.encodeWithSignature(
      'sweepNativeToken(address,uint256)',
      bridgeProxy,
      amount
    );
    bool success = _execWithOwner(
      ownerToUse,
      deployerAddress,
      portal,
      sweepData,
      label
    );
    remaining_ = success ? remaining - amount : remaining;
  }

  function _execSweepUsdc(
    address ownerToUse,
    address deployerAddress,
    address l1UsdcBridge,
    address bridgeProxy,
    uint256 amount
  ) internal returns (bool) {
    bytes memory sweepData = abi.encodeWithSignature(
      'sweepUSDC(address,uint256)',
      bridgeProxy,
      amount
    );
    return
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
  ) internal returns (bool) {
    bool isDryRun = vm.envOr('DRY_RUN', false);
    if (SafeUtils.isContract(ownerToUse)) {
      if (isDryRun) {
        console.log('[INFO] DRY_RUN=true, simulating Safe call for:', label);
        vm.prank(ownerToUse);
        (bool ok, ) = target.call(data);
        if (!ok) {
          console.log(
            '[WARN] DRY_RUN: call failed (contract may not be upgraded yet), skipping:',
            label
          );
          return false;
        }
        return true;
      }
      bool executed = SafeUtils.execViaSafeFromEnv(
        IGnosisSafe(ownerToUse),
        target,
        data,
        label
      );
      require(executed, 'D1: safe execution failed');
      return true;
    }

    require(
      ownerToUse == deployerAddress || isDryRun,
      'D1: caller is not owner'
    );
    (bool success, ) = target.call(data);
    if (isDryRun && !success) {
      console.log(
        '[WARN] DRY_RUN: call failed (contract may not be upgraded yet), skipping:',
        label
      );
      return false;
    }
    require(success, 'D1: sweep call failed');
    return true;
  }

  function _getEip1967Admin(
    address _proxy
  ) internal view returns (address admin_) {
    bytes32 data = vm.load(_proxy, ShutdownConfig.PROXY_ADMIN_SLOT);
    admin_ = address(uint160(uint256(data)));
  }
}
