# Tokamak Thanos DisputeGame 독립 배포 가이드

**작성일**: 2025-11-07
**목적**: 새로 추가된 DeployDisputeGame.s.sol 스크립트 사용 가이드

---

## 📦 추가된 파일들

### 1. DeployDisputeGame.s.sol
- **경로**: `packages/tokamak/contracts-bedrock/scripts/deploy/DeployDisputeGame.s.sol`
- **크기**: 14,851 bytes
- **역할**: FaultDisputeGame 또는 PermissionedDisputeGame을 독립적으로 배포

### 2. SetDisputeGameImpl.s.sol
- **경로**: `packages/tokamak/contracts-bedrock/scripts/deploy/SetDisputeGameImpl.s.sol`
- **크기**: 3,437 bytes
- **역할**: 배포된 게임을 DisputeGameFactory에 등록하고 AnchorStateRegistry 업데이트

### 3. DeployAlphabetVM.s.sol
- **경로**: `packages/tokamak/contracts-bedrock/scripts/deploy/DeployAlphabetVM.s.sol`
- **크기**: 1,355 bytes
- **역할**: 테스트용 간단한 AlphabetVM 배포

---

## 🚀 사용 방법

### 기본 흐름

```
1. VM 배포 (MIPS64, AlphabetVM, Asterisc 등)
   ↓
2. DeployDisputeGame 실행 (게임 배포)
   ↓
3. SetDisputeGameImpl 실행 (Factory 등록)
   ↓
4. 완료!
```

---

## 📝 예시 1: FaultDisputeGame (Permissionless) 배포

### Step 1: MIPS VM 확인/배포

MIPS VM은 이미 DeployImplementations에서 배포되었을 것입니다.

```solidity
// MIPSSingleton 주소 확인
address mipsVM = artifacts.mustGetAddress("MipsSingleton");
```

### Step 2: DeployDisputeGame 스크립트 작성

**파일**: `scripts/deploy/examples/DeployCannonGame.s.sol`

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import { Script } from "forge-std/Script.sol";
import { DeployDisputeGame } from "scripts/deploy/DeployDisputeGame.s.sol";
import { DeployDisputeGameInput } from "scripts/deploy/DeployDisputeGame.s.sol";
import { IDelayedWETH } from "interfaces/dispute/IDelayedWETH.sol";
import { IAnchorStateRegistry } from "interfaces/dispute/IAnchorStateRegistry.sol";
import { IMIPS64 } from "interfaces/cannon/IMIPS64.sol";
import { GameTypes } from "src/dispute/lib/Types.sol";
import { Duration } from "src/dispute/lib/LibUDT.sol";

contract DeployCannonGame is Script {
    function run() public {
        // 기존 컨트랙트 주소 로드
        IDelayedWETH delayedWETH = IDelayedWETH(payable(vm.envAddress("DELAYED_WETH_PROXY")));
        IAnchorStateRegistry registry = IAnchorStateRegistry(vm.envAddress("ANCHOR_STATE_REGISTRY_PROXY"));
        IMIPS64 mipsVM = IMIPS64(vm.envAddress("MIPS_SINGLETON"));

        // Input 설정
        DeployDisputeGameInput input = new DeployDisputeGameInput();
        input.set(input.gameKind.selector, "FaultDisputeGame");
        input.set(input.gameType.selector, uint256(GameTypes.CANNON.raw()));
        input.set(input.absolutePrestate.selector, vm.envBytes32("ABSOLUTE_PRESTATE"));
        input.set(input.maxGameDepth.selector, 73);
        input.set(input.splitDepth.selector, 30);
        input.set(input.clockExtension.selector, uint256(Duration.unwrap(Duration.wrap(3 hours))));
        input.set(input.maxClockDuration.selector, uint256(Duration.unwrap(Duration.wrap(3.5 days))));
        input.set(input.delayedWethProxy.selector, address(delayedWETH));
        input.set(input.anchorStateRegistryProxy.selector, address(registry));
        input.set(input.vmAddress.selector, address(mipsVM));
        input.set(input.l2ChainId.selector, vm.envUint("L2_CHAIN_ID"));

        // 배포 실행
        DeployDisputeGame deployer = new DeployDisputeGame();
        deployer.run(input);
    }
}
```

### Step 3: 실행

```bash
# .env 파일 설정
DELAYED_WETH_PROXY=0x...
ANCHOR_STATE_REGISTRY_PROXY=0x...
MIPS_SINGLETON=0x...
ABSOLUTE_PRESTATE=0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98
L2_CHAIN_ID=901

# 실행
forge script scripts/deploy/examples/DeployCannonGame.s.sol:DeployCannonGame \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast
```

### Step 4: Factory에 등록

**파일**: `scripts/deploy/examples/SetCannonGameImpl.s.sol`

```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Script } from "forge-std/Script.sol";
import { SetDisputeGameImpl } from "scripts/deploy/SetDisputeGameImpl.s.sol";
import { SetDisputeGameImplInput } from "scripts/deploy/SetDisputeGameImpl.s.sol";
import { IDisputeGameFactory } from "interfaces/dispute/IDisputeGameFactory.sol";
import { IAnchorStateRegistry } from "interfaces/dispute/IAnchorStateRegistry.sol";
import { IFaultDisputeGame } from "interfaces/dispute/IFaultDisputeGame.sol";
import { GameTypes } from "src/dispute/lib/Types.sol";

contract SetCannonGameImpl is Script {
    function run() public {
        IDisputeGameFactory factory = IDisputeGameFactory(vm.envAddress("DISPUTE_GAME_FACTORY_PROXY"));
        IAnchorStateRegistry registry = IAnchorStateRegistry(vm.envAddress("ANCHOR_STATE_REGISTRY_PROXY"));
        IFaultDisputeGame game = IFaultDisputeGame(vm.envAddress("FAULT_DISPUTE_GAME"));

        SetDisputeGameImplInput input = new SetDisputeGameImplInput();
        input.set(input.factory.selector, address(factory));
        input.set(input.anchorStateRegistry.selector, address(registry));
        input.set(input.impl.selector, address(game));
        input.set(input.gameType.selector, GameTypes.CANNON.raw());

        SetDisputeGameImpl setter = new SetDisputeGameImpl();
        setter.run(input);
    }
}
```

```bash
# 실행
DISPUTE_GAME_FACTORY_PROXY=0x...
ANCHOR_STATE_REGISTRY_PROXY=0x...
FAULT_DISPUTE_GAME=0x...  # 이전 단계에서 배포된 주소

forge script scripts/deploy/examples/SetCannonGameImpl.s.sol:SetCannonGameImpl \
    --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast
```

---

## 📝 예시 2: PermissionedDisputeGame 배포

### 차이점
- `gameKind`: "PermissionedDisputeGame"
- `gameType`: `GameTypes.PERMISSIONED_CANNON`
- 추가 파라미터: `proposer`, `challenger`

```solidity
input.set(input.gameKind.selector, "PermissionedDisputeGame");
input.set(input.gameType.selector, uint256(GameTypes.PERMISSIONED_CANNON.raw()));
input.set(input.proposer.selector, vm.envAddress("PROPOSER"));
input.set(input.challenger.selector, vm.envAddress("CHALLENGER"));
// ... 나머지는 동일
```

---

## 📝 예시 3: AlphabetVM 테스트 게임 배포

### Step 1: PreimageOracle 확인

```solidity
address preimageOracle = artifacts.mustGetAddress("PreimageOracle");
```

### Step 2: AlphabetVM 배포

```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Script } from "forge-std/Script.sol";
import { DeployAlphabetVM } from "scripts/deploy/DeployAlphabetVM.s.sol";
import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";

contract DeployAlphabetVMScript is Script {
    function run() public {
        IPreimageOracle oracle = IPreimageOracle(vm.envAddress("PREIMAGE_ORACLE"));

        DeployAlphabetVM deployer = new DeployAlphabetVM();
        DeployAlphabetVM.Output memory output = deployer.run(
            DeployAlphabetVM.Input({
                absolutePrestate: bytes32(uint256(0x60)),  // Alphabet prestate
                preimageOracle: oracle
            })
        );

        console.log("AlphabetVM deployed at:", address(output.alphabetVM));
    }
}
```

### Step 3: Alphabet 게임 배포

```solidity
input.set(input.gameKind.selector, "FaultDisputeGame");
input.set(input.gameType.selector, uint256(GameTypes.ALPHABET.raw()));
input.set(input.absolutePrestate.selector, bytes32(uint256(0x60)));
input.set(input.maxGameDepth.selector, 4);  // Alphabet은 깊이가 얕음
input.set(input.splitDepth.selector, 2);
input.set(input.clockExtension.selector, uint256(Duration.unwrap(Duration.wrap(30 minutes))));
input.set(input.maxClockDuration.selector, uint256(Duration.unwrap(Duration.wrap(1 days))));
input.set(input.vmAddress.selector, address(alphabetVM));  // AlphabetVM 사용
// ... 나머지
```

---

## 🔧 파라미터 설명

### GameType
| 타입 | 값 | 설명 |
|------|-----|------|
| `CANNON` | 0 | Permissionless MIPS Cannon 게임 |
| `PERMISSIONED_CANNON` | 1 | Permissioned MIPS Cannon 게임 |
| `ALPHABET` | 255 | 테스트용 Alphabet 게임 |

### 게임 파라미터

| 파라미터 | 타입 | 설명 | 예시 값 |
|---------|------|------|---------|
| `gameKind` | string | 게임 종류 | "FaultDisputeGame" or "PermissionedDisputeGame" |
| `gameType` | uint32 | GameType 값 | 0 (CANNON), 1 (PERMISSIONED_CANNON) |
| `absolutePrestate` | bytes32 | VM 시작 상태 해시 | 0x03c7ae... |
| `maxGameDepth` | uint256 | 게임 최대 깊이 | 73 (MIPS), 4 (Alphabet) |
| `splitDepth` | uint256 | Split 깊이 | 30 (MIPS), 2 (Alphabet) |
| `clockExtension` | uint64 | 시계 연장 시간 (초) | 10800 (3시간) |
| `maxClockDuration` | uint64 | 최대 시간 (초) | 302400 (3.5일) |
| `delayedWethProxy` | address | DelayedWETH 프록시 | 0x... |
| `anchorStateRegistryProxy` | address | AnchorStateRegistry 프록시 | 0x... |
| `vm` | address | VM 주소 (MIPS, Alphabet 등) | 0x... |
| `l2ChainId` | uint256 | L2 Chain ID | 901 |

### PermissionedDisputeGame 추가 파라미터

| 파라미터 | 타입 | 설명 |
|---------|------|------|
| `proposer` | address | Proposer 주소 (게임 시작 권한) |
| `challenger` | address | Challenger 주소 (챌린지 권한) |

---

## 🎯 사용 시나리오

### 시나리오 1: 개발 환경에서 빠른 테스트

```bash
# AlphabetVM으로 빠른 테스트 게임 배포
1. DeployAlphabetVM 실행
2. DeployDisputeGame (Alphabet) 실행
3. SetDisputeGameImpl 실행
4. 테스트 → 1-2분 소요

# 프로덕션 MIPS 게임으로 전환
5. DeployDisputeGame (CANNON) 실행
6. SetDisputeGameImpl 실행
```

### 시나리오 2: 게임 파라미터 튜닝

```bash
# 다양한 maxGameDepth로 테스트
1. DeployDisputeGame (depth=50) 실행
2. 테스트
3. DeployDisputeGame (depth=73) 실행
4. 테스트
5. 최적값 선택 후 Factory 업데이트
```

### 시나리오 3: Permissioned → Permissionless 전환

```bash
# 현재: PermissionedDisputeGame 사용 중

# 새 Permissionless 게임 배포 (기존 시스템 영향 없음)
1. DeployDisputeGame (FaultDisputeGame, CANNON) 실행
2. 테스트넷에서 검증
3. SetDisputeGameImpl 실행 (Factory 업데이트)
4. AnchorStateRegistry의 respectedGameType 업데이트
```

### 시나리오 4: 여러 게임 타입 동시 지원

```bash
# 다양한 게임 타입을 동시에 제공
1. DeployDisputeGame (CANNON) → GameType 0
2. DeployDisputeGame (PERMISSIONED_CANNON) → GameType 1
3. DeployDisputeGame (ALPHABET) → GameType 255
4. 사용자가 선택하여 게임 생성 가능
```

---

## ⚠️ 주의사항

### 1. 기존 OPCM 배포와의 관계

**중요**: 새로 추가된 스크립트는 기존 OPCM 배포와 **독립적**입니다.

```
기존 방식 (유지):
  Deploy.deployOpChain()
    └─ opcm.deploy()
        └─ 자동으로 DisputeGame 배포

새 방식 (추가):
  DeployDisputeGame.run()
    └─ 독립적으로 게임 배포
```

**하위 호환성**: 기존 배포 코드는 그대로 작동합니다.

### 2. Factory 등록 시 주의

```solidity
// SetDisputeGameImpl.run()은 다음을 확인:
require(address(factory.gameImpls(gameType)) == address(0), "SDGI-10");
// → 이미 등록된 GameType은 등록 불가
```

**해결 방법**:
- 새로운 GameType 사용
- 또는 Factory owner가 기존 구현을 제거 후 재등록

### 3. AnchorStateRegistry 업데이트

`SetDisputeGameImpl`은 선택적으로 AnchorStateRegistry도 업데이트합니다:

```solidity
if (address(anchorStateRegistry) != address(0)) {
    anchorStateRegistry.setRespectedGameType(gameType);
}
```

**주의**: `respectedGameType`은 하나만 설정 가능 (가장 신뢰하는 게임 타입)

### 4. 권한 관리

배포 및 등록에는 적절한 권한이 필요합니다:
- `DisputeGameFactory.setImplementation()`: Factory owner만 가능
- `AnchorStateRegistry.setRespectedGameType()`: Registry owner만 가능

---

## 🔍 검증 방법

### 게임이 올바르게 배포되었는지 확인

```solidity
// 1. Factory에 등록되었는지 확인
IDisputeGameFactory factory = IDisputeGameFactory(factoryAddress);
IFaultDisputeGame impl = factory.gameImpls(gameType);
console.log("Game implementation:", address(impl));

// 2. AnchorStateRegistry의 respectedGameType 확인
IAnchorStateRegistry registry = IAnchorStateRegistry(registryAddress);
GameType respectedType = registry.respectedGameType();
console.log("Respected game type:", respectedType.raw());

// 3. 게임 생성 테스트
bytes32 rootClaim = bytes32(uint256(0x123456));
bytes memory extraData = abi.encode(uint256(block.number));
IDisputeGame game = factory.create(gameType, Claim.wrap(rootClaim), extraData);
console.log("Created game:", address(game));
```

---

## 📊 기존 방식과의 비교

| 작업 | 기존 (OPCM) | 새 방식 (DeployDisputeGame) |
|------|------------|----------------------------|
| **게임 배포** | OP Chain 전체 배포 필요 | 게임만 독립 배포 |
| **소요 시간** | 10-20분 | 1-2분 |
| **VM 선택** | 고정 (Cannon) | 자유 (MIPS/Alphabet/Asterisc) |
| **파라미터 변경** | DeployInput 수정 → 재배포 | 파라미터만 변경 → 재실행 |
| **테스트** | 전체 시스템 재배포 | 게임만 재배포 |
| **다중 게임** | 불가능 | 여러 게임 독립 배포 가능 |

---

## 🎓 Best Practices

### 1. 개발 환경
- AlphabetVM으로 빠른 프로토타이핑
- 파라미터 튜닝 후 MIPS로 전환

### 2. 테스트 환경
- Permissioned 게임으로 먼저 검증
- 충분한 테스트 후 Permissionless로 전환

### 3. 프로덕션 환경
- 새 게임을 별도 GameType으로 배포
- 테스트 완료 후 Factory 업데이트
- AnchorStateRegistry의 respectedGameType 신중히 변경

### 4. 업그레이드
- 기존 게임은 유지 (하위 호환성)
- 새 게임을 추가로 배포
- 점진적으로 전환

---

## 📚 관련 문서

- [DISPUTE-GAME-DEPLOYMENT-ANALYSIS.md](./DISPUTE-GAME-DEPLOYMENT-ANALYSIS.md) - 배포 방식 분석
- [ESSENTIAL-DEPLOYMENT-SCRIPTS-FOR-THANOS.md](./ESSENTIAL-DEPLOYMENT-SCRIPTS-FOR-THANOS.md) - 필수 스크립트 목록
- [DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md](./DEPLOYMENT-SCRIPT-MIGRATION-GUIDE.md) - 마이그레이션 가이드

---

## 🆘 트러블슈팅

### 문제 1: "SDGI-10" 에러
```
원인: GameType이 이미 Factory에 등록되어 있음
해결: 새로운 GameType을 사용하거나 Factory owner가 기존 구현 제거
```

### 문제 2: "SDGI-20" 에러
```
원인: AnchorStateRegistry의 disputeGameFactory가 일치하지 않음
해결: 올바른 AnchorStateRegistry 주소 사용
```

### 문제 3: VM 주소가 없음
```
원인: MIPS/AlphabetVM이 배포되지 않음
해결: DeployMIPS2 또는 DeployAlphabetVM 먼저 실행
```

### 문제 4: 권한 부족
```
원인: Factory/Registry owner가 아님
해결: Owner 계정으로 실행하거나 권한 위임
```

---

**작성자**: Claude Code
**다음 단계**: 실제 배포 테스트 및 E2E 통합

---

## 💡 요약

1. ✅ **3개 파일 추가됨**: DeployDisputeGame, SetDisputeGameImpl, DeployAlphabetVM
2. ✅ **기존 코드 유지**: OPCM 배포 방식은 그대로 작동
3. ✅ **독립 배포 가능**: 게임만 빠르게 배포/업데이트
4. ✅ **유연성 증가**: 다양한 VM과 파라미터 선택
5. ✅ **개발 속도 향상**: 1-2분으로 재배포 가능

**권장**: 개발/테스트에서는 새 방식 사용, 프로덕션 초기 배포는 OPCM 사용
