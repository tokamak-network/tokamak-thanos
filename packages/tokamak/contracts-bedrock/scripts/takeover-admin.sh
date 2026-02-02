#!/bin/bash
# Script to transfer ProxyAdmin ownership to account #0 on devnet

RPC_URL="http://localhost:8545"
PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
ACCOUNT_0="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

PROXY_ADMIN="0xcf27F781841484d5CF7e155b44954D7224caF1dD"
SYSTEM_SAFE="0x410fbB1364c5A3EE879C28673cC21EE5AA8204A5"

echo "📍 Current Account: $ACCOUNT_0"
echo "🛡️  Target Safe:    $SYSTEM_SAFE"
echo "📦 Proxy Admin:   $PROXY_ADMIN"

# 1. Generate ProxyAdmin.transferOwnership(ACCOUNT_0) data
DATA=$(cast calldata "transferOwnership(address)" $ACCOUNT_0)

# 2. Generate Safe Signature (Owner 0, Threshold 1)
SIG="0x000000000000000000000000$(echo $ACCOUNT_0 | cut -c 3-)000000000000000000000000000000000000000000000000000000000000000001"

# 3. Execute execTransaction (Safe -> ProxyAdmin)
echo "🚀 Executing Safe transaction to transfer ProxyAdmin ownership..."
cast send $SYSTEM_SAFE "execTransaction(address,uint256,bytes,uint8,uint256,uint256,uint256,address,bytes)" \
    $PROXY_ADMIN 0 $DATA 0 0 0 0 0x0000000000000000000000000000000000000000 $SIG \
    --rpc-url $RPC_URL --private-key $PRIVATE_KEY

# 4. Verify result
echo -e "\n🔍 Verifying new ProxyAdmin owner..."
NEW_OWNER=$(cast call $PROXY_ADMIN "owner()(address)" --rpc-url $RPC_URL)
echo "✅ New ProxyAdmin Owner: $NEW_OWNER"
