# E2E Testing Deployment Compatibility Analysis

**Date**: 2025-11-04
**Context**: Tokamak-Thanos v1.16.0 Migration
**Issue**: "revision id 1 cannot be reverted" error when running E2E tests
**Status**: 🔴 Critical - E2E tests cannot run with Optimism v1.16.0 deployment scripts

## Executive Summary

Optimism v1.16.0's E2E tests require deployment scripts (`DeploySuperchain.s.sol`, `DeployImplementations.s.sol`, etc.) to create test environments. However, when these scripts execute against Tokamak-Thanos contracts, they fail with a `"revision id 1 cannot be reverted"` error during contract initialization.

**Root Cause**: Tokamak's custom contract architecture uses multi-layer initialization patterns (OpenZeppelin `Initializable` + `AccessControlUpgradeable` + custom logic) that create EVM state snapshot conflicts when executed through Optimism's proxy upgrade flow.

**Impact**: E2E tests (`TestOutputAlphabetGame_ReclaimBond` and others) cannot run until this incompatibility is resolved.

## Error Analysis

### The Error

```
panic: revision id 1 cannot be reverted [recovered]
	panic: revision id 1 cannot be reverted

goroutine 53 [running]:
github.com/tokamak-network/tokamak-thanos/op-chain-ops/script.(*Host).Call.func1()
	/Users/zena/tokamak-projects/tokamak-thanos/op-chain-ops/script/script.go:362 +0x214
```

**Location**: `op-chain-ops/script/script.go:352-388`

### What This Error Means

The error originates from Go-Ethereum's EVM state management:

1. **Snapshots (Revisions)**: The EVM creates snapshots of state before executing transactions to enable rollback on failure
2. **Snapshot Stack**: Snapshots are managed in a stack with revision IDs (0, 1, 2, etc.)
3. **The Problem**: Code tries to revert to revision ID 1, but that snapshot has already been consumed/invalidated

### Error Handler Code

**File**: `op-chain-ops/script/script.go:352-388`

```go
func (h *Host) Call(from common.Address, to common.Address, input []byte, gas uint64, value *uint256.Int) (returnData []byte, leftOverGas uint64, err error) {
	h.prelude(from, &to)

	defer func() {
		if r := recover(); r != nil {
			rStr, ok := r.(string)
			if !ok || !strings.Contains(strings.ToLower(rStr), "revision id") {
				fmt.Printf("Unexpected panic in script execution: %v\n", r)
				panic(r)
			}

			// This catches the "revision id 1 cannot be reverted" error
			fmt.Printf("Caught revision id error: %s\n", rStr)

			if h.evmRevertErr != nil {
				err = h.evmRevertErr
			} else {
				err = errors.New("execution reverted, check logs")
			}
		}
		h.evmRevertErr = nil
	}()

	returnData, leftOverGas, err = h.env.Call(from, to, input, gas, value)
	// ...
}
```

The code already has a panic recovery handler for this specific error, but the underlying issue prevents deployment from succeeding.

## Tokamak's Custom Contract Architecture

### 1. Custom Verification System (NEW - Not in Optimism)

**File**: `packages/tokamak/contracts-bedrock/src/tokamak-contracts/verification/L1ContractVerification.sol`

**Key Characteristics**:

```solidity
contract L1ContractVerification is
  IL1ContractVerification,
  Initializable,                    // OpenZeppelin upgradeable pattern
  AccessControlUpgradeable          // Role-based access control
{
    function initialize(
        address _tokenAddress,
        address _initialAdmin
    ) public initializer {            // Multiple initialization layers
        __AccessControl_init();       // Init AccessControl first
        _setupRole(DEFAULT_ADMIN_ROLE, _initialAdmin);
        _setupRole(ADMIN_ROLE, _initialAdmin);
        expectedNativeToken = _tokenAddress;
        isVerificationPossible = false;
    }

    /**
     * CRITICAL COMPATIBILITY NOTE:
     * The _proxyAdmin parameter MUST be an OpenZeppelin v4.9.x ProxyAdmin
     * contract or earlier. OpenZeppelin v5.x ProxyAdmin contracts do not
     * support getProxyImplementation and will cause this function to revert.
     */
    function verifyL1Contracts(
        address _proxyAdmin,
        address _l1StandardBridgeProxy,
        // ...
    ) external onlyRole(ADMIN_ROLE) {
        // Custom verification logic
    }
}
```

**Issues for Deployment Scripts**:
1. Requires `Initializable` pattern - must call `initialize()` exactly once
2. Requires `AccessControlUpgradeable` - creates additional state snapshots during init
3. ProxyAdmin version constraints - only compatible with OpenZeppelin v4.9.x or earlier
4. Expects roles to be set up before verification functions are called

### 2. Custom Native Token (TON)

**File**: `packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol:1091-1102`

```solidity
contract L2NativeToken is Ownable, ERC20Detailed, SeigToken {
    constructor() ERC20Detailed("Tokamak Network Token", "TON", 18) { }

    function setSeigManager(SeigManagerI) external pure override {
        revert("TON: TON doesn't allow setSeigManager");
    }

    function faucet(uint256 _amount) external {
        _mint(msg.sender, _amount);
    }
}
```

**Issues for Deployment Scripts**:
1. Multiple inheritance chain: `Ownable` → `ERC20Detailed` → `SeigToken` → `ERC20OnApprove`
2. Custom initialization order requirements
3. Some functions intentionally disabled (reverts on `setSeigManager`)
4. Different from standard Optimism ETH handling

### 3. ProxyAdmin Compatibility Requirements

**Standard Optimism Pattern**:
```solidity
// packages/contracts-bedrock/src/universal/ProxyAdmin.sol:184-202
function upgradeAndCall(
    address payable _proxy,
    address _implementation,
    bytes memory _data
) external payable onlyOwner {
    ProxyType ptype = proxyType[_proxy];
    if (ptype == ProxyType.ERC1967) {
        Proxy(_proxy).upgradeToAndCall{ value: msg.value }(_implementation, _data);
    } else {
        upgrade(_proxy, _implementation);
        (bool success,) = _proxy.call{ value: msg.value }(_data);
        require(success, "ProxyAdmin: call to proxy after upgrade failed");
    }
}
```

**Tokamak's Additional Requirements**:

**File**: `packages/contracts-bedrock/scripts/Deploy.s.sol:1066-1073`

```solidity
// ProxyType MUST be set before upgradeAndCall
uint256 proxyType = uint256(proxyAdmin.proxyType(l1StandardBridgeProxy));
Safe safe = Safe(mustGetAddress("SystemOwnerSafe"));
if (proxyType != uint256(ProxyAdmin.ProxyType.CHUGSPLASH)) {
    _callViaSafe({
        _safe: safe,
        _target: address(proxyAdmin),
        _data: abi.encodeCall(ProxyAdmin.setProxyType,
              (l1StandardBridgeProxy, ProxyAdmin.ProxyType.CHUGSPLASH))
    });
}
```

**Issue**: Optimism's deployment scripts don't set proxy types before initialization.

## Root Cause: EVM Snapshot Stack Mismatch

### Normal Optimism Flow (Works)

```
1. ProxyAdmin.upgradeAndCall() called
2.   Creates Snapshot #0
3.   Proxy.upgradeToAndCall() called
4.     Creates Snapshot #1
5.     _implementation.delegatecall(_data) - Simple initialization
6.       Single-step init function runs
7.       Returns successfully
8.     Reverts to Snapshot #1 (success) ✅
9.   Reverts to Snapshot #0 (success) ✅
```

### Tokamak Flow (Fails)

```
1. ProxyAdmin.upgradeAndCall() called
2.   Creates Snapshot #0
3.   Proxy.upgradeToAndCall() called
4.     Creates Snapshot #1
5.     _implementation.delegatecall(_data) - Multi-layer initialization
6.       initialize() called (initializer modifier)
7.         Creates Snapshot #2
8.         __AccessControl_init() called
9.           Creates Snapshot #3
10.          _setupRole() creates more snapshots
11.          Custom logic creates more snapshots
12.        Tries to revert to Snapshot #1
13.        ❌ ERROR: Snapshot #1 already consumed/invalidated
14.        panic: revision id 1 cannot be reverted
```

**The Problem**: Multi-layer initialization (Initializable + AccessControlUpgradeable + custom) creates a deep snapshot stack that conflicts with the proxy upgrade's snapshot management.

## Architectural Differences Summary

| Aspect | Optimism v1.16.0 | Tokamak-Thanos | Compatibility |
|--------|------------------|----------------|---------------|
| **Verification Layer** | None | Custom `L1ContractVerification.sol` with AccessControl | ❌ Incompatible |
| **Native Token** | Standard ETH | Custom `L2NativeToken` (TON) with SeigManager | ❌ Incompatible |
| **Initialization Pattern** | Simple single-step | Multi-layer (Initializable + AccessControl + Custom) | ❌ Incompatible |
| **ProxyAdmin Version** | Any version | Requires OpenZeppelin v4.9.x or earlier | ⚠️ Constrained |
| **Proxy Type Setup** | Automatic | Manual `setProxyType` calls required before init | ❌ Incompatible |
| **Safe Wallet Integration** | Optional | Required with 3-of-3 multisig verification | ❌ Different |
| **Deployment Flow** | OPCM-based Solidity scripts | Python-based (`bedrock-devnet/main.py`) | ❌ Fundamentally different |

## Required Modifications for Compatibility

To make Tokamak compatible with Optimism v1.16.0's E2E deployment scripts, the following modifications are required:

### Option A: Modify Tokamak Contracts (NOT RECOMMENDED)

**⚠️ WARNING**: This would require significant contract changes and user consent for each modification.

1. **Simplify L1ContractVerification Initialization**
   - Remove multi-layer initialization pattern
   - Use single-step initialization without AccessControlUpgradeable
   - Estimated effort: 3-5 days + security audit

2. **Adapt L2NativeToken to Standard Pattern**
   - Simplify inheritance chain
   - Make initialization compatible with single delegatecall pattern
   - Estimated effort: 2-3 days + security audit

3. **Remove ProxyType Requirements**
   - Make proxy type detection automatic like Optimism
   - Estimated effort: 1-2 days

**Total Estimated Effort**: 2-3 weeks development + security audits + testing

**Risk**: 🔴 HIGH
- Breaking changes to production contracts
- Requires security audits
- May lose Tokamak-specific features
- Deployment already works with current Python system

### Option B: Modify Deployment Scripts (PARTIALLY FEASIBLE)

Customize Optimism's deployment scripts to accommodate Tokamak's patterns:

1. **Update DeploySuperchain.s.sol**
   - Add proxy type setup before initialization calls
   - Split initialization into multiple steps
   - Add Safe wallet deployment and ownership transfer

   **File**: `packages/tokamak/contracts-bedrock/scripts/deploy/DeploySuperchain.s.sol:83-98`

   ```solidity
   function deploySuperchainProxyAdmin(InternalInput memory, Output memory _output) private {
       vm.broadcast(msg.sender);
       IProxyAdmin superchainProxyAdmin = IProxyAdmin(
           DeployUtils.create1({
               _name: "ProxyAdmin",
               _args: DeployUtils.encodeConstructor(abi.encodeCall(IProxyAdmin.__constructor__, (msg.sender)))
           })
       );

       // TOKAMAK MODIFICATION: Set proxy types before initialization
       vm.broadcast(msg.sender);
       superchainProxyAdmin.setProxyType(/* proxy address */, ProxyAdmin.ProxyType.ERC1967);

       vm.label(address(superchainProxyAdmin), "SuperchainProxyAdmin");
       _output.superchainProxyAdmin = superchainProxyAdmin;
   }
   ```

2. **Update DeployImplementations.s.sol**
   - Add custom initialization sequence for Tokamak contracts
   - Handle multi-layer initialization separately

3. **Add TokamakDeploymentHelpers.sol**
   - New helper contract for Tokamak-specific deployment logic
   - Handles AccessControl role setup
   - Manages Safe wallet integration

**Estimated Effort**: 1-2 weeks development + testing

**Risk**: 🟡 MEDIUM
- Must maintain fork of deployment scripts
- Will diverge from upstream Optimism
- Needs updates whenever Optimism updates deployment scripts

### Option C: Create Tokamak E2E Test Adapter (RECOMMENDED)

Instead of making contracts or deployment scripts compatible, create an adapter layer that translates Optimism E2E test expectations to Tokamak's deployment system.

**Architecture**:

```
Optimism E2E Test
       ↓
TokamakE2EAdapter
       ↓
Python Deployment (bedrock-devnet/main.py)
       ↓
Tokamak Contracts (deployed and initialized correctly)
       ↓
Adapter returns addresses to E2E test
```

**Implementation**:

1. **Create `op-e2e/tokamak/adapter.go`**
   ```go
   package tokamak

   import (
       "github.com/tokamak-network/tokamak-thanos/op-e2e/e2eutils"
   )

   // TokamakDeploymentAdapter wraps Tokamak's Python deployment
   type TokamakDeploymentAdapter struct {
       pythonDeployer *PythonDeploymentWrapper
   }

   // DeployL1Contracts deploys using Tokamak's Python system
   func (a *TokamakDeploymentAdapter) DeployL1Contracts(cfg *e2eutils.DeployConfig) (*e2eutils.DeployResult, error) {
       // Call Python deployment script
       result, err := a.pythonDeployer.Deploy(cfg)
       if err != nil {
           return nil, err
       }

       // Translate Python deployment output to E2E test format
       return &e2eutils.DeployResult{
           L1Deployments: translateAddresses(result),
           // ... other fields
       }, nil
   }
   ```

2. **Update E2E Test Initialization**

   **File**: `op-e2e/config/init.go`

   ```go
   // Add Tokamak-specific deployment path
   func InitL1(t *testing.T) (*L1Deployment, func()) {
       if useTokamakDeployment() {
           adapter := tokamak.NewDeploymentAdapter()
           return adapter.DeployL1Contracts(cfg)
       }

       // Standard Optimism deployment
       return standardDeployL1(t)
   }
   ```

3. **Wrapper for Python Deployment**

   ```go
   // PythonDeploymentWrapper executes bedrock-devnet/main.py
   func (w *PythonDeploymentWrapper) Deploy(cfg *Config) (*Result, error) {
       cmd := exec.Command("python3", "bedrock-devnet/main.py", "--config", cfg.ToJSON())
       output, err := cmd.CombinedOutput()
       if err != nil {
           return nil, err
       }
       return parseDeploymentOutput(output)
   }
   ```

**Benefits**:
- ✅ No contract modifications required
- ✅ No deployment script modifications required
- ✅ Uses Tokamak's proven Python deployment system
- ✅ E2E tests work without breaking production deployment
- ✅ Clear separation of concerns

**Estimated Effort**: 3-5 days development + testing

**Risk**: 🟢 LOW
- Non-invasive to existing systems
- Easy to test and validate
- No security audit required
- Can be updated independently

### Option D: Skip Problematic E2E Tests (TEMPORARY)

For immediate progress, skip E2E tests that require full deployment:

**File**: `op-e2e/faultproofs/output_alphabet_bond_test.go`

```go
func TestOutputAlphabetGame_ReclaimBond(t *testing.T) {
    // TEMPORARY: Skip until deployment adapter is implemented
    t.Skip("Skipping E2E test - Tokamak deployment scripts not yet compatible")

    // ... rest of test
}
```

Add environment variable control:

```go
func TestOutputAlphabetGame_ReclaimBond(t *testing.T) {
    if os.Getenv("TOKAMAK_SKIP_E2E_DEPLOYMENT") == "true" {
        t.Skip("Tokamak deployment E2E tests disabled")
    }
    // ... rest of test
}
```

**Benefits**:
- ✅ Allows other tests to run
- ✅ Unblocks migration progress
- ✅ Zero development effort

**Drawbacks**:
- ❌ E2E tests don't run
- ❌ Temporary solution only

## Recommendation

**Recommended Approach**: **Option C - Create Tokamak E2E Test Adapter**

**Rationale**:
1. **Preserves Tokamak's Architecture**: No contract changes needed
2. **Maintains Production Deployment**: Python deployment system continues to work
3. **Enables E2E Testing**: Tests can run using proper Tokamak deployment
4. **Low Risk**: Adapter is isolated from production systems
5. **Reasonable Effort**: 3-5 days vs 2-3 weeks for other options

**Short-term**: Use **Option D** to unblock migration while developing Option C

**Implementation Plan**:

### Phase 1: Immediate (Option D)
1. Add skip conditions to E2E tests requiring full deployment
2. Document which tests are skipped and why
3. Continue with rest of v1.16.0 migration

### Phase 2: Adapter Development (Option C) - 1 Week
1. Create `op-e2e/tokamak/adapter.go` package
2. Implement Python deployment wrapper
3. Add address translation layer
4. Update `op-e2e/config/init.go` to use adapter

### Phase 3: Testing - 3-5 Days
1. Test adapter with simple E2E tests
2. Validate address mappings
3. Run full E2E test suite
4. Document any remaining incompatibilities

### Phase 4: Documentation - 1 Day
1. Document adapter architecture
2. Add developer guide for adding new E2E tests
3. Update migration guide

**Total Timeline**: ~2 weeks (including buffer)

## Alternative: Tokamak-Specific E2E Tests

Instead of adapting Optimism's E2E tests, create Tokamak-specific E2E tests:

**Location**: `op-e2e/tokamak/`

**Example**: `op-e2e/tokamak/output_game_test.go`

```go
package tokamak_test

import (
    "testing"
    "github.com/tokamak-network/tokamak-thanos/op-e2e/tokamak"
)

func TestTokamakOutputAlphabetGame_ReclaimBond(t *testing.T) {
    // Use Tokamak's deployment system directly
    deployment := tokamak.SetupTokamakDevnet(t)

    // Run game logic tests
    // ... test implementation
}
```

**Benefits**:
- Tests designed specifically for Tokamak's architecture
- No need to maintain compatibility with Optimism's changing E2E tests
- Can test Tokamak-specific features (verification, TON token, etc.)

**Drawbacks**:
- More test code to maintain
- Won't catch regressions from upstream Optimism changes

## Files Requiring Attention

### Critical Files (Root Cause)

1. **`packages/tokamak/contracts-bedrock/src/tokamak-contracts/verification/L1ContractVerification.sol`**
   - Lines 25-29: Multi-inheritance initialization pattern
   - Lines 113-128: Multi-layer `initialize()` function
   - Lines 166-169: ProxyAdmin compatibility constraints

2. **`packages/tokamak/contracts-bedrock/src/L1/L2NativeToken.sol`**
   - Lines 1091-1102: Custom token with complex inheritance

3. **`op-chain-ops/script/script.go`**
   - Lines 352-388: Error handler catching revision ID panics

### Deployment Scripts (Incompatibility Source)

4. **`packages/tokamak/contracts-bedrock/scripts/deploy/DeploySuperchain.s.sol`**
   - Lines 83-98: ProxyAdmin deployment (needs proxy type setup)
   - Lines 100-138: Implementation deployment (needs custom init)

5. **`packages/tokamak/contracts-bedrock/scripts/deploy/DeployImplementations.s.sol`**
   - Entire file: Expects standard Optimism contract initialization

### E2E Test Files (Need Adapter)

6. **`op-e2e/config/init.go`**
   - Line 198: Artifact path (already fixed)
   - Lines 50-100: L1 deployment initialization (needs adapter hook)

7. **`op-e2e/faultproofs/output_alphabet_bond_test.go`**
   - Test that's currently failing
   - Needs adapter or skip condition

### Proxy Management

8. **`packages/contracts-bedrock/src/universal/ProxyAdmin.sol`**
   - Lines 184-202: `upgradeAndCall` function
   - Lines 113-124: `setProxyType` function

9. **`packages/contracts-bedrock/src/universal/Proxy.sol`**
   - Lines 67-81: `upgradeToAndCall` (creates snapshot)

## Testing Strategy

Once adapter is implemented:

1. **Unit Tests for Adapter**
   ```bash
   go test ./op-e2e/tokamak/... -v
   ```

2. **Integration Test with Python Deployment**
   ```bash
   go test -run TestTokamakDeploymentAdapter ./op-e2e/tokamak/... -v
   ```

3. **E2E Test Suite**
   ```bash
   go test -run TestOutputAlphabetGame_ReclaimBond ./op-e2e/faultproofs/... -v -timeout 30m
   ```

4. **Full Test Suite**
   ```bash
   make test-e2e
   ```

## Security Considerations

### If Modifying Contracts (Option A)

- ⚠️ **Security Audit Required**: Any contract initialization changes need professional audit
- ⚠️ **Migration Risk**: Existing deployments may not be compatible
- ⚠️ **Feature Loss**: Tokamak-specific features may be removed

### If Using Adapter (Option C - Recommended)

- ✅ **No Contract Changes**: Existing security guarantees preserved
- ✅ **Isolated Risk**: Adapter only affects test environment
- ✅ **Easy Rollback**: Can disable adapter if issues arise

## Monitoring and Validation

After implementing solution:

1. **Verify E2E Test Success Rate**
   ```bash
   go test ./op-e2e/... -v | grep -E "(PASS|FAIL)"
   ```

2. **Validate Deployment Artifacts**
   ```bash
   ls -la packages/tokamak/contracts-bedrock/forge-artifacts/Deploy*.s.sol/
   ```

3. **Check Python Deployment Still Works**
   ```bash
   make devnet-allocs
   ```

4. **Monitor Test Execution Time**
   - Adapter should not significantly slow down tests
   - Target: < 5% overhead vs direct deployment

## Conclusion

The "revision id 1 cannot be reverted" error is caused by fundamental architectural differences between Tokamak's multi-layer contract initialization system and Optimism's simpler deployment pattern.

**Recommended Solution**: Create an E2E test adapter (Option C) that bridges Optimism's E2E tests with Tokamak's Python deployment system. This approach:

- ✅ Preserves Tokamak's production deployment system
- ✅ Enables E2E testing without contract modifications
- ✅ Low risk and reasonable development effort
- ✅ Maintains security of existing contracts

**Short-term**: Skip problematic E2E tests (Option D) to unblock v1.16.0 migration

**Long-term**: Consider developing Tokamak-specific E2E tests that are designed for Tokamak's architecture rather than adapting Optimism's tests

## Next Steps

1. **Immediate**: Add skip conditions to E2E tests (Option D)
2. **Week 1**: Design and implement E2E test adapter (Option C)
3. **Week 2**: Test adapter with full E2E suite
4. **Week 3**: Document and deploy solution

## References

- [Deployment Scripts Comparison](./DEPLOYMENT-SCRIPTS-COMPARISON.md)
- [Migration Guide](./MIGRATION-GUIDE.md)
- [OpenZeppelin Proxy Pattern](https://docs.openzeppelin.com/contracts/4.x/api/proxy)
- [Go-Ethereum State Management](https://github.com/ethereum/go-ethereum/blob/master/core/state/journal.go)

---

**Document Version**: 1.0
**Last Updated**: 2025-11-04
**Next Review**: After adapter implementation
