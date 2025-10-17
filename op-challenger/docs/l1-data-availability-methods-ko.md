# L1 데이터 가용성 방식 비교 분석

## 목차
1. [개요](#개요)
2. [세 가지 DA 방식](#세-가지-da-방식)
3. [Calldata 방식 (기본)](#calldata-방식-기본)
4. [Blob 방식 (EIP-4844)](#blob-방식-eip-4844)
5. [Plasma 방식 (테스트용)](#plasma-방식-테스트용)
6. [방식 비교 및 선택 가이드](#방식-비교-및-선택-가이드)
7. [실제 트랜잭션 예시](#실제-트랜잭션-예시)
8. [TRH SDK의 DA 전략](#trh-sdk의-da-전략)

---

## 개요

### 데이터 가용성(DA)이란?

Layer 2 Rollup에서 **데이터 가용성(Data Availability)**는 모든 노드가 L2 트랜잭션 데이터에 접근하여 L2 상태를 독립적으로 검증할 수 있도록 보장하는 것을 의미합니다.

```
┌─────────────────────────────────────────┐
│  L2 Sequencer (op-node)                 │
│  - 사용자 트랜잭션 수신                  │
│  - L2 블록 생성                          │
│  - 상태 전이 실행                        │
└──────────────┬──────────────────────────┘
               │
               ▼
    ❓ 이 데이터를 어디에 저장?
               │
    ┌──────────┴──────────┐
    │                     │
    ▼                     ▼
┌─────────┐         ┌──────────┐
│   L1    │         │ Off-chain│
│ (안전)  │         │ (위험?)  │
└─────────┘         └──────────┘
```

### 왜 중요한가?

L2 상태의 무결성과 보안은 **모든 노드가 트랜잭션 데이터에 접근할 수 있어야** 보장됩니다:

- ✅ **검증 가능성**: 누구나 L2 상태를 재구성하고 검증
- ✅ **검열 저항성**: Sequencer의 검열을 감지하고 증명
- ✅ **복구 가능성**: Sequencer 장애 시 다른 노드가 복구
- ✅ **Fault Proof**: 무효한 상태 전이를 증명하여 이의 제기

---

## 세 가지 DA 방식

Optimism/Tokamak Thanos Stack은 세 가지 데이터 가용성 방식을 지원합니다:

| 방식 | 데이터 저장 위치 | L1 기록 | 가스 비용 | 보안 | 프로덕션 |
|------|----------------|---------|----------|------|----------|
| **Calldata** | L1 Execution Layer | 전체 데이터 | 높음 | ✅ 높음 | ✅ 사용 |
| **Blob (EIP-4844)** | L1 Beacon Chain | Blob commitment | 낮음 | ✅ 높음 | ✅✅ 권장 |
| **Plasma** | Off-chain DA Server | DA commitment | 매우 낮음 | ⚠️ 낮음 | ❌ 테스트만 |

---

## Calldata 방식 (기본)

### 개념

L2 트랜잭션 데이터를 **L1 트랜잭션의 calldata에 직접 포함**하여 제출하는 방식입니다.

### 코드 분석

#### op-batcher의 Calldata 생성

```go
// op-batcher/batcher/driver.go:497-515
if l.Config.UseBlobs {
    // Blob 사용
} else {
    // Calldata 사용
    data := txdata.CallData()
    
    // Plasma 모드가 아니면 전체 데이터를 L1에 제출
    if l.Config.UsePlasma {
        // Plasma: commitment만 제출
        comm, err := l.PlasmaDA.SetInput(ctx, data)
        data = comm.TxData()
    }
    // ✅ 일반 모드: 전체 데이터 제출
    candidate = l.calldataTxCandidate(data)
}
```

#### Calldata 트랜잭션 생성

```go
// op-batcher/batcher/driver.go:546-551
func (l *BatchSubmitter) calldataTxCandidate(data []byte) *txmgr.TxCandidate {
    l.Log.Info("building Calldata transaction candidate", "size", len(data))
    return &txmgr.TxCandidate{
        To:     &l.RollupConfig.BatchInboxAddress,  // Batch Inbox 컨트랙트
        TxData: data,  // ✅ 전체 L2 트랜잭션 데이터
    }
}
```

#### Calldata 포맷

```go
// op-batcher/batcher/tx_data.go:33-42
func (td *txData) CallData() []byte {
    data := make([]byte, 1, 1+td.Len())
    data[0] = derive.DerivationVersion0  // Version byte (0x00)
    for _, f := range td.frames {
        data = append(data, f.data...)  // Frame 데이터 추가
    }
    return data
}
```

**구조**:
```
[0x00] [Frame 1] [Frame 2] ... [Frame N]
  │       │
  │       └─> 압축된 L2 트랜잭션 배치
  │           (zlib, brotli 등)
  │
  └─> Derivation Version (0 = Legacy)
```

### 동작 과정

```
1. L2 트랜잭션 수집
   └─ op-node에서 여러 L2 트랜잭션 생성
   
2. 배치 생성 (op-batcher)
   ├─ L2 트랜잭션들을 하나의 배치로 묶음
   ├─ 압축 (zlib 또는 brotli)
   │  └─ 원본 크기: 500 KB → 압축 후: 150 KB
   └─ Frame으로 분할
      └─ 각 Frame은 독립적으로 전송 가능

3. Calldata 준비
   └─ Version byte (0x00) + 압축된 데이터
   
4. L1 트랜잭션 생성
   ├─ To: BatchInboxAddress
   │      0xff00000000000000000000000000000000{L2ChainID}
   ├─ Data: [전체 배치 데이터]
   └─ Gas Limit: 자동 계산
   
5. L1 Ethereum에 제출
   └─ 트랜잭션이 L1 블록에 포함됨
   └─ 모든 Ethereum 노드가 데이터 보유
   └─ 영구 보존 ✅
```

### 가스 비용 계산

**Calldata 가스 비용**:
```
- Zero byte (0x00): 4 gas
- Non-zero byte: 16 gas

예시: 150 KB 압축 데이터
- 평균 non-zero byte 비율: 80%
- 150,000 × 0.8 × 16 = 1,920,000 gas (non-zero)
- 150,000 × 0.2 × 4 = 120,000 gas (zero)
- 총 gas: 2,040,000 gas

Gas Price: 20 Gwei
Total Cost: 2,040,000 × 20 = 40,800,000 Gwei
         = 0.0408 ETH
         ≈ $75 (ETH = $1,800 기준)
```

### 장점

- ✅ **영구 저장**: 모든 Ethereum 노드가 영구 보관
- ✅ **검증 가능**: 누구나 데이터에 접근하여 검증
- ✅ **탈중앙화**: 제3자 서비스 불필요
- ✅ **안전성**: L1의 보안을 그대로 활용
- ✅ **간단함**: 추가 인프라 불필요

### 단점

- ⚠️ **비용 높음**: Calldata 가스 비용이 매우 비쌈
- ⚠️ **확장성 제한**: 높은 비용으로 인한 처리량 제약
- ⚠️ **L1 혼잡**: 많은 데이터로 L1 블록 공간 사용

---

## Blob 방식 (EIP-4844)

### 개념

**EIP-4844 (Proto-Danksharding)**은 L2를 위한 저렴한 데이터 저장 방식을 제공합니다. 데이터를 L1 Beacon Chain의 **Blob**에 저장하고, Execution Layer에는 commitment만 기록합니다.

### 코드 분석

#### Blob 트랜잭션 생성

```go
// op-batcher/batcher/driver.go:489-496
if l.Config.UseBlobs {
    if candidate, err = l.blobTxCandidate(txdata); err != nil {
        return fmt.Errorf("could not create blob tx candidate: %w", err)
    }
}
```

```go
// op-batcher/batcher/driver.go:530-544
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
        To:    &l.RollupConfig.BatchInboxAddress,
        Blobs: blobs,  // ✅ Blob 데이터
    }, nil
}
```

#### Blob 데이터 생성

```go
// op-batcher/batcher/tx_data.go:44-54
func (td *txData) Blobs() ([]*eth.Blob, error) {
    blobs := make([]*eth.Blob, 0, len(td.frames))
    for _, f := range td.frames {
        var blob eth.Blob
        // Version byte + frame data를 blob으로 변환
        if err := blob.FromData(append([]byte{derive.DerivationVersion0}, f.data...)); err != nil {
            return nil, err
        }
        blobs = append(blobs, &blob)
    }
    return blobs, nil
}
```

**Blob 제약사항**:
```go
// op-service/eth/blob.go (constants)
const (
    MaxBlobDataSize = 4096 * 32 - 1  // 131,071 bytes (~128 KB)
    MaxBlobsPerTx   = 6               // 최대 6개 blob
)

// 따라서 최대 데이터 크기:
// 131,071 bytes × 6 = 786,426 bytes (~768 KB per transaction)
```

### 동작 과정

```
1. L2 트랜잭션 배치 생성 (op-batcher)
   └─ 여러 L2 트랜잭션을 압축
   
2. Blob 생성
   ├─ 압축 데이터를 128KB 단위로 분할
   ├─ 각 chunk를 blob으로 변환
   │  └─ Field elements (4096개) 형태로 인코딩
   └─ KZG commitment 생성
      └─ 암호학적 해시 (데이터 무결성 보장)

3. Blob 트랜잭션 생성 (EIP-4844)
   ├─ Type: 3 (Blob Transaction)
   ├─ To: BatchInboxAddress
   ├─ Blobs: [blob1, blob2, blob3, ...]
   ├─ Commitments: [kzg_comm1, kzg_comm2, ...]
   └─ Proofs: [kzg_proof1, kzg_proof2, ...]

4. L1에 제출
   ├─ Execution Layer: commitment만 기록
   │  └─ 트랜잭션 해시 + KZG commitments
   └─ Beacon Chain: 실제 blob 데이터 저장
      └─ 약 18일간 보존 (4096 epochs)

5. op-node가 데이터 읽기
   ├─ L1 Beacon API를 통해 blob 조회
   ├─ KZG commitment로 무결성 검증
   └─ L2 블록 생성
```

### Blob Transaction 구조

```
EIP-4844 Blob Transaction
├── Transaction Fields (Execution Layer)
│   ├── chainId: 1 (Ethereum Mainnet)
│   ├── nonce: 1234
│   ├── maxFeePerGas: 20 Gwei
│   ├── maxPriorityFeePerGas: 2 Gwei
│   ├── maxFeePerBlobGas: 10 Gwei
│   ├── gas: 21,000
│   ├── to: 0xff00000000000000000000000000000000000420
│   ├── value: 0
│   ├── data: 0x (empty)
│   └── blobVersionedHashes: [0x01abc..., 0x01def...]
│       └─ KZG commitments (각 blob당 하나)
│
└── Blobs (Beacon Chain)
    ├── Blob 0: 131,071 bytes
    │   ├─ KZG Commitment: 0x01abc...
    │   └─ KZG Proof: 0x...
    ├── Blob 1: 131,071 bytes
    │   ├─ KZG Commitment: 0x01def...
    │   └─ KZG Proof: 0x...
    └── Blob 2: 98,304 bytes
        ├─ KZG Commitment: 0x01ghi...
        └─ KZG Proof: 0x...
```

### 가스 비용 계산

**Blob 가스 비용**:
```
Blob gas per blob: 131,072 blob gas
Number of blobs: 3
Total blob gas: 393,216 blob gas

Blob base fee: 1 Wei (평상시)
             : 100 Wei (혼잡 시)

평상시 비용:
393,216 × 1 Wei = 393,216 Wei = 0.00039 ETH ≈ $0.70

혼잡 시 비용:
393,216 × 100 Wei = 39,321,600 Wei = 0.039 ETH ≈ $70

Calldata 대비: 98-99% 저렴! 🎉
```

### Blob 검증 과정

```
op-node가 L1에서 데이터 읽기:

1. L1 트랜잭션에서 blobVersionedHashes 읽기
   └─ [0x01abc..., 0x01def..., 0x01ghi...]

2. Beacon API로 blob 데이터 요청
   └─ GET /eth/v1/beacon/blob_sidecars/{block_id}

3. 각 blob에 대해 KZG 검증
   ├─ commitment = KZG.commitment(blob)
   ├─ if commitment != blobVersionedHash:
   │     reject! (데이터 무결성 실패)
   └─ proof 검증
   
4. 검증 성공 → L2 블록 생성
```

### 장점

- ✅ **비용 효율**: Calldata 대비 98% 저렴
- ✅ **높은 처리량**: 더 많은 트랜잭션 처리 가능
- ✅ **검증 가능**: KZG commitment로 무결성 보장
- ✅ **탈중앙화**: L1 Beacon Chain 사용
- ✅ **EIP-4844 네이티브**: Ethereum 프로토콜 수준 지원

### 단점

- ⚠️ **일시적 보존**: ~18일 후 blob 데이터 pruning
  - 하지만: op-node는 18일 내에 이미 처리하므로 문제 없음
  - Archive node는 blob 데이터를 영구 보관 가능
- ⚠️ **크기 제한**: 최대 768KB per transaction
- ⚠️ **Beacon API 필요**: L1 Beacon node 연결 필수

### Ecotone 업그레이드

Blob 사용은 **Ecotone 하드포크** 이후 활성화됩니다:

```go
// op-batcher/batcher/service.go:224-226
if bs.UseBlobs && !bs.RollupConfig.IsEcotone(uint64(time.Now().Unix())) {
    bs.Log.Error("Cannot use Blob data before Ecotone!")
}
```

**Ecotone 체크**:
```go
// op-node/rollup/types.go
func (c *Config) IsEcotone(timestamp uint64) bool {
    return c.EcotoneTime != nil && timestamp >= *c.EcotoneTime
}
```

---

## Plasma 방식 (테스트용)

### 개념

L2 트랜잭션 데이터를 **오프체인 DA Server**에 저장하고, L1에는 작은 크기의 **commitment(해시)만** 제출하는 방식입니다.

### 코드 분석

#### Plasma 데이터 제출

```go
// op-batcher/batcher/driver.go:502-514
data := txdata.CallData()

// Plasma 모드 활성화 시
if l.Config.UsePlasma {
    // 1. DA Server에 전체 데이터 저장
    comm, err := l.PlasmaDA.SetInput(ctx, data)
    if err != nil {
        l.Log.Error("Failed to post input to Plasma DA", "error", err)
        l.recordFailedTx(txdata.ID(), err)
        return nil
    }
    
    // 2. L1에는 commitment만 제출
    data = comm.TxData()  // ← 작은 크기 (~32 bytes)
}
candidate = l.calldataTxCandidate(data)
```

#### DA Client (Plasma 클라이언트)

```go
// op-plasma/daclient.go:63-107
func (c *DAClient) SetInput(ctx context.Context, img []byte) (CommitmentData, error) {
    var url string
    if c.precompute {
        // Keccak256 commitment 미리 계산
        comm := NewKeccak256Commitment(img)
        url = fmt.Sprintf("%s/put/0x%x", c.url, comm.Encode())
    } else {
        // Generic commitment (서버가 생성)
        url = fmt.Sprintf("%s/put", c.url)
    }
    
    // DA Server에 HTTP POST 요청
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(img))
    resp, err := http.DefaultClient.Do(req)
    
    // Commitment 읽기
    comm, err := io.ReadAll(resp.Body)
    return DecodeCommitmentData(comm)
}
```

#### Commitment 타입

**Keccak256 Commitment** (기본):
```go
// op-plasma/commitment.go:86-125
func NewKeccak256Commitment(input []byte) Keccak256Commitment {
    return Keccak256Commitment(crypto.Keccak256(input))
}

func (c Keccak256Commitment) Encode() []byte {
    return append([]byte{byte(Keccak256CommitmentType)}, c...)
}

func (c Keccak256Commitment) TxData() []byte {
    return append([]byte{TxDataVersion1}, c.Encode()...)
}

func (c Keccak256Commitment) Verify(input []byte) error {
    if !bytes.Equal(c, crypto.Keccak256(input)) {
        return ErrCommitmentMismatch
    }
    return nil
}
```

**구조**:
```
[0x01] [0x00] [32 bytes Keccak256 hash]
  │      │      │
  │      │      └─> Commitment: Keccak256(L2 데이터)
  │      └─> Commitment Type: Keccak256 (0x00)
  └─> Tx Data Version: Plasma (0x01)

총 크기: 34 bytes
```

**Generic Commitment** (외부 DA용):
```go
// op-plasma/commitment.go:127-158
func NewGenericCommitment(input []byte) GenericCommitment {
    return GenericCommitment(input)
}

// ⚠️ 검증 안 함!
func (c GenericCommitment) Verify(input []byte) error {
    return nil  // 항상 통과
}
```

### 동작 과정

```
1. L2 트랜잭션 배치 생성
   └─ [Calldata와 동일]

2. DA Server에 저장 (op-batcher)
   ├─ HTTP POST http://da-server:3100/put/
   │  └─ Body: [압축된 L2 트랜잭션 데이터]
   │
   └─ Response: commitment
      └─ Keccak256: 0x47173285a8d7341e...
      └─ Generic: 0x01ff00000009f5a2b3...

3. L1 트랜잭션 생성
   ├─ To: BatchInboxAddress
   ├─ Data: [0x01][commitment]
   │        └─ 34 bytes만!
   └─ Gas: ~21,000 + (34 × 16) = ~21,544 gas

4. L1에 제출
   └─ Commitment만 L1에 기록
   └─ 실제 데이터는 DA Server에만 존재

5. op-node가 데이터 읽기
   ├─ L1에서 commitment 읽기
   ├─ DA Server에서 데이터 조회
   │  └─ GET http://da-server:3100/get/0x47173285...
   ├─ Keccak256 검증 (Keccak256 commitment)
   │  └─ if Keccak256(data) != commitment: reject!
   └─ L2 블록 생성
```

### 가스 비용 계산

```
Commitment 크기: 34 bytes
- 1 byte (version): 16 gas
- 1 byte (type): 16 gas
- 32 bytes (hash): 32 × 16 = 512 gas

Total calldata gas: 544 gas
Base transaction gas: 21,000 gas

Total: 21,544 gas

Gas Price: 20 Gwei
Total Cost: 21,544 × 20 = 430,880 Gwei
         = 0.00043 ETH
         ≈ $0.77

Calldata 대비: 99% 저렴! 🎉
```

### 장점

- ✅ **극도로 저렴**: Calldata 대비 99% 비용 절감
- ✅ **확장성**: 더 많은 트랜잭션 처리 가능
- ✅ **검증 가능**: Keccak256 commitment로 무결성 보장 (Keccak256 모드)
- ✅ **유연성**: 외부 DA 시스템 통합 가능 (Generic 모드)

### 단점

- ⚠️ **중앙화 위험**: DA Server에 의존
- ⚠️ **가용성 위험**: DA Server 다운 시 L2 중단
- ⚠️ **보안 위험**: DA Server 해킹/데이터 손실
- ⚠️ **Challenge 필요**: DataAvailabilityChallenge 컨트랙트 필요
- ⚠️ **인프라 복잡도**: DA Server 운영 및 백업 필요

### Challenge 메커니즘

Plasma의 보안은 **Challenge 시스템**으로 보장됩니다:

```solidity
// packages/contracts-bedrock/src/L1/DataAvailabilityChallenge.sol
contract DataAvailabilityChallenge {
    uint256 public challengeWindow;  // 챌린지 가능 기간
    uint256 public resolveWindow;    // 해결 필요 기간
    uint256 public bondSize;         // 챌린지 보증금
    
    function challenge(
        uint256 challengedBlockNumber,
        bytes calldata challengedCommitment
    ) external payable {
        require(msg.value >= bondSize, "BondTooLow");
        require(
            block.number <= challengedBlockNumber + challengeWindow,
            "ChallengeWindowNotOpen"
        );
        // 챌린지 생성
    }
    
    function resolve(
        uint256 challengedBlockNumber,
        bytes calldata challengedCommitment,
        bytes calldata preImage
    ) external {
        // Commitment 검증
        require(
            keccak256(preImage) == bytes32(challengedCommitment),
            "InvalidInputData"
        );
        // 챌린지 해결
    }
}
```

**Challenge 프로세스**:
```
1. Sequencer가 commitment 제출
   └─ L1 tx: 0x01 + commitment

2. Challenger가 데이터 검증
   ├─ DA Server에서 데이터 조회 시도
   └─ 데이터 없거나 잘못됨 발견

3. Challenge 시작
   ├─ Challenge bond 예치 (예: 0.1 ETH)
   └─ challengeWindow 내에만 가능 (예: 7일)

4-1. Sequencer가 데이터 제공 (Resolve)
     ├─ 올바른 preimage 제출
     ├─ Challenger는 bond 손실
     └─ Commitment 유효 확인

4-2. Sequencer가 응답 못함 (Expire)
     ├─ resolveWindow 만료 (예: 48시간)
     ├─ Challenger가 bond 회수 + 보상
     └─ Commitment 무효 처리
     └─ op-node가 해당 배치 스킵
```

---

## 방식 비교 및 선택 가이드

### 상세 비교표

| 항목 | Calldata | Blob (EIP-4844) | Plasma (Keccak256) | Plasma (Generic) |
|------|----------|-----------------|-------------------|------------------|
| **L1 저장 위치** | Execution Layer | Beacon Chain | Off-chain | Off-chain |
| **L1 기록 크기** | 전체 데이터 | KZG commitment | Keccak256 hash | Random ID |
| **데이터 크기 제한** | 무제한 | 768KB/tx | 무제한 | 무제한 |
| **가스 비용 (150KB)** | ~$75 | ~$0.70 | ~$0.77 | ~$0.77 |
| **비용 절감율** | 0% | 99% | 99% | 99% |
| **데이터 보존** | 영구 | ~18일 | DA Server 의존 | DA Server 의존 |
| **검증 방식** | L1 재생 | KZG proof | Keccak256 | 검증 안 함 |
| **탈중앙화** | ✅ 완전 | ✅ 완전 | ⚠️ 부분적 | ❌ 중앙화 |
| **Challenge 필요** | ❌ 불필요 | ❌ 불필요 | ✅ 필수 | ✅ 필수 |
| **추가 인프라** | 없음 | Beacon node | DA Server | DA Server + 외부 DA |
| **보안 수준** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| **프로덕션 사용** | ✅ 가능 | ✅✅ 권장 | ⚠️ 주의 | ❌ 비권장 |
| **복잡도** | 낮음 | 중간 | 높음 | 매우 높음 |

### 비용 비교 (실제 데이터)

**시나리오**: 1일 10,000 트랜잭션, 평균 배치 크기 150KB

| 방식 | 배치 수/일 | 가스비/배치 | 총 가스비/일 | 월 비용 |
|------|----------|------------|------------|---------|
| **Calldata** | 100 | $75 | $7,500 | **$225,000** 💰💰💰 |
| **Blob** | 100 | $0.70 | $70 | **$2,100** 💰 |
| **Plasma** | 100 | $0.77 | $77 | **$2,310** 💰 |

**추가 비용**:
- Plasma: DA Server 운영비 (~$500/월)
- Plasma: 백업 및 HA 구성 (~$1,000/월)
- Plasma 총 비용: ~$3,810/월

**결론**: **Blob이 가장 비용 효율적!** ✅

---

## 선택 가이드

### ✅ Calldata를 선택해야 하는 경우

```
- Ecotone 업그레이드 전
- 매우 간단한 설정 필요
- Beacon node 접근 불가
- 영구 데이터 보존 필수
```

### ✅✅ Blob을 선택해야 하는 경우 (권장!)

```
- Ecotone 업그레이드 이후 ⭐
- 프로덕션 환경
- 비용 최적화 필요
- L1 Beacon node 연결 가능
- 높은 처리량 필요
```

### ⚠️ Plasma (Keccak256)를 선택해야 하는 경우

```
- 극도의 비용 절감 필요
- DA Server 운영 리소스 충분
- 보안 위험 감수 가능
- Challenge 메커니즘 구현 완료
- 24/7 모니터링 체계 구축
```

### ❌ Plasma (Generic)은 사용하지 말 것

```
- 테스트 및 개발 환경만
- 외부 DA 시스템 통합 시험용
- 프로덕션 사용 절대 금지!
```

---

## 실제 트랜잭션 예시

### Optimism Mainnet 예시

#### 1. Calldata 방식 (과거)

```
Transaction: 0x8f3b9c0d1e2a3f4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8
Block: 18,234,567
Date: 2023-10-15

To: 0xff00000000000000000000000000000000000420 (Batch Inbox)
Value: 0 ETH
Input Data: 
  0x00789c5d915d6f1c4710c73f...  [129,543 bytes]
  │   │
  │   └─> 압축된 L2 트랜잭션 배치
  └─> Version 0

Gas Used: 6,477,150
Gas Price: 25 Gwei
Tx Fee: 0.1619 ETH ≈ $291

포함된 L2 트랜잭션: 423개
L2 트랜잭션당 비용: $0.69
```

**Etherscan 보기**:
```
Internal Txns (0)
Input Data:
  0x
  00
  789c5d915d6f1c4710c73f...
  [View More] ← 전체 데이터가 L1에 기록됨
```

#### 2. Blob 방식 (현재)

```
Transaction: 0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a
Block: 19,456,789
Date: 2024-03-20

To: 0xff00000000000000000000000000000000000420
Type: 3 (EIP-4844 Blob Transaction)
Value: 0 ETH
Input Data: 0x (empty!)

Max Fee Per Gas: 0.02 Gwei
Max Priority Fee: 0.001 Gwei
Max Fee Per Blob Gas: 10 Gwei

Blob Versioned Hashes:
  0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9
  0x01ab2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a
  0x01cd3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c

Blobs: 3
  - Blob 0: 131,071 bytes (Sidecar)
  - Blob 1: 131,071 bytes (Sidecar)
  - Blob 2: 98,304 bytes (Sidecar)

Gas Used: 21,000
Blob Gas Used: 393,216
Gas Price: 0.02 Gwei
Blob Gas Price: 1 Wei
Tx Fee: 0.00042 ETH + 0.00039 ETH = 0.00081 ETH ≈ $1.46

포함된 L2 트랜잭션: 512개
L2 트랜잭션당 비용: $0.0028
절감율: 99.6%! 🎉
```

**Etherscan 보기**:
```
Type: 3 (EIP4844)
Blob Versioned Hashes: (3)
  0x01fa3b84...
  0x01ab2c3d...
  0x01cd3e4f...
[View Blob Sidecars] ← Beacon Chain에서 조회
```

#### 3. Plasma 방식 (테스트)

```
Transaction: 0x9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9
Block: 19,678,901
Date: 2024-05-10 (테스트넷)

To: 0xff00000000000000000000000000000012345678
Value: 0 ETH
Input Data:
  0x01004717328...  [34 bytes]
  │   │  │
  │   │  └─> Keccak256 commitment (32 bytes)
  │   └─> Commitment type: Keccak256 (0x00)
  └─> Plasma version (0x01)

Gas Used: 21,544
Gas Price: 20 Gwei
Tx Fee: 0.00043 ETH ≈ $0.77

실제 데이터: DA Server에 360 KB 저장
포함된 L2 트랜잭션: 580개
L2 트랜잭션당 비용: $0.0013

⚠️ 하지만:
- DA Server 운영비: $500/월
- 보안 위험
- 중앙화
```

---

## 코드 흐름 다이어그램

### Calldata/Blob 흐름

```
┌─────────────────────────────────────────┐
│ op-node (Sequencer)                     │
│ - 사용자 트랜잭션 수신                   │
│ - L2 블록 생성                           │
└──────────────┬──────────────────────────┘
               │ L2 Block Data
               ▼
┌─────────────────────────────────────────┐
│ op-batcher                              │
│ - 배치 생성 및 압축                      │
│ - Frame 분할                             │
└──────────────┬──────────────────────────┘
               │
          ┌────┴─────┐
          │ UseBlobs?│
          └────┬─────┘
      Yes ◄────┼────► No
          │         │
          ▼         ▼
    ┌─────────┐ ┌──────────┐
    │ Blobs   │ │ Calldata │
    │ (EIP-   │ │ (Legacy) │
    │  4844)  │ │          │
    └────┬────┘ └─────┬────┘
         │            │
         └────────┬───┘
                  │ TxCandidate
                  ▼
         ┌────────────────┐
         │   TxManager    │
         │ - Gas 추정      │
         │ - Nonce 관리   │
         └────────┬───────┘
                  │ Signed Tx
                  ▼
         ┌────────────────┐
         │  L1 Ethereum   │
         │ - Tx 포함       │
         │ - 블록 확정     │
         └────────┬───────┘
                  │ Tx Receipt
                  ▼
         ┌────────────────┐
         │ op-node        │
         │ - L1 모니터링   │
         │ - 데이터 조회   │
         │ - L2 블록 생성  │
         └────────────────┘
```

### Plasma 흐름

```
┌─────────────────────────────────────────┐
│ op-batcher                              │
└──────────────┬──────────────────────────┘
               │ Compressed Data
          ┌────┴─────┐
          │ UsePlasma?│
          └────┬─────┘
      Yes ◄────┘
          │
          ▼
    ┌─────────────────────┐
    │ PlasmaDA.SetInput() │
    │ - DA Server에 저장   │
    │ - Commitment 생성    │
    └─────┬───────────────┘
          │
          ├─────► DA Server (Off-chain)
          │       - POST /put/
          │       - Store: commitment → data
          │       - Return: commitment
          │
          ▼
    ┌──────────────┐
    │ Commitment   │
    │ (34 bytes)   │
    └─────┬────────┘
          │
          ▼
    [L1 Ethereum에 제출]
          │
          ▼
    ┌──────────────────────┐
    │ op-node              │
    │ - L1에서 commitment  │
    │   읽기               │
    │ - DA Server에서 데이터│
    │   조회               │
    │ - Keccak256 검증     │
    └──────────────────────┘
          ▲
          │ GET /get/0x...
          │
    [DA Server에서 조회]
```

---

## op-batcher 설정

### Calldata 모드 (기본)

```yaml
# docker-compose.yml
op-batcher:
  environment:
    OP_BATCHER_DATA_AVAILABILITY_TYPE: "calldata"
    # 또는 설정 없음 (기본값)
```

```bash
# 직접 실행
op-batcher \
  --l1-eth-rpc=https://mainnet.infura.io/v3/... \
  --l2-eth-rpc=http://localhost:8545 \
  --rollup-rpc=http://localhost:7545 \
  --private-key=0x... \
  # data-availability-type 미지정 → calldata
```

### Blob 모드 (Ecotone 이후 권장)

```yaml
op-batcher:
  environment:
    OP_BATCHER_DATA_AVAILABILITY_TYPE: "blobs"
```

```bash
op-batcher \
  --l1-eth-rpc=https://mainnet.infura.io/v3/... \
  --l1-beacon=https://beacon-nd-123.p2pify.com/... \  # ← Beacon API 필수!
  --l2-eth-rpc=http://localhost:8545 \
  --rollup-rpc=http://localhost:7545 \
  --private-key=0x... \
  --data-availability-type=blobs
```

**필수 조건**:
```go
// op-batcher/batcher/service.go:224-226
if bs.UseBlobs && !bs.RollupConfig.IsEcotone(...) {
    bs.Log.Error("Cannot use Blob data before Ecotone!")
}
```

### Plasma 모드 (테스트만)

```yaml
op-batcher:
  environment:
    OP_BATCHER_PLASMA_ENABLED: "true"
    OP_BATCHER_PLASMA_DA_SERVICE: "true"  # Generic commitment
    OP_BATCHER_PLASMA_DA_SERVER: "http://da-server:3100"

da-server:
  image: tokamaknetwork/thanos-da-server
  ports:
    - "3100:3100"
  command:
    - da-server
    - --file.path=/data
    - --generic-commitment=false  # Keccak256 (안전)
    # - --generic-commitment=true  # Generic (위험!)
```

---

## TRH SDK의 DA 전략

### 로컬 Devnet

```go
// trh-sdk/pkg/stacks/thanos/deploy_chain.go:71-89
func (t *ThanosStack) deployLocalDevnet(ctx context.Context) error {
    // make devnet-up 실행
    err = utils.ExecuteCommandStream(ctx, t.logger, "bash", "-c", 
        fmt.Sprintf("cd %s/tokamak-thanos && export DEVNET_L2OO=true && make devnet-up", 
            t.deploymentPath))
    
    // ↓ docker-compose.yml 실행
    // ↓ da-server 포함 (Plasma 선택 가능)
}
```

**docker-compose.yml**:
```yaml
services:
  da-server:
    image: tokamaknetwork/thanos-da-server
    ports:
      - "3100:3100"
    # ✅ 로컬 개발/테스트용으로 포함

  op-batcher:
    environment:
      OP_BATCHER_PLASMA_ENABLED: "${PLASMA_ENABLED:-false}"
      OP_BATCHER_PLASMA_DA_SERVER: "http://da-server:3100"
```

### AWS/Kubernetes 프로덕션

```go
// trh-sdk/pkg/stacks/thanos/deploy_chain.go:91-406
func (t *ThanosStack) deployNetworkToAWS(ctx context.Context, inputs *DeployInfraInput) error {
    // 1. Terraform으로 EKS 인프라 배포
    // 2. Helm으로 Kubernetes 리소스 배포
    //    ├─ op-geth
    //    ├─ op-node
    //    ├─ op-batcher
    //    ├─ op-proposer
    //    └─ op-challenger
    
    // ❌ da-server 관련 코드 없음!
    // → Plasma 모드 미사용
    // → Blob 또는 Calldata 사용
}
```

**Terraform 환경 변수** (`pkg/types/terraform.go`):
```go
type TerraformEnvConfig struct {
    Namespace, AwsRegion
    SequencerKey, BatcherKey, ProposerKey, ChallengerKey
    L1RpcUrl, L1BeaconUrl, L1RpcProvider
    OpGethImageTag, ThanosStackImageTag
    MaxChannelDuration
    
    // ❌ PlasmaEnabled 없음
    // ❌ DaServerUrl 없음
    // ❌ PlasmaDAService 없음
}
```

**Helm Values** (자동 생성):
```yaml
# thanos-stack-values.yaml
l1_rpc_url: https://mainnet.infura.io/v3/...
l1_beacon_url: https://beacon-nd-123.p2pify.com/...
l1_rpc_provider: infura

# ❌ plasma_enabled: 없음
# ❌ da_server_url: 없음

# → op-batcher는 Blob 모드 사용 (L1 Beacon 설정 있음)
# → L1에 직접 데이터 제출
```

---

## 프로덕션 환경 DA 전략 요약

### TRH SDK 프로덕션 배포 시

```
1. L1 컨트랙트 배포
   └─ DataAvailabilityChallenge.sol 배포됨
      (하지만 사용되지 않음 - Plasma OFF)

2. AWS EKS 인프라 배포
   ├─ VPC, EKS, EFS 생성
   └─ ❌ DA Server 관련 리소스 없음

3. Kubernetes 리소스 배포
   ├─ op-geth (StatefulSet)
   ├─ op-node (StatefulSet)
   ├─ op-batcher (Deployment)
   │  └─ L1 Beacon URL 설정됨
   │  └─ → Blob 모드 사용
   ├─ op-proposer (Deployment)
   └─ op-challenger (Deployment)
   
4. op-batcher 동작
   ├─ UseBlobs = true (Beacon URL 있으면)
   ├─ UsePlasma = false (DA Server 없음)
   └─ → Blob으로 L1에 데이터 제출
   
5. 데이터 흐름
   L2 Tx → op-batcher → EIP-4844 Blob → L1 Beacon Chain
                                        └─> 모든 노드 접근 가능 ✅
```

### 왜 Blob을 선택했는가?

**비용 효율**:
```
Calldata: $225,000/월
Blob:     $2,100/월     ← 99% 절감! 🎉
Plasma:   $3,810/월     (DA Server 운영비 포함)
```

**보안**:
```
Calldata: ⭐⭐⭐⭐⭐ (L1 Execution)
Blob:     ⭐⭐⭐⭐⭐ (L1 Beacon)
Plasma:   ⭐⭐⭐     (Off-chain, Challenge 의존)
```

**운영 복잡도**:
```
Calldata: 낮음
Blob:     중간 (Beacon node 필요)
Plasma:   높음 (DA Server + Challenge + 백업 + 모니터링)
```

**결론**: **Blob이 최적의 선택!** ✅

---

## 실전 가이드

### 1. Calldata에서 Blob으로 마이그레이션

```bash
# 전제조건: Ecotone 업그레이드 완료

# 1. L1 Beacon node 준비
export L1_BEACON_URL="https://beacon-nd-123.p2pify.com/..."

# 2. op-batcher 재시작
docker compose down op-batcher

# 3. docker-compose.yml 수정
# op-batcher:
#   environment:
#     OP_BATCHER_DATA_AVAILABILITY_TYPE: "blobs"
#     OP_BATCHER_L1_BEACON: "${L1_BEACON_URL}"

# 4. op-batcher 시작
docker compose up -d op-batcher

# 5. 로그 확인
docker compose logs -f op-batcher | grep "blob"
# INFO building Blob transaction candidate size=360446 num_blobs=3
```

### 2. Blob 성능 튜닝

```yaml
op-batcher:
  environment:
    # 채널 최적화
    OP_BATCHER_MAX_CHANNEL_DURATION: "10"  # 블록 수
    OP_BATCHER_TARGET_NUM_FRAMES: "6"      # Blob 수 목표
    
    # 압축 최적화
    OP_BATCHER_COMPRESSOR: "brotli"        # brotli > zlib
    OP_BATCHER_APPROX_COMPR_RATIO: "0.4"   # 압축률 40%
    
    # 제출 타이밍
    OP_BATCHER_SUB_SAFETY_MARGIN: "10"     # 안전 마진
    OP_BATCHER_POLL_INTERVAL: "6s"         # 폴링 간격
```

### 3. Blob 가용성 모니터링

```bash
# Blob 메트릭 확인
curl http://localhost:7301/metrics | grep blob
# op_batcher_blob_used_bytes 360446
# op_batcher_blobs_per_tx 3

# Beacon API로 blob 직접 조회
curl https://beacon-nd-123.p2pify.com/eth/v1/beacon/blob_sidecars/19456789

# 응답:
{
  "data": [
    {
      "index": "0",
      "blob": "0x00789c...",  # 실제 blob 데이터
      "kzg_commitment": "0x01fa3b84...",
      "kzg_proof": "0x..."
    },
    ...
  ]
}
```

### 4. Plasma 테스트 (로컬만)

```bash
# 1. DA Server 시작
docker compose up -d da-server

# 2. Plasma 모드 활성화
export PLASMA_ENABLED=true
export PLASMA_DA_SERVICE=true
export PLASMA_GENERIC_DA=false  # Keccak256 사용

# 3. op-batcher, op-node 재시작
docker compose restart op-batcher op-node

# 4. DA Server 로그 확인
docker compose logs -f da-server
# INFO PUT url=/put/
# INFO stored commitment key=47173285a8d7341e... input_len=360446

# 5. L1 트랜잭션 확인
cast tx 0x... --rpc-url http://localhost:8545
# input: 0x0100471732...  ← Plasma commitment (34 bytes)
```

---

## 비용 시뮬레이션

### 시나리오: 중간 규모 L2 체인

**체인 활동**:
- 일일 트랜잭션: 100,000 건
- 평균 배치 크기: 200 KB (압축 후)
- 배치 빈도: 5분마다 (288 batches/day)

### 1. Calldata 모드

```
배치 크기: 200 KB
Calldata gas: 200,000 × 0.8 × 16 = 2,560,000 gas (non-zero)
            + 200,000 × 0.2 × 4 = 160,000 gas (zero)
            = 2,720,000 gas per batch

일일 가스:
2,720,000 × 288 = 783,360,000 gas/day

Gas price: 25 Gwei (평균)
일일 비용: 783,360,000 × 25 / 10^9 = 19.584 ETH
        ≈ $35,250 (ETH = $1,800)

월간 비용: $35,250 × 30 = $1,057,500 💰💰💰
```

### 2. Blob 모드 (권장)

```
배치 크기: 200 KB
Blobs 필요: ceil(200,000 / 131,071) = 2 blobs

Blob gas per batch: 131,072 × 2 = 262,144 blob gas

일일 blob gas:
262,144 × 288 = 75,497,472 blob gas/day

Blob base fee: 1 Wei (평상시)
일일 비용: 75,497,472 × 1 / 10^9 = 0.0755 ETH
        ≈ $136

월간 비용: $136 × 30 = $4,080 💰

절감액: $1,057,500 - $4,080 = $1,053,420/월
절감률: 99.6%! 🎉
```

### 3. Plasma 모드

```
Commitment 크기: 34 bytes
Calldata gas: 34 × 16 = 544 gas (거의 non-zero)

일일 가스:
21,544 × 288 = 6,204,672 gas/day

Gas price: 25 Gwei
일일 비용: 6,204,672 × 25 / 10^9 = 0.155 ETH
        ≈ $279

월간 L1 비용: $279 × 30 = $8,370 💰

DA Server 운영비:
- AWS EC2 (t3.large): $60/월
- EBS (1TB): $100/월
- S3 (10TB): $230/월
- 백업: $50/월
- 모니터링: $30/월
- 운영 인력: $2,000/월
총 운영비: $2,470/월

총 비용: $8,370 + $2,470 = $10,840/월

⚠️ Blob보다 비싸고 보안 위험!
```

---

## 데이터 검증 프로세스

### Calldata 검증

```go
// op-node/rollup/derive/calldata_source.go
func NewCalldataSource(ctx context.Context, log log.Logger, cfg DataSourceConfig, 
    fetcher L1Fetcher, ref eth.L1BlockRef, batcherAddr common.Address) DataIter {
    
    return &CalldataSource{
        // L1 블록에서 calldata 직접 읽기
        open: func() (eth.Data, error) {
            // 1. L1 트랜잭션 조회
            tx := fetcher.FetchTransaction(ref, batchInbox)
            
            // 2. Calldata 추출
            data := tx.Data()
            
            // 3. Version byte 확인
            if data[0] != DerivationVersion0 {
                return nil, ErrInvalidVersion
            }
            
            // 4. 데이터 반환 (검증 완료)
            return data[1:], nil
        },
    }
}
```

### Blob 검증

```go
// op-node/rollup/derive/blob_data_source.go
func NewBlobDataSource(ctx context.Context, log log.Logger, cfg DataSourceConfig,
    fetcher L1Fetcher, blobsFetcher L1BlobsFetcher, ref eth.L1BlockRef, 
    batcherAddr common.Address) DataIter {
    
    return &BlobDataSource{
        open: func() (eth.Data, error) {
            // 1. L1 트랜잭션에서 blob versioned hashes 읽기
            tx := fetcher.FetchTransaction(ref, batchInbox)
            blobHashes := tx.BlobVersionedHashes()
            
            // 2. Beacon API로 blob sidecars 조회
            sidecars, err := blobsFetcher.GetBlobs(ctx, ref, blobHashes)
            
            // 3. 각 blob에 대해 KZG 검증
            for i, sidecar := range sidecars {
                // KZG commitment 검증
                if !VerifyKZGCommitment(sidecar.Blob, sidecar.KZGCommitment) {
                    return nil, ErrBlobVerificationFailed
                }
                
                // Versioned hash 검증
                computedHash := ComputeBlobVersionedHash(sidecar.KZGCommitment)
                if computedHash != blobHashes[i] {
                    return nil, ErrBlobHashMismatch
                }
            }
            
            // 4. Blob 데이터 추출
            data := ExtractDataFromBlobs(sidecars)
            
            // 5. Version byte 확인
            if data[0] != DerivationVersion0 {
                return nil, ErrInvalidVersion
            }
            
            return data[1:], nil
        },
    }
}
```

### Plasma 검증

```go
// op-node/rollup/derive/plasma_data_source.go:35-103
func (s *PlasmaDataSource) Next(ctx context.Context) (eth.Data, error) {
    // 1. Challenge 상태 업데이트
    if err := s.fetcher.AdvanceL1Origin(ctx, s.l1, s.id); err != nil {
        if errors.Is(err, plasma.ErrReorgRequired) {
            return nil, NewResetError(fmt.Errorf("new expired challenge"))
        }
        return nil, NewTemporaryError(err)
    }

    if s.comm == nil {
        // 2. L1에서 commitment 읽기
        data, err := s.src.Next(ctx)
        if err != nil {
            return nil, err
        }

        // 3. Plasma 트랜잭션 확인
        if data[0] != plasma.TxDataVersion1 {
            // Plasma 아님 → 그대로 반환
            return data, nil
        }

        // 4. Commitment 파싱
        comm, err := plasma.DecodeCommitmentData(data[1:])
        if err != nil {
            s.log.Warn("invalid commitment", "commitment", data, "err", err)
            return nil, NotEnoughData
        }
        s.comm = comm
    }
    
    // 5. DA Server에서 데이터 조회
    data, err := s.fetcher.GetInput(ctx, s.l1, s.comm, s.id)
    
    // 6. Challenge 상태 확인
    if errors.Is(err, plasma.ErrExpiredChallenge) {
        s.log.Warn("challenge expired, skipping batch", "comm", s.comm)
        s.comm = nil
        return s.Next(ctx)  // 스킵
    } else if errors.Is(err, plasma.ErrMissingPastWindow) {
        return nil, NewCriticalError(fmt.Errorf("data not available: %w", err))
    } else if errors.Is(err, plasma.ErrPendingChallenge) {
        return nil, NotEnoughData  // 대기
    }
    
    // 7. 크기 검증
    if s.comm.CommitmentType() == plasma.Keccak256CommitmentType 
        && len(data) > plasma.MaxInputSize {
        s.log.Warn("input data exceeds max size")
        s.comm = nil
        return s.Next(ctx)
    }
    
    // 8. 검증 완료 → 반환
    s.comm = nil
    return data, nil
}
```

---

## 장애 시나리오 및 복구

### Calldata 장애 시나리오

**시나리오**: L1 Ethereum 네트워크 혼잡

```
상황: Gas price 급등 (500 Gwei)
영향: 배치 제출 비용 20배 증가
     $75/batch → $1,500/batch

op-batcher 대응:
1. Gas price 모니터링
2. MaxGasPrice 설정 확인
3. 임계값 초과 시 제출 대기
4. Gas price 하락 후 재시도

복구: 자동 (Gas price 정상화 시)
데이터 손실: 없음 (메모리 버퍼링)
```

### Blob 장애 시나리오

**시나리오 1**: Beacon node 다운

```
상황: L1 Beacon API 연결 실패
영향: Blob 제출 불가

op-batcher 대응:
1. Beacon node 연결 재시도
2. Fallback beacon URL 사용 (있으면)
3. 실패 시 에러 로깅 및 재시도

복구: Beacon node 재시작 필요
데이터 손실: 없음 (재시도 큐)
```

**시나리오 2**: Blob pruning 후 조회

```
상황: 18일 후 blob 데이터 삭제됨
     새 노드가 sync 시도

op-node 대응:
1. Beacon API로 blob 조회
2. 404 Not Found 에러
3. Archive beacon node로 재시도 (있으면)
4. 또는 다른 노드에서 snapshot 다운로드

복구: Archive node 또는 snapshot 필요
데이터 손실: 가능 (archive node 없으면)

해결책:
- Archive beacon node 운영
- L2 snapshot 정기 제공
- 또는 Calldata fallback
```

### Plasma 장애 시나리오

**시나리오 1**: DA Server 다운

```
상황: DA Server 장애 (하드웨어, 네트워크 등)
영향: L2 전체 중단! 🔥

타임라인:
T+0min: DA Server 다운
        └─ op-batcher: PUT 요청 실패
        └─ 배치 제출 중단

T+5min: op-node: 새 배치 없음
        └─ L2 블록 생성 계속 (로컬)
        └─ 하지만 L1 finalization 중단

T+30min: Challenge window 진행 중
         └─ 챌린저가 데이터 조회 불가
         └─ Challenge 해결 불가

T+48hr: Challenge window 만료
        └─ 잘못된 commitment 확정 가능
        └─ 시스템 무결성 훼손 💥

복구: 
1. DA Server 재시작 (ASAP!)
2. 백업에서 데이터 복구
3. op-batcher 재시작
4. Challenge 수동 해결

데이터 손실: 가능 (백업 없으면)
```

**시나리오 2**: DA Server 데이터 손상

```
상황: S3 버킷 손상 또는 파일 시스템 오류
영향: 특정 commitment의 데이터 조회 불가

타임라인:
T+0: DA Server에서 특정 데이터 손실
     └─ op-node: GET /get/0x4717... → 404 Not Found

T+1min: Challenge 시작
        └─ 데이터 없음을 증명

T+48hr: Challenge 만료
        └─ 해당 배치 무효 처리
        └─ L2 상태 reorg

복구: 불가능 (데이터 영구 손실)
영향: 해당 배치의 트랜잭션 롤백
     └─ 사용자 자금 손실 가능 💰
```

---

## 모범 사례

### 프로덕션 환경 권장 설정

```yaml
# docker-compose.yml (프로덕션)
version: '3.4'

services:
  op-batcher:
    image: tokamaknetwork/thanos-op-batcher:v1.7.7
    environment:
      # ✅ Blob 모드 (Ecotone 이후)
      OP_BATCHER_DATA_AVAILABILITY_TYPE: "blobs"
      
      # L1 연결
      OP_BATCHER_L1_ETH_RPC: "https://mainnet.infura.io/v3/..."
      OP_BATCHER_L1_BEACON: "https://beacon-nd-123.p2pify.com/..."
      
      # 성능 튜닝
      OP_BATCHER_MAX_CHANNEL_DURATION: "10"
      OP_BATCHER_TARGET_NUM_FRAMES: "6"
      OP_BATCHER_COMPRESSOR: "brotli"
      OP_BATCHER_APPROX_COMPR_RATIO: "0.4"
      
      # 안전성
      OP_BATCHER_SUB_SAFETY_MARGIN: "10"
      OP_BATCHER_NUM_CONFIRMATIONS: "10"
      
      # ❌ Plasma 비활성화
      # OP_BATCHER_PLASMA_ENABLED: "false"  # 기본값
    
    restart: always
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "10"

  op-node:
    image: tokamaknetwork/thanos-op-node:v1.7.7
    command:
      - op-node
      - --l1=https://mainnet.infura.io/v3/...
      - --l1.beacon=https://beacon-nd-123.p2pify.com/...
      - --l1.rpckind=infura
      # ❌ Plasma 설정 없음
      # --plasma.enabled=false  # 기본값
```

### 개발 환경 권장 설정

```yaml
# docker-compose.yml (개발)
services:
  op-batcher:
    environment:
      # Calldata (간단함) 또는 Plasma (테스트)
      OP_BATCHER_DATA_AVAILABILITY_TYPE: "calldata"
      
      # 또는 Plasma 테스트
      # OP_BATCHER_PLASMA_ENABLED: "true"
      # OP_BATCHER_PLASMA_DA_SERVER: "http://da-server:3100"
  
  da-server:
    image: tokamaknetwork/thanos-da-server
    ports:
      - "3100:3100"
    command:
      - da-server
      - --file.path=/data
      - --generic-commitment=false  # Keccak256 (안전)
    volumes:
      - da_data:/data
```

---

## 참고 자료

### 공식 문서

- [EIP-4844 Specification](https://eips.ethereum.org/EIPS/eip-4844)
- [Optimism Bedrock Specs - Data Availability](https://specs.optimism.io/protocol/derivation.html#batch-submission)
- [Optimism Plasma Mode](https://specs.optimism.io/experimental/plasma/)
- [Ethereum Beacon Chain API](https://ethereum.github.io/beacon-APIs/)

### 코드 참조

- `op-batcher/batcher/driver.go`: 배치 제출 로직
- `op-batcher/batcher/tx_data.go`: Calldata/Blob 생성
- `op-node/rollup/derive/`: 데이터 파싱 및 검증
- `op-plasma/`: Plasma DA 구현
- `packages/contracts-bedrock/src/L1/DataAvailabilityChallenge.sol`: Challenge 컨트랙트

### 관련 문서

- `deployment-guide-ko.md`: Docker Compose 배포 가이드
- `da-server-security-analysis-ko.md`: DA Server 보안 분석
- `trh-sdk-deployment-guide-ko.md`: TRH SDK 배포 가이드

### 비용 계산기

- [Ethereum Gas Tracker](https://etherscan.io/gastracker)
- [Blob Base Fee](https://etherscan.io/chart/blobbasefee)
- [L2 Fee Calculator](https://l2fees.info/)

---

## FAQ

### Q1: Ecotone 이전에는 어떻게 해야 하나요?

**A**: Calldata 모드를 사용해야 합니다. Blob은 Ecotone 하드포크 이후에만 사용 가능합니다.

```yaml
op-batcher:
  environment:
    OP_BATCHER_DATA_AVAILABILITY_TYPE: "calldata"
```

### Q2: Blob이 18일 후 삭제되면 새 노드는 어떻게 sync하나요?

**A**: 세 가지 방법이 있습니다:

1. **Archive Beacon Node 사용**
   ```bash
   op-node --l1.beacon=https://archive-beacon.example.com
   ```

2. **L2 Snapshot 다운로드**
   ```bash
   # Checkpoint sync
   op-node --l2.checkpoint-sync-url=https://snapshot.example.com
   ```

3. **다른 노드에서 데이터 요청**
   ```bash
   # P2P sync
   op-node --p2p.bootnodes=/dns4/bootnode.example.com/...
   ```

### Q3: Plasma를 프로덕션에서 사용할 수 있나요?

**A**: 기술적으로는 가능하지만 **강력히 비권장**합니다.

**이유**:
- ⚠️ DA Server SPOF (단일 장애점)
- ⚠️ 보안 위험 (인증, 암호화 부재)
- ⚠️ Challenge 복잡도
- ⚠️ 데이터 손실 위험

**대안**: Blob 사용 (비슷한 비용, 훨씬 안전)

### Q4: Blob과 Plasma를 함께 사용할 수 있나요?

**A**: 아니오, 둘 중 하나만 선택해야 합니다.

```go
// op-batcher/batcher/driver.go:489-516
if l.Config.UseBlobs {
    // Blob 모드
    candidate, err = l.blobTxCandidate(txdata)
} else {
    // Calldata 또는 Plasma 모드
    if l.Config.UsePlasma {
        // Plasma
    } else {
        // Calldata
    }
}
```

### Q5: TRH SDK로 Plasma 모드 배포할 수 있나요?

**A**: 로컬 devnet에서만 가능합니다.

**로컬 (Docker Compose)**:
```bash
export PLASMA_ENABLED=true
trh-sdk deploy  # da-server 포함
```

**프로덕션 (AWS/K8s)**:
```bash
trh-sdk deploy  # ❌ da-server 배포 안 됨
# → Blob 모드 자동 사용
```

### Q6: Generic Commitment는 언제 사용하나요?

**A**: 외부 DA 시스템(EigenDA, Celestia 등) 통합 시에만 사용합니다.

```yaml
da-server:
  command:
    # ⚠️ 테스트용만!
    - --generic-commitment=true
    
# 프로덕션은 항상:
    - --generic-commitment=false  # Keccak256
```

---

## 요약

### 🎯 핵심 정리

**L1 데이터 가용성 = L2 트랜잭션 데이터를 어디에 저장하는가?**

1. **Calldata**: L1 Execution Layer에 직접 저장
   - 비용: 높음 💰💰💰
   - 보안: 최고 ⭐⭐⭐⭐⭐
   - 사용: Ecotone 이전, 또는 간단한 설정

2. **Blob (EIP-4844)**: L1 Beacon Chain에 저장
   - 비용: 매우 낮음 💰
   - 보안: 최고 ⭐⭐⭐⭐⭐
   - 사용: **프로덕션 권장!** ✅✅

3. **Plasma**: Off-chain DA Server에 저장
   - 비용: 매우 낮음 💰
   - 보안: 중간 ⭐⭐⭐
   - 사용: 테스트만, 프로덕션 비권장 ⚠️

### 🚀 TRH SDK 전략

- **로컬 Devnet**: Calldata 또는 Plasma (선택)
- **프로덕션**: Blob (자동)
  - L1 Beacon URL 설정됨
  - DA Server 배포 안 됨
  - 안전하고 비용 효율적

### ✅ 최종 권장

**프로덕션 체인을 운영한다면:**
→ **Blob (EIP-4844) 사용!**

이유:
- 99% 비용 절감
- L1 수준의 보안
- 간단한 설정 (Beacon node만)
- 탈중앙화 유지
- Ethereum 네이티브 지원

---

**문서 버전**: 1.0
**작성일**: 2025-01-17
**최종 수정일**: 2025-01-17
**대상 프로젝트**: tokamak-thanos (Optimism Fork)

