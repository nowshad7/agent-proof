# Roadmap

This reflects the delivery plan. It is dependency-ordered, not date-driven, and
it distinguishes what is **shipped** from what is **planned**. Planned items are
not implemented and should not be relied upon.

## Shipped in the 0.1 preview

- Strict JSON acceptance contract with full validation.
- Commit resolution and safe `git archive` snapshot extraction.
- Trusted local command execution with timeouts, output limits and a minimal
  environment.
- Go (`go test -json`) and JUnit XML evidence adapters, failing closed on
  missing, incomplete, duplicate or malformed evidence.
- Requirement aggregation with five evidence states.
- JSON, Markdown and self-contained HTML receipts.
- Stable CLI exit codes.
- Runnable Go and JUnit examples.
- This documentation.

## Planned — P1 (after the core is dependable)

- A reusable GitHub Action with explicit trust boundaries and reviewed policy
  inputs, validated from a separate sample repository.
- Platform release binaries, checksums and installation instructions.
- Machine-readable contract and receipt schemas with compatibility tests.
- Better diagnostic codes and a report-only CI mode.
- More framework fixtures based on real emitted JUnit reports.
- Published sample receipts generated from the bundled examples.
- Protected-base policy examples and acceptance-test guidance.

## Planned — P2 (later)

- A restricted container runner with explicit network and resource controls and
  a dedicated threat model.
- Signed attestations, identity verification and provenance documentation.
- Additional adapters, only when there are concrete users and fixtures.
- Policy bundles, richer receipt comparison and historical views.
- Optional integrations with agent ecosystems, without making the core
  provider-specific.

::: warning
Container isolation and signed attestations are **future work**. The current
preview's local execution is not a sandbox, and its receipts are unsigned. See
the [security model](/guide/security).
:::
