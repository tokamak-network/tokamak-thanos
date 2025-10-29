# Challenger System Architecture

> **Complete System Configuration and Fault Proof Challenger System Detailed Guide**

---

## 📚 Table of Contents

1. [Overview](#overview)
2. [Overall System Configuration](#overall-system-configuration)
3. [Component Detailed Description](#component-detailed-description)
4. [GameType-Specific Architecture](#gametype-specific-architecture)
5. [Data Flow](#data-flow)
6. [Deployment Architecture](#deployment-architecture)
7. [Security and Independence](#security-and-independence)

---

## Overview

The optimism-based L2 Rollup system verifies L2 state correctness through a **Fault Proof System**. This document describes the architecture of the entire system and especially the configuration of the **Challenger System** in detail.

### Key Features

- ✅ **Multi-GameType Support**: Cannon (MIPS), Asterisc (RISC-V), AsteriscKona (RISC-V + Rust)
- ✅ **Independent Challenger Stack**: Verification system completely separated from Sequencer
- ✅ **Modular Architecture**: Selectively deploy VMs and Executors by GameType
- ✅ **Reproducible Builds**: Deterministic builds via Docker

---

## Overall System Configuration

### 1. High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                         L1 (Ethereum)                                         │
│  ┌────────────────┐  ┌──────────────────┐   ┌────────────────────────────┐   │
│  │ Batcher Inbox  │  │ OptimismPortal2  │   │ DisputeGameFactory         │   │
│  │ (Batch Data)   │  │ (User Deposits/  │   │  (Fault Proof Games)        │   │
│  │                │  │  Withdrawals)    │   │ - Output Verification       │   │
│  └────────────────┘  └─────────┬────────┘   │ - GameType 0,1,2,3...      │   │
│         ▲                  ▲   │   │        └────────────────────────────┘   │
│         │                  │   │   │                   ▲           ▲         │
│         │                  │   │   └────References────┘           │         │
│         │                  │   │     (respectedGameType validation)│         │
└─────────┼──────────────────┼───┼───────────────────────────────────┼──────────┘
          │                  │   │                                   │
          │ Batch Submit L1→L2 │ L2→L1                      Output/Challenge
          │              Deposit│ Withdrawal                 Submit/Dispute
          │                  │   │                                   │
┌─────────┴─────────┐     [Users]      ┌─────────────────────────────┴─────┐
│ Sequencer Stack   │                   │  op-proposer / op-challenger      │
│                   │                   │                                   │
├───────────────────┤                   │  ┌──────────────┐  ┌────────────┐ │
│ sequencer-l2      │                   │  │ op-proposer  │  │ Challenger │ │
│ (L2 geth)         │                   │  │ (Game Create)│  │   Stack    │ │
├───────────────────┤                   │  └──────────────┘  │(Independent│ │
│ sequencer-op-node │                   │                    │     )⭐    │ │
│ (Sequencer Mode)  │                   │                    └────────────┘ │
├───────────────────┤                   └───────────────────────────────────┘
│ op-batcher        │                          │                      │
│ (Batch Submit)    │                          │                      │
└───────────────────┘                   DisputeGameFactory    challenger-l2
                                            .create()        challenger-op-node
                                         (Create New Game)    op-challenger
```

**Main Component Roles**:

| Component | Role | Users | Interaction |
|---------|------|--------|---------|
| **Batcher Inbox** | Store L2 transaction batch data | op-batcher | Submit batch data |
| **OptimismPortal2** | L1↔L2 message/asset transfer | General users | Deposit (`depositTransaction`), Prove withdrawal (`proveWithdrawalTransaction`), Finalize withdrawal (`finalizeWithdrawalTransaction`) |
| **DisputeGameFactory** | Fault Proof game management | op-proposer, op-challenger | Create game (`create`), Participate in game (`move`, `step`) |

**Key Points**:
- ✅ **OptimismPortal2 ← DisputeGameFactory**: OptimismPortal2 relies on `respectedGameType` games for withdrawal validation
- ✅ **op-proposer → DisputeGameFactory**: Creates new games via `.create()` (L2 Output proposal)
- ✅ **op-challenger → DisputeGameFactory**: Participates in games to challenge incorrect outputs
- ✅ **Users → OptimismPortal2**: Deposits/withdrawals (indirectly connected to DisputeGameFactory via game validation)

### 2. System Layer Structure

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ op-batcher   │  │ op-proposer  │  │ op-challenger    │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │
┌─────────────────────────────────────────────────────────────┐
│                    Rollup Protocol Layer                     │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              op-node (Rollup Consensus)                │ │
│  │  - Derivation: L1 batches → L2 blocks                 │ │
│  │  - Sync: State sync with other nodes                  │ │
│  │  - P2P: Block propagation and receipt                 │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │
┌─────────────────────────────────────────────────────────────┐
│                    Execution Layer                           │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              L2 geth (EVM Execution)                   │ │
│  │  - Transaction execution                               │ │
│  │  - State storage (independent DB) ⭐                   │ │
│  │  - RPC service                                         │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                          ▲
                          │
┌─────────────────────────────────────────────────────────────┐
│                    Data Availability Layer                   │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                L1 Ethereum                             │ │
│  │  - Batch Data (calldata/blobs)                        │ │
│  │  - Output Roots                                        │ │
│  │  - Dispute Game State                                  │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## Component Detailed Description

### L1 Layer

#### 1. **Batcher Inbox (Managed by SystemConfig)**
- **Role**: Storage address for L2 transaction batch data
- **Data**: Compressed L2 block data (calldata or blobs)
- **Access**: op-node reads transaction data sent to the Batcher Inbox address on L1 and reconstructs L2 blocks
- **Note**: OptimismPortal is for L1↔L2 message passing, different from Batcher Inbox

#### 2. **OptimismPortal2**
- **Role**: Message and asset transfer between L1 and L2 (bridge)
- **Users**: General users (depositors and withdrawers)
- **Main Functions**:
  - `depositTransaction()`: L1 → L2 asset deposit (ETH/ERC20)
  - `proveWithdrawalTransaction(_disputeGameIndex)`: L2 → L1 withdrawal proof
  - `finalizeWithdrawalTransaction()`: Complete withdrawal (after proof and delay period)

- **Relationship with DisputeGameFactory** ⭐:
  ```solidity
  // Inside OptimismPortal2 code
  GameType public respectedGameType;  // Trusted GameType (e.g., 3)

  function proveWithdrawalTransaction(..., uint256 _disputeGameIndex, ...) {
      // Get game from DisputeGameFactory
      (GameType gameType,, IDisputeGame gameProxy) =
          disputeGameFactory.gameAtIndex(_disputeGameIndex);

      // Only trust respectedGameType
      require(gameType.raw() == respectedGameType.raw());

      // Cannot use CHALLENGER_WINS games
      require(gameProxy.status() != GameStatus.CHALLENGER_WINS);

      // Use game's rootClaim for withdrawal proof
      Claim outputRoot = gameProxy.rootClaim();
  }
  ```

- **Key Point**: OptimismPortal2 **relies on** DisputeGameFactory for withdrawal validation
  - ✅ Users deposit/withdraw via OptimismPortal2
  - ✅ OptimismPortal2 trusts only `respectedGameType` games from DisputeGameFactory
  - ✅ For withdrawals, OptimismPortal2 validates against game results (rootClaim, status)

#### 3. **Output Oracle (DisputeGameFactory-based)**
- **Role**: L2 state root submission and verification
- **Interval**: Proposer submits periodically (e.g., every 30 seconds)
- **Verification**: Fault Proof games via DisputeGameFactory
- **Note**: Modern system uses DisputeGameFactory instead of L2OutputOracle

#### 4. **Dispute Games (DisputeGameFactory)**
- **Role**: Fault Proof game management
- **GameType**: Supports 0, 1, 2, 3, 254, 255
- **Implementation**: Different VMs for each GameType (MIPS.sol, RISCV.sol)
- **Game Creation**: `create(GameType, Claim, extraData)`
- **Game Verification**: Implementation contracts per GameType (`gameImpls` mapping)

### Sequencer Stack

#### 1. **sequencer-l2 (L2 geth)**
```yaml
Role: L2 transaction execution and block creation
Mode: Sequencer (Leader)
DB: sequencer_l2_data (independent volume)
Port: 9545 (RPC)
```

#### 2. **sequencer-op-node**
```yaml
Role: Rollup consensus and L1 sync
Mode: Sequencer (Leader)
Connections:
  - L1 (http://l1:8545)
  - sequencer-l2 (Engine API)
Port: 7545 (RPC)
```

#### 3. **op-batcher**
```yaml
Role: Submit L2 blocks to L1 as batches
Connections:
  - L1 (batch submission)
  - sequencer-op-node (fetch L2 blocks)
Key: BATCHER_PRIVATE_KEY
```

#### 4. **op-proposer**
```yaml
Role: Submit L2 Output to DisputeGameFactory
Connections:
  - L1 DisputeGameFactory (game creation)
  - sequencer-op-node (fetch state roots)
Key: PROPOSER_PRIVATE_KEY
Interval: PROPOSAL_INTERVAL (e.g., 30s)
Method: Create new game via DisputeGameFactory.create()
```

### Challenger Stack (Independent!) ⭐

#### 1. **challenger-l2 (L2 geth)**
```yaml
Role: Independent L2 execution verification
Mode: Follower (read-only)
DB: challenger_l2_data (independent volume) ⭐
Port: 9546 (RPC)
Important: Completely separate DB from Sequencer!
```

**Why must it be independent?**
- Prepare for cases where Sequencer is malicious
- Independently reconstruct L2 from L1
- True verification (trustless)

#### 2. **challenger-op-node**
```yaml
Role: Independent rollup consensus
Mode: Follower (verification only) ⭐
Connections:
  - L1 (http://l1:8545)
  - challenger-l2 (Engine API)
Port: 7546 (RPC)
Important: Read batches directly from L1 for reconstruction!
```

#### 3. **op-challenger**
```yaml
Role: Execute challenges against incorrect outputs
Connections:
  - L1 (game participation)
  - challenger-op-node (verified state)
  - challenger-l2 (independent execution results)
Key: CHALLENGER_PRIVATE_KEY
TraceType: cannon, asterisc, asterisc-kona
```

**op-challenger Operation**:
1. Monitor Dispute Games on L1
2. Detect incorrect outputs
3. Find first mismatch via binary search
4. Execute VM (Cannon/Asterisc) to generate proof
5. Submit proof to L1

---

## GameType-Specific Architecture

### GameType 0/1: Cannon (MIPS VM)

```
┌─────────────────────────────────────────────────────────┐
│                  GameType 0/1 (Cannon)                  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  On-Chain (L1):                                         │
│  ┌────────────────────────────────────────────────┐    │
│  │           MIPS.sol (MIPS VM)                   │    │
│  │  - Single-step instruction verification        │    │
│  │  - Memory merkle proof verification            │    │
│  │  - PreimageOracle support                      │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Off-Chain (Challenger):                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  Cannon VM (MIPS Emulator)                     │    │
│  │  - Execute op-program (Go)                     │    │
│  │  - Generate trace (each instruction)           │    │
│  │  - Prestate: prestate.json                     │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Build:                                                  │
│  ┌────────────────────────────────────────────────┐    │
│  │  make reproducible-prestate                    │    │
│  │  → Generate Linux binary via Docker            │    │
│  │  → cannon binary                               │    │
│  │  → op-program binary                           │    │
│  │  → prestate-proof.json                         │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  File Structure:                                        │
│  /cannon/bin/cannon             ← MIPS VM              │
│  /op-program/bin/op-program     ← Server (Go)          │
│  /op-program/bin/prestate.json  ← Prestate             │
└─────────────────────────────────────────────────────────┘
```

### GameType 2: Asterisc (RISC-V VM)

```
┌─────────────────────────────────────────────────────────┐
│                  GameType 2 (Asterisc)                  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  On-Chain (L1):                                         │
│  ┌────────────────────────────────────────────────┐    │
│  │          RISCV.sol (RISC-V VM) ⭐               │    │
│  │  - RV64GC instruction verification             │    │
│  │  - Memory merkle proof verification            │    │
│  │  - PreimageOracle support                      │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Off-Chain (Challenger):                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  Asterisc VM (RISC-V Emulator)                 │    │
│  │  - Execute op-program (Go, RISC-V compiled)   │    │
│  │  - Generate trace (each instruction)           │    │
│  │  - Prestate: prestate-proof-rv64.json          │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Build:                                                  │
│  ┌────────────────────────────────────────────────┐    │
│  │  Clone Asterisc source:                        │    │
│  │  git clone github.com/ethereum-optimism/       │    │
│  │             asterisc                           │    │
│  │                                                 │    │
│  │  make reproducible-prestate                    │    │
│  │  → Generate Linux binary via Docker            │    │
│  │  → asterisc binary (RISC-V VM)                 │    │
│  │  → op-program binary (RISC-V target)          │    │
│  │  → prestate-proof.json (.stateHash)            │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  File Structure:                                        │
│  /asterisc/bin/asterisc                ← RISC-V VM     │
│  /op-program/bin/op-program            ← Server (Go)   │
│  /asterisc/bin/prestate-proof.json     ← Prestate      │
│  /asterisc/bin/prestate.json → prestate-proof.json     │
│                                   (symbolic link)       │
│                                                          │
│  Features:                                               │
│  - RISC-V architecture (replaces Cannon's MIPS)       │
│  - More modern ISA                                     │
│  - Uses same op-program (different target)            │
└─────────────────────────────────────────────────────────┘
```

### GameType 3: AsteriscKona (RISC-V + Rust) 🆕

```
┌─────────────────────────────────────────────────────────┐
│              GameType 3 (AsteriscKona) 🆕               │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  On-Chain (L1):                                         │
│  ┌────────────────────────────────────────────────┐    │
│  │     RISCV.sol (Same as GameType 2!) ⭐⭐        │    │
│  │  - Reuses same implementation as GameType 2    │    │
│  │  - 100% identical on-chain verification logic  │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Off-Chain (Challenger):                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  Asterisc VM (RISC-V Emulator)                 │    │
│  │  - Execute kona-client (Rust!) ⭐⭐            │    │
│  │  - Generate trace (each instruction)           │    │
│  │  - Prestate: prestate-kona.json                │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Build:                                                  │
│  ┌────────────────────────────────────────────────┐    │
│  │  Clone Kona source:                            │    │
│  │  git clone github.com/op-rs/kona               │    │
│  │                                                 │    │
│  │  Docker Build (automated):                     │    │
│  │  docker run ghcr.io/op-rs/kona/              │    │
│  │    asterisc-builder:0.3.0                     │    │
│  │    cargo build -Zbuild-std=core,alloc         │    │
│  │      -p kona-client                            │    │
│  │      --profile release-client-lto             │    │
│  │  → target/riscv64imac-unknown-none-elf/       │    │
│  │            release-client-lto/kona-client     │    │
│  │                                                 │    │
│  │  Prestate Generation (inside Docker):          │    │
│  │  docker run golang:1.22-bookworm              │    │
│  │    ./asterisc load-elf                         │    │
│  │      --path /kona/.../kona-client             │    │
│  │      --out prestate-kona.json                 │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  File Structure (simplified):                           │
│  /asterisc/bin/asterisc             ← RISC-V VM        │
│  /bin/kona-client                   ← Server (Rust) ⭐ │
│  /op-program/bin/prestate-kona.json ← Prestate         │
│                                                          │
│  Differences (vs GameType 2):                          │
│  - On-chain VM: Identical (RISCV.sol)                 │
│  - Server: kona-client (Rust) vs op-program (Go)      │
│  - Benefits: ~80% size reduction, ZK proof integration │
│                                                          │
│  Fallback:                                              │
│  - If kona prestate unavailable, use Asterisc prestate │
│    (compatible due to same RISCV.sol)                  │
└─────────────────────────────────────────────────────────┘
```

---

## 📂 GameType 2 vs GameType 3 Detailed Folder Structure

### GameType 2 (Asterisc - Go-based)

```
Project Directory Structure:
tokamak-projects/                        # Working directory
│
└── tokamak-thanos/                      # Project root (this repo)
    │
    ├── op-program/                      # ← GameType 2 source (Go, internal)
    │   ├── client/                      # op-program client
    │   │   ├── main.go                  # Written in Go
    │   │   ├── driver.go
    │   │   └── ...
    │   │
    │   ├── host/                        # op-program host
    │   │
    │   └── bin/
    │       ├── op-program               # ← GameType 2 binary (Go build)
    │       └── prestate.json            # Cannon prestate
    │
    └── asterisc/                        # ← VM (GameType 2, 3 common)
        └── bin/
            ├── asterisc                 # RISC-V VM executable (12MB)
            └── prestate-proof.json      # ← GameType 2 prestate (63MB)
```

**Build Method**:
```bash
# Build Go-based op-program
cd tokamak-thanos/op-program
go build -o bin/op-program ./client

# Build Asterisc VM (Docker)
cd tokamak-thanos/asterisc
make reproducible-prestate
```

---

### GameType 3 (AsteriscKona - Rust-based)

```
Project Directory Structure:
tokamak-projects/                        # Working directory
│
├── kona/                                # ← GameType 3 source (separate repo!)
│   │                                    # GitHub: https://github.com/op-rs/kona
│   │                                    # Auto-cloned by deploy-modular.sh
│   │
│   ├── bin/
│   │   ├── client/                      # ← kona-client source (Rust)
│   │   │   ├── src/
│   │   │   │   ├── main.rs            # Written in Rust ⭐
│   │   │   │   └── lib.rs
│   │   │   └── Cargo.toml             # Rust package config
│   │   │
│   │   ├── host/                        # kona-host (runtime env)
│   │   ├── node/                        # kona-node
│   │   └── ...
│   │
│   ├── crates/                          # Rust crates
│   │   ├── executor/
│   │   ├── preimage/
│   │   └── ...
│   │
│   └── target/                          # Rust build artifacts
│       └── riscv64imac-unknown-none-elf/  # ← RISC-V bare metal target
│           └── release-client-lto/       # ← LTO optimization profile
│               └── kona-client           # ← RISC-V ELF (for prestate + execution)
│
└── tokamak-thanos/                      # Project root
    ├── bin/                             # ← Created on GameType 3 deploy
    │   └── kona-client                  # ← GameType 3 binary (copied from kona, persisted)
    │
    ├── op-program/
    │   └── bin/
    │       ├── prestate.json            # Cannon
    │       └── prestate-kona.json       # ← GameType 3 prestate (created on deploy, persisted)
    │
    └── asterisc/                        # ← VM (GameType 2, 3 common)
        └── bin/
            ├── asterisc                 # RISC-V VM (same as GameType 2!)
            └── prestate-proof.json      # GameType 2 prestate
```


---

### Key Differences Summary

| Item | GameType 2 (Asterisc) | GameType 3 (AsteriscKona) |
|------|----------------------|--------------------------|
| **Program Language** | Go | Rust ⭐ |
| **Source Code Location** | `tokamak-thanos/op-program/` | `../kona/` (separate repo) ⭐ |
| **Built Binary** | `tokamak-thanos/op-program/bin/op-program` | `tokamak-thanos/bin/kona-client` |
| **Build Tool** | Go compiler | Cargo (Rust) |
| **VM Executable** | `tokamak-thanos/asterisc/bin/asterisc` | `tokamak-thanos/asterisc/bin/asterisc` (same!) ✅ |
| **Prestate Location** | `tokamak-thanos/asterisc/bin/prestate-proof.json` | `tokamak-thanos/op-program/bin/prestate-kona.json` |
| **Binary Size** | ~100MB | ~20MB (~80% reduction) ⭐ |

**Location Relationship**:
```
tokamak-projects/          # Working directory
├── kona/                  # GameType 3 source (separate repo)
└── tokamak-thanos/        # Project root
    ├── op-program/        # GameType 2 source
    ├── bin/kona-client    # GameType 3 binary (copied from kona)
    └── asterisc/bin/      # VM (common)
```

**Automatic Deployment**:
```bash
# deploy-modular.sh handles automatically
cd tokamak-thanos
./op-challenger/scripts/deploy-modular.sh --dg-type 3

# Files created after deployment:
# ✅ bin/kona-client (copied from kona, persisted)
# ✅ op-program/bin/prestate-kona.json (generated by asterisc, persisted)
# ✅ ../kona/ (repo cloned, persisted)
```

**Before/After Deployment Comparison**:
```
Before Deployment:              After Deployment:
tokamak-thanos/                 tokamak-thanos/
├── op-program/                 ├── bin/                    ← Created
│   └── bin/                    │   └── kona-client         ← Created
│       └── prestate.json       ├── op-program/
└── asterisc/                   │   └── bin/
    └── bin/                    │       ├── prestate.json
        └── asterisc            │       └── prestate-kona.json  ← Created
                                └── asterisc/
                                    └── bin/
                                        └── asterisc

kona/ does not exist            ../kona/                    ← Cloned
                                ├── bin/client/
                                └── target/riscv64imac-unknown-none-elf/
                                    └── release-client-lto/
                                        └── kona-client
```

---

### GameType Comparison Table

| Component | GameType 0/1 (Cannon) | GameType 2 (Asterisc) | GameType 3 (AsteriscKona) |
|----------|---------------------|---------------------|-------------------------|
| **On-chain VM** | MIPS.sol | RISCV.sol | RISCV.sol (same) ⭐ |
| **ISA** | MIPS | RISC-V (RV64GC) | RISC-V (RV64GC) |
| **Off-chain VM** | cannon | asterisc | asterisc (same) |
| **Server** | op-program (Go) | op-program (Go) | kona-client (Rust) ⭐ |
| **Server Language** | Go | Go | Rust |
| **Prestate Field** | `.pre` | `.stateHash` | `.stateHash` |
| **Build Tools** | Go + Docker | Go + Docker | Rust + Cargo + Docker |
| **Binary Size** | ~100MB | ~100MB | ~20MB (~80% reduction) ⭐ |
| **Advantages** | Stability, Go ecosystem | Modern ISA, RISC-V standard | Lightweight, ZK integration ready |
| **Status** | ✅ Production | ✅ Production | 🆕 New Integration |

---

## Data Flow

### 1. Happy Path (Normal Flow)

```
┌─────────────────────────────────────────────────────────────┐
│                   1. L2 Transaction Execution                │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        User Transaction → Sequencer L2 → Block Creation
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              2. Submit Batch Data to L1 (op-batcher)         │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
              Sequencer op-node → Collect block data
                          │
                          ▼
              op-batcher → Submit to L1 BatcherInbox
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│           3. Submit L2 Output (op-proposer)                  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
              op-proposer → Submit to DisputeGameFactory
                          (Create new game: proposed Output)
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│           4. Challenger Independent Verification (Parallel)  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
┌──────────────────┐              ┌──────────────────┐
│ Read batches     │              │ Read submitted   │
│ from L1          │              │ Output           │
│ (challenger-     │              │                  │
│  op-node)        │              │                  │
└──────────────────┘              └──────────────────┘
        │                                   │
        ▼                                   ▼
┌──────────────────┐              ┌──────────────────┐
│ Reconstruct L2   │              │ Compare with     │
│ (challenger-l2)  │              │ locally computed │
│                  │              │ Output           │
└──────────────────┘              └──────────────────┘
        │                                   │
        └─────────────────┬─────────────────┘
                          │
                          ▼
                    ┌─────────────┐
                    │  Match?      │
                    └─────────────┘
                      │         │
              Match ✅ │         │ Mismatch ❌
                      │         │
                      ▼         ▼
                  (Normal)   Start Challenge!
```

### 2. Challenge Flow (Dispute)

```
┌─────────────────────────────────────────────────────────────┐
│              1. Detect Incorrect Output (op-challenger)      │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        Challenger's local state ≠ L1 submitted Output
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   2. Start Dispute Game                      │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        L1 DisputeGameFactory.create()
          → Create new game instance
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              3. Binary Search (Bisection)                    │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────────┐
        │  Output Root (full block range)     │
        │  ↓                                   │
        │  Claim at depth 1 (half)            │
        │  ↓                                   │
        │  Claim at depth 2 (1/4)             │
        │  ↓                                   │
        │  ...                                 │
        │  ↓                                   │
        │  Single Instruction (final depth)   │
        └─────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│         4. VM Execution and Proof Generation (Single Step)   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        Different VM based on GameType:
        ┌───────────────────────────────────┐
        │ GameType 0/1: cannon (MIPS)       │
        │ GameType 2:   asterisc (RISC-V)   │
        │ GameType 3:   asterisc + kona     │
        └───────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────────┐
        │  1. VM executes that instruction    │
        │  2. Prepare pre-state               │
        │  3. Compute post-state              │
        │  4. Generate memory merkle proof    │
        │  5. Prepare preimage data           │
        └─────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              5. Submit Proof to L1 (step())                  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        L1 FaultDisputeGame.step()
          → Verify in VM (MIPS.sol / RISCV.sol)
          → Compare with on-chain execution result
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   6. Resolve Game (resolve)                  │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────┴─────────────────┐
        │                                   │
        ▼                                   ▼
  Challenger Wins ✅             Proposer Wins ❌
  (Remove incorrect Output)      (Keep Output)
```

### 3. Challenger Independent Verification Details

```
┌─────────────────────────────────────────────────────────────┐
│              Challenger Independent Verification Process     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Step 1: Read Batch Data from L1                            │
│  ┌────────────────────────────────────────────────────┐    │
│  │  challenger-op-node                                │    │
│  │  → Monitor L1 BatcherInbox                         │    │
│  │  → Download batch data                             │    │
│  │  → Decompress                                      │    │
│  │  → Restore L2 transactions                         │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  Step 2: Independent L2 Reconstruction                      │
│  ┌────────────────────────────────────────────────────┐    │
│  │  challenger-l2 (independent geth DB!)              │    │
│  │  → Execute restored transactions in order          │    │
│  │  → Store in independent state DB                   │    │
│  │  → Compute Output Root                             │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  Step 3: Compare Outputs                                    │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Locally computed Output Root                      │    │
│  │         vs                                          │    │
│  │  L1 submitted Output Root (DisputeGameFactory)    │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  Step 4: On Mismatch Detection                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │  1. Create Dispute Game                            │    │
│  │  2. Find first mismatched instruction via binary   │    │
│  │     search                                          │    │
│  │  3. Execute VM to generate proof                   │    │
│  │  4. Submit proof to L1                             │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
│  Important: 0% sharing with Sequencer components! ⭐⭐      │
│  - Independent L2 geth                                      │
│  - Independent op-node                                      │
│  - Independent DB                                           │
│  → True verification (trustless)                            │
└─────────────────────────────────────────────────────────────┘
```

---

## Deployment Architecture

### Docker Compose Structure

**Conceptual Structure** (actual configuration is more complex):

```yaml
services:
  # L1 Layer
  l1:
    ports: ["8545:8545"]
    volumes: [l1_data:/db]

  # ============================================================
  # Sequencer Stack
  # ============================================================
  sequencer-l2:
    ports: ["9545:8545"]              # Sequencer L2 RPC
    volumes: [sequencer_l2_data:/db]  # Sequencer DB
    environment:
      ROLLUP_CLIENT_HTTP: http://sequencer-op-node:8545

  sequencer-op-node:
    ports: ["7545:8545"]              # Sequencer op-node RPC
    environment:
      OP_NODE_L1_ETH_RPC: http://l1:8545
      OP_NODE_L2_ENGINE_RPC: http://sequencer-l2:8551
      OP_NODE_SEQUENCER_ENABLED: true  # ← Sequencer mode

  op-batcher:
    environment:
      OP_BATCHER_ROLLUP_RPC: http://sequencer-op-node:8545

  op-proposer:
    environment:
      OP_PROPOSER_ROLLUP_RPC: http://sequencer-op-node:8545

  # ============================================================
  # Challenger Stack (Completely Independent!) ⭐
  # ============================================================
  challenger-l2:
    ports: ["9546:8545"]                  # Challenger L2 RPC (different port!)
    volumes: [challenger_l2_data:/db]    # ← Independent DB ⭐⭐
    environment:
      ROLLUP_CLIENT_HTTP: http://challenger-op-node:8545

  challenger-op-node:
    ports: ["7546:8545"]                  # Challenger op-node RPC (different port!)
    environment:
      OP_NODE_L1_ETH_RPC: http://l1:8545  # ← Read directly from L1
      OP_NODE_L2_ENGINE_RPC: http://challenger-l2:8551
      OP_NODE_SEQUENCER_ENABLED: false    # ← Follower mode ⭐

  op-challenger:
    volumes:
      # Mount all GameType VM binaries
      - ${PROJECT_ROOT}/cannon/bin:/cannon        # GameType 0/1
      - ${PROJECT_ROOT}/asterisc/bin:/asterisc    # GameType 2, 3
      - ${PROJECT_ROOT}/bin:/kona                 # GameType 3 (kona-client)
      - ${PROJECT_ROOT}/op-program/bin:/op-program
    environment:
      OP_CHALLENGER_L1_ETH_RPC: http://l1:8545
      OP_CHALLENGER_ROLLUP_RPC: http://challenger-op-node:8545  # ⭐ Own op-node
      OP_CHALLENGER_L2_ETH_RPC: http://challenger-l2:8545        # ⭐ Own L2 geth
      OP_CHALLENGER_TRACE_TYPE: ${CHALLENGER_TRACE_TYPE}  # Supports all GameTypes
      # ... GameType-specific environment variables ...

volumes:
  l1_data:
  sequencer_l2_data:      # ← Sequencer L2 DB
  challenger_l2_data:     # ← Challenger L2 DB (independent!) ⭐⭐
  challenger_data:        # ← Challenger working directory
```

**Key Points**:
- ✅ Challenger uses **separate L2 geth** (`challenger-l2`)
- ✅ Challenger uses **separate op-node** (`challenger-op-node`, Follower mode)
- ✅ Challenger uses **separate DB** (`challenger_l2_data`)
- ✅ Mounts all GameType VM binaries to **support all game types**
- ✅ **0% sharing** with Sequencer → true independent verification

> 📚 **Actual Docker Compose Configuration**: [docker-compose-full.yml](../scripts/docker-compose-full.yml)

### Port Mapping

| Service | Port | Description |
|--------|------|------|
| **L1** | 8545 | L1 Ethereum RPC |
| **Sequencer L2** | 9545 | User transaction submission |
| **Sequencer op-node** | 7545 | Rollup info query |
| **Challenger L2** | 9546 | Independent verification RPC ⭐ |
| **Challenger op-node** | 7546 | Independent Rollup info ⭐ |
| **op-batcher** | 6545 | Batcher management |
| **op-proposer** | 6546 | Proposer management |
| **op-challenger** | 6547 | Challenger management |

---

## Security and Independence

### 1. Importance of Challenger Independence

#### Why Must It Be Independent?

```
Scenario: Sequencer Acts Maliciously

┌─────────────────────────────────────────────────────────┐
│  If Challenger shares node/DB with Sequencer...         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Malicious Sequencer:                                   │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. Create incorrect L2 blocks                 │    │
│  │  2. Store incorrect state in shared DB         │    │
│  │  3. Shared op-node provides incorrect info     │    │
│  └────────────────────────────────────────────────┘    │
│            ▼                                             │
│  Challenger (if shared):                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. Read state from shared DB                  │    │
│  │     → Already corrupted data!                  │    │
│  │  2. Get info from shared op-node               │    │
│  │     → Already corrupted info!                  │    │
│  │  3. Judges incorrect Output as "normal" ❌     │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Result: Challenger is deceived! ❌❌❌                 │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  With independent Challenger Stack...                   │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Malicious Sequencer:                                   │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. Create incorrect L2 blocks                 │    │
│  │  2. Submit incorrect Output to L1              │    │
│  └────────────────────────────────────────────────┘    │
│            ▼                                             │
│  Independent Challenger:                                │
│  ┌────────────────────────────────────────────────┐    │
│  │  1. Read batch data directly from L1 ⭐        │    │
│  │     → Bypass Sequencer, trust only L1         │    │
│  │  2. Reconstruct L2 in independent DB ⭐        │    │
│  │     → Recalculate from scratch with clean     │    │
│  │       state                                     │    │
│  │  3. Compare local Output vs L1 submitted       │    │
│  │     Output                                      │    │
│  │     → Mismatch detected! ✅                     │    │
│  │  4. Start Dispute Game                         │    │
│  │     → Remove incorrect Output ✅                │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Result: Challenger verifies correctly! ✅✅✅          │
└─────────────────────────────────────────────────────────┘
```

### 2. Independence Verification Checklist

#### Must Verify on Deployment:

```bash
# ✅ Verify independent DB
docker volume ls | grep l2_data
# Output:
# scripts_sequencer_l2_data    ← Sequencer DB
# scripts_challenger_l2_data   ← Challenger DB (different volume!)

# ✅ Verify independent op-node
docker exec challenger-op-node ps aux | grep op-node
# Check env vars:
# OP_NODE_SEQUENCER_ENABLED=false  ← Follower mode

# ✅ Verify connections
docker logs challenger-op-node | grep "L1 source"
# Output:
# L1 source: http://l1:8545  ← Reading directly from L1

# ✅ Verify DB isolation
docker exec sequencer-l2 ls /db      # Sequencer DB
docker exec challenger-l2 ls /db     # Challenger DB (different directory!)
```

### 3. Trust Model

```
┌─────────────────────────────────────────────────────────┐
│                    Trust Model                           │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Trust Required:                                        │
│  ┌────────────────────────────────────────────────┐    │
│  │  ✅ L1 Ethereum (public blockchain)            │    │
│  │  ✅ Consensus layer (verified)                 │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Trust Not Required (verifiable):                      │
│  ┌────────────────────────────────────────────────┐    │
│  │  ⭐ Sequencer (Challenger verifies)            │    │
│  │  ⭐ Proposer (Challenger verifies)             │    │
│  │  ⭐ Batcher (Recorded on L1, anyone can read)  │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Fault Proof Core Principle:                            │
│  ┌────────────────────────────────────────────────┐    │
│  │  "Don't trust, verify"                         │    │
│  │                                                 │    │
│  │  → Challenger independently verifies           │    │
│  │    everything                                   │    │
│  │  → Incorrect behavior is removed with proof    │    │
│  └────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## Summary

### Core Architecture Principles

1. **Independence**: Challenger is completely separated from Sequencer
2. **Verifiability**: All data is directly verified from L1
3. **Modularity**: Different VMs selectable by GameType
4. **Reproducibility**: Deterministic builds via Docker

### GameType-Specific Features

- **GameType 0/1 (Cannon)**: MIPS VM, stable, Go-based
- **GameType 2 (Asterisc)**: RISC-V VM, modern, Go-based
- **GameType 3 (AsteriscKona)**: RISC-V VM, lightweight, Rust-based 🆕

### Deployment Precautions

1. ✅ Verify Challenger independent DB
2. ✅ Verify Follower mode
3. ✅ Verify L1 direct connection
4. ✅ Prepare VM binaries by GameType
5. ✅ Generate Prestate correctly

---

## References

- [Optimism Specs](https://specs.optimism.io)
- [OP Stack GitHub](https://github.com/ethereum-optimism/optimism)
- [Fault Proof Documentation](https://specs.optimism.io/fault-proof/index.html)
- [Deployment Guide](../scripts/README-KR.md)
- [GameType 3 Integration Plan](./gametype3-integration-plan-ko.md)

---
