# Optimism v1.16.0 Migration Guide

## Overview

This guide documents the migration process from Optimism v1.7.7 to v1.16.0 for the Tokamak-Thanos codebase. This migration was necessary to support the op-geth upgrade from v1.101315.2 to v1.101601.0-rc.1.

## Migration Context

**Trigger**: op-geth upgrade to v1.101601.0-rc.1
**Source Version**: Optimism v1.7.7
**Target Version**: Optimism v1.16.0
**Date**: 2025-11-04

## Prerequisites

Before starting the migration:

1. Ensure you have the source Optimism v1.16.0 repository available:
   ```bash
   /Users/zena/tokamak-projects/optimism
   ```

2. Required tools:
   - Go 1.21+
   - Foundry (forge, cast)
   - Python 3 (for devnet-allocs)
   - sed (for batch import path updates)

## Migration Steps

### 1. Identify Required Packages

The migration requires updating approximately 90+ packages. Key areas include:

- **Core Challenger Packages**:
  - `op-challenger/game/fault/types`
  - `op-challenger/game/types`
  - `op-challenger/metrics`

- **Chain Operations**:
  - `op-chain-ops/crossdomain`
  - `op-chain-ops/genesis`

- **E2E Testing**:
  - `op-e2e/actions/proofs`
  - `op-e2e/e2eutils/disputegame`

- **Proposer**:
  - `op-proposer/proposer`
  - `op-proposer/flags`

### 2. Copy Required Packages

Use the following approach to copy packages from Optimism v1.16.0:

```bash
# Example: Copy a package directory
cp -r /Users/zena/tokamak-projects/optimism/op-e2e/actions/proofs \
      /Users/zena/tokamak-projects/tokamak-thanos/op-e2e/actions/

# Fix import paths automatically
find /Users/zena/tokamak-projects/tokamak-thanos/op-e2e/actions/proofs \
     -name "*.go" \
     -exec sed -i '' 's|github.com/ethereum-optimism/optimism|github.com/tokamak-network/tokamak-thanos|g' {} \;
```

### 3. Update Type Definitions

#### GameType Enumeration

Add new GameType constants to `op-challenger/game/types/types.go`:

```go
const (
    GameTypeAlphabet            GameType = 255
    GameTypeFast                GameType = 254
    GameTypeCannon              GameType = 0
    GameTypePermissionedCannon  GameType = 1
    GameTypeAsterisc            GameType = 2
    GameTypeAsteriscKona        GameType = 3
    GameTypeSuperCannon         GameType = 5
    GameTypeSuperPermissioned   GameType = 6
    GameTypeSuperAsteriscKona   GameType = 7
    GameTypeOPSuccinct          GameType = 100
    GameTypeKailua              GameType = 101
    GameTypeFast                GameType = 102
)
```

#### BondDistributionMode

Add bond distribution mode enumeration:

```go
type BondDistributionMode uint8

const (
    BondDistributionModeUnknown BondDistributionMode = iota
    BondDistributionModeIncentiveSystem
    BondDistributionModeRewardSystem
)
```

### 4. Update Interfaces

#### Metrics Interface

Update `op-challenger/metrics/noop.go` to implement the new `RecordLargePreimageCount` method:

```go
func (*NoopMetricsImpl) RecordLargePreimageCount(count int) {}
```

### 5. Fix Import Path Conflicts

**Problem**: TypesWithdrawalTransaction type mismatch between `op-bindings/bindings` and `op-node/bindings`.

**Solution**: Update import in `op-chain-ops/crossdomain/withdrawal.go`:

```go
// Before
import "github.com/tokamak-network/tokamak-thanos/op-bindings/bindings"

// After
import "github.com/tokamak-network/tokamak-thanos/op-node/bindings"
```

### 6. Update Struct Definitions

#### PreimageOracleData

Update the preimage oracle data structure:

```go
// v1.7.7
type PreimageOracleData struct {
    BlobFieldIndex uint64
    // ...
}

// v1.16.0
type PreimageOracleData struct {
    ZPoint [32]byte  // Replaced BlobFieldIndex
    // ...
}
```

#### Add DecodeClock Function

Add clock decoding utility:

```go
func DecodeClock(clock uint128) types.Clock {
    duration := uint64(clock >> 64)
    timestamp := uint64(clock & 0xFFFFFFFFFFFFFFFF)
    return types.Clock{
        Duration:  time.Duration(duration) * time.Second,
        Timestamp: uint64(timestamp),
    }
}
```

### 7. Update op-node Command Files

The op-node command structure changed in v1.16.0. Copy updated files:

```bash
# Copy genesis command files
cp /path/to/optimism/op-node/cmd/genesis/cmd.go \
   op-node/cmd/genesis/
cp /path/to/optimism/op-node/cmd/genesis/systemconfig.go \
   op-node/cmd/genesis/

# Copy updated main.go
cp /path/to/optimism/op-node/cmd/main.go \
   op-node/cmd/

# Copy new interop command
cp -r /path/to/optimism/op-node/cmd/interop \
      op-node/cmd/

# Fix import paths
find op-node/cmd/ -name "*.go" \
  -exec sed -i '' 's|github.com/ethereum-optimism/optimism|github.com/tokamak-network/tokamak-thanos|g' {} \;
```

**Key Changes**:
- `config.Check()` now requires `log.Logger` parameter
- `genesis.BuildL2Genesis` signature updated
- `node.New` call signature changed (removed snapshotLog parameter)
- New `interop` command directory added

### 8. Update Contracts Path in E2E Tests

Tokamak-Thanos uses its own contracts located at `packages/tokamak/contracts-bedrock` instead of Optimism's `packages/contracts-bedrock`. Update the E2E test files to use the correct path:

```bash
# Update all contracts-bedrock references in E2E tests
find op-e2e/faultproofs -name "*.go" -exec sed -i '' 's|packages/contracts-bedrock|packages/tokamak/contracts-bedrock|g' {} \;
find op-e2e/interop -name "*.go" -exec sed -i '' 's|packages/contracts-bedrock|packages/tokamak/contracts-bedrock|g' {} \;
find op-e2e/system/contracts -name "*.go" -exec sed -i '' 's|packages/contracts-bedrock|packages/tokamak/contracts-bedrock|g' {} \;

# Verify the changes
grep -r "packages/contracts-bedrock" op-e2e/ --include="*.go"
# Should return no results
```

**Updated Files**:
- `op-e2e/faultproofs/util_interop.go`
- `op-e2e/interop/interop_recipe_test.go`
- `op-e2e/interop/interop_test.go`
- `op-e2e/system/contracts/artifactsfs_test.go`

**Note**: The Tokamak contracts at `packages/tokamak/contracts-bedrock` should never be modified during migration. They are maintained separately from Optimism contracts.

### 9. Handle RAT Test Files

RAT (Resolve Attestation Token) tests require bindings not yet migrated. Temporarily exclude these files:

```bash
mv op-e2e/faultproofs/rat_e2e_test.go \
   op-e2e/faultproofs/rat_e2e_test.go.bak

mv op-e2e/faultproofs/rat_simple_test.go \
   op-e2e/faultproofs/rat_simple_test.go.bak

mv op-e2e/faultproofs/rat_unit_test.go \
   op-e2e/faultproofs/rat_unit_test.go.bak
```

### 10. Build Solidity Artifacts

E2E tests require Solidity contract artifacts. Note that Tokamak-Thanos contracts are already pre-built and located at `packages/tokamak/contracts-bedrock/`:

```bash
# Verify Tokamak contracts artifacts exist
ls packages/tokamak/contracts-bedrock/forge-artifacts/ | head -10
ls packages/tokamak/contracts-bedrock/deployments/31337-deploy.json
```

These artifacts are maintained separately and should not be regenerated during this migration.

**Note**: Unlike Optimism's migration process which requires running `make devnet-allocs`, Tokamak-Thanos uses pre-existing contract artifacts that are part of the repository.

### 11. Run Tests

After migration, compile and run tests:

```bash
# Compile the test package
go build ./op-e2e/faultproofs

# Run specific E2E test
go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./op-e2e/faultproofs
```

## Common Issues and Solutions

### Issue 1: Missing Package Dependencies

**Error**: `no required module provides package X`

**Solution**: Copy the missing package from Optimism v1.16.0 and fix import paths using sed.

### Issue 2: Interface Method Mismatch

**Error**: `does not implement Interface (missing method X)`

**Solution**: Check the interface definition in v1.16.0 and add the missing method to your implementation.

### Issue 3: Type Compatibility

**Error**: `cannot use value of type X as type Y`

**Solution**: Verify the type definitions match v1.16.0. Common issues:
- `TypesWithdrawalTransaction` import path changes
- Struct field changes (e.g., `BlobFieldIndex` → `ZPoint`)

### Issue 4: Missing Solidity Artifacts

**Error**: `failed to open artifact "DeployImplementations.s.sol/DeployImplementations.json"`

**Solution**: Run `make devnet-allocs` to generate all required artifacts.

## Verification

After completing the migration:

1. **Compilation Check**:
   ```bash
   go build ./op-e2e/faultproofs
   ```
   Should complete without errors.

2. **Test Execution**:
   ```bash
   go test -v -run TestOutputAlphabetGame_ReclaimBond -timeout 10m ./op-e2e/faultproofs
   ```
   Should compile and start executing (runtime errors related to missing devnet are expected in development environment).

3. **Import Path Consistency**:
   ```bash
   # Verify no ethereum-optimism imports remain
   grep -r "github.com/ethereum-optimism/optimism" --include="*.go" .
   ```
   Should return no results in migrated packages.

## Rollback Plan

If issues arise:

1. **Restore RAT Test Files**:
   ```bash
   mv op-e2e/faultproofs/rat_e2e_test.go.bak \
      op-e2e/faultproofs/rat_e2e_test.go
   ```

2. **Revert Package Changes**:
   Use git to revert specific package directories:
   ```bash
   git checkout HEAD -- op-challenger/
   ```

3. **Clean Build Artifacts**:
   ```bash
   make clean
   rm -rf packages/tokamak/contracts-bedrock/forge-artifacts/
   ```

## Next Steps

After successful migration:

1. Complete RAT bindings migration when available
2. Update remaining test files to v1.16.0 patterns
3. Review and update documentation
4. Run full test suite: `make test`

## References

- [Optimism v1.16.0 Release Notes](https://github.com/ethereum-optimism/optimism/releases/tag/v1.16.0)
- [op-geth v1.101601.0-rc.1](https://github.com/ethereum-optimism/op-geth/releases/tag/v1.101601.0-rc.1)
- [Tokamak-Thanos Repository](https://github.com/tokamak-network/tokamak-thanos)

## Support

For issues or questions:
- Create an issue in the Tokamak-Thanos repository
- Refer to `CHANGELOG.md` for detailed list of changes
- See `TESTING-GUIDE.md` for E2E testing instructions
