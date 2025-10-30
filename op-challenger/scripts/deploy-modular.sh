#!/usr/bin/env bash

##############################################################################
# Modular Deployment Script
#
# This is a refactored version of deploy-full-stack.sh that uses modularized
# components. It demonstrates how to orchestrate the deployment using the
# new module system.
#
# Usage:
#   ./deploy-modular.sh [options]
#
# Options:
#   --mode local|existing      Deployment mode (default: local)
#   --dg-type 0|1|2|3|254|255  Dispute game type (default: 0)
#   --help                     Show this help
#
##############################################################################

set -euo pipefail

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
export PROJECT_ROOT

# Optional: Auto-logging to file
# Uncomment to enable automatic log file creation
# LOG_FILE="${SCRIPT_DIR}/deploy-$(date +%Y%m%d-%H%M%S).log"
# exec > >(tee -a "$LOG_FILE") 2>&1
# echo "[INFO] Logging to: $LOG_FILE"

# Source common library
source "${SCRIPT_DIR}/lib/common.sh"

# Source modules
source "${SCRIPT_DIR}/modules/cleanup.sh"
# vm-build.sh is NOT sourced - only used for manual local builds
# deploy-modular.sh expects pre-built binaries (from pull-vm-images.sh)
source "${SCRIPT_DIR}/modules/genesis.sh"
source "${SCRIPT_DIR}/modules/env-setup.sh"
source "${SCRIPT_DIR}/modules/docker-deploy.sh"

##############################################################################
# Configuration
##############################################################################

DEPLOY_MODE="${DEPLOY_MODE:-local}"
DG_TYPE="${DG_TYPE:-0}"
# CHALLENGER_TRACE_TYPE will be set automatically based on DG_TYPE
# by validate_gametype_tracetype() function
CHALLENGER_TRACE_TYPE="${CHALLENGER_TRACE_TYPE:-}"

# Docker compose command - check which version is available
if docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
elif command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
else
    echo "ERROR: Neither 'docker compose' nor 'docker-compose' is available"
    exit 1
fi

##############################################################################
# Argument Parsing
##############################################################################

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --mode)
                DEPLOY_MODE="$2"
                # Validate deploy mode
                if [[ "$DEPLOY_MODE" != "local" && "$DEPLOY_MODE" != "existing" ]]; then
                    log_error "Invalid deploy mode: $DEPLOY_MODE"
                    log_error "Valid options: local, existing"
                    exit 1
                fi
                shift 2
                ;;
            --dg-type)
                DG_TYPE="$2"
                # Validate GameType
                if [[ ! "$DG_TYPE" =~ ^(0|1|2|3|254|255)$ ]]; then
                    log_error "Invalid GameType: $DG_TYPE"
                    log_error "Valid options: 0 (Cannon), 1 (Permissioned), 2 (Asterisc), 3 (AsteriscKona), 254 (Fast), 255 (Alphabet)"
                    exit 1
                fi
                shift 2
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --mode local|existing      Deployment mode (default: local)"
                echo "  --dg-type 0|1|2|3|254|255  Dispute game type (default: 0)"
                echo "  --help                     Show this help"
                echo ""
                echo "Environment Variables:"
                echo "  FAULT_GAME_MAX_CLOCK_DURATION    Max clock duration in seconds"
                echo "  FAULT_GAME_WITHDRAWAL_DELAY      Withdrawal delay in seconds"
                echo "  PROPOSAL_INTERVAL                Proposal interval"
                echo ""
                echo "Examples:"
                echo "  $0                              # Deploy with default settings (GameType 0)"
                echo "  $0 --dg-type 2                  # Deploy with GameType 2 (Asterisc)"
                echo "  $0 --dg-type 3                  # Deploy with GameType 3 (AsteriscKona)"
                echo "  FAULT_GAME_MAX_CLOCK_DURATION=120 $0  # Custom clock duration"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
}

##############################################################################
# Main Deployment Flow
##############################################################################

##############################################################################
# Prestate Validation Functions
##############################################################################

# Validate prestate hash format (0x + 64 hex characters)
validate_prestate_hash_format() {
    local hash="$1"
    if [[ "$hash" =~ ^0x[0-9a-fA-F]{64}$ ]]; then
        return 0
    else
        return 1
    fi
}

# Validate prestate file for a specific GameType
validate_prestate_file() {
    local dg_type="$1"
    local file_path="$2"
    local field_name="$3"
    local game_name="$4"

    # 1. Check file exists
    if [ ! -f "$file_path" ]; then
        log_error "  ✗ Prestate file not found: $file_path"
        return 1
    fi
    log_success "  ✓ File exists: $file_path"

    # 2. Check JSON validity
    if ! jq empty "$file_path" 2>/dev/null; then
        log_error "  ✗ Invalid JSON format"
        return 1
    fi
    log_success "  ✓ Valid JSON format"

    # 3. Check field exists
    local hash=$(jq -r ".$field_name" "$file_path" 2>/dev/null)
    if [ -z "$hash" ] || [ "$hash" = "null" ]; then
        log_error "  ✗ Field '.$field_name' not found or null"
        return 1
    fi
    log_success "  ✓ Field '.$field_name' exists"

    # 4. Validate hash format
    if ! validate_prestate_hash_format "$hash"; then
        log_error "  ✗ Invalid hash format: $hash"
        log_error "     Expected: 0x[64 hex characters]"
        return 1
    fi
    log_success "  ✓ Valid hash format: ${hash:0:10}...${hash: -8}"

    # 5. Check file size (should be reasonable, not empty or corrupted)
    # Note: prestate-proof.json (deployment format) is ~70-90 bytes, which is normal
    local file_size=$(stat -f%z "$file_path" 2>/dev/null || stat -c%s "$file_path" 2>/dev/null)
    if [ "$file_size" -lt 50 ]; then
        log_error "  ✗ Prestate file too small ($file_size bytes), may be corrupted"
        return 1
    fi
    log_success "  ✓ File size OK: $file_size bytes"

    return 0
}

get_prestate_hash_for_gametype() {
    local dg_type="$1"
    local prestate_hash=""

    case "$dg_type" in
        0|1|254|255)
            # Cannon (MIPS) prestate
            local cannon_prestate="${PROJECT_ROOT}/op-program/bin/prestate-proof.json"
            if [ -f "$cannon_prestate" ]; then
                prestate_hash=$(jq -r '.pre' "$cannon_prestate" 2>/dev/null || echo "")
                if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                    # Log to stderr to avoid polluting stdout (function return value)
                    log_info "Using Cannon (MIPS) prestate: $prestate_hash" >&2
                else
                    log_error "Failed to extract prestate from $cannon_prestate" >&2
                fi
            else
                log_error "Cannon prestate file not found: $cannon_prestate" >&2
            fi
            ;;
        2)
            # Asterisc (RISC-V) prestate
            local asterisc_prestate="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
            if [ -f "$asterisc_prestate" ]; then
                # ✅ Read .pre field (모든 GameType 공통!)
                prestate_hash=$(jq -r '.pre' "$asterisc_prestate" 2>/dev/null || echo "")
                if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                    # Log to stderr to avoid polluting stdout (function return value)
                    log_info "Using Asterisc (RISC-V) prestate: $prestate_hash" >&2
                else
                    log_error "Failed to extract .pre from $asterisc_prestate" >&2
                fi
            else
                log_error "Asterisc prestate not found: $asterisc_prestate" >&2
                log_error "❌ GameType 2 requires RISC-V prestate!" >&2
                log_error "   Cannot use MIPS prestate for RISC-V architecture." >&2
                echo "" >&2
                log_error "Please ensure Asterisc VM build succeeded:" >&2
                log_error "  - Check VM build logs above" >&2
                log_error "  - Verify asterisc/bin/prestate-proof.json exists" >&2
                log_error "  - Or use GameType 0/1 (Cannon) instead" >&2
                # Do NOT fallback to Cannon prestate - return empty
            fi
            ;;
        3)
            # AsteriscKona (RISC-V + Rust) prestate
            # Note: GameType 3 shares RISCV.sol with GameType 2 but uses kona-client
            # First, try kona-specific prestate if it exists
            local kona_prestate="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
            if [ -f "$kona_prestate" ]; then
                # ✅ Read .pre field (모든 GameType 공통!)
                prestate_hash=$(jq -r '.pre' "$kona_prestate" 2>/dev/null || echo "")
                if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                    log_info "Using AsteriscKona (RISC-V + Rust) prestate: $prestate_hash" >&2
                else
                    log_error "Failed to extract .pre from $kona_prestate" >&2
                fi
            else
                # Fallback to Asterisc prestate (same RISCV.sol)
                log_warn "Kona-specific prestate not found: $kona_prestate" >&2
                log_info "Falling back to Asterisc (RISC-V) prestate (same RISCV.sol)" >&2
                local asterisc_prestate="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
                if [ -f "$asterisc_prestate" ]; then
                    # ✅ Read .pre field (모든 GameType 공통!)
                    prestate_hash=$(jq -r '.pre' "$asterisc_prestate" 2>/dev/null || echo "")
                    if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                        log_info "Using Asterisc (RISC-V) prestate for GameType 3: $prestate_hash" >&2
                    else
                        log_error "Failed to extract .pre from $asterisc_prestate" >&2
                    fi
                else
                    log_error "Neither kona nor asterisc prestate found for GameType 3!" >&2
                    log_error "Please ensure either:" >&2
                    log_error "  - op-program/bin/prestate-kona.json exists, or" >&2
                    log_error "  - asterisc/bin/prestate-proof.json exists" >&2
                fi
            fi
            ;;
    esac

    echo "$prestate_hash"
}

main() {
    parse_arguments "$@"

    # Initialize logging
    init_logging

    log_info "=========================================="
    log_info "Tokamak L2 Full Stack Deployment (Modular)"
    log_info "=========================================="
    echo ""

    # Check required dependencies
    log_info "Checking required dependencies..."
    local missing_deps=0

    # Check jq (required for JSON parsing)
    if ! command -v jq &> /dev/null; then
        log_error "jq is not installed (required for JSON parsing)"
        log_error "Install: brew install jq (macOS) or apt-get install jq (Linux)"
        missing_deps=1
    fi

    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        log_error "Install from: https://docs.docker.com/get-docker/"
        missing_deps=1
    elif ! docker info &> /dev/null; then
        log_error "Docker is not running"
        log_error "Please start Docker and try again"
        missing_deps=1
    fi

    # Check make
    if ! command -v make &> /dev/null; then
        log_error "make is not installed"
        missing_deps=1
    fi

    # Check openssl (for JWT generation)
    if ! command -v openssl &> /dev/null; then
        log_error "openssl is not installed (required for JWT secret generation)"
        missing_deps=1
    fi

    if [ $missing_deps -eq 1 ]; then
        log_error "Missing required dependencies. Please install them and try again."
        exit 1
    fi

    log_success "All required dependencies are available"
    echo ""

    log_info "Deployment Mode: $DEPLOY_MODE"
    log_info "Dispute Game Type: $DG_TYPE ($(get_gametype_name $DG_TYPE))"
    echo ""

    # Validate GameType and TraceType compatibility (must be before logging trace type!)
    validate_gametype_tracetype "$DG_TYPE" "$CHALLENGER_TRACE_TYPE"

    log_info "Trace Type: $CHALLENGER_TRACE_TYPE"
    echo ""

    # Note: CHALLENGER_TRACE_TYPE will be exported by export_docker_env_vars()
    # in docker-deploy.sh, so .env file update is not needed.
    # Exported shell variables take precedence over .env file in Docker Compose.

    # Step 0: Cleanup existing containers and devnet files
    cleanup_containers
    cleanup_devnet_files

    # Step 1: Check VM Binaries (or Build if needed)
    log_info "=========================================="
    log_info "Step 1: Check VM Binaries"
    log_info "=========================================="
    echo ""

    # Check required VM binaries
    local missing_binaries=0
    local cannon_bin="${PROJECT_ROOT}/cannon/bin/cannon"
    local asterisc_bin="${PROJECT_ROOT}/asterisc/bin/asterisc"
    local op_program_bin="${PROJECT_ROOT}/op-program/bin/op-program"
    local kona_bin="${PROJECT_ROOT}/bin/kona-client"

    # Check Cannon (required for all GameTypes)
    if [ -f "$cannon_bin" ]; then
        log_success "✓ Cannon binary found: $cannon_bin"
    else
        log_error "✗ Cannon binary missing: $cannon_bin"
        missing_binaries=1
    fi

    # Check Asterisc (required for GameType 2, 3)
    if [ -f "$asterisc_bin" ]; then
        log_success "✓ Asterisc binary found: $asterisc_bin"
    else
        log_error "✗ Asterisc binary missing: $asterisc_bin"
        missing_binaries=1
    fi

    # Check op-program (required for all GameTypes)
    if [ -f "$op_program_bin" ]; then
        log_success "✓ op-program binary found: $op_program_bin"
    else
        log_error "✗ op-program binary missing: $op_program_bin"
        missing_binaries=1
    fi

    # Check kona-client (required for Challenger to support GameType 3)
    # Note: Challenger must support ALL GameTypes, not just the proposer's choice!
    if [ -f "$kona_bin" ]; then
        log_success "✓ kona-client binary found: $kona_bin"
    else
        log_warn "✗ kona-client binary missing: $kona_bin"
        log_warn "  ⚠️  Challenger will NOT be able to respond to GameType 3 games!"
        log_warn "  Proposer GameType: $DG_TYPE"

        # Only fail if we're explicitly deploying GameType 3
        if [ "$DG_TYPE" = "3" ]; then
            log_error "  ❌ Cannot deploy GameType 3 without kona-client!"
            missing_binaries=1
        else
            log_warn "  ℹ️  This is OK if you don't plan to use GameType 3"
            log_warn "  Download if needed: ./op-challenger/scripts/pull-vm-images.sh --tag latest"
        fi
    fi

    echo ""

    # If binaries are missing, show download instructions and exit
    if [ $missing_binaries -eq 1 ]; then
        log_error "❌ Required VM binaries are missing!"
        echo ""
        log_error "Please download pre-built binaries using one of these methods:"
        echo ""
        log_error "Method 1: Download from registry (recommended, fast)"
        log_error "  ./op-challenger/scripts/pull-vm-images.sh --tag latest"
        echo ""
        log_error "Method 2: Build locally (slow, 20-30 minutes)"
        log_error "  cd ${PROJECT_ROOT}"
        log_error "  cd cannon && make cannon"
        log_error "  cd ../asterisc && make asterisc"
        log_error "  cd ../op-program && make op-program"
        if [ "$DG_TYPE" = "3" ]; then
            log_error "  # For GameType 3:"
            log_error "  ./op-challenger/scripts/modules/vm-build.sh --kona-only"
        fi
        echo ""
        exit 1
    fi

    log_success "All required VM binaries are present!"
    echo ""

    # Step 1.5: Validate Prestate Files
    log_info "=========================================="
    log_info "Step 1.5: Validate Prestate Files"
    log_info "=========================================="
    log_info "Note: Challenger supports ALL GameTypes, checking all prestates..."
    echo ""

    local validation_failed=0
    local proposer_gametype_validated=false

    # 1. Validate Cannon (GameType 0/1) prestate
    log_info "1️⃣  Cannon (GameType 0/1) prestate..."
    if validate_prestate_file "0" \
        "${PROJECT_ROOT}/op-program/bin/prestate-proof.json" \
        "pre" \
        "Cannon"; then
        log_success "✅ Cannon prestate validation passed"
        if [ "$DG_TYPE" = "0" ] || [ "$DG_TYPE" = "1" ] || [ "$DG_TYPE" = "254" ] || [ "$DG_TYPE" = "255" ]; then
            proposer_gametype_validated=true
        fi
    else
        if [ "$DG_TYPE" = "0" ] || [ "$DG_TYPE" = "1" ] || [ "$DG_TYPE" = "254" ] || [ "$DG_TYPE" = "255" ]; then
            log_error "❌ Cannon prestate validation failed (REQUIRED for GameType $DG_TYPE)"
            validation_failed=1
            proposer_gametype_validated=true
        else
            log_warn "⚠️  Cannon prestate validation failed (Challenger won't support GameType 0/1)"
        fi
    fi
    echo ""

    # 2. Validate Asterisc (GameType 2) prestate
    log_info "2️⃣  Asterisc (GameType 2) prestate..."
    if validate_prestate_file "2" \
        "${PROJECT_ROOT}/asterisc/bin/prestate-proof.json" \
        "pre" \
        "Asterisc"; then
        log_success "✅ Asterisc prestate validation passed"
        if [ "$DG_TYPE" = "2" ]; then
            proposer_gametype_validated=true
        fi
    else
        if [ "$DG_TYPE" = "2" ]; then
            log_error "❌ Asterisc prestate validation failed (REQUIRED for GameType $DG_TYPE)"
            validation_failed=1
            proposer_gametype_validated=true
        else
            log_warn "⚠️  Asterisc prestate validation failed (Challenger won't support GameType 2)"
        fi
    fi
    echo ""

    # 3. Validate AsteriscKona (GameType 3) prestate
    log_info "3️⃣  AsteriscKona (GameType 3) prestate..."
    local kona_prestate="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
    if [ -f "$kona_prestate" ]; then
        log_info "Validating AsteriscKona (RISC-V + Rust) prestate..."
        if validate_prestate_file "3" \
            "$kona_prestate" \
            "pre" \
            "AsteriscKona"; then
            log_success "✅ AsteriscKona prestate validation passed"
            if [ "$DG_TYPE" = "3" ]; then
                proposer_gametype_validated=true
            fi
        else
            if [ "$DG_TYPE" = "3" ]; then
                log_error "❌ AsteriscKona prestate validation failed"
                validation_failed=1
                proposer_gametype_validated=true
            else
                log_warn "⚠️  AsteriscKona prestate validation failed (Challenger won't support GameType 3)"
            fi
        fi
    else
        # Fallback to Asterisc prestate
        log_warn "Kona-specific prestate not found, checking Asterisc prestate (fallback)..."
        if [ -f "${PROJECT_ROOT}/asterisc/bin/prestate-proof.json" ]; then
            log_success "✅ Asterisc prestate available (fallback for GameType 3)"
            log_info "   Note: GameType 2/3 share the same RISCV.sol, so this is safe"
            if [ "$DG_TYPE" = "3" ]; then
                proposer_gametype_validated=true
            fi
        else
            if [ "$DG_TYPE" = "3" ]; then
                log_error "❌ Neither kona nor asterisc prestate found (REQUIRED for GameType 3)"
                validation_failed=1
                proposer_gametype_validated=true
            else
                log_warn "⚠️  No prestate for GameType 3 (Challenger won't support GameType 3)"
            fi
        fi
    fi
    echo ""

    # Summary
    log_info "=========================================="
    log_info "Prestate Validation Summary"
    log_info "=========================================="
    log_info "✅ Proposer will use: GameType $DG_TYPE"
    log_info "✅ Challenger will support: ALL available GameTypes"
    log_info "   (based on which prestates are present)"
    echo ""

    if [ $validation_failed -eq 1 ]; then
        log_error "=========================================="
        log_error "Prestate validation failed!"
        log_error "=========================================="
        echo ""
        log_error "Possible solutions:"
        log_error "  1. Download pre-built binaries (includes prestates):"
        log_error "     ./op-challenger/scripts/pull-vm-images.sh --tag latest"
        echo ""
        log_error "  2. Build prestates locally:"
        case "$DG_TYPE" in
            0|1|254|255)
                log_error "     cd op-program && make reproducible-prestate"
                ;;
            2)
                log_error "     ./op-challenger/scripts/modules/vm-build.sh --asterisc-only"
                ;;
            3)
                log_error "     ./op-challenger/scripts/modules/vm-build.sh --kona-only"
                log_error "     # Or use Asterisc prestate (fallback):"
                log_error "     ./op-challenger/scripts/modules/vm-build.sh --asterisc-only"
                ;;
        esac
        echo ""
        exit 1
    fi

    log_success "=========================================="
    log_success "All prestate validations passed!"
    log_success "=========================================="
    echo ""

    # Get the correct prestate hash for this GameType
    local prestate_hash=$(get_prestate_hash_for_gametype "$DG_TYPE")

    if [ -z "$prestate_hash" ] || [ "$prestate_hash" = "null" ] || [ "$prestate_hash" = "error" ]; then
        log_error "Failed to get prestate hash for GameType $DG_TYPE"
        log_error "Please check VM build logs above for details"
        exit 1
    fi

    log_success "Prestate hash for GameType $DG_TYPE: $prestate_hash"
    echo ""

    # Step 2: Prepare environment FILES FIRST (before Docker starts!)
    # This prevents Docker from creating directories when files don't exist
    log_info "=========================================="
    log_info "Step 2: Prepare Required Files"
    log_info "=========================================="
    echo ""

    if ! prepare_required_files; then
        log_error "Failed to prepare required files"
        exit 1
    fi

    # Step 3: Generate Genesis (if needed)
    # WARNING: This will start Docker (make devnet-up), so files must exist first!
    if [ "$DEPLOY_MODE" = "local" ]; then
        log_info "=========================================="
        log_info "Step 3: Generate Genesis"
        log_info "=========================================="
        echo ""

        if ! generate_genesis "$prestate_hash" "$DG_TYPE"; then
            log_error "Genesis generation failed"
            exit 1
        fi

        # Step 3.5: Verify Prestate and GameType in Contract Deployment Config
        log_info "=========================================="
        log_info "Step 3.5: Verify Config for Contract Deployment"
        log_info "=========================================="
        echo ""

        # Use the same config file that genesis.sh modified
        local deploy_config="${PROJECT_ROOT}/packages/tokamak/contracts-bedrock/deploy-config/devnetL1-template.json"
        if [ ! -f "$deploy_config" ]; then
            log_error "Deployment config not found: $deploy_config"
            log_error "This should have been modified by genesis generation"
            exit 1
        fi

        # 1. Verify respectedGameType
        log_info "1️⃣  Checking respectedGameType in deployment config..."
        local config_gametype=$(jq -r '.respectedGameType' "$deploy_config" 2>/dev/null || echo "")

        if [ -z "$config_gametype" ] || [ "$config_gametype" = "null" ]; then
            log_error "❌ respectedGameType not found in deployment config!"
            log_error "   Config: $deploy_config"
            exit 1
        fi

        log_success "  ✓ respectedGameType found in config"
        log_info "    Config value: $config_gametype"
        log_info "    Expected value: $DG_TYPE"

        # Verify it matches our GameType
        if [ "$config_gametype" = "$DG_TYPE" ]; then
            log_success "✅ respectedGameType correctly set!"
            log_success "   OptimismPortal2 will accept GameType $DG_TYPE"
        else
            log_error "❌ respectedGameType mismatch!"
            log_error "   Config has:  $config_gametype"
            log_error "   Expected:    $DG_TYPE"
            log_error ""
            log_error "This means OptimismPortal2 will reject proofs from GameType $DG_TYPE!"
            log_error "Deployment cannot continue."
            exit 1
        fi
        echo ""

        # 2. Verify faultGameAbsolutePrestate
        log_info "2️⃣  Checking faultGameAbsolutePrestate in deployment config..."
        local config_prestate=$(jq -r '.faultGameAbsolutePrestate' "$deploy_config" 2>/dev/null || echo "")

        if [ -z "$config_prestate" ] || [ "$config_prestate" = "null" ]; then
            log_error "❌ faultGameAbsolutePrestate not found in deployment config!"
            log_error "   Config: $deploy_config"
            exit 1
        fi

        log_success "  ✓ faultGameAbsolutePrestate found in config"
        log_info "    Config value: $config_prestate"
        log_info "    Expected value: $prestate_hash"

        # Verify it matches our prestate hash
        if [ "$config_prestate" = "$prestate_hash" ]; then
            log_success "✅ Prestate correctly set in deployment config!"
            log_success "   GameType $DG_TYPE contract will use the correct prestate hash"
        else
            log_error "❌ Prestate mismatch!"
            log_error "   Config has:  $config_prestate"
            log_error "   Expected:    $prestate_hash"
            log_error ""
            log_error "This means the contract will be deployed with WRONG prestate!"
            log_error "Deployment cannot continue."
            exit 1
        fi
        echo ""

        log_success "=========================================="
        log_success "Configuration Verification Complete!"
        log_success "=========================================="
        log_success "✅ GameType: $DG_TYPE"
        log_success "✅ respectedGameType: $config_gametype"
        log_success "✅ Prestate: ${prestate_hash:0:10}...${prestate_hash: -8}"
        log_success ""
        log_success "OptimismPortal2 and GameType $DG_TYPE are correctly aligned!"
        echo ""
    else
        log_info "=========================================="
        log_info "Step 3: Genesis Generation"
        log_info "=========================================="
        echo ""
        log_info "Skipping Genesis generation (existing mode)"
        echo ""
    fi

    # Step 4: Load environment variables
    log_info "=========================================="
    log_info "Step 4: Load Environment Variables"
    log_info "=========================================="
    echo ""

    if ! load_environment_variables; then
        log_error "Failed to load environment variables"
        exit 1
    fi

    # Step 5: Deploy services and run health check
    log_info "=========================================="
    log_info "Step 5: Deploy Docker Services"
    log_info "=========================================="
    echo ""

    if ! deploy_and_check "$DG_TYPE" true; then
        log_error "Deployment failed"
        exit 1
    fi
}

##############################################################################
# Execute Main
##############################################################################

main "$@"
