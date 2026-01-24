# Shutdown Unit Tests

This document summarizes the unit tests implemented under `test/shutdown` and the latest local results.

## Test Coverage Overview

### ForceWithdrawBridge
File: `test/shutdown/ForceWithdrawBridge.t.sol`

| Test | Content/Purpose | Success Criteria |
| --- | --- | --- |
| `test_step3_registerPositions` | Verify position registration flow | Position registration status is true |
| `test_step4_activateForceWithdrawal` | Verify force withdrawal activation flow | `active()` returns true |
| `test_step5_forceWithdrawClaim_ERC20_succeeds` | Verify ERC20 force withdrawal success path | Receiver balance increases and `claimState` is true |
| `test_step6_forceWithdrawClaim_ETH_succeeds` | Verify ETH force withdrawal success path | Receiver balance increases and `claimState` is true |
| `test_error_invalidHash_reverts` | Verify rejection of invalid hash input | Reverts with `FW_INVALID_HASH` |
| `test_error_doubleClaimreverts` | Verify rejection of duplicate claims for a hash | Reverts with `already claim Hash` |
| `test_error_unregisteredPosition_reverts` | Verify rejection of unregistered positions | Reverts with `FW_NOT_AVAILABLE_POSITION` |
| `test_step7_batchClaim_succeeds` | Verify batch claim success path | All receiver balances increase and `claimState` is true |

Result:
- `forge test --match-path test/shutdown/ForceWithdrawBridge.t.sol`
- 8 passed, 0 failed, 0 skipped

### OptimismPortalClosing
File: `test/shutdown/OptimismPortalClosing.t.sol`

| Test | Content/Purpose | Success Criteria |
| --- | --- | --- |
| `test_constructor_revertsOnZeroImplementation` | Block deployment when implementation address is `address(0)` | Fails with specified revert message |
| `test_version_isClosing` | Check closing version string | `version()` returns `2.8.1-closing` |
| `test_depositTransaction_reverts` | Verify blocking of deposit path | Fails with specified revert message |
| `test_onApprove_reverts` | Verify blocking of approve-and-call deposit path | Fails with specified revert message |
| `test_receive_reverts` | Verify blocking of direct ETH transfers | Fails with specified revert message |
| `test_fallback_delegatesToImplementation` | Check normal delegation of fallback delegatecall | Reflects implementation contract state/return values |

Result:
- `forge test --match-path test/shutdown/OptimismPortalClosing.t.sol`
- 6 passed, 0 failed, 0 skipped

### ShutdownL1UsdcBridge
File: `test/shutdown/ShutdownL1UsdcBridge.t.sol`

| Test | Content/Purpose | Success Criteria |
| --- | --- | --- |
| `test_sweepUSDC_revertsWhenAdminNotSet` | Block calls if proxy admin is not set | Fails with specified revert message |
| `test_sweepUSDC_revertsWhenL1UsdcUnset` | Block calls if `l1Usdc` is not set | Fails with specified revert message |
| `test_sweepUSDC_succeedsForAdminEOA` | Allow sweep with EOA admin permissions | Receiver balance increases, bridge balance decreases |
| `test_sweepUSDC_succeedsForAdminOwnerContract` | Succeed when admin is a contract and called by owner | Receiver balance increases, bridge balance decreases |
| `test_sweepUSDC_revertsForNonOwner` | Block callers who are not the owner | Fails with specified revert message |
| `test_sweepUSDC_revertsWhenAdminOwnerLookupFails` | Block if `admin.owner()` lookup fails | Fails with specified revert message |

Result:
- `forge test --match-path test/shutdown/ShutdownL1UsdcBridge.t.sol`
- 6 passed, 0 failed, 0 skipped

### ShutdownOptimismPortal
File: `test/shutdown/ShutdownOptimismPortal.t.sol`

| Test | Content/Purpose | Success Criteria |
| --- | --- | --- |
| `test_sweepNativeToken_revertsWhenAdminNotSet` | Block calls if proxy admin is not set | Fails with specified revert message |
| `test_sweepNativeToken_revertsWhenNativeTokenIsETH` | Block sweep if native token is ETH | Fails with specified revert message |
| `test_sweepNativeToken_succeedsForAdminEOA` | Allow sweep with EOA admin permissions | Receiver balance increases, portal balance decreases |
| `test_sweepNativeToken_succeedsForAdminOwnerContract` | Succeed when admin is a contract and called by owner | Receiver balance increases, portal balance decreases |
| `test_sweepNativeToken_revertsForNonOwner` | Block callers who are not the owner | Fails with specified revert message |
| `test_sweepNativeToken_revertsWhenAdminOwnerLookupFails` | Block if `admin.owner()` lookup fails | Fails with specified revert message |

Result:
- `forge test --match-path test/shutdown/ShutdownOptimismPortal.t.sol`
- 6 passed, 0 failed, 0 skipped

### ShutdownOptimismPortal2
File: `test/shutdown/ShutdownOptimismPortal2.t.sol`

| Test | Content/Purpose | Success Criteria |
| --- | --- | --- |
| `test_sweepNativeToken_revertsWhenAdminNotSet` | Block calls if proxy admin is not set | Fails with specified revert message |
| `test_sweepNativeToken_revertsWhenNativeTokenIsETH` | Block sweep if native token is ETH | Fails with specified revert message |
| `test_sweepNativeToken_succeedsForAdminEOA` | Allow sweep with EOA admin permissions | Receiver balance increases, portal balance decreases |
| `test_sweepNativeToken_succeedsForAdminOwnerContract` | Succeed when admin is a contract and called by owner | Receiver balance increases, portal balance decreases |
| `test_sweepNativeToken_revertsForNonOwner` | Block callers who are not the owner | Fails with specified revert message |
| `test_sweepNativeToken_revertsWhenAdminOwnerLookupFails` | Block if `admin.owner()` lookup fails | Fails with specified revert message |

Result:
- `forge test --match-path test/shutdown/ShutdownOptimismPortal2.t.sol`
- 6 passed, 0 failed, 0 skipped
