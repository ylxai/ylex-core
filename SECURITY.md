# Security Policy

## Reporting a vulnerability

Please do **not** open a public issue for a security vulnerability.

Report it privately to the maintainers via GitHub's private vulnerability
reporting on `ylxai/ylex-core` (Security → Report a vulnerability). Include:

- the affected version or commit,
- a description of the issue and its impact,
- a proof of concept or reproduction steps if available.

We aim to acknowledge reports within 72 hours and to publish a coordinated
disclosure once a fix is available.

## Scope

In scope:

- consensus bugs: reward calculation, difficulty adjustment, block
  validation, chain reorganisation
- transaction validity, signature handling and the transaction pool
- the EVM/CVM implementation (`core/vm`)
- P2P protocol handling
- wallet, keystore and key handling code
- JSON-RPC endpoints exposed by `ylexd`

Out of scope:

- issues in third-party forks or services
- denial of service through rate-unlimited public RPC endpoints operated by
  others
- anything that requires physical access to a node

## Supported versions

Security fixes land on `master` and are backported to the latest release tag.
Older releases are not supported.

## What we ask of node operators

- Keep `ylexd` up to date.
- Do not expose the personal, admin or debug RPC namespaces to the internet.
- Back up keystore files and, if you run it, the treasury multisig seed.

## Bounties

There is no bug bounty programme at launch. The council may allocate treasury
funds to reward responsible disclosure; such a payment is discretionary and
published under the transparency rules in [GOVERNANCE.md](GOVERNANCE.md).
