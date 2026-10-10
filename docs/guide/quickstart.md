# Quickstart

This runs your first verification in a few minutes against a sample project.

## Prerequisites

- **Go 1.23 or newer** to build the CLI (`go version`).
- **Git**, because AgentProof verifies a committed snapshot.
- Whatever your checks need — for the Go example, just Go; for the JUnit
  example, `python3`.

## 1. Install the CLI

**Prebuilt binary** (no Go toolchain needed) — download the archive for your
platform from the [latest release](https://github.com/nowshad7/agent-proof/releases),
verify it against `checksums.txt`, extract, and put `agentproof` on your `PATH`:

```sh
# Example for Linux x86-64; adjust the version and platform.
VER=0.1.0
curl -fsSLO "https://github.com/nowshad7/agent-proof/releases/download/v${VER}/agentproof_${VER}_linux_amd64.tar.gz"
tar -xzf "agentproof_${VER}_linux_amd64.tar.gz"
./agentproof version
```

**With Go**, into your `GOBIN`:

```sh
go install github.com/nowshad7/agent-proof/cmd/agentproof@latest
agentproof version
```

Or build a local binary from a clone:

```sh
git clone https://github.com/nowshad7/agent-proof
cd agent-proof
go build -o agentproof ./cmd/agentproof
```

## 2. Add a contract

In the root of a git repository, create `agent-proof.json`. Here is a minimal
Go contract:

```json
{
  "version": 1,
  "checks": [
    {
      "id": "unit",
      "command": ["go", "test", "-json", "./..."],
      "parser": "go",
      "timeout_seconds": 120
    }
  ],
  "requirements": [
    {
      "id": "unicode-export",
      "description": "CSV exports preserve Unicode text",
      "check": "unit",
      "test": "example.com/project::TestUnicode"
    }
  ]
}
```

Commit it. AgentProof reads the contract from the **committed** snapshot, so
uncommitted edits to the policy or the code are not verified.

```sh
git add agent-proof.json
git commit -m "Add acceptance contract"
```

## 3. Verify

```sh
agentproof verify --trust-local
```

`--trust-local` is required: it is your explicit acknowledgement that checks run
locally with your privileges. This is **not** a sandbox — see the
[security model](/guide/security).

You will see a summary and three receipt files in `.agent-proof/`:

```text
AgentProof PASS
  commit 96e440c08f9390dcaff0bce7985d29ee6dde5f58
  [PASS] unicode-export
Receipts written to .agent-proof/
```

## 4. Read the receipt

```sh
agentproof report --input .agent-proof/receipt.json --format md
open .agent-proof/receipt.html   # or your platform's equivalent
```

See a complete example on the [sample receipt](/guide/sample-receipt) page.

## Try the bundled examples

The repository ships runnable Go and JUnit examples. The helper script copies
each into a throwaway git repository and verifies it:

```sh
./examples/run-examples.sh          # both examples
./examples/run-examples.sh export   # just the Go one
```

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | All required acceptance passed. |
| `1` | Acceptance failed, or required evidence was unverified/skipped. |
| `2` | Configuration, execution, parsing or report error. |

These are stable so automation can distinguish a real failure from an
operational error. Next, learn how to write the [contract](/guide/configuration).
