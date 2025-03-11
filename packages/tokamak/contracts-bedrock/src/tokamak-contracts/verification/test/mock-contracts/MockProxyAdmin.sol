// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import '../../interface/IProxyAdmin.sol';

contract MockProxyAdmin is IProxyAdmin {
  mapping(address => address) private implementations;
  mapping(address => address) private admins;
  mapping(address => ProxyType) private proxyTypes;
  mapping(address => string) private implNames;
  address private _owner;
  bool private _isUpgrading;

  constructor(address owner_) {
    _owner = owner_;
  }

  function getProxyImplementation(address proxy) external view returns (address) {
    // Return the stored implementation or revert if not set
    address impl = implementations[proxy];
    require(impl != address(0), 'Implementation not set for proxy');
    return impl;
  }

  function getProxyAdmin(address payable proxy) external view returns (address) {
    // Return the stored admin or revert if not set
    address admin = admins[proxy];
    require(admin != address(0), 'Admin not set for proxy');
    return admin;
  }

  function setImplementation(address proxy, address implementation) external {
    require(implementation != address(0), 'Implementation cannot be zero address');
    implementations[proxy] = implementation;
  }

  function setAdmin(address proxy, address admin) external {
    require(admin != address(0), 'Admin cannot be zero address');
    admins[proxy] = admin;
  }

  function owner() external view returns (address) {
    return _owner;
  }

  function setProxyType(address _address, ProxyType _type) external {
    proxyTypes[_address] = _type;
  }

  function setImplementationName(address _address, string memory _name) external {
    implNames[_address] = _name;
  }

  function setUpgrading(bool _upgrading) external {
    _isUpgrading = _upgrading;
  }

  function isUpgrading() external view returns (bool) {
    return _isUpgrading;
  }

  function changeProxyAdmin(address payable _proxy, address _newAdmin) external {
    admins[_proxy] = _newAdmin;
  }

  function upgrade(address payable _proxy, address _implementation) external {
    implementations[_proxy] = _implementation;
  }

  function upgradeAndCall(
    address payable _proxy,
    address _implementation,
    bytes memory _data
  ) external payable {
    implementations[_proxy] = _implementation;
    // We don't actually call the implementation in this mock
  }

  function proxyType(address _address) external view returns (ProxyType) {
    return proxyTypes[_address];
  }

  function implementationName(address _address) external view returns (string memory) {
    return implNames[_address];
  }
}