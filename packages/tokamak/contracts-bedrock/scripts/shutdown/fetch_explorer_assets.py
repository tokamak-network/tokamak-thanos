import urllib.request
import urllib.parse
import json
import os
import time
import ssl
import sys

# LF: FFI compatibility - Log to stderr so stdout corresponds to ABI-encoded return values (if any)
def log(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)

# SSL Context
ctx = ssl.create_default_context()
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE

# 설정
BASE_V2_URL = "https://explorer.thanos-sepolia.tokamak.network/api/v2"
OUTPUT_DIR = "/Users/theo/workspace_tokamak/tokamak-thanos/packages/tokamak/contracts-bedrock/data"

# L2ToL1MessagePasser MessagePassed topic (from indexer bindings)
MESSAGE_PASSED_TOPIC = "0x02a52367d10742d8032712c1bb8e0144ff1ec5ffda1ed7d70bb05a2744955054"

# L1StandardBridge selectors (from forge artifacts)
FINALIZE_ERC20_SELECTOR = "a9f9e675"
FINALIZE_BRIDGE_ERC20_SELECTOR = "0166a07a"
FINALIZE_ETH_WITHDRAWAL_SELECTOR = "1532ec34"  # finalizeETHWithdrawal
FINALIZE_BRIDGE_ETH_SELECTOR = "1635f5fd"     # finalizeBridgeETH
FINALIZE_NATIVE_TOKEN_SELECTOR = "6580297d" # Thanos Native Token bridge selector
RELAY_MESSAGE_SELECTOR = "d764ad0b"         # Thanos Messenger selector

# OptimismPortal.finalizedWithdrawals(bytes32) selector
FINALIZED_WITHDRAWALS_SELECTOR = "a14238e7"

# L2 predeploy defaults
DEFAULT_L2_MESSAGE_PASSER = "0x4200000000000000000000000000000000000016"
DEFAULT_L2_ETH_TOKEN = "0x4200000000000000000000000000000000000486"
DEFAULT_L1_NATIVE_TOKEN = "0xa30fe40285B8f5c0457DbC3B7C8A280373c40044"

def fetch_v2_all_pages(description, endpoint, address_key="address"):
    """
    Explorer V2 API에서 모든 페이지의 데이터를 수집합니다.
    """
    all_addresses = set()
    next_page_params = {}

    log(f"🚀 {description} 수집 시작 (V2 API)...")

    while True:
        # Encode parameters
        query_string = urllib.parse.urlencode(next_page_params)
        url = f"{BASE_V2_URL}{endpoint}"
        if query_string:
            url += f"?{query_string}"

        try:
            req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
            with urllib.request.urlopen(req, context=ctx, timeout=20) as response:
                if response.getcode() != 200:
                    log(f"  ⚠️ HTTP Error: {response.getcode()}")
                    break

                data = json.loads(response.read().decode("utf-8"))
                items = data.get('items', [])

                if not items:
                    break

                for item in items:
                    # V2 API typically uses 'address' or 'hash' for address strings
                    # For tokens, 'address' is a string.
                    # For contracts, 'address' is a dict containing 'hash'.
                    # For addresses, 'hash' is a string.
                    addr = item.get(address_key)

                    if isinstance(addr, dict) and 'hash' in addr:
                        addr = addr['hash']

                    if not addr and 'hash' in item:
                        addr = item['hash'] # Common fallback for /addresses endpoint

                    if addr and isinstance(addr, str):
                        all_addresses.add(addr.lower())

                log(f"  📦 {len(items)}개 항목 처리됨. (누적 고유 주소: {len(all_addresses)})")

                # Pagination
                next_params = data.get('next_page_params')
                if not next_params:
                    break

                next_page_params = next_params
                time.sleep(0.1)  # Rate limit

        except Exception as e:
            log(f"  ❌ 오류 발생: {e}")
            break

    return sorted(list(all_addresses))

def fetch_token_holders_for_tokens(tokens):
    """
    Fetch holders for each token and return a union set.
    """
    all_holders = set()
    if not tokens:
        return []

    log("\n📌 토큰 홀더 목록 수집 시작")
    for idx, token in enumerate(tokens, start=1):
        endpoint = f"/tokens/{token}/holders"
        holders = fetch_v2_all_pages(
            f"Token holders ({idx}/{len(tokens)}) {token}",
            endpoint,
            address_key="address",
        )
        for h in holders:
            all_holders.add(h)

    log(f"✅ 토큰 홀더 수집 완료 (고유 주소: {len(all_holders)})")
    return sorted(list(all_holders))

def rpc_call(url, method, params):
    payload = json.dumps({
        "jsonrpc": "2.0",
        "id": 1,
        "method": method,
        "params": params,
    }).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=payload,
        headers={"Content-Type": "application/json"},
    )
    with urllib.request.urlopen(req, context=ctx, timeout=30) as response:
        data = json.loads(response.read().decode("utf-8"))
        if "error" in data:
            raise RuntimeError(data["error"])
        return data.get("result")

def to_hex_block(num):
    return hex(num)

def normalize_address(addr):
    if not addr:
        return None
    if addr.startswith("0x"):
        return addr.lower()
    return ("0x" + addr).lower()

def decode_address(word_hex):
    return "0x" + word_hex[-40:]

def parse_message_passed_data(entry):
    """
    MessagePassed(uint256 indexed nonce, address indexed sender, address indexed target,
                  uint256 value, uint256 gasLimit, bytes data, bytes32 withdrawalHash)
    """
    topics = entry.get("topics", [])
    if len(topics) < 4:
        return None, None, None

    # indexed params
    nonce = int(topics[1], 16)
    sender = decode_address(topics[2])
    target = decode_address(topics[3])

    data_hex = entry.get("data", "0x")
    if not data_hex or data_hex == "0x":
        return None, None, None

    raw = bytes.fromhex(data_hex[2:])
    # Bedrock MessagePassed data encoding:
    # [0:32] value
    # [32:64] gasLimit
    # [64:96] data offset
    # [96:128] withdrawalHash  <-- Correct position for the hash
    if len(raw) < 128:
        return None, None, None

    withdrawal_hash = "0x" + raw[96:128].hex()

    data_offset = int.from_bytes(raw[64:96], "big")
    if data_offset + 32 > len(raw):
        return None, None, None

    data_len = int.from_bytes(raw[data_offset:data_offset + 32], "big")
    data_start = data_offset + 32
    data_end = data_start + data_len
    if data_end > len(raw):
        # Fallback for some variations where length/data might be packed differently
        bridge_data = raw[data_start:]
    else:
        bridge_data = raw[data_start:data_end]

    return withdrawal_hash, bridge_data, target

def parse_bridge_call(data_bytes, default_l2_eth_token, bridge_addresses, l1_native_token):
    if not data_bytes or len(data_bytes) < 4:
        return None

    selector = data_bytes[0:4].hex()
    payload = data_bytes[4:]

    # Case 1: Nested in CrossDomainMessenger.relayMessage(uint256,address,address,uint256,uint256,bytes)
    # selector: 0xd764ad0b (Thanos) or 0xd1f3e792 (Standard Bedrock)
    if (selector == RELAY_MESSAGE_SELECTOR or selector == "d1f3e792") and len(payload) >= 32 * 6:
        # relayMessage args: uint256 _nonce, address _sender, address _target, uint256 _value, uint256 _minGasLimit, bytes _message
        inner_target = decode_address(payload[64:96].hex())
        # Check if the inner target is any of our bridges
        if any(inner_target.lower() == addr.lower() for addr in bridge_addresses):
            # Get the offset of the 'message' bytes (it's the 6th argument, so at 32 * 5 = 160)
            msg_offset = int.from_bytes(payload[160:192], "big")
            msg_len_start = msg_offset
            if len(payload) >= msg_len_start + 32:
                msg_len = int.from_bytes(payload[msg_len_start:msg_len_start + 32], "big")
                inner_message = payload[msg_len_start + 32:msg_len_start + 32 + msg_len]
                return parse_bridge_call(inner_message, default_l2_eth_token, bridge_addresses, l1_native_token)

    # Case 2: Direct L1StandardBridge calls or Standard Bridge finalizers
    if selector in (FINALIZE_ERC20_SELECTOR, FINALIZE_BRIDGE_ERC20_SELECTOR) and len(payload) >= 32 * 5:
        l1_token = decode_address(payload[0:32].hex())
        l2_token = decode_address(payload[32:64].hex())
        holder = decode_address(payload[96:128].hex()) # _to is the 4th param
        amount = int.from_bytes(payload[128:160], "big")
        return {
            "selector": selector,
            "l1Token": l1_token,
            "l2Token": l2_token,
            "holder": holder,
            "amount": amount,
        }

    if selector in (FINALIZE_ETH_WITHDRAWAL_SELECTOR, FINALIZE_BRIDGE_ETH_SELECTOR):
        # finalizeETHWithdrawal(address _from, address _to, uint256 _amount, bytes _extraData)
        if len(payload) >= 32 * 3:
            holder = decode_address(payload[32:64].hex()) # _to is the 2nd param
            amount = int.from_bytes(payload[64:96], "big")
            return {
                "selector": selector,
                "l1Token": "0x0000000000000000000000000000000000000000",
                "l2Token": normalize_address(default_l2_eth_token),
                "holder": holder,
                "amount": amount,
            }

    if selector == FINALIZE_NATIVE_TOKEN_SELECTOR:
        # finalizeBridgeNativeToken(address _from, address _to, uint256 _amount, bytes _extraData)
        if len(payload) >= 32 * 3:
            holder = decode_address(payload[32:64].hex()) # _to is the 2nd param
            amount = int.from_bytes(payload[64:96], "big")
            return {
                "selector": selector,
                "l1Token": normalize_address(l1_native_token),
                "l2Token": "0x0000000000000000000000000000000000000000",
                "holder": holder,
                "amount": amount,
            }

    return None

def is_finalized(l1_rpc_url, portal_address, withdrawal_hash):
    call_data = "0x" + FINALIZED_WITHDRAWALS_SELECTOR + withdrawal_hash[2:].rjust(64, "0")
    result = rpc_call(
        l1_rpc_url,
        "eth_call",
        [{"to": portal_address, "data": call_data}, "latest"],
    )
    return int(result, 16) == 1

def fetch_unclaimed_withdrawals(chain_id_suffix):
    l1_rpc_url = os.environ.get("L1_RPC_URL")
    l2_rpc_url = os.environ.get("L2_RPC_URL")
    l1_bridge = os.environ.get("BRIDGE_PROXY")
    optimism_portal = os.environ.get("OPTIMISM_PORTAL_PROXY")
    l1_usdc_bridge = os.environ.get("L1_USDC_BRIDGE_PROXY")

    if not l1_rpc_url or not l2_rpc_url or not l1_bridge or not optimism_portal:
        log("⚠️ unclaimed withdrawals: missing L1_RPC_URL/L2_RPC_URL/BRIDGE_PROXY/OPTIMISM_PORTAL_PROXY")
        return

    # All bridge addresses to scan (Standard + Custom)
    bridge_addresses = [l1_bridge]
    if l1_usdc_bridge:
        bridge_addresses.append(l1_usdc_bridge)
        log(f"  - Including USDC Bridge: {l1_usdc_bridge}")

    message_passer = normalize_address(os.environ.get("L2_TO_L1_MESSAGE_PASSER", DEFAULT_L2_MESSAGE_PASSER))
    l2_eth_token = normalize_address(os.environ.get("L2_ETH_TOKEN", DEFAULT_L2_ETH_TOKEN))
    l1_native_token = normalize_address(os.environ.get("L1_NATIVE_TOKEN", DEFAULT_L1_NATIVE_TOKEN))

    # Range configuration
    latest_hex = rpc_call(l2_rpc_url, "eth_blockNumber", [])
    latest_block = int(latest_hex, 16)

    # Default to scanning only recent blocks to find "unclaimed" ones quickly
    default_start = max(0, latest_block - 100000)
    start_block = int(os.environ.get("L2_START_BLOCK", str(default_start)))
    end_block_env = int(os.environ.get("L2_END_BLOCK", "0"))
    step = int(os.environ.get("L2_LOGS_STEP", "10000"))

    if end_block_env == 0:
        end_block = latest_block
    else:
        end_block = end_block_env

    log("\n📌 미완료 출금(MessagePassed) 수집 시작")
    log(f"  - L2 range: {start_block} ~ {end_block}")
    log(f"  - MessagePasser: {message_passer}")
    log(f"  - Filter: all MessagePassed events (scanning for bridge calls)")

    l1_tokens = []
    l2_tokens = []
    holders = []
    amounts = []
    withdrawal_hashes = []

    current = start_block
    while current <= end_block:
        to_block = min(current + step - 1, end_block)
        params = [{
            "address": message_passer,
            "fromBlock": to_hex_block(current),
            "toBlock": to_hex_block(to_block),
            "topics": [MESSAGE_PASSED_TOPIC], # Scan all message passing
        }]

        try:
            logs = rpc_call(l2_rpc_url, "eth_getLogs", params)
        except Exception as e:
            log(f"  ❌ eth_getLogs 실패: {e}")
            break

        for entry in logs:
            withdrawal_hash, bridge_data, target = parse_message_passed_data(entry)
            if not withdrawal_hash or not bridge_data:
                continue

            # Parse to see if it's a bridge call (possibly nested in Messenger)
            parsed = parse_bridge_call(bridge_data, l2_eth_token, bridge_addresses, l1_native_token)
            if not parsed:
                continue

            try:
                finalized = is_finalized(l1_rpc_url, optimism_portal, withdrawal_hash)
            except Exception as e:
                log(f"  ❌ finalizedWithdrawals 조회 실패: {e}")
                continue

            if finalized:
                continue

            l1_tokens.append(parsed["l1Token"])
            l2_tokens.append(parsed["l2Token"])
            holders.append(parsed["holder"])
            amounts.append(parsed["amount"])
            withdrawal_hashes.append(withdrawal_hash)
            log(f"  ✨ Found unclaimed: {withdrawal_hash[:10]}... Holder: {parsed['holder']} Amount: {parsed['amount']}")

        log(f"  - processed blocks: {current} ~ {to_block} (logs: {len(logs)})")
        current = to_block + 1

    items = []
    for i in range(len(withdrawal_hashes)):
        items.append({
            "withdrawalHash": withdrawal_hashes[i],
            "holder": holders[i],
            "l1Token": l1_tokens[i],
            "l2Token": l2_tokens[i],
            "amount": str(amounts[i]) # Store as string to avoid JSON overflow issues
        })

    filename = f"unclaimed-withdrawals{chain_id_suffix}.json"
    out_path = os.path.join(OUTPUT_DIR, filename)
    with open(out_path, "w") as f:
        json.dump(items, f, indent=2)

    log(f"✅ unclaimed withdrawals 저장: {filename} (count: {len(withdrawal_hashes)})")
    return filename, len(withdrawal_hashes)

def main():
    if not os.path.exists(OUTPUT_DIR):
        os.makedirs(OUTPUT_DIR)

    chain_id_suffix = ""
    if len(sys.argv) > 1:
        chain_id_suffix = f"-{sys.argv[1]}"

    log("====================================================")
    log("   Thanos L2 Explorer Asset Data Fetcher (V2)")
    log("====================================================\n")

    # 1. Accounts (Holders) - /addresses
    holders = fetch_v2_all_pages(
        "전체 계정(Holders) 목록",
        "/addresses",
        address_key="hash",
    )
    filename_holders = f"l2-holders{chain_id_suffix}.json"
    with open(os.path.join(OUTPUT_DIR, filename_holders), "w") as f:
        json.dump(holders, f, indent=2)

    # 2. Contracts (CA List) - /smart-contracts
    contracts = fetch_v2_all_pages("배포된 컨트랙트(CA) 목록", "/smart-contracts", address_key="address")
    filename_contracts = f"l2-contracts{chain_id_suffix}.json"
    with open(os.path.join(OUTPUT_DIR, filename_contracts), "w") as f:
        json.dump(contracts, f, indent=2)

    # 3. Tokens (ERC20 List) - /tokens
    tokens = fetch_v2_all_pages("ERC20 토큰 목록", "/tokens", address_key="address")
    filename_tokens = f"l2-tokens{chain_id_suffix}.json"
    with open(os.path.join(OUTPUT_DIR, filename_tokens), "w") as f:
        json.dump(tokens, f, indent=2)

    # 4. Token Holders (union)
    token_holders = fetch_token_holders_for_tokens(tokens)
    merged_holders = sorted(set(holders) | set(token_holders))

    # Overwrite holders file with merged holders
    with open(os.path.join(OUTPUT_DIR, filename_holders), "w") as f:
        json.dump(merged_holders, f, indent=2)

    # 5. Unclaimed withdrawals (MessagePassed -> OptimismPortal.finalizedWithdrawals)
    filename_unclaimed, unclaimed_count = fetch_unclaimed_withdrawals(chain_id_suffix)

    log("\n====================================================")
    log("✅ 데이터 수집 완료 및 파일 저장 성공!")
    log(f"📂 경로: {OUTPUT_DIR}")
    log(f"  - {filename_holders:<20} : {len(merged_holders):>5} 주소 (accounts + token holders)")
    log(f"  - {filename_contracts:<20} : {len(contracts):>5} 주소 (스마트 컨트랙트)")
    log(f"  - {filename_tokens:<20} : {len(tokens):>5} 주소 (ERC20 토큰)")
    log(f"  - {filename_unclaimed:<20} : {unclaimed_count:>5} 항목 (미완료 출금)")
    log("====================================================")

    # FFI requirement: output hex string to stdout
    print("0x")

if __name__ == "__main__":
    main()
