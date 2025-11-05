# op-geth v1.101601.0 업그레이드 분석 및 수정 계획

> **작성일**: 2025-11-04
> **기준**: op-geth v1.101601.0-rc.1, go-ethereum v1.15.11
> **목적**: tokamak-thanos의 op-geth 업그레이드 과정에서 발생한 호환성 문제 분석 및 수정 방안

## 📋 목차

1. [문제 배경](#1-문제-배경)
2. [근본 원인 분석](#2-근본-원인-분석)
3. [호환성 문제 목록](#3-호환성-문제-목록)
4. [수정 계획](#4-수정-계획)
5. [파일별 변경 사항](#5-파일별-변경-사항)
6. [남은 작업 체크리스트](#6-남은-작업-체크리스트)
7. [테스트 계획](#7-테스트-계획)

---

## 1. 문제 배경

### 1.1 초기 증상

E2E 테스트 `TestOutputAlphabetGame_ReclaimBond` 실행 시 다음 오류 발생:

```
Error: missing required field 'sourceHash' in transaction
```

**발생 위치**: Sequencer가 Block #1을 생성하려고 할 때 L1 attributes deposit 트랜잭션 처리 중

### 1.2 문제 발견 경로

1. 초기: Genesis 설정 문제로 의심 → `make devnet-allocs` 실행
2. 조사: `op-chain-ops/genesis/config.go` 수정 사항 확인
3. 비교: Optimism 코드와 비교 → 동일한 코드로 Optimism은 성공
4. **핵심 발견**: geth 버전 차이 확인

| 프로젝트 | op-geth 버전 | go-ethereum 버전 | 테스트 결과 |
|---------|-------------|-----------------|-----------|
| Optimism | v1.101601.0-rc.1 | v1.15.11 | ✅ 통과 |
| tokamak-thanos (이전) | tokamak-thanos-geth | v1.13.15 | ❌ 실패 |

### 1.3 해결 방향

커스텀 `tokamak-thanos-geth` 사용 중단 → 공식 Optimism `op-geth v1.101601.0-rc.1` 사용

**장점**:
- Optimism 메인넷과 동일한 환경
- 최신 EIP 지원 (EIP-4844 blob transactions 등)
- 커뮤니티 지원 및 업데이트

**단점**:
- API breaking changes로 인한 코드 수정 필요
- 의존성 업데이트 (Go 1.21 → 1.23)

---

## 2. 근본 원인 분석

### 2.1 sourceHash 오류의 진짜 이유

```
missing required field 'sourceHash' in transaction
```

이 오류는 genesis 설정 문제가 **아니라** geth 버전 불일치 문제였습니다.

**메커니즘**:
1. L1 attributes는 L2 블록 생성 시 deposit transaction으로 전달됨
2. 새로운 op-geth는 `sourceHash` 필드를 **필수**로 요구
3. 구버전 tokamak-thanos-geth는 이 필드를 인식하지 못함
4. 결과: 트랜잭션 검증 실패

### 2.2 API Breaking Changes

op-geth v1.101601.0으로 업그레이드하면서 여러 API가 변경됨:

#### 2.2.1 EIP-4844 관련 변경 (Blob Transactions)

**CalcBlobFee 함수**:
```go
// 이전 (v1.13.x)
func CalcBlobFee(excessBlobGas uint64) *big.Int

// 현재 (v1.15.x)
func CalcBlobFee(cfg *params.ChainConfig, header *types.Header) *big.Int
```

**KZG Blob 타입**:
```go
// 이전
func BlobToCommitment(blob *kzg4844.Blob) (kzg4844.Commitment, error)

// 현재
func BlobToCommitment(blob kzg4844.Blob) (kzg4844.Commitment, error)
```

#### 2.2.2 ChainConfig 변경

```go
// 제거된 필드
type ChainConfig struct {
    // ...
    TerminalTotalDifficultyPassed bool  // ❌ 제거됨
}
```

#### 2.2.3 StateDB 인터페이스 변경

```go
// 새로 추가된 메서드
type StateDB interface {
    // ... 기존 메서드들
    AccessEvents() *types.AccessEvents  // ✅ 신규
}
```

#### 2.2.4 로깅 패키지 변경

```go
// 이전: experimental slog
import "golang.org/x/exp/slog"

// 현재: 표준 라이브러리
import "log/slog"
```

---

## 3. 호환성 문제 목록

### 3.1 빌드 오류 전체 목록

`build_errors.log`에서 추출한 모든 오류:

#### 카테고리 A: Log 패키지 (op-service/log)

| 오류 | 파일 | 라인 | 심각도 |
|-----|------|-----|-------|
| `undefined: LogfmtMsHandlerWithLevel` | cli.go | 152 | 🔴 높음 |
| `undefined: JSONMsHandler` | cli.go | 157 | 🔴 높음 |
| `DynamicLogHandler` 타입 충돌 | cli.go | 230 | 🔴 높음 |
| `slog.Level` vs `exp/slog.Level` 불일치 | cli.go | 230, 15, 16 | 🔴 높음 |
| Level switch case 타입 불일치 | writer.go | 20-28 | 🔴 높음 |

#### 카테고리 B: State 패키지 (op-chain-ops/state)

| 오류 | 파일 | 라인 | 심각도 |
|-----|------|-----|-------|
| `AccessEvents()` 메서드 누락 | memory_db.go | 18 | 🟡 중간 |

#### 카테고리 C: TxMgr 패키지 (op-service/txmgr)

| 오류 | 파일 | 라인 | 심각도 |
|-----|------|-----|-------|
| `kzg4844.Blob` 포인터/값 타입 불일치 | txmgr.go | 330, 335 | 🟡 중간 |
| `CalcBlobFee` 인자 개수 불일치 | txmgr.go | 602, 758 | 🟡 중간 |

### 3.2 이미 수정된 항목 ✅

| 항목 | 파일 | 수정 내용 |
|-----|------|----------|
| CalcBlobFee (block_info) | op-service/eth/block_info.go | `CalcBlobFeeDefault` 래퍼 함수 사용 |
| CalcBlobFee 래퍼 생성 | op-service/eth/blob.go | 신규 파일 생성 (Optimism 코드 복사) |
| TerminalTotalDifficultyPassed | op-chain-ops/deployer/deployer.go | 필드 제거 |
| Go 버전 | go.mod | 1.21 → 1.23 |
| op-geth 의존성 | go.mod | tokamak-thanos-geth → op-geth v1.101601.0-rc.1 |
| 로그 패키지 (부분) | op-service/log/cli.go | Optimism 버전으로 교체 |
| 로그 서브패키지 | op-service/log{filter,mods}/*, op-service/tri/ | Optimism에서 복사 |

---

## 4. 수정 계획

### 4.1 우선순위 정의

| 우선순위 | 카테고리 | 이유 |
|---------|---------|------|
| P0 (긴급) | Log 패키지 | 가장 많은 오류, 다른 패키지에 영향 |
| P1 (높음) | State 패키지 | 핵심 인터페이스 구현 |
| P1 (높음) | TxMgr 패키지 | 트랜잭션 처리 핵심 기능 |
| P2 (중간) | 전체 빌드 검증 | 통합 테스트 |
| P3 (낮음) | E2E 테스트 | sourceHash 오류 재현 및 해결 확인 |

### 4.2 단계별 작업 계획

#### Phase 1: 누락된 Log 파일 복사
**예상 시간**: 10분
**난이도**: ⭐ 낮음

```bash
# Optimism 저장소에서 누락된 파일 확인 및 복사
cd /Users/zena/tokamak-projects/optimism/op-service/log
# 필요한 파일: dynamic_handler.go, ms_handler.go, writer.go 등
```

**필요한 파일**:
- `dynamic_handler.go` - `DynamicLogHandler` 정의
- `ms_handler.go` - `LogfmtMsHandlerWithLevel`, `JSONMsHandler` 정의
- `writer.go` - slog 호환 버전
- 기타 누락된 유틸리티 파일

#### Phase 2: Import Path 일괄 변경
**예상 시간**: 5분
**난이도**: ⭐ 낮음

```bash
# Optimism → tokamak-thanos 경로 변경
find ./op-service/log -type f -name "*.go" -exec sed -i '' \
  's|github.com/ethereum-optimism/optimism|github.com/tokamak-network/tokamak-thanos|g' {} +
```

#### Phase 3: MemoryStateDB에 AccessEvents() 추가
**예상 시간**: 15분
**난이도**: ⭐⭐ 중간

```go
// op-chain-ops/state/memory_db.go
func (db *MemoryStateDB) AccessEvents() *types.AccessEvents {
    // EIP-2930/7702 access events 구현
    // Optimism 구현 참고
}
```

**참고 파일**:
- `/Users/zena/tokamak-projects/optimism/op-chain-ops/state/memory_db.go`

#### Phase 4: TxMgr의 kzg4844.Blob 타입 수정
**예상 시간**: 10분
**난이도**: ⭐⭐ 중간

```go
// op-service/txmgr/txmgr.go

// 수정 전
commitment, err := kzg4844.BlobToCommitment(rawBlob)  // rawBlob: kzg4844.Blob

// 수정 후 (포인터 제거)
commitment, err := kzg4844.BlobToCommitment(rawBlob)  // 이미 값 타입
// 또는
commitment, err := kzg4844.BlobToCommitment(*rawBlob)  // 포인터였다면
```

#### Phase 5: TxMgr의 CalcBlobFee 호출 수정
**예상 시간**: 10분
**난이도**: ⭐ 낮음

```go
// op-service/txmgr/txmgr.go:602, 758

// 수정 전
blobFee := eip4844.CalcBlobFee(excessBlobGas)

// 수정 후
blobFee := eth.CalcBlobFeeDefault(header)
```

**import 추가**:
```go
import (
    "github.com/tokamak-network/tokamak-thanos/op-service/eth"
)
```

#### Phase 6: 전체 빌드 검증
**예상 시간**: 5분
**난이도**: ⭐ 낮음

```bash
make build-go 2>&1 | tee build_verification.log
# 오류 없음 확인
echo $? # 0이어야 함
```

#### Phase 7: E2E 테스트 실행
**예상 시간**: 10분
**난이도**: ⭐⭐ 중간

```bash
cd op-e2e
go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./faultproofs
```

**성공 기준**:
- ✅ `sourceHash` 오류 없음
- ✅ 테스트 통과

---

## 5. 파일별 변경 사항

### 5.1 이미 수정 완료된 파일

#### `go.mod`
```diff
- go 1.21
+ go 1.23

- replace github.com/ethereum/go-ethereum => github.com/tokamak-network/tokamak-thanos-geth v0.0.0-20250316144452-ffef43a7e0ca
+ replace github.com/ethereum/go-ethereum => github.com/ethereum-optimism/op-geth v1.101601.0-rc.1
```

#### `op-service/eth/blob.go` (신규 생성)
```go
// CalcBlobFeeDefault는 op-geth v1.15.x의 새로운 API를 위한 래퍼
func CalcBlobFeeDefault(header *types.Header) *big.Int {
    dummyChainCfg := &params.ChainConfig{
        LondonBlock:        common.Big0,
        CancunTime:         ptr(uint64(0)),
        BlobScheduleConfig: params.DefaultBlobSchedule,
    }
    if header.RequestsHash != nil {
        dummyChainCfg.PragueTime = ptr(uint64(0))
    }
    return eip4844.CalcBlobFee(dummyChainCfg, header)
}
```

#### `op-service/eth/block_info.go`
```diff
  func (b blockInfo) BlobBaseFee() *big.Int {
      ebg := b.ExcessBlobGas()
      if ebg == nil {
          return nil
      }
-     return eip4844.CalcBlobFee(*ebg)
+     return CalcBlobFeeDefault(b.Header())
  }
```

#### `op-chain-ops/deployer/deployer.go`
```diff
  MergeNetsplitBlock:      big.NewInt(0),
  TerminalTotalDifficulty: big.NewInt(-1),
- TerminalTotalDifficultyPassed: true,
```

#### `op-service/log/cli.go`
- 전체 파일 Optimism 버전으로 교체
- `golang.org/x/exp/slog` → `log/slog`
- Import paths 업데이트

### 5.2 수정 예정 파일

#### `op-service/log/` (누락된 파일들)
- [ ] `dynamic_handler.go`
- [ ] `ms_handler.go`
- [ ] `writer.go` (slog 버전 충돌 수정)
- [ ] 기타 의존 파일

#### `op-chain-ops/state/memory_db.go`
```go
// 추가 필요
func (db *MemoryStateDB) AccessEvents() *types.AccessEvents {
    // TODO: Optimism 구현 참고하여 구현
    return &types.AccessEvents{}
}
```

#### `op-service/txmgr/txmgr.go`
```diff
  // Line 330, 335: kzg4844.Blob 포인터 제거
- commitment, err := kzg4844.BlobToCommitment(&rawBlob)
+ commitment, err := kzg4844.BlobToCommitment(rawBlob)

- proof, err := kzg4844.ComputeBlobProof(&rawBlob, commitment)
+ proof, err := kzg4844.ComputeBlobProof(rawBlob, commitment)

  // Line 602, 758: CalcBlobFee 호출 변경
- blobFee := eip4844.CalcBlobFee(excessBlobGas)
+ blobFee := eth.CalcBlobFeeDefault(header)
```

---

## 6. 남은 작업 체크리스트

### 6.1 작업 상태

| # | 작업 | 우선순위 | 예상 시간 | 난이도 | 상태 |
|---|------|---------|----------|-------|------|
| 1 | Optimism에서 남은 log 파일 복사 | P0 | 10분 | ⭐ | 🔄 대기 중 |
| 2 | DynamicLogHandler slog 버전 충돌 수정 | P0 | 5분 | ⭐ | 🔄 대기 중 |
| 3 | MemoryStateDB에 AccessEvents() 추가 | P1 | 15분 | ⭐⭐ | 🔄 대기 중 |
| 4 | kzg4844.Blob 포인터/값 타입 수정 | P1 | 10분 | ⭐⭐ | 🔄 대기 중 |
| 5 | txmgr의 CalcBlobFee 호출 수정 | P1 | 10분 | ⭐ | 🔄 대기 중 |
| 6 | 전체 빌드 및 검증 | P2 | 5분 | ⭐ | 🔄 대기 중 |
| 7 | 업데이트된 geth로 E2E 테스트 | P3 | 10분 | ⭐⭐ | 🔄 대기 중 |

**총 예상 시간**: ~65분

### 6.2 의존성 그래프

```
Task 1 (Log 파일 복사)
  │
  ├─> Task 2 (slog 버전 충돌)
  │     │
  │     └─> Task 6 (빌드 검증)
  │           │
  ├───────────┘
  │
Task 3 (AccessEvents) ────┐
  │                        │
Task 4 (kzg4844.Blob) ────┼─> Task 6 (빌드 검증)
  │                        │       │
Task 5 (CalcBlobFee) ─────┘       │
                                   │
                                   └─> Task 7 (E2E 테스트)
```

**병렬 처리 가능**:
- Task 1-2 (Log 패키지)는 독립적
- Task 3-5는 서로 독립적 (병렬 수정 가능)
- Task 6은 모든 작업 완료 후
- Task 7은 Task 6 성공 후

### 6.3 위험 요소

| 위험 | 확률 | 영향 | 완화 방안 |
|-----|------|------|----------|
| Optimism 코드에 추가 의존성 존재 | 중간 | 높음 | 빌드하면서 점진적으로 추가 |
| AccessEvents 구현 복잡도 | 낮음 | 중간 | Optimism 구현 그대로 복사 |
| kzg4844 API 추가 변경 | 낮음 | 높음 | Optimism txmgr.go 참고 |
| E2E 테스트에서 새로운 오류 | 중간 | 중간 | 로그 분석 후 단계적 해결 |

---

## 7. 테스트 계획

### 7.1 단위 테스트

```bash
# 각 패키지별 테스트
go test ./op-service/log/...
go test ./op-service/eth/...
go test ./op-chain-ops/state/...
go test ./op-service/txmgr/...
```

### 7.2 통합 빌드 테스트

```bash
# 전체 빌드
make build-go

# 특정 타겟 빌드
make -C ./op-node op-node
make -C ./op-batcher op-batcher
make -C ./op-proposer op-proposer
make -C ./op-challenger op-challenger
```

### 7.3 E2E 테스트

#### 주요 테스트 케이스

| 테스트 | 목적 | 기대 결과 |
|-------|------|----------|
| `TestOutputAlphabetGame_ReclaimBond` | sourceHash 오류 해결 확인 | ✅ 통과 |
| `TestBlockBuilding` | L1→L2 deposit 처리 | ✅ sourceHash 필드 존재 |
| `TestBlobTransactions` | EIP-4844 호환성 | ✅ Blob fee 계산 정상 |
| `TestStateTransition` | StateDB 인터페이스 | ✅ AccessEvents 호출 성공 |

#### 실행 명령어

```bash
# 단일 테스트
cd op-e2e
go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./faultproofs

# Fault Proof 전체 테스트
go test -v -timeout 30m ./faultproofs/...

# 전체 E2E 테스트 (시간 소요 많음)
go test -v -timeout 60m ./...
```

### 7.4 회귀 테스트

기존 기능이 정상 작동하는지 확인:

```bash
# Challenger 핵심 기능
go test -v ./op-challenger/game/fault/...

# Chain ops
go test -v ./op-chain-ops/genesis/...

# 서비스 레이어
go test -v ./op-service/...
```

### 7.5 성공 기준

✅ **빌드 성공**:
- `make build-go` 오류 없이 완료
- 모든 바이너리 생성 확인

✅ **테스트 통과**:
- 단위 테스트 100% 통과
- E2E 테스트 `TestOutputAlphabetGame_ReclaimBond` 통과
- sourceHash 오류 발생하지 않음

✅ **기능 검증**:
- L1→L2 deposit 트랜잭션 정상 처리
- Blob 트랜잭션 fee 계산 정확
- 로깅 시스템 정상 작동

---

## 8. 참고 자료

### 8.1 주요 파일 경로

**Optimism 참조**:
```
/Users/zena/tokamak-projects/optimism/
├── op-service/log/
│   ├── cli.go
│   ├── dynamic_handler.go
│   ├── ms_handler.go
│   └── writer.go
├── op-service/eth/blob.go
├── op-chain-ops/state/memory_db.go
└── op-service/txmgr/txmgr.go
```

**tokamak-thanos 수정**:
```
/Users/zena/tokamak-projects/tokamak-thanos/
├── go.mod
├── op-service/
│   ├── log/
│   ├── eth/
│   └── txmgr/
└── op-chain-ops/
    ├── deployer/
    └── state/
```

### 8.2 관련 문서

- [Optimism op-geth 릴리스 노트](https://github.com/ethereum-optimism/op-geth/releases/tag/v1.101601.0-rc.1)
- [EIP-4844: Shard Blob Transactions](https://eips.ethereum.org/EIPS/eip-4844)
- [go-ethereum v1.15.x Changes](https://github.com/ethereum/go-ethereum/releases/tag/v1.15.11)

### 8.3 빌드 로그

- `build_errors.log` - 초기 빌드 오류 전체 목록
- `build_verification.log` - 수정 후 검증 빌드 (예정)

---

## 9. 결론

### 9.1 핵심 요약

1. **문제**: `sourceHash` 오류는 geth 버전 불일치가 원인
2. **해결**: tokamak-thanos-geth → op-geth v1.101601.0-rc.1 전환
3. **영향**: API breaking changes로 7개 작업 필요
4. **예상 시간**: 약 65분
5. **위험도**: 낮음 (Optimism 코드 참조 가능)

### 9.2 다음 단계

1. ✅ 이 문서를 기반으로 팀 리뷰
2. 🔄 Phase 1부터 순차적으로 작업 진행
3. 🔄 각 Phase마다 빌드 및 테스트 확인
4. 🔄 최종 E2E 테스트로 sourceHash 오류 해결 검증

### 9.3 추가 고려사항

- **향후 업그레이드**: Optimism 업스트림 변경사항을 정기적으로 반영하는 프로세스 수립 필요
- **테스트 자동화**: CI/CD에 op-geth 버전 호환성 테스트 추가
- **문서화**: 각 수정 사항을 코드 주석으로 명확히 표시

---

> **작성자**: Claude Code
> **문서 버전**: 1.0
> **최종 수정**: 2025-11-04
