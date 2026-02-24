# Sepolia L2 Deployment Runbook

## 1. 목적

이 문서는 `tokamak-thanos`에서 **Sepolia(L1 chain id: `11155111`)에 Thanos L2를 배포**하는 표준 절차를 정리한다.

- 배포 엔진: Tokamak `Deploy.s.sol`
- 권장 실행 경로: `ops-bedrock/scripts/sepolia-oneclick.sh`
- 주요 산출물: `.deployer/l1-deployments.json`, `.deployer/genesis-l2.json`, `.deployer/rollup.json`

## 2. 권장 경로 (One-click)

### 2.1 사전 준비

- 실행 위치: 리포지토리 루트 (`tokamak-thanos`)
- 필수 도구:
  - `docker` (+ `docker compose`)
  - `jq`
  - `cast` (Foundry)
  - `forge` (Foundry)

### 2.2 환경파일 생성

```bash
cp ops-bedrock/scripts/sepolia-oneclick.env.example .env.sepolia.oneclick
```

필수 입력값:

- `L1_RPC_URL`
- `L1_BEACON_URL`
- `DEPLOYER_PRIVATE_KEY`
- `BATCHER_PRIVATE_KEY`
- `PROPOSER_PRIVATE_KEY`
- `CHALLENGER_PRIVATE_KEY`

권장 확인값:

- `DEPLOYER_L1_CHAIN_ID=11155111`
- `TOKAMAK_CONTRACTS_DIR=packages/tokamak/contracts-bedrock`
- `TOKAMAK_DEPLOY_CONFIG_TEMPLATE=packages/tokamak/contracts-bedrock/deploy-config/thanos-sepolia.json`
- `DEPLOYER_WORKDIR=.deployer`

### 2.3 Dry-run

```bash
./ops-bedrock/scripts/sepolia-oneclick.sh --env-file .env.sepolia.oneclick --dry-run
```

체크 포인트:

- `forge script scripts/Deploy.s.sol:Deploy` 실행 라인이 출력되는지
- `op-node genesis l2` 실행 라인이 출력되는지
- 산출물 경로가 `.deployer/*`로 기대와 일치하는지

### 2.4 실제 실행

```bash
./ops-bedrock/scripts/sepolia-oneclick.sh --env-file .env.sepolia.oneclick
```

## 3. 실행 결과

배포 후 기본 산출물:

- L1 배포 결과: `.deployer/l1-deployments.json`
- L2 genesis: `.deployer/genesis-l2.json`
- Rollup config: `.deployer/rollup.json`
- Runtime env: `ops-bedrock/.env.sepolia.thanos.generated`

컨테이너 확인:

```bash
docker compose -f ops-bedrock/docker-compose.sepolia.thanos.yml --env-file ops-bedrock/.env.sepolia.thanos.generated ps
```

핵심 로그 확인:

```bash
docker compose -f ops-bedrock/docker-compose.sepolia.thanos.yml --env-file ops-bedrock/.env.sepolia.thanos.generated logs -f op-node op-batcher op-proposer op-challenger
```

## 4. 수동 배포 경로 (필요 시)

### 4.1 Deploy config 생성

- 템플릿: `packages/tokamak/contracts-bedrock/deploy-config/thanos-sepolia.json`
- 운영 주소(`batchSenderAddress`, `l2OutputOracleProposer`, `l2OutputOracleChallenger`, `p2pSequencerAddress`)를 실제 키 주소로 반영

### 4.2 L1 컨트랙트 배포

```bash
cd packages/tokamak/contracts-bedrock
DEPLOY_CONFIG_PATH="$REPO_ROOT/.deployer/deploy-config.generated.json" \
DEPLOYMENT_OUTFILE="$REPO_ROOT/.deployer/l1-deployments.json" \
IMPL_SALT="$(openssl rand -hex 32)" \
forge script scripts/Deploy.s.sol:Deploy \
  --rpc-url "$L1_RPC_URL" \
  --broadcast \
  --private-key "$DEPLOYER_PRIVATE_KEY" \
  --slow \
  --legacy \
  --non-interactive
```

### 4.3 Genesis/Rollup 생성

```bash
op-node genesis l2 \
  --deploy-config .deployer/deploy-config.generated.json \
  --l1-deployments .deployer/l1-deployments.json \
  --outfile.l2 .deployer/genesis-l2.json \
  --outfile.rollup .deployer/rollup.json \
  --l1-rpc "$L1_RPC_URL"
```

### 4.4 런타임 실행

```bash
docker compose -f ops-bedrock/docker-compose.sepolia.thanos.yml --env-file ops-bedrock/.env.sepolia.thanos.generated up -d
```

## 5. 자주 발생하는 오류

### 5.1 `Missing required env var`

원인:

- `.env.sepolia.oneclick` 필수 항목 누락

조치:

- 2.2의 필수 키 6개를 모두 설정

### 5.2 `Invalid ... address`

원인:

- private key 형식 오류

조치:

- `0x` + 64 hex 형식 사용

### 5.3 `Incorrect Usage: flag provided but not defined: -l1.fork-public-network`

원인:

- `op-node`에 지원하지 않는 플래그 주입

조치:

- runtime 옵션에서 `l1.fork-public-network` 주입 제거
- 필요 시 `--l1.beacon.ignore` 사용

### 5.4 `invalid challenger address for chain`

원인:

- 배포/런타임 role 주소와 private key 주소 불일치

조치:

- `.env.sepolia.oneclick`의 proposer/challenger/batcher 키를 실제 운영 키로 통일
- 이전 산출물 재사용 중이면 `.deployer`를 정리하고 재배포
