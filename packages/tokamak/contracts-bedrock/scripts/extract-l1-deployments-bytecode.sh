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

# Function to extract from input JSON file (legacy format support)
extract_from_input_json() {
    local input_file=$1

    if [ ! -f "$input_file" ]; then
        echo "Error: Input file $input_file not found"
        return 1
    fi

    echo "Extracting from input JSON file: $input_file"

    # Check if input file is a single contract JSON or a list of contracts
    local file_type=$(jq -r 'type' "$input_file" 2>/dev/null)

    if [ "$file_type" = "object" ]; then
        # Check if it's a deploy.json format (has contract names as keys)
        local has_contract_keys=$(jq -r 'keys[]' "$input_file" 2>/dev/null | head -1)
        if [ -n "$has_contract_keys" ]; then
            # It's a deploy.json format
            extract_from_deploy_json "$input_file"
            return
        fi

        # Single contract JSON file
        local contract_address=$(jq -r '.address // empty' "$input_file" 2>/dev/null)
        local contract_name=$(jq -r '.name // empty' "$input_file" 2>/dev/null)

        if [ -z "$contract_name" ]; then
            # Use filename as contract name if name field is not present
            contract_name=$(basename "$input_file" .json)
        fi

        if [ -n "$contract_address" ] && [ "$contract_address" != "null" ]; then
            extract_contract_json "$contract_address" "$contract_name"
        else
            echo "  ✗ Error: No address found in $input_file"
        fi

    elif [ "$file_type" = "array" ]; then
        # Array of contracts JSON file
        local contract_count=$(jq '. | length' "$input_file" 2>/dev/null)
        echo "Found $contract_count contracts in array"

        for i in $(seq 0 $((contract_count - 1))); do
            local contract_address=$(jq -r ".[$i].address // empty" "$input_file" 2>/dev/null)
            local contract_name=$(jq -r ".[$i].name // empty" "$input_file" 2>/dev/null)

            if [ -z "$contract_name" ]; then
                contract_name="Contract_$i"
            fi

            if [ -n "$contract_address" ] && [ "$contract_address" != "null" ]; then
                extract_contract_json "$contract_address" "$contract_name"
            else
                echo "  ⚠ Warning: No address found at index $i"
            fi
        done

    else
        echo "  ✗ Error: Invalid JSON format in $input_file"
        return 1
    fi
}

# Function to extract from deployments directory
extract_from_deployments() {
    local deployments_dir="deployments"

    if [ ! -d "$deployments_dir" ]; then
        echo "Warning: $deployments_dir directory not found"
        return
    fi

    echo "Scanning $deployments_dir for deployment files..."

    # Find all JSON files in deployments directory
    find "$deployments_dir" -name "*.json" -type f | while read -r file; do
        # Extract contract name from filename (remove .json extension)
        local contract_name=$(basename "$file" .json)

        # Extract address from JSON file using jq
        local contract_address=$(jq -r '.address // empty' "$file" 2>/dev/null)

        if [ -n "$contract_address" ] && [ "$contract_address" != "null" ]; then
            extract_contract_json "$contract_address" "$contract_name"
        else
            echo "  ⚠ Warning: No address found in $file"
        fi
    done
}

# Function to extract from specific network deployments
extract_from_network_deployments() {
    local network=$1
    local network_dir="deployments/$network"

    if [ ! -d "$network_dir" ]; then
        echo "Warning: $network_dir directory not found"
        return
    fi

    echo "Scanning $network_dir for deployment files..."

    # Find all JSON files in network deployments directory
    find "$network_dir" -name "*.json" -type f | while read -r file; do
        # Extract contract name from filename (remove .json extension)
        local contract_name=$(basename "$file" .json)

        # Extract address from JSON file using jq
        local contract_address=$(jq -r '.address // empty' "$file" 2>/dev/null)

        if [ -n "$contract_address" ] && [ "$contract_address" != "null" ]; then
            extract_contract_json "$contract_address" "$contract_name"
        else
            echo "  ⚠ Warning: No address found in $file"
        fi
    done
}

# Function to extract from manual address list
extract_from_address_list() {
    # Array of L1 contract addresses (example - modify as needed)
    declare -A l1_contracts=(
        ["0x1234567890123456789012345678901234567890"]="ExampleContract1"
        ["0x2345678901234567890123456789012345678901"]="ExampleContract2"
        # Add more contracts as needed
    )

    echo "Extracting from manual address list..."

    for addr in "${!l1_contracts[@]}"; do
        contract_name="${l1_contracts[$addr]}"
        extract_contract_json "$addr" "$contract_name"
    done
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [INPUT_JSON_FILE]"
    echo ""
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -o, --output DIR        Set output directory (default: l1_onchain_bytecodes)"
    echo "  -n, --network NETWORK   Extract from specific network deployments"
    echo ""
    echo "Examples:"
    echo "  # Extract from deploy.json file"
    echo "  L1_RPC_URL=https://your-rpc ./$0 deployments/11155111-deploy.json"
    echo ""
    echo "  # Extract from network deployments"
    echo "  L1_RPC_URL=https://your-rpc ./$0 -n thanos-sepolia"
    echo ""
    echo "  # Extract from custom output directory"
    echo "  L1_RPC_URL=https://your-rpc ./$0 -o custom_output contracts.json"
    echo ""
    echo "  # Auto-detect from deployments directory"
    echo "  L1_RPC_URL=https://your-rpc ./$0"
    echo ""
    echo "Supported JSON formats:"
    echo "  1. Deploy.json format: {\"ContractName\": \"0x...\", ...}"
    echo "  2. Single contract: {\"address\": \"0x...\", \"name\": \"ContractName\"}"
    echo "  3. Contract array: [{\"address\": \"0x...\", \"name\": \"ContractName\"}, ...]"
}

# Parse command line arguments
INPUT_FILE=""
NETWORK=""
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -n|--network)
            NETWORK="$2"
            shift 2
            ;;
        -*)
            echo "Error: Unknown option $1"
            show_usage
            exit 1
            ;;
        *)
            INPUT_FILE="$1"
            shift
            ;;
    esac
done

# Create output directory
mkdir -p $OUTPUT_DIR

# Main execution
echo "Starting L1 contract onchain bytecode extraction..."

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

# Extract from input JSON file if provided
if [ -n "$INPUT_FILE" ]; then
    extract_from_input_json "$INPUT_FILE"
elif [ -n "$NETWORK" ]; then
    # Extract from specific network
    extract_from_network_deployments "$NETWORK"
else
    # Extract from deployments directory (if exists)
    if [ -d "deployments" ]; then
        extract_from_deployments
    fi

    # Try common network names if no deployments found
    if [ ! -f "$OUTPUT_DIR"/*.json ]; then
        for network in "mainnet" "sepolia" "goerli" "thanos-sepolia"; do
            if [ -d "deployments/$network" ]; then
                extract_from_network_deployments "$network"
                break
            fi
        done
    fi

    # If no deployments found, extract from manual list
    if [ ! -f "$OUTPUT_DIR"/*.json ]; then
        echo "No deployment files found, using manual address list..."
        extract_from_address_list
    fi
fi

echo "Extraction complete! Results saved in $OUTPUT_DIR folder."
echo "Files created:"
ls -la "$OUTPUT_DIR"/*.json 2>/dev/null || echo "No JSON files created"