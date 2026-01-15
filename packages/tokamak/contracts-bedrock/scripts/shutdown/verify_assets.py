#!/usr/bin/env python3
"""
Asset Snapshot Verification Script

Verifies that the generated asset snapshot matches on-chain balances.

Usage:
    python3 verify_assets.py <snapshot_file> <rpc_url>

Example:
    python3 verify_assets.py data/generate-assets-111551119090.json https://rpc.thanos-sepolia.tokamak.network
"""

import json
import sys
from web3 import Web3
from eth_abi import decode

# Constants
NATIVE_TOKEN = "0xdeaDDeADDEaDdeaDdEAddEADDEAdDeadDEADDEaD"
ERC20_BALANCE_OF_SELECTOR = "0x70a08231"  # balanceOf(address)

# ANSI color codes
GREEN = '\033[92m'
YELLOW = '\033[93m'
RED = '\033[91m'
BLUE = '\033[94m'
RESET = '\033[0m'


class VerificationStats:
    def __init__(self):
        self.total_tokens = 0
        self.total_claims = 0
        self.total_amount = 0
        self.mismatch_count = 0
        self.hash_mismatch_count = 0
        self.balance_query_failures = 0


def print_header(snapshot_file, chain_id, block_number):
    print("=" * 60)
    print(f"{BLUE}       Asset Snapshot Verification Report{RESET}")
    print("=" * 60)
    print(f"Target File:  {snapshot_file}")
    print(f"Chain ID:     {chain_id}")
    print(f"Block Number: {block_number}")
    print("-" * 60 + "\n")


def print_summary(stats):
    print("\n" + "=" * 60)
    print(f"{BLUE}              Verification Summary{RESET}")
    print("=" * 60)
    print(f"Total Tokens:          {stats.total_tokens}")
    print(f"Total Claims:          {stats.total_claims}")
    print(f"Total Amount:          {stats.total_amount} wei")
    print(f"Balance Mismatches:    {stats.mismatch_count}")
    print(f"Hash Mismatches:       {stats.hash_mismatch_count}")
    print(f"Balance Query Failures: {stats.balance_query_failures}")

    if stats.mismatch_count == 0 and stats.hash_mismatch_count == 0:
        print(f"\n{GREEN}[SUCCESS] All verifications passed!{RESET}")
    else:
        print(f"\n{YELLOW}[WARNING] Verification issues detected!{RESET}")
        print("Please review the mismatches above.")

    print("=" * 60 + "\n")


def get_balance(w3, token_address, account_address):
    """Get token balance for an account"""
    try:
        if token_address.lower() == NATIVE_TOKEN.lower():
            # Native token balance
            return w3.eth.get_balance(account_address), True
        else:
            # ERC20 token balance
            data = ERC20_BALANCE_OF_SELECTOR + account_address[2:].zfill(64)
            result = w3.eth.call({
                'to': token_address,
                'data': data
            })
            balance = int(result.hex(), 16)
            return balance, True
    except Exception as e:
        return 0, False


def verify_hash(l1_token, claimer, amount, stored_hash):
    """Verify the hash matches keccak256(abi.encodePacked(l1Token, claimer, amount))"""
    from eth_utils import keccak

    # Convert addresses to bytes (remove 0x and pad to 20 bytes)
    l1_token_bytes = bytes.fromhex(l1_token[2:].zfill(40))
    claimer_bytes = bytes.fromhex(claimer[2:].zfill(40))

    # Convert amount to bytes32 (big-endian)
    amount_bytes = int(amount).to_bytes(32, byteorder='big')

    # Concatenate and hash
    packed = l1_token_bytes + claimer_bytes + amount_bytes
    calculated_hash = '0x' + keccak(packed).hex()

    return calculated_hash.lower() == stored_hash.lower()


def verify_claim(w3, claim, l1_token, l2_token, stats):
    """Verify a single claim"""
    claimer = claim['claimer']
    amount = int(claim['amount'])
    stored_hash = claim['hash']

    # Verify hash
    if not verify_hash(l1_token, claimer, amount, stored_hash):
        print(f"    {YELLOW}[WARN]{RESET} Hash mismatch: {claimer}")
        print(f"      Stored: {stored_hash}")
        stats.hash_mismatch_count += 1

    # Get actual balance
    actual_balance, success = get_balance(w3, l2_token, claimer)

    if not success:
        print(f"    {YELLOW}[WARN]{RESET} Balance query failed: {claimer}")
        stats.balance_query_failures += 1
    elif actual_balance != amount:
        diff = abs(actual_balance - amount)
        print(f"    {RED}[MISMATCH]{RESET} {claimer}")
        print(f"      Snapshot: {amount}")
        print(f"      Actual:   {actual_balance}")
        print(f"      Diff:     {diff}")
        stats.mismatch_count += 1
    else:
        print(f"    {GREEN}[OK]{RESET} {claimer} | {amount}")

    stats.total_amount += amount


def verify_token(w3, token_data, token_index, stats):
    """Verify all claims for a single token"""
    l1_token = token_data['l1Token']
    l2_token = token_data['l2Token']
    token_name = token_data['tokenName']
    claims = token_data['data']

    print(f"\n{BLUE}[Token {token_index + 1}]{RESET} {token_name}")
    print(f"  L1: {l1_token}")
    print(f"  L2: {l2_token}")
    print(f"  Claims: {len(claims)}")

    stats.total_claims += len(claims)

    for claim in claims:
        verify_claim(w3, claim, l1_token, l2_token, stats)


def main():
    if len(sys.argv) != 3:
        print("Usage: python3 verify_assets.py <snapshot_file> <rpc_url>")
        print("\nExample:")
        print("  python3 verify_assets.py data/generate-assets-111551119090.json https://rpc.thanos-sepolia.tokamak.network")
        sys.exit(1)

    snapshot_file = sys.argv[1]
    rpc_url = sys.argv[2]

    # Connect to RPC
    try:
        w3 = Web3(Web3.HTTPProvider(rpc_url))
        if not w3.is_connected():
            print(f"{RED}Error: Could not connect to RPC: {rpc_url}{RESET}")
            sys.exit(1)
    except Exception as e:
        print(f"{RED}Error connecting to RPC: {e}{RESET}")
        sys.exit(1)

    # Get chain info
    chain_id = w3.eth.chain_id
    block_number = w3.eth.block_number

    # Load snapshot
    try:
        with open(snapshot_file, 'r') as f:
            snapshot_data = json.load(f)
    except FileNotFoundError:
        print(f"{RED}Error: Snapshot file not found: {snapshot_file}{RESET}")
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"{RED}Error: Invalid JSON in snapshot file: {e}{RESET}")
        sys.exit(1)

    # Print header
    print_header(snapshot_file, chain_id, block_number)

    # Verify all tokens
    stats = VerificationStats()
    stats.total_tokens = len(snapshot_data)

    for i, token_data in enumerate(snapshot_data):
        verify_token(w3, token_data, i, stats)

    # Print summary
    print_summary(stats)

    # Exit with appropriate code
    if stats.mismatch_count > 0 or stats.hash_mismatch_count > 0:
        sys.exit(1)
    else:
        sys.exit(0)


if __name__ == "__main__":
    main()
