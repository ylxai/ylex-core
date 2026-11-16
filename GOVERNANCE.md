# YLEX Governance

## What is on chain, and what is not

Everything that makes YLEX a functioning network is enforced by consensus
rules: the block reward, the 70 / 20 / 10 split, the bootstrap ramp, and the
protocol addresses that receive the treasury and liquidity shares. No
governance process is required for the chain to keep running, and none can
change these rules without a coordinated network upgrade.

Governance exists for exactly one thing: deciding how the treasury and
liquidity funds are spent.

## The treasury and liquidity accounts

The treasury (`796c7875010000000000000000000000000000000000`) receives 20 % of
every block reward. The liquidity account
(`796c7861020000000000000000000000000000000000`) receives 10 %. Both are
ordinary accounts with no protocol privileges.

Keys for both accounts are held by a multisig controlled by the **YLEX
council**, a five-seat body:

| Seat | Mandate |
| --- | --- |
| 1 – 2 | Founding maintainers of the YLEX protocol |
| 3 | Treasury and finance |
| 4 | Community and ecosystem |
| 5 | Independent member, chosen by the other four, replaced annually |

Spending thresholds:

| Amount | Approval required |
| --- | --- |
| ≤ 500 YLX | 2 of 5 seats |
| ≤ 5,000 YLX | 3 of 5 seats |
| > 5,000 YLX | 4 of 5 seats, with a published rationale |

## Transparency rules

The council is bound by the same rules the chain is: every spend is published,
and the multisig threshold is a property of the account, not a promise.

- Every treasury or liquidity movement is published within 7 days with the
  destination, amount and purpose.
- A quarterly report lists all outflows since the previous report.
- Any seat holder may publish a dissenting opinion on any spend, and it must
  be included in the report verbatim.

## Liquidity mandate

The liquidity account exists to make YLEX tradable where people already are:
providing or pairing liquidity on exchanges and swap venues, funding market
making during the bootstrap phase, and covering listing costs. The council
does not trade the liquidity account for profit; unused liquidity funds roll
over.

## Protocol changes

Changes to consensus rules — reward, block time, split percentages, genesis
constants — require:

1. A public specification and reference implementation,
2. A minimum 30-day comment period,
3. A supermajority (4 of 5) council sign-off on the release that activates it,
   plus adoption by miners executing the upgrade.

Governance cannot rewrite history, freeze accounts, or confiscate balances;
the protocol has no such capability, and adding one is out of scope for any
proposal.
