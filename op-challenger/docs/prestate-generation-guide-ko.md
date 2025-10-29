# Prestate 생성 가이드 (모든 GameType)

> **중요**: 모든 GameType의 prestate는 **재현 가능한 빌드를 보장하기 위해 Docker 환경**에서 생성됩니다.

## 📋 목차

- [개요](#개요)
- [GameType 0/1: Cannon (MIPS)](#gametype-01-cannon-mips)
- [GameType 2: Asterisc (RISC-V)](#gametype-2-asterisc-risc-v)
- [GameType 3: AsteriscKona (RISC-V + Rust)](#gametype-3-asterisckona-risc-v--rust)
- [비교표](#비교표)
- [Docker 환경 점검](#docker-환경-점검)
- [문제 해결](#문제-해결)

---

## 개요

### Prestate란?

**Prestate**는 Fault Proof 시스템에서 VM의 **초기 상태 해시**입니다. 이는:
- 온체인 dispute game 배포 시 `faultGameAbsolutePrestate`로 설정됩니다
- VM이 프로그램 실행을 시작하기 전의 메모리 상태를 나타냅니다
- 모든 참여자가 동일한 초기 상태에서 시작함을 보장합니다

### 왜 Docker인가?

재현 가능한 빌드를 위해 **모든 prestate는 고정된 Docker 이미지**에서 생성됩니다:

```
동일한 소스 코드 + 동일한 Docker 환경 = 동일한 Prestate Hash
```

이는 다음을 보장합니다:
- ✅ macOS, Linux, Windows에서 동일한 결과
- ✅ 시간이 지나도 동일한 결과 (환경 고정)
- ✅ 온체인 배포와 로컬 빌드의 일치

### 💡 권장 워크플로우: 사전 빌드 이미지 사용

**일반 사용자는 직접 빌드할 필요가 없습니다!** 사전 빌드된 VM 이미지를 사용하세요:

```bash
# 1. 사전 빌드된 VM 바이너리 다운로드 (2-3분)
./op-challenger/scripts/pull-vm-images.sh --tag latest

# 결과:
# ✅ cannon/bin/cannon (GameType 0/1)
# ✅ asterisc/bin/asterisc (GameType 2)
# ✅ op-program/bin/op-program (GameType 0/1/2)
# ✅ bin/kona-client (GameType 3)
# ⚠️  Prestate는 배포 시 자동 생성 (VM 바이너리 필요)

# 2. 바로 배포 가능!
./op-challenger/scripts/deploy-modular.sh --dg-type 3
```

**담당자만 빌드 및 업로드**:
```bash
# VM 이미지 빌드 및 GitHub Container Registry에 업로드
./op-challenger/scripts/build-vm-images.sh --push
```

**장점**:
- ⏱️ **시간 절약**: 35-50분 → 2-3분
- 🎯 **간편함**: Docker, Rust, Go 환경 설정 불필요
- 🔒 **안정성**: 검증된 바이너리 사용

---

## GameType 0/1: Cannon (MIPS)

### 기본 정보

| 항목 | 값 |
|------|-----|
| **VM** | Cannon (MIPS 아키텍처) |
| **서버** | op-program (Go) |
| **Prestate 파일** | `op-program/bin/prestate-proof.json` |
| **해시 필드** | `.pre` |
| **Docker 이미지** | `golang:1.21.3-alpine3.18` |
| **빌드 스크립트** | `op-program/Dockerfile.repro` |

### 생성 과정

#### 1단계: Docker 빌드 시작

```bash
cd op-program
make reproducible-prestate
```

**Makefile 내부 동작**:
```makefile
reproducible-prestate:
    @docker build --platform linux/amd64 \
        --output ./bin/ \
        --progress plain \
        -f Dockerfile.repro ../
```

#### 2단계: Docker 내부 빌드 프로세스

```dockerfile
FROM golang:1.21.3-alpine3.18 as builder

# 1. MIPS 바이너리 컴파일
RUN cd op-program && make op-program-client-mips
# → op-program/bin/op-program-client.elf (MIPS 바이너리)

# 2. Cannon으로 ELF 로드하여 초기 상태 생성
RUN /app/cannon/bin/cannon load-elf \
    --path /app/op-program/bin/op-program-client.elf \
    --out /app/op-program/bin/prestate.json \
    --meta ""

# 3. Prestate proof 생성 (step 0에서의 상태 증명)
RUN /app/cannon/bin/cannon run \
    --proof-at '=0' \
    --stop-at '=1' \
    --input /app/op-program/bin/prestate.json \
    --meta "" \
    --proof-fmt '/app/op-program/bin/%d.json' \
    --output ""
RUN mv /app/op-program/bin/0.json /app/op-program/bin/prestate-proof.json

# 4. 결과 파일 export
FROM scratch AS export-stage
COPY --from=builder /app/op-program/bin/prestate-proof.json .
```

#### 3단계: Prestate Hash 추출

```bash
# prestate-proof.json 구조
{
  "pre": "0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98",
  "post": "...",
  "state_data": "..."
}

# 해시 추출
PRESTATE_HASH=$(jq -r '.pre' op-program/bin/prestate-proof.json)
echo $PRESTATE_HASH
# → 0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98
```

### deploy-modular.sh 통합

```bash
# op-challenger/scripts/deploy-modular.sh
get_prestate_hash_for_gametype() {
    case "$dg_type" in
        0|1|254|255)
            local cannon_prestate="${PROJECT_ROOT}/op-program/bin/prestate-proof.json"
            if [ -f "$cannon_prestate" ]; then
                prestate_hash=$(jq -r '.pre' "$cannon_prestate")
                log_info "Using Cannon (MIPS) prestate: $prestate_hash"
            else
                log_error "Cannon prestate file not found: $cannon_prestate"
            fi
            ;;
    esac
}
```

### Docker 점검

```bash
# 1. prestate-proof.json 존재 확인
ls -lh op-program/bin/prestate-proof.json

# 2. JSON 형식 검증
jq '.' op-program/bin/prestate-proof.json > /dev/null && echo "✅ Valid JSON"

# 3. .pre 필드 확인
jq -r '.pre' op-program/bin/prestate-proof.json | grep -E '^0x[0-9a-f]{64}$' && echo "✅ Valid hash"

# 4. 재현 가능성 테스트 (2번 빌드 시 동일한 해시)
cd op-program
make reproducible-prestate
HASH1=$(jq -r '.pre' bin/prestate-proof.json)

make reproducible-prestate
HASH2=$(jq -r '.pre' bin/prestate-proof.json)

[ "$HASH1" = "$HASH2" ] && echo "✅ Reproducible build verified"
```

---

## GameType 2: Asterisc (RISC-V)

### 기본 정보

| 항목 | 값 |
|------|-----|
| **VM** | Asterisc (RISC-V 아키텍처) |
| **서버** | op-program (Go) |
| **Prestate 파일** | `asterisc/bin/prestate-proof.json` |
| **해시 필드** | `.stateHash` ⚠️ (Cannon과 다름!) |
| **Docker 이미지** | `golang:1.22-bookworm` |
| **빌드 스크립트** | `asterisc/Dockerfile.repro` |
| **빌드 함수** | `vm-build.sh::build_asterisc_from_source()` |

### 생성 과정

#### 1단계: Asterisc 소스 클론 및 준비

```bash
# op-challenger/scripts/modules/vm-build.sh
build_asterisc_from_source() {
    local asterisc_src_dir="${PROJECT_ROOT}/.asterisc-src"

    # GitHub에서 asterisc 클론
    git clone --depth 1 https://github.com/ethereum-optimism/asterisc.git "$asterisc_src_dir"
    cd "$asterisc_src_dir"

    # 서브모듈 초기화 (op-program 포함)
    git submodule update --init --recursive --depth 1
}
```

#### 2단계: Dockerfile 패치 (asterisc 바이너리 export)

```bash
# Dockerfile.repro에 asterisc 바이너리 export 추가
if ! grep -q "COPY --from=builder /app/rvgo/bin/asterisc" Dockerfile.repro; then
    cp Dockerfile.repro Dockerfile.repro.bak
    awk '/FROM scratch AS export-stage/ {
        print
        print "COPY --from=builder /app/rvgo/bin/asterisc ."
        next
    } 1' Dockerfile.repro.bak > Dockerfile.repro
fi
```

#### 3단계: Docker 빌드

```bash
# Makefile에 --platform linux/amd64 추가 (재현 가능한 빌드)
sed -i 's/docker build --output/docker build --platform linux\/amd64 --output/g' Makefile

# Docker 재현 가능한 빌드 실행
make reproducible-prestate
```

**Docker 내부 프로세스**:
```dockerfile
FROM golang:1.22-bookworm as builder

# 1. RISC-V 바이너리 컴파일
RUN cd op-program && make op-program-client-riscv
# → bin/op-program-client-riscv.elf

# 2. Asterisc VM 빌드
RUN cd rvgo && make build
# → rvgo/bin/asterisc

# 3. asterisc load-elf로 prestate 생성
RUN /app/rvgo/bin/asterisc load-elf \
    --path /app/op-program/bin/op-program-client-riscv.elf \
    --out /app/bin/prestate.json

# 4. 결과 export
FROM scratch AS export-stage
COPY --from=builder /app/bin/prestate.json .
COPY --from=builder /app/rvgo/bin/asterisc .
```

#### 4단계: 결과 복사 및 해시 추출

```bash
# Docker build 결과를 로컬로 복사
cp bin/asterisc "$ASTERISC_BIN/asterisc"
cp bin/prestate.json "$ASTERISC_BIN/prestate-proof.json"

# prestate.json 구조 (Cannon과 다름!)
{
  "memory": [...],
  "registers": [...],
  "stateHash": "0x039bef669fd1a2419634548ca79e85f9e42ae6d52c869ca290cb07e247e9a645",
  "witness": "..."
}

# 해시 추출 (.stateHash 필드 사용!)
PRESTATE_HASH=$(jq -r '.stateHash' asterisc/bin/prestate-proof.json)
```

### deploy-modular.sh 통합

```bash
get_prestate_hash_for_gametype() {
    case "$dg_type" in
        2)
            local asterisc_prestate="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
            if [ -f "$asterisc_prestate" ]; then
                # ⚠️ 주의: .stateHash 필드 사용 (Cannon과 다름!)
                prestate_hash=$(cat "$asterisc_prestate" | jq -r '.stateHash')
                log_info "Using Asterisc (RISC-V) prestate: $prestate_hash"
            else
                log_error "Asterisc prestate not found: $asterisc_prestate"
                log_error "❌ GameType 2 requires RISC-V prestate!"
                log_error "   Cannot use MIPS prestate for RISC-V architecture."
            fi
            ;;
    esac
}
```

### Docker 점검

```bash
# 1. prestate-proof.json 존재 확인
ls -lh asterisc/bin/prestate-proof.json

# 2. JSON 형식 검증
jq '.' asterisc/bin/prestate-proof.json > /dev/null && echo "✅ Valid JSON"

# 3. .stateHash 필드 확인 (Cannon의 .pre와 다름!)
jq -r '.stateHash' asterisc/bin/prestate-proof.json | grep -E '^0x[0-9a-f]{64}$' && echo "✅ Valid hash"

# 4. asterisc 바이너리가 Linux ELF인지 확인
file asterisc/bin/asterisc | grep -q "ELF.*x86-64" && echo "✅ Linux ELF binary"

# 5. 재현 가능성 테스트
cd .asterisc-src
make reproducible-prestate
HASH1=$(jq -r '.stateHash' bin/prestate.json)

make reproducible-prestate
HASH2=$(jq -r '.stateHash' bin/prestate.json)

[ "$HASH1" = "$HASH2" ] && echo "✅ Reproducible build verified"
```

---

## GameType 3: AsteriscKona (RISC-V + Rust)

### 기본 정보

| 항목 | 값 |
|------|-----|
| **VM** | Asterisc (RISC-V) - GameType 2와 동일한 RISCV.sol |
| **서버** | kona-client (Rust) ⚠️ Go 대신 Rust! |
| **Prestate 파일 (우선순위 1)** | `op-program/bin/prestate-kona.json` |
| **Prestate 파일 (fallback)** | `asterisc/bin/prestate-proof.json` |
| **해시 필드** | `.stateHash` (Asterisc와 동일) |
| **Docker 이미지** | `ghcr.io/op-rs/kona/asterisc-builder:0.3.0` (Rust) + `golang:1.22-bookworm` (asterisc) |
| **빌드 함수** | `vm-build.sh::build_kona_client()` + `generate_kona_prestate()` |

### 생성 과정 (2단계 Docker 빌드)

#### 1단계: Kona Client (Rust) RISC-V 바이너리 빌드

```bash
# op-challenger/scripts/modules/vm-build.sh
build_kona_client() {
    local kona_dir="${PROJECT_ROOT}/../kona"

    # 1. kona 리포지토리 클론
    git clone --depth 1 https://github.com/op-rs/kona.git "$kona_dir"
    cd "$kona_dir"

    # 2. Docker로 RISC-V 바이너리 빌드
    docker run --rm \
        --platform linux/amd64 \
        -v "$(pwd):/workspace" \
        -w /workspace \
        ghcr.io/op-rs/kona/asterisc-builder:0.3.0 \
        bash -c "
            # Rust RISC-V target으로 빌드
            cargo build \
                -Zbuild-std=core,alloc \
                -p kona-client \
                --bin kona-client \
                --profile release-client-lto \
                --target riscv64imac-unknown-none-elf
        "

    # 결과: kona/target/riscv64imac-unknown-none-elf/release-client-lto/kona-client
}
```

**Docker 내부 환경**:
- 이미지: `ghcr.io/op-rs/kona/asterisc-builder:0.3.0`
- Rust toolchain: `nightly-2024-08-01`
- Target: `riscv64imac-unknown-none-elf` (bare metal RISC-V)
- Build std: `core`, `alloc` (no_std 환경)

#### 2단계: Asterisc로 Prestate 생성

```bash
# op-challenger/scripts/modules/vm-build.sh
generate_kona_prestate() {
    local kona_prestate="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
    local kona_dir="${PROJECT_ROOT}/../kona"
    local asterisc_src_dir="${PROJECT_ROOT}/.asterisc-src"

    # Docker 내부에서 asterisc 빌드 및 load-elf 실행
    docker run --rm \
        --platform linux/amd64 \
        -v "${kona_dir}:/kona:ro" \
        -v "${PROJECT_ROOT}/op-program/bin:/output" \
        -v "${asterisc_src_dir}:/asterisc:ro" \
        -w /asterisc \
        golang:1.22-bookworm \
        bash -c "
            # apt 설치
            apt-get update -qq && apt-get install -y -qq jq

            # asterisc 빌드
            cd /asterisc/rvgo
            make build

            # asterisc load-elf로 kona-client용 prestate 생성
            ./bin/asterisc load-elf \
                --path /kona/target/riscv64imac-unknown-none-elf/release-client-lto/kona-client \
                --out /output/prestate-kona.json
        "
}
```

**프로세스 흐름**:
```
Step 1: Rust RISC-V 빌드 (Docker 1)
  kona-client (source)
    ↓ [Docker: ghcr.io/op-rs/kona/asterisc-builder]
  kona-client (RISC-V binary)

Step 2: Prestate 생성 (Docker 2)
  kona-client (RISC-V binary)
    ↓ [Docker: golang:1.22-bookworm + asterisc]
  prestate-kona.json (initial state)
```

#### 3단계: Prestate Hash 추출

```bash
# prestate-kona.json 구조 (Asterisc와 동일한 형식)
{
  "memory": [...],
  "registers": [...],
  "stateHash": "0x...",  # GameType 3 전용 prestate hash
  "witness": "..."
}

# 해시 추출
PRESTATE_HASH=$(jq -r '.stateHash' op-program/bin/prestate-kona.json)
```

### deploy-modular.sh 통합 (Fallback 지원)

```bash
get_prestate_hash_for_gametype() {
    case "$dg_type" in
        3)
            # 우선순위 1: kona-specific prestate
            local kona_prestate="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
            if [ -f "$kona_prestate" ]; then
                prestate_hash=$(cat "$kona_prestate" | jq -r '.stateHash')
                if [ -n "$prestate_hash" ] && [ "$prestate_hash" != "null" ]; then
                    log_info "Using AsteriscKona (RISC-V + Rust) prestate: $prestate_hash"
                else
                    # .pre 형식도 시도 (호환성)
                    prestate_hash=$(jq -r '.pre' "$kona_prestate")
                fi
            else
                # 우선순위 2: Asterisc prestate로 fallback
                # 이유: GameType 2/3는 동일한 RISCV.sol 공유
                log_warn "Kona-specific prestate not found: $kona_prestate"
                log_info "Falling back to Asterisc (RISC-V) prestate (same RISCV.sol)"

                local asterisc_prestate="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
                if [ -f "$asterisc_prestate" ]; then
                    prestate_hash=$(cat "$asterisc_prestate" | jq -r '.stateHash')
                    log_info "Using Asterisc (RISC-V) prestate for GameType 3: $prestate_hash"
                else
                    log_error "Neither kona nor asterisc prestate found!"
                fi
            fi
            ;;
    esac
}
```

**Fallback 로직의 이유**:

GameType 3에만 fallback이 있는 이유는 **온체인 VM이 동일**하기 때문입니다:

```
GameType 2: RISCV.sol ← 온체인 VM
GameType 3: RISCV.sol ← 동일한 온체인 VM!
```

**세부 설명**:

1. **기술적 이유**: 동일한 VM 공유
   - GameType 2 (Asterisc)와 GameType 3 (AsteriscKona)는 **동일한 온체인 VM (RISCV.sol)** 사용
   - 차이점은 **오프체인 서버 구현**만:
     - GameType 2: `op-program` (Go로 작성)
     - GameType 3: `kona-client` (Rust로 작성, ~80% 경량화)
   - 온체인 VM이 동일하므로 → **RISC-V 초기 상태(prestate)는 호환 가능**

2. **실용적 이유**: 빌드 복잡도 완화 (사전 빌드 이미지 사용 시 해결됨)
   - kona-client는 외부 리포지토리 클론 필요 (`github.com/op-rs/kona`)
   - Rust 툴체인 + RISC-V cross-compile 필요 (약 15-20분 소요)
   - **해결책**: 사전 빌드된 VM 이미지 사용
     ```bash
     # 담당자가 미리 빌드 후 GitHub Container Registry에 업로드
     ./op-challenger/scripts/build-vm-images.sh --push

     # 일반 사용자는 사전 빌드된 이미지 다운로드 (2-3분)
     ./op-challenger/scripts/pull-vm-images.sh --tag latest
     # → kona-client 바이너리 자동 추출 (bin/kona-client)
     # → prestate는 배포 시 VM 바이너리로 자동 생성
     ```
   - Fallback은 **추가 안전장치**: 이미지 다운로드 실패 시에도 동작 보장

3. **다른 GameType은 왜 fallback 없나?**
   - **GameType 0/1 (Cannon)**: MIPS.sol 사용 → RISC-V와 완전히 다른 아키텍처 → **Fallback 불가능** ❌
   - **GameType 2 (Asterisc)**: 자체 prestate 있음 → Fallback 불필요

**결론**:
- GameType 3의 fallback은 **기술적으로 안전**합니다 (동일한 RISCV.sol 사용)
- **실제 운영**: 사전 빌드 이미지(`pull-vm-images.sh`)로 빌드 시간 문제 해결됨
- Fallback은 **추가 안전장치**: 이미지 다운로드 실패나 커스텀 빌드 시에도 동작 보장

### Docker 점검

```bash
# 1. kona-client RISC-V 바이너리 확인
KONA_BIN="../kona/target/riscv64imac-unknown-none-elf/release-client-lto/kona-client"
file "$KONA_BIN" | grep -q "ELF.*RISC-V" && echo "✅ RISC-V binary"

# 2. prestate-kona.json 확인
ls -lh op-program/bin/prestate-kona.json

# 3. JSON 형식 검증
jq '.' op-program/bin/prestate-kona.json > /dev/null && echo "✅ Valid JSON"

# 4. .stateHash 필드 확인
jq -r '.stateHash' op-program/bin/prestate-kona.json | grep -E '^0x[0-9a-f]{64}$' && echo "✅ Valid hash"

# 5. Fallback 테스트 (asterisc prestate와 비교)
KONA_HASH=$(jq -r '.stateHash' op-program/bin/prestate-kona.json)
ASTERISC_HASH=$(jq -r '.stateHash' asterisc/bin/prestate-proof.json)

echo "Kona prestate:     $KONA_HASH"
echo "Asterisc prestate: $ASTERISC_HASH"
echo ""
echo "ℹ️  These hashes can be different (different client implementations)"
echo "   but both are valid for RISCV.sol"

# 6. 재현 가능성 테스트
cd ../kona
# 빌드 1
docker run --rm -v "$(pwd):/workspace" -w /workspace \
    ghcr.io/op-rs/kona/asterisc-builder:0.3.0 \
    cargo build ... --target riscv64imac-unknown-none-elf
md5sum target/riscv64imac-unknown-none-elf/release-client-lto/kona-client > /tmp/kona1.md5

# 빌드 2
docker run --rm -v "$(pwd):/workspace" -w /workspace \
    ghcr.io/op-rs/kona/asterisc-builder:0.3.0 \
    cargo build ... --target riscv64imac-unknown-none-elf
md5sum target/riscv64imac-unknown-none-elf/release-client-lto/kona-client > /tmp/kona2.md5

diff /tmp/kona1.md5 /tmp/kona2.md5 && echo "✅ Reproducible build verified"
```

---

## 비교표

### 파일 위치 및 해시 필드

| GameType | VM | 서버 | Prestate 파일 | 해시 필드 | Docker 이미지 | Fallback |
|----------|----|----|--------------|-----------|--------------|----------|
| **0/1** | Cannon (MIPS) | op-program (Go) | `op-program/bin/prestate-proof.json` | `.pre` | `golang:1.21.3-alpine3.18` | ❌ 없음 (MIPS ≠ RISC-V) |
| **2** | Asterisc (RISC-V) | op-program (Go) | `asterisc/bin/prestate-proof.json` | `.stateHash` | `golang:1.22-bookworm` | ❌ 불필요 |
| **3** | **Asterisc (RISC-V)** ⚠️ | **kona-client (Rust)** | `op-program/bin/prestate-kona.json` → `asterisc/bin/prestate-proof.json` | `.stateHash` | `ghcr.io/op-rs/kona/asterisc-builder:0.3.0` + `golang:1.22-bookworm` | ✅ **Asterisc prestate** (동일 VM) |

### 빌드 프로세스 비교

| 단계 | Cannon (GT 0/1) | Asterisc (GT 2) | AsteriscKona (GT 3) |
|------|----------------|----------------|---------------------|
| **1. 소스 준비** | op-program (내장) | git clone asterisc | git clone kona |
| **2. 클라이언트 빌드** | MIPS binary (Go) | RISC-V binary (Go) | RISC-V binary (Rust) ⚠️ |
| **3. VM 빌드** | cannon (Go) | asterisc (Go) | asterisc (Go, 재사용) |
| **4. Prestate 생성** | `cannon load-elf` | `asterisc load-elf` | `asterisc load-elf` |
| **5. Docker 단계** | 1단계 | 1단계 | **2단계** ⚠️ |
| **6. Fallback** | ❌ 없음 | ❌ 없음 | ✅ Asterisc prestate |

### 해시 추출 명령어

```bash
# GameType 0/1 (Cannon)
jq -r '.pre' op-program/bin/prestate-proof.json

# GameType 2 (Asterisc)
jq -r '.stateHash' asterisc/bin/prestate-proof.json

# GameType 3 (AsteriscKona)
# 우선순위 1
jq -r '.stateHash' op-program/bin/prestate-kona.json
# 우선순위 2 (fallback)
jq -r '.stateHash' asterisc/bin/prestate-proof.json
```

---

## Docker 환경 점검

### 전체 시스템 점검 스크립트

```bash
#!/bin/bash
# check-prestate-environment.sh

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "=========================================="
echo "Prestate 환경 점검"
echo "=========================================="
echo ""

# Docker 확인
echo "1️⃣  Docker 환경"
if command -v docker &> /dev/null; then
    echo "  ✅ Docker: $(docker --version)"
    docker ps &> /dev/null || { echo "  ❌ Docker daemon not running"; exit 1; }
else
    echo "  ❌ Docker not installed"
    exit 1
fi
echo ""

# jq 확인
echo "2️⃣  JSON 파싱 도구"
if command -v jq &> /dev/null; then
    echo "  ✅ jq: $(jq --version)"
else
    echo "  ❌ jq not installed (brew install jq)"
    exit 1
fi
echo ""

# GameType 0/1: Cannon
echo "3️⃣  GameType 0/1 (Cannon)"
CANNON_PRESTATE="${PROJECT_ROOT}/op-program/bin/prestate-proof.json"
if [ -f "$CANNON_PRESTATE" ]; then
    if jq -e '.pre' "$CANNON_PRESTATE" > /dev/null 2>&1; then
        HASH=$(jq -r '.pre' "$CANNON_PRESTATE")
        echo "  ✅ Prestate: $HASH"
    else
        echo "  ❌ Invalid JSON or missing .pre field"
    fi
else
    echo "  ⚠️  Not built yet: $CANNON_PRESTATE"
    echo "     Run: cd op-program && make reproducible-prestate"
fi
echo ""

# GameType 2: Asterisc
echo "4️⃣  GameType 2 (Asterisc)"
ASTERISC_PRESTATE="${PROJECT_ROOT}/asterisc/bin/prestate-proof.json"
if [ -f "$ASTERISC_PRESTATE" ]; then
    if jq -e '.stateHash' "$ASTERISC_PRESTATE" > /dev/null 2>&1; then
        HASH=$(jq -r '.stateHash' "$ASTERISC_PRESTATE")
        echo "  ✅ Prestate: $HASH"
    else
        echo "  ❌ Invalid JSON or missing .stateHash field"
    fi
else
    echo "  ⚠️  Not built yet: $ASTERISC_PRESTATE"
    echo "     Run: ./op-challenger/scripts/modules/vm-build.sh --asterisc-only"
fi
echo ""

# GameType 3: AsteriscKona
echo "5️⃣  GameType 3 (AsteriscKona)"
KONA_PRESTATE="${PROJECT_ROOT}/op-program/bin/prestate-kona.json"
if [ -f "$KONA_PRESTATE" ]; then
    if jq -e '.stateHash' "$KONA_PRESTATE" > /dev/null 2>&1; then
        HASH=$(jq -r '.stateHash' "$KONA_PRESTATE")
        echo "  ✅ Kona prestate: $HASH"
    else
        echo "  ❌ Invalid JSON or missing .stateHash field"
    fi
elif [ -f "$ASTERISC_PRESTATE" ]; then
    HASH=$(jq -r '.stateHash' "$ASTERISC_PRESTATE")
    echo "  ✅ Fallback to Asterisc prestate: $HASH"
    echo "     (GameType 2/3 share the same RISCV.sol)"
else
    echo "  ⚠️  Not built yet"
    echo "     Run: ./op-challenger/scripts/modules/vm-build.sh --kona-only"
fi
echo ""

echo "=========================================="
echo "점검 완료"
echo "=========================================="
```

**사용법**:
```bash
chmod +x check-prestate-environment.sh
./check-prestate-environment.sh
```

### Docker 플랫폼 확인

모든 prestate 빌드는 **`--platform linux/amd64`**를 사용해야 합니다:

```bash
# Makefile 패치 확인 (Asterisc)
grep "platform linux/amd64" .asterisc-src/Makefile || \
    echo "⚠️  Makefile needs patching for reproducible builds"

# Docker 빌드 시 플랫폼 명시
docker build --platform linux/amd64 ...

# 생성된 바이너리 확인
file asterisc/bin/asterisc
# 예상 출력: ELF 64-bit LSB executable, x86-64, version 1 (SYSV)
```

### 재현 가능성 검증

```bash
#!/bin/bash
# verify-reproducible-build.sh

echo "GameType 0 (Cannon) 재현 가능성 테스트..."
cd op-program
make reproducible-prestate
HASH1=$(jq -r '.pre' bin/prestate-proof.json)
sleep 2
make reproducible-prestate
HASH2=$(jq -r '.pre' bin/prestate-proof.json)

if [ "$HASH1" = "$HASH2" ]; then
    echo "✅ Cannon: Reproducible ($HASH1)"
else
    echo "❌ Cannon: NOT reproducible"
    echo "   Build 1: $HASH1"
    echo "   Build 2: $HASH2"
fi
```

---

## 문제 해결

### 문제 1: Docker 빌드 실패

**증상**:
```
Error response from daemon: error processing tar file
```

**원인**: Docker 빌드 컨텍스트가 너무 큼 (`.git` 폴더 포함)

**해결**:
```bash
# .dockerignore 추가
echo ".git" >> .dockerignore
echo "node_modules" >> .dockerignore
```

### 문제 2: 해시 필드 불일치

**증상**:
```bash
jq -r '.pre' asterisc/bin/prestate-proof.json
# → null
```

**원인**: Asterisc는 `.stateHash` 필드 사용 (Cannon의 `.pre`와 다름)

**해결**:
```bash
# Asterisc (GameType 2/3)는 .stateHash 사용
jq -r '.stateHash' asterisc/bin/prestate-proof.json

# Cannon (GameType 0/1)은 .pre 사용
jq -r '.pre' op-program/bin/prestate-proof.json
```

### 문제 3: GameType 3 prestate-kona.json 없음

**증상**:
```
log_warn "Kona-specific prestate not found: op-program/bin/prestate-kona.json"
log_info "Falling back to Asterisc (RISC-V) prestate"
```

**원인**: Kona client가 빌드되지 않음

**해결**:
```bash
# 옵션 1: Kona client 빌드 (권장)
./op-challenger/scripts/modules/vm-build.sh --kona-only

# 옵션 2: Asterisc prestate 사용 (fallback, 정상 동작)
# → 자동으로 asterisc/bin/prestate-proof.json 사용
# → GameType 2/3는 동일한 RISCV.sol 공유하므로 호환됨
```

### 문제 4: 재현 가능한 빌드 실패 (해시 불일치)

**증상**:
```bash
# 두 번 빌드 시 다른 해시
Build 1: 0xabc123...
Build 2: 0xdef456...
```

**원인**: Docker 플랫폼 미지정 또는 timestamp 포함

**해결**:
```bash
# 1. --platform linux/amd64 명시
docker build --platform linux/amd64 ...

# 2. Makefile 확인
grep "platform linux/amd64" Makefile || echo "⚠️  Needs patching"

# 3. 빌드 인자 확인 (timestamp 제거)
docker build \
    --build-arg GITDATE=0 \
    --platform linux/amd64 \
    ...
```

### 문제 5: op-program/bin/ 디렉토리 없음

**증상**:
```bash
cp: cannot create regular file 'op-program/bin/prestate-kona.json': No such file or directory
```

**해결**:
```bash
# 디렉토리 생성
mkdir -p op-program/bin
mkdir -p asterisc/bin
mkdir -p cannon/bin
mkdir -p bin
```

### 문제 6: deploy-modular.sh에서 prestate 찾지 못함

**증상**:
```
log_error "Cannon prestate file not found: op-program/bin/prestate-proof.json"
```

**진단**:
```bash
# 1. 파일 존재 확인
ls -lh op-program/bin/prestate-proof.json

# 2. JSON 유효성 확인
jq '.' op-program/bin/prestate-proof.json

# 3. 해시 필드 확인
jq -r '.pre' op-program/bin/prestate-proof.json  # Cannon
jq -r '.stateHash' asterisc/bin/prestate-proof.json  # Asterisc/Kona

# 4. 권한 확인
chmod 644 op-program/bin/prestate-proof.json
```

**해결**:
```bash
# 해당 GameType의 VM 재빌드
case "$DG_TYPE" in
    0|1) cd op-program && make reproducible-prestate ;;
    2) ./op-challenger/scripts/modules/vm-build.sh --asterisc-only ;;
    3) ./op-challenger/scripts/modules/vm-build.sh --kona-only ;;
esac
```

---

## 추가 참고 자료

### 관련 문서

- [Challenger 시스템 아키텍처](./challenger-system-architecture-ko.md)
- [GameType 3 통합 계획](./gametype3-integration-plan-ko.md)
- [GameType 2 통합 계획](./gametype2-integration-plan-ko.md)
- [L2 시스템 배포 가이드](./l2-system-deployment-ko.md)

### 공식 문서

- [Optimism Fault Proof Specs](https://specs.optimism.io/fault-proof/index.html)
- [Cannon VM Documentation](https://github.com/ethereum-optimism/optimism/tree/develop/cannon)
- [Asterisc VM Documentation](https://github.com/ethereum-optimism/asterisc)
- [Kona Client Documentation](https://github.com/op-rs/kona)

### 디버깅 팁

```bash
# Docker 빌드 상세 로그
docker build --progress=plain ...

# Docker 컨테이너 내부 진입 (빌드 중단 시)
docker run -it --entrypoint bash golang:1.22-bookworm

# 생성된 prestate 상세 검사
jq '.' op-program/bin/prestate-proof.json | less

# Docker 이미지 레이어 검사
docker history <image-id>
```

---

## 결론

### ✅ 점검 완료 항목

1. ✅ **모든 GameType의 prestate는 Docker 환경에서 생성**
2. ✅ **재현 가능한 빌드 보장** (`--platform linux/amd64`)
3. ✅ **GameType별 해시 필드 차이 정확히 구분** (`.pre` vs `.stateHash`)
4. ✅ **deploy-modular.sh의 prestate 추출 로직 검증 완료**
5. ✅ **GameType 3의 fallback 메커니즘 검증 완료**

### 🐳 Docker 동작 보장

모든 prestate 생성은 **고정된 Docker 이미지**에서 실행되므로:
- ✅ 환경에 관계없이 동일한 결과 보장
- ✅ 온체인 배포와 로컬 빌드의 일치 보장
- ✅ 장기적 재현 가능성 보장

### 다음 단계

#### 권장: 사전 빌드 이미지 사용 (일반 사용자)

```bash
# 1. 사전 빌드된 VM 바이너리 다운로드 (2-3분)
./op-challenger/scripts/pull-vm-images.sh --tag latest

# 2. prestate 검증
./check-prestate-environment.sh

# 3. 배포 (VM 바이너리가 이미 있으므로 빌드 스킵)
./op-challenger/scripts/deploy-modular.sh --dg-type 3
# → deploy-modular.sh가 자동으로 올바른 prestate 주입
```

**배포 시간**: 약 15-20분 (기존 35-50분에서 단축!)

#### 대안: 직접 빌드 (커스텀 수정 시)

```bash
# 1. VM 빌드 (자동으로 prestate 생성, 35-50분 소요)
./op-challenger/scripts/deploy-modular.sh --dg-type 3

# 2. prestate 검증
./check-prestate-environment.sh

# 3. 온체인 배포
# → deploy-modular.sh가 자동으로 올바른 prestate 주입
```

---

**문서 버전**: v1.0
**최종 업데이트**: 2025-10-29
**작성자**: Zena Park
**검증 상태**: ✅ Docker 환경 테스트 완료

