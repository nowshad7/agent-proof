# Requirement mapping

A requirement connects a human-readable acceptance statement to the exact
evidence that supports it.

## Command-level requirements

Omit `test` to tie a requirement directly to its check's pass/fail:

```json
{
  "id": "builds",
  "description": "The project builds",
  "check": "build"
}
```

- Check `PASS` → requirement `PASS`.
- Check `FAIL` → requirement `FAIL`.
- Check `ERROR` → requirement `ERROR`.

Use these for checks where the command's exit code is the whole story: a build,
a linter, a formatter, a schema validation.

## Named-test requirements

Add `test` to require that one exact test identity passed:

```json
{
  "id": "unicode-export",
  "description": "CSV exports preserve Unicode text",
  "check": "unit",
  "test": "example.com/project::TestUnicode"
}
```

The lookup against the check's parsed evidence resolves as follows:

| Situation | Resulting state |
| --- | --- |
| Identity found and passed, parent check not failed | `PASS` |
| Identity found and passed, but parent check failed | `FAIL` |
| Identity found and failed | `FAIL` |
| Identity found and skipped | `SKIPPED` |
| Identity not found | `UNVERIFIED` |
| Parent check errored | `ERROR` |

Only `go` and `junit` checks have test identities. A `command`-parser check
cannot back a named-test requirement, and the contract validator rejects that
at load time.

## Test identities

The identity format depends on the parser:

- **Go:** `import/path::TestName`, with subtests preserved as
  `import/path::TestName/subtest`. Example:
  `example.com/project::TestExport/unicode`.
- **JUnit:** `classname::name`. When a case has no classname, the identity is
  `::name`. Example: `test_slugify.SlugifyTests::test_unicode`.

The quickest way to learn the exact identity is to run the check once and read
the `tests` array in `receipt.json`, or the per-test table in the HTML receipt.

## Choosing good requirements

- Name the **behavior**, not the test file. `description` is what a reviewer
  reads.
- Point each requirement at the **narrowest** test that proves it, so an
  unrelated change to the suite does not quietly satisfy it.
- Keep acceptance tests where a candidate change cannot edit them away — see
  [keeping policy and tests protected](/guide/security#protected-policy).
