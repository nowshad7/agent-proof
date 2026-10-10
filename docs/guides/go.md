# Go test JSON

The `go` parser reads a complete `go test -json` stream from the check's
standard output.

## Contract

```json
{
  "id": "unit",
  "command": ["go", "test", "-json", "./..."],
  "parser": "go",
  "timeout_seconds": 120
}
```

The `-json` flag is essential: without it, `go test` prints human-readable text
that the parser will reject.

## Identities

Go identities are `import/path::TestName`. Subtests preserve their path:

```text
example.com/project::TestExport
example.com/project::TestExport/unicode
example.com/project::TestExport/quoting
```

Reference them from a requirement's `test` field exactly as they appear in the
receipt's `tests` array.

## What the parser guarantees

The parser **fails closed**. The following produce an `ERROR`, never a pass:

- A malformed line that is not a JSON object.
- An event without a package.
- A test that started (`run`) but never reported a terminal event
  (`pass`/`fail`/`skip`) — an incomplete or truncated stream.
- A package that emitted tests but never reported its own final result.
- A duplicate terminal event for the same identity, which would make the
  outcome ambiguous (for example a test re-run under `-count` or retries).
- Zero test results.

These rules mean a build that crashes mid-suite, or output captured only
partially, cannot be mistaken for success.

## Notes and limits

- **Coverage is not inferred.** AgentProof does not read coverage percentages
  from ordinary test output or impose thresholds.
- **Caching.** Checks run with a fresh `GOCACHE` and `GOMODCACHE` under a
  scratch `HOME`, so results are not served from your personal build cache.
  The first run therefore compiles from scratch; size `timeout_seconds`
  accordingly.
- **Modules.** The snapshot excludes `.git` and any uncommitted files. If your
  build needs module downloads, the runner inherits `PATH` and common proxy
  variables; see [trusted local execution](/guide/execution).

## Example

See [`examples/export`](https://github.com/nowshad7/agent-proof/tree/main/examples/export)
for a complete, runnable Go project with a contract, and the
[sample receipt](/guide/sample-receipt) it produces.
