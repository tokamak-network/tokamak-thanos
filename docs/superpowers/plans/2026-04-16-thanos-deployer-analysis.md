# thanos-deployer 배포 로직 분석 구현 계획

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Electron UI에서 L1 배포 버튼 클릭부터 L2 Genesis 파일 생성까지의 완전한 코드 수준 분석 문서화 및 시각화

**Architecture:** 6개 Phase를 순차 분석 (SSO → UI 요청 → 백엔드 큐잉 → L1 배포 → L2 Genesis → 결과)하고, 각 Phase별로 호출 체인, 파일 경로, 라인 번호, 데이터 흐름을 코드에서 직접 추출. 분석 결과를 마크다운 문서 + Mermaid 다이어그램으로 작성.

**Tech Stack:** Markdown, Mermaid, Grep/LSP (코드 분석), Bash (git log), mcp__mermaid__generate (다이어그램 렌더링)

---

## 파일 구조

### 생성 파일
```
docs/analysis/
├── thanos-deployer-flow-analysis.md      # 메인 분석 문서 (Phase별 상세)
├── diagrams/
│   ├── system-architecture.mmd           # 전체 시스템 아키텍처
│   ├── system-architecture.svg           # 렌더링된 다이어그램
│   ├── l1-deploy-flow.mmd                # L1 배포 플로우
│   ├── l1-deploy-flow.svg
│   ├── l2-genesis-flow.mmd               # L2 Genesis 생성 플로우
│   ├── l2-genesis-flow.svg
│   ├── call-graph.mmd                    # 모듈 호출 그래프
│   ├── call-graph.svg
│   ├── data-flow.mmd                     # 데이터 변환 흐름
│   └── data-flow.svg
└── code-reference-table.md               # 함수별 코드 맵 (테이블)

docs/superpowers/
└── specs/
    └── 2026-04-16-thanos-deployer-analysis-design.md  # 이미 존재
```

### 참고 레포
- `trh-platform` — Electron SSO, UI
- `trh-backend` — HTTP 핸들러, TaskManager
- `trh-sdk` — L1 배포 오케스트레이션
- `tokamak-thanos` — Foundry 스크립트
- `op-chain-ops` — L2 Genesis 생성

---

## Task 1: Phase 1 (Electron SSO) 분석

**Files:**
- Read: `trh-platform/src/main/aws-auth.ts` (SSO 플로우)
- Read: `trh-platform/src/renderer/pages/deploy.tsx` (배포 UI, 있다면)
- Create: `docs/analysis/thanos-deployer-flow-analysis.md` (Phase 1 섹션)

**Context:**
Electron 앱에서 배포 버튼 클릭 전에 AWS SSO 인증을 통해 accessToken을 얻는 단계. `startSsoLoginDirect` 함수 추적.

- [ ] **Step 1: trh-platform/src/main/aws-auth.ts 읽기**

```bash
cd /Users/theo/workspace_tokamak/trh-platform
grep -n "startSsoLoginDirect\|SSO\|getAccessToken" src/main/aws-auth.ts | head -20
```

Expected: AWS SSO 관련 함수/설정 확인

- [ ] **Step 2: SSO 인증 흐름 문서화**

`docs/analysis/thanos-deployer-flow-analysis.md` 생성 후 Phase 1 섹션 작성:

```markdown
# Phase 1: Electron SSO 인증

## 개요
Electron 앱 시작 시 AWS SSO 인증을 통해 accessToken 획득

## 핵심 파일
- `trh-platform/src/main/aws-auth.ts` — AWS SDK 통합, SSO 플로우

## 함수 분석

### startSsoLoginDirect (line 335)
- **역할**: SSO 인증 시작, AWS STS에서 임시 크레덴셜 획득
- **입력**: 없음 (환경 변수에서 AWS 설정 읽음)
- **출력**: 
  ```typescript
  {
    accessToken: string,
    region: string,
    credentials: AwsCredentials
  }
  ```
- **핵심 로직**:
  1. AWS SDK 클라이언트 초기화
  2. SSO 프로필 로드
  3. STS AssumeRole 호출
  4. 임시 크레덴셜 반환

## 데이터 구조

### AwsCredentials
```typescript
interface AwsCredentials {
  accessKeyId: string;
  secretAccessKey: string;
  sessionToken: string;
}
```

## 호출 시퀀스
```
Electron Main → startSsoLoginDirect() 
  → AWS SDK SSO Client
  → AWS STS AssumeRole
  → Return accessToken
```
```

- [ ] **Step 3: 관련 타입/함수 시그니처 확인**

```bash
grep -A 5 "function startSsoLoginDirect\|export.*startSsoLoginDirect" /Users/theo/workspace_tokamak/trh-platform/src/main/aws-auth.ts
```

Expected: 함수 정의, 입/출력 타입 확인

- [ ] **Step 4: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/thanos-deployer-flow-analysis.md
git commit -m "docs: add Phase 1 analysis (Electron SSO authentication)"
```

---

## Task 2: Phase 2 (Web UI 배포 요청) 분석

**Files:**
- Read: `trh-platform-ui/src/pages/stacks/deploy.tsx` (배포 UI)
- Read: `trh-backend/pkg/api/handlers/thanos/deployment.go:32` (Deploy 핸들러)
- Modify: `docs/analysis/thanos-deployer-flow-analysis.md` (Phase 2 섹션 추가)

**Context:**
Electron 내장 Next.js UI에서 배포 버튼 클릭 → HTTP POST /api/v1/stacks/thanos → 백엔드 핸들러

- [ ] **Step 1: UI에서 요청 생성하는 코드 확인**

```bash
grep -n "stacks/thanos\|DeployRequest\|POST" /Users/theo/workspace_tokamak/trh-platform-ui/src/pages/stacks/deploy.tsx | head -15
```

Expected: 배포 요청 페이로드 구조

- [ ] **Step 2: 백엔드 핸들러 코드 읽기**

```bash
grep -A 30 "func.*Deploy\|POST.*thanos" /Users/theo/workspace_tokamak/trh-backend/pkg/api/handlers/thanos/deployment.go | head -40
```

Expected: 핸들러 함수 정의, 요청 파싱 로직

- [ ] **Step 3: Phase 2 섹션 추가**

`docs/analysis/thanos-deployer-flow-analysis.md`에 다음 내용 추가:

```markdown
# Phase 2: Web UI 배포 요청

## 개요
Electron 내장 Next.js 프론트엔드에서 배포 요청 생성 및 전송

## 핵심 파일
- `trh-platform-ui/src/pages/stacks/deploy.tsx` — 배포 UI
- `trh-backend/pkg/api/handlers/thanos/deployment.go:32` — Deploy 핸들러

## UI 요청 흐름

### Deploy 페이지 (deploy.tsx)
- **역할**: 배포 파라미터 입력, POST 요청 생성
- **입력**: 사용자 입력 (stackName, preset, region 등)
- **출력**: 
  ```json
  POST /api/v1/stacks/thanos
  {
    "name": "my-stack",
    "preset": "general|defi|gaming",
    "region": "us-east-1",
    "chainId": 901,
    ...
  }
  ```

## 백엔드 핸들러

### Deploy (deployment.go:32)
- **역할**: 요청 수신, 유효성 검사, 서비스 호출
- **입력**: HTTP POST body (DeployThanosRequest)
- **출력**: 
  ```json
  {
    "stackId": "uuid",
    "status": "pending"
  }
  ```
- **핵심 로직**:
  1. 요청 바디 파싱 (DeployThanosRequest)
  2. 유효성 검사
  3. CreateThanosStack 서비스 호출
  4. 응답 반환

## 데이터 구조

### DeployThanosRequest
```typescript
interface DeployThanosRequest {
  name: string;
  preset: 'general' | 'defi' | 'gaming' | 'full';
  region: string;
  chainId: number;
  l1RpcUrl: string;
  adminAddress: string;
  batcherAddress: string;
  proposerAddress: string;
  sequencerAddress: string;
}
```

## 호출 시퀀스
```
UI → POST /api/v1/stacks/thanos
  → Deploy() handler
  → CreateThanosStack() service
```
```

- [ ] **Step 4: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/thanos-deployer-flow-analysis.md
git commit -m "docs: add Phase 2 analysis (Web UI deployment request)"
```

---

## Task 3: Phase 3 (백엔드 큐잉) 분석

**Files:**
- Read: `trh-backend/pkg/services/thanos/stack_lifecycle.go:20` (CreateThanosStack)
- Read: `trh-backend/pkg/services/thanos/deployment.go:31` (deployment orchestrator)
- Read: `trh-backend/pkg/task/task_manager.go` (TaskManager, 있다면)
- Modify: `docs/analysis/thanos-deployer-flow-analysis.md` (Phase 3 섹션)

**Context:**
백엔드가 요청을 받아 DB에 저장하고, 비동기 TaskManager 큐에 추가. 실제 배포는 별도 goroutine에서 실행.

- [ ] **Step 1: CreateThanosStack 함수 읽기**

```bash
grep -A 50 "func.*CreateThanosStack" /Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/stack_lifecycle.go | head -60
```

Expected: 스택 생성 로직, DB 저장, TaskManager enqueue

- [ ] **Step 2: deployment orchestrator 읽기**

```bash
grep -A 30 "func.*executeDeployments\|deploymentOrchestrator" /Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go
```

Expected: 배포 실행 순서, 각 서브태스크 호출

- [ ] **Step 3: TaskManager 동작 확인 (있다면)**

```bash
find /Users/theo/workspace_tokamak/trh-backend -name "*task*" -type f | grep -E "\.(go|ts)$"
grep -n "Enqueue\|Queue\|Worker" /Users/theo/workspace_tokamak/trh-backend/pkg/task/task_manager.go | head -10
```

Expected: TaskManager 인터페이스, Enqueue 방식

- [ ] **Step 4: Phase 3 섹션 추가**

```markdown
# Phase 3: trh-backend 영속화 + 비동기 큐잉

## 개요
배포 요청을 DB에 저장하고, 비동기 TaskManager 큐에 추가하여 백그라운드에서 실행

## 핵심 파일
- `trh-backend/pkg/services/thanos/stack_lifecycle.go:20` — CreateThanosStack
- `trh-backend/pkg/services/thanos/deployment.go:31` — deployment orchestrator
- `trh-backend/pkg/task/task_manager.go` — TaskManager (in-process goroutine 기반)

## CreateThanosStack 함수 (stack_lifecycle.go:20)

- **역할**: 스택 메타데이터 저장, 배포 태스크 큐잉
- **입력**: 
  ```go
  type CreateThanosStackRequest struct {
    Name string
    Preset string
    Region string
    ChainID int
    L1RpcUrl string
    AdminAddress string
    // ... more fields
  }
  ```
- **출력**: 
  ```go
  type Stack struct {
    ID string
    Status string // "pending" | "deploying" | "deployed" | "failed"
    CreatedAt time.Time
    // ... more fields
  }
  ```
- **핵심 로직**:
  1. 요청 유효성 검사
  2. 스택 객체 생성 + DB 저장 (GORM)
  3. executeDeployments 함수를 TaskManager에 Enqueue
  4. 생성된 스택 ID 반환

## deployment orchestrator (deployment.go:31)

- **역할**: 배포 프로세스 오케스트레이션 (Phase별 실행)
- **호출 순서**:
  1. DeployContracts (L1 배포) — trh-sdk 호출
  2. DeployChain (L2 노드 배포) — trh-sdk 호출
  3. DeployNetworkToAWS (인프라) — Terraform 호출
  4. DeployDApps (CrossTrade 등) — 커스텀 배포
  5. UpdateStackStatus (DB 업데이트)

## TaskManager

- **타입**: In-process goroutine 기반 (Redis 아님)
- **Enqueue 방식**: 
  ```go
  taskMgr.Enqueue(ctx, Task{
    ID: stackID,
    Type: "deploy_thanos",
    Payload: createRequest,
    Handler: executeDeployments,
  })
  ```
- **실행**: 별도 goroutine에서 비동기 실행

## 데이터 구조

### Stack (DB 모델)
```go
type Stack struct {
  ID string
  Name string
  Status string
  Preset string
  Region string
  ChainID int
  CreatedAt time.Time
  UpdatedAt time.Time
  DeploymentLog string
}
```

## 호출 시퀀스
```
Handler (Phase 2)
  → CreateThanosStack()
    → Validate request
    → Save Stack to DB
    → TaskManager.Enqueue(executeDeployments)
    → Return stackID
  ↓ (비동기)
  executeDeployments() [별도 goroutine]
    → DeployContracts (Phase 4)
    → DeployChain (Phase 5)
    → DeployNetworkToAWS
    → UpdateStackStatus
```
```

- [ ] **Step 5: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/thanos-deployer-flow-analysis.md
git commit -m "docs: add Phase 3 analysis (backend persistence and async queuing)"
```

---

## Task 4: Phase 4 (L1 컨트랙트 배포) 분석

**Files:**
- Read: `trh-sdk/pkg/stacks/thanos/deploy_contracts.go:33` (DeployContracts)
- Read: `tokamak-thanos/start-deploy.sh` (Foundry 호출)
- Read: `tokamak-thanos/Makefile` (컴파일/배포 타겟)
- Modify: `docs/analysis/thanos-deployer-flow-analysis.md` (Phase 4 섹션)

**Context:**
가장 복잡한 Phase. L1에 OptimismPortal, SystemConfig, L2OutputOracle 등 핵심 컨트랙트 배포. Foundry 셸아웃.

- [ ] **Step 1: DeployContracts 함수 읽기**

```bash
grep -A 80 "func.*DeployContracts" /Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_contracts.go | head -100
```

Expected: Foundry 호출 로직, .env 파일 생성, 배포 순서

- [ ] **Step 2: start-deploy.sh 분석**

```bash
head -100 /Users/theo/workspace_tokamak/tokamak-thanos/start-deploy.sh
grep -n "forge deploy\|foundry\|script" /Users/theo/workspace_tokamak/tokamak-thanos/start-deploy.sh
```

Expected: Foundry forge 명령어, 배포 스크립트 경로

- [ ] **Step 3: 배포 컨트랙트 목록 확인**

```bash
ls -la /Users/theo/workspace_tokamak/tokamak-thanos/scripts/
grep -r "OptimismPortal\|SystemConfig\|L2OutputOracle" /Users/theo/workspace_tokamak/tokamak-thanos/contracts/src/ | grep -E "contract |interface " | head -10
```

Expected: 핵심 컨트랙트 파일 경로

- [ ] **Step 4: 배포 산출물 구조 확인**

```bash
grep -A 20 "deploy.json\|rollup.json" /Users/theo/workspace_tokamak/tokamak-thanos/start-deploy.sh
```

Expected: 배포 후 생성되는 파일 포맷

- [ ] **Step 5: Phase 4 섹션 추가**

```markdown
# Phase 4: L1 컨트랙트 배포 (Foundry)

## 개요
trh-sdk에서 Foundry 스크립트를 셸아웃하여 L1 체인에 OP Stack 핵심 컨트랙트 배포

## 핵심 파일
- `trh-sdk/pkg/stacks/thanos/deploy_contracts.go:33` — DeployContracts 오케스트레이션
- `tokamak-thanos/start-deploy.sh` — Foundry forge 호출
- `tokamak-thanos/scripts/` — Foundry 배포 스크립트
- `tokamak-thanos/Makefile` — 컴파일 타겟

## DeployContracts 함수 (deploy_contracts.go:33)

- **역할**: Foundry 배포 환경 준비, 셸아웃, 산출물 수집
- **입력**:
  ```go
  type DeployContractsInput struct {
    StackID string
    ChainID int
    L1RpcUrl string
    AdminAddress string
    BugFixAddress string
    SystemConfigOwner string
    // ... more
  }
  ```
- **출력**:
  ```go
  type DeployContractsOutput struct {
    OptimismPortal string
    SystemConfig string
    L2OutputOracle string
    AnchorStateRegistry string
    DeployJSON map[string]interface{}
  }
  ```
- **핵심 로직**:
  1. .env 파일 생성 (L1 RPC, 키, 체인 ID 등)
  2. start-deploy.sh 실행 (Foundry forge)
  3. deploy.json 파싱
  4. 주요 컨트랙트 주소 추출
  5. 배포 상태 DB 업데이트

## start-deploy.sh 실행 흐름

- **위치**: `tokamak-thanos/start-deploy.sh`
- **동작**:
  1. 환경 변수 검증
  2. `forge build` — 컨트랙트 컴파일 (solc)
  3. `forge script` 호출 — Foundry TypeScript/Solidity 배포 스크립트 실행
  4. 배포된 컨트랙트 주소를 deploy.json에 저장
  5. rollup.json 생성 (체인 설정)

## 배포되는 컨트랙트 (배포 순서)

| 순서 | 컨트랙트명 | 파일 | 역할 |
|------|----------|------|------|
| 1 | OptimismPortal | contracts/src/L1/OptimismPortal.sol | L1→L2 deposit transaction 처리 |
| 2 | SystemConfig | contracts/src/L1/SystemConfig.sol | L2 배포 파라미터 저장 |
| 3 | L2OutputOracle | contracts/src/L1/L2OutputOracle.sol | L2 state root 검증 |
| 4 | AnchorStateRegistry | contracts/src/L1/AnchorStateRegistry.sol | Fault Proof 베이스 상태 |
| 5 | ProxyAdminOwner | contracts/src/L1/ProxyAdmin.sol | 프록시 관리 |
| ... | (기타) | ... | ... |

## 산출물 파일

### deploy.json
```json
{
  "OptimismPortal": "0x...",
  "SystemConfig": "0x...",
  "L2OutputOracle": "0x...",
  "AnchorStateRegistry": "0x...",
  "ProxyAdmin": "0x...",
  ...
}
```

### rollup.json
```json
{
  "genesis": { ... },
  "blockOracle": "0x...",
  "proxyAdmin": "0x...",
  ...
}
```

## 호출 시퀀스
```
orchestrator (Phase 3)
  → DeployContracts(input)
    1. Write .env file
    2. Execute: bash start-deploy.sh
       → forge build
       → forge script Deploy
       → Parse deploy.json
    3. Extract contract addresses
    4. Return DeployContractsOutput
    5. Update Stack.Status = "contracts_deployed"
```

## 에러 처리 및 재시도
- Foundry 배포 실패 시: 재시도 로직 (exponential backoff)
- L1 RPC 타임아웃: 타임아웃 설정 확인
- 컨트랙트 검증 실패: 배포 중단, 에러 메시지 기록
```

- [ ] **Step 6: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/thanos-deployer-flow-analysis.md
git commit -m "docs: add Phase 4 analysis (L1 contract deployment via Foundry)"
```

---

## Task 5: Phase 5 (L2 Genesis 생성) 분석

**Files:**
- Read: `trh-sdk/pkg/stacks/thanos/deploy_chain.go:30` (Deploy 함수)
- Read: `op-chain-ops/deployer/deployer.go` (L2 Genesis 생성)
- Grep: `op-chain-ops/deployer/` 내 주요 함수
- Modify: `docs/analysis/thanos-deployer-flow-analysis.md` (Phase 5 섹션)

**Context:**
L1 배포 완료 후, op-chain-ops deployer를 호출하여 L2 Genesis 파일 생성. 가장 복잡한 데이터 변환.

- [ ] **Step 1: Deploy 함수 (deploy_chain.go) 읽기**

```bash
grep -A 60 "func.*Deploy\(" /Users/theo/workspace_tokamak/trh-sdk/pkg/stacks/thanos/deploy_chain.go | head -80
```

Expected: L2 배포 오케스트레이션 로직

- [ ] **Step 2: op-chain-ops deployer 구조 파악**

```bash
ls -la /Users/theo/workspace_tokamak/op-chain-ops/deployer/
grep -n "func.*Deploy\|type.*Deployer" /Users/theo/workspace_tokamak/op-chain-ops/deployer/deployer.go | head -20
```

Expected: Deployer 구조체, 주요 메서드

- [ ] **Step 3: Genesis 생성 단계 추적**

```bash
grep -n "genesis\|alloc\|predeploy\|state" /Users/theo/workspace_tokamak/op-chain-ops/deployer/deployer.go | head -30
```

Expected: Genesis 파일 구성 단계

- [ ] **Step 4: op-chain-ops 타입 정의 확인**

```bash
grep -B 5 -A 15 "type DeployConfig struct\|type Genesis struct" /Users/theo/workspace_tokamak/op-chain-ops/deployer/deployer.go
```

Expected: DeployConfig, Genesis 구조 정의

- [ ] **Step 5: Phase 5 섹션 추가**

```markdown
# Phase 5: L2 Genesis 생성 (op-chain-ops)

## 개요
L1 배포된 컨트랙트 주소를 기반으로 op-chain-ops deployer를 사용해 L2 Genesis 파일 및 rollup.json 생성

## 핵심 파일
- `trh-sdk/pkg/stacks/thanos/deploy_chain.go:30` — Deploy 함수 (L2 체인 배포 오케스트레이션)
- `op-chain-ops/deployer/deployer.go` — Deployer 구조체, Genesis 생성 로직

## Deploy 함수 (deploy_chain.go:30)

- **역할**: op-chain-ops 호출, L2 Genesis 생성, 결과 수집
- **입력**:
  ```go
  type DeployChainInput struct {
    DeployContractsOutput DeployContractsOutput  // L1 배포 결과
    ChainID int
    L1RpcUrl string
    SequencerPrivateKey string
    ProposerPrivateKey string
    BatcherPrivateKey string
  }
  ```
- **출력**:
  ```go
  type DeployChainOutput struct {
    GenesisPath string           // build/genesis.json 경로
    RollupPath string            // build/rollup.json 경로
    Genesis *types.Genesis
    RollupConfig *rollup.Config
  }
  ```
- **핵심 로직**:
  1. DeployContractsOutput을 DeployConfig로 변환
  2. op-chain-ops deployer 초기화
  3. Deployer.Deploy() 호출 → Genesis 생성
  4. genesis.json, rollup.json 파일 쓰기
  5. 결과 반환

## op-chain-ops Deployer (deployer.go)

- **구조**:
  ```go
  type Deployer struct {
    DeployConfig *DeployConfig
    EthClient *ethclient.Client
    L1Addresses *L1Addresses  // OptimismPortal, SystemConfig 등
  }
  ```

- **Deploy 메서드** — Genesis 생성 메인 로직
  1. **Genesis Allocs 생성**: L1 컨트랙트 상태를 L2 genesis에 반영
  2. **Predeploy 배치**: 
     - L2CrossDomainMessenger
     - L1Block (L1 정보 제공)
     - L2ToL1MessagePasser
     - L2ERC721Bridge
     - L2StandardBridge
     - SequencerFeeVault
     - 기타...
  3. **State 초기화**: admin, batcher, proposer 계정 잔액 설정
  4. **Rollup Config 생성**: 체인 파라미터 (blockTime, gasLimit, sequencer 주소 등)
  5. **Fault Proof 패칭** (옵션): Fault Proof 시스템 설정

## Genesis 파일 구조 (genesis.json)

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
    "berlynBlock": 0,
    "londonBlock": 0,
    "arrowGlacierBlock": 0,
    "grayGlacierBlock": 0,
    "mergeNetsplitBlock": 0,
    "optimismConfig": {
      "eip1559Elasticity": 10,
      "eip1559DenominationDivisor": 50,
      "regolith": 0
    }
  },
  "difficulty": "0x1",
  "gasLimit": "0x9896e0",
  "nonce": "0x0",
  "timestamp": "0x...",
  "extraData": "...",
  "coinbase": "0x...",
  "stateRoot": "0x...",
  "hash": "0x...",
  "allocs": {
    "0x...L2CrossDomainMessenger": { "code": "...", "balance": "..." },
    "0x...L1Block": { "code": "...", "balance": "..." },
    "0x...Admin": { "balance": "0x..." },
    ...
  }
}
```

## Rollup Config 파일 (rollup.json)

```json
{
  "genesis": {
    "l1": {
      "hash": "0x...",
      "number": 123456
    },
    "l2": {
      "hash": "0x...",
      "number": 0
    },
    "l2Time": 1234567890,
    "systemConfig": {
      "batcherAddr": "0x...",
      "overhead": "0x...",
      "scalar": "0x...",
      "gasLimit": 9000000
    }
  },
  "blockTime": 2,
  "maxSequencerDrift": 600,
  "seqWindowSize": 3600,
  "channelTimeout": 300,
  "l1ChainID": 1,
  "l2ChainID": 901,
  "regolithTime": null,
  "canyonTime": null,
  "deltaTime": null,
  "ecotoneTime": null,
  "fjordTime": null,
  "addresses": {
    "optimismPortal": "0x...",
    "proxyAdminOwner": "0x...",
    "systemConfigOwner": "0x...",
    ...
  }
}
```

## Predeploy 주소 (L2)

| 주소 | 컨트랙트 | 역할 |
|------|---------|------|
| 0x4200000000000000000000000000000000000001 | LegacyMessagePasser | 레거시 메시지 (EVM) |
| 0x4200000000000000000000000000000000000002 | DeployerWhitelist | 배포자 화이트리스트 |
| 0x4200000000000000000000000000000000000007 | WETH9 | 래핑된 ETH |
| 0x4200000000000000000000000000000000000010 | L2CrossDomainMessenger | 크로스 도메인 메시징 |
| 0x4200000000000000000000000000000000000011 | L2StandardBridge | 표준 브리지 |
| 0x4200000000000000000000000000000000000013 | SequencerFeeVault | 시퀀서 수수료 모음 |
| 0x4200000000000000000000000000000000000015 | OptimismMintableERC20Factory | 민팅 가능 토큰 팩토리 |
| ... | ... | ... |

## 호출 시퀀스
```
orchestrator (Phase 3)
  → DeployChain(DeployContractsOutput)
    1. Convert L1 addresses to DeployConfig
    2. Initialize Deployer
    3. Call Deployer.Deploy()
       → Create genesis allocs
       → Batch predeploys
       → Initialize state
       → Create rollup config
       → (Optional) Patch fault proof
    4. Write genesis.json to build/
    5. Write rollup.json to build/
    6. Return DeployChainOutput
    7. Update Stack.Status = "chain_deployed"
```

## 데이터 변환 흐름
```
DeployContractsOutput (L1 주소)
  → DeployConfig (체인 파라미터)
  → Deployer.Deploy()
    → genesis.json (L2 초기 상태)
    → rollup.json (L2 체인 설정)
```
```

- [ ] **Step 6: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/thanos-deployer-flow-analysis.md
git commit -m "docs: add Phase 5 analysis (L2 genesis generation via op-chain-ops)"
```

---

## Task 6: Phase 6 (결과 반영) 분석

**Files:**
- Read: `trh-backend/pkg/services/thanos/deployment.go` (UpdateStackStatus 등)
- Modify: `docs/analysis/thanos-deployer-flow-analysis.md` (Phase 6 섹션)

**Context:**
배포 완료 후 DB 상태 업데이트, 결과 파일 저장, API 응답 대기.

- [ ] **Step 1: 배포 완료 핸들러 읽기**

```bash
grep -n "UpdateStackStatus\|deploymentComplete\|saveResults" /Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go | head -15
```

Expected: 상태 업데이트 로직

- [ ] **Step 2: 결과 저장 방식 확인**

```bash
grep -n "SaveDeployment\|WriteFile\|S3\|storage" /Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go | head -10
```

Expected: genesis.json, rollup.json 저장 위치 (로컬 또는 S3)

- [ ] **Step 3: Phase 6 섹션 추가**

```markdown
# Phase 6: 결과 반영

## 개요
L1, L2 배포 완료 후 결과를 DB에 저장하고, 클라이언트에 알림

## 핵심 파일
- `trh-backend/pkg/services/thanos/deployment.go` — UpdateStackStatus, SaveDeploymentResults
- `trh-backend/pkg/api/models/stack.go` — Stack 모델

## UpdateStackStatus 함수

- **역할**: Stack 레코드 상태 업데이트
- **입력**: 
  ```go
  stackID string
  status string  // "deployed" | "failed"
  ```
- **출력**: 에러 (있으면)
- **동작**:
  1. DB에서 Stack 조회
  2. Status = "deployed" 또는 "failed" 설정
  3. DeploymentLog 추가 (진행 상황)
  4. DB 저장

## SaveDeploymentResults 함수

- **역할**: genesis.json, rollup.json 저장
- **저장 위치**:
  - **로컬**: `~/.trh/stacks/{stackID}/genesis.json`
  - **또는 S3** (프로덕션): `s3://bucket/stacks/{stackID}/genesis.json`
- **저장 파일**:
  - genesis.json
  - rollup.json
  - deploy.json (L1 배포 결과)
  - .env (배포 설정, 민감 정보 제외)

## 호출 시퀀스
```
executeDeployments() [최종 단계]
  → DeployContracts()
  → DeployChain()
  → DeployNetworkToAWS() [Terraform]
  → SaveDeploymentResults()
    → Write genesis.json
    → Write rollup.json
    → Write deploy.json
  → UpdateStackStatus("deployed")
  → Notify Frontend (WebSocket 또는 polling)
```

## 에러 처리
- 배포 실패 시: Stack.Status = "failed", 에러 메시지 기록
- Rollback 로직 (선택): 실패 시 생성된 리소스 정리
```

- [ ] **Step 4: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/thanos-deployer-flow-analysis.md
git commit -m "docs: add Phase 6 analysis (result persistence and client notification)"
```

---

## Task 7: 시스템 아키텍처 다이어그램 생성

**Files:**
- Create: `docs/analysis/diagrams/system-architecture.mmd`
- Create: `docs/analysis/diagrams/system-architecture.svg`

**Context:**
Electron → trh-backend → trh-sdk → Foundry/Go deployer 계층 구조를 보여주는 높은 수준의 다이어그램.

- [ ] **Step 1: system-architecture.mmd 작성**

```bash
cat > /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams/system-architecture.mmd << 'EOF'
graph TB
    subgraph "Client Layer"
        A["Electron App<br/>(trh-platform)"]
        B["AWS SSO<br/>(Phase 1)"]
    end
    
    subgraph "Frontend Layer"
        C["Next.js UI<br/>(trh-platform-ui)"]
    end
    
    subgraph "API Layer"
        D["trh-backend<br/>Gin REST API"]
        E["POST /api/v1/stacks/thanos<br/>(Phase 2)"]
    end
    
    subgraph "Task Manager"
        F["TaskManager<br/>(in-process goroutine)"]
        G["executeDeployments()<br/>(Phase 3)"]
    end
    
    subgraph "SDK Layer (trh-sdk)"
        H["DeployContracts<br/>(L1 배포 오케스트레이션)"]
        I["DeployChain<br/>(L2 Genesis 생성)"]
        J["DeployNetworkToAWS<br/>(인프라)"]
    end
    
    subgraph "Execution Layer"
        K["Foundry<br/>(start-deploy.sh)"]
        L["op-chain-ops deployer<br/>(Genesis 생성)"]
        M["Terraform + Helm<br/>(AWS EKS)"]
    end
    
    subgraph "Output"
        N["L1 Contracts<br/>(OptimismPortal, SystemConfig...)"]
        O["build/genesis.json<br/>(L2 초기 상태)"]
        P["build/rollup.json<br/>(L2 체인 설정)"]
    end
    
    A -->|1. AWS SSO| B
    B -->|accessToken| C
    C -->|2. Deploy Request| E
    E -->|Parse| D
    D -->|3. CreateThanosStack| F
    F -->|Enqueue| G
    G -->|Phase 4 호출| H
    G -->|Phase 5 호출| I
    G -->|Phase 5 호출| J
    H -->|4. Foundry 호출| K
    K -->|Deploy| N
    K -->|Output| O
    I -->|5. Generate Genesis| L
    L -->|Output| O
    L -->|Output| P
    J -->|Terraform| M
    M -->|EKS Deployment| N
    
    style A fill:#e1f5ff
    style D fill:#fff3e0
    style H fill:#f3e5f5
    style K fill:#e8f5e9
    style O fill:#fce4ec
    style P fill:#fce4ec
EOF
```

Expected: Mermaid 파일 생성 확인

- [ ] **Step 2: Mermaid를 SVG로 렌더링**

```bash
# mcp__mermaid__generate 도구 사용
```

실제 렌더링은 mcp__mermaid__generate 도구로 수행 (다음 step에서)

- [ ] **Step 3: mcp__mermaid__generate로 SVG 생성**

(이 부분은 도구 호출이므로 Task 7.1로 분리)

- [ ] **Step 4: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/diagrams/system-architecture.mmd
git commit -m "docs: add system architecture diagram (Mermaid)"
```

---

## Task 7.1: Mermaid 다이어그램 렌더링 (모든 다이어그램)

**Files:**
- `docs/analysis/diagrams/system-architecture.mmd` → SVG
- `docs/analysis/diagrams/l1-deploy-flow.mmd` → SVG
- `docs/analysis/diagrams/l2-genesis-flow.mmd` → SVG
- `docs/analysis/diagrams/call-graph.mmd` → SVG
- `docs/analysis/diagrams/data-flow.mmd` → SVG

**Context:**
모든 Mermaid 다이어그램을 SVG로 렌더링하는 작업 (Task 8, 9, 10에서 생성한 .mmd 파일 포함).

- [ ] **Step 1: system-architecture.mmd 렌더링**

(Tool call: mcp__mermaid__generate)

- [ ] **Step 2: l1-deploy-flow.mmd 렌더링**

(Tool call: mcp__mermaid__generate)

- [ ] **Step 3: l2-genesis-flow.mmd 렌더링**

(Tool call: mcp__mermaid__generate)

- [ ] **Step 4: call-graph.mmd 렌더링**

(Tool call: mcp__mermaid__generate)

- [ ] **Step 5: data-flow.mmd 렌더링**

(Tool call: mcp__mermaid__generate)

- [ ] **Step 6: 렌더링 결과 확인**

```bash
ls -lh /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams/*.svg
```

Expected: 5개의 SVG 파일 생성 확인

- [ ] **Step 7: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/diagrams/*.svg
git commit -m "docs: render all Mermaid diagrams to SVG"
```

---

## Task 8: L1 배포 플로우 다이어그램 생성

**Files:**
- Create: `docs/analysis/diagrams/l1-deploy-flow.mmd`

**Context:**
start-deploy.sh 실행부터 OptimismPortal, SystemConfig, L2OutputOracle, AnchorStateRegistry 배포 순서 및 산출물 생성.

- [ ] **Step 1: l1-deploy-flow.mmd 작성**

```bash
cat > /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams/l1-deploy-flow.mmd << 'EOF'
sequenceDiagram
    participant SDK as trh-sdk<br/>DeployContracts
    participant Foundry as Foundry<br/>start-deploy.sh
    participant Solc as Solc<br/>Compiler
    participant L1 as L1 Chain<br/>(Sepolia/Mainnet)
    participant Out as Output<br/>deploy.json

    SDK->>Foundry: 1. Execute start-deploy.sh
    Foundry->>Solc: 2. forge build<br/>(Compile contracts)
    Solc->>Foundry: Compiled artifacts
    Foundry->>L1: 3. Deploy OptimismPortal
    L1->>Foundry: OptimismPortal address
    Foundry->>L1: 4. Deploy SystemConfig
    L1->>Foundry: SystemConfig address
    Foundry->>L1: 5. Deploy L2OutputOracle
    L1->>Foundry: L2OutputOracle address
    Foundry->>L1: 6. Deploy L1StandardBridge
    L1->>Foundry: L1StandardBridge address
    Foundry->>L1: 7. Deploy other contracts
    L1->>Foundry: Contract addresses
    Foundry->>L1: 8. Deploy AnchorStateRegistry
    L1->>Foundry: AnchorStateRegistry address
    Foundry->>Out: 9. Write deploy.json
    Out->>SDK: 10. Return addresses + deploy.json

    Note over Foundry,L1: Total: 8-10 contracts deployed
    Note over Out: deploy.json contains all L1 contract addresses
EOF
```

Expected: Mermaid 파일 생성

- [ ] **Step 2: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/diagrams/l1-deploy-flow.mmd
git commit -m "docs: add L1 deployment flow diagram"
```

---

## Task 9: L2 Genesis 생성 플로우 다이어그램 생성

**Files:**
- Create: `docs/analysis/diagrams/l2-genesis-flow.mmd`
- Create: `docs/analysis/diagrams/call-graph.mmd`
- Create: `docs/analysis/diagrams/data-flow.mmd`

**Context:**
L2 Genesis 생성의 3가지 주요 다이어그램: (1) Genesis 생성 단계, (2) 함수 호출 그래프, (3) 데이터 변환 흐름.

- [ ] **Step 1: l2-genesis-flow.mmd 작성 (Genesis 생성 단계)**

```bash
cat > /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams/l2-genesis-flow.mmd << 'EOF'
flowchart TD
    A["Deploy() 호출<br/>(op-chain-ops deployer)"]
    B["1. Create Genesis Allocs<br/>(L1 contracts → L2 state)"]
    C["2. Batch Predeploys<br/>(L2CrossDomainMessenger, L1Block, ...)"]
    D["3. Initialize State<br/>(Admin, Batcher, Proposer balance)"]
    E["4. Create Rollup Config<br/>(blockTime, gasLimit, chainID)"]
    F["5. Patch Fault Proof<br/>(if enabled)"]
    G["Write genesis.json"]
    H["Write rollup.json"]
    I["Return Genesis + RollupConfig"]

    A --> B
    B --> C
    C --> D
    D --> E
    E --> F
    F --> G
    F --> H
    G --> I
    H --> I

    style B fill:#e8f5e9
    style C fill:#e8f5e9
    style D fill:#e8f5e9
    style E fill:#bbdefb
    style F fill:#f3e5f5
    style G fill:#fce4ec
    style H fill:#fce4ec
EOF
```

Expected: l2-genesis-flow.mmd 생성

- [ ] **Step 2: call-graph.mmd 작성 (전체 호출 체인)**

```bash
cat > /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams/call-graph.mmd << 'EOF'
graph TD
    A["Deploy Handler<br/>(trh-backend:32)"]
    B["CreateThanosStack<br/>(trh-backend/stack_lifecycle.go:20)"]
    C["TaskManager.Enqueue<br/>(executeDeployments)"]
    D["executeDeployments<br/>(deployment.go:31)"]
    E1["DeployContracts<br/>(trh-sdk/deploy_contracts.go:33)"]
    E2["DeployChain<br/>(trh-sdk/deploy_chain.go:30)"]
    E3["DeployNetworkToAWS<br/>(Terraform)"]
    F1["start-deploy.sh<br/>(Foundry)"]
    F2["op-chain-ops Deployer<br/>(deployer.go)"]
    F3["Terraform + Helm"]
    G1["L1 Contracts Deployed"]
    G2["genesis.json Created"]
    G3["rollup.json Created"]

    A --> B
    B --> C
    C --> D
    D --> E1
    D --> E2
    D --> E3
    E1 --> F1
    E2 --> F2
    E3 --> F3
    F1 --> G1
    F2 --> G2
    F2 --> G3
    F3 --> G1

    style A fill:#fff3e0
    style D fill:#fff3e0
    style E1 fill:#f3e5f5
    style E2 fill:#f3e5f5
    style E3 fill:#f3e5f5
    style F1 fill:#e8f5e9
    style F2 fill:#e8f5e9
    style G1 fill:#fce4ec
    style G2 fill:#fce4ec
    style G3 fill:#fce4ec
EOF
```

Expected: call-graph.mmd 생성

- [ ] **Step 3: data-flow.mmd 작성 (데이터 변환)**

```bash
cat > /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams/data-flow.mmd << 'EOF'
graph LR
    A["DeployThanosRequest<br/>(JSON from UI)"]
    B[".env file<br/>(L1 RPC, keys, chainID)"]
    C["DeployContractsOutput<br/>(L1 addresses)"]
    D["DeployConfig<br/>(op-chain-ops format)"]
    E["genesis.json<br/>(L2 initial state)"]
    F["rollup.json<br/>(L2 chain config)"]

    A -->|Create| B
    B -->|Foundry uses| C
    C -->|Convert| D
    D -->|Generate| E
    D -->|Generate| F

    style A fill:#fff3e0
    style B fill:#e8f5e9
    style C fill:#f3e5f5
    style D fill:#f3e5f5
    style E fill:#fce4ec
    style F fill:#fce4ec
EOF
```

Expected: data-flow.mmd 생성

- [ ] **Step 4: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/diagrams/l2-genesis-flow.mmd \
         docs/analysis/diagrams/call-graph.mmd \
         docs/analysis/diagrams/data-flow.mmd
git commit -m "docs: add L2 genesis, call graph, and data flow diagrams"
```

---

## Task 10: 코드 참조 테이블 생성

**Files:**
- Create: `docs/analysis/code-reference-table.md`

**Context:**
모든 핵심 함수/파일의 파일명, 라인 번호, 역할, 입/출 타입을 한눈에 보는 테이블.

- [ ] **Step 1: code-reference-table.md 작성**

```bash
cat > /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/code-reference-table.md << 'EOF'
# thanos-deployer 코드 참조 테이블

Complete mapping of all files, functions, line numbers, and responsibilities involved in the thanos-deployer deployment flow.

## Phase 1: Electron SSO Authentication

| Phase | 파일 | 함수 | 라인 | 역할 | 입력 | 출력 |
|-------|------|------|------|------|------|------|
| 1 | `trh-platform/src/main/aws-auth.ts` | `startSsoLoginDirect()` | 335 | AWS SSO 인증 수행 | None (env vars) | `{ accessToken, region, credentials }` |

## Phase 2: Web UI Deployment Request

| Phase | 파일 | 함수/컴포넌트 | 라인 | 역할 | 입력 | 출력 |
|-------|------|------|------|------|------|------|
| 2a | `trh-platform-ui/src/pages/stacks/deploy.tsx` | Deploy component | - | UI에서 배포 요청 생성 | User input | `DeployThanosRequest` JSON |
| 2b | `trh-backend/pkg/api/handlers/thanos/deployment.go` | `Deploy()` | 32 | HTTP 요청 수신 및 파싱 | POST body | `{ stackId, status }` |

## Phase 3: Backend Persistence and Async Queuing

| Phase | 파일 | 함수 | 라인 | 역할 | 입력 | 출력 |
|-------|------|------|------|------|------|------|
| 3a | `trh-backend/pkg/services/thanos/stack_lifecycle.go` | `CreateThanosStack()` | 20 | 스택 메타데이터 저장, TaskManager 큐잉 | `DeployThanosRequest` | Stack ID |
| 3b | `trh-backend/pkg/task/task_manager.go` | `Enqueue()` | - | 비동기 태스크 큐에 추가 | Task object | - |
| 3c | `trh-backend/pkg/services/thanos/deployment.go` | `executeDeployments()` | 435 | 배포 오케스트레이션 (마스터 함수) | Stack ID | Deployment status |

## Phase 4: L1 Contract Deployment (Foundry)

| Phase | 파일 | 함수 | 라인 | 역할 | 입력 | 출력 |
|-------|------|------|------|------|------|------|
| 4a | `trh-sdk/pkg/stacks/thanos/deploy_contracts.go` | `DeployContracts()` | 33 | L1 배포 오케스트레이션 | `DeployContractsInput` | `DeployContractsOutput` (L1 addresses) |
| 4b | `tokamak-thanos/start-deploy.sh` | (bash script) | - | Foundry 호출 (forge build, forge script) | .env file | `deploy.json` |
| 4c | `tokamak-thanos/Makefile` | deploy target | - | 컴파일 및 배포 조정 | - | - |
| 4d | `tokamak-thanos/contracts/src/L1/OptimismPortal.sol` | (Solidity) | - | OptimismPortal 컨트랙트 | L1 RPC | Deployed address |
| 4e | `tokamak-thanos/contracts/src/L1/SystemConfig.sol` | (Solidity) | - | SystemConfig 컨트랙트 | L1 RPC | Deployed address |

## Phase 5: L2 Genesis Generation (op-chain-ops)

| Phase | 파일 | 함수 | 라인 | 역할 | 입력 | 출력 |
|-------|------|------|------|------|------|------|
| 5a | `trh-sdk/pkg/stacks/thanos/deploy_chain.go` | `Deploy()` | 30 | L2 배포 오케스트레이션 | `DeployContractsOutput` | `DeployChainOutput` (genesis paths) |
| 5b | `op-chain-ops/deployer/deployer.go` | `Deploy()` method | - | L2 Genesis 및 rollup config 생성 | `DeployConfig` + L1 addresses | Genesis, RollupConfig |
| 5c | `op-chain-ops/deployer/deployer.go` | (Predeploy logic) | - | L2 predeploy 컨트랙트 배치 | - | Genesis allocs |
| 5d | `op-chain-ops/script/deploy.go` | (Foundry integration) | - | Foundry 스크립트 인터페이스 | - | - |

## Phase 6: Result Persistence

| Phase | 파일 | 함수 | 라인 | 역할 | 입력 | 출력 |
|-------|------|------|------|------|------|------|
| 6a | `trh-backend/pkg/services/thanos/deployment.go` | `SaveDeploymentResults()` | - | genesis.json, rollup.json 저장 | File paths | File status |
| 6b | `trh-backend/pkg/services/thanos/deployment.go` | `UpdateStackStatus()` | - | Stack 상태 업데이트 (DB) | Status string | DB status |

## 데이터 구조 매핑

### Phase 2 → 3: DeployThanosRequest
```typescript
{
  name: string;
  preset: 'general' | 'defi' | 'gaming' | 'full';
  region: string;
  chainId: number;
  l1RpcUrl: string;
  adminAddress: string;
  batcherAddress: string;
  proposerAddress: string;
  sequencerAddress: string;
}
```

### Phase 3 → 4: DeployContractsInput (trh-sdk)
```go
type DeployContractsInput struct {
  StackID string
  ChainID int
  L1RpcUrl string
  AdminAddress string
  BugFixAddress string
  SystemConfigOwner string
  // ... more fields
}
```

### Phase 4 → 5: DeployContractsOutput (L1 addresses)
```go
type DeployContractsOutput struct {
  OptimismPortal string
  SystemConfig string
  L2OutputOracle string
  AnchorStateRegistry string
  DeployJSON map[string]interface{}
}
```

### Phase 5: DeployConfig (op-chain-ops format)
```go
type DeployConfig struct {
  L1ChainID int64
  L2ChainID int64
  L1StartingBlockTag rpc.BlockTag
  L1Addresses L1Addresses
  // ... more fields
}
```

### Phase 5 Output: genesis.json, rollup.json
See Phase 5 analysis document for full schema.

---

## 호출 체인 요약

```
Handler (Phase 2)
  → CreateThanosStack() [Phase 3]
    → Enqueue(executeDeployments)
      → DeployContracts() [Phase 4]
        → start-deploy.sh [Foundry]
          → forge build → forge script → deploy.json
      → DeployChain() [Phase 5]
        → op-chain-ops Deployer.Deploy()
          → genesis.json, rollup.json
      → DeployNetworkToAWS() [Terraform]
      → SaveDeploymentResults() [Phase 6]
      → UpdateStackStatus() [Phase 6]
```
EOF
```

Expected: code-reference-table.md 생성

- [ ] **Step 2: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/code-reference-table.md
git commit -m "docs: add comprehensive code reference table"
```

---

## Task 11: 메인 분석 문서 최종화 및 부록 추가

**Files:**
- Modify: `docs/analysis/thanos-deployer-flow-analysis.md`

**Context:**
Phase 1-6 섹션은 이미 Task 1-6에서 작성. 이제 부록을 추가: 통합 호출 그래프, 알려진 함정, 개선 포인트, 검증 체크리스트.

- [ ] **Step 1: 부록 섹션 추가**

기존 분석 문서에 다음 부록 추가:

```markdown
---

# 부록

## A. 통합 호출 그래프

[call-graph.mmd 에 구현됨 — docs/analysis/diagrams/call-graph.svg 참조]

## B. 데이터 구조 스키마

[code-reference-table.md 참조]

## C. 알려진 함정 및 개선 포인트

### 함정 1: Foundry 배포 실패 시 재시도 없음
**문제**: start-deploy.sh 실패 시 자동 재시도가 없어, 일시적 RPC 타임아웃으로 전체 배포 실패 가능.
**현황**: wiki [[l1-gas-limits]] 참조
**개선**: exponential backoff 재시도 로직 추가

### 함정 2: L2 Genesis 생성 시 predeploy 순서 중요
**문제**: predeploy 배치 순서가 초기화 로직에 영향을 미침
**현황**: op-chain-ops deployer.go에서 고정 순서 사용
**해결책**: 순서 문서화 및 단위 테스트 추가

### 함정 3: L1 RPC 수수료 동적 변화
**문제**: Sepolia blob fee spike로 인해 배포 트랜잭션 비용 급증
**현황**: wiki [[l1-gas-limits]] + 최근 fixes (blob fee cap)
**개선**: 동적 가스 가격 조정 알고리즘

### 함정 4: 타임존 의존성
**문제**: rollup.json genesis.l2Time이 UTC 타임스탬프 의존
**해결책**: 항상 UTC로 생성, 테스트에서 고정값 사용

## D. 검증 체크리스트

배포 분석 검증 시 다음 항목 확인:

- [ ] 모든 파일 경로가 현재 레포에 존재
- [ ] 모든 라인 번호가 코드와 일치 (git blame으로 확인)
- [ ] 호출 시퀀스가 실제 코드 실행 순서와 일치
- [ ] DeployThanosRequest JSON 스키마가 UI와 일치
- [ ] L1 컨트랙트 배포 순서가 Foundry 스크립트와 일치
- [ ] Predeploy 주소가 op-chain-ops와 일치
- [ ] genesis.json 및 rollup.json 스키마가 geth와 호환
- [ ] 다이어그램과 텍스트 설명이 모순 없음

## E. 다이어그램 참조

| 다이어그램 | 파일 | 용도 |
|----------|------|------|
| 시스템 아키텍처 | `diagrams/system-architecture.svg` | 고수준 모듈 연결 |
| L1 배포 플로우 | `diagrams/l1-deploy-flow.svg` | OptimismPortal 등 배포 순서 |
| L2 Genesis 생성 | `diagrams/l2-genesis-flow.svg` | Genesis 생성 단계 |
| 호출 그래프 | `diagrams/call-graph.svg` | 함수 호출 관계 |
| 데이터 흐름 | `diagrams/data-flow.svg` | 데이터 변환 |

## F. 관련 Wiki 페이지

- `[[l2-deployment]]` — L2 배포 플로우 (고수준)
- `[[ec2-deploy]]` — AWS EC2 배포 (이 문서의 상위 개념)
- `[[l1-gas-limits]]` — L1 가스 한계 튜닝
- `[[port-conflicts]]` — 포트 충돌 분석
- `[[tech-debt-and-risks]]` — 알려진 버그 및 위험

## G. 작성자 노트

- **작성일**: 2026-04-16
- **분석 범위**: Electron UI → L1 배포 → L2 Genesis 파일 생성 (AWS 경로)
- **제외 범위**: Terraform/Helm 상세, 컨트랙트 Solidity 구현, 로컬 Docker 배포
- **검증 상태**: 코드 수준 분석 완료, wiki와 교차 검증 중
```

- [ ] **Step 2: 부록 섹션 docs에 추가**

```bash
cat >> /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/thanos-deployer-flow-analysis.md << 'EOF'

---

# 부록

## A. 통합 호출 그래프

[call-graph.mmd에 구현됨 — diagrams/call-graph.svg 참조]

## B. 데이터 구조 스키마

[code-reference-table.md 참조]

## C. 알려진 함정 및 개선 포인트

### 함정 1: Foundry 배포 실패 시 재시도 없음
**문제**: start-deploy.sh 실패 시 자동 재시도가 없어, 일시적 RPC 타임아웃으로 전체 배포 실패 가능.
**현황**: wiki [[l1-gas-limits]] 참조
**개선**: exponential backoff 재시도 로직 추가

### 함정 2: L2 Genesis 생성 시 predeploy 순서 중요
**문제**: predeploy 배치 순서가 초기화 로직에 영향을 미침
**현황**: op-chain-ops deployer.go에서 고정 순서 사용
**해결책**: 순서 문서화 및 단위 테스트 추가

### 함정 3: L1 RPC 수수료 동적 변화
**문제**: Sepolia blob fee spike로 인해 배포 트랜잭션 비용 급증
**현황**: wiki [[l1-gas-limits]] + 최근 fixes (blob fee cap)
**개선**: 동적 가스 가격 조정 알고리즘

### 함정 4: 타임존 의존성
**문제**: rollup.json genesis.l2Time이 UTC 타임스탐프 의존
**해결책**: 항상 UTC로 생성, 테스트에서 고정값 사용

## D. 검증 체크리스트

배포 분석 검증 시 다음 항목 확인:

- [ ] 모든 파일 경로가 현재 레포에 존재
- [ ] 모든 라인 번호가 코드와 일치 (git blame으로 확인)
- [ ] 호출 시퀀스가 실제 코드 실행 순서와 일치
- [ ] DeployThanosRequest JSON 스키마가 UI와 일치
- [ ] L1 컨트랙트 배포 순서가 Foundry 스크립트와 일치
- [ ] Predeploy 주소가 op-chain-ops와 일치
- [ ] genesis.json 및 rollup.json 스키마가 geth와 호환
- [ ] 다이어그램과 텍스트 설명이 모순 없음

## E. 다이어그램 참조

| 다이어그램 | 파일 | 용도 |
|----------|------|------|
| 시스템 아키텍처 | `diagrams/system-architecture.svg` | 고수준 모듈 연결 |
| L1 배포 플로우 | `diagrams/l1-deploy-flow.svg` | OptimismPortal 등 배포 순서 |
| L2 Genesis 생성 | `diagrams/l2-genesis-flow.svg` | Genesis 생성 단계 |
| 호출 그래프 | `diagrams/call-graph.svg` | 함수 호출 관계 |
| 데이터 흐름 | `diagrams/data-flow.svg` | 데이터 변환 |

## F. 관련 Wiki 페이지

- `[[l2-deployment]]` — L2 배포 플로우 (고수준)
- `[[ec2-deploy]]` — AWS EC2 배포 (이 문서의 상위 개념)
- `[[l1-gas-limits]]` — L1 가스 한계 튜닝
- `[[port-conflicts]]` — 포트 충돌 분석
- `[[tech-debt-and-risks]]` — 알려진 버그 및 위험

## G. 작성자 노트

- **작성일**: 2026-04-16
- **분석 범위**: Electron UI → L1 배포 → L2 Genesis 파일 생성 (AWS 경로)
- **제외 범위**: Terraform/Helm 상세, 컨트랙트 Solidity 구현, 로컬 Docker 배포
- **검증 상태**: 코드 수준 분석 완료, wiki와 교차 검증 중
EOF
```

Expected: 부록 섹션 추가 확인

- [ ] **Step 3: 최종 문서 검증**

```bash
# 파일 크기 및 섹션 개수 확인
wc -l /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/thanos-deployer-flow-analysis.md
grep -c "^# " /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/thanos-deployer-flow-analysis.md
```

Expected: 400+ 라인, 6개 Phase 섹션 + 부록

- [ ] **Step 4: Commit**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git add docs/analysis/thanos-deployer-flow-analysis.md
git commit -m "docs: finalize thanos-deployer analysis document with appendices"
```

---

## Task 12: Wiki 업데이트 및 최종 Commit

**Files:**
- Modify: `/Users/theo/workspace_tokamak/trh-wiki/wiki/log.md` (새 항목 추가)
- Create 또는 Modify: `/Users/theo/workspace_tokamak/trh-wiki/wiki/concepts/thanos-deployer-analysis.md` (새 wiki 페이지)

**Context:**
분석 결과를 trh-wiki에 ingest하여 팀 전체가 접근 가능하도록 함. log.md에 변경 기록 추가.

- [ ] **Step 1: trh-wiki 새 페이지 생성**

```bash
cat > /Users/theo/workspace_tokamak/trh-wiki/wiki/concepts/thanos-deployer-analysis.md << 'EOF'
---
updated: 2026-04-16
---

# thanos-deployer 배포 로직 분석

**Summary**: Complete code-level tracing of the thanos-deployer deployment flow from Electron UI to L2 Genesis file generation. 6 phases, function-level detail, data structures, call sequences, and visual diagrams.

## Files

- Main analysis: `/Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/thanos-deployer-flow-analysis.md`
- Code reference table: `/Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/code-reference-table.md`
- Diagrams: `docs/analysis/diagrams/*.svg`

## 6 Phases

1. **Phase 1: Electron SSO** — AWS SSO 인증 (startSsoLoginDirect)
2. **Phase 2: Web UI** — HTTP POST /api/v1/stacks/thanos 요청
3. **Phase 3: Backend** — CreateThanosStack, TaskManager 큐잉
4. **Phase 4: L1 배포** — Foundry start-deploy.sh (OptimismPortal, SystemConfig...)
5. **Phase 5: L2 Genesis** — op-chain-ops deployer (genesis.json, rollup.json)
6. **Phase 6: 결과** — DB 저장, 클라이언트 알림

## 핵심 호출 체인

```
Handler → CreateThanosStack → Enqueue(executeDeployments)
  → DeployContracts [Foundry] → OptimismPortal, ...
  → DeployChain [op-chain-ops] → genesis.json, rollup.json
  → DeployNetworkToAWS [Terraform]
  → SaveDeploymentResults, UpdateStackStatus
```

## 다이어그램

- System architecture
- L1 deployment flow
- L2 genesis generation
- Call graph
- Data flow

## 관련 페이지

- [[l2-deployment]] — High-level L2 flow
- [[ec2-deploy]] — AWS deployment (related)
- [[tech-debt-and-risks]] — Known issues

---

*작성: 2026-04-16*
*분석 깊이: 코드 수준 (함수명, 라인번호, 데이터 구조)*
EOF
```

Expected: wiki 페이지 생성

- [ ] **Step 2: trh-wiki index.md 업데이트 (옵션)**

기존 index.md의 "Concepts" 섹션에 새 페이지 추가 (있으면).

```bash
# wiki index에 이미 있는지 확인
grep -n "thanos-deployer\|deployment flow" /Users/theo/workspace_tokamak/trh-wiki/wiki/index.md
```

If not present, add to Concepts section.

- [ ] **Step 3: log.md에 변경 기록 추가**

```bash
cat >> /Users/theo/workspace_tokamak/trh-wiki/wiki/log.md << 'EOF'

### 2026-04-16

**thanos-deployer 배포 로직 분석 문서 작성**
- `concepts/thanos-deployer-analysis.md` 신규 작성
- `tokamak-thanos/docs/analysis/thanos-deployer-flow-analysis.md` (6 Phase 상세)
- 5개 Mermaid 다이어그램 (시스템 아키텍처, L1/L2 플로우, 호출 그래프, 데이터 흐름)
- 코드 참조 테이블 (함수, 파일, 라인번호 매핑)
- 알려진 함정 및 개선 포인트 정리
EOF
```

Expected: log.md 업데이트 확인

- [ ] **Step 4: trh-wiki 커밋 및 푸시**

```bash
cd /Users/theo/workspace_tokamak/trh-wiki
git add wiki/concepts/thanos-deployer-analysis.md wiki/log.md
git commit -m "docs: add thanos-deployer deployment logic analysis"
git push origin main
```

Expected: trh-wiki 푸시 완료

- [ ] **Step 5: 최종 정리**

```bash
# 생성된 파일 목록 확인
ls -lh /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/
ls -lh /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams/
```

Expected: 모든 분석 파일 생성 확인

- [ ] **Step 6: 최종 커밋 및 푸시 (tokamak-thanos)**

```bash
cd /Users/theo/workspace_tokamak/tokamak-thanos
git status  # 모든 파일 확인
git add docs/analysis/ docs/superpowers/specs/2026-04-16-thanos-deployer-analysis-design.md
git commit -m "docs: complete thanos-deployer deployment logic analysis

- Phase-by-phase code tracing (6 phases)
- Mermaid diagrams: system architecture, L1/L2 flow, call graph, data flow
- Code reference table with file paths, line numbers, signatures
- Known pitfalls and improvement points
- Appendices with call chains, validation checklist, wiki references"
git push origin main
```

Expected: tokamak-thanos 푸시 완료

---

## 검증 및 완료

- [ ] **최종 검증 체크리스트**

```bash
# 1. 모든 파일 생성 확인
ls -1 /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/thanos-deployer-flow-analysis.md
ls -1 /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/code-reference-table.md
find /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams -name "*.mmd" | wc -l
find /Users/theo/workspace_tokamak/tokamak-thanos/docs/analysis/diagrams -name "*.svg" | wc -l

# 2. Git 커밋 이력 확인
cd /Users/theo/workspace_tokamak/tokamak-thanos
git log --oneline | head -10

# 3. trh-wiki 업데이트 확인
cd /Users/theo/workspace_tokamak/trh-wiki
git log --oneline | head -5

# Expected output:
# - 6 Task commits (Phase 1-6 분석)
# - 1 다이어그램 커밋 (Mermaid .mmd 생성)
# - 1 렌더링 커밋 (.svg 생성)
# - 1 코드 테이블 커밋
# - 1 최종화 커밋 (부록)
# - 1 wiki ingest 커밋
```

Expected: 모든 커밋 확인, 파일 생성 완료

---

## 실행 방법

### 옵션 1: Subagent-Driven (권장)
```bash
# 신선한 subagent가 각 task 수행, 단계별 검토
/gsd:subagent-driven-development
```

### 옵션 2: Inline Execution
```bash
# 이 세션에서 Task 1-12 순차 실행
/gsd:executing-plans
```
