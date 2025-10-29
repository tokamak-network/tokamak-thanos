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

    # Patch op-program Dockerfile.repro to export cannon binary (if not already patched)
    local dockerfile="${PROJECT_ROOT}/op-program/Dockerfile.repro"
    if ! grep -q "COPY --from=builder /app/cannon/bin/cannon" "$dockerfile"; then
        log_info "Patching Dockerfile.repro to export cannon binary..."
        # Backup original
        cp "$dockerfile" "${dockerfile}.bak"
        # Add cannon to export-stage using awk
        awk '/FROM scratch AS export-stage/ {print; print "COPY --from=builder /app/cannon/bin/cannon ."; next} 1' \
            "${dockerfile}.bak" > "$dockerfile"
        log_success "  ✓ Dockerfile.repro patched"
    fi

    # Patch op-program Makefile to add --platform linux/amd64 for Apple Silicon compatibility
    local makefile="${PROJECT_ROOT}/op-program/Makefile"
    if ! grep -q "platform linux/amd64" "$makefile"; then
        log_info "Patching op-program Makefile to use linux/amd64 platform..."
        # Backup original
        cp "$makefile" "${makefile}.bak"
        # Add --platform linux/amd64 to docker build command
        sed -i.tmp 's/docker build --output/docker build --platform linux\/amd64 --output/g' "$makefile"
        rm -f "${makefile}.tmp"
        log_success "  ✓ op-program Makefile patched for linux/amd64"
    fi

    log_info "Building op-program and Cannon using Docker (reproducible build)..."
    log_info "   This ensures Linux binaries compatible with Docker containers"
    log_info "   (approximately 3-5 minutes)"
    echo ""

    # Use Docker reproducible build (generates Linux binaries)
    # Note: make reproducible-prestate runs in op-program directory
    # Docker export goes to op-program/bin/
    if make reproducible-prestate; then
        log_success "Cannon reproducible prestate build completed"

        # Check if cannon binary was exported to op-program/bin/
        local cannon_binary=""
        if [ -f "op-program/bin/cannon" ]; then
            cannon_binary="op-program/bin/cannon"
        else
            log_error "Cannon binary not found in op-program/bin/ after Docker build"
            log_error "Dockerfile patch may have failed"
            return 1
        fi

        # Copy cannon binary to expected location
        ensure_dir "${CANNON_BIN}"
        cp "$cannon_binary" "${CANNON_BIN}/cannon"
        chmod +x "${CANNON_BIN}/cannon"
        log_success "  ✓ Copied: cannon binary from Docker export"

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

            # Initialize submodules (required for op-program)
            log_info "Initializing Asterisc submodules..."
            cd "$asterisc_src_dir"
            if git submodule update --init --recursive --depth 1; then
                log_success "Submodules initialized"
            else
                log_error "Failed to initialize submodules"
                return 1
            fi
        else
            log_error "Failed to clone Asterisc repository"
            return 1
        fi
    else
        log_info "Asterisc repository already exists, pulling latest..."
        cd "$asterisc_src_dir"
        git pull || log_warn "Failed to pull latest changes (continuing with existing)"

        # Update submodules
        log_info "Updating submodules..."
        git submodule update --init --recursive --depth 1 || log_warn "Failed to update submodules (continuing)"
    fi

    # Build Asterisc using Docker (reproducible build)
    cd "$asterisc_src_dir"
    log_info "Building Asterisc using Docker (reproducible build)..."
    log_info "   This ensures Linux binaries compatible with Docker containers"
    log_info "   (approximately 5-10 minutes)"
    echo ""

    # Patch Dockerfile.repro to export asterisc binary (if not already patched)
    if ! grep -q "COPY --from=builder /app/rvgo/bin/asterisc" Dockerfile.repro; then
        log_info "Patching Dockerfile.repro to export asterisc binary..."
        # Backup original
        cp Dockerfile.repro Dockerfile.repro.bak
        # Add asterisc to export-stage using awk
        awk '/FROM scratch AS export-stage/ {print; print "COPY --from=builder /app/rvgo/bin/asterisc ."; next} 1' \
            Dockerfile.repro.bak > Dockerfile.repro
        log_success "  ✓ Dockerfile.repro patched"
    fi

    # Patch Makefile to add --platform linux/amd64 for reproducible builds
    if ! grep -q "platform linux/amd64" Makefile; then
        log_info "Patching Makefile to use linux/amd64 platform..."
        # Backup original
        cp Makefile Makefile.bak
        # Add --platform linux/amd64 to docker build command
        sed -i.tmp 's/docker build --output/docker build --platform linux\/amd64 --output/g' Makefile
        rm -f Makefile.tmp
        log_success "  ✓ Makefile patched for linux/amd64"
    fi

    # Use Docker reproducible build (generates Linux binaries + prestate)
    if make reproducible-prestate; then
        log_success "Asterisc reproducible prestate build completed"

        # Copy generated files to expected location
        ensure_dir "$ASTERISC_BIN"

        # Check for asterisc binary (should be exported now)
        local asterisc_binary=""
        if [ -f "bin/asterisc" ]; then
            asterisc_binary="bin/asterisc"
            log_info "  Found asterisc in: bin/asterisc (Docker export)"
        else
            log_error "Asterisc binary not found in bin/ after Docker build"
            log_error "Dockerfile patch may have failed"
            return 1
        fi

        cp "$asterisc_binary" "$ASTERISC_BIN/asterisc"
        chmod +x "$ASTERISC_BIN/asterisc"
        log_success "  ✓ Copied: asterisc binary from Docker export"

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

    # Always rebuild to ensure latest code changes are included
    # (same behavior as GameType 2 Asterisc build)

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

    # Build kona-client using Docker (reproducible build)
    log_info "Building kona-client using Docker (reproducible build)..."
    log_info "   This ensures Linux binaries compatible with Docker containers"
    log_info "   (approximately 5-10 minutes on first build)"
    echo ""

    cd "$kona_dir"

    # Create Dockerfile for kona-client if it doesn't exist
    local dockerfile_kona="${kona_dir}/Dockerfile.kona-client"
    if [ ! -f "$dockerfile_kona" ]; then
        log_info "Creating Dockerfile for kona-client reproducible build..."
        cat > "$dockerfile_kona" <<'EOF'
# Multi-stage Docker build for kona-client
FROM rust:1.82-bookworm AS builder

WORKDIR /workspace

# Copy the entire kona repository
COPY . .

# Build kona-client for Linux (release mode)
RUN cargo build --release -p kona-client

# Export stage - copy binaries to /out for easy extraction
FROM scratch AS export
COPY --from=builder /workspace/target/release/kona-client /kona-client
EOF
        log_success "Created: $dockerfile_kona"
    fi

    # Build using Docker (with platform specification for Apple Silicon compatibility)
    log_info "Running Docker build for kona-client..."
    if docker build --platform linux/amd64 -f "$dockerfile_kona" --target export --output type=local,dest=./bin-docker . ; then
        log_success "Docker build completed successfully"

        # Copy binary to expected location
        ensure_dir "${PROJECT_ROOT}/bin"

        if [ -f "./bin-docker/kona-client" ]; then
            cp "./bin-docker/kona-client" "${PROJECT_ROOT}/bin/kona-client"
            chmod +x "${PROJECT_ROOT}/bin/kona-client"
            log_success "  ✓ Copied: kona-client -> ${PROJECT_ROOT}/bin/kona-client"

            # Verify it's a Linux binary
            if command -v file &> /dev/null; then
                if file "${PROJECT_ROOT}/bin/kona-client" | grep -q "ELF"; then
                    log_success "  ✓ Verified: Linux ELF binary (compatible with Docker)"
                else
                    log_warn "  ⚠️ Binary format: $(file ${PROJECT_ROOT}/bin/kona-client)"
                fi
            fi

            # Cleanup Docker build artifacts
            rm -rf ./bin-docker

            cd "${PROJECT_ROOT}"
            return 0
        else
            log_error "kona-client binary not found in Docker build output"
            cd "${PROJECT_ROOT}"
            return 1
        fi
    else
        log_error "Docker build failed for kona-client"
        log_error ""
        log_error "Common issues:"
        log_error "  - Docker not running"
        log_error "  - Insufficient disk space"
        log_error "  - Network issues (cargo dependencies)"
        log_error ""
        log_error "You can provide kona-client binary manually to:"
        log_error "  ${PROJECT_ROOT}/bin/kona-client"

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

    # Always regenerate to ensure it matches the latest RISC-V binary
    # (same behavior as GameType 2 Asterisc build)

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

    log_info "Building kona op-program-client (RISC-V target) using Docker..."
    log_info "   This ensures reproducible RISC-V builds"
    log_info "   (approximately 10-15 minutes on first build)"
    echo ""

    cd "$kona_dir"

    # Use kona's official asterisc-builder Docker image
    # No need to create custom Dockerfile - use their proven build process
    log_info "Using kona's official asterisc-builder for RISC-V build..."
    log_info "   Image: ghcr.io/op-rs/kona/asterisc-builder:0.3.0"

    # Build RISC-V binary using kona's official asterisc-builder
    log_info "Running Docker build for RISC-V kona-client..."

    # Use kona's official build command from justfile
    if docker run --rm \
        -v "$(pwd):/workdir" \
        -w="/workdir" \
        ghcr.io/op-rs/kona/asterisc-builder:0.3.0 \
        cargo build -Zbuild-std=core,alloc -p kona-client --bin kona-client --profile release-client-lto; then

        log_success "Docker build completed successfully"

        # The binary is built in target/riscv64imac-unknown-none-elf/release-client-lto/kona-client
        # Note: kona's asterisc-builder Docker image uses riscv64imac target (not riscv64gc)
        local rv_binary="target/riscv64imac-unknown-none-elf/release-client-lto/kona-client"

        if [ ! -f "$rv_binary" ]; then
            log_error "RISC-V binary not found at expected location: $rv_binary"
            log_error "Checking other possible locations..."
            find target -name "kona-client" -type f 2>/dev/null || true
            cd "${PROJECT_ROOT}"
            return 1
        fi

        log_success "  ✓ Generated: RISC-V kona-client"
        log_info "    Location: $rv_binary"

        # Verify it's a RISC-V binary
        if command -v file &> /dev/null; then
            if file "$rv_binary" | grep -q "RISC-V"; then
                log_success "  ✓ Verified: RISC-V ELF binary"
            else
                log_warn "  ⚠️ Binary format: $(file $rv_binary)"
            fi
        fi
    else
        log_error "Failed to build kona-client using Docker"
        log_error "This is needed to generate the kona prestate"
        cd "${PROJECT_ROOT}"
        return 1
    fi

    # Generate prestate using asterisc (in Docker)
    log_info "Generating kona prestate from RISC-V ELF using Docker..."
    cd "${PROJECT_ROOT}"

    ensure_dir "${PROJECT_ROOT}/op-program/bin"

    # Use the RISC-V kona-client binary built above
    # Note: kona's asterisc-builder uses riscv64imac target
    local rv_binary_path="${kona_dir}/target/riscv64imac-unknown-none-elf/release-client-lto/kona-client"

    # Check if asterisc source directory exists (needed for Docker image)
    local asterisc_src_dir="${PROJECT_ROOT}/.asterisc-src"
    if [ ! -d "$asterisc_src_dir" ]; then
        log_error "Asterisc source directory not found: $asterisc_src_dir"
        log_error "This should have been created during Asterisc build"

        # Fallback to Asterisc prestate
        log_warn "Creating fallback: using Asterisc prestate for GameType 3"
        local asterisc_prestate="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
        if [ -f "$asterisc_prestate" ]; then
            cp "$asterisc_prestate" "$kona_prestate"
            local prestate_hash=$(cat "$kona_prestate" | jq -r '.stateHash' 2>/dev/null || echo "")
            if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                log_info "   Prestate Hash: $prestate_hash"
                export KONA_ABSOLUTE_PRESTATE="$prestate_hash"
            fi
            return 0
        else
            return 1
        fi
    fi

    # Run asterisc load-elf in Docker (same environment as asterisc build)
    # This ensures Linux binary compatibility
    log_info "   Running asterisc load-elf in Docker container..."

    if docker run --rm \
        --platform linux/amd64 \
        -v "${kona_dir}:/kona:ro" \
        -v "${PROJECT_ROOT}/op-program/bin:/output" \
        -v "${asterisc_src_dir}:/asterisc:ro" \
        -w /asterisc \
        golang:1.22-bookworm \
        bash -c "
            # Install required tools
            apt-get update -qq && apt-get install -y -qq jq > /dev/null 2>&1

            # Build asterisc from rvgo directory
            cd /asterisc/rvgo
            if [ -f Makefile ]; then
                # Use Makefile build (simple go build)
                make build 2>&1 | grep -v 'go: downloading' || true

                # asterisc binary is in rvgo/bin/asterisc
                if [ -f bin/asterisc ]; then
                    # Generate prestate
                    ./bin/asterisc load-elf \
                        --path /kona/target/riscv64imac-unknown-none-elf/release-client-lto/kona-client \
                        --out /output/prestate-kona.json
                else
                    echo 'ERROR: asterisc binary not found after build'
                    exit 1
                fi
            else
                echo 'ERROR: Makefile not found in rvgo directory'
                exit 1
            fi
        "; then

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
        log_error "Failed to generate kona prestate using Docker"
        log_error ""
        log_error "GameType 3 requires proper Kona prestate!"
        log_error "Cannot use Asterisc prestate as fallback - different binaries (op-program vs kona-client)"
        log_error ""
        log_error "Troubleshooting:"
        log_error "  1. Check Docker has Go 1.22+ (asterisc requires go >= 1.22.0)"
        log_error "  2. Verify asterisc source directory: ${asterisc_src_dir}"
        log_error "  3. Verify kona RISC-V binary: ${rv_binary_path}"
        log_error ""
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

    local kona_bin="${PROJECT_ROOT}/bin/kona-client"
    if [ -f "$kona_bin" ]; then
        log_info "  ✓ Kona (RISC-V + Rust):"
        log_info "    - Binary: $kona_bin"

        local kona_prestate="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
        if [ -f "$kona_prestate" ]; then
            local kona_hash=$(cat "$kona_prestate" | jq -r '.stateHash' 2>/dev/null || echo "")
            if [ -n "$kona_hash" ] && [ "$kona_hash" != "null" ]; then
                log_info "    - Prestate: $kona_hash ✅"
            else
                log_warn "    - Prestate: MISSING ❌"
            fi
        else
            log_warn "    - Prestate: NOT GENERATED ❌"
            log_warn "    - GameType 3 will fallback to Asterisc prestate"
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

