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
