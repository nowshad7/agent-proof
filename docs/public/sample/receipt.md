# AgentProof receipt

- **Status:** PASS
- **Commit:** `96e440c08f9390dcaff0bce7985d29ee6dde5f58`
- **Config:** `sha256:1e5302987f79c47e1da5c1d68ed58d05f6ac56db30dc8f3468c99e853fc2a7c4`
- **Started:** 2026-10-10T04:41:31Z
- **Runner:** examples/export
- **AgentProof:** 0.1.0-dev

## Requirements

| Requirement | State | Check | Detail |
| --- | --- | --- | --- |
| unicode-export | PASS | unit | — |
| csv-quoting | PASS | unit | — |

## Checks

| Check | Parser | Status | Exit | Duration (ms) |
| --- | --- | --- | --- | --- |
| unit | go | PASS | 0 | 8547 |

## Trust and scope

- Trusted local execution is not a sandbox: checks run with the caller's privileges.
- This receipt is unsigned and editable; it is not an authenticated attestation.
- AgentProof verifies the declared evidence workflow, not the honesty or completeness of the tests.
