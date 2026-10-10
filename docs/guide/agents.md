# Agent workflows

AgentProof is agent-independent. Any agent that can run a command-line tool can
produce a receipt, and any human can read it. There is no provider-specific
integration required.

## Give the agent a task it can prove

A good agent task names the behavior *and* the evidence that will prove it:

> Implement CSV export that preserves Unicode text. Add a test
> `TestUnicode` in the `export` package that asserts non-ASCII content round
> trips. Then add a requirement `unicode-export` to `agent-proof.json` pointing
> at `example.com/project::TestUnicode`, run `agentproof verify --trust-local`,
> and include the receipt in your summary.

Because the contract names the exact test, the agent cannot claim success
without that test existing and passing.

## A typical loop

1. The agent edits code and tests.
2. The agent commits (AgentProof verifies committed state, so uncommitted work
   is not counted).
3. The agent runs `agentproof verify --trust-local`.
4. If the gate is not `PASS`, the agent reads the per-requirement detail in the
   receipt and iterates.
5. The agent attaches `receipt.md` (or a link to `receipt.html`) to its summary.

## Keep the human in the loop

AgentProof narrows what a reviewer must check; it does not remove the reviewer.
The receipt tells you which requirement was checked and how, but it cannot tell
you whether the test's assertions are meaningful. Review the policy, the tests
and the source together. A test with weak assertions still produces a passing
receipt.

## Keep acceptance outside the candidate's control

If the agent can edit the acceptance contract or the acceptance tests, it can
make the gate pass without doing the work. For higher-trust reviews:

- Keep the contract on a protected base branch and apply it with `--config`.
- Keep acceptance tests in a location the candidate change does not modify.

See [keeping policy and tests protected](/guide/security#protected-policy) and
the [CI guide](/guide/ci).
