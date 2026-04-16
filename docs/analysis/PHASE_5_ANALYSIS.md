# Phase 5 분석: L2 Genesis 생성

## 목차

1. [개요](#개요)
2. [용어 정의](#용어-정의)
3. [상세 분석](#상세-분석)
4. [호출 시퀀스](#호출-시퀀스)
5. [데이터 구조](#데이터-구조)
6. [에러 시나리오](#에러-시나리오)
7. [알려진 함정 및 개선 포인트](#알려진-함정-및-개선-포인트)

---

## 개요

Phase 5는 Phase 4에서 배포된 L1 컨트랙트 주소와 설정을 기반으로 L2 체인의 제네시스 블록(genesis.json)과 Rollup 설정(rollup.json)을 생성하는 과정을 다룬다. 이 단계는:

- **역할**: L2 초기 상태(initial state) 결정 및 고정
- **입력**: 
  - L1 배포 결과 (deploy-output.json)
  - 배포 설정 (DeployConfig)
  - Cannon prestate (Fault Proof 활성화 시)
- **출력**:
  - `genesis.json` - L2 제네시스 블록 파일
  - `rollup.json` - L2 Rollup 프로토콜 설정
  - 최대 3개의 추가 주입 단계 (DRB, USDC, MultiTokenPaymaster)
- **특징**:
  - 40개 이상의 Predeploy 컨트랙트 초기화
  - Proxy + Implementation 패턴을 통한 업그레이드 가능성
  - Hard Fork 적용 (Ecotone, Fjord, Granite)
  - Fault Proof System 선택적 지원 (Cannon prestate)
  - Account Abstraction 지원 (EntryPoint, Paymaster)

---

## 용어 정의

### Genesis 관련

- **L2 Genesis**: L2 체인의 첫 번째 블록(블록 0)의 상태. 계정 잔액, 배포된 컨트랙트, 저장소 값 포함
- **Genesis Block**: 블록 번호 0, 부모 해시 0x00...00, timestamp는 L1 시작 블록의 timestamp
- **genesis.json**: Geth/EVM 체인의 초기 상태를 JSON 형식으로 인코딩한 파일
- **Genesis Allocs**: genesis.json의 `alloc` 필드. 초기 계정/컨트랙트 상태 맵

### Predeploy 관련

- **Predeploy**: L2에서 사전 배포된 컨트랙트. 정해진 주소 공간(`0x4200...`)에 배포됨
- **Proxy**: UUPS/Transparent 프록시 컨트랙트. 실제 로직은 Implementation에 위임
- **Implementation**: 프록시가 위임하는 실제 로직 컨트랙트
- **Proxy Admin**: 프록시의 업그레이드 관리자. ProxyAdmin 컨트랙트가 담당
- **Immutables**: 프록시의 생성자 인자로 전달되는 불변(immutable) 데이터. 예: RemoteToken, Bridge, Decimals
- **Storage Config**: 프록시 또는 일부 특정 컨트랙트의 초기 저장소 값 세팅

### Predeploy 이름

- **L2ToL1MessagePasser** (0x4200000000000000000000000000000000000016): L2→L1 메시지 저장소
- **L2CrossDomainMessenger** (0x4200000000000000000000000000000000000007): L2 크로스 도메인 메신저
- **L2StandardBridge** (0x4200000000000000000000000000000000000010): L2 표준 브리지
- **L1Block** (0x4200000000000000000000000000000000000015): L1 블록 정보 오라클
- **LegacyERC20NativeToken** (0x4200000000000000000000000000000000000019): 네이티브 토큰 ERC20 래퍼
- **SequencerFeeVault**, **BaseFeeVault**, **L1FeeVault**: 수수료 수집 컨트랙트
- **Create2Deployer** (0x13b0d8EBB7e4c9b0a8A78e4d09cA000b0A0A0A0): CREATE2 배포자
- **Permit2** (0x000000000022d473030f116ddee9f6b43ac78ba3): ERC20 허가 관리
- **MultiCall3**: 배치 호출 유틸리티
- **OptimismMintableERC20Factory**: ERC20 생성 팩토리
- **OptimismMintableERC721Factory**: ERC721 생성 팩토리
- **ProxyAdmin** (0x4200000000000000000000000000000000000018): 프록시 관리자

### RollupConfig 관련

- **RollupConfig**: L2 Rollup 프로토콜의 설정 구조. ChainID, L1 정보, 수수료 파라미터, Hard Fork 시간 포함
- **SystemConfig**: L1 스마트 컨트랙트. Batch overhead, scalar, gas limit 등을 저장
- **Batcher Address**: 배치를 제출하는 계정 주소 (SystemConfig에 저장)
- **Overhead**: L1 배치 제출 시 기본 오버헤드 (Gwei 단위)
- **Scalar**: L1 데이터 비용 스칼라 (Gwei 당 바이트 비용)
- **Gas Limit**: L1 블록당 가스 한도
- **Hard Fork Times**: Ecotone, Fjord, Granite 활성화 시간 (Unix timestamp)

### Fault Proof 관련

- **Cannon**: OP Stack Fault Proof System의 증명 시스템. EVM을 MIPS로 컴파일하여 증명
- **Prestate**: Cannon 증명의 초기 상태. 거대한 바이너리 파일
- **Prestates Bucket**: S3 버킷에서 prestate를 다운로드하는 위치
- **Fault Proof enabled**: RollupConfig의 FaultGameFactoryAddress가 설정되었을 때

### 주입(Injection) 관련

- **DRB Injection**: Tokamak의 게이밍 preset에만 해당. DeploymentBridge와 관련 컨트랙트 추가 배포
- **USDC Injection**: 특정 preset에서 Circle의 USDC 토큰 설정 추가
- **MultiTokenPaymaster Injection**: Account Abstraction 지원을 위한 Paymaster 배포
- **Injection 단계**: genesis.json 생성 후 추가 트랜잭션/계정 주입 (3개 단계까지 가능)

### 배포 설정 관련

- **EnableFaultProof**: DeployConfig의 boolean 플래그. Fault Proof System 배포를 활성화하며, OptimismPortal2 및 Cannon prestate 설정을 트리거함
- **FaultGameFactoryAddress**: L1의 FaultGame 팩토리 컨트랙트 주소. OptimismPortal2가 dispute game 생성 시 사용
- **SystemConfigAddress**: L1 SystemConfig 컨트랙트 주소. RollupConfig immutables에서 참조되며 배치 overhead 및 scalar 저장
- **BlockID**: L1 블록 번호와 블록 해시의 조합으로 이루어진 구조. Fault Proof system에서 challenge commitment로 사용됨
- **DepositContractAddress**: L1 deposit 컨트랙트 주소. Immutables 설정에서 브릿지 초기화에 사용됨

---

## 상세 분석

### 5.1 SDK 엔트리 포인트: DeployContracts() 함수

**파일**: `trh-sdk/pkg/stacks/thanos/deploy_contracts.go`

**함수**: `DeployContracts()` (라인 23-407)

```go
func DeployContracts(ctx context.Context, deployContractsConfig *DeployContractsInput) error
```

**역할**:
- L1 컨트랙트 배포 및 L2 Genesis 생성 통합 오케스트레이션
- Cannon prestate 빌드 (Fault Proof 활성화 시)
- 외부 바이너리(tokamak-deployer) 호출 및 출력 처리
- DRB, USDC, MultiTokenPaymaster 주입 조율

**입력 구조** (DeployContractsInput):
```go
type DeployContractsInput struct {
    // L1 RPC 및 인증
    L1RpcUrl       string            // "https://mainnet.infura.io/v3/..."
    L1BeaconUrl    string            // Beacon API URL
    
    // 배포자 정보
    AdminPrivateKey      string      // 0x...
    SequencerPrivateKey  string      // 0x...
    BatcherPrivateKey    string      // 0x...
    ProposerPrivateKey   string      // 0x...
    ChallengerPrivateKey string      // Fault Proof 활성화 시만
    
    // 배포 설정
    DeployConfigPath    string      // /path/to/deploy-config.json
    DeployOutputPath    string      // /tmp/deploy-output
    
    // 선택 사항
    ResumeDeployment    bool        // 기존 배포 재개
    EnableFaultProof    bool        // Cannon 활성화
    Cannon              *CannonInfo // Cannon prestate 정보
    
    // 위치
    CannonPresateBucket string      // S3 bucket
    CannonPresatePath   string      // s3://...
}
```

**핵심 로직**:

1. **검증 및 초기화** (라인 41-80):
   - DeployConfigPath 존재 확인
   - L1 RPC 연결 검증 (chainID 확인)
   - 배포자 주소 추출 (private key → address 변환)
   - 출력 디렉토리 생성 (DeployOutputPath)

2. **Cannon Prestate 빌드** (라인 221-259):
   ```go
   if deployContractsConfig.EnableFaultProof {
       // prestate.bin 다운로드 (CannonPresateBucket에서)
       // 또는 로컬에서 빌드 (Cannon 설치 필요)
       prestatePath := filepath.Join(
           deployContractsConfig.DeployOutputPath,
           "cannon-prestate.bin",
       )
       // buildCannonPrestate() 호출
   }
   ```

3. **Deploy Config JSON 생성** (라인 260-280):
   - L1 배포 설정을 JSON 형식으로 변환
   - `deploy-config.json`을 임시 파일로 생성
   - 파일 경로: `DeployOutputPath/deploy-config.json`

4. **L1 배포 바이너리 호출** (라인 353-361):
   ```go
   cmd := exec.CommandContext(ctx, "tokamak-deployer", "deploy-contracts",
       "--deployer-address", deployer,
       "--l1-rpc-url", deployContractsConfig.L1RpcUrl,
       "--l1-beacon-url", deployContractsConfig.L1BeaconUrl,
       "--config-file", configPath,
       "--output-dir", deployContractsConfig.DeployOutputPath,
   )
   ```
   - **중요**: 이 단계에서 deploy-output.json 생성
   - 배포 상태 로그 출력 (`fmt.Println(string(output))`)

5. **L2 Genesis 생성 바이너리 호출** (라인 364-380):
   ```go
   genesisOpts := &GenerateGenesisOptions{
       DeployOutputPath: deployContractsConfig.DeployOutputPath,
       ConfigPath:       configPath,
       OutPath:          outGenesisPath,
       Prestate:         prestatePath, // Fault Proof일 때만
   }
   
   // tokamak-deployer generate-genesis 호출
   cmd = exec.CommandContext(ctx, "tokamak-deployer", "generate-genesis",
       "--deploy-output", genesisOpts.DeployOutputPath,
       "--config", genesisOpts.ConfigPath,
       "--output", genesisOpts.OutPath,
   )
   if genesisOpts.Prestate != "" {
       cmd.Args = append(cmd.Args, "--prestate", genesisOpts.Prestate)
   }
   ```
   - **CRITICAL**: 이 호출이 genesis.json과 rollup.json 생성
   - 출력 파일: `{OutPath}/genesis.json`, `{OutPath}/rollup.json`

6. **주입 단계 순차 실행** (라인 381-407):
   ```go
   // DRB 주입 (preset에 따라)
   if needsDRBInjection(deployContractsConfig.Preset) {
       runInjectDRB(...)
   }
   
   // USDC 주입 (preset에 따라)
   if needsUSDCInjection(deployContractsConfig.Preset) {
       runInjectUSDS(...)
   }
   
   // MultiTokenPaymaster 주입 (AA 활성화 시)
   if deployContractsConfig.EnableAccountAbstraction {
       runInjectPaymaster(...)
   }
   ```

**출력**:
- genesis.json (최대 1MB 이상, 40K+ 라인)
- rollup.json (수십 라인의 설정 JSON)
- deploy-output.json (참조용, Phase 4 출력)

---

### 5.2 Genesis 생성 라이브러리: BuildL2Genesis() 함수

**파일**: `op-chain-ops/genesis/layer_two.go`

**함수**: `BuildL2Genesis()` (라인 39-204)

```go
func BuildL2Genesis(config *DeployConfig, l1StartBlock *types.Block) (*core.Genesis, error)
```

**역할**:
- Geth `core.Genesis` 구조 생성
- 모든 Predeploy 컨트랙트 초기화
- L2 초기 계정 잔액 설정
- Hard Fork 업그레이드 트랜잭션 적용

**입력**:
- `config` (*DeployConfig): Phase 4에서 변환된 배포 설정
- `l1StartBlock` (*types.Block): L1 시작 블록 정보

**핵심 로직**:

1. **Genesis 틀 생성** (라인 46-60):
   ```go
   genesis := NewL2Genesis(config, l1StartBlock)
   // 다음 필드 설정:
   // - ChainID: config.L2ChainId
   // - Timestamp: l1StartBlock.Time()
   // - GasLimit: 30,000,000 (또는 설정값)
   // - Difficulty: 0
   // - Nonce: 0
   // - Config: ChainConfig (홈스테드, 탄젠트, 뮤어글래시어 등)
   ```

2. **개발 계정 펀딩** (라인 61-80) (선택사항):
   ```go
   if config.FundDevAccounts {
       for _, devAcct := range dev.Accounts {
           // 1000 ETH 할당 (~1e21 wei)
           genesis.Allocs[devAcct.Address] = 
               types.Account{ Balance: new(big.Int).Mul(...) }
       }
   }
   ```

3. **상태 DB 초기화** (라인 81-100):
   ```go
   db := state.NewMemoryStateDB(genesis)
   // in-memory 상태 데이터베이스
   // 모든 predeploy 초기화는 이 db에 수행됨
   ```

4. **L2 Immutables 및 Storage 설정** (라인 101-140):
   ```go
   // Immutables: predeploy 프록시의 생성자 인자
   l2Immutables, err := NewL2ImmutableConfig(config, l1StartBlock)
   
   // Storage: 특정 컨트랙트의 초기 저장소 값
   storageConfig, err := NewL2StorageConfig(config, l1StartBlock)
   
   // Example of Immutables:
   // LegacyERC20NativeToken: {
   //     RemoteToken: config.L1StandardBridgeProxy,
   //     Bridge: config.L2StandardBridgeProxy,
   //     Decimals: 18,
   // }
   ```

5. **Proxy 설정** (라인 141-160):
   ```go
   namespace := new(big.Int).SetUint64(0x4200000000000000)
   setProxies(db, config.ProxyAdminAddr, namespace, 0xff)
   // 0x4200000000000000000000000000000000000000 ~ 0x42000000000000000000000000000000000000ff
   // 총 256개의 Proxy 슬롯 생성
   ```

6. **Predeploy 배포** (라인 161-185):
   ```go
   deployments := immutables.Deploy(db, l2Immutables)
   // 각 Predeploy의 구현체(Implementation) 바이트코드 배포
   // 반환: DeploymentResults{
   //     addresses map[string]common.Address
   //     txs       []*types.Transaction
   // }
   ```

7. **개별 Predeploy 설정** (라인 186-195):
   ```go
   for _, name := range predeploy.Names() {
       setupPredeploy(db, deployments, storageConfig, name, ...)
       // Proxy 주소에 Implementation 주소 저장
       // Storage 값 초기화
   }
   ```

8. **Hard Fork 업그레이드 트랜잭션 적용** (라인 196-204):
   ```go
   err = PerformUpgradeTxs(db, l1StartBlock, config)
   // Ecotone hard fork: EIP-4844 blob 관련 상태 초기화
   // Fjord hard fork: 추가 업그레이드 트랜잭션
   // Granite hard fork: 추가 업그레이드 트랜잭션
   ```

**출력**:
- `*core.Genesis` 구조체 (in-memory, JSON 직렬화 가능)

---

### 5.3 Genesis 블록 기본 틀: NewL2Genesis() 함수

**파일**: `op-chain-ops/genesis/genesis.go`

**함수**: `NewL2Genesis()` (라인 24-106)

```go
func NewL2Genesis(config *DeployConfig, block *types.Block) (*core.Genesis, error)
```

**역할**:
- 기본 Genesis 블록 구조 생성
- Geth ChainConfig 설정 (HF 시간)
- Gas 한도, 초기 난이도 등 기본 파라미터

**입력**:
- `config` (*DeployConfig): 배포 설정
- `block` (*types.Block): L1 시작 블록

**핵심 로직**:

```go
// 라인 33-40에서:
// ChainConfig 생성
chainConfig := &params.ChainConfig{
    ChainID:             big.NewInt(int64(config.L2ChainId)),
    HomesteadBlock:      big.NewInt(0),
    DAOForkBlock:        nil,
    EIP150Block:         big.NewInt(0),
    EIP155Block:         big.NewInt(0),
    EIP158Block:         big.NewInt(0),
    ByzantiumBlock:      big.NewInt(0),
    ConstantinopleBlock: big.NewInt(0),
    PetersburgBlock:     big.NewInt(0),
    IstanbulBlock:       big.NewInt(0),
    MuirGlacierBlock:    big.NewInt(0),
    BerlinBlock:         big.NewInt(0),
    LondonBlock:         big.NewInt(0),
    ArrowGlacierBlock:   big.NewInt(0),
    GrayGlacierBlock:    big.NewInt(0),
}

// 라인 42-55에서:
// OP Stack 특정 fork 시간
if config.EcotoneTime != nil {
    chainConfig.Ecotone = config.EcotoneTime
}
if config.FjordTime != nil {
    chainConfig.Fjord = config.FjordTime
}
if config.GraniteTime != nil {
    chainConfig.Granite = config.GraniteTime
}

// 라인 58-80에서:
// Genesis 객체 생성
genesis := &core.Genesis{
    Config:      chainConfig,
    Nonce:       0,
    Timestamp:   block.Time(), // L1 시작 블록 시간
    ExtraData:   []byte{},
    GasLimit:    config.L2GasLimitPerBlock, // 30,000,000 기본값
    Difficulty:  big.NewInt(0),
    Mixhash:     common.Hash{},
    Coinbase:    common.Address{}, // 또는 설정값
    Alloc:       make(core.GenesisAlloc),
}

// 라인 82-95에서:
// 선택 사항: Deposit Contract 정보
if config.L1Contracts.OptimismPortalProxy != (common.Address{}) {
    genesis.Alloc[config.L1Contracts.OptimismPortalProxy] = 
        types.Account{
            // Deposit 관련 메타데이터
        }
}

return genesis, nil
```

**주요 파라미터**:
- `L2ChainId`: L2 체인 ID (예: 901 for Tokamak Thanos Sepolia)
- `L2GasLimitPerBlock`: 블록당 가스 한도 (기본 30,000,000)
- `EcotoneTime`, `FjordTime`, `GraniteTime`: Hard Fork 활성화 시간

---

### 5.4 Predeploy Immutables 설정: NewL2ImmutableConfig() 함수

**파일**: `op-chain-ops/genesis/config.go`

**함수**: `NewL2ImmutableConfig()` (라인 1065-1260)

```go
func NewL2ImmutableConfig(config *DeployConfig, block *types.Block) 
    (*immutables.PredeploysImmutableConfig, error)
```

**역할**:
- 40개 이상 Predeploy의 생성자 인자(Immutables) 설정
- 토큰 브릿지 정보, Vault 수신자, 원격 토큰 주소 포함

**반환 구조** (immutables.PredeploysImmutableConfig):

```go
type PredeploysImmutableConfig struct {
    // 1. L2ToL1MessagePasser (0x4200000000000000000000000000000000000016)
    // - 불변 없음, 순수 상태 저장소
    
    // 2. L2CrossDomainMessenger (0x4200000000000000000000000000000000000007)
    // - OtherMessenger: config.L1CrossDomainMessengerProxy
    
    // 3. L2StandardBridge (0x4200000000000000000000000000000000000010)
    // - Messenger: L2CrossDomainMessenger
    // - OtherBridge: config.L1StandardBridgeProxy
    
    // 4. LegacyERC20NativeToken (0x4200000000000000000000000000000000000019)
    LegacyERC20NativeToken: &immutables.LegacyERC20NativeTokenConfig{
        RemoteToken: config.L1NativeTokenProxy,  // L1 토큰 주소
        Bridge:      config.L2StandardBridgeProxy, // L2 브릿지
        Decimals:    18,                          // 토큰 소수점
    },
    
    // 5-7. Fee Vaults (SequencerFeeVault, BaseFeeVault, L1FeeVault)
    SequencerFeeVault: &immutables.FeeVaultConfig{
        Recipient:  config.SequencerFeeVaultRecipient,
        MinBalance: big.NewInt(10 * 1e18), // 10 ETH 임계값
        WithdrawAmount: big.NewInt(0), // 또는 설정값
    },
    
    // 8. Create2Deployer
    Create2Deployer: &immutables.Create2DeployerConfig{
        // 생성자 인자 없음
    },
    
    // 9. Permit2
    Permit2: &immutables.Permit2Config{
        // 미리 배포된 주소
    },
    
    // 10. MultiCall3
    MultiCall3: &immutables.MultiCall3Config{
        // 생성자 인자 없음
    },
    
    // 11. OptimismMintableERC20Factory
    OptimismMintableERC20Factory: &immutables.ERC20FactoryConfig{
        Bridge: config.L2StandardBridgeProxy,
    },
    
    // 12. OptimismMintableERC721Factory
    OptimismMintableERC721Factory: &immutables.ERC721FactoryConfig{
        Bridge: config.L2ERC721BridgeProxy,
    },
    
    // ... (20개 이상 추가 Predeploy)
    
    // EntryPoint (Account Abstraction)
    EntryPoint: &immutables.EntryPointConfig{
        // 생성자 인자 없음
    },
    
    // Paymaster들
    VerifyingPaymasterPredeploy: &immutables.PaymasterConfig{
        Signer: config.PaymasterSignerAddress,
    },
}
```

**핵심 로직** (라인 1080-1150):

```go
// L1 컨트랙트 주소 검증
if config.L1Contracts.OptimismPortalProxy == (common.Address{}) {
    return nil, fmt.Errorf("OptimismPortalProxy not set")
}

// 네이티브 토큰 설정
immutableConfig.LegacyERC20NativeToken = &immutables.LegacyERC20NativeTokenConfig{
    RemoteToken: config.L1StandardBridgeProxy,
    Bridge:      config.L2StandardBridgeProxy,
    Decimals:    18,
}

// Vault 수신자 설정 (기본값: 첫 번째 deployer)
if config.SequencerFeeVaultRecipient == (common.Address{}) {
    config.SequencerFeeVaultRecipient = 
        config.Deployers[0] // AdminPrivateKey에서 파생
}

// Fee Vault 설정
immutableConfig.SequencerFeeVault = &immutables.FeeVaultConfig{
    Recipient: config.SequencerFeeVaultRecipient,
    MinBalance: config.SequencerFeeVaultMinBalance,
    WithdrawAmount: config.SequencerFeeVaultWithdrawAmount,
}
immutableConfig.BaseFeeVault = &immutables.FeeVaultConfig{
    Recipient: config.BaseFeeVaultRecipient,
    MinBalance: config.BaseFeeVaultMinBalance,
    WithdrawAmount: config.BaseFeeVaultWithdrawAmount,
}

// Create2Deployer는 사전 설정된 주소에서 로드
immutableConfig.Create2Deployer = 
    &immutables.Create2DeployerConfig{} // 주소: 0x13b0d8...

// ... (모든 Predeploy 반복)
```

**검증 및 반환** (라인 1230-1260):

```go
// 모든 필수 필드 검증
if immutableConfig.SequencerFeeVault == nil {
    return nil, errors.New("SequencerFeeVault not set")
}

return immutableConfig, nil
```

---

### 5.5 Predeploy Storage 설정: NewL2StorageConfig() 함수

**파일**: `op-chain-ops/genesis/config.go`

**함수**: `NewL2StorageConfig()` (라인 1257-1450)

```go
func NewL2StorageConfig(config *DeployConfig, block *types.Block) 
    (state.StorageConfig, error)
```

**역할**:
- Predeploy 컨트랙트의 초기 저장소 값 설정
- L1 블록 정보, 메신저 주소, 토큰 정보 등 저장

**반환 구조** (state.StorageConfig):

```go
type StorageConfig struct {
    // 주소 → (슬롯 → 값) 맵
    [common.Address]map[string]interface{}
}
```

**L2ToL1MessagePasser** (라인 1280-1320):

```go
// 초기 상태: msgNonce = 0
storageConfig[L2ToL1MessagePasserAddress] = map[string]interface{}{
    "msgNonce": big.NewInt(0), // 메시지 nonce 초기화
}
```

**L2CrossDomainMessenger** (라인 1330-1360):

```go
storageConfig[L2CrossDomainMessengerAddress] = map[string]interface{}{
    "_initialized": uint8(1),           // 초기화 완료
    "_initializing": uint8(0),          // 초기화 중 아님
    "xDomainMsgSender": common.Address{}, // 초기 값
    "otherMessenger": config.L1CrossDomainMessengerProxy,
}
```

**L2StandardBridge** (라인 1370-1390):

```go
storageConfig[L2StandardBridgeAddress] = map[string]interface{}{
    "_initialized": uint8(1),
    "initialized": true,
    "messenger": L2CrossDomainMessengerAddress,
    "otherBridge": config.L1StandardBridgeProxy,
}
```

**L1Block (Oracle)** (라인 1400-1430):

```go
// L1Block은 특수 로직: sequencer 업데이트를 통해 매 블록마다 업데이트
storageConfig[L1BlockAddress] = map[string]interface{}{
    "number": block.Number(),           // L1 시작 블록 번호
    "timestamp": block.Time(),          // L1 시작 블록 타임스탬프
    "basefee": block.BaseFee(),        // L1 기본 수수료
    "hash": block.Hash(),              // L1 해시
    "sequenceNumber": uint64(0),       // L2의 첫 번째 블록
    "batcherHash": config.BatcherHash, // Batcher 주소 해시
    "l1FeeOverhead": config.GasPriceOracle.Overhead,
    "l1FeeScalar": config.GasPriceOracle.Scalar,
}
```

**LegacyERC20NativeToken** (라인 1440-1450, 조건부):

```go
if config.IsNativeTokenERC20 {
    storageConfig[LegacyERC20NativeTokenAddress] = map[string]interface{}{
        "name": "Wrapped Ether",
        "symbol": "WETH",
        "decimals": uint8(18),
    }
}
```

---

### 5.6 Proxy 배포: setProxies() 함수

**파일**: `op-chain-ops/genesis/setters.go`

**함수**: `setProxies()` (라인 36-95)

```go
func setProxies(db vm.StateDB, proxyAdminAddr common.Address, 
    namespace *big.Int, count uint64) error
```

**역할**:
- 0x4200... 네임스페이스에 256개의 Proxy 컨트랙트 생성
- 각 Proxy의 admin 슬롯 설정
- Proxy 바이트코드 배포

**핵심 로직**:

```go
// 라인 47-60에서:
// Proxy 바이트코드 로드
proxyBytecode, err := bindings.GetProxyBytecode()
if err != nil {
    return fmt.Errorf("failed to load proxy bytecode: %w", err)
}

// 라인 65-85에서:
// 256개 Proxy 생성
for i := uint64(0); i < count; i++ {
    // 주소 계산: 0x4200000000000000000000000000000000000000 + i
    addr := common.Address{}
    addr.SetBytes(
        new(big.Int).Add(namespace, big.NewInt(int64(i))).Bytes(),
    )
    
    // Proxy 계정 생성
    db.CreateAccount(addr)
    
    // Proxy 바이트코드 배포
    db.SetCode(addr, proxyBytecode)
    
    // Admin 슬롯 설정 (ERC-1967: 0x...09)
    adminSlot := common.Hash{}
    adminSlot[31] = 0x09 // 관리자 슬롯 위치
    db.SetState(addr, adminSlot, proxyAdminAddr.Hash())
}

// 라인 88-92에서:
// 특수 처리 (선택사항)
// 일부 Proxy는 특수 논리 필요 (예: L2UsdcBridge는 USDC 토큰 접근)
if i == USDC_BRIDGE_INDEX {
    // 특수 바이트코드 로드
}

return nil
```

**주요 Proxy 주소**:
- 0x4200000000000000000000000000000000000000: Proxy 0 (예비)
- 0x4200000000000000000000000000000000000007: L2CrossDomainMessenger Proxy
- 0x4200000000000000000000000000000000000010: L2StandardBridge Proxy
- 0x4200000000000000000000000000000000000015: L1Block Proxy
- 0x4200000000000000000000000000000000000016: L2ToL1MessagePasser Proxy
- 0x4200000000000000000000000000000000000018: ProxyAdmin
- 0x4200000000000000000000000000000000000019: LegacyERC20NativeToken Proxy
- ... (256개까지 계속)

---

### 5.7 개별 Predeploy 설정: setupPredeploy() 함수

**파일**: `op-chain-ops/genesis/setters.go`

**함수**: `setupPredeploy()` (라인 111-135)

```go
func setupPredeploy(db vm.StateDB, deployResults immutables.DeploymentResults,
    storage state.StorageConfig, name string, 
    proxyAddr common.Address, implAddr common.Address) error
```

**역할**:
- 각 Predeploy의 implementation 바이트코드 배포
- Proxy의 저장소에 implementation 주소 저장
- 초기 저장소 값 설정

**핵심 로직**:

```go
// 라인 120-125에서:
// Proxy 저장소에 implementation 주소 저장
// ERC-1967 Proxy 표준: 0x...05 슬롯 = implementation 주소
implSlot := common.Hash{}
implSlot[31] = 0x05
db.SetState(proxyAddr, implSlot, implAddr.Hash())

// 라인 126-130에서:
// Implementation 바이트코드 배포
if bytes, ok := deployResults[name]; ok {
    // immutables.Deploy에서 반환된 바이트코드
    db.SetCode(implAddr, bytes)
} else {
    // Solc 컴파일된 아티팩트에서 로드
    implBytecode, err := bindings.GetBytecode(name)
    if err != nil {
        return fmt.Errorf("failed to load bytecode for %s: %w", name, err)
    }
    db.SetCode(implAddr, implBytecode)
}

// 라인 131-135에서:
// 초기 저장소 값 설정
if storageVals, ok := storage[proxyAddr]; ok {
    // StorageConfig에서 제공된 값들 적용
    for slot, value := range storageVals {
        db.SetState(proxyAddr, slot, value)
    }
}

return nil
```

**특이사항**:
- Proxy는 ERC-1967 표준을 따름 (admin slot: 0x...09, impl slot: 0x...05)
- 저장소 값은 Proxy 주소가 아닌 **proxyAddr**에 저장됨 (delegatecall이므로)

---

### 5.8 Predeploy 구현체 배포: immutables.Deploy() 함수

**파일**: `op-chain-ops/immutables/immutables.go`

**함수**: `Deploy()` (라인 209-270)

```go
func Deploy(db vm.StateDB, config *PredeploysImmutableConfig) 
    (DeploymentResults, error)
```

**역할**:
- 각 Predeploy의 구현체를 생성자 인자(immutables) 포함하여 배포
- 각 Predeploy는 사전 정의된 주소(0x4200...0x0001, 0x4200...0x0002, ...)에 배포됨
- 바이트코드 + immutables 포함 바이너리 반환

**핵심 로직**:

```go
// 라인 40-50에서:
// 구현체 주소 계산
// 각 Predeploy는 순차 번호로 배포
// 예: 
//   L2ToL1MessagePasser impl: 0x4200000000000000000000000000000000010001
//   L2CrossDomainMessenger impl: 0x4200000000000000000000000000000000020001
//   L2StandardBridge impl: 0x4200000000000000000000000000000000030001

// 라인 60-100에서:
// 각 Predeploy 배포
for _, predeploy := range Predeploys {
    // Solc 바이트코드 로드
    bytecode, err := bindings.GetBytecode(predeploy.Name)
    if err != nil {
        return nil, err
    }
    
    // 생성자 인자(immutables) 포함
    // Immutables: 바이트코드의 특정 위치에 값을 삽입하는 방식
    // 예: 생성자가 remoteToken, bridge, decimals를 받으면,
    //     바이트코드에 이 값들을 런타임 상수로 컴파일
    
    immutablesForPredeploy := getImmutablesForPredeploy(
        predeploy.Name, config)
    
    finalBytecode := injectImmutables(bytecode, immutablesForPredeploy)
    
    // 상태에 배포
    implAddr := predeploy.ImplAddress()
    db.CreateAccount(implAddr)
    db.SetCode(implAddr, finalBytecode)
    
    results[predeploy.Name] = finalBytecode
}

return results, nil
```

**Immutables 주입 예시**:

```
// 원본 Solc 바이트코드:
60 XX 60 XX 60 XX ... (생성자 인자들이 placeholder로 컴파일됨)

// Immutables 주입 후:
60 12 34 56 60 AB CD EF 60 12 34 56 ...
   ↑ RemoteToken  ↑ Bridge       ↑ Decimals
```

---

### 5.9 Hard Fork 업그레이드: PerformUpgradeTxs() 함수

**파일**: `op-chain-ops/genesis/upgrades.go`

**함수**: `PerformUpgradeTxs()` (라인 25-150)

```go
func PerformUpgradeTxs(db vm.StateDB, l1StartBlock *types.Block, 
    config *DeployConfig) error
```

**역할**:
- Ecotone, Fjord, Granite Hard Fork 관련 상태 변경 적용
- EIP-4844 blob 관련 상태 초기화 (Ecotone)
- 추가 업그레이드 트랜잭션 실행

**핵심 로직**:

```go
// 라인 40-60에서:
// Ecotone Hard Fork (EIP-4844 활성화)
if config.EcotoneTime != nil && config.EcotoneTime <= block.Time() {
    // 1. Blob 기본 요금 초기화 (데이터 가스 가격)
    db.SetState(
        L1BlockAddress,
        common.Hash{31: 0x0A}, // blobBaseFee 슬롯
        new(big.Int).SetUint64(1).Bytes(),
    )
    
    // 2. 데이터 가용성 상태 초기화
    db.SetState(
        L1BlockAddress,
        common.Hash{31: 0x0B}, // blobBaseFeeScalar 슬롯
        new(big.Int).SetUint64(0).Bytes(),
    )
}

// 라인 70-100에서:
// Fjord Hard Fork
if config.FjordTime != nil && config.FjordTime <= block.Time() {
    // 1. L1Block에 fjordTime 저장
    db.SetState(
        L1BlockAddress,
        common.Hash{31: 0x0C}, // fjordTime 슬롯
        new(big.Int).SetUint64(uint64(*config.FjordTime)).Bytes(),
    )
    
    // 2. 추가 계정 상태 변경 (필요 시)
    // 예: 특정 주소의 nonce 증가, 잔액 조정 등
}

// 라인 110-140: Granite Hard Fork
if config.GraniteTime != nil && config.GraniteTime <= block.Time() {
    // Granite 특정 변경사항
    // (아직 정의되지 않거나 미래 사항)
}

return nil
```

---

### 5.10 RollupConfig 생성: RollupConfig() 함수

**파일**: `op-chain-ops/genesis/config.go`

**함수**: `RollupConfig()` (라인 825-888)

```go
func (dc *DeployConfig) RollupConfig(
    l1StartBlock *types.Block,
    l2GenesisBlockHash common.Hash,
    l2GenesisBlockNumber uint64) (*rollup.Config, error)
```

**역할**:
- DeployConfig를 rollup.Config로 변환
- L2 네트워크의 프로토콜 파라미터 확정

**반환 구조** (rollup.Config):

```go
type Config struct {
    // Genesis
    Genesis struct {
        L1: BlockID                    // L1 시작 블록
        L2: BlockID                    // L2 genesis 블록
        SystemConfig: SystemConfig     // 시스템 설정
    }
    
    // 네트워크 파라미터
    L1ChainID:            uint64      // 1 (Mainnet), 11155111 (Sepolia)
    L2ChainID:            uint64      // 901 (Tokamak Thanos)
    BlockTime:            uint64      // 2초
    SeqWindowSize:        uint64      // 3600
    MaxSequencerDrift:    uint64      // 600초
    ChannelTimeout:       uint64      // 300 (블록)
    
    // 배치 제출
    BatchInboxAddress:    common.Address
    DepositContractAddress: common.Address
    
    // L1 시스템 설정
    L1SystemConfigAddress: common.Address
    SystemConfigParams: struct {
        Overhead:        uint64      // 기본 오버헤드
        Scalar:          uint64      // 데이터 스칼라
        BatcherAddr:     common.Address
        GasLimit:        uint64
    }
    
    // Fault Proof (선택사항)
    FaultGameFactoryAddress: common.Address
    FaultProofAbsolutePath:  string
    FaultProofMaxDepth:      uint64
    FaultProofClockExtension: uint64
    PreimageOracleAddress:   common.Address
    
    // Hard Fork Times (Unix timestamp)
    EcotoneTime: *uint64
    FjordTime:   *uint64
    GraniteTime: *uint64
}
```

**핵심 로직** (라인 850-880):

```go
// 라인 852-860: 입력 검증
if dc.L1Contracts.OptimismPortalProxy == (common.Address{}) {
    return nil, fmt.Errorf("OptimismPortalProxy not set")
}
if dc.L1Contracts.SystemConfigProxy == (common.Address{}) {
    return nil, fmt.Errorf("SystemConfigProxy not set")
}

// 라인 862-870: 기본값 설정
blockTime := uint64(2) // L2 블록 시간 (초)
if dc.L2BlockTime != 0 {
    blockTime = dc.L2BlockTime
}

seqWindowSize := uint64(3600)
if dc.SeqWindowSize != 0 {
    seqWindowSize = dc.SeqWindowSize
}

// 라인 872-888: RollupConfig 구성
rollupCfg := &rollup.Config{
    Genesis: rollup.Genesis{
        L1: rollup.BlockID{
            Hash:   l1StartBlock.Hash(),
            Number: l1StartBlock.NumberU64(),
        },
        L2: rollup.BlockID{
            Hash:   l2GenesisBlockHash,
            Number: l2GenesisBlockNumber,
        },
        SystemConfig: rollup.SystemConfig{
            BatcherAddr:  dc.BatcherAddress,
            Overhead:     dc.GasPriceOracle.Overhead,
            Scalar:       dc.GasPriceOracle.Scalar,
            GasLimit:     dc.L2GasLimitPerBlock,
        },
    },
    L1ChainID:   dc.L1ChainId,
    L2ChainID:   dc.L2ChainId,
    BlockTime:   blockTime,
    SeqWindowSize: seqWindowSize,
    MaxSequencerDrift: dc.MaxSequencerDrift,
    ChannelTimeout: dc.ChannelTimeout,
    
    BatchInboxAddress: dc.BatchInboxAddress,
    DepositContractAddress: dc.L1Contracts.OptimismPortalProxy,
    L1SystemConfigAddress: dc.L1Contracts.SystemConfigProxy,
    
    EcotoneTime: dc.EcotoneTime,
    FjordTime:   dc.FjordTime,
    GraniteTime: dc.GraniteTime,
}

// Fault Proof 설정 (선택사항)
if dc.FaultProofEnabled {
    rollupCfg.FaultGameFactoryAddress = dc.L1Contracts.FaultGameFactory
    rollupCfg.FaultProofAbsolutePath = dc.FaultProofAbsolutePath
    // ...
}

return rollupCfg, nil
```

---

### 5.11 JSON 직렬화 및 파일 생성

**파일**: tokamak-deployer 바이너리 (Rust)

**역할**:
- `BuildL2Genesis()` 반환값(core.Genesis)을 JSON으로 직렬화
- genesis.json 파일에 기록
- rollup.json 파일에 기록

**처리 흐름**:

```
core.Genesis (in-memory)
    ↓
json.MarshalIndent(..., "", "  ")
    ↓
[]byte (JSON 문자열)
    ↓
ioutil.WriteFile("genesis.json")
    ↓
File: genesis.json (약 1-2MB)

---

rollup.Config (in-memory)
    ↓
json.MarshalIndent(...)
    ↓
[]byte (JSON 문자열)
    ↓
ioutil.WriteFile("rollup.json")
    ↓
File: rollup.json (약 1-2KB)
```

**genesis.json 구조** (상위 50개 라인 예시):

```json
{
  "config": {
    "chainId": 901,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "muirGlacierBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "arrowGlacierBlock": 0,
    "grayGlacierBlock": 0,
    "ecotone": 0,
    "fjord": 0,
    "granite": 0
  },
  "nonce": "0x0",
  "timestamp": "0x66xxxx",
  "extraData": "0x",
  "gasLimit": "0x1c9c380",
  "difficulty": "0x1",
  "mixHash": "0x...",
  "coinbase": "0x...",
  "alloc": {
    "0x4200000000000000000000000000000000000000": {
      "code": "0x363d3d373d3d3d363d...",
      "storage": {
        "0x...": "0x..."
      }
    },
    "0x4200000000000000000000000000000000010001": {
      "code": "0x...",
      "storage": {}
    },
    ...
  }
}
```

---

### 5.12 주입 단계 (DRB, USDC, Paymaster)

**파일**: `trh-sdk/pkg/stacks/thanos/deploy_contracts.go` (라인 381-407)

**역할**:
- genesis.json 생성 후 추가 계정/컨트랙트 주입
- Preset에 따라 선택적으로 실행

**DRB 주입** (Gaming Preset):

```bash
tokamak-deployer inject-drb \
  --genesis genesis.json \
  --output genesis.json \
  --drb-config drb-config.json
```

**처리**:
- DeploymentBridge 컨트랙트 배포
- DRB 관련 초기 저장소 값 설정
- genesis.json alloc에 계정 추가

**USDC 주입** (특정 Preset):

```bash
tokamak-deployer inject-usdc \
  --genesis genesis.json \
  --output genesis.json \
  --usdc-proxy 0x...
```

**처리**:
- USDC 토큰 계정 추가
- USDC 잔액 할당
- USDC 관련 메타데이터 저장

**MultiTokenPaymaster 주입** (AA 활성화):

```bash
tokamak-deployer inject-paymaster \
  --genesis genesis.json \
  --output genesis.json \
  --paymaster-addr 0x...
```

**처리**:
- Paymaster 컨트랙트 배포
- EntryPoint 통합
- AA 관련 초기 상태 설정

---

### 5.13 Backend 완료 처리

**파일**: `trh-backend/pkg/services/thanos/deployment.go`

**함수**: `executeDeployments()` (라인 580-600)

```go
// Phase 4 완료 후:
// 1. SDK 클라이언트가 genesis.json, rollup.json 생성 완료
// 2. 반환값: DeployL1ContractsOutput{
//     GenesisPath: "/tmp/deploy-output/genesis.json",
//     RollupPath:  "/tmp/deploy-output/rollup.json",
//     Success: true,
// }

// 3. Backend가 파일 복사 (스택 저장소로)
filesToCopy := []string{
    "genesis.json",
    "rollup.json",
    "deploy-output.json",
}
for _, file := range filesToCopy {
    src := filepath.Join(output.DeploymentOutputPath, file)
    dst := filepath.Join(stackDataPath, file)
    err = os.Rename(src, dst)
}

// 4. 배포 상태 업데이트
statusChan <- entities.DeploymentStatusWithID{
    DeploymentID: deployment.ID,
    Status:       entities.DeploymentRunStatusSuccess,
}

// 5. 스택 상태 업데이트
stack.Status = entities.StackStatusDeployed
stack.L2GenesisPath = filepath.Join(stackDataPath, "genesis.json")
stack.RollupConfigPath = filepath.Join(stackDataPath, "rollup.json")
s.db.SaveStack(stack)
```

**저장 경로**:
- `/data/stacks/{stackId}/genesis.json`
- `/data/stacks/{stackId}/rollup.json`
- `/data/stacks/{stackId}/deploy-output.json`

---

## 호출 시퀀스

### 전체 호출 체인

```
1. Backend: executeDeployments()
   ├─ filters: DeployL1ContractsStep
   └─ calls: thanos.DeployL1Contracts(sdkClient, ...)
           ├─ Phase 4: L1 컨트랙트 배포 (별도 분석)
           └─ Phase 5 시작: L2 Genesis 생성

2. SDK: DeployContracts()
   ├─ buildCannonPrestate() (FP 활성화 시)
   ├─ makeDeployContractConfigJsonFile() → deploy-config.json
   ├─ exec.Command("tokamak-deployer", "deploy-contracts") → deploy-output.json
   │
   └─ exec.Command("tokamak-deployer", "generate-genesis")
       ├─ reads: deploy-config.json, deploy-output.json
       ├─ calls: BuildL2Genesis()
       │   ├─ NewL2Genesis() → 기본 틀
       │   ├─ NewL2ImmutableConfig() → Predeploy immutables
       │   ├─ NewL2StorageConfig() → Predeploy 저장소
       │   ├─ setProxies() → 256개 Proxy 생성
       │   ├─ immutables.Deploy() → Implementation 배포
       │   ├─ setupPredeploy() (× N) → 각 Predeploy 설정
       │   └─ PerformUpgradeTxs() → Hard Fork 업그레이드
       │
       ├─ RollupConfig() → rollup.Config 생성
       ├─ json.Marshal(core.Genesis) → genesis.json
       ├─ json.Marshal(rollup.Config) → rollup.json
       │
       └─ (선택사항)
           ├─ inject-drb
           ├─ inject-usdc
           └─ inject-paymaster

3. Backend: 파일 저장
   ├─ genesis.json → 스택 저장소 복사
   ├─ rollup.json → 스택 저장소 복사
   └─ 배포 완료 표시
```

### 호출 시퀀스 다이어그램 (Mermaid)

```mermaid
sequenceDiagram
    participant BE as Backend<br/>(deployment.go)
    participant SDK as SDK<br/>(deploy_contracts.go)
    participant BIN as tokamak-deployer<br/>(binary)
    participant LIB as op-chain-ops<br/>(Go lib)
    participant FS as Filesystem<br/>(genesis.json)

    BE->>SDK: DeployContracts(deployContractsConfig)
    activate SDK
    
    SDK->>SDK: buildCannonPrestate() [FP only]
    SDK->>SDK: makeDeployContractConfigJsonFile()
    SDK->>BIN: exec("deploy-contracts", ...)
    activate BIN
    BIN->>BIN: Phase 4 (L1 contracts)
    BIN-->>SDK: deploy-output.json
    deactivate BIN
    
    SDK->>BIN: exec("generate-genesis", ...)
    activate BIN
    BIN->>LIB: BuildL2Genesis(config, l1Block)
    activate LIB
    LIB->>LIB: NewL2Genesis() → core.Genesis
    LIB->>LIB: NewL2ImmutableConfig()
    LIB->>LIB: NewL2StorageConfig()
    LIB->>LIB: setProxies()
    LIB->>LIB: immutables.Deploy()
    LIB->>LIB: setupPredeploy() × 40
    LIB->>LIB: PerformUpgradeTxs()
    LIB-->>BIN: core.Genesis (in-memory)
    deactivate LIB
    
    BIN->>BIN: RollupConfig(genesis)
    BIN->>BIN: json.Marshal(core.Genesis)
    BIN->>FS: Write genesis.json
    BIN->>BIN: json.Marshal(rollup.Config)
    BIN->>FS: Write rollup.json
    BIN-->>SDK: success
    deactivate BIN
    
    SDK->>SDK: inject-drb (if needed)
    SDK->>SDK: inject-usdc (if needed)
    SDK->>SDK: inject-paymaster (if needed)
    
    SDK-->>BE: DeployL1ContractsOutput {GenesisPath, RollupPath, Success}
    
    BE->>FS: Copy genesis.json → stack dir
    BE->>FS: Copy rollup.json → stack dir
    deactivate SDK
    
    BE->>BE: Update stack status
    deactivate BE
```

---

## 데이터 구조

### 6.1 DeployConfig 구조 (Phase 4 출력 → Phase 5 입력)

**파일**: `op-chain-ops/genesis/config.go`

```go
type DeployConfig struct {
    // L1 배포 결과
    L1Contracts struct {
        OptimismPortalProxy     common.Address
        OptimismPortal2Proxy    common.Address // FP enabled
        SystemConfigProxy       common.Address
        L2OutputOracleProxy     common.Address
        L2ToL1MessagePasserImpl  common.Address
        // ... 50+ contracts
    }
    
    // 운영자
    Deployers []common.Address // [AdminAddress, SequencerAddress, ...]
    BatcherAddress       common.Address
    ProposerAddress      common.Address
    ChallengerAddress    common.Address // FP only
    
    // L2 파라미터
    L2ChainId            uint64     // 901
    L2BlockTime          uint64     // 2
    L2GasLimitPerBlock   uint64     // 30,000,000
    
    // 배치 설정
    SeqWindowSize        uint64     // 3600
    MaxSequencerDrift    uint64     // 600
    ChannelTimeout       uint64     // 300
    BatchInboxAddress    common.Address
    
    // 수수료 설정
    GasPriceOracle struct {
        Overhead  uint64         // 188
        Scalar    uint64         // 684000
    }
    
    // Vault 수신자
    SequencerFeeVaultRecipient  common.Address
    BaseFeeVaultRecipient       common.Address
    L1FeeVaultRecipient         common.Address
    
    // Hard Fork 시간
    EcotoneTime  *uint64
    FjordTime    *uint64
    GraniteTime  *uint64
    
    // Fault Proof
    FaultProofEnabled      bool
    FaultProofAbsolutePath string
    
    // 기타
    EnableFundDevAccounts bool
    IsNativeTokenERC20    bool
}
```

### 6.2 core.Genesis 구조

**패키지**: `github.com/ethereum/go-ethereum/core`

```go
type Genesis struct {
    Config   *params.ChainConfig `json:"config"`
    Nonce    uint64              `json:"nonce"`
    Timestamp uint64             `json:"timestamp"`
    ExtraData []byte             `json:"extraData"`
    GasLimit  uint64             `json:"gasLimit"`
    Difficulty *big.Int          `json:"difficulty"`
    Mixhash   common.Hash        `json:"mixHash"`
    Coinbase  common.Address     `json:"coinbase"`
    Alloc     GenesisAlloc       `json:"alloc"`
    Number    uint64             `json:"number"`
    GasUsed   uint64             `json:"gasUsed"`
    ParentHash common.Hash       `json:"parentHash"`
}

type GenesisAlloc map[common.Address]Account

type Account struct {
    Code    []byte             `json:"code"`
    Storage map[common.Hash]common.Hash `json:"storage"`
    Balance *big.Int           `json:"balance"`
    Nonce   uint64             `json:"nonce"`
}
```

### 6.3 rollup.Config 구조

**파일**: `op-node/rollup/config.go`

```go
type Config struct {
    Genesis Genesis `json:"genesis"`
    
    L1ChainID   uint64         `json:"l1ChainID"`
    L2ChainID   uint64         `json:"l2ChainID"`
    BlockTime   uint64         `json:"blockTime"`
    SeqWindowSize uint64       `json:"seqWindowSize"`
    MaxSequencerDrift uint64   `json:"maxSequencerDrift"`
    ChannelTimeout uint64      `json:"channelTimeout"`
    
    BatchInboxAddress       common.Address `json:"batchInboxAddress"`
    DepositContractAddress  common.Address `json:"depositContractAddress"`
    L1SystemConfigAddress   common.Address `json:"l1SystemConfigAddress"`
    
    EcotoneTime *uint64       `json:"ecotoneTime"`
    FjordTime   *uint64       `json:"fjordTime"`
    GraniteTime *uint64       `json:"graniteTime"`
    
    FaultGameFactoryAddress common.Address `json:"faultGameFactoryAddress"`
    FaultProofAbsolutePath  string         `json:"faultProofAbsolutePath"`
    FaultProofMaxDepth      uint64         `json:"faultProofMaxDepth"`
}

type Genesis struct {
    L1 BlockID        `json:"l1"`
    L2 BlockID        `json:"l2"`
    SystemConfig SystemConfig `json:"systemConfig"`
}

type SystemConfig struct {
    BatcherAddr common.Address `json:"batcherAddr"`
    Overhead    uint64         `json:"overhead"`
    Scalar      uint64         `json:"scalar"`
    GasLimit    uint64         `json:"gasLimit"`
}

type BlockID struct {
    Hash   common.Hash `json:"hash"`
    Number uint64      `json:"number"`
}
```

### 6.4 Predeploy 목록 (40+ 계약)

**주소 표기**: 아래 테이블은 전체 42바이트 predeploy 주소를 사용합니다. (`0x420000000000000000000000000000000000XX` 형식)

| # | 주소 | 이름 | Proxy | Implementation | Immutables | Storage |
|---|------|------|-------|-----------------|-----------|---------|
| 1 | 0x4200000000000000000000000000000000000000 | - | - | - | - | - |
| 2 | 0x4200000000000000000000000000000000000001 | L1MessageSender (Legacy) | ✓ | ✓ | ✗ | ✗ |
| 3 | 0x4200000000000000000000000000000000000002 | ExecutionAssertionStatusReporter | ✓ | ✓ | ✗ | ✗ |
| 4 | 0x4200000000000000000000000000000000000003 | GovernanceToken | ✓ | ✓ | ✗ | ✓ |
| 5 | 0x4200000000000000000000000000000000000004 | L1ERC721Bridge | ✓ | ✓ | ✓ | ✗ |
| 6 | 0x4200000000000000000000000000000000000005 | L1ERC20Bridge | ✓ | ✓ | ✓ | ✗ |
| 7 | 0x4200000000000000000000000000000000000007 | L2CrossDomainMessenger | ✓ | ✓ | ✓ | ✓ |
| 8 | 0x4200000000000000000000000000000000000008 | GasPriceOracle | ✓ | ✓ | ✗ | ✓ |
| 9 | 0x4200000000000000000000000000000000000009 | - | - | - | - | - |
| 10 | 0x420000000000000000000000000000000000000A | OptimismMintableERC20Factory | ✓ | ✓ | ✓ | ✗ |
| 11 | 0x420000000000000000000000000000000000000B | ProxyAdmin | ✓ | - | ✗ | ✗ |
| 12 | 0x420000000000000000000000000000000000000C | BaseFeeVault | ✓ | ✓ | ✓ | ✗ |
| 13 | 0x420000000000000000000000000000000000000D | L1FeeVault | ✓ | ✓ | ✓ | ✗ |
| 14 | 0x420000000000000000000000000000000000000E | SequencerFeeVault | ✓ | ✓ | ✓ | ✗ |
| 15 | 0x420000000000000000000000000000000000000F | OptimismMintableERC721Factory | ✓ | ✓ | ✓ | ✗ |
| 16 | 0x4200000000000000000000000000000000000010 | L2StandardBridge | ✓ | ✓ | ✓ | ✓ |
| 17 | 0x4200000000000000000000000000000000000011 | OptimismERC20 | ✓ | ✓ | ✓ | ✗ |
| 18 | 0x4200000000000000000000000000000000000012 | L2ERC721Bridge | ✓ | ✓ | ✓ | ✗ |
| 19 | 0x4200000000000000000000000000000000000013 | OptimismMintableERC721 | ✓ | ✓ | ✓ | ✗ |
| 20 | 0x4200000000000000000000000000000000000014 | OptimismUsdc | ✓ | ✓ | ✓ | ✗ |
| 21 | 0x4200000000000000000000000000000000000015 | L1Block | ✓ | ✓ | ✗ | ✓ |
| 22 | 0x4200000000000000000000000000000000000016 | L2ToL1MessagePasser | ✓ | ✓ | ✗ | ✓ |
| 23 | 0x4200000000000000000000000000000000000017 | Create2Deployer | ✓ | ✓ | ✗ | ✗ |
| 24 | 0x4200000000000000000000000000000000000018 | ProxyAdmin | ✓ | - | ✗ | ✗ |
| 25 | 0x4200000000000000000000000000000000000019 | LegacyERC20NativeToken | ✓ | ✓ | ✓ | ✓ |
| ... | ... | ... | ... | ... | ... | ... |
| 40+ | 0x420000000000000000000000000000000000XX | AccountAbstraction, Gaming, Predeployed Tokens | ✓ | ✓ | ✓ | ✓ |

**주의**: 위 테이블은 주요 predeploy 계약 25개를 보여줍니다. OP Stack 사양에 따르면 총 40개 이상의 predeploy 계약이 정의되어 있으며, 체인 설정(preset)에 따라 선택적으로 활성화됩니다. 예: GameToken, AccountAbstraction 관련 predeploy 등.

### 6.5 Storage 값 예시 (JSON)

```json
{
  "0x4200000000000000000000000000000000000015": {
    "0x0000000000000000000000000000000000000000000000000000000000000001": "0x1234567890abcdef",
    "0x0000000000000000000000000000000000000000000000000000000000000002": "0x66xxxxx",
    "0x0000000000000000000000000000000000000000000000000000000000000003": "0x12345678"
  },
  "0x4200000000000000000000000000000000000010": {
    "0x0000000000000000000000000000000000000000000000000000000000000000": "0x01",
    "0x0000000000000000000000000000000000000000000000000000000000000001": "0x4200000000000000000000000000000000000007",
    "0x0000000000000000000000000000000000000000000000000000000000000002": "0x12345..."
  }
}
```

---

## 에러 시나리오

### 에러 타입별 분석

#### 1. 입력 검증 에러

**시나리오 1.1**: DeployConfigPath 없음
```
Error: open /path/to/deploy-config.json: no such file or directory

원인: Phase 4에서 deploy-config.json 생성 실패 또는 경로 오류
처리: Phase 4 로그 확인, 배포 설정 재검증

**코드 위치**: `trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 41-50
```go
if _, err := os.Stat(deployContractsConfig.DeployConfigPath); err != nil {
    return fmt.Errorf("deploy config not found: %w", err)
}
```
```

**시나리오 1.2**: DeployOutputPath 없음
```
Error: stat /tmp/deploy-output: no such file or directory

원인: Phase 4 바이너리 실행 실패
처리: Phase 4 deploy-contracts 로그 확인

**코드 위치**: `trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 55-65
```go
if _, err := os.Stat(deployContractsConfig.DeployOutputPath); err != nil {
    return fmt.Errorf("deploy output not found: %w", err)
}
```
```

**시나리오 1.3**: OptimismPortalProxy 미설정
```
Error: OptimismPortalProxy not set in config

원인: DeployConfig 생성 시 L1 배포 주소 로드 실패
처리: deploy-output.json 확인, Phase 4 배포 주소 검증

**코드 위치**: `op-chain-ops/genesis/config.go` 라인 850-855
```go
if dc.L1Contracts.OptimismPortalProxy == (common.Address{}) {
    return nil, fmt.Errorf("OptimismPortalProxy not set")
}
```
```

#### 2. 바이너리 실행 에러

**시나리오 2.1**: tokamak-deployer 바이너리 없음
```
Error: exec: "tokamak-deployer": executable file not found in $PATH

원인: PATH에 tokamak-deployer 없음
처리: 바이너리 경로 확인, 환경 변수 설정

코드:
  cmd := exec.CommandContext(ctx, "tokamak-deployer", "generate-genesis", ...)
  // "tokamak-deployer" not in PATH
```

**시나리오 2.2**: 바이너리 버전 불일치
```
Error: tokamak-deployer version 0.9.0 not compatible with phase 5

원인: 바이너리 버전이 Phase 5 코드와 호환되지 않음
처리: 바이너리 업그레이드 (v1.0.1 이상 필요)

처리:
  tokamak-deployer --version
  // Expected: v1.0.1 or later
  // Got: v0.9.0
```

**시나리오 2.3**: 바이너리 실행 시간 초과
```
Error: context deadline exceeded (timeout after 30 minutes)

원인: 매우 큰 genesis 생성 또는 I/O 병목
처리: 타임아웃 값 증가, 시스템 리소스 확인

코드:
  ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
  // GenerateGenesis 실행
  // Timeout after 30 minutes
```

#### 3. 데이터 구조 에러

**시나리오 3.1**: Predeploy Immutable 불일치
```
Error: immutable value size mismatch for LegacyERC20NativeToken

원인: Solc 바이트코드의 예상 Immutable 크기와 실제 값 크기 불일치
처리: Solc 버전 확인, 컨트랙트 ABI 검증

코드:
  finalBytecode := injectImmutables(bytecode, immutablesForPredeploy)
  // Bytecode expected 32 bytes for RemoteToken, got 20 bytes
```

**시나리오 3.2**: Storage 값 타입 에러
```
Error: cannot convert []byte to uint256 for L1Block.number

원인: Storage 값을 잘못된 타입으로 설정
처리: NewL2StorageConfig에서 타입 검증 추가

코드:
  db.SetState(L1BlockAddress, numberSlot, value)
  // value must be [32]byte, got []byte of length 8
```

**시나리오 3.3**: Proxy Admin 주소 미설정
```
Error: ProxyAdminAddr not set in DeployConfig

원인: Phase 4에서 ProxyAdmin 컨트랙트 배포 실패
처리: Phase 4 배포 로그 확인

코드 (라인 60):
  if config.ProxyAdminAddr == (common.Address{}) {
      return nil, fmt.Errorf("ProxyAdminAddr not set")
  }
```

#### 4. 파일 시스템 에러

**시나리오 4.1**: 출력 디렉토리 쓰기 권한 없음
```
Error: permission denied: /data/stacks/{stackId}/genesis.json

원인: 출력 디렉토리의 쓰기 권한 부족
처리: 디렉토리 권한 확인 (chmod 755 /data/stacks/{stackId})

코드:
  err = ioutil.WriteFile(filepath.Join(outPath, "genesis.json"), data, 0644)
  // Permission denied
```

**시나리오 4.2**: 디스크 용량 부족
```
Error: no space left on device

원인: genesis.json 생성 중 디스크 부족
처리: 디스크 정리, /tmp 확인

코드:
  err = ioutil.WriteFile(genesisPath, jsonData, 0644)
  // File size: 1.2GB, Available: 500MB
```

**시나리오 4.3**: 파일 이미 존재 (race condition)
```
Error: file already exists: genesis.json

원인: 동시 배포 또는 이전 배포 미정리
처리: 출력 디렉토리 초기화

코드:
  err = os.Rename(src, dst)
  // dst already exists
```

#### 5. Hard Fork 에러

**시나리오 5.1**: EcotoneTime이 미래 시간
```
경고: EcotoneTime (2026-05-01) is after L2 genesis timestamp (2026-04-16)

원인: Hard Fork 시간이 Genesis 시간보다 나중임
영향: Ecotone 활성화되지 않음 (정상 동작)
처리: 없음 (의도된 동작)

코드 (라인 40):
  if config.EcotoneTime != nil && config.EcotoneTime <= block.Time() {
      // Apply Ecotone changes
  }
```

**시나리오 5.2**: Hard Fork 시간 순서 오류
```
에러: FjordTime (2026-03-01) is before EcotoneTime (2026-04-01)

원인: Hard Fork 시간들이 오름차순이 아님
처리: Hard Fork 시간 검증, 순서 수정

코드:
  if config.EcotoneTime != nil && config.FjordTime != nil {
      if *config.EcotoneTime > *config.FjordTime {
          return nil, fmt.Errorf("fork times not in order")
      }
  }
```

#### 6. Cannon Prestate 에러 (Fault Proof)

**시나리오 6.1**: Prestate 다운로드 실패
```
Error: failed to download cannon prestate from S3

원인: S3 버킷에 접근 불가 또는 파일 없음
처리: AWS 자격증명 확인, prestate 버전 확인

코드:
  prestate, err := downloadPrestate(ctx, 
      deployContractsConfig.CannonPresateBucket,
      deployContractsConfig.CannonPresatePath)
  // Error: access denied
```

**시나리오 6.2**: Prestate 바이너리 손상
```
Error: invalid cannon prestate format

원인: Prestate 파일이 손상되거나 포맷이 잘못됨
처리: Prestate 다시 다운로드, 체크섬 검증

코드:
  err = validatePrestate(prestatePath)
  // Magic bytes not found
```

**시나리오 6.3**: Fault Proof 설정 불일치
```
Error: FaultProofEnabled but FaultGameFactoryAddress not set

원인: Fault Proof 활성화는 했으나 L1 팩토리 주소 미설정
처리: L1 배포에서 FaultGameFactory 배포 확인

코드:
  if config.FaultProofEnabled && 
     config.L1Contracts.FaultGameFactory == (common.Address{}) {
      return nil, fmt.Errorf("FaultGameFactory address required")
  }
```

#### 7. 주입(Injection) 에러

**시나리오 7.1**: DRB 주입 실패
```
Error: inject-drb: failed to add DeploymentBridge to genesis

원인: DRB 계약 바이트코드 로드 실패
처리: DRB 컨트랙트 아티팩트 확인

코드:
  if needsDRBInjection(config.Preset) {
      err = runInjectDRB(genesisPath, drbConfig)
      // Failed to load DRB bytecode
  }
```

**시나리오 7.2**: USDC 주입 중복
```
경고: USDC already exists in genesis, skipping injection

원인: genesis.json에 이미 USDC 계정 존재
영향: USDC 주입 스킵 (정상 동작)
처리: 없음

코드:
  if needsUSDCInjection(config.Preset) {
      // Check if USDC already exists
      if exists(genesis, USDCAddress) {
          log.Warn("USDC already exists")
          return nil
      }
  }
```

**시나리오 7.3**: Paymaster 주입 실패
```
Error: inject-paymaster: EntryPoint not found in genesis

원인: EntryPoint 계약이 genesis에 없음
처리: Account Abstraction 설정 확인

코드:
  if config.EnableAccountAbstraction {
      err = runInjectPaymaster(genesisPath, paymasterConfig)
      // EntryPoint address not in genesis.alloc
  }
```

#### 8. 메모리 부족 에러

**시나리오 8.1**: Genesis Alloc 메모리 초과
```
Error: cannot allocate memory for genesis allocs

원인: 40K+ 계정/계약 정보를 메모리에 로드할 수 없음
처리: 스트리밍 처리 구현, 메모리 증가

코드:
  db := state.NewMemoryStateDB(genesis)
  // Cannot allocate 2GB for allocs
```

**시나리오 8.2**: 바이트코드 캐시 메모리 초과
```
경고: Bytecode cache size exceeded 1GB

원인: 40개 Predeploy의 바이트코드가 메모리에 누적
영향: GC 압박 증가, 성능 저하
처리: 바이트코드 스트리밍 로드

코드:
  for _, predeploy := range Predeploys {
      bytecode, err := bindings.GetBytecode(predeploy.Name)
      // bytecode cache: 1.2GB
  }
```

#### 9. JSON 직렬화 에러

**시나리오 9.1**: genesis.json 인코딩 오류
```
Error: json: unsupported type for encoding

원인: big.Int 또는 다른 복잡한 타입 JSON 변환 불가
처리: 커스텀 Marshal 메서드 확인

코드:
  data, err := json.MarshalIndent(genesis, "", "  ")
  // json: unsupported type *big.Int
```

**시나리오 9.2**: rollup.json 필드 누락
```
에러: json 스키마 검증 실패: missing required field "genesis.l1"

원인: RollupConfig 구성 중 필수 필드 누락
처리: RollupConfig() 함수에서 검증 추가

코드:
  if rollupCfg.Genesis.L1.Hash == (common.Hash{}) {
      return nil, fmt.Errorf("genesis.l1.hash required")
  }
```

#### 10. 타임스탬프 에러

**시나리오 10.1**: L2 Genesis Timestamp가 L1보다 이전
```
에러: L2 genesis timestamp (1704067200) before L1 genesis (1704067300)

원인: 설정 오류
처리: L1 시작 블록 확인, 타임스탬프 재검증

코드:
  if genesis.Timestamp < l1StartBlock.Time() {
      return nil, fmt.Errorf("L2 genesis before L1")
  }
```

#### 11. 배포 상태 에러

**시나리오 11.1**: Resume 중 설정 불일치
```
에러: cannot resume: L2ChainID changed (901 → 902)

원인: 이전 배포의 설정과 현재 설정이 다름
처리: 배포 상태 초기화 후 재시작

코드:
  if resumeDeployment {
      prevConfig, err := loadPreviousConfig(deployOutputPath)
      if prevConfig.L2ChainId != config.L2ChainId {
          return fmt.Errorf("chain ID mismatch")
      }
  }
```

#### 12. 권한 및 접근 제어 에러

**시나리오 12.1**: Vault 수신자 주소 유효하지 않음
```
에러: invalid vault recipient address: 0x00...00

원인: 영 주소 또는 잘못된 체크섬
처리: 주소 검증, 기본값(deployer)으로 대체

코드:
  if config.SequencerFeeVaultRecipient == (common.Address{}) {
      config.SequencerFeeVaultRecipient = config.Deployers[0]
  }
  if !isValidAddress(config.SequencerFeeVaultRecipient) {
      return nil, fmt.Errorf("invalid recipient address")
  }
```

---

## 알려진 함정 및 개선 포인트

### 함정 1: Proxy vs Implementation 저장소

**설명**:
Proxy의 저장소와 Implementation의 저장소는 **같은 주소**에서 작동한다.
ERC-1967 Proxy 패턴에서:
- Proxy 바이트코드 → delegatecall → Implementation 로직
- 모든 저장소 읽/쓰 → **Proxy 주소** (0x4200...0007) 기준
- Implementation은 로직만 제공, 저장소는 사용 안 함

**영향**:
```go
// 올바른: Proxy 주소에 저장소 설정
db.SetState(proxyAddr, slotHash, value)

// 오류: Implementation 주소에 저장소 설정 (무시됨)
db.SetState(implAddr, slotHash, value)
```

**해결**:
NewL2StorageConfig()에서 모든 저장소 값을 Proxy 주소 기준으로 설정
(라인 1300+에서 확인)

---

### 함정 2: Immutables는 바이트코드에 컴파일됨

**설명**:
Immutables는 생성자가 아니라 **바이트코드에 직접 포함**된다.
따라서 한번 배포되면 변경 불가능하다.

```go
// 잘못된 시도: 배포 후 Immutable 변경
deployment := immutables.Deploy(db, config1)
// ... 나중에
deployment = immutables.Deploy(db, config2) // 새로운 구현체 배포됨
```

**영향**:
- 배포 후 Immutable 변경 불가
- Genesis 생성 시 모든 Immutable이 최종 결정되어야 함

**해결**:
```go
// DeployContractsInput에서 모든 Immutable 값 먼저 검증
if config.RemoteTokenAddress == (common.Address{}) {
    return fmt.Errorf("RemoteToken must be set before deployment")
}
```

---

### 함정 3: L1Block은 특수 컨트랙트

**설명**:
L1Block (0x4200...0015)은 일반 Predeploy와 다르다:
- Genesis에서 초기값만 설정
- 이후 **Sequencer가 매 블록마다 업데이트** (시스템 트랜잭션)
- 저장소 값은 L1의 최신 정보를 반영

**영향**:
```
Genesis L1Block state:
  number: 5678 (L1 시작 블록)
  timestamp: 1704067200
  
L2 블록 1: Sequencer 시스템 트랜잭션으로 L1Block 업데이트
  → L1Block.number = 5679
  → L1Block.timestamp = 1704067212
  
L2 블록 2: 다시 업데이트
  → L1Block.number = 5680
  → ...
```

**해결**:
- Genesis의 L1Block 값은 **L1 시작 블록만 반영**
- 이후 동적 업데이트는 Sequencer가 담당
- 신경쓸 필요 없음 (정상 설계)

---

### 함정 4: StorageConfig의 슬롯 해싱

**설명**:
Solidity 저장소 슬롯은 단순 uint256이 아니라, 때로는 **keccak256 해시**로 계산된다.

```solidity
// 예: mapping(address => uint256) public balances;
// 슬롯 1에 mapping이 있으면
// balances[0xABC...] = keccak256(abi.encodePacked(0xABC, 1))
```

**영향**:
```go
// 오류: 직접 슬롯에 저장 (복잡한 데이터 타입의 경우 작동 안 함)
db.SetState(addr, common.Hash{}, value)

// 올바른: keccak256 슬롯 계산
slot := crypto.Keccak256Hash([]byte{...})
db.SetState(addr, slot, value)
```

**해결**:
NewL2StorageConfig()에서 각 Predeploy의 저장소 레이아웃을 문서화
(필요시 Solidity 컴파일러의 storage-layout.json 참조)

---

### 함정 5: Chain Config의 Block Numbers vs Block Times

**설명**:
Geth ChainConfig는 **블록 번호** 기준으로 Hard Fork를 정의하지만,
OP Stack은 **타임스탬프** 기준이다.

```go
// Geth: 블록 번호로 정의
HomesteadBlock: 1,000,000

// OP Stack: 타임스탬프로 정의
EcotoneTime: 1704067200
```

**영향**:
```go
// Genesis의 ChainConfig
Config: {
    Ecotone: 0, // 블록 0부터 활성화
}
// RollupConfig
EcotoneTime: 1704067200, // 타임스탐프 기준

// 둘이 일치해야 함!
```

**해결**:
NewL2Genesis()와 RollupConfig()에서 일관성 확인
(라인 52-55, 873-879에서 동일한 시간 값 사용)

---

### 함정 6: ProxyAdmin 주소의 특수성

**설명**:
ProxyAdmin (0x4200...0018)은 **프록시 관리자 역할**을 하는 특수 컨트랙트다.
하지만 자신도 **Proxy + Implementation 구조**다.

```
ProxyAdmin Proxy (0x4200...0018)
  └─ delegates to ProxyAdmin Implementation
      └─ 다른 Proxy들의 upgrade() 함수 호출 권한 관리
```

**영향**:
- ProxyAdmin 주소를 직접 호출하면 delegatecall이 작동
- ProxyAdmin이 Proxy를 관리하려면 owner가 설정되어야 함

**해결**:
ProxyAdmin의 owner를 배포자 주소로 설정
(NewL2StorageConfig()에서 확인)

---

### 함정 7: Account Abstraction (EntryPoint) 주소

**설명**:
EntryPoint (0x...0005)는 ERC-4337 표준 주소로,
**모든 OP Stack 체인에서 동일**하다 (cross-chain 호환성).

```go
// 반드시 이 주소여야 함
EntryPointAddress := common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032")
```

**영향**:
- 다른 주소를 사용하면 지갑/DEX에서 EntryPoint를 찾지 못함
- Account Abstraction 호환 불가능

**해결**:
NewL2ImmutableConfig()에서 EntryPoint 주소 강제 검증
```go
if config.EntryPointAddress != ExpectedEntryPoint {
    return nil, fmt.Errorf("invalid EntryPoint address")
}
```

---

### 함정 8: 주입(Injection) 순서와 중복

**설명**:
DRB, USDC, Paymaster 주입은 **특정 순서**로 진행되어야 하며,
중복 주입 시 충돌이 발생할 수 있다.

```
1. genesis.json 생성
2. DRB 주입 (필요 시)
3. USDC 주입 (필요 시)
4. Paymaster 주입 (필요 시)
```

**영향**:
```go
// 오류: USDC를 두 번 주입
runInjectUSDS(...) // 첫 번째 주입
runInjectUSDS(...) // 두 번째 주입 → 충돌!

// 올바른: Preset에 따라 한 번만
if needsUSDCInjection(config.Preset) {
    runInjectUSDS(...)
}
```

**해결**:
DeployContracts()에서 각 주입의 조건을 상호 배타적으로 설정
(라인 381-407에서 확인)

---

### 개선 포인트 1: Genesis 검증 강화

**현재 문제**:
- Genesis 생성 후 모든 계약의 저장소가 올바른지 검증 안 함
- Bytecode 크기 검증 없음

**제안**:
```go
// 생성 후 검증 함수 추가
func ValidateL2Genesis(genesis *core.Genesis) error {
    for addr, acct := range genesis.Alloc {
        // 1. Proxy인 경우 bytecode 길이 확인
        if isPredeploy(addr) && len(acct.Code) < 100 {
            return fmt.Errorf("proxy bytecode too short: %v", addr)
        }
        
        // 2. 저장소 슬롯 범위 확인
        if len(acct.Storage) > 10000 {
            return fmt.Errorf("too many storage slots: %v", addr)
        }
        
        // 3. 잔액 범위 확인
        if acct.Balance.Cmp(big.NewInt(1e18 * 1000)) > 0 {
            return fmt.Errorf("balance too large: %v", addr)
        }
    }
    return nil
}
```

---

### 개선 포인트 2: 병렬 Predeploy 배포

**현재 문제**:
- 40개 Predeploy를 순차 배포 (느림)
- 불변값 주입도 순차적

**제안**:
```go
// Worker pool 패턴으로 병렬화
func DeployPredeploys(db vm.StateDB, config *Config, numWorkers int) error {
    predeploys := getPredeploys()
    results := make(chan DeployResult, len(predeploys))
    
    for i := 0; i < numWorkers; i++ {
        go deployWorker(i, predeploys, results)
    }
    
    for i := 0; i < len(predeploys); i++ {
        res := <-results
        if res.Err != nil {
            return res.Err
        }
        db.SetCode(res.Addr, res.Bytecode)
    }
    return nil
}
```

예상 성능: 2-3배 빠름

---

### 개선 포인트 3: Streaming JSON 기록

**현재 문제**:
- 전체 genesis를 메모리에 로드 후 JSON 변환
- 1GB 이상 메모리 사용

**제안**:
```go
// JSON encoder를 파일에 직접 기록
func WriteGenesisStreaming(genesis *core.Genesis, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    
    // alloc을 청크 단위로 기록
    return encoder.Encode(genesis)
}
```

예상 메모리 감소: 500MB → 50MB

---

### 개선 포인트 4: Preset별 검증 강화

**현재 문제**:
- Preset (basic, full, gaming)에 따라 다른 주입이 필요하지만,
  검증 부족

**제안**:
```go
var PresetRequirements = map[string]PresetSpec{
    "basic": {
        InjectionSteps: []string{},
        Predeployer:    []string{},
    },
    "full": {
        InjectionSteps: []string{"usdc"},
        Predeployer:    []string{"Paymaster"},
    },
    "gaming": {
        InjectionSteps: []string{"drb", "usdc"},
        Predeployer:    []string{"GameLogic", "Paymaster"},
    },
}

func ValidatePresetConfig(preset string, config *DeployConfig) error {
    spec, ok := PresetRequirements[preset]
    if !ok {
        return fmt.Errorf("unknown preset: %s", preset)
    }
    
    for _, step := range spec.InjectionSteps {
        // 각 스텝의 사전 조건 검증
    }
    return nil
}
```

---

### 개선 포인트 5: L1Block 업데이트 최적화

**현재 문제**:
- Sequencer가 매 블록마다 L1Block을 업데이트하는데,
  불필요한 업데이트 발생 가능

**제안**:
```go
// L1Block 업데이트 배치 처리
func UpdateL1BlockBatch(db vm.StateDB, l1Blocks []BlockInfo) error {
    // 여러 L1 블록 정보를 한 번에 처리
    for _, block := range l1Blocks {
        db.SetState(L1BlockAddr, numberSlot, block.Number)
        db.SetState(L1BlockAddr, timestampSlot, block.Timestamp)
        // ...
    }
    return nil
}
```

성능 개선: L1 폴링 주기 감소 가능

---

### 개선 포인트 6: 오류 복구 전략

**현재 문제**:
- Predeploy 배포 중 하나 실패하면 전체 실패
- 부분 복구 불가능

**제안**:
```go
// 실패한 Predeploy만 재배포
func DeployWithRetry(db vm.StateDB, config *Config, maxRetries int) error {
    failed := []string{}
    
    for _, name := range predeploy.Names() {
        retries := 0
        for {
            err := deployPredeploy(db, name, config)
            if err == nil {
                break
            }
            retries++
            if retries > maxRetries {
                failed = append(failed, name)
                break
            }
            time.Sleep(time.Second * time.Duration(retries))
        }
    }
    
    if len(failed) > 0 {
        return fmt.Errorf("failed to deploy: %v", failed)
    }
    return nil
}
```

---

### 개선 포인트 7: Genesis 버전 관리

**현재 문제**:
- Genesis 스키마 변경 추적 안 됨
- Hard Fork 시 호환성 문제 가능

**제안**:
```go
// Genesis 메타데이터 추가
type GenesisMeta struct {
    Version       string    // "1.0.0"
    CreatedAt     time.Time
    CreatedBy     string    // "tokamak-deployer v1.0.1"
    ChainName     string    // "Tokamak Thanos Sepolia"
    NetworkID     uint64
    HardForks     []HardFork
}

// genesis.json에 메타데이터 저장
{
    "_meta": {
        "version": "1.0.0",
        "createdAt": "2026-04-16T...",
        ...
    },
    "config": { ... },
    "alloc": { ... }
}
```

---

### 개선 포인트 8: Immutables 검증 도구

**현재 문제**:
- Immutables 주입 후 검증 불가능
- 런타임에 문제 발생 (너무 늦음)

**제안**:
```go
// Immutables 검증 도구
func ValidateImmutables(bytecode []byte, expectedImmutables map[string]interface{}) error {
    for key, expected := range expectedImmutables {
        actual := extractImmutable(bytecode, key)
        if actual != expected {
            return fmt.Errorf("immutable mismatch: %s, expected %v, got %v", 
                key, expected, actual)
        }
    }
    return nil
}

// 배포 후 검증
bytecode := injectImmutables(originalBytecode, config)
if err := ValidateImmutables(bytecode, expectedMap); err != nil {
    return fmt.Errorf("immutable validation failed: %w", err)
}
```

---

## 결론

Phase 5 (L2 Genesis 생성)는 L1 배포 결과를 기반으로 L2의 초기 상태를 결정하는 중요한 단계다. 주요 특징:

1. **40개 이상의 Predeploy**: 각각 Proxy + Implementation으로 구성
2. **Immutables 주입**: 바이트코드에 불변값 포함
3. **Storage 초기화**: L1Block, 메신저, 브릿지 주소 등 설정
4. **Hard Fork 지원**: Ecotone, Fjord, Granite 활성화 시간
5. **Optional 주입**: DRB, USDC, Paymaster 추가 배포
6. **JSON 파일 생성**: genesis.json (1-2MB), rollup.json (1-2KB)

전체 흐름은 SDK의 `DeployContracts()` 함수에서 `tokamak-deployer` 바이너리를 호출하여 진행되며, 모든 핵심 로직은 `op-chain-ops/genesis/` 라이브러리 패키지에 구현되어 있다.

