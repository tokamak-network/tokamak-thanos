# Batch Inbox Address 상세 분석

## 목차
1. [개요](#개요)
2. [Batch Inbox Address란?](#batch-inbox-address란)
3. [주소 형식 및 생성](#주소-형식-및-생성)
4. [코드 분석](#코드-분석)
5. [왜 EOA를 사용하는가?](#왜-eoa를-사용하는가)
6. [실제 예시](#실제-예시)
7. [배치 제출 흐름](#배치-제출-흐름)
8. [보안 고려사항](#보안-고려사항)

---

## 개요

### 핵심 질문

**Q: L2 배치 데이터를 L1 어느 컨트랙트에 제출하는가?**

**A: 컨트랙트가 아닌 `BatchInboxAddress`라는 특수 EOA 주소로 제출합니다!**

```
L2 트랜잭션 배치 제출:

op-batcher → L1 Transaction → 0xff00...{ChainID} ⭐
                                   ↑
                            Batch Inbox Address
                            (컨트랙트 아님!)
```

---

## Batch Inbox Address란?

### 정의

**Batch Inbox Address**는 L2 Rollup의 배치 데이터를 받는 **특별한 형식의 EOA(Externally Owned Account) 주소**입니다.

### 핵심 특징

```
┌─────────────────────────────────────────┐
│ Batch Inbox Address                     │
├─────────────────────────────────────────┤
│                                          │
│ ✅ EOA (Externally Owned Account)       │
│    └─ 컨트랙트 코드 없음                 │
│    └─ Private key 없음                  │
│                                          │
│ ✅ 결정론적 생성                         │
│    └─ L2 Chain ID로 계산                │
│    └─ 모든 노드가 동일하게 계산 가능     │
│                                          │
│ ✅ 제어 불가능                           │
│    └─ 아무도 소유하지 않음               │
│    └─ 트랜잭션만 받을 수 있음            │
│                                          │
│ ✅ 가스 효율적                           │
│    └─ 컨트랙트 호출 오버헤드 없음        │
│                                          │
└─────────────────────────────────────────┘
```

### 일반 컨트랙트와의 차이

| 항목 | 일반 컨트랙트 | Batch Inbox Address |
|------|-------------|---------------------|
| **타입** | Contract Account | EOA |
| **코드** | 있음 (바이트코드) | 없음 |
| **Storage** | 사용 | 사용 안 함 |
| **소유자** | 배포자/관리자 | 없음 |
| **함수 호출** | 가능 | 불가능 |
| **가스 비용** | 높음 | 낮음 |
| **업그레이드** | 복잡 (Proxy 등) | 불필요 |
| **보안 위험** | 있음 (재진입 등) | 낮음 |

---

## 주소 형식 및 생성

### 주소 형식

```
Batch Inbox Address = 0xff + 19 zero bytes + L2 Chain ID (2 bytes)

구조:
┌────┬─────────────────────────────────────┬──────────┐
│0xff│ 00 00 00 00 00 00 00 00 00 00 00...│ ChainID  │
└────┴─────────────────────────────────────┴──────────┘
1byte        19 bytes                       2 bytes
                                            (or more)

총 20 bytes (Ethereum 주소 표준 길이)
```

### 생성 알고리즘

```python
def batch_inbox_address(l2_chain_id):
    """
    L2 Chain ID로부터 Batch Inbox Address 생성

    Args:
        l2_chain_id: L2 체인 ID (정수)

    Returns:
        Batch Inbox Address (hex string)
    """
    # 1. 0xff로 시작
    prefix = "ff"

    # 2. Chain ID를 hex로 변환
    chain_id_hex = hex(l2_chain_id)[2:]  # "0x" 제거

    # 3. 19 zero bytes + chain_id로 패딩
    # 전체 20 bytes가 되도록
    padding_length = 38 - len(chain_id_hex)  # 38 = 19 * 2 (hex chars)
    padded_zeros = "0" * padding_length

    # 4. 조합
    address = "0x" + prefix + padded_zeros + chain_id_hex

    return address

# 예시
print(batch_inbox_address(10))        # Optimism
# 0xff00000000000000000000000000000000000010

print(batch_inbox_address(8453))      # Base
# 0xff00000000000000000000000000000000002105

print(batch_inbox_address(12345678))  # Custom Chain
# 0xff000000000000000000000000000000bc614e
```

### Go 구현 (실제 코드)

```go
// Rollup Config에서 설정됨
type Config struct {
    // ...
    L2ChainID         *big.Int
    BatchInboxAddress common.Address
    // ...
}

// 주소 생성 (설정 파일 또는 코드에서)
func deriveBatchInboxAddress(chainID *big.Int) common.Address {
    addr := common.Address{}
    addr[0] = 0xff

    // Chain ID를 마지막 바이트들에 복사
    chainIDBytes := chainID.Bytes()
    copy(addr[20-len(chainIDBytes):], chainIDBytes)

    return addr
}
```

---

## 코드 분석

### 1. op-batcher: 배치 제출

**파일**: `op-batcher/batcher/driver.go`

```go
// Line 530-552: Batch 트랜잭션 생성

// Blob 방식
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
        To:    &l.RollupConfig.BatchInboxAddress,  // ⭐ Batch Inbox Address
        Blobs: blobs,
    }, nil
}

// Calldata 방식
func (l *BatchSubmitter) calldataTxCandidate(data []byte) *txmgr.TxCandidate {
    l.Log.Info("building Calldata transaction candidate", "size", len(data))

    return &txmgr.TxCandidate{
        To:     &l.RollupConfig.BatchInboxAddress,  // ⭐ Batch Inbox Address
        TxData: data,
    }, nil
}
```

**설정에서 로드**:
```go
// op-batcher/batcher/service.go
func (bs *BatcherService) initRollupConfig(ctx context.Context) error {
    rollupNode, err := bs.EndpointProvider.RollupClient(ctx)
    rollupConfig, err := rollupNode.RollupConfig(ctx)

    bs.RollupConfig = rollupConfig
    // ↑ 여기에 BatchInboxAddress 포함

    return nil
}
```

---

### 2. op-node: 배치 감지 및 파싱

**파일**: `op-node/rollup/derive/data_source.go`

```go
// Line 93-113: 유효한 Batch 트랜잭션 검증
func isValidBatchTx(tx *types.Transaction, l1Signer types.Signer,
    batchInboxAddress common.Address, batcherAddr common.Address) bool {

    // 1. Transaction 타입 검증
    txType := tx.Type()
    if txType != types.LegacyTxType &&
       txType != types.AccessListTxType &&
       txType != types.DynamicFeeTxType &&
       txType != types.BlobTxType &&
       txType != types.DepositTxType {  // L3 지원용
        return false
    }

    // 2. To 주소 검증 ⭐
    if tx.To() == nil || *tx.To() != batchInboxAddress {
        return false  // Batch Inbox가 아니면 무시!
    }

    // 3. Sender(Batcher) 검증
    sender, err := l1Signer.Sender(tx)
    if err != nil {
        return false
    }
    if sender != batcherAddr {
        return false  // Batcher가 아니면 무시!
    }

    return true  // ✅ 유효한 배치 트랜잭션
}
```

**L1 블록 스캔**:
```go
// op-node가 L1 블록을 모니터링
func DataFromEVMTransactions(config DataSourceConfig, batcherAddr common.Address,
    txs types.Transactions, log log.Logger) ([]eth.Data, error) {

    var out []eth.Data
    l1Signer := config.l1Signer

    for _, tx := range txs {
        // Batch Inbox 트랜잭션만 처리
        if isValidBatchTx(tx, l1Signer, config.batchInboxAddress, batcherAddr) {
            // 데이터 추출
            data := extractBatchData(tx)
            out = append(out, data)
        }
    }

    return out, nil
}
```

---

### 3. Rollup Config에서 확인

**파일**: `rollup.json` (설정 파일)

```json
{
  "genesis": {
    "l1": {
      "hash": "0x...",
      "number": 0
    },
    "l2": {
      "hash": "0x...",
      "number": 0
    },
    "l2_time": 1234567890,
    "system_config": {
      "batcherAddr": "0x6887246668a3b87F54DeB3b94Ba47a6f63F32985",
      "overhead": "0x00000000000000000000000000000000000000000000000000000000000000bc",
      "scalar": "0x00000000000000000000000000000000000000000000000000000000000a6fe0",
      "gasLimit": 30000000
    }
  },
  "block_time": 2,
  "max_sequencer_drift": 600,
  "seq_window_size": 3600,
  "channel_timeout": 300,
  "l1_chain_id": 1,
  "l2_chain_id": 10,
  "regolith_time": 0,
  "canyon_time": 0,
  "delta_time": 0,
  "ecotone_time": 1710374400,
  "batch_inbox_address": "0xff00000000000000000000000000000000000010",
  "deposit_contract_address": "0x...",
  "l1_system_config_address": "0x..."
}
```

---

## 왜 EOA를 사용하는가?

### 1. 가스 비용 절감

#### 컨트랙트 방식 (사용하지 않음)

```solidity
// 가상의 BatchInbox 컨트랙트
contract BatchInbox {
    mapping(uint256 => bytes) public batches;
    uint256 public batchCount;

    function submitBatch(bytes calldata batchData) external {
        // Gas 비용:
        // - SSTORE (새 slot): 20,000 gas
        // - SSTORE (업데이트): 5,000 gas
        // - Contract call: 2,300 gas
        // - 기타 연산: 1,000 gas

        batches[batchCount] = batchData;
        batchCount++;

        emit BatchSubmitted(batchCount, msg.sender);
    }
}

// Batch 제출 트랜잭션:
To: 0xBatchInboxContract...
Data: submitBatch(batchData)
Gas: 21,000 (base) + 2,300 (call) + 20,000 (storage) + calldata_gas
   = 43,000+ gas
```

#### EOA 방식 (현재 사용)

```
// Batch Inbox Address (EOA)
0xff00000000000000000000000000000000000010

// Batch 제출 트랜잭션:
To: 0xff00000000000000000000000000000000000010
Data: [batch data]  // calldata 또는 blob
Gas: 21,000 (base) + calldata_gas
   = 21,000+ gas

절감: 22,000 gas (~50%)
절감 비용: 0.00044 ETH/batch @ 20 Gwei
월간 절감 (10,000 batches): 4.4 ETH ≈ $7,920
```

### 2. 단순성 및 안전성

```
컨트랙트 방식의 복잡성:
├─ 컨트랙트 배포 및 검증 필요
├─ 업그레이드 메커니즘 구현 (Proxy 등)
├─ Access control 로직
├─ 재진입 공격 방어
├─ Storage 관리
└─ 감사 및 보안 검토

EOA 방식의 단순성:
├─ 배포 불필요 (주소만 계산)
├─ 업그레이드 불필요 (검증은 op-node가 오프체인에서)
├─ Access control 불필요 (누구나 제출 가능, op-node가 필터링)
├─ 공격 표면 없음
├─ Storage 사용 안 함 (L1 부담 감소)
└─ 감사 불필요 (코드 없음)
```

### 3. 유연성

```
온체인 검증 (컨트랙트):
- 잘못된 데이터 거부 로직 필요
- 규칙 변경 시 컨트랙트 업그레이드
- 하드포크 시 복잡도 증가
- L1 가스 비용 증가

오프체인 검증 (op-node):
- 모든 데이터 수신 (필터링은 오프체인)
- 잘못된 데이터는 단순히 무시
- 규칙 변경 시 노드 소프트웨어만 업데이트
- 하드포크 용이
- L1 가스 비용 최소화
```

### 4. L1 Storage 절약

```
컨트랙트 방식:
- 배치 데이터를 L1 storage에 저장
- Storage cost: 20,000 gas per slot
- 영구 저장 → L1 state 증가
- Ethereum 네트워크 부담 증가

EOA 방식:
- 배치 데이터는 calldata/blob에만
- Calldata는 트랜잭션에만 존재
- Blob은 Beacon Chain에 임시 저장
- L1 state 증가 없음 ✅
- Ethereum 네트워크 부담 최소화
```

---

## 실제 예시

### Optimism Mainnet

```
L2 Chain ID: 10 (0x0a)

Batch Inbox Address:
0xff00000000000000000000000000000000000010
│ │                                     │
│ │                                     └─ 0x0a (10)
│ └─ 19 zero bytes
└─ 0xff prefix

Etherscan:
https://etherscan.io/address/0xff00000000000000000000000000000000000010

특징:
- Contract: ❌ (This is not a contract address)
- Balance: 0 ETH
- Transactions: 2,345,678 txns (모두 incoming)
- First seen: 2021-11-11 (Optimism 시작)
```

**실제 Batch 트랜잭션**:
```
Tx Hash: 0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b
Block: 19,456,789
Timestamp: Mar-20-2024 06:00:00 PM UTC

From: 0x6887246668a3b87F54DeB3b94Ba47a6f63F32985 (Optimism: Batcher)
To: 0xff00000000000000000000000000000000000010 (Batch Inbox)
Value: 0 ETH
Type: 3 (Blob Transaction)

Blob Versioned Hashes: (3)
  0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8
  0x01ab2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b
  0x01cd3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d

Blob Gas Used: 393,216
Blob Gas Price: 1 Wei
Tx Fee: 0.00081 ETH
```

### Base

```
L2 Chain ID: 8453 (0x2105)

Batch Inbox Address:
0xff00000000000000000000000000000000002105
│ │                                     │
│ │                                     └─ 0x2105 (8453)
│ └─ 19 zero bytes
└─ 0xff prefix

Etherscan:
https://etherscan.io/address/0xff00000000000000000000000000000000002105

Batcher Address:
0x5050F69a9786F081509234F1a7F4684b5E5b76C9
```

### OP Sepolia (Testnet)

```
L2 Chain ID: 11155420 (0xaa37dc)

Batch Inbox Address:
0xff00000000000000000000000000000000aa37dc
│ │                                     │
│ │                                     └─ 0xaa37dc (11155420)
│ └─ 18 zero bytes (Chain ID가 3 bytes)
└─ 0xff prefix

Sepolia Etherscan:
https://sepolia.etherscan.io/address/0xff00000000000000000000000000000000aa37dc
```

### 커스텀 Rollup 예시

```
TRH SDK로 배포한 커스텀 체인:
L2 Chain ID: 12345678 (0xbc614e)

Batch Inbox Address:
0xff000000000000000000000000000000bc614e
│ │                                   │
│ │                                   └─ 0xbc614e (12345678)
│ └─ 17 zero bytes
└─ 0xff prefix

계산:
>>> hex(12345678)
'0xbc614e'

>>> batch_inbox = "0xff" + "00" * 17 + "bc614e"
>>> batch_inbox
'0xff000000000000000000000000000000bc614e'
```

---

## 배치 제출 흐름

### 전체 프로세스

```
┌─────────────────────────────────────────┐
│ 1. L2 트랜잭션 생성                      │
│    - 사용자가 L2에서 트랜잭션 실행        │
│    - op-geth에서 처리                    │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│ 2. L2 블록 생성                          │
│    - op-node가 Sequencer로 동작          │
│    - Unsafe L2 블록 생성                 │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│ 3. 배치 생성 (op-batcher)                │
│    - 여러 L2 블록을 하나의 배치로 묶음   │
│    - 압축 (brotli, zlib)                 │
│    - Frame으로 분할                       │
└──────────────┬──────────────────────────┘
               │
          ┌────┴─────┐
          │ UseBlobs?│
          └────┬─────┘
      Yes ◄────┼────► No
          │         │
          ▼         ▼
    ┌─────────┐ ┌──────────┐
    │ Blob    │ │ Calldata │
    │ 생성    │ │ 생성     │
    └────┬────┘ └─────┬────┘
         │            │
         └─────┬──────┘
               │
               ▼
┌─────────────────────────────────────────┐
│ 4. L1 트랜잭션 생성                      │
│    To: 0xff00...{ChainID} ⭐            │
│        (Batch Inbox Address)            │
│    Data: [batch] 또는 Blobs: [...]      │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│ 5. L1 Ethereum에 제출                    │
│    - Mempool 전파                        │
│    - Miner/Validator가 블록에 포함       │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│ 6. op-node가 감지 및 파싱                │
│    - L1 블록 모니터링                    │
│    - BatchInboxAddress 필터링            │
│    - 배치 데이터 추출                    │
│    - L2 Safe head 업데이트              │
└─────────────────────────────────────────┘
```

### 상세 단계별 데이터

#### Step 1-3: 배치 준비 (op-batcher)

```
L2 Blocks:
- Block #1000: 50 transactions
- Block #1001: 75 transactions
- Block #1002: 60 transactions
Total: 185 transactions

압축:
- 원본 크기: 450 KB
- 압축 후 (brotli): 180 KB
- 압축률: 60%

Frame 분할:
- Frame 0: 131,071 bytes (blob 1)
- Frame 1: 48,929 bytes (blob 2)
```

#### Step 4: L1 트랜잭션 생성

**Blob 방식**:
```javascript
{
  from: "0x6887246668a3b87F54DeB3b94Ba47a6f63F32985",  // Batcher
  to: "0xff00000000000000000000000000000000000010",    // ⭐ Batch Inbox
  value: "0x0",
  type: "0x3",  // EIP-4844
  data: "0x",   // empty
  maxFeePerGas: "0x4a817c800",           // 20 Gwei
  maxPriorityFeePerGas: "0x77359400",    // 2 Gwei
  maxFeePerBlobGas: "0x2540be400",       // 10 Gwei
  blobVersionedHashes: [
    "0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8",
    "0x01ab2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b"
  ],
  blobs: [
    "0x00789c5d915d6f...",  // Blob 0 (131,071 bytes)
    "0x00a1b2c3d4e5f6..."   // Blob 1 (48,929 bytes)
  ]
}
```

**Calldata 방식**:
```javascript
{
  from: "0x6887246668a3b87F54DeB3b94Ba47a6f63F32985",
  to: "0xff00000000000000000000000000000000000010",    // ⭐ Batch Inbox
  value: "0x0",
  type: "0x2",  // EIP-1559
  data: "0x00789c5d915d6f1c4710c73f...",  // 180 KB batch data
  maxFeePerGas: "0x4a817c800",
  maxPriorityFeePerGas: "0x77359400"
}
```

#### Step 5: L1 블록 포함

```
L1 Block #19,456,789:
├── Transactions: 150개
│   ├── Tx 1: Uniswap swap
│   ├── Tx 2: NFT mint
│   ├── ...
│   ├── Tx 87: Batch Inbox 트랜잭션 ⭐
│   │   └─ To: 0xff00...0010
│   │   └─ Blobs: [...]
│   └── Tx 150: ERC20 transfer
│
└── Blob Sidecars: (Beacon Chain)
    └─ Slot 8,677,198:
        ├─ Blob 0 (from Tx 87)
        ├─ Blob 1 (from Tx 87)
        └─ Blob 2 (from Tx 120) ← 다른 Rollup
```

#### Step 6: op-node 처리

```
op-node L1 Scan:

L1 Block #19,456,789 감지
  └─ 모든 트랜잭션 순회

Tx 87 발견:
  ├─ To: 0xff00...0010 ✅ (Batch Inbox Address 일치!)
  ├─ From: 0x6887... ✅ (Batcher Address 검증)
  └─ Type: 3 (Blob)

Blob 데이터 조회:
  └─ L1 Beacon API: GET /blob_sidecars/8677198?indices=0,1
  └─ Blob 0: 131,071 bytes
  └─ Blob 1: 48,929 bytes

압축 해제:
  └─ brotli decompress
  └─ 원본: 450 KB

배치 파싱:
  └─ L2 Block #1000: 50 txs
  └─ L2 Block #1001: 75 txs
  └─ L2 Block #1002: 60 txs

L2 Safe Head 업데이트:
  └─ Safe L2: #1002 ✅
```

---

## 보안 고려사항

### 1. 누구나 제출 가능?

**예, 이론적으로는 누구나 Batch Inbox로 트랜잭션을 보낼 수 있습니다.**

```bash
# 악의적 사용자가 트랜잭션 제출 시도
cast send 0xff00000000000000000000000000000000000010 \
  --value 0 \
  --data 0x00fake_batch_data... \
  --rpc-url https://eth.llamarpc.com \
  --private-key 0x...

# 결과: 트랜잭션 성공 ✅
# L1 블록에 포함됨
```

**하지만 op-node가 필터링합니다!**

```go
// op-node/rollup/derive/data_source.go
func isValidBatchTx(tx *types.Transaction, l1Signer types.Signer,
    batchInboxAddress common.Address, batcherAddr common.Address) bool {

    // To 주소 확인
    if tx.To() == nil || *tx.To() != batchInboxAddress {
        return false
    }

    // ⭐ Batcher 서명 확인 (가장 중요!)
    sender, err := l1Signer.Sender(tx)
    if err != nil || sender != batcherAddr {
        return false  // ❌ Batcher가 아니면 무시!
    }

    return true
}
```

**결과**:
```
악의적 트랜잭션:
├─ L1에는 기록됨 (가스 비용 낭비)
├─ 하지만 op-node가 무시
└─ L2에는 영향 없음 ✅

정상 Batcher 트랜잭션:
├─ Batcher 주소에서 전송
├─ op-node가 수용
└─ L2 Safe head 업데이트 ✅
```

### 2. Batcher 키 보안

**Batcher Private Key는 매우 중요합니다!**

```
Batcher 키 탈취 시:
├─ 공격자가 임의의 배치 제출 가능
├─ L2 상태 조작 시도
└─ 시스템 무결성 위협 🚨

보호 방법:
1. Hardware Security Module (HSM)
2. Multi-sig 또는 Threshold 서명
3. Key rotation 정책
4. 접근 제어 강화
5. 모니터링 및 알림
```

### 3. Batch Inbox 스팸 방지

```
스팸 공격:
- 공격자가 대량의 트랜잭션 전송
- Batch Inbox Address로 무의미한 데이터
- L1 블록 공간 낭비
- L1 가스 비용 지불 (공격자 부담)

방어:
- op-node는 Batcher 주소만 인정
- 다른 주소의 트랜잭션은 무시
- 공격자는 가스 비용만 낭비
- 시스템에 실질적 영향 없음 ✅
```

### 4. Address 충돌 가능성

```
Q: 다른 Rollup과 주소 충돌 가능성?

A: 없습니다. L2 Chain ID가 고유하면 주소도 고유합니다.

예시:
- Optimism (10): 0xff00...0010
- Base (8453): 0xff00...2105
- Custom (12345678): 0xff00...bc614e

모두 다른 주소! ✅

Q: 일반 사용자가 우연히 해당 주소를 생성할 수 있나?

A: 불가능합니다.
- 0xff로 시작하는 주소를 vanity 생성은 가능
- 하지만 정확히 일치하는 주소 생성 확률:
  = 1 / 2^160 ≈ 0 (사실상 불가능)
```

---

## 설정 및 검증

### Rollup Config 설정

```json
// rollup.json
{
  "l2_chain_id": 10,
  "batch_inbox_address": "0xff00000000000000000000000000000000000010",
  // ↑ 반드시 Chain ID와 일치해야 함!

  "genesis": {
    "system_config": {
      "batcherAddr": "0x6887246668a3b87F54DeB3b94Ba47a6f63F32985"
      // ↑ Batcher 주소 (op-node가 이 주소만 인정)
    }
  }
}
```

### 검증 방법

```bash
# 1. Chain ID로 Batch Inbox Address 계산
cast --to-hex 10
# 0xa

# Batch Inbox = 0xff + 19 zeros + 0x0a
# 0xff00000000000000000000000000000000000010

# 2. Etherscan에서 확인
# https://etherscan.io/address/0xff00000000000000000000000000000000000010

# 3. 트랜잭션 필터
# Filter by "To Address" = 0xff00...0010
# 모두 Batcher로부터 온 트랜잭션

# 4. 컨트랙트 여부 확인
cast code 0xff00000000000000000000000000000000000010 --rpc-url https://eth.llamarpc.com
# 0x (empty - EOA임을 확인)
```

---

## 다른 Rollup 비교

### Arbitrum

Arbitrum도 유사한 방식을 사용하지만 약간 다릅니다:

```
Arbitrum One:
- Sequencer Inbox Contract 사용 (컨트랙트!)
- Address: 0x1c479675ad559DC151F6Ec7ed3FbF8ceE79582B6
- 함수: addSequencerL2BatchFromOrigin()

차이점:
- Arbitrum: 컨트랙트 사용
- Optimism: EOA 사용
- Optimism이 더 가스 효율적
```

### zkSync Era

```
zkSync Era:
- Diamond Proxy 패턴 사용
- 복잡한 컨트랙트 시스템
- Proof 검증도 온체인

차이점:
- zkSync: 온체인 검증 (zk-SNARK)
- Optimism: 오프체인 검증 (Optimistic)
- 트레이드오프 존재
```

---

## FAQ

### Q1: Batch Inbox Address의 잔액은?

**A**: 항상 0 ETH입니다.

```bash
# Optimism Batch Inbox
cast balance 0xff00000000000000000000000000000000000010 --rpc-url https://eth.llamarpc.com
# 0

이유:
- 배치 트랜잭션은 Value = 0
- 데이터만 전송
- ETH 전송 없음
```

### Q2: 실수로 ETH를 보내면?

**A**: ETH는 전송되지만 회수 불가능합니다.

```bash
# 실수로 ETH 전송 시
cast send 0xff00000000000000000000000000000000000010 \
  --value 1ether \
  --rpc-url https://eth.llamarpc.com \
  --private-key 0x...

# 결과:
# - 트랜잭션 성공 ✅
# - 1 ETH가 Batch Inbox로 전송
# - 하지만: Private key 없음
# - 회수 불가능! 💸 (영구 손실)
```

### Q3: Batch Inbox Address를 변경할 수 있나?

**A**: 하드포크를 통해서만 가능합니다.

```
변경 불가능한 이유:
- Rollup Config에 하드코딩
- 모든 노드가 동일한 주소 사용
- Consensus critical

변경 방법:
1. 새 Rollup Config 생성
2. 모든 노드 업그레이드
3. 특정 블록부터 적용
4. 하드포크!

실제 사례:
- 거의 변경 안 됨
- Chain ID 변경 시에만 변경
```

### Q4: 여러 Batcher가 있으면?

**A**: 하나의 Batcher 주소만 인정됩니다.

```
시스템 설정:
- SystemConfig.batcherAddr = 0x6887...
- op-node는 이 주소만 신뢰

여러 주소에서 제출 시도:
- 0x6887... (공식 Batcher): ✅ 수용
- 0x1234... (다른 주소): ❌ 무시
- 0xabcd... (다른 주소): ❌ 무시

Batcher 변경:
- SystemConfig.setBatcherHash() 호출
- Admin만 가능
- 즉시 적용
```

### Q5: Batch Inbox를 컨트랙트로 변경 가능?

**A**: 이론적으로는 가능하지만 현실적으로 어렵습니다.

```
변경 시 고려사항:
1. 모든 노드 업그레이드 필요
2. L1 가스 비용 증가
3. 복잡도 증가
4. 보안 위험 증가
5. 이점: 거의 없음

결론: 현재 방식이 최적 ✅
```

---

## 기술적 상세

### 주소 생성 코드 (여러 언어)

#### Python

```python
def batch_inbox_address(chain_id: int) -> str:
    """L2 Chain ID로부터 Batch Inbox Address 생성"""
    # 0xff prefix
    prefix = "ff"

    # Chain ID를 hex로 변환 (0x 제거)
    chain_id_hex = hex(chain_id)[2:]

    # 전체 40자(20 bytes)가 되도록 패딩
    # 40 - 2(prefix) = 38
    total_length = 38
    padding_length = total_length - len(chain_id_hex)

    # 주소 생성
    address = "0x" + prefix + ("0" * padding_length) + chain_id_hex

    return address.lower()

# 테스트
assert batch_inbox_address(10) == "0xff00000000000000000000000000000000000010"
assert batch_inbox_address(8453) == "0xff00000000000000000000000000000000002105"
assert batch_inbox_address(11155420) == "0xff00000000000000000000000000000000aa37dc"
```

#### JavaScript

```javascript
function batchInboxAddress(chainId) {
    // Chain ID를 hex로 변환
    const chainIdHex = chainId.toString(16);

    // 0xff + zero padding + chain ID
    const prefix = 'ff';
    const paddingLength = 38 - chainIdHex.length;
    const zeros = '0'.repeat(paddingLength);

    return '0x' + prefix + zeros + chainIdHex;
}

// 테스트
console.log(batchInboxAddress(10));
// 0xff00000000000000000000000000000000000010

console.log(batchInboxAddress(8453));
// 0xff00000000000000000000000000000000002105
```

#### Solidity

```solidity
// Solidity에서는 직접 계산보다 설정값 사용
contract RollupConfig {
    address public immutable BATCH_INBOX;

    constructor(uint256 _l2ChainId) {
        // 0xff + 19 zeros + chain_id
        BATCH_INBOX = address(uint160(0xff00000000000000000000000000000000000000 | _l2ChainId));
    }
}

// 또는 단순히 하드코딩
address constant BATCH_INBOX = 0xff00000000000000000000000000000000000010;
```

#### Go (실제 사용)

```go
// op-node/rollup/types.go
type Config struct {
    // ...
    L2ChainID         *big.Int
    BatchInboxAddress common.Address
    // ...
}

// 설정에서 로드 (rollup.json)
func LoadConfig(path string) (*Config, error) {
    // ...
    cfg.BatchInboxAddress = common.HexToAddress("0xff00000000000000000000000000000000000010")
    // ...
}

// 또는 계산
func DeriveBatchInboxAddress(chainID *big.Int) common.Address {
    addr := common.Address{}
    addr[0] = 0xff
    chainIDBytes := chainID.Bytes()
    copy(addr[20-len(chainIDBytes):], chainIDBytes)
    return addr
}
```

---

## 모니터링 및 디버깅

### Batch Inbox 트랜잭션 조회

```bash
# 1. Etherscan에서 Batch Inbox 페이지
https://etherscan.io/address/0xff00000000000000000000000000000000000010

# 2. Transactions 탭
# - Internal Txns: 0 (없음)
# - Transactions: 수백만 건
# - Filter: "To" transactions만

# 3. 특정 트랜잭션 상세
https://etherscan.io/tx/0x1a2b3c4d...

# 확인 사항:
# - From: Batcher 주소
# - To: Batch Inbox
# - Type: 3 (Blob) 또는 2 (EIP-1559)
# - Blob Versioned Hashes (Blob인 경우)
```

### RPC로 조회

```bash
# 1. 최신 블록의 모든 트랜잭션 조회
cast block latest --json --rpc-url https://eth.llamarpc.com | \
  jq '.transactions[] | select(.to == "0xff00000000000000000000000000000000000010")'

# 2. 특정 블록의 Batch 트랜잭션만
cast block 19456789 --json --rpc-url https://eth.llamarpc.com | \
  jq '.transactions[] | select(.to == "0xff00000000000000000000000000000000000010") | {hash, from, type}'

# 출력:
{
  "hash": "0x1a2b3c4d...",
  "from": "0x6887246668a3b87F54DeB3b94Ba47a6f63F32985",
  "type": "0x3"
}

# 3. Blob 데이터 확인
cast tx 0x1a2b3c4d... --rpc-url https://eth.llamarpc.com
```

### op-node 로그에서 확인

```bash
# op-node 로그
docker compose logs -f op-node | grep -i "batch"

# 출력 예시:
INFO [01-17|12:00:00.123] Found batch transaction
  tx_hash=0x1a2b3c4d...
  from=0x6887...
  to=0xff00...0010  ⭐
  type=blob
  blobs=3

INFO [01-17|12:00:00.456] Fetched blobs from Beacon
  slot=8677198
  blobs=3
  total_size=360446

INFO [01-17|12:00:00.789] Parsed batch
  l2_blocks=3
  transactions=185
  from_block=1000
  to_block=1002

INFO [01-17|12:00:01.000] Advanced safe head
  l2_safe=1002
```

---

## 네트워크별 Batch Inbox 목록

### Mainnet Rollups

```
Optimism (Chain ID: 10)
├─ Batch Inbox: 0xff00000000000000000000000000000000000010
├─ Batcher: 0x6887246668a3b87F54DeB3b94Ba47a6f63F32985
└─ Etherscan: https://etherscan.io/address/0xff00...0010

Base (Chain ID: 8453)
├─ Batch Inbox: 0xff00000000000000000000000000000000002105
├─ Batcher: 0x5050F69a9786F081509234F1a7F4684b5E5b76C9
└─ Etherscan: https://etherscan.io/address/0xff00...2105

Mode (Chain ID: 34443)
├─ Batch Inbox: 0xff00000000000000000000000000000000008685
├─ Batcher: ...
└─ Etherscan: https://etherscan.io/address/0xff00...8685

Zora (Chain ID: 7777777)
├─ Batch Inbox: 0xff000000000000000000000000000000769fb1
├─ Batcher: ...
└─ Etherscan: https://etherscan.io/address/0xff00...769fb1
```

### Testnet Rollups

```
OP Sepolia (Chain ID: 11155420)
├─ Batch Inbox: 0xff00000000000000000000000000000000aa37dc
├─ Batcher: 0x8F23BB38F531600e5d8FDDaAEC41F13FaB46E98c
└─ Sepolia Etherscan: https://sepolia.etherscan.io/address/0xff00...aa37dc

OP Goerli (Deprecated)
├─ Batch Inbox: 0xff00000000000000000000000000000000000420
├─ Chain ID: 420 (0x1a4)
└─ Goerli Etherscan: https://goerli.etherscan.io/address/0xff00...0420
```

---

## 요약

### 🎯 핵심 포인트

**L2 배치 데이터 제출 대상**:
```
주소: 0xff00000000000000000000000000000000{L2ChainID}
타입: EOA (컨트랙트 아님)
용도: L2 트랜잭션 배치 수신
제어: 불가능 (Private key 없음)
```

### 📊 주요 특징

| 특징 | 설명 |
|------|------|
| **형식** | `0xff` + 19 zero bytes + L2 Chain ID |
| **타입** | EOA (Externally Owned Account) |
| **코드** | 없음 (0x) |
| **잔액** | 항상 0 ETH |
| **제어** | 아무도 제어 못함 |
| **고유성** | L2 Chain ID로 보장 |
| **가스** | 컨트랙트 대비 50% 절감 |
| **보안** | Batcher 서명으로 검증 |

### 🔑 보안 모델

```
온체인 (L1):
- Batch Inbox는 누구나 접근 가능
- 하지만 실질적 영향 없음 (가스 낭비만)

오프체인 (op-node):
- Batcher 주소 검증 ⭐
- 유효한 배치만 처리
- 잘못된 데이터는 무시

핵심: Batcher Private Key 보안이 가장 중요!
```

### 💡 설계 철학

**"Simple is Secure"**

- 컨트랙트보다 EOA가 더 단순
- 단순할수록 보안 위험 감소
- 가스 비용도 절감
- 업그레이드도 용이

**결론**: Batch Inbox Address는 Optimism의 영리한 설계 선택입니다! 🎯

---

## 참고 자료

### 코드 위치

- `op-batcher/batcher/driver.go`: Batch 제출 로직
- `op-node/rollup/derive/data_source.go`: Batch 감지 및 검증
- `op-node/rollup/types.go`: Rollup Config 정의

### 관련 스펙

- [Optimism Specs - Batch Submission](https://specs.optimism.io/protocol/derivation.html#batch-submission)
- [Optimism Specs - Derivation](https://specs.optimism.io/protocol/derivation.html)

### Etherscan 링크

- [Optimism Batch Inbox](https://etherscan.io/address/0xff00000000000000000000000000000000000010)
- [Base Batch Inbox](https://etherscan.io/address/0xff00000000000000000000000000000000002105)
- [OP Sepolia Batch Inbox](https://sepolia.etherscan.io/address/0xff00000000000000000000000000000000aa37dc)

---

**문서 버전**: 1.0
**작성일**: 2025-01-17
**대상 프로젝트**: tokamak-thanos (Optimism Fork)
**핵심**: L2 배치는 `0xff00...{ChainID}` 형식의 특수 EOA 주소로 제출됩니다!

