# Trusted local execution

AgentProof runs your checks against a snapshot of a committed commit. This page
describes exactly how, and what the boundaries are.

::: warning Not a sandbox
Local execution is **not** isolation. A temporary `HOME` and an exported
checkout do not stop code from reading other files or using the network with
your privileges. For untrusted changes, run inside a disposable, unprivileged
environment. See the [security model](/guide/security).
:::

## The flow

1. Resolve the revision (default `HEAD`) to a full commit SHA.
2. Export the committed tree with `git archive` and extract it into a fresh
   temporary directory. The working checkout is never touched.
3. Load and validate the contract.
4. For each check, run its command in the snapshot directory with a timeout.
5. Parse evidence, aggregate requirements, write receipts.
6. Remove the temporary snapshot.

## Snapshot safety

The extractor applies strict checks to the archive:

- Rejects symlinks, hard links, path traversal and special files.
- Accepts safe tar metadata records, such as git's `pax_global_header`, without
  treating them as files.
- Limits extraction to 256 MiB and 50,000 entries.

## What the snapshot does and does not contain

- Dirty and untracked working-tree files are **excluded** — only committed
  content is verified.
- The snapshot does **not** contain `.git`.
- `git archive` attributes (for example `export-ignore`) affect exported
  content.
- Submodule content and Git LFS hydration are **not** supported in this preview.
- An exact input commit does not make execution bit-for-bit reproducible:
  runtime tools, network services and unpinned dependencies can still vary.

## Process environment

Each check runs with a deliberately minimal environment:

- A fresh `HOME`, `TMPDIR`, `GOCACHE`, `GOPATH` and `GOMODCACHE`.
- `PATH` is inherited so installed toolchains remain discoverable.
- Common proxy and TLS variables (`HTTPS_PROXY`, `SSL_CERT_FILE`, …) are
  forwarded so module downloads can work behind a proxy.
- The rest of your environment is **not** forwarded, which reduces accidental
  secret leakage into check processes.

## Limits and timeouts

- `stdout` and `stderr` are each bounded to 1 MiB; overflow is an `ERROR`.
- JUnit report input is bounded to 8 MiB.
- On timeout, AgentProof cancels the child's process group on Unix. Descendant
  handling on Windows is weaker and documented as a limitation.
- Checks run **sequentially** in the same exported snapshot.
