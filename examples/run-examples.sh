#!/usr/bin/env bash
#
# Build agentproof and verify each bundled example.
#
# The examples live inside this repository, so they are not their own git
# repositories. AgentProof snapshots a committed git tree, so this script
# copies each example into a fresh temporary git repository, commits it, and
# runs `agentproof verify` against that. This is also how a real user would
# run AgentProof: from the root of their own project's repository.
#
# Usage: examples/run-examples.sh [example-name ...]
#   With no arguments, every example is verified.

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
examples_dir="$repo_root/examples"

bin="$(mktemp -d)/agentproof"
echo ">> building agentproof"
(cd "$repo_root" && go build -o "$bin" ./cmd/agentproof)

examples=("$@")
if [ ${#examples[@]} -eq 0 ]; then
  examples=(export junit)
fi

overall=0
for name in "${examples[@]}"; do
  src="$examples_dir/$name"
  if [ ! -d "$src" ]; then
    echo "!! no such example: $name" >&2
    overall=1
    continue
  fi

  echo
  echo ">> verifying example: $name"
  work="$(mktemp -d)"
  cp -R "$src/." "$work/"
  (
    cd "$work"
    git init -q
    git config user.email "examples@agent-proof.local"
    git config user.name "AgentProof examples"
    git add -A
    git commit -q -m "example: $name"
  )

  out="$work/.agent-proof"
  if "$bin" verify --trust-local --repo "$work" --out "$out" --runner "examples/$name"; then
    echo ">> $name: PASS"
  else
    code=$?
    echo "!! $name: agentproof exited $code" >&2
    overall=$code
  fi
done

echo
if [ "$overall" -eq 0 ]; then
  echo ">> all examples passed"
else
  echo "!! one or more examples did not pass" >&2
fi
exit "$overall"
