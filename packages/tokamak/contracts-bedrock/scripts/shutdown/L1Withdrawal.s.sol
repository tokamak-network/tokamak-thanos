// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {ForceWithdrawBridge} from '../../src/shutdown/ForceWithdrawBridge.sol';
import {
  ShutdownOptimismPortal
} from '../../src/shutdown/ShutdownOptimismPortal.sol';
import {
  ShutdownOptimismPortal2
} from '../../src/shutdown/ShutdownOptimismPortal2.sol';
import {GenFWStorage} from '../../src/shutdown/GenFWStorage.sol';
import {ProxyAdmin} from '../../src/universal/ProxyAdmin.sol';
import {IERC20} from '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import {IGnosisSafe, Enum} from '../interfaces/IGnosisSafe.sol';
import {
  ShutdownL1UsdcBridge
} from '../../src/shutdown/ShutdownL1UsdcBridge.sol';

interface IL1UsdcBridge {
  function l1Usdc() external view returns (address);
  function sweepUSDC(address _to, uint256 _amount) external;
}

interface IOptimismPortal {
  function sweepNativeToken(address _to, uint256 _amount) external;
}

interface IOptimismPortal2 {
  function proofMaturityDelaySeconds() external view returns (uint256);
  function disputeGameFinalityDelaySeconds() external view returns (uint256);
}

/**
 * @title L1Withdrawal
 * @notice Forge script to L1 force withdrawals during L2 shutdown.
 *         Executes all steps required to prepare the L1 bridge for force withdrawals:
 *         1. Deploy new ForceWithdrawBridge implementation
 *         2. Upgrade bridge proxy via ProxyAdmin
 *         3. Deploy GenFWStorage contract with snapshot data
 *         4. Register storage contract with bridge
 *         5. Activate force withdrawal mode
 *         6. Execute force withdrawal claims (optional)
 *
 * @dev Requirements:
 *      - EOA (Externally Owned Account) must own the ProxyAdmin
 *      - Private key of the EOA must be available in PRIVATE_KEY env variable
 *      - Assets JSON file must contain properly formatted snapshot data
 *
 */
contract L1Withdrawal is Script {
  using stdJson for string;

  struct ClaimEntry {
    string amount;
    address claimer;
    string hash;
  }

  struct TokenClaims {
    ClaimEntry[] data;
    address l1Token;
    address l2Token;
    string tokenName;
  }

  // Deployment state variables
  address public deployedImpl;
  address public deployedPortalImpl;
  address public deployedUsdcBridgeImpl;
  address public deployedStorage;
  string public assetsJsonContent;

  /**
   * @notice Default entry point - reads configuration from .env
   * @dev This is called when no function signature is specified.
   *      Designed for L2 Operators to prepare the L1 bridge for shutdown.
   */
  function run() public {
    // Load configuration from .env
    address bridgeProxy = vm.envAddress('BRIDGE_PROXY');
    address proxyAdmin = vm.envAddress('PROXY_ADMIN');
    string memory dataPath = vm.envString('DATA_PATH');
    bool executeClaims = vm.envOr('EXECUTE_CLAIMS', false);
    address systemOwnerSafe = vm.envOr('SYSTEM_OWNER_SAFE', address(0));
    address closerPrivateKeyAddress = vm.envOr(
      'CLOSER_ADDRESS_PRIVATE_KEY',
      uint256(0)
    ) == 0
      ? address(0)
      : vm.addr(vm.envUint('CLOSER_ADDRESS_PRIVATE_KEY'));
    address optimismPortal = vm.envOr('OPTIMISM_PORTAL_PROXY', address(0));
    bool optimismPortalUseV2 = vm.envOr('OPTIMISM_PORTAL_USE_V2', false);
    address l1NativeToken = vm.envOr('L1_NATIVE_TOKEN', address(0));
    address l1UsdcBridge = vm.envOr('L1_USDC_BRIDGE_PROXY', address(0));
    address l1UsdcBridgeAdmin = vm.envOr('L1_USDC_BRIDGE_ADMIN', address(0));
    address l1UsdcBridgeProxyAdmin = vm.envOr(
      'L1_USDC_BRIDGE_PROXY_ADMIN',
      address(0)
    );

    console.log('==============================================');
    console.log('L2 OPERATOR: ENABLE L1 WITHDRAWAL');
    console.log('==============================================');
    console.log('[ENV] Configuration loaded from .env:');
    console.log('  BRIDGE_PROXY:', bridgeProxy);
    console.log('  PROXY_ADMIN:', proxyAdmin);
    console.log('  SYSTEM_OWNER_SAFE:', systemOwnerSafe);
    console.log('  CLOSER_ADDRESS_PRIVATE_KEY:', closerPrivateKeyAddress);
    console.log('  OPTIMISM_PORTAL_PROXY:', optimismPortal);
    console.log('  OPTIMISM_PORTAL_USE_V2:', optimismPortalUseV2);
    console.log('  L1_NATIVE_TOKEN:', l1NativeToken);
    console.log('  L1_USDC_BRIDGE_PROXY:', l1UsdcBridge);
    console.log('  L1_USDC_BRIDGE_ADMIN:', l1UsdcBridgeAdmin);
    console.log('  L1_USDC_BRIDGE_PROXY_ADMIN:', l1UsdcBridgeProxyAdmin);
    console.log('  DATA_PATH:', dataPath);
    console.log('  EXECUTE_CLAIMS (Verification):', executeClaims);
    console.log('==============================================\n');

    // Execute shutdown flow
    _executeShutdown(
      bridgeProxy,
      proxyAdmin,
      dataPath,
      executeClaims,
      systemOwnerSafe,
      closerPrivateKeyAddress
    );
  }

  /**
   * @notice Run shutdown flow with explicit parameters (legacy)
   */
  function runWithParams(
    address _bridgeProxy,
    address _proxyAdmin,
    string memory _assetsPath,
    bool _executeClaims,
    address _systemOwnerSafe
  ) public {
    console.log('==============================================');
    console.log('L2 OPERATOR: ENABLE L1 WITHDRAWAL');
    console.log('==============================================\n');

    _executeShutdown(
      _bridgeProxy,
      _proxyAdmin,
      _assetsPath,
      _executeClaims,
      _systemOwnerSafe,
      address(0)
    );
  }

  function runWithParamsAndCloser(
    address _bridgeProxy,
    address _proxyAdmin,
    string memory _assetsPath,
    bool _executeClaims,
    address _systemOwnerSafe,
    address _closerPrivateKeyAddress
  ) public {
    console.log('==============================================');
    console.log('L2 OPERATOR: ENABLE L1 WITHDRAWAL');
    console.log('==============================================\n');

    _executeShutdown(
      _bridgeProxy,
      _proxyAdmin,
      _assetsPath,
      _executeClaims,
      _systemOwnerSafe,
      _closerPrivateKeyAddress
    );
  }

  // ============================================
  // Main Execution Flow
  // ============================================

  function _executeShutdown(
    address _bridgeProxy,
    address _proxyAdmin,
    string memory _assetsPath,
    bool _executeClaims,
    address _systemOwnerSafe,
    address _closerPrivateKeyAddress
  ) internal {
    _validateInputs(_bridgeProxy, _proxyAdmin, _assetsPath);
    _loadAssetsJson(_assetsPath);

    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');
    address deployerAddress = vm.addr(deployerPrivateKey);
    address closerToSet = _closerPrivateKeyAddress != address(0)
      ? _closerPrivateKeyAddress
      : deployerAddress;
    require(closerToSet != address(0), 'D1: closer address required');

    vm.startBroadcast(deployerPrivateKey);

    // Step 1: Deploy new implementation
    _step1_deployImplementation();

    // Step 2: Upgrade bridge proxy
    _step2_upgradeProxy(_bridgeProxy, _proxyAdmin, _systemOwnerSafe);

    // Step 3: Upgrade OptimismPortal
    _step3_upgradeOptimismPortal(_proxyAdmin, _systemOwnerSafe);

    // Step 4: Upgrade L1 USDC bridge
    _step4_upgradeL1UsdcBridge(_proxyAdmin, _systemOwnerSafe, deployerAddress);

    // Step 5: Set closer and activate
    _step5_setCloserAndActive(
      _bridgeProxy,
      _proxyAdmin,
      _systemOwnerSafe,
      closerToSet,
      deployerAddress
    );

    // Step 6: Deploy GenFWStorage with snapshot data
    _step6_deployStorage();

    // Step 7: Register storage with bridge
    _step7_registerStorage(_bridgeProxy);

    // Step 8: Activate force withdrawal mode
    _step8_activateForceWithdraw(_bridgeProxy);

    // Step 9: Sweep liquidity to bridge
    _step9_sweepLiquidity(_bridgeProxy, _systemOwnerSafe, deployerAddress);

    // Step 10: Execute verification claims (optional)
    if (_executeClaims) {
      _step10_executeClaims(_bridgeProxy);
    }

    vm.stopBroadcast();

    _printSummary(_executeClaims, _assetsPath);
  }

  // ============================================
  // Internal Step Functions
  // ============================================

  function _step1_deployImplementation() internal {
    console.log('------------------------------------------');
    console.log('[STEP 1] Deploying ForceWithdrawBridge Implementation');
    console.log('------------------------------------------');

    ForceWithdrawBridge impl = new ForceWithdrawBridge();
    deployedImpl = address(impl);

    console.log('[SUCCESS] Implementation deployed at:', deployedImpl);
    console.log('');
  }

  function _step2_upgradeProxy(
    address _bridgeProxy,
    address _proxyAdmin,
    address _systemOwnerSafe
  ) internal {
    console.log('------------------------------------------');
    console.log('[STEP 2] Upgrading Bridge Proxy');
    console.log('------------------------------------------');
    _upgradeProxyWithSafe(
      _bridgeProxy,
      _proxyAdmin,
      _systemOwnerSafe,
      deployedImpl,
      'Bridge Proxy Upgrade'
    );
    console.log('------------------------------------------\n');
  }

  function _step3_upgradeOptimismPortal(
    address _proxyAdmin,
    address _systemOwnerSafe
  ) internal {
    console.log('------------------------------------------');
    console.log('[STEP 3] Upgrading OptimismPortal');
    console.log('------------------------------------------');

    address optimismPortal = vm.envOr('OPTIMISM_PORTAL_PROXY', address(0));
    bool optimismPortalUseV2 = vm.envOr('OPTIMISM_PORTAL_USE_V2', false);
    if (optimismPortal == address(0)) {
      console.log('[WARN] OPTIMISM_PORTAL_PROXY not set, skipping');
      console.log('------------------------------------------\n');
      return;
    }

    if (optimismPortalUseV2) {
      uint256 proofDelay = IOptimismPortal2(optimismPortal)
        .proofMaturityDelaySeconds();
      uint256 gameDelay = IOptimismPortal2(optimismPortal)
        .disputeGameFinalityDelaySeconds();
      ShutdownOptimismPortal2 impl = new ShutdownOptimismPortal2(
        proofDelay,
        gameDelay
      );
      deployedPortalImpl = address(impl);
    } else {
      ShutdownOptimismPortal impl = new ShutdownOptimismPortal();
      deployedPortalImpl = address(impl);
    }

    console.log('[INFO] OptimismPortal Proxy:', optimismPortal);
    console.log('[INFO] New Implementation:', deployedPortalImpl);

    _upgradeProxyWithSafe(
      optimismPortal,
      _proxyAdmin,
      _systemOwnerSafe,
      deployedPortalImpl,
      'OptimismPortal Upgrade'
    );

    console.log('------------------------------------------\n');
  }

  function _step4_upgradeL1UsdcBridge(
    address _proxyAdmin,
    address _systemOwnerSafe,
    address _deployerAddress
  ) internal {
    console.log('------------------------------------------');
    console.log('[STEP 4] Upgrading L1 USDC Bridge');
    console.log('------------------------------------------');

    address l1UsdcBridge = vm.envOr('L1_USDC_BRIDGE_PROXY', address(0));
    address l1UsdcBridgeAdmin = vm.envOr('L1_USDC_BRIDGE_ADMIN', address(0));
    address l1UsdcBridgeProxyAdmin = vm.envOr(
      'L1_USDC_BRIDGE_PROXY_ADMIN',
      address(0)
    );
    if (l1UsdcBridge == address(0)) {
      console.log('[WARN] L1_USDC_BRIDGE_PROXY not set, skipping');
      console.log('------------------------------------------\n');
      return;
    }

    ShutdownL1UsdcBridge impl = new ShutdownL1UsdcBridge();
    deployedUsdcBridgeImpl = address(impl);

    console.log('[INFO] L1 USDC Bridge Proxy:', l1UsdcBridge);
    console.log('[INFO] New Implementation:', deployedUsdcBridgeImpl);

    address adminFromSlot = _getEip1967Admin(l1UsdcBridge);
    address proxyAdminToUse = l1UsdcBridgeProxyAdmin != address(0)
      ? l1UsdcBridgeProxyAdmin
      : _proxyAdmin;
    if (adminFromSlot != address(0) && adminFromSlot != proxyAdminToUse) {
      if (l1UsdcBridgeAdmin == address(0)) {
        l1UsdcBridgeAdmin = adminFromSlot;
      }
    }

    if (
      l1UsdcBridgeAdmin != address(0) && l1UsdcBridgeAdmin == _deployerAddress
    ) {
      console.log('[INFO] Using EOA admin for L1 USDC bridge upgrade');
      _upgradeProxyByEoa(l1UsdcBridge, deployedUsdcBridgeImpl);
    } else {
      _upgradeProxyWithSafe(
        l1UsdcBridge,
        proxyAdminToUse,
        _systemOwnerSafe,
        deployedUsdcBridgeImpl,
        'L1 USDC Bridge Upgrade'
      );
    }

    console.log('------------------------------------------\n');
  }

  function _step5_setCloserAndActive(
    address _bridgeProxy,
    address _proxyAdmin,
    address _systemOwnerSafe,
    address _closerAddress,
    address _deployerAddress
  ) internal {
    console.log('------------------------------------------');
    console.log('[STEP 5] Setting closer and activating');
    console.log('------------------------------------------');

    bytes memory setCloserData = abi.encodeWithSignature(
      'setCloserAndActive(address,bool)',
      _closerAddress,
      true
    );

    address adminOwner = ProxyAdmin(_proxyAdmin).owner();
    address ownerToUse = _systemOwnerSafe != address(0)
      ? _systemOwnerSafe
      : adminOwner;

    console.log('[INFO] Bridge Proxy:', _bridgeProxy);
    console.log('[INFO] Closer Address:', _closerAddress);
    console.log('[INFO] Owner to Use:', ownerToUse);

    if (_isContract(ownerToUse)) {
      console.log('[INFO] Owner is a contract (Safe). Preparing Safe TX...');
      IGnosisSafe safe = IGnosisSafe(ownerToUse);

      _logSafeTx(safe, _bridgeProxy, setCloserData, 'Set Closer And Active');

      bool success = _execViaSafe(
        safe,
        _bridgeProxy,
        setCloserData,
        'Set Closer And Active'
      );
      if (success) {
        console.log('[SUCCESS] Closer set via Safe');
      } else {
        console.log(
          '[WARN] Safe execution skipped or failed. Proposal must be signed.'
        );
      }
    } else {
      require(
        ownerToUse == _deployerAddress,
        'D1: caller is not owner for setCloser'
      );
      console.log('[ACTION] Executing direct setCloser via EOA...');
      (bool success, ) = _bridgeProxy.call(setCloserData);
      require(success, 'D1: setCloserAndActive failed');
      console.log('[SUCCESS] Closer set successfully');
    }

    console.log('------------------------------------------\n');
  }

  function _logSafeTx(
    IGnosisSafe _safe,
    address _target,
    bytes memory _data,
    string memory _label
  ) internal view {
    uint256 nonce = _safe.nonce();
    bytes32 txHash = _safe.getTransactionHash(
      _target,
      0,
      _data,
      Enum.Operation.Call,
      0,
      0,
      0,
      address(0),
      address(0),
      nonce
    );
    console.log('--- Safe TX Preview ---');
    console.log('Label:', _label);
    console.log('Safe:', address(_safe));
    console.log('Target:', _target);
    console.log('Nonce:', vm.toString(nonce));
    console.log('SafeTxHash:');
    console.logBytes32(txHash);
  }

  function _execViaSafe(
    IGnosisSafe _safe,
    address _target,
    bytes memory _data,
    string memory _label
  ) internal returns (bool executed_) {
    uint256 threshold = _safe.getThreshold();
    address caller = vm.addr(vm.envUint('PRIVATE_KEY'));

    if (threshold != 1) {
      console.log('Safe execution skipped (threshold > 1) for:', _label);
      return false;
    }

    if (!_safe.isOwner(caller)) {
      console.log('Safe execution skipped (caller not owner) for:', _label);
      return false;
    }

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
  }

  function _isContract(address _addr) internal view returns (bool) {
    return _addr.code.length > 0;
  }

  function _step6_deployStorage() internal {
    console.log('------------------------------------------');
    console.log('[STEP 6] Deploying GenFWStorage Contract');
    console.log('------------------------------------------');

    GenFWStorage storage1 = new GenFWStorage();
    deployedStorage = address(storage1);

    console.log('[INFO] Deployed GenFWStorage at:', deployedStorage);
    console.log('[INFO] Registering hashes from snapshot data...');

    uint256 totalClaims = 0;
    // Manual parsing to avoid abi.decode issues
    uint256 tokenCount = _countTokens(assetsJsonContent);

    for (uint256 i = 0; i < tokenCount; i++) {
      string memory dataPath = string.concat('[', vm.toString(i), '].data');
      bytes memory dataRaw = vm.parseJson(assetsJsonContent, dataPath);

      // We can still use abi.decode for simple arrays if they are consistent,
      // but let is be even safer and use count
      uint256 claimCount = _countClaimsInToken(assetsJsonContent, i);

      for (uint256 j = 0; j < claimCount; j++) {
        string memory hashPath = string.concat(
          '[',
          vm.toString(i),
          '].data[',
          vm.toString(j),
          '].hash'
        );
        string memory hashString = vm.parseJsonString(
          assetsJsonContent,
          hashPath
        );

        bytes32 hashValue = vm.parseBytes32(hashString);
        bytes4 sig = _hashToSelector(hashString);
        storage1.setHash(sig, hashValue);
        totalClaims++;
      }
    }

    console.log('[SUCCESS] Total hashes registered:', totalClaims);
    console.log('------------------------------------------\n');
  }

  function _countTokens(string memory json) internal pure returns (uint256) {
    // Crude but effective way to count objects in array for this specific format
    // Search for "l1Token" occurrences
    return _countOccurrences(json, '"l1Token"');
  }

  function _countClaimsInToken(
    string memory json,
    uint256 tokenIdx
  ) internal pure returns (uint256) {
    // Each claim has a "hash" key.
    // We extract the chunk of JSON for this token and count hashes.
    // Searching for unique "hash" keys within data array
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
          // Found our token's start. Now count patterns until END of this token or next l1Token
          for (
            uint256 j = i + bL1Token.length;
            j <= bJson.length - bPattern.length;
            j++
          ) {
            // If we find next token's l1Token, stop
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
    uint256 start,
    bytes memory pattern
  ) internal pure returns (bool) {
    for (uint256 i = 0; i < pattern.length; i++) {
      if (data[start + i] != pattern[i]) return false;
    }
    return true;
  }

  function _step7_registerStorage(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 7] Registering Storage with Bridge');
    console.log('------------------------------------------');

    ForceWithdrawBridge bridge = ForceWithdrawBridge(payable(_bridgeProxy));

    address[] memory positions = new address[](1);
    positions[0] = deployedStorage;

    console.log('[INFO] Registering position:', deployedStorage);
    bridge.forceRegistry(positions);

    // Verify registration
    bool isRegistered = bridge.position(deployedStorage);
    require(isRegistered, 'D1: storage registration failed');

    console.log('[SUCCESS] Storage contract registered successfully');
    console.log('------------------------------------------\n');
  }

  function _step8_activateForceWithdraw(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 8] Activating Force Withdrawal Mode');
    console.log('------------------------------------------');

    ForceWithdrawBridge bridge = ForceWithdrawBridge(payable(_bridgeProxy));

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

  function _step10_executeClaims(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 10] Executing Force Withdrawal Claims');
    console.log('------------------------------------------');

    ForceWithdrawBridge bridge = ForceWithdrawBridge(payable(_bridgeProxy));
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

  // ============================================
  // Helper Functions
  // ============================================

  function _validateInputs(
    address _bridgeProxy,
    address _proxyAdmin,
    string memory _assetsPath
  ) internal pure {
    require(_bridgeProxy != address(0), 'D1: bridge proxy required');
    require(_proxyAdmin != address(0), 'D1: proxy admin required');
    require(bytes(_assetsPath).length > 0, 'D1: assets path required');
  }

  function _loadAssetsJson(string memory _assetsPath) internal {
    console.log('[INFO] Loading assets from:', _assetsPath);

    string memory json = vm.readFile(_assetsPath);
    require(bytes(json).length > 0, 'D1: assets json empty');
    assetsJsonContent = json;

    console.log('[SUCCESS] Assets file loaded successfully\n');
  }

  function _printSummary(
    bool _executedClaims,
    string memory _assetsPath
  ) internal {
    console.log('==============================================');
    console.log('ENABLE L1 WITHDRAWAL: SUMMARY');
    console.log('==============================================');
    console.log('Implementation deployed:', deployedImpl);
    console.log('GenFWStorage deployed:', deployedStorage);
    console.log('Force withdrawal mode: ACTIVE');

    if (_executedClaims) {
      uint256 totalClaims = _countOccurrences(assetsJsonContent, '"hash"');
      console.log('Claims executed:', totalClaims);
    } else {
      console.log('Claims executed: SKIPPED (execute_claims=false)');
    }
    console.log('Assets Path:', _assetsPath);
    console.log('==============================================\n');
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
      _execViaSafe(IGnosisSafe(ownerToUse), target, data, label);
      return;
    }

    require(ownerToUse == deployerAddress, 'D1: caller is not owner');
    (bool success, ) = target.call(data);
    require(success, 'D1: sweep call failed');
  }

  function _upgradeProxyWithSafe(
    address _proxy,
    address _proxyAdmin,
    address _systemOwnerSafe,
    address _newImpl,
    string memory _label
  ) internal {
    console.log('[INFO] Proxy:', _proxy);
    console.log('[INFO] Proxy Admin:', _proxyAdmin);
    console.log('[INFO] New Implementation:', _newImpl);

    bytes memory upgradeData = abi.encodeWithSignature(
      'upgrade(address,address)',
      _proxy,
      _newImpl
    );

    address adminOwner = ProxyAdmin(_proxyAdmin).owner();
    address ownerToUse = _systemOwnerSafe != address(0)
      ? _systemOwnerSafe
      : adminOwner;
    console.log('[INFO] ProxyAdmin Owner:', adminOwner);

    if (_isContract(ownerToUse)) {
      console.log('[INFO] Owner is a contract (Safe). Preparing Safe TX...');
      IGnosisSafe safe = IGnosisSafe(ownerToUse);

      _logSafeTx(safe, _proxyAdmin, upgradeData, _label);

      bool success = _execViaSafe(safe, _proxyAdmin, upgradeData, _label);
      if (success) {
        console.log('[SUCCESS] Proxy upgraded via Safe');
      } else {
        console.log(
          '[WARN] Safe execution skipped or failed. Proposal must be signed.'
        );
      }
    } else {
      require(
        ownerToUse == vm.addr(vm.envUint('PRIVATE_KEY')),
        'D1: caller is not owner'
      );
      console.log('[ACTION] Executing direct upgrade via EOA...');
      (bool success, ) = _proxyAdmin.call(upgradeData);
      require(success, 'D1: proxy upgrade failed');
      console.log('[SUCCESS] Proxy upgraded successfully');
    }
  }

  function _upgradeProxyByEoa(address _proxy, address _newImpl) internal {
    bytes memory upgradeData = abi.encodeWithSignature(
      'upgradeTo(address)',
      _newImpl
    );
    (bool success, ) = _proxy.call(upgradeData);
    require(success, 'D1: proxy upgradeTo failed');
    console.log('[SUCCESS] Proxy upgraded via admin EOA');
  }

  function _getEip1967Admin(
    address _proxy
  ) internal view returns (address admin_) {
    bytes32 slot = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;
    bytes32 data = vm.load(_proxy, slot);
    admin_ = address(uint160(uint256(data)));
  }

  function _hashToSelector(
    string memory hashString
  ) internal pure returns (bytes4) {
    string memory stripped = _strip0x(hashString);
    return bytes4(keccak256(abi.encodePacked('_', stripped, '()')));
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
}
