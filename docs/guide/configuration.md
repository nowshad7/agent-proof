# Configuration

The contract is a small, reviewable JSON document named `agent-proof.json`. It
maps **requirements** to the **checks** that produce their evidence.

## Shape

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

## Checks

A check is one executable step whose output becomes evidence.

| Field | Required | Meaning |
| --- | --- | --- |
| `id` | yes | Unique, nonempty identifier referenced by requirements. |
| `command` | yes | Argument array. No shell is involved unless you invoke one explicitly. |
| `parser` | yes | `command`, `go` or `junit`. |
| `timeout_seconds` | yes | Between 1 and 3600. |
| `report_path` | junit only | Safe relative path to the JUnit XML the command writes. |

### Parsers

- **`command`** — pass/fail is the process exit code. There are no test
  identities, so a `command` check cannot satisfy a requirement that names a
  `test`.
- **`go`** — reads a complete `go test -json` stream from stdout. See the
  [Go guide](/guides/go).
- **`junit`** — reads a JUnit XML report from `report_path`. See the
  [JUnit guide](/guides/junit).

### Commands are argument arrays

Commands are executed directly as an argv array; there is no implicit shell
interpolation, so `"command": ["go", "test", "-json", "./..."]` runs exactly
those arguments. You *may* invoke a shell explicitly
(`["sh", "-c", "…"]`) — it simply receives the same trust as any other
executable.

## Requirements

A requirement is an acceptance statement.

| Field | Required | Meaning |
| --- | --- | --- |
| `id` | yes | Unique, nonempty identifier. |
| `description` | recommended | Human-readable statement of the behavior. |
| `check` | yes | The `id` of an existing check. |
| `test` | optional | An exact test identity within that check's evidence. |

If `test` is omitted, the requirement is **command-level**: it passes when its
check passes. If `test` is present, the requirement passes only when that exact
test identity passed. A missing named test becomes
[`UNVERIFIED`](/guide/evidence).

## Validation rules

An invalid contract exits with an actionable error **before any command runs**:

- `version` must be `1`; there must be at least one check and one requirement.
- Check and requirement IDs are unique and nonempty.
- Each command is a nonempty argument array with no empty arguments.
- `parser` is one of `command`, `go`, `junit`.
- Each `timeout_seconds` is between 1 and 3600.
- A `junit` check requires a safe, relative `report_path`; non-junit checks
  must not set it.
- Each requirement references an existing check.
- A `command`-parser check cannot back a requirement that names a `test`.
- Unknown fields and trailing JSON content are rejected.

## Where the contract comes from

By default the contract is read from `agent-proof.json` in the **committed
snapshot**. To apply a separately reviewed policy — for example one kept on a
protected base branch so a candidate change cannot weaken its own gate — pass
`--config`:

```sh
agentproof verify --trust-local --config /path/to/reviewed-policy.json
```

AgentProof records the SHA-256 of the exact policy bytes in the receipt, so a
whitespace-only change still changes the recorded hash. See
[keeping policy outside the candidate](/guide/security#protected-policy).
