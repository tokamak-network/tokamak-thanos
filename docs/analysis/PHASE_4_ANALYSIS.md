# Phase 4 분석: L1 컨트랙트 배포

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

Phase 4는 L1 체인(Ethereum Mainnet, Testnet/Sepolia)에 Thanos Stack의 핵심 컨트랙트를 배포하는 과정을 다룬다. 이 단계는:

- **역할**: Foundry 기반 `tokamak-deployer` 바이너리를 통해 L1 컨트랙트 체계적 배포
- **입력**: L1 RPC URL, L1 Beacon URL, 운영자 개인키(5개), 배포 설정
- **출력**: 
  - `deploy-output.json` (배포된 컨트랙트 주소들)
  - `rollup.json` (Rollup 설정 메타데이터)
  - `genesis.json` (L2 제네시스 파일)
- **특징**:
  - Fault Proof 활성화 옵션 지원
  - 배포 상태 추적 (Pending → InProgress → Completed)
  - 배포 중단 가능 (resume 기능)
  - 자동 기본값 생성 (L2 Chain ID, 배포 비용 추정)

---

## 용어 정의

### 배포 생명주기 관련

- **DeployContractState**: `settings.json`에 저장된 L1 배포 상태
  - `Status: DeployContractStatusInProgress` - 배포 진행 중
  - `Status: DeployContractStatusCompleted` - 배포 완료
- **Resume**: 중단된 배포를 재개하는 기능. 이전 배포 상태 파일이 존재할 때 활성화됨
- **Reuse Deployment**: 기존 구현체 주소를 재사용하는 옵션 (구현체가 여러 프록시에서 사용될 때)

### 컨트랙트 배포 관련

- **Operators**: L1 배포에 필요한 5개 운영자 계정
  - `AdminPrivateKey` - 배포 수행 및 초기화
  - `SequencerPrivateKey` - 트랜잭션 생성자
  - `BatcherPrivateKey` - 배치 제출자
  - `ProposerPrivateKey` - 상태 루트 제안자
  - `ChallengerPrivateKey` (선택) - 결함 증명 활성화 시만 필요
- **ChainConfiguration**: L1 블록 시간, L2 블록 시간, 배치 제출 주기 등
- **Fault Proof**: OptimismPortal의 결함 증명 메커니즘. 활성화 시 AnchorStateRegistry 배포 추가
- **Register Candidate**: 토카막 네트워크의 검증자 후보 등록 (별도 트랜잭션)

### tokamak-deployer 바이너리

- **Binary Version**: v1.0.0 (핀포인트)
- **Deploy Output**: `deploy-output.json` - 배포된 컨트랙트 주소 및 트랜잭션 해시
- **Genesis Generation**: 제네시스 파일 및 rollup.json 생성

---

## 상세 분석

### 4.1 백엔드 진입점: executeDeployments() 함수

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go`

**함수**: `executeDeployments()` (라인 439-653)

```go
func (s *ThanosStackDeploymentService) executeDeployments(ctx context.Context, stackId uuid.UUID) error
```

**역할**: 
- 배포 엔티티 필터링 및 순서 결정
- L1 배포 → AWS 인프라 배포 순차 실행
- 배포 상태 추적 및 데이터베이스 업데이트

**핵심 로직**:

1. **상태 업데이트**: 스택 상태를 `StackStatusDeploying`으로 설정 (라인 442)
2. **배포 필터링** (라인 476-501):
   - 대기 중인 배포 엔티티 조회
   - `DeployL1ContractsStep` ("deploy-l1-contracts") 필터링
   - `DeployInfraStep` ("deploy-aws-infra") 필터링
   - 순서 보장: L1 먼저, AWS 인프라 나중
3. **상태 채널 설정** (라인 449-519):
   - `statusChan` (배포 상태 업데이트)
   - `errChan` (최종 에러 반환)
4. **배포 루프** (라인 521-649):
   ```go
   for _, deployment := range pendingDeployments {
       // SDK 클라이언트 생성
       sdkClient, err := thanos.NewThanosSDKClient(...)
       
       // 상태 업데이트: InProgress
       statusChan <- entities.DeploymentStatusWithID{
           DeploymentID: deployment.ID,
           Status:       entities.DeploymentRunStatusInProgress,
       }
       
       switch deployment.Step {
       case "deploy-l1-contracts":
           // DeployL1Contracts 호출
       case "deploy-aws-infra":
           // DeployLocalInfrastructure 또는 DeployAWSInfrastructure 호출
       }
   }
   ```

**입력 구조** (`dtos.DeployL1ContractsRequest`):
```json
{
  "l1_rpc_url": "https://sepolia.infura.io/...",
  "admin_account": "0x...",          // private key (hex)
  "sequencer_account": "0x...",
  "batcher_account": "0x...",
  "proposer_account": "0x...",
  "challenger_account": "0x...",    // optional, for fault proof
  "batch_submission_frequency": 3600,
  "challenge_period": 604800,
  "output_root_frequency": 900,
  "l2_block_time": 12,
  "register_candidate": false,
  "reuse_deployment": false,
  "enable_fault_proof": false,
  "preset": "custom",
  "fee_token": "ETH"
}
```

**상태 머신** (라인 528-600):
```
Pending ──→ InProgress ──→ Success
             ↓
             Failed (error returned, status updated)
```

---

### 4.2 SDK 래퍼: DeployL1Contracts 함수

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/stacks/thanos/thanos_stack.go`

**함수**: `DeployL1Contracts()` (라인 129-173)

```go
func DeployL1Contracts(ctx context.Context, sdkClient *thanosStack.ThanosStack, 
    req *dtos.DeployL1ContractsRequest) error
```

**역할**:
- `dtos.DeployL1ContractsRequest` (백엔드 DTO) → `thanosStack.DeployContractsInput` (SDK 구조) 변환
- SDK의 `DeployContracts()` 메서드 호출

**핵심 변환**:

1. **ChainConfiguration 생성** (라인 132-138):
   ```go
   chainConfig := thanosTypes.ChainConfiguration{
       BatchSubmissionFrequency: uint64(req.BatchSubmissionFrequency),  // 초 단위
       ChallengePeriod:          uint64(req.ChallengePeriod),          // 초 단위
       OutputRootFrequency:      uint64(req.OutputRootFrequency),      // 초 단위
       L2BlockTime:              uint64(req.L2BlockTime),              // 초 단위
       L1BlockTime:              12,  // Ethereum 블록 시간 고정
   }
   ```

2. **Operators 생성** (라인 139-145):
   ```go
   operators := thanosTypes.Operators{
       AdminPrivateKey:      req.AdminAccount,
       SequencerPrivateKey:  req.SequencerAccount,
       BatcherPrivateKey:    req.BatcherAccount,
       ProposerPrivateKey:   req.ProposerAccount,
       ChallengerPrivateKey: req.ChallengerAccount,  // Fault Proof 필수
   }
   ```

3. **RegisterCandidate 조건부 설정** (라인 157-164):
   ```go
   if req.RegisterCandidate && req.RegisterCandidateParams != nil {
       contractDeploymentInput.RegisterCandidate = &thanosStack.RegisterCandidateInput{
           Amount:   req.RegisterCandidateParams.Amount,    // 보증금 (TON)
           Memo:     req.RegisterCandidateParams.Memo,
           NameInfo: req.RegisterCandidateParams.NameInfo,
           UseTon:   true,  // 현재 TON만 지원
       }
   }
   ```

4. **SDK 메서드 호출** (라인 166):
   ```go
   err := sdkClient.DeployContracts(ctx, &contractDeploymentInput)
   ```

---

### 4.3 SDK 핵심: ThanosStack.DeployContracts 메서드

**파일**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go`

**함수**: `DeployContracts()` (라인 23-407)

**역할**: 
- 네트워크 검증 (LocalDevnet 제외, Testnet/Mainnet만)
- 배포 상태 확인 (이미 완료된 경우 처리)
- 배포 설정 자동 생성 (L2 Chain ID, 배포 비용 추정)
- tokamak-deployer 바이너리 다운로드 및 실행
- Post-processing (Genesis 및 Rollup 파일 생성)

**네트워크 검증** (라인 24-31):
```go
if t.network == constants.LocalDevnet {
    return fmt.Errorf("LocalDevnet does not require contract deployment")
}
if t.network != constants.Testnet && t.network != constants.Mainnet {
    return fmt.Errorf("network %s does not support", t.network)
}
```

**ChainConfiguration 자동 생성** (라인 38-64):
```go
if deployContractsConfig.ChainConfiguration == nil {
    l1Client, err := ethclient.DialContext(ctx, deployContractsConfig.L1RPCurl)
    l1ChainID, err := l1Client.ChainID(ctx)
    
    // L1 체인 설정 상수에서 조회
    finalzationPeriodSeconds := constants.L1ChainConfigurations[l1ChainID.Uint64()].FinalizationPeriodSeconds
    l2OutputSubmissionInterval := constants.L1ChainConfigurations[l1ChainID.Uint64()].L2OutputOracleSubmissionInterval
    maxChannelDuration := constants.L1ChainConfigurations[l1ChainID.Uint64()].MaxChannelDuration
    
    deployContractsConfig.ChainConfiguration = &types.ChainConfiguration{
        BatchSubmissionFrequency: maxChannelDuration * l1BlockTime,
        ChallengePeriod:          finalzationPeriodSeconds,
        OutputRootFrequency:      l2BlockTime * l2OutputSubmissionInterval,
        L2BlockTime:              l2BlockTime,
        L1BlockTime:              l1BlockTime,
    }
}
```

**L2 Chain ID 생성** (라인 172):
```go
l2ChainID, err := utils.GenerateL2ChainId()  // 무작위 고유 ID 생성
```

**배포 상태 확인** (라인 99-129):
```go
if t.deployConfig.DeployContractState != nil {
    switch t.deployConfig.DeployContractState.Status {
    case types.DeployContractStatusCompleted:
        // 재배포 여부 사용자에게 확인
        if !t.usePromptInput {
            return nil  // 이미 완료됨 → 스킵
        }
    case types.DeployContractStatusInProgress:
        // resume 옵션 활성화
        isResume = true
    }
}
```

**Fault Proof 처리** (라인 224-259):
```go
if deployContractsConfig.EnableFaultProof {
    // 1. tokamak-thanos 저장소 클론
    err = t.cloneSourcecode(ctx, "tokamak-thanos", "https://github.com/tokamak-network/tokamak-thanos.git")
    
    // 2. AnchorStateRegistry.sol 패치
    if patchErr := patchAnchorStateRegistry(tokamakThanosDir); patchErr != nil {
        return fmt.Errorf("failed to patch AnchorStateRegistry.sol: %w", patchErr)
    }
    
    // 3. Cannon prestate 빌드 (병렬)
    g, gctx := errgroup.WithContext(ctx)
    g.Go(func() error {
        if buildErr := buildCannonPrestate(gctx, t.logger, tokamakThanosDir); buildErr != nil {
            return fmt.Errorf("failed to build cannon prestate: %w", buildErr)
        }
        prestateHash, hashErr := readPrestateHash(prestatePath)
        return hashErr
    })
    if err = g.Wait(); err != nil {
        return err
    }
}
```

**Admin 계정 잔액 확인** (라인 292-318):
```go
balance, err := l1Client.BalanceAt(ctx, adminAccount, nil)
gasPriceWei, err := l1Client.SuggestGasPrice(ctx)

estimatedCost := new(big.Int).Mul(gasPriceWei, estimatedDeployContracts)
estimatedCost.Mul(estimatedCost, big.NewInt(2))  // 2배 여유

if balance.Cmp(estimatedCost) < 0 {
    return fmt.Errorf("insufficient balance: %.4f ETH < %.4f ETH",
        utils.WeiToEther(balance), utils.WeiToEther(estimatedCost))
}
```

**Deploy 상태 파일 생성** (라인 329-336):
```go
t.deployConfig.DeployContractState = &types.DeployContractState{
    Status: types.DeployContractStatusInProgress,
}
err = t.deployConfig.WriteToJSONFile(t.deploymentPath)
// → settings.json에 저장됨
```

---

### 4.4 tokamak-deployer 바이너리 호출

**파일**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deployer_binary.go`

#### 4.4.1 바이너리 다운로드 (라인 20-49)

**함수**: `ensureTokamakDeployer()`

```go
func ensureTokamakDeployer(cacheDir string) (string, error) {
    // 캐시 위치: ~/.trh/bin/tokamak-deployer-v1.0.0
    binaryPath := filepath.Join(cacheDir, "tokamak-deployer-v1.0.0")
    
    if _, err := os.Stat(binaryPath); err == nil {
        return binaryPath, nil  // 캐시 히트
    }
    
    // 다운로드 URL 생성
    downloadURL := fmt.Sprintf(
        "https://github.com/tokamak-network/tokamak-thanos/releases/download/v1.0.0/tokamak-deployer-%s-%s",
        runtime.GOOS,  // "linux" 또는 "darwin"
        runtime.GOARCH, // "amd64" 또는 "arm64"
    )
    // 예: https://github.com/tokamak-network/tokamak-thanos/releases/download/v1.0.0/tokamak-deployer-darwin-arm64
    
    if err := downloadFile(downloadURL, binaryPath); err != nil {
        return "", fmt.Errorf("failed to download tokamak-deployer: %w", err)
    }
    os.Chmod(binaryPath, 0755)
    return binaryPath, nil
}
```

#### 4.4.2 Foundry 스크립트 실행

**tokamak-thanos 내 배포 스크립트**:
- 진입점: `tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/start-deploy.sh`
- 파일: `tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/Deploy.s.sol`
- 스크립트 위치: `packages/tokamak/contracts-bedrock/scripts/Deploy.s.sol` (Foundry 배포 로직)

**실행 흐름** (start-deploy.sh에서):

```bash
# 1. 배포 설정 파일 및 환경 변수 로드
handleScriptInput "$@"  # -c configPath -e envFilePath

# 2. 소스 코드 빌드 (필요시)
buildSource  # pnpm install, make submodules, forge build 등

# 3. 스마트 컨트랙트 배포
deployContracts() {
    export IMPL_SALT=$(openssl rand -hex 32)
    cd packages/tokamak/contracts-bedrock
    
    if [[ -n "$GAS_PRICE" && "$GAS_PRICE" -gt 0 ]]; then
        forge script scripts/Deploy.s.sol:Deploy \
          --private-key $GS_ADMIN_PRIVATE_KEY \
          --broadcast \
          --rpc-url $L1_RPC_URL \
          --slow --legacy --non-interactive \
          --with-gas-price $GAS_PRICE
    else
        forge script scripts/Deploy.s.sol:Deploy \
          --private-key $GS_ADMIN_PRIVATE_KEY \
          --broadcast \
          --rpc-url $L1_RPC_URL \
          --slow --legacy --non-interactive
    fi
}

# 4. L2 제네시스 및 롤업 설정 생성
generateL2Genesis() {
    forge script scripts/L2Genesis.s.sol:L2Genesis --rpc-url $L1_RPC_URL
    op-node genesis l2 ...
}
```

**script 호출 방식**:
```bash
# start-deploy.sh 실행
./start-deploy.sh deploy \
  -c ./deploy-config.json \
  -e ./deploy.env

# 또는 재배포
./start-deploy.sh redeploy \
  -c ./deploy-config.json \
  -e ./deploy.env

# 또는 전체 프로세스 (설치 + 빌드 + 배포 + 제네시스)
./start-deploy.sh all \
  -c ./deploy-config.json \
  -e ./deploy.env
```

**.env 파일 설정** (tokamak-thanos 루트에서 필요):

환경 변수가 필요합니다:

| 변수명 | 타입 | 필수여부 | 설명 | 예시 |
|-------|------|---------|------|------|
| `L1_RPC_URL` | URL | 필수 | L1 노드 RPC 엔드포인트 | https://sepolia.infura.io/v3/... |
| `GS_ADMIN_PRIVATE_KEY` | Hex | 필수 | 배포 계정의 Private Key (0x 포함) | 0x1234567890abcdef... |
| `L1_CHAIN_ID` | Number | 필수 | L1 체인 ID (deploy-config.json에서 읽음) | 11155111 (Sepolia) |
| `DEPLOY_CONFIG_PATH` | Path | 필수 | 배포 설정 JSON 파일 경로 | /tmp/deploy-config.json |
| `GAS_PRICE` | Number | 선택 | Gas 가격 (wei) | 20000000000 |

**deploy-config.json 예시** (deploy/deploy-config.json에서):
```json
{
  "l1ChainID": 11155111,
  "l2ChainID": 1729,
  "l1BlockTime": 12,
  "l2BlockTime": 12,
  "maxSequencerDrift": 600,
  "sequencerWindowSize": 3600,
  "channelTimeout": 300,
  "p2pSequencerAddress": "0x...",
  "batchInboxAddress": "0x",
  "batchSenderAddress": "0x...",
  "l2OutputOracleSubmissionInterval": 900,
  "l2OutputOracleStartingTimestamp": 1234567890,
  "l2OutputOracleProposer": "0x...",
  "l2OutputOracleChallenger": "0x...",
  "finalSystemOwner": "0x...",
  "superchainConfigGuardian": "0x...",
  "proxyAdminOwner": "0x...",
  "baseFeeVaultRecipient": "0x...",
  "l1FeeVaultRecipient": "0x...",
  "sequencerFeeVaultRecipient": "0x...",
  "governanceTokenOwner": "0x...",
  "governanceTokenSymbol": "THANOS",
  "governanceTokenName": "Thanos Token"
}
```

**검증**:
- .env 파일이 .gitignore에 포함되어야 함 (Private Key 보안)
- deploy-config.json은 배포 전에 유효성 검사됨
- Foundry 스크립트는 `Deploy.s.sol:Deploy` 컨트랙트의 `run()` 함수 실행

**출력 파일**:
- `deployments/{l1ChainID}-deploy.json`: 배포된 컨트랙트 주소 및 트랜잭션 해시
- `build/genesis.json`: L2 제네시스 파일
- `build/rollup.json`: Rollup 설정 메타데이터

#### 4.4.3 Genesis/Rollup 파일 생성 (라인 104-118)

**함수**: `runGenerateGenesis()`

```go
func runGenerateGenesis(ctx context.Context, binaryPath string, opts genesisOpts) error {
    args := []string{
        "generate-genesis",
        "--deploy-output", opts.DeployOutputPath,
        "--config", opts.ConfigPath,
        "--out", opts.OutPath,
    }
    if opts.Preset != "" {
        args = append(args, "--preset", opts.Preset)
    }
    return runBinaryCommand(ctx, binaryPath, args)
}
```

**실행 명령어 예**:
```bash
/Users/theo/.trh/bin/tokamak-deployer-v1.0.0 \
  generate-genesis \
  --deploy-output "/path/to/deployment/deploy-output.json" \
  --config "/path/to/deployment/deploy-config.json" \
  --out "/path/to/deployment/genesis.json"
```

**출력 파일**:
- `genesis.json`: L2 제네시스 파일
- `rollup.json`: Rollup 설정 및 메타데이터
- Binary는 자동으로 다음을 주입:
  - DRB (Distributed Randomness Beacon) 설정
  - USDC 토큰 설정 (필요 시)
  - MultiTokenPaymaster 설정 (필요 시)
  - L1Block Isthmus 바이트코드 패치

---

### 4.5 배포 구성 생성

**파일**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go`

**함수**: `makeDeployContractConfigJsonFile()` (라인 284)

```go
err = makeDeployContractConfigJsonFile(ctx, l1Client, operators, deployContractsTemplate, deployConfigFilePath)
```

**생성 파일**: `deploy-config.json`

**샘플 구조**:
```json
{
  "adminAddress": "0x...",
  "sequencerAddress": "0x...",
  "batcherAddress": "0x...",
  "proposerAddress": "0x...",
  "challengerAddress": "0x...",
  "l1ChainId": 11155111,          // Sepolia
  "l2ChainId": 123456,            // 자동 생성됨
  "l1BlockTime": 12,
  "l2BlockTime": 12,
  "batchSubmissionFrequency": 3600,
  "challengePeriod": 604800,
  "outputRootFrequency": 900,
  "reuseDeployment": false,
  "enableFaultProof": false,
  "preset": "custom",
  "feeToken": "ETH"
}
```

---

### 4.6 배포 설정 검증

**파일**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go`

**함수**: `deployContractsConfig.Validate()` (라인 66)

**검증 항목**:
1. L1 RPC URL 연결 가능성
2. Admin 계정 잔액 충분성 (RegisterCandidate 시)
3. Operator 개인키 유효성
4. ChainConfiguration 파라미터 범위 검증

---

### 4.7 배포 상태 추적 및 최종화

**파일**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/types/configuration.go`

**구조**: `Config` 및 `DeployContractState`

```go
type Config struct {
    Stack                    string  `json:"stack"`                  // "thanos"
    Network                  string  `json:"network"`                // "testnet", "mainnet"
    AdminPrivateKey          string  `json:"admin_private_key"`
    SequencerPrivateKey      string  `json:"sequencer_private_key"`
    BatcherPrivateKey        string  `json:"batcher_private_key"`
    ProposerPrivateKey       string  `json:"proposer_private_key"`
    ChallengerPrivateKey     string  `json:"challenger_private_key,omitempty"`
    DeploymentFilePath       string  `json:"deployment_file_path"`
    DeployContractState      *DeployContractState `json:"deploy_contract_state,omitempty"`
    EnableFraudProof         bool    `json:"enable_fraud_proof"`
    FeeToken                 string  `json:"fee_token"`
    Preset                   string  `json:"preset"`
}

type DeployContractState struct {
    Status DeployContractStatus `json:"status"`  // InProgress, Completed
}

const (
    DeployContractStatusInProgress DeployContractStatus = iota + 1  // 1
    DeployContractStatusCompleted                                    // 2
)
```

**저장 위치**: `{deploymentPath}/settings.json`

**업데이트 시점**:
1. 배포 시작 전: `Status = DeployContractStatusInProgress`
2. 배포 완료 후: `Status = DeployContractStatusCompleted`

**Resume 메커니즘**:
- 배포 중 실패 → `Status`는 `InProgress`로 유지
- 재배포 시: `isResume = true` → `runDeployContracts` 다시 실행
- 기존 `deploy-output.json` 있으면 마지막 상태부터 재개

---

### 4.8 Candidate 등록 (RegisterCandidate)

**파일**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go`

**함수**: `verifyRegisterCandidates()` (라인 395)

**조건**:
- `req.RegisterCandidate == true`
- `req.RegisterCandidateParams != nil`

**처리 과정**:
1. Safe Wallet 설정 (라인 391):
   ```go
   if err := t.setupSafeWallet(ctx, t.deploymentPath); err != nil {
       return err
   }
   ```

2. Candidate 검증 및 등록 (라인 395):
   ```go
   verifyRegisterError := t.verifyRegisterCandidates(ctx, registerCandidate)
   ```

3. 추가 정보 표시 (라인 401):
   ```go
   t.DisplayRegistrationAdditionalInfo(ctx, registerCandidate)
   ```

**입력 구조** (`RegisterCandidateInput`):
```go
type RegisterCandidateInput struct {
    Amount    *big.Int  // 보증금 (TON)
    Memo      string
    NameInfo  interface{}
    UseTon    bool      // 항상 true
}
```

---

### 4.9 배포 컨트랙트 순서 및 목록

**배포 순서 및 의존성 분석**:

Foundry 스크립트는 다음 계층적 순서로 컨트랙트를 배포합니다. 배포 순서는 `Deploy.s.sol`의 `_run()` 함수에 의해 정의됩니다:

#### Phase 1: Superchain 설정 (setupSuperchain)

| # | 컨트랙트명 | 종류 | 파일:라인 | 설명 |
|---|-----------|------|----------|------|
| 1 | AddressManager | 구현체 | src/legacy/AddressManager.sol:547 | 주소 관리 레지스트리 (레거시) |
| 2 | ProxyAdmin | 구현체 | src/universal/ProxyAdmin.sol:558 | 프록시 업그레이드 관리자 |
| 3 | SuperchainConfigProxy | 프록시 | src/universal/Proxy.sol:360 | SuperchainConfig 프록시 (ERC1967) |
| 4 | SuperchainConfig | 구현체 | src/L1/SuperchainConfig.sol:662 | Superchain 전역 설정 |
| 5 | ProtocolVersionsProxy | 프록시 | src/universal/Proxy.sol:360 | ProtocolVersions 프록시 (ERC1967) |
| 6 | ProtocolVersions | 구현체 | src/L1/ProtocolVersions.sol:816 | 프로토콜 버전 관리 |

#### Phase 2: OP Chain 프록시 배포 (deployProxies)

| # | 컨트랙트명 | 종류 | 파일:라인 | 설명 |
|---|-----------|------|----------|------|
| 7 | OptimismPortalProxy | 프록시 | src/universal/Proxy.sol:617 | OptimismPortal 프록시 |
| 8 | SystemConfigProxy | 프록시 | src/universal/Proxy.sol:617 | SystemConfig 프록시 |
| 9 | L1StandardBridgeProxy | 프록시 | src/legacy/L1ChugSplashProxy.sol:590 | L1StandardBridge 프록시 (ChugSplash) |
| 10 | L1CrossDomainMessengerProxy | 프록시 | src/legacy/ResolvedDelegateProxy.sol:603 | L1CrossDomainMessenger 프록시 (ResolvedDelegate) |
| 11 | OptimismMintableERC20FactoryProxy | 프록시 | src/universal/Proxy.sol:617 | OptimismMintableERC20Factory 프록시 |
| 12 | L1ERC721BridgeProxy | 프록시 | src/universal/Proxy.sol:617 | L1ERC721Bridge 프록시 |
| 13 | DisputeGameFactoryProxy | 프록시 | src/universal/Proxy.sol:617 | DisputeGameFactory 프록시 |
| 14 | L2OutputOracleProxy | 프록시 | src/universal/Proxy.sol:617 | L2OutputOracle 프록시 (레거시) |
| 15 | DelayedWETHProxy | 프록시 | src/universal/Proxy.sol:617 | DelayedWETH 프록시 |
| 16 | PermissionedDelayedWETHProxy | 프록시 | src/universal/Proxy.sol:617 | PermissionedDelayedWETH 프록시 |
| 17 | AnchorStateRegistryProxy | 프록시 | src/universal/Proxy.sol:617 | AnchorStateRegistry 프록시 |

#### Phase 3: OP Chain 구현체 배포 (deployImplementations)

| # | 컨트랙트명 | 종류 | 파일:라인 | 설명 | 조건 |
|---|-----------|------|----------|------|------|
| 18 | SystemConfig | 구현체 | src/L1/SystemConfig.sol:932 | 시스템 설정 저장소 | reuseDeployment=false |
| 19 | L1StandardBridge | 구현체 | src/L1/L1StandardBridge.sol:951 | L1-L2 자산 브릿지 | reuseDeployment=false |
| 20 | L1ERC721Bridge | 구현체 | src/L1/L1ERC721Bridge.sol:970 | L1-L2 ERC721 브릿지 | reuseDeployment=false |
| 21 | OptimismMintableERC20Factory | 구현체 | src/universal/OptimismMintableERC20Factory.sol:761 | 토큰 팩토리 | reuseDeployment=false |
| 22 | L1CrossDomainMessenger | 구현체 | src/L1/L1CrossDomainMessenger.sol:674 | 크로스 도메인 메시징 | reuseDeployment=false |
| 23 | L2OutputOracle | 구현체 | src/L1/L2OutputOracle.sol:738 | L2 상태 출력 오라클 (레거시) | reuseDeployment=false |
| 24 | OptimismPortal | 구현체 | src/L1/OptimismPortal.sol:692 | 메시지/출금 포털 | reuseDeployment=false |
| 25 | OptimismPortal2 | 구현체 | src/L1/OptimismPortal2.sol:711 | 메시지/출금 포털 v2 (Fault Proof 활성화) | reuseDeployment=false |
| 26 | DisputeGameFactory | 구현체 | src/dispute/DisputeGameFactory.sol:779 | 분쟁 게임 생성 팩토리 | reuseDeployment=false |
| 27 | DelayedWETH | 구현체 | src/dispute/weth/DelayedWETH.sol:794 | Wrapped ETH (지연 출금) | reuseDeployment=false |
| 28 | AnchorStateRegistry | 구현체 | src/dispute/AnchorStateRegistry.sol:921 | Fault Proof 앵커 상태 레지스트리 | reuseDeployment=false |
| 29 | PreimageOracle | 구현체 | src/cannon/PreimageOracle.sol:898 | Cannon Preimage 오라클 | 항상 배포 |
| 30 | MIPS | 구현체 | src/cannon/MIPS.sol:911 | Cannon MIPS VM 구현 | 항상 배포 |
| 31 | L1UsdcBridge | 구현체 | src/USDC/L1/tokamak-UsdcBridge/L1UsdcBridge.sol:852 | USDC 브릿지 | 항상 배포 |

#### Phase 4: 추가 배포 및 초기화

| # | 컨트랙트명 | 종류 | 파일:라인 | 설명 |
|---|-----------|------|----------|------|
| 32 | L1UsdcBridgeProxy | 프록시 | src/USDC/L1/tokamak-UsdcBridge/L1UsdcBridgeProxy.sol:867 | USDC 브릿지 프록시 |
| 33 | L2NativeToken | 구현체 | src/L1/L2NativeToken.sol:842 | L2 네이티브 토큰 (devnet만) |
| 34 | DataAvailabilityChallengeProxy | 프록시 | src/universal/Proxy.sol:644 | DAC 프록시 (Plasma 활성화 시) |
| 35 | DataAvailabilityChallenge | 구현체 | src/L1/DataAvailabilityChallenge.sol:1003 | Data Availability Challenge (Plasma 활성화 시) |

**주요 배포 패턴**:

1. **프록시 우선 배포**: 모든 프록시는 구현체 배포 전에 먼저 배포됨
2. **구현체 생성**: 각 구현체는 `create2` salt를 사용하여 결정적 배포 주소 생성
3. **의존성 관리**:
   - `OptimismPortal` → `OptimismPortalProxy` 업그레이드 (또는 `OptimismPortal2`가 Fault Proof 활성화 시)
   - `SystemConfig` → `SystemConfigProxy` 업그레이드
   - `L1StandardBridge` → `L1StandardBridgeProxy` 업그레이드
   - 기타 브릿지/팩토리 → 해당 프록시 업그레이드
4. **조건부 배포**:
   - `reuseDeployment=true`: 구현체 배포 스킵, 기존 주소 재사용
   - `useFaultProofs()=true`: `OptimismPortal2` 배포 및 초기화
   - `usePlasma()=true`: DataAvailabilityChallenge 배포 및 초기화
   - `devnet()=true`: `L2NativeToken` 배포

**초기화 순서** (initializeImplementations):

프록시 업그레이드 및 초기화는 SafeProxy를 통해 수행됩니다:

1. OptimismPortal 또는 OptimismPortal2 (Fault Proof에 따라)
2. SystemConfig
3. L1StandardBridge
4. L1ERC721Bridge
5. OptimismMintableERC20Factory
6. L1CrossDomainMessenger
7. L2OutputOracle (레거시)
8. DisputeGameFactory
9. DelayedWETH
10. PermissionedDelayedWETH
11. AnchorStateRegistry

**생성 파일**:
- `deployments/{chainId}-deploy.json`: 모든 배포된 컨트랙트의 주소 및 배포 트랜잭션 정보
  ```json
  {
    "AddressManager": "0x...",
    "ProxyAdmin": "0x...",
    "SuperchainConfigProxy": "0x...",
    "SuperchainConfig": "0x...",
    ...
    "AnchorStateRegistry": "0x..."
  }
  ```

---

## 호출 시퀀스

```
executeDeployments (Phase 3)
    ↓
[Deploy 엔티티 필터링 및 순서 정렬]
    ↓
for deployment := range [L1Contracts, AWSInfra] {
    ↓
    NewThanosSDKClient
        ↓
        ThanosStack.NewThanosStack
            ├─ 로그 초기화
            ├─ 네트워크 설정
            └─ AWS 설정 (선택)
    ↓
    DeployL1Contracts (백엔드 래퍼)
        ↓
        DTO 변환 (DeployL1ContractsRequest → DeployContractsInput)
        ↓
        sdkClient.DeployContracts (SDK 핵심)
            ├─ 1. 네트워크 검증 (LocalDevnet 제외)
            ├─ 2. 배포 상태 확인 (이미 완료 시 스킵 또는 재배포)
            ├─ 3. ChainConfiguration 자동 생성 (L1 RPC 조회)
            ├─ 4. L2 Chain ID 생성
            ├─ 5. Fault Proof 처리 (활성화 시)
            │   ├─ tokamak-thanos 클론
            │   ├─ AnchorStateRegistry 패치
            │   └─ Cannon prestate 빌드
            ├─ 6. Deploy 설정 파일 생성 (deploy-config.json)
            ├─ 7. Admin 계정 잔액 확인
            ├─ 8. 배포 상태 → InProgress 업데이트 (settings.json)
            ├─ 9. tokamak-deployer 바이너리 다운로드/캐시
            ├─ 10. runDeployContracts 실행
            │   ↓
            │   tokamak-deployer deploy-contracts \
            │     --l1-rpc <URL> \
            │     --private-key <AdminKey> \
            │     --chain-id <L2ChainID> \
            │     --out deploy-output.json
            │   ↓
            │   [Foundry를 통해 L1 컨트랙트 배포]
            │   ├─ OptimismPortal, L1StandardBridge, L1CrossDomainMessenger 등
            │   └─ deploy-output.json 생성
            ├─ 11. runGenerateGenesis 실행
            │   ↓
            │   tokamak-deployer generate-genesis \
            │     --deploy-output deploy-output.json \
            │     --config deploy-config.json \
            │     --out genesis.json
            │   ↓
            │   [Post-processing]
            │   ├─ DRB 설정 주입
            │   ├─ USDC 토큰 주입
            │   ├─ MultiTokenPaymaster 주입
            │   └─ Rollup 해시 업데이트
            └─ 12. RegisterCandidate 처리 (선택)
                ├─ Safe Wallet 설정
                └─ Candidate 등록 트랜잭션

    ↓
    statusChan ← DeploymentRunStatusSuccess
    ↓
[배포 상태 저장 (데이터베이스)]
}

↓
deploy() 함수 호출 (Phase 3)
    ├─ Chain Information 조회
    ├─ Bridge/BlockExplorer 통합 메타데이터 업데이트
    ├─ RegisterCandidate 통합 상태 업데이트
    └─ 스택 상태 → StackStatusDeployed
```

---

## 데이터 구조

### 4.1 입력 구조체

#### dtos.DeployL1ContractsRequest (백엔드 API)

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/api/dtos/`

```go
type DeployL1ContractsRequest struct {
    L1RpcUrl                    string                    `json:"l1_rpc_url" binding:"required"`
    AdminAccount                string                    `json:"admin_account" binding:"required"`      // private key
    SequencerAccount            string                    `json:"sequencer_account" binding:"required"`  // private key
    BatcherAccount              string                    `json:"batcher_account" binding:"required"`    // private key
    ProposerAccount             string                    `json:"proposer_account" binding:"required"`   // private key
    ChallengerAccount           string                    `json:"challenger_account"`                    // optional
    BatchSubmissionFrequency    int64                     `json:"batch_submission_frequency"`            // seconds
    ChallengePeriod             int64                     `json:"challenge_period"`                      // seconds
    OutputRootFrequency         int64                     `json:"output_root_frequency"`                 // seconds
    L2BlockTime                 int64                     `json:"l2_block_time"`                         // seconds
    RegisterCandidate           bool                      `json:"register_candidate"`
    RegisterCandidateParams     *RegisterCandidateParams  `json:"register_candidate_params,omitempty"`
    ReuseDeployment             bool                      `json:"reuse_deployment"`
    EnableFaultProof            bool                      `json:"enable_fault_proof"`
    Preset                      string                    `json:"preset"`                                // custom, defi, full
    FeeToken                    string                    `json:"fee_token"`                             // ETH, USDC, TON
}

type RegisterCandidateParams struct {
    Amount   *big.Int    `json:"amount"`
    Memo     string      `json:"memo"`
    NameInfo interface{} `json:"name_info"`
}
```

#### thanosStack.DeployContractsInput (SDK API)

**파일**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/`

```go
type DeployContractsInput struct {
    L1RPCurl             string
    ChainConfiguration   *types.ChainConfiguration
    Operators            *types.Operators
    RegisterCandidate    *RegisterCandidateInput
    ReuseDeployment      bool
    EnableFaultProof     bool
    Preset               string
    FeeToken             string
}

type ChainConfiguration struct {
    BatchSubmissionFrequency uint64
    ChallengePeriod          uint64
    OutputRootFrequency      uint64
    L2BlockTime              uint64
    L1BlockTime              uint64
}

type Operators struct {
    AdminPrivateKey      string
    SequencerPrivateKey  string
    BatcherPrivateKey    string
    ProposerPrivateKey   string
    ChallengerPrivateKey string
}

type RegisterCandidateInput struct {
    Amount   *big.Int
    Memo     string
    NameInfo interface{}
    UseTon   bool
}
```

### 4.2 출력 구조체

#### deploy-output.json

**위치**: `{deploymentPath}/deploy-output.json`

**샘플** (Foundry 배포 결과):
```json
{
  "1155111": {
    "OptimismPortal": "0x...",
    "OptimismPortalProxy": "0x...",
    "L1StandardBridge": "0x...",
    "L1StandardBridgeProxy": "0x...",
    "L1CrossDomainMessenger": "0x...",
    "L1CrossDomainMessengerProxy": "0x...",
    "L1ERC721Bridge": "0x...",
    "L1ERC721BridgeProxy": "0x...",
    "SystemConfig": "0x...",
    "SystemConfigProxy": "0x...",
    "L2OutputOracle": "0x...",
    "L2OutputOracleProxy": "0x...",
    "AddressManager": "0x...",
    "ProxyAdmin": "0x...",
    "AnchorStateRegistry": "0x...",
    "DelayedWETHProxy": "0x...",
    "WETH98": "0x...",
    "PermissionedDisputeGame": "0x...",
    "PermissionedDisputeGameProxy": "0x...",
    "FaultDisputeGame": "0x...",
    "FaultDisputeGameProxy": "0x...",
    "Cannon": "0x...",
    "CannonProxy": "0x...",
    "MIPS": "0x...",
    "MIPSProxy": "0x...",
    "transactionHash": "0x..."  // 배포 트랜잭션
  }
}
```

#### rollup.json

**위치**: `{deploymentPath}/rollup.json`

**샘플**:
```json
{
  "rollup": {
    "version": 0,
    "basefeeScalar": 1368,
    "blobbasefeeScalar": 810949,
    "maxSequencerDrift": 600,
    "sequencerWindowSize": 3600,
    "channelTimeout": 300,
    "l1BlockTime": 12,
    "blockTime": 12,
    "l2OutputOracleSubmissionInterval": 75,
    "maxOutputRootSize": 134217728,
    "faultGameAbsolutePrestatusSize": 33554432,
    "faultGameMaxDepth": 73,
    "faultGameClockExtension": 10800,
    "faultGameMaxClockDuration": 302400,
    "faultGameGenesisBlock": 0,
    "faultGameGenesisOutputRoot": "0x...",
    "faultGameSplitDepth": 30,
    "faultGameWithdrawalDelay": 604800
  },
  "genesis": {
    "l1StartingBlockTag": "...",
    "systemConfig": {
      "batcherAddr": "0x...",
      "overhead": 188,
      "scalar": 0.684,
      "gasLimit": 30000000
    },
    "l2Time": ...,
    "l2Block": {
      "number": 0,
      "hash": "0x..."
    }
  },
  "blockOracle": {
    "sources": [
      {
        "name": "optimism",
        "origin": "0x...",
        "blockOffset": 0
      }
    ]
  },
  "l1ChainID": 11155111,
  "l2ChainID": 123456,
  "depositContractAddress": "0x...",
  "l1StandardBridgeAddress": "0x...",
  "l1CrossDomainMessengerAddress": "0x...",
  "l2ToL1MessagePasser": "0x4200000000000000000000000000000000000016",
  "l2StandardBridgeAddress": "0x4200000000000000000000000000000000000010",
  "l2CrossDomainMessengerAddress": "0x4200000000000000000000000000000000000007",
  "l2ERC721BridgeAddress": "0x4200000000000000000000000000000000000014",
  "baseFeeVaultRecipient": "0x4200000000000000000000000000000000000019",
  "l1FeeVaultRecipient": "0x420000000000000000000000000000000000001a",
  "sequencerFeeVaultRecipient": "0x4200000000000000000000000000000000000011",
  "optimismBaseFeeRecipient": "0x420000000000000000000000000000000000001f",
  "optimismL1FeeRecipient": "0x4200000000000000000000000000000000000020",
  "baseFeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "l1FeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "sequencerFeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "baseFeeVaultWithdrawalNetwork": 0,
  "l1FeeVaultWithdrawalNetwork": 0,
  "sequencerFeeVaultWithdrawalNetwork": 1
}
```

#### genesis.json

**위치**: `{deploymentPath}/genesis.json`

**샘플**:
```json
{
  "config": {
    "chainId": 123456,
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
    "mergeNetsplitBlock": 0,
    "shanghaiBlock": 0,
    "cancunBlock": 0,
    "optimismSupernetwork": {
      "ejectFinalizedBlocks": false
    }
  },
  "difficulty": "0x1",
  "gasLimit": "0x1c9c380",
  "alloc": {
    "0x...": {
      "balance": "0x...",
      "code": "0x..."
    }
  },
  "number": "0x0",
  "gasUsed": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "baseFeePerGas": "0x3b9aca00"
}
```

---

## 에러 시나리오

### E4.1 네트워크 검증 실패

**원인**: LocalDevnet에 대해 L1 배포 시도

**코드 위치**: `deploy_contracts.go` 라인 24-31

```go
if t.network == constants.LocalDevnet {
    return fmt.Errorf("network %s does not require contract deployment", constants.LocalDevnet)
}
```

**증상**: 요청 즉시 실패, 명확한 에러 메시지

**복구**:
- 네트워크 파라미터 확인 (testnet/mainnet 선택)
- 로컬 배포는 L2 인프라 배포 단계에서만 실행

---

### E4.2 L1 RPC 연결 실패

**원인**: 잘못된 RPC URL, RPC 서비스 다운, 네트워크 불안정

**코드 위치**: `deploy_contracts.go` 라인 39-43, 93-96

```go
l1Client, err := ethclient.DialContext(ctx, deployContractsConfig.L1RPCurl)
if err != nil {
    t.logger.Error("Failed to get the L1 client", "err", err)
    return err
}
```

**증상**:
- 배포 시작 전 즉시 실패
- 에러 메시지: `"Failed to get the L1 client"`

**복구**:
1. RPC URL 유효성 확인 (https 프로토콜)
2. RPC 서비스 상태 확인 (예: Infura, Alchemy 대시보드)
3. 네트워크 연결 확인
4. Rate limit 초과 시 다른 RPC 제공자 사용

---

### E4.3 L1 체인 ID 조회 실패

**원인**: RPC 연결 불안정, RPC 메서드 미지원

**코드 위치**: `deploy_contracts.go` 라인 45-49

```go
l1ChainID, err := l1Client.ChainID(ctx)
if err != nil {
    t.logger.Error("Failed to get the L1 ChainID", "err", err)
    return err
}
```

**증상**:
- 에러 메시지: `"Failed to get the L1 ChainID"`
- Admin 계정 잔액 확인 전 실패

**복구**:
- RPC URL이 기본 메서드(eth_chainId)를 지원하는지 확인
- 다른 RPC 제공자로 재시도

---

### E4.4 관리자 계정 잔액 부족

**원인**: Admin 계정 ETH 잔액이 배포 비용 추정치보다 적음

**코드 위치**: `deploy_contracts.go` 라인 313-318

```go
if balance.Cmp(estimatedCost) < 0 {
    return fmt.Errorf("admin account balance (%.4f ETH) is less than estimated deployment cost (%.4f ETH)",
        utils.WeiToEther(balance), utils.WeiToEther(estimatedCost))
}
```

**증상**:
- 명확한 에러 메시지: `"Insufficient balance for deployment"`
- 현재 잔액과 필요 금액 표시

**복구**:
1. Admin 계정에 ETH 입금 (Sepolia Faucet 사용)
2. 배포 비용 추정 재계산:
   - 현재 가스비 확인
   - 큰 배포 비용 예상 시 gas price 낮은 시간대 재시도

---

### E4.5 배포 상태 파일 쓰기 실패

**원인**: 디스크 권한 부족, 경로 불존재, 디스크 용량 부족

**코드 위치**: `deploy_contracts.go` 라인 332-336

```go
err = t.deployConfig.WriteToJSONFile(t.deploymentPath)
if err != nil {
    t.logger.Error("Failed to write settings file", "err", err)
    return err
}
```

**증상**:
- 에러 메시지: `"Failed to write settings file"`
- 배포 상태 저장 불가능

**복구**:
1. 배포 경로 권한 확인: `chmod 755 {deploymentPath}`
2. 상위 디렉토리 생성 필요 시: `mkdir -p {deploymentPath}`
3. 디스크 공간 확인

---

### E4.6 tokamak-deployer 바이너리 다운로드 실패

**원인**: GitHub 다운로드 실패, 네트워크 제한, 잘못된 바이너리 버전

**코드 위치**: `deployer_binary.go` 라인 42-44

```go
if err := downloadFile(downloadURL, binaryPath); err != nil {
    return "", fmt.Errorf("failed to download tokamak-deployer %s: %w", version, err)
}
```

**증상**:
- 에러 메시지: `"failed to download tokamak-deployer v1.0.0"`
- HTTP 상태 코드 제시 (404, 403 등)

**복구**:
1. GitHub Releases 페이지 확인:
   - `https://github.com/tokamak-network/tokamak-thanos/releases/tag/v1.0.0`
2. 바이너리 존재 확인 (OS/아키텍처 매칭)
3. 캐시 삭제 후 재시도:
   ```bash
   rm -rf ~/.trh/bin/tokamak-deployer-v1.0.0
   ```
4. 네트워크 방화벽/프록시 확인

---

### E4.7 Foundry 배포 트랜잭션 실패

**원인**: 가스 부족, Nonce 충돌, 계정 잠금, L1 혼잡

**코드 위치**: `deployer_binary.go` 라인 120-126

```go
cmd := exec.CommandContext(ctx, binaryPath, args...)
if err := cmd.Run(); err != nil {
    return fmt.Errorf("tokamak-deployer deploy-contracts: %w", err)
}
```

**증상**:
- 에러 메시지: `"tokamak-deployer deploy-contracts: <binary error>"`
- Binary stdout/stderr에 Foundry 에러 메시지 포함

**일반적인 원인들**:

#### E4.7.1 가스 부족
```
Error: Not enough ETH to send transaction
```
**복구**: Admin 계정에 추가 ETH 입금

#### E4.7.2 Nonce 충돌
```
Error: nonce has already been used
```
**원인**: 동일 계정에서 동시에 다른 트랜잭션 전송

**복구**:
- Admin 계정 다른 pending 트랜잭션 확인
- 트랜잭션 완료 후 배포 재시도

#### E4.7.3 Tx 실패 (Revert)
```
Error: Transaction failed with status 0
```
**복구**:
- 배포 파라미터 검증 (ChainConfiguration, Operators)
- L1 상태 확인 (스마트 컨트랙트 버그 등)

---

### E4.8 Resume 중 상태 불일치

**원인**: 부분 배포 후 설정 파일 손상, 수동 개입

**코드 위치**: `deploy_contracts.go` 라인 99-129

```go
if t.deployConfig.DeployContractState != nil {
    switch t.deployConfig.DeployContractState.Status {
    case types.DeployContractStatusInProgress:
        isResume = true
    }
}
```

**증상**:
- `settings.json`에 `Status: InProgress` 존재
- `deploy-output.json` 부분적으로만 존재 또는 손상

**복구**:
1. 배포 상태 파일 확인:
   ```bash
   cat {deploymentPath}/settings.json | jq .deploy_contract_state
   ```
2. 부분 파일 정리:
   ```bash
   rm {deploymentPath}/deploy-output.json
   ```
3. 상태 초기화:
   ```bash
   # settings.json의 deploy_contract_state 제거 또는 Status를 Completed로 변경
   ```
4. 배포 재시도

---

### E4.9 Canon Prestate 빌드 실패 (Fault Proof)

**원인**: Fault Proof 활성화 시 tokamak-thanos 클론 또는 컴파일 실패

**코드 위치**: `deploy_contracts.go` 라인 224-259

```go
if deployContractsConfig.EnableFaultProof {
    err = t.cloneSourcecode(ctx, "tokamak-thanos", ...)
    if patchErr := patchAnchorStateRegistry(tokamakThanosDir); patchErr != nil {
        return fmt.Errorf("failed to patch AnchorStateRegistry.sol: %w", patchErr)
    }
    if buildErr := buildCannonPrestate(gctx, t.logger, tokamakThanosDir); buildErr != nil {
        return fmt.Errorf("failed to build cannon prestate: %w", buildErr)
    }
}
```

**증상**:
- 에러 메시지: `"failed to patch AnchorStateRegistry.sol"` 또는 `"failed to build cannon prestate"`
- 배포 30분 이상 소요 가능

**복구**:
1. Fault Proof 미활성화로 재시도 (개발 환경)
2. 필수 도구 확인:
   - Go 1.21+
   - Rust (cannon 빌드용)
   - Foundry (forge/cast)
3. 토카막 저장소 상태 확인:
   ```bash
   cd {deploymentPath}/tokamak-thanos
   git status
   git log --oneline -5
   ```

---

### E4.10 Genesis/Rollup 파일 생성 실패

**원인**: `deploy-output.json` 손상, 바이너리 버전 불일치, 설정 파일 누락

**코드 위치**: `deploy_contracts.go` 라인 369-380

```go
if err = runGenerateGenesis(ctx, binaryPath, genesisOpts{
    DeployOutputPath: deployOutputPath,
    ConfigPath:       deployConfigFilePath,
    OutPath:          genesisPath,
}); err != nil {
    t.logger.Error("❌ Failed to generate rollup and genesis files!")
    return err
}
```

**증상**:
- 에러 메시지: `"Failed to generate rollup and genesis files"`
- `deploy-output.json` 또는 `deploy-config.json` 검증 실패

**복구**:
1. 입력 파일 검증:
   ```bash
   jq . {deploymentPath}/deploy-output.json
   jq . {deploymentPath}/deploy-config.json
   ```
2. 바이너리 버전 확인:
   ```bash
   {binaryPath} version
   ```
3. 필요 시 바이너리 재다운로드:
   ```bash
   rm ~/.trh/bin/tokamak-deployer-v1.0.0
   ```

---

### E4.11 RegisterCandidate 실패

**원인**: Safe Wallet 설정 실패, 토카막 네트워크 접근 불가, 보증금 부족

**코드 위치**: `deploy_contracts.go` 라인 389-404

```go
if t.registerCandidate {
    if err := t.setupSafeWallet(ctx, t.deploymentPath); err != nil {
        return err
    }
    if verifyRegisterError := t.verifyRegisterCandidates(ctx, registerCandidate); verifyRegisterError != nil {
        return fmt.Errorf("candidate registration failed: %v", verifyRegisterError)
    }
}
```

**증상**:
- 에러 메시지: `"candidate registration failed"` + 상세 원인
- L1 배포는 성공했으나 Candidate 등록 실패

**일반적 원인**:

#### E4.11.1 Safe Wallet 설정 실패
```
Error: Failed to initialize Safe Wallet
```
**복구**:
- 토카막 Safe 주소 확인
- Safe 공개 RPC 접근성 확인

#### E4.11.2 보증금 부족
```
Error: Insufficient TON balance for registration
```
**복구**:
- Admin 계정 TON 잔액 확인
- `RegisterCandidateParams.Amount` 감소

#### E4.11.3 Candidate 검증 실패
```
Error: Candidate verification failed: <reason>
```
**복구**:
- Candidate 메타데이터 유효성 확인
- Candidate 이미 등록되었는지 확인

---

### E4.12 Context Cancelled (배포 중단)

**원인**: 사용자 수동 중단, 상위 컨텍스트 타임아웃

**코드 위치**: `deployer_binary.go` 라인 120-126

```go
cmd := exec.CommandContext(ctx, binaryPath, args...)
if err := cmd.Run(); err != nil {
    if errors.Is(err, context.Canceled) {
        t.logger.Warn("Deployment canceled")
        return err
    }
    return fmt.Errorf("deployment failed: %w", err)
}
```

**증상**:
- 배포 중 진행 상황 멈춤
- 에러 메시지: `"context canceled"`
- `deploy-output.json` 부분적 존재 가능

**복구**:
1. Resume 기능 활용:
   - `settings.json`에 `Status: InProgress` 유지
   - 재배포 시 이전 진행 상황에서 재개
2. 부분 배포 파일 정리 (필요 시):
   ```bash
   rm {deploymentPath}/deploy-output.json
   ```

---

## 알려진 함정 및 개선 포인트

### F4.1 Operator 계정과 Candidate 등록 계정 혼동

**함정**: Admin 계정이 배포와 Candidate 등록 모두에 사용됨

**현재 동작**:
- Admin 계정: L1 배포 수행 및 초기화 + Candidate 등록
- 보증금: Admin 계정에서 차감

**위험성**:
- Admin 계정 개인키 노출 시 배포 및 등록 모두 손상
- 다중 시그니처 지갑 불가능

**개선 제안**:
```go
type Operators struct {
    AdminPrivateKey      string  // 배포용
    CandidateOperatorKey string  // Candidate 등록용 (분리)
    SequencerPrivateKey  string
    BatcherPrivateKey    string
    ProposerPrivateKey   string
    ChallengerPrivateKey string
}
```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 195-199

---

### F4.2 Resume 시 상태 파일 검증 부재

**함정**: `settings.json`의 `DeployContractState.Status`가 손상되면 Resume 메커니즘이 작동 불가능

**현재 동작**:
```go
if t.deployConfig.DeployContractState != nil {
    switch t.deployConfig.DeployContractState.Status {
    case types.DeployContractStatusCompleted:
        return nil  // 완료 가정
    case types.DeployContractStatusInProgress:
        isResume = true  // Resume 활성화
    }
}
```

**문제점**:
- 상태 값이 `0` (uninitialized)이면 Resume 미활성화
- 사용자가 수동으로 파일을 편집한 경우 불일치 가능
- `deploy-output.json` 존재 여부를 검증하지 않음

**개선 제안**:
```go
// Resume 여부를 상태뿐만 아니라 파일 존재로도 판단
if _, err := os.Stat(deployOutputPath); err == nil {
    // deploy-output.json 존재 → Resume 가능
    isResume = true
}
```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 99-129

---

### F4.3 배포 비용 추정 부정확

**함정**: 예상 배포 비용이 실제 비용과 크게 차이날 수 있음

**현재 동작**:
```go
// 추정값 (하드코딩)
estimatedCost := new(big.Int).Mul(gasPriceWei, estimatedDeployContracts)
estimatedCost.Mul(estimatedCost, big.NewInt(2))  // 2배 여유
```

**문제점**:
- `estimatedDeployContracts` 상수값은 기존 배포 기반
- Fault Proof, Register Candidate 등 옵션별 비용 차이 미반영
- L1 혼잡도에 따라 실제 비용 크게 변동
- 2배 여유는 경우에 따라 부족할 수 있음

**개선 제안**:
```go
// 옵션별 동적 비용 계산
baseGasEstimate := uint64(5000000)  // 기본 배포
if req.EnableFaultProof {
    baseGasEstimate += 1000000      // Fault Proof 추가
}
if req.RegisterCandidate {
    baseGasEstimate += 500000       // Candidate 등록 추가
}
estimatedCost := new(big.Int).Mul(gasPriceWei, new(big.Int).SetUint64(baseGasEstimate))
```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 299-310

---

### F4.4 Fault Proof 활성화 시 긴 처리 시간

**함정**: Fault Proof 활성화 시 Cannon prestate 빌드로 30분 이상 소요

**현재 동작**:
```go
if deployContractsConfig.EnableFaultProof {
    // 1. tokamak-thanos 클론 (5-10분, 네트워크 환경에 따라 변동)
    // 2. AnchorStateRegistry 패치 (1초)
    // 3. Cannon prestate 빌드 (20-30분, CPU 처리량 따라 변동)
}
```

**문제점**:
- 배포 진행 상황 표시 부족
- 중단/재개 메커니즘 미지원 (빌드 중 ctx cancel 시 부분 캐시 남을 수 있음)
- 캐시 재사용 안 함 (매번 새로 빌드)

**개선 제안**:
```go
// 캐시 위치: ~/.trh/cannon-prestate/{version}/prestate.json
prestateCache := filepath.Join(homeDir, ".trh", "cannon-prestate", "v1", "prestate.json")
if _, err := os.Stat(prestateCache); err == nil {
    // 캐시된 prestate 사용
    prestateHash, _ := readPrestateHash(prestateCache)
} else {
    // 새로 빌드
    buildCannonPrestate(ctx, ...)
}
```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 224-259

---

### F4.5 Private Key 평문 저장

**함정**: Admin 및 Operator 개인키가 `settings.json`에 평문으로 저장됨

**현재 동작**:
```go
t.deployConfig.AdminPrivateKey = operators.AdminPrivateKey
t.deployConfig.SequencerPrivateKey = operators.SequencerPrivateKey
// ...
err = t.deployConfig.WriteToJSONFile(t.deploymentPath)
```

**문제점**:
- `{deploymentPath}/settings.json`에 개인키 평문 저장
- 파일 권한 설정 부재 (기본 644)
- 로그에 개인키 노출 위험

**개선 제안**:
1. **파일 권한 설정**:
   ```go
   err = os.WriteFile(jsonPath, jsonData, 0600)  // 소유자만 읽기/쓰기
   ```

2. **개인키 암호화**:
   ```go
   type EncryptedConfig struct {
       EncryptedAdminKey string `json:"encrypted_admin_key"`  // AES-256 암호화
       Salt              string `json:"salt"`
   }
   ```

3. **로그 필터링**:
   ```go
   // 로그에서 개인키 마스킹
   maskedKey := key[:4] + "***" + key[len(key)-4:]
   logger.Info("Using admin key: " + maskedKey)
   ```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 264-273

---

### F4.6 설정 검증 시 반복 작업

**함정**: ChainConfiguration 검증이 배포 초기가 아닌 SDK 메서드에서 수행되어 두 번 검증

**현재 동작**:
1. 백엔드 HTTP 핸들러: 기본 검증 (필수 필드)
2. 서비스 레이어: 없음
3. SDK 메서드: `validate()` 호출 (라인 66)

**문제점**:
- 배포 파라미터 유효성이 늦게 확인됨
- 프론트엔드와 백엔드 간 검증 중복 가능

**개선 제안**:
```go
// 서비스 레이어에서 조기 검증
func (s *ThanosStackDeploymentService) ValidateL1DeploymentConfig(req *dtos.DeployL1ContractsRequest) error {
    if req.L1RpcUrl == "" {
        return fmt.Errorf("L1 RPC URL required")
    }
    if req.AdminAccount == "" {
        return fmt.Errorf("Admin account private key required")
    }
    // L1 RPC 연결 테스트
    l1Client, err := ethclient.Dial(req.L1RpcUrl)
    if err != nil {
        return fmt.Errorf("failed to connect to L1 RPC: %w", err)
    }
    return nil
}
```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go` 라인 66-70

---

### F4.7 토카막 바이너리 버전 핀포인트의 유연성 부족

**함�**: `TokamakDeployerVersion = "v1.0.0"` 하드코딩으로 버전 업그레이드가 어려움

**현재 동작**:
```go
const TokamakDeployerVersion = "v1.0.0"
binaryName := fmt.Sprintf("tokamak-deployer-%s", version)
```

**문제점**:
- 새 버전 배포 시 코드 변경 필수
- 다양한 버전 동시 지원 불가능
- 롤백 시 코드 재컴파일 필요

**개선 제안**:
```go
// 환경 변수로 버전 지정
var TokamakDeployerVersion = func() string {
    if v := os.Getenv("TOKAMAK_DEPLOYER_VERSION"); v != "" {
        return v
    }
    return "v1.0.0"
}()

// 또는 배포 요청 시 버전 지정
type DeployContractsInput struct {
    // ...
    DeployerVersion string  // optional, defaults to v1.0.0
}
```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deployer_binary.go` 라인 15

---

### F4.8 배포 진행 상황 추적 부족

**함정**: 배포 중 사용자에게 진행 상황 피드백이 거의 없음

**현재 동작**:
```go
// 진행 상황 표시 전무 (Binary가 stdout/stderr로 출력)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
if err := cmd.Run(); err != nil {
    return fmt.Errorf("deployment failed: %w", err)
}
```

**문제점**:
- 30분 이상 소요되는 Fault Proof 배포 시 진행 상황 불명
- 배포 중단 여부 판단 어려움
- 로그 수집 불충분

**개선 제안**:
```go
// 바이너리 stdout 파싱하여 진행 상황 추출
type DeploymentProgress struct {
    Step     string  // "deploying_portal", "deploying_bridges", ...
    Progress int     // 0-100
    Duration time.Duration
}

func parseDeployerOutput(line string) *DeploymentProgress {
    // Binary가 "Step X/Y" 형식 출력하면 파싱
    if strings.Contains(line, "Deploying") {
        return &DeploymentProgress{
            Step:     extractStep(line),
            Progress: extractProgress(line),
        }
    }
    return nil
}
```

**참고 코드**: `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deployer_binary.go` 라인 120-128

---

## 참고 문헌

### 관련 파일
- `/Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go` - 배포 오케스트레이션
- `/Users/theo/workspace_tokamak/trh-backend/pkg/stacks/thanos/thanos_stack.go` - SDK 래퍼
- `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go` - SDK 핵심 로직
- `/Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deployer_binary.go` - 바이너리 관리

### 외부 참고
- Foundry: https://github.com/foundry-rs/foundry
- Optimism Bedrock: https://github.com/ethereum-optimism/optimism
- Tokamak Network: https://github.com/tokamak-network/

---

**작성일**: 2026-04-16  
**Phase**: 4 (L1 컨트랙트 배포)  
**문서 버전**: 1.0
