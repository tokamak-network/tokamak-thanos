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
import {IGnosisSafe, Enum} from '../interfaces/IGnosisSafe.sol';
import {
  ShutdownL1UsdcBridge
} from '../../src/shutdown/ShutdownL1UsdcBridge.sol';

import {ShutdownConfig} from './lib/ShutdownConfig.sol';
import {JsonUtils} from './lib/JsonUtils.sol';
import {SafeUtils} from './lib/SafeUtils.sol';

interface IOptimismPortal2 {
  function proofMaturityDelaySeconds() external view returns (uint256);
  function disputeGameFinalityDelaySeconds() external view returns (uint256);
}

/**
 * @title PrepareL1Withdrawal
 * @notice Phase 1 of L2 shutdown: upgrades, registration, and activation (steps 1-8).
 */
contract PrepareL1Withdrawal is Script {
  using stdJson for string;

  address public deployedImpl;
  address public deployedPortalImpl;
  address public deployedUsdcBridgeImpl;
  address public deployedStorage;
  string public assetsJsonContent;
  uint256 public safeNonceOffset;

  function run() public {
    address bridgeProxy = vm.envAddress('BRIDGE_PROXY');
    address proxyAdmin = vm.envAddress('PROXY_ADMIN');
    string memory dataPath = vm.envString('DATA_PATH');
    address systemOwnerSafe = vm.envOr('SYSTEM_OWNER_SAFE', address(0));

    _loadAssetsJson(dataPath);

    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');
    address deployerAddress = vm.addr(deployerPrivateKey);

    uint256 closerPrivateKey = vm.envOr(
      'CLOSER_ADDRESS_PRIVATE_KEY',
      uint256(0)
    );
    if (closerPrivateKey == 0) {
      closerPrivateKey = deployerPrivateKey;
    }
    address closerAddress = vm.addr(closerPrivateKey);

    vm.startBroadcast(deployerPrivateKey);

    // Step 1: Deploy new implementation
    _step1_deployImplementation();

    // Step 2: Upgrade bridge proxy
    _step2_upgradeProxy(bridgeProxy, proxyAdmin, systemOwnerSafe);

    // Step 3: Upgrade OptimismPortal
    _step3_upgradeOptimismPortal(proxyAdmin, systemOwnerSafe);

    // Step 4: Upgrade L1 USDC bridge
    _step4_upgradeL1UsdcBridge(proxyAdmin, systemOwnerSafe, deployerAddress);

    // Step 5: Set closer
    _step5_setCloser(
      bridgeProxy,
      proxyAdmin,
      systemOwnerSafe,
      closerAddress,
      deployerAddress
    );

    // Step 6: Deploy GenFWStorage + register hashes
    _step6_deployStorage();

    vm.stopBroadcast();

    vm.startBroadcast(closerPrivateKey);

    // Step 7: Register storage with bridge
    _step7_registerStorage(bridgeProxy);

    // Step 8: Activate force withdrawal mode
    _step8_activate(bridgeProxy, closerAddress);

    vm.stopBroadcast();
  }

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

    require(
      _systemOwnerSafe != address(0),
      'D1: SYSTEM_OWNER_SAFE required for OptimismPortal upgrade'
    );
    require(
      SafeUtils.isContract(_systemOwnerSafe),
      'D1: SYSTEM_OWNER_SAFE must be a contract (Safe)'
    );

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

    require(
      _systemOwnerSafe != address(0),
      'D1: SYSTEM_OWNER_SAFE required for L1 USDC bridge upgrade'
    );
    require(
      SafeUtils.isContract(_systemOwnerSafe),
      'D1: SYSTEM_OWNER_SAFE must be a contract (Safe)'
    );

    address l1UsdcBridge = vm.envOr('L1_USDC_BRIDGE_PROXY', address(0));
    address l1UsdcBridgeAdmin = _getEip1967Admin(l1UsdcBridge);
    if (l1UsdcBridge == address(0)) {
      console.log('[WARN] L1_USDC_BRIDGE_PROXY not set, skipping');
      console.log('------------------------------------------\n');
      return;
    }

    require(l1UsdcBridgeAdmin != address(0), 'L1 USDC bridge admin not found');

    bool isDryRun = vm.envOr('DRY_RUN', false);
    if (l1UsdcBridgeAdmin != _deployerAddress && !isDryRun) {
      revert('L1 USDC bridge admin must be deployer (EOA)');
    }

    ShutdownL1UsdcBridge impl = new ShutdownL1UsdcBridge();
    deployedUsdcBridgeImpl = address(impl);

    console.log('[INFO] L1 USDC Bridge Proxy:', l1UsdcBridge);
    console.log('[INFO] L1 USDC Bridge Admin (EOA):', l1UsdcBridgeAdmin);
    console.log('[INFO] New Implementation:', deployedUsdcBridgeImpl);

    _upgradeProxyByEoa(l1UsdcBridge, deployedUsdcBridgeImpl);

    console.log('------------------------------------------\n');
  }

  function _step5_setCloser(
    address _bridgeProxy,
    address _proxyAdmin,
    address _systemOwnerSafe,
    address _closerAddress,
    address _deployerAddress
  ) internal {
    console.log('------------------------------------------');
    console.log('[STEP 5] Setting closer');
    console.log('------------------------------------------');

    // Check if closer is already set to the same address
    (bool success, bytes memory data) = _bridgeProxy.staticcall(
      abi.encodeWithSignature('closer()')
    );
    if (success && data.length >= 32) {
      address currentCloser = abi.decode(data, (address));
      console.log('[INFO] Current closer:', currentCloser);
      if (currentCloser == _closerAddress) {
        console.log('[INFO] Closer already set to the same address, skipping');
        console.log('------------------------------------------\n');
        return;
      }
    }

    bytes memory setCloserData = abi.encodeWithSignature(
      'setCloser(address)',
      _closerAddress
    );

    // Always use the actual ProxyAdmin owner
    address adminOwner = ProxyAdmin(_proxyAdmin).owner();

    console.log('[INFO] Bridge Proxy:', _bridgeProxy);
    console.log('[INFO] Closer Address:', _closerAddress);
    console.log('[INFO] ProxyAdmin Owner:', adminOwner);

    if (SafeUtils.isContract(adminOwner)) {
      console.log('[INFO] Owner is a contract (Safe). Preparing Safe TX...');
      IGnosisSafe safe = IGnosisSafe(adminOwner);

      SafeUtils.logSafeTxWithNonce(
        safe,
        _bridgeProxy,
        setCloserData,
        'Set Closer',
        safeNonceOffset
      );

      bool success = SafeUtils.execViaSafeFromEnvWithNonce(
        safe,
        _bridgeProxy,
        setCloserData,
        'Set Closer',
        safeNonceOffset
      );
      if (success) {
        console.log('[SUCCESS] Closer set via Safe');
        safeNonceOffset++;
      } else {
        console.log(
          '[WARN] Safe execution skipped (Threshold > 1 or simulation). Proposing TX...'
        );
      }
    } else {
      bool isDryRun = vm.envOr('DRY_RUN', false);
      if (adminOwner != _deployerAddress && !isDryRun) {
        revert('D1: caller is not owner for setCloser');
      }
      console.log('[ACTION] Executing direct setCloser via EOA...');
      (bool success, ) = _bridgeProxy.call(setCloserData);
      require(success, 'D1: setCloser failed');
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
    uint256 tokenCount = JsonUtils.countTokens(assetsJsonContent);

    for (uint256 i = 0; i < tokenCount; i++) {
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

    bool isRegistered = bridge.position(deployedStorage);
    require(isRegistered, 'D1: storage registration failed');

    console.log('[SUCCESS] Storage contract registered successfully');
    console.log('------------------------------------------\n');
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

    // Always use the actual ProxyAdmin owner
    address adminOwner = ProxyAdmin(_proxyAdmin).owner();
    console.log('[INFO] ProxyAdmin Owner:', adminOwner);

    if (SafeUtils.isContract(adminOwner)) {
      console.log('[INFO] Owner is a contract (Safe). Preparing Safe TX...');
      IGnosisSafe safe = IGnosisSafe(adminOwner);

      SafeUtils.logSafeTxWithNonce(
        safe,
        _proxyAdmin,
        upgradeData,
        _label,
        safeNonceOffset
      );

      bool success = SafeUtils.execViaSafeFromEnvWithNonce(
        safe,
        _proxyAdmin,
        upgradeData,
        _label,
        safeNonceOffset
      );
      if (success) {
        console.log('[SUCCESS] Proxy upgraded via Safe');
        safeNonceOffset++;
      } else {
        console.log(
          '[WARN] Safe execution skipped (Threshold > 1 or simulation). Proposing TX...'
        );
      }
    } else {
      bool isDryRun = vm.envOr('DRY_RUN', false);
      if (adminOwner != vm.addr(vm.envUint('PRIVATE_KEY')) && !isDryRun) {
        revert('D1: caller is not owner');
      }
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

  function _loadAssetsJson(string memory _path) internal {
    assetsJsonContent = vm.readFile(_path);
  }
}
