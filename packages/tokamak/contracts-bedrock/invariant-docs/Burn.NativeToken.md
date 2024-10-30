# `Burn.NativeToken` Invariants

## `nativeToken(uint256)` always burns the exact amount of native token passed.
**Test:** [`Burn.NativeToken.t.sol#L66`](../test/invariants/Burn.NativeToken.t.sol#L66)

Asserts that when `Burn.nativeToken(uint256)` is called, it always burns the exact amount of native token passed to the function.
