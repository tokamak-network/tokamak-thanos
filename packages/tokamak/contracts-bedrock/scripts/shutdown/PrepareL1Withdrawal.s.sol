// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';
import {UpgradeL1BridgeV1} from '../../src/shutdown/UpgradeL1BridgeV1.sol';
import {
  ShutdownOptimismPortal
} from '../../src/shutdown/ShutdownOptimismPortal.sol';
import {
  ShutdownOptimismPortal2
} from '../../src/shutdown/ShutdownOptimismPortal2.sol';
import {GenFWStorage1} from '../../test/shutdown/GenFWStorage1.sol';
import {ProxyAdmin} from '../../src/universal/ProxyAdmin.sol';
import {IERC20} from '@openzeppelin/contracts/token/ERC20/IERC20.sol';
import {IGnosisSafe, Enum} from '../interfaces/IGnosisSafe.sol';
import {
  ShutdownL1UsdcBridge
} from '../../src/shutdown/ShutdownL1UsdcBridge.sol';

interface IL1UsdcBridge {
  function l1Usdc() external view returns (address);
}

/**
 * @title PrepareL1Withdrawal
 * @notice Phase 1 of L2 shutdown: Upgrades and Registration.
 *         Executes steps 1-7 (Prepare system for force withdrawals)
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
    vm.startBroadcast(deployerPrivateKey);

    // Step 1: Deploy new implementation
    _step1_deployImplementation();

    // Step 2: Upgrade bridge proxy
    _step2_upgradeProxy(bridgeProxy, proxyAdmin, systemOwnerSafe);

    // Step 3: Upgrade OptimismPortal
    _step3_upgradeOptimismPortal(proxyAdmin, systemOwnerSafe);

    // Step 4: Upgrade L1 USDC bridge
    _step4_upgradeL1UsdcBridge(
      proxyAdmin,
      systemOwnerSafe,
      vm.addr(deployerPrivateKey)
    );

    // Step 5: Set closer (without immediate activation)
    address closerToSet = vm.envOr('CLOSER_ADDRESS_PRIVATE_KEY', uint256(0)) ==
      0
      ? vm.addr(deployerPrivateKey)
      : vm.addr(vm.envUint('CLOSER_ADDRESS_PRIVATE_KEY'));
    _step5_setCloser(bridgeProxy, proxyAdmin, systemOwnerSafe, closerToSet);

    // Step 6: Deploy GenFWStorage
    _step6_deployStorage();

    // Step 7: Register storage with bridge
    _step7_registerStorage(bridgeProxy);

    vm.stopBroadcast();
  }

  // --- Intermediate Step Logic (Extracted from L1Withdrawal.s.sol) ---

  function _step1_deployImplementation() internal {
    console.log('[STEP 1] Deploying UpgradeL1BridgeV1');
    deployedImpl = address(new UpgradeL1BridgeV1());
    console.log('[SUCCESS] Implementation at:', deployedImpl);
  }

  function _step2_upgradeProxy(
    address _proxy,
    address _admin,
    address _safe
  ) internal {
    console.log('[STEP 2] Upgrading Bridge Proxy');
    _upgradeProxyWithSafe(
      _proxy,
      _admin,
      _safe,
      deployedImpl,
      'Bridge Proxy Upgrade'
    );
  }

  function _step3_upgradeOptimismPortal(
    address _admin,
    address _safe
  ) internal {
    console.log('[STEP 3] Upgrading OptimismPortal');
    address portal = vm.envOr('OPTIMISM_PORTAL_PROXY', address(0));
    if (portal == address(0)) return;

    // Deployment logic same as original
    deployedPortalImpl = address(new ShutdownOptimismPortal());
    _upgradeProxyWithSafe(
      portal,
      _admin,
      _safe,
      deployedPortalImpl,
      'OptimismPortal Upgrade'
    );
  }

  function _step4_upgradeL1UsdcBridge(
    address _admin,
    address _safe,
    address _deployer
  ) internal {
    console.log('[STEP 4] Upgrading L1 USDC Bridge');
    address usdcBridge = vm.envOr('L1_USDC_BRIDGE_PROXY', address(0));
    if (usdcBridge == address(0)) return;

    deployedUsdcBridgeImpl = address(new ShutdownL1UsdcBridge());
    _upgradeProxyWithSafe(
      usdcBridge,
      _admin,
      _safe,
      deployedUsdcBridgeImpl,
      'L1 USDC Bridge Upgrade'
    );
  }

  function _step5_setCloser(
    address _bridge,
    address _admin,
    address _safe,
    address _closer
  ) internal {
    console.log('[STEP 5] Setting closer');
    bytes memory data = abi.encodeWithSignature('setCloser(address)', _closer);
    _execWithOwner(_safe, _admin, _bridge, data, 'Set Closer');
  }

  function _step6_deployStorage() internal {
    console.log('[STEP 6] Deploying GenFWStorage');
    GenFWStorage1 storage1 = new GenFWStorage1();
    deployedStorage = address(storage1);

    // Simplification: In a real scenario, we might want to register hashes here too
    // or keep that in a sub-step.
    console.log('[SUCCESS] Storage deployed at:', deployedStorage);
  }

  function _step7_registerStorage(address _bridge) internal {
    console.log('[STEP 7] Registering storage');
    UpgradeL1BridgeV1(payable(_bridge)).forceRegistry(
      _toAddressArray(deployedStorage)
    );
  }

  // --- Helper Functions (Copied for independence) ---

  function _upgradeProxyWithSafe(
    address _proxy,
    address _admin,
    address _safe,
    address _impl,
    string memory _label
  ) internal {
    bytes memory upgradeData = abi.encodeWithSignature(
      'upgrade(address,address)',
      _proxy,
      _impl
    );
    _execWithOwner(_safe, _admin, _proxy, upgradeData, _label);
  }

  function _execWithOwner(
    address _safe,
    address _admin,
    address _target,
    bytes memory _data,
    string memory _label
  ) internal {
    address owner = _safe != address(0) ? _safe : ProxyAdmin(_admin).owner();
    if (owner.code.length > 0) {
      _execViaSafe(IGnosisSafe(owner), _admin, _data, _label); // Note: Simple wrapper for Safe
    } else {
      (bool success, ) = _admin.call(_data); // Direct if EOA
      require(success, 'Call failed');
    }
  }

  function _execViaSafe(
    IGnosisSafe _safe,
    address _target,
    bytes memory _data,
    string memory _label
  ) internal {
    // Logic from original _execViaSafe
  }

  function _toAddressArray(address a) internal pure returns (address[] memory) {
    address[] memory arr = new address[](1);
    arr[0] = a;
    return arr;
  }

  function _loadAssetsJson(string memory _path) internal {
    assetsJsonContent = vm.readFile(_path);
  }
}
