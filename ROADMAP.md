# Roadmap

Dependency-ordered, not date-driven. This separates what is **shipped** from
what is **planned**. Planned items are not implemented and must not be relied
upon. There are no committed delivery dates.

## Shipped — 0.1 developer preview

- Strict JSON acceptance contract with full validation.
- Commit resolution and safe `git archive` snapshot extraction.
- Trusted local execution with timeouts, output limits and a minimal environment.
- Go (`go test -json`) and JUnit XML adapters, failing closed.
- Requirement aggregation with five evidence states.
- JSON / Markdown / self-contained HTML receipts.
- Stable CLI exit codes and a `report` subcommand.
- Runnable Go and JUnit examples.
- Contract and receipt JSON Schemas.
- Documentation site.

## P1 — after the core is dependable

- Validate the reusable GitHub Action from a separate sample repository, with
  explicit trust boundaries and reviewed policy inputs.
- Platform release binaries, checksums and installation instructions.
- Schema compatibility tests across contract/receipt versions.
- Better diagnostic codes and a report-only CI mode.
- More framework fixtures based on real emitted JUnit reports.
- Published sample receipts generated from the bundled examples.
- Protected-base policy examples and acceptance-test guidance.

## P2 — later

- A restricted container runner with explicit network and resource controls and
  a dedicated threat model (no network by default, no Docker socket, dropped
  capabilities, resource limits, explicit image policy).
- Signed attestations, identity verification and provenance documentation.
- Additional adapters, only with concrete users and real fixtures.
- Policy bundles, richer receipt comparison and historical views.
- Optional integrations with agent ecosystems, keeping the core
  provider-independent.

## Explicitly out of scope

- Proving complete correctness, security or the absence of bugs.
- Using an LLM to infer whether arbitrary prose requirements are satisfied.
- Replacing human code review or scoring test quality automatically.
- A hosted service, team accounts or billing.
- Any claim of authenticated provenance from an unsigned receipt.
