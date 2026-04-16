---
title: Thanos Deployer - Complete Flow Analysis
date: 2026-04-16
status: Final
description: End-to-end analysis of Thanos Stack deployment flow with 6 appendices synthesizing all prior analysis work
---

# Thanos Deployer - Complete Flow Analysis

## Table of Contents

1. [Overview](#overview)
2. [Document Structure](#document-structure)
3. [Quick Start: Finding Information](#quick-start-finding-information)
4. [Phase Analyses](#phase-analyses)
5. [Appendices](#appendices)
   - [Appendix A: Unified Call Graph](#appendix-a-unified-call-graph)
   - [Appendix B: Data Structures & Transformations](#appendix-b-data-structures--transformations)
   - [Appendix C: Known Pitfalls & Improvement Points](#appendix-c-known-pitfalls--improvement-points)
   - [Appendix D: Verification Checklist](#appendix-d-verification-checklist)
   - [Appendix E: Diagram Reference & Index](#appendix-e-diagram-reference--index)
   - [Appendix F: Related Wiki & Documentation](#appendix-f-related-wiki--documentation)
6. [References](#references)

---

## Overview

This document provides a **complete, code-level analysis** of how a Thanos Stack deployment request flows from the Electron UI through the backend, contract deployment, and L2 genesis generation. It synthesizes 6 detailed phase analyses, a comprehensive code reference table, and visual diagrams into a single authoritative reference.

### Scope

**Start**: HTTP `POST /api/v1/stacks/thanos` request from Electron UI
**End**: `genesis.json` + `rollup.json` written to disk, client notified of completion

**Repositories Covered**:
- `trh-platform` — Electron shell and SSO
- `trh-backend` — HTTP handlers, task manager, orchestration
- `trh-sdk` — L1 deployment wrapper
- `tokamak-thanos` — Foundry scripts, L2 genesis generation
- `op-chain-ops` — Optimism stack genesis builder

**What's Included**:
- Complete call chain from HTTP handler to file output
- All data structure transformations
- Critical pitfalls and known issues
- Hard fork compatibility considerations
- Predeploy contract injection mechanisms

**What's NOT Included**:
- Solidity contract implementations (only deployment mechanism)
- Terraform/Helm deployment (AWS infrastructure, covered separately)
- Local Docker deployment (AWS focus)
- Cross-Trade module details (Preset layer)

---

## Document Structure

This analysis is organized in **6 sequential phases** plus **6 synthesizing appendices**:

### Phases 1-6

Each phase represents a distinct deployment stage:

| Phase | Name | Focus | Primary Files |
|-------|------|-------|----------------|
| 1 | Electron SSO | User authentication | trh-platform |
| 2 | Web UI Request | Frontend to backend HTTP | trh-platform-ui, API layer |
| 3 | Backend Queuing | Task manager, database persistence | trh-backend services, DB |
| 4 | L1 Contract Deploy | Foundry script execution | tokamak-deployer, start-deploy.sh |
| 5 | L2 Genesis Generation | Predeploy initialization, hard fork application | op-chain-ops |
| 6 | Result Persistence | Database updates, client notification | trh-backend, long-polling |

**Full Phase Analyses**: See `PHASE_1_ANALYSIS.md` through `PHASE_6_ANALYSIS.md`

---

## Quick Start: Finding Information

**I want to understand...**

- **The big picture**: See [Appendix E: Diagram Reference](#appendix-e-diagram-reference--index)
- **How function X calls function Y**: See [Appendix A: Unified Call Graph](#appendix-a-unified-call-graph)
- **What data structure is used at stage X**: See [Appendix B: Data Structures](#appendix-b-data-structures--transformations)
- **Line-by-line code details**: See `code-reference-table.md` (functions with file paths and line numbers)
- **Why thing X fails or behaves unexpectedly**: See [Appendix C: Known Pitfalls](#appendix-c-known-pitfalls--improvement-points)
- **How to verify this analysis is correct**: See [Appendix D: Verification Checklist](#appendix-d-verification-checklist)

---

## Phase Analyses

For complete details on each phase, see the dedicated analysis documents:

1. **[PHASE_1_ANALYSIS.md](./PHASE_1_ANALYSIS.md)** — Electron SSO authentication workflow
2. **[PHASE_2_ANALYSIS.md](./PHASE_2_ANALYSIS.md)** — Frontend preset wizard and HTTP request submission
3. **[PHASE_3_ANALYSIS.md](./PHASE_3_ANALYSIS.md)** — Backend queuing, stack creation, task orchestration
4. **[PHASE_4_ANALYSIS.md](./PHASE_4_ANALYSIS.md)** — L1 contract deployment via Foundry and tokamak-deployer
5. **[PHASE_5_ANALYSIS.md](./PHASE_5_ANALYSIS.md)** — L2 genesis generation, predeploy initialization, hard forks
6. **[PHASE_6_ANALYSIS.md](./PHASE_6_ANALYSIS.md)** — Result persistence, metadata updates, client notification

---

## Appendices

---

## Appendix A: Unified Call Graph

### Complete Call Chain: HTTP Request to Client Notification

This appendix shows the **complete end-to-end function call hierarchy** from the initial HTTP POST request through all deployment phases to final client notification.

```
HTTP POST /api/v1/stacks/thanos
│
├─── PHASE 2: Web UI Request Handler
│    │
│    └─ PresetDeploy() [trh-backend/pkg/api/handlers/thanos/presets.go:78-104]
│       │  ├─ Input: JSON body (DeployWithPresetRequest)
│       │  └─ Output: HTTP 200 + {stackId, taskId, status}
│       │
│       └─ PresetDeployRequest → Service validation
│
├─── PHASE 3: Backend Queuing & Orchestration
│    │
│    └─ CreateThanosStackFromPreset() [trh-backend/pkg/services/thanos/preset_deploy.go]
│       │  ├─ Input: PresetDeployRequest
│       │  └─ Output: Stack entity created in DB
│       │
│       ├─ createStackFromPreset()
│       │  └─ Expand Preset → DeployConfig
│       │
│       └─ CreateThanosStack() [trh-backend/pkg/services/thanos/stack_lifecycle.go:20-65]
│          │  ├─ Atomic transaction: Create Stack, Deployments, IntegrationEntities
│          │  └─ Persist to PostgreSQL
│          │
│          └─ startPresetDeployment()
│             │
│             └─ TaskManager.AddTask() [trh-backend/pkg/task/task_manager.go]
│                │
│                └─ WorkerPool.Execute() [trh-backend/pkg/task/worker_pool.go]
│                   │
│                   └─ [ASYNC EXECUTION BEGINS]
│
├─── PHASE 4: L1 Contract Deployment
│    │
│    └─ deploy() [trh-backend/pkg/services/thanos/deployment.go:31-437]
│       │  ├─ Main orchestrator for Phases 4-5
│       │  └─ Called by: WorkerPool.Execute()
│       │
│       └─ executeDeployments() [trh-backend/pkg/services/thanos/deployment.go:439-653]
│          │  ├─ Filter deployments by type
│          │  └─ Execute sequentially: deploy_l1_contracts → deploy_aws_infrastructure
│          │
│          └─ (For "deploy_l1_contracts" deployment)
│             │
│             ├─ DeployL1Contracts() [trh-backend/pkg/sdk/thanos_stack.go:129-173]
│             │  │  ├─ Input: DeployConfig
│             │  │  └─ Output: L1DeployOutput (contract addresses)
│             │  │
│             │  └─ Download tokamak-deployer binary (version-pinned)
│             │     │
│             │     └─ DeployContracts() [tokamak-deployer/pkg/deployer/deploy_contracts.go:23-407]
│             │        │  ├─ Input: DeployContractsInput (L1RPC, keys, config)
│             │        │  └─ Output: DeployOutput (contract addresses, txhashes)
│             │        │
│             │        ├─ GetDeploymentOrder() [contract_order.go]
│             │        │  └─ Returns ordered list: [SuperchainConfig, Proxies, Implementations, ...]
│             │        │
│             │        ├─ For each contract in order:
│             │        │  │
│             │        │  ├─ Create2Deploy() [create2_deployer.go]
│             │        │  │  └─ Deploy with deterministic address (CREATE2)
│             │        │  │
│             │        │  └─ Track address, timestamp, block hash
│             │        │
│             │        ├─ [IF EnableFaultProof]
│             │        │  │
│             │        │  └─ BuildCannonPrestate() [fault_proof/cannon_prestate.go]
│             │        │     └─ Generate Cannon prestate for fault proof system
│             │        │
│             │        └─ WriteDeployOutput() [utils/file_utils.go]
│             │           └─ Output: deploy-output.json with all contract addresses
│             │
│             └─ [Return to executeDeployments with L1DeployOutput]
│
├─── PHASE 5: L2 Genesis Generation
│    │
│    └─ (Continued in executeDeployments)
│       │
│       ├─ BuildL2Genesis() [op-chain-ops/genesis/layer_two.go:39-204]
│       │  │  ├─ Input: L1DeployOutput, DeployConfig, optional Cannon prestate
│       │  │  └─ Output: core.Genesis structure (all predeploys initialized)
│       │  │
│       │  ├─ NewL2Genesis() [layer_two.go]
│       │  │  └─ Create base genesis structure with:
│       │  │     • ChainID, Timestamp, GasLimit
│       │  │     • Block root, state root placeholders
│       │  │
│       │  ├─ NewL2ImmutableConfig() [layer_two.go]
│       │  │  └─ Extract from L1DeployOutput:
│       │  │     • 40+ predeploy immutable values
│       │  │     • Bridge addresses, oracle refs, etc.
│       │  │
│       │  ├─ NewL2StorageConfig() [layer_two.go]
│       │  │  └─ Initialize storage for specific contracts:
│       │  │     • ProxyAdmin permissions
│       │  │     • Fee vault recipients
│       │  │     • System parameters
│       │  │
│       │  ├─ setProxies() [layer_two.go]
│       │  │  └─ Deploy 256 transparent proxy contracts @ 0x4200...
│       │  │
│       │  ├─ For each of 40+ predeploy contracts:
│       │  │  │
│       │  │  ├─ Deploy() [immutables/immutables.go]
│       │  │  │  └─ Inject immutable config into bytecode (bytecode patching)
│       │  │  │     • Replace hardcoded offsets with actual addresses
│       │  │  │     • Validate bytecode hash before patching
│       │  │  │
│       │  │  └─ setupPredeploy() [layer_two.go]
│       │  │     └─ Add to genesis.alloc:
│       │  │        • Bytecode with injected immutables
│       │  │        • Storage values
│       │  │        • Balance (usually 0)
│       │  │
│       │  ├─ PerformUpgradeTxs() [layer_two.go]
│       │  │  └─ Apply hard fork changes:
│       │  │     • Ecotone: Update specific storage slots
│       │  │     • Fjord: Enable new transaction format
│       │  │     • Granite: Additional contract updates
│       │  │
│       │  ├─ NewRollupConfig() [rollup/config.go]
│       │  │  └─ Generate rollup.json:
│       │  │     • L1 chain info, block time, finality period
│       │  │     • Hard fork activation times
│       │  │     • Batch overhead, scalar, gas limit
│       │  │
│       │  └─ WriteGenesisJSON() [utils/genesis_writer.go]
│       │     └─ Output: genesis.json (20-100MB) + rollup.json (5-50KB)
│       │
│       └─ [Return with genesis.json and rollup.json written]
│
└─── PHASE 6: Result Persistence & Client Notification
     │
     ├─ UpdateStatus() [deployment.go]
     │  └─ Update Stack.status = "Deployed" in PostgreSQL
     │
     ├─ UpdateMetadata() [deployment.go]
     │  └─ Store in Stack.metadata (jsonb):
     │     • L1 chain info, RPC URLs, contract addresses
     │     • L2 chain info, genesis hash, L1 deposit contract
     │     • Bridge setup details
     │
     ├─ SetupBridge() [deployment.go]
     │  └─ Initialize bridge integration
     │
     ├─ [IF local deployment]
     │  │
     │  └─ autoInstallCrossTradeLocal() [deployment.go:672-732]
     │     └─ Auto-install cross-trade module
     │
     └─ [Database transaction committed]

CLIENT NOTIFICATION LOOP (Parallel, Long-Polling):
│
└─ useDeploymentStatus(stackId) [trh-platform-ui/src/hooks/useDeploymentStatus.ts]
   │
   ├─ Poll GetStackStatus() endpoint every N seconds
   │  │
   │  └─ GetStackStatus() [trh-backend/pkg/api/handlers.go]
   │     │
   │     └─ StackRepository.FindByID(stackId) [db/stack_repository.go]
   │        └─ Query Stack table from PostgreSQL
   │           └─ Return: {status, progress, metadata, error}
   │
   └─ Frontend updates UI with:
      • Deployment status (queued → in_progress → completed)
      • Contract addresses, RPC URLs
      • Transaction hashes, block numbers
      • Any errors or warnings
```

### Call Graph Legend

| Symbol | Meaning |
|--------|---------|
| `→` | Synchronous function call (blocking) |
| `└─` | Call chain continues |
| `[ASYNC]` | Execution enters async task queue |
| `[IF condition]` | Conditional branch based on config |
| `[Input/Output]` | Data structure transformation |

### Key Synchronization Points

1. **HTTP Request → Response**: Phase 2
   - Client receives `{stackId, taskId}` immediately
   - Deployment continues asynchronously in background

2. **Task Queue → Worker Pool**: Phase 3
   - TaskManager enqueues task
   - WorkerPool goroutine executes async

3. **Foundry Execution → Output**: Phase 4
   - Binary download (cached if already present)
   - Shell execution with environment variables
   - Wait for `deploy-output.json` file creation

4. **L1 Output → L2 Genesis**: Phase 5
   - Read `deploy-output.json`
   - Extract contract addresses
   - Inject into genesis.json via bytecode patching

5. **Database Commit → Client Notification**: Phase 6
   - Status update transaction committed
   - Client polls and sees updated status
   - Long-polling loop continues until completion

---

## Appendix B: Data Structures & Transformations

### Overview

This appendix documents all major data structures and shows how data transforms as it flows through the system.

### Phase 2→3: HTTP Request to Task Queue

**Input (Phase 2 - HTTP Body)**:
```json
{
  "stackName": "my-stack-1",
  "ethereumRpcUrl": "https://eth-sepolia.g.alchemy.com/v2/...",
  "presetName": "thanos-default",
  "enableFaultProof": true,
  "usePlasma": false,
  "reuseDeployment": false,
  "adminAddress": "0x...",
  "sequencerPrivateKey": "0x...",
  "batcherPrivateKey": "0x...",
  "proposerPrivateKey": "0x...",
  "challengerPrivateKey": "0x..."
}
```

**Processing (Phase 3)**:
```go
// File: trh-backend/pkg/dtos/deploy_request.go
type DeployWithPresetRequest struct {
    StackName            string `json:"stack_name" binding:"required"`
    PresetName           string `json:"preset_name" binding:"required"`
    EthereumRpcUrl       string `json:"ethereum_rpc_url" binding:"required"`
    AdminAddress         string `json:"admin_address" binding:"required"`
    SequencerPrivateKey  string `json:"sequencer_private_key" binding:"required"`
    BatcherPrivateKey    string `json:"batcher_private_key" binding:"required"`
    ProposerPrivateKey   string `json:"proposer_private_key" binding:"required"`
    ChallengerPrivateKey string `json:"challenger_private_key"`
    EnableFaultProof     bool   `json:"enable_fault_proof"`
    UsePlasma            bool   `json:"use_plasma"`
    ReuseDeployment      bool   `json:"reuse_deployment"`
}

type DeployConfig struct {
    StackName            string
    L1ChainID            uint64
    L1RPC                string
    L2ChainID            uint64
    AdminAddress         common.Address
    SequencerPrivateKey  *ecdsa.PrivateKey
    BatcherPrivateKey    *ecdsa.PrivateKey
    ProposerPrivateKey   *ecdsa.PrivateKey
    ChallengerPrivateKey *ecdsa.PrivateKey
    GasPrice             *big.Int
    EnableFaultProof     bool
    UsePlasma            bool
    ReuseDeployment      bool
}
```

**Output (Phase 3 - Task Queue)**:
```go
type Task struct {
    ID            uuid.UUID
    Type          TaskType        // "deploy_l1_contracts", "deploy_aws_infrastructure", etc.
    Stack         *Stack          // Full Stack entity with ID
    Config        *DeployConfig   // Serialized into task.Config JSON
    Status        TaskStatus      // Pending → InProgress → Completed → Failed
    StartedAt     *time.Time
    CompletedAt   *time.Time
    Error         *string
}
```

### Phase 4: L1 Contract Deployment

**Input (Phase 4)**:
```bash
# Environment variables written to .env file
export L1_RPC_URL="https://eth-sepolia.g.alchemy.com/v2/..."
export ADMIN_ADDRESS="0x..."
export PRIVATE_KEY="0x..."
export GAS_PRICE="100000000000"  # Wei
```

**Foundry Script Execution** (`start-deploy.sh`):
```bash
forge script scripts/Deploy.sol:Deploy \
  --fork-url $L1_RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --slow \
  --verify
```

**Output: deploy-output.json**:
```json
{
  "contracts": {
    "SuperchainConfig": "0xAddress1",
    "OptimismPortal": "0xAddress2",
    "L1StandardBridge": "0xAddress3",
    "L1CrossDomainMessenger": "0xAddress4",
    "OptimismPortalProxy": "0xAddress5",
    "SystemConfig": "0xAddress6",
    ... (30+ more contracts)
  },
  "deployment_order": [
    "SuperchainConfig",
    "OptimismPortal",
    "L1StandardBridge",
    ... (in order)
  ],
  "timestamps": {
    "SuperchainConfig": 1713283200,
    "OptimismPortal": 1713283210,
    ... (block timestamps)
  },
  "block_numbers": [
    17500000,
    17500001,
    17500002,
    ... (L1 block numbers)
  ]
}
```

### Phase 5: L2 Genesis Generation

**Input: L1 Contract Addresses**:
```go
// Extracted from deploy-output.json
type L1Contracts struct {
    SuperchainConfig        common.Address
    OptimismPortal          common.Address
    OptimismPortalProxy     common.Address
    L1StandardBridge        common.Address
    L1StandardBridgeProxy   common.Address
    L1CrossDomainMessenger  common.Address
    L1ERC721Bridge          common.Address
    SystemConfig            common.Address
    ProxyAdmin              common.Address
    Create2Deployer         common.Address
    FaultDisputeGame        common.Address
    // ... 30+ more contracts
}
```

**Predeploy Immutable Injection**:

For each predeploy contract (40+ total), immutables are injected:

```go
type ImmutableConfig struct {
    // StandardBridge contracts
    L1StandardBridge           common.Address
    L2StandardBridge           common.Address
    L1CrossDomainMessenger     common.Address
    L2CrossDomainMessenger     common.Address
    
    // Oracle & system
    L1BlockAddress             common.Address
    SystemConfigAddress        common.Address
    
    // Sequencer & operators
    SequencerAddress           common.Address
    BatchInboxAddress          common.Address
    
    // Fee vaults
    L1FeeVaultAddress          common.Address
    BaseFeeVaultAddress        common.Address
    SequencerFeeVaultAddress   common.Address
    
    // Fault proof (if enabled)
    FaultGameFactoryAddress    common.Address
    FaultDisputeGameAddress    common.Address
    
    // Account abstraction
    EntryPointAddress          common.Address
    MultiTokenPaymasterAddress common.Address
    
    // ... 20+ more immutable configurations
}
```

**Bytecode Patching Process**:

For a contract like `L2StandardBridge`:

```
Original bytecode:
  [PUSH32 0x00000000000000000000000000000000] ← placeholder for L1StandardBridge
  [PUSH32 0x00000000000000000000000000000000] ← placeholder for L1CrossDomainMessenger
  [... rest of code ...]

After injection:
  [PUSH32 0xAddress1Address1Address1Address1] ← actual L1StandardBridge address
  [PUSH32 0xAddress2Address2Address2Address2] ← actual L1CrossDomainMessenger address
  [... rest of code ...]
```

The `Deploy()` function (immutables.go) performs this replacement by:
1. Loading contract bytecode
2. Finding hardcoded placeholder offsets
3. Replacing with actual addresses from L1 deployment
4. Validating bytecode hash matches expected value

### Phase 5 Output: genesis.json

**Structure**:
```json
{
  "config": {
    "chainId": 11155420,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "londonBlock": 0,
    "terminalTotalDifficulty": 0,
    "terminalTotalDifficultyPassed": true,
    "optimismConfig": {
      "eip1559Elasticity": 6,
      "eip1559Denominator": 50
    }
  },
  "nonce": "0x0",
  "timestamp": "0x661a1000",
  "extraData": "0x",
  "gasLimit": "0x1c9c380",
  "difficulty": "0x1",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc": {
    "0x4200000000000000000000000000000000000000": {
      "nonce": "0x0",
      "code": "[L2ToL1MessagePasser contract bytecode with injected immutables]",
      "storage": {
        "[slot1]": "[value1]",
        "[slot2]": "[value2]"
      },
      "balance": "0x0"
    },
    "0x4200000000000000000000000000000000000001": {
      "code": "[L2CrossDomainMessenger bytecode]",
      "storage": { ... },
      "balance": "0x0"
    },
    // ... 38+ more predeploys
    "0x4200000000000000000000000000000f0d6d8f": {
      "code": "[SequencerFeeVault bytecode]",
      "storage": {
        "[RECIPIENT_SLOT]": "[feeVaultRecipient]"
      },
      "balance": "0x0"
    }
  }
}
```

**Size**: 20-100 MB depending on:
- Number of predeploys
- Size of each contract bytecode
- Storage initialization data
- Hard fork applied

### Phase 6 Output: StackMetadata (Database)

**Stored in PostgreSQL Stack table (metadata JSONB column)**:
```json
{
  "layer1": {
    "chain_id": 11155111,
    "rpc_url": "https://eth-sepolia.g.alchemy.com/v2/...",
    "explorer_url": "https://sepolia.etherscan.io",
    "contracts": {
      "SuperchainConfig": "0xAddress1",
      "OptimismPortal": "0xAddress2",
      "L1StandardBridge": "0xAddress3",
      "SystemConfig": "0xAddress6",
      "ProxyAdmin": "0xAddressN"
    },
    "block_number": 17500000,
    "timestamp": 1713283200
  },
  "layer2": {
    "chain_id": 11155420,
    "rpc_url": "https://rpc.l2-sepolia.thanos.network",
    "explorer_url": "https://sepolia-explorer.thanos.network",
    "genesis_hash": "0x...",
    "genesis_file_path": "build/genesis.json",
    "rollup_config_path": "build/rollup.json"
  },
  "bridge": {
    "l1_bridge_address": "0xAddress3",
    "l2_bridge_address": "0x4200000000000000000000000000000000000010",
    "status": "active"
  },
  "deployment_timing": {
    "started_at": "2024-04-16T10:30:00Z",
    "completed_at": "2024-04-16T11:45:00Z",
    "total_seconds": 4500
  }
}
```

### Data Transformation Summary Table

| Phase | Input Type | Transformation | Output Type | Location |
|-------|-----------|----------------|------------|----------|
| 2→3 | JSON (HTTP body) | Deserialization + validation | DeployWithPresetRequest | presets.go:PresetDeploy |
| 3 | Preset | Template expansion + variable substitution | DeployConfig | stack_lifecycle.go:CreateThanosStack |
| 3 | DeployConfig | Task serialization | Task (queued) | task_manager.go:AddTask |
| 4 | DeployConfig | Environment variable generation | .env file | thanos_stack.go:DeployL1Contracts |
| 4 | .env + Foundry script | Shell execution | deploy-output.json | start-deploy.sh execution |
| 4→5 | deploy-output.json | Address extraction | L1Contracts struct | deployment.go:executeDeployments |
| 5 | L1Contracts | Bytecode patching | Modified bytecode | immutables.go:Deploy |
| 5 | Modified bytecode + storage | Genesis serialization | genesis.json file | genesis_writer.go:WriteGenesisJSON |
| 5 | DeployConfig + genesis | Config serialization | rollup.json file | rollup/config.go:NewRollupConfig |
| 6 | Stack + outputs | Metadata extraction | StackMetadata (JSONB) | deployment.go:UpdateMetadata |

---

## Appendix C: Known Pitfalls & Improvement Points

### Critical Pitfalls (🔴 High Risk)

#### 1. Blob Fee Spike on Layer 1

**Problem**: L1 can experience rapid spikes in blob fees (EIP-4844 dynamic pricing), causing gas estimation to fail or become wildly inaccurate during contract deployment.

**Impact**: Deployment hangs or fails with "transaction underpriced" or "insufficient balance" errors.

**Current Mitigation**:
- `fix(txmgr): restore real excessBlobGas in suggestGasPriceCaps` (commit d8202223de)
- Configurable cap multiplier and threshold in tx manager
- Treat excessBlobGas as 0 during Sepolia deployments (commit 2a9e294c32)

**What can go wrong**:
- Fee estimation multiplier set too low → reverts
- Multiplier set too high → wasteful gas spending
- Threshold not matched to network → alternates between failure and overpaying

**Recommendation**: Monitor blob base fee trends before deployment; consider queuing if base fee > 100 Gwei.

---

#### 2. Preset Variable Substitution Errors

**Problem**: Environment variables in preset `.env` files are subject to shell expansion. Typos in variable names or missing secrets cause silent failures or incorrect deployments.

**Example**:
```bash
# .env template
export ADMIN_KEY=$ADMIN_PRIVATE_KEY
export WRONG_VAR=$TYPO_VARIABLE    # ← Typo: should be TYPO_VARIABLE_NAME

# If TYPO_VARIABLE_NAME is not set, shell expands to empty string
# Result: WRONG_VAR="" (empty private key)
```

**Impact**: Contract deployment with wrong addresses or zero keys → fund loss or uncontrolled accounts.

**Current Mitigation**: Validation in `CreateThanosStack()` checks required fields.

**What can go wrong**:
- Typo in environment variable name
- Missing secret in CI/CD environment
- Case sensitivity: `$ADMIN_KEY` vs `$admin_key`
- Expansion of special characters (e.g., `$` in password)

**Recommendation**:
1. Pre-validate all required variables before shell-out
2. Use strict mode in bash: `set -u` (fail on undefined variables)
3. Log expanded .env file (masked) for debugging

---

#### 3. Cross-Trade Auto-Install Flag Scope

**Problem**: `flags.CrossTradeAutoInstall` is a global configuration flag, but should only apply to local deployments. If enabled during production deployment, cross-trade auto-installs unexpectedly.

**Location**: `deployment.go:672-732` (`autoInstallCrossTradeLocal()`)

**Impact**: Unexpected module auto-installation, potentially breaking production stacks.

**Current Mitigation**:
```go
// Lines ~675-680
if s.config.Flags.CrossTradeAutoInstall && isLocalDeployment(stack) {
    // only install if local
}
```

**What can go wrong**:
- Flag checked only at deployment time, not at preset creation time
- If flag flipped between preset creation and deployment → unexpected behavior
- No audit trail of flag state at deployment time

**Recommendation**:
1. Store flag in Stack entity (immutable after creation)
2. Validate flag during preset creation, not just deployment
3. Add explicit "auto_install" field to preset definition

---

#### 4. Immutable Injection Bytecode Patching

**Problem**: Immutable offsets are hardcoded per contract version. If Openzeppelin or other dependencies upgrade, bytecode layout shifts → injected addresses end up at wrong offsets → incorrect contract behavior.

**Location**: `op-chain-ops/genesis/immutables/immutables.go`

**Example**: StandardBridge requires injecting `_l1Token` and `_l2Token` at specific bytes32 offsets. If the contract adds a new immutable parameter, all offsets shift.

**Impact**: Genesis.json contains contract bytecode with addresses injected at wrong memory locations → unpredictable contract behavior at runtime.

**Current Mitigation**:
- Bytecode hash validation before patching
- Hardcoded offset constants with version comments

**What can go wrong**:
- Dependency update without updating offset constants
- Multiple contract versions in same deployment
- Offset validation passes but addresses injected incorrectly (rare, but possible if bytecode changed in non-obvious way)

**Recommendation**:
1. Add bytecode hash signatures for each version
2. Create predeploy bytecode patching test suite with before/after validation
3. Use AST analysis to compute offsets dynamically (instead of hardcoding)
4. Pin all dependencies (openzeppelin, optimism) to exact versions

---

#### 5. Hard Fork Compatibility Version Mismatch

**Problem**: Different genesis allocations and storage modifications apply for Ecotone, Fjord, and Granite. If RollupConfig hard fork time doesn't match genesis modifications, L2 execution diverges from expectations.

**Location**: `layer_two.go:PerformUpgradeTxs()`

**Example**: Fjord changes transaction format (RLP encoding). If genesis allocates storage for pre-Fjord format but RollupConfig enables Fjord at block 0, sequencer can't parse first block.

**Impact**: L2 chain halts or processes transactions incorrectly after hard fork block.

**Current Mitigation**:
- Hard fork times specified in DeployConfig
- PerformUpgradeTxs applies changes before genesis is finalized

**What can go wrong**:
- Hard fork time in RollupConfig ≠ actual block number when upgrade applied
- Different clients (op-geth, op-node) interpret hard fork times differently
- Predeploy bytecode updated for Fjord but storage not updated

**Recommendation**:
1. Add hard fork version compatibility matrix to validation checklist
2. Cross-reference RollupConfig hard fork times with genesis block allocations
3. Include hard fork version in deployment log output
4. Test genesis against all supported hard fork versions before finalization

---

### Warnings (🟡 Medium Risk)

#### 1. Large genesis.json File Size

**Problem**: genesis.json can grow to 20-100 MB depending on predeploy bytecode and storage. This causes:
- Slow disk I/O during write/read
- Memory pressure during serialization
- Network bottlenecks if transferred to nodes

**Current Approach**: Serialize entire genesis in memory, then write to disk.

**Recommendation**: Implement streaming JSON writer for genesis generation.

---

#### 2. Create2 Determinism Assumptions

**Problem**: Multiple contracts deployed with CREATE2. If salt calculation changes or collides with existing deployments, address divergence occurs.

**Recommendation**: Validate all CREATE2 salts are unique before deployment; log all salts used.

---

#### 3. Fault Proof State Initialization

**Problem**: Cannon prestate must match L1 oracle data. If prestate is stale or doesn't account for recent L1 state changes, fault proofs fail.

**Recommendation**: Validate prestate freshness before deployment; compare against L1 block height.

---

#### 4. ReuseDeployment Flag Restart Safety

**Problem**: If deployment fails partway through, reusing existing contracts may resume incorrectly if some contracts are partially initialized.

**Current Mitigation**: Resume only checks if `settings.json` exists, not if all contracts are fully initialized.

**Recommendation**: Add state validation to check all existing contracts before resuming.

---

#### 5. Context Cancellation in ExecuteTask

**Problem**: Long-running deployments (4+ hours) may not handle context cancellation gracefully. If user cancels, partial state cleanup may leave database inconsistencies.

**Current Approach**: Context passed through all function calls, but cleanup logic not guaranteed.

**Recommendation**: Add transaction rollback and cleanup handlers for cancellation.

---

### Improvement Opportunities (💡 Low Risk, High Value)

#### 1. Separate Preset Validation Layer

**Current State**: Validation happens in CreateThanosStack(), late in the flow.

**Improvement**: Create PresetValidator that checks:
- All required variables are defined
- Variable values are valid (addresses, keys, URLs)
- Preset is compatible with target network
- Hard fork settings match network state

This runs **before** shell-out, catching errors early.

---

#### 2. Predeploy Bytecode Patching Test Suite

**Current State**: No automated tests for bytecode injection correctness.

**Improvement**: Create test for each predeploy:
- Load original bytecode
- Inject immutables
- Verify addresses are at correct offsets in memory
- Execute getter functions to confirm values readable

---

#### 3. Genesis Validation Schema

**Current State**: Genesis.json is validated only at contract execution time.

**Improvement**: Add JSON schema or custom validator that checks:
- All predeploy addresses in valid range (0x4200...)
- All storage slots are valid
- All bytecode hashes match expected
- genesis.alloc has no duplicates

Run before file write.

---

#### 4. Error Recovery Strategy for Multi-Phase Failure

**Current State**: If Phase 4 fails, entire deployment fails; Phase 5 never starts.

**Improvement**: Implement recovery strategy:
- Save partial outputs (deploy-output.json) even on failure
- Allow retry from specific phase
- Log which phase failed for faster debugging

---

#### 5. Observability: Phase Timings & Logging

**Current State**: Timings only in database (started_at, completed_at).

**Improvement**: Add per-phase instrumentation:
- Log each phase start/end with duration
- Log contract deployment order and individual times
- Log predeploy addresses as injected (with immutable values)
- Log hard fork application with before/after state

Example output:
```
[PHASE_4_START] L1 deployment started
[L1_CONTRACT] SuperchainConfig deployed at 0x... in 5.3s
[L1_CONTRACT] OptimismPortal deployed at 0x... in 4.1s
...
[PHASE_4_END] L1 deployment completed in 45.2s

[PHASE_5_START] L2 genesis generation started
[L2_PREDEPLOY] L2ToL1MessagePasser (0x4200000000000000000000000000000000000016) with L1Bridge=0x...
[L2_PREDEPLOY] L2StandardBridge (0x4200000000000000000000000000000000000010) with L1Bridge=0x..., L1Messenger=0x...
...
[PHASE_5_UPGRADE] Applying Fjord hard fork
[PHASE_5_END] L2 genesis generation completed in 12.5s
```

---

## Appendix D: Verification Checklist

### Accuracy Verification

Use this checklist to verify that this analysis matches the current codebase.

#### Phase Completeness

- [ ] All 6 Phase analyses (PHASE_1_ANALYSIS.md through PHASE_6_ANALYSIS.md) reviewed
- [ ] All HTTP handlers listed in code-reference-table.md exist and match handler signatures
- [ ] All SDK functions in trh-sdk/stacks/thanos exist and match function signatures
- [ ] All Foundry scripts in tokamak-thanos/scripts/deploy exist and match calling convention
- [ ] All op-chain-ops deployer functions exist and match code structure

#### Code Accuracy

- [ ] File paths in code-reference-table.md are correct (verify with `ls -la`)
- [ ] Function line numbers ±5 lines of actual definition
  - Command: `grep -n "func FunctionName" file.go`
- [ ] All function signatures match code exactly:
  - Parameter count
  - Parameter types
  - Return types
- [ ] No invented functions or files (all exist in current codebase)

**Verification commands**:
```bash
# Check if file exists
ls -la /Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go

# Check function line number
grep -n "func deploy(" /Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go

# Verify function signature
sed -n '31,40p' /Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go
```

#### Cross-Reference Consistency

- [ ] Function names consistent across all documents
  - `PresetDeploy` not `PresetDeployment` or `DeployPreset`
  - `CreateThanosStack` not `NewThanosStack` or `ThanosStackCreate`
  - Same spelling in code-reference-table.md and phase analyses

- [ ] File paths consistent (relative from repo root)
  - trh-backend: `/pkg/services/thanos/...` not `pkg/...`
  - tokamak-thanos: `/cmd/...` not `cmd/...`

- [ ] Phase numbering consistent (1-6)
  - Phase 1 = SSO
  - Phase 2 = HTTP request
  - Phase 3 = Backend queuing
  - Phase 4 = L1 deployment
  - Phase 5 = L2 genesis
  - Phase 6 = Result persistence

#### Data Structure Consistency

- [ ] DeployConfig fields match across all phases
- [ ] L1DeployOutput matches expected keys in deploy-output.json
- [ ] ImmutableConfig includes all predeploy addresses needed
- [ ] StackMetadata JSONB schema matches database schema

#### Call Chain Consistency

- [ ] Appendix A call graph matches code execution order
- [ ] No forward references (Phase 4 doesn't call Phase 2 functions)
- [ ] All synchronous calls show sequential execution (→)
- [ ] All async queuing shows correct symbol ([ASYNC])

**Verification by tracing**:
1. Start at PresetDeploy() handler
2. Follow each function call listed in Appendix A
3. Verify actual code matches each step
4. Confirm no skipped steps or extra steps

---

## Appendix E: Diagram Reference & Index

### Diagrams Location

All diagrams are in `docs/analysis/diagrams/` directory (SVG format, generated from Mermaid).

### Available Diagrams

#### 1. System Architecture

**File**: `docs/analysis/diagrams/thanos-system-architecture.svg`

**Purpose**: Show 8-layer system decomposition from Electron client to database

**Layers**:
1. Electron Shell (UI)
2. Next.js Frontend (trh-platform-ui)
3. HTTP API Gateway (trh-backend HTTP handlers)
4. Service Layer (orchestration, TaskManager)
5. SDK Layer (trh-sdk wrappers)
6. Deployment Executors (tokamak-deployer, op-chain-ops)
7. External Services (L1 RPC, S3, EC2)
8. Data Layer (PostgreSQL, JSON files)

**Communication Channels**:
- HTTP between layers 2-3
- Go module imports between layers 3-5
- Shell execution (layer 5 → 6)
- File I/O for deployment outputs

**Use this to**:
- Understand overall system decomposition
- Identify service boundaries
- Understand data flow across layers
- Locate which layer handles which concern

---

#### 2. L1 Deploy Flow

**File**: `docs/analysis/diagrams/l1-deploy-flow.svg`

**Purpose**: Detailed flowchart of L1 contract deployment (Foundry execution)

**Key Stages**:
1. Environment setup (L1 RPC URL, private keys)
2. Contract deployment order decision:
   - System contracts (SuperchainConfig, OptimismPortal)
   - Proxy pattern contracts
   - Conditional branches: EnableFaultProof, UsePlasma, ReuseDeployment
3. Contract deployment execution:
   - Each contract → CREATE2 address generation
   - Fund account if needed
   - Submit deployment transaction
   - Wait for confirmation
4. Output generation (deploy-output.json, rollup.json)

**Conditional Branches**:
- **ReuseDeployment**: If true, skip already-deployed contracts
- **EnableFaultProof**: If true, additionally deploy Cannon prestate
- **UsePlasma**: If true, include Plasma-specific contracts

**Use this to**:
- Understand contract deployment order and why
- See which contracts are optional vs. required
- Understand conditional logic in deployment
- Debug individual contract deployment failures

---

#### 3. L2 Genesis Flow

**File**: `docs/analysis/diagrams/l2-genesis-flow.svg`

**Purpose**: Detailed op-chain-ops Deployer workflow for genesis construction

**Major Steps** (14+):
1. Read L1 deployment artifacts (deploy-output.json)
2. Initialize base genesis (chain ID, timestamp, gas limit)
3. Create ImmutableConfig from L1 addresses
4. Create StorageConfig for initial storage values
5. Deploy transparent proxy contracts (256 @ 0x4200...)
6. For each of 40+ predeploys:
   a. Load bytecode
   b. Apply immutable injection (bytecode patching)
   c. Add to genesis.alloc with storage values
   d. Calculate state root
7. Apply hard fork changes (Ecotone/Fjord/Granite)
8. Generate RollupConfig (rollup.json)
9. Validate genesis (all addresses, storage slots)
10. Write genesis.json (20-100 MB)

**Key Data Structures**:
- ImmutableConfig: 40+ address mappings
- StorageConfig: Storage slot overrides
- core.Genesis: Final structure with all allocs
- RollupConfig: Protocol parameters

**Use this to**:
- Understand L2 genesis construction details
- See predeploy initialization order
- Understand hard fork application
- Debug genesis validation errors

---

#### 4. Call Hierarchy

**File**: `docs/analysis/diagrams/thanos-call-hierarchy.svg`

**Purpose**: Vertical function call nesting showing parent-child relationships

**Top-Level Entry**:
```
PresetDeploy (HTTP handler)
├─ PresetDeployRequest validation
└─ CreateThanosStackFromPreset
   ├─ createStackFromPreset
   │  └─ Preset → DeployConfig expansion
   └─ CreateThanosStack
      ├─ Atomic transaction (create Stack, Deployments, IntegrationEntities)
      └─ startPresetDeployment
         └─ TaskManager.AddTask → WorkerPool.Execute
            └─ deploy() [async execution begins]
               ├─ executeDeployments
               │  ├─ DeployL1Contracts → tokamak-deployer binary
               │  │  └─ DeployContracts → Foundry shell-out
               │  │     ├─ For each contract: Create2Deploy
               │  │     └─ WriteDeployOutput → deploy-output.json
               │  └─ BuildL2Genesis
               │     ├─ NewL2Genesis
               │     ├─ NewL2ImmutableConfig
               │     ├─ NewL2StorageConfig
               │     ├─ For each predeploy: Deploy (bytecode injection)
               │     ├─ PerformUpgradeTxs (hard fork changes)
               │     ├─ NewRollupConfig
               │     └─ WriteGenesisJSON → genesis.json + rollup.json
               ├─ UpdateStatus
               ├─ UpdateMetadata
               └─ SetupBridge
```

**Use this to**:
- Trace function call dependencies
- Understand nesting levels (how deep calls go)
- Identify phase boundaries
- See where async/sync boundaries exist

---

#### 5. Data Flow

**File**: `docs/analysis/diagrams/thanos-data-flow.svg`

**Purpose**: Horizontal transformation flow showing data structure changes

**Main Flow**:
```
HTTP Request (JSON)
    ↓
DeployWithPresetRequest (validated)
    ↓
Preset Template (expanded)
    ↓
DeployConfig (expanded with variables)
    ↓
.env File (environment variables)
    ↓
Foundry Script Execution
    ↓
deploy-output.json (L1 contract addresses)
    ↓
L1Contracts struct (extracted addresses)
    ↓
ImmutableConfig (immutable injections)
    ↓
Bytecode Patching (for 40+ predeploys)
    ↓
core.Genesis struct (all predeploys initialized)
    ↓
RollupConfig (protocol parameters)
    ↓
genesis.json + rollup.json (file output)
    ↓
StackMetadata (JSONB) → PostgreSQL
    ↓
HTTP Response (client notification)
```

**Key Transformations**:
- JSON deserialization
- String interpolation (variable expansion)
- Type conversion (string → Address)
- Bytecode patching (binary modification)
- JSON serialization (to files and database)

**Use this to**:
- Understand data transformations at each stage
- Trace where data gets lost or corrupted
- Debug JSON serialization issues
- Understand dependencies between data structures

---

### How to Use Diagrams in Documentation

| Diagram | Use Case |
|---------|----------|
| System Architecture | Onboard new engineers; explain module boundaries |
| L1 Deploy Flow | Understand contract deployment order; debug deployment failures |
| L2 Genesis Flow | Understand predeploy initialization; debug genesis issues |
| Call Hierarchy | Trace function dependencies; understand call nesting |
| Data Flow | Understand data transformations; debug data corruption |

### Generating Diagrams

Diagrams are generated from Mermaid source files:

```bash
# Example: regenerate System Architecture diagram
cd /Users/theo/workspace_tokamak/tokamak-thanos
npx mermaid docs/diagrams/thanos-system-architecture.mmd -o docs/analysis/diagrams/
```

---

## Appendix F: Related Wiki & Documentation

### Internal Documentation

#### TRH Wiki

- **[wiki/ec2-deploy.md](../../trh-wiki/wiki/ec2-deploy.md)** — AWS EC2 deployment procedures, infrastructure setup, Terraform configuration
- **[wiki/tokamak-thanos-stack.md](../../trh-wiki/wiki/tokamak-thanos-stack.md)** — Infrastructure details, networking, security groups, RDS configuration
- **[wiki/troubleshooting.md](../../trh-wiki/wiki/troubleshooting.md)** — Known issues, common errors, recovery procedures
- **[wiki/index.md](../../trh-wiki/wiki/index.md)** — Wiki index and navigation

#### This Analysis Document

- **[code-reference-table.md](./code-reference-table.md)** — Complete function-by-function mapping with file paths and line numbers
- **[PHASE_1_ANALYSIS.md](./PHASE_1_ANALYSIS.md)** — Electron SSO authentication workflow
- **[PHASE_2_ANALYSIS.md](./PHASE_2_ANALYSIS.md)** — Frontend preset wizard and HTTP request submission
- **[PHASE_3_ANALYSIS.md](./PHASE_3_ANALYSIS.md)** — Backend queuing, stack creation, task orchestration
- **[PHASE_4_ANALYSIS.md](./PHASE_4_ANALYSIS.md)** — L1 contract deployment via Foundry
- **[PHASE_5_ANALYSIS.md](./PHASE_5_ANALYSIS.md)** — L2 genesis generation and hard fork application
- **[PHASE_6_ANALYSIS.md](./PHASE_6_ANALYSIS.md)** — Result persistence and client notification

#### Design Documentation

- **[Design Spec](../superpowers/specs/2026-04-16-thanos-deployer-analysis-design.md)** — Original project requirements and success criteria

---

### Code Repositories

#### Core Repositories

1. **trh-backend** — Go backend with HTTP handlers, task manager, database layer
   - HTTP API: `/pkg/api/handlers/thanos/`
   - Services: `/pkg/services/thanos/`
   - Task management: `/pkg/task/`

2. **trh-sdk** — Go SDK for stack deployment orchestration
   - Thanos stack wrapper: `/pkg/stacks/thanos/deploy_contracts.go`

3. **tokamak-thanos** — Foundry scripts and deployer binary
   - Deployment scripts: `/scripts/deploy/`
   - Deployer binary: `/cmd/`

4. **op-chain-ops** — Optimism stack genesis builder
   - Genesis layer two: `/genesis/layer_two.go`
   - Immutable injection: `/genesis/immutables/`
   - Rollup config: `/rollup/config.go`

5. **trh-platform-ui** — Next.js frontend with Electron integration
   - Preset wizard: `/src/hooks/usePresetWizard.ts`
   - Deployment form: `/src/components/DeploymentForm.tsx`
   - Status polling: `/src/hooks/useDeploymentStatus.ts`

---

### External References

#### Optimism OP Stack

- **[OP Stack Documentation](https://docs.optimism.io/)** — Official documentation for OP Stack deployment, configuration, and operations
- **[op-chain-ops Repo](https://github.com/ethereum-optimism/optimism/tree/develop/op-chain-ops)** — Optimism's chain operations tooling (genesis generation, rollup config)

#### Development Tools

- **[Foundry Book](https://book.getfoundry.sh/)** — Foundry smart contract development and deployment framework
- **[Solidity Docs](https://docs.soliditylang.org/)** — Solidity programming language reference

#### Standards & Protocols

- **[EIP-1967: Proxy Storage Slots](https://eips.ethereum.org/EIPS/eip-1967)** — Standard for upgradeable proxy implementation
- **[EIP-4844: Proto-Danksharding](https://eips.ethereum.org/EIPS/eip-4844)** — Blob fee mechanism (relevant to pitfalls)

---

### Cross-References by Role

#### For Deployment Engineers

- Start with: **System Architecture diagram** (Appendix E)
- Then read: **L1 Deploy Flow** and **L2 Genesis Flow** (Appendix E)
- Deep dive: **PHASE_4_ANALYSIS.md** and **PHASE_5_ANALYSIS.md**
- Troubleshooting: **Appendix C (Known Pitfalls)** and **wiki/troubleshooting.md**

#### For Backend Developers

- Start with: **PHASE_3_ANALYSIS.md** and **Call Hierarchy diagram**
- Reference: **code-reference-table.md** (Phase 3 section)
- Implementation: **trh-backend** repository

#### For Smart Contract Developers

- Note: This analysis does NOT cover contract implementations
- For contract details: See individual contract repos (optimism, tokamak)
- For bytecode injection: **PHASE_5_ANALYSIS.md** → "Immutable Injection" section

#### For DevOps/Infrastructure

- Start with: **wiki/ec2-deploy.md** and **wiki/tokamak-thanos-stack.md**
- Then read: **PHASE_6_ANALYSIS.md** (result persistence, metadata storage)

---

## References

### Phase Analysis Documents
- PHASE_1_ANALYSIS.md — Electron SSO authentication
- PHASE_2_ANALYSIS.md — Frontend to backend HTTP flow
- PHASE_3_ANALYSIS.md — Backend queuing and orchestration
- PHASE_4_ANALYSIS.md — L1 contract deployment
- PHASE_5_ANALYSIS.md — L2 genesis generation
- PHASE_6_ANALYSIS.md — Result persistence and notification

### Code Reference
- code-reference-table.md — Complete function mapping with line numbers

### Diagrams
- docs/analysis/diagrams/thanos-system-architecture.svg
- docs/analysis/diagrams/l1-deploy-flow.svg
- docs/analysis/diagrams/l2-genesis-flow.svg
- docs/analysis/diagrams/thanos-call-hierarchy.svg
- docs/analysis/diagrams/thanos-data-flow.svg

### Design & Requirements
- docs/superpowers/specs/2026-04-16-thanos-deployer-analysis-design.md

### External Documentation
- TRH Wiki (../../trh-wiki/wiki/)
- OP Stack Documentation (https://docs.optimism.io/)
- Foundry Book (https://book.getfoundry.sh/)

---

## Document Info

**Created**: 2026-04-16
**Status**: Final
**Appendices**: A-F (6 total)
**Total Sections**: 16
**Cross-References**: 50+
**Diagrams**: 5
**Pitfalls Identified**: 5 critical + 5 warnings
**Data Structures Documented**: 15+
**Functions Traced**: 30+

---

*This document synthesizes all prior analysis work into a single comprehensive reference. For detailed code information, see the phase analyses and code-reference-table.md. For infrastructure details, see the TRH wiki. For troubleshooting, see Appendix C (Known Pitfalls) and wiki/troubleshooting.md.*
