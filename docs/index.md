---
layout: home

hero:
  name: AgentProof
  text: Reproducible verification for agent-written changes
  tagline: Explicit requirements. Real test evidence. Portable receipts. No model API, account or hosted dashboard required.
  image:
    src: /logo.svg
    alt: AgentProof
  actions:
    - theme: brand
      text: Run your first verification
      link: /guide/quickstart
    - theme: alt
      text: Explore a receipt
      link: /guide/sample-receipt
    - theme: alt
      text: GitHub
      link: https://github.com/nowshad7/agent-proof

features:
  - icon: 📝
    title: Declare
    details: Write a small, reviewable acceptance contract that points each requirement at a command or an exact test identity.
  - icon: ▶️
    title: Execute
    details: AgentProof exports the committed snapshot of a resolved commit and runs your checks against it, excluding uncommitted changes.
  - icon: 🔎
    title: Inspect
    details: Get a portable receipt in JSON, Markdown and self-contained HTML showing which requirement passed and which remains unverified.
  - icon: 🧩
    title: Agent-independent
    details: Any agent or human that can invoke a command-line tool can produce a receipt. There is no provider lock-in.
  - icon: 🚫
    title: Fails closed
    details: Missing, failed and skipped evidence is visible. A missing named test becomes UNVERIFIED — it never silently passes.
  - icon: 🔐
    title: Honest boundaries
    details: Trusted local execution is not a sandbox, and unsigned receipts are not attestations. The limits are documented, not hidden.
---

## The problem

Coding agents can produce convincing explanations without proving that the
requested behavior works. A green test suite can still leave an acceptance
requirement untested. Reviewers need to know **which requirement was checked,
which command or named test supports it, which revision was executed, and what
remains unverified.**

AgentProof connects explicit acceptance requirements to executable evidence and
produces a receipt you can review alongside the code and test changes.

```sh
agentproof verify --trust-local
```

```text
AgentProof FAIL
  commit a1b2c3d4e5f6...
  [PASS]       unicode-export
  [UNVERIFIED] error-handling   named test not found in the evidence
Receipts written to .agent-proof/
```

The run above fails closed: one requirement passed, but a declared test was
never found, so the gate does not report success.

::: tip Developer preview
AgentProof is a 0.1 developer preview. The core workflow — contracts, Go and
JUnit evidence, and portable receipts — works today. Container isolation and
signed attestations are [future work](/community/roadmap) and are not
implemented.
:::

## What it is not

AgentProof does not prove correctness, security or the absence of bugs, and it
does not use an LLM to judge whether prose requirements are met. It verifies the
*declared evidence workflow* — not the honesty or completeness of your tests.
Review the policy, the tests and the source together. See the
[security model](/guide/security).
