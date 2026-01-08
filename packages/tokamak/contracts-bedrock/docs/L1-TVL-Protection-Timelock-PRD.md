# [PRD] L1 TVL 보호를 위한 Timelock 시스템 도입

## 1. 개요 (Introduction)
Thanos(Bedrock) 아키텍처는 L1 브릿지에 막대한 양의 자산(TVL)을 보관합니다. 현재의 관리자 권한(Operator/Guardian) 모델은 빠른 운영을 가능하게 하지만, 권한 탈취나 운영자의 변심 시 자산을 즉시 유출시킬 수 있는 리스크가 존재합니다. 본 제품 요건서(PRD)는 이를 해결하기 위해 **Timelock(시간 지연)** 메커니즘을 도입하여 시스템의 보안 등급을 Stage 2 수준으로 격상시키는 것을 목표로 합니다.

## 2. 문제 정의 (Problem Statement)
*   **보안 계층의 불균형**: Fault Proof(결함 증명)는 데이터 조작을 막지만, 코드 업그레이드 권한을 가진 운영자가 브릿지 로직 자체를 조작하는 행위(Backdoor 등)는 막을 수 없습니다.
*   **사용자의 대응 시간 부재**: 운영자가 악의적인 업그레이드를 실행할 때 지연 시간이 없다면, 사용자는 자신의 자산이 탈취되는 것을 인지하더라도 L1으로 대피할 기회를 얻지 못합니다.
*   **권한의 중앙화**: 소수(또는 1인)의 관리자 키가 노출될 경우 전체 시스템의 자산 안전성이 즉각적으로 무너집니다.

## 3. 목표 (Goals)
*   **자산 탈취의 물리적 불가능성**: TVL 조작이 가능한 핵심 함수에 최소 7일의 실행 지연을 강제합니다.
*   **사용자 주권 확보**: 사용자가 악의적 업그레이드를 인지하고 **Emergency Exit(Fault Proof 기반 인출)**을 완료할 수 있는 유예 기간을 실질적으로 보장합니다.
*   **권한 분리(Separation of Powers)**: Guardian은 시스템을 멈출 수는 있으나, Timelock을 무력화하거나 규칙을 바꿀 수는 없어야 합니다.
*   **투명성 및 신뢰도 제고**: 모든 인프라 변경 사항을 사전에 온체인 공시하여 커뮤니티의 검토를 받습니다.

## 4. 구현 범위 (Scope)

### 4.1 대상 컨트랙트 및 함수 (5대 핵심 함수)
| 계층 | 컨트랙트 | 대상 함수 (Functions) |
| :--- | :--- | :--- |
| **프록시 제어** | `ProxyAdmin` | `upgrade`, `upgradeAndCall`, `changeProxyAdmin`, `transferOwnership` |
| **인출 검증** | `OptimismPortal2` | `setDisputeGameFactory` (신규 구현), `setRespectedGameType` |

### 4.2 인프라 컴포넌트
*   **L1TimelockController**: OpenZeppelin 표준 라이브러리를 기반으로 한 시간 지연 제어기.
*   **모니터링 알림 봇**: Timelock 이벤트(CallScheduled) 감지 및 커뮤니티 전파 시스템.

## 5. 기능 요구사항 (Functional Requirements)

### 5.1 지연 실행 로직
*   운영자가 위험 함수를 호출할 때 즉시 실행되지 않고 `schedule` 단계를 거쳐야 함.
*   설정된 `minDelay` (7일)가 지나기 전에는 어떤 방식으로도 실제 실행(`execute`)이 불가능해야 함.

### 5.2 권한 분리 및 예외 처리 (Separation of Powers)
*   **Governance 계층 (Timelock 전용 - 7일 지연)**: 업그레이드, 주소 변경 등 '규칙'을 바꾸는 함수(`setRespectedGameType`, `setDisputeGameFactory`, `upgrade` 등) 관리. Guardian은 이 함수들을 실행할 수 없음.
*   **Guardian 계층 (즉시 실행)**: `pause()`, `blacklistDisputeGame()` 등 자산 보호를 위한 '중단' 및 '차단' 함수만 실행 가능. 이는 사고 발생 시 즉각적인 대응을 위함임.
*   **보안 설계 원칙**: Guardian이 `setRespectedGameType`을 즉시 호출하여 사용자의 Emergency Exit을 차단하는 것을 원천 방지함.

### 5.3 정보 공시 및 가독성
*   Timelock에 예약된 비트(Calldata)를 사용자가 이해할 수 있는 텍스트로 변환하여 제공해야 함.

## 6. 기술 요구사항 (Technical Requirements)

### 6.1 `OptimismPortal2` 수정
*   현재 코드에는 `setDisputeGameFactory` 함수가 부재하므로, 이를 명시적으로 추가하여 Factory 교체 시에도 Timelock이 적용되도록 함.
*   기존 Guardian 전용 함수들을 `onlyGovernance` (Timelock 포함) 수식어로 확장.

### 6.2 `ProxyAdmin` 통합
*   `ProxyAdmin`의 `owner`를 새로 배포된 `TimelockController`의 주소로 이전.

### 6.3 권한 계층 구조 (Privilege Hierarchy)
*   **구조**: `SystemOwnerSafe` (Proposer) -> `GovernanceTimelock` (Owner) -> `System Contracts`.
*   **설정**:
    - `GovernanceTimelock`의 `PROPOSER_ROLE`과 `TIMELOCK_ADMIN_ROLE`을 `SystemOwnerSafe`에 부여.
    - 배포자(Deployer)의 모든 권한은 배포 직후 폐기(Renounce).
    - 핵심 컨트랙트(`ProxyAdmin`, `OptimismPortal2` 등)의 소유권은 단일 `GovernanceTimelock`이 보유.

### 6.4 보안 표준 준수
*   OpenZeppelin `TimelockController`의 최신 버전을 사용하여 이미 검증된 로직을 재사용함.

## 7. 구현 마일스톤 (Milestones)
1.  **M1. 코드 개발 및 유닛 테스트**: `OptimismPortal2` 수정 및 Timelock 연동 테스트 완료.
2.  **M2. 테스트넷(Sepolia) 배포**: 실제 7일 대기 기간을 포함한 업그레이드 시나리오 검증.
3.  **M3. 모니터링 대시보드 구축**: 예약된 제안을 실시간으로 보여주는 웹 인터페이스 개발.
4.  **M4. 메인넷 활성화**: 소유권 이전을 위한 7일 공시 후 최종 거버넌스 권한 위임.

## 8. 성공 지표 (Success Metrics)
*   **보안성**: 외부 감사를 통해 위 5대 함수에 대한 직접적인 호출 가능성이 0임이 입증되어야 함.
*   **투명성**: 모든 업그레이드 제안이 예약된 시점부터 1시간 이내에 커뮤니티 알림 채널에 공시되어야 함.
*   **신뢰도**: L2 BEATS 등의 롤업 평가 기준에서 'Code Upgrade' 항목의 리스크가 'Yellow' 이하로 개선되어야 함.
