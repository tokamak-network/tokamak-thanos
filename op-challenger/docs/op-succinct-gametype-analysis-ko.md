# OP Succinct (GameType 6) 분석 및 ZK Proof 통합

## 개요

Optimism 프로젝트에서 **GameType 6 (OP_SUCCINCT)**는 2025년 1월 17일에 추가된 최신 게임 타입입니다. 이는 **ZK proof (영지식 증명)**를 사용하여 dispute game을 처리하기 위해 예약된 타입이지만, 현재는 **타입 정의만 존재하고 실제 구현은 없는 상태**입니다.

## 현재 상태

### 추가된 내용 (커밋 477dbc668, 2025-01-17)

#### 1. Go 코드 정의
**파일**: `op-challenger/game/fault/types/types.go:36`

```go
const (
    CannonGameType            GameType = 0
    PermissionedGameType      GameType = 1
    AsteriscGameType          GameType = 2
    AsteriscKonaGameType      GameType = 3
    SuperCannonGameType       GameType = 4
    SuperPermissionedGameType GameType = 5
    OPSuccinctGameType        GameType = 6  // 새로 추가
    SuperAsteriscKonaGameType GameType = 7
    FastGameType              GameType = 254
    AlphabetGameType          GameType = 255
)

func (t GameType) String() string {
    // ...
    case OPSuccinctGameType:
        return "op-succinct"
    // ...
}
```

#### 2. Solidity 컨트랙트 정의
**파일**: `packages/contracts-bedrock/src/dispute/lib/Types.sol:70-71`

```solidity
library GameTypes {
    /// @notice A dispute game type that uses OP Succinct
    GameType internal constant OP_SUCCINCT = GameType.wrap(6);
}
```

#### 3. 배포 스크립트 지원
**파일**: `packages/contracts-bedrock/scripts/deploy/Deploy.s.sol`

```solidity
if (rawGameType == GameTypes.OP_SUCCINCT.raw()) {
    gameTypeString = "OP Succinct";
}
```

### 구현되지 않은 내용

다음 핵심 구성 요소들이 **아직 구현되지 않았습니다**:

1. **TraceProvider 구현 없음**
   - `op-challenger/game/fault/trace/` 디렉토리에 OP Succinct 관련 구현 없음
   - cannon, asterisc와 같은 trace provider가 필요하지만 존재하지 않음

2. **게임 등록 코드 없음**
   - `op-challenger/game/fault/register.go`에 GameType 6 등록 코드 없음
   - 다른 게임 타입(0,1,2,3,4,5,7,254,255)은 모두 등록되어 있음
   - **GameType 6만 건너뛰어짐**

3. **VM 통합 없음**
   - ZK proof 생성을 위한 VM 통합 코드 없음
   - SP1 zkVM과의 연동 로직 없음

4. **설정 및 테스트 없음**
   - OP Succinct 관련 설정 파일 없음
   - 테스트 코드 없음

## OP Succinct란?

### 개념

**OP Succinct**는 Optimism의 fault proof 시스템에 **ZK proof (영지식 증명)**를 통합하는 프로젝트입니다. 기존의 fraud proof 대신 validity proof를 사용하여 L2 상태 전환을 검증합니다.

### 핵심 특징

#### 1. Fraud Proof vs ZK Proof 비교

| 특징 | Fraud Proof (기존) | ZK Proof (OP Succinct) |
|------|-------------------|----------------------|
| **검증 방식** | 이의 제기 + 반증 | 수학적 증명 |
| **검증 시간** | 7일 challenge period | 즉시 (증명 생성 후) |
| **가스 비용** | 높음 (bisection game) | 낮음 (증명 검증만) |
| **보안 가정** | 1-of-N 정직한 검증자 | 암호학적 보안 |
| **최종성** | 지연됨 (challenge period) | 빠름 |

#### 2. SP1 zkVM

OP Succinct는 **Succinct Labs**의 **SP1 zkVM**을 사용할 것으로 예상됩니다:

- **아키텍처**: RISC-V 기반 zkVM
- **특징**: 범용 프로그램을 ZK proof로 실행
- **성능**: 높은 증명 생성 속도
- **호환성**: Rust 코드를 직접 실행 가능

## 다른 ZK GameType: KAILUA (GameType 1337)

Optimism은 OP Succinct 외에도 **RISC Zero의 Kailua**를 지원하기 위한 GameType도 정의했습니다.

**파일**: `packages/contracts-bedrock/src/dispute/lib/Types.sol:81-82`

```solidity
/// @notice A dispute game type that uses RISC Zero's Kailua
GameType internal constant KAILUA = GameType.wrap(1337);
```

### Kailua vs OP Succinct

| 특징 | OP Succinct (Type 6) | Kailua (Type 1337) |
|------|---------------------|-------------------|
| **zkVM** | SP1 (Succinct Labs) | RISC Zero |
| **아키텍처** | RISC-V | RISC-V |
| **개발사** | Succinct Labs | RISC Zero |
| **상태** | 정의만 존재 | 정의만 존재 |

## 왜 구현이 없는가?

### 1. 외부 프로젝트로 개발 중

OP Succinct는 **Succinct Labs**에서 별도 프로젝트로 개발 중일 가능성이 높습니다:
- Optimism 메인 레포지토리와 독립적으로 개발
- 성숙해지면 통합될 예정

### 2. 실험적 단계

ZK proof 통합은 아직 실험적 단계:
- 성능 최적화 필요
- 경제성 검증 필요
- 보안 감사 필요

### 3. 점진적 통합 전략

Optimism은 여러 증명 시스템을 지원하는 전략:
- 먼저 타입 정의를 추가 (인터페이스 확장)
- 나중에 구현 통합 (모듈식 아키텍처)

## 예상 구현 아키텍처

OP Succinct가 구현될 때 필요한 구성 요소:

### 1. Trace Provider
```
op-challenger/game/fault/trace/succinct/
├── provider.go          # TraceProvider 인터페이스 구현
├── vm.go                # SP1 zkVM 통합
├── prover.go            # ZK proof 생성
└── config.go            # OP Succinct 설정
```

### 2. 게임 등록
```go
// op-challenger/game/fault/register.go
func RegisterGameTypes(...) {
    // ...
    registry.RegisterGameType(types.OPSuccinctGameType,
        &succinct.TraceProvider{...})
}
```

### 3. 증명 생성 플로우
```
1. L2 상태 전환 발생
2. OP Succinct provider가 실행 trace 수집
3. SP1 zkVM이 trace로부터 증명 생성
4. 증명을 온체인에 제출
5. 스마트 컨트랙트가 증명 검증
6. 검증 성공 시 즉시 최종성 확보
```

## 현재 사용 가능한 GameType

Optimism에서 **실제로 구현되어 사용 가능한** GameType:

| GameType | 이름 | VM | 구현 상태 |
|----------|------|-----|----------|
| 0 | CANNON | MIPS | 완전 구현 |
| 1 | PERMISSIONED_CANNON | MIPS | 완전 구현 |
| 2 | ASTERISC | RISC-V | 완전 구현 |
| 3 | ASTERISC_KONA | RISC-V + Kona | 완전 구현 |
| 4 | SUPER_CANNON | MIPS (Super Roots) | 완전 구현 |
| 5 | SUPER_PERMISSIONED_CANNON | MIPS (Super Roots) | 완전 구현 |
| **6** | **OP_SUCCINCT** | **SP1 zkVM** | **정의만 존재** |
| 7 | SUPER_ASTERISC_KONA | RISC-V (Super Roots) | 완전 구현 |
| 254 | FAST | 테스트용 | 완전 구현 |
| 255 | ALPHABET | 테스트용 | 완전 구현 |
| **1337** | **KAILUA** | **RISC Zero** | **정의만 존재** |

## 향후 전망

### 1. ZK Proof 통합의 이점

OP Succinct가 구현되면:
- **빠른 최종성**: 7일 challenge period 불필요
- **낮은 비용**: bisection game 대신 증명 검증만
- **높은 보안**: 암호학적 보안 보장
- **확장성**: 더 많은 트랜잭션 처리 가능

### 2. 하이브리드 접근

Optimism은 여러 증명 시스템을 동시에 지원할 수 있습니다:
- **Fraud Proof** (Cannon, Asterisc): 성숙하고 안정적
- **ZK Proof** (OP Succinct, Kailua): 빠르고 효율적
- 사용자/체인이 선택 가능한 유연성

### 3. 생태계 발전

ZK 통합은 L2 생태계 전체의 발전을 촉진:
- **상호운용성**: 다양한 zkVM 지원
- **혁신**: 새로운 증명 시스템 실험 가능
- **경쟁**: 여러 솔루션이 경쟁하며 발전

## 결론

**GameType 6 (OP_SUCCINCT)**는:

1. **정의만 존재**: 2025년 1월 17일 추가
2. **구현 없음**: TraceProvider, VM 통합, 등록 코드 모두 없음
3. **미래를 위한 예약**: ZK proof 통합을 위한 인터페이스 확장
4. **외부 개발**: Succinct Labs에서 별도 프로젝트로 개발 중일 가능성

현재 Optimism은 **fraud proof 기반 시스템**(Cannon, Asterisc)만 실제로 사용 가능하며, ZK proof 통합은 **향후 계획**입니다.

이는 Optimism이 **점진적이고 모듈식 아키텍처**를 채택하여, 다양한 증명 시스템을 유연하게 통합할 수 있도록 설계되었음을 보여줍니다.

## 참고 자료

### 커밋 히스토리
- **477dbc668**: feat: add OP_SUCCINCT game type (#13780)
  - 날짜: 2025-01-17
  - 작성자: Aurélien
  - 내용: GameType 6 정의 추가 (Go, Solidity, Deploy script)

### 관련 파일
- `op-challenger/game/fault/types/types.go` - GameType 정의
- `packages/contracts-bedrock/src/dispute/lib/Types.sol` - Solidity 타입 정의
- `op-challenger/game/fault/register.go` - 게임 등록 (OP Succinct 없음)

### 외부 프로젝트
- **Succinct Labs**: https://succinct.xyz/
- **SP1 zkVM**: RISC-V 기반 영지식 증명 가상 머신
- **RISC Zero**: https://www.risczero.com/
- **Kailua**: RISC Zero의 zkVM 프로젝트
