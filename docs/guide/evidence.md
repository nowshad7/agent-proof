# Evidence states

Every requirement resolves to one of five states. Missing results are never
silently treated as success.

| State | Meaning | Passes the gate? |
| --- | --- | --- |
| <span class="ap-badge ap-pass">PASS</span> | Declared command or test supplied passing evidence. | Yes |
| <span class="ap-badge ap-fail">FAIL</span> | Command returned nonzero, or a test failed. | No |
| <span class="ap-badge ap-warn">UNVERIFIED</span> | No usable test evidence, or a named test was missing. | No |
| <span class="ap-badge ap-fail">ERROR</span> | Execution, timeout, parsing or report handling failed. | No |
| <span class="ap-badge ap-warn">SKIPPED</span> | A required test was reported as skipped. | No |

## How the gate aggregates

- **All configured checks participate.** A passing named test cannot rescue a
  sibling check that failed.
- **An error takes precedence** over an ordinary acceptance failure, because an
  operational problem means the evidence is untrustworthy, not merely negative.
- **An unrelated skipped test does not fail a check** unless a requirement names
  it.

The overall status is computed as:

1. If any check or requirement is `ERROR` → overall `ERROR` (exit code `2`).
2. Otherwise, if any check is not `PASS`, or any requirement is not `PASS`
   (that is, `FAIL`, `UNVERIFIED` or `SKIPPED`) → overall `FAIL` (exit code `1`).
3. Otherwise → overall `PASS` (exit code `0`).

## Worked examples

**A named test that passed, but its parent check failed.** The check runs the
whole suite; another test in it failed, so the command exited nonzero. The
check is `FAIL`. Even though your named test passed, the requirement is reported
as `FAIL` ("the named test passed but its parent check failed") and the gate
does not pass. This prevents a green requirement from hiding a broken suite.

**A requirement names a test that was never emitted.** The evidence parsed
cleanly but contains no such identity. The requirement is `UNVERIFIED`. This is
the common agent failure mode: the change claims a behavior that no test
actually exercises.

**A command timed out.** The check is `ERROR`, every requirement pointing at it
is `ERROR`, and the overall status is `ERROR`. Automation sees exit code `2` and
can distinguish this from a genuine acceptance failure.

See [requirement mapping](/guide/requirements) for the exact lookup rules.
