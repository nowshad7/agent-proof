# AGENTS.md

Guidance for coding agents (and their humans) working in this repository.

## What this project is

AgentProof maps explicit acceptance requirements to executable evidence from an
exact committed snapshot and produces portable receipts. Its single most
important property is that it **fails closed**: missing, truncated, duplicate or
malformed evidence must never be reported as a pass.

## Hard rules

- **Never let missing evidence pass.** If you touch an adapter
  (`internal/proof/adapter_*.go`) or aggregation (`verify.go`), preserve the
  fail-closed behavior and add tests for the missing/incomplete/duplicate cases.
- **Do not claim unimplemented capabilities.** Container isolation and signed
  attestations are future work. Do not describe local execution as a sandbox or
  a receipt as an attestation, in code, docs or commit messages.
- **Prefer the standard library.** The core avoids third-party dependencies.
- **Keep the security framing honest.** See `docs/guide/security.md` and
  `SECURITY.md`.

## Build, test, verify

```sh
go build ./...
go vet ./...
gofmt -l .                 # must be empty
go test ./...               # or: go test -race ./...
./examples/run-examples.sh  # runs the bundled examples through the CLI
```

This repository verifies itself. `agent-proof.json` points at the core
regression tests:

```sh
go build -o /tmp/agentproof ./cmd/agentproof
/tmp/agentproof verify --trust-local
```

A change is not done until `go test ./...` passes, the examples pass, and
self-verification reports `PASS`.

## Layout

| Path | Responsibility |
| --- | --- |
| `cmd/agentproof` | CLI arguments, exit codes, output locations |
| `internal/proof` | Contract types, validation, snapshot, execution, adapters, aggregation, reports |
| `examples/` | Small, deterministic, runnable examples |
| `schemas/` | Contract and receipt JSON Schemas |
| `docs/` | VitePress documentation |
| `.github/workflows/` | CI (verification) and Pages (deployment), kept separate |

## When adding evidence states or parsers

- Test identities: Go is `import/path::TestName` (subtests preserved); JUnit is
  `classname::name` (`::name` with no classname).
- A `command`-parser check cannot back a named-test requirement — the validator
  enforces this.
- An operational error (`ERROR`, exit `2`) takes precedence over an ordinary
  acceptance failure (`FAIL`/`UNVERIFIED`/`SKIPPED`, exit `1`).

## Docs

```sh
npm ci && npm run docs:build   # dead links fail the build
```

Keep docs in step with behavior. Do not add badges or links for things that do
not exist yet.
