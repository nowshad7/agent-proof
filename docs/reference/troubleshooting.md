# Troubleshooting

## `verify requires --trust-local`

`verify` will not run checks without `--trust-local`. This is intentional — it is
your acknowledgement that checks run locally with your privileges. Add the flag
once you understand the [security model](/guide/security).

## `reading agent-proof.json from the committed snapshot`

The contract was not found in the committed snapshot. Remember that AgentProof
reads **committed** content: if you just created `agent-proof.json`, commit it
first, or pass an explicit `--config` path.

## A requirement is `UNVERIFIED`

The named test identity was not found in the parsed evidence. Common causes:

- The identity is misspelled. Run the check once and copy the exact name from
  the `tests` array in `receipt.json`.
- For Go, the identity must be `import/path::TestName` — the full import path,
  not the directory.
- For JUnit, the identity is `classname::name`; a case with no classname is
  `::name`.
- The test genuinely does not exist yet.

## A check is `ERROR`

`ERROR` means an operational failure, not an acceptance failure. The `diagnostic`
field explains it:

- **"command timed out"** — raise `timeout_seconds`. Note the first Go run
  compiles from scratch under a fresh cache.
- **"command output exceeded 1048576 bytes"** — the check printed more than
  1 MiB to one stream. Reduce verbosity, or write large output to a file.
- **"incomplete go test json"** — the suite crashed or output was truncated.
  Run the command directly to see why.
- **"reading junit report … no such file"** — the command did not write
  `report_path`. AgentProof removes any stale committed report before the run,
  so the command must regenerate it.
- **"is a symlink, which is rejected"** — `report_path` resolved to a symlink.

## `go test` is not using `-json`

The `go` parser needs machine-readable output. Your command must include
`-json`, for example `["go", "test", "-json", "./..."]`.

## The gate fails even though a named test passed

All checks participate in the gate. If the check's command exited nonzero
because *another* test failed, the check is `FAIL`, and a passing named test
cannot rescue it. Fix the failing test, or narrow the check's scope.

## Module downloads fail during a check

Checks run with a minimal environment. `PATH` and common proxy variables
(`HTTPS_PROXY`, `SSL_CERT_FILE`, …) are forwarded, but private module
credentials in other variables are not. Vendor dependencies, or run in an
environment where the forwarded proxy settings suffice.

## Still stuck?

Open an issue with your tool version (`agentproof version`), OS, a minimal
contract and a sanitized receipt. See [contributing](/community/contributing).
