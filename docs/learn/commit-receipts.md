# Generating a commit-specific test receipt

A receipt is only meaningful if you know *exactly what it describes*. AgentProof
ties every receipt to a single resolved commit and a single policy hash.

## Why the commit matters

If verification runs against your working tree, an uncommitted edit — or a file
the agent forgot to commit — can change the result without changing the record.
AgentProof avoids this by snapshotting a **committed** commit:

1. The revision you pass (default `HEAD`) is resolved to a full 40-character SHA.
2. That commit's tree is exported with `git archive` into a temporary directory.
3. Checks run there. Dirty and untracked files are excluded; `.git` is absent.

The receipt records the full SHA, so the receipt means: *these checks, against
exactly this committed code.*

## Why the policy hash matters

The receipt also records `config_hash`, the SHA-256 of the exact policy bytes
applied. Two consequences:

- A whitespace-only change to the contract produces a different hash, so you can
  never be fooled about which policy ran.
- When you apply a reviewed policy from a protected branch with `--config`, the
  hash lets a reviewer confirm the protected bytes were used — not a version the
  candidate change could edit.

## What a commit receipt does *not* promise

An exact input commit does not make execution bit-for-bit reproducible. Runtime
tool versions, network services and unpinned dependencies can still vary between
runs. The receipt identifies the **inputs** precisely; it does not freeze the
world around them. See [execution limits](/guide/execution).

## Using it in review

Attach `receipt.md` to the pull request, or link the `receipt.html`. A reviewer
can read, at a glance:

- which commit was verified,
- which policy (by hash) was applied,
- which requirements passed and which did not, and why.

That is a far stronger artifact than "the agent says it works."
