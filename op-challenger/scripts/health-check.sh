#!/usr/bin/env bash
set -euo pipefail

##############################################################################
# Service Health Check Script (Including Independent Challenger Stack)
#
# Checks the status of all deployed L2 system components.
#
# Reference: /op-challenger/docs/l2-system-deployment-ko.md
##############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCKER_COMPOSE_FILE="${SCRIPT_DIR}/docker-compose-full.yml"

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Log functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[!]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Health check counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

##############################################################################
# RPC health check function
##############################################################################

check_rpc() {
    local service_name="$1"
    local rpc_url="$2"
    local timeout="${3:-5}"

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    if curl -s -m "$timeout" -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "$rpc_url" > /dev/null 2>&1; then

        # Get block number
        block_hex=$(curl -s -X POST -H "Content-Type: application/json" \
            --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
            "$rpc_url" | jq -r '.result')

        if [ "$block_hex" != "null" ] && [ -n "$block_hex" ]; then
            block_num=$((16#${block_hex#0x}))
            log_success "$service_name - RPC responding (block: $block_num)"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        else
            log_warn "$service_name - RPC responding but no block number"
            FAILED_CHECKS=$((FAILED_CHECKS + 1))
            return 1
        fi
    else
        log_error "$service_name - No RPC response"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

# op-node specific health check (uses optimism_syncStatus instead of eth_blockNumber)
check_opnode_rpc() {
    local service_name="$1"
    local rpc_url="$2"
    local timeout="${3:-5}"

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    # Check if RPC is responding using optimism_syncStatus
    sync_status=$(curl -s -m "$timeout" -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"optimism_syncStatus","params":[],"id":1}' \
        "$rpc_url" 2>/dev/null)

    if [ -n "$sync_status" ]; then
        # Get unsafe L2 block number
        unsafe_l2_num=$(echo "$sync_status" | jq -r '.result.unsafe_l2.number // 0' 2>/dev/null)

        if [ "$unsafe_l2_num" != "null" ] && [ -n "$unsafe_l2_num" ] && [ "$unsafe_l2_num" != "0" ]; then
            log_success "$service_name - RPC responding (unsafe L2: $unsafe_l2_num)"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        else
            log_warn "$service_name - RPC responding but L2 not synced yet"
            FAILED_CHECKS=$((FAILED_CHECKS + 1))
            return 1
        fi
    else
        log_error "$service_name - No RPC response"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

##############################################################################
# Docker container status check
##############################################################################

check_docker_container() {
    local container_pattern="$1"
    local service_name="$2"

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    # Find container name pattern (including stopped containers)
    local container_name=$(docker ps -a --format '{{.Names}}' | grep "scripts-.*$container_pattern" | head -1)

    if [ -n "$container_name" ]; then
        local status=$(docker inspect --format='{{.State.Status}}' "$container_name" 2>/dev/null || echo "unknown")

        if [ "$status" = "running" ]; then
            log_success "Container: $service_name - Running"
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
            return 0
        elif [ "$status" = "exited" ]; then
            local exit_code=$(docker inspect --format='{{.State.ExitCode}}' "$container_name" 2>/dev/null || echo "?")
            log_error "Container: $service_name - Exited (code: $exit_code)"
            FAILED_CHECKS=$((FAILED_CHECKS + 1))

            # Show last 3 lines of logs for failed containers
            echo "  └─ Last error:"
            docker logs "$container_name" 2>&1 | tail -3 | sed 's/^/     /'
            return 1
        else
            log_error "Container: $service_name - Status: $status"
            FAILED_CHECKS=$((FAILED_CHECKS + 1))
            return 1
        fi
    else
        log_error "Container: $service_name - Not found (pattern: scripts-.*$container_pattern)"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

##############################################################################
# Synchronization status check
##############################################################################

check_sync_status() {
    local service_name="$1"
    local rpc_url="$2"

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    # Call eth_syncing
    sync_status=$(curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
        "$rpc_url" | jq -r '.result')

    if [ "$sync_status" = "false" ]; then
        log_success "$service_name - Synchronized"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    elif [ "$sync_status" = "null" ]; then
        log_warn "$service_name - Cannot check sync status"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    else
        # Syncing
        current=$(echo "$sync_status" | jq -r '.currentBlock // 0' 2>/dev/null || echo 0)
        highest=$(echo "$sync_status" | jq -r '.highestBlock // 0' 2>/dev/null || echo 0)

        if [ "$current" != "0" ] && [ "$highest" != "0" ]; then
            log_warn "$service_name - Syncing ($current / $highest)"
        else
            log_warn "$service_name - Syncing"
        fi
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

##############################################################################
# Main health check
##############################################################################

echo ""
log_info "=========================================="
log_info "L2 System Health Check Started"
log_info "=========================================="
echo ""

# 1. Check Docker container status
log_info "━━━ Docker Container Status ━━━"
check_docker_container "l1" "L1 Ethereum" || true
check_docker_container "sequencer-l2" "Sequencer L2 geth" || true
check_docker_container "sequencer-op-node" "Sequencer op-node" || true
check_docker_container "op-batcher" "op-batcher" || true
check_docker_container "op-proposer" "op-proposer" || true
echo ""
log_info "━━━ Challenger Stack Container Status ⭐ ━━━"
check_docker_container "challenger-l2" "Challenger L2 geth (Independent)" || true
check_docker_container "challenger-op-node" "Challenger op-node (Follower)" || true
check_docker_container "op-challenger" "op-challenger" || true
echo ""

# 2. L1 RPC check
log_info "━━━ L1 Ethereum RPC ━━━"
check_rpc "L1 geth" "http://localhost:8545" 5 || true
check_sync_status "L1 geth" "http://localhost:8545" || true
echo ""

# 3. Sequencer stack check
log_info "━━━ Sequencer Stack ━━━"
check_rpc "Sequencer L2 op-geth" "http://localhost:9545" 5 || true
check_sync_status "Sequencer L2" "http://localhost:9545" || true
check_opnode_rpc "Sequencer op-node" "http://localhost:7545" 5 || true
echo ""

# 4. Challenger stack check (Independent!) ⭐
log_info "━━━ Challenger Stack (Independent!) ⭐ ━━━"
check_rpc "Challenger L2 op-geth (Separate DB)" "http://localhost:9546" 5 || true
check_sync_status "Challenger L2" "http://localhost:9546" || true
check_opnode_rpc "Challenger op-node (Follower)" "http://localhost:7546" 5 || true

# Check Challenger logs
if docker ps --format '{{.Names}}' | grep -q "op-challenger"; then
    challenger_container=$(docker ps --format '{{.Names}}' | grep "op-challenger" | head -1)
    log_info "Checking op-challenger logs..."

    # Check for errors in recent logs
    error_count=$(docker logs --tail=50 "$challenger_container" 2>&1 | grep -i "error" | wc -l | tr -d ' \n' || echo "0")
    error_count=${error_count:-0}  # 빈 값 방지

    if [ "$error_count" -lt 5 ]; then
        log_success "op-challenger - Operating normally (error logs: $error_count)"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
    else
        log_warn "op-challenger - Multiple errors found ($error_count)"
        log_info "Detailed logs: docker logs $challenger_container"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
    fi
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
fi
echo ""

##############################################################################
# GameType Configuration Check
##############################################################################

log_info "━━━ GameType Configuration ━━━"

# Load addresses
ADDRESSES_FILE="${PROJECT_ROOT:-.}/.devnet/addresses.json"
if [ -f "$ADDRESSES_FILE" ]; then
    DGF_ADDRESS=$(jq -r '.DisputeGameFactoryProxy // .DisputeGameFactory' "$ADDRESSES_FILE" 2>/dev/null || echo "")

    if [ -n "$DGF_ADDRESS" ] && [ "$DGF_ADDRESS" != "null" ]; then
        log_info "DisputeGameFactory: $DGF_ADDRESS"

        # Check GameType 0 (CANNON)
        GT0_IMPL=$(cast call --rpc-url "http://localhost:8545" "$DGF_ADDRESS" "gameImpls(uint32)(address)" 0 2>/dev/null || echo "")
        if [ -n "$GT0_IMPL" ] && [ "$GT0_IMPL" != "0x0000000000000000000000000000000000000000" ]; then
            log_success "GameType 0 (CANNON) - Deployed: $GT0_IMPL"
            TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
        else
            log_warn "GameType 0 (CANNON) - Not deployed"
            TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
            FAILED_CHECKS=$((FAILED_CHECKS + 1))
        fi

        # Check GameType 1 (PERMISSIONED_CANNON)
        GT1_IMPL=$(cast call --rpc-url "http://localhost:8545" "$DGF_ADDRESS" "gameImpls(uint32)(address)" 1 2>/dev/null || echo "")
        if [ -n "$GT1_IMPL" ] && [ "$GT1_IMPL" != "0x0000000000000000000000000000000000000000" ]; then
            log_success "GameType 1 (PERMISSIONED_CANNON) - Deployed: $GT1_IMPL"
            TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
            PASSED_CHECKS=$((PASSED_CHECKS + 1))
        else
            log_info "GameType 1 (PERMISSIONED_CANNON) - Not deployed (optional)"
        fi

        # Check GameType 2 (ASTERISC) ⭐
        GT2_IMPL=$(cast call --rpc-url "http://localhost:8545" "$DGF_ADDRESS" "gameImpls(uint32)(address)" 2 2>/dev/null || echo "")
        if [ -n "$GT2_IMPL" ] && [ "$GT2_IMPL" != "0x0000000000000000000000000000000000000000" ]; then
            log_success "GameType 2 (ASTERISC/RISCV) - Deployed: $GT2_IMPL ⭐"
            TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
            PASSED_CHECKS=$((PASSED_CHECKS + 1))

            # Get RISCV VM address
            RISCV_VM=$(cast call --rpc-url "http://localhost:8545" "$GT2_IMPL" "vm()(address)" 2>/dev/null || echo "")
            if [ -n "$RISCV_VM" ] && [ "$RISCV_VM" != "0x0000000000000000000000000000000000000000" ]; then
                log_info "  └─ RISCV VM: $RISCV_VM"
            fi
        else
            log_info "GameType 2 (ASTERISC/RISCV) - Not deployed (optional) ⭐"
        fi

        # Check GameType 3 (ASTERISC_KONA) 🆕
        GT3_IMPL=$(cast call --rpc-url "http://localhost:8545" "$DGF_ADDRESS" "gameImpls(uint32)(address)" 3 2>/dev/null || echo "")
        if [ -n "$GT3_IMPL" ] && [ "$GT3_IMPL" != "0x0000000000000000000000000000000000000000" ]; then
            log_success "GameType 3 (ASTERISC_KONA/RISCV + Rust) - Deployed: $GT3_IMPL 🆕"
            TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
            PASSED_CHECKS=$((PASSED_CHECKS + 1))

            # Get RISCV VM address (should be same as GameType 2)
            RISCV_VM_GT3=$(cast call --rpc-url "http://localhost:8545" "$GT3_IMPL" "vm()(address)" 2>/dev/null || echo "")
            if [ -n "$RISCV_VM_GT3" ] && [ "$RISCV_VM_GT3" != "0x0000000000000000000000000000000000000000" ]; then
                log_info "  └─ RISCV VM: $RISCV_VM_GT3"

                # Check if same as GameType 2
                if [ -n "$RISCV_VM" ] && [ "$RISCV_VM" = "$RISCV_VM_GT3" ]; then
                    log_info "  └─ ✅ Shares same RISCV.sol with GameType 2"
                fi
            fi
        else
            log_info "GameType 3 (ASTERISC_KONA) - Not deployed (optional) 🆕"
        fi

        # Check GameType 254 (FAST)
        GT254_IMPL=$(cast call --rpc-url "http://localhost:8545" "$DGF_ADDRESS" "gameImpls(uint32)(address)" 254 2>/dev/null || echo "")
        if [ -n "$GT254_IMPL" ] && [ "$GT254_IMPL" != "0x0000000000000000000000000000000000000000" ]; then
            log_info "GameType 254 (FAST) - Deployed: $GT254_IMPL"
        fi

        # Check GameType 255 (ALPHABET)
        GT255_IMPL=$(cast call --rpc-url "http://localhost:8545" "$DGF_ADDRESS" "gameImpls(uint32)(address)" 255 2>/dev/null || echo "")
        if [ -n "$GT255_IMPL" ] && [ "$GT255_IMPL" != "0x0000000000000000000000000000000000000000" ]; then
            log_info "GameType 255 (ALPHABET) - Deployed: $GT255_IMPL"
        fi
    else
        log_warn "DisputeGameFactory address not found in addresses.json"
    fi
else
    log_warn "addresses.json not found: $ADDRESSES_FILE"
fi

echo ""

##############################################################################
# Collect additional information
##############################################################################

log_info "━━━ Block Height Comparison ━━━"

# L1 and Sequencer L2 block heights
l1_block=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "http://localhost:8545" 2>/dev/null | jq -r '.result' || echo "0x0")

sequencer_l2_block=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "http://localhost:9545" 2>/dev/null | jq -r '.result' || echo "0x0")

# Challenger L2 block height
challenger_l2_block=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "http://localhost:9546" 2>/dev/null | jq -r '.result' || echo "0x0")

if [ "$l1_block" != "0x0" ]; then
    l1_num=$((16#${l1_block#0x}))
    echo "L1 Block Height:            $l1_num"
fi

if [ "$sequencer_l2_block" != "0x0" ]; then
    seq_l2_num=$((16#${sequencer_l2_block#0x}))
    echo "Sequencer L2 Block Height:  $seq_l2_num"
fi

if [ "$challenger_l2_block" != "0x0" ]; then
    chal_l2_num=$((16#${challenger_l2_block#0x}))
    echo "Challenger L2 Block Height: $chal_l2_num ⭐"

    # Compare Sequencer and Challenger L2
    if [ "$sequencer_l2_block" != "0x0" ] && [ "$challenger_l2_block" != "0x0" ]; then
        if [ "$seq_l2_num" = "$chal_l2_num" ]; then
            log_success "Sequencer and Challenger L2 are synchronized!"
        else
            diff=$((seq_l2_num - chal_l2_num))
            if [ $diff -gt 0 ]; then
                log_info "Challenger L2 is ${diff} blocks behind Sequencer (normal)"
            else
                log_warn "Challenger L2 is ahead of Sequencer (abnormal)"
            fi
        fi
    fi
fi

echo ""

# Check peer count
log_info "━━━ P2P Peer Count ━━━"

seq_peer_count=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' \
    "http://localhost:7545" 2>/dev/null | jq -r '.result' || echo "0x0")

chal_peer_count=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' \
    "http://localhost:7546" 2>/dev/null | jq -r '.result' || echo "0x0")

if [ "$seq_peer_count" != "0x0" ] && [ "$seq_peer_count" != "null" ] && [ -n "$seq_peer_count" ]; then
    peers=$((16#${seq_peer_count#0x}))
    echo "Sequencer op-node peers: $peers"
else
    echo "Sequencer op-node peers: unavailable"
fi

if [ "$chal_peer_count" != "0x0" ] && [ "$chal_peer_count" != "null" ] && [ -n "$chal_peer_count" ]; then
    peers=$((16#${chal_peer_count#0x}))
    echo "Challenger op-node peers: $peers"
else
    echo "Challenger op-node peers: unavailable"
fi

echo ""

##############################################################################
# Independence verification
##############################################################################

log_info "━━━ Challenger Independence Verification ⭐ ━━━"

# Check if Challenger uses its own op-node and op-geth
challenger_container=$(docker ps -a --format '{{.Names}}' | grep "scripts-.*op-challenger" | head -1 || echo "")

if [ -n "$challenger_container" ]; then
    # Check op-challenger environment variables
    challenger_rollup_rpc=$(docker inspect "$challenger_container" | jq -r '.[0].Config.Env[]' | grep "OP_CHALLENGER_ROLLUP_RPC" || echo "")
    challenger_l2_rpc=$(docker inspect "$challenger_container" | jq -r '.[0].Config.Env[]' | grep "OP_CHALLENGER_L2_ETH_RPC" || echo "")

    if echo "$challenger_rollup_rpc" | grep -q "challenger-op-node"; then
        log_success "op-challenger is using independent challenger-op-node ✅"
    else
        log_error "op-challenger is using sequencer op-node! ❌"
        echo "  $challenger_rollup_rpc"
    fi

    if echo "$challenger_l2_rpc" | grep -q "challenger-l2"; then
        log_success "op-challenger is using independent challenger-l2 ✅"
    else
        log_error "op-challenger is using sequencer l2! ❌"
        echo "  $challenger_l2_rpc"
    fi
else
    log_warn "Cannot find op-challenger container"
fi

echo ""

##############################################################################
# Result summary
##############################################################################

echo ""
log_info "=========================================="
log_info "Health Check Completed"
log_info "=========================================="
echo ""

echo "Total Checks: $TOTAL_CHECKS"
echo -e "${GREEN}Passed: $PASSED_CHECKS${NC}"
echo -e "${RED}Failed: $FAILED_CHECKS${NC}"
echo ""

if [ $FAILED_CHECKS -eq 0 ]; then
    log_success "All services are operating normally!"
    echo ""
    log_info "⭐ Challenger stack is operating independently from Sequencer!"
    log_info "   - challenger-l2: Independently reconstructs L2 from L1 with separate DB"
    log_info "   - challenger-op-node: Verifies L1 data in Follower mode"
    log_info "   - op-challenger: Performs challenges with independently verified data"
    exit 0
elif [ $PASSED_CHECKS -gt $FAILED_CHECKS ]; then
    log_warn "Some services have issues. Check the logs."
    log_info "Check logs: docker-compose -f ${DOCKER_COMPOSE_FILE} logs -f [service-name]"
    exit 1
else
    log_error "Multiple services have encountered problems!"
    log_info "All logs: docker-compose -f ${DOCKER_COMPOSE_FILE} logs -f"
    log_info "Specific service: docker-compose -f ${DOCKER_COMPOSE_FILE} logs -f sequencer-op-node"
    exit 2
fi
