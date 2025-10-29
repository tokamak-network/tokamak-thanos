# L2 System Deployment Scripts (with Independent Challenger Stack)

This directory contains scripts for complete deployment of the Optimism L2 Rollup system.

## ⭐ Important: Independent Challenger Stack

These deployment scripts implement **fully independent Challenger stack as recommended in the documentation**:

```
Sequencer Stack:                Challenger Stack (Independent!):
├─ sequencer-l2 (L2 geth)       ├─ challenger-l2 (separate L2 geth!) ⭐
├─ sequencer-op-node            ├─ challenger-op-node (Follower) ⭐
├─ op-batcher                   └─ op-challenger
└─ op-proposer
```

**Why must it be independent?**
- Sharing nodes with Sequencer means Challenger can be fooled when Sequencer is malicious ❌
- True Fault Proof requires independent L2 reconstruction from L1 with separate DB ✅

## 📚 Reference Documentation

Please read these documents before deployment:

- **[Challenger System Architecture](../docs/challenger-system-architecture-ko.md)** ⭐ - Full system architecture and GameType details
- **[L2 System Deployment Guide](../docs/l2-system-deployment-ko.md)** - Complete deployment architecture
- **[L2 System Architecture](../docs/l2-system-architecture-ko.md)** - L2 system fundamentals


## 🚀 Quick Start

```bash
# 1. Clone repository
git clone https://github.com/tokamak-network/tokamak-thanos.git
cd tokamak-thanos
git checkout feature/challenger-gametype3

# 2. Download VM binaries (2-3 minutes)
./op-challenger/scripts/pull-vm-images.sh --tag latest
# ✅ Downloads: cannon, asterisc, op-program, kona-client

# 3. Cleanup and deploy
./op-challenger/scripts/cleanup.sh
./op-challenger/scripts/deploy-modular.sh --dg-type 3

# 4. Check status
./op-challenger/scripts/health-check.sh
./op-challenger/scripts/monitor-challenger.sh
```

**Total deployment time**: ~15-20 minutes

> 📚 **GameType Selection**: [GameType Guide](#-gametype-selection-guide)
> 📚 **VM Image Build**: [VM Image Build and Share Guide](../docs/vm-image-build-and-share-guide.md)
> 📚 **Understanding Prestate**: [Prestate Generation Guide](../docs/prestate-generation-guide-ko.md)

## 📋 Script Descriptions

### 🆕 deploy-modular.sh (Main Deployment Script)

Modularized L2 system deployment script with automatic GameType configuration and validation.

```bash
# Usage
./deploy-modular.sh [options]

# Options
--dg-type TYPE          # GameType (0, 1, 2, 3, 254, 255, default: 0)
--mode MODE             # Deployment mode: local|existing (default: local)

# Examples
./deploy-modular.sh                    # GameType 0 (default)
./deploy-modular.sh --dg-type 3        # GameType 3
FAULT_GAME_MAX_CLOCK_DURATION=150 ./deploy-modular.sh --dg-type 3  # Custom settings
```

**Key features**:
- ✅ Automatic VM and prestate validation per GameType
- ✅ Auto-set respectedGameType (sync with OptimismPortal2)
- ✅ Validate all GameType prestates (Challenger supports all types)
- ✅ Pre-deployment configuration integrity check
- ✅ Modular structure (lib/ + modules/)

> 📚 **Details**: [Script Modularization Summary](../docs/SCRIPT-MODULARIZATION-SUMMARY-ko.md)

---

### cleanup.sh

Cleanup and reset deployment state.

```bash
# Usage
./cleanup.sh [options]

# Options
--all-containers       # Clean all related containers (devnet, kurtosis, etc.)
--rebuild              # Remove Docker images and cache (force full rebuild) ⭐
--keep-env             # Keep environment file (.env)
--keep-genesis         # Keep Genesis files
--clean-vm-builds      # Remove VM binaries (default: protected) 🆕
--clean-pulled-images  # Remove registry images (default: protected) 🆕
--help                 # Show help

# Examples

# Basic cleanup (VM binaries and pulled images are protected) ⭐
./cleanup.sh

# Full rebuild (only removes local build images, protects ghcr.io images) ⭐
./cleanup.sh --rebuild

# Remove VM binaries (including those from pull-vm-images.sh)
./cleanup.sh --clean-vm-builds

# Remove registry images (ghcr.io images)
./cleanup.sh --clean-pulled-images

# Complete removal (VM + images + containers)
./cleanup.sh --clean-vm-builds --clean-pulled-images --all-containers --rebuild

# Keep .env, delete Genesis only
./cleanup.sh --keep-env

# Keep Genesis, delete .env only
./cleanup.sh --keep-genesis

# Delete containers and volumes only (keep both .env and Genesis)
./cleanup.sh --keep-env --keep-genesis
```

**⭐ Default behavior (Important)**:
- VM binaries **protected** (downloaded via pull-vm-images.sh)
- Registry images **protected** (ghcr.io images)
- Use explicit flags to remove: `--clean-vm-builds`, `--clean-pulled-images`

### setup-env.sh

Generate environment variables and wallets required for deployment.

```bash
# Usage
./setup-env.sh [options]

# Options
--mode MODE             # Deployment mode: local|sepolia|mainnet (default: local)
--output FILE           # Environment file path (default: .env)
--help                  # Show help

# Examples
./setup-env.sh --mode local
./setup-env.sh --mode sepolia --output .env.sepolia
```

**Generated wallets:**
- Admin: System administrator account
- Batcher: Submits batch data to L1
- Proposer: Submits outputs to L1
- Sequencer: Generates L2 blocks
- Challenger: Challenges invalid outputs

---

## 🎮 GameType Selection Guide

| GameType | VM | Server | Features |
|----------|-----|--------|----------|
| **0** | MIPS | op-program (Go) | Stable (default) |
| **2** | RISC-V | op-program (Go) | Latest architecture |
| **3** | RISC-V | kona-client (Rust) | Lightweight (~80%), ZK ready 🆕 |

```bash
./deploy-modular.sh                    # GameType 0 (default)
./deploy-modular.sh --dg-type 2        # GameType 2
./deploy-modular.sh --dg-type 3        # GameType 3

# Custom settings
FAULT_GAME_MAX_CLOCK_DURATION=150 PROPOSAL_INTERVAL=30s \
./deploy-modular.sh --dg-type 3
```

> 📚 **Detailed Comparison**: [Challenger System Architecture](../docs/challenger-system-architecture.md)

---

### health-check.sh

Verify the status of all deployed services.

```bash
# Usage
./health-check.sh

# Checks:
# - Docker container status
# - L1/L2 RPC responses
# - Synchronization status
# - Block heights
# - Peer counts
# - Error logs
```

### monitor-challenger.sh

Monitor Challenger's real-time activity, game participation, and sync status.

```bash
# Usage
./monitor-challenger.sh [mode]

# Modes: summary | config | sync | logs | games | errors | metrics
# Examples
./monitor-challenger.sh              # Full dashboard
./monitor-challenger.sh config       # System config and Prestate verification ⭐
```

**Key monitoring items**:
- Container status and sync
- GameType settings and Prestate verification (on-chain vs local)
- Game participation statistics and Challenger actions
- Independent stack verification (Sequencer separation)

---

## 📊 Service Endpoints

Available endpoints after deployment:

| Service | Port | URL | Description |
|---------|------|-----|-------------|
| **L1 Stack** ||||
| L1 RPC | 8545 | http://localhost:8545 | L1 Ethereum RPC |
| L1 WebSocket | 8546 | ws://localhost:8546 | L1 WebSocket |
| **Sequencer Stack** ||||
| Sequencer L2 RPC | 9545 | http://localhost:9545 | L2 User RPC |
| Sequencer op-node RPC | 7545 | http://localhost:7545 | Sequencer op-node RPC |
| op-batcher RPC | 6545 | http://localhost:6545 | op-batcher RPC |
| op-proposer RPC | 6546 | http://localhost:6546 | op-proposer RPC |
| **Challenger Stack ⭐** ||||
| Challenger L2 RPC | 9546 | http://localhost:9546 | Challenger L2 RPC (Independent!) |
| Challenger op-node RPC | 7546 | http://localhost:7546 | Challenger op-node RPC (Follower) |
| **Optional** ||||
| Indexer API | 8100 | http://localhost:8100 | Indexer API (optional) |
| Grafana | 3000 | http://localhost:3000 | Grafana (optional) |

## 🔍 Status Check Commands

```bash
# All service status
docker-compose -f op-challenger/scripts/docker-compose-full.yml ps

# Specific service logs
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f sequencer-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f challenger-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-challenger

# L1 block height
cast block-number --rpc-url http://localhost:8545

# Sequencer L2 block height
cast block-number --rpc-url http://localhost:9545

# Challenger L2 block height
cast block-number --rpc-url http://localhost:9546

# Sequencer op-node sync status
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
  http://localhost:7545

# Challenger op-node sync status
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
  http://localhost:7546
```

## 🛑 System Stop and Cleanup

```bash
# Stop services
docker-compose -f op-challenger/scripts/docker-compose-full.yml down

# Delete data as well
docker-compose -f op-challenger/scripts/docker-compose-full.yml down -v

# Delete images as well
docker-compose -f op-challenger/scripts/docker-compose-full.yml down -v --rmi all
```

## 🔧 Troubleshooting

### 1. Genesis Hash Mismatch (Most Common Issue)

**Symptoms:**
```
Error: incorrect L1 genesis block hash 0x..., expected 0x...
```

**Cause:**
- Existing Docker volumes contain blockchain data from previous deployment
- Geth ignores genesis.json when existing data is found

**Solution:**
```bash
# 1. Complete cleanup
./op-challenger/scripts/cleanup.sh

# 2. Redeploy
./op-challenger/scripts/deploy-modular.sh --dg-type 0  # or any GameType
```

### 2. No RPC Response

```bash
# Check container status
docker ps -a

# Check logs
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f [service-name]

# Restart
docker-compose -f op-challenger/scripts/docker-compose-full.yml restart [service-name]
```

### 3. Synchronization Issues

```bash
# If L2 is not following L1
# 1. Check Sequencer op-node logs
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f sequencer-op-node

# 2. Check op-batcher logs (batch submission)
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-batcher

# 3. Check L1 connection status
curl http://localhost:8545
```

### 4. Challenger Issues

```bash
# Check Challenger logs
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-challenger

# Verify Challenger uses independent stack
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f challenger-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f challenger-l2

# Check network connection status
docker network ls | grep l2-network

# Inspect network details
docker network inspect <network_name>
# Example: docker network inspect scripts_l2-network
```

### 5. Port Conflicts

```bash
# Check processes using port (macOS)
lsof -i :8545

# Kill process
kill -9 [PID]

# Or change ports in docker-compose-full.yml
```

## 📚 Additional Resources

### Official Documentation
- [Optimism Specs](https://specs.optimism.io)
- [OP Stack GitHub](https://github.com/ethereum-optimism/optimism)
- [Optimism Docs](https://docs.optimism.io)

### Project Documentation

**Deployment & Operations**:
- [L2 System Deployment Guide](../docs/l2-system-deployment-ko.md) - Complete deployment architecture
- [Script Modularization Summary](../docs/SCRIPT-MODULARIZATION-SUMMARY-ko.md) - Module structure and usage
- [Prestate Generation Guide](../docs/prestate-generation-guide-ko.md) ⭐ - Prestate generation for all GameTypes
- [VM Image Build and Share Guide](../docs/vm-image-build-and-share-guide.md) ⭐ - VM image management workflow

**GameType & Architecture**:
- [Challenger System Architecture](../docs/challenger-system-architecture.md) ⭐ - Detailed architecture per GameType
- [RISC-V GameTypes Comparison](../docs/risc-v-gametypes-comparison-ko.md) - GameType 2 vs 3 comparison
- [Game Types and VMs](../docs/game-types-and-vms-ko.md) - VM architecture explanation

**Security & Analysis**:
- [Blob Pruning Risk Analysis](../docs/blob-pruning-risk-analysis-ko.md)
- [Challenge Game Vulnerability Analysis](../docs/challenge-game-vulnerability-ko.md)
- [DA System Analysis](../docs/data-availability-analysis-ko.md)
