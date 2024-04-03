SHELL := /usr/bin/env bash

pkg := bindings
monorepo-base := $(shell dirname $(realpath .))
contracts-dir := $(monorepo-base)/packages/tokamak/contracts-bedrock
hardhat-artifacts-dir := $(monorepo-base)/packages/tokamak/contracts-bedrock


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

mkdir:
	mkdir -p $(pkg)

clean-contracts:
	cd $(contracts-dir) && \
		pnpm clean

clean:
	rm -rf $(pkg)

test:
	go test ./...
