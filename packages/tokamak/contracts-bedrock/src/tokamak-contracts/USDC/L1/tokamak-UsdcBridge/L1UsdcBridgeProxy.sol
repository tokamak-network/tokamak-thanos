// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import { ERC1967Proxy } from "@openzeppelin/contracts_v5.0.1/proxy/ERC1967/ERC1967Proxy.sol";
import { ERC1967Utils } from "@openzeppelin/contracts_v5.0.1/proxy/ERC1967/ERC1967Utils.sol";
import { L1UsdcBridgeStorage } from "./L1UsdcBridgeStorage.sol";

contract L1UsdcBridgeProxy is L1UsdcBridgeStorage, ERC1967Proxy {
    modifier onlyProxyOwner() {
        require(msg.sender == owner(), "not owner");
        _;
    }

    modifier nonZeroAddress(address addr) {
        require(addr != address(0), "zero address");
        _;
    }

    constructor(address _logic, address initialOwner, bytes memory _data) payable ERC1967Proxy(_logic, _data) {
        ERC1967Utils.changeAdmin(initialOwner);
    }

    receive() external payable {
        revert("cannot receive Ether");
    }

    function proxyChangeOwner(address newAdmin) external onlyProxyOwner {
        ERC1967Utils.changeAdmin(newAdmin);
    }

    function setAddress(
        address _messenger,
        address _otherBridge,
        address _l1Usdc,
        address _l2Usdc
    )
        external
        onlyProxyOwner
        nonZeroAddress(_messenger)
        nonZeroAddress(_otherBridge)
        nonZeroAddress(_l1Usdc)
        nonZeroAddress(_l2Usdc)
    {
        messenger = _messenger;
        otherBridge = _otherBridge;
        l1Usdc = _l1Usdc;
        l2Usdc = _l2Usdc;
    }

    function upgradeTo(address newImplementation) external onlyProxyOwner {
        ERC1967Utils.upgradeToAndCall(newImplementation, bytes(""));
    }

    function upgradeToAndCall(address newImplementation, bytes memory data) external onlyProxyOwner {
        ERC1967Utils.upgradeToAndCall(newImplementation, data);
    }

    function owner() public view returns (address) {
        return ERC1967Utils.getAdmin();
    }

    function implementation() external view returns (address) {
        return _implementation();
    }
}
