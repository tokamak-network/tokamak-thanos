# Sepolia 원클릭 스크립트 실행 가이드

## 1. 목적

이 문서는 `ops-bedrock/scripts/sepolia-oneclick.sh`를 사용해
`.env` 설정만으로 아래를 한 번에 실행하는 절차를 설명한다.

- Tokamak `Deploy.s.sol`로 L1 컨트랙트 배포
- `op-node genesis l2`로 `genesis-l2.json`/`rollup.json` 생성
- `docker compose up`으로 런타임(`op-node`, `op-batcher`, `op-proposer`, `op-challenger`, `l2`) 기동

## 2. 사전 조건

- 실행 위치: 리포지토리 루트 (`tokamak-thanos`)
- 필수 도구:
  - `docker` (+ `docker compose`)
  - `jq`
  - `cast` (Foundry)
  - `forge` (Foundry)

## 3. 환경파일 준비

```bash
cp ops-bedrock/scripts/sepolia-oneclick.env.example .env.sepolia.oneclick
```

최소 필수 수정값:

- `L1_RPC_URL`
- `L1_BEACON_URL`
- `DEPLOYER_PRIVATE_KEY`
- `BATCHER_PRIVATE_KEY`
- `PROPOSER_PRIVATE_KEY`
- `CHALLENGER_PRIVATE_KEY`

주요 기본값:

- `DEPLOYER_WORKDIR=.deployer`
- `TOKAMAK_CONTRACTS_DIR=packages/tokamak/contracts-bedrock`
- `TOKAMAK_DEPLOY_CONFIG_TEMPLATE=packages/tokamak/contracts-bedrock/deploy-config/thanos-sepolia.json`
- `TOKAMAK_OP_NODE_RUNNER=docker`

## 4. Dry-run

온체인 트랜잭션 없이 실행 흐름만 확인한다.

```bash
./ops-bedrock/scripts/sepolia-oneclick.sh --env-file .env.sepolia.oneclick --dry-run
```

확인 포인트:

- `forge script scripts/Deploy.s.sol:Deploy` 명령 출력
- `op-node genesis l2` 명령 출력
- runtime env 파일 출력 경로 확인 (`RUNTIME_ENV_OUT`)

## 5. 실제 실행

```bash
./ops-bedrock/scripts/sepolia-oneclick.sh --env-file .env.sepolia.oneclick
```

실행 결과(기본값):

- L1 배포 결과: `.deployer/l1-deployments.json`
- Rollup 설정: `.deployer/rollup.json`
- L2 genesis: `.deployer/genesis-l2.json`
- 런타임 env: `ops-bedrock/.env.sepolia.thanos.generated`

## 6. 실행 후 확인

컨테이너 상태:

```bash
docker compose -f ops-bedrock/docker-compose.sepolia.thanos.yml --env-file ops-bedrock/.env.sepolia.thanos.generated ps
```

핵심 로그:

```bash
docker compose -f ops-bedrock/docker-compose.sepolia.thanos.yml --env-file ops-bedrock/.env.sepolia.thanos.generated logs -f op-node op-batcher op-proposer op-challenger
```

## 7. 부분 실행 옵션

### 7.1 L1 배포 스킵

기존 `l1-deployments.json`을 재사용할 때:

```bash
TOKAMAK_SKIP_L1_DEPLOY=true
```

### 7.2 Genesis 생성 스킵

기존 `genesis-l2.json`/`rollup.json`을 재사용할 때:

```bash
TOKAMAK_SKIP_GENESIS=true
```

### 7.3 런타임만 끄고 다시 올리기

```bash
RUNTIME_DOWN_FIRST=true
```

### 7.4 런타임 기동 생략

```bash
RUNTIME_UP=false
```

## 8. 자주 발생하는 문제

### 8.1 `Missing required env var`

- 원인: 필수 변수가 `.env`에 누락
- 조치: 3번 필수 키 6개를 모두 설정

### 8.2 `docker compose is not available`

- 원인: Compose 플러그인 미설치
- 조치: `docker compose version` 확인 후 설치

### 8.3 `Invalid ... address`

- 원인: private key 포맷 오류
- 조치: `0x` + 64 hex 형식 사용

### 8.4 `Incorrect Usage: flag provided but not defined: -l1.fork-public-network`

- 원인: `op-node`에 미지원 플래그 주입
- 조치: 해당 플래그 제거, 필요 시 `--l1.beacon.ignore` 사용

### 8.5 `missing bin: ./cannon/bin/cannon`

- 원인: `op-challenger` cannon 바이너리 경로 불일치
- 조치: runtime env에서 `OP_CHALLENGER_CANNON_BIN`을 이미지 내부 경로로 고정

## 9. 참고 파일

- 실행 스크립트: `ops-bedrock/scripts/sepolia-oneclick.sh`
- env 템플릿: `ops-bedrock/scripts/sepolia-oneclick.env.example`
- 런타임 compose: `ops-bedrock/docker-compose.sepolia.thanos.yml`
