#!/bin/sh
set -exu

GENESIS_FILE_PATH="${GENESIS_FILE_PATH:-/genesis.json}"
CHAIN_ID=$(jq -r .config.chainId "$GENESIS_FILE_PATH")
GAS_LIMIT=$(jq -r .gasLimit "$GENESIS_FILE_PATH")
GAS_LIMIT_VALUE=$(echo "GAS_LIMIT" | awk '{ printf("%d", $0) }')

RPC_PORT="${RPC_PORT:-8545}"
L1_RPC="${L1_RPC}"
BLOCK_NUMBER="${BLOCK_NUMBER}"

exec anvil \
  --fork-url "$L1_RPC" \
  --fork-block-number "$BLOCK_NUMBER" \
  --host "0.0.0.0" \
  --port "$RPC_PORT" \
  --base-fee "1" \
  --block-time "12" \
  --gas-limit "$GAS_LIMIT_VALUE" \
  --chain-id "$CHAIN_ID" \
  "$@" &

LONG_LIVED_PID=$!

# Wait for the long-lived service to be ready
while ! nc -z localhost "$RPC_PORT"; do
  sleep 1
done
echo "Long-lived service is up and running."

python3 main.py --l1-rpc http://localhost:${RPC_PORT} &

wait $LONG_LIVED_PID
