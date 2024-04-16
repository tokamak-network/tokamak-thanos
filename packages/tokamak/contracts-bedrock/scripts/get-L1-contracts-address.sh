#!/usr/bin/env bash

if [ -z $1 ]; then
  echo "Please input network name..."
fi

DEPLOYMENTS_PATH=$(cd $(dirname $0)/../deployments && pwd -P)
ADDRESS_FILE_PATH=$(cd $(dirname $0)/../genesis && pwd -P)
ARTIFACTS_PATH=${DEPLOYMENTS_PATH}/$1

if [ ! -d "$ARTIFACTS_PATH" ]; then
  echo "$1 is not exist..."
fi

AddressManager=$(cat ${ARTIFACTS_PATH}/AddressManager.json | jq '.address')
L1CrossDomainMessenger=$(cat ${ARTIFACTS_PATH}/L1CrossDomainMessenger.json | jq '.address')
L1CrossDomainMessengerProxy=$(cat ${ARTIFACTS_PATH}/L1CrossDomainMessengerProxy.json | jq '.address')
L1ERC721Bridge=$(cat ${ARTIFACTS_PATH}/L1ERC721Bridge.json | jq '.address')
L1ERC721BridgeProxy=$(cat ${ARTIFACTS_PATH}/L1ERC721BridgeProxy.json | jq '.address')
L1StandardBridge=$(cat ${ARTIFACTS_PATH}/L1StandardBridge.json | jq '.address')
L1StandardBridgeProxy=$(cat ${ARTIFACTS_PATH}/L1StandardBridgeProxy.json | jq '.address')
L2OutputOracle=$(cat ${ARTIFACTS_PATH}/L2OutputOracle.json | jq '.address')
L2OutputOracleProxy=$(cat ${ARTIFACTS_PATH}/L2OutputOracleProxy.json | jq '.address')
OptimismMintableERC20Factory=$(cat ${ARTIFACTS_PATH}/OptimismMintableERC20Factory.json | jq '.address')
OptimismMintableERC20FactoryProxy=$(cat ${ARTIFACTS_PATH}/OptimismMintableERC20FactoryProxy.json | jq '.address')
OptimismPortal=$(cat ${ARTIFACTS_PATH}/OptimismPortal.json | jq '.address')
OptimismPortalProxy=$(cat ${ARTIFACTS_PATH}/OptimismPortalProxy.json | jq '.address')
ProtocolVersions=$(cat ${ARTIFACTS_PATH}/ProtocolVersions.json | jq '.address')
ProtocolVersionsProxy=$(cat ${ARTIFACTS_PATH}/ProtocolVersionsProxy.json | jq '.address')
ProxyAdmin=$(cat ${ARTIFACTS_PATH}/ProxyAdmin.json | jq '.address')
SystemConfig=$(cat ${ARTIFACTS_PATH}/SystemConfig.json | jq '.address')
SystemConfigProxy=$(cat ${ARTIFACTS_PATH}/SystemConfigProxy.json | jq '.address')
L1UsdcBridge=$(cat ${ARTIFACTS_PATH}/L1UsdcBrige.json | jq '.address')
L1UsdcBridgeProxy=$(cat ${ARTIFACTS_PATH}/L1UsdcBridgeProxy.json | jq '.address')

echo "AddressManager: $AddressManager"
echo "L1CrossDomainMessenger: $L1CrossDomainMessenger"
echo "L1CrossDomainMessengerProxy: $L1CrossDomainMessengerProxy"
echo "L1ERC721Bridge: $L1ERC721Bridge"
echo "L1ERC721BridgeProxy: $L1ERC721BridgeProxy"
echo "L1StandardBridge: $L1StandardBridge"
echo "L1StandardBridgeProxy: $L1StandardBridgeProxy"
echo "L2OutputOracle: $L2OutputOracle"
echo "L2OutputOracleProxy: $L2OutputOracleProxy"
echo "OptimismMintableERC20Factory: $OptimismMintableERC20Factory"
echo "OptimismMintableERC20FactoryProxy: $OptimismMintableERC20FactoryProxy"
echo "OptimismPortal: $OptimismPortal"
echo "OptimismPortalProxy: $OptimismPortalProxy"
echo "ProtocolVersions: $ProtocolVersions"
echo "ProtocolVersionsProxy: $ProtocolVersionsProxy"
echo "ProxyAdmin: $ProxyAdmin"
echo "SystemConfig: $SystemConfig"
echo "SystemConfigProxy: $SystemConfigProxy"
echo "L1UsdcBridge: $L1UsdcBridge"
echo "L1UsdcBridgeProxy: $L1UsdcBridgeProxy"

cat << EOF > ${ADDRESS_FILE_PATH}/$1/address.json
{
  "AddressManager": $AddressManager,
  "L1CrossDomainMessenger": $L1CrossDomainMessenger,
  "L1CrossDomainMessengerProxy": $L1CrossDomainMessengerProxy,
  "L1ERC721Bridge": $L1ERC721Bridge,
  "L1ERC721BridgeProxy": $L1ERC721BridgeProxy,
  "L1StandardBridge": $L1StandardBridge,
  "L1StandardBridgeProxy": $L1StandardBridgeProxy,
  "L2OutputOracle": $L2OutputOracle,
  "L2OutputOracleProxy": $L2OutputOracleProxy,
  "OptimismMintableERC20Factory": $OptimismMintableERC20Factory,
  "OptimismMintableERC20FactoryProxy": $OptimismMintableERC20FactoryProxy,
  "OptimismPortal": $OptimismPortal,
  "OptimismPortalProxy": $OptimismPortalProxy,
  "ProtocolVersions": $ProtocolVersions,
  "ProtocolVersionsProxy": $ProtocolVersionsProxy,
  "ProxyAdmin": $ProxyAdmin,
  "SystemConfig": $SystemConfig,
  "SystemConfigProxy": $SystemConfigProxy,
  "L1UsdcBridge": $L1UsdcBridge,
  "L1UsdcBridgeProxy": $L1UsdcBridgeProxy,
}
EOF
