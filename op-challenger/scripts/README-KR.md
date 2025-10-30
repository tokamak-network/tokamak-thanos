# L2 시스템 배포 스크립트 (두 개의 Challenger 포함)

이 디렉토리는 Optimism L2 Rollup 시스템의 완전한 배포를 위한 스크립트들을 포함합니다.

## ⭐ L2 시스템 아키텍처 (두 개의 Challenger 포함)

이 배포 스크립트는 **시퀀서용 Challenger와 독립적인 제3자 Challenger**를 모두 구현합니다:

```
Sequencer 스택:                 독립적인 Challenger 스택:
├─ sequencer-l2 (L2 geth)       ├─ challenger-l2 (별도 L2 geth!) ⭐
├─ sequencer-op-node            ├─ challenger-op-node (Follower) ⭐
├─ op-batcher                   └─ op-challenger (독립 검증자)
├─ op-proposer
└─ sequencer-challenger ⭐
```

**왜 두 개의 Challenger가 필요한가요?**

### 1️⃣ sequencer-challenger (시퀀서용)
```yaml
연결: sequencer-op-node + sequencer-l2 (Sequencer 스택 공유)
역할: 사용자 인출 보장을 위한 게임 close
설정: --selective-claim-resolution 없음 (모든 게임 close)
```

**역할:**
- ✅ 자신이 제안한 Output을 직접 방어
- ✅ **게임을 적시에 resolve하고 close → 사용자 인출 보장** ⭐
- ✅ 별도 L2 노드 불필요 (sequencer-l2를 공유 사용)
- ⚠️ 게임이 close되지 않으면 사용자 자금 인출 불가!

### 2️⃣ op-challenger (독립적인 제3자 검증자)
```yaml
연결: challenger-op-node + challenger-l2 (완전히 독립적인 스택!)
역할: 신뢰 최소화 검증 및 잘못된 Output 챌린지
설정: --selective-claim-resolution=true
```

**역할:**
- ✅ 잘못된 Output에 대한 검증 및 챌린지
- ✅ 자신이 참여한 게임 방어
- ✅ **별도 DB로 L1에서 독립적으로 L2 재구성 (신뢰 최소화)** ⭐

**⚠️ 현재 동작 (코드 이슈):**
- `--selective-claim-resolution=true` 설정에도 불구하고:
  - ✅ Individual claims resolve: 자신이 만든 claim만 resolve (정상)
  - ❌ Game resolve: 모든 게임 resolve 시도 (설정 무시됨)
- 결과: Sequencer Challenger와 경쟁하여 revert 발생 (불필요한 가스비 소모)
- 원인: `op-challenger/game/fault/agent.go`의 `tryResolve()` 함수에 selective 체크 누락
- 📝 **향후 수정 필요**: Game resolve에도 selective 체크 추가 필요

**왜 독립적인 스택이 필요한가요?**
- Sequencer와 노드를 공유하면 Sequencer가 악의적일 때 Challenger도 속을 수 있음 ❌
- 제3자 검증자로서 진정한 Fault Proof 검증을 위해 필수 ✅
- Challenger 시스템의 올바른 동작을 검증하고 테스트하기 위한 환경 ⭐

> 💡 **요약**: `sequencer-challenger`는 사용자 인출을 위해 필수이고, `op-challenger`는 탈중앙화 검증을 위해 필수입니다.

## 📚 참고 문서

배포 전에 다음 문서를 반드시 읽어보시기 바랍니다:

- **[Challenger 시스템 아키텍처](../docs/challenger-system-architecture-ko.md)** ⭐ - 전체 시스템 구성 및 GameType별 아키텍처 상세 설명
- **[L2 시스템 배포 가이드](../docs/l2-system-deployment-ko.md)** - 전체 배포 아키텍처 및 상세 설명
- **[L2 시스템 아키텍처](../docs/l2-system-architecture-ko.md)** - L2 시스템의 기본 구조


## 🔨 이미지 빌드 및 공유 (개발담당자용)

**일반 사용자는 건너뛰세요!** 이미 빌드된 이미지를 사용하려면 [빠른 시작](#-빠른-시작)으로 이동하세요.

### 빌드되는 이미지:

**VM 이미지** (prestates 포함):
- `vm-cannon` - Cannon VM + prestate (GameType 0, 1)
- `vm-asterisc` - Asterisc VM + prestate (GameType 2)
- `vm-op-program` - op-program + clients + prestates
- `vm-kona-client` - kona-client + prestate (GameType 3) ⭐

**OP Stack 이미지** (수정된 코드 반영):
- `op-node`, `op-batcher`, `op-proposer`, `op-challenger`

### 빌드 & 푸시:

```bash
# VM 이미지 빌드 & 푸시 (30-40분)
./op-challenger/scripts/build-vm-images.sh --no-cache --push

# 또는 op-challenger 코드 수정 후 OP Stack만 재빌드:
docker build -f ops/docker/op-stack-go/Dockerfile \
  --target op-challenger-target --platform linux/amd64 \
  -t ghcr.io/zena-park/op-challenger:latest .
docker push ghcr.io/zena-park/op-challenger:latest
```

> 📚 **자세한 내용**: [VM 이미지 빌드 및 공유](../docs/vm-image-build-and-share-guide.md)

---

## 🚀 빠른 시작

```bash
# 1. 리포지토리 클론
git clone https://github.com/tokamak-network/tokamak-thanos.git
cd tokamak-thanos
git checkout feature/challenger-gametype3


# 2. 초기화
./op-challenger/scripts/cleanup.sh

# 3. VM 바이너리 다운로드 (2-3분) (vm 삭제하지 않았으면 계속 재사용함. 한번만 실행가능)
./op-challenger/scripts/pull-vm-images.sh --tag latest
# ✅ cannon, asterisc, op-program, kona-client 다운로드 완료
# 또는 특정 태그: ./op-challenger/scripts/pull-vm-images.sh --tag v1.0.0

# 4. 배포
./op-challenger/scripts/deploy-modular.sh --dg-type 3

# 커스터마이징 설정 (개발 및 테스트)
FAULT_GAME_MAX_CLOCK_DURATION=150 \
FAULT_GAME_WITHDRAWAL_DELAY=3600 \
PROPOSAL_INTERVAL=30s \
./op-challenger/scripts/deploy-modular.sh --dg-type 3

# 4. 상태 확인
./op-challenger/scripts/health-check.sh
./op-challenger/scripts/monitor-challenger.sh
```

**총 배포 시간**: 약 15-20분

> 📚 **GameType 선택**: [GameType 가이드](#-gametype-선택-가이드)
> 📚 **VM 이미지 빌드**: [VM 이미지 빌드 및 공유](../docs/vm-image-build-and-share-guide.md)
> 📚 **Prestate 이해**: [Prestate 생성 가이드](../docs/prestate-generation-guide-ko.md)

## 📋 스크립트 설명

### 🆕 deploy-modular.sh (메인 배포 스크립트)

모듈화된 L2 시스템 배포 스크립트. GameType별 자동 설정 및 검증을 수행합니다.

```bash
# 사용법
./deploy-modular.sh [옵션]

# 옵션
--dg-type TYPE          # GameType (0, 1, 2, 3, 254, 255, 기본: 0)
--mode MODE             # 배포 모드: local|existing (기본: local)

# 예시
./deploy-modular.sh                    # GameType 0 (기본)
./deploy-modular.sh --dg-type 3        # GameType 3

FAULT_GAME_MAX_CLOCK_DURATION=150 \
FAULT_GAME_WITHDRAWAL_DELAY=3600 \
PROPOSAL_INTERVAL=30s \
./op-challenger/scripts/deploy-modular.sh --dg-type 3  2>&1 | tee deploy-debug.log
# 커스텀 설정
```

**주요 기능**:
- ✅ GameType별 자동 VM 및 prestate 검증
- ✅ respectedGameType 자동 설정 (OptimismPortal2와 동기화)
- ✅ 모든 GameType의 prestate 검증 (Challenger는 모든 타입 지원)
- ✅ 배포 전 설정 무결성 검증
- ✅ 모듈식 구조 (lib/ + modules/)

> 📚 **자세한 내용**: [스크립트 모듈화 요약](../docs/SCRIPT-MODULARIZATION-SUMMARY-ko.md)

---

### cleanup.sh

배포된 시스템을 정리하고 초기화합니다.

**⚠️ 중요**: 기본 동작이 `--all`로 설정되어 있어 `.env`와 Genesis 파일도 함께 삭제됩니다.

```bash
# 사용법
./cleanup.sh [옵션]

# 옵션
--all-containers       # 모든 관련 컨테이너 정리 (devnet, kurtosis 등 포함)
--rebuild              # Docker 이미지 및 빌드 캐시 제거 (완전 재빌드 강제) ⭐
--keep-env             # 환경 변수 파일 (.env) 유지
--keep-genesis         # Genesis 파일 유지
--clean-vm-builds      # VM 바이너리 삭제 (기본: 보호됨) 🆕
--clean-pulled-images  # Registry에서 다운로드한 이미지 삭제 (기본: 보호됨) 🆕
--help                 # 도움말 표시

# 예시

# 기본 정리 (VM 바이너리와 pulled 이미지는 보호됨) ⭐
./cleanup.sh

# 완전 재빌드 (로컬 빌드 이미지만 삭제, ghcr.io 이미지는 보호) ⭐
./cleanup.sh --rebuild

# VM 바이너리까지 삭제 (pull-vm-images.sh로 다운로드한 것도 삭제)
./cleanup.sh --clean-vm-builds

# Registry 이미지까지 삭제 (ghcr.io 이미지 삭제)
./cleanup.sh --clean-pulled-images

# 완전 삭제 (VM + 이미지 + 컨테이너)
./cleanup.sh --clean-vm-builds --clean-pulled-images --all-containers --rebuild

# .env는 유지하고 Genesis만 삭제
./cleanup.sh --keep-env

# Genesis는 유지하고 .env만 삭제
./cleanup.sh --keep-genesis

# 컨테이너와 볼륨만 삭제 (.env와 Genesis 모두 유지)
./cleanup.sh --keep-env --keep-genesis
```

**⭐ 기본값 변경 (중요)**:
- VM 바이너리 **보호** (pull-vm-images.sh로 다운로드한 것)
- Registry 이미지 **보호** (ghcr.io 이미지)
- 명시적으로 삭제하려면: `--clean-vm-builds`, `--clean-pulled-images`

---

### build-vm-images.sh ⭐

VM 바이너리와 prestate를 Docker에서 재현 가능하게 빌드하고 레지스트리에 푸시합니다 (개발자용).

```bash
# 사용법
./build-vm-images.sh [옵션]

# 옵션
--push              # 빌드 후 레지스트리에 푸시
--registry URL      # 레지스트리 URL (기본: ghcr.io/zena-park)
--tag TAG           # 이미지 태그 (기본: 현재 Git 커밋)
--no-cache          # 캐시 없이 강제 재빌드
--asterisc-only     # Asterisc 이미지만 빌드
--kona-only         # Kona 이미지만 빌드
--help              # 도움말 표시

# 예시
./build-vm-images.sh                          # 모든 이미지 로컬 빌드
./build-vm-images.sh --push                   # 빌드 + 푸시 (latest 태그 자동 추가!)
./build-vm-images.sh --tag v1.0.0 --push      # 특정 태그로 푸시
./build-vm-images.sh --no-cache --push        # 캐시 없이 재빌드 + 푸시
./build-vm-images.sh --kona-only --push       # Kona만 빌드 + 푸시
./build-vm-images.sh --asterisc-only          # Asterisc만 로컬 빌드
```

**빌드되는 이미지:**
- `vm-cannon` - Cannon (MIPS VM) + Cannon prestate
- `vm-asterisc` - Asterisc (RISC-V VM) + Asterisc prestate
- `vm-op-program` - op-program + MIPS/RISC-V clients + prestates
- `vm-kona-client` - kona-client (Rust) + kona prestate ⭐
- `op-challenger`, `op-node`, `op-batcher`, `op-proposer` - OP Stack 서비스

**자동 태그 생성:**
- Primary 태그: 지정한 태그 또는 Git 커밋 해시
- `latest` 태그: **자동으로 추가됨** ⭐ (가장 최근 빌드가 latest가 됨)

**주요 기능:**
- ✅ Docker 재현 가능한 빌드 (reproducible builds)
- ✅ **모든 GameType의 prestate 자동 생성** ⭐
  - Cannon prestate (GameType 0, 1)
  - Asterisc prestate (GameType 2)
  - Kona prestate (GameType 3)
- ✅ latest 태그 자동 업데이트
- ✅ 레지스트리 푸시 지원

**빌드 시간:**
- 전체 빌드: 약 30-40분 (첫 빌드)
- 증분 빌드: 약 5-10분 (캐시 활용)

**레지스트리 권한:**
```bash
# GitHub Container Registry 로그인 (푸시 전 필요)
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# 또는 Personal Access Token 사용
# Settings → Developer settings → Personal access tokens
# Scopes: write:packages, read:packages, delete:packages
```

**푸시 후 Public 설정:**
1. https://github.com/zena-park?tab=packages
2. 각 패키지 선택
3. Settings → Change visibility → Make public

> 📚 **자세한 내용**: [VM 이미지 빌드 및 공유](../docs/vm-image-build-and-share-guide.md)

---

### pull-vm-images.sh

레지스트리에서 사전 빌드된 VM 바이너리와 prestate를 다운로드합니다 (권장).

```bash
# 사용법
./pull-vm-images.sh [옵션]

# 옵션
--tag TAG           # 이미지 태그 (기본: 현재 Git 커밋)
--commit HASH       # 특정 커밋의 이미지
--registry URL      # 레지스트리 URL (기본: ghcr.io/zena-park)
--help              # 도움말 표시

# 예시
./pull-vm-images.sh --tag latest          # latest 태그 사용 (권장) ⭐
./pull-vm-images.sh --tag v1.0.0          # 특정 버전
./pull-vm-images.sh --commit 47efca11     # 특정 커밋
./pull-vm-images.sh                       # 현재 커밋 (자동)
```

**다운로드되는 파일:**
- `cannon/bin/cannon` - MIPS VM
- `asterisc/bin/asterisc` - RISC-V VM
- `op-program/bin/op-program` - Fault Proof 서버
- `op-program/bin/prestate-proof.json` - Cannon prestate (배포용) ⭐
- `asterisc/bin/prestate-proof.json` - Asterisc prestate (배포용) ⭐
- `bin/kona-client` - Rust 서버 (GameType 3)
- `op-program/bin/prestate-kona.json` - Kona prestate (배포용) ⭐
- `bin/prestate.bin.gz` - Kona prestate (런타임용, GameType 3 필수) 🆕

**다운로드 시간:** 약 2-3분

---

### setup-env.sh

환경 변수와 배포에 필요한 지갑을 생성합니다.

```bash
# 사용법
./setup-env.sh [옵션]

# 옵션
--mode MODE             # 배포 모드: local|sepolia|mainnet (기본: local)
--output FILE           # 환경 변수 파일 경로 (기본: .env)
--help                  # 도움말 표시

# 예시
./setup-env.sh --mode local
./setup-env.sh --mode sepolia --output .env.sepolia
```

**생성되는 지갑:**
- Admin: 시스템 관리자 계정
- Batcher: L1에 배치 데이터를 제출하는 계정
- Proposer: L1에 Output을 제출하는 계정
- Sequencer: L2 블록을 생성하는 계정
- Challenger: 독립적인 제3자 검증자 계정 ⭐
- Sequencer Challenger: 시퀀서 스택용 Challenger 계정 (사용자 인출 보장) ⭐

---

## 🎮 GameType 선택 가이드

| GameType | VM | 서버 | 특징 |
|----------|-----|------|------|
| **0** | MIPS | op-program (Go) | 안정적 (기본값) |
| **2** | RISC-V | op-program (Go) | 최신 아키텍처 |
| **3** | RISC-V | kona-client (Rust) | 경량화, ZK 준비 🆕 |

```bash
./deploy-modular.sh                    # GameType 0 (기본)
./deploy-modular.sh --dg-type 2        # GameType 2
./deploy-modular.sh --dg-type 3        # GameType 3

# 커스텀 설정
FAULT_GAME_MAX_CLOCK_DURATION=150 PROPOSAL_INTERVAL=30s \
./deploy-modular.sh --dg-type 3
```

> 📚 **상세 비교**: [Challenger 시스템 아키텍처](../docs/challenger-system-architecture-ko.md)

---

### health-check.sh

배포된 모든 서비스의 상태를 확인합니다.

```bash
# 사용법
./health-check.sh

# 확인 항목
# - Docker 컨테이너 상태
# - L1/L2 RPC 응답
# - 동기화 상태
# - 블록 높이
# - 피어 수
# - 에러 로그
```

### monitor-challenger.sh

Challenger의 실시간 활동, 게임 참여, 동기화 상태를 모니터링합니다.

```bash
# 사용법
./monitor-challenger.sh [mode]

# 모드: summary | config | sync | logs | games | errors | metrics
# 예시
./monitor-challenger.sh              # 전체 대시보드
./monitor-challenger.sh config       # 시스템 설정 및 Prestate 검증 ⭐
```

**주요 모니터링 항목**:
- 컨테이너 상태 및 동기화
- GameType 설정 및 Prestate 검증 (온체인 vs 로컬)
- 게임 참여 통계 및 Challenger 액션
- 독립 스택 검증 (Sequencer 분리 확인)

---

## 📊 서비스 엔드포인트

배포 완료 후 다음 엔드포인트를 사용할 수 있습니다:

| 서비스 | 포트 | URL | 설명 |
|--------|------|-----|------|
| **L1 스택** ||||
| L1 RPC | 8545 | http://localhost:8545 | L1 Ethereum RPC |
| L1 WebSocket | 8546 | ws://localhost:8546 | L1 WebSocket |
| **Sequencer 스택** ||||
| Sequencer L2 RPC | 9545 | http://localhost:9545 | L2 사용자 RPC |
| Sequencer op-node RPC | 7545 | http://localhost:7545 | Sequencer op-node RPC |
| op-batcher RPC | 6545 | http://localhost:6545 | op-batcher RPC |
| op-proposer RPC | 6546 | http://localhost:6546 | op-proposer RPC |
| sequencer-challenger RPC | 6548 | http://localhost:6548 | Sequencer Challenger ⭐ |
| **독립 Challenger 스택 ⭐** ||||
| Challenger L2 RPC | 9546 | http://localhost:9546 | Challenger L2 RPC (독립적!) |
| Challenger op-node RPC | 7546 | http://localhost:7546 | Challenger op-node RPC (Follower) |
| op-challenger RPC | 6547 | http://localhost:6547 | 독립 Challenger (제3자 검증자) |
| **선택사항** ||||
| Indexer API | 8100 | http://localhost:8100 | Indexer API (선택) |
| Grafana | 3000 | http://localhost:3000 | Grafana (선택) |

## 🔍 상태 확인 명령어

```bash
# 전체 서비스 상태
docker-compose -f op-challenger/scripts/docker-compose-full.yml ps

# 특정 서비스 로그
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f sequencer-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f sequencer-challenger  # ⭐ 시퀀서용
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f challenger-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-challenger  # ⭐ 독립 검증자

# L1 블록 높이
cast block-number --rpc-url http://localhost:8545

# Sequencer L2 블록 높이
cast block-number --rpc-url http://localhost:9545

# Challenger L2 블록 높이
cast block-number --rpc-url http://localhost:9546

# Sequencer op-node 동기화 상태
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
  http://localhost:7545

# Challenger op-node 동기화 상태
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
  http://localhost:7546
```

## 🛑 시스템 중지 및 삭제

```bash
# 서비스 중지
docker-compose -f op-challenger/scripts/docker-compose-full.yml down

# 데이터까지 삭제
docker-compose -f op-challenger/scripts/docker-compose-full.yml down -v

# 이미지까지 삭제
docker-compose -f op-challenger/scripts/docker-compose-full.yml down -v --rmi all
```

## 🔧 트러블슈팅

### 1. Genesis 해시 불일치 (가장 흔한 문제)

**증상:**
```
Error: incorrect L1 genesis block hash 0x..., expected 0x...
```

**원인:**
- 기존 Docker volumes에 이전 배포의 블록체인 데이터가 남아있음
- geth는 기존 데이터가 있으면 genesis.json을 무시함

**해결:**
```bash
# 1. 완전 초기화
./op-challenger/scripts/cleanup.sh

# 2. 재배포
./op-challenger/scripts/deploy-modular.sh --dg-type 0  # 또는 --dg-type 2 (Asterisc)
```

### 2. RPC 응답 없음

```bash
# 컨테이너 상태 확인
docker ps -a

# 로그 확인
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f [service-name]

# 재시작
docker-compose -f op-challenger/scripts/docker-compose-full.yml restart [service-name]
```

### 3. 동기화 문제

```bash
# L2가 L1을 따라가지 못하는 경우
# 1. Sequencer op-node 로그 확인
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f sequencer-op-node

# 2. op-batcher 로그 확인 (배치 제출 여부)
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-batcher

# 3. L1 연결 상태 확인
curl http://localhost:8545
```

### 4. Challenger 문제

```bash
# Challenger 로그 확인
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-challenger

# Challenger가 독립적인 스택을 사용하는지 확인
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f challenger-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f challenger-l2

# 네트워크 연결 상태 확인
# 네트워크 이름 찾기 (실행 디렉토리에 따라 다를 수 있음)
docker network ls | grep l2-network

# 찾은 네트워크 이름으로 상세 정보 확인
docker network inspect <네트워크_이름>
# 예: docker network inspect scripts_l2-network
```

### 5. 포트 충돌

```bash
# 포트 사용 중인 프로세스 확인 (macOS)
lsof -i :8545

# 프로세스 종료
kill -9 [PID]

# 또는 docker-compose-full.yml의 포트 변경
```

## 📚 추가 리소스

### 공식 문서
- [Optimism Specs](https://specs.optimism.io)
- [OP Stack GitHub](https://github.com/ethereum-optimism/optimism)
- [Optimism Docs](https://docs.optimism.io)

### 프로젝트 문서

**배포 및 운영**:
- [L2 시스템 배포 가이드](../docs/l2-system-deployment-ko.md) - 전체 배포 아키텍처
- [L2 시스템 아키텍처](../docs/l2-system-architecture-ko.md) - L2 기본 구조
- [Prestate 생성 가이드](../docs/prestate-generation-guide-ko.md) ⭐ - 모든 GameType의 prestate 생성 과정

**Challenger & GameType**:
- [Challenger 시스템 아키텍처](../docs/challenger-system-architecture-ko.md) ⭐ - GameType별 상세 아키텍처
- [op-challenger 아키텍처](../docs/op-challenger-architecture-ko.md) - Challenger 내부 구조

**리서치 & 분석** ([research/](../docs/research/) 폴더):
- [RISC-V GameTypes 비교](../docs/research/risc-v-gametypes-comparison-ko.md) - GameType 2 vs 3 상세 비교
- [Optimism 코드 비교](../docs/research/asterisc-comparison-optimism-vs-tokamak-ko.md) - Asterisc 구현 비교
- [Data Availability 분석](../docs/research/challenger-data-sources-ko.md) - DA 메커니즘 분석
- [보안 분석](../docs/research/blob-pruning-risk-analysis-ko.md) - Blob pruning 위험성
- [미래 기술 (ZK, RISC-V)](../docs/research/riscv-and-zk-future-ko.md) - ZK 통합 로드맵
