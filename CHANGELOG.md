# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project aims
to follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html) once it
reaches a stable release.

## [0.1.0] — 2026-10-10 (developer preview)

The first public developer preview.

### Added

- `agentproof report --format badge` renders a Shields.io endpoint JSON for a
  live README verification badge, and `verify` writes it to `badge.json` on
  every run. The badge fails closed: only an explicit `PASS` is coloured green,
  and an unknown or empty status reads `unknown`, never a pass. A new
  `badge-fails-closed` self-verification requirement and the
  [badge guide](docs/guide/badge.md) cover it; the Action exposes it as a
  `badge` output.
- `agentproof report --format gh` renders a GitHub-flavored Markdown body for a
  pull-request comment, with a sticky `<!-- agentproof-receipt -->` marker so it
  can be updated in place.
- The GitHub Action posts the receipt as a single sticky pull-request comment
  (created once, updated on every run, including on failure), gated by a
  `comment-on-pr` input and the caller's `pull-requests: write` permission.
- The GitHub Action downloads a prebuilt binary when `version` is a release tag,
  falling back to building from source for a branch or SHA.
- GoReleaser configuration and a tag-triggered release workflow produce
  cross-platform binaries (linux/macOS/windows, amd64/arm64) and a
  `checksums.txt`, with the build version injected via ldflags.
- CI validates the GoReleaser config (`goreleaser check`) and exercises the
  Action end to end against this repository on pull requests.
- Strict `agent-proof.json` acceptance contract (version 1) with full validation
  of checks, requirements, parsers, timeouts and references; unknown fields and
  trailing content are rejected.
- Commit resolution and safe `git archive` snapshot extraction, including correct
  handling of git's `pax_global_header` metadata and rejection of symlinks, hard
  links, path traversal and oversized archives (256 MiB / 50,000 entries).
- Trusted local command execution (behind `--trust-local`) with per-check
  timeouts, 1 MiB per-stream output bounds, process-group cancellation on Unix
  and a minimal, non-inherited environment.
- Go `go test -json` adapter and JUnit XML adapter, both failing closed on
  missing, incomplete, duplicate or malformed evidence.
- Requirement aggregation across five states (`PASS`, `FAIL`, `UNVERIFIED`,
  `ERROR`, `SKIPPED`) with a fail-closed overall gate.
- JSON, Markdown and self-contained, HTML-escaped receipts with a restrictive
  content-security policy.
- Stable CLI exit codes (`0` pass, `1` acceptance failure, `2` operational
  error) and a `report` subcommand to re-render stored receipts.
- Runnable Go and JUnit examples and an example runner script.
- Machine-readable contract and receipt JSON Schemas.
- A composite GitHub Action (`action.yml`) — preview quality.
- VitePress documentation with a graphite-and-teal theme.

### Known limitations

- Local execution is not a sandbox; container isolation is future work.
- Receipts are unsigned; signed attestations are future work.
- Submodules and Git LFS hydration are not supported.
- The bundled Action is early; validate it from a separate repository first.

[0.1.0]: https://github.com/nowshad7/agent-proof/releases/tag/v0.1.0
