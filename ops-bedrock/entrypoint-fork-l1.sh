#!/bin/sh
set -exu

GENESIS_FILE_PATH="${GENESIS_FILE_PATH:-/genesis.json}"
CHAIN_ID=$(jq -r .config.chainId "$GENESIS_FILE_PATH" || echo "1")  # This line reads chainId from genesis file or defaults to "1" if the file is not present
RPC_PORT="${RPC_PORT:-8545}"
WS_PORT="${WS_PORT:-8546}"
L1_RPC="${L1_RPC:-https://eth.llamarpc.com}"

exec anvil \
  --fork-url "$L1_RPC" \
  --host "0.0.0.0" \
  --port "$RPC_PORT" \
  --chain-id "$CHAIN_ID" \
  "$@"

