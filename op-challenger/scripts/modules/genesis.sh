#!/usr/bin/env bash

##############################################################################
# Genesis Generation Module
#
# This module handles Genesis file generation and game configuration:
# - Cleanup of old Genesis files and deployment cache
# - Game configuration (clock duration, withdrawal delay)
# - Prestate hash injection into deploy config template
# - devnet-up orchestration
# - Contract address extraction and verification
##############################################################################

set -euo pipefail

# Source common library
# Note: Use local variable to avoid polluting parent script's SCRIPT_DIR
_MODULE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${_MODULE_DIR}/../lib/common.sh"

##############################################################################
# Configuration
##############################################################################

DEVNET_DIR="${PROJECT_ROOT}/.devnet"
GENESIS_L1="${DEVNET_DIR}/genesis-l1.json"
GENESIS_L2="${DEVNET_DIR}/genesis-l2.json"
ROLLUP_CONFIG="${DEVNET_DIR}/rollup.json"
ADDRESSES_FILE="${DEVNET_DIR}/addresses.json"

DEPLOY_CONFIG_TEMPLATE="${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deploy-config/devnetL1-template.json"
DEPLOY_CONFIG="${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deploy-config/devnetL1.json"
DEPLOYMENT_CACHE="${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deployments/devnetL1"

##############################################################################
# Cleanup Functions
##############################################################################

cleanup_docker_volume_directories() {
    log_info "Checking for directories incorrectly created by Docker volume mounts..."

    local cleanup_needed=false
    for genesis_file in "$GENESIS_L1" "$GENESIS_L2" "$ROLLUP_CONFIG"; do
        if [ -d "$genesis_file" ]; then
            log_warn "Found directory instead of file: $genesis_file"
            log_warn "Removing incorrectly created directory..."
            rm -rf "$genesis_file"
            cleanup_needed=true
        fi
    done

    if [ "$cleanup_needed" = true ]; then
        log_success "Cleaned up incorrectly created directories"
        log_info "This usually happens when Docker volume mount creates directories for missing files"
        echo ""
    fi
}

cleanup_for_config_change() {
    log_warn "⚙️  Game configuration environment variables detected!"
    log_info "  FAULT_GAME_MAX_CLOCK_DURATION: ${FAULT_GAME_MAX_CLOCK_DURATION:-default}"
    log_info "  FAULT_GAME_WITHDRAWAL_DELAY: ${FAULT_GAME_WITHDRAWAL_DELAY:-default}"
    echo ""
    log_warn "Will apply custom game settings during Genesis generation."
    echo ""

    log_info "🧹 Removing old files to apply new game settings..."
    echo ""

    # Remove deployment cache (critical!)
    log_info "  1️⃣  Removing deployment cache..."
    if [ -d "$DEPLOYMENT_CACHE" ]; then
        rm -rf "$DEPLOYMENT_CACHE"
        echo "     ✅ Removed: packages/tokamak/contracts-bedrock/deployments/devnetL1"
    else
        echo "     ℹ️  Not found: deployments/devnetL1"
    fi

    # Remove devnetL1.json (critical!) - forces template to be used
    log_info "  2️⃣  Removing devnetL1.json (forces template to be used)..."
    if [ -f "$DEPLOY_CONFIG" ]; then
        local old_clock=$(jq -r '.faultGameMaxClockDuration' "$DEPLOY_CONFIG" 2>/dev/null || echo "unknown")
        local old_delay=$(jq -r '.faultGameWithdrawalDelay' "$DEPLOY_CONFIG" 2>/dev/null || echo "unknown")
        echo "     📋 Old config: clock=${old_clock}s, delay=${old_delay}s"
        rm -f "$DEPLOY_CONFIG"
        echo "     ✅ Removed: devnetL1.json"
    else
        echo "     ℹ️  Not found: devnetL1.json"
    fi

    # Remove existing Genesis files
    log_info "  3️⃣  Removing Genesis and contract files..."
    rm -f "$GENESIS_L1" "$GENESIS_L2" "$ROLLUP_CONFIG" "$ADDRESSES_FILE"
    rm -f "${DEVNET_DIR}/allocs"*.json
    echo "     ✅ Removed: genesis*.json, rollup.json, addresses.json, allocs*.json"

    # Remove Docker volumes to prevent Genesis hash mismatch
    log_info "  4️⃣  Removing Docker volumes..."
    local volumes_removed=0
    for vol in scripts_l1_data scripts_sequencer_l2_data scripts_challenger_l2_data \
               scripts_sequencer_safedb_data scripts_challenger_safedb_data \
               scripts_challenger_data scripts_op_log; do
        if docker volume rm "$vol" 2>/dev/null; then
            volumes_removed=$((volumes_removed + 1))
        fi
    done
    echo "     ✅ Removed ${volumes_removed} Docker volume(s)"
    echo ""

    log_success "✅ Cleanup complete. Will regenerate with custom game settings."
    echo ""
}

##############################################################################
# Configuration Functions
##############################################################################

configure_game_settings() {
    local prestate_hash="$1"
    local dg_type="${2:-0}"  # ← GameType 추가 (기본값: 0)

    log_config "=== Template 파일 경로 ==="
    log_config "DEPLOY_CONFIG: $DEPLOY_CONFIG_TEMPLATE"
    log_config "File exists: $([ -f "$DEPLOY_CONFIG_TEMPLATE" ] && echo 'YES' || echo 'NO')"
    log_config "GameType to configure: $dg_type"

    if [ ! -f "$DEPLOY_CONFIG_TEMPLATE" ]; then
        log_error "Deploy config template not found: $DEPLOY_CONFIG_TEMPLATE"
        return 1
    fi

    # Backup original config
    if [ ! -f "${DEPLOY_CONFIG_TEMPLATE}.backup" ]; then
        cp "$DEPLOY_CONFIG_TEMPLATE" "${DEPLOY_CONFIG_TEMPLATE}.backup"
        log_info "Backup created: ${DEPLOY_CONFIG_TEMPLATE}.backup"
    fi

    local config_modified=false

    log_config "=== Template 수정 전 상태 ==="
    log_config "Original values:"
    local orig_clock=$(jq -r '.faultGameMaxClockDuration' "$DEPLOY_CONFIG_TEMPLATE" 2>/dev/null || echo "unknown")
    local orig_delay=$(jq -r '.faultGameWithdrawalDelay' "$DEPLOY_CONFIG_TEMPLATE" 2>/dev/null || echo "unknown")
    local orig_gametype=$(jq -r '.respectedGameType' "$DEPLOY_CONFIG_TEMPLATE" 2>/dev/null || echo "unknown")
    log_config "  faultGameMaxClockDuration: ${orig_clock}s"
    log_config "  faultGameWithdrawalDelay: ${orig_delay}s"
    log_config "  respectedGameType: ${orig_gametype}"

    # Update clock duration if specified
    if [ -n "${FAULT_GAME_MAX_CLOCK_DURATION:-}" ]; then
        log_info "Configuring faultGameMaxClockDuration: ${FAULT_GAME_MAX_CLOCK_DURATION}s ($((FAULT_GAME_MAX_CLOCK_DURATION / 60)) min)"
        log_config "Modifying faultGameMaxClockDuration: ${orig_clock} -> ${FAULT_GAME_MAX_CLOCK_DURATION}"

        if command -v jq >/dev/null 2>&1; then
            local temp_config=$(mktemp)
            jq ".faultGameMaxClockDuration = ${FAULT_GAME_MAX_CLOCK_DURATION}" "$DEPLOY_CONFIG_TEMPLATE" > "$temp_config"
            mv "$temp_config" "$DEPLOY_CONFIG_TEMPLATE"
            config_modified=true

            # Verify modification
            local new_clock=$(jq -r '.faultGameMaxClockDuration' "$DEPLOY_CONFIG_TEMPLATE")
            log_config "Verification: faultGameMaxClockDuration = ${new_clock}s"
            if [ "$new_clock" = "${FAULT_GAME_MAX_CLOCK_DURATION}" ]; then
                log_config "✅ faultGameMaxClockDuration successfully updated"
            else
                log_config "❌ faultGameMaxClockDuration update FAILED!"
            fi
        else
            log_warn "jq not found. Cannot modify faultGameMaxClockDuration"
            log_config "❌ jq not available - modification skipped"
        fi
    fi

    # Update withdrawal delay if specified
    if [ -n "${FAULT_GAME_WITHDRAWAL_DELAY:-}" ]; then
        log_info "Configuring faultGameWithdrawalDelay: ${FAULT_GAME_WITHDRAWAL_DELAY}s ($((FAULT_GAME_WITHDRAWAL_DELAY / 86400)) days)"
        log_config "Modifying faultGameWithdrawalDelay: ${orig_delay} -> ${FAULT_GAME_WITHDRAWAL_DELAY}"

        if command -v jq >/dev/null 2>&1; then
            local temp_config=$(mktemp)
            jq ".faultGameWithdrawalDelay = ${FAULT_GAME_WITHDRAWAL_DELAY}" "$DEPLOY_CONFIG_TEMPLATE" > "$temp_config"
            mv "$temp_config" "$DEPLOY_CONFIG_TEMPLATE"
            config_modified=true

            # Verify modification
            local new_delay=$(jq -r '.faultGameWithdrawalDelay' "$DEPLOY_CONFIG_TEMPLATE")
            log_config "Verification: faultGameWithdrawalDelay = ${new_delay}s"
            if [ "$new_delay" = "${FAULT_GAME_WITHDRAWAL_DELAY}" ]; then
                log_config "✅ faultGameWithdrawalDelay successfully updated"
            else
                log_config "❌ faultGameWithdrawalDelay update FAILED!"
            fi
        fi
    fi

    # Log proposal interval if specified
    if [ -n "${PROPOSAL_INTERVAL:-}" ]; then
        log_info "Proposal interval will be set to: ${PROPOSAL_INTERVAL}"
        log_config "PROPOSAL_INTERVAL: ${PROPOSAL_INTERVAL}"
    fi

    # ⭐ Update respectedGameType (OptimismPortal2가 사용할 GameType)
    if [ -n "$dg_type" ]; then
        log_info "Updating respectedGameType for OptimismPortal2..."
        log_config "Updating respectedGameType: ${orig_gametype} -> ${dg_type}"

        local temp_config=$(mktemp)
        jq ".respectedGameType = ${dg_type}" "$DEPLOY_CONFIG_TEMPLATE" > "$temp_config"
        mv "$temp_config" "$DEPLOY_CONFIG_TEMPLATE"

        # Verify
        local new_gametype=$(jq -r '.respectedGameType' "$DEPLOY_CONFIG_TEMPLATE")
        log_config "Verification: respectedGameType = ${new_gametype}"
        if [ "$new_gametype" = "$dg_type" ]; then
            log_config "✅ respectedGameType successfully updated"
            log_success "OptimismPortal2 will use GameType $dg_type"
        else
            log_config "❌ respectedGameType update FAILED!"
            log_error "Failed to update respectedGameType"
            return 1
        fi
        config_modified=true
    fi

    # ⭐ Update template with generated prestate hash
    log_info "=========================================="
    log_info "Setting Prestate Hash for GameType $dg_type"
    log_info "=========================================="

    if [ -z "$prestate_hash" ] || [ "$prestate_hash" = "error" ]; then
        log_error "❌ Prestate hash is empty or invalid!"
        log_error "   Provided: '$prestate_hash'"
        log_error "   GameType: $dg_type"
        return 1
    fi

    log_info "Prestate hash to set: $prestate_hash"
    log_info "GameType: $dg_type"

    # Verify prestate file exists
    case "$dg_type" in
        0|1|254|255)
            local prestate_file="${PROJECT_ROOT}/op-program/bin/prestate-proof.json"
            log_info "Source file: $prestate_file (Cannon)"
            ;;
        2)
            local prestate_file="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
            log_info "Source file: $prestate_file (Asterisc)"
            ;;
        3)
            local prestate_file="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
            log_info "Source file: $prestate_file (Kona)"
            if [ ! -f "$prestate_file" ]; then
                log_warn "Kona prestate not found, using Asterisc fallback"
                prestate_file="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
                log_info "Fallback file: $prestate_file"
            fi
            ;;
    esac

    if [ -f "$prestate_file" ]; then
        log_success "✅ Prestate source file exists"
        local file_hash=$(cat "$prestate_file" | jq -r '.pre' 2>/dev/null || echo "")
        log_info "   File hash: ${file_hash:0:18}...${file_hash: -6}"

        if [ "$file_hash" = "$prestate_hash" ]; then
            log_success "✅ Hash matches expected value"
        else
            log_error "❌ Hash mismatch!"
            log_error "   Expected: $prestate_hash"
            log_error "   File has: $file_hash"
            return 1
        fi
    else
        log_error "❌ Prestate source file not found: $prestate_file"
        return 1
    fi
    echo ""

    log_info "Updating deploy config template..."
    log_config "Updating faultGameAbsolutePrestate: $prestate_hash"

    local temp_config=$(mktemp)
    jq ".faultGameAbsolutePrestate = \"${prestate_hash}\"" "$DEPLOY_CONFIG_TEMPLATE" > "$temp_config"
    mv "$temp_config" "$DEPLOY_CONFIG_TEMPLATE"

    # Verify
    local new_prestate=$(jq -r '.faultGameAbsolutePrestate' "$DEPLOY_CONFIG_TEMPLATE")
    log_config "Verification: faultGameAbsolutePrestate = ${new_prestate}"
    if [ "$new_prestate" = "$prestate_hash" ]; then
        log_config "✅ faultGameAbsolutePrestate successfully updated"
        log_success "Template updated with GameType $dg_type prestate"
    else
        log_config "❌ faultGameAbsolutePrestate update FAILED!"
        log_error "Failed to update template with prestate"
        log_error "   Expected: $prestate_hash"
        log_error "   Got: $new_prestate"
        return 1
    fi
    config_modified=true
    echo ""

    if [ "$config_modified" = true ]; then
        log_success "Game configuration updated in template"
        log_config "=== Template 수정 완료 ==="
        log_config "Modified template file: $DEPLOY_CONFIG_TEMPLATE"
        log_config "GameType ${dg_type}: respectedGameType and prestate are now aligned"
        echo ""
    fi

    return 0
}

##############################################################################
# Genesis Generation Functions
##############################################################################

run_devnet_up() {
    log_info "Starting devnet (make devnet-up)..."
    log_info "This will take 2-3 minutes. Please wait..."
    echo ""

    cd "${PROJECT_ROOT}"

    # ⭐ Pre-flight checks: Verify prestate files exist before starting devnet
    log_info "=========================================="
    log_info "Pre-flight: Verifying prestate files..."
    log_info "=========================================="

    local preflight_failed=0

    # Check asterisc prestate (required for GameType 2/3)
    local asterisc_prestate="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
    if [ -f "$asterisc_prestate" ]; then
        local asterisc_hash=$(cat "$asterisc_prestate" | jq -r '.pre' 2>/dev/null || echo "")
        if [ -n "$asterisc_hash" ] && [ "$asterisc_hash" != "null" ]; then
            log_success "✅ Asterisc prestate: ${asterisc_hash:0:18}...${asterisc_hash: -6}"
        else
            log_error "❌ Asterisc prestate file exists but .pre field is invalid!"
            preflight_failed=1
        fi
    else
        log_error "❌ Asterisc prestate file not found: $asterisc_prestate"
        preflight_failed=1
    fi

    # Check kona prestate (if GameType 3)
    local kona_prestate="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
    if [ -f "$kona_prestate" ]; then
        local kona_hash=$(cat "$kona_prestate" | jq -r '.pre' 2>/dev/null || echo "")
        if [ -n "$kona_hash" ] && [ "$kona_hash" != "null" ]; then
            log_success "✅ Kona prestate: ${kona_hash:0:18}...${kona_hash: -6}"
        else
            log_warn "⚠️  Kona prestate file exists but .pre field is invalid"
        fi
    else
        log_warn "⚠️  Kona prestate not found (will fallback to Asterisc for GameType 3)"
    fi

    if [ $preflight_failed -eq 1 ]; then
        log_error "=========================================="
        log_error "Pre-flight checks FAILED!"
        log_error "=========================================="
        log_error "Prestate files are missing or invalid."
        log_error "Deploy.s.sol will fail when trying to load prestate!"
        echo ""
        log_error "Please run:"
        log_error "  ./op-challenger/scripts/pull-vm-images.sh --tag latest"
        echo ""
        return 1
    fi

    log_success "✅ Pre-flight checks passed!"
    echo ""

    # Temporarily patch Dockerfile to add missing asterisc-builder stage
    # This is needed because we build op-challenger but provide asterisc via volume mount
    local dockerfile="${PROJECT_ROOT}/ops/docker/op-stack-go/Dockerfile"
    local dockerfile_backup="${dockerfile}.backup"

    # Check if asterisc-builder is missing
    if ! grep -q "as asterisc-builder" "$dockerfile"; then
        log_info "Patching Dockerfile to add asterisc-builder stage..."

        # Backup original Dockerfile
        cp "$dockerfile" "$dockerfile_backup"

        # Add asterisc-builder stage after cannon-builder
        sed -i '' '/^FROM.*as op-program-builder/i\
\
# Dummy asterisc-builder stage (actual binary provided via volume mount)\
FROM --platform=$BUILDPLATFORM builder as asterisc-builder\
RUN mkdir -p /app/asterisc/bin && \\\
    echo "#!/bin/sh" > /app/asterisc/bin/asterisc && \\\
    echo "echo asterisc-placeholder" >> /app/asterisc/bin/asterisc && \\\
    chmod +x /app/asterisc/bin/asterisc\
' "$dockerfile"

        log_success "Dockerfile patched for asterisc-builder"
    fi

    # Run devnet-up
    log_info "=========================================="
    log_info "Running: make devnet-up"
    log_info "=========================================="
    echo ""

    local result=0
    if make devnet-up 2>&1 | tee /tmp/devnet-up.log; then
        log_success "Devnet started successfully"
        echo ""
    else
        log_error "=========================================="
        log_error "Devnet startup FAILED!"
        log_error "=========================================="
        log_error "Last 50 lines of devnet-up output:"
        echo ""
        tail -50 /tmp/devnet-up.log
        echo ""
        log_error "Full log saved to: /tmp/devnet-up.log"
        result=1
    fi

    # Restore original Dockerfile if we patched it
    if [ -f "$dockerfile_backup" ]; then
        log_info "Restoring original Dockerfile..."
        mv "$dockerfile_backup" "$dockerfile"
    fi

    return $result
}

wait_for_genesis_files() {
    log_info "Waiting for Genesis and contract files to be created..."

    local max_wait=180  # 3 minutes
    local elapsed=0

    while [ $elapsed -lt $max_wait ]; do
        if [ -f "$GENESIS_L1" ] && [ -f "$GENESIS_L2" ] && [ -f "$ROLLUP_CONFIG" ] && [ -f "$ADDRESSES_FILE" ]; then
            log_success "Genesis and contract files detected!"
            echo ""
            return 0
        fi
        sleep 5
        elapsed=$((elapsed + 5))
        echo -n "."
    done
    echo ""

    log_error "Genesis files were not created within ${max_wait} seconds"
    return 1
}

verify_genesis_files() {
    log_success "All Genesis and contract files created successfully!"
    echo ""
    log_info "Generated files:"
    ls -lh "$GENESIS_L1" "$GENESIS_L2" "$ROLLUP_CONFIG" "$ADDRESSES_FILE" 2>/dev/null || true
    echo ""

    # Verify devnetL1.json was created with correct values
    log_config "=== devnetL1.json 생성 확인 (make devnet-up 후) ==="
    if [ -f "$DEPLOY_CONFIG" ]; then
        local actual_clock=$(jq -r '.faultGameMaxClockDuration' "$DEPLOY_CONFIG" 2>/dev/null || echo "unknown")
        local actual_delay=$(jq -r '.faultGameWithdrawalDelay' "$DEPLOY_CONFIG" 2>/dev/null || echo "unknown")
        log_config "devnetL1.json values:"
        log_config "  faultGameMaxClockDuration: ${actual_clock}s"
        log_config "  faultGameWithdrawalDelay: ${actual_delay}s"

        # Compare with expected values
        if [ -n "${FAULT_GAME_MAX_CLOCK_DURATION:-}" ]; then
            if [ "$actual_clock" = "${FAULT_GAME_MAX_CLOCK_DURATION}" ]; then
                log_config "✅ Clock duration matches expected: ${FAULT_GAME_MAX_CLOCK_DURATION}s"
            else
                log_config "❌ Clock duration mismatch! Expected: ${FAULT_GAME_MAX_CLOCK_DURATION}s, Actual: ${actual_clock}s"
            fi
        fi

        if [ -n "${FAULT_GAME_WITHDRAWAL_DELAY:-}" ]; then
            if [ "$actual_delay" = "${FAULT_GAME_WITHDRAWAL_DELAY}" ]; then
                log_config "✅ Withdrawal delay matches expected: ${FAULT_GAME_WITHDRAWAL_DELAY}s"
            else
                log_config "❌ Withdrawal delay mismatch! Expected: ${FAULT_GAME_WITHDRAWAL_DELAY}s, Actual: ${actual_delay}s"
            fi
        fi
    else
        log_config "❌ devnetL1.json not found at: $DEPLOY_CONFIG"
    fi

    # Verify deployed contract settings
    log_config "=== Deployed Contract 설정 확인 (addresses.json) ==="
    if [ -f "$ADDRESSES_FILE" ]; then
        local dgf_addr=$(jq -r '.DisputeGameFactoryProxy // empty' "$ADDRESSES_FILE" 2>/dev/null || echo "")
        log_config "DisputeGameFactoryProxy: ${dgf_addr}"
        if [ -n "$dgf_addr" ] && [ "$dgf_addr" != "null" ]; then
            log_config "✅ Contract deployed successfully"
        else
            log_config "❌ DisputeGameFactoryProxy address not found!"
        fi
    fi

    return 0
}

stop_devnet() {
    log_info "Stopping devnet and removing volumes to free ports for scripts stack..."
    log_warn "Removing devnet volumes to ensure L1 uses the same genesis hash..."

    cd "${PROJECT_ROOT}/ops-bedrock"

    # Determine docker compose command
    local docker_compose
    if docker compose version &> /dev/null; then
        docker_compose="docker compose"
    elif command -v docker-compose &> /dev/null; then
        docker_compose="docker-compose"
    else
        log_error "Neither 'docker compose' nor 'docker-compose' is available"
        return 1
    fi

    # Try to stop using docker-compose.yml
    if [ -f "docker-compose.yml" ]; then
        log_info "Attempting to stop ops-bedrock using docker compose..."

        # Don't suppress errors - we need to see them
        if ! $docker_compose -f docker-compose.yml down -v 2>&1 | tee -a "${DEPLOY_LOG:-/dev/null}"; then
            log_warn "docker compose down failed, attempting manual cleanup..."

            # Manual cleanup as fallback
            log_info "Stopping containers manually..."
            docker ps -a --format '{{.Names}}' | grep "^ops-bedrock-" | xargs -r docker stop 2>/dev/null || true
            docker ps -a --format '{{.Names}}' | grep "^ops-bedrock-" | xargs -r docker rm 2>/dev/null || true

            log_info "Removing volumes manually..."
            docker volume ls --format '{{.Name}}' | grep "^ops-bedrock_" | xargs -r docker volume rm 2>/dev/null || true
        fi
    fi

    cd "${PROJECT_ROOT}"

    # ⭐ CRITICAL: Verify containers are actually stopped
    log_info "Verifying ops-bedrock containers are stopped..."
    if docker ps --format '{{.Names}}' | grep "^ops-bedrock-" > /dev/null 2>&1; then
        log_error "❌ ops-bedrock containers are still running!"
        log_error "This will cause port conflicts with scripts stack"
        echo ""
        log_error "Running containers:"
        docker ps --filter "name=ops-bedrock-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
        echo ""
        log_error "Please stop them manually:"
        log_error "  cd ${PROJECT_ROOT}/ops-bedrock"
        log_error "  docker compose down -v"
        echo ""
        return 1
    fi

    # ⭐ Verify critical ports are freed
    log_info "Verifying critical ports are freed..."
    local port_conflict=false

    for port in 8545 7060; do
        if lsof -i :${port} 2>/dev/null | grep -v "PID" > /dev/null; then
            log_error "❌ Port ${port} is still in use!"
            lsof -i :${port} 2>/dev/null || true
            port_conflict=true
        fi
    done

    if [ "$port_conflict" = true ]; then
        log_error "Port conflicts detected. Cannot proceed with deployment."
        return 1
    fi

    log_success "✅ ops-bedrock fully stopped and verified"
    log_success "✅ All critical ports freed (8545, 7060)"
    log_info "Ready to start scripts stack without conflicts"
    echo ""

    return 0
}

##############################################################################
# Main Genesis Generation Function
##############################################################################

generate_genesis() {
    local prestate_hash="${1:-}"
    local dg_type="${2:-0}"  # ← GameType 추가 (기본값 0)

    log_info "=========================================="
    log_info "Genesis Generation Module"
    log_info "=========================================="
    log_info "GameType: $dg_type"
    echo ""

    # Ensure devnet directory exists
    ensure_dir "$DEVNET_DIR"

    # Step 1: Cleanup Docker volume directories
    cleanup_docker_volume_directories

    # Step 2: Check if custom game config requires cleanup
    if [ -n "${FAULT_GAME_MAX_CLOCK_DURATION:-}" ] || [ -n "${FAULT_GAME_WITHDRAWAL_DELAY:-}" ]; then
        cleanup_for_config_change
    fi

    # Step 3: Check if Genesis files exist
    if [ -f "$GENESIS_L1" ] && [ -f "$GENESIS_L2" ] && [ -f "$ROLLUP_CONFIG" ] && [ -f "$ADDRESSES_FILE" ]; then
        log_success "Genesis and contract files already exist"
        log_info "Existing files:"
        ls -lh "$GENESIS_L1" "$GENESIS_L2" "$ROLLUP_CONFIG" "$ADDRESSES_FILE"
        echo ""
        log_info "Delete these files to force regeneration"
        return 0
    fi

    # Step 4: Genesis files missing, need to generate
    log_warn "Required Genesis and contract files are missing!"
    echo ""
    log_info "Missing files:"
    [ ! -f "$GENESIS_L1" ] && echo "  ❌ ${GENESIS_L1}"
    [ ! -f "$GENESIS_L2" ] && echo "  ❌ ${GENESIS_L2}"
    [ ! -f "$ROLLUP_CONFIG" ] && echo "  ❌ ${ROLLUP_CONFIG}"
    [ ! -f "$ADDRESSES_FILE" ] && echo "  ❌ ${ADDRESSES_FILE} (contract addresses)"
    echo ""

    log_info "Automatically generating Genesis and contract files..."
    log_warn "This will:"
    echo "  1. Configure game settings with prestate hash"
    echo "  2. Start devnet (make devnet-up)"
    echo "  3. Deploy contracts and generate Genesis (~2-3 minutes)"
    echo "  4. Stop devnet (make devnet-down)"
    echo ""

    # Validate prestate hash
    if [ -z "$prestate_hash" ]; then
        log_error "Prestate hash is required for Genesis generation"
        log_error "Please build VMs first to generate prestate hash"
        return 1
    fi

    # Step 5: Configure game settings
    if ! configure_game_settings "$prestate_hash" "$dg_type"; then
        log_error "Failed to configure game settings"
        return 1
    fi

    # Step 6: Run devnet-up
    if ! run_devnet_up; then
        log_error "Failed to start devnet"
        return 1
    fi

    # Step 7: Wait for Genesis files
    if ! wait_for_genesis_files; then
        log_error "Genesis files were not created"
        return 1
    fi

    # Step 8: Verify Genesis files
    verify_genesis_files

    # Step 9: Stop devnet
    stop_devnet

    log_success "=========================================="
    log_success "Genesis Generation Complete"
    log_success "=========================================="
    echo ""

    return 0
}

##############################################################################
# CLI Interface (if run directly)
##############################################################################

if [ "${BASH_SOURCE[0]}" -ef "$0" ]; then
    # Script is being run directly, not sourced

    PRESTATE_HASH=""

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --prestate)
                PRESTATE_HASH="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --prestate HASH    Absolute prestate hash (required)"
                echo "  --help             Show this help"
                echo ""
                echo "Environment Variables:"
                echo "  FAULT_GAME_MAX_CLOCK_DURATION    Max clock duration in seconds"
                echo "  FAULT_GAME_WITHDRAWAL_DELAY      Withdrawal delay in seconds"
                echo "  PROPOSAL_INTERVAL                Proposal interval"
                echo ""
                echo "Examples:"
                echo "  $0 --prestate 0x1234..."
                echo "  FAULT_GAME_MAX_CLOCK_DURATION=120 $0 --prestate 0x1234..."
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done

    # Initialize logging
    init_logging

    # Run Genesis generation
    if generate_genesis "$PRESTATE_HASH"; then
        log_success "Genesis generation completed successfully!"
        exit 0
    else
        log_error "Genesis generation failed"
        exit 1
    fi
fi
