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
# Format: "address name proxy_usage"
# proxy_usage: "proxy" for proxied contracts, "direct" for non-proxied contracts
# Based on Predeploys.sol notProxied() function: only WETH and GOVERNANCE_TOKEN are not proxied
# Removed legacy/deprecated and conditional contracts
predeploys=(
  "0x4200000000000000000000000000000000000000 LegacyMessagePasser proxy"
  "0x4200000000000000000000000000000000000002 DeployerWhitelist proxy"
  "0x4200000000000000000000000000000000000006 WETH direct"
  "0x4200000000000000000000000000000000000007 L2CrossDomainMessenger proxy"
  "0x420000000000000000000000000000000000000F GasPriceOracle proxy"
  "0x4200000000000000000000000000000000000010 L2StandardBridge proxy"
  "0x4200000000000000000000000000000000000011 SequencerFeeVault proxy"
  "0x4200000000000000000000000000000000000012 OptimismMintableERC20Factory proxy"
  "0x4200000000000000000000000000000000000013 L1BlockNumber proxy"
  "0x4200000000000000000000000000000000000014 L2ERC721Bridge proxy"
  "0x4200000000000000000000000000000000000015 L1Block proxy"
  "0x4200000000000000000000000000000000000016 L2ToL1MessagePasser proxy"
  "0x4200000000000000000000000000000000000017 OptimismMintableERC721Factory proxy"
  "0x4200000000000000000000000000000000000018 ProxyAdmin proxy"
  "0x4200000000000000000000000000000000000019 BaseFeeVault proxy"
  "0x420000000000000000000000000000000000001A L1FeeVault proxy"
  "0x4200000000000000000000000000000000000020 SchemaRegistry proxy"
  "0x4200000000000000000000000000000000000021 EAS proxy"
  "0x4200000000000000000000000000000000000042 GovernanceToken direct"
  "0x4200000000000000000000000000000000000486 ETH proxy"
  "0x4200000000000000000000000000000000000500 QuoterV2 proxy"
  "0x4200000000000000000000000000000000000501 SwapRouter02 proxy"
  "0x4200000000000000000000000000000000000502 UniswapV3Factory proxy"
  "0x4200000000000000000000000000000000000503 NFTDescriptor proxy"
  "0x4200000000000000000000000000000000000504 NonfungiblePositionManager proxy"
  "0x4200000000000000000000000000000000000505 NonfungibleTokenPositionDescriptor proxy"
  "0x4200000000000000000000000000000000000506 TickLens proxy"
  "0x4200000000000000000000000000000000000507 UniswapInterfaceMulticall proxy"
  "0x4200000000000000000000000000000000000508 UniversalRouter proxy"
  "0x4200000000000000000000000000000000000509 UnsupportedProtocol proxy"
  "0x4200000000000000000000000000000000000775 L2UsdcBridge proxy"
  "0x4200000000000000000000000000000000000776 SignatureChecker proxy"
  "0x4200000000000000000000000000000000000777 MasterMinter proxy"
  "0x4200000000000000000000000000000000000778 FiatTokenV2_2 proxy"
  "0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000 LegacyERC20NativeToken proxy"
)

# Function to check if contract is actually a proxy by checking storage slots
is_actual_proxy() {
    local contract_address=$1
    # EIP-1967 slot
    local implementation_slot="0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
    # OpenZeppelin pattern
    local implementation_slot_alt="0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3"

    # Check EIP-1967 slot
    local implementation_hex=$(cast storage $contract_address $implementation_slot --rpc-url $L2_RPC_URL 2>/dev/null)
    if [ "$implementation_hex" != "0x0000000000000000000000000000000000000000000000000000000000000000" ] && [ -n "$implementation_hex" ]; then
        return 0  # true - is proxy
    fi

    # Check alternative slot
    local implementation_hex_alt=$(cast storage $contract_address $implementation_slot_alt --rpc-url $L2_RPC_URL 2>/dev/null)
    if [ "$implementation_hex_alt" != "0x0000000000000000000000000000000000000000000000000000000000000000" ] && [ -n "$implementation_hex_alt" ]; then
        return 0  # true - is proxy
    fi

    return 1  # false - is direct
}

# Function to get implementation address from proxy
get_implementation_address() {
    local proxy_address=$1
    local implementation_slot="0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
    local implementation_slot_alt="0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3"

    # Get the implementation address from storage slot (EIP-1967)
    local implementation_hex=$(cast storage $proxy_address $implementation_slot --rpc-url $L2_RPC_URL)

    if [ "$implementation_hex" != "0x0000000000000000000000000000000000000000000000000000000000000000" ]; then
        # Convert from hex to address format (extract last 20 bytes and add 0x)
        local implementation_addr="0x$(echo $implementation_hex | sed 's/^0x//' | tail -c 41)"
        echo $implementation_addr
    else
        # Try alternative slot (OpenZeppelin pattern)
        local implementation_hex_alt=$(cast storage $proxy_address $implementation_slot_alt --rpc-url $L2_RPC_URL)
        if [ "$implementation_hex_alt" != "0x0000000000000000000000000000000000000000000000000000000000000000" ]; then
            local implementation_addr="0x$(echo $implementation_hex_alt | sed 's/^0x//' | tail -c 41)"
            echo $implementation_addr
        else
            echo ""
        fi
    fi
}

# Function to extract contract bytecode and create JSON
extract_contract_json() {
    local contract_address=$1
    local contract_name=$2
    local proxy_type=$3

    echo "Extracting: $contract_name ($contract_address) - $proxy_type"

    # Get the bytecode
    local bytecode=$(cast code $contract_address --rpc-url $L2_RPC_URL)

    if [ ${#bytecode} -gt 2 ] && [ "$bytecode" != "0x" ]; then
        # Create JSON content
        local json_content=$(cat << EOF
{
  "address": "$contract_address",
  "name": "$contract_name",
  "bytecode": "$bytecode",
  "extracted_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "network": "L2",
  "type": "$proxy_type"
}
EOF
)
        # Save to JSON file with appropriate naming
        if [ "$proxy_type" = "proxy" ]; then
            # For proxy contracts, add "Proxy" suffix
            echo "$json_content" > "$OUTPUT_DIR/${contract_name}Proxy.json"
            echo "  ✓ Success: ${contract_name}Proxy JSON created"
        else
            # For direct contracts, use original name
            echo "$json_content" > "$OUTPUT_DIR/${contract_name}.json"
            echo "  ✓ Success: $contract_name JSON created"
        fi

        # If this is a proxy contract, also extract implementation bytecode
        if [ "$proxy_type" = "proxy" ]; then
            local implementation_addr=$(get_implementation_address $contract_address)
            if [ -n "$implementation_addr" ] && [ "$implementation_addr" != "0x0" ]; then
                echo "  Extracting implementation for: $contract_name ($implementation_addr)"
                local impl_bytecode=$(cast code $implementation_addr --rpc-url $L2_RPC_URL)

                if [ ${#impl_bytecode} -gt 2 ] && [ "$impl_bytecode" != "0x" ]; then
                    local impl_json_content=$(cat << EOF
{
  "address": "$implementation_addr",
  "name": "$contract_name",
  "bytecode": "$impl_bytecode",
  "extracted_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "network": "L2",
  "type": "implementation",
  "proxy_address": "$contract_address"
}
EOF
)
                    echo "$impl_json_content" > "$OUTPUT_DIR/${contract_name}.json"
                    echo "  ✓ Success: $contract_name JSON created (implementation)"
                else
                    echo "  ✗ Failed: $contract_name implementation no bytecode found"
                fi
            else
                echo "  ⚠ Warning: No implementation address found for $contract_name"
            fi
        fi
    else
        echo "  ✗ Failed: $contract_name no bytecode found"
    fi
}

echo "Starting L2 contract onchain bytecode extraction..."
echo "Checking actual storage slots to determine proxy/direct deployment"
echo "Total contracts to extract: ${#predeploys[@]}"

# Check if cast is installed
if ! command -v cast &> /dev/null; then
    echo "Error: cast is required but not installed. Please install Foundry first."
    exit 1
fi

for entry in "${predeploys[@]}"; do
  addr=$(echo "$entry" | awk '{print $1}')
  contract_name=$(echo "$entry" | awk '{print $2}')
  original_proxy_type=$(echo "$entry" | awk '{print $3}')

  # Check actual storage slots to determine if it's really a proxy
  if is_actual_proxy "$addr"; then
    actual_proxy_type="proxy"
    echo "  ✓ $contract_name is actually a proxy (storage check confirmed)"
  else
    actual_proxy_type="direct"
    echo "  ✓ $contract_name is actually direct deployment (storage check confirmed)"
  fi

  extract_contract_json "$addr" "$contract_name" "$actual_proxy_type"
done

echo ""
echo "Extraction complete! Results saved in $OUTPUT_DIR folder."
