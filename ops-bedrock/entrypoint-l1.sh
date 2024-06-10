#!/bin/bash
set -exu

VERBOSITY=${GETH_VERBOSITY:-3}
GETH_DATA_DIR=/db
GETH_CHAINDATA_DIR="$GETH_DATA_DIR/geth/chaindata"
GETH_KEYSTORE_DIR="$GETH_DATA_DIR/keystore"
GENESIS_FILE_PATH="${GENESIS_FILE_PATH:-/genesis.json}"
CHAIN_ID=$(cat "$GENESIS_FILE_PATH" | jq -r .config.chainId)
BLOCK_SIGNER_PRIVATE_KEY="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
BLOCK_SIGNER_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
RPC_PORT="${RPC_PORT:-8545}"
WS_PORT="${WS_PORT:-8546}"
EL_BOOTNODE_ENODE="enode://51ea9bb34d31efc3491a842ed13b8cab70e753af108526b57916d716978b380ed713f4336a80cdb85ec2a115d5a8c0ae9f3247bed3c84d3cb025c6bab311062c@127.0.0.1:0?discport=30301"

if [ ! -d "$GETH_KEYSTORE_DIR" ]; then
	echo "$GETH_KEYSTORE_DIR missing, running account import"
	echo -n "pwd" > "$GETH_DATA_DIR"/password
  echo -n "$BLOCK_SIGNER_PRIVATE_KEY" | sed 's/0x//' > "$GETH_DATA_DIR"/block-signer-key
  geth account import \
		--datadir="$GETH_DATA_DIR" \
		--password="$GETH_DATA_DIR"/password \
		"$GETH_DATA_DIR"/block-signer-key
else
	echo "$GETH_KEYSTORE_DIR exists."
fi

if [ ! -d "$GETH_CHAINDATA_DIR" ]; then
	echo "$GETH_CHAINDATA_DIR missing, running init"
	echo "Initializing genesis."
	geth --verbosity="$VERBOSITY" init \
		--datadir="$GETH_DATA_DIR" \
		"$GENESIS_FILE_PATH"
else
	echo "$GETH_CHAINDATA_DIR exists."
fi

# Warning: Archive mode is required, otherwise old trie nodes will be
# pruned within minutes of starting the devnet.

exec geth \
	--verbosity="$VERBOSITY" \
    --datadir=$GETH_DATA_DIR \
	--http \
	--http.corsdomain="*" \
	--http.vhosts="*" \
	--http.addr=0.0.0.0 \
	--http.port=$RPC_PORT \
	--http.api="engine,eth,web3,net,debug" \
	--networkid=$CHAIN_ID \
	--ws \
	--ws.addr=0.0.0.0 \
	--ws.port="$WS_PORT" \
	--ws.origins="*" \
	--ws.api=debug,eth,txpool,net,engine \
	--syncmode=full \
	--bootnodes=$EL_BOOTNODE_ENODE \
	--unlock="$BLOCK_SIGNER_ADDRESS" \
	--mine \
	--miner.etherbase="$BLOCK_SIGNER_ADDRESS" \
	--password="$GETH_DATA_DIR"/password \
	--allow-insecure-unlock \
	--rpc.allow-unprotected-txs \
	--port=7001 \
	--authrpc.addr=0.0.0.0 \
  	--authrpc.vhosts="*" \
	--authrpc.port=8551 \
	--authrpc.jwtsecret=/config/test-jwt-secret.txt \
	--gcmode=archive \
  	--metrics \
  	--metrics.addr=0.0.0.0 \
  	--metrics.port=6060
	"$@"