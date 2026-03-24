// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

// import {Owned} from "@solmate/src/auth/Owned.sol";
import {Ownable} from "@solady/auth/Ownable.sol";

contract OperatorManager is Ownable {
    // * State Variables

    mapping(address operator => uint256) public s_depositAmount;
    mapping(address operator => uint256) public s_activatedOperatorIndex1Based;

    uint256 public s_slashRewardPerOperatorX8;
    mapping(address => uint256) public s_slashRewardPerOperatorPaidX8;

    uint256 public s_isInProcess = COMPLETED;
    uint256 public s_activationThreshold;

    // ** internal variables
    /**
     * @notice Holds the list of currently active operator addresses.
     * @dev
     *   - Operators appear in this array once they call `activate()` successfully.
     *   - The 1-based index for each operator is stored in `s_activatedOperatorIndex1Based`.
     *   - When an operator is deactivated, it is removed from this array (replacing
     *     its slot with the last operator in the array).
     */
    address[] internal s_activatedOperators;

    // ** constants
    uint256 internal constant IN_PROGRESS = 1;
    uint256 internal constant COMPLETED = 2;
    uint256 internal constant HALTED = 3;
    uint256 internal constant MAX_OPERATOR_INDEX = 31;
    uint256 public constant MAX_ACTIVATED_OPERATORS = 32;

    // ** Events
    event Activated(address operator); // 0x0cc43938d137e7efade6a531f663e78c1fc75257b0d65ffda2fdaf70cb49cdf9
    event DeActivated(address operator); // 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d

    // * Errors
    error TransferFailed();
    error InProcess(); // 0x0f56c325
    error OnlyActivatedOperatorCanClaim(); // 0x111fa29f
    error OwnerCannotActivate(); // 0x4534ad7f
    error LessThanActivationThreshold(); // 0x5af30906
    error AlreadyActivated(); // 0xef65161f
    error ActivatedOperatorsLimitReached(); // 0x3e8fbd5f
    error WithdrawAmountIsZero(); // 0xa393d14b
    error PendingOwnerCannotBeActivatedOperator(); // 0x5df6bf29

    constructor() {
        _initializeOwner(msg.sender);
    }

    /**
     * @notice Ensures that no actions can be taken while the contract is in an ongoing process.
     * @dev
     *   - Reverts with {InProcess} if `s_isInProcess == IN_PROGRESS`.
     *   - Commonly used to protect functions that should not execute while the system is ongoing
     *     with a round of operations or an uncompleted flow.
     */
    modifier notInProcess() {
        assembly ("memory-safe") {
            if eq(sload(s_isInProcess.slot), IN_PROGRESS) {
                mstore(0x00, 0x0f56c325) // `InProcess()`.
                revert(0x1c, 0x04)
            }
        }
        _;
    }

    // ** Override Ownable Functions
    function transferOwnership(address newOwner) public payable override onlyOwner {
        assembly ("memory-safe") {
            mstore(0x00, newOwner)
            mstore(0x20, s_activatedOperatorIndex1Based.slot)
            if gt(sload(keccak256(0x00, 0x40)), 0) {
                mstore(0x00, 0x9279dd8e) // NewOwnerCannotBeActivatedOperator()
                revert(0x1c, 0x04)
            }
            if iszero(shl(96, newOwner)) {
                mstore(0x00, 0x7448fbae) // `NewOwnerIsZeroAddress()`.
                revert(0x1c, 0x04)
            }
        }
        _setOwner(newOwner);
    }

    function requestOwnershipHandover() public payable override {
        unchecked {
            uint256 expires = block.timestamp + _ownershipHandoverValidFor();
            /// @solidity memory-safe-assembly
            assembly {
                mstore(0x00, caller())
                mstore(0x20, s_activatedOperatorIndex1Based.slot)
                if gt(sload(keccak256(0x00, 0x40)), 0) {
                    mstore(0x00, 0x9279dd8e) // NewOwnerCannotBeActivatedOperator()
                    revert(0x1c, 0x04)
                }
                // Compute and set the handover slot to `expires`.
                mstore(0x0c, 0x389a75e1) // _HANDOVER_SLOT_SEED
                mstore(0x00, caller())
                sstore(keccak256(0x0c, 0x20), expires)
                // Emit the {OwnershipHandoverRequested} event.
                log2(0, 0, 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d, caller()) // _OWNERSHIP_HANDOVER_REQUESTED_EVENT_SIGNATURE
            }
        }
    }

    function completeOwnershipHandover(address pendingOwner) public payable override onlyOwner {
        /// @solidity memory-safe-assembly
        assembly {
            mstore(0x00, pendingOwner)
            mstore(0x20, s_activatedOperatorIndex1Based.slot)
            if gt(sload(keccak256(0x00, 0x40)), 0) {
                mstore(0x00, 0x5df6bf29) // PendingOwnerCannotBeActivatedOperator()
                revert(0x1c, 0x04)
            }
            // Compute and set the handover slot to 0.
            mstore(0x0c, 0x389a75e1) // _HANDOVER_SLOT_SEED
            mstore(0x00, pendingOwner)
            let handoverSlot := keccak256(0x0c, 0x20)
            // If the handover does not exist, or has expired.
            if gt(timestamp(), sload(handoverSlot)) {
                mstore(0x00, 0x6f5e8818) // `NoHandoverRequest()`.
                revert(0x1c, 0x04)
            }
            // Set the handover slot to 0.
            sstore(handoverSlot, 0)
        }
        _setOwner(pendingOwner);
    }

    function deposit() public payable {
        assembly ("memory-safe") {
            mstore(0x00, caller())
            mstore(0x20, s_depositAmount.slot)
            let depositAmountSlot := keccak256(0x00, 0x40)
            sstore(depositAmountSlot, add(sload(depositAmountSlot), callvalue()))
        }
    }

    function activate() public notInProcess {
        assembly ("memory-safe") {
            // Check if the caller's deposit amount is less than the activation threshold and revert if true.
            mstore(0x00, caller())
            mstore(0x20, s_depositAmount.slot)
            if lt(sload(keccak256(0x00, 0x40)), sload(s_activationThreshold.slot)) {
                mstore(0x00, 0x5af30906) // `LessThanActivationThreshold()`.
                revert(0x1c, 0x04)
            }
        }
        _activate();
    }

    function _activate() internal {
        assembly ("memory-safe") {
            // Check if the caller is the owner (leader node) and revert if true.
            if eq(caller(), sload(_OWNER_SLOT)) {
                mstore(0x00, 0x4534ad7f) // `OwnerCannotActivate()`.
                revert(0x1c, 0x04)
            }

            // Check if the caller is already an activated operator and revert if true.
            let newLength := add(sload(s_activatedOperators.slot), 1)
            if gt(newLength, MAX_ACTIVATED_OPERATORS) {
                mstore(0x00, 0x3e8fbd5f) // `ActivatedOperatorsLimitReached()`.
                revert(0x1c, 0x04)
            }
            mstore(0x00, caller())
            mstore(0x20, s_activatedOperatorIndex1Based.slot)
            let activatedOperatorIndex1BasedSlot := keccak256(0x00, 0x40)
            if gt(sload(activatedOperatorIndex1BasedSlot), 0) {
                mstore(0x00, 0xef65161f) // `AlreadyActivated()`.
                revert(0x1c, 0x04)
            }
            sstore(activatedOperatorIndex1BasedSlot, newLength)
            // Update the length of the activated operators array.
            sstore(s_activatedOperators.slot, newLength)
            // Store the caller's address in the activated operators array.
            mstore(0x20, s_activatedOperators.slot)
            sstore(add(keccak256(0x20, 0x20), sub(newLength, 1)), caller())

            // initialize slashRewardPerOperatorPaid
            mstore(0x20, s_slashRewardPerOperatorPaidX8.slot) // caller is already in memory 0x00
            sstore(keccak256(0x00, 0x40), sload(s_slashRewardPerOperatorX8.slot))
            log1(0x00, 0x20, 0x0cc43938d137e7efade6a531f663e78c1fc75257b0d65ffda2fdaf70cb49cdf9)
        }
    }

    function depositAndActivate() external payable virtual notInProcess {
        assembly ("memory-safe") {
            mstore(0x00, caller())
            mstore(0x20, s_depositAmount.slot)
            let depositAmountSlot := keccak256(0x00, 0x40)
            let updatedDepositAmount := add(sload(depositAmountSlot), callvalue())
            if lt(updatedDepositAmount, sload(s_activationThreshold.slot)) {
                mstore(0x00, 0x5af30906) // `LessThanActivationThreshold()`.
                revert(0x1c, 0x04)
            }
            sstore(depositAmountSlot, updatedDepositAmount)
        }
        _activate();
    }

    function withdraw() external notInProcess {
        assembly ("memory-safe") {
            mstore(0x00, caller())
            mstore(0x20, s_depositAmount.slot)
            let withdrawAmount := sload(keccak256(0x00, 0x40))
            mstore(0x20, s_activatedOperatorIndex1Based.slot)
            let activatedOperatorIndex1Based := sload(keccak256(0x00, 0x40))
            let currentSlashRewardPerOperatorX8 := sload(s_slashRewardPerOperatorX8.slot)
            mstore(0x20, s_slashRewardPerOperatorPaidX8.slot)

            if gt(activatedOperatorIndex1Based, 0) {
                // ** update withdraw amount
                withdrawAmount :=
                    add(withdrawAmount, shr(8, sub(currentSlashRewardPerOperatorX8, sload(keccak256(0x00, 0x40)))))
                // ** deactivate msg.sender
                mstore(0x00, s_activatedOperators.slot)
                let firstActivatedOperatorSlot := keccak256(0x00, 0x20)
                let lastOperatorIndex := sub(sload(s_activatedOperators.slot), 1)
                let lastOperatorAddress := sload(add(firstActivatedOperatorSlot, lastOperatorIndex))
                // swap the operator to remove with the last operator in the array
                mstore(0x20, s_activatedOperatorIndex1Based.slot)
                if iszero(eq(lastOperatorAddress, caller())) {
                    sstore(add(firstActivatedOperatorSlot, sub(activatedOperatorIndex1Based, 1)), lastOperatorAddress)
                    mstore(0x00, lastOperatorAddress)
                    sstore(keccak256(0x00, 0x40), activatedOperatorIndex1Based)
                }
                // pop the last operator by setting the length to the last index
                sstore(s_activatedOperators.slot, lastOperatorIndex)
                mstore(0x00, caller())
                sstore(keccak256(0x00, 0x40), 0)
                log1(0x00, 0x20, 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d) // `DeActivated(address operator)`.
            }
            if eq(caller(), sload(_OWNER_SLOT)) {
                // If the caller is the owner (leader node) but not an operator,
                // they can still withdraw deposit plus slash reward.
                withdrawAmount :=
                    add(withdrawAmount, shr(8, sub(currentSlashRewardPerOperatorX8, sload(keccak256(0x00, 0x40)))))
            }
            if iszero(withdrawAmount) {
                mstore(0x00, 0xa393d14b) // `WithdrawAmountIsZero()`.
                revert(0x1c, 0x04)
            }
            mstore(0x20, s_slashRewardPerOperatorPaidX8.slot)
            sstore(keccak256(0x00, 0x40), currentSlashRewardPerOperatorX8)
            // Reset deposit to zero and attempt transfer
            mstore(0x20, s_depositAmount.slot)
            sstore(keccak256(0x00, 0x40), 0)
            // Transfer the ETH and check if it succeeded or not.
            if iszero(call(gas(), caller(), withdrawAmount, 0x00, 0x00, 0x00, 0x00)) {
                mstore(0x00, 0xb12d13eb) // `ETHTransferFailed()`.
                revert(0x1c, 0x04)
            }
        }
    }

    function deactivate() external notInProcess {
        // Note: Intentionally no operator activation check for gas optimization.
        // Non-activated operators have s_activatedOperatorIndex1Based[msg.sender] = 0,
        // causing underflow (0 - 1) which serves as implicit validation and reverts.
        // This design prioritizes gas efficiency over verbose error messages.
        _deactivate(s_activatedOperatorIndex1Based[msg.sender] - 1, msg.sender);
        assembly ("memory-safe") {
            let currentSlashRewardPerOperatorX8 := sload(s_slashRewardPerOperatorX8.slot)
            mstore(0x00, caller())
            mstore(0x20, s_depositAmount.slot)
            let depositAmountSlot := keccak256(0x00, 0x40)
            mstore(0x20, s_slashRewardPerOperatorPaidX8.slot)
            let slashRewardPerOperatorPaidX8Slot := keccak256(0x00, 0x40)
            sstore(
                depositAmountSlot,
                add(
                    sload(depositAmountSlot),
                    shr(8, sub(currentSlashRewardPerOperatorX8, sload(slashRewardPerOperatorPaidX8Slot)))
                )
            )
            sstore(slashRewardPerOperatorPaidX8Slot, currentSlashRewardPerOperatorX8)
        }
    }

    function claimSlashReward() external {
        assembly ("memory-safe") {
            mstore(0x00, caller())
            mstore(0x20, s_activatedOperatorIndex1Based.slot)
            if iszero(sload(keccak256(0x00, 0x40))) {
                mstore(0x00, 0x111fa29f) // `OnlyActivatedOperatorCanClaim()`.
                revert(0x1c, 0x04)
            }
            let currentSlashRewardPerOperatorX8 := sload(s_slashRewardPerOperatorX8.slot)
            mstore(0x20, s_slashRewardPerOperatorPaidX8.slot)
            let slashRewardPerOperatorPaidX8Slot := keccak256(0x00, 0x40)
            let slashRewardAmount :=
                shr(8, sub(currentSlashRewardPerOperatorX8, sload(slashRewardPerOperatorPaidX8Slot)))
            if iszero(slashRewardAmount) {
                mstore(0x00, 0xa393d14b) // `WithdrawAmountIsZero()`.
                revert(0x1c, 0x04)
            }
            sstore(slashRewardPerOperatorPaidX8Slot, currentSlashRewardPerOperatorX8)
            // Transfer the ETH and check if it succeeded or not.
            if iszero(call(gas(), caller(), slashRewardAmount, 0x00, 0x00, 0x00, 0x00)) {
                mstore(0x00, 0xb12d13eb) // `ETHTransferFailed()`.
                revert(0x1c, 0x04)
            }
        }
    }

    // ** getter
    function getActivatedOperators() external view returns (address[] memory) {
        return s_activatedOperators;
    }

    function getActivatedOperatorsLength() external view returns (uint256) {
        return s_activatedOperators.length;
    }

    // ** For Testing

    function getDepositPlusSlashReward(address operator) external view returns (uint256) {
        if (owner() != operator && s_activatedOperatorIndex1Based[operator] == 0) {
            return s_depositAmount[operator];
        }
        return
            s_depositAmount[operator] + ((s_slashRewardPerOperatorX8 - s_slashRewardPerOperatorPaidX8[operator]) >> 8);
    }

    /**
     * @notice Removes an operator from the `s_activatedOperators` array and resets their mapping indices.
     * @dev
     *   - Internal function used by either `deactivate()` or `withdraw()`.
     *   - Swaps the operator to remove with the last operator in the array, preserving continuous storage,
     *     then `pop()`s the array end.
     *   - Clears `s_activatedOperatorIndex1Based[operator]` to indicate the operator is no longer active.
     *   - Emits {DeActivated} with the removed operator’s address.
     * @param activatedOperatorIndex The zero-based index of the operator in `s_activatedOperators`.
     * @param operator The address of the operator to remove.
     */
    function _deactivate(uint256 activatedOperatorIndex, address operator) internal {
        assembly ("memory-safe") {
            mstore(0x00, s_activatedOperators.slot)
            let firstActivatedOperatorSlot := keccak256(0x00, 0x20)
            let lastOperatorIndex := sub(sload(s_activatedOperators.slot), 1)
            let lastOperatorAddress := sload(add(firstActivatedOperatorSlot, lastOperatorIndex))
            mstore(0x20, s_activatedOperatorIndex1Based.slot)
            // swap the operator to remove with the last operator in the array
            if iszero(eq(lastOperatorAddress, operator)) {
                sstore(add(firstActivatedOperatorSlot, activatedOperatorIndex), lastOperatorAddress)
                mstore(0x00, lastOperatorAddress)
                sstore(keccak256(0x00, 0x40), add(activatedOperatorIndex, 1))
            }
            // pop the last operator by setting the length to the last index
            sstore(s_activatedOperators.slot, lastOperatorIndex)
            mstore(0x00, operator)
            sstore(keccak256(0x00, 0x40), 0)
            log1(0x00, 0x20, 0x5d10eb48d8c00fb4cc9120533a99e2eac5eb9d0f8ec06216b2e4d5b1ff175a4d) // `DeActivated(address operator)`.
        }
    }
}
