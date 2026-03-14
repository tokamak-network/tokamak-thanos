# Justfiles Dependency Graph

## Overview

```
default.just (terminal)
    └── git.just
        └── go.just (7개 패키지에서 import)
```

## File Structure

### 1. **default.just** (Terminal)
- **Purpose**: Base configuration for all justfiles
- **Exports**:
  - `PARALLEL_JOBS` - CPU 코어 수 기반 병렬 작업 수
  - `MAP_JUST` - 병렬 작업 실행을 위한 헬퍼
- **Dependencies**: None

### 2. **git.just**
- **Purpose**: Git 관련 메타데이터 및 버전 관리
- **Imports**: `default.just`
- **Exports**:
  - `GITCOMMIT` - 현재 커밋 해시
  - `GITDATE` - 커밋 날짜 (Unix timestamp)
  - `VERSION` - 버전 태그 (non-rc 우선)
  - `VERSION_META` - 추가 버전 정보
- **Key Functions**: Git 태그 파싱, 버전 결정

### 3. **go.just** (Target)
- **Purpose**: Go 빌드 시스템 핵심 utilities
- **Imports**: `git.just`
- **Exports**:
  - **Environment Variables**:
    - `GOOS`, `GOARCH` - 기본 OS/아키텍처
    - `TARGETOS`, `TARGETARCH` - 빌드 대상 (우선순위: env > GOOS/GOARCH)
    - `GORACE` - Race detector 활성화 플래그

  - **Functions** (모두 `[private]`):
    - `go_build(BIN, PKG, *FLAGS)` - Go 바이너리 컴파일
    - `go_test(SELECTOR, *FLAGS)` - Go 테스트 실행
    - `go_fuzz(FUZZ, TIME, PKG)` - Fuzz 테스트
    - `go_generate(SELECTOR, *FLAGS)` - Code generation (go generate)

  - **Internal**:
    - `_EXTRALDFLAGS` - macOS 호환성 링커 플래그
    - `_GORACE_FLAG` - Race detector 플래그

## Dependents (go.just를 import하는 패키지)

```
op-node/justfile           ✓ imports ../justfiles/go.just
op-batcher/justfile        ✓ imports ../justfiles/go.just
op-proposer/justfile       ✓ imports ../justfiles/go.just
op-challenger/justfile     ✓ imports ../justfiles/go.just
op-dispute-mon/justfile    ✓ imports ../justfiles/go.just
op-program/justfile        ✓ imports ../justfiles/go.just
cannon/justfile            ✓ imports ../justfiles/go.just
```

## Import Chain Details

```
┌─────────────────────────────────────────────────────────┐
│ op-node/justfile (and 6 others)                         │
└─────────────────────────────────────────────────────────┘
                     │
                     │ import '../justfiles/go.just'
                     ↓
┌─────────────────────────────────────────────────────────┐
│ go.just                                                 │
│                                                         │
│ Provides:                                               │
│ - go_build(BIN, PKG, *FLAGS)                            │
│ - go_test(SELECTOR, *FLAGS)                             │
│ - go_fuzz(FUZZ, TIME, PKG)                              │
│ - go_generate(SELECTOR, *FLAGS)                         │
│ - GOOS, GOARCH, TARGETOS, TARGETARCH                    │
│ - VERSION, GITCOMMIT, GITDATE (via git.just)            │
└─────────────────────────────────────────────────────────┘
                     │
                     │ import 'git.just'
                     ↓
┌─────────────────────────────────────────────────────────┐
│ git.just                                                │
│                                                         │
│ Provides:                                               │
│ - GITCOMMIT, GITDATE                                    │
│ - VERSION (semantic versioning from tags)               │
│ - VERSION_META                                          │
└─────────────────────────────────────────────────────────┘
                     │
                     │ import 'default.just'
                     ↓
┌─────────────────────────────────────────────────────────┐
│ default.just                                            │
│                                                         │
│ Provides:                                               │
│ - PARALLEL_JOBS                                         │
│ - MAP_JUST                                              │
│ - Shell configuration                                   │
└─────────────────────────────────────────────────────────┘
```

## Transitive Dependencies

| 패키지 | 직접 의존 | 간접 의존 | 최종 의존 |
|--------|---------|----------|----------|
| op-node | go.just | git.just, default.just | 3개 파일 |
| op-batcher | go.just | git.just, default.just | 3개 파일 |
| op-proposer | go.just | git.just, default.just | 3개 파일 |
| op-challenger | go.just | git.just, default.just | 3개 파일 |
| op-dispute-mon | go.just | git.just, default.just | 3개 파일 |
| op-program | go.just | git.just, default.just | 3개 파일 |
| cannon | go.just | git.just, default.just | 3개 파일 |

## Build Flow

```
any_go_package/justfile
    ↓ [호출]
go_build(BIN, PKG, *FLAGS)
    ↓ [사용]
TARGETOS, TARGETARCH (from go.just)
    ↓ [요구]
git.just
    ↓ [요구]
default.just (shell config)
    ↓ [실행]
go build -o {{BIN}} {{PKG}}
```

## Environment Variable Resolution

```
TARGETOS 결정 순서:
1. env('TARGETOS')        ← 외부 환경변수 (최우선)
2. GOOS                   ← env('GOOS', `go env GOOS`)
3. go env GOOS의 기본값

TARGETARCH 결정 순서:
1. env('TARGETARCH')      ← 외부 환경변수 (최우선)
2. GOARCH                 ← env('GOARCH', `go env GOARCH`)
3. go env GOARCH의 기본값
```

## Usage Example

```bash
# op-node 빌드
cd op-node
just op-node

# 다른 아키텍처로 빌드
TARGETOS=linux TARGETARCH=arm64 just op-node

# Race detector 활성화
GORACE=1 just op-node
```

## Critical Symbols

| Symbol | Source | Used By | Purpose |
|--------|--------|---------|---------|
| `go_build` | go.just | op-node, op-batcher, ... | Go 바이너리 컴파일 |
| `go_test` | go.just | op-node test target | 테스트 실행 |
| `GITCOMMIT` | git.just | 모든 Op 패키지 | 버전 정보 임베딩 |
| `VERSION` | git.just | 모든 Op 패키지 | 의미있는 버전 태그 |
| `PARALLEL_JOBS` | default.just | 병렬 작업 | CPU 기반 병렬화 |

## Missing Files (Restored)

✅ **Created**:
- `justfiles/default.just`
- `justfiles/git.just`
- `justfiles/go.just`

이 파일들은 upstream Optimism에서 가져온 것이며, contracts-bedrock 패키지 이동 시에 제거되었습니다.
