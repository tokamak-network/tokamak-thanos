# Asterisc 온체인 컨트랙트 파일 체크리스트

## 개요

Optimism의 Asterisc (GameType 2) 온체인 검증에 필요한 모든 컨트랙트 파일 목록 및 통합 체크리스트입니다.

**작성일**: 2025년 10월 21일
**목적**: Tokamak-Thanos에 누락된 파일 확인 및 통합

## 파일 비교 요약

### Optimism (✅ 완전)

| # | 파일 | 크기 | 라인 수 | 타입 | 상태 |
|---|------|------|---------|------|------|
| 1 | `src/vendor/asterisc/RISCV.sol` | 70K | 1,707 | 컨트랙트 | ✅ |
| 2 | `interfaces/vendor/asterisc/IRISCV.sol` | 530B | 14 | 인터페이스 | ✅ |
| 3 | `scripts/deploy/DeployAsterisc.s.sol` | 1.7K | 57 | 배포 스크립트 | ✅ |
| 4 | `test/opcm/DeployAsterisc.t.sol` | 1.5K | 44 | 테스트 | ✅ |

**총 4개 파일**, **1,822 라인**

### Tokamak-Thanos (❌ 없음)

| # | 파일 | 상태 |
|---|------|------|
| 1 | `src/vendor/asterisc/RISCV.sol` | ❌ 없음 |
| 2 | `interfaces/vendor/asterisc/IRISCV.sol` | ❌ 없음 |
| 3 | `scripts/deploy/DeployAsterisc.s.sol` | ❌ 없음 |
| 4 | `test/opcm/DeployAsterisc.t.sol` | ❌ 없음 |

**디렉토리 자체가 존재하지 않음**

## 상세 파일 목록

### 1. RISCV.sol (컨트랙트)

**경로**: `src/vendor/asterisc/RISCV.sol`
**크기**: 70K (1,707 라인)
**Solidity 버전**: 0.8.25

#### 내용
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
            // RISC-V 명령어 집합 구현
        }
    }
}
```

#### 의존성
- ✅ **IPreimageOracle**: `interfaces/cannon/IPreimageOracle.sol`
- ✅ **IBigStepper**: `interfaces/dispute/IBigStepper.sol`

#### 기능
- RISC-V RV64GC 명령어 온체인 실행
- Merkle proof 기반 메모리 검증
- Preimage oracle 통합
- Gas 최적화를 위한 Yul 구현

#### 체크리스트
- [ ] 파일 복사 완료
- [ ] import 경로 수정 (필요 시)
- [ ] IPreimageOracle 존재 확인
- [ ] IBigStepper 존재 확인
- [ ] 컴파일 성공
- [ ] 가스 사용량 프로파일링

---

### 2. IRISCV.sol (인터페이스)

**경로**: `interfaces/vendor/asterisc/IRISCV.sol`
**크기**: 530B (14 라인)
**Solidity 버전**: ^0.8.0

#### 내용
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { ISemver } from "interfaces/universal/ISemver.sol";
import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";

interface IRISCV is ISemver {
    function oracle() external view returns (IPreimageOracle);
    function step(bytes memory _stateData, bytes memory _proof, bytes32 _localContext)
        external returns (bytes32);
    function __constructor__(IPreimageOracle _oracle) external;
}
```

#### 의존성
- ✅ **ISemver**: `interfaces/universal/ISemver.sol`
- ✅ **IPreimageOracle**: `interfaces/cannon/IPreimageOracle.sol`

#### 기능
- RISCV 컨트랙트의 공식 인터페이스
- 타입 안전성 제공
- 프록시 패턴 지원 (`__constructor__`)

#### 체크리스트
- [ ] 파일 복사 완료
- [ ] import 경로 수정 (필요 시)
- [ ] ISemver 존재 확인
- [ ] IPreimageOracle 존재 확인
- [ ] 컴파일 성공

---

### 3. DeployAsterisc.s.sol (배포 스크립트)

**경로**: `scripts/deploy/DeployAsterisc.s.sol`
**크기**: 1.7K (57 라인)
**Solidity 버전**: 0.8.15

#### 내용
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

#### 의존성
- ✅ **Script**: `forge-std/Script.sol`
- ✅ **DeployUtils**: `scripts/libraries/DeployUtils.sol`
- ✅ **IPreimageOracle**: `interfaces/cannon/IPreimageOracle.sol`
- ✅ **IRISCV**: `interfaces/vendor/asterisc/IRISCV.sol`

#### 기능
- RISCV 컨트랙트 배포
- PreimageOracle 연결
- 입력 검증
- 출력 검증

#### 사용법
```bash
forge script scripts/deploy/DeployAsterisc.s.sol:DeployAsterisc \
  --rpc-url $RPC_URL \
  --private-key $DEPLOYER_KEY \
  --broadcast \
  --sig "run((address))" \
  "($PREIMAGE_ORACLE_ADDRESS)"
```

#### 체크리스트
- [ ] 파일 복사 완료
- [ ] import 경로 수정 (필요 시)
- [ ] DeployUtils 존재 확인
- [ ] forge-std 설치 확인
- [ ] 컴파일 성공
- [ ] 배포 테스트 (DevNet)
- [ ] 배포 성공 확인

---

### 4. DeployAsterisc.t.sol (테스트)

**경로**: `test/opcm/DeployAsterisc.t.sol`
**크기**: 1.5K (44 라인)
**Solidity 버전**: 0.8.15

#### 내용
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { DeployUtils } from "scripts/libraries/DeployUtils.sol";
import { IPreimageOracle } from "interfaces/cannon/IPreimageOracle.sol";
import { DeployAsterisc } from "scripts/deploy/DeployAsterisc.s.sol";

contract DeployAsterisc_Test is Test {
    DeployAsterisc deployAsterisc;
    IPreimageOracle defaultPreimageOracle = IPreimageOracle(makeAddr("preimageOracle"));

    function setUp() public {
        deployAsterisc = new DeployAsterisc();
    }

    function test_run_succeeds(DeployAsterisc.Input memory _input) public {
        vm.assume(address(_input.preimageOracle) != address(0));

        DeployAsterisc.Output memory output = deployAsterisc.run(_input);

        DeployUtils.assertValidContractAddress(address(output.asteriscSingleton));
        assertEq(address(output.asteriscSingleton.oracle()),
                 address(_input.preimageOracle), "100");
    }

    function test_run_nullInput_reverts() public {
        DeployAsterisc.Input memory input;
        input = defaultInput();
        input.preimageOracle = IPreimageOracle(address(0));

        vm.expectRevert("DeployAsterisc: preimageOracle not set");
        deployAsterisc.run(input);
    }

    function defaultInput() internal view returns (DeployAsterisc.Input memory input_) {
        input_ = DeployAsterisc.Input(defaultPreimageOracle);
    }
}
```

#### 의존성
- ✅ **Test**: `forge-std/Test.sol`
- ✅ **DeployUtils**: `scripts/libraries/DeployUtils.sol`
- ✅ **IPreimageOracle**: `interfaces/cannon/IPreimageOracle.sol`
- ✅ **DeployAsterisc**: `scripts/deploy/DeployAsterisc.s.sol`

#### 테스트 케이스
1. **test_run_succeeds**: 정상 배포 테스트 (fuzz testing)
2. **test_run_nullInput_reverts**: null input 검증

#### 실행 방법
```bash
forge test --match-contract DeployAsterisc_Test -vvv
```

#### 체크리스트
- [ ] 파일 복사 완료
- [ ] import 경로 수정 (필요 시)
- [ ] test/opcm 디렉토리 생성
- [ ] 컴파일 성공
- [ ] 테스트 실행 성공
- [ ] 모든 테스트 통과

---

## 의존성 파일 확인

Asterisc 컨트랙트가 의존하는 다른 파일들도 확인 필요:

### 필수 인터페이스

| 파일 | 경로 | 확인 필요 |
|------|------|----------|
| **IPreimageOracle** | `interfaces/cannon/IPreimageOracle.sol` | ✅ 필수 |
| **IBigStepper** | `interfaces/dispute/IBigStepper.sol` | ✅ 필수 |
| **ISemver** | `interfaces/universal/ISemver.sol` | ✅ 필수 |

### 필수 라이브러리

| 파일 | 경로 | 확인 필요 |
|------|------|----------|
| **DeployUtils** | `scripts/libraries/DeployUtils.sol` | ✅ 필수 |

### Foundry 도구

| 도구 | 확인 필요 |
|------|----------|
| **forge** | ✅ 필수 |
| **cast** | ✅ 필수 |
| **forge-std** | ✅ 필수 |

---

## 통합 체크리스트

### Phase 1: 파일 준비

#### 1.1 디렉토리 생성

⚠️ **중요**: Tokamak 커스텀 경로 사용!

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

mkdir -p src/vendor/asterisc
mkdir -p interfaces/vendor/asterisc
mkdir -p test/opcm
```

- [ ] `src/vendor/asterisc` 디렉토리 생성
- [ ] `interfaces/vendor/asterisc` 디렉토리 생성
- [ ] `test/opcm` 디렉토리 존재 확인

#### 1.2 파일 복사

⚠️ **중요**: 대상 경로가 `packages/tokamak/contracts-bedrock/`입니다!

```bash
# RISCV.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/vendor/asterisc/RISCV.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/src/vendor/asterisc/

# IRISCV.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/vendor/asterisc/IRISCV.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/interfaces/vendor/asterisc/

# DeployAsterisc.s.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAsterisc.s.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# DeployAsterisc.t.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/test/opcm/DeployAsterisc.t.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/test/opcm/
```

- [ ] RISCV.sol 복사 완료
- [ ] IRISCV.sol 복사 완료
- [ ] DeployAsterisc.s.sol 복사 완료
- [ ] DeployAsterisc.t.sol 복사 완료

#### 1.3 파일 확인

```bash
ls -lh src/vendor/asterisc/RISCV.sol
ls -lh interfaces/vendor/asterisc/IRISCV.sol
ls -lh scripts/deploy/DeployAsterisc.s.sol
ls -lh test/opcm/DeployAsterisc.t.sol
```

- [ ] 모든 파일 존재 확인
- [ ] 파일 크기 확인 (RISCV.sol ~70K)

---

### Phase 2: 의존성 확인

#### 2.1 필수 인터페이스 확인

```bash
# IPreimageOracle 확인
ls interfaces/cannon/IPreimageOracle.sol

# IBigStepper 확인
ls interfaces/dispute/IBigStepper.sol

# ISemver 확인
ls interfaces/universal/ISemver.sol
```

- [ ] IPreimageOracle.sol 존재
- [ ] IBigStepper.sol 존재
- [ ] ISemver.sol 존재

**없으면**: Optimism에서 복사 필요

#### 2.2 라이브러리 확인

```bash
# DeployUtils 확인
ls scripts/libraries/DeployUtils.sol
```

- [ ] DeployUtils.sol 존재

**없으면**: Optimism에서 복사 또는 대체 구현 필요

#### 2.3 Solidity 버전 확인

```bash
# foundry.toml 확인
cat foundry.toml | grep solc_version
```

- [ ] Solidity 0.8.25 지원 확인
- [ ] Solidity 0.8.15 지원 확인 (스크립트용)

**필요 시**: foundry.toml 수정

---

### Phase 3: 컴파일 및 테스트

#### 3.1 개별 컴파일

⚠️ **중요**: Tokamak 커스텀 경로로 이동!

```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# RISCV.sol 컴파일
forge build --contracts src/vendor/asterisc/RISCV.sol

# IRISCV.sol 컴파일
forge build --contracts interfaces/vendor/asterisc/IRISCV.sol

# DeployAsterisc.s.sol 컴파일
forge build --contracts scripts/deploy/DeployAsterisc.s.sol

# DeployAsterisc.t.sol 컴파일
forge build --contracts test/opcm/DeployAsterisc.t.sol
```

- [ ] RISCV.sol 컴파일 성공
- [ ] IRISCV.sol 컴파일 성공
- [ ] DeployAsterisc.s.sol 컴파일 성공
- [ ] DeployAsterisc.t.sol 컴파일 성공

#### 3.2 전체 빌드

```bash
forge build
```

- [ ] 전체 빌드 성공
- [ ] 경고 없음 (또는 경고 검토)

#### 3.3 테스트 실행

```bash
# DeployAsterisc 테스트
forge test --match-contract DeployAsterisc_Test -vvv

# 모든 테스트
forge test
```

- [ ] DeployAsterisc_Test 통과
- [ ] test_run_succeeds 통과
- [ ] test_run_nullInput_reverts 통과
- [ ] 기존 테스트 회귀 없음

---

### Phase 4: 배포

#### 4.1 DevNet 준비

```bash
# DevNet 시작
cd /Users/zena/tokamak-projects/tokamak-thanos
make devnet-up
```

- [ ] DevNet 실행 중

#### 4.2 PreimageOracle 주소 확인

```bash
# PreimageOracle 주소 조회 (예시)
cast call <DisputeGameFactory 주소> "preimageOracle()" --rpc-url http://localhost:8545
```

- [ ] PreimageOracle 주소 확보
- [ ] 주소 환경 변수 저장

#### 4.3 RISCV 배포

```bash
cd packages/tokamak/contracts-bedrock

# 환경 변수 설정
export PREIMAGE_ORACLE=<주소>
export RPC_URL=http://localhost:8545
export DEPLOYER_KEY=<개발용 키>

# 배포
forge script scripts/deploy/DeployAsterisc.s.sol:DeployAsterisc \
  --rpc-url $RPC_URL \
  --private-key $DEPLOYER_KEY \
  --broadcast \
  --sig "run((address))" \
  "($PREIMAGE_ORACLE)"
```

- [ ] 배포 트랜잭션 성공
- [ ] RISCV 컨트랙트 주소 확보
- [ ] 배포 로그 저장

#### 4.4 배포 검증

```bash
export RISCV_ADDRESS=<배포된 주소>

# version 확인
cast call $RISCV_ADDRESS "version()" --rpc-url $RPC_URL

# oracle 확인
cast call $RISCV_ADDRESS "oracle()" --rpc-url $RPC_URL
```

- [ ] version() 반환: "1.2.0-rc.1"
- [ ] oracle() 반환: PreimageOracle 주소 일치

---

### Phase 5: DisputeGameFactory 통합

#### 5.1 GameType 2 등록

```bash
export DISPUTE_GAME_FACTORY=<주소>
export ADMIN_KEY=<관리자 키>

# GameType 2 설정
cast send $DISPUTE_GAME_FACTORY \
  "setImplementation(uint32,address)" \
  2 \
  $RISCV_ADDRESS \
  --rpc-url $RPC_URL \
  --private-key $ADMIN_KEY
```

- [ ] setImplementation 트랜잭션 성공

#### 5.2 등록 확인

```bash
# GameType 2 구현 확인
cast call $DISPUTE_GAME_FACTORY "gameImpls(uint32)" 2 --rpc-url $RPC_URL
```

- [ ] gameImpls(2) 반환: RISCV 주소 일치

---

### Phase 6: E2E 테스트

#### 6.1 Dispute Game 생성

```bash
# op-challenger로 dispute game 생성
./bin/op-challenger \
  --trace-type asterisc \
  --game-type 2 \
  # ... 기타 설정
```

- [ ] Dispute game 생성 성공
- [ ] GameType 2 사용 확인

#### 6.2 Step 실행 테스트

```bash
# 온체인 step 실행
cast send $RISCV_ADDRESS \
  "step(bytes,bytes,bytes32)" \
  <stateData> \
  <proof> \
  <localContext> \
  --gas-limit 30000000 \
  --rpc-url $RPC_URL
```

- [ ] step() 트랜잭션 성공
- [ ] Gas 사용량 30M 이하
- [ ] 올바른 상태 해시 반환

#### 6.3 전체 Dispute 해결

```bash
# E2E 테스트 (Go)
cd /Users/zena/tokamak-projects/tokamak-thanos/op-e2e
go test -run TestAsteriscOnChain -v
```

- [ ] 전체 dispute 해결 성공
- [ ] 온체인 검증 성공
- [ ] 최종 결과 정확함

---

## 문제 해결 가이드

### 컴파일 에러

#### 에러: "Solidity version mismatch"
```
Solution: foundry.toml에서 solc_version 확인
solc_version = "0.8.25"
```

#### 에러: "Interface not found"
```
Solution: 필수 인터페이스 복사
- IPreimageOracle.sol
- IBigStepper.sol
- ISemver.sol
```

#### 에러: "DeployUtils not found"
```
Solution: Optimism에서 DeployUtils.sol 복사
cp <optimism>/scripts/libraries/DeployUtils.sol \
   <tokamak>/scripts/libraries/
```

### 배포 에러

#### 에러: "PreimageOracle not set"
```
Solution: PreimageOracle 주소 확인
cast call <DisputeGameFactory> "preimageOracle()"
```

#### 에러: "Out of gas"
```
Solution: Gas limit 증가
--gas-limit 30000000
```

#### 에러: "Create1 failed"
```
Solution: Deployer 계정 잔액 확인
cast balance $DEPLOYER --rpc-url $RPC_URL
```

### 테스트 에러

#### 에러: "Test failed: assertEq"
```
Solution: PreimageOracle 주소 불일치 확인
배포 시 사용한 주소와 테스트 주소 일치 필요
```

#### 에러: "VM execution reverted"
```
Solution: RISCV step() 로직 문제
- stateData 형식 확인
- proof 유효성 확인
- localContext 값 확인
```

---

## 참고 자료

### Optimism 원본 파일 위치

```
/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/
├── src/vendor/asterisc/RISCV.sol
├── interfaces/vendor/asterisc/IRISCV.sol
├── scripts/deploy/DeployAsterisc.s.sol
└── test/opcm/DeployAsterisc.t.sol
```

### Tokamak-Thanos 대상 위치 (⚠️ 정정됨!)

**중요**: Tokamak 커스텀 경로 사용

```
/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/
├── src/vendor/asterisc/RISCV.sol           (복사 필요)
├── interfaces/vendor/asterisc/IRISCV.sol   (복사 필요)
├── scripts/deploy/DeployAsterisc.s.sol     (복사 필요)
└── test/opcm/DeployAsterisc.t.sol          (복사 필요)
```

**현재 상태**:
```
packages/tokamak/contracts-bedrock/src/vendor/
├── AddressAliasHelper.sol                  # ✅ 있음
├── WNativeToken.sol                        # ✅ 있음
└── asterisc/                               # ❌ 없음 (생성 필요!)
```

### 관련 문서

- **온체인 검증 통합 가이드**: [onchain-contracts-integration-ko.md](./onchain-contracts-integration-ko.md)
- **GameType 2 통합 계획**: [gametype2-integration-plan-ko.md](./gametype2-integration-plan-ko.md)
- **Asterisc 비교**: [asterisc-comparison-optimism-vs-tokamak-ko.md](./asterisc-comparison-optimism-vs-tokamak-ko.md)

---

## 마스터 체크리스트

완료 현황을 추적하세요:

### 준비 단계
- [ ] Optimism 프로젝트 최신 상태 확인
- [ ] Tokamak-Thanos 프로젝트 클린 상태 확인
- [ ] forge, cast 설치 확인
- [ ] 작업 브랜치 생성

### 파일 복사 (4개)
- [ ] RISCV.sol 복사
- [ ] IRISCV.sol 복사
- [ ] DeployAsterisc.s.sol 복사
- [ ] DeployAsterisc.t.sol 복사

### 의존성 확인 (4개)
- [ ] IPreimageOracle.sol 존재
- [ ] IBigStepper.sol 존재
- [ ] ISemver.sol 존재
- [ ] DeployUtils.sol 존재

### 컴파일 (4개)
- [ ] RISCV.sol 컴파일 성공
- [ ] IRISCV.sol 컴파일 성공
- [ ] DeployAsterisc.s.sol 컴파일 성공
- [ ] DeployAsterisc.t.sol 컴파일 성공

### 테스트 (2개)
- [ ] test_run_succeeds 통과
- [ ] test_run_nullInput_reverts 통과

### 배포 (3개)
- [ ] DevNet 배포 성공
- [ ] DisputeGameFactory 등록 성공
- [ ] 배포 검증 완료

### E2E 테스트 (3개)
- [ ] Dispute game 생성 성공
- [ ] step() 실행 성공
- [ ] 전체 dispute 해결 성공

### 정리 (3개)
- [ ] 문서 업데이트
- [ ] 커밋 및 PR
- [ ] 배포 정보 기록

**총 진행률**: __ / 30 (0%)

완료되면 GameType 2 (Asterisc) 온체인 검증이 완전히 작동합니다!
