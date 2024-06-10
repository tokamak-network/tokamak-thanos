#!/usr/bin/env bash

set -Eeuo pipefail

source /vars.env

# Setup file permissions
chmod +x /setup_time.sh
chmod u+w /genesis-l1.json

# Function to output SLOT_PER_EPOCH for mainnet or minimal
get_spec_preset_value() {
  case "$SPEC_PRESET" in
    mainnet)   echo 32 ;;
    minimal)   echo 8  ;;
    gnosis)    echo 16 ;;
    *)         echo "Unsupported preset: $SPEC_PRESET" >&2; exit 1 ;;
  esac
}

SLOT_PER_EPOCH=$(get_spec_preset_value $SPEC_PRESET)
echo "slot_per_epoch=$SLOT_PER_EPOCH"

# Update future hardforks time in the EL genesis file based on the CL genesis time
GENESIS_TIME=$(lcli pretty-ssz --spec $SPEC_PRESET --testnet-dir $TESTNET_DIR BeaconState $TESTNET_DIR/genesis.ssz | jq | grep -Po 'genesis_time": "\K.*\d')
echo $GENESIS_TIME
CAPELLA_TIME=$((GENESIS_TIME + (CAPELLA_FORK_EPOCH * $SLOT_PER_EPOCH * SECONDS_PER_SLOT)))
echo $CAPELLA_TIME
echo "$(sed 's/"shanghaiTime".*$/"shanghaiTime": '"$CAPELLA_TIME"',/g' /genesis-l1.json)" > /genesis-l1.json
# sed -i 's/"shanghaiTime".*$/"shanghaiTime": '"$CAPELLA_TIME"',/g' /genesis-l1.json
CANCUN_TIME=$((GENESIS_TIME + (DENEB_FORK_EPOCH * $SLOT_PER_EPOCH * SECONDS_PER_SLOT)))
echo $CANCUN_TIME
# sed -i 's/"cancunTime".*$/"cancunTime": '"$CANCUN_TIME"',/g' /genesis-l1.json
echo "$(sed 's/"cancunTime".*$/"cancunTime": '"$CANCUN_TIME"',/g' /genesis-l1.json)" > /genesis-l1.json

cat /genesis-l1.json

# Move the modified genesis-l1.json to the host directory
mv /genesis-l1.json ./../.devnet/genesis-l1.json