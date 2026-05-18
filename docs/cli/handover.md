---
title: ctx handover
icon: lucide/scroll-text
---

![ctx](../images/ctx-banner.png)

## `ctx handover`

Writes the per-session handover under
`.context/handovers/<TS>-<slug>.md`: a former-agent-to-next-agent
note created at session end by `/ctx-wrap-up` and read at
session start by `/ctx-remember`. When `.context/kb/` exists,
the writer additionally folds postdated closeouts into the
handover's `## Folded Closeouts` section and archives them.

### `ctx handover write <title>`

```bash
ctx handover write "Cursor Hooks deep dive" \
  --summary "Drafted topic-page; minted EV-018..EV-024; cold-reader passed." \
  --next "Re-ingest the v1.1 release notes URL once you have it."
```

**Required flags**:

| Flag        | Description                                                                                              |
|-------------|----------------------------------------------------------------------------------------------------------|
| `--summary` | What happened this session (past tense). Placeholder values (`TBD`, `see chat`, `n/a`) are rejected.     |
| `--next`    | What the next agent should do FIRST (future tense, specific). Same placeholder rejection.                |

**Optional flags**:

| Flag               | Description                                                                                       |
|--------------------|---------------------------------------------------------------------------------------------------|
| `--highlights`     | Notable artifacts produced this session.                                                          |
| `--open-questions` | Things that remain undecided.                                                                     |
| `--commit`         | Override resolved git HEAD for the Provenance line (CI replay; honors `CTX_TASK_COMMIT`).         |
| `--no-fold`        | Skip closeout consumption (mid-session checkpoint).                                               |

**Writes**: `.context/handovers/<TS>-<slug>.md` with frontmatter
(`sha`, `branch`, `generated-at`, `title`) and body sections
(`## Summary`, `## Next Session`, optionally `## Highlights`,
`## Open Questions`, `## Folded Closeouts`). The
`<TS>-<slug>.md` filename is timestamped so multiple concurrent
agent runs never overwrite one another's handover.

**Side effect** (when `--no-fold` is absent and `.context/kb/`
exists): closeouts that postdate the latest handover are
folded into the new handover and **physically archived** under
`.context/archive/closeouts/`.

### How to Trigger

In ordinary sessions you do not invoke `ctx handover write`
directly. The user-facing trigger is `/ctx-wrap-up`:

```text
/ctx-wrap-up "session title"
```

`/ctx-wrap-up` owns session-end and always delegates to
`/ctx-handover` as its final step. Direct invocation of
`/ctx-handover` is reserved for two cases:

- `--no-fold` mid-session checkpoint.
- Recovery, when a prior session aborted before wrap-up.

See [`/ctx-wrap-up`](../reference/skills.md#ctx-wrap-up) and
[`/ctx-handover`](../reference/skills.md#ctx-handover).

## Reference

- Recipe: [Session Lifecycle](../recipes/session-lifecycle.md)
- Recipe: [Recover an Aborted Session](../recipes/recover-aborted-session.md)
- Skill: [`/ctx-wrap-up`](../reference/skills.md#ctx-wrap-up)
- Skill: [`/ctx-handover`](../reference/skills.md#ctx-handover)
- Skill: [`/ctx-remember`](../reference/skills.md#ctx-remember)
