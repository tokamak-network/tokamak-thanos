#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
OPS_BEDROCK_DIR="${REPO_ROOT}/ops-bedrock"

ENV_FILE="${REPO_ROOT}/.env.sepolia.oneclick"
DRY_RUN=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --env-file)
      ENV_FILE="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    *)
      echo "Unknown argument: $1"
      echo "Usage: $0 [--env-file <path>] [--dry-run]"
      exit 1
      ;;
  esac
done

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "Environment file not found: ${ENV_FILE}"
  exit 1
fi

set -a
# shellcheck disable=SC1090
source "${ENV_FILE}"
set +a

log() {
  echo "[oneclick] $*"
}

run() {
  if [[ "${DRY_RUN}" == "true" ]]; then
    echo "[dry-run] $*"
  else
    "$@"
  fi
}

run_in_dir() {
  local dir="$1"
  shift
  if [[ "${DRY_RUN}" == "true" ]]; then
    echo "[dry-run] (cd ${dir} && $*)"
  else
    (
      cd "${dir}"
      "$@"
    )
  fi
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Required command not found: $1"
    exit 1
  fi
}

require_env() {
  local name="$1"
  if [[ -z "${!name:-}" ]]; then
    echo "Missing required env var: ${name}"
    exit 1
  fi
}

resolve_path_from_repo() {
  local path="$1"
  if [[ "${path}" = /* ]]; then
    echo "${path}"
  else
    echo "${REPO_ROOT}/${path}"
  fi
}

ensure_address() {
  local name="$1"
  local value="$2"
  if [[ ! "${value}" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    echo "Invalid ${name} address: ${value}"
    exit 1
  fi
}

to_decimal_chain_id() {
  local raw="$1"
  if [[ "${raw}" =~ ^0x[0-9a-fA-F]+$ ]]; then
    cast to-dec "${raw}"
    return
  fi
  if [[ "${raw}" =~ ^[0-9]+$ ]]; then
    echo "${raw}"
    return
  fi
  echo "Invalid chain id format: ${raw}" >&2
  exit 1
}

require_cmd jq
require_cmd docker
require_cmd cast
require_cmd forge
if ! docker compose version >/dev/null 2>&1; then
  echo "docker compose is not available"
  exit 1
fi

require_env L1_RPC_URL
require_env L1_BEACON_URL
require_env DEPLOYER_PRIVATE_KEY
require_env BATCHER_PRIVATE_KEY
require_env PROPOSER_PRIVATE_KEY
require_env CHALLENGER_PRIVATE_KEY

if [[ ! -f "${OPS_BEDROCK_DIR}/docker-compose.sepolia.thanos.yml" ]]; then
  echo "Missing runtime compose file: ${OPS_BEDROCK_DIR}/docker-compose.sepolia.thanos.yml"
  exit 1
fi

DEPLOYER_L1_CHAIN_ID="${DEPLOYER_L1_CHAIN_ID:-11155111}"
DEPLOYER_WORKDIR="${DEPLOYER_WORKDIR:-.deployer}"
DEPLOYER_SKIP_APPLY="${DEPLOYER_SKIP_APPLY:-false}"
DEPLOYER_L2_CHAIN_ID="${DEPLOYER_L2_CHAIN_ID:-}"

TOKAMAK_CONTRACTS_DIR="${TOKAMAK_CONTRACTS_DIR:-packages/tokamak/contracts-bedrock}"
TOKAMAK_DEPLOY_CONFIG_TEMPLATE="${TOKAMAK_DEPLOY_CONFIG_TEMPLATE:-${TOKAMAK_CONTRACTS_DIR}/deploy-config/thanos-sepolia.json}"
TOKAMAK_DEPLOY_CONFIG_GENERATED="${TOKAMAK_DEPLOY_CONFIG_GENERATED:-${DEPLOYER_WORKDIR}/deploy-config.generated.json}"
TOKAMAK_DEPLOYMENT_OUTFILE="${TOKAMAK_DEPLOYMENT_OUTFILE:-${DEPLOYER_WORKDIR}/l1-deployments.json}"
TOKAMAK_SKIP_L1_DEPLOY="${TOKAMAK_SKIP_L1_DEPLOY:-${DEPLOYER_SKIP_APPLY}}"
TOKAMAK_DEPLOY_RESUME="${TOKAMAK_DEPLOY_RESUME:-false}"
TOKAMAK_SKIP_GENESIS="${TOKAMAK_SKIP_GENESIS:-false}"
TOKAMAK_DEPLOY_SLOW="${TOKAMAK_DEPLOY_SLOW:-true}"
TOKAMAK_DEPLOY_LEGACY="${TOKAMAK_DEPLOY_LEGACY:-true}"
TOKAMAK_DEPLOY_NON_INTERACTIVE="${TOKAMAK_DEPLOY_NON_INTERACTIVE:-true}"
TOKAMAK_DEPLOY_GAS_PRICE="${TOKAMAK_DEPLOY_GAS_PRICE:-${GAS_PRICE:-}}"
TOKAMAK_OP_NODE_RUNNER="${TOKAMAK_OP_NODE_RUNNER:-docker}"

RUNTIME_ENV_OUT="${RUNTIME_ENV_OUT:-${OPS_BEDROCK_DIR}/.env.sepolia.thanos.generated}"
RUNTIME_UP="${RUNTIME_UP:-true}"
RUNTIME_DOWN_FIRST="${RUNTIME_DOWN_FIRST:-false}"
RUNTIME_COMPOSE_FILE="${RUNTIME_COMPOSE_FILE:-${OPS_BEDROCK_DIR}/docker-compose.sepolia.thanos.yml}"
RUNTIME_GENESIS_PATH_REL="${RUNTIME_GENESIS_PATH_REL:-../.deployer/genesis-l2.json}"
RUNTIME_ROLLUP_PATH_REL="${RUNTIME_ROLLUP_PATH_REL:-../.deployer/rollup.json}"

IMAGE_TAG="${IMAGE_TAG:-latest}"
TOKAMAK_OP_NODE_IMAGE="${TOKAMAK_OP_NODE_IMAGE:-tokamaknetwork/thanos-op-node:${IMAGE_TAG}}"
if [[ -x "${REPO_ROOT}/op-node/bin/op-node" ]]; then
  TOKAMAK_OP_NODE_LOCAL_BIN="${TOKAMAK_OP_NODE_LOCAL_BIN:-${REPO_ROOT}/op-node/bin/op-node}"
else
  TOKAMAK_OP_NODE_LOCAL_BIN="${TOKAMAK_OP_NODE_LOCAL_BIN:-op-node}"
fi

DEPLOYER_WORKDIR="$(resolve_path_from_repo "${DEPLOYER_WORKDIR}")"
TOKAMAK_CONTRACTS_DIR="$(resolve_path_from_repo "${TOKAMAK_CONTRACTS_DIR}")"
TOKAMAK_DEPLOY_CONFIG_TEMPLATE="$(resolve_path_from_repo "${TOKAMAK_DEPLOY_CONFIG_TEMPLATE}")"
TOKAMAK_DEPLOY_CONFIG_GENERATED="$(resolve_path_from_repo "${TOKAMAK_DEPLOY_CONFIG_GENERATED}")"
TOKAMAK_DEPLOYMENT_OUTFILE="$(resolve_path_from_repo "${TOKAMAK_DEPLOYMENT_OUTFILE}")"
RUNTIME_COMPOSE_FILE="$(resolve_path_from_repo "${RUNTIME_COMPOSE_FILE}")"
RUNTIME_ENV_OUT="$(resolve_path_from_repo "${RUNTIME_ENV_OUT}")"

mkdir -p "${DEPLOYER_WORKDIR}"

if [[ "${DRY_RUN}" != "true" && ! -d "${TOKAMAK_CONTRACTS_DIR}" ]]; then
  echo "Tokamak contracts directory not found: ${TOKAMAK_CONTRACTS_DIR}"
  exit 1
fi
if [[ "${DRY_RUN}" != "true" && ! -f "${TOKAMAK_DEPLOY_CONFIG_TEMPLATE}" ]]; then
  echo "Deploy config template not found: ${TOKAMAK_DEPLOY_CONFIG_TEMPLATE}"
  exit 1
fi

DEPLOYER_ADDRESS="$(cast wallet address --private-key "${DEPLOYER_PRIVATE_KEY}" | tr -d ' ')"
BATCHER_ADDRESS="$(cast wallet address --private-key "${BATCHER_PRIVATE_KEY}" | tr -d ' ')"
PROPOSER_ADDRESS="$(cast wallet address --private-key "${PROPOSER_PRIVATE_KEY}" | tr -d ' ')"
CHALLENGER_ADDRESS="$(cast wallet address --private-key "${CHALLENGER_PRIVATE_KEY}" | tr -d ' ')"

ensure_address "DEPLOYER_ADDRESS" "${DEPLOYER_ADDRESS}"
ensure_address "BATCHER_ADDRESS" "${BATCHER_ADDRESS}"
ensure_address "PROPOSER_ADDRESS" "${PROPOSER_ADDRESS}"
ensure_address "CHALLENGER_ADDRESS" "${CHALLENGER_ADDRESS}"

if [[ "${DRY_RUN}" == "true" ]]; then
  L2_CHAIN_ID_DEC="${DEPLOYER_L2_CHAIN_ID:-111551119090}"
else
  if [[ -n "${DEPLOYER_L2_CHAIN_ID}" ]]; then
    L2_CHAIN_ID_DEC="$(to_decimal_chain_id "${DEPLOYER_L2_CHAIN_ID}")"
  else
    L2_CHAIN_ID_DEC="$(jq -er '.l2ChainID' "${TOKAMAK_DEPLOY_CONFIG_TEMPLATE}")"
    if [[ ! "${L2_CHAIN_ID_DEC}" =~ ^[0-9]+$ ]]; then
      echo "Invalid l2ChainID in deploy config template: ${L2_CHAIN_ID_DEC}"
      exit 1
    fi
  fi
fi

if [[ -z "${TOKAMAK_IMPL_SALT:-}" ]]; then
  if command -v openssl >/dev/null 2>&1; then
    TOKAMAK_IMPL_SALT="$(openssl rand -hex 32)"
  else
    TOKAMAK_IMPL_SALT="$(date +%s)"
  fi
fi

log "Using Tokamak contracts dir: ${TOKAMAK_CONTRACTS_DIR}"
log "Using deploy config template: ${TOKAMAK_DEPLOY_CONFIG_TEMPLATE}"
log "Generated deploy config path: ${TOKAMAK_DEPLOY_CONFIG_GENERATED}"
log "Deployment output path: ${TOKAMAK_DEPLOYMENT_OUTFILE}"

if [[ "${DRY_RUN}" == "true" ]]; then
  echo "[dry-run] would generate deploy config with proposer=${PROPOSER_ADDRESS}, challenger=${CHALLENGER_ADDRESS}, batcher=${BATCHER_ADDRESS}"
else
  jq \
    --argjson l1_chain_id "${DEPLOYER_L1_CHAIN_ID}" \
    --argjson l2_chain_id "${L2_CHAIN_ID_DEC}" \
    --arg batcher "${BATCHER_ADDRESS}" \
    --arg proposer "${PROPOSER_ADDRESS}" \
    --arg challenger "${CHALLENGER_ADDRESS}" \
    --arg sequencer "${DEPLOYER_ADDRESS}" \
    '.l1ChainID = $l1_chain_id
    | .l2ChainID = $l2_chain_id
    | .batchSenderAddress = $batcher
    | .l2OutputOracleProposer = $proposer
    | .l2OutputOracleChallenger = $challenger
    | .p2pSequencerAddress = $sequencer
    | .reuseDeployment = false' \
    "${TOKAMAK_DEPLOY_CONFIG_TEMPLATE}" > "${TOKAMAK_DEPLOY_CONFIG_GENERATED}"
fi

if [[ "${TOKAMAK_SKIP_L1_DEPLOY}" != "true" ]]; then
  log "Deploying L1 contracts with Tokamak Deploy.s.sol"
  FORGE_CMD=(forge script scripts/Deploy.s.sol:Deploy --rpc-url "${L1_RPC_URL}" --broadcast --private-key "${DEPLOYER_PRIVATE_KEY}")
  if [[ "${TOKAMAK_DEPLOY_SLOW}" == "true" ]]; then
    FORGE_CMD+=(--slow)
  fi
  if [[ "${TOKAMAK_DEPLOY_LEGACY}" == "true" ]]; then
    FORGE_CMD+=(--legacy)
  fi
  if [[ "${TOKAMAK_DEPLOY_NON_INTERACTIVE}" == "true" ]]; then
    FORGE_CMD+=(--non-interactive)
  fi
  if [[ -n "${TOKAMAK_DEPLOY_GAS_PRICE}" ]]; then
    FORGE_CMD+=(--with-gas-price "${TOKAMAK_DEPLOY_GAS_PRICE}")
  fi
  if [[ "${TOKAMAK_DEPLOY_RESUME}" == "true" ]]; then
    FORGE_CMD+=(--resume)
  fi

  run_in_dir "${TOKAMAK_CONTRACTS_DIR}" \
    env \
      DEPLOY_CONFIG_PATH="${TOKAMAK_DEPLOY_CONFIG_GENERATED}" \
      DEPLOYMENT_OUTFILE="${TOKAMAK_DEPLOYMENT_OUTFILE}" \
      IMPL_SALT="${TOKAMAK_IMPL_SALT}" \
      "${FORGE_CMD[@]}"
else
  log "Skipping L1 deployment because TOKAMAK_SKIP_L1_DEPLOY=true"
fi

if [[ "${DRY_RUN}" != "true" && ! -f "${TOKAMAK_DEPLOYMENT_OUTFILE}" ]]; then
  echo "Deployment output file not found: ${TOKAMAK_DEPLOYMENT_OUTFILE}"
  exit 1
fi

ROLLUP_PATH="${DEPLOYER_WORKDIR}/rollup.json"
GENESIS_L2_PATH="${DEPLOYER_WORKDIR}/genesis-l2.json"

if [[ "${TOKAMAK_SKIP_GENESIS}" != "true" ]]; then
  log "Generating rollup/genesis with op-node genesis l2"
  if [[ "${TOKAMAK_OP_NODE_RUNNER}" == "docker" ]]; then
    GENESIS_CMD=(
      docker run --rm
      -v "${REPO_ROOT}:${REPO_ROOT}"
      -w "${REPO_ROOT}"
      "${TOKAMAK_OP_NODE_IMAGE}"
      op-node genesis l2
      --deploy-config "${TOKAMAK_DEPLOY_CONFIG_GENERATED}"
      --l1-deployments "${TOKAMAK_DEPLOYMENT_OUTFILE}"
      --outfile.l2 "${GENESIS_L2_PATH}"
      --outfile.rollup "${ROLLUP_PATH}"
      --l1-rpc "${L1_RPC_URL}"
    )
    run "${GENESIS_CMD[@]}"
  elif [[ "${TOKAMAK_OP_NODE_RUNNER}" == "local" ]]; then
    if [[ "${TOKAMAK_OP_NODE_LOCAL_BIN}" = /* ]]; then
      if [[ "${DRY_RUN}" != "true" && ! -x "${TOKAMAK_OP_NODE_LOCAL_BIN}" ]]; then
        echo "Local op-node binary not executable: ${TOKAMAK_OP_NODE_LOCAL_BIN}"
        exit 1
      fi
    else
      require_cmd "${TOKAMAK_OP_NODE_LOCAL_BIN}"
    fi
    GENESIS_CMD=(
      "${TOKAMAK_OP_NODE_LOCAL_BIN}" genesis l2
      --deploy-config "${TOKAMAK_DEPLOY_CONFIG_GENERATED}"
      --l1-deployments "${TOKAMAK_DEPLOYMENT_OUTFILE}"
      --outfile.l2 "${GENESIS_L2_PATH}"
      --outfile.rollup "${ROLLUP_PATH}"
      --l1-rpc "${L1_RPC_URL}"
    )
    run "${GENESIS_CMD[@]}"
  else
    echo "Unsupported TOKAMAK_OP_NODE_RUNNER: ${TOKAMAK_OP_NODE_RUNNER} (expected: docker|local)"
    exit 1
  fi
else
  log "Skipping genesis generation because TOKAMAK_SKIP_GENESIS=true"
fi

if [[ "${DRY_RUN}" != "true" ]]; then
  if [[ ! -f "${ROLLUP_PATH}" ]]; then
    echo "Rollup file not found: ${ROLLUP_PATH}"
    exit 1
  fi
  if [[ ! -f "${GENESIS_L2_PATH}" ]]; then
    echo "Genesis file not found: ${GENESIS_L2_PATH}"
    exit 1
  fi
fi

if [[ "${DRY_RUN}" == "true" ]]; then
  DGF_ADDRESS="${DGF_ADDRESS:-0x0000000000000000000000000000000000000000}"
  SYSTEM_CONFIG_PROXY="${SYSTEM_CONFIG_PROXY:-0x0000000000000000000000000000000000000000}"
  OPTIMISM_PORTAL_PROXY="${OPTIMISM_PORTAL_PROXY:-0x0000000000000000000000000000000000000000}"
  L1_STANDARD_BRIDGE_PROXY="${L1_STANDARD_BRIDGE_PROXY:-0x0000000000000000000000000000000000000000}"
  L1_XDM_PROXY="${L1_XDM_PROXY:-0x0000000000000000000000000000000000000000}"
  MINTABLE_ERC20_FACTORY_PROXY="${MINTABLE_ERC20_FACTORY_PROXY:-0x0000000000000000000000000000000000000000}"
else
  DGF_ADDRESS="$(jq -er '.DisputeGameFactoryProxy' "${TOKAMAK_DEPLOYMENT_OUTFILE}")"
  SYSTEM_CONFIG_PROXY="$(jq -er '.SystemConfigProxy' "${TOKAMAK_DEPLOYMENT_OUTFILE}")"
  OPTIMISM_PORTAL_PROXY="$(jq -er '.OptimismPortalProxy' "${TOKAMAK_DEPLOYMENT_OUTFILE}")"
  L1_STANDARD_BRIDGE_PROXY="$(jq -er '.L1StandardBridgeProxy' "${TOKAMAK_DEPLOYMENT_OUTFILE}")"
  L1_XDM_PROXY="$(jq -er '.L1CrossDomainMessengerProxy' "${TOKAMAK_DEPLOYMENT_OUTFILE}")"
  MINTABLE_ERC20_FACTORY_PROXY="$(jq -er '.OptimismMintableERC20FactoryProxy' "${TOKAMAK_DEPLOYMENT_OUTFILE}")"
  L2_CHAIN_ID_DEC="$(jq -er '.l2_chain_id' "${ROLLUP_PATH}")"
fi

L2_IMAGE="${L2_IMAGE:-tokamaknetwork/thanos-op-geth:nightly}"

L2_HTTP_PORT="${L2_HTTP_PORT:-9545}"
L2_WS_PORT="${L2_WS_PORT:-9546}"
L2_AUTH_PORT="${L2_AUTH_PORT:-8551}"
L2_METRICS_PORT="${L2_METRICS_PORT:-8060}"
OP_NODE_RPC_PORT="${OP_NODE_RPC_PORT:-7545}"
OP_NODE_METRICS_PORT="${OP_NODE_METRICS_PORT:-7300}"
OP_NODE_PPROF_PORT="${OP_NODE_PPROF_PORT:-6060}"
OP_BATCHER_RPC_PORT="${OP_BATCHER_RPC_PORT:-8548}"
OP_BATCHER_METRICS_PORT="${OP_BATCHER_METRICS_PORT:-7301}"
OP_BATCHER_PPROF_PORT="${OP_BATCHER_PPROF_PORT:-6061}"
OP_PROPOSER_RPC_PORT="${OP_PROPOSER_RPC_PORT:-8560}"
OP_PROPOSER_METRICS_PORT="${OP_PROPOSER_METRICS_PORT:-7302}"
OP_PROPOSER_PPROF_PORT="${OP_PROPOSER_PPROF_PORT:-6062}"

SEQUENCER_L1_CONFS="${SEQUENCER_L1_CONFS:-4}"
VERIFIER_L1_CONFS="${VERIFIER_L1_CONFS:-4}"
GETH_VERBOSITY="${GETH_VERBOSITY:-3}"
BATCHER_POLL_INTERVAL="${BATCHER_POLL_INTERVAL:-6s}"
BATCHER_SUB_SAFETY_MARGIN="${BATCHER_SUB_SAFETY_MARGIN:-6}"
BATCHER_NUM_CONFIRMATIONS="${BATCHER_NUM_CONFIRMATIONS:-1}"
BATCHER_SAFE_ABORT_NONCE_TOO_LOW_COUNT="${BATCHER_SAFE_ABORT_NONCE_TOO_LOW_COUNT:-3}"
BATCHER_RESUBMISSION_TIMEOUT="${BATCHER_RESUBMISSION_TIMEOUT:-30s}"
BATCHER_MAX_CHANNEL_DURATION="${BATCHER_MAX_CHANNEL_DURATION:-300}"
BATCHER_TARGET_NUM_FRAMES="${BATCHER_TARGET_NUM_FRAMES:-1}"
BATCHER_APPROX_COMPR_RATIO="${BATCHER_APPROX_COMPR_RATIO:-0.4}"
OP_BATCHER_BATCH_TYPE="${OP_BATCHER_BATCH_TYPE:-0}"
PROPOSER_POLL_INTERVAL="${PROPOSER_POLL_INTERVAL:-12s}"
PROPOSER_ALLOW_NON_FINALIZED="${PROPOSER_ALLOW_NON_FINALIZED:-false}"
PROPOSAL_INTERVAL="${PROPOSAL_INTERVAL:-10m}"
DG_TYPE="${DG_TYPE:-0}"
CHALLENGER_GAME_TYPES="${CHALLENGER_GAME_TYPES:-cannon}"
CHALLENGER_NUM_CONFIRMATIONS="${CHALLENGER_NUM_CONFIRMATIONS:-1}"

JWT_SECRET_PATH="${JWT_SECRET_PATH:-./test-jwt-secret.txt}"
CANNON_PRESTATE_PATH="${CANNON_PRESTATE_PATH:-../op-program/bin/prestate.json}"

log "Writing runtime env file: ${RUNTIME_ENV_OUT}"
if [[ "${DRY_RUN}" == "true" ]]; then
  echo "[dry-run] would write ${RUNTIME_ENV_OUT}"
else
  mkdir -p "$(dirname "${RUNTIME_ENV_OUT}")"
  cat > "${RUNTIME_ENV_OUT}" <<EOF
L1_RPC_URL=${L1_RPC_URL}
L1_BEACON_URL=${L1_BEACON_URL}
L1_CHAIN_ID=${DEPLOYER_L1_CHAIN_ID}
L2_CHAIN_ID=${L2_CHAIN_ID_DEC}

IMAGE_TAG=${IMAGE_TAG}
L2_IMAGE=${L2_IMAGE}

L2_HTTP_PORT=${L2_HTTP_PORT}
L2_WS_PORT=${L2_WS_PORT}
L2_AUTH_PORT=${L2_AUTH_PORT}
L2_METRICS_PORT=${L2_METRICS_PORT}
OP_NODE_RPC_PORT=${OP_NODE_RPC_PORT}
OP_NODE_METRICS_PORT=${OP_NODE_METRICS_PORT}
OP_NODE_PPROF_PORT=${OP_NODE_PPROF_PORT}
OP_BATCHER_RPC_PORT=${OP_BATCHER_RPC_PORT}
OP_BATCHER_METRICS_PORT=${OP_BATCHER_METRICS_PORT}
OP_BATCHER_PPROF_PORT=${OP_BATCHER_PPROF_PORT}
OP_PROPOSER_RPC_PORT=${OP_PROPOSER_RPC_PORT}
OP_PROPOSER_METRICS_PORT=${OP_PROPOSER_METRICS_PORT}
OP_PROPOSER_PPROF_PORT=${OP_PROPOSER_PPROF_PORT}

GENESIS_L2_PATH=${RUNTIME_GENESIS_PATH_REL}
ROLLUP_CONFIG_PATH=${RUNTIME_ROLLUP_PATH_REL}
JWT_SECRET_PATH=${JWT_SECRET_PATH}
CANNON_PRESTATE_PATH=${CANNON_PRESTATE_PATH}

GETH_VERBOSITY=${GETH_VERBOSITY}
SEQUENCER_L1_CONFS=${SEQUENCER_L1_CONFS}
VERIFIER_L1_CONFS=${VERIFIER_L1_CONFS}

BATCHER_POLL_INTERVAL=${BATCHER_POLL_INTERVAL}
BATCHER_SUB_SAFETY_MARGIN=${BATCHER_SUB_SAFETY_MARGIN}
BATCHER_NUM_CONFIRMATIONS=${BATCHER_NUM_CONFIRMATIONS}
BATCHER_SAFE_ABORT_NONCE_TOO_LOW_COUNT=${BATCHER_SAFE_ABORT_NONCE_TOO_LOW_COUNT}
BATCHER_RESUBMISSION_TIMEOUT=${BATCHER_RESUBMISSION_TIMEOUT}
BATCHER_MAX_CHANNEL_DURATION=${BATCHER_MAX_CHANNEL_DURATION}
BATCHER_TARGET_NUM_FRAMES=${BATCHER_TARGET_NUM_FRAMES}
BATCHER_APPROX_COMPR_RATIO=${BATCHER_APPROX_COMPR_RATIO}
OP_BATCHER_BATCH_TYPE=${OP_BATCHER_BATCH_TYPE}

PROPOSER_POLL_INTERVAL=${PROPOSER_POLL_INTERVAL}
PROPOSER_ALLOW_NON_FINALIZED=${PROPOSER_ALLOW_NON_FINALIZED}
PROPOSAL_INTERVAL=${PROPOSAL_INTERVAL}
DG_TYPE=${DG_TYPE}

CHALLENGER_GAME_TYPES=${CHALLENGER_GAME_TYPES}
CHALLENGER_NUM_CONFIRMATIONS=${CHALLENGER_NUM_CONFIRMATIONS}

BATCHER_PRIVATE_KEY=${BATCHER_PRIVATE_KEY}
PROPOSER_PRIVATE_KEY=${PROPOSER_PRIVATE_KEY}
CHALLENGER_PRIVATE_KEY=${CHALLENGER_PRIVATE_KEY}

DGF_ADDRESS=${DGF_ADDRESS}
SYSTEM_CONFIG_PROXY=${SYSTEM_CONFIG_PROXY}
OPTIMISM_PORTAL_PROXY=${OPTIMISM_PORTAL_PROXY}
L1_STANDARD_BRIDGE_PROXY=${L1_STANDARD_BRIDGE_PROXY}
L1_XDM_PROXY=${L1_XDM_PROXY}
MINTABLE_ERC20_FACTORY_PROXY=${MINTABLE_ERC20_FACTORY_PROXY}

EXPECTED_BATCHER_ADDRESS=${BATCHER_ADDRESS}
EXPECTED_PROPOSER_ADDRESS=${PROPOSER_ADDRESS}
EXPECTED_CHALLENGER_ADDRESS=${CHALLENGER_ADDRESS}
EOF
fi

if [[ "${DRY_RUN}" != "true" ]]; then
  for addr in \
    "${DGF_ADDRESS}" \
    "${SYSTEM_CONFIG_PROXY}" \
    "${OPTIMISM_PORTAL_PROXY}" \
    "${L1_STANDARD_BRIDGE_PROXY}" \
    "${L1_XDM_PROXY}" \
    "${MINTABLE_ERC20_FACTORY_PROXY}"; do
    ensure_address "deployment address" "${addr}"
  done
fi

if [[ "${RUNTIME_UP}" == "true" ]]; then
  log "Starting runtime services with docker compose"
  if [[ "${RUNTIME_DOWN_FIRST}" == "true" ]]; then
    run docker compose -f "${RUNTIME_COMPOSE_FILE}" --env-file "${RUNTIME_ENV_OUT}" down
  fi
  run docker compose -f "${RUNTIME_COMPOSE_FILE}" --env-file "${RUNTIME_ENV_OUT}" up -d
fi

log "Completed"
log "Tokamak deploy config: ${TOKAMAK_DEPLOY_CONFIG_GENERATED}"
log "L1 deployments file: ${TOKAMAK_DEPLOYMENT_OUTFILE}"
log "Runtime env: ${RUNTIME_ENV_OUT}"
log "Rollup file: ${ROLLUP_PATH}"
log "Genesis file: ${GENESIS_L2_PATH}"
