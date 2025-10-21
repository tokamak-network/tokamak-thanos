# 온체인 검증 컨트랙트 통합 가이드 (Asterisc RISC-V)

## 개요

Optimism의 Asterisc (GameType 2) 온체인 검증 컨트랙트를 Tokamak-Thanos에 통합하기 위한 가이드입니다.

**작성일**: 2025년 10월 21일
**대상**: Tokamak-Thanos 개발팀

## 현재 상태

### Optimism (✅ 완전 구현)

| 구성 요소 | 상태 | 위치 | 크기/버전 |
|----------|------|------|-----------|
| **RISCV.sol** | ✅ 있음 | `src/vendor/asterisc/RISCV.sol` | 1,707 라인 |
| **IRISCV.sol** | ✅ 있음 | `interfaces/vendor/asterisc/IRISCV.sol` | 인터페이스 |
| **DeployAsterisc.s.sol** | ✅ 있음 | `scripts/deploy/DeployAsterisc.s.sol` | 배포 스크립트 |
| **DeployAsterisc.t.sol** | ✅ 있음 | `test/opcm/DeployAsterisc.t.sol` | 테스트 |
| **버전** | v1.2.0-rc.1 | - | 최신 |

### Tokamak-Thanos (❌ 없음)

| 구성 요소 | 상태 | 비고 |
|----------|------|------|
| **RISCV.sol** | ❌ 없음 | vendor/asterisc 디렉토리 자체 없음 |
| **IRISCV.sol** | ❌ 없음 | 인터페이스 없음 |
| **DeployAsterisc.s.sol** | ❌ 없음 | 배포 스크립트 없음 |
| **DeployAsterisc.t.sol** | ❌ 없음 | 테스트 없음 |

## 파일 구조

### Optimism 구조

```
packages/contracts-bedrock/
├── src/
│   └── vendor/
│       └── asterisc/
│           └── RISCV.sol                    # 1,707 라인 (온체인 VM)
├── interfaces/
│   └── vendor/
│       └── asterisc/
│           └── IRISCV.sol                   # 인터페이스
├── scripts/
│   └── deploy/
│       └── DeployAsterisc.s.sol            # 배포 스크립트
└── test/
    └── opcm/
        └── DeployAsterisc.t.sol            # 테스트
```

### Tokamak-Thanos 현재 구조

**⚠️ 중요**: Tokamak-Thanos는 **2개의 contracts-bedrock**을 가지고 있습니다:
- `packages/contracts-bedrock/` - Optimism 원본 fork
- `packages/tokamak/contracts-bedrock/` - **Tokamak 커스텀 버전 (사용해야 함!)** ⭐

```
packages/tokamak/contracts-bedrock/
├── src/
│   ├── cannon/
│   │   ├── interfaces/
│   │   │   └── IPreimageOracle.sol        # ✅ 존재
│   │   └── MIPS.sol                        # ✅ 존재
│   ├── dispute/
│   │   └── interfaces/
│   │       └── IBigStepper.sol             # ✅ 존재
│   └── vendor/
│       ├── AddressAliasHelper.sol
│       └── # asterisc/ 디렉토리 없음 ❌
├── interfaces/
│   └── vendor/
│       # asterisc/ 디렉토리 없음 ❌
└── scripts/
    ├── libraries/
    │   # DeployUtils.sol 없음 ❌
    └── deploy/
        # DeployAsterisc.s.sol 없음 ❌
```

## RISCV.sol 상세 분석

### 기본 정보

```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";
import { IBigStepper } from "interfaces/dispute/IBigStepper.sol";

contract RISCV is IBigStepper {
    IPreimageOracle public oracle;
    string public constant version = "1.2.0-rc.1";

    constructor(IPreimageOracle _oracle) {
        oracle = _oracle;
    }

    function step(
        bytes calldata _stateData,
        bytes calldata _proof,
        bytes32 _localContext
    ) public returns (bytes32) {
        assembly {
            // 1,600+ 라인의 Yul 코드
            // RISC-V 명령어 실행 로직
        }
    }
}
```

### 주요 기능

#### 1. RISC-V Instruction Execution

**컨트랙트 역할**: 단일 RISC-V 명령어를 stateless하게 실행하고 검증

**입력**:
- `_stateData`: VM 상태 (PC, 레지스터, 메모리 등)
- `_proof`: Merkle proof (메모리 접근 증명)
- `_localContext`: 로컬 컨텍스트 (L1 head, L2 output root 등)

**출력**:
- `bytes32`: 실행 후 상태 해시

#### 2. Yul 구현

**최적화 이유**: Gas 효율성을 위해 저수준 Yul로 작성

**구현 범위**:
- RV64GC instruction set 지원
- 64비트 연산
- 메모리 Merkle proof 검증
- Preimage oracle 통합

### 코드 구조 (1,707 라인)

```
RISCV.sol (1,707 lines)
├── Helper Functions (Yul64)        # ~100 lines
│   ├── u64Mask(), u32Mask()
│   ├── signExtend64()
│   └── 산술/논리 연산
│
├── Memory Access Functions         # ~200 lines
│   ├── readMem()
│   ├── writeMem()
│   └── Merkle proof 검증
│
├── Register Operations             # ~100 lines
│   ├── readReg()
│   └── writeReg()
│
├── Preimage Oracle                 # ~150 lines
│   ├── readPreimageKey()
│   └── readPreimageValue()
│
└── Instruction Execution           # ~1,200 lines
    ├── ALU 명령어 (ADD, SUB, AND, OR, XOR, ...)
    ├── Shift 명령어 (SLL, SRL, SRA, ...)
    ├── Branch 명령어 (BEQ, BNE, BLT, BGE, ...)
    ├── Load/Store 명령어 (LB, LH, LW, LD, SB, SH, SW, SD)
    ├── Multiply/Divide (MUL, DIV, REM, ...)
    └── System 명령어 (ECALL, EBREAK, ...)
```

## IRISCV 인터페이스

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { ISemver } from "interfaces/universal/ISemver.sol";
import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";

interface IRISCV is ISemver {
    /// @notice Preimage oracle 컨트랙트 주소
    function oracle() external view returns (IPreimageOracle);

    /// @notice 단일 RISC-V 명령어 실행
    /// @param _stateData VM 상태 데이터
    /// @param _proof Merkle proof
    /// @param _localContext 로컬 컨텍스트
    /// @return 실행 후 상태 해시
    function step(
        bytes memory _stateData,
        bytes memory _proof,
        bytes32 _localContext
    ) external returns (bytes32);

    /// @notice 생성자 (프록시 패턴용)
    /// @param _oracle Preimage oracle 컨트랙트
    function __constructor__(IPreimageOracle _oracle) external;
}
```

## DeployAsterisc.s.sol 배포 스크립트

```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Script } from "forge-std/Script.sol";
import { DeployUtils } from "scripts/libraries/DeployUtils.sol";
import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";
import { IRISCV } from "interfaces/vendor/asterisc/IRISCV.sol";

contract DeployAsterisc is Script {
    struct Input {
        IPreimageOracle preimageOracle;
    }

    struct Output {
        IRISCV asteriscSingleton;
    }

    function run(Input memory _input) public returns (Output memory output_) {
        assertValidInput(_input);

        deployAsteriscSingleton(_input, output_);

        assertValidOutput(_input, output_);
    }

    function deployAsteriscSingleton(Input memory _input, Output memory _output) internal {
        vm.broadcast(msg.sender);
        IRISCV singleton = IRISCV(
            DeployUtils.create1({
                _name: "RISCV",
                _args: DeployUtils.encodeConstructor(
                    abi.encodeCall(IRISCV.__constructor__, (_input.preimageOracle))
                )
            })
        );

        vm.label(address(singleton), "AsteriscSingleton");
        _output.asteriscSingleton = singleton;
    }

    function assertValidInput(Input memory _input) internal pure {
        require(
            address(_input.preimageOracle) != address(0),
            "DeployAsterisc: preimageOracle not set"
        );
    }

    function assertValidOutput(Input memory _input, Output memory _output) internal view {
        DeployUtils.assertValidContractAddress(address(_output.asteriscSingleton));

        require(
            _output.asteriscSingleton.oracle() == _input.preimageOracle,
            "DeployAsterisc: preimageOracle does not match input"
        );
    }
}
```

### 배포 프로세스

1. **PreimageOracle 확인**: 기존 PreimageOracle 컨트랙트 주소 확인
2. **RISCV 배포**: CREATE1 opcode로 RISCV 컨트랙트 배포
3. **검증**: oracle() 함수로 올바르게 설정되었는지 확인

## 통합 계획

### Phase 0: foundry.toml 설정 (1시간) 🔴 선행 작업

#### 0.1 remapping 추가

**문제**: Optimism은 `interfaces/` remapping을 사용하지만 Tokamak에는 없음

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# foundry.toml 편집
```

**추가할 내용**:
```toml
remappings = [
  # ... 기존 내용 유지 ...
  'interfaces/=interfaces'  # ← 이 줄 추가
]
```

**검증**:
```bash
# remapping 확인
forge remappings | grep interfaces
```

---

### Phase 1: 파일 복사 및 디렉토리 구조 생성 (1일)

#### 1.1 디렉토리 생성

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# vendor/asterisc 디렉토리 생성
mkdir -p src/vendor/asterisc
mkdir -p interfaces/vendor/asterisc
mkdir -p scripts/libraries
```

#### 1.2 RISCV.sol 복사

```bash
# RISCV.sol 컨트랙트 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/vendor/asterisc/RISCV.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/src/vendor/asterisc/

# 파일 확인
ls -lh /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/src/vendor/asterisc/RISCV.sol
```

#### 1.3 IRISCV.sol 인터페이스 복사

```bash
# IRISCV.sol 인터페이스 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/vendor/asterisc/IRISCV.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/interfaces/vendor/asterisc/
```

#### 1.4 DeployUtils.sol 복사 (⭐ 중요!)

```bash
# DeployUtils.sol 라이브러리 복사 (배포 스크립트 의존성)
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/libraries/DeployUtils.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/libraries/

# 확인
ls -lh /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/libraries/DeployUtils.sol
```

#### 1.5 배포 스크립트 복사

```bash
# DeployAsterisc.s.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAsterisc.s.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/
```

#### 1.6 테스트 파일 복사 (선택)

```bash
# DeployAsterisc.t.sol 복사 (선택)
mkdir -p test/opcm
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/test/opcm/DeployAsterisc.t.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/test/opcm/
```

---

### Phase 2: 의존성 확인 (1시간)

#### 2.1 Import 경로 확인

**✅ Phase 0에서 remapping 추가로 import 경로 수정 불필요!**

**RISCV.sol 의존성**:
```solidity
import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";
import { IBigStepper } from "interfaces/dispute/IBigStepper.sol";
```

**실제 파일 위치 확인**:
```bash
# IPreimageOracle 존재 확인
ls /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/src/cannon/interfaces/IPreimageOracle.sol
# ✅ 존재함

# IBigStepper 존재 확인
ls /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/src/dispute/interfaces/IBigStepper.sol
# ✅ 존재함
```

**Remapping 동작 확인**:
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# interfaces/ remapping이 있는지 확인
forge remappings | grep "interfaces"
# 출력: interfaces/=interfaces

# 이로 인해 다음과 같이 해석됨:
# "interfaces/cannon/IPreimageOracle.sol"
# → "interfaces/cannon/IPreimageOracle.sol" (그대로)
```

**⚠️ 주의**: Tokamak의 실제 파일은 `src/cannon/interfaces/`에 있지만, `interfaces/` 디렉토리도 심볼릭 링크 또는 복사본으로 존재할 수 있습니다. 확인 필요!

---

### Phase 3: 컴파일 및 테스트 (2-3일)

#### 3.1 Solidity 버전 업그레이드 (0.8.15 → 0.8.25) ⭐

**목표**: Tokamak-Thanos 전체를 Optimism과 동일한 **0.8.25로 통일**

**현재 상태**:
- Optimism RISCV.sol: `pragma solidity 0.8.25;`
- Tokamak MIPS.sol: `pragma solidity 0.8.15;`

**전략**: 전체 프로젝트 업그레이드

**Step 1: Solidity 0.8.15 → 0.8.25 변경사항 확인**

주요 변경사항 (0.8.16 ~ 0.8.25):
- 0.8.19: `abi.encodeCall` 타입 체크 강화
- 0.8.20: Immutable 함수 개선
- 0.8.21: Contract size 최적화
- 0.8.22: 버그 수정
- 0.8.24: Yul optimizer 개선
- 0.8.25: 추가 버그 수정

**Step 2: 모든 컨트랙트 pragma 업데이트**

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# 모든 .sol 파일에서 0.8.15를 0.8.25로 변경
find src -name "*.sol" -exec sed -i '' 's/pragma solidity 0.8.15;/pragma solidity 0.8.25;/g' {} \;
find test -name "*.sol" -exec sed -i '' 's/pragma solidity 0.8.15;/pragma solidity 0.8.25;/g' {} \;
find scripts -name "*.sol" -exec sed -i '' 's/pragma solidity 0.8.15;/pragma solidity 0.8.25;/g' {} \;

# ^0.8.15 형태도 변경
find src -name "*.sol" -exec sed -i '' 's/pragma solidity \^0.8.15;/pragma solidity 0.8.25;/g' {} \;

# 변경 내용 확인
git diff src/cannon/MIPS.sol | head -20
```

**Step 3: foundry.toml 확인**

```bash
# Solidity 버전이 명시되어 있는지 확인
grep "solc" foundry.toml

# 명시되어 있다면 0.8.25로 변경
# [profile.default]
# solc = "0.8.25"
```

**Step 4: 전체 컴파일 및 검증**

```bash
# 전체 clean build
forge clean
forge build

# 경고 및 에러 확인
forge build 2>&1 | grep -i "error\|warning"
```

**Step 5: 전체 테스트 실행**

```bash
# 모든 테스트 실행
forge test

# 실패한 테스트 확인
forge test 2>&1 | grep -i "fail"
```

**장점**:
- ✅ Optimism과 완전 통일
- ✅ 최신 Solidity 기능 및 최적화 활용
- ✅ 향후 업스트림 동기화 간소화
- ✅ RISCV.sol 수정 불필요 (그대로 사용)

**주의사항**:
- ⚠️ 전체 테스트 재실행 필수
- ⚠️ 버전 변경으로 인한 동작 차이 검증
- ⚠️ 배포 전 철저한 감사 필요
- ⚠️ 기존 배포된 컨트랙트와 호환성 확인

#### 3.2 컴파일

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# RISCV.sol만 컴파일
forge build --contracts src/vendor/asterisc/RISCV.sol

# 전체 컴파일
forge build
```

**예상 이슈**:
- ✅ ~~Import 경로 문제~~ (remapping으로 해결)
- ⚠️ Solidity 버전 불일치 (위에서 해결)
- ⚠️ interfaces/ 디렉토리 누락 (확인 필요)

#### 3.2 Unit 테스트 작성

```solidity
// test/vendor/asterisc/RISCV.t.sol
pragma solidity 0.8.25;

import { Test } from "forge-std/Test.sol";
import { RISCV } from "src/vendor/asterisc/RISCV.sol";
import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";

contract RISCVTest is Test {
    RISCV public riscv;
    IPreimageOracle public oracle;

    function setUp() public {
        // PreimageOracle 모의 객체 생성
        oracle = IPreimageOracle(address(0x1234));
        riscv = new RISCV(oracle);
    }

    function test_constructor() public {
        assertEq(address(riscv.oracle()), address(oracle));
        assertEq(riscv.version(), "1.2.0-rc.1");
    }

    function test_step_simpleInstruction() public {
        // 간단한 RISC-V 명령어 테스트
        bytes memory stateData = hex"...";  // 초기 상태
        bytes memory proof = hex"...";      // Merkle proof
        bytes32 localContext = bytes32(0);

        bytes32 result = riscv.step(stateData, proof, localContext);

        // 결과 검증
        assertNotEq(result, bytes32(0));
    }
}
```

#### 3.3 배포 스크립트 테스트

```bash
# DeployAsterisc 테스트
forge test --match-contract DeployAsteriscTest -vvv
```

---

### Phase 4: 배포 (3-5일)

#### 4.1 로컬 DevNet 배포

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos

# DevNet 시작
make devnet-up

# PreimageOracle 주소 확인
PREIMAGE_ORACLE=$(cast call <DisputeGameFactory 주소> "preimageOracle()")

# RISCV 배포
forge script scripts/deploy/DeployAsterisc.s.sol:DeployAsterisc \
  --rpc-url http://localhost:8545 \
  --private-key <개발용 키> \
  --broadcast \
  --sig "run((address))" \
  "($PREIMAGE_ORACLE)"
```

#### 4.2 배포 검증

```bash
# RISCV 컨트랙트 주소 저장
RISCV_ADDRESS=<배포된 주소>

# oracle() 확인
cast call $RISCV_ADDRESS "oracle()" --rpc-url http://localhost:8545

# version() 확인
cast call $RISCV_ADDRESS "version()" --rpc-url http://localhost:8545
```

#### 4.3 DisputeGameFactory 연결

GameType 2 (Asterisc)에 RISCV 컨트랙트 등록:

```bash
# DisputeGameFactory.setImplementation() 호출
cast send <DisputeGameFactory 주소> \
  "setImplementation(uint32,address)" \
  2 \
  $RISCV_ADDRESS \
  --rpc-url http://localhost:8545 \
  --private-key <관리자 키>
```

#### 4.4 Testnet 배포

```bash
# Testnet RPC URL 설정
TESTNET_RPC="https://rpc.testnet.tokamak.network"

# Testnet에 배포
forge script scripts/deploy/DeployAsterisc.s.sol:DeployAsterisc \
  --rpc-url $TESTNET_RPC \
  --private-key $DEPLOYER_KEY \
  --broadcast \
  --verify \
  --etherscan-api-key $ETHERSCAN_API_KEY
```

---

### Phase 5: E2E 테스트 (3-5일)

#### 5.1 Dispute Game 생성

```bash
# op-challenger로 dispute game 생성 (GameType 2)
./bin/op-challenger \
  --trace-type asterisc \
  --game-type 2 \
  --l1-eth-rpc http://localhost:8545 \
  --rollup-rpc http://localhost:9545 \
  # ... 기타 설정
```

#### 5.2 Step 실행 테스트

```bash
# 온체인에서 step() 함수 호출
cast send $RISCV_ADDRESS \
  "step(bytes,bytes,bytes32)" \
  <stateData> \
  <proof> \
  <localContext> \
  --gas-limit 30000000
```

#### 5.3 Gas 사용량 측정

```bash
# Gas 사용량 프로파일링
forge test --match-contract RISCVTest --gas-report
```

**예상 Gas 사용량**:
- 간단한 ALU 명령어: ~50,000 gas
- Load/Store 명령어: ~100,000 gas
- Branch 명령어: ~60,000 gas
- Multiply/Divide: ~150,000 gas

#### 5.4 전체 Dispute 시뮬레이션

```bash
# E2E 테스트 실행
cd /Users/zena/tokamak-projects/tokamak-thanos/op-e2e

go test -run TestAsteriscOnChain -v
```

---

### Phase 6: 문서화 및 모니터링 (1-2일)

#### 6.1 배포 문서 작성

```markdown
# RISCV 컨트랙트 배포 정보

## Testnet
- **RISCV Address**: 0x...
- **PreimageOracle**: 0x...
- **DisputeGameFactory**: 0x...
- **GameType**: 2
- **Version**: 1.2.0-rc.1
- **Deployment Date**: 2025-10-22
- **Deployer**: 0x...

## Mainnet
- **Status**: Not deployed
- **Planned Date**: TBD
```

#### 6.2 모니터링 설정

```yaml
# monitoring/alerts/asterisc.yaml
alerts:
  - name: AsteriscStepGasUsage
    condition: gas_used > 25000000
    severity: warning
    description: RISCV step() gas usage is high

  - name: AsteriscStepFailure
    condition: tx_reverted == true
    severity: critical
    description: RISCV step() transaction reverted
```

#### 6.3 변경사항 문서화

```markdown
# CHANGELOG.md

## [Unreleased]

### Added
- Asterisc (RISC-V) on-chain verification contract (RISCV.sol)
- IRISCV interface
- DeployAsterisc deployment script
- RISCV unit tests

### Changed
- GameType 2 now supports on-chain verification

### Security
- RISCV.sol audited by [Auditor Name] on [Date]
```

---

## 기술적 세부사항

### RISCV.sol 핵심 로직

#### 1. State Format

```
State Data Format (362 bytes):
├── exitCode (1 byte)        # VM 종료 상태
├── exited (1 byte)           # 종료 여부
├── step (8 bytes)            # 현재 step 번호
├── PC (8 bytes)              # Program Counter
├── Registers (32 * 8 bytes)  # 32개 레지스터 (64-bit)
└── Memory Root (32 bytes)    # Merkle root
```

#### 2. Memory Proof Format

Sparse Merkle Tree로 메모리 증명:

```
Proof:
├── Leaf Value (8 bytes)      # 메모리 값
├── Leaf Index (8 bytes)      # 메모리 주소
└── Siblings (32 * 32 bytes)  # Merkle proof 경로
```

#### 3. Step Execution Flow

```
1. Decode State
   ↓
2. Fetch Instruction (from memory with proof)
   ↓
3. Decode Instruction
   ↓
4. Execute Instruction
   - Read registers
   - Perform operation
   - Write result
   ↓
5. Update PC
   ↓
6. Encode New State
   ↓
7. Return State Hash
```

### Gas 최적화 전략

#### 1. Yul 사용

**이유**: Solidity보다 15-30% gas 절감

```yul
// Yul 예시 (RISCV.sol에서)
function add64(a, b) -> out {
    out := and(add(a, b), u64Mask())
}
```

#### 2. Memory Layout 최적화

32바이트 단위로 정렬하여 SLOAD/SSTORE 최소화

#### 3. Merkle Proof 배치 처리

한 번에 여러 메모리 접근을 검증하여 중복 계산 제거

### 보안 고려사항

#### 1. Prestate 검증

```solidity
// RISCV.sol에서 prestate 검증
require(
    keccak256(stateData) == expectedPrestateHash,
    "Invalid prestate"
);
```

#### 2. Proof 검증

```solidity
// Merkle proof 검증
require(
    verifyMerkleProof(leaf, proof, memoryRoot),
    "Invalid memory proof"
);
```

#### 3. Gas Limit

```solidity
// Gas limit 확인
require(
    gasleft() > MIN_GAS_REQUIRED,
    "Insufficient gas"
);
```

## 위험 요소 및 대응

### 1. 경로 오류 (가장 흔한 실수)

**위험**: `packages/contracts-bedrock/` vs `packages/tokamak/contracts-bedrock/` 혼동

**대응**:
- ✅ 항상 `packages/tokamak/contracts-bedrock/` 사용
- ❌ `packages/contracts-bedrock/` 사용 금지 (Optimism 원본 fork)
- 모든 명령 실행 전 `pwd` 확인

### 2. 컴파일 실패

**위험**: Solidity 버전 불일치, remapping 누락, 의존성 누락

**대응**:
- ✅ foundry.toml에 `interfaces/=interfaces` remapping 추가 (Phase 0)
- ✅ RISCV.sol을 0.8.15로 다운그레이드 (Phase 3.1)
- ✅ DeployUtils.sol 복사 필수 (Phase 1.4)
- ⚠️ 컴파일 에러 시 remapping 재확인

### 2. Gas Limit 초과

**위험**: RISC-V 명령어 실행 시 30M gas limit 초과

**대응**:
- Gas profiling 실행
- 복잡한 명령어는 여러 step으로 분할
- L1 gas price 모니터링

### 3. 보안 취약점

**위험**: 온체인 검증 로직에 버그

**대응**:
- Optimism 감사 보고서 확인
- 추가 감사 실시
- Bug bounty 프로그램 운영

### 4. PreimageOracle 호환성

**위험**: Tokamak의 PreimageOracle이 Optimism과 다를 수 있음

**대응**:
- IPreimageOracle 인터페이스 확인
- 호환성 테스트 실행
- 필요 시 어댑터 작성

## 성공 기준

### 필수 요구사항

- ✅ RISCV.sol 컴파일 성공
- ✅ 모든 unit 테스트 통과
- ✅ DevNet 배포 성공
- ✅ GameType 2 dispute game 생성 및 해결 성공
- ✅ Gas 사용량이 예산 내 (< 30M per step)

### 선택 요구사항

- 감사 완료
- Mainnet 배포
- 모니터링 대시보드 구축
- 문서화 완료

## 타임라인

| Phase | 작업 | 예상 기간 | 변경 |
|-------|------|----------|------|
| **Phase 0** | foundry.toml 설정 | 1시간 | 🆕 추가 |
| Phase 1 | 파일 복사 (DeployUtils 포함) | 1일 | ✏️ 수정 |
| Phase 2 | 의존성 확인 | 1시간 | ✏️ 간소화 |
| Phase 3 | 컴파일 및 테스트 (버전 조정) | 2-3일 | ✏️ 수정 |
| Phase 4 | 배포 | 3-5일 | - |
| Phase 5 | E2E 테스트 | 3-5일 | - |
| Phase 6 | 문서화 및 모니터링 | 1-2일 | - |

**총 예상 기간**: 10-16일 (약 2-3주)
**단축된 시간**: 1-2일 (remapping으로 import 경로 수정 불필요)

## 참고 자료

### Optimism

- **RISCV.sol**: `/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/vendor/asterisc/RISCV.sol`
- **Asterisc Repo**: https://github.com/ethereum-optimism/asterisc
- **Specs**: https://specs.optimism.io/experimental/fault-proof/

### Tokamak-Thanos

- **contracts-bedrock**: `/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/` ⭐
- **관련 문서**:
  - [GameType 2 통합 계획](./gametype2-integration-plan-ko.md)
  - [Asterisc 비교](./asterisc-comparison-optimism-vs-tokamak-ko.md)

### Audits

- Optimism Fault Proof Audit (확인 필요)
- Asterisc Security Review (확인 필요)

## 다음 단계

1. ✅ **Phase 1 시작**: 디렉토리 생성 및 파일 복사
2. **의존성 확인**: IPreimageOracle, IBigStepper 존재 확인
3. **컴파일 테스트**: forge build 실행
4. **배포 준비**: PreimageOracle 주소 확인

온체인 검증 컨트랙트 통합을 시작하시겠습니까?
