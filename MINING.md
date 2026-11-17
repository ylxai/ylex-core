# Mining YLEX

YLEX is CPU-mineable proof-of-work. This guide covers solo mining on the test
network (Devin) and what to expect on mainnet.

## What you need

- A 64-bit Linux, macOS or Windows machine with at least 4 cores.
- The `ylexd` binary (see the README to build it, or use a release build).
- A YLEX address to receive rewards. Any address works; create one with:

```shell
ylexd account new
```

The address is printed in the `ylx...` format (hex starts with `796c78`).

## Mine on the test network (Devin)

```shell
ylexd --devin \
  --datadir ~/.ylex-devin \
  --mine \
  --miner.threads 4 \
  --miner.corebase <your-address>
```

- `--miner.threads` sets how many CPU threads search for a seal. Start with
  one thread per physical core.
- `--miner.corebase` is the address that receives the miner share of every
  block you find. Mining does not require unlocking the account.
- The node finds its own difficulty. On a fresh network the first blocks may
  come faster or slower than 10 seconds until the retarget settles.

Check progress:

```shell
ylexd attach --datadir ~/.ylex-devin
> xcb.blockNumber
> xcb.getBalance("<your-address>")
```

## Rewards

Every block pays three accounts at once:

| Recipient | Share | Steady state (0.8 YLX block) |
| --- | --- | --- |
| Miner | 70 % | 0.56 YLX |
| Treasury | 20 % | 0.16 YLX |
| Liquidity | 10 % | 0.08 YLX |

During the first 30 days a bootstrap multiplier starts at 3× and decays
linearly to 1×; see [ECONOMICS.md](ECONOMICS.md) for the full schedule.

## Estimating income

At steady state there are about 8,640 blocks per day. If your hashrate is
`h` and the network hashrate is `H`, your expected daily income is:

```
0.56 YLX * 8640 * (h / H)  =  4,838 YLX * (h / H)
```

During the bootstrap phase multiply by up to 3. Because difficulty adjusts to
total network hashrate, the only number that matters is your share of the
network, not your absolute hashrate.

## Running a pool-friendly setup

The client exposes the same mining interface as the upstream Core client, so
standard PoDE pool software works with it. If you operate a pool, point it at
a `ylexd --mine` node with `--miner.corebase` set to the pool's payout
address.

## Troubleshooting

| Symptom | Cause |
| --- | --- |
| `corebase must be explicitly specified` | pass `--miner.corebase <address>` |
| no blocks for a long time | difficulty is still settling, or peer count is zero (check `--verbosity 3` logs) |
| `invalid network id prefix in address` | the address is from another chain; YLEX addresses start with `796c78` |
| CPU at 100 % but no seals | expected; sealing is a search. Watch `Successfully sealed new block` lines |
