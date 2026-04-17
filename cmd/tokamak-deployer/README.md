# tokamak-deployer

CLI that drives OP Stack L1 contract deployment for tokamak-thanos and applies
tokamak-specific post-processing to the L2 genesis. The binary embeds every
contract artifact (`//go:embed deploy-artifacts/`) so the L1 deploy itself runs
over plain JSON-RPC without Foundry.

## Commands

| Command | Purpose |
|---------|---------|
| `deploy-contracts` | Deploy the 14 L1 contracts (`AddressManager`, `ProxyAdmin`, 6 proxies + 6 implementations) plus their upgrade / setup calls — **32 on-chain steps total** — and write `deploy-output.json`. |
| `generate-genesis` | Apply tokamak-specific post-processing (DRB inject, USDC inject, MultiTokenPaymaster inject, L1Block Isthmus bytecode patch, rollup hash update) on top of a base L2 genesis. See §3 — **this is not a one-shot; two external steps must run first**. |

## Sepolia prerequisites

| Resource | Recommended |
|----------|-------------|
| L1 RPC | A stable Sepolia endpoint (Alchemy/Infura/public). Public Alchemy works but may rate-limit. |
| Deployer balance | **≥ 0.5 ETH**. A real run at 13-20 Gwei cost **0.35 ETH** for 26 mined txs; the 0.07 ETH figure from the design doc assumes 1 Gwei and is rarely achievable on modern Sepolia. |
| Deployer key | **Exclusive to this deployment** — the deployer pins the nonce sequence, so any other pending tx from the same address will break it. |
| L2 chain ID | Any 32-bit unused ID (e.g. `111551143645`). Used verbatim in `batchInboxAddress` (`0xff00…` + zero-padded decimal). |

## 1. Get the binary

### Option A: pre-built release (recommended)

```bash
VERSION=v0.0.2
PLATFORM=linux-amd64   # or: darwin-arm64 | darwin-amd64 | linux-arm64

curl -L -o tokamak-deployer.tgz \
  "https://github.com/tokamak-network/tokamak-thanos/releases/download/tokamak-deployer/${VERSION}/tokamak-deployer-${PLATFORM}.tar.gz"
tar xzf tokamak-deployer.tgz
./tokamak-deployer --help
```

All `tokamak-deployer/v*` releases: https://github.com/tokamak-network/tokamak-thanos/releases

### Option B: build from source

Required only when you changed contracts or need a custom build.
Go **1.24+** is required.

```bash
cd /path/to/tokamak-thanos

# 1. Compile Solidity → forge-artifacts/
cd packages/tokamak/contracts-bedrock && forge build && cd -

# 2. Extract the 17 artifact JSONs the binary embeds
bash cmd/tokamak-deployer/scripts/extract-artifacts.sh \
  packages/tokamak/contracts-bedrock/forge-artifacts \
  cmd/tokamak-deployer/deploy-artifacts

# 3. Build (CGO off + GOWORK off keeps the binary static & avoids workspace imports)
cd cmd/tokamak-deployer
CGO_ENABLED=0 GOWORK=off go build -o tokamak-deployer .
```

The `deploy-artifacts/` directory must exist before `go build` — it's embedded
at compile time via `//go:embed`. Step 2 is therefore mandatory for a from-source
build.

## 2. Deploy L1 contracts

```bash
export SEPOLIA_RPC="https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY"
export DEPLOYER_KEY="0x<64-hex-chars>"          # no 0x prefix also accepted
export L2_CHAIN_ID=111551143645

./tokamak-deployer deploy-contracts \
  --l1-rpc      "$SEPOLIA_RPC" \
  --private-key "$DEPLOYER_KEY" \
  --chain-id    "$L2_CHAIN_ID" \
  --out         deploy-output.json
```

Typical output (excerpt from a real Sepolia run):

```
[deployer] Starting contract deployment for L2 chain 111551143645
[deployer] Connected to L1 RPC
[deployer] L1 chain ID: 11155111
[deployer] Starting nonce: 1787, deployer address: 0x7220c7...
[deployer] Step 1/32: Deploying AddressManager
[deployer] Suggested gas price: 15 Gwei
[deployer] deploy(nonce=1787, 1334 bytes): broadcasting (attempt 1/5, hash: 0x7a85...)
[deployer] Transaction mined in block 10676204 (status: 1)
[deployer] ✓ AddressManager deployed: 0xbfB3990F2bB579FC18fd04A585Dd6021c9E0d4aE
...
[deployer] Step 32/32: Upgrading L2OutputOracleProxy
[deployer] ✅ All contracts deployed successfully!
```

Step count: **32** (14 `Deploying` steps + 18 `Upgrading`/setup calls).

Expected wall-clock time on Sepolia: **8–15 min** at ~15 Gwei, roughly 26 txs.
`deploy-output.json` contains every Proxy/AddressManager/ProxyAdmin address plus
`l1ChainId`/`l2ChainId` metadata.

### Stuck-transaction recovery (v0.0.2+)

If a tx does not confirm within 90 s, the deployer re-signs the **same nonce**
with a bumped gas price (max of previous × 1.25 and current `eth_gasPrice`)
and re-broadcasts as a replacement. Up to 5 attempts per contract — so ~7.5 min
worst case, and a total gas escalation of ~2.44× from the initial suggestion.

You'll see:

```
[deployer] deploy(nonce=1788, 5442 bytes): tx 0x48d5... not mined within 1m30s, will retry with bumped gas
[deployer] deploy(nonce=1788, 5442 bytes): attempt 2 bumping gas price to 19015977006 wei
[deployer] deploy(nonce=1788, 5442 bytes): broadcasting (attempt 2/5, hash: 0x9fa1...)
```

Observed on a real run: 3 separate txs needed gas bumps, all succeeded on the
2nd attempt.

If all retries fail you get `transaction not mined after 5 attempts (last hash: …)` —
the nonce is still held by the last broadcast, so wait for that tx to clear the
mempool (or drop) before retrying the whole command.

## 3. Generate L2 genesis

> ⚠️ **`generate-genesis` is not standalone.** It runs tokamak-specific
> post-processing on top of a base L2 genesis produced by `op-node genesis l2`,
> which in turn needs an L2 allocs file produced by Foundry's `L2Genesis.s.sol`.
> Two external steps must run first. The `--base-genesis` flag exists precisely
> so you can bring your own base genesis after running those two steps.
>
> (When invoked from trh-sdk, the wrapper runs the forge script and op-node for
> you; the standalone flow documented below mirrors `start-deploy.sh`.)

### 3a. Build the L2 allocs (forge)

Foundry's FFI sandbox only reads files under the contracts-bedrock project
root — **not `/tmp`** — and `L2Genesis.s.sol` tries to parse every key in the
addresses file as an `address`, so you also need a copy without the
`l1ChainId`/`l2ChainId` metadata. Stage the two inputs first:

```bash
REPO=/path/to/tokamak-thanos
CB=$REPO/packages/tokamak/contracts-bedrock
mkdir -p "$CB/deploy-config" "$CB/deployments"

# Strip metadata fields so forge's parseJsonAddress doesn't choke
jq 'del(.l1ChainId, .l2ChainId)' deploy-output.json \
  > "$CB/deployments/${L2_CHAIN_ID}-addresses.json"

# Place deploy-config under the forge project root (FFI sandbox requirement)
cp deploy-config.json "$CB/deploy-config/${L2_CHAIN_ID}.json"
```

Then run the script:

```bash
cd "$CB"
CONTRACT_ADDRESSES_PATH="$CB/deployments/${L2_CHAIN_ID}-addresses.json" \
DEPLOY_CONFIG_PATH="$CB/deploy-config/${L2_CHAIN_ID}.json" \
forge script scripts/L2Genesis.s.sol:L2Genesis \
  --rpc-url "$SEPOLIA_RPC"
```

Output: `$CB/state-dump-${L2_CHAIN_ID}.json` (~7 MB).

### 3b. Build base genesis + rollup (op-node)

```bash
$REPO/op-node/bin/op-node genesis l2 \
  --deploy-config   "$CB/deploy-config/${L2_CHAIN_ID}.json" \
  --l1-deployments  "$CB/deployments/${L2_CHAIN_ID}-addresses.json" \
  --l2-allocs       "$CB/state-dump-${L2_CHAIN_ID}.json" \
  --outfile.l2      genesis-base.json \
  --outfile.rollup  rollup.json \
  --l1-rpc          "$SEPOLIA_RPC"
```

Output: `genesis-base.json` (~9.6 MB), `rollup.json` (~1.4 KB).

### 3c. Apply tokamak post-processing (this binary)

Pass `--base-genesis` so the binary skips its own `op-node` invocation and only
runs the five post-processing steps (DRB inject, USDC inject, MultiTokenPaymaster
inject, L1Block Isthmus bytecode patch, rollup hash update):

```bash
./tokamak-deployer generate-genesis \
  --deploy-output deploy-output.json \
  --config        deploy-config.json \
  --base-genesis  genesis-base.json \
  --out           genesis.json \
  --rollup-out    rollup.json \
  --preset        defi
```

You'll see a short post-processing log and `genesis.json` / `rollup.json` in place.
The rollup's `genesis.l2.hash` is re-derived to match the patched allocs.

`--preset` selects which predeploys ship in the L2 genesis:

| Preset | Extras on top of baseline OP Stack |
|--------|-----------------------------------|
| `general` | (baseline, no extras) |
| `defi`    | CrossTrade contracts, paymaster for non-TON fee tokens |
| `gaming`  | gaming-specific predeploys + DRB inject |
| `full`    | everything |

Examples of `deploy-config.json` live under
[`packages/tokamak/contracts-bedrock/deploy-config/`](../../packages/tokamak/contracts-bedrock/deploy-config/)
— `preset-<name>.json` is the starting point; patch `l2ChainID`,
`l1StartingBlockTag` (current block hash), operator addresses, and
`batchInboxAddress` for your chain.

## Troubleshooting

| Symptom | Cause / fix |
|---------|-------------|
| `Suggested gas price: 0 Gwei` in logs | Display rounding only. `new(big.Int).Div(gasPrice, 1e9)` truncates sub-1-Gwei values to `0` for the log line; the actual tx is signed with the full wei value. |
| `insufficient funds for gas × price + value` | Deployer balance too low. Budget ≥ 0.5 ETH; mainnet-like gas on Sepolia burns ~0.35 ETH. |
| Tx hash accepted but `blockNumber` stays `null` for > 2 min | Mempool stall (priced out, or Alchemy-only mempool). v0.0.2+ retries automatically; on v0.0.1 the process will hang forever — upgrade. |
| `deployment reverted` on step X | Contract constructor failed. Inspect the receipt with the printed tx hash. Often caused by a prior partial run or a bad `--chain-id` collision. |
| Binary hangs with no log output at all | Usually a broken L1 RPC. Sanity-check with `curl … eth_chainId`. |
| `tx already known` on retry | The node saw your previous broadcast at the same gas. Treated as a retry signal and ignored — execution continues. |
| Piping deployer output (`... \| tee` / `\| tail`) produces an empty log | The pipe's reader can die before the deployer flushes, swallowing everything. Redirect to a file directly (`> deploy.log 2>&1`) — don't pipe. |
| `generate-genesis` fails with `Required flag "l1-rpc" not set` | `generate-genesis` calls `op-node` but doesn't pass `--l1-rpc`. Use the 3-step flow above (3a/3b/3c) instead of running `generate-genesis` without `--base-genesis`. |
| `generate-genesis` fails with `missing l2-allocs` | Same root cause — the binary doesn't know about the forge state-dump step. Run 3a → 3b → 3c. |
| `forge script` fails on `parseJsonAddress` | Your `deploy-output.json` contains `l1ChainId`/`l2ChainId` — strip them with `jq 'del(.l1ChainId, .l2ChainId)'` (see §3a). |
| `forge script` fails with `path … is not allowed to be accessed for read operations` | Foundry's FFI sandbox. Place the deploy config + addresses under the contracts-bedrock project root, not `/tmp` (see §3a). |

## Log & debug tips

- Every deploy step writes a `[deployer] Step X/32: <action> <Name>` line. The
  last such line tells you exactly where execution is.
- Every broadcast logs the tx hash before waiting. Paste it into
  https://sepolia.etherscan.io/tx/<hash> to see mempool / block status.
- When run from trh-backend, stdout is tee'd into
  `storage/logs/<stackId>/<timestamp>_deploy-l1-contracts_logs.txt`.
- If you're re-running after a crashed / killed session, first check the pending
  nonce (`cast nonce --rpc-url <rpc> <addr>`) — the old run's last unmined tx is
  still holding a slot and a new run will start at that nonce.

## Related

- [`internal/deployer/contracts.go`](internal/deployer/contracts.go) — deploy loop, `sendAndWaitMined` retry helper
- [`internal/genesis/`](internal/genesis/) — genesis post-processing
- [`scripts/extract-artifacts.sh`](scripts/extract-artifacts.sh) — artifact extraction for `//go:embed`
- [`../../.github/workflows/release-deployer.yml`](../../.github/workflows/release-deployer.yml) — release CI (goreleaser → cross-platform tarballs → GitHub pre-release)
- [`../../packages/tokamak/contracts-bedrock/scripts/start-deploy.sh`](../../packages/tokamak/contracts-bedrock/scripts/start-deploy.sh) — reference for how trh-sdk wires forge + op-node + this binary together
