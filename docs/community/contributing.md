# Contributing

AgentProof is an open-source developer preview and contributions are welcome.
The authoritative guide is
[`CONTRIBUTING.md`](https://github.com/nowshad7/agent-proof/blob/main/CONTRIBUTING.md)
in the repository; this page is a summary.

## Develop

```sh
git clone https://github.com/nowshad7/agent-proof
cd agent-proof

go build ./...
go vet ./...
go test ./...

# Run the bundled examples end to end.
./examples/run-examples.sh
```

## What makes a good contribution

- **Fail-closed first.** Any adapter or aggregation change must keep missing,
  truncated or ambiguous evidence from passing. Add tests that prove it.
- **Standard library by default.** The core avoids dependencies; prefer the
  standard library unless there is a strong reason.
- **Honest scope.** Do not add placeholder integrations or claims for
  capabilities that are not implemented and tested.

## Adapter requests

A new framework adapter needs a **real emitted report fixture** and the exact
test identities it should produce. An emitter is supported only once its real
output is covered by fixtures — see the [JUnit guide](/guides/junit).

## Filing issues

- **Bugs:** include your tool version (`agentproof version`), OS, a minimal
  contract and sanitized evidence.
- **Adapters:** include a real report fixture and expected identities.

Issues scoped small enough may be labeled `good first issue` or `help wanted`.

## Docs

The documentation is VitePress. Build it locally with:

```sh
npm ci
npm run docs:dev      # local preview
npm run docs:build    # production build
```

See the [roadmap](/community/roadmap) for where help is most useful.
