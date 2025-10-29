#!/usr/bin/env bash

##############################################################################
# Common Library - Logging, Colors, and Utility Functions
#
# This library provides reusable functions for all deployment scripts.
##############################################################################

# Color definitions
export RED='\033[0;31m'
export GREEN='\033[0;32m'
export YELLOW='\033[1;33m'
export BLUE='\033[0;34m'
export CYAN='\033[0;36m'
export NC='\033[0m' # No Color

# Log file setup (can be overridden by caller)
export LOG_DIR="${LOG_DIR:-${PROJECT_ROOT:-.}/.devnet/logs}"
export DEPLOY_LOG="${DEPLOY_LOG:-}"
export CONFIG_TRACE_LOG="${CONFIG_TRACE_LOG:-}"

##############################################################################
# Logging Functions
##############################################################################

# Initialize log files
init_logging() {
    local timestamp=$(date +%Y%m%d-%H%M%S)
    mkdir -p "$LOG_DIR"

    if [ -z "$DEPLOY_LOG" ]; then
        export DEPLOY_LOG="${LOG_DIR}/deploy-${timestamp}.log"
    fi

    if [ -z "$CONFIG_TRACE_LOG" ]; then
        export CONFIG_TRACE_LOG="${LOG_DIR}/config-trace-${timestamp}.log"
    fi

    log_info "Log files initialized:"
    log_info "  Deployment: $DEPLOY_LOG"
    log_info "  Config trace: $CONFIG_TRACE_LOG"
}

# Log to both console and file
log_info() {
    local msg="$1"
    echo -e "${BLUE}[INFO]${NC} $msg"
    if [ -n "$DEPLOY_LOG" ]; then
        echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
    fi
}

log_success() {
    local msg="$1"
    echo -e "${GREEN}[SUCCESS]${NC} $msg"
    if [ -n "$DEPLOY_LOG" ]; then
        echo "[SUCCESS] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
    fi
}

log_warn() {
    local msg="$1"
    echo -e "${YELLOW}[WARN]${NC} $msg"
    if [ -n "$DEPLOY_LOG" ]; then
        echo "[WARN] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
    fi
}

log_error() {
    local msg="$1"
    echo -e "${RED}[ERROR]${NC} $msg"
    if [ -n "$DEPLOY_LOG" ]; then
        echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') $msg" >> "$DEPLOY_LOG"
    fi
}

# Config trace log (config-trace.log only, no console)
log_config() {
    local msg="$1"
    if [ -n "$CONFIG_TRACE_LOG" ]; then
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] $msg" >> "$CONFIG_TRACE_LOG"
    fi
}

##############################################################################
# Command Validation
##############################################################################

check_command() {
    local cmd="$1"
    if ! command -v "$cmd" &> /dev/null; then
        log_error "$cmd is not installed."
        return 1
    fi
    return 0
}

check_required_commands() {
    local missing=0
    for cmd in "$@"; do
        if ! check_command "$cmd"; then
            missing=1
        fi
    done

    if [ $missing -eq 1 ]; then
        log_error "Required commands are missing. Please install them and try again."
        exit 1
    fi

    log_success "All required commands are available."
}

##############################################################################
# Utility Functions
##############################################################################

# Wait for a command to succeed with timeout
wait_for_command() {
    local timeout="$1"
    shift
    local cmd="$@"
    local elapsed=0

    while [ $elapsed -lt $timeout ]; do
        if eval "$cmd" &> /dev/null; then
            return 0
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done

    return 1
}

# Wait for RPC endpoint to be ready
wait_for_rpc() {
    local rpc_url="$1"
    local timeout="${2:-120}"
    local service_name="${3:-RPC}"

    log_info "Waiting for $service_name to be ready ($rpc_url)..."

    local elapsed=0
    while [ $elapsed -lt $timeout ]; do
        if curl -s -X POST -H "Content-Type: application/json" \
            --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
            "$rpc_url" > /dev/null 2>&1; then
            log_success "$service_name is ready"
            return 0
        fi
        sleep 2
        elapsed=$((elapsed + 2))
        echo -n "."
    done
    echo ""

    log_error "$service_name failed to start within ${timeout} seconds"
    return 1
}

# Load environment variables from file
load_env_file() {
    local env_file="$1"

    if [ ! -f "$env_file" ]; then
        log_error "Environment file not found: $env_file"
        return 1
    fi

    log_info "Loading environment variables: $env_file"
    set -a
    source "$env_file"
    set +a
    log_success "Environment variables loaded"
    return 0
}

# Validate required environment variables
validate_env_vars() {
    local missing=0
    for var in "$@"; do
        if [ -z "${!var:-}" ]; then
            log_error "Required environment variable not set: $var"
            missing=1
        fi
    done

    if [ $missing -eq 1 ]; then
        log_error "Required environment variables are missing."
        return 1
    fi

    return 0
}

##############################################################################
# File Operations
##############################################################################

# Create directory if it doesn't exist
ensure_dir() {
    local dir="$1"
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        log_info "Created directory: $dir"
    fi
}

# Backup file if it exists
backup_file() {
    local file="$1"
    local backup_suffix="${2:-.backup}"

    if [ -f "$file" ] && [ ! -f "${file}${backup_suffix}" ]; then
        cp "$file" "${file}${backup_suffix}"
        log_info "Backup created: ${file}${backup_suffix}"
    fi
}

##############################################################################
# GameType Utilities
##############################################################################

# Get GameType name from number
get_gametype_name() {
    local game_type="$1"

    case "$game_type" in
        0) echo "CANNON (MIPS VM)" ;;
        1) echo "PERMISSIONED_CANNON" ;;
        2) echo "ASTERISC (RISC-V VM)" ;;
        3) echo "ASTERISC_KONA (RISC-V + Rust)" ;;
        254) echo "FAST (Test)" ;;
        255) echo "ALPHABET (Test)" ;;
        *) echo "UNKNOWN ($game_type)" ;;
    esac
}

# Validate GameType and TraceType compatibility
validate_gametype_tracetype() {
    local game_type="${1:-0}"
    local trace_type="${2:-}"  # Allow empty string to trigger auto-detection

    # Auto-detect trace type based on GameType if not set or empty
    if [ -z "$trace_type" ]; then
        case "$game_type" in
            0|1) trace_type="cannon" ;;
            2) trace_type="asterisc" ;;
            3) trace_type="asterisc-kona" ;;
            254) trace_type="fast" ;;
            255) trace_type="alphabet" ;;
            *) trace_type="cannon" ;;  # Default fallback
        esac
        log_info "Auto-detected trace_type=$trace_type for GameType $game_type"
    fi

    # Validate and correct if needed
    case "$game_type" in
        0|1)
            if [ "$trace_type" != "cannon" ] && [ "$trace_type" != "alphabet" ]; then
                log_warn "GameType $game_type requires trace_type=cannon, but got: $trace_type"
                log_info "Auto-correcting to trace_type=cannon"
                trace_type="cannon"
            fi
            ;;
        2)
            if [ "$trace_type" != "asterisc" ]; then
                log_warn "GameType 2 requires trace_type=asterisc, but got: $trace_type"
                log_info "Auto-correcting to trace_type=asterisc"
                trace_type="asterisc"
            fi
            ;;
        3)
            if [ "$trace_type" != "asterisc-kona" ]; then
                log_warn "GameType 3 (AsteriscKona) requires trace_type=asterisc-kona, but got: $trace_type"
                log_info "✅ Auto-correcting to trace_type=asterisc-kona"
                trace_type="asterisc-kona"
            fi
            ;;
        254)
            if [ "$trace_type" != "fast" ]; then
                log_warn "GameType 254 requires trace_type=fast, but got: $trace_type"
                log_info "Auto-correcting to trace_type=fast"
                trace_type="fast"
            fi
            ;;
        255)
            if [ "$trace_type" != "alphabet" ]; then
                log_warn "GameType 255 requires trace_type=alphabet, but got: $trace_type"
                log_info "Auto-correcting to trace_type=alphabet"
                trace_type="alphabet"
            fi
            ;;
    esac

    # IMPORTANT: Always export the validated/corrected value
    export CHALLENGER_TRACE_TYPE="$trace_type"
}

##############################################################################
# Initialization
##############################################################################

# This function is called when the library is sourced
_common_lib_init() {
    # Detect PROJECT_ROOT if not set
    # Note: This is calculated from lib/common.sh location: op-challenger/scripts/lib -> op-challenger/scripts -> op-challenger -> project_root
    if [ -z "${PROJECT_ROOT:-}" ]; then
        # Get the directory where common.sh is located (op-challenger/scripts/lib/)
        local lib_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        # Go up three levels: lib -> scripts -> op-challenger -> project_root
        export PROJECT_ROOT="$(cd "${lib_dir}/../../.." && pwd)"
    fi
}

# Run initialization
_common_lib_init
