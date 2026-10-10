<div align="center">

<img src="docs/public/logo.svg" alt="AgentProof" width="88" height="88">

# AgentProof

**Reproducible verification for agent-written changes.**
Explicit requirements. Real test evidence. Portable receipts.

[![CI](https://github.com/nowshad7/agent-proof/actions/workflows/ci.yml/badge.svg)](https://github.com/nowshad7/agent-proof/actions/workflows/ci.yml)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Status: developer preview](https://img.shields.io/badge/status-developer%20preview-f59e0b.svg)](ROADMAP.md)

[Documentation](https://nowshad7.github.io/agent-proof/) ·
[Quickstart](https://nowshad7.github.io/agent-proof/guide/quickstart) ·
[Sample receipt](https://nowshad7.github.io/agent-proof/guide/sample-receipt) ·
[Security model](https://nowshad7.github.io/agent-proof/guide/security)

</div>

> **Developer preview (0.1).** The core workflow — contracts, Go and JUnit
> evidence, and portable receipts — works today. Container isolation and signed
> attestations are [future work](ROADMAP.md) and are **not** implemented. The
> documentation site is a deployment target; verify it is live before relying on
> the link above.

---

AgentProof connects explicit acceptance requirements to executable evidence and
produces portable verification receipts for agent-written changes. Coding agents
can produce convincing explanations without proving the requested behavior works,
and a green test suite can still leave a requirement untested. AgentProof asks you
to write the mapping from **requirement → evidence** down, then executes it
against a committed snapshot and reports what passed, what failed, and **what
remains unverified**.

No model API, account or hosted dashboard is required. Any agent or person that
can run a command-line tool can use it.

## Demonstration

> The output below is a **real** run of the bundled `examples/export` project.
> Regenerate it with `./examples/run-examples.sh export`.

```text
$ agentproof verify --trust-local
AgentProof PASS
  commit 96e440c08f9390dcaff0bce7985d29ee6dde5f58
  [PASS] unicode-export
  [PASS] csv-quoting
Receipts written to .agent-proof/
```

When a declared test is missing, the gate **fails closed** instead of pretending:

```text
AgentProof FAIL
  commit a1b2c3d4e5f6...
  [PASS]       unicode-export
  [UNVERIFIED] error-handling   named test not found in the evidence
```

## Quickstart

```sh
# 1. Install — a prebuilt binary from the releases page, or with Go:
go install github.com/nowshad7/agent-proof/cmd/agentproof@latest

# 2. Add agent-proof.json to your repo root and commit it (see below)

# 3. Verify the committed snapshot
agentproof verify --trust-local

# 4. Read the receipt
agentproof report --input .agent-proof/receipt.json --format md
```

Prebuilt binaries (linux/macOS/windows, amd64/arm64) and checksums are attached
to each [release](https://github.com/nowshad7/agent-proof/releases).

### In CI

The bundled [GitHub Action](https://nowshad7.github.io/agent-proof/guide/ci)
runs the verification and posts the receipt as a single sticky pull-request
comment — updated in place on every run, including on failure — so reviewers see
exactly which requirement passed and which stayed unverified:

```yaml
permissions:
  contents: read
  pull-requests: write
steps:
  - uses: actions/checkout@v4
    with: { fetch-depth: 0 }
  - uses: nowshad7/agent-proof@v0.1.0
    with:
      version: v0.1.0
```

A minimal contract:

```json
{
  "version": 1,
  "checks": [
    { "id": "unit", "command": ["go", "test", "-json", "./..."], "parser": "go", "timeout_seconds": 120 }
  ],
  "requirements": [
    { "id": "unicode-export", "description": "CSV exports preserve Unicode text",
      "check": "unit", "test": "example.com/project::TestUnicode" }
  ]
}
```

AgentProof reads the contract from the **committed** snapshot, resolves the
commit, runs the checks in an exported copy of that commit, and writes
`receipt.json`, `receipt.md` and a self-contained `receipt.html`.

## Supported adapters

| Parser | Evidence | Identity format |
| --- | --- | --- |
| `go` | a complete `go test -json` stream | `import/path::TestName` (subtests preserved) |
| `junit` | a JUnit XML report the command writes | `classname::name` (`::name` with no classname) |
| `command` | process exit code only | — (cannot back a named-test requirement) |

All adapters **fail closed**: missing, incomplete, duplicate or malformed
evidence becomes `UNVERIFIED` or `ERROR`, never a pass.

## Security boundary

⚠️ **Trusted local execution is not a sandbox.** Checks run with your
privileges. For untrusted changes, run inside a disposable, unprivileged
environment, and do not expose production credentials, tokens or Docker sockets.
Receipts are **unsigned and editable** — they prove only that the declared checks
ran, not who ran them, and not that your tests are meaningful. Review the policy,
the tests and the source together. See the
[security model](https://nowshad7.github.io/agent-proof/guide/security).

## Documentation, examples and roadmap

- 📚 [Full documentation](https://nowshad7.github.io/agent-proof/)
- 🧪 Runnable examples: [`examples/export`](examples/export) (Go),
  [`examples/junit`](examples/junit) (Python/JUnit) — run all with
  [`examples/run-examples.sh`](examples/run-examples.sh)
- 🗺️ [Roadmap](ROADMAP.md) · [Changelog](CHANGELOG.md)

## Contributing

Contributions are welcome. Build and test with:

```sh
go build ./... && go vet ./... && go test ./...
./examples/run-examples.sh
```

See [CONTRIBUTING.md](CONTRIBUTING.md), [AGENTS.md](AGENTS.md) for working with
coding agents in this repo, and [SECURITY.md](SECURITY.md) to report a
vulnerability.

## License

[Apache-2.0](LICENSE).
