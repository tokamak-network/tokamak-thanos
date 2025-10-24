# RISC-V와 Zero-Knowledge Proof: 현재와 미래

## 목차
1. [개요](#개요)
2. [Fraud Proof vs ZK Proof](#fraud-proof-vs-zk-proof)
3. [RISC-V 기반 zkVM 프로젝트](#risc-v-기반-zkvm-프로젝트)
4. [Optimistic Rollup → ZK Rollup 전환](#optimistic-rollup--zk-rollup-전환)
5. [OP Stack + ZK 통합 시나리오](#op-stack--zk-통합-시나리오)
6. [실제 통합 예시](#실제-통합-예시)
7. [성능 및 비용 비교](#성능-및-비용-비교)
8. [로드맵 및 미래 전망](#로드맵-및-미래-전망)
9. [FAQ](#faq)

---

## 개요

### 현재 상태

**tokamak-thanos의 Asterisc (RISC-V)**는 **Optimistic Rollup의 fraud proof** 방식을 사용합니다. 이는 ZK proof와는 다른 메커니즘입니다.

```
┌─────────────────────────────────────────────────────────────┐
│                현재: Optimistic Rollup (Fraud Proof)         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Asterisc (RISC-V VM)                                       │
│  ├─ 오프체인: RISC-V 프로그램 실행 (op-program-rv64)        │
│  ├─ 증명: 실행 trace 저장 (.json.gz)                       │
│  └─ 온체인: 단일 instruction 재실행 (RiscV.sol)             │
│                                                              │
│  특징:                                                       │
│  ✅ 빠른 증명 생성 (5-30분)                                 │
│  ✅ 저렴한 오프체인 비용                                     │
│  ❌ 온체인 검증 비용 높음 (~5M gas)                         │
│  ❌ 챌린지 기간 필요 (7일)                                  │
│  ❌ Dispute 발생 시에만 온체인 실행                          │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 미래 가능성

**RISC-V + Zero-Knowledge Proof**를 결합하면 **ZK Rollup**이 가능합니다.

```
┌─────────────────────────────────────────────────────────────┐
│              미래: ZK Rollup (Validity Proof)                │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  RISC-V zkVM (예: RISC Zero, SP1)                           │
│  ├─ 오프체인: RISC-V 프로그램 실행                          │
│  ├─ 증명: ZK proof 생성 (STARK/SNARK)                       │
│  └─ 온체인: ZK proof 검증 (~300K gas)                       │
│                                                              │
│  특징:                                                       │
│  ✅ 낮은 온체인 검증 비용 (~300K gas)                       │
│  ✅ 즉각적인 finality (챌린지 기간 없음)                    │
│  ✅ 수학적 보장 (암호학적 증명)                             │
│  ❌ 느린 증명 생성 (1-10시간)                               │
│  ❌ 높은 오프체인 비용 (강력한 하드웨어)                    │
│  ❌ 프로그램 크기 제한                                       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Fraud Proof vs ZK Proof

### 메커니즘 비교

| 항목 | Fraud Proof (Optimistic) | ZK Proof (Validity) |
|-----|--------------------------|---------------------|
| **검증 방식** | 챌린지 발생 시 온체인 재실행 | 항상 ZK proof 검증 |
| **신뢰 가정** | "1-of-N 정직" (1명만 정직하면 됨) | 수학적 보장 (신뢰 불필요) |
| **챌린지 기간** | 7일 (일반적) | 없음 (즉시 finality) |
| **온체인 비용** | 높음 (~5M gas/dispute) | 낮음 (~300K gas) |
| **오프체인 비용** | 낮음 (일반 서버) | 높음 (GPU/FPGA) |
| **증명 생성 시간** | 빠름 (분 단위) | 느림 (시간 단위) |
| **프로그램 제약** | 없음 (임의 크기) | 있음 (proof 크기 제한) |

### Fraud Proof 동작 방식 (현재)

```
1. Proposer: Output root 제안
   └─> "Block 1000의 state root는 0xabcd..."

2. 챌린지 기간 대기 (7일)
   ├─ 아무도 이의 제기 안 함 → 확정
   └─ Challenger 발견 → Dispute game 시작

3. Dispute Game (Bisection)
   ├─ Depth 0-29: L2 output 검증
   └─ Depth 30-73: VM 실행 검증
       ├─ Asterisc 실행 (오프체인)
       ├─ Trace 생성
       └─ Claim 제출

4. Step Execution (온체인)
   ├─ RiscV.sol에서 단일 instruction 실행
   ├─ Prestate → Poststate 검증
   └─ 올바른 주장자 승리

5. Resolution
   └─> 7일 + dispute 시간 후 finality
```

### ZK Proof 동작 방식 (미래)

```
1. Sequencer: 트랜잭션 실행
   └─> State transition: S1 → S2

2. Prover: ZK proof 생성
   ├─ RISC-V 프로그램 실행 (op-program-rv64)
   ├─ Execution trace 생성
   └─ ZK proof 생성 (STARK/SNARK)
       π = Prove(S1, S2, transitions)

3. Verifier (온체인)
   ├─ Verify(π, S1, S2) → true/false
   └─ 수학적으로 보장됨

4. Finality
   └─> 즉시 확정 (챌린지 기간 없음)
```

### 트레이드오프

```
Fraud Proof:
  장점: 빠른 증명, 저렴한 운영, 큰 프로그램 지원
  단점: 긴 finality, 높은 온체인 비용 (dispute 발생 시)

ZK Proof:
  장점: 즉각적 finality, 낮은 온체인 비용, 수학적 보장
  단점: 느린 증명, 비싼 하드웨어, 프로그램 크기 제한
```

---

## RISC-V 기반 zkVM 프로젝트

### 1. RISC Zero

**개요**:
- 가장 성숙한 RISC-V zkVM
- zkSTARK 기반 증명 시스템
- Rust로 작성된 프로그램을 ZK proof로 실행

**아키텍처**:

```
┌────────────────────────────────────────────────────────┐
│                    RISC Zero zkVM                      │
├────────────────────────────────────────────────────────┤
│                                                         │
│  Guest Program (RISC-V)                                │
│  ├─ Rust 코드 작성                                     │
│  ├─ risc0-build로 컴파일                               │
│  └─ RISC-V ELF 바이너리 생성                           │
│                                                         │
│  Prover (Host)                                         │
│  ├─ RISC-V VM 시뮬레이션                               │
│  ├─ Execution trace 생성                               │
│  └─ zkSTARK proof 생성                                 │
│                                                         │
│  Verifier (온체인 또는 오프체인)                        │
│  ├─ Proof 검증                                         │
│  └─ Public inputs/outputs 확인                         │
│                                                         │
└────────────────────────────────────────────────────────┘
```

**코드 예시**:

```rust
// Guest 프로그램 (RISC-V)
use risc0_zkvm::guest::env;

pub fn main() {
    // Input 읽기
    let a: u64 = env::read();
    let b: u64 = env::read();

    // 연산
    let sum = a + b;

    // Output 커밋 (공개)
    env::commit(&sum);
}

// Host 프로그램
use risc0_zkvm::{default_prover, ExecutorEnv};

fn main() {
    // Executor 생성
    let env = ExecutorEnv::builder()
        .write(&10u64).unwrap()
        .write(&20u64).unwrap()
        .build().unwrap();

    // Proof 생성
    let prover = default_prover();
    let receipt = prover.prove(env, GUEST_ELF).unwrap();

    // 검증
    receipt.verify(GUEST_ID).unwrap();

    // Public output 읽기
    let sum: u64 = receipt.journal.decode().unwrap();
    println!("Sum: {}", sum); // 30
}
```

**성능**:
- 증명 생성: ~1-10분 (프로그램 크기에 따라)
- Proof 크기: ~100KB - 1MB
- 검증 시간: ~수 초 (오프체인), ~300K gas (온체인)

**Bonsai 네트워크**:
```
오프체인 증명 생성 서비스
├─ API를 통해 증명 요청
├─ 분산 prover 네트워크
└─ 결과 반환
```

---

### 2. SP1 (Succinct)

**개요**:
- Succinct Labs의 RISC-V zkVM
- Plonky3 증명 시스템 사용
- RISC Zero보다 10-100배 빠른 증명 생성

**특징**:
- **빠른 성능**: 최첨단 증명 알고리즘
- **모듈화**: 다양한 증명 시스템 지원 (STARK, SNARK, Groth16)
- **OP Stack 통합**: Optimism과 협력 중

**아키텍처**:

```
┌────────────────────────────────────────────────────────┐
│                      SP1 zkVM                          │
├────────────────────────────────────────────────────────┤
│                                                         │
│  Core VM                                               │
│  ├─ RISC-V RV32IM 지원                                 │
│  ├─ Execution trace 생성                               │
│  └─ Constraint 생성                                    │
│                                                         │
│  Prover                                                │
│  ├─ Plonky3 backend                                    │
│  ├─ STARK → SNARK 변환                                 │
│  └─ Groth16 최종 proof                                 │
│                                                         │
│  Verifier                                              │
│  ├─ EVM-compatible (Solidity)                          │
│  └─ ~300K gas                                          │
│                                                         │
└────────────────────────────────────────────────────────┘
```

**코드 예시**:

```rust
// Guest 프로그램
#![no_main]
sp1_zkvm::entrypoint!(main);

pub fn main() {
    let a = sp1_zkvm::io::read::<u64>();
    let b = sp1_zkvm::io::read::<u64>();

    let sum = a + b;

    sp1_zkvm::io::commit(&sum);
}

// Host 프로그램
use sp1_sdk::{ProverClient, SP1Stdin};

fn main() {
    let client = ProverClient::new();

    let mut stdin = SP1Stdin::new();
    stdin.write(&10u64);
    stdin.write(&20u64);

    // Proof 생성 (빠름!)
    let (pk, vk) = client.setup(ELF);
    let proof = client.prove(&pk, stdin).unwrap();

    // 검증
    client.verify(&proof, &vk).unwrap();
}
```

**성능 비교**:

| 프로그램 크기 | RISC Zero | SP1 |
|-------------|-----------|-----|
| 10K instructions | ~30초 | ~3초 |
| 100K instructions | ~5분 | ~30초 |
| 1M instructions | ~1시간 | ~5분 |

**OP Stack 통합 (OP Succinct)**:

```bash
# Optimism의 fault proof를 SP1으로 대체
op-program (RISC-V) → SP1 prover → ZK proof → 온체인 검증
```

---

### 3. Valida

**개요**:
- Lita Foundation의 RISC-V zkVM
- STARK 기반
- 단순하고 이해하기 쉬운 구조

**특징**:
- 교육용으로도 적합
- 투명한 구현
- 활발한 개발 중

---

### 4. Jolt (a16z crypto)

**개요**:
- Lookup-based zkVM
- Lasso/Jolt 기법 사용
- 매우 빠른 prover 성능

**특징**:
- 전통적 STARK보다 10-100배 빠름
- Lookup table 기반 최적화
- 연구 단계

---

## Optimistic Rollup → ZK Rollup 전환

### 하이브리드 접근

**단계적 전환**:

```
Phase 1: Optimistic Rollup (현재)
  └─> Fraud proof (Cannon/Asterisc)

Phase 2: 선택적 ZK
  ├─> 일반 트랜잭션: Fraud proof
  └─> 고가치 트랜잭션: ZK proof

Phase 3: 하이브리드
  ├─> Fraud proof (백업)
  └─> ZK proof (메인)

Phase 4: 완전 ZK Rollup
  └─> ZK proof만 사용
```

### OP Stack의 ZK 통합 방향

**Optimism의 접근**:

```
┌────────────────────────────────────────────────────────┐
│              OP Stack Modular Architecture             │
├────────────────────────────────────────────────────────┤
│                                                         │
│  Execution Layer (op-geth)                             │
│  ├─ EVM 트랜잭션 실행                                  │
│  └─ State transition                                   │
│                                                         │
│  Consensus Layer (op-node)                             │
│  ├─ Rollup 로직                                        │
│  └─ Derivation pipeline                                │
│                                                         │
│  Settlement Layer (L1)                                 │
│  ├─ [현재] Fault proof (Cannon/Asterisc)               │
│  └─ [미래] ZK proof (SP1/RISC Zero)                    │
│                                                         │
│  Proof Generation                                      │
│  ├─ op-program (RISC-V)                                │
│  ├─ [선택] Fraud proof trace                           │
│  └─ [선택] ZK proof                                    │
│                                                         │
└────────────────────────────────────────────────────────┘
```

**설정 예시**:

```bash
# 현재
--trace-type cannon

# 미래 (ZK 추가)
--trace-type cannon,sp1

# 완전 ZK
--trace-type sp1
```

---

## OP Stack + ZK 통합 시나리오

### 시나리오 1: Fraud Proof 백업

```
정상 동작: ZK proof 생성 및 제출
  ├─ op-program-rv64 실행
  ├─> SP1 prover로 ZK proof 생성
  └─> 온체인 검증 (~300K gas)

예외 상황: ZK prover 실패 시
  ├─> Fraud proof로 fallback
  ├─> Dispute game 시작
  └─> 기존 메커니즘 사용
```

**구현**:

```go
// game/fault/register.go
func RegisterGameTypes(...) {
    // ZK 지원 게임 타입
    if cfg.TraceTypeEnabled(config.TraceTypeSP1) {
        registerSP1ZK(faultTypes.SP1GameType, registry, ...)
    }

    // Fraud proof 백업
    if cfg.TraceTypeEnabled(config.TraceTypeCannon) {
        registerCannon(faultTypes.CannonGameType, registry, ...)
    }
}

// Hybrid TraceAccessor
type HybridTraceAccessor struct {
    zkProver   ZKProver      // SP1 또는 RISC Zero
    fallback   TraceAccessor // Cannon 또는 Asterisc
}

func (h *HybridTraceAccessor) Get(ctx, pos Position) (common.Hash, error) {
    // 1. ZK proof 시도
    proof, err := h.zkProver.GenerateProof(ctx, pos)
    if err == nil {
        return proof.PublicOutput, nil
    }

    // 2. 실패 시 fraud proof
    log.Warn("ZK proof failed, using fraud proof", "error", err)
    return h.fallback.Get(ctx, pos)
}
```

---

### 시나리오 2: 선택적 ZK (고가치 트랜잭션)

```
트랜잭션 분류:
├─ 일반 트랜잭션 (< 1 ETH)
│  └─> Fraud proof (빠르고 저렴)
│
└─ 고가치 트랜잭션 (>= 1 ETH)
   └─> ZK proof (즉각 finality)
```

**온체인 로직**:

```solidity
// FaultDisputeGame.sol
contract FaultDisputeGame {
    enum ProofType {
        FRAUD,
        ZK_STARK,
        ZK_SNARK
    }

    struct Claim {
        bytes32 value;
        Position position;
        ProofType proofType;  // 증명 타입
    }

    function resolve() external {
        if (rootClaim.proofType == ProofType.ZK_STARK) {
            // ZK proof 검증
            require(zkVerifier.verify(rootClaim.proof), "Invalid ZK proof");
            // 즉시 finality
            status = GameStatus.DEFENDER_WON;
        } else {
            // Fraud proof 로직 (기존)
            // 챌린지 기간 대기
            // ...
        }
    }
}
```

---

### 시나리오 3: Fast Finality (선택적 빠른 확정)

```
사용자 선택:
├─ 일반 모드 (7일 대기)
│  └─> 저렴함, Fraud proof
│
└─ Fast 모드 (즉시 확정)
   ├─> 비쌈 (ZK proof 비용 지불)
   └─> ZK proof 제출
```

**UI/UX**:

```typescript
// Frontend
const withdrawOptions = [
  {
    mode: 'standard',
    delay: '7 days',
    fee: '0.001 ETH',
    proofType: 'fraud'
  },
  {
    mode: 'fast',
    delay: '~30 minutes', // ZK proof 생성 시간
    fee: '0.05 ETH',      // ZK proof 비용
    proofType: 'zk'
  }
];
```

---

## 실제 통합 예시

### SP1을 사용한 ZK Fault Proof

#### 1. op-program을 SP1 guest로 실행

```rust
// sp1-guest/src/main.rs
#![no_main]
sp1_zkvm::entrypoint!(main);

use op_program::run_fault_proof;

pub fn main() {
    // SP1 input에서 파라미터 읽기
    let l1_head: [u8; 32] = sp1_zkvm::io::read();
    let l2_head: [u8; 32] = sp1_zkvm::io::read();
    let l2_claim: [u8; 32] = sp1_zkvm::io::read();
    let l2_block_number: u64 = sp1_zkvm::io::read();

    // op-program 로직 실행 (RISC-V)
    let result = run_fault_proof(
        l1_head,
        l2_head,
        l2_claim,
        l2_block_number,
    );

    // 결과 커밋 (공개)
    sp1_zkvm::io::commit(&result.output_root);
    sp1_zkvm::io::commit(&result.is_valid);
}
```

#### 2. Prover 실행

```rust
// sp1-host/src/main.rs
use sp1_sdk::{ProverClient, SP1Stdin};

fn main() {
    let client = ProverClient::new();

    // Input 준비
    let mut stdin = SP1Stdin::new();
    stdin.write(&l1_head);
    stdin.write(&l2_head);
    stdin.write(&l2_claim);
    stdin.write(&l2_block_number);

    // ZK proof 생성
    println!("Generating ZK proof...");
    let (pk, vk) = client.setup(GUEST_ELF);
    let proof = client.prove(&pk, stdin).expect("Proof generation failed");

    println!("Proof generated!");
    println!("Proof size: {} bytes", proof.bytes().len());

    // 검증 (로컬)
    client.verify(&proof, &vk).expect("Verification failed");

    // 온체인 제출용 proof 준비
    let groth16_proof = client.groth16_prove(&proof).expect("Groth16 conversion failed");

    // 저장
    std::fs::write("proof.bin", groth16_proof.bytes()).expect("Failed to write proof");
}
```

#### 3. 온체인 검증

```solidity
// contracts/zkFaultProof.sol
pragma solidity ^0.8.15;

import {SP1Verifier} from "@sp1-contracts/SP1Verifier.sol";

contract ZKFaultProof {
    SP1Verifier public immutable VERIFIER;
    bytes32 public immutable PROGRAM_VKEY;

    struct ZKProof {
        bytes proof;
        bytes publicInputs;
    }

    function verifyFaultProof(
        ZKProof calldata zkProof,
        bytes32 l1Head,
        bytes32 l2Head,
        bytes32 l2Claim,
        uint256 l2BlockNumber
    ) external view returns (bool) {
        // Public inputs 인코딩
        bytes memory publicInputs = abi.encode(
            l1Head,
            l2Head,
            l2Claim,
            l2BlockNumber
        );

        // ZK proof 검증
        require(
            VERIFIER.verify(
                zkProof.proof,
                publicInputs,
                PROGRAM_VKEY
            ),
            "Invalid ZK proof"
        );

        // Public output 디코딩
        (bytes32 outputRoot, bool isValid) = abi.decode(
            zkProof.publicInputs,
            (bytes32, bool)
        );

        return isValid;
    }
}
```

#### 4. op-challenger 통합

```go
// op-challenger/game/fault/trace/sp1/provider.go
package sp1

import (
    "context"
    "github.com/ethereum/go-ethereum/common"
    "github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/types"
)

type SP1TraceProvider struct {
    logger         log.Logger
    proverClient   *SP1Client
    gameDepth      types.Depth
    inputs         utils.LocalGameInputs
}

func (p *SP1TraceProvider) Get(ctx context.Context, pos types.Position) (common.Hash, error) {
    traceIndex := pos.TraceIndex(p.gameDepth)

    // ZK proof 생성 요청
    proof, err := p.proverClient.GenerateProof(ctx, ProofRequest{
        L1Head:         p.inputs.L1Head,
        L2Head:         p.inputs.L2Head,
        L2Claim:        p.inputs.L2Claim,
        L2BlockNumber:  p.inputs.L2BlockNumber,
        TraceIndex:     traceIndex.Uint64(),
    })
    if err != nil {
        return common.Hash{}, fmt.Errorf("failed to generate ZK proof: %w", err)
    }

    // Public output 반환
    return proof.OutputRoot, nil
}

func (p *SP1TraceProvider) GetStepData(ctx context.Context, pos types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
    // ZK proof에서는 step data가 단순함
    proof, err := p.proverClient.GetProof(ctx, pos)
    if err != nil {
        return nil, nil, nil, err
    }

    return proof.ProofBytes, proof.PublicInputs, nil, nil
}
```

#### 5. 설정 및 실행

```bash
# op-challenger 실행 (SP1 ZK mode)
./op-challenger \
  --trace-type sp1 \
  --sp1-prover-url http://localhost:3000 \  # SP1 prover 서비스
  --sp1-program-vkey 0x...                \  # Program verification key
  --l1-eth-rpc http://localhost:8545 \
  --rollup-rpc http://localhost:9546 \
  --game-factory-address $DGF_ADDRESS \
  --datadir /data/zk-challenger \
  --private-key $PRIVATE_KEY
```

---

## 성능 및 비용 비교

### 증명 생성 시간

| 프로그램 크기 | Fraud Proof (Asterisc) | SP1 (ZK) | RISC Zero (ZK) |
|-------------|----------------------|----------|----------------|
| 1M instructions | 5분 | 5분 | 30분 |
| 10M instructions | 30분 | 30분 | 5시간 |
| 100M instructions | 5시간 | 5시간 | 50시간 |

*op-program은 보통 10M-100M instructions*

### 온체인 검증 비용

| 증명 타입 | Gas 비용 | USD (gas=20 gwei, ETH=$2000) |
|---------|---------|------------------------------|
| Fraud Proof (Step) | ~5,000,000 gas | ~$200 |
| ZK STARK | ~1,000,000 gas | ~$40 |
| ZK SNARK (Groth16) | ~300,000 gas | ~$12 |

### Finality 시간

| 방식 | Finality |
|-----|----------|
| Fraud Proof | 7일 (챌린지 기간) |
| ZK Proof | ~30분 (증명 생성 시간) |

### 하드웨어 요구사항

#### Fraud Proof (현재)

```
CPU: 4 cores
RAM: 8GB
Disk: 100GB SSD
Network: 100 Mbps

비용: ~$50/월 (클라우드)
```

#### ZK Proof (미래)

```
CPU: 32+ cores (또는 GPU)
RAM: 128GB+
Disk: 1TB NVMe SSD
Network: 1 Gbps

비용: ~$500-2000/월 (클라우드)
또는 GPU 서버: ~$1000-5000/월
```

### 총 비용 비교 (월간)

| 항목 | Fraud Proof | ZK Proof |
|-----|-------------|----------|
| 서버 비용 | $50 | $1000 |
| 온체인 비용 (100 disputes/월) | $20,000 | $1,200 |
| **총합** | **$20,050** | **$2,200** |

**결론**:
- Dispute가 자주 발생하면 ZK가 저렴
- Dispute가 거의 없으면 Fraud proof가 저렴

---

## 로드맵 및 미래 전망

### 단계별 발전

```
2024 Q4 - 2025 Q1 (현재)
├─ Optimism Fault Proof 메인넷 출시
├─ Cannon (MIPS) 프로덕션
└─ Asterisc (RISC-V) 테스트

2025 Q2 - Q3
├─ SP1 통합 연구 (OP Succinct)
├─ RISC Zero Boundless 프로토타입
└─ 하이브리드 접근 테스트

2025 Q4 - 2026 Q1
├─ 선택적 ZK proof 도입
├─ Fast finality 옵션
└─ 비용 최적화

2026+
├─ 완전 ZK Rollup 전환 (선택적)
├─ zkEVM 네이티브 지원
└─ 범용 ZK 인프라
```

### 기술 발전 예상

#### Prover 성능 향상

```
현재 (2025):
  10M instructions → 30분 (SP1)

2026:
  10M instructions → 5분 (하드웨어 가속)

2027:
  10M instructions → 1분 (ASIC/FPGA)
```

#### 비용 절감

```
현재:
  ZK proof 생성: $10-100

2026:
  분산 prover 네트워크: $1-10

2027:
  대량 생산 ASIC: $0.1-1
```

### Optimism의 전략

**모듈화 (Modular Approach)**:

```
OP Stack = Execution + Consensus + Settlement

Settlement Layer (교체 가능):
├─ Fraud Proof (Cannon/Asterisc)
├─ ZK Proof (SP1/RISC Zero)
├─ Hybrid (둘 다)
└─ 커스텀 (네트워크 선택)
```

**예시**:

```bash
# 네트워크 A: 빠른 finality 중요 → ZK
op-node --settlement-type sp1

# 네트워크 B: 비용 최소화 → Fraud proof
op-node --settlement-type cannon

# 네트워크 C: 하이브리드
op-node --settlement-type hybrid --primary sp1 --fallback cannon
```

---

## FAQ

### Q1: tokamak-thanos는 언제 ZK를 지원하나요?

**A**: 현재는 지원 계획이 없지만, OP Stack의 로드맵을 따를 것으로 예상됩니다.

**예상 일정**:
- 2025 Q2: 연구 및 프로토타입
- 2025 Q4: 테스트넷 실험
- 2026+: 선택적 ZK 도입

**준비 사항**:
```bash
# 코드 구조는 이미 모듈화되어 있음
git grep "TraceType" # 새 타입 추가 용이
git grep "TraceAccessor" # 인터페이스 통일
```

---

### Q2: Fraud proof와 ZK proof를 동시에 사용할 수 있나요?

**A**: 네, 가능합니다!

**시나리오**:
1. **백업**: ZK 실패 시 fraud proof로 fallback
2. **검증**: 두 방식으로 동시 검증 (보안 강화)
3. **선택**: 트랜잭션마다 다른 방식 사용

**구현**:
```go
type HybridProofSystem struct {
    zkProver     ZKProver
    fraudProver  FraudProver
    strategy     ProofStrategy
}

func (h *HybridProofSystem) GenerateProof(ctx, claim) (Proof, error) {
    switch h.strategy {
    case StrategyZKFirst:
        proof, err := h.zkProver.Prove(ctx, claim)
        if err != nil {
            return h.fraudProver.Prove(ctx, claim) // fallback
        }
        return proof, nil

    case StrategyBoth:
        zkProof, _ := h.zkProver.Prove(ctx, claim)
        fraudProof, _ := h.fraudProver.Prove(ctx, claim)
        return CombinedProof{zkProof, fraudProof}, nil
    }
}
```

---

### Q3: ZK proof가 fraud proof를 완전히 대체할까요?

**A**: 완전 대체보다는 **공존**할 가능성이 높습니다.

**이유**:
1. **프로그램 크기 제한**: ZK는 큰 프로그램에 불리
2. **비용**: 작은 네트워크는 fraud proof가 경제적
3. **유연성**: 다양한 use case 지원
4. **백업**: 한 시스템 장애 시 대안

**미래 모습**:
```
대형 네트워크 (예: Optimism Mainnet):
  └─> ZK proof (빠른 finality, 많은 사용자)

소형 네트워크 (예: 게임 전용 롤업):
  └─> Fraud proof (저렴, 충분한 성능)

하이브리드 네트워크:
  ├─ 일반 트랜잭션: Fraud proof
  └─- 고가치 트랜잭션: ZK proof
```

---

### Q4: RISC-V zkVM을 직접 개발할 수 있나요?

**A**: 가능하지만 매우 어렵습니다.

**필요 지식**:
- RISC-V ISA 이해
- Zero-knowledge 암호학 (STARK/SNARK)
- Constraint system 설계
- 최적화 기법

**대안**:
- 기존 zkVM 사용 (RISC Zero, SP1)
- 커뮤니티 기여
- 특정 부분만 최적화

**학습 자료**:
```
1. RISC-V Specs: https://riscv.org/specifications/
2. ZK Proofs: "Proofs, Arguments, and Zero-Knowledge" (Justin Thaler)
3. SP1 Docs: https://docs.succinct.xyz/
4. RISC Zero Docs: https://dev.risczero.com/
```

---

### Q5: ZK proof를 사용하면 프라이버시가 보장되나요?

**A**: **아닙니다**. ZK proof ≠ Privacy.

**구분**:
```
ZK Proof (Zero-Knowledge):
  └─> "나는 X를 알고 있다" 증명 (X를 공개하지 않고)

Privacy:
  └─> 트랜잭션 내용 숨김
```

**OP Stack + ZK**:
- 실행 정확성 증명 (validity)
- 트랜잭션 내용은 공개 (투명성 유지)

**프라이버시가 필요하면**:
- Aztec, Aleo 같은 Privacy rollup 사용
- zk-SNARK + encryption 결합

---

### Q6: SP1과 RISC Zero 중 어떤 것이 더 좋나요?

**A**: 상황에 따라 다릅니다.

| 항목 | SP1 | RISC Zero |
|-----|-----|-----------|
| **증명 속도** | 매우 빠름 (10-100배) | 빠름 |
| **성숙도** | 개발 중 | 프로덕션 ready |
| **생태계** | Optimism 협력 | 광범위한 사용 |
| **학습 곡선** | 중간 | 쉬움 |
| **커뮤니티** | 성장 중 | 활발함 |

**권장**:
```
프로덕션 배포: RISC Zero (안정성)
OP Stack 통합: SP1 (공식 협력)
연구 프로젝트: 둘 다 시도
```

---

### Q7: ZK proof 비용은 어떻게 청구되나요?

**A**: 여러 모델이 가능합니다.

**1. Sequencer 부담**:
```
Sequencer가 ZK proof 비용 지불
└─> 트랜잭션 수수료에 포함
```

**2. 사용자 부담 (선택적)**:
```
일반: 저렴 (fraud proof) + 7일 대기
Fast: 비쌈 (ZK proof) + 즉시 finality
```

**3. 분산 Prover 네트워크**:
```
Prover 경쟁 시장
├─ 가장 저렴한 prover 선택
└─ 토큰 인센티브
```

**예시 (Bonsai)**:
```
RISC Zero Bonsai API:
  - Pay-as-you-go
  - ~$0.01-1 per proof
  - Credit 시스템
```

---

### Q8: 어떤 프로그램이든 ZK proof로 증명할 수 있나요?

**A**: 대부분 가능하지만 제약이 있습니다.

**지원**:
- ✅ 결정론적 계산
- ✅ RISC-V instruction
- ✅ 메모리 접근 (제한적)

**제약**:
- ❌ 비결정론적 동작 (random, time)
- ❌ 무한 루프
- ❌ 너무 큰 프로그램 (>수백MB 메모리)

**우회 방법**:
```rust
// 비결정론적 입력은 public input으로
let random_seed = sp1_zkvm::io::read::<[u8; 32]>();

// 무한 루프는 bounded로
for i in 0..MAX_ITERATIONS {
    if done { break; }
}

// 큰 데이터는 Merkle root만 증명
let data_commitment = merkle_root(large_data);
sp1_zkvm::io::commit(&data_commitment);
```

---

## 참고 자료

### 공식 문서

- [OP Stack Specs](https://specs.optimism.io/)
- [RISC Zero Docs](https://dev.risczero.com/)
- [SP1 Docs](https://docs.succinct.xyz/)
- [Optimism ZK Research](https://github.com/ethereum-optimism/optimism/discussions)

### 프로젝트

- [RISC Zero](https://github.com/risc0/risc0)
- [SP1](https://github.com/succinctlabs/sp1)
- [OP Succinct](https://github.com/succinctlabs/op-succinct)
- [Valida](https://github.com/lita-xyz/valida)

### 논문 및 자료

- [zkSNARKs in a Nutshell](https://chriseth.github.io/notes/articles/zksnarks/zksnarks.pdf)
- [STARKs Paper](https://eprint.iacr.org/2018/046.pdf)
- [Plonky3](https://github.com/Plonky3/Plonky3)
- [Lasso/Jolt](https://eprint.iacr.org/2023/1217)

### 커뮤니티

- [Optimism Discord](https://discord.optimism.io)
- [RISC Zero Discord](https://discord.gg/risczero)
- [ZK Research Forum](https://zkresear.ch/)

### 관련 문서

- [Asterisc (RISC-V) 가이드](./asterisc-riscv-guide-ko.md)
- [게임 타입과 VM 매핑](./game-types-and-vms-ko.md)
- [op-challenger 아키텍처](./op-challenger-architecture-ko.md)

---

**마지막 업데이트**: 2025-01-21
**버전**: 연구 단계 (ZK 통합은 미래 기능)
**참고**: 이 문서는 미래 가능성을 다루며, 현재 tokamak-thanos는 fraud proof만 지원합니다.
