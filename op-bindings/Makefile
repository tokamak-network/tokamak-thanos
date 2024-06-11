pkg := bindings
monorepo-base := $(shell dirname $(realpath .))
contracts-dir := $(monorepo-base)/packages/tokamak/contracts-bedrock
contracts-list := ./artifacts.json
contracts-hardhat-list := ./artifacts-hardhat.json
log-level := info
ETHERSCAN_APIKEY_ETH ?=
ETHERSCAN_APIKEY_OP ?=
RPC_URL_ETH ?=
RPC_URL_OP ?=

all: version mkdir bindings

version:
	forge --version
	abigen --version

compile:
	cd $(contracts-dir) && \
		forge clean && \
		pnpm build

bindings: bindgen-local

bindgen-local: compile bindgen-generate-local

bindgen-generate-local:
	go run ./cmd/ \
		generate \
		--metadata-out ./$(pkg) \
		--bindings-package $(pkg) \
		--contracts-list $(contracts-list) \
		--log.level $(log-level) \
		local \
		--forge-artifacts $(contracts-dir)/forge-artifacts \
		--hardhat-artifacts $(contracts-dir)/hardhat-artifacts \
		--hardhat-contracts-list $(contracts-hardhat-list)

bindgen-remote:
	go run ./cmd/ \
		generate \
		--metadata-out ./$(pkg) \
		--bindings-package $(pkg) \
		--contracts-list $(contracts-list) \
		--log.level $(log-level) \
		remote \
		--etherscan.apikey.eth $(ETHERSCAN_APIKEY_ETH) \
		--etherscan.apikey.op $(ETHERSCAN_APIKEY_OP) \
		--rpc.url.eth $(RPC_URL_ETH) \
		--rpc.url.op $(RPC_URL_OP)

bindings-legacy:
	go run ./gen/main.go \
    -hardhat-artifacts $(contracts-dir)/hardhat-artifacts \
    -out ./bindings \
    -hardhat-contracts ./artifacts-hardhat.json \
    -source-maps MIPS,PreimageOracle \
    -package $(pkg) \
    -monorepo-base $(monorepo-base)

hardhat-bindings:
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/openzeppeline/TransparentUpgradeableProxy.sol/TransparentUpgradeableProxy.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/swap-router-contracts/QuoterV2.sol/QuoterV2.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/swap-router-contracts/SwapRouter02.sol/SwapRouter02.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-core/UniswapV3Factory.sol/UniswapV3Factory.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/NFTDescriptor.sol/NFTDescriptor.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/NonfungiblePositionManager.sol/NonfungiblePositionManager.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/NonfungibleTokenPositionDescriptor.sol/NonfungibleTokenPositionDescriptor.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/TickLens.sol/TickLens.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/v3-periphery/UniswapInterfaceMulticall.sol/UniswapInterfaceMulticall.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/universal-router/UniversalRouter.sol/UniversalRouter.json $(pkg)
	./gen_bindings.sh $(contracts-dir)/hardhat-artifacts/universal-router/UnsupportedProtocol.sol/UnsupportedProtocol.json $(pkg)

mkdir:
	mkdir -p $(pkg)

clean-contracts:
	cd $(contracts-dir) && \
    pnpm clean

clean:
	rm -rf $(pkg)

test:
	go test ./...
