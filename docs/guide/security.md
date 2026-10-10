# Security model

AgentProof is honest about what it does and does not protect. Read this before
verifying changes you do not trust.

## Local mode is not a sandbox

Trusted local execution runs your checks with **your** privileges. A temporary
`HOME` and an exported checkout reduce accidental mess and some secret leakage,
but they do **not** prevent code from:

- reading other files your user can read,
- making network requests,
- consuming CPU, memory and disk.

For untrusted changes — including agent-authored changes you have not reviewed —
run AgentProof inside a **disposable, unprivileged environment** that you or your
CI owner control. Do not expose privileged tokens, production credentials,
Docker sockets or valuable host mounts to it.

## Threats and boundaries

| Threat | Mitigation or documented boundary |
| --- | --- |
| Candidate weakens its own policy | Apply a policy from a protected base revision with `--config`. |
| Candidate changes acceptance tests | Keep acceptance tests outside the candidate's control. |
| Stale test report "passes" | The expected JUnit report is removed before the check runs. |
| Missing or truncated evidence | Fail closed with `ERROR` or `UNVERIFIED`. |
| Archive escape | Paths and entry types are validated before extraction. |
| HTML injection in the receipt | All evidence is escaped; a restrictive CSP is applied. |
| Secret leakage | Minimal environment forwarding; review artifacts before sharing. |
| Forged receipt | Receipts are unsigned and editable; treat them accordingly. |
| Runaway process | Timeouts and process-group cancellation (weaker on Windows). |
| Dangerous CI permissions | Separate verification and deployment; minimize permissions. |

## Protected policy and tests {#protected-policy}

If a candidate change can edit the contract or the acceptance tests, it can make
the gate pass without doing the work. To raise the bar:

- **Policy:** keep `agent-proof.json` on a protected base branch and apply it
  with `--config`. The receipt records the policy's SHA-256 so you can confirm
  which bytes ran.
- **Tests:** keep acceptance tests where the candidate change does not modify
  them, and point requirements at those tests.

See the [CI guide](/guide/ci) for a concrete GitHub Actions setup.

## Provenance and signatures {#provenance}

A receipt proves only that the declared checks were run and what they reported.
It is **unsigned and editable**:

- A plain content hash, if present, is for integrity, not authentication.
  AgentProof never describes a hash as a signature or an independent attestation.
- Signed attestations and identity verification are [future work](/community/roadmap),
  not implemented in this preview.

## Evidence honesty

AgentProof verifies the *declared evidence workflow*. It cannot certify that your
tests have adequate assertions, or that a command did not fabricate its output.
A weak test still produces a passing receipt. **Review the policy, the tests and
the source together** — the receipt tells a reviewer where to look, not that
looking is unnecessary.

## Reporting a vulnerability

See [`SECURITY.md`](https://github.com/nowshad7/agent-proof/blob/main/SECURITY.md)
in the repository for how to report security issues privately.
