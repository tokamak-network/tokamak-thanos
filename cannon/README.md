<!--![cannon](https://upload.wikimedia.org/wikipedia/commons/8/80/Cannon%2C_Château_du_Haut-Koenigsbourg%2C_France.jpg)-->
<!--![cannon](https://cdn1.epicgames.com/ue/product/Featured/SCIFIWEAPONBUNDLE_featured-894x488-83fbc936b6d86edcbbe892b1a6780224.png)-->
<!--![cannon](https://static.wikia.nocookie.net/ageofempires/images/8/80/Bombard_cannon_aoe2DE.png/revision/latest/top-crop/width/360/height/360?cb=20200331021834)-->
<!--![cannon](https://paradacreativa.es/wp-content/uploads/2021/05/Canon-orbital-GTA-01.jpg)-->

---

Cannon *(cannon cannon cannon)* is an onchain MIPS instruction emulator.
Cannon supports EVM-equivalent fault proofs by enabling Geth to run onchain,
one instruction at a time, as part of an interactive dispute game.

* It's Go code
* ...that runs an EVM
* ...emulating a MIPS machine
* ...running compiled Go code
* ...that runs an EVM

For more information, see [Docs](./docs/README.md).

## Usage

```shell
# Build op-program server-mode and MIPS-client binaries.
cd ../op-program
make op-program # build

# Switch back to cannon, and build the CLI
cd ../cannon
make cannon

# Transform MIPS op-program client binary into first VM state.
# This outputs state.json (VM state) and meta.json (for debug symbols).
./bin/cannon load-elf --path=../op-program/bin/op-program-client.elf

# Run cannon emulator (with example inputs)
# Note that the server-mode op-program command is passed into cannon (after the --),
# it runs as sub-process to provide the pre-image data.
#
# Note:
#  - The L2 RPC is an archive L2 node on OP goerli.
#  - The L1 RPC is a non-archive RPC, also change `--l1.rpckind` to reflect the correct L1 RPC type.
./bin/cannon run \
    --pprof.cpu \
    --info-at '%10000000' \
    --proof-at never \
    --input ./state.json \
    -- \
    ../op-program/bin/op-program \
    --l1 http://127.0.0.1:8545 \
    --l2 http://127.0.0.1:9545 \
    --l1.head 0x9365329285630a16702cb6f9a8928a103bda9c514de641c9bfa088444d388e68 \
    --l2.head 0x726f3bf1688fa0f25c08d6649140f0c3087eb64c940748d4c6052cd25ad0efcf \
    --l2.outputroot 0xd6ef62caf6bc52b7d522b62867673cada7d64972c61b5557af68571b66f421d3 \
    --l2.claim 0xfcb85c6fd7537626bdc4fb16ad34049773b12ee0a1ee94e03cf1a13d6eab1bb2 \
    --l2.blocknumber 17 \
    --datadir /tmp/fpp-database \
    --log.format terminal \
    --rollup.config /home/nguyen/Projects/tokamak/tokamak-thanos/.devnet/rollup.json \
    --l2.genesis /home/nguyen/Projects/tokamak/tokamak-thanos/.devnet/rollup.json \
    --server

# Add --proof-at '=12345' (or pick other pattern, see --help)
# to pick a step to build a proof for (e.g. exact step, every N steps, etc.)

# Also see `./bin/cannon run --help` for more options
```

## Contracts

The Cannon contracts:
- `MIPS.sol`: A MIPS emulator implementation, to run a single instruction onchain, with merkleized VM memory.
- `PreimageOracle.sol`: implements the pre-image oracle ABI, to support the instruction execution pre-image requests.

The smart-contracts are integrated into the Optimism monorepo contracts:
[`../packages/contracts-bedrock/src/cannon`](../packages/contracts-bedrock/src/cannon)

## `mipsevm`

`mipsevm` is Go tooling to test the onchain MIPS implementation, and generate proof data.

## `example`

Example programs that can be run and proven with Cannon.
Optional dependency, but required for `mipsevm` Go tests.
See [`example/Makefile`](./example/Makefile) for building the example MIPS binaries.

## License

MIT, see [`LICENSE`](./LICENSE) file.

**Note: This code is unaudited.**
In NO WAY should it be used to secure any monetary value before testing and auditing.
This is experimental software, and should be treated as such.
The authors of this project make no guarantees of security of ANY KIND.
