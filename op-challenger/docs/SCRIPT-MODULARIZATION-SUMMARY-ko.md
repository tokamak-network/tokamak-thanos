# 배포 스크립트 모듈화 완료 요약

**목적**: GameType 2 (Asterisc) VM 자동 빌드 및 배포 스크립트 모듈화

---

## 개요

배포 스크립트를 재사용 가능한 모듈 구조로 분리하여 유지보수성을 개선했습니다.

특히 **GameType 2 (Asterisc) VM 빌드를 자동화**하여, Docker reproducible build를 통해 Linux 호환 바이너리를 자동으로 생성하도록 개선했습니다.

---

## 모듈화 구조

### 디렉토리 구조

```
op-challenger/scripts/
├── lib/
│   └── common.sh                    # 공통 라이브러리 (로깅, 유틸리티)
├── modules/
│   ├── vm-build.sh                  # VM 빌드 모듈 (Cannon + Asterisc)
│   └── genesis.sh                   # Genesis 생성 모듈
├── deploy-modular.sh                # 모듈식 배포 스크립트
├── health-check.sh                  # 헬스 체크 (GameType 2 지원)
├── monitor-challenger.sh            # 모니터링 (GameType 2 지원)
└── docker-compose-full.yml          # Docker 구성 (GameType 2 지원)
```

---

## 생성된 모듈

### 1. lib/common.sh (7.6KB)

**기능**: 모든 스크립트에서 공통으로 사용하는 유틸리티 함수 제공

**주요 함수**:
- **로깅**: `log_info`, `log_success`, `log_warn`, `log_error`, `log_config`
- **명령어 검증**: `check_command`, `check_required_commands`
- **대기 함수**: `wait_for_command`, `wait_for_rpc`
- **환경 변수**: `load_env_file`, `validate_env_vars`
- **파일 작업**: `ensure_dir`, `backup_file`
- **GameType 유틸리티**:
  - `get_gametype_name()` - GameType 번호 → 이름 변환
  - `validate_gametype_tracetype()` - GameType과 TraceType 호환성 검증

**사용 예시**:
```bash
source "${SCRIPT_DIR}/lib/common.sh"

log_info "Starting deployment..."
check_required_commands jq docker cast
validate_gametype_tracetype 2 "asterisc"  # GameType 2는 asterisc만 허용
```

---

### 2. modules/vm-build.sh (11KB)

**기능**: Cannon (MIPS) 및 Asterisc (RISC-V) VM 자동 빌드

**주요 함수**:

#### `build_op_program()`
- op-program 빌드 (모든 VM에서 공유)
- 출력: `op-program/bin/op-program`

#### `build_cannon()`
- Cannon VM (MIPS) 빌드
- Cannon prestate 생성
- 출력:
  - `cannon/bin/cannon`
  - `op-program/bin/prestate-proof.json` (prestate hash 포함)

#### `build_asterisc()` ⭐ **핵심 기능**
Asterisc VM 자동 빌드 (Docker reproducible build):

1. **소스에서 빌드** (`build_asterisc_from_source`) ✅
   - GitHub에서 asterisc 저장소 클론
   - 저장소: `https://github.com/ethereum-optimism/asterisc.git`
   - 빌드 위치: `.asterisc-src/`
   - **Docker reproducible build 사용** (`make reproducible-prestate`)
     - Linux 호환 바이너리 생성 (Docker 컨테이너에서 실행 가능)
     - Prestate 자동 생성 및 검증
   - 설치 위치: `asterisc/bin/asterisc`
   - Prestate 위치: `asterisc/bin/prestate.json`

#### `build_vms(build_cannon, build_asterisc)`
- 전체 빌드 오케스트레이터
- GameType에 따라 필요한 VM만 선택적 빌드

**사용 예시**:
```bash
# Cannon만 빌드 (GameType 0, 1)
build_vms true false

# Asterisc만 빌드 (GameType 2)
build_vms false true

# 둘 다 빌드
build_vms true true
```

**독립 실행**:
```bash
# Cannon만 빌드
./modules/vm-build.sh --cannon-only

# Asterisc만 빌드
./modules/vm-build.sh --asterisc-only

# 둘 다 빌드
./modules/vm-build.sh --all
```

---

### 3. modules/genesis.sh (18KB)

**기능**: Genesis 파일 및 계약 배포 자동화

**주요 함수**:

#### `cleanup_docker_volume_directories()`
- Docker volume mount로 인한 잘못된 디렉토리 정리

#### `cleanup_for_config_change()`
- Game 설정 변경 시 필요한 전체 정리:
  - 배포 캐시 삭제 (`deployments/devnetL1/`)
  - devnetL1.json 삭제 (템플릿 강제 사용)
  - Genesis 파일 삭제
  - Docker 볼륨 삭제

#### `configure_game_settings(prestate_hash)`
- Game 설정 템플릿 수정:
  - `faultGameMaxClockDuration` (환경 변수로 설정 가능)
  - `faultGameWithdrawalDelay` (환경 변수로 설정 가능)
  - `faultGameAbsolutePrestate` (VM에서 생성된 prestate hash 주입)

#### `run_devnet_up()`
- `make devnet-up` 실행하여 Genesis 생성 및 계약 배포

#### `wait_for_genesis_files()`
- Genesis 파일 생성 대기 (최대 3분):
  - `genesis-l1.json`
  - `genesis-l2.json`
  - `rollup.json`
  - `addresses.json`

#### `verify_genesis_files()`
- 생성된 Genesis 및 설정 검증
- 계약 주소 확인 (DisputeGameFactoryProxy)

#### `stop_devnet()`
- devnet 종료 및 볼륨 정리

#### `generate_genesis(prestate_hash)`
- 전체 Genesis 생성 프로세스 오케스트레이터

**사용 예시**:
```bash
# Prestate hash와 함께 Genesis 생성
generate_genesis "0x038512e02c4c3f7bdaec27d00edf55b7155e0905301e1a88083e4e0a6764d54c"

# 커스텀 게임 설정과 함께 Genesis 생성
FAULT_GAME_MAX_CLOCK_DURATION=120 \
FAULT_GAME_WITHDRAWAL_DELAY=604800 \
generate_genesis "0x038512..."
```

**독립 실행**:
```bash
./modules/genesis.sh --prestate 0x038512e02c4c3f7bdaec27d00edf55b7155e0905301e1a88083e4e0a6764d54c
```

---

### 4. deploy-modular.sh (새 배포 스크립트)

**기능**: 모듈을 사용한 간결한 배포 스크립트

**배포 흐름**:
```bash
1. 인자 파싱 (--mode, --dg-type)
2. 환경 설정 로깅
3. GameType/TraceType 호환성 검증
4. VM 빌드 (GameType에 따라 Cannon 또는 Asterisc)
5. Genesis 생성 (local 모드만)
6. Docker 서비스 시작
```

**사용 예시**:
```bash
# GameType 0 (Cannon) 배포
./deploy-modular.sh

# GameType 2 (Asterisc) 배포 ⭐
./deploy-modular.sh --dg-type 2

# 커스텀 게임 설정과 함께 배포
FAULT_GAME_MAX_CLOCK_DURATION=120 ./deploy-modular.sh --dg-type 2
```

**장점**:
- 기존 1392줄 → 약 200줄로 축소
- 모듈 재사용으로 유지보수 용이
- GameType 2 (Asterisc) 자동 빌드 지원

---

## GameType 2 (Asterisc) 자동 빌드 흐름

### 전체 프로세스

```
사용자: ./deploy-modular.sh --dg-type 2
   ↓
1. GameType 검증 (2 = Asterisc)
   ↓
2. build_vms(false, true) 호출
   ↓
3. build_asterisc() 실행 ✅
   ├─→ git clone https://github.com/ethereum-optimism/asterisc.git → .asterisc-src/
   ├─→ cd .asterisc-src
   ├─→ make reproducible-prestate (Docker reproducible build)
   │   ├─→ Docker로 Linux 호환 바이너리 생성
   │   ├─→ Prestate 자동 생성 및 검증
   │   └─→ bin/ 디렉토리에 결과물 출력
   ├─→ cp bin/asterisc → asterisc/bin/asterisc
   ├─→ cp bin/prestate.json → asterisc/bin/prestate.json
   └─→ Prestate hash 추출 (.stateHash 필드)
   ↓
4. generate_genesis(prestate_hash) 실행
   - 템플릿에 prestate hash 주입
   - make devnet-up (계약 배포 + Genesis 생성)
   - make devnet-down (정리)
   ↓
5. Docker 서비스 시작
   - CHALLENGER_TRACE_TYPE=asterisc 환경 변수 설정
   - Asterisc 바이너리 경로 마운트
   ↓
6. 배포 완료
```

---


## 환경 변수

### GameType 설정
```bash
# GameType 선택 (0=Cannon, 1=PermissionedCannon, 2=Asterisc, 254=Fast, 255=Alphabet)
DG_TYPE=2

# TraceType (GameType에 따라 자동 설정됨)
# GameType 0,1 → cannon
# GameType 2 → asterisc
# GameType 255 → alphabet
CHALLENGER_TRACE_TYPE=asterisc
```

### Game 설정 (옵션)
```bash
# Clock duration (초 단위, 기본값: 302400 = 3.5일)
FAULT_GAME_MAX_CLOCK_DURATION=120

# Withdrawal delay (초 단위, 기본값: 604800 = 7일)
FAULT_GAME_WITHDRAWAL_DELAY=86400

# Proposal interval
PROPOSAL_INTERVAL=30
```

### Asterisc 경로 (자동 설정됨)
```bash
ASTERISC_BIN=/path/to/asterisc/bin
ASTERISC_PRESTATE=/path/to/asterisc/bin/prestate.json
```

---

## 사용 시나리오

### 시나리오 1: GameType 0 (Cannon) 배포
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/op-challenger/scripts

# GameType 0 (Cannon) 배포
./deploy-modular.sh --dg-type 0
```

### 시나리오 2: GameType 2 (Asterisc) 배포 ⭐
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/op-challenger/scripts

# Asterisc VM 자동 빌드 및 배포
./deploy-modular.sh --dg-type 2

# 또는 환경 변수로
DG_TYPE=2 ./deploy-modular.sh
```

### 시나리오 3: 커스텀 게임 설정
```bash
# 2분 clock duration, GameType 2 배포
FAULT_GAME_MAX_CLOCK_DURATION=120 \
./deploy-modular.sh --dg-type 2
```

### 시나리오 4: VM만 빌드 (배포 없이)
```bash
# Asterisc만 빌드
./modules/vm-build.sh --asterisc-only

# 모든 VM 빌드
./modules/vm-build.sh --all
```

### 시나리오 5: Genesis만 재생성
```bash
# Prestate hash 얻기
PRESTATE=$(jq -r '.pre' op-program/bin/prestate-proof.json)

# Genesis 재생성
./modules/genesis.sh --prestate "$PRESTATE"
```

---

## 모니터링 및 검증

### 1. 헬스 체크 (GameType 2 지원)
```bash
./health-check.sh

# 출력 예시:
# ━━━ GameType Configuration ━━━
# GameType 0 (CANNON) - Deployed: 0x1234...
# GameType 2 (ASTERISC/RISCV) - Deployed: 0x5678... ⭐
#   └─ RISCV VM: 0x9abc...
```

### 2. 챌린저 모니터링 (GameType 2 지원)
```bash
./monitor-challenger.sh

# 출력 예시:
# DisputeGameFactory Configuration:
#   ┌─ GameType 0 (CANNON - MIPS VM)
#   │  Address: 0x1234...
#   │
#   └─ GameType 2 (ASTERISC - RISC-V VM) ⭐
#      Address: 0x5678...
#      RISCV VM: 0x9abc...
```

### 3. Docker 서비스 확인
```bash
docker ps

# Asterisc 환경 변수 확인
docker inspect scripts-op-challenger-1 | grep -A 10 ASTERISC
```

---

## 문제 해결

### 문제 1: Asterisc 빌드 실패
**증상**: `build_asterisc_from_source()` 실패 (Docker reproducible build 에러)

**원인**:
- Docker가 실행되지 않음
- 네트워크 문제로 asterisc 저장소 클론 실패
- Dockerfile.repro가 asterisc 저장소에 없음

**해결책**:
```bash
# 1. Docker 실행 확인
docker info

# 2. 수동으로 Asterisc 빌드 시도 (Docker reproducible build)
cd /Users/zena/tokamak-projects/tokamak-thanos
rm -rf .asterisc-src
git clone https://github.com/ethereum-optimism/asterisc.git .asterisc-src
cd .asterisc-src

# Docker reproducible build 실행
make reproducible-prestate

# 빌드된 파일들 복사
mkdir -p ../asterisc/bin
cp bin/asterisc ../asterisc/bin/
cp bin/prestate.json ../asterisc/bin/

# Prestate hash 확인
cat bin/prestate.json | jq '.stateHash'
```

### 문제 2: Genesis 파일 생성 안됨
**증상**: `wait_for_genesis_files()` 타임아웃

**해결책**:
```bash
# devnet 로그 확인
cd ops-bedrock
docker compose logs -f

# 수동 Genesis 생성
cd /Users/zena/tokamak-projects/tokamak-thanos
make devnet-up

# Genesis 파일 확인
ls -lh .devnet/*.json
```

### 문제 3: GameType 2 계약 배포 안됨
**증상**: `health-check.sh`에서 GameType 2 주소가 0x0

**원인**: devnetL1-template.json 설정 문제 또는 배포 스크립트 미실행

**해결책**:
```bash
# 계약 배포 상태 확인
cat .devnet/addresses.json | jq '.DisputeGameFactoryProxy'

# GameType 2 구현체 확인
cast call --rpc-url http://localhost:8545 \
  $(jq -r '.DisputeGameFactoryProxy' .devnet/addresses.json) \
  "gameImpls(uint32)(address)" 2
```

---

## 향후 개선 사항

### 1. 추가 모듈 분리
현재는 VM 빌드와 Genesis만 모듈화했습니다. 향후:
- `modules/l1-deploy.sh` - L1 계약 배포
- `modules/sequencer-deploy.sh` - Sequencer 배포
- `modules/challenger-deploy.sh` - Challenger 배포

### 2. Asterisc 사전 빌드 바이너리
현재는 Docker reproducible build를 사용하므로 첫 빌드 시 시간이 소요됩니다:
- GitHub Releases에 사전 빌드 바이너리 업로드
- 다운로드 우선 시도, 실패 시 Docker build로 fallback

### 3. 설정 파일 지원
현재는 환경 변수로만 설정:
- YAML/TOML 설정 파일 지원
- 여러 환경(dev, staging, prod) 프로필 지원

---

## 요약

### 달성한 것 ✅
1. **모듈화**: 재사용 가능한 모듈 구조 (lib/ + modules/)
2. **Asterisc 자동 빌드**: GameType 2 선택 시 Docker reproducible build로 자동 빌드
   - Linux 호환 바이너리 생성
   - Prestate 자동 생성 및 검증
3. **개선된 로깅**: 색상 코딩 및 config trace 로그
4. **GameType 호환성 검증**: GameType과 TraceType 자동 검증
5. **모니터링 개선**: health-check, monitor 스크립트에 GameType 2 지원 추가

### 주요 파일
- `lib/common.sh` - 공통 라이브러리
- `modules/vm-build.sh` - VM 빌드 (Cannon + Asterisc Docker reproducible build)
- `modules/genesis.sh` - Genesis 생성
- `deploy-modular.sh` - 모듈식 배포 스크립트
- `health-check.sh` - GameType 2 지원 추가
- `monitor-challenger.sh` - GameType 2 표시 및 Prestate 검증 추가

### 명령어 요약
```bash
# GameType 2 (Asterisc) 전체 배포
./deploy-modular.sh --dg-type 2

# Asterisc VM만 빌드
./modules/vm-build.sh --asterisc-only

# 헬스 체크
./health-check.sh

# 모니터링
./monitor-challenger.sh
```

---
