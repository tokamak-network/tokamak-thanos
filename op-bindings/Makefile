SHELL := /usr/bin/env bash

pkg := bindings
monorepo-base := $(shell dirname $(realpath .))
contracts-dir := $(monorepo-base)/packages/tokamak/contracts-bedrock
hardhat-artifacts-dir := $(monorepo-base)/packages/tokamak/contracts-bedrock
ARTIFACTS_PATH := $(monorepo-base)/packages/tokamak/contracts-bedrock/hardhat-artifacts

all: version mkdir bindings

version:
	forge --version
	abigen --version

compile:
	cd $(contracts-dir) && \
		forge clean && \
		pnpm build

bindings: compile bindings-build

bindings-build:
	go run ./gen/main.go \
    -forge-artifacts $(contracts-dir)/forge-artifacts \
    -hardhat-artifacts $(hardhat-artifacts-dir)/hardhat-artifacts \
    -out ./bindings \
    -forge-contracts ./artifacts.json \
    -hardhat-contracts ./artifacts-hardhat.json \
    -source-maps MIPS,PreimageOracle \
    -package $(pkg) \
    -monorepo-base $(monorepo-base)

go-bindings:
	./gen_bindings.sh $(ARTIFACTS_PATH)/openzeppeline/TransparentUpgradeableProxy.sol/TransparentUpgradeableProxy.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/permit2/Permit2.sol/Permit2.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/swap-router-contracts/QuoterV2.sol/QuoterV2.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/swap-router-contracts/SwapRouter02.sol/SwapRouter02.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/v3-core/UniswapV3Factory.sol/UniswapV3Factory.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/v3-periphery/NFTDescriptor.sol/NFTDescriptor.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/v3-periphery/NonfungiblePositionManager.sol/NonfungiblePositionManager.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/v3-periphery/NonfungibleTokenPositionDescriptor.sol/NonfungibleTokenPositionDescriptor.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/v3-periphery/TickLens.sol/TickLens.json $(pkg)
	./gen_bindings.sh $(ARTIFACTS_PATH)/v3-periphery/UniswapInterfaceMulticall.sol/UniswapInterfaceMulticall.json $(pkg)

mkdir:
	mkdir -p $(pkg)

clean-contracts:
	cd $(contracts-dir) && \
    pnpm clean

clean:
	rm -rf $(pkg)

test:
	go test ./...
