// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Initializable } from "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import { VRFConsumerBase } from "./VRFConsumerBase.sol";

/// @title VRFCoordinator
/// @notice Manages DRB nodes and fulfills random word requests.
///         Nodes register via registerNode(), submit results via fulfillRandomWords().
///         Admin = SystemConfig owner.
contract VRFCoordinator is Initializable {
    address public admin;
    address public vrfPredeploy;
    uint256 private _requestCounter;

    struct Request {
        address requester;
        uint32 numWords;
        uint256 callbackGasLimit;
        bool fulfilled;
        uint256[] randomWords;
    }

    mapping(address => bool) public registeredNodes;
    mapping(uint256 => Request) public requests;

    event NodeRegistered(address indexed node);
    event NodeDeregistered(address indexed node);
    event RandomWordsRequested(uint256 indexed requestId, address indexed requester, uint32 numWords);
    event RandomWordsFulfilled(uint256 indexed requestId, uint256[] randomWords);

    modifier onlyAdmin() {
        require(msg.sender == admin, "VRFCoordinator: only admin");
        _;
    }

    modifier onlyNode() {
        require(registeredNodes[msg.sender], "VRFCoordinator: only registered node");
        _;
    }

    modifier onlyVRFPredeploy() {
        require(msg.sender == vrfPredeploy, "VRFCoordinator: only VRFPredeploy");
        _;
    }

    /// @notice Initializer (called once by proxy).
    function initialize(address _admin) external initializer {
        require(_admin != address(0), "VRFCoordinator: zero admin");
        admin = _admin;
    }

    /// @notice Set the VRFPredeploy address (admin only). Resolves circular init dependency.
    function setPredeploy(address _vrfPredeploy) external onlyAdmin {
        require(_vrfPredeploy != address(0), "VRFCoordinator: zero predeploy");
        vrfPredeploy = _vrfPredeploy;
    }

    /// @notice Register a DRB node (admin only).
    function registerNode(address node) external onlyAdmin {
        registeredNodes[node] = true;
        emit NodeRegistered(node);
    }

    /// @notice Deregister a DRB node (admin only).
    function deregisterNode(address node) external onlyAdmin {
        registeredNodes[node] = false;
        emit NodeDeregistered(node);
    }

    /// @notice Called by VRFPredeploy to record a randomness request.
    function requestRandomWords(address requester, uint32 numWords, uint256 callbackGasLimit)
        external
        onlyVRFPredeploy
        returns (uint256 requestId)
    {
        require(requester != address(0), "VRFCoordinator: zero requester");
        require(numWords > 0 && numWords <= 10, "VRFCoordinator: invalid numWords");

        requestId = ++_requestCounter;
        requests[requestId] = Request({
            requester: requester,
            numWords: numWords,
            callbackGasLimit: callbackGasLimit,
            fulfilled: false,
            randomWords: new uint256[](0)
        });
        emit RandomWordsRequested(requestId, requester, numWords);
    }

    /// @notice Called by DRB nodes to fulfill a randomness request.
    function fulfillRandomWords(uint256 requestId, uint256[] calldata randomWords) external onlyNode {
        Request storage req = requests[requestId];
        require(!req.fulfilled, "VRFCoordinator: already fulfilled");
        require(randomWords.length == req.numWords, "VRFCoordinator: wrong word count");

        req.fulfilled = true;
        req.randomWords = randomWords;

        emit RandomWordsFulfilled(requestId, randomWords);
        VRFConsumerBase(req.requester).rawFulfillRandomWords{gas: req.callbackGasLimit}(requestId, randomWords);
    }

    /// @notice Returns the status and result of a randomness request.
    function getRequestStatus(uint256 requestId)
        external
        view
        returns (bool fulfilled, uint256[] memory randomWords)
    {
        Request storage req = requests[requestId];
        return (req.fulfilled, req.randomWords);
    }
}
