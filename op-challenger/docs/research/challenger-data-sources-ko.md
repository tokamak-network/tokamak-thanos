# op-challenger 데이터 소스 분석

## 목차
1. [핵심 질문과 답변](#핵심-질문과-답변)
2. [데이터 소스 개요](#데이터-소스-개요)
3. [코드 분석](#코드-분석)
4. [Fraud Proof 생성 전체 흐름](#fraud-proof-생성-전체-흐름)
5. [각 데이터 소스별 상세](#각-데이터-소스별-상세)
6. [실제 예시](#실제-예시)

---

## 핵심 질문과 답변

### Q: Challenger가 Fraud Proof를 생성할 때 원본 데이터를 어디서 가져오는가?

### A: 세 곳에서 가져옵니다!

```
┌──────────────────────────────────────────┐
│ op-challenger                            │
│ └─ op-program (Fault Proof 실행)        │
└────────┬─────────────────────────────────┘
         │
         ├──► ① L1 RPC (Execution Layer)
         │     - L1 블록 헤더
         │     - L1 트랜잭션 (Batch Inbox)
         │     - L1 Receipts (Deposits)
         │
         ├──► ② L1 Beacon API
         │     - Blob sidecars
         │     - L2 트랜잭션 배치 데이터 ⭐
         │
         └──► ③ L2 RPC
               - L2 블록 헤더
               - L2 트랜잭션
               - L2 상태 (계정, 스토리지)
               - L2 컨트랙트 코드
```

**핵심**: L2 트랜잭션 원본 데이터는 **L1 Beacon Chain의 Blob**에서 가져옵니다!

---

## 데이터 소스 개요

### 전체 아키텍처

```
┌─────────────────────────────────────────────────────────┐
│                     op-challenger                       │
├─────────────────────────────────────────────────────────┤
│  DisputeGame 모니터링 → 의심스러운 claim 발견           │
└──────────────────┬──────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│              Cannon/Asterisc Executor                   │
├─────────────────────────────────────────────────────────┤
│  MIPS VM 실행 → op-program 바이너리 실행                │
└──────────────────┬──────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│                    op-program                           │
├─────────────────────────────────────────────────────────┤
│  L2 상태 전이 재실행 (op-node와 동일한 derivation 로직) │
│  ↓                                                       │
│  필요한 데이터를 Hint로 요청                             │
└──────────────────┬──────────────────────────────────────┘
                   │ Hint: "l1-blob 0x..."
                   │ Hint: "l2-state-node 0x..."
                   ▼
┌─────────────────────────────────────────────────────────┐
│                    Prefetcher (Host)                    │
├─────────────────────────────────────────────────────────┤
│  Hint 해석 → 적절한 데이터 소스에서 조회                 │
└──────────────────┬──────────────────────────────────────┘
                   │
      ┌────────────┼────────────┐
      │            │            │
      ▼            ▼            ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│ L1 RPC   │ │L1 Beacon │ │ L2 RPC   │
│          │ │   API    │ │          │
└──────────┘ └──────────┘ └──────────┘
```

---

## 코드 분석

### 1. op-program 설정 (데이터 소스 지정)

**파일**: `op-program/host/config/config.go`

```go
type Config struct {
    // ✅ L1 Execution Layer (RPC)
    L1Head      common.Hash
    L1URL       string  // 예: "https://mainnet.infura.io/v3/..."
    L1TrustRPC  bool
    L1RPCKind   sources.RPCProviderKind

    // ✅ L1 Beacon Chain (API)
    L1BeaconURL string  // 예: "https://beacon-nd-123.p2pify.com"

    // ✅ L2 Execution Layer (RPC)
    L2Head       common.Hash
    L2OutputRoot common.Hash
    L2URL        string  // 예: "http://op-geth:8545"
    L2Claim      common.Hash
    L2ClaimBlockNumber uint64
    L2ChainConfig *params.ChainConfig

    DataDir string  // 로컬 preimage 캐시
}
```

**실제 설정** (Docker Compose):
```yaml
# ops-bedrock/docker-compose.yml
op-challenger:
  environment:
    OP_CHALLENGER_L1_ETH_RPC: http://l1:8545           # ← L1 RPC
    OP_CHALLENGER_L1_BEACON: http://beacon:5052        # ← L1 Beacon (있으면)
    OP_CHALLENGER_L2_ETH_RPC: http://l2:8545          # ← L2 RPC
    OP_CHALLENGER_ROLLUP_RPC: http://op-node:8545
    OP_CHALLENGER_CANNON_SERVER: /op-program/op-program  # op-program 경로
```

**Cannon 실행 시 전달**:
```go
// op-challenger/game/fault/trace/cannon/executor.go:92-111
args := []string{
    "run",
    "--",
    e.server,  // op-program 바이너리
    "--server",
    "--l1", e.l1,                    // ← L1 RPC URL
    "--l1.beacon", e.l1Beacon,       // ← L1 Beacon URL
    "--l2", e.l2,                    // ← L2 RPC URL
    "--datadir", dataDir,
    "--l1.head", e.inputs.L1Head.Hex(),
    "--l2.head", e.inputs.L2Head.Hex(),
    "--l2.outputroot", e.inputs.L2OutputRoot.Hex(),
    "--l2.claim", e.inputs.L2Claim.Hex(),
    "--l2.blocknumber", e.inputs.L2BlockNumber.Text(10),
}
```

---

### 2. Prefetcher 초기화 (데이터 소스 연결)

**파일**: `op-program/host/host.go`

```go
// Line 185-212: Prefetcher 생성
func makePrefetcher(ctx context.Context, logger log.Logger, kv kvstore.KV, cfg *config.Config) (*prefetcher.Prefetcher, error) {

    // ① L1 RPC 연결
    logger.Info("Connecting to L1 node", "l1", cfg.L1URL)
    l1RPC, err := client.NewRPC(ctx, logger, cfg.L1URL, client.WithDialBackoff(10))
    if err != nil {
        return nil, fmt.Errorf("failed to setup L1 RPC: %w", err)
    }

    // ② L2 RPC 연결
    logger.Info("Connecting to L2 node", "l2", cfg.L2URL)
    l2RPC, err := client.NewRPC(ctx, logger, cfg.L2URL, client.WithDialBackoff(10))
    if err != nil {
        return nil, fmt.Errorf("failed to setup L2 RPC: %w", err)
    }

    // L1 클라이언트 생성 (Execution Layer)
    l1ClCfg := sources.L1ClientDefaultConfig(cfg.Rollup, cfg.L1TrustRPC, cfg.L1RPCKind)
    l1Cl, err := sources.NewL1Client(l1RPC, logger, nil, l1ClCfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create L1 client: %w", err)
    }

    // ③ L1 Beacon API 클라이언트 생성
    l1Beacon := sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(cfg.L1BeaconURL, logger))
    l1BlobFetcher := sources.NewL1BeaconClient(l1Beacon, sources.L1BeaconClientConfig{FetchAllSidecars: false})

    // L2 클라이언트 생성
    l2Cl, err := NewL2Client(l2RPC, logger, nil, &L2ClientConfig{L2ClientConfig: l2ClCfg, L2Head: cfg.L2Head})
    if err != nil {
        return nil, fmt.Errorf("failed to create L2 client: %w", err)
    }
    l2DebugCl := &L2Source{L2Client: l2Cl, DebugClient: sources.NewDebugClient(l2RPC.CallContext)}

    // ✅ 세 가지 데이터 소스를 모두 Prefetcher에 주입!
    return prefetcher.NewPrefetcher(logger, l1Cl, l1BlobFetcher, l2DebugCl, kv), nil
}
```

---

### 3. Prefetcher 구조

**파일**: `op-program/host/prefetcher/prefetcher.go`

```go
// Line 37-54: 인터페이스 정의
type L1Source interface {
    InfoByHash(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, error)
    InfoAndTxsByHash(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, types.Transactions, error)
    FetchReceipts(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, types.Receipts, error)
}

type L1BlobSource interface {
    GetBlobSidecars(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.BlobSidecar, error)
    GetBlobs(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.Blob, error)
}

type L2Source interface {
    InfoAndTxsByHash(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, types.Transactions, error)
    NodeByHash(ctx context.Context, hash common.Hash) ([]byte, error)
    CodeByHash(ctx context.Context, hash common.Hash) ([]byte, error)
    OutputByRoot(ctx context.Context, root common.Hash) (eth.Output, error)
}

// Line 55-62: Prefetcher 구조체
type Prefetcher struct {
    logger        log.Logger
    l1Fetcher     L1Source        // ← L1 RPC 클라이언트
    l1BlobFetcher L1BlobSource    // ← L1 Beacon 클라이언트
    l2Fetcher     L2Source        // ← L2 RPC 클라이언트
    lastHint      string
    kvStore       kvstore.KV      // 로컬 캐시 (Key-Value Store)
}
```

---

### 4. 데이터 조회 로직 (Hint 기반)

**파일**: `op-program/host/prefetcher/prefetcher.go`

```go
// Line 99-258: prefetch 함수
func (p *Prefetcher) prefetch(ctx context.Context, hint string) error {
    hintType, hintBytes, err := parseHint(hint)

    p.logger.Debug("Prefetching", "type", hintType, "bytes", hexutil.Bytes(hintBytes))

    switch hintType {

    // ==========================================
    // ① L1 RPC에서 가져오는 데이터
    // ==========================================

    case l1.HintL1BlockHeader:
        // L1 블록 헤더 조회
        hash := common.Hash(hintBytes)
        header, err := p.l1Fetcher.InfoByHash(ctx, hash)  // ⭐ L1 RPC 호출!
        if err != nil {
            return fmt.Errorf("failed to fetch L1 block %s header: %w", hash, err)
        }
        data, err := header.HeaderRLP()
        return p.kvStore.Put(preimage.Keccak256Key(hash).PreimageKey(), data)

    case l1.HintL1Transactions:
        // L1 트랜잭션 조회 (Batch Inbox 트랜잭션 포함!)
        hash := common.Hash(hintBytes)
        _, txs, err := p.l1Fetcher.InfoAndTxsByHash(ctx, hash)  // ⭐ L1 RPC 호출!
        if err != nil {
            return fmt.Errorf("failed to fetch L1 block %s txs: %w", hash, err)
        }
        return p.storeTransactions(txs)

    case l1.HintL1Receipts:
        // L1 Receipt 조회 (Deposit 이벤트 등)
        hash := common.Hash(hintBytes)
        _, receipts, err := p.l1Fetcher.FetchReceipts(ctx, hash)  // ⭐ L1 RPC 호출!
        if err != nil {
            return fmt.Errorf("failed to fetch L1 block %s receipts: %w", hash, err)
        }
        return p.storeReceipts(receipts)

    // ==========================================
    // ② L1 Beacon API에서 가져오는 데이터
    // ==========================================

    case l1.HintL1Blob:
        // Blob 데이터 조회 (가장 중요!)
        if len(hintBytes) != 48 {
            return fmt.Errorf("invalid blob hint: %x", hint)
        }

        // Hint 파싱
        blobVersionHash := common.Hash(hintBytes[:32])   // Blob versioned hash
        blobHashIndex := binary.BigEndian.Uint64(hintBytes[32:40])  // Blob index
        refTimestamp := binary.BigEndian.Uint64(hintBytes[40:48])   // L1 timestamp

        indexedBlobHash := eth.IndexedBlobHash{
            Hash:  blobVersionHash,
            Index: blobHashIndex,
        }

        // ⭐ L1 Beacon API 호출!
        // GET /eth/v1/beacon/blob_sidecars/{slot}
        sidecars, err := p.l1BlobFetcher.GetBlobSidecars(ctx,
            eth.L1BlockRef{Time: refTimestamp},
            []eth.IndexedBlobHash{indexedBlobHash})

        if err != nil || len(sidecars) != 1 {
            return fmt.Errorf("failed to fetch blob sidecars for %s %d: %w",
                blobVersionHash, blobHashIndex, err)
        }
        sidecar := sidecars[0]

        // KZG commitment 저장
        if err = p.kvStore.Put(preimage.Sha256Key(blobVersionHash).PreimageKey(),
            sidecar.KZGCommitment[:]); err != nil {
            return err
        }

        // Blob의 4096개 field elements를 모두 저장
        blobKey := make([]byte, 80)
        copy(blobKey[:48], sidecar.KZGCommitment[:])
        for i := 0; i < params.BlobTxFieldElementsPerBlob; i++ {  // 4096
            binary.BigEndian.PutUint64(blobKey[72:], uint64(i))
            blobKeyHash := crypto.Keccak256Hash(blobKey)

            // Field element 키 저장
            if err := p.kvStore.Put(preimage.Keccak256Key(blobKeyHash).PreimageKey(), blobKey); err != nil {
                return err
            }

            // Field element 데이터 저장 (32 bytes씩)
            if err = p.kvStore.Put(preimage.BlobKey(blobKeyHash).PreimageKey(),
                sidecar.Blob[i<<5:(i+1)<<5]); err != nil {
                return err
            }
        }
        return nil

    // ==========================================
    // ③ L2 RPC에서 가져오는 데이터
    // ==========================================

    case l2.HintL2BlockHeader, l2.HintL2Transactions:
        // L2 블록 헤더 및 트랜잭션 조회
        hash := common.Hash(hintBytes)
        header, txs, err := p.l2Fetcher.InfoAndTxsByHash(ctx, hash)  // ⭐ L2 RPC 호출!
        if err != nil {
            return fmt.Errorf("failed to fetch L2 block %s: %w", hash, err)
        }

        // 헤더 저장
        data, err := header.HeaderRLP()
        err = p.kvStore.Put(preimage.Keccak256Key(hash).PreimageKey(), data)

        // 트랜잭션 저장
        return p.storeTransactions(txs)

    case l2.HintL2StateNode:
        // L2 상태 트리 노드 조회
        hash := common.Hash(hintBytes)
        node, err := p.l2Fetcher.NodeByHash(ctx, hash)  // ⭐ L2 RPC 호출!
        if err != nil {
            return fmt.Errorf("failed to fetch L2 state node %s: %w", hash, err)
        }
        return p.kvStore.Put(preimage.Keccak256Key(hash).PreimageKey(), node)

    case l2.HintL2Code:
        // L2 컨트랙트 코드 조회
        hash := common.Hash(hintBytes)
        code, err := p.l2Fetcher.CodeByHash(ctx, hash)  // ⭐ L2 RPC 호출!
        if err != nil {
            return fmt.Errorf("failed to fetch L2 contract code %s: %w", hash, err)
        }
        return p.kvStore.Put(preimage.Keccak256Key(hash).PreimageKey(), code)

    case l2.HintL2Output:
        // L2 Output Root 조회
        hash := common.Hash(hintBytes)
        output, err := p.l2Fetcher.OutputByRoot(ctx, hash)  // ⭐ L2 RPC 호출!
        if err != nil {
            return fmt.Errorf("failed to fetch L2 output root %s: %w", hash, err)
        }
        return p.kvStore.Put(preimage.Keccak256Key(hash).PreimageKey(), output.Marshal())
    }

    return fmt.Errorf("unknown hint type: %v", hintType)
}
```

---

### 5. op-program의 Derivation 실행

**파일**: `op-program/client/program.go`

```go
// Line 41-83: op-program 메인 로직
func RunProgram(logger log.Logger, preimageOracle io.ReadWriter, preimageHinter io.ReadWriter) error {
    // Oracle 클라이언트 생성
    pClient := preimage.NewOracleClient(preimageOracle)
    hClient := preimage.NewHintWriter(preimageHinter)

    // L1/L2 Oracle 래퍼 생성
    l1PreimageOracle := l1.NewCachingOracle(l1.NewPreimageOracle(pClient, hClient))
    l2PreimageOracle := l2.NewCachingOracle(l2.NewPreimageOracle(pClient, hClient))

    // Bootstrap 정보 로드
    bootInfo := NewBootstrapClient(pClient).BootInfo()
    logger.Info("Program Bootstrapped", "bootInfo", bootInfo)

    return runDerivation(
        logger,
        bootInfo.RollupConfig,
        bootInfo.L2ChainConfig,
        bootInfo.L1Head,
        bootInfo.L2OutputRoot,
        bootInfo.L2Claim,
        bootInfo.L2ClaimBlockNumber,
        l1PreimageOracle,  // ← L1 데이터 요청 인터페이스
        l2PreimageOracle,  // ← L2 데이터 요청 인터페이스
    )
}

func runDerivation(logger log.Logger, cfg *rollup.Config, l2Cfg *params.ChainConfig,
    l1Head common.Hash, l2OutputRoot common.Hash, l2Claim common.Hash,
    l2ClaimBlockNum uint64, l1Oracle l1.Oracle, l2Oracle l2.Oracle) error {

    // L1 데이터 소스 생성
    l1Source := l1.NewOracleL1Client(logger, l1Oracle, l1Head)
    l1BlobsSource := l1.NewBlobFetcher(logger, l1Oracle)

    // L2 데이터 소스 생성
    engineBackend, err := l2.NewOracleBackedL2Chain(logger, l2Oracle, l1Oracle, l2Cfg, l2OutputRoot)
    l2Source := l2.NewOracleEngine(cfg, logger, engineBackend)

    // ✅ Derivation 파이프라인 실행 (op-node와 동일!)
    logger.Info("Starting derivation")
    d := cldr.NewDriver(logger, cfg, l1Source, l1BlobsSource, l2Source, l2ClaimBlockNum)

    // L1에서 배치 데이터를 읽어 L2 블록 재구성
    for {
        if err = d.Step(context.Background()); errors.Is(err, io.EOF) {
            break
        } else if err != nil {
            return err
        }
    }

    // 최종 Claim 검증
    return d.ValidateClaim(l2ClaimBlockNum, eth.Bytes32(l2Claim))
}
```

---

### 6. L1 Blob 조회 상세

**파일**: `op-program/client/l1/blob_fetcher.go`

```go
// Line 14-36: BlobFetcher 구현
type BlobFetcher struct {
    logger log.Logger
    oracle Oracle  // Preimage Oracle (Host에서 제공)
}

func (b *BlobFetcher) GetBlobs(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.Blob, error) {
    blobs := make([]*eth.Blob, len(hashes))
    for i := 0; i < len(hashes); i++ {
        b.logger.Info("Fetching blob",
            "l1_ref", ref.Hash,
            "blob_versioned_hash", hashes[i].Hash,
            "index", hashes[i].Index)

        // ⭐ Oracle를 통해 blob 조회
        // 내부적으로 Hint를 Host에 전송
        // Host의 Prefetcher가 Beacon API 호출
        blobs[i] = b.oracle.GetBlob(ref, hashes[i])
    }
    return blobs, nil
}
```

**실제 Beacon API 호출** (`op-service/sources/l1_beacon_client.go`):

```go
// Line 240-278: Beacon API로 blob sidecars 조회
func (cl *L1BeaconClient) GetBlobSidecars(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.BlobSidecar, error) {
    // L1 블록 timestamp → Beacon slot 변환
    slotFn, err := cl.GetTimeToSlotFn(ctx)
    slot, err := slotFn(ref.Time)

    // ⭐ Beacon API HTTP 요청!
    // GET /eth/v1/beacon/blob_sidecars/{slot}
    resp, err := cl.fetchSidecars(ctx, slot, hashes)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch blob sidecars for slot %v block %v: %w", slot, ref, err)
    }

    // 응답 필터링 및 정렬
    apiscs := make([]*eth.APIBlobSidecar, 0, len(hashes))
    for _, h := range hashes {
        for _, apisc := range resp.Data {
            if h.Index == uint64(apisc.Index) {
                apiscs = append(apiscs, apisc)
                break
            }
        }
    }

    // BlobSidecar 객체로 변환
    bscs := make([]*eth.BlobSidecar, 0, len(hashes))
    for _, apisc := range apiscs {
        bscs = append(bscs, apisc.BlobSidecar())
    }

    return bscs, nil
}
```

---

## Fraud Proof 생성 전체 흐름

### 상황: 잘못된 Output Root 발견

```
Proposer가 제안:
  L2 Block #1,000,000
  Output Root: 0xBADROOT...

op-challenger 검증:
  자체 계산: 0xGOODROOT...
  → 불일치 발견! 🚨
  → Fraud Proof 생성 시작
```

### 단계별 데이터 조회

#### Phase 1: 게임 설정

```bash
# op-challenger가 Cannon 실행
cannon run \
  --input /data/games/{addr}/prestate.json \
  -- \
  /op-program/op-program --server \
    --l1 https://mainnet.infura.io/v3/...          # ← L1 RPC
    --l1.beacon https://beacon-nd-123.p2pify.com   # ← L1 Beacon
    --l2 http://op-geth:8545                       # ← L2 RPC
    --l1.head 0xL1HEAD... \
    --l2.head 0xL2HEAD... \
    --l2.claim 0xBADROOT... \
    --l2.blocknumber 1000000
```

#### Phase 2: L1 배치 데이터 조회

```
1. op-program이 L1 블록 조회 요청
   └─ Hint: "l1-block-header 0xL1HEAD..."

2. Prefetcher가 L1 RPC 호출
   └─ POST https://mainnet.infura.io/v3/...
   └─ Method: eth_getBlockByHash
   └─ Params: ["0xL1HEAD...", false]

3. L1 블록 헤더 수신
   └─ {
       "number": "0x129d20b",
       "hash": "0xL1HEAD...",
       "timestamp": "0x65f4a8c0",
       ...
     }

4. KV Store에 저장
   └─ Key: Keccak256("0xL1HEAD...")
   └─ Value: RLP(header)
```

#### Phase 3: L1 트랜잭션 조회 (Batch Inbox)

```
1. op-program이 L1 트랜잭션 요청
   └─ Hint: "l1-transactions 0xL1HEAD..."

2. Prefetcher가 L1 RPC 호출
   └─ POST https://mainnet.infura.io/v3/...
   └─ Method: eth_getBlockByHash
   └─ Params: ["0xL1HEAD...", true]  # true = 트랜잭션 포함

3. L1 블록 및 트랜잭션 수신
   └─ {
       "transactions": [
         {
           "hash": "0xTX1...",
           "to": "0x1234...",  # 일반 트랜잭션
           "type": "0x2",
           "input": "0x..."
         },
         {
           "hash": "0xBATCH...",
           "to": "0xff00000000000000000000000000000000000420",  # ⭐ Batch Inbox!
           "type": "0x3",  # EIP-4844 Blob Transaction
           "blobVersionedHashes": [
             "0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8",
             "0x01ab2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b"
           ]
         }
       ]
     }

4. Batch 트랜잭션 식별
   └─ To == BatchInboxAddress
   └─ Type == 3 (Blob)
   └─ BlobVersionedHashes 추출
```

#### Phase 4: Blob 데이터 조회 (핵심!)

```
1. op-program이 Blob 요청
   └─ Hint: "l1-blob 0x01fa3b84...{index}{timestamp}"
   └─ blobVersionedHash: 0x01fa3b84...
   └─ index: 0
   └─ timestamp: 1710950400

2. Prefetcher가 Timestamp → Slot 변환
   └─ Genesis Time: 1606824023
   └─ Slot = (1710950400 - 1606824023) / 12 = 8,677,198

3. Prefetcher가 L1 Beacon API 호출 ⭐
   └─ GET https://beacon-nd-123.p2pify.com/eth/v1/beacon/blob_sidecars/8677198

4. Beacon API 응답 수신
   └─ {
       "data": [
         {
           "index": "0",
           "blob": "0x00789c5d915d6f1c4710c73f...",  # ← L2 트랜잭션 배치 데이터!
                    └─ 131,071 bytes (압축된 L2 트랜잭션들)
           "kzg_commitment": "0x01fa3b84e98e6f3c...",
           "kzg_proof": "0x9a8b7c6d5e4f3a2b..."
         }
       ]
     }

5. Blob 데이터를 4096개 Field Elements로 분할
   └─ 각 32 bytes씩 KV Store에 저장
   └─ Key: Keccak256(KZGCommitment || index)
   └─ Value: blob[i*32 : (i+1)*32]

6. op-program이 Blob 데이터 수신
   └─ 4096개 field elements 조합
   └─ 압축 해제 (brotli/zlib)
   └─ L2 트랜잭션 400개 추출 ✅
```

#### Phase 5: L2 상태 재실행

```
1. op-program이 L2 블록 조회 요청
   └─ Hint: "l2-block-header 0xL2BLOCK..."

2. Prefetcher가 L2 RPC 호출
   └─ POST http://op-geth:8545
   └─ Method: eth_getBlockByHash
   └─ Params: ["0xL2BLOCK...", true]

3. L2 블록 데이터 수신
   └─ {
       "number": "0xf4240",  # 1,000,000
       "hash": "0xL2BLOCK...",
       "stateRoot": "0xSTATE...",
       "transactions": [...]
     }

4. 각 L2 트랜잭션 실행
   ├─ 컨트랙트 호출?
   │  └─ Hint: "l2-code 0xCODEHASH..."
   │  └─ Prefetcher → L2 RPC: eth_getCode
   │  └─ 컨트랙트 바이트코드 수신
   │
   ├─ Storage 읽기?
   │  └─ Hint: "l2-state-node 0xSTATEHASH..."
   │  └─ Prefetcher → L2 RPC: debug_dbGet
   │  └─ 상태 트리 노드 수신
   │
   └─ 트랜잭션 실행
      └─ EVM 실행 (MIPS VM 내부에서!)

5. 모든 트랜잭션 실행 완료
   └─ 최종 State Root 계산
   └─ Output Root 생성
   └─ Output Root: 0xGOODROOT...

6. Claim 검증
   └─ Claimed: 0xBADROOT...
   └─ Computed: 0xGOODROOT...
   └─ Mismatch! → Fraud 증명 완료 ✅
```

---

## 실제 예시: 구체적인 데이터 흐름

### 시나리오: Block #1,000,000 검증

```
게임 정보:
- L2 Block Number: 1,000,000
- Claimed Output Root: 0xBAD123...
- L1 Head: Block #19,456,789
- L2 Head: Block #999,999
```

### 단계별 데이터 조회

#### 1. L1 블록 헤더 조회

```
Hint 요청:
  "l1-block-header 0xabcd1234..."

Prefetcher 동작:
  → POST https://mainnet.infura.io/v3/YOUR_KEY
  → {
      "method": "eth_getBlockByHash",
      "params": ["0xabcd1234...", false]
    }

응답:
  {
    "number": "0x129d20b",        # 19,456,779
    "hash": "0xabcd1234...",
    "timestamp": "0x65f4a8c0",    # 1710950592
    "parentHash": "0x...",
    "stateRoot": "0x...",
    "receiptsRoot": "0x...",
    "baseFeePerGas": "0x4a817c800"  # 20 Gwei
  }

저장:
  Key: Keccak256("0xabcd1234...")
  Value: RLP(header)
```

#### 2. L1 트랜잭션 조회 (Batch 포함)

```
Hint 요청:
  "l1-transactions 0xabcd1234..."

Prefetcher 동작:
  → POST https://mainnet.infura.io/v3/YOUR_KEY
  → {
      "method": "eth_getBlockByHash",
      "params": ["0xabcd1234...", true]  # 트랜잭션 포함
    }

응답:
  {
    "transactions": [
      ... 일반 트랜잭션들 ...
      {
        "hash": "0x1a2b3c4d...",
        "from": "0xBATCHER...",
        "to": "0xff00000000000000000000000000000000000420",  # ⭐ Batch Inbox!
        "type": "0x3",  # Blob Transaction
        "blobVersionedHashes": [
          "0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8",
          "0x01ab2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b"
        ],
        "maxFeePerBlobGas": "0x2540be400",
        "blockNumber": "0x129d20b"
      }
    ]
  }

처리:
  → Batch Inbox 트랜잭션 식별
  → BlobVersionedHashes 추출: [0x01fa3b84..., 0x01ab2c3d...]
  → 각 Blob에 대해 별도 조회 필요!
```

#### 3. Blob 데이터 조회 (L2 트랜잭션 배치)

```
Hint 요청 #1:
  "l1-blob 0x01fa3b84...{0}{1710950400}"
          └─ blob hash
                     └─ index
                           └─ timestamp

Prefetcher 동작:
  1. Timestamp → Slot 변환
     └─ slot = (1710950400 - 1606824023) / 12 = 8,677,198

  2. Beacon API 호출
     → GET https://beacon-nd-123.p2pify.com/eth/v1/beacon/blob_sidecars/8677198
     → 또는: GET .../blob_sidecars/8677198?indices=0

응답 (Blob Sidecar):
  {
    "data": [
      {
        "index": "0",
        "blob": "0x00789c5d915d6f1c4710c73f...",  # ⭐ 131,071 bytes
        "kzg_commitment": "0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4",
        "kzg_proof": "0x9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d3e2f1a0b9c8d7e6f5a4b3c2d1e0f9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d"
      }
    ]
  }

저장:
  1. KZG Commitment 저장
     Key: SHA256("0x01fa3b84...")
     Value: kzg_commitment (48 bytes)

  2. Blob을 4096개 field elements로 분할
     for i in 0..4095:
       blobKey = KZGCommitment || i
       Key: Keccak256(blobKey)
       Value: blob[i*32 : (i+1)*32]  # 32 bytes씩

Hint 요청 #2:
  "l1-blob 0x01ab2c3d...{1}{1710950400}"

  → 동일한 과정 반복 (Blob index 1)
```

#### 4. Blob 데이터 압축 해제 및 파싱

```
op-program이 Blob field elements 조회:
  for i in 0..4095:
    preimage := GetPreimage(Keccak256(KZGCommitment || i))
    blob[i*32 : (i+1)*32] = preimage

Blob 데이터 완성:
  blob = [131,071 bytes]

압축 해제:
  └─ Version byte: 0x00 (첫 바이트)
  └─ 압축 데이터: 0x789c5d915d6f... (brotli 또는 zlib)
  └─ Decompress(blob[1:])
  └─ 원본 데이터: [L2 트랜잭션 배치]

L2 트랜잭션 파싱:
  └─ Frame 형식 파싱
  └─ Channel 재구성
  └─ Batch 디코딩
  └─ L2 트랜잭션 400개 추출 ✅

획득한 L2 트랜잭션 예시:
  Transaction 1:
    From: 0xUSER1...
    To: 0xCONTRACT...
    Value: 1 ETH
    Data: 0x...

  Transaction 2:
    From: 0xUSER2...
    To: 0xUNISWAP...
    Data: swap(...)

  ... (총 400개)
```

#### 5. L2 상태 조회

```
L2 트랜잭션 실행 중:
  Transaction: transferFrom(0xUSER1, 0xUSER2, 100 tokens)

1. 컨트랙트 코드 필요
   └─ Hint: "l2-code 0xCODEHASH..."

2. Prefetcher가 L2 RPC 호출
   → POST http://op-geth:8545
   → {
       "method": "eth_getCode",
       "params": ["0xCONTRACT...", "0xL2BLOCK..."]
     }

3. 컨트랙트 바이트코드 수신
   → "0x608060405234801561001057600080fd5b50..."
   → ERC20 컨트랙트 바이트코드

4. 잔액 조회 필요 (Storage)
   └─ Hint: "l2-state-node 0xSTATEHASH..."

5. Prefetcher가 L2 RPC 호출
   → POST http://op-geth:8545
   → {
       "method": "debug_dbGet",
       "params": ["0xSTATEHASH..."]
     }

6. 상태 트리 노드 수신
   → RLP encoded state trie node
   → USER1 balance: 1000 tokens

7. 트랜잭션 실행
   → EVM 실행 (MIPS VM 내부)
   → Storage 업데이트
   → USER1 balance: 1000 - 100 = 900
   → USER2 balance: 100
```

#### 6. 최종 검증

```
모든 L2 트랜잭션 (Block #999,999 → #1,000,000) 실행 완료:

계산된 Output:
  StateRoot: 0xGOODSTATE...
  WithdrawalStorageRoot: 0x...
  BlockHash: 0xL2BLOCK...

  Output Root = Keccak256(
    version || StateRoot || WithdrawalStorageRoot || BlockHash
  )
  = 0xGOODROOT...

Claimed Output Root: 0xBADROOT...

비교:
  0xGOODROOT... != 0xBADROOT...

결론: FRAUD 발견! 🚨
```

---

## 데이터 소스별 상세 분석

### 1. L1 RPC (Execution Layer)

**제공 데이터**:
- ✅ L1 블록 헤더 (번호, 해시, timestamp, baseFee 등)
- ✅ L1 트랜잭션 (Batch Inbox로 전송된 트랜잭션)
- ✅ L1 Receipts (Deposit 이벤트, 로그 등)
- ✅ Blob versioned hashes (Blob 트랜잭션의 경우)

**RPC 메서드**:
```bash
# 블록 조회
eth_getBlockByHash(hash, includeTransactions)
eth_getBlockByNumber(number, includeTransactions)

# 트랜잭션 조회
eth_getTransactionByHash(hash)
eth_getTransactionReceipt(hash)

# 상태 조회 (일반적으로 사용 안 함)
eth_getBalance(address, blockHash)
eth_getCode(address, blockHash)
```

**코드 구현**:
```go
// op-service/sources/l1_client.go
type L1Client struct {
    client client.RPC
    // ...
}

func (s *L1Client) InfoByHash(ctx context.Context, hash common.Hash) (eth.BlockInfo, error) {
    // eth_getBlockByHash 호출
}

func (s *L1Client) InfoAndTxsByHash(ctx context.Context, hash common.Hash) (eth.BlockInfo, types.Transactions, error) {
    // eth_getBlockByHash(hash, true) 호출
}

func (s *L1Client) FetchReceipts(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, types.Receipts, error) {
    // 각 트랜잭션에 대해 eth_getTransactionReceipt 호출
}
```

---

### 2. L1 Beacon API

**제공 데이터**:
- ✅ Blob sidecars (L2 트랜잭션 배치 데이터)
- ✅ KZG commitments (Blob 검증용)
- ✅ KZG proofs (Blob 검증용)
- ✅ Beacon block headers

**API 엔드포인트**:
```bash
# Node 정보
GET /eth/v1/node/version

# Chain 설정
GET /eth/v1/config/spec

# Genesis
GET /eth/v1/beacon/genesis

# Blob sidecars (핵심!)
GET /eth/v1/beacon/blob_sidecars/{block_id}
GET /eth/v1/beacon/blob_sidecars/{block_id}?indices=0,1,2
```

**코드 구현**:
```go
// op-service/sources/l1_beacon_client.go
type L1BeaconClient struct {
    cl   BeaconClient
    pool *ClientPool[BlobSideCarsFetcher]
    cfg  L1BeaconClientConfig
    timeToSlotFn TimeToSlotFn
}

func (cl *L1BeaconClient) GetBlobSidecars(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.BlobSidecar, error) {
    // 1. Timestamp → Slot 변환
    slotFn, err := cl.GetTimeToSlotFn(ctx)
    slot, err := slotFn(ref.Time)

    // 2. Beacon API HTTP 요청
    resp, err := cl.fetchSidecars(ctx, slot, hashes)

    // 3. 응답 파싱 및 반환
    return blobSidecars, nil
}

func (cl *L1BeaconClient) GetBlobs(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.Blob, error) {
    // GetBlobSidecars 호출 후 검증
    blobSidecars, err := cl.GetBlobSidecars(ctx, ref, hashes)

    // KZG 검증
    for i, sidecar := range blobSidecars {
        // Versioned hash 검증
        hash := eth.KZGToVersionedHash(kzg4844.Commitment(sidecar.KZGCommitment))
        if hash != hashes[i].Hash {
            return nil, fmt.Errorf("hash mismatch")
        }

        // KZG Proof 검증
        if err := eth.VerifyBlobProof(&sidecar.Blob, sidecar.KZGCommitment, sidecar.KZGProof); err != nil {
            return nil, fmt.Errorf("blob verification failed: %w", err)
        }
    }

    return blobs, nil
}
```

---

### 3. L2 RPC

**제공 데이터**:
- ✅ L2 블록 헤더 및 트랜잭션
- ✅ L2 상태 트리 노드 (계정 잔액, 스토리지)
- ✅ L2 컨트랙트 바이트코드
- ✅ L2 Output Root

**RPC 메서드**:
```bash
# 표준 메서드
eth_getBlockByHash(hash, includeTransactions)
eth_getBlockByNumber(number, includeTransactions)
eth_getCode(address, blockHash)

# Debug 메서드 (필수!)
debug_dbGet(key)              # 상태 트리 노드 직접 조회
optimism_outputAtBlock(number)  # L2 Output Root 조회
```

**코드 구현**:
```go
// op-program/host/l2_client.go
type L2Client struct {
    *sources.L2Client
    l2Head common.Hash
}

func (s *L2Client) OutputByRoot(ctx context.Context, l2OutputRoot common.Hash) (eth.Output, error) {
    // optimism_outputAtBlock 호출
    output, err := s.OutputV0AtBlock(ctx, s.l2Head)

    // Output Root 검증
    actualOutputRoot := eth.OutputRoot(output)
    if actualOutputRoot != eth.Bytes32(l2OutputRoot) {
        panic(fmt.Errorf("output root mismatch"))
    }

    return output, nil
}

// op-service/sources/l2_client.go
func (s *L2Client) InfoAndTxsByHash(ctx context.Context, hash common.Hash) (eth.BlockInfo, types.Transactions, error) {
    // eth_getBlockByHash(hash, true) 호출
}

// op-service/sources/debug_client.go
func (s *DebugClient) NodeByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
    // debug_dbGet(hash) 호출
}

func (s *DebugClient) CodeByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
    // eth_getCode로 코드 조회
}
```

---

## 실제 RPC 호출 추적

### L1 RPC 호출 예시

```bash
# 1. L1 블록 헤더 조회
curl https://mainnet.infura.io/v3/YOUR_KEY \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockByHash",
    "params": ["0xabcd1234567890abcdef...", false],
    "id": 1
  }'

# 응답:
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "number": "0x129d20b",
    "hash": "0xabcd1234567890abcdef...",
    "timestamp": "0x65f4a8c0",
    "baseFeePerGas": "0x4a817c800",
    ...
  }
}

# 2. L1 트랜잭션 조회 (Batch Inbox)
curl https://mainnet.infura.io/v3/YOUR_KEY \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{
    "method": "eth_getBlockByHash",
    "params": ["0xabcd1234567890abcdef...", true]
  }'

# 응답:
{
  "result": {
    "transactions": [
      {
        "hash": "0x1a2b3c4d...",
        "to": "0xff00000000000000000000000000000000000420",
        "type": "0x3",
        "blobVersionedHashes": [
          "0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8",
          "0x01ab2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b"
        ]
      }
    ]
  }
}
```

### L1 Beacon API 호출 예시

```bash
# 1. Beacon node 버전 확인
curl https://beacon-nd-123.p2pify.com/eth/v1/node/version

# 응답:
{
  "data": {
    "version": "Lighthouse/v4.5.0-1a0223c/x86_64-linux"
  }
}

# 2. Genesis 정보 조회
curl https://beacon-nd-123.p2pify.com/eth/v1/beacon/genesis

# 응답:
{
  "data": {
    "genesis_time": "1606824023",
    "genesis_validators_root": "0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95"
  }
}

# 3. Spec 조회 (Slot 계산용)
curl https://beacon-nd-123.p2pify.com/eth/v1/config/spec

# 응답:
{
  "data": {
    "SECONDS_PER_SLOT": "12",
    "SLOTS_PER_EPOCH": "32",
    "MIN_EPOCHS_FOR_BLOB_SIDECARS_REQUESTS": "4096"
  }
}

# 4. Blob sidecars 조회 (핵심!)
curl https://beacon-nd-123.p2pify.com/eth/v1/beacon/blob_sidecars/8677198

# 또는 특정 indices만:
curl https://beacon-nd-123.p2pify.com/eth/v1/beacon/blob_sidecars/8677198?indices=0,1

# 응답:
{
  "data": [
    {
      "index": "0",
      "blob": "0x00789c5d915d6f1c4710c73f9c5f88d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7b6c5d4e3f2a1...",
      "kzg_commitment": "0x01fa3b84e98e6f3c2d1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4",
      "kzg_proof": "0x9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d3e2f1a0b9c8d7e6f5a4b3c2d1e0f9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d",
      "signed_block_header": { ... },
      "kzg_commitment_inclusion_proof": [ ... ]
    },
    {
      "index": "1",
      "blob": "0x00a1b2c3d4e5f6...",
      "kzg_commitment": "0x01ab2c3d4e5f6a7b...",
      "kzg_proof": "0x...",
      ...
    }
  ]
}
```

### L2 RPC 호출 예시

```bash
# 1. L2 블록 조회
curl http://op-geth:8545 \
  -X POST \
  -H "Content-Type: application/json" \
  --data '{
    "method": "eth_getBlockByHash",
    "params": ["0xL2BLOCK...", true]
  }'

# 응답:
{
  "result": {
    "number": "0xf4240",  # 1,000,000
    "hash": "0xL2BLOCK...",
    "stateRoot": "0xSTATE...",
    "transactionsRoot": "0x...",
    "receiptsRoot": "0x...",
    "transactions": [
      {
        "hash": "0xTX1...",
        "from": "0xUSER1...",
        "to": "0xCONTRACT...",
        "input": "0x...",
        "value": "0xde0b6b3a7640000"
      },
      ...
    ]
  }
}

# 2. 상태 트리 노드 조회 (Debug API)
curl http://op-geth:8545 \
  -X POST \
  --data '{
    "method": "debug_dbGet",
    "params": ["0xSTATEHASH..."]
  }'

# 응답:
{
  "result": "0xf84d80893635c9adc5dea00000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
  # ← RLP encoded state trie node
}

# 3. 컨트랙트 코드 조회
curl http://op-geth:8545 \
  -X POST \
  --data '{
    "method": "eth_getCode",
    "params": ["0xCONTRACT...", "0xL2BLOCK..."]
  }'

# 응답:
{
  "result": "0x608060405234801561001057600080fd5b50600436106100..."
  # ← 컨트랙트 바이트코드 (ERC20 등)
}

# 4. L2 Output Root 조회
curl http://op-geth:8545 \
  -X POST \
  --data '{
    "method": "optimism_outputAtBlock",
    "params": ["0xf4240"]  # Block #1,000,000
  }'

# 응답:
{
  "result": {
    "version": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "stateRoot": "0xGOODSTATE...",
    "withdrawalStorageRoot": "0x...",
    "blockHash": "0xL2BLOCK..."
  }
}
```

---

## 로컬 캐시 구조

### 디렉토리 구조

```
/data/games/{game_address}/
├── cannon/                      # Cannon VM 데이터
│   ├── meta.json               # 메타데이터
│   ├── state.json.gz           # 마지막 실행 상태
│   │
│   ├── snapshots/              # VM 스냅샷
│   │   ├── 0.json.gz           # 초기 상태
│   │   ├── 1000000.json.gz     # 1M step
│   │   ├── 2000000.json.gz     # 2M step
│   │   └── ...
│   │
│   ├── proofs/                 # 생성된 Proof
│   │   ├── 0.json.gz
│   │   ├── 500000.json.gz
│   │   ├── 1000000.json.gz
│   │   └── ...
│   │
│   └── preimages/              # Preimage 캐시
│       ├── 0x12/               # 첫 바이트별 디렉토리
│       │   ├── 0x1234...       # L1 블록 헤더
│       │   ├── 0x12ab...       # L1 트랜잭션
│       │   └── ...
│       ├── 0xfa/
│       │   ├── 0xfa3b...       # Blob field element
│       │   └── ...
│       └── ...
```

### Preimage 저장 형식

```
Key 타입:
1. Keccak256 (Type 1): L1/L2 블록, 트랜잭션, 상태
   └─ Key: 0x01 || Keccak256(data)

2. SHA256 (Type 2): Blob KZG commitments
   └─ Key: 0x02 || SHA256(commitment)

3. Blob (Type 3): Blob field elements
   └─ Key: 0x03 || Keccak256(KZGCommitment || index)

4. Precompile (Type 5): Precompile 결과
   └─ Key: 0x05 || Keccak256(input)
```

**Proof 파일 예시**:
```json
// proofs/1000000.json.gz (압축 해제 후)
{
  "post": "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad",
  "state-data": "0x0000000000000000000000000000000000000000000000000000000000000000...",
  "proof-data": "0x0100000000f4240000000000...",
  "oracle-key": "0x0301fa3b84e98e6f3c2d1b0a9f8e7d6c5b...",
  "oracle-value": "0x00789c5d915d6f1c4710c73f...",
  "oracle-offset": 0
}
```

---

## 정리: 데이터 조회 요약

### 데이터 흐름 다이어그램

```
┌─────────────────────────────────────────────────────────┐
│                   Challenger 시작                        │
│  DisputeGame에서 의심스러운 claim 발견                   │
└──────────────────┬──────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│             Cannon/Asterisc 실행 시작                   │
│  MIPS VM 초기화 → op-program 바이너리 로드               │
└──────────────────┬──────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│                 op-program 실행                         │
│  L2 상태 전이 재실행 (Derivation Pipeline)              │
└──────────────────┬──────────────────────────────────────┘
                   │
          ┌────────┴────────┐
          │ 필요한 데이터:   │
          │ 1. L1 배치 데이터│
          │ 2. L2 상태 데이터│
          └────────┬────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│              Hint-Preimage 프로토콜                     │
│  op-program → Host에 Hint 전송                          │
│  Host → 데이터 조회 → Preimage 반환                     │
└──────────────────┬──────────────────────────────────────┘
                   │
      ┌────────────┼────────────┐
      │            │            │
      ▼            ▼            ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│ L1 RPC   │ │L1 Beacon │ │ L2 RPC   │
│          │ │   API    │ │          │
│ Infura   │ │ p2pify   │ │ op-geth  │
│ Alchemy  │ │ AllNodes │ │ (local)  │
└─────┬────┘ └─────┬────┘ └─────┬────┘
      │            │            │
      └────────────┴────────────┘
                   │
                   ▼
         ┌─────────────────┐
         │   Preimage      │
         │   KV Store      │
         │  (로컬 캐시)     │
         └─────────────────┘
```

### 데이터 종류별 소스

| 데이터 종류 | 소스 | RPC/API | 용도 |
|------------|------|---------|------|
| **L1 블록 헤더** | L1 RPC | `eth_getBlockByHash` | Timestamp, BaseFee 등 |
| **L1 트랜잭션** | L1 RPC | `eth_getBlockByHash` | Batch Inbox 트랜잭션 식별 |
| **Blob Hashes** | L1 RPC | (트랜잭션 내부) | Blob 조회를 위한 참조 |
| **Blob 데이터** | L1 Beacon | `GET /blob_sidecars/{slot}` | ⭐ L2 트랜잭션 배치 원본 |
| **L1 Receipts** | L1 RPC | `eth_getTransactionReceipt` | Deposit 이벤트 |
| **L2 블록** | L2 RPC | `eth_getBlockByHash` | L2 블록 정보 |
| **L2 상태** | L2 RPC | `debug_dbGet` | 계정, 스토리지 |
| **L2 코드** | L2 RPC | `eth_getCode` | 컨트랙트 바이트코드 |
| **L2 Output** | L2 RPC | `optimism_outputAtBlock` | Output Root 검증 |

### 가장 중요한 데이터: L2 트랜잭션 배치

**결론**:
```
L2 트랜잭션 원본 데이터는 어디에?

→ L1 Beacon Chain의 Blob Sidecars! ⭐

조회 방법:
1. L1 RPC로 Batch Inbox 트랜잭션 조회
   → BlobVersionedHashes 추출

2. L1 Beacon API로 Blob 조회
   → GET /eth/v1/beacon/blob_sidecars/{slot}
   → Blob 데이터 수신 (131KB × N개)

3. Blob 압축 해제
   → brotli/zlib decompress
   → L2 트랜잭션 수백~수천 개 추출

4. L2 트랜잭션 재실행
   → 필요 시 L2 RPC로 상태 조회
   → EVM 실행하여 Output Root 계산
```

---

## FAQ

### Q1: L1에서 Calldata 대신 Blob을 사용하는 이유는?

**A**: 비용 때문입니다.

```
Calldata: L1 Execution Layer에 저장
- 비용: 16 gas/byte × 150,000 bytes = 2,400,000 gas
- Gas Price: 20 Gwei
- 총 비용: 0.048 ETH ≈ $86

Blob: L1 Beacon Chain에 저장
- 비용: 131,072 blob gas × 2 blobs = 262,144 blob gas
- Blob Gas Price: 1 Wei
- 총 비용: 0.00026 ETH ≈ $0.47

절감: 99.5%! 🎉
```

### Q2: Blob이 18일 후 삭제되면 Challenger는 어떻게 하나요?

**A**: Challenger는 게임 시작 시점에 데이터를 조회합니다.

```
Timeline:
Day 0: Proposer가 Output Root 제안
       └─ L1에 Blob 트랜잭션 제출

Day 1: Challenger가 게임 참여
       └─ Blob 데이터 조회 (Beacon API)
       └─ 로컬 캐시에 저장 (/data/games/{addr}/preimages/)

Day 7-30: Challenge 게임 진행
          └─ 로컬 캐시에서 데이터 사용
          └─ Beacon API는 사용 안 함 (이미 저장됨)

Day 18: Beacon Chain에서 Blob pruning
        └─ Challenger에는 영향 없음 (이미 캐시됨)
```

**핵심**: Challenger는 게임 초기에 모든 필요한 데이터를 조회하여 로컬에 캐시합니다!

### Q3: L2 RPC에 접근할 수 없으면 어떻게 되나요?

**A**: Fraud Proof 생성 불가능합니다.

```
op-challenger 실행 조건:
- ✅ L1 RPC 필수
- ✅ L1 Beacon API 필수 (Blob 사용 시)
- ✅ L2 RPC 필수

L2 RPC 없으면:
- L2 상태 데이터 조회 불가
- 컨트랙트 코드 조회 불가
- Output Root 계산 불가
- Fraud Proof 생성 실패 ❌

해결 방법:
- 자체 L2 노드 운영 (권장)
- 또는 신뢰할 수 있는 L2 RPC 제공자 사용
```

### Q4: Plasma 모드에서는 어떻게 되나요?

**A**: DA Server도 데이터 소스가 됩니다.

```
Plasma 모드 (사용 시):
  ① L1 RPC: L1 블록, Batch Inbox 트랜잭션
     └─ Plasma commitment 읽기

  ② DA Server: L2 트랜잭션 배치 데이터
     └─ GET http://da-server:3100/get/0x{commitment}

  ③ L2 RPC: L2 상태 데이터

일반 모드 (Blob):
  ① L1 RPC: L1 블록, Batch Inbox 트랜잭션
  ② L1 Beacon: L2 트랜잭션 배치 데이터
  ③ L2 RPC: L2 상태 데이터
```

### Q5: 모든 데이터를 매번 조회하나요?

**A**: 아니오, 로컬 캐시를 적극 활용합니다.

```go
// op-program/host/prefetcher/prefetcher.go:80-96
func (p *Prefetcher) GetPreimage(ctx context.Context, key common.Hash) ([]byte, error) {
    // 1. 먼저 로컬 캐시 확인
    pre, err := p.kvStore.Get(key)

    // 2. 캐시에 없으면 prefetch
    for errors.Is(err, kvstore.ErrNotFound) && p.lastHint != "" {
        hint := p.lastHint
        if err := p.prefetch(ctx, hint); err != nil {
            return nil, fmt.Errorf("prefetch failed: %w", err)
        }
        pre, err = p.kvStore.Get(key)
    }

    return pre, err
}
```

**캐시 효과**:
```
첫 번째 Proof 생성:
- L1 RPC 호출: 100회
- Beacon API 호출: 50회
- L2 RPC 호출: 1,000회
- 시간: 5분

두 번째 Proof 생성 (같은 게임):
- 모든 데이터가 캐시됨
- 외부 호출: 0회
- 시간: 10초 ✨
```

---

## 핵심 정리

### ✅ Challenger의 데이터 소스 (요약)

```
op-challenger가 Fraud Proof 생성 시:

1. L1 RPC (Execution Layer)
   └─ L1 블록, 트랜잭션 (Batch Inbox)
   └─ Blob versioned hashes 추출

2. L1 Beacon API ⭐⭐⭐
   └─ Blob sidecars 조회
   └─ L2 트랜잭션 배치 원본 데이터!
   └─ 가장 중요한 데이터 소스!

3. L2 RPC
   └─ L2 블록, 상태, 코드
   └─ Output Root 검증
```

### 🎯 원본 데이터의 위치

**Q: L2 트랜잭션 원본 데이터는 어디에?**

**A: L1 Beacon Chain의 Blob에 저장되어 있습니다!**

```
op-batcher가 제출:
  L2 Transactions → Compress → Blob → L1 Beacon Chain

op-challenger가 조회:
  L1 Beacon Chain → Blob → Decompress → L2 Transactions

데이터 흐름:
  L2 (생성) → L1 Beacon (저장) → Challenger (조회) → Fraud Proof
```

### 🔑 핵심 코드 위치

```
데이터 소스 설정:
- op-program/host/config/config.go (L1URL, L1BeaconURL, L2URL)

데이터 조회 로직:
- op-program/host/prefetcher/prefetcher.go (prefetch 함수)
- op-program/host/host.go (makePrefetcher)

Beacon API 클라이언트:
- op-service/sources/l1_beacon_client.go (GetBlobSidecars)

L1/L2 RPC 클라이언트:
- op-service/sources/l1_client.go (L1 조회)
- op-service/sources/l2_client.go (L2 조회)
- op-service/sources/debug_client.go (L2 상태 조회)
```

---

**문서 버전**: 1.0
**작성일**: 2025-01-17
**대상 프로젝트**: tokamak-thanos (Optimism Fork)
**핵심 발견**: Challenger는 L1 Beacon Chain의 Blob에서 L2 트랜잭션 원본 데이터를 가져옵니다!

