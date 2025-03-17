// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { IAnchorStateRegistry } from "interfaces/dispute/IAnchorStateRegistry.sol";
import { IFaultDisputeGame } from "interfaces/dispute/IFaultDisputeGame.sol";
import { ISuperchainConfig } from "interfaces/L1/ISuperchainConfig.sol";
import { IOptimismPortal2 } from "interfaces/L1/IOptimismPortal2.sol";
import { IDisputeGame } from "interfaces/dispute/IDisputeGame.sol";
import { ISemver } from "interfaces/universal/ISemver.sol";
import { GameType, Timestamp } from "src/dispute/lib/Types.sol";
import { GnosisSafe as Safe } from "safe-contracts/GnosisSafe.sol";

interface IDeputyGuardianModule is ISemver {
    error ExecutionFailed(string);
    error Unauthorized();

    event Paused(string identifier);
    event Unpaused();
    event DisputeGameBlacklisted(IDisputeGame indexed game);
    event RespectedGameTypeSet(GameType indexed gameType, Timestamp indexed updatedAt);

    function version() external view returns (string memory);
    function __constructor__(Safe _safe, ISuperchainConfig _superchainConfig, address _deputyGuardian) external;
    function safe() external view returns (Safe safe_);
    function superchainConfig() external view returns (ISuperchainConfig superchainConfig_);
    function deputyGuardian() external view returns (address deputyGuardian_);
    function pause() external;
    function unpause() external;
    function setAnchorState(IAnchorStateRegistry _registry, IFaultDisputeGame _game) external;
    function blacklistDisputeGame(IOptimismPortal2 _portal, IDisputeGame _game) external;
    function setRespectedGameType(IOptimismPortal2 _portal, GameType _gameType) external;
}
