#!/usr/bin/env bash

set -e

if ! command -v forge &>/dev/null; then
  echo "forge could not be found. Please install forge by running:"
  echo "curl -L https://foundry.paradigm.xyz | bash"
  exit
fi

contracts=(
  src/L1/L1CrossDomainMessenger.sol:L1CrossDomainMessenger
  src/L1/L1StandardBridge.sol:L1StandardBridge
  src/L1/L2OutputOracle.sol:L2OutputOracle
  src/L1/OptimismPortal.sol:OptimismPortal
  src/L1/SystemConfig.sol:SystemConfig
  src/L1/L1ERC721Bridge.sol:L1ERC721Bridge
  src/legacy/DeployerWhitelist.sol:DeployerWhitelist
  src/L2/L1Block.sol:L1Block
  src/legacy/L1BlockNumber.sol:L1BlockNumber
  src/L2/L2CrossDomainMessenger.sol:L2CrossDomainMessenger
  src/L2/L2StandardBridge.sol:L2StandardBridge
  src/L2/L2ToL1MessagePasser.sol:L2ToL1MessagePasser
  src/legacy/LegacyERC20NativeToken.sol:LegacyERC20NativeToken
  src/L2/SequencerFeeVault.sol:SequencerFeeVault
  src/L2/BaseFeeVault.sol:BaseFeeVault
  src/L2/L1FeeVault.sol:L1FeeVault
  src/L2/L2ERC721Bridge.sol:L2ERC721Bridge
  src/vendor/WNativeToken.sol:WNativeToken
  src/universal/ProxyAdmin.sol:ProxyAdmin
  src/universal/Proxy.sol:Proxy
  src/legacy/L1ChugSplashProxy.sol:L1ChugSplashProxy
  src/universal/OptimismMintableERC20.sol:OptimismMintableERC20
  src/universal/OptimismMintableERC20Factory.sol:OptimismMintableERC20Factory
  src/dispute/DisputeGameFactory.sol:DisputeGameFactory
  src/tokamak-contracts/USDC/L1/tokamak-UsdcBridge/L1UsdcBridge.sol:L1UsdcBridge
  src/tokamak-contracts/USDC/L1/tokamak-UsdcBridge/L1UsdcBridgeProxy.sol:L1UsdcBridgeProxy
  src/tokamak-contracts/USDC/L2/tokamak-UsdcBridge/L2UsdcBridge.sol:L2UsdcBridge
  src/tokamak-contracts/USDC/L2/tokamak-UsdcBridge/L2UsdcBridgeProxy.sol:L2UsdcBridgeProxy
  src/tokamak-contracts/USDC/L2/tokamak-USDC/minting/MasterMinter.sol:MasterMinter
  src/tokamak-contracts/USDC/L2/tokamak-USDC/util/SignatureChecker.sol:SignatureChecker
  src/tokamak-contracts/USDC/L2/tokamak-USDC/v2/FiatTokenV2_2.sol:FiatTokenV2_2
  src/tokamak-contracts/USDC/L2/tokamak-USDC/v1/FiatTokenProxy.sol:FiatTokenProxy
)

dir=$(dirname "$0")

echo "Creating storage layout diagrams.."

echo "=======================" >$dir/../.storage-layout
echo "👁👁 STORAGE LAYOUT snapshot 👁👁" >>$dir/../.storage-layout
echo "=======================" >>$dir/../.storage-layout

for contract in ${contracts[@]}; do
  echo -e "\n=======================" >>$dir/../.storage-layout
  echo "➡ $contract" >>$dir/../.storage-layout
  echo -e "=======================\n" >>$dir/../.storage-layout
  forge inspect --pretty $contract storage-layout >>$dir/../.storage-layout
done
echo "Storage layout snapshot stored at $dir/../.storage-layout"
