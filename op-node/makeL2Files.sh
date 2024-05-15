go run cmd/main.go genesis l2 \
    --deploy-config ../packages/tokamak/contracts-bedrock/deploy-config/getting-started.json \
    --deployment-dir ../packages/tokamak/contracts-bedrock/deployments/getting-started/ \
    --outfile.l2 genesis.json \
    --outfile.rollup rollup.json \
    --l1-rpc $L1_RPC_URL

