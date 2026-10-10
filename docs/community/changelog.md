# Changelog

The authoritative changelog lives at
[`CHANGELOG.md`](https://github.com/nowshad7/agent-proof/blob/main/CHANGELOG.md)
in the repository.

## 0.1.0 — Developer preview

The first public developer preview.

### Added

- Strict `agent-proof.json` acceptance contract (version 1) with full
  validation of checks, requirements, parsers, timeouts and references.
- Commit resolution and safe `git archive` snapshot extraction, including
  correct handling of git's `pax_global_header` metadata and rejection of
  symlinks, hard links, traversal and oversized archives.
- Trusted local command execution (behind `--trust-local`) with per-check
  timeouts, 1 MiB output bounds, process-group cancellation on Unix and a
  minimal, non-inherited environment.
- Go `go test -json` adapter and JUnit XML adapter, both failing closed on
  missing, incomplete, duplicate or malformed evidence.
- Requirement aggregation across five states (`PASS`, `FAIL`, `UNVERIFIED`,
  `ERROR`, `SKIPPED`) with a fail-closed overall gate.
- JSON, Markdown and self-contained, escaped HTML receipts.
- Stable CLI exit codes (`0` pass, `1` acceptance failure, `2` operational
  error) and a `report` subcommand to re-render stored receipts.
- Runnable Go and JUnit examples and an example runner script.
- VitePress documentation.

### Known limitations

- Local execution is not a sandbox. Container isolation is future work.
- Receipts are unsigned; signed attestations are future work.
- Submodules and Git LFS hydration are not supported.
- The bundled GitHub Action is early and should be validated from a separate
  repository before use.

This is a preview: interfaces may change before a stable release.
