// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { VRFCoordinator } from "../../../src/tokamak/DRB/VRFCoordinator.sol";
import { VRFPredeploy } from "../../../src/tokamak/DRB/VRFPredeploy.sol";
import { VRFConsumerBase } from "../../../src/tokamak/DRB/VRFConsumerBase.sol";

contract MockConsumer is VRFConsumerBase {
    uint256 public lastRequestId;
    uint256[] public lastRandomWords;

    constructor(address coordinator) VRFConsumerBase(coordinator) {}

    function fulfillRandomWords(uint256 requestId, uint256[] memory randomWords) internal override {
        lastRequestId = requestId;
        lastRandomWords = randomWords;
    }
}

contract VRFCoordinatorTest is Test {
    VRFCoordinator coordinator;
    VRFPredeploy vrfPredeploy;
    MockConsumer consumer;
    address admin = address(0xAD);
    address node = address(0xB0DE);

    function setUp() public {
        coordinator = new VRFCoordinator();
        coordinator.initialize(admin);

        vrfPredeploy = new VRFPredeploy();
        vrfPredeploy.initialize(address(coordinator));

        consumer = new MockConsumer(address(coordinator));
    }

    function test_requestAndFulfill() public {
        vm.prank(admin);
        coordinator.registerNode(node);

        uint256 reqId = vrfPredeploy.requestRandomWords(1, 200_000);
        // Note: consumer.address is the requester in the real flow;
        // here we request via vrfPredeploy which records msg.sender as requester.
        // For simplicity test fulfillment via coordinator directly.
        assertEq(reqId, 1);

        uint256[] memory words = new uint256[](1);
        words[0] = 42;

        // Fulfill via the coordinator (node calling)
        vm.prank(node);
        // The requester is vrfPredeploy itself in this test (msg.sender of requestRandomWords call)
        // We need coordinator to call rawFulfillRandomWords on requester
        // So make consumer the requester by requesting from coordinator directly
        uint256 reqId2 = coordinator.requestRandomWords(address(consumer), 1, 200_000);
        vm.prank(node);
        coordinator.fulfillRandomWords(reqId2, words);

        assertEq(consumer.lastRequestId(), reqId2);
        assertEq(consumer.lastRandomWords(0), 42);
    }

    function test_onlyNodeCanFulfill() public {
        uint256 reqId = coordinator.requestRandomWords(address(consumer), 1, 200_000);

        uint256[] memory words = new uint256[](1);
        words[0] = 42;

        vm.expectRevert("VRFCoordinator: only registered node");
        coordinator.fulfillRandomWords(reqId, words);
    }

    function test_onlyAdminCanRegisterNode() public {
        vm.expectRevert("VRFCoordinator: only admin");
        coordinator.registerNode(address(0x1234));
    }

    function test_cannotFulfillTwice() public {
        vm.prank(admin);
        coordinator.registerNode(node);

        uint256 reqId = coordinator.requestRandomWords(address(consumer), 1, 200_000);
        uint256[] memory words = new uint256[](1);
        words[0] = 42;

        vm.prank(node);
        coordinator.fulfillRandomWords(reqId, words);

        vm.expectRevert("VRFCoordinator: already fulfilled");
        vm.prank(node);
        coordinator.fulfillRandomWords(reqId, words);
    }
}
