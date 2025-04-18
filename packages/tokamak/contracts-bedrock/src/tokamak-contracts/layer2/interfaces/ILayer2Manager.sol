// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ILayer2Manager {
    // Events
    event SetAddresses(
        address _l2Register,
        address _operatorManagerFactory,
        address _ton,
        address _wton,
        address _dao,
        address _depositManager,
        address _seigManager,
        address _swapProxy
    );
    event SetMinimumInitialDepositAmount(uint256 _minimumInitialDepositAmount);
    event RegisteredCandidateAddOn(address rollupConfig, uint256 wtonAmount, string memo, address operator, address candidateAddOn);
    event PausedCandidateAddOn(address rollupConfig, address candidateAddOn);
    event UnpausedCandidateAddOn(address rollupConfig, address candidateAddOn);
    event SetOperatorManagerFactory(address _operatorManagerFactory);
    event TransferWTON(address layer2, address to, uint256 amount);

    // Admin functions
    function setAddresses(
        address _l1BridgeRegistry,
        address _operatorManagerFactory,
        address _ton,
        address _wton,
        address _dao,
        address _depositManager,
        address _seigManager,
        address _swapProxy
    ) external;

    function setOperatorManagerFactory(address _operatorManagerFactory) external;
    function setMinimumInitialDepositAmount(uint256 _minimumInitialDepositAmount) external;

    // L1BridgeRegistry functions
    function pauseCandidateAddOn(address rollupConfig) external;
    function unpauseCandidateAddOn(address rollupConfig) external;

    // SeigManager functions
    function transferL2Seigniorage(address layer2, uint256 amount) external;

    // Public functions
    function registerCandidateAddOn(
        address rollupConfig,
        uint256 amount,
        bool flagTon,
        string calldata memo
    ) external;

    function onApprove(
        address owner,
        address spender,
        uint256 amount,
        bytes calldata data
    ) external returns (bool);

    // View functions
    function rollupConfigOfOperator(address _oper) external view returns (address);
    function operatorOfRollupConfig(address _rollupConfig) external view returns (address);
    function candidateAddOnOfOperator(address _oper) external view returns (address);
    function statusLayer2(address _rollupConfig) external view returns (uint8);
    function checkLayer2TVL(address _rollupConfig) external view returns (bool result, uint256 amount);
    function checkL1Bridge(address _rollupConfig) external view returns (bool result, address l1Bridge, address portal, address l2Ton);
    function availableRegister(address _rollupConfig) external view returns (bool result);
    function verifyOperator(address layer2, address _rollupConfig, address _operator) external view returns (bool verified);
    function checkL1BridgeDetail(address _rollupConfig) external view returns (
        bool result,
        address l1Bridge,
        address portal,
        address l2Ton,
        uint8 _type,
        uint8 status,
        bool rejectedSeigs,
        bool rejectedL2Deposit
    );
    function layerInfo(address layer2) external view returns (address rollupConfig, address operator);
}
