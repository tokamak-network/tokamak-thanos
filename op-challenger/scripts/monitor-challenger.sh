#!/usr/bin/env bash
set -euo pipefail

##############################################################################
# Challenger Monitoring Script
#
# Monitors op-challenger activity, game participation, and performance
#
# Reference: /op-challenger/docs/op-challenger-architecture-ko.md
##############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DOCKER_COMPOSE_FILE="${SCRIPT_DIR}/docker-compose-full.yml"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

log_title() {
    echo -e "${CYAN}━━━ $1 ━━━${NC}"
}

##############################################################################
# Helper functions
##############################################################################

get_challenger_container() {
    docker ps --format '{{.Names}}' | grep "op-challenger" | head -1 || echo ""
}

check_challenger_running() {
    local container=$(get_challenger_container)
    if [ -z "$container" ]; then
        log_error "op-challenger container is not running"
        exit 1
    fi
    echo "$container"
}

##############################################################################
# Main monitoring functions
##############################################################################

# 1. Challenger Container Status
show_container_status() {
    log_title "Challenger Container Status"

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        local status=$(docker inspect --format='{{.State.Status}}' "$container" 2>/dev/null || echo "unknown")
        local uptime=$(docker inspect --format='{{.State.StartedAt}}' "$container" 2>/dev/null || echo "unknown")

        if [ "$status" = "running" ]; then
            log_success "Container: $container"
            echo "  Status:  $status"
            echo "  Started: $uptime"
        else
            log_error "Container: $container"
            echo "  Status: $status"
        fi
    else
        log_error "op-challenger container not found"
    fi
    echo ""
}

# 2. Recent Challenger Logs
show_recent_logs() {
    log_title "Recent Challenger Activity (Last 20 lines)"

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        docker logs --tail=20 "$container" 2>&1 | while IFS= read -r line; do
            # Color-code different log levels
            if echo "$line" | grep -q "lvl=error\|lvl=crit"; then
                echo -e "${RED}$line${NC}"
            elif echo "$line" | grep -q "lvl=warn"; then
                echo -e "${YELLOW}$line${NC}"
            elif echo "$line" | grep -q "Game info\|Performing action\|Resolving"; then
                echo -e "${GREEN}$line${NC}"
            else
                echo "$line"
            fi
        done
    fi
    echo ""
}

# 3. Game Participation Summary
show_game_summary() {
    log_title "Game Participation Summary"

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        echo "Analyzing logs for game activity..."

        # Get all logs once for efficiency
        local all_logs=$(docker logs "$container" 2>&1)

        # Count resolved games first (unique)
        local games_resolved=$(echo "$all_logs" | grep "Game resolved" | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' | sort -u | wc -l | tr -d ' ')
        local defender_won=$(echo "$all_logs" | grep "Game resolved" | grep 'status="Defender Won"' | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' | sort -u | wc -l | tr -d ' ')
        local challenger_won=$(echo "$all_logs" | grep "Game resolved" | grep 'status="Challenger Won"' | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' | sort -u | wc -l | tr -d ' ')

        # Get list of resolved game addresses
        local resolved_games=$(echo "$all_logs" | grep "Game resolved" | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' | sort -u)

        # Count all unique games
        local all_games=$(echo "$all_logs" | grep -E "Game info|Game resolved" | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' | sort -u)
        local games_detected=$(echo "$all_games" | wc -l | tr -d ' ')

        # Count in-progress games (all games minus resolved games)
        local games_in_progress=0
        if [ -n "$all_games" ]; then
            if [ -n "$resolved_games" ]; then
                games_in_progress=$(echo "$all_games" | grep -v -F -f <(echo "$resolved_games") | wc -l | tr -d ' ')
            else
                # No resolved games, so all games are in progress
                games_in_progress=$games_detected
            fi
        fi

        # Count moves made (reuse all_logs)
        local moves_made=$(echo "$all_logs" | grep "Performing action" | wc -l | tr -d ' ')

        # Count errors (last 100 lines of all_logs)
        local recent_logs=$(echo "$all_logs" | tail -100)
        local errors=$(echo "$recent_logs" | grep "lvl=error" | wc -l | tr -d ' ')

        # Count invalid prestate warnings
        local invalid_prestate=$(echo "$recent_logs" | grep "Invalid prestate" | wc -l | tr -d ' ')

        # Verify totals
        local calculated_total=$((games_in_progress + games_resolved))
        local match_status=""
        if [ "$calculated_total" -eq "$games_detected" ]; then
            match_status=" ✓"
        else
            match_status=" ⚠️ (calc: $calculated_total)"
        fi

        echo ""
        echo "  📊 Total Games:           $games_detected$match_status"
        echo "     ├─ 🟢 Active (In Progress):  $games_in_progress"
        echo "     └─ ✅ Resolved:              $games_resolved"
        if [ "$games_resolved" -gt 0 ]; then
            # Calculate percentages
            local total_resolved=$((defender_won + challenger_won))
            if [ "$total_resolved" -gt 0 ]; then
                local defender_pct=$((defender_won * 100 / total_resolved))
                local challenger_pct=$((challenger_won * 100 / total_resolved))
                echo "        ├─ 🛡️  Defender Won:     $defender_won ($defender_pct%)"
                echo "        └─ ⚔️  Challenger Won:   $challenger_won ($challenger_pct%)"
            fi
        fi
        echo ""
        echo "  🎯 Moves Made:            $moves_made"
        echo "  ⚠️  Recent Errors:         $errors"

        if [ "$invalid_prestate" -gt 0 ]; then
            echo ""
            log_warn "Invalid Prestate Warnings: $invalid_prestate"
            echo "     └─ This may indicate prestate mismatch with DisputeGameFactory"
        fi
    fi
    echo ""
}

# 4. Active Games List
show_active_games() {
    log_title "Active Games (Last 10)"

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        docker logs --tail=200 "$container" 2>&1 | grep "Game info" | tail -10 | while IFS= read -r line; do
            # Extract game address and status
            local game=$(echo "$line" | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' || echo "")
            local claims=$(echo "$line" | sed -n 's/.*claims=\([0-9]*\).*/\1/p' || echo "")
            local status=$(echo "$line" | sed -n 's/.*status="\([^"]*\)".*/\1/p' || echo "")

            if [ -n "$game" ]; then
                if [ "$status" = "In Progress" ]; then
                    echo -e "  ${CYAN}$game${NC} - Claims: $claims - ${YELLOW}$status${NC}"
                elif [ "$status" = "Challenger Wins" ]; then
                    echo -e "  ${CYAN}$game${NC} - Claims: $claims - ${GREEN}$status${NC}"
                elif [ "$status" = "Defender Wins" ]; then
                    echo -e "  ${CYAN}$game${NC} - Claims: $claims - ${RED}$status${NC}"
                else
                    echo -e "  ${CYAN}$game${NC} - Claims: $claims - $status"
                fi
            fi
        done
    fi
    echo ""
}

# 5. Recently Resolved Games
show_resolved_games() {
    log_title "Recently Resolved Games (Last 10)"

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        docker logs --tail=500 "$container" 2>&1 | grep "Game resolved" | tail -10 | while IFS= read -r line; do
            # Extract game address, status, and timestamp
            local timestamp=$(echo "$line" | sed -n 's/^t=\([^ ]*\).*/\1/p' || echo "")
            local game=$(echo "$line" | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' || echo "")
            local status=$(echo "$line" | sed -n 's/.*status="\([^"]*\)".*/\1/p' || echo "")

            # Format time (remove timezone for brevity)
            local time_short=$(echo "$timestamp" | sed 's/\+.*//' | sed 's/T/ /')

            if [ -n "$game" ]; then
                if [ "$status" = "Challenger Won" ]; then
                    echo -e "  ${CYAN}$game${NC} - ${GREEN}✓ $status${NC} - $time_short"
                elif [ "$status" = "Defender Won" ]; then
                    echo -e "  ${CYAN}$game${NC} - ${BLUE}✓ $status${NC} - $time_short"
                else
                    echo -e "  ${CYAN}$game${NC} - ✓ $status - $time_short"
                fi
            fi
        done
    fi
    echo ""
}

# 6. Challenger Actions Log
show_actions() {
    log_title "Recent Challenger Actions (Last 10)"

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        docker logs --tail=200 "$container" 2>&1 | grep "Performing action" | tail -10 | while IFS= read -r line; do
            local game=$(echo "$line" | sed -n 's/.*game=\(0x[a-fA-F0-9]*\).*/\1/p' || echo "")
            local action=$(echo "$line" | sed -n 's/.*action=\([a-z]*\).*/\1/p' || echo "")
            local attack=$(echo "$line" | sed -n 's/.*is_attack=\(true\|false\).*/\1/p' || echo "")

            if [ "$attack" = "true" ]; then
                echo -e "  ${CYAN}$game${NC} - ${RED}Attack${NC} ($action)"
            else
                echo -e "  ${CYAN}$game${NC} - ${GREEN}Defend${NC} ($action)"
            fi
        done
    fi
    echo ""
}

# 6. Error Analysis
show_errors() {
    log_title "Recent Errors (Last 10)"

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        local error_lines=$(docker logs --tail=100 "$container" 2>&1 | grep -E "lvl=error|lvl=crit" | tail -10)

        if [ -n "$error_lines" ]; then
            echo "$error_lines" | while IFS= read -r line; do
                echo -e "${RED}$line${NC}"
            done
        else
            log_success "No recent errors found!"
        fi
    fi
    echo ""
}

# 7. Metrics (if available)
show_metrics() {
    log_title "Challenger Metrics"

    # Try to fetch Prometheus metrics
    local metrics_url="http://localhost:7304/metrics"

    if curl -s -m 2 "$metrics_url" > /dev/null 2>&1; then
        log_info "Fetching metrics from $metrics_url..."
        echo ""

        # Games metrics
        local games_in_progress=$(curl -s "$metrics_url" | grep "^op_challenger_games_in_progress" | awk '{print $2}' || echo "0")
        local games_won=$(curl -s "$metrics_url" | grep "^op_challenger_games_challenger_won" | awk '{print $2}' || echo "0")
        local games_lost=$(curl -s "$metrics_url" | grep "^op_challenger_games_defender_won" | awk '{print $2}' || echo "0")

        echo "  Games In Progress: $games_in_progress"
        echo "  Games Won:         $games_won"
        echo "  Games Lost:        $games_lost"
    else
        log_warn "Metrics endpoint not accessible: $metrics_url"
        echo "  Make sure Challenger is running with --metrics.enabled"
    fi
    echo ""
}

# 8. System Configuration
show_system_config() {
    log_title "System Configuration"

    # Load addresses
    local addresses_file="${PROJECT_ROOT}/.devnet/addresses.json"
    if [ ! -f "$addresses_file" ]; then
        log_warn "addresses.json not found"
        echo ""
        return
    fi

    local dgf_address=$(jq -r '.DisputeGameFactoryProxy // .DisputeGameFactory' "$addresses_file" 2>/dev/null || echo "")

    if [ -z "$dgf_address" ] || [ "$dgf_address" = "null" ]; then
        log_warn "DisputeGameFactory address not found"
        echo ""
        return
    fi

    echo "  DisputeGameFactory: $dgf_address"
    echo ""

    # Get game implementation for type 0
    local game_impl=$(cast call "$dgf_address" "gameImpls(uint32)(address)" 0 --rpc-url http://localhost:8545 2>/dev/null || echo "")

    if [ -n "$game_impl" ] && [ "$game_impl" != "0x0000000000000000000000000000000000000000" ]; then
        echo "  Game Type 0 (CANNON) Implementation:"
        echo "    Address: $game_impl"

        # Get game configuration
        local max_clock=$(cast call "$game_impl" "maxClockDuration()(uint64)" --rpc-url http://localhost:8545 2>/dev/null || echo "0")
        local max_clock_min=$((max_clock / 60))
        local max_clock_hours=$((max_clock / 3600))

        local absolute_prestate=$(cast call "$game_impl" "absolutePrestate()(bytes32)" --rpc-url http://localhost:8545 2>/dev/null || echo "")

        echo "    Max Clock Duration: ${max_clock}s (${max_clock_min} min)"
        echo "    Absolute Prestate: ${absolute_prestate:0:18}...${absolute_prestate: -6}"

        # Get DelayedWETH address and withdrawal delay
        local delayed_weth=$(jq -r '.DelayedWETHProxy // .DelayedWETH' "$addresses_file" 2>/dev/null || echo "")
        if [ -n "$delayed_weth" ] && [ "$delayed_weth" != "null" ]; then
            local withdrawal_delay=$(cast call "$delayed_weth" "delay()(uint256)" --rpc-url http://localhost:8545 2>/dev/null | sed 's/\x1b\[[0-9;]*m//g' | awk '{print $1}' || echo "0")
            if [ -n "$withdrawal_delay" ] && [ "$withdrawal_delay" != "0" ]; then
                local withdrawal_days=$((withdrawal_delay / 86400))
                local withdrawal_hours=$(( (withdrawal_delay % 86400) / 3600 ))
                echo "    Withdrawal Delay: ${withdrawal_delay}s (${withdrawal_days}d ${withdrawal_hours}h)"
            fi
        fi

        # Get finality delay from DisputeGameFactory
        local finality_delay=$(cast call "$dgf_address" "initBonds(uint32)(uint256)" 0 --rpc-url http://localhost:8545 2>/dev/null || echo "0")
        if [ "$finality_delay" = "0" ]; then
            # Try alternative method - check if there's a respectedGameType function
            finality_delay="6s (default)"
            echo "    Finality Delay: $finality_delay"
        fi
    fi

    echo ""

    # Proposer configuration
    local proposer_container=$(docker ps --format '{{.Names}}' | grep "op-proposer" | head -1 || echo "")
    if [ -n "$proposer_container" ]; then
        echo "  Proposer Configuration:"
        local game_type=$(docker inspect "$proposer_container" 2>/dev/null | jq -r '.[0].Config.Env[]' | grep "OP_PROPOSER_GAME_TYPE" | cut -d= -f2 || echo "unknown")
        local proposal_interval=$(docker inspect "$proposer_container" 2>/dev/null | jq -r '.[0].Config.Env[]' | grep "OP_PROPOSER_PROPOSAL_INTERVAL" | cut -d= -f2 || echo "unknown")

        echo "    Game Type: $game_type"
        echo "    Proposal Interval: $proposal_interval"
    fi

    echo ""

    # Challenger configuration
    local challenger_container=$(docker ps --format '{{.Names}}' | grep "op-challenger" | head -1 || echo "")
    if [ -n "$challenger_container" ]; then
        echo "  Challenger Configuration:"
        local challenger_addr=$(docker inspect "$challenger_container" 2>/dev/null | jq -r '.[0].Config.Env[]' | grep "OP_CHALLENGER_PRIVATE_KEY" | cut -d= -f2 || echo "")
        if [ -n "$challenger_addr" ] && command -v cast >/dev/null 2>&1; then
            challenger_addr=$(cast wallet address "$challenger_addr" 2>/dev/null || echo "unknown")
            echo "    Address: $challenger_addr"
        fi
    fi

    echo ""
}

# 9. Synchronization Status (L1, Sequencer, Challenger)
show_sync_status() {
    log_title "Blockchain Synchronization Status"

    # L1 Block Height
    local l1_block=$(curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "http://localhost:8545" 2>/dev/null | jq -r '.result' || echo "0x0")
    local l1_num=$((16#${l1_block#0x} || 0))

    # Sequencer L2 Block Height
    local seq_l2_block=$(curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "http://localhost:9545" 2>/dev/null | jq -r '.result' || echo "0x0")
    local seq_l2_num=$((16#${seq_l2_block#0x} || 0))

    # Challenger L2 Block Height
    local chal_l2_block=$(curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "http://localhost:9546" 2>/dev/null | jq -r '.result' || echo "0x0")
    local chal_l2_num=$((16#${chal_l2_block#0x} || 0))

    # Sequencer op-node sync status
    local seq_sync=$(curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"optimism_syncStatus","params":[],"id":1}' \
        "http://localhost:7545" 2>/dev/null | jq -r '.result' || echo "{}")

    # Challenger op-node sync status
    local chal_sync=$(curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"optimism_syncStatus","params":[],"id":1}' \
        "http://localhost:7546" 2>/dev/null | jq -r '.result' || echo "{}")

    echo ""
    echo "  ┌─ L1 Ethereum"
    if [ "$l1_num" -gt 0 ]; then
        echo -e "  │  └─ Block: ${GREEN}$l1_num${NC}"
    else
        echo -e "  │  └─ Block: ${RED}Not available${NC}"
    fi

    echo "  │"
    echo "  ├─ Sequencer Stack"
    if [ "$seq_l2_num" -gt 0 ]; then
        echo -e "  │  ├─ L2 Block: ${GREEN}$seq_l2_num${NC}"

        # Sequencer unsafe L2
        local seq_unsafe=$(echo "$seq_sync" | jq -r '.unsafe_l2.number // 0' 2>/dev/null || echo "0")
        if [ "$seq_unsafe" != "0" ] && [ "$seq_unsafe" != "null" ]; then
            echo "  │  └─ op-node Unsafe L2: $seq_unsafe"
        fi
    else
        echo -e "  │  └─ L2 Block: ${RED}Not synced${NC}"
    fi

    echo "  │"
    echo "  └─ Challenger Stack ⭐ (Verifier - L1 Data Only)"
    if [ "$chal_l2_num" -gt 0 ]; then
        echo -e "     ├─ L2 Block: ${GREEN}$chal_l2_num${NC}"

        # Compare with Sequencer
        if [ "$seq_l2_num" -gt 0 ] && [ "$chal_l2_num" -gt 0 ]; then
            local diff=$((seq_l2_num - chal_l2_num))
            if [ $diff -eq 0 ]; then
                echo -e "     │  └─ ${GREEN}✅ In sync with Sequencer${NC}"
            elif [ $diff -gt 0 ] && [ $diff -lt 50 ]; then
                echo -e "     │  └─ ${YELLOW}⏱  $diff blocks behind Sequencer (normal)${NC}"
            elif [ $diff -gt 0 ]; then
                echo -e "     │  └─ ${RED}⚠️  $diff blocks behind Sequencer (syncing)${NC}"
            else
                echo -e "     │  └─ ${RED}⚠️  Ahead of Sequencer (abnormal!)${NC}"
            fi
        fi

        # Challenger op-node status (Follower mode - Safe blocks only)
        local chal_safe=$(echo "$chal_sync" | jq -r '.safe_l2.number // 0' 2>/dev/null || echo "0")

        if [ "$chal_safe" != "0" ] && [ "$chal_safe" != "null" ]; then
            echo "     └─ op-node Safe L2: $chal_safe (verified from L1)"
        fi
    else
        echo -e "     └─ L2 Block: ${RED}Not synced${NC}"
    fi

    echo ""
}

# 9. L1 Batch Submission Status
show_batch_status() {
    log_title "L1 Batch Submission Status"

    local container=$(docker ps --format '{{.Names}}' | grep "op-batcher" | head -1 || echo "")

    if [ -n "$container" ]; then
        log_info "op-batcher container: $container"

        # Check recent batch submissions
        local recent_batches=$(docker logs --tail=50 "$container" 2>&1 | grep -i "batch submitted\|published" | tail -5)

        if [ -n "$recent_batches" ]; then
            echo ""
            echo "Recent batch submissions:"
            echo "$recent_batches" | sed 's/^/  /'
        else
            log_warn "No recent batch submission logs found"
            echo "  Check if op-batcher is actively submitting batches"
        fi

        # Check for errors
        local batcher_errors=$(docker logs --tail=50 "$container" 2>&1 | grep -i "error" | wc -l | tr -d ' ')

        if [ "$batcher_errors" -gt 0 ]; then
            log_warn "Batcher has $batcher_errors recent errors"
        else
            log_success "No recent batcher errors"
        fi
    else
        log_error "op-batcher container not found"
    fi

    echo ""
}

# 10. Independent Verification Status
show_independence_check() {
    log_title "Challenger Independence Verification"

    local challenger_container=$(docker ps -a --format '{{.Names}}' | grep "op-challenger" | head -1 || echo "")

    if [ -n "$challenger_container" ]; then
        # Check environment variables
        local rollup_rpc=$(docker inspect "$challenger_container" 2>/dev/null | jq -r '.[0].Config.Env[]' | grep "OP_CHALLENGER_ROLLUP_RPC" || echo "")
        local l2_rpc=$(docker inspect "$challenger_container" 2>/dev/null | jq -r '.[0].Config.Env[]' | grep "OP_CHALLENGER_L2_ETH_RPC" || echo "")

        if echo "$rollup_rpc" | grep -q "challenger-op-node"; then
            log_success "Using independent challenger-op-node ✅"
        else
            log_error "Using sequencer op-node! ❌"
            echo "  $rollup_rpc"
        fi

        if echo "$l2_rpc" | grep -q "challenger-l2"; then
            log_success "Using independent challenger-l2 ✅"
        else
            log_error "Using sequencer l2! ❌"
            echo "  $l2_rpc"
        fi

        # Check data directories (volumes)
        local volumes=$(docker inspect "$challenger_container" 2>/dev/null | jq -r '.[0].Mounts[] | select(.Destination=="/db" or .Destination=="/op-program" or .Destination=="/cannon") | "\(.Source) → \(.Destination)"')

        if [ -n "$volumes" ]; then
            echo ""
            echo "Mounted volumes:"
            echo "$volumes" | sed 's/^/  /'
        fi
    fi

    echo ""
}

# 9. Live Log Tail (interactive mode)
tail_logs() {
    log_title "Challenger Live Logs (Ctrl+C to exit)"
    echo ""

    local container=$(get_challenger_container)

    if [ -n "$container" ]; then
        docker logs -f --tail=50 "$container" 2>&1 | while IFS= read -r line; do
            # Color-code logs
            if echo "$line" | grep -q "lvl=error\|lvl=crit"; then
                echo -e "${RED}$line${NC}"
            elif echo "$line" | grep -q "lvl=warn"; then
                echo -e "${YELLOW}$line${NC}"
            elif echo "$line" | grep -q "Game info"; then
                echo -e "${CYAN}$line${NC}"
            elif echo "$line" | grep -q "Performing action"; then
                echo -e "${GREEN}$line${NC}"
            elif echo "$line" | grep -q "Resolving"; then
                echo -e "${GREEN}✅ $line${NC}"
            else
                echo "$line"
            fi
        done
    fi
}

##############################################################################
# Main script
##############################################################################

# Parse arguments
MODE="${1:-summary}"

case "$MODE" in
    summary)
        echo ""
        log_info "=========================================="
        log_info "Challenger Monitoring Dashboard"
        log_info "=========================================="
        echo ""

        show_container_status
        show_system_config
        show_sync_status
        show_batch_status
        show_independence_check
        show_game_summary
        show_active_games
        show_resolved_games
        show_actions
        show_metrics
        show_errors

        echo ""
        log_info "=========================================="
        log_info "Monitoring Options"
        log_info "=========================================="
        echo ""
        echo "  Live logs:      $0 logs"
        echo "  Games only:     $0 games"
        echo "  Errors only:    $0 errors"
        echo "  Metrics:        $0 metrics"
        echo "  Sync status:    $0 sync"
        echo "  Config:         $0 config"
        echo "  Full summary:   $0 summary (default)"
        echo ""
        ;;

    config)
        echo ""
        show_system_config
        ;;

    sync)
        echo ""
        show_sync_status
        show_batch_status
        ;;

    logs)
        tail_logs
        ;;

    games)
        echo ""
        show_game_summary
        show_active_games
        show_resolved_games
        show_actions
        ;;

    errors)
        echo ""
        show_errors
        ;;

    metrics)
        echo ""
        show_metrics
        ;;

    *)
        log_error "Unknown mode: $MODE"
        echo ""
        echo "Usage: $0 [summary|config|sync|logs|games|errors|metrics]"
        echo ""
        echo "Modes:"
        echo "  summary  - Full dashboard with all information (default)"
        echo "  config   - System configuration (game settings, proposer, challenger)"
        echo "  sync     - Blockchain sync status (L1, Sequencer, Challenger)"
        echo "  logs     - Live log tail (Ctrl+C to exit)"
        echo "  games    - Game participation details"
        echo "  errors   - Error analysis"
        echo "  metrics  - Prometheus metrics"
        echo ""
        echo "Examples:"
        echo "  $0              # Show full dashboard"
        echo "  $0 config       # Check system configuration"
        echo "  $0 sync         # Check sync status only"
        echo "  $0 logs         # Watch live logs"
        echo "  $0 games        # See game activity"
        echo ""
        exit 1
        ;;
esac

exit 0

