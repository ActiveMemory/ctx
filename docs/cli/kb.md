---
title: ctx kb
icon: lucide/library
---

![ctx](../images/ctx-banner.png)

## `ctx kb`

Knowledge-base editorial pipeline (Phase KB). Manages the
`.context/kb/` knowledge base via mode-aware skills and a small
set of supporting CLI commands. The editorial constitution
lives at `.context/ingest/KB-RULES.md` (laid down by
`ctx init`).

```bash
ctx kb [subcommand]
```

| Subcommand                       | Type             | Purpose                                                             |
|----------------------------------|------------------|---------------------------------------------------------------------|
| `ctx kb topic new "<name>"`      | CLI (real)       | Sole writer of topic-page scaffolds. Creates `.context/kb/topics/<slug>/index.md` from the embedded template. Refuses when the topic exists. |
| `ctx kb note "<text>"`           | CLI (real)       | Appends a one-liner to `.context/ingest/findings.md`. Never touches a topic page.   |
| `ctx kb reindex`                 | CLI (real)       | Refreshes the `CTX:KB:TOPICS` managed block in `.context/kb/index.md`.              |
| `ctx kb ingest <folder\|paths>`  | Skill-driven     | Mode-aware editorial pass. CLI form refuses on empty input and points at the `/ctx-kb-ingest` skill. |
| `ctx kb ask "<question>"`        | Skill-driven     | Q&A grounded in the kb. CLI form refuses on empty input and points at the `/ctx-kb-ask` skill.  |
| `ctx kb site-review`             | Skill-driven     | Mechanical structural audit. Points at `/ctx-kb-site-review`.                       |
| `ctx kb ground`                  | Skill-driven     | External grounding via `grounding-sources.md`. Refuses when the file is empty.      |

!!! note "Skill-driven vs real CLI"
    The mode skills (`ingest`, `ask`, `site-review`, `ground`)
    do the editorial work themselves: the agent reads
    `.context/ingest/30-INGEST.md` (etc.) and executes the
    pass per the pass-mode contract. The CLI form for those
    subcommands validates input and prints the canonical skill
    invocation. The real CLI commands (`topic new`, `note`,
    `reindex`) own concrete state changes.

### `ctx kb topic new "<name>"`

Scaffolds a folder-shaped topic at `.context/kb/topics/<slug>/index.md`
from the embedded template.

**Slug**: lowercase + kebab-case. Slashes are preserved for
vendor-namespaced topology (e.g. `cursor/hooks`,
`cursor/skills`, `cursor/rules` under a shared `cursor/`
folder).

**Refuses** when the topic folder already exists. Use the
existing folder instead; the editorial pass extends pages,
it doesn't reset them.

### `ctx kb note "<text>"`

Appends a timestamped one-liner to
`.context/ingest/findings.md`. Use for parking findings the
next ingest pass should absorb.

```bash
ctx kb note "follow-up: chase the v1.2 release notes for the SIGTERM change"
```

### `ctx kb reindex`

Refreshes the `CTX:KB:TOPICS` managed block inside
`.context/kb/index.md` so the kb landing page enumerates
current topic folders. Run after `ctx kb topic new` to update
the landing.

### Skill-Driven Subcommands

`ingest`, `ask`, `site-review`, `ground` exist as CLI surfaces
so the editorial workflow is **drivable from outside Claude
Code** (via the fallback `PROMPT.md` auto-router). In Claude
Code, prefer the skills:

```text
/ctx-kb-ingest ./inputs/2026-05-15-call.md "cursor hooks"
/ctx-kb-ask "does the kb say hooks fire async?"
/ctx-kb-site-review
/ctx-kb-ground
```

See the
[Build a Knowledge Base recipe](../recipes/build-a-knowledge-base.md)
for the full workflow.

## Reference

- Recipe: [Build a Knowledge Base](../recipes/build-a-knowledge-base.md)
- Recipe: [Typical KB Session](../recipes/typical-kb-session.md)
- Editorial constitution: `.context/ingest/KB-RULES.md`
