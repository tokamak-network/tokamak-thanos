# L2 자산 보호를 위한 Escape Hatch(비상 탈출구) 설계 및 구현 가이드

## 1. 개요 (Introduction)
Thanos(Bedrock) 아키텍처에서 L2 시퀀서가 중단되거나 특정 사용자의 트랜잭션을 의도적으로 검열할 경우, 사용자는 일반적인 방법으로 자산을 인출할 수 없습니다. 이를 해결하기 위해 L1에서 직접 명령을 내려 L2 자산을 안전하게 회수하는 **Escape Hatch(비상 탈출구)** 메커니즘을 정의합니다.

## 2. 설계 원칙 (Design Principles)
*   **Forced Inclusion 활용**: L1 `OptimismPortal2.depositTransaction`을 통해 시퀀서를 거치지 않고 L2 컨트랙트에 명령을 전달합니다.
*   **Address Aliasing 대응**: 호출자가 L1 컨트랙트일 경우 변조되는 주소(Alias)를 고려하여 권한을 검증합니다.
*   **표준 브릿지 연동**: 자체 인출 로직 대신 검증된 `L2StandardBridge`를 호출하여 자산 탈출의 안정성을 확보합니다.

---

## 3. 핵심 인터페이스 (`IEscapeHatchVault.sol`)

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IEscapeHatchVault {
    /**
     * @notice L2 시퀀서를 거치지 않고 L1 강제 메시지를 통해 자산을 인출합니다.
     * @dev msg.sender가 L2의 자산 소유자여야 합니다.
     */
    function emergencyExit() external;

    /**
     * @notice L1 컨트랙트(Multisig 등)가 소유한 자산을 Aliasing된 주소를 통해 인출합니다.
     * @param _l1OwnerAddress 실제 자산을 소유한 L1의 컨트랙트 주소
     */
    function emergencyExitForL1Contract(address _l1OwnerAddress) external;

    /**
     * @notice 비상 탈출구가 성공적으로 실행되었을 때 발생하는 이벤트
     */
    event EscapeHatchTriggered(address indexed user, uint256 amount);
}
```

---

## 4. 상세 구현 (`L2EscapeHatchVault.sol`)

```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { IL2StandardBridge } from "src/L2/interfaces/IL2StandardBridge.sol";
import { AddressAliasHelper } from "src/vendor/AddressAliasHelper.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";

contract L2EscapeHatchVault is IEscapeHatchVault {
    address public constant L2_BRIDGE = Predeploys.L2_STANDARD_BRIDGE;
    mapping(address => uint256) public balances;
    IERC20 public immutable l2Token;

    constructor(address _l2Token) {
        l2Token = IERC20(_l2Token);
    }

    function emergencyExit() external override {
        _performExit(msg.sender);
    }

    function emergencyExitForL1Contract(address _l1OwnerAddress) external override {
        address expectedAliasedSender = AddressAliasHelper.applyL1ToL2Alias(_l1OwnerAddress);
        require(msg.sender == expectedAliasedSender, "EscapeHatch: sender is not aliased owner");
        _performExit(_l1OwnerAddress);
    }

    function _performExit(address _owner) internal {
        uint256 amount = balances[_owner];
        require(amount > 0, "EscapeHatch: no balance");

        balances[_owner] = 0;
        l2Token.approve(L2_BRIDGE, amount);

        // L1으로 전송하기 위해 L2 표준 브릿지 호출
        IL2StandardBridge(L2_BRIDGE).withdraw(
            address(l2Token),
            amount,
            200_000,
            ""
        );

        emit EscapeHatchTriggered(_owner, amount);
    }
}
```

---

## 5. 실행 절차 (Operational Flow)

### Step 1: L1에서 강제 트랜잭션 삽입
시퀀서가 멈췄을 때 사용자는 L1 Ethereum 상에서 `OptimismPortal2.depositTransaction`을 호출합니다.

| 파라미터 | 값 | 설명 |
| :--- | :--- | :--- |
| **`_to`** | `Vault_Address` | 탈출 기능을 지원하는 L2 컨트랙트 주소 |
| **`_gasLimit`** | `150,000` | L2 실행에 필요한 가스양 |
| **`_data`** | `0xbe4a2e55` | `emergencyExit()` 함수 호출 데이터 |

### Step 2: L1 인출 증명 및 확정
1.  **시중 지연(Waiting)**: Thanos 노드가 L1 이벤트를 감지하여 L2 블록에 포함시킬 때까지 대기합니다.
2.  **증명(Prove)**: L2에서 탈출 메시지가 생성되면 L1에서 `proveWithdrawalTransaction`을 호출합니다.
3.  **확정(Finalize)**: 7일의 지연 기간(Fault Proof 대기)이 지난 후 `finalizeWithdrawalTransaction`을 통해 L1 자산을 수령합니다.

---

## 6. 보안 및 기대 효과
1.  **검열 저항성(Censorship Resistance)**: 특정 주소의 트랜잭션을 거부하는 악의적인 시퀀서로부터 자산을 보호합니다.
2.  **보안 등급 격상**: L2의 가용성(Availability) 장애 상황에서도 사용자 스스로 자산을 회수할 수 있게 함으로써 Stage 2 수준의 탈중앙화 보안을 달성합니다.
3.  **사용자 신뢰 확보**: 코드 레벨에서 보장되는 탈출 경로는 대규모 TVL 유치를 위한 필수적인 신뢰 기반이 됩니다.
