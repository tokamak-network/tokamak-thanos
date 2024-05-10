SHELL := /usr/bin/env bash

pkg := bindings
monorepo-base := $(shell dirname $(realpath .))
contracts-dir := $(monorepo-base)/packages/tokamak/contracts-bedrock
contracts-list := ./artifacts.json
log-level := info

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
		--forge-artifacts $(contracts-dir)/forge-artifacts

mkdir:
	mkdir -p $(pkg)

clean-contracts:
	cd $(contracts-dir) && \
    pnpm clean

clean:
	rm -rf $(pkg)

test:
	go test ./...