# Thanos Smart Contracts

This package contains the L1 and L2 contracts and components to build the Thanos. We can use ERC20 token as L2 native token after modifying configuration in <a href="./deploy-config/">deploy-config</a>.

## Directory Structure
<pre>
├── <a href="./deploy-config/">deploy-config</a>: Pre-defined deployment configuration on each network
├── <a href="./deployments/">deployments</a>: Information about the contracts deployed on each network
├── <a href="./genesis/">genesis</a>: Deployed contract address list, L2 Genesis file, Rollup configuration on each network
├── <a href="./uniswap-v3-artifacts/">uniswap-v3-artifacts</a>: Hardhat artifacts for Uniswap V3
├── <a href="./invariant-docs/">invariant-docs</a>: Documentation for all defined invariant tests
├── <a href="./lib/">lib</a>: External libraries with Git submodules
├── <a href="./scripts/">scripts</a>: Deploy Scripts
├── <a href="./src/">src</a>: Contracts
│   ├── <a href="./src/L1/">L1</a>: Contracts deployed on the L1
│   ├── <a href="./src/L2/">L2</a>: Contracts deployed on the L2
│   ├── <a href="./src/cannon/">cannon</a>: Contracts for cannon
│   ├── <a href="./src/dispute/">dispute</a>: Contracts for dispute game
│   ├── <a href="./src/libraries/">libraries</a>: Libraries
│   │    └── <a href="./src/libraries/Predeploys.sol">Predeploys.sol</a>: Pre-deployed contract addresses on L2 Genesis
│   ├── <a href="./src/tokamak-contracts/">tokamak-contracts</a>
│   │    └── <a href="./src/tokamak-contracts/USDC/">USDC</a>: Contract for USDC bridge
│   └── <a href="./src/universal/">universal</a>: Universal contracts
├── <a href="./test/">test</a>: Contracts for unit test
├── <a href="./foundry.toml">foundry.toml</a>: Foundry configuration
├── <a href="./hardhat.config.ts">hardhat.config.ts</a>: Hardhat configuration
...
</pre>

## Contracts Overview

### Contracts deployed to L1

| Name                                                                                     | Proxy Type                                                              | Description                                                                                         |
| ---------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [`L1CrossDomainMessenger`](./src/L1/L1CrossDomainMessenger.sol)                                    | [`ResolvedDelegateProxy`](./src/legacy/ResolvedDelegateProxy.sol) | High-level interface for sending messages to and receiving messages from Thanos                   |
| [`L1StandardBridge`](./src/L1/L1StandardBridge.sol)                                             | [`L1ChugSplashProxy`](./src/legacy/L1ChugSplashProxy.sol)         | Standardized system for transferring ERC20 tokens to/from Thanos                                   |
| [`L2OutputOracle`](./src/L1/L2OutputOracle.sol)             | [`Proxy`](./src/universal/Proxy.sol)                              | Stores commitments to the state of Thanos which can be used by contracts on L1 to access L2 state |
| [`OptimismPortal`](./src/L1/OptimismPortal.sol)                             | [`Proxy`](./src/universal/Proxy.sol)                              | Low-level message passing interface                                                                 |
| [`OptimismMintableERC20Factory`](./src/universal/OptimismMintableERC20Factory.sol) | [`Proxy`](./src/universal/Proxy.sol)                              | Deploys standard `OptimismMintableERC20` tokens that are compatible with either `StandardBridge`    |
| [`ProxyAdmin`](./src/universal/ProxyAdmin.sol)                                                         | -                                                                       | Contract that can upgrade L1 contracts                                                              |

### Contracts deployed to L2

| Name                                                                                     | Proxy Type                                 | Description                                                                                      |
| ---------------------------------------------------------------------------------------- | ------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| [`GasPriceOracle`](./src/L2/GasPriceOracle.sol)                         | [`Proxy`](./src/universal/Proxy.sol) | Stores L2 gas price configuration values                                                         |
| [`L1Block`](./src/L2/L1Block.sol)                                           | [`Proxy`](./src/universal/Proxy.sol) | Stores L1 block context information (e.g., latest known L1 block hash)                           |
| [`L2CrossDomainMessenger`](./src/L2/L2CrossDomainMessenger.sol)             | [`Proxy`](./src/universal/Proxy.sol) | High-level interface for sending messages to and receiving messages from L1                      |
| [`L2StandardBridge`](./src/L2/L2StandardBridge.sol)                         | [`Proxy`](./src/universal/Proxy.sol) | Standardized system for transferring ERC20 tokens to/from L1                                     |
| [`L2ToL1MessagePasser`](./src/L2/L2ToL1MessagePasser.sol)               | [`Proxy`](./src/universal/Proxy.sol) | Low-level message passing interface                                                              |
| [`SequencerFeeVault`](./src/L2/SequencerFeeVault.sol)                       | [`Proxy`](./src/universal/Proxy.sol) | Vault for L2 transaction fees                                                                    |
| [`OptimismMintableERC20Factory`](./src/universal/OptimismMintableERC20Factory.sol) | [`Proxy`](./src/universal/Proxy.sol) | Deploys standard `OptimismMintableERC20` tokens that are compatible with either `StandardBridge` |
| [`ProxyAdmin`](./src/universal/ProxyAdmin.sol)                                                       | -                                          | Contract that can upgrade L2 contracts when sent a transaction from L1                           |

## Installation

We export contract ABIs, contract source code, and contract deployment information for this package via `npm`:

```shell
npm install @tokamak-network/thanos-contracts
```

## Deployment

The smart contracts are deployed using `foundry` with a `hardhat-deploy` compatibility layer. When the contracts are deployed,
they will write a temp file to disk that can then be formatted into a `hardhat-deploy` style artifact by calling another script.

### Configuration

Create or modify a file `<network-name>.json` inside of the [`deploy-config`](./deploy-config/) folder.
By default, the network name will be selected automatically based on the chainid. Alternatively, the `DEPLOYMENT_CONTEXT` env var can be used to override the network name.
The `IMPL_SALT` env var can be used to set the `create2` salt for deploying the implementation contracts.

### Execution

1. Set the env vars `ETH_RPC_URL`, `PRIVATE_KEY` and `ETHERSCAN_API_KEY` if contract verification is desired
1. Deploy the contracts with `forge script -vvv scripts/Deploy.s.sol:Deploy --rpc-url $ETH_RPC_URL --broadcast --private-key $PRIVATE_KEY`
   Pass the `--verify` flag to verify the deployments automatically with Etherscan.
1. Generate the hardhat deploy artifacts with `forge script -vvv scripts/Deploy.s.sol:Deploy --sig 'sync()' --rpc-url $ETH_RPC_URL --broadcast --private-key $PRIVATE_KEY`

### Deploying a single contract

All of the functions for deploying a single contract are `public` meaning that the `--sig` argument to `forge script` can be used to
target the deployment of a single contract.

## Static Analysis

`contracts-bedrock` uses [slither](https://github.com/crytic/slither) as its primary static analysis tool. Slither will be run against PRs as part of CI, and new findings will be reported as a comment on the PR.