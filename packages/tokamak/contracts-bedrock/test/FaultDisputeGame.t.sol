// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Test } from "forge-std/Test.sol";
import { Vm } from "forge-std/Vm.sol";
import { DisputeGameFactory_Init } from "test/DisputeGameFactory.t.sol";
import { DisputeGameFactory } from "src/dispute/DisputeGameFactory.sol";
import { FaultDisputeGame } from "src/dispute/FaultDisputeGame.sol";
import { PreimageOracle } from "src/cannon/PreimageOracle.sol";

import "src/libraries/DisputeTypes.sol";
import "src/libraries/DisputeErrors.sol";
import { LibClock } from "src/dispute/lib/LibUDT.sol";
import { LibPosition } from "src/dispute/lib/LibPosition.sol";
import { IPreimageOracle } from "src/dispute/interfaces/IBigStepper.sol";
import { AlphabetVM } from "test/mocks/AlphabetVM.sol";

import { DisputeActor, HonestDisputeActor } from "test/actors/FaultDisputeActors.sol";

contract FaultDisputeGame_Init is DisputeGameFactory_Init {
    /// @dev The type of the game being tested.
    GameType internal constant GAME_TYPE = GameType.wrap(0);

    /// @dev The implementation of the game.
    FaultDisputeGame internal gameImpl;
    /// @dev The `Clone` proxy of the game.
    FaultDisputeGame internal gameProxy;

    /// @dev The extra data passed to the game for initialization.
    bytes internal extraData;

    event Move(uint256 indexed parentIndex, Claim indexed pivot, address indexed claimant);

    function init(
        Claim rootClaim,
        Claim absolutePrestate,
        uint256 l2BlockNumber,
        uint256 genesisBlockNumber,
        Hash genesisOutputRoot
    )
        public
    {
        // Set the time to a realistic date.
        vm.warp(1690906994);

        // Set the extra data for the game creation
        extraData = abi.encode(l2BlockNumber);

        AlphabetVM _vm = new AlphabetVM(absolutePrestate, new PreimageOracle(0, 0, 0));

        // Deploy an implementation of the fault game
        gameImpl = new FaultDisputeGame({
            _gameType: GAME_TYPE,
            _absolutePrestate: absolutePrestate,
            _genesisBlockNumber: genesisBlockNumber,
            _genesisOutputRoot: genesisOutputRoot,
            _maxGameDepth: 2 ** 3,
            _splitDepth: 2 ** 2,
            _gameDuration: Duration.wrap(7 days),
            _vm: _vm
        });
        // Register the game implementation with the factory.
        disputeGameFactory.setImplementation(GAME_TYPE, gameImpl);
        // Create a new game.
        gameProxy = FaultDisputeGame(address(disputeGameFactory.create(GAME_TYPE, rootClaim, extraData)));

        // Check immutables
        assertEq(gameProxy.gameType().raw(), GAME_TYPE.raw());
        assertEq(gameProxy.absolutePrestate().raw(), absolutePrestate.raw());
        assertEq(gameProxy.genesisBlockNumber(), genesisBlockNumber);
        assertEq(gameProxy.genesisOutputRoot().raw(), genesisOutputRoot.raw());
        assertEq(gameProxy.maxGameDepth(), 2 ** 3);
        assertEq(gameProxy.splitDepth(), 2 ** 2);
        assertEq(gameProxy.gameDuration().raw(), 7 days);
        assertEq(address(gameProxy.vm()), address(_vm));

        // Label the proxy
        vm.label(address(gameProxy), "FaultDisputeGame_Clone");
    }

    fallback() external payable { }

    receive() external payable { }
}

contract FaultDisputeGame_Test is FaultDisputeGame_Init {
    /// @dev The root claim of the game.
    Claim internal constant ROOT_CLAIM = Claim.wrap(bytes32((uint256(1) << 248) | uint256(10)));

    /// @dev The preimage of the absolute prestate claim
    bytes internal absolutePrestateData;
    /// @dev The absolute prestate of the trace.
    Claim internal absolutePrestate;

    /// @dev Minimum bond value that covers all possible moves.
    uint256 internal constant MIN_BOND = 0.01 ether;

    function setUp() public override {
        absolutePrestateData = abi.encode(0);
        absolutePrestate = _changeClaimStatus(Claim.wrap(keccak256(absolutePrestateData)), VMStatuses.UNFINISHED);

        super.setUp();
        super.init({
            rootClaim: ROOT_CLAIM,
            absolutePrestate: absolutePrestate,
            l2BlockNumber: 0x10,
            genesisBlockNumber: 0,
            genesisOutputRoot: Hash.wrap(bytes32(0))
        });
    }

    ////////////////////////////////////////////////////////////////
    //            `IDisputeGame` Implementation Tests             //
    ////////////////////////////////////////////////////////////////

    /// @dev Tests that the constructor of the `FaultDisputeGame` reverts when the `_splitDepth`
    ///      parameter is greater than or equal to the `MAX_GAME_DEPTH`
    function test_constructor_wrongArgs_reverts(uint256 _splitDepth) public {
        AlphabetVM alphabetVM = new AlphabetVM(absolutePrestate, new PreimageOracle(0, 0, 0));

        // Test that the constructor reverts when the `_splitDepth` parameter is greater than or equal
        // to the `MAX_GAME_DEPTH` parameter.
        _splitDepth = bound(_splitDepth, 2 ** 3, type(uint256).max);
        vm.expectRevert(InvalidSplitDepth.selector);
        new FaultDisputeGame({
            _gameType: GAME_TYPE,
            _absolutePrestate: absolutePrestate,
            _genesisBlockNumber: 0,
            _genesisOutputRoot: Hash.wrap(bytes32(0)),
            _maxGameDepth: 2 ** 3,
            _splitDepth: _splitDepth,
            _gameDuration: Duration.wrap(7 days),
            _vm: alphabetVM
        });
    }

    /// @dev Tests that the game's root claim is set correctly.
    function test_rootClaim_succeeds() public {
        assertEq(gameProxy.rootClaim().raw(), ROOT_CLAIM.raw());
    }

    /// @dev Tests that the game's extra data is set correctly.
    function test_extraData_succeeds() public {
        assertEq(gameProxy.extraData(), extraData);
    }

    /// @dev Tests that the game's starting timestamp is set correctly.
    function test_createdAt_succeeds() public {
        assertEq(gameProxy.createdAt().raw(), block.timestamp);
    }

    /// @dev Tests that the game's type is set correctly.
    function test_gameType_succeeds() public {
        assertEq(gameProxy.gameType().raw(), GAME_TYPE.raw());
    }

    /// @dev Tests that the game's data is set correctly.
    function test_gameData_succeeds() public {
        (GameType gameType, Claim rootClaim, bytes memory _extraData) = gameProxy.gameData();

        assertEq(gameType.raw(), GAME_TYPE.raw());
        assertEq(rootClaim.raw(), ROOT_CLAIM.raw());
        assertEq(_extraData, extraData);
    }

    ////////////////////////////////////////////////////////////////
    //          `IFaultDisputeGame` Implementation Tests       //
    ////////////////////////////////////////////////////////////////

    /// @dev Tests that the game cannot be initialized with an output root that commits to <= the configured genesis
    ///      block number
    function testFuzz_initialize_cannotProposeGenesis_reverts(uint256 _blockNumber) public {
        _blockNumber = bound(_blockNumber, 0, gameProxy.genesisBlockNumber());

        Claim claim = _dummyClaim();
        vm.expectRevert(abi.encodeWithSelector(UnexpectedRootClaim.selector, claim));
        gameProxy = FaultDisputeGame(address(disputeGameFactory.create(GAME_TYPE, claim, abi.encode(_blockNumber))));
    }

    /// @dev Tests that the proxy receives ETH from the dispute game factory.
    function test_initialize_receivesETH_succeeds(uint128 _value) public {
        _value = uint128(bound(_value, gameProxy.getRequiredBond(Position.wrap(1)), type(uint128).max));
        vm.deal(address(this), _value);

        assertEq(address(gameProxy).balance, 0);
        gameProxy =
            FaultDisputeGame(address(disputeGameFactory.create{ value: _value }(GAME_TYPE, ROOT_CLAIM, abi.encode(1))));
        assertEq(address(gameProxy).balance, _value);
    }

    /// @dev Tests that the game cannot be initialized with extra data > 64 bytes long (root claim + l2 block number
    ///      concatenated)
    function testFuzz_initialize_extraDataTooLong_reverts(uint256 _extraDataLen) public {
        // The `DisputeGameFactory` will pack the root claim and the extra data into a single array, which is enforced
        // to be at least 64 bytes long.
        // We bound the upper end to 23.5KB to ensure that the minimal proxy never surpasses the contract size limit
        // in this test, as CWIA proxies store the immutable args in their bytecode.
        // [33 bytes, 23.5 KB]
        _extraDataLen = bound(_extraDataLen, 33, 23_500);
        bytes memory _extraData = new bytes(_extraDataLen);

        // Assign the first 32 bytes in `extraData` to a valid L2 block number passed genesis.
        uint256 genesisBlockNumber = gameProxy.genesisBlockNumber();
        assembly {
            mstore(add(_extraData, 0x20), add(genesisBlockNumber, 1))
        }

        Claim claim = _dummyClaim();
        vm.expectRevert(abi.encodeWithSelector(ExtraDataTooLong.selector));
        gameProxy = FaultDisputeGame(address(disputeGameFactory.create(GAME_TYPE, claim, _extraData)));
    }

    /// @dev Tests that the game is initialized with the correct data.
    function test_initialize_correctData_succeeds() public {
        // Assert that the root claim is initialized correctly.
        (
            uint32 parentIndex,
            address counteredBy,
            address claimant,
            uint128 bond,
            Claim claim,
            Position position,
            Clock clock
        ) = gameProxy.claimData(0);
        assertEq(parentIndex, type(uint32).max);
        assertEq(counteredBy, address(0));
        assertEq(claimant, DEFAULT_SENDER);
        assertEq(bond, 0);
        assertEq(claim.raw(), ROOT_CLAIM.raw());
        assertEq(position.raw(), 1);
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(0), Timestamp.wrap(uint64(block.timestamp))).raw());

        // Assert that the `createdAt` timestamp is correct.
        assertEq(gameProxy.createdAt().raw(), block.timestamp);

        // Assert that the blockhash provided is correct.
        assertEq(gameProxy.l1Head().raw(), blockhash(block.number - 1));
    }

    /// @dev Tests that the game cannot be initialized twice.
    function test_initialize_onlyOnce_succeeds() public {
        vm.expectRevert(AlreadyInitialized.selector);
        gameProxy.initialize();
    }

    /// @dev Tests that a move while the game status is not `IN_PROGRESS` causes the call to revert
    ///      with the `GameNotInProgress` error
    function test_move_gameNotInProgress_reverts() public {
        uint256 chalWins = uint256(GameStatus.CHALLENGER_WINS);

        // Replace the game status in storage. It exists in slot 0 at offset 16.
        uint256 slot = uint256(vm.load(address(gameProxy), bytes32(0)));
        uint256 offset = 16 << 3;
        uint256 mask = 0xFF << offset;
        // Replace the byte in the slot value with the challenger wins status.
        slot = (slot & ~mask) | (chalWins << offset);
        vm.store(address(gameProxy), bytes32(0), bytes32(slot));

        // Ensure that the game status was properly updated.
        GameStatus status = gameProxy.status();
        assertEq(uint256(status), chalWins);

        // Attempt to make a move. Should revert.
        vm.expectRevert(GameNotInProgress.selector);
        gameProxy.attack(0, Claim.wrap(0));
    }

    /// @dev Tests that an attempt to defend the root claim reverts with the `CannotDefendRootClaim` error.
    function test_move_defendRoot_reverts() public {
        vm.expectRevert(CannotDefendRootClaim.selector);
        gameProxy.defend(0, _dummyClaim());
    }

    /// @dev Tests that an attempt to move against a claim that does not exist reverts with the
    ///      `ParentDoesNotExist` error.
    function test_move_nonExistentParent_reverts() public {
        Claim claim = _dummyClaim();

        // Expect an out of bounds revert for an attack
        vm.expectRevert(abi.encodeWithSignature("Panic(uint256)", 0x32));
        gameProxy.attack(1, claim);

        // Expect an out of bounds revert for a defense
        vm.expectRevert(abi.encodeWithSignature("Panic(uint256)", 0x32));
        gameProxy.defend(1, claim);
    }

    /// @dev Tests that an attempt to move at the maximum game depth reverts with the
    ///      `GameDepthExceeded` error.
    function test_move_gameDepthExceeded_reverts() public {
        Claim claim = _changeClaimStatus(_dummyClaim(), VMStatuses.PANIC);

        uint256 maxDepth = gameProxy.maxGameDepth();

        for (uint256 i = 0; i <= maxDepth; i++) {
            // At the max game depth, the `_move` function should revert with
            // the `GameDepthExceeded` error.
            if (i == maxDepth) {
                vm.expectRevert(GameDepthExceeded.selector);
            }
            gameProxy.attack{ value: MIN_BOND }(i, claim);
        }
    }

    /// @dev Tests that a move made after the clock time has exceeded reverts with the
    ///      `ClockTimeExceeded` error.
    function test_move_clockTimeExceeded_reverts() public {
        // Warp ahead past the clock time for the first move (3 1/2 days)
        vm.warp(block.timestamp + 3 days + 12 hours + 1);
        vm.expectRevert(ClockTimeExceeded.selector);
        gameProxy.attack{ value: MIN_BOND }(0, _dummyClaim());
    }

    /// @notice Static unit test for the correctness of the chess clock incrementation.
    function test_move_clockCorrectness_succeeds() public {
        (,,,,,, Clock clock) = gameProxy.claimData(0);
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(0), Timestamp.wrap(uint64(block.timestamp))).raw());

        Claim claim = _dummyClaim();

        vm.warp(block.timestamp + 15);
        gameProxy.attack{ value: MIN_BOND }(0, claim);
        (,,,,,, clock) = gameProxy.claimData(1);
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(15), Timestamp.wrap(uint64(block.timestamp))).raw());

        vm.warp(block.timestamp + 10);
        gameProxy.attack{ value: MIN_BOND }(1, claim);
        (,,,,,, clock) = gameProxy.claimData(2);
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(10), Timestamp.wrap(uint64(block.timestamp))).raw());

        // We are at the split depth, so we need to set the status byte of the claim
        // for the next move.
        claim = _changeClaimStatus(claim, VMStatuses.PANIC);

        vm.warp(block.timestamp + 10);
        gameProxy.attack{ value: MIN_BOND }(2, claim);
        (,,,,,, clock) = gameProxy.claimData(3);
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(25), Timestamp.wrap(uint64(block.timestamp))).raw());

        vm.warp(block.timestamp + 10);
        gameProxy.attack{ value: MIN_BOND }(3, claim);
        (,,,,,, clock) = gameProxy.claimData(4);
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(20), Timestamp.wrap(uint64(block.timestamp))).raw());
    }

    /// @dev Tests that an identical claim cannot be made twice. The duplicate claim attempt should
    ///      revert with the `ClaimAlreadyExists` error.
    function test_move_duplicateClaim_reverts() public {
        Claim claim = _dummyClaim();

        // Make the first move. This should succeed.
        gameProxy.attack{ value: MIN_BOND }(0, claim);

        // Attempt to make the same move again.
        vm.expectRevert(ClaimAlreadyExists.selector);
        gameProxy.attack{ value: MIN_BOND }(0, claim);
    }

    /// @dev Static unit test asserting that identical claims at the same position can be made in different subgames.
    function test_move_duplicateClaimsDifferentSubgames_succeeds() public {
        Claim claimA = _dummyClaim();
        Claim claimB = _dummyClaim();

        // Make the first moves. This should succeed.
        gameProxy.attack{ value: MIN_BOND }(0, claimA);
        gameProxy.attack{ value: MIN_BOND }(0, claimB);

        // Perform an attack at the same position with the same claim value in both subgames.
        // These both should succeed.
        gameProxy.attack{ value: MIN_BOND }(1, claimA);
        gameProxy.attack{ value: MIN_BOND }(2, claimA);
    }

    /// @dev Static unit test for the correctness of an opening attack.
    function test_move_simpleAttack_succeeds() public {
        // Warp ahead 5 seconds.
        vm.warp(block.timestamp + 5);

        Claim counter = _dummyClaim();

        // Perform the attack.
        vm.expectEmit(true, true, true, false);
        emit Move(0, counter, address(this));
        gameProxy.attack{ value: MIN_BOND }(0, counter);

        // Grab the claim data of the attack.
        (
            uint32 parentIndex,
            address counteredBy,
            address claimant,
            uint128 bond,
            Claim claim,
            Position position,
            Clock clock
        ) = gameProxy.claimData(1);

        // Assert correctness of the attack claim's data.
        assertEq(parentIndex, 0);
        assertEq(counteredBy, address(0));
        assertEq(claimant, address(this));
        assertEq(bond, MIN_BOND);
        assertEq(claim.raw(), counter.raw());
        assertEq(position.raw(), Position.wrap(1).move(true).raw());
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(5), Timestamp.wrap(uint64(block.timestamp))).raw());

        // Grab the claim data of the parent.
        (parentIndex, counteredBy, claimant, bond, claim, position, clock) = gameProxy.claimData(0);

        // Assert correctness of the parent claim's data.
        assertEq(parentIndex, type(uint32).max);
        assertEq(counteredBy, address(0));
        assertEq(claimant, DEFAULT_SENDER);
        assertEq(bond, 0);
        assertEq(claim.raw(), ROOT_CLAIM.raw());
        assertEq(position.raw(), 1);
        assertEq(clock.raw(), LibClock.wrap(Duration.wrap(0), Timestamp.wrap(uint64(block.timestamp - 5))).raw());
    }

    /// @dev Tests that making a claim at the execution trace bisection root level with an invalid status
    ///      byte reverts with the `UnexpectedRootClaim` error.
    function test_move_incorrectStatusExecRoot_reverts() public {
        for (uint256 i; i < 4; i++) {
            gameProxy.attack{ value: MIN_BOND }(i, _dummyClaim());
        }

        vm.expectRevert(abi.encodeWithSelector(UnexpectedRootClaim.selector, bytes32(0)));
        gameProxy.attack{ value: MIN_BOND }(4, Claim.wrap(bytes32(0)));
    }

    /// @dev Tests that making a claim at the execution trace bisection root level with a valid status
    ///      byte succeeds.
    function test_move_correctStatusExecRoot_succeeds() public {
        for (uint256 i; i < 4; i++) {
            gameProxy.attack{ value: MIN_BOND }(i, _dummyClaim());
        }
        gameProxy.attack{ value: MIN_BOND }(4, _changeClaimStatus(_dummyClaim(), VMStatuses.PANIC));
    }

    /// @dev Static unit test asserting that a move reverts when the bond is insufficient.
    function test_move_insufficientBond_reverts() public {
        vm.expectRevert(InsufficientBond.selector);
        gameProxy.attack{ value: 0 }(0, _dummyClaim());
    }

    /// @dev Static unit test for the correctness an uncontested root resolution.
    function test_resolve_rootUncontested_succeeds() public {
        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);
        gameProxy.resolveClaim(0);
        assertEq(uint8(gameProxy.resolve()), uint8(GameStatus.DEFENDER_WINS));
    }

    /// @dev Static unit test for the correctness an uncontested root resolution.
    function test_resolve_rootUncontestedClockNotExpired_succeeds() public {
        vm.warp(block.timestamp + 3 days + 12 hours);
        vm.expectRevert(ClockNotExpired.selector);
        gameProxy.resolveClaim(0);
    }

    /// @dev Static unit test asserting that resolve reverts when the absolute root
    ///      subgame has not been resolved.
    function test_resolve_rootUncontestedButUnresolved_reverts() public {
        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);
        vm.expectRevert(OutOfOrderResolution.selector);
        gameProxy.resolve();
    }

    /// @dev Static unit test asserting that resolve reverts when the game state is
    ///      not in progress.
    function test_resolve_notInProgress_reverts() public {
        uint256 chalWins = uint256(GameStatus.CHALLENGER_WINS);

        // Replace the game status in storage. It exists in slot 0 at offset 16.
        uint256 slot = uint256(vm.load(address(gameProxy), bytes32(0)));
        uint256 offset = 16 << 3;
        uint256 mask = 0xFF << offset;
        // Replace the byte in the slot value with the challenger wins status.
        slot = (slot & ~mask) | (chalWins << offset);

        vm.store(address(gameProxy), bytes32(uint256(0)), bytes32(slot));
        vm.expectRevert(GameNotInProgress.selector);
        gameProxy.resolveClaim(0);
    }

    /// @dev Static unit test for the correctness of resolving a single attack game state.
    function test_resolve_rootContested_succeeds() public {
        gameProxy.attack{ value: MIN_BOND }(0, _dummyClaim());

        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        gameProxy.resolveClaim(0);
        assertEq(uint8(gameProxy.resolve()), uint8(GameStatus.CHALLENGER_WINS));
    }

    /// @dev Static unit test for the correctness of resolving a game with a contested challenge claim.
    function test_resolve_challengeContested_succeeds() public {
        gameProxy.attack{ value: MIN_BOND }(0, _dummyClaim());
        gameProxy.defend{ value: MIN_BOND }(1, _dummyClaim());

        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        gameProxy.resolveClaim(1);
        gameProxy.resolveClaim(0);
        assertEq(uint8(gameProxy.resolve()), uint8(GameStatus.DEFENDER_WINS));
    }

    /// @dev Static unit test for the correctness of resolving a game with multiplayer moves.
    function test_resolve_teamDeathmatch_succeeds() public {
        gameProxy.attack{ value: MIN_BOND }(0, _dummyClaim());
        gameProxy.attack{ value: MIN_BOND }(0, _dummyClaim());
        gameProxy.defend{ value: MIN_BOND }(1, _dummyClaim());
        gameProxy.defend{ value: MIN_BOND }(1, _dummyClaim());

        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        gameProxy.resolveClaim(1);
        gameProxy.resolveClaim(0);
        assertEq(uint8(gameProxy.resolve()), uint8(GameStatus.CHALLENGER_WINS));
    }

    /// @dev Static unit test for the correctness of resolving a game that reaches max game depth.
    function test_resolve_stepReached_succeeds() public {
        Claim claim = _dummyClaim();
        for (uint256 i; i < gameProxy.splitDepth(); i++) {
            gameProxy.attack{ value: MIN_BOND }(i, claim);
        }

        claim = _changeClaimStatus(claim, VMStatuses.PANIC);
        for (uint256 i = gameProxy.claimDataLen() - 1; i < gameProxy.maxGameDepth(); i++) {
            gameProxy.attack{ value: MIN_BOND }(i, claim);
        }

        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        // resolving claim at 8 isn't necessary
        for (uint256 i = 8; i > 0; i--) {
            gameProxy.resolveClaim(i - 1);
        }
        assertEq(uint8(gameProxy.resolve()), uint8(GameStatus.DEFENDER_WINS));
    }

    /// @dev Static unit test asserting that resolve reverts when attempting to resolve a subgame multiple times
    function test_resolve_claimAlreadyResolved_reverts() public {
        vm.deal(address(this), 2 * MIN_BOND);

        Claim claim = _dummyClaim();
        gameProxy.attack{ value: MIN_BOND }(0, claim);
        gameProxy.attack{ value: MIN_BOND }(1, claim);

        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        assertEq(address(this).balance, 0);
        gameProxy.resolveClaim(1);
        gameProxy.claimCredit(address(this));
        assertEq(address(this).balance, MIN_BOND);

        vm.expectRevert(ClaimAlreadyResolved.selector);
        gameProxy.resolveClaim(1);
        assertEq(address(this).balance, MIN_BOND);
    }

    /// @dev Static unit test asserting that resolve reverts when attempting to resolve a subgame at max depth
    function test_resolve_claimAtMaxDepthAlreadyResolved_reverts() public {
        Claim claim = _dummyClaim();
        for (uint256 i; i < gameProxy.splitDepth(); i++) {
            gameProxy.attack{ value: MIN_BOND }(i, claim);
        }

        vm.deal(address(this), 100 ether);
        claim = _changeClaimStatus(claim, VMStatuses.PANIC);
        for (uint256 i = gameProxy.claimDataLen() - 1; i < gameProxy.maxGameDepth(); i++) {
            gameProxy.attack{ value: MIN_BOND }(i, claim);
        }

        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        // Resolve to claim bond
        uint256 balanceBefore = address(this).balance;
        gameProxy.resolveClaim(8);
        gameProxy.claimCredit(address(this));
        assertEq(address(this).balance, balanceBefore + MIN_BOND);

        vm.expectRevert(ClaimAlreadyResolved.selector);
        gameProxy.resolveClaim(8);
    }

    /// @dev Static unit test asserting that resolve reverts when attempting to resolve subgames out of order
    function test_resolve_outOfOrderResolution_reverts() public {
        gameProxy.attack{ value: MIN_BOND }(0, _dummyClaim());
        gameProxy.attack{ value: MIN_BOND }(1, _dummyClaim());

        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        vm.expectRevert(OutOfOrderResolution.selector);
        gameProxy.resolveClaim(0);
    }

    /// @dev Static unit test asserting that resolve pays out bonds on step, output bisection, and execution trace
    /// moves.
    function test_resolve_bondPayouts_succeeds() public {
        // Give the test contract some ether
        vm.deal(address(this), 100 ether);

        // Make claims all the way down the tree.
        gameProxy.attack{ value: 1 ether }(0, _dummyClaim());
        gameProxy.attack{ value: 1 ether }(1, _dummyClaim());
        gameProxy.attack{ value: 1 ether }(2, _dummyClaim());
        gameProxy.attack{ value: 1 ether }(3, _dummyClaim());
        gameProxy.attack{ value: 1 ether }(4, _changeClaimStatus(_dummyClaim(), VMStatuses.PANIC));
        gameProxy.attack{ value: 1 ether }(5, _dummyClaim());
        gameProxy.attack{ value: 1 ether }(6, _dummyClaim());
        gameProxy.attack{ value: 1 ether }(7, _dummyClaim());
        gameProxy.addLocalData(LocalPreimageKey.STARTING_L2_BLOCK_NUMBER, 8, 0);
        gameProxy.step(8, true, absolutePrestateData, hex"");

        // Ensure that the step successfully countered the leaf claim.
        (, address counteredBy,,,,,) = gameProxy.claimData(8);
        assertEq(counteredBy, address(this));

        // Ensure we bonded the correct amounts
        uint256 bonded = (gameProxy.claimDataLen() - 1) * 1 ether;
        assertEq(address(this).balance, 100 ether - bonded);
        assertEq(address(gameProxy).balance, bonded);

        // Resolve all claims
        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);
        for (uint256 i = gameProxy.claimDataLen(); i > 0; i--) {
            (bool success,) = address(gameProxy).call(abi.encodeCall(gameProxy.resolveClaim, (i - 1)));
            assertTrue(success);
        }
        gameProxy.resolve();

        gameProxy.claimCredit(address(this));

        // Ensure that bonds were paid out correctly.
        assertEq(address(this).balance, 100 ether);
        assertEq(address(gameProxy).balance, 0);

        // Ensure that the init bond for the game is 0, in case we change it in the test suite in the future.
        assertEq(disputeGameFactory.initBonds(GAME_TYPE), 0);
    }

    /// @dev Static unit test asserting that resolve pays out bonds on step, output bisection, and execution trace
    /// moves with 2 actors and a dishonest root claim.
    function test_resolve_bondPayoutsSeveralActors_succeeds() public {
        // Give the test contract and bob some ether
        address bob = address(0xb0b);
        vm.deal(address(this), 100 ether);
        vm.deal(bob, 100 ether);

        // Make claims all the way down the tree, trading off between bob and the test contract.
        gameProxy.attack{ value: 1 ether }(0, _dummyClaim());

        vm.prank(bob);
        gameProxy.attack{ value: 1 ether }(1, _dummyClaim());

        gameProxy.attack{ value: 1 ether }(2, _dummyClaim());

        vm.prank(bob);
        gameProxy.attack{ value: 1 ether }(3, _dummyClaim());

        gameProxy.attack{ value: 1 ether }(4, _changeClaimStatus(_dummyClaim(), VMStatuses.PANIC));

        vm.prank(bob);
        gameProxy.attack{ value: 1 ether }(5, _dummyClaim());

        gameProxy.attack{ value: 1 ether }(6, _dummyClaim());

        vm.prank(bob);
        gameProxy.attack{ value: 1 ether }(7, _dummyClaim());

        gameProxy.addLocalData(LocalPreimageKey.STARTING_L2_BLOCK_NUMBER, 8, 0);
        gameProxy.step(8, true, absolutePrestateData, hex"");

        // Ensure that the step successfully countered the leaf claim.
        (, address counteredBy,,,,,) = gameProxy.claimData(8);
        assertEq(counteredBy, address(this));

        // Ensure we bonded the correct amounts
        uint256 bonded = ((gameProxy.claimDataLen() - 1) / 2) * 1 ether;
        assertEq(address(this).balance, 100 ether - bonded);
        assertEq(bob.balance, 100 ether - bonded);
        assertEq(address(gameProxy).balance, bonded * 2);

        // Resolve all claims
        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);
        for (uint256 i = gameProxy.claimDataLen(); i > 0; i--) {
            (bool success,) = address(gameProxy).call(abi.encodeCall(gameProxy.resolveClaim, (i - 1)));
            assertTrue(success);
        }
        gameProxy.resolve();

        gameProxy.claimCredit(address(this));
        gameProxy.claimCredit(bob);

        // Ensure that bonds were paid out correctly.
        assertEq(address(this).balance, 100 ether + bonded);
        assertEq(bob.balance, 100 ether - bonded);
        assertEq(address(gameProxy).balance, 0);

        // Ensure that the init bond for the game is 0, in case we change it in the test suite in the future.
        assertEq(disputeGameFactory.initBonds(GAME_TYPE), 0);
    }

    /// @dev Static unit test asserting that resolve pays out bonds on moves to the leftmost actor
    /// in subgames containing successful counters.
    function test_resolve_leftmostBondPayout_succeeds() public {
        address alice = address(0xa11ce);
        address bob = address(0xb0b);
        address charlie = address(0xc0c);
        vm.deal(address(this), 100 ether);
        vm.deal(alice, 100 ether);
        vm.deal(bob, 100 ether);
        vm.deal(charlie, 100 ether);

        // Make claims with bob, charlie and the test contract on defense, and alice as the challenger
        // charlie is successfully countered by alice
        // alice is successfully countered by both bob and the test contract
        vm.prank(alice);
        gameProxy.attack{ value: 1 ether }(0, _dummyClaim());

        vm.prank(bob);
        gameProxy.defend{ value: 1 ether }(1, _dummyClaim());
        vm.prank(charlie);
        gameProxy.attack{ value: 1 ether }(1, _dummyClaim());
        gameProxy.attack{ value: 1 ether }(1, _dummyClaim());

        vm.prank(alice);
        gameProxy.attack{ value: 1 ether }(3, _dummyClaim());

        // Resolve all claims
        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);
        for (uint256 i = gameProxy.claimDataLen(); i > 0; i--) {
            (bool success,) = address(gameProxy).call(abi.encodeCall(gameProxy.resolveClaim, (i - 1)));
            assertTrue(success);
        }
        gameProxy.resolve();

        gameProxy.claimCredit(address(this));
        gameProxy.claimCredit(alice);
        gameProxy.claimCredit(bob);
        gameProxy.claimCredit(charlie);

        // Ensure that bonds were paid out correctly.
        uint256 aliceLosses = 1 ether;
        uint256 charlieLosses = 1 ether;
        assertEq(address(this).balance, 100 ether + aliceLosses, "incorrect this balance");
        assertEq(alice.balance, 100 ether - aliceLosses + charlieLosses, "incorrect alice balance");
        assertEq(bob.balance, 100 ether, "incorrect bob balance");
        assertEq(charlie.balance, 100 ether - charlieLosses, "incorrect charlie balance");
        assertEq(address(gameProxy).balance, 0);

        // Ensure that the init bond for the game is 0, in case we change it in the test suite in the future.
        assertEq(disputeGameFactory.initBonds(GAME_TYPE), 0);
    }

    /// @dev Static unit test asserting that credit may not be drained past allowance through reentrancy.
    function test_claimCredit_claimAlreadyResolved_reverts() public {
        ClaimCreditReenter reenter = new ClaimCreditReenter(gameProxy);
        vm.startPrank(address(reenter));

        // Give the test contract some ether to bond.
        vm.deal(address(reenter), MIN_BOND * 2);
        // Give the game proxy 1 extra ether, unregistered.
        vm.deal(address(gameProxy), 1 ether);

        // Perform a bonded move.
        Claim claim = _dummyClaim();
        gameProxy.attack{ value: MIN_BOND }(0, claim);
        gameProxy.attack{ value: MIN_BOND }(1, claim);

        // Warp past the finalization period
        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        // Ensure that we bonded all the test contract's ETH
        assertEq(address(reenter).balance, 0);
        // Ensure the game proxy has 1 + 2 * MIN_BOND, 2*MIN_BOND from bonding and 1 unregistered.
        assertEq(address(gameProxy).balance, 1 ether + MIN_BOND * 2);

        // Resolve the claim at gindex 1 and claim the reenter contract's credit.
        gameProxy.resolveClaim(1);

        // Ensure that the game registered the `reenter` contract's credit.
        assertEq(gameProxy.credit(address(reenter)), MIN_BOND);

        // Initiate the reentrant credit claim.
        reenter.claimCredit(address(reenter));

        // The reenter contract should have performed 5 calls to `claimCredit`, but only received the amount bonded for
        // the gindex 1 subgame.
        // The root claim bond and the unregistered ETH should still exist in the game proxy.
        assertEq(reenter.numCalls(), 5);
        assertEq(address(reenter).balance, MIN_BOND);
        assertEq(address(gameProxy).balance, 1 ether + MIN_BOND);

        vm.stopPrank();
    }

    /// @dev Tests that adding local data with an out of bounds identifier reverts.
    function testFuzz_addLocalData_oob_reverts(uint256 _ident) public {
        // Get a claim below the split depth so that we can add local data for an execution trace subgame.
        for (uint256 i; i < 4; i++) {
            gameProxy.attack{ value: MIN_BOND }(i, _dummyClaim());
        }
        gameProxy.attack{ value: MIN_BOND }(4, _changeClaimStatus(_dummyClaim(), VMStatuses.PANIC));

        // [1, 5] are valid local data identifiers.
        if (_ident <= 5) _ident = 0;

        vm.expectRevert(InvalidLocalIdent.selector);
        gameProxy.addLocalData(_ident, 5, 0);
    }

    /// @dev Tests that local data is loaded into the preimage oracle correctly in the subgame
    ///      that is disputing the transition from `GENESIS -> GENESIS + 1`
    function test_addLocalDataGenesisTransition_static_succeeds() public {
        IPreimageOracle oracle = IPreimageOracle(address(gameProxy.vm().oracle()));

        // Get a claim below the split depth so that we can add local data for an execution trace subgame.
        for (uint256 i; i < 4; i++) {
            gameProxy.attack{ value: MIN_BOND }(i, Claim.wrap(bytes32(i)));
        }
        gameProxy.attack{ value: MIN_BOND }(4, _changeClaimStatus(_dummyClaim(), VMStatuses.PANIC));

        // Expected start/disputed claims
        bytes32 startingClaim = gameProxy.genesisOutputRoot().raw();
        bytes32 disputedClaim = bytes32(uint256(3));
        Position disputedPos = LibPosition.wrap(4, 0);

        // Expected local data
        bytes32[5] memory data =
            [gameProxy.l1Head().raw(), startingClaim, disputedClaim, bytes32(0), bytes32(block.chainid << 0xC0)];

        for (uint256 i = 1; i <= 5; i++) {
            uint256 expectedLen = i > 3 ? 8 : 32;
            bytes32 key = _getKey(i, keccak256(abi.encode(disputedClaim, disputedPos)));

            gameProxy.addLocalData(i, 5, 0);
            (bytes32 dat, uint256 datLen) = oracle.readPreimage(key, 0);
            assertEq(dat >> 0xC0, bytes32(expectedLen));
            // Account for the length prefix if i > 3 (the data stored
            // at identifiers i <= 3 are 32 bytes long, so the expected
            // length is already correct. If i > 3, the data is only 8
            // bytes long, so the length prefix + the data is 16 bytes
            // total.)
            assertEq(datLen, expectedLen + (i > 3 ? 8 : 0));

            gameProxy.addLocalData(i, 5, 8);
            (dat, datLen) = oracle.readPreimage(key, 8);
            assertEq(dat, data[i - 1]);
            assertEq(datLen, expectedLen);
        }
    }

    /// @dev Tests that local data is loaded into the preimage oracle correctly.
    function test_addLocalDataMiddle_static_succeeds() public {
        IPreimageOracle oracle = IPreimageOracle(address(gameProxy.vm().oracle()));

        // Get a claim below the split depth so that we can add local data for an execution trace subgame.
        for (uint256 i; i < 4; i++) {
            gameProxy.attack{ value: MIN_BOND }(i, Claim.wrap(bytes32(i)));
        }
        gameProxy.defend{ value: MIN_BOND }(4, _changeClaimStatus(ROOT_CLAIM, VMStatuses.VALID));

        // Expected start/disputed claims
        bytes32 startingClaim = bytes32(uint256(3));
        Position startingPos = LibPosition.wrap(4, 0);
        bytes32 disputedClaim = bytes32(uint256(2));
        Position disputedPos = LibPosition.wrap(3, 0);

        // Expected local data
        bytes32[5] memory data = [
            gameProxy.l1Head().raw(),
            startingClaim,
            disputedClaim,
            bytes32(uint256(1) << 0xC0),
            bytes32(block.chainid << 0xC0)
        ];

        for (uint256 i = 1; i <= 5; i++) {
            uint256 expectedLen = i > 3 ? 8 : 32;
            bytes32 key = _getKey(i, keccak256(abi.encode(startingClaim, startingPos, disputedClaim, disputedPos)));

            gameProxy.addLocalData(i, 5, 0);
            (bytes32 dat, uint256 datLen) = oracle.readPreimage(key, 0);
            assertEq(dat >> 0xC0, bytes32(expectedLen));
            // Account for the length prefix if i > 3 (the data stored
            // at identifiers i <= 3 are 32 bytes long, so the expected
            // length is already correct. If i > 3, the data is only 8
            // bytes long, so the length prefix + the data is 16 bytes
            // total.)
            assertEq(datLen, expectedLen + (i > 3 ? 8 : 0));

            gameProxy.addLocalData(i, 5, 8);
            (dat, datLen) = oracle.readPreimage(key, 8);
            assertEq(dat, data[i - 1]);
            assertEq(datLen, expectedLen);
        }
    }

    /// @dev Helper to return a pseudo-random claim
    function _dummyClaim() internal view returns (Claim) {
        return Claim.wrap(keccak256(abi.encode(gasleft())));
    }

    /// @dev Helper to get the localized key for an identifier in the context of the game proxy.
    function _getKey(uint256 _ident, bytes32 _localContext) internal view returns (bytes32) {
        bytes32 h = keccak256(abi.encode(_ident | (1 << 248), address(gameProxy), _localContext));
        return bytes32((uint256(h) & ~uint256(0xFF << 248)) | (1 << 248));
    }
}

contract FaultDispute_1v1_Actors_Test is FaultDisputeGame_Init {
    /// @dev The honest actor
    DisputeActor internal honest;
    /// @dev The dishonest actor
    DisputeActor internal dishonest;

    function setUp() public override {
        // Setup the `FaultDisputeGame`
        super.setUp();
    }

    /// @notice Fuzz test for a 1v1 output bisection dispute.
    /// @dev The alphabet game has a constant status byte, and is not safe from someone being dishonest in
    ///      output bisection and then posting a correct execution trace bisection root claim. This test
    ///      does not cover this case (i.e. root claim of output bisection is dishonest, root claim of
    ///      execution trace bisection is made by the dishonest actor but is honest, honest actor cannot
    ///      attack it without risk of losing).
    function testFuzz_outputBisection1v1honestRoot_succeeds(uint8 _divergeOutput, uint8 _divergeStep) public {
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        uint256 divergeAtOutput = bound(_divergeOutput, 0, 15);
        uint256 divergeAtStep = bound(_divergeStep, 0, 7);
        uint256 divergeStepOffset = (divergeAtOutput << 4) + divergeAtStep;

        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i >= divergeAtOutput ? 0xFF : i + 1;
        }
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = i >= divergeStepOffset ? bytes1(uint8(0xFF)) : bytes1(uint8(i));
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 16,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.DEFENDER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1honestRootGenesisAbsolutePrestate_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are from [2, 17] in this game.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i + 2;
        }
        // The dishonest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of all set bits.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = bytes1(0xFF);
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 16,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.DEFENDER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1dishonestRootGenesisAbsolutePrestate_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are from [2, 17] in this game.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i + 2;
        }
        // The dishonest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of all set bits.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = bytes1(0xFF);
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 17,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.CHALLENGER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1honestRoot_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are from [2, 17] in this game.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i + 2;
        }
        // The dishonest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of all zeros.
        bytes memory dishonestTrace = new bytes(256);

        // Run the actor test
        _actorTest({
            _rootClaim: 16,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.DEFENDER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1dishonestRoot_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are from [2, 17] in this game.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i + 2;
        }
        // The dishonest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of all zeros.
        bytes memory dishonestTrace = new bytes(256);

        // Run the actor test
        _actorTest({
            _rootClaim: 17,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.CHALLENGER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1correctRootHalfWay_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are half correct, half incorrect.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i > 7 ? 0xFF : i + 1;
        }
        // The dishonest trace is half correct, half incorrect.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = i > (127 + 4) ? bytes1(0xFF) : bytes1(uint8(i));
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 16,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.DEFENDER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1dishonestRootHalfWay_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are half correct, half incorrect.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i > 7 ? 0xFF : i + 1;
        }
        // The dishonest trace is half correct, half incorrect.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = i > (127 + 4) ? bytes1(0xFF) : bytes1(uint8(i));
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 0xFF,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.CHALLENGER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1correctAbsolutePrestate_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are half correct, half incorrect.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i > 7 ? 0xFF : i + 1;
        }
        // The dishonest trace correct is half correct, half incorrect.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = i > 127 ? bytes1(0xFF) : bytes1(uint8(i));
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 16,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.DEFENDER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1dishonestAbsolutePrestate_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are half correct, half incorrect.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i > 7 ? 0xFF : i + 1;
        }
        // The dishonest trace correct is half correct, half incorrect.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = i > 127 ? bytes1(0xFF) : bytes1(uint8(i));
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 0xFF,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.CHALLENGER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1honestRootFinalInstruction_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are half correct, half incorrect.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i > 7 ? 0xFF : i + 1;
        }
        // The dishonest trace is half correct, and correct all the way up to the final instruction of the exec
        // subgame.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = i > (127 + 7) ? bytes1(0xFF) : bytes1(uint8(i));
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 16,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.DEFENDER_WINS
        });
    }

    /// @notice Static unit test for a 1v1 output bisection dispute.
    function test_static_1v1dishonestRootFinalInstruction_succeeds() public {
        // The honest l2 outputs are from [1, 16] in this game.
        uint256[] memory honestL2Outputs = new uint256[](16);
        for (uint256 i; i < honestL2Outputs.length; i++) {
            honestL2Outputs[i] = i + 1;
        }
        // The honest trace covers all block -> block + 1 transitions, and is 256 bytes long, consisting
        // of bytes [0, 255].
        bytes memory honestTrace = new bytes(256);
        for (uint256 i; i < honestTrace.length; i++) {
            honestTrace[i] = bytes1(uint8(i));
        }

        // The dishonest l2 outputs are half correct, half incorrect.
        uint256[] memory dishonestL2Outputs = new uint256[](16);
        for (uint256 i; i < dishonestL2Outputs.length; i++) {
            dishonestL2Outputs[i] = i > 7 ? 0xFF : i + 1;
        }
        // The dishonest trace is half correct, and correct all the way up to the final instruction of the exec
        // subgame.
        bytes memory dishonestTrace = new bytes(256);
        for (uint256 i; i < dishonestTrace.length; i++) {
            dishonestTrace[i] = i > (127 + 7) ? bytes1(0xFF) : bytes1(uint8(i));
        }

        // Run the actor test
        _actorTest({
            _rootClaim: 0xFF,
            _absolutePrestateData: 0,
            _honestTrace: honestTrace,
            _honestL2Outputs: honestL2Outputs,
            _dishonestTrace: dishonestTrace,
            _dishonestL2Outputs: dishonestL2Outputs,
            _expectedStatus: GameStatus.CHALLENGER_WINS
        });
    }

    ////////////////////////////////////////////////////////////////
    //                          HELPERS                           //
    ////////////////////////////////////////////////////////////////

    /// @dev Helper to run a 1v1 actor test
    function _actorTest(
        uint256 _rootClaim,
        uint256 _absolutePrestateData,
        bytes memory _honestTrace,
        uint256[] memory _honestL2Outputs,
        bytes memory _dishonestTrace,
        uint256[] memory _dishonestL2Outputs,
        GameStatus _expectedStatus
    )
        internal
    {
        // Setup the environment
        bytes memory absolutePrestateData =
            _setup({ _absolutePrestateData: _absolutePrestateData, _rootClaim: _rootClaim });

        // Create actors
        _createActors({
            _honestTrace: _honestTrace,
            _honestPreStateData: absolutePrestateData,
            _honestL2Outputs: _honestL2Outputs,
            _dishonestTrace: _dishonestTrace,
            _dishonestPreStateData: absolutePrestateData,
            _dishonestL2Outputs: _dishonestL2Outputs
        });

        // Exhaust all moves from both actors
        _exhaustMoves();

        // Resolve the game and assert that the defender won
        _warpAndResolve();
        assertEq(uint8(gameProxy.status()), uint8(_expectedStatus));
    }

    /// @dev Helper to setup the 1v1 test
    function _setup(
        uint256 _absolutePrestateData,
        uint256 _rootClaim
    )
        internal
        returns (bytes memory absolutePrestateData_)
    {
        absolutePrestateData_ = abi.encode(_absolutePrestateData);
        Claim absolutePrestateExec =
            _changeClaimStatus(Claim.wrap(keccak256(absolutePrestateData_)), VMStatuses.UNFINISHED);
        Claim rootClaim = Claim.wrap(bytes32(uint256(_rootClaim)));
        super.init({
            rootClaim: rootClaim,
            absolutePrestate: absolutePrestateExec,
            l2BlockNumber: _rootClaim,
            genesisBlockNumber: 0,
            genesisOutputRoot: Hash.wrap(bytes32(0))
        });
    }

    /// @dev Helper to create actors for the 1v1 dispute.
    function _createActors(
        bytes memory _honestTrace,
        bytes memory _honestPreStateData,
        uint256[] memory _honestL2Outputs,
        bytes memory _dishonestTrace,
        bytes memory _dishonestPreStateData,
        uint256[] memory _dishonestL2Outputs
    )
        internal
    {
        honest = new HonestDisputeActor({
            _gameProxy: gameProxy,
            _l2Outputs: _honestL2Outputs,
            _trace: _honestTrace,
            _preStateData: _honestPreStateData
        });
        dishonest = new HonestDisputeActor({
            _gameProxy: gameProxy,
            _l2Outputs: _dishonestL2Outputs,
            _trace: _dishonestTrace,
            _preStateData: _dishonestPreStateData
        });

        vm.deal(address(honest), 100 ether);
        vm.deal(address(dishonest), 100 ether);
        vm.label(address(honest), "HonestActor");
        vm.label(address(dishonest), "DishonestActor");
    }

    /// @dev Helper to exhaust all moves from both actors.
    function _exhaustMoves() internal {
        while (true) {
            // Allow the dishonest actor to make their moves, and then the honest actor.
            (uint256 numMovesA,) = dishonest.move();
            (uint256 numMovesB, bool success) = honest.move();

            require(success, "Honest actor's moves should always be successful");

            // If both actors have run out of moves, we're done.
            if (numMovesA == 0 && numMovesB == 0) break;
        }
    }

    /// @dev Helper to warp past the chess clock and resolve all claims within the dispute game.
    function _warpAndResolve() internal {
        // Warp past the chess clock
        vm.warp(block.timestamp + 3 days + 12 hours + 1 seconds);

        // Resolve all claims in reverse order. We allow `resolveClaim` calls to fail due to
        // the check that prevents claims with no subgames attached from being passed to
        // `resolveClaim`. There's also a check in `resolve` to ensure all children have been
        // resolved before global resolution, which catches any unresolved subgames here.
        for (uint256 i = gameProxy.claimDataLen(); i > 0; i--) {
            (bool success,) = address(gameProxy).call(abi.encodeCall(gameProxy.resolveClaim, (i - 1)));
            assertTrue(success);
        }
        gameProxy.resolve();
    }
}

contract ClaimCreditReenter {
    FaultDisputeGame internal immutable GAME;
    uint256 public numCalls;

    constructor(FaultDisputeGame _gameProxy) {
        GAME = _gameProxy;
    }

    function claimCredit(address _recipient) public {
        numCalls += 1;
        GAME.claimCredit(_recipient);
    }

    receive() external payable {
        if (numCalls == 5) {
            return;
        }
        claimCredit(address(this));
    }
}

/// @dev Helper to change the VM status byte of a claim.
function _changeClaimStatus(Claim _claim, VMStatus _status) pure returns (Claim out_) {
    assembly {
        out_ := or(and(not(shl(248, 0xFF)), _claim), shl(248, _status))
    }
}
