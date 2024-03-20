#!/bin/sh
set -exu

GENESIS_FILE_PATH="${GENESIS_FILE_PATH:-/genesis.json}"
CHAIN_ID=$(jq -r .config.chainId "$GENESIS_FILE_PATH")
GAS_LIMIT=$(jq -r .gasLimit "$GENESIS_FILE_PATH")
GAS_LIMIT_VALUE=$(echo "GAS_LIMIT" | awk '{ printf("%d", $0) }')

RPC_PORT="${RPC_PORT:-8545}"
WS_PORT="${WS_PORT:-8546}"
L1_RPC="${L1_RPC:-https://eth-pokt.nodies.app}"
BLOCK_NUMBER=${BLOCK_NUMBER:-19472106}

exec anvil \
  --fork-url "$L1_RPC" \
  --fork-block-number "$BLOCK_NUMBER" \
  --host "0.0.0.0" \
  --port "$RPC_PORT" \
  --base-fee "1" \
  --gas-limit "$GAS_LIMIT_VALUE" \
  --chain-id "$CHAIN_ID" \
  "$@"

