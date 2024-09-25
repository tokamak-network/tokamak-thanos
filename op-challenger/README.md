# Tokamak Thanos Local Devnet FPS Operation Guide
## op-challenger

The `op-challenger` is a modular **op-stack** challenge agent written in Go for various dispute games like attestation, fault, and validity games. For more details, visit the [fault proof specs][proof-specs].

[proof-specs]: https://specs.optimism.io/experimental/fault-proof/index.html
## 1. Clone the Monorepo
To clone the Tokamak Thanos repository, use the following command:

```bash
git clone https://github.com/tokamak-network/tokamak-thanos.git
```
<br/>

## 2. Running with Cannon on Local Devnet
To run the `op-challenger` on the local devnet, begin by cleaning and starting the devnet from the root of the repository:

```bash
make devnet-clean
make devnet-up
```

> **Note:** After the FPS update, the `ops-bedrock-op-challenger-1` container is automatically created, meaning the challenger node is ready for use.

<br/>

## 3. Running the op-challenger Manually (Optional)
If you prefer to run your own challenger node for interactive proof (attack/defense), use the following configuration with `/tokamak-thanos/op-challenger/docker-compose.yml`:

```yaml
version: '3.8'

services:
  challenger:
    image: us-docker.pkg.dev/oplabs-tools-artifacts/images/op-challenger:latest
    volumes:
      - "./challenger-data:/data"
      - "../op-program/bin:/op-program"
    environment:
      OP_CHALLENGER_L1_ETH_RPC: http://l1:8545
      OP_CHALLENGER_L1_BEACON: 'unset'
      OP_CHALLENGER_ROLLUP_RPC: http://op-node:8545
      OP_CHALLENGER_TRACE_TYPE: cannon,fast
      OP_CHALLENGER_GAME_FACTORY_ADDRESS: "0x11c81c1A7979cdd309096D1ea53F887EA9f8D14d"
      OP_CHALLENGER_UNSAFE_ALLOW_INVALID_PRESTATE: "true"
      OP_CHALLENGER_DATADIR: /db
      OP_CHALLENGER_CANNON_ROLLUP_CONFIG: ./.devnet/rollup.json
      OP_CHALLENGER_CANNON_L2_GENESIS: ./.devnet/genesis-l2.json
      OP_CHALLENGER_CANNON_BIN: ./cannon/bin/cannon
      OP_CHALLENGER_CANNON_SERVER: /op-program/op-program
      OP_CHALLENGER_CANNON_PRESTATE: /op-program/prestate.json
      OP_CHALLENGER_L2_ETH_RPC: http://l2:8545
      OP_CHALLENGER_PRIVATE_KEY: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
      OP_CHALLENGER_NUM_CONFIRMATIONS: 1
    networks:
      - ops-bedrock_default  

networks:
  ops-bedrock_default:
    external: true  
```
> **Note:** This configuration sets up the `op-challenger` service for Optimism dispute games, allowing interaction with L1 and L2 Ethereum networks.


### Key Details:
- **Image**: Uses the latest `op-challenger` Docker image from oplabs.
- **Volumes**: 
  - `challenger-data` stores local challenger data.
  - `op-program` contains the binary for the fault proof system.
- **Environment Variables**:
  - **L1 and L2 RPC Endpoints**: Defines connections to L1 (`http://l1:8545`) and L2 (`http://op-node:8545`).
  - **Cannon Configuration**: Specifies files for Cannon VM setup, including rollup config, genesis, and prestate.
  - **Game Factory Address**: Connects to the contract for managing dispute games.
  - **Private Key**: Used for signing transactions.
  - **Trace Type**: Set to use Cannon for fast trace validation.
- **Networks**: The service uses the `ops-bedrock_default` external network for connectivity.


To run the challenger:

```bash
cd op-challenger
docker-compose up
```
<br/>

## 4. Building the op-program
Navigate to the `op-program` directory and run the following commands:

```bash
cd op-program
make
```

After building, two important output files will be located in the `bin` folder:

- `op-program`: Executable file intended for use on the host machine.
- `op-program-client.elf`: MIPS-compiled executable used inside a MIPS emulator or virtual machine as part of Optimism's fault-proof system. It generates prestate and proof data necessary for validating claims in Optimism’s fault dispute game.
> Note: Cannon System use Prestate and Proof Generation, Claim Verification, and Dispute Resolution
<br/>

## 5. Generate Prestate and Proof Files

To generate the prestate and proof files for the Optimism Cannon system, follow these steps:

```bash
make cannon-prestate
cd cannon && make
```

### Prestate Generation and Verification

Once the commands are executed, the following steps will be performed automatically:
The output is verified using readelf on the op-program-client.elf file.
```bash
readelf -h bin/op-program-client.elf
```

Next, rename the generated JSON files: The generated 0.json file is renamed to prestate-proof.json for proper identification.

```bash
mv op-program/bin/0.json op-program/bin/prestate-proof.json
```

### Summary of the Process

- This process compiles the `op-program` into a MIPS architecture executable (`op-program-client.elf`).
- The compiled program is loaded using Cannon and run inside a MIPS-based virtual machine.
- This execution generates two key files:
  - **prestate.json**: Represents the initial state.
  - **prestate-proof.json**: Corresponding proof file for the generated state.

These files are critical for establishing the initial state and proof for Optimism’s Cannon-based fault dispute resolution system.

---
<br/>

## 6. Verify Claims with `op-program` (Challenger Side)

To verify claims on the challenger side using `op-program`, use the following command:

```bash
./op-program \
  --l1 http://127.0.0.1:8545 \   # L1 Ethereum RPC URL
  --l2 http://127.0.0.1:9545 \   # L2 RPC URL
  --l1.head 0x2d21d3cbe1d3042263cf7ff8ec834c37ff693be8fdeef3b14d464d3036ccb808 \ # Latest L1 block hash
  --l2.head 0xe79893e3e165e1f1bcdc7482fe8ce282cdc207ed9f5e5708538200b8a48e5de4 \ # Latest L2 block hash
  --l2.outputroot 0x7e99eb01431a318395e9969a4bc84db73b7d7032a3f749a39e2e71d89ccbd6f0 \ # Previous block’s output root
  --l2.claim 0x00685eeb8cd764a60a58784ba10f56743be033c1990ba65bee35e0ffed73bfc2 \ # Claimed output root of the block
  --l2.blocknumber 18   # The block number being validated
```

### Output Root Verification

To verify the output root of a specific block, use the following commands:

```bash
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"optimism_outputAtBlock","params":["0x11"],"id":1}' \
  http://localhost:7545

curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"optimism_outputAtBlock","params":["0x12"],"id":1}' \
  http://localhost:7545

curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", true],"id":1}' \
  -H "Content-Type: application/json" \
  http://127.0.0.1:8545
```

These requests allow you to fetch the output root and block information for further verification of the challenger claims.

---
<br/>

## 7. Verifying Claims with Cannon

### Fault Proof Virtual Machine (FPVM) State Validation

Using Cannon, you can validate the state of the Fault Proof Virtual Machine (FPVM) by checking the VM’s state after executing specific instructions:

```bash
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"optimism_outputAtBlock","params":["0x11"],"id":1}' \
  http://localhost:7545

curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"optimism_outputAtBlock","params":["0x12"],"id":1}' \
  http://localhost:7545

curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", true],"id":1}' \
  -H "Content-Type: application/json" \
  http://127.0.0.1:8545
```

### Cannon Run Command for Validation

To validate the claim using Cannon, execute the following:

```bash
./bin/cannon run \
  --pprof.cpu \
  --info-at '%10000000' \
  --proof-at never \
  --input ./state.json \
  -- \
  ../op-program/bin/op-program \
  --l2.genesis /path/to/genesis-l2.json \
  --rollup.config /path/to/rollup.json \
  --l1.trustrpc \
  --l1.rpckind debug_geth \
  --l1 http://127.0.0.1:8545 \
  --l2 http://127.0.0.1:9545 \
  --l1.head 0x1efe8ddeb25e41ce2d37afbde433b82a1e3ad5a117f28da51a3130382b22552b \
  --l2.head 0x7b6449d1b2f2eb41aa74a11e7c415d19bc0e9040cee5d9c5a181beba4e48cfb9 \
  --l2.outputroot 0x315d68d1a075bf74cf549d0eb741546ff83f2a927278bfeb4c9a4a919e2cf09c \
  --l2.claim 0x4f65f4c41291764702b14a1ba4ffeb4fab19025791609109a86a7d73e67131ad \
  --l2.blocknumber 18 \
  --datadir /tmp/fpp-database \
  --log.format terminal \
  --server
```

### Verifying the Execution Steps

During the Cannon run, several details are printed to track and verify the execution progress. These key steps ensure that the correct program instructions are being executed, and the state transitions are valid:

- **Step**: The current instruction number being executed (e.g., 10,000,000). This helps you trace which instruction is being processed in the program.
- **PC (Program Counter)**: The memory address of the current instruction. This helps you understand where the program is in memory.
- **Insn**: The actual machine-level instruction being executed at the moment.
- **IPS(Instructions Per Second)**: The rate at which instructions are being processed. This is useful for performance monitoring.
- **Pages**: The number of memory pages used, which helps track memory consumption.
- **Mem**: Total memory usage at that point in the run.
- **Name**: The function or operation currently being executed. Knowing which operation is being processed can help in debugging.

By monitoring these values, you can ensure that the target block’s batch and the transactions within it have been properly fetched, executed, and validated according to the rollup's rules.

---
<br/>

## 8. Creating Proof for Dispute Resolution

When a dispute arises regarding the correctness of a step or transaction, you may need to generate proof of specific steps to resolve the dispute. To generate proof for a specific instruction step, use the following command:

```bash
./bin/cannon run --proof-at '=2792' --stop-at '=2793' --input state.json
```

### Explanation of the Command:

- `--proof-at '=2792'`: Specifies the exact instruction (2792) where proof should be generated. This allows you to provide verifiable evidence of the correctness of this specific instruction.
- `--stop-at '=2793'`: Stops the execution right after this instruction. This is useful to isolate the step being disputed.
- `--input state.json`: Loads the state from a saved JSON file, which contains the pre-existing state of the rollup chain.

### Logs Generated by the Command:

- **Step**: The index of the instruction being executed. This helps you identify where in the program the proof is generated.
- **Pre-State**: The state before the step is executed. This includes information such as memory values, registers, etc., before the instruction runs.
- **Post-State**: The state after the step is executed, showing the changes made by the instruction.
- **State Data**: Detailed state information such as memory layout, register values, etc., for a complete picture of the program’s state.
- **Proof Data**: Information proving the correctness of the step. This proof is used for on-chain verification and dispute resolution, ensuring that the execution follows the rules of the rollup and is verifiable on-chain.

<br/>

## Subcommands Overview

### create-game
Starts a new fault dispute game to challenge the latest output proposal.

```shell
./bin/op-challenger create-game \
  --l1-eth-rpc <L1_ETH_RPC> \
  --game-address <GAME_FACTORY_ADDRESS> \
  --output-root <OUTPUT_ROOT> \
  --l2-block-num <L2_BLOCK_NUM> \
  <SIGNER_ARGS>
```

### move
Makes a move (attack or defend) in an existing game.

```shell
./bin/op-challenger move \
  --l1-eth-rpc <L1_ETH_RPC> \
  --game-address <GAME_ADDRESS> \
  --attack | --defend \
  --parent-index <PARENT_INDEX> \
  --claim <CLAIM> \
  <SIGNER_ARGS>
```

### resolve-claim
Resolves a specific claim within a dispute game.

```shell
./bin/op-challenger resolve-claim \
  --l1-eth-rpc <L1_ETH_RPC> \
  --game-address <GAME_ADDRESS> \
  --claim <CLAIM_INDEX> \
  <SIGNER_ARGS>
```

### resolve
Resolves a full dispute game once eligible.

```shell
./bin/op-challenger resolve \
  --l1-eth-rpc <L1_ETH_RPC> \
  --game-address <GAME_ADDRESS> \
  <SIGNER_ARGS>
```

### list-games
Lists all dispute games created by the game factory.

```shell
./bin/op-challenger list-games \
  --l1-eth-rpc <L1_ETH_RPC> \
  --game-factory-address <GAME_FACTORY_ADDRESS>
```

### list-claims
Displays current claims in a specific dispute game.

```shell
./bin/op-challenger list-claims \
  --l1-eth-rpc <L1_ETH_RPC> \
  --game-address <GAME_ADDRESS>
```
