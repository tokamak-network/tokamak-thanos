# RISC-V 기반 GameType 비교 (2, 3, 7)

## 개요

Optimism은 RISC-V 아키텍처를 사용하는 **3가지 GameType**을 지원합니다:

| GameType | 이름 | VM | 프로그램 | 대상 체인 | 구현 언어 |
|----------|------|-----|---------|----------|----------|
| **2** | **ASTERISC** | RISC-V | op-program | 단일 L2 | Go |
| **3** | **ASTERISC_KONA** | RISC-V | kona-client | 단일 L2 | Rust |
| **7** | **SUPER_ASTERISC_KONA** | RISC-V | kona-client | Superchain (다중 L2) | Rust |

세 가지 모두 **동일한 RISC-V VM (Asterisc)**을 사용하지만, **실행하는 프로그램과 아키텍처가 다릅니다**.

## 상세 비교

### 1. GameType 2: ASTERISC (기본)

#### 특징
- **VM**: Asterisc (RISC-V RV64GC)
- **프로그램**: **op-program** (Go 기반)
- **대상**: 단일 L2 체인
- **Executor**: `OpProgramServerExecutor`

#### 등록 코드
**파일**: `op-challenger/game/fault/register.go:91-96`

```go
if cfg.TraceTypeEnabled(faultTypes.TraceTypeAsterisc) {
    l2HeaderSource, rollupClient, syncValidator, err := clients.SingleChainClients()
    if err != nil {
        return nil, err
    }
    registerTasks = append(registerTasks, NewAsteriscRegisterTask(
        faultTypes.AsteriscGameType,
        cfg,
        m,
        vm.NewOpProgramServerExecutor(logger),  // Go 기반 op-program 사용
        l2HeaderSource,
        rollupClient,
        syncValidator
    ))
}
```

#### 실행 인자 예시
```bash
op-program --server \
  --l1 <L1_RPC_URL> \
  --l1.beacon <BEACON_URL> \
  --l2 <L2_RPC_URL> \
  --datadir <DATA_DIR> \
  --l1.head <L1_HEAD_HASH> \
  --l2.claim <L2_CLAIM_HASH> \
  --l2.blocknumber <BLOCK_NUMBER>
```

#### 장점
- 성숙하고 안정적 (Optimism의 기본 구현)
- 많은 설정 옵션 지원
- 광범위하게 테스트됨

#### 단점
- Go 언어로 작성 (Rust보다 크기가 큼)
- 단일 체인만 지원

---

### 2. GameType 3: ASTERISC_KONA

#### 특징
- **VM**: Asterisc (RISC-V RV64GC)
- **프로그램**: **kona-client** (Rust 기반)
- **대상**: 단일 L2 체인
- **Executor**: `KonaExecutor`

#### 등록 코드
**파일**: `op-challenger/game/fault/register.go:98-103`

```go
if cfg.TraceTypeEnabled(faultTypes.TraceTypeAsteriscKona) {
    l2HeaderSource, rollupClient, syncValidator, err := clients.SingleChainClients()
    if err != nil {
        return nil, err
    }
    registerTasks = append(registerTasks, NewAsteriscKonaRegisterTask(
        faultTypes.AsteriscKonaGameType,
        cfg,
        m,
        vm.NewKonaExecutor(),  // Rust 기반 kona 사용
        l2HeaderSource,
        rollupClient,
        syncValidator
    ))
}
```

#### 실행 인자 예시
**파일**: `op-challenger/game/fault/trace/vm/kona_server_executor.go:29-40`

```bash
kona-client single \
  --l1-node-address <L1_RPC_URL> \
  --l1-beacon-address <BEACON_URL> \
  --l2-node-address <L2_RPC_URL> \
  --l1-head <L1_HEAD_HASH> \
  --l2-head <L2_HEAD_HASH> \
  --l2-output-root <OUTPUT_ROOT> \
  --l2-claim <L2_CLAIM_HASH> \
  --l2-block-number <BLOCK_NUMBER> \
  --server \
  --data-dir <DATA_DIR>
```

#### Kona란?

**Kona**는 Optimism OP Stack의 **Rust 기반 대체 구현체**입니다:

- **코드 크기**: ~8,000 라인 (Go 구현보다 훨씬 작음)
- **언어**: Rust (`no_std` 지원)
- **특징**:
  - 경량화 및 최적화
  - 모듈식 아키텍처
  - ZK proof 시스템과의 통합 용이 (SP1, RISC Zero 등)
  - 포터블 (다양한 환경에서 실행 가능)

**공식 문서**: https://devdocs.optimism.io/book/

#### 장점
- **경량**: Go 구현보다 훨씬 작은 코드베이스
- **포터블**: `no_std` 지원으로 다양한 환경 지원
- **ZK 친화적**: zkVM (SP1, RISC Zero)과 쉽게 통합
- **모던**: Rust의 안전성과 성능

#### 단점
- 상대적으로 새로운 구현 (2024-2025년 추가)
- 단일 L2만 지원 (코드에서 명시적으로 체크)

#### 제약 사항
**파일**: `op-challenger/game/fault/trace/vm/kona_server_executor.go:26-28`

```go
if len(cfg.L2s) != 1 || len(cfg.RollupConfigPaths) > 1 || len(cfg.Networks) > 1 {
    return nil, errors.New("multiple L2s specified but only one supported")
}
```

---

### 3. GameType 7: SUPER_ASTERISC_KONA

#### 특징
- **VM**: Asterisc (RISC-V RV64GC)
- **프로그램**: **kona-client** (Rust 기반)
- **대상**: **Superchain** (여러 L2 체인)
- **Executor**: `KonaSuperExecutor`

#### 등록 코드
**파일**: `op-challenger/game/fault/register.go:105-110`

```go
if cfg.TraceTypeEnabled(faultTypes.TraceTypeSuperAsteriscKona) {
    rootProvider, syncValidator, err := clients.SuperchainClients()  // Superchain 클라이언트
    if err != nil {
        return nil, err
    }
    registerTasks = append(registerTasks, NewSuperAsteriscKonaRegisterTask(
        faultTypes.SuperAsteriscKonaGameType,
        cfg,
        m,
        vm.NewKonaSuperExecutor(),  // Super 버전 사용
        rootProvider,
        syncValidator
    ))
}
```

#### 실행 인자 예시
**파일**: `op-challenger/game/fault/trace/vm/kona_super_server_executor.go:30-40`

```bash
kona-client super \
  --l1-node-address <L1_RPC_URL> \
  --l1-beacon-address <BEACON_URL> \
  --l2-node-addresses <L2_RPC_1>,<L2_RPC_2>,<L2_RPC_3> \  # 여러 L2 지원
  --l1-head <L1_HEAD_HASH> \
  --agreed-l2-pre-state <AGREED_PRESTATE> \  # 필수
  --claimed-l2-post-state <L2_CLAIM_HASH> \
  --claimed-l2-timestamp <TIMESTAMP> \
  --server \
  --data-dir <DATA_DIR>
```

#### Superchain이란?

**Superchain**은 Optimism의 **다중 L2 체인 아키텍처**입니다:

- 하나의 L1에 여러 L2 체인이 연결
- 체인 간 상호운용성 지원
- 공유 보안 모델
- 통합 dispute resolution

#### 핵심 차이점

1. **여러 L2 지원**: `--l2-node-addresses`에 쉼표로 구분된 여러 RPC URL 전달
2. **AgreedPreState 필수**: Super 모드에서는 사전 합의된 상태가 필수

**파일**: `op-challenger/game/fault/trace/vm/kona_super_server_executor.go:26-28`

```go
if len(inputs.AgreedPreState) == 0 {
    return nil, errors.New("agreed pre-state is not defined")
}
```

3. **SuperchainClients 사용**: 단일 체인 대신 슈퍼체인 클라이언트 사용

#### 장점
- **다중 L2 지원**: 여러 체인을 동시에 처리
- **Superchain 아키텍처**: 미래 지향적 설계
- **Rust 기반**: Kona의 모든 장점 계승

#### 단점
- 더 복잡한 설정 필요
- AgreedPreState 관리 필요

---

## 핵심 차이점 요약

### 1. 프로그램 구현

| 특징 | GameType 2 | GameType 3, 7 |
|------|-----------|--------------|
| **언어** | Go | Rust |
| **프로그램** | op-program | kona-client |
| **코드 크기** | 크다 | 작다 (~8K 라인) |
| **성숙도** | 매우 높음 | 중간 (2024-2025년 추가) |

### 2. 명령어 형식

```bash
# GameType 2: op-program
op-program --server --l1 <URL> --l2 <URL> ...

# GameType 3: kona single
kona-client single --l1-node-address <URL> --l2-node-address <URL> ...

# GameType 7: kona super
kona-client super --l1-node-address <URL> --l2-node-addresses <URL1>,<URL2> ...
```

### 3. 아키텍처

| GameType | 클라이언트 타입 | L2 지원 | 필수 입력 |
|----------|--------------|---------|----------|
| 2 | `SingleChainClients()` | 단일 | L1Head, L2Claim, L2BlockNumber |
| 3 | `SingleChainClients()` | 단일 | L1Head, L2Head, L2OutputRoot, L2Claim |
| 7 | `SuperchainClients()` | **다중** | L1Head, **AgreedPreState**, L2Claim |

### 4. 사용 시나리오

| GameType | 사용 시나리오 |
|----------|-------------|
| **2** | 전통적인 Optimism 단일 체인 (안정성 중시) |
| **3** | 단일 체인 + Rust 최적화 (경량화, ZK 준비) |
| **7** | Superchain 환경 (다중 L2, 상호운용성) |

## 개발 히스토리

### Kona 추가 과정

주요 커밋:
```
c36de049f - chore(ops): Support kona + asterisc in the op-challenger (2024)
7457c5689 - feat(op-challenger): Add TraceTypeAsteriscKona to default (2024)
d164b6d06 - fix: bump kona-client version
3cc36be2e - chore(op-challenger): Update kona executor to use subcommand
bd7c16b87 - feat(op-challenger): Kona interop executor (2025)
4720bc7da - chore: Bump kona version (2025)
```

### 진화 과정

```
2023-2024: GameType 2 (ASTERISC) - Go 기반 op-program
    ↓
2024: GameType 3 (ASTERISC_KONA) - Rust 기반 kona 단일 체인
    ↓
2024-2025: GameType 7 (SUPER_ASTERISC_KONA) - Rust 기반 kona 슈퍼체인
```

## 왜 3가지 버전이 필요한가?

### 1. 언어 다양성 (GameType 2 vs 3)
- **Go (op-program)**: 안정성, 성숙도
- **Rust (kona)**: 경량화, 성능, ZK 친화성

### 2. 구현 다양성 (Redundancy)
- 두 가지 독립적인 구현으로 **클라이언트 다양성** 확보
- 한 구현에 버그가 있어도 다른 구현으로 검증 가능

### 3. 아키텍처 진화 (GameType 3 vs 7)
- **GameType 3**: 전통적인 단일 L2 체인
- **GameType 7**: 미래의 Superchain 아키텍처

### 4. ZK 통합 준비
Kona (Rust)는 **zkVM과의 통합이 쉬움**:
- SP1 zkVM으로 Kona 실행 → **OP Succinct** (GameType 6)
- RISC Zero zkVM으로 Kona 실행 → **Kailua** (GameType 1337)

## Tokamak-Thanos의 지원 현황

Tokamak-Thanos는 **GameType 2 (ASTERISC)만 지원**합니다:

```go
// tokamak-thanos/op-challenger에서는
// GameType 3, 7이 없음
```

### 비교

| 프로젝트 | GameType 2 | GameType 3 | GameType 7 |
|---------|-----------|-----------|-----------|
| **Optimism** | ✅ 완전 구현 | ✅ 완전 구현 | ✅ 완전 구현 |
| **Tokamak-Thanos** | ✅ 구현됨 | ❌ 없음 | ❌ 없음 |

Tokamak-Thanos에 GameType 3, 7을 추가하려면:
1. Kona 의존성 추가
2. `KonaExecutor`, `KonaSuperExecutor` 구현
3. 등록 코드 추가
4. 테스트 및 검증

## 선택 가이드

### 어떤 GameType을 사용해야 하나?

#### GameType 2 (ASTERISC) 추천 대상
- ✅ 안정성이 최우선
- ✅ 성숙한 구현 필요
- ✅ 단일 L2 체인
- ✅ Optimism 메인넷 호환

#### GameType 3 (ASTERISC_KONA) 추천 대상
- ✅ 경량화된 구현 필요
- ✅ Rust 생태계 선호
- ✅ 향후 ZK 통합 계획
- ✅ 단일 L2 체인

#### GameType 7 (SUPER_ASTERISC_KONA) 추천 대상
- ✅ Superchain 아키텍처
- ✅ 여러 L2 체인 관리
- ✅ 체인 간 상호운용성 필요
- ✅ 최신 기술 스택

## 결론

**동일한 RISC-V VM (Asterisc)**을 사용하지만:

1. **GameType 2**: Go 기반 전통적 구현 (안정성)
2. **GameType 3**: Rust 기반 현대적 구현 (경량화)
3. **GameType 7**: Rust 기반 미래 지향적 구현 (Superchain)

모두 **fraud proof** 기반이며, ZK proof는 아닙니다. 하지만 **Kona (GameType 3, 7)**는 향후 **ZK 통합의 기반**이 됩니다:
- Kona + SP1 → OP Succinct (GameType 6)
- Kona + RISC Zero → Kailua (GameType 1337)

Optimism은 **점진적이고 모듈식 접근**을 통해:
- 안정적인 기존 구현 유지 (GameType 2)
- 새로운 경량 구현 도입 (GameType 3)
- 미래 아키텍처 준비 (GameType 7)
- ZK 통합 실험 (GameType 6, 1337)

이런 다양성이 **Optimism의 강력함**입니다.

## 참고 자료

### 공식 문서
- **Kona Book**: https://devdocs.optimism.io/book/
- **Kona GitHub**: https://github.com/ethereum-optimism/kona
- **Asterisc GitHub**: https://github.com/ethereum-optimism/asterisc
- **Optimism 블로그**: https://www.optimism.io/blog/introducing-the-kona-node-a-rust-powered-leap-for-the-op-stack

### 코드 위치
- `op-challenger/game/fault/register.go` - 게임 타입 등록
- `op-challenger/game/fault/trace/vm/op_program_server_executor.go` - GameType 2 executor
- `op-challenger/game/fault/trace/vm/kona_server_executor.go` - GameType 3 executor
- `op-challenger/game/fault/trace/vm/kona_super_server_executor.go` - GameType 7 executor
- `op-challenger/game/fault/types/types.go` - GameType 정의

### 관련 문서
- [게임 타입과 VM 가이드](game-types-and-vms-ko.md)
- [Asterisc (RISC-V) 가이드](asterisc-riscv-guide-ko.md)
- [Optimism vs Tokamak-Thanos 비교](asterisc-comparison-optimism-vs-tokamak-ko.md)
- [OP Succinct (GameType 6) 분석](op-succinct-gametype-analysis-ko.md)
