# Introduction

AgentProof connects explicit acceptance requirements to executable evidence and
produces portable verification receipts for agent-written changes.

## Why this exists

When a coding agent opens a pull request, the explanation is often more
polished than the proof. A passing test suite tells you *something* ran, but it
does not tell you which acceptance requirement each test was meant to satisfy,
or whether the requirement you care about was tested at all.

AgentProof asks you to write that mapping down — a short contract from
requirements to the exact commands and test identities that support them — and
then executes it against a precise, committed snapshot of the code. The result
is a receipt that states, for every requirement:

- whether declared evidence **passed**,
- whether it **failed**,
- or whether it is **unverified** because the evidence was missing.

Missing evidence is never silently treated as success.

## The workflow in three steps

1. **Declare** a small acceptance contract (`agent-proof.json`). Each
   requirement points at a check, and optionally at an exact test identity.
2. **Execute** the checks against the committed snapshot of a resolved commit.
   Uncommitted and untracked files are excluded.
3. **Inspect** the receipt — JSON for automation, Markdown for review, and a
   self-contained HTML page for sharing.

## What AgentProof is

- A small, portable Go CLI with no required model API, account or hosted
  service.
- Agent-independent: any agent or person that can run a command can use it.
- Honest about its boundaries: trusted local execution is explicitly *not* a
  sandbox, and unsigned receipts are *not* authenticated attestations.

## What AgentProof is not

- It does not prove correctness, security or the absence of bugs.
- It does not use an LLM to infer whether prose requirements are satisfied.
- It does not replace human code review or judge test quality for you.
- It does not, in this preview, contain hostile code — local execution runs
  with your privileges. For untrusted changes, use a disposable, unprivileged
  environment. See the [security model](/guide/security).

## Next steps

- [Quickstart](/guide/quickstart) — run your first verification in a few
  minutes.
- [Sample receipt](/guide/sample-receipt) — see what a receipt looks like.
- [Configuration](/guide/configuration) — the contract in detail.
