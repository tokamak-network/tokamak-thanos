// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {OptimismPortalClosing} from 'src/shutdown/OptimismPortalClosing.sol';
import {ProxyAdmin} from 'src/universal/ProxyAdmin.sol';
import {SuperchainConfig} from 'src/L1/SuperchainConfig.sol';
import {IGnosisSafe, Enum} from '../interfaces/IGnosisSafe.sol';

import {SafeUtils} from './lib/SafeUtils.sol';

/**
 * @title BlockDepositsWithdrawals
 * @notice Forge script to block deposits and withdrawals during shutdown.
 */
contract BlockDepositsWithdrawals is Script {
  function run() external {
    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');
    address portalProxyAddr = vm.envAddress('OPTIMISM_PORTAL_PROXY');
    address superchainConfigAddr = vm.envAddress('SUPERCHAIN_CONFIG_PROXY');
    address proxyAdminAddr = vm.envAddress('PROXY_ADMIN');
    address systemOwnerSafeAddr = vm.envAddress('SYSTEM_OWNER_SAFE');
    address guardianSafeAddr = vm.envOr('GUARDIAN_SAFE', systemOwnerSafeAddr);

    address derivedCaller = vm.addr(deployerPrivateKey);
    console.log('Derived caller from PRIVATE_KEY:', derivedCaller);

    vm.startBroadcast(deployerPrivateKey);

    ProxyAdmin proxyAdmin = ProxyAdmin(proxyAdminAddr);
    address currentImplementation = proxyAdmin.getProxyImplementation(
      portalProxyAddr
    );
    console.log('Current implementation:', currentImplementation);

    // 1. Deploy the Closing Implementation
    // 1. Deploy the Closing Implementation
    console.log('Deploying OptimismPortalClosing...');
    OptimismPortalClosing closingImpl = new OptimismPortalClosing{
      salt: bytes32(block.timestamp)
    }(currentImplementation);
    console.log('New Implementation deployed at:', address(closingImpl));

    // 2. Prepare Safe transaction for upgrading the Portal Proxy
    bytes memory upgradeData = abi.encodeCall(
      ProxyAdmin.upgrade,
      (payable(portalProxyAddr), address(closingImpl))
    );
    console.log('Prepared upgrade calldata:', vm.toString(upgradeData));
    if (!SafeUtils.isContract(systemOwnerSafeAddr)) {
      revert('SYSTEM_OWNER_SAFE must be a Safe contract');
    }
    SafeUtils.logSafeOwnership(
      IGnosisSafe(systemOwnerSafeAddr),
      'SYSTEM_OWNER_SAFE',
      derivedCaller
    );
    SafeUtils.logSafeTx(
      IGnosisSafe(systemOwnerSafeAddr),
      proxyAdminAddr,
      upgradeData,
      'Portal upgrade'
    );

    bool upgradeExecuted = SafeUtils.execViaSafeFromEnv(
      IGnosisSafe(systemOwnerSafeAddr),
      proxyAdminAddr,
      upgradeData,
      'Portal upgrade'
    );

    // 3. Prepare Safe transaction for pausing SuperchainConfig
    bytes memory pauseData = abi.encodeCall(
      SuperchainConfig.pause,
      ('Chain shutdown initiated')
    );
    console.log('Prepared pause calldata:', vm.toString(pauseData));
    bool pauseExecuted = false;
    if (SafeUtils.isContract(guardianSafeAddr)) {
      SafeUtils.logSafeOwnership(
        IGnosisSafe(guardianSafeAddr),
        'GUARDIAN_SAFE',
        derivedCaller
      );
      SafeUtils.logSafeTx(
        IGnosisSafe(guardianSafeAddr),
        superchainConfigAddr,
        pauseData,
        'Superchain pause'
      );

      pauseExecuted = SafeUtils.execViaSafeFromEnv(
        IGnosisSafe(guardianSafeAddr),
        superchainConfigAddr,
        pauseData,
        'Superchain pause'
      );
    } else {
      if (derivedCaller != guardianSafeAddr) {
        console.log(
          'GUARDIAN_SAFE is EOA. Caller is not guardian, skip pause.'
        );
      } else {
        console.log('Executing pause via guardian EOA...');
        SuperchainConfig(superchainConfigAddr).pause(
          'Chain shutdown initiated'
        );
        pauseExecuted = true;
      }
    }

    vm.stopBroadcast();

    // 4. Verification Check
    if (upgradeExecuted && pauseExecuted) {
      _verify(portalProxyAddr, superchainConfigAddr, derivedCaller);
    } else {
      console.log('Skipping verification. Safe execution was not completed.');
    }
  }

  function _verify(
    address _portal,
    address _superchain,
    address _caller
  ) internal {
    bool isPaused = SuperchainConfig(_superchain).paused();
    string memory ver = OptimismPortalClosing(payable(_portal)).version();

    console.log('\n--- Verification Report ---');
    console.log('Portal Version:', ver);
    console.log('Superchain Paused:', isPaused ? 'Yes' : 'No');

    require(isPaused, 'Verification Failed: Superchain not paused');
    require(
      keccak256(bytes(ver)) == keccak256(bytes('2.8.1-closing')),
      'Verification Failed: Portal version mismatch'
    );
    _assertDepositBlocked(_portal, _caller);
    _assertReceiveBlocked(_portal);
    _assertOnApproveBlocked(_portal, _caller);
    console.log('All systems successfully shutdown.');
  }

  function _assertDepositBlocked(address _portal, address _caller) internal {
    bytes memory data = abi.encodeWithSignature(
      'depositTransaction(address,uint256,uint256,uint64,bool,bytes)',
      _caller,
      uint256(0),
      uint256(0),
      uint64(21000),
      false,
      bytes('')
    );
    (bool success, bytes memory returndata) = _portal.call(data);
    require(!success, 'Verification Failed: deposit is not blocked');

    _requireShutdownRevert('depositTransaction', returndata);
  }

  function _assertReceiveBlocked(address _portal) internal {
    (bool success, bytes memory returndata) = payable(_portal).call{value: 1}(
      ''
    );
    require(!success, 'Verification Failed: receive is not blocked');
    _requireShutdownRevert('receive', returndata);
  }

  function _assertOnApproveBlocked(address _portal, address _caller) internal {
    bytes memory data = abi.encodeWithSignature(
      'onApprove(address,address,uint256,bytes)',
      _caller,
      _caller,
      uint256(0),
      bytes('')
    );
    (bool success, bytes memory returndata) = _portal.call(data);
    require(!success, 'Verification Failed: onApprove is not blocked');
    _requireShutdownRevert('onApprove', returndata);
  }

  function _requireShutdownRevert(
    string memory _label,
    bytes memory _returndata
  ) internal {
    string memory reason = _decodeRevertReason(_returndata);
    console.log('Deposit block reason for', _label, ':', reason);
    require(
      bytes(reason).length > 0,
      'Verification Failed: missing revert reason'
    );
    require(
      keccak256(bytes(reason)) ==
        keccak256(
          bytes('OptimismPortal: deposits are disabled due to chain shutdown')
        ),
      'Verification Failed: unexpected deposit revert reason'
    );
  }

  function _decodeRevertReason(
    bytes memory _data
  ) internal pure returns (string memory reason) {
    if (_data.length < 4) {
      return '';
    }
    bytes4 selector;
    assembly {
      selector := mload(add(_data, 0x20))
    }
    if (selector != 0x08c379a0) {
      return '';
    }
    assembly {
      _data := add(_data, 0x04)
    }
    reason = abi.decode(_data, (string));
  }
}
