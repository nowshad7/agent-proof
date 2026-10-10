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

## About the bundled Action

A reusable GitHub Action (`action.yml`) ships in the repository, but it is
early. Test it from a separate sample repository before relying on it, and read
its inputs and trust notes first. Treat the Action as preview-quality until the
[roadmap](/community/roadmap) marks it validated.
