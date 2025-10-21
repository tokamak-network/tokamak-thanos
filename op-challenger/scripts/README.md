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

- **[L2 시스템 배포 가이드](../docs/l2-system-deployment-ko.md)** - 전체 배포 아키텍처 및 상세 설명
- **[L2 시스템 아키텍처](../docs/l2-system-architecture-ko.md)** - L2 시스템의 기본 구조


## 🚀 빠른 시작

```bash
# 1. 완전 초기화
./op-challenger/scripts/cleanup.sh

# 1. 전체 자동 배포 (처음 배포 시)
./op-challenger/scripts/deploy-full-stack.sh --mode local
# → 약 5-10분 소요

# 2. 배포 상태 확인
./op-challenger/scripts/health-check.sh

# 3. 상세 모니터링 (게임 통계, 동기화 상태 등)
./op-challenger/scripts/monitor-challenger.sh
```


**⚠️ 중요:**
- `deploy-full-stack.sh`는 `.env`와 Genesis를 자동 생성합니다
- `cleanup.sh`는 기본값으로 모든 것을 삭제합니다 (`.env`, Genesis, volumes)
- `.env`를 유지하려면: `./cleanup.sh --keep-env`

## 📋 스크립트 설명

### cleanup.sh

배포된 시스템을 정리하고 초기화합니다.

**⚠️ 중요**: 기본 동작이 `--all`로 설정되어 있어 `.env`와 Genesis 파일도 함께 삭제됩니다.

```bash
# 사용법
./cleanup.sh [옵션]

# 옵션
--all               # 모든 것 정리 (컨테이너, 볼륨, .devnet, .env) [기본값]
--all-containers    # 모든 관련 컨테이너 정리 (devnet, kurtosis 등 포함)
--keep-env          # 환경 변수 파일 (.env) 유지
--keep-genesis      # Genesis 파일 유지
--help              # 도움말 표시

# 예시

# 기본 정리 (.env와 Genesis 파일도 삭제됨)
./cleanup.sh

# .env는 유지하고 Genesis만 삭제
./cleanup.sh --keep-env

# Genesis는 유지하고 .env만 삭제
./cleanup.sh --keep-genesis

# 컨테이너와 볼륨만 삭제 (.env와 Genesis 모두 유지)
./cleanup.sh --keep-env --keep-genesis

# 완전 초기화 (모든 관련 컨테이너 포함)
./cleanup.sh --all-containers
```

**정리되는 항목:**

**기본 동작 (--all):**
- scripts-* 컨테이너 (docker-compose-full.yml)
- 관련 볼륨 및 네트워크
- .env 파일
- Genesis 파일 (genesis-*.json, rollup.json)
- .devnet의 모든 allocation 파일

**추가 옵션 (--all-containers):**
- 위 항목 +
- ops-bedrock-* 컨테이너 (devnet)
- kurtosis-* 컨테이너
- op-* 관련 모든 컨테이너
- dangling 볼륨

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

### deploy-full-stack.sh

전체 L2 시스템을 배포합니다.

**✨ 자동 기능:**
- 이전 배포의 컨테이너 자동 중지
- Genesis 파일 자동 생성 (없는 경우)
- Docker volumes 자동 정리 (Genesis 해시 일관성 보장)
- Cannon VM 및 op-program 자동 빌드
- L1 계정 자동 펀딩 (로컬 모드)

```bash
# 사용법
./deploy-full-stack.sh [옵션]

# 옵션
--env-file FILE         # 환경 변수 파일 경로 (기본: .env)
--mode MODE             # 배포 모드: local|sepolia|mainnet (기본: local)
--skip-build            # Docker 이미지 빌드 건너뛰기
--skip-l1               # L1 배포 건너뛰기 (외부 RPC 사용 시)
--with-indexer          # Indexer 스택도 함께 배포
--help                  # 도움말 표시

# 예시
./deploy-full-stack.sh --mode local
./deploy-full-stack.sh --mode sepolia --skip-l1
./deploy-full-stack.sh --mode local --with-indexer
```

**배포 프로세스:**

1. **사전 검증**: 필수 도구 확인 (docker, jq, openssl, cast)
2. **컨테이너 정리**: 이전 배포 컨테이너 자동 중지
3. **환경 변수**: 없으면 setup-env.sh 자동 실행
4. **Genesis 생성**: 없으면 make devnet-up으로 자동 생성
5. **Volumes 정리**: L1, L2 volumes 자동 삭제 (Genesis 일관성 보장)
6. **Cannon VM 빌드**: op-program, cannon, prestate 자동 빌드
7. **서비스 배포**: L1 → Sequencer → Challenger 순차 배포
8. **헬스 체크**: 배포 완료 후 자동 상태 확인

**배포되는 컴포넌트:**

**L1 스택:**
- L1 geth (로컬 모드만, 또는 외부 RPC)

**Sequencer 스택:**
- sequencer-l2 (L2 op-geth, Sequencer용)
- sequencer-op-node (Sequencer 모드)
- op-batcher (배치 제출)
- op-proposer (Output 제출)

**Challenger 스택 (독립적!) ⭐:**
- challenger-l2 (L2 op-geth, **별도 데이터베이스!**)
- challenger-op-node (Follower 모드, **독립적 검증!**)
- op-challenger (챌린지 수행)

**선택사항:**
- Indexer API
- Grafana/Prometheus

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

### monitor-challenger.sh ⭐ (신규)

Challenger의 실시간 활동, 게임 참여, 동기화 상태를 모니터링합니다.

```bash
# 사용법
./monitor-challenger.sh [mode]

# 모드
summary   # 전체 대시보드 (기본값)
sync      # 블록체인 동기화 상태만 표시
logs      # 실시간 로그 (Ctrl+C로 종료)
games     # 게임 참여 상세 정보
errors    # 에러 분석
metrics   # Prometheus 메트릭

# 예시
./monitor-challenger.sh              # 전체 대시보드
./monitor-challenger.sh sync         # 동기화 상태만
./monitor-challenger.sh logs         # 실시간 로그
./monitor-challenger.sh games        # 게임 활동만
```

**모니터링 항목:**
- ✅ Challenger 컨테이너 상태
- ✅ **시스템 설정 (게임 시간, Proposer 간격 등)** ⭐
- ✅ L1/Sequencer/Challenger 블록 높이 비교 및 동기화 차이
- ✅ L1 batch 제출 상태 (op-batcher)
- ✅ Challenger 독립성 검증 (독립적인 op-node, L2 geth 사용 확인)
- ✅ 게임 참여 통계 (감지, 진행 중, 완료)
- ✅ Challenger 액션 로그 (Attack/Defend)
- ✅ 에러 및 경고 분석
- ✅ Prometheus 메트릭 (활성화 시)

**주요 기능:**
- 색상 코딩으로 로그 가독성 향상 (에러=빨강, 경고=노랑, 액션=초록)
- Sequencer와 Challenger의 동기화 차이 확인 (정상/지연 자동 판단)
- L1 배치 제출 및 Challenger의 배치 읽기 상태 모니터링
- 독립적인 스택 사용 검증 (Sequencer와 분리 확인)
- **온체인 설정 확인 (실제 배포된 게임 설정 조회)** ⭐

## ⚙️ 게임 설정

### 기본 게임 설정 (devnetL1-template.json)

```json
{
  "faultGameMaxClockDuration": 1200,        // 20분 (초 단위)
  "faultGameClockExtension": 0,             // 클락 연장 시간
  "faultGameMaxDepth": 50,                  // 최대 게임 깊이
  "faultGameSplitDepth": 14,                // 분할 깊이
  "faultGameWithdrawalDelay": 604800,       // 인출 지연 (7일)
  "disputeGameFinalityDelaySeconds": 6,     // finality 지연 (6초)
  "faultGameAbsolutePrestate": "0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98",
  "faultGameGenesisOutputRoot": "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF"
}
```

**주요 설정:**

- **faultGameMaxClockDuration**: 1200초 (20분) - 각 팀의 최대 체스 클락 시간
  - ⚠️ 최소값: 120초 (`preimageOracleChallengePeriod` 때문)
- **faultGameWithdrawalDelay**: 604800초 (7일) - 게임 종료 후 인출 대기 시간
- **faultGameSplitDepth**: 14 - 16,384 L2 블록 커버 (약 9.1시간)
- **faultGameMaxDepth**: 50 - 최대 bisection 깊이 (687억 instructions)

**환경 변수로 빠른 테스트 환경 설정:**

```bash
FAULT_GAME_MAX_CLOCK_DURATION=60 \
FAULT_GAME_WITHDRAWAL_DELAY=3600 \
PROPOSAL_INTERVAL=30s \
./deploy-full-stack.sh --mode local
```

**설정값:**
- 게임 진행: 1분
- 인출 대기: 1시간
- 제안 간격: 30초

**⚠️ 주의:** 이 설정들은 개발/테스트 전용입니다. 프로덕션에서는 충분한 시간을 설정하세요.

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

# 2. 재배포 (volumes 자동 정리됨)
./op-challenger/scripts/deploy-full-stack.sh --mode local
```

**참고**: deploy-full-stack.sh는 자동으로 L1/L2 volumes를 정리하지만, 이전에 수동으로 생성한 volumes가 있다면 cleanup.sh를 먼저 실행하세요.

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
- [Blob Pruning 위험 분석](../docs/blob-pruning-risk-analysis-ko.md)
- [Challenge Game 취약점 분석](../docs/challenge-game-vulnerability-ko.md)
- [DA 시스템 분석](../docs/data-availability-analysis-ko.md)
