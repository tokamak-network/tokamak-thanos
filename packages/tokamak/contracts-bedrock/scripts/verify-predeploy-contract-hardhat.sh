#!/usr/bin/env bash

# solidity version: https://etherscan.io/solcversions

RED="\e[31m"
RESET="\e[0m"
BOLD=$(tput bold)
NORMAL=$(tput sgr0)

declare -A network_list

network_list+=(
  ["thanos-sepolia-test"]=111551118080
  ["devnetL1"]=901
)

function help() {
  echo -e "L2 Pre-deployment Contract Verification\n"
  echo -e "${BOLD}Usage: /path/to/script/verify-predeploy-contracts.sh${NORMAL} <ARGS>\n"
  echo -e "${BOLD}Arguments:${NORMAL}\n"
  echo -e "  ${BOLD}-n,    --network      (required)${NORMAL}               Specify target network"
  echo -e "  ${BOLD}-gURL, --genesis-url  (required)${NORMAL}               Set genesis file URL to fetch contract addresses"
  echo -e "  ${BOLD}-eURL, --explorer-url (required)${NORMAL}               Set block explorer URL for contract verification (blockscout only)"
  echo -e "\n${BOLD}List of networks:${NORMAL}\n"
  for network in ${!network_list[@]}; do
    echo -e "  ${network}"
  done
  echo -e "\n${BOLD}Example:${NORMAL}\n"
  echo -e "  ./verify-predeploy-contracts.sh \\"
  echo -e "  --network thanos-sepolia-test \\"
  echo -e "  --genesis-url https://tokamak-thanos.s3.ap-northeast-2.amazonaws.com/thanos-sepolia-test/genesis.json \\"
  echo -e "  --explorer-url https://explorer.thanos-sepolia-test.tokamak.network"
}

while [ $# -gt 0 ]; do
  ARG=$1
  case $ARG in
  -n | --network)
    if [ -z $2 ]; then
      echo -e "${RED}Error${RESET}: Please enter a network value (${BOLD}required${NORMAL})\n"
      echo -e "${BOLD}Arguments:${NORMAL}\n"
      echo -e "  ${BOLD}-n, --network${NORMAL}"
      echo -e "\n${BOLD}List of networks:${NORMAL}\n"
      for network in ${!network_list[@]}; do
        echo -e "  ${network}"
      done
      echo -e "\nFor more information try '${BOLD}--help${NORMAL}'"
      exit 1
    fi
    NETWORK_NAME=$2
    shift 2
    ;;
  -gURL | --genesis-url)
    if [ -z $2 ]; then
      echo -e "${RED}Error${RESET}: Please enter a genesis URL value (${BOLD}required${NORMAL})\n"
      echo -e "${BOLD}Arguments:${NORMAL}\n"
      echo -e "  ${BOLD}-gURL, --genesis-url${NORMAL}"
      echo -e "\n${BOLD}Genesis URL:${NORMAL}\n"
      echo -e "  thanos-sepolia-test : https://tokamak-thanos.s3.ap-northeast-2.amazonaws.com/thanos-sepolia-test/genesis.json"
      echo -e "\nFor more information try '${BOLD}--help${NORMAL}'"
      exit 1
    fi
    GENESIS_URL=$2
    shift 2
    ;;
  -eURL | --explorer-url)
    if [ -z $2 ]; then
      echo -e "${RED}Error${RESET}: Please enter an explorer URL value (${BOLD}required${NORMAL})\n"
      echo -e "${BOLD}Arguments:${NORMAL}\n"
      echo -e "  ${BOLD}-eURL, --explorer-url${NORMAL}"
      echo -e "\n${BOLD}Explorer URL:${NORMAL}\n"
      echo -e "  thanos-sepolia-test : https://explorer.thanos-sepolia-test.tokamak.network"
      echo -e "\nFor more information try '${BOLD}--help${NORMAL}'"
      exit 1
    fi
    EXPLORER_URL=$2
    shift 2
    ;;
  -h | --help)
    help
    exit 0
    ;;
  *)
    POSITION+=($1)
    shift
    ;;
  esac
done

set -- "${POSITION[@]}"

for network in ${!network_list[@]}; do
  if [ ${network} == ${NETWORK_NAME} ]; then
    CHAIN_ID=${network_list[${network}]}
    break
  fi
done

if [ -z $CHAIN_ID ]; then
  echo -e "${RED}Error${RESET}: Please enter a network value (${BOLD}required${NORMAL})\n"
  echo -e "${BOLD}Arguments:${NORMAL}\n"
  echo -e "  ${BOLD}-n, --network${NORMAL}"
  echo -e "\n${BOLD}List of networks:${NORMAL}\n"
  for network in ${!network_list[@]}; do
    echo -e "  ${network}"
  done
  echo -e "\nFor more information try '${BOLD}--help${NORMAL}'"
  exit 1
fi

if [ -z $GENESIS_URL ]; then
  echo -e "${RED}Error${RESET}: Please enter a genesis URL value (${BOLD}required${NORMAL})\n"
  echo -e "${BOLD}Arguments:${NORMAL}\n"
  echo -e "  ${BOLD}-gURL, --genesis-url${NORMAL}"
  echo -e "\n${BOLD}Genesis URL:${NORMAL}\n"
  echo -e "  thanos-sepolia-test : https://tokamak-thanos.s3.ap-northeast-2.amazonaws.com/thanos-sepolia-test/genesis.json"
  echo -e "\nFor more information try '${BOLD}--help${NORMAL}'"
  exit 1
fi

if [ -z $EXPLORER_URL ]; then
  echo -e "${RED}Error${RESET}: Please enter an explorer URL value (${BOLD}required${NORMAL})\n"
  echo -e "${BOLD}Arguments:${NORMAL}\n"
  echo -e "  ${BOLD}-eURL, --explorer-url${NORMAL}"
  echo -e "\n${BOLD}Explorer URL:${NORMAL}\n"
  echo -e "  thanos-sepolia-test : https://explorer.thanos-sepolia-test.tokamak.network"
  echo -e "\nFor more information try '${BOLD}--help${NORMAL}'"
  exit 1
fi

# Contract address mapping
# Key (string) : Proxy contract address
# Value (string) : Implementation contract address
declare -A contracts

contracts+=(
  ["0x4200000000000000000000000000000000000501"]="" # Permit2
  ["0x4200000000000000000000000000000000000502"]="" # QuoterV2
  ["0x4200000000000000000000000000000000000503"]="" # SwapRouter02
  ["0x4200000000000000000000000000000000000504"]="" # UniswapV3Factory
  ["0x4200000000000000000000000000000000000505"]="" # NFTDescriptor
  ["0x4200000000000000000000000000000000000506"]="" # NonfungiblePositionManager
  ["0x4200000000000000000000000000000000000507"]="" # NonfungibleTokenPositionDescriptor
  ["0x4200000000000000000000000000000000000508"]="" # TickLens
  ["0x4200000000000000000000000000000000000509"]="" # UniswapInterfaceMulticall
)

IMPLEMENTATION_SLOT="0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"

echo "Fetching contract addresses..."
for address in ${!contracts[@]}; do
  QUERY=.alloc.\"$address\".storage.\"$IMPLEMENTATION_SLOT\"
  contracts[${address}]=$(echo $(curl -fsSL $GENESIS_URL | jq $QUERY) | cut -c 28-67)

  if [ -z ${contracts[${address}]} ]; then
    contracts[${address}]=${address}
  fi
done
echo "Successfully fetched contract addresses!"

# Contract path
BASE_PATH=$(cd $(dirname $0)/.. && pwd -P)


function run() {
  verify_Permit2
  verify_QuoterV2
  verify_SwapRouter02
  verify_UniswapV3Factory
  verify_NFTDescriptor
  verify_NonfungiblePositionManager
  verify_NonfungibleTokenPositionDescriptor
  verify_TickLens
  verify_UniswapInterfaceMulticall
}

# Verify contract
# $1 : Compiler version
# $2 : Constructor arguments
# $3 : Contract address
# $4 : Contract path
function verify_contract() {
  npx hardhat verify --network $NETWORK_NAME $3 $4 $([[ -n $2 ]] && echo "$2")
}

function verify_Permit2() {
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000501"]}
  COMPILER_VERSION=v0.8.17+commit.8df45f5f
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${PERMIT2_PATH}:Permit2"
}

function verify_QuoterV2() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,address)" 0x4200000000000000000000000000000000000504 0x4200000000000000000000000000000000000022)
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000502"]}
  COMPILER_VERSION=v0.7.6+commit.7338295f
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/swap-router-contracts/QuoterV2.sol/QuoterV2.sol:QuoterV2"
}

function verify_SwapRouter02() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,address,address,address)" 0x0000000000000000000000000000000000000000 0x4200000000000000000000000000000000000504 0x4200000000000000000000000000000000000506 0x4200000000000000000000000000000000000022)
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000503"]}
  COMPILER_VERSION=v0.7.6+commit.7338295f
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/swap-router-contracts/SwapRouter02.sol/SwapRouter02.sol:SwapRouter02"
}

function verify_UniswapV3Factory() {
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000504"]}
  COMPILER_VERSION=v0.7.6+commit.7338295f
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/v3-core/UniswapV3Factory.sol/UniswapV3Factory.sol:UniswapV3Factory"
}

function verify_NFTDescriptor() {
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000505"]}
  COMPILER_VERSION=v0.7.0+commit.9e61f92b
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/v3-periphery/NFTDescriptor.sol/NFTDescriptor.sol:NFTDescriptor"
}

function verify_NonfungiblePositionManager() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,address,address)" 0x4200000000000000000000000000000000000504 0x4200000000000000000000000000000000000022 0x4200000000000000000000000000000000000507)
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000506"]}
  COMPILER_VERSION=v0.7.6+commit.7338295f
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/v3-periphery/NonfungibleTokenPositionDescriptor.sol/NonfungibleTokenPositionDescriptor.sol:NonfungiblePositionManager"
}

function verify_NonfungibleTokenPositionDescriptor() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,bytes32)" 0x4200000000000000000000000000000000000022 0x54574F4E00000000000000000000000000000000000000000000000000000000)
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000507"]}
  COMPILER_VERSION=v0.7.6+commit.7338295f
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/v3-periphery/NonfungibleTokenPositionDescriptor.sol:NonfungibleTokenPositionDescriptor"
}

function verify_TickLens() {
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000508"]}
  COMPILER_VERSION=v0.5.0+commit.1d4f565a
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/v3-periphery/TickLens.sol/TickLens.sol:TickLens"
}

function verify_UniswapInterfaceMulticall() {
  CONTRACT_ADDR=${contracts["0x4200000000000000000000000000000000000509"]}
  COMPILER_VERSION=v0.7.6+commit.7338295f
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${BASE_PATH}/hardhat-artifacts/v3-periphery/UniswapInterfaceMulticall.sol/UniswapInterfaceMulticall.sol:UniswapInterfaceMulticall"
}

run
