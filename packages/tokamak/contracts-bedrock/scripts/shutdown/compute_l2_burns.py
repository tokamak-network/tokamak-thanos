#!/usr/bin/env python3
import json
import os
import subprocess
import sys


L2_STANDARD_BRIDGE = "0x4200000000000000000000000000000000000010"
L2_USDC_PREDEPLOY = "0x4200000000000000000000000000000000000778"
L2_PREDEPLOY_ETH = "0x4200000000000000000000000000000000000486"
WITHDRAWAL_INIT_TOPIC0 = "0x73d170910aba9e6d50b102db522b1dbcd796216f5128b445aa2135272886497e"
BURN_TOPIC0 = "0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5"


def pad_topic(addr):
    return "0x" + addr.lower().replace("0x", "").rjust(64, "0")


def load_tokens(chain_id):
    data_dir = os.getenv("DATA_DIR", "data")
    path = os.path.join(data_dir, f"l2-tokens-{chain_id}.json")
    if not os.path.exists(path):
        path = os.path.join(data_dir, "l2-tokens.json")
    if not os.path.exists(path):
        sys.stderr.write("Token list file not found\n")
        sys.exit(1)
    with open(path, "r") as handle:
        return json.load(handle)


def cast_get_logs(rpc_url, filter_obj):
    cmd = [
        "cast",
        "rpc",
        "--rpc-url",
        rpc_url,
        "eth_getLogs",
        json.dumps(filter_obj),
    ]
    try:
        raw = subprocess.check_output(cmd)
    except subprocess.CalledProcessError as exc:
        sys.stderr.write(exc.output.decode("utf-8", errors="ignore"))
        raise
    payload = json.loads(raw.decode("utf-8"))
    if isinstance(payload, dict) and "result" in payload:
        return payload["result"]
    return payload


def sum_burn_amount(logs):
    total = 0
    for log in logs:
        data = log.get("data", "0x")
        if data.startswith("0x"):
            data = data[2:]
        if len(data) < 64:
            continue
        total += int(data[0:64], 16)
    return total


def sum_withdrawal_amount(logs):
    total = 0
    for log in logs:
        data = log.get("data", "0x")
        if data.startswith("0x"):
            data = data[2:]
        if len(data) < 128:
            continue
        total += int(data[64:128], 16)
    return total


def main():
    if len(sys.argv) < 3:
        print("Usage: compute_l2_burns.py <l2_rpc_url> <chain_id> [start_block]")
        sys.exit(1)

    rpc_url = sys.argv[1]
    chain_id = sys.argv[2]
    start_block = int(os.environ.get("L2_START_BLOCK", "0"))
    if len(sys.argv) >= 4:
        start_block = int(sys.argv[3])

    tokens = load_tokens(chain_id)
    results = []
    l2_usdc_bridge = os.environ.get("L2_USDC_BRIDGE_PROXY")

    for token in tokens:
        if token.lower() == L2_PREDEPLOY_ETH.lower():
            results.append(
                {
                    "l2Token": token,
                    "burnTotal": "0",
                    "withdrawalTotal": "0",
                    "extraBurn": "0",
                }
            )
            continue

        burn_filter = {
            "fromBlock": hex(start_block),
            "toBlock": "latest",
            "address": token,
            "topics": [BURN_TOPIC0],
        }
        bridge_address = L2_STANDARD_BRIDGE
        if l2_usdc_bridge and token.lower() == L2_USDC_PREDEPLOY.lower():
            bridge_address = l2_usdc_bridge

        withdrawal_filter = {
            "fromBlock": hex(start_block),
            "toBlock": "latest",
            "address": bridge_address,
            "topics": [WITHDRAWAL_INIT_TOPIC0, None, pad_topic(token)],
        }

        try:
            burn_logs = cast_get_logs(rpc_url, burn_filter)
            withdrawal_logs = cast_get_logs(rpc_url, withdrawal_filter)
        except subprocess.CalledProcessError as exc:
            sys.stderr.write(f"Failed to fetch logs for token {token}\n")
            sys.exit(1)

        burn_total = sum_burn_amount(burn_logs)
        withdrawal_total = sum_withdrawal_amount(withdrawal_logs)
        extra_burn = burn_total - withdrawal_total
        if extra_burn < 0:
            extra_burn = 0

        results.append(
            {
                "l2Token": token,
                "burnTotal": str(burn_total),
                "withdrawalTotal": str(withdrawal_total),
                "extraBurn": str(extra_burn),
            }
        )

    data_dir = os.getenv("DATA_DIR", "data")
    out_path = os.path.join(data_dir, f"l2-burns-{chain_id}.json")
    with open(out_path, "w") as handle:
        json.dump(results, handle, indent=2)
    print(f"Wrote {out_path} ({len(results)} tokens)")


if __name__ == "__main__":
    main()
