# YLEX Core

**YLEX (YLX)** — an open, CPU-mineable proof-of-work network for everyday payments and long-horizon value storage.

YLEX Core is the official Go implementation of the YLEX protocol. It is a fork
of the Core Blockchain client (`core-coin/go-core`), rebranded and re-tuned for
the YLEX network: a new genesis with **zero premine**, an on-chain
**treasury and liquidity split**, addresses that carry the `ylx` prefix, and a
**10 second** target block time.

## Key properties

| Property | Value |
| --- | --- |
| Ticker | YLX |
| Consensus | PoDE (proof-of-difficulty, EtHash-family), CPU friendly |
| Target block time | 10 seconds (retarget band 7–13 s) |
| Block reward (steady state) | 0.8 YLX |
| Reward split | 70 % miner / 20 % treasury / 10 % liquidity |
| Bootstrap ramp | 3× decaying linearly to 1× over 30 days |
| Premine | **None.** The genesis allocation is empty. |
| Mainnet chain id | 828 |
| Testnet (Devin) chain id | 829 |
| Address prefix | `ylx` (every YLEX address starts with `796c78` in hex) |

Monetary policy, emission schedule and protocol addresses are documented in
[ECONOMICS.md](ECONOMICS.md).

## Building

Building `ylexd` requires Go 1.22 or later and a C/C++ compiler.

### Linux / macOS

```shell
make gocore
```

or, to build the full suite of utilities:

```shell
make all
```

### Windows

`ylexd` requires mingw to build on Windows.

```shell
choco install mingw
```

## Running

### Mainnet

```shell
ylexd
```

### Testnet (Devin)

```shell
ylexd --devin
```

### Mining

Mining needs no unlocked account — the proof-of-diculty coinbase does not
sign transactions. Point the miner at your payout address:

```shell
ylexd --mine --miner.threads 4 --miner.corebase ylxyouraddress
```

On a modern CPU the network targets roughly one block every ten seconds.

## Addresses

YLEX addresses are 22 bytes, laid out as:

```
prefix (3 bytes, ASCII "ylx" = 796c78)
checksum (1 byte)
payload (18 bytes)
```

Every valid YLEX address therefore begins with `796c78` in hex, which reads as
`ylx` in ASCII. The checksum is computed the same way as upstream ICAN
checksums, over the payload followed by the prefix. Addresses from other
networks (including legacy `cb`/`ab` prefixed ones) are rejected by the client.

The three protocol addresses are:

| Address | Purpose |
| --- | --- |
| `796c7875010000000000000000000000000000000000` | Treasury (20 % of every block reward) |
| `796c7861020000000000000000000000000000000000` | Liquidity (10 % of every block reward) |
| `796c7889000000000000000000000000000000000000` | Genesis coinbase (unspendable) |

## Executables

Command | Description
--- | ---
`ylexd` | The main YLEX CLI client. Full node (default), archive node, or light node. Exposes JSON RPC over HTTP, WebSocket and IPC. (`gocore` upstream name, renamed.)
`abigen` | Generates Go bindings from contract ABI definitions.
`bootnode` | Stripped-down client used for network discovery nodes.
`evm` | Developer utility to run bytecode snippets in a configurable environment.
`ylexkey` | Key management and wallet tool. (`xcbkey` upstream name, renamed.)
`cvm` | Core Virtual Machine — runs smart contracts on YLEX.
`pode` | PoDE mining client.

Renaming of the remaining command-line binaries is in progress; see
[ROADMAP.md](ROADMAP.md).

## Networks

| Network | Chain id | Notes |
| --- | --- | --- |
| Mainnet | 828 | YLEX production network |
| Devin | 829 | Pre-configured proof-of-work test network (`--devin`) |
| Private | any | `--networkid <n>` plus `--datadir` |

Bootnodes for the YLEX networks are being provisioned; until they are
announced, the client ships with an empty bootnode list, so private or
standalone operation is the default.

## Contribute / report security issues

See [SECURITY.md](SECURITY.md) for responsible disclosure, and
[GOVERNANCE.md](GOVERNANCE.md) for how the treasury is run.

## License

Licensed under the [YLEX License v1.0](LICENSE), a source-disclosure license
derived from the Core Blockchain client's CORE License v1.0. YLEX is built on
the work of the Core Blockchain Foundation and the go-ethereum community.
