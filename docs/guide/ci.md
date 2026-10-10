# GitHub Actions / CI

AgentProof's exit codes make it a natural CI gate. This page shows a minimal
workflow and the trust considerations that matter in CI.

## A minimal gate

```yaml
name: acceptance
on: [push, pull_request]

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0 # AgentProof snapshots a commit, so history is needed

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Build agentproof
        run: go build -o /usr/local/bin/agentproof ./cmd/agentproof

      - name: Verify
        run: agentproof verify --trust-local --runner "ci/${{ github.sha }}"

      - name: Upload receipt
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: agent-proof-receipt
          path: .agent-proof/
```

The job fails when `agentproof` returns a nonzero exit code. The receipt is
uploaded on every run (`if: always()`) so you can inspect failures.

## Trust in CI

::: warning
In CI, your checks run with whatever the job can access. Treat an
agent-authored change as untrusted:

- Do **not** run untrusted PR code under `pull_request_target` with privileged
  credentials or secrets.
- Keep deployment permissions out of the verification job.
- Prefer a disposable, unprivileged runner with no access to production
  credentials, Docker sockets or valuable mounts.
:::

## Protecting the policy

So a candidate PR cannot weaken its own gate, keep the contract on a protected
base branch and apply it explicitly:

```yaml
- name: Fetch reviewed policy from base
  run: git show origin/main:agent-proof.json > /tmp/policy.json
- name: Verify against protected policy
  run: agentproof verify --trust-local --config /tmp/policy.json
```

The receipt records the SHA-256 of the exact policy bytes applied, so reviewers
can confirm which policy ran.

## The bundled Action

A reusable GitHub Action (`action.yml`) runs the verification and, on pull
requests, posts the receipt as a single sticky comment that it updates in place
on every run — including when the gate fails, so reviewers always see why.

```yaml
name: acceptance
on: [pull_request]

permissions:
  contents: read
  pull-requests: write # required so the Action can post its receipt comment

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: nowshad7/agent-proof@v0.1.0
        with:
          version: v0.1.0 # a release tag downloads a prebuilt binary
```

Key inputs:

| Input | Default | Meaning |
| --- | --- | --- |
| `version` | `main` | A release tag (e.g. `v0.1.0`) downloads a prebuilt binary; a branch or SHA builds from source. |
| `config` | — | Apply a separately reviewed policy (for example from a protected base). |
| `comment-on-pr` | `true` | Post/update the receipt comment on `pull_request` events. |
| `commit`, `repo`, `out` | `HEAD`, `.`, `.agent-proof` | As for the CLI. |

Outputs `status` and `receipt` are available for later steps. The Action still
fails the job with the CLI's exit code, so it gates the PR.

::: warning
The Action is a developer preview. It runs the target repository's checks with
the job's privileges — keep the trust notes above in mind, and do not run it
under `pull_request_target` with secrets.
:::
