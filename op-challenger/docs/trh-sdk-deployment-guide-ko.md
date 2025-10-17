# TRH SDK 배포 동작 상세 가이드

## 목차
1. [TRH SDK 개요](#trh-sdk-개요)
2. [SDK 설치 과정](#sdk-설치-과정)
3. [로컬 Devnet 배포](#로컬-devnet-배포)
4. [Testnet/Mainnet 배포](#testnetmainnet-배포)
5. [플러그인 시스템](#플러그인-시스템)
6. [모니터링 및 관리](#모니터링-및-관리)
7. [배포 플로우 다이어그램](#배포-플로우-다이어그램)

---

## TRH SDK 개요

**Tokamak Rollup Hub (TRH) SDK**는 Ethereum 네트워크 위에 커스터마이징된 Layer 2 Rollup을 빠르게 배포할 수 있게 해주는 CLI 도구입니다.

### 주요 특징

- **원클릭 배포**: 로컬 devnet부터 프로덕션 환경까지 단일 명령어로 배포
- **자동화된 인프라**: AWS EKS + Terraform을 통한 자동 인프라 프로비저닝
- **플러그인 시스템**: Bridge, Block Explorer, Monitoring 등 선택적 설치
- **다중 네트워크 지원**: Local Devnet, Testnet(Sepolia), Mainnet
- **Thanos Stack 기반**: Optimism Bedrock 기반의 Tokamak 커스터마이즈 스택

### 지원 스택

- **Thanos Stack**: Tokamak Network의 Optimism 포크
  - Fault Proof 지원
  - Plasma DA 모드 지원
  - 커스터마이즈된 경제 모델

### 기술 스택

```
┌─────────────────────────────────────────┐
│           TRH SDK (Go CLI)              │
├─────────────────────────────────────────┤
│  ┌────────────┐  ┌────────────────────┐ │
│  │  Commands  │  │  Thanos Stack Impl │ │
│  └─────┬──────┘  └─────┬──────────────┘ │
│        │                │                 │
├────────┼────────────────┼────────────────┤
│  ┌─────▼─────┐    ┌─────▼──────┐        │
│  │ Docker    │    │ Terraform  │        │
│  │ Compose   │    │ (AWS EKS)  │        │
│  └───────────┘    └────────────┘        │
├─────────────────────────────────────────┤
│       Tokamak Thanos Stack              │
│  (L1 Contracts + L2 Nodes + Services)   │
└─────────────────────────────────────────┘
```

---

## SDK 설치 과정

### setup.sh 스크립트 분석

`setup.sh`는 SDK 설치를 자동화하는 스크립트입니다.

#### 설치 모드

```bash
# 1. 최신 릴리스 버전 설치 (기본)
./setup.sh

# 2. 특정 커밋 해시 설치
./setup.sh -c <commit_hash>

# 3. main 브랜치 최신 버전 설치
./setup.sh --latest
```

#### 설치 프로세스

```
1. 환경 감지
   ├─ OS 타입 (Darwin/Linux)
   ├─ 아키텍처 (amd64/arm64)
   └─ 쉘 타입 (zsh/bash)

2. 패키지 설치
   ├─ Homebrew (macOS)
   ├─ Git
   ├─ Xcode Command Line Tools
   └─ Build Essential (Linux)

3. Go 1.22.6 설치
   ├─ 기존 버전 확인
   ├─ 버전 불일치 시 재설치
   └─ PATH 설정

4. TRH SDK CLI 설치
   ├─ go install github.com/tokamak-network/trh-sdk@<version>
   └─ $GOPATH/bin에 바이너리 설치

5. 추가 의존성 설치
   └─ install-all-packages.sh 실행
      ├─ Node.js 20.16.0
      ├─ pnpm
      ├─ Foundry (forge, cast, anvil)
      ├─ Terraform
      ├─ Helm
      ├─ kubectl
      └─ AWS CLI v2
```

#### PATH 설정

```bash
# .zshrc 또는 .bashrc에 추가됨
export PATH="$PATH:/usr/local/go/bin"
export PATH="$HOME/go/bin:$PATH"
```

#### 설치 확인

```bash
# SDK 버전 확인
trh-sdk version

# 의존성 확인
go version           # 1.22.6
node --version       # 20.16.0
terraform --version  # >= 1.0
helm version         # >= 3.0
kubectl version      # >= 1.20
aws --version        # >= 2.0
forge --version      # Foundry
```

---

## 로컬 Devnet 배포

로컬 환경에서 전체 스택을 Docker Compose로 실행합니다.

### 배포 명령어

```bash
trh-sdk deploy
```

### 내부 동작 과정

#### 1. 설정 파일 확인 (`commands/deploy.go`)

```go
// settings.json 파일 읽기
config, err := utils.ReadConfigFromJSONFile(deploymentPath)

// 없으면 devnet으로 간주
if config == nil {
    network = constants.LocalDevnet
    stack = constants.ThanosStack
}
```

#### 2. Thanos Stack 초기화

```go
thanosStack, err := thanos.NewThanosStack(
    ctx,
    logger,
    network,      // "local-devnet"
    true,         // usePromptInput
    deploymentPath,
    nil,          // awsConfig (devnet에선 불필요)
)
```

#### 3. 로컬 Devnet 배포 (`thanos/deploy_chain.go`)

```go
func (t *ThanosStack) deployLocalDevnet(ctx context.Context) error {
    // STEP 1: tokamak-thanos 저장소 클론
    err := t.cloneSourcecode(
        ctx,
        "tokamak-thanos",
        "https://github.com/tokamak-network/tokamak-thanos.git",
    )

    // STEP 2: devnet 시작
    cmd := fmt.Sprintf(
        "cd %s/tokamak-thanos && export DEVNET_L2OO=true && make devnet-up",
        t.deploymentPath,
    )
    err = utils.ExecuteCommandStream(ctx, t.logger, "bash", "-c", cmd)

    return nil
}
```

#### 4. make devnet-up 실행

이 명령어는 tokamak-thanos의 Makefile을 사용합니다:

```makefile
# tokamak-thanos/Makefile
devnet-up: pre-devnet
    PYTHONPATH=./bedrock-devnet python3 ./bedrock-devnet/main.py --monorepo-dir=.

pre-devnet: submodules
    # geth 설치 확인
    # cannon-prestate 생성
```

Python 스크립트 (`bedrock-devnet/main.py`)가:
1. L1/L2 제네시스 생성
2. 스마트 컨트랙트 배포
3. `.devnet/` 설정 파일 생성
4. Docker Compose 시작 (`ops-bedrock/docker-compose.yml`)

#### 5. Docker Compose 서비스 시작

```yaml
# ops-bedrock/docker-compose.yml (로컬 Devnet용)
services:
  l1:              # L1 geth
  l2:              # L2 geth
  op-node:         # Rollup node
  op-proposer:     # Output proposer
  op-batcher:      # Batch submitter
  op-challenger:   # Fault proof challenger
  da-server:       # Data availability (Plasma 모드, 로컬 전용)
  artifact-server: # Config file server

# 주의: AWS/Kubernetes 배포에는 da-server가 포함되지 않습니다.
# 프로덕션 환경은 Plasma 대신 L1을 직접 사용하여 데이터 가용성을 보장합니다.
```

#### 6. 배포 완료 출력

```
Container ops-bedrock-l1-1  Running
Container ops-bedrock-l2-1  Running
Container ops-bedrock-op-node-1  Running
Container ops-bedrock-op-challenger-1  Started
✅ Devnet up!
```

### 로컬 vs 프로덕션 배포 차이점

#### 로컬 Devnet (Docker Compose)

**8개 서비스 구성:**
- l1, l2 (geth)
- op-node, op-proposer, op-batcher, op-challenger
- **da-server** (Plasma Data Availability)
- artifact-server

**특징:**
- Plasma 모드 선택 가능 (`PLASMA_ENABLED` 환경 변수)
- da-server가 L2 트랜잭션 데이터를 로컬에 저장
- 개발 및 테스트 용도

#### AWS/Kubernetes 프로덕션 배포

**5개 주요 서비스 구성:**
- op-geth (L2 StatefulSet)
- op-node (StatefulSet)
- op-batcher, op-proposer, op-challenger (Deployments)

**특징:**
- ❌ **da-server 없음**
- ✅ L1 Ethereum을 직접 사용하여 데이터 가용성 보장
- Plasma 모드 미사용 (프로덕션 보안 고려)
- 모든 트랜잭션 데이터가 L1에 직접 게시됨

**데이터 가용성 비교:**

| 모드 | Local Devnet | AWS Production |
|-----|-------------|----------------|
| **DA 방식** | Plasma (da-server) | L1 Direct |
| **비용** | 무료 (로컬) | L1 Gas 비용 |
| **보안** | 낮음 (테스트용) | 높음 (L1 보장) |
| **da-server** | ✅ 필요 | ❌ 불필요 |

### 로그 확인

```bash
# 전체 로그 (로컬 devnet)
cd {deploymentPath}/tokamak-thanos/ops-bedrock
docker compose logs -f

# 특정 서비스
docker compose logs -f op-challenger
```

### 제거

```bash
trh-sdk destroy
```

내부적으로 실행되는 명령어:

```bash
cd {deploymentPath}/tokamak-thanos
make devnet-clean
```

---

## Testnet/Mainnet 배포

프로덕션 환경은 AWS EKS에 배포됩니다.

### 전제 조건

```json
{
  "L1 RPC URL": "https://sepolia.infura.io/v3/YOUR_KEY",
  "L1 Beacon URL": "https://beacon.sepolia.io",
  "AWS Credentials": {
    "access_key": "AKIA...",
    "secret_key": "...",
    "region": "ap-northeast-2"
  },
  "Operators": {
    "admin_private_key": "0x...",
    "sequencer_private_key": "0x...",
    "batcher_private_key": "0x...",
    "proposer_private_key": "0x..."
  }
}
```

### 2단계 배포 프로세스

```
Phase 1: L1 컨트랙트 배포
   └─> trh-sdk deploy-contracts

Phase 2: 인프라 및 L2 체인 배포
   └─> trh-sdk deploy
```

---

## Phase 1: L1 컨트랙트 배포

### 명령어

```bash
trh-sdk deploy-contracts --network testnet --stack thanos
```

### 상세 동작 (`thanos/deploy_contracts.go`)

#### STEP 1: 사용자 입력 수집

```go
deployContractsConfig, err := thanos.InputDeployContracts(ctx)
```

**수집 정보**:
```go
type DeployContractsInput struct {
    L1RPCurl            string
    L1BeaconURL         string
    Operators           *OperatorsConfig
    ChainConfiguration  *ChainConfiguration
    RegisterCandidate   *RegisterCandidateInput
}

type OperatorsConfig struct {
    AdminPrivateKey      string
    SequencerPrivateKey  string
    BatcherPrivateKey    string
    ProposerPrivateKey   string
    ChallengerPrivateKey string
}

type ChainConfiguration struct {
    BatchSubmissionFrequency uint64  // 배치 제출 주기
    ChallengePeriod          uint64  // 챌린지 기간 (7일)
    OutputRootFrequency      uint64  // Output root 제출 주기
    L2BlockTime              uint64  // L2 블록 타임 (2초)
    L1BlockTime              uint64  // L1 블록 타임
}
```

#### STEP 2: L2 Chain ID 생성

```go
l2ChainID, err := utils.GenerateL2ChainId()
// 랜덤 Chain ID 생성 (범위: 10000-99999999)
```

#### STEP 3: 잔액 확인

```go
// Admin 계정 잔액 조회
balance, err := l1Client.BalanceAt(ctx, adminAccount, nil)

// Gas Price 조회
gasPriceWei, err := l1Client.SuggestGasPrice(ctx)

// 예상 비용 계산 (실제 비용의 2배로 여유 확보)
estimatedCost := gasPriceWei * estimatedDeployContracts * 2

// 잔액 충분성 검증
if balance < estimatedCost {
    return fmt.Errorf("insufficient balance")
}
```

#### STEP 4: 소스코드 클론

```go
err := t.cloneSourcecode(
    ctx,
    "tokamak-thanos",
    "https://github.com/tokamak-network/tokamak-thanos.git",
)
```

#### STEP 5: deploy-config.json 생성

```go
deployConfigFilePath := fmt.Sprintf(
    "%s/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy-config.json",
    t.deploymentPath,
)

err = makeDeployContractConfigJsonFile(
    ctx,
    l1Client,
    operators,
    deployContractsTemplate,
    deployConfigFilePath,
)
```

**deploy-config.json 구조**:
```json
{
  "l1ChainID": 11155111,
  "l2ChainID": 12345678,
  "l2BlockTime": 2,
  "finalizationPeriodSeconds": 604800,
  "l2OutputOracleSubmissionInterval": 120,
  "maxSequencerDrift": 600,
  "sequencerWindowSize": 3600,
  "channelTimeout": 300,
  "batchInboxAddress": "0xff00000000000000000000000000000012345678",
  "batchSenderAddress": "0x...",
  "l2OutputOracleProposer": "0x...",
  "l2OutputOracleChallenger": "0x...",
  "baseFeeVaultRecipient": "0x...",
  "l1FeeVaultRecipient": "0x...",
  "sequencerFeeVaultRecipient": "0x...",
  "governanceTokenOwner": "0x...",
  "governanceTokenSymbol": "TRH",
  "governanceTokenName": "Tokamak Rollup Hub Token",
  "l2GenesisBlockGasLimit": "0x1c9c380",
  "gasPriceOracleOverhead": 2100,
  "gasPriceOracleScalar": 1000000,
  "enableGovernance": true,
  "eip1559Denominator": 50,
  "eip1559Elasticity": 10
}
```

#### STEP 6: 컨트랙트 빌드

```bash
cd tokamak-thanos/packages/tokamak/contracts-bedrock/scripts
bash ./start-deploy.sh build
```

**빌드 내용**:
- Solidity 컨트랙트 컴파일 (forge build)
- ABI 및 바이트코드 생성
- 배포 스크립트 준비

#### STEP 7: .env 파일 생성

```bash
export GS_ADMIN_PRIVATE_KEY={admin_private_key}
export L1_RPC_URL={l1_rpc_url}
export GAS_PRICE={suggested_gas_price * 2}
```

#### STEP 8: 컨트랙트 배포

```bash
cd tokamak-thanos/packages/tokamak/contracts-bedrock/scripts
bash ./start-deploy.sh deploy -e .env -c deploy-config.json
```

**배포되는 컨트랙트** (`start-deploy.sh` 내부):

```bash
# 1. ProxyAdmin
forge script --broadcast --rpc-url $L1_RPC_URL \
  --private-key $GS_ADMIN_PRIVATE_KEY \
  Deploy_ProxyAdmin.s.sol

# 2. AddressManager
forge script --broadcast ... Deploy_AddressManager.s.sol

# 3. SystemConfig
forge script --broadcast ... Deploy_SystemConfig.s.sol

# 4. L1CrossDomainMessenger
forge script --broadcast ... Deploy_L1CrossDomainMessenger.s.sol

# 5. L1StandardBridge
forge script --broadcast ... Deploy_L1StandardBridge.s.sol

# 6. L2OutputOracle
forge script --broadcast ... Deploy_L2OutputOracle.s.sol

# 7. OptimismPortal
forge script --broadcast ... Deploy_OptimismPortal.s.sol

# 8. DisputeGameFactory (Fault Proof)
forge script --broadcast ... Deploy_DisputeGameFactory.s.sol

# 9. ProtocolVersions
forge script --broadcast ... Deploy_ProtocolVersions.s.sol

# 10. DataAvailabilityChallenge (Plasma)
forge script --broadcast ... Deploy_DataAvailabilityChallenge.s.sol
```

배포 결과는 다음 파일에 저장됩니다:
```
tokamak-thanos/packages/tokamak/contracts-bedrock/deployments/{l1_chain_id}-deploy.json
```

#### STEP 9: Rollup 및 Genesis 파일 생성

```bash
cd tokamak-thanos/packages/tokamak/contracts-bedrock/scripts
bash ./start-deploy.sh generate -e .env -c deploy-config.json
```

**생성 파일**:

1. **rollup.json** (`tokamak-thanos/build/rollup.json`):
```json
{
  "genesis": {
    "l1": {
      "hash": "0x...",
      "number": 1234567
    },
    "l2": {
      "hash": "0x...",
      "number": 0
    },
    "l2_time": 1234567890,
    "system_config": {
      "batcherAddr": "0x...",
      "overhead": "0x...",
      "scalar": "0x...",
      "gasLimit": 30000000
    }
  },
  "block_time": 2,
  "max_sequencer_drift": 600,
  "seq_window_size": 3600,
  "channel_timeout": 300,
  "l1_chain_id": 11155111,
  "l2_chain_id": 12345678,
  "batch_inbox_address": "0xff00...",
  "deposit_contract_address": "0x...",
  "l1_system_config_address": "0x..."
}
```

2. **genesis.json** (`tokamak-thanos/build/genesis.json`):
```json
{
  "config": {
    "chainId": 12345678,
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
    "optimism": {
      "eip1559Denominator": 50,
      "eip1559Elasticity": 10
    }
  },
  "nonce": "0x0",
  "timestamp": "0x...",
  "extraData": "0x",
  "gasLimit": "0x1c9c380",
  "difficulty": "0x1",
  "mixHash": "0x...",
  "coinbase": "0x...",
  "alloc": {
    "0x4200...0006": {
      "code": "0x...",
      "storage": {...}
    },
    ...
  }
}
```

#### STEP 10: Candidate 등록 (선택)

`--no-candidate` 플래그가 없으면 후보로 등록:

```go
if t.registerCandidate {
    // Safe Wallet 설정
    t.setupSafeWallet(ctx, t.deploymentPath)

    // Candidate 검증 및 등록
    t.verifyRegisterCandidates(ctx, registerCandidate)
}
```

#### STEP 11: settings.json 생성

```json
{
  "admin_private_key": "0x...",
  "sequencer_private_key": "0x...",
  "batcher_private_key": "0x...",
  "proposer_private_key": "0x...",
  "deployment_path": "./tokamak-thanos/packages/tokamak/contracts-bedrock/deployments/11155111-deploy.json",
  "l1_rpc_url": "https://sepolia.infura.io/v3/...",
  "l1_beacon_url": "https://beacon.sepolia.io",
  "l1_rpc_provider": "infura",
  "l1_chain_id": 11155111,
  "l2_chain_id": 12345678,
  "stack": "thanos",
  "network": "testnet",
  "enable_fraud_proof": false,
  "chain_configuration": {
    "batch_submission_frequency": 300,
    "challenge_period": 604800,
    "output_root_frequency": 240,
    "l2_block_time": 2,
    "l1_block_time": 12
  },
  "deploy_contract_state": {
    "status": "completed"
  }
}
```

---

## Phase 2: 인프라 및 L2 체인 배포

### 명령어

```bash
trh-sdk deploy
```

### 전제 조건 확인

```go
// settings.json 존재 확인
config, err := utils.ReadConfigFromJSONFile(deploymentPath)

// L1 컨트랙트 배포 완료 확인
if config.DeployContractState.Status != types.DeployContractStatusCompleted {
    return fmt.Errorf("contracts are not deployed")
}
```

### 상세 동작 (`thanos/deploy_chain.go`)

#### STEP 1: 인프라 제공자 선택

```go
fmt.Print("Please select your infrastructure provider [AWS] (default: AWS): ")
input, err := scanner.ScanString()
infraOpt := strings.ToLower(input)
if infraOpt == "" {
    infraOpt = constants.AWS
}
```

#### STEP 2: AWS 로그인

```go
awsConfig, err := thanos.InputAWSLogin()
```

**입력 정보**:
```go
type AWSConfig struct {
    AccessKey      string
    SecretKey      string
    Region         string
    DefaultFormat  string  // "json"
}
```

**AWS CLI 설정**:
```bash
aws configure set aws_access_key_id {access_key}
aws configure set aws_secret_access_key {secret_key}
aws configure set region {region}
aws configure set output json
```

#### STEP 3: 배포 파라미터 입력

```go
inputs, err := thanos.InputDeployInfra()
```

**입력 정보**:
```go
type DeployInfraInput struct {
    ChainName          string
    L1BeaconURL        string
    IgnoreInstallBridge bool
    GithubCredentials  *GitHubCredentials  // Mainnet만
    MetadataInfo       *MetadataInfo       // Mainnet만
}
```

#### STEP 4: 의존성 확인

```go
dependencies.CheckTerraformInstallation(ctx)
dependencies.CheckHelmInstallation(ctx)
dependencies.CheckAwsCLIInstallation(ctx)
dependencies.CheckK8sInstallation(ctx)
```

#### STEP 5: 소스코드 클론

```go
err := t.cloneSourcecode(
    ctx,
    "tokamak-thanos-stack",
    "https://github.com/tokamak-network/tokamak-thanos-stack.git",
)
```

**tokamak-thanos-stack 구조**:
```
tokamak-thanos-stack/
├── terraform/
│   ├── backend/         # S3 + DynamoDB state backend
│   └── thanos-stack/    # EKS + 네트워크 리소스
└── charts/
    └── thanos-stack/    # Helm chart
```

#### STEP 6: AWS 인증

```go
// IAM 사용자 정보 조회
accountProfile := t.awsProfile.AccountProfile

// ~/.aws/credentials 및 config 설정 완료
```

#### STEP 7: Terraform 환경 파일 생성

```go
namespace := utils.ConvertChainNameToNamespace(inputs.ChainName)

err = makeTerraformEnvFile(
    fmt.Sprintf("%s/tokamak-thanos-stack/terraform", t.deploymentPath),
    types.TerraformEnvConfig{
        Namespace:           namespace,
        AwsRegion:           awsLoginInputs.Region,
        SequencerKey:        t.deployConfig.SequencerPrivateKey,
        BatcherKey:          t.deployConfig.BatcherPrivateKey,
        ProposerKey:         t.deployConfig.ProposerPrivateKey,
        ChallengerKey:       t.deployConfig.ChallengerPrivateKey,
        EksClusterAdmins:    awsAccountProfile.Arn,
        DeploymentFilePath:  t.deployConfig.DeploymentFilePath,
        L1BeaconUrl:         inputs.L1BeaconURL,
        L1RpcUrl:            t.deployConfig.L1RPCURL,
        L1RpcProvider:       t.deployConfig.L1RPCProvider,
        Azs:                 awsAccountProfile.AvailabilityZones,
        ThanosStackImageTag: "v1.7.7",
        OpGethImageTag:      "v1.101411.0",
        MaxChannelDuration:  chainConfiguration.GetMaxChannelDuration(),
    },
)
```

**.envrc 파일 내용**:
```bash
export TF_VAR_namespace="my-chain-testnet"
export TF_VAR_aws_region="ap-northeast-2"
export TF_VAR_sequencer_key="0x..."
export TF_VAR_batcher_key="0x..."
export TF_VAR_proposer_key="0x..."
export TF_VAR_challenger_key="0x..."
export TF_VAR_eks_cluster_admins="arn:aws:iam::123456789012:user/admin"
export TF_VAR_deployment_file_path="./tokamak-thanos/packages/.../11155111-deploy.json"
export TF_VAR_l1_beacon_url="https://beacon.sepolia.io"
export TF_VAR_l1_rpc_url="https://sepolia.infura.io/v3/..."
export TF_VAR_l1_rpc_provider="infura"
export TF_VAR_azs='["ap-northeast-2a","ap-northeast-2b","ap-northeast-2c"]'
export TF_VAR_thanos_stack_image_tag="v1.7.7"
export TF_VAR_op_geth_image_tag="v1.101411.0"
export TF_VAR_max_channel_duration="300"
```

#### STEP 8: 설정 파일 복사

```go
// rollup.json 복사
utils.CopyFile(
    fmt.Sprintf("%s/tokamak-thanos/build/rollup.json", t.deploymentPath),
    fmt.Sprintf("%s/tokamak-thanos-stack/terraform/thanos-stack/config-files/rollup.json", t.deploymentPath),
)

// genesis.json 복사
utils.CopyFile(
    fmt.Sprintf("%s/tokamak-thanos/build/genesis.json", t.deploymentPath),
    fmt.Sprintf("%s/tokamak-thanos-stack/terraform/thanos-stack/config-files/genesis.json", t.deploymentPath),
)
```

#### STEP 9: Terraform Backend 배포

```bash
cd tokamak-thanos-stack/terraform
source .envrc
cd backend
terraform init
terraform plan
terraform apply -auto-approve
```

**생성되는 리소스** (`backend/main.tf`):
```hcl
# S3 Bucket (Terraform State 저장)
resource "aws_s3_bucket" "terraform_state" {
  bucket = "${var.namespace}-terraform-state"

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}

# DynamoDB Table (State Lock)
resource "aws_dynamodb_table" "terraform_locks" {
  name         = "${var.namespace}-terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
```

#### STEP 10: Thanos Stack 인프라 배포

```bash
cd tokamak-thanos-stack/terraform
source .envrc
cd thanos-stack
terraform init
terraform plan
terraform apply -auto-approve
```

**생성되는 리소스** (`thanos-stack/main.tf`):

```hcl
# 1. VPC
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.namespace}-vpc"
  cidr = "10.0.0.0/16"

  azs             = var.azs
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = true
  enable_dns_hostnames = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = "1"
  }
}

# 2. EKS Cluster
module "eks" {
  source = "terraform-aws-modules/eks/aws"

  cluster_name    = "${var.namespace}-eks"
  cluster_version = "1.28"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  # Node Groups
  eks_managed_node_groups = {
    default = {
      desired_size = 3
      min_size     = 2
      max_size     = 5

      instance_types = ["t3.xlarge"]
      capacity_type  = "ON_DEMAND"

      labels = {
        role = "thanos-stack"
      }

      taints = []
    }
  }

  # Cluster access
  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = true

  enable_irsa = true
}

# 3. EFS (Persistent Storage)
resource "aws_efs_file_system" "thanos_storage" {
  creation_token = "${var.namespace}-efs"
  encrypted      = true

  tags = {
    Name = "${var.namespace}-thanos-storage"
  }
}

resource "aws_efs_mount_target" "thanos_storage" {
  count           = length(var.azs)
  file_system_id  = aws_efs_file_system.thanos_storage.id
  subnet_id       = module.vpc.private_subnets[count.index]
  security_groups = [aws_security_group.efs.id]
}

# 4. Security Groups
resource "aws_security_group" "efs" {
  name        = "${var.namespace}-efs-sg"
  description = "Allow EFS traffic from EKS nodes"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# 5. IAM Roles for Service Accounts (IRSA)
resource "aws_iam_role" "ebs_csi_driver" {
  name = "${var.namespace}-ebs-csi-driver"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = module.eks.oidc_provider_arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
    }]
  })
}

# 6. Add-ons
resource "aws_eks_addon" "ebs_csi_driver" {
  cluster_name = module.eks.cluster_name
  addon_name   = "aws-ebs-csi-driver"

  service_account_role_arn = aws_iam_role.ebs_csi_driver.arn
}

resource "aws_eks_addon" "vpc_cni" {
  cluster_name = module.eks.cluster_name
  addon_name   = "vpc-cni"
}

resource "aws_eks_addon" "coredns" {
  cluster_name = module.eks.cluster_name
  addon_name   = "coredns"
}

resource "aws_eks_addon" "kube_proxy" {
  cluster_name = module.eks.cluster_name
  addon_name   = "kube-proxy"
}

# 7. AWS Load Balancer Controller
resource "aws_iam_role" "aws_load_balancer_controller" {
  name = "${var.namespace}-aws-load-balancer-controller"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = module.eks.oidc_provider_arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "${replace(module.eks.cluster_oidc_issuer_url, "https://", "")}:sub": "system:serviceaccount:kube-system:aws-load-balancer-controller"
        }
      }
    }]
  })
}

# 8. Outputs for Helm values
output "vpc_id" {
  value = module.vpc.vpc_id
}

output "eks_cluster_name" {
  value = module.eks.cluster_name
}

output "efs_file_system_id" {
  value = aws_efs_file_system.thanos_storage.id
}
```

**리소스 생성 시간**: 약 15-20분

#### STEP 11: Helm Values 파일 생성

Terraform이 `thanos-stack-values.yaml`을 자동 생성합니다:

```yaml
# thanos-stack-values.yaml
namespace: my-chain-testnet

# Infrastructure
vpc_id: vpc-0123456789abcdef0
efs_file_system_id: fs-0123456789abcdef0

# Images
images:
  op_geth:
    repository: tokamaknetwork/thanos-geth
    tag: v1.101411.0
  op_node:
    repository: tokamaknetwork/thanos-op-node
    tag: v1.7.7
  op_batcher:
    repository: tokamaknetwork/thanos-op-batcher
    tag: v1.7.7
  op_proposer:
    repository: tokamaknetwork/thanos-op-proposer
    tag: v1.7.7
  op_challenger:
    repository: tokamaknetwork/thanos-op-challenger
    tag: v1.7.7

# Network Configuration
l1_rpc_url: https://sepolia.infura.io/v3/...
l1_beacon_url: https://beacon.sepolia.io
l1_rpc_provider: infura
max_channel_duration: "300"

# Operator Keys (Base64 encoded)
sequencer_private_key: "0x..."
batcher_private_key: "0x..."
proposer_private_key: "0x..."
challenger_private_key: "0x..."

# Deployment Configuration
deployment_json: |
  {
    "AddressManager": "0x...",
    "L1CrossDomainMessengerProxy": "0x...",
    ...
  }

rollup_json: |
  {
    "genesis": {...},
    "block_time": 2,
    ...
  }

genesis_json: |
  {
    "config": {...},
    "alloc": {...}
  }

# Storage
storage:
  op_geth:
    size: 100Gi
    storageClassName: gp3
  op_node:
    size: 50Gi
    storageClassName: gp3

# Resources
resources:
  op_geth:
    requests:
      cpu: 2000m
      memory: 8Gi
    limits:
      cpu: 4000m
      memory: 16Gi
  op_node:
    requests:
      cpu: 1000m
      memory: 4Gi
    limits:
      cpu: 2000m
      memory: 8Gi
  op_batcher:
    requests:
      cpu: 500m
      memory: 2Gi
    limits:
      cpu: 1000m
      memory: 4Gi
  op_proposer:
    requests:
      cpu: 500m
      memory: 2Gi
    limits:
      cpu: 1000m
      memory: 4Gi
  op_challenger:
    requests:
      cpu: 500m
      memory: 2Gi
    limits:
      cpu: 1000m
      memory: 4Gi

# Ingress
ingress:
  enabled: true
  className: alb
  annotations:
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}]'
  hosts:
    - host: rpc.my-chain-testnet.tokamak.network
      paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: op-geth
              port:
                number: 8545

# Enable deployment stages
enable_vpc: false       # 첫 번째 단계
enable_deployment: false  # 두 번째 단계
```

#### STEP 12: Kubernetes Context 전환

```go
err = utils.SwitchKubernetesContext(ctx, namespace, awsLoginInputs.Region)
```

```bash
# kubeconfig 업데이트
aws eks update-kubeconfig \
  --region {region} \
  --name {namespace}-eks \
  --alias {namespace}

# Context 전환
kubectl config use-context {namespace}
```

#### STEP 13: Kubernetes 클러스터 준비 대기

```go
k8sReady, err := utils.CheckK8sReady(ctx, namespace)
```

```bash
# 노드 준비 상태 확인
kubectl get nodes
# NAME                                           STATUS   ROLES    AGE   VERSION
# ip-10-0-1-123.ap-northeast-2.compute.internal   Ready    <none>   5m    v1.28.0
# ip-10-0-2-123.ap-northeast-2.compute.internal   Ready    <none>   5m    v1.28.0
# ip-10-0-3-123.ap-northeast-2.compute.internal   Ready    <none>   5m    v1.28.0
```

#### STEP 14: Helm Repository 추가

```bash
helm repo add thanos-stack https://tokamak-network.github.io/tokamak-thanos-stack
helm repo update
helm search repo thanos-stack
```

#### STEP 15: PVC 먼저 생성

```go
// enable_vpc = true 설정
err = utils.UpdateYAMLField(valueFile, "enable_vpc", true)

// Helm 설치
err = utils.InstallHelmRelease(ctx, helmReleaseName, chartFile, valueFile, namespace)

// PVC 준비 대기
err = utils.WaitPVCReady(ctx, namespace)
```

**생성되는 PVC**:
```bash
kubectl get pvc -n {namespace}
# NAME                    STATUS   VOLUME                   CAPACITY   ACCESS MODES
# op-geth-data-pvc        Bound    pvc-123456789abcdef0     100Gi      RWO
# op-node-data-pvc        Bound    pvc-123456789abcdef1     50Gi       RWO
```

#### STEP 16: 전체 배포 실행

```go
// enable_deployment = true 설정
err = utils.UpdateYAMLField(valueFile, "enable_deployment", true)

// Helm 재설치 (업그레이드)
err = utils.InstallHelmRelease(ctx, helmReleaseName, chartFile, valueFile, namespace)
```

**Helm chart가 배포하는 Kubernetes 리소스**:

```yaml
# 1. ConfigMaps
apiVersion: v1
kind: ConfigMap
metadata:
  name: thanos-config
  namespace: {{ .Values.namespace }}
data:
  rollup.json: |
    {{ .Values.rollup_json | nindent 4 }}
  genesis.json: |
    {{ .Values.genesis_json | nindent 4 }}
  deployment.json: |
    {{ .Values.deployment_json | nindent 4 }}

# 2. Secrets
apiVersion: v1
kind: Secret
metadata:
  name: operator-keys
  namespace: {{ .Values.namespace }}
type: Opaque
data:
  sequencer-key: {{ .Values.sequencer_private_key | b64enc }}
  batcher-key: {{ .Values.batcher_private_key | b64enc }}
  proposer-key: {{ .Values.proposer_private_key | b64enc }}
  challenger-key: {{ .Values.challenger_private_key | b64enc }}

# 3. StatefulSets
---
# op-geth
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: op-geth
  namespace: {{ .Values.namespace }}
spec:
  serviceName: op-geth
  replicas: 1
  selector:
    matchLabels:
      app: op-geth
  template:
    metadata:
      labels:
        app: op-geth
    spec:
      initContainers:
      - name: init-genesis
        image: {{ .Values.images.op_geth.repository }}:{{ .Values.images.op_geth.tag }}
        command:
        - sh
        - -c
        - |
          if [ ! -d /data/geth ]; then
            geth init --datadir=/data /config/genesis.json
          fi
        volumeMounts:
        - name: data
          mountPath: /data
        - name: config
          mountPath: /config
      containers:
      - name: op-geth
        image: {{ .Values.images.op_geth.repository }}:{{ .Values.images.op_geth.tag }}
        args:
        - --datadir=/data
        - --http
        - --http.addr=0.0.0.0
        - --http.port=8545
        - --http.api=eth,net,web3,debug,txpool
        - --http.vhosts=*
        - --http.corsdomain=*
        - --ws
        - --ws.addr=0.0.0.0
        - --ws.port=8546
        - --ws.api=eth,net,web3,debug,txpool
        - --authrpc.addr=0.0.0.0
        - --authrpc.port=8551
        - --authrpc.vhosts=*
        - --authrpc.jwtsecret=/secrets/jwt.txt
        - --rollup.disabletxpoolgossip
        - --rollup.sequencerhttp=http://op-node:8545
        - --nodiscover
        - --maxpeers=0
        ports:
        - containerPort: 8545
          name: http
        - containerPort: 8546
          name: ws
        - containerPort: 8551
          name: authrpc
        volumeMounts:
        - name: data
          mountPath: /data
        - name: config
          mountPath: /config
        - name: jwt
          mountPath: /secrets
        resources:
          {{ toYaml .Values.resources.op_geth | nindent 10 }}
      volumes:
      - name: config
        configMap:
          name: thanos-config
      - name: jwt
        secret:
          secretName: jwt-secret
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: {{ .Values.storage.op_geth.storageClassName }}
      resources:
        requests:
          storage: {{ .Values.storage.op_geth.size }}

---
# op-node
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: op-node
  namespace: {{ .Values.namespace }}
spec:
  serviceName: op-node
  replicas: 1
  template:
    spec:
      containers:
      - name: op-node
        image: {{ .Values.images.op_node.repository }}:{{ .Values.images.op_node.tag }}
        args:
        - op-node
        - --l1={{ .Values.l1_rpc_url }}
        - --l1.beacon={{ .Values.l1_beacon_url }}
        - --l1.rpckind={{ .Values.l1_rpc_provider }}
        - --l2=http://op-geth:8551
        - --l2.jwt-secret=/secrets/jwt.txt
        - --rollup.config=/config/rollup.json
        - --rpc.addr=0.0.0.0
        - --rpc.port=8545
        - --p2p.listen.ip=0.0.0.0
        - --p2p.listen.tcp=9003
        - --p2p.listen.udp=9003
        - --sequencer.enabled
        - --sequencer.l1-confs=4
        - --verifier.l1-confs=4
        - --p2p.sequencer.key={{ .Values.sequencer_private_key }}
        - --metrics.enabled
        - --metrics.addr=0.0.0.0
        - --metrics.port=7300
        - --pprof.enabled
        - --pprof.addr=0.0.0.0
        - --pprof.port=6060
        ports:
        - containerPort: 8545
          name: rpc
        - containerPort: 9003
          name: p2p-tcp
        - containerPort: 9003
          name: p2p-udp
          protocol: UDP
        - containerPort: 7300
          name: metrics
        - containerPort: 6060
          name: pprof

---
# op-batcher
apiVersion: apps/v1
kind: Deployment
metadata:
  name: op-batcher
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: op-batcher
        image: {{ .Values.images.op_batcher.repository }}:{{ .Values.images.op_batcher.tag }}
        env:
        - name: OP_BATCHER_L1_ETH_RPC
          value: {{ .Values.l1_rpc_url }}
        - name: OP_BATCHER_L2_ETH_RPC
          value: http://op-geth:8545
        - name: OP_BATCHER_ROLLUP_RPC
          value: http://op-node:8545
        - name: OP_BATCHER_MAX_CHANNEL_DURATION
          value: "{{ .Values.max_channel_duration }}"
        - name: OP_BATCHER_POLL_INTERVAL
          value: "6s"
        - name: OP_BATCHER_SUB_SAFETY_MARGIN
          value: "10"
        - name: OP_BATCHER_NUM_CONFIRMATIONS
          value: "10"
        - name: OP_BATCHER_PRIVATE_KEY
          valueFrom:
            secretKeyRef:
              name: operator-keys
              key: batcher-key

---
# op-proposer
apiVersion: apps/v1
kind: Deployment
metadata:
  name: op-proposer
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: op-proposer
        image: {{ .Values.images.op_proposer.repository }}:{{ .Values.images.op_proposer.tag }}
        env:
        - name: OP_PROPOSER_L1_ETH_RPC
          value: {{ .Values.l1_rpc_url }}
        - name: OP_PROPOSER_ROLLUP_RPC
          value: http://op-node:8545
        - name: OP_PROPOSER_POLL_INTERVAL
          value: "12s"
        - name: OP_PROPOSER_NUM_CONFIRMATIONS
          value: "10"
        - name: OP_PROPOSER_PRIVATE_KEY
          valueFrom:
            secretKeyRef:
              name: operator-keys
              key: proposer-key
        - name: OP_PROPOSER_L2OO_ADDRESS
          value: "{{ .Values.l2oo_address }}"

---
# op-challenger (Fault Proof)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: op-challenger
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: op-challenger
        image: {{ .Values.images.op_challenger.repository }}:{{ .Values.images.op_challenger.tag }}
        env:
        - name: OP_CHALLENGER_L1_ETH_RPC
          value: {{ .Values.l1_rpc_url }}
        - name: OP_CHALLENGER_L1_BEACON
          value: {{ .Values.l1_beacon_url }}
        - name: OP_CHALLENGER_L2_ETH_RPC
          value: http://op-geth:8545
        - name: OP_CHALLENGER_ROLLUP_RPC
          value: http://op-node:8545
        - name: OP_CHALLENGER_TRACE_TYPE
          value: "cannon"
        - name: OP_CHALLENGER_GAME_FACTORY_ADDRESS
          value: "{{ .Values.game_factory_address }}"
        - name: OP_CHALLENGER_DATADIR
          value: /data
        - name: OP_CHALLENGER_PRIVATE_KEY
          valueFrom:
            secretKeyRef:
              name: operator-keys
              key: challenger-key
        volumeMounts:
        - name: data
          mountPath: /data
        - name: op-program
          mountPath: /op-program

# 4. Services
---
apiVersion: v1
kind: Service
metadata:
  name: op-geth
spec:
  selector:
    app: op-geth
  ports:
  - name: http
    port: 8545
    targetPort: 8545
  - name: ws
    port: 8546
    targetPort: 8546
  - name: authrpc
    port: 8551
    targetPort: 8551

---
apiVersion: v1
kind: Service
metadata:
  name: op-node
spec:
  selector:
    app: op-node
  ports:
  - name: rpc
    port: 8545
    targetPort: 8545

# 5. Ingress
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: op-geth-ingress
  annotations:
    {{ toYaml .Values.ingress.annotations | nindent 4 }}
spec:
  ingressClassName: {{ .Values.ingress.className }}
  rules:
  - host: {{ .Values.ingress.hosts[0].host }}
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: op-geth
            port:
              number: 8545
```

#### STEP 17: RPC 엔드포인트 확인

```go
for {
    k8sIngresses, err := utils.GetAddressByIngress(ctx, namespace, helmReleaseName)
    if len(k8sIngresses) > 0 {
        l2RPCUrl = "http://" + k8sIngresses[0]
        break
    }
    time.Sleep(15 * time.Second)
}
```

```bash
kubectl get ingress -n {namespace}
# NAME              CLASS   HOSTS                          ADDRESS                         PORTS   AGE
# op-geth-ingress   alb     rpc.my-chain-testnet...        k8s-mychain-xxx.elb.amazonaws   80      5m

# L2 RPC URL
http://k8s-mychain-xxx.ap-northeast-2.elb.amazonaws.com
```

#### STEP 18: Backup 시스템 초기화

```go
err = t.initializeBackupSystem(ctx, inputs.ChainName)
```

EFS 백업을 위한 AWS Backup 설정

#### STEP 19: Bridge 플러그인 설치

```go
if !inputs.IgnoreInstallBridge {
    _, err = t.InstallBridge(ctx)
}
```

#### STEP 20: settings.json 업데이트

```json
{
  ...기존 설정...,
  "l2_rpc_url": "http://k8s-mychain-xxx.elb.amazonaws.com",
  "aws": {
    "access_key": "AKIA...",
    "secret_key": "...",
    "region": "ap-northeast-2",
    "vpc_id": "vpc-0123456789abcdef0"
  },
  "k8s": {
    "namespace": "my-chain-testnet"
  },
  "chain_name": "My Chain"
}
```

#### 배포 완료 메시지

```
✅ Network deployment completed successfully!
🌐 RPC endpoint: http://k8s-mychain-xxx.elb.amazonaws.com
✅ Backup system initialized successfully
✅ Bridge installed successfully
🎉 Thanos Stack installation completed successfully!
🚀 Your network is now up and running.
🔧 You can start interacting with your deployed infrastructure.
```

---

## 플러그인 시스템

TRH SDK는 선택적으로 설치 가능한 플러그인을 제공합니다.

### 지원 플러그인

```go
const (
    PluginBridge         = "bridge"
    PluginBlockExplorer  = "block-explorer"
    PluginMonitoring     = "monitoring"
)
```

### Bridge 플러그인

공식 Tokamak Bridge UI를 배포합니다.

#### 설치

```bash
trh-sdk install bridge
```

#### 동작 (`thanos/bridge.go`)

```go
func (t *ThanosStack) InstallBridge(ctx context.Context) (string, error) {
    // 1. Helm 릴리스 이름 생성
    helmReleaseName := fmt.Sprintf("%s-bridge", t.deployConfig.K8s.Namespace)

    // 2. Values 파일 생성
    valueFile := fmt.Sprintf("%s/bridge-values.yaml", t.deploymentPath)
    err := makeBridgeValuesFile(valueFile, t.deployConfig)

    // 3. Helm 차트 설치
    chartFile := "thanos-stack/bridge"
    err = utils.InstallHelmRelease(ctx, helmReleaseName, chartFile, valueFile, t.deployConfig.K8s.Namespace)

    // 4. Ingress 주소 조회
    bridgeUrl, err := utils.GetAddressByIngress(ctx, t.deployConfig.K8s.Namespace, helmReleaseName)

    return bridgeUrl, nil
}
```

**bridge-values.yaml**:
```yaml
l1_chain_id: 11155111
l2_chain_id: 12345678
l1_rpc_url: https://sepolia.infura.io/v3/...
l2_rpc_url: http://k8s-mychain-xxx.elb.amazonaws.com
l1_explorer_url: https://sepolia.etherscan.io
l2_explorer_url: http://explorer.my-chain.com

contracts:
  l1_standard_bridge: "0x..."
  l1_cross_domain_messenger: "0x..."
  optimism_portal: "0x..."

ingress:
  enabled: true
  className: alb
  host: bridge.my-chain-testnet.tokamak.network
```

**배포되는 리소스**:
- Deployment: Bridge frontend (React app)
- Service: ClusterIP
- Ingress: ALB를 통한 외부 노출

#### 제거

```bash
trh-sdk uninstall bridge
```

### Block Explorer 플러그인

Blockscout 블록 탐색기를 배포합니다.

#### 설치

```bash
trh-sdk install block-explorer
```

#### 동작 (`thanos/block_explorer.go`)

```go
func (t *ThanosStack) InstallBlockExplorer(ctx context.Context) (string, error) {
    // 1. PostgreSQL 먼저 배포
    err := t.deployBlockExplorerDB(ctx)

    // 2. Blockscout 배포
    helmReleaseName := fmt.Sprintf("%s-blockscout", t.deployConfig.K8s.Namespace)
    valueFile := fmt.Sprintf("%s/blockscout-values.yaml", t.deploymentPath)

    err = makeBlockscoutValuesFile(valueFile, t.deployConfig)
    err = utils.InstallHelmRelease(ctx, helmReleaseName, "thanos-stack/blockscout", valueFile, namespace)

    // 3. Ingress 주소 조회
    explorerUrl, err := utils.GetAddressByIngress(ctx, namespace, helmReleaseName)

    return explorerUrl, nil
}
```

**blockscout-values.yaml**:
```yaml
database:
  host: postgresql
  port: 5432
  name: blockscout
  username: blockscout
  password: {generated}

network:
  name: My Chain
  chain_id: 12345678
  rpc_url: http://op-geth:8545
  verification_url: http://op-geth:8545

indexer:
  block_interval: 5s
  first_block: 0
  trace_first_block: 0

resources:
  backend:
    requests:
      cpu: 1000m
      memory: 2Gi
  frontend:
    requests:
      cpu: 500m
      memory: 1Gi
```

### Monitoring 플러그인

Prometheus + Grafana 모니터링 스택을 배포합니다.

#### 설치

```bash
trh-sdk install monitoring
```

#### 동작 (`thanos/monitoring.go`)

```go
func (t *ThanosStack) InstallMonitoring(ctx context.Context) error {
    // 1. Prometheus 배포
    err := t.deployPrometheus(ctx)

    // 2. Grafana 배포
    err = t.deployGrafana(ctx)

    // 3. AlertManager 배포
    err = t.deployAlertManager(ctx)

    // 4. CloudWatch Exporter 배포 (AWS)
    err = t.deployCloudWatchExporter(ctx)

    return nil
}
```

**수집 메트릭**:
```yaml
# op-geth metrics (port 6060)
- geth_chain_head_block
- geth_chain_head_header
- geth_txpool_pending
- geth_txpool_queued

# op-node metrics (port 7300)
- op_node_l1_head
- op_node_l2_head
- op_node_l2_safe
- op_node_l2_finalized
- op_node_sequencer_active

# op-batcher metrics (port 7301)
- op_batcher_batch_size
- op_batcher_channel_size
- op_batcher_tx_count

# op-proposer metrics (port 7302)
- op_proposer_proposals
- op_proposer_last_proposal_time

# op-challenger metrics
- op_challenger_games_in_progress
- op_challenger_games_won
- op_challenger_games_lost
```

**Grafana 대시보드**:
- Thanos Stack Overview
- L1 → L2 Sync Status
- Transaction Performance
- Resource Usage (CPU, Memory, Disk)
- Alert History

#### Alert 설정

```bash
# Alert 상태 확인
trh-sdk alert-config --status

# 이메일 채널 설정
trh-sdk alert-config --channel email --configure

# Telegram 채널 설정
trh-sdk alert-config --channel telegram --configure

# 규칙 재설정
trh-sdk alert-config --rule reset
```

---

## 모니터링 및 관리

### 체인 정보 조회

```bash
trh-sdk info
```

**출력 예시**:
```
Chain Information
=================
Chain Name: My Chain
Network: testnet
Stack: thanos

L1 Configuration
----------------
Chain ID: 11155111
RPC URL: https://sepolia.infura.io/v3/...
Beacon URL: https://beacon.sepolia.io

L2 Configuration
----------------
Chain ID: 12345678
RPC URL: http://k8s-mychain-xxx.elb.amazonaws.com
Block Time: 2s

Contracts
---------
AddressManager: 0x...
L1CrossDomainMessenger: 0x...
L1StandardBridge: 0x...
L2OutputOracle: 0x...
OptimismPortal: 0x...
DisputeGameFactory: 0x...
SystemConfig: 0x...

Infrastructure
--------------
Provider: AWS
Region: ap-northeast-2
Namespace: my-chain-testnet
VPC ID: vpc-0123456789abcdef0
EKS Cluster: my-chain-testnet-eks

Plugins
-------
✅ Bridge: http://bridge.my-chain-testnet.com
✅ Block Explorer: http://explorer.my-chain-testnet.com
✅ Monitoring: http://grafana.my-chain-testnet.com
```

### 로그 확인

```bash
# 컴포넌트 로그
trh-sdk logs --component op-node

# 트러블슈팅 모드
trh-sdk logs --component op-challenger --troubleshoot
```

### 네트워크 업데이트

```bash
# settings.json 수정 후
trh-sdk update
```

### 백업 관리

```bash
# 백업 상태
trh-sdk backup-manager --status

# 스냅샷 생성
trh-sdk backup-manager --snapshot

# 백업 목록
trh-sdk backup-manager --list --limit 10

# 복구
trh-sdk backup-manager --restore
```

### SDK 업그레이드

```bash
trh-sdk upgrade
```

---

## 배포 플로우 다이어그램

### Local Devnet 배포

```
[trh-sdk deploy]
        │
        ├─> Config 확인 (settings.json 없음)
        │   └─> network=local-devnet, stack=thanos
        │
        ├─> ThanosStack 초기화
        │
        ├─> cloneSourcecode("tokamak-thanos")
        │
        └─> deployLocalDevnet()
            │
            └─> make devnet-up
                │
                ├─> Python: bedrock-devnet/main.py
                │   ├─> L1/L2 제네시스 생성
                │   ├─> 스마트 컨트랙트 배포
                │   └─> .devnet/ 설정 파일 생성
                │
                └─> Docker Compose up
                    ├─> l1 (geth)
                    ├─> l2 (geth)
                    ├─> op-node
                    ├─> op-proposer
                    ├─> op-batcher
                    ├─> op-challenger ⭐
                    ├─> da-server
                    └─> artifact-server
```

### Testnet/Mainnet 배포

```
Phase 1: L1 컨트랙트 배포
========================
[trh-sdk deploy-contracts --network testnet --stack thanos]
        │
        ├─> 사용자 입력 수집
        │   ├─> L1 RPC URL
        │   ├─> Operators (private keys)
        │   └─> Chain Configuration
        │
        ├─> L2 Chain ID 생성
        │
        ├─> 잔액 및 Gas 확인
        │
        ├─> cloneSourcecode("tokamak-thanos")
        │
        ├─> deploy-config.json 생성
        │
        ├─> bash start-deploy.sh build
        │   └─> Solidity 컴파일
        │
        ├─> bash start-deploy.sh deploy
        │   ├─> ProxyAdmin
        │   ├─> AddressManager
        │   ├─> SystemConfig
        │   ├─> L1CrossDomainMessenger
        │   ├─> L1StandardBridge
        │   ├─> L2OutputOracle
        │   ├─> OptimismPortal
        │   ├─> DisputeGameFactory ⭐
        │   └─> DataAvailabilityChallenge
        │
        ├─> bash start-deploy.sh generate
        │   ├─> rollup.json
        │   └─> genesis.json
        │
        ├─> Candidate 등록 (선택)
        │
        └─> settings.json 생성
            └─> deploy_contract_state.status = "completed"


Phase 2: 인프라 및 L2 체인 배포
==============================
[trh-sdk deploy]
        │
        ├─> settings.json 읽기
        │   └─> L1 컨트랙트 배포 완료 확인
        │
        ├─> 인프라 제공자 선택 (AWS)
        │
        ├─> AWS 로그인
        │   └─> aws configure
        │
        ├─> 배포 파라미터 입력
        │   ├─> Chain Name
        │   ├─> L1 Beacon URL
        │   └─> Ignore Bridge?
        │
        ├─> 의존성 확인
        │   ├─> Terraform
        │   ├─> Helm
        │   ├─> kubectl
        │   └─> AWS CLI
        │
        ├─> cloneSourcecode("tokamak-thanos-stack")
        │
        ├─> .envrc 생성 (Terraform 변수)
        │
        ├─> rollup.json, genesis.json 복사
        │
        ├─> Terraform Backend 배포
        │   ├─> S3 Bucket (state)
        │   └─> DynamoDB Table (lock)
        │
        ├─> Terraform Thanos Stack 배포
        │   ├─> VPC + Subnets
        │   ├─> EKS Cluster
        │   ├─> Node Groups
        │   ├─> EFS (persistent storage)
        │   ├─> Security Groups
        │   ├─> IAM Roles (IRSA)
        │   ├─> EBS CSI Driver
        │   └─> AWS Load Balancer Controller
        │
        ├─> kubectl context 전환
        │
        ├─> Helm repository 추가
        │
        ├─> Helm 배포 (2단계)
        │   ├─> Phase 1: PVC 생성
        │   │   └─> enable_vpc=true
        │   │
        │   └─> Phase 2: 전체 배포
        │       ├─> enable_deployment=true
        │       ├─> ConfigMaps (rollup.json, genesis.json)
        │       ├─> Secrets (operator keys)
        │       ├─> StatefulSets
        │       │   ├─> op-geth
        │       │   └─> op-node
        │       ├─> Deployments
        │       │   ├─> op-batcher
        │       │   ├─> op-proposer
        │       │   └─> op-challenger ⭐
        │       ├─> Services
        │       └─> Ingress (ALB)
        │
        ├─> RPC 엔드포인트 확인
        │   └─> L2 RPC URL 획득
        │
        ├─> Backup 시스템 초기화
        │   └─> AWS Backup (EFS)
        │
        ├─> Bridge 플러그인 설치
        │   └─> Helm install bridge
        │
        └─> settings.json 업데이트
            ├─> l2_rpc_url
            ├─> aws.vpc_id
            ├─> k8s.namespace
            └─> chain_name
```

### Kubernetes 리소스 배포 순서

```
Helm Install (Phase 1: PVC)
    │
    ├─> PVC 생성
    │   ├─> op-geth-data-pvc (100Gi)
    │   └─> op-node-data-pvc (50Gi)
    │
    └─> 준비 대기
        └─> PVC Bound 확인

Helm Install (Phase 2: Deployment)
    │
    ├─> ConfigMaps
    │   └─> thanos-config (rollup.json, genesis.json)
    │
    ├─> Secrets
    │   └─> operator-keys (private keys)
    │
    ├─> StatefulSet: op-geth
    │   ├─> Init Container: geth init
    │   └─> Container: geth run
    │       └─> Port 8545, 8546, 8551 READY
    │
    ├─> StatefulSet: op-node
    │   └─> Container: op-node
    │       ├─> L1 연결 대기
    │       ├─> L2 (op-geth:8551) 연결
    │       └─> Port 8545, 9003, 7300 READY
    │
    ├─> Deployment: op-batcher
    │   └─> Container: op-batcher
    │       ├─> op-node 대기
    │       ├─> L1에 배치 제출 시작
    │       └─> Port 7301 READY
    │
    ├─> Deployment: op-proposer
    │   └─> Container: op-proposer
    │       ├─> op-node 대기
    │       ├─> L1에 output 제안 시작
    │       └─> Port 7302 READY
    │
    ├─> Deployment: op-challenger ⭐
    │   └─> Container: op-challenger
    │       ├─> op-node, op-geth 대기
    │       ├─> DisputeGameFactory 모니터링 시작
    │       └─> Port 7303 READY
    │
    ├─> Services 생성
    │   ├─> op-geth (ClusterIP)
    │   └─> op-node (ClusterIP)
    │
    └─> Ingress 생성
        └─> ALB 프로비저닝
            ├─> DNS: k8s-mychain-xxx.elb.amazonaws.com
            └─> Target: op-geth:8545
```

---

## 참고 자료

### 공식 문서
- [TRH SDK GitHub](https://github.com/tokamak-network/trh-sdk)
- [Tokamak Thanos GitHub](https://github.com/tokamak-network/tokamak-thanos)
- [Tokamak Thanos Stack GitHub](https://github.com/tokamak-network/tokamak-thanos-stack)
- [Optimism Bedrock Specs](https://specs.optimism.io)

### 관련 파일
- `setup.sh`: SDK 설치 스크립트
- `cli.go`: CLI 명령어 정의
- `commands/*.go`: 각 명령어 구현
- `pkg/stacks/thanos/*.go`: Thanos Stack 구현

### AWS 리소스
- [AWS EKS 문서](https://docs.aws.amazon.com/eks/)
- [AWS EFS 문서](https://docs.aws.amazon.com/efs/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)

### Kubernetes 리소스
- [Helm 문서](https://helm.sh/docs/)
- [Kubernetes 공식 문서](https://kubernetes.io/docs/)
- [AWS Load Balancer Controller](https://kubernetes-sigs.github.io/aws-load-balancer-controller/)

---

**대상 SDK**: trh-sdk (Tokamak Rollup Hub SDK)
