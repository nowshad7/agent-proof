# Security policy

## Reporting a vulnerability

Please report security issues **privately**. Do not open a public issue for a
vulnerability.

- Preferred: open a private advisory at
  <https://github.com/nowshad7/agent-proof/security/advisories/new>.

Please include a description, affected version (`agentproof version`), and steps
to reproduce. We will acknowledge your report and keep you updated as we
investigate. This is a volunteer-maintained preview, so we do not promise a fixed
response-time SLA, but we take security reports seriously.

## Scope and threat model

Read the full [security model](https://nowshad7.github.io/agent-proof/guide/security)
before relying on AgentProof. The essentials:

- **Trusted local execution is not a sandbox.** Checks run with the invoking
  user's privileges and can read files and use the network. For untrusted
  changes, run inside a disposable, unprivileged environment.
- **Receipts are unsigned and editable.** They prove only that the declared
  checks were run and what they reported. They are not authenticated
  attestations, and any content hash is for integrity, not authentication.
- **Evidence honesty is out of scope.** AgentProof verifies the declared
  evidence workflow, not whether your tests have adequate assertions.

## What counts as a vulnerability here

Examples of in-scope issues:

- Snapshot extraction that writes outside the snapshot root, follows a symlink,
  or otherwise escapes containment.
- Evidence that should fail closed but is accepted as a pass (for example a
  stale report, truncated stream, or missing named test reported as `PASS`).
- HTML receipt output that is not correctly escaped, enabling injection.
- The CLI forwarding more of the environment to checks than documented.

Out of scope: the documented boundaries above (local execution not being a
sandbox, receipts being unsigned) are known limitations, not vulnerabilities.
