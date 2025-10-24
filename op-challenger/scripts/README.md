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

- **[L2 System Deployment Guide](../docs/l2-system-deployment-ko.md)** - Full deployment architecture and details
- **[L2 System Architecture](../docs/l2-system-architecture-ko.md)** - L2 system fundamentals


## 🚀 Quick Start

```bash
# 0. Clone & Checkout
git clone https://github.com/tokamak-network/tokamak-thanos.git
cd tokamak-thanos
git checkout feature/challenger-risc-gametype2

# 1. Complete cleanup
./op-challenger/scripts/cleanup.sh

# 2. Full automated deployment (first time)

# GameType 0 (Cannon - MIPS VM) - default
./op-challenger/scripts/deploy-modular.sh

# GameType 2 (Asterisc - RISC-V VM) - Auto-builds Asterisc VM! ⭐
./op-challenger/scripts/deploy-modular.sh --dg-type 2

# With custom game settings (GameType 2)
FAULT_GAME_MAX_CLOCK_DURATION=150 \
FAULT_GAME_WITHDRAWAL_DELAY=3600 \
PROPOSAL_INTERVAL=30s \
./op-challenger/scripts/deploy-modular.sh --dg-type 2

# → Takes about 5-10 minutes

# 3. Verify deployment status
./op-challenger/scripts/health-check.sh

# 4. Detailed monitoring (game statistics, sync status, etc.)
./op-challenger/scripts/monitor-challenger.sh
```


**⚠️ Important:**
- `cleanup.sh` defaults to deleting everything (`.env`, Genesis, volumes)
- To keep `.env`: `./cleanup.sh --keep-env`

## 📋 Script Descriptions

### 🆕 deploy-modular.sh (New Modular Script) ⭐

**Modularized deployment script with automatic VM build support for all GameTypes**

```bash
# Usage
./deploy-modular.sh [options]

# Options
--dg-type TYPE          # GameType (0, 1, 2, 254, 255, default: 0)
--mode MODE             # Deployment mode: local|existing (default: local)
--help                  # Show help

# Examples
./deploy-modular.sh --dg-type 0    # GameType 0 (Cannon)
./deploy-modular.sh --dg-type 2    # GameType 2 (Asterisc) - Auto-builds VM! ⭐
```

**Key Features:**
- ✅ Automatic VM build (Cannon + Asterisc)
- ✅ Modular structure (lib/ + modules/)
- ✅ Automatic TraceType validation
- ✅ Automatic prestate injection

**📚 Detailed Documentation:**
- [Script Modularization Summary](../docs/SCRIPT-MODULARIZATION-SUMMARY-ko.md) - Full guide with module structure, usage, and examples

---

### cleanup.sh

Cleanup and reset deployment state.

**⚠️ Important**: Default behavior is `--all`, which also deletes `.env` and Genesis files.

```bash
# Usage
./cleanup.sh [options]

# Options
--all-containers    # Clean all related containers (including devnet, kurtosis, etc.)
--rebuild           # Remove Docker images and build cache (forces full rebuild) ⭐
--keep-env          # Keep environment variable file (.env)
--keep-genesis      # Keep Genesis files
--help              # Show help

# Examples

# Default cleanup (.env and Genesis files also deleted)
./cleanup.sh

# Full rebuild (remove Docker images and cache) ⭐
./cleanup.sh --rebuild

# Keep .env, delete Genesis only
./cleanup.sh --keep-env

# Keep Genesis, delete .env only
./cleanup.sh --keep-genesis

# Delete containers and volumes only (keep both .env and Genesis)
./cleanup.sh --keep-env --keep-genesis

# Complete cleanup (including all related containers)
./cleanup.sh --all-containers

# Complete cleanup + rebuild
./cleanup.sh --all-containers --rebuild
```

**Items cleaned:**

**Default behavior (--all):**
- scripts-* containers (docker-compose-full.yml)
- Related volumes and networks
- .env file
- Genesis files (genesis-*.json, rollup.json)
- All allocation files in .devnet

**Additional option (--all-containers):**
- Above items +
- ops-bedrock-* containers (devnet)
- kurtosis-* containers
- All op-* related containers
- Dangling volumes

### setup-env.sh

Generate environment variables and wallets required for deployment.

```bash
# Usage
./setup-env.sh [options]

# Options
--mode MODE             # Deployment mode: local|sepolia|mainnet (default: local)
--output FILE           # Environment variable file path (default: .env)
--help                  # Show help

# Examples
./setup-env.sh --mode local
./setup-env.sh --mode sepolia --output .env.sepolia
```

**Generated wallets:**
- Admin: System administrator account
- Batcher: Account that submits batch data to L1
- Proposer: Account that submits outputs to L1
- Sequencer: Account that generates L2 blocks
- Challenger: Account that challenges invalid outputs

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

### monitor-challenger.sh ⭐ (New)

Monitor Challenger's real-time activity, game participation, and sync status.

```bash
# Usage
./monitor-challenger.sh [mode]

# Modes
summary   # Full dashboard (default)
config    # System configuration (GameType settings, Prestate verification, etc.) ⭐
sync      # Blockchain sync status only
logs      # Real-time logs (Ctrl+C to exit)
games     # Game participation details
errors    # Error analysis (Prestate mismatch detection included)
metrics   # Prometheus metrics

# Examples
./monitor-challenger.sh              # Full dashboard
./monitor-challenger.sh config       # System config and Prestate verification ⭐
./monitor-challenger.sh sync         # Sync status only
./monitor-challenger.sh logs         # Real-time logs
./monitor-challenger.sh games        # Game activity only
./monitor-challenger.sh errors       # Errors and Prestate verification
```

**Monitoring items:**
- ✅ Challenger container status
- ✅ **System configuration (GameType, game time, proposer interval, etc.)** ⭐
- ✅ **Prestate verification (on-chain vs local Docker build comparison)** ⭐
- ✅ L1/Sequencer/Challenger block height comparison and sync differences
- ✅ L1 batch submission status (op-batcher)
- ✅ Challenger independence verification (using independent op-node, L2 geth)
- ✅ Game participation statistics (detected, in progress, resolved)
- ✅ Challenger action logs (Attack/Defend)
- ✅ Error and warning analysis (Prestate mismatch detection)
- ✅ Prometheus metrics (when enabled)

**Key features:**
- Color-coded logs for readability (errors=red, warnings=yellow, actions=green)
- Sequencer vs Challenger sync difference monitoring (auto-detect normal/delayed)
- L1 batch submission and Challenger batch reading status
- Independent stack usage verification (separation from Sequencer)
- **On-chain configuration verification (query deployed game settings)** ⭐
- **Prestate verification (on-chain absolute prestate vs local build prestate auto-comparison)** ⭐
- **Prestate mismatch error auto-detection and solution guide** ⭐

## ⚙️ Game Settings

### Default Game Settings (devnetL1-template.json)

```json
{
  "faultGameMaxClockDuration": 1200,        // 20 minutes (in seconds)
  "faultGameClockExtension": 0,             // Clock extension time
  "faultGameMaxDepth": 50,                  // Maximum game depth
  "faultGameSplitDepth": 14,                // Split depth
  "faultGameWithdrawalDelay": 604800,       // Withdrawal delay (7 days)
  "disputeGameFinalityDelaySeconds": 6,     // Finality delay (6 seconds)
  "faultGameAbsolutePrestate": "0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98",
  "faultGameGenesisOutputRoot": "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF"
}
```

**Key settings:**

- **faultGameMaxClockDuration**: 1200s (20 min) - Maximum chess clock time per team
  - ⚠️ Minimum value: 120s (due to `preimageOracleChallengePeriod`)
- **faultGameWithdrawalDelay**: 604800s (7 days) - Withdrawal waiting time after game ends
- **faultGameSplitDepth**: 14 - Covers 16,384 L2 blocks (~9.1 hours)
- **faultGameMaxDepth**: 50 - Maximum bisection depth (68.7B instructions)

**Quick test environment via environment variables:**

```bash
FAULT_GAME_MAX_CLOCK_DURATION=150 \
FAULT_GAME_WITHDRAWAL_DELAY=3600 \
PROPOSAL_INTERVAL=30s \
./deploy-modular.sh --dg-type 0  # or --dg-type 2 for Asterisc
```

**Settings:**
- Game duration: 2.5 minutes
- Withdrawal delay: 1 hour
- Proposal interval: 30 seconds

**⚠️ Note:** These settings are for development/testing only. Use sufficient time for production.

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
./op-challenger/scripts/deploy-modular.sh --dg-type 0  # or --dg-type 2 for Asterisc
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

# Inspect network details (network name may vary by execution directory)
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
- [Blob Pruning Risk Analysis](../docs/blob-pruning-risk-analysis-ko.md)
- [Challenge Game Vulnerability Analysis](../docs/challenge-game-vulnerability-ko.md)
- [Data Availability Analysis](../docs/data-availability-analysis-ko.md)
