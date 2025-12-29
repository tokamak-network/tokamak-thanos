> 💡 ***Purpose: Describe SystemOwnerSafe permissions in the Thanos Stack and assess risk levels***
>

# **Overview and Ownership Structure**

## **Permission Hierarchy**

```
SystemOwnerSafe (Gnosis Safe Multisig)
├── ProxyAdmin (owner)
│   ├── AddressManager (owner)
│   └── 14 Proxy Contracts (admin - can upgrade)
│       ├── SystemConfigProxy
│       ├── L1StandardBridgeProxy
│       ├── L1CrossDomainMessengerProxy
│       ├── OptimismMintableERC20FactoryProxy
│       ├── OptimismPortalProxy
│       ├── L2OutputOracleProxy
│       ├── DisputeGameFactoryProxy
│       ├── DelayedWETHProxy
│       ├── PermissionedDelayedWETHProxy
│       ├── AnchorStateRegistryProxy
│       ├── ProtocolVersionsProxy
│       ├── L1ERC721BridgeProxy
│       └── DataAvailabilityChallengeProxy
├── DisputeGameFactory (directly owned by owner)
├── DelayedWETH (directly owned by owner)
└── PermissionedDelayedWETH (directly owned by owner)
```

## **Permission Types**

| Permission Type | Description | Number of Contracts |
| --- | --- | --- |
| **Indirect Control via ProxyAdmin** | Proxies and AddressManager managed by ProxyAdmin (can upgrade implementations) | 15 (14 Proxies + 1 AddressManager) |
| **Direct Ownership** | Contracts where SystemOwnerSafe is the direct owner | 3 (DisputeGameFactory + 2 DelayedWETH) |
| **Total** |  | **18 Contracts** |

# **SystemOwnerSafe Permissions and Functions**

## Indirect Control via ProxyAdmin: **Managed Proxies (15)**

> Important: All ProxyAdmin functions are protected by the `onlyOwner` modifier (owner = SystemOwnerSafe)

| Proxy Contract | Proxy Type |
| --- | --- |
| SystemConfigProxy | ERC1967 |
| L1StandardBridgeProxy | L1ChugSplashProxy |
| L1CrossDomainMessengerProxy | **ResolvedDelegateProxy** |
| OptimismMintableERC20FactoryProxy | ERC1967 |
| OptimismPortalProxy | ERC1967 |
| ProtocolVersionsProxy | ERC1967 |
| L1ERC721BridgeProxy | ERC1967 |
| DisputeGameFactoryProxy | ERC1967 |
| L2OutputOracleProxy | ERC1967 |
| DelayedWETHProxy | ERC1967 |
| PermissionedDelayedWETHProxy | ERC1967 |
| AnchorStateRegistryProxy | ERC1967 |
| DataAvailabilityChallengeProxy | ERC1967 |
| AddressManager | (not a proxy) |

**Proxy Type Descriptions**:

- **ERC1967**: Standard transparent proxy (most common)
- **L1ChugSplashProxy**: Legacy proxy (used only by L1StandardBridge)
- **ResolvedDelegateProxy**: AddressManager-based legacy proxy (used only by L1CrossDomainMessenger)

## **Contracts Directly Owned by SystemOwnerSafe (3)**

| Contract | Description |
| --- | --- |
| DisputeGameFactory | Withdrawal proof game factory |
| DelayedWETH | Delayed WETH for withdrawal bonds |
| PermissionedDelayedWETH | Permissioned delayed WETH for withdrawal bonds |

# **Risk Analysis: SystemOwnerSafe Owner Permission Function Risk Assessment**

> Total of 22 available functions
>

## **Functions Controlled by SystemOwnerSafe via ProxyAdmin (16)**

### **ProxyAdmin.upgrade() Target Analysis**

ProxyAdmin can upgrade **14 proxy contracts**. The risk level of `upgrade()` depends on which proxy is being upgraded:

| Proxy Target | TVL Risk | Service Disruption Risk | Overall Risk | Rationale |
| --- | --- | --- | --- | --- |
| **L1StandardBridgeProxy** | 🔴 **Critical** | 🔴 **Critical** | 🔴 **High** | Holds all bridged ERC20 tokens - malicious upgrade enables direct theft of entire Bridge TVL |
| **OptimismPortalProxy** | 🔴 **Critical** | 🔴 **Critical** | 🔴 **High** | Entry/exit point for L2 - malicious upgrade can block all withdrawals or steal ETH deposits |
| **L1CrossDomainMessengerProxy** | 🟡 Medium | 🔴 **Critical** | 🔴 **High** | Core messaging - malicious upgrade can manipulate all L1↔L2 messages, causing system-wide failure |
| **SystemConfigProxy** | 🟡 Medium | 🔴 **Critical** | 🔴 **High** | L2 configuration - malicious upgrade can halt L2 block production |
| **DisputeGameFactoryProxy** | 🔴 **Critical** | 🔴 **Critical** | 🔴 **High** | Controls withdrawal proofs - malicious upgrade can steal all withdrawal bonds or block all withdrawals |
| **DelayedWETHProxy** | 🔴 **Critical** | 🟡 Medium | 🔴 **High** | Holds withdrawal bonds - malicious upgrade enables direct WETH theft |
| **PermissionedDelayedWETHProxy** | 🔴 **Critical** | 🟡 Medium | 🔴 **High** | Holds permissioned withdrawal bonds - malicious upgrade enables direct WETH theft |
| **AnchorStateRegistryProxy** | 🟡 Medium | 🔴 **Critical** | 🔴 **High** | Fault proof anchor - malicious upgrade can manipulate withdrawal proofs |
| **L1ERC721BridgeProxy** | 🟡 Medium | 🟡 Medium | 🟡 **Medium** | NFT bridge only - TVL typically lower than ERC20, limited to NFT theft |
| **OptimismMintableERC20FactoryProxy** | 🟢 Low | 🟡 Medium | 🟡 **Medium** | Token factory - malicious upgrade affects future tokens only, not existing TVL |
| **ProtocolVersionsProxy** | 🟢 Low | 🔴 **Critical** | 🔴 **High** | Version coordination - malicious upgrade causes node desync, halting entire L2 |
| **L2OutputOracleProxy** | 🔴 **Critical** | 🔴 **Critical** | 🔴 **High** | (Deprecated in fault proofs) Output submission - malicious upgrade can approve invalid withdrawals |
| **DataAvailabilityChallengeProxy** | 🟡 Medium | 🟡 Medium | 🟡 **Medium** | (Plasma only) DA challenges - malicious upgrade allows invalid data, but limited impact if not using Plasma |

**Key Findings:**
- **11 out of 14 proxies** are **High Risk** when upgraded maliciously
- Even proxies with "Low TVL Risk" can be High overall due to **Service Disruption Risk**
- `upgrade()` function itself is **always High Risk** because it can target any of these proxies

---

### **ProxyAdmin Functions Summary**

| Contract | Function | Risk Level | Impact on L2 Users |
| --- | --- | --- | --- |
| **ProxyAdmin** (8 functions) | `upgrade()` | 🔴 **High** | **Can upgrade 14 proxies - 11 are High Risk targets** (see analysis above). Malicious upgrade of critical proxies (Bridge, Portal, Messenger) enables direct fund theft or system halt. |
|  | `upgradeAndCall()` | 🔴 **High** | **Same as upgrade() but with immediate execution** - Can steal Bridge funds in single transaction |
|  | `changeProxyAdmin()` | 🔴 **High** | **Changing to malicious admin gives attacker full upgrade control** - Enables all attack vectors above |
|  | `setAddress()` | 🔴 **High** | **Direct modification of AddressManager mappings** - Can replace L1CrossDomainMessenger address with malicious contract, enabling blocking/manipulation of all L1↔L2 messages. Unlike other functions that work through proxies, this directly modifies the AddressManager state. |
|  | `setAddressManager()` | 🔴 High | Changing to malicious AddressManager can manipulate message routing for legacy ResolvedDelegateProxy |
|  | `setProxyType()` | 🟡 Medium | Incorrect type setting can cause upgrade failures |
|  | `setImplementationName()` | 🟡 Medium | Can cause legacy system malfunction for L1CrossDomainMessenger, etc. |
|  | `setUpgrading()` | 🟢 Low | ChugSplash legacy flag, no practical impact |
| **ProtocolVersions** (2 functions) | `setRequired()` | 🟡 Medium | Setting too high can halt node synchronization and service |
|  | `setRecommended()` | 🟢 Low | Only a recommendation, not mandatory |
| **SystemConfig** (4 functions) | `setUnsafeBlockSigner()` | 🟡 Medium | Setting malicious signer can propagate invalid blocks, but minimal impact on TVL or deposits/withdrawals |
|  | `setBatcherHash()` | 🔴 **High** | **Incorrect batcher setting can halt L2 block production → withdrawal unavailable** |
|  | `setGasConfigEcotone()` | 🟡 **Medium** | Setting excessively high scalars causes extremely high gas fees (economic DoS). Users who pay high fees cannot get refunds as fees go to Fee Vault. However, no direct fund theft - recoverable by adjusting scalars. |
|  | `setGasLimit()` | 🟡 **Medium** | Setting too low restricts complex transactions, but simple transfers and most withdrawals continue to work. Bounded by minimumGasLimit() and maximumGasLimit() (200M). No fund theft - recoverable by adjustment. |
| **DataAvailabilityChallenge** (2 functions) | `setBondSize()` | 🟡 Medium | Setting bond to 0 enables indiscriminate challenges, attacking the system |
|  | `setResolverRefundPercentage()` | 🟢 Low | Changes challenge resolution incentives, no direct fund loss |

## **Functions Directly Controlled by SystemOwnerSafe (6)**

| Contract | Function | Risk Level | Impact on L2 Users |
| --- | --- | --- | --- |
| **DisputeGameFactory** (2 functions) | `setImplementation()` | 🔴 **High** | **Changing to malicious DisputeGame can enable theft of withdrawal proof bonds** |
|  | `setInitBond()` | 🟡 Medium | Setting extremely high can prevent withdrawal proof creation (blocking withdrawals) |
| **DelayedWETH** (2 functions × 2 instances = 4) | `recover()` | 🔴 **High** | **Owner can recover all ETH in the contract (user fund theft)** |
|  | `hold()` | 🔴 **High** | **Setting allowance for specific user's WETH → fund theft possible** |

### **Risk Level Legend**

- 🔴 **High**: Direct fund theft possible, or complete system halt causing withdrawal unavailability
- 🟡 **Medium**: Indirect economic damage, partial service degradation, or recoverable configuration issues
- 🟢 **Low**: Minimal impact, recommendation level, no direct harm

# Conclusion

## Key Summary

### 1. SystemOwnerSafe Permission Scope

SystemOwnerSafe controls a **total of 18 contracts** in the Thanos Stack:

- **Indirect Control via ProxyAdmin**: 15 (14 Proxies + AddressManager)
- **Direct Ownership**: 3 (DisputeGameFactory, DelayedWETH × 2)

### 2. Permission Function Risk Classification

The **total of 22 functions** available to SystemOwnerSafe are classified by risk level:

**🔴 High (11 functions)**: Direct fund theft, or complete system halt

- ProxyAdmin: `upgrade()`, `upgradeAndCall()`, `changeProxyAdmin()`, `setAddress()`, `setAddressManager()` (5 functions)
- SystemConfig: `setBatcherHash()` (1 function)
- DisputeGameFactory: `setImplementation()` (1 function)
- DelayedWETH (× 2): `recover()`, `hold()` (4 functions total)

**🟡 Medium (8 functions)**: Indirect economic damage, partial service degradation, recoverable issues

- ProxyAdmin: `setProxyType()`, `setImplementationName()` (2 functions)
- ProtocolVersions: `setRequired()` (1 function)
- SystemConfig: `setUnsafeBlockSigner()`, `setGasConfigEcotone()`, `setGasLimit()` (3 functions)
- DataAvailabilityChallenge: `setBondSize()` (1 function)
- DisputeGameFactory: `setInitBond()` (1 function)

**🟢 Low (3 functions)**: Minimal impact

- ProxyAdmin: `setUpgrading()`
- ProtocolVersions: `setRecommended()`
- DataAvailabilityChallenge: `setResolverRefundPercentage()`

> **Note**: `setGasConfig()` is deprecated after Ecotone upgrade and not counted in the above totals.

## Issues When Registering as DAO Candidate

When an operator deploying L2 through TRH registers as a DAO candidate, **High Risk (🔴) functions** pose risks of manipulating the seigniorage distribution algorithm or stealing user funds, requiring a **secure permission transfer mechanism**.

[[WIP] Brainstorming: How to deal with high-risk functions](https://www.notion.so/WIP-Brainstorming-How-to-deal-with-high-risk-functions-2d8d96a400a380ed9637c36131199bbc?pvs=21)