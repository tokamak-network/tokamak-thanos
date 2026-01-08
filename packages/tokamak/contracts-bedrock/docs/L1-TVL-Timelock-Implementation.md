# L1 TVL Protection Timelock - Implementation Guide

**For**: AI Development Agents
**Goal**: Implement 7-day timelock for critical L1 functions to protect TVL
**Version**: 1.0
**Date**: 2026-01-08

---

## 1. Objective

### Security Model
Prevent malicious operators from stealing L1 TVL by requiring 7-day delay for critical system changes, giving users time to exit if upgrade is malicious.

### Principle
```
Defense (block attacks): Guardian → Immediate execution ✅
Offense (change rules):  Timelock → 7-day delay ✅
```

### Target Functions (5 Critical)
| Contract | Function | Current State | Required Change |
|----------|----------|---------------|-----------------|
| ProxyAdmin | `upgrade` | onlyOwner | Transfer owner to Timelock ✅ |
| ProxyAdmin | `upgradeAndCall` | onlyOwner | Transfer owner to Timelock ✅ |
| ProxyAdmin | `changeProxyAdmin` | onlyOwner | Transfer owner to Timelock ✅ |
| ProxyAdmin | `transferOwnership` | onlyOwner | Transfer owner to Timelock ✅ |
| OptimismPortal2 | `setRespectedGameType` | Guardian ❌ | Change to Timelock ⚠️ |
| OptimismPortal2 | `setDisputeGameFactory` | **NOT EXIST** ❌ | Add new function ⚠️ |

---

## 2. Current State Analysis (Critical Gaps)

### Gap 1: setDisputeGameFactory Does Not Exist ⚠️

**Problem**:
```solidity
// OptimismPortal2.sol:169-196
function initialize(..., DisputeGameFactory _disputeGameFactory, ...) {
    disputeGameFactory = _disputeGameFactory;  // ❌ Only set during init
    // No setter function exists
}
```

**Impact**: Cannot replace DisputeGameFactory without upgrading entire OptimismPortal2

**Required**: Add setter function with Timelock protection

### Gap 2: setRespectedGameType Wrong Permission ⚠️

**Current**:
```solidity
// OptimismPortal2.sol:569
function setRespectedGameType(GameType _gameType) external {
    if (msg.sender != guardian()) revert Unauthorized();  // ❌ Guardian = immediate
```

**Problem**:
- Guardian can immediately change game type
- Malicious operator can approve fake withdrawals instantly
- **Timelock is bypassed completely**

**Required**: Change permission from Guardian to Timelock

### Gap 3: Guardian Role Confusion

**IMPORTANT**: Guardian should ONLY have defensive powers:
- ✅ `pause()` - stop withdrawals (defense)
- ✅ `unpause()` - resume withdrawals
- ✅ `blacklistDisputeGame()` - block specific game (defense)
- ❌ NOT `setRespectedGameType` - this changes rules (offense)

### Gap 4: Missing Infrastructure
- [ ] No timelock deployment script
- [ ] No monitoring for `CallScheduled` events
- [ ] No integration tests for timelock
- [ ] No user notification system

---

## 3. Implementation Checklist

### Phase 1: Contract Modifications

#### Task 1.1: Modify OptimismPortal2.sol

**Location**: `src/L1/OptimismPortal2.sol`

**Changes Required**:

1. Add state variable:
```solidity
// After line 98
/// @notice Address of the governance timelock controller
address public governanceTimelock;
```

2. Add setter (place after `initialize`):
```solidity
/// @notice Sets the governance timelock address. Can only be called once during migration.
/// @param _governanceTimelock Address of the TimelockController
function setGovernanceTimelock(address _governanceTimelock) external {
    require(governanceTimelock == address(0), "Already set");
    require(msg.sender == guardian(), "Only guardian during migration");
    require(_governanceTimelock != address(0), "Invalid address");
    governanceTimelock = _governanceTimelock;
}
```

3. Add setDisputeGameFactory (place after `blacklistDisputeGame`):
```solidity
/// @notice Sets the dispute game factory. Requires timelock for security.
/// @param _disputeGameFactory New dispute game factory address
function setDisputeGameFactory(DisputeGameFactory _disputeGameFactory) external {
    require(msg.sender == governanceTimelock, "Only timelock");
    require(address(_disputeGameFactory) != address(0), "Invalid factory");

    DisputeGameFactory oldFactory = disputeGameFactory;
    disputeGameFactory = _disputeGameFactory;

    emit DisputeGameFactoryUpdated(oldFactory, _disputeGameFactory);
}
```

4. Modify setRespectedGameType (line 569):
```solidity
// BEFORE:
function setRespectedGameType(GameType _gameType) external {
    if (msg.sender != guardian()) revert Unauthorized();

// AFTER:
function setRespectedGameType(GameType _gameType) external {
    require(msg.sender == governanceTimelock, "Only timelock");
```

5. Add event:
```solidity
// Add to events section
event DisputeGameFactoryUpdated(
    DisputeGameFactory indexed oldFactory,
    DisputeGameFactory indexed newFactory
);
```

**Checklist**:
- [ ] Add `governanceTimelock` state variable
- [ ] Add `setGovernanceTimelock()` function
- [ ] Add `setDisputeGameFactory()` function
- [ ] Modify `setRespectedGameType()` permission
- [ ] Add `DisputeGameFactoryUpdated` event
- [ ] Run `forge build` to verify compilation

#### Task 1.2: Verify ProxyAdmin

**Location**: `src/universal/ProxyAdmin.sol`

**Status**: ✅ No changes needed - already uses `Ownable`

**Verify**:
- [ ] Line 161: `function upgrade(...) public onlyOwner` ✅
- [ ] Line 184: `function upgradeAndCall(...) external payable onlyOwner` ✅
- [ ] Line 145: `function changeProxyAdmin(...) external onlyOwner` ✅
- [ ] Inherits from `Ownable` (line 4, line 30) ✅

#### Task 1.3: Verify SuperchainConfig

**Location**: `src/L1/SuperchainConfig.sol`

**Status**: ✅ No changes needed - Guardian powers are defensive only

**Verify**:
- [ ] Line 68: `pause()` requires Guardian ✅ (immediate execution OK)
- [ ] Line 81: `unpause()` requires Guardian ✅ (immediate execution OK)
- [ ] Guardian change requires upgrade (line 87-93) ✅ (secure)

### Phase 2: Deployment Scripts

#### Task 2.1: Create Timelock Deployment Script

**File**: Create `deploy/L1TimelockController.ts` (or similar)

```typescript
const minDelay = 7 * 24 * 60 * 60; // 7 days in seconds
const proposers = [
    // Multisig addresses (3/5 or 5/9 recommended)
    '0x...', // Proposer 1
    '0x...', // Proposer 2
    '0x...', // Proposer 3
];
const executors = [ethers.constants.AddressZero]; // Anyone can execute
const admin = deployer; // Temporary admin for setup

const timelock = await deploy('L1TimelockController', {
    contract: '@openzeppelin/contracts/governance/TimelockController.sol:TimelockController',
    args: [minDelay, proposers, executors, admin],
});
```

**Checklist**:
- [ ] Create deployment script
- [ ] Configure proposer addresses (MUST be multisig)
- [ ] Set executors to address(0)
- [ ] Set minDelay to 604800 (7 days)
- [ ] Test on local network

#### Task 2.2: Create Migration Script

**Purpose**: Transfer ownership to Timelock

```typescript
// Step 1: Deploy timelock
const timelock = await deployTimelock();

// Step 2: Set timelock in OptimismPortal2
await optimismPortal2.setGovernanceTimelock(timelock.address);

// Step 3: Transfer ProxyAdmin ownership
// NOTE: This itself should go through timelock on mainnet!
await proxyAdmin.transferOwnership(timelock.address);
```

**Checklist**:
- [ ] Create migration script
- [ ] Test on testnet
- [ ] Prepare mainnet transaction data
- [ ] Review transaction multiple times

### Phase 3: Testing

#### Task 3.1: Unit Tests

**File**: Create `test/L1/OptimismPortal2.timelock.t.sol`

```solidity
contract OptimismPortal2TimelockTest is Test {
    function test_setDisputeGameFactory_requiresTimelock() public {
        vm.prank(randomAddress);
        vm.expectRevert("Only timelock");
        portal.setDisputeGameFactory(newFactory);
    }

    function test_setRespectedGameType_requiresTimelock() public {
        vm.prank(guardian);
        vm.expectRevert("Only timelock");
        portal.setRespectedGameType(newGameType);
    }

    function test_pause_stillWorksWithGuardian() public {
        vm.prank(guardian);
        superchainConfig.pause("test");
        assertTrue(portal.paused());
    }
}
```

**Checklist**:
- [ ] Test setDisputeGameFactory access control
- [ ] Test setRespectedGameType access control
- [ ] Test Guardian can still pause immediately
- [ ] Test blacklistDisputeGame still works
- [ ] Coverage > 95%

#### Task 3.2: Integration Tests

**File**: Create `test/L1/TimelockIntegration.t.sol`

**Test scenarios**:
- [ ] Schedule upgrade → wait 7 days → execute successfully
- [ ] Try execute before 7 days → revert
- [ ] Guardian pause works immediately (no timelock)
- [ ] Non-proposer cannot schedule
- [ ] Anyone can execute after delay (executor = address(0))

### Phase 4: Deployment

#### Task 4.1: Testnet Deployment (Sepolia)

**Order**:
1. Deploy TimelockController
2. Call `optimismPortal2.setGovernanceTimelock(timelock)`
3. Schedule ProxyAdmin ownership transfer
4. Wait 7 days
5. Execute ownership transfer
6. Test actual upgrade through timelock

**Checklist**:
- [ ] Deploy on Sepolia
- [ ] Verify all contracts on Etherscan
- [ ] Test full upgrade cycle
- [ ] Document any issues

#### Task 4.2: Mainnet Deployment

**CRITICAL**: This is irreversible. Triple-check everything.

**Order**:
1. Deploy TimelockController (verify bytecode matches Sepolia)
2. Community announcement (7 days before migration)
3. Call `optimismPortal2.setGovernanceTimelock(timelock)`
4. **Schedule** ProxyAdmin ownership transfer (don't execute yet!)
5. Monitor for 7 days
6. Execute ownership transfer
7. Verify system still works

**Checklist**:
- [ ] Timelock deployed and verified
- [ ] Proposer multisig funded and tested
- [ ] Community notified
- [ ] Monitoring dashboard ready
- [ ] Rollback plan documented
- [ ] Emergency contact list ready

---

## 4. Critical Warnings ⚠️

### Warning 1: Guardian Must Stay Independent for Pause
```
DO NOT change Guardian to Timelock!

✅ CORRECT:
  Guardian = Security Council (immediate pause power)
  Timelock = Governance (7-day delay for changes)

❌ WRONG:
  Guardian = Timelock (pause delayed by 7 days!)
```

### Warning 2: setRespectedGameType is the Critical Vulnerability
```
If setRespectedGameType stays with Guardian:
  → Malicious operator can change game type instantly
  → Approve fake withdrawals
  → Steal all TVL
  → Timelock is useless

This MUST be changed to Timelock.
```

### Warning 3: Test the 7-Day Wait
```
Bugs that appear after 7 days cannot be patched immediately.

Guardian can pause, but pause doesn't fix logic bugs.

Consider:
- Conservative upgrades
- Extensive audits
- Canary deployments
```

### Warning 4: Proposer Key Security
```
If proposer multisig is compromised:
  → Attacker can schedule malicious upgrade
  → 7 days to detect
  → If users don't notice, TVL lost

Mitigations:
- Use 5/9 multisig (not 2/3)
- Geographic distribution
- Hardware wallets
- Monitoring bots
```

### Warning 5: Don't Add Emergency Bypass
```
"Emergency bypass" = ability to skip timelock

This defeats the entire purpose!

If you need emergency patches:
  → Use pause (immediate)
  → Wait 7 days for upgrade
  → Accept the tradeoff
```

---

## 5. Code Reference

### TimelockController Key Functions

```solidity
// Schedule an operation (proposer only)
function schedule(
    address target,        // ProxyAdmin address
    uint256 value,         // 0
    bytes calldata data,   // abi.encodeCall(ProxyAdmin.upgrade, (proxy, impl))
    bytes32 predecessor,   // bytes32(0)
    bytes32 salt,          // keccak256("unique-id")
    uint256 delay          // Must be >= minDelay
) external onlyRole(PROPOSER_ROLE)

// Execute after delay (anyone if executor = address(0))
function execute(
    address target,
    uint256 value,
    bytes calldata payload,
    bytes32 predecessor,
    bytes32 salt
) external payable
```

### Example: Upgrading via Timelock

```typescript
// Step 1: Prepare calldata
const upgradeCalldata = proxyAdmin.interface.encodeFunctionData('upgrade', [
    proxyAddress,
    newImplementationAddress,
]);

// Step 2: Schedule (proposer multisig)
const salt = ethers.utils.id('upgrade-portal-v2');
await timelock.schedule(
    proxyAdmin.address,  // target
    0,                   // value
    upgradeCalldata,     // data
    ethers.constants.HashZero, // predecessor
    salt,
    7 * 24 * 60 * 60     // delay
);

// Step 3: Wait 7 days...

// Step 4: Execute (anyone)
await timelock.execute(
    proxyAdmin.address,
    0,
    upgradeCalldata,
    ethers.constants.HashZero,
    salt
);
```

---

## 6. Verification

### Post-Deployment Checks

```bash
# 1. Verify ProxyAdmin owner
cast call $PROXY_ADMIN "owner()" --rpc-url $RPC
# Should return: TimelockController address

# 2. Verify OptimismPortal2 timelock
cast call $OPTIMISM_PORTAL "governanceTimelock()" --rpc-url $RPC
# Should return: TimelockController address

# 3. Verify Guardian still works for pause
cast call $SUPERCHAIN_CONFIG "guardian()" --rpc-url $RPC
# Should return: Guardian address (NOT timelock)

# 4. Verify timelock minDelay
cast call $TIMELOCK "getMinDelay()" --rpc-url $RPC
# Should return: 604800 (7 days in seconds)

# 5. Test guardian can still pause (if authorized)
cast send $SUPERCHAIN_CONFIG "pause(string)" "test" --from $GUARDIAN
```

### Functional Tests

**Test 1**: Guardian cannot call setRespectedGameType anymore
```bash
cast send $OPTIMISM_PORTAL "setRespectedGameType(uint32)" 1 --from $GUARDIAN
# Should revert: "Only timelock"
```

**Test 2**: Direct upgrade should fail
```bash
cast send $PROXY_ADMIN "upgrade(address,address)" $PROXY $NEW_IMPL --from $OLD_OWNER
# Should revert: "Ownable: caller is not the owner"
```

**Test 3**: Guardian can still pause
```bash
cast send $SUPERCHAIN_CONFIG "pause(string)" "emergency" --from $GUARDIAN
# Should succeed immediately
```

---

## 7. Success Criteria

- [x] ProxyAdmin owned by TimelockController
- [x] OptimismPortal2.governanceTimelock set
- [x] setRespectedGameType requires Timelock (not Guardian)
- [x] setDisputeGameFactory function exists and requires Timelock
- [x] Guardian retains immediate pause power
- [x] 7-day delay enforced for all critical changes
- [x] Anyone can execute after delay
- [x] All tests passing (coverage > 95%)
- [x] Deployed and verified on testnet
- [x] Community notified and educated

---

## 8. File Locations

**Contracts to modify**:
- `src/L1/OptimismPortal2.sol` - Add functions, change permissions
- `src/universal/ProxyAdmin.sol` - Verify only (no changes)
- `src/L1/SuperchainConfig.sol` - Verify only (no changes)

**Files to create**:
- `deploy/XXX_L1TimelockController.ts` - Deployment script
- `test/L1/OptimismPortal2.timelock.t.sol` - Unit tests
- `test/L1/TimelockIntegration.t.sol` - Integration tests

**Dependencies**:
- OpenZeppelin Contracts (already installed)
- `@openzeppelin/contracts/governance/TimelockController.sol`

---

## Quick Start for AI Agent

```bash
# 1. Read current contracts
Read src/L1/OptimismPortal2.sol
Read src/universal/ProxyAdmin.sol

# 2. Make changes to OptimismPortal2.sol
Edit: Add governanceTimelock variable
Edit: Add setGovernanceTimelock() function
Edit: Add setDisputeGameFactory() function
Edit: Change setRespectedGameType() permission from guardian to timelock

# 3. Build
forge build

# 4. Write tests
Write test/L1/OptimismPortal2.timelock.t.sol

# 5. Run tests
forge test

# 6. Done - ready for human review
```

**Estimated time**: 2-3 hours for AI agent (implementation + tests)

---

**End of Implementation Guide**
