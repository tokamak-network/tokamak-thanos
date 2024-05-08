SHELL := /usr/bin/env bash

pkg := bindings
monorepo-base := $(shell dirname $(realpath .))
contracts-dir := $(monorepo-base)/packages/tokamak/contracts-bedrock

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
    -hardhat-artifacts $(contracts-dir)/hardhat-artifacts \
    -out ./bindings \
    -forge-contracts ./artifacts.json \
    -hardhat-contracts ./artifacts-hardhat.json \
    -source-maps MIPS,PreimageOracle \
    -package $(pkg) \
    -monorepo-base $(monorepo-base)

hardhat-bindings:
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/openzeppeline/TransparentUpgradeableProxy.sol/TransparentUpgradeableProxy.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/permit2/Permit2.sol/Permit2.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/swap-router-contracts/QuoterV2.sol/QuoterV2.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/swap-router-contracts/SwapRouter02.sol/SwapRouter02.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-core/UniswapV3Factory.sol/UniswapV3Factory.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/NFTDescriptor.sol/NFTDescriptor.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/NonfungiblePositionManager.sol/NonfungiblePositionManager.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/NonfungibleTokenPositionDescriptor.sol/NonfungibleTokenPositionDescriptor.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/TickLens.sol/TickLens.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/UniswapInterfaceMulticall.sol/UniswapInterfaceMulticall.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/universal-router/UniversalRouter.sol/UniversalRouter.json $(pkg)


mkdir:
	mkdir -p $(pkg)

clean-contracts:
	cd $(contracts-dir) && \
    pnpm clean

clean:
	rm -rf $(pkg)

test:
	go test ./...
