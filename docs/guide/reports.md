# Reports

Every verification writes its receipts to the output directory (default
`.agent-proof/`). They all describe the same result.

| File | Purpose |
| --- | --- |
| `receipt.json` | Full structured result for automation and archival. |
| `receipt.md` | Compact requirement summary for code review. |
| `receipt.html` | Self-contained, portable page for sharing. |
| `badge.json` | Shields.io endpoint for a live README [badge](/guide/badge). |

## Regenerating a view

You can re-render any format from a stored `receipt.json` without re-running the
checks:

```sh
agentproof report --input .agent-proof/receipt.json --format md
agentproof report --input .agent-proof/receipt.json --format html > receipt.html
agentproof report --input .agent-proof/receipt.json --format json
```

## JSON

The canonical result. It contains the receipt and AgentProof versions, the
resolved commit, the config hash, the start time and runner, the overall status,
every check (command, exit code, duration, per-test outcomes, captured output)
and every requirement (with its resulting state and detail). See the
[receipt format reference](/reference/receipt).

## Markdown

A short table of requirements and checks, plus the trust warnings. This is the
form to paste into a pull request or review thread. See a real one on the
[sample receipt](/guide/sample-receipt) page.

## HTML

A single self-contained file — no external assets — with expandable check
details and per-test tables. It is safe to share:

- All evidence (test names, captured output, diagnostics) is HTML-escaped.
- A restrictive `Content-Security-Policy` is applied: no scripts, no external
  resources.

## Badge

A [Shields.io endpoint](https://shields.io/badges/endpoint-badge) JSON document
for a live README badge. `verify` writes it to `badge.json` automatically, and
`report --format badge` re-renders it from a stored receipt. It fails closed —
only an explicit `PASS` is green. See [Verification badge](/guide/badge).

## Interpreting results

- The **overall status** is the gate. Only `PASS` is a pass.
- Each **requirement** shows its [state](/guide/evidence) and a plain-language
  detail when it is not a pass.
- Each **check** shows its parser, status, exit code and duration; expand it in
  HTML to see the tests it produced and any captured `stderr`.

## On receipt integrity

Receipts are **unsigned and editable**. They prove only that the declared
checks were run and what they reported. A plain content hash, if added, is for
integrity, not authentication — AgentProof never describes a hash as a signature
or an independent attestation. See [provenance limits](/guide/security#provenance).
