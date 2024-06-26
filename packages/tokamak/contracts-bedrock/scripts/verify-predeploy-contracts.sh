#!/usr/bin/env bash

# solidity Versions: https://etherscan.io/solcversions

RED="\e[31m"
RESET="\e[0m"
BOLD=$(tput bold)
NORMAL=$(tput sgr0)

declare -A network_list

network_list+=(
  ["thanos-sepolia-test"]=111551118080
  ["thanos-sepolia"]=111551119090
  ["thanos-sepolia-nightly"]=111551118282
  ["devnetL1"]=901
)

function help() {
  echo -e "Verify L2 Pre-deploy contracts\n"
  echo -e "${BOLD}Usage: /path/to/script/verify-predeploy-contracts.sh${NORMAL} <ARGS>\n"
  echo -e "${BOLD}Arguments:${NORMAL}\n"
  echo -e "  ${BOLD}-n,    --network      (required)${NORMAL}               Set target network"
  echo -e "  ${BOLD}-gURL, --genesis-url  (required)${NORMAL}               Set genesis file for get contracts addresses"
  echo -e "  ${BOLD}-eURL, --explorer-url (required)${NORMAL}               Block explorer URL to verify contracts(only blockscout)"
  echo -e "\n${BOLD}Network list:${NORMAL}\n"
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
      echo -e "${RED}error${RESET}: please enter an network value(${BOLD}required${NORMAL})\n"
      echo -e "${BOLD}argument:${NORMAL}\n"
      echo -e "  ${BOLD}-n, --network${NORMAL}"
      echo -e "\n${BOLD}network list:${NORMAL}\n"
      for network in ${!network_list[@]}; do
        echo -e "  ${network}"
      done
      echo -e "\nFor more information, try '${BOLD}--help${NORMAL}'"
      exit 1
    fi
    NETWORK_NAME=$2
    shift 2
    ;;
  -gURL | --genesis-url)
    if [ -z $2 ]; then
      echo -e "${RED}error${RESET}: please enter an genesis url value(${BOLD}required${NORMAL})\n"
      echo -e "${BOLD}argument:${NORMAL}\n"
      echo -e "  ${BOLD}-gURL, --genesis-url${NORMAL}"
      echo -e "\n${BOLD}genesis url:${NORMAL}\n"
      echo -e "  thanos-sepolia-test : https://tokamak-thanos.s3.ap-northeast-2.amazonaws.com/thanos-sepolia-test/genesis.json"
      echo -e "\nFor more information, try '${BOLD}--help${NORMAL}'"
      exit 1
    fi
    GENESIS_URL=$2
    shift 2
    ;;
  -eURL | --explorer-url)
    if [ -z $2 ]; then
      echo -e "${RED}error${RESET}: please enter an explorer url value(${BOLD}required${NORMAL})\n"
      echo -e "${BOLD}argument:${NORMAL}\n"
      echo -e "  ${BOLD}-eURL, --explorer-url${NORMAL}"
      echo -e "\n${BOLD}explorer url:${NORMAL}\n"
      echo -e "  thanos-sepolia-test : https://explorer.thanos-sepolia-test.tokamak.network"
      echo -e "\nFor more information, try '${BOLD}--help${NORMAL}'"
      exit 1
    fi
    VERIFIER_URL=$(echo $2/api?)
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
  echo -e "${RED}error${RESET}: please enter an network value(${BOLD}required${NORMAL})\n"
  echo -e "${BOLD}argument:${NORMAL}\n"
  echo -e "  ${BOLD}-n, --network${NORMAL}"
  echo -e "\n${BOLD}network list:${NORMAL}\n"
  for network in ${!network_list[@]}; do
    echo -e "  ${network}"
  done
  echo -e "\nFor more information, try '${BOLD}--help${NORMAL}'"
  exit 1
fi

if [ -z $GENESIS_URL ]; then
  echo -e "${RED}error${RESET}: please enter an genesis url value(${BOLD}required${NORMAL})\n"
  echo -e "${BOLD}argument:${NORMAL}\n"
  echo -e "  ${BOLD}-gURL, --genesis-url${NORMAL}"
  echo -e "\n${BOLD}genesis url:${NORMAL}\n"
  echo -e "  thanos-sepolia-test : https://tokamak-thanos.s3.ap-northeast-2.amazonaws.com/thanos-sepolia-test/genesis.json"
  echo -e "\nFor more information, try '${BOLD}--help${NORMAL}'"
  exit 1
fi

if [ -z $VERIFIER_URL ]; then
  echo -e "${RED}error${RESET}: please enter an explorer url value(${BOLD}required${NORMAL})\n"
  echo -e "${BOLD}argument:${NORMAL}\n"
  echo -e "  ${BOLD}-eURL, --explorer-url${NORMAL}"
  echo -e "\n${BOLD}explorer url:${NORMAL}\n"
  echo -e "  thanos-sepolia-test : https://explorer.thanos-sepolia-test.tokamak.network"
  echo -e "\nFor more information, try '${BOLD}--help${NORMAL}'"
  exit 1
fi

# map array for contracts
# key(string) : proxy contract address
# value(string) : implementation contract address
declare -A contracts

contracts+=(
  ["deaddeaddeaddeaddeaddeaddeaddeaddead0000"]="" # LegacyERC20NativeToken
  ["4200000000000000000000000000000000000001"]="" # Proxy
  ["4200000000000000000000000000000000000016"]="" # L2ToL1MessagePasser
  ["4200000000000000000000000000000000000002"]="" # DeployerWhitelist
  ["4200000000000000000000000000000000000006"]="" # WNativeToken : No Proxy
  ["4200000000000000000000000000000000000007"]="" # L2CrossDomainMessenger
  ["4200000000000000000000000000000000000010"]="" # L2StandardBridge
  ["4200000000000000000000000000000000000011"]="" # SequencerFeeVault
  ["4200000000000000000000000000000000000012"]="" # OptimismMintableERC20Factory
  ["4200000000000000000000000000000000000013"]="" # L1BlockNumber
  ["420000000000000000000000000000000000000f"]="" # GasPriceOracle
  ["4200000000000000000000000000000000000015"]="" # L1Block
  ["4200000000000000000000000000000000000000"]="" # LegacyMessagePasser
  ["4200000000000000000000000000000000000014"]="" # L2ERC721Bridge
  ["4200000000000000000000000000000000000017"]="" # OptimismMintableERC721Factory
  ["4200000000000000000000000000000000000018"]="" # ProxyAdmin
  ["4200000000000000000000000000000000000019"]="" # BaseFeeVault
  ["420000000000000000000000000000000000001a"]="" # L1FeeVault
  ["4200000000000000000000000000000000000020"]="" # SchemaRegistry
  ["4200000000000000000000000000000000000021"]="" # EAS
  ["4200000000000000000000000000000000000486"]="" # ETH (TOKAMAK)
  ["4200000000000000000000000000000000000775"]="" # L2UsdcBridge (TOKAMAK)
  ["4200000000000000000000000000000000000776"]="" # SignatureChecker : No Proxy
  ["4200000000000000000000000000000000000777"]="" # MasterMinter : No Proxy
  ["4200000000000000000000000000000000000778"]="" # FiatTokenV2_2 (TOKAMAK)
)

IMPLEMENTATION_SLOT="0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
IMPLEMENTATIONSLOTFORZEPPLIN="0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3"

echo "Getting contracts addresses ImplementationSlot …"
for address in ${!contracts[@]}; do
  QUERY=.alloc.\"$address\".storage.\"$IMPLEMENTATION_SLOT\"
  contracts[${address}]=$(echo $(curl -fsSL $GENESIS_URL | jq $QUERY) | cut -c 28-67)

  if [ -z ${contracts[${address}]} ]; then
    contracts[${address}]=${address}
  fi
done

echo "Getting contracts addresses for ImplementationSlotForZepplin ...."
contract="4200000000000000000000000000000000000778" # FiatTokenV2_2
QUERY=.alloc.\"$contract\".storage.\"$IMPLEMENTATIONSLOTFORZEPPLIN\"
contracts[${contract}]=$(echo $(curl -fsSL $GENESIS_URL | jq $QUERY) | cut -c 28-67)

echo "Successfully getting contracts addresses …."

# Path of contracts
BASE_PATH=$(cd $(dirname $0)/.. && pwd -P)
LEGACY_ERC20_NATIVE_TOKEN_PATH=${BASE_PATH}/src/legacy/LegacyERC20NativeToken.sol
PROXY_PATH=${BASE_PATH}/src/universal/Proxy.sol
L2_TO_L1_MESSAGE_PASSER_PATH=${BASE_PATH}/src/L2/L2ToL1MessagePasser.sol
DEPLOYER_WHITE_LIST_PATH=${BASE_PATH}/src/legacy/DeployerWhitelist.sol
WNATIVE_TOKEN_PATH=${BASE_PATH}/src/vendor/WNativeToken.sol
L2_CROSS_DOMAIN_MESSENGER_PATH=${BASE_PATH}/src/L2/L2CrossDomainMessenger.sol
L2_STANDARD_BRIDGE_PATH=${BASE_PATH}/src/L2/L2StandardBridge.sol
SEQUENCER_FEE_VAULT_PATH=${BASE_PATH}/src/L2/SequencerFeeVault.sol
OPTIMISM_MINTABLE_ERC20_FACTORY_PATH=${BASE_PATH}/src/universal/OptimismMintableERC20Factory.sol
L1_BLOCK_NUMBER_PATH=${BASE_PATH}/src/legacy/L1BlockNumber.sol
GAS_PRICE_ORACLE_PATH=${BASE_PATH}/src/L2/GasPriceOracle.sol
L1_BLOCK_PATH=${BASE_PATH}/src/L2/L1Block.sol
LEGACY_MESSAGE_PASSER_PATH=${BASE_PATH}/src/legacy/LegacyMessagePasser.sol
L2_ERC721_BRIDGE_PATH=${BASE_PATH}/src/L2/L2ERC721Bridge.sol
OPTIMISM_MINTABLE_ERC721_FACTORY_PATH=${BASE_PATH}/src/universal/OptimismMintableERC721Factory.sol
PROXY_ADMIN_PATH=${BASE_PATH}/src/universal/ProxyAdmin.sol
BASE_FEE_VAULT_PATH=${BASE_PATH}/src/L2/BaseFeeVault.sol
L1_FEE_VAULT_PATH=${BASE_PATH}/src/L2/L1FeeVault.sol
SCHEMA_REGISTRY_PATH=${BASE_PATH}/src/EAS/SchemaRegistry.sol
EAS_PATH=${BASE_PATH}/src/EAS/EAS.sol
ETH_PATH=${BASE_PATH}/src/L2/ETH.sol
L2_USDC_BRIDGE_PROXY_PATH=${BASE_PATH}/src/tokamak-contracts/USDC/L2/tokamak-UsdcBridge/L2UsdcBridgeProxy.sol
L2_USDC_BRIDGE_PATH=${BASE_PATH}/src/tokamak-contracts/USDC/L2/tokamak-UsdcBridge/L2UsdcBridge.sol
SIGNATURECHECKER_PATH=${BASE_PATH}/src/tokamak-contracts/USDC/L2/tokamak-USDC/util/SignatureChecker.sol
MASTERMINTER_PATH=${BASE_PATH}/src/tokamak-contracts/USDC/L2/tokamak-USDC/minting/MasterMinter.sol
FIATTOKENPROXY_PATH=${BASE_PATH}/src/tokamak-contracts/USDC/L2/tokamak-USDC/v1/FiatTokenProxy.sol
FIATTOKENV2_2_PATH=${BASE_PATH}/src/tokamak-contracts/USDC/L2/tokamak-USDC/v2/FiatTokenV2_2.sol

function run() {
  verify_LegacyERC20NativeToken
  verify_proxy
  verify_L2ToL1MessagePasser
  verify_DeployerWhitelist
  verify_WNativeToken
  verify_L2CrossDomainMessenger
  verify_L2StandardBridge
  verify_SequencerFeeVault
  verify_OptimismMintableERC20Factory
  verify_L1BlockNumber
  verify_GasPriceOracle
  verify_L1Block
  verify_LegacyMessagePasser
  verify_L2ERC721Bridge
  verify_OptimismMintableERC721Factory
  verify_ProxyAdmin
  verify_BaseFeeVault
  verify_L1FeeVault
  verify_SchemaRegistry
  verify_EAS
  verify_ETH
  verify_L2UsdcBridgeProxy
  verify_L2UsdcBridge
  verify_SignatureChecker
  verify_MasterMinter
  verify_FiatTokenProxy
  verify_FiatTokenV2_2
}

# Verify contract
# $1 : compiler version
# $2 : constructor args
# $3 : contract address
# $4 : contract path
function verify_contract() {
  forge verify-contract \
    --chain-id $CHAIN_ID \
    --verifier blockscout \
    --verifier-url $VERIFIER_URL \
    --root $BASE_PATH \
    --compiler-version $1 \
    $([[ -n $2 ]] && echo "--constructor-args $2") \
    $3 \
    $4
}

function verify_LegacyERC20NativeToken() {
  CONTRACT_ADDR=${contracts["deaddeaddeaddeaddeaddeaddeaddeaddead0000"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${LEGACY_ERC20_NATIVE_TOKEN_PATH}:LegacyERC20NativeToken"
}

function verify_proxy() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0x4200000000000000000000000000000000000001)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000001"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${PROXY_PATH}:Proxy"
}

function verify_L2ToL1MessagePasser() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000016"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${L2_TO_L1_MESSAGE_PASSER_PATH}:L2ToL1MessagePasser"
}

function verify_DeployerWhitelist() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000002"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${DEPLOYER_WHITE_LIST_PATH}:DeployerWhitelist"
}

function verify_WNativeToken() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000006"]}
  COMPILER_VERSION=v0.5.17+commit.d19bba13
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${WNATIVE_TOKEN_PATH}:WNativeToken"
}

function verify_L2CrossDomainMessenger() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0x4200000000000000000000000000000000000007)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000007"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${L2_CROSS_DOMAIN_MESSENGER_PATH}:L2CrossDomainMessenger"
}

function verify_L2StandardBridge() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0x4200000000000000000000000000000000000010)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000010"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${L2_STANDARD_BRIDGE_PATH}:L2StandardBridge"
}

function verify_SequencerFeeVault() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,uint256,uint8)" 0x4200000000000000000000000000000000000011 100 0)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000011"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${SEQUENCER_FEE_VAULT_PATH}:SequencerFeeVault"
}

function verify_OptimismMintableERC20Factory() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000012"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${OPTIMISM_MINTABLE_ERC20_FACTORY_PATH}:OptimismMintableERC20Factory"
}

function verify_L1BlockNumber() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000013"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${L1_BLOCK_NUMBER_PATH}:L1BlockNumber"
}

function verify_GasPriceOracle() {
  CONTRACT_ADDR=${contracts["420000000000000000000000000000000000000f"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${GAS_PRICE_ORACLE_PATH}:GasPriceOracle"
}

function verify_L1Block() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000015"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${L1_BLOCK_PATH}:L1Block"
}

function verify_LegacyMessagePasser() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000000"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${LEGACY_MESSAGE_PASSER_PATH}:LegacyMessagePasser"
}

function verify_L2ERC721Bridge() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0xc38696FB1509776e2b97eE624410313267bf6Fa6)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000014"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${L2_ERC721_BRIDGE_PATH}:L2ERC721Bridge"
}

function verify_OptimismMintableERC721Factory() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,uint256)" 0x4200000000000000000000000000000000000014 11155111)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000017"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${OPTIMISM_MINTABLE_ERC721_FACTORY_PATH}:OptimismMintableERC721Factory"
}

function verify_ProxyAdmin() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0x4200000000000000000000000000000000000018)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000018"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${PROXY_ADMIN_PATH}:ProxyAdmin"
}

function verify_BaseFeeVault() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,uint256,uint8)" 0x4200000000000000000000000000000000000019 1 1)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000019"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${BASE_FEE_VAULT_PATH}:BaseFeeVault"
}

function verify_L1FeeVault() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,uint256,uint8)" 420000000000000000000000000000000000001a 1 1)
  CONTRACT_ADDR=${contracts["420000000000000000000000000000000000001a"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${L1_FEE_VAULT_PATH}:L1FeeVault"
}

function verify_SchemaRegistry() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000020"]}
  COMPILER_VERSION=v0.8.19+commit.7dd6d404
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${SCHEMA_REGISTRY_PATH}:SchemaRegistry"
}

function verify_EAS() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor()")
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000021"]}
  COMPILER_VERSION=v0.8.19+commit.7dd6d404
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${EAS_PATH}:EAS"
}

function verify_ETH() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor()")
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000486"]}
  COMPILER_VERSION=v0.8.15+commit.e14f2714
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${ETH_PATH}:ETH"
}

function verify_L2UsdcBridgeProxy() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address,address,bytes)" 0xC0D3c0D3c0D3C0D3c0D3c0d3C0D3c0d3c0D30775 0xa0Ee7A142d267C1f36714E4a8F75612F20a79720 0x)
  CONTRACT_ADDR="4200000000000000000000000000000000000775"
  COMPILER_VERSION=v0.8.20+commit.a1b79de6
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${L2_USDC_BRIDGE_PROXY_PATH}:L2UsdcBridgeProxy"
}

function verify_L2UsdcBridge() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000775"]}
  COMPILER_VERSION=v0.8.20+commit.a1b79de6
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${L2_USDC_BRIDGE_PATH}:L2UsdcBridge"
}

function verify_SignatureChecker() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000776"]}
  COMPILER_VERSION=v0.6.12+commit.27d51765
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${SIGNATURECHECKER_PATH}:SignatureChecker"
}

function verify_MasterMinter() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0x4200000000000000000000000000000000000778)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000777"]}
  COMPILER_VERSION=v0.6.12+commit.27d51765
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${MASTERMINTER_PATH}:MasterMinter"
}

function verify_FiatTokenProxy() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0xc0D3c0d3C0D3C0d3C0D3C0D3C0D3c0D3C0D30778)
  CONTRACT_ADDR="4200000000000000000000000000000000000778"
  COMPILER_VERSION=v0.6.12+commit.27d51765
  verify_contract $COMPILER_VERSION $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${FIATTOKENPROXY_PATH}:FiatTokenProxy"
}

function verify_FiatTokenV2_2() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000778"]}
  COMPILER_VERSION=v0.6.12+commit.27d51765
  verify_contract $COMPILER_VERSION "" $CONTRACT_ADDR "${FIATTOKENV2_2_PATH}:FiatTokenV2_2"
}

run
