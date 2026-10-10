# Verification badge

A badge in your README tells everyone who lands on your repository that its
agent-written changes are verified against an explicit acceptance contract — and
links them straight to the receipt. Every `agentproof verify` writes a
`badge.json` next to the other receipt files, ready to publish.

::: info Honest by construction
The badge reflects the receipt and **fails closed**. Only an explicit `PASS`
is ever coloured green. `FAIL` and `ERROR` are red, `UNVERIFIED` is orange,
`SKIPPED` is grey, and any unrecognised or empty status reads `unknown` in
grey. A broken, truncated or tampered receipt can never show as verified.
:::

## What gets written

```sh
agentproof verify --trust-local
# .agent-proof/badge.json  → {"schemaVersion":1,"label":"agent-proof","message":"4/4 passed","color":"brightgreen"}
```

You can also render it from any stored receipt without re-running the checks:

```sh
agentproof report --input .agent-proof/receipt.json --format badge
```

The document follows the [Shields.io endpoint](https://shields.io/badges/endpoint-badge)
schema, so Shields renders it for you from a URL.

## Publish it, then point a badge at it

The badge needs `badge.json` at a stable public URL. Pick whichever fits your
project; all three work the same way.

### Option A — a build artifact committed to a branch

A small CI job writes the badge to an orphan `badges` branch on every push to
your default branch:

```yaml
name: agent-proof-badge
on:
  push:
    branches: [main]

permissions:
  contents: write

jobs:
  badge:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: nowshad7/agent-proof@v0.1.0
        with:
          version: v0.1.0
          comment-on-pr: 'false'
      - name: Publish badge
        run: |
          git fetch origin badges || true
          git switch badges 2>/dev/null || git switch --orphan badges
          cp .agent-proof/badge.json badge.json
          git add badge.json
          git -c user.name=agent-proof -c user.email=noreply@github.com \
            commit -m "update agent-proof badge" || exit 0
          git push origin badges
```

Then in your README:

```md
![AgentProof](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/OWNER/REPO/badges/badge.json)
```

### Option B — GitHub Pages

If you already deploy a site, copy `badge.json` into the published output (for
a VitePress or static site, anything under `public/` is served at the site
root). The badge URL is then:

```md
![AgentProof](https://img.shields.io/endpoint?url=https://OWNER.github.io/REPO/badge.json)
```

### Option C — a gist

Write `badge.json` to a public gist from CI (for example with the `gh` CLI) and
point the endpoint at its raw URL. Useful when you do not want an extra branch.

## Link the badge to the receipt

Make the badge clickable so a reader goes straight to the evidence — the latest
CI run, the uploaded receipt artifact, or a published `receipt.html`:

```md
[![AgentProof](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/OWNER/REPO/badges/badge.json)](https://github.com/OWNER/REPO/actions/workflows/acceptance.yml)
```

## A static fallback

If you do not want to publish a live endpoint, a static badge is honest as long
as your CI gate actually blocks merges on a failed receipt (so the badge cannot
drift from reality):

```md
[![AgentProof](https://img.shields.io/badge/agent--proof-verified-brightgreen)](https://github.com/OWNER/REPO/actions)
```

A static badge does not update on its own — only use one when a required status
check keeps `main` from going red underneath it.

## See also

- [GitHub Actions / CI](/guide/ci) — the Action that produces the receipt and
  `badge.json` in CI.
- [Reports](/guide/reports) — the other receipt formats.
- [Security model](/guide/security) — what a receipt does and does not prove.
