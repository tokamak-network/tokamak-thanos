# L2 시스템 배포 스크립트 (독립적인 Challenger 스택 포함)

이 디렉토리는 Optimism L2 Rollup 시스템의 완전한 배포를 위한 스크립트들을 포함합니다.

## ⭐ 중요: 독립적인 Challenger 스택

이 배포 스크립트는 **문서에서 권장하는 완전히 독립적인 Challenger 스택**을 구현합니다:

```
Sequencer 스택:                 Challenger 스택 (독립적!):
├─ sequencer-l2 (L2 geth)       ├─ challenger-l2 (별도 L2 geth!) ⭐
├─ sequencer-op-node            ├─ challenger-op-node (Follower) ⭐
├─ op-batcher                   └─ op-challenger
└─ op-proposer
```

**왜 독립적이어야 하나요?**
- Sequencer와 노드를 공유하면 Sequencer가 악의적일 때 Challenger도 속을 수 있음 ❌
- 별도 DB로 L1에서 독립적으로 L2를 재구성해야 진정한 Fault Proof 가능 ✅

## 📚 참고 문서

배포 전에 다음 문서를 반드시 읽어보시기 바랍니다:

- **[Challenger 시스템 아키텍처](../docs/challenger-system-architecture-ko.md)** ⭐ - 전체 시스템 구성 및 GameType별 아키텍처 상세 설명
- **[L2 시스템 배포 가이드](../docs/l2-system-deployment-ko.md)** - 전체 배포 아키텍처 및 상세 설명
- **[L2 시스템 아키텍처](../docs/l2-system-architecture-ko.md)** - L2 시스템의 기본 구조


## 🚀 빠른 시작

```bash
# 1. 리포지토리 클론
git clone https://github.com/tokamak-network/tokamak-thanos.git
cd tokamak-thanos
git checkout feature/challenger-gametype3

# 2. VM 바이너리 다운로드 (2-3분)
./op-challenger/scripts/pull-vm-images.sh --tag latest
# ✅ cannon, asterisc, op-program, kona-client 다운로드 완료

# 3. 초기화 및 배포
./op-challenger/scripts/cleanup.sh
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
FAULT_GAME_MAX_CLOCK_DURATION=150 ./deploy-modular.sh --dg-type 3  # 커스텀 설정
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
- Challenger: 잘못된 Output에 대해 챌린지하는 계정

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
| **Challenger 스택 ⭐** ||||
| Challenger L2 RPC | 9546 | http://localhost:9546 | Challenger L2 RPC (독립적!) |
| Challenger op-node RPC | 7546 | http://localhost:7546 | Challenger op-node RPC (Follower) |
| **선택사항** ||||
| Indexer API | 8100 | http://localhost:8100 | Indexer API (선택) |
| Grafana | 3000 | http://localhost:3000 | Grafana (선택) |

## 🔍 상태 확인 명령어

```bash
# 전체 서비스 상태
docker-compose -f op-challenger/scripts/docker-compose-full.yml ps

# 특정 서비스 로그
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f sequencer-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f challenger-op-node
docker-compose -f op-challenger/scripts/docker-compose-full.yml logs -f op-challenger

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
- [스크립트 모듈화 요약](../docs/SCRIPT-MODULARIZATION-SUMMARY-ko.md) - 모듈 구조 및 사용법
- [Prestate 생성 가이드](../docs/prestate-generation-guide-ko.md) ⭐ - 모든 GameType의 prestate 생성 과정
- [VM 이미지 빌드 및 공유](../docs/vm-image-build-and-share-guide.md) ⭐ - VM 이미지 관리 워크플로우

**GameType 및 아키텍처**:
- [Challenger 시스템 아키텍처](../docs/challenger-system-architecture-ko.md) ⭐ - GameType별 상세 아키텍처
- [RISC-V GameType 비교](../docs/risc-v-gametypes-comparison-ko.md) - GameType 2 vs 3 비교
- [게임 타입 및 VM](../docs/game-types-and-vms-ko.md) - VM 아키텍처 설명

**보안 및 분석**:
- [Blob Pruning 위험 분석](../docs/blob-pruning-risk-analysis-ko.md)
- [Challenge Game 취약점 분석](../docs/challenge-game-vulnerability-ko.md)
- [DA 시스템 분석](../docs/data-availability-analysis-ko.md)
