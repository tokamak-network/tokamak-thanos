// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { ISemver } from "src/universal/ISemver.sol";
import { ResourceMetering } from "src/L1/ResourceMetering.sol";
import { Storage } from "src/libraries/Storage.sol";
import { Constants } from "src/libraries/Constants.sol";

interface ISystemConfig {

    enum UpdateType {
        BATCHER,
        GAS_CONFIG,
        GAS_LIMIT,
        UNSAFE_BLOCK_SIGNER
    }

    struct Addresses {
        address l1CrossDomainMessenger;
        address l1ERC721Bridge;
        address l1StandardBridge;
        address l2OutputOracle;
        address optimismPortal;
        address optimismMintableERC20Factory;
    }

    // uint256 public constant VERSION = 0;

    // bytes32 public constant UNSAFE_BLOCK_SIGNER_SLOT = keccak256("systemconfig.unsafeblocksigner");

    // bytes32 public constant L1_CROSS_DOMAIN_MESSENGER_SLOT =
    //     bytes32(uint256(keccak256("systemconfig.l1crossdomainmessenger")) - 1);

    // bytes32 public constant L1_ERC_721_BRIDGE_SLOT = bytes32(uint256(keccak256("systemconfig.l1erc721bridge")) - 1);

    // bytes32 public constant L1_STANDARD_BRIDGE_SLOT = bytes32(uint256(keccak256("systemconfig.l1standardbridge")) - 1);

    // bytes32 public constant L2_OUTPUT_ORACLE_SLOT = bytes32(uint256(keccak256("systemconfig.l2outputoracle")) - 1);

    // bytes32 public constant OPTIMISM_PORTAL_SLOT = bytes32(uint256(keccak256("systemconfig.optimismportal")) - 1);

    // bytes32 public constant OPTIMISM_MINTABLE_ERC20_FACTORY_SLOT =
    //     bytes32(uint256(keccak256("systemconfig.optimismmintableerc20factory")) - 1);

    // bytes32 public constant BATCH_INBOX_SLOT = bytes32(uint256(keccak256("systemconfig.batchinbox")) - 1);

    // uint256 public overhead;

    // uint256 public scalar;

    // bytes32 public batcherHash;

    // uint64 public gasLimit;

    // uint256 public startBlock;

    // string public constant version = "1.10.0";
   function initialize(
        address _owner,
        uint256 _overhead,
        uint256 _scalar,
        bytes32 _batcherHash,
        uint64 _gasLimit,
        address _unsafeBlockSigner,
        ResourceMetering.ResourceConfig memory _config,
        uint256 _startBlock,
        address _batchInbox,
        ISystemConfig.Addresses memory _addresses
    )
        external;

    function minimumGasLimit() external view returns (uint64)   ;

    function unsafeBlockSigner() external view returns (address addr_) ;

    /// @notice Getter for the L1CrossDomainMessenger address.
    function l1CrossDomainMessenger() external view returns (address addr_) ;

    /// @notice Getter for the L1ERC721Bridge address.
    function l1ERC721Bridge() external view returns (address addr_)  ;

    /// @notice Getter for the L1StandardBridge address.
    function l1StandardBridge() external view returns (address addr_) ;

    /// @notice Getter for the L2OutputOracle address.
    function l2OutputOracle() external view returns (address addr_) ;

    /// @notice Getter for the OptimismPortal address.
    function optimismPortal() external view returns (address addr_) ;

    /// @notice Getter for the OptimismMintableERC20Factory address.
    function optimismMintableERC20Factory() external view returns (address addr_) ;

    /// @notice Getter for the BatchInbox address.
    function batchInbox() external view returns (address addr_);

    function setUnsafeBlockSigner(address _unsafeBlockSigner) external  ;

    /// @notice Updates the batcher hash. Can only be called by the owner.
    /// @param _batcherHash New batcher hash.
    function setBatcherHash(bytes32 _batcherHash) external ;

    function setGasConfig(uint256 _overhead, uint256 _scalar) external ;


    function setGasLimit(uint64 _gasLimit) external ;


    function resourceConfig() external view returns (ResourceMetering.ResourceConfig memory);

    function setResourceConfig(ResourceMetering.ResourceConfig memory _config) external ;

}
