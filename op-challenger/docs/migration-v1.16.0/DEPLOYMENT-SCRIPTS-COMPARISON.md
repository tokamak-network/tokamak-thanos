# Deployment Scripts Comparison: Optimism vs Tokamak-Thanos

**Date**: 2025-11-04
**Optimism Version**: v1.16.0
**Tokamak-Thanos**: Based on v1.7.7 (migrating to v1.16.0)

## Executive Summary

Optimism v1.16.0 introduces a significantly expanded deployment system with 33 Solidity deployment scripts totaling 6,525 lines of code, while Tokamak-Thanos currently maintains only 1 deployment script with 57 lines. This document analyzes the differences and provides recommendations for migration.

## Directory Structure Comparison

### Optimism v1.16.0
**Location**: `packages/contracts-bedrock/scripts/deploy/`
**Total Files**: 34 files (33 `.sol` + 1 `.sh`)
**Total Lines**: 6,525 lines of Solidity code

### Tokamak-Thanos
**Location**: `packages/tokamak/contracts-bedrock/scripts/deploy/`
**Total Files**: 1 file
**Total Lines**: 57 lines of Solidity code

## Detailed File Comparison

| File Name | Optimism Lines | Tokamak Lines | Status | Purpose |
|-----------|----------------|---------------|--------|---------|
| **DeployAsterisc.s.sol** | 57 | 57 | ✅ **Present** | Asterisc VM deployment |
| **DeployImplementations.s.sol** | 674 | ❌ Missing | 🔴 **Critical** | Core implementations deployment (NEW in v1.16.0) |
| **DeployOPChain.s.sol** | 733 | ❌ Missing | 🔴 **Critical** | OP Chain deployment |
| **VerifyOPCM.s.sol** | 943 | ❌ Missing | 🟡 Important | OPCM verification |
| **DeploySuperchain.s.sol** | 299 | ❌ Missing | 🟡 Important | Superchain setup |
| **DeployOwnership.s.sol** | 347 | ❌ Missing | 🟡 Important | Ownership management |
| **DeployConfig.s.sol** | 300 | ❌ Missing | 🟡 Important | Configuration management |
| **ChainAssertions.sol** | 537 | ❌ Missing | 🟡 Important | Chain state validation |
| **Deploy.s.sol** | 523 | ❌ Missing | 🟡 Important | Main deployment entry point |
| **DeployDisputeGame.s.sol** | 314 | ❌ Missing | 🟢 Optional | Dispute game deployment |
| **InteropMigration.s.sol** | 240 | ❌ Missing | 🟢 Optional | Interoperability migration |
| **DeployMIPS.s.sol** | 105 | ❌ Missing | 🟢 Optional | MIPS VM deployment |
| **AddGameType.s.sol** | 106 | ❌ Missing | 🟢 Optional | Game type addition |
| Others (20 files) | ~2,387 | ❌ Missing | 🟢 Optional | Various utilities |

## Critical Missing Components

### 1. **DeployImplementations.s.sol** (674 lines) ⚠️ CRITICAL

**Purpose**: Deploys all L1 contract implementations for the OP Stack

**Key Components**:
```solidity
- IOPContractsManager opcm
- IDelayedWETH delayedWETHImpl
- IOptimismPortal optimismPortalImpl
- ISystemConfig systemConfigImpl
- IL1CrossDomainMessenger l1CrossDomainMessengerImpl
- IL1StandardBridge l1StandardBridgeImpl
- IDisputeGameFactory disputeGameFactoryImpl
- IAnchorStateRegistry anchorStateRegistryImpl
- IRAT ratImpl (RAT implementation)
- ISuperchainConfig superchainConfigImpl
- IProtocolVersions protocolVersionsImpl
```

**Impact on Migration**:
- **Required for E2E tests**: The `op-deployer` package loads this script to deploy test environments
- **Artifact dependency**: `op-e2e/config/init.go` references `DeployImplementations.s.sol/DeployImplementations.json`
- **Current workaround**: Copied compiled artifacts from Optimism to `packages/tokamak/contracts-bedrock/forge-artifacts/`

**Action Required**:
- ✅ Artifacts copied (temporary solution)
- ⚠️ Source script migration needed for production deployments

### 2. **DeployOPChain.s.sol** (733 lines) ⚠️ CRITICAL

**Purpose**: Deploys an individual OP Chain with all necessary contracts

**Key Components**:
- Chain-specific deployment logic
- Proxy contract deployment
- Implementation contract initialization
- Integration with OPCM (OP Contracts Manager)

**Action Required**: Analyze if Tokamak uses custom chain deployment logic

### 3. **Deploy.s.sol** (523 lines) 🟡 IMPORTANT

**Purpose**: Main deployment orchestration script

**Integration**: Coordinates DeployImplementations, DeployOPChain, and DeploySuperchain

**Action Required**: Review Tokamak's deployment workflow

## New v1.16.0 Architecture: OPCM (OP Contracts Manager)

### What is OPCM?

OPCM is a new contract management system introduced in Optimism v1.16.0 that standardizes OP Stack deployments.

**Related Scripts**:
- `DeployImplementations.s.sol` - Deploys OPCM and implementations
- `VerifyOPCM.s.sol` - Validates OPCM deployment
- `GenerateOPCMMigrateCalldata.sol` - Migration utilities
- `AddGameType.s.sol` - Adds dispute game types to OPCM

**OPCM Components**:
```solidity
IOPContractsManager
IOPContractsManagerGameTypeAdder
IOPContractsManagerDeployer
IOPContractsManagerUpgrader
IOPContractsManagerInteropMigrator
IOPContractsManagerStandardValidator
IOPContractsManagerContractsContainer
```

## Tokamak-Specific Deployment Strategy

### Current Approach

Tokamak-Thanos appears to use:
1. **Minimal deployment scripts**: Only `DeployAsterisc.s.sol` is present
2. **Python-based deployment**: `make devnet-allocs` uses Python scripts in `bedrock-devnet/`
3. **Pre-built artifacts**: Contracts are pre-compiled and artifacts are versioned

### Advantages

1. **Simplified maintenance**: Fewer scripts to update
2. **Custom workflow**: Tailored to Tokamak's specific needs
3. **Version control**: Artifacts are committed and versioned

### Challenges for v1.16.0 Migration

1. **op-deployer dependency**: v1.16.0 E2E tests require `DeployImplementations.s.sol`
2. **OPCM integration**: New architecture may conflict with custom deployment
3. **Artifact compatibility**: Need to ensure artifacts match v1.16.0 interfaces

## Migration Recommendations

### Phase 1: Immediate Actions (Current Status)

✅ **Completed**:
- Copied `DeployImplementations.s.sol` artifacts to `forge-artifacts/`
- Updated `op-e2e/config/init.go` to use Tokamak's contracts path
- E2E tests can now find required artifacts

### Phase 2: Short-term (For v1.16.0 Compatibility)

**Option A: Minimal Approach** (Recommended)
1. Keep Tokamak's custom deployment workflow
2. Maintain copied artifacts for E2E testing only
3. Add artifacts to `.gitignore` as they're reference-only

**Option B: Partial Integration**
1. Copy only critical scripts:
   - `DeployImplementations.s.sol`
   - `DeployOPChain.s.sol`
   - `ChainAssertions.sol`
2. Modify scripts to work with Tokamak's customizations
3. Use for E2E testing, not production deployments

### Phase 3: Long-term (Future Consideration)

**Option C: Full OPCM Integration**
1. Adopt Optimism's full deployment system
2. Implement Tokamak customizations within OPCM framework
3. Migrate from Python deployment to Solidity-based deployment
4. **Estimated effort**: 2-4 weeks development + testing

## Specific Recommendations

### For Current Migration (v1.16.0)

1. **Keep current artifacts approach**: ✅ Already implemented
   - Artifacts copied to `packages/tokamak/contracts-bedrock/forge-artifacts/`
   - No source code changes to Tokamak contracts
   - E2E tests work with compiled artifacts

2. **Document the difference**:
   ```
   Tokamak-Thanos uses pre-built contract artifacts and custom Python-based
   deployment scripts, while Optimism v1.16.0 uses Solidity deployment scripts
   with the OPCM system. For E2E testing compatibility, Optimism's compiled
   DeployImplementations artifacts are used as reference-only files.
   ```

3. **Add to `.gitignore`** (if artifacts change frequently):
   ```
   # Optimism reference artifacts for E2E testing (do not modify)
   packages/tokamak/contracts-bedrock/forge-artifacts/DeployImplementations.s.sol/
   ```

### For Future Development

1. **Monitor OPCM evolution**: Optimism is actively developing OPCM
2. **Evaluate integration timeline**: Consider OPCM when planning next major version
3. **Maintain compatibility layer**: Keep artifacts updated with Optimism releases

## Risk Assessment

| Risk | Severity | Mitigation |
|------|----------|------------|
| **Artifact version mismatch** | 🔴 High | Automated CI checks, version tracking |
| **OPCM API changes** | 🟡 Medium | Monitor Optimism releases, update artifacts |
| **Tokamak customizations lost** | 🟡 Medium | Keep separate deployment workflow |
| **E2E test failures** | 🟢 Low | Artifacts are reference-only, tests isolated |

## Action Items

### Immediate (Current Migration)
- [x] Copy DeployImplementations artifacts
- [x] Update op-e2e/config/init.go paths
- [ ] Test E2E suite with new artifacts
- [ ] Document artifact sourcing in README

### Short-term (Next Month)
- [ ] Create artifact update process
- [ ] Set up CI to validate artifact compatibility
- [ ] Review Tokamak deployment vs. OPCM differences

### Long-term (Next Quarter)
- [ ] Evaluate OPCM integration benefits
- [ ] Plan migration strategy if beneficial
- [ ] Prototype OPCM integration in test environment

## Conclusion

Tokamak-Thanos and Optimism v1.16.0 have significantly different deployment architectures. For the current migration:

1. **Keep Tokamak's deployment system**: It works well and is battle-tested
2. **Use Optimism artifacts for E2E tests only**: Minimal integration approach
3. **Monitor OPCM development**: Consider integration in future major version
4. **Document the architectural difference**: Ensure future maintainers understand the separation

The current approach (copied artifacts) is pragmatic and low-risk, allowing v1.16.0 migration without disrupting Tokamak's proven deployment workflow.

## References

- [Optimism DeployImplementations.s.sol](https://github.com/ethereum-optimism/optimism/blob/v1.16.0/packages/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol)
- [OPCM Documentation](https://docs.optimism.io/builders/chain-operators/tools/opcm)
- Tokamak-Thanos: `bedrock-devnet/main.py` (custom deployment)
- Migration Guide: `./MIGRATION-GUIDE.md`

---

**Document Version**: 1.0
**Last Updated**: 2025-11-04
**Next Review**: After v1.16.0 migration completion
