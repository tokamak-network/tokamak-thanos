
# @tokamak-network/thanos-sdk

The `@tokamak-network/thanos-sdk` package provides a set of tools for interacting with Thanos.

## Installation

```
npm install @tokamak-network/thanos-sdk
```

## Docs

You can find auto-generated API documentation over at [tokamak-network.github.io/thanos-sdk-docs](https://tokamak-network.github.io/thanos-sdk-docs/).

## Using the SDK

### CrossChainMessenger

The [`CrossChainMessenger`](https://github.com/tokamak-network/tokamak-thanos/blob/main/packages/tokamak/sdk/src/cross-chain-messenger.ts) class simplifies the process of moving assets and data between Ethereum and Optimism.
You can use this class to, for example, initiate a withdrawal of ERC20 tokens from Optimism back to Ethereum, accurately track when the withdrawal is ready to be finalized on Ethereum, and execute the finalization transaction after the challenge period has elapsed.
The `CrossChainMessenger` can handle deposits and withdrawals of ETH and any ERC20-compatible token.
Detailed API descriptions can be found at [tokamak-network.github.io/thanos-sdk-docs](https://tokamak-network.github.io/thanos-sdk-docs/classes/cross_chain_messenger.CrossChainMessenger.html).
The `CrossChainMessenger` automatically connects to all relevant contracts so complex configuration is not necessary.

### L2Provider and related utilities

The Optimism SDK includes [various utilities](https://github.com/tokamak-network/tokamak-thanos/blob/main/packages/tokamak/sdk/src/l2-provider.ts) for handling Optimism's [transaction fee model](https://community.optimism.io/docs/developers/build/transaction-fees/).
For instance, [`estimateTotalGasCost`](https://tokamak-network.github.io/thanos-sdk-docs/functions/l2_provider.estimateL2GasCost.html) will estimate the total cost (in wei) to send at transaction on Optimism including both the L2 execution cost and the L1 data cost.
You can also use the [`asL2Provider`](https://tokamak-network.github.io/thanos-sdk-docs/functions/l2_provider.asL2Provider.html) function to wrap an ethers Provider object into an `L2Provider` which will have all of these helper functions attached.

### Other utilities

The SDK contains other useful helper functions and constants.
For a complete list, refer to the auto-generated [SDK documentation](https://tokamak-network.github.io/thanos-sdk-docs/)
