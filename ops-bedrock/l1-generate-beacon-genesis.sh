#!/bin/sh

set -eu

echo "eth2-testnet-genesis path: $(which eth2-testnet-genesis)"

eth2-testnet-genesis electra \
  --config=./beacon-data/config.yaml \
  --preset-phase0=minimal \
  --preset-altair=minimal \
  --preset-bellatrix=minimal \
  --preset-capella=minimal \
  --preset-deneb=minimal \
  --preset-electra=minimal \
  --eth1-config=../.devnet/genesis-l1.json \
  --state-output=../.devnet/genesis-l1.ssz \
  --tranches-dir=../.devnet/tranches \
  --mnemonics=./beacon-data/mnemonics.yaml \
  --eth1-withdrawal-address=0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa \
  --eth1-match-genesis-time


# Check if eth-genesis-state-generator is available
# if ! command -v eth-genesis-state-generator >/dev/null 2>&1; then
#   echo "Error: eth-genesis-state-generator not found"
#   echo "Please install it with: go install github.com/ethpandaops/eth-beacon-genesis@latest"
#   exit 1
# fi

# Generate beacon genesis state using eth-beacon-genesis
# eth-genesis-state-generator beaconchain \
#   --eth1-config=../.devnet/genesis-l1.json \
#   --config=./beacon-data/config.yaml \
#   --mnemonics=./beacon-data/mnemonics.yaml \
#   --state-output=../.devnet/genesis-l1.ssz



