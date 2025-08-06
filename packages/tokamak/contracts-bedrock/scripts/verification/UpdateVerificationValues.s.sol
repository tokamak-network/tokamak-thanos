// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {Script, console} from 'forge-std/Script.sol';
import {L1ContractVerification} from '../../src/tokamak-contracts/verification/L1ContractVerification.sol';
import {IMultiSigWallet} from '../../src/tokamak-contracts/verification/interface/IMultiSigWallet.sol';

interface IGnosisSafe {
  function masterCopy() external view returns (address);
}

contract L1ContractVerificationMultisigConfigurationCalls is Script {
  // Address of the deployed L1ContractVerification contract
  address private _proxyAddress;
  address private _verificationContractProxyAdminAddress;
  uint256 private _multisigOwner1;
  uint256 private _multisigOwner2;
  uint256 private _multisigOwner3;
  address private _multisigWallet;
  address private safeWalletProxy;
  address private safeWalletImplementation;
  L1ContractVerification.SafeWalletInfo private newSafeWallet;
  address private newTokamakDAO;

  // Parameters for setSafeConfig
  address tokamakDAO;
  address foundation;
  uint256 safeThreshold;
  bytes32 safeImplementationCodehash;
  bytes32 safeProxyCodehash;
  uint256 safeOwnersCount;

  // Parameters for setLogicContractInfo
  address systemConfigProxyAddress;
  address proxyAdminAddress;

  function setUp() public {
    _proxyAddress = vm.envAddress('L1_CONTRACT_VERIFICATION_PROXY');
    _verificationContractProxyAdminAddress = vm.envAddress(
      'L1_CONTRACT_VERIFICATION_PROXY_ADMIN'
    );
    _multisigOwner1 = vm.envUint('MULTISIG_OWNER_1');
    _multisigOwner2 = vm.envUint('MULTISIG_OWNER_2');
    _multisigOwner3 = vm.envUint('MULTISIG_OWNER_3');
    _multisigWallet = vm.envAddress('MULTISIG_WALLET');

    tokamakDAO = vm.envAddress('TOKAMAK_DAO_ADDRESS');
    foundation = vm.envAddress('FOUNDATION_ADDRESS');
    safeThreshold = 3;
    safeWalletProxy = vm.envAddress('SAFE_WALLET');
    safeWalletImplementation = IGnosisSafe(safeWalletProxy).masterCopy();
    safeOwnersCount = 3;

    systemConfigProxyAddress = vm.envAddress('SYSTEM_CONFIG_PROXY');
    proxyAdminAddress = vm.envAddress('PROXY_ADMIN_ADDRESS');
  }

  function run() public {
    setUp();

    L1ContractVerification l1Verification = L1ContractVerification(
      _proxyAddress
    );
    IMultiSigWallet multiSig = IMultiSigWallet(_multisigWallet);

    console.log('Preparing to call setSafeConfig via multisig...');
    bytes memory setSafeConfigData = abi.encodeWithSelector(
      L1ContractVerification.setSafeConfig.selector,
      tokamakDAO,
      foundation,
      3,
      safeWalletImplementation.codehash,
      safeWalletProxy.codehash,
      3
    );

    vm.broadcast(_multisigOwner1);
    multiSig.submitTransaction(_proxyAddress, 0, setSafeConfigData);
    uint txIndexSetSafeConfig = multiSig.getTransactionCount() - 1;
    console.log(
      'setSafeConfig transaction submitted by owner 1, txIndex:',
      txIndexSetSafeConfig
    );

    vm.broadcast(_multisigOwner2);
    multiSig.confirmTransaction(txIndexSetSafeConfig);
    console.log('setSafeConfig transaction confirmed by owner 2');

    vm.broadcast(_multisigOwner3);
    multiSig.confirmTransaction(txIndexSetSafeConfig);
    console.log('setSafeConfig transaction confirmed by owner 3');

    vm.broadcast(_multisigOwner1);
    multiSig.executeTransaction(txIndexSetSafeConfig);
    console.log('setSafeConfig transaction executed by owner 1');

    console.log('Preparing to call setLogicContractInfo via multisig...');
    bytes memory setLogicContractInfoData = abi.encodeWithSelector(
      L1ContractVerification.setLogicContractInfo.selector,
      systemConfigProxyAddress,
      proxyAdminAddress
    );

    vm.broadcast(_multisigOwner1);
    multiSig.submitTransaction(_proxyAddress, 0, setLogicContractInfoData);
    uint txIndexSetLogic = multiSig.getTransactionCount() - 1;
    console.log(
      'setLogicContractInfo transaction submitted by owner 1, txIndex:',
      txIndexSetLogic
    );

    vm.broadcast(_multisigOwner2);
    multiSig.confirmTransaction(txIndexSetLogic);
    console.log('setLogicContractInfo transaction confirmed by owner 2');

    vm.broadcast(_multisigOwner3);
    multiSig.confirmTransaction(txIndexSetLogic);
    console.log('setLogicContractInfo transaction confirmed by owner 3');

    vm.broadcast(_multisigOwner1);
    multiSig.executeTransaction(txIndexSetLogic);
    console.log('setLogicContractInfo transaction executed by owner 1');


    console.log('Verification factor update completed');
  }
}
