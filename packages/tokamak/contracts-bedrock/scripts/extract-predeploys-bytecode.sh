#!/bin/bash

# L2 RPC URL from environment
L2_RPC_URL="${L2_RPC_URL:-}"
OUTPUT_DIR="bytecode/l2"

# Create output directory
mkdir -p $OUTPUT_DIR

# Check if L2_RPC_URL is provided
if [ -z "$L2_RPC_URL" ]; then
    echo "Error: Please set L2_RPC_URL environment variable"
    echo "Usage: L2_RPC_URL=https://your-l2-rpc ./extract-predeploys-bytecode.sh"
    exit 1
fi

# Predeploy address-name pairs (from Predeploys.sol)
predeploys=(
  "0x4200000000000000000000000000000000000000 LegacyMessagePasser"
  "0x4200000000000000000000000000000000000001 L1MessageSender"
  "0x4200000000000000000000000000000000000002 DeployerWhitelist"
  "0x4200000000000000000000000000000000000006 WETH"
  "0x4200000000000000000000000000000000000007 L2CrossDomainMessenger"
  "0x420000000000000000000000000000000000000F GasPriceOracle"
  "0x4200000000000000000000000000000000000010 L2StandardBridge"
  "0x4200000000000000000000000000000000000011 SequencerFeeVault"
  "0x4200000000000000000000000000000000000012 OptimismMintableERC20Factory"
  "0x4200000000000000000000000000000000000013 L1BlockNumber"
  "0x4200000000000000000000000000000000000014 L2ERC721Bridge"
  "0x4200000000000000000000000000000000000015 L1Block"
  "0x4200000000000000000000000000000000000016 L2ToL1MessagePasser"
  "0x4200000000000000000000000000000000000017 OptimismMintableERC721Factory"
  "0x4200000000000000000000000000000000000018 ProxyAdmin"
  "0x4200000000000000000000000000000000000019 BaseFeeVault"
  "0x420000000000000000000000000000000000001A L1FeeVault"
  "0x4200000000000000000000000000000000000020 SchemaRegistry"
  "0x4200000000000000000000000000000000000021 EAS"
  "0x4200000000000000000000000000000000000486 ETH"
  "0x4200000000000000000000000000000000000500 QuoterV2"
  "0x4200000000000000000000000000000000000501 SwapRouter02"
  "0x4200000000000000000000000000000000000502 UniswapV3Factory"
  "0x4200000000000000000000000000000000000503 NFTDescriptor"
  "0x4200000000000000000000000000000000000504 NonfungiblePositionManager"
  "0x4200000000000000000000000000000000000505 NonfungibleTokenPositionDescriptor"
  "0x4200000000000000000000000000000000000506 TickLens"
  "0x4200000000000000000000000000000000000507 UniswapInterfaceMulticall"
  "0x4200000000000000000000000000000000000508 UniversalRouter"
  "0x4200000000000000000000000000000000000509 UnsupportedProtocol"
  "0x4200000000000000000000000000000000000775 L2UsdcBridge"
  "0x4200000000000000000000000000000000000776 SignatureChecker"
  "0x4200000000000000000000000000000000000777 MasterMinter"
  "0x4200000000000000000000000000000000000778 FiatTokenV2_2"
  "0x4200000000000000000000000000000000000042 GovernanceToken"
  "0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000 LegacyERC20NativeToken"
  "0x4200000000000000000000000000000000000022 CrossL2Inbox"
  "0x4200000000000000000000000000000000000023 L2ToL2CrossDomainMessenger"
)

echo "Starting L2 contract onchain bytecode extraction..."

# Check if cast is installed
if ! command -v cast &> /dev/null; then
    echo "Error: cast is required but not installed. Please install Foundry first."
    exit 1
fi

for entry in "${predeploys[@]}"; do
  addr=$(echo "$entry" | awk '{print $1}')
  contract_name=$(echo "$entry" | awk '{print $2}')
  echo "Extracting: $contract_name ($addr)"

  bytecode=$(cast code $addr --rpc-url $L2_RPC_URL)

  if [ ${#bytecode} -gt 2 ] && [ "$bytecode" != "0x" ]; then
    json_content=$(cat << EOF
{
  "address": "$addr",
  "name": "$contract_name",
  "bytecode": "$bytecode",
  "extracted_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "network": "L2"
}
EOF
)
    echo "$json_content" > "$OUTPUT_DIR/${contract_name}.json"
    echo "  ✓ Success: $contract_name JSON created"
  else
    echo "  ✗ Failed: $contract_name no bytecode found"
  fi
done

echo "Extraction complete! Results saved in $OUTPUT_DIR folder."
echo "Files created:"
ls -la "$OUTPUT_DIR"/*.json 2>/dev/null || echo "No JSON files created"
