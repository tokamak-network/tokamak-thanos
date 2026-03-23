// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { VRFCoordinator } from "../../src/DRB/VRFCoordinator.sol";
import { VRFPredeploy } from "../../src/DRB/VRFPredeploy.sol";
import { VRFConsumerBase } from "../../src/DRB/VRFConsumerBase.sol";

contract MockConsumer is VRFConsumerBase {
    uint256 public lastRequestId;
    uint256[] public lastRandomWords;
    uint256 public pendingRequestId;

    constructor(address coordinator) VRFConsumerBase(coordinator) {}

    function requestViaVRF(VRFPredeploy vrfPredeploy, uint32 numWords, uint256 gasLimit)
        external
        returns (uint256)
    {
        pendingRequestId = vrfPredeploy.requestRandomWords(numWords, gasLimit);
        return pendingRequestId;
    }

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

        vm.prank(admin);
        coordinator.setPredeploy(address(vrfPredeploy));

        consumer = new MockConsumer(address(coordinator));
    }

    function test_requestAndFulfill() public {
        vm.prank(admin);
        coordinator.registerNode(node);

        // Consumer requests via VRFPredeploy (end-to-end)
        uint256 reqId = consumer.requestViaVRF(vrfPredeploy, 1, 200_000);
        assertEq(reqId, 1);

        uint256[] memory words = new uint256[](1);
        words[0] = 42;

        vm.prank(node);
        coordinator.fulfillRandomWords(reqId, words);

        assertEq(consumer.lastRequestId(), reqId);
        assertEq(consumer.lastRandomWords(0), 42);
    }

    function test_onlyNodeCanFulfill() public {
        uint256 reqId = consumer.requestViaVRF(vrfPredeploy, 1, 200_000);

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

        uint256 reqId = consumer.requestViaVRF(vrfPredeploy, 1, 200_000);
        uint256[] memory words = new uint256[](1);
        words[0] = 42;

        vm.prank(node);
        coordinator.fulfillRandomWords(reqId, words);

        vm.expectRevert("VRFCoordinator: already fulfilled");
        vm.prank(node);
        coordinator.fulfillRandomWords(reqId, words);
    }

    function test_onlyVRFPredeployCanRequest() public {
        vm.expectRevert("VRFCoordinator: only VRFPredeploy");
        coordinator.requestRandomWords(address(consumer), 1, 200_000);
    }

    function test_setPredeploy_onlyAdmin() public {
        vm.expectRevert("VRFCoordinator: only admin");
        coordinator.setPredeploy(address(vrfPredeploy));
    }
}
