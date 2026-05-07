# tokamak-deployer Registry

Each `{l1ChainId}.json` file lists pre-deployed L1 implementation addresses
that can be reused (skipping fresh deployment) when the user passes
`--reuse-deployment` to `deploy-contracts`.

## Schema

```json
{
  "tokamakDeployerVersion": "v0.0.6",
  "l1ChainId": 11155111,
  "comment": "Sepolia — implementations from <date>, <commit>",
  "implementations": {
    "SuperchainConfig":             "0x...",
    "OptimismPortal":               "0x...",
    "SystemConfig":                 "0x...",
    "L1StandardBridge":             "0x...",
    "L1CrossDomainMessenger":       "0x...",
    "OptimismMintableERC20Factory": "0x...",
    "L1ERC721Bridge":               "0x...",
    "L2OutputOracle":               "0x...",
    "DisputeGameFactory":           "0x..."
  }
}
```

## Update procedure

1. Run a fresh deploy without `--reuse-deployment` against the target L1.
2. Read the resulting `deploy-output.json` — its `implementations` field is the
   exact map shape required here.
3. Update this directory's `{l1ChainId}.json`.
4. Bump `tokamakDeployerVersion` if shipping a new release alongside the update.
5. Commit. The `//go:embed registry` directive in `cmd/registry.go` pulls these
   files into every binary built thereafter.

## Safety

- Bytecode of every listed address is verified at deploy preflight against the
  binary's embedded artifact `deployedBytecode.object`. Mismatched entries fall
  back to fresh deploy in default mode, or abort with `--reuse-strict`.
- Adding wrong addresses cannot silently corrupt a deployment — only delay it.

## Why pre-existing addresses are not always reusable

Addresses from prior deployments (e.g. the Foundry-era files under
`packages/tokamak/contracts-bedrock/deployments/<network>/address.json`) often
**cannot** be used as registry entries directly. The reuse mechanism compares
on-chain runtime bytecode keccak256-equal against the binary's embedded
`deployedBytecode.object`, which is regenerated each time the contracts are
recompiled. Any of the following invalidates a previously-deployed impl:

- Source revision change (added/removed functions, changed selectors)
- Compiler version change (solc 0.8.15 → 0.8.20)
- Optimizer setting change (runs=200 → runs=10000)
- Different metadata bytes (compile flags, source paths)

When this happens, every entry fails preflight and the deploy silently falls
back to fresh deploys — the registry becomes useless noise (9 WARN lines per
deploy) without any reuse benefit. Verify before populating:

```bash
# Per-entry check: on-chain code keccak256 vs artifact deployedBytecode keccak256
ONCHAIN=$(cast code <addr> --rpc-url <rpc>)
ARTIFACT=$(jq -r '.deployedBytecode.object' \
  cmd/tokamak-deployer/cmd/deploy-artifacts/<Name>.json)
[ "$(cast keccak "$ONCHAIN")" = "$(cast keccak "$ARTIFACT")" ] \
  && echo MATCH || echo MISMATCH
```

The reliable way to populate a registry is to run a fresh deploy with this
exact binary version (no `--reuse-deployment`), then copy the
`deploy-output.json:implementations` map verbatim. That guarantees byte-for-byte
equality on the next reuse-enabled deploy.
