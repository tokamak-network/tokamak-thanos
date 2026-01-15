#!/bin/bash
# 브릿지 프록시의 관리자를 ProxyAdmin에서 계정 #0으로 직접 탈취하는 스크립트

RPC_URL="http://localhost:8545"
PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
ACCOUNT_0="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

PROXY_ADMIN="0xcf27F781841484d5CF7e155b44954D7224caF1dD"
SYSTEM_SAFE="0x410fbB1364c5A3EE879C28673cC21EE5AA8204A5"
BRIDGE_PROXY="0x072B5bdBFC5e66B55317Ef4B4d1AE7d61592ebB2"

# 1. ProxyAdmin.changeProxyAdmin(BRIDGE_PROXY, ACCOUNT_0) 데이터 생성
DATA=$(cast calldata "changeProxyAdmin(address,address)" $BRIDGE_PROXY $ACCOUNT_0)

# 2. Safe Signature (Threshold 1)
SIG="0x000000000000000000000000$(echo $ACCOUNT_0 | cut -c 3-)000000000000000000000000000000000000000000000000000000000000000001"

# 3. execTransaction 실행 (Safe -> ProxyAdmin)
echo "🚀 Calling ProxyAdmin.changeProxyAdmin via Safe..."
cast send $SYSTEM_SAFE "execTransaction(address,uint256,bytes,uint8,uint256,uint256,uint256,address,bytes)" \
    $PROXY_ADMIN 0 $DATA 0 0 0 0 0x0000000000000000000000000000000000000000 $SIG \
    --rpc-url $RPC_URL --private-key $PRIVATE_KEY

# 4. 결과 확인 (Admin 슬롯 조회)
# admin 슬롯: 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103
echo -e "\n🔍 Verifying new Proxy Admin..."
NEW_ADMIN=$(cast storage $BRIDGE_PROXY 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103 --rpc-url $RPC_URL)
echo "✅ New Proxy Admin (Storage Slot): $NEW_ADMIN"
