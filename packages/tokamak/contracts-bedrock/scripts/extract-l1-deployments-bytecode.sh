#!/bin/bash

# L1 RPC URL and input JSON file from environment
L1_RPC_URL="${L1_RPC_URL:-}"  # e.g. https://eth-sepolia.g.alchemy.com/v2/xxx
INPUT_JSON="${L1_DEPLOY_JSON:-}"  # e.g. deployments/11155111-deploy.json
OUTPUT_DIR="bytecode/l1"

# Create output directory
mkdir -p $OUTPUT_DIR

# Check if L1_RPC_URL is provided
if [ -z "$L1_RPC_URL" ]; then
    echo "Error: Please set L1_RPC_URL environment variable"
    echo "Usage: L1_RPC_URL=https://your-rpc-url L1_DEPLOY_JSON=deployments/11155111-deploy.json ./extract-l1-deployments-bytecode.sh"
    exit 1
fi

# Check if INPUT_JSON is provided
if [ -z "$INPUT_JSON" ]; then
    echo "Error: Please set L1_DEPLOY_JSON environment variable to the deploy.json file path"
    echo "Usage: L1_RPC_URL=... L1_DEPLOY_JSON=deployments/11155111-deploy.json ./extract-l1-deployments-bytecode.sh"
    exit 1
fi

# Function to extract contract bytecode and create JSON
extract_contract_json() {
    local contract_address=$1
    local contract_name=$2

    echo "Extracting: $contract_name ($contract_address)"

    # Extract bytecode using cast
    local bytecode=$(cast code $contract_address --rpc-url $L1_RPC_URL)

    # Check if bytecode was retrieved successfully
    if [ ${#bytecode} -gt 2 ] && [ "$bytecode" != "0x" ]; then
        # Create JSON structure
        local json_content=$(cat << EOF
{
  "address": "$contract_address",
  "name": "$contract_name",
  "bytecode": "$bytecode",
  "extracted_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "network": "L1"
}
EOF
)
        # Save to JSON file
        echo "$json_content" > "$OUTPUT_DIR/${contract_name}.json"
        echo "  ✓ Success: $contract_name JSON created"
    else
        echo "  ✗ Failed: $contract_name no bytecode found"
    fi
}

# Function to extract from deploy.json format
extract_from_deploy_json() {
    local input_file=$1

    if [ ! -f "$input_file" ]; then
        echo "Error: Input file $input_file not found"
        return 1
    fi

    echo "Extracting from deploy.json file: $input_file"

    # Check if input file is a deploy.json format (object with contract names as keys)
    local file_type=$(jq -r 'type' "$input_file" 2>/dev/null)

    if [ "$file_type" = "object" ]; then
        # Get all contract names (keys) from the JSON
        local contract_names=$(jq -r 'keys[]' "$input_file" 2>/dev/null)

        if [ -z "$contract_names" ]; then
            echo "  ✗ Error: No contracts found in $input_file"
            return 1
        fi

        echo "Found contracts:"
        echo "$contract_names" | while read -r contract_name; do
            # Get address for this contract
            local contract_address=$(jq -r ".$contract_name" "$input_file" 2>/dev/null)

            if [ -n "$contract_address" ] && [ "$contract_address" != "null" ]; then
                extract_contract_json "$contract_address" "$contract_name"
            else
                echo "  ⚠ Warning: No address found for $contract_name"
            fi
        done

    else
        echo "  ✗ Error: Invalid deploy.json format in $input_file"
        return 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: L1_RPC_URL=<L1_RPC_URL> L1_DEPLOY_JSON=<deploy.json> ./extract-l1-deployments-bytecode.sh"
    echo ""
    echo "Required environment variables:"
    echo "  L1_RPC_URL      RPC endpoint for L1 (e.g. https://eth-sepolia.g.alchemy.com/v2/xxx)"
    echo "  L1_DEPLOY_JSON  Path to deploy.json (e.g. deployments/11155111-deploy.json)"
    echo ""
    echo "Example:"
    echo "  L1_RPC_URL=https://... L1_DEPLOY_JSON=deployments/11155111-deploy.json ./extract-l1-deployments-bytecode.sh"
    echo ""
    echo "The deploy.json file must be an object: {\"ContractName\": \"0x...\", ...}"
}

# Parse command line arguments
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_usage
    exit 0
fi

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Please install jq first."
    exit 1
fi

# Check if cast is installed
if ! command -v cast &> /dev/null; then
    echo "Error: cast is required but not installed. Please install Foundry first."
    exit 1
fi

echo "Starting L1 contract onchain bytecode extraction..."

extract_from_deploy_json "$INPUT_JSON"

echo "Extraction complete! Results saved in $OUTPUT_DIR folder."
echo "Files created:"
ls -la "$OUTPUT_DIR"/*.json 2>/dev/null || echo "No JSON files created"