// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {UpgradeL1BridgeV1} from '../../src/shutdown/ForceWithdrawBridge.sol';
import {
  ShutdownOptimismPortal
} from '../../src/shutdown/ShutdownOptimismPortal.sol';
import {
  ShutdownOptimismPortal2
} from '../../src/shutdown/ShutdownOptimismPortal2.sol';
import {GenFWStorage1} from '../../test/shutdown/GenFWStorage1.sol';
import {ProxyAdmin} from '../../src/universal/ProxyAdmin.sol';
import {IGnosisSafe, Enum} from '../interfaces/IGnosisSafe.sol';
import {
  ShutdownL1UsdcBridge
} from '../../src/shutdown/ShutdownL1UsdcBridge.sol';

interface IOptimismPortal2 {
  function proofMaturityDelaySeconds() external view returns (uint256);
  function disputeGameFinalityDelaySeconds() external view returns (uint256);
}

/**
 * @title PrepareL1Withdrawal
 * @notice Phase 1 of L2 shutdown: upgrades and registration (steps 1-7).
 */
contract PrepareL1Withdrawal is Script {
  using stdJson for string;

  address public deployedImpl;
  address public deployedPortalImpl;
  address public deployedUsdcBridgeImpl;
  address public deployedStorage;
  string public assetsJsonContent;

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

    // Step 7: Activate force withdrawal mode
    _step7_activateForceWithdraw(bridgeProxy);

    vm.stopBroadcast();
  }

  function _step1_deployImplementation() internal {
    console.log('------------------------------------------');
    console.log('[STEP 1] Deploying UpgradeL1BridgeV1 Implementation');
    console.log('------------------------------------------');

    UpgradeL1BridgeV1 impl = new UpgradeL1BridgeV1();
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
      _isContract(_systemOwnerSafe),
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
      _isContract(_systemOwnerSafe),
      'D1: SYSTEM_OWNER_SAFE must be a contract (Safe)'
    );

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
    bool useProxyAdmin = false;
    address proxyAdminToUse = address(0);
    if (adminFromSlot != address(0)) {
      if (
        l1UsdcBridgeProxyAdmin != address(0) &&
        adminFromSlot == l1UsdcBridgeProxyAdmin
      ) {
        useProxyAdmin = true;
        proxyAdminToUse = l1UsdcBridgeProxyAdmin;
      } else if (adminFromSlot == _proxyAdmin) {
        useProxyAdmin = true;
        proxyAdminToUse = _proxyAdmin;
      }
    } else if (l1UsdcBridgeProxyAdmin != address(0)) {
      useProxyAdmin = true;
      proxyAdminToUse = l1UsdcBridgeProxyAdmin;
    } else {
      useProxyAdmin = true;
      proxyAdminToUse = _proxyAdmin;
    }

    if (
      l1UsdcBridgeAdmin != address(0) && l1UsdcBridgeAdmin == _deployerAddress
    ) {
      console.log('[INFO] L1 USDC bridge admin matches deployer');
      console.log('[INFO] Safe upgrade enforced; EOA path disabled');
    }

    if (!useProxyAdmin) {
      address adminToUse = l1UsdcBridgeAdmin != address(0)
        ? l1UsdcBridgeAdmin
        : adminFromSlot;
      require(adminToUse != address(0), 'D1: L1 USDC bridge admin not set');

      bool isDryRun = vm.envOr('DRY_RUN', false);
      if (adminToUse != _deployerAddress && !isDryRun) {
        revert('D1: caller is not admin');
      }

      _upgradeProxyDirectWithSafe(
        l1UsdcBridge,
        adminToUse,
        deployedUsdcBridgeImpl,
        _deployerAddress,
        'L1 USDC Bridge Upgrade'
      );
    }

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

    bytes memory setCloserData = abi.encodeWithSignature(
      'setCloser(address)',
      _closerAddress
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

      _logSafeTx(safe, _bridgeProxy, setCloserData, 'Set Closer');

      bool success = _execViaSafe(
        safe,
        _bridgeProxy,
        setCloserData,
        'Set Closer'
      );
      if (success) {
        console.log('[SUCCESS] Closer set via Safe');
      } else {
        console.log(
          '[WARN] Safe execution skipped (Threshold > 1 or simulation). Proposing TX...'
        );
      }
    } else {
      bool isDryRun = vm.envOr('DRY_RUN', false);
      if (ownerToUse != _deployerAddress && !isDryRun) {
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

    GenFWStorage1 storage1 = new GenFWStorage1();
    deployedStorage = address(storage1);

    console.log('[INFO] Deployed GenFWStorage1 at:', deployedStorage);
    console.log('[INFO] Registering hashes from snapshot data...');

    uint256 totalClaims = 0;
    uint256 tokenCount = _countTokens(assetsJsonContent);

    for (uint256 i = 0; i < tokenCount; i++) {
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

  function _step7_registerStorage(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 6] Registering Storage with Bridge');
    console.log('------------------------------------------');

    UpgradeL1BridgeV1 bridge = UpgradeL1BridgeV1(payable(_bridgeProxy));

    address[] memory positions = new address[](1);
    positions[0] = deployedStorage;

    console.log('[INFO] Registering position:', deployedStorage);
    bridge.forceRegistry(positions);

    bool isRegistered = bridge.position(deployedStorage);
    require(isRegistered, 'D1: storage registration failed');

    console.log('[SUCCESS] Storage contract registered successfully');
    console.log('------------------------------------------\n');
  }

  function _step7_activateForceWithdraw(address _bridgeProxy) internal {
    console.log('------------------------------------------');
    console.log('[STEP 7] Activating Force Withdrawal Mode');
    console.log('------------------------------------------');

    UpgradeL1BridgeV1 bridge = UpgradeL1BridgeV1(payable(_bridgeProxy));

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
          '[WARN] Safe execution skipped (Threshold > 1 or simulation). Proposing TX...'
        );
      }
    } else {
      bool isDryRun = vm.envOr('DRY_RUN', false);
      if (ownerToUse != vm.addr(vm.envUint('PRIVATE_KEY')) && !isDryRun) {
        revert('D1: caller is not owner');
      }
      console.log('[ACTION] Executing direct upgrade via EOA...');
      (bool success, ) = _proxyAdmin.call(upgradeData);
      require(success, 'D1: proxy upgrade failed');
      console.log('[SUCCESS] Proxy upgraded successfully');
    }
  }

  function _upgradeProxyDirectWithSafe(
    address _proxy,
    address _admin,
    address _newImpl,
    address _deployerAddress,
    string memory _label
  ) internal {
    console.log('[INFO] Proxy:', _proxy);
    console.log('[INFO] Admin (Safe):', _admin);
    console.log('[INFO] New Implementation:', _newImpl);

    bytes memory upgradeData = abi.encodeWithSignature(
      'upgradeTo(address)',
      _newImpl
    );

    if (_isContract(_admin)) {
      IGnosisSafe safe = IGnosisSafe(_admin);
      _logSafeTx(safe, _proxy, upgradeData, _label);

      bool success = _execViaSafe(safe, _proxy, upgradeData, _label);
      require(success, 'D1: safe execTransaction failed');
      console.log('[SUCCESS] Proxy upgraded via Safe (direct admin)');
      return;
    }

    require(_admin == _deployerAddress, 'D1: caller is not admin');
    console.log('[ACTION] Executing direct upgrade via EOA...');
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

  function _loadAssetsJson(string memory _path) internal {
    assetsJsonContent = vm.readFile(_path);
  }
}
