#!/bin/bash
set -eu

if [ "$#" -ne 2 ]; then
  echo "This script requires 2 arguments: path to JSON artifact and package name."
  exit 1
fi

ARTIFACT_PATH=$1
PACKAGE=$2
FILENAME=$(basename -- "$ARTIFACT_PATH")
CONTRACT_NAME="${FILENAME%.*}"
TYPE_LOWER=$(echo "${CONTRACT_NAME}" | tr '[:upper:]' '[:lower:]')

# Ensure jq and abigen are available
need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Need '$1' (command not found)"
    exit 1
  fi
}

need_cmd jq
need_cmd abigen

ABI_PATH=$(mktemp)
BIN_PATH=$(mktemp)

jq -r '.abi' "$ARTIFACT_PATH" >"$ABI_PATH"
jq -r '.bytecode' "$ARTIFACT_PATH" >"$BIN_PATH"

# Generate Go bindings
abigen --abi "$ABI_PATH" --bin "$BIN_PATH" --pkg "$PACKAGE" --type "$CONTRACT_NAME" --out "./$PACKAGE/$TYPE_LOWER.go"

echo "Bindings for $CONTRACT_NAME generated in ./$PACKAGE/$TYPE_LOWER.go"
