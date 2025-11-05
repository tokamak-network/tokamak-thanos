#!/bin/bash

# Fix AnchorStateRegistry by creating and resolving first valid dispute game
# This resolves the "0xdead cold start" problem

set -e

echo "🔧 Fixing AnchorStateRegistry cold start issue..."

# Get contract addresses from devnet descriptor
if [[ ! -f "/tmp/devnet-desc/env.json" ]]; then
    echo "❌ Devnet descriptor not found at /tmp/devnet-desc/env.json"
    echo "Make sure devnet is running and try again"
    exit 1
fi

DISPUTE_GAME_FACTORY=$(grep -o '"DisputeGameFactoryProxy": *"[^"]*"' /tmp/devnet-desc/env.json | cut -d'"' -f4)
ANCHOR_STATE_REGISTRY=$(grep -o '"AnchorStateRegistryProxy": *"[^"]*"' /tmp/devnet-desc/env.json | cut -d'"' -f4)

echo "📋 Contract addresses:"
echo "  DisputeGameFactory: $DISPUTE_GAME_FACTORY"
echo "  AnchorStateRegistry: $ANCHOR_STATE_REGISTRY"

# Get dynamic L1 RPC port from Kurtosis (actual running ports)
L1_PORT=$(kurtosis enclave inspect simple-devnet | grep "el-1-geth-teku" -A10 | grep "rpc: 8545" | grep -o '127.0.0.1:[0-9]*' | cut -d':' -f2)
if [[ -z "$L1_PORT" ]]; then
    echo "❌ Failed to find L1 RPC port from Kurtosis"
    exit 1
fi
L1_RPC="http://127.0.0.1:${L1_PORT}"

echo "🌐 Using L1 RPC: $L1_RPC"

# Check current anchor state
echo "🔍 Checking current anchor state..."
CURRENT_ANCHOR=$(cast call $ANCHOR_STATE_REGISTRY "getAnchorRoot()" --rpc-url $L1_RPC)
CURRENT_ROOT=$(echo $CURRENT_ANCHOR | cut -c1-66)

if [[ "$CURRENT_ROOT" != "0xdead000000000000000000000000000000000000000000000000000000000000" ]]; then
    echo "✅ AnchorStateRegistry already has a valid anchor state: $CURRENT_ROOT"
    exit 0
fi

echo "⚠️  AnchorStateRegistry is using default 0xdead value!"
echo "🔧 Creating first valid dispute game to fix this..."

# Get L2 genesis root (safest approach for cold start)
echo "🔍 Getting L2 genesis state root..."
# Get dynamic L2 RPC port from Kurtosis (actual running ports)
L2_PORT=$(kurtosis enclave inspect simple-devnet | grep "op-el.*-op-geth" -A10 | grep "rpc:" | grep -o '127.0.0.1:[0-9]*' | cut -d':' -f2)
if [[ -z "$L2_PORT" ]]; then
    echo "❌ Failed to find L2 RPC port from Kurtosis"
    exit 1
fi
L2_RPC="http://127.0.0.1:${L2_PORT}"
echo "🌐 Using L2 RPC: $L2_RPC"

L2_GENESIS_ROOT=$(cast block 0 --rpc-url $L2_RPC -f stateRoot 2>/dev/null)

if [[ -z "$L2_GENESIS_ROOT" || "$L2_GENESIS_ROOT" == "null" ]]; then
    echo "❌ Failed to get L2 genesis root from $L2_RPC"
    echo "Please check L2 node connectivity."
    exit 1
fi

echo "  L2 Genesis Root: $L2_GENESIS_ROOT"

# Create dispute game with genesis root
GAME_TYPE=0  # CANNON
L2_BLOCK_NUM=0  # Genesis block
EXTRA_DATA=$(printf "0x%064x" $L2_BLOCK_NUM)

echo "🎮 Creating dispute game with correct genesis root..."
echo "  GameType: $GAME_TYPE (CANNON)"
echo "  RootClaim: $L2_GENESIS_ROOT"
echo "  L2BlockNum: $L2_BLOCK_NUM"
echo "  ExtraData: $EXTRA_DATA"

# Use admin wallet for game creation
ADMIN_WALLET="0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"  # First account from test mnemonic

echo "🚀 Creating dispute game..."
GAME_TX=$(cast send $DISPUTE_GAME_FACTORY "create(uint32,bytes32,bytes)" \
    $GAME_TYPE $L2_GENESIS_ROOT $EXTRA_DATA \
    --rpc-url $L1_RPC \
    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
    --json 2>/dev/null || echo "")

if [[ -z "$GAME_TX" ]]; then
    echo "❌ Failed to create dispute game. This might be expected if:"
    echo "   1. Game with same parameters already exists"
    echo "   2. Account doesn't have enough ETH for gas"
    echo "   3. Genesis root is invalid"
    echo ""
    echo "💡 You can manually create a game using op-challenger:"
    echo "   ./op-challenger/bin/op-challenger create-game \\"
    echo "     --l1-eth-rpc http://127.0.0.1:58524 \\"  
    echo "     --game-factory-address $DISPUTE_GAME_FACTORY \\"
    echo "     --output-root $L2_GENESIS_ROOT \\"
    echo "     --l2-block-num 0 \\"
    echo "     --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
    exit 1
fi

echo "✅ Dispute game created successfully!"

# Wait a moment for game to be created
sleep 2

echo "⏳ Waiting for game to resolve naturally (defender wins with correct root)..."
echo "💡 If the root is correct, challenger won't challenge it, and it will resolve as DEFENDER_WINS"
echo "   This will automatically call setAnchorState() and update the registry."
echo ""
echo "🔍 You can monitor the game status with:"
echo "   cast call \$(cast call $DISPUTE_GAME_FACTORY \"games(uint32,bytes32,bytes)\" $GAME_TYPE $L2_GENESIS_ROOT $EXTRA_DATA --rpc-url http://127.0.0.1:58524) \"status()\" --rpc-url http://127.0.0.1:58524"
echo ""
echo "✅ Fix initiated! The AnchorStateRegistry will be updated when the game resolves."