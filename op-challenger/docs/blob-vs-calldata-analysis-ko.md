# TRH-SDK의 데이터 가용성 방식: Blob vs Calldata 분석

## 목차
1. [개요](#개요)
2. [trh-sdk의 기본 설정](#trh-sdk의-기본-설정)
3. [코드 분석](#코드-분석)
4. [Blob vs Calldata 비교](#blob-vs-calldata-비교)
5. [실제 동작 흐름](#실제-동작-흐름)
6. [Blob 활성화 방법](#blob-활성화-방법)
7. [결론](#결론)

---

## 개요

op-batcher는 L2 배치 데이터를 L1에 게시할 때 두 가지 방식을 지원합니다:

1. **Calldata 방식**: 전통적인 트랜잭션 calldata를 통한 데이터 게시
2. **Blob 방식**: EIP-4844 blob 트랜잭션을 통한 데이터 게시 (Ecotone 업그레이드 이후)

이 문서는 **trh-sdk가 어떤 방식을 사용하는지**, 그리고 **그 이유**를 코드 분석을 통해 설명합니다.

---

## trh-sdk의 기본 설정

### 결론부터 말하자면

**trh-sdk는 Calldata 방식을 사용합니다 (Blob 미사용)**

### 근거

#### 1. 기본값 설정 (op-batcher/flags/flags.go:124-133)

```go
DataAvailabilityTypeFlag = &cli.GenericFlag{
    Name: "data-availability-type",
    Value: func() *DataAvailabilityType {
        out := CalldataType  // ⭐️ 기본값: calldata
        return &out
    }(),
    Usage: "The data availability type to use for posting batches " +
           "Options: \"calldata\", \"blobs\" (blobs requires Ecotone)",
    EnvVars: prefixEnvVars("DATA_AVAILABILITY_TYPE"),
    Category: flags.BatcherCategory,
}
```

**핵심 포인트:**
- 기본값이 `CalldataType`로 설정됨
- 환경 변수: `OP_BATCHER_DATA_AVAILABILITY_TYPE`
- Blob 사용을 위해서는 명시적으로 "blobs"로 설정 필요

#### 2. UseBlobs 플래그 설정 (op-batcher/batcher/service.go:204-216)

```go
func (bs *BatcherService) initFromCLIConfig(ctx context.Context,
    cfg *CLIConfig) error {

    // ... (생략) ...

    switch cfg.DataAvailabilityType {
    case flags.BlobsType:
        bs.Log.Info("Using EIP-4844 blobs for data availability")
        bs.UseBlobs = true
    case flags.CalldataType:  // ⭐️ 기본값
        bs.Log.Info("Using calldata for data availability")
        bs.UseBlobs = false
    default:
        return fmt.Errorf("unknown data availability type: %s",
                         cfg.DataAvailabilityType)
    }

    // ... (생략) ...
}
```

**핵심 포인트:**
- `DataAvailabilityType`이 `CalldataType`일 때 `UseBlobs = false`
- 이 설정이 실제 트랜잭션 생성 로직을 제어

#### 3. trh-sdk 배포 설정 확인

**로컬 Devnet (ops-bedrock/docker-compose.yml):**
```yaml
op-batcher:
  environment:
    # ⭐️ DATA_AVAILABILITY_TYPE 환경 변수가 설정되어 있지 않음
    # = 기본값(calldata) 사용
```

**AWS/Kubernetes 배포 (trh-sdk):**
- trh-sdk 코드에서 `OP_BATCHER_DATA_AVAILABILITY_TYPE` 설정을 찾을 수 없음
- 따라서 기본값(calldata) 사용

---

## 코드 분석

### 1. 플래그 타입 정의 (op-batcher/flags/types.go)

```go
// DataAvailabilityType represents the type of data availability
type DataAvailabilityType string

const (
    CalldataType DataAvailabilityType = "calldata"
    BlobsType    DataAvailabilityType = "blobs"
)
```

### 2. 트랜잭션 전송 로직 (op-batcher/batcher/driver.go:484-528)

```go
func (l *BatchSubmitter) sendTransaction(ctx context.Context,
    txdata txData, queue *txmgr.Queue[txID],
    receiptsCh chan txmgr.TxReceipt[txID]) error {

    // ... (생략) ...

    var candidate *txmgr.TxCandidate

    // ⭐️ Blob vs Calldata 결정
    if l.Config.UseBlobs {  // false (trh-sdk)
        // EIP-4844 Blob 트랜잭션 생성
        candidate, err = l.blobTxCandidate(txdata)
        if err != nil {
            return fmt.Errorf("failed to create blob tx candidate: %w", err)
        }
    } else {
        // ⭐️ trh-sdk는 이 경로를 실행
        // Calldata 트랜잭션 생성
        data := txdata.CallData()

        // Plasma 체크 (trh-sdk는 false)
        if l.Config.UsePlasma {
            comm, err := l.PlasmaDA.SetInput(ctx, data)
            if err != nil {
                return fmt.Errorf("failed to set data on plasma da: %w", err)
            }
            data = comm.TxData()
        }

        // ⭐️ 최종적으로 실행되는 코드
        candidate = l.calldataTxCandidate(data)
    }

    // L1에 트랜잭션 전송
    queue.Send(txdata.ID(), *candidate, receiptsCh)
    return nil
}
```

### 3. Calldata 트랜잭션 생성 (op-batcher/batcher/driver.go:546-552)

```go
func (l *BatchSubmitter) calldataTxCandidate(data []byte)
    *txmgr.TxCandidate {
    return &txmgr.TxCandidate{
        To:     &l.RollupConfig.BatchInboxAddress,  // L1 주소
        TxData: data,  // 전체 배치 데이터를 calldata에 포함
    }
}
```

**동작 방식:**
- 전체 배치 데이터를 트랜잭션의 `calldata` 필드에 포함
- `BatchInboxAddress`(L1 컨트랙트)로 전송
- 데이터가 L1에 영구적으로 저장됨

### 4. Blob 트랜잭션 생성 (op-batcher/batcher/driver.go:530-544)

```go
func (l *BatchSubmitter) blobTxCandidate(data txData)
    (*txmgr.TxCandidate, error) {

    var err error
    blobs, err := data.Blobs()
    if err != nil {
        return nil, err
    }

    return &txmgr.TxCandidate{
        To:    &l.RollupConfig.BatchInboxAddress,
        Blobs: blobs,  // EIP-4844 blob 사용
    }, nil
}
```

**동작 방식 (trh-sdk는 사용 안 함):**
- 데이터를 EIP-4844 blob으로 변환
- Blob 트랜잭션으로 L1에 전송
- Calldata보다 저렴하지만 일정 기간 후 삭제됨

---

## Blob vs Calldata 비교

### 특성 비교표

| 항목 | Calldata | Blob (EIP-4844) |
|------|----------|-----------------|
| **trh-sdk 사용 여부** | ✅ 사용 중 | ❌ 미사용 |
| **기본값** | ✅ 기본값 | ❌ 명시적 설정 필요 |
| **필수 업그레이드** | 없음 | Ecotone 하드포크 필요 |
| **가스 비용** | 높음 (16 gas/byte) | 낮음 (~1 gas/byte) |
| **데이터 보관 기간** | 영구 | 임시 (~18일) |
| **데이터 위치** | Transaction calldata | Blob sidecar |
| **L1 저장소** | 영구 블록체인 저장 | 임시 blob 저장 |
| **검증 가능 기간** | 무제한 | ~18일 (이후 재구성 불가) |
| **적합한 용도** | 장기 데이터 보관 필요 | 단기 검증 후 폐기 가능 |

### 비용 예시

**가정:**
- 배치 크기: 100KB
- L1 gas price: 50 gwei

**Calldata 방식:**
```
Gas 사용량 = 100,000 bytes × 16 gas/byte = 1,600,000 gas
비용 = 1,600,000 × 50 gwei = 0.08 ETH (~$200)
```

**Blob 방식 (미사용):**
```
Gas 사용량 = 100,000 bytes × ~1 gas/byte = ~100,000 gas
비용 = 100,000 × 50 gwei = 0.005 ETH (~$12.5)
```

⚠️ **trh-sdk는 Calldata를 사용하므로 더 높은 비용이 발생하지만, 데이터 영구 보관이 보장됩니다.**

---

## 실제 동작 흐름

### trh-sdk의 실제 실행 경로

```
┌─────────────────────────────────────────────────────────────┐
│                    BatchSubmitter 시작                        │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│          initFromCLIConfig() - 설정 로드                      │
│                                                               │
│  DataAvailabilityType = CalldataType (기본값)                │
│  UseBlobs = false ⭐️                                         │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│              sendTransaction() - TX 생성                      │
│                                                               │
│  if l.Config.UseBlobs {        // false ⭐️                  │
│      candidate = blobTxCandidate()    // 실행 안 됨          │
│  } else {                                                     │
│      data = txdata.CallData()         // ⭐️ 실행됨          │
│                                                               │
│      if l.Config.UsePlasma {          // false               │
│          data = PlasmaDA.SetInput()   // 실행 안 됨          │
│      }                                                        │
│                                                               │
│      candidate = calldataTxCandidate(data) // ⭐️ 실행됨     │
│  }                                                            │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│         calldataTxCandidate() - Calldata TX 생성             │
│                                                               │
│  TxCandidate {                                                │
│      To: BatchInboxAddress (L1 컨트랙트)                     │
│      TxData: [전체 배치 데이터] ⭐️                           │
│  }                                                            │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│              queue.Send() - L1에 전송                         │
│                                                               │
│  → Transaction calldata에 전체 데이터 포함                    │
│  → L1 BatchInboxAddress로 전송                               │
│  → 데이터가 L1 블록체인에 영구 저장 ⭐️                       │
└─────────────────────────────────────────────────────────────┘
```

### 코드 실행 요약

1. **설정 단계**: `UseBlobs = false` (기본값)
2. **분기 선택**: `else` 블록 실행 (Calldata 경로)
3. **Plasma 체크**: 건너뜀 (`UsePlasma = false`)
4. **TX 생성**: `calldataTxCandidate()` 호출
5. **데이터 위치**: Transaction의 `calldata` 필드
6. **최종 결과**: L1에 영구 저장

---

## Blob 활성화 방법

### trh-sdk에서 Blob을 사용하려면 (현재 미사용)

⚠️ **주의: Ecotone 하드포크가 활성화된 네트워크에서만 사용 가능합니다.**

#### 1. 로컬 Devnet에서 Blob 활성화

**ops-bedrock/docker-compose.yml 수정:**

```yaml
op-batcher:
  environment:
    # 기존 환경 변수들...

    # Blob 활성화 추가
    OP_BATCHER_DATA_AVAILABILITY_TYPE: "blobs"  # ⭐️ 추가
```

#### 2. AWS/Kubernetes 배포에서 Blob 활성화

**trh-sdk에서 환경 변수 추가 필요 (코드 수정 필요):**

op-batcher deployment에 환경 변수 추가:
```yaml
env:
  - name: OP_BATCHER_DATA_AVAILABILITY_TYPE
    value: "blobs"
```

#### 3. 직접 CLI로 실행 시

```bash
op-batcher \
  --data-availability-type=blobs \
  # ... 기타 옵션들 ...
```

### Blob 사용 시 주의사항

1. **Ecotone 하드포크 필수**
   - L1과 L2 모두 Ecotone 업그레이드 완료 필요
   - 업그레이드 전에는 사용 불가

2. **데이터 보관 기간 제한**
   - Blob 데이터는 약 18일 후 L1에서 삭제됨
   - 장기 보관이 필요한 경우 별도 아카이브 필요

3. **재구성(Reorg) 제한**
   - Blob 데이터 삭제 후에는 체인 재구성 불가
   - 장기적인 체인 검증에 제약

4. **비용 절감 vs 영구성 트레이드오프**
   - 비용: Blob이 ~16배 저렴
   - 영구성: Calldata가 무제한 보관

---

## 결론

### 현재 상태 (trh-sdk)

| 항목 | 설정값 | 근거 |
|------|--------|------|
| **데이터 가용성 방식** | Calldata | `DataAvailabilityTypeFlag` 기본값 |
| **UseBlobs 플래그** | `false` | 명시적 설정 없음 |
| **실행 경로** | `calldataTxCandidate()` | driver.go:546-552 |
| **데이터 위치** | Transaction calldata | L1 블록체인에 영구 저장 |
| **비용** | 높음 (16 gas/byte) | Calldata 비용 |
| **보관 기간** | 영구 | L1 블록체인 특성 |

### 핵심 코드 경로

```go
// 1. 설정 (service.go:204-216)
cfg.DataAvailabilityType = CalldataType  // 기본값
bs.UseBlobs = false  // ⭐️

// 2. 트랜잭션 생성 (driver.go:484-528)
if l.Config.UseBlobs {  // false ⭐️
    // 실행 안 됨
} else {
    // ⭐️ 실제 실행
    data := txdata.CallData()
    candidate = l.calldataTxCandidate(data)  // ⭐️
}

// 3. Calldata TX 생성 (driver.go:546-552)
TxCandidate {
    To: BatchInboxAddress,
    TxData: data,  // 전체 배치 데이터 ⭐️
}
```

### 왜 trh-sdk는 Calldata를 사용하는가?

1. **안정성 우선**
   - Calldata는 검증된 방식
   - 데이터 영구 보관 보장

2. **호환성**
   - Ecotone 이전 네트워크와도 호환
   - 별도 하드포크 불필요

3. **장기 검증 가능**
   - 무제한 체인 재구성 가능
   - 영구적인 데이터 검증

4. **기본 설정 유지**
   - 명시적 설정 없이 동작
   - 운영 복잡도 감소

### 향후 고려사항

**Blob 전환을 고려해야 하는 경우:**
- L1 가스 비용이 매우 높을 때
- Ecotone 업그레이드가 안정적으로 완료된 후
- 단기 데이터 검증만으로 충분한 경우
- 별도 데이터 아카이브 시스템 구축 가능한 경우

**현재 Calldata를 유지해야 하는 이유:**
- 데이터 영구 보관이 중요
- 장기적인 체인 재구성 가능성 필요
- 운영 안정성 우선
- 검증된 방식 선호

---

## 참고 자료

### 관련 파일

1. **플래그 정의**
   - `op-batcher/flags/flags.go:124-133` - DataAvailabilityTypeFlag 정의
   - `op-batcher/flags/types.go` - CalldataType, BlobsType 상수

2. **설정 로직**
   - `op-batcher/batcher/service.go:204-216` - UseBlobs 설정

3. **트랜잭션 생성**
   - `op-batcher/batcher/driver.go:484-528` - sendTransaction()
   - `op-batcher/batcher/driver.go:530-544` - blobTxCandidate()
   - `op-batcher/batcher/driver.go:546-552` - calldataTxCandidate()

4. **배포 설정**
   - `ops-bedrock/docker-compose.yml` - 로컬 devnet 설정
   - `trh-sdk` - AWS/Kubernetes 배포 설정

### 관련 문서

- [EIP-4844: Shard Blob Transactions](https://eips.ethereum.org/EIPS/eip-4844)
- [Optimism Ecotone Upgrade](https://docs.optimism.io/builders/node-operators/network-upgrades#ecotone)
- [op-batcher 문서](https://github.com/ethereum-optimism/optimism/tree/develop/op-batcher)

---

**기반 코드**: tokamak-thanos (commit: ef63c0e65)
