# CLI reference

```sh
agentproof version
agentproof verify [flags]
agentproof report --input <receipt.json> [--format json|md|html]
```

## `version`

Prints the AgentProof version and exits `0`. (`--version` and `-v` are
accepted too.)

## `verify`

Resolves a commit, exports its snapshot, runs the contract's checks, parses
evidence, aggregates requirements and writes receipts.

| Flag | Default | Meaning |
| --- | --- | --- |
| `--trust-local` | — | **Required.** Acknowledges that checks run locally with your privileges. Not a sandbox. |
| `--repo <dir>` | `.` | Repository to verify. |
| `--commit <rev>` | `HEAD` | Revision to resolve and snapshot. Any git revision works. |
| `--config <path>` | — | Use a separately reviewed policy instead of the committed `agent-proof.json`. |
| `--out <dir>` | `.agent-proof` | Directory for the receipt files. |
| `--runner <name>` | — | Optional runner identity recorded in the receipt. |

Without `--trust-local`, `verify` exits `2` and does nothing. This is a
deliberate, explicit gate on local execution.

## `report`

Re-renders a stored receipt in another format. It does not run checks.

| Flag | Default | Meaning |
| --- | --- | --- |
| `--input <path>` | — | **Required.** Path to a `receipt.json`. |
| `--format <fmt>` | `md` | `json`, `md` (or `markdown`), or `html`. |

Output is written to standard output.

## Exit codes

| Code | Meaning |
| --- | --- |
| `0` | All required acceptance passed. |
| `1` | Acceptance failed, or required evidence was unverified/skipped. |
| `2` | Configuration, execution, parsing or report error. |

These are stable across releases so automation can distinguish a genuine
acceptance failure (`1`) from an operational error (`2`).

## Examples

```sh
# Verify HEAD of the current repository.
agentproof verify --trust-local

# Verify a specific commit, writing receipts elsewhere.
agentproof verify --trust-local --commit 1a2b3c4 --out build/receipts

# Apply a reviewed policy kept on a protected branch.
agentproof verify --trust-local --config /tmp/policy.json

# Render a shareable HTML page from a stored receipt.
agentproof report --input .agent-proof/receipt.json --format html > receipt.html
```
