#!/usr/bin/env bash

# Deployed contracts list
# L2ToL1MessagePasser           = "0x4200000000000000000000000000000000000016"
# DeployerWhitelist             = "0x4200000000000000000000000000000000000002"
# WTON                          = "0x4200000000000000000000000000000000000006" No Proxy
# L2CrossDomainMessenger        = "0x4200000000000000000000000000000000000007"
# L2StandardBridge              = "0x4200000000000000000000000000000000000010"
# SequencerFeeVault             = "0x4200000000000000000000000000000000000011"
# OptimismMintableERC20Factory  = "0x4200000000000000000000000000000000000012"
# L1BlockNumber                 = "0x4200000000000000000000000000000000000013"
# GasPriceOracle                = "0x420000000000000000000000000000000000000F"
# L1Block                       = "0x4200000000000000000000000000000000000015"
# GovernanceToken               = "0x4200000000000000000000000000000000000042" No Proxy
# LegacyMessagePasser           = "0x4200000000000000000000000000000000000000"
# L2ERC721Bridge                = "0x4200000000000000000000000000000000000014"
# OptimismMintableERC721Factory = "0x4200000000000000000000000000000000000017"
# ProxyAdmin                    = "0x4200000000000000000000000000000000000018"
# BaseFeeVault                  = "0x4200000000000000000000000000000000000019"
# L1FeeVault                    = "0x420000000000000000000000000000000000001a"
# SchemaRegistry                = "0x4200000000000000000000000000000000000020"
# EAS                           = "0x4200000000000000000000000000000000000021"
# WETH                          = "0x4200000000000000000000000000000000000022"


# implement slot
# ImplementationSlot = common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")

# codeNamespace
# codeNamespace = common.HexToAddress("0xc0D3C0d3C0d3C0D3c0d3C0d3c0D3C0d3c0d30000")

while [ $# -gt 0 ]; do
  OPTION=$1
  case $OPTION in
  -c | --chain-id)
    if [ -z $2 ]; then
        echo "Error: Please Enter an chain id value(required)"
        echo "chain id : "
        echo "  titan-sepolia-test : 111551115050"
        exit 1
    fi
    CHAIN_ID=$2
    shift 2
    ;;
  -gURL | --genesis-url)
    if [ -z $2 ]; then
        echo "Error: Please Enter an genesis url value(required)"
        echo "genesis url : "
        echo "  titan-sepolia-test : https://tokamak-titan-canyon.s3.ap-northeast-2.amazonaws.com/titan-sepolia-test/genesis.json"
        exit 1
    fi
    GENESIS_URL=$2
    shift 2
    ;;
  -eURL | --explorer-url)
    if [ -z $2 ]; then
        echo "Error: Please Enter an explorer url value(required)"
        echo "explorer url : "
        echo "  titan-sepolia-test : https://explorer.titan-sepolia-test.tokamak.network"
        exit 1
    fi
    VERIFIER_URL=$(echo $2/api?)
    shift 2
    ;;
  *)
    POSITION+=($1)
    shift
    ;;
  esac
done

set -- "${POSITION[@]}"

if [ -z $CHAIN_ID ]; then
  echo "Error: Please Enter an chain id value(required)"
  echo "option   : -c, --chain-id"
  echo "chain id :"
  echo "  titan-sepolia-test : 111551115050"
  exit 1
fi

if [ -z $GENESIS_URL ]; then
  echo "Error: Please Enter an genesis url value(required)"
  echo "option      : -gURL, --genesis-url"
  echo "genesis url : "
  echo "  titan-sepolia-test : https://tokamak-titan-canyon.s3.ap-northeast-2.amazonaws.com/titan-sepolia-test/genesis.json"
  exit 1
fi

if [ -z $VERIFIER_URL ]; then
  echo "Error: Please Enter an explorer url value(required)"
  echo "option       : -eURL, --explorer-url"
  echo "explorer url : "
  echo "  titan-sepolia-test : https://explorer.titan-sepolia-test.tokamak.network"
  exit 1
fi

# map array for contracts
# key(string) : contract address
# value(string) : implementation contract address
declare -A contracts

contracts+=( ["4200000000000000000000000000000000000001"]=""  # Proxy
             ["4200000000000000000000000000000000000016"]=""  # L2ToL1MessagePasser
             ["4200000000000000000000000000000000000002"]=""  # DeployerWhitelist
             ["4200000000000000000000000000000000000006"]=""  # WTON !
             ["4200000000000000000000000000000000000007"]=""  # L2CrossDomainMessenger
             ["4200000000000000000000000000000000000010"]=""  # L2StandardBridge
             ["4200000000000000000000000000000000000011"]=""  # SequencerFeeVault
             ["4200000000000000000000000000000000000012"]=""  # OptimismMintableERC20Factory
             ["4200000000000000000000000000000000000013"]=""  # L1BlockNumber
             ["420000000000000000000000000000000000000f"]=""  # GasPriceOracle
             ["4200000000000000000000000000000000000015"]=""  # L1Block
             ["4200000000000000000000000000000000000042"]=""  # GovernanceToken !
             ["4200000000000000000000000000000000000000"]=""  # LegacyMessagePasser
             ["4200000000000000000000000000000000000014"]=""  # L2ERC721Bridge
             ["4200000000000000000000000000000000000017"]=""  # OptimismMintableERC721Factory
             ["4200000000000000000000000000000000000018"]=""  # ProxyAdmin
             ["4200000000000000000000000000000000000019"]=""  # BaseFeeVault
             ["420000000000000000000000000000000000001a"]=""  # L1FeeVault
             ["4200000000000000000000000000000000000020"]=""  # SchemaRegistry
             ["4200000000000000000000000000000000000021"]=""  # EAS
             ["4200000000000000000000000000000000000022"]=""  # WETH ! is proxy or not?
)

IMPLEMENTATION_SLOT="0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"

echo "Getting contracts addresses..."
for address in ${!contracts[@]}; do
  QUERY=.alloc.\"$address\".storage.\"$IMPLEMENTATION_SLOT\"
  contracts[${address}]=$(echo $(curl -fsSL $GENESIS_URL | jq $QUERY) | cut -c 28-67)

  if [ -z ${contracts[${address}]} ]; then
    contracts[${address}]=${address}
  fi
done
echo "Successfully getting contracts addresses!"

# Path of contracts
BASE_PATH=$(cd $(dirname $0)/.. && pwd -P)
PROXY_PATH=${BASE_PATH}/src/universal/Proxy.sol
L2_TO_L1_MESSAGE_PASSER_PATH=${BASE_PATH}/src/L2/L2ToL1MessagePasser.sol
DEPLOYER_WHITE_LIST_PATH=${BASE_PATH}/src/legacy/DeployerWhitelist.sol
WTON_PATH=${BASE_PATH} # where is it
L2_CROSS_DOMAIN_MESSENGER_PATH=${BASE_PATH}/src/L2/L2CrossDomainMessenger.sol
L2_STANDARD_BRIDGE_PATH=${BASE_PATH}/src/L2/L2StandardBridge.sol
SEQUENCER_FEE_VAULT_PATH=${BASE_PATH}/src/L2/SequencerFeeVault.sol
OPTIMISM_MINTABLE_ERC20_FACTORY_PATH=${BASE_PATH}/src/universal/OptimismMintableERC20Factory.sol
L1_BLOCK_NUMBER_PATH=${BASE_PATH}/src/legacy/L1BlockNumber.sol
GAS_PRICE_ORACLE_PATH=${BASE_PATH}/src/L2/GasPriceOracle.sol
L1_BLOCK_PATH=${BASE_PATH}/src/L2/L1Block.sol
GOVERNANCE_TOKEN_PATH=${BASE_PATH}/src/governance/GovernanceToken.sol
LEGACY_MESSAGE_PASSER_PATH=${BASE_PATH}/src/legacy/LegacyMessagePasser.sol
L2_ERC721_BRIDGE_PATH=${BASE_PATH}/src/L2/L2ERC721Bridge.sol
OPTIMISM_MINTABLE_ERC721_FACTORY_PATH=${BASE_PATH}/src/universal/OptimismMintableERC721Factory.sol
PROXY_ADMIN_PATH=${BASE_PATH}/src/universal/ProxyAdmin.sol
BASE_FEE_VAULT_PATH=${BASE_PATH}/src/L2/BaseFeeVault.sol
L1_FEE_VAULT_PATH=${BASE_PATH}/src/L2/L1FeeVault.sol
SCHEMA_REGISTRY_PATH=${BASE_PATH}/src/EAS/SchemaRegistry.sol
EAS_PATH=${BASE_PATH}/src/EAS/EAS.sol
WETH_PATH=${BASE_PATH}/src/L2/WETH.sol

function verify_contract() {
  forge verify-contract \
    --chain-id $CHAIN_ID \
    --verifier blockscout \
    --verifier-url $VERIFIER_URL \
    --compiler-version v0.8.23+commit.f704f362 \
    $([[ -n $1 ]] && echo "--constructor-args $1") \
    $2 \
    $3
}

function verify_proxy() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0x4200000000000000000000000000000000000001)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000001"]}
  verify_contract $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${PROXY_PATH}:Proxy"
}

function verify_L2ToL1MessagePasser() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000016"]}
  verify_contract "" $CONTRACT_ADDR "${L2_TO_L1_MESSAGE_PASSER_PATH}:L2ToL1MessagePasser"
}

function verify_DeployerWhitelist() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000002"]}
  verify_contract "" $CONTRACT_ADDR "${DEPLOYER_WHITE_LIST_PATH}:DeployerWhitelist"
}

function verify_WTON() {
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000006"]}
  # TODO : Verify WTON
}

function verify_L2CrossDomainMessenger() {
  CONSTRUCTOR_ARGS=$(cast abi-encode "constructor(address)" 0x4200000000000000000000000000000000000007)
  CONTRACT_ADDR=${contracts["4200000000000000000000000000000000000007"]}
  verify_contract $CONSTRUCTOR_ARGS $CONTRACT_ADDR "${L2_CROSS_DOMAIN_MESSENGER_PATH}:L2CrossDomainMessenger"
}

verify_proxy
verify_L2ToL1MessagePasser
verify_DeployerWhitelist
verify_L2CrossDomainMessenger
