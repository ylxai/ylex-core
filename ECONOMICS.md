# YLEX Economics

This document describes the monetary policy of the YLEX network as enforced by
the consensus engine itself. These rules are code, not promises: the block
reward, the split and the addresses they pay are part of
`consensus/cryptore`, verified by every full node.

## Monetary policy at a glance

| Parameter | Value | Source |
| --- | --- | --- |
| Target block time | 10 seconds | `makeDifficultyCalculator` |
| Base block reward (steady state) | 0.8 YLX | `BaseBlockReward` |
| Miner share | 70 % (7000 BPS) | `MinerShareBPS` |
| Treasury share | 20 % (2000 BPS) | `TreasuryShareBPS` |
| Liquidity share | 10 % (1000 BPS) | `LiquidityShareBPS` |
| Premine / genesis allocation | **0 YLX** | empty `GenesisAlloc` |
| Maximum supply | None (no hard cap) | — |

The split is applied on chain inside `accumulateRewards`. Every block credits
three accounts:

| Account | Address | Share |
| --- | --- | --- |
| Block miner | the address that sealed the block | 70 % |
| Treasury | `796c7875010000000000000000000000000000000000` | 20 % |
| Liquidity | `796c7861020000000000000000000000000000000000` | 10 % |

Both protocol addresses are derived from fixed payloads (tags `0x01` and
`0x02`) and are ordinary accounts: their balances are public and verifiable on
chain, and nothing in the protocol gives them special privileges.

## Bootstrap ramp

For the first months the reward is multiplied to seed mining participation:

| Phase | Blocks | Wall clock (10 s blocks) | Multiplier | Reward per block |
| --- | --- | --- | --- | --- |
| Full rate | 0 – 1,440 | first ~4 hours | 3.0× | 2.4 YLX |
| Linear taper | 1,441 – 259,200 | ~4 hours – 30 days | 3.0× → 1.0× | 2.4 → 0.8 YLX |
| Steady state | 259,201+ | day 30 onwards | 1.0× | 0.8 YLX |

The taper is linear in block height. `SplitReward` always pays the miner,
treasury and liquidity in exact proportion, in every phase.

## Emission schedule

At the 10 second target there are 3,153,600 blocks per year.

| Period | Blocks | Miner (70 %) | Treasury (20 %) | Liquidity (10 %) | Total |
| --- | --- | --- | --- | --- | --- |
| Per block (steady state) | 1 | 0.56 YLX | 0.16 YLX | 0.08 YLX | 0.8 YLX |
| Per day | 8,640 | 4,838 YLX | 1,382 YLX | 691 YLX | 6,912 YLX |
| Per year | 3,153,600 | 1,766,016 YLX | 504,576 YLX | 252,288 YLX | 2,522,880 YLX |

Yearly figures exclude the 30-day bootstrap ramp, which adds roughly
208,000 YLX in total before steady state begins.

There is no hard supply cap. Emission is continuous and predictable: anyone
can compute the exact supply at any future height from the constants in
`consensus/cryptore/reward.go`. Inflation is approximately 100 % of the
year-1 issuance each year (constant issuance against a growing supply base),
meaning the inflation *rate* declines over time.

## Why these numbers

- **0.8 YLX per block at 10 seconds** keeps annual issuance (~2.52M YLX) in
  the same band as the original 60 s / 5 YLX design (~2.63M YLX), while giving
  users one-minute confirmations instead of six-minute ones.
- **70 / 20 / 10** keeps the miner incentive dominant — the chain only lives if
  mining is worth doing — while guaranteeing the treasury and liquidity budgets
  are funded by code rather than by donations.
- **No premine** means the founders' only YLEX comes from mining on the same
  terms as everyone else.

## Block time evaluation

The 10 second target was chosen as a conservative, proven operating point
(upstream XCB's retarget formula was designed for a 5–10 second band).
Shorter intervals are possible in principle; the roadmap defers any move
toward 5 seconds until real propagation and orphan-rate data exists. See
[ROADMAP.md](ROADMAP.md).
