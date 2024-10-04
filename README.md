THis PR is for analyzing gas cost in L1CrossDomainMessenger.relayMessage.

Testing method:
- Run Devnet
- Unittest

What we want here is the call L1CrossDomainMessenger.relayMessage from OptimismPortal cannot be reverted.
If the call reverts, token will be locked in OptimismPortal and relayMessage cannot be successful. It could happen if L1CrossDomainMessenger.relayMessage run out-of-gas in its context. By this test, we will check how many gas is used for relayMessage()

The tests show that relayMessage only takes less than 100K gas and gas cost after the external call is executed is ~25K gas (RELAY_RESERVED_GAS = 40_000). So we don't have any issues relates to gas here.

However, RELAY_GAS_CHECK_BUFFER needs to be increased to 35_000 from 5000 because the gas cost after SafeCall.hasMinGas until the external call is increasing by ~27K gas (The call still succeeds because the total gas provided by OptimismPortal ~300K gas)

Let explain why there is not much extra gas in relayMessage? Because Thanos has used gas optimization techniques, such as resetting storage values to their original state, or clear the storage values or leverage warm accessed objects.

For example the code below:

```
    if (_value != 0 && _target != address(0)) {
        IERC20(_nativeTokenAddress).approve(_target, 0);
    }
```
- The contract L2 native token is called again so it is a warm access.
- And we try to clear the approval storage. If the storage value is 0 already, the gas cost is low.

We can visit the doc that explain opcode SSTORE at https://ethereum.org/en/developers/docs/evm/opcodes/
And: https://github.com/wolflo/evm-opcodes/blob/main/gas.md#a7-sstore

