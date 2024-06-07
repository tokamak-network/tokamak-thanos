#!/usr/bin/env bash

#
# Produces a testnet specification and a genesis state where the genesis time
# is now + $GENESIS_DELAY.
#
# Generates datadirs for multiple validator keys according to the
# $VALIDATOR_COUNT and $BN_COUNT variables.
#

set -o nounset -o errexit -o pipefail

source /vars.env

NOW=`date +%s`
GENESIS_TIME=`expr $NOW + $GENESIS_DELAY`
ETH1_BLOCK_HASH=`jq -r '.genesis.l1.hash' /rollup.json`
echo $ETH1_BLOCK_HASH

lcli \
	new-testnet \
	--spec $SPEC_PRESET \
	--deposit-contract-address $DEPOSIT_CONTRACT_ADDRESS \
	--testnet-dir $TESTNET_DIR \
	--min-genesis-active-validator-count $GENESIS_VALIDATOR_COUNT \
	--min-genesis-time $GENESIS_TIME \
	--genesis-delay $GENESIS_DELAY \
	--genesis-fork-version $GENESIS_FORK_VERSION \
	--altair-fork-epoch $ALTAIR_FORK_EPOCH \
	--bellatrix-fork-epoch $BELLATRIX_FORK_EPOCH \
	--capella-fork-epoch $CAPELLA_FORK_EPOCH \
	--deneb-fork-epoch $DENEB_FORK_EPOCH \
	--ttd $TTD \
	--eth1-block-hash $ETH1_BLOCK_HASH \
	--eth1-id $CHAIN_ID \
	--eth1-follow-distance 128 \
	--seconds-per-slot $SECONDS_PER_SLOT \
	--seconds-per-eth1-block $SECONDS_PER_ETH1_BLOCK \
	--proposer-score-boost "$PROPOSER_SCORE_BOOST" \
	--validator-count $GENESIS_VALIDATOR_COUNT \
	--interop-genesis-state \
	--force

echo Specification and genesis.ssz generated at $TESTNET_DIR.

cd /beacon_data/local_testnet/testnet
ls -al
echo "Generating $VALIDATOR_COUNT validators concurrently... (this may take a while)"

lcli \
	insecure-validators \
	--count $VALIDATOR_COUNT \
	--base-dir $DATADIR \
	--node-count $VC_COUNT

echo Validators generated with keystore passwords at $DATADIR.

# Function to output SLOT_PER_EPOCH for mainnet or minimal
# get_spec_preset_value() {
#   case "$SPEC_PRESET" in
#     mainnet)   echo 32 ;;
#     minimal)   echo 8  ;;
#     gnosis)    echo 16 ;;
#     *)         echo "Unsupported preset: $SPEC_PRESET" >&2; exit 1 ;;
#   esac
# }

# SLOT_PER_EPOCH=32
# echo "slot_per_epoch=$SLOT_PER_EPOCH"

# # Update future hardforks time in the EL genesis file based on the CL genesis time
# GENESIS_TIME=$(lcli pretty-ssz --spec $SPEC_PRESET --testnet-dir $TESTNET_DIR BeaconState $TESTNET_DIR/genesis.ssz | jq | grep -Po 'genesis_time": "\K.*\d')
# echo $GENESIS_TIME
# CAPELLA_TIME=$((GENESIS_TIME + (CAPELLA_FORK_EPOCH * $SLOT_PER_EPOCH * SECONDS_PER_SLOT)))
# echo $CAPELLA_TIME
# sed -i 's/"shanghaiTime".*$/"shanghaiTime": '"$CAPELLA_TIME"',/g' /genesis-l1.json
# CANCUN_TIME=$((GENESIS_TIME + (DENEB_FORK_EPOCH * $SLOT_PER_EPOCH * SECONDS_PER_SLOT)))
# echo $CANCUN_TIME
# sed -i 's/"cancunTime".*$/"cancunTime": '"$CANCUN_TIME"',/g' /genesis-l1.json
# cat /genesis-l1.json