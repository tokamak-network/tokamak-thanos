#!/usr/bin/env python3
import json
import subprocess
import sys


def main():
    if len(sys.argv) != 3:
        print("Usage: compute_finalized_native_withdrawals.py <rpc_url> <bridge_address>")
        sys.exit(1)

    rpc_url = sys.argv[1]
    bridge_address = sys.argv[2]

    sig = "NativeTokenWithdrawalFinalized(address,address,uint256,bytes)"
    import os
    l1_start_block = os.environ.get("L1_START_BLOCK", "earliest")

    # Alchemy Free Tier workaround: don't use 0 to latest if start block isn't specified.
    # If using Alchemy, 'earliest' might fail. Better to use a specific number.
    if l1_start_block == "earliest":
        # Using 'latest' is the safest default for Alchemy Free Tier (range = 0)
        l1_start_block = "latest"

    cmd = [
        "cast",
        "logs",
        sig,
        "--address",
        bridge_address,
        "--from-block",
        l1_start_block,
        "--to-block",
        "latest",
        "--rpc-url",
        rpc_url,
        "--json",
    ]

    try:
        raw = subprocess.check_output(cmd)
    except subprocess.CalledProcessError as exc:
        sys.stderr.write(exc.output.decode("utf-8", errors="ignore"))
        sys.exit(1)

    logs = json.loads(raw.decode("utf-8"))
    total = 0
    for log in logs:
        data = log.get("data", "0x")
        if data.startswith("0x"):
            data = data[2:]
        if len(data) < 64:
            continue
        amount = int(data[0:64], 16)
        total += amount

    print(f"DEC:{total}")


if __name__ == "__main__":
    main()
