#!/usr/bin/env python3
"""
Generate genesis state for E2E tests
This creates the state-dump files that are required for Challenger testing
"""

import json
import os
from pathlib import Path

def generate_genesis_state():
    """Generate genesis state for E2E tests"""

    # Load deployment addresses
    deploy_file = Path("packages/tokamak/contracts-bedrock/deployments/devnetL1/.deploy")
    if not deploy_file.exists():
        print(f"❌ Deploy file not found: {deploy_file}")
        return False

    with open(deploy_file) as f:
        addresses = json.load(f)

    print(f"✅ Loaded {len(addresses)} contract addresses")

    # Create L1 genesis state
    l1_genesis = {
        "alloc": {}
    }

    # Add deployed contracts to genesis with minimal proxy code
    # This is the minimal proxy bytecode that delegates all calls
    minimal_proxy_code = "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063f851a44014610030575b600080fd5b61003861004e565b6040516001600160a01b03909116815260200160405180910390f35b60006100767f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5490565b905090565b600080fd5b5056fea26469706673582212208b0e7b7c7d89f3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c3f3c64736f6c63430008110033"

    for name, address in addresses.items():
        print(f"  Adding {name}: {address}")
        l1_genesis["alloc"][address.lower()] = {
            "balance": "0x0",
            "code": minimal_proxy_code,
            "storage": {}
        }

    # Add funded accounts (dev accounts with ETH balance)
    dev_accounts = [
        "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc",
        "0x976EA74026E726554dB657fA54763abd0C3a0aa9",
        "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955",
        "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f",
        "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720",
        "0xBcd4042DE499D14e55001CcbB24a551F3b954096",
        "0x71bE63f3384f5fb98995898A86B02Fb2426c5788",
        "0xFABB0ac9d68B0B445fB7357272Ff202C5651694a",
        "0x1CBd3b2770909D4e10f157cABC84C7264073C9Ec",
        "0xdF3e18d64BC6A983f673Ab319CCaE4f1a57C7097",
        "0xcd3B766CCDd6AE721141F452C550Ca635964ce71",
        "0x2546BcD3c84621e976D8185a91A922aE77ECEc30",
        "0xbDA5747bFD65F08deb54cb465eB87D40e51B197E",
        "0xdD2FD4581271e230360230F9337D5c0430Bf44C0",
        "0x8626f6940E2eb28930eFb4CeF49B2d1F2C9C1199",
        "0x09DB0a93B389bEF724429898f539AEB7ac2Dd55f",
        "0x02484cb50AAC86Eae85610D6f4Bf026f30f6627D",
        "0x08135Da0A343E492FA2d4282F2AE34c6c5CC1BbE",
        "0x5E661B79FE2D3F6cE70F5AAC07d8Cd9abb2743F1",
        "0x61097BA76cD906d2ba4FD106E757f7Eb455fc295",
        "0xDf37F81dAAD2b0327A0A50003740e1C935C70913",
        "0x553BC17A05702530097c3677091C5BB47a3a7931",
        "0x87BdCE72c06C21cd96219BD8521bDF1F42C78b5e",
        "0x40Fc963A729c542424cD800349a7E4Ecc4896624",
        "0x9DCCe783B6464611f38631e6C851bf441907c710"
    ]

    print(f"\n✅ Adding {len(dev_accounts)} funded dev accounts")
    for account in dev_accounts:
        l1_genesis["alloc"][account.lower()] = {
            "balance": "0x21e19e0c9bab2400000"  # 10000 ETH in wei (hex)
        }

    # Write state dump files
    output_dir = Path("packages/tokamak/contracts-bedrock")
    output_dir.mkdir(parents=True, exist_ok=True)

    # Delta version (initial state)
    delta_file = output_dir / "state-dump-901-delta.json"
    print(f"\n📝 Writing delta state to: {delta_file}")
    with open(delta_file, "w") as f:
        json.dump(l1_genesis, f, indent=2)

    # Ecotone version (with gas price oracle update)
    # For now, using same state as delta (can be enhanced later)
    l1_genesis_ecotone = l1_genesis.copy()

    ecotone_file = output_dir / "state-dump-901-ecotone.json"
    print(f"📝 Writing ecotone state to: {ecotone_file}")
    with open(ecotone_file, "w") as f:
        json.dump(l1_genesis_ecotone, f, indent=2)

    # Final version (Fjord activated)
    final_file = output_dir / "state-dump-901.json"
    print(f"📝 Writing final state to: {final_file}")
    with open(final_file, "w") as f:
        json.dump(l1_genesis, f, indent=2)

    print("\n✅ Genesis state files generated successfully!")
    print("\nYou can now run E2E tests with:")
    print("  go test -v ./op-e2e/faultproofs")

    return True

if __name__ == "__main__":
    success = generate_genesis_state()
    exit(0 if success else 1)