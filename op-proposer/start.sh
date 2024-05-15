./bin/op-proposer \
    --poll-interval=2s \
    --rpc.port=8560 \
    --log.level=TRACE \
    --rollup-rpc=http://localhost:8547 \
    --l2oo-address=$(cat ../packages/tokamak/contracts-bedrock/deployments/getting-started/L2OutputOracleProxy.json | jq -r .address) \
    --private-key=$GS_PROPOSER_PRIVATE_KEY \
    --allow-non-finalized=true \
    --num-confirmations=1 \
    --l1-eth-rpc=$L1_RPC_URL

