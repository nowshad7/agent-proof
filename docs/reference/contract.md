# Contract schema

The contract (`agent-proof.json`) is strict JSON. Unknown fields and trailing
content are rejected. A machine-readable JSON Schema ships in the repository at
[`schemas/contract.schema.json`](https://github.com/nowshad7/agent-proof/blob/main/schemas/contract.schema.json).

## Top level

| Field | Type | Rules |
| --- | --- | --- |
| `version` | integer | Must be `1`. |
| `checks` | array | At least one check. |
| `requirements` | array | At least one requirement. |

## `checks[]`

| Field | Type | Rules |
| --- | --- | --- |
| `id` | string | Nonempty, unique across checks. |
| `command` | string[] | Nonempty; no empty arguments. Executed as argv, no shell. |
| `parser` | string | `command`, `go` or `junit`. |
| `timeout_seconds` | integer | 1–3600 inclusive. |
| `report_path` | string | Required for `junit`; forbidden otherwise. Safe, relative path. |

## `requirements[]`

| Field | Type | Rules |
| --- | --- | --- |
| `id` | string | Nonempty, unique across requirements. |
| `description` | string | Optional but recommended. |
| `check` | string | Must reference an existing check `id`. |
| `test` | string | Optional. If set, the check's parser must be `go` or `junit`. |

## Full example

```json
{
  "version": 1,
  "checks": [
    {
      "id": "unit",
      "command": ["go", "test", "-json", "./..."],
      "parser": "go",
      "timeout_seconds": 120
    },
    {
      "id": "e2e",
      "command": ["python3", "junit_runner.py"],
      "parser": "junit",
      "timeout_seconds": 300,
      "report_path": "build/report.xml"
    },
    {
      "id": "lint",
      "command": ["gofmt", "-l", "."],
      "parser": "command",
      "timeout_seconds": 60
    }
  ],
  "requirements": [
    {
      "id": "unicode-export",
      "description": "CSV exports preserve Unicode text",
      "check": "unit",
      "test": "example.com/project::TestUnicode"
    },
    {
      "id": "signup-flow",
      "description": "End-to-end signup succeeds",
      "check": "e2e",
      "test": "e2e.Signup::test_happy_path"
    },
    {
      "id": "formatted",
      "description": "Source is gofmt-clean",
      "check": "lint"
    }
  ]
}
```

## Config hashing

AgentProof records the SHA-256 of the exact policy bytes applied, prefixed with
`sha256:`. Any change to the bytes — including whitespace — changes the recorded
hash, so the receipt always identifies the precise policy that ran.
