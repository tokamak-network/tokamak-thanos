#!/usr/bin/env bash

##############################################################################
# VM Build Module - Cannon (MIPS) and Asterisc (RISC-V)
#
# This module handles building all VM components needed for dispute games:
# - op-program (shared by both VMs)
# - Cannon VM (MIPS) + prestate
# - Asterisc VM (RISC-V) + prestate (auto-download or build)
##############################################################################

set -euo pipefail

# Source common library
# Note: Use local variable to avoid polluting parent script's SCRIPT_DIR
_MODULE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${_MODULE_DIR}/../lib/common.sh"

##############################################################################
# Configuration
##############################################################################

OP_PROGRAM_BIN="${PROJECT_ROOT}/op-program/bin"
CANNON_BIN="${PROJECT_ROOT}/cannon/bin"
ASTERISC_BIN="${PROJECT_ROOT}/asterisc/bin"
ASTERISC_REPO="https://github.com/ethereum-optimism/asterisc.git"

##############################################################################
# op-program Build
##############################################################################

build_op_program() {
    log_info "=========================================="
    log_info "Building op-program (shared by all VMs)"
    log_info "=========================================="
    echo ""

    cd "${PROJECT_ROOT}"

    log_info "Building op-program... (approximately 2-5 minutes)"
    if make op-program; then
        log_success "op-program build completed"

        # Verify binaries exist
        if [ -f "$OP_PROGRAM_BIN/op-program" ]; then
            log_success "  ✓ op-program binary: $OP_PROGRAM_BIN/op-program"
        else
            log_error "op-program binary not found after build!"
            return 1
        fi

        return 0
    else
        log_error "op-program build failed"
        return 1
    fi
}

##############################################################################
# Cannon VM Build (MIPS)
##############################################################################

build_cannon() {
    log_info "=========================================="
    log_info "Building Cannon VM (MIPS)"
    log_info "=========================================="
    echo ""

    cd "${PROJECT_ROOT}"

    log_info "Building op-program and Cannon using Docker (reproducible build)..."
    log_info "   This ensures Linux binaries compatible with Docker containers"
    log_info "   (approximately 3-5 minutes)"
    echo ""

    # Use Docker reproducible build (generates Linux binaries)
    if make reproducible-prestate; then
        log_success "Cannon reproducible prestate build completed"

        # Verify generated files
        if [ -f "${OP_PROGRAM_BIN}/prestate-proof.json" ] && \
           [ -f "${CANNON_BIN}/cannon" ]; then
            log_success "  ✓ Generated: cannon binary (Linux)"
            log_success "  ✓ Generated: prestate-proof.json"

            # Verify it's a Linux binary
            if command -v file &> /dev/null; then
                file "${CANNON_BIN}/cannon" | grep -q "ELF.*x86-64" && \
                    log_success "  ✓ Verified: Linux ELF binary" || \
                    log_warn "  ⚠️ Binary format: $(file ${CANNON_BIN}/cannon)"
            fi
        else
            log_error "Expected files not found after reproducible build"
            return 1
        fi

        # Get generated prestate hash
        local prestate_proof_file="${OP_PROGRAM_BIN}/prestate-proof.json"
        if [ -f "$prestate_proof_file" ]; then
            local prestate_hash=$(cat "$prestate_proof_file" | jq -r '.pre' 2>/dev/null || echo "")
            if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                log_success "✅ Cannon Absolute Prestate Hash:"
                log_info "   $prestate_hash"
                log_config "Cannon Absolute Prestate: $prestate_hash"
                export CANNON_ABSOLUTE_PRESTATE="$prestate_hash"
            fi

            # Create prestate.json symlink for Docker compatibility
            local prestate_file="${OP_PROGRAM_BIN}/prestate.json"
            if [ ! -f "$prestate_file" ] && [ ! -L "$prestate_file" ]; then
                # Create symlink in the same directory (relative path)
                cd "${OP_PROGRAM_BIN}"
                ln -sf prestate-proof.json prestate.json
                cd "${PROJECT_ROOT}"
                log_info "Created symlink: ${OP_PROGRAM_BIN}/prestate.json -> prestate-proof.json"
            fi
        fi

        return 0
    else
        log_error "Cannon prestate generation failed"
        return 1
    fi
}

##############################################################################
# Asterisc VM Build/Download (RISC-V)
##############################################################################

download_asterisc_binary() {
    log_info "Attempting to download pre-built Asterisc binary..."

    # TODO: Add actual download URL when available
    # For now, this is a placeholder
    log_warn "Pre-built Asterisc binaries are not available yet."
    return 1
}

build_asterisc_from_source() {
    log_info "Building Asterisc VM from source..."

    local asterisc_src_dir="${PROJECT_ROOT}/.asterisc-src"

    # Clone if not exists
    if [ ! -d "$asterisc_src_dir" ]; then
        log_info "Cloning Asterisc repository..."
        if git clone --depth 1 "$ASTERISC_REPO" "$asterisc_src_dir"; then
            log_success "Asterisc repository cloned"
        else
            log_error "Failed to clone Asterisc repository"
            return 1
        fi
    else
        log_info "Asterisc repository already exists, pulling latest..."
        cd "$asterisc_src_dir"
        git pull || log_warn "Failed to pull latest changes (continuing with existing)"
    fi

    # Build Asterisc using Docker (reproducible build)
    cd "$asterisc_src_dir"
    log_info "Building Asterisc using Docker (reproducible build)..."
    log_info "   This ensures Linux binaries compatible with Docker containers"
    log_info "   (approximately 5-10 minutes)"
    echo ""

    # Use Docker reproducible build (generates Linux binaries + prestate)
    if make reproducible-prestate; then
        log_success "Asterisc reproducible prestate build completed"

        # Copy generated files to expected location
        ensure_dir "$ASTERISC_BIN"

        # IMPORTANT: Only use Docker reproducible build output (bin/)
        # Docker build should export asterisc binary to bin/
        if [ ! -f "bin/asterisc" ]; then
            log_error "Asterisc binary not found in bin/ after Docker build"
            log_error "Docker reproducible build must generate asterisc binary"
            log_error "Check Dockerfile.repro export-stage section"
            return 1
        fi

        cp "bin/asterisc" "$ASTERISC_BIN/asterisc"
        chmod +x "$ASTERISC_BIN/asterisc"
        log_success "  ✓ Copied: asterisc binary from Docker reproducible build"

        # Verify it's a Linux ELF binary (required for Docker containers)
        if command -v file &> /dev/null; then
            if file "$ASTERISC_BIN/asterisc" | grep -q "ELF"; then
                log_success "  ✓ Verified: Linux ELF binary (compatible with Docker)"
            else
                log_error "  ❌ Binary is not ELF format: $(file $ASTERISC_BIN/asterisc)"
                log_error "  Docker reproducible build should generate Linux ELF binaries"
                return 1
            fi
        fi

        # Copy prestate files from Docker build output (bin/)
        # IMPORTANT: Only use Docker reproducible build results, no local fallback!
        if [ ! -f "bin/prestate.json" ]; then
            log_error "Asterisc prestate.json not found in bin/ after Docker build"
            log_error "Docker reproducible build must generate prestate.json"
            return 1
        fi

        cp "bin/prestate.json" "$ASTERISC_BIN/prestate-proof.json"
        log_success "  ✓ Copied: prestate.json from Docker reproducible build"

        # Extract prestate hash from Docker build prestate.json
        # Note: Asterisc uses .stateHash field (not .pre like Cannon)
        local prestate_hash=$(cat "bin/prestate.json" | jq -r '.stateHash' 2>/dev/null || echo "")
        if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
            log_success "✅ Asterisc Absolute Prestate Hash (Docker reproducible build):"
            log_info "   $prestate_hash"
            log_config "Asterisc Absolute Prestate: $prestate_hash"
            export ASTERISC_ABSOLUTE_PRESTATE="$prestate_hash"
        else
            log_error "Failed to extract stateHash from Docker build prestate.json"
            return 1
        fi

        # Create symlink for backward compatibility
        if [ -f "$ASTERISC_BIN/prestate-proof.json" ]; then
            cd "$ASTERISC_BIN"
            ln -sf prestate-proof.json prestate.json
            cd "$asterisc_src_dir"
            log_success "  ✓ Created symlink: prestate.json -> prestate-proof.json"
        fi

        # Copy meta.json if it exists (optional, from Docker build)
        if [ -f "bin/meta.json" ]; then
            cp "bin/meta.json" "$ASTERISC_BIN/meta.json"
            log_success "  ✓ Copied: meta.json from Docker build"
        fi

        return 0
    else
        log_error "Asterisc Docker reproducible build failed"
        log_error "Make sure Docker is running and Dockerfile.repro exists in asterisc repo"
        return 1
    fi
}

setup_asterisc_fallback() {
    log_warn "⚠️  Setting up Asterisc fallback..."
    log_warn "   This is NOT recommended for production!"
    log_warn "   MIPS and RISC-V are different architectures."
    echo ""

    ensure_dir "$ASTERISC_BIN"

    # Create symlink to op-program as fallback
    if [ ! -f "$ASTERISC_BIN/asterisc" ] && [ -f "$OP_PROGRAM_BIN/op-program" ]; then
        ln -sf "$OP_PROGRAM_BIN/op-program" "$ASTERISC_BIN/asterisc"
        log_info "Created symlink: asterisc -> op-program"
    fi

    # ❌ DO NOT copy MIPS prestate for RISC-V!
    # MIPS (Cannon) and RISC-V (Asterisc) use different instruction sets
    # Using wrong prestate will cause verification failures

    log_warn "⚠️  No RISC-V prestate available with fallback method"
    log_warn "   Asterisc games will NOT work correctly"
    log_warn "   Please build Asterisc from source for proper prestate"
    echo ""

    log_success "Asterisc fallback setup completed (with limitations)"
    return 0
}

build_asterisc() {
    log_info "=========================================="
    log_info "Building Asterisc VM (RISC-V)"
    log_info "=========================================="
    echo ""

    # Try methods in order of preference (same as Cannon - no skip logic):
    # 1. Download pre-built binary (fastest) - TODO
    # 2. Build from source (recommended, generates proper RISC-V prestate)
    # 3. Fallback to op-program symlink (NOT RECOMMENDED - no proper prestate)

    if download_asterisc_binary; then
        log_success "Asterisc binary downloaded successfully"
        return 0
    fi

    log_warn "Download not available, building from source..."

    if build_asterisc_from_source; then
        log_success "✅ Asterisc built from source successfully"
        log_success "✅ RISC-V prestate generated correctly"
        return 0
    fi

    log_error "❌ Build from source failed!"
    log_warn ""
    log_warn "⚠️  CRITICAL: Fallback method does NOT generate proper RISC-V prestate"
    log_warn "   MIPS and RISC-V are different architectures!"
    log_warn "   Asterisc games will NOT work without proper prestate."
    log_warn ""
    log_warn "Options:"
    log_warn "  1. Fix build errors and retry"
    log_warn "  2. Use only GameType 0/1 (Cannon) instead of GameType 2 (Asterisc)"
    log_warn ""

    read -p "Continue with fallback anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_error "Asterisc setup cancelled"
        return 1
    fi

    if setup_asterisc_fallback; then
        log_warn "Asterisc fallback setup completed (NOT suitable for production)"
        return 0
    else
        log_error "All Asterisc setup methods failed"
        return 1
    fi
}

##############################################################################
# Kona Client Build (Rust-based for GameType 3)
##############################################################################

build_kona_client() {
    log_info "=========================================="
    log_info "Building kona-client (Rust) for GameType 3"
    log_info "=========================================="
    echo ""

    # Check if kona-client already exists
    local kona_bin="${PROJECT_ROOT}/bin/kona-client"
    if [ -f "$kona_bin" ]; then
        log_success "kona-client binary already exists: $kona_bin"
        log_info "To rebuild, delete the existing binary first"
        return 0
    fi

    # Check for Rust/Cargo
    if ! command -v cargo &> /dev/null; then
        log_error "Rust/Cargo not installed!"
        log_error ""
        log_error "Please install Rust toolchain:"
        log_error "  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
        log_error ""
        log_error "Or skip kona build and provide kona-client binary manually to:"
        log_error "  ${PROJECT_ROOT}/bin/kona-client"
        return 1
    fi

    # Check for kona repository
    local kona_dir="${KONA_DIR:-${PROJECT_ROOT}/../kona}"

    if [ ! -d "$kona_dir" ]; then
        log_warn "kona repository not found at: $kona_dir"
        log_info "Attempting to clone kona repository..."

        local parent_dir="$(dirname ${PROJECT_ROOT})"
        cd "$parent_dir"

        if git clone --depth 1 https://github.com/op-rs/kona.git; then
            log_success "kona repository cloned successfully"
            kona_dir="${parent_dir}/kona"
        else
            log_error "Failed to clone kona repository"
            log_error ""
            log_error "Please clone manually:"
            log_error "  cd $(dirname ${PROJECT_ROOT})"
            log_error "  git clone https://github.com/op-rs/kona.git"
            return 1
        fi

        cd "${PROJECT_ROOT}"
    fi

    # Build kona-client
    log_info "Building kona-client from source..."
    log_info "Repository: $kona_dir"
    log_info "(this may take 5-10 minutes on first build)"
    echo ""

    cd "$kona_dir"

    if cargo build --release -p kona-client; then
        log_success "kona-client build completed"

        # Copy binary to expected location
        ensure_dir "${PROJECT_ROOT}/bin"
        cp target/release/kona-client "${PROJECT_ROOT}/bin/kona-client"
        chmod +x "${PROJECT_ROOT}/bin/kona-client"

        log_success "  ✓ Copied: kona-client -> ${PROJECT_ROOT}/bin/kona-client"

        cd "${PROJECT_ROOT}"
        return 0
    else
        log_error "kona-client build failed"
        log_error ""
        log_error "Common issues:"
        log_error "  - Outdated Rust version (update with: rustup update)"
        log_error "  - Missing dependencies (check kona README)"
        log_error ""
        log_error "You can still use GameType 3 if you provide kona-client binary manually"

        cd "${PROJECT_ROOT}"
        return 1
    fi
}

# Generate kona prestate
generate_kona_prestate() {
    log_info "=========================================="
    log_info "Generating Kona Prestate for GameType 3"
    log_info "=========================================="
    echo ""

    local kona_prestate="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"

    # Check if prestate already exists
    if [ -f "$kona_prestate" ]; then
        log_success "Kona prestate already exists: $kona_prestate"
        log_info "To regenerate, delete the existing file first"
        return 0
    fi

    # Check for asterisc binary (needed to generate prestate from ELF)
    if [ ! -f "${ASTERISC_BIN}/asterisc" ]; then
        log_error "Asterisc binary not found!"
        log_error "Please build Asterisc VM first (it's needed to generate kona prestate)"
        return 1
    fi

    # Check for kona repository
    local kona_dir="${KONA_DIR:-${PROJECT_ROOT}/../kona}"
    if [ ! -d "$kona_dir" ]; then
        log_error "Kona repository not found at: $kona_dir"
        log_error "Please run build_kona_client first to clone the repository"
        return 1
    fi

    log_info "Building kona op-program-client (RISC-V target)..."
    log_info "(this may take 10-15 minutes on first build)"
    echo ""

    cd "$kona_dir"

    # Add RISC-V target if not installed
    if ! rustup target list | grep -q "riscv64gc-unknown-linux-gnu (installed)"; then
        log_info "Adding RISC-V target to rustup..."
        rustup target add riscv64gc-unknown-linux-gnu
    fi

    # Build RISC-V op-program-client
    if cargo build --release --target riscv64gc-unknown-linux-gnu -p op-program-client 2>/dev/null || \
       make build-rv64 2>/dev/null; then
        log_success "kona op-program-client (RISC-V) build completed"
    else
        log_error "Failed to build kona op-program-client"
        log_error "This is needed to generate the kona prestate"
        cd "${PROJECT_ROOT}"
        return 1
    fi

    # Generate prestate using asterisc
    local rv_binary="target/riscv64gc-unknown-linux-gnu/release/op-program-client"
    if [ ! -f "$rv_binary" ]; then
        log_error "RISC-V binary not found: $rv_binary"
        cd "${PROJECT_ROOT}"
        return 1
    fi

    log_info "Generating kona prestate from RISC-V ELF..."
    cd "${PROJECT_ROOT}"

    if "${ASTERISC_BIN}/asterisc" load-elf \
        --path "${kona_dir}/${rv_binary}" \
        --out "$kona_prestate"; then
        log_success "✅ Kona prestate generated: $kona_prestate"

        # Extract and display prestate hash
        local prestate_hash=$(cat "$kona_prestate" | jq -r '.stateHash' 2>/dev/null || echo "")
        if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
            log_success "✅ Kona Absolute Prestate Hash:"
            log_info "   $prestate_hash"
            log_config "Kona Absolute Prestate: $prestate_hash"
            export KONA_ABSOLUTE_PRESTATE="$prestate_hash"
        fi

        return 0
    else
        log_error "Failed to generate kona prestate"
        return 1
    fi
}

##############################################################################
# Main Build Function
##############################################################################

build_vms() {
    local build_cannon="${1:-true}"
    local build_asterisc="${2:-false}"

    log_info "=========================================="
    log_info "VM Build Module"
    log_info "=========================================="
    log_info "Build Cannon: $build_cannon"
    log_info "Build Asterisc: $build_asterisc"
    echo ""

    # Always build op-program first (shared by all VMs)
    if ! build_op_program; then
        log_error "Failed to build op-program"
        return 1
    fi

    # Build Cannon if requested
    if [ "$build_cannon" = "true" ]; then
        if ! build_cannon; then
            log_error "Failed to build Cannon VM"
            return 1
        fi
    else
        log_info "Skipping Cannon VM build"
    fi

    # Build Asterisc if requested
    if [ "$build_asterisc" = "true" ]; then
        if ! build_asterisc; then
            log_error "Failed to build Asterisc VM"
            return 1
        fi
    else
        log_info "Skipping Asterisc VM build"
    fi

    log_success "=========================================="
    log_success "VM Build Module Completed"
    log_success "=========================================="
    echo ""

    # Summary
    log_info "Built VM Summary:"
    echo ""

    if [ -f "$OP_PROGRAM_BIN/op-program" ]; then
        log_info "  ✓ op-program: $OP_PROGRAM_BIN/op-program"
    fi

    if [ -f "$CANNON_BIN/cannon" ]; then
        log_info "  ✓ Cannon (MIPS):"
        log_info "    - Binary: $CANNON_BIN/cannon"
        if [ -f "$OP_PROGRAM_BIN/prestate-proof.json" ]; then
            local cannon_hash=$(cat "$OP_PROGRAM_BIN/prestate-proof.json" | jq -r '.pre' 2>/dev/null || echo "")
            [ -n "$cannon_hash" ] && log_info "    - Prestate: $cannon_hash"
        fi
    fi

    if [ -f "$ASTERISC_BIN/asterisc" ]; then
        log_info "  ✓ Asterisc (RISC-V):"
        log_info "    - Binary: $ASTERISC_BIN/asterisc"

        if [ -f "$ASTERISC_BIN/prestate-proof.json" ]; then
            # Read stateHash directly from prestate-proof.json (Docker build output)
            # Note: Asterisc uses .stateHash field (not .pre like Cannon)
            local asterisc_hash=$(cat "$ASTERISC_BIN/prestate-proof.json" | jq -r '.stateHash' 2>/dev/null || echo "")
            if [ -n "$asterisc_hash" ] && [ "$asterisc_hash" != "null" ]; then
                log_info "    - Prestate: $asterisc_hash ✅"
            else
                log_warn "    - Prestate: MISSING ❌"
            fi
        else
            log_warn "    - Prestate: NOT GENERATED ❌"
            log_warn "    - GameType 2 will NOT work correctly!"
        fi
    fi

    echo ""

    return 0
}

##############################################################################
# CLI Interface (if run directly)
##############################################################################

if [ "${BASH_SOURCE[0]}" -ef "$0" ]; then
    # Script is being run directly, not sourced

    BUILD_CANNON="true"
    BUILD_ASTERISC="false"

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --cannon-only)
                BUILD_ASTERISC="false"
                shift
                ;;
            --asterisc-only)
                BUILD_CANNON="false"
                BUILD_ASTERISC="true"
                shift
                ;;
            --all)
                BUILD_CANNON="true"
                BUILD_ASTERISC="true"
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --cannon-only     Build only Cannon VM (default)"
                echo "  --asterisc-only   Build only Asterisc VM"
                echo "  --all             Build both Cannon and Asterisc"
                echo "  --help            Show this help"
                echo ""
                echo "Examples:"
                echo "  $0                    # Build Cannon only"
                echo "  $0 --all              # Build both VMs"
                echo "  $0 --asterisc-only    # Build Asterisc only"
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

    # Run build
    if build_vms "$BUILD_CANNON" "$BUILD_ASTERISC"; then
        log_success "All VM builds completed successfully!"
        exit 0
    else
        log_error "VM build failed"
        exit 1
    fi
fi

