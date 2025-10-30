# Data Availability (데이터 가용성) 코드 분석

## 목차
1. [개요](#개요)
2. [L1 직접 사용 코드 분석](#l1-직접-사용-코드-분석)
3. [Plasma vs L1 직접 사용 비교](#plasma-vs-l1-직접-사용-비교)
4. [실제 트랜잭션 흐름](#실제-트랜잭션-흐름)
5. [코드 위치 참조](#코드-위치-참조)

---

## 개요

Tokamak Thanos Stack에서 L2 데이터를 L1에 게시하는 방식은 두 가지입니다:

1. **Plasma 모드** (로컬 devnet 전용)
   - da-server에 데이터 저장
   - L1에는 commitment(해시)만 게시
   - 개발/테스트 용도

2. **L1 직접 사용** (프로덕션)
   - 전체 배치 데이터를 L1에 직접 게시
   - Ethereum의 데이터 가용성 보장 활용
   - da-server 불필요

---

## L1 직접 사용 코드 분석

### 1. 핵심 트랜잭션 전송 로직

**파일**: `op-batcher/batcher/driver.go:484-528`

```go
// sendTransaction creates & queues for sending a transaction to the batch inbox address
func (l *BatchSubmitter) sendTransaction(ctx context.Context, txdata txData,
    queue *txmgr.Queue[txID], receiptsCh chan txmgr.TxReceipt[txID]) error {

    var err error
    var candidate *txmgr.TxCandidate

    // 방법 1: EIP-4844 Blob 트랜잭션 사용
    if l.Config.UseBlobs {
        if candidate, err = l.blobTxCandidate(txdata); err != nil {
            return fmt.Errorf("could not create blob tx candidate: %w", err)
        }
    } else {
        // 방법 2: Calldata 트랜잭션 사용
        data := txdata.CallData()

        // ⭐️ Plasma 모드 체크 - 프로덕션에서는 false
        if l.Config.UsePlasma {
            // Plasma 모드: da-server에 데이터를 저장하고 commitment만 L1에 게시
            comm, err := l.PlasmaDA.SetInput(ctx, data)
            if err != nil {
                l.Log.Error("Failed to post input to Plasma DA", "error", err)
                l.recordFailedTx(txdata.ID(), err)
                return nil
            }
            // signal plasma commitment tx with TxDataVersion1
            data = comm.TxData()
        }
        // Plasma 미사용 (프로덕션): 전체 데이터를 L1에 직접 게시
        candidate = l.calldataTxCandidate(data)
    }

    // Gas 추정
    intrinsicGas, err := core.IntrinsicGas(candidate.TxData, nil, false, true, true, false)
    if err != nil {
        l.Log.Error("Failed to calculate intrinsic gas", "err", err)
    } else {
        candidate.GasLimit = intrinsicGas
    }

    // L1에 트랜잭션 전송 큐에 추가
    queue.Send(txdata.ID(), *candidate, receiptsCh)
    return nil
}
```

**핵심 포인트**:
- `l.Config.UsePlasma` 체크로 Plasma 모드 여부 결정
- 프로덕션: `UsePlasma = false` → 504-514줄 블록 실행 안됨
- 515줄의 `calldataTxCandidate(data)` 호출로 전체 데이터가 L1에 게시됨

---

### 2. Calldata 트랜잭션 생성

**파일**: `op-batcher/batcher/driver.go:546-552`

```go
func (l *BatchSubmitter) calldataTxCandidate(data []byte) *txmgr.TxCandidate {
    l.Log.Info("building Calldata transaction candidate", "size", len(data))
    return &txmgr.TxCandidate{
        To:     &l.RollupConfig.BatchInboxAddress,  // ⭐️ L1의 BatchInbox 컨트랙트 주소
        TxData: data,  // ⭐️ 전체 배치 데이터가 calldata로 포함됨
    }
}
```

**설명**:
- `BatchInboxAddress`: L1 Ethereum에 배포된 특수 컨트랙트 주소
- `TxData`: 압축된 L2 배치 데이터 (수십~수백 KB)
- 이 데이터가 L1 트랜잭션의 calldata로 영구 저장됨

---

### 3. Blob 트랜잭션 생성

**파일**: `op-batcher/batcher/driver.go:530-544`

```go
func (l *BatchSubmitter) blobTxCandidate(data txData) (*txmgr.TxCandidate, error) {
    blobs, err := data.Blobs()
    if err != nil {
        return nil, fmt.Errorf("generating blobs for tx data: %w", err)
    }
    size := data.Len()
    lastSize := len(data.frames[len(data.frames)-1].data)
    l.Log.Info("building Blob transaction candidate",
        "size", size, "last_size", lastSize, "num_blobs", len(blobs))
    l.Metr.RecordBlobUsedBytes(lastSize)
    return &txmgr.TxCandidate{
        To:    &l.RollupConfig.BatchInboxAddress,  // ⭐️ L1의 BatchInbox 주소
        Blobs: blobs,  // ⭐️ EIP-4844 blob으로 데이터 게시
    }, nil
}
```

**설명**:
- EIP-4844 blob 트랜잭션 사용 (가스 비용 저렴)
- Blob은 L1에 일시적으로 저장 (~18일)
- Blob commitment는 영구 저장

---

### 4. BatchInboxAddress 정의

**파일**: `op-node/rollup/types.go:119`

```go
type Config struct {
    // ... 다른 필드들 ...

    // L1 address that batches are sent to.
    // ⭐️ L1 Ethereum에 배포된 특수 컨트랙트 주소
    BatchInboxAddress common.Address `json:"batch_inbox_address"`

    // L1 Deposit Contract Address
    DepositContractAddress common.Address `json:"deposit_contract_address"`

    // ...
}
```

**예시 주소**:
```
0xff00000000000000000000000000000000000010  // 일반적인 형식
```

이 주소는:
- L1 컨트랙트 주소
- 모든 배치 데이터가 이 주소로 전송됨
- 네트워크별로 고정된 주소 사용

---

## Plasma vs L1 직접 사용 비교

### Plasma 모드 (로컬 devnet)

```go
// op-batcher/batcher/driver.go:504-513
if l.Config.UsePlasma {
    // 1. da-server에 전체 데이터 저장
    comm, err := l.PlasmaDA.SetInput(ctx, data)
    if err != nil {
        l.Log.Error("Failed to post input to Plasma DA", "error", err)
        return nil
    }

    // 2. L1에는 commitment(해시)만 게시
    data = comm.TxData()  // 작은 크기 (32-64 bytes)

    // 3. L1 트랜잭션 크기 작음 → Gas 비용 절약
}
```

**데이터 흐름**:
```
L2 Batch Data (예: 1MB)
    │
    ├─→ da-server (로컬 파일 시스템에 저장)
    │   - 전체 1MB 데이터 저장
    │   - HTTP 서버로 데이터 제공
    │
    └─→ L1 Ethereum
        - commitment (32 bytes만 게시)
        - 매우 저렴한 Gas 비용
```

**장점**:
- Gas 비용 대폭 절감 (commitment만 게시)
- 빠른 테스트 사이클

**단점**:
- da-server 장애 시 데이터 손실 위험
- 신뢰 가정 필요 (da-server가 정직하게 데이터 제공)
- Challenge 메커니즘 필요
- 프로덕션에 부적합

---

### L1 직접 사용 (프로덕션)

```go
// op-batcher/batcher/driver.go:502, 515
data := txdata.CallData()
// l.Config.UsePlasma == false이므로 Plasma 블록 건너뜀
candidate = l.calldataTxCandidate(data)  // 전체 데이터 포함

// 또는 Blob 사용
candidate = l.blobTxCandidate(txdata)  // Blob으로 전체 데이터 게시
```

**데이터 흐름**:
```
L2 Batch Data (예: 1MB)
    │
    └─→ L1 Ethereum BatchInboxAddress
        - 전체 1MB 데이터 게시
        - Calldata 또는 EIP-4844 Blob 형태
        - L1 블록체인에 영구 저장
        - 누구나 검증 및 접근 가능
```

**장점**:
- ✅ **완전한 데이터 가용성**: Ethereum L1의 보안성 활용
- ✅ **신뢰 최소화**: 추가 신뢰 가정 불필요
- ✅ **검증 가능**: 누구나 언제든지 데이터 접근 및 검증 가능
- ✅ **단순성**: 외부 da-server 인프라 불필요
- ✅ **영구 보존**: L1에 영구 저장

**단점**:
- 높은 Gas 비용 (전체 데이터 게시)
- EIP-4844 Blob 사용으로 비용 절감 가능

---

## 실제 트랜잭션 흐름

### 전체 아키텍처

```
┌──────────────────────────────────────────────────────┐
│ op-batcher (L2 Batch Submitter)                      │
├──────────────────────────────────────────────────────┤
│                                                       │
│ STEP 1: L2 블록 수집                                  │
│ ┌─────────────────────────────────────────────────┐  │
│ │ loadBlocksIntoState()                           │  │
│ │ - L2 geth에서 블록 가져오기                     │  │
│ │ - 블록 1000-1100 수집                           │  │
│ └─────────────────────────────────────────────────┘  │
│                   ↓                                   │
│ STEP 2: 배치 데이터 생성                              │
│ ┌─────────────────────────────────────────────────┐  │
│ │ state.TxData()                                  │  │
│ │ - 트랜잭션 압축 (zlib/brotli)                   │  │
│ │ - 배치 인코딩                                    │  │
│ │ - 1MB → 100KB (압축률 ~10:1)                    │  │
│ └─────────────────────────────────────────────────┘  │
│                   ↓                                   │
│ STEP 3: 트랜잭션 후보 생성                            │
│ ┌─────────────────────────────────────────────────┐  │
│ │ sendTransaction()                               │  │
│ │                                                  │  │
│ │ if UsePlasma == false:                          │  │
│ │   ├─ UseBlobs?                                  │  │
│ │   │  ├─ Yes → blobTxCandidate()                │  │
│ │   │  │         (EIP-4844 Blob)                 │  │
│ │   │  └─ No  → calldataTxCandidate()            │  │
│ │   │           (Calldata)                        │  │
│ │   │                                              │  │
│ │   └─ To: BatchInboxAddress                      │  │
│ │      Data: [압축된 100KB 배치 데이터]           │  │
│ └─────────────────────────────────────────────────┘  │
│                   ↓                                   │
│ STEP 4: L1 트랜잭션 전송                              │
│ ┌─────────────────────────────────────────────────┐  │
│ │ queue.Send()                                     │  │
│ │ - TxManager를 통해 전송                          │  │
│ │ - Gas 추정 및 논스 관리                          │  │
│ └─────────────────────────────────────────────────┘  │
└───────────────────────────┬──────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────┐
│ Ethereum L1 Network                                  │
├──────────────────────────────────────────────────────┤
│                                                       │
│ Transaction Details:                                 │
│ ┌────────────────────────────────────────────────┐   │
│ │ From:  0xBatcher...                            │   │
│ │ To:    0xff00...0010 (BatchInboxAddress)       │   │
│ │ Type:  2 (EIP-1559) or 3 (EIP-4844 Blob)      │   │
│ │ Gas:   ~500,000                                │   │
│ │ Data:  [100KB compressed batch data]           │   │
│ │                                                 │   │
│ │ Calldata/Blob contains:                        │   │
│ │ - L2 Block 1000-1100 transactions              │   │
│ │ - Block headers                                │   │
│ │ - State transition data                        │   │
│ │ - Compressed with zlib/brotli                  │   │
│ └────────────────────────────────────────────────┘   │
│                                                       │
│ ✅ Transaction mined in L1 block #5000000            │
│ ✅ Data permanently stored on Ethereum L1            │
│ ✅ Anyone can retrieve and verify this data          │
│                                                       │
└───────────────────────────┬──────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────┐
│ op-node (L2 Derivation)                              │
├──────────────────────────────────────────────────────┤
│                                                       │
│ STEP 1: L1 모니터링                                   │
│ ┌─────────────────────────────────────────────────┐  │
│ │ - L1 블록 #5000000 감지                         │  │
│ │ - BatchInboxAddress로 전송된 트랜잭션 필터링    │  │
│ └─────────────────────────────────────────────────┘  │
│                   ↓                                   │
│ STEP 2: 데이터 추출                                   │
│ ┌─────────────────────────────────────────────────┐  │
│ │ DataFromEVMTransactions()                       │  │
│ │ - Calldata 또는 Blob에서 데이터 추출            │  │
│ │ - 압축 해제 (100KB → 1MB)                       │  │
│ └─────────────────────────────────────────────────┘  │
│                   ↓                                   │
│ STEP 3: L2 블록 재구성                                │
│ ┌─────────────────────────────────────────────────┐  │
│ │ ProcessBlock()                                   │  │
│ │ - L2 블록 1000-1100 재생성                      │  │
│ │ - 트랜잭션 실행                                  │  │
│ │ - 상태 업데이트                                  │  │
│ └─────────────────────────────────────────────────┘  │
│                                                       │
└──────────────────────────────────────────────────────┘
```

---

### 상세 단계별 설명

#### op-batcher 단계

**1. L2 블록 수집** (`driver.go:152-185`)
```go
func (l *BatchSubmitter) loadBlocksIntoState(ctx context.Context) error {
    start, end, err := l.calculateL2BlockRangeToStore(ctx)
    // start.Number = 1000, end.Number = 1100

    for i := start.Number + 1; i < end.Number+1; i++ {
        block, err := l.loadBlockIntoState(ctx, i)
        // L2 블록 1001, 1002, ..., 1100 순차 로드
        l.lastStoredBlock = eth.ToBlockID(block)
    }
    return nil
}
```

**2. 배치 데이터 압축 및 인코딩** (`channel_manager.go`)
```go
// 블록 데이터 → RLP 인코딩 → 압축 → 프레임 분할
compressed_data = compress(encode(blocks))
frames = split_into_frames(compressed_data)
```

**3. 트랜잭션 생성 및 전송** (`driver.go:484`)
```go
candidate := &txmgr.TxCandidate{
    To:     &l.RollupConfig.BatchInboxAddress,
    TxData: compressed_data,  // 100KB
}
queue.Send(txdata.ID(), *candidate, receiptsCh)
```

#### L1 Ethereum 단계

**트랜잭션이 L1 블록에 포함됨**:
```
Block #5000000
├─ Transaction Hash: 0xabc...def
├─ From: 0xBatcher...
├─ To: 0xff00...0010
├─ Input Data: 0x00...ff (100KB compressed batch)
└─ Status: Success
```

#### op-node 단계

**L1에서 데이터 추출 및 L2 재구성** (`data_source.go`)
```go
func DataFromEVMTransactions(cfg DataSourceConfig, batcherAddr common.Address,
    txs types.Transactions, log log.Logger) []eth.Data {

    var out []eth.Data
    for _, tx := range txs {
        // BatchInboxAddress로 전송된 트랜잭션만 필터링
        if tx.To() != nil && *tx.To() == cfg.batchInboxAddress {
            // Calldata에서 배치 데이터 추출
            out = append(out, tx.Data())
        }
    }
    return out
}
```

---

## 코드 위치 참조

### 주요 파일 및 함수

| 컴포넌트 | 파일 경로 | 라인 | 함수/구조체 | 설명 |
|---------|----------|------|------------|------|
| **배치 전송 핵심** | `op-batcher/batcher/driver.go` | 484-528 | `sendTransaction()` | L1 트랜잭션 생성 및 전송 |
| **Plasma 체크** | `op-batcher/batcher/driver.go` | 504-514 | Plasma 분기 로직 | Plasma 모드 여부에 따라 분기 |
| **Calldata TX 생성** | `op-batcher/batcher/driver.go` | 546-552 | `calldataTxCandidate()` | Calldata 트랜잭션 생성 |
| **Blob TX 생성** | `op-batcher/batcher/driver.go` | 530-544 | `blobTxCandidate()` | EIP-4844 Blob 트랜잭션 생성 |
| **BatchInbox 주소** | `op-node/rollup/types.go` | 119 | `Config.BatchInboxAddress` | L1 BatchInbox 컨트랙트 주소 |
| **Config 구조체** | `op-batcher/batcher/config.go` | 1-104 | `CLIConfig`, `BatcherConfig` | Batcher 설정 구조 |
| **데이터 추출** | `op-node/rollup/derive/data_source.go` | 50-56 | `NewDataSourceFactory()` | L1에서 배치 데이터 추출 |
| **Calldata 파싱** | `op-node/rollup/derive/calldata_source.go` | 60-127 | `DataFromEVMTransactions()` | EVM 트랜잭션에서 데이터 추출 |

---

### 설정 플래그

**op-batcher 실행 시 설정**:

```bash
# Plasma 모드 (로컬 devnet)
op-batcher \
  --plasma.enabled=true \
  --plasma.da-server=http://da-server:3100

# L1 직접 사용 (프로덕션)
op-batcher \
  --plasma.enabled=false  # 또는 플래그 생략 (기본값)

# Blob 사용 (프로덕션, EIP-4844)
op-batcher \
  --data-availability-type=blobs
```

**Docker Compose 환경 변수**:

```yaml
# 로컬 devnet
environment:
  OP_BATCHER_PLASMA_ENABLED: 'true'
  OP_BATCHER_PLASMA_DA_SERVER: 'http://da-server:3100'

# 프로덕션
# Plasma 관련 환경 변수 없음 (기본값 false)
```

**Kubernetes Helm Values**:

```yaml
# AWS/K8s 배포에는 Plasma 설정 없음
op_batcher:
  env:
    - name: OP_BATCHER_L1_ETH_RPC
      value: "https://ethereum-rpc.com"
    # Plasma 관련 설정 없음
```

---

## 데이터 가용성 보장 비교

### Plasma 모드

```
┌─────────────────────────────────────┐
│ 데이터 가용성 체인                   │
├─────────────────────────────────────┤
│                                      │
│ L2 Data (1MB)                       │
│      │                               │
│      ├─→ da-server (로컬 저장)      │
│      │   └─ 단일 장애점 (SPOF)     │
│      │                               │
│      └─→ L1 (commitment: 32 bytes)  │
│          └─ 데이터 자체는 없음      │
│                                      │
│ 검증 과정:                           │
│ 1. L1에서 commitment 읽기            │
│ 2. da-server에 데이터 요청           │
│ 3. 해시 검증                         │
│                                      │
│ 위험:                                │
│ - da-server 다운 → 데이터 손실      │
│ - 악의적 da-server → 잘못된 데이터  │
│                                      │
└─────────────────────────────────────┘
```

### L1 직접 사용

```
┌─────────────────────────────────────┐
│ 데이터 가용성 체인                   │
├─────────────────────────────────────┤
│                                      │
│ L2 Data (1MB)                       │
│      │                               │
│      └─→ L1 (전체 1MB 게시)         │
│          ├─ Ethereum 검증자 네트워크│
│          ├─ 수천 개 노드에 복제     │
│          └─ 영구 보존               │
│                                      │
│ 검증 과정:                           │
│ 1. L1에서 트랜잭션 calldata 읽기     │
│ 2. 압축 해제                         │
│ 3. L2 블록 재구성                    │
│                                      │
│ 보장:                                │
│ ✅ 단일 장애점 없음                  │
│ ✅ Ethereum 보안성 상속              │
│ ✅ 누구나 검증 가능                  │
│                                      │
└─────────────────────────────────────┘
```

---

## 실전 예시

### 프로덕션 트랜잭션 예시

**L1 Etherscan에서 확인 가능한 실제 트랜잭션**:

```
Transaction: 0xabc123...def789
Block: 5000000
From: 0xBatcher... (Batcher EOA)
To: 0xff00000000000000000000000000000000000010 (BatchInboxAddress)
Value: 0 ETH
Gas Used: 487,234
Gas Price: 50 gwei
Input Data: 0x00...ff (102,400 bytes)

Input Data 내용:
- Version byte: 0x00
- Compressed batch data (zlib):
  - L2 Block 1000-1100
  - 1,234 transactions
  - State diffs
```

**비용 계산**:

```
Calldata 방식:
- Data size: 100KB
- Calldata gas: 16 gas/byte (non-zero)
- Total gas: 100,000 * 16 = 1,600,000 gas
- Cost @ 50 gwei: 0.08 ETH (~$200)

EIP-4844 Blob 방식:
- Blob size: 128KB (최대)
- Blob gas: ~0.01 ETH (~$25)
- 약 87.5% 비용 절감
```

---

## 결론

### 왜 프로덕션에서 L1 직접 사용을 선택하는가?

1. **보안성**
   - Ethereum L1의 검증자 네트워크 활용
   - 51% 공격 저항성 상속
   - 수천 개 노드에 데이터 복제

2. **검증 가능성**
   - 누구나 L1에서 데이터 추출 가능
   - 독립적 검증 가능
   - 신뢰 최소화

3. **단순성**
   - da-server 인프라 불필요
   - 추가 Challenge 메커니즘 불필요
   - 운영 복잡도 감소

4. **영구성**
   - L1에 영구 저장
   - 히스토리 재구성 보장
   - Archive node로 언제든 접근

### 트레이드오프

| 항목 | Plasma | L1 직접 |
|-----|--------|---------|
| **Gas 비용** | 매우 저렴 | 높음 (Blob으로 완화) |
| **보안** | 낮음 (신뢰 가정) | 높음 (L1 보장) |
| **복잡도** | 높음 | 낮음 |
| **용도** | 개발/테스트 | 프로덕션 |

### 코드로 확인하는 방법

```bash
# op-batcher 로그 확인
docker compose logs op-batcher | grep "building.*transaction candidate"

# 로컬 devnet (Plasma)
# 출력: building Calldata transaction candidate size=64

# 프로덕션 (L1 직접)
# 출력: building Calldata transaction candidate size=102400

# L1 트랜잭션 확인
cast tx <tx_hash> --rpc-url $L1_RPC_URL
# To 필드가 BatchInboxAddress인지 확인
# Input data 크기 확인
```

---

**대상 프로젝트**: tokamak-thanos (Optimism Fork)
**코드 베이스**: op-batcher, op-node, op-plasma
