# Contributing to AgentProof

Thanks for your interest. AgentProof is an open-source developer preview, and
contributions — bug reports, adapters, docs, tests — are welcome.

## Principles

1. **Fail closed.** Missing, truncated, duplicate or malformed evidence must
   never pass. Any change to an adapter or to aggregation must preserve this and
   add tests that prove it.
2. **Standard library by default.** The core avoids third-party dependencies.
   Prefer the standard library unless there is a strong, stated reason.
3. **Honest scope.** Do not add placeholder integrations, badges, or claims for
   capabilities that are not implemented and tested.

## Development

```sh
git clone https://github.com/nowshad7/agent-proof
cd agent-proof

go build ./...
go vet ./...
gofmt -l .          # should print nothing
go test ./...        # or: go test -race ./...

# Run the bundled examples end to end.
./examples/run-examples.sh
```

The repository verifies itself: `agent-proof.json` at the root points at the
core regression tests. You can run it with:

```sh
go build -o /tmp/agentproof ./cmd/agentproof
/tmp/agentproof verify --trust-local
```

## Documentation

```sh
npm ci
npm run docs:dev      # local preview at http://localhost:5173
npm run docs:build    # production build (must pass; dead links fail the build)
```

## Adding an adapter

An adapter is considered supported only once its **real emitted output** is
covered by fixtures. To propose one:

1. Open an [adapter request](https://github.com/nowshad7/agent-proof/issues/new?template=adapter_request.yml)
   with a genuine, sanitized report and the exact test identities expected.
2. In the PR, add the parser under `internal/proof`, fixtures drawn from real
   output, and tests for the pass, fail, skip, missing and malformed cases.

## Pull requests

- Keep changes focused and tested.
- Update docs when behavior or interfaces change.
- Fill in the PR template, including how you verified the change.
- By contributing, you agree your contribution is licensed under
  [Apache-2.0](LICENSE).

## Reporting bugs

Use the [bug report template](https://github.com/nowshad7/agent-proof/issues/new?template=bug_report.yml).
Include your `agentproof version`, OS, a minimal contract and sanitized evidence.
Never paste secrets.

## Code of conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). By
participating you are expected to uphold it.
