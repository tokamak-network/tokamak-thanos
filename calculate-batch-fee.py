#!/usr/bin/env python3
"""
L1 Batch Transaction Fee Calculator
L1 트랜잭션 해시로부터 배치 비용 및 L2 트랜잭션 수수료 계산
Supports Fjord fee calculation
"""

import json
import sys
import subprocess
from typing import Dict, Optional

# RPC endpoints
L1_RPC = "http://localhost:8545"
L2_RPC = "http://localhost:9545"

# Fjord scalar values (from deploy config)
BASE_FEE_SCALAR = 1368
BLOB_BASE_FEE_SCALAR = 810949
DECIMALS = 6

# EIP-4844 constants
MIN_BLOB_BASE_FEE = 1  # wei
BLOB_BASE_FEE_UPDATE_FRACTION = 3338477


def rpc_call(endpoint: str, method: str, params: list) -> dict:
    """RPC 호출"""
    payload = {
        "jsonrpc": "2.0",
        "method": method,
        "params": params,
        "id": 1
    }

    result = subprocess.run(
        ["curl", "-s", endpoint, "-X", "POST", "-H", "Content-Type: application/json",
         "-d", json.dumps(payload)],
        capture_output=True,
        text=True
    )

    return json.loads(result.stdout)


def get_transaction(tx_hash: str) -> Optional[Dict]:
    """L1 트랜잭션 정보 조회"""
    response = rpc_call(L1_RPC, "eth_getTransactionByHash", [tx_hash])
    return response.get("result")


def get_receipt(tx_hash: str) -> Optional[Dict]:
    """L1 트랜잭션 receipt 조회"""
    response = rpc_call(L1_RPC, "eth_getTransactionReceipt", [tx_hash])
    return response.get("result")


def get_block(block_number: str) -> Optional[Dict]:
    """L1 블록 정보 조회"""
    response = rpc_call(L1_RPC, "eth_getBlockByNumber", [block_number, False])
    return response.get("result")


def calculate_blob_fee(receipt: Dict, block: Dict) -> Dict:
    """Blob 트랜잭션 수수료 계산"""
    blob_gas_used = int(receipt.get("blobGasUsed", "0x0"), 16)
    blob_gas_price = int(receipt.get("blobGasPrice", "0x0"), 16)

    blob_count = blob_gas_used // 131072 if blob_gas_used > 0 else 0
    blob_fee_wei = blob_gas_used * blob_gas_price
    blob_fee_gwei = blob_fee_wei / 1e9
    blob_fee_eth = blob_fee_wei / 1e18

    return {
        "type": "blob",
        "blob_count": blob_count,
        "blob_gas_used": blob_gas_used,
        "blob_gas_price_wei": blob_gas_price,
        "blob_gas_price_gwei": blob_gas_price / 1e9,
        "blob_fee_wei": blob_fee_wei,
        "blob_fee_gwei": blob_fee_gwei,
        "blob_fee_eth": blob_fee_eth,
    }


def calculate_calldata_fee(tx: Dict, block: Dict) -> Dict:
    """Calldata 트랜잭션 수수료 계산"""
    input_data = tx.get("input", "0x")
    data_size = (len(input_data) - 2) // 2  # hex string → bytes

    # Calldata gas: 16 gas per non-zero byte, 4 gas per zero byte
    # 간단화: 평균 16 gas/byte 사용
    calldata_gas = data_size * 16

    gas_price = int(tx.get("gasPrice", "0x0"), 16)
    calldata_fee_wei = calldata_gas * gas_price
    calldata_fee_gwei = calldata_fee_wei / 1e9
    calldata_fee_eth = calldata_fee_wei / 1e18

    return {
        "type": "calldata",
        "data_size_bytes": data_size,
        "data_size_kb": data_size / 1024,
        "calldata_gas": calldata_gas,
        "gas_price_wei": gas_price,
        "gas_price_gwei": gas_price / 1e9,
        "calldata_fee_wei": calldata_fee_wei,
        "calldata_fee_gwei": calldata_fee_gwei,
        "calldata_fee_eth": calldata_fee_eth,
    }


def calculate_total_fee(receipt: Dict) -> Dict:
    """전체 트랜잭션 수수료 계산"""
    gas_used = int(receipt.get("gasUsed", "0x0"), 16)
    effective_gas_price = int(receipt.get("effectiveGasPrice", "0x0"), 16)

    execution_fee_wei = gas_used * effective_gas_price
    execution_fee_gwei = execution_fee_wei / 1e9
    execution_fee_eth = execution_fee_wei / 1e18

    return {
        "gas_used": gas_used,
        "effective_gas_price_wei": effective_gas_price,
        "effective_gas_price_gwei": effective_gas_price / 1e9,
        "execution_fee_wei": execution_fee_wei,
        "execution_fee_gwei": execution_fee_gwei,
        "execution_fee_eth": execution_fee_eth,
    }


def fake_exponential(factor: int, numerator: int, denominator: int) -> int:
    """
    EIP-4844 fake exponential function
    Used to calculate blob base fee from excess blob gas
    """
    i = 1
    output = 0
    numerator_accum = factor * denominator
    while numerator_accum > 0:
        output += numerator_accum
        numerator_accum = (numerator_accum * numerator) // (denominator * i)
        i += 1
    return output // denominator


def get_blob_base_fee(excess_blob_gas: int) -> int:
    """
    Calculate blob base fee from excess blob gas (EIP-4844)

    Args:
        excess_blob_gas: Excess blob gas from block header

    Returns:
        Blob base fee in wei
    """
    return fake_exponential(
        MIN_BLOB_BASE_FEE,
        excess_blob_gas,
        BLOB_BASE_FEE_UPDATE_FRACTION
    )


def fjord_linear_regression(fast_lz_size: int) -> int:
    """
    Fjord Linear Regression to estimate Brotli size from FastLZ size
    Based on GasPriceOracle.sol _fjordLinearRegression
    """
    # Constants from Solidity
    INTERCEPT = -31_962_044
    COEFFICIENT = 1_033_413
    DIVISOR = 1_000_000
    MIN_SIZE = 100_000  # Minimum 100KB

    estimated = (INTERCEPT + COEFFICIENT * fast_lz_size) // DIVISOR
    return max(estimated, MIN_SIZE)


def calculate_fjord_l1_fee(tx_data_size: int, l1_base_fee: int, blob_base_fee: int, is_blob: bool) -> Dict:
    """
    Calculate L1 fee using Fjord formula
    Based on GasPriceOracle.sol _fjordL1Cost

    Args:
        tx_data_size: Compressed transaction data size (bytes)
        l1_base_fee: L1 base fee per gas (wei)
        blob_base_fee: Blob base fee per gas (wei) - 0 for calldata
        is_blob: Whether this is a blob transaction
    """
    # For Fjord, we estimate based on FastLZ compression
    # Assume FastLZ achieves ~50% of original size (rough estimate)
    # Since we have compressed size, multiply by 2 for original estimate
    fast_lz_size = tx_data_size * 2 if is_blob else tx_data_size

    # Estimate Brotli-compressed size using linear regression
    estimated_size = fjord_linear_regression(fast_lz_size)

    # Calculate fee: estimatedSize * (baseFeeScalar * 16 * l1BaseFee + blobBaseFeeScalar * blobBaseFee) / (10 ** (DECIMALS * 2))
    fee_scaled = BASE_FEE_SCALAR * 16 * l1_base_fee + BLOB_BASE_FEE_SCALAR * blob_base_fee

    # Use float division to preserve precision for small values
    fee_wei_precise = (estimated_size * fee_scaled) / (10 ** (DECIMALS * 2))
    fee_wei = int(fee_wei_precise)  # For display as integer wei

    return {
        "fast_lz_size_estimate": fast_lz_size,
        "brotli_size_estimate": estimated_size,
        "base_fee_scalar": BASE_FEE_SCALAR,
        "blob_base_fee_scalar": BLOB_BASE_FEE_SCALAR,
        "l1_base_fee_wei": l1_base_fee,
        "l1_base_fee_gwei": l1_base_fee / 1e9,
        "blob_base_fee_wei": blob_base_fee,
        "blob_base_fee_gwei": blob_base_fee / 1e9,
        "fjord_l1_fee_wei": fee_wei,
        "fjord_l1_fee_wei_precise": fee_wei_precise,  # Precise value with decimals
        "fjord_l1_fee_gwei": fee_wei_precise / 1e9,
        "fjord_l1_fee_eth": fee_wei_precise / 1e18,
    }




def print_results(tx: Dict, receipt: Dict, block: Dict):
    """결과 출력"""
    print("=" * 80)
    print("📊 L1 Batch Transaction Fee Analysis")
    print("=" * 80)
    print()

    # 기본 정보
    print(f"🔹 Transaction Hash: {tx['hash']}")
    print(f"🔹 Block Number: {int(receipt['blockNumber'], 16)}")
    print(f"🔹 From: {tx['from']}")
    print(f"🔹 To: {tx['to']}")
    print(f"🔹 Transaction Type: {tx.get('type', '0x0')}")
    print()

    # DA 타입 구분
    is_blob = tx.get("type") == "0x3"

    if is_blob:
        print("=" * 80)
        print("📦 Blob Transaction (EIP-4844)")
        print("=" * 80)
        blob_info = calculate_blob_fee(receipt, block)
        print(f"Blob Count: {blob_info['blob_count']}")
        print(f"Blob Gas Used: {blob_info['blob_gas_used']:,} gas")
        print(f"Blob Gas Price: {blob_info['blob_gas_price_wei']:,} wei")
        print(f"Blob Data Fee: {blob_info['blob_fee_gwei']:.6f} Gwei")
        print()
        data_fee_info = blob_info
    else:
        print("=" * 80)
        print("📄 Calldata Transaction")
        print("=" * 80)
        calldata_info = calculate_calldata_fee(tx, block)
        print(f"Data Size: {calldata_info['data_size_kb']:.2f} KB ({calldata_info['data_size_bytes']:,} bytes)")
        print(f"Calldata Gas: {calldata_info['calldata_gas']:,} gas")
        print(f"Gas Price: {calldata_info['gas_price_gwei']:.4f} Gwei")
        print(f"Calldata Fee: {calldata_info['calldata_fee_gwei']:.6f} Gwei ({calldata_info['calldata_fee_eth']:.9f} ETH)")
        print()
        data_fee_info = calldata_info

    # 전체 수수료
    print("=" * 80)
    print("💰 Execution Fee")
    print("=" * 80)
    total_fee = calculate_total_fee(receipt)
    print(f"Gas Used: {total_fee['gas_used']:,}")
    print(f"Effective Gas Price: {total_fee['effective_gas_price_gwei']:.4f} Gwei")
    print(f"Execution Fee: {total_fee['execution_fee_gwei']:.6f} Gwei")
    print()

    # L1 Data Fee 계산 (blob transaction인 경우만)
    if is_blob:
        print("=" * 80)
        print("💎 L1 Data Fee (Execution Fee + Blob Data Fee)")
        print("=" * 80)

        # Execution Fee = gasUsed * effectiveGasPrice
        execution_fee_wei = total_fee['execution_fee_wei']
        execution_fee_gwei = total_fee['execution_fee_gwei']
        execution_fee_eth = total_fee['execution_fee_eth']

        # Blob Data Fee = blobGasUsed * blobGasPrice
        blob_data_fee_wei = blob_info['blob_fee_wei']
        blob_data_fee_gwei = blob_info['blob_fee_gwei']
        blob_data_fee_eth = blob_info['blob_fee_eth']

        # L1 Data Fee = Execution Fee + Blob Data Fee
        l1_data_fee_wei = execution_fee_wei + blob_data_fee_wei
        l1_data_fee_gwei = execution_fee_gwei + blob_data_fee_gwei
        l1_data_fee_eth = execution_fee_eth + blob_data_fee_eth

        print(f"Components:")
        print(f"  - Execution Fee: {execution_fee_gwei:.6f} Gwei")
        print(f"  - Blob Data Fee: {blob_data_fee_gwei:.6f} Gwei")
        print()
        print(f"Total L1 Data Fee: {l1_data_fee_eth:.9f} ETH")
        print()

    print("=" * 80)


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 calculate-batch-fee.py <L1_TX_HASH>")
        print()
        print("Example:")
        print("  python3 calculate-batch-fee.py 0x83d5b70ff265eb49d3b5ad4862feecc19fddd9f7479b35b96b7b55bf319741a8")
        sys.exit(1)

    tx_hash = sys.argv[1]

    print(f"\n🔍 Fetching transaction {tx_hash}...\n")

    # 트랜잭션 정보 조회
    tx = get_transaction(tx_hash)
    if not tx:
        print(f"❌ Transaction not found: {tx_hash}")
        sys.exit(1)

    receipt = get_receipt(tx_hash)
    if not receipt:
        print(f"❌ Receipt not found: {tx_hash}")
        sys.exit(1)

    block = get_block(receipt["blockNumber"])
    if not block:
        print(f"❌ Block not found: {receipt['blockNumber']}")
        sys.exit(1)

    # 결과 출력
    print_results(tx, receipt, block)


if __name__ == "__main__":
    main()

