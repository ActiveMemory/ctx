# Domain Decisions Schema

Shape for `.context/kb/domain-decisions.md`. Kb-scoped
decisions about the subject matter under study, distinct
from project-level `.context/DECISIONS.md`, which records
decisions about how this project is built.

The two files have different schemas, different write
authority, and different lifecycles. `domain-decisions.md`
is written by `/ctx-kb-ingest`; `DECISIONS.md` is written by
`/ctx-decision-add`.

## Fields

| Field           | Description                                                       |
|-----------------|-------------------------------------------------------------------|
| `id`            | Stable `DD-###` identifier; zero-padded; never renumbered.        |
| `date`          | ISO date the decision was recorded in the kb.                     |
| `context`       | What in the domain prompted the decision; observable facts only.  |
| `decision`      | The position taken, in one sentence.                              |
| `rationale`     | Why this position over the alternatives that were on the table.   |
| `consequence`   | What now changes for topic pages, glossary, or downstream claims. |
| `supporting-ev` | Comma-separated `EV-###` references that ground the decision.     |

## Example

```markdown
## DD-004: Treat Widget Bundles as Opaque to the Consumer

- **Date:** 2026-05-16.
- **Context.** Two upstream sources disagreed on whether a
  widget consumer should inspect bundle internals (EV-042,
  EV-051); the disagreement was resolved in `C-007`.
- **Decision.** The kb treats widget bundles as opaque from
  the consumer's perspective; inspection is an implementation
  concern of the producer.
- **Rationale.** Opacity preserves the producer's freedom to
  evolve internals without breaking downstream call sites.
- **Consequence.** The `widget-composition` topic page is
  rewritten to lead with opacity; the glossary's `widget`
  entry gains an `opacity` cross-reference.
- **Supporting EV:** EV-042, EV-043.
```

## Rules

- Do not confuse this with `.context/DECISIONS.md`; that file
  records project-build decisions and is owned by a different
  CLI surface.
- A domain decision is a kb-scoped position on the subject
  matter; it cites `EV-###` rows, not commit SHAs.
- Every entry cites at least one supporting `EV-###`; a
  domain decision without evidence belongs in
  `outstanding-questions.md` instead.
- Decisions are append-only and never renumbered; supersede
  in place by recording a new `DD-###` that cites the prior.
- Do not write to `.context/DECISIONS.md` from
  `/ctx-kb-ingest`; the authority boundary is strict.
