# JUnit XML

The `junit` parser reads a JUnit XML report that your check writes to disk. This
makes AgentProof usable with any framework that can emit JUnit XML — a very
broad, language-independent interchange format.

## Contract

```json
{
  "id": "unit",
  "command": ["python3", "junit_runner.py"],
  "parser": "junit",
  "timeout_seconds": 120,
  "report_path": "report.xml"
}
```

`report_path` is a **safe, relative** path. Absolute paths and parent-directory
traversal are rejected by the contract validator.

## Identities

JUnit identities are `classname::name`. A case with no `classname` yields
`::name`:

```text
test_slugify.SlugifyTests::test_unicode
com.example.ExportTest::preservesUnicode
::standalone_case
```

Both a top-level `<testsuites>` wrapper and a bare `<testsuite>` are supported,
and nested `<testsuite>` elements are traversed. A case counts as failed if it
contains a `<failure>` or `<error>`, skipped if it contains `<skipped>`, and
passed otherwise.

## Stale report protection

A tracked, committed `report.xml` could otherwise "pass" without the tests ever
running. AgentProof prevents this:

1. **Before** the check runs, AgentProof removes the file at `report_path` from
   the snapshot. A stale, committed report cannot survive into the run.
2. **After** the check runs, AgentProof reads the freshly written report. If the
   command did not produce it, the check is `ERROR`.
3. The produced path is validated and **symlinks are rejected**, so a report
   cannot redirect reads outside the snapshot.

The report input is also bounded to 8 MiB.

## Emitters are supported only once fixtured

AgentProof treats a given framework's JUnit output as officially supported only
after its *real* emitted XML is covered by fixtures in the test suite. The
bundled Python example ships a tiny, dependency-free runner
([`junit_runner.py`](https://github.com/nowshad7/agent-proof/blob/main/examples/junit/junit_runner.py))
that produces genuine JUnit XML from `unittest`. If you use pytest, Jest,
Maven Surefire or similar, point `command` at your usual report-producing
invocation and `report_path` at the file it writes.

## Example

See [`examples/junit`](https://github.com/nowshad7/agent-proof/tree/main/examples/junit)
for a complete, runnable project.
