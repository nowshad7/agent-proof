# Verifying AI-generated code with acceptance tests

Coding agents are good at producing code *and* a confident explanation of why it
works. The explanation is not evidence. This page is about closing that gap with
explicit acceptance tests and a receipt.

## Why a green suite is not enough

"All tests pass" answers a different question than "the behavior I asked for
works." A suite can be green while:

- the requested behavior has no test at all,
- the test exists but asserts almost nothing,
- a test was added and then skipped,
- the relevant test was deleted along with the bug it caught.

What a reviewer actually needs is a mapping from **the requirement** to **the
specific evidence** that proves it — and a signal when that evidence is missing.

## The acceptance-contract approach

1. Write the requirement down in plain language.
2. Point it at the exact test identity that proves it.
3. Verify against a committed snapshot.
4. Read the state: `PASS`, `FAIL`, `UNVERIFIED`, `SKIPPED` or `ERROR`.

```json
{
  "id": "unicode-export",
  "description": "CSV exports preserve Unicode text",
  "check": "unit",
  "test": "example.com/project::TestUnicode"
}
```

If `TestUnicode` does not exist, the requirement is `UNVERIFIED` — the receipt
says, in effect, "you claimed this, but nothing proves it."

## What this does and does not buy you

It buys you a **fail-closed** signal: the absence of evidence is visible, not
assumed-good. It makes a reviewer's job narrower and faster.

It does **not** judge whether your assertions are meaningful. A test that asserts
`true == true` will pass. That is why AgentProof is a complement to human review,
not a replacement — see the [security model](/guide/security).

## Where to go next

- [Detecting missing evidence](/learn/missing-evidence)
- [Commit-specific receipts](/learn/commit-receipts)
- [Agent workflows](/guide/agents)
