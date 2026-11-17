# YLEX Roadmap

The roadmap is organised as phases with concrete, verifiable exit criteria. It
is a plan, not a promise; it will be updated as the network produces real
data.

## Phase 0 — Rebrand and tooling (current)

- [x] Fork the Core Blockchain Go client and rebrand the module, build and
      address prefix to YLEX
- [x] New genesis: chain ids 828 / 829, zero premine, on-chain treasury and
      liquidity split
- [x] Monetary policy in consensus code: 0.8 YLX steady state, 70 / 20 / 10,
      30-day bootstrap ramp
- [x] 10 second target block time with a 7–13 s retarget band, retarget
      divisor 512
- [x] Verified end to end on the test network: blocks mined, split paid
      exactly 7 : 2 : 1
- [x] `ylx` address prefix and rebranded documentation
- [x] Rename remaining binaries (`gocore` → `ylexd`, `xcbkey` → `ylexkey`)
- [x] Remove the inherited dead genesis allocation and let the test build
      suite pass cleanly

## Phase 1 — Public test network

- [ ] Publish bootnodes, seed a public Devin network
- [ ] Faucet for test YLX
- [ ] Block explorer (Blockscout or a lightweight indexer)
- [ ] Live difficulty and orphan-rate dashboard
- [x] Miner documentation: solo mining, pool mining, hardware expectations
      ([MINING.md](MINING.md))

## Phase 2 — Mainnet launch

- [ ] Final genesis image, published with reproducible build hashes
- [ ] Mining client binaries for Linux, Windows, macOS
- [ ] Listing on at least one community exchange
- [ ] Runtime for the treasury multisig (council tooling, published address
      list)

## Phase 3 — Ecosystem

- [ ] Wallet support with the `ylx` address prefix
- [ ] Bridge or wrapped-YLX on one major chain
- [ ] Contract-hosting guidance (the CVM is fully EVM-compatible)
- [ ] Self-hosted mining pool option (open source) so small miners do not
      depend on third-party pools

## Phase 4 — Protocol review

- [ ] Publish six months of propagation, orphan and difficulty data
- [ ] Evaluate shortening the target block time from 10 s toward 5 s; adopt
      only if measured orphan rate stays low
- [ ] Evaluate the treasury share once the bootstrap phase is over
- [ ] Decide whether the council size should grow from 5 to 7

## Long-term principles

1. **No premine, ever.** Founders hold only mined coins.
2. **No hard supply cap.** Issuance stays predictable and auditable in
   `reward.go`; inflation declines as a rate because issuance is constant.
3. **CPU mineable for as long as possible.** Any proposal that requires
   specialised hardware is rejected as a matter of policy.
4. **Changes are published before they are activated.**
