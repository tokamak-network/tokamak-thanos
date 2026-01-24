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

import {ShutdownConfig} from './lib/ShutdownConfig.sol';
import {JsonUtils} from './lib/JsonUtils.sol';
import {SafeUtils} from './lib/SafeUtils.sol';

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

    if (SafeUtils.isContract(ownerToUse)) {
      console.log('[INFO] Owner is a contract (Safe). Preparing Safe TX...');
      IGnosisSafe safe = IGnosisSafe(ownerToUse);

      SafeUtils.logSafeTx(safe, _bridgeProxy, setCloserData, 'Set Closer And Active');

      bool success = SafeUtils.execViaSafeFromEnv(
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
    uint256 tokenCount = JsonUtils.countTokens(assetsJsonContent);

    for (uint256 i = 0; i < tokenCount; i++) {
      string memory dataPath = string.concat('[', vm.toString(i), '].data');
      bytes memory dataRaw = vm.parseJson(assetsJsonContent, dataPath);

      // We can still use abi.decode for simple arrays if they are consistent,
      // but let is be even safer and use count
      uint256 claimCount = JsonUtils.countClaimsInToken(assetsJsonContent, i);

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
        bytes4 sig = JsonUtils.hashToSelector(hashString);
        storage1.setHash(sig, hashValue);
        totalClaims++;
      }
    }

    console.log('[SUCCESS] Total hashes registered:', totalClaims);
    console.log('------------------------------------------\n');
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
    console.log('[INFO] Bridge proxy:', _bridgeProxy);
    console.log('[INFO] Owner to Use:', _resolveOwner(_systemOwnerSafe));

    _sweepNativeLiquidity(_bridgeProxy, _systemOwnerSafe, _deployerAddress);
    _sweepUsdcLiquidity(_bridgeProxy, _systemOwnerSafe, _deployerAddress);

    console.log('------------------------------------------\n');
  }

  function _resolveOwner(address _systemOwnerSafe) internal view returns (address) {
    return _systemOwnerSafe != address(0)
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

    uint256 amount = _min(
      remaining,
      _getBalance(l1UsdcBridge, l1UsdcToken)
    );
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

  function _step10_executeClaims(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 10] Executing Force Withdrawal Claims');
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
      uint256 totalClaims = JsonUtils.countOccurrences(assetsJsonContent, '"hash"');
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
      SafeUtils.execViaSafeFromEnv(IGnosisSafe(ownerToUse), target, data, label);
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

    if (SafeUtils.isContract(ownerToUse)) {
      console.log('[INFO] Owner is a contract (Safe). Preparing Safe TX...');
      IGnosisSafe safe = IGnosisSafe(ownerToUse);

      SafeUtils.logSafeTx(safe, _proxyAdmin, upgradeData, _label);

      bool success = SafeUtils.execViaSafeFromEnv(safe, _proxyAdmin, upgradeData, _label);
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
    bytes32 data = vm.load(_proxy, ShutdownConfig.PROXY_ADMIN_SLOT);
    admin_ = address(uint160(uint256(data)));
  }
}
