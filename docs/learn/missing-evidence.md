# Detecting missing evidence in agent-written changes

The most useful thing a verification tool can do with agent-written changes is
notice when a claim has **no evidence behind it**. Many tools treat the absence
of a result as success. AgentProof treats it as `UNVERIFIED` and fails the gate.

## The failure modes that "pass" elsewhere

| What happened | Naive result | AgentProof |
| --- | --- | --- |
| The requirement's test was never written | green suite | `UNVERIFIED` |
| The test exists but was skipped | green suite | `SKIPPED` |
| The suite crashed after some tests | partial / green | `ERROR` |
| Output was truncated | parsed as pass | `ERROR` |
| A stale committed report is present | "passes" | removed before run → `ERROR` if not regenerated |
| A test identity was duplicated ambiguously | last-wins | `ERROR` |

Each row is a way a change can *look* verified without *being* verified.

## How AgentProof fails closed

- **Named test not found** → `UNVERIFIED`. The identity you required is not in
  the evidence.
- **Incomplete `go test -json`** → `ERROR`. Every started test must have a
  terminal event, and every package must finish.
- **Missing JUnit report** → `ERROR`. The stale report is removed first, so the
  command must produce a fresh one.
- **Zero results** → `UNVERIFIED`/`ERROR`. An empty suite proves nothing.

## Putting it to work

Write one requirement per behavior you care about, each pointing at the narrowest
test that proves it. Then the receipt becomes a checklist: every `PASS` has named
evidence, and every non-pass tells you exactly what is missing.

```text
[PASS]       unicode-export
[UNVERIFIED] error-handling   named test "pkg::TestErrors" not found in the evidence
[SKIPPED]    migration        the named test was reported as skipped
```

See [evidence states](/guide/evidence) for the precise aggregation rules.
