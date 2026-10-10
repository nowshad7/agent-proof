# Receipt format

A receipt is the portable record of one verification. The JSON form is
canonical; the Markdown and HTML forms are renderings of it. A machine-readable
JSON Schema ships at
[`schemas/receipt.schema.json`](https://github.com/nowshad7/agent-proof/blob/main/schemas/receipt.schema.json).

## Top level

| Field | Type | Meaning |
| --- | --- | --- |
| `receipt_version` | integer | Receipt document format version (currently `1`). |
| `agentproof_version` | string | The AgentProof build that produced it. |
| `commit` | string | Resolved full commit SHA that was executed. |
| `config_hash` | string | `sha256:` of the exact policy bytes. |
| `started_at` | string | RFC 3339 UTC start timestamp. |
| `runner` | string | Optional runner identity from `--runner`. |
| `status` | string | Overall gate: `PASS`, `FAIL`, `UNVERIFIED`, `ERROR` or `SKIPPED`. |
| `checks` | array | Per-check results. |
| `requirements` | array | Per-requirement results. |
| `warnings` | string[] | Trust and scope warnings. |

## `checks[]`

| Field | Type | Meaning |
| --- | --- | --- |
| `id` | string | The check id. |
| `command` | string[] | The argv that ran. |
| `parser` | string | `command`, `go` or `junit`. |
| `exit_code` | integer | Process exit code (`-1` if it never produced one). |
| `duration_ms` | integer | Wall-clock duration in milliseconds. |
| `status` | string | Check status. |
| `tests` | array | `{name, status}` per test identity (go/junit only). |
| `stdout` / `stderr` | string | Captured output (bounded to 1 MiB each). |
| `diagnostic` | string | Present on `ERROR`, explaining the operational failure. |

## `requirements[]`

| Field | Type | Meaning |
| --- | --- | --- |
| `id` | string | The requirement id. |
| `description` | string | The acceptance statement. |
| `check` | string | The check it references. |
| `test` | string | The named test identity, if any. |
| `state` | string | Resulting [state](/guide/evidence). |
| `detail` | string | Plain-language reason when not a pass. |

## Provenance

A receipt is **unsigned and editable**. `config_hash` identifies the policy
bytes and `commit` identifies the code, but neither authenticates *who* produced
the receipt. Do not treat a receipt as an attestation. See
[provenance limits](/guide/security#provenance).

See the [sample receipt](/guide/sample-receipt) for a real, complete document.
